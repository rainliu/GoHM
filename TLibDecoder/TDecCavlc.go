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

func (this *TDecCavlc)  ParseShortTermRefPicSet            (sps *TLibCommon.TComSPS, rps *TLibCommon.TComReferencePictureSet, idx int){
  var code, interRPSPred uint;
  
//#if SPS_INTER_REF_SET_PRED
  if idx > 0{
//#endif 
  	this.READ_FLAG(&interRPSPred, "inter_ref_pic_set_prediction_flag");  
  	rps.SetInterRPSPrediction(interRPSPred!=0);
//#if SPS_INTER_REF_SET_PRED
  }else{
    interRPSPred = 0;
    rps.SetInterRPSPrediction(false);
  }
//#endif

  if interRPSPred!=0 {
    var bit uint;
    if idx == sps.GetRPSList().GetNumberOfReferencePictureSets(){
      this.READ_UVLC(&code, "delta_idx_minus1" ); // delta index of the Reference Picture Set used for prediction minus 1
    }else{
      code = 0;
    }
    //assert(code <= idx-1); // delta_idx_minus1 shall not be larger than idx-1, otherwise we will predict from a negative row position that does not exist. When idx equals 0 there is no legal value and interRPSPred must be zero. See J0185-r2
    rIdx :=  idx - 1 - int(code);
    //assert (rIdx <= idx-1 && rIdx >= 0); // Made assert tighter; if rIdx = idx then prediction is done from itself. rIdx must belong to range 0, idx-1, inclusive, see J0185-r2
    rpsRef := sps.GetRPSList().GetReferencePictureSet(rIdx);
    k := 0;
    k0 := 0;
    k1 := 0;
    this.READ_CODE(1, &bit, "delta_rps_sign"); // delta_RPS_sign
    this.READ_UVLC(&code, "abs_delta_rps_minus1");  // absolute delta RPS minus 1
    deltaRPS := (1 - (bit<<1)) * (code + 1); // delta_RPS
    for j:=0 ; j <= rpsRef.GetNumberOfPictures(); j++ {
      this.READ_CODE(1, &bit, "used_by_curr_pic_flag" ); //first bit is "1" if Idc is 1 
      refIdc := bit;
      if refIdc == 0 {
        this.READ_CODE(1, &bit, "use_delta_flag" ); //second bit is "1" if Idc is 2, "0" otherwise.
        refIdc = bit<<1; //second bit is "1" if refIdc is 2, "0" if refIdc = 0.
      }
      
      if refIdc == 1 || refIdc == 2 {
      	var deltaPOC int
     	if j < rpsRef.GetNumberOfPictures(){
        	deltaPOC = int(deltaRPS) + rpsRef.GetDeltaPOC(j);
        }else{
        	deltaPOC = int(deltaRPS);
        }
        rps.SetDeltaPOC(k, deltaPOC);
        rps.SetUsed(k, (refIdc == 1));

        if deltaPOC < 0 {
          k0++;
        }else{
          k1++;
        }
        k++;
      }  
      rps.SetRefIdc(j, int(refIdc));  
    }
    rps.SetNumRefIdc(rpsRef.GetNumberOfPictures()+1);  
    rps.SetNumberOfPictures(k);
    rps.SetNumberOfNegativePictures(k0);
    rps.SetNumberOfPositivePictures(k1);
    rps.SortDeltaPOC();
  }else{
    this.READ_UVLC(&code, "num_negative_pics");           
    rps.SetNumberOfNegativePictures(int(code));
    this.READ_UVLC(&code, "num_positive_pics");           
    rps.SetNumberOfPositivePictures(int(code));
    prev := 0;
    var poc int;
    for j:=0 ; j < rps.GetNumberOfNegativePictures(); j++ {
      this.READ_UVLC(&code, "delta_poc_s0_minus1");
      poc = prev-int(code)-1;
      prev = poc;
      rps.SetDeltaPOC(j,poc);
      this.READ_FLAG(&code, "used_by_curr_pic_s0_flag");  
      rps.SetUsed(j,code!=0);
    }
    prev = 0;
    for j:=rps.GetNumberOfNegativePictures(); j < rps.GetNumberOfNegativePictures()+rps.GetNumberOfPositivePictures(); j++ {
      this.READ_UVLC(&code, "delta_poc_s1_minus1");
      poc = prev+int(code)+1;
      prev = poc;
      rps.SetDeltaPOC(j,poc);
      this.READ_FLAG(&code, "used_by_curr_pic_s1_flag");  
      rps.SetUsed(j,code!=0);
    }
    rps.SetNumberOfPictures(rps.GetNumberOfNegativePictures()+rps.GetNumberOfPositivePictures());
  }
//#if PRINT_RPS_INFO
//  rps->printDeltaPOC();
//#endif
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
//#if ENC_DEC_TRACE  
  this.xTraceSPSHeader (pcSPS);
//#endif

  var uiCode uint;
  this.READ_CODE( 4,  &uiCode, "video_parameter_set_id");              
  pcSPS.SetVPSId        ( int(uiCode) );
  this.READ_CODE( 3,  &uiCode, "sps_max_sub_layers_minus1" );          
  pcSPS.SetMaxTLayers   ( uiCode+1 );
//#if MOVE_SPS_TEMPORAL_ID_NESTING_FLAG
  this.READ_FLAG( &uiCode, "sps_temporal_id_nesting_flag" );  
  if uiCode > 0 {             
  	pcSPS.SetTemporalIdNestingFlag ( true );
  }else{
  	pcSPS.SetTemporalIdNestingFlag ( false );
  }
//#else
//  READ_FLAG(     &uiCode, "sps_reserved_zero_bit");               assert(uiCode == 0);
//#endif
  this.ParsePTL(pcSPS.GetPTL(), true, int(pcSPS.GetMaxTLayers()) - 1);
  this.READ_UVLC(     &uiCode, "seq_parameter_set_id" );               
  pcSPS.SetSPSId( int(uiCode) );
  this.READ_UVLC(     &uiCode, "chroma_format_idc" );                  
  pcSPS.SetChromaFormatIdc( int(uiCode) );
  // in the first version we only support chroma_format_idc equal to 1 (4:2:0), so separate_colour_plane_flag cannot appear in the bitstream
  //assert (uiCode == 1);
  if uiCode == 3 {
    this.READ_FLAG(     &uiCode, "separate_colour_plane_flag");        //assert(uiCode == 0);
  }

  this.READ_UVLC (    &uiCode, "pic_width_in_luma_samples" );          
  pcSPS.SetPicWidthInLumaSamples ( uiCode    );
  this.READ_UVLC (    &uiCode, "pic_height_in_luma_samples" );         
  pcSPS.SetPicHeightInLumaSamples( uiCode    );
  this.READ_FLAG(     &uiCode, "pic_cropping_flag");
  if uiCode != 0 {
    crop := pcSPS.GetPicCroppingWindow();
    this.READ_UVLC(   &uiCode, "pic_crop_left_offset" );               
    crop.SetPicCropLeftOffset  ( int(uiCode) * pcSPS.GetCropUnitX( pcSPS.GetChromaFormatIdc() ) );
    this.READ_UVLC(   &uiCode, "pic_crop_right_offset" );              
    crop.SetPicCropRightOffset ( int(uiCode) * pcSPS.GetCropUnitX( pcSPS.GetChromaFormatIdc() ) );
    this.READ_UVLC(   &uiCode, "pic_crop_top_offset" );                
    crop.SetPicCropTopOffset   ( int(uiCode) * pcSPS.GetCropUnitY( pcSPS.GetChromaFormatIdc() ) );
    this.READ_UVLC(   &uiCode, "pic_crop_bottom_offset" );             
    crop.SetPicCropBottomOffset( int(uiCode) * pcSPS.GetCropUnitY( pcSPS.GetChromaFormatIdc() ) );
  }

  this.READ_UVLC(     &uiCode, "bit_depth_luma_minus8" );
  TLibCommon.G_bitDepthY = 8 + int(uiCode);
  pcSPS.SetBitDepthY(TLibCommon.G_bitDepthY);
  pcSPS.SetQpBDOffsetY( int(6*uiCode) );

  this.READ_UVLC( &uiCode,    "bit_depth_chroma_minus8" );
  TLibCommon.G_bitDepthC = 8 + int(uiCode);
  pcSPS.SetBitDepthC(TLibCommon.G_bitDepthC);
  pcSPS.SetQpBDOffsetC( int(6*uiCode) );
/*
#if !HLS_GROUP_SPS_PCM_FLAGS
  this.READ_FLAG( &uiCode, "pcm_enabled_flag" ); pcSPS.SetUsePCM( uiCode ? true : false );

  if( pcSPS->getUsePCM() )
  {
    this.READ_CODE( 4, &uiCode, "pcm_bit_depth_luma_minus1" );           pcSPS.SetPCMBitDepthLuma   ( 1 + uiCode );
    this.READ_CODE( 4, &uiCode, "pcm_bit_depth_chroma_minus1" );         pcSPS.SetPCMBitDepthChroma ( 1 + uiCode );
  }

#endif // !HLS_GROUP_SPS_PCM_FLAGS 
*/
  this.READ_UVLC( &uiCode,    "log2_max_pic_order_cnt_lsb_minus4" );   pcSPS.SetBitsForPOC( 4 + uiCode );

//#if HLS_ADD_SUBLAYER_ORDERING_INFO_PRESENT_FLAG
  var subLayerOrderingInfoPresentFlag uint;
  this.READ_FLAG(&subLayerOrderingInfoPresentFlag, "sps_sub_layer_ordering_info_present_flag");
//#endif // HLS_ADD_SUBLAYER_ORDERING_INFO_PRESENT_FLAG
  for i:=uint(0); i <= pcSPS.GetMaxTLayers()-1; i++ {
    this.READ_UVLC ( &uiCode, "max_dec_pic_buffering");
    pcSPS.SetMaxDecPicBuffering( uiCode, i);
    this.READ_UVLC ( &uiCode, "num_reorder_pics" );
    pcSPS.SetNumReorderPics(int(uiCode), i);
    this.READ_UVLC ( &uiCode, "max_latency_increase");
    pcSPS.SetMaxLatencyIncrease( uiCode, i );

//#if HLS_ADD_SUBLAYER_ORDERING_INFO_PRESENT_FLAG
    if subLayerOrderingInfoPresentFlag==0{
      for i++; i <= pcSPS.GetMaxTLayers()-1; i++ {
        pcSPS.SetMaxDecPicBuffering(pcSPS.GetMaxDecPicBuffering(0), i);
        pcSPS.SetNumReorderPics(pcSPS.GetNumReorderPics(0), i);
        pcSPS.SetMaxLatencyIncrease(pcSPS.GetMaxLatencyIncrease(0), i);
      }
      break;
    }
//#endif // HLS_ADD_SUBLAYER_ORDERING_INFO_PRESENT_FLAG 
  }
/*
#if !HLS_MOVE_SPS_PICLIST_FLAGS
  this.READ_FLAG( &uiCode, "restricted_ref_pic_lists_flag" );
  pcSPS.SetRestrictedRefPicListsFlag( uiCode );
  if( pcSPS->getRestrictedRefPicListsFlag() )
  {
    this.READ_FLAG( &uiCode, "lists_modification_present_flag" );
    pcSPS.SetListsModificationPresentFlag(uiCode);
  }
  else 
  {
    pcSPS.SetListsModificationPresentFlag(true);
  }
#endif // !HLS_MOVE_SPS_PICLIST_FLAGS
*/ 
  this.READ_UVLC( &uiCode, "log2_min_coding_block_size_minus3" );
  log2MinCUSize := uiCode + 3;
  this.READ_UVLC( &uiCode, "log2_diff_max_min_coding_block_size" );
  uiMaxCUDepthCorrect := uiCode;
  pcSPS.SetMaxCUWidth  ( 1<<(log2MinCUSize + uiMaxCUDepthCorrect) ); 
  TLibCommon.G_uiMaxCUWidth  = 1<<(log2MinCUSize + uiMaxCUDepthCorrect);
  pcSPS.SetMaxCUHeight ( 1<<(log2MinCUSize + uiMaxCUDepthCorrect) ); 
  TLibCommon.G_uiMaxCUHeight = 1<<(log2MinCUSize + uiMaxCUDepthCorrect);
  this.READ_UVLC( &uiCode, "log2_min_transform_block_size_minus2" );   
  pcSPS.SetQuadtreeTULog2MinSize( uiCode + 2 );

  this.READ_UVLC( &uiCode, "log2_diff_max_min_transform_block_size" ); 
  pcSPS.SetQuadtreeTULog2MaxSize( uiCode + pcSPS.GetQuadtreeTULog2MinSize() );
  pcSPS.SetMaxTrSize( 1<<(uiCode + pcSPS.GetQuadtreeTULog2MinSize()) );
/*#if !HLS_GROUP_SPS_PCM_FLAGS
  if( pcSPS->getUsePCM() )
  {
    this.READ_UVLC( &uiCode, "log2_min_pcm_coding_block_size_minus3" );  pcSPS.SetPCMLog2MinSize (uiCode+3); 
    this.READ_UVLC( &uiCode, "log2_diff_max_min_pcm_coding_block_size" ); pcSPS.SetPCMLog2MaxSize ( uiCode+pcSPS->getPCMLog2MinSize() );
  }
#endif*/ /* !HLS_GROUP_SPS_PCM_FLAGS */

  this.READ_UVLC( &uiCode, "max_transform_hierarchy_depth_inter" );    
  pcSPS.SetQuadtreeTUMaxDepthInter( uiCode+1 );
  this.READ_UVLC( &uiCode, "max_transform_hierarchy_depth_intra" );    
  pcSPS.SetQuadtreeTUMaxDepthIntra( uiCode+1 );
  TLibCommon.G_uiAddCUDepth = 0;
  for ( pcSPS.GetMaxCUWidth() >> uiMaxCUDepthCorrect ) > ( 1 << ( pcSPS.GetQuadtreeTULog2MinSize() + TLibCommon.G_uiAddCUDepth )  ) {
    TLibCommon.G_uiAddCUDepth++;
  }
  pcSPS.SetMaxCUDepth( uiMaxCUDepthCorrect+TLibCommon.G_uiAddCUDepth  ); 
  TLibCommon.G_uiMaxCUDepth  = uiMaxCUDepthCorrect+TLibCommon.G_uiAddCUDepth;
  // BB: these parameters may be removed completly and replaced by the fixed values
  pcSPS.SetMinTrDepth( 0 );
  pcSPS.SetMaxTrDepth( 1 );
  this.READ_FLAG( &uiCode, "scaling_list_enabled_flag" );                 
  pcSPS.SetScalingListFlag ( uiCode==1 );
  if pcSPS.GetScalingListFlag() {
    this.READ_FLAG( &uiCode, "sps_scaling_list_data_present_flag" );                 
    pcSPS.SetScalingListPresentFlag ( uiCode==1 );
    if pcSPS.GetScalingListPresentFlag () {
      this.ParseScalingList( pcSPS.GetScalingList() );
    }
  }
  this.READ_FLAG( &uiCode, "asymmetric_motion_partitions_enabled_flag" ); 
  pcSPS.SetUseAMP( uiCode==1 );
  this.READ_FLAG( &uiCode, "sample_adaptive_offset_enabled_flag" );       
  if uiCode!=0 {
  	pcSPS.SetUseSAO ( true );
  }else{
    pcSPS.SetUseSAO ( false );
  }
//#if HLS_GROUP_SPS_PCM_FLAGS
  this.READ_FLAG( &uiCode, "pcm_enabled_flag" ); 
  if uiCode!=0 {
  	pcSPS.SetUsePCM( true );
  }else{
    pcSPS.SetUsePCM( false );
  }
//#endif /* HLS_GROUP_SPS_PCM_FLAGS */
  if pcSPS.GetUsePCM() {
//#if !HLS_GROUP_SPS_PCM_FLAGS
//    this.READ_FLAG( &uiCode, "pcm_loop_filter_disable_flag" );           pcSPS.SetPCMFilterDisableFlag ( uiCode ? true : false );
//#else /* HLS_GROUP_SPS_PCM_FLAGS */
    this.READ_CODE( 4, &uiCode, "pcm_sample_bit_depth_luma_minus1" );          
    pcSPS.SetPCMBitDepthLuma   ( 1 + uiCode );
    this.READ_CODE( 4, &uiCode, "pcm_sample_bit_depth_chroma_minus1" );        
    pcSPS.SetPCMBitDepthChroma ( 1 + uiCode );
    this.READ_UVLC( &uiCode, "log2_min_pcm_luma_coding_block_size_minus3" );   
    pcSPS.SetPCMLog2MinSize (uiCode+3);
    this.READ_UVLC( &uiCode, "log2_diff_max_min_pcm_luma_coding_block_size" ); 
    pcSPS.SetPCMLog2MaxSize ( uiCode+pcSPS.GetPCMLog2MinSize() );
    this.READ_FLAG( &uiCode, "pcm_loop_filter_disable_flag" );  
    if uiCode!=0 {               
    	pcSPS.SetPCMFilterDisableFlag ( true );
    }else{
    	pcSPS.SetPCMFilterDisableFlag ( false );
    }
//#endif /* HLS_GROUP_SPS_PCM_FLAGS */
  }

/*#if !MOVE_SPS_TEMPORAL_ID_NESTING_FLAG
  this.READ_FLAG( &uiCode, "temporal_id_nesting_flag" );    
  if uiCode > 0 {            
  	pcSPS.SetTemporalIdNestingFlag (true );
  }else{
  	pcSPS.SetTemporalIdNestingFlag (false );
  }
#endif*/

  this.READ_UVLC( &uiCode, "num_short_term_ref_pic_sets" );
  pcSPS.CreateRPSList(int(uiCode));

  rpsList := pcSPS.GetRPSList();
  //var rps *TComReferencePictureSet;

  for i:=0; i< rpsList.GetNumberOfReferencePictureSets(); i++ {
    rps := rpsList.GetReferencePictureSet(i);
    this.ParseShortTermRefPicSet(pcSPS,rps,i);
  }
  this.READ_FLAG( &uiCode, "long_term_ref_pics_present_flag" );          
  pcSPS.SetLongTermRefsPresent(uiCode!=0);
  if pcSPS.GetLongTermRefsPresent() {
    this.READ_UVLC( &uiCode, "num_long_term_ref_pic_sps" );
    pcSPS.SetNumLongTermRefPicSPS(uiCode);
    for k := 0; k < int(pcSPS.GetNumLongTermRefPicSPS()); k++ {
      this.READ_CODE( pcSPS.GetBitsForPOC(), &uiCode, "lt_ref_pic_poc_lsb_sps" );
      pcSPS.SetLtRefPicPocLsbSps(uint(k), uiCode);
      this.READ_FLAG( &uiCode,  "used_by_curr_pic_lt_sps_flag[i]");
      if uiCode !=0 {
      	pcSPS.SetUsedByCurrPicLtSPSFlag(k, true);
      }else{
        pcSPS.SetUsedByCurrPicLtSPSFlag(k, false);
      }
    }
  }
  this.READ_FLAG( &uiCode, "sps_temporal_mvp_enable_flag" );            
  pcSPS.SetTMVPFlagsPresent(uiCode!=0);

//#if STRONG_INTRA_SMOOTHING
  this.READ_FLAG( &uiCode, "sps_strong_intra_smoothing_enable_flag" );  
  pcSPS.SetUseStrongIntraSmoothing(uiCode!=0);
//#endif

  this.READ_FLAG( &uiCode, "vui_parameters_present_flag" );             
  pcSPS.SetVuiParametersPresentFlag(uiCode!=0);

  if pcSPS.GetVuiParametersPresentFlag() {
    this.ParseVUI(pcSPS.GetVuiParameters(), pcSPS);
  }

  this.READ_FLAG( &uiCode, "sps_extension_flag");
  if uiCode!=0 {
    for this.xMoreRbspData() {
      this.READ_FLAG( &uiCode, "sps_extension_data_flag");
    }
  }
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
//    READ_FLAG( uiCode, "sub_layer_profile_present_flag[i]" ); rpcPTL.SetSubLayerProfilePresentFlag(i, uiCode);
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
  bitsLeft := this.m_pcBitstream.GetNumBitsLeft();

  // if there are more than 8 bits, it cannot be rbsp_trailing_bits
  if bitsLeft > 8{
    return true;
  }

  lastByte := this.m_pcBitstream.PeekBits(bitsLeft);
  cnt := bitsLeft;

  // remove trailing bits equal to zero
  for (cnt>0) && ((lastByte & 1) == 0) {
    lastByte >>= 1;
    cnt--;
  }
  // remove bit equal to one
  cnt--;

  // we should not have a negative number of bits
  //assert (cnt>=0);

  // we have more data, if cnt is not zero
  return cnt>0;
}
  
