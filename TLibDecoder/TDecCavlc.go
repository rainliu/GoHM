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

/*
#if ENC_DEC_TRACE

#define READ_CODE(length, code, name)     xReadCodeTr ( length, code, name )
#define READ_UVLC(        code, name)     xReadUvlcTr (         code, name )
#define READ_SVLC(        code, name)     xReadSvlcTr (         code, name )
#define READ_FLAG(        code, name)     xReadFlagTr (         code, name )

#else

#define READ_CODE(length, code, name)     xReadCode ( length, code )
#define READ_UVLC(        code, name)     xReadUvlc (         code )
#define READ_SVLC(        code, name)     xReadSvlc (         code )
#define READ_FLAG(        code, name)     xReadFlag (         code )

#endif
*/
//! \ingroup TLibDecoder
//! \{

// ====================================================================================================================
// Class definition
// ====================================================================================================================

type SyntaxElementParser struct {
    //protected:
    m_pTraceFile  io.Writer
    m_bSliceTrace bool
    m_pcBitstream *TLibCommon.TComInputBitstream
}

func NewSyntaxElementParser() *SyntaxElementParser {
    return &SyntaxElementParser{}
}

func (this *SyntaxElementParser) SetTraceFile(traceFile io.Writer) {
    this.m_pTraceFile = traceFile
}

func (this *SyntaxElementParser) GetTraceFile() io.Writer {
    return this.m_pTraceFile
}

func (this *SyntaxElementParser) SetSliceTrace(bSliceTrace bool) {
    this.m_bSliceTrace = bSliceTrace
}

func (this *SyntaxElementParser) GetSliceTrace() bool {
    return this.m_bSliceTrace
}

func (this *SyntaxElementParser) xTraceVUIHeader(pVUI *TLibCommon.TComVUI) {
    if this.GetTraceFile() != nil {
        io.WriteString(this.m_pTraceFile, fmt.Sprintf("========= VUI Parameter Set ===============================================\n")) //, pVPS.GetVPSId() );
    }
}

func (this *SyntaxElementParser) xTraceVPSHeader(pVPS *TLibCommon.TComVPS) {
    if this.GetTraceFile() != nil {
        io.WriteString(this.m_pTraceFile, fmt.Sprintf("========= Video Parameter Set =============================================\n")) //, pVPS.GetVPSId() );
    }
}

func (this *SyntaxElementParser) xTraceSPSHeader(pSPS *TLibCommon.TComSPS) {
    if this.GetTraceFile() != nil {
        io.WriteString(this.m_pTraceFile, fmt.Sprintf("========= Sequence Parameter Set ==========================================\n")) //, pSPS.GetSPSId() );
    }
}

func (this *SyntaxElementParser) xTracePPSHeader(pPPS *TLibCommon.TComPPS) {
    if this.GetTraceFile() != nil {
        io.WriteString(this.m_pTraceFile, fmt.Sprintf("========= Picture Parameter Set ===========================================\n")) //, pPPS.GetPPSId() );
    }
}

func (this *SyntaxElementParser) xTraceSliceHeader(pSlice *TLibCommon.TComSlice) {
    if this.GetTraceFile() != nil {
        if this.GetSliceTrace() {
            io.WriteString(this.m_pTraceFile, fmt.Sprintf("========= Slice Parameter Set =============================================\n"))
        }
    }
}

func (this *SyntaxElementParser) XTraceLCUHeader(traceLevel uint) {
}
func (this *SyntaxElementParser) XTraceCUHeader(traceLevel uint) {
}
func (this *SyntaxElementParser) XTracePUHeader(traceLevel uint) {
}
func (this *SyntaxElementParser) XTraceTUHeader(traceLevel uint) {
}
func (this *SyntaxElementParser) XTraceCoefHeader(traceLevel uint) {
}
func (this *SyntaxElementParser) XTraceResiHeader(traceLevel uint) {
}
func (this *SyntaxElementParser) XTracePredHeader(traceLevel uint) {
}
func (this *SyntaxElementParser) XTraceRecoHeader(traceLevel uint) {
}
func (this *SyntaxElementParser) XReadAeTr(Value int, pSymbolName string, traceLevel uint) {
}
func (this *SyntaxElementParser) XReadCeofTr(pCoeff []TLibCommon.TCoeff, uiWidth, traceLevel uint) {
}
func (this *SyntaxElementParser) XReadResiTr(pPel []TLibCommon.Pel, uiWidth, traceLevel uint) {
}
func (this *SyntaxElementParser) XReadPredTr(pPel []TLibCommon.Pel, uiWidth, traceLevel uint) {
}
func (this *SyntaxElementParser) XReadRecoTr(pPel []TLibCommon.Pel, uiWidth, traceLevel uint) {
}
/*
func (this *SyntaxElementParser) DTRACE_CABAC_F(x float32) {
    if this.GetTraceFile() != nil {
        io.WriteString(this.m_pTraceFile, fmt.Sprintf("%f", x))
    }
}
func (this *SyntaxElementParser) DTRACE_CABAC_V(x uint) {
    if this.GetTraceFile() != nil {
        io.WriteString(this.m_pTraceFile, fmt.Sprintf("%d", x))
    }
}
func (this *SyntaxElementParser) DTRACE_CABAC_VL(x uint) {
    if this.GetTraceFile() != nil {
        io.WriteString(this.m_pTraceFile, fmt.Sprintf("%lld", x))
    }
}
func (this *SyntaxElementParser) DTRACE_CABAC_T(x string) {
    if this.GetTraceFile() != nil {
        io.WriteString(this.m_pTraceFile, fmt.Sprintf("%s", x))
    }
}
func (this *SyntaxElementParser) DTRACE_CABAC_X(x uint) {
    if this.GetTraceFile() != nil {
        io.WriteString(this.m_pTraceFile, fmt.Sprintf("%x", x))
    }
}
func (this *SyntaxElementParser) DTRACE_CABAC_N() {
    if this.GetTraceFile() != nil {
        io.WriteString(this.m_pTraceFile, "\n")
    }
}
*/

func (this *SyntaxElementParser) xReadCode(length uint, val *uint) {
    //assert ( uiLength > 0 );
    this.m_pcBitstream.Read(length, val)
}
func (this *SyntaxElementParser) xReadUvlc(val *uint) {
    uiVal := uint(0)
    uiCode := uint(0)
    uiLength := uint(0)

    this.m_pcBitstream.Read(1, &uiCode)

    if 0 == uiCode {
        uiLength = 0

        for (uiCode & 1) == 0 {
            this.m_pcBitstream.Read(1, &uiCode)
            uiLength++
        }

        this.m_pcBitstream.Read(uiLength, &uiVal)

        uiVal += (1 << uiLength) - 1
    }

    *val = uiVal
}
func (this *SyntaxElementParser) xReadSvlc(val *int) {
    uiBits := uint(0)
    this.m_pcBitstream.Read(1, &uiBits)
    if 0 == uiBits {
        uiLength := uint(0)

        for (uiBits & 1) == 0 {
            this.m_pcBitstream.Read(1, &uiBits)
            uiLength++
        }

        this.m_pcBitstream.Read(uiLength, &uiBits)

        uiBits += (1 << uiLength)

        if (uiBits & 1) == 1 {
            *val = -int(uiBits >> 1)
        } else {
            *val = int(uiBits >> 1)
        }
    } else {
        *val = 0
    }
}
func (this *SyntaxElementParser) xReadFlag(val *uint) {
    this.m_pcBitstream.Read(1, val)
}

//#if ENC_DEC_TRACE
func (this *SyntaxElementParser) xReadCodeTr(length uint, rValue *uint, pSymbolName string) {
    this.xReadCode(length, rValue)
    if this.GetSliceTrace() {
        //fprintf( g_hTrace, "%8lld  ", g_nSymbolCounter++ );
        if length < 10 {
            io.WriteString(this.m_pTraceFile, fmt.Sprintf("%-62s u(%d)  : %4d\n", pSymbolName, length, *rValue))
        } else {
            io.WriteString(this.m_pTraceFile, fmt.Sprintf("%-62s u(%d) : %4d\n", pSymbolName, length, *rValue))
        }
        //fflush ( g_hTrace );
    }
}
func (this *SyntaxElementParser) xReadUvlcTr(rValue *uint, pSymbolName string) {
    this.xReadUvlc(rValue)
    if this.GetSliceTrace() {
        //fprintf( g_hTrace, "%8lld  ", g_nSymbolCounter++ );
        io.WriteString(this.m_pTraceFile, fmt.Sprintf("%-62s ue(v) : %4d\n", pSymbolName, *rValue))
        //fflush ( g_hTrace );
    }
}
func (this *SyntaxElementParser) xReadSvlcTr(rValue *int, pSymbolName string) {
    this.xReadSvlc(rValue)
    if this.GetSliceTrace() {
        //fprintf( g_hTrace, "%8lld  ", g_nSymbolCounter++ );
        io.WriteString(this.m_pTraceFile, fmt.Sprintf("%-62s se(v) : %4d\n", pSymbolName, *rValue))
        //fflush ( g_hTrace );
    }
}
func (this *SyntaxElementParser) xReadFlagTr(rValue *uint, pSymbolName string) {
    this.xReadFlag(rValue)
    if this.GetSliceTrace() {
        //fprintf( g_hTrace, "%8lld  ", g_nSymbolCounter++ );
        io.WriteString(this.m_pTraceFile, fmt.Sprintf("%-62s u(1)  : %4d\n", pSymbolName, *rValue))
        //fflush ( g_hTrace );
    }
}

//#endif
//public:
func (this *SyntaxElementParser) SetBitstream(p *TLibCommon.TComInputBitstream) {
    this.m_pcBitstream = p
}
func (this *SyntaxElementParser) GetBitstream() *TLibCommon.TComInputBitstream {
    return this.m_pcBitstream
}

//};

//class SEImessages;

/// CAVLC decoder class
type TDecCavlc struct {
    SyntaxElementParser //, public TDecEntropyIf
}

func NewTDecCavlc() *TDecCavlc {
    return &TDecCavlc{}
}

func (this *TDecCavlc) READ_CODE(length uint, rValue *uint, pSymbolName string) {
    this.xReadCodeTr(length, rValue, pSymbolName)
}

func (this *TDecCavlc) READ_UVLC(rValue *uint, pSymbolName string) {
    this.xReadUvlcTr(rValue, pSymbolName)
}
func (this *TDecCavlc) READ_SVLC(rValue *int, pSymbolName string) {
    this.xReadSvlcTr(rValue, pSymbolName)
}
func (this *TDecCavlc) READ_FLAG(rValue *uint, pSymbolName string) {
    this.xReadFlagTr(rValue, pSymbolName)
}

//protected:
func (this *TDecCavlc) xReadEpExGolomb(ruiSymbol *uint, uiCount uint) {
}
func (this *TDecCavlc) xReadExGolombLevel(ruiSymbol *uint) {
}
func (this *TDecCavlc) xReadUnaryMaxSymbol(ruiSymbol, uiMaxSymbol uint) {
}

func (this *TDecCavlc) xReadPCMAlignZero() {
    uiNumberOfBits := this.m_pcBitstream.GetNumBitsUntilByteAligned()

    if uiNumberOfBits != 0 {
        var uiBits, uiSymbol uint

        for uiBits = 0; uiBits < uiNumberOfBits; uiBits++ {
            this.xReadFlag(&uiSymbol)
            //assert( uiSymbol == 0 );
        }
    }
}

func (this *TDecCavlc) xGetBit() uint {
    return 0
}

