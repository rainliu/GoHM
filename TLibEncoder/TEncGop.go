/* The copyright in this software is being made available under the BSD
 * License, included below. This software may be subject to other third party
 * and contributor rights, including patent rights, and no such rights are
 * granted under this license.
 *
 * Copyright (c) 2012-2013, H265.net
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 *  * Redistributions of source code must retain the above copyright notice,
 *    this list of conditions and the following disclaimer.
 *  * Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *  * Neither the name of the H265.net nor the names of its contributors may
 *    be used to endorse or promote products derived from this software without
 *    specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS
 * BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF
 * THE POSSIBILITY OF SUCH DAMAGE.
 */

package TLibEncoder

import (
	"io"
    "container/list"
    "fmt"
    "gohm/TLibCommon"
    "math"
    "time"
)

// ====================================================================================================================
// Enumeration
// ====================================================================================================================
type PROCESSING_STATE uint8

const (
    EXECUTE_INLOOPFILTER = iota
    ENCODE_SLICE
)

type SCALING_LIST_PARAMETER uint8

const (
    SCALING_LIST_OFF = iota
    SCALING_LIST_DEFAULT
    SCALING_LIST_FILE_READ
)

// ====================================================================================================================
// Class definition
// ====================================================================================================================

/// GOP encoder class
type TEncGOP struct {
    //  Data
    m_bLongtermTestPictureHasBeenCoded  bool
    m_bLongtermTestPictureHasBeenCoded2 bool
    m_numLongTermRefPicSPS              uint
    m_ltRefPicPocLsbSps                 [33]uint
    m_ltRefPicUsedByCurrPicFlag         [33]uint
    m_iLastIDR                          int
    m_iGopSize                          int
    m_iNumPicCoded                      int
    m_bFirst                            bool

    //  Access channel
    m_pcEncTop       *TEncTop
    m_pcCfg          *TEncCfg
    m_pcSliceEncoder *TEncSlice
    m_pcListPic      *list.List

	m_pTraceFile	  io.Writer
    m_pcEntropyCoder *TEncEntropy
    m_pcCavlcCoder   *TEncCavlc
    m_pcSbacCoder    *TEncSbac
    m_pcBinCABAC     *TEncBinCABAC
    m_pcLoopFilter   *TLibCommon.TComLoopFilter

    // m_seiWriter		SEIWriter;

    //--Adaptive Loop filter
    m_pcSAO        *TEncSampleAdaptiveOffset
    m_pcBitCounter *TLibCommon.TComBitCounter
    m_pcRateCtrl   *TEncRateCtrl
    // indicate sequence first
    m_bSeqFirst bool

    // clean decoding refresh
    m_bRefreshPending                          bool
    m_pocCRA                                   int
    m_storedStartCUAddrForEncodingSlice        map[int]int //vector<int>
    m_storedStartCUAddrForEncodingSliceSegment map[int]int

    m_vRVM_RP         map[int]int
    m_lastBPSEI       uint
    m_totalCoded      uint
    m_cpbRemovalDelay uint
    m_tl0Idx          uint
    m_rapIdx          uint
    //#if L0045_NON_NESTED_SEI_RESTRICTIONS
    m_activeParameterSetSEIPresentInAU bool
    m_bufferingPeriodSEIPresentInAU    bool
    m_pictureTimingSEIPresentInAU      bool
    //#endif
    m_gcAnalyzeAll TEncAnalyze
    m_gcAnalyzeI   TEncAnalyze
    m_gcAnalyzeP   TEncAnalyze
    m_gcAnalyzeB   TEncAnalyze
}

func NewTEncGOP() *TEncGOP {
    return &TEncGOP{m_bFirst: true, m_bSeqFirst: true}
    //m_storedStartCUAddrForEncodingSlice:list.New(),
    //m_storedStartCUAddrForEncodingSliceSegment:list.New()};
}

func (this *TEncGOP) create()  {}
func (this *TEncGOP) destroy() {}

func (this *TEncGOP) init(pcTEncTop *TEncTop) {
    this.m_pcEncTop = pcTEncTop
    this.m_pcCfg = pcTEncTop.GetEncCfg()
    this.m_pcSliceEncoder = pcTEncTop.getSliceEncoder()
    this.m_pcListPic = pcTEncTop.getListPic()

	this.m_pTraceFile = pcTEncTop.getTraceFile()
    this.m_pcEntropyCoder = pcTEncTop.getEntropyCoder()
    this.m_pcCavlcCoder = pcTEncTop.getCavlcCoder()
    this.m_pcSbacCoder = pcTEncTop.getSbacCoder()
    this.m_pcBinCABAC = pcTEncTop.getBinCABAC()
    this.m_pcLoopFilter = pcTEncTop.getLoopFilter()
    this.m_pcBitCounter = pcTEncTop.getBitCounter()

    //--Adaptive Loop filter
    this.m_pcSAO = pcTEncTop.getSAO()
    this.m_pcRateCtrl = pcTEncTop.getRateCtrl()

    this.m_lastBPSEI = 0
    this.m_totalCoded = 0
    
    this.m_vRVM_RP = make(map[int]int)
}

func (this *TEncGOP) compressGOP(iPOCLast, iNumPicRcvd int, rcListPic, rcListPicYuvRecOut *list.List, accessUnitsInGOP *AccessUnits) {
    var pcPic *TLibCommon.TComPic
    var pcPicYuvRecOut *TLibCommon.TComPicYuv
    var pcSlice *TLibCommon.TComSlice
    var pcBitstreamRedirect *TLibCommon.TComOutputBitstream
    pcBitstreamRedirect = TLibCommon.NewTComOutputBitstream()
    //AccessUnit::iterator  itLocationToPushSliceHeaderNALU; // used to store location where NALU containing slice header is to be inserted
    uiOneBitstreamPerSliceLength := uint(0)
    var pcSbacCoders []*TEncSbac
    var pcSubstreamsOut []*TLibCommon.TComOutputBitstream

    this.xInitGOP(iPOCLast, iNumPicRcvd, rcListPic, rcListPicYuvRecOut)

    this.m_iNumPicCoded = 0
    //var pictureTimingSEI SEIPictureTiming;
    //#if L0044_DU_DPB_OUTPUT_DELAY_HRD
    //picSptDpbOutputDuDelay := 0;
    //#endif
    //var accumBitsDU,accumNalsDU  *uint;
    //var decodingUnitInfoSEI SEIDecodingUnitInfo;
    for iGOPid := 0; iGOPid < this.m_iGopSize; iGOPid++ {
        uiColDir := uint(1)
        //-- For time output for each slice
        iBeforeTime := time.Now()

        //select uiColDir
        iCloseLeft := 1
        iCloseRight := -1
        for i := 0; i < this.m_pcCfg.GetGOPEntry(iGOPid).m_numRefPics; i++ {
            iRef := this.m_pcCfg.GetGOPEntry(iGOPid).m_referencePics[i]
            if iRef > 0 && (iRef < iCloseRight || iCloseRight == -1) {
                iCloseRight = iRef
            } else if iRef < 0 && (iRef > iCloseLeft || iCloseLeft == 1) {
                iCloseLeft = iRef
            }
        }
        if iCloseRight > -1 {
            iCloseRight = iCloseRight + this.m_pcCfg.GetGOPEntry(iGOPid).m_POC - 1
        }
        if iCloseLeft < 1 {
            iCloseLeft = iCloseLeft + this.m_pcCfg.GetGOPEntry(iGOPid).m_POC - 1
            for iCloseLeft < 0 {
                iCloseLeft += this.m_iGopSize
            }
        }
        iLeftQP := 0
        iRightQP := 0
        for i := 0; i < this.m_iGopSize; i++ {
            if this.m_pcCfg.GetGOPEntry(i).m_POC == (iCloseLeft%this.m_iGopSize)+1 {
                iLeftQP = this.m_pcCfg.GetGOPEntry(i).m_QPOffset
            }
            if this.m_pcCfg.GetGOPEntry(i).m_POC == (iCloseRight%this.m_iGopSize)+1 {
                iRightQP = this.m_pcCfg.GetGOPEntry(i).m_QPOffset
            }
        }
        if iCloseRight > -1 && iRightQP < iLeftQP {
            uiColDir = 0
        }

        /////////////////////////////////////////////////////////////////////////////////////////////////// Initial to start encoding
        pocCurr := iPOCLast - iNumPicRcvd + this.m_pcCfg.GetGOPEntry(iGOPid).m_POC
        iTimeOffset := this.m_pcCfg.GetGOPEntry(iGOPid).m_POC
        if iPOCLast == 0 {
            pocCurr = 0
            iTimeOffset = 1
        }
        if pocCurr >= this.m_pcCfg.GetFramesToBeEncoded() {
            continue
        }

        if this.getNalUnitType(pocCurr) == TLibCommon.NAL_UNIT_CODED_SLICE_IDR || this.getNalUnitType(pocCurr) == TLibCommon.NAL_UNIT_CODED_SLICE_IDR_N_LP {
            this.m_iLastIDR = pocCurr
        }
        // start a new access unit: create an entry in the list of output access units
        accessUnitsInGOP.PushBack(NewAccessUnit())
        accessUnit := accessUnitsInGOP.Back().Value.(*AccessUnit)
        pcPic, pcPicYuvRecOut = this.xGetBuffer(rcListPic, rcListPicYuvRecOut, iNumPicRcvd, iTimeOffset, pocCurr)
		
        //  Slice data initialization
        pcPic.ClearSliceBuffer()
        //assert(pcPic.GetNumAllocatedSlice() == 1);
        this.m_pcSliceEncoder.setSliceIdx(0)
        pcPic.SetCurrSliceIdx(0)
		
        pcSlice = this.m_pcSliceEncoder.initEncSlice(pcPic, iPOCLast, pocCurr, iNumPicRcvd, iGOPid, this.m_pcEncTop.getSPS(), this.m_pcEncTop.getPPS())
		//fmt.Printf("getSliceType1=%d\n", pcSlice.GetSliceType());
		
        pcSlice.SetLastIDR(this.m_iLastIDR)
        pcSlice.SetSliceIdx(0)

        //set default slice level flag to the same as SPS level flag
        pcSlice.SetLFCrossSliceBoundaryFlag(pcSlice.GetPPS().GetLoopFilterAcrossSlicesEnabledFlag())
        pcSlice.SetScalingList(this.m_pcEncTop.getScalingList())
        pcSlice.GetScalingList().SetUseTransformSkip(this.m_pcEncTop.getPPS().GetUseTransformSkip())
        if this.m_pcCfg.GetUseScalingListId() == SCALING_LIST_OFF {
            this.m_pcEncTop.getTrQuant().SetFlatScalingList()
            this.m_pcEncTop.getTrQuant().SetUseScalingList(false)
            this.m_pcEncTop.getSPS().SetScalingListPresentFlag(false)
            this.m_pcEncTop.getPPS().SetScalingListPresentFlag(false)
        } else if this.m_pcCfg.GetUseScalingListId() == SCALING_LIST_DEFAULT {
            pcSlice.SetDefaultScalingList()
            this.m_pcEncTop.getSPS().SetScalingListPresentFlag(false)
            this.m_pcEncTop.getPPS().SetScalingListPresentFlag(false)
            this.m_pcEncTop.getTrQuant().SetScalingList(pcSlice.GetScalingList())
            this.m_pcEncTop.getTrQuant().SetUseScalingList(true)
        } else if this.m_pcCfg.GetUseScalingListId() == SCALING_LIST_FILE_READ {
            if pcSlice.GetScalingList().XParseScalingList(this.m_pcCfg.GetScalingListFile()) {
                pcSlice.SetDefaultScalingList()
            }
            pcSlice.GetScalingList().CheckDcOfMatrix()
            this.m_pcEncTop.getSPS().SetScalingListPresentFlag(pcSlice.CheckDefaultScalingList())
            this.m_pcEncTop.getPPS().SetScalingListPresentFlag(false)
            this.m_pcEncTop.getTrQuant().SetScalingList(pcSlice.GetScalingList())
            this.m_pcEncTop.getTrQuant().SetUseScalingList(true)
        } else {
            fmt.Printf("error : ScalingList == %d no support\n", this.m_pcCfg.GetUseScalingListId())
            return //assert(0);
        }

        if pcSlice.GetSliceType() == TLibCommon.B_SLICE && this.m_pcCfg.GetGOPEntry(iGOPid).m_sliceType == "P" {
            pcSlice.SetSliceType(TLibCommon.P_SLICE)
        }
        //fmt.Printf("getSliceType2=%d\n", pcSlice.GetSliceType());
        
        // Set the nal unit type
        pcSlice.SetNalUnitType(this.getNalUnitType(pocCurr))
        if pcSlice.GetNalUnitType() == TLibCommon.NAL_UNIT_CODED_SLICE_TRAIL_R {
            if pcSlice.GetTemporalLayerNonReferenceFlag() {
                pcSlice.SetNalUnitType(TLibCommon.NAL_UNIT_CODED_SLICE_TRAIL_N)
            }
        }

        // Do decoding refresh marking if any
        pcSlice.DecodingRefreshMarking(&this.m_pocCRA, &this.m_bRefreshPending, rcListPic)
        this.m_pcEncTop.selectReferencePictureSet(pcSlice, pocCurr, iGOPid)
        pcSlice.GetRPS().SetNumberOfLongtermPictures(0)

        if pcSlice.CheckThatAllRefPicsAreAvailable(rcListPic, pcSlice.GetRPS(), false, 0) != 0 {
            pcSlice.CreateExplicitReferencePictureSetFromReference(rcListPic, pcSlice.GetRPS())
        }
        pcSlice.ApplyReferencePictureSet(rcListPic, pcSlice.GetRPS())

        if pcSlice.GetTLayer() > 0 {
            if pcSlice.IsTemporalLayerSwitchingPoint(rcListPic) || pcSlice.GetSPS().GetTemporalIdNestingFlag() {
                if pcSlice.GetTemporalLayerNonReferenceFlag() {
                    pcSlice.SetNalUnitType(TLibCommon.NAL_UNIT_CODED_SLICE_TSA_N)
                } else {
                    pcSlice.SetNalUnitType(TLibCommon.NAL_UNIT_CODED_SLICE_TLA)
                }
            } else if pcSlice.IsStepwiseTemporalLayerSwitchingPointCandidate(rcListPic) {
                isSTSA := true
                for ii := iGOPid + 1; ii < this.m_pcCfg.GetGOPSize() && isSTSA == true; ii++ {
                    lTid := this.m_pcCfg.GetGOPEntry(ii).m_temporalId
                    if lTid == int(pcSlice.GetTLayer()) {
                        nRPS := pcSlice.GetSPS().GetRPSList().GetReferencePictureSet(ii)
                        for jj := 0; jj < nRPS.GetNumberOfPictures(); jj++ {
                            if nRPS.GetUsed(jj) {
                                tPoc := this.m_pcCfg.GetGOPEntry(ii).m_POC + nRPS.GetDeltaPOC(jj)
                                kk := 0
                                for kk = 0; kk < this.m_pcCfg.GetGOPSize(); kk++ {
                                    if this.m_pcCfg.GetGOPEntry(kk).m_POC == tPoc {
                                        break
                                    }
                                }
                                tTid := this.m_pcCfg.GetGOPEntry(kk).m_temporalId
                                if tTid >= int(pcSlice.GetTLayer()) {
                                    isSTSA = false
                                    break
                                }
                            }
                        }
                    }
                }
                if isSTSA == true {
                    if pcSlice.GetTemporalLayerNonReferenceFlag() {
                        pcSlice.SetNalUnitType(TLibCommon.NAL_UNIT_CODED_SLICE_STSA_N)
                    } else {
                        pcSlice.SetNalUnitType(TLibCommon.NAL_UNIT_CODED_SLICE_STSA_R)
                    }
                }
            }
        }
        this.arrangeLongtermPicturesInRPS(pcSlice, rcListPic)
        refPicListModification := pcSlice.GetRefPicListModification()
        refPicListModification.SetRefPicListModificationFlagL0(false)
        refPicListModification.SetRefPicListModificationFlagL1(false)
        pcSlice.SetNumRefIdx(TLibCommon.REF_PIC_LIST_0, TLibCommon.MIN(this.m_pcCfg.GetGOPEntry(iGOPid).m_numRefPicsActive, pcSlice.GetRPS().GetNumberOfPictures()).(int))
        pcSlice.SetNumRefIdx(TLibCommon.REF_PIC_LIST_1, TLibCommon.MIN(this.m_pcCfg.GetGOPEntry(iGOPid).m_numRefPicsActive, pcSlice.GetRPS().GetNumberOfPictures()).(int))

        //#if ADAPTIVE_QP_SELECTION
        pcSlice.SetTrQuant(this.m_pcEncTop.getTrQuant())
        //#endif

        //  Set reference list
        pcSlice.SetRefPicList(rcListPic)

        //  Slice info. refinement
        if (pcSlice.GetSliceType() == TLibCommon.B_SLICE) && (pcSlice.GetNumRefIdx(TLibCommon.REF_PIC_LIST_1) == 0) {
            pcSlice.SetSliceType(TLibCommon.P_SLICE)
        }

        if pcSlice.GetSliceType() != TLibCommon.B_SLICE || !pcSlice.GetSPS().GetUseLComb() {
            pcSlice.SetNumRefIdx(TLibCommon.REF_PIC_LIST_C, 0)
            pcSlice.SetRefPicListCombinationFlag(false)
            pcSlice.SetRefPicListModificationFlagLC(false)
        } else {
            pcSlice.SetRefPicListCombinationFlag(pcSlice.GetSPS().GetUseLComb())
            pcSlice.SetNumRefIdx(TLibCommon.REF_PIC_LIST_C, pcSlice.GetNumRefIdx(TLibCommon.REF_PIC_LIST_0))
        }

        if pcSlice.GetSliceType() == TLibCommon.B_SLICE {
            pcSlice.SetColFromL0Flag(1 - uiColDir)
            bLowDelay := true
            iCurrPOC := pcSlice.GetPOC()
            iRefIdx := 0

            for iRefIdx = 0; iRefIdx < pcSlice.GetNumRefIdx(TLibCommon.REF_PIC_LIST_0) && bLowDelay; iRefIdx++ {
                if int(pcSlice.GetRefPic(TLibCommon.REF_PIC_LIST_0, iRefIdx).GetPOC()) > iCurrPOC {
                    bLowDelay = false
                }
            }
            for iRefIdx = 0; iRefIdx < pcSlice.GetNumRefIdx(TLibCommon.REF_PIC_LIST_1) && bLowDelay; iRefIdx++ {
                if int(pcSlice.GetRefPic(TLibCommon.REF_PIC_LIST_1, iRefIdx).GetPOC()) > iCurrPOC {
                    bLowDelay = false
                }
            }

            pcSlice.SetCheckLDC(bLowDelay)
        }

        uiColDir = 1 - uiColDir

        //-------------------------------------------------------------
        pcSlice.SetRefPOCList()

        pcSlice.SetNoBackPredFlag(false)
        if pcSlice.GetSliceType() == TLibCommon.B_SLICE && !pcSlice.GetRefPicListCombinationFlag() {
            if pcSlice.GetNumRefIdx(TLibCommon.RefPicList(0)) == pcSlice.GetNumRefIdx(TLibCommon.RefPicList(1)) {
                pcSlice.SetNoBackPredFlag(true)
                var i int
                for i = 0; i < pcSlice.GetNumRefIdx(TLibCommon.RefPicList(1)); i++ {
                    if pcSlice.GetRefPOC(TLibCommon.RefPicList(1), i) != pcSlice.GetRefPOC(TLibCommon.RefPicList(0), i) {
                        pcSlice.SetNoBackPredFlag(false)
                        break
                    }
                }
            }
        }

        if pcSlice.GetNoBackPredFlag() {
            pcSlice.SetNumRefIdx(TLibCommon.REF_PIC_LIST_C, 0)
        }
        pcSlice.GenerateCombinedList()

        if this.m_pcCfg.GetTMVPModeId() == 2 {
            if iGOPid == 0 { // first picture in SOP (i.e. forward B)
                pcSlice.SetEnableTMVPFlag(false)
            } else {
                // Note: pcSlice.GetColFromL0Flag() is assumed to be always 0 and getcolRefIdx() is always 0.
                pcSlice.SetEnableTMVPFlag(true)
            }
            pcSlice.GetSPS().SetTMVPFlagsPresent(true)
        } else if this.m_pcCfg.GetTMVPModeId() == 1 {
            pcSlice.GetSPS().SetTMVPFlagsPresent(true)
            pcSlice.SetEnableTMVPFlag(true)
        } else {
            pcSlice.GetSPS().SetTMVPFlagsPresent(false)
            pcSlice.SetEnableTMVPFlag(false)
        }
        /////////////////////////////////////////////////////////////////////////////////////////////////// Compress a slice
        //  Slice compression
        if this.m_pcCfg.GetUseASR() {
            this.m_pcSliceEncoder.setSearchRange(pcSlice)
        }

        bGPBcheck := false
        if pcSlice.GetSliceType() == TLibCommon.B_SLICE {
            if pcSlice.GetNumRefIdx(TLibCommon.RefPicList(0)) == pcSlice.GetNumRefIdx(TLibCommon.RefPicList(1)) {
                bGPBcheck = true
                var i int
                for i = 0; i < pcSlice.GetNumRefIdx(TLibCommon.RefPicList(1)); i++ {
                    if pcSlice.GetRefPOC(TLibCommon.RefPicList(1), i) != pcSlice.GetRefPOC(TLibCommon.RefPicList(0), i) {
                        bGPBcheck = false
                        break
                    }
                }
            }
        }
        if bGPBcheck {
            pcSlice.SetMvdL1ZeroFlag(true)
        } else {
            pcSlice.SetMvdL1ZeroFlag(false)
        }
        pcPic.GetSlice(pcSlice.GetSliceIdx()).SetMvdL1ZeroFlag(pcSlice.GetMvdL1ZeroFlag())

        //#if RATE_CONTROL_LAMBDA_DOMAIN
        sliceQP := pcSlice.GetSliceQp()
        lambda := float64(0.0)
        actualHeadBits := 0
        actualTotalBits := 0
        estimatedBits := 0
        tmpBitsBeforeWriting := 0
        if this.m_pcCfg.GetUseRateCtrl() {
            frameLevel := this.m_pcRateCtrl.getRCSeq().getGOPID2Level1(iGOPid)
            if pcPic.GetSlice(0).GetSliceType() == TLibCommon.I_SLICE {
                frameLevel = 0
            }
            this.m_pcRateCtrl.initRCPic(frameLevel)
            estimatedBits = this.m_pcRateCtrl.getRCPic().getTargetBits()

            if (pcSlice.GetPOC() == 0 && this.m_pcCfg.GetInitialQP() > 0) || (frameLevel == 0 && this.m_pcCfg.GetForceIntraQP()) { // QP is specified
                sliceQP = this.m_pcCfg.GetInitialQP()
                NumberBFrames := (this.m_pcCfg.GetGOPSize() - 1)
                dLambda_scale := 1.0 - TLibCommon.CLIP3(0.0, 0.5, 0.05*float64(NumberBFrames)).(float64)
                dQPFactor := 0.57 * dLambda_scale
                SHIFT_QP := 12
                bitdepth_luma_qp_scale := 0
                qp_temp := float64(sliceQP + bitdepth_luma_qp_scale - SHIFT_QP)
                lambda = dQPFactor * math.Pow(2.0, qp_temp/3.0)
            } else if frameLevel == 0 { // intra case, but use the model
                if this.m_pcCfg.GetIntraPeriod() != 1 { // do not refine allocated bits for all intra case
                    bits := this.m_pcRateCtrl.getRCSeq().getLeftAverageBits()
                    bits = this.m_pcRateCtrl.getRCSeq().getRefineBitsForIntra(bits)
                    if bits < 200 {
                        bits = 200
                    }
                    this.m_pcRateCtrl.getRCPic().setTargetBits(bits)
                }

                listPreviousPicture := this.m_pcRateCtrl.getPicList()
                lambda = this.m_pcRateCtrl.getRCPic().estimatePicLambda(listPreviousPicture)
                sliceQP = this.m_pcRateCtrl.getRCPic().estimatePicQP(lambda, listPreviousPicture)
            } else { // normal case
                listPreviousPicture := this.m_pcRateCtrl.getPicList()
                lambda = this.m_pcRateCtrl.getRCPic().estimatePicLambda(listPreviousPicture)
                sliceQP = this.m_pcRateCtrl.getRCPic().estimatePicQP(lambda, listPreviousPicture)
            }

            sliceQP = TLibCommon.CLIP3(-pcSlice.GetSPS().GetQpBDOffsetY(), TLibCommon.MAX_QP, sliceQP).(int)
            this.m_pcRateCtrl.getRCPic().setPicEstQP(sliceQP)

            this.m_pcSliceEncoder.resetQP(pcPic, sliceQP, lambda)
        }
        //#endif

        uiNumSlices := uint(1)

        uiInternalAddress := pcPic.GetNumPartInCU() - 4
        uiExternalAddress := pcPic.GetPicSym().GetNumberOfCUsInFrame() - 1
        uiPosX := (uiExternalAddress%pcPic.GetFrameWidthInCU())*pcPic.GetSlice(0).GetSPS().GetMaxCUWidth() + TLibCommon.G_auiRasterToPelX[TLibCommon.G_auiZscanToRaster[uiInternalAddress]]
        uiPosY := (uiExternalAddress/pcPic.GetFrameWidthInCU())*pcPic.GetSlice(0).GetSPS().GetMaxCUHeight() + TLibCommon.G_auiRasterToPelY[TLibCommon.G_auiZscanToRaster[uiInternalAddress]]
        uiWidth := pcSlice.GetSPS().GetPicWidthInLumaSamples()
        uiHeight := pcSlice.GetSPS().GetPicHeightInLumaSamples()
        for uiPosX >= uiWidth || uiPosY >= uiHeight {
            uiInternalAddress--
            uiPosX = (uiExternalAddress%pcPic.GetFrameWidthInCU())*pcPic.GetSlice(0).GetSPS().GetMaxCUWidth() + TLibCommon.G_auiRasterToPelX[TLibCommon.G_auiZscanToRaster[uiInternalAddress]]
            uiPosY = (uiExternalAddress/pcPic.GetFrameWidthInCU())*pcPic.GetSlice(0).GetSPS().GetMaxCUHeight() + TLibCommon.G_auiRasterToPelY[TLibCommon.G_auiZscanToRaster[uiInternalAddress]]
        }
        uiInternalAddress++
        if uiInternalAddress == pcPic.GetNumPartInCU() {
            uiInternalAddress = 0
            uiExternalAddress++
        }
        uiRealEndAddress := uiExternalAddress*pcPic.GetNumPartInCU() + uiInternalAddress

        var uiCummulativeTileWidth, uiCummulativeTileHeight uint
        var p, j int
        var uiEncCUAddr uint

        //set NumColumnsMinus1 and NumRowsMinus1
        pcPic.GetPicSym().SetNumColumnsMinus1(pcSlice.GetPPS().GetNumColumnsMinus1())
        pcPic.GetPicSym().SetNumRowsMinus1(pcSlice.GetPPS().GetNumRowsMinus1())

        //create the TLibCommon.TComTileArray
        pcPic.GetPicSym().XCreateTComTileArray()

        if pcSlice.GetPPS().GetUniformSpacingFlag() == true {
            //set the width for each tile
            for j = 0; j < pcPic.GetPicSym().GetNumRowsMinus1()+1; j++ {
                for p = 0; p < pcPic.GetPicSym().GetNumColumnsMinus1()+1; p++ {
                    pcPic.GetPicSym().GetTComTile(uint(j*(pcPic.GetPicSym().GetNumColumnsMinus1()+1) + p)).SetTileWidth(uint((p+1)*int(pcPic.GetPicSym().GetFrameWidthInCU())/(pcPic.GetPicSym().GetNumColumnsMinus1()+1) -
                        (p*int(pcPic.GetPicSym().GetFrameWidthInCU()))/(pcPic.GetPicSym().GetNumColumnsMinus1()+1)))
                }
            }

            //set the height for each tile
            for j = 0; j < pcPic.GetPicSym().GetNumColumnsMinus1()+1; j++ {
                for p = 0; p < pcPic.GetPicSym().GetNumRowsMinus1()+1; p++ {
                    pcPic.GetPicSym().GetTComTile(uint(p*(pcPic.GetPicSym().GetNumColumnsMinus1()+1) + j)).SetTileHeight(uint((p+1)*int(pcPic.GetPicSym().GetFrameHeightInCU())/(pcPic.GetPicSym().GetNumRowsMinus1()+1) -
                        (p*int(pcPic.GetPicSym().GetFrameHeightInCU()))/(pcPic.GetPicSym().GetNumRowsMinus1()+1)))
                }
            }
        } else {
            //set the width for each tile
            for j = 0; j < pcPic.GetPicSym().GetNumRowsMinus1()+1; j++ {
                uiCummulativeTileWidth = 0
                for p = 0; p < pcPic.GetPicSym().GetNumColumnsMinus1(); p++ {
                    pcPic.GetPicSym().GetTComTile(uint(j*(pcPic.GetPicSym().GetNumColumnsMinus1()+1) + p)).SetTileWidth(uint(pcSlice.GetPPS().GetColumnWidth(p)))
                    uiCummulativeTileWidth += uint(pcSlice.GetPPS().GetColumnWidth(p))
                }
                pcPic.GetPicSym().GetTComTile(uint(j*(pcPic.GetPicSym().GetNumColumnsMinus1()+1) + p)).SetTileWidth(uint(pcPic.GetPicSym().GetFrameWidthInCU()) - uiCummulativeTileWidth)
            }

            //set the height for each tile
            for j = 0; j < pcPic.GetPicSym().GetNumColumnsMinus1()+1; j++ {
                uiCummulativeTileHeight = 0
                for p = 0; p < pcPic.GetPicSym().GetNumRowsMinus1(); p++ {
                    pcPic.GetPicSym().GetTComTile(uint(p*(pcPic.GetPicSym().GetNumColumnsMinus1()+1) + j)).SetTileHeight(uint(pcSlice.GetPPS().GetRowHeight(p)))
                    uiCummulativeTileHeight += uint(pcSlice.GetPPS().GetRowHeight(p))
                }
                pcPic.GetPicSym().GetTComTile(uint(p*(pcPic.GetPicSym().GetNumColumnsMinus1()+1) + j)).SetTileHeight(uint(pcPic.GetPicSym().GetFrameHeightInCU()) - uiCummulativeTileHeight)
            }
        }
        //intialize each tile of the current picture
        pcPic.GetPicSym().XInitTiles()

        // Allocate some coders, now we know how many tiles there are.
        iNumSubstreams := pcSlice.GetPPS().GetNumSubstreams()

        //generate the Coding Order Map and Inverse Coding Order Map
        uiEncCUAddr = 0
        for p = 0; p < int(pcPic.GetPicSym().GetNumberOfCUsInFrame()); p++ {
            pcPic.GetPicSym().SetCUOrderMap(p, int(uiEncCUAddr))
            pcPic.GetPicSym().SetInverseCUOrderMap(int(uiEncCUAddr), p)
            uiEncCUAddr = pcPic.GetPicSym().XCalculateNxtCUAddr(uiEncCUAddr)
        }
        pcPic.GetPicSym().SetCUOrderMap(int(pcPic.GetPicSym().GetNumberOfCUsInFrame()), int(pcPic.GetPicSym().GetNumberOfCUsInFrame()))
        pcPic.GetPicSym().SetInverseCUOrderMap(int(pcPic.GetPicSym().GetNumberOfCUsInFrame()), int(pcPic.GetPicSym().GetNumberOfCUsInFrame()))

        // Allocate some coders, now we know how many tiles there are.
        this.m_pcEncTop.CreateWPPCoders(iNumSubstreams)
        pcSbacCoders = this.m_pcEncTop.getSbacCoders()
        pcSubstreamsOut = make([]*TLibCommon.TComOutputBitstream, iNumSubstreams)
        for i := 0; i < iNumSubstreams; i++ {
            pcSubstreamsOut[i] = TLibCommon.NewTComOutputBitstream()
        }

        startCUAddrSliceIdx := uint(0)                   // used to index "m_uiStoredStartCUAddrForEncodingSlice" containing locations of slice boundaries
        startCUAddrSlice := uint(0)                      // used to keep track of current slice's starting CU addr.
        pcSlice.SetSliceCurStartCUAddr(startCUAddrSlice) // Setting "start CU addr" for current slice
        this.m_storedStartCUAddrForEncodingSlice = make(map[int]int)

        startCUAddrSliceSegmentIdx := uint(0)                          // used to index "m_uiStoredStartCUAddrForEntropyEncodingSlice" containing locations of slice boundaries
        startCUAddrSliceSegment := uint(0)                             // used to keep track of current Dependent slice's starting CU addr.
        pcSlice.SetSliceSegmentCurStartCUAddr(startCUAddrSliceSegment) // Setting "start CU addr" for current Dependent slice

        this.m_storedStartCUAddrForEncodingSliceSegment = make(map[int]int)
        nextCUAddr := uint(0)
        this.m_storedStartCUAddrForEncodingSlice[len(this.m_storedStartCUAddrForEncodingSlice)] = int(nextCUAddr)
        startCUAddrSliceIdx++
        this.m_storedStartCUAddrForEncodingSliceSegment[len(this.m_storedStartCUAddrForEncodingSliceSegment)] = int(nextCUAddr)
        startCUAddrSliceSegmentIdx++

        for nextCUAddr < uiRealEndAddress { // determine slice boundaries
            pcSlice.SetNextSlice(false)
            pcSlice.SetNextSliceSegment(false)
            //assert(pcPic->getNumAllocatedSlice() == startCUAddrSliceIdx);
            this.m_pcSliceEncoder.precompressSlice(pcPic)
            this.m_pcSliceEncoder.compressSlice(pcPic)

            bNoBinBitConstraintViolated := (!pcSlice.IsNextSlice() && !pcSlice.IsNextSliceSegment())
            if pcSlice.IsNextSlice() || (bNoBinBitConstraintViolated && this.m_pcCfg.GetSliceMode() == TLibCommon.FIXED_NUMBER_OF_LCU) {
                startCUAddrSlice = pcSlice.GetSliceCurEndCUAddr()
                // Reconstruction slice
                this.m_storedStartCUAddrForEncodingSlice[len(this.m_storedStartCUAddrForEncodingSlice)] = int(startCUAddrSlice)
                startCUAddrSliceIdx++
                // Dependent slice
                if startCUAddrSliceSegmentIdx > 0 && this.m_storedStartCUAddrForEncodingSliceSegment[int(startCUAddrSliceSegmentIdx)-1] != int(startCUAddrSlice) {
                    this.m_storedStartCUAddrForEncodingSliceSegment[len(this.m_storedStartCUAddrForEncodingSliceSegment)] = int(startCUAddrSlice)
                    startCUAddrSliceSegmentIdx++
                }

                if startCUAddrSlice < uiRealEndAddress {
                    pcPic.AllocateNewSlice()
                    pcPic.SetCurrSliceIdx(startCUAddrSliceIdx - 1)
                    this.m_pcSliceEncoder.setSliceIdx(startCUAddrSliceIdx - 1)
                    pcSlice = pcPic.GetSlice(startCUAddrSliceIdx - 1)
                    pcSlice.CopySliceInfo(pcPic.GetSlice(0))
                    pcSlice.SetSliceIdx(startCUAddrSliceIdx - 1)
                    pcSlice.SetSliceCurStartCUAddr(startCUAddrSlice)
                    pcSlice.SetSliceSegmentCurStartCUAddr(startCUAddrSlice)
                    pcSlice.SetSliceBits(0)
                    uiNumSlices++
                }
            } else if pcSlice.IsNextSliceSegment() || (bNoBinBitConstraintViolated && this.m_pcCfg.GetSliceSegmentMode() == TLibCommon.FIXED_NUMBER_OF_LCU) {
                startCUAddrSliceSegment = pcSlice.GetSliceSegmentCurEndCUAddr()
                this.m_storedStartCUAddrForEncodingSliceSegment[len(this.m_storedStartCUAddrForEncodingSliceSegment)] = int(startCUAddrSliceSegment)
                startCUAddrSliceSegmentIdx++
                pcSlice.SetSliceSegmentCurStartCUAddr(startCUAddrSliceSegment)
            } else {
                startCUAddrSlice = pcSlice.GetSliceCurEndCUAddr()
                startCUAddrSliceSegment = pcSlice.GetSliceSegmentCurEndCUAddr()
            }
            if startCUAddrSlice > startCUAddrSliceSegment {
                nextCUAddr = startCUAddrSlice
            } else {
                nextCUAddr = startCUAddrSliceSegment
            }
        }
        this.m_storedStartCUAddrForEncodingSlice[len(this.m_storedStartCUAddrForEncodingSlice)] = int(pcSlice.GetSliceCurEndCUAddr())
        startCUAddrSliceIdx++
        this.m_storedStartCUAddrForEncodingSliceSegment[len(this.m_storedStartCUAddrForEncodingSliceSegment)] = int(pcSlice.GetSliceCurEndCUAddr())
        startCUAddrSliceSegmentIdx++

        pcSlice = pcPic.GetSlice(0)

        // SAO parameter estimation using non-deblocked pixels for LCU bottom and right boundary areas
        if this.m_pcCfg.GetSaoLcuBasedOptimization() && this.m_pcCfg.GetSaoLcuBoundary() {
            fmt.Printf("sao not implement\n")
            /*
               this.m_pcSAO.resetStats();
               this.m_pcSAO.calcSaoStatsCu_BeforeDblk( pcPic );
            */
        }

        //-- Loop filter
        bLFCrossTileBoundary := pcSlice.GetPPS().GetLoopFilterAcrossTilesEnabledFlag()
        this.m_pcLoopFilter.SetCfg(bLFCrossTileBoundary)
        this.m_pcLoopFilter.LoopFilterPic(pcPic)

        pcSlice = pcPic.GetSlice(0)
        if pcSlice.GetSPS().GetUseSAO() {
            var LFCrossSliceBoundaryFlag map[int]bool
            LFCrossSliceBoundaryFlag = make(map[int]bool, uiNumSlices)
            for s := int(0); s < int(uiNumSlices); s++ {
                if uiNumSlices == 1 {
                    LFCrossSliceBoundaryFlag[s] = true //:pcPic.GetSlice(s).getLFCrossSliceBoundaryFlag()) );
                } else {
                    LFCrossSliceBoundaryFlag[s] = pcPic.GetSlice(uint(s)).GetLFCrossSliceBoundaryFlag()
                }
            }
            //this.m_storedStartCUAddrForEncodingSlice.resize(uiNumSlices+1);
            pcPic.CreateNonDBFilterInfo(this.m_storedStartCUAddrForEncodingSlice, 0, LFCrossSliceBoundaryFlag, pcPic.GetPicSym().GetNumTiles(), bLFCrossTileBoundary)
        }

        pcSlice = pcPic.GetSlice(0)

        if pcSlice.GetSPS().GetUseSAO() {
            fmt.Printf("sao not implemented\n")
            //this.m_pcSAO.createPicSaoInfo(pcPic);
        }

        /////////////////////////////////////////////////////////////////////////////////////////////////// File writing
        // Set entropy coder
        this.m_pcEntropyCoder.setEntropyCoder(this.m_pcCavlcCoder, pcSlice, this.m_pTraceFile)
        
        // write various header sets.
        if this.m_bSeqFirst {
            nalu := NewOutputNALUnit(TLibCommon.NAL_UNIT_VPS, 0, 0)
            this.m_pcEntropyCoder.setBitstream(nalu.m_Bitstream)
            
            //fmt.Printf("vps=%d\n",nalu.m_Bitstream.GetNumberOfWrittenBits());
            
            this.m_pcEntropyCoder.encodeVPS(this.m_pcCfg.GetVPS())
            nalu.m_Bitstream.WriteRBSPTrailingBits()
            
            //fmt.Printf("vps=%d\n",nalu.m_Bitstream.GetNumberOfWrittenBits());
      		naluEbsp := NewNALUnitEBSP(nalu);
            accessUnit.PushBack(naluEbsp);
            //#if RATE_CONTROL_LAMBDA_DOMAIN
            actualTotalBits += int(naluEbsp.m_Bitstream.GetByteStreamLength() * 8)
            //#endif
			//fmt.Printf("vps=%d\n",naluEbsp.m_Bitstream.GetByteStreamLength()*8);
			
            nalu = NewOutputNALUnit(TLibCommon.NAL_UNIT_SPS, 0, 0)
            this.m_pcEntropyCoder.setBitstream(nalu.m_Bitstream)
            if this.m_bSeqFirst {
                pcSlice.GetSPS().SetNumLongTermRefPicSPS(this.m_numLongTermRefPicSPS)
                for k := uint(0); k < this.m_numLongTermRefPicSPS; k++ {
                    pcSlice.GetSPS().SetLtRefPicPocLsbSps(k, this.m_ltRefPicPocLsbSps[k])
                    pcSlice.GetSPS().SetUsedByCurrPicLtSPSFlag(int(k), this.m_ltRefPicUsedByCurrPicFlag[k] != 0)
                }
            }
            if this.m_pcCfg.GetPictureTimingSEIEnabled() != 0 || this.m_pcCfg.GetDecodingUnitInfoSEIEnabled() != 0 {
                maxCU := this.m_pcCfg.GetSliceArgument() >> (pcSlice.GetSPS().GetMaxCUDepth() << 1)
                var numDU uint
                if this.m_pcCfg.GetSliceMode() == 1 {
                    numDU = (pcPic.GetNumCUsInFrame() / uint(maxCU))
                } else {
                    numDU = (0)
                }

                if pcPic.GetNumCUsInFrame()%uint(maxCU) != 0 {
                    numDU++
                }
                pcSlice.GetSPS().GetVuiParameters().GetHrdParameters().SetNumDU(numDU)
                pcSlice.GetSPS().SetHrdParameters(uint(this.m_pcCfg.GetFrameRate()), numDU, uint(this.m_pcCfg.GetTargetBitrate()), (this.m_pcCfg.GetIntraPeriod() > 0))
            }
            if this.m_pcCfg.GetBufferingPeriodSEIEnabled() != 0 || this.m_pcCfg.GetPictureTimingSEIEnabled() != 0 || this.m_pcCfg.GetDecodingUnitInfoSEIEnabled() != 0 {
                pcSlice.GetSPS().GetVuiParameters().SetHrdParametersPresentFlag(true)
            }
            this.m_pcEntropyCoder.encodeSPS(pcSlice.GetSPS())
            nalu.m_Bitstream.WriteRBSPTrailingBits()
            naluEbsp = NewNALUnitEBSP(nalu);
            accessUnit.PushBack(naluEbsp);
            //#if RATE_CONTROL_LAMBDA_DOMAIN
            actualTotalBits += int(naluEbsp.m_Bitstream.GetByteStreamLength() * 8)
            //#endif

            nalu = NewOutputNALUnit(TLibCommon.NAL_UNIT_PPS, 0, 0)
            this.m_pcEntropyCoder.setBitstream(nalu.m_Bitstream)
            this.m_pcEntropyCoder.encodePPS(pcSlice.GetPPS())
            nalu.m_Bitstream.WriteRBSPTrailingBits()
            naluEbsp = NewNALUnitEBSP(nalu);
            accessUnit.PushBack(naluEbsp);
            //#if RATE_CONTROL_LAMBDA_DOMAIN
            actualTotalBits += int(naluEbsp.m_Bitstream.GetByteStreamLength() * 8)
            //#endif

            this.xCreateLeadingSEIMessages(accessUnit, pcSlice.GetSPS())

            this.m_bSeqFirst = false
        }

        if (this.m_pcCfg.GetPictureTimingSEIEnabled() != 0 || this.m_pcCfg.GetDecodingUnitInfoSEIEnabled() != 0) &&
            (pcSlice.GetSPS().GetVuiParametersPresentFlag()) &&
            ((pcSlice.GetSPS().GetVuiParameters().GetHrdParameters().GetNalHrdParametersPresentFlag()) ||
                (pcSlice.GetSPS().GetVuiParameters().GetHrdParameters().GetVclHrdParametersPresentFlag())) {
            /*
               if pcSlice.GetSPS().GetVuiParameters().GetSubPicCpbParamsPresentFlag() {
                 numDU := pcSlice.GetSPS().GetVuiParameters().GetNumDU();
                 pictureTimingSEI.m_numDecodingUnitsMinus1     = ( numDU - 1 );
                 pictureTimingSEI.m_duCommonCpbRemovalDelayFlag = 0;

                 if pictureTimingSEI.m_numNalusInDuMinus1 == nil {
                   pictureTimingSEI.m_numNalusInDuMinus1       = make([]uint, numDU );
                 }
                 if pictureTimingSEI.m_duCpbRemovalDelayMinus1  == nil {
                   pictureTimingSEI.m_duCpbRemovalDelayMinus1  = make([]uint, numDU );
                 }
                 if accumBitsDU == nil {
                   accumBitsDU                                  = make([]uint, numDU );
                 }
                 if accumNalsDU == nil {
                   accumNalsDU                                  = make([]uint, numDU );
                 }
               }
               pictureTimingSEI.m_auCpbRemovalDelay = TLibCommon.MAX(1,this.m_totalCoded - this.m_lastBPSEI);
               pictureTimingSEI.m_picDpbOutputDelay = pcSlice.GetSPS().getNumReorderPics(0) + pcSlice.GetPOC() - this.m_totalCoded;
            */
        }
        if (this.m_pcCfg.GetBufferingPeriodSEIEnabled() != 0) && (pcSlice.GetSliceType() == TLibCommon.I_SLICE) &&
            (pcSlice.GetSPS().GetVuiParametersPresentFlag()) &&
            ((pcSlice.GetSPS().GetVuiParameters().GetHrdParameters().GetNalHrdParametersPresentFlag()) ||
                (pcSlice.GetSPS().GetVuiParameters().GetHrdParameters().GetVclHrdParametersPresentFlag())) {
            /*
                   nalu = NewOutputNALUnit(TLibCommon.NAL_UNIT_SEI);
                  this.m_pcEntropyCoder.setEntropyCoder(this.m_pcCavlcCoder, pcSlice);
                  this.m_pcEntropyCoder.setTraceFile(this.m_pTraceFile)
                  this.m_pcEntropyCoder.setBitstream(&nalu.m_Bitstream);

                  var sei_buffering_period SEIBufferingPeriod;

                  uiInitialCpbRemovalDelay := uint(90000/2);                      // 0.5 sec
                  sei_buffering_period.m_initialCpbRemovalDelay      [0][0]     = uiInitialCpbRemovalDelay;
                  sei_buffering_period.m_initialCpbRemovalDelayOffset[0][0]     = uiInitialCpbRemovalDelay;
                  sei_buffering_period.m_initialCpbRemovalDelay      [0][1]     = uiInitialCpbRemovalDelay;
                  sei_buffering_period.m_initialCpbRemovalDelayOffset[0][1]     = uiInitialCpbRemovalDelay;

            #if L0043_TIMING_INFO
                  Double dTmp = (Double)pcSlice->getSPS()->getVuiParameters()->getTimingInfo()->getNumUnitsInTick() / (Double)pcSlice->getSPS()->getVuiParameters()->getTimingInfo()->getTimeScale();
            #else
                  Double dTmp = (Double)pcSlice->getSPS()->getVuiParameters()->getHrdParameters()->getNumUnitsInTick() / (Double)pcSlice->getSPS()->getVuiParameters()->getHrdParameters()->getTimeScale();
            #endif

                  uiTmp := uint( dTmp * 90000.0 );
                  uiInitialCpbRemovalDelay -= uiTmp;
                  uiInitialCpbRemovalDelay -= uiTmp / ( pcSlice.GetSPS().getVuiParameters().getTickDivisorMinus2() + 2 );
                  sei_buffering_period.m_initialAltCpbRemovalDelay      [0][0]  = uiInitialCpbRemovalDelay;
                  sei_buffering_period.m_initialAltCpbRemovalDelayOffset[0][0]  = uiInitialCpbRemovalDelay;
                  sei_buffering_period.m_initialAltCpbRemovalDelay      [0][1]  = uiInitialCpbRemovalDelay;
                  sei_buffering_period.m_initialAltCpbRemovalDelayOffset[0][1]  = uiInitialCpbRemovalDelay;

                  sei_buffering_period.m_altCpbParamsPresentFlag              = 0;
                  sei_buffering_period.m_sps                                  = pcSlice.GetSPS();

                  this.m_seiWriter.writeSEImessage( nalu.m_Bitstream, sei_buffering_period );
                  writeRBSPTrailingBits(nalu.m_Bitstream);
                  accessUnit.push_back(NewNALUnitEBSP(nalu));

                  this.m_lastBPSEI = this.m_totalCoded;
                  this.m_cpbRemovalDelay = 0;
            */
        }
        this.m_cpbRemovalDelay++
        if (this.m_pcCfg.GetRecoveryPointSEIEnabled() != 0) && (pcSlice.GetSliceType() == TLibCommon.I_SLICE) {
            // Recovery point SEI
            /*
                nalu := NewOutputNALUnit(TLibCommon.NAL_UNIT_SEI);
               this.m_pcEntropyCoder.setEntropyCoder(this.m_pcCavlcCoder, pcSlice);
               this.m_pcEntropyCoder.setTraceFile(this.m_pTraceFile)
               this.m_pcEntropyCoder.setBitstream(&nalu.m_Bitstream);

               var sei_recovery_point SEIRecoveryPoint;
               sei_recovery_point.m_recoveryPocCnt    = 0;
               sei_recovery_point.m_exactMatchingFlag =  pcSlice.GetPOC() == 0 ;
               sei_recovery_point.m_brokenLinkFlag    = false;

               this.m_seiWriter.writeSEImessage( nalu.m_Bitstream, sei_recovery_point );
               writeRBSPTrailingBits(nalu.m_Bitstream);
               accessUnit.push_back(NewNALUnitEBSP(nalu));
            */
        }
        // use the main bitstream buffer for storing the marshalled picture
        this.m_pcEntropyCoder.setBitstream(nil)

        startCUAddrSliceIdx = 0
        startCUAddrSlice = 0

        startCUAddrSliceSegmentIdx = 0
        startCUAddrSliceSegment = 0
        nextCUAddr = 0
        pcSlice = pcPic.GetSlice(startCUAddrSliceIdx)

        var processingState PROCESSING_STATE
        if pcSlice.GetSPS().GetUseSAO() {
            processingState = EXECUTE_INLOOPFILTER
        } else {
            processingState = ENCODE_SLICE
        }

        skippedSlice := false
        for nextCUAddr < uiRealEndAddress { // Iterate over all slices
            switch processingState {
            case ENCODE_SLICE:
                pcSlice.SetNextSlice(false)
                pcSlice.SetNextSliceSegment(false)
                if nextCUAddr == uint(this.m_storedStartCUAddrForEncodingSlice[int(startCUAddrSliceIdx)]) {
                    pcSlice = pcPic.GetSlice(startCUAddrSliceIdx)
                    if startCUAddrSliceIdx > 0 && pcSlice.GetSliceType() != TLibCommon.I_SLICE {
                        pcSlice.CheckColRefIdx(startCUAddrSliceIdx, pcPic)
                    }
                    pcPic.SetCurrSliceIdx(startCUAddrSliceIdx)
                    this.m_pcSliceEncoder.setSliceIdx(startCUAddrSliceIdx)
                    //assert(startCUAddrSliceIdx == pcSlice.GetSliceIdx());
                    // Reconstruction slice
                    pcSlice.SetSliceCurStartCUAddr(nextCUAddr) // to be used in encodeSlice() + context restriction
                    pcSlice.SetSliceCurEndCUAddr(uint(this.m_storedStartCUAddrForEncodingSlice[int(startCUAddrSliceIdx+1)]))
                    // Dependent slice
                    pcSlice.SetSliceSegmentCurStartCUAddr(nextCUAddr) // to be used in encodeSlice() + context restriction
                    pcSlice.SetSliceSegmentCurEndCUAddr(uint(this.m_storedStartCUAddrForEncodingSliceSegment[int(startCUAddrSliceSegmentIdx+1)]))

                    pcSlice.SetNextSlice(true)

                    startCUAddrSliceIdx++
                    startCUAddrSliceSegmentIdx++
                } else if nextCUAddr == uint(this.m_storedStartCUAddrForEncodingSliceSegment[int(startCUAddrSliceSegmentIdx)]) {
                    // Dependent slice
                    pcSlice.SetSliceSegmentCurStartCUAddr(nextCUAddr) // to be used in encodeSlice() + context restriction
                    pcSlice.SetSliceSegmentCurEndCUAddr(uint(this.m_storedStartCUAddrForEncodingSliceSegment[int(startCUAddrSliceSegmentIdx+1)]))

                    pcSlice.SetNextSliceSegment(true)

                    startCUAddrSliceSegmentIdx++
                }

                pcSlice.SetRPS(pcPic.GetSlice(0).GetRPS())
                pcSlice.SetRPSidx(pcPic.GetSlice(0).GetRPSidx())
                var uiDummyStartCUAddr, uiDummyBoundingCUAddr uint
                this.m_pcSliceEncoder.xDetermineStartAndBoundingCUAddr(&uiDummyStartCUAddr, &uiDummyBoundingCUAddr, pcPic, true)

                uiInternalAddress = pcPic.GetPicSym().GetPicSCUAddr(pcSlice.GetSliceSegmentCurEndCUAddr()-1) % pcPic.GetNumPartInCU()
                uiExternalAddress = pcPic.GetPicSym().GetPicSCUAddr(pcSlice.GetSliceSegmentCurEndCUAddr()-1) / pcPic.GetNumPartInCU()
                uiPosX = (uiExternalAddress%pcPic.GetFrameWidthInCU())*pcPic.GetSlice(0).GetSPS().GetMaxCUWidth() + TLibCommon.G_auiRasterToPelX[TLibCommon.G_auiZscanToRaster[uiInternalAddress]]
                uiPosY = (uiExternalAddress/pcPic.GetFrameWidthInCU())*pcPic.GetSlice(0).GetSPS().GetMaxCUHeight() + TLibCommon.G_auiRasterToPelY[TLibCommon.G_auiZscanToRaster[uiInternalAddress]]
                uiWidth = pcSlice.GetSPS().GetPicWidthInLumaSamples()
                uiHeight = pcSlice.GetSPS().GetPicHeightInLumaSamples()
                for uiPosX >= uiWidth || uiPosY >= uiHeight {
                    uiInternalAddress--
                    uiPosX = (uiExternalAddress%pcPic.GetFrameWidthInCU())*pcPic.GetSlice(0).GetSPS().GetMaxCUWidth() + TLibCommon.G_auiRasterToPelX[TLibCommon.G_auiZscanToRaster[uiInternalAddress]]
                    uiPosY = (uiExternalAddress/pcPic.GetFrameWidthInCU())*pcPic.GetSlice(0).GetSPS().GetMaxCUHeight() + TLibCommon.G_auiRasterToPelY[TLibCommon.G_auiZscanToRaster[uiInternalAddress]]
                }
                uiInternalAddress++
                if uiInternalAddress == pcPic.GetNumPartInCU() {
                    uiInternalAddress = 0
                    uiExternalAddress = pcPic.GetPicSym().GetCUOrderMap(int(pcPic.GetPicSym().GetInverseCUOrderMap(int(uiExternalAddress)) + 1))
                }
                endAddress := pcPic.GetPicSym().GetPicSCUEncOrder(uiExternalAddress*pcPic.GetNumPartInCU() + uiInternalAddress)
                if endAddress <= pcSlice.GetSliceSegmentCurStartCUAddr() {
                    var boundingAddrSlice, boundingAddrSliceSegment uint
                    boundingAddrSlice = uint(this.m_storedStartCUAddrForEncodingSlice[int(startCUAddrSliceIdx)])
                    boundingAddrSliceSegment = uint(this.m_storedStartCUAddrForEncodingSliceSegment[int(startCUAddrSliceSegmentIdx)])
                    nextCUAddr = TLibCommon.MIN(boundingAddrSlice, boundingAddrSliceSegment).(uint)
                    if pcSlice.IsNextSlice() {
                        skippedSlice = true
                    }
                    continue
                }
                if skippedSlice {
                    pcSlice.SetNextSlice(true)
                    pcSlice.SetNextSliceSegment(false)
                }
                skippedSlice = false
                pcSlice.AllocSubstreamSizes(uint(iNumSubstreams))
                for ui := 0; ui < iNumSubstreams; ui++ {
                    pcSubstreamsOut[ui].Clear()
                }

                this.m_pcEntropyCoder.setEntropyCoder(this.m_pcCavlcCoder, pcSlice, this.m_pTraceFile)
                this.m_pcEntropyCoder.resetEntropy()
                // start slice NALunit
                nalu := NewOutputNALUnit(pcSlice.GetNalUnitType(), pcSlice.GetTLayer(), 0)
                sliceSegment := (!pcSlice.IsNextSlice())
                if !sliceSegment {
                    uiOneBitstreamPerSliceLength = 0 // start of a new slice
                }
                this.m_pcEntropyCoder.setBitstream(nalu.m_Bitstream)
                //#if RATE_CONTROL_LAMBDA_DOMAIN
                tmpBitsBeforeWriting = int(this.m_pcEntropyCoder.getNumberOfWrittenBits())
                //#endif
                this.m_pcEntropyCoder.encodeSliceHeader(pcSlice)
                //#if RATE_CONTROL_LAMBDA_DOMAIN
                actualHeadBits += (int(this.m_pcEntropyCoder.getNumberOfWrittenBits()) - tmpBitsBeforeWriting)
                //#endif

                // is it needed?
                {
                    if !sliceSegment {
                        pcBitstreamRedirect.WriteAlignOne()
                    } else {
                        // We've not completed our slice header info yet, do the alignment later.
                    }
                    this.m_pcSbacCoder.init(this.m_pcBinCABAC)
                    this.m_pcEntropyCoder.setEntropyCoder(this.m_pcSbacCoder, pcSlice, this.m_pTraceFile)
                    this.m_pcEntropyCoder.resetEntropy()
                    for ui := 0; ui < pcSlice.GetPPS().GetNumSubstreams(); ui++ {
                        this.m_pcEntropyCoder.setEntropyCoder(pcSbacCoders[ui], pcSlice, this.m_pTraceFile)
                        this.m_pcEntropyCoder.resetEntropy()
                    }
                }

                if pcSlice.IsNextSlice() {
                    // set entropy coder for writing
                    this.m_pcSbacCoder.init(this.m_pcBinCABAC)
                    {
                        for ui := 0; ui < pcSlice.GetPPS().GetNumSubstreams(); ui++ {
                            this.m_pcEntropyCoder.setEntropyCoder(pcSbacCoders[ui], pcSlice, this.m_pTraceFile)
                            this.m_pcEntropyCoder.resetEntropy()
                        }
                        pcSbacCoders[0].load(this.m_pcSbacCoder)
                        this.m_pcEntropyCoder.setEntropyCoder(pcSbacCoders[0], pcSlice, this.m_pTraceFile) //ALF is written in substream #0 with CABAC coder #0 (see ALF param encoding below)
                    }
                    this.m_pcEntropyCoder.resetEntropy()
                    // File writing
                    if !sliceSegment {
                        this.m_pcEntropyCoder.setBitstream(pcBitstreamRedirect)
                    } else {
                        this.m_pcEntropyCoder.setBitstream(nalu.m_Bitstream)
                    }
                    // for now, override the TILES_DECODER setting in order to write substreams.
                    this.m_pcEntropyCoder.setBitstream(pcSubstreamsOut[0])

                }
                pcSlice.SetFinalized(true)

                this.m_pcSbacCoder.load(pcSbacCoders[0])

                pcSlice.SetTileOffstForMultES(uiOneBitstreamPerSliceLength)
                if !sliceSegment {
                    pcSlice.SetTileLocationCount1(0)
                    this.m_pcSliceEncoder.encodeSlice(pcPic, pcBitstreamRedirect, pcSubstreamsOut) // redirect is only used for CAVLC tile position info.
                } else {
                    this.m_pcSliceEncoder.encodeSlice(pcPic, nalu.m_Bitstream, pcSubstreamsOut) // nalu.m_Bitstream is only used for CAVLC tile position info.
                }

                {
                    // Construct the final bitstream by flushing and concatenating substreams.
                    // The final bitstream is either nalu.m_Bitstream or pcBitstreamRedirect;
                    puiSubstreamSizes := pcSlice.GetSubstreamSizes()
                    uiTotalCodedSize := uint(0) // for padding calcs.
                    uiNumSubstreamsPerTile := iNumSubstreams
                    if iNumSubstreams > 1 {
                        uiNumSubstreamsPerTile /= pcPic.GetPicSym().GetNumTiles()
                    }
                    for ui := 0; ui < iNumSubstreams; ui++ {
                        // Flush all substreams -- this includes empty ones.
                        // Terminating bit and flush.
                        this.m_pcEntropyCoder.setEntropyCoder(pcSbacCoders[ui], pcSlice, this.m_pTraceFile)
                        this.m_pcEntropyCoder.setBitstream(pcSubstreamsOut[ui])
                        this.m_pcEntropyCoder.encodeTerminatingBit(1)
                        this.m_pcEntropyCoder.encodeSliceFinish()

                        pcSubstreamsOut[ui].WriteByteAlignment() // Byte-alignment in slice_data() at end of sub-stream
                        // Byte alignment is necessary between tiles when tiles are independent.
                        uiTotalCodedSize += pcSubstreamsOut[ui].GetNumberOfWrittenBits()

                        bNextSubstreamInNewTile := ((ui + 1) < iNumSubstreams) && ((ui+1)%uiNumSubstreamsPerTile == 0)
                        if bNextSubstreamInNewTile {
                            pcSlice.SetTileLocation(ui/uiNumSubstreamsPerTile, pcSlice.GetTileOffstForMultES()+(uiTotalCodedSize>>3))
                        }
                        if ui+1 < pcSlice.GetPPS().GetNumSubstreams() {
                            puiSubstreamSizes[ui] = pcSubstreamsOut[ui].GetNumberOfWrittenBits()
                        }
                    }

                    // Complete the slice header info.
                    this.m_pcEntropyCoder.setEntropyCoder(this.m_pcCavlcCoder, pcSlice, this.m_pTraceFile)
                    this.m_pcEntropyCoder.setBitstream(nalu.m_Bitstream)
                    this.m_pcEntropyCoder.encodeTilesWPPEntryPoint(pcSlice)

                    // Substreams...
                    pcOut := pcBitstreamRedirect
                    offs := uint(0)
                    nss := pcSlice.GetPPS().GetNumSubstreams()
                    if pcSlice.GetPPS().GetEntropyCodingSyncEnabledFlag() {
                        // 1st line present for WPP.
                        offs = pcSlice.GetSliceSegmentCurStartCUAddr() / pcSlice.GetPic().GetNumPartInCU() / pcSlice.GetPic().GetFrameWidthInCU()
                        nss = pcSlice.GetNumEntryPointOffsets() + 1
                    }
                    for ui := 0; ui < nss; ui++ {
                        pcOut.AddSubstream(pcSubstreamsOut[ui+int(offs)])
                    }
                }

                var boundingAddrSlice, boundingAddrSliceSegment uint
                boundingAddrSlice = uint(this.m_storedStartCUAddrForEncodingSlice[int(startCUAddrSliceIdx)])
                boundingAddrSliceSegment = uint(this.m_storedStartCUAddrForEncodingSliceSegment[int(startCUAddrSliceSegmentIdx)])
                nextCUAddr = TLibCommon.MIN(boundingAddrSlice, boundingAddrSliceSegment).(uint)
                // If current NALU is the first NALU of slice (containing slice header) and more NALUs exist (due to multiple dependent slices) then buffer it.
                // If current NALU is the last NALU of slice and a NALU was buffered, then (a) Write current NALU (b) Update an write buffered NALU at approproate location in NALU list.
                bNALUAlignedWrittenToList := false // used to ensure current NALU is not written more than once to the NALU list.
                this.xWriteTileLocationToSliceHeader(nalu, &pcBitstreamRedirect, pcSlice)
                naluEbsp := NewNALUnitEBSP(nalu);
            	accessUnit.PushBack(naluEbsp);
                //#if RATE_CONTROL_LAMBDA_DOMAIN
                actualTotalBits += int(naluEbsp.m_Bitstream.GetByteStreamLength() * 8)
                //#endif
                bNALUAlignedWrittenToList = true
                uiOneBitstreamPerSliceLength += nalu.m_Bitstream.GetNumberOfWrittenBits() // length of bitstream after byte-alignment
				//fmt.Printf("nalu.m_Bitstream.GetNumberOfWrittenBits()=%d\n", nalu.m_Bitstream.GetNumberOfWrittenBits());
				
                if !bNALUAlignedWrittenToList {
                    nalu.m_Bitstream.WriteAlignZero()
					
					naluEbsp = NewNALUnitEBSP(nalu);
            		accessUnit.PushBack(naluEbsp);
                    uiOneBitstreamPerSliceLength += nalu.m_Bitstream.GetNumberOfWrittenBits() + 24 // length of bitstream after byte-alignment + 3 byte startcode 0x000001
                }

                if (this.m_pcCfg.GetPictureTimingSEIEnabled() != 0 || this.m_pcCfg.GetDecodingUnitInfoSEIEnabled() != 0) &&
                    (pcSlice.GetSPS().GetVuiParametersPresentFlag()) &&
                    ((pcSlice.GetSPS().GetVuiParameters().GetHrdParameters().GetNalHrdParametersPresentFlag()) ||
                        (pcSlice.GetSPS().GetVuiParameters().GetHrdParameters().GetVclHrdParametersPresentFlag())) &&
                    (pcSlice.GetSPS().GetVuiParameters().GetHrdParameters().GetSubPicCpbParamsPresentFlag()) {
                    /*
                                 numNalus := uint(0);
                                 numRBSPBytes := uint(0);
                                for it := accessUnit.Front(); it != nil; it=it.Next() {
                                   //v:= it.Value.();
                                   numRBSPBytes_nal := uint((*it).m_nalUnitData.str().size());
                    //#if HM9_NALU_TYPES
                                  if (*it).m_nalUnitType != TLibCommon.NAL_UNIT_SEI && (*it).m_nalUnitType != TLibCommon.NAL_UNIT_SEI_SUFFIX {
                    //#else
                    //              if ((*it).m_nalUnitType != TLibCommon.NAL_UNIT_SEI)
                    //#endif

                                    numRBSPBytes += numRBSPBytes_nal;
                                  }
                                }
                                accumBitsDU[ pcSlice.GetSliceIdx() ] = ( numRBSPBytes << 3 );
                                accumNalsDU[ pcSlice.GetSliceIdx() ] = uint(accessUnit.Len());
                    */
                }
                processingState = ENCODE_SLICE

            case EXECUTE_INLOOPFILTER:
                // set entropy coder for RD
                this.m_pcEntropyCoder.setEntropyCoder(this.m_pcSbacCoder, pcSlice, this.m_pTraceFile)
                if pcSlice.GetSPS().GetUseSAO() {
                    fmt.Printf("sao not enabled\n")
                    /*
                                  this.m_pcEntropyCoder.resetEntropy();
                                  this.m_pcEntropyCoder.setBitstream( this.m_pcBitCounter );
                                  this.m_pcSAO.startSaoEnc(pcPic, this.m_pcEntropyCoder, this.m_pcEncTop.getRDSbacCoder(), m_pcEncTop.getRDGoOnSbacCoder());
                                  cSaoParam := pcSlice.GetPic().getPicSym().GetSaoParam();

                    //#if SAO_CHROMA_LAMBDA
                    //#if SAO_ENCODING_CHOICE
                                  this.m_pcSAO.SAOProcess(&cSaoParam, pcPic.GetSlice(0).getLambdaLuma(), pcPic.GetSlice(0).getLambdaChroma(), pcPic.GetSlice(0).getDepth());
                    //#else
                    //              this.m_pcSAO.SAOProcess(&cSaoParam, pcPic.GetSlice(0).getLambdaLuma(), pcPic.GetSlice(0).getLambdaChroma());
                    //#endif
                    //#else
                    //              this.m_pcSAO.SAOProcess(&cSaoParam, pcPic.GetSlice(0).getLambda());
                    //#endif
                                  this.m_pcSAO.endSaoEnc();
                                  this.m_pcSAO.PCMLFDisableProcess(pcPic);
                    */
                }
                //#if SAO_RDO
                this.m_pcEntropyCoder.setEntropyCoder(this.m_pcCavlcCoder, pcSlice, this.m_pTraceFile)
                //#endif
                processingState = ENCODE_SLICE

                for s := uint(0); s < uiNumSlices; s++ {
                    if pcSlice.GetSPS().GetUseSAO() {
                        pcPic.GetSlice(s).SetSaoEnabledFlag((pcSlice.GetPic().GetPicSym().GetSaoParam().SaoFlag[0] == true))
                    }
                }
            default:
                fmt.Printf("Not a supported encoding state\n")

            }
        }   // end iteration over slices

        if pcSlice.GetSPS().GetUseSAO() {
            if pcSlice.GetSPS().GetUseSAO() {
                fmt.Printf("sao not enabled\n")
                //this.m_pcSAO.destroyPicSaoInfo();
            }
            pcPic.DestroyNonDBFilterInfo()
        }

        pcPic.CompressMotion()

        //-- For time output for each slice
        dEncTime := time.Now().Sub(iBeforeTime) //Double dEncTime = (Double)(clock()-iBeforeTime) / CLOCKS_PER_SEC;

        var digestStr string
        if this.m_pcCfg.GetDecodedPictureHashSEIEnabled() != 0 {
            //calculate MD5sum for entire reconstructed picture
            /*
                    var sei_recon_picture_digest	SEIDecodedPictureHash;
                    if this.m_pcCfg.GetDecodedPictureHashSEIEnabled() == 1 {
                      sei_recon_picture_digest.method = SEIDecodedPictureHash::MD5;
                      calcMD5(*pcPic.GetPicYuvRec(), sei_recon_picture_digest.digest);
                      digestStr = digestToString(sei_recon_picture_digest.digest, 16);
                    }else if(this.m_pcCfg.GetDecodedPictureHashSEIEnabled() == 2{
                      sei_recon_picture_digest.method = SEIDecodedPictureHash::CRC;
                      calcCRC(*pcPic.GetPicYuvRec(), sei_recon_picture_digest.digest);
                      digestStr = digestToString(sei_recon_picture_digest.digest, 2);
                    }
                    else if(this.m_pcCfg.GetDecodedPictureHashSEIEnabled() == 3)
                    {
                      sei_recon_picture_digest.method = SEIDecodedPictureHash::CHECKSUM;
                      calcChecksum(*pcPic.GetPicYuvRec(), sei_recon_picture_digest.digest);
                      digestStr = digestToString(sei_recon_picture_digest.digest, 4);
                    }
            //#if SUFFIX_SEI_NUT_DECODED_HASH_SEI
                    OutputNALUnit nalu(TLibCommon.NAL_UNIT_SEI_SUFFIX, pcSlice.GetTLayer());
            //#else
            //        OutputNALUnit nalu(TLibCommon.NAL_UNIT_SEI, pcSlice.GetTLayer());
            //#endif

                    //write the SEI messages
                    this.m_pcEntropyCoder.setEntropyCoder(this.m_pcCavlcCoder, pcSlice);
                    this.m_pcEntropyCoder.setTraceFile(this.m_pTraceFile)
                    this.m_seiWriter.writeSEImessage(nalu.m_Bitstream, sei_recon_picture_digest);
                    writeRBSPTrailingBits(nalu.m_Bitstream);

            //#if SUFFIX_SEI_NUT_DECODED_HASH_SEI
                    accessUnit.insert(accessUnit.end(), new NALUnitEBSP(nalu));
            //#else
                    // insert the SEI message NALUnit before any Slice NALUnits
            //        AccessUnit::iterator it = find_if(accessUnit.begin(), accessUnit.end(), mem_fun(&NALUnit::isSlice));
            //        accessUnit.insert(it, new NALUnitEBSP(nalu));
            //#endif
            */
        }
        //#if SEI_TEMPORAL_LEVEL0_INDEX
        if this.m_pcCfg.GetTemporalLevel0IndexSEIEnabled() != 0 {
            /*
               var sei_temporal_level0_index	SEITemporalLevel0Index;
               if pcSlice.GetRapPicFlag() {
                 this.m_tl0Idx = 0;
                 this.m_rapIdx = (this.m_rapIdx + 1) & 0xFF;
               }else{
                 if pcSlice.GetTLayer()!=0 {
                 	this.m_tl0Idx = (this.m_tl0Idx + 0) & 0xFF;
                 }else{
                 	this.m_tl0Idx = (this.m_tl0Idx + 1) & 0xFF;
                 }
               }
               sei_temporal_level0_index.tl0Idx = this.m_tl0Idx;
               sei_temporal_level0_index.rapIdx = this.m_rapIdx;

                nalu := NewOutputNALUnit(TLibCommon.NAL_UNIT_SEI);

               // write the SEI messages
               this.m_pcEntropyCoder.setEntropyCoder(this.m_pcCavlcCoder, pcSlice);
               this.m_pcEntropyCoder.setTraceFile(this.m_pTraceFile)
               this.m_seiWriter.writeSEImessage(nalu.m_Bitstream, sei_temporal_level0_index);
               writeRBSPTrailingBits(nalu.m_Bitstream);

               // insert the SEI message NALUnit before any Slice NALUnits
               //???it := find_if(accessUnit.begin(), accessUnit.end(), (&NALUnit::isSlice));
               accessUnit.insert(it, NewNALUnitEBSP(nalu));
            */
        }
        //#endif

        this.xCalculateAddPSNR(pcPic, pcPic.GetPicYuvRec(), accessUnit, dEncTime)

        if digestStr != "" {
            if this.m_pcCfg.GetDecodedPictureHashSEIEnabled() == 1 {
                fmt.Printf(" [MD5:%s]", digestStr)
            } else if this.m_pcCfg.GetDecodedPictureHashSEIEnabled() == 2 {
                fmt.Printf(" [CRC:%s]", digestStr)
            } else if this.m_pcCfg.GetDecodedPictureHashSEIEnabled() == 3 {
                fmt.Printf(" [Checksum:%s]", digestStr)
            }
        }
        //#if RATE_CONTROL_LAMBDA_DOMAIN
        if this.m_pcCfg.GetUseRateCtrl() {
            effectivePercentage := this.m_pcRateCtrl.getRCPic().getEffectivePercentage()
            avgQP := this.m_pcRateCtrl.getRCPic().calAverageQP()
            avgLambda := this.m_pcRateCtrl.getRCPic().calAverageLambda()
            if avgLambda < 0.0 {
                avgLambda = lambda
            }
            this.m_pcRateCtrl.getRCPic().updateAfterPicture(actualHeadBits, actualTotalBits, avgQP, avgLambda, effectivePercentage)
            this.m_pcRateCtrl.getRCPic().addToPictureLsit(this.m_pcRateCtrl.getPicList())

            this.m_pcRateCtrl.getRCSeq().updateAfterPic(int64(actualTotalBits))
            if pcSlice.GetSliceType() != TLibCommon.I_SLICE {
                this.m_pcRateCtrl.getRCGOP().updateAfterPicture(actualTotalBits)
            } else { // for intra picture, the estimated bits are used to update the current status in the GOP
                this.m_pcRateCtrl.getRCGOP().updateAfterPicture(estimatedBits)
            }
        }
        //#else
        //      if(this.m_pcCfg.GetUseRateCtrl())
        //      {
        //        UInt  frameBits = this.m_vRVM_RP[this.m_vRVM_RP.size()-1];
        //        this.m_pcRateCtrl.updataRCFrameStatus((Int)frameBits, pcSlice.GetSliceType());
        //      }
        //#endif
        if (this.m_pcCfg.GetPictureTimingSEIEnabled() != 0 || this.m_pcCfg.GetDecodingUnitInfoSEIEnabled() != 0) &&
            (pcSlice.GetSPS().GetVuiParametersPresentFlag()) &&
            ((pcSlice.GetSPS().GetVuiParameters().GetHrdParameters().GetNalHrdParametersPresentFlag()) ||
                (pcSlice.GetSPS().GetVuiParameters().GetHrdParameters().GetVclHrdParametersPresentFlag())) {
            /*
                    nalu := NewOutputNALUnit(TLibCommon.NAL_UNIT_SEI, pcSlice.GetTLayer());
                    vui := pcSlice.GetSPS().getVuiParameters();

                    if vui.getSubPicCpbParamsPresentFlag() {
                      var i	int;
                      var ui64Tmp uint64;
                      var uiTmp, uiPrev, uiCurr uint;

                      uiPrev = 0;
                      for i = 0; i < ( pictureTimingSEI.m_numDecodingUnitsMinus1 + 1 ); i ++ {
                      	if i == 0 {
                        	pictureTimingSEI.m_numNalusInDuMinus1[ i ]       =  ( accumNalsDU[ i ] ) ;
                        }else{
                        	pictureTimingSEI.m_numNalusInDuMinus1[ i ]       =  ( accumNalsDU[ i ] - accumNalsDU[ i - 1] - 1 );
                        }
                        ui64Tmp = ( ( ( accumBitsDU[ pictureTimingSEI.m_numDecodingUnitsMinus1 ] - accumBitsDU[ i ] ) * ( vui.getTimeScale() / vui.getNumUnitsInTick() ) * ( vui.getTickDivisorMinus2() + 2 ) ) /
                                     ( this.m_pcCfg.GetTargetBitrate() << 10 ) );

                        uiTmp = uint(ui64Tmp);
                        if uiTmp >= ( vui.getTickDivisorMinus2() + 2 ) {
                              uiCurr = 0;
                        }else{
                        	uiCurr = ( vui.getTickDivisorMinus2() + 2 ) - uiTmp;
            			}

                        if i == pictureTimingSEI.m_numDecodingUnitsMinus1 {
                         uiCurr = vui.getTickDivisorMinus2() + 2;
                        }
                        if uiCurr <= uiPrev {
                         uiCurr = uiPrev + 1;
            			}
                        pictureTimingSEI.m_duCpbRemovalDelayMinus1[ i ]              = (uiCurr - uiPrev) - 1;
                        uiPrev = uiCurr;
                      }
                    }
                    this.m_pcEntropyCoder.setEntropyCoder(m_pcCavlcCoder, pcSlice);
                    this.m_pcEntropyCoder.setTraceFile(this.m_pTraceFile)
                    pictureTimingSEI.m_sps = pcSlice.GetSPS();
                    this.m_seiWriter.writeSEImessage(nalu.m_Bitstream, pictureTimingSEI);
                    writeRBSPTrailingBits();

                    //??? AccessUnit::iterator it = find_if(accessUnit.begin(), accessUnit.end(), mem_fun(&NALUnit::isSlice));
                    accessUnit.insert(it, NewNALUnitEBSP(nalu));
            */
        }
        //#if L0045_NON_NESTED_SEI_RESTRICTIONS
        this.xResetNonNestedSEIPresentFlags()
        //#endif
        pcPic.GetPicYuvRec().CopyToPic(pcPicYuvRecOut)

        pcPic.SetReconMark(true)
        this.m_bFirst = false
        this.m_iNumPicCoded++
        this.m_totalCoded++
        //logging: insert a newline at end of picture period
        fmt.Printf("\n")
        //fflush(stdout);

        //delete[] pcSubstreamsOut;
    }

    /*#if !RATE_CONTROL_LAMBDA_DOMAIN
      if(m_pcCfg.GetUseRateCtrl())
      {
        this.m_pcRateCtrl.updateRCGOPStatus();
      }
    #endif
      delete pcBitstreamRedirect;

      if( accumBitsDU != NULL) delete accumBitsDU;
      if( accumNalsDU != NULL) delete accumNalsDU;

      assert ( this.m_iNumPicCoded == iNumPicRcvd );
    */
}

func (this *TEncGOP) xWriteTileLocationToSliceHeader(rNalu *OutputNALUnit, rpcBitstreamRedirect **TLibCommon.TComOutputBitstream, rpcSlice *TLibCommon.TComSlice) {
    // Byte-align
    rNalu.m_Bitstream.WriteByteAlignment() // Slice header byte-alignment

    // Perform bitstream concatenation
    if (*rpcBitstreamRedirect).GetNumberOfWrittenBits() > 0 {
        uiBitCount := (*rpcBitstreamRedirect).GetNumberOfWrittenBits()
        if (*rpcBitstreamRedirect).GetByteStreamLength() > 0 {
            pucStart := (*rpcBitstreamRedirect).GetFIFO().Front()
            uiWriteByteCount := uint(0)
            for uiWriteByteCount < (uiBitCount >> 3) {
                v := pucStart.Value.(byte)
                uiBits := uint(v)
                rNalu.m_Bitstream.Write(uiBits, 8)
                pucStart = pucStart.Next()
                uiWriteByteCount++
            }
        }
        uiBitsHeld := uint(uiBitCount & 0x07)
        for uiIdx := uint(0); uiIdx < uiBitsHeld; uiIdx++ {
            rNalu.m_Bitstream.Write(uint(((*rpcBitstreamRedirect).GetHeldBits()&(1<<(7-uiIdx)))>>(7-uiIdx)), 1)
        }
    }

    this.m_pcEntropyCoder.setBitstream(rNalu.m_Bitstream)

    //delete rpcBitstreamRedirect;
    (*rpcBitstreamRedirect) = TLibCommon.NewTComOutputBitstream()
}

