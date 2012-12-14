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
    //do nothing
}
func (this *TDecSbac)   ParseVPS                  ( pcVPS *TLibCommon.TComVPS )  {
    //do nothing
}
func (this *TDecSbac)   ParseSPS                  ( pcSPS *TLibCommon.TComSPS         ) {
    //do nothing
}
func (this *TDecSbac)   ParsePPS                  ( pcPPS *TLibCommon.TComPPS        ) {
    //do nothing
}

func (this *TDecSbac)   ParseSliceHeader          ( rpcSlice *TLibCommon.TComSlice, parameterSetManager *TLibCommon.ParameterSetManager) {
    //do nothing
}
func (this *TDecSbac)   ParseTerminatingBit       ( ruiBit *uint){
    this.m_pcTDecBinIf.DecodeBinTrm( ruiBit );
}
func (this *TDecSbac)   ParseMVPIdx               ( riMVPIdx  *int ){
}
func (this *TDecSbac)   ParseSaoMaxUvlc           ( val *uint,  maxSymbol uint){
    if maxSymbol == 0 {
      *val = 0;
      return;
    }

    var code uint;
    var i uint;
    this.m_pcTDecBinIf.DecodeBinEP( &code );
    if code == 0 {
      *val = 0;
      return;
    }

    i=1;
    for {
      this.m_pcTDecBinIf.DecodeBinEP( &code );
      if code == 0 {
        break;
      }
      i++;
      if i == maxSymbol{
        break;
      }
    }

    *val = i;
}
func (this *TDecSbac)   ParseSaoMerge         	  ( ruiVal *uint  ){
    var uiCode uint;
    this.m_pcTDecBinIf.DecodeBin( &uiCode, this.m_cSaoMergeSCModel.Get3( 0, 0, 0 ) );
    *ruiVal = uiCode;
}
func (this *TDecSbac)   ParseSaoTypeIdx           ( ruiVal *uint ){
    var uiCode uint;
    this.m_pcTDecBinIf.DecodeBin( &uiCode, this.m_cSaoTypeIdxSCModel.Get3( 0, 0, 0 ) );
    if uiCode == 0{
      *ruiVal = 0;
    }else{
      this.m_pcTDecBinIf.DecodeBinEP( &uiCode );
      if uiCode == 0{
        *ruiVal = 5;
      }else{
        *ruiVal = 1;
      }
    }
}
func (this *TDecSbac)   ParseSaoUflc              ( uiLength uint, ruiVal *uint    ){
    this.m_pcTDecBinIf.DecodeBinsEP ( ruiVal, int(uiLength) );
}

func (this *TDecSbac)   CopySaoOneLcuParam(psDst *TLibCommon.SaoLcuParam,  psSrc *TLibCommon.SaoLcuParam){
  var i int;
  psDst.PartIdx = psSrc.PartIdx;
  psDst.TypeIdx = psSrc.TypeIdx;
  if psDst.TypeIdx != -1 {
    psDst.SubTypeIdx = psSrc.SubTypeIdx ;
    psDst.Length  = psSrc.Length;
    for i=0;i<psDst.Length;i++ {
      psDst.Offset[i] = psSrc.Offset[i];
    }
  }else{
    psDst.Length  = 0;
    for i=0;i<TLibCommon.SAO_BO_LEN;i++{
      psDst.Offset[i] = 0;
    }
  }
}