func (this *TDecCavlc) ParseShortTermRefPicSet(sps *TLibCommon.TComSPS, rps *TLibCommon.TComReferencePictureSet, idx int) {
    var code, interRPSPred uint

    if idx > 0 {
        this.READ_FLAG(&interRPSPred, "inter_ref_pic_set_prediction_flag")
        rps.SetInterRPSPrediction(interRPSPred != 0)
    } else {
        interRPSPred = 0
        rps.SetInterRPSPrediction(false)
    }

    if interRPSPred != 0 {
        var bit uint
        if idx == sps.GetRPSList().GetNumberOfReferencePictureSets() {
            this.READ_UVLC(&code, "delta_idx_minus1") // delta index of the Reference Picture Set used for prediction minus 1
        } else {
            code = 0
        }
        //assert(code <= idx-1); // delta_idx_minus1 shall not be larger than idx-1, otherwise we will predict from a negative row position that does not exist. When idx equals 0 there is no legal value and interRPSPred must be zero. See J0185-r2
        rIdx := idx - 1 - int(code)
        //assert (rIdx <= idx-1 && rIdx >= 0); // Made assert tighter; if rIdx = idx then prediction is done from itself. rIdx must belong to range 0, idx-1, inclusive, see J0185-r2
        rpsRef := sps.GetRPSList().GetReferencePictureSet(rIdx)
        k := 0
        k0 := 0
        k1 := 0
        this.READ_CODE(1, &bit, "delta_rps_sign")     // delta_RPS_sign
        this.READ_UVLC(&code, "abs_delta_rps_minus1") // absolute delta RPS minus 1
        deltaRPS := (1 - 2*bit) * (code + 1)          // delta_RPS
        for j := 0; j <= (rpsRef.GetNumberOfNegativePictures()+rpsRef.GetNumberOfPositivePictures()); j++ {
            this.READ_CODE(1, &bit, "used_by_curr_pic_flag") //first bit is "1" if Idc is 1
            refIdc := bit
            if refIdc == 0 {
                this.READ_CODE(1, &bit, "use_delta_flag") //second bit is "1" if Idc is 2, "0" otherwise.
                refIdc = bit << 1                         //second bit is "1" if refIdc is 2, "0" if refIdc = 0.
            }

            if refIdc == 1 || refIdc == 2 {
                var deltaPOC int
                if j < (rpsRef.GetNumberOfNegativePictures()+rpsRef.GetNumberOfPositivePictures()) {
                    deltaPOC = int(deltaRPS) + rpsRef.GetDeltaPOC(j)
                } else {
                    deltaPOC = int(deltaRPS)
                }
                rps.SetDeltaPOC(k, deltaPOC)
                rps.SetUsed(k, (refIdc == 1))

                if deltaPOC < 0 {
                    k0++
                } else {
                    k1++
                }
                k++
            }
            rps.SetRefIdc(j, int(refIdc))
        }
        rps.SetNumRefIdc(rpsRef.GetNumberOfPictures() + 1)
        rps.SetNumberOfPictures(k)
        rps.SetNumberOfNegativePictures(k0)
        rps.SetNumberOfPositivePictures(k1)
        rps.SortDeltaPOC()
    } else {
        this.READ_UVLC(&code, "num_negative_pics")
        rps.SetNumberOfNegativePictures(int(code))
        this.READ_UVLC(&code, "num_positive_pics")
        rps.SetNumberOfPositivePictures(int(code))
        prev := 0
        var poc int
        for j := 0; j < rps.GetNumberOfNegativePictures(); j++ {
            this.READ_UVLC(&code, "delta_poc_s0_minus1")
            poc = prev - int(code) - 1
            prev = poc
            rps.SetDeltaPOC(j, poc)
            this.READ_FLAG(&code, "used_by_curr_pic_s0_flag")
            rps.SetUsed(j, code != 0)
        }
        prev = 0
        for j := rps.GetNumberOfNegativePictures(); j < rps.GetNumberOfNegativePictures()+rps.GetNumberOfPositivePictures(); j++ {
            this.READ_UVLC(&code, "delta_poc_s1_minus1")
            poc = prev + int(code) + 1
            prev = poc
            rps.SetDeltaPOC(j, poc)
            this.READ_FLAG(&code, "used_by_curr_pic_s1_flag")
            rps.SetUsed(j, code != 0)
        }
        rps.SetNumberOfPictures(rps.GetNumberOfNegativePictures() + rps.GetNumberOfPositivePictures())
    }
    //#if PRINT_RPS_INFO
    //  rps->printDeltaPOC();
    //#endif
}

//public:

/// rest entropy coder by intial QP and IDC in CABAC
func (this *TDecCavlc) ResetEntropy(pcSlice *TLibCommon.TComSlice) {
    //assert(0);
}
func (this *TDecCavlc) SetBitstream(p *TLibCommon.TComInputBitstream) {
    this.m_pcBitstream = p
}
func (this *TDecCavlc) ParseTransformSubdivFlag(ruiSubdivFlag *uint, uiLog2TransformBlockSize uint) {
}
func (this *TDecCavlc) ParseQtCbf(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint, eType TLibCommon.TextType, uiTrDepth, uiDepth uint) {
}
func (this *TDecCavlc) ParseQtRootCbf(uiAbsPartIdx uint, uiQtRootCbf *uint) {
}
func (this *TDecCavlc) ParseVPS(pcVPS *TLibCommon.TComVPS) {

    var uiCode uint
    //#if ENC_DEC_TRACE
    this.xTraceVPSHeader(pcVPS)
    //#endif

    this.READ_CODE(4, &uiCode, "vps_video_parameter_set_id")
    pcVPS.SetVPSId(int(uiCode))

    this.READ_CODE(2, &uiCode, "vps_reserved_three_2bits") //assert(uiCode == 3);
    this.READ_CODE(6, &uiCode, "vps_reserved_zero_6bits")  //assert(uiCode == 0);
    this.READ_CODE(3, &uiCode, "vps_max_sub_layers_minus1")
    pcVPS.SetMaxTLayers(uiCode + 1)

    this.READ_FLAG(&uiCode, "vps_temporal_id_nesting_flag")
    pcVPS.SetTemporalNestingFlag(uiCode != 0)

    this.READ_CODE(16, &uiCode, "vps_reserved_ffff_16bits") //assert(uiCode == 0xffff);
    this.ParsePTL(pcVPS.GetPTL(), true, int(pcVPS.GetMaxTLayers())-1)

    var subLayerOrderingInfoPresentFlag uint
    this.READ_FLAG(&subLayerOrderingInfoPresentFlag, "vps_sub_layer_ordering_info_present_flag")

    for i := uint(0); i <= pcVPS.GetMaxTLayers()-1; i++ {
        this.READ_UVLC(&uiCode, "vps_max_dec_pic_buffering_minus1[i]")
        pcVPS.SetMaxDecPicBuffering(uiCode+1, i)
        this.READ_UVLC(&uiCode, "vps_num_reorder_pics[i]")
        pcVPS.SetNumReorderPics(uiCode, i)
        this.READ_UVLC(&uiCode, "vps_max_latency_increase[i]")
        pcVPS.SetMaxLatencyIncrease(uiCode, i)

        if subLayerOrderingInfoPresentFlag == 0 {
            for i++; i <= pcVPS.GetMaxTLayers()-1; i++ {
                pcVPS.SetMaxDecPicBuffering(pcVPS.GetMaxDecPicBuffering(0), i)
                pcVPS.SetNumReorderPics(pcVPS.GetNumReorderPics(0), i)
                pcVPS.SetMaxLatencyIncrease(pcVPS.GetMaxLatencyIncrease(0), i)
            }
            break
        }
    }

    this.READ_CODE(6, &uiCode, "vps_max_nuh_reserved_zero_layer_id")
    pcVPS.SetMaxNuhReservedZeroLayerId(uiCode)
    this.READ_UVLC(&uiCode, "vps_max_op_sets_minus1")
    pcVPS.SetMaxOpSets(uiCode + 1)
    for opsIdx := uint(1); opsIdx <= (pcVPS.GetMaxOpSets() - 1); opsIdx++ {
        // Operation point set
        for i := uint(0); i <= pcVPS.GetMaxNuhReservedZeroLayerId(); i++ {
            this.READ_FLAG(&uiCode, "layer_id_included_flag[opsIdx][i]")
            pcVPS.SetLayerIdIncludedFlag(uiCode == 1, opsIdx, i)
        }
    }

    timingInfo := pcVPS.GetTimingInfo()
    this.READ_FLAG(&uiCode, "vps_timing_info_present_flag")
    timingInfo.SetTimingInfoPresentFlag(uiCode != 0)
    if timingInfo.GetTimingInfoPresentFlag() {
        this.READ_CODE(32, &uiCode, "vps_num_units_in_tick")
        timingInfo.SetNumUnitsInTick(uiCode)
        this.READ_CODE(32, &uiCode, "vps_time_scale")
        timingInfo.SetTimeScale(uiCode)
        this.READ_FLAG(&uiCode, "vps_poc_proportional_to_timing_flag")
        timingInfo.SetPocProportionalToTimingFlag(uiCode != 0)
        if timingInfo.GetPocProportionalToTimingFlag() {
            this.READ_UVLC(&uiCode, "vps_num_ticks_poc_diff_one_minus1")
            timingInfo.SetNumTicksPocDiffOneMinus1(int(uiCode))
        }

        this.READ_UVLC(&uiCode, "vps_num_hrd_parameters")
        pcVPS.SetNumHrdParameters(uiCode)

        if pcVPS.GetNumHrdParameters() > 0 {
            pcVPS.CreateHrdParamBuffer()
        }
        for i := uint(0); i < pcVPS.GetNumHrdParameters(); i++ {
            this.READ_UVLC(&uiCode, "hrd_op_set_idx")
            pcVPS.SetHrdOpSetIdx(uiCode, i)
            if i > 0 {
                this.READ_FLAG(&uiCode, "cprms_present_flag[i]")
                pcVPS.SetCprmsPresentFlag(uiCode == 1, i)
            }
            this.ParseHrdParameters(pcVPS.GetHrdParameters(i), pcVPS.GetCprmsPresentFlag(i), pcVPS.GetMaxTLayers()-1)
        }
    }

    this.READ_FLAG(&uiCode, "vps_extension_flag") //assert(!uiCode);
    if uiCode != 0 {
        for this.xMoreRbspData() {
            this.READ_FLAG(&uiCode, "vps_extension_data_flag")
        }
    }
}
func (this *TDecCavlc) ParseSPS(pcSPS *TLibCommon.TComSPS) {
    //#if ENC_DEC_TRACE
    this.xTraceSPSHeader(pcSPS)
    //#endif

    var uiCode uint
    this.READ_CODE(4, &uiCode, "sps_video_parameter_set_id")
    pcSPS.SetVPSId(int(uiCode))
    this.READ_CODE(3, &uiCode, "sps_max_sub_layers_minus1")
    pcSPS.SetMaxTLayers(uiCode + 1)

    this.READ_FLAG(&uiCode, "sps_temporal_id_nesting_flag")
    if uiCode > 0 {
        pcSPS.SetTemporalIdNestingFlag(true)
    } else {
        pcSPS.SetTemporalIdNestingFlag(false)
    }

    this.ParsePTL(pcSPS.GetPTL(), true, int(pcSPS.GetMaxTLayers())-1)
    this.READ_UVLC(&uiCode, "sps_seq_parameter_set_id")
    pcSPS.SetSPSId(int(uiCode))
    this.READ_UVLC(&uiCode, "chroma_format_idc")
    pcSPS.SetChromaFormatIdc(int(uiCode))
    // in the first version we only support chroma_format_idc equal to 1 (4:2:0), so separate_colour_plane_flag cannot appear in the bitstream
    //assert (uiCode == 1);
    if uiCode == 3 {
        this.READ_FLAG(&uiCode, "separate_colour_plane_flag") //assert(uiCode == 0);
    }

    this.READ_UVLC(&uiCode, "pic_width_in_luma_samples")
    pcSPS.SetPicWidthInLumaSamples(uiCode)
    this.READ_UVLC(&uiCode, "pic_height_in_luma_samples")
    pcSPS.SetPicHeightInLumaSamples(uiCode)
    this.READ_FLAG(&uiCode, "conformance_window_flag")
    if uiCode != 0 {
        crop := pcSPS.GetConformanceWindow()
        this.READ_UVLC(&uiCode, "conf_win_left_offset")
        crop.SetWindowLeftOffset(int(uiCode) * pcSPS.GetWinUnitX(pcSPS.GetChromaFormatIdc()))
        this.READ_UVLC(&uiCode, "conf_win_right_offset")
        crop.SetWindowRightOffset(int(uiCode) * pcSPS.GetWinUnitX(pcSPS.GetChromaFormatIdc()))
        this.READ_UVLC(&uiCode, "conf_win_top_offset")
        crop.SetWindowTopOffset(int(uiCode) * pcSPS.GetWinUnitY(pcSPS.GetChromaFormatIdc()))
        this.READ_UVLC(&uiCode, "conf_win_bottom_offset")
        crop.SetWindowBottomOffset(int(uiCode) * pcSPS.GetWinUnitY(pcSPS.GetChromaFormatIdc()))
    }

    this.READ_UVLC(&uiCode, "bit_depth_luma_minus8")
    TLibCommon.G_bitDepthY = 8 + int(uiCode)
    pcSPS.SetBitDepthY(TLibCommon.G_bitDepthY)
    pcSPS.SetQpBDOffsetY(int(6 * uiCode))

    this.READ_UVLC(&uiCode, "bit_depth_chroma_minus8")
    TLibCommon.G_bitDepthC = 8 + int(uiCode)
    pcSPS.SetBitDepthC(TLibCommon.G_bitDepthC)
    pcSPS.SetQpBDOffsetC(int(6 * uiCode))

    this.READ_UVLC(&uiCode, "log2_max_pic_order_cnt_lsb_minus4")
    pcSPS.SetBitsForPOC(4 + uiCode)

    var subLayerOrderingInfoPresentFlag uint
    this.READ_FLAG(&subLayerOrderingInfoPresentFlag, "sps_sub_layer_ordering_info_present_flag")

    for i := uint(0); i <= pcSPS.GetMaxTLayers()-1; i++ {
        this.READ_UVLC(&uiCode, "sps_max_dec_pic_buffering_minus1")
        pcSPS.SetMaxDecPicBuffering(uiCode+1, i)
        this.READ_UVLC(&uiCode, "sps_num_reorder_pics")
        pcSPS.SetNumReorderPics(int(uiCode), i)
        this.READ_UVLC(&uiCode, "sps_max_latency_increase")
        pcSPS.SetMaxLatencyIncrease(uiCode, i)

        if subLayerOrderingInfoPresentFlag == 0 {
            for i++; i <= pcSPS.GetMaxTLayers()-1; i++ {
                pcSPS.SetMaxDecPicBuffering(pcSPS.GetMaxDecPicBuffering(0), i)
                pcSPS.SetNumReorderPics(pcSPS.GetNumReorderPics(0), i)
                pcSPS.SetMaxLatencyIncrease(pcSPS.GetMaxLatencyIncrease(0), i)
            }
            break
        }
    }

    this.READ_UVLC(&uiCode, "log2_min_coding_block_size_minus3")
    log2MinCUSize := uiCode + 3
    this.READ_UVLC(&uiCode, "log2_diff_max_min_coding_block_size")
    uiMaxCUDepthCorrect := uiCode
    pcSPS.SetMaxCUWidth(1 << (log2MinCUSize + uiMaxCUDepthCorrect))
    //TLibCommon.G_uiMaxCUWidth = 1 << (log2MinCUSize + uiMaxCUDepthCorrect)
    pcSPS.SetMaxCUHeight(1 << (log2MinCUSize + uiMaxCUDepthCorrect))
    //TLibCommon.G_uiMaxCUHeight = 1 << (log2MinCUSize + uiMaxCUDepthCorrect)
    this.READ_UVLC(&uiCode, "log2_min_transform_block_size_minus2")
    pcSPS.SetQuadtreeTULog2MinSize(uiCode + 2)

    this.READ_UVLC(&uiCode, "log2_diff_max_min_transform_block_size")
    pcSPS.SetQuadtreeTULog2MaxSize(uiCode + pcSPS.GetQuadtreeTULog2MinSize())
    pcSPS.SetMaxTrSize(1 << (uiCode + pcSPS.GetQuadtreeTULog2MinSize()))

    this.READ_UVLC(&uiCode, "max_transform_hierarchy_depth_inter")
    pcSPS.SetQuadtreeTUMaxDepthInter(uiCode + 1)
    this.READ_UVLC(&uiCode, "max_transform_hierarchy_depth_intra")
    pcSPS.SetQuadtreeTUMaxDepthIntra(uiCode + 1)
    uiAddCUDepth := uint(0)
    for (pcSPS.GetMaxCUWidth() >> uiMaxCUDepthCorrect) > (1 << (pcSPS.GetQuadtreeTULog2MinSize() + uiAddCUDepth)) {
        uiAddCUDepth++
    }
    pcSPS.SetMaxCUDepth(uiMaxCUDepthCorrect + uiAddCUDepth)
    pcSPS.SetAddCUDepth( uiAddCUDepth );
    //TLibCommon.G_uiMaxCUDepth = uiMaxCUDepthCorrect + TLibCommon.G_uiAddCUDepth
    // BB: these parameters may be removed completly and replaced by the fixed values
    pcSPS.SetMinTrDepth(0)
    pcSPS.SetMaxTrDepth(1)
    this.READ_FLAG(&uiCode, "scaling_list_enabled_flag")
    pcSPS.SetScalingListFlag(uiCode == 1)
    if pcSPS.GetScalingListFlag() {
        this.READ_FLAG(&uiCode, "sps_scaling_list_data_present_flag")
        pcSPS.SetScalingListPresentFlag(uiCode == 1)
        if pcSPS.GetScalingListPresentFlag() {
            this.ParseScalingList(pcSPS.GetScalingList())
        }
    }
    this.READ_FLAG(&uiCode, "amp_enabled_flag")
    pcSPS.SetUseAMP(uiCode == 1)
    this.READ_FLAG(&uiCode, "sample_adaptive_offset_enabled_flag")
    pcSPS.SetUseSAO(uiCode != 0)

    this.READ_FLAG(&uiCode, "pcm_enabled_flag")
    pcSPS.SetUsePCM(uiCode != 0)
    if pcSPS.GetUsePCM() {
        this.READ_CODE(4, &uiCode, "pcm_sample_bit_depth_luma_minus1")
        pcSPS.SetPCMBitDepthLuma(1 + uiCode)
        this.READ_CODE(4, &uiCode, "pcm_sample_bit_depth_chroma_minus1")
        pcSPS.SetPCMBitDepthChroma(1 + uiCode)
        this.READ_UVLC(&uiCode, "log2_min_pcm_luma_coding_block_size_minus3")
        pcSPS.SetPCMLog2MinSize(uiCode + 3)
        this.READ_UVLC(&uiCode, "log2_diff_max_min_pcm_luma_coding_block_size")
        pcSPS.SetPCMLog2MaxSize(uiCode + pcSPS.GetPCMLog2MinSize())
        this.READ_FLAG(&uiCode, "pcm_loop_filter_disable_flag")
        pcSPS.SetPCMFilterDisableFlag(uiCode != 0)
    }

    this.READ_UVLC(&uiCode, "num_short_term_ref_pic_sets")
    pcSPS.CreateRPSList(int(uiCode))

    rpsList := pcSPS.GetRPSList()

    for i := 0; i < rpsList.GetNumberOfReferencePictureSets(); i++ {
        rps := rpsList.GetReferencePictureSet(i)
        this.ParseShortTermRefPicSet(pcSPS, rps, i)
    }
    this.READ_FLAG(&uiCode, "long_term_ref_pics_present_flag")
    pcSPS.SetLongTermRefsPresent(uiCode != 0)
    if pcSPS.GetLongTermRefsPresent() {
        this.READ_UVLC(&uiCode, "num_long_term_ref_pic_sps")
        pcSPS.SetNumLongTermRefPicSPS(uiCode)
        for k := 0; k < int(pcSPS.GetNumLongTermRefPicSPS()); k++ {
            this.READ_CODE(pcSPS.GetBitsForPOC(), &uiCode, "lt_ref_pic_poc_lsb_sps")
            pcSPS.SetLtRefPicPocLsbSps(uint(k), uiCode)
            this.READ_FLAG(&uiCode, "used_by_curr_pic_lt_sps_flag[i]")
            if uiCode != 0 {
                pcSPS.SetUsedByCurrPicLtSPSFlag(k, true)
            } else {
                pcSPS.SetUsedByCurrPicLtSPSFlag(k, false)
            }
        }
    }
    this.READ_FLAG(&uiCode, "sps_temporal_mvp_enable_flag")
    pcSPS.SetTMVPFlagsPresent(uiCode != 0)

    this.READ_FLAG(&uiCode, "sps_strong_intra_smoothing_enable_flag")
    pcSPS.SetUseStrongIntraSmoothing(uiCode != 0)

    this.READ_FLAG(&uiCode, "vui_parameters_present_flag")
    pcSPS.SetVuiParametersPresentFlag(uiCode != 0)

    if pcSPS.GetVuiParametersPresentFlag() {
        this.ParseVUI(pcSPS.GetVuiParameters(), pcSPS)
    }

    this.READ_FLAG(&uiCode, "sps_extension_flag")
    if uiCode != 0 {
        for this.xMoreRbspData() {
            this.READ_FLAG(&uiCode, "sps_extension_data_flag")
        }
    }
}

func (this *TDecCavlc) ParsePPS(pcPPS *TLibCommon.TComPPS) {
    //#if ENC_DEC_TRACE
    this.xTracePPSHeader(pcPPS)
    //#endif
    var uiCode uint
    var iCode int

    this.READ_UVLC(&uiCode, "pps_pic_parameter_set_id")
    pcPPS.SetPPSId(int(uiCode))
    this.READ_UVLC(&uiCode, "pps_seq_parameter_set_id")
    pcPPS.SetSPSId(int(uiCode))
    this.READ_FLAG(&uiCode, "dependent_slice_segments_enabled_flag")
    pcPPS.SetDependentSliceSegmentsEnabledFlag(uiCode == 1)

    this.READ_FLAG(&uiCode, "output_flag_present_flag")
    pcPPS.SetOutputFlagPresentFlag(uiCode == 1)

    this.READ_CODE(3, &uiCode, "num_extra_slice_header_bits")
    pcPPS.SetNumExtraSliceHeaderBits(int(uiCode))

    this.READ_FLAG(&uiCode, "sign_data_hiding_flag")
    pcPPS.SetSignHideFlag(uiCode != 0)

    this.READ_FLAG(&uiCode, "cabac_init_present_flag")
    pcPPS.SetCabacInitPresentFlag(uiCode != 0)

    this.READ_UVLC(&uiCode, "num_ref_idx_l0_default_active_minus1")
    //assert(uiCode <= 14);
    pcPPS.SetNumRefIdxL0DefaultActive(uiCode + 1)

    this.READ_UVLC(&uiCode, "num_ref_idx_l1_default_active_minus1")
    //assert(uiCode <= 14);
    pcPPS.SetNumRefIdxL1DefaultActive(uiCode + 1)

    this.READ_SVLC(&iCode, "init_qp_minus26")
    pcPPS.SetPicInitQPMinus26(iCode)
    this.READ_FLAG(&uiCode, "constrained_intra_pred_flag")
    pcPPS.SetConstrainedIntraPred(uiCode != 0)
    this.READ_FLAG(&uiCode, "transform_skip_enabled_flag")
    pcPPS.SetUseTransformSkip(uiCode != 0)

    this.READ_FLAG(&uiCode, "cu_qp_delta_enabled_flag")
    pcPPS.SetUseDQP(uiCode != 0)
    if pcPPS.GetUseDQP() {
        this.READ_UVLC(&uiCode, "diff_cu_qp_delta_depth")
        pcPPS.SetMaxCuDQPDepth(uiCode)
    } else {
        pcPPS.SetMaxCuDQPDepth(0)
    }
    this.READ_SVLC(&iCode, "pps_cb_qp_offset")
    pcPPS.SetChromaCbQpOffset(iCode)
    //assert( pcPPS.GetChromaCbQpOffset() >= -12 );
    //assert( pcPPS.GetChromaCbQpOffset() <=  12 );

    this.READ_SVLC(&iCode, "pps_cr_qp_offset")
    pcPPS.SetChromaCrQpOffset(iCode)
    //assert( pcPPS.GetChromaCrQpOffset() >= -12 );
    //assert( pcPPS.GetChromaCrQpOffset() <=  12 );

    this.READ_FLAG(&uiCode, "pps_slice_chroma_qp_offsets_present_flag")
    pcPPS.SetSliceChromaQpFlag(uiCode != 0)

    this.READ_FLAG(&uiCode, "weighted_pred_flag") // Use of Weighting Prediction (P_SLICE)
    pcPPS.SetUseWP(uiCode == 1)
    this.READ_FLAG(&uiCode, "weighted_bipred_flag") // Use of Bi-Directional Weighting Prediction (B_SLICE)
    pcPPS.SetWPBiPred(uiCode == 1)

    //this.READ_FLAG(&uiCode, "output_flag_present_flag")
    //pcPPS.SetOutputFlagPresentFlag(uiCode != 0)

    this.READ_FLAG(&uiCode, "transquant_bypass_enable_flag")
    pcPPS.SetTransquantBypassEnableFlag(uiCode != 0)

    this.READ_FLAG(&uiCode, "tiles_enabled_flag")
    pcPPS.SetTilesEnabledFlag(uiCode == 1)
    this.READ_FLAG(&uiCode, "entropy_coding_sync_enabled_flag")
    pcPPS.SetEntropyCodingSyncEnabledFlag(uiCode == 1)

    if pcPPS.GetTilesEnabledFlag() {
        this.READ_UVLC(&uiCode, "num_tile_columns_minus1")
        pcPPS.SetNumColumnsMinus1(int(uiCode))
        this.READ_UVLC(&uiCode, "num_tile_rows_minus1")
        pcPPS.SetNumRowsMinus1(int(uiCode))
        this.READ_FLAG(&uiCode, "uniform_spacing_flag")
        pcPPS.SetUniformSpacingFlag(uiCode != 0)

        if !pcPPS.GetUniformSpacingFlag() {
            columnWidth := make([]int, pcPPS.GetNumColumnsMinus1()) //UInt* columnWidth = (UInt*)malloc(pcPPS.GetNumColumnsMinus1()*sizeof(UInt));
            for i := 0; i < pcPPS.GetNumColumnsMinus1(); i++ {
                this.READ_UVLC(&uiCode, "column_width_minus1")
                columnWidth[i] = int(uiCode) + 1
            }
            pcPPS.SetColumnWidth(columnWidth)
            //free(columnWidth);

            rowHeight := make([]int, pcPPS.GetNumRowsMinus1()) //UInt* rowHeight = (UInt*)malloc(pcPPS.GetNumRowsMinus1()*sizeof(UInt));
            for i := 0; i < pcPPS.GetNumRowsMinus1(); i++ {
                this.READ_UVLC(&uiCode, "row_height_minus1")
                rowHeight[i] = int(uiCode) + 1
            }
            pcPPS.SetRowHeight(rowHeight)
            //free(rowHeight);
        }

        if pcPPS.GetNumColumnsMinus1() != 0 || pcPPS.GetNumRowsMinus1() != 0 {
            this.READ_FLAG(&uiCode, "loop_filter_across_tiles_enabled_flag")
            pcPPS.SetLoopFilterAcrossTilesEnabledFlag(uiCode != 0)
        }
    }
    this.READ_FLAG(&uiCode, "loop_filter_across_slices_enabled_flag")
    pcPPS.SetLoopFilterAcrossSlicesEnabledFlag(uiCode != 0)
    this.READ_FLAG(&uiCode, "deblocking_filter_control_present_flag")
    pcPPS.SetDeblockingFilterControlPresentFlag(uiCode != 0)
    if pcPPS.GetDeblockingFilterControlPresentFlag() {
        this.READ_FLAG(&uiCode, "deblocking_filter_override_enabled_flag")
        pcPPS.SetDeblockingFilterOverrideEnabledFlag(uiCode != 0)
        this.READ_FLAG(&uiCode, "pps_disable_deblocking_filter_flag")
        pcPPS.SetPicDisableDeblockingFilterFlag(uiCode != 0)
        if !pcPPS.GetPicDisableDeblockingFilterFlag() {
            this.READ_SVLC(&iCode, "pps_beta_offset_div2")
            pcPPS.SetDeblockingFilterBetaOffsetDiv2(iCode)
            this.READ_SVLC(&iCode, "pps_tc_offset_div2")
            pcPPS.SetDeblockingFilterTcOffsetDiv2(iCode)
        }
    }
    this.READ_FLAG(&uiCode, "pps_scaling_list_data_present_flag")
    pcPPS.SetScalingListPresentFlag(uiCode != 0)
    if pcPPS.GetScalingListPresentFlag() {
        this.ParseScalingList(pcPPS.GetScalingList())
    }

    this.READ_FLAG(&uiCode, "lists_modification_present_flag")
    pcPPS.SetListsModificationPresentFlag(uiCode != 0)

    this.READ_UVLC(&uiCode, "log2_parallel_merge_level_minus2")
    pcPPS.SetLog2ParallelMergeLevelMinus2(uiCode)

    //this.READ_CODE(3, &uiCode, "num_extra_slice_header_bits")
    //pcPPS.SetNumExtraSliceHeaderBits(int(uiCode))

    this.READ_FLAG(&uiCode, "slice_segment_header_extension_present_flag")
    pcPPS.SetSliceHeaderExtensionPresentFlag(uiCode != 0)

    this.READ_FLAG(&uiCode, "pps_extension_flag")
    if uiCode != 0 {
        for this.xMoreRbspData() {
            this.READ_FLAG(&uiCode, "pps_extension_data_flag")
        }
    }
}
func (this *TDecCavlc) ParseVUI(pcVUI *TLibCommon.TComVUI, pcSPS *TLibCommon.TComSPS) {
    //#if ENC_DEC_TRACE
    //fprintf( g_hTrace, "----------- vui_parameters -----------\n");
    this.xTraceVUIHeader(pcVUI)
    //#endif
    var uiCode uint

    this.READ_FLAG(&uiCode, "aspect_ratio_info_present_flag")
    pcVUI.SetAspectRatioInfoPresentFlag(uiCode != 0)
    if pcVUI.GetAspectRatioInfoPresentFlag() {
        this.READ_CODE(8, &uiCode, "aspect_ratio_idc")
        pcVUI.SetAspectRatioIdc(int(uiCode))
        if pcVUI.GetAspectRatioIdc() == 255 {
            this.READ_CODE(16, &uiCode, "sar_width")
            pcVUI.SetSarWidth(int(uiCode))
            this.READ_CODE(16, &uiCode, "sar_height")
            pcVUI.SetSarHeight(int(uiCode))
        }
    }

    this.READ_FLAG(&uiCode, "overscan_info_present_flag")
    pcVUI.SetOverscanInfoPresentFlag(uiCode != 0)
    if pcVUI.GetOverscanInfoPresentFlag() {
        this.READ_FLAG(&uiCode, "overscan_appropriate_flag")
        pcVUI.SetOverscanAppropriateFlag(uiCode != 0)
    }

    this.READ_FLAG(&uiCode, "video_signal_type_present_flag")
    pcVUI.SetVideoSignalTypePresentFlag(uiCode != 0)
    if pcVUI.GetVideoSignalTypePresentFlag() {
        this.READ_CODE(3, &uiCode, "video_format")
        pcVUI.SetVideoFormat(int(uiCode))
        this.READ_FLAG(&uiCode, "video_full_range_flag")
        pcVUI.SetVideoFullRangeFlag(uiCode != 0)
        this.READ_FLAG(&uiCode, "colour_description_present_flag")
        pcVUI.SetColourDescriptionPresentFlag(uiCode != 0)
        if pcVUI.GetColourDescriptionPresentFlag() {
            this.READ_CODE(8, &uiCode, "colour_primaries")
            pcVUI.SetColourPrimaries(int(uiCode))
            this.READ_CODE(8, &uiCode, "transfer_characteristics")
            pcVUI.SetTransferCharacteristics(int(uiCode))
            this.READ_CODE(8, &uiCode, "matrix_coefficients")
            pcVUI.SetMatrixCoefficients(int(uiCode))
        }
    }

    this.READ_FLAG(&uiCode, "chroma_loc_info_present_flag")
    pcVUI.SetChromaLocInfoPresentFlag(uiCode != 0)
    if pcVUI.GetChromaLocInfoPresentFlag() {
        this.READ_UVLC(&uiCode, "chroma_sample_loc_type_top_field")
        pcVUI.SetChromaSampleLocTypeTopField(int(uiCode))
        this.READ_UVLC(&uiCode, "chroma_sample_loc_type_bottom_field")
        pcVUI.SetChromaSampleLocTypeBottomField(int(uiCode))
    }

    this.READ_FLAG(&uiCode, "neutral_chroma_indication_flag")
    pcVUI.SetNeutralChromaIndicationFlag(uiCode != 0)

    this.READ_FLAG(&uiCode, "field_seq_flag")
    pcVUI.SetFieldSeqFlag(uiCode != 0)
    //assert(pcVUI.GetFieldSeqFlag() == false);        // not supported yet

    this.READ_FLAG(&uiCode, "frame_field_info_present_flag")
    pcVUI.SetFrameFieldInfoPresentFlag(uiCode != 0)

    this.READ_FLAG(&uiCode, "default_display_window_flag") //assert(uiCode == 0);
    if uiCode != 0 {
        defDisp := pcVUI.GetDefaultDisplayWindow()
        this.READ_UVLC(&uiCode, "def_disp_win_left_offset")
        defDisp.SetWindowLeftOffset(int(uiCode * uint(pcSPS.GetWinUnitX(pcSPS.GetChromaFormatIdc()))))
        this.READ_UVLC(&uiCode, "def_disp_win_right_offset")
        defDisp.SetWindowRightOffset(int(uiCode * uint(pcSPS.GetWinUnitX(pcSPS.GetChromaFormatIdc()))))
        this.READ_UVLC(&uiCode, "def_disp_win_top_offset")
        defDisp.SetWindowTopOffset(int(uiCode * uint(pcSPS.GetWinUnitY(pcSPS.GetChromaFormatIdc()))))
        this.READ_UVLC(&uiCode, "def_disp_win_bottom_offset")
        defDisp.SetWindowBottomOffset(int(uiCode * uint(pcSPS.GetWinUnitY(pcSPS.GetChromaFormatIdc()))))
    }

    timingInfo := pcVUI.GetTimingInfo()
    this.READ_FLAG(&uiCode, "vui_timing_info_present_flag")
    timingInfo.SetTimingInfoPresentFlag(uiCode != 0)
    if timingInfo.GetTimingInfoPresentFlag() {
        this.READ_CODE(32, &uiCode, "vui_num_units_in_tick")
        timingInfo.SetNumUnitsInTick(uiCode)
        this.READ_CODE(32, &uiCode, "vui_time_scale")
        timingInfo.SetTimeScale(uiCode)
        this.READ_FLAG(&uiCode, "vui_poc_proportional_to_timing_flag")
        timingInfo.SetPocProportionalToTimingFlag(uiCode != 0)
        if timingInfo.GetPocProportionalToTimingFlag() {
            this.READ_UVLC(&uiCode, "vui_num_ticks_poc_diff_one_minus1")
            timingInfo.SetNumTicksPocDiffOneMinus1(int(uiCode))
        }
        this.READ_FLAG(&uiCode, "hrd_parameters_present_flag")
        pcVUI.SetHrdParametersPresentFlag(uiCode != 0)
        if pcVUI.GetHrdParametersPresentFlag() {
            this.ParseHrdParameters(pcVUI.GetHrdParameters(), true, pcSPS.GetMaxTLayers()-1)
        }
    }

    this.READ_FLAG(&uiCode, "bitstream_restriction_flag")
    pcVUI.SetBitstreamRestrictionFlag(uiCode != 0)
    if pcVUI.GetBitstreamRestrictionFlag() {
        this.READ_FLAG(&uiCode, "tiles_fixed_structure_flag")
        pcVUI.SetTilesFixedStructureFlag(uiCode != 0)
        this.READ_FLAG(&uiCode, "motion_vectors_over_pic_boundaries_flag")
        pcVUI.SetMotionVectorsOverPicBoundariesFlag(uiCode != 0)

        this.READ_FLAG(&uiCode, "restricted_ref_pic_lists_flag")
        pcVUI.SetRestrictedRefPicListsFlag(uiCode != 0)

        this.READ_UVLC(&uiCode, "min_spatial_segmentation_idc")
        pcVUI.SetMinSpatialSegmentationIdc(int(uiCode))

        this.READ_UVLC(&uiCode, "max_bytes_per_pic_denom")
        pcVUI.SetMaxBytesPerPicDenom(int(uiCode))
        this.READ_UVLC(&uiCode, "max_bits_per_mincu_denom")
        pcVUI.SetMaxBitsPerMinCuDenom(int(uiCode))
        this.READ_UVLC(&uiCode, "log2_max_mv_length_horizontal")
        pcVUI.SetLog2MaxMvLengthHorizontal(int(uiCode))
        this.READ_UVLC(&uiCode, "log2_max_mv_length_vertical")
        pcVUI.SetLog2MaxMvLengthVertical(int(uiCode))
    }
}
func (this *TDecCavlc) ParseSEI(sei *TLibCommon.SEImessages) {
}
func (this *TDecCavlc) ParsePTL(rpcPTL *TLibCommon.TComPTL, profilePresentFlag bool, maxNumSubLayersMinus1 int) {
    var uiCode uint
    if profilePresentFlag {
        this.ParseProfileTier(rpcPTL.GetGeneralPTL())
    }
    this.READ_CODE(8, &uiCode, "general_level_idc")
    rpcPTL.GetGeneralPTL().SetLevelIdc(int(uiCode))

    for i := 0; i < maxNumSubLayersMinus1; i++ {
        if profilePresentFlag {
            this.READ_FLAG(&uiCode, "sub_layer_profile_present_flag[i]")
            rpcPTL.SetSubLayerProfilePresentFlag(i, uiCode != 0)
        }
        this.READ_FLAG(&uiCode, "sub_layer_level_present_flag[i]")
        rpcPTL.SetSubLayerLevelPresentFlag(i, uiCode != 0)
    }
    if maxNumSubLayersMinus1 > 0 {
        for i := maxNumSubLayersMinus1; i < 8; i++ {
            this.READ_CODE(2, &uiCode, "reserved_zero_2bits")
            //assert(uiCode == 0);
        }
    }

    for i := 0; i < maxNumSubLayersMinus1; i++ {
        if profilePresentFlag && rpcPTL.GetSubLayerProfilePresentFlag(i) {
            this.ParseProfileTier(rpcPTL.GetSubLayerPTL(i))
        }
        if rpcPTL.GetSubLayerLevelPresentFlag(i) {
            this.READ_CODE(8, &uiCode, "sub_layer_level_idc[i]")
            rpcPTL.GetSubLayerPTL(i).SetLevelIdc(int(uiCode))
        }
    }
}
func (this *TDecCavlc) ParseProfileTier(ptl *TLibCommon.ProfileTierLevel) {
    var uiCode uint
    this.READ_CODE(2, &uiCode, "XXX_profile_space[]")
    ptl.SetProfileSpace(int(uiCode))
    this.READ_FLAG(&uiCode, "XXX_tier_flag[]")
    ptl.SetTierFlag(uiCode != 0)
    this.READ_CODE(5, &uiCode, "XXX_profile_idc[]")
    ptl.SetProfileIdc(int(uiCode))
    for j := 0; j < 32; j++ {
        this.READ_FLAG(&uiCode, "XXX_profile_compatibility_flag[][j]")
        ptl.SetProfileCompatibilityFlag(j, uiCode != 0)
    }

    this.READ_FLAG(&uiCode, "general_progressive_source_flag")
    ptl.SetProgressiveSourceFlag(uiCode != 0)

    this.READ_FLAG(&uiCode, "general_interlaced_source_flag")
    ptl.SetInterlacedSourceFlag(uiCode != 0)

    this.READ_FLAG(&uiCode, "general_non_packed_constraint_flag")
    ptl.SetNonPackedConstraintFlag(uiCode != 0)

    this.READ_FLAG(&uiCode, "general_frame_only_constraint_flag")
    ptl.SetFrameOnlyConstraintFlag(uiCode != 0)

    this.READ_CODE(16, &uiCode, "XXX_reserved_zero_44bits[0..15]")
    this.READ_CODE(16, &uiCode, "XXX_reserved_zero_44bits[16..31]")
    this.READ_CODE(12, &uiCode, "XXX_reserved_zero_44bits[32..43]")
}

func (this *TDecCavlc) ParseHrdParameters(hrd *TLibCommon.TComHRD, commonInfPresentFlag bool, maxNumSubLayersMinus1 uint) {
    var uiCode uint
    if commonInfPresentFlag {
        this.READ_FLAG(&uiCode, "nal_hrd_parameters_present_flag")
        hrd.SetNalHrdParametersPresentFlag(uiCode == 1)
        this.READ_FLAG(&uiCode, "vcl_hrd_parameters_present_flag")
        hrd.SetVclHrdParametersPresentFlag(uiCode == 1)
        if hrd.GetNalHrdParametersPresentFlag() || hrd.GetVclHrdParametersPresentFlag() {
            this.READ_FLAG(&uiCode, "sub_pic_cpb_params_present_flag")
            hrd.SetSubPicCpbParamsPresentFlag(uiCode == 1)
            if hrd.GetSubPicCpbParamsPresentFlag() {
                this.READ_CODE(8, &uiCode, "tick_divisor_minus2")
                hrd.SetTickDivisorMinus2(uiCode)
                this.READ_CODE(5, &uiCode, "du_cpb_removal_delay_length_minus1")
                hrd.SetDuCpbRemovalDelayLengthMinus1(uiCode)
                this.READ_FLAG(&uiCode, "sub_pic_cpb_params_in_pic_timing_sei_flag")
                hrd.SetSubPicCpbParamsInPicTimingSEIFlag(uiCode == 1)
                this.READ_CODE(5, &uiCode, "dpb_output_delay_du_length_minus1")
                hrd.SetDpbOutputDelayDuLengthMinus1(uiCode)
            }
            this.READ_CODE(4, &uiCode, "bit_rate_scale")
            hrd.SetBitRateScale(uiCode)
            this.READ_CODE(4, &uiCode, "cpb_size_scale")
            hrd.SetCpbSizeScale(uiCode)
            if hrd.GetSubPicCpbParamsPresentFlag() {
                this.READ_CODE(4, &uiCode, "cpb_size_du_scale")
                hrd.SetDuCpbSizeScale(uiCode)
            }
            this.READ_CODE(5, &uiCode, "initial_cpb_removal_delay_length_minus1")
            hrd.SetInitialCpbRemovalDelayLengthMinus1(uiCode)
            this.READ_CODE(5, &uiCode, "au_cpb_removal_delay_length_minus1")
            hrd.SetCpbRemovalDelayLengthMinus1(uiCode)
            this.READ_CODE(5, &uiCode, "dpb_output_delay_length_minus1")
            hrd.SetDpbOutputDelayLengthMinus1(uiCode)
        }
    }
    var i, j, nalOrVcl int
    for i = 0; i <= int(maxNumSubLayersMinus1); i++ {
        this.READ_FLAG(&uiCode, "fixed_pic_rate_general_flag")
        hrd.SetFixedPicRateFlag(i, uiCode == 1)
        if !hrd.GetFixedPicRateFlag(i) {
            this.READ_FLAG(&uiCode, "fixed_pic_rate_within_cvs_flag")
            hrd.SetFixedPicRateWithinCvsFlag(i, uiCode == 1)
        } else {
            hrd.SetFixedPicRateWithinCvsFlag(i, true)
        }

        hrd.SetLowDelayHrdFlag(i, false) // Infered to be 0 when not present
        hrd.SetCpbCntMinus1(i, 0)        // Infered to be 0 when not present

        if hrd.GetFixedPicRateWithinCvsFlag(i) {
            this.READ_UVLC(&uiCode, "elemental_duration_in_tc_minus1")
            hrd.SetPicDurationInTcMinus1(i, uiCode)
        } else {
            this.READ_FLAG(&uiCode, "low_delay_hrd_flag")
            hrd.SetLowDelayHrdFlag(i, uiCode == 1)
        }
        if !hrd.GetLowDelayHrdFlag(i) {
            this.READ_UVLC(&uiCode, "cpb_cnt_minus1")
            hrd.SetCpbCntMinus1(i, uiCode)
        }

        for nalOrVcl = 0; nalOrVcl < 2; nalOrVcl++ {
            if ((nalOrVcl == 0) && (hrd.GetNalHrdParametersPresentFlag())) ||
                ((nalOrVcl == 1) && (hrd.GetVclHrdParametersPresentFlag())) {
                for j = 0; j <= int(hrd.GetCpbCntMinus1(i)); j++ {
                    this.READ_UVLC(&uiCode, "bit_rate_value_minus1")
                    hrd.SetBitRateValueMinus1(i, j, nalOrVcl, uiCode)
                    this.READ_UVLC(&uiCode, "cpb_size_value_minus1")
                    hrd.SetCpbSizeValueMinus1(i, j, nalOrVcl, uiCode)
                    if hrd.GetSubPicCpbParamsPresentFlag() {
                        this.READ_UVLC(&uiCode, "bit_rate_du_value_minus1")
                        hrd.SetDuBitRateValueMinus1(i, j, nalOrVcl, uiCode)
                        this.READ_UVLC(&uiCode, "cpb_size_du_value_minus1")
                        hrd.SetDuCpbSizeValueMinus1(i, j, nalOrVcl, uiCode)
                    }
                    this.READ_FLAG(&uiCode, "cbr_flag")
                    hrd.SetCbrFlag(i, j, nalOrVcl, uiCode == 1)
                }
            }
        }
    }
}

func (this *TDecCavlc) ParseSliceHeader(rpcSlice *TLibCommon.TComSlice, parameterSetManager *TLibCommon.ParameterSetManager) bool { //Decoder
    var uiCode uint
    var iCode int

    //#if ENC_DEC_TRACE
    this.xTraceSliceHeader(rpcSlice)
    //#endif
    //TComPPS* pps = NULL;
    //TComSPS* sps = NULL;

    var firstSliceSegmentInPic uint
    this.READ_FLAG(&firstSliceSegmentInPic, "first_slice_segment_in_pic_flag")
    if rpcSlice.GetRapPicFlag() {
        this.READ_FLAG(&uiCode, "no_output_of_prior_pics_flag") //ignored
    }
    this.READ_UVLC(&uiCode, "slice_pic_parameter_set_id")
    rpcSlice.SetPPSId(int(uiCode))
    pps := parameterSetManager.GetPPS(int(uiCode))
    //!KS: need to add error handling code here, if PPS is not available
    //assert(pps!=0);
    sps := parameterSetManager.GetSPS(pps.GetSPSId())
    //!KS: need to add error handling code here, if SPS is not available
    //assert(sps!=0);
    rpcSlice.SetSPS(sps)
    rpcSlice.SetPPS(pps)

    if pps.GetDependentSliceSegmentsEnabledFlag() && firstSliceSegmentInPic == 0 {
        this.READ_FLAG(&uiCode, "dependent_slice_segment_flag")
        rpcSlice.SetDependentSliceSegmentFlag(uiCode != 0)
    } else {
        rpcSlice.SetDependentSliceSegmentFlag(false)
    }

    numCTUs := ((sps.GetPicWidthInLumaSamples() + sps.GetMaxCUWidth() - 1) / sps.GetMaxCUWidth()) * ((sps.GetPicHeightInLumaSamples() + sps.GetMaxCUHeight() - 1) / sps.GetMaxCUHeight())
    maxParts := uint(1 << (sps.GetMaxCUDepth() << 1))
    sliceSegmentAddress := uint(0)
    bitsSliceSegmentAddress := uint(0)
    for numCTUs > (1 << bitsSliceSegmentAddress) {
        bitsSliceSegmentAddress++
    }

    if firstSliceSegmentInPic == 0 {
        this.READ_CODE(bitsSliceSegmentAddress, &sliceSegmentAddress, "slice_segment_address")
    }
    //set uiCode to equal slice start address (or dependent slice start address)
    startCuAddress := maxParts * sliceSegmentAddress
    rpcSlice.SetSliceSegmentCurStartCUAddr(startCuAddress)
    rpcSlice.SetSliceSegmentCurEndCUAddr(numCTUs * maxParts)

    if rpcSlice.GetDependentSliceSegmentFlag() {
        rpcSlice.SetNextSlice(false)
        rpcSlice.SetNextSliceSegment(true)
    } else {
        rpcSlice.SetNextSlice(true)
        rpcSlice.SetNextSliceSegment(false)

        rpcSlice.SetSliceCurStartCUAddr(startCuAddress)
        rpcSlice.SetSliceCurEndCUAddr(numCTUs * maxParts)
    }

    if !rpcSlice.GetDependentSliceSegmentFlag() {
        for i := 0; i < rpcSlice.GetPPS().GetNumExtraSliceHeaderBits(); i++ {
            this.READ_FLAG(&uiCode, "slice_reserved_undetermined_flag[]") // ignored
        }

        this.READ_UVLC(&uiCode, "slice_type")
        rpcSlice.SetSliceType(TLibCommon.SliceType(uiCode))
        if pps.GetOutputFlagPresentFlag() {
            this.READ_FLAG(&uiCode, "pic_output_flag")
            rpcSlice.SetPicOutputFlag(uiCode != 0)
        } else {
            rpcSlice.SetPicOutputFlag(true)
        }
        // in the first version chroma_format_idc is equal to one, thus colour_plane_id will not be present
        //assert (sps.GetChromaFormatIdc() == 1 );
        // if( separate_colour_plane_flag  ==  1 )
        //   colour_plane_id                                      u(2)

        if rpcSlice.GetIdrPicFlag() {
            rpcSlice.SetPOC(0)
            rps := rpcSlice.GetLocalRPS()
            rps.SetNumberOfNegativePictures(0)
            rps.SetNumberOfPositivePictures(0)
            rps.SetNumberOfLongtermPictures(0)
            rps.SetNumberOfPictures(0)
            rpcSlice.SetRPS(rps)
        } else {
            this.READ_CODE(sps.GetBitsForPOC(), &uiCode, "pic_order_cnt_lsb")
            iPOClsb := int(uiCode)
            iPrevPOC := int(rpcSlice.GetPrevPOC())
            iMaxPOClsb := int(1 << sps.GetBitsForPOC())
            //iPrevPOClsb := int(iPrevPOC % iMaxPOClsb)
            iPrevPOClsb := int(-1);
      		multipler := int(0);
		    for iPrevPOClsb<0 {
		        iPrevPOClsb = (iPrevPOC+iMaxPOClsb*multipler)%iMaxPOClsb;
		        multipler++;
		    }
            iPrevPOCmsb := int(iPrevPOC - iPrevPOClsb)
            var iPOCmsb int
            //fmt.Printf("iPrevPOC %d, iPOClsb %d, iPrevPOClsb %d, iMaxPOClsb %d\n", iPrevPOC, iPOClsb, iPrevPOClsb, iMaxPOClsb);
            if (iPOClsb < iPrevPOClsb) && ((iPrevPOClsb - iPOClsb) >= (iMaxPOClsb / 2)) {
                iPOCmsb = iPrevPOCmsb + iMaxPOClsb
            } else if (iPOClsb > iPrevPOClsb) && ((iPOClsb - iPrevPOClsb) > (iMaxPOClsb / 2)) {
                iPOCmsb = iPrevPOCmsb - iMaxPOClsb
            } else {
                iPOCmsb = iPrevPOCmsb
            }
            //fmt.Printf("iPOCmsb %d + iPOClsb %d\n", iPOCmsb, iPOClsb);
            if rpcSlice.GetNalUnitType() == TLibCommon.NAL_UNIT_CODED_SLICE_BLA ||
                rpcSlice.GetNalUnitType() == TLibCommon.NAL_UNIT_CODED_SLICE_BLANT ||
                rpcSlice.GetNalUnitType() == TLibCommon.NAL_UNIT_CODED_SLICE_BLA_N_LP {
                // For BLA picture types, POCmsb is set to 0.
                iPOCmsb = 0
            }
            rpcSlice.SetPOC(iPOCmsb + iPOClsb)
			//fmt.Printf("iPOC %d = iPOCmsb %d + iPOClsb %d\n", rpcSlice.GetPOC(), iPOCmsb, iPOClsb);

            var rps *TLibCommon.TComReferencePictureSet
            this.READ_FLAG(&uiCode, "short_term_ref_pic_set_sps_flag")
            if uiCode == 0 { // use short-term reference picture set explicitly signalled in slice header
                rps = rpcSlice.GetLocalRPS()
                this.ParseShortTermRefPicSet(sps, rps, sps.GetRPSList().GetNumberOfReferencePictureSets())
                rpcSlice.SetRPS(rps)
            } else { // use reference to short-term reference picture set in PPS
                numBits := uint(0)
                for (1 << numBits) < rpcSlice.GetSPS().GetRPSList().GetNumberOfReferencePictureSets() {
                    numBits++
                }
                if numBits > 0 {
                    this.READ_CODE(numBits, &uiCode, "short_term_ref_pic_set_idx")
                } else {
                    uiCode = 0
                }
                rpcSlice.SetRPS(sps.GetRPSList().GetReferencePictureSet(int(uiCode)))

                rps = rpcSlice.GetRPS()
            }
            if sps.GetLongTermRefsPresent() {
                offset := rps.GetNumberOfNegativePictures() + rps.GetNumberOfPositivePictures()
                numOfLtrp := uint(0)
                numLtrpInSPS := uint(0)
                if rpcSlice.GetSPS().GetNumLongTermRefPicSPS() > 0 {
                    this.READ_UVLC(&uiCode, "num_long_term_sps")
                    numLtrpInSPS = uiCode
                    numOfLtrp += numLtrpInSPS
                    rps.SetNumberOfLongtermPictures(int(numOfLtrp))
                }
                bitsForLtrpInSPS := uint(0)
                for rpcSlice.GetSPS().GetNumLongTermRefPicSPS() > (1 << bitsForLtrpInSPS) {
                    bitsForLtrpInSPS++
                }
                this.READ_UVLC(&uiCode, "num_long_term_pics")
                rps.SetNumberOfLongtermPictures(int(uiCode))
                numOfLtrp += uiCode
                rps.SetNumberOfLongtermPictures(int(numOfLtrp))
                maxPicOrderCntLSB := 1 << rpcSlice.GetSPS().GetBitsForPOC()
                //prevLSB := 0
                prevDeltaMSB := 0
                deltaPocMSBCycleLT := 0
                j := offset + rps.GetNumberOfLongtermPictures() - 1
                for k := uint(0); k < numOfLtrp; k++ {
                    var pocLsbLt int
                    if k < numLtrpInSPS {
                        uiCode = 0
                        if bitsForLtrpInSPS > 0 {
                            this.READ_CODE(bitsForLtrpInSPS, &uiCode, "lt_idx_sps[i]")
                        }

                        usedByCurrFromSPS := rpcSlice.GetSPS().GetUsedByCurrPicLtSPSFlag(int(uiCode))

                        pocLsbLt = int(rpcSlice.GetSPS().GetLtRefPicPocLsbSps(uiCode))
                        rps.SetUsed(j, usedByCurrFromSPS)
                    } else {
                        this.READ_CODE(rpcSlice.GetSPS().GetBitsForPOC(), &uiCode, "poc_lsb_lt")
                        pocLsbLt = int(uiCode)
                        this.READ_FLAG(&uiCode, "used_by_curr_pic_lt_flag")
                        rps.SetUsed(j, uiCode != 0)
                    }
                    this.READ_FLAG(&uiCode, "delta_poc_msb_present_flag")
                    mSBPresentFlag := uiCode != 0
                    if mSBPresentFlag {
                        this.READ_UVLC(&uiCode, "delta_poc_msb_cycle_lt[i]")
                        deltaFlag := false
                        //            First LTRP                               || First LTRP from SH           || curr LSB    != prev LSB
                        if (j == offset+rps.GetNumberOfLongtermPictures()-1) || (j == offset+int(numOfLtrp-numLtrpInSPS)-1) {//|| (pocLsbLt != prevLSB) {
                            deltaFlag = true
                        }
                        if deltaFlag {
                            deltaPocMSBCycleLT = int(uiCode)
                        } else {
                            deltaPocMSBCycleLT = int(uiCode) + prevDeltaMSB
                        }

                        pocLTCurr := rpcSlice.GetPOC() - deltaPocMSBCycleLT*maxPicOrderCntLSB - iPOClsb + pocLsbLt
                        rps.SetPOC(j, pocLTCurr)
                        rps.SetDeltaPOC(j, -rpcSlice.GetPOC()+pocLTCurr)
                        rps.SetCheckLTMSBPresent(j, true)
                    } else {
                        if (j == offset+rps.GetNumberOfLongtermPictures()-1) || (j == offset+int(numOfLtrp-numLtrpInSPS)-1) {//|| (pocLsbLt != prevLSB) {
                            deltaPocMSBCycleLT = 0
                        }
                        rps.SetPOC(j, pocLsbLt)
                        rps.SetDeltaPOC(j, -rpcSlice.GetPOC()+pocLsbLt)
                        rps.SetCheckLTMSBPresent(j, false)
                    }
                    //prevLSB = pocLsbLt
                    prevDeltaMSB = deltaPocMSBCycleLT
                    j--
                }
                offset += rps.GetNumberOfLongtermPictures()
                rps.SetNumberOfPictures(offset)
            }
            /*if TLibCommon.G_uiPicNo>380 && firstSliceSegmentInPic==1 {
		        for i:=0;i<rps.GetNumberOfPictures();i++ {
		          fmt.Printf("%d ", rps.GetPOC(i));
		        }
		      }*/

            if rpcSlice.GetNalUnitType() == TLibCommon.NAL_UNIT_CODED_SLICE_BLA ||
                rpcSlice.GetNalUnitType() == TLibCommon.NAL_UNIT_CODED_SLICE_BLANT ||
                rpcSlice.GetNalUnitType() == TLibCommon.NAL_UNIT_CODED_SLICE_BLA_N_LP {
                // In the case of BLA picture types, rps data is read from slice header but ignored
                rps = rpcSlice.GetLocalRPS()
                rps.SetNumberOfNegativePictures(0)
                rps.SetNumberOfPositivePictures(0)
                rps.SetNumberOfLongtermPictures(0)
                rps.SetNumberOfPictures(0)
                rpcSlice.SetRPS(rps)
            }

            if rpcSlice.GetSPS().GetTMVPFlagsPresent() {
                this.READ_FLAG(&uiCode, "slice_temporal_mvp_enable_flag")
                rpcSlice.SetEnableTMVPFlag(uiCode != 0)
            } else {
                rpcSlice.SetEnableTMVPFlag(false)
            }
        }
        if sps.GetUseSAO() {
            this.READ_FLAG(&uiCode, "slice_sao_luma_flag")
            rpcSlice.SetSaoEnabledFlag(uiCode != 0)
            this.READ_FLAG(&uiCode, "slice_sao_chroma_flag")
            rpcSlice.SetSaoEnabledFlagChroma(uiCode != 0)
        }

        if rpcSlice.GetIdrPicFlag() {
            rpcSlice.SetEnableTMVPFlag(false)
        }

        if !rpcSlice.IsIntra() {
            this.READ_FLAG(&uiCode, "num_ref_idx_active_override_flag")
            if uiCode != 0 {
                this.READ_UVLC(&uiCode, "num_ref_idx_l0_active_minus1")
                rpcSlice.SetNumRefIdx(TLibCommon.REF_PIC_LIST_0, int(uiCode)+1)
                if rpcSlice.IsInterB() {
                    this.READ_UVLC(&uiCode, "num_ref_idx_l1_active_minus1")
                    rpcSlice.SetNumRefIdx(TLibCommon.REF_PIC_LIST_1, int(uiCode)+1)
                } else {
                    rpcSlice.SetNumRefIdx(TLibCommon.REF_PIC_LIST_1, 0)
                }
            } else {
                rpcSlice.SetNumRefIdx(TLibCommon.REF_PIC_LIST_0, int(rpcSlice.GetPPS().GetNumRefIdxL0DefaultActive()))
                if rpcSlice.IsInterB() {
                    rpcSlice.SetNumRefIdx(TLibCommon.REF_PIC_LIST_1, int(rpcSlice.GetPPS().GetNumRefIdxL1DefaultActive()))
                } else {
                    rpcSlice.SetNumRefIdx(TLibCommon.REF_PIC_LIST_1, 0)
                }
            }
        }
        // }
        refPicListModification := rpcSlice.GetRefPicListModification()
        if !rpcSlice.IsIntra() {
            if !rpcSlice.GetPPS().GetListsModificationPresentFlag() || rpcSlice.GetNumRpsCurrTempList() <= 1 {
                refPicListModification.SetRefPicListModificationFlagL0(false)
            } else {
                this.READ_FLAG(&uiCode, "ref_pic_list_modification_flag_l0")
                refPicListModification.SetRefPicListModificationFlagL0(uiCode != 0)
            }

            if refPicListModification.GetRefPicListModificationFlagL0() {
                uiCode = 0
                i := 0
                numRpsCurrTempList0 := rpcSlice.GetNumRpsCurrTempList()
                if numRpsCurrTempList0 > 1 {
                    length := 1
                    numRpsCurrTempList0--
                    numRpsCurrTempList0 >>= 1
                    for numRpsCurrTempList0 != 0 {
                        length++
                        numRpsCurrTempList0 >>= 1
                    }
                    for i = 0; i < rpcSlice.GetNumRefIdx(TLibCommon.REF_PIC_LIST_0); i++ {
                        this.READ_CODE(uint(length), &uiCode, "list_entry_l0")
                        refPicListModification.SetRefPicSetIdxL0(uint(i), uiCode)
                    }
                } else {
                    for i = 0; i < rpcSlice.GetNumRefIdx(TLibCommon.REF_PIC_LIST_0); i++ {
                        refPicListModification.SetRefPicSetIdxL0(uint(i), 0)
                    }
                }
            }
        } else {
            refPicListModification.SetRefPicListModificationFlagL0(false)
        }
        if rpcSlice.IsInterB() {
            if !rpcSlice.GetPPS().GetListsModificationPresentFlag() || rpcSlice.GetNumRpsCurrTempList() <= 1 {
                refPicListModification.SetRefPicListModificationFlagL1(false)
            } else {
                this.READ_FLAG(&uiCode, "ref_pic_list_modification_flag_l1")
                refPicListModification.SetRefPicListModificationFlagL1(uiCode != 0)
            }
            if refPicListModification.GetRefPicListModificationFlagL1() {
                uiCode = 0
                i := 0
                numRpsCurrTempList1 := rpcSlice.GetNumRpsCurrTempList()
                if numRpsCurrTempList1 > 1 {
                    length := 1
                    numRpsCurrTempList1--
                    numRpsCurrTempList1 >>= 1
                    for numRpsCurrTempList1 != 0 {
                        length++
                        numRpsCurrTempList1 >>= 1
                    }
                    for i = 0; i < rpcSlice.GetNumRefIdx(TLibCommon.REF_PIC_LIST_1); i++ {
                        this.READ_CODE(uint(length), &uiCode, "list_entry_l1")
                        refPicListModification.SetRefPicSetIdxL1(uint(i), uiCode)
                    }
                } else {
                    for i = 0; i < rpcSlice.GetNumRefIdx(TLibCommon.REF_PIC_LIST_1); i++ {
                        refPicListModification.SetRefPicSetIdxL1(uint(i), 0)
                    }
                }
            }
        } else {
            refPicListModification.SetRefPicListModificationFlagL1(false)
        }
        if rpcSlice.IsInterB() {
            this.READ_FLAG(&uiCode, "mvd_l1_zero_flag")
            rpcSlice.SetMvdL1ZeroFlag(uiCode != 0)
        }

        rpcSlice.SetCabacInitFlag(false) // default
        if pps.GetCabacInitPresentFlag() && !rpcSlice.IsIntra() {
            this.READ_FLAG(&uiCode, "cabac_init_flag")
            rpcSlice.SetCabacInitFlag(uiCode != 0)
        }

        if rpcSlice.GetEnableTMVPFlag() {
            if rpcSlice.GetSliceType() == TLibCommon.B_SLICE {
                this.READ_FLAG(&uiCode, "collocated_from_l0_flag")
                rpcSlice.SetColFromL0Flag(uiCode)
            } else {
                rpcSlice.SetColFromL0Flag(1)
            }

            if rpcSlice.GetSliceType() != TLibCommon.I_SLICE &&
                ((rpcSlice.GetColFromL0Flag() == 1 && rpcSlice.GetNumRefIdx(TLibCommon.REF_PIC_LIST_0) > 1) ||
                    (rpcSlice.GetColFromL0Flag() == 0 && rpcSlice.GetNumRefIdx(TLibCommon.REF_PIC_LIST_1) > 1)) {
                this.READ_UVLC(&uiCode, "collocated_ref_idx")
                rpcSlice.SetColRefIdx(uiCode)
            } else {
                rpcSlice.SetColRefIdx(0)
            }
        }
        if (pps.GetUseWP() && rpcSlice.GetSliceType() == TLibCommon.P_SLICE) ||
            (pps.GetWPBiPred() && rpcSlice.GetSliceType() == TLibCommon.B_SLICE) {
            this.xParsePredWeightTable(rpcSlice)
            rpcSlice.InitWpScaling()
        }
        if !rpcSlice.IsIntra() {
            this.READ_UVLC(&uiCode, "five_minus_max_num_merge_cand")
            rpcSlice.SetMaxNumMergeCand(TLibCommon.MRG_MAX_NUM_CANDS - uiCode)
        }

        this.READ_SVLC(&iCode, "slice_qp_delta")
        rpcSlice.SetSliceQp(26 + pps.GetPicInitQPMinus26() + iCode)

        //assert( rpcSlice.GetSliceQp() >= -sps.GetQpBDOffsetY() );
        //assert( rpcSlice.GetSliceQp() <=  51 );

        if rpcSlice.GetPPS().GetSliceChromaQpFlag() {
            this.READ_SVLC(&iCode, "slice_qp_delta_cb")
            rpcSlice.SetSliceQpDeltaCb(iCode)
            //assert( rpcSlice.GetSliceQpDeltaCb() >= -12 );
            //assert( rpcSlice.GetSliceQpDeltaCb() <=  12 );
            //assert( (rpcSlice.GetPPS().GetChromaCbQpOffset() + rpcSlice.GetSliceQpDeltaCb()) >= -12 );
            //assert( (rpcSlice.GetPPS().GetChromaCbQpOffset() + rpcSlice.GetSliceQpDeltaCb()) <=  12 );

            this.READ_SVLC(&iCode, "slice_qp_delta_cr")
            rpcSlice.SetSliceQpDeltaCr(iCode)
            //assert( rpcSlice.GetSliceQpDeltaCr() >= -12 );
            //assert( rpcSlice.GetSliceQpDeltaCr() <=  12 );
            //assert( (rpcSlice.GetPPS().GetChromaCrQpOffset() + rpcSlice.GetSliceQpDeltaCr()) >= -12 );
            //assert( (rpcSlice.GetPPS().GetChromaCrQpOffset() + rpcSlice.GetSliceQpDeltaCr()) <=  12 );
        }

        if rpcSlice.GetPPS().GetDeblockingFilterControlPresentFlag() {
            if rpcSlice.GetPPS().GetDeblockingFilterOverrideEnabledFlag() {
                this.READ_FLAG(&uiCode, "deblocking_filter_override_flag")
                rpcSlice.SetDeblockingFilterOverrideFlag(uiCode != 0)
            } else {
                rpcSlice.SetDeblockingFilterOverrideFlag(false)
            }
            if rpcSlice.GetDeblockingFilterOverrideFlag() {
                this.READ_FLAG(&uiCode, "slice_disable_deblocking_filter_flag")
                rpcSlice.SetDeblockingFilterDisable(uiCode != 0)
                if !rpcSlice.GetDeblockingFilterDisable() {
                    this.READ_SVLC(&iCode, "beta_offset_div2")
                    rpcSlice.SetDeblockingFilterBetaOffsetDiv2(iCode)
                    this.READ_SVLC(&iCode, "tc_offset_div2")
                    rpcSlice.SetDeblockingFilterTcOffsetDiv2(iCode)
                }
            } else {
                rpcSlice.SetDeblockingFilterDisable(rpcSlice.GetPPS().GetPicDisableDeblockingFilterFlag())
                rpcSlice.SetDeblockingFilterBetaOffsetDiv2(rpcSlice.GetPPS().GetDeblockingFilterBetaOffsetDiv2())
                rpcSlice.SetDeblockingFilterTcOffsetDiv2(rpcSlice.GetPPS().GetDeblockingFilterTcOffsetDiv2())
            }
        } else {
            rpcSlice.SetDeblockingFilterDisable(false)
            rpcSlice.SetDeblockingFilterBetaOffsetDiv2(0)
            rpcSlice.SetDeblockingFilterTcOffsetDiv2(0)
        }

        var isSAOEnabled bool
        if !rpcSlice.GetSPS().GetUseSAO() {
            isSAOEnabled = false
        } else {
            isSAOEnabled = rpcSlice.GetSaoEnabledFlag() || rpcSlice.GetSaoEnabledFlagChroma()
        }
        isDBFEnabled := (!rpcSlice.GetDeblockingFilterDisable())

        if rpcSlice.GetPPS().GetLoopFilterAcrossSlicesEnabledFlag() && (isSAOEnabled || isDBFEnabled) {
            this.READ_FLAG(&uiCode, "slice_loop_filter_across_slices_enabled_flag")
        } else {
            uiCode = uint(TLibCommon.B2U(rpcSlice.GetPPS().GetLoopFilterAcrossSlicesEnabledFlag()))
        }
        rpcSlice.SetLFCrossSliceBoundaryFlag(uiCode != 0)
    }

    var entryPointOffset []uint;
    var numEntryPointOffsets, offsetLenMinus1 uint
    if pps.GetTilesEnabledFlag() || pps.GetEntropyCodingSyncEnabledFlag() {
        this.READ_UVLC(&numEntryPointOffsets, "num_entry_point_offsets")
        rpcSlice.SetNumEntryPointOffsets(int(numEntryPointOffsets))
        if numEntryPointOffsets > 0 {
            this.READ_UVLC(&offsetLenMinus1, "offset_len_minus1")
        }
        entryPointOffset = make([]uint, numEntryPointOffsets)
        for idx := uint(0); idx < numEntryPointOffsets; idx++ {
            this.READ_CODE(offsetLenMinus1+1, &uiCode, "entry_point_offset_minus1")
            entryPointOffset[idx] = uiCode + 1
        }
    }else{
        rpcSlice.SetNumEntryPointOffsets ( 0 );
    }

    if pps.GetSliceHeaderExtensionPresentFlag(){
        this.READ_UVLC(&uiCode,"slice_header_extension_length");
        for i:=uint(0); i<uiCode; i++ {
            var ignore uint;
            this.READ_CODE(8,&ignore,"slice_header_extension_data_byte");
        }
    }
    this.m_pcBitstream.ReadByteAlignment();

    if pps.GetTilesEnabledFlag() || pps.GetEntropyCodingSyncEnabledFlag() {
        endOfSliceHeaderLocation := uint(this.m_pcBitstream.GetByteLocation());
        curEntryPointOffset      := uint(0);
        prevEntryPointOffset     := uint(0);
        for idx:=uint(0); idx<numEntryPointOffsets; idx++{
            curEntryPointOffset += entryPointOffset[ idx ];

            emulationPreventionByteCount := uint(0);
            for curByteIdx := uint(0); curByteIdx<this.m_pcBitstream.NumEmulationPreventionBytesRead(); curByteIdx++ {
               if this.m_pcBitstream.GetEmulationPreventionByteLocationIdx( curByteIdx ) >= ( prevEntryPointOffset + endOfSliceHeaderLocation ) &&
                  this.m_pcBitstream.GetEmulationPreventionByteLocationIdx( curByteIdx ) <  ( curEntryPointOffset  + endOfSliceHeaderLocation )  {
                  emulationPreventionByteCount++;
               }
            }

            entryPointOffset[ idx ] -= emulationPreventionByteCount;
            prevEntryPointOffset = curEntryPointOffset;
        }

        if pps.GetTilesEnabledFlag() {
            //rpcSlice.SetTileLocationCount( numEntryPointOffsets );
            prevPos := uint(0)
            for idx := uint(0); idx < numEntryPointOffsets; /*rpcSlice.GetTileLocationCount()*/ idx++ {
                rpcSlice.SetTileLocation(int(idx), prevPos+entryPointOffset[idx])
                prevPos += entryPointOffset[idx]
            }
        } else if pps.GetEntropyCodingSyncEnabledFlag() {
            numSubstreams := rpcSlice.GetNumEntryPointOffsets() + 1
            rpcSlice.AllocSubstreamSizes(uint(numSubstreams))
            pSubstreamSizes := rpcSlice.GetSubstreamSizes()
            for idx := 0; idx < numSubstreams-1; idx++ {
                if idx < int(numEntryPointOffsets) {
                    pSubstreamSizes[idx] = (entryPointOffset[idx] << 3)
                } else {
                    pSubstreamSizes[idx] = 0
                }
            }
        }
    }

    return firstSliceSegmentInPic==1;
}
func (this *TDecCavlc) ParseTerminatingBit(ruiBit *uint) {
}
func (this *TDecCavlc) ParseMVPIdx(riMVPIdx *int) {
}
func (this *TDecCavlc) ParseSkipFlag(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint) {
}
func (this *TDecCavlc) ParseCUTransquantBypassFlag(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint) {
}
func (this *TDecCavlc) ParseMergeFlag(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth, uiPUIdx uint) {
}
func (this *TDecCavlc) ParseMergeIndex(pcCU *TLibCommon.TComDataCU, ruiMergeIndex *uint) {
}
func (this *TDecCavlc) ParseSplitFlag(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint) {
}
func (this *TDecCavlc) ParsePartSize(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint) {
}
func (this *TDecCavlc) ParsePredMode(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint) {
}
func (this *TDecCavlc) ParseIntraDirLumaAng(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint) {
}
func (this *TDecCavlc) ParseIntraDirChroma(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint) {
}
func (this *TDecCavlc) ParseInterDir(pcCU *TLibCommon.TComDataCU, ruiInterDir *uint, uiAbsPartIdx uint) {
}
func (this *TDecCavlc) ParseRefFrmIdx(pcCU *TLibCommon.TComDataCU, riRefFrmIdx *int, eRefList TLibCommon.RefPicList) {
}
func (this *TDecCavlc) ParseMvd(pcCU *TLibCommon.TComDataCU, uiAbsPartAddr, uiPartIdx, uiDepth uint, eRefList TLibCommon.RefPicList) {
}
func (this *TDecCavlc) ParseDeltaQP(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint) {
  var qp, iDQp int;

  this.xReadSvlc( &iDQp );

  qpBdOffsetY := pcCU.GetSlice().GetSPS().GetQpBDOffsetY();
  qp = (( int(pcCU.GetRefQP( uiAbsPartIdx )) + iDQp + 52 + 2*qpBdOffsetY )%(52+ qpBdOffsetY)) -  qpBdOffsetY;

  uiAbsQpCUPartIdx := (uiAbsPartIdx>>((pcCU.GetSlice().GetSPS().GetMaxCUDepth() - pcCU.GetSlice().GetPPS().GetMaxCuDQPDepth())<<1))<<((pcCU.GetSlice().GetSPS().GetMaxCUDepth() - pcCU.GetSlice().GetPPS().GetMaxCuDQPDepth())<<1) ;

  uiQpCUDepth := TLibCommon.MIN(uint(uiDepth), uint(pcCU.GetSlice().GetPPS().GetMaxCuDQPDepth())).(uint) ;

  pcCU.SetQPSubParts( qp, uiAbsQpCUPartIdx, uiQpCUDepth );
}
func (this *TDecCavlc) ParseCoeffNxN(pcCU *TLibCommon.TComDataCU, pcCoef []TLibCommon.TCoeff, uiAbsPartIdx, uiWidth, uiHeight, uiDepth uint, eTType TLibCommon.TextType) {
}
func (this *TDecCavlc) ParseTransformSkipFlags(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, width, height, uiDepth uint, eTType TLibCommon.TextType) {
}
func (this *TDecCavlc) ParseIPCMInfo(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint) {
}
func (this *TDecCavlc) UpdateContextTables(eSliceType TLibCommon.SliceType, iQp int) {
}

func (this *TDecCavlc) xParsePredWeightTable(pcSlice *TLibCommon.TComSlice) {
    var wp []TLibCommon.WpScalingParam
    bChroma := true // color always present in HEVC ?

    eSliceType := pcSlice.GetSliceType()
    iNbRef := int(1 + TLibCommon.B2U(eSliceType == TLibCommon.B_SLICE))
    var uiLog2WeightDenomLuma, uiLog2WeightDenomChroma uint

    uiTotalSignalledWeightFlags := uint(0)

    var iDeltaDenom int
    // decode delta_luma_log2_weight_denom :
    this.READ_UVLC(&uiLog2WeightDenomLuma, "luma_log2_weight_denom") // ue(v): luma_log2_weight_denom
    if bChroma {
        this.READ_SVLC(&iDeltaDenom, "delta_chroma_log2_weight_denom") // se(v): delta_chroma_log2_weight_denom
        //assert((iDeltaDenom + (Int)uiLog2WeightDenomLuma)>=0);
        uiLog2WeightDenomChroma = uint(iDeltaDenom + int(uiLog2WeightDenomLuma))
    }

    for iNumRef := 0; iNumRef < iNbRef; iNumRef++ {
        var eRefPicList TLibCommon.RefPicList
        if iNumRef != 0 {
            eRefPicList = TLibCommon.REF_PIC_LIST_1
        } else {
            eRefPicList = TLibCommon.REF_PIC_LIST_0
        }
        for iRefIdx := 0; iRefIdx < pcSlice.GetNumRefIdx(eRefPicList); iRefIdx++ {
            wp = pcSlice.GetWpScaling(eRefPicList, iRefIdx)

            wp[0].SetLog2WeightDenom(uiLog2WeightDenomLuma)
            wp[1].SetLog2WeightDenom(uiLog2WeightDenomChroma)
            wp[2].SetLog2WeightDenom(uiLog2WeightDenomChroma)

            var uiCode uint
            this.READ_FLAG(&uiCode, "luma_weight_lX_flag") // u(1): luma_weight_l0_flag
            wp[0].SetPresentFlag(uiCode == 1)
            uiTotalSignalledWeightFlags += uint(TLibCommon.B2U(wp[0].GetPresentFlag()))
        }
        if bChroma {
            var uiCode uint
            for iRefIdx := 0; iRefIdx < pcSlice.GetNumRefIdx(eRefPicList); iRefIdx++ {
                wp = pcSlice.GetWpScaling(eRefPicList, iRefIdx)
                this.READ_FLAG(&uiCode, "chroma_weight_lX_flag") // u(1): chroma_weight_l0_flag
                wp[1].SetPresentFlag(uiCode == 1)
                wp[2].SetPresentFlag(uiCode == 1)
                uiTotalSignalledWeightFlags += 2 * uint(TLibCommon.B2U(wp[1].GetPresentFlag()))
            }
        }
        for iRefIdx := 0; iRefIdx < pcSlice.GetNumRefIdx(eRefPicList); iRefIdx++ {
            wp = pcSlice.GetWpScaling(eRefPicList, iRefIdx)
            if wp[0].GetPresentFlag() {
                var iDeltaWeight, iOffset int
                this.READ_SVLC(&iDeltaWeight, "delta_luma_weight_lX") // se(v): delta_luma_weight_l0[i]
                wp[0].SetWeight(iDeltaWeight + (1 << wp[0].GetLog2WeightDenom()))
                this.READ_SVLC(&iOffset, "luma_offset_lX") // se(v): luma_offset_l0[i]
                wp[0].SetOffset(iOffset)
            } else {
                wp[0].SetWeight(1 << wp[0].GetLog2WeightDenom())
                wp[0].SetOffset(0)
            }
            if bChroma {
                if wp[1].GetPresentFlag() {
                    for j := 1; j < 3; j++ {
                        var iDeltaWeight int
                        this.READ_SVLC(&iDeltaWeight, "delta_chroma_weight_lX") // se(v): chroma_weight_l0[i][j]
                        wp[j].SetWeight(iDeltaWeight + (1 << wp[1].GetLog2WeightDenom()))

                        var iDeltaChroma int
                        this.READ_SVLC(&iDeltaChroma, "delta_chroma_offset_lX") // se(v): delta_chroma_offset_l0[i][j]
                        pred := (128 - ((128 * wp[j].GetWeight()) >> (wp[j].GetLog2WeightDenom())))
                        wp[j].SetOffset(TLibCommon.CLIP3(-128, 127, (iDeltaChroma + pred)).(int))
                    }
                } else {
                    for j := 1; j < 3; j++ {
                        wp[j].SetWeight(1 << wp[j].GetLog2WeightDenom())
                        wp[j].SetOffset(0)
                    }
                }
            }
        }

        for iRefIdx := pcSlice.GetNumRefIdx(eRefPicList); iRefIdx < TLibCommon.MAX_NUM_REF; iRefIdx++ {
            wp = pcSlice.GetWpScaling(eRefPicList, iRefIdx)

            wp[0].SetPresentFlag(false)
            wp[1].SetPresentFlag(false)
            wp[2].SetPresentFlag(false)
        }
    }
}
func (this *TDecCavlc) ParseScalingList(scalingList *TLibCommon.TComScalingList) {
    var code, sizeId, listId uint
    var scalingListPredModeFlag bool
    //for each size
    for sizeId = 0; sizeId < TLibCommon.SCALING_LIST_SIZE_NUM; sizeId++ {
        for listId = 0; listId < TLibCommon.G_scalingListNum[sizeId]; listId++ {
            this.READ_FLAG(&code, "scaling_list_pred_mode_flag")
            scalingListPredModeFlag = code != 0
            if !scalingListPredModeFlag { //Copy Mode
                this.READ_UVLC(&code, "scaling_list_pred_matrix_id_delta")
                scalingList.SetRefMatrixId(sizeId, listId, listId-code)
                if sizeId > TLibCommon.SCALING_LIST_8x8 {
                    if listId == scalingList.GetRefMatrixId(sizeId, listId) {
                        scalingList.SetScalingListDC(sizeId, listId, 16)
                    } else {
                        scalingList.SetScalingListDC(sizeId, listId, uint(scalingList.GetScalingListDC(sizeId, scalingList.GetRefMatrixId(sizeId, listId))))
                    }
                }
                scalingList.ProcessRefMatrix(sizeId, listId, scalingList.GetRefMatrixId(sizeId, listId))
            } else { //DPCM Mode
                this.xDecodeScalingList(scalingList, sizeId, listId)
            }
        }
    }
}
func (this *TDecCavlc) xDecodeScalingList(scalingList *TLibCommon.TComScalingList, sizeId, listId uint) {
    var i, coefNum int
    coefNum = TLibCommon.MIN(int(TLibCommon.MAX_MATRIX_COEF_NUM), int(TLibCommon.G_scalingListSize[sizeId])).(int)
    var data int
    scalingListDcCoefMinus8 := int(0)
    nextCoef := int(TLibCommon.SCALING_LIST_START_VALUE)
    var scan []uint
    if sizeId == 0 {
        scan = TLibCommon.G_auiSigLastScan[TLibCommon.SCAN_DIAG][1][:]
    } else {
        scan = TLibCommon.G_sigLastScanCG32x32[:]
    }
    dst := scalingList.GetScalingListAddress(sizeId, listId)

    if sizeId > TLibCommon.SCALING_LIST_8x8 {
        this.READ_SVLC(&scalingListDcCoefMinus8, "scaling_list_dc_coef_minus8")
        scalingList.SetScalingListDC(sizeId, listId, uint(scalingListDcCoefMinus8+8))
        nextCoef = scalingList.GetScalingListDC(sizeId, listId)
    }

    for i = 0; i < coefNum; i++ {
        this.READ_SVLC(&data, "scaling_list_delta_coef")
        nextCoef = (nextCoef + data + 256) % 256
        dst[scan[i]] = nextCoef
    }
}

//protected:
func (this *TDecCavlc) xMoreRbspData() bool {
    bitsLeft := this.m_pcBitstream.GetNumBitsLeft()

    // if there are more than 8 bits, it cannot be rbsp_trailing_bits
    if bitsLeft > 8 {
        return true
    }

    lastByte := this.m_pcBitstream.PeekBits(bitsLeft)
    cnt := bitsLeft

    // remove trailing bits equal to zero
    for (cnt > 0) && ((lastByte & 1) == 0) {
        lastByte >>= 1
        cnt--
    }
    // remove bit equal to one
    cnt--

    // we should not have a negative number of bits
    //assert (cnt>=0);

    // we have more data, if cnt is not zero
    return cnt > 0
}
