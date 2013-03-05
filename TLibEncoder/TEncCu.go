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
    //"fmt"
    "gohm/TLibCommon"
    "math"
)

// ====================================================================================================================
// Class definition
// ====================================================================================================================

/// CU encoder class
type TEncCu struct {
    m_ppcBestCU    []*TLibCommon.TComDataCU ///< Best CUs in each depth
    m_ppcTempCU    []*TLibCommon.TComDataCU ///< Temporary CUs in each depth
    m_uhTotalDepth byte

    m_ppcPredYuvBest []*TLibCommon.TComYuv ///< Best Prediction Yuv for each depth
    m_ppcResiYuvBest []*TLibCommon.TComYuv ///< Best Residual Yuv for each depth
    m_ppcRecoYuvBest []*TLibCommon.TComYuv ///< Best Reconstruction Yuv for each depth
    m_ppcPredYuvTemp []*TLibCommon.TComYuv ///< Temporary Prediction Yuv for each depth
    m_ppcResiYuvTemp []*TLibCommon.TComYuv ///< Temporary Residual Yuv for each depth
    m_ppcRecoYuvTemp []*TLibCommon.TComYuv ///< Temporary Reconstruction Yuv for each depth
    m_ppcOrigYuv     []*TLibCommon.TComYuv ///< Original Yuv for each depth

    //  Data : encoder control
    m_bEncodeDQP bool

    //  Access channel
    m_pcEncCfg     *TEncCfg
    m_pcPredSearch *TEncSearch
    m_pcTrQuant    *TLibCommon.TComTrQuant
    m_pcBitCounter *TLibCommon.TComBitCounter
    m_pcRdCost     *TEncRdCost

    m_pcEntropyCoder *TEncEntropy
    m_pcCavlcCoder   *TEncCavlc
    m_pcSbacCoder    *TEncSbac
    m_pcBinCABAC     *TEncBinCABAC

    // SBAC RD
    m_pppcRDSbacCoder   [][]*TEncSbac
    m_pcRDGoOnSbacCoder *TEncSbac
    m_bUseSBACRD        bool
    m_pcRateCtrl        *TEncRateCtrl
    //#if RATE_CONTROL_LAMBDA_DOMAIN
    m_LCUPredictionSAD uint
    m_addSADDepth      int
    m_temporalSAD      int
    //#endif
}

func NewTEncCu() *TEncCu {
    return &TEncCu{}
}

/// copy parameters from encoder class
func (this *TEncCu) init(pcEncTop *TEncTop) {
    //fmt.Printf("not added yet\n");

    this.m_pcEncCfg = pcEncTop.GetEncCfg()
    this.m_pcPredSearch = pcEncTop.getPredSearch()
    this.m_pcTrQuant = pcEncTop.getTrQuant()
    this.m_pcBitCounter = pcEncTop.getBitCounter()
    this.m_pcRdCost = pcEncTop.getRdCost()

    this.m_pcEntropyCoder = pcEncTop.getEntropyCoder()
    this.m_pcCavlcCoder = pcEncTop.getCavlcCoder()
    this.m_pcSbacCoder = pcEncTop.getSbacCoder()
    this.m_pcBinCABAC = pcEncTop.getBinCABAC()

    this.m_pppcRDSbacCoder = pcEncTop.getRDSbacCoder()
    this.m_pcRDGoOnSbacCoder = pcEncTop.getRDGoOnSbacCoder()

    this.m_bUseSBACRD = pcEncTop.GetEncCfg().GetUseSBACRD()
    this.m_pcRateCtrl = pcEncTop.getRateCtrl()
}

/// create internal buffers
func (this *TEncCu) create(uhTotalDepth byte, uiMaxWidth, uiMaxHeight uint) {
    var i uint

    this.m_uhTotalDepth = uhTotalDepth + 1
    this.m_ppcBestCU = make([]*TLibCommon.TComDataCU, this.m_uhTotalDepth-1)
    this.m_ppcTempCU = make([]*TLibCommon.TComDataCU, this.m_uhTotalDepth-1)

    this.m_ppcPredYuvBest = make([]*TLibCommon.TComYuv, this.m_uhTotalDepth-1)
    this.m_ppcResiYuvBest = make([]*TLibCommon.TComYuv, this.m_uhTotalDepth-1)
    this.m_ppcRecoYuvBest = make([]*TLibCommon.TComYuv, this.m_uhTotalDepth-1)
    this.m_ppcPredYuvTemp = make([]*TLibCommon.TComYuv, this.m_uhTotalDepth-1)
    this.m_ppcResiYuvTemp = make([]*TLibCommon.TComYuv, this.m_uhTotalDepth-1)
    this.m_ppcRecoYuvTemp = make([]*TLibCommon.TComYuv, this.m_uhTotalDepth-1)
    this.m_ppcOrigYuv = make([]*TLibCommon.TComYuv, this.m_uhTotalDepth-1)

    var uiNumPartitions uint
    for i = 0; i < uint(this.m_uhTotalDepth)-1; i++ {
        uiNumPartitions = 1 << uint((uint(this.m_uhTotalDepth)-i-1)<<1)
        uiWidth := uiMaxWidth >> i
        uiHeight := uiMaxHeight >> i

        this.m_ppcBestCU[i] = TLibCommon.NewTComDataCU()
        this.m_ppcBestCU[i].Create(uiNumPartitions, uiWidth, uiHeight, false, int(uiMaxWidth>>(this.m_uhTotalDepth-1)), false)
        this.m_ppcTempCU[i] = TLibCommon.NewTComDataCU()
        this.m_ppcTempCU[i].Create(uiNumPartitions, uiWidth, uiHeight, false, int(uiMaxWidth>>(this.m_uhTotalDepth-1)), false)

        this.m_ppcPredYuvBest[i] = TLibCommon.NewTComYuv()
        this.m_ppcPredYuvBest[i].Create(uiWidth, uiHeight)
        this.m_ppcResiYuvBest[i] = TLibCommon.NewTComYuv()
        this.m_ppcResiYuvBest[i].Create(uiWidth, uiHeight)
        this.m_ppcRecoYuvBest[i] = TLibCommon.NewTComYuv()
        this.m_ppcRecoYuvBest[i].Create(uiWidth, uiHeight)

        this.m_ppcPredYuvTemp[i] = TLibCommon.NewTComYuv()
        this.m_ppcPredYuvTemp[i].Create(uiWidth, uiHeight)
        this.m_ppcResiYuvTemp[i] = TLibCommon.NewTComYuv()
        this.m_ppcResiYuvTemp[i].Create(uiWidth, uiHeight)
        this.m_ppcRecoYuvTemp[i] = TLibCommon.NewTComYuv()
        this.m_ppcRecoYuvTemp[i].Create(uiWidth, uiHeight)

        this.m_ppcOrigYuv[i] = TLibCommon.NewTComYuv()
        this.m_ppcOrigYuv[i].Create(uiWidth, uiHeight)
    }

    this.m_bEncodeDQP = false
    //#if RATE_CONTROL_LAMBDA_DOMAIN
    this.m_LCUPredictionSAD = 0
    this.m_addSADDepth = 0
    this.m_temporalSAD = 0
    //#endif

    // initialize partition order.
    rpIdx := uint(0)
    piTmp := TLibCommon.G_auiZscanToRaster[0:]
    TLibCommon.InitZscanToRaster(int(this.m_uhTotalDepth), 1, 0, piTmp, &rpIdx)
    TLibCommon.InitRasterToZscan(uiMaxWidth, uiMaxHeight, uint(this.m_uhTotalDepth))

    // initialize conversion matrix from partition index to pel
    TLibCommon.InitRasterToPelXY(uiMaxWidth, uiMaxHeight, uint(this.m_uhTotalDepth))
}

/// destroy internal buffers
func (this *TEncCu) destroy() {}

/// CU analysis function
func (this *TEncCu) compressCU(rpcCU *TLibCommon.TComDataCU) {
    // initialize CU data
    this.m_ppcBestCU[0].InitCU(rpcCU.GetPic(), rpcCU.GetAddr())
    this.m_ppcTempCU[0].InitCU(rpcCU.GetPic(), rpcCU.GetAddr())

    //#if RATE_CONTROL_LAMBDA_DOMAIN
    this.m_addSADDepth = 0
    this.m_LCUPredictionSAD = 0
    this.m_temporalSAD = 0
    //#endif

    // analysis of CU
    this.xCompressCU(&(this.m_ppcBestCU[0]), &(this.m_ppcTempCU[0]), 0, TLibCommon.SIZE_NONE)

    //#if ADAPTIVE_QP_SELECTION
    if this.m_pcEncCfg.GetUseAdaptQpSelect() {
        if rpcCU.GetSlice().GetSliceType() != TLibCommon.I_SLICE { //IIII
            this.xLcuCollectARLStats(rpcCU)
        }
    }
    //#endif
}

/// CU encoding function
func (this *TEncCu) encodeCU(pcCU *TLibCommon.TComDataCU) {
    if pcCU.GetSlice().GetPPS().GetUseDQP() {
        this.setdQPFlag(true)
    }

    // Encode CU data
    this.xEncodeCU(pcCU, 0, 0)
}

func (this *TEncCu) setBitCounter(pcBitCounter *TLibCommon.TComBitCounter) {
    this.m_pcBitCounter = pcBitCounter
}

//#if RATE_CONTROL_LAMBDA_DOMAIN
func (this *TEncCu) getLCUPredictionSAD() uint { return this.m_LCUPredictionSAD }

//#endif

func (this *TEncCu) finishCU(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint) {
    pcPic := pcCU.GetPic()
    pcSlice := pcCU.GetPic().GetSlice(pcCU.GetPic().GetCurrSliceIdx())

    //Calculate end address
    uiCUAddr := pcCU.GetSCUAddr() + uiAbsPartIdx

    uiInternalAddress := pcPic.GetPicSym().GetPicSCUAddr(pcSlice.GetSliceSegmentCurEndCUAddr()-1) % pcPic.GetNumPartInCU()
    uiExternalAddress := pcPic.GetPicSym().GetPicSCUAddr(pcSlice.GetSliceSegmentCurEndCUAddr()-1) / pcPic.GetNumPartInCU()
    uiPosX := (uiExternalAddress%pcPic.GetFrameWidthInCU())*pcCU.GetSlice().GetSPS().GetMaxCUWidth() + TLibCommon.G_auiRasterToPelX[TLibCommon.G_auiZscanToRaster[uiInternalAddress]]
    uiPosY := (uiExternalAddress/pcPic.GetFrameWidthInCU())*pcCU.GetSlice().GetSPS().GetMaxCUHeight() + TLibCommon.G_auiRasterToPelY[TLibCommon.G_auiZscanToRaster[uiInternalAddress]]
    uiWidth := pcSlice.GetSPS().GetPicWidthInLumaSamples()
    uiHeight := pcSlice.GetSPS().GetPicHeightInLumaSamples()
    for uiPosX >= uiWidth || uiPosY >= uiHeight {
        uiInternalAddress--
        uiPosX = (uiExternalAddress%pcPic.GetFrameWidthInCU())*pcCU.GetSlice().GetSPS().GetMaxCUWidth() + TLibCommon.G_auiRasterToPelX[TLibCommon.G_auiZscanToRaster[uiInternalAddress]]
        uiPosY = (uiExternalAddress/pcPic.GetFrameWidthInCU())*pcCU.GetSlice().GetSPS().GetMaxCUHeight() + TLibCommon.G_auiRasterToPelY[TLibCommon.G_auiZscanToRaster[uiInternalAddress]]
    }
    uiInternalAddress++
    if uiInternalAddress == pcCU.GetPic().GetNumPartInCU() {
        uiInternalAddress = 0
        uiExternalAddress = pcPic.GetPicSym().GetCUOrderMap(int(pcPic.GetPicSym().GetInverseCUOrderMap(int(uiExternalAddress))) + 1)
    }
    uiRealEndAddress := pcPic.GetPicSym().GetPicSCUEncOrder(uiExternalAddress*pcPic.GetNumPartInCU() + uiInternalAddress)

    // Encode slice finish
    bTerminateSlice := false
    if uiCUAddr+(pcCU.GetPic().GetNumPartInCU()>>(uiDepth<<1)) == uiRealEndAddress {
        bTerminateSlice = true
    }
    uiGranularityWidth := pcCU.GetSlice().GetSPS().GetMaxCUWidth()
    uiPosX = pcCU.GetCUPelX() + TLibCommon.G_auiRasterToPelX[TLibCommon.G_auiZscanToRaster[uiAbsPartIdx]]
    uiPosY = pcCU.GetCUPelY() + TLibCommon.G_auiRasterToPelY[TLibCommon.G_auiZscanToRaster[uiAbsPartIdx]]
    granularityBoundary := ((uiPosX+uint(pcCU.GetWidth1(uiAbsPartIdx)))%uiGranularityWidth == 0 || (uiPosX+uint(pcCU.GetWidth1(uiAbsPartIdx)) == uiWidth)) &&
        ((uiPosY+uint(pcCU.GetHeight1(uiAbsPartIdx)))%uiGranularityWidth == 0 || (uiPosY+uint(pcCU.GetHeight1(uiAbsPartIdx)) == uiHeight))

    if granularityBoundary {
        // The 1-terminating bit is added to all streams, so don't add it here when it's 1.
        if !bTerminateSlice {
            this.m_pcEntropyCoder.encodeTerminatingBit(uint(TLibCommon.B2U(bTerminateSlice)))
        }
    }

    numberOfWrittenBits := 0
    if this.m_pcBitCounter != nil {
        numberOfWrittenBits = int(this.m_pcEntropyCoder.getNumberOfWrittenBits())
    }

    // Calculate slice end IF this CU puts us over slice bit size.
    iGranularitySize := int(pcCU.GetPic().GetNumPartInCU())
    iGranularityEnd := (int(pcCU.GetSCUAddr()+uiAbsPartIdx) / iGranularitySize) * iGranularitySize
    if iGranularityEnd <= int(pcSlice.GetSliceSegmentCurStartCUAddr()) {
        iGranularityEnd += TLibCommon.MAX(int(iGranularitySize), int(pcCU.GetPic().GetNumPartInCU() >> (uiDepth << 1))).(int)
    }
    // Set slice end parameter
    if pcSlice.GetSliceMode() == TLibCommon.FIXED_NUMBER_OF_BYTES && !pcSlice.GetFinalized() && int(pcSlice.GetSliceBits())+numberOfWrittenBits > int(pcSlice.GetSliceArgument())<<3 {
        pcSlice.SetSliceSegmentCurEndCUAddr(uint(iGranularityEnd))
        pcSlice.SetSliceCurEndCUAddr(uint(iGranularityEnd))
        return
    }
    // Set dependent slice end parameter
    if pcSlice.GetSliceSegmentMode() == TLibCommon.FIXED_NUMBER_OF_BYTES && !pcSlice.GetFinalized() && int(pcSlice.GetSliceSegmentBits())+numberOfWrittenBits > int(pcSlice.GetSliceSegmentArgument())<<3 {
        pcSlice.SetSliceSegmentCurEndCUAddr(uint(iGranularityEnd))
        return
    }
    if granularityBoundary {
        pcSlice.SetSliceBits(pcSlice.GetSliceBits() + uint(numberOfWrittenBits))
        pcSlice.SetSliceSegmentBits(pcSlice.GetSliceSegmentBits() + uint(numberOfWrittenBits))
        if this.m_pcBitCounter != nil {
            this.m_pcEntropyCoder.resetBits()
        }
    }
}

