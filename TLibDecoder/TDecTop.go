package TLibDecoder

import (
	"container/list"
	"gohm/TLibCommon"
)

/// decoder class
type TDecTop struct{
//private:
  m_iGopSize		int;
  m_bGopSizeSet		bool;
  m_iMaxRefPicNum	int;
  
  m_bRefreshPending	bool;    ///< refresh pending flag
  m_pocCRA			int;            ///< POC number of the latest CRA picture
  m_prevRAPisBLA	bool;      ///< true if the previous RAP (CRA/CRANT/BLA/BLANT/IDR) picture is a BLA/BLANT picture
  m_pocRandomAccess	int;   ///< POC number of the random access point (the first IDR or CRA picture)

  m_cListPic		*list.List;         //  Dynamic buffer
  m_parameterSetManagerDecoder				ParameterSetManagerDecoder;  // storage for parameter sets 
  m_apcSlicePilot	*TLibCommon.TComSlice;
  
  //SEImessages *m_SEIs; ///< "all" SEI messages.  If not NULL, we own the object.

  // functional classes
 /* 
  TComPrediction          m_cPrediction;
  TComTrQuant             m_cTrQuant;
  TDecGop                 m_cGopDecoder;
  TDecSlice               m_cSliceDecoder;
  TDecCu                  m_cCuDecoder;
  TDecEntropy             m_cEntropyDecoder;
  TDecCavlc               m_cCavlcDecoder;
  TDecSbac                m_cSbacDecoder;
  TDecBinCABAC            m_cBinCABAC;
  SEIReader               m_seiReader;
  TComLoopFilter          m_cLoopFilter;
  TComSampleAdaptiveOffset m_cSAO;
 */

  m_pcPic		*TLibCommon.TComPic;
  m_uiSliceIdx	uint;
  m_prevPOC		int;
  m_bFirstSliceInPicture	bool;
  m_bFirstSliceInSequence	bool;
}


//public:
func NewTDecTop() *TDecTop{
  return &TDecTop{ m_pcPic : nil,
  m_iGopSize : 0,
  m_bGopSizeSet : false,
  m_iMaxRefPicNum : 0,
//#if ENC_DEC_TRACE
//  g_hTrace = fopen( "TraceDec.txt", "wb" );
//  g_bJustDoIt = g_bEncDecTraceDisable;
//  g_nSymbolCounter = 0;
//#endif
  m_bRefreshPending : false,
  m_pocCRA : 0,
  m_prevRAPisBLA : false,
  m_pocRandomAccess : TLibCommon.MAX_INT,          
  m_prevPOC                 : TLibCommon.MAX_INT,
  m_bFirstSliceInPicture    : true,
  m_bFirstSliceInSequence   : true}
}
  
func (this *TDecTop)  Create  (){
  //this.m_cGopDecoder.Create();
  this.m_apcSlicePilot = TLibCommon.NewTComSlice();
  this.m_uiSliceIdx = 0;
}
func (this *TDecTop)  Destroy (){
}
  
func (this *TDecTop)  IsSkipPictureForBLA(iPOCLastDisplay *int) bool{
	return true;
}
func (this *TDecTop)  IsRandomAccessSkipPicture(iSkipFrame *int,  iPOCLastDisplay *int) bool {
	return true;
}

func (this *TDecTop)  SetDecodedPictureHashSEIEnabled(enabled int) { 
	//this.m_cGopDecoder.SetDecodedPictureHashSEIEnabled(enabled); 
}

func (this *TDecTop)  Init(){
}
func (this *TDecTop)  Decode(nalu *TLibCommon.InputNALUnit, iSkipFrame *int, iPOCLastDisplay *int) bool {
	return true;
}
  
func (this *TDecTop)  DeletePicBuffer(){
}

func (this *TDecTop)  ExecuteDeblockAndAlf(poc *int, rpcListPic *list.List, iSkipFrame *int,  iPOCLastDisplay *int){
}

//protected:
func (this *TDecTop)  xGetNewPicBuffer  (pcSlice *TLibCommon.TComSlice, rpcPic *TLibCommon.TComPic){
}
func (this *TDecTop)  xUpdateGopSize    (pcSlice *TLibCommon.TComSlice){
}
func (this *TDecTop)  xCreateLostPicture (iLostPOC int){
}

func (this *TDecTop)  xActivateParameterSets(){
}
func (this *TDecTop)  xDecodeSlice(nalu *TLibCommon.InputNALUnit, iSkipFrame *int, iPOCLastDisplay int) bool{
	return true;
}
func (this *TDecTop)  xDecodeVPS(){
}
func (this *TDecTop)  xDecodeSPS(){
}
func (this *TDecTop)  xDecodePPS(){
}
func (this *TDecTop)  xDecodeSEI( bs *TLibCommon.TComInputBitstream ){
}

