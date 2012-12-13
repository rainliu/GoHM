package TLibDecoder

import (
	"fmt"
	"io"
    "gohm/TLibCommon"
)

// ====================================================================================================================
// Class definition
// ====================================================================================================================

type TDecBinIf interface {
    Init              (pcTComBitstream * TLibCommon.TComInputBitstream );
    Uninit            ();
    
    Start             ();
    Finish            ();
    Flush             ();

    DecodeBin         ( ruiBin *uint, rcCtxModel *TLibCommon.ContextModel);
    DecodeBinEP       ( ruiBin *uint                     );
    DecodeBinsEP      ( ruiBins *uint,  numBins  int          );
    DecodeBinTrm      ( ruiBin *uint                     );

    ResetBac          ()                                        ;
    //#if !REMOVE_BURST_IPCM
    //  virtual Void  decodeNumSubseqIPCM( Int& numSubseqIPCM );
    //#endif
    DecodePCMAlignBits()                                         ;
    xReadPCMCode      (  uiLength uint, ruiCode *uint);

    CopyState         ( pcTDecBinIf TDecBinIf);
    GetTDecBinCABAC   () *TDecBinCabac;
    
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
	this.m_pcTComBitstream = pcTComBitstream;
}
func (this *TDecBinCabac)  Uninit            (){
	this.m_pcTComBitstream = nil;
}

func (this *TDecBinCabac)  Start             (){
  //assert( m_pcTComBitstream->getNumBitsUntilByteAligned() == 0 );
  this.m_uiRange    = 510;
  this.m_bitsNeeded = -8;
  this.m_uiValue    = this.m_pcTComBitstream.ReadByte() << 8;
  this.m_uiValue   |= this.m_pcTComBitstream.ReadByte();
}
func (this *TDecBinCabac)  Finish            (){
	//do nothing
}
func (this *TDecBinCabac)  Flush             (){
  for this.m_pcTComBitstream.GetNumBitsLeft() > 0 && 
  	  this.m_pcTComBitstream.GetNumBitsUntilByteAligned() != 0 {
    var uiBits uint;
    this.m_pcTComBitstream.Read ( 1, &uiBits );
  }
  this.Start();
}