func (this *TEncGOP) getGOPSize() int { return this.m_iGopSize }

func (this *TEncGOP) getListPic() *list.List { return this.m_pcListPic }

func (this *TEncGOP) printOutSummary(uiNumAllPicCoded uint) {
    //assert (uiNumAllPicCoded == this.m_gcAnalyzeAll.getNumPic());

    //--CFG_KDY
    this.m_gcAnalyzeAll.setFrmRate(float64(this.m_pcCfg.GetFrameRate()))
    this.m_gcAnalyzeI.setFrmRate(float64(this.m_pcCfg.GetFrameRate()))
    this.m_gcAnalyzeP.setFrmRate(float64(this.m_pcCfg.GetFrameRate()))
    this.m_gcAnalyzeB.setFrmRate(float64(this.m_pcCfg.GetFrameRate()))

    //-- all
    fmt.Printf("\n\nSUMMARY --------------------------------------------------------\n")
    this.m_gcAnalyzeAll.printOut("a")

    fmt.Printf("\n\nI Slices--------------------------------------------------------\n")
    this.m_gcAnalyzeI.printOut("i")

    fmt.Printf("\n\nP Slices--------------------------------------------------------\n")
    this.m_gcAnalyzeP.printOut("p")

    fmt.Printf("\n\nB Slices--------------------------------------------------------\n")
    this.m_gcAnalyzeB.printOut("b")

    //#if _SUMMARY_OUT_
    //  this.m_gcAnalyzeAll.printSummaryOut();
    //#endif
    /*#if _SUMMARY_PIC_
      this.m_gcAnalyzeI.printSummary("I");
      this.m_gcAnalyzeP.printSummary("P");
      this.m_gcAnalyzeB.printSummary("B");
    #endif
    */
    fmt.Printf("\nRVM: %.3f\n", this.xCalculateRVM())
}

