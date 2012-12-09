package TLibDecoder

import (
    "gohm/TLibCommon"
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
    m_pcBitstream *TLibCommon.TComInputBitstream
}

/*
  SyntaxElementParser()
  : m_pcBitstream (NULL)
  {};
  virtual ~SyntaxElementParser() {};

  Void  xReadCode    ( UInt   length, UInt& val );
  Void  xReadUvlc    ( UInt&  val );
  Void  xReadSvlc    ( Int&   val );
  Void  xReadFlag    ( UInt&  val );
#if ENC_DEC_TRACE
  Void  xReadCodeTr  (UInt  length, UInt& rValue, const Char *pSymbolName);
  Void  xReadUvlcTr  (              UInt& rValue, const Char *pSymbolName);
  Void  xReadSvlcTr  (               Int& rValue, const Char *pSymbolName);
  Void  xReadFlagTr  (              UInt& rValue, const Char *pSymbolName);
#endif
public:
  Void  setBitstream ( TComInputBitstream* p )   { m_pcBitstream = p; }
  TComInputBitstream* getBitstream() { return m_pcBitstream; }
};*/

//class SEImessages;

/// CAVLC decoder class
type TDecCavlc struct {
    SyntaxElementParser //, public TDecEntropyIf
}


func NewTDecCavlc() *TDecCavlc{
	return &TDecCavlc{}
}

//protected:
func (this *TDecCavlc)  xReadEpExGolomb       ( ruiSymbol *uint, uiCount uint){
}
func (this *TDecCavlc)  xReadExGolombLevel    ( ruiSymbol *uint){
}
func (this *TDecCavlc)  xReadUnaryMaxSymbol   ( ruiSymbol, uiMaxSymbol uint){
}

func (this *TDecCavlc)  xReadPCMAlignZero     (){
}

func (this *TDecCavlc)  xGetBit             () uint{
	return 0;
}

func (this *TDecCavlc)  ParseShortTermRefPicSet            (pcSPS *TLibCommon.TComSPS, pcRPS *TLibCommon.TComReferencePictureSet, idx int){
}


//public:

  /// rest entropy coder by intial QP and IDC in CABAC
func (this *TDecCavlc)  ResetEntropy        ( pcSlice *TLibCommon.TComSlice )     { 
	//assert(0); 
};
func (this *TDecCavlc)  SetBitstream        ( p *TLibCommon.TComInputBitstream)   { 
	this.m_pcBitstream = p; 
}
func (this *TDecCavlc)  ParseTransformSubdivFlag( ruiSubdivFlag *uint, uiLog2TransformBlockSize uint ){
}
func (this *TDecCavlc)  ParseQtCbf          ( pcCU *TLibCommon.TComDataCU,  uiAbsPartIdx uint,  eType TLibCommon.TextType,  uiTrDepth,  uiDepth uint ){
}
func (this *TDecCavlc)  ParseQtRootCbf      ( pcCU *TLibCommon.TComDataCU,  uiAbsPartIdx,  uiDepth uint, uiQtRootCbf *uint ){
}
func (this *TDecCavlc)  ParseVPS            ( pcVPS *TLibCommon.TComVPS){
/*
  UInt  uiCode;
#if ENC_DEC_TRACE  
  xTraceVPSHeader (pcVPS);
#endif  
  
  READ_CODE( 4,  uiCode,  "video_parameter_set_id" );             pcVPS->setVPSId( uiCode );
  READ_FLAG(     uiCode,  "vps_temporal_id_nesting_flag" );       pcVPS->setTemporalNestingFlag( uiCode ? true:false );
#if VPS_REARRANGE
  READ_CODE( 2,  uiCode,  "vps_reserved_three_2bits" );           assert(uiCode == 3);
#else
  READ_CODE( 2,  uiCode,  "vps_reserved_zero_2bits" );            assert(uiCode == 0);
#endif
  READ_CODE( 6,  uiCode,  "vps_reserved_zero_6bits" );            assert(uiCode == 0);
  READ_CODE( 3,  uiCode,  "vps_max_sub_layers_minus1" );          pcVPS->setMaxTLayers( uiCode + 1 );
#if VPS_REARRANGE
  READ_CODE( 16, uiCode,  "vps_reserved_ffff_16bits" );           assert(uiCode == 0xffff);
  parsePTL ( pcVPS->getPTL(), true, pcVPS->getMaxTLayers()-1);
#else
  parsePTL ( pcVPS->getPTL(), true, pcVPS->getMaxTLayers()-1);
  READ_CODE( 12, uiCode,  "vps_reserved_zero_12bits" );           assert(uiCode == 0);
#endif
#if SIGNAL_BITRATE_PICRATE_IN_VPS
  parseBitratePicRateInfo( pcVPS->getBitratePicrateInfo(), 0, pcVPS->getMaxTLayers() - 1);
#endif
#if HLS_ADD_SUBLAYER_ORDERING_INFO_PRESENT_FLAG
  UInt subLayerOrderingInfoPresentFlag;
  READ_FLAG(subLayerOrderingInfoPresentFlag, "vps_sub_layer_ordering_info_present_flag");
#endif // HLS_ADD_SUBLAYER_ORDERING_INFO_PRESENT_FLAG 
  for(UInt i = 0; i <= pcVPS->getMaxTLayers()-1; i++)
  {
    READ_UVLC( uiCode,  "vps_max_dec_pic_buffering[i]" );     pcVPS->setMaxDecPicBuffering( uiCode, i );
    READ_UVLC( uiCode,  "vps_num_reorder_pics[i]" );          pcVPS->setNumReorderPics( uiCode, i );
    READ_UVLC( uiCode,  "vps_max_latency_increase[i]" );      pcVPS->setMaxLatencyIncrease( uiCode, i );

#if HLS_ADD_SUBLAYER_ORDERING_INFO_PRESENT_FLAG
    if (!subLayerOrderingInfoPresentFlag)
    {
      for (i++; i <= pcVPS->getMaxTLayers()-1; i++)
      {
        pcVPS->setMaxDecPicBuffering(pcVPS->getMaxDecPicBuffering(0), i);
        pcVPS->setNumReorderPics(pcVPS->getNumReorderPics(0), i);
        pcVPS->setMaxLatencyIncrease(pcVPS->getMaxLatencyIncrease(0), i);
      }
      break;
    }
#endif // HLS_ADD_SUBLAYER_ORDERING_INFO_PRESENT_FLAG 
  }

#if VPS_OPERATING_POINT
  READ_UVLC(    uiCode, "vps_num_hrd_parameters" );               pcVPS->setNumHrdParameters( uiCode );
  READ_CODE( 6, uiCode, "vps_max_nuh_reserved_zero_layer_id" );   pcVPS->setMaxNuhReservedZeroLayerId( uiCode );
  assert( pcVPS->getNumHrdParameters() < MAX_VPS_NUM_HRD_PARAMETERS_ALLOWED_PLUS1 );
  assert( pcVPS->getMaxNuhReservedZeroLayerId() < MAX_VPS_NUH_RESERVED_ZERO_LAYER_ID_PLUS1 );
  for( UInt opIdx = 0; opIdx < pcVPS->getNumHrdParameters(); opIdx++ )
  {
    if( opIdx > 0 )
    {
      // operation_point_layer_id_flag( opIdx )
      for( UInt i = 0; i <= pcVPS->getMaxNuhReservedZeroLayerId(); i++ )
      {
        READ_FLAG( uiCode, "op_layer_id_included_flag[opIdx][i]" ); pcVPS->setOpLayerIdIncludedFlag( uiCode, opIdx, i );
      }
    }
    // TODO: add hrd_parameters()
  }
#else
  READ_UVLC( uiCode,    "vps_num_hrd_parameters" );           assert(uiCode == 0);
  // hrd_parameters
#endif
  READ_FLAG( uiCode,  "vps_extension_flag" );          assert(!uiCode);
  //future extensions go here..
  
  return;
  */
}
func (this *TDecCavlc)  ParseSPS            ( pcSPS *TLibCommon.TComSPS){
}
func (this *TDecCavlc)  ParsePPS            ( pcPPS	*TLibCommon.TComPPS){
}
func (this *TDecCavlc)  ParseVUI            ( pcVUI *TLibCommon.TComVUI, pcSPS *TLibCommon.TComSPS){
}
func (this *TDecCavlc)  ParseSEI			( sei   *TLibCommon.SEImessages){
}
func (this *TDecCavlc)  ParsePTL            ( rpcPTL *TLibCommon.TComPTL, profilePresentFlag bool, maxNumSubLayersMinus1 int ){
}
func (this *TDecCavlc)  ParseProfileTier    ( ptl	*TLibCommon.ProfileTierLevel){
}
//#if SIGNAL_BITRATE_PICRATE_IN_VPS
//func (this *TDecCavlc)  ParseBitratePicRateInfo(info *TLibCommon.TComBitRatePicRateInfo,  tempLevelLow,  tempLevelHigh int){
//}
//#endif
func (this *TDecCavlc)  ParseSliceHeader    ( rpcSlice *TLibCommon.TComSlice, parameterSetManager *TLibCommon.ParameterSetManager){//Decoder
}
func (this *TDecCavlc)  ParseTerminatingBit ( ruiBit *uint){
}

func (this *TDecCavlc)  ParseMVPIdx         ( riMVPIdx *int){
}

func (this *TDecCavlc)  ParseSkipFlag        ( pcCU *TLibCommon.TComDataCU,  uiAbsPartIdx,  uiDepth uint){
}
func (this *TDecCavlc)  ParseCUTransquantBypassFlag( pcCU *TLibCommon.TComDataCU,  uiAbsPartIdx,  uiDepth uint ){
}
func (this *TDecCavlc)  ParseMergeFlag       ( pcCU *TLibCommon.TComDataCU,  uiAbsPartIdx,  uiDepth, uiPUIdx uint ){
}
func (this *TDecCavlc)  ParseMergeIndex      ( pcCU *TLibCommon.TComDataCU,  ruiMergeIndex *uint,  uiAbsPartIdx,  uiDepth uint ){
}
func (this *TDecCavlc)  ParseSplitFlag       ( pcCU *TLibCommon.TComDataCU,  uiAbsPartIdx,  uiDepth uint ){
}
func (this *TDecCavlc)  ParsePartSize        ( pcCU *TLibCommon.TComDataCU,  uiAbsPartIdx,  uiDepth uint ){
}
func (this *TDecCavlc)  ParsePredMode        ( pcCU *TLibCommon.TComDataCU,  uiAbsPartIdx,  uiDepth uint ){
}

func (this *TDecCavlc)  ParseIntraDirLumaAng ( pcCU *TLibCommon.TComDataCU,  uiAbsPartIdx,  uiDepth uint ){
}

func (this *TDecCavlc)  ParseIntraDirChroma  ( pcCU *TLibCommon.TComDataCU,  uiAbsPartIdx,  uiDepth uint ){
}

func (this *TDecCavlc)  ParseInterDir        ( pcCU *TLibCommon.TComDataCU, ruiInterDir *uint,  uiAbsPartIdx,  uiDepth uint){
}
func (this *TDecCavlc)  ParseRefFrmIdx       ( pcCU *TLibCommon.TComDataCU, riRefFrmIdx *int,   uiAbsPartIdx,  uiDepth uint,  eRefList TLibCommon.RefPicList){
}
func (this *TDecCavlc)  ParseMvd             ( pcCU *TLibCommon.TComDataCU, uiAbsPartAddr, uiPartIdx, uiDepth uint,  eRefList TLibCommon.RefPicList){
}

func (this *TDecCavlc)  ParseDeltaQP         ( pcCU *TLibCommon.TComDataCU,  uiAbsPartIdx,  uiDepth uint){
}
func (this *TDecCavlc)  ParseCoeffNxN        ( pcCU *TLibCommon.TComDataCU, pcCoef *TLibCommon.TCoeff,  uiAbsPartIdx,  uiWidth,  uiHeight,  uiDepth uint,  eTType TLibCommon.TextType){
}
func (this *TDecCavlc)  ParseTransformSkipFlags ( pcCU *TLibCommon.TComDataCU,  uiAbsPartIdx,  width,  height,  uiDepth uint,  eTType TLibCommon.TextType){
}

func (this *TDecCavlc)  ParseIPCMInfo        ( pcCU *TLibCommon.TComDataCU,  uiAbsPartIdx,  uiDepth uint){
}

func (this *TDecCavlc)  UpdateContextTables  (  eSliceType TLibCommon.SliceType,  iQp int) { 
	return; 
}

func (this *TDecCavlc)  xParsePredWeightTable ( pcSlice *TLibCommon.TComSlice){
}
func (this *TDecCavlc)  ParseScalingList               ( scalingList *TLibCommon.TComScalingList){
}
func (this *TDecCavlc)  xDecodeScalingList    ( scalingList *TLibCommon.TComScalingList,  sizeId,  listId uint){
}
//protected:
func (this *TDecCavlc)  xMoreRbspData() bool{
	return true;
}
  
