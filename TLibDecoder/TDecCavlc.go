package TLibDecoder

import (
	"os"
	"fmt"
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
    m_pTraceFile	*os.File;
    m_pcBitstream 	*TLibCommon.TComInputBitstream
}

func NewSyntaxElementParser() *SyntaxElementParser{
	return &SyntaxElementParser{};
}

func (this *SyntaxElementParser)  SetTraceFile (traceFile *os.File){
	this.m_pTraceFile = traceFile;
}

func (this *SyntaxElementParser)  GetTraceFile () *os.File{
	return this.m_pTraceFile;
}

func (this *SyntaxElementParser)  xTraceVPSHeader (pVPS *TLibCommon.TComVPS){
	if this.GetTraceFile()!=nil { 
  		this.GetTraceFile().WriteString(fmt.Sprintf ("========= Video Parameter Set =============================================\n"));//, pVPS->getVPSId() );
	}
}

func (this *SyntaxElementParser)  xTraceSPSHeader (pSPS *TLibCommon.TComSPS){
	if this.GetTraceFile()!=nil { 
  		this.GetTraceFile().WriteString(fmt.Sprintf ("========= Sequence Parameter Set ==========================================\n"));//, pSPS->getSPSId() );
	}
}

func (this *SyntaxElementParser)  xTracePPSHeader (pPPS *TLibCommon.TComPPS){
	if this.GetTraceFile()!=nil { 
  		this.GetTraceFile().WriteString(fmt.Sprintf ("========= Picture Parameter Set ===========================================\n"));//, pPPS->getPPSId() );
	}
}

func (this *SyntaxElementParser)  xTraceSliceHeader (pSlice *TLibCommon.TComSlice){
	if this.GetTraceFile()!=nil { 
  		//if(g_bSliceTrace)
  		this.GetTraceFile().WriteString(fmt.Sprintf ("========= Slice Parameter Set =============================================\n"));
	}
}
/*
func (this *SyntaxElementParser)  xTraceLCUHeader (pLCU *TComDataCU, traceLevel uint)
{
  if(traceLevel & TRACE_LEVEL)
  fprintf( g_hTrace, "========= LCU Parameter Set ===============================================\n");//, pLCU->getAddr());
}

func (this *SyntaxElementParser)  xTraceCUHeader (pLCU *TComDataCU, traceLevel uint)
{
  if(traceLevel & TRACE_LEVEL)
  fprintf( g_hTrace, "========= CU Parameter Set ================================================\n");//, pCU->getCUPelX(), pCU->getCUPelY());
}

func (this *SyntaxElementParser)  xTracePUHeader (traceLevel uint)
{
  if(traceLevel & TRACE_LEVEL)
    fprintf( g_hTrace, "========= PU Parameter Set ================================================\n");//, pCU->getCUPelX(), pCU->getCUPelY());
}

func (this *SyntaxElementParser)  xTraceTUHeader (traceLevel uint)
{
  if(traceLevel & TRACE_LEVEL)
    fprintf( g_hTrace, "========= TU Parameter Set ================================================\n");//, pCU->getCUPelX(), pCU->getCUPelY());
}

func (this *SyntaxElementParser)  xTraceCoefHeader (traceLevel uint)
{
  if(traceLevel & TRACE_LEVEL)
    fprintf( g_hTrace, "========= Coefficient Parameter Set =======================================\n");//, pCU->getCUPelX(), pCU->getCUPelY());
}

func (this *SyntaxElementParser)  xTraceResiHeader (traceLevel uint)
{
  if(traceLevel & TRACE_LEVEL)
    fprintf( g_hTrace, "========= Residual Parameter Set ==========================================\n");//, pCU->getCUPelX(), pCU->getCUPelY());
}

func (this *SyntaxElementParser)  xTracePredHeader (traceLevel uint)
{
  if(traceLevel & TRACE_LEVEL)
    fprintf( g_hTrace, "========= Prediction Parameter Set ========================================\n");//, pCU->getCUPelX(), pCU->getCUPelY());
}

func (this *SyntaxElementParser)  xTraceRecoHeader (traceLevel uint)
{
  if(traceLevel & TRACE_LEVEL)
    fprintf( g_hTrace, "========= Reconstruction Parameter Set ====================================\n");//, pCU->getCUPelX(), pCU->getCUPelY());
}
*/

func (this *SyntaxElementParser)  xReadCode    ( length uint, val *uint ){
	//assert ( uiLength > 0 );
  	this.m_pcBitstream.Read (length, val);
}
func (this *SyntaxElementParser)  xReadUvlc    ( val *uint ){
  uiVal := uint(0);
  uiCode := uint(0);
  uiLength := uint(0);
  
  this.m_pcBitstream.Read( 1, &uiCode );

  if 0 == uiCode {
    uiLength = 0;

    for ( uiCode & 1 ) ==0 {
      this.m_pcBitstream.Read( 1, &uiCode );
      uiLength++;
    }

    this.m_pcBitstream.Read( uiLength, &uiVal );

    uiVal += (1 << uiLength)-1;
  }

  *val = uiVal;
}
func (this *SyntaxElementParser)  xReadSvlc    ( val *int ){
  uiBits := uint(0);
  this.m_pcBitstream.Read( 1, &uiBits );
  if 0 == uiBits {
    uiLength := uint(0);

    for ( uiBits & 1 )==0 {
      this.m_pcBitstream.Read( 1, &uiBits );
      uiLength++;
    }

    this.m_pcBitstream.Read( uiLength, &uiBits );

    uiBits += (1 << uiLength);
    
    if ( uiBits & 1)==1{
    	*val =-int(uiBits>>1); 
    }else{
    	*val = int(uiBits>>1);
    }
  }else{
    *val = 0;
  }
}
func (this *SyntaxElementParser)  xReadFlag    ( val *uint ){
	this.m_pcBitstream.Read( 1, val );
}
//#if ENC_DEC_TRACE
func (this *SyntaxElementParser)  xReadCodeTr  ( length uint, rValue *uint, pSymbolName string){
  this.xReadCode (length, rValue);
  //if(g_bSliceTrace)
  //{
    //fprintf( g_hTrace, "%8lld  ", g_nSymbolCounter++ );
    if length < 10 {
      this.GetTraceFile().WriteString(fmt.Sprintf ("%-62s u(%d)  : %4d\n", pSymbolName, length, *rValue ));
    }else{
      this.GetTraceFile().WriteString(fmt.Sprintf ("%-62s u(%d) : %4d\n", pSymbolName, length, *rValue )); 
    }
    //fflush ( g_hTrace );
  //}
}
func (this *SyntaxElementParser)  xReadUvlcTr  (              rValue *uint, pSymbolName string){
  this.xReadUvlc (rValue);
  //if(g_bSliceTrace)
  //{
    //fprintf( g_hTrace, "%8lld  ", g_nSymbolCounter++ );
    this.GetTraceFile().WriteString(fmt.Sprintf ("%-62s ue(v) : %4d\n", pSymbolName, *rValue )); 
    //fflush ( g_hTrace );
  //}
}
func (this *SyntaxElementParser)  xReadSvlcTr  (              rValue *int,  pSymbolName string){
  this.xReadSvlc(rValue);
  //if(g_bSliceTrace)
  //{
    //fprintf( g_hTrace, "%8lld  ", g_nSymbolCounter++ );
    this.GetTraceFile().WriteString(fmt.Sprintf ("%-62s se(v) : %4d\n", pSymbolName, *rValue )); 
    //fflush ( g_hTrace );
  //}
}
func (this *SyntaxElementParser)  xReadFlagTr  (              rValue *uint, pSymbolName string){
  this.xReadFlag(rValue);
  //if(g_bSliceTrace)
  //{
    //fprintf( g_hTrace, "%8lld  ", g_nSymbolCounter++ );
    this.GetTraceFile().WriteString(fmt.Sprintf ("%-62s u(1)  : %4d\n", pSymbolName, *rValue )); 
    //fflush ( g_hTrace );
  //}
}
//#endif
//public:
func (this *SyntaxElementParser)  SetBitstream ( p *TLibCommon.TComInputBitstream)   { 
	this.m_pcBitstream = p; 
}
func (this *SyntaxElementParser)  GetBitstream() *TLibCommon.TComInputBitstream{ 
	return this.m_pcBitstream; 
}
//};

//class SEImessages;

/// CAVLC decoder class
type TDecCavlc struct {
    SyntaxElementParser //, public TDecEntropyIf
}


func NewTDecCavlc() *TDecCavlc{
	return &TDecCavlc{}
}

func (this *TDecCavlc)  READ_CODE(length uint, rValue *uint, pSymbolName string){   
	this.xReadCodeTr ( length, rValue, pSymbolName );
}

func (this *TDecCavlc)  READ_UVLC(			   rValue *uint, pSymbolName string){
	this.xReadUvlcTr (         rValue, pSymbolName );
}
func (this *TDecCavlc)  READ_SVLC(             rValue *int,  pSymbolName string){
	this.xReadSvlcTr (         rValue, pSymbolName );
}
func (this *TDecCavlc)  READ_FLAG(             rValue *uint, pSymbolName string){
	this.xReadFlagTr (         rValue, pSymbolName );
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

  var	uiCode uint;
//#if ENC_DEC_TRACE  
  this.xTraceVPSHeader (pcVPS);
//#endif  
  
  this.READ_CODE( 4,  &uiCode,  "video_parameter_set_id" );             
  pcVPS.SetVPSId( int(uiCode) );
  this.READ_FLAG(     &uiCode,  "vps_temporal_id_nesting_flag" );       
  pcVPS.SetTemporalNestingFlag( TLibCommon.U2B(uint8(uiCode)) );
//#if VPS_REARRANGE
  this.READ_CODE( 2,  &uiCode,  "vps_reserved_three_2bits" );           //assert(uiCode == 3);
//#else
//  READ_CODE( 2,  uiCode,  "vps_reserved_zero_2bits" );            assert(uiCode == 0);
//#endif
  this.READ_CODE( 6,  &uiCode,  "vps_reserved_zero_6bits" );            //assert(uiCode == 0);
  this.READ_CODE( 3,  &uiCode,  "vps_max_sub_layers_minus1" );          
  pcVPS.SetMaxTLayers( uiCode + 1 );
//#if VPS_REARRANGE
  this.READ_CODE( 16, &uiCode,  "vps_reserved_ffff_16bits" );           //assert(uiCode == 0xffff);
  this.ParsePTL ( pcVPS.GetPTL(), true, int(pcVPS.GetMaxTLayers())-1);
//#else
//  parsePTL ( pcVPS->getPTL(), true, pcVPS->getMaxTLayers()-1);
//  READ_CODE( 12, uiCode,  "vps_reserved_zero_12bits" );           assert(uiCode == 0);
//#endif
//#if SIGNAL_BITRATE_PICRATE_IN_VPS
  this.ParseBitratePicRateInfo( pcVPS.GetBitratePicrateInfo(), 0, int(pcVPS.GetMaxTLayers()) - 1);
//#endif
//#if HLS_ADD_SUBLAYER_ORDERING_INFO_PRESENT_FLAG
  var subLayerOrderingInfoPresentFlag uint;
  this.READ_FLAG(&subLayerOrderingInfoPresentFlag, "vps_sub_layer_ordering_info_present_flag");
//#endif // HLS_ADD_SUBLAYER_ORDERING_INFO_PRESENT_FLAG 
  for i := uint(0); i <= pcVPS.GetMaxTLayers()-1; i++ {
    this.READ_UVLC( &uiCode,  "vps_max_dec_pic_buffering[i]" );     
    pcVPS.SetMaxDecPicBuffering( uiCode, i );
    this.READ_UVLC( &uiCode,  "vps_num_reorder_pics[i]" );          
    pcVPS.SetNumReorderPics	 ( uiCode, i );
    this.READ_UVLC( &uiCode,  "vps_max_latency_increase[i]" );      
    pcVPS.SetMaxLatencyIncrease( uiCode, i );

//#if HLS_ADD_SUBLAYER_ORDERING_INFO_PRESENT_FLAG
    if subLayerOrderingInfoPresentFlag==0 {
      for i++; i <= pcVPS.GetMaxTLayers()-1; i++ {
        pcVPS.SetMaxDecPicBuffering(pcVPS.GetMaxDecPicBuffering(0), i);
        pcVPS.SetNumReorderPics	   (pcVPS.GetNumReorderPics(0), 	i);
        pcVPS.SetMaxLatencyIncrease(pcVPS.GetMaxLatencyIncrease(0), i);
      }
      break;
    }
//#endif // HLS_ADD_SUBLAYER_ORDERING_INFO_PRESENT_FLAG 
  }

//#if VPS_OPERATING_POINT
  this.READ_UVLC(    &uiCode, "vps_num_hrd_parameters" );               
  pcVPS.SetNumHrdParameters( uiCode );
  this.READ_CODE( 6, &uiCode, "vps_max_nuh_reserved_zero_layer_id" );   
  pcVPS.SetMaxNuhReservedZeroLayerId( uiCode );
  //assert( pcVPS->getNumHrdParameters() < MAX_VPS_NUM_HRD_PARAMETERS_ALLOWED_PLUS1 );
  //assert( pcVPS->getMaxNuhReservedZeroLayerId() < MAX_VPS_NUH_RESERVED_ZERO_LAYER_ID_PLUS1 );
  for opIdx := uint(0); opIdx < pcVPS.GetNumHrdParameters(); opIdx++ {
    if opIdx > 0 {
      // operation_point_layer_id_flag( opIdx )
      for i := uint(0); i <= pcVPS.GetMaxNuhReservedZeroLayerId(); i++ {
        this.READ_FLAG( &uiCode, "op_layer_id_included_flag[opIdx][i]" ); 
        pcVPS.SetOpLayerIdIncludedFlag( TLibCommon.U2B(uint8(uiCode)), opIdx, i );
      }
    }
    // TODO: add hrd_parameters()
  }
//#else
//  READ_UVLC( uiCode,    "vps_num_hrd_parameters" );           assert(uiCode == 0);
  // hrd_parameters
//#endif
  this.READ_FLAG( &uiCode,  "vps_extension_flag" );          //assert(!uiCode);
  //future extensions go here..
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
  var uiCode uint;
  if profilePresentFlag {
    this.ParseProfileTier(rpcPTL.GetGeneralPTL());
  }
  this.READ_CODE( 8, &uiCode, "general_level_idc" );    
  rpcPTL.GetGeneralPTL().SetLevelIdc(int(uiCode));

  for i := 0; i < maxNumSubLayersMinus1; i++ {
//#if CONDITION_SUBLAYERPROFILEPRESENTFLAG
    if profilePresentFlag {
      this.READ_FLAG( &uiCode, "sub_layer_profile_present_flag[i]" ); 
      rpcPTL.SetSubLayerProfilePresentFlag(i, TLibCommon.U2B(uint8(uiCode)));
    }
//#else
//    READ_FLAG( uiCode, "sub_layer_profile_present_flag[i]" ); rpcPTL->setSubLayerProfilePresentFlag(i, uiCode);
//#endif
    this.READ_FLAG( &uiCode, "sub_layer_level_present_flag[i]"   ); 
    rpcPTL.SetSubLayerLevelPresentFlag  (i, TLibCommon.U2B(uint8(uiCode)));
    if profilePresentFlag && rpcPTL.GetSubLayerProfilePresentFlag(i) {
      this.ParseProfileTier(rpcPTL.GetSubLayerPTL(i));
    }
    if rpcPTL.GetSubLayerLevelPresentFlag(i){
      this.READ_CODE( 8, &uiCode, "sub_layer_level_idc[i]" );   
      rpcPTL.GetSubLayerPTL(i).SetLevelIdc(int(uiCode));
    }
  }
}
func (this *TDecCavlc)  ParseProfileTier    ( ptl	*TLibCommon.ProfileTierLevel){
  var uiCode uint;
  this.READ_CODE(2 , &uiCode, "XXX_profile_space[]");   
  ptl.SetProfileSpace(int(uiCode));
  this.READ_FLAG(    &uiCode, "XXX_tier_flag[]"    ); 
  if uiCode!=0 {  
  	ptl.SetTierFlag    (true);
  }else{
  	ptl.SetTierFlag    (false);
  }
  this.READ_CODE(5 , &uiCode, "XXX_profile_idc[]"  );   
  ptl.SetProfileIdc  (int(uiCode));
  for j := 0; j < 32; j++ {
    this.READ_FLAG(  &uiCode, "XXX_profile_compatibility_flag[][j]"); 
    if uiCode!=0 {    
    	ptl.SetProfileCompatibilityFlag(j, true);
    }else{
    	ptl.SetProfileCompatibilityFlag(j, false);
    }
  }
  this.READ_CODE(16, &uiCode, "XXX_reserved_zero_16bits[]");  //assert( uiCode == 0 );  
	
}
//#if SIGNAL_BITRATE_PICRATE_IN_VPS
func (this *TDecCavlc)  ParseBitratePicRateInfo(info *TLibCommon.TComBitRatePicRateInfo,  tempLevelLow,  tempLevelHigh int){
  var uiCode uint;
  for i := tempLevelLow; i <= tempLevelHigh; i++ {
    this.READ_FLAG( &uiCode, "bit_rate_info_present_flag[i]" ); 
    if uiCode!=0{
    	info.SetBitRateInfoPresentFlag(i, true);
    }else{
    	info.SetBitRateInfoPresentFlag(i, false);
    }
    this.READ_FLAG( &uiCode, "pic_rate_info_present_flag[i]" ); 
    if uiCode!=0{
    	info.SetPicRateInfoPresentFlag(i, true);
    }else{
    	info.SetPicRateInfoPresentFlag(i, false);
    }
    if info.GetBitRateInfoPresentFlag(i){
      this.READ_CODE( 16, &uiCode, "avg_bit_rate[i]" ); 
      info.SetAvgBitRate(i, int(uiCode));
      this.READ_CODE( 16, &uiCode, "max_bit_rate[i]" ); 
      info.SetMaxBitRate(i, int(uiCode));
    }
    if info.GetPicRateInfoPresentFlag(i) {
      this.READ_CODE(  2, &uiCode,  "constant_pic_rate_idc[i]" ); 
      info.SetConstantPicRateIdc(i, int(uiCode));
      this.READ_CODE( 16, &uiCode,  "avg_pic_rate[i]"          ); 
      info.SetAvgPicRate(i, int(uiCode));
    }
  }
}
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
  