func (this *TEncGOP) preLoopFilterPicAll(pcPic *TLibCommon.TComPic, ruiDist *uint64, ruiBits *uint64) {
    pcSlice := pcPic.GetSlice(pcPic.GetCurrSliceIdx())
    bCalcDist := false

    this.m_pcLoopFilter.SetCfg(this.m_pcCfg.GetLFCrossTileBoundaryFlag())

    this.m_pcLoopFilter.LoopFilterPic(pcPic)

    this.m_pcEntropyCoder.setEntropyCoder(this.m_pcEncTop.getRDGoOnSbacCoder(), pcSlice, this.m_pTraceFile)
    this.m_pcEntropyCoder.resetEntropy()
    this.m_pcEntropyCoder.setBitstream(this.m_pcBitCounter)
    pcSlice = pcPic.GetSlice(0)
    if pcSlice.GetSPS().GetUseSAO() {
        var LFCrossSliceBoundaryFlag map[int]bool //(1, true); //std::vector<Bool>
        var sliceStartAddress map[int]int         //std::vector<Int>
        sliceStartAddress[0] = 0
        sliceStartAddress[1] = int(pcPic.GetNumCUsInFrame() * pcPic.GetNumPartInCU())
        pcPic.CreateNonDBFilterInfo(sliceStartAddress, 0, LFCrossSliceBoundaryFlag, 1, true)
    }

    if pcSlice.GetSPS().GetUseSAO() {
        pcPic.DestroyNonDBFilterInfo()
    }

    this.m_pcEntropyCoder.resetEntropy()
    *ruiBits += uint64(this.m_pcEntropyCoder.getNumberOfWrittenBits())

    if !bCalcDist {
        *ruiDist = this.xFindDistortionFrame(pcPic.GetPicYuvOrg(), pcPic.GetPicYuvRec())
    }
}

