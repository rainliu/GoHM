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
    "io"
)

// ====================================================================================================================
// Class definition
// ====================================================================================================================

type TDecBinIf interface {
    Init(pcTComBitstream *TLibCommon.TComInputBitstream)
    Uninit()

    Start()
    Finish()
    Flush()

    DecodeBin(ruiBin *uint, rcCtxModel *TLibCommon.ContextModel)
    DecodeBinEP(ruiBin *uint)
    DecodeBinsEP(ruiBins *uint, numBins int)
    DecodeBinTrm(ruiBin *uint)

    ResetBac()

    DecodePCMAlignBits()
    xReadPCMCode(uiLength uint, ruiCode *uint)

    CopyState(pcTDecBinIf TDecBinIf)
    GetTDecBinCABAC() *TDecBinCabac

    SetSbac(pDecSbac *TDecSbac)
}

type TDecBinCabac struct { //: public TDecBinIf
    m_pTDecSbac *TDecSbac
    //private:
    m_pcTComBitstream *TLibCommon.TComInputBitstream
    m_uiRange         uint
    m_uiValue         uint
    m_bitsNeeded      int
}

func NewTDecBinCabac() *TDecBinCabac {
    return &TDecBinCabac{}
}

func (this *TDecBinCabac) SetSbac(pDecSbac *TDecSbac) {
    this.m_pTDecSbac = pDecSbac
}

func (this *TDecBinCabac) Init(pcTComBitstream *TLibCommon.TComInputBitstream) {
    this.m_pcTComBitstream = pcTComBitstream
}
func (this *TDecBinCabac) Uninit() {
    this.m_pcTComBitstream = nil
}

func (this *TDecBinCabac) Start() {
    //assert( m_pcTComBitstream->getNumBitsUntilByteAligned() == 0 );
    this.m_uiRange = 510
    this.m_bitsNeeded = -8
    this.m_uiValue = this.m_pcTComBitstream.ReadByte() << 8
    this.m_uiValue |= this.m_pcTComBitstream.ReadByte()
}
func (this *TDecBinCabac) Finish() {
    //do nothing
}
func (this *TDecBinCabac) Flush() {
    for this.m_pcTComBitstream.GetNumBitsLeft() > 0 &&
        this.m_pcTComBitstream.GetNumBitsUntilByteAligned() != 0 {
        var uiBits uint
        this.m_pcTComBitstream.Read(1, &uiBits)
    }
    this.Start()
}

func (this *TDecBinCabac) DecodeBin(ruiBin *uint, rcCtxModel *TLibCommon.ContextModel) {
    /*this.DTRACE_CABAC_VL( g_nSymbolCounter++ )*/
    /*this.m_pTDecSbac.DTRACE_CABAC_T("\tDecodeBin()")
      this.m_pTDecSbac.DTRACE_CABAC_T("\tm_uiRange=")
      this.m_pTDecSbac.DTRACE_CABAC_V(this.m_uiRange)
      this.m_pTDecSbac.DTRACE_CABAC_T("\tm_uiValue=")
      this.m_pTDecSbac.DTRACE_CABAC_V(this.m_uiValue)
      //this.m_pTDecSbac.DTRACE_CABAC_T("\tm_bitsNeeded=")
      //this.m_pTDecSbac.DTRACE_CABAC_V(uint(this.m_bitsNeeded))
      this.m_pTDecSbac.DTRACE_CABAC_T("\n")*/

    uiLPS := uint(TLibCommon.TComCABACTables_sm_aucLPSTable[rcCtxModel.GetState()][(this.m_uiRange>>6)-4])
    this.m_uiRange -= uiLPS
    scaledRange := this.m_uiRange << 7

    if this.m_uiValue < scaledRange {
        // MPS path
        *ruiBin = uint(rcCtxModel.GetMps())
        rcCtxModel.UpdateMPS()

        if scaledRange >= (256 << 7) {
            return
        }

        this.m_uiRange = scaledRange >> 6
        this.m_uiValue += this.m_uiValue
        this.m_bitsNeeded++
        if this.m_bitsNeeded == 0 {
            this.m_bitsNeeded = -8
            this.m_uiValue += this.m_pcTComBitstream.ReadByte()
        }
    } else {
        // LPS path
        numBits := TLibCommon.TComCABACTables_sm_aucRenormTable[uiLPS>>3]
        this.m_uiValue = (this.m_uiValue - scaledRange) << numBits
        this.m_uiRange = uiLPS << numBits
        *ruiBin = uint(1 - rcCtxModel.GetMps())
        rcCtxModel.UpdateLPS()

        this.m_bitsNeeded += int(numBits)

        if this.m_bitsNeeded >= 0 {
            this.m_uiValue += this.m_pcTComBitstream.ReadByte() << uint(this.m_bitsNeeded)
            this.m_bitsNeeded -= 8
        }
    }
}
func (this *TDecBinCabac) DecodeBinEP(ruiBin *uint) {
    /*this.DTRACE_CABAC_VL( g_nSymbolCounter++ )*/
    /*this.m_pTDecSbac.DTRACE_CABAC_T("\tDecodeBinEP()")
      this.m_pTDecSbac.DTRACE_CABAC_T("\tm_uiRange=")
      this.m_pTDecSbac.DTRACE_CABAC_V(this.m_uiRange)
      this.m_pTDecSbac.DTRACE_CABAC_T("\tm_uiValue=")
      this.m_pTDecSbac.DTRACE_CABAC_V(this.m_uiValue)
      //this.m_pTDecSbac.DTRACE_CABAC_T("\tm_bitsNeeded=")
      //this.m_pTDecSbac.DTRACE_CABAC_V(uint(this.m_bitsNeeded))
      this.m_pTDecSbac.DTRACE_CABAC_T("\n")*/

    this.m_uiValue += this.m_uiValue
    this.m_bitsNeeded++
    if this.m_bitsNeeded >= 0 {
        this.m_bitsNeeded = -8
        this.m_uiValue += this.m_pcTComBitstream.ReadByte()
    }

    *ruiBin = 0
    scaledRange := this.m_uiRange << 7
    if this.m_uiValue >= scaledRange {
        *ruiBin = 1
        this.m_uiValue -= scaledRange
    }
}
func (this *TDecBinCabac) DecodeBinsEP(ruiBin *uint, numBins int) {
    /*this.DTRACE_CABAC_VL( g_nSymbolCounter++ )*/
    /*this.m_pTDecSbac.DTRACE_CABAC_T("\tDecodeBinsEP()")
      this.m_pTDecSbac.DTRACE_CABAC_T("\tm_uiRange=")
      this.m_pTDecSbac.DTRACE_CABAC_V(this.m_uiRange)
      this.m_pTDecSbac.DTRACE_CABAC_T("\tm_uiValue=")
      this.m_pTDecSbac.DTRACE_CABAC_V(this.m_uiValue)
      this.m_pTDecSbac.DTRACE_CABAC_T("\tm_bitsNeeded=")
      this.m_pTDecSbac.DTRACE_CABAC_V(uint(this.m_bitsNeeded))*/
    //this.m_pTDecSbac.DTRACE_CABAC_T("\n")

    bins := uint(0)

    for numBins > 8 {
        byteTmp := this.m_pcTComBitstream.ReadByte()
        this.m_uiValue = (this.m_uiValue << 8) + (byteTmp << uint(8+this.m_bitsNeeded))
        /*this.m_pTDecSbac.DTRACE_CABAC_T("\tbyteTmp=")
        this.m_pTDecSbac.DTRACE_CABAC_V(byteTmp)
        this.m_pTDecSbac.DTRACE_CABAC_T("\tm_uiValue=")
        this.m_pTDecSbac.DTRACE_CABAC_V(this.m_uiValue)*/

        scaledRange := this.m_uiRange << 15
        for i := 0; i < 8; i++ {
            bins += bins
            scaledRange >>= 1
            if this.m_uiValue >= scaledRange {
                bins++
                this.m_uiValue -= scaledRange
            }
        }
        numBins -= 8
    }

    this.m_bitsNeeded += numBins
    this.m_uiValue <<= uint(numBins)

    if this.m_bitsNeeded >= 0 {
        this.m_uiValue += this.m_pcTComBitstream.ReadByte() << uint(this.m_bitsNeeded)

        //this.m_pTDecSbac.DTRACE_CABAC_T("\tm_uiValue=")
        //this.m_pTDecSbac.DTRACE_CABAC_V(this.m_uiValue)

        this.m_bitsNeeded -= 8
    }

    scaledRange := this.m_uiRange << uint(numBins+7)
    for i := 0; i < numBins; i++ {
        bins += bins
        scaledRange >>= 1
        if this.m_uiValue >= scaledRange {
            bins++
            this.m_uiValue -= scaledRange
        }
    }

    //this.m_pTDecSbac.DTRACE_CABAC_T("\n")

    *ruiBin = bins
}
func (this *TDecBinCabac) DecodeBinTrm(ruiBin *uint) {
    this.m_uiRange -= 2
    scaledRange := this.m_uiRange << 7
    if this.m_uiValue >= scaledRange {
        *ruiBin = 1
    } else {
        *ruiBin = 0
        if scaledRange < (256 << 7) {
            this.m_uiRange = scaledRange >> 6
            this.m_uiValue += this.m_uiValue
            this.m_bitsNeeded++
            if this.m_bitsNeeded == 0 {
                this.m_bitsNeeded = -8
                this.m_uiValue += this.m_pcTComBitstream.ReadByte()
            }
        }
    }
}

func (this *TDecBinCabac) ResetBac() {
    this.m_uiRange = 510
    this.m_bitsNeeded = -8
    this.m_uiValue = this.m_pcTComBitstream.ReadBits(16)
}

func (this *TDecBinCabac) DecodePCMAlignBits() {
    iNum := this.m_pcTComBitstream.GetNumBitsUntilByteAligned()

    uiBit := uint(0)
    this.m_pcTComBitstream.Read(iNum, &uiBit)
}
func (this *TDecBinCabac) xReadPCMCode(uiLength uint, ruiCode *uint) {
    //assert ( uiLength > 0 );
    this.m_pcTComBitstream.Read(uiLength, ruiCode)
}

func (this *TDecBinCabac) CopyState(pcTDecBinIf TDecBinIf) {
    pcTDecBinCABAC := pcTDecBinIf.GetTDecBinCABAC()
    this.m_uiRange = pcTDecBinCABAC.m_uiRange
    this.m_uiValue = pcTDecBinCABAC.m_uiValue
    this.m_bitsNeeded = pcTDecBinCABAC.m_bitsNeeded
}
func (this *TDecBinCabac) GetTDecBinCABAC() *TDecBinCabac {
    return this
}

//class SEImessages;

/// SBAC decoder class
type TDecSbac struct { //: public TDecEntropyIf
    //private:
    m_pTraceFile  io.Writer
    m_pcBitstream *TLibCommon.TComInputBitstream
    m_pcTDecBinIf TDecBinIf

    //private:
    m_uiLastDQpNonZero uint
    m_uiLastQp         uint

    m_contextModels    [TLibCommon.MAX_NUM_CTX_MOD]TLibCommon.ContextModel
    m_numContextModels int

    m_cCUSplitFlagSCModel       *TLibCommon.ContextModel3DBuffer
    m_cCUSkipFlagSCModel        *TLibCommon.ContextModel3DBuffer
    m_cCUMergeFlagExtSCModel    *TLibCommon.ContextModel3DBuffer
    m_cCUMergeIdxExtSCModel     *TLibCommon.ContextModel3DBuffer
    m_cCUPartSizeSCModel        *TLibCommon.ContextModel3DBuffer
    m_cCUPredModeSCModel        *TLibCommon.ContextModel3DBuffer
    m_cCUIntraPredSCModel       *TLibCommon.ContextModel3DBuffer
    m_cCUChromaPredSCModel      *TLibCommon.ContextModel3DBuffer
    m_cCUDeltaQpSCModel         *TLibCommon.ContextModel3DBuffer
    m_cCUInterDirSCModel        *TLibCommon.ContextModel3DBuffer
    m_cCURefPicSCModel          *TLibCommon.ContextModel3DBuffer
    m_cCUMvdSCModel             *TLibCommon.ContextModel3DBuffer
    m_cCUQtCbfSCModel           *TLibCommon.ContextModel3DBuffer
    m_cCUTransSubdivFlagSCModel *TLibCommon.ContextModel3DBuffer
    m_cCUQtRootCbfSCModel       *TLibCommon.ContextModel3DBuffer

    m_cCUSigCoeffGroupSCModel *TLibCommon.ContextModel3DBuffer
    m_cCUSigSCModel           *TLibCommon.ContextModel3DBuffer
    m_cCuCtxLastX             *TLibCommon.ContextModel3DBuffer
    m_cCuCtxLastY             *TLibCommon.ContextModel3DBuffer
    m_cCUOneSCModel           *TLibCommon.ContextModel3DBuffer
    m_cCUAbsSCModel           *TLibCommon.ContextModel3DBuffer

    m_cMVPIdxSCModel *TLibCommon.ContextModel3DBuffer

    m_cCUAMPSCModel                 *TLibCommon.ContextModel3DBuffer
    m_cSaoMergeSCModel              *TLibCommon.ContextModel3DBuffer
    m_cSaoTypeIdxSCModel            *TLibCommon.ContextModel3DBuffer
    m_cTransformSkipSCModel         *TLibCommon.ContextModel3DBuffer
    m_CUTransquantBypassFlagSCModel *TLibCommon.ContextModel3DBuffer
}

