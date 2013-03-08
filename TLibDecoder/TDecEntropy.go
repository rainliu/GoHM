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
    //"fmt"
    "gohm/TLibCommon"
    "io"
)

// ====================================================================================================================
// Class definition
// ====================================================================================================================

/// entropy decoder pure class
type TDecEntropyIf interface {
    //public:
    //  Virtual list for SBAC/CAVLC
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

    /*DTRACE_CABAC_F(x float32)
    DTRACE_CABAC_V(x uint)
    DTRACE_CABAC_VL(x uint)
    DTRACE_CABAC_T(x string)
    DTRACE_CABAC_X(x uint)
    DTRACE_CABAC_N()*/

    ResetEntropy(pcSlice *TLibCommon.TComSlice)
    SetBitstream(p *TLibCommon.TComInputBitstream)
    SetTraceFile(traceFile io.Writer)
    SetSliceTrace(bSliceTrace bool)

    ParseVPS(pcVPS *TLibCommon.TComVPS)
    ParseSPS(pcSPS *TLibCommon.TComSPS)
    ParsePPS(pcPPS *TLibCommon.TComPPS)

    ParseSliceHeader(rpcSlice *TLibCommon.TComSlice, parameterSetManager *TLibCommon.ParameterSetManager) bool

    ParseTerminatingBit(ruilsLast *uint)

    ParseMVPIdx(riMVPIdx *int)

    ParseSkipFlag(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint)
    ParseCUTransquantBypassFlag(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint)
    ParseSplitFlag(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint)
    ParseMergeFlag(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth, uiPUIdx uint)
    ParseMergeIndex(pcCU *TLibCommon.TComDataCU, ruiMergeIndex *uint)
    ParsePartSize(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint)
    ParsePredMode(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint)

    ParseIntraDirLumaAng(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint)

    ParseIntraDirChroma(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint)

    ParseInterDir(pcCU *TLibCommon.TComDataCU, ruiInterDir *uint, uiAbsPartIdx uint)
    ParseRefFrmIdx(pcCU *TLibCommon.TComDataCU, riRefFrmIdx *int, eRefList TLibCommon.RefPicList)
    ParseMvd(pcCU *TLibCommon.TComDataCU, uiAbsPartAddr, uiPartIdx, uiDepth uint, eRefList TLibCommon.RefPicList)

    ParseTransformSubdivFlag(ruiSubdivFlag *uint, uiLog2TransformBlockSize uint)
    ParseQtCbf(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint, eType TLibCommon.TextType, uiTrDepth, uiDepth uint)
    ParseQtRootCbf(uiAbsPartIdx uint, uiQtRootCbf *uint)

    ParseDeltaQP(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint)

    ParseIPCMInfo(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint)

    ParseCoeffNxN(pcCU *TLibCommon.TComDataCU, pcCoef []TLibCommon.TCoeff, uiAbsPartIdx, uiWidth, uiHeight, uiDepth uint, eTType TLibCommon.TextType)
    ParseTransformSkipFlags(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, width, height, uiDepth uint, eTType TLibCommon.TextType)
    UpdateContextTables(eSliceType TLibCommon.SliceType, iQp int)
}

/// entropy decoder class
type TDecEntropy struct {
    //private:
    m_pcEntropyDecoderIf TDecEntropyIf
    m_pcPrediction       *TLibCommon.TComPrediction
    m_uiBakAbsPartIdx    uint
    m_uiBakChromaOffset  uint
    m_bakAbsPartIdxCU    uint
}

func NewTDecEntropy() *TDecEntropy {
    return &TDecEntropy{}
}

func (this *TDecEntropy) Init(p *TLibCommon.TComPrediction) {
    this.m_pcPrediction = p
}

