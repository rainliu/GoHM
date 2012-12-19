package TLibCommon

import ()

// ====================================================================================================================
// Class definition
// ====================================================================================================================

/// prediction class
type TComPrediction struct {
    TComWeightPrediction
    //protected:
    m_piYuvExt     []int
    m_iYuvExtStride int
    m_iYuvExtHeight int

    m_acYuvPred        [2]TComYuv
    m_cYuvPredTemp     TComYuv
    m_filteredBlock    [4][4]TComYuv
    m_filteredBlockTmp [4]TComYuv

    m_if TComInterpolationFilter

    m_pLumaRecBuffer *Pel ///< array for downsampled reconstructed luma sample
    m_iLumaRecStride int  ///< stride of #m_pLumaRecBuffer array
}

/*
  Void xPredIntraAng            (Int bitDepth, Int* pSrc, Int srcStride, Pel*& rpDst, Int dstStride, UInt width, UInt height, UInt dirMode, Bool blkAboveAvailable, Bool blkLeftAvailable, Bool bFilter );
  Void xPredIntraPlanar         ( Int* pSrc, Int srcStride, Pel* rpDst, Int dstStride, UInt width, UInt height );

  // motion compensation functions
  Void xPredInterUni            ( TComDataCU* pcCU,                          UInt uiPartAddr,               Int iWidth, Int iHeight, RefPicList eRefPicList, TComYuv*& rpcYuvPred, Int iPartIdx, Bool bi=false          );
  Void xPredInterBi             ( TComDataCU* pcCU,                          UInt uiPartAddr,               Int iWidth, Int iHeight,                         TComYuv*& rpcYuvPred, Int iPartIdx          );
  Void xPredInterLumaBlk  ( TComDataCU *cu, TComPicYuv *refPic, UInt partAddr, TComMv *mv, Int width, Int height, TComYuv *&dstPic, Bool bi );
  Void xPredInterChromaBlk( TComDataCU *cu, TComPicYuv *refPic, UInt partAddr, TComMv *mv, Int width, Int height, TComYuv *&dstPic, Bool bi );
  Void xWeightedAverage         ( TComDataCU* pcCU, TComYuv* pcYuvSrc0, TComYuv* pcYuvSrc1, Int iRefIdx0, Int iRefIdx1, UInt uiPartAddr, Int iWidth, Int iHeight, TComYuv*& rpcYuvDst );

  Void xGetLLSPrediction ( TComPattern* pcPattern, Int* pSrc0, Int iSrcStride, Pel* pDst0, Int iDstStride, UInt uiWidth, UInt uiHeight, UInt uiExt0 );

  Void xDCPredFiltering( Int* pSrc, Int iSrcStride, Pel*& rpDst, Int iDstStride, Int iWidth, Int iHeight );
  Bool xCheckIdenticalMotion    ( TComDataCU* pcCU, UInt PartAddr);
*/
func NewTComPrediction() *TComPrediction{
	return &TComPrediction{};
}

func (this *TComPrediction) InitTempBuff(){
}

  // inter
func (this *TComPrediction) MotionCompensation         ( pcCU *TComDataCU, pcYuvPred *TComYuv,  eRefPicList RefPicList,  iPartIdx int ){
}

  // motion vector prediction
func (this *TComPrediction) GetMvPredAMVP              ( pcCU *TComDataCU,  uiPartIdx,  uiPartAddr uint,  eRefPicList RefPicList,  iRefIdx int, rcMvPred *TComMv ){
}

  // Angular Intra
func (this *TComPrediction) PredIntraLumaAng           ( pcTComPattern *TComPattern,  uiDirMode uint, piPred []Pel,  uiStride uint,  iWidth,  iHeight int,  pcCU *TComDataCU,  bAbove,  bLeft bool){
}
func (this *TComPrediction) PredIntraChromaAng         ( pcTComPattern *TComPattern, piSrc []int,  uiDirMode uint, piPred []Pel,  uiStride uint,  iWidth,  iHeight int, pcCU *TComDataCU,  bAbove,  bLeft bool){
}

func (this *TComPrediction) PredIntraGetPredValDC      ( pSrc []int,  iSrcStride int,  iWidth,  iHeight uint,  bAbove,  bLeft bool) Pel{
	return 0;
}

func (this *TComPrediction) GetPredicBuf()         []int   { 
	return this.m_piYuvExt;      
}
func (this *TComPrediction) GetPredicBufWidth()     int   { 
	return this.m_iYuvExtStride; 
}
func (this *TComPrediction) GetPredicBufHeight()    int   { 
	return this.m_iYuvExtHeight; 
}


// ====================================================================================================================
// Class definition
// ====================================================================================================================
/// weighting prediction class
type TComWeightPrediction struct {
    m_wp0 [3]wpScalingParam
    m_wp1 [3]wpScalingParam
}

func NewTComWeightPrediction() *TComWeightPrediction {
    return &TComWeightPrediction{}
}

/*
func (this *TComWeightPrediction)  GetWpScaling(TComDataCU*  pcCU , Int iRefIdx0, Int iRefIdx1, wpScalingParam *&wp0 , wpScalingParam *&wp1){
}

func (this *TComWeightPrediction)  AddWeightBi( TComYuv* pcYuvSrc0, TComYuv* pcYuvSrc1, UInt iPartUnitIdx, UInt iWidth, UInt iHeight, wpScalingParam *wp0, wpScalingParam *wp1, TComYuv* rpcYuvDst, Bool bRound=true ){
}
func (this *TComWeightPrediction)  AddWeightUni( TComYuv* pcYuvSrc0, UInt iPartUnitIdx, UInt iWidth, UInt iHeight, wpScalingParam *wp0, TComYuv* rpcYuvDst ){
}

func (this *TComWeightPrediction)  xWeightedPredictionUni( TComDataCU* pcCU, TComYuv* pcYuvSrc, UInt uiPartAddr, Int iWidth, Int iHeight, RefPicList eRefPicList, TComYuv*& rpcYuvPred, Int iPartIdx, Int iRefIdx=-1 ){
}
func (this *TComWeightPrediction)  xWeightedPredictionBi( TComDataCU* pcCU, TComYuv* pcYuvSrc0, TComYuv* pcYuvSrc1, Int iRefIdx0, Int iRefIdx1, UInt uiPartIdx, Int iWidth, Int iHeight, TComYuv* rpcYuvDst ){
}
*/