//#if AMP_ENC_SPEEDUP
func (this *TEncCu) xCompressCU(rpcBestCU **TLibCommon.TComDataCU, rpcTempCU **TLibCommon.TComDataCU, uiDepth uint, eParentPartSize TLibCommon.PartSize) {
    pcPic := rpcBestCU.GetPic()

    // get Original YUV data from picture
    this.m_ppcOrigYuv[uiDepth].CopyFromPicYuv(pcPic.GetPicYuvOrg(), rpcBestCU.GetAddr(), rpcBestCU.GetZorderIdxInCU())

    // variables for fast encoder decision
    bEarlySkip := false
    bTrySplit := true
    fRD_Skip := TLibCommon.MAX_DOUBLE

    // variable for Early CU determination
    bSubBranch := true

    // variable for Cbf fast mode PU decision
    doNotBlockPu := true
    earlyDetectionSkipMode := false

    bTrySplitDQP := true

    var afCost [TLibCommon.MAX_CU_DEPTH]float64
    var aiNum [TLibCommon.MAX_CU_DEPTH]int

    if rpcBestCU.GetAddr() == 0 {
        for i := 0; i < TLibCommon.MAX_CU_DEPTH; i++ {
            afCost[i] = 0 //, sizeof( afCost ) );
            aiNum[i] = 0  //, sizeof( aiNum  ) );
        }
    }

    bBoundary := false
    uiLPelX := rpcBestCU.GetCUPelX()
    uiRPelX := uiLPelX + uint(rpcBestCU.GetWidth1(0)) - 1
    uiTPelY := rpcBestCU.GetCUPelY()
    uiBPelY := uiTPelY + uint(rpcBestCU.GetHeight1(0)) - 1

    iBaseQP := this.xComputeQP(*rpcBestCU, uiDepth)
    var iMinQP, iMaxQP int
    isAddLowestQP := false
    lowestQP := -rpcTempCU.GetSlice().GetSPS().GetQpBDOffsetY()

    if (rpcTempCU.GetSlice().GetSPS().GetMaxCUWidth() >> uiDepth) >= rpcTempCU.GetSlice().GetPPS().GetMinCuDQPSize() {
        idQP := this.m_pcEncCfg.GetMaxDeltaQP()
        iMinQP = TLibCommon.CLIP3(-rpcTempCU.GetSlice().GetSPS().GetQpBDOffsetY(), TLibCommon.MAX_QP, iBaseQP-idQP).(int)
        iMaxQP = TLibCommon.CLIP3(-rpcTempCU.GetSlice().GetSPS().GetQpBDOffsetY(), TLibCommon.MAX_QP, iBaseQP+idQP).(int)
        if (rpcTempCU.GetSlice().GetSPS().GetUseLossless()) && (lowestQP < iMinQP) && rpcTempCU.GetSlice().GetPPS().GetUseDQP() {
            isAddLowestQP = true
            iMinQP = iMinQP - 1
        }
    } else {
        iMinQP = int(rpcTempCU.GetQP1(0))
        iMaxQP = int(rpcTempCU.GetQP1(0))
    }

    //#if RATE_CONTROL_LAMBDA_DOMAIN
    if this.m_pcEncCfg.GetUseRateCtrl() {
        iMinQP = this.m_pcRateCtrl.getRCQP()
        iMaxQP = this.m_pcRateCtrl.getRCQP()
    }
    /*#else
      if this.m_pcEncCfg.GetUseRateCtrl(){
        qp := this.m_pcRateCtrl.getUnitQP();
        iMinQP  = TLibCommon.CLIP3( TLibCommon.MIN_QP, TLibCommon.MAX_QP, qp).(int);
        iMaxQP  = TLibCommon.CLIP3( TLibCommon.MIN_QP, TLibCommon.MAX_QP, qp).(int);
      }
    //#endif*/

    // If slice start or slice end is within this cu...
    pcSlice := rpcTempCU.GetPic().GetSlice(rpcTempCU.GetPic().GetCurrSliceIdx())
    bSliceStart := pcSlice.GetSliceSegmentCurStartCUAddr() > rpcTempCU.GetSCUAddr() && pcSlice.GetSliceSegmentCurStartCUAddr() < rpcTempCU.GetSCUAddr()+rpcTempCU.GetTotalNumPart()
    bSliceEnd := (pcSlice.GetSliceSegmentCurEndCUAddr() > rpcTempCU.GetSCUAddr() && pcSlice.GetSliceSegmentCurEndCUAddr() < rpcTempCU.GetSCUAddr()+rpcTempCU.GetTotalNumPart())
    bInsidePicture := (uiRPelX < rpcBestCU.GetSlice().GetSPS().GetPicWidthInLumaSamples()) && (uiBPelY < rpcBestCU.GetSlice().GetSPS().GetPicHeightInLumaSamples())
    // We need to split, so don't try these modes.
    if !bSliceEnd && !bSliceStart && bInsidePicture {
        for iQP := iMinQP; iQP <= iMaxQP; iQP++ {
            if isAddLowestQP && (iQP == iMinQP) {
                iQP = lowestQP
            }
            // variables for fast encoder decision
            bEarlySkip = false
            bTrySplit = true
            fRD_Skip = TLibCommon.MAX_DOUBLE

            rpcTempCU.InitEstData(uiDepth, iQP)

            // do inter modes, SKIP and 2Nx2N
            if rpcBestCU.GetSlice().GetSliceType() != TLibCommon.I_SLICE {
                // 2Nx2N
                if this.m_pcEncCfg.GetUseEarlySkipDetection() {
                    this.xCheckRDCostInter(rpcBestCU, rpcTempCU, TLibCommon.SIZE_2Nx2N, false)
                    rpcTempCU.InitEstData(uiDepth, iQP) //by Competition for inter_2Nx2N
                }
                // SKIP
                this.xCheckRDCostMerge2Nx2N(rpcBestCU, rpcTempCU, &earlyDetectionSkipMode) //by Merge for inter_2Nx2N
                rpcTempCU.InitEstData(uiDepth, iQP)

                // fast encoder decision for early skip
                if this.m_pcEncCfg.GetUseFastEnc() {
                    iIdx := TLibCommon.G_aucConvertToBit[rpcBestCU.GetWidth1(0)]
                    if aiNum[iIdx] > 5 && fRD_Skip < TLibCommon.EARLY_SKIP_THRES*afCost[iIdx]/float64(aiNum[iIdx]) {
                        bEarlySkip = true
                        bTrySplit = false
                    }
                }

                if !this.m_pcEncCfg.GetUseEarlySkipDetection() {
                    // 2Nx2N, NxN
                    if !bEarlySkip {
                        this.xCheckRDCostInter(rpcBestCU, rpcTempCU, TLibCommon.SIZE_2Nx2N, false)
                        rpcTempCU.InitEstData(uiDepth, iQP)
                        if this.m_pcEncCfg.GetUseCbfFastMode() {
                            doNotBlockPu = rpcBestCU.GetQtRootCbf(0) != false
                        }
                    }
                }
            }

            if (rpcTempCU.GetSlice().GetSPS().GetMaxCUWidth() >> uiDepth) >= rpcTempCU.GetSlice().GetPPS().GetMinCuDQPSize() {
                if iQP == iBaseQP {
                    bTrySplitDQP = bTrySplit
                }
            } else {
                bTrySplitDQP = bTrySplit
            }
            if isAddLowestQP && (iQP == lowestQP) {
                iQP = iMinQP
            }
        }

        //#if RATE_CONTROL_LAMBDA_DOMAIN
        if uiDepth <= uint(this.m_addSADDepth) {
            this.m_LCUPredictionSAD += uint(this.m_temporalSAD)
            this.m_addSADDepth = int(uiDepth)
        }
        //#endif

        if !earlyDetectionSkipMode {
            for iQP := iMinQP; iQP <= iMaxQP; iQP++ {
                if isAddLowestQP && (iQP == iMinQP) {
                    iQP = lowestQP
                }
                rpcTempCU.InitEstData(uiDepth, iQP)

                // do inter modes, NxN, 2NxN, and Nx2N
                if rpcBestCU.GetSlice().GetSliceType() != TLibCommon.I_SLICE {
                    // 2Nx2N, NxN
                    if !bEarlySkip {
                        if !((rpcBestCU.GetWidth1(0) == 8) && (rpcBestCU.GetHeight1(0) == 8)) {
                            if uiDepth == rpcTempCU.GetSlice().GetSPS().GetMaxCUDepth()-rpcTempCU.GetSlice().GetSPS().GetAddCUDepth() && doNotBlockPu {
                                this.xCheckRDCostInter(rpcBestCU, rpcTempCU, TLibCommon.SIZE_NxN, false)
                                rpcTempCU.InitEstData(uiDepth, iQP)
                            }
                        }
                    }

                    // 2NxN, Nx2N
                    if doNotBlockPu {
                        this.xCheckRDCostInter(rpcBestCU, rpcTempCU, TLibCommon.SIZE_Nx2N, false)
                        rpcTempCU.InitEstData(uiDepth, iQP)
                        if this.m_pcEncCfg.GetUseCbfFastMode() && rpcBestCU.GetPartitionSize1(0) == TLibCommon.SIZE_Nx2N {
                            doNotBlockPu = rpcBestCU.GetQtRootCbf(0) != false
                        }
                    }
                    if doNotBlockPu {
                        this.xCheckRDCostInter(rpcBestCU, rpcTempCU, TLibCommon.SIZE_2NxN, false)
                        rpcTempCU.InitEstData(uiDepth, iQP)
                        if this.m_pcEncCfg.GetUseCbfFastMode() && rpcBestCU.GetPartitionSize1(0) == TLibCommon.SIZE_2NxN {
                            doNotBlockPu = rpcBestCU.GetQtRootCbf(0) != false
                        }
                    }
                    //#if 1
                    //! Try AMP (SIZE_2NxnU, SIZE_2NxnD, SIZE_nLx2N, SIZE_nRx2N)
                    if pcPic.GetSlice(0).GetSPS().GetAMPAcc(uiDepth) != 0 {
                        //#if AMP_ENC_SPEEDUP
                        bTestAMP_Hor := false
                        bTestAMP_Ver := false

                        //#if AMP_MRG
                        bTestMergeAMP_Hor := false
                        bTestMergeAMP_Ver := false

                        this.deriveTestModeAMP(*rpcBestCU, eParentPartSize, &bTestAMP_Hor, &bTestAMP_Ver, &bTestMergeAMP_Hor, &bTestMergeAMP_Ver)
                        //#else
                        //          deriveTestModeAMP (rpcBestCU, eParentPartSize, bTestAMP_Hor, bTestAMP_Ver);
                        //#endif

                        //! Do horizontal AMP
                        if bTestAMP_Hor {
                            if doNotBlockPu {
                                this.xCheckRDCostInter(rpcBestCU, rpcTempCU, TLibCommon.SIZE_2NxnU, false)
                                rpcTempCU.InitEstData(uiDepth, iQP)
                                if this.m_pcEncCfg.GetUseCbfFastMode() && rpcBestCU.GetPartitionSize1(0) == TLibCommon.SIZE_2NxnU {
                                    doNotBlockPu = rpcBestCU.GetQtRootCbf(0) != false
                                }
                            }
                            if doNotBlockPu {
                                this.xCheckRDCostInter(rpcBestCU, rpcTempCU, TLibCommon.SIZE_2NxnD, false)
                                rpcTempCU.InitEstData(uiDepth, iQP)
                                if this.m_pcEncCfg.GetUseCbfFastMode() && rpcBestCU.GetPartitionSize1(0) == TLibCommon.SIZE_2NxnD {
                                    doNotBlockPu = rpcBestCU.GetQtRootCbf(0) != false
                                }
                            }
                            //#if AMP_MRG
                        } else if bTestMergeAMP_Hor {
                            if doNotBlockPu {
                                this.xCheckRDCostInter(rpcBestCU, rpcTempCU, TLibCommon.SIZE_2NxnU, true)
                                rpcTempCU.InitEstData(uiDepth, iQP)
                                if this.m_pcEncCfg.GetUseCbfFastMode() && rpcBestCU.GetPartitionSize1(0) == TLibCommon.SIZE_2NxnU {
                                    doNotBlockPu = rpcBestCU.GetQtRootCbf(0) != false
                                }
                            }
                            if doNotBlockPu {
                                this.xCheckRDCostInter(rpcBestCU, rpcTempCU, TLibCommon.SIZE_2NxnD, true)
                                rpcTempCU.InitEstData(uiDepth, iQP)
                                if this.m_pcEncCfg.GetUseCbfFastMode() && rpcBestCU.GetPartitionSize1(0) == TLibCommon.SIZE_2NxnD {
                                    doNotBlockPu = rpcBestCU.GetQtRootCbf(0) != false
                                }
                            }
                        }
                        //#endif

                        //! Do horizontal AMP
                        if bTestAMP_Ver {
                            if doNotBlockPu {
                                this.xCheckRDCostInter(rpcBestCU, rpcTempCU, TLibCommon.SIZE_nLx2N, false)
                                rpcTempCU.InitEstData(uiDepth, iQP)
                                if this.m_pcEncCfg.GetUseCbfFastMode() && rpcBestCU.GetPartitionSize1(0) == TLibCommon.SIZE_nLx2N {
                                    doNotBlockPu = rpcBestCU.GetQtRootCbf(0) != false
                                }
                            }
                            if doNotBlockPu {
                                this.xCheckRDCostInter(rpcBestCU, rpcTempCU, TLibCommon.SIZE_nRx2N, false)
                                rpcTempCU.InitEstData(uiDepth, iQP)
                            }
                            //#if AMP_MRG
                        } else if bTestMergeAMP_Ver {
                            if doNotBlockPu {
                                this.xCheckRDCostInter(rpcBestCU, rpcTempCU, TLibCommon.SIZE_nLx2N, true)
                                rpcTempCU.InitEstData(uiDepth, iQP)
                                if this.m_pcEncCfg.GetUseCbfFastMode() && rpcBestCU.GetPartitionSize1(0) == TLibCommon.SIZE_nLx2N {
                                    doNotBlockPu = rpcBestCU.GetQtRootCbf(0) != false
                                }
                            }
                            if doNotBlockPu {
                                this.xCheckRDCostInter(rpcBestCU, rpcTempCU, TLibCommon.SIZE_nRx2N, true)
                                rpcTempCU.InitEstData(uiDepth, iQP)
                            }
                        }
                        //#endif

                        /*#else
                                  this.xCheckRDCostInter( rpcBestCU, rpcTempCU, SIZE_2NxnU );
                                  rpcTempCU.InitEstData( uiDepth, iQP );
                                  this.xCheckRDCostInter( rpcBestCU, rpcTempCU, SIZE_2NxnD );
                                  rpcTempCU.InitEstData( uiDepth, iQP );
                                  this.xCheckRDCostInter( rpcBestCU, rpcTempCU, SIZE_nLx2N );
                                  rpcTempCU.InitEstData( uiDepth, iQP );

                                  this.xCheckRDCostInter( rpcBestCU, rpcTempCU, SIZE_nRx2N );
                                  rpcTempCU.InitEstData( uiDepth, iQP );

                        #endif*/
                    }
                    //#endif
                }

                /*#if !REMOVE_BURST_IPCM
                      // initialize PCM flag
                      rpcTempCU.setIPCMFlag( 0, false);
                      rpcTempCU.setIPCMFlagSubParts ( false, 0, uiDepth); //SUB_LCU_DQP
                #endif*/

                // do normal intra modes
                if !bEarlySkip {
                    // speedup for inter frames
                    if rpcBestCU.GetSlice().GetSliceType() == TLibCommon.I_SLICE ||
                        rpcBestCU.GetCbf2(0, TLibCommon.TEXT_LUMA) != 0 ||
                        rpcBestCU.GetCbf2(0, TLibCommon.TEXT_CHROMA_U) != 0 ||
                        rpcBestCU.GetCbf2(0, TLibCommon.TEXT_CHROMA_V) != 0 { // avoid very complex intra if it is unlikely
                        this.xCheckRDCostIntra(rpcBestCU, rpcTempCU, TLibCommon.SIZE_2Nx2N)
                        rpcTempCU.InitEstData(uiDepth, iQP)
                        if uiDepth == rpcTempCU.GetSlice().GetSPS().GetMaxCUDepth()-rpcTempCU.GetSlice().GetSPS().GetAddCUDepth() {
                            if rpcTempCU.GetWidth1(0) > (1 << rpcTempCU.GetSlice().GetSPS().GetQuadtreeTULog2MinSize()) {
                                this.xCheckRDCostIntra(rpcBestCU, rpcTempCU, TLibCommon.SIZE_NxN)
                                rpcTempCU.InitEstData(uiDepth, iQP)
                            }
                        }
                    }
                }

                // test PCM
                if pcPic.GetSlice(0).GetSPS().GetUsePCM() &&
                    rpcTempCU.GetWidth1(0) <= (1<<pcPic.GetSlice(0).GetSPS().GetPCMLog2MaxSize()) &&
                    rpcTempCU.GetWidth1(0) >= (1<<pcPic.GetSlice(0).GetSPS().GetPCMLog2MinSize()) {
                    uiRawBits := uint(2*TLibCommon.G_bitDepthY+TLibCommon.G_bitDepthC) * uint(rpcBestCU.GetWidth1(0)) * uint(rpcBestCU.GetHeight1(0)) / 2
                    uiBestBits := rpcBestCU.GetTotalBits()
                    if (uiBestBits > uiRawBits) || (rpcBestCU.GetTotalCost() > this.m_pcRdCost.calcRdCost(uiRawBits, 0, false, TLibCommon.DF_DEFAULT)) {
                        this.xCheckIntraPCM(rpcBestCU, rpcTempCU)
                        rpcTempCU.InitEstData(uiDepth, iQP)
                    }
                }
                if isAddLowestQP && (iQP == lowestQP) {
                    iQP = iMinQP
                }
            }

        }

        this.m_pcEntropyCoder.resetBits()
        this.m_pcEntropyCoder.encodeSplitFlag(*rpcBestCU, 0, uiDepth, true)
        rpcBestCU.SetTotalBits(rpcBestCU.GetTotalBits() + this.m_pcEntropyCoder.getNumberOfWrittenBits()) // split bits
        if this.m_pcEncCfg.GetUseSBACRD() {
            rpcBestCU.SetTotalBins(rpcBestCU.GetTotalBins() + this.m_pcEntropyCoder.m_pcEntropyCoderIf.getEncBinIf().getTEncBinCABAC().getBinsCoded())
        }
        rpcBestCU.SetTotalCost(this.m_pcRdCost.calcRdCost(rpcBestCU.GetTotalBits(), rpcBestCU.GetTotalDistortion(), false, TLibCommon.DF_DEFAULT))

        // accumulate statistics for early skip
        if this.m_pcEncCfg.GetUseFastEnc() {
            if rpcBestCU.IsSkipped(0) {
                iIdx := TLibCommon.G_aucConvertToBit[rpcBestCU.GetWidth1(0)]
                afCost[iIdx] += rpcBestCU.GetTotalCost()
                aiNum[iIdx]++
            }
        }

        // Early CU determination
        if this.m_pcEncCfg.GetUseEarlyCU() && rpcBestCU.IsSkipped(0) {
            bSubBranch = false
        } else {
            bSubBranch = true
        }
    } else if !(bSliceEnd && bInsidePicture) {
        bBoundary = true
        //#if RATE_CONTROL_LAMBDA_DOMAIN
        this.m_addSADDepth++
        //#endif
    }

    // copy orginal YUV samples to PCM buffer
    if rpcBestCU.IsLosslessCoded(0) && (rpcBestCU.GetIPCMFlag1(0) == false) {
        this.xFillPCMBuffer(*rpcBestCU, this.m_ppcOrigYuv[uiDepth])
    }
    if (rpcTempCU.GetSlice().GetSPS().GetMaxCUWidth() >> uiDepth) == rpcTempCU.GetSlice().GetPPS().GetMinCuDQPSize() {
        idQP := this.m_pcEncCfg.GetMaxDeltaQP()
        iMinQP = TLibCommon.CLIP3(-rpcTempCU.GetSlice().GetSPS().GetQpBDOffsetY(), TLibCommon.MAX_QP, iBaseQP-idQP).(int)
        iMaxQP = TLibCommon.CLIP3(-rpcTempCU.GetSlice().GetSPS().GetQpBDOffsetY(), TLibCommon.MAX_QP, iBaseQP+idQP).(int)
        if (rpcTempCU.GetSlice().GetSPS().GetUseLossless()) && (lowestQP < iMinQP) && rpcTempCU.GetSlice().GetPPS().GetUseDQP() {
            isAddLowestQP = true
            iMinQP = iMinQP - 1
        }
    } else if (rpcTempCU.GetSlice().GetSPS().GetMaxCUWidth() >> uiDepth) > rpcTempCU.GetSlice().GetPPS().GetMinCuDQPSize() {
        iMinQP = iBaseQP
        iMaxQP = iBaseQP
    } else {
        var iStartQP int
        if pcPic.GetCU(rpcTempCU.GetAddr()).GetSliceSegmentStartCU(rpcTempCU.GetZorderIdxInCU()) == pcSlice.GetSliceSegmentCurStartCUAddr() {
            iStartQP = int(rpcTempCU.GetQP1(0))
        } else {
            uiCurSliceStartPartIdx := pcSlice.GetSliceSegmentCurStartCUAddr()%pcPic.GetNumPartInCU() - rpcTempCU.GetZorderIdxInCU()
            iStartQP = int(rpcTempCU.GetQP1(uiCurSliceStartPartIdx))
        }
        iMinQP = iStartQP
        iMaxQP = iStartQP
    }
    //#if RATE_CONTROL_LAMBDA_DOMAIN
    if this.m_pcEncCfg.GetUseRateCtrl() {
        iMinQP = this.m_pcRateCtrl.getRCQP()
        iMaxQP = this.m_pcRateCtrl.getRCQP()
    }
    /*#else
      if(this.m_pcEncCfg.GetUseRateCtrl())
      {
        Int qp = this.m_pcRateCtrl.GetUnitQP();
        iMinQP  = Clip3( MIN_QP, MAX_QP, qp);
        iMaxQP  = Clip3( MIN_QP, MAX_QP, qp);
      }
    #endif*/
    for iQP := iMinQP; iQP <= iMaxQP; iQP++ {
        if isAddLowestQP && (iQP == iMinQP) {
            iQP = lowestQP
        }
        rpcTempCU.InitEstData(uiDepth, iQP)

        // further split
        if bSubBranch && bTrySplitDQP && uiDepth < rpcTempCU.GetSlice().GetSPS().GetMaxCUDepth()-rpcTempCU.GetSlice().GetSPS().GetAddCUDepth() {
            uhNextDepth := uiDepth + 1
            pcSubBestPartCU := this.m_ppcBestCU[uhNextDepth]
            pcSubTempPartCU := this.m_ppcTempCU[uhNextDepth]

            for uiPartUnitIdx := uint(0); uiPartUnitIdx < 4; uiPartUnitIdx++ {
                pcSubBestPartCU.InitSubCU(*rpcTempCU, uiPartUnitIdx, uhNextDepth, iQP) // clear sub partition datas or init.
                pcSubTempPartCU.InitSubCU(*rpcTempCU, uiPartUnitIdx, uhNextDepth, iQP) // clear sub partition datas or init.

                bInSlice := pcSubBestPartCU.GetSCUAddr()+pcSubBestPartCU.GetTotalNumPart() > pcSlice.GetSliceSegmentCurStartCUAddr() && pcSubBestPartCU.GetSCUAddr() < pcSlice.GetSliceSegmentCurEndCUAddr()
                if bInSlice && (pcSubBestPartCU.GetCUPelX() < pcSlice.GetSPS().GetPicWidthInLumaSamples()) && (pcSubBestPartCU.GetCUPelY() < pcSlice.GetSPS().GetPicHeightInLumaSamples()) {
                    if this.m_bUseSBACRD {
                        if 0 == uiPartUnitIdx { //initialize RD with previous depth buffer
                            this.m_pppcRDSbacCoder[uhNextDepth][TLibCommon.CI_CURR_BEST].load(this.m_pppcRDSbacCoder[uiDepth][TLibCommon.CI_CURR_BEST])
                        } else {
                            this.m_pppcRDSbacCoder[uhNextDepth][TLibCommon.CI_CURR_BEST].load(this.m_pppcRDSbacCoder[uhNextDepth][TLibCommon.CI_NEXT_BEST])
                        }
                    }

                    //#if AMP_ENC_SPEEDUP
                    if rpcBestCU.IsIntra(0) {
                        this.xCompressCU(&pcSubBestPartCU, &pcSubTempPartCU, uhNextDepth, TLibCommon.SIZE_NONE)
                    } else {
                        this.xCompressCU(&pcSubBestPartCU, &pcSubTempPartCU, uhNextDepth, rpcBestCU.GetPartitionSize1(0))
                    }
                    //#else
                    //          xCompressCU( pcSubBestPartCU, pcSubTempPartCU, uhNextDepth );
                    //#endif

                    rpcTempCU.CopyPartFrom(pcSubBestPartCU, uiPartUnitIdx, uhNextDepth) // Keep best part data to current temporary data.
                    this.xCopyYuv2Tmp(pcSubBestPartCU.GetTotalNumPart()*uiPartUnitIdx, uhNextDepth)
                } else if bInSlice {
                    pcSubBestPartCU.CopyToPic1(uhNextDepth)
                    rpcTempCU.CopyPartFrom(pcSubBestPartCU, uiPartUnitIdx, uhNextDepth)
                }
            }

            if !bBoundary {
                this.m_pcEntropyCoder.resetBits()
                this.m_pcEntropyCoder.encodeSplitFlag(*rpcTempCU, 0, uiDepth, true)

                rpcTempCU.SetTotalBits(rpcTempCU.GetTotalBits() + this.m_pcEntropyCoder.getNumberOfWrittenBits()) // split bits
                if this.m_pcEncCfg.GetUseSBACRD() {
                    rpcTempCU.SetTotalBins(rpcTempCU.GetTotalBins() + this.m_pcEntropyCoder.m_pcEntropyCoderIf.getEncBinIf().getTEncBinCABAC().getBinsCoded())
                }
            }
            rpcTempCU.SetTotalCost(this.m_pcRdCost.calcRdCost(rpcTempCU.GetTotalBits(), rpcTempCU.GetTotalDistortion(), false, TLibCommon.DF_DEFAULT))

            if (rpcTempCU.GetSlice().GetSPS().GetMaxCUWidth()>>uiDepth) == rpcTempCU.GetSlice().GetPPS().GetMinCuDQPSize() && rpcTempCU.GetSlice().GetPPS().GetUseDQP() {
                hasResidual := false
                for uiBlkIdx := uint(0); uiBlkIdx < rpcTempCU.GetTotalNumPart(); uiBlkIdx++ {
                    if (pcPic.GetCU(rpcTempCU.GetAddr()).GetSliceSegmentStartCU(uiBlkIdx+rpcTempCU.GetZorderIdxInCU()) == rpcTempCU.GetSlice().GetSliceSegmentCurStartCUAddr()) &&
                        (rpcTempCU.GetCbf2(uiBlkIdx, TLibCommon.TEXT_LUMA) != 0 || rpcTempCU.GetCbf2(uiBlkIdx, TLibCommon.TEXT_CHROMA_U) != 0 || rpcTempCU.GetCbf2(uiBlkIdx, TLibCommon.TEXT_CHROMA_V) != 0) {
                        hasResidual = true
                        break
                    }
                }

                var uiTargetPartIdx uint
                if pcPic.GetCU(rpcTempCU.GetAddr()).GetSliceSegmentStartCU(rpcTempCU.GetZorderIdxInCU()) != pcSlice.GetSliceSegmentCurStartCUAddr() {
                    uiTargetPartIdx = pcSlice.GetSliceSegmentCurStartCUAddr()%pcPic.GetNumPartInCU() - rpcTempCU.GetZorderIdxInCU()
                } else {
                    uiTargetPartIdx = 0
                }
                if hasResidual {
                    //#if !RDO_WITHOUT_DQP_BITS
                    this.m_pcEntropyCoder.resetBits()
                    this.m_pcEntropyCoder.encodeQP(*rpcTempCU, uiTargetPartIdx, false)
                    rpcTempCU.SetTotalBits(rpcTempCU.GetTotalBits() + this.m_pcEntropyCoder.getNumberOfWrittenBits()) // dQP bits
                    if this.m_pcEncCfg.GetUseSBACRD() {
                        rpcTempCU.SetTotalBins(rpcTempCU.GetTotalBins() + this.m_pcEntropyCoder.m_pcEntropyCoderIf.getEncBinIf().getTEncBinCABAC().getBinsCoded())
                    }
                    rpcTempCU.SetTotalCost(this.m_pcRdCost.calcRdCost(rpcTempCU.GetTotalBits(), rpcTempCU.GetTotalDistortion(), false, TLibCommon.DF_DEFAULT))
                    //#endif

                    foundNonZeroCbf := false
                    rpcTempCU.SetQPSubCUs(int((*rpcTempCU).GetRefQP(uiTargetPartIdx)), *rpcTempCU, 0, uiDepth, &foundNonZeroCbf)
                    //assert( foundNonZeroCbf );
                } else {
                    rpcTempCU.SetQPSubParts(int((*rpcTempCU).GetRefQP(uiTargetPartIdx)), 0, uiDepth) // set QP to default QP
                }
            }

            if this.m_bUseSBACRD {
                this.m_pppcRDSbacCoder[uhNextDepth][TLibCommon.CI_NEXT_BEST].store(this.m_pppcRDSbacCoder[uiDepth][TLibCommon.CI_TEMP_BEST])
            }
            isEndOfSlice := rpcBestCU.GetSlice().GetSliceMode() == TLibCommon.FIXED_NUMBER_OF_BYTES &&
                (rpcBestCU.GetTotalBits() > rpcBestCU.GetSlice().GetSliceArgument()<<3)
            isEndOfSliceSegment := rpcBestCU.GetSlice().GetSliceSegmentMode() == TLibCommon.FIXED_NUMBER_OF_BYTES &&
                (rpcBestCU.GetTotalBits() > rpcBestCU.GetSlice().GetSliceSegmentArgument()<<3)
            if isEndOfSlice || isEndOfSliceSegment {
                rpcBestCU.SetTotalCost(rpcTempCU.GetTotalCost() + 1)
            }
            this.xCheckBestMode3(rpcBestCU, rpcTempCU, uiDepth) // RD compare current larger prediction
        }   // with sub partitioned prediction.
        if isAddLowestQP && (iQP == lowestQP) {
            iQP = iMinQP
        }
    }

    rpcBestCU.CopyToPic1(uiDepth) // Copy Best data to Picture for next partition prediction.

    this.xCopyYuv2Pic((*rpcBestCU).GetPic(), (*rpcBestCU).GetAddr(), (*rpcBestCU).GetZorderIdxInCU(), uiDepth, uiDepth, *rpcBestCU, uiLPelX, uiTPelY) // Copy Yuv data to picture Yuv
    if bBoundary || (bSliceEnd && bInsidePicture) {
        return
    }

    // Assert if Best prediction mode is NONE
    // Selected mode's RD-cost must be not MAX_DOUBLE.
    //assert( rpcBestCU.GetPartitionSize ( 0 ) != SIZE_NONE  );
    //assert( rpcBestCU.GetPredictionMode( 0 ) != MODE_NONE  );
    //assert( rpcBestCU.GetTotalCost     (   ) != MAX_DOUBLE );
}

