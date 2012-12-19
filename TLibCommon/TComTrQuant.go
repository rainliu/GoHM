package TLibCommon

import (
	"math"
)

// ====================================================================================================================
// Constants
// ====================================================================================================================

const QP_BITS = 15

// ====================================================================================================================
// Type definition
// ====================================================================================================================

type estBitsSbacStruct struct {
    significantCoeffGroupBits [NUM_SIG_CG_FLAG_CTX][2]int
    significantBits           [NUM_SIG_FLAG_CTX][2]int
    lastXBits                 [32]int
    lastYBits                 [32]int
    m_greaterOneBits          [NUM_ONE_FLAG_CTX][2]int
    m_levelAbsBits            [NUM_ABS_FLAG_CTX][2]int

    blockCbpBits     [3 * NUM_QT_CBF_CTX][2]int
    blockRootCbpBits [4][2]int
    scanZigzag       [2]int ///< flag for zigzag scan
    scanNonZigzag    [2]int ///< flag for non zigzag scan
}

// ====================================================================================================================
// Class definition
// ====================================================================================================================

/// QP class
type QpParam struct {
    m_iQP   int
    m_iPer  int
    m_iRem  int
    m_iBits int
}

//public:


func NewQpParam() *QpParam{ 
	return &QpParam{};
}

func (this *QpParam) SetQpParam(  qpScaled int) {
    this.m_iQP   = qpScaled;
    this.m_iPer  = qpScaled / 6;
    this.m_iRem  = qpScaled % 6;
    this.m_iBits = QP_BITS + this.m_iPer;
}

func (this *QpParam) Clear(){
    this.m_iQP   = 0;
    this.m_iPer  = 0;
    this.m_iRem  = 0;
    this.m_iBits = 0;
}


func (this *QpParam)  GetPer()   int { 
	return this.m_iPer; 
}
func (this *QpParam)  GetRem()   int { 
	return this.m_iRem;
}
func (this *QpParam)  GetBits()  int { 
	return this.m_iBits; 
}

func (this *QpParam)  GetQp() int{
	return this.m_iQP;
}


/// transform and quantization class
type TComTrQuant struct {
    //protected:
    //#if ADAPTIVE_QP_SELECTION
    m_qpDelta       [MAX_QP + 1]int
    m_sliceNsamples [LEVEL_RANGE + 1]int
    m_sliceSumC     [LEVEL_RANGE + 1]float64
    //#endif
    m_plTempCoeff *int

    m_cQP QpParam
    //#if RDOQ_CHROMA_LAMBDA
    m_dLambdaLuma   float64
    m_dLambdaChroma float64
    //#endif
    m_dLambda      float64
    m_uiRDOQOffset uint
    m_uiMaxTrSize  uint
    m_bEnc         bool
    m_useRDOQ      bool
    //#if RDOQ_TRANSFORMSKIP
    m_useRDOQTS bool
    //#endif
    //#if ADAPTIVE_QP_SELECTION
    m_bUseAdaptQpSelect bool
    //#endif
    m_useTransformSkipFast   bool
    m_scalingListEnabledFlag bool
    m_quantCoef              [SCALING_LIST_SIZE_NUM][SCALING_LIST_NUM][SCALING_LIST_REM_NUM][]int     ///< array of quantization matrix coefficient 4x4
    m_dequantCoef            [SCALING_LIST_SIZE_NUM][SCALING_LIST_NUM][SCALING_LIST_REM_NUM][]int     ///< array of dequantization matrix coefficient 4x4
    m_errScale               [SCALING_LIST_SIZE_NUM][SCALING_LIST_NUM][SCALING_LIST_REM_NUM][]float64 ///< array of quantization matrix coefficient 4x4
	m_pcEstBitsSbac	*estBitsSbacStruct;
}

func NewTComTrQuant() *TComTrQuant{
	return &TComTrQuant{}
}

  // initialize class
func (this *TComTrQuant) Init ( uiMaxWidth, uiMaxHeight, uiMaxTrSize uint, iSymbolMode int, 
	aTable4 []uint, aTable8 []uint, aTableLastPosVlcIndex []uint, useRDOQ bool,  
//#if RDOQ_TRANSFORMSKIP
    useRDOQTS bool,  
//#endif
    bEnc bool, useTransformSkipFast bool,
//#if ADAPTIVE_QP_SELECTION
    bUseAdaptQpSelect bool ){
//#endif 
}

  // transform & inverse transform functions
func (this *TComTrQuant)  transformNxN( pcCU *TComDataCU, 
                     pcResidual *Pel, 
                     uiStride	uint, 
                     rpcCoeff *TCoeff, 
//#if ADAPTIVE_QP_SELECTION
                     rpcArlCoeff *int, 
//#endif
                           uiWidth uint, 
                             uiHeight uint, 
                            uiAbsSum *uint, 
                        eTType TextType, 
                             uiAbsPartIdx uint,
                             useTransformSkip bool ){
}

func (this *TComTrQuant)  InvtransformNxN( transQuantBypass bool,  eText TextType,  uiMode uint, rpcResidual *Pel,  uiStride uint, pcCoeff *TCoeff,  uiWidth,  uiHeight uint,   scalingListType int,  useTransformSkip bool ){
}
func (this *TComTrQuant)  InvRecurTransformNxN ( pcCU *TComDataCU, pcYuvPred *TComYuv,  uiAbsPartIdx uint,  eTxt TextType, rpcResidual *Pel,  uiAddr, uiStride,  uiWidth,  uiHeight, uiMaxTrMode, uiTrMode uint, rpcCoeff *TCoeff){
}

  // Misc functions
func (this *TComTrQuant)   SetQPforQuant(  qpy int,  eTxtType TextType,  qpBdOffset,  chromaQPOffset int){
}

//#if RDOQ_CHROMA_LAMBDA 
func (this *TComTrQuant)   SetLambda( dLambdaLuma,  dLambdaChroma float64) { 
	this.m_dLambdaLuma = dLambdaLuma;
	this.m_dLambdaChroma = dLambdaChroma; 
}
func (this *TComTrQuant)   SelectLambda( eTType TextType) { 
	if 	eTType == TEXT_LUMA {
		this.m_dLambda = this.m_dLambdaLuma;
	}else{
		this.m_dLambda = this.m_dLambdaChroma; 
	}
}
//#else
//  Void setLambda(Double dLambda) { m_dLambda = dLambda;}
//#endif
func (this *TComTrQuant)   SetRDOQOffset(  uiRDOQOffset uint) { 
	this.m_uiRDOQOffset = uiRDOQOffset; 
}

func CalcPatternSigCtx( sigCoeffGroupFlag []uint,  posXCG,  posYCG uint,  width,  height int) int{
	return 0;
}

func GetSigCtxInc     (
                                                                  patternSigCtx int,
                                                                 scanIdx uint,
                                                                  posX int,
                                                                  posY int,
                                                                  blockType int,
                                                                  width int,
                                                                 height int,
                                                            textureType  TextType) int{
	return 0;
}
func GetSigCoeffGroupCtxInc  ( uiSigCoeffGroupFlag []uint,
                                       uiCGPosX uint,
                                       uiCGPosY uint,
                                       scanIdx uint,
                                        width,  height int) uint{
   return 0;
}
func (this *TComTrQuant)  InitScalingList                      (){
}
func (this *TComTrQuant)  DestroyScalingList                   (){
}

func (this *TComTrQuant)  SetErrScaleCoeff    	   ( list, size, qp uint ){
  uiLog2TrSize := int(G_aucConvertToBit[ G_scalingListSizeX[size] ]) + 2;
  var bitDepth int;
  if size < SCALING_LIST_32x32 && list != 0 && list != 3 {
  	bitDepth =  G_bitDepthC;
  }else{
  	bitDepth =  G_bitDepthY;
  }
  iTransformShift := MAX_TR_DYNAMIC_RANGE - bitDepth - uiLog2TrSize;  // Represents scaling through forward transform

  uiMaxNumCoeff  := G_scalingListSize[size];
  piQuantcoeff   := this.GetQuantCoeff(list, qp,size);
  pdErrScale     := this.GetErrScaleCoeff(list, size, qp);

  dErrScale := float64(1<<SCALE_BITS);                              // Compensate for scaling of bitcount in Lagrange cost function
  dErrScale = dErrScale*math.Pow(2.0,-2.0*float64(iTransformShift));                     // Compensate for scaling through forward transform
  for i:=uint(0);i<uiMaxNumCoeff;i++ {
  	a := 1<<uint(2*(bitDepth-8))
    pdErrScale[i] = dErrScale / float64(piQuantcoeff[i]) / float64(piQuantcoeff[i]) / float64(a);//DISTORTION_PRECISION_ADJUSTMENT
  }
}
func (this *TComTrQuant)  GetErrScaleCoeff ( list, size, qp uint) []float64{
	return this.m_errScale[size][list][qp];
}    //!< get Error Scale Coefficent
func (this *TComTrQuant)  GetQuantCoeff       ( list, qp, size uint) []int{
	return this.m_quantCoef[size][list][qp];
}   //!< get Quant Coefficent
func (this *TComTrQuant)  GetDequantCoeff     ( list, qp, size uint) []int{
	return this.m_dequantCoef[size][list][qp];
} //!< get DeQuant Coefficent
func (this *TComTrQuant)  SetUseScalingList   ( bUseScalingList bool){ 
	this.m_scalingListEnabledFlag = bUseScalingList; 
}
func (this *TComTrQuant)  GetUseScalingList   () bool{ 
	return this.m_scalingListEnabledFlag; 
}
func (this *TComTrQuant)  SetFlatScalingList  (){
}
func (this *TComTrQuant)  xSetFlatScalingList ( list, size, qp uint){
}

func (this *TComTrQuant) xSetScalingListEnc  ( scalingList *TComScalingList, list, size, qp uint){
}
func (this *TComTrQuant) xSetScalingListDec  ( scalingList *TComScalingList, list, size, qp uint){
}

func (this *TComTrQuant) SetScalingList      ( scalingList *TComScalingList){
  var size,list,qp uint;

  for size=0;size<SCALING_LIST_SIZE_NUM;size++ {
    for list = 0; list < G_scalingListNum[size]; list++ {
      for qp=0;qp<SCALING_LIST_REM_NUM;qp++ {
        this.xSetScalingListEnc(scalingList,list,size,qp);
        this.xSetScalingListDec(scalingList,list,size,qp);
        this.SetErrScaleCoeff(list,size,qp);
      }
    }
  }
}
func (this *TComTrQuant) SetScalingListDec   ( scalingList *TComScalingList){
  var size,list,qp uint;

  for size=0;size<SCALING_LIST_SIZE_NUM;size++ {
    for list = 0; list < G_scalingListNum[size]; list++ {
      for qp=0;qp<SCALING_LIST_REM_NUM;qp++ {
        this.xSetScalingListDec(scalingList,list,size,qp);
      }
    }
  }
}
/*
  Void processScalingListEnc( Int *coeff, Int *quantcoeff, Int quantScales, UInt height, UInt width, UInt ratio, Int sizuNum, UInt dc);
  Void processScalingListDec( Int *coeff, Int *dequantcoeff, Int invQuantScales, UInt height, UInt width, UInt ratio, Int sizuNum, UInt dc);
#if ADAPTIVE_QP_SELECTION
  Void    initSliceQpDelta() ;
  Void    storeSliceQpNext(TComSlice* pcSlice);
  Void    clearSliceARLCnt();
  Int     getQpDelta(Int qp) { return m_qpDelta[qp]; } 
  Int*    getSliceNSamples(){ return m_sliceNsamples ;} 
  Double* getSliceSumC()    { return m_sliceSumC; }
#endif
private:
  // forward Transform
  Void xT   (Int bitDepth, UInt uiMode,Pel* pResidual, UInt uiStride, Int* plCoeff, Int iWidth, Int iHeight );

  // skipping Transform
  Void xTransformSkip (Int bitDepth, Pel* piBlkResi, UInt uiStride, Int* psCoeff, Int width, Int height );

  Void signBitHidingHDQ( TComDataCU* pcCU, TCoeff* pQCoef, TCoeff* pCoef, UInt const *scan, Int* deltaU, Int width, Int height );

  // quantization
  Void xQuant( TComDataCU* pcCU, 
               Int*        pSrc, 
               TCoeff*     pDes, 
#if ADAPTIVE_QP_SELECTION
               Int*&       pArlDes,
#endif
               Int         iWidth, 
               Int         iHeight, 
               UInt&       uiAcSum, 
               TextType    eTType, 
               UInt        uiAbsPartIdx );

  // RDOQ functions

  Void           xRateDistOptQuant ( TComDataCU*                     pcCU,
                                     Int*                            plSrcCoeff,
                                     TCoeff*                         piDstCoeff,
#if ADAPTIVE_QP_SELECTION
                                     Int*&                           piArlDstCoeff,
#endif
                                     UInt                            uiWidth,
                                     UInt                            uiHeight,
                                     UInt&                           uiAbsSum,
                                     TextType                        eTType,
                                     UInt                            uiAbsPartIdx );
__inline UInt              xGetCodedLevel  ( Double&                         rd64CodedCost,
                                             Double&                         rd64CodedCost0,
                                             Double&                         rd64CodedCostSig,
                                             Int                             lLevelDouble,
                                             UInt                            uiMaxAbsLevel,
                                             UShort                          ui16CtxNumSig,
                                             UShort                          ui16CtxNumOne,
                                             UShort                          ui16CtxNumAbs,
                                             UShort                          ui16AbsGoRice,
                                             UInt                            c1Idx,  
                                             UInt                            c2Idx,  
                                             Int                             iQBits,
                                             Double                          dTemp,
                                             Bool                            bLast        ) const;
  __inline Double xGetICRateCost   ( UInt                            uiAbsLevel,
                                     UShort                          ui16CtxNumOne,
                                     UShort                          ui16CtxNumAbs,
                                     UShort                          ui16AbsGoRice 
                                   , UInt                            c1Idx,
                                     UInt                            c2Idx
                                     ) const;
__inline Int xGetICRate  ( UInt                            uiAbsLevel,
                           UShort                          ui16CtxNumOne,
                           UShort                          ui16CtxNumAbs,
                           UShort                          ui16AbsGoRice
                         , UInt                            c1Idx,
                           UInt                            c2Idx
                         ) const;
  __inline Double xGetRateLast     ( const UInt                      uiPosX,
                                     const UInt                      uiPosY,
                                     const UInt                      uiBlkWdth     ) const;
  __inline Double xGetRateSigCoeffGroup (  UShort                    uiSignificanceCoeffGroup,
                                     UShort                          ui16CtxNumSig ) const;
  __inline Double xGetRateSigCoef (  UShort                          uiSignificance,
                                     UShort                          ui16CtxNumSig ) const;
  __inline Double xGetICost        ( Double                          dRate         ) const; 
  __inline Double xGetIEPRate      (                                               ) const;


  // dequantization
  Void xDeQuant(Int bitDepth, const TCoeff* pSrc, Int* pDes, Int iWidth, Int iHeight, Int scalingListType );

  // inverse transform
  Void xIT    (Int bitDepth, UInt uiMode, Int* plCoef, Pel* pResidual, UInt uiStride, Int iWidth, Int iHeight );

  // inverse skipping transform
  Void xITransformSkip (Int bitDepth, Int* plCoef, Pel* pResidual, UInt uiStride, Int width, Int height );
*/
