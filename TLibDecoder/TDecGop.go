package TLibDecoder

import (
	"fmt"
	"time"
	"io"
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
    m_dDecTime                   time.Duration;//float64
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
	//do nothing
}
func (this *TDecGop) Destroy() {
	//do nothing
}
func (this *TDecGop) DecompressSlice(pcBitstream *TLibCommon.TComInputBitstream, rpcPic *TLibCommon.TComPic, pTraceFile io.Writer) {
  pcSlice := rpcPic.GetSlice(rpcPic.GetCurrSliceIdx());
  // Table of extracted substreams.
  // These must be deallocated AND their internal fifos, too.
  //TComInputBitstream **ppcSubstreams = NULL;

  //-- For time output for each slice
  iBeforeTime := time.Now()
  
  uiStartCUAddr   := pcSlice.GetDependentSliceCurStartCUAddr();

  uiSliceStartCuAddr := pcSlice.GetSliceCurStartCUAddr();
  if uiSliceStartCuAddr == uiStartCUAddr {
    this.m_sliceStartCUAddress.PushBack(uiSliceStartCuAddr);
  }

  this.m_pcSbacDecoder.Init( this.m_pcBinCABAC );//(TDecBinIf*)
  this.m_pcEntropyDecoder.SetEntropyDecoder (this.m_pcSbacDecoder);
  this.m_pcEntropyDecoder.SetTraceFile(pTraceFile);

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
    this.m_pcEntropyDecoder.SetTraceFile	  (pTraceFile);
    this.m_pcEntropyDecoder.SetBitstream      ( ppcSubstreams   [uiNumSubstreams - 1 - ui] );
    this.m_pcEntropyDecoder.ResetEntropy      ( pcSlice);
  }

  this.m_pcEntropyDecoder.SetEntropyDecoder ( this.m_pcSbacDecoder  );
  this.m_pcEntropyDecoder.SetTraceFile		(pTraceFile);
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
  
  lAfterTime := time.Now()
  this.m_dDecTime += lAfterTime.Sub(iBeforeTime);
}