func (this *TDecSbac)   ParseSaoOneLcuInterleaving( rx,  ry int, pSaoParam *TLibCommon.SAOParam, pcCU *TLibCommon.TComDataCU,  iCUAddrInSlice,  iCUAddrUpInSlice int,  allowMergeLeft,  allowMergeUp bool){
  iAddr := int(pcCU.GetAddr());
  var uiSymbol uint;
  for iCompIdx:=0; iCompIdx<3; iCompIdx++{
    pSaoParam.SaoLcuParam[iCompIdx][iAddr].MergeUpFlag    = false;
    pSaoParam.SaoLcuParam[iCompIdx][iAddr].MergeLeftFlag  = false;
    pSaoParam.SaoLcuParam[iCompIdx][iAddr].SubTypeIdx     = 0;
    pSaoParam.SaoLcuParam[iCompIdx][iAddr].TypeIdx        = -1;
    pSaoParam.SaoLcuParam[iCompIdx][iAddr].Offset[0]      = 0;
    pSaoParam.SaoLcuParam[iCompIdx][iAddr].Offset[1]      = 0;
    pSaoParam.SaoLcuParam[iCompIdx][iAddr].Offset[2]      = 0;
    pSaoParam.SaoLcuParam[iCompIdx][iAddr].Offset[3]      = 0;

  }
  if pSaoParam.SaoFlag[0] || pSaoParam.SaoFlag[1]  {
    if rx>0 && iCUAddrInSlice!=0 && allowMergeLeft {
      this.ParseSaoMerge(&uiSymbol); 
      pSaoParam.SaoLcuParam[0][iAddr].MergeLeftFlag = uiSymbol!=0;  
//#ifdef ENC_DEC_TRACE
      this.XReadAeTr(int(uiSymbol), "sao_merge_left_flag", TLibCommon.TRACE_LCU);
//#endif
    }
    if pSaoParam.SaoLcuParam[0][iAddr].MergeLeftFlag==false{
      if (ry > 0) && (iCUAddrUpInSlice>=0) && allowMergeUp {
        this.ParseSaoMerge(&uiSymbol);
        pSaoParam.SaoLcuParam[0][iAddr].MergeUpFlag = uiSymbol!=0;  
//#ifdef ENC_DEC_TRACE
        this.XReadAeTr(int(uiSymbol), "sao_merge_up_flag", TLibCommon.TRACE_LCU);
//#endif
      }
    }
  }

  for iCompIdx:=0; iCompIdx<3; iCompIdx++{
    if (iCompIdx == 0  && pSaoParam.SaoFlag[0]) || (iCompIdx > 0  && pSaoParam.SaoFlag[1]) {
      if rx>0 && iCUAddrInSlice!=0 && allowMergeLeft{
        pSaoParam.SaoLcuParam[iCompIdx][iAddr].MergeLeftFlag = pSaoParam.SaoLcuParam[0][iAddr].MergeLeftFlag;
      }else{
        pSaoParam.SaoLcuParam[iCompIdx][iAddr].MergeLeftFlag = false;
      }

      if pSaoParam.SaoLcuParam[iCompIdx][iAddr].MergeLeftFlag==false{
        if (ry > 0) && (iCUAddrUpInSlice>=0) && allowMergeUp{
          pSaoParam.SaoLcuParam[iCompIdx][iAddr].MergeUpFlag = pSaoParam.SaoLcuParam[0][iAddr].MergeUpFlag;
        }else{
          pSaoParam.SaoLcuParam[iCompIdx][iAddr].MergeUpFlag = false;
        }
        if !pSaoParam.SaoLcuParam[iCompIdx][iAddr].MergeUpFlag{
          pSaoParam.SaoLcuParam[2][iAddr].TypeIdx = pSaoParam.SaoLcuParam[1][iAddr].TypeIdx;
          this.ParseSaoOffset(&(pSaoParam.SaoLcuParam[iCompIdx][iAddr]), uint(iCompIdx));
        }else{
          this.CopySaoOneLcuParam(&pSaoParam.SaoLcuParam[iCompIdx][iAddr], &pSaoParam.SaoLcuParam[iCompIdx][iAddr-pSaoParam.NumCuInWidth]);
        }
      }else{
        this.CopySaoOneLcuParam(&pSaoParam.SaoLcuParam[iCompIdx][iAddr],  &pSaoParam.SaoLcuParam[iCompIdx][iAddr-1]);
      }
    }else{
      pSaoParam.SaoLcuParam[iCompIdx][iAddr].TypeIdx = -1;
      pSaoParam.SaoLcuParam[iCompIdx][iAddr].SubTypeIdx = 0;
    }
  }
}

var iTypeLength = [TLibCommon.MAX_NUM_SAO_TYPE]int{
    TLibCommon.SAO_EO_LEN,
    TLibCommon.SAO_EO_LEN,
    TLibCommon.SAO_EO_LEN,
    TLibCommon.SAO_EO_LEN,
    TLibCommon.SAO_BO_LEN,
};

func (this *TDecSbac)   ParseSaoOffset            (psSaoLcuParam *TLibCommon.SaoLcuParam,  compIdx uint){
  var uiSymbol uint;

  if compIdx==2{
    uiSymbol = uint( psSaoLcuParam.TypeIdx + 1);
  }else{
    this.ParseSaoTypeIdx(&uiSymbol);
//#ifdef ENC_DEC_TRACE
    if compIdx==0{
      this.XReadAeTr(int(uiSymbol), "sao_type_idx_luma", TLibCommon.TRACE_LCU);
    }else{
      this.XReadAeTr(int(uiSymbol), "sao_type_idx_chroma", TLibCommon.TRACE_LCU);
    }
//#endif
  }
  psSaoLcuParam.TypeIdx = int(uiSymbol) - 1;
  if uiSymbol!=0{
    psSaoLcuParam.Length = iTypeLength[psSaoLcuParam.TypeIdx];

    var bitDepth, offsetTh int;
    if compIdx!=0 {
        bitDepth = TLibCommon.G_bitDepthC;
    }else{
        bitDepth = TLibCommon.G_bitDepthY;
    }
    if bitDepth - 5 < 5 {
        offsetTh = 1 << 5;
    }else{
        offsetTh = 1 << uint(bitDepth - 5);
    }

    if psSaoLcuParam.TypeIdx == TLibCommon.SAO_BO {
      for i:=0; i< psSaoLcuParam.Length; i++{
        this.ParseSaoMaxUvlc(&uiSymbol, uint(offsetTh -1) );
        psSaoLcuParam.Offset[i] = int(uiSymbol);
//#ifdef ENC_DEC_TRACE
        this.XReadAeTr(int(uiSymbol), "sao_offset_abs", TLibCommon.TRACE_LCU);
//#endif
      }
      for i:=0; i< psSaoLcuParam.Length; i++{
        if psSaoLcuParam.Offset[i] != 0 {
          this.m_pcTDecBinIf.DecodeBinEP ( &uiSymbol);
//#ifdef ENC_DEC_TRACE
          this.XReadAeTr(int(uiSymbol), "sao_offset_sign", TLibCommon.TRACE_LCU);
//#endif
          if uiSymbol!=0{
            psSaoLcuParam.Offset[i] = -psSaoLcuParam.Offset[i] ;
          }
        }
      }
      this.ParseSaoUflc(5, &uiSymbol );
      psSaoLcuParam.SubTypeIdx = int(uiSymbol);
//#ifdef ENC_DEC_TRACE
      this.XReadAeTr(int(uiSymbol), "sao_band_position", TLibCommon.TRACE_LCU);
//#endif
    }else if psSaoLcuParam.TypeIdx < 4 {
      this.ParseSaoMaxUvlc(&uiSymbol, uint(offsetTh -1) );
      psSaoLcuParam.Offset[0] = int(uiSymbol);
//#ifdef ENC_DEC_TRACE
      this.XReadAeTr(int(uiSymbol), "sao_offset_abs", TLibCommon.TRACE_LCU);
//#endif
      this.ParseSaoMaxUvlc(&uiSymbol, uint(offsetTh -1) );
      psSaoLcuParam.Offset[1] = int(uiSymbol);
//#ifdef ENC_DEC_TRACE
      this.XReadAeTr(int(uiSymbol), "sao_offset_abs", TLibCommon.TRACE_LCU);
//#endif
      this.ParseSaoMaxUvlc(&uiSymbol, uint(offsetTh -1) );
      psSaoLcuParam.Offset[2] = -int(uiSymbol);
//#ifdef ENC_DEC_TRACE
      this.XReadAeTr(int(uiSymbol), "sao_offset_abs", TLibCommon.TRACE_LCU);
//#endif
      this.ParseSaoMaxUvlc(&uiSymbol, uint(offsetTh -1) );
      psSaoLcuParam.Offset[3] = -int(uiSymbol);
//#ifdef ENC_DEC_TRACE
      this.XReadAeTr(int(uiSymbol), "sao_offset_abs", TLibCommon.TRACE_LCU);
//#endif
     if compIdx != 2 {
       this.ParseSaoUflc(2, &uiSymbol );
       psSaoLcuParam.SubTypeIdx = int(uiSymbol);
       psSaoLcuParam.TypeIdx += psSaoLcuParam.SubTypeIdx;
//#ifdef ENC_DEC_TRACE
       if compIdx==0{
         this.XReadAeTr(int(uiSymbol), "sao_eo_class_luma", TLibCommon.TRACE_LCU);
       }else{
         this.XReadAeTr(int(uiSymbol), "sao_eo_class_chroma", TLibCommon.TRACE_LCU);
       }
//#endif
     }
   }
  }else{
    psSaoLcuParam.Length = 0;
  }
}
//private:
func (this *TDecSbac)   xReadUnarySymbol    ( ruiSymbol *uint, pcSCModel []TLibCommon.ContextModel,  iOffset int){
    this.m_pcTDecBinIf.DecodeBin( ruiSymbol, &pcSCModel[0] );

    if *ruiSymbol!=0{
      return;
    }

    uiSymbol := uint(0);
    uiCont := uint(1);

    for uiCont!=0 {
      this.m_pcTDecBinIf.DecodeBin( &uiCont, &pcSCModel[ iOffset ] );
      uiSymbol++;
    }

    *ruiSymbol = uiSymbol;
}
func (this *TDecSbac)   xReadUnaryMaxSymbol ( ruiSymbol *uint, pcSCModel []TLibCommon.ContextModel,  iOffset,  uiMaxSymbol uint ){
    if uiMaxSymbol == 0 {
      *ruiSymbol = 0;
      return;
    }

    this.m_pcTDecBinIf.DecodeBin( ruiSymbol, &pcSCModel[0] );

    if *ruiSymbol == 0 || uiMaxSymbol == 1 {
      return;
    }

    uiSymbol := uint(0);
    uiCont   := uint(1);

    for uiCont!=0 && ( uiSymbol < uiMaxSymbol - 1 ) {
      this.m_pcTDecBinIf.DecodeBin( &uiCont, &pcSCModel[ iOffset ] );
      uiSymbol++;
    }


    if uiCont!=0 && ( uiSymbol == uiMaxSymbol - 1 ) {
      uiSymbol++;
    }

    *ruiSymbol = uiSymbol;
}
func (this *TDecSbac)   xReadEpExGolomb     ( ruiSymbol *uint,  uiCount uint){
    uiSymbol := uint(0);
    uiBit := uint(1);

    for uiBit!=0 {
      this.m_pcTDecBinIf.DecodeBinEP( &uiBit );
      uiSymbol += uiBit << uiCount;
      uiCount++;
    }

    uiCount--;
    if uiCount !=0 {
      var bins uint;
      this.m_pcTDecBinIf.DecodeBinsEP( &bins, int(uiCount) );
      uiSymbol += bins;
    }

    *ruiSymbol = uiSymbol;
}
func (this *TDecSbac)   xReadCoefRemainExGolomb ( rSymbol *uint, rParam uint){
    prefix   := uint(0);
    codeWord := uint(1);
    for codeWord!=0{
      prefix++;
      this.m_pcTDecBinIf.DecodeBinEP( &codeWord );
    }

    codeWord  = 1 - codeWord;
    prefix -= codeWord;
    codeWord=0;
    if prefix < TLibCommon.COEF_REMAIN_BIN_REDUCTION {
      this.m_pcTDecBinIf.DecodeBinsEP(&codeWord, int(rParam));
      *rSymbol = (prefix<<rParam) + codeWord;
    }else{
      this.m_pcTDecBinIf.DecodeBinsEP(&codeWord,int(prefix-TLibCommon.COEF_REMAIN_BIN_REDUCTION+rParam));
      *rSymbol = (((1<<(prefix-TLibCommon.COEF_REMAIN_BIN_REDUCTION))+TLibCommon.COEF_REMAIN_BIN_REDUCTION-1)<<rParam)+codeWord;
    }
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
