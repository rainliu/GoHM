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
	"io"
    "gohm/TLibCommon"
)

// ====================================================================================================================
// Class definition
// ====================================================================================================================

/// entropy encoder pure class
type TEncEntropyIf interface {
	XTraceLCUHeader(traceLevel uint)
    XTraceCUHeader(traceLevel uint)
    XTracePUHeader(traceLevel uint)
    XTraceTUHeader(traceLevel uint)
    XTraceCoefHeader(traceLevel uint)
    XTraceResiHeader(traceLevel uint)
    XTracePredHeader(traceLevel uint)
    XTraceRecoHeader(traceLevel uint)
    XReadAeTr(Value int, pSymbolName string, traceLevel uint)
    XReadCeofTr(pCoeff []TLibCommon.TCoeff, uiWidth, traceLevel uint)
    XReadResiTr(pPel []TLibCommon.Pel, uiWidth, traceLevel uint)
    XReadPredTr(pPel []TLibCommon.Pel, uiWidth, traceLevel uint)
    XReadRecoTr(pPel []TLibCommon.Pel, uiWidth, traceLevel uint)
    
	DTRACE_CABAC_F(x float32)
    DTRACE_CABAC_V(x uint)
    DTRACE_CABAC_VL(x uint)
    DTRACE_CABAC_T(x string)
    DTRACE_CABAC_X(x uint)
    DTRACE_CABAC_N()
    
    resetEntropy()
    determineCabacInitIdx()
    setBitstream(p TLibCommon.TComBitIf)
    setTraceFile(traceFile io.Writer)
    setSlice(p *TLibCommon.TComSlice)
    resetBits()
    resetCoeffCost()
    getNumberOfWrittenBits() uint
    getCoeffCost() uint
    codeVPS(pcVPS *TLibCommon.TComVPS)
    codeSPS(pcSPS *TLibCommon.TComSPS)
    codePPS(pcPPS *TLibCommon.TComPPS)
    codeSliceHeader(pcSlice *TLibCommon.TComSlice)
    codeTilesWPPEntryPoint(pSlice *TLibCommon.TComSlice)
    codeTerminatingBit(uilsLast uint)
    codeSliceFinish()
    codeMVPIdx(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint, eRefList TLibCommon.RefPicList)
    codeScalingList(scalingList *TLibCommon.TComScalingList)
    codeCUTransquantBypassFlag(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint)
    codeSkipFlag(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint)
    codeMergeFlag(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint)
    codeMergeIndex(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint)
    codeSplitFlag(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint, uiDepth uint)
    codePartSize(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint, uiDepth uint)
    codePredMode(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint)
    codeIPCMInfo(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint)
    codeTransformSubdivFlag(uiSymbol, uiCtx uint)
    codeQtCbf(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint, eType TLibCommon.TextType, uiTrDepth uint)
    codeQtRootCbf(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint)
    codeQtCbfZero(pcCU *TLibCommon.TComDataCU, eType TLibCommon.TextType, uiTrDepth uint)
    codeQtRootCbfZero(pcCU *TLibCommon.TComDataCU)
    codeIntraDirLumaAng(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint, isMultiplePU bool)
    codeIntraDirChroma(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint)
    codeInterDir(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint)
    codeRefFrmIdx(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint, eRefList TLibCommon.RefPicList)
    codeMvd(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint, eRefList TLibCommon.RefPicList)
    codeDeltaQP(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint)
    codeCoeffNxN(pcCU *TLibCommon.TComDataCU, pcCoef []TLibCommon.TCoeff, uiAbsPartIdx, uiWidth, uiHeight, uiDepth uint, eTType TLibCommon.TextType)
    codeTransformSkipFlags(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint, width, height uint, eTType TLibCommon.TextType)
    codeSAOSign(code uint)
    codeSaoMaxUvlc(code, maxSymbol uint)
    codeSaoMerge(uiCode uint)
    codeSaoTypeIdx(uiCode uint)
    codeSaoUflc(uiLength, uiCode uint)
    estBit(pcEstBitsSbac *TLibCommon.EstBitsSbacStruct, width, height int, eTType TLibCommon.TextType)
    updateContextTables3(eSliceType TLibCommon.SliceType, iQp int, bExecuteFinish bool)
    updateContextTables2(eSliceType TLibCommon.SliceType, iQp int)
    codeDFFlag(uiCode uint, pSymbolName string)
    codeDFSvlc(iCode int, pSymbolName string)
    getEncBinIf() TEncBinIf
}