func (this *TDecGop) FilterPicture(rpcPic *TLibCommon.TComPic) {
  pcSlice := rpcPic.GetSlice(rpcPic.GetCurrSliceIdx());

  //-- For time output for each slice
  iBeforeTime := time.Now();

  // deblocking filter
  bLFCrossTileBoundary := pcSlice.GetPPS().GetLoopFilterAcrossTilesEnabledFlag();
  this.m_pcLoopFilter.SetCfg(bLFCrossTileBoundary);
  //this.m_pcLoopFilter.LoopFilterPic( rpcPic );

  pcSlice = rpcPic.GetSlice(0);
  if pcSlice.GetSPS().GetUseSAO() {
    this.m_sliceStartCUAddress.PushBack(rpcPic.GetNumCUsInFrame()* rpcPic.GetNumPartInCU());
    rpcPic.CreateNonDBFilterInfo(this.m_sliceStartCUAddress, 0, this.m_LFCrossSliceBoundaryFlag, rpcPic.GetPicSym().GetNumTiles(), bLFCrossTileBoundary);
  }

  if pcSlice.GetSPS().GetUseSAO() {
    if pcSlice.GetSaoEnabledFlag()||pcSlice.GetSaoEnabledFlagChroma() {
      saoParam := rpcPic.GetPicSym().GetSaoParam();
      saoParam.SaoFlag[0] = pcSlice.GetSaoEnabledFlag();
      saoParam.SaoFlag[1] = pcSlice.GetSaoEnabledFlagChroma();
      this.m_pcSAO.SetSaoLcuBasedOptimization(true);
      this.m_pcSAO.CreatePicSaoInfo(rpcPic, this.m_sliceStartCUAddress.Len() - 1);
      this.m_pcSAO.SAOProcess(saoParam);
      this.m_pcSAO.PCMLFDisableProcess(rpcPic);
      this.m_pcSAO.DestroyPicSaoInfo();
    }
  }

  if pcSlice.GetSPS().GetUseSAO() {
    rpcPic.DestroyNonDBFilterInfo();
  }

  rpcPic.CompressMotion(); 
  
  var c string;
  
  if pcSlice.IsIntra() {
  	c = "I";
  }else if pcSlice.IsInterP(){
  	if pcSlice.IsReferenced() {
  		c = "P";
  	}else{
  		c = "p";
  	}
  }else{
  	if pcSlice.IsReferenced() {
    	c = "B";
    }else{
    	c = "b"
    }
  }
  
  //-- For time output for each slice
  fmt.Printf("\nPOC %4d TId: %1d ( %s-SLICE, QP%3d ) ", pcSlice.GetPOC(), pcSlice.GetTLayer(), c, pcSlice.GetSliceQp() );

  this.m_dDecTime += time.Now().Sub(iBeforeTime);
  fmt.Printf ("[DT %10v] ", this.m_dDecTime );
  this.m_dDecTime  = 0;

  for iRefList := 0; iRefList < 2; iRefList++ {
    fmt.Printf ("[L%d ", iRefList);
    for iRefIndex := 0; iRefIndex < pcSlice.GetNumRefIdx(TLibCommon.RefPicList(iRefList)); iRefIndex++ {
      fmt.Printf ("%d ", pcSlice.GetRefPOC(TLibCommon.RefPicList(iRefList), iRefIndex));
    }
    fmt.Printf ("] ");
  }
  if this.m_decodedPictureHashSEIEnabled>0 {
    this.CalcAndPrintHashStatus(rpcPic.GetPicYuvRec(), rpcPic.GetSEIs());
  }

  rpcPic.SetOutputMark(true);
  rpcPic.SetReconMark(true);
  
  this.m_sliceStartCUAddress.Init();
  this.m_LFCrossSliceBoundaryFlag.Init();
}
func (this *TDecGop) SetGopSize(i int) {
    this.m_iGopSize = i
}

func (this *TDecGop) SetDecodedPictureHashSEIEnabled(enabled int) {
    this.m_decodedPictureHashSEIEnabled = enabled
}

func (this *TDecGop) CalcAndPrintHashStatus(pic *TLibCommon.TComPicYuv, seis *TLibCommon.SEImessages){
  /*
  // calculate MD5sum for entire reconstructed picture 
  UChar recon_digest[3][16];
  Int numChar=0;
  const Char* hashType = "\0";

  if (seis && seis->picture_digest)
  {
    switch (seis->picture_digest->method)
    {
    case SEIDecodedPictureHash::MD5:
      {
        hashType = "MD5";
        calcMD5(pic, recon_digest);
        numChar = 16;
        break;
      }
    case SEIDecodedPictureHash::CRC:
      {
        hashType = "CRC";
        calcCRC(pic, recon_digest);
        numChar = 2;
        break;
      }
    case SEIDecodedPictureHash::CHECKSUM:
      {
        hashType = "Checksum";
        calcChecksum(pic, recon_digest);
        numChar = 4;
        break;
      }
    default:
      {
        assert (!"unknown hash type");
      }
    }
  }

  // compare digest against received version 
  const Char* ok = "(unk)";
  Bool mismatch = false;

  if (seis && seis->picture_digest)
  {
    ok = "(OK)";
    for(Int yuvIdx = 0; yuvIdx < 3; yuvIdx++)
    {
      for (UInt i = 0; i < numChar; i++)
      {
        if (recon_digest[yuvIdx][i] != seis->picture_digest->digest[yuvIdx][i])
        {
          ok = "(***ERROR***)";
          mismatch = true;
        }
      }
    }
  }

  //printf("[%s:%s,%s] ", hashType, digestToString(recon_digest, numChar), ok);

  if (mismatch)
  {
    g_md5_mismatch = true;
    printf("[rx%s:%s] ", hashType, digestToString(seis->picture_digest->digest, numChar));
  }
  */
}