//#else
//func (this *TEncCu)  xCompressCU         ( TLibCommon.TComDataCU*& rpcBestCU, TLibCommon.TComDataCU*& rpcTempCU, UInt uiDepth        );
//#endif
func (this *TEncCu) xEncodeCU(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint) {
    //fmt.Printf("Enter xEncodeCU with %d\n", pcCU.GetAddr());
    
    pcPic := pcCU.GetPic()

    bBoundary := false
    uiLPelX := pcCU.GetCUPelX() + TLibCommon.G_auiRasterToPelX[TLibCommon.G_auiZscanToRaster[uiAbsPartIdx]]
    uiRPelX := uiLPelX + (pcCU.GetSlice().GetSPS().GetMaxCUWidth() >> uiDepth) - 1
    uiTPelY := pcCU.GetCUPelY() + TLibCommon.G_auiRasterToPelY[TLibCommon.G_auiZscanToRaster[uiAbsPartIdx]]
    uiBPelY := uiTPelY + (pcCU.GetSlice().GetSPS().GetMaxCUHeight() >> uiDepth) - 1

    pcSlice := pcCU.GetPic().GetSlice(pcCU.GetPic().GetCurrSliceIdx())
    // If slice start is within this cu...
    bSliceStart := pcSlice.GetSliceSegmentCurStartCUAddr() > pcPic.GetPicSym().GetInverseCUOrderMap(int(pcCU.GetAddr()))*pcCU.GetPic().GetNumPartInCU()+uiAbsPartIdx &&
        pcSlice.GetSliceSegmentCurStartCUAddr() < pcPic.GetPicSym().GetInverseCUOrderMap(int(pcCU.GetAddr()))*pcCU.GetPic().GetNumPartInCU()+uiAbsPartIdx+(pcPic.GetNumPartInCU()>>(uiDepth<<1))
    // We need to split, so don't try these modes.
    if !bSliceStart && (uiRPelX < pcSlice.GetSPS().GetPicWidthInLumaSamples()) && (uiBPelY < pcSlice.GetSPS().GetPicHeightInLumaSamples()) {
        this.m_pcEntropyCoder.encodeSplitFlag(pcCU, uiAbsPartIdx, uiDepth, false)
    } else {
        bBoundary = true
    }

    if ((uiDepth < uint(pcCU.GetDepth1(uiAbsPartIdx))) && (uiDepth < (pcCU.GetSlice().GetSPS().GetMaxCUDepth() - pcCU.GetSlice().GetSPS().GetAddCUDepth()))) || bBoundary {
        uiQNumParts := (pcPic.GetNumPartInCU() >> (uiDepth << 1)) >> 2
        if (pcCU.GetSlice().GetSPS().GetMaxCUWidth()>>uiDepth) == pcCU.GetSlice().GetPPS().GetMinCuDQPSize() && pcCU.GetSlice().GetPPS().GetUseDQP() {
            this.setdQPFlag(true)
        }

        for uiPartUnitIdx := uint(0); uiPartUnitIdx < 4; uiPartUnitIdx++ {
            uiLPelX = pcCU.GetCUPelX() + TLibCommon.G_auiRasterToPelX[TLibCommon.G_auiZscanToRaster[uiAbsPartIdx]]
            uiTPelY = pcCU.GetCUPelY() + TLibCommon.G_auiRasterToPelY[TLibCommon.G_auiZscanToRaster[uiAbsPartIdx]]
            bInSlice := pcCU.GetSCUAddr()+uiAbsPartIdx+uiQNumParts > pcSlice.GetSliceSegmentCurStartCUAddr() && pcCU.GetSCUAddr()+uiAbsPartIdx < pcSlice.GetSliceSegmentCurEndCUAddr()
            if bInSlice && (uiLPelX < pcSlice.GetSPS().GetPicWidthInLumaSamples()) && (uiTPelY < pcSlice.GetSPS().GetPicHeightInLumaSamples()) {
                this.xEncodeCU(pcCU, uiAbsPartIdx, uiDepth+1)
            }

            uiAbsPartIdx += uiQNumParts
        }
        //fmt.Printf("Exit xEncodeCU\n");
        return
    }

    if (pcCU.GetSlice().GetSPS().GetMaxCUWidth()>>uiDepth) >= pcCU.GetSlice().GetPPS().GetMinCuDQPSize() && pcCU.GetSlice().GetPPS().GetUseDQP() {
        this.setdQPFlag(true)
    }
    if pcCU.GetSlice().GetPPS().GetTransquantBypassEnableFlag() {
        this.m_pcEntropyCoder.encodeCUTransquantBypassFlag(pcCU, uiAbsPartIdx, false)
    }
    if !pcCU.GetSlice().IsIntra() {
        this.m_pcEntropyCoder.encodeSkipFlag(pcCU, uiAbsPartIdx, false)
    }

    if pcCU.IsSkipped(uiAbsPartIdx) {
        this.m_pcEntropyCoder.encodeMergeIndex(pcCU, uiAbsPartIdx, false)
        this.finishCU(pcCU, uiAbsPartIdx, uiDepth)
        
        //fmt.Printf("Exit xEncodeCU\n");
        
        return
    }
    this.m_pcEntropyCoder.encodePredMode(pcCU, uiAbsPartIdx, false)

    this.m_pcEntropyCoder.encodePartSize(pcCU, uiAbsPartIdx, uiDepth, false)

    if pcCU.IsIntra(uiAbsPartIdx) && pcCU.GetPartitionSize1(uiAbsPartIdx) == TLibCommon.SIZE_2Nx2N {
        this.m_pcEntropyCoder.encodeIPCMInfo(pcCU, uiAbsPartIdx, false)

        if pcCU.GetIPCMFlag1(uiAbsPartIdx) {
            // Encode slice finish
            this.finishCU(pcCU, uiAbsPartIdx, uiDepth)
            
            //fmt.Printf("Exit xEncodeCU\n");
            
            return
        }
    }

    // prediction Info ( Intra : direction mode, Inter : Mv, reference idx )
    this.m_pcEntropyCoder.encodePredInfo(pcCU, uiAbsPartIdx, false)

    // Encode Coefficients
    bCodeDQP := this.getdQPFlag()
    this.m_pcEntropyCoder.encodeCoeff(pcCU, uiAbsPartIdx, uiDepth, uint(pcCU.GetWidth1(uiAbsPartIdx)), uint(pcCU.GetHeight1(uiAbsPartIdx)), &bCodeDQP)
    this.setdQPFlag(bCodeDQP)

    // --- write terminating bit ---
    this.finishCU(pcCU, uiAbsPartIdx, uiDepth)
    
    //fmt.Printf("Exit xEncodeCU\n");
}

func (this *TEncCu) xComputeQP(pcCU *TLibCommon.TComDataCU, uiDepth uint) int {
    iBaseQp := pcCU.GetSlice().GetSliceQp()
    iQpOffset := 0
    if this.m_pcEncCfg.GetUseAdaptiveQP() {
        pcEPic := pcCU.GetPic()
        uiAQDepth := uint(TLibCommon.MIN(int(uiDepth), int(pcEPic.GetMaxAQDepth()-1)).(int))
        pcAQLayer := pcEPic.GetAQLayer(uiAQDepth)
        uiAQUPosX := pcCU.GetCUPelX() / pcAQLayer.GetAQPartWidth()
        uiAQUPosY := pcCU.GetCUPelY() / pcAQLayer.GetAQPartHeight()
        uiAQUStride := pcAQLayer.GetAQPartStride()
        acAQU := pcAQLayer.GetQPAdaptationUnit()

        dMaxQScale := math.Pow(2.0, float64(this.m_pcEncCfg.GetQPAdaptationRange())/6.0)
        dAvgAct := pcAQLayer.GetAvgActivity()
        dCUAct := acAQU[uiAQUPosY*uiAQUStride+uiAQUPosX].GetActivity()
        dNormAct := (dMaxQScale*dCUAct + dAvgAct) / (dCUAct + dMaxQScale*dAvgAct)
        dQpOffset := math.Log(dNormAct) / math.Log(2.0) * 6.0
        iQpOffset = int(math.Floor(dQpOffset + 0.49999))
    }
    return TLibCommon.CLIP3(-pcCU.GetSlice().GetSPS().GetQpBDOffsetY(), TLibCommon.MAX_QP, iBaseQp+iQpOffset).(int)
}

func (this *TEncCu) xCheckBestMode3(rpcBestCU **TLibCommon.TComDataCU, rpcTempCU **TLibCommon.TComDataCU, uiDepth uint) {
    if rpcTempCU.GetTotalCost() < rpcBestCU.GetTotalCost() {
        var pcYuv *TLibCommon.TComYuv
        // Change Information data
        pcCU := *rpcBestCU
        *rpcBestCU = *rpcTempCU
        *rpcTempCU = pcCU

        // Change Prediction data
        pcYuv = this.m_ppcPredYuvBest[uiDepth]
        this.m_ppcPredYuvBest[uiDepth] = this.m_ppcPredYuvTemp[uiDepth]
        this.m_ppcPredYuvTemp[uiDepth] = pcYuv

        // Change Reconstruction data
        pcYuv = this.m_ppcRecoYuvBest[uiDepth]
        this.m_ppcRecoYuvBest[uiDepth] = this.m_ppcRecoYuvTemp[uiDepth]
        this.m_ppcRecoYuvTemp[uiDepth] = pcYuv

        pcYuv = nil
        pcCU = nil

        if this.m_bUseSBACRD { // store temp best CI for next CU coding
            this.m_pppcRDSbacCoder[uiDepth][TLibCommon.CI_TEMP_BEST].store(this.m_pppcRDSbacCoder[uiDepth][TLibCommon.CI_NEXT_BEST])
        }
    }
}

func (this *TEncCu) xCheckRDCostMerge2Nx2N(rpcBestCU **TLibCommon.TComDataCU, rpcTempCU **TLibCommon.TComDataCU, earlyDetectionSkipMode *bool) {
    //assert( rpcTempCU.GetSlice().GetSliceType() != TLibCommon.I_SLICE );
    var cMvFieldNeighbours [TLibCommon.MRG_MAX_NUM_CANDS << 1]TLibCommon.TComMvField // double length for mv of both lists
    var uhInterDirNeighbours [TLibCommon.MRG_MAX_NUM_CANDS]byte
    numValidMergeCand := 0

    for ui := uint(0); ui < rpcTempCU.GetSlice().GetMaxNumMergeCand(); ui++ {
        uhInterDirNeighbours[ui] = 0
    }
    uhDepth := uint(rpcTempCU.GetDepth1(0))
    rpcTempCU.SetPartSizeSubParts(TLibCommon.SIZE_2Nx2N, 0, uhDepth) // interprets depth relative to LCU level
    rpcTempCU.SetCUTransquantBypassSubParts(this.m_pcEncCfg.GetCUTransquantBypassFlagValue(), 0, uhDepth)
    rpcTempCU.GetInterMergeCandidates(0, 0, cMvFieldNeighbours[:], uhInterDirNeighbours[:], &numValidMergeCand, -1)

    var mergeCandBuffer [TLibCommon.MRG_MAX_NUM_CANDS]int
    for ui := uint(0); ui < rpcTempCU.GetSlice().GetMaxNumMergeCand(); ui++ {
        mergeCandBuffer[ui] = 0
    }

    bestIsSkip := false

    var iteration uint
    if rpcTempCU.IsLosslessCoded(0) {
        iteration = 1
    } else {
        iteration = 2
    }

    for uiNoResidual := uint(0); uiNoResidual < iteration; uiNoResidual++ {
        for uiMergeCand := uint(0); uiMergeCand < uint(numValidMergeCand); uiMergeCand++ {
            {
                if !(uiNoResidual == 1 && mergeCandBuffer[uiMergeCand] == 1) {

                    if !(bestIsSkip && uiNoResidual == 0) {
                        // set MC parameters
                        rpcTempCU.SetPredModeSubParts(TLibCommon.MODE_INTER, 0, uhDepth) // interprets depth relative to LCU level
                        rpcTempCU.SetCUTransquantBypassSubParts(this.m_pcEncCfg.GetCUTransquantBypassFlagValue(), 0, uhDepth)
                        rpcTempCU.SetPartSizeSubParts(TLibCommon.SIZE_2Nx2N, 0, uhDepth)                                                                      // interprets depth relative to LCU level
                        rpcTempCU.SetMergeFlagSubParts(true, 0, 0, uhDepth)                                                                                   // interprets depth relative to LCU level
                        rpcTempCU.SetMergeIndexSubParts(uiMergeCand, 0, 0, uhDepth)                                                                           // interprets depth relative to LCU level
                        rpcTempCU.SetInterDirSubParts(uint(uhInterDirNeighbours[uiMergeCand]), 0, 0, uhDepth)                                                 // interprets depth relative to LCU level
                        rpcTempCU.GetCUMvField(TLibCommon.REF_PIC_LIST_0).SetAllMvField(&cMvFieldNeighbours[0+2*uiMergeCand], TLibCommon.SIZE_2Nx2N, 0, 0, 0) // interprets depth relative to rpcTempCU level
                        rpcTempCU.GetCUMvField(TLibCommon.REF_PIC_LIST_1).SetAllMvField(&cMvFieldNeighbours[1+2*uiMergeCand], TLibCommon.SIZE_2Nx2N, 0, 0, 0) // interprets depth relative to rpcTempCU level

                        // do MC
                        this.m_pcPredSearch.MotionCompensation(*rpcTempCU, this.m_ppcPredYuvTemp[uhDepth], TLibCommon.REF_PIC_LIST_X, -1)
                        // estimate residual and encode everything
                        this.m_pcPredSearch.encodeResAndCalcRdInterCU(*rpcTempCU,
                            this.m_ppcOrigYuv[uhDepth],
                            this.m_ppcPredYuvTemp[uhDepth],
                            this.m_ppcResiYuvTemp[uhDepth],
                            this.m_ppcResiYuvBest[uhDepth],
                            this.m_ppcRecoYuvTemp[uhDepth],
                            uiNoResidual != 0)

                        if uiNoResidual == 0 {
                            if rpcTempCU.GetQtRootCbf(0) == false {
                                mergeCandBuffer[uiMergeCand] = 1
                            }
                        }

                        rpcTempCU.SetSkipFlagSubParts(rpcTempCU.GetQtRootCbf(0) == false, 0, uhDepth)
                        orgQP := int(rpcTempCU.GetQP1(0))
                        this.xCheckDQP(*rpcTempCU)
                        this.xCheckBestMode3(rpcBestCU, rpcTempCU, uhDepth)
                        rpcTempCU.InitEstData(uhDepth, orgQP)

                        if this.m_pcEncCfg.GetUseFastDecisionForMerge() && !bestIsSkip {
                            bestIsSkip = rpcBestCU.GetQtRootCbf(0) == false
                        }

                    }
                }
            }
        }

        if uiNoResidual == 0 && this.m_pcEncCfg.GetUseEarlySkipDetection() {
            if rpcBestCU.GetQtRootCbf(0) == false {
                if rpcBestCU.GetMergeFlag1(0) {
                    *earlyDetectionSkipMode = true
                } else {
                    absoulte_MV := 0
                    for uiRefListIdx := uint(0); uiRefListIdx < 2; uiRefListIdx++ {
                        if rpcBestCU.GetSlice().GetNumRefIdx(TLibCommon.RefPicList(uiRefListIdx)) > 0 {
                            pcCUMvField := rpcBestCU.GetCUMvField(TLibCommon.RefPicList(uiRefListIdx))
                            mvd := pcCUMvField.GetMvd(0)
                            iHor := int(mvd.GetAbsHor())
                            iVer := int(mvd.GetAbsVer())
                            absoulte_MV += iHor + iVer
                        }
                    }

                    if absoulte_MV == 0 {
                        *earlyDetectionSkipMode = true
                    }
                }
            }
        }
    }
}

//#if AMP_MRG
func (this *TEncCu) xCheckRDCostInter(rpcBestCU **TLibCommon.TComDataCU, rpcTempCU **TLibCommon.TComDataCU, ePartSize TLibCommon.PartSize, bUseMRG bool) {
    uhDepth := uint(rpcTempCU.GetDepth1(0))

    rpcTempCU.SetDepthSubParts(uhDepth, 0)

    rpcTempCU.SetSkipFlagSubParts(false, 0, uhDepth)

    rpcTempCU.SetPartSizeSubParts(ePartSize, 0, uhDepth)
    rpcTempCU.SetPredModeSubParts(TLibCommon.MODE_INTER, 0, uhDepth)
    rpcTempCU.SetCUTransquantBypassSubParts(this.m_pcEncCfg.GetCUTransquantBypassFlagValue(), 0, uhDepth)

    //#if AMP_MRG
    rpcTempCU.SetMergeAMP(true)
    this.m_pcPredSearch.predInterSearch(*rpcTempCU, this.m_ppcOrigYuv[uhDepth], this.m_ppcPredYuvTemp[uhDepth], this.m_ppcResiYuvTemp[uhDepth], this.m_ppcRecoYuvTemp[uhDepth], false, bUseMRG)
    //#else
    //  this.m_pcPredSearch.predInterSearch ( rpcTempCU, this.m_ppcOrigYuv[uhDepth], this.m_ppcPredYuvTemp[uhDepth], this.m_ppcResiYuvTemp[uhDepth], this.m_ppcRecoYuvTemp[uhDepth] );
    //#endif

    //#if AMP_MRG
    if !rpcTempCU.GetMergeAMP() {
        return
    }
    //#endif

    //#if RATE_CONTROL_LAMBDA_DOMAIN
    if this.m_pcEncCfg.GetUseRateCtrl() && this.m_pcEncCfg.GetLCULevelRC() && ePartSize == TLibCommon.SIZE_2Nx2N && uhDepth <= uint(this.m_addSADDepth) {
        SAD := this.m_pcRdCost.getSADPart(TLibCommon.G_bitDepthY, this.m_ppcPredYuvTemp[uhDepth].GetLumaAddr(), int(this.m_ppcPredYuvTemp[uhDepth].GetStride()),
            this.m_ppcOrigYuv[uhDepth].GetLumaAddr(), int(this.m_ppcOrigYuv[uhDepth].GetStride()),
            int(rpcTempCU.GetWidth1(0)), int(rpcTempCU.GetHeight1(0)))
        this.m_temporalSAD = int(SAD)
    }
    //#endif

    this.m_pcPredSearch.encodeResAndCalcRdInterCU(*rpcTempCU, this.m_ppcOrigYuv[uhDepth], this.m_ppcPredYuvTemp[uhDepth], this.m_ppcResiYuvTemp[uhDepth], this.m_ppcResiYuvBest[uhDepth], this.m_ppcRecoYuvTemp[uhDepth], false)
    rpcTempCU.SetTotalCost(this.m_pcRdCost.calcRdCost(rpcTempCU.GetTotalBits(), rpcTempCU.GetTotalDistortion(), false, TLibCommon.DF_DEFAULT))

    this.xCheckDQP(*rpcTempCU)
    this.xCheckBestMode3(rpcBestCU, rpcTempCU, uhDepth)
}

//#else
//func (this *TEncCu)  xCheckRDCostInter   ( TLibCommon.TComDataCU*& rpcBestCU, TLibCommon.TComDataCU*& rpcTempCU, PartSize ePartSize  );
//#endif
func (this *TEncCu) xCheckRDCostIntra(rpcBestCU **TLibCommon.TComDataCU, rpcTempCU **TLibCommon.TComDataCU, eSize TLibCommon.PartSize) {
    //fmt.Printf("Enter xCheckRDCostIntra\n");
    
    uiDepth := uint(rpcTempCU.GetDepth1(0))

    rpcTempCU.SetSkipFlagSubParts(false, 0, uiDepth)

    rpcTempCU.SetPartSizeSubParts(eSize, 0, uiDepth)
    rpcTempCU.SetPredModeSubParts(TLibCommon.MODE_INTRA, 0, uiDepth)
    rpcTempCU.SetCUTransquantBypassSubParts(this.m_pcEncCfg.GetCUTransquantBypassFlagValue(), 0, uiDepth)

    bSeparateLumaChroma := true // choose estimation mode
    uiPreCalcDistC := uint(0)
    if !bSeparateLumaChroma {
        this.m_pcPredSearch.preestChromaPredMode(*rpcTempCU, this.m_ppcOrigYuv[uiDepth], this.m_ppcPredYuvTemp[uiDepth], 0)
    }
    this.m_pcPredSearch.estIntraPredQT(*rpcTempCU, this.m_ppcOrigYuv[uiDepth], this.m_ppcPredYuvTemp[uiDepth], this.m_ppcResiYuvTemp[uiDepth], this.m_ppcRecoYuvTemp[uiDepth], &uiPreCalcDistC, bSeparateLumaChroma)

    this.m_ppcRecoYuvTemp[uiDepth].CopyToPicLuma(rpcTempCU.GetPic().GetPicYuvRec(), rpcTempCU.GetAddr(), rpcTempCU.GetZorderIdxInCU(), 0, 0)
	//fmt.Printf("enter estIntraPredChromaQT\n")
    this.m_pcPredSearch.estIntraPredChromaQT(*rpcTempCU, this.m_ppcOrigYuv[uiDepth], this.m_ppcPredYuvTemp[uiDepth], this.m_ppcResiYuvTemp[uiDepth], this.m_ppcRecoYuvTemp[uiDepth], uiPreCalcDistC)
	//fmt.Printf("outof estIntraPredChromaQT\n")
    this.m_pcEntropyCoder.resetBits()
    if rpcTempCU.GetSlice().GetPPS().GetTransquantBypassEnableFlag() {
        this.m_pcEntropyCoder.encodeCUTransquantBypassFlag(*rpcTempCU, 0, true)
    }
    this.m_pcEntropyCoder.encodeSkipFlag(*rpcTempCU, 0, true)
    this.m_pcEntropyCoder.encodePredMode(*rpcTempCU, 0, true)
    this.m_pcEntropyCoder.encodePartSize(*rpcTempCU, 0, uiDepth, true)
    this.m_pcEntropyCoder.encodePredInfo(*rpcTempCU, 0, true)
    this.m_pcEntropyCoder.encodeIPCMInfo(*rpcTempCU, 0, true)
	//fmt.Printf("outof encodeIPCMInfo\n")
	
    // Encode Coefficients
    bCodeDQP := this.getdQPFlag()
    this.m_pcEntropyCoder.encodeCoeff(*rpcTempCU, 0, uiDepth, uint((*rpcTempCU).GetWidth1(0)), uint(rpcTempCU.GetHeight1(0)), &bCodeDQP)
    this.setdQPFlag(bCodeDQP)
	
    if this.m_bUseSBACRD {
        this.m_pcRDGoOnSbacCoder.store(this.m_pppcRDSbacCoder[uiDepth][TLibCommon.CI_TEMP_BEST])
    }
    rpcTempCU.SetTotalBits(this.m_pcEntropyCoder.getNumberOfWrittenBits())
    if this.m_pcEncCfg.GetUseSBACRD() {
        rpcTempCU.SetTotalBins(this.m_pcEntropyCoder.m_pcEntropyCoderIf.getEncBinIf().getTEncBinCABAC().getBinsCoded())
    }
    rpcTempCU.SetTotalCost(this.m_pcRdCost.calcRdCost(rpcTempCU.GetTotalBits(), rpcTempCU.GetTotalDistortion(), false, TLibCommon.DF_DEFAULT))

    this.xCheckDQP(*rpcTempCU)
    this.xCheckBestMode3(rpcBestCU, rpcTempCU, uiDepth)
    
    //fmt.Printf("Exit xCheckRDCostIntra\n");
}

func (this *TEncCu) xCheckDQP(pcCU *TLibCommon.TComDataCU) {
    uiDepth := uint(pcCU.GetDepth1(0))

    if pcCU.GetSlice().GetPPS().GetUseDQP() && (pcCU.GetSlice().GetSPS().GetMaxCUWidth()>>uiDepth) >= pcCU.GetSlice().GetPPS().GetMinCuDQPSize() {
        if pcCU.GetCbf3(0, TLibCommon.TEXT_LUMA, 0) != 0 || pcCU.GetCbf3(0, TLibCommon.TEXT_CHROMA_U, 0) != 0 || pcCU.GetCbf3(0, TLibCommon.TEXT_CHROMA_V, 0) != 0 {
            //#if !RDO_WITHOUT_DQP_BITS
            this.m_pcEntropyCoder.resetBits()
            this.m_pcEntropyCoder.encodeQP(pcCU, 0, false)
            pcCU.SetTotalBits(pcCU.GetTotalBits() + this.m_pcEntropyCoder.getNumberOfWrittenBits()) // dQP bits
            if this.m_pcEncCfg.GetUseSBACRD() {
                pcCU.SetTotalBins(pcCU.GetTotalBins() + this.m_pcEntropyCoder.m_pcEntropyCoderIf.getEncBinIf().getTEncBinCABAC().getBinsCoded())
            }
            pcCU.SetTotalCost(this.m_pcRdCost.calcRdCost(pcCU.GetTotalBits(), pcCU.GetTotalDistortion(), false, TLibCommon.DF_DEFAULT))
            //#endif
        } else {
            pcCU.SetQPSubParts(int(pcCU.GetRefQP(0)), 0, uiDepth) // set QP to default QP
        }
    }
}

func (this *TEncCu) xCheckIntraPCM(rpcBestCU **TLibCommon.TComDataCU, rpcTempCU **TLibCommon.TComDataCU) {
    uiDepth := uint(rpcTempCU.GetDepth1(0))

    rpcTempCU.SetSkipFlagSubParts(false, 0, uiDepth)
    rpcTempCU.SetIPCMFlag(0, true)
    rpcTempCU.SetIPCMFlagSubParts(true, 0, uint(rpcTempCU.GetDepth1(0)))
    rpcTempCU.SetPartSizeSubParts(TLibCommon.SIZE_2Nx2N, 0, uiDepth)
    rpcTempCU.SetPredModeSubParts(TLibCommon.MODE_INTRA, 0, uiDepth)
    rpcTempCU.SetTrIdxSubParts(0, 0, uiDepth)
    rpcTempCU.SetCUTransquantBypassSubParts(this.m_pcEncCfg.GetCUTransquantBypassFlagValue(), 0, uiDepth)

    this.m_pcPredSearch.IPCMSearch(*rpcTempCU, this.m_ppcOrigYuv[uiDepth], this.m_ppcPredYuvTemp[uiDepth], this.m_ppcResiYuvTemp[uiDepth], this.m_ppcRecoYuvTemp[uiDepth])

    if this.m_bUseSBACRD {
        this.m_pcRDGoOnSbacCoder.load(this.m_pppcRDSbacCoder[uiDepth][TLibCommon.CI_CURR_BEST])
    }

    this.m_pcEntropyCoder.resetBits()
    if rpcTempCU.GetSlice().GetPPS().GetTransquantBypassEnableFlag() {
        this.m_pcEntropyCoder.encodeCUTransquantBypassFlag(*rpcTempCU, 0, true)
    }
    this.m_pcEntropyCoder.encodeSkipFlag(*rpcTempCU, 0, true)
    this.m_pcEntropyCoder.encodePredMode(*rpcTempCU, 0, true)
    this.m_pcEntropyCoder.encodePartSize(*rpcTempCU, 0, uiDepth, true)
    this.m_pcEntropyCoder.encodeIPCMInfo(*rpcTempCU, 0, true)

    if this.m_bUseSBACRD {
        this.m_pcRDGoOnSbacCoder.store(this.m_pppcRDSbacCoder[uiDepth][TLibCommon.CI_TEMP_BEST])
    }

    rpcTempCU.SetTotalBits(this.m_pcEntropyCoder.getNumberOfWrittenBits())
    if this.m_pcEncCfg.GetUseSBACRD() {
        rpcTempCU.SetTotalBins(this.m_pcEntropyCoder.m_pcEntropyCoderIf.getEncBinIf().getTEncBinCABAC().getBinsCoded())
    }
    rpcTempCU.SetTotalCost(this.m_pcRdCost.calcRdCost(rpcTempCU.GetTotalBits(), rpcTempCU.GetTotalDistortion(), false, TLibCommon.DF_DEFAULT))

    this.xCheckDQP(*rpcTempCU)
    this.xCheckBestMode3(rpcBestCU, rpcTempCU, uiDepth)
}

func (this *TEncCu) xCopyAMVPInfo(pSrc, pDst *TLibCommon.AMVPInfo) {
    pDst.IN = pSrc.IN
    for i := int(0); i < pSrc.IN; i++ {
        pDst.MvCand[i] = pSrc.MvCand[i]
    }
}