func (this *TDecEntropy) DecodePUWise(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint, pcSubCU *TLibCommon.TComDataCU) {
    ePartSize := pcCU.GetPartitionSize1(uiAbsPartIdx)
    var uiNumPU uint
    if ePartSize == TLibCommon.SIZE_2Nx2N {
        uiNumPU = 1
    } else if ePartSize == TLibCommon.SIZE_NxN {
        uiNumPU = 4
    } else {
        uiNumPU = 2
    }

    uiPUOffset := (TLibCommon.G_auiPUOffset[uint(ePartSize)] << ((pcCU.GetSlice().GetSPS().GetMaxCUDepth() - uiDepth) << 1)) >> 4

    var cMvFieldNeighbours [TLibCommon.MRG_MAX_NUM_CANDS << 1]TLibCommon.TComMvField // double length for mv of both lists
    var uhInterDirNeighbours [TLibCommon.MRG_MAX_NUM_CANDS]byte

    for ui := uint(0); ui < pcCU.GetSlice().GetMaxNumMergeCand(); ui++ {
        uhInterDirNeighbours[ui] = 0
    }
    numValidMergeCand := 0
    isMerged := false

    pcSubCU.CopyInterPredInfoFrom(pcCU, uiAbsPartIdx, TLibCommon.REF_PIC_LIST_0)
    pcSubCU.CopyInterPredInfoFrom(pcCU, uiAbsPartIdx, TLibCommon.REF_PIC_LIST_1)
    uiSubPartIdx := uiAbsPartIdx
    for uiPartIdx := uint(0); uiPartIdx < uiNumPU; uiPartIdx++ {
        this.DecodeMergeFlag(pcCU, uiSubPartIdx, uiDepth, uiPartIdx)
        if pcCU.GetMergeFlag1(uiSubPartIdx) {
            this.DecodeMergeIndex(pcCU, uiPartIdx, uiSubPartIdx, uiDepth)
            uiMergeIndex := pcCU.GetMergeIndex1(uiSubPartIdx)
            if pcCU.GetSlice().GetPPS().GetLog2ParallelMergeLevelMinus2() != 0 && ePartSize != TLibCommon.SIZE_2Nx2N && pcSubCU.GetWidth1(0) <= 8 {
                pcSubCU.SetPartSizeSubParts(TLibCommon.SIZE_2Nx2N, 0, uiDepth)
                if !isMerged {
                    pcSubCU.GetInterMergeCandidates(0, 0, cMvFieldNeighbours[:], uhInterDirNeighbours[:], &numValidMergeCand, -1)
                    isMerged = true
                }
                pcSubCU.SetPartSizeSubParts(ePartSize, 0, uiDepth)
            } else {
                uiMergeIndex = pcCU.GetMergeIndex1(uiSubPartIdx)
                pcSubCU.GetInterMergeCandidates(uiSubPartIdx-uiAbsPartIdx, uiPartIdx, cMvFieldNeighbours[:], uhInterDirNeighbours[:], &numValidMergeCand, int(uiMergeIndex))
            }
            pcCU.SetInterDirSubParts(uint(uhInterDirNeighbours[uiMergeIndex]), uiSubPartIdx, uiPartIdx, uiDepth)

            cTmpMv := TLibCommon.NewTComMv(0, 0)
            for uiRefListIdx := 0; uiRefListIdx < 2; uiRefListIdx++ {
                if pcCU.GetSlice().GetNumRefIdx(TLibCommon.RefPicList(uiRefListIdx)) > 0 {
                    pcCU.SetMVPIdxSubParts(0, TLibCommon.RefPicList(uiRefListIdx), uiSubPartIdx, uiPartIdx, uiDepth)
                    pcCU.SetMVPNumSubParts(0, TLibCommon.RefPicList(uiRefListIdx), uiSubPartIdx, uiPartIdx, uiDepth)
                    pcCU.GetCUMvField(TLibCommon.RefPicList(uiRefListIdx)).SetAllMvd(*cTmpMv, ePartSize, int(uiSubPartIdx), uiDepth, int(uiPartIdx))
                    pcCU.GetCUMvField(TLibCommon.RefPicList(uiRefListIdx)).SetAllMvField(&cMvFieldNeighbours[2*int(uiMergeIndex)+uiRefListIdx], ePartSize, int(uiSubPartIdx), uiDepth, int(uiPartIdx))
                }
            }
        } else {
            this.DecodeInterDirPU(pcCU, uiSubPartIdx, uiDepth, uiPartIdx)
            for uiRefListIdx := 0; uiRefListIdx < 2; uiRefListIdx++ {
                if pcCU.GetSlice().GetNumRefIdx(TLibCommon.RefPicList(uiRefListIdx)) > 0 {
                    //fmt.Printf("%d \n",uiRefListIdx);
                    this.DecodeRefFrmIdxPU(pcCU, uiSubPartIdx, uiDepth, uiPartIdx, TLibCommon.RefPicList(uiRefListIdx))
                    this.DecodeMvdPU(pcCU, uiSubPartIdx, uiDepth, uiPartIdx, TLibCommon.RefPicList(uiRefListIdx))
                    this.DecodeMVPIdxPU(pcSubCU, uiSubPartIdx-uiAbsPartIdx, uiDepth, uiPartIdx, TLibCommon.RefPicList(uiRefListIdx))
                }
            }
        }
        if (pcCU.GetInterDir1(uiSubPartIdx) == 3) && pcSubCU.IsBipredRestriction(uiPartIdx) {
            pcCU.GetCUMvField(TLibCommon.REF_PIC_LIST_1).SetAllMv(*TLibCommon.NewTComMv(0, 0), ePartSize, int(uiSubPartIdx), uiDepth, int(uiPartIdx))
            pcCU.GetCUMvField(TLibCommon.REF_PIC_LIST_1).SetAllRefIdx(-1, ePartSize, int(uiSubPartIdx), uiDepth, int(uiPartIdx))
            pcCU.SetInterDirSubParts(1, uiSubPartIdx, uiPartIdx, uiDepth)
        }

        uiSubPartIdx += uiPUOffset
    }
    return
}
func (this *TDecEntropy) DecodeInterDirPU(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth, uiPartIdx uint) {
    var uiInterDir uint

    if pcCU.GetSlice().IsInterP() {
        uiInterDir = 1
    } else {
        this.m_pcEntropyDecoderIf.ParseInterDir(pcCU, &uiInterDir, uiAbsPartIdx)
    }

    pcCU.SetInterDirSubParts(uiInterDir, uiAbsPartIdx, uiPartIdx, uiDepth)

}

func (this *TDecEntropy) DecodeRefFrmIdxPU(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth, uiPartIdx uint, eRefList TLibCommon.RefPicList) {
    iRefFrmIdx := 0
    iParseRefFrmIdx := pcCU.GetInterDir1(uiAbsPartIdx) & (1 << uint(eRefList))
    //fmt.Printf("iParseRefFrmIdx=%d\n", iParseRefFrmIdx);

    if pcCU.GetSlice().GetNumRefIdx(eRefList) > 1 && iParseRefFrmIdx != 0 {
        this.m_pcEntropyDecoderIf.ParseRefFrmIdx(pcCU, &iRefFrmIdx, eRefList)
        //fmt.Printf("0iRefFrmIdx=%d\n", iRefFrmIdx);
    } else if iParseRefFrmIdx == 0 {
        iRefFrmIdx = TLibCommon.NOT_VALID
        //fmt.Printf("1iRefFrmIdx=%d\n", iRefFrmIdx);
    } else {
        iRefFrmIdx = 0
        //fmt.Printf("2iRefFrmIdx=%d\n", iRefFrmIdx);
    }

    ePartSize := pcCU.GetPartitionSize1(uiAbsPartIdx)
    pcCU.GetCUMvField(eRefList).SetAllRefIdx(int8(iRefFrmIdx), ePartSize, int(uiAbsPartIdx), uiDepth, int(uiPartIdx))
}
func (this *TDecEntropy) DecodeMvdPU(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth, uiPartIdx uint, eRefList TLibCommon.RefPicList) {
    if (pcCU.GetInterDir1(uiAbsPartIdx) & (1 << eRefList)) != 0 {
        this.m_pcEntropyDecoderIf.ParseMvd(pcCU, uiAbsPartIdx, uiPartIdx, uiDepth, eRefList)
    }
}
func (this *TDecEntropy) DecodeMVPIdxPU(pcSubCU *TLibCommon.TComDataCU, uiPartAddr, uiDepth, uiPartIdx uint, eRefList TLibCommon.RefPicList) {
    iMVPIdx := -1

    cZeroMv := TLibCommon.NewTComMv(0, 0)
    cMv := cZeroMv
    iRefIdx := -1

    pcSubCUMvField := pcSubCU.GetCUMvField(eRefList)
    pAMVPInfo := pcSubCUMvField.GetAMVPInfo()

    iRefIdx = int(pcSubCUMvField.GetRefIdx(int(uiPartAddr)))
    cMv = cZeroMv

    if (pcSubCU.GetInterDir1(uiPartAddr) & (1 << eRefList)) != 0 {
        this.m_pcEntropyDecoderIf.ParseMVPIdx(&iMVPIdx)
    }
    pcSubCU.FillMvpCand(uiPartIdx, uiPartAddr, eRefList, iRefIdx, pAMVPInfo)
    pcSubCU.SetMVPNumSubParts(pAMVPInfo.IN, eRefList, uiPartAddr, uiPartIdx, uiDepth)
    pcSubCU.SetMVPIdxSubParts(iMVPIdx, eRefList, uiPartAddr, uiPartIdx, uiDepth)
    if iRefIdx >= 0 {
        cMv = this.m_pcPrediction.GetMvPredAMVP(pcSubCU, uiPartIdx, uiPartAddr, eRefList)
        cMvd := pcSubCUMvField.GetMvd(int(uiPartAddr))
        //fmt.Printf("%d=(%d,%d)=(%d,%d)\n", iRefIdx, cMv.GetHor(), cMv.GetVer(), cMvd.GetHor(), cMvd.GetVer());
        cMv.Set(cMv.GetHor()+cMvd.GetHor(), cMv.GetVer()+cMvd.GetVer())
    }

    ePartSize := pcSubCU.GetPartitionSize1(uiPartAddr)
    pcSubCU.GetCUMvField(eRefList).SetAllMv(*cMv, ePartSize, int(uiPartAddr), 0, int(uiPartIdx))
    //fmt.Printf("%d(%d,%d)=%d,%d,%d",eRefList,cMv.GetHor(), cMv.GetVer(),ePartSize, uiPartAddr, uiPartIdx);
}

func (this *TDecEntropy) SetEntropyDecoder(p TDecEntropyIf) {
    this.m_pcEntropyDecoderIf = p
}
func (this *TDecEntropy) SetBitstream(p *TLibCommon.TComInputBitstream) {
    this.m_pcEntropyDecoderIf.SetBitstream(p)
}
func (this *TDecEntropy) SetTraceFile(traceFile io.Writer) {
    this.m_pcEntropyDecoderIf.SetTraceFile(traceFile)
}
func (this *TDecEntropy) SetSliceTrace(bSliceTrace bool) {
    this.m_pcEntropyDecoderIf.SetSliceTrace(bSliceTrace)
}
func (this *TDecEntropy) ResetEntropy(p *TLibCommon.TComSlice) {
    this.m_pcEntropyDecoderIf.ResetEntropy(p)
}
func (this *TDecEntropy) DecodeVPS(pcVPS *TLibCommon.TComVPS) {
    this.m_pcEntropyDecoderIf.ParseVPS(pcVPS)
}
func (this *TDecEntropy) DecodeSPS(pcSPS *TLibCommon.TComSPS) {
    this.m_pcEntropyDecoderIf.ParseSPS(pcSPS)
}
func (this *TDecEntropy) DecodePPS(pcPPS *TLibCommon.TComPPS) {
    this.m_pcEntropyDecoderIf.ParsePPS(pcPPS)
}
func (this *TDecEntropy) DecodeSliceHeader(rpcSlice *TLibCommon.TComSlice, parameterSetManager *TLibCommon.ParameterSetManager) bool {
    return this.m_pcEntropyDecoderIf.ParseSliceHeader(rpcSlice, parameterSetManager)
}

func (this *TDecEntropy) DecodeTerminatingBit(ruiIsLast *uint) {
    this.m_pcEntropyDecoderIf.ParseTerminatingBit(ruiIsLast)
}

func (this *TDecEntropy) GetEntropyDecoder() TDecEntropyIf {
    return this.m_pcEntropyDecoderIf
}

