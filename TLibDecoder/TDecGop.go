package TLibDecoder

import (
    "container/list"
    "gohm/TLibCommon"
)

// ====================================================================================================================
// Class definition
// ====================================================================================================================

/// GOP decoder class
type TDecGop struct {
    //private:
    m_iGopSize int
    m_cListPic *list.List //  Dynamic buffer

    //  Access channel    
    m_pcEntropyDecoder			*TDecEntropy;
    m_pcSbacDecoder				*TDecSbac;
    m_pcBinCABAC				*TDecBinCabac;
    m_pcSbacDecoders			[]TDecSbac; // independant CABAC decoders
    m_pcBinCABACs				[]TDecBinCabac;
    m_pcCavlcDecoder			*TDecCavlc;
    m_pcSliceDecoder			*TDecSlice;
    m_pcLoopFilter				*TLibCommon.TComLoopFilter;
    m_pcSAO					    *TLibCommon.TComSampleAdaptiveOffset;
    m_dDecTime                     float64
    m_decodedPictureHashSEIEnabled int ///< Checksum(3)/CRC(2)/MD5(1)/disable(0) acting on decoded picture hash SEI message

    //! list that contains the CU address of each slice plus the end address 
    m_sliceStartCUAddress      *list.List
    m_LFCrossSliceBoundaryFlag *list.List
}

//public:
func NewTDecGop() *TDecGop {
    return &TDecGop{m_sliceStartCUAddress:list.New(), m_LFCrossSliceBoundaryFlag:list.New()}
}

func (this *TDecGop) Init(pcEntropyDecoder *TDecEntropy,
    pcSbacDecoder *TDecSbac,
    pcBinCabac *TDecBinCabac,
    pcCavlcDecoder *TDecCavlc,
    pcSliceDecoder *TDecSlice,
    pcLoopFilter *TLibCommon.TComLoopFilter,
    pcSAO *TLibCommon.TComSampleAdaptiveOffset) {
  this.m_pcEntropyDecoder      = pcEntropyDecoder;
  this.m_pcSbacDecoder         = pcSbacDecoder;
  this.m_pcBinCABAC            = pcBinCabac;
  this.m_pcCavlcDecoder        = pcCavlcDecoder;
  this.m_pcSliceDecoder        = pcSliceDecoder;
  this.m_pcLoopFilter          = pcLoopFilter;
  this.m_pcSAO  			   = pcSAO;    
}
func (this *TDecGop) Create() {
}
func (this *TDecGop) Destroy() {
}
func (this *TDecGop) DecompressSlice(pcBitstream *TLibCommon.TComInputBitstream, rpcPic *TLibCommon.TComPic) {
  pcSlice := rpcPic.GetSlice(rpcPic.GetCurrSliceIdx());
  // Table of extracted substreams.
  // These must be deallocated AND their internal fifos, too.
  //TComInputBitstream **ppcSubstreams = NULL;

  //-- For time output for each slice
  //long iBeforeTime = clock();
  
  uiStartCUAddr   := pcSlice.GetDependentSliceCurStartCUAddr();

  uiSliceStartCuAddr := pcSlice.GetSliceCurStartCUAddr();
  if uiSliceStartCuAddr == uiStartCUAddr {
    this.m_sliceStartCUAddress.PushBack(uiSliceStartCuAddr);
  }

  this.m_pcSbacDecoder.Init( this.m_pcBinCABAC );//(TDecBinIf*)
  this.m_pcEntropyDecoder.SetEntropyDecoder (this.m_pcSbacDecoder);

  var uiNumSubstreams uint;
  
  if pcSlice.GetPPS().GetEntropyCodingSyncEnabledFlag() {
  	uiNumSubstreams = uint(pcSlice.GetNumEntryPointOffsets()+1);
  }else{
  	uiNumSubstreams = uint(pcSlice.GetPPS().GetNumSubstreams());
  }
  
  // init each couple {EntropyDecoder, Substream}
  puiSubstreamSizes := pcSlice.GetSubstreamSizes();
  ppcSubstreams     := make([]*TLibCommon.TComInputBitstream, uiNumSubstreams);
  this.m_pcSbacDecoders = make([]TDecSbac, uiNumSubstreams);
  this.m_pcBinCABACs    = make([]TDecBinCabac,uiNumSubstreams);
  for ui := uint(0) ; ui < uiNumSubstreams ; ui++ {
    this.m_pcSbacDecoders[ui].Init(&this.m_pcBinCABACs[ui]);
    if ui+1 < uiNumSubstreams {
    	ppcSubstreams[ui] = pcBitstream.ExtractSubstream( puiSubstreamSizes[ui]);
    }else{
    	ppcSubstreams[ui] = pcBitstream.ExtractSubstream( pcBitstream.GetNumBitsLeft());
    }
  }

  for ui := uint(0) ; ui+1 < uiNumSubstreams; ui++ {
    this.m_pcEntropyDecoder.SetEntropyDecoder ( &this.m_pcSbacDecoders[uiNumSubstreams - 1 - ui] );
    this.m_pcEntropyDecoder.SetBitstream      ( ppcSubstreams   [uiNumSubstreams - 1 - ui] );
    this.m_pcEntropyDecoder.ResetEntropy      ( pcSlice);
  }

  this.m_pcEntropyDecoder.SetEntropyDecoder ( this.m_pcSbacDecoder  );
  this.m_pcEntropyDecoder.SetBitstream      ( ppcSubstreams[0] );
  this.m_pcEntropyDecoder.ResetEntropy      (pcSlice);

  if uiSliceStartCuAddr == uiStartCUAddr {
    this.m_LFCrossSliceBoundaryFlag.PushBack( pcSlice.GetLFCrossSliceBoundaryFlag());
  }
  this.m_pcSbacDecoders[0].Load(this.m_pcSbacDecoder);
  this.m_pcSliceDecoder.DecompressSlice( pcBitstream, ppcSubstreams, rpcPic, this.m_pcSbacDecoder, this.m_pcSbacDecoders);
  this.m_pcEntropyDecoder.SetBitstream(  ppcSubstreams[uiNumSubstreams-1] );
  // deallocate all created substreams, including internal buffers.
  /*for ui := uint(0); ui < uiNumSubstreams; ui++ {
    ppcSubstreams[ui]->deleteFifo();
    delete ppcSubstreams[ui];
  }
  delete[] ppcSubstreams;
  delete[] m_pcSbacDecoders; 
  delete[] m_pcBinCABACs; 
  */
  this.m_pcSbacDecoders = nil;
  this.m_pcBinCABACs = nil;
  //m_dDecTime += (Double)(clock()-iBeforeTime) / CLOCKS_PER_SEC;
}
func (this *TDecGop) FilterPicture(rpcPic *TLibCommon.TComPic) {
}
func (this *TDecGop) SetGopSize(i int) {
    this.m_iGopSize = i
}

func (this *TDecGop) SetDecodedPictureHashSEIEnabled(enabled int) {
    this.m_decodedPictureHashSEIEnabled = enabled
}
