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
	"fmt"
    "gohm/TLibCommon"
)

type TEncBinIf interface {
    init(pcTComBitIf TLibCommon.TComBitIf)
    uninit()
    start()
    finish()
    copyState(pcTEncBinIf TEncBinIf)
    flush()
    resetBac()

    encodePCMAlignBits()
    xWritePCMCode(uiCode, uiLength uint)
    resetBits()
    getNumWrittenBits() uint

    encodeBin(uiBin uint, rcCtxModel *TLibCommon.ContextModel)
    encodeBinEP(uiBin uint)
    encodeBinsEP(uiBins uint, numBins int)
    encodeBinTrm(uiBin uint)

    getTEncBinCABAC() *TEncBinCABAC
    SetSbac(pEncSbac *TEncSbac)
}

type TEncBinCABAC struct { //: public TEncBinIf
	m_pTEncSbac 		*TEncSbac

    m_pcTComBitIf       TLibCommon.TComBitIf
    m_uiLow             uint
    m_uiRange           uint
    m_bufferedByte      uint
    m_numBufferedBytes  int
    m_bitsLeft          int
    m_uiBinsCoded       uint
    m_binCountIncrement int
    //#if FAST_BIT_EST
    m_fracBits uint64
    //#endif
}

//public:

func NewTEncBinCABAC() *TEncBinCABAC {
    return &TEncBinCABAC{}
}

func (this *TEncBinCABAC) SetSbac(pEncSbac *TEncSbac) {
    this.m_pTEncSbac = pEncSbac
}

func (this *TEncBinCABAC) init(pcTComBitIf TLibCommon.TComBitIf) {
    this.m_pcTComBitIf = pcTComBitIf
}
func (this *TEncBinCABAC) uninit() {
    this.m_pcTComBitIf = nil
}

func (this *TEncBinCABAC) start() {
    this.m_uiLow = 0
    this.m_uiRange = 510
    this.m_bitsLeft = 23
    this.m_numBufferedBytes = 0
    this.m_bufferedByte = 0xff
}

func (this *TEncBinCABAC) finish() {
    if (this.m_uiLow >> uint(32-this.m_bitsLeft)) != 0 {
        //assert( this.m_numBufferedBytes > 0 );
        //assert( this.m_bufferedByte != 0xff );
        this.m_pcTComBitIf.Write(this.m_bufferedByte+1, 8)
        for this.m_numBufferedBytes > 1 {
            this.m_pcTComBitIf.Write(0x00, 8)
            this.m_numBufferedBytes--
        }
        this.m_uiLow -= 1 << uint(32-this.m_bitsLeft)
    } else {
        if this.m_numBufferedBytes > 0 {
            this.m_pcTComBitIf.Write(this.m_bufferedByte, 8)
        }
        for this.m_numBufferedBytes > 1 {
            this.m_pcTComBitIf.Write(0xff, 8)
            this.m_numBufferedBytes--
        }
    }
    this.m_pcTComBitIf.Write(this.m_uiLow>>8, uint(24-this.m_bitsLeft))
}

func (this *TEncBinCABAC) copyState(pcTEncBinIf TEncBinIf) {
    pcTEncBinCABAC := pcTEncBinIf.getTEncBinCABAC()
    this.m_uiLow = pcTEncBinCABAC.m_uiLow
    this.m_uiRange = pcTEncBinCABAC.m_uiRange
    this.m_bitsLeft = pcTEncBinCABAC.m_bitsLeft
    this.m_bufferedByte = pcTEncBinCABAC.m_bufferedByte
    this.m_numBufferedBytes = pcTEncBinCABAC.m_numBufferedBytes
    //#if FAST_BIT_EST
    this.m_fracBits = pcTEncBinCABAC.m_fracBits
    //#endif
}

func (this *TEncBinCABAC) flush() {
    this.encodeBinTrm(1)
    this.finish()
    this.m_pcTComBitIf.Write(1, 1)
    this.m_pcTComBitIf.WriteAlignZero()

    this.start()
}

func (this *TEncBinCABAC) resetBac() {
    this.start()
}

func (this *TEncBinCABAC) encodePCMAlignBits() {
    this.finish()
    this.m_pcTComBitIf.Write(1, 1)
    this.m_pcTComBitIf.WriteAlignZero() // pcm align zero
}

func (this *TEncBinCABAC) xWritePCMCode(uiCode, uiLength uint) {
    this.m_pcTComBitIf.Write(uiCode, uiLength)
}

func (this *TEncBinCABAC) resetBits() {
    this.m_uiLow = 0
    this.m_bitsLeft = 23
    this.m_numBufferedBytes = 0
    this.m_bufferedByte = 0xff
    if this.m_binCountIncrement != 0 {
        this.m_uiBinsCoded = 0
    }
    //#if FAST_BIT_EST
    this.m_fracBits &= 32767
    //#endif
}

func (this *TEncBinCABAC) getNumWrittenBits() uint {
	//fmt.Printf("TEncBinCABAC:%d,%d,%d\n",this.m_pcTComBitIf.GetNumberOfWrittenBits(),this.m_numBufferedBytes,this.m_bitsLeft);
	
    return this.m_pcTComBitIf.GetNumberOfWrittenBits() + 8*uint(this.m_numBufferedBytes) + 23 - uint(this.m_bitsLeft)
}

func (this *TEncBinCABAC) encodeBin(binValue uint, rcCtxModel *TLibCommon.ContextModel) {
    /*DTRACE_CABAC_VL( g_nSymbolCounter++ )*/
    /*this.DTRACE_CABAC_T( "\tstate=" )
    this.DTRACE_CABAC_V( ( rcCtxModel.GetState() << 1 ) + rcCtxModel.GetMps() )
    this.DTRACE_CABAC_T( "\tsymbol=" )
    this.DTRACE_CABAC_V( binValue )
    this.DTRACE_CABAC_T( "\n" )*/
    
    this.m_uiBinsCoded += uint(this.m_binCountIncrement)
    rcCtxModel.SetBinsCoded(1)

    uiLPS := uint(TLibCommon.TComCABACTables_sm_aucLPSTable[rcCtxModel.GetState()][(this.m_uiRange>>6)&3])
    this.m_uiRange -= uint(uiLPS)

    if binValue != uint(rcCtxModel.GetMps()) {
        numBits := uint(TLibCommon.TComCABACTables_sm_aucRenormTable[uiLPS>>3])
        this.m_uiLow = (this.m_uiLow + this.m_uiRange) << numBits
        this.m_uiRange = uiLPS << numBits
        rcCtxModel.UpdateLPS()

        this.m_bitsLeft -= int(numBits)
    } else {
        rcCtxModel.UpdateMPS()
        if this.m_uiRange >= 256 {
            return
        }

        this.m_uiLow <<= 1
        this.m_uiRange <<= 1
        this.m_bitsLeft--
    }

    this.testAndWriteOut()
}

func (this *TEncBinCABAC) encodeBinEP(binValue uint) {
    /*DTRACE_CABAC_VL( g_nSymbolCounter++ )*/
    /*this.DTRACE_CABAC_T( "\tEPsymbol=" )
    this.DTRACE_CABAC_V( binValue )
    this.DTRACE_CABAC_T( "\n" )*/
  
    this.m_uiBinsCoded += uint(this.m_binCountIncrement)
    this.m_uiLow <<= 1
    if binValue != 0 {
        this.m_uiLow += this.m_uiRange
    }
    this.m_bitsLeft--

    this.testAndWriteOut()
}

func (this *TEncBinCABAC) encodeBinsEP(binValues uint, numBins int) {
    this.m_uiBinsCoded += uint(numBins & (-this.m_binCountIncrement))

    for i := 0; i < numBins; i++ {
        /*DTRACE_CABAC_VL( g_nSymbolCounter++ )*/
        /*this.DTRACE_CABAC_T( "\tEPsymbol=" )
        this.DTRACE_CABAC_V( ( binValues >> ( numBins - 1 - i ) ) & 1 )
        this.DTRACE_CABAC_T( "\n" )*/
    }

    for numBins > 8 {
        numBins -= 8
        pattern := binValues >> uint(numBins)
        this.m_uiLow <<= 8
        this.m_uiLow += this.m_uiRange * pattern
        binValues -= pattern << uint(numBins)
        this.m_bitsLeft -= 8

        this.testAndWriteOut()
    }

    this.m_uiLow <<= uint(numBins)
    this.m_uiLow += this.m_uiRange * binValues
    this.m_bitsLeft -= numBins

    this.testAndWriteOut()
}

func (this *TEncBinCABAC) encodeBinTrm(binValue uint) {
    this.m_uiBinsCoded += uint(this.m_binCountIncrement)
    this.m_uiRange -= 2
    if binValue != 0 {
        this.m_uiLow += this.m_uiRange
        this.m_uiLow <<= 7
        this.m_uiRange = 2 << 7
        this.m_bitsLeft -= 7
    } else if this.m_uiRange >= 256 {
        return
    } else {
        this.m_uiLow <<= 1
        this.m_uiRange <<= 1
        this.m_bitsLeft--
    }

    this.testAndWriteOut()
}

func (this *TEncBinCABAC) getTEncBinCABAC() *TEncBinCABAC { return this }

func (this *TEncBinCABAC) setBinsCoded(uiVal uint) { this.m_uiBinsCoded = uiVal }
func (this *TEncBinCABAC) getBinsCoded() uint      { return this.m_uiBinsCoded }
func (this *TEncBinCABAC) setBinCountingEnableFlag(bFlag bool) {
    this.m_binCountIncrement = int(TLibCommon.B2U(bFlag))
}
func (this *TEncBinCABAC) getBinCountingEnableFlag() bool { return this.m_binCountIncrement != 0 }

func (this *TEncBinCABAC) testAndWriteOut() {
    if this.m_bitsLeft < 12 {
        this.writeOut()
    }
}

func (this *TEncBinCABAC) writeOut() {
    leadByte := this.m_uiLow >> uint(24-this.m_bitsLeft)
    this.m_bitsLeft += 8
    this.m_uiLow &= (0xffffffff >> uint(this.m_bitsLeft))

    if leadByte == 0xff {
        this.m_numBufferedBytes++
    } else {
        if this.m_numBufferedBytes > 0 {
            carry := leadByte >> 8
            byte1 := this.m_bufferedByte + carry
            this.m_bufferedByte = leadByte & 0xff
            this.m_pcTComBitIf.Write(byte1, 8)

            byte1 = (0xff + carry) & 0xff
            for this.m_numBufferedBytes > 1 {
                this.m_pcTComBitIf.Write(byte1, 8)
                this.m_numBufferedBytes--
            }
        } else {
            this.m_numBufferedBytes = 1
            this.m_bufferedByte = leadByte
        }
    }
}

type TEncBinCABACCounter struct {
    TEncBinCABAC
}

func NewTEncBinCABACCounter() *TEncBinCABACCounter {
    return &TEncBinCABACCounter{}
}

func (this *TEncBinCABACCounter) finish() {
    this.m_pcTComBitIf.Write(0, uint(this.m_fracBits>>15))
    this.m_fracBits &= 32767
}

func (this *TEncBinCABACCounter) getNumWrittenBits() uint {
	//fmt.Printf("TEncBinCABACCounter:%d,%d\n",this.m_pcTComBitIf.GetNumberOfWrittenBits(),this.m_fracBits);
	
    return this.m_pcTComBitIf.GetNumberOfWrittenBits() + uint(this.m_fracBits>>15)
}

func (this *TEncBinCABACCounter) encodeBin(binValue uint, rcCtxModel *TLibCommon.ContextModel) {
    /*DTRACE_CABAC_VL( g_nSymbolCounter++ )*/
    /*this.m_pTEncSbac.DTRACE_CABAC_T( "\tstate=" )
    this.m_pTEncSbac.DTRACE_CABAC_V( uint(rcCtxModel.GetState()) )
    this.m_pTEncSbac.DTRACE_CABAC_T( "\tsymbol=" )
    this.m_pTEncSbac.DTRACE_CABAC_V( binValue )
    this.m_pTEncSbac.DTRACE_CABAC_T( "\n" )*/
    
    this.m_uiBinsCoded += uint(this.m_binCountIncrement)

    this.m_fracBits += uint64(rcCtxModel.GetEntropyBits(int16(binValue)))
    rcCtxModel.Update(int(binValue))
}

func (this *TEncBinCABACCounter) encodeBinEP(binValue uint) {
    this.m_uiBinsCoded += uint(this.m_binCountIncrement)
    this.m_fracBits += 32768
}

func (this *TEncBinCABACCounter) encodeBinsEP(binValues uint, numBins int) {
    this.m_uiBinsCoded += uint(numBins & (-this.m_binCountIncrement))
    this.m_fracBits += 32768 * uint64(numBins)
}

func (this *TEncBinCABACCounter) encodeBinTrm(binValue uint) {
    this.m_uiBinsCoded += uint(this.m_binCountIncrement)
    this.m_fracBits += uint64(TLibCommon.ContextModel_GetEntropyBitsTrm(int(binValue)))
}

// ====================================================================================================================
// Class definition
// ====================================================================================================================

/// SBAC encoder class
type TEncSbac struct { //: public TEncEntropyIf
	m_pTraceFile  io.Writer

    m_pcBitIf TLibCommon.TComBitIf
    m_pcSlice *TLibCommon.TComSlice
    m_pcBinIf TEncBinIf
    //SBAC RD
    m_uiCoeffCost uint
    m_uiLastQp    uint

    m_contextModels                 [TLibCommon.MAX_NUM_CTX_MOD]TLibCommon.ContextModel
    m_numContextModels              int
    m_cCUSplitFlagSCModel           *TLibCommon.ContextModel3DBuffer
    m_cCUSkipFlagSCModel            *TLibCommon.ContextModel3DBuffer
    m_cCUMergeFlagExtSCModel        *TLibCommon.ContextModel3DBuffer
    m_cCUMergeIdxExtSCModel         *TLibCommon.ContextModel3DBuffer
    m_cCUPartSizeSCModel            *TLibCommon.ContextModel3DBuffer
    m_cCUPredModeSCModel            *TLibCommon.ContextModel3DBuffer
    m_cCUIntraPredSCModel           *TLibCommon.ContextModel3DBuffer
    m_cCUChromaPredSCModel          *TLibCommon.ContextModel3DBuffer
    m_cCUDeltaQpSCModel             *TLibCommon.ContextModel3DBuffer
    m_cCUInterDirSCModel            *TLibCommon.ContextModel3DBuffer
    m_cCURefPicSCModel              *TLibCommon.ContextModel3DBuffer
    m_cCUMvdSCModel                 *TLibCommon.ContextModel3DBuffer
    m_cCUQtCbfSCModel               *TLibCommon.ContextModel3DBuffer
    m_cCUTransSubdivFlagSCModel     *TLibCommon.ContextModel3DBuffer
    m_cCUQtRootCbfSCModel           *TLibCommon.ContextModel3DBuffer
    m_cCUSigCoeffGroupSCModel       *TLibCommon.ContextModel3DBuffer
    m_cCUSigSCModel                 *TLibCommon.ContextModel3DBuffer
    m_cCuCtxLastX                   *TLibCommon.ContextModel3DBuffer
    m_cCuCtxLastY                   *TLibCommon.ContextModel3DBuffer
    m_cCUOneSCModel                 *TLibCommon.ContextModel3DBuffer
    m_cCUAbsSCModel                 *TLibCommon.ContextModel3DBuffer
    m_cMVPIdxSCModel                *TLibCommon.ContextModel3DBuffer
    m_cCUAMPSCModel                 *TLibCommon.ContextModel3DBuffer
    m_cSaoMergeSCModel              *TLibCommon.ContextModel3DBuffer
    m_cSaoTypeIdxSCModel            *TLibCommon.ContextModel3DBuffer
    m_cTransformSkipSCModel         *TLibCommon.ContextModel3DBuffer
    m_CUTransquantBypassFlagSCModel *TLibCommon.ContextModel3DBuffer
}

func (this *TEncSbac) SetTraceFile(traceFile io.Writer) {
    this.m_pTraceFile = traceFile
}

func (this *TEncSbac) GetTraceFile() io.Writer {
    return this.m_pTraceFile
}


func (this *TEncSbac) XTraceLCUHeader(traceLevel uint) {
    if this.GetTraceFile() != nil && (traceLevel&TLibCommon.TRACE_LEVEL) != 0 {
        io.WriteString(this.m_pTraceFile, "========= LCU Parameter Set ===============================================\n") //, pLCU.GetAddr());
    }
}

func (this *TEncSbac) XTraceCUHeader(traceLevel uint) {
    if this.GetTraceFile() != nil && (traceLevel&TLibCommon.TRACE_LEVEL) != 0 {
        io.WriteString(this.m_pTraceFile, "========= CU Parameter Set ================================================\n") //, pCU.GetCUPelX(), pCU.GetCUPelY());
    }
}

func (this *TEncSbac) XTracePUHeader(traceLevel uint) {
    if this.GetTraceFile() != nil && (traceLevel&TLibCommon.TRACE_LEVEL) != 0 {
        io.WriteString(this.m_pTraceFile, "========= PU Parameter Set ================================================\n") //, pCU.GetCUPelX(), pCU.GetCUPelY());
    }
}

func (this *TEncSbac) XTraceTUHeader(traceLevel uint) {
    if this.GetTraceFile() != nil && (traceLevel&TLibCommon.TRACE_LEVEL) != 0 {
        io.WriteString(this.m_pTraceFile, "========= TU Parameter Set ================================================\n") //, pCU.GetCUPelX(), pCU.GetCUPelY());
    }
}

func (this *TEncSbac) XTraceCoefHeader(traceLevel uint) {
    if this.GetTraceFile() != nil && (traceLevel&TLibCommon.TRACE_LEVEL) != 0 {
        io.WriteString(this.m_pTraceFile, "========= Coefficient Parameter Set =======================================\n") //, pCU.GetCUPelX(), pCU.GetCUPelY());
    }
}

func (this *TEncSbac) XTraceResiHeader(traceLevel uint) {
    if this.GetTraceFile() != nil && (traceLevel&TLibCommon.TRACE_LEVEL) != 0 {
        io.WriteString(this.m_pTraceFile, "========= Residual Parameter Set ==========================================\n") //, pCU.GetCUPelX(), pCU.GetCUPelY());
    }
}

func (this *TEncSbac) XTracePredHeader(traceLevel uint) {
    if this.GetTraceFile() != nil && (traceLevel&TLibCommon.TRACE_LEVEL) != 0 {
        io.WriteString(this.m_pTraceFile, "========= Prediction Parameter Set ========================================\n") //, pCU.GetCUPelX(), pCU.GetCUPelY());
    }
}

func (this *TEncSbac) XTraceRecoHeader(traceLevel uint) {
    if this.GetTraceFile() != nil && (traceLevel&TLibCommon.TRACE_LEVEL) != 0 {
        io.WriteString(this.m_pTraceFile, "========= Reconstruction Parameter Set ====================================\n") //, pCU.GetCUPelX(), pCU.GetCUPelY());
    }
}

func (this *TEncSbac) XReadAeTr(Value int, pSymbolName string, traceLevel uint) {
    if this.GetTraceFile() != nil && (traceLevel&TLibCommon.TRACE_LEVEL) != 0 {
        //fprintf( g_hTrace, "%8lld  ", g_nSymbolCounter++ );
        io.WriteString(this.m_pTraceFile, fmt.Sprintf("%-62s ae(v) : %4d\n", pSymbolName, Value))
        //fflush ( g_hTrace );
    }
}

func (this *TEncSbac) XReadCeofTr(pCoeff []TLibCommon.TCoeff, uiWidth, traceLevel uint) {
    if this.GetTraceFile() != nil && (traceLevel&TLibCommon.TRACE_LEVEL) != 0 {
        for i := uint(0); i < uiWidth; i++ {
            io.WriteString(this.m_pTraceFile, fmt.Sprintf("%04x ", uint16(pCoeff[i])))
        }
        io.WriteString(this.m_pTraceFile, "\n")
    }
}

func (this *TEncSbac) XReadResiTr(pPel []TLibCommon.Pel, uiWidth, traceLevel uint) {
    if this.GetTraceFile() != nil && (traceLevel&TLibCommon.TRACE_LEVEL) != 0 {
        for i := uint(0); i < uiWidth; i++ {
            io.WriteString(this.m_pTraceFile, fmt.Sprintf("%04x ", uint16(pPel[i])))
        }
        io.WriteString(this.m_pTraceFile, "\n")
    }
}

func (this *TEncSbac) XReadPredTr(pPel []TLibCommon.Pel, uiWidth, traceLevel uint) {
    if this.GetTraceFile() != nil && (traceLevel&TLibCommon.TRACE_LEVEL) != 0 {
        for i := uint(0); i < uiWidth; i++ {
            io.WriteString(this.m_pTraceFile, fmt.Sprintf("%02x ", TLibCommon.Pxl(pPel[i])))
        }
        io.WriteString(this.m_pTraceFile, "\n")
    }
}

func (this *TEncSbac) XReadRecoTr(pPel []TLibCommon.Pel, uiWidth, traceLevel uint) {
    if this.GetTraceFile() != nil && (traceLevel&TLibCommon.TRACE_LEVEL) != 0 {
        for i := uint(0); i < uiWidth; i++ {
            io.WriteString(this.m_pTraceFile, fmt.Sprintf("%02x ", TLibCommon.Pxl(pPel[i])))
        }
        io.WriteString(this.m_pTraceFile, "\n")
    }
}

func (this *TEncSbac) DTRACE_CABAC_F(x float32) {
    if this.GetTraceFile() != nil && TLibCommon.TRACE_CABAC {
        //fmt.Printf("%f", x)
        io.WriteString(this.m_pTraceFile, fmt.Sprintf("%f", x))
    }
}
func (this *TEncSbac) DTRACE_CABAC_V(x uint) {
    if this.GetTraceFile() != nil && TLibCommon.TRACE_CABAC {
        //fmt.Printf ("%d", x )
        io.WriteString(this.m_pTraceFile, fmt.Sprintf("%x", x))
    }
}
func (this *TEncSbac) DTRACE_CABAC_VL(x uint) {
    if this.GetTraceFile() != nil && TLibCommon.TRACE_CABAC {
        //fmt.Printf ("%lld", x )
        io.WriteString(this.m_pTraceFile, fmt.Sprintf("%lld", x))
    }
}
func (this *TEncSbac) DTRACE_CABAC_T(x string) {
    if this.GetTraceFile() != nil && TLibCommon.TRACE_CABAC {
        //fmt.Printf ("%s", x )
        io.WriteString(this.m_pTraceFile, fmt.Sprintf("%s", x))
    }
}
func (this *TEncSbac) DTRACE_CABAC_X(x uint) {
    if this.GetTraceFile() != nil && TLibCommon.TRACE_CABAC {
        //fmt.Printf ("%x", x )
        io.WriteString(this.m_pTraceFile, fmt.Sprintf("%x", x))
    }
}
func (this *TEncSbac) DTRACE_CABAC_N() {
    if this.GetTraceFile() != nil && TLibCommon.TRACE_CABAC {
        //fmt.Printf ("\n" )
        io.WriteString(this.m_pTraceFile, "\n")
    }
}

/*func (this *TEncSbac) DTRACE_CABAC_R(x,y) {
	if this.GetTraceFile()!=nil {
		io.WriteString(this.m_pTraceFile, fmt.Sprintf (x,    y ));
	}
}*/

func NewTEncSbac() *TEncSbac {
    pTEncSbac := &TEncSbac{m_pcBitIf: nil, m_pcSlice: nil, m_pcBinIf: nil, m_uiCoeffCost: 0, m_numContextModels: 0}
    pTEncSbac.xInit()

    return pTEncSbac
}

func (this *TEncSbac) xInit() {
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
func (this *TEncSbac) init(p TEncBinIf) { this.m_pcBinIf = p }
func (this *TEncSbac) uninit()          { this.m_pcBinIf = nil }

//  Virtual list
func (this *TEncSbac) resetEntropy() {
    iQp := this.m_pcSlice.GetSliceQp()
    eSliceType := this.m_pcSlice.GetSliceType()

    encCABACTableIdx := this.m_pcSlice.GetPPS().GetEncCABACTableIdx()
    if !this.m_pcSlice.IsIntra() && (encCABACTableIdx == TLibCommon.B_SLICE || encCABACTableIdx == TLibCommon.P_SLICE) && this.m_pcSlice.GetPPS().GetCabacInitPresentFlag() {
        eSliceType = TLibCommon.SliceType(encCABACTableIdx)
    }

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
    this.m_cCUTransSubdivFlagSCModel.InitBuffer(eSliceType, iQp, TLibCommon.INIT_TRANS_SUBDIV_FLAG[:])
    this.m_cSaoMergeSCModel.InitBuffer(eSliceType, iQp, TLibCommon.INIT_SAO_MERGE_FLAG[:])
    this.m_cSaoTypeIdxSCModel.InitBuffer(eSliceType, iQp, TLibCommon.INIT_SAO_TYPE_IDX[:])
    this.m_cTransformSkipSCModel.InitBuffer(eSliceType, iQp, TLibCommon.INIT_TRANSFORMSKIP_FLAG[:])
    this.m_CUTransquantBypassFlagSCModel.InitBuffer(eSliceType, iQp, TLibCommon.INIT_CU_TRANSQUANT_BYPASS_FLAG[:])
    // new structure
    this.m_uiLastQp = uint(iQp)

    this.m_pcBinIf.start()

    return
}
func (this *TEncSbac) determineCabacInitIdx() {
    qp := this.m_pcSlice.GetSliceQp()

    if !this.m_pcSlice.IsIntra() {
        var aSliceTypeChoices = []TLibCommon.SliceType{TLibCommon.B_SLICE, TLibCommon.P_SLICE}

        bestCost := uint(TLibCommon.MAX_UINT)
        bestSliceType := aSliceTypeChoices[0]
        for idx := uint(0); idx < 2; idx++ {
            curCost := uint(0)
            curSliceType := aSliceTypeChoices[idx]

            curCost = this.m_cCUSplitFlagSCModel.CalcCost(curSliceType, qp, TLibCommon.INIT_SPLIT_FLAG[:])
            curCost += this.m_cCUSkipFlagSCModel.CalcCost(curSliceType, qp, TLibCommon.INIT_SKIP_FLAG[:])
            curCost += this.m_cCUMergeFlagExtSCModel.CalcCost(curSliceType, qp, TLibCommon.INIT_MERGE_FLAG_EXT[:])
            curCost += this.m_cCUMergeIdxExtSCModel.CalcCost(curSliceType, qp, TLibCommon.INIT_MERGE_IDX_EXT[:])
            curCost += this.m_cCUPartSizeSCModel.CalcCost(curSliceType, qp, TLibCommon.INIT_PART_SIZE[:])
            curCost += this.m_cCUAMPSCModel.CalcCost(curSliceType, qp, TLibCommon.INIT_CU_AMP_POS[:])
            curCost += this.m_cCUPredModeSCModel.CalcCost(curSliceType, qp, TLibCommon.INIT_PRED_MODE[:])
            curCost += this.m_cCUIntraPredSCModel.CalcCost(curSliceType, qp, TLibCommon.INIT_INTRA_PRED_MODE[:])
            curCost += this.m_cCUChromaPredSCModel.CalcCost(curSliceType, qp, TLibCommon.INIT_CHROMA_PRED_MODE[:])
            curCost += this.m_cCUInterDirSCModel.CalcCost(curSliceType, qp, TLibCommon.INIT_INTER_DIR[:])
            curCost += this.m_cCUMvdSCModel.CalcCost(curSliceType, qp, TLibCommon.INIT_MVD[:])
            curCost += this.m_cCURefPicSCModel.CalcCost(curSliceType, qp, TLibCommon.INIT_REF_PIC[:])
            curCost += this.m_cCUDeltaQpSCModel.CalcCost(curSliceType, qp, TLibCommon.INIT_DQP[:])
            curCost += this.m_cCUQtCbfSCModel.CalcCost(curSliceType, qp, TLibCommon.INIT_QT_CBF[:])
            curCost += this.m_cCUQtRootCbfSCModel.CalcCost(curSliceType, qp, TLibCommon.INIT_QT_ROOT_CBF[:])
            curCost += this.m_cCUSigCoeffGroupSCModel.CalcCost(curSliceType, qp, TLibCommon.INIT_SIG_CG_FLAG[:])
            curCost += this.m_cCUSigSCModel.CalcCost(curSliceType, qp, TLibCommon.INIT_SIG_FLAG[:])
            curCost += this.m_cCuCtxLastX.CalcCost(curSliceType, qp, TLibCommon.INIT_LAST[:])
            curCost += this.m_cCuCtxLastY.CalcCost(curSliceType, qp, TLibCommon.INIT_LAST[:])
            curCost += this.m_cCUOneSCModel.CalcCost(curSliceType, qp, TLibCommon.INIT_ONE_FLAG[:])
            curCost += this.m_cCUAbsSCModel.CalcCost(curSliceType, qp, TLibCommon.INIT_ABS_FLAG[:])
            curCost += this.m_cMVPIdxSCModel.CalcCost(curSliceType, qp, TLibCommon.INIT_MVP_IDX[:])
            curCost += this.m_cCUTransSubdivFlagSCModel.CalcCost(curSliceType, qp, TLibCommon.INIT_TRANS_SUBDIV_FLAG[:])
            curCost += this.m_cSaoMergeSCModel.CalcCost(curSliceType, qp, TLibCommon.INIT_SAO_MERGE_FLAG[:])
            curCost += this.m_cSaoTypeIdxSCModel.CalcCost(curSliceType, qp, TLibCommon.INIT_SAO_TYPE_IDX[:])
            curCost += this.m_cTransformSkipSCModel.CalcCost(curSliceType, qp, TLibCommon.INIT_TRANSFORMSKIP_FLAG[:])
            curCost += this.m_CUTransquantBypassFlagSCModel.CalcCost(curSliceType, qp, TLibCommon.INIT_CU_TRANSQUANT_BYPASS_FLAG[:])
            if curCost < bestCost {
                bestSliceType = curSliceType
                bestCost = curCost
            }
        }
        this.m_pcSlice.GetPPS().SetEncCABACTableIdx(uint(bestSliceType))
    } else {
        this.m_pcSlice.GetPPS().SetEncCABACTableIdx(TLibCommon.I_SLICE)
    }
}

func (this *TEncSbac) setBitstream(p TLibCommon.TComBitIf) {
    this.m_pcBitIf = p
    this.m_pcBinIf.init(p)
}
func (this *TEncSbac) setTraceFile(traceFile io.Writer) {
    this.m_pTraceFile = traceFile
}
func (this *TEncSbac) setSlice(p *TLibCommon.TComSlice) { this.m_pcSlice = p }

// SBAC RD
func (this *TEncSbac) resetCoeffCost()    { this.m_uiCoeffCost = 0 }
func (this *TEncSbac) getCoeffCost() uint { return this.m_uiCoeffCost }

func (this *TEncSbac) load(pSrc *TEncSbac) { this.xCopyFrom(pSrc) }
func (this *TEncSbac) loadIntraDirModeLuma(pSrc *TEncSbac) {
    this.m_pcBinIf.copyState(pSrc.m_pcBinIf)

    this.m_cCUIntraPredSCModel.CopyFrom(pSrc.m_cCUIntraPredSCModel)
}

func (this *TEncSbac) store(pDest *TEncSbac)       { pDest.xCopyFrom(this) }
func (this *TEncSbac) loadContexts(pScr *TEncSbac) { this.xCopyContextsFrom(pScr) }
func (this *TEncSbac) resetBits() {
    this.m_pcBinIf.resetBits()
    this.m_pcBitIf.ResetBits()
}
func (this *TEncSbac) getNumberOfWrittenBits() uint { return this.m_pcBinIf.getNumWrittenBits() }

//--SBAC RD

func (this *TEncSbac) codeVPS(pcVPS *TLibCommon.TComVPS)                    {}
func (this *TEncSbac) codeSPS(pcSPS *TLibCommon.TComSPS)                    {}
func (this *TEncSbac) codePPS(pcPPS *TLibCommon.TComPPS)                    {}
func (this *TEncSbac) codeSliceHeader(pcSlice *TLibCommon.TComSlice)        {}
func (this *TEncSbac) codeTilesWPPEntryPoint(pcSlice *TLibCommon.TComSlice) {}
func (this *TEncSbac) codeTerminatingBit(uilsLast uint) {
    this.m_pcBinIf.encodeBinTrm(uilsLast)
}
func (this *TEncSbac) codeSliceFinish() { this.m_pcBinIf.finish() }
func (this *TEncSbac) codeSaoMaxUvlc(code, maxSymbol uint) {
    if maxSymbol == 0 {
        return
    }

    var i int
    bCodeLast := (maxSymbol > code)

    if code == 0 {
        this.m_pcBinIf.encodeBinEP(0)
    } else {
        this.m_pcBinIf.encodeBinEP(1)
        for i = 0; i < int(code)-1; i++ {
            this.m_pcBinIf.encodeBinEP(1)
        }
        if bCodeLast {
            this.m_pcBinIf.encodeBinEP(0)
        }
    }
}
func (this *TEncSbac) codeSaoMerge(uiCode uint) {
    if uiCode == 0 {
        this.m_pcBinIf.encodeBin(0, this.m_cSaoMergeSCModel.Get3(0, 0, 0))
    } else {
        this.m_pcBinIf.encodeBin(1, this.m_cSaoMergeSCModel.Get3(0, 0, 0))
    }
}

func (this *TEncSbac) codeSaoTypeIdx(uiCode uint) {
    if uiCode == 0 {
        this.m_pcBinIf.encodeBin(0, this.m_cSaoTypeIdxSCModel.Get3(0, 0, 0))
    } else {
        this.m_pcBinIf.encodeBin(1, this.m_cSaoTypeIdxSCModel.Get3(0, 0, 0))
        this.m_pcBinIf.encodeBinEP(uint(TLibCommon.B2U(uiCode <= 4)))
    }
}
func (this *TEncSbac) codeSaoUflc(uiLength, uiCode uint) {
    this.m_pcBinIf.encodeBinsEP(uiCode, int(uiLength))
}
func (this *TEncSbac) codeSAOSign(uiCode uint)                                 { this.m_pcBinIf.encodeBinEP(uiCode) }
func (this *TEncSbac) codeScalingList(scalingList *TLibCommon.TComScalingList) {}
func (this *TEncSbac) xWriteUnarySymbol(uiSymbol uint, pcSCModel []TLibCommon.ContextModel, iOffset int) {
    if uiSymbol != 0 {
        this.m_pcBinIf.encodeBin(1, &pcSCModel[0])
    } else {
        this.m_pcBinIf.encodeBin(0, &pcSCModel[0])
    }

    if 0 == uiSymbol {
        return
    }

    for uiSymbol != 0 {
        uiSymbol--

        if uiSymbol != 0 {
            this.m_pcBinIf.encodeBin(1, &pcSCModel[iOffset])
        } else {
            this.m_pcBinIf.encodeBin(0, &pcSCModel[iOffset])
        }
    }

    return
}

func (this *TEncSbac) xWriteUnaryMaxSymbol(uiSymbol uint, pcSCModel []TLibCommon.ContextModel, iOffset int, uiMaxSymbol uint) {
    if uiMaxSymbol == 0 {
        return
    }
    if uiSymbol != 0 {
        this.m_pcBinIf.encodeBin(1, &pcSCModel[0])
    } else {
        this.m_pcBinIf.encodeBin(0, &pcSCModel[0])
    }

    if uiSymbol == 0 {
        return
    }

    bCodeLast := (uiMaxSymbol > uiSymbol)

    uiSymbol--
    for uiSymbol != 0 {
        this.m_pcBinIf.encodeBin(1, &pcSCModel[iOffset])
        uiSymbol--
    }
    if bCodeLast {
        this.m_pcBinIf.encodeBin(0, &pcSCModel[iOffset])
    }

    return
}

func (this *TEncSbac) xWriteEpExGolomb(uiSymbol, uiCount uint) {
    bins := uint(0)
    numBins := 0

    for uiSymbol >= uint(1<<uiCount) {
        bins = 2*bins + 1
        numBins++
        uiSymbol -= 1 << uiCount
        uiCount++
    }
    bins = 2*bins + 0
    numBins++

    bins = (bins << uiCount) | uiSymbol
    numBins += int(uiCount)

    //assert( numBins <= 32 );
    this.m_pcBinIf.encodeBinsEP(bins, numBins)
}

func (this *TEncSbac) xWriteCoefRemainExGolomb(symbol uint, rParam uint) {
    codeNumber := int(symbol)
    var length uint
    if codeNumber < (TLibCommon.COEF_REMAIN_BIN_REDUCTION << rParam) {
        length = uint(codeNumber) >> rParam
        this.m_pcBinIf.encodeBinsEP((1<<(length+1))-2, int(length+1))
        this.m_pcBinIf.encodeBinsEP(uint(codeNumber%(1<<rParam)), int(rParam))
    } else {
        length = rParam
        codeNumber = codeNumber - (TLibCommon.COEF_REMAIN_BIN_REDUCTION << rParam)
        for codeNumber >= (1 << length) {
            codeNumber -= (1 << (length))
            length++
        }
        this.m_pcBinIf.encodeBinsEP((1<<(TLibCommon.COEF_REMAIN_BIN_REDUCTION+length+1-rParam))-2, int(TLibCommon.COEF_REMAIN_BIN_REDUCTION+length+1-rParam))
        this.m_pcBinIf.encodeBinsEP(uint(codeNumber), int(length))
    }
}

func (this *TEncSbac) xWriteTerminatingBit(uiBit uint) {}

func (this *TEncSbac) xCopyFrom(pSrc *TEncSbac) {
    this.m_pcBinIf.copyState(pSrc.m_pcBinIf)

    this.m_uiCoeffCost = pSrc.m_uiCoeffCost
    this.m_uiLastQp = pSrc.m_uiLastQp

    for i := 0; i < this.m_numContextModels; i++ {
        this.m_contextModels[i] = pSrc.m_contextModels[i]
    }
    //memcpy( this.m_contextModels, pSrc->m_contextModels, this.m_numContextModels * sizeof( ContextModel ) );
}

func (this *TEncSbac) xCopyContextsFrom(pSrc *TEncSbac) {
    for i := 0; i < this.m_numContextModels; i++ {
        this.m_contextModels[i] = pSrc.m_contextModels[i]
    }
    //memcpy(this.m_contextModels, pSrc->m_contextModels, this.m_numContextModels*sizeof(this.m_contextModels[0]));
}

func (this *TEncSbac) codeDFFlag(uiCode uint, pSymbolName string) {}
func (this *TEncSbac) codeDFSvlc(iCode int, pSymbolName string)   {}

func (this *TEncSbac) codeCUTransquantBypassFlag(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint) {
    uiSymbol := uint(TLibCommon.B2U(pcCU.GetCUTransquantBypass1(uiAbsPartIdx)))
    this.m_pcBinIf.encodeBin(uiSymbol, this.m_CUTransquantBypassFlagSCModel.Get3(0, 0, 0))
}

func (this *TEncSbac) codeSkipFlag(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint) {
    // get context function is here
    uiSymbol := uint(TLibCommon.B2U(pcCU.IsSkipped(uiAbsPartIdx)))
    uiCtxSkip := pcCU.GetCtxSkipFlag(uiAbsPartIdx)
    this.m_pcBinIf.encodeBin(uiSymbol, this.m_cCUSkipFlagSCModel.Get3(0, 0, uiCtxSkip))
    /*DTRACE_CABAC_VL( g_nSymbolCounter++ );*/
    this.DTRACE_CABAC_T( "\tSkipFlag" );
    this.DTRACE_CABAC_T( "\tuiCtxSkip: ");
    this.DTRACE_CABAC_V( uiCtxSkip );
    this.DTRACE_CABAC_T( "\tuiSymbol: ");
    this.DTRACE_CABAC_V( uiSymbol );
    this.DTRACE_CABAC_T( "\n");
}

func (this *TEncSbac) codeMergeFlag(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint) {
    uiSymbol := uint(TLibCommon.B2U(pcCU.GetMergeFlag1(uiAbsPartIdx)))
    this.m_pcBinIf.encodeBin(uiSymbol, this.m_cCUMergeFlagExtSCModel.Get3(0, 0, 0))

    /*DTRACE_CABAC_VL( g_nSymbolCounter++ );*/
    this.DTRACE_CABAC_T( "\tMergeFlag: " );
    this.DTRACE_CABAC_V( uiSymbol );
    this.DTRACE_CABAC_T( "\tAddress: " );
    this.DTRACE_CABAC_V( pcCU.GetAddr() );
    this.DTRACE_CABAC_T( "\tuiAbsPartIdx: " );
    this.DTRACE_CABAC_V( uiAbsPartIdx );
    this.DTRACE_CABAC_T( "\n" );
}

func (this *TEncSbac) codeMergeIndex(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint) {
    uiUnaryIdx := uint(pcCU.GetMergeIndex1(uiAbsPartIdx))
    uiNumCand := pcCU.GetSlice().GetMaxNumMergeCand()
    if uiNumCand > 1 {
        for ui := uint(0); ui < uiNumCand-1; ui++ {
            var uiSymbol uint
            if ui == uiUnaryIdx {
                uiSymbol = 0
            } else {
                uiSymbol = 1
            }

            if ui == 0 {
                this.m_pcBinIf.encodeBin(uiSymbol, this.m_cCUMergeIdxExtSCModel.Get3(0, 0, 0))
            } else {
                this.m_pcBinIf.encodeBinEP(uiSymbol)
            }
            if uiSymbol == 0 {
                break
            }
        }
    }
    /*DTRACE_CABAC_VL( g_nSymbolCounter++ );*/
    this.DTRACE_CABAC_T( "\tparseMergeIndex()" );
    this.DTRACE_CABAC_T( "\tuiMRGIdx= " );
    this.DTRACE_CABAC_V( uint(pcCU.GetMergeIndex1( uiAbsPartIdx )) );
    this.DTRACE_CABAC_T( "\n" );
}

func (this *TEncSbac) codeSplitFlag(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint, uiDepth uint) {
    if uiDepth == pcCU.GetSlice().GetSPS().GetMaxCUDepth()-pcCU.GetSlice().GetSPS().GetAddCUDepth() {
        return
    }

    uiCtx := uint(pcCU.GetCtxSplitFlag(uiAbsPartIdx, uiDepth))
    uiCurrSplitFlag := uint(TLibCommon.B2U(uint(pcCU.GetDepth1(uiAbsPartIdx)) > uiDepth))

    //assert( uiCtx < 3 );
    this.m_pcBinIf.encodeBin(uiCurrSplitFlag, this.m_cCUSplitFlagSCModel.Get3(0, 0, uiCtx))
    //DTRACE_CABAC_VL( g_nSymbolCounter++ )
    this.DTRACE_CABAC_T( "\tSplitFlag\n" )
    return
}

func (this *TEncSbac) codeMVPIdx(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint, eRefList TLibCommon.RefPicList) {
    iSymbol := int(pcCU.GetMVPIdx2(eRefList, uiAbsPartIdx))
    iNum := TLibCommon.AMVP_MAX_NUM_CANDS

    this.xWriteUnaryMaxSymbol(uint(iSymbol), this.m_cMVPIdxSCModel.Get1(0), 1, uint(iNum-1))
}

func (this *TEncSbac) codePartSize(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint, uiDepth uint) {
    eSize := pcCU.GetPartitionSize1(uiAbsPartIdx)
    if pcCU.IsIntra(uiAbsPartIdx) {
        if uiDepth == pcCU.GetSlice().GetSPS().GetMaxCUDepth()-pcCU.GetSlice().GetSPS().GetAddCUDepth() {
            this.m_pcBinIf.encodeBin(uint(TLibCommon.B2U(eSize == TLibCommon.SIZE_2Nx2N)), this.m_cCUPartSizeSCModel.Get3(0, 0, 0))
        }
        return
    }

    switch eSize {
    case TLibCommon.SIZE_2Nx2N:
        this.m_pcBinIf.encodeBin(1, this.m_cCUPartSizeSCModel.Get3(0, 0, 0))
    case TLibCommon.SIZE_2NxN:
        fallthrough
    case TLibCommon.SIZE_2NxnU:
        fallthrough
    case TLibCommon.SIZE_2NxnD:
        this.m_pcBinIf.encodeBin(0, this.m_cCUPartSizeSCModel.Get3(0, 0, 0))
        this.m_pcBinIf.encodeBin(1, this.m_cCUPartSizeSCModel.Get3(0, 0, 1))
        if pcCU.GetSlice().GetSPS().GetAMPAcc(uiDepth) != 0 {
            if eSize == TLibCommon.SIZE_2NxN {
                this.m_pcBinIf.encodeBin(1, this.m_cCUAMPSCModel.Get3(0, 0, 0))
            } else {
                this.m_pcBinIf.encodeBin(0, this.m_cCUAMPSCModel.Get3(0, 0, 0))
                if eSize == TLibCommon.SIZE_2NxnU {
                    this.m_pcBinIf.encodeBinEP(0)
                } else {
                    this.m_pcBinIf.encodeBinEP(1)
                }
            }
        }
    case TLibCommon.SIZE_Nx2N:
        fallthrough
    case TLibCommon.SIZE_nLx2N:
        fallthrough
    case TLibCommon.SIZE_nRx2N:
        this.m_pcBinIf.encodeBin(0, this.m_cCUPartSizeSCModel.Get3(0, 0, 0))
        this.m_pcBinIf.encodeBin(0, this.m_cCUPartSizeSCModel.Get3(0, 0, 1))
        if uiDepth == pcCU.GetSlice().GetSPS().GetMaxCUDepth()-pcCU.GetSlice().GetSPS().GetAddCUDepth() && !(pcCU.GetWidth1(uiAbsPartIdx) == 8 && pcCU.GetHeight1(uiAbsPartIdx) == 8) {
            this.m_pcBinIf.encodeBin(1, this.m_cCUPartSizeSCModel.Get3(0, 0, 2))
        }
        if pcCU.GetSlice().GetSPS().GetAMPAcc(uiDepth) != 0 {
            if eSize == TLibCommon.SIZE_Nx2N {
                this.m_pcBinIf.encodeBin(1, this.m_cCUAMPSCModel.Get3(0, 0, 0))
            } else {
                this.m_pcBinIf.encodeBin(0, this.m_cCUAMPSCModel.Get3(0, 0, 0))

                if eSize == TLibCommon.SIZE_nLx2N {
                    this.m_pcBinIf.encodeBinEP(0)
                } else {
                    this.m_pcBinIf.encodeBinEP(1)
                }
            }
        }
    case TLibCommon.SIZE_NxN:
        if uiDepth == pcCU.GetSlice().GetSPS().GetMaxCUDepth()-pcCU.GetSlice().GetSPS().GetAddCUDepth() && !(pcCU.GetWidth1(uiAbsPartIdx) == 8 && pcCU.GetHeight1(uiAbsPartIdx) == 8) {
            this.m_pcBinIf.encodeBin(0, this.m_cCUPartSizeSCModel.Get3(0, 0, 0))
            this.m_pcBinIf.encodeBin(0, this.m_cCUPartSizeSCModel.Get3(0, 0, 1))
            this.m_pcBinIf.encodeBin(0, this.m_cCUPartSizeSCModel.Get3(0, 0, 2))
        }
    default:
        {
            //  assert(0);
        }
    }
}

func (this *TEncSbac) codePredMode(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint) {
    // get context function is here
    iPredMode := pcCU.GetPredictionMode1(uiAbsPartIdx)
    if iPredMode == TLibCommon.MODE_INTER {
        this.m_pcBinIf.encodeBin(0, this.m_cCUPredModeSCModel.Get3(0, 0, 0))
    } else {
        this.m_pcBinIf.encodeBin(1, this.m_cCUPredModeSCModel.Get3(0, 0, 0))
    }
}

func (this *TEncSbac) codeIPCMInfo(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint) {
    uiIPCM := uint(TLibCommon.B2U(pcCU.GetIPCMFlag1(uiAbsPartIdx)))

    writePCMSampleFlag := pcCU.GetIPCMFlag1(uiAbsPartIdx)

    this.m_pcBinIf.encodeBinTrm(uiIPCM)

    if writePCMSampleFlag {
        this.m_pcBinIf.encodePCMAlignBits()

        uiMinCoeffSize := pcCU.GetPic().GetMinCUWidth() * pcCU.GetPic().GetMinCUHeight()
        uiLumaOffset := uiMinCoeffSize * uiAbsPartIdx
        uiChromaOffset := uiLumaOffset >> 2
        var piPCMSample []TLibCommon.Pel
        var uiWidth, uiHeight, uiSampleBits, uiX, uiY uint

        piPCMSample = pcCU.GetPCMSampleY()[uiLumaOffset:]
        uiWidth = uint(pcCU.GetWidth1(uiAbsPartIdx))
        uiHeight = uint(pcCU.GetHeight1(uiAbsPartIdx))
        uiSampleBits = pcCU.GetSlice().GetSPS().GetPCMBitDepthLuma()

        for uiY = 0; uiY < uiHeight; uiY++ {
            for uiX = 0; uiX < uiWidth; uiX++ {
                uiSample := uint(piPCMSample[uiX])

                this.m_pcBinIf.xWritePCMCode(uiSample, uiSampleBits)
            }
            piPCMSample = piPCMSample[uiWidth:]
        }

        piPCMSample = pcCU.GetPCMSampleCb()[uiChromaOffset:]
        uiWidth = uint(pcCU.GetWidth1(uiAbsPartIdx) / 2)
        uiHeight = uint(pcCU.GetHeight1(uiAbsPartIdx) / 2)
        uiSampleBits = pcCU.GetSlice().GetSPS().GetPCMBitDepthChroma()

        for uiY = 0; uiY < uiHeight; uiY++ {
            for uiX = 0; uiX < uiWidth; uiX++ {
                uiSample := uint(piPCMSample[uiX])

                this.m_pcBinIf.xWritePCMCode(uiSample, uiSampleBits)
            }
            piPCMSample = piPCMSample[uiWidth:]
        }

        piPCMSample = pcCU.GetPCMSampleCr()[uiChromaOffset:]
        uiWidth = uint(pcCU.GetWidth1(uiAbsPartIdx) / 2)
        uiHeight = uint(pcCU.GetHeight1(uiAbsPartIdx) / 2)
        uiSampleBits = pcCU.GetSlice().GetSPS().GetPCMBitDepthChroma()

        for uiY = 0; uiY < uiHeight; uiY++ {
            for uiX = 0; uiX < uiWidth; uiX++ {
                uiSample := uint(piPCMSample[uiX])

                this.m_pcBinIf.xWritePCMCode(uiSample, uiSampleBits)
            }
            piPCMSample = piPCMSample[uiWidth:]
        }
        this.m_pcBinIf.resetBac()
    }
}
func (this *TEncSbac) codeTransformSubdivFlag(uiSymbol, uiCtx uint) {
    this.m_pcBinIf.encodeBin(uiSymbol, this.m_cCUTransSubdivFlagSCModel.Get3(0, 0, uiCtx))
    /*DTRACE_CABAC_VL( g_nSymbolCounter++ )*/
    this.DTRACE_CABAC_T( "\tparseTransformSubdivFlag()" )
    this.DTRACE_CABAC_T( "\tsymbol=" )
    this.DTRACE_CABAC_V( uiSymbol )
    this.DTRACE_CABAC_T( "\tctx=" )
    this.DTRACE_CABAC_V( uiCtx )
    this.DTRACE_CABAC_T( "\n" )
}

func (this *TEncSbac) codeQtCbf(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint, eType TLibCommon.TextType, uiTrDepth uint) {
    uiCbf := uint(pcCU.GetCbf3(uiAbsPartIdx, eType, uiTrDepth))
    uiCtx := pcCU.GetCtxQtCbf(eType, uiTrDepth)
    if eType != 0 {
        this.m_pcBinIf.encodeBin(uiCbf, this.m_cCUQtCbfSCModel.Get3(0, TLibCommon.TEXT_CHROMA, uiCtx))
    } else {
        this.m_pcBinIf.encodeBin(uiCbf, this.m_cCUQtCbfSCModel.Get3(0, uint(eType), uiCtx))
    }
    /*DTRACE_CABAC_VL( g_nSymbolCounter++ )*/
    this.DTRACE_CABAC_T( "\tparseQtCbf()" )
    this.DTRACE_CABAC_T( "\tsymbol=" )
    this.DTRACE_CABAC_V( uiCbf )
    this.DTRACE_CABAC_T( "\tctx=" )
    this.DTRACE_CABAC_V( uiCtx )
    this.DTRACE_CABAC_T( "\tetype=" )
    this.DTRACE_CABAC_V( uint(eType) )
    this.DTRACE_CABAC_T( "\tuiAbsPartIdx=" )
    this.DTRACE_CABAC_V( uiAbsPartIdx )
    this.DTRACE_CABAC_T( "\n" )
}

func (this *TEncSbac) codeQtRootCbf(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint) {
    uiCbf := uint(TLibCommon.B2U(pcCU.GetQtRootCbf(uiAbsPartIdx)))
    uiCtx := uint(0)
    this.m_pcBinIf.encodeBin(uiCbf, this.m_cCUQtRootCbfSCModel.Get3(0, 0, uiCtx))
    /*  DTRACE_CABAC_VL( g_nSymbolCounter++ )*/
    this.DTRACE_CABAC_T( "\tparseQtRootCbf()" )
    this.DTRACE_CABAC_T( "\tsymbol=" )
    this.DTRACE_CABAC_V( uiCbf )
    this.DTRACE_CABAC_T( "\tctx=" )
    this.DTRACE_CABAC_V( uiCtx )
    this.DTRACE_CABAC_T( "\tuiAbsPartIdx=" )
    this.DTRACE_CABAC_V( uiAbsPartIdx )
    this.DTRACE_CABAC_T( "\n" )
}

func (this *TEncSbac) codeQtCbfZero(pcCU *TLibCommon.TComDataCU, eType TLibCommon.TextType, uiTrDepth uint) {
    // this function is only used to estimate the bits when cbf is 0
    // and will never be called when writing the bistream. do not need to write log
    uiCbf := uint(0)
    uiCtx := pcCU.GetCtxQtCbf(eType, uiTrDepth)
    if eType != 0 {
        this.m_pcBinIf.encodeBin(uiCbf, this.m_cCUQtCbfSCModel.Get3(0, TLibCommon.TEXT_CHROMA, uiCtx))
    } else {
        this.m_pcBinIf.encodeBin(uiCbf, this.m_cCUQtCbfSCModel.Get3(0, uint(eType), uiCtx))
    }
}

func (this *TEncSbac) codeQtRootCbfZero(pcCU *TLibCommon.TComDataCU) {
    // this function is only used to estimate the bits when cbf is 0
    // and will never be called when writing the bistream. do not need to write log
    uiCbf := uint(0)
    uiCtx := uint(0)
    this.m_pcBinIf.encodeBin(uiCbf, this.m_cCUQtRootCbfSCModel.Get3(0, 0, uiCtx))
}

func (this *TEncSbac) codeIntraDirLumaAng(pcCU *TLibCommon.TComDataCU, absPartIdx uint, isMultiple bool) {
    var dir [4]uint
    var j uint
    var preds = [4][3]int{{-1, -1, -1}, {-1, -1, -1}, {-1, -1, -1}, {-1, -1, -1}}
    var predNum [4]int
    var predIdx = [4]int{-1, -1, -1, -1}
    mode := pcCU.GetPartitionSize1(absPartIdx)
    var partNum uint
    if isMultiple {
        if mode == TLibCommon.SIZE_NxN {
            partNum = 4
        } else {
            partNum = 1
        }
    } else {
        partNum = 1
    }

    partOffset := (pcCU.GetPic().GetNumPartInCU() >> (pcCU.GetDepth1(absPartIdx) << 1)) >> 2
    for j = 0; j < partNum; j++ {
        dir[j] = uint(pcCU.GetLumaIntraDir1(absPartIdx + partOffset*j))
        predNum[j] = pcCU.GetIntraDirLumaPredictor(absPartIdx+partOffset*j, preds[j][:], nil)
        for i := int(0); i < predNum[j]; i++ {
            if dir[j] == uint(preds[j][i]) {
                predIdx[j] = i
            }
        }
        if predIdx[j] != -1 {
            this.m_pcBinIf.encodeBin(1, this.m_cCUIntraPredSCModel.Get3(0, 0, 0))
        } else {
            this.m_pcBinIf.encodeBin(0, this.m_cCUIntraPredSCModel.Get3(0, 0, 0))
        }
    }
    for j = 0; j < partNum; j++ {
        if predIdx[j] != -1 {
            this.m_pcBinIf.encodeBinEP(uint(TLibCommon.B2U(predIdx[j] != 0)))
            if predIdx[j] != 0 {
                this.m_pcBinIf.encodeBinEP(uint(predIdx[j] - 1))
            }
        } else {
            if preds[j][0] > preds[j][1] {
                tmp := preds[j][0]
                preds[j][0] = preds[j][1]
                preds[j][1] = tmp
                //std::swap(preds[j][0], preds[j][1]);
            }
            if preds[j][0] > preds[j][2] {
                tmp := preds[j][0]
                preds[j][0] = preds[j][2]
                preds[j][2] = tmp
                //std::swap(preds[j][0], preds[j][2]);
            }
            if preds[j][1] > preds[j][2] {
                tmp := preds[j][1]
                preds[j][1] = preds[j][2]
                preds[j][2] = tmp
                //std::swap(preds[j][1], preds[j][2]);
            }
            for i := int(predNum[j] - 1); i >= 0; i-- {
                if dir[j] > uint(preds[j][i]) {
                    dir[j] = dir[j] - 1
                } else {
                    dir[j] = dir[j]
                }
            }
            this.m_pcBinIf.encodeBinsEP(dir[j], 5)
        }
    }
    return
}

func (this *TEncSbac) codeIntraDirChroma(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint) {
    uiIntraDirChroma := uint(pcCU.GetChromaIntraDir1(uiAbsPartIdx))

    if uiIntraDirChroma == TLibCommon.DM_CHROMA_IDX {
        this.m_pcBinIf.encodeBin(0, this.m_cCUChromaPredSCModel.Get3(0, 0, 0))
    } else {
        var uiAllowedChromaDir [TLibCommon.NUM_CHROMA_MODE]uint
        pcCU.GetAllowedChromaDir(uiAbsPartIdx, uiAllowedChromaDir[:])

        for i := uint(0); i < TLibCommon.NUM_CHROMA_MODE-1; i++ {
            if uiIntraDirChroma == uiAllowedChromaDir[i] {
                uiIntraDirChroma = i
                break
            }
        }
        this.m_pcBinIf.encodeBin(1, this.m_cCUChromaPredSCModel.Get3(0, 0, 0))

        this.m_pcBinIf.encodeBinsEP(uiIntraDirChroma, 2)
    }
    return
}

func (this *TEncSbac) codeInterDir(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint) {
    uiInterDir := uint(pcCU.GetInterDir1(uiAbsPartIdx) - 1)
    uiCtx := pcCU.GetCtxInterDir(uiAbsPartIdx)
    pCtx := this.m_cCUInterDirSCModel.Get1(0)
    if pcCU.GetPartitionSize1(uiAbsPartIdx) == TLibCommon.SIZE_2Nx2N || pcCU.GetHeight1(uiAbsPartIdx) != 8 {
        if uiInterDir == 2 {
            this.m_pcBinIf.encodeBin(1, &pCtx[uiCtx])
        } else {
            this.m_pcBinIf.encodeBin(0, &pCtx[uiCtx])
        }
    }
    if uiInterDir < 2 {
        this.m_pcBinIf.encodeBin(uiInterDir, &pCtx[4])
    }
    return
}

func (this *TEncSbac) codeRefFrmIdx(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint, eRefList TLibCommon.RefPicList) {
    {
        iRefFrame := pcCU.GetCUMvField(eRefList).GetRefIdx(int(uiAbsPartIdx))
        pCtx := this.m_cCURefPicSCModel.Get1(0)
        if iRefFrame == 0 {
            this.m_pcBinIf.encodeBin(0, &pCtx[0])
        } else {
            this.m_pcBinIf.encodeBin(1, &pCtx[0])
        }

        if iRefFrame > 0 {
            uiRefNum := uint(pcCU.GetSlice().GetNumRefIdx(eRefList) - 2)
            pCtx = pCtx[1:]
            iRefFrame--
            for ui := uint(0); ui < uiRefNum; ui++ {
                var uiSymbol uint
                if ui == uint(iRefFrame) {
                    uiSymbol = 0
                } else {
                    uiSymbol = 1
                }
                if ui == 0 {
                    this.m_pcBinIf.encodeBin(uiSymbol, &pCtx[0])
                } else {
                    this.m_pcBinIf.encodeBinEP(uiSymbol)
                }
                if uiSymbol == 0 {
                    break
                }
            }
        }
    }
    return
}

func (this *TEncSbac) codeMvd(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint, eRefList TLibCommon.RefPicList) {
    if pcCU.GetSlice().GetMvdL1ZeroFlag() && eRefList == TLibCommon.REF_PIC_LIST_1 && pcCU.GetInterDir1(uiAbsPartIdx) == 3 {
        return
    }

    pcCUMvField := pcCU.GetCUMvField(eRefList)
    mvd := pcCUMvField.GetMvd(int(uiAbsPartIdx))
    iHor := int(mvd.GetHor())
    iVer := int(mvd.GetVer())
    pCtx := this.m_cCUMvdSCModel.Get1(0)

    this.m_pcBinIf.encodeBin(uint(TLibCommon.B2U(iHor != 0)), &pCtx[0])
    this.m_pcBinIf.encodeBin(uint(TLibCommon.B2U(iVer != 0)), &pCtx[0])

    bHorAbsGr0 := iHor != 0
    bVerAbsGr0 := iVer != 0
    uiHorAbs := TLibCommon.ABS(iHor).(int) // ? -iHor : iHor;
    uiVerAbs := TLibCommon.ABS(iVer).(int) // ? -iVer : iVer;
    pCtx = pCtx[1:]

    if bHorAbsGr0 {
        this.m_pcBinIf.encodeBin(uint(TLibCommon.B2U(uiHorAbs > 1)), &pCtx[0])
    }

    if bVerAbsGr0 {
        this.m_pcBinIf.encodeBin(uint(TLibCommon.B2U(uiVerAbs > 1)), &pCtx[0])
    }

    if bHorAbsGr0 {
        if uiHorAbs > 1 {
            this.xWriteEpExGolomb(uint(uiHorAbs-2), 1)
        }

        this.m_pcBinIf.encodeBinEP(uint(TLibCommon.B2U(0 > iHor)))
    }

    if bVerAbsGr0 {
        if uiVerAbs > 1 {
            this.xWriteEpExGolomb(uint(uiVerAbs-2), 1)
        }

        this.m_pcBinIf.encodeBinEP(uint(TLibCommon.B2U(0 > iVer)))
    }

    return
}

func (this *TEncSbac) codeDeltaQP(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint) {
    iDQp := int(pcCU.GetQP1(uiAbsPartIdx) - pcCU.GetRefQP(uiAbsPartIdx))

    qpBdOffsetY := pcCU.GetSlice().GetSPS().GetQpBDOffsetY()
    iDQp = (iDQp+78+qpBdOffsetY+(qpBdOffsetY/2))%(52+qpBdOffsetY) - 26 - (qpBdOffsetY / 2)

    uiAbsDQp := TLibCommon.ABS(iDQp).(int) // > 0)? iDQp  : (-iDQp));
    TUValue := TLibCommon.MIN(int(uiAbsDQp), int(TLibCommon.CU_DQP_TU_CMAX)).(int)
    this.xWriteUnaryMaxSymbol(uint(TUValue), this.m_cCUDeltaQpSCModel.Get1(0), 1, TLibCommon.CU_DQP_TU_CMAX)
    if uiAbsDQp >= TLibCommon.CU_DQP_TU_CMAX {
        this.xWriteEpExGolomb(uint(uiAbsDQp-TLibCommon.CU_DQP_TU_CMAX), TLibCommon.CU_DQP_EG_k)
    }

    if uiAbsDQp > 0 {
        var uiSign uint
        if iDQp > 0 {
            uiSign = 0
        } else {
            uiSign = 1
        }
        this.m_pcBinIf.encodeBinEP(uiSign)
    }

    return
}

func (this *TEncSbac) codeLastSignificantXY(uiPosX, uiPosY uint, width, height int, eTType TLibCommon.TextType, uiScanIdx uint) {
    // swap
    if uiScanIdx == TLibCommon.SCAN_VER {
        tmp := uiPosX
        uiPosX = uiPosY
        uiPosY = tmp
        //swap( uiPosX, uiPosY );
    }

    var uiCtxLast uint
    pCtxX := this.m_cCuCtxLastX.Get2(0, uint(eTType))
    pCtxY := this.m_cCuCtxLastY.Get2(0, uint(eTType))
    uiGroupIdxX := TLibCommon.G_uiGroupIdx[uiPosX]
    uiGroupIdxY := TLibCommon.G_uiGroupIdx[uiPosY]

    var blkSizeOffsetX, blkSizeOffsetY int
    var shiftX, shiftY uint
    if eTType != 0 {
        blkSizeOffsetX = 0
        blkSizeOffsetY = 0
        shiftX = uint(TLibCommon.G_aucConvertToBit[width])
        shiftY = uint(TLibCommon.G_aucConvertToBit[height])
    } else {
        blkSizeOffsetX = int(TLibCommon.G_aucConvertToBit[width]*3 + ((TLibCommon.G_aucConvertToBit[width] + 1) >> 2))
        blkSizeOffsetY = int(TLibCommon.G_aucConvertToBit[height]*3 + ((TLibCommon.G_aucConvertToBit[height] + 1) >> 2))
        shiftX = uint((TLibCommon.G_aucConvertToBit[width] + 3) >> 2)
        shiftY = uint((TLibCommon.G_aucConvertToBit[height] + 3) >> 2)

    }
    // posX
    for uiCtxLast = 0; uiCtxLast < uiGroupIdxX; uiCtxLast++ {
        this.m_pcBinIf.encodeBin(1, &pCtxX[blkSizeOffsetX+int(uiCtxLast>>shiftX)])
        //DTRACE_CABAC_VL( g_nSymbolCounter++ )
    	/*this.DTRACE_CABAC_T( "\txsymbol=" )
    	this.DTRACE_CABAC_V( 1 )
    	this.DTRACE_CABAC_T( "\tctx=" )
    	this.DTRACE_CABAC_V( uint(blkSizeOffsetX + int(uiCtxLast >>shiftX)) )
    	this.DTRACE_CABAC_T("\tuiBits: ")
        this.DTRACE_CABAC_V(this.getNumberOfWrittenBits())
    	this.DTRACE_CABAC_T( "\n" )*/
    }
    if uiGroupIdxX < TLibCommon.G_uiGroupIdx[width-1] {
        this.m_pcBinIf.encodeBin(0, &pCtxX[blkSizeOffsetX+int(uiCtxLast>>shiftX)])
        //DTRACE_CABAC_VL( g_nSymbolCounter++ )
        /*this.DTRACE_CABAC_T( "\tsymbol=" )
    	this.DTRACE_CABAC_V( 0 )
    	this.DTRACE_CABAC_T( "\tctx=" )
    	this.DTRACE_CABAC_V( uint(blkSizeOffsetX + int(uiCtxLast >>shiftX)) )
    	this.DTRACE_CABAC_T("\tuiBits: ")
        this.DTRACE_CABAC_V(this.getNumberOfWrittenBits())
    	this.DTRACE_CABAC_T( "\n" )*/
    }

    // posY
    for uiCtxLast = 0; uiCtxLast < uiGroupIdxY; uiCtxLast++ {
        this.m_pcBinIf.encodeBin(1, &pCtxY[blkSizeOffsetY+int(uiCtxLast>>shiftY)])
        //DTRACE_CABAC_VL( g_nSymbolCounter++ )
        /*this.DTRACE_CABAC_T( "\tysymbol=" )
    	this.DTRACE_CABAC_V( 1 )
    	this.DTRACE_CABAC_T( "\tctx=" )
    	this.DTRACE_CABAC_V( uint(blkSizeOffsetY+int(uiCtxLast>>shiftY)) )
    	this.DTRACE_CABAC_T("\tuiBits: ")
        this.DTRACE_CABAC_V(this.getNumberOfWrittenBits())
    	this.DTRACE_CABAC_T( "\n" )*/
    }
    if uiGroupIdxY < TLibCommon.G_uiGroupIdx[height-1] {
        this.m_pcBinIf.encodeBin(0, &pCtxY[blkSizeOffsetY+int(uiCtxLast>>shiftY)])
        //DTRACE_CABAC_VL( g_nSymbolCounter++ )
        /*this.DTRACE_CABAC_T( "\tsymbol=" )
    	this.DTRACE_CABAC_V( 0 )
    	this.DTRACE_CABAC_T( "\tctx=" )
    	this.DTRACE_CABAC_V( uint(blkSizeOffsetY+int(uiCtxLast>>shiftY)) )
    	this.DTRACE_CABAC_T("\tuiBits: ")
        this.DTRACE_CABAC_V(this.getNumberOfWrittenBits())
    	this.DTRACE_CABAC_T( "\n" )*/
    }
    if uiGroupIdxX > 3 {
        uiCount := (uiGroupIdxX - 2) >> 1
        uiPosX = uiPosX - TLibCommon.G_uiMinInGroup[uiGroupIdxX]
        for i := int(uiCount) - 1; i >= 0; i-- {
            this.m_pcBinIf.encodeBinEP((uiPosX >> uint(i)) & 1)
            //DTRACE_CABAC_VL( g_nSymbolCounter++ )
		    /*this.DTRACE_CABAC_T( "\tuiPosX=" )
		    this.DTRACE_CABAC_V( uint( uiPosX >> uint(i) ) & 1 )
		    this.DTRACE_CABAC_T("\tuiBits: ")
       		this.DTRACE_CABAC_V(this.getNumberOfWrittenBits())
		    this.DTRACE_CABAC_T( "\n" )*/
        }
    }
    if uiGroupIdxY > 3 {
        uiCount := (uiGroupIdxY - 2) >> 1
        uiPosY = uiPosY - TLibCommon.G_uiMinInGroup[uiGroupIdxY]
        for i := int(uiCount) - 1; i >= 0; i-- {
            this.m_pcBinIf.encodeBinEP((uiPosY >> uint(i)) & 1)
            //DTRACE_CABAC_VL( g_nSymbolCounter++ )
      		/*this.DTRACE_CABAC_T( "\tuiPosY=" )
      		this.DTRACE_CABAC_V( uint( uiPosY >> uint(i) ) & 1  )
      		this.DTRACE_CABAC_T("\tuiBits: ")
        	this.DTRACE_CABAC_V(this.getNumberOfWrittenBits())
      		this.DTRACE_CABAC_T( "\n" )*/
        }
    }
}

func (this *TEncSbac) codeCoeffNxN(pcCU *TLibCommon.TComDataCU, pcCoef []TLibCommon.TCoeff, uiAbsPartIdx, uiWidth, uiHeight, uiDepth uint, eTType TLibCommon.TextType) {
    /*DTRACE_CABAC_VL( g_nSymbolCounter++ )*/
    this.DTRACE_CABAC_T( "\tparseCoeffNxN()\teType=" )
    this.DTRACE_CABAC_V( uint(eTType) )
    this.DTRACE_CABAC_T( "\twidth=" )
    this.DTRACE_CABAC_V( uiWidth )
    this.DTRACE_CABAC_T( "\theight=" )
    this.DTRACE_CABAC_V( uiHeight )
    this.DTRACE_CABAC_T( "\tdepth=" )
    this.DTRACE_CABAC_V( uiDepth )
    this.DTRACE_CABAC_T( "\tabspartidx=" )
    this.DTRACE_CABAC_V( uiAbsPartIdx )
    this.DTRACE_CABAC_T( "\ttoCU-X=" )
    this.DTRACE_CABAC_V( pcCU.GetCUPelX() )
    this.DTRACE_CABAC_T( "\ttoCU-Y=" )
    this.DTRACE_CABAC_V( pcCU.GetCUPelY() )
    this.DTRACE_CABAC_T( "\tCU-addr=" )
    this.DTRACE_CABAC_V(  pcCU.GetAddr() )
    this.DTRACE_CABAC_T( "\tinCU-X=" )
    this.DTRACE_CABAC_V( TLibCommon.G_auiRasterToPelX[ TLibCommon.G_auiZscanToRaster[uiAbsPartIdx] ] )
    this.DTRACE_CABAC_T( "\tinCU-Y=" )
    this.DTRACE_CABAC_V( TLibCommon.G_auiRasterToPelY[ TLibCommon.G_auiZscanToRaster[uiAbsPartIdx] ] )
    this.DTRACE_CABAC_T( "\tpredmode=" )
    this.DTRACE_CABAC_V( uint(pcCU.GetPredictionMode1( uiAbsPartIdx )) )
    this.DTRACE_CABAC_T("\tuiBits: ")
    this.DTRACE_CABAC_V(this.getNumberOfWrittenBits())
    this.DTRACE_CABAC_T( "\n" )

    if uiWidth > this.m_pcSlice.GetSPS().GetMaxTrSize() {
        uiWidth = this.m_pcSlice.GetSPS().GetMaxTrSize()
        uiHeight = this.m_pcSlice.GetSPS().GetMaxTrSize()
    }

    uiNumSig := uint(0)

    // compute number of significant coefficients
    uiNumSig = uint(TEncEntropy_countNonZeroCoeffs(pcCoef, uiWidth*uiHeight))

    if uiNumSig == 0 {
        return
    }
    if pcCU.GetSlice().GetPPS().GetUseTransformSkip() {
        this.codeTransformSkipFlags(pcCU, uiAbsPartIdx, uiWidth, uiHeight, eTType)
    }
    if eTType == TLibCommon.TEXT_LUMA {
        eTType = TLibCommon.TEXT_LUMA
    } else {
        if eTType == TLibCommon.TEXT_NONE {
            eTType = TLibCommon.TEXT_NONE
        } else {
            eTType = TLibCommon.TEXT_CHROMA
        }
    }
    //----- encode significance map -----
    uiLog2BlockSize := uint(TLibCommon.G_aucConvertToBit[uiWidth]) + 2
    uiScanIdx := pcCU.GetCoefScanIdx(uiAbsPartIdx, uiWidth, eTType == TLibCommon.TEXT_LUMA, pcCU.IsIntra(uiAbsPartIdx))
    scan := TLibCommon.G_auiSigLastScan[uiScanIdx][uiLog2BlockSize-1]

    var beValid bool
    if pcCU.GetCUTransquantBypass1(uiAbsPartIdx) {
        beValid = false
    } else {
        beValid = pcCU.GetSlice().GetPPS().GetSignHideFlag()
    }

    // Find position of last coefficient
    scanPosLast := -1
    var posLast int

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
    var uiSigCoeffGroupFlag [TLibCommon.MLS_GRP_NUM]uint
    uiShift := uint(TLibCommon.MLS_CG_SIZE >> 1)
    uiNumBlkSide := uiWidth >> uiShift

    //::memset( uiSigCoeffGroupFlag, 0, sizeof(UInt) * MLS_GRP_NUM );

    //do
    {
        scanPosLast++
        posLast = int(scan[scanPosLast])

        // get L1 sig map
        uiPosY := posLast >> uiLog2BlockSize
        uiPosX := posLast - (uiPosY << uiLog2BlockSize)
        uiBlkIdx := uiNumBlkSide*uint(uiPosY>>uiShift) + uint(uiPosX>>uiShift)
        if pcCoef[posLast] != 0 {
            uiSigCoeffGroupFlag[uiBlkIdx] = 1
        }

        uiNumSig -= uint(TLibCommon.B2U(pcCoef[posLast] != 0))
    }

    for uiNumSig > 0 {
        scanPosLast++
        posLast = int(scan[scanPosLast])

        // get L1 sig map
        uiPosY := posLast >> uiLog2BlockSize
        uiPosX := posLast - (uiPosY << uiLog2BlockSize)
        uiBlkIdx := uiNumBlkSide*uint(uiPosY>>uiShift) + uint(uiPosX>>uiShift)
        if pcCoef[posLast] != 0 {
            uiSigCoeffGroupFlag[uiBlkIdx] = 1
        }

        uiNumSig -= uint(TLibCommon.B2U(pcCoef[posLast] != 0))
    }

    // Code position of last coefficient
    posLastY := posLast >> uiLog2BlockSize
    posLastX := posLast - (posLastY << uiLog2BlockSize)
    //fmt.Printf("posLastX%d, posLastY%d, uiWidth%d, uiHeight%d, eTType%d, uiScanIdx%d\n",posLastX, posLastY, uiWidth, uiHeight, eTType, uiScanIdx)
    //fmt.Printf("before uiBits=%d\n", this.getNumberOfWrittenBits());
    this.codeLastSignificantXY(uint(posLastX), uint(posLastY), int(uiWidth), int(uiHeight), eTType, uiScanIdx)
	//fmt.Printf("after uiBits=%d\n", this.getNumberOfWrittenBits());
	
    //===== code significance flag =====
    baseCoeffGroupCtx := this.m_cCUSigCoeffGroupSCModel.Get2(0, uint(eTType))
    var baseCtx []TLibCommon.ContextModel
    if eTType == TLibCommon.TEXT_LUMA {
        baseCtx = this.m_cCUSigSCModel.Get2(0, 0)
    } else {
        baseCtx = this.m_cCUSigSCModel.Get2(0, 0)[TLibCommon.NUM_SIG_FLAG_CTX_LUMA:]
    }

    iLastScanSet := int(scanPosLast >> TLibCommon.LOG2_SCAN_SET_SIZE)
    c1 := uint(1)
    uiGoRiceParam := uint(0)
    iScanPosSig := int(scanPosLast)

    for iSubSet := int(iLastScanSet); iSubSet >= 0; iSubSet-- {
        numNonZero := 0
        iSubPos := iSubSet << TLibCommon.LOG2_SCAN_SET_SIZE
        uiGoRiceParam = 0
        var absCoeff [16]int
        coeffSigns := uint(0)

        lastNZPosInCG := -1
        firstNZPosInCG := int(TLibCommon.SCAN_SET_SIZE)

        if iScanPosSig == int(scanPosLast) {
            absCoeff[0] = int(TLibCommon.ABS(pcCoef[posLast]).(TLibCommon.TCoeff))
            coeffSigns = uint(TLibCommon.B2U(pcCoef[posLast] < 0))
            numNonZero = 1
            lastNZPosInCG = iScanPosSig
            firstNZPosInCG = iScanPosSig
            iScanPosSig--
        }

        // encode significant_coeffgroup_flag
        iCGBlkPos := scanCG[iSubSet]
        iCGPosY := iCGBlkPos / uiNumBlkSide
        iCGPosX := iCGBlkPos - (iCGPosY * uiNumBlkSide)
        if iSubSet == int(iLastScanSet) || iSubSet == 0 {
            uiSigCoeffGroupFlag[iCGBlkPos] = 1
        } else {
            uiSigCoeffGroup := uint(TLibCommon.B2U(uiSigCoeffGroupFlag[iCGBlkPos] != 0))
            uiCtxSig := TLibCommon.GetSigCoeffGroupCtxInc(uiSigCoeffGroupFlag[:], iCGPosX, iCGPosY, int(uiWidth), int(uiHeight))
            this.m_pcBinIf.encodeBin(uiSigCoeffGroup, &baseCoeffGroupCtx[uiCtxSig])
            //this.DTRACE_CABAC_VL( g_nSymbolCounter++ );
            /*this.DTRACE_CABAC_T("\tuiSigCoeffGroup")
            this.DTRACE_CABAC_V(uiSigCoeffGroup)
            this.DTRACE_CABAC_T("\tuiCtxSig: ")
            this.DTRACE_CABAC_V(uiCtxSig)
            this.DTRACE_CABAC_T("\tuiBits: ")
            this.DTRACE_CABAC_V(this.getNumberOfWrittenBits())
            this.DTRACE_CABAC_T("\n")*/
        }

        // encode significant_coeff_flag
        if uiSigCoeffGroupFlag[iCGBlkPos] != 0 {
            patternSigCtx := TLibCommon.CalcPatternSigCtx(uiSigCoeffGroupFlag[:], iCGPosX, iCGPosY, int(uiWidth), int(uiHeight))
            var uiBlkPos, uiPosY, uiPosX, uiSig, uiCtxSig uint
            for ; iScanPosSig >= int(iSubPos); iScanPosSig-- {
                uiBlkPos = scan[iScanPosSig]
                uiPosY = uiBlkPos >> uiLog2BlockSize
                uiPosX = uiBlkPos - (uiPosY << uiLog2BlockSize)
                uiSig = uint(TLibCommon.B2U(pcCoef[uiBlkPos] != 0))
                if iScanPosSig > int(iSubPos) || iSubSet == 0 || numNonZero != 0 {
                    uiCtxSig = uint(TLibCommon.GetSigCtxInc(patternSigCtx, uiScanIdx, int(uiPosX), int(uiPosY), int(uiLog2BlockSize), eTType))
                    this.m_pcBinIf.encodeBin(uiSig, &baseCtx[uiCtxSig])
                    
                    //this.DTRACE_CABAC_VL( g_nSymbolCounter++ );
                    /*this.DTRACE_CABAC_T("\tuiSig")
                    this.DTRACE_CABAC_V(uiSig)
                    this.DTRACE_CABAC_T("\tuiCtxSig: ")
                    this.DTRACE_CABAC_V(uiCtxSig)
                    this.DTRACE_CABAC_T("\tuiBits: ")
            		this.DTRACE_CABAC_V(this.getNumberOfWrittenBits())
                    this.DTRACE_CABAC_T("\n")*/
                }
                if uiSig != 0 {
                    absCoeff[numNonZero] = int(TLibCommon.ABS(pcCoef[uiBlkPos]).(TLibCommon.TCoeff))
                    coeffSigns = 2*coeffSigns + uint(TLibCommon.B2U(pcCoef[uiBlkPos] < 0))
                    numNonZero++
                    if lastNZPosInCG == -1 {
                        lastNZPosInCG = iScanPosSig
                    }
                    firstNZPosInCG = iScanPosSig
                }
            }
        } else {
            iScanPosSig = iSubPos - 1
        }

        if numNonZero > 0 {
            signHidden := (lastNZPosInCG-firstNZPosInCG >= TLibCommon.SBH_THRESHOLD)
            var uiCtxSet uint
            if iSubSet > 0 && eTType == TLibCommon.TEXT_LUMA {
                uiCtxSet = 2
            } else {
                uiCtxSet = 0
            }

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

            numC1Flag := TLibCommon.MIN(int(numNonZero), int(TLibCommon.C1FLAG_NUMBER)).(int)
            firstC2FlagIdx := -1
            for idx := int(0); idx < numC1Flag; idx++ {
                uiSymbol := uint(TLibCommon.B2U(absCoeff[idx] > 1))
                this.m_pcBinIf.encodeBin(uiSymbol, &baseCtxMod[c1])
                
                //this.DTRACE_CABAC_VL( g_nSymbolCounter++ );
                /*this.DTRACE_CABAC_T("\tuiBin")
                this.DTRACE_CABAC_V(uiSymbol)
                this.DTRACE_CABAC_T("\tc1: ")
                this.DTRACE_CABAC_V(uint(c1))
                this.DTRACE_CABAC_T("\tuiBits: ")
            	this.DTRACE_CABAC_V(this.getNumberOfWrittenBits())
                this.DTRACE_CABAC_T("\n")*/
                
                if uiSymbol != 0 {
                    c1 = 0

                    if firstC2FlagIdx == -1 {
                        firstC2FlagIdx = idx
                    }
                } else if (c1 < 3) && (c1 > 0) {
                    c1++
                }
            }

            if c1 == 0 {
                if eTType == TLibCommon.TEXT_LUMA {
                    baseCtxMod = this.m_cCUAbsSCModel.Get2(0, 0)[uiCtxSet:]
                } else {
                    baseCtxMod = this.m_cCUAbsSCModel.Get2(0, 0)[TLibCommon.NUM_ABS_FLAG_CTX_LUMA+uiCtxSet:]
                }
                if firstC2FlagIdx != -1 {
                    symbol := uint(TLibCommon.B2U(absCoeff[firstC2FlagIdx] > 2))
                    this.m_pcBinIf.encodeBin(symbol, &baseCtxMod[0])
                    
                    //this.DTRACE_CABAC_VL( g_nSymbolCounter++ );
                    /*this.DTRACE_CABAC_T("\tuiBin")
                    this.DTRACE_CABAC_V(symbol)
                    this.DTRACE_CABAC_T("\tc1: ")
                    this.DTRACE_CABAC_V(0)
                    this.DTRACE_CABAC_T("\tuiBits: ")
            		this.DTRACE_CABAC_V(this.getNumberOfWrittenBits())
                    this.DTRACE_CABAC_T("\n")*/
                }
            }

            if beValid && signHidden {
                this.m_pcBinIf.encodeBinsEP((coeffSigns >> 1), numNonZero-1)
                
                //this.DTRACE_CABAC_VL( g_nSymbolCounter++ );
                /*this.DTRACE_CABAC_T("\tcoeffSigns")
                this.DTRACE_CABAC_V((coeffSigns >> 1))
                this.DTRACE_CABAC_T("\tnumNonZero-1: ")
                this.DTRACE_CABAC_V(uint(numNonZero - 1))
                this.DTRACE_CABAC_T("\tuiBits: ")
            	this.DTRACE_CABAC_V(this.getNumberOfWrittenBits())
                this.DTRACE_CABAC_T("\n")*/
            } else {
                this.m_pcBinIf.encodeBinsEP(coeffSigns, numNonZero)
                
                //this.DTRACE_CABAC_VL( g_nSymbolCounter++ );
                /*this.DTRACE_CABAC_T("\tcoeffSigns")
                this.DTRACE_CABAC_V(coeffSigns)
                this.DTRACE_CABAC_T("\tnumNonZero: ")
                this.DTRACE_CABAC_V(uint(numNonZero))
                this.DTRACE_CABAC_T("\tuiBits: ")
            	this.DTRACE_CABAC_V(this.getNumberOfWrittenBits())
                this.DTRACE_CABAC_T("\n")*/
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

                    if absCoeff[idx] >= int(baseLevel) {
                        this.xWriteCoefRemainExGolomb(uint(absCoeff[idx])-baseLevel, uiGoRiceParam)
                        //this.DTRACE_CABAC_VL( g_nSymbolCounter++ );
                        /*this.DTRACE_CABAC_T("\tuiLevel")
                        this.DTRACE_CABAC_V(uint(absCoeff[idx])-baseLevel)
                        this.DTRACE_CABAC_T("\tuiGoRiceParam: ")
                        this.DTRACE_CABAC_V(uint(uiGoRiceParam))
                        this.DTRACE_CABAC_T("\tuiBits: ")
            			this.DTRACE_CABAC_V(this.getNumberOfWrittenBits())
                        this.DTRACE_CABAC_T("\n")*/
                        
                        if absCoeff[idx] > 3*(1<<uiGoRiceParam) {
                            uiGoRiceParam = uint(TLibCommon.MIN(int(uiGoRiceParam+1), int(4)).(int))
                        }
                    }
                    if absCoeff[idx] >= 2 {
                        iFirstCoeff2 = 0
                    }
                }
            }
        }
    }

    return
}

func (this *TEncSbac) codeTransformSkipFlags(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint, width, height uint, eTType TLibCommon.TextType) {
    if pcCU.GetCUTransquantBypass1(uiAbsPartIdx) {
        return
    }
    if width != 4 || height != 4 {
        return
    }

    useTransformSkip := uint(TLibCommon.B2U(pcCU.GetTransformSkip2(uiAbsPartIdx, eTType)))
    if eTType != 0 {
        this.m_pcBinIf.encodeBin(useTransformSkip, this.m_cTransformSkipSCModel.Get3(0, TLibCommon.TEXT_CHROMA, 0))
    } else {
        this.m_pcBinIf.encodeBin(useTransformSkip, this.m_cTransformSkipSCModel.Get3(0, TLibCommon.TEXT_LUMA, 0))
    }
    /*DTRACE_CABAC_VL( g_nSymbolCounter++ )*/
    this.DTRACE_CABAC_T("\tparseTransformSkip()");
    this.DTRACE_CABAC_T( "\tsymbol=" )
    this.DTRACE_CABAC_V( useTransformSkip )
    this.DTRACE_CABAC_T( "\tAddr=" )
    this.DTRACE_CABAC_V( pcCU.GetAddr() )
    this.DTRACE_CABAC_T( "\tetype=" )
    this.DTRACE_CABAC_V( uint(eTType) )
    this.DTRACE_CABAC_T( "\tuiAbsPartIdx=" )
    this.DTRACE_CABAC_V( uiAbsPartIdx )
    this.DTRACE_CABAC_T( "\n" )
}