//public:
func (this *TDecEntropy) DecodeSplitFlag(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint) {
    this.m_pcEntropyDecoderIf.ParseSplitFlag(pcCU, uiAbsPartIdx, uiDepth)
}
func (this *TDecEntropy) DecodeSkipFlag(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint) {
    this.m_pcEntropyDecoderIf.ParseSkipFlag(pcCU, uiAbsPartIdx, uiDepth)
}
func (this *TDecEntropy) DecodeCUTransquantBypassFlag(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint) {
    this.m_pcEntropyDecoderIf.ParseCUTransquantBypassFlag(pcCU, uiAbsPartIdx, uiDepth)
}
func (this *TDecEntropy) DecodeMergeFlag(pcSubCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth, uiPUIdx uint) {
    // at least one merge candidate exists
    this.m_pcEntropyDecoderIf.ParseMergeFlag(pcSubCU, uiAbsPartIdx, uiDepth, uiPUIdx)
}
func (this *TDecEntropy) DecodeMergeIndex(pcCU *TLibCommon.TComDataCU, uiPartIdx, uiAbsPartIdx uint, uiDepth uint) {
    uiMergeIndex := uint(0)
    this.m_pcEntropyDecoderIf.ParseMergeIndex(pcCU, &uiMergeIndex)
    pcCU.SetMergeIndexSubParts(uiMergeIndex, uiAbsPartIdx, uiPartIdx, uiDepth)

}
func (this *TDecEntropy) DecodePredMode(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint) {
    this.m_pcEntropyDecoderIf.ParsePredMode(pcCU, uiAbsPartIdx, uiDepth)
}
func (this *TDecEntropy) DecodePartSize(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint) {
    this.m_pcEntropyDecoderIf.ParsePartSize(pcCU, uiAbsPartIdx, uiDepth)
}

func (this *TDecEntropy) DecodeIPCMInfo(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint) {
    if !pcCU.GetSlice().GetSPS().GetUsePCM() ||
        pcCU.GetWidth1(uiAbsPartIdx) > (1<<pcCU.GetSlice().GetSPS().GetPCMLog2MaxSize()) ||
        pcCU.GetWidth1(uiAbsPartIdx) < (1<<pcCU.GetSlice().GetSPS().GetPCMLog2MinSize()) {
        return
    }

    this.m_pcEntropyDecoderIf.ParseIPCMInfo(pcCU, uiAbsPartIdx, uiDepth)
}

func (this *TDecEntropy) DecodePredInfo(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint, pcSubCU *TLibCommon.TComDataCU) {
    if pcCU.IsIntra(uiAbsPartIdx) {
        this.DecodeIntraDirModeLuma(pcCU, uiAbsPartIdx, uiDepth)
        this.DecodeIntraDirModeChroma(pcCU, uiAbsPartIdx, uiDepth)
    } else {
        this.DecodePUWise(pcCU, uiAbsPartIdx, uiDepth, pcSubCU)
    }
}

func (this *TDecEntropy) DecodeIntraDirModeLuma(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint) {
    this.m_pcEntropyDecoderIf.ParseIntraDirLumaAng(pcCU, uiAbsPartIdx, uiDepth)
}
func (this *TDecEntropy) DecodeIntraDirModeChroma(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint) {
    this.m_pcEntropyDecoderIf.ParseIntraDirChroma(pcCU, uiAbsPartIdx, uiDepth)
}

func (this *TDecEntropy) DecodeQP(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint) {
    if pcCU.GetSlice().GetPPS().GetUseDQP() {
        this.m_pcEntropyDecoderIf.ParseDeltaQP(pcCU, uiAbsPartIdx, uint(pcCU.GetDepth1(uiAbsPartIdx)))
    }
}

func (this *TDecEntropy) UpdateContextTables(eSliceType TLibCommon.SliceType, iQp int) {
    this.m_pcEntropyDecoderIf.UpdateContextTables(eSliceType, iQp)
}
func (this *TDecEntropy) DecodeCoeff(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth, uiWidth, uiHeight uint, bCodeDQP *bool) {
    uiMinCoeffSize := pcCU.GetPic().GetMinCUWidth() * pcCU.GetPic().GetMinCUHeight()
    uiLumaOffset := uiMinCoeffSize * uiAbsPartIdx
    uiChromaOffset := uiLumaOffset >> 2

    if !pcCU.IsIntra(uiAbsPartIdx) {
        uiQtRootCbf := uint(1)
        if !(pcCU.GetPartitionSize1(uiAbsPartIdx) == TLibCommon.SIZE_2Nx2N && pcCU.GetMergeFlag1(uiAbsPartIdx)) {
            this.m_pcEntropyDecoderIf.ParseQtRootCbf(uiAbsPartIdx, &uiQtRootCbf)
        }
        if uiQtRootCbf == 0 {
            pcCU.SetCbfSubParts(0, 0, 0, uiAbsPartIdx, uiDepth)
            pcCU.SetTrIdxSubParts(0, uiAbsPartIdx, uiDepth)
            return
        }

    }
    this.xDecodeTransform(pcCU, uiLumaOffset, uiChromaOffset, uiAbsPartIdx, uiDepth, uiWidth, uiHeight, 0, bCodeDQP)
}