func (this *TEncGOP) getSliceEncoder() *TEncSlice { return this.m_pcSliceEncoder }

func (this *TEncGOP) getNalUnitType(pocCurr int) TLibCommon.NalUnitType {
    if pocCurr == 0 {
        return TLibCommon.NAL_UNIT_CODED_SLICE_IDR
    }
    if pocCurr%int(this.m_pcCfg.GetIntraPeriod()) == 0 {
        if this.m_pcCfg.GetDecodingRefreshType() == 1 {
            return TLibCommon.NAL_UNIT_CODED_SLICE_CRA
        } else if this.m_pcCfg.GetDecodingRefreshType() == 2 {
            return TLibCommon.NAL_UNIT_CODED_SLICE_IDR
        }
    }
    if this.m_pocCRA > 0 {
        if pocCurr < this.m_pocCRA {
            // All leading pictures are being marked as TFD pictures here since current encoder uses all
            // reference pictures while encoding leading pictures. An encoder can ensure that a leading
            // picture can be still decodable when random accessing to a CRA/CRANT/BLA/BLANT picture by
            // controlling the reference pictures used for encoding that leading picture. Such a leading
            // picture need not be marked as a TFD picture.
            return TLibCommon.NAL_UNIT_CODED_SLICE_TFD
        }
    }
    return TLibCommon.NAL_UNIT_CODED_SLICE_TRAIL_R
}

