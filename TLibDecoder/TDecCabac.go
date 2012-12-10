package TLibDecoder

import (
	"io"
    "gohm/TLibCommon"
)

// ====================================================================================================================
// Class definition
// ====================================================================================================================

type TDecBinIf interface {
    Init              (pcTComBitstream * TLibCommon.TComInputBitstream );
    Uninit            ();
    /*
      virtual Void  start             ()                                          = 0;
      virtual Void  finish            ()                                          = 0;
      virtual Void  flush            ()                                           = 0;

      virtual Void  decodeBin         ( UInt& ruiBin, ContextModel& rcCtxModel )  = 0;
      virtual Void  decodeBinEP       ( UInt& ruiBin                           )  = 0;
      virtual Void  decodeBinsEP      ( UInt& ruiBins, Int numBins             )  = 0;
      virtual Void  decodeBinTrm      ( UInt& ruiBin                           )  = 0;

      virtual Void  resetBac          ()                                          = 0;
    #if !REMOVE_BURST_IPCM
      virtual Void  decodeNumSubseqIPCM( Int& numSubseqIPCM )                  = 0;
    #endif
      virtual Void  decodePCMAlignBits()                                          = 0;
      virtual Void  xReadPCMCode      ( UInt uiLength, UInt& ruiCode)              = 0;

      virtual ~TDecBinIf() {}

      virtual Void  copyState         ( TDecBinIf* pcTDecBinIf )                  = 0;
      virtual TDecBinCABAC*   getTDecBinCABAC   ()  { return 0; }
    */
}

type TDecBinCabac struct { //: public TDecBinIf
    //private:
    m_pcTComBitstream *TLibCommon.TComInputBitstream
    m_uiRange         uint
    m_uiValue         uint
    m_bitsNeeded      int
}

func NewTDecBinCabac () *TDecBinCabac{
	return &TDecBinCabac{}
}

func (this *TDecBinCabac)  Init              ( pcTComBitstream *TLibCommon.TComInputBitstream){
}
func (this *TDecBinCabac)  Uninit            (){
}
/*
  Void  start             ();
  Void  finish            ();
  Void  flush             ();

  Void  decodeBin         ( UInt& ruiBin, ContextModel& rcCtxModel );
  Void  decodeBinEP       ( UInt& ruiBin                           );
  Void  decodeBinsEP      ( UInt& ruiBin, Int numBins              );
  Void  decodeBinTrm      ( UInt& ruiBin                           );

  Void  resetBac          ();
#if !REMOVE_BURST_IPCM
  Void  decodeNumSubseqIPCM( Int& numSubseqIPCM ) ;
#endif
  Void  decodePCMAlignBits();
  Void  xReadPCMCode      ( UInt uiLength, UInt& ruiCode );

  Void  copyState         ( TDecBinIf* pcTDecBinIf );
  TDecBinCABAC* getTDecBinCABAC()  { return this; }


};*/

//class SEImessages;

/// SBAC decoder class
type TDecSbac struct { //: public TDecEntropyIf
    //private:
    m_pcBitstream *TLibCommon.TComInputBitstream
    m_pcTDecBinIf		TDecBinIf;

    //private:
    m_uiLastDQpNonZero uint
    m_uiLastQp         uint

    m_contextModels             [TLibCommon.MAX_NUM_CTX_MOD]TLibCommon.ContextModel
    m_numContextModels          int
    m_cCUSplitFlagSCModel       TLibCommon.ContextModel3DBuffer
    m_cCUSkipFlagSCModel        TLibCommon.ContextModel3DBuffer
    m_cCUMergeFlagExtSCModel    TLibCommon.ContextModel3DBuffer
    m_cCUMergeIdxExtSCModel     TLibCommon.ContextModel3DBuffer
    m_cCUPartSizeSCModel        TLibCommon.ContextModel3DBuffer
    m_cCUPredModeSCModel        TLibCommon.ContextModel3DBuffer
    m_cCUIntraPredSCModel       TLibCommon.ContextModel3DBuffer
    m_cCUChromaPredSCModel      TLibCommon.ContextModel3DBuffer
    m_cCUDeltaQpSCModel         TLibCommon.ContextModel3DBuffer
    m_cCUInterDirSCModel        TLibCommon.ContextModel3DBuffer
    m_cCURefPicSCModel          TLibCommon.ContextModel3DBuffer
    m_cCUMvdSCModel             TLibCommon.ContextModel3DBuffer
    m_cCUQtCbfSCModel           TLibCommon.ContextModel3DBuffer
    m_cCUTransSubdivFlagSCModel TLibCommon.ContextModel3DBuffer
    m_cCUQtRootCbfSCModel       TLibCommon.ContextModel3DBuffer

    m_cCUSigCoeffGroupSCModel TLibCommon.ContextModel3DBuffer
    m_cCUSigSCModel           TLibCommon.ContextModel3DBuffer
    m_cCuCtxLastX             TLibCommon.ContextModel3DBuffer
    m_cCuCtxLastY             TLibCommon.ContextModel3DBuffer
    m_cCUOneSCModel           TLibCommon.ContextModel3DBuffer
    m_cCUAbsSCModel           TLibCommon.ContextModel3DBuffer

    m_cMVPIdxSCModel TLibCommon.ContextModel3DBuffer

    m_cCUAMPSCModel                 TLibCommon.ContextModel3DBuffer
    m_cSaoMergeSCModel              TLibCommon.ContextModel3DBuffer
    m_cSaoTypeIdxSCModel            TLibCommon.ContextModel3DBuffer
    m_cTransformSkipSCModel         TLibCommon.ContextModel3DBuffer
    m_CUTransquantBypassFlagSCModel TLibCommon.ContextModel3DBuffer
}


func NewTDecSbac() *TDecSbac{
	return &TDecSbac{}
}

  

func (this *TDecSbac) Init ( p TDecBinIf)    { 
	this.m_pcTDecBinIf = p; 
}

func (this *TDecSbac)   Uninit                  ( )    { 
	this.m_pcTDecBinIf = nil; 
}

func (this *TDecSbac)   Load                    ( pScr *TDecSbac){
}
func (this *TDecSbac)   LoadContexts            ( pScr *TDecSbac){
}
func (this *TDecSbac)   xCopyFrom           	( pSrc *TDecSbac){
}
func (this *TDecSbac)   xCopyContextsFrom       ( pSrc *TDecSbac){
}

func (this *TDecSbac)   ResetEntropy 			( pSlice *TLibCommon.TComSlice){
}
func (this *TDecSbac)   SetBitstream            ( p  *TLibCommon.TComInputBitstream) { 
	this.m_pcBitstream = p; 
	this.m_pcTDecBinIf.Init( p ); 
}
func (this *TDecSbac)   SetTraceFile 		      ( traceFile io.Writer){
}
func (this *TDecSbac)   ParseVPS                  ( pcVPS *TLibCommon.TComVPS )  {
}
func (this *TDecSbac)   ParseSPS                  ( pcSPS *TLibCommon.TComSPS         ) {
}
func (this *TDecSbac)   ParsePPS                  ( pcPPS *TLibCommon.TComPPS        ) {
}

func (this *TDecSbac)   ParseSliceHeader          ( rpcSlice *TLibCommon.TComSlice, parameterSetManager *TLibCommon.ParameterSetManager) {
}
func (this *TDecSbac)   ParseTerminatingBit       ( ruiBit *uint){
}
func (this *TDecSbac)   ParseMVPIdx               ( riMVPIdx  *int ){
}
func (this *TDecSbac)   ParseSaoMaxUvlc           ( val *uint,  maxSymbol uint){
}
func (this *TDecSbac)   ParseSaoMerge         	  ( ruiVal *uint  ){
}
func (this *TDecSbac)   ParseSaoTypeIdx           ( ruiVal *uint ){
}
func (this *TDecSbac)   ParseSaoUflc              ( uiLength uint, ruiVal *uint    ){
}
func (this *TDecSbac)   ParseSaoOneLcuInterleaving( rx,  ry int, pSaoParam *TLibCommon.SAOParam, pcCU *TLibCommon.TComDataCU,  iCUAddrInSlice,  iCUAddrUpInSlice,  allowMergeLeft,  allowMergeUp int){
}
func (this *TDecSbac)   ParseSaoOffset            (psSaoLcuParam *TLibCommon.SaoLcuParam,  compIdx uint){
}
//private:
func (this *TDecSbac)   xReadUnarySymbol    ( ruiSymbol *uint, pcSCModel *TLibCommon.ContextModel,  iOffset int){
}
func (this *TDecSbac)   xReadUnaryMaxSymbol ( ruiSymbol *uint, pcSCModel *TLibCommon.ContextModel,  iOffset,  uiMaxSymbol uint ){
}
func (this *TDecSbac)   xReadEpExGolomb     ( ruiSymbol *uint,  uiCount uint){
}
func (this *TDecSbac)   xReadCoefRemainExGolomb ( ruiSymbol *uint, rParam *uint){
}

//public:
func (this *TDecSbac)  ParseSkipFlag      		  ( pcCU *TLibCommon.TComDataCU,  uiAbsPartIdx,  uiDepth uint){
}
func (this *TDecSbac)  ParseCUTransquantBypassFlag( pcCU *TLibCommon.TComDataCU,  uiAbsPartIdx,  uiDepth uint ){
}
func (this *TDecSbac)  ParseSplitFlag     ( pcCU *TLibCommon.TComDataCU,  uiAbsPartIdx,  uiDepth uint ){
}
func (this *TDecSbac)  ParseMergeFlag     ( pcCU *TLibCommon.TComDataCU,  uiAbsPartIdx,  uiDepth, uiPUIdx uint  ){
}
func (this *TDecSbac)  ParseMergeIndex    ( pcCU *TLibCommon.TComDataCU, ruiMergeIndex *uint,  uiAbsPartIdx,  uiDepth uint){
}
func (this *TDecSbac)  ParsePartSize      ( pcCU *TLibCommon.TComDataCU, uiAbsPartIdx,  uiDepth uint ){
}
func (this *TDecSbac)  ParsePredMode      ( pcCU *TLibCommon.TComDataCU, uiAbsPartIdx,  uiDepth uint ){
}

func (this *TDecSbac)  ParseIntraDirLumaAng( pcCU *TLibCommon.TComDataCU, uiAbsPartIdx,  uiDepth uint ){
}

func (this *TDecSbac)  ParseIntraDirChroma( pcCU *TLibCommon.TComDataCU, uiAbsPartIdx,  uiDepth uint ){
}

func (this *TDecSbac)  ParseInterDir      ( pcCU *TLibCommon.TComDataCU, ruiInterDir *uint, uiAbsPartIdx,  uiDepth uint ){
}
func (this *TDecSbac)  ParseRefFrmIdx     ( pcCU *TLibCommon.TComDataCU, riRefFrmIdx *int, uiAbsPartIdx,  uiDepth uint,  eRefList TLibCommon.RefPicList){
}
func (this *TDecSbac)  ParseMvd           ( pcCU *TLibCommon.TComDataCU,  uiAbsPartIdx,  uiPartIdx,  uiDepth uint, eRefList TLibCommon.RefPicList){
}

func (this *TDecSbac)  ParseTransformSubdivFlag(  ruiSubdivFlag *uint,  uiLog2TransformBlockSize uint){
}
func (this *TDecSbac)  ParseQtCbf         ( pcCU *TLibCommon.TComDataCU,  uiAbsPartIdx uint,  eType TLibCommon.TextType,  uiTrDepth, uiDepth uint  ){
}
func (this *TDecSbac)  ParseQtRootCbf     ( pcCU *TLibCommon.TComDataCU,  uiAbsPartIdx,  uiDepth uint, uiQtRootCbf *uint ){
}

func (this *TDecSbac)  ParseDeltaQP       ( pcCU *TLibCommon.TComDataCU,  uiAbsPartIdx,  uiDepth uint){
}

func (this *TDecSbac)  ParseIPCMInfo      ( pcCU *TLibCommon.TComDataCU, uiAbsPartIdx,  uiDepth uint){
}

func (this *TDecSbac)  ParseLastSignificantXY( uiPosLastX *uint, uiPosLastY *uint,  width,  height int,  eTType TLibCommon.TextType,  uiScanIdx uint){
}
func (this *TDecSbac)  ParseCoeffNxN      ( pcCU *TLibCommon.TComDataCU, pcCoef *TLibCommon.TCoeff,  uiAbsPartIdx,  uiWidth,  uiHeight,  uiDepth uint,  eTType TLibCommon.TextType){
}
func (this *TDecSbac)  ParseTransformSkipFlags ( pcCU *TLibCommon.TComDataCU,  uiAbsPartIdx,  width,  height,  uiDepth uint,  eTType TLibCommon.TextType){
}

func (this *TDecSbac)  UpdateContextTables(  eSliceType TLibCommon.SliceType,  iQp int ){
}

func (this *TDecSbac)  ParseScalingList ( scalingList *TLibCommon.TComScalingList) {
}