//private:
func (this *TDecEntropy) xDecodeTransform(pcCU *TLibCommon.TComDataCU, offsetLuma, offsetChroma, uiAbsPartIdx, uiDepth, width, height, uiTrIdx uint, bCodeDQP *bool) {
    var uiSubdiv uint
    uiLog2TrafoSize := uint(TLibCommon.G_aucConvertToBit[pcCU.GetSlice().GetSPS().GetMaxCUWidth()]) + 2 - uiDepth

    if uiTrIdx == 0 {
        this.m_bakAbsPartIdxCU = uiAbsPartIdx
    }
    if uiLog2TrafoSize == 2 {
        partNum := pcCU.GetPic().GetNumPartInCU() >> ((uiDepth - 1) << 1)
        if (uiAbsPartIdx % partNum) == 0 {
            this.m_uiBakAbsPartIdx = uiAbsPartIdx
            this.m_uiBakChromaOffset = offsetChroma
        }
    }
    if pcCU.GetPredictionMode1(uiAbsPartIdx) == TLibCommon.MODE_INTRA && pcCU.GetPartitionSize1(uiAbsPartIdx) == TLibCommon.SIZE_NxN && uiDepth == uint(pcCU.GetDepth1(uiAbsPartIdx)) {
        uiSubdiv = 1
    } else if (pcCU.GetSlice().GetSPS().GetQuadtreeTUMaxDepthInter() == 1) && (pcCU.GetPredictionMode1(uiAbsPartIdx) == TLibCommon.MODE_INTER) && (pcCU.GetPartitionSize1(uiAbsPartIdx) != TLibCommon.SIZE_2Nx2N) && (uiDepth == uint(pcCU.GetDepth1(uiAbsPartIdx))) {
        uiSubdiv = uint(TLibCommon.B2U(uiLog2TrafoSize > pcCU.GetQuadtreeTULog2MinSizeInCU(uiAbsPartIdx)))
    } else if uiLog2TrafoSize > pcCU.GetSlice().GetSPS().GetQuadtreeTULog2MaxSize() {
        uiSubdiv = 1
    } else if uiLog2TrafoSize == pcCU.GetSlice().GetSPS().GetQuadtreeTULog2MinSize() {
        uiSubdiv = 0
    } else if uiLog2TrafoSize == pcCU.GetQuadtreeTULog2MinSizeInCU(uiAbsPartIdx) {
        uiSubdiv = 0
    } else {
        //assert( uiLog2TrafoSize > pcCU.GetQuadtreeTULog2MinSizeInCU(uiAbsPartIdx) );
        this.m_pcEntropyDecoderIf.ParseTransformSubdivFlag(&uiSubdiv, 5-uiLog2TrafoSize)
    }

    uiTrDepth := uiDepth - uint(pcCU.GetDepth1(uiAbsPartIdx))
    {
        bFirstCbfOfCU := uiTrDepth == 0
        if bFirstCbfOfCU {
            pcCU.SetCbfSubParts4(0, TLibCommon.TEXT_CHROMA_U, uiAbsPartIdx, uiDepth)
            pcCU.SetCbfSubParts4(0, TLibCommon.TEXT_CHROMA_V, uiAbsPartIdx, uiDepth)
        }
        if bFirstCbfOfCU || uiLog2TrafoSize > 2 {
            if bFirstCbfOfCU || pcCU.GetCbf3(uiAbsPartIdx, TLibCommon.TEXT_CHROMA_U, uiTrDepth-1) != 0 {
                this.m_pcEntropyDecoderIf.ParseQtCbf(pcCU, uiAbsPartIdx, TLibCommon.TEXT_CHROMA_U, uiTrDepth, uiDepth)
            }
            if bFirstCbfOfCU || pcCU.GetCbf3(uiAbsPartIdx, TLibCommon.TEXT_CHROMA_V, uiTrDepth-1) != 0 {
                this.m_pcEntropyDecoderIf.ParseQtCbf(pcCU, uiAbsPartIdx, TLibCommon.TEXT_CHROMA_V, uiTrDepth, uiDepth)
            }
        } else {
            pcCU.SetCbfSubParts4(byte(pcCU.GetCbf3(uiAbsPartIdx, TLibCommon.TEXT_CHROMA_U, uiTrDepth-1)<<uiTrDepth), TLibCommon.TEXT_CHROMA_U, uiAbsPartIdx, uiDepth)
            pcCU.SetCbfSubParts4(byte(pcCU.GetCbf3(uiAbsPartIdx, TLibCommon.TEXT_CHROMA_V, uiTrDepth-1)<<uiTrDepth), TLibCommon.TEXT_CHROMA_V, uiAbsPartIdx, uiDepth)
        }
    }

    if uiSubdiv != 0 {
        var size uint
        width >>= 1
        height >>= 1
        size = width * height
        uiTrIdx++
        uiDepth++
        uiQPartNum := pcCU.GetPic().GetNumPartInCU() >> (uiDepth << 1)
        uiStartAbsPartIdx := uiAbsPartIdx
        uiYCbf := uint(0)
        uiUCbf := uint(0)
        uiVCbf := uint(0)

        for i := uint(0); i < 4; i++ {
            this.xDecodeTransform(pcCU, offsetLuma, offsetChroma, uiAbsPartIdx, uiDepth, width, height, uiTrIdx, bCodeDQP)
            uiYCbf |= uint(pcCU.GetCbf3(uiAbsPartIdx, TLibCommon.TEXT_LUMA, uiTrDepth+1))
            uiUCbf |= uint(pcCU.GetCbf3(uiAbsPartIdx, TLibCommon.TEXT_CHROMA_U, uiTrDepth+1))
            uiVCbf |= uint(pcCU.GetCbf3(uiAbsPartIdx, TLibCommon.TEXT_CHROMA_V, uiTrDepth+1))
            uiAbsPartIdx += uiQPartNum
            offsetLuma += size
            offsetChroma += (size >> 2)
        }

        for ui := uint(0); ui < 4*uiQPartNum; ui++ {
            pcCU.GetCbf1(TLibCommon.TEXT_LUMA)[uiStartAbsPartIdx+ui] |= byte(uiYCbf << uiTrDepth)
            pcCU.GetCbf1(TLibCommon.TEXT_CHROMA_U)[uiStartAbsPartIdx+ui] |= byte(uiUCbf << uiTrDepth)
            pcCU.GetCbf1(TLibCommon.TEXT_CHROMA_V)[uiStartAbsPartIdx+ui] |= byte(uiVCbf << uiTrDepth)
        }
    } else {
        //assert( uiDepth >= pcCU.GetDepth( uiAbsPartIdx ) );
        pcCU.SetTrIdxSubParts(uiTrDepth, uiAbsPartIdx, uiDepth)

        {
            //DTRACE_CABAC_VL( TLibCommon.G_nSymbolCounter++ );
            /*this.m_pcEntropyDecoderIf.DTRACE_CABAC_T("\tTrIdx: abspart=")
            this.m_pcEntropyDecoderIf.DTRACE_CABAC_V(uiAbsPartIdx)
            this.m_pcEntropyDecoderIf.DTRACE_CABAC_T("\tdepth=")
            this.m_pcEntropyDecoderIf.DTRACE_CABAC_V(uiDepth)
            this.m_pcEntropyDecoderIf.DTRACE_CABAC_T("\ttrdepth=")
            this.m_pcEntropyDecoderIf.DTRACE_CABAC_V(uiTrDepth)
            this.m_pcEntropyDecoderIf.DTRACE_CABAC_T("\n")*/
        }

        pcCU.SetCbfSubParts4(0, TLibCommon.TEXT_LUMA, uiAbsPartIdx, uiDepth)
        if pcCU.GetPredictionMode1(uiAbsPartIdx) != TLibCommon.MODE_INTRA && uiDepth == uint(pcCU.GetDepth1(uiAbsPartIdx)) && pcCU.GetCbf3(uiAbsPartIdx, TLibCommon.TEXT_CHROMA_U, 0) == 0 && pcCU.GetCbf3(uiAbsPartIdx, TLibCommon.TEXT_CHROMA_V, 0) == 0 {
            pcCU.SetCbfSubParts4(1<<uiTrDepth, TLibCommon.TEXT_LUMA, uiAbsPartIdx, uiDepth)
        } else {
            this.m_pcEntropyDecoderIf.ParseQtCbf(pcCU, uiAbsPartIdx, TLibCommon.TEXT_LUMA, uiTrDepth, uiDepth)
        }

        // transforthis.m_unit begin
        cbfY := pcCU.GetCbf3(uiAbsPartIdx, TLibCommon.TEXT_LUMA, uiTrIdx)
        cbfU := pcCU.GetCbf3(uiAbsPartIdx, TLibCommon.TEXT_CHROMA_U, uiTrIdx)
        cbfV := pcCU.GetCbf3(uiAbsPartIdx, TLibCommon.TEXT_CHROMA_V, uiTrIdx)
        if uiLog2TrafoSize == 2 {
            partNum := pcCU.GetPic().GetNumPartInCU() >> ((uiDepth - 1) << 1)
            if (uiAbsPartIdx % partNum) == (partNum - 1) {
                cbfU = pcCU.GetCbf3(this.m_uiBakAbsPartIdx, TLibCommon.TEXT_CHROMA_U, uiTrIdx)
                cbfV = pcCU.GetCbf3(this.m_uiBakAbsPartIdx, TLibCommon.TEXT_CHROMA_V, uiTrIdx)
            }
        }
        if cbfY != 0 || cbfU != 0 || cbfV != 0 {
            // dQP: only for LCU
            if pcCU.GetSlice().GetPPS().GetUseDQP() {
                if *bCodeDQP {
                    this.DecodeQP(pcCU, this.m_bakAbsPartIdxCU)
                    *bCodeDQP = false
                }
            }
        }
        if cbfY != 0 {
            trWidth := width
            trHeight := height
            this.m_pcEntropyDecoderIf.ParseCoeffNxN(pcCU, pcCU.GetCoeffY()[offsetLuma:], uiAbsPartIdx, trWidth, trHeight, uiDepth, TLibCommon.TEXT_LUMA)
        }
        if uiLog2TrafoSize > 2 {
            trWidth := width >> 1
            trHeight := height >> 1
            if cbfU != 0 {
                this.m_pcEntropyDecoderIf.ParseCoeffNxN(pcCU, pcCU.GetCoeffCb()[offsetChroma:], uiAbsPartIdx, trWidth, trHeight, uiDepth, TLibCommon.TEXT_CHROMA_U)
            }
            if cbfV != 0 {
                this.m_pcEntropyDecoderIf.ParseCoeffNxN(pcCU, pcCU.GetCoeffCr()[offsetChroma:], uiAbsPartIdx, trWidth, trHeight, uiDepth, TLibCommon.TEXT_CHROMA_V)
            }
        } else {
            partNum := pcCU.GetPic().GetNumPartInCU() >> ((uiDepth - 1) << 1)
            if (uiAbsPartIdx % partNum) == (partNum - 1) {
                trWidth := width
                trHeight := height
                if cbfU != 0 {
                    this.m_pcEntropyDecoderIf.ParseCoeffNxN(pcCU, pcCU.GetCoeffCb()[this.m_uiBakChromaOffset:], this.m_uiBakAbsPartIdx, trWidth, trHeight, uiDepth, TLibCommon.TEXT_CHROMA_U)
                }
                if cbfV != 0 {
                    this.m_pcEntropyDecoderIf.ParseCoeffNxN(pcCU, pcCU.GetCoeffCr()[this.m_uiBakChromaOffset:], this.m_uiBakAbsPartIdx, trWidth, trHeight, uiDepth, TLibCommon.TEXT_CHROMA_V)
                }
            }
        }
        // transform_unit end
    }
}
