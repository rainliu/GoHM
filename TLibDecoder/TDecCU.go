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

package TLibDecoder

import (
    "fmt"
    "gohm/TLibCommon"
)

// ====================================================================================================================
// Class definition
// ====================================================================================================================

/// CU decoder class
type TDecCu struct {
    //private:
    m_uiMaxDepth uint                     ///< max. number of depth
    m_ppcYuvResi []*TLibCommon.TComYuv    ///< array of residual buffer
    m_ppcYuvReco []*TLibCommon.TComYuv    ///< array of prediction & reconstruction buffer
    m_ppcCU      []*TLibCommon.TComDataCU ///< CU data array

    // access channel
    m_pcTrQuant        *TLibCommon.TComTrQuant
    m_pcPrediction     *TLibCommon.TComPrediction
    m_pcEntropyDecoder *TDecEntropy

    m_bDecodeDQP bool
}

func NewTDecCu() *TDecCu {
    return &TDecCu{}
}

/// initialize access channels
func (this *TDecCu) Init(pcEntropyDecoder *TDecEntropy, pcTrQuant *TLibCommon.TComTrQuant, pcPrediction *TLibCommon.TComPrediction) {
    this.m_pcEntropyDecoder = pcEntropyDecoder
    this.m_pcTrQuant = pcTrQuant
    this.m_pcPrediction = pcPrediction
}

/// create internal buffers
func (this *TDecCu) Create(uiMaxDepth, uiMaxWidth, uiMaxHeight uint) {
	//fmt.Printf("uiMaxDepth=%d, uiMaxWidth=%d,uiMaxHeight=%d\n",uiMaxDepth,  uiMaxWidth,  uiMaxHeight);
	
    this.m_uiMaxDepth = uiMaxDepth + 1

    this.m_ppcYuvResi = make([]*TLibCommon.TComYuv, this.m_uiMaxDepth-1)
    this.m_ppcYuvReco = make([]*TLibCommon.TComYuv, this.m_uiMaxDepth-1)
    this.m_ppcCU = make([]*TLibCommon.TComDataCU, this.m_uiMaxDepth-1)

    var uiNumPartitions uint
    for ui := uint(0); ui < this.m_uiMaxDepth-1; ui++ {
        uiNumPartitions = 1 << ((this.m_uiMaxDepth - ui - 1) << 1)
        uiWidth := (uiMaxWidth >> ui)
        uiHeight := (uiMaxHeight >> ui)
		//fmt.Printf("uiWidth=%d, uiHeight=%d\n",uiWidth,  uiHeight);
	
        this.m_ppcYuvResi[ui] = TLibCommon.NewTComYuv()
        this.m_ppcYuvResi[ui].Create(uiWidth, uiHeight)
        this.m_ppcYuvReco[ui] = TLibCommon.NewTComYuv()
        this.m_ppcYuvReco[ui].Create(uiWidth, uiHeight)
        this.m_ppcCU[ui] = TLibCommon.NewTComDataCU()
        this.m_ppcCU[ui].Create(uiNumPartitions, uiWidth, uiHeight, true, int(uiMaxWidth>>(this.m_uiMaxDepth-1)), false)
    }

    this.m_bDecodeDQP = false

    // initialize partition order.
    piTmp := uint(0)
    TLibCommon.InitZscanToRaster(int(this.m_uiMaxDepth), 1, 0, TLibCommon.G_auiZscanToRaster[:], &piTmp)
    TLibCommon.InitRasterToZscan(uiMaxWidth, uiMaxHeight, this.m_uiMaxDepth)

    // initialize conversion matrix from partition index to pel
    TLibCommon.InitRasterToPelXY(uiMaxWidth, uiMaxHeight, this.m_uiMaxDepth)
}

/// destroy internal buffers
func (this *TDecCu) Destroy() {
    for ui := uint(0); ui < this.m_uiMaxDepth-1; ui++ {
        this.m_ppcYuvResi[ui].Destroy()
        //delete m_ppcYuvResi[ui];
        this.m_ppcYuvResi[ui] = nil
        this.m_ppcYuvReco[ui].Destroy()
        //delete m_ppcYuvReco[ui];
        this.m_ppcYuvReco[ui] = nil
        this.m_ppcCU[ui].Destroy()
        //delete m_ppcCU     [ui];
        this.m_ppcCU[ui] = nil
    }

    //delete [] m_ppcYuvResi;
    this.m_ppcYuvResi = nil
    //delete [] m_ppcYuvReco;
    this.m_ppcYuvReco = nil
    //delete [] m_ppcCU     ;
    this.m_ppcCU = nil
}

/// decode CU information
func (this *TDecCu) DecodeCU(pcCU *TLibCommon.TComDataCU, ruiIsLast *uint) {
    if pcCU.GetSlice().GetPPS().GetUseDQP() {
        this.SetdQPFlag(true)
    }

    // start from the top level CU
    this.xDecodeCU(pcCU, 0, 0, ruiIsLast)
}

/// reconstruct CU information
func (this *TDecCu) DecompressCU(pcCU *TLibCommon.TComDataCU) {
    this.xDecompressCU(pcCU, 0, 0)
}