func (this *TEncGOP) getLSB(poc, maxLSB int) int {
    if poc >= 0 {
        return poc % maxLSB
    }

    return (maxLSB - ((-poc) % maxLSB)) % maxLSB
}

func (this *TEncGOP) arrangeLongtermPicturesInRPS(pcSlice *TLibCommon.TComSlice, rcListPic *list.List) {
    rps := pcSlice.GetRPS()
    if rps.GetNumberOfLongtermPictures() == 0 {
        return
    }

    // Arrange long-term reference pictures in the correct order of LSB and MSB,
    // and assign values for pocLSBLT and MSB present flag
    var longtermPicsPoc, longtermPicsLSB, indices [TLibCommon.MAX_NUM_REF_PICS]int

    var longtermPicsMSB [TLibCommon.MAX_NUM_REF_PICS]int

    var mSBPresentFlag [TLibCommon.MAX_NUM_REF_PICS]bool
    //  ::memset(longtermPicsPoc, 0, sizeof(longtermPicsPoc));    // Store POC values of LTRP
    //  ::memset(longtermPicsLSB, 0, sizeof(longtermPicsLSB));    // Store POC LSB values of LTRP
    //  ::memset(longtermPicsMSB, 0, sizeof(longtermPicsMSB));    // Store POC LSB values of LTRP
    //  ::memset(indices        , 0, sizeof(indices));            // Indices to aid in tracking sorted LTRPs
    //  ::memset(mSBPresentFlag , 0, sizeof(mSBPresentFlag));     // Indicate if MSB needs to be present

    // Get the long-term reference pictures
    offset := rps.GetNumberOfNegativePictures() + rps.GetNumberOfPositivePictures()
    var i, j, ctr int
    maxPicOrderCntLSB := 1 << pcSlice.GetSPS().GetBitsForPOC()
    for i = rps.GetNumberOfPictures() - 1; i >= offset; i-- {
        longtermPicsPoc[ctr] = rps.GetPOC(i)                                        // LTRP POC
        longtermPicsLSB[ctr] = this.getLSB(longtermPicsPoc[ctr], maxPicOrderCntLSB) // LTRP POC LSB
        indices[ctr] = i
        longtermPicsMSB[ctr] = longtermPicsPoc[ctr] - longtermPicsLSB[ctr]
        ctr++
    }
    numLongPics := rps.GetNumberOfLongtermPictures()
    //assert(ctr == numLongPics);

    // Arrange pictures in decreasing order of MSB;
    for i = 0; i < numLongPics; i++ {
        for j = 0; j < numLongPics-1; j++ {
            if longtermPicsMSB[j] < longtermPicsMSB[j+1] {
                var tmp int
                tmp = longtermPicsPoc[j]
                longtermPicsPoc[j] = longtermPicsPoc[j+1]
                longtermPicsPoc[j+1] = tmp

                tmp = longtermPicsLSB[j]
                longtermPicsLSB[j] = longtermPicsLSB[j+1]
                longtermPicsLSB[j+1] = tmp

                tmp = longtermPicsMSB[j]
                longtermPicsMSB[j] = longtermPicsMSB[j+1]
                longtermPicsMSB[j+1] = tmp

                tmp = indices[j]
                indices[j] = indices[j+1]
                indices[j+1] = tmp
            }
        }
    }

    for i = 0; i < numLongPics; i++ {
        // Check if MSB present flag should be enabled.
        // Check if the buffer contains any pictures that have the same LSB.
        iterPic := rcListPic.Front()
        var pcPic *TLibCommon.TComPic
        for iterPic != nil {
            pcPic = iterPic.Value.(*TLibCommon.TComPic)
            if (this.getLSB(int(pcPic.GetPOC()), maxPicOrderCntLSB) == longtermPicsLSB[i]) && // Same LSB
                (pcPic.GetSlice(0).IsReferenced()) && // Reference picture
                (int(pcPic.GetPOC()) != longtermPicsPoc[i]) { // Not the LTRP itself
                mSBPresentFlag[i] = true
                break
            }
            iterPic = iterPic.Next()
        }
    }

    // tempArray for usedByCurr flag
    var tempArray [TLibCommon.MAX_NUM_REF_PICS]bool //::memset(tempArray, 0, sizeof(tempArray));
    for i = 0; i < numLongPics; i++ {
        tempArray[i] = rps.GetUsed(indices[i])
    }
    // Now write the final values;
    ctr = 0
    currMSB := 0
    currLSB := 0
    // currPicPoc = currMSB + currLSB
    currLSB = this.getLSB(pcSlice.GetPOC(), maxPicOrderCntLSB)
    currMSB = pcSlice.GetPOC() - currLSB

    for i = int(rps.GetNumberOfPictures()) - 1; i >= offset; i-- {
        rps.SetPOC(i, longtermPicsPoc[ctr])
        rps.SetDeltaPOC(i, -pcSlice.GetPOC()+longtermPicsPoc[ctr])
        rps.SetUsed(i, tempArray[ctr])
        rps.SetPocLSBLT(i, longtermPicsLSB[ctr])
        rps.SetDeltaPocMSBCycleLT(i, (currMSB-(longtermPicsPoc[ctr]-longtermPicsLSB[ctr]))/maxPicOrderCntLSB)
        rps.SetDeltaPocMSBPresentFlag(i, mSBPresentFlag[ctr])

        //assert(rps.GetDeltaPocMSBCycleLT(i) >= 0);   // Non-negative value
        ctr++
    }

    ctr = 1
    for i = rps.GetNumberOfPictures() - 1; i >= offset; i-- {
        for j = rps.GetNumberOfPictures() - 1 - ctr; j >= offset; j-- {
            // Here at the encoder we know that we have set the full POC value for the LTRPs, hence we
            // don't have to check the MSB present flag values for this constraint.
            //assert( rps.GetPOC(i) != rps.GetPOC(j) ); // If assert fails, LTRP entry repeated in RPS!!!
        }
        ctr++
    }
}

func (this *TEncGOP) getRateCtrl() *TEncRateCtrl { return this.m_pcRateCtrl }

func (this *TEncGOP) xInitGOP(iPOCLast, iNumPicRcvd int, rcListPic *list.List, rcListPicYuvRecOut *list.List) {
    //assert( iNumPicRcvd > 0 );
    //  Exception for the first frame
    if iPOCLast == 0 {
        this.m_iGopSize = 1
    } else {
        this.m_iGopSize = this.m_pcCfg.GetGOPSize()
    }
    //assert (m_iGopSize > 0);

    return
}

func (this *TEncGOP) xGetBuffer(rcListPic, rcListPicYuvRecOut *list.List, iNumPicRcvd, iTimeOffset int, pocCurr int) (rpcPic *TLibCommon.TComPic, rpcPicYuvRecOut *TLibCommon.TComPicYuv) {
    var i int
    //  Rec. output
    iterPicYuvRec := rcListPicYuvRecOut.Back()
    for i = 0; i < iNumPicRcvd-iTimeOffset; i++ {
        iterPicYuvRec = iterPicYuvRec.Prev()
    }
    //fmt.Printf("len=%d, iNumPicRcvd=%d, iTimeOffset=%d\n",rcListPicYuvRecOut.Len(), iNumPicRcvd, iTimeOffset);

    rpcPicYuvRecOut = iterPicYuvRec.Value.(*TLibCommon.TComPicYuv)

    //  Current pic.
    iterPic := rcListPic.Front()
    for iterPic != nil {
        rpcPic = iterPic.Value.(*TLibCommon.TComPic)
        rpcPic.SetCurrSliceIdx(0)
        if int(rpcPic.GetPOC()) == pocCurr {
            break
        }
        iterPic = iterPic.Next()
    }

    //assert (rpcPic.GetPOC() == pocCurr);

    return
}

func (this *TEncGOP) xCalculateAddPSNR(pcPic *TLibCommon.TComPic, pcPicD *TLibCommon.TComPicYuv, accessUnit *AccessUnit, dEncTime time.Duration) {
    var x, y int
    uiSSDY := uint64(0)
    uiSSDU := uint64(0)
    uiSSDV := uint64(0)

    dYPSNR := float64(0.0)
    dUPSNR := float64(0.0)
    dVPSNR := float64(0.0)

    //===== calculate PSNR =====
    pOrg := pcPic.GetPicYuvOrg().GetLumaAddr()
    pRec := pcPicD.GetLumaAddr()
    iStride := pcPicD.GetStride()

    var iWidth, iHeight int

    iWidth = pcPicD.GetWidth() - this.m_pcEncTop.m_pcEncCfg.GetPad(0)
    iHeight = pcPicD.GetHeight() - this.m_pcEncTop.m_pcEncCfg.GetPad(1)

    iSize := iWidth * iHeight

    for y = 0; y < iHeight; y++ {
        for x = 0; x < iWidth; x++ {
            iDiff := int(pOrg[x] - pRec[x])
            uiSSDY += uint64(iDiff * iDiff)
        }
        pOrg = pOrg[iStride:]
        pRec = pRec[iStride:]
    }

    iHeight >>= 1
    iWidth >>= 1
    iStride >>= 1
    pOrg = pcPic.GetPicYuvOrg().GetCbAddr()
    pRec = pcPicD.GetCbAddr()

    for y = 0; y < iHeight; y++ {
        for x = 0; x < iWidth; x++ {
            iDiff := int(pOrg[x] - pRec[x])
            uiSSDU += uint64(iDiff * iDiff)
        }
        pOrg = pOrg[iStride:]
        pRec = pRec[iStride:]
    }

    pOrg = pcPic.GetPicYuvOrg().GetCrAddr()
    pRec = pcPicD.GetCrAddr()

    for y = 0; y < iHeight; y++ {
        for x = 0; x < iWidth; x++ {
            iDiff := int(pOrg[x] - pRec[x])
            uiSSDV += uint64(iDiff * iDiff)
        }
        pOrg = pOrg[iStride:]
        pRec = pRec[iStride:]
    }

    maxvalY := 255 << uint(TLibCommon.G_bitDepthY-8)
    maxvalC := 255 << uint(TLibCommon.G_bitDepthC-8)
    fRefValueY := float64(maxvalY * maxvalY * iSize)
    fRefValueC := float64(maxvalC*maxvalC*iSize) / 4.0
    if uiSSDY != 0 {
        dYPSNR = 10.0 * math.Log10(fRefValueY/float64(uiSSDY))
    } else {
        dYPSNR = 99.99
    }
    if uiSSDU != 0 {
        dUPSNR = 10.0 * math.Log10(fRefValueC/float64(uiSSDU))
    } else {
        dUPSNR = 99.99
    }
    if uiSSDV != 0 {
        dVPSNR = 10.0 * math.Log10(fRefValueC/float64(uiSSDV))
    } else {
        dVPSNR = 99.99
    }

    /* calculate the size of the access unit, excluding:
     *  - any AnnexB contributions (start_code_prefix, zero_byte, etc.,)
     *  - SEI NAL units
     */
    //fmt.Printf("not implement yet xCalculateAddPSNR\n")

    numRBSPBytes := uint(0)
    
      for it := accessUnit.Front(); it != nil; it=it.Next() {
        nalu := it.Value.(*NALUnitEBSP)
        numRBSPBytes_nal := uint(nalu.m_Bitstream.GetByteStreamLength());
    //#if VERBOSE_RATE
        fmt.Printf("*** %d numBytesInNALunit: %d\n", nalu.GetNalUnitType(), numRBSPBytes_nal);
    //#endif
        if nalu.GetNalUnitType() != TLibCommon.NAL_UNIT_SEI && nalu.GetNalUnitType() != TLibCommon.NAL_UNIT_SEI_SUFFIX {
          numRBSPBytes += numRBSPBytes_nal;
        }
      }
    

    uibits := numRBSPBytes * 8
    this.m_vRVM_RP[len(this.m_vRVM_RP)] = int(uibits)

    //===== add PSNR =====
    this.m_gcAnalyzeAll.addResult(dYPSNR, dUPSNR, dVPSNR, float64(uibits))
    pcSlice := pcPic.GetSlice(0)
    if pcSlice.IsIntra() {
        this.m_gcAnalyzeI.addResult(dYPSNR, dUPSNR, dVPSNR, float64(uibits))
    }
    if pcSlice.IsInterP() {
        this.m_gcAnalyzeP.addResult(dYPSNR, dUPSNR, dVPSNR, float64(uibits))
    }
    if pcSlice.IsInterB() {
        this.m_gcAnalyzeB.addResult(dYPSNR, dUPSNR, dVPSNR, float64(uibits))
    }

    var c string
    if pcSlice.IsIntra() {
        c = "I"
    } else if pcSlice.IsInterP() {
        c = "P"
    } else {
        c = "B"
    }
    if !pcSlice.IsReferenced() {
        //c += 32;
        if pcSlice.IsIntra() {
            c = "i"
        } else if pcSlice.IsInterP() {
            c = "p"
        } else {
            c = "b"
        }
    }
    //#if ADAPTIVE_QP_SELECTION
    fmt.Printf("POC %4d TId: %1d ( %s-SLICE, nQP %d QP %d ) %10d bits",
        pcSlice.GetPOC(),
        pcSlice.GetTLayer(),
        c,
        pcSlice.GetSliceQpBase(),
        pcSlice.GetSliceQp(),
        uibits)
    /*#else
      printf("POC %4d TId: %1d ( %c-SLICE, QP %d ) %10d bits",
             pcSlice.GetPOC()-pcSlice.GetLastIDR(),
             pcSlice.GetTLayer(),
             c,
             pcSlice.GetSliceQp(),
             uibits );
    #endif*/

    fmt.Printf(" [Y %6.4f dB    U %6.4f dB    V %6.4f dB]", dYPSNR, dUPSNR, dVPSNR)
    fmt.Printf(" [ET %v ]", dEncTime)

    for iRefList := 0; iRefList < 2; iRefList++ {
        fmt.Printf(" [L%d ", iRefList)
        for iRefIndex := 0; iRefIndex < pcSlice.GetNumRefIdx(TLibCommon.RefPicList(iRefList)); iRefIndex++ {
            fmt.Printf("%d ", pcSlice.GetRefPOC(TLibCommon.RefPicList(iRefList), iRefIndex)-pcSlice.GetLastIDR())
        }
        fmt.Printf("]")
    }
}

func (this *TEncGOP) xFindDistortionFrame(pcPic0 *TLibCommon.TComPicYuv, pcPic1 *TLibCommon.TComPicYuv) uint64 {
    var x, y int
    pSrc0 := pcPic0.GetLumaAddr()
    pSrc1 := pcPic1.GetLumaAddr()
    uiShift := 2 * TLibCommon.DISTORTION_PRECISION_ADJUSTMENT(uint(TLibCommon.G_bitDepthY-8)).(uint)
    var iTemp int

    iStride := pcPic0.GetStride()
    iWidth := pcPic0.GetWidth()
    iHeight := pcPic0.GetHeight()

    uiTotalDiff := uint64(0)

    for y = 0; y < iHeight; y++ {
        for x = 0; x < iWidth; x++ {
            iTemp = int(pSrc0[x] - pSrc1[x])
            uiTotalDiff += uint64(iTemp*iTemp) >> uiShift
        }
        pSrc0 = pSrc0[iStride:]
        pSrc1 = pSrc1[iStride:]
    }

    uiShift = 2 * TLibCommon.DISTORTION_PRECISION_ADJUSTMENT(uint(TLibCommon.G_bitDepthC-8)).(uint)
    iHeight >>= 1
    iWidth >>= 1
    iStride >>= 1

    pSrc0 = pcPic0.GetCbAddr()
    pSrc1 = pcPic1.GetCbAddr()

    for y = 0; y < iHeight; y++ {
        for x = 0; x < iWidth; x++ {
            iTemp = int(pSrc0[x] - pSrc1[x])
            uiTotalDiff += uint64(iTemp*iTemp) >> uiShift
        }
        pSrc0 = pSrc0[iStride:]
        pSrc1 = pSrc1[iStride:]
    }

    pSrc0 = pcPic0.GetCrAddr()
    pSrc1 = pcPic1.GetCrAddr()

    for y = 0; y < iHeight; y++ {
        for x = 0; x < iWidth; x++ {
            iTemp = int(pSrc0[x] - pSrc1[x])
            uiTotalDiff += uint64(iTemp*iTemp) >> uiShift
        }
        pSrc0 = pSrc0[iStride:]
        pSrc1 = pSrc1[iStride:]
    }

    return uiTotalDiff
}

func (this *TEncGOP) xCalculateRVM() float64 {
    dRVM := float64(0)
    //fmt.Printf("not implement yet xCalculateRVM\n")
    
      if this.m_pcCfg.GetGOPSize() == 1 && this.m_pcCfg.GetIntraPeriod() != 1 && this.m_pcCfg.GetFramesToBeEncoded() > TLibCommon.RVM_VCEGAM10_M * 2 {
        // calculate RVM only for lowdelay configurations
        //std::vector<Double> 
        var vRL map[int]float64;
        var vB	map[int]float64;
        N := len(this.m_vRVM_RP);
        vRL = make(map[int]float64, N );
        vB  = make(map[int]float64, N );

        var i, j int;
        var dRavg, dBavg float64;
        vB[TLibCommon.RVM_VCEGAM10_M] = 0;
        dRavg = 0;
        dBavg = 0;
        for i = TLibCommon.RVM_VCEGAM10_M + 1 ; i < N - TLibCommon.RVM_VCEGAM10_M + 1 ; i++ {
          vRL[i] = 0;
          for j = i - TLibCommon.RVM_VCEGAM10_M ; j <= i + TLibCommon.RVM_VCEGAM10_M - 1 ; j++ {
            vRL[i] += float64(this.m_vRVM_RP[j]);
          }
          vRL[i] /= ( 2 * TLibCommon.RVM_VCEGAM10_M );
          vB[i] = vB[i-1] + float64(this.m_vRVM_RP[i]) - vRL[i];
          dRavg += float64(this.m_vRVM_RP[i]);
          dBavg += vB[i];
        }

        dRavg /= float64( N - 2 * TLibCommon.RVM_VCEGAM10_M );
        dBavg /= float64( N - 2 * TLibCommon.RVM_VCEGAM10_M );

        var dSigamB float64;
        dSigamB = 0;
        for i = TLibCommon.RVM_VCEGAM10_M + 1 ; i < N - TLibCommon.RVM_VCEGAM10_M + 1 ; i++  {
          tmp := vB[i] - dBavg;
          dSigamB += tmp * tmp;
        }
        dSigamB = math.Sqrt( dSigamB / float64( N - 2 * TLibCommon.RVM_VCEGAM10_M ) );

        f := math.Sqrt( 12.0 * ( TLibCommon.RVM_VCEGAM10_M - 1 ) / ( TLibCommon.RVM_VCEGAM10_M + 1 ) );

        dRVM = dSigamB / dRavg * f;
      }
    
    return (dRVM)
}

/*
func (this *TEncGOP)    SEIActiveParameterSets* xCreateSEIActiveParameterSets (TComSPS *sps);
func (this *TEncGOP)    SEIFramePacking*        xCreateSEIFramePacking();
func (this *TEncGOP)    SEIDisplayOrientation*  xCreateSEIDisplayOrientation();
*/
func (this *TEncGOP) xCreateLeadingSEIMessages(accessUnit *AccessUnit, sps *TLibCommon.TComSPS) {
}

//#if L0045_NON_NESTED_SEI_RESTRICTIONS
func (this *TEncGOP) xGetFirstSeiLocation(accessUnit *AccessUnit) int {
    // Find the location of the first SEI message
    //AccessUnit::iterator it;
    seiStartPos := 0
    /*for(it = accessUnit.begin(); it != accessUnit.end(); it++, seiStartPos++)
      {
         if ((*it)->isSei() || (*it)->isVcl())
         {
           break;
         }
      }
      assert(it != accessUnit.end());*/
    return seiStartPos
}
func (this *TEncGOP) xResetNonNestedSEIPresentFlags() {
    this.m_activeParameterSetSEIPresentInAU = false
    this.m_bufferingPeriodSEIPresentInAU = false
    this.m_pictureTimingSEIPresentInAU = false
}

//#endif
