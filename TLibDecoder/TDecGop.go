package TLibDecoder

import (
	"container/list"
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
/*
public:
  TDecGop();
  virtual ~TDecGop();
  
  Void  init    ( TDecEntropy*            pcEntropyDecoder, 
                 TDecSbac*               pcSbacDecoder, 
                 TDecBinCABAC*           pcBinCABAC,
                 TDecCavlc*              pcCavlcDecoder, 
                 TDecSlice*              pcSliceDecoder, 
                 TComLoopFilter*         pcLoopFilter,
                 TComSampleAdaptiveOffset* pcSAO
                 );
  Void  create  ();
  Void  destroy ();
  Void  decompressSlice(TComInputBitstream* pcBitstream, TComPic*& rpcPic );
  Void  filterPicture  (TComPic*& rpcPic );
  Void  setGopSize( Int i) { m_iGopSize = i; }

  void setDecodedPictureHashSEIEnabled(Int enabled) { m_decodedPictureHashSEIEnabled = enabled; }

};*/