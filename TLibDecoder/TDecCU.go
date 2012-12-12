package TLibDecoder

import (
    "gohm/TLibCommon"
)

// ====================================================================================================================
// Class definition
// ====================================================================================================================

/// CU decoder class
type TDecCu struct {
    //private:
    m_uiMaxDepth uint ///< max. number of depth
    //TComYuv**           m_ppcYuvResi;       ///< array of residual buffer
    //TComYuv**           m_ppcYuvReco;       ///< array of prediction & reconstruction buffer
    //TComDataCU**        m_ppcCU;            ///< CU data array

    // access channel
    m_pcTrQuant        *TLibCommon.TComTrQuant
    m_pcPrediction     *TLibCommon.TComPrediction
    m_pcEntropyDecoder *TDecEntropy

    m_bDecodeDQP bool
}

func NewTDecCu() *TDecCu{
	return &TDecCu{};
}

  /// initialize access channels
func (this *TDecCu) Init  ( pcEntropyDecoder *TDecEntropy, pcTrQuant *TLibCommon.TComTrQuant, pcPrediction *TLibCommon.TComPrediction){
}

  /// create internal buffers
func (this *TDecCu) Create  (  uiMaxDepth,  uiMaxWidth,  uiMaxHeight uint){
}

/// destroy internal buffers
func (this *TDecCu) Destroy() {
}

 
  /// decode CU information
func (this *TDecCu)  DecodeCU                ( pcCU *TLibCommon.TComDataCU, ruiIsLast *uint){
}

  /// reconstruct CU information
func (this *TDecCu)  DecompressCU            ( pcCU *TLibCommon.TComDataCU ){
}


func (this *TDecCu)  xDecodeCU               ( pcCU *TLibCommon.TComDataCU,                        uiAbsPartIdx,  uiDepth uint,  ruiIsLast *uint){
}
func (this *TDecCu)  xFinishDecodeCU         ( pcCU *TLibCommon.TComDataCU,                        uiAbsPartIdx,  uiDepth uint,  ruiIsLast *uint){
}
func (this *TDecCu)  xDecodeSliceEnd         ( pcCU *TLibCommon.TComDataCU,                        uiAbsPartIdx,  uiDepth uint) bool{
	return true
}
func (this *TDecCu)  xDecompressCU           ( pcCU *TLibCommon.TComDataCU, pcCUCur *TLibCommon.TComDataCU,   uiAbsPartIdx,  uiDepth uint){
}

func (this *TDecCu)  xReconInter             ( pcCU *TLibCommon.TComDataCU, uiAbsPartIdx,  uiDepth uint ){
}

func (this *TDecCu)  xReconIntraQT           ( pcCU *TLibCommon.TComDataCU, uiAbsPartIdx,  uiDepth uint ){
}
func (this *TDecCu)  xIntraRecLumaBlk        ( pcCU *TLibCommon.TComDataCU, uiTrDepth,  uiAbsPartIdx uint, pcRecoYuv *TLibCommon.TComYuv, pcPredYuv *TLibCommon.TComYuv, pcResiYuv *TLibCommon.TComYuv){
}
func (this *TDecCu)  xIntraRecChromaBlk      ( pcCU *TLibCommon.TComDataCU, uiTrDepth,  uiAbsPartIdx uint, pcRecoYuv *TLibCommon.TComYuv, pcPredYuv *TLibCommon.TComYuv, pcResiYuv *TLibCommon.TComYuv,  uiChromaId uint){
}

func (this *TDecCu)  xReconPCM               ( pcCU *TLibCommon.TComDataCU, uiAbsPartIdx,  uiDepth uint ){
}

func (this *TDecCu)  xDecodeInterTexture     ( pcCU *TLibCommon.TComDataCU, uiAbsPartIdx,  uiDepth uint ){
}
func (this *TDecCu)  xDecodePCMTexture       ( pcCU *TLibCommon.TComDataCU,  uiPartIdx uint, piPCM *TLibCommon.Pel, piReco *TLibCommon.Pel,  uiStride,  uiWidth,  uiHeight uint,  ttText TLibCommon.TextType){
}

func (this *TDecCu)  xCopyToPic              ( pcCU *TLibCommon.TComDataCU, pcPic *TLibCommon.TComPic,  uiZorderIdx,  uiDepth uint){
}

func (this *TDecCu)  xIntraLumaRecQT         ( pcCU *TLibCommon.TComDataCU,  uiTrDepth,  uiAbsPartIdx uint, pcRecoYuv *TLibCommon.TComYuv, pcPredYuv *TLibCommon.TComYuv, pcResiYuv *TLibCommon.TComYuv ){
}
func (this *TDecCu)  xIntraChromaRecQT       ( pcCU *TLibCommon.TComDataCU,  uiTrDepth,  uiAbsPartIdx uint, pcRecoYuv *TLibCommon.TComYuv, pcPredYuv *TLibCommon.TComYuv, pcResiYuv *TLibCommon.TComYuv ){
}

func (this *TDecCu) GetdQPFlag               ()         bool               { 
	return this.m_bDecodeDQP;        
}
func (this *TDecCu) SetdQPFlag               (  b bool)                { 
	this.m_bDecodeDQP = b;           
}
func (this *TDecCu) xFillPCMBuffer           ( pCU *TLibCommon.TComDataCU,  absPartIdx,  depth uint){
}