func (this *TDecBinCabac)  DecodeBin         ( ruiBin *uint, rcCtxModel *TLibCommon.ContextModel ){
  uiLPS := uint(TLibCommon.TComCABACTables_sm_aucLPSTable[ rcCtxModel.GetState() ][ ( this.m_uiRange >> 6 ) - 4 ]);
  this.m_uiRange -= uiLPS;
  scaledRange := this.m_uiRange << 7;
  
  if this.m_uiValue < scaledRange {
    // MPS path
    *ruiBin = uint(rcCtxModel.GetMps());
    rcCtxModel.UpdateMPS();
    
    if scaledRange >= ( 256 << 7 ) {
      return;
    }
    
    this.m_uiRange = scaledRange >> 6;
    this.m_uiValue += this.m_uiValue;
    this.m_bitsNeeded++;
    if this.m_bitsNeeded == 0 {
      this.m_bitsNeeded = -8;
      this.m_uiValue += this.m_pcTComBitstream.ReadByte();      
    }
  }else{
    // LPS path
    numBits := TLibCommon.TComCABACTables_sm_aucRenormTable[ uiLPS >> 3 ];
    this.m_uiValue   = ( this.m_uiValue - scaledRange ) << numBits;
    this.m_uiRange   = uiLPS << numBits;
    *ruiBin      = uint(1 - rcCtxModel.GetMps());
    rcCtxModel.UpdateLPS();
    
    this.m_bitsNeeded += int(numBits);
    
    if this.m_bitsNeeded >= 0 {
      this.m_uiValue += this.m_pcTComBitstream.ReadByte() << uint(this.m_bitsNeeded);
      this.m_bitsNeeded -= 8;
    }
  }
}
func (this *TDecBinCabac)  DecodeBinEP       ( ruiBin *uint                           ){
  this.m_uiValue += this.m_uiValue;
  this.m_bitsNeeded++
  if this.m_bitsNeeded >= 0 {
    this.m_bitsNeeded = -8;
    this.m_uiValue += this.m_pcTComBitstream.ReadByte();
  }
  
  *ruiBin = 0;
  scaledRange := this.m_uiRange << 7;
  if this.m_uiValue >= scaledRange {
    *ruiBin = 1;
    this.m_uiValue -= scaledRange;
  }
}
func (this *TDecBinCabac)  DecodeBinsEP      ( ruiBin *uint, numBins int              ){
  bins := uint(0);
  
  for numBins > 8 {
    this.m_uiValue = ( this.m_uiValue << 8 ) + ( this.m_pcTComBitstream.ReadByte() << uint( 8 + this.m_bitsNeeded ) );
    
    scaledRange := this.m_uiRange << 15;
    for i := 0; i < 8; i++ {
      bins += bins;
      scaledRange >>= 1;
      if this.m_uiValue >= scaledRange {
        bins++;
        this.m_uiValue -= scaledRange;
      }
    }
    numBins -= 8;
  }
  
  this.m_bitsNeeded += numBins;
  this.m_uiValue <<= uint(numBins);
  
  if this.m_bitsNeeded >= 0 {
    this.m_uiValue += this.m_pcTComBitstream.ReadByte() << uint(this.m_bitsNeeded);
    this.m_bitsNeeded -= 8;
  }
  
  scaledRange := this.m_uiRange << uint( numBins + 7 );
  for i := 0; i < numBins; i++ {
    bins += bins;
    scaledRange >>= 1;
    if this.m_uiValue >= scaledRange {
      bins++;
      this.m_uiValue -= scaledRange;
    }
  }
  
  *ruiBin = bins;
}
func (this *TDecBinCabac)  DecodeBinTrm      ( ruiBin *uint                           ){
  this.m_uiRange -= 2;
  scaledRange := this.m_uiRange << 7;
  if this.m_uiValue >= scaledRange {
    *ruiBin = 1;
  }else{
    *ruiBin = 0;
    if scaledRange < ( 256 << 7 ) {
      this.m_uiRange = scaledRange >> 6;
      this.m_uiValue += this.m_uiValue;
      this.m_bitsNeeded++;
      if this.m_bitsNeeded == 0 {
        this.m_bitsNeeded = -8;
        this.m_uiValue += this.m_pcTComBitstream.ReadByte();      
      }
    }
  }
}

func (this *TDecBinCabac)  ResetBac          (){
  this.m_uiRange    = 510;
  this.m_bitsNeeded = -8;
  this.m_uiValue    = this.m_pcTComBitstream.ReadBits( 16 );
}
//#if !REMOVE_BURST_IPCM
//  Void  decodeNumSubseqIPCM( Int& numSubseqIPCM ) ;
//#endif
func (this *TDecBinCabac)  DecodePCMAlignBits(){
  iNum := this.m_pcTComBitstream.GetNumBitsUntilByteAligned();
  
  uiBit := uint(0);
  this.m_pcTComBitstream.Read( iNum, &uiBit );
}
func (this *TDecBinCabac)  xReadPCMCode      (  uiLength uint, ruiCode *uint){
  //assert ( uiLength > 0 );
  this.m_pcTComBitstream.Read (uiLength, ruiCode);
}

func (this *TDecBinCabac)  CopyState         ( pcTDecBinIf TDecBinIf){
  pcTDecBinCABAC := pcTDecBinIf.GetTDecBinCABAC();
  this.m_uiRange   = pcTDecBinCABAC.m_uiRange;
  this.m_uiValue   = pcTDecBinCABAC.m_uiValue;
  this.m_bitsNeeded= pcTDecBinCABAC.m_bitsNeeded;
}
func (this *TDecBinCabac)  GetTDecBinCABAC()  *TDecBinCabac{ 
	return this; 
}


//class SEImessages;

/// SBAC decoder class
type TDecSbac struct { //: public TDecEntropyIf
    //private:
    m_pTraceFile	io.Writer;
    m_pcBitstream 	*TLibCommon.TComInputBitstream
    m_pcTDecBinIf		TDecBinIf;

    //private:
    m_uiLastDQpNonZero uint
    m_uiLastQp         uint

    m_contextModels             [TLibCommon.MAX_NUM_CTX_MOD]TLibCommon.ContextModel
    m_numContextModels          int
    
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

    m_cCUSigCoeffGroupSCModel 	*TLibCommon.ContextModel3DBuffer
    m_cCUSigSCModel           	*TLibCommon.ContextModel3DBuffer
    m_cCuCtxLastX             	*TLibCommon.ContextModel3DBuffer
    m_cCuCtxLastY             	*TLibCommon.ContextModel3DBuffer
    m_cCUOneSCModel           	*TLibCommon.ContextModel3DBuffer
    m_cCUAbsSCModel           	*TLibCommon.ContextModel3DBuffer

    m_cMVPIdxSCModel 			*TLibCommon.ContextModel3DBuffer

    m_cCUAMPSCModel                 *TLibCommon.ContextModel3DBuffer
    m_cSaoMergeSCModel              *TLibCommon.ContextModel3DBuffer
    m_cSaoTypeIdxSCModel            *TLibCommon.ContextModel3DBuffer
    m_cTransformSkipSCModel         *TLibCommon.ContextModel3DBuffer
    m_CUTransquantBypassFlagSCModel *TLibCommon.ContextModel3DBuffer
}

func (this *TDecSbac)  XTraceLCUHeader (traceLevel uint){
  if (traceLevel & TLibCommon.TRACE_LEVEL) !=0 {
  	io.WriteString(this.m_pTraceFile, "========= LCU Parameter Set ===============================================\n");//, pLCU.GetAddr());
  }
}

func (this *TDecSbac)  xTraceCUHeader (traceLevel uint){
  if (traceLevel & TLibCommon.TRACE_LEVEL) !=0 {
  	io.WriteString(this.m_pTraceFile, "========= CU Parameter Set ================================================\n");//, pCU.GetCUPelX(), pCU.GetCUPelY());
  }
}

func (this *TDecSbac)  xTracePUHeader (traceLevel uint){
  if (traceLevel & TLibCommon.TRACE_LEVEL) !=0 {
    io.WriteString(this.m_pTraceFile, "========= PU Parameter Set ================================================\n");//, pCU.GetCUPelX(), pCU.GetCUPelY());
  }
}

func (this *TDecSbac)  xTraceTUHeader (traceLevel uint){
  if (traceLevel & TLibCommon.TRACE_LEVEL) !=0 {
    io.WriteString(this.m_pTraceFile, "========= TU Parameter Set ================================================\n");//, pCU.GetCUPelX(), pCU.GetCUPelY());
  }
}

func (this *TDecSbac)  xTraceCoefHeader (traceLevel uint){
  if (traceLevel & TLibCommon.TRACE_LEVEL) !=0 {
    io.WriteString(this.m_pTraceFile, "========= Coefficient Parameter Set =======================================\n");//, pCU.GetCUPelX(), pCU.GetCUPelY());
  }
}

func (this *TDecSbac)  xTraceResiHeader (traceLevel uint){
  if (traceLevel & TLibCommon.TRACE_LEVEL) !=0 {
    io.WriteString(this.m_pTraceFile, "========= Residual Parameter Set ==========================================\n");//, pCU.GetCUPelX(), pCU.GetCUPelY());
  }
}

func (this *TDecSbac) xTracePredHeader (traceLevel uint){
  if (traceLevel & TLibCommon.TRACE_LEVEL) !=0 {
    io.WriteString(this.m_pTraceFile, "========= Prediction Parameter Set ========================================\n");//, pCU.GetCUPelX(), pCU.GetCUPelY());
  }
}

func (this *TDecSbac)  xTraceRecoHeader (traceLevel uint){
  if (traceLevel & TLibCommon.TRACE_LEVEL) !=0 {
    io.WriteString(this.m_pTraceFile, "========= Reconstruction Parameter Set ====================================\n");//, pCU.GetCUPelX(), pCU.GetCUPelY());
  }	
}

func (this *TDecSbac)  XReadAeTr ( Value int, pSymbolName string,  traceLevel uint){
  if (traceLevel & TLibCommon.TRACE_LEVEL) !=0 {
    //fprintf( g_hTrace, "%8lld  ", g_nSymbolCounter++ );
    io.WriteString(this.m_pTraceFile, fmt.Sprintf ("%-62s ae(v) : %4d\n", pSymbolName, Value )); 
    //fflush ( g_hTrace );
  }
}


func NewTDecSbac() *TDecSbac{
	pTDecSbac := &TDecSbac{ m_pcBitstream : nil, m_pcTDecBinIf : nil, m_numContextModels : 0};
	pTDecSbac.xInit();
	
	return pTDecSbac;
}

func (this *TDecSbac) xInit(){
	this.m_cCUSplitFlagSCModel 		 	= TLibCommon.NewContextModel3DBuffer( 1, 1, TLibCommon.NUM_SPLIT_FLAG_CTX            	, this.m_contextModels[this.m_numContextModels:], &this.m_numContextModels);
	this.m_cCUSkipFlagSCModel        	= TLibCommon.NewContextModel3DBuffer( 1, 1, TLibCommon.NUM_SKIP_FLAG_CTX             	, this.m_contextModels[this.m_numContextModels:], &this.m_numContextModels)
	this.m_cCUMergeFlagExtSCModel    	= TLibCommon.NewContextModel3DBuffer( 1, 1, TLibCommon.NUM_MERGE_FLAG_EXT_CTX        	, this.m_contextModels[this.m_numContextModels:], &this.m_numContextModels)
	this.m_cCUMergeIdxExtSCModel     	= TLibCommon.NewContextModel3DBuffer( 1, 1, TLibCommon.NUM_MERGE_IDX_EXT_CTX         	, this.m_contextModels[this.m_numContextModels:], &this.m_numContextModels)
	this.m_cCUPartSizeSCModel        	= TLibCommon.NewContextModel3DBuffer( 1, 1, TLibCommon.NUM_PART_SIZE_CTX             	, this.m_contextModels[this.m_numContextModels:], &this.m_numContextModels)
	this.m_cCUPredModeSCModel        	= TLibCommon.NewContextModel3DBuffer( 1, 1, TLibCommon.NUM_PRED_MODE_CTX             	, this.m_contextModels[this.m_numContextModels:], &this.m_numContextModels)
	this.m_cCUIntraPredSCModel       	= TLibCommon.NewContextModel3DBuffer( 1, 1, TLibCommon.NUM_ADI_CTX                   	, this.m_contextModels[this.m_numContextModels:], &this.m_numContextModels)
	this.m_cCUChromaPredSCModel      	= TLibCommon.NewContextModel3DBuffer( 1, 1, TLibCommon.NUM_CHROMA_PRED_CTX           	, this.m_contextModels[this.m_numContextModels:], &this.m_numContextModels)
	this.m_cCUDeltaQpSCModel         	= TLibCommon.NewContextModel3DBuffer( 1, 1, TLibCommon.NUM_DELTA_QP_CTX              	, this.m_contextModels[this.m_numContextModels:], &this.m_numContextModels)
	this.m_cCUInterDirSCModel        	= TLibCommon.NewContextModel3DBuffer( 1, 1, TLibCommon.NUM_INTER_DIR_CTX             	, this.m_contextModels[this.m_numContextModels:], &this.m_numContextModels)
	this.m_cCURefPicSCModel          	= TLibCommon.NewContextModel3DBuffer( 1, 1, TLibCommon.NUM_REF_NO_CTX                	, this.m_contextModels[this.m_numContextModels:], &this.m_numContextModels)
	this.m_cCUMvdSCModel             	= TLibCommon.NewContextModel3DBuffer( 1, 1, TLibCommon.NUM_MV_RES_CTX                	, this.m_contextModels[this.m_numContextModels:], &this.m_numContextModels)
	this.m_cCUQtCbfSCModel           	= TLibCommon.NewContextModel3DBuffer( 1, 2, TLibCommon.NUM_QT_CBF_CTX                	, this.m_contextModels[this.m_numContextModels:], &this.m_numContextModels)
	this.m_cCUTransSubdivFlagSCModel 	= TLibCommon.NewContextModel3DBuffer( 1, 1, TLibCommon.NUM_TRANS_SUBDIV_FLAG_CTX     	, this.m_contextModels[this.m_numContextModels:], &this.m_numContextModels)
	this.m_cCUQtRootCbfSCModel       	= TLibCommon.NewContextModel3DBuffer( 1, 1, TLibCommon.NUM_QT_ROOT_CBF_CTX           	, this.m_contextModels[this.m_numContextModels:], &this.m_numContextModels)
	this.m_cCUSigCoeffGroupSCModel   	= TLibCommon.NewContextModel3DBuffer( 1, 2, TLibCommon.NUM_SIG_CG_FLAG_CTX           	, this.m_contextModels[this.m_numContextModels:], &this.m_numContextModels)
	this.m_cCUSigSCModel             	= TLibCommon.NewContextModel3DBuffer( 1, 1, TLibCommon.NUM_SIG_FLAG_CTX              	, this.m_contextModels[this.m_numContextModels:], &this.m_numContextModels)
	this.m_cCuCtxLastX               	= TLibCommon.NewContextModel3DBuffer( 1, 2, TLibCommon.NUM_CTX_LAST_FLAG_XY          	, this.m_contextModels[this.m_numContextModels:], &this.m_numContextModels)
	this.m_cCuCtxLastY               	= TLibCommon.NewContextModel3DBuffer( 1, 2, TLibCommon.NUM_CTX_LAST_FLAG_XY          	, this.m_contextModels[this.m_numContextModels:], &this.m_numContextModels)
	this.m_cCUOneSCModel             	= TLibCommon.NewContextModel3DBuffer( 1, 1, TLibCommon.NUM_ONE_FLAG_CTX              	, this.m_contextModels[this.m_numContextModels:], &this.m_numContextModels)
	this.m_cCUAbsSCModel             	= TLibCommon.NewContextModel3DBuffer( 1, 1, TLibCommon.NUM_ABS_FLAG_CTX              	, this.m_contextModels[this.m_numContextModels:], &this.m_numContextModels)
	this.m_cMVPIdxSCModel            	= TLibCommon.NewContextModel3DBuffer( 1, 1, TLibCommon.NUM_MVP_IDX_CTX               	, this.m_contextModels[this.m_numContextModels:], &this.m_numContextModels)
	this.m_cCUAMPSCModel             	= TLibCommon.NewContextModel3DBuffer( 1, 1, TLibCommon.NUM_CU_AMP_CTX                	, this.m_contextModels[this.m_numContextModels:], &this.m_numContextModels)
	this.m_cSaoMergeSCModel          	= TLibCommon.NewContextModel3DBuffer( 1, 1, TLibCommon.NUM_SAO_MERGE_FLAG_CTX   	 	, this.m_contextModels[this.m_numContextModels:], &this.m_numContextModels)
	this.m_cSaoTypeIdxSCModel        	= TLibCommon.NewContextModel3DBuffer( 1, 1, TLibCommon.NUM_SAO_TYPE_IDX_CTX          	, this.m_contextModels[this.m_numContextModels:], &this.m_numContextModels)
	this.m_cTransformSkipSCModel     	= TLibCommon.NewContextModel3DBuffer( 1, 2, TLibCommon.NUM_TRANSFORMSKIP_FLAG_CTX    	, this.m_contextModels[this.m_numContextModels:], &this.m_numContextModels)
	this.m_CUTransquantBypassFlagSCModel= TLibCommon.NewContextModel3DBuffer( 1, 1, TLibCommon.NUM_CU_TRANSQUANT_BYPASS_FLAG_CTX, this.m_contextModels[this.m_numContextModels:], &this.m_numContextModels)
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
  sliceType  := pSlice.GetSliceType();
  qp         := pSlice.GetSliceQp();

  if pSlice.GetPPS().GetCabacInitPresentFlag() && pSlice.GetCabacInitFlag() {
    switch sliceType {
    case TLibCommon.P_SLICE:           // change initialization table to B_SLICE initialization
      sliceType = TLibCommon.B_SLICE; 
      //break;
    case TLibCommon.B_SLICE:           // change initialization table to P_SLICE initialization
      sliceType = TLibCommon.P_SLICE; 
      //break;
    //default     :           // should not occur
      //assert(0);
    }
  }

  this.m_cCUSplitFlagSCModel.InitBuffer        	 ( sliceType, qp, TLibCommon.INIT_SPLIT_FLAG[:] );
  this.m_cCUSkipFlagSCModel.InitBuffer        	 ( sliceType, qp, TLibCommon.INIT_SKIP_FLAG[:] );
  this.m_cCUMergeFlagExtSCModel.InitBuffer    	 ( sliceType, qp, TLibCommon.INIT_MERGE_FLAG_EXT[:] );
  this.m_cCUMergeIdxExtSCModel.InitBuffer     	 ( sliceType, qp, TLibCommon.INIT_MERGE_IDX_EXT[:] );
  this.m_cCUPartSizeSCModel.InitBuffer        	 ( sliceType, qp, TLibCommon.INIT_PART_SIZE[:] );
  this.m_cCUAMPSCModel.InitBuffer             	 ( sliceType, qp, TLibCommon.INIT_CU_AMP_POS[:] );
  this.m_cCUPredModeSCModel.InitBuffer        	 ( sliceType, qp, TLibCommon.INIT_PRED_MODE[:] );
  this.m_cCUIntraPredSCModel.InitBuffer       	 ( sliceType, qp, TLibCommon.INIT_INTRA_PRED_MODE[:] );
  this.m_cCUChromaPredSCModel.InitBuffer      	 ( sliceType, qp, TLibCommon.INIT_CHROMA_PRED_MODE[:] );
  this.m_cCUInterDirSCModel.InitBuffer        	 ( sliceType, qp, TLibCommon.INIT_INTER_DIR[:] );
  this.m_cCUMvdSCModel.InitBuffer             	 ( sliceType, qp, TLibCommon.INIT_MVD[:] );
  this.m_cCURefPicSCModel.InitBuffer          	 ( sliceType, qp, TLibCommon.INIT_REF_PIC[:] );
  this.m_cCUDeltaQpSCModel.InitBuffer         	 ( sliceType, qp, TLibCommon.INIT_DQP[:] );
  this.m_cCUQtCbfSCModel.InitBuffer           	 ( sliceType, qp, TLibCommon.INIT_QT_CBF[:] );
  this.m_cCUQtRootCbfSCModel.InitBuffer       	 ( sliceType, qp, TLibCommon.INIT_QT_ROOT_CBF[:] );
  this.m_cCUSigCoeffGroupSCModel.InitBuffer   	 ( sliceType, qp, TLibCommon.INIT_SIG_CG_FLAG[:] );
  this.m_cCUSigSCModel.InitBuffer              	 ( sliceType, qp, TLibCommon.INIT_SIG_FLAG[:] );
  this.m_cCuCtxLastX.InitBuffer               	 ( sliceType, qp, TLibCommon.INIT_LAST[:] );
  this.m_cCuCtxLastY.InitBuffer               	 ( sliceType, qp, TLibCommon.INIT_LAST[:] );
  this.m_cCUOneSCModel.InitBuffer             	 ( sliceType, qp, TLibCommon.INIT_ONE_FLAG[:] );
  this.m_cCUAbsSCModel.InitBuffer             	 ( sliceType, qp, TLibCommon.INIT_ABS_FLAG[:] );
  this.m_cMVPIdxSCModel.InitBuffer            	 ( sliceType, qp, TLibCommon.INIT_MVP_IDX[:] );
  this.m_cSaoMergeSCModel.InitBuffer     	  	 ( sliceType, qp, TLibCommon.INIT_SAO_MERGE_FLAG[:] );
  this.m_cSaoTypeIdxSCModel.InitBuffer        	 ( sliceType, qp, TLibCommon.INIT_SAO_TYPE_IDX[:] );
  this.m_cCUTransSubdivFlagSCModel.InitBuffer 	 ( sliceType, qp, TLibCommon.INIT_TRANS_SUBDIV_FLAG[:] );
  this.m_cTransformSkipSCModel.InitBuffer     	 ( sliceType, qp, TLibCommon.INIT_TRANSFORMSKIP_FLAG[:] );
  this.m_CUTransquantBypassFlagSCModel.InitBuffer( sliceType, qp, TLibCommon.INIT_CU_TRANSQUANT_BYPASS_FLAG[:] );
  
  this.m_uiLastDQpNonZero  = 0;
  
  // new structure
  this.m_uiLastQp          = uint(qp);
  
  this.m_pcTDecBinIf.Start();
}
func (this *TDecSbac)   SetBitstream            ( p  *TLibCommon.TComInputBitstream) {
	this.m_pcBitstream = p;
	this.m_pcTDecBinIf.Init( p );
}
func (this *TDecSbac)   SetTraceFile 		      ( traceFile io.Writer){
	this.m_pTraceFile = traceFile;
}
func (this *TDecSbac)   SetSliceTrace 		      ( bSliceTrace bool){
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