func (this *TDecCu) xDecodeCU(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint, ruiIsLast *uint) {
    pcPic := pcCU.GetPic()
    uiCurNumParts := pcPic.GetNumPartInCU() >> (uiDepth << 1)
    uiQNumParts := uiCurNumParts >> 2

    bBoundary := false
    uiLPelX := pcCU.GetCUPelX() + TLibCommon.G_auiRasterToPelX[TLibCommon.G_auiZscanToRaster[uiAbsPartIdx]]
    uiRPelX := uiLPelX + (pcCU.GetSlice().GetSPS().GetMaxCUWidth() >> uiDepth) - 1
    uiTPelY := pcCU.GetCUPelY() + TLibCommon.G_auiRasterToPelY[TLibCommon.G_auiZscanToRaster[uiAbsPartIdx]]
    uiBPelY := uiTPelY + (pcCU.GetSlice().GetSPS().GetMaxCUHeight() >> uiDepth) - 1

    pcSlice := pcCU.GetPic().GetSlice(pcCU.GetPic().GetCurrSliceIdx())
    bStartInCU := pcCU.GetSCUAddr()+uiAbsPartIdx+uiCurNumParts > pcSlice.GetSliceSegmentCurStartCUAddr() && pcCU.GetSCUAddr()+uiAbsPartIdx < pcSlice.GetSliceSegmentCurStartCUAddr()
    if (!bStartInCU) && (uiRPelX < pcSlice.GetSPS().GetPicWidthInLumaSamples()) && (uiBPelY < pcSlice.GetSPS().GetPicHeightInLumaSamples()) {
        this.m_pcEntropyDecoder.DecodeSplitFlag(pcCU, uiAbsPartIdx, uiDepth)
    } else {
        bBoundary = true
    }

    if ((uiDepth < uint(pcCU.GetDepth1(uiAbsPartIdx))) && (uiDepth < pcCU.GetSlice().GetSPS().GetMaxCUDepth()-pcCU.GetSlice().GetSPS().GetAddCUDepth())) || bBoundary {
        uiIdx := uiAbsPartIdx
        if (pcCU.GetSlice().GetSPS().GetMaxCUWidth()>>uiDepth) == pcCU.GetSlice().GetPPS().GetMinCuDQPSize() && pcCU.GetSlice().GetPPS().GetUseDQP() {
            this.SetdQPFlag(true)
            pcCU.SetQPSubParts(int(pcCU.GetRefQP(uiAbsPartIdx)), uiAbsPartIdx, uiDepth) // set QP to default QP
        }

        for uiPartUnitIdx := uint(0); uiPartUnitIdx < 4; uiPartUnitIdx++ {
            uiLPelX = pcCU.GetCUPelX() + TLibCommon.G_auiRasterToPelX[TLibCommon.G_auiZscanToRaster[uiIdx]]
            uiTPelY = pcCU.GetCUPelY() + TLibCommon.G_auiRasterToPelY[TLibCommon.G_auiZscanToRaster[uiIdx]]

            bSubInSlice := pcCU.GetSCUAddr()+uiIdx+uiQNumParts > pcSlice.GetSliceSegmentCurStartCUAddr()
            //fmt.Printf("pcCU.GetSCUAddr()%d+uiIdx%d+uiQNumParts%d>pcSlice.GetDependentSliceCurStartCUAddr()%d\n",pcCU.GetSCUAddr(),uiIdx,uiQNumParts,pcSlice.GetDependentSliceCurStartCUAddr());
            if bSubInSlice {
                if *ruiIsLast == 0 && (uiLPelX < pcCU.GetSlice().GetSPS().GetPicWidthInLumaSamples()) && (uiTPelY < pcCU.GetSlice().GetSPS().GetPicHeightInLumaSamples()) {
                    this.xDecodeCU(pcCU, uiIdx, uiDepth+1, ruiIsLast)
                } else {
                    pcCU.SetOutsideCUPart(uiIdx, uiDepth+1)
                }
            }

            uiIdx += uiQNumParts
        }
        if (pcCU.GetSlice().GetSPS().GetMaxCUWidth()>>uiDepth) == pcCU.GetSlice().GetPPS().GetMinCuDQPSize() && pcCU.GetSlice().GetPPS().GetUseDQP() {
            if this.GetdQPFlag() {
                var uiQPSrcPartIdx uint
                if pcPic.GetCU(pcCU.GetAddr()).GetSliceSegmentStartCU(uiAbsPartIdx) != pcSlice.GetSliceSegmentCurStartCUAddr() {
                    uiQPSrcPartIdx = pcSlice.GetSliceSegmentCurStartCUAddr() % pcPic.GetNumPartInCU()
                } else {
                    uiQPSrcPartIdx = uiAbsPartIdx
                }
                pcCU.SetQPSubParts(int(pcCU.GetRefQP(uiQPSrcPartIdx)), uiAbsPartIdx, uiDepth) // set QP to default QP
            }
        }
        return
    }

    if (pcCU.GetSlice().GetSPS().GetMaxCUWidth()>>uiDepth) >= pcCU.GetSlice().GetPPS().GetMinCuDQPSize() && pcCU.GetSlice().GetPPS().GetUseDQP() {
        this.SetdQPFlag(true)
        pcCU.SetQPSubParts(int(pcCU.GetRefQP(uiAbsPartIdx)), uiAbsPartIdx, uiDepth) // set QP to default QP
    }

    if pcCU.GetSlice().GetPPS().GetTransquantBypassEnableFlag() {
        this.m_pcEntropyDecoder.DecodeCUTransquantBypassFlag(pcCU, uiAbsPartIdx, uiDepth)
    }

    if !pcCU.GetSlice().IsIntra() {
        this.m_pcEntropyDecoder.DecodeSkipFlag(pcCU, uiAbsPartIdx, uiDepth)
    }

    if pcCU.IsSkipped(uiAbsPartIdx) {
        this.m_ppcCU[uiDepth].CopyInterPredInfoFrom(pcCU, uiAbsPartIdx, TLibCommon.REF_PIC_LIST_0)
        this.m_ppcCU[uiDepth].CopyInterPredInfoFrom(pcCU, uiAbsPartIdx, TLibCommon.REF_PIC_LIST_1)
        var cMvFieldNeighbours [TLibCommon.MRG_MAX_NUM_CANDS << 1]TLibCommon.TComMvField // double length for mv of both lists
        var uhInterDirNeighbours [TLibCommon.MRG_MAX_NUM_CANDS]byte
        numValidMergeCand := 0
        for ui := uint(0); ui < this.m_ppcCU[uiDepth].GetSlice().GetMaxNumMergeCand(); ui++ {
            uhInterDirNeighbours[ui] = 0
        }
        this.m_pcEntropyDecoder.DecodeMergeIndex(pcCU, 0, uiAbsPartIdx, uiDepth)
        uiMergeIndex := pcCU.GetMergeIndex1(uiAbsPartIdx)
        this.m_ppcCU[uiDepth].GetInterMergeCandidates(0, 0, cMvFieldNeighbours[:], uhInterDirNeighbours[:], &numValidMergeCand, int(uiMergeIndex))
        pcCU.SetInterDirSubParts(uint(uhInterDirNeighbours[uiMergeIndex]), uiAbsPartIdx, 0, uiDepth)

        cTmpMv := TLibCommon.NewTComMv(0, 0)
        for uiRefListIdx := 0; uiRefListIdx < 2; uiRefListIdx++ {
            if pcCU.GetSlice().GetNumRefIdx(TLibCommon.RefPicList(uiRefListIdx)) > 0 {
                pcCU.SetMVPIdxSubParts(0, TLibCommon.RefPicList(uiRefListIdx), uiAbsPartIdx, 0, uiDepth)
                pcCU.SetMVPNumSubParts(0, TLibCommon.RefPicList(uiRefListIdx), uiAbsPartIdx, 0, uiDepth)
                pcCU.GetCUMvField(TLibCommon.RefPicList(uiRefListIdx)).SetAllMvd(*cTmpMv, TLibCommon.SIZE_2Nx2N, int(uiAbsPartIdx), uiDepth, 0)
                pcCU.GetCUMvField(TLibCommon.RefPicList(uiRefListIdx)).SetAllMvField(&cMvFieldNeighbours[2*int(uiMergeIndex)+uiRefListIdx], TLibCommon.SIZE_2Nx2N, int(uiAbsPartIdx), uiDepth, 0)
            }
        }
        this.xFinishDecodeCU(pcCU, uiAbsPartIdx, uiDepth, ruiIsLast)
        return
    }

    this.m_pcEntropyDecoder.DecodePredMode(pcCU, uiAbsPartIdx, uiDepth)
    this.m_pcEntropyDecoder.DecodePartSize(pcCU, uiAbsPartIdx, uiDepth)

    if pcCU.IsIntra(uiAbsPartIdx) && pcCU.GetPartitionSize1(uiAbsPartIdx) == TLibCommon.SIZE_2Nx2N {
        this.m_pcEntropyDecoder.DecodeIPCMInfo(pcCU, uiAbsPartIdx, uiDepth)

        if pcCU.GetIPCMFlag1(uiAbsPartIdx) {
            this.xFinishDecodeCU(pcCU, uiAbsPartIdx, uiDepth, ruiIsLast)
            return
        }
    }

    uiCurrWidth := pcCU.GetWidth1(uiAbsPartIdx)
    uiCurrHeight := pcCU.GetHeight1(uiAbsPartIdx)

    // prediction mode ( Intra : direction mode, Inter : Mv, reference idx )
    this.m_pcEntropyDecoder.DecodePredInfo(pcCU, uiAbsPartIdx, uiDepth, this.m_ppcCU[uiDepth])

    // Coefficient decoding
    bCodeDQP := this.GetdQPFlag()
    this.m_pcEntropyDecoder.DecodeCoeff(pcCU, uiAbsPartIdx, uiDepth, uint(uiCurrWidth), uint(uiCurrHeight), &bCodeDQP)
    this.SetdQPFlag(bCodeDQP)
    this.xFinishDecodeCU(pcCU, uiAbsPartIdx, uiDepth, ruiIsLast)
}

func (this *TDecCu) xFinishDecodeCU(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint, ruiIsLast *uint) {
    if pcCU.GetSlice().GetPPS().GetUseDQP() {
        if this.GetdQPFlag() {
            pcCU.SetQPSubParts(int(pcCU.GetRefQP(uiAbsPartIdx)), uiAbsPartIdx, uiDepth) // set QP
        } else {
            pcCU.SetQPSubParts(int(pcCU.GetCodedQP()), uiAbsPartIdx, uiDepth) // set QP
        }
    }

    *ruiIsLast = uint(TLibCommon.B2U(this.xDecodeSliceEnd(pcCU, uiAbsPartIdx, uiDepth)))
}

func (this *TDecCu) xDecodeSliceEnd(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint) bool {
    var uiIsLast uint
    pcPic := pcCU.GetPic()
    pcSlice := pcPic.GetSlice(pcPic.GetCurrSliceIdx())
    uiCurNumParts := pcPic.GetNumPartInCU() >> (uiDepth << 1)
    uiWidth := pcSlice.GetSPS().GetPicWidthInLumaSamples()
    uiHeight := pcSlice.GetSPS().GetPicHeightInLumaSamples()
    uiGranularityWidth := pcCU.GetSlice().GetSPS().GetMaxCUWidth()
    uiPosX := pcCU.GetCUPelX() + TLibCommon.G_auiRasterToPelX[TLibCommon.G_auiZscanToRaster[uiAbsPartIdx]]
    uiPosY := pcCU.GetCUPelY() + TLibCommon.G_auiRasterToPelY[TLibCommon.G_auiZscanToRaster[uiAbsPartIdx]]

    if ((uiPosX+uint(pcCU.GetWidth1(uiAbsPartIdx)))%uiGranularityWidth == 0 || (uiPosX+uint(pcCU.GetWidth1(uiAbsPartIdx)) == uiWidth)) &&
        ((uiPosY+uint(pcCU.GetHeight1(uiAbsPartIdx)))%uiGranularityWidth == 0 || (uiPosY+uint(pcCU.GetHeight1(uiAbsPartIdx)) == uiHeight)) {
        this.m_pcEntropyDecoder.DecodeTerminatingBit(&uiIsLast)
    } else {
        uiIsLast = 0
    }

    if uiIsLast != 0 {
        if pcSlice.IsNextSliceSegment() && !pcSlice.IsNextSlice() {
            pcSlice.SetSliceSegmentCurEndCUAddr(pcCU.GetSCUAddr() + uiAbsPartIdx + uiCurNumParts)
        } else {
            pcSlice.SetSliceCurEndCUAddr(pcCU.GetSCUAddr() + uiAbsPartIdx + uiCurNumParts)
            pcSlice.SetSliceSegmentCurEndCUAddr(pcCU.GetSCUAddr() + uiAbsPartIdx + uiCurNumParts)
        }
    }

    return uiIsLast > 0
}
func (this *TDecCu) xDecompressCU(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint) {
    pcPic := pcCU.GetPic()

    bBoundary := false
    uiLPelX := pcCU.GetCUPelX() + TLibCommon.G_auiRasterToPelX[TLibCommon.G_auiZscanToRaster[uiAbsPartIdx]]
    uiRPelX := uiLPelX + (pcCU.GetSlice().GetSPS().GetMaxCUWidth() >> uiDepth) - 1
    uiTPelY := pcCU.GetCUPelY() + TLibCommon.G_auiRasterToPelY[TLibCommon.G_auiZscanToRaster[uiAbsPartIdx]]
    uiBPelY := uiTPelY + (pcCU.GetSlice().GetSPS().GetMaxCUHeight() >> uiDepth) - 1

    uiCurNumParts := pcPic.GetNumPartInCU() >> (uiDepth << 1)
    pcSlice := pcCU.GetPic().GetSlice(pcCU.GetPic().GetCurrSliceIdx())
    bStartInCU := pcCU.GetSCUAddr()+uiAbsPartIdx+uiCurNumParts > pcSlice.GetSliceSegmentCurStartCUAddr() && pcCU.GetSCUAddr()+uiAbsPartIdx < pcSlice.GetSliceSegmentCurStartCUAddr()
    if bStartInCU || (uiRPelX >= pcSlice.GetSPS().GetPicWidthInLumaSamples()) || (uiBPelY >= pcSlice.GetSPS().GetPicHeightInLumaSamples()) {
        bBoundary = true
    }

    if ((uiDepth < uint(pcCU.GetDepth1(uiAbsPartIdx))) && (uiDepth < pcCU.GetSlice().GetSPS().GetMaxCUDepth()-pcCU.GetSlice().GetSPS().GetAddCUDepth())) || bBoundary {
        uiNextDepth := uiDepth + 1
        uiQNumParts := pcCU.GetTotalNumPart() >> (uiNextDepth << 1)
        uiIdx := uiAbsPartIdx
        for uiPartIdx := 0; uiPartIdx < 4; uiPartIdx++ {
            uiLPelX = pcCU.GetCUPelX() + TLibCommon.G_auiRasterToPelX[TLibCommon.G_auiZscanToRaster[uiIdx]]
            uiTPelY = pcCU.GetCUPelY() + TLibCommon.G_auiRasterToPelY[TLibCommon.G_auiZscanToRaster[uiIdx]]

            binSlice := (pcCU.GetSCUAddr()+uiIdx+uiQNumParts > pcSlice.GetSliceSegmentCurStartCUAddr()) && (pcCU.GetSCUAddr()+uiIdx < pcSlice.GetSliceSegmentCurEndCUAddr())
            if binSlice && (uiLPelX < pcSlice.GetSPS().GetPicWidthInLumaSamples()) && (uiTPelY < pcSlice.GetSPS().GetPicHeightInLumaSamples()) {
                this.xDecompressCU(pcCU, uiIdx, uiNextDepth)
            }

            uiIdx += uiQNumParts
        }
        return
    }

    // Residual reconstruction
    this.m_ppcYuvResi[uiDepth].Clear()

    this.m_ppcCU[uiDepth].CopySubCU(pcCU, uiAbsPartIdx, uiDepth)

    /*#ifdef ENC_DEC_TRACE*/
    this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XTraceCUHeader(TLibCommon.TRACE_CU)

    this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(int(this.m_ppcCU[uiDepth].GetCUPelX()), "cu_x", TLibCommon.TRACE_CU)
    this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(int(this.m_ppcCU[uiDepth].GetCUPelY()), "cu_y", TLibCommon.TRACE_CU)
    this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(int(this.m_ppcCU[uiDepth].GetWidth1(0)), "cu_size", TLibCommon.TRACE_CU)
    this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(int(this.m_ppcCU[uiDepth].GetPredictionMode1(0)), "cu_type", TLibCommon.TRACE_CU)

    if this.m_ppcCU[uiDepth].GetSlice().GetPPS().GetTransquantBypassEnableFlag() {
        this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(int(TLibCommon.B2U(this.m_ppcCU[uiDepth].GetCUTransquantBypass1(0))), "cu_transquant_bypass_flag", TLibCommon.TRACE_CU)
    }

    if this.m_ppcCU[uiDepth].GetPredictionMode1(0) == TLibCommon.MODE_INTRA {
        this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(int(TLibCommon.B2U(this.m_ppcCU[uiDepth].GetIPCMFlag1(0))), "cu_pcm_skip_flag", TLibCommon.TRACE_CU)

        for iPartIdx := byte(0); iPartIdx < this.m_ppcCU[uiDepth].GetNumPartInter(); iPartIdx++ {
            var iWidth, iHeight, iPosX, iPosY int
            var uiPartAddr uint

            this.m_ppcCU[uiDepth].GetPartIndexAndSizePos(uint(iPartIdx), &uiPartAddr, &iWidth, &iHeight, &iPosX, &iPosY)

            this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XTracePUHeader(TLibCommon.TRACE_PU)

            this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(iPosX, "pu_x", TLibCommon.TRACE_PU)
            this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(iPosY, "pu_y", TLibCommon.TRACE_PU)
            this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(iWidth, "pu_width", TLibCommon.TRACE_PU)
            this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(iHeight, "pu_height", TLibCommon.TRACE_PU)
            this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(int(this.m_ppcCU[uiDepth].GetPartitionSize1(0)), "pu_shape", TLibCommon.TRACE_PU)
            this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(int(this.m_ppcCU[uiDepth].GetLumaIntraDir1(uiPartAddr)), "pu_intra_pred_mode_luma", TLibCommon.TRACE_PU)
            this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(int(this.m_ppcCU[uiDepth].GetChromaIntraDir1(0)), "pu_intra_pred_mode_chroma", TLibCommon.TRACE_PU)
        }
    } else {
        this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(int(TLibCommon.B2U(this.m_ppcCU[uiDepth].GetSkipFlag1(0))), "cu_pcm_skip_flag", TLibCommon.TRACE_CU)

        for iPartIdx := byte(0); iPartIdx < this.m_ppcCU[uiDepth].GetNumPartInter(); iPartIdx++ {
            var iWidth, iHeight, iPosX, iPosY int
            var uiPartAddr uint

            this.m_ppcCU[uiDepth].GetPartIndexAndSizePos(uint(iPartIdx), &uiPartAddr, &iWidth, &iHeight, &iPosX, &iPosY)

            this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XTracePUHeader(TLibCommon.TRACE_PU)

            this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(iPosX, "pu_x", TLibCommon.TRACE_PU)
            this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(iPosY, "pu_y", TLibCommon.TRACE_PU)
            this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(iWidth, "pu_width", TLibCommon.TRACE_PU)
            this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(iHeight, "pu_height", TLibCommon.TRACE_PU)
            this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(int(this.m_ppcCU[uiDepth].GetPartitionSize1(0)), "pu_shape", TLibCommon.TRACE_PU)

            //fmt.Printf("(%d,%d)\n", this.m_ppcCU[uiDepth].GetCUMvField(TLibCommon.REF_PIC_LIST_0).GetRefIdx(int(uiPartAddr)), this.m_ppcCU[uiDepth].GetCUMvField(TLibCommon.REF_PIC_LIST_1).GetRefIdx(int(uiPartAddr)));
            if this.m_ppcCU[uiDepth].GetCUMvField(TLibCommon.REF_PIC_LIST_0).GetRefIdx(int(uiPartAddr)) >= 0 &&
                this.m_ppcCU[uiDepth].GetCUMvField(TLibCommon.REF_PIC_LIST_1).GetRefIdx(int(uiPartAddr)) >= 0 {
                this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(2, "pu_inter_pred_mode", TLibCommon.TRACE_PU)

                refIdx := this.m_ppcCU[uiDepth].GetCUMvField(TLibCommon.REF_PIC_LIST_0).GetRefIdx(int(uiPartAddr))
                cMv := this.m_ppcCU[uiDepth].GetCUMvField(TLibCommon.REF_PIC_LIST_0).GetMv(int(uiPartAddr))
                this.m_ppcCU[uiDepth].ClipMv(&cMv)
                this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(int(refIdx), "pu_ref_id", TLibCommon.TRACE_PU)
                this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(int(cMv.GetHor()), "pu_mv_x", TLibCommon.TRACE_PU)
                this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(int(cMv.GetVer()), "pu_mv_y", TLibCommon.TRACE_PU)

                refIdx = this.m_ppcCU[uiDepth].GetCUMvField(TLibCommon.REF_PIC_LIST_1).GetRefIdx(int(uiPartAddr))
                cMv = this.m_ppcCU[uiDepth].GetCUMvField(TLibCommon.REF_PIC_LIST_1).GetMv(int(uiPartAddr))
                this.m_ppcCU[uiDepth].ClipMv(&cMv)
                this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(int(refIdx), "pu_ref_id", TLibCommon.TRACE_PU)
                this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(int(cMv.GetHor()), "pu_mv_x", TLibCommon.TRACE_PU)
                this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(int(cMv.GetVer()), "pu_mv_y", TLibCommon.TRACE_PU)
            } else if this.m_ppcCU[uiDepth].GetCUMvField(TLibCommon.REF_PIC_LIST_0).GetRefIdx(int(uiPartAddr)) >= 0 {
                this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(0, "pu_inter_pred_mode", TLibCommon.TRACE_PU)

                refIdx := this.m_ppcCU[uiDepth].GetCUMvField(TLibCommon.REF_PIC_LIST_0).GetRefIdx(int(uiPartAddr))
                cMv := this.m_ppcCU[uiDepth].GetCUMvField(TLibCommon.REF_PIC_LIST_0).GetMv(int(uiPartAddr))
                this.m_ppcCU[uiDepth].ClipMv(&cMv)
                this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(int(refIdx), "pu_ref_id", TLibCommon.TRACE_PU)
                this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(int(cMv.GetHor()), "pu_mv_x", TLibCommon.TRACE_PU)
                this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(int(cMv.GetVer()), "pu_mv_y", TLibCommon.TRACE_PU)
            } else if this.m_ppcCU[uiDepth].GetCUMvField(TLibCommon.REF_PIC_LIST_1).GetRefIdx(int(uiPartAddr)) >= 0 {
                this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(1, "pu_inter_pred_mode", TLibCommon.TRACE_PU)

                refIdx := this.m_ppcCU[uiDepth].GetCUMvField(TLibCommon.REF_PIC_LIST_1).GetRefIdx(int(uiPartAddr))
                cMv := this.m_ppcCU[uiDepth].GetCUMvField(TLibCommon.REF_PIC_LIST_1).GetMv(int(uiPartAddr))
                this.m_ppcCU[uiDepth].ClipMv(&cMv)
                this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(int(refIdx), "pu_ref_id", TLibCommon.TRACE_PU)
                this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(int(cMv.GetHor()), "pu_mv_x", TLibCommon.TRACE_PU)
                this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(int(cMv.GetVer()), "pu_mv_y", TLibCommon.TRACE_PU)
            } else {
                //assert(0);
                fmt.Printf("pu_inter_pred_mode error\n")
            }
        }
    }
    /*#endif*/

    switch this.m_ppcCU[uiDepth].GetPredictionMode1(0) {
    case TLibCommon.MODE_INTER:
        this.xReconInter(this.m_ppcCU[uiDepth], uiDepth)
        //break;
    case TLibCommon.MODE_INTRA:
        this.xReconIntraQT(this.m_ppcCU[uiDepth], uiDepth)
        //break;
    default:
        //assert(0);
        //break;
    }
    if this.m_ppcCU[uiDepth].IsLosslessCoded(0) && (this.m_ppcCU[uiDepth].GetIPCMFlag1(0) == false) {
        this.xFillPCMBuffer(this.m_ppcCU[uiDepth], uiDepth)
    }

    this.xCopyToPic(this.m_ppcCU[uiDepth], pcPic, uiAbsPartIdx, uiDepth)
}

func (this *TDecCu) xReconInter(pcCU *TLibCommon.TComDataCU, uiDepth uint) {
	//fmt.Printf("xReconInter %d, %d, uiDepth=%d\n", this.m_ppcYuvReco[uiDepth].GetWidth(),this.m_ppcYuvReco[uiDepth].GetHeight(), uiDepth);
    // inter prediction
    this.m_pcPrediction.MotionCompensation(pcCU, this.m_ppcYuvReco[uiDepth], TLibCommon.REF_PIC_LIST_X, -1)

    // inter recon
    this.xDecodeInterTexture(pcCU, this.m_ppcYuvReco[uiDepth], 0, uiDepth)

    // clip for only non-zero cbp case
    if (pcCU.GetCbf2(0, TLibCommon.TEXT_LUMA) != 0) || (pcCU.GetCbf2(0, TLibCommon.TEXT_CHROMA_U) != 0) || (pcCU.GetCbf2(0, TLibCommon.TEXT_CHROMA_V) != 0) {
        this.m_ppcYuvReco[uiDepth].AddClip(this.m_ppcYuvReco[uiDepth], this.m_ppcYuvResi[uiDepth], 0, uint(pcCU.GetWidth1(0)))
    } else {
        this.m_ppcYuvReco[uiDepth].CopyPartToPartYuv(this.m_ppcYuvReco[uiDepth], 0, uint(pcCU.GetWidth1(0)), uint(pcCU.GetHeight1(0)))
    }
}

func (this *TDecCu) xReconIntraQT(pcCU *TLibCommon.TComDataCU, uiDepth uint) {
    var uiInitTrDepth uint
    if pcCU.GetPartitionSize1(0) == TLibCommon.SIZE_2Nx2N {
        uiInitTrDepth = 0
    } else {
        uiInitTrDepth = 1
    }
    uiNumPart := uint(pcCU.GetNumPartInter())
    uiNumQParts := pcCU.GetTotalNumPart() >> 2

    if pcCU.GetIPCMFlag1(0) {
        this.xReconPCM(pcCU, uiDepth)
        return
    }

    for uiPU := uint(0); uiPU < uiNumPart; uiPU++ {
        this.xIntraLumaRecQT(pcCU, uiInitTrDepth, uiPU*uiNumQParts, this.m_ppcYuvReco[uiDepth], this.m_ppcYuvReco[uiDepth], this.m_ppcYuvResi[uiDepth])
    }

    for uiPU := uint(0); uiPU < uiNumPart; uiPU++ {
        this.xIntraChromaRecQT(pcCU, uiInitTrDepth, uiPU*uiNumQParts, this.m_ppcYuvReco[uiDepth], this.m_ppcYuvReco[uiDepth], this.m_ppcYuvResi[uiDepth])
    }
}
func (this *TDecCu) xIntraRecLumaBlk(pcCU *TLibCommon.TComDataCU, uiTrDepth, uiAbsPartIdx uint, pcRecoYuv *TLibCommon.TComYuv, pcPredYuv *TLibCommon.TComYuv, pcResiYuv *TLibCommon.TComYuv) {
    uiWidth := uint(pcCU.GetWidth1(0)) >> uiTrDepth
    uiHeight := uint(pcCU.GetHeight1(0)) >> uiTrDepth
    uiStride := pcRecoYuv.GetStride()
    piReco := pcRecoYuv.GetLumaAddr1(uiAbsPartIdx)
    piPred := pcPredYuv.GetLumaAddr1(uiAbsPartIdx)
    piResi := pcResiYuv.GetLumaAddr1(uiAbsPartIdx)

    uiNumCoeffInc := (pcCU.GetSlice().GetSPS().GetMaxCUWidth() * pcCU.GetSlice().GetSPS().GetMaxCUHeight()) >> (pcCU.GetSlice().GetSPS().GetMaxCUDepth() << 1)
    pcCoeff := pcCU.GetCoeffY()[(uiNumCoeffInc * uiAbsPartIdx):]

    uiLumaPredMode := pcCU.GetLumaIntraDir1(uiAbsPartIdx)

    /*#ifdef ENC_DEC_TRACE*/
    blkX := int(TLibCommon.G_auiRasterToPelX[TLibCommon.G_auiZscanToRaster[uiAbsPartIdx]])
    blkY := int(TLibCommon.G_auiRasterToPelY[TLibCommon.G_auiZscanToRaster[uiAbsPartIdx]])

    this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XTraceTUHeader(TLibCommon.TRACE_TU)

    this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(0, "tu_color", TLibCommon.TRACE_TU)
    this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(blkX, "tu_x", TLibCommon.TRACE_TU)
    this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(blkY, "tu_y", TLibCommon.TRACE_TU)
    this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(int(uiWidth), "tu_width", TLibCommon.TRACE_TU)
    this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(int(uiHeight), "tu_height", TLibCommon.TRACE_TU)
    /*#endif*/

    uiZOrder := pcCU.GetZorderIdxInCU() + uiAbsPartIdx
    piRecIPred := pcCU.GetPic().GetPicYuvRec().GetLumaAddr2(int(pcCU.GetAddr()), int(uiZOrder))
    uiRecIPredStride := uint(pcCU.GetPic().GetPicYuvRec().GetStride())
    useTransformSkip := pcCU.GetTransformSkip2(uiAbsPartIdx, TLibCommon.TEXT_LUMA)
    //===== init availability pattern =====
    bAboveAvail := false
    bLeftAvail := false
    pcCU.GetPattern().InitPattern3(pcCU, uiTrDepth, uiAbsPartIdx)
    pcCU.GetPattern().InitAdiPattern(pcCU, uiAbsPartIdx, uiTrDepth,
        this.m_pcPrediction.GetPredicBuf(),
        this.m_pcPrediction.GetPredicBufWidth(),
        this.m_pcPrediction.GetPredicBufHeight(),
        &bAboveAvail, &bLeftAvail, false)

    //===== get prediction signal =====
    this.m_pcPrediction.PredIntraLumaAng(pcCU.GetPattern(), uint(uiLumaPredMode), piPred, uiStride, int(uiWidth), int(uiHeight), bAboveAvail, bLeftAvail)

    //===== inverse transform =====
    this.m_pcTrQuant.SetQPforQuant(int(pcCU.GetQP1(0)), TLibCommon.TEXT_LUMA, pcCU.GetSlice().GetSPS().GetQpBDOffsetY(), 0)

    var scalingListType int
    if pcCU.IsIntra(uiAbsPartIdx) {
        scalingListType = 0 + TLibCommon.G_eTTable[TLibCommon.TEXT_LUMA]
    } else {
        scalingListType = 3 + TLibCommon.G_eTTable[TLibCommon.TEXT_LUMA]
    }
    //assert(scalingListType < 6);

    /*#ifdef ENC_DEC_TRACE*/
    this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XTraceCoefHeader(TLibCommon.TRACE_COEF)

    for k := uint(0); k < uiHeight; k++ {
        this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadCeofTr(pcCoeff[k*uiWidth:], uiWidth, TLibCommon.TRACE_COEF)
    }
    /*#endif*/
    this.m_pcTrQuant.InvtransformNxN(pcCU.GetCUTransquantBypass1(uiAbsPartIdx), TLibCommon.TEXT_LUMA, uint(pcCU.GetLumaIntraDir1(uiAbsPartIdx)), piResi, uiStride, pcCoeff, uiWidth, uiHeight, scalingListType, useTransformSkip)
    /*#ifdef ENC_DEC_TRACE*/
    this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XTraceResiHeader(TLibCommon.TRACE_RESI)

    for k := uint(0); k < uiHeight; k++ {
        this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadResiTr(piResi[k*uiStride:], uiWidth, TLibCommon.TRACE_RESI)
    }
    /*#endif*/

    /*#ifdef ENC_DEC_TRACE*/
    {
        this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XTracePredHeader(TLibCommon.TRACE_PRED)

        pPred := piPred
        for uiY := uint(0); uiY < uiHeight; uiY++ {
            this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadPredTr(pPred[uiY*uiStride:], uiWidth, TLibCommon.TRACE_PRED)
            //pPred += uiStride;
        }
    }

    this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XTraceRecoHeader(TLibCommon.TRACE_RECON)
    /*#endif*/
    //===== reconstruction =====
    pPred := piPred
    pResi := piResi
    pReco := piReco
    pRecIPred := piRecIPred
    for uiY := uint(0); uiY < uiHeight; uiY++ {
        for uiX := uint(0); uiX < uiWidth; uiX++ {
            pReco[uiY*uiStride+uiX] = TLibCommon.ClipY(pPred[uiY*uiStride+uiX] + pResi[uiY*uiStride+uiX])
            pRecIPred[uiY*uiRecIPredStride+uiX] = pReco[uiY*uiStride+uiX]
        }
        /*#ifdef ENC_DEC_TRACE*/
        this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadRecoTr(pReco[uiY*uiStride:], uiWidth, TLibCommon.TRACE_RECON)
        /*#endif*/
        /*pPred     += uiStride;
          pResi     += uiStride;
          pReco     += uiStride;
          pRecIPred += uiRecIPredStride;*/
    }
}
func (this *TDecCu) xIntraRecChromaBlk(pcCU *TLibCommon.TComDataCU, uiTrDepth, uiAbsPartIdx uint, pcRecoYuv *TLibCommon.TComYuv, pcPredYuv *TLibCommon.TComYuv, pcResiYuv *TLibCommon.TComYuv, uiChromaId uint) {
    uiFullDepth := uint(pcCU.GetDepth1(0)) + uiTrDepth
    uiLog2TrSize := TLibCommon.G_aucConvertToBit[pcCU.GetSlice().GetSPS().GetMaxCUWidth()>>uiFullDepth] + 2

    if uiLog2TrSize == 2 {
        //assert( uiTrDepth > 0 );
        uiTrDepth--
        uiQPDiv := pcCU.GetPic().GetNumPartInCU() >> ((uint(pcCU.GetDepth1(0)) + uiTrDepth) << 1)
        bFirstQ := ((uiAbsPartIdx % uiQPDiv) == 0)
        if !bFirstQ {
            return
        }
    }

    uiWidth := uint(pcCU.GetWidth1(0)) >> (uiTrDepth + 1)
    uiHeight := uint(pcCU.GetHeight1(0)) >> (uiTrDepth + 1)
    uiStride := pcRecoYuv.GetCStride()
    uiNumCoeffInc := ((pcCU.GetSlice().GetSPS().GetMaxCUWidth() * pcCU.GetSlice().GetSPS().GetMaxCUHeight()) >> (pcCU.GetSlice().GetSPS().GetMaxCUDepth() << 1)) >> 2

    var eText TLibCommon.TextType
    var piReco, piPred, piResi []TLibCommon.Pel
    var pcCoeff []TLibCommon.TCoeff
    if uiChromaId > 0 {
        eText = TLibCommon.TEXT_CHROMA_V
        piReco = pcRecoYuv.GetCrAddr1(uiAbsPartIdx)
        piPred = pcPredYuv.GetCrAddr1(uiAbsPartIdx)
        piResi = pcResiYuv.GetCrAddr1(uiAbsPartIdx)
        pcCoeff = pcCU.GetCoeffCr()[(uiNumCoeffInc * uiAbsPartIdx):]
    } else {
        eText = TLibCommon.TEXT_CHROMA_U
        piReco = pcRecoYuv.GetCbAddr1(uiAbsPartIdx)
        piPred = pcPredYuv.GetCbAddr1(uiAbsPartIdx)
        piResi = pcResiYuv.GetCbAddr1(uiAbsPartIdx)
        pcCoeff = pcCU.GetCoeffCb()[(uiNumCoeffInc * uiAbsPartIdx):]
    }
    uiChromaPredMode := pcCU.GetChromaIntraDir1(0)

    /*#ifdef ENC_DEC_TRACE*/
    //if(uiChromaId==1){
    blkX := int(TLibCommon.G_auiRasterToPelX[TLibCommon.G_auiZscanToRaster[uiAbsPartIdx]]) >> 1
    blkY := int(TLibCommon.G_auiRasterToPelY[TLibCommon.G_auiZscanToRaster[uiAbsPartIdx]]) >> 1

    this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XTraceTUHeader(TLibCommon.TRACE_TU)

    this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(int(uiChromaId)+1, "tu_color", TLibCommon.TRACE_TU)
    this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(blkX, "tu_x", TLibCommon.TRACE_TU)
    this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(blkY, "tu_y", TLibCommon.TRACE_TU)
    this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(int(uiWidth), "tu_width", TLibCommon.TRACE_TU)
    this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(int(uiHeight), "tu_height", TLibCommon.TRACE_TU)
    //}
    /*#endif*/

    uiZOrder := pcCU.GetZorderIdxInCU() + uiAbsPartIdx
    var piRecIPred []TLibCommon.Pel
    if uiChromaId > 0 {
        piRecIPred = pcCU.GetPic().GetPicYuvRec().GetCrAddr2(int(pcCU.GetAddr()), int(uiZOrder))
    } else {
        piRecIPred = pcCU.GetPic().GetPicYuvRec().GetCbAddr2(int(pcCU.GetAddr()), int(uiZOrder))
    }
    uiRecIPredStride := uint(pcCU.GetPic().GetPicYuvRec().GetCStride())
    useTransformSkipChroma := pcCU.GetTransformSkip2(uiAbsPartIdx, eText)
    //===== init availability pattern =====
    bAboveAvail := false
    bLeftAvail := false
    pcCU.GetPattern().InitPattern3(pcCU, uiTrDepth, uiAbsPartIdx)

    pcCU.GetPattern().InitAdiPatternChroma(pcCU, uiAbsPartIdx, uiTrDepth,
        this.m_pcPrediction.GetPredicBuf(),
        this.m_pcPrediction.GetPredicBufWidth(),
        this.m_pcPrediction.GetPredicBufHeight(),
        &bAboveAvail, &bLeftAvail, uiChromaId)
    var pPatChroma []TLibCommon.Pel
    if uiChromaId > 0 {
        pPatChroma = pcCU.GetPattern().GetAdiCrBuf(int(uiWidth), int(uiHeight), this.m_pcPrediction.GetPredicBuf())
    } else {
        pPatChroma = pcCU.GetPattern().GetAdiCbBuf(int(uiWidth), int(uiHeight), this.m_pcPrediction.GetPredicBuf())
    }
    //===== get prediction signal =====
    {
        if uiChromaPredMode == TLibCommon.DM_CHROMA_IDX {
            uiChromaPredMode = pcCU.GetLumaIntraDir1(0)
        }
        this.m_pcPrediction.PredIntraChromaAng(pPatChroma, uint(uiChromaPredMode), piPred, uiStride, int(uiWidth), int(uiHeight), bAboveAvail, bLeftAvail)
    }

    //===== inverse transform =====
    var curChromaQpOffset int
    if eText == TLibCommon.TEXT_CHROMA_U {
        curChromaQpOffset = pcCU.GetSlice().GetPPS().GetChromaCbQpOffset() + pcCU.GetSlice().GetSliceQpDeltaCb()
    } else {
        curChromaQpOffset = pcCU.GetSlice().GetPPS().GetChromaCrQpOffset() + pcCU.GetSlice().GetSliceQpDeltaCr()
    }
    this.m_pcTrQuant.SetQPforQuant(int(pcCU.GetQP1(0)), eText, pcCU.GetSlice().GetSPS().GetQpBDOffsetC(), curChromaQpOffset)

    var scalingListType int
    if pcCU.IsIntra(uiAbsPartIdx) {
        scalingListType = 0 + TLibCommon.G_eTTable[eText]
    } else {
        scalingListType = 3 + TLibCommon.G_eTTable[eText]
    }
    //assert(scalingListType < 6);
    /*#ifdef ENC_DEC_TRACE*/
    this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XTraceCoefHeader(TLibCommon.TRACE_COEF)

    for k := uint(0); k < uiHeight; k++ {
        this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadCeofTr(pcCoeff[k*uiWidth:], uiWidth, TLibCommon.TRACE_COEF)
    }
    /*#endif*/
    this.m_pcTrQuant.InvtransformNxN(pcCU.GetCUTransquantBypass1(uiAbsPartIdx), eText, TLibCommon.REG_DCT, piResi, uiStride, pcCoeff, uiWidth, uiHeight, scalingListType, useTransformSkipChroma)
    /*#ifdef ENC_DEC_TRACE*/
    this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XTraceResiHeader(TLibCommon.TRACE_RESI)

    for k := uint(0); k < uiHeight; k++ {
        this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadResiTr(piResi[k*uiStride:], uiWidth, TLibCommon.TRACE_RESI)
    }
    /*#endif*/
    /*#ifdef ENC_DEC_TRACE*/
    {
        this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XTracePredHeader(TLibCommon.TRACE_PRED)

        pPred := piPred
        for uiY := uint(0); uiY < uiHeight; uiY++ {
            this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadPredTr(pPred[uiY*uiStride:], uiWidth, TLibCommon.TRACE_PRED)
            //pPred += uiStride;
        }
    }

    this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XTraceRecoHeader(TLibCommon.TRACE_RECON)
    /*#endif*/
    //===== reconstruction =====
    pPred := piPred
    pResi := piResi
    pReco := piReco
    pRecIPred := piRecIPred
    for uiY := uint(0); uiY < uiHeight; uiY++ {
        for uiX := uint(0); uiX < uiWidth; uiX++ {
            pReco[uiY*uiStride+uiX] = TLibCommon.ClipC(pPred[uiY*uiStride+uiX] + pResi[uiY*uiStride+uiX])
            pRecIPred[uiY*uiRecIPredStride+uiX] = pReco[uiY*uiStride+uiX]
        }
        /*#ifdef ENC_DEC_TRACE*/
        this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadRecoTr(pReco[uiY*uiStride:], uiWidth, TLibCommon.TRACE_RECON)
        /*#endif*/
        /*pPred     += uiStride;
          pResi     += uiStride;
          pReco     += uiStride;
          pRecIPred += uiRecIPredStride;*/
    }
}

func (this *TDecCu) xReconPCM(pcCU *TLibCommon.TComDataCU, uiDepth uint) {
    // Luma
    uiWidth := (pcCU.GetSlice().GetSPS().GetMaxCUWidth() >> uiDepth)
    uiHeight := (pcCU.GetSlice().GetSPS().GetMaxCUHeight() >> uiDepth)

    piPcmY := pcCU.GetPCMSampleY()
    piRecoY := this.m_ppcYuvReco[uiDepth].GetLumaAddr2(0, uiWidth)

    uiStride := this.m_ppcYuvResi[uiDepth].GetStride()

    /*#ifdef ENC_DEC_TRACE*/
    blkX := 0 //TLibCommon.G_auiRasterToPelX[ TLibCommon.G_auiZscanToRaster[ uiAbsPartIdx ] ];
    blkY := 0 //TLibCommon.G_auiRasterToPelY[ TLibCommon.G_auiZscanToRaster[ uiAbsPartIdx ] ];

    this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XTraceTUHeader(TLibCommon.TRACE_TU)

    this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(0, "tu_color", TLibCommon.TRACE_TU)
    this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(blkX, "tu_x", TLibCommon.TRACE_TU)
    this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(blkY, "tu_y", TLibCommon.TRACE_TU)
    this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(int(uiWidth), "tu_width", TLibCommon.TRACE_TU)
    this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(int(uiHeight), "tu_height", TLibCommon.TRACE_TU)
    /*#endif*/

    this.xDecodePCMTexture(pcCU, 0, piPcmY, piRecoY, uiStride, uiWidth, uiHeight, TLibCommon.TEXT_LUMA)

    // Cb and Cr
    uiCWidth := (uiWidth >> 1)
    uiCHeight := (uiHeight >> 1)

    piPcmCb := pcCU.GetPCMSampleCb()
    piPcmCr := pcCU.GetPCMSampleCr()
    pRecoCb := this.m_ppcYuvReco[uiDepth].GetCbAddr()
    pRecoCr := this.m_ppcYuvReco[uiDepth].GetCrAddr()

    uiCStride := this.m_ppcYuvReco[uiDepth].GetCStride()

    /*#ifdef ENC_DEC_TRACE*/
    CblkX := 0 //TLibCommon.G_auiRasterToPelX[ TLibCommon.G_auiZscanToRaster[ uiAbsPartIdx ] ]>>1;
    CblkY := 0 //TLibCommon.G_auiRasterToPelY[ TLibCommon.G_auiZscanToRaster[ uiAbsPartIdx ] ]>>1;

    this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XTraceTUHeader(TLibCommon.TRACE_TU)

    this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(1, "tu_color", TLibCommon.TRACE_TU)
    this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(CblkX, "tu_x", TLibCommon.TRACE_TU)
    this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(CblkY, "tu_y", TLibCommon.TRACE_TU)
    this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(int(uiCWidth), "tu_width", TLibCommon.TRACE_TU)
    this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(int(uiCHeight), "tu_height", TLibCommon.TRACE_TU)
    /*#endif*/

    this.xDecodePCMTexture(pcCU, 0, piPcmCb, pRecoCb, uiCStride, uiCWidth, uiCHeight, TLibCommon.TEXT_CHROMA_U)

    /*#ifdef ENC_DEC_TRACE*/
    this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XTraceTUHeader(TLibCommon.TRACE_TU)

    this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(2, "tu_color", TLibCommon.TRACE_TU)
    this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(CblkX, "tu_x", TLibCommon.TRACE_TU)
    this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(CblkY, "tu_y", TLibCommon.TRACE_TU)
    this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(int(uiCWidth), "tu_width", TLibCommon.TRACE_TU)
    this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(int(uiCHeight), "tu_height", TLibCommon.TRACE_TU)
    /*#endif*/
    this.xDecodePCMTexture(pcCU, 0, piPcmCr, pRecoCr, uiCStride, uiCWidth, uiCHeight, TLibCommon.TEXT_CHROMA_V)
}

func (this *TDecCu) xDecodeInterTexture(pcCU *TLibCommon.TComDataCU, pcYuvPred *TLibCommon.TComYuv, uiAbsPartIdx, uiDepth uint) {
    uiWidth := uint(pcCU.GetWidth1(uiAbsPartIdx))
    uiHeight := uint(pcCU.GetHeight1(uiAbsPartIdx))
    var piCoeff []TLibCommon.TCoeff

    var pResi []TLibCommon.Pel
    trMode := uint(pcCU.GetTransformIdx1(uiAbsPartIdx))

    // Y
    piCoeff = pcCU.GetCoeffY()
    pResi = this.m_ppcYuvResi[uiDepth].GetLumaAddr()

    this.m_pcTrQuant.SetQPforQuant(int(pcCU.GetQP1(uiAbsPartIdx)), TLibCommon.TEXT_LUMA, pcCU.GetSlice().GetSPS().GetQpBDOffsetY(), 0)

    this.InvRecurTransformNxN(pcCU, pcYuvPred, 0, TLibCommon.TEXT_LUMA, pResi, 0, this.m_ppcYuvResi[uiDepth].GetStride(), uiWidth, uiHeight, trMode, 0, piCoeff)

    // Cb and Cr
    curChromaQpOffset := pcCU.GetSlice().GetPPS().GetChromaCbQpOffset() + pcCU.GetSlice().GetSliceQpDeltaCb()
    this.m_pcTrQuant.SetQPforQuant(int(pcCU.GetQP1(uiAbsPartIdx)), TLibCommon.TEXT_CHROMA, pcCU.GetSlice().GetSPS().GetQpBDOffsetC(), curChromaQpOffset)

    uiWidth >>= 1
    uiHeight >>= 1
    piCoeff = pcCU.GetCoeffCb()
    pResi = this.m_ppcYuvResi[uiDepth].GetCbAddr()
    this.InvRecurTransformNxN(pcCU, pcYuvPred, 0, TLibCommon.TEXT_CHROMA_U, pResi, 0, this.m_ppcYuvResi[uiDepth].GetCStride(), uiWidth, uiHeight, trMode, 0, piCoeff)

    curChromaQpOffset = pcCU.GetSlice().GetPPS().GetChromaCrQpOffset() + pcCU.GetSlice().GetSliceQpDeltaCr()
    this.m_pcTrQuant.SetQPforQuant(int(pcCU.GetQP1(uiAbsPartIdx)), TLibCommon.TEXT_CHROMA, pcCU.GetSlice().GetSPS().GetQpBDOffsetC(), curChromaQpOffset)

    piCoeff = pcCU.GetCoeffCr()
    pResi = this.m_ppcYuvResi[uiDepth].GetCrAddr()
    this.InvRecurTransformNxN(pcCU, pcYuvPred, 0, TLibCommon.TEXT_CHROMA_V, pResi, 0, this.m_ppcYuvResi[uiDepth].GetCStride(), uiWidth, uiHeight, trMode, 0, piCoeff)
}
func (this *TDecCu) xDecodePCMTexture(pcCU *TLibCommon.TComDataCU, uiPartIdx uint, piPCM, piReco []TLibCommon.Pel, uiStride, uiWidth, uiHeight uint, ttText TLibCommon.TextType) {
    var uiX, uiY uint
    var piPicReco []TLibCommon.Pel
    var uiPicStride, uiPcmLeftShiftBit uint

    /* #ifdef ENC_DEC_TRACE*/
    this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XTraceCoefHeader(TLibCommon.TRACE_COEF)

    var coeffs [64]TLibCommon.TCoeff
    for uiY = 0; uiY < uiHeight; uiY++ {
        for uiX = 0; uiX < uiWidth; uiX++ {
            coeffs[uiX] = TLibCommon.TCoeff(piPCM[uiY*uiWidth+uiX])
        }
        this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadCeofTr(coeffs[:], uiWidth, TLibCommon.TRACE_COEF)
    }

    this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XTraceResiHeader(TLibCommon.TRACE_RESI)

    for uiY = 0; uiY < uiHeight; uiY++ {
        this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadResiTr(piPCM[uiY*uiWidth:], uiWidth, TLibCommon.TRACE_RESI)
    }

    this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XTracePredHeader(TLibCommon.TRACE_PRED)

    var pred [64]TLibCommon.Pel
    for uiY = 0; uiY < uiHeight; uiY++ {
        for uiX = 0; uiX < uiWidth; uiX++ {
            pred[uiX] = 0
        }
        this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadPredTr(pred[:], uiWidth, TLibCommon.TRACE_PRED)
    }

    this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XTraceRecoHeader(TLibCommon.TRACE_RECON)
    for uiY = 0; uiY < uiHeight; uiY++ {
        this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadRecoTr(piPCM[uiY*uiWidth:], uiWidth, TLibCommon.TRACE_RECON)
    }
    /*#endif*/

    if ttText == TLibCommon.TEXT_LUMA {
        uiPicStride = uint(pcCU.GetPic().GetPicYuvRec().GetStride())
        piPicReco = pcCU.GetPic().GetPicYuvRec().GetLumaAddr2(int(pcCU.GetAddr()), int(pcCU.GetZorderIdxInCU()+uiPartIdx))
        uiPcmLeftShiftBit = uint(TLibCommon.G_bitDepthY) - pcCU.GetSlice().GetSPS().GetPCMBitDepthLuma()
    } else {
        uiPicStride = uint(pcCU.GetPic().GetPicYuvRec().GetCStride())

        if ttText == TLibCommon.TEXT_CHROMA_U {
            piPicReco = pcCU.GetPic().GetPicYuvRec().GetCbAddr2(int(pcCU.GetAddr()), int(pcCU.GetZorderIdxInCU()+uiPartIdx))
        } else {
            piPicReco = pcCU.GetPic().GetPicYuvRec().GetCrAddr2(int(pcCU.GetAddr()), int(pcCU.GetZorderIdxInCU()+uiPartIdx))
        }
        uiPcmLeftShiftBit = uint(TLibCommon.G_bitDepthC) - pcCU.GetSlice().GetSPS().GetPCMBitDepthChroma()
    }

    for uiY = 0; uiY < uiHeight; uiY++ {
        for uiX = 0; uiX < uiWidth; uiX++ {
            piReco[uiY*uiStride+uiX] = (piPCM[uiY*uiWidth+uiX] << uiPcmLeftShiftBit)
            piPicReco[uiY*uiPicStride+uiX] = piReco[uiY*uiStride+uiX]
        }
        /*piPCM += uiWidth;
          piReco += uiStride;
          piPicReco += uiPicStride;*/
    }
}

func (this *TDecCu) xCopyToPic(pcCU *TLibCommon.TComDataCU, pcPic *TLibCommon.TComPic, uiZorderIdx, uiDepth uint) {
    uiCUAddr := pcCU.GetAddr()

    this.m_ppcYuvReco[uiDepth].CopyToPicYuv(pcPic.GetPicYuvRec(), uiCUAddr, uiZorderIdx, 0, 0)

    return
}

func (this *TDecCu) xIntraLumaRecQT(pcCU *TLibCommon.TComDataCU, uiTrDepth, uiAbsPartIdx uint, pcRecoYuv *TLibCommon.TComYuv, pcPredYuv *TLibCommon.TComYuv, pcResiYuv *TLibCommon.TComYuv) {
    uiFullDepth := uint(pcCU.GetDepth1(0)) + uiTrDepth
    uiTrMode := uint(pcCU.GetTransformIdx1(uiAbsPartIdx))
    if uiTrMode == uiTrDepth {
        this.xIntraRecLumaBlk(pcCU, uiTrDepth, uiAbsPartIdx, pcRecoYuv, pcPredYuv, pcResiYuv)
    } else {
        uiNumQPart := pcCU.GetPic().GetNumPartInCU() >> ((uiFullDepth + 1) << 1)
        for uiPart := uint(0); uiPart < 4; uiPart++ {
            this.xIntraLumaRecQT(pcCU, uiTrDepth+1, uiAbsPartIdx+uiPart*uiNumQPart, pcRecoYuv, pcPredYuv, pcResiYuv)
        }
    }
}
func (this *TDecCu) xIntraChromaRecQT(pcCU *TLibCommon.TComDataCU, uiTrDepth, uiAbsPartIdx uint, pcRecoYuv *TLibCommon.TComYuv, pcPredYuv *TLibCommon.TComYuv, pcResiYuv *TLibCommon.TComYuv) {
    uiFullDepth := uint(pcCU.GetDepth1(0)) + uiTrDepth
    uiTrMode := uint(pcCU.GetTransformIdx1(uiAbsPartIdx))
    if uiTrMode == uiTrDepth {
        this.xIntraRecChromaBlk(pcCU, uiTrDepth, uiAbsPartIdx, pcRecoYuv, pcPredYuv, pcResiYuv, 0)
        this.xIntraRecChromaBlk(pcCU, uiTrDepth, uiAbsPartIdx, pcRecoYuv, pcPredYuv, pcResiYuv, 1)
    } else {
        uiNumQPart := pcCU.GetPic().GetNumPartInCU() >> ((uiFullDepth + 1) << 1)
        for uiPart := uint(0); uiPart < 4; uiPart++ {
            this.xIntraChromaRecQT(pcCU, uiTrDepth+1, uiAbsPartIdx+uiPart*uiNumQPart, pcRecoYuv, pcPredYuv, pcResiYuv)
        }
    }
}

func (this *TDecCu) GetdQPFlag() bool {
    return this.m_bDecodeDQP
}
func (this *TDecCu) SetdQPFlag(b bool) {
    this.m_bDecodeDQP = b
}
func (this *TDecCu) xFillPCMBuffer(pCU *TLibCommon.TComDataCU, depth uint) {
    // Luma
    width := (pCU.GetSlice().GetSPS().GetMaxCUWidth() >> depth)
    height := (pCU.GetSlice().GetSPS().GetMaxCUHeight() >> depth)

    pPcmY := pCU.GetPCMSampleY()
    pRecoY := this.m_ppcYuvReco[depth].GetLumaAddr2(0, width)

    stride := this.m_ppcYuvReco[depth].GetStride()

    for y := uint(0); y < height; y++ {
        for x := uint(0); x < width; x++ {
            pPcmY[y*width+x] = pRecoY[y*stride+x]
        }
        //pPcmY += width;
        //pRecoY += stride;
    }

    // Cb and Cr
    widthC := (width >> 1)
    heightC := (height >> 1)

    pPcmCb := pCU.GetPCMSampleCb()
    pPcmCr := pCU.GetPCMSampleCr()
    pRecoCb := this.m_ppcYuvReco[depth].GetCbAddr()
    pRecoCr := this.m_ppcYuvReco[depth].GetCrAddr()

    strideC := this.m_ppcYuvReco[depth].GetCStride()

    for y := uint(0); y < heightC; y++ {
        for x := uint(0); x < widthC; x++ {
            pPcmCb[y*widthC+x] = pRecoCb[y*strideC+x]
            pPcmCr[y*widthC+x] = pRecoCr[y*strideC+x]
        }
        /*pPcmCr += widthC;
          pPcmCb += widthC;
          pRecoCb += strideC;
          pRecoCr += strideC;*/
    }
}

func (this *TDecCu) InvRecurTransformNxN(pcCU *TLibCommon.TComDataCU, pcYuvPred *TLibCommon.TComYuv, uiAbsPartIdx uint, eTxt TLibCommon.TextType, rpcResidual []TLibCommon.Pel,
    uiAddr, uiStride, uiWidth, uiHeight, uiMaxTrMode, uiTrMode uint, rpcCoeff []TLibCommon.TCoeff) {
    if pcCU.GetCbf3(uiAbsPartIdx, eTxt, uiTrMode) == 0 {
        /*#ifdef ENC_DEC_TRACE*/
        chroma := eTxt != TLibCommon.TEXT_LUMA
        blkX := int(TLibCommon.G_auiRasterToPelX[TLibCommon.G_auiZscanToRaster[uiAbsPartIdx]]) >> TLibCommon.B2U(chroma)
        blkY := int(TLibCommon.G_auiRasterToPelY[TLibCommon.G_auiZscanToRaster[uiAbsPartIdx]]) >> TLibCommon.B2U(chroma)

        this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XTraceTUHeader(TLibCommon.TRACE_TU)

        if eTxt != TLibCommon.TEXT_LUMA {
            this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(int(eTxt)-1, "tu_color", TLibCommon.TRACE_TU)
        } else {
            this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(0, "tu_color", TLibCommon.TRACE_TU)
        }
        this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(blkX, "tu_x", TLibCommon.TRACE_TU)
        this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(blkY, "tu_y", TLibCommon.TRACE_TU)
        this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(int(uiWidth), "tu_width", TLibCommon.TRACE_TU)
        this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(int(uiHeight), "tu_height", TLibCommon.TRACE_TU)

        this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XTraceCoefHeader(TLibCommon.TRACE_COEF)

        var piCoef [64]TLibCommon.TCoeff
        for k := uint(0); k < uiHeight; k++ {
            for i := uint(0); i < uiWidth; i++ {
                piCoef[i] = 0
            }
            this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadCeofTr(piCoef[:], uiWidth, TLibCommon.TRACE_COEF)
        }
        this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XTraceResiHeader(TLibCommon.TRACE_RESI)

        var piResi [64]TLibCommon.Pel
        for k := uint(0); k < uiHeight; k++ {
            for i := uint(0); i < uiWidth; i++ {
                piResi[i] = 0
            }
            this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadResiTr(piResi[:], uiWidth, TLibCommon.TRACE_RESI)
        }

        var piPred []TLibCommon.Pel
        var uiPredStride uint
        if eTxt == TLibCommon.TEXT_LUMA {
            piPred = pcYuvPred.GetLumaAddr1(uiAbsPartIdx)
            uiPredStride = pcYuvPred.GetStride()
        } else if eTxt == TLibCommon.TEXT_CHROMA_U {
            piPred = pcYuvPred.GetCbAddr1(uiAbsPartIdx)
            uiPredStride = pcYuvPred.GetCStride()
        } else {
            piPred = pcYuvPred.GetCrAddr1(uiAbsPartIdx)
            uiPredStride = pcYuvPred.GetCStride()
        }

        this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XTracePredHeader(TLibCommon.TRACE_PRED)

        pPred := piPred
        for k := uint(0); k < uiHeight; k++ {
            this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadPredTr(pPred[k*uiPredStride:], uiWidth, TLibCommon.TRACE_PRED)
            //pPred += uiStride;
        }

        this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XTraceRecoHeader(TLibCommon.TRACE_RECON)

        pReco := piPred
        for k := uint(0); k < uiHeight; k++ {
            this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadRecoTr(pReco[k*uiPredStride:], uiWidth, TLibCommon.TRACE_RECON)
            //pReco += uiStride;
        }
        /*#endif*/
        return
    }
    stopTrMode := uint(pcCU.GetTransformIdx1(uiAbsPartIdx))

    if uiTrMode == stopTrMode {
        uiDepth := uint(pcCU.GetDepth1(uiAbsPartIdx)) + uiTrMode
        uiLog2TrSize := TLibCommon.G_aucConvertToBit[pcCU.GetSlice().GetSPS().GetMaxCUWidth()>>uiDepth] + 2
        if eTxt != TLibCommon.TEXT_LUMA && uiLog2TrSize == 2 {
            uiQPDiv := pcCU.GetPic().GetNumPartInCU() >> ((uiDepth - 1) << 1)
            if (uiAbsPartIdx % uiQPDiv) != 0 {
                return
            }
            uiWidth <<= 1
            uiHeight <<= 1
        }
        pResi := rpcResidual[uiAddr:]
        /*#ifdef ENC_DEC_TRACE*/
        chroma := eTxt != TLibCommon.TEXT_LUMA
        blkX := int(TLibCommon.G_auiRasterToPelX[TLibCommon.G_auiZscanToRaster[uiAbsPartIdx]]) >> TLibCommon.B2U(chroma)
        blkY := int(TLibCommon.G_auiRasterToPelY[TLibCommon.G_auiZscanToRaster[uiAbsPartIdx]]) >> TLibCommon.B2U(chroma)

        this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XTraceTUHeader(TLibCommon.TRACE_TU)

        if eTxt != TLibCommon.TEXT_LUMA {
            this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(int(eTxt)-1, "tu_color", TLibCommon.TRACE_TU)
        } else {
            this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(0, "tu_color", TLibCommon.TRACE_TU)
        }
        this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(blkX, "tu_x", TLibCommon.TRACE_TU)
        this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(blkY, "tu_y", TLibCommon.TRACE_TU)
        this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(int(uiWidth), "tu_width", TLibCommon.TRACE_TU)
        this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadAeTr(int(uiHeight), "tu_height", TLibCommon.TRACE_TU)
        /*#endif*/

        var scalingListType int
        if pcCU.IsIntra(uiAbsPartIdx) {
            scalingListType = 0 + TLibCommon.G_eTTable[int(eTxt)]
        } else {
            scalingListType = 3 + TLibCommon.G_eTTable[int(eTxt)]
        }
        //assert(scalingListType < 6);
        /*#ifdef ENC_DEC_TRACE*/
        this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XTraceCoefHeader(TLibCommon.TRACE_COEF)

        for k := uint(0); k < uiHeight; k++ {
            this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadCeofTr(rpcCoeff[k*uiWidth:], uiWidth, TLibCommon.TRACE_COEF)
        }
        /*#endif*/
        this.m_pcTrQuant.InvtransformNxN(pcCU.GetCUTransquantBypass1(uiAbsPartIdx), eTxt, TLibCommon.REG_DCT, pResi, uiStride, rpcCoeff, uiWidth, uiHeight, scalingListType, pcCU.GetTransformSkip2(uiAbsPartIdx, eTxt))
        /*#ifdef ENC_DEC_TRACE*/
        this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XTraceResiHeader(TLibCommon.TRACE_RESI)

        for k := uint(0); k < uiHeight; k++ {
            this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadResiTr(pResi[k*uiStride:], uiWidth, TLibCommon.TRACE_RESI)
        }
        /*#endif*/
        /*#ifdef ENC_DEC_TRACE*/
        var piPred []TLibCommon.Pel
        if eTxt == TLibCommon.TEXT_LUMA {
            piPred = pcYuvPred.GetLumaAddr1(uiAbsPartIdx)
        } else if eTxt == TLibCommon.TEXT_CHROMA_U {
            piPred = pcYuvPred.GetCbAddr1(uiAbsPartIdx)
        } else {
            piPred = pcYuvPred.GetCrAddr1(uiAbsPartIdx)
        }
        this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XTracePredHeader(TLibCommon.TRACE_PRED)

        pPred := piPred
        for k := uint(0); k < uiHeight; k++ {
            this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadPredTr(pPred[k*uiStride:], uiWidth, TLibCommon.TRACE_PRED)
            //pPred += uiStride;
        }

        this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XTraceRecoHeader(TLibCommon.TRACE_RECON)

        var pReco [64]TLibCommon.Pel
        for k := uint(0); k < uiHeight; k++ {
            for j := uint(0); j < uiWidth; j++ {
                if !chroma {
                    pReco[j] = TLibCommon.Pel(TLibCommon.ClipY(piPred[k*uiStride+j] + pResi[k*uiStride+j]))
                } else {
                    pReco[j] = TLibCommon.Pel(TLibCommon.ClipC(piPred[k*uiStride+j] + pResi[k*uiStride+j]))
                }
            }
            this.m_pcEntropyDecoder.m_pcEntropyDecoderIf.XReadRecoTr(pReco[:], uiWidth, TLibCommon.TRACE_RECON)
        }
        /*#endif*/
    } else {
        uiTrMode++
        uiWidth >>= 1
        uiHeight >>= 1
        trWidth := uiWidth
        trHeight := uiHeight
        uiAddrOffset := trHeight * uiStride
        uiCoefOffset := trWidth * trHeight
        uiPartOffset := pcCU.GetTotalNumPart() >> (uiTrMode << 1)
        {
            this.InvRecurTransformNxN(pcCU, pcYuvPred, uiAbsPartIdx, eTxt, rpcResidual, uiAddr, uiStride, uiWidth, uiHeight, uiMaxTrMode, uiTrMode, rpcCoeff)
            uiAbsPartIdx += uiPartOffset //rpcCoeff += uiCoefOffset;
            this.InvRecurTransformNxN(pcCU, pcYuvPred, uiAbsPartIdx, eTxt, rpcResidual, uiAddr+trWidth, uiStride, uiWidth, uiHeight, uiMaxTrMode, uiTrMode, rpcCoeff[uiCoefOffset:])
            uiAbsPartIdx += uiPartOffset //rpcCoeff += uiCoefOffset;
            this.InvRecurTransformNxN(pcCU, pcYuvPred, uiAbsPartIdx, eTxt, rpcResidual, uiAddr+uiAddrOffset, uiStride, uiWidth, uiHeight, uiMaxTrMode, uiTrMode, rpcCoeff[uiCoefOffset*2:])
            uiAbsPartIdx += uiPartOffset //rpcCoeff += uiCoefOffset;
            this.InvRecurTransformNxN(pcCU, pcYuvPred, uiAbsPartIdx, eTxt, rpcResidual, uiAddr+uiAddrOffset+trWidth, uiStride, uiWidth, uiHeight, uiMaxTrMode, uiTrMode, rpcCoeff[uiCoefOffset*3:])
        }
    }
}