/// entropy encoder class
type TEncEntropy struct {
    //private:
    m_uiBakAbsPartIdx   uint
    m_uiBakChromaOffset uint
    m_bakAbsPartIdxCU   uint
    m_pcEntropyCoderIf  TEncEntropyIf
}

func NewTEncEntropy() *TEncEntropy{
	return &TEncEntropy{};
}

//public:
func (this *TEncEntropy) setEntropyCoder(e TEncEntropyIf, pcSlice *TLibCommon.TComSlice, traceFile io.Writer) {
    this.m_pcEntropyCoderIf = e
    this.m_pcEntropyCoderIf.setSlice(pcSlice)
    this.m_pcEntropyCoderIf.setTraceFile(traceFile)
}
func (this *TEncEntropy) setBitstream(p TLibCommon.TComBitIf) {
    this.m_pcEntropyCoderIf.setBitstream(p)
}

func (this *TEncEntropy) setTraceFile(traceFile io.Writer) {
    this.m_pcEntropyCoderIf.setTraceFile(traceFile)
}
func (this *TEncEntropy) resetBits()      { this.m_pcEntropyCoderIf.resetBits() }
func (this *TEncEntropy) resetCoeffCost() { this.m_pcEntropyCoderIf.resetCoeffCost() }
func (this *TEncEntropy) getNumberOfWrittenBits() uint {
    return this.m_pcEntropyCoderIf.getNumberOfWrittenBits()
}
func (this *TEncEntropy) getCoeffCost() uint {
    return this.m_pcEntropyCoderIf.getCoeffCost()
}
func (this *TEncEntropy) resetEntropy()          { this.m_pcEntropyCoderIf.resetEntropy() }
func (this *TEncEntropy) determineCabacInitIdx() { this.m_pcEntropyCoderIf.determineCabacInitIdx() }
func (this *TEncEntropy) encodeSliceHeader(pcSlice *TLibCommon.TComSlice) {
    if pcSlice.GetSPS().GetUseSAO() {
        saoParam := pcSlice.GetPic().GetPicSym().GetSaoParam()
        pcSlice.SetSaoEnabledFlag(saoParam.SaoFlag[0])
        pcSlice.SetSaoEnabledFlagChroma(saoParam.SaoFlag[1])
    }

    this.m_pcEntropyCoderIf.codeSliceHeader(pcSlice)
    return
}
func (this *TEncEntropy) encodeTilesWPPEntryPoint(pcSlice *TLibCommon.TComSlice) {
    this.m_pcEntropyCoderIf.codeTilesWPPEntryPoint(pcSlice)
}
func (this *TEncEntropy) encodeTerminatingBit(uiIsLast uint) {
    this.m_pcEntropyCoderIf.codeTerminatingBit(uiIsLast)
}
func (this *TEncEntropy) encodeSliceFinish() { this.m_pcEntropyCoderIf.codeSliceFinish() }
func (this *TEncEntropy) encodeVPS(pcVPS *TLibCommon.TComVPS) {
    this.m_pcEntropyCoderIf.codeVPS(pcVPS)
}
func (this *TEncEntropy) encodeSPS(pcSPS *TLibCommon.TComSPS) {
    this.m_pcEntropyCoderIf.codeSPS(pcSPS)
}
func (this *TEncEntropy) encodePPS(pcPPS *TLibCommon.TComPPS) {
    this.m_pcEntropyCoderIf.codePPS(pcPPS)
}
func (this *TEncEntropy) encodeSplitFlag(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint, uiDepth uint, bRD bool) { //= false );
    if bRD {
        uiAbsPartIdx = 0
    }

    this.m_pcEntropyCoderIf.codeSplitFlag(pcCU, uiAbsPartIdx, uiDepth)
}
func (this *TEncEntropy) encodeCUTransquantBypassFlag(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint, bRD bool) { //= false );
    if bRD {
        uiAbsPartIdx = 0
    }

    this.m_pcEntropyCoderIf.codeCUTransquantBypassFlag(pcCU, uiAbsPartIdx)
}
func (this *TEncEntropy) encodeSkipFlag(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint, bRD bool) { //= false );
    if pcCU.GetSlice().IsIntra() {
        return
    }
    if bRD {
        uiAbsPartIdx = 0
    }

    this.m_pcEntropyCoderIf.codeSkipFlag(pcCU, uiAbsPartIdx)
}
func (this *TEncEntropy) encodePUWise(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint, bRD bool) { //= false );
    if bRD {
        uiAbsPartIdx = 0
    }

    ePartSize := pcCU.GetPartitionSize1(uiAbsPartIdx)
    var uiNumPU uint
    if ePartSize == TLibCommon.SIZE_2Nx2N {
        uiNumPU = 1
    } else if ePartSize == TLibCommon.SIZE_NxN {
        uiNumPU = 4
    } else {
        uiNumPU = 2
    }
    uiDepth := uint(pcCU.GetDepth1(uiAbsPartIdx))
    //fmt.Printf("ePartSize=%d\n",ePartSize);
    uiPUOffset := (TLibCommon.G_auiPUOffset[uint(ePartSize)] << ((pcCU.GetSlice().GetSPS().GetMaxCUDepth() - uiDepth) << 1)) >> 4

    uiSubPartIdx := uiAbsPartIdx
    for uiPartIdx := uint(0); uiPartIdx < uiNumPU; uiPartIdx++ {
        this.encodeMergeFlag(pcCU, uiSubPartIdx)
        if pcCU.GetMergeFlag1(uiSubPartIdx) {
            this.encodeMergeIndex(pcCU, uiSubPartIdx, false)
        } else {
            this.encodeInterDirPU(pcCU, uiSubPartIdx)
            for uiRefListIdx := 0; uiRefListIdx < 2; uiRefListIdx++ {
                if pcCU.GetSlice().GetNumRefIdx(TLibCommon.RefPicList(uiRefListIdx)) > 0 {
                    this.encodeRefFrmIdxPU(pcCU, uiSubPartIdx, TLibCommon.RefPicList(uiRefListIdx))
                    this.encodeMvdPU(pcCU, uiSubPartIdx, TLibCommon.RefPicList(uiRefListIdx))
                    this.encodeMVPIdxPU(pcCU, uiSubPartIdx, TLibCommon.RefPicList(uiRefListIdx))
                }
            }
        }
        uiSubPartIdx += uiPUOffset
    }

    return
}
func (this *TEncEntropy) encodeInterDirPU(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint) {
    if !pcCU.GetSlice().IsInterB() {
        return
    }

    this.m_pcEntropyCoderIf.codeInterDir(pcCU, uiAbsPartIdx)
    return
}
func (this *TEncEntropy) encodeRefFrmIdxPU(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint, eRefList TLibCommon.RefPicList) {
    if pcCU.GetSlice().GetNumRefIdx(eRefList) == 1 {
        return
    }

    if pcCU.GetInterDir1(uiAbsPartIdx)&(1<<eRefList) != 0 {
        this.m_pcEntropyCoderIf.codeRefFrmIdx(pcCU, uiAbsPartIdx, eRefList)
    }

    return
}
func (this *TEncEntropy) encodeMvdPU(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint, eRefList TLibCommon.RefPicList) {
    if pcCU.GetInterDir1(uiAbsPartIdx)&(1<<eRefList) != 0 {
        this.m_pcEntropyCoderIf.codeMvd(pcCU, uiAbsPartIdx, eRefList)
    }
    return
}
func (this *TEncEntropy) encodeMVPIdxPU(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint, eRefList TLibCommon.RefPicList) {
    if (pcCU.GetInterDir1(uiAbsPartIdx) & (1 << eRefList)) != 0 {
        this.m_pcEntropyCoderIf.codeMVPIdx(pcCU, uiAbsPartIdx, eRefList)
    }

    return
}
func (this *TEncEntropy) encodeMergeFlag(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint) {
    this.m_pcEntropyCoderIf.codeMergeFlag(pcCU, uiAbsPartIdx)
}
func (this *TEncEntropy) encodeMergeIndex(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint, bRD bool) { //= false );
    if bRD {
        uiAbsPartIdx = 0
        //assert( pcCU.GetPartitionSize(uiAbsPartIdx) == SIZE_2Nx2N );
    }
    this.m_pcEntropyCoderIf.codeMergeIndex(pcCU, uiAbsPartIdx)
}
func (this *TEncEntropy) encodePredMode(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint, bRD bool) { // = false );
    if bRD {
        uiAbsPartIdx = 0
    }

    if pcCU.GetSlice().IsIntra() {
        return
    }

    this.m_pcEntropyCoderIf.codePredMode(pcCU, uiAbsPartIdx)
}
func (this *TEncEntropy) encodePartSize(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint, uiDepth uint, bRD bool) { // = false );
    if bRD {
        uiAbsPartIdx = 0
    }

    this.m_pcEntropyCoderIf.codePartSize(pcCU, uiAbsPartIdx, uiDepth)
}
func (this *TEncEntropy) encodeIPCMInfo(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint, bRD bool) { // = false );
    if !pcCU.GetSlice().GetSPS().GetUsePCM() ||
        pcCU.GetWidth1(uiAbsPartIdx) > (1<<pcCU.GetSlice().GetSPS().GetPCMLog2MaxSize()) ||
        pcCU.GetWidth1(uiAbsPartIdx) < (1<<pcCU.GetSlice().GetSPS().GetPCMLog2MinSize()) {
        return
    }

    if bRD {
        uiAbsPartIdx = 0
    }

    this.m_pcEntropyCoderIf.codeIPCMInfo(pcCU, uiAbsPartIdx)
}
func (this *TEncEntropy) encodePredInfo(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint, bRD bool) { // = false );
    if bRD {
        uiAbsPartIdx = 0
    }
    //fmt.Printf("uiAbsPartIdx=%d\n", uiAbsPartIdx);
    if pcCU.IsIntra(uiAbsPartIdx) { // If it is Intra mode, encode intra prediction mode.
        this.encodeIntraDirModeLuma(pcCU, uiAbsPartIdx, true)
        this.encodeIntraDirModeChroma(pcCU, uiAbsPartIdx, bRD)
    } else {
        this.encodePUWise(pcCU, uiAbsPartIdx, bRD)
    }
}
func (this *TEncEntropy) encodeIntraDirModeLuma(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint, isMultiplePU bool) { // = false );
    this.m_pcEntropyCoderIf.codeIntraDirLumaAng(pcCU, uiAbsPartIdx, isMultiplePU)
}
func (this *TEncEntropy) encodeIntraDirModeChroma(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint, bRD bool) { // = false );
    if bRD {
        uiAbsPartIdx = 0
    }

    this.m_pcEntropyCoderIf.codeIntraDirChroma(pcCU, uiAbsPartIdx)
}
func (this *TEncEntropy) encodeTransformSubdivFlag(uiSymbol, uiCtx uint) {
    this.m_pcEntropyCoderIf.codeTransformSubdivFlag(uiSymbol, uiCtx)
}
func (this *TEncEntropy) encodeQtCbf(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint, eType TLibCommon.TextType, uiTrDepth uint) {
    this.m_pcEntropyCoderIf.codeQtCbf( pcCU, uiAbsPartIdx, eType, uiTrDepth )
}
func (this *TEncEntropy) encodeQtCbfZero(pcCU *TLibCommon.TComDataCU, eType TLibCommon.TextType, uiTrDepth uint) {
    this.m_pcEntropyCoderIf.codeQtCbfZero(pcCU, eType, uiTrDepth)
}
func (this *TEncEntropy) encodeQtRootCbfZero(pcCU *TLibCommon.TComDataCU) {
    this.m_pcEntropyCoderIf.codeQtRootCbfZero(pcCU)
}
func (this *TEncEntropy) encodeQtRootCbf(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint) {
    this.m_pcEntropyCoderIf.codeQtRootCbf(pcCU, uiAbsPartIdx)
}
func (this *TEncEntropy) encodeQP(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint, bRD bool) { //= false );
    if bRD {
        uiAbsPartIdx = 0
    }

    if pcCU.GetSlice().GetPPS().GetUseDQP() {
        this.m_pcEntropyCoderIf.codeDeltaQP(pcCU, uiAbsPartIdx)
    }
}
func (this *TEncEntropy) updateContextTables3(eSliceType TLibCommon.SliceType, iQp int, bExecuteFinish bool) {
    this.m_pcEntropyCoderIf.updateContextTables3(eSliceType, iQp, bExecuteFinish)
}
func (this *TEncEntropy) updateContextTables2(eSliceType TLibCommon.SliceType, iQp int) {
    this.m_pcEntropyCoderIf.updateContextTables3(eSliceType, iQp, true)
}
func (this *TEncEntropy) encodeScalingList(scalingList *TLibCommon.TComScalingList) {
    this.m_pcEntropyCoderIf.codeScalingList(scalingList)
}
func (this *TEncEntropy) xEncodeTransform(pcCU *TLibCommon.TComDataCU, offsetLuma, offsetChroma, uiAbsPartIdx, uiDepth, width, height, uiTrIdx uint, bCodeDQP *bool) {
    uiSubdiv := uint(TLibCommon.B2U(uint(pcCU.GetTransformIdx1(uiAbsPartIdx) + pcCU.GetDepth1(uiAbsPartIdx)) > uiDepth))
    uiLog2TrafoSize := uint(TLibCommon.G_aucConvertToBit[pcCU.GetSlice().GetSPS().GetMaxCUWidth()]) + 2 - uiDepth
    cbfY := pcCU.GetCbf3(uiAbsPartIdx, TLibCommon.TEXT_LUMA, uiTrIdx)
    cbfU := pcCU.GetCbf3(uiAbsPartIdx, TLibCommon.TEXT_CHROMA_U, uiTrIdx)
    cbfV := pcCU.GetCbf3(uiAbsPartIdx, TLibCommon.TEXT_CHROMA_V, uiTrIdx)
	
	//fmt.Print("Enter xEncodeTransform\n");// with uiSubdiv=%d, uiAbsPartIdx=%d, uiDepth=%d\n", uiSubdiv, uiAbsPartIdx, uiDepth);
  
    if uiTrIdx == 0 {
        this.m_bakAbsPartIdxCU = uiAbsPartIdx
    }
    if uiLog2TrafoSize == 2 {
        partNum := pcCU.GetPic().GetNumPartInCU() >> ((uiDepth - 1) << 1)
        if (uiAbsPartIdx % partNum) == 0 {
            this.m_uiBakAbsPartIdx = uiAbsPartIdx
            this.m_uiBakChromaOffset = offsetChroma
        } else if (uiAbsPartIdx % partNum) == (partNum - 1) {
            cbfU = pcCU.GetCbf3(this.m_uiBakAbsPartIdx, TLibCommon.TEXT_CHROMA_U, uiTrIdx)
            cbfV = pcCU.GetCbf3(this.m_uiBakAbsPartIdx, TLibCommon.TEXT_CHROMA_V, uiTrIdx)
        }
    }

    if pcCU.GetPredictionMode1(uiAbsPartIdx) == TLibCommon.MODE_INTRA && pcCU.GetPartitionSize1(uiAbsPartIdx) == TLibCommon.SIZE_NxN && uiDepth == uint(pcCU.GetDepth1(uiAbsPartIdx)) {
        //assert( uiSubdiv );
    } else if pcCU.GetPredictionMode1(uiAbsPartIdx) == TLibCommon.MODE_INTER && (pcCU.GetPartitionSize1(uiAbsPartIdx) != TLibCommon.SIZE_2Nx2N) && uiDepth == uint(pcCU.GetDepth1(uiAbsPartIdx)) && (pcCU.GetSlice().GetSPS().GetQuadtreeTUMaxDepthInter() == 1) {
        if uiLog2TrafoSize > pcCU.GetQuadtreeTULog2MinSizeInCU(uiAbsPartIdx) {
            //assert( uiSubdiv );
        } else {
            //assert(!uiSubdiv );
        }
    } else if uiLog2TrafoSize > pcCU.GetSlice().GetSPS().GetQuadtreeTULog2MaxSize() {
        //assert( uiSubdiv );
    } else if uiLog2TrafoSize == pcCU.GetSlice().GetSPS().GetQuadtreeTULog2MinSize() {
        //assert( !uiSubdiv );
    } else if uiLog2TrafoSize == pcCU.GetQuadtreeTULog2MinSizeInCU(uiAbsPartIdx) {
        //assert( !uiSubdiv );
    } else {
        //assert( uiLog2TrafoSize > pcCU.GetQuadtreeTULog2MinSizeInCU(uiAbsPartIdx) );
        this.m_pcEntropyCoderIf.codeTransformSubdivFlag(uiSubdiv, 5-uiLog2TrafoSize)
    }

    uiTrDepthCurr := uiDepth - uint(pcCU.GetDepth1(uiAbsPartIdx))
    bFirstCbfOfCU := uiTrDepthCurr == 0
    if bFirstCbfOfCU || uiLog2TrafoSize > 2 {
        if bFirstCbfOfCU || pcCU.GetCbf3(uiAbsPartIdx, TLibCommon.TEXT_CHROMA_U, uiTrDepthCurr-1) != 0 {
            this.m_pcEntropyCoderIf.codeQtCbf(pcCU, uiAbsPartIdx, TLibCommon.TEXT_CHROMA_U, uiTrDepthCurr)
        }
        if bFirstCbfOfCU || pcCU.GetCbf3(uiAbsPartIdx, TLibCommon.TEXT_CHROMA_V, uiTrDepthCurr-1) != 0 {
            this.m_pcEntropyCoderIf.codeQtCbf(pcCU, uiAbsPartIdx, TLibCommon.TEXT_CHROMA_V, uiTrDepthCurr)
        }
    } else if uiLog2TrafoSize == 2 {
        //assert( pcCU.GetCbf( uiAbsPartIdx, TEXT_CHROMA_U, uiTrDepthCurr ) == pcCU.GetCbf( uiAbsPartIdx, TEXT_CHROMA_U, uiTrDepthCurr - 1 ) );
        //assert( pcCU.GetCbf( uiAbsPartIdx, TEXT_CHROMA_V, uiTrDepthCurr ) == pcCU.GetCbf( uiAbsPartIdx, TEXT_CHROMA_V, uiTrDepthCurr - 1 ) );
    }

    if uiSubdiv != 0 {
        var size uint
        width >>= 1
        height >>= 1
        size = width * height
        uiTrIdx++
        uiDepth++
        partNum := pcCU.GetPic().GetNumPartInCU() >> (uiDepth << 1)

        this.xEncodeTransform(pcCU, offsetLuma, offsetChroma, uiAbsPartIdx, uiDepth, width, height, uiTrIdx, bCodeDQP)

        uiAbsPartIdx += partNum
        offsetLuma += size
        offsetChroma += (size >> 2)
        this.xEncodeTransform(pcCU, offsetLuma, offsetChroma, uiAbsPartIdx, uiDepth, width, height, uiTrIdx, bCodeDQP)

        uiAbsPartIdx += partNum
        offsetLuma += size
        offsetChroma += (size >> 2)
        this.xEncodeTransform(pcCU, offsetLuma, offsetChroma, uiAbsPartIdx, uiDepth, width, height, uiTrIdx, bCodeDQP)

        uiAbsPartIdx += partNum
        offsetLuma += size
        offsetChroma += (size >> 2)
        this.xEncodeTransform(pcCU, offsetLuma, offsetChroma, uiAbsPartIdx, uiDepth, width, height, uiTrIdx, bCodeDQP)
    } else {
        /*DTRACE_CABAC_VL( g_nSymbolCounter++ );*/
        this.m_pcEntropyCoderIf.DTRACE_CABAC_T( "\tTrIdx: abspart=" );
        this.m_pcEntropyCoderIf.DTRACE_CABAC_V( uiAbsPartIdx );
        this.m_pcEntropyCoderIf.DTRACE_CABAC_T( "\tdepth=" );
        this.m_pcEntropyCoderIf.DTRACE_CABAC_V( uiDepth );
        this.m_pcEntropyCoderIf.DTRACE_CABAC_T( "\ttrdepth=" );
        this.m_pcEntropyCoderIf.DTRACE_CABAC_V( uint(pcCU.GetTransformIdx1( uiAbsPartIdx )) );
        this.m_pcEntropyCoderIf.DTRACE_CABAC_T( "\n" );
        

        if pcCU.GetPredictionMode1(uiAbsPartIdx) != TLibCommon.MODE_INTRA && uiDepth == uint(pcCU.GetDepth1(uiAbsPartIdx)) && pcCU.GetCbf3(uiAbsPartIdx, TLibCommon.TEXT_CHROMA_U, 0) == 0 && pcCU.GetCbf3(uiAbsPartIdx, TLibCommon.TEXT_CHROMA_V, 0) == 0 {
            //assert( pcCU.GetCbf( uiAbsPartIdx, TLibCommon.TEXT_LUMA, 0 ) );
        } else {
            this.m_pcEntropyCoderIf.codeQtCbf(pcCU, uiAbsPartIdx, TLibCommon.TEXT_LUMA, uint(pcCU.GetTransformIdx1(uiAbsPartIdx)))
        }

        if cbfY != 0 || cbfU != 0 || cbfV != 0 {
            // dQP: only for LCU once
            if pcCU.GetSlice().GetPPS().GetUseDQP() {
                if *bCodeDQP {
                    this.encodeQP(pcCU, this.m_bakAbsPartIdxCU, false)
                    *bCodeDQP = false
                }
            }
        }
        if cbfY != 0 {
            trWidth := width
            trHeight := height
            this.m_pcEntropyCoderIf.codeCoeffNxN(pcCU, pcCU.GetCoeffY()[offsetLuma:], uiAbsPartIdx, trWidth, trHeight, uiDepth, TLibCommon.TEXT_LUMA)
        }
        if uiLog2TrafoSize > 2 {
            trWidth := width >> 1
            trHeight := height >> 1
            if cbfU != 0 {
                this.m_pcEntropyCoderIf.codeCoeffNxN(pcCU, pcCU.GetCoeffCb()[offsetChroma:], uiAbsPartIdx, trWidth, trHeight, uiDepth, TLibCommon.TEXT_CHROMA_U)
            }
            if cbfV != 0 {
                this.m_pcEntropyCoderIf.codeCoeffNxN(pcCU, pcCU.GetCoeffCr()[offsetChroma:], uiAbsPartIdx, trWidth, trHeight, uiDepth, TLibCommon.TEXT_CHROMA_V)
            }
        } else {
            partNum := pcCU.GetPic().GetNumPartInCU() >> ((uiDepth - 1) << 1)
            if (uiAbsPartIdx % partNum) == (partNum - 1) {
                trWidth := width
                trHeight := height
                if cbfU != 0 {
                    this.m_pcEntropyCoderIf.codeCoeffNxN(pcCU, pcCU.GetCoeffCb()[this.m_uiBakChromaOffset:], this.m_uiBakAbsPartIdx, trWidth, trHeight, uiDepth, TLibCommon.TEXT_CHROMA_U)
                }
                if cbfV != 0 {
                    this.m_pcEntropyCoderIf.codeCoeffNxN(pcCU, pcCU.GetCoeffCr()[this.m_uiBakChromaOffset:], this.m_uiBakAbsPartIdx, trWidth, trHeight, uiDepth, TLibCommon.TEXT_CHROMA_V)
                }
            }
        }
    }
    //fmt.Print("Exit xEncodeTransform\n");
}
func (this *TEncEntropy) encodeCoeff(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth, uiWidth, uiHeight uint, bCodeDQP *bool) {
    uiMinCoeffSize := pcCU.GetPic().GetMinCUWidth() * pcCU.GetPic().GetMinCUHeight()
    uiLumaOffset := uiMinCoeffSize * uiAbsPartIdx
    uiChromaOffset := uiLumaOffset >> 2

    if pcCU.IsIntra(uiAbsPartIdx) {
        /*DTRACE_CABAC_VL( g_nSymbolCounter++ )*/
        this.m_pcEntropyCoderIf.DTRACE_CABAC_T( "\tdecodeTransformIdx()\tCUDepth=" )
        this.m_pcEntropyCoderIf.DTRACE_CABAC_V( uiDepth )
        this.m_pcEntropyCoderIf.DTRACE_CABAC_T( "\n" )
    } else {
        if !(pcCU.GetMergeFlag1(uiAbsPartIdx) && pcCU.GetPartitionSize1(uiAbsPartIdx) == TLibCommon.SIZE_2Nx2N) {
            this.m_pcEntropyCoderIf.codeQtRootCbf(pcCU, uiAbsPartIdx)
        }
        if !pcCU.GetQtRootCbf(uiAbsPartIdx) {
            return
        }
    }

    this.xEncodeTransform(pcCU, uiLumaOffset, uiChromaOffset, uiAbsPartIdx, uiDepth, uiWidth, uiHeight, 0, bCodeDQP)
}
func (this *TEncEntropy) encodeCoeffNxN(pcCU *TLibCommon.TComDataCU, pcCoeff []TLibCommon.TCoeff, uiAbsPartIdx, uiTrWidth, uiTrHeight, uiDepth uint, eType TLibCommon.TextType) {
    // This is for Transform unit processing. This may be used at mode selection stage for Inter.
    this.m_pcEntropyCoderIf.codeCoeffNxN(pcCU, pcCoeff, uiAbsPartIdx, uiTrWidth, uiTrHeight, uiDepth, eType)
}
func (this *TEncEntropy) estimateBit(pcEstBitsSbac *TLibCommon.EstBitsSbacStruct, width, height int, eTType TLibCommon.TextType) {
    if eTType == TLibCommon.TEXT_LUMA {
        eTType = TLibCommon.TEXT_LUMA
    } else {
        eTType = TLibCommon.TEXT_CHROMA
    }
    this.m_pcEntropyCoderIf.estBit(pcEstBitsSbac, width, height, eTType)
}
func (this *TEncEntropy) encodeSaoOffset(saoLcuParam *TLibCommon.SaoLcuParam, compIdx uint) {
    var uiSymbol uint
    var i int

    uiSymbol = uint(saoLcuParam.TypeIdx) + 1
    if compIdx != 2 {
        this.m_pcEntropyCoderIf.codeSaoTypeIdx(uiSymbol)
    }
    if uiSymbol != 0 {
        if saoLcuParam.TypeIdx < 4 && compIdx != 2 {
            saoLcuParam.SubTypeIdx = saoLcuParam.TypeIdx
        }
        var bitDepth int
        if compIdx != 0 {
            bitDepth = TLibCommon.G_bitDepthC
        } else {
            bitDepth = TLibCommon.G_bitDepthY
        }
        offsetTh := 1 << uint(TLibCommon.MIN(int(bitDepth-5), int(5)).(int))
        if saoLcuParam.TypeIdx == TLibCommon.SAO_BO {
            for i = 0; i < saoLcuParam.Length; i++ {
                var absOffset int
                if saoLcuParam.Offset[i] < 0 {
                    absOffset = -saoLcuParam.Offset[i]
                } else {
                    absOffset = saoLcuParam.Offset[i]
                }
                this.m_pcEntropyCoderIf.codeSaoMaxUvlc(uint(absOffset), uint(offsetTh-1))
            }
            for i = 0; i < saoLcuParam.Length; i++ {
                if saoLcuParam.Offset[i] != 0 {
                    var sign uint
                    if saoLcuParam.Offset[i] < 0 {
                        sign = 1
                    } else {
                        sign = 0
                    }
                    this.m_pcEntropyCoderIf.codeSAOSign(sign)
                }
            }
            uiSymbol = uint(saoLcuParam.SubTypeIdx)
            this.m_pcEntropyCoderIf.codeSaoUflc(5, uiSymbol)
        } else if saoLcuParam.TypeIdx < 4 {
            this.m_pcEntropyCoderIf.codeSaoMaxUvlc(uint(saoLcuParam.Offset[0]), uint(offsetTh-1))
            this.m_pcEntropyCoderIf.codeSaoMaxUvlc(uint(saoLcuParam.Offset[1]), uint(offsetTh-1))
            this.m_pcEntropyCoderIf.codeSaoMaxUvlc(uint(-saoLcuParam.Offset[2]), uint(offsetTh-1))
            this.m_pcEntropyCoderIf.codeSaoMaxUvlc(uint(-saoLcuParam.Offset[3]), uint(offsetTh-1))
            if compIdx != 2 {
                uiSymbol = uint(saoLcuParam.SubTypeIdx)
                this.m_pcEntropyCoderIf.codeSaoUflc(2, uiSymbol)
            }
        }
    }
}
func (this *TEncEntropy) encodeSaoUnitInterleaving(compIdx int, saoFlag bool, rx, ry int, saoLcuParam *TLibCommon.SaoLcuParam, cuAddrInSlice, cuAddrUpInSlice int, allowMergeLeft, allowMergeUp bool) {
    if saoFlag {
        if rx > 0 && cuAddrInSlice != 0 && allowMergeLeft {
            this.m_pcEntropyCoderIf.codeSaoMerge(uint(TLibCommon.B2U(saoLcuParam.MergeLeftFlag)))
        } else {
            saoLcuParam.MergeLeftFlag = false
        }
        if saoLcuParam.MergeLeftFlag == false {
            if (ry > 0) && (cuAddrUpInSlice >= 0) && allowMergeUp {
                this.m_pcEntropyCoderIf.codeSaoMerge(uint(TLibCommon.B2U(saoLcuParam.MergeUpFlag)))
            } else {
                saoLcuParam.MergeUpFlag = false
            }
            if saoLcuParam.MergeUpFlag == false {
                this.encodeSaoOffset(saoLcuParam, uint(compIdx))
            }
        }
    }
}
func TEncEntropy_countNonZeroCoeffs(pcCoef []TLibCommon.TCoeff, uiSize uint) int {
    count := 0

    for i := uint(0); i < uiSize; i++ {
        count += int(TLibCommon.B2U(pcCoef[i] != 0))
    }

    return count
}