func (this *TDecSbac) XTraceLCUHeader(traceLevel uint) {
    if this.GetTraceFile() != nil && (traceLevel&TLibCommon.TRACE_LEVEL) != 0 {
        io.WriteString(this.m_pTraceFile, "========= LCU Parameter Set ===============================================\n") //, pLCU.GetAddr());
    }
}

func (this *TDecSbac) XTraceCUHeader(traceLevel uint) {
    if this.GetTraceFile() != nil && (traceLevel&TLibCommon.TRACE_LEVEL) != 0 {
        io.WriteString(this.m_pTraceFile, "========= CU Parameter Set ================================================\n") //, pCU.GetCUPelX(), pCU.GetCUPelY());
    }
}

func (this *TDecSbac) XTracePUHeader(traceLevel uint) {
    if this.GetTraceFile() != nil && (traceLevel&TLibCommon.TRACE_LEVEL) != 0 {
        io.WriteString(this.m_pTraceFile, "========= PU Parameter Set ================================================\n") //, pCU.GetCUPelX(), pCU.GetCUPelY());
    }
}

func (this *TDecSbac) XTraceTUHeader(traceLevel uint) {
    if this.GetTraceFile() != nil && (traceLevel&TLibCommon.TRACE_LEVEL) != 0 {
        io.WriteString(this.m_pTraceFile, "========= TU Parameter Set ================================================\n") //, pCU.GetCUPelX(), pCU.GetCUPelY());
    }
}

func (this *TDecSbac) XTraceCoefHeader(traceLevel uint) {
    if this.GetTraceFile() != nil && (traceLevel&TLibCommon.TRACE_LEVEL) != 0 {
        io.WriteString(this.m_pTraceFile, "========= Coefficient Parameter Set =======================================\n") //, pCU.GetCUPelX(), pCU.GetCUPelY());
    }
}

func (this *TDecSbac) XTraceResiHeader(traceLevel uint) {
    if this.GetTraceFile() != nil && (traceLevel&TLibCommon.TRACE_LEVEL) != 0 {
        io.WriteString(this.m_pTraceFile, "========= Residual Parameter Set ==========================================\n") //, pCU.GetCUPelX(), pCU.GetCUPelY());
    }
}

func (this *TDecSbac) XTracePredHeader(traceLevel uint) {
    if this.GetTraceFile() != nil && (traceLevel&TLibCommon.TRACE_LEVEL) != 0 {
        io.WriteString(this.m_pTraceFile, "========= Prediction Parameter Set ========================================\n") //, pCU.GetCUPelX(), pCU.GetCUPelY());
    }
}

func (this *TDecSbac) XTraceRecoHeader(traceLevel uint) {
    if this.GetTraceFile() != nil && (traceLevel&TLibCommon.TRACE_LEVEL) != 0 {
        io.WriteString(this.m_pTraceFile, "========= Reconstruction Parameter Set ====================================\n") //, pCU.GetCUPelX(), pCU.GetCUPelY());
    }
}

func (this *TDecSbac) XReadAeTr(Value int, pSymbolName string, traceLevel uint) {
    if this.GetTraceFile() != nil && (traceLevel&TLibCommon.TRACE_LEVEL) != 0 {
        //fprintf( g_hTrace, "%8lld  ", g_nSymbolCounter++ );
        io.WriteString(this.m_pTraceFile, fmt.Sprintf("%-62s ae(v) : %4d\n", pSymbolName, Value))
        //fflush ( g_hTrace );
    }
}

func (this *TDecSbac) XReadCeofTr(pCoeff []TLibCommon.TCoeff, uiWidth, traceLevel uint) {
    if this.GetTraceFile() != nil && (traceLevel&TLibCommon.TRACE_LEVEL) != 0 {
    	if TLibCommon.G_uiPicNo==70 {
    	for i := uint(0); i < uiWidth; i++ {
            io.WriteString(this.m_pTraceFile, fmt.Sprintf("%04x ", uint16(pCoeff[i])))
            //if uiWidth==4 {
      		//fmt.Printf("%8d ",pCoeff[i]);
      		//}
        }
        io.WriteString(this.m_pTraceFile, "\n")
        	//if uiWidth==4 {
      		//fmt.Printf("\n");
      		//}
        }
    }
}

func (this *TDecSbac) XReadResiTr(pPel []TLibCommon.Pel, uiWidth, traceLevel uint) {
    if this.GetTraceFile() != nil && (traceLevel&TLibCommon.TRACE_LEVEL) != 0 {
        if TLibCommon.G_uiPicNo==70 {
        for i := uint(0); i < uiWidth; i++ {
            io.WriteString(this.m_pTraceFile, fmt.Sprintf("%04x ", uint16(pPel[i])))
            /*if	uiWidth==4 {
      			fmt.Printf("%4d ",pPel[i]);
      			}*/
        }
        io.WriteString(this.m_pTraceFile, "\n")
        	/*if uiWidth==4 {
      		fmt.Printf("\n");
      		}*/
        }
    }
}

func (this *TDecSbac) XReadPredTr(pPel []TLibCommon.Pel, uiWidth, traceLevel uint) {
    if this.GetTraceFile() != nil && (traceLevel&TLibCommon.TRACE_LEVEL) != 0 {
        for i := uint(0); i < uiWidth; i++ {
            io.WriteString(this.m_pTraceFile, fmt.Sprintf("%02x ", TLibCommon.Pxl(pPel[i])))
        }
        io.WriteString(this.m_pTraceFile, "\n")
    }
}

func (this *TDecSbac) XReadRecoTr(pPel []TLibCommon.Pel, uiWidth, traceLevel uint) {
    if this.GetTraceFile() != nil && (traceLevel&TLibCommon.TRACE_LEVEL) != 0 {
        for i := uint(0); i < uiWidth; i++ {
            io.WriteString(this.m_pTraceFile, fmt.Sprintf("%02x ", TLibCommon.Pxl(pPel[i])))
        }
        io.WriteString(this.m_pTraceFile, "\n")
    }
}
/*
func (this *TDecSbac) DTRACE_CABAC_F(x float32) {
    if this.GetTraceFile() != nil && TLibCommon.TRACE_CABAC {
        //fmt.Printf("%f", x)
        io.WriteString(this.m_pTraceFile, fmt.Sprintf("%f", x))
    }
}
func (this *TDecSbac) DTRACE_CABAC_V(x uint) {
    if this.GetTraceFile() != nil && TLibCommon.TRACE_CABAC {
        //fmt.Printf ("%d", x )
        io.WriteString(this.m_pTraceFile, fmt.Sprintf("%x", x))
    }
}
func (this *TDecSbac) DTRACE_CABAC_VL(x uint) {
    if this.GetTraceFile() != nil && TLibCommon.TRACE_CABAC {
        //fmt.Printf ("%lld", x )
        io.WriteString(this.m_pTraceFile, fmt.Sprintf("%lld", x))
    }
}
func (this *TDecSbac) DTRACE_CABAC_T(x string) {
    if this.GetTraceFile() != nil && TLibCommon.TRACE_CABAC {
        //fmt.Printf ("%s", x )
        io.WriteString(this.m_pTraceFile, fmt.Sprintf("%s", x))
    }
}
func (this *TDecSbac) DTRACE_CABAC_X(x uint) {
    if this.GetTraceFile() != nil && TLibCommon.TRACE_CABAC {
        //fmt.Printf ("%x", x )
        io.WriteString(this.m_pTraceFile, fmt.Sprintf("%x", x))
    }
}
func (this *TDecSbac) DTRACE_CABAC_N() {
    if this.GetTraceFile() != nil && TLibCommon.TRACE_CABAC {
        //fmt.Printf ("\n" )
        io.WriteString(this.m_pTraceFile, "\n")
    }
}
*/
/*func (this *TDecSbac) DTRACE_CABAC_R(x,y) {
	if this.GetTraceFile()!=nil {
		io.WriteString(this.m_pTraceFile, fmt.Sprintf (x,    y ));
	}
}*/

func NewTDecSbac() *TDecSbac {
    pTDecSbac := &TDecSbac{m_pcBitstream: nil, m_pcTDecBinIf: nil, m_numContextModels: 0}
    pTDecSbac.xInit()

    return pTDecSbac
}

func (this *TDecSbac) xInit() {
    this.m_cCUSplitFlagSCModel = TLibCommon.NewContextModel3DBuffer(1, 1, TLibCommon.NUM_SPLIT_FLAG_CTX, this.m_contextModels[this.m_numContextModels:], &this.m_numContextModels)
    this.m_cCUSkipFlagSCModel = TLibCommon.NewContextModel3DBuffer(1, 1, TLibCommon.NUM_SKIP_FLAG_CTX, this.m_contextModels[this.m_numContextModels:], &this.m_numContextModels)
    this.m_cCUMergeFlagExtSCModel = TLibCommon.NewContextModel3DBuffer(1, 1, TLibCommon.NUM_MERGE_FLAG_EXT_CTX, this.m_contextModels[this.m_numContextModels:], &this.m_numContextModels)
    this.m_cCUMergeIdxExtSCModel = TLibCommon.NewContextModel3DBuffer(1, 1, TLibCommon.NUM_MERGE_IDX_EXT_CTX, this.m_contextModels[this.m_numContextModels:], &this.m_numContextModels)
    this.m_cCUPartSizeSCModel = TLibCommon.NewContextModel3DBuffer(1, 1, TLibCommon.NUM_PART_SIZE_CTX, this.m_contextModels[this.m_numContextModels:], &this.m_numContextModels)
    this.m_cCUPredModeSCModel = TLibCommon.NewContextModel3DBuffer(1, 1, TLibCommon.NUM_PRED_MODE_CTX, this.m_contextModels[this.m_numContextModels:], &this.m_numContextModels)
    this.m_cCUIntraPredSCModel = TLibCommon.NewContextModel3DBuffer(1, 1, TLibCommon.NUM_ADI_CTX, this.m_contextModels[this.m_numContextModels:], &this.m_numContextModels)
    this.m_cCUChromaPredSCModel = TLibCommon.NewContextModel3DBuffer(1, 1, TLibCommon.NUM_CHROMA_PRED_CTX, this.m_contextModels[this.m_numContextModels:], &this.m_numContextModels)
    this.m_cCUDeltaQpSCModel = TLibCommon.NewContextModel3DBuffer(1, 1, TLibCommon.NUM_DELTA_QP_CTX, this.m_contextModels[this.m_numContextModels:], &this.m_numContextModels)
    this.m_cCUInterDirSCModel = TLibCommon.NewContextModel3DBuffer(1, 1, TLibCommon.NUM_INTER_DIR_CTX, this.m_contextModels[this.m_numContextModels:], &this.m_numContextModels)
    this.m_cCURefPicSCModel = TLibCommon.NewContextModel3DBuffer(1, 1, TLibCommon.NUM_REF_NO_CTX, this.m_contextModels[this.m_numContextModels:], &this.m_numContextModels)
    this.m_cCUMvdSCModel = TLibCommon.NewContextModel3DBuffer(1, 1, TLibCommon.NUM_MV_RES_CTX, this.m_contextModels[this.m_numContextModels:], &this.m_numContextModels)
    this.m_cCUQtCbfSCModel = TLibCommon.NewContextModel3DBuffer(1, 2, TLibCommon.NUM_QT_CBF_CTX, this.m_contextModels[this.m_numContextModels:], &this.m_numContextModels)
    this.m_cCUTransSubdivFlagSCModel = TLibCommon.NewContextModel3DBuffer(1, 1, TLibCommon.NUM_TRANS_SUBDIV_FLAG_CTX, this.m_contextModels[this.m_numContextModels:], &this.m_numContextModels)
    this.m_cCUQtRootCbfSCModel = TLibCommon.NewContextModel3DBuffer(1, 1, TLibCommon.NUM_QT_ROOT_CBF_CTX, this.m_contextModels[this.m_numContextModels:], &this.m_numContextModels)
    this.m_cCUSigCoeffGroupSCModel = TLibCommon.NewContextModel3DBuffer(1, 2, TLibCommon.NUM_SIG_CG_FLAG_CTX, this.m_contextModels[this.m_numContextModels:], &this.m_numContextModels)
    this.m_cCUSigSCModel = TLibCommon.NewContextModel3DBuffer(1, 1, TLibCommon.NUM_SIG_FLAG_CTX, this.m_contextModels[this.m_numContextModels:], &this.m_numContextModels)
    this.m_cCuCtxLastX = TLibCommon.NewContextModel3DBuffer(1, 2, TLibCommon.NUM_CTX_LAST_FLAG_XY, this.m_contextModels[this.m_numContextModels:], &this.m_numContextModels)
    this.m_cCuCtxLastY = TLibCommon.NewContextModel3DBuffer(1, 2, TLibCommon.NUM_CTX_LAST_FLAG_XY, this.m_contextModels[this.m_numContextModels:], &this.m_numContextModels)
    this.m_cCUOneSCModel = TLibCommon.NewContextModel3DBuffer(1, 1, TLibCommon.NUM_ONE_FLAG_CTX, this.m_contextModels[this.m_numContextModels:], &this.m_numContextModels)
    this.m_cCUAbsSCModel = TLibCommon.NewContextModel3DBuffer(1, 1, TLibCommon.NUM_ABS_FLAG_CTX, this.m_contextModels[this.m_numContextModels:], &this.m_numContextModels)
    this.m_cMVPIdxSCModel = TLibCommon.NewContextModel3DBuffer(1, 1, TLibCommon.NUM_MVP_IDX_CTX, this.m_contextModels[this.m_numContextModels:], &this.m_numContextModels)
    this.m_cCUAMPSCModel = TLibCommon.NewContextModel3DBuffer(1, 1, TLibCommon.NUM_CU_AMP_CTX, this.m_contextModels[this.m_numContextModels:], &this.m_numContextModels)
    this.m_cSaoMergeSCModel = TLibCommon.NewContextModel3DBuffer(1, 1, TLibCommon.NUM_SAO_MERGE_FLAG_CTX, this.m_contextModels[this.m_numContextModels:], &this.m_numContextModels)
    this.m_cSaoTypeIdxSCModel = TLibCommon.NewContextModel3DBuffer(1, 1, TLibCommon.NUM_SAO_TYPE_IDX_CTX, this.m_contextModels[this.m_numContextModels:], &this.m_numContextModels)
    this.m_cTransformSkipSCModel = TLibCommon.NewContextModel3DBuffer(1, 2, TLibCommon.NUM_TRANSFORMSKIP_FLAG_CTX, this.m_contextModels[this.m_numContextModels:], &this.m_numContextModels)
    this.m_CUTransquantBypassFlagSCModel = TLibCommon.NewContextModel3DBuffer(1, 1, TLibCommon.NUM_CU_TRANSQUANT_BYPASS_FLAG_CTX, this.m_contextModels[this.m_numContextModels:], &this.m_numContextModels)
}

func (this *TDecSbac) Init(p TDecBinIf) {
    this.m_pcTDecBinIf = p
    this.m_pcTDecBinIf.SetSbac(this)
}

func (this *TDecSbac) Uninit() {
    this.m_pcTDecBinIf = nil
}

func (this *TDecSbac) Load(pScr *TDecSbac) {
    this.xCopyFrom(pScr)
}
func (this *TDecSbac) LoadContexts(pScr *TDecSbac) {
    this.xCopyContextsFrom(pScr)
}
func (this *TDecSbac) xCopyFrom(pSrc *TDecSbac) {
    this.m_pcTDecBinIf.CopyState(pSrc.m_pcTDecBinIf)

    this.m_uiLastQp = pSrc.m_uiLastQp
    this.xCopyContextsFrom(pSrc)
}
func (this *TDecSbac) xCopyContextsFrom(pSrc *TDecSbac) {
    for i := 0; i < this.m_numContextModels; i++ {
        this.m_contextModels[i] = pSrc.m_contextModels[i] //, m_numContextModels*sizeof(m_contextModels[0]));
    }
}

func (this *TDecSbac) ResetEntropy(pSlice *TLibCommon.TComSlice) {
    sliceType := pSlice.GetSliceType()
    qp := pSlice.GetSliceQp()

    if pSlice.GetPPS().GetCabacInitPresentFlag() && pSlice.GetCabacInitFlag() {
        switch sliceType {
        case TLibCommon.P_SLICE: // change initialization table to B_SLICE initialization
            sliceType = TLibCommon.B_SLICE
            //break;
        case TLibCommon.B_SLICE: // change initialization table to P_SLICE initialization
            sliceType = TLibCommon.P_SLICE
            //break;
            //default     :           // should not occur
            //assert(0);
        }
    }
    
    this.m_cCUSplitFlagSCModel.InitBuffer(sliceType, qp, TLibCommon.INIT_SPLIT_FLAG[:])
    this.m_cCUSkipFlagSCModel.InitBuffer(sliceType, qp, TLibCommon.INIT_SKIP_FLAG[:])
    this.m_cCUMergeFlagExtSCModel.InitBuffer(sliceType, qp, TLibCommon.INIT_MERGE_FLAG_EXT[:])
    this.m_cCUMergeIdxExtSCModel.InitBuffer(sliceType, qp, TLibCommon.INIT_MERGE_IDX_EXT[:])
    this.m_cCUPartSizeSCModel.InitBuffer(sliceType, qp, TLibCommon.INIT_PART_SIZE[:])
    this.m_cCUAMPSCModel.InitBuffer(sliceType, qp, TLibCommon.INIT_CU_AMP_POS[:])
    this.m_cCUPredModeSCModel.InitBuffer(sliceType, qp, TLibCommon.INIT_PRED_MODE[:])
    this.m_cCUIntraPredSCModel.InitBuffer(sliceType, qp, TLibCommon.INIT_INTRA_PRED_MODE[:])
    this.m_cCUChromaPredSCModel.InitBuffer(sliceType, qp, TLibCommon.INIT_CHROMA_PRED_MODE[:])
    this.m_cCUInterDirSCModel.InitBuffer(sliceType, qp, TLibCommon.INIT_INTER_DIR[:])
    this.m_cCUMvdSCModel.InitBuffer(sliceType, qp, TLibCommon.INIT_MVD[:])
    this.m_cCURefPicSCModel.InitBuffer(sliceType, qp, TLibCommon.INIT_REF_PIC[:])
    this.m_cCUDeltaQpSCModel.InitBuffer(sliceType, qp, TLibCommon.INIT_DQP[:])
    this.m_cCUQtCbfSCModel.InitBuffer(sliceType, qp, TLibCommon.INIT_QT_CBF[:])
    this.m_cCUQtRootCbfSCModel.InitBuffer(sliceType, qp, TLibCommon.INIT_QT_ROOT_CBF[:])
    this.m_cCUSigCoeffGroupSCModel.InitBuffer(sliceType, qp, TLibCommon.INIT_SIG_CG_FLAG[:])
    this.m_cCUSigSCModel.InitBuffer(sliceType, qp, TLibCommon.INIT_SIG_FLAG[:])
    this.m_cCuCtxLastX.InitBuffer(sliceType, qp, TLibCommon.INIT_LAST[:])
    this.m_cCuCtxLastY.InitBuffer(sliceType, qp, TLibCommon.INIT_LAST[:])
    this.m_cCUOneSCModel.InitBuffer(sliceType, qp, TLibCommon.INIT_ONE_FLAG[:])
    this.m_cCUAbsSCModel.InitBuffer(sliceType, qp, TLibCommon.INIT_ABS_FLAG[:])
    this.m_cMVPIdxSCModel.InitBuffer(sliceType, qp, TLibCommon.INIT_MVP_IDX[:])
    this.m_cSaoMergeSCModel.InitBuffer(sliceType, qp, TLibCommon.INIT_SAO_MERGE_FLAG[:])
    this.m_cSaoTypeIdxSCModel.InitBuffer(sliceType, qp, TLibCommon.INIT_SAO_TYPE_IDX[:])
    this.m_cCUTransSubdivFlagSCModel.InitBuffer(sliceType, qp, TLibCommon.INIT_TRANS_SUBDIV_FLAG[:])
    this.m_cTransformSkipSCModel.InitBuffer(sliceType, qp, TLibCommon.INIT_TRANSFORMSKIP_FLAG[:])
    this.m_CUTransquantBypassFlagSCModel.InitBuffer(sliceType, qp, TLibCommon.INIT_CU_TRANSQUANT_BYPASS_FLAG[:])

    this.m_uiLastDQpNonZero = 0

    // new structure
    this.m_uiLastQp = uint(qp)

    this.m_pcTDecBinIf.Start()
}
func (this *TDecSbac) SetBitstream(p *TLibCommon.TComInputBitstream) {
    this.m_pcBitstream = p
    this.m_pcTDecBinIf.Init(p)
}
func (this *TDecSbac) SetTraceFile(traceFile io.Writer) {
    this.m_pTraceFile = traceFile
}

func (this *TDecSbac) GetTraceFile() io.Writer {
    return this.m_pTraceFile
}

func (this *TDecSbac) SetSliceTrace(bSliceTrace bool) {
    //do nothing
}
func (this *TDecSbac) ParseVPS(pcVPS *TLibCommon.TComVPS) {
    //do nothing
}
func (this *TDecSbac) ParseSPS(pcSPS *TLibCommon.TComSPS) {
    //do nothing
}
func (this *TDecSbac) ParsePPS(pcPPS *TLibCommon.TComPPS) {
    //do nothing
}

func (this *TDecSbac) ParseSliceHeader(rpcSlice *TLibCommon.TComSlice, parameterSetManager *TLibCommon.ParameterSetManager) bool {
    //do nothing
    return false;
}
func (this *TDecSbac) ParseTerminatingBit(ruiBit *uint) {
    this.m_pcTDecBinIf.DecodeBinTrm(ruiBit)
}
func (this *TDecSbac) ParseMVPIdx(riMVPIdx *int) {
    var uiSymbol uint
    this.xReadUnaryMaxSymbol(&uiSymbol, this.m_cMVPIdxSCModel.Get1(0), 1, TLibCommon.AMVP_MAX_NUM_CANDS-1)
    *riMVPIdx = int(uiSymbol)
}
func (this *TDecSbac) ParseSaoMaxUvlc(val *uint, maxSymbol uint) {
    if maxSymbol == 0 {
        *val = 0
        return
    }

    var code uint
    var i uint
    this.m_pcTDecBinIf.DecodeBinEP(&code)
    if code == 0 {
        *val = 0
        return
    }

    i = 1
    for {
        this.m_pcTDecBinIf.DecodeBinEP(&code)
        if code == 0 {
            break
        }
        i++
        if i == maxSymbol {
            break
        }
    }

    *val = i
}
func (this *TDecSbac) ParseSaoMerge(ruiVal *uint) {
    var uiCode uint
    this.m_pcTDecBinIf.DecodeBin(&uiCode, this.m_cSaoMergeSCModel.Get3(0, 0, 0))
    *ruiVal = uiCode
}
func (this *TDecSbac) ParseSaoTypeIdx(ruiVal *uint) {
    var uiCode uint
    this.m_pcTDecBinIf.DecodeBin(&uiCode, this.m_cSaoTypeIdxSCModel.Get3(0, 0, 0))
    if uiCode == 0 {
        *ruiVal = 0
    } else {
        this.m_pcTDecBinIf.DecodeBinEP(&uiCode)
        if uiCode == 0 {
            *ruiVal = 5
        } else {
            *ruiVal = 1
        }
    }
}
func (this *TDecSbac) ParseSaoUflc(uiLength uint, ruiVal *uint) {
    this.m_pcTDecBinIf.DecodeBinsEP(ruiVal, int(uiLength))
}

func (this *TDecSbac) CopySaoOneLcuParam(psDst *TLibCommon.SaoLcuParam, psSrc *TLibCommon.SaoLcuParam) {
    var i int
    psDst.PartIdx = psSrc.PartIdx
    psDst.TypeIdx = psSrc.TypeIdx
    if psDst.TypeIdx != -1 {
        psDst.SubTypeIdx = psSrc.SubTypeIdx
        psDst.Length = psSrc.Length
        for i = 0; i < psDst.Length; i++ {
            psDst.Offset[i] = psSrc.Offset[i]
        }
    } else {
        psDst.Length = 0
        for i = 0; i < TLibCommon.SAO_BO_LEN; i++ {
            psDst.Offset[i] = 0
        }
    }
}

func (this *TDecSbac) ParseSaoOneLcuInterleaving(rx, ry int, pSaoParam *TLibCommon.SAOParam, pcCU *TLibCommon.TComDataCU, iCUAddrInSlice, iCUAddrUpInSlice int, allowMergeLeft, allowMergeUp bool) {
    iAddr := int(pcCU.GetAddr())
    var uiSymbol uint
    for iCompIdx := 0; iCompIdx < 3; iCompIdx++ {
        pSaoParam.SaoLcuParam[iCompIdx][iAddr].MergeUpFlag = false
        pSaoParam.SaoLcuParam[iCompIdx][iAddr].MergeLeftFlag = false
        pSaoParam.SaoLcuParam[iCompIdx][iAddr].SubTypeIdx = 0
        pSaoParam.SaoLcuParam[iCompIdx][iAddr].TypeIdx = -1
        pSaoParam.SaoLcuParam[iCompIdx][iAddr].Offset[0] = 0
        pSaoParam.SaoLcuParam[iCompIdx][iAddr].Offset[1] = 0
        pSaoParam.SaoLcuParam[iCompIdx][iAddr].Offset[2] = 0
        pSaoParam.SaoLcuParam[iCompIdx][iAddr].Offset[3] = 0

    }
    if pSaoParam.SaoFlag[0] || pSaoParam.SaoFlag[1] {
        if rx > 0 && iCUAddrInSlice != 0 && allowMergeLeft {
            this.ParseSaoMerge(&uiSymbol)
            pSaoParam.SaoLcuParam[0][iAddr].MergeLeftFlag = uiSymbol != 0
            //#ifdef ENC_DEC_TRACE
            this.XReadAeTr(int(uiSymbol), "sao_merge_left_flag", TLibCommon.TRACE_LCU)
            //#endif
        }
        if pSaoParam.SaoLcuParam[0][iAddr].MergeLeftFlag == false {
            if (ry > 0) && (iCUAddrUpInSlice >= 0) && allowMergeUp {
                this.ParseSaoMerge(&uiSymbol)
                pSaoParam.SaoLcuParam[0][iAddr].MergeUpFlag = uiSymbol != 0
                //#ifdef ENC_DEC_TRACE
                this.XReadAeTr(int(uiSymbol), "sao_merge_up_flag", TLibCommon.TRACE_LCU)
                //#endif
            }
        }
    }

    for iCompIdx := 0; iCompIdx < 3; iCompIdx++ {
        if (iCompIdx == 0 && pSaoParam.SaoFlag[0]) || (iCompIdx > 0 && pSaoParam.SaoFlag[1]) {
            if rx > 0 && iCUAddrInSlice != 0 && allowMergeLeft {
                pSaoParam.SaoLcuParam[iCompIdx][iAddr].MergeLeftFlag = pSaoParam.SaoLcuParam[0][iAddr].MergeLeftFlag
            } else {
                pSaoParam.SaoLcuParam[iCompIdx][iAddr].MergeLeftFlag = false
            }

            if pSaoParam.SaoLcuParam[iCompIdx][iAddr].MergeLeftFlag == false {
                if (ry > 0) && (iCUAddrUpInSlice >= 0) && allowMergeUp {
                    pSaoParam.SaoLcuParam[iCompIdx][iAddr].MergeUpFlag = pSaoParam.SaoLcuParam[0][iAddr].MergeUpFlag
                } else {
                    pSaoParam.SaoLcuParam[iCompIdx][iAddr].MergeUpFlag = false
                }
                if !pSaoParam.SaoLcuParam[iCompIdx][iAddr].MergeUpFlag {
                    pSaoParam.SaoLcuParam[2][iAddr].TypeIdx = pSaoParam.SaoLcuParam[1][iAddr].TypeIdx
                    this.ParseSaoOffset(&(pSaoParam.SaoLcuParam[iCompIdx][iAddr]), uint(iCompIdx))
                } else {
                    this.CopySaoOneLcuParam(&pSaoParam.SaoLcuParam[iCompIdx][iAddr], &pSaoParam.SaoLcuParam[iCompIdx][iAddr-pSaoParam.NumCuInWidth])
                }
            } else {
                this.CopySaoOneLcuParam(&pSaoParam.SaoLcuParam[iCompIdx][iAddr], &pSaoParam.SaoLcuParam[iCompIdx][iAddr-1])
            }
        } else {
            pSaoParam.SaoLcuParam[iCompIdx][iAddr].TypeIdx = -1
            pSaoParam.SaoLcuParam[iCompIdx][iAddr].SubTypeIdx = 0
        }
    }
}

var iTypeLength = [TLibCommon.MAX_NUM_SAO_TYPE]int{
    TLibCommon.SAO_EO_LEN,
    TLibCommon.SAO_EO_LEN,
    TLibCommon.SAO_EO_LEN,
    TLibCommon.SAO_EO_LEN,
    TLibCommon.SAO_BO_LEN,
}

func (this *TDecSbac) ParseSaoOffset(psSaoLcuParam *TLibCommon.SaoLcuParam, compIdx uint) {
    var uiSymbol uint

    if compIdx == 2 {
        uiSymbol = uint(psSaoLcuParam.TypeIdx + 1)
    } else {
        this.ParseSaoTypeIdx(&uiSymbol)
        //#ifdef ENC_DEC_TRACE
        if compIdx == 0 {
            this.XReadAeTr(int(uiSymbol), "sao_type_idx_luma", TLibCommon.TRACE_LCU)
        } else {
            this.XReadAeTr(int(uiSymbol), "sao_type_idx_chroma", TLibCommon.TRACE_LCU)
        }
        //#endif
    }
    psSaoLcuParam.TypeIdx = int(uiSymbol) - 1
    if uiSymbol != 0 {
        psSaoLcuParam.Length = iTypeLength[psSaoLcuParam.TypeIdx]

        var bitDepth, offsetTh int
        if compIdx != 0 {
            bitDepth = TLibCommon.G_bitDepthC
        } else {
            bitDepth = TLibCommon.G_bitDepthY
        }
        offsetTh = 1 << uint(TLibCommon.MIN(bitDepth-5, 5).(int))

        if psSaoLcuParam.TypeIdx == TLibCommon.SAO_BO {
            for i := 0; i < psSaoLcuParam.Length; i++ {
                this.ParseSaoMaxUvlc(&uiSymbol, uint(offsetTh-1))
                psSaoLcuParam.Offset[i] = int(uiSymbol)
                //#ifdef ENC_DEC_TRACE
                this.XReadAeTr(int(uiSymbol), "sao_offset_abs", TLibCommon.TRACE_LCU)
                //#endif
            }
            for i := 0; i < psSaoLcuParam.Length; i++ {
                if psSaoLcuParam.Offset[i] != 0 {
                    this.m_pcTDecBinIf.DecodeBinEP(&uiSymbol)
                    //#ifdef ENC_DEC_TRACE
                    this.XReadAeTr(int(uiSymbol), "sao_offset_sign", TLibCommon.TRACE_LCU)
                    //#endif
                    if uiSymbol != 0 {
                        psSaoLcuParam.Offset[i] = -psSaoLcuParam.Offset[i]
                    }
                }
            }
            this.ParseSaoUflc(5, &uiSymbol)
            psSaoLcuParam.SubTypeIdx = int(uiSymbol)
            //#ifdef ENC_DEC_TRACE
            this.XReadAeTr(int(uiSymbol), "sao_band_position", TLibCommon.TRACE_LCU)
            //#endif
        } else if psSaoLcuParam.TypeIdx < 4 {
            this.ParseSaoMaxUvlc(&uiSymbol, uint(offsetTh-1))
            psSaoLcuParam.Offset[0] = int(uiSymbol)
            //#ifdef ENC_DEC_TRACE
            this.XReadAeTr(int(uiSymbol), "sao_offset_abs", TLibCommon.TRACE_LCU)
            //#endif
            this.ParseSaoMaxUvlc(&uiSymbol, uint(offsetTh-1))
            psSaoLcuParam.Offset[1] = int(uiSymbol)
            //#ifdef ENC_DEC_TRACE
            this.XReadAeTr(int(uiSymbol), "sao_offset_abs", TLibCommon.TRACE_LCU)
            //#endif
            this.ParseSaoMaxUvlc(&uiSymbol, uint(offsetTh-1))
            psSaoLcuParam.Offset[2] = -int(uiSymbol)
            //#ifdef ENC_DEC_TRACE
            this.XReadAeTr(int(uiSymbol), "sao_offset_abs", TLibCommon.TRACE_LCU)
            //#endif
            this.ParseSaoMaxUvlc(&uiSymbol, uint(offsetTh-1))
            psSaoLcuParam.Offset[3] = -int(uiSymbol)
            //#ifdef ENC_DEC_TRACE
            this.XReadAeTr(int(uiSymbol), "sao_offset_abs", TLibCommon.TRACE_LCU)
            //#endif
            if compIdx != 2 {
                this.ParseSaoUflc(2, &uiSymbol)
                psSaoLcuParam.SubTypeIdx = int(uiSymbol)
                psSaoLcuParam.TypeIdx += psSaoLcuParam.SubTypeIdx
                //#ifdef ENC_DEC_TRACE
                if compIdx == 0 {
                    this.XReadAeTr(int(uiSymbol), "sao_eo_class_luma", TLibCommon.TRACE_LCU)
                } else {
                    this.XReadAeTr(int(uiSymbol), "sao_eo_class_chroma", TLibCommon.TRACE_LCU)
                }
                //#endif
            }
        }
    } else {
        psSaoLcuParam.Length = 0
    }
}

//private:
func (this *TDecSbac) xReadUnarySymbol(ruiSymbol *uint, pcSCModel []TLibCommon.ContextModel, iOffset int) {
    this.m_pcTDecBinIf.DecodeBin(ruiSymbol, &pcSCModel[0])

    if *ruiSymbol != 0 {
        return
    }

    uiSymbol := uint(0)
    uiCont := uint(1)

    for uiCont != 0 {
        this.m_pcTDecBinIf.DecodeBin(&uiCont, &pcSCModel[iOffset])
        uiSymbol++
    }

    *ruiSymbol = uiSymbol
}
func (this *TDecSbac) xReadUnaryMaxSymbol(ruiSymbol *uint, pcSCModel []TLibCommon.ContextModel, iOffset, uiMaxSymbol uint) {
    if uiMaxSymbol == 0 {
        *ruiSymbol = 0
        return
    }

    this.m_pcTDecBinIf.DecodeBin(ruiSymbol, &pcSCModel[0])

    if *ruiSymbol == 0 || uiMaxSymbol == 1 {
        return
    }

    uiSymbol := uint(0)
    uiCont := uint(1)

    for uiCont != 0 && (uiSymbol < uiMaxSymbol-1) {
        this.m_pcTDecBinIf.DecodeBin(&uiCont, &pcSCModel[iOffset])
        uiSymbol++
    }

    if uiCont != 0 && (uiSymbol == uiMaxSymbol-1) {
        uiSymbol++
    }

    *ruiSymbol = uiSymbol
}
func (this *TDecSbac) xReadEpExGolomb(ruiSymbol *uint, uiCount uint) {
    uiSymbol := uint(0)
    uiBit := uint(1)

    for uiBit != 0 {
        this.m_pcTDecBinIf.DecodeBinEP(&uiBit)
        uiSymbol += uiBit << uiCount
        uiCount++
    }

    uiCount--
    if uiCount != 0 {
        var bins uint
        this.m_pcTDecBinIf.DecodeBinsEP(&bins, int(uiCount))
        uiSymbol += bins
    }

    *ruiSymbol = uiSymbol
}
func (this *TDecSbac) xReadCoefRemainExGolomb(rSymbol *uint, rParam uint) {
    prefix := uint(0)
    codeWord := uint(1)
    for codeWord != 0 {
        prefix++
        this.m_pcTDecBinIf.DecodeBinEP(&codeWord)
    }

    codeWord = 1 - codeWord
    prefix -= codeWord
    codeWord = 0
    if prefix < TLibCommon.COEF_REMAIN_BIN_REDUCTION {
        this.m_pcTDecBinIf.DecodeBinsEP(&codeWord, int(rParam))
        *rSymbol = (prefix << rParam) + codeWord
    } else {
        this.m_pcTDecBinIf.DecodeBinsEP(&codeWord, int(prefix-TLibCommon.COEF_REMAIN_BIN_REDUCTION+rParam))
        *rSymbol = (((1 << (prefix - TLibCommon.COEF_REMAIN_BIN_REDUCTION)) + TLibCommon.COEF_REMAIN_BIN_REDUCTION - 1) << rParam) + codeWord
    }
}

