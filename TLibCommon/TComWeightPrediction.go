package TLibCommon

import (

)


// ====================================================================================================================
// Class definition
// ====================================================================================================================
/// weighting prediction class
type TComWeightPrediction struct{
    m_wp0	[3]wpScalingParam;
    m_wp1	[3]wpScalingParam;
}

func NewTComWeightPrediction() *TComWeightPrediction{
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