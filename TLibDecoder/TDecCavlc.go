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
func (this *TDecCavlc)  ParseSliceHeader    ( rpcSlice *TLibCommon.TComSlice, parameterSetManager *ParameterSetManagerDecoder){
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
  