func (this *TEncCu) xCopyYuv2Pic(rpcPic *TLibCommon.TComPic, uiCUAddr, uiAbsPartIdx, uiDepth, uiSrcDepth uint, pcCU *TLibCommon.TComDataCU, uiLPelX, uiTPelY uint) {
    uiRPelX := uiLPelX + (pcCU.GetSlice().GetSPS().GetMaxCUWidth() >> uiDepth) - 1
    uiBPelY := uiTPelY + (pcCU.GetSlice().GetSPS().GetMaxCUHeight() >> uiDepth) - 1
    pcSlice := pcCU.GetPic().GetSlice(pcCU.GetPic().GetCurrSliceIdx())
    bSliceStart := pcSlice.GetSliceSegmentCurStartCUAddr() > rpcPic.GetPicSym().GetInverseCUOrderMap(int(pcCU.GetAddr()))*pcCU.GetPic().GetNumPartInCU()+uiAbsPartIdx &&
        pcSlice.GetSliceSegmentCurStartCUAddr() < rpcPic.GetPicSym().GetInverseCUOrderMap(int(pcCU.GetAddr()))*pcCU.GetPic().GetNumPartInCU()+uiAbsPartIdx+(pcCU.GetPic().GetNumPartInCU()>>(uiDepth<<1))
    bSliceEnd := pcSlice.GetSliceSegmentCurEndCUAddr() > rpcPic.GetPicSym().GetInverseCUOrderMap(int(pcCU.GetAddr()))*pcCU.GetPic().GetNumPartInCU()+uiAbsPartIdx &&
        pcSlice.GetSliceSegmentCurEndCUAddr() < rpcPic.GetPicSym().GetInverseCUOrderMap(int(pcCU.GetAddr()))*pcCU.GetPic().GetNumPartInCU()+uiAbsPartIdx+(pcCU.GetPic().GetNumPartInCU()>>(uiDepth<<1))
    if !bSliceEnd && !bSliceStart && (uiRPelX < pcSlice.GetSPS().GetPicWidthInLumaSamples()) && (uiBPelY < pcSlice.GetSPS().GetPicHeightInLumaSamples()) {
        uiAbsPartIdxInRaster := TLibCommon.G_auiZscanToRaster[uiAbsPartIdx]
        uiSrcBlkWidth := rpcPic.GetNumPartInWidth() >> (uiSrcDepth)
        uiBlkWidth := rpcPic.GetNumPartInWidth() >> (uiDepth)
        uiPartIdxX := ((uiAbsPartIdxInRaster % rpcPic.GetNumPartInWidth()) % uiSrcBlkWidth) / uiBlkWidth
        uiPartIdxY := ((uiAbsPartIdxInRaster / rpcPic.GetNumPartInWidth()) % uiSrcBlkWidth) / uiBlkWidth
        uiPartIdx := uiPartIdxY*(uiSrcBlkWidth/uiBlkWidth) + uiPartIdxX
        this.m_ppcRecoYuvBest[uiSrcDepth].CopyToPicYuv(rpcPic.GetPicYuvRec(), uiCUAddr, uiAbsPartIdx, uiDepth-uiSrcDepth, uiPartIdx)
    } else {
        uiQNumParts := (pcCU.GetPic().GetNumPartInCU() >> (uiDepth << 1)) >> 2

        for uiPartUnitIdx := uint(0); uiPartUnitIdx < 4; uiPartUnitIdx++ {
            uiSubCULPelX := uiLPelX + (pcCU.GetSlice().GetSPS().GetMaxCUWidth()>>(uiDepth+1))*(uiPartUnitIdx&1)
            uiSubCUTPelY := uiTPelY + (pcCU.GetSlice().GetSPS().GetMaxCUHeight()>>(uiDepth+1))*(uiPartUnitIdx>>1)

            bInSlice := rpcPic.GetPicSym().GetInverseCUOrderMap(int(pcCU.GetAddr()))*pcCU.GetPic().GetNumPartInCU()+uiAbsPartIdx+uiQNumParts > pcSlice.GetSliceSegmentCurStartCUAddr() &&
                rpcPic.GetPicSym().GetInverseCUOrderMap(int(pcCU.GetAddr()))*pcCU.GetPic().GetNumPartInCU()+uiAbsPartIdx < pcSlice.GetSliceSegmentCurEndCUAddr()
            if bInSlice && (uiSubCULPelX < pcSlice.GetSPS().GetPicWidthInLumaSamples()) && (uiSubCUTPelY < pcSlice.GetSPS().GetPicHeightInLumaSamples()) {
                this.xCopyYuv2Pic(rpcPic, uiCUAddr, uiAbsPartIdx, uiDepth+1, uiSrcDepth, pcCU, uiSubCULPelX, uiSubCUTPelY) // Copy Yuv data to picture Yuv
            }
            uiAbsPartIdx += uiQNumParts
        }
    }
}

func (this *TEncCu) xCopyYuv2Tmp(uiPartUnitIdx, uiNextDepth uint) {
    uiCurrDepth := uiNextDepth - 1
    this.m_ppcRecoYuvBest[uiNextDepth].CopyToPartYuv(this.m_ppcRecoYuvTemp[uiCurrDepth], uiPartUnitIdx)
}

func (this *TEncCu) getdQPFlag() bool  { return this.m_bEncodeDQP }
func (this *TEncCu) setdQPFlag(b bool) { this.m_bEncodeDQP = b }

//#if ADAPTIVE_QP_SELECTION
// Adaptive reconstruction level (ARL) statistics collection functions
func (this *TEncCu) xLcuCollectARLStats(rpcCU *TLibCommon.TComDataCU) {
    var cSum [TLibCommon.LEVEL_RANGE + 1]float64    //: the sum of DCT coefficients corresponding to datatype and quantization output
    var numSamples [TLibCommon.LEVEL_RANGE + 1]uint //: the number of coefficients corresponding to datatype and quantization output

    pCoeffY := rpcCU.GetCoeffY()
    pArlCoeffY := rpcCU.GetArlCoeffY()

    uiMinCUWidth := rpcCU.GetSlice().GetSPS().GetMaxCUWidth() >> rpcCU.GetSlice().GetSPS().GetMaxCUDepth()
    uiMinNumCoeffInCU := 1 << uiMinCUWidth

    //memset( cSum, 0, sizeof( Double )*(LEVEL_RANGE+1) );
    //memset( numSamples, 0, sizeof( UInt )*(LEVEL_RANGE+1) );

    // Collect stats to cSum[][] and numSamples[][]
    for i := uint(0); i < rpcCU.GetTotalNumPart(); i++ {
        uiTrIdx := uint(rpcCU.GetTransformIdx1(i))

        if rpcCU.GetPredictionMode1(i) == TLibCommon.MODE_INTER {
            if rpcCU.GetCbf3(i, TLibCommon.TEXT_LUMA, uiTrIdx) != 0 {
                this.xTuCollectARLStats(pCoeffY, pArlCoeffY, uiMinNumCoeffInCU, cSum[:], numSamples[:])
            }   //Note that only InterY is processed. QP rounding is based on InterY data only.
        }
        pCoeffY = pCoeffY[uiMinNumCoeffInCU:]
        pArlCoeffY = pArlCoeffY[uiMinNumCoeffInCU:]
    }

    for u := 1; u < TLibCommon.LEVEL_RANGE; u++ {
        this.m_pcTrQuant.GetSliceSumC()[u] += cSum[u]
        this.m_pcTrQuant.GetSliceNSamples()[u] += int(numSamples[u])
    }
    this.m_pcTrQuant.GetSliceSumC()[TLibCommon.LEVEL_RANGE] += cSum[TLibCommon.LEVEL_RANGE]
    this.m_pcTrQuant.GetSliceNSamples()[TLibCommon.LEVEL_RANGE] += int(numSamples[TLibCommon.LEVEL_RANGE])
}

func (this *TEncCu) xTuCollectARLStats(rpcCoeff []TLibCommon.TCoeff, rpcArlCoeff []TLibCommon.TCoeff, NumCoeffInCU int, cSum []float64, numSamples []uint) int {
    for n := 0; n < NumCoeffInCU; n++ {
        u := int(TLibCommon.ABS(rpcCoeff[n]).(TLibCommon.TCoeff))
        absc := rpcArlCoeff[n]

        if u != 0 {
            if u < TLibCommon.LEVEL_RANGE {
                cSum[u] += float64(absc)
                numSamples[u]++
            } else {
                cSum[TLibCommon.LEVEL_RANGE] += float64(absc) - float64(u<<TLibCommon.ARL_C_PRECISION)
                numSamples[TLibCommon.LEVEL_RANGE]++
            }
        }
    }

    return 0
}

//#endif

//#if AMP_ENC_SPEEDUP
//#if AMP_MRG
func (this *TEncCu) deriveTestModeAMP(rpcBestCU *TLibCommon.TComDataCU, eParentPartSize TLibCommon.PartSize, bTestAMP_Hor, bTestAMP_Ver, bTestMergeAMP_Hor, bTestMergeAMP_Ver *bool) {
    if rpcBestCU.GetPartitionSize1(0) == TLibCommon.SIZE_2NxN {
        *bTestAMP_Hor = true
    } else if rpcBestCU.GetPartitionSize1(0) == TLibCommon.SIZE_Nx2N {
        *bTestAMP_Ver = true
    } else if rpcBestCU.GetPartitionSize1(0) == TLibCommon.SIZE_2Nx2N && rpcBestCU.GetMergeFlag1(0) == false && rpcBestCU.IsSkipped(0) == false {
        *bTestAMP_Hor = true
        *bTestAMP_Ver = true
    }

    //#if AMP_MRG
    //! Utilizing the partition size of parent PU
    if eParentPartSize >= TLibCommon.SIZE_2NxnU && eParentPartSize <= TLibCommon.SIZE_nRx2N {
        *bTestMergeAMP_Hor = true
        *bTestMergeAMP_Ver = true
    }

    if eParentPartSize == TLibCommon.SIZE_NONE { //! if parent is intra
        if rpcBestCU.GetPartitionSize1(0) == TLibCommon.SIZE_2NxN {
            *bTestMergeAMP_Hor = true
        } else if rpcBestCU.GetPartitionSize1(0) == TLibCommon.SIZE_Nx2N {
            *bTestMergeAMP_Ver = true
        }
    }

    if rpcBestCU.GetPartitionSize1(0) == TLibCommon.SIZE_2Nx2N && rpcBestCU.IsSkipped(0) == false {
        *bTestMergeAMP_Hor = true
        *bTestMergeAMP_Ver = true
    }

    if rpcBestCU.GetWidth1(0) == 64 {
        *bTestAMP_Hor = false
        *bTestAMP_Ver = false
    }
    /*#else
      //! Utilizing the partition size of parent PU
      if ( eParentPartSize >= SIZE_2NxnU && eParentPartSize <= SIZE_nRx2N )
      {
        bTestAMP_Hor = true;
        bTestAMP_Ver = true;
      }

      if ( eParentPartSize == SIZE_2Nx2N )
      {
        bTestAMP_Hor = false;
        bTestAMP_Ver = false;
      }
    #endif*/
}

//#else
//func (this *TEncCu) deriveTestModeAMP (TLibCommon.TComDataCU *&rpcBestCU, PartSize eParentPartSize, Bool &bTestAMP_Hor, Bool &bTestAMP_Ver);
//#endif
//#endif

func (this *TEncCu) xFillPCMBuffer(pCU *TLibCommon.TComDataCU, pOrgYuv *TLibCommon.TComYuv) {
    width := uint(pCU.GetWidth1(0))
    height := uint(pCU.GetHeight1(0))

    pSrcY := pOrgYuv.GetLumaAddr2(0, width)
    pDstY := pCU.GetPCMSampleY()
    srcStride := pOrgYuv.GetStride()

    for y := uint(0); y < height; y++ {
        for x := uint(0); x < width; x++ {
            pDstY[x] = pSrcY[x]
        }
        pDstY = pDstY[width:]
        pSrcY = pSrcY[srcStride:]
    }

    pSrcCb := pOrgYuv.GetCbAddr()
    pSrcCr := pOrgYuv.GetCrAddr()


    pDstCb := pCU.GetPCMSampleCb()
    pDstCr := pCU.GetPCMSampleCr()


    srcStrideC := pOrgYuv.GetCStride()
    heightC := height >> 1
    widthC := width >> 1

    for y := uint(0); y < heightC; y++ {
        for x := uint(0); x < widthC; x++ {
            pDstCb[x] = pSrcCb[x]
            pDstCr[x] = pSrcCr[x]
        }
        pDstCb = pDstCb[widthC:]
        pDstCr = pDstCr[widthC:]
        pSrcCb = pSrcCb[srcStrideC:]
        pSrcCr = pSrcCr[srcStrideC:]
    }
}
