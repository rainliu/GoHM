package TLibDecoder

import (

)


/// slice decoder class
type TDecSlice struct{
//private:
  // access channel
  //TDecEntropy*    m_pcEntropyDecoder;
  //TDecCu*         m_pcCuDecoder;
  m_uiCurrSliceIdx uint;

  //TDecSbac*       m_pcBufferSbacDecoders;   ///< line to store temporary contexts, one per column of tiles.
  //TDecBinCABAC*   m_pcBufferBinCABACs;
  //TDecSbac*       m_pcBufferLowLatSbacDecoders;   ///< dependent tiles: line to store temporary contexts, one per column of tiles.
  //TDecBinCABAC*   m_pcBufferLowLatBinCABACs;
//#if DEPENDENT_SLICES
  //std::vector<TDecSbac*> CTXMem;
//#endif
}

//public:
func NewTDecSlice() *TDecSlice{
	return &TDecSlice{}
}
  
//func (this *TDecSlice) Init( TDecEntropy* pcEntropyDecoder, TDecCu* pcMbDecoder ){
//}

//func (this *TDecSlice) Create            ( TComSlice* pcSlice, Int iWidth, Int iHeight, UInt uiMaxWidth, UInt uiMaxHeight, UInt uiMaxDepth );
//func (this *TDecSlice) Destroy           ();
  
//func (this *TDecSlice) DecompressSlice   ( TComInputBitstream* pcBitstream, TComInputBitstream** ppcSubstreams,   TComPic*& rpcPic, TDecSbac* pcSbacDecoder, TDecSbac* pcSbacDecoders );

//#if DEPENDENT_SLICES
//  Void      initCtxMem(  UInt i );
//  Void      setCtxMem( TDecSbac* sb, Int b )   { CTXMem[b] = sb; }
//#endif
//};

/*
class ParameterSetManagerDecoder:public ParameterSetManager
{
public:
  ParameterSetManagerDecoder();
  virtual ~ParameterSetManagerDecoder();
  Void     storePrefetchedVPS(TComVPS *vps)  { m_vpsBuffer.storePS( vps->getVPSId(), vps); };
  TComVPS* getPrefetchedVPS  (Int vpsId);
  Void     storePrefetchedSPS(TComSPS *sps)  { m_spsBuffer.storePS( sps->getSPSId(), sps); };
  TComSPS* getPrefetchedSPS  (Int spsId);
  Void     storePrefetchedPPS(TComPPS *pps)  { m_ppsBuffer.storePS( pps->getPPSId(), pps); };
  TComPPS* getPrefetchedPPS  (Int ppsId);
  Void     applyPrefetchedPS();

private:
  ParameterSetMap<TComVPS> m_vpsBuffer;
  ParameterSetMap<TComSPS> m_spsBuffer; 
  ParameterSetMap<TComPPS> m_ppsBuffer;
};
*/