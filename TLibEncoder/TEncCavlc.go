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
  m_pcBitIf TLibCommon.TComBitIf;
}

func NewSyntaxElementWriter() *SyntaxElementWriter{
	return &SyntaxElementWriter{};
}

func (this *SyntaxElementWriter)  setBitstream          ( p TLibCommon.TComBitIf)  { this.m_pcBitIf = p;  }

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
/*

func (this *TEncCavlc)  xWritePCMAlignZero    () {m_pcBitIf->writeAlignZero();}
func (this *TEncCavlc)  xWriteEpExGolomb      ( UInt uiSymbol, UInt uiCount ){
  while( uiSymbol >= (UInt)(1<<uiCount) )
  {
    xWriteFlag( 1 );
    uiSymbol -= 1<<uiCount;
    uiCount  ++;
  }
  xWriteFlag( 0 );
  while( uiCount-- )
  {
    xWriteFlag( (uiSymbol>>uiCount) & 1 );
  }
  return;
}
func (this *TEncCavlc)  xWriteExGolombLevel   ( UInt uiSymbol ){
  if( uiSymbol )
  {
    xWriteFlag( 1 );
    UInt uiCount = 0;
    Bool bNoExGo = (uiSymbol < 13);
    
    while( --uiSymbol && ++uiCount < 13 )
    {
      xWriteFlag( 1 );
    }
    if( bNoExGo )
    {
      xWriteFlag( 0 );
    }
    else
    {
      xWriteEpExGolomb( uiSymbol, 0 );
    }
  }
  else
  {
    xWriteFlag( 0 );
  }
  return;
}
func (this *TEncCavlc)  xWriteUnaryMaxSymbol  ( UInt uiSymbol, UInt uiMaxSymbol ){
  if (uiMaxSymbol == 0)
  {
    return;
  }
  xWriteFlag( uiSymbol ? 1 : 0 );
  if ( uiSymbol == 0 )
  {
    return;
  }
  
  Bool bCodeLast = ( uiMaxSymbol > uiSymbol );
  
  while( --uiSymbol )
  {
    xWriteFlag( 1 );
  }
  if( bCodeLast )
  {
    xWriteFlag( 0 );
  }
  return;
}
//#if SPS_INTER_REF_SET_PRED
func (this *TEncCavlc) codeShortTermRefPicSet              ( TComSPS* pcSPS, TComReferencePictureSet* pcRPS, Bool calledFromSliceHeader, Int idx ){
//#else
//func (this *TEncCavlc) codeShortTermRefPicSet              ( TComSPS* pcSPS, TComReferencePictureSet* pcRPS, Bool calledFromSliceHeader );
//#endif
#if PRINT_RPS_INFO
  Int lastBits = getNumberOfWrittenBits();
#endif
#if SPS_INTER_REF_SET_PRED
  if (idx > 0)
  {
#endif
  WRITE_FLAG( rps->getInterRPSPrediction(), "inter_ref_pic_set_prediction_flag" ); // inter_RPS_prediction_flag
#if SPS_INTER_REF_SET_PRED
  }
#endif
  if (rps->getInterRPSPrediction()) 
  {
    Int deltaRPS = rps->getDeltaRPS();
    if(calledFromSliceHeader)
    {
      WRITE_UVLC( rps->getDeltaRIdxMinus1(), "delta_idx_minus1" ); // delta index of the Reference Picture Set used for prediction minus 1
    }

    WRITE_CODE( (deltaRPS >=0 ? 0: 1), 1, "delta_rps_sign" ); //delta_rps_sign
    WRITE_UVLC( abs(deltaRPS) - 1, "abs_delta_rps_minus1"); // absolute delta RPS minus 1

    for(Int j=0; j < rps->getNumRefIdc(); j++)
    {
      Int refIdc = rps->getRefIdc(j);
      WRITE_CODE( (refIdc==1? 1: 0), 1, "used_by_curr_pic_flag" ); //first bit is "1" if Idc is 1 
      if (refIdc != 1) 
      {
        WRITE_CODE( refIdc>>1, 1, "use_delta_flag" ); //second bit is "1" if Idc is 2, "0" otherwise.
      }
    }
  }
  else
  {
    WRITE_UVLC( rps->getNumberOfNegativePictures(), "num_negative_pics" );
    WRITE_UVLC( rps->getNumberOfPositivePictures(), "num_positive_pics" );
    Int prev = 0;
    for(Int j=0 ; j < rps->getNumberOfNegativePictures(); j++)
    {
      WRITE_UVLC( prev-rps->getDeltaPOC(j)-1, "delta_poc_s0_minus1" );
      prev = rps->getDeltaPOC(j);
      WRITE_FLAG( rps->getUsed(j), "used_by_curr_pic_s0_flag"); 
    }
    prev = 0;
    for(Int j=rps->getNumberOfNegativePictures(); j < rps->getNumberOfNegativePictures()+rps->getNumberOfPositivePictures(); j++)
    {
      WRITE_UVLC( rps->getDeltaPOC(j)-prev-1, "delta_poc_s1_minus1" );
      prev = rps->getDeltaPOC(j);
      WRITE_FLAG( rps->getUsed(j), "used_by_curr_pic_s1_flag" ); 
    }
  }

#if PRINT_RPS_INFO
  printf("irps=%d (%2d bits) ", rps->getInterRPSPrediction(), getNumberOfWrittenBits() - lastBits);
  rps->printDeltaPOC();
#endif
}
func (this *TEncCavlc)  findMatchingLTRP ( TComSlice* pcSlice, UInt *ltrpsIndex, Int ltrpPOC, Bool usedFlag ) bool{
  // Bool state = true, state2 = false;
  Int lsb = ltrpPOC % (1<<pcSlice->getSPS()->getBitsForPOC());
  for (Int k = 0; k < pcSlice->getSPS()->getNumLongTermRefPicSPS(); k++)
  {
    if ( (lsb == pcSlice->getSPS()->getLtRefPicPocLsbSps(k)) && (usedFlag == pcSlice->getSPS()->getUsedByCurrPicLtSPSFlag(k)) )
    {
      *ltrpsIndex = k;
      return true;
    }
  }
  return false;
} 
func (this *TEncCavlc)  resetEntropy          () {}
func (this *TEncCavlc)  determineCabacInitIdx () {}
func (this *TEncCavlc)  setBitstream          ( TComBitIf* p )  { m_pcBitIf = p;  }
func (this *TEncCavlc)  setSlice              ( TComSlice* p )  { m_pcSlice = p;  }
func (this *TEncCavlc)  resetBits             ()                { m_pcBitIf->resetBits(); }
func (this *TEncCavlc)  resetCoeffCost        ()                { m_uiCoeffCost = 0;  }
func (this *TEncCavlc)  getNumberOfWrittenBits() uint               { return  m_pcBitIf->getNumberOfWrittenBits();  }
func (this *TEncCavlc)  getCoeffCost          () uint               { return  m_uiCoeffCost;  }
func (this *TEncCavlc)  codeVPS                 ( TComVPS* pcVPS ){
  WRITE_CODE( pcVPS->getVPSId(),                    4,        "video_parameter_set_id" );
  WRITE_FLAG( pcVPS->getTemporalNestingFlag(),                "vps_temporal_id_nesting_flag" );
#if VPS_REARRANGE
  WRITE_CODE( 3,                                    2,        "vps_reserved_three_2bits" );
#else
  WRITE_CODE( 0,                                    2,        "vps_reserved_zero_2bits" );
#endif
  WRITE_CODE( 0,                                    6,        "vps_reserved_zero_6bits" );
  WRITE_CODE( pcVPS->getMaxTLayers() - 1,           3,        "vps_max_sub_layers_minus1" );
#if VPS_REARRANGE
  WRITE_CODE( 0xffff,                              16,        "vps_reserved_ffff_16bits" );
  codePTL( pcVPS->getPTL(), true, pcVPS->getMaxTLayers() - 1 );
#else
  codePTL( pcVPS->getPTL(), true, pcVPS->getMaxTLayers() - 1 );
  WRITE_CODE( 0,                                   12,        "vps_reserved_zero_12bits" );
#endif
#if SIGNAL_BITRATE_PICRATE_IN_VPS
  codeBitratePicRateInfo(pcVPS->getBitratePicrateInfo(), 0, pcVPS->getMaxTLayers() - 1);
#endif  
#if HLS_ADD_SUBLAYER_ORDERING_INFO_PRESENT_FLAG
  const Bool subLayerOrderingInfoPresentFlag = 1;
  WRITE_FLAG(subLayerOrderingInfoPresentFlag,              "vps_sub_layer_ordering_info_present_flag");
#endif // HLS_ADD_SUBLAYER_ORDERING_INFO_PRESENT_FLAG
  for(UInt i=0; i <= pcVPS->getMaxTLayers()-1; i++)
  {
    WRITE_UVLC( pcVPS->getMaxDecPicBuffering(i),           "vps_max_dec_pic_buffering[i]" );
    WRITE_UVLC( pcVPS->getNumReorderPics(i),               "vps_num_reorder_pics[i]" );
    WRITE_UVLC( pcVPS->getMaxLatencyIncrease(i),           "vps_max_latency_increase[i]" );
#if HLS_ADD_SUBLAYER_ORDERING_INFO_PRESENT_FLAG
    if (!subLayerOrderingInfoPresentFlag)
    {
      break;
    }
#endif // HLS_ADD_SUBLAYER_ORDERING_INFO_PRESENT_FLAG
  }

#if VPS_OPERATING_POINT
  assert( pcVPS->getNumHrdParameters() <= MAX_VPS_NUM_HRD_PARAMETERS );
  assert( pcVPS->getMaxNuhReservedZeroLayerId() < MAX_VPS_NUH_RESERVED_ZERO_LAYER_ID_PLUS1 );
  WRITE_UVLC( pcVPS->getNumHrdParameters(),                 "vps_num_hrd_parameters" );
  WRITE_CODE( pcVPS->getMaxNuhReservedZeroLayerId(), 6,     "vps_max_nuh_reserved_zero_layer_id" );
  for( UInt opIdx = 0; opIdx < pcVPS->getNumHrdParameters(); opIdx++ )
  {
    if( opIdx > 0 )
    {
      // operation_point_layer_id_flag( opIdx )
      for( UInt i = 0; i <= pcVPS->getMaxNuhReservedZeroLayerId(); i++ )
      {
        WRITE_FLAG( pcVPS->getOpLayerIdIncludedFlag( opIdx, i ), "op_layer_id_included_flag[opIdx][i]" );
      }
    }
    // TODO: add hrd_parameters()
  }
#else
  WRITE_UVLC( 0,                                           "vps_num_hrd_parameters" );
  // hrd_parameters
#endif
  WRITE_FLAG( 0,                     "vps_extension_flag" );
  
  //future extensions here..
  
  return;
}
func (this *TEncCavlc)  codeVUI                 ( TComVUI *pcVUI, TComSPS* pcSPS );
func (this *TEncCavlc)  codeSPS                 ( TComSPS* pcSPS );
func (this *TEncCavlc)  codePPS                 ( TComPPS* pcPPS );
func (this *TEncCavlc)  codeSliceHeader         ( TComSlice* pcSlice );
func (this *TEncCavlc)  codePTL                 ( TComPTL* pcPTL, Bool profilePresentFlag, Int maxNumSubLayersMinus1);
func (this *TEncCavlc)  codeProfileTier         ( ProfileTierLevel* ptl );
//#if SIGNAL_BITRATE_PICRATE_IN_VPS
func (this *TEncCavlc) codeBitratePicRateInfo(TComBitRatePicRateInfo *info, Int tempLevelLow, Int tempLevelHigh);
//#endif
func (this *TEncCavlc)  codeTilesWPPEntryPoint( TComSlice* pSlice );
func (this *TEncCavlc)  codeTerminatingBit      ( UInt uilsLast );
func (this *TEncCavlc)  codeSliceFinish         ();
func (this *TEncCavlc)  encodeStart             () {}
  
func (this *TEncCavlc) codeMVPIdx ( TComDataCU* pcCU, UInt uiAbsPartIdx, RefPicList eRefList );
func (this *TEncCavlc) codeSAOSign       ( UInt code   ) { printf("Not supported\n"); assert (0); }
func (this *TEncCavlc) codeSaoMaxUvlc    ( UInt   code, UInt maxSymbol ){printf("Not supported\n"); assert (0);}
func (this *TEncCavlc) codeSaoMerge  ( UInt uiCode ){printf("Not supported\n"); assert (0);}
func (this *TEncCavlc) codeSaoTypeIdx    ( UInt uiCode ){printf("Not supported\n"); assert (0);}
func (this *TEncCavlc) codeSaoUflc       ( UInt uiLength, UInt   uiCode ){ assert(uiCode < 32); printf("Not supported\n"); assert (0);}

func (this *TEncCavlc) codeCUTransquantBypassFlag( TComDataCU* pcCU, UInt uiAbsPartIdx );
func (this *TEncCavlc) codeSkipFlag      ( TComDataCU* pcCU, UInt uiAbsPartIdx );
func (this *TEncCavlc) codeMergeFlag     ( TComDataCU* pcCU, UInt uiAbsPartIdx );
func (this *TEncCavlc) codeMergeIndex    ( TComDataCU* pcCU, UInt uiAbsPartIdx );
 
func (this *TEncCavlc) codeInterModeFlag( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiDepth, UInt uiEncMode );
func (this *TEncCavlc) codeSplitFlag     ( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiDepth );
  
func (this *TEncCavlc) codePartSize      ( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiDepth );
func (this *TEncCavlc) codePredMode      ( TComDataCU* pcCU, UInt uiAbsPartIdx );
  
//#if !REMOVE_BURST_IPCM
//func (this *TEncCavlc) codeIPCMInfo      ( TComDataCU* pcCU, UInt uiAbsPartIdx, Int numIPCM, Bool firstIPCMFlag);
//#else
func (this *TEncCavlc) codeIPCMInfo      ( TComDataCU* pcCU, UInt uiAbsPartIdx );
//#endif

func (this *TEncCavlc) codeTransformSubdivFlag( UInt uiSymbol, UInt uiCtx );
func (this *TEncCavlc) codeQtCbf         ( TComDataCU* pcCU, UInt uiAbsPartIdx, TextType eType, UInt uiTrDepth );
func (this *TEncCavlc) codeQtRootCbf     ( TComDataCU* pcCU, UInt uiAbsPartIdx );
func (this *TEncCavlc) codeQtCbfZero     ( TComDataCU* pcCU, UInt uiAbsPartIdx, TextType eType, UInt uiTrDepth );
func (this *TEncCavlc) codeQtRootCbfZero ( TComDataCU* pcCU, UInt uiAbsPartIdx );
func (this *TEncCavlc) codeIntraDirLumaAng( TComDataCU* pcCU, UInt absPartIdx, Bool isMultiple);
func (this *TEncCavlc) codeIntraDirChroma( TComDataCU* pcCU, UInt uiAbsPartIdx );
func (this *TEncCavlc) codeInterDir      ( TComDataCU* pcCU, UInt uiAbsPartIdx );
func (this *TEncCavlc) codeRefFrmIdx     ( TComDataCU* pcCU, UInt uiAbsPartIdx, RefPicList eRefList );
func (this *TEncCavlc) codeMvd           ( TComDataCU* pcCU, UInt uiAbsPartIdx, RefPicList eRefList );
  
func (this *TEncCavlc) codeDeltaQP       ( TComDataCU* pcCU, UInt uiAbsPartIdx );
  
func (this *TEncCavlc) codeCoeffNxN      ( TComDataCU* pcCU, TCoeff* pcCoef, UInt uiAbsPartIdx, UInt uiWidth, UInt uiHeight, UInt uiDepth, TextType eTType );
func (this *TEncCavlc) codeTransformSkipFlags ( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt width, UInt height, UInt uiDepth, TextType eTType );

func (this *TEncCavlc) estBit               (estBitsSbacStruct* pcEstBitsSbac, Int width, Int height, TextType eTType);
  
func (this *TEncCavlc) xCodePredWeightTable          ( TComSlice* pcSlice );
func (this *TEncCavlc) updateContextTables           ( SliceType eSliceType, Int iQp, Bool bExecuteFinish=true ) { return;   }
func (this *TEncCavlc) updateContextTables           ( SliceType eSliceType, Int iQp  )                          { return;   }

func (this *TEncCavlc) codeScalingList  ( TComScalingList* scalingList );
func (this *TEncCavlc) xCodeScalingList ( TComScalingList* scalingList, UInt sizeId, UInt listId);
func (this *TEncCavlc) codeDFFlag       ( UInt uiCode, const Char *pSymbolName );
func (this *TEncCavlc) codeDFSvlc       ( Int   iCode, const Char *pSymbolName );
*/