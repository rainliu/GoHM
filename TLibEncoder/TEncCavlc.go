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

/*#if ENC_DEC_TRACE

#define WRITE_CODE( value, length, name)    xWriteCodeTr ( value, length, name )
#define WRITE_UVLC( value,         name)    xWriteUvlcTr ( value,         name )
#define WRITE_SVLC( value,         name)    xWriteSvlcTr ( value,         name )
#define WRITE_FLAG( value,         name)    xWriteFlagTr ( value,         name )

#else

#define WRITE_CODE( value, length, name)     xWriteCode ( value, length )
#define WRITE_UVLC( value,         name)     xWriteUvlc ( value )
#define WRITE_SVLC( value,         name)     xWriteSvlc ( value )
#define WRITE_FLAG( value,         name)     xWriteFlag ( value )

#endif*/

type SyntaxElementWriter struct{
  m_pTraceFile  io.Writer
  m_pcBitIf TLibCommon.TComBitIf;
}

func NewSyntaxElementWriter() *SyntaxElementWriter{
	return &SyntaxElementWriter{};
}

func (this *SyntaxElementWriter) setBitstream          ( p TLibCommon.TComBitIf)  { this.m_pcBitIf = p;  }


func (this *SyntaxElementWriter) SetTraceFile(traceFile io.Writer) {
    this.m_pTraceFile = traceFile
}

func (this *SyntaxElementWriter) GetTraceFile() io.Writer {
    return this.m_pTraceFile
}


func (this *SyntaxElementWriter) xTraceVUIHeader(pVUI *TLibCommon.TComVUI) {
    if this.GetTraceFile() != nil {
        io.WriteString(this.m_pTraceFile, fmt.Sprintf("========= VUI Parameter Set ===============================================\n")) //, pVPS.GetVPSId() );
    }
}

func (this *SyntaxElementWriter) xTraceVPSHeader(pVPS *TLibCommon.TComVPS) {
    if this.GetTraceFile() != nil {
        io.WriteString(this.m_pTraceFile, fmt.Sprintf("========= Video Parameter Set =============================================\n")) //, pVPS.GetVPSId() );
    }
}

func (this *SyntaxElementWriter) xTraceSPSHeader(pSPS *TLibCommon.TComSPS) {
    if this.GetTraceFile() != nil {
        io.WriteString(this.m_pTraceFile, fmt.Sprintf("========= Sequence Parameter Set ==========================================\n")) //, pSPS.GetSPSId() );
    }
}

func (this *SyntaxElementWriter) xTracePPSHeader(pPPS *TLibCommon.TComPPS) {
    if this.GetTraceFile() != nil {
        io.WriteString(this.m_pTraceFile, fmt.Sprintf("========= Picture Parameter Set ===========================================\n")) //, pPPS.GetPPSId() );
    }
}

func (this *SyntaxElementWriter) xTraceSliceHeader(pSlice *TLibCommon.TComSlice) {
    if this.GetTraceFile() != nil {
        io.WriteString(this.m_pTraceFile, fmt.Sprintf("========= Slice Parameter Set =============================================\n"))
    }
}

func (this *SyntaxElementWriter) WRITE_CODE ( value, length uint, pSymbolName string) {
  this.xWriteCode (value,length);
  if this.GetTraceFile() != nil {
    //io.WriteString(this.m_pTraceFile, fmt.Sprintf("%8lld  ", g_nSymbolCounter++ );
    if length<10 {
      io.WriteString(this.m_pTraceFile, fmt.Sprintf("%-50s u(%d)  : %d\n", pSymbolName, length, value ));
    }else{
      io.WriteString(this.m_pTraceFile, fmt.Sprintf("%-50s u(%d) : %d\n", pSymbolName, length, value ));
    }
  }
}

func (this *SyntaxElementWriter) WRITE_UVLC ( value uint, pSymbolName string) {
  this.xWriteUvlc (value);
  if this.GetTraceFile() != nil {
    //io.WriteString(this.m_pTraceFile, fmt.Sprintf("%8lld  ", g_nSymbolCounter++ );
    io.WriteString(this.m_pTraceFile, fmt.Sprintf("%-50s ue(v) : %d\n", pSymbolName, value ));
  }
}

func (this *SyntaxElementWriter) WRITE_SVLC ( value int, pSymbolName string) {
  this.xWriteSvlc(value);
  if this.GetTraceFile() != nil {
    //fprintf( g_hTrace, "%8lld  ", g_nSymbolCounter++ );
    io.WriteString(this.m_pTraceFile, fmt.Sprintf("%-50s se(v) : %d\n", pSymbolName, value ));
  }
}

func (this *SyntaxElementWriter) WRITE_FLAG( value uint, pSymbolName string) {
  this.xWriteFlag(value);
  if this.GetTraceFile() != nil {
    //fprintf( g_hTrace, "%8lld  ", g_nSymbolCounter++ );
    io.WriteString(this.m_pTraceFile, fmt.Sprintf("%-50s u(1)  : %d\n", pSymbolName, value ));
  }
}

func (this *SyntaxElementWriter)  xWriteCode            ( uiCode, uiLength uint)   { this.m_pcBitIf.Write( uiCode, uiLength );}
func (this *SyntaxElementWriter)  xWriteUvlc            ( uiCode uint){
  uiLength := uint(1);
  var uiTemp uint;
  uiCode++;

  uiTemp = uiCode;
  //assert ( uiTemp );

  for 1 != uiTemp {
    uiTemp >>= 1;
    uiLength += 2;
  }
  // Take care of cases where uiLength > 32
  this.m_pcBitIf.Write( 0, uiLength >> 1);
  this.m_pcBitIf.Write( uiCode, (uiLength+1) >> 1);
}
func (this *SyntaxElementWriter)  xWriteSvlc            (  iCode int ){
  var uiCode uint;

  uiCode = this.xConvertToUInt( iCode );
  this.xWriteUvlc( uiCode );
}
func (this *SyntaxElementWriter)  xWriteFlag            ( uiCode uint){  this.m_pcBitIf.Write( uiCode, 1 );}
/*#if ENC_DEC_TRACE
  Void  xWriteCodeTr          ( UInt value, UInt  length, const Char *pSymbolName);
  Void  xWriteUvlcTr          ( UInt value,               const Char *pSymbolName);
  Void  xWriteSvlcTr          ( Int  value,               const Char *pSymbolName);
  Void  xWriteFlagTr          ( UInt value,               const Char *pSymbolName);
#endif*/

func (this *SyntaxElementWriter)  xConvertToUInt        ( iValue int ) uint {
	if iValue <= 0 {
		return  uint(-iValue)<<1;
	}

	return uint(iValue<<1)-1;
}


// ====================================================================================================================
// Class definition
// ====================================================================================================================

/// CAVLC encoder class
type TEncCavlc struct{
  SyntaxElementWriter

  m_pcSlice 	*TLibCommon.TComSlice;
  m_uiCoeffCost	uint;
}

func NewTEncCavlc() *TEncCavlc{
	return &TEncCavlc{};
}


func (this *TEncCavlc)  xWritePCMAlignZero    () {this.m_pcBitIf.WriteAlignZero();}
func (this *TEncCavlc)  xWriteEpExGolomb      ( uiSymbol, uiCount uint){
  for uiSymbol >= (1<<uiCount) {
    this.xWriteFlag( 1 );
    uiSymbol -= 1<<uiCount;
    uiCount  ++;
  }
  this.xWriteFlag( 0 );
  for uiCount!=0 {
    uiCount--
    this.xWriteFlag( (uiSymbol>>uiCount) & 1 );
  }
  return;
}
func (this *TEncCavlc)  xWriteExGolombLevel   ( uiSymbol uint){
  if uiSymbol!=0 {
    this.xWriteFlag( 1 );
    uiCount := uint(0);
    bNoExGo := (uiSymbol < 13);

    uiSymbol--;
    uiCount++;
    for uiSymbol!=0 && uiCount < 13 {
      this.xWriteFlag( 1 );

      uiSymbol--;
      uiCount++;
    }
    if bNoExGo {
      this.xWriteFlag( 0 );
    }else{
      this.xWriteEpExGolomb( uiSymbol, 0 );
    }
  }else{
    this.xWriteFlag( 0 );
  }
  return;
}
func (this *TEncCavlc)  xWriteUnaryMaxSymbol  ( uiSymbol, uiMaxSymbol uint){
  if uiMaxSymbol == 0 {
    return;
  }
  this.xWriteFlag( uint(TLibCommon.B2U(uiSymbol!=0)));
  if uiSymbol == 0 {
    return;
  }

  bCodeLast := ( uiMaxSymbol > uiSymbol );

  uiSymbol--;
  for uiSymbol!=0 {
    this.xWriteFlag( 1 );
  	uiSymbol--;
  }
  if bCodeLast {
    this.xWriteFlag( 0 );
  }
  return;
}
//#if SPS_INTER_REF_SET_PRED
func (this *TEncCavlc)  codeShortTermRefPicSet              ( pcSPS *TLibCommon.TComSPS, rps *TLibCommon.TComReferencePictureSet, calledFromSliceHeader bool, idx int){
//#else
//func (this *TEncCavlc) codeShortTermRefPicSet              ( TComSPS* pcSPS, TComReferencePictureSet* pcRPS, Bool calledFromSliceHeader );
//#endif
//#if PRINT_RPS_INFO
//  Int lastBits = getNumberOfWrittenBits();
//#endif
//#if SPS_INTER_REF_SET_PRED
  if idx > 0 {
//#endif
  	this.WRITE_FLAG( uint(TLibCommon.B2U(rps.GetInterRPSPrediction())), "inter_ref_pic_set_prediction_flag" ); // inter_RPS_prediction_flag
//#if SPS_INTER_REF_SET_PRED
  }
//#endif
  if rps.GetInterRPSPrediction() {
    deltaRPS := rps.GetDeltaRPS();
    if calledFromSliceHeader {
      this.WRITE_UVLC( uint(rps.GetDeltaRIdxMinus1()), "delta_idx_minus1" ); // delta index of the Reference Picture Set used for prediction minus 1
    }

	if deltaRPS >=0 {
    	this.WRITE_CODE( 0, 1, "delta_rps_sign" ); //delta_rps_sign
    }else{
    	this.WRITE_CODE( 1, 1, "delta_rps_sign" ); //delta_rps_sign
    }
    this.WRITE_UVLC( uint(TLibCommon.ABS(deltaRPS).(int) - 1), "abs_delta_rps_minus1"); // absolute delta RPS minus 1

    for j:=0; j < rps.GetNumRefIdc(); j++ {
      refIdc := rps.GetRefIdc(j);
      this.WRITE_CODE( uint(TLibCommon.B2U(refIdc==1)), 1, "used_by_curr_pic_flag" ); //first bit is "1" if Idc is 1
      if refIdc != 1 {
        this.WRITE_CODE(uint(refIdc>>1), 1, "use_delta_flag" ); //second bit is "1" if Idc is 2, "0" otherwise.
      }
    }
  }else{
    this.WRITE_UVLC( uint(rps.GetNumberOfNegativePictures()), "num_negative_pics" );
    this.WRITE_UVLC( uint(rps.GetNumberOfPositivePictures()), "num_positive_pics" );
    prev := 0;
    for j:=0 ; j < rps.GetNumberOfNegativePictures(); j++ {
      this.WRITE_UVLC( uint(prev-rps.GetDeltaPOC(j)-1), "delta_poc_s0_minus1" );
      prev = rps.GetDeltaPOC(j);
      this.WRITE_FLAG( uint(TLibCommon.B2U(rps.GetUsed(j))), "used_by_curr_pic_s0_flag");
    }
    prev = 0;
    for j:=rps.GetNumberOfNegativePictures(); j < rps.GetNumberOfNegativePictures()+rps.GetNumberOfPositivePictures(); j++ {
      this.WRITE_UVLC( uint(rps.GetDeltaPOC(j)-prev-1), "delta_poc_s1_minus1" );
      prev = rps.GetDeltaPOC(j);
      this.WRITE_FLAG( uint(TLibCommon.B2U(rps.GetUsed(j))), "used_by_curr_pic_s1_flag" );
    }
  }

//#if PRINT_RPS_INFO
//  printf("irps=%d (%2d bits) ", rps.GetInterRPSPrediction(), getNumberOfWrittenBits() - lastBits);
//  rps->printDeltaPOC();
//#endif
}
func (this *TEncCavlc)  findMatchingLTRP ( pcSlice *TLibCommon.TComSlice, ltrpsIndex *uint, ltrpPOC int, usedFlag bool ) bool{
  // Bool state = true, state2 = false;
  lsb := ltrpPOC % (1<<pcSlice.GetSPS().GetBitsForPOC());
  for k := uint(0); k < pcSlice.GetSPS().GetNumLongTermRefPicSPS(); k++ {
    if (lsb == int(pcSlice.GetSPS().GetLtRefPicPocLsbSps(uint(k)))) && (usedFlag == pcSlice.GetSPS().GetUsedByCurrPicLtSPSFlag(int(k))) {
      *ltrpsIndex = k;
      return true;
    }
  }
  return false;
}
func (this *TEncCavlc)  resetEntropy          () {}
func (this *TEncCavlc)  determineCabacInitIdx () {}
func (this *TEncCavlc)  setBitstream          ( p TLibCommon.TComBitIf)  { this.m_pcBitIf = p;  }
func (this *TEncCavlc)  setSlice              ( p *TLibCommon.TComSlice)  { this.m_pcSlice = p;  }
func (this *TEncCavlc)  resetBits             ()                { this.m_pcBitIf.ResetBits(); }
func (this *TEncCavlc)  resetCoeffCost        ()                { this.m_uiCoeffCost = 0;  }
func (this *TEncCavlc)  getNumberOfWrittenBits() uint               { return  this.m_pcBitIf.GetNumberOfWrittenBits();  }
func (this *TEncCavlc)  getCoeffCost          () uint               { return  this.m_uiCoeffCost;  }
func (this *TEncCavlc)  codeVPS                 ( pcVPS *TLibCommon.TComVPS){
  this.WRITE_CODE( uint(pcVPS.GetVPSId()),                    4,        "video_parameter_set_id" );
  this.WRITE_FLAG( uint(TLibCommon.B2U(pcVPS.GetTemporalNestingFlag())),                "vps_temporal_id_nesting_flag" );
//#if VPS_REARRANGE
  this.WRITE_CODE( 3,                                    2,        "vps_reserved_three_2bits" );
//#else
//  WRITE_CODE( 0,                                    2,        "vps_reserved_zero_2bits" );
//#endif
  this.WRITE_CODE( 0,                                    6,        "vps_reserved_zero_6bits" );
  this.WRITE_CODE( pcVPS.GetMaxTLayers() - 1,           3,        "vps_max_sub_layers_minus1" );
//#if VPS_REARRANGE
  this.WRITE_CODE( 0xffff,                              16,        "vps_reserved_ffff_16bits" );
  this.codePTL( pcVPS.GetPTL(), true, int(pcVPS.GetMaxTLayers()) - 1 );
//#else
//  codePTL( pcVPS.GetPTL(), true, pcVPS.GetMaxTLayers() - 1 );
//  WRITE_CODE( 0,                                   12,        "vps_reserved_zero_12bits" );
//#endif
//#if SIGNAL_BITRATE_PICRATE_IN_VPS
  this.codeBitratePicRateInfo(pcVPS.GetBitratePicrateInfo(), 0, int(pcVPS.GetMaxTLayers()) - 1);
//#endif
//#if HLS_ADD_SUBLAYER_ORDERING_INFO_PRESENT_FLAG
  subLayerOrderingInfoPresentFlag := uint(1);
  this.WRITE_FLAG(subLayerOrderingInfoPresentFlag,              "vps_sub_layer_ordering_info_present_flag");
//#endif // HLS_ADD_SUBLAYER_ORDERING_INFO_PRESENT_FLAG
  for i:=uint(0); i <= pcVPS.GetMaxTLayers()-1; i++ {
    this.WRITE_UVLC( pcVPS.GetMaxDecPicBuffering(i),           "vps_max_dec_pic_buffering[i]" );
    this.WRITE_UVLC( pcVPS.GetNumReorderPics(i),               "vps_num_reorder_pics[i]" );
    this.WRITE_UVLC( pcVPS.GetMaxLatencyIncrease(i),           "vps_max_latency_increase[i]" );
//#if HLS_ADD_SUBLAYER_ORDERING_INFO_PRESENT_FLAG
    if subLayerOrderingInfoPresentFlag==0 {
      break;
    }
//#endif // HLS_ADD_SUBLAYER_ORDERING_INFO_PRESENT_FLAG
  }

//#if VPS_OPERATING_POINT
  //assert( pcVPS.GetNumHrdParameters() <= MAX_VPS_NUM_HRD_PARAMETERS );
  //assert( pcVPS.GetMaxNuhReservedZeroLayerId() < MAX_VPS_NUH_RESERVED_ZERO_LAYER_ID_PLUS1 );
  this.WRITE_UVLC( pcVPS.GetNumHrdParameters(),                 "vps_num_hrd_parameters" );
  this.WRITE_CODE( pcVPS.GetMaxNuhReservedZeroLayerId(), 6,     "vps_max_nuh_reserved_zero_layer_id" );
  for opIdx := uint(0); opIdx < pcVPS.GetNumHrdParameters(); opIdx++ {
    if opIdx > 0 {
      // operation_point_layer_id_flag( opIdx )
      for i := uint(0); i <= pcVPS.GetMaxNuhReservedZeroLayerId(); i++ {
        //this.WRITE_FLAG( uint(TLibCommon.B2U(pcVPS.GetOpLayerIdIncludedFlag( opIdx, i ))), "op_layer_id_included_flag[opIdx][i]" );
      }
    }
    // TODO: add hrd_parameters()
  }
//#else
//  WRITE_UVLC( 0,                                           "vps_num_hrd_parameters" );
  // hrd_parameters
//#endif
  this.WRITE_FLAG( 0,                     "vps_extension_flag" );

  //future extensions here..

  return;
}

func (this *TEncCavlc)  codeVUI                 ( pcVUI *TLibCommon.TComVUI, pcSPS *TLibCommon.TComSPS){
//#if ENC_DEC_TRACE
  //fprintf( g_hTrace, "----------- vui_parameters -----------\n");
  this.xTraceVUIHeader (pcVUI);
//#endif
  this.WRITE_FLAG(uint(TLibCommon.B2U(pcVUI.GetAspectRatioInfoPresentFlag())),            "aspect_ratio_info_present_flag");
  if pcVUI.GetAspectRatioInfoPresentFlag() {
    this.WRITE_CODE(uint(pcVUI.GetAspectRatioIdc()), 8,                   "aspect_ratio_idc" );
    if pcVUI.GetAspectRatioIdc() == 255 {
      this.WRITE_CODE(uint(pcVUI.GetSarWidth()), 16,                      "sar_width");
      this.WRITE_CODE(uint(pcVUI.GetSarHeight()), 16,                     "sar_height");
    }
  }
  this.WRITE_FLAG(uint(TLibCommon.B2U(pcVUI.GetOverscanInfoPresentFlag())),               "overscan_info_present_flag");
  if pcVUI.GetOverscanInfoPresentFlag() {
    this.WRITE_FLAG(uint(TLibCommon.B2U(pcVUI.GetOverscanAppropriateFlag())),             "overscan_appropriate_flag");
  }
  this.WRITE_FLAG(uint(TLibCommon.B2U(pcVUI.GetVideoSignalTypePresentFlag())),            "video_signal_type_present_flag");
  if pcVUI.GetVideoSignalTypePresentFlag() {
    this.WRITE_CODE(uint(pcVUI.GetVideoFormat()), 3,                      "video_format");
    this.WRITE_FLAG(uint(TLibCommon.B2U(pcVUI.GetVideoFullRangeFlag())),                  "video_full_range_flag");
    this.WRITE_FLAG(uint(TLibCommon.B2U(pcVUI.GetColourDescriptionPresentFlag())),        "colour_description_present_flag");
    if pcVUI.GetColourDescriptionPresentFlag() {
      this.WRITE_CODE(uint(pcVUI.GetColourPrimaries()), 8,                "colour_primaries");
      this.WRITE_CODE(uint(pcVUI.GetTransferCharacteristics()), 8,        "transfer_characteristics");
      this.WRITE_CODE(uint(pcVUI.GetMatrixCoefficients()), 8,             "matrix_coefficients");
    }
  }

  this.WRITE_FLAG(uint(TLibCommon.B2U(pcVUI.GetChromaLocInfoPresentFlag())),              "chroma_loc_info_present_flag");
  if pcVUI.GetChromaLocInfoPresentFlag() {
    this.WRITE_UVLC(uint(pcVUI.GetChromaSampleLocTypeTopField()),         "chroma_sample_loc_type_top_field");
    this.WRITE_UVLC(uint(pcVUI.GetChromaSampleLocTypeBottomField()),      "chroma_sample_loc_type_bottom_field");
  }

  this.WRITE_FLAG(uint(TLibCommon.B2U(pcVUI.GetNeutralChromaIndicationFlag())),           "neutral_chroma_indication_flag");
  this.WRITE_FLAG(uint(TLibCommon.B2U(pcVUI.GetFieldSeqFlag())),                          "field_seq_flag");
  //assert(pcVUI.GetFieldSeqFlag() == 0);                        // not currently supported
//#if HLS_DISPLAY_WINDOW_PLACEHOLDER
  this.WRITE_FLAG(0,                                                 "default_display_window_flag");
//#endif
//#if HLS_ADD_VUI_PICSTRUCT_PRESENT_FLAG
  //this.WRITE_FLAG(uint(TLibCommon.B2U(pcVUI.GetPicStructPresentFlag())),                  "pic_struct_present_flag");
//#endif /* HLS_ADD_VUI_PICSTRUCT_PRESENT_FLAG */
/*  this.WRITE_FLAG(uint(TLibCommon.B2U(pcVUI.GetHrdParametersPresentFlag())),              "hrd_parameters_present_flag");
  if pcVUI.GetHrdParametersPresentFlag() {
    this.WRITE_FLAG(uint(TLibCommon.B2U(pcVUI.GetTimingInfoPresentFlag())),               "timing_info_present_flag");
    if pcVUI.GetTimingInfoPresentFlag() {
      this.WRITE_CODE(pcVUI.GetNumUnitsInTick(), 32,                "num_units_in_tick");
      this.WRITE_CODE(pcVUI.GetTimeScale(),      32,                "time_scale");
    }
    this.WRITE_FLAG(uint(TLibCommon.B2U(pcVUI.GetNalHrdParametersPresentFlag())),         "nal_hrd_parameters_present_flag");
    this.WRITE_FLAG(uint(TLibCommon.B2U(pcVUI.GetVclHrdParametersPresentFlag())),         "vcl_hrd_parameters_present_flag");
    if pcVUI.GetNalHrdParametersPresentFlag() || pcVUI.GetVclHrdParametersPresentFlag() {
      this.WRITE_FLAG(uint(TLibCommon.B2U(pcVUI.GetSubPicCpbParamsPresentFlag())),        "sub_pic_cpb_params_present_flag");
      if pcVUI.GetSubPicCpbParamsPresentFlag() {
        this.WRITE_CODE(pcVUI.GetTickDivisorMinus2(), 8,            "tick_divisor_minus2");
        this.WRITE_CODE(pcVUI.GetDuCpbRemovalDelayLengthMinus1(), 5,  "du_cpb_removal_delay_length_minus1");
      }
      this.WRITE_CODE(pcVUI.GetBitRateScale(), 4,                   "bit_rate_scale");
      this.WRITE_CODE(pcVUI.GetCpbSizeScale(), 4,                   "cpb_size_scale");
//#if HRD_BUFFER
      if pcVUI.GetSubPicCpbParamsPresentFlag() {
        this.WRITE_CODE(pcVUI.GetDuCpbSizeScale(), 4,           "du_cpb_size_scale");
      }
//#endif
      this.WRITE_CODE(pcVUI.GetInitialCpbRemovalDelayLengthMinus1(), 5, "initial_cpb_removal_delay_length_minus1");
      this.WRITE_CODE(pcVUI.GetCpbRemovalDelayLengthMinus1(),        5, "cpb_removal_delay_length_minus1");
      this.WRITE_CODE(pcVUI.GetDpbOutputDelayLengthMinus1(),         5, "dpb_output_delay_length_minus1");
    }
    var i, j, nalOrVcl int;
    for i = 0; i < int(pcSPS.GetMaxTLayers()); i ++ {
      this.WRITE_FLAG(uint(TLibCommon.B2U(pcVUI.GetFixedPicRateFlag(i))),                 "fixed_pic_rate_flag");
      if pcVUI.GetFixedPicRateFlag( i ) {
        this.WRITE_UVLC(pcVUI.GetPicDurationInTcMinus1(i),          "pic_duration_in_tc_minus1");
      }
      this.WRITE_FLAG(uint(TLibCommon.B2U(pcVUI.GetLowDelayHrdFlag(i))),                  "low_delay_hrd_flag");
      this.WRITE_UVLC(pcVUI.GetCpbCntMinus1(i),                     "cpb_cnt_minus1");

      for nalOrVcl = 0; nalOrVcl < 2; nalOrVcl ++ {
        if  ( ( nalOrVcl == 0 ) && ( pcVUI.GetNalHrdParametersPresentFlag() ) ) ||
            ( ( nalOrVcl == 1 ) && ( pcVUI.GetVclHrdParametersPresentFlag() ) ) {
          for j = 0; j < int( pcVUI.GetCpbCntMinus1( i ) + 1 ); j ++ {
            this.WRITE_UVLC(pcVUI.GetBitRateValueMinus1(i, j, nalOrVcl), "bit_size_value_minus1");
            this.WRITE_UVLC(pcVUI.GetCpbSizeValueMinus1(i, j, nalOrVcl), "cpb_size_value_minus1");
//#if HRD_BUFFER
            if pcVUI.GetSubPicCpbParamsPresentFlag() {
              this.WRITE_UVLC(pcVUI.GetDuCpbSizeValueMinus1(i, j, nalOrVcl), "du_cpb_size_value_minus1");
            }
//#endif

            this.WRITE_FLAG(uint(TLibCommon.B2U(pcVUI.GetCbrFlag(i, j, nalOrVcl))),            "cbr_flag");
          }
        }
      }
    }
  }
//#if POC_TEMPORAL_RELATIONSHIP
  this.WRITE_FLAG( uint(TLibCommon.B2U(pcVUI.GetPocProportionalToTimingFlag())), "poc_proportional_to_timing_flag" );
  if pcVUI.GetPocProportionalToTimingFlag() && pcVUI.GetTimingInfoPresentFlag() {
    this.WRITE_UVLC( uint(pcVUI.GetNumTicksPocDiffOneMinus1()), "num_ticks_poc_diff_one_minus1" );
  }
//#endif
*/
  this.WRITE_FLAG(uint(TLibCommon.B2U(pcVUI.GetBitstreamRestrictionFlag())),              "bitstream_restriction_flag");
  if pcVUI.GetBitstreamRestrictionFlag() {
    this.WRITE_FLAG(uint(TLibCommon.B2U(pcVUI.GetTilesFixedStructureFlag())),             "tiles_fixed_structure_flag");
    this.WRITE_FLAG(uint(TLibCommon.B2U(pcVUI.GetMotionVectorsOverPicBoundariesFlag())),  "motion_vectors_over_pic_boundaries_flag");
//#if HLS_MOVE_SPS_PICLIST_FLAGS
    this.WRITE_FLAG(uint(TLibCommon.B2U(pcVUI.GetRestrictedRefPicListsFlag())),           "restricted_ref_pic_lists_flag");
//#endif /* HLS_MOVE_SPS_PICLIST_FLAGS */
//#if MIN_SPATIAL_SEGMENTATION
    this.WRITE_CODE(uint(pcVUI.GetMinSpatialSegmentationIdc()),        8, "min_spatial_segmentation_idc");
//#endif
    this.WRITE_UVLC(uint(pcVUI.GetMaxBytesPerPicDenom()),                 "max_bytes_per_pic_denom");
    this.WRITE_UVLC(uint(pcVUI.GetMaxBitsPerMinCuDenom()),                "max_bits_per_mincu_denom");
    this.WRITE_UVLC(uint(pcVUI.GetLog2MaxMvLengthHorizontal()),           "log2_max_mv_length_horizontal");
    this.WRITE_UVLC(uint(pcVUI.GetLog2MaxMvLengthVertical()),             "log2_max_mv_length_vertical");
  }
}
func (this *TEncCavlc)  codeSPS                 ( pcSPS *TLibCommon.TComSPS){
//#if ENC_DEC_TRACE
  this.xTraceSPSHeader (pcSPS);
//#endif
  this.WRITE_CODE( uint(pcSPS.GetVPSId ()),          4,       "video_parameter_set_id" );
  this.WRITE_CODE( pcSPS.GetMaxTLayers() - 1,  3,       "sps_max_sub_layers_minus1" );
//#if MOVE_SPS_TEMPORAL_ID_NESTING_FLAG
  this.WRITE_FLAG( uint(TLibCommon.B2U(pcSPS.GetTemporalIdNestingFlag())),                             "sps_temporal_id_nesting_flag" );
/*#else
  this.WRITE_FLAG( 0,                                    "sps_reserved_zero_bit");
#endif*/
  this.codePTL(pcSPS.GetPTL(), true, int(pcSPS.GetMaxTLayers()) - 1);
  this.WRITE_UVLC( uint(pcSPS.GetSPSId ()),                   "seq_parameter_set_id" );
  this.WRITE_UVLC( uint(pcSPS.GetChromaFormatIdc ()),         "chroma_format_idc" );
  //assert(pcSPS.GetChromaFormatIdc () == 1);
  // in the first version chroma_format_idc can only be equal to 1 (4:2:0)
  if pcSPS.GetChromaFormatIdc () == 3 {
    this.WRITE_FLAG( 0,                                  "separate_colour_plane_flag");
  }

  this.WRITE_UVLC( pcSPS.GetPicWidthInLumaSamples (),   "pic_width_in_luma_samples" );
  this.WRITE_UVLC( pcSPS.GetPicHeightInLumaSamples(),   "pic_height_in_luma_samples" );
/*  crop := pcSPS.GetPicCroppingWindow();
  this.WRITE_FLAG( uint(TLibCommon.B2U(crop.GetPicCroppingFlag())),            "pic_cropping_flag" );
  if crop.GetPicCroppingFlag() {
    this.WRITE_UVLC( uint(crop.GetPicCropLeftOffset()   / pcSPS.GetCropUnitX(pcSPS.GetChromaFormatIdc() )), "pic_crop_left_offset" );
    this.WRITE_UVLC( uint(crop.GetPicCropRightOffset()  / pcSPS.GetCropUnitX(pcSPS.GetChromaFormatIdc() )), "pic_crop_right_offset" );
    this.WRITE_UVLC( uint(crop.GetPicCropTopOffset()    / pcSPS.GetCropUnitY(pcSPS.GetChromaFormatIdc() )), "pic_crop_top_offset" );
    this.WRITE_UVLC( uint(crop.GetPicCropBottomOffset() / pcSPS.GetCropUnitY(pcSPS.GetChromaFormatIdc() )), "pic_crop_bottom_offset" );
  }
*/
  this.WRITE_UVLC( uint(pcSPS.GetBitDepthY() - 8),             "bit_depth_luma_minus8" );
  this.WRITE_UVLC( uint(pcSPS.GetBitDepthC() - 8),             "bit_depth_chroma_minus8" );

/*#if !HLS_GROUP_SPS_PCM_FLAGS
  this.WRITE_FLAG( pcSPS.GetUsePCM() ? 1 : 0,                   "pcm_enabled_flag");

  if( pcSPS.GetUsePCM() )
  {
  this.WRITE_CODE( pcSPS.GetPCMBitDepthLuma() - 1, 4,   "pcm_bit_depth_luma_minus1" );
  this.WRITE_CODE( pcSPS.GetPCMBitDepthChroma() - 1, 4, "pcm_bit_depth_chroma_minus1" );
  }

#endif*/ /* !HLS_GROUP_SPS_PCM_FLAGS */
  this.WRITE_UVLC( pcSPS.GetBitsForPOC()-4,                 "log2_max_pic_order_cnt_lsb_minus4" );

//#if HLS_ADD_SUBLAYER_ORDERING_INFO_PRESENT_FLAG
  subLayerOrderingInfoPresentFlag := uint(1);
  this.WRITE_FLAG(subLayerOrderingInfoPresentFlag,       "sps_sub_layer_ordering_info_present_flag");
//#endif /* HLS_ADD_SUBLAYER_ORDERING_INFO_PRESENT_FLAG */
  for i:=uint(0); i <= pcSPS.GetMaxTLayers()-1; i++ {
    this.WRITE_UVLC( pcSPS.GetMaxDecPicBuffering(i),           "max_dec_pic_buffering[i]" );
    this.WRITE_UVLC( uint(pcSPS.GetNumReorderPics(i)),               "num_reorder_pics[i]" );
    this.WRITE_UVLC( pcSPS.GetMaxLatencyIncrease(i),           "max_latency_increase[i]" );
//#if HLS_ADD_SUBLAYER_ORDERING_INFO_PRESENT_FLAG
    if subLayerOrderingInfoPresentFlag==0 {
      break;
    }
//#endif /* HLS_ADD_SUBLAYER_ORDERING_INFO_PRESENT_FLAG */
  }
  //assert( pcSPS.GetMaxCUWidth() == pcSPS.GetMaxCUHeight() );

  MinCUSize := pcSPS.GetMaxCUWidth() >> ( pcSPS.GetMaxCUDepth()-TLibCommon.G_uiAddCUDepth );
  log2MinCUSize := uint(0);
  for MinCUSize > 1 {
    MinCUSize >>= 1;
    log2MinCUSize++;
  }

/*#if !HLS_MOVE_SPS_PICLIST_FLAGS
  this.WRITE_FLAG( pcSPS.GetRestrictedRefPicListsFlag(),                                 "restricted_ref_pic_lists_flag" );
  if( pcSPS.GetRestrictedRefPicListsFlag() )
  {
    this.WRITE_FLAG( pcSPS.GetListsModificationPresentFlag(),                            "lists_modification_present_flag" );
  }
#endif*/ /* !HLS_MOVE_SPS_PICLIST_FLAGS */
  this.WRITE_UVLC( log2MinCUSize - 3,                                                     "log2_min_coding_block_size_minus3" );
  this.WRITE_UVLC( pcSPS.GetMaxCUDepth()-TLibCommon.G_uiAddCUDepth,                                 "log2_diff_max_min_coding_block_size" );
  this.WRITE_UVLC( pcSPS.GetQuadtreeTULog2MinSize() - 2,                                 "log2_min_transform_block_size_minus2" );
  this.WRITE_UVLC( pcSPS.GetQuadtreeTULog2MaxSize() - pcSPS.GetQuadtreeTULog2MinSize(), "log2_diff_max_min_transform_block_size" );
/*#if !HLS_GROUP_SPS_PCM_FLAGS
  if( pcSPS.GetUsePCM() )
  {
    this.WRITE_UVLC( pcSPS.GetPCMLog2MinSize() - 3,                                      "log2_min_pcm_coding_block_size_minus3" );
    this.WRITE_UVLC( pcSPS.GetPCMLog2MaxSize() - pcSPS.GetPCMLog2MinSize(),             "log2_diff_max_min_pcm_coding_block_size" );
  }
#endif*/ /* !HLS_GROUP_SPS_PCM_FLAGS */
  this.WRITE_UVLC( pcSPS.GetQuadtreeTUMaxDepthInter() - 1,                               "max_transform_hierarchy_depth_inter" );
  this.WRITE_UVLC( pcSPS.GetQuadtreeTUMaxDepthIntra() - 1,                               "max_transform_hierarchy_depth_intra" );
  this.WRITE_FLAG( uint(TLibCommon.B2U(pcSPS.GetScalingListFlag())),                                   "scaling_list_enabled_flag" );
  if pcSPS.GetScalingListFlag() {
    this.WRITE_FLAG( uint(TLibCommon.B2U(pcSPS.GetScalingListPresentFlag())),                          "sps_scaling_list_data_present_flag" );
    if pcSPS.GetScalingListPresentFlag() {
/*#if SCALING_LIST_OUTPUT_RESULT
    printf("SPS\n");
#endif*/
      this.codeScalingList( this.m_pcSlice.GetScalingList() );
    }
  }
  this.WRITE_FLAG( uint(TLibCommon.B2U(pcSPS.GetUseAMP())),                                                    "asymmetric_motion_partitions_enabled_flag" );
  this.WRITE_FLAG( uint(TLibCommon.B2U(pcSPS.GetUseSAO())),                                            "sample_adaptive_offset_enabled_flag");

//#if HLS_GROUP_SPS_PCM_FLAGS
  this.WRITE_FLAG( uint(TLibCommon.B2U(pcSPS.GetUsePCM())),                                            "pcm_enabled_flag");
//#endif /* HLS_GROUP_SPS_PCM_FLAGS */
  if pcSPS.GetUsePCM() {
/*#if !HLS_GROUP_SPS_PCM_FLAGS
  this.WRITE_FLAG( pcSPS.GetPCMFilterDisableFlag()?1 : 0,                                "pcm_loop_filter_disable_flag");
#else*/ /* HLS_GROUP_SPS_PCM_FLAGS */
    this.WRITE_CODE( pcSPS.GetPCMBitDepthLuma() - 1, 4,                                  "pcm_sample_bit_depth_luma_minus1" );
    this.WRITE_CODE( pcSPS.GetPCMBitDepthChroma() - 1, 4,                                "pcm_sample_bit_depth_chroma_minus1" );
    this.WRITE_UVLC( pcSPS.GetPCMLog2MinSize() - 3,                                      "log2_min_pcm_luma_coding_block_size_minus3" );
    this.WRITE_UVLC( pcSPS.GetPCMLog2MaxSize() - pcSPS.GetPCMLog2MinSize(),             "log2_diff_max_min_pcm_luma_coding_block_size" );
    this.WRITE_FLAG( uint(TLibCommon.B2U(pcSPS.GetPCMFilterDisableFlag())),                              "pcm_loop_filter_disable_flag");
//#endif /* HLS_GROUP_SPS_PCM_FLAGS */
  }

  //assert( pcSPS.GetMaxTLayers() > 0 );
/*#if !MOVE_SPS_TEMPORAL_ID_NESTING_FLAG
  this.WRITE_FLAG( pcSPS.GetTemporalIdNestingFlag() ? 1 : 0,                             "temporal_id_nesting_flag" );
#endif*/

  rpsList := pcSPS.GetRPSList();
  var rps *TLibCommon.TComReferencePictureSet;

  this.WRITE_UVLC(uint(rpsList.GetNumberOfReferencePictureSets()), "num_short_term_ref_pic_sets" );
  for i:=0; i < rpsList.GetNumberOfReferencePictureSets(); i++ {
    rps = rpsList.GetReferencePictureSet(i);
//#if SPS_INTER_REF_SET_PRED
    this.codeShortTermRefPicSet(pcSPS,rps,false, i);
/*#else
    codeShortTermRefPicSet(pcSPS,rps,false);
#endif*/
  }
  this.WRITE_FLAG( uint(TLibCommon.B2U(pcSPS.GetLongTermRefsPresent())),         "long_term_ref_pics_present_flag" );
  if pcSPS.GetLongTermRefsPresent() {
    this.WRITE_UVLC(pcSPS.GetNumLongTermRefPicSPS(), "num_long_term_ref_pic_sps" );
    for k := uint(0); k < pcSPS.GetNumLongTermRefPicSPS(); k++ {
      this.WRITE_CODE( pcSPS.GetLtRefPicPocLsbSps(k), pcSPS.GetBitsForPOC(), "lt_ref_pic_poc_lsb_sps");
      this.WRITE_FLAG( uint(TLibCommon.B2U(pcSPS.GetUsedByCurrPicLtSPSFlag(int(k)))), "used_by_curr_pic_lt_sps_flag");
    }
  }
  this.WRITE_FLAG( uint(TLibCommon.B2U(pcSPS.GetTMVPFlagsPresent())),           "sps_temporal_mvp_enable_flag" );

//#if STRONG_INTRA_SMOOTHING
  this.WRITE_FLAG( uint(TLibCommon.B2U(pcSPS.GetUseStrongIntraSmoothing())),             "sps_strong_intra_smoothing_enable_flag" );
//#endif

  this.WRITE_FLAG( uint(TLibCommon.B2U(pcSPS.GetVuiParametersPresentFlag())),             "vui_parameters_present_flag" );
  if pcSPS.GetVuiParametersPresentFlag() {
      this.codeVUI(pcSPS.GetVuiParameters(), pcSPS);
  }

  this.WRITE_FLAG( 0, "sps_extension_flag" );
}
func (this *TEncCavlc)  codePPS                 ( pcPPS *TLibCommon.TComPPS){
//#if ENC_DEC_TRACE
  this.xTracePPSHeader (pcPPS);
//#endif

  this.WRITE_UVLC( uint(pcPPS.GetPPSId()),                             "pic_parameter_set_id" );
  this.WRITE_UVLC( uint(pcPPS.GetSPSId()),                             "seq_parameter_set_id" );
//#if DEPENDENT_SLICE_SEGMENT_FLAGS
///#if DEPENDENT_SLICES
//  this.WRITE_FLAG( uint(TLibCommon.B2U(pcPPS.GetDependentSliceSegmentsEnabledFlag())), "dependent_slice_enabled_flag" );
//#endif
//#endif
  this.WRITE_FLAG( uint(TLibCommon.B2U(pcPPS.GetSignHideFlag())), "sign_data_hiding_flag" );
  this.WRITE_FLAG( uint(TLibCommon.B2U(pcPPS.GetCabacInitPresentFlag())),   "cabac_init_present_flag" );
  this.WRITE_UVLC( pcPPS.GetNumRefIdxL0DefaultActive()-1,     "num_ref_idx_l0_default_active_minus1");
  this.WRITE_UVLC( pcPPS.GetNumRefIdxL1DefaultActive()-1,     "num_ref_idx_l1_default_active_minus1");

  this.WRITE_SVLC( pcPPS.GetPicInitQPMinus26(),                  "pic_init_qp_minus26");
  this.WRITE_FLAG( uint(TLibCommon.B2U(pcPPS.GetConstrainedIntraPred())),      "constrained_intra_pred_flag" );
  this.WRITE_FLAG( uint(TLibCommon.B2U(pcPPS.GetUseTransformSkip())),  "transform_skip_enabled_flag" );
  this.WRITE_FLAG( uint(TLibCommon.B2U(pcPPS.GetUseDQP())), "cu_qp_delta_enabled_flag" );
  if pcPPS.GetUseDQP() {
    this.WRITE_UVLC( pcPPS.GetMaxCuDQPDepth(), "diff_cu_qp_delta_depth" );
  }
  this.WRITE_SVLC( pcPPS.GetChromaCbQpOffset(),                   "cb_qp_offset" );
  this.WRITE_SVLC( pcPPS.GetChromaCrQpOffset(),                   "cr_qp_offset" );
  this.WRITE_FLAG( uint(TLibCommon.B2U(pcPPS.GetSliceChromaQpFlag())),          "slicelevel_chroma_qp_flag" );

  this.WRITE_FLAG( uint(TLibCommon.B2U(pcPPS.GetUseWP())),  "weighted_pred_flag" );   // Use of Weighting Prediction (P_SLICE)
  this.WRITE_FLAG( uint(TLibCommon.B2U(pcPPS.GetWPBiPred())), "weighted_bipred_flag" );  // Use of Weighting Bi-Prediction (B_SLICE)
  this.WRITE_FLAG( uint(TLibCommon.B2U(pcPPS.GetOutputFlagPresentFlag())),  "output_flag_present_flag" );
  this.WRITE_FLAG( uint(TLibCommon.B2U(pcPPS.GetTransquantBypassEnableFlag())), "transquant_bypass_enable_flag" );
/*#if !DEPENDENT_SLICE_SEGMENT_FLAGS
#if DEPENDENT_SLICES
  this.WRITE_FLAG( pcPPS.GetDependentSliceSegmentsEnabledFlag()    ? 1 : 0, "dependent_slice_enabled_flag" );
#endif
#endif*/
  this.WRITE_FLAG( uint(TLibCommon.B2U(pcPPS.GetTilesEnabledFlag())), "tiles_enabled_flag" );
  this.WRITE_FLAG( uint(TLibCommon.B2U(pcPPS.GetEntropyCodingSyncEnabledFlag())), "entropy_coding_sync_enabled_flag" );
/*#if !REMOVE_ENTROPY_SLICES
  this.WRITE_FLAG( pcPPS.GetEntropySliceEnabledFlag()      ? 1 : 0, "entropy_slice_enabled_flag" );
#endif*/
  if pcPPS.GetTilesEnabledFlag() {
    this.WRITE_UVLC( uint(pcPPS.GetNumColumnsMinus1()),                                    "num_tile_columns_minus1" );
    this.WRITE_UVLC( uint(pcPPS.GetNumRowsMinus1()),                                       "num_tile_rows_minus1" );
    this.WRITE_FLAG( uint(TLibCommon.B2U(pcPPS.GetUniformSpacingFlag())),                                  "uniform_spacing_flag" );
    if pcPPS.GetUniformSpacingFlag() == false {
      for i:=uint(0); i<uint(pcPPS.GetNumColumnsMinus1()); i++ {
        this.WRITE_UVLC( uint(pcPPS.GetColumnWidth(int(i))-1),                                  "column_width_minus1" );
      }
      for i:=uint(0); i<uint(pcPPS.GetNumRowsMinus1()); i++ {
        this.WRITE_UVLC( uint(pcPPS.GetRowHeight(int(i))-1),                                    "row_height_minus1" );
      }
    }
    if pcPPS.GetNumColumnsMinus1() !=0 || pcPPS.GetNumRowsMinus1() !=0 {
      this.WRITE_FLAG( uint(TLibCommon.B2U(pcPPS.GetLoopFilterAcrossTilesEnabledFlag())),          "loop_filter_across_tiles_enabled_flag");
    }
  }
  this.WRITE_FLAG( uint(TLibCommon.B2U(pcPPS.GetLoopFilterAcrossSlicesEnabledFlag())),        "loop_filter_across_slices_enabled_flag");
  this.WRITE_FLAG( uint(TLibCommon.B2U(pcPPS.GetDeblockingFilterControlPresentFlag())),       "deblocking_filter_control_present_flag");
  if pcPPS.GetDeblockingFilterControlPresentFlag() {
    this.WRITE_FLAG( uint(TLibCommon.B2U(pcPPS.GetDeblockingFilterOverrideEnabledFlag())),  "deblocking_filter_override_enabled_flag" );
    this.WRITE_FLAG( uint(TLibCommon.B2U(pcPPS.GetPicDisableDeblockingFilterFlag())),               "pic_disable_deblocking_filter_flag" );
    if !pcPPS.GetPicDisableDeblockingFilterFlag() {
      this.WRITE_SVLC( pcPPS.GetDeblockingFilterBetaOffsetDiv2(),             "pps_beta_offset_div2" );
      this.WRITE_SVLC( pcPPS.GetDeblockingFilterTcOffsetDiv2(),               "pps_tc_offset_div2" );
    }
  }
  this.WRITE_FLAG( uint(TLibCommon.B2U(pcPPS.GetScalingListPresentFlag())),                          "pps_scaling_list_data_present_flag" );
  if pcPPS.GetScalingListPresentFlag() {
/*#if SCALING_LIST_OUTPUT_RESULT
    printf("PPS\n");
#endif*/
    this.codeScalingList( this.m_pcSlice.GetScalingList() );
  }
//#if HLS_MOVE_SPS_PICLIST_FLAGS
  this.WRITE_FLAG( uint(TLibCommon.B2U(pcPPS.GetListsModificationPresentFlag())), "lists_modification_present_flag");
//#endif /* HLS_MOVE_SPS_PICLIST_FLAGS */
  this.WRITE_UVLC( pcPPS.GetLog2ParallelMergeLevelMinus2(), "log2_parallel_merge_level_minus2");
//#if HLS_EXTRA_SLICE_HEADER_BITS
  this.WRITE_CODE( uint(pcPPS.GetNumExtraSliceHeaderBits()), 3, "num_extra_slice_header_bits");
//#endif /* HLS_EXTRA_SLICE_HEADER_BITS */
  this.WRITE_FLAG( uint(TLibCommon.B2U(pcPPS.GetSliceHeaderExtensionPresentFlag())), "slice_header_extension_present_flag");
  this.WRITE_FLAG( 0, "pps_extension_flag" );
}
func (this *TEncCavlc)  codeSliceHeader         ( pcSlice *TLibCommon.TComSlice){
//#if ENC_DEC_TRACE
  this.xTraceSliceHeader (pcSlice);
//#endif

  //calculate number of bits required for slice address
  maxAddrOuter := pcSlice.GetPic().GetNumCUsInFrame();
  reqBitsOuter := 0;
  for maxAddrOuter>(1<<uint(reqBitsOuter)) {
    reqBitsOuter++;
  }
  maxAddrInner := pcSlice.GetPic().GetNumPartInCU()>>(2);
  maxAddrInner = 1;
  reqBitsInner := 0;

  for maxAddrInner>(1<<uint(reqBitsInner)) {
    reqBitsInner++;
  }
  var lCUAddress, innerAddress int;
  if pcSlice.IsNextSlice() {
    // Calculate slice address
    lCUAddress = int(pcSlice.GetSliceCurStartCUAddr()/pcSlice.GetPic().GetNumPartInCU());
    innerAddress = 0;
  }else{
    // Calculate slice address
    lCUAddress = int(pcSlice.GetSliceSegmentCurStartCUAddr()/pcSlice.GetPic().GetNumPartInCU());
    innerAddress = 0;
  }
  //write slice address
  address := uint(pcSlice.GetPic().GetPicSym().GetCUOrderMap(lCUAddress) << uint(reqBitsInner)) + uint(innerAddress);
  this.WRITE_FLAG( uint(TLibCommon.B2U(address==0)), "first_slice_in_pic_flag" );
  if pcSlice.GetRapPicFlag() {
    this.WRITE_FLAG( 0, "no_output_of_prior_pics_flag" );
  }
  this.WRITE_UVLC( uint(pcSlice.GetPPS().GetPPSId()), "pic_parameter_set_id" );
//#if DEPENDENT_SLICE_SEGMENT_FLAGS
  pcSlice.SetDependentSliceSegmentFlag(!pcSlice.IsNextSlice());
  if pcSlice.GetPPS().GetDependentSliceSegmentsEnabledFlag() && (address!=0) {
    this.WRITE_FLAG( uint(TLibCommon.B2U(pcSlice.GetDependentSliceSegmentFlag())), "dependent_slice_flag" );
  }
//#endif
  if address>0 {
    this.WRITE_CODE( address, uint(reqBitsOuter+reqBitsInner), "slice_address" );
  }
/*#if !DEPENDENT_SLICE_SEGMENT_FLAGS
  pcSlice->setDependentSliceFlag(!pcSlice->isNextSlice());
  if ( pcSlice.GetPPS().GetDependentSliceSegmentsEnabledFlag() && (address!=0) )
  {
    this.WRITE_FLAG( pcSlice.GetDependentSliceSegmentFlag() ? 1 : 0, "dependent_slice_flag" );
  }
#endif*/
  if !pcSlice.GetDependentSliceSegmentFlag() {
//#if HLS_EXTRA_SLICE_HEADER_BITS
    for i := 0; i < pcSlice.GetPPS().GetNumExtraSliceHeaderBits(); i++ {
      //assert(!!"slice_reserved_undetermined_flag[]");
      this.WRITE_FLAG(0, "slice_reserved_undetermined_flag[]");
    }
//#endif /* HLS_EXTRA_SLICE_HEADER_BITS */

    this.WRITE_UVLC( uint(pcSlice.GetSliceType()),       "slice_type" );

    if pcSlice.GetPPS().GetOutputFlagPresentFlag() {
      this.WRITE_FLAG( uint(TLibCommon.B2U(pcSlice.GetPicOutputFlag())), "pic_output_flag" );
    }

    // in the first version chroma_format_idc is equal to one, thus colour_plane_id will not be present
    //assert (pcSlice.GetSPS().GetChromaFormatIdc() == 1 );
    // if( separate_colour_plane_flag  ==  1 )
    //   colour_plane_id                                      u(2)

    if !pcSlice.GetIdrPicFlag() {
      picOrderCntLSB := uint(pcSlice.GetPOC()-pcSlice.GetLastIDR()+(1<<pcSlice.GetSPS().GetBitsForPOC()))%(1<<pcSlice.GetSPS().GetBitsForPOC());
      this.WRITE_CODE( picOrderCntLSB, pcSlice.GetSPS().GetBitsForPOC(), "pic_order_cnt_lsb");
      rps := pcSlice.GetRPS();
      if pcSlice.GetRPSidx() < 0 {
        this.WRITE_FLAG( 0, "short_term_ref_pic_set_sps_flag");
//#if SPS_INTER_REF_SET_PRED
        this.codeShortTermRefPicSet(pcSlice.GetSPS(), rps, true, pcSlice.GetSPS().GetRPSList().GetNumberOfReferencePictureSets());
/*#else
        codeShortTermRefPicSet(pcSlice.GetSPS(), rps, true);
#endif*/
      }else{
        this.WRITE_FLAG( 1, "short_term_ref_pic_set_sps_flag");
        numBits := uint(0);
        for (1 << numBits) < pcSlice.GetSPS().GetRPSList().GetNumberOfReferencePictureSets() {
          numBits++;
        }
        if numBits > 0 {
          this.WRITE_CODE( uint(pcSlice.GetRPSidx()), numBits, "short_term_ref_pic_set_idx" );
        }
      }
      if pcSlice.GetSPS().GetLongTermRefsPresent() {
        numLtrpInSH := rps.GetNumberOfLongtermPictures();
        var ltrpInSPS	[TLibCommon.MAX_NUM_REF_PICS]int;
        numLtrpInSPS := 0;
        var ltrpIndex uint;
        counter := 0;
        for k := rps.GetNumberOfPictures()-1; k > rps.GetNumberOfPictures()-rps.GetNumberOfLongtermPictures()-1; k-- {
          if this.findMatchingLTRP(pcSlice, &ltrpIndex, rps.GetPOC(k), rps.GetUsed(k)) {
            ltrpInSPS[numLtrpInSPS] = int(ltrpIndex);
            numLtrpInSPS++;
          }else{
            counter++;
          }
        }
        numLtrpInSH -= numLtrpInSPS;

        bitsForLtrpInSPS := uint(1);
        for pcSlice.GetSPS().GetNumLongTermRefPicSPS() > (1 << bitsForLtrpInSPS) {
          bitsForLtrpInSPS++;
        }
        if pcSlice.GetSPS().GetNumLongTermRefPicSPS() > 0 {
          this.WRITE_UVLC( uint(numLtrpInSPS), "num_long_term_sps");
        }
        this.WRITE_UVLC( uint(numLtrpInSH), "num_long_term_pics");
        // Note that the LSBs of the LT ref. pic. POCs must be sorted before.
        // Not sorted here because LT ref indices will be used in setRefPicList()
        prevDeltaMSB := 0;
        prevLSB := 0;
        offset := rps.GetNumberOfNegativePictures() + rps.GetNumberOfPositivePictures();
        for i:=rps.GetNumberOfPictures()-1 ; i > offset-1; i-- {
          if counter < numLtrpInSPS {
            this.WRITE_CODE( uint(ltrpInSPS[counter]), bitsForLtrpInSPS, "lt_idx_sps[i]");
          }else{
            this.WRITE_CODE( uint(rps.GetPocLSBLT(i)), pcSlice.GetSPS().GetBitsForPOC(), "poc_lsb_lt");
            this.WRITE_FLAG( uint(TLibCommon.B2U(rps.GetUsed(i))), "used_by_curr_pic_lt_flag");
          }
          this.WRITE_FLAG( uint(TLibCommon.B2U(rps.GetDeltaPocMSBPresentFlag(i))), "delta_poc_msb_present_flag");

          if rps.GetDeltaPocMSBPresentFlag(i) {
            deltaFlag := false;
            //  First LTRP from SPS                 ||  First LTRP from SH                              || curr LSB            != prev LSB
            if (i == rps.GetNumberOfPictures()-1) || (i == rps.GetNumberOfPictures()-1-numLtrpInSPS) || (rps.GetPocLSBLT(i) != prevLSB) {
              deltaFlag = true;
            }
            if deltaFlag {
              this.WRITE_UVLC( uint(rps.GetDeltaPocMSBCycleLT(i)), "delta_poc_msb_cycle_lt[i]" );
            }else{
              differenceInDeltaMSB := rps.GetDeltaPocMSBCycleLT(i) - prevDeltaMSB;
              //assert(differenceInDeltaMSB >= 0);
              this.WRITE_UVLC( uint(differenceInDeltaMSB), "delta_poc_msb_cycle_lt[i]" );
            }
            prevLSB = rps.GetPocLSBLT(i);
            prevDeltaMSB = rps.GetDeltaPocMSBCycleLT(i);
          }
        }
      }
    }
    if pcSlice.GetSPS().GetUseSAO() {
      if pcSlice.GetSPS().GetUseSAO() {
         this.WRITE_FLAG( uint(TLibCommon.B2U(pcSlice.GetSaoEnabledFlag())), "slice_sao_luma_flag" );
         {
           saoParam := pcSlice.GetPic().GetPicSym().GetSaoParam();
           this.WRITE_FLAG( uint(TLibCommon.B2U(saoParam.SaoFlag[1])), "slice_sao_chroma_flag" );
         }
      }
    }

    //check if numrefidxes match the defaults. If not, override

//#if K0251
    if !pcSlice.GetIdrPicFlag() {
/*#else
    if (!pcSlice->isIntra())
#endif*/
      if pcSlice.GetSPS().GetTMVPFlagsPresent() {
        this.WRITE_FLAG( uint(TLibCommon.B2U(pcSlice.GetEnableTMVPFlag())), "enable_temporal_mvp_flag" );
      }
//#if K0251
    }

    if !pcSlice.IsIntra() {
//#endif
      overrideFlag := (uint(pcSlice.GetNumRefIdx( TLibCommon.REF_PIC_LIST_0 ))!=pcSlice.GetPPS().GetNumRefIdxL0DefaultActive()||(pcSlice.IsInterB()&&uint(pcSlice.GetNumRefIdx( TLibCommon.REF_PIC_LIST_1 ))!=pcSlice.GetPPS().GetNumRefIdxL1DefaultActive()));
      this.WRITE_FLAG( uint(TLibCommon.B2U(overrideFlag)),                               "num_ref_idx_active_override_flag");
      if overrideFlag {
        this.WRITE_UVLC( uint(pcSlice.GetNumRefIdx( TLibCommon.REF_PIC_LIST_0 ) - 1),      "num_ref_idx_l0_active_minus1" );
        if pcSlice.IsInterB() {
          this.WRITE_UVLC( uint(pcSlice.GetNumRefIdx( TLibCommon.REF_PIC_LIST_1 ) - 1),    "num_ref_idx_l1_active_minus1" );
        }else{
          pcSlice.SetNumRefIdx(TLibCommon.REF_PIC_LIST_1, 0);
        }
      }
    }else{
      pcSlice.SetNumRefIdx(TLibCommon.REF_PIC_LIST_0, 0);
      pcSlice.SetNumRefIdx(TLibCommon.REF_PIC_LIST_1, 0);
    }

//#if SAVE_BITS_REFPICLIST_MOD_FLAG
/*#if !HLS_MOVE_SPS_PICLIST_FLAGS
    if( pcSlice.GetSPS().GetListsModificationPresentFlag() && pcSlice.GetNumRpsCurrTempList() > 1)
#else*/ /* HLS_MOVE_SPS_PICLIST_FLAGS */
    if pcSlice.GetPPS().GetListsModificationPresentFlag() && pcSlice.GetNumRpsCurrTempList() > 1 {
//#endif /* HLS_MOVE_SPS_PICLIST_FLAGS */
/*#else
#if !HLS_MOVE_SPS_PICLIST_FLAGS
    if( pcSlice.GetSPS().GetListsModificationPresentFlag() )
#else*/ /* HLS_MOVE_SPS_PICLIST_FLAGS */
//    if( pcSlice.GetPPS().GetListsModificationPresentFlag() )
//#endif /* HLS_MOVE_SPS_PICLIST_FLAGS */
//#endif
      refPicListModification := pcSlice.GetRefPicListModification();
      if !pcSlice.IsIntra() {
        this.WRITE_FLAG( uint(TLibCommon.B2U(pcSlice.GetRefPicListModification().GetRefPicListModificationFlagL0())),       "ref_pic_list_modification_flag_l0" );
        if pcSlice.GetRefPicListModification().GetRefPicListModificationFlagL0() {
          numRpsCurrTempList0 := pcSlice.GetNumRpsCurrTempList();
          if numRpsCurrTempList0 > 1 {
            length := 1;
            numRpsCurrTempList0 --;
            numRpsCurrTempList0 >>= 1;
            for numRpsCurrTempList0 !=0 {
              length ++;
              numRpsCurrTempList0 >>= 1
            }
            for i := 0; i < pcSlice.GetNumRefIdx( TLibCommon.REF_PIC_LIST_0 ); i++ {
              this.WRITE_CODE( refPicListModification.GetRefPicSetIdxL0(uint(i)), uint(length), "list_entry_l0");
            }
          }
        }
      }
      if pcSlice.IsInterB() {
        this.WRITE_FLAG( uint(TLibCommon.B2U(pcSlice.GetRefPicListModification().GetRefPicListModificationFlagL1())),       "ref_pic_list_modification_flag_l1" );
        if pcSlice.GetRefPicListModification().GetRefPicListModificationFlagL1() {
          numRpsCurrTempList1 := pcSlice.GetNumRpsCurrTempList();
          if numRpsCurrTempList1 > 1 {
            length := 1;
            numRpsCurrTempList1 --;
            numRpsCurrTempList1 >>= 1
            for numRpsCurrTempList1 != 0 {
              length ++;
              numRpsCurrTempList1 >>= 1
            }
            for i := 0; i < pcSlice.GetNumRefIdx( TLibCommon.REF_PIC_LIST_1 ); i++ {
              this.WRITE_CODE( refPicListModification.GetRefPicSetIdxL1(uint(i)), uint(length), "list_entry_l1");
            }
          }
        }
      }
    }

    if pcSlice.IsInterB() {
      this.WRITE_FLAG( uint(TLibCommon.B2U(pcSlice.GetMvdL1ZeroFlag())),   "mvd_l1_zero_flag");
    }

    if !pcSlice.IsIntra() {
      if !pcSlice.IsIntra() && pcSlice.GetPPS().GetCabacInitPresentFlag() {
        sliceType   := pcSlice.GetSliceType();
        encCABACTableIdx := pcSlice.GetPPS().GetEncCABACTableIdx();
        encCabacInitFlag := (uint(sliceType)!=encCABACTableIdx && encCABACTableIdx!=TLibCommon.I_SLICE);
        pcSlice.SetCabacInitFlag( encCabacInitFlag );
        this.WRITE_FLAG( uint(TLibCommon.B2U(encCabacInitFlag)), "cabac_init_flag" );
      }
    }

    if pcSlice.GetEnableTMVPFlag() {
      if pcSlice.GetSliceType() == TLibCommon.B_SLICE {
        this.WRITE_FLAG( pcSlice.GetColFromL0Flag(), "collocated_from_l0_flag" );
      }

      if  pcSlice.GetSliceType() != TLibCommon.I_SLICE &&
        ((pcSlice.GetColFromL0Flag()==1 && pcSlice.GetNumRefIdx(TLibCommon.REF_PIC_LIST_0)>1)||
        (pcSlice.GetColFromL0Flag()==0  && pcSlice.GetNumRefIdx(TLibCommon.REF_PIC_LIST_1)>1)) {
        this.WRITE_UVLC( pcSlice.GetColRefIdx(), "collocated_ref_idx" );
      }
    }
    if (pcSlice.GetPPS().GetUseWP() && pcSlice.GetSliceType()==TLibCommon.P_SLICE) || (pcSlice.GetPPS().GetWPBiPred() && pcSlice.GetSliceType()==TLibCommon.B_SLICE) {
      this.xCodePredWeightTable( pcSlice );
    }
    //assert(pcSlice.GetMaxNumMergeCand()<=MRG_MAX_NUM_CANDS);
    if !pcSlice.IsIntra() {
      this.WRITE_UVLC(TLibCommon.MRG_MAX_NUM_CANDS - pcSlice.GetMaxNumMergeCand(), "five_minus_max_num_merge_cand");
    }
    iCode := pcSlice.GetSliceQp() - ( pcSlice.GetPPS().GetPicInitQPMinus26() + 26 );
    this.WRITE_SVLC( iCode, "slice_qp_delta" );
    if pcSlice.GetPPS().GetSliceChromaQpFlag() {
      iCode = pcSlice.GetSliceQpDeltaCb();
      this.WRITE_SVLC( iCode, "slice_qp_delta_cb" );
      iCode = pcSlice.GetSliceQpDeltaCr();
      this.WRITE_SVLC( iCode, "slice_qp_delta_cr" );
    }
    if pcSlice.GetPPS().GetDeblockingFilterControlPresentFlag() {
      if pcSlice.GetPPS().GetDeblockingFilterOverrideEnabledFlag() {
        this.WRITE_FLAG(uint(TLibCommon.B2U(pcSlice.GetDeblockingFilterOverrideFlag())), "deblocking_filter_override_flag");
      }
      if pcSlice.GetDeblockingFilterOverrideFlag() {
        this.WRITE_FLAG(uint(TLibCommon.B2U(pcSlice.GetDeblockingFilterDisable())), "slice_disable_deblocking_filter_flag");
        if !pcSlice.GetDeblockingFilterDisable() {
          this.WRITE_SVLC (pcSlice.GetDeblockingFilterBetaOffsetDiv2(), "beta_offset_div2");
          this.WRITE_SVLC (pcSlice.GetDeblockingFilterTcOffsetDiv2(),   "tc_offset_div2");
        }
      }
    }

    var isSAOEnabled bool;
    if !pcSlice.GetSPS().GetUseSAO() {
    	isSAOEnabled = false;
    }else{
    	isSAOEnabled = (pcSlice.GetSaoEnabledFlag()||pcSlice.GetSaoEnabledFlagChroma());
    }
    isDBFEnabled := (!pcSlice.GetDeblockingFilterDisable());

    if pcSlice.GetPPS().GetLoopFilterAcrossSlicesEnabledFlag() && ( isSAOEnabled || isDBFEnabled ) {
      this.WRITE_FLAG(uint(TLibCommon.B2U(pcSlice.GetLFCrossSliceBoundaryFlag())), "slice_loop_filter_across_slices_enabled_flag");
    }
  }
  if pcSlice.GetPPS().GetSliceHeaderExtensionPresentFlag() {
    this.WRITE_UVLC(0,"slice_header_extension_length");
  }
}

func (this *TEncCavlc)  codePTL                 ( pcPTL *TLibCommon.TComPTL, profilePresentFlag bool, maxNumSubLayersMinus1 int){
  if profilePresentFlag {
    this.codeProfileTier(pcPTL.GetGeneralPTL());    // general_...
  }
  this.WRITE_CODE( uint(pcPTL.GetGeneralPTL().GetLevelIdc()), 8, "general_level_idc" );

  for i := 0; i < maxNumSubLayersMinus1; i++ {
//#if CONDITION_SUBLAYERPROFILEPRESENTFLAG
    if profilePresentFlag {
      this.WRITE_FLAG( uint(TLibCommon.B2U(pcPTL.GetSubLayerProfilePresentFlag(i))), "sub_layer_profile_present_flag[i]" );
    }
//#else
//    this.WRITE_FLAG( pcPTL.GetSubLayerProfilePresentFlag(i), "sub_layer_profile_present_flag[i]" );
//#endif

    this.WRITE_FLAG( uint(TLibCommon.B2U(pcPTL.GetSubLayerLevelPresentFlag(i))),   "sub_layer_level_present_flag[i]" );
    if profilePresentFlag && pcPTL.GetSubLayerProfilePresentFlag(i) {
      this.codeProfileTier(pcPTL.GetSubLayerPTL(i));  // sub_layer_...
    }
    if pcPTL.GetSubLayerLevelPresentFlag(i) {
      this.WRITE_CODE( uint(pcPTL.GetSubLayerPTL(i).GetLevelIdc()), 8, "sub_layer_level_idc[i]" );
    }
  }
}

func (this *TEncCavlc)  codeProfileTier         ( ptl *TLibCommon.ProfileTierLevel){
  this.WRITE_CODE( uint(ptl.GetProfileSpace()), 2 ,     "XXX_profile_space[]");
  this.WRITE_FLAG( uint(TLibCommon.B2U(ptl.GetTierFlag    ())),         "XXX_tier_flag[]"    );
  this.WRITE_CODE( uint(ptl.GetProfileIdc  ()), 5 ,     "XXX_profile_idc[]"  );
  for j := 0; j < 32; j++ {
    this.WRITE_FLAG( uint(TLibCommon.B2U(ptl.GetProfileCompatibilityFlag(j))), "XXX_profile_compatibility_flag[][j]");
  }
  this.WRITE_CODE(0 , 16, "XXX_reserved_zero_16bits[]");
}

//#if SIGNAL_BITRATE_PICRATE_IN_VPS
func (this *TEncCavlc)  codeBitratePicRateInfo(info *TLibCommon.TComBitRatePicRateInfo, tempLevelLow, tempLevelHigh int){
  for i := tempLevelLow; i <= tempLevelHigh; i++ {
    this.WRITE_FLAG( uint(TLibCommon.B2U(info.GetBitRateInfoPresentFlag(i))),  "bit_rate_info_present_flag[i]" );
    this.WRITE_FLAG( uint(TLibCommon.B2U(info.GetPicRateInfoPresentFlag(i))),  "pic_rate_info_present_flag[i]" );
    if info.GetBitRateInfoPresentFlag(i) {
      this.WRITE_CODE( uint(info.GetAvgBitRate(i)), 16, "avg_bit_rate[i]" );
      this.WRITE_CODE( uint(info.GetMaxBitRate(i)), 16, "max_bit_rate[i]" );
    }
    if info.GetPicRateInfoPresentFlag(i) {
      this.WRITE_CODE( uint(info.GetConstantPicRateIdc(i)),  2, "constant_pic_rate_idc[i]" );
      this.WRITE_CODE( uint(info.GetAvgPicRate(i)),         16, "avg_pic_rate[i]"          );
    }
  }
}
//#endif

func (this *TEncCavlc)  codeTilesWPPEntryPoint( pSlice *TLibCommon.TComSlice){
  if !pSlice.GetPPS().GetTilesEnabledFlag() && !pSlice.GetPPS().GetEntropyCodingSyncEnabledFlag() {
    return;
  }
  numEntryPointOffsets := uint(0);
  offsetLenMinus1 := uint(0);
  maxOffset := uint(0);
  numZeroSubstreamsAtStartOfSlice := 0;
  var entryPointOffset []uint;
  if pSlice.GetPPS().GetTilesEnabledFlag() {
    numEntryPointOffsets = pSlice.GetTileLocationCount();
    entryPointOffset     = make([]uint, numEntryPointOffsets);
    for idx:=0; idx<int(pSlice.GetTileLocationCount()); idx++ {
      if idx == 0 {
        entryPointOffset [ idx ] = pSlice.GetTileLocation( 0 );
      }else{
        entryPointOffset [ idx ] = pSlice.GetTileLocation( idx ) - pSlice.GetTileLocation( idx-1 );
      }

      if entryPointOffset[ idx ] > maxOffset {
        maxOffset = entryPointOffset[ idx ];
      }
    }
  }else if pSlice.GetPPS().GetEntropyCodingSyncEnabledFlag() {
    pSubstreamSizes                  := pSlice.GetSubstreamSizes();
    maxNumParts                      := pSlice.GetPic().GetNumPartInCU();
//#if DEPENDENT_SLICES
    numZeroSubstreamsAtStartOfSlice   = int(pSlice.GetSliceSegmentCurStartCUAddr()/maxNumParts/pSlice.GetPic().GetFrameWidthInCU());
    numZeroSubstreamsAtEndOfSlice    := int(pSlice.GetPic().GetFrameHeightInCU()-1 - ((pSlice.GetSliceSegmentCurEndCUAddr()-1)/maxNumParts/pSlice.GetPic().GetFrameWidthInCU()));
/*#else
    numZeroSubstreamsAtStartOfSlice       = pSlice.GetSliceCurStartCUAddr()/maxNumParts/pSlice.GetPic().GetFrameWidthInCU();
    Int  numZeroSubstreamsAtEndOfSlice    = pSlice.GetPic().GetFrameHeightInCU()-1 - ((pSlice.GetSliceCurEndCUAddr()-1)/maxNumParts/pSlice.GetPic().GetFrameWidthInCU());
#endif*/
    numEntryPointOffsets                  = uint(pSlice.GetPPS().GetNumSubstreams() - numZeroSubstreamsAtStartOfSlice - numZeroSubstreamsAtEndOfSlice - 1);
    pSlice.SetNumEntryPointOffsets(int(numEntryPointOffsets));
    entryPointOffset           = make([]uint, numEntryPointOffsets);
    for idx:=uint(0); idx<numEntryPointOffsets; idx++ {
      entryPointOffset[ idx ] = ( pSubstreamSizes[ int(idx)+numZeroSubstreamsAtStartOfSlice ] >> 3 ) ;
      if entryPointOffset[ idx ] > maxOffset {
        maxOffset = entryPointOffset[ idx ];
      }
    }
  }
  // Determine number of bits "offsetLenMinus1+1" required for entry point information
  offsetLenMinus1 = 0;
  for maxOffset >= (1 << (offsetLenMinus1 + 1)) {
    offsetLenMinus1++;
    //assert(offsetLenMinus1 + 1 < 32);
  }

  this.WRITE_UVLC(numEntryPointOffsets, "num_entry_point_offsets");
  if numEntryPointOffsets>0 {
    this.WRITE_UVLC(offsetLenMinus1, "offset_len_minus1");
  }

  for idx:=uint(0); idx<numEntryPointOffsets; idx++ {
    this.WRITE_CODE(entryPointOffset[ idx ], offsetLenMinus1+1, "entry_point_offset");
  }

  //delete [] entryPointOffset;
}

func (this *TEncCavlc)  codeTerminatingBit      ( uilsLast uint){}
func (this *TEncCavlc)  codeSliceFinish         () {}
func (this *TEncCavlc)  encodeStart             () {}

func (this *TEncCavlc) codeMVPIdx ( pcCU *TLibCommon.TComDataCU,  uiAbsPartIdx uint,  eRefList TLibCommon.RefPicList) {}
func (this *TEncCavlc) codeSAOSign       ( code uint  ) {  }
func (this *TEncCavlc) codeSaoMaxUvlc    ( code, maxSymbol uint){}
func (this *TEncCavlc) codeSaoMerge  	 ( uiCode uint){}
func (this *TEncCavlc) codeSaoTypeIdx    ( uiCode uint){}
func (this *TEncCavlc) codeSaoUflc       ( uiLength, uiCode uint){ }

func (this *TEncCavlc) codeCUTransquantBypassFlag( pcCU *TLibCommon.TComDataCU,  uiAbsPartIdx uint ) {}
func (this *TEncCavlc) codeSkipFlag      ( pcCU *TLibCommon.TComDataCU,  uiAbsPartIdx uint ) {}
func (this *TEncCavlc) codeMergeFlag     ( pcCU *TLibCommon.TComDataCU,  uiAbsPartIdx uint ) {}
func (this *TEncCavlc) codeMergeIndex    ( pcCU *TLibCommon.TComDataCU,  uiAbsPartIdx uint ) {}

func (this *TEncCavlc) codeInterModeFlag( pcCU *TLibCommon.TComDataCU,  uiAbsPartIdx uint,  uiDepth,  uiEncMode uint) {}
func (this *TEncCavlc) codeSplitFlag     ( pcCU *TLibCommon.TComDataCU,  uiAbsPartIdx uint, uiDepth uint) {}

func (this *TEncCavlc) codePartSize      ( pcCU *TLibCommon.TComDataCU,  uiAbsPartIdx uint, uiDepth uint) {}
func (this *TEncCavlc) codePredMode      ( pcCU *TLibCommon.TComDataCU,  uiAbsPartIdx uint ) {}

//#if !REMOVE_BURST_IPCM
//func (this *TEncCavlc) codeIPCMInfo      ( pcCU *TLibCommon.TComDataCU,  uiAbsPartIdx uint, Int numIPCM, Bool firstIPCMFlag);
//#else
func (this *TEncCavlc) codeIPCMInfo      ( pcCU *TLibCommon.TComDataCU,  uiAbsPartIdx uint ) {}
//#endif

func (this *TEncCavlc) codeTransformSubdivFlag( uiSymbol, uiCtx uint ) {}
func (this *TEncCavlc) codeQtCbf         ( pcCU *TLibCommon.TComDataCU,  uiAbsPartIdx uint, eType TLibCommon.TextType, uiTrDepth uint) {}
func (this *TEncCavlc) codeQtRootCbf     ( pcCU *TLibCommon.TComDataCU,  uiAbsPartIdx uint ) {}
func (this *TEncCavlc) codeQtCbfZero     ( pcCU *TLibCommon.TComDataCU,  uiAbsPartIdx uint, eType TLibCommon.TextType, uiTrDepth uint ) {}
func (this *TEncCavlc) codeQtRootCbfZero  ( pcCU *TLibCommon.TComDataCU,  uiAbsPartIdx uint ) {}
func (this *TEncCavlc) codeIntraDirLumaAng( pcCU *TLibCommon.TComDataCU,  uiAbsPartIdx uint, isMultiple bool) {}
func (this *TEncCavlc) codeIntraDirChroma( pcCU *TLibCommon.TComDataCU,  uiAbsPartIdx uint ) {}
func (this *TEncCavlc) codeInterDir      ( pcCU *TLibCommon.TComDataCU,  uiAbsPartIdx uint ) {}
func (this *TEncCavlc) codeRefFrmIdx     ( pcCU *TLibCommon.TComDataCU,  uiAbsPartIdx uint, eRefList TLibCommon.RefPicList ) {}
func (this *TEncCavlc) codeMvd           ( pcCU *TLibCommon.TComDataCU,  uiAbsPartIdx uint, eRefList TLibCommon.RefPicList ) {}
func (this *TEncCavlc) codeDeltaQP       ( pcCU *TLibCommon.TComDataCU,  uiAbsPartIdx uint ){}
func (this *TEncCavlc) codeCoeffNxN      ( pcCU *TLibCommon.TComDataCU, pcCoef []TLibCommon.TCoeff, uiAbsPartIdx, uiWidth, uiHeight, uiDepth uint, eTType TLibCommon.TextType ) {}
func (this *TEncCavlc) codeTransformSkipFlags ( pcCU *TLibCommon.TComDataCU,  uiAbsPartIdx uint, width, height, uiDepth uint, eTType TLibCommon.TextType ) {}
func (this *TEncCavlc) estBit               (pcEstBitsSbac *TLibCommon.EstBitsSbacStruct, width, height int, eTType TLibCommon.TextType ) {}

func (this *TEncCavlc) xCodePredWeightTable          ( pcSlice *TLibCommon.TComSlice) {
  var	wp	[]TLibCommon.WpScalingParam;
  bChroma      := true; // color always present in HEVC ?
  var             iNbRef   int;
  if pcSlice.GetSliceType() == TLibCommon.B_SLICE {
  	iNbRef = 2;
  }else{
  	iNbRef = 1;
  }
  bDenomCoded  := false;
  uiMode := uint(0);
  uiTotalSignalledWeightFlags := uint(0);
  if (pcSlice.GetSliceType()==TLibCommon.P_SLICE && pcSlice.GetPPS().GetUseWP()) || (pcSlice.GetSliceType()==TLibCommon.B_SLICE && pcSlice.GetPPS().GetWPBiPred()) {
    uiMode = 1; // explicit
  }
  if uiMode == 1 {

    for iNumRef:=0 ; iNumRef<iNbRef ; iNumRef++ {
      var eRefPicList TLibCommon.RefPicList;
      if iNumRef!=0 {
      	eRefPicList = TLibCommon.REF_PIC_LIST_1;
      }else{
      	eRefPicList = TLibCommon.REF_PIC_LIST_0;
      }

      for iRefIdx:=0 ; iRefIdx<pcSlice.GetNumRefIdx(eRefPicList) ; iRefIdx++ {
        wp = pcSlice.GetWpScaling(eRefPicList, iRefIdx);
        if !bDenomCoded {
          var iDeltaDenom int;
          this.WRITE_UVLC( wp[0].GetLog2WeightDenom(), "luma_log2_weight_denom" );     // ue(v): luma_log2_weight_denom

          if bChroma {
            iDeltaDenom = int(wp[1].GetLog2WeightDenom() - wp[0].GetLog2WeightDenom());
            this.WRITE_SVLC( iDeltaDenom, "delta_chroma_log2_weight_denom" );       // se(v): delta_chroma_log2_weight_denom
          }
          bDenomCoded = true;
        }
        this.WRITE_FLAG( uint(TLibCommon.B2U(wp[0].GetPresentFlag())), "luma_weight_lX_flag" );               // u(1): luma_weight_lX_flag
        uiTotalSignalledWeightFlags += uint(TLibCommon.B2U(wp[0].GetPresentFlag()));
      }
      if bChroma {
        for iRefIdx:=0 ; iRefIdx<pcSlice.GetNumRefIdx(eRefPicList) ; iRefIdx++ {
          wp = pcSlice.GetWpScaling(eRefPicList, iRefIdx);
          this.WRITE_FLAG( uint(TLibCommon.B2U(wp[1].GetPresentFlag())), "chroma_weight_lX_flag" );           // u(1): chroma_weight_lX_flag
          uiTotalSignalledWeightFlags += 2*uint(TLibCommon.B2U(wp[1].GetPresentFlag()));
        }
      }

      for iRefIdx:=0 ; iRefIdx<pcSlice.GetNumRefIdx(eRefPicList) ; iRefIdx++ {
        wp = pcSlice.GetWpScaling(eRefPicList, iRefIdx);
        if wp[0].GetPresentFlag() {
          iDeltaWeight := int(wp[0].GetWeight() - (1<<wp[0].GetLog2WeightDenom()));
          this.WRITE_SVLC( iDeltaWeight, "delta_luma_weight_lX" );                  // se(v): delta_luma_weight_lX
          this.WRITE_SVLC( wp[0].GetOffset(), "luma_offset_lX" );                       // se(v): luma_offset_lX
        }

        if bChroma {
          if wp[1].GetPresentFlag() {
            for j:=1 ; j<3 ; j++ {
              iDeltaWeight := int(wp[j].GetWeight() - (1<<wp[1].GetLog2WeightDenom()));
              this.WRITE_SVLC( iDeltaWeight, "delta_chroma_weight_lX" );            // se(v): delta_chroma_weight_lX

              pred := int( 128 - ( ( 128*wp[j].GetWeight())>>(wp[j].GetLog2WeightDenom()) ) );
              iDeltaChroma := int(wp[j].GetOffset() - pred);
              this.WRITE_SVLC( iDeltaChroma, "delta_chroma_offset_lX" );            // se(v): delta_chroma_offset_lX
            }
          }
        }
      }
    }
    //assert(uiTotalSignalledWeightFlags<=24);
  }
}

func (this *TEncCavlc) updateContextTables3           ( eSliceType TLibCommon.SliceType, iQp int, bExecuteFinish bool ) { return;   }
func (this *TEncCavlc) updateContextTables2           ( eSliceType TLibCommon.SliceType, iQp int)                       { return;   }

func (this *TEncCavlc) codeScalingList  (scalingList *TLibCommon.TComScalingList) {
  var listId,sizeId uint;
  var scalingListPredModeFlag bool;

/*#if SCALING_LIST_OUTPUT_RESULT
  Int startBit;
  Int startTotalBit;
  startBit = m_pcBitIf.GetNumberOfWrittenBits();
  startTotalBit = m_pcBitIf.GetNumberOfWrittenBits();
#endif*/

    //for each size
    for sizeId = 0; sizeId < TLibCommon.SCALING_LIST_SIZE_NUM; sizeId++ {
      for listId = 0; listId < TLibCommon.G_scalingListNum[sizeId]; listId++ {
/*#if SCALING_LIST_OUTPUT_RESULT
        startBit = m_pcBitIf.GetNumberOfWrittenBits();
#endif*/
        scalingListPredModeFlag = scalingList.CheckPredMode( sizeId, listId );
        this.WRITE_FLAG( uint(TLibCommon.B2U(scalingListPredModeFlag)), "scaling_list_pred_mode_flag" );
        if !scalingListPredModeFlag {// Copy Mode
          this.WRITE_UVLC( uint(int(listId) - int(scalingList.GetRefMatrixId (sizeId,listId))), "scaling_list_pred_matrix_id_delta");
        }else{// DPCM Mode
          this.xCodeScalingList(scalingList, sizeId, listId);
        }
/*#if SCALING_LIST_OUTPUT_RESULT
        printf("Matrix [%d][%d] Bit %d\n",sizeId,listId,m_pcBitIf.GetNumberOfWrittenBits() - startBit);
#endif*/
      }
    }
/*#if SCALING_LIST_OUTPUT_RESULT
  printf("Total Bit %d\n",m_pcBitIf.GetNumberOfWrittenBits()-startTotalBit);
#endif*/
  return;
}

func (this *TEncCavlc) xCodeScalingList ( scalingList *TLibCommon.TComScalingList, sizeId, listId uint) {
  coefNum := TLibCommon.MIN(int(TLibCommon.MAX_MATRIX_COEF_NUM), int(TLibCommon.G_scalingListSize[sizeId])).(int);
  var scan []uint;
  if sizeId == 0 {
  	scan = TLibCommon.G_auiSigLastScan [ TLibCommon.SCAN_DIAG ] [ 1 ][:];
  }else{
  	scan = TLibCommon.G_sigLastScanCG32x32[:];
  }
  nextCoef := TLibCommon.SCALING_LIST_START_VALUE;
  var data int;
  src := scalingList.GetScalingListAddress(sizeId, listId);
    if sizeId > TLibCommon.SCALING_LIST_8x8 {
      this.WRITE_SVLC( scalingList.GetScalingListDC(sizeId,listId) - 8, "scaling_list_dc_coef_minus8");
      nextCoef = scalingList.GetScalingListDC(sizeId,listId);
    }
    for i:=0;i<coefNum;i++ {
      data = src[scan[i]] - nextCoef;
      nextCoef = src[scan[i]];
      if data > 127 {
        data = data - 256;
      }
      if data < -128 {
        data = data + 256;
      }

      this.WRITE_SVLC( data,  "scaling_list_delta_coef");
    }
}

func (this *TEncCavlc) codeDFFlag       ( uiCode uint, pSymbolName string) {
	this.WRITE_FLAG(uiCode, pSymbolName);
}
func (this *TEncCavlc) codeDFSvlc       ( iCode int, pSymbolName string) {
	this.WRITE_SVLC(iCode, pSymbolName);
}

func (this *TEncCavlc) getEncBinIf() TEncBinIf{
	return nil;
}
