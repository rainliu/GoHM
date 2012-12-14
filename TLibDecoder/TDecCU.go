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
    m_uiMaxDepth 	uint 			  ///< max. number of depth
    m_ppcYuvResi 	[]*TLibCommon.TComYuv;       ///< array of residual buffer
    m_ppcYuvReco	[]*TLibCommon.TComYuv;       ///< array of prediction & reconstruction buffer
    m_ppcCU		    []*TLibCommon.TComDataCU;    ///< CU data array

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
  this.m_pcEntropyDecoder  = pcEntropyDecoder;
  this.m_pcTrQuant         = pcTrQuant;
  this.m_pcPrediction      = pcPrediction;
}

  /// create internal buffers
func (this *TDecCu) Create  (  uiMaxDepth,  uiMaxWidth,  uiMaxHeight uint){
  this.m_uiMaxDepth = uiMaxDepth+1;
  
  this.m_ppcYuvResi = make([]*TLibCommon.TComYuv, 	  this.m_uiMaxDepth-1);
  this.m_ppcYuvReco = make([]*TLibCommon.TComYuv,	  this.m_uiMaxDepth-1);
  this.m_ppcCU      = make([]*TLibCommon.TComDataCU, this.m_uiMaxDepth-1);
  
  var uiNumPartitions uint;
  for ui := uint(0); ui < this.m_uiMaxDepth-1; ui++ {
    uiNumPartitions = 1<<( ( this.m_uiMaxDepth - ui - 1 )<<1 );
    uiWidth  := uiMaxWidth  >> ui;
    uiHeight := uiMaxHeight >> ui;
    
    this.m_ppcYuvResi[ui] = TLibCommon.NewTComYuv();    
    this.m_ppcYuvResi[ui].Create( uiWidth, uiHeight );
    this.m_ppcYuvReco[ui] = TLibCommon.NewTComYuv();    
    this.m_ppcYuvReco[ui].Create( uiWidth, uiHeight );
    this.m_ppcCU     [ui] = TLibCommon.NewTComDataCU(); 
    this.m_ppcCU     [ui].Create( uiNumPartitions, uiWidth, uiHeight, true, int(uiMaxWidth >> (this.m_uiMaxDepth - 1)), false );
  }
  
  this.m_bDecodeDQP = false;

  // initialize partition order.
  piTmp := uint(0);
  TLibCommon.InitZscanToRaster(int(this.m_uiMaxDepth), 1, 0, TLibCommon.G_auiZscanToRaster[:], &piTmp);
  TLibCommon.InitRasterToZscan( uiMaxWidth, uiMaxHeight, this.m_uiMaxDepth );
  
  // initialize conversion matrix from partition index to pel
  TLibCommon.InitRasterToPelXY( uiMaxWidth, uiMaxHeight, this.m_uiMaxDepth );
//#if !LINEBUF_CLEANUP
//  initMotionReferIdx ( uiMaxWidth, uiMaxHeight, this.m_uiMaxDepth );
//#endif
}

/// destroy internal buffers
func (this *TDecCu) Destroy() {
  for ui := uint(0); ui < this.m_uiMaxDepth-1; ui++ {
    this.m_ppcYuvResi[ui].Destroy(); 
    //delete m_ppcYuvResi[ui]; 
    this.m_ppcYuvResi[ui] = nil;
    this.m_ppcYuvReco[ui].Destroy(); 
    //delete m_ppcYuvReco[ui]; 
    this.m_ppcYuvReco[ui] = nil;
    this.m_ppcCU     [ui].Destroy(); 
    //delete m_ppcCU     [ui]; 
    this.m_ppcCU     [ui] = nil;
  }
  
  //delete [] m_ppcYuvResi; 
  this.m_ppcYuvResi = nil;
  //delete [] m_ppcYuvReco; 
  this.m_ppcYuvReco = nil;
  //delete [] m_ppcCU     ; 
  this.m_ppcCU      = nil;
}

 
  /// decode CU information
func (this *TDecCu)  DecodeCU                ( pcCU *TLibCommon.TComDataCU, ruiIsLast *uint){
  if pcCU.GetSlice().GetPPS().GetUseDQP() {
    this.SetdQPFlag(true);
  }

//#if !REMOVE_BURST_IPCM
//  pcCU->setNumSucIPCM(0);
//#endif

  // start from the top level CU
  this.xDecodeCU( pcCU, 0, 0, ruiIsLast);
}

  /// reconstruct CU information
func (this *TDecCu)  DecompressCU            ( pcCU *TLibCommon.TComDataCU ){
	this.xDecompressCU( pcCU, pcCU, 0,  0 );
}


func (this *TDecCu)  xDecodeCU               ( pcCU *TLibCommon.TComDataCU,                        uiAbsPartIdx,  uiDepth uint,  ruiIsLast *uint){
}
func (this *TDecCu)  xFinishDecodeCU         ( pcCU *TLibCommon.TComDataCU,                        uiAbsPartIdx,  uiDepth uint,  ruiIsLast *uint){
}
func (this *TDecCu)  xDecodeSliceEnd         ( pcCU *TLibCommon.TComDataCU,                        uiAbsPartIdx,  uiDepth uint) bool{
  var uiIsLast uint;
  pcPic := pcCU.GetPic();
  pcSlice := pcPic.GetSlice(pcPic.GetCurrSliceIdx());
  uiCurNumParts := pcPic.GetNumPartInCU() >> (uiDepth<<1);
  uiWidth := pcSlice.GetSPS().GetPicWidthInLumaSamples();
  uiHeight := pcSlice.GetSPS().GetPicHeightInLumaSamples();
  uiGranularityWidth := TLibCommon.G_uiMaxCUWidth;
  uiPosX := pcCU.GetCUPelX() + TLibCommon.G_auiRasterToPelX[ TLibCommon.G_auiZscanToRaster[uiAbsPartIdx] ];
  uiPosY := pcCU.GetCUPelY() + TLibCommon.G_auiRasterToPelY[ TLibCommon.G_auiZscanToRaster[uiAbsPartIdx] ];

  if ((uiPosX+uint(pcCU.GetWidth1 (uiAbsPartIdx)))%uiGranularityWidth==0||(uiPosX+uint(pcCU.GetWidth1 (uiAbsPartIdx))==uiWidth)) &&
     ((uiPosY+uint(pcCU.GetHeight1(uiAbsPartIdx)))%uiGranularityWidth==0||(uiPosY+uint(pcCU.GetHeight1(uiAbsPartIdx))==uiHeight)) {
    this.m_pcEntropyDecoder.DecodeTerminatingBit( &uiIsLast );
  }else{
    uiIsLast=0;
  }
  
  if uiIsLast!=0 {
    if pcSlice.IsNextDependentSlice()&&!pcSlice.IsNextSlice() {
      pcSlice.SetDependentSliceCurEndCUAddr(pcCU.GetSCUAddr()+uiAbsPartIdx+uiCurNumParts);
    }else{
      pcSlice.SetSliceCurEndCUAddr(pcCU.GetSCUAddr()+uiAbsPartIdx+uiCurNumParts);
      pcSlice.SetDependentSliceCurEndCUAddr(pcCU.GetSCUAddr()+uiAbsPartIdx+uiCurNumParts);
    }
  }

  return uiIsLast>0;
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