//public:
func (this *TDecSbac) ParseSkipFlag(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint) {
    if pcCU.GetSlice().IsIntra() {
        return
    }

    uiSymbol := uint(0)
    uiCtxSkip := pcCU.GetCtxSkipFlag(uiAbsPartIdx)
    this.m_pcTDecBinIf.DecodeBin(&uiSymbol, this.m_cCUSkipFlagSCModel.Get3(0, 0, uiCtxSkip))
    //this.DTRACE_CABAC_VL( g_nSymbolCounter++ );
    /*this.DTRACE_CABAC_T("\tSkipFlag")
    this.DTRACE_CABAC_T("\tuiCtxSkip: ")
    this.DTRACE_CABAC_V(uiCtxSkip)
    this.DTRACE_CABAC_T("\tuiSymbol: ")
    this.DTRACE_CABAC_V(uiSymbol)
    this.DTRACE_CABAC_T("\n")*/

    if uiSymbol != 0 {
        pcCU.SetSkipFlagSubParts(true, uiAbsPartIdx, uiDepth)
        pcCU.SetPredModeSubParts(TLibCommon.MODE_INTER, uiAbsPartIdx, uiDepth)
        pcCU.SetPartSizeSubParts(TLibCommon.SIZE_2Nx2N, uiAbsPartIdx, uiDepth)
        pcCU.SetSizeSubParts(pcCU.GetSlice().GetSPS().GetMaxCUWidth()>>uiDepth, pcCU.GetSlice().GetSPS().GetMaxCUHeight()>>uiDepth, uiAbsPartIdx, uiDepth)
        pcCU.SetMergeFlagSubParts(true, uiAbsPartIdx, 0, uiDepth)
    }
}
func (this *TDecSbac) ParseCUTransquantBypassFlag(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint) {
    var uiSymbol uint
    this.m_pcTDecBinIf.DecodeBin(&uiSymbol, this.m_CUTransquantBypassFlagSCModel.Get3(0, 0, 0))
    pcCU.SetCUTransquantBypassSubParts(uiSymbol != 0, uiAbsPartIdx, uiDepth)
}
func (this *TDecSbac) ParseSplitFlag(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint) {
    if uiDepth == pcCU.GetSlice().GetSPS().GetMaxCUDepth()-pcCU.GetSlice().GetSPS().GetAddCUDepth() {
        pcCU.SetDepthSubParts(uiDepth, uiAbsPartIdx)
        return
    }

    var uiSymbol uint
    this.m_pcTDecBinIf.DecodeBin(&uiSymbol, this.m_cCUSplitFlagSCModel.Get3(0, 0, pcCU.GetCtxSplitFlag(uiAbsPartIdx, uiDepth)))
    //this.DTRACE_CABAC_VL( g_nSymbolCounter++ )
    //this.DTRACE_CABAC_T("\tSplitFlag\n")
    pcCU.SetDepthSubParts(uiDepth+uiSymbol, uiAbsPartIdx)

    return
}
func (this *TDecSbac) ParseMergeFlag(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth, uiPUIdx uint) {
    var uiSymbol uint
    this.m_pcTDecBinIf.DecodeBin(&uiSymbol, this.m_cCUMergeFlagExtSCModel.Get3(0, 0, 0)) // &this.m_cCUMergeFlagExtSCModel.Get1( 0 )[0] );
    pcCU.SetMergeFlagSubParts(uiSymbol != 0, uiAbsPartIdx, uiPUIdx, uiDepth)

    /*this.DTRACE_CABAC_VL( g_nSymbolCounter++ );*/
    /*this.DTRACE_CABAC_T("\tMergeFlag: ")
    this.DTRACE_CABAC_V(uiSymbol)
    this.DTRACE_CABAC_T("\tAddress: ")
    this.DTRACE_CABAC_V(pcCU.GetAddr())
    this.DTRACE_CABAC_T("\tuiAbsPartIdx: ")
    this.DTRACE_CABAC_V(uiAbsPartIdx)
    this.DTRACE_CABAC_T("\n")*/
}
func (this *TDecSbac) ParseMergeIndex(pcCU *TLibCommon.TComDataCU, ruiMergeIndex *uint) {
    uiUnaryIdx := uint(0)
    uiNumCand := pcCU.GetSlice().GetMaxNumMergeCand()
    if uiNumCand > 1 {
        for ; uiUnaryIdx < uiNumCand-1; uiUnaryIdx++ {
            uiSymbol := uint(0)
            if uiUnaryIdx == 0 {
                this.m_pcTDecBinIf.DecodeBin(&uiSymbol, this.m_cCUMergeIdxExtSCModel.Get3(0, 0, 0))
            } else {
                this.m_pcTDecBinIf.DecodeBinEP(&uiSymbol)
            }

            if uiSymbol == 0 {
                break
            }
        }
    }
    *ruiMergeIndex = uiUnaryIdx

    /*this.DTRACE_CABAC_VL( g_nSymbolCounter++ )*/
    /*this.DTRACE_CABAC_T("\tparseMergeIndex()")
    this.DTRACE_CABAC_T("\tuiMRGIdx= ")
    this.DTRACE_CABAC_V(*ruiMergeIndex)
    this.DTRACE_CABAC_T("\n")*/
}
func (this *TDecSbac) ParsePartSize(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint) {
    var uiSymbol, uiMode uint
    uiMode = 0
    var eMode TLibCommon.PartSize

    if pcCU.IsIntra(uiAbsPartIdx) {
        uiSymbol = 1
        if uiDepth == pcCU.GetSlice().GetSPS().GetMaxCUDepth()-pcCU.GetSlice().GetSPS().GetAddCUDepth() {
            this.m_pcTDecBinIf.DecodeBin(&uiSymbol, this.m_cCUPartSizeSCModel.Get3(0, 0, 0))
        }
        if uiSymbol != 0 {
            eMode = TLibCommon.SIZE_2Nx2N
        } else {
            eMode = TLibCommon.SIZE_NxN
        }

        uiTrLevel := uint(0)
        uiWidthInBit := uint(TLibCommon.G_aucConvertToBit[pcCU.GetWidth1(uiAbsPartIdx)] + 2)
        uiTrSizeInBit := uint(TLibCommon.G_aucConvertToBit[pcCU.GetSlice().GetSPS().GetMaxTrSize()] + 2)
        if uiWidthInBit >= uiTrSizeInBit {
            uiTrLevel = uiWidthInBit - uiTrSizeInBit
        } else {
            uiTrLevel = 0
        }
        if eMode == TLibCommon.SIZE_NxN {
            pcCU.SetTrIdxSubParts(1+uiTrLevel, uiAbsPartIdx, uiDepth)
        } else {
            pcCU.SetTrIdxSubParts(uiTrLevel, uiAbsPartIdx, uiDepth)
        }
    } else {
        uiMaxNumBits := uint(2)
        if uiDepth == pcCU.GetSlice().GetSPS().GetMaxCUDepth()-pcCU.GetSlice().GetSPS().GetAddCUDepth() && !((pcCU.GetSlice().GetSPS().GetMaxCUWidth()>>uiDepth) == 8 && (pcCU.GetSlice().GetSPS().GetMaxCUHeight()>>uiDepth) == 8) {
            uiMaxNumBits++
        }
        for ui := uint(0); ui < uiMaxNumBits; ui++ {
            this.m_pcTDecBinIf.DecodeBin(&uiSymbol, this.m_cCUPartSizeSCModel.Get3(0, 0, ui))
            if uiSymbol != 0 {
                break
            }
            uiMode++
        }
        eMode = TLibCommon.PartSize(uiMode)
        if pcCU.GetSlice().GetSPS().GetAMPAcc(uiDepth) != 0 {
            if eMode == TLibCommon.SIZE_2NxN {
                this.m_pcTDecBinIf.DecodeBin(&uiSymbol, this.m_cCUAMPSCModel.Get3(0, 0, 0))
                if uiSymbol == 0 {
                    this.m_pcTDecBinIf.DecodeBinEP(&uiSymbol)
                    if uiSymbol == 0 {
                        eMode = TLibCommon.SIZE_2NxnU
                    } else {
                        eMode = TLibCommon.SIZE_2NxnD
                    }
                }
            } else if eMode == TLibCommon.SIZE_Nx2N {
                this.m_pcTDecBinIf.DecodeBin(&uiSymbol, this.m_cCUAMPSCModel.Get3(0, 0, 0))
                if uiSymbol == 0 {
                    this.m_pcTDecBinIf.DecodeBinEP(&uiSymbol)
                    if uiSymbol == 0 {
                        eMode = TLibCommon.SIZE_nLx2N
                    } else {
                        eMode = TLibCommon.SIZE_nRx2N
                    }
                }
            }
        }
    }
    pcCU.SetPartSizeSubParts(eMode, uiAbsPartIdx, uiDepth)
    pcCU.SetSizeSubParts(pcCU.GetSlice().GetSPS().GetMaxCUWidth()>>uiDepth, pcCU.GetSlice().GetSPS().GetMaxCUHeight()>>uiDepth, uiAbsPartIdx, uiDepth)
}
func (this *TDecSbac) ParsePredMode(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint) {
    if pcCU.GetSlice().IsIntra() {
        pcCU.SetPredModeSubParts(TLibCommon.MODE_INTRA, uiAbsPartIdx, uiDepth)
        return
    }

    var uiSymbol uint
    iPredMode := TLibCommon.MODE_INTER
    this.m_pcTDecBinIf.DecodeBin(&uiSymbol, this.m_cCUPredModeSCModel.Get3(0, 0, 0))
    iPredMode += int(uiSymbol)
    pcCU.SetPredModeSubParts(TLibCommon.PredMode(iPredMode), uiAbsPartIdx, uiDepth)
}

func (this *TDecSbac) ParseIntraDirLumaAng(pcCU *TLibCommon.TComDataCU, absPartIdx, depth uint) {
    mode := pcCU.GetPartitionSize1(absPartIdx)
    var partNum uint
    if mode == TLibCommon.SIZE_NxN {
        partNum = 4
    } else {
        partNum = 1
    }

    partOffset := uint(pcCU.GetPic().GetNumPartInCU()>>(pcCU.GetDepth1(absPartIdx)<<1)) >> 2
    var mpmPred [4]uint
    var symbol uint
    var j, intraPredMode int
    if mode == TLibCommon.SIZE_NxN {
        depth++
    }
    for j = 0; j < int(partNum); j++ {
        this.m_pcTDecBinIf.DecodeBin(&symbol, this.m_cCUIntraPredSCModel.Get3(0, 0, 0))
        mpmPred[j] = symbol
    }
    for j = 0; j < int(partNum); j++ {
        var preds = [3]int{-1, -1, -1}
        predNum := pcCU.GetIntraDirLumaPredictor(absPartIdx+partOffset*uint(j), preds[:], nil)
        if mpmPred[j] != 0 {
            this.m_pcTDecBinIf.DecodeBinEP(&symbol)
            if symbol != 0 {
                this.m_pcTDecBinIf.DecodeBinEP(&symbol)
                symbol++
            }
            intraPredMode = preds[symbol]
        } else {
            intraPredMode = 0
            this.m_pcTDecBinIf.DecodeBinsEP(&symbol, 5)
            intraPredMode = int(symbol)

            //postponed sorting of MPMs (only in remaining branch)
            if preds[0] > preds[1] {
                tmp := preds[0]
                preds[0] = preds[1]
                preds[1] = tmp
                //std::swap(preds[0], preds[1]);
            }
            if preds[0] > preds[2] {
                tmp := preds[0]
                preds[0] = preds[2]
                preds[2] = tmp
                //std::swap(preds[0], preds[2]);
            }
            if preds[1] > preds[2] {
                tmp := preds[1]
                preds[1] = preds[2]
                preds[2] = tmp
                //std::swap(preds[1], preds[2]);
            }
            for i := int(0); i < predNum; i++ {
                intraPredMode += int(TLibCommon.B2U(intraPredMode >= preds[i]))
            }
        }
        pcCU.SetLumaIntraDirSubParts(uint(intraPredMode), absPartIdx+partOffset*uint(j), depth)
    }
}

func (this *TDecSbac) ParseIntraDirChroma(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint) {
    var uiSymbol uint

    this.m_pcTDecBinIf.DecodeBin(&uiSymbol, this.m_cCUChromaPredSCModel.Get3(0, 0, 0))

    if uiSymbol == 0 {
        uiSymbol = TLibCommon.DM_CHROMA_IDX
    } else {
        var uiIPredMode uint
        this.m_pcTDecBinIf.DecodeBinsEP(&uiIPredMode, 2)
        var uiAllowedChromaDir [TLibCommon.NUM_CHROMA_MODE]uint
        pcCU.GetAllowedChromaDir(uiAbsPartIdx, uiAllowedChromaDir[:])
        uiSymbol = uiAllowedChromaDir[uiIPredMode]
    }
    pcCU.SetChromIntraDirSubParts(uiSymbol, uiAbsPartIdx, uiDepth)
    return
}

func (this *TDecSbac) ParseInterDir(pcCU *TLibCommon.TComDataCU, ruiInterDir *uint, uiAbsPartIdx uint) {
    var uiSymbol uint
    uiCtx := pcCU.GetCtxInterDir(uiAbsPartIdx)
    pCtx := this.m_cCUInterDirSCModel.Get1(0)
    uiSymbol = 0
    if pcCU.GetPartitionSize1(uiAbsPartIdx) == TLibCommon.SIZE_2Nx2N || pcCU.GetHeight1(uiAbsPartIdx) != 8 {
        this.m_pcTDecBinIf.DecodeBin(&uiSymbol, &pCtx[uiCtx])
    }

    if uiSymbol != 0 {
        uiSymbol = 2
    } else {
        this.m_pcTDecBinIf.DecodeBin(&uiSymbol, &pCtx[4])
        //assert(uiSymbol == 0 || uiSymbol == 1);
    }

    uiSymbol++
    *ruiInterDir = uiSymbol
    return
}
func (this *TDecSbac) ParseRefFrmIdx(pcCU *TLibCommon.TComDataCU, riRefFrmIdx *int, eRefList TLibCommon.RefPicList) {
    var uiSymbol uint
    {
        pCtx := this.m_cCURefPicSCModel.Get1(0)
        this.m_pcTDecBinIf.DecodeBin(&uiSymbol, &pCtx[0])

        if uiSymbol != 0 {
            uiRefNum := pcCU.GetSlice().GetNumRefIdx(eRefList) - 2
            //pCtx++;
            var ui uint
            for ui = 0; ui < uint(uiRefNum); ui++ {
                if ui == 0 {
                    this.m_pcTDecBinIf.DecodeBin(&uiSymbol, &pCtx[1])
                } else {
                    this.m_pcTDecBinIf.DecodeBinEP(&uiSymbol)
                }
                if uiSymbol == 0 {
                    break
                }
            }
            uiSymbol = ui + 1
        }
        *riRefFrmIdx = int(uiSymbol)
    }

    return
}
func (this *TDecSbac) ParseMvd(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiPartIdx, uiDepth uint, eRefList TLibCommon.RefPicList) {
    var uiSymbol uint
    var uiHorAbs, uiVerAbs, uiHorSign, uiVerSign uint
    uiHorSign = 0
    uiVerSign = 0
    pCtx := this.m_cCUMvdSCModel.Get1(0)

    if pcCU.GetSlice().GetMvdL1ZeroFlag() && eRefList == TLibCommon.REF_PIC_LIST_1 && pcCU.GetInterDir1(uiAbsPartIdx) == 3 {
        uiHorAbs = 0
        uiVerAbs = 0
    } else {
        this.m_pcTDecBinIf.DecodeBin(&uiHorAbs, &pCtx[0])
        this.m_pcTDecBinIf.DecodeBin(&uiVerAbs, &pCtx[0])

        bHorAbsGr0 := uiHorAbs != 0
        bVerAbsGr0 := uiVerAbs != 0
        //pCtx++;

        if bHorAbsGr0 {
            this.m_pcTDecBinIf.DecodeBin(&uiSymbol, &pCtx[1])
            uiHorAbs += uiSymbol
        }

        if bVerAbsGr0 {
            this.m_pcTDecBinIf.DecodeBin(&uiSymbol, &pCtx[1])
            uiVerAbs += uiSymbol
        }

        if bHorAbsGr0 {
            if 2 == uiHorAbs {
                this.xReadEpExGolomb(&uiSymbol, 1)
                uiHorAbs += uiSymbol
            }

            this.m_pcTDecBinIf.DecodeBinEP(&uiHorSign)
        }

        if bVerAbsGr0 {
            if 2 == uiVerAbs {
                this.xReadEpExGolomb(&uiSymbol, 1)
                uiVerAbs += uiSymbol
            }

            this.m_pcTDecBinIf.DecodeBinEP(&uiVerSign)
        }

    }

    var mv_x, mv_y int16
    if uiHorSign != 0 {
        mv_x = -int16(uiHorAbs)
    } else {
        mv_x = int16(uiHorAbs)
    }
    if uiVerSign != 0 {
        mv_y = -int16(uiVerAbs)
    } else {
        mv_y = int16(uiVerAbs)
    }
    //const TComMv cMv( uiHorSign ? -Int( uiHorAbs ): uiHorAbs, uiVerSign ? -Int( uiVerAbs ) : uiVerAbs );
    cMv := TLibCommon.NewTComMv(mv_x, mv_y)

    pcCU.GetCUMvField(eRefList).SetAllMvd(*cMv, pcCU.GetPartitionSize1(uiAbsPartIdx), int(uiAbsPartIdx), uiDepth, int(uiPartIdx))
    return
}

func (this *TDecSbac) ParseTransformSubdivFlag(ruiSubdivFlag *uint, uiLog2TransformBlockSize uint) {
    this.m_pcTDecBinIf.DecodeBin(ruiSubdivFlag, this.m_cCUTransSubdivFlagSCModel.Get3(0, 0, uiLog2TransformBlockSize))
    /*this.DTRACE_CABAC_VL( g_nSymbolCounter++ )*/
    /*this.DTRACE_CABAC_T("\tparseTransformSubdivFlag()")
    this.DTRACE_CABAC_T("\tsymbol=")
    this.DTRACE_CABAC_V(*ruiSubdivFlag)
    this.DTRACE_CABAC_T("\tctx=")
    this.DTRACE_CABAC_V(uiLog2TransformBlockSize)
    this.DTRACE_CABAC_T("\n")*/
}
func (this *TDecSbac) ParseQtCbf(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint, eType TLibCommon.TextType, uiTrDepth, uiDepth uint) {
    var uiSymbol uint
    uiCtx := pcCU.GetCtxQtCbf(eType, uiTrDepth)
    //fmt.Printf("uiCtx=%d\n",uiCtx);
    if eType != 0 {
        this.m_pcTDecBinIf.DecodeBin(&uiSymbol, this.m_cCUQtCbfSCModel.Get3(0, TLibCommon.TEXT_CHROMA, uiCtx))
    } else {
        this.m_pcTDecBinIf.DecodeBin(&uiSymbol, this.m_cCUQtCbfSCModel.Get3(0, uint(eType), uiCtx))
    }
    /*this.DTRACE_CABAC_VL( g_nSymbolCounter++ )*/
    /*this.DTRACE_CABAC_T("\tparseQtCbf()")
    this.DTRACE_CABAC_T("\tsymbol=")
    this.DTRACE_CABAC_V(uiSymbol)
    this.DTRACE_CABAC_T("\tctx=")
    this.DTRACE_CABAC_V(uiCtx)
    this.DTRACE_CABAC_T("\tetype=")
    this.DTRACE_CABAC_V(uint(eType))
    this.DTRACE_CABAC_T("\tuiAbsPartIdx=")
    this.DTRACE_CABAC_V(uiAbsPartIdx)
    this.DTRACE_CABAC_T("\n")*/

    pcCU.SetCbfSubParts4(byte(uiSymbol<<uiTrDepth), eType, uiAbsPartIdx, uiDepth)
}
func (this *TDecSbac) ParseQtRootCbf(uiAbsPartIdx uint, uiQtRootCbf *uint) {
    var uiSymbol uint
    uiCtx := uint(0)
    this.m_pcTDecBinIf.DecodeBin(&uiSymbol, this.m_cCUQtRootCbfSCModel.Get3(0, 0, uiCtx))
    /*this.DTRACE_CABAC_VL( g_nSymbolCounter++ )*/
    /*this.DTRACE_CABAC_T("\tparseQtRootCbf()")
    this.DTRACE_CABAC_T("\tsymbol=")
    this.DTRACE_CABAC_V(uiSymbol)
    this.DTRACE_CABAC_T("\tctx=")
    this.DTRACE_CABAC_V(uiCtx)
    this.DTRACE_CABAC_T("\tuiAbsPartIdx=")
    this.DTRACE_CABAC_V(uiAbsPartIdx)
    this.DTRACE_CABAC_T("\n")*/

    *uiQtRootCbf = uiSymbol
}

func (this *TDecSbac) ParseDeltaQP(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint) {
    var qp, iDQp int
    var uiDQp uint
    var uiSymbol uint

    this.xReadUnaryMaxSymbol(&uiDQp, this.m_cCUDeltaQpSCModel.Get1(0), 1, TLibCommon.CU_DQP_TU_CMAX)

    if uiDQp >= TLibCommon.CU_DQP_TU_CMAX {
        this.xReadEpExGolomb(&uiSymbol, TLibCommon.CU_DQP_EG_k)
        uiDQp += uiSymbol
    }

    if uiDQp > 0 {
        var uiSign uint
        qpBdOffsetY := pcCU.GetSlice().GetSPS().GetQpBDOffsetY()
        this.m_pcTDecBinIf.DecodeBinEP(&uiSign)
        iDQp = int(uiDQp)
        if uiSign != 0 {
            iDQp = -iDQp
        }
        qp = ((int(pcCU.GetRefQP(uiAbsPartIdx)) + iDQp + 52 + 2*qpBdOffsetY) % (52 + qpBdOffsetY)) - qpBdOffsetY
    } else {
        iDQp = 0
        qp = int(pcCU.GetRefQP(uiAbsPartIdx))
    }
    pcCU.SetQPSubParts(qp, uiAbsPartIdx, uiDepth)
    pcCU.SetCodedQP(int8(qp))
}

func (this *TDecSbac) ParseIPCMInfo(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint) {
    var uiSymbol uint

    readPCMSampleFlag := false

    this.m_pcTDecBinIf.DecodeBinTrm(&uiSymbol)

    if uiSymbol != 0 {
        readPCMSampleFlag = true

        this.m_pcTDecBinIf.DecodePCMAlignBits()
    }

    if readPCMSampleFlag == true {
        bIpcmFlag := true

        pcCU.SetPartSizeSubParts(TLibCommon.SIZE_2Nx2N, uiAbsPartIdx, uiDepth)
        pcCU.SetSizeSubParts(pcCU.GetSlice().GetSPS().GetMaxCUWidth()>>uiDepth, pcCU.GetSlice().GetSPS().GetMaxCUHeight()>>uiDepth, uiAbsPartIdx, uiDepth)
        pcCU.SetTrIdxSubParts(0, uiAbsPartIdx, uiDepth)
        pcCU.SetIPCMFlagSubParts(bIpcmFlag, uiAbsPartIdx, uiDepth)

        uiMinCoeffSize := pcCU.GetPic().GetMinCUWidth() * pcCU.GetPic().GetMinCUHeight()
        uiLumaOffset := uiMinCoeffSize * uiAbsPartIdx
        uiChromaOffset := uiLumaOffset >> 2

        //Pel* piPCMSample;
        var uiWidth, uiHeight, uiSampleBits, uiX, uiY uint

        piPCMSample := pcCU.GetPCMSampleY()[uiLumaOffset:]
        uiWidth = uint(pcCU.GetWidth1(uiAbsPartIdx))
        uiHeight = uint(pcCU.GetHeight1(uiAbsPartIdx))
        uiSampleBits = uint(pcCU.GetSlice().GetSPS().GetPCMBitDepthLuma())

        for uiY = 0; uiY < uiHeight; uiY++ {
            for uiX = 0; uiX < uiWidth; uiX++ {
                var uiSample uint
                this.m_pcTDecBinIf.xReadPCMCode(uiSampleBits, &uiSample)
                piPCMSample[uiY*uiWidth+uiX] = TLibCommon.Pel(uiSample)
            }
            //piPCMSample += uiWidth;
        }

        piPCMSample = pcCU.GetPCMSampleCb()[uiChromaOffset:]
        uiWidth = uint(pcCU.GetWidth1(uiAbsPartIdx)) / 2
        uiHeight = uint(pcCU.GetHeight1(uiAbsPartIdx)) / 2
        uiSampleBits = uint(pcCU.GetSlice().GetSPS().GetPCMBitDepthChroma())

        for uiY = 0; uiY < uiHeight; uiY++ {
            for uiX = 0; uiX < uiWidth; uiX++ {
                var uiSample uint
                this.m_pcTDecBinIf.xReadPCMCode(uiSampleBits, &uiSample)
                piPCMSample[uiY*uiWidth+uiX] = TLibCommon.Pel(uiSample)
            }
            //piPCMSample += uiWidth;
        }

        piPCMSample = pcCU.GetPCMSampleCr()[uiChromaOffset:]
        uiWidth = uint(pcCU.GetWidth1(uiAbsPartIdx)) / 2
        uiHeight = uint(pcCU.GetHeight1(uiAbsPartIdx)) / 2
        uiSampleBits = uint(pcCU.GetSlice().GetSPS().GetPCMBitDepthChroma())

        for uiY = 0; uiY < uiHeight; uiY++ {
            for uiX = 0; uiX < uiWidth; uiX++ {
                var uiSample uint
                this.m_pcTDecBinIf.xReadPCMCode(uiSampleBits, &uiSample)
                piPCMSample[uiY*uiWidth+uiX] = TLibCommon.Pel(uiSample)
            }
            //piPCMSample += uiWidth;
        }
        this.m_pcTDecBinIf.ResetBac()
    }
}

func (this *TDecSbac) ParseLastSignificantXY(uiPosLastX *uint, uiPosLastY *uint, width, height int, eTType TLibCommon.TextType, uiScanIdx uint) {
    var uiLast uint
    pCtxX := this.m_cCuCtxLastX.Get2(0, uint(eTType))
    pCtxY := this.m_cCuCtxLastY.Get2(0, uint(eTType))

    var blkSizeOffsetX, blkSizeOffsetY, shiftX, shiftY int
    if eTType != 0 {
        blkSizeOffsetX = 0
        blkSizeOffsetY = 0
        shiftX = int(TLibCommon.G_aucConvertToBit[width])
        shiftY = int(TLibCommon.G_aucConvertToBit[height])
    } else {
        blkSizeOffsetX = int(TLibCommon.G_aucConvertToBit[width]*3 + ((TLibCommon.G_aucConvertToBit[width] + 1) >> 2))
        blkSizeOffsetY = int(TLibCommon.G_aucConvertToBit[height]*3 + ((TLibCommon.G_aucConvertToBit[height] + 1) >> 2))
        shiftX = int((TLibCommon.G_aucConvertToBit[width] + 3) >> 2)
        shiftY = int((TLibCommon.G_aucConvertToBit[height] + 3) >> 2)
    }
    // posX
    for *uiPosLastX = 0; *uiPosLastX < TLibCommon.G_uiGroupIdx[width-1]; (*uiPosLastX)++ {
        this.m_pcTDecBinIf.DecodeBin(&uiLast, &pCtxX[blkSizeOffsetX+int(*uiPosLastX>>uint(shiftX))])
        //fmt.Printf("uiLast=%d\n", uiLast);
        if uiLast == 0 {
            break
        }
    }

    // posY
    for *uiPosLastY = 0; *uiPosLastY < TLibCommon.G_uiGroupIdx[height-1]; (*uiPosLastY)++ {
        this.m_pcTDecBinIf.DecodeBin(&uiLast, &pCtxY[blkSizeOffsetY+int(*uiPosLastY>>uint(shiftY))])
        if uiLast == 0 {
            break
        }
    }
    if *uiPosLastX > 3 {
        uiTemp := uint(0)
        uiCount := uint(*uiPosLastX-2) >> 1
        for i := int(uiCount) - 1; i >= 0; i-- {
            this.m_pcTDecBinIf.DecodeBinEP(&uiLast)
            uiTemp += uiLast << uint(i)
        }
        *uiPosLastX = TLibCommon.G_uiMinInGroup[*uiPosLastX] + uiTemp
    }
    if *uiPosLastY > 3 {
        uiTemp := uint(0)
        uiCount := uint(*uiPosLastY-2) >> 1
        for i := int(uiCount) - 1; i >= 0; i-- {
            this.m_pcTDecBinIf.DecodeBinEP(&uiLast)
            uiTemp += uiLast << uint(i)
        }
        *uiPosLastY = TLibCommon.G_uiMinInGroup[*uiPosLastY] + uiTemp
    }

    if uiScanIdx == TLibCommon.SCAN_VER {
        tmp := *uiPosLastX
        *uiPosLastX = *uiPosLastY
        *uiPosLastY = tmp
        //swap( uiPosLastX, uiPosLastY );
    }
}
func (this *TDecSbac) ParseCoeffNxN(pcCU *TLibCommon.TComDataCU, pcCoef []TLibCommon.TCoeff, uiAbsPartIdx, uiWidth, uiHeight, uiDepth uint, eTType TLibCommon.TextType) {
    /*this.DTRACE_CABAC_VL( TLibCommon.GnSymbolCounter++ )*/
    /*this.DTRACE_CABAC_T("\tparseCoeffNxN()\teType=")
    this.DTRACE_CABAC_V(uint(eTType))
    this.DTRACE_CABAC_T("\twidth=")
    this.DTRACE_CABAC_V(uiWidth)
    this.DTRACE_CABAC_T("\theight=")
    this.DTRACE_CABAC_V(uiHeight)
    this.DTRACE_CABAC_T("\tdepth=")
    this.DTRACE_CABAC_V(uiDepth)
    this.DTRACE_CABAC_T("\tabspartidx=")
    this.DTRACE_CABAC_V(uiAbsPartIdx)
    this.DTRACE_CABAC_T("\ttoCU-X=")
    this.DTRACE_CABAC_V(pcCU.GetCUPelX())
    this.DTRACE_CABAC_T("\ttoCU-Y=")
    this.DTRACE_CABAC_V(pcCU.GetCUPelY())
    this.DTRACE_CABAC_T("\tCU-addr=")
    this.DTRACE_CABAC_V(pcCU.GetAddr())
    this.DTRACE_CABAC_T("\tinCU-X=")
    this.DTRACE_CABAC_V(TLibCommon.G_auiRasterToPelX[TLibCommon.G_auiZscanToRaster[uiAbsPartIdx]])
    this.DTRACE_CABAC_T("\tinCU-Y=")
    this.DTRACE_CABAC_V(TLibCommon.G_auiRasterToPelY[TLibCommon.G_auiZscanToRaster[uiAbsPartIdx]])
    this.DTRACE_CABAC_T("\tpredmode=")
    this.DTRACE_CABAC_V(uint(pcCU.GetPredictionMode1(uiAbsPartIdx)))
    this.DTRACE_CABAC_T("\n")*/

    if uiWidth > pcCU.GetSlice().GetSPS().GetMaxTrSize() {
        uiWidth = pcCU.GetSlice().GetSPS().GetMaxTrSize()
        uiHeight = pcCU.GetSlice().GetSPS().GetMaxTrSize()
    }
    if pcCU.GetSlice().GetPPS().GetUseTransformSkip() {
        this.ParseTransformSkipFlags(pcCU, uiAbsPartIdx, uiWidth, uiHeight, uiDepth, eTType)
    }

    if eTType == TLibCommon.TEXT_LUMA {
        eTType = TLibCommon.TEXT_LUMA
    } else if eTType == TLibCommon.TEXT_NONE {
        eTType = TLibCommon.TEXT_NONE
    } else {
        eTType = TLibCommon.TEXT_CHROMA
    }

    //----- parse significance map -----
    uiLog2BlockSize := TLibCommon.G_aucConvertToBit[uiWidth] + 2
    uiMaxNumCoeff := uiWidth * uiHeight
    uiMaxNumCoeffM1 := uiMaxNumCoeff - 1
    uiScanIdx := pcCU.GetCoefScanIdx(uiAbsPartIdx, uiWidth, eTType == TLibCommon.TEXT_LUMA, pcCU.IsIntra(uiAbsPartIdx))

    //===== decode last significant =====
    var uiPosLastX, uiPosLastY uint
    this.ParseLastSignificantXY(&uiPosLastX, &uiPosLastY, int(uiWidth), int(uiHeight), eTType, uiScanIdx)
    uiBlkPosLast := uiPosLastX + (uiPosLastY << uint(uiLog2BlockSize))
    pcCoef[uiBlkPosLast] = 1

    //===== decode significance flags =====
    uiScanPosLast := uiBlkPosLast
    scan := TLibCommon.G_auiSigLastScan[uiScanIdx][uiLog2BlockSize-1]
    for uiScanPosLast = 0; uiScanPosLast < uiMaxNumCoeffM1; uiScanPosLast++ {
        uiBlkPos := scan[uiScanPosLast]
        if uiBlkPosLast == uiBlkPos {
            break
        }
    }

    baseCoeffGroupCtx := this.m_cCUSigCoeffGroupSCModel.Get2(0, uint(eTType))
    var baseCtx []TLibCommon.ContextModel
    if eTType == TLibCommon.TEXT_LUMA {
        baseCtx = this.m_cCUSigSCModel.Get2(0, 0)
    } else {
        baseCtx = this.m_cCUSigSCModel.Get2(0, 0)[TLibCommon.NUM_SIG_FLAG_CTX_LUMA:]
    }

    iLastScanSet := uiScanPosLast >> TLibCommon.LOG2_SCAN_SET_SIZE
    c1 := 1
    uiGoRiceParam := 0

    var beValid bool
    if pcCU.GetCUTransquantBypass1(uiAbsPartIdx) {
        beValid = false
    } else {
        beValid = pcCU.GetSlice().GetPPS().GetSignHideFlag()
    }
    absSum := uint(0)

    var uiSigCoeffGroupFlag [TLibCommon.MLS_GRP_NUM]uint
    //::memset( uiSigCoeffGroupFlag, 0, sizeof(UInt) * MLS_GRP_NUM );
    uiNumBlkSide := uiWidth >> (TLibCommon.MLS_CG_SIZE >> 1)
    var scanCG []uint
    {
        if uiLog2BlockSize > 3 {
            scanCG = TLibCommon.G_auiSigLastScan[uiScanIdx][uiLog2BlockSize-2-1]
        } else {
            scanCG = TLibCommon.G_auiSigLastScan[uiScanIdx][0]
        }
        if uiLog2BlockSize == 3 {
            scanCG = TLibCommon.G_sigLastScan8x8[uiScanIdx][:]
        } else if uiLog2BlockSize == 5 {
            scanCG = TLibCommon.G_sigLastScanCG32x32[:]
        }
    }

    iScanPosSig := int(uiScanPosLast)
    for iSubSet := int(iLastScanSet); iSubSet >= 0; iSubSet-- {
        iSubPos := iSubSet << TLibCommon.LOG2_SCAN_SET_SIZE
        uiGoRiceParam = 0
        numNonZero := 0

        lastNZPosInCG := -1
        firstNZPosInCG := TLibCommon.SCAN_SET_SIZE

        var pos [TLibCommon.SCAN_SET_SIZE]int
        if iScanPosSig == int(uiScanPosLast) {
            lastNZPosInCG = iScanPosSig
            firstNZPosInCG = iScanPosSig
            iScanPosSig--
            pos[numNonZero] = int(uiBlkPosLast)
            numNonZero = 1
        }

        // decode significant_coeffgroup_flag
        iCGBlkPos := scanCG[iSubSet]
        iCGPosY := iCGBlkPos / uiNumBlkSide
        iCGPosX := iCGBlkPos - (iCGPosY * uiNumBlkSide)
        if iSubSet == int(iLastScanSet) || iSubSet == 0 {
            uiSigCoeffGroupFlag[iCGBlkPos] = 1
        } else {
            var uiSigCoeffGroup uint
            uiCtxSig := TLibCommon.GetSigCoeffGroupCtxInc(uiSigCoeffGroupFlag[:], iCGPosX, iCGPosY, int(uiWidth), int(uiHeight))
            this.m_pcTDecBinIf.DecodeBin(&uiSigCoeffGroup, &baseCoeffGroupCtx[uiCtxSig])

            //this.DTRACE_CABAC_VL( g_nSymbolCounter++ );
            /*this.DTRACE_CABAC_T("\tuiSigCoeffGroup")
            this.DTRACE_CABAC_V(uiSigCoeffGroup)
            this.DTRACE_CABAC_T("\tuiCtxSig: ")
            this.DTRACE_CABAC_V(uiCtxSig)
            this.DTRACE_CABAC_T("\n")*/

            uiSigCoeffGroupFlag[iCGBlkPos] = uiSigCoeffGroup
        }

        // decode significant_coeff_flag
        patternSigCtx := TLibCommon.CalcPatternSigCtx(uiSigCoeffGroupFlag[:], iCGPosX, iCGPosY, int(uiWidth), int(uiHeight))
        var uiBlkPos, uiPosY, uiPosX, uiSig, uiCtxSig uint
        for ; iScanPosSig >= int(iSubPos); iScanPosSig-- {
            uiBlkPos = scan[iScanPosSig]
            uiPosY = uiBlkPos >> uint(uiLog2BlockSize)
            uiPosX = uiBlkPos - (uiPosY << uint(uiLog2BlockSize))
            uiSig = 0

            if uiSigCoeffGroupFlag[iCGBlkPos] != 0 {
                if iScanPosSig > int(iSubPos) || iSubSet == 0 || numNonZero != 0 {
                    uiCtxSig = uint(TLibCommon.GetSigCtxInc(patternSigCtx, uiScanIdx, int(uiPosX), int(uiPosY), int(uiLog2BlockSize), eTType))
                    this.m_pcTDecBinIf.DecodeBin(&uiSig, &baseCtx[uiCtxSig])

                    //this.DTRACE_CABAC_VL( g_nSymbolCounter++ );
                    /*this.DTRACE_CABAC_T("\tuiSig")
                    this.DTRACE_CABAC_V(uiSig)
                    this.DTRACE_CABAC_T("\tuiCtxSig: ")
                    this.DTRACE_CABAC_V(uiCtxSig)
                    this.DTRACE_CABAC_T("\n")*/
                } else {
                    uiSig = 1
                }
            }
            pcCoef[uiBlkPos] = TLibCommon.TCoeff(uiSig)
            if uiSig != 0 {
                pos[numNonZero] = int(uiBlkPos)
                numNonZero++
                if lastNZPosInCG == -1 {
                    lastNZPosInCG = iScanPosSig
                }
                firstNZPosInCG = iScanPosSig
            }
        }

        if numNonZero != 0 {
            signHidden := (lastNZPosInCG-firstNZPosInCG >= TLibCommon.SBH_THRESHOLD)
            absSum = 0

            var uiCtxSet uint
            if iSubSet > 0 && eTType == TLibCommon.TEXT_LUMA {
                uiCtxSet = 2
            } else {
                uiCtxSet = 0
            }

            var uiBin uint
            if c1 == 0 {
                uiCtxSet++
            }
            c1 = 1

            var baseCtxMod []TLibCommon.ContextModel

            if eTType == TLibCommon.TEXT_LUMA {
                baseCtxMod = this.m_cCUOneSCModel.Get2(0, 0)[4*uiCtxSet:]
            } else {
                baseCtxMod = this.m_cCUOneSCModel.Get2(0, 0)[TLibCommon.NUM_ONE_FLAG_CTX_LUMA+4*uiCtxSet:]
            }
            var absCoeff [TLibCommon.SCAN_SET_SIZE]int

            for i := int(0); i < numNonZero; i++ {
                absCoeff[i] = 1
            }

            numC1Flag := TLibCommon.MIN(numNonZero, TLibCommon.C1FLAG_NUMBER).(int)
            firstC2FlagIdx := -1

            for idx := int(0); idx < numC1Flag; idx++ {
                this.m_pcTDecBinIf.DecodeBin(&uiBin, &baseCtxMod[c1])
                //this.DTRACE_CABAC_VL( g_nSymbolCounter++ );
                /*this.DTRACE_CABAC_T("\tuiBin")
                this.DTRACE_CABAC_V(uiBin)
                this.DTRACE_CABAC_T("\tc1: ")
                this.DTRACE_CABAC_V(uint(c1))
                this.DTRACE_CABAC_T("\n")*/

                if uiBin == 1 {
                    c1 = 0
                    if firstC2FlagIdx == -1 {
                        firstC2FlagIdx = idx
                    }
                } else if (c1 < 3) && (c1 > 0) {
                    c1++
                }
                absCoeff[idx] = int(uiBin) + 1
            }

            if c1 == 0 {
                if eTType == TLibCommon.TEXT_LUMA {
                    baseCtxMod = this.m_cCUAbsSCModel.Get2(0, 0)[uiCtxSet:]
                } else {
                    baseCtxMod = this.m_cCUAbsSCModel.Get2(0, 0)[TLibCommon.NUM_ABS_FLAG_CTX_LUMA+uiCtxSet:]
                }

                if firstC2FlagIdx != -1 {
                    this.m_pcTDecBinIf.DecodeBin(&uiBin, &baseCtxMod[0])

                    //this.DTRACE_CABAC_VL( g_nSymbolCounter++ );
                    /*this.DTRACE_CABAC_T("\tuiBin")
                    this.DTRACE_CABAC_V(uiBin)
                    this.DTRACE_CABAC_T("\tc1: ")
                    this.DTRACE_CABAC_V(0)
                    this.DTRACE_CABAC_T("\n")*/

                    absCoeff[firstC2FlagIdx] = int(uiBin) + 2
                }
            }

            var coeffSigns uint
            if signHidden && beValid {
                this.m_pcTDecBinIf.DecodeBinsEP(&coeffSigns, numNonZero-1)
                //this.DTRACE_CABAC_VL( g_nSymbolCounter++ );
                /*this.DTRACE_CABAC_T("\tcoeffSigns")
                this.DTRACE_CABAC_V(coeffSigns)
                this.DTRACE_CABAC_T("\tnumNonZero-1: ")
                this.DTRACE_CABAC_V(uint(numNonZero - 1))
                this.DTRACE_CABAC_T("\n")*/

                coeffSigns <<= uint(32 - (numNonZero - 1))
            } else {
                this.m_pcTDecBinIf.DecodeBinsEP(&coeffSigns, numNonZero)
                //this.DTRACE_CABAC_VL( g_nSymbolCounter++ );
                /*this.DTRACE_CABAC_T("\tcoeffSigns")
                this.DTRACE_CABAC_V(coeffSigns)
                this.DTRACE_CABAC_T("\tnumNonZero: ")
                this.DTRACE_CABAC_V(uint(numNonZero))
                this.DTRACE_CABAC_T("\n")*/

                coeffSigns <<= uint(32 - numNonZero)
            }

            iFirstCoeff2 := int(1)
            if c1 == 0 || numNonZero > TLibCommon.C1FLAG_NUMBER {
                for idx := int(0); idx < numNonZero; idx++ {
                    var baseLevel uint
                    if idx < TLibCommon.C1FLAG_NUMBER {
                        baseLevel = uint(2 + iFirstCoeff2)
                    } else {
                        baseLevel = 1
                    }

                    if absCoeff[idx] == int(baseLevel) {
                        var uiLevel uint
                        this.xReadCoefRemainExGolomb(&uiLevel, uint(uiGoRiceParam))
                        //this.DTRACE_CABAC_VL( g_nSymbolCounter++ );
                        /*this.DTRACE_CABAC_T("\tuiLevel")
                        this.DTRACE_CABAC_V(uiLevel)
                        this.DTRACE_CABAC_T("\tuiGoRiceParam: ")
                        this.DTRACE_CABAC_V(uint(uiGoRiceParam))
                        this.DTRACE_CABAC_T("\n")*/

                        absCoeff[idx] = int(uiLevel + baseLevel)
                        if absCoeff[idx] > 3*(1<<uint(uiGoRiceParam)) {
                            uiGoRiceParam = TLibCommon.MIN(uiGoRiceParam+1, 4).(int)
                        }
                    }

                    if absCoeff[idx] >= 2 {
                        iFirstCoeff2 = 0
                    }
                }
            }

            for idx := int(0); idx < numNonZero; idx++ {
                blkPos := pos[idx]
                // Signs applied later.
                pcCoef[blkPos] = TLibCommon.TCoeff(absCoeff[idx])
                absSum += uint(absCoeff[idx])

                if idx == numNonZero-1 && signHidden && beValid {
                    // Infer sign of 1st element.
                    if (absSum & 0x1) != 0 {
                        pcCoef[blkPos] = -pcCoef[blkPos]
                    }
                } else {
                    sign := int(coeffSigns) >> 31
                    pcCoef[blkPos] = (pcCoef[blkPos] ^ TLibCommon.TCoeff(sign)) - TLibCommon.TCoeff(sign)
                    coeffSigns <<= 1
                }
            }
        }
    }

    return
}
func (this *TDecSbac) ParseTransformSkipFlags(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, width, height, uiDepth uint, eTType TLibCommon.TextType) {
    if pcCU.GetCUTransquantBypass1(uiAbsPartIdx) {
        return
    }
    if width != 4 || height != 4 {
        return
    }

    var useTransformSkip uint
    if eTType != 0 {
        this.m_pcTDecBinIf.DecodeBin(&useTransformSkip, this.m_cTransformSkipSCModel.Get3(0, TLibCommon.TEXT_CHROMA, 0))
    } else {
        this.m_pcTDecBinIf.DecodeBin(&useTransformSkip, this.m_cTransformSkipSCModel.Get3(0, TLibCommon.TEXT_LUMA, 0))
    }
    if eTType != TLibCommon.TEXT_LUMA {
        uiLog2TrafoSize := uint(TLibCommon.G_aucConvertToBit[pcCU.GetSlice().GetSPS().GetMaxCUWidth()]) + 2 - uiDepth
        if uiLog2TrafoSize == 2 {
            uiDepth--
        }
    }
    /*this.DTRACE_CABAC_VL( TLibCommon.GnSymbolCounter++ )*/
    /*this.DTRACE_CABAC_T("\tparseTransformSkip()")
    this.DTRACE_CABAC_T("\tsymbol=")
    this.DTRACE_CABAC_V(useTransformSkip)
    this.DTRACE_CABAC_T("\tAddr=")
    this.DTRACE_CABAC_V(pcCU.GetAddr())
    this.DTRACE_CABAC_T("\tetype=")
    this.DTRACE_CABAC_V(uint(eTType))
    this.DTRACE_CABAC_T("\tuiAbsPartIdx=")
    this.DTRACE_CABAC_V(uiAbsPartIdx)
    this.DTRACE_CABAC_T("\n")*/

    pcCU.SetTransformSkipSubParts4(useTransformSkip != 0, eTType, uiAbsPartIdx, uiDepth)
}

func (this *TDecSbac) UpdateContextTables(eSliceType TLibCommon.SliceType, iQp int) {
    var uiBit uint
    this.m_pcTDecBinIf.DecodeBinTrm(&uiBit)
    this.m_pcTDecBinIf.Finish()
    this.m_pcBitstream.ReadOutTrailingBits()
    
    this.m_cCUSplitFlagSCModel.InitBuffer(eSliceType, iQp, TLibCommon.INIT_SPLIT_FLAG[:])
    this.m_cCUSkipFlagSCModel.InitBuffer(eSliceType, iQp, TLibCommon.INIT_SKIP_FLAG[:])
    this.m_cCUMergeFlagExtSCModel.InitBuffer(eSliceType, iQp, TLibCommon.INIT_MERGE_FLAG_EXT[:])
    this.m_cCUMergeIdxExtSCModel.InitBuffer(eSliceType, iQp, TLibCommon.INIT_MERGE_IDX_EXT[:])
    this.m_cCUPartSizeSCModel.InitBuffer(eSliceType, iQp, TLibCommon.INIT_PART_SIZE[:])
    this.m_cCUAMPSCModel.InitBuffer(eSliceType, iQp, TLibCommon.INIT_CU_AMP_POS[:])
    this.m_cCUPredModeSCModel.InitBuffer(eSliceType, iQp, TLibCommon.INIT_PRED_MODE[:])
    this.m_cCUIntraPredSCModel.InitBuffer(eSliceType, iQp, TLibCommon.INIT_INTRA_PRED_MODE[:])
    this.m_cCUChromaPredSCModel.InitBuffer(eSliceType, iQp, TLibCommon.INIT_CHROMA_PRED_MODE[:])
    this.m_cCUInterDirSCModel.InitBuffer(eSliceType, iQp, TLibCommon.INIT_INTER_DIR[:])
    this.m_cCUMvdSCModel.InitBuffer(eSliceType, iQp, TLibCommon.INIT_MVD[:])
    this.m_cCURefPicSCModel.InitBuffer(eSliceType, iQp, TLibCommon.INIT_REF_PIC[:])
    this.m_cCUDeltaQpSCModel.InitBuffer(eSliceType, iQp, TLibCommon.INIT_DQP[:])
    this.m_cCUQtCbfSCModel.InitBuffer(eSliceType, iQp, TLibCommon.INIT_QT_CBF[:])
    this.m_cCUQtRootCbfSCModel.InitBuffer(eSliceType, iQp, TLibCommon.INIT_QT_ROOT_CBF[:])
    this.m_cCUSigCoeffGroupSCModel.InitBuffer(eSliceType, iQp, TLibCommon.INIT_SIG_CG_FLAG[:])
    this.m_cCUSigSCModel.InitBuffer(eSliceType, iQp, TLibCommon.INIT_SIG_FLAG[:])
    this.m_cCuCtxLastX.InitBuffer(eSliceType, iQp, TLibCommon.INIT_LAST[:])
    this.m_cCuCtxLastY.InitBuffer(eSliceType, iQp, TLibCommon.INIT_LAST[:])
    this.m_cCUOneSCModel.InitBuffer(eSliceType, iQp, TLibCommon.INIT_ONE_FLAG[:])
    this.m_cCUAbsSCModel.InitBuffer(eSliceType, iQp, TLibCommon.INIT_ABS_FLAG[:])
    this.m_cMVPIdxSCModel.InitBuffer(eSliceType, iQp, TLibCommon.INIT_MVP_IDX[:])
    this.m_cSaoMergeSCModel.InitBuffer(eSliceType, iQp, TLibCommon.INIT_SAO_MERGE_FLAG[:])
    this.m_cSaoTypeIdxSCModel.InitBuffer(eSliceType, iQp, TLibCommon.INIT_SAO_TYPE_IDX[:])
    this.m_cCUTransSubdivFlagSCModel.InitBuffer(eSliceType, iQp, TLibCommon.INIT_TRANS_SUBDIV_FLAG[:])
    this.m_cTransformSkipSCModel.InitBuffer(eSliceType, iQp, TLibCommon.INIT_TRANSFORMSKIP_FLAG[:])
    this.m_CUTransquantBypassFlagSCModel.InitBuffer(eSliceType, iQp, TLibCommon.INIT_CU_TRANSQUANT_BYPASS_FLAG[:])
    
    this.m_pcTDecBinIf.Start()
}

func (this *TDecSbac) ParseScalingList(scalingList *TLibCommon.TComScalingList) {
    //do nothing
}
