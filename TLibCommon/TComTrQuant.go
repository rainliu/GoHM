package TLibCommon

import ()

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

/*
  QpParam();  
  Void setQpParam( Int qpScaled )
  {
    m_iQP   = qpScaled;
    m_iPer  = qpScaled / 6;
    m_iRem  = qpScaled % 6;
    m_iBits = QP_BITS + m_iPer;
  }

  Void clear()
  {
    m_iQP   = 0;
    m_iPer  = 0;
    m_iRem  = 0;
    m_iBits = 0;
  }


  Int per()   const { return m_iPer; }
  Int rem()   const { return m_iRem; }
  Int bits()  const { return m_iBits; }

  Int qp() {return m_iQP;}
}; // END CLASS DEFINITION QpParam
*/
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
    m_quantCoef              [SCALING_LIST_SIZE_NUM][SCALING_LIST_NUM][SCALING_LIST_REM_NUM]*int     ///< array of quantization matrix coefficient 4x4
    m_dequantCoef            [SCALING_LIST_SIZE_NUM][SCALING_LIST_NUM][SCALING_LIST_REM_NUM]*int     ///< array of dequantization matrix coefficient 4x4
    m_errScale               [SCALING_LIST_SIZE_NUM][SCALING_LIST_NUM][SCALING_LIST_REM_NUM]*float64 ///< array of quantization matrix coefficient 4x4
}

/*
public:
  TComTrQuant();
  ~TComTrQuant();

  // initialize class
  Void init                 ( UInt uiMaxWidth, UInt uiMaxHeight, UInt uiMaxTrSize, Int iSymbolMode = 0, UInt *aTable4 = NULL, UInt *aTable8 = NULL, UInt *aTableLastPosVlcIndex=NULL, Bool useRDOQ = false,  
#if RDOQ_TRANSFORMSKIP
    Bool useRDOQTS = false,  
#endif
    Bool bEnc = false, Bool useTransformSkipFast = false
#if ADAPTIVE_QP_SELECTION
    , Bool bUseAdaptQpSelect = false
#endif 
    );

  // transform & inverse transform functions
  Void transformNxN( TComDataCU* pcCU, 
                     Pel*        pcResidual, 
                     UInt        uiStride, 
                     TCoeff*     rpcCoeff, 
#if ADAPTIVE_QP_SELECTION
                     Int*&       rpcArlCoeff, 
#endif
                     UInt        uiWidth, 
                     UInt        uiHeight, 
                     UInt&       uiAbsSum, 
                     TextType    eTType, 
                     UInt        uiAbsPartIdx,
                     Bool        useTransformSkip = false );

  Void invtransformNxN( Bool transQuantBypass, TextType eText, UInt uiMode,Pel* rpcResidual, UInt uiStride, TCoeff*   pcCoeff, UInt uiWidth, UInt uiHeight,  Int scalingListType, Bool useTransformSkip = false );
  Void invRecurTransformNxN ( TComDataCU* pcCU, TComYuv* pcYuvPred, UInt uiAbsPartIdx, TextType eTxt, Pel* rpcResidual, UInt uiAddr,   UInt uiStride, UInt uiWidth, UInt uiHeight,
                             UInt uiMaxTrMode,  UInt uiTrMode, TCoeff* rpcCoeff );

  // Misc functions
  Void setQPforQuant( Int qpy, TextType eTxtType, Int qpBdOffset, Int chromaQPOffset);

#if RDOQ_CHROMA_LAMBDA 
  Void setLambda(Double dLambdaLuma, Double dLambdaChroma) { m_dLambdaLuma = dLambdaLuma; m_dLambdaChroma = dLambdaChroma; }
  Void selectLambda(TextType eTType) { m_dLambda = (eTType == TEXT_LUMA) ? m_dLambdaLuma : m_dLambdaChroma; }
#else
  Void setLambda(Double dLambda) { m_dLambda = dLambda;}
#endif
  Void setRDOQOffset( UInt uiRDOQOffset ) { m_uiRDOQOffset = uiRDOQOffset; }

  estBitsSbacStruct* m_pcEstBitsSbac;

  static Int      calcPatternSigCtx( const UInt* sigCoeffGroupFlag, UInt posXCG, UInt posYCG, Int width, Int height );

  static Int      getSigCtxInc     (
                                     Int                             patternSigCtx,
                                     UInt                            scanIdx,
                                     Int                             posX,
                                     Int                             posY,
                                     Int                             blockType,
                                     Int                             width
                                    ,Int                             height
                                    ,TextType                        textureType
                                    );
  static UInt getSigCoeffGroupCtxInc  ( const UInt*                   uiSigCoeffGroupFlag,
                                       const UInt                       uiCGPosX,
                                       const UInt                       uiCGPosY,
                                       const UInt                     scanIdx,
                                       Int width, Int height);
  Void initScalingList                      ();
  Void destroyScalingList                   ();
  Void setErrScaleCoeff    ( UInt list, UInt size, UInt qp);
  Double* getErrScaleCoeff ( UInt list, UInt size, UInt qp) {return m_errScale[size][list][qp];};    //!< get Error Scale Coefficent
  Int* getQuantCoeff       ( UInt list, UInt qp, UInt size) {return m_quantCoef[size][list][qp];};   //!< get Quant Coefficent
  Int* getDequantCoeff     ( UInt list, UInt qp, UInt size) {return m_dequantCoef[size][list][qp];}; //!< get DeQuant Coefficent
  Void setUseScalingList   ( Bool bUseScalingList){ m_scalingListEnabledFlag = bUseScalingList; };
  Bool getUseScalingList   (){ return m_scalingListEnabledFlag; };
  Void setFlatScalingList  ();
  Void xsetFlatScalingList ( UInt list, UInt size, UInt qp);
  Void xSetScalingListEnc  ( TComScalingList *scalingList, UInt list, UInt size, UInt qp);
  Void xSetScalingListDec  ( TComScalingList *scalingList, UInt list, UInt size, UInt qp);
  Void setScalingList      ( TComScalingList *scalingList);
  Void setScalingListDec   ( TComScalingList *scalingList);
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
