package TLibDecoder

import (
	"container/list"
	"gohm/TLibCommon"
)


// ====================================================================================================================
// Class definition
// ====================================================================================================================

/// GOP decoder class
type TDecGop struct{
//private:
  m_iGopSize	int;
  m_cListPic	*list.List;         //  Dynamic buffer
  
  //  Access channel
 /* 
  TDecEntropy*          m_pcEntropyDecoder;
  TDecSbac*             m_pcSbacDecoder;
  TDecBinCABAC*         m_pcBinCABAC;
  TDecSbac*             m_pcSbacDecoders; // independant CABAC decoders
  TDecBinCABAC*         m_pcBinCABACs;
  TDecCavlc*            m_pcCavlcDecoder;
  TDecSlice*            m_pcSliceDecoder;
  TComLoopFilter*       m_pcLoopFilter;
  
  TComSampleAdaptiveOffset*     m_pcSAO;*/
  m_dDecTime	float64;
  m_decodedPictureHashSEIEnabled	int;  ///< Checksum(3)/CRC(2)/MD5(1)/disable(0) acting on decoded picture hash SEI message

  //! list that contains the CU address of each slice plus the end address 
  m_sliceStartCUAddress *list.List;
  m_LFCrossSliceBoundaryFlag	*list.List;
}

//public:
func NewTDecGop() *TDecGop{
	return &TDecGop{}
}
  
func (this *TDecGop)  Init    ( pcEntropyDecoder	*TDecEntropy, 
                 pcSbacDecoder 			*TDecSbac, 
                 pcBinCabac				*TDecBinCabac,
                 pcCavlcDecoder			*TDecCavlc, 
                 pcSliceDecoder			*TDecSlice, 
                 pcLoopFilter			*TLibCommon.TComLoopFilter,
                 pcSAO					*TLibCommon.TComSampleAdaptiveOffset){
}
func (this *TDecGop)   Create  (){
}
func (this *TDecGop)   Destroy (){
}
func (this *TDecGop)   DecompressSlice(pcBitstream *TLibCommon.TComInputBitstream, rpcPic *TLibCommon.TComPic){
}
func (this *TDecGop)   FilterPicture  (rpcPic *TLibCommon.TComPic){
}
func (this *TDecGop)   SetGopSize( i int ) { 
	this.m_iGopSize = i; 
}

func (this *TDecGop)   SetDecodedPictureHashSEIEnabled(enabled int) { 
	this.m_decodedPictureHashSEIEnabled = enabled; 
}
