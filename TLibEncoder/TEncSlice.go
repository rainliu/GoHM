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
)

// ====================================================================================================================
// Class definition
// ====================================================================================================================

type WeightPredAnalysis struct {
    m_weighted_pred_flag   bool
    m_weighted_bipred_flag bool
    m_wp                   [2][TLibCommon.MAX_NUM_REF][3]TLibCommon.WpScalingParam
}

func NewWeightPredAnalysis() *WeightPredAnalysis {
    return &WeightPredAnalysis{}
}

// WP analysis :
func (this *WeightPredAnalysis) xCalcACDCParamSlice(slice *TLibCommon.TComSlice) bool {
    return false
}
func (this *WeightPredAnalysis) xEstimateWPParamSlice(slice *TLibCommon.TComSlice) bool {
    return false
}
func (this *WeightPredAnalysis) xStoreWPparam(weighted_pred_flag, weighted_bipred_flag bool) {
}
func (this *WeightPredAnalysis) xRestoreWPparam(slice *TLibCommon.TComSlice) {
}
func (this *WeightPredAnalysis) xCheckWPEnable(slice *TLibCommon.TComSlice) {
}

/*
  Int64   xCalcDCValueSlice(TComSlice *slice, Pel *pPel,Int *iSample);
  Int64   xCalcACValueSlice(TComSlice *slice, Pel *pPel, Int64 iDC);
  Int64   xCalcDCValueUVSlice(TComSlice *slice, Pel *pPel, Int *iSample);
  Int64   xCalcACValueUVSlice(TComSlice *slice, Pel *pPel, Int64 iDC);
  Int64   xCalcSADvalueWPSlice(TComSlice *slice, Pel *pOrgPel, Pel *pRefPel, Int iDenom, Int iWeight, Int iOffset);

  Int64   xCalcDCValue(Pel *pPel, Int iWidth, Int iHeight, Int iStride);
  Int64   xCalcACValue(Pel *pPel, Int iWidth, Int iHeight, Int iStride, Int64 iDC);
  Int64   xCalcSADvalueWP(Int bitDepth, Pel *pOrgPel, Pel *pRefPel, Int iWidth, Int iHeight, Int iOrgStride, Int iRefStride, Int iDenom, Int iWeight, Int iOffset);
  Bool    xSelectWP(TComSlice *slice, wpScalingParam weightPredTable[2][MAX_NUM_REF][3], Int iDenom);
  Bool    xUpdatingWPParameters(TComSlice *slice, wpScalingParam weightPredTable[2][MAX_NUM_REF][3], Int log2Denom);
*/

/// slice encoder class
type TEncSlice struct {
    WeightPredAnalysis

    // encoder configuration
    m_pcEncTop *TEncTop ///< encoder configuration class
    m_pcCfg    *TEncCfg ///< encoder configuration class

    // pictures
    m_pcListPic     *list.List             ///< list of pictures
    m_apcPicYuvPred *TLibCommon.TComPicYuv ///< prediction picture buffer
    m_apcPicYuvResi *TLibCommon.TComPicYuv ///< residual picture buffer

    // processing units
    m_pcGOPEncoder *TEncGOP ///< GOP encoder
    m_pcCuEncoder  *TEncCu  ///< CU encoder

    // encoder search
    m_pcPredSearch *TEncSearch ///< encoder search class

    // coding tools
    m_pTraceFile	  io.Writer
    m_pcEntropyCoder *TEncEntropy            ///< entropy encoder
    m_pcCavlcCoder   *TEncCavlc              ///< CAVLC encoder
    m_pcSbacCoder    *TEncSbac               ///< SBAC encoder
    m_pcBinCABAC     *TEncBinCABAC           ///< Bin encoder CABAC
    m_pcTrQuant      *TLibCommon.TComTrQuant ///< transform & quantization

    // RD optimization
    m_pcBitCounter                 *TLibCommon.TComBitCounter ///< bit counter
    m_pcRdCost                     *TEncRdCost                ///< RD cost computation
    m_pppcRDSbacCoder              [][]*TEncSbac              ///< storage for SBAC-based RD optimization
    m_pcRDGoOnSbacCoder            *TEncSbac                  ///< go-on SBAC encoder
    m_uiPicTotalBits               uint64                     ///< total bits for the picture
    m_uiPicDist                    uint64                     ///< total distortion for the picture
    m_dPicRdCost                   float64                    ///< picture-level RD cost
    m_pdRdPicLambda                []float64                  ///< array of lambda candidates
    m_pdRdPicQp                    []float64                  ///< array of picture QP candidates (double-type for lambda)
    m_piRdPicQp                    []int                      ///< array of picture QP candidates (Int-type)
    m_pcBufferBinCoderCABACs       []*TEncBinCABAC            ///< line of bin coder CABAC
    m_pcBufferSbacCoders           []*TEncSbac                ///< line to store temporary contexts
    m_pcBufferLowLatBinCoderCABACs []*TEncBinCABAC            ///< dependent tiles: line of bin coder CABAC
    m_pcBufferLowLatSbacCoders     []*TEncSbac                ///< dependent tiles: line to store temporary contexts
    m_pcRateCtrl                   *TEncRateCtrl              ///< Rate control manager
    m_uiSliceIdx                   uint

    CTXMem map[int]*TEncSbac
}

func NewTEncSlice() *TEncSlice {
    return &TEncSlice{}
}

func (this *TEncSlice) create(iWidth, iHeight int, iMaxCUWidth, iMaxCUHeight uint, uhTotalDepth byte) {
    // create prediction picture
    if this.m_apcPicYuvPred == nil {
        this.m_apcPicYuvPred = TLibCommon.NewTComPicYuv()
        this.m_apcPicYuvPred.Create(iWidth, iHeight, iMaxCUWidth, iMaxCUHeight, uint(uhTotalDepth))
    }

    // create residual picture
    if this.m_apcPicYuvResi == nil {
        this.m_apcPicYuvResi = TLibCommon.NewTComPicYuv()
        this.m_apcPicYuvResi.Create(iWidth, iHeight, iMaxCUWidth, iMaxCUHeight, uint(uhTotalDepth))
    }
}

func (this *TEncSlice) destroy() {
    // destroy prediction picture
    if this.m_apcPicYuvPred != nil {
        this.m_apcPicYuvPred.Destroy()
        //delete this.m_apcPicYuvPred;
        this.m_apcPicYuvPred = nil
    }

    // destroy residual picture
    if this.m_apcPicYuvResi != nil {
        this.m_apcPicYuvResi.Destroy()
        //delete this.m_apcPicYuvResi;
        this.m_apcPicYuvResi = nil
    }

    // free lambda and QP arrays
    if this.m_pdRdPicLambda != nil {
        //xFree( this.m_pdRdPicLambda );
        this.m_pdRdPicLambda = nil
    }
    if this.m_pdRdPicQp != nil {
        //xFree( this.m_pdRdPicQp     );
        this.m_pdRdPicQp = nil
    }
    if this.m_piRdPicQp != nil {
        //xFree( this.m_piRdPicQp     );
        this.m_piRdPicQp = nil
    }

    if this.m_pcBufferSbacCoders != nil {
        //delete[] this.m_pcBufferSbacCoders;
    }
    if this.m_pcBufferBinCoderCABACs != nil {
        //delete[] this.m_pcBufferBinCoderCABACs;
    }
    if this.m_pcBufferLowLatSbacCoders != nil {
        //delete[] this.m_pcBufferLowLatSbacCoders;
    }
    if this.m_pcBufferLowLatBinCoderCABACs != nil {
        //delete[] this.m_pcBufferLowLatBinCoderCABACs;
    }
}

func (this *TEncSlice) init(pcEncTop *TEncTop) {
    this.m_pcEncTop = pcEncTop
    this.m_pcCfg = pcEncTop.GetEncCfg()
    this.m_pcListPic = pcEncTop.getListPic()

    this.m_pcGOPEncoder = pcEncTop.getGOPEncoder()
    this.m_pcCuEncoder = pcEncTop.getCuEncoder()
    this.m_pcPredSearch = pcEncTop.getPredSearch()

	this.m_pTraceFile = pcEncTop.getTraceFile()
    this.m_pcEntropyCoder = pcEncTop.getEntropyCoder()
    this.m_pcCavlcCoder = pcEncTop.getCavlcCoder()
    this.m_pcSbacCoder = pcEncTop.getSbacCoder()
    this.m_pcBinCABAC = pcEncTop.getBinCABAC()
    this.m_pcTrQuant = pcEncTop.getTrQuant()

    this.m_pcBitCounter = pcEncTop.getBitCounter()
    this.m_pcRdCost = pcEncTop.getRdCost()
    this.m_pppcRDSbacCoder = pcEncTop.getRDSbacCoder()
    this.m_pcRDGoOnSbacCoder = pcEncTop.getRDGoOnSbacCoder()

    // create lambda and QP arrays
    this.m_pdRdPicLambda = make([]float64, this.m_pcCfg.GetDeltaQpRD()*2+1)
    this.m_pdRdPicQp = make([]float64, this.m_pcCfg.GetDeltaQpRD()*2+1)
    this.m_piRdPicQp = make([]int, this.m_pcCfg.GetDeltaQpRD()*2+1)
    this.m_pcRateCtrl = pcEncTop.getRateCtrl()
}

/// preparation of slice encoding (reference marking, QP and lambda)
func (this *TEncSlice) initEncSlice(pcPic *TLibCommon.TComPic, pocLast, pocCurr, iNumPicRcvd, iGOPid int, pSPS *TLibCommon.TComSPS, pPPS *TLibCommon.TComPPS) (rpcSlice *TLibCommon.TComSlice) {
    var dQP, dLambda float64

    rpcSlice = pcPic.GetSlice(0)
    rpcSlice.SetSPS(pSPS)
    rpcSlice.SetPPS(pPPS)
    rpcSlice.SetSliceBits(0)
    rpcSlice.SetPic(pcPic)
    rpcSlice.InitSlice()
    rpcSlice.SetPicOutputFlag(true)
    rpcSlice.SetPOC(pocCurr)

    // depth computation based on GOP size
    var depth int
    {
        poc := rpcSlice.GetPOC() % this.m_pcCfg.GetGOPSize()
        if poc == 0 {
            depth = 0
        } else {
            step := this.m_pcCfg.GetGOPSize()
            depth = 0
            for i := step >> 1; i >= 1; i >>= 1 {
                for j := i; j < this.m_pcCfg.GetGOPSize(); j += step {
                    if j == poc {
                        i = 0
                        break
                    }
                }
                step >>= 1
                depth++
            }
        }
    }

    // slice type
    var eSliceType TLibCommon.SliceType

    eSliceType = TLibCommon.B_SLICE
    if pocLast == 0 || int64(pocCurr)%int64(this.m_pcCfg.GetIntraPeriod()) == 0 || this.m_pcGOPEncoder.getGOPSize() == 0 {
        eSliceType = TLibCommon.I_SLICE
    }

    rpcSlice.SetSliceType(eSliceType)
    //fmt.Printf("getSliceType=%d, %d, %d mod %d=%d, %d\n", rpcSlice.GetSliceType(), pocLast, int64(pocCurr), int64(this.m_pcCfg.GetIntraPeriod()), int64(pocCurr)%int64(this.m_pcCfg.GetIntraPeriod()), this.m_pcGOPEncoder.getGOPSize());
	
    // ------------------------------------------------------------------------------------------------------------------
    // Non-referenced frame marking
    // ------------------------------------------------------------------------------------------------------------------

    if pocLast == 0 {
        rpcSlice.SetTemporalLayerNonReferenceFlag(false)
    } else {
        rpcSlice.SetTemporalLayerNonReferenceFlag(!this.m_pcCfg.GetGOPEntry(iGOPid).m_refPic)
    }
    rpcSlice.SetReferenced(true)

    // ------------------------------------------------------------------------------------------------------------------
    // QP setting
    // ------------------------------------------------------------------------------------------------------------------

    dQP = float64(this.m_pcCfg.GetQP())
    if eSliceType != TLibCommon.I_SLICE {
        if !((this.m_pcCfg.GetMaxDeltaQP() == 0) && (dQP == float64(-rpcSlice.GetSPS().GetQpBDOffsetY())) && (rpcSlice.GetSPS().GetUseLossless())) {
            dQP += float64(this.m_pcCfg.GetGOPEntry(iGOPid).m_QPOffset)
        }
    }

    // modify QP
    pdQPs := this.m_pcCfg.GetdQPs()
    if pdQPs != nil {
        dQP += float64(pdQPs[rpcSlice.GetPOC()])
    }
    //#if !RATE_CONTROL_LAMBDA_DOMAIN
    //  if ( this.m_pcCfg.GetUseRateCtrl())
    //  {
    //    dQP = this.m_pcRateCtrl.getFrameQP(rpcSlice.isReferenced(), rpcSlice.GetPOC());
    //  }
    //#endif
    // ------------------------------------------------------------------------------------------------------------------
    // Lambda computation
    // ------------------------------------------------------------------------------------------------------------------

    var iQP int
    dOrigQP := dQP

    // pre-compute lambda and QP values for all possible QP candidates
    for iDQpIdx := 0; iDQpIdx < 2*int(this.m_pcCfg.GetDeltaQpRD())+1; iDQpIdx++ {
        // compute QP value
        if iDQpIdx%2 != 0 {
            dQP = dOrigQP + float64((iDQpIdx+1)>>1)*(-1)
        } else {
            dQP = dOrigQP + float64((iDQpIdx+1)>>1)*(1)
        }

        // compute lambda value
        NumberBFrames := (this.m_pcCfg.GetGOPSize() - 1)
        SHIFT_QP := 12
        dLambda_scale := 1.0 - TLibCommon.CLIP3(0.0, 0.5, 0.05*float64(NumberBFrames)).(float64)
        //#if FULL_NBIT
        //        bitdepth_luma_qp_scale := 6 * (TLibCommon.G_bitDepth - 8);
        //#else
        bitdepth_luma_qp_scale := 0
        //#endif
        qp_temp := dQP + float64(bitdepth_luma_qp_scale-SHIFT_QP)

        //#if FULL_NBIT
        //    Double qp_temp_orig = (Double) dQP - SHIFT_QP;
        //#endif
        // Case #1: I or P-slices (key-frame)
        dQPFactor := this.m_pcCfg.GetGOPEntry(iGOPid).m_QPFactor

        if eSliceType == TLibCommon.I_SLICE {
            dQPFactor = 0.57 * dLambda_scale
        }
        dLambda = dQPFactor * math.Pow(2.0, qp_temp/3.0)

        if depth > 0 {
            //#if FULL_NBIT
            //        dLambda *= TLibCommon.CLIP3( 2.00, 4.00, (qp_temp_orig / 6.0) ).(float64); // (j == TLibCommon.B_SLICE && p_cur_frm.layer != 0 )
            //#else
            dLambda *= TLibCommon.CLIP3(2.00, 4.00, (qp_temp / 6.0)).(float64) // (j == TLibCommon.B_SLICE && p_cur_frm.layer != 0 )
            //#endif
        }

        // if hadamard is used in ME process
        if !this.m_pcCfg.GetUseHADME() && rpcSlice.GetSliceType() != TLibCommon.I_SLICE {
            dLambda *= 0.95
        }

        iQP = TLibCommon.MAX(int(-pSPS.GetQpBDOffsetY()), TLibCommon.MIN(int(TLibCommon.MAX_QP), int(math.Floor(dQP+0.5))).(int)).(int)
		
        this.m_pdRdPicLambda[iDQpIdx] = dLambda
        this.m_pdRdPicQp[iDQpIdx] = dQP
        this.m_piRdPicQp[iDQpIdx] = iQP
    }

    // obtain dQP = 0 case
    dLambda = this.m_pdRdPicLambda[0]
    dQP = this.m_pdRdPicQp[0]
    iQP = this.m_piRdPicQp[0]

    if rpcSlice.GetSliceType() != TLibCommon.I_SLICE {
        dLambda *= this.m_pcCfg.GetLambdaModifier(uint(this.m_pcCfg.GetGOPEntry(iGOPid).m_temporalId))
    }

    // store lambda
    this.m_pcRdCost.setLambda(dLambda)
    //#if WEIGHTED_CHROMA_DISTORTION
    // for RDO
    // in RdCost there is only one lambda because the luma and chroma bits are not separated, instead we weight the distortion of chroma.
    weight := float64(1.0)
    var qpc, chromaQPOffset int

    chromaQPOffset = rpcSlice.GetPPS().GetChromaCbQpOffset() + rpcSlice.GetSliceQpDeltaCb()
    qpc = TLibCommon.CLIP3(0, 57, iQP+chromaQPOffset).(int)
    weight = math.Pow(2.0, float64(iQP-int(TLibCommon.G_aucChromaScale[qpc]))/3.0) // takes into account of the chroma qp mapping and chroma qp Offset
    this.m_pcRdCost.setCbDistortionWeight(weight)

    chromaQPOffset = rpcSlice.GetPPS().GetChromaCrQpOffset() + rpcSlice.GetSliceQpDeltaCr()
    qpc = TLibCommon.CLIP3(0, 57, iQP+chromaQPOffset).(int)
    weight = math.Pow(2.0, float64(iQP-int(TLibCommon.G_aucChromaScale[qpc]))/3.0) // takes into account of the chroma qp mapping and chroma qp Offset
    this.m_pcRdCost.setCrDistortionWeight(weight)
    //#endif

    //#if RDOQ_CHROMA_LAMBDA
    // for RDOQ
    this.m_pcTrQuant.SetLambda(dLambda, dLambda/weight)
    //#else
    //  this.m_pcTrQuant.setLambda( dLambda );
    //#endif

    //#if SAO_CHROMA_LAMBDA
    // For SAO
    rpcSlice.SetLambda(dLambda, dLambda/weight)
    //#else
    //  rpcSlice   .setLambda( dLambda );
    //#endif

    //#if HB_LAMBDA_FOR_LDC
    // restore original slice type
    if pocLast == 0 || int64(pocCurr)%int64(this.m_pcCfg.GetIntraPeriod()) == 0 || this.m_pcGOPEncoder.getGOPSize() == 0 {
        eSliceType = TLibCommon.I_SLICE
    } else {
        eSliceType = eSliceType
    }

    rpcSlice.SetSliceType(eSliceType)
    //#endif

    if this.m_pcCfg.GetUseRecalculateQPAccordingToLambda() {
        dQP = this.xGetQPValueAccordingToLambda(dLambda)
        iQP = TLibCommon.MAX(int(-pSPS.GetQpBDOffsetY()), TLibCommon.MIN(int(TLibCommon.MAX_QP), int(math.Floor(dQP+0.5))).(int)).(int)
    }

    rpcSlice.SetSliceQp(iQP)
    //#if ADAPTIVE_QP_SELECTION
    rpcSlice.SetSliceQpBase(iQP)
    //#endif
    rpcSlice.SetSliceQpDelta(0)
    rpcSlice.SetSliceQpDeltaCb(0)
    rpcSlice.SetSliceQpDeltaCr(0)
    rpcSlice.SetNumRefIdx(TLibCommon.REF_PIC_LIST_0, this.m_pcCfg.GetGOPEntry(iGOPid).m_numRefPicsActive)
    rpcSlice.SetNumRefIdx(TLibCommon.REF_PIC_LIST_1, this.m_pcCfg.GetGOPEntry(iGOPid).m_numRefPicsActive)

    if rpcSlice.GetPPS().GetDeblockingFilterControlPresentFlag() {
        rpcSlice.GetPPS().SetDeblockingFilterOverrideEnabledFlag(!this.m_pcCfg.GetLoopFilterOffsetInPPS())
        rpcSlice.SetDeblockingFilterOverrideFlag(!this.m_pcCfg.GetLoopFilterOffsetInPPS())
        rpcSlice.GetPPS().SetPicDisableDeblockingFilterFlag(this.m_pcCfg.GetLoopFilterDisable())
        rpcSlice.SetDeblockingFilterDisable(this.m_pcCfg.GetLoopFilterDisable())
        if !rpcSlice.GetDeblockingFilterDisable() {
            if !this.m_pcCfg.GetLoopFilterOffsetInPPS() && eSliceType != TLibCommon.I_SLICE {
                rpcSlice.GetPPS().SetDeblockingFilterBetaOffsetDiv2(this.m_pcCfg.GetGOPEntry(iGOPid).m_betaOffsetDiv2 + this.m_pcCfg.GetLoopFilterBetaOffset())
                rpcSlice.GetPPS().SetDeblockingFilterTcOffsetDiv2(this.m_pcCfg.GetGOPEntry(iGOPid).m_tcOffsetDiv2 + this.m_pcCfg.GetLoopFilterTcOffset())
                rpcSlice.SetDeblockingFilterBetaOffsetDiv2(this.m_pcCfg.GetGOPEntry(iGOPid).m_betaOffsetDiv2 + this.m_pcCfg.GetLoopFilterBetaOffset())
                rpcSlice.SetDeblockingFilterTcOffsetDiv2(this.m_pcCfg.GetGOPEntry(iGOPid).m_tcOffsetDiv2 + this.m_pcCfg.GetLoopFilterTcOffset())
            } else {
                rpcSlice.GetPPS().SetDeblockingFilterBetaOffsetDiv2(this.m_pcCfg.GetLoopFilterBetaOffset())
                rpcSlice.GetPPS().SetDeblockingFilterTcOffsetDiv2(this.m_pcCfg.GetLoopFilterTcOffset())
                rpcSlice.SetDeblockingFilterBetaOffsetDiv2(this.m_pcCfg.GetLoopFilterBetaOffset())
                rpcSlice.SetDeblockingFilterTcOffsetDiv2(this.m_pcCfg.GetLoopFilterTcOffset())
            }
        }
    } else {
        rpcSlice.SetDeblockingFilterOverrideFlag(false)
        rpcSlice.SetDeblockingFilterDisable(false)
        rpcSlice.SetDeblockingFilterBetaOffsetDiv2(0)
        rpcSlice.SetDeblockingFilterTcOffsetDiv2(0)
    }

    rpcSlice.SetDepth(depth)

    pcPic.SetTLayer(uint(this.m_pcCfg.GetGOPEntry(iGOPid).m_temporalId))
    if eSliceType == TLibCommon.I_SLICE {
        pcPic.SetTLayer(0)
    }
    rpcSlice.SetTLayer(pcPic.GetTLayer())

    //assert( this.m_apcPicYuvPred );
    //assert( this.m_apcPicYuvResi );

    pcPic.SetPicYuvPred(this.m_apcPicYuvPred)
    pcPic.SetPicYuvResi(this.m_apcPicYuvResi)
    rpcSlice.SetSliceMode(uint(this.m_pcCfg.GetSliceMode()))
    rpcSlice.SetSliceArgument(uint(this.m_pcCfg.GetSliceArgument()))
    rpcSlice.SetSliceSegmentMode(uint(this.m_pcCfg.GetSliceSegmentMode()))
    rpcSlice.SetSliceSegmentArgument(uint(this.m_pcCfg.GetSliceSegmentArgument()))
    rpcSlice.SetMaxNumMergeCand(this.m_pcCfg.GetMaxNumMergeCand())
    this.xStoreWPparam(pPPS.GetUseWP(), pPPS.GetWPBiPred())
    
    return
}

//#if RATE_CONTROL_LAMBDA_DOMAIN
func (this *TEncSlice) resetQP(pic *TLibCommon.TComPic, sliceQP int, lambda float64) {
    slice := pic.GetSlice(0)

    // store lambda
    slice.SetSliceQp(sliceQP)
    //#if L0033_RC_BUGFIX
    slice.SetSliceQpBase(sliceQP)
    //#endif
    this.m_pcRdCost.setLambda(lambda)
    //#if WEIGHTED_CHROMA_DISTORTION
    // for RDO
    // in RdCost there is only one lambda because the luma and chroma bits are not separated, instead we weight the distortion of chroma.
    var weight float64
    var qpc, chromaQPOffset int

    chromaQPOffset = slice.GetPPS().GetChromaCbQpOffset() + slice.GetSliceQpDeltaCb()
    qpc = TLibCommon.CLIP3(0, 57, sliceQP+chromaQPOffset).(int)
    weight = math.Pow(2.0, float64(sliceQP-int(TLibCommon.G_aucChromaScale[qpc]))/3.0) // takes into account of the chroma qp mapping and chroma qp Offset
    this.m_pcRdCost.setCbDistortionWeight(weight)

    chromaQPOffset = slice.GetPPS().GetChromaCrQpOffset() + slice.GetSliceQpDeltaCr()
    qpc = TLibCommon.CLIP3(0, 57, sliceQP+chromaQPOffset).(int)
    weight = math.Pow(2.0, float64(sliceQP-int(TLibCommon.G_aucChromaScale[qpc]))/3.0) // takes into account of the chroma qp mapping and chroma qp Offset
    this.m_pcRdCost.setCrDistortionWeight(weight)
    //#endif

    //#if RDOQ_CHROMA_LAMBDA
    // for RDOQ
    this.m_pcTrQuant.SetLambda(lambda, lambda/weight)
    //#else
    //  this.m_pcTrQuant.setLambda( lambda );
    //#endif

    //#if SAO_CHROMA_LAMBDA
    // For SAO
    slice.SetLambda(lambda, lambda/weight)
    //#else
    //  slice   .setLambda( lambda );
    //#endif
}

//#else
//func (this *TEncSlice)    xLamdaRecalculation ( Int changeQP, Int idGOP, Int depth, SliceType eSliceType, TComSPS* pcSPS, TComSlice* pcSlice);
//#endif
// compress and encode slice
func (this *TEncSlice) precompressSlice(rpcPic *TLibCommon.TComPic) { ///< precompress slice for multi-loop opt.
    // if deltaQP RD is not used, simply return
    if this.m_pcCfg.GetDeltaQpRD() == 0 {
        return
    }

    //#if RATE_CONTROL_LAMBDA_DOMAIN
    if this.m_pcCfg.GetUseRateCtrl() {
        fmt.Printf("\nMultiple QP optimization is not allowed when rate control is enabled.")
        return //assert(0);
    }
    //#endif

    pcSlice := rpcPic.GetSlice(this.getSliceIdx())
    dPicRdCostBest := TLibCommon.MAX_DOUBLE
    uiQpIdxBest := uint(0)

    var dFrameLambda float64
    //#if FULL_NBIT
    //  SHIFT_QP := 12 + 6 * (TLibCommon.G_bitDepth - 8);
    //#else
    SHIFT_QP := 12
    //#endif

    // set frame lambda
    if this.m_pcCfg.GetGOPSize() > 1 {
        dFrameLambda = 0.68 * math.Pow(2, float64(this.m_piRdPicQp[0]-SHIFT_QP)/3.0) * float64(1+TLibCommon.B2U(pcSlice.IsInterB()))
    } else {
        dFrameLambda = 0.68 * math.Pow(2, float64(this.m_piRdPicQp[0]-SHIFT_QP)/3.0)
    }
    this.m_pcRdCost.setFrameLambda(dFrameLambda)

    // for each QP candidate
    for uiQpIdx := uint(0); uiQpIdx < 2*this.m_pcCfg.GetDeltaQpRD()+1; uiQpIdx++ {
        pcSlice.SetSliceQp(this.m_piRdPicQp[uiQpIdx])
        //#if ADAPTIVE_QP_SELECTION
        pcSlice.SetSliceQpBase(this.m_piRdPicQp[uiQpIdx])
        //#endif
        this.m_pcRdCost.setLambda(this.m_pdRdPicLambda[uiQpIdx])
        //#if WEIGHTED_CHROMA_DISTORTION
        // for RDO
        // in RdCost there is only one lambda because the luma and chroma bits are not separated, instead we weight the distortion of chroma.
        iQP := this.m_piRdPicQp[uiQpIdx]
        weight := float64(1.0)

        var qpc, chromaQPOffset int

        chromaQPOffset = pcSlice.GetPPS().GetChromaCbQpOffset() + pcSlice.GetSliceQpDeltaCb()
        qpc = TLibCommon.CLIP3(0, 57, iQP+chromaQPOffset).(int)
        weight = math.Pow(2.0, float64(iQP-int(TLibCommon.G_aucChromaScale[qpc]))/3.0) // takes into account of the chroma qp mapping and chroma qp Offset
        this.m_pcRdCost.setCbDistortionWeight(weight)

        chromaQPOffset = pcSlice.GetPPS().GetChromaCrQpOffset() + pcSlice.GetSliceQpDeltaCr()
        qpc = TLibCommon.CLIP3(0, 57, iQP+chromaQPOffset).(int)
        weight = math.Pow(2.0, float64(iQP-int(TLibCommon.G_aucChromaScale[qpc]))/3.0) // takes into account of the chroma qp mapping and chroma qp Offset
        this.m_pcRdCost.setCrDistortionWeight(weight)
        //#endif

        //#if RDOQ_CHROMA_LAMBDA
        // for RDOQ
        this.m_pcTrQuant.SetLambda(this.m_pdRdPicLambda[uiQpIdx], this.m_pdRdPicLambda[uiQpIdx]/weight)
        //#else
        //    this.m_pcTrQuant   .setLambda              ( this.m_pdRdPicLambda[uiQpIdx] );
        //#endif
        //#if SAO_CHROMA_LAMBDA
        // For SAO
        pcSlice.SetLambda(this.m_pdRdPicLambda[uiQpIdx], this.m_pdRdPicLambda[uiQpIdx]/weight)
        //#else
        //    pcSlice       .setLambda              ( this.m_pdRdPicLambda[uiQpIdx] );
        //#endif

        // try compress
        this.compressSlice(rpcPic)

        var dPicRdCost float64
        uiPicDist := uint64(this.m_uiPicDist)
        uiALFBits := uint64(0)

        this.m_pcGOPEncoder.preLoopFilterPicAll(rpcPic, &uiPicDist, &uiALFBits)

        // compute RD cost and choose the best
        dPicRdCost = this.m_pcRdCost.calcRdCost64(this.m_uiPicTotalBits+uiALFBits, uiPicDist, true, TLibCommon.DF_SSE_FRAME)

        if dPicRdCost < dPicRdCostBest {
            uiQpIdxBest = uiQpIdx
            dPicRdCostBest = dPicRdCost
        }
    }

    // set best values
    pcSlice.SetSliceQp(this.m_piRdPicQp[uiQpIdxBest])
    //#if ADAPTIVE_QP_SELECTION
    pcSlice.SetSliceQpBase(this.m_piRdPicQp[uiQpIdxBest])
    //#endif
    this.m_pcRdCost.setLambda(this.m_pdRdPicLambda[uiQpIdxBest])
    //#if WEIGHTED_CHROMA_DISTORTION
    // in RdCost there is only one lambda because the luma and chroma bits are not separated, instead we weight the distortion of chroma.
    iQP := this.m_piRdPicQp[uiQpIdxBest]
    weight := float64(1.0)

    var qpc, chromaQPOffset int

    chromaQPOffset = pcSlice.GetPPS().GetChromaCbQpOffset() + pcSlice.GetSliceQpDeltaCb()
    qpc = TLibCommon.CLIP3(0, 57, iQP+chromaQPOffset).(int)
    weight = math.Pow(2.0, float64(iQP-int(TLibCommon.G_aucChromaScale[qpc]))/3.0) // takes into account of the chroma qp mapping and chroma qp Offset
    this.m_pcRdCost.setCbDistortionWeight(weight)

    chromaQPOffset = pcSlice.GetPPS().GetChromaCrQpOffset() + pcSlice.GetSliceQpDeltaCr()
    qpc = TLibCommon.CLIP3(0, 57, iQP+chromaQPOffset).(int)
    weight = math.Pow(2.0, float64(iQP-int(TLibCommon.G_aucChromaScale[qpc]))/3.0) // takes into account of the chroma qp mapping and chroma qp Offset
    this.m_pcRdCost.setCrDistortionWeight(weight)
    //#endif

    //#if RDOQ_CHROMA_LAMBDA
    // for RDOQ
    this.m_pcTrQuant.SetLambda(this.m_pdRdPicLambda[uiQpIdxBest], this.m_pdRdPicLambda[uiQpIdxBest]/weight)
    //#else
    //  this.m_pcTrQuant   .setLambda              ( this.m_pdRdPicLambda[uiQpIdxBest] );
    //#endif
    //#if SAO_CHROMA_LAMBDA
    // For SAO
    pcSlice.SetLambda(this.m_pdRdPicLambda[uiQpIdxBest], this.m_pdRdPicLambda[uiQpIdxBest]/weight)
    //#else
    //  pcSlice       .setLambda              ( this.m_pdRdPicLambda[uiQpIdxBest] );
    //#endif
}

func (this *TEncSlice) compressSlice(rpcPic *TLibCommon.TComPic) { ///< analysis stage of slice
    var uiCUAddr, uiStartCUAddr, uiBoundingCUAddr uint
    rpcPic.GetSlice(this.getSliceIdx()).SetSliceSegmentBits(0)
    var pppcRDSbacCoder *TEncBinCABAC
    pcSlice := rpcPic.GetSlice(this.getSliceIdx())
    this.xDetermineStartAndBoundingCUAddr(&uiStartCUAddr, &uiBoundingCUAddr, rpcPic, false)

    // initialize cost values
    this.m_uiPicTotalBits = 0
    this.m_dPicRdCost = 0
    this.m_uiPicDist = 0

    // set entropy coder
    if this.m_pcCfg.GetUseSBACRD() {
        this.m_pcSbacCoder.init(this.m_pcBinCABAC)
        this.m_pcEntropyCoder.setEntropyCoder(this.m_pcSbacCoder, pcSlice, this.m_pTraceFile)
        this.m_pcEntropyCoder.resetEntropy()
        this.m_pppcRDSbacCoder[0][TLibCommon.CI_CURR_BEST].load(this.m_pcSbacCoder)
        pppcRDSbacCoder = this.m_pppcRDSbacCoder[0][TLibCommon.CI_CURR_BEST].getEncBinIf().getTEncBinCABAC()
        pppcRDSbacCoder.setBinCountingEnableFlag(false)
        pppcRDSbacCoder.setBinsCoded(0)
    } else {
        this.m_pcEntropyCoder.setEntropyCoder(this.m_pcCavlcCoder, pcSlice, this.m_pTraceFile)
        this.m_pcEntropyCoder.resetEntropy()
        this.m_pcEntropyCoder.setBitstream(this.m_pcBitCounter)
    }

    //------------------------------------------------------------------------------
    //  Weighted Prediction parameters estimation.
    //------------------------------------------------------------------------------
    // calculate AC/DC values for current picture
    if pcSlice.GetPPS().GetUseWP() || pcSlice.GetPPS().GetWPBiPred() {
        this.xCalcACDCParamSlice(pcSlice)
    }

    bWp_explicit := (pcSlice.GetSliceType() == TLibCommon.P_SLICE && pcSlice.GetPPS().GetUseWP()) || (pcSlice.GetSliceType() == TLibCommon.B_SLICE && pcSlice.GetPPS().GetWPBiPred())

    if bWp_explicit {
        //------------------------------------------------------------------------------
        //  Weighted Prediction implemented at Slice level. SliceMode=2 is not supported yet.
        //------------------------------------------------------------------------------
        if pcSlice.GetSliceMode() == 2 || pcSlice.GetSliceSegmentMode() == 2 {
            fmt.Printf("Weighted Prediction is not supported with slice mode determined by max number of bins.\n")
            return
        }

        this.xEstimateWPParamSlice(pcSlice)
        pcSlice.InitWpScaling()

        // check WP on/off
        this.xCheckWPEnable(pcSlice)
    }

    //#if ADAPTIVE_QP_SELECTION
    if this.m_pcCfg.GetUseAdaptQpSelect() {
        this.m_pcTrQuant.ClearSliceARLCnt()
        if pcSlice.GetSliceType() != TLibCommon.I_SLICE {
            qpBase := pcSlice.GetSliceQpBase()
            pcSlice.SetSliceQp(qpBase + this.m_pcTrQuant.GetQpDelta(qpBase))
        }
    }
    //#endif
    pcEncTop := this.m_pcEncTop
    ppppcRDSbacCoders := pcEncTop.getRDSbacCoders()
    pcBitCounters := pcEncTop.getBitCounters()
    iNumSubstreams := 1
    uiTilesAcross := uint(0)

    if this.m_pcCfg.GetUseSBACRD() {
        iNumSubstreams = pcSlice.GetPPS().GetNumSubstreams()
        uiTilesAcross = uint(rpcPic.GetPicSym().GetNumColumnsMinus1()) + 1
        //delete[] this.m_pcBufferSbacCoders;
        //delete[] this.m_pcBufferBinCoderCABACs;
        this.m_pcBufferSbacCoders = make([]*TEncSbac, uiTilesAcross)
        this.m_pcBufferBinCoderCABACs = make([]*TEncBinCABAC, uiTilesAcross)
        for ui := uint(0); ui < uiTilesAcross; ui++ {
            this.m_pcBufferSbacCoders[ui] = NewTEncSbac()
            this.m_pcBufferBinCoderCABACs[ui] = NewTEncBinCABAC()

            this.m_pcBufferSbacCoders[ui].init(this.m_pcBufferBinCoderCABACs[ui])
        }
        for ui := uint(0); ui < uiTilesAcross; ui++ {
            this.m_pcBufferSbacCoders[ui].load(this.m_pppcRDSbacCoder[0][TLibCommon.CI_CURR_BEST]) //init. state
        }

        for ui := int(0); ui < iNumSubstreams; ui++ { //init all sbac coders for RD optimization
            ppppcRDSbacCoders[ui][0][TLibCommon.CI_CURR_BEST].load(this.m_pppcRDSbacCoder[0][TLibCommon.CI_CURR_BEST])
        }
    }
    //if( this.m_pcCfg.GetUseSBACRD() )
    {
        //delete[] this.m_pcBufferLowLatSbacCoders;
        //delete[] this.m_pcBufferLowLatBinCoderCABACs;
        this.m_pcBufferLowLatSbacCoders = make([]*TEncSbac, uiTilesAcross)
        this.m_pcBufferLowLatBinCoderCABACs = make([]*TEncBinCABAC, uiTilesAcross)
        for ui := uint(0); ui < uiTilesAcross; ui++ {
            this.m_pcBufferLowLatSbacCoders[ui] = NewTEncSbac()
            this.m_pcBufferLowLatBinCoderCABACs[ui] = NewTEncBinCABAC()

            this.m_pcBufferLowLatSbacCoders[ui].init(this.m_pcBufferLowLatBinCoderCABACs[ui])
        }
        for ui := uint(0); ui < uiTilesAcross; ui++ {
            this.m_pcBufferLowLatSbacCoders[ui].load(this.m_pppcRDSbacCoder[0][TLibCommon.CI_CURR_BEST]) //init. state
        }
    }
    uiWidthInLCUs := rpcPic.GetPicSym().GetFrameWidthInCU()
    //UInt uiHeightInLCUs = rpcPic.getPicSym().getFrameHeightInCU();
    var uiCol, uiLin, uiSubStrm, uiTileCol, uiTileStartLCU, uiTileLCUX uint
    depSliceSegmentsEnabled := pcSlice.GetPPS().GetDependentSliceSegmentsEnabledFlag()
    uiCUAddr = rpcPic.GetPicSym().GetCUOrderMap(int(uiStartCUAddr / rpcPic.GetNumPartInCU()))
    uiTileStartLCU = rpcPic.GetPicSym().GetTComTile(rpcPic.GetPicSym().GetTileIdxMap(int(uiCUAddr))).GetFirstCUAddr()
    if depSliceSegmentsEnabled {
        if (pcSlice.GetSliceSegmentCurStartCUAddr() != pcSlice.GetSliceCurStartCUAddr()) && (uiCUAddr != uiTileStartLCU) {
            if this.m_pcCfg.GetWaveFrontsynchro() != 0 {
                uiTileCol = rpcPic.GetPicSym().GetTileIdxMap(int(uiCUAddr)) % uint(rpcPic.GetPicSym().GetNumColumnsMinus1()+1)
                this.m_pcBufferSbacCoders[uiTileCol].loadContexts(this.CTXMem[1])

                iNumSubstreamsPerTile := iNumSubstreams / rpcPic.GetPicSym().GetNumTiles()
                uiCUAddr = rpcPic.GetPicSym().GetCUOrderMap(int(uiStartCUAddr / rpcPic.GetNumPartInCU()))
                uiLin = uiCUAddr / uiWidthInLCUs
                uiSubStrm = rpcPic.GetPicSym().GetTileIdxMap(int(rpcPic.GetPicSym().GetCUOrderMap(int(uiCUAddr))))*uint(iNumSubstreamsPerTile) + uiLin%uint(iNumSubstreamsPerTile)
                if (uiCUAddr%uiWidthInLCUs + 1) >= uiWidthInLCUs {
                    uiTileLCUX = uiTileStartLCU % uiWidthInLCUs
                    uiCol = uiCUAddr % uiWidthInLCUs
                    if uiCol == uiTileStartLCU {
                        this.CTXMem[0].loadContexts(this.m_pcSbacCoder)
                    }
                }
            }
            ppppcRDSbacCoders[uiSubStrm][0][TLibCommon.CI_CURR_BEST].loadContexts(this.CTXMem[0])
        } else {
            if this.m_pcCfg.GetWaveFrontsynchro() != 0 {
                this.CTXMem[1].loadContexts(this.m_pcSbacCoder)
            }
            this.CTXMem[0].loadContexts(this.m_pcSbacCoder)
        }
    }

    // for every CU in slice
    var uiEncCUOrder uint

    for uiEncCUOrder = uiStartCUAddr / rpcPic.GetNumPartInCU(); uiEncCUOrder < (uiBoundingCUAddr+(rpcPic.GetNumPartInCU()-1))/rpcPic.GetNumPartInCU(); uiCUAddr = rpcPic.GetPicSym().GetCUOrderMap(int(uiEncCUOrder)) {
        uiEncCUOrder++
        
        // initialize CU encoder
        pcCU := rpcPic.GetCU(uiCUAddr)
        pcCU.InitCU(rpcPic, uiCUAddr)
        
        fmt.Printf("Compress uiCUAddr=%d\n", pcCU.GetAddr());

        /*#if !RATE_CONTROL_LAMBDA_DOMAIN
            if(this.m_pcCfg.GetUseRateCtrl())
            {
              if(this.m_pcRateCtrl.calculateUnitQP())
              {
                xLamdaRecalculation(this.m_pcRateCtrl.getUnitQP(), this.m_pcRateCtrl.getGOPId(), pcSlice.GetDepth(), pcSlice.GetSliceType(), pcSlice.GetSPS(), pcSlice );
              }
            }
        #endif*/
        // inherit from TR if necessary, select substream to use.
        if this.m_pcCfg.GetUseSBACRD() {
            uiTileCol = rpcPic.GetPicSym().GetTileIdxMap(int(uiCUAddr)) % uint(rpcPic.GetPicSym().GetNumColumnsMinus1()+1) // what column of tiles are we in?
            uiTileStartLCU = rpcPic.GetPicSym().GetTComTile(rpcPic.GetPicSym().GetTileIdxMap(int(uiCUAddr))).GetFirstCUAddr()
            uiTileLCUX = uiTileStartLCU % uiWidthInLCUs
            //UInt uiSliceStartLCU = pcSlice.GetSliceCurStartCUAddr();
            uiCol = uiCUAddr % uiWidthInLCUs
            uiLin = uiCUAddr / uiWidthInLCUs
            if pcSlice.GetPPS().GetNumSubstreams() > 1 {
                // independent tiles => substreams are "per tile".  iNumSubstreams has already been multiplied.
                iNumSubstreamsPerTile := iNumSubstreams / rpcPic.GetPicSym().GetNumTiles()
                uiSubStrm = rpcPic.GetPicSym().GetTileIdxMap(int(uiCUAddr))*uint(iNumSubstreamsPerTile) + uiLin%uint(iNumSubstreamsPerTile)
            } else {
                // dependent tiles => substreams are "per frame".
                uiSubStrm = uiLin % uint(iNumSubstreams)
            }

            if ((pcSlice.GetPPS().GetNumSubstreams() > 1) || depSliceSegmentsEnabled) && (uiCol == uiTileLCUX) && this.m_pcCfg.GetWaveFrontsynchro() != 0 {
                // We'll sync if the TR is available.
                pcCUUp := pcCU.GetCUAbove()
                uiWidthInCU := rpcPic.GetFrameWidthInCU()
                uiMaxParts := uint(1) << (pcSlice.GetSPS().GetMaxCUDepth() << 1)
                var pcCUTR *TLibCommon.TComDataCU
                if pcCUUp != nil && ((uiCUAddr%uiWidthInCU + 1) < uiWidthInCU) {
                    pcCUTR = rpcPic.GetCU(uiCUAddr - uiWidthInCU + 1)
                }
                if (pcCUTR == nil) || (pcCUTR.GetSlice() == nil) ||
                    (pcCUTR.GetSCUAddr()+uiMaxParts-1 < pcSlice.GetSliceCurStartCUAddr()) ||
                    (rpcPic.GetPicSym().GetTileIdxMap(int(pcCUTR.GetAddr())) != rpcPic.GetPicSym().GetTileIdxMap(int(uiCUAddr))) {
                    // TR not available.
                } else {
                    // TR is available, we use it.
                    ppppcRDSbacCoders[uiSubStrm][0][TLibCommon.CI_CURR_BEST].loadContexts(this.m_pcBufferSbacCoders[uiTileCol])
                }
            }
            this.m_pppcRDSbacCoder[0][TLibCommon.CI_CURR_BEST].load(ppppcRDSbacCoders[uiSubStrm][0][TLibCommon.CI_CURR_BEST]) //this load is used to simplify the code
        }

        // reset the entropy coder
        if uiCUAddr == rpcPic.GetPicSym().GetTComTile(rpcPic.GetPicSym().GetTileIdxMap(int(uiCUAddr))).GetFirstCUAddr() && // must be first CU of tile
            uiCUAddr != 0 && // cannot be first CU of picture
            uiCUAddr != rpcPic.GetPicSym().GetPicSCUAddr(rpcPic.GetSlice(rpcPic.GetCurrSliceIdx()).GetSliceSegmentCurStartCUAddr())/rpcPic.GetNumPartInCU() &&
            uiCUAddr != rpcPic.GetPicSym().GetPicSCUAddr(rpcPic.GetSlice(rpcPic.GetCurrSliceIdx()).GetSliceCurStartCUAddr())/rpcPic.GetNumPartInCU() { // cannot be first CU of slice
            sliceType := pcSlice.GetSliceType()
            if !pcSlice.IsIntra() && pcSlice.GetPPS().GetCabacInitPresentFlag() && pcSlice.GetPPS().GetEncCABACTableIdx() != TLibCommon.I_SLICE {
                sliceType = TLibCommon.SliceType(pcSlice.GetPPS().GetEncCABACTableIdx())
            }
            this.m_pcEntropyCoder.updateContextTables3(sliceType, pcSlice.GetSliceQp(), false)
            this.m_pcEntropyCoder.setEntropyCoder(this.m_pppcRDSbacCoder[0][TLibCommon.CI_CURR_BEST], pcSlice, this.m_pTraceFile)
            this.m_pcEntropyCoder.updateContextTables2(sliceType, pcSlice.GetSliceQp())
            this.m_pcEntropyCoder.setEntropyCoder(this.m_pcSbacCoder, pcSlice, this.m_pTraceFile)
        }
        // if RD based on SBAC is used
        if this.m_pcCfg.GetUseSBACRD() {
            // set go-on entropy coder
            this.m_pcEntropyCoder.setEntropyCoder(this.m_pcRDGoOnSbacCoder, pcSlice, this.m_pTraceFile)
            this.m_pcEntropyCoder.setBitstream(pcBitCounters[uiSubStrm])

            this.m_pcRDGoOnSbacCoder.getEncBinIf().getTEncBinCABAC().setBinCountingEnableFlag(true)

            //#if RATE_CONTROL_LAMBDA_DOMAIN
            oldLambda := this.m_pcRdCost.getLambda()
            if this.m_pcCfg.GetUseRateCtrl() {
                estQP := pcSlice.GetSliceQp()
                estLambda := float64(-1.0)
                bpp := float64(-1.0)

                if rpcPic.GetSlice(0).GetSliceType() == TLibCommon.I_SLICE || !this.m_pcCfg.GetLCULevelRC() {
                    estQP = pcSlice.GetSliceQp()
                } else {
                    bpp = this.m_pcRateCtrl.getRCPic().getLCUTargetBpp()
                    estLambda = this.m_pcRateCtrl.getRCPic().getLCUEstLambda(bpp)
                    estQP = this.m_pcRateCtrl.getRCPic().getLCUEstQP(estLambda, pcSlice.GetSliceQp())
                    estQP = TLibCommon.CLIP3(-pcSlice.GetSPS().GetQpBDOffsetY(), TLibCommon.MAX_QP, estQP).(int)

                    this.m_pcRdCost.setLambda(estLambda)
                }

                this.m_pcRateCtrl.setRCQP(estQP)
                //#if L0033_RC_BUGFIX
                pcCU.GetSlice().SetSliceQpBase(estQP)
                //#endif
            }
            //#endif

            // run CU encoder
            this.m_pcCuEncoder.compressCU(pcCU)

            //#if RATE_CONTROL_LAMBDA_DOMAIN
            if this.m_pcCfg.GetUseRateCtrl() {
                SAD := this.m_pcCuEncoder.getLCUPredictionSAD()
                height := TLibCommon.MIN(pcSlice.GetSPS().GetMaxCUHeight(), pcSlice.GetSPS().GetPicHeightInLumaSamples()-uiCUAddr/rpcPic.GetFrameWidthInCU()*pcSlice.GetSPS().GetMaxCUHeight()).(int)
                width := TLibCommon.MIN(pcSlice.GetSPS().GetMaxCUWidth(), pcSlice.GetSPS().GetPicWidthInLumaSamples()-uiCUAddr%rpcPic.GetFrameWidthInCU()*pcSlice.GetSPS().GetMaxCUWidth()).(int)
                MAD := float64(SAD) / float64(height*width)
                MAD = MAD * MAD
                this.m_pcRateCtrl.getRCPic().getLCU(int(uiCUAddr)).m_MAD = MAD

                actualQP := g_RCInvalidQPValue
                actualLambda := float64(this.m_pcRdCost.getLambda())
                actualBits := pcCU.GetTotalBits()
                numberOfEffectivePixels := 0
                for idx := uint(0); idx < rpcPic.GetNumPartInCU(); idx++ {
                    if pcCU.GetPredictionMode1(idx) != TLibCommon.MODE_NONE && (!pcCU.IsSkipped(idx)) {
                        numberOfEffectivePixels = numberOfEffectivePixels + 16
                        break
                    }
                }

                if numberOfEffectivePixels == 0 {
                    actualQP = g_RCInvalidQPValue
                } else {
                    actualQP = int(pcCU.GetQP1(0))
                }
                this.m_pcRdCost.setLambda(oldLambda)

                this.m_pcRateCtrl.getRCPic().updateAfterLCU(this.m_pcRateCtrl.getRCPic().getLCUCoded(), int(actualBits), actualQP, actualLambda, this.m_pcCfg.GetLCULevelRC())
            }
            //#endif

            // restore entropy coder to an initial stage
            this.m_pcEntropyCoder.setEntropyCoder(this.m_pppcRDSbacCoder[0][TLibCommon.CI_CURR_BEST], pcSlice, this.m_pTraceFile)
            this.m_pcEntropyCoder.setBitstream(pcBitCounters[uiSubStrm])
            this.m_pcCuEncoder.setBitCounter(pcBitCounters[uiSubStrm])
            this.m_pcBitCounter = pcBitCounters[uiSubStrm]
            pppcRDSbacCoder.setBinCountingEnableFlag(true)
            this.m_pcBitCounter.ResetBits()
            pppcRDSbacCoder.setBinsCoded(0)
            this.m_pcCuEncoder.encodeCU(pcCU)

            pppcRDSbacCoder.setBinCountingEnableFlag(false)
            if this.m_pcCfg.GetSliceMode() == TLibCommon.FIXED_NUMBER_OF_BYTES && (pcSlice.GetSliceBits() + uint(this.m_pcEntropyCoder.getNumberOfWrittenBits())) > uint(this.m_pcCfg.GetSliceArgument())<<3 {
                pcSlice.SetNextSlice(true)
                break
            }
            if this.m_pcCfg.GetSliceSegmentMode() == TLibCommon.FIXED_NUMBER_OF_BYTES && pcSlice.GetSliceSegmentBits()+this.m_pcEntropyCoder.getNumberOfWrittenBits() > uint(this.m_pcCfg.GetSliceSegmentArgument()<<3) && pcSlice.GetSliceCurEndCUAddr() != pcSlice.GetSliceSegmentCurEndCUAddr() {
                pcSlice.SetNextSliceSegment(true)
                break
            }
            if this.m_pcCfg.GetUseSBACRD() {
                ppppcRDSbacCoders[uiSubStrm][0][TLibCommon.CI_CURR_BEST].load(this.m_pppcRDSbacCoder[0][TLibCommon.CI_CURR_BEST])

                //Store probabilties of second LCU in line into buffer
                if (uiCol == uiTileLCUX+1) && (depSliceSegmentsEnabled || (pcSlice.GetPPS().GetNumSubstreams() > 1)) && this.m_pcCfg.GetWaveFrontsynchro() != 0 {
                    this.m_pcBufferSbacCoders[uiTileCol].loadContexts(ppppcRDSbacCoders[uiSubStrm][0][TLibCommon.CI_CURR_BEST])
                }
            }
        } else { // other case: encodeCU is not called
            this.m_pcCuEncoder.compressCU(pcCU)
            this.m_pcCuEncoder.encodeCU(pcCU)
            if this.m_pcCfg.GetSliceMode() == TLibCommon.FIXED_NUMBER_OF_BYTES && (pcSlice.GetSliceBits() + this.m_pcEntropyCoder.getNumberOfWrittenBits()) > uint(this.m_pcCfg.GetSliceArgument())<<3 {
                pcSlice.SetNextSlice(true)
                break
            }
            if this.m_pcCfg.GetSliceSegmentMode() == TLibCommon.FIXED_NUMBER_OF_BYTES && pcSlice.GetSliceSegmentBits()+this.m_pcEntropyCoder.getNumberOfWrittenBits() > uint(this.m_pcCfg.GetSliceSegmentArgument())<<3 && pcSlice.GetSliceCurEndCUAddr() != pcSlice.GetSliceSegmentCurEndCUAddr() {
                pcSlice.SetNextSliceSegment(true)
                break
            }
        }

        this.m_uiPicTotalBits += uint64(pcCU.GetTotalBits())
        this.m_dPicRdCost += pcCU.GetTotalCost()
        this.m_uiPicDist += uint64(pcCU.GetTotalDistortion())
        /*#if !RATE_CONTROL_LAMBDA_DOMAIN
            if(this.m_pcCfg.GetUseRateCtrl())
            {
              this.m_pcRateCtrl.updateLCUData(pcCU, pcCU.getTotalBits(), pcCU.getQP(0));
              this.m_pcRateCtrl.updataRCUnitStatus();
            }
        #endif*/
    }

    if (pcSlice.GetPPS().GetNumSubstreams() > 1) && !depSliceSegmentsEnabled {
        pcSlice.SetNextSlice(true)
    }

    if depSliceSegmentsEnabled {
        if this.m_pcCfg.GetWaveFrontsynchro() != 0 {
            this.CTXMem[1].loadContexts(this.m_pcBufferSbacCoders[uiTileCol]) //ctx 2.LCU
        }
        this.CTXMem[0].loadContexts(this.m_pppcRDSbacCoder[0][TLibCommon.CI_CURR_BEST]) //ctx end of dep.slice
    }

    this.xRestoreWPparam(pcSlice)
    /*#if !RATE_CONTROL_LAMBDA_DOMAIN
      if(this.m_pcCfg.GetUseRateCtrl())
      {
        this.m_pcRateCtrl.updateFrameData(this.m_uiPicTotalBits);
      }
    #endif*/
}

func (this *TEncSlice) encodeSlice(rpcPic *TLibCommon.TComPic, rpcBitstream *TLibCommon.TComOutputBitstream, pcSubstreams []*TLibCommon.TComOutputBitstream) {
    var uiCUAddr, uiStartCUAddr, uiBoundingCUAddr uint
    pcSlice := rpcPic.GetSlice(this.getSliceIdx())

    uiStartCUAddr = pcSlice.GetSliceSegmentCurStartCUAddr()
    uiBoundingCUAddr = pcSlice.GetSliceSegmentCurEndCUAddr()
    // choose entropy coder
    {
        this.m_pcSbacCoder.init(this.m_pcBinCABAC)
        this.m_pcEntropyCoder.setEntropyCoder(this.m_pcSbacCoder, pcSlice, this.m_pTraceFile)
    }

    this.m_pcCuEncoder.setBitCounter(nil)
    this.m_pcBitCounter = nil
    // Appropriate substream bitstream is switched later.
    // for every CU
    /*#if ENC_DEC_TRACE
      TLibCommon.G_bJustDoIt = TLibCommon.G_bEncDecTraceEnable;
    #endif
      DTRACE_CABAC_VL( TLibCommon.G_nSymbolCounter++ );
      DTRACE_CABAC_T( "\tPOC: " );
      DTRACE_CABAC_V( rpcPic.GetPOC() );
      DTRACE_CABAC_T( "\n" );
    #if ENC_DEC_TRACE
      TLibCommon.G_bJustDoIt = TLibCommon.G_bEncDecTraceDisable;
    #endif*/

    pcEncTop := this.m_pcEncTop
    pcSbacCoders := pcEncTop.getSbacCoders() //coder for each substream
    iNumSubstreams := pcSlice.GetPPS().GetNumSubstreams()
    uiBitsOriginallyInSubstreams := uint(0)
    {
        uiTilesAcross := rpcPic.GetPicSym().GetNumColumnsMinus1() + 1
        for ui := 0; ui < uiTilesAcross; ui++ {
            this.m_pcBufferSbacCoders[ui].load(this.m_pcSbacCoder) //init. state
        }

        for iSubstrmIdx := 0; iSubstrmIdx < iNumSubstreams; iSubstrmIdx++ {
            uiBitsOriginallyInSubstreams += pcSubstreams[iSubstrmIdx].GetNumberOfWrittenBits()
        }

        for ui := 0; ui < uiTilesAcross; ui++ {
            this.m_pcBufferLowLatSbacCoders[ui].load(this.m_pcSbacCoder) //init. state
        }
    }

    uiWidthInLCUs := rpcPic.GetPicSym().GetFrameWidthInCU()
    var uiCol, uiLin, uiSubStrm, uiTileCol, uiTileStartLCU, uiTileLCUX uint
    depSliceSegmentsEnabled := pcSlice.GetPPS().GetDependentSliceSegmentsEnabledFlag()
    uiCUAddr = rpcPic.GetPicSym().GetCUOrderMap(int(uiStartCUAddr / rpcPic.GetNumPartInCU())) /* for tiles, uiStartCUAddr is NOT the real raster scan address, it is actually
       an encoding order index, so we need to convert the index (uiStartCUAddr)
       into the real raster scan address (uiCUAddr) via the CUOrderMap */
    uiTileStartLCU = rpcPic.GetPicSym().GetTComTile(rpcPic.GetPicSym().GetTileIdxMap(int(uiCUAddr))).GetFirstCUAddr()
    if depSliceSegmentsEnabled {
        if pcSlice.IsNextSlice() ||
            uiCUAddr == rpcPic.GetPicSym().GetTComTile(rpcPic.GetPicSym().GetTileIdxMap(int(uiCUAddr))).GetFirstCUAddr() {
            if this.m_pcCfg.GetWaveFrontsynchro() != 0 {
                this.CTXMem[1].loadContexts(this.m_pcSbacCoder)
            }
            this.CTXMem[0].loadContexts(this.m_pcSbacCoder)
        } else {
            if this.m_pcCfg.GetWaveFrontsynchro() != 0 {
                uiTileCol = rpcPic.GetPicSym().GetTileIdxMap(int(uiCUAddr)) % uint(rpcPic.GetPicSym().GetNumColumnsMinus1()+1)
                this.m_pcBufferSbacCoders[uiTileCol].loadContexts(this.CTXMem[1])

                iNumSubstreamsPerTile := iNumSubstreams / rpcPic.GetPicSym().GetNumTiles()
                uiLin = uiCUAddr / uiWidthInLCUs
                uiSubStrm = rpcPic.GetPicSym().GetTileIdxMap(int(rpcPic.GetPicSym().GetCUOrderMap(int(uiCUAddr))))*uint(iNumSubstreamsPerTile) + uiLin%uint(iNumSubstreamsPerTile)
                if (uiCUAddr%uiWidthInLCUs + 1) >= uiWidthInLCUs {
                    uiCol = uiCUAddr % uiWidthInLCUs
                    uiTileLCUX = uiTileStartLCU % uiWidthInLCUs
                    if uiCol == uiTileLCUX {
                        this.CTXMem[0].loadContexts(this.m_pcSbacCoder)
                    }
                }
            }
            pcSbacCoders[uiSubStrm].loadContexts(this.CTXMem[0])
        }
    }

    var uiEncCUOrder uint
    for uiEncCUOrder = uiStartCUAddr / rpcPic.GetNumPartInCU(); 
    	uiEncCUOrder < (uiBoundingCUAddr+rpcPic.GetNumPartInCU()-1)/rpcPic.GetNumPartInCU(); 
    	uiCUAddr = rpcPic.GetPicSym().GetCUOrderMap(int(uiEncCUOrder)) {
        uiEncCUOrder++ 
        
        fmt.Printf("Encoding uiCUAddr=%d\n", uiCUAddr);
        
        
        if this.m_pcCfg.GetUseSBACRD() {
            uiTileCol = rpcPic.GetPicSym().GetTileIdxMap(int(uiCUAddr)) % uint(rpcPic.GetPicSym().GetNumColumnsMinus1()+1) // what column of tiles are we in?
            uiTileStartLCU = rpcPic.GetPicSym().GetTComTile(rpcPic.GetPicSym().GetTileIdxMap(int(uiCUAddr))).GetFirstCUAddr()
            uiTileLCUX = uiTileStartLCU % uiWidthInLCUs
            //UInt uiSliceStartLCU = pcSlice.GetSliceCurStartCUAddr();
            uiCol = uiCUAddr % uiWidthInLCUs
            uiLin = uiCUAddr / uiWidthInLCUs
            if pcSlice.GetPPS().GetNumSubstreams() > 1 {
                // independent tiles => substreams are "per tile".  iNumSubstreams has already been multiplied.
                iNumSubstreamsPerTile := iNumSubstreams / rpcPic.GetPicSym().GetNumTiles()
                uiSubStrm = rpcPic.GetPicSym().GetTileIdxMap(int(uiCUAddr))*uint(iNumSubstreamsPerTile) + uiLin%uint(iNumSubstreamsPerTile)
            } else {
                // dependent tiles => substreams are "per frame".
                uiSubStrm = uiLin % uint(iNumSubstreams)
            }

            this.m_pcEntropyCoder.setBitstream(pcSubstreams[uiSubStrm])
            // Synchronize cabac probabilities with upper-right LCU if it's available and we're at the start of a line.
            if ((pcSlice.GetPPS().GetNumSubstreams() > 1) || depSliceSegmentsEnabled) && (uiCol == uiTileLCUX) && this.m_pcCfg.GetWaveFrontsynchro() != 0 {
                // We'll sync if the TR is available.
                pcCUUp := rpcPic.GetCU(uiCUAddr).GetCUAbove()
                uiWidthInCU := rpcPic.GetFrameWidthInCU()
                uiMaxParts := uint(1) << (pcSlice.GetSPS().GetMaxCUDepth() << 1)
                var pcCUTR *TLibCommon.TComDataCU
                if pcCUUp != nil && ((uiCUAddr%uiWidthInCU + 1) < uiWidthInCU) {
                    pcCUTR = rpcPic.GetCU(uiCUAddr - uiWidthInCU + 1)
                }
                if true /*bEnforceSliceRestriction*/ &&
                    ((pcCUTR == nil) || (pcCUTR.GetSlice() == nil) ||
                        (pcCUTR.GetSCUAddr()+uiMaxParts-1 < pcSlice.GetSliceCurStartCUAddr()) ||
                        (rpcPic.GetPicSym().GetTileIdxMap(int(pcCUTR.GetAddr())) != rpcPic.GetPicSym().GetTileIdxMap(int(uiCUAddr)))) {
                    // TR not available.
                } else {
                    // TR is available, we use it.
                    pcSbacCoders[uiSubStrm].loadContexts(this.m_pcBufferSbacCoders[uiTileCol])
                }
            }
            this.m_pcSbacCoder.load(pcSbacCoders[uiSubStrm]) //this load is used to simplify the code (avoid to change all the call to this.m_pcSbacCoder)
        }
        // reset the entropy coder
        if uiCUAddr == rpcPic.GetPicSym().GetTComTile(uint(rpcPic.GetPicSym().GetTileIdxMap(int(uiCUAddr)))).GetFirstCUAddr() && // must be first CU of tile
            uiCUAddr != 0 && // cannot be first CU of picture
            uiCUAddr != rpcPic.GetPicSym().GetPicSCUAddr(rpcPic.GetSlice(rpcPic.GetCurrSliceIdx()).GetSliceSegmentCurStartCUAddr())/rpcPic.GetNumPartInCU() &&
            uiCUAddr != rpcPic.GetPicSym().GetPicSCUAddr(rpcPic.GetSlice(rpcPic.GetCurrSliceIdx()).GetSliceCurStartCUAddr())/rpcPic.GetNumPartInCU() { // cannot be first CU of slice
            {
                // We're crossing into another tile, tiles are independent.
                // When tiles are independent, we have "substreams per tile".  Each substream has already been terminated, and we no longer
                // have to perform it here.
                if pcSlice.GetPPS().GetNumSubstreams() > 1 {
                    // do nothing.
                } else {
                    sliceType := pcSlice.GetSliceType()
                    if !pcSlice.IsIntra() && pcSlice.GetPPS().GetCabacInitPresentFlag() && pcSlice.GetPPS().GetEncCABACTableIdx() != TLibCommon.I_SLICE {
                        sliceType = TLibCommon.SliceType(pcSlice.GetPPS().GetEncCABACTableIdx())
                    }
                    this.m_pcEntropyCoder.updateContextTables2(sliceType, pcSlice.GetSliceQp())
                    // Byte-alignment in slice_data() when new tile
                    pcSubstreams[uiSubStrm].WriteByteAlignment()
                }
            }
            {
                uiCounter := uint(0)
                rbsp := pcSubstreams[uiSubStrm].GetFIFO()
                var v0, v1 byte
                for it := rbsp.Front(); it != nil; it = it.Next() {
                    /* 1) find the next emulated 00 00 {00,01,02,03}
                     * 2a) if not found, write all remaining bytes out, stop.
                     * 2b) otherwise, write all non-emulated bytes out
                     * 3) insert emulation_prevention_three_byte
                     */
                    found := it
                    for {
                        /* NB, end()-1, prevents finding a trailing two byte sequence */
                        //found = search_n(found, rbsp.end()-1, 2, 0);
                        for found != rbsp.Back() {
                            v0 = found.Value.(byte)
                            if found.Next() != nil {
                                v1 = found.Next().Value.(byte)
                            } else {
                                v1 = 0xFF
                            }

                            if v0 == 0 && v1 == 0 {
                                break
                            }
                            
                            found = found.Next()
                        }

                        found = found.Next()

                        /* if not found, found == end, otherwise found = second zero byte */
                        if found == nil {
                            break
                        }

                        found = found.Next()

                        if found.Value.(byte) <= 3 {
                            break
                        }
                    }

                    it = found
                    if found != nil {
                        it = rbsp.InsertBefore(emulation_prevention_three_byte[0], found)
                    }else{
                    	break;
                    }
                }

                uiAccumulatedSubstreamLength := uint(0)
                for iSubstrmIdx := 0; iSubstrmIdx < iNumSubstreams; iSubstrmIdx++ {
                    uiAccumulatedSubstreamLength += pcSubstreams[iSubstrmIdx].GetNumberOfWrittenBits()
                }
                // add bits coded in previous dependent slices + bits coded so far
                // add number of emulation prevention byte count in the tile
                pcSlice.AddTileLocation(((pcSlice.GetTileOffstForMultES() + uiAccumulatedSubstreamLength - uiBitsOriginallyInSubstreams) >> 3) + uiCounter)
            }
        }

        pcCU := rpcPic.GetCU(uiCUAddr)
        if pcSlice.GetSPS().GetUseSAO() && (pcSlice.GetSaoEnabledFlag() || pcSlice.GetSaoEnabledFlagChroma()) {
            saoParam := pcSlice.GetPic().GetPicSym().GetSaoParam()
            iNumCuInWidth := saoParam.NumCuInWidth
            iCUAddrInSlice := int(uiCUAddr - rpcPic.GetPicSym().GetCUOrderMap(int(pcSlice.GetSliceCurStartCUAddr()/rpcPic.GetNumPartInCU())))
            iCUAddrUpInSlice := iCUAddrInSlice - iNumCuInWidth
            rx := int(uiCUAddr) % iNumCuInWidth
            ry := int(uiCUAddr) / iNumCuInWidth
            allowMergeLeft := true
            allowMergeUp := true
            if rx != 0 {
                if rpcPic.GetPicSym().GetTileIdxMap(int(uiCUAddr-1)) != rpcPic.GetPicSym().GetTileIdxMap(int(uiCUAddr)) {
                    allowMergeLeft = false
                }
            }
            if ry != 0 {
                if rpcPic.GetPicSym().GetTileIdxMap(int(uiCUAddr)-iNumCuInWidth) != rpcPic.GetPicSym().GetTileIdxMap(int(uiCUAddr)) {
                    allowMergeUp = false
                }
            }
            addr := pcCU.GetAddr()
            allowMergeLeft = allowMergeLeft && (rx > 0) && (iCUAddrInSlice != 0)
            allowMergeUp = allowMergeUp && (ry > 0) && (iCUAddrUpInSlice >= 0)
            if saoParam.SaoFlag[0] || saoParam.SaoFlag[1] {
                mergeLeft := saoParam.SaoLcuParam[0][addr].MergeLeftFlag
                mergeUp := saoParam.SaoLcuParam[0][addr].MergeUpFlag
                if allowMergeLeft {
                    this.m_pcEntropyCoder.m_pcEntropyCoderIf.codeSaoMerge(uint(TLibCommon.B2U(mergeLeft)))
                } else {
                    mergeLeft = false
                }
                if mergeLeft == false {
                    if allowMergeUp {
                        this.m_pcEntropyCoder.m_pcEntropyCoderIf.codeSaoMerge(uint(TLibCommon.B2U(mergeUp)))
                    } else {
                        mergeUp = false
                    }
                    if mergeUp == false {
                        for compIdx := uint(0); compIdx < 3; compIdx++ {
                            if (compIdx == 0 && saoParam.SaoFlag[0]) || (compIdx > 0 && saoParam.SaoFlag[1]) {
                                this.m_pcEntropyCoder.encodeSaoOffset(&saoParam.SaoLcuParam[compIdx][addr], compIdx)
                            }
                        }
                    }
                }
            }
        } else if pcSlice.GetSPS().GetUseSAO() {
            addr := pcCU.GetAddr()
            saoParam := pcSlice.GetPic().GetPicSym().GetSaoParam()
            for cIdx := 0; cIdx < 3; cIdx++ {
                saoLcuParam := &(saoParam.SaoLcuParam[cIdx][addr])
                if ((cIdx == 0) && !pcSlice.GetSaoEnabledFlag()) || ((cIdx == 1 || cIdx == 2) && !pcSlice.GetSaoEnabledFlagChroma()) {
                    saoLcuParam.MergeUpFlag = false
                    saoLcuParam.MergeLeftFlag = false
                    saoLcuParam.SubTypeIdx = 0
                    saoLcuParam.TypeIdx = -1
                    saoLcuParam.Offset[0] = 0
                    saoLcuParam.Offset[1] = 0
                    saoLcuParam.Offset[2] = 0
                    saoLcuParam.Offset[3] = 0
                }
            }
        }
        //#if ENC_DEC_TRACE
        //    TLibCommon.G_bJustDoIt = TLibCommon.G_bEncDecTraceEnable;
        //#endif
        if (this.m_pcCfg.GetSliceMode() != 0 || this.m_pcCfg.GetSliceSegmentMode() != 0) &&
            uiCUAddr == rpcPic.GetPicSym().GetCUOrderMap(int((uiBoundingCUAddr+rpcPic.GetNumPartInCU()-1)/rpcPic.GetNumPartInCU()-1)) {
            this.m_pcCuEncoder.encodeCU(pcCU)
        } else {
            this.m_pcCuEncoder.encodeCU(pcCU)
        }
        //#if ENC_DEC_TRACE
        //    TLibCommon.G_bJustDoIt = TLibCommon.G_bEncDecTraceDisable;
        //#endif
        if this.m_pcCfg.GetUseSBACRD() {
            pcSbacCoders[uiSubStrm].load(this.m_pcSbacCoder) //load back status of the entropy coder after encoding the LCU into relevant bitstream entropy coder

            //Store probabilties of second LCU in line into buffer
            if (depSliceSegmentsEnabled || (pcSlice.GetPPS().GetNumSubstreams() > 1)) && (uiCol == uiTileLCUX+1) && this.m_pcCfg.GetWaveFrontsynchro() != 0 {
                this.m_pcBufferSbacCoders[uiTileCol].loadContexts(pcSbacCoders[uiSubStrm])
            }
        }        
    }

    if depSliceSegmentsEnabled {
        if this.m_pcCfg.GetWaveFrontsynchro() != 0 {
            this.CTXMem[1].loadContexts(this.m_pcBufferSbacCoders[uiTileCol]) //ctx 2.LCU
        }
        this.CTXMem[0].loadContexts(this.m_pcSbacCoder) //ctx end of dep.slice
    }

    //#if ADAPTIVE_QP_SELECTION
    if this.m_pcCfg.GetUseAdaptQpSelect() {
        this.m_pcTrQuant.StoreSliceQpNext(pcSlice)
    }
    //#endif
    if pcSlice.GetPPS().GetCabacInitPresentFlag() {
        if pcSlice.GetPPS().GetDependentSliceSegmentsEnabledFlag() {
            pcSlice.GetPPS().SetEncCABACTableIdx(uint(pcSlice.GetSliceType()))
        } else {
            this.m_pcEntropyCoder.determineCabacInitIdx()
        }
    }
}

// misc. functions
func (this *TEncSlice) setSearchRange(pcSlice *TLibCommon.TComSlice) { ///< set ME range adaptively
    iCurrPOC := pcSlice.GetPOC()
    iRefPOC := 0
    iGOPSize := this.m_pcCfg.GetGOPSize()
    iOffset := (iGOPSize >> 1)
    iMaxSR := this.m_pcCfg.GetSearchRange()
    var iNumPredDir int
    if pcSlice.IsInterP() {
        iNumPredDir = 1
    } else {
        iNumPredDir = 2
    }

    for iDir := 0; iDir <= iNumPredDir; iDir++ {
        //RefPicList e = (RefPicList)iDir;
        var e TLibCommon.RefPicList
        if iDir != 0 {
            e = TLibCommon.REF_PIC_LIST_1
        } else {
            e = TLibCommon.REF_PIC_LIST_0
        }
        for iRefIdx := 0; iRefIdx < pcSlice.GetNumRefIdx(e); iRefIdx++ {
            iRefPOC = int(pcSlice.GetRefPic(e, iRefIdx).GetPOC())
            iNewSR := TLibCommon.CLIP3(8, iMaxSR, (iMaxSR*TLibCommon.ADAPT_SR_SCALE*TLibCommon.ABS(iCurrPOC-iRefPOC).(int)+iOffset)/iGOPSize).(int)
            this.m_pcPredSearch.setAdaptiveSearchRange(iDir, iRefIdx, iNewSR)
        }
    }
}
func (this *TEncSlice) getTotalBits() uint64 { return this.m_uiPicTotalBits }

func (this *TEncSlice) getCUEncoder() *TEncCu { return this.m_pcCuEncoder } ///< CU encoder
func (this *TEncSlice) xDetermineStartAndBoundingCUAddr(startCUAddr, boundingCUAddr *uint, rpcPic *TLibCommon.TComPic, bEncodeSlice bool) {
    pcSlice := rpcPic.GetSlice(this.getSliceIdx())
    var uiStartCUAddrSlice, uiBoundingCUAddrSlice uint
    var tileIdxIncrement, tileIdx, tileWidthInLcu, tileHeightInLcu, tileTotalCount uint

    uiStartCUAddrSlice = pcSlice.GetSliceCurStartCUAddr()
    uiNumberOfCUsInFrame := rpcPic.GetNumCUsInFrame()
    uiBoundingCUAddrSlice = uiNumberOfCUsInFrame
    if bEncodeSlice {
        var uiCUAddrIncrement uint
        switch this.m_pcCfg.GetSliceMode() {
        case TLibCommon.FIXED_NUMBER_OF_LCU:
            uiCUAddrIncrement = uint(this.m_pcCfg.GetSliceArgument())
            if (uiStartCUAddrSlice + uiCUAddrIncrement) < uiNumberOfCUsInFrame*rpcPic.GetNumPartInCU() {
                uiBoundingCUAddrSlice = (uiStartCUAddrSlice + uiCUAddrIncrement)
            } else {
                uiBoundingCUAddrSlice = uiNumberOfCUsInFrame * rpcPic.GetNumPartInCU()
            }

        case TLibCommon.FIXED_NUMBER_OF_BYTES:
            uiCUAddrIncrement = rpcPic.GetNumCUsInFrame()
            uiBoundingCUAddrSlice = pcSlice.GetSliceCurEndCUAddr()

        case TLibCommon.FIXED_NUMBER_OF_TILES:
            tileIdx = rpcPic.GetPicSym().GetTileIdxMap(int(rpcPic.GetPicSym().GetCUOrderMap(int(uiStartCUAddrSlice / rpcPic.GetNumPartInCU()))))
            uiCUAddrIncrement = 0
            tileTotalCount = uint(rpcPic.GetPicSym().GetNumColumnsMinus1()+1) * uint(rpcPic.GetPicSym().GetNumRowsMinus1()+1)

            for tileIdxIncrement = 0; tileIdxIncrement < uint(this.m_pcCfg.GetSliceArgument()); tileIdxIncrement++ {
                if (tileIdx + tileIdxIncrement) < tileTotalCount {
                    tileWidthInLcu = rpcPic.GetPicSym().GetTComTile(tileIdx + tileIdxIncrement).GetTileWidth()
                    tileHeightInLcu = rpcPic.GetPicSym().GetTComTile(tileIdx + tileIdxIncrement).GetTileHeight()
                    uiCUAddrIncrement += (tileWidthInLcu * tileHeightInLcu * rpcPic.GetNumPartInCU())
                }
            }
            if (uiStartCUAddrSlice + uiCUAddrIncrement) < uiNumberOfCUsInFrame*rpcPic.GetNumPartInCU() {
                uiBoundingCUAddrSlice = (uiStartCUAddrSlice + uiCUAddrIncrement)
            } else {
                uiBoundingCUAddrSlice = uiNumberOfCUsInFrame * rpcPic.GetNumPartInCU()
            }
        default:
            uiCUAddrIncrement = rpcPic.GetNumCUsInFrame()
            uiBoundingCUAddrSlice = uiNumberOfCUsInFrame * rpcPic.GetNumPartInCU()

        }
        // WPP: if a slice does not start at the beginning of a CTB row, it must end within the same CTB row
        if pcSlice.GetPPS().GetNumSubstreams() > 1 && (uiStartCUAddrSlice%(rpcPic.GetFrameWidthInCU()*rpcPic.GetNumPartInCU()) != 0) {
            uiBoundingCUAddrSlice = TLibCommon.MIN(uiBoundingCUAddrSlice, uiStartCUAddrSlice-(uiStartCUAddrSlice%(rpcPic.GetFrameWidthInCU()*rpcPic.GetNumPartInCU()))+(rpcPic.GetFrameWidthInCU()*rpcPic.GetNumPartInCU())).(uint)
        }
        pcSlice.SetSliceCurEndCUAddr(uiBoundingCUAddrSlice)
    } else {
        var uiCUAddrIncrement uint
        switch this.m_pcCfg.GetSliceMode() {
        case TLibCommon.FIXED_NUMBER_OF_LCU:
            uiCUAddrIncrement = uint(this.m_pcCfg.GetSliceArgument())
            if (uiStartCUAddrSlice + uiCUAddrIncrement) < uiNumberOfCUsInFrame*rpcPic.GetNumPartInCU() {
                uiBoundingCUAddrSlice = (uiStartCUAddrSlice + uiCUAddrIncrement)
            } else {
                uiBoundingCUAddrSlice = uiNumberOfCUsInFrame * rpcPic.GetNumPartInCU()
            }
        case TLibCommon.FIXED_NUMBER_OF_TILES:
            tileIdx = rpcPic.GetPicSym().GetTileIdxMap(int(rpcPic.GetPicSym().GetCUOrderMap(int(uiStartCUAddrSlice / rpcPic.GetNumPartInCU()))))
            uiCUAddrIncrement = 0
            tileTotalCount = uint(rpcPic.GetPicSym().GetNumColumnsMinus1()+1) * uint(rpcPic.GetPicSym().GetNumRowsMinus1()+1)

            for tileIdxIncrement = 0; tileIdxIncrement < uint(this.m_pcCfg.GetSliceArgument()); tileIdxIncrement++ {
                if (tileIdx + tileIdxIncrement) < tileTotalCount {
                    tileWidthInLcu = rpcPic.GetPicSym().GetTComTile(tileIdx + tileIdxIncrement).GetTileWidth()
                    tileHeightInLcu = rpcPic.GetPicSym().GetTComTile(tileIdx + tileIdxIncrement).GetTileHeight()
                    uiCUAddrIncrement += (tileWidthInLcu * tileHeightInLcu * rpcPic.GetNumPartInCU())
                }
            }
            if (uiStartCUAddrSlice + uiCUAddrIncrement) < uiNumberOfCUsInFrame*rpcPic.GetNumPartInCU() {
                uiBoundingCUAddrSlice = (uiStartCUAddrSlice + uiCUAddrIncrement)
            } else {
                uiBoundingCUAddrSlice = uiNumberOfCUsInFrame * rpcPic.GetNumPartInCU()
            }

        default:
            uiCUAddrIncrement = rpcPic.GetNumCUsInFrame()
            uiBoundingCUAddrSlice = uiNumberOfCUsInFrame * rpcPic.GetNumPartInCU()

        }
        // WPP: if a slice does not start at the beginning of a CTB row, it must end within the same CTB row
        if pcSlice.GetPPS().GetNumSubstreams() > 1 && (uiStartCUAddrSlice%(rpcPic.GetFrameWidthInCU()*rpcPic.GetNumPartInCU()) != 0) {
            uiBoundingCUAddrSlice = TLibCommon.MIN(uiBoundingCUAddrSlice, uiStartCUAddrSlice-(uiStartCUAddrSlice%(rpcPic.GetFrameWidthInCU()*rpcPic.GetNumPartInCU()))+(rpcPic.GetFrameWidthInCU()*rpcPic.GetNumPartInCU())).(uint)
        }
        pcSlice.SetSliceCurEndCUAddr(uiBoundingCUAddrSlice)
    }

    tileBoundary := false
    if (this.m_pcCfg.GetSliceMode() == TLibCommon.FIXED_NUMBER_OF_LCU ||
        this.m_pcCfg.GetSliceMode() == TLibCommon.FIXED_NUMBER_OF_BYTES) &&
        (this.m_pcCfg.GetNumRowsMinus1() > 0 || this.m_pcCfg.GetNumColumnsMinus1() > 0) {
        lcuEncAddr := (uiStartCUAddrSlice + rpcPic.GetNumPartInCU() - 1) / rpcPic.GetNumPartInCU()
        lcuAddr := rpcPic.GetPicSym().GetCUOrderMap(int(lcuEncAddr))
        startTileIdx := rpcPic.GetPicSym().GetTileIdxMap(int(lcuAddr))
        tileBoundingCUAddrSlice := uint(0)
        for lcuEncAddr < uiNumberOfCUsInFrame && rpcPic.GetPicSym().GetTileIdxMap(int(lcuAddr)) == startTileIdx {
            lcuEncAddr++
            lcuAddr = rpcPic.GetPicSym().GetCUOrderMap(int(lcuEncAddr))
        }
        tileBoundingCUAddrSlice = lcuEncAddr * rpcPic.GetNumPartInCU()

        if tileBoundingCUAddrSlice < uiBoundingCUAddrSlice {
            uiBoundingCUAddrSlice = tileBoundingCUAddrSlice
            pcSlice.SetSliceCurEndCUAddr(uiBoundingCUAddrSlice)
            tileBoundary = true
        }
    }

    // Dependent slice
    var startCUAddrSliceSegment, boundingCUAddrSliceSegment uint
    startCUAddrSliceSegment = pcSlice.GetSliceSegmentCurStartCUAddr()
    boundingCUAddrSliceSegment = uiNumberOfCUsInFrame
    if bEncodeSlice {
        var uiCUAddrIncrement uint
        switch this.m_pcCfg.GetSliceSegmentMode() {
        case TLibCommon.FIXED_NUMBER_OF_LCU:
            uiCUAddrIncrement = uint(this.m_pcCfg.GetSliceSegmentArgument())
            if (startCUAddrSliceSegment + uiCUAddrIncrement) < uiNumberOfCUsInFrame*rpcPic.GetNumPartInCU() {
                boundingCUAddrSliceSegment = (startCUAddrSliceSegment + uiCUAddrIncrement)
            } else {
                boundingCUAddrSliceSegment = uiNumberOfCUsInFrame * rpcPic.GetNumPartInCU()
            }

        case TLibCommon.FIXED_NUMBER_OF_BYTES:
            uiCUAddrIncrement = rpcPic.GetNumCUsInFrame()
            boundingCUAddrSliceSegment = pcSlice.GetSliceSegmentCurEndCUAddr()

        case TLibCommon.FIXED_NUMBER_OF_TILES:
            tileIdx = rpcPic.GetPicSym().GetTileIdxMap(int(rpcPic.GetPicSym().GetCUOrderMap(int(pcSlice.GetSliceSegmentCurStartCUAddr() / rpcPic.GetNumPartInCU()))))
            uiCUAddrIncrement = 0
            tileTotalCount = uint(rpcPic.GetPicSym().GetNumColumnsMinus1()+1) * uint(rpcPic.GetPicSym().GetNumRowsMinus1()+1)

            for tileIdxIncrement = 0; tileIdxIncrement < uint(this.m_pcCfg.GetSliceSegmentArgument()); tileIdxIncrement++ {
                if (tileIdx + tileIdxIncrement) < tileTotalCount {
                    tileWidthInLcu = rpcPic.GetPicSym().GetTComTile(tileIdx + tileIdxIncrement).GetTileWidth()
                    tileHeightInLcu = rpcPic.GetPicSym().GetTComTile(tileIdx + tileIdxIncrement).GetTileHeight()
                    uiCUAddrIncrement += (tileWidthInLcu * tileHeightInLcu * rpcPic.GetNumPartInCU())
                }
            }
            if (startCUAddrSliceSegment + uiCUAddrIncrement) < uiNumberOfCUsInFrame*rpcPic.GetNumPartInCU() {
                boundingCUAddrSliceSegment = (startCUAddrSliceSegment + uiCUAddrIncrement)
            } else {
                boundingCUAddrSliceSegment = uiNumberOfCUsInFrame * rpcPic.GetNumPartInCU()
            }
        default:
            uiCUAddrIncrement = rpcPic.GetNumCUsInFrame()
            boundingCUAddrSliceSegment = uiNumberOfCUsInFrame * rpcPic.GetNumPartInCU()

        }
        // WPP: if a slice segment does not start at the beginning of a CTB row, it must end within the same CTB row
        if pcSlice.GetPPS().GetNumSubstreams() > 1 && (startCUAddrSliceSegment%(rpcPic.GetFrameWidthInCU()*rpcPic.GetNumPartInCU()) != 0) {
            boundingCUAddrSliceSegment = TLibCommon.MIN(boundingCUAddrSliceSegment, startCUAddrSliceSegment-(startCUAddrSliceSegment%(rpcPic.GetFrameWidthInCU()*rpcPic.GetNumPartInCU()))+(rpcPic.GetFrameWidthInCU()*rpcPic.GetNumPartInCU())).(uint)
        }
        pcSlice.SetSliceSegmentCurEndCUAddr(boundingCUAddrSliceSegment)
    } else {
        var uiCUAddrIncrement uint
        switch this.m_pcCfg.GetSliceSegmentMode() {
        case TLibCommon.FIXED_NUMBER_OF_LCU:
            uiCUAddrIncrement = uint(this.m_pcCfg.GetSliceSegmentArgument())
            if (startCUAddrSliceSegment + uiCUAddrIncrement) < uiNumberOfCUsInFrame*rpcPic.GetNumPartInCU() {
                boundingCUAddrSliceSegment = (startCUAddrSliceSegment + uiCUAddrIncrement)
            } else {
                boundingCUAddrSliceSegment = uiNumberOfCUsInFrame * rpcPic.GetNumPartInCU()
            }
        case TLibCommon.FIXED_NUMBER_OF_TILES:
            tileIdx = rpcPic.GetPicSym().GetTileIdxMap(int(rpcPic.GetPicSym().GetCUOrderMap(int(pcSlice.GetSliceSegmentCurStartCUAddr() / rpcPic.GetNumPartInCU()))))
            uiCUAddrIncrement = 0
            tileTotalCount = uint(rpcPic.GetPicSym().GetNumColumnsMinus1()+1) * uint(rpcPic.GetPicSym().GetNumRowsMinus1()+1)

            for tileIdxIncrement = 0; tileIdxIncrement < uint(this.m_pcCfg.GetSliceSegmentArgument()); tileIdxIncrement++ {
                if (tileIdx + tileIdxIncrement) < tileTotalCount {
                    tileWidthInLcu = rpcPic.GetPicSym().GetTComTile(tileIdx + tileIdxIncrement).GetTileWidth()
                    tileHeightInLcu = rpcPic.GetPicSym().GetTComTile(tileIdx + tileIdxIncrement).GetTileHeight()
                    uiCUAddrIncrement += (tileWidthInLcu * tileHeightInLcu * rpcPic.GetNumPartInCU())
                }
            }
            if (startCUAddrSliceSegment + uiCUAddrIncrement) < uiNumberOfCUsInFrame*rpcPic.GetNumPartInCU() {
                boundingCUAddrSliceSegment = (startCUAddrSliceSegment + uiCUAddrIncrement)
            } else {
                boundingCUAddrSliceSegment = uiNumberOfCUsInFrame * rpcPic.GetNumPartInCU()
            }
        default:
            uiCUAddrIncrement = rpcPic.GetNumCUsInFrame()
            boundingCUAddrSliceSegment = uiNumberOfCUsInFrame * rpcPic.GetNumPartInCU()
        }
        // WPP: if a slice segment does not start at the beginning of a CTB row, it must end within the same CTB row
        if pcSlice.GetPPS().GetNumSubstreams() > 1 && (startCUAddrSliceSegment%(rpcPic.GetFrameWidthInCU()*rpcPic.GetNumPartInCU()) != 0) {
            boundingCUAddrSliceSegment = uint(TLibCommon.MIN(int(boundingCUAddrSliceSegment), int(startCUAddrSliceSegment-(startCUAddrSliceSegment%(rpcPic.GetFrameWidthInCU()*rpcPic.GetNumPartInCU()))+(rpcPic.GetFrameWidthInCU()*rpcPic.GetNumPartInCU()))).(int))
        }
        pcSlice.SetSliceSegmentCurEndCUAddr(boundingCUAddrSliceSegment)
    }
    if (this.m_pcCfg.GetSliceSegmentMode() == TLibCommon.FIXED_NUMBER_OF_LCU ||
        this.m_pcCfg.GetSliceSegmentMode() == TLibCommon.FIXED_NUMBER_OF_BYTES) &&
        (this.m_pcCfg.GetNumRowsMinus1() > 0 ||
            this.m_pcCfg.GetNumColumnsMinus1() > 0) {
        lcuEncAddr := uint(startCUAddrSliceSegment+rpcPic.GetNumPartInCU()-1) / rpcPic.GetNumPartInCU()
        lcuAddr := uint(rpcPic.GetPicSym().GetCUOrderMap(int(lcuEncAddr)))
        startTileIdx := uint(rpcPic.GetPicSym().GetTileIdxMap(int(lcuAddr)))
        tileBoundingCUAddrSlice := uint(0)
        for lcuEncAddr < uiNumberOfCUsInFrame && rpcPic.GetPicSym().GetTileIdxMap(int(lcuAddr)) == startTileIdx {
            lcuEncAddr++
            lcuAddr = rpcPic.GetPicSym().GetCUOrderMap(int(lcuEncAddr))
        }
        tileBoundingCUAddrSlice = lcuEncAddr * rpcPic.GetNumPartInCU()

        if tileBoundingCUAddrSlice < boundingCUAddrSliceSegment {
            boundingCUAddrSliceSegment = tileBoundingCUAddrSlice
            pcSlice.SetSliceSegmentCurEndCUAddr(boundingCUAddrSliceSegment)
            tileBoundary = true
        }
    }

    if boundingCUAddrSliceSegment > uiBoundingCUAddrSlice {
        boundingCUAddrSliceSegment = uiBoundingCUAddrSlice
        pcSlice.SetSliceSegmentCurEndCUAddr(uiBoundingCUAddrSlice)
    }
    //calculate real dependent slice start address
    uiInternalAddress := rpcPic.GetPicSym().GetPicSCUAddr(pcSlice.GetSliceSegmentCurStartCUAddr()) % rpcPic.GetNumPartInCU()
    uiExternalAddress := rpcPic.GetPicSym().GetPicSCUAddr(pcSlice.GetSliceSegmentCurStartCUAddr()) / rpcPic.GetNumPartInCU()
    uiPosX := (uiExternalAddress%rpcPic.GetFrameWidthInCU())*pcSlice.GetSPS().GetMaxCUWidth() + TLibCommon.G_auiRasterToPelX[TLibCommon.G_auiZscanToRaster[uiInternalAddress]]
    uiPosY := (uiExternalAddress/rpcPic.GetFrameWidthInCU())*pcSlice.GetSPS().GetMaxCUHeight() + TLibCommon.G_auiRasterToPelY[TLibCommon.G_auiZscanToRaster[uiInternalAddress]]
    uiWidth := pcSlice.GetSPS().GetPicWidthInLumaSamples()
    uiHeight := pcSlice.GetSPS().GetPicHeightInLumaSamples()
    for (uiPosX >= uiWidth || uiPosY >= uiHeight) && !(uiPosX >= uiWidth && uiPosY >= uiHeight) {
        uiInternalAddress++
        if uiInternalAddress >= rpcPic.GetNumPartInCU() {
            uiInternalAddress = 0
            uiExternalAddress = rpcPic.GetPicSym().GetCUOrderMap(int(rpcPic.GetPicSym().GetInverseCUOrderMap(int(uiExternalAddress)) + 1))
        }
        uiPosX = (uiExternalAddress%rpcPic.GetFrameWidthInCU())*pcSlice.GetSPS().GetMaxCUWidth() + TLibCommon.G_auiRasterToPelX[TLibCommon.G_auiZscanToRaster[uiInternalAddress]]
        uiPosY = (uiExternalAddress/rpcPic.GetFrameWidthInCU())*pcSlice.GetSPS().GetMaxCUHeight() + TLibCommon.G_auiRasterToPelY[TLibCommon.G_auiZscanToRaster[uiInternalAddress]]
    }
    uiRealStartAddress := rpcPic.GetPicSym().GetPicSCUEncOrder(uiExternalAddress*rpcPic.GetNumPartInCU() + uiInternalAddress)

    pcSlice.SetSliceSegmentCurStartCUAddr(uiRealStartAddress)
    startCUAddrSliceSegment = uiRealStartAddress

    //calculate real slice start address
    uiInternalAddress = rpcPic.GetPicSym().GetPicSCUAddr(pcSlice.GetSliceSegmentCurStartCUAddr()) % rpcPic.GetNumPartInCU()
    uiExternalAddress = rpcPic.GetPicSym().GetPicSCUAddr(pcSlice.GetSliceSegmentCurStartCUAddr()) / rpcPic.GetNumPartInCU()
    uiPosX = (uiExternalAddress%rpcPic.GetFrameWidthInCU())*pcSlice.GetSPS().GetMaxCUWidth() + TLibCommon.G_auiRasterToPelX[TLibCommon.G_auiZscanToRaster[uiInternalAddress]]
    uiPosY = (uiExternalAddress/rpcPic.GetFrameWidthInCU())*pcSlice.GetSPS().GetMaxCUHeight() + TLibCommon.G_auiRasterToPelY[TLibCommon.G_auiZscanToRaster[uiInternalAddress]]
    uiWidth = pcSlice.GetSPS().GetPicWidthInLumaSamples()
    uiHeight = pcSlice.GetSPS().GetPicHeightInLumaSamples()
    for (uiPosX >= uiWidth || uiPosY >= uiHeight) && !(uiPosX >= uiWidth && uiPosY >= uiHeight) {
        uiInternalAddress++
        if uiInternalAddress >= rpcPic.GetNumPartInCU() {
            uiInternalAddress = 0
            uiExternalAddress = rpcPic.GetPicSym().GetCUOrderMap(int(rpcPic.GetPicSym().GetInverseCUOrderMap(int(uiExternalAddress)) + 1))
        }
        uiPosX = (uiExternalAddress%rpcPic.GetFrameWidthInCU())*pcSlice.GetSPS().GetMaxCUWidth() + TLibCommon.G_auiRasterToPelX[TLibCommon.G_auiZscanToRaster[uiInternalAddress]]
        uiPosY = (uiExternalAddress/rpcPic.GetFrameWidthInCU())*pcSlice.GetSPS().GetMaxCUHeight() + TLibCommon.G_auiRasterToPelY[TLibCommon.G_auiZscanToRaster[uiInternalAddress]]
    }
    uiRealStartAddress = rpcPic.GetPicSym().GetPicSCUEncOrder(uiExternalAddress*rpcPic.GetNumPartInCU() + uiInternalAddress)

    pcSlice.SetSliceSegmentCurStartCUAddr(uiRealStartAddress)
    startCUAddrSliceSegment = uiRealStartAddress

    // Make a joint decision based on reconstruction and dependent slice bounds
    *startCUAddr = TLibCommon.MAX(uiStartCUAddrSlice, startCUAddrSliceSegment).(uint)
    *boundingCUAddr = TLibCommon.MIN(uiBoundingCUAddrSlice, boundingCUAddrSliceSegment).(uint)

    if !bEncodeSlice {
        // For fixed number of LCU within an entropy and reconstruction slice we already know whether we will encounter end of entropy and/or reconstruction slice
        // first. Set the flags accordingly.
        if (this.m_pcCfg.GetSliceMode() == TLibCommon.FIXED_NUMBER_OF_LCU && this.m_pcCfg.GetSliceSegmentMode() == TLibCommon.FIXED_NUMBER_OF_LCU) ||
            (this.m_pcCfg.GetSliceMode() == 0 && this.m_pcCfg.GetSliceSegmentMode() == TLibCommon.FIXED_NUMBER_OF_LCU) ||
            (this.m_pcCfg.GetSliceMode() == TLibCommon.FIXED_NUMBER_OF_LCU && this.m_pcCfg.GetSliceSegmentMode() == 0) ||
            (this.m_pcCfg.GetSliceMode() == TLibCommon.FIXED_NUMBER_OF_TILES && this.m_pcCfg.GetSliceSegmentMode() == TLibCommon.FIXED_NUMBER_OF_LCU) ||
            (this.m_pcCfg.GetSliceMode() == TLibCommon.FIXED_NUMBER_OF_TILES && this.m_pcCfg.GetSliceSegmentMode() == 0) ||
            (this.m_pcCfg.GetSliceSegmentMode() == TLibCommon.FIXED_NUMBER_OF_TILES && this.m_pcCfg.GetSliceMode() == 0) ||
            tileBoundary {
            if uiBoundingCUAddrSlice < boundingCUAddrSliceSegment {
                pcSlice.SetNextSlice(true)
                pcSlice.SetNextSliceSegment(false)
            } else if uiBoundingCUAddrSlice > boundingCUAddrSliceSegment {
                pcSlice.SetNextSlice(false)
                pcSlice.SetNextSliceSegment(true)
            } else {
                pcSlice.SetNextSlice(true)
                pcSlice.SetNextSliceSegment(true)
            }
        } else {
            pcSlice.SetNextSlice(false)
            pcSlice.SetNextSliceSegment(false)
        }
    }
}

func (this *TEncSlice) getSliceIdx() uint  { return this.m_uiSliceIdx }
func (this *TEncSlice) setSliceIdx(i uint) { this.m_uiSliceIdx = i }

func (this *TEncSlice) initCtxMem(i uint) {
    //for j := CTXMem.Front(); j != nil; j=j.Next() {
    //i := j.Value.(*TEncSbac);
    //delete (*j);
    //}
    //this.CTXMem.Init();
    this.CTXMem = make(map[int]*TEncSbac, i)
}

func (this *TEncSlice) setCtxMem(sb *TEncSbac, b int) { this.CTXMem[b] = sb }

//private:
func (this *TEncSlice) xGetQPValueAccordingToLambda(lambda float64) float64 {
    return 4.2005*math.Log(lambda) + 13.7122
}