// -------------------------------------------------------------------------------------------------------------------
// for RD-optimizatioon
// -------------------------------------------------------------------------------------------------------------------

func (this *TEncSbac) estBit(pcEstBitsSbac *TLibCommon.EstBitsSbacStruct, width, height int, eTType TLibCommon.TextType) {
    this.estCBFBit(pcEstBitsSbac)

    this.estSignificantCoeffGroupMapBit(pcEstBitsSbac, eTType)

    // encode significance map
    this.estSignificantMapBit(pcEstBitsSbac, width, height, eTType)

    // encode significant coefficients
    this.estSignificantCoefficientsBit(pcEstBitsSbac, eTType)
}

func (this *TEncSbac) estCBFBit(pcEstBitsSbac *TLibCommon.EstBitsSbacStruct) {
    pCtx := this.m_cCUQtCbfSCModel.Get1(0)

    for uiCtxInc := uint(0); uiCtxInc < 3*TLibCommon.NUM_QT_CBF_CTX; uiCtxInc++ {
        pcEstBitsSbac.BlockCbpBits[uiCtxInc][0] = pCtx[uiCtxInc].GetEntropyBits(0)
        pcEstBitsSbac.BlockCbpBits[uiCtxInc][1] = pCtx[uiCtxInc].GetEntropyBits(1)
    }

    pCtx = this.m_cCUQtRootCbfSCModel.Get1(0)

    for uiCtxInc := 0; uiCtxInc < 4; uiCtxInc++ {
        pcEstBitsSbac.BlockRootCbpBits[uiCtxInc][0] = pCtx[uiCtxInc].GetEntropyBits(0)
        pcEstBitsSbac.BlockRootCbpBits[uiCtxInc][1] = pCtx[uiCtxInc].GetEntropyBits(1)
    }
}

func (this *TEncSbac) estSignificantCoeffGroupMapBit(pcEstBitsSbac *TLibCommon.EstBitsSbacStruct, eTType TLibCommon.TextType) {
    firstCtx := uint(0)
    numCtx := uint(TLibCommon.NUM_SIG_CG_FLAG_CTX)

    for ctxIdx := firstCtx; ctxIdx < firstCtx+numCtx; ctxIdx++ {
        for uiBin := int16(0); uiBin < 2; uiBin++ {
            pcEstBitsSbac.SignificantCoeffGroupBits[ctxIdx][uiBin] = this.m_cCUSigCoeffGroupSCModel.Get3(0, uint(eTType), ctxIdx).GetEntropyBits(int16(uiBin))
        }
    }
}

func (this *TEncSbac) estSignificantMapBit(pcEstBitsSbac *TLibCommon.EstBitsSbacStruct, width, height int, eTType TLibCommon.TextType) {
    firstCtx := 1
    numCtx := 8
    if TLibCommon.MAX(width, height).(int) >= 16 {
        if eTType == TLibCommon.TEXT_LUMA {
            firstCtx = 21
            numCtx = 6
        } else {
            firstCtx = 12
            numCtx = 3
        }
    } else if width == 8 {
        firstCtx = 9
        if eTType == TLibCommon.TEXT_LUMA {
            numCtx = 12
        } else {
            numCtx = 3
        }
    }

    if eTType == TLibCommon.TEXT_LUMA {
        for bin := int16(0); bin < 2; bin++ {
            pcEstBitsSbac.SignificantBits[0][bin] = this.m_cCUSigSCModel.Get3(0, 0, 0).GetEntropyBits(bin)
        }

        for ctxIdx := firstCtx; ctxIdx < firstCtx+numCtx; ctxIdx++ {
            for uiBin := int16(0); uiBin < 2; uiBin++ {
                pcEstBitsSbac.SignificantBits[ctxIdx][uiBin] = this.m_cCUSigSCModel.Get3(0, 0, uint(ctxIdx)).GetEntropyBits(uiBin)
            }
        }
    } else {
        for bin := int16(0); bin < 2; bin++ {
            pcEstBitsSbac.SignificantBits[0][bin] = this.m_cCUSigSCModel.Get3(0, 0, TLibCommon.NUM_SIG_FLAG_CTX_LUMA+0).GetEntropyBits(bin)
        }
        for ctxIdx := firstCtx; ctxIdx < firstCtx+numCtx; ctxIdx++ {
            for uiBin := int16(0); uiBin < 2; uiBin++ {
                pcEstBitsSbac.SignificantBits[ctxIdx][uiBin] = this.m_cCUSigSCModel.Get3(0, 0, uint(TLibCommon.NUM_SIG_FLAG_CTX_LUMA+ctxIdx)).GetEntropyBits(uiBin)
            }
        }
    }
    iBitsX := 0
    iBitsY := 0
    var blkSizeOffsetX, blkSizeOffsetY int
    var shiftX, shiftY uint

    if eTType != 0 {
        blkSizeOffsetX = 0
        blkSizeOffsetY = 0
        shiftX = uint(TLibCommon.G_aucConvertToBit[width])
        shiftY = uint(TLibCommon.G_aucConvertToBit[height])
    } else {
        blkSizeOffsetX = int(TLibCommon.G_aucConvertToBit[width]*3 + ((TLibCommon.G_aucConvertToBit[width] + 1) >> 2))
        blkSizeOffsetY = int(TLibCommon.G_aucConvertToBit[height]*3 + ((TLibCommon.G_aucConvertToBit[height] + 1) >> 2))
        shiftX = uint((TLibCommon.G_aucConvertToBit[width] + 3) >> 2)
        shiftY = uint((TLibCommon.G_aucConvertToBit[height] + 3) >> 2)
    }

    var ctx int
    pCtxX := this.m_cCuCtxLastX.Get2(0, uint(eTType))
    for ctx = 0; ctx < int(TLibCommon.G_uiGroupIdx[width-1]); ctx++ {
        ctxOffset := blkSizeOffsetX + (ctx >> shiftX)
        pcEstBitsSbac.LastXBits[ctx] = iBitsX + pCtxX[ctxOffset].GetEntropyBits(0)
        iBitsX += pCtxX[ctxOffset].GetEntropyBits(1)
    }
    pcEstBitsSbac.LastXBits[ctx] = iBitsX
    pCtxY := this.m_cCuCtxLastY.Get2(0, uint(eTType))
    for ctx = 0; ctx < int(TLibCommon.G_uiGroupIdx[height-1]); ctx++ {
        ctxOffset := blkSizeOffsetY + (ctx >> shiftY)
        pcEstBitsSbac.LastYBits[ctx] = iBitsY + pCtxY[ctxOffset].GetEntropyBits(0)
        iBitsY += pCtxY[ctxOffset].GetEntropyBits(1)
    }
    pcEstBitsSbac.LastYBits[ctx] = iBitsY
}

func (this *TEncSbac) estSignificantCoefficientsBit(pcEstBitsSbac *TLibCommon.EstBitsSbacStruct, eTType TLibCommon.TextType) {
    if eTType == TLibCommon.TEXT_LUMA {
        ctxOne := this.m_cCUOneSCModel.Get2(0, 0)
        ctxAbs := this.m_cCUAbsSCModel.Get2(0, 0)

        for ctxIdx := 0; ctxIdx < TLibCommon.NUM_ONE_FLAG_CTX_LUMA; ctxIdx++ {
            pcEstBitsSbac.GreaterOneBits[ctxIdx][0] = ctxOne[ctxIdx].GetEntropyBits(0)
            pcEstBitsSbac.GreaterOneBits[ctxIdx][1] = ctxOne[ctxIdx].GetEntropyBits(1)
        }

        for ctxIdx := 0; ctxIdx < TLibCommon.NUM_ABS_FLAG_CTX_LUMA; ctxIdx++ {
            pcEstBitsSbac.LevelAbsBits[ctxIdx][0] = ctxAbs[ctxIdx].GetEntropyBits(0)
            pcEstBitsSbac.LevelAbsBits[ctxIdx][1] = ctxAbs[ctxIdx].GetEntropyBits(1)
        }
    } else {
        ctxOne := this.m_cCUOneSCModel.Get2(0, 0)[TLibCommon.NUM_ONE_FLAG_CTX_LUMA:]
        ctxAbs := this.m_cCUAbsSCModel.Get2(0, 0)[TLibCommon.NUM_ABS_FLAG_CTX_LUMA:]

        for ctxIdx := 0; ctxIdx < TLibCommon.NUM_ONE_FLAG_CTX_CHROMA; ctxIdx++ {
            pcEstBitsSbac.GreaterOneBits[ctxIdx][0] = ctxOne[ctxIdx].GetEntropyBits(0)
            pcEstBitsSbac.GreaterOneBits[ctxIdx][1] = ctxOne[ctxIdx].GetEntropyBits(1)
        }

        for ctxIdx := 0; ctxIdx < TLibCommon.NUM_ABS_FLAG_CTX_CHROMA; ctxIdx++ {
            pcEstBitsSbac.LevelAbsBits[ctxIdx][0] = ctxAbs[ctxIdx].GetEntropyBits(0)
            pcEstBitsSbac.LevelAbsBits[ctxIdx][1] = ctxAbs[ctxIdx].GetEntropyBits(1)
        }
    }
}

func (this *TEncSbac) updateContextTables3(eSliceType TLibCommon.SliceType, iQp int, bExecuteFinish bool) {
    this.m_pcBinIf.encodeBinTrm(1)
    if bExecuteFinish {
        this.m_pcBinIf.finish()
    }
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
    this.m_cCUTransSubdivFlagSCModel.InitBuffer(eSliceType, iQp, TLibCommon.INIT_TRANS_SUBDIV_FLAG[:])
    this.m_cSaoMergeSCModel.InitBuffer(eSliceType, iQp, TLibCommon.INIT_SAO_MERGE_FLAG[:])
    this.m_cSaoTypeIdxSCModel.InitBuffer(eSliceType, iQp, TLibCommon.INIT_SAO_TYPE_IDX[:])
    this.m_cTransformSkipSCModel.InitBuffer(eSliceType, iQp, TLibCommon.INIT_TRANSFORMSKIP_FLAG[:])
    this.m_CUTransquantBypassFlagSCModel.InitBuffer(eSliceType, iQp, TLibCommon.INIT_CU_TRANSQUANT_BYPASS_FLAG[:])
    this.m_pcBinIf.start()
}

func (this *TEncSbac) updateContextTables2(eSliceType TLibCommon.SliceType, iQp int) {
    this.updateContextTables3(eSliceType, iQp, true)
}

func (this *TEncSbac) getEncBinIf() TEncBinIf {
    return this.m_pcBinIf
}
