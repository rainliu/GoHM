package TLibDecoder

import (
	"os"
    "container/list"
    "fmt"
    "gohm/TLibCommon"

)

/// decoder class
type TDecTop struct {
    //private:
    m_iGopSize      int
    m_bGopSizeSet   bool
    m_iMaxRefPicNum int

    m_bRefreshPending bool ///< refresh pending flag
    m_pocCRA          int  ///< POC number of the latest CRA picture
    m_prevRAPisBLA    bool ///< true if the previous RAP (CRA/CRANT/BLA/BLANT/IDR) picture is a BLA/BLANT picture
    m_pocRandomAccess int  ///< POC number of the random access point (the first IDR or CRA picture)

    m_pcListPic                  *list.List                 //  Dynamic buffer
    m_parameterSetManagerDecoder *TLibCommon.ParameterSetManager // storage for parameter sets
    m_apcSlicePilot              *TLibCommon.TComSlice

    m_SEIs *TLibCommon.SEImessages ///< "all" SEI messages.  If not NULL, we own the object.

    // functional classes

    m_cPrediction     *TLibCommon.TComPrediction
    m_cTrQuant        *TLibCommon.TComTrQuant
    m_cGopDecoder     *TDecGop
    m_cSliceDecoder   *TDecSlice
    m_cCuDecoder      *TDecCu
    m_cEntropyDecoder *TDecEntropy
    m_cCavlcDecoder   *TDecCavlc
    m_cSbacDecoder    *TDecSbac
    m_cBinCabac       *TDecBinCabac
    m_cSeiReader      *TDecSeiReader
    m_cLoopFilter     *TLibCommon.TComLoopFilter
    m_cSAO            *TLibCommon.TComSampleAdaptiveOffset

    m_pcPic                 *TLibCommon.TComPic
    m_uiSliceIdx            uint
    m_prevPOC               int
    m_bFirstSliceInPicture  bool
    m_bFirstSliceInSequence bool

    //static
    warningMessage bool
    m_pTraceFile	*os.File
}

//public:
func NewTDecTop() *TDecTop {
    return &TDecTop{m_pcPic: nil,
        m_iGopSize:      0,
        m_bGopSizeSet:   false,
        m_iMaxRefPicNum: 0,
        m_bRefreshPending:       false,
        m_pocCRA:                0,
        m_prevRAPisBLA:          false,
        m_pocRandomAccess:       TLibCommon.MAX_INT,
        m_prevPOC:               TLibCommon.MAX_INT,
        m_bFirstSliceInPicture:  true,
        m_bFirstSliceInSequence: true,
        m_pcListPic:             list.New(),
        m_parameterSetManagerDecoder: TLibCommon.NewParameterSetManager(),
        warningMessage:          false,
        m_cPrediction     :TLibCommon.NewTComPrediction(),
    	m_cTrQuant        :TLibCommon.NewTComTrQuant(),
    	m_cGopDecoder     :NewTDecGop(),
    	m_cSliceDecoder   :NewTDecSlice(),
    	m_cCuDecoder      :NewTDecCu(),
    	m_cEntropyDecoder :NewTDecEntropy(),
    	m_cCavlcDecoder   :NewTDecCavlc(),
    	m_cSbacDecoder    :NewTDecSbac(),
    	m_cBinCabac       :NewTDecBinCabac(),
    	m_cSeiReader      :NewTDecSeiReader(),
    	m_cLoopFilter     :nil,//TLibCommon.NewTComLoopFilter(),
    	m_cSAO            :nil}//TLibCommon.NewTComSampleAdaptiveOffset()}
}

func (this *TDecTop) Create(pchTraceFile string) {
    this.m_cGopDecoder.Create()
    this.m_apcSlicePilot = TLibCommon.NewTComSlice()
    this.m_uiSliceIdx = 0

    if pchTraceFile!=""{
    	this.m_pTraceFile, _ = os.Create(pchTraceFile)
    }else{
    	this.m_pTraceFile = nil
    }
}
func (this *TDecTop) Destroy() {
	if this.m_pTraceFile!=nil{
		this.m_pTraceFile.Close();
	}

    this.m_cGopDecoder.Destroy()
    this.m_apcSlicePilot = nil
    this.m_cSliceDecoder.Destroy()
}

func (this *TDecTop) IsSkipPictureForBLA(iPOCLastDisplay *int) bool {
    if this.m_prevRAPisBLA &&
        this.m_apcSlicePilot.GetPOC() < this.m_pocCRA &&
        this.m_apcSlicePilot.GetNalUnitType() == TLibCommon.NAL_UNIT_CODED_SLICE_TFD {
        (*iPOCLastDisplay)++
        return true
    }
    return false
}
func (this *TDecTop) IsRandomAccessSkipPicture(iSkipFrame *int, iPOCLastDisplay *int) bool {
    if *iSkipFrame != 0 {
        *iSkipFrame-- // decrement the counter
        return true
    } else if this.m_pocRandomAccess == TLibCommon.MAX_INT { // start of random access point, m_pocRandomAccess has not been set yet.
        if this.m_apcSlicePilot.GetNalUnitType() == TLibCommon.NAL_UNIT_CODED_SLICE_CRA ||
            this.m_apcSlicePilot.GetNalUnitType() == TLibCommon.NAL_UNIT_CODED_SLICE_BLA ||
            this.m_apcSlicePilot.GetNalUnitType() == TLibCommon.NAL_UNIT_CODED_SLICE_BLA_N_LP ||
            this.m_apcSlicePilot.GetNalUnitType() == TLibCommon.NAL_UNIT_CODED_SLICE_BLANT {
            // set the POC random access since we need to skip the reordered pictures in the case of CRA/CRANT/BLA/BLANT.
            this.m_pocRandomAccess = this.m_apcSlicePilot.GetPOC()
        } else if this.m_apcSlicePilot.GetNalUnitType() == TLibCommon.NAL_UNIT_CODED_SLICE_IDR ||
            this.m_apcSlicePilot.GetNalUnitType() == TLibCommon.NAL_UNIT_CODED_SLICE_IDR_N_LP {
            this.m_pocRandomAccess = -TLibCommon.MAX_INT // no need to skip the reordered pictures in IDR, they are decodable.
        } else {
            //static Bool warningMessage = false;
            if !this.warningMessage {
                fmt.Printf("\nWarning: this is not a valid random access point and the data is discarded until the first CRA picture")
                this.warningMessage = true
            }
            return true
        }
    } else if this.m_apcSlicePilot.GetPOC() < this.m_pocRandomAccess &&
        this.m_apcSlicePilot.GetNalUnitType() == TLibCommon.NAL_UNIT_CODED_SLICE_TFD { // skip the reordered pictures, if necessary
        *iPOCLastDisplay++
        return true
    }
    // if we reach here, then the picture is not skipped.
    return false
}

func (this *TDecTop) SetDecodedPictureHashSEIEnabled(enabled int) {
    this.m_cGopDecoder.SetDecodedPictureHashSEIEnabled(enabled)
}

func (this *TDecTop) Init() {
    // initialize ROM
    TLibCommon.InitROM()
    this.m_cGopDecoder.Init(this.m_cEntropyDecoder, this.m_cSbacDecoder, this.m_cBinCabac, this.m_cCavlcDecoder, this.m_cSliceDecoder, this.m_cLoopFilter, this.m_cSAO)
    this.m_cSliceDecoder.Init(this.m_cEntropyDecoder, this.m_cCuDecoder)
    this.m_cEntropyDecoder.Init(this.m_cPrediction)
}

func (this *TDecTop) Decode(nalu *InputNALUnit, iSkipFrame *int, iPOCLastDisplay *int, bSliceTrace bool) bool {
    // Initialize entropy decoder
    this.m_cEntropyDecoder.SetEntropyDecoder (this.m_cCavlcDecoder);
    this.m_cEntropyDecoder.SetBitstream(nalu.GetBitstream());
    this.m_cEntropyDecoder.SetTraceFile(this.m_pTraceFile);
    this.m_cEntropyDecoder.SetSliceTrace(bSliceTrace);

    switch nalu.GetNalUnitType() {
    case TLibCommon.NAL_UNIT_VPS:
        this.xDecodeVPS()
        return false

    case TLibCommon.NAL_UNIT_SPS:
        this.xDecodeSPS()
        return false

    case TLibCommon.NAL_UNIT_PPS:
        this.xDecodePPS()
        return false

    case TLibCommon.NAL_UNIT_SEI:
        fallthrough
        //#if SUFFIX_SEI_NUT_DECODED_HASH_SEI
    case TLibCommon.NAL_UNIT_SEI_SUFFIX:
        this.xDecodeSEI(nalu.GetBitstream(), nalu.GetNalUnitType())
        //#else
        //      xDecodeSEI( nalu.m_Bitstream );
        //#endif
        return false

    case TLibCommon.NAL_UNIT_CODED_SLICE_TRAIL_R:
        fallthrough
    case TLibCommon.NAL_UNIT_CODED_SLICE_TRAIL_N:
        fallthrough
    case TLibCommon.NAL_UNIT_CODED_SLICE_TLA:
        fallthrough
    case TLibCommon.NAL_UNIT_CODED_SLICE_TSA_N:
        fallthrough
    case TLibCommon.NAL_UNIT_CODED_SLICE_STSA_R:
        fallthrough
    case TLibCommon.NAL_UNIT_CODED_SLICE_STSA_N:
        fallthrough
    case TLibCommon.NAL_UNIT_CODED_SLICE_BLA:
        fallthrough
    case TLibCommon.NAL_UNIT_CODED_SLICE_BLANT:
        fallthrough
    case TLibCommon.NAL_UNIT_CODED_SLICE_BLA_N_LP:
        fallthrough
    case TLibCommon.NAL_UNIT_CODED_SLICE_IDR:
        fallthrough
    case TLibCommon.NAL_UNIT_CODED_SLICE_IDR_N_LP:
        fallthrough
    case TLibCommon.NAL_UNIT_CODED_SLICE_CRA:
        fallthrough
    case TLibCommon.NAL_UNIT_CODED_SLICE_DLP:
        fallthrough
    case TLibCommon.NAL_UNIT_CODED_SLICE_TFD:
        return this.xDecodeSlice(nalu, iSkipFrame, *iPOCLastDisplay)
        break
    default:
        //assert (1);
    }

    return false
}

func (this *TDecTop) DeletePicBuffer() {
	for e := this.m_pcListPic.Front(); e != nil; e = e.Next() {
		pcPic := e.Value.(*TLibCommon.TComPic)
		pcPic.Destroy();
		this.m_pcListPic.Remove(e)
	}

    this.m_cSAO.Destroy()
    this.m_cLoopFilter.Destroy()

    // destroy ROM
    TLibCommon.DestroyROM()
}

func (this *TDecTop) ExecuteLoopFilters(poc *int, iSkipFrame *int, iPOCLastDisplay *int) *list.List {
    if this.m_pcPic == nil {
        /* nothing to deblock */
        return nil
    }

    //pcPic := this.m_pcPic

    // Execute Deblock + Cleanup
    this.m_cGopDecoder.FilterPicture(this.m_pcPic)

    TLibCommon.SortPicList(this.m_pcListPic)// sorting for application output
    *poc = this.m_pcPic.GetSlice(this.m_uiSliceIdx - 1).GetPOC()
    this.m_cCuDecoder.Destroy()
    this.m_bFirstSliceInPicture = true

    return this.m_pcListPic
}

//protected:
func (this *TDecTop) xGetNewPicBuffer(pcSlice *TLibCommon.TComSlice) (rpcPic *TLibCommon.TComPic) {
    var numReorderPics [TLibCommon.MAX_TLAYER]int
    picCroppingWindow := pcSlice.GetSPS().GetPicCroppingWindow()

    for temporalLayer := uint(0); temporalLayer < TLibCommon.MAX_TLAYER; temporalLayer++ {
        numReorderPics[temporalLayer] = pcSlice.GetSPS().GetNumReorderPics(temporalLayer)
    }

    this.xUpdateGopSize(pcSlice)

    this.m_iMaxRefPicNum = int(pcSlice.GetSPS().GetMaxDecPicBuffering(pcSlice.GetTLayer())) + pcSlice.GetSPS().GetNumReorderPics(pcSlice.GetTLayer()) + 1 // +1 to have space for the picture currently being decoded
    if this.m_pcListPic.Len() < this.m_iMaxRefPicNum {
        rpcPic = TLibCommon.NewTComPic()

        rpcPic.Create(int(pcSlice.GetSPS().GetPicWidthInLumaSamples()), int(pcSlice.GetSPS().GetPicHeightInLumaSamples()),
            TLibCommon.G_uiMaxCUWidth, TLibCommon.G_uiMaxCUHeight, TLibCommon.G_uiMaxCUDepth,
            picCroppingWindow, numReorderPics[:], true)
        rpcPic.GetPicSym().AllocSaoParam(this.m_cSAO)
        this.m_pcListPic.PushBack(rpcPic)

        return rpcPic
    }

    bBufferIsAvailable := false;
	for e := this.m_pcListPic.Front(); e != nil; e = e.Next() {
		rpcPic = e.Value.(*TLibCommon.TComPic)
		if rpcPic.GetReconMark() == false && rpcPic.GetOutputMark() == false {
          rpcPic.SetOutputMark(false);
          bBufferIsAvailable = true;
          break;
        }

        if rpcPic.GetSlice( 0 ).IsReferenced() == false  && rpcPic.GetOutputMark() == false{
          rpcPic.SetOutputMark(false);
          rpcPic.SetReconMark( false );
          rpcPic.GetPicYuvRec().SetBorderExtension( false );
          bBufferIsAvailable = true;
          break;
        }
	}

    if !bBufferIsAvailable {
        //There is no room for this picture, either because of faulty encoder or dropped NAL. Extend the buffer.
        this.m_iMaxRefPicNum++;
        rpcPic := TLibCommon.NewTComPic();
        this.m_pcListPic.PushBack( rpcPic );
    }

    rpcPic.Destroy()
    rpcPic.Create(int(pcSlice.GetSPS().GetPicWidthInLumaSamples()), int(pcSlice.GetSPS().GetPicHeightInLumaSamples()),
        TLibCommon.G_uiMaxCUWidth, TLibCommon.G_uiMaxCUHeight, TLibCommon.G_uiMaxCUDepth,
        picCroppingWindow, numReorderPics[:], true)
    rpcPic.GetPicSym().AllocSaoParam(this.m_cSAO)

    return rpcPic;
}
func (this *TDecTop) xUpdateGopSize(pcSlice *TLibCommon.TComSlice) {
    if !pcSlice.IsIntra() && !this.m_bGopSizeSet {
        this.m_iGopSize = pcSlice.GetPOC()
        this.m_bGopSizeSet = true
        this.m_cGopDecoder.SetGopSize(this.m_iGopSize)
    }
}
func (this *TDecTop) xCreateLostPicture(iLostPOC int) {
}

func (this *TDecTop) xActivateParameterSets() {
    this.m_parameterSetManagerDecoder.ApplyPS()

    pps := this.m_parameterSetManagerDecoder.GetPPS(this.m_apcSlicePilot.GetPPSId())
    //assert (pps != 0);

    sps := this.m_parameterSetManagerDecoder.GetSPS(pps.GetSPSId())
    //assert (sps != 0);

    this.m_apcSlicePilot.SetPPS(pps)
    this.m_apcSlicePilot.SetSPS(sps)
    pps.SetSPS(sps)

    if pps.GetEntropyCodingSyncEnabledFlag() {
        pps.SetNumSubstreams(int((sps.GetPicHeightInLumaSamples()+sps.GetMaxCUHeight()-1)/sps.GetMaxCUHeight()) * (pps.GetNumColumnsMinus1() + 1))
    } else {
        pps.SetNumSubstreams(1)
    }

    pps.SetMinCuDQPSize(sps.GetMaxCUWidth() >> (pps.GetMaxCuDQPDepth()))

    for i := uint(0); i < sps.GetMaxCUDepth()-TLibCommon.G_uiAddCUDepth; i++ {
        sps.SetAMPAcc(i, int(TLibCommon.B2U(sps.GetUseAMP())))
    }

    for i := sps.GetMaxCUDepth() - TLibCommon.G_uiAddCUDepth; i < sps.GetMaxCUDepth(); i++ {
        sps.SetAMPAcc(i, 0)
    }

    this.m_cSAO.Destroy()
    this.m_cSAO.Create(sps.GetPicWidthInLumaSamples(), sps.GetPicHeightInLumaSamples(), TLibCommon.G_uiMaxCUWidth, TLibCommon.G_uiMaxCUHeight)
    this.m_cLoopFilter.Create(TLibCommon.G_uiMaxCUDepth)
}

func (this *TDecTop) xDecodeSlice(nalu *InputNALUnit, iSkipFrame *int, iPOCLastDisplay int) bool {
  this.m_apcSlicePilot.InitSlice();

  if this.m_bFirstSliceInPicture {
    this.m_uiSliceIdx     = 0;
  }
  this.m_apcSlicePilot.SetSliceIdx(this.m_uiSliceIdx);
  if !this.m_bFirstSliceInPicture {
    this.m_apcSlicePilot.CopySliceInfo( this.m_pcPic.GetPicSym().GetSlice(this.m_uiSliceIdx-1) );
  }

  this.m_apcSlicePilot.SetNalUnitType(nalu.GetNalUnitType());
  if (this.m_apcSlicePilot.GetNalUnitType() == TLibCommon.NAL_UNIT_CODED_SLICE_TRAIL_N) ||
     (this.m_apcSlicePilot.GetNalUnitType() == TLibCommon.NAL_UNIT_CODED_SLICE_TSA_N)   ||
     (this.m_apcSlicePilot.GetNalUnitType() == TLibCommon.NAL_UNIT_CODED_SLICE_STSA_N)  {
    this.m_apcSlicePilot.SetTemporalLayerNonReferenceFlag(true);
  }
  this.m_apcSlicePilot.SetReferenced(true); // Putting this as true ensures that picture is referenced the first time it is in an RPS
  this.m_apcSlicePilot.SetTLayerInfo(nalu.GetTemporalId());

  this.m_cEntropyDecoder.DecodeSliceHeader (this.m_apcSlicePilot, this.m_parameterSetManagerDecoder);

  //fmt.Printf("POC=%d, bNextSlice=%v\n", this.m_apcSlicePilot.GetPOC(),this.m_apcSlicePilot.IsNextSlice());

  // exit when a new picture is found
  if this.m_apcSlicePilot.IsNextSlice() && this.m_apcSlicePilot.GetPOC()!=this.m_prevPOC && !this.m_bFirstSliceInSequence {
    if this.m_prevPOC >= this.m_pocRandomAccess {
      this.m_prevPOC = this.m_apcSlicePilot.GetPOC();
      return true;
    }
    this.m_prevPOC = this.m_apcSlicePilot.GetPOC();
  }

  // actual decoding starts here
  this.xActivateParameterSets();

  if this.m_apcSlicePilot.IsNextSlice() {
    this.m_prevPOC = this.m_apcSlicePilot.GetPOC();
  }
  this.m_bFirstSliceInSequence = false;
  if this.m_apcSlicePilot.IsNextSlice() {
    // Skip pictures due to random access
    if this.IsRandomAccessSkipPicture(iSkipFrame, &iPOCLastDisplay) {
      return false;
    }
    // Skip TFD pictures associated with BLA/BLANT pictures
    if this.IsSkipPictureForBLA(&iPOCLastDisplay) {
      return false;
    }
  }
  //detect lost reference picture and insert copy of earlier frame.
  /*lostPoc := this.m_apcSlicePilot.CheckThatAllRefPicsAreAvailable(this.m_pcListPic, this.m_apcSlicePilot.GetRPS(), true, this.m_pocRandomAccess);
  for lostPoc > 0 {
    this.xCreateLostPicture(lostPoc-1);
    lostPoc = this.m_apcSlicePilot.CheckThatAllRefPicsAreAvailable(this.m_pcListPic, this.m_apcSlicePilot.GetRPS(), true, this.m_pocRandomAccess);
  }*/
  if this.m_bFirstSliceInPicture {
    // Buffer initialize for prediction.
    this.m_cPrediction.InitTempBuff();
    this.m_apcSlicePilot.ApplyReferencePictureSet(this.m_pcListPic, this.m_apcSlicePilot.GetRPS());
    //  Get a new picture buffer
    this.m_pcPic = this.xGetNewPicBuffer (this.m_apcSlicePilot);

    // transfer any SEI messages that have been received to the picture
    this.m_pcPic.SetSEIs(this.m_SEIs);
    this.m_SEIs = nil;

    // Recursive structure
    this.m_cCuDecoder.Create ( TLibCommon.G_uiMaxCUDepth, TLibCommon.G_uiMaxCUWidth, TLibCommon.G_uiMaxCUHeight );
    this.m_cCuDecoder.Init   ( this.m_cEntropyDecoder, this.m_cTrQuant, this.m_cPrediction );
    this.m_cTrQuant.Init     ( TLibCommon.G_uiMaxCUWidth, TLibCommon.G_uiMaxCUHeight, this.m_apcSlicePilot.GetSPS().GetMaxTrSize(),
    						   0, nil, nil, nil, false, false, false, false, false);

    this.m_cSliceDecoder.Create( this.m_apcSlicePilot, int(this.m_apcSlicePilot.GetSPS().GetPicWidthInLumaSamples()),
    							 int(this.m_apcSlicePilot.GetSPS().GetPicHeightInLumaSamples()), TLibCommon.G_uiMaxCUWidth, TLibCommon.G_uiMaxCUHeight, TLibCommon.G_uiMaxCUDepth );
  }

  //  Set picture slice pointer
  pcSlice := this.m_apcSlicePilot;
  bNextSlice := pcSlice.IsNextSlice();

  var uiCummulativeTileWidth, uiCummulativeTileHeight uint;
  var i, j, p int;

  //set NumColumnsMins1 and NumRowsMinus1
  this.m_pcPic.GetPicSym().SetNumColumnsMinus1( pcSlice.GetPPS().GetNumColumnsMinus1() );
  this.m_pcPic.GetPicSym().SetNumRowsMinus1	  ( pcSlice.GetPPS().GetNumRowsMinus1() );

  //create the TComTileArray
  this.m_pcPic.GetPicSym().XCreateTComTileArray();

  if pcSlice.GetPPS().GetUniformSpacingFlag() {
    //set the width for each tile
    for j=0; j < this.m_pcPic.GetPicSym().GetNumRowsMinus1()+1; j++ {
      for p=0; p < this.m_pcPic.GetPicSym().GetNumColumnsMinus1()+1; p++ {
      	a := (p+1)*int(this.m_pcPic.GetPicSym().GetFrameWidthInCU())/(this.m_pcPic.GetPicSym().GetNumColumnsMinus1()+1)
      	b := (p*int(this.m_pcPic.GetPicSym().GetFrameWidthInCU()))/(this.m_pcPic.GetPicSym().GetNumColumnsMinus1()+1)
        this.m_pcPic.GetPicSym().GetTComTile(uint( j * (this.m_pcPic.GetPicSym().GetNumColumnsMinus1()+1) + p )).SetTileWidth( uint(a-b) );
      }
    }

    //set the height for each tile
    for j=0; j < this.m_pcPic.GetPicSym().GetNumColumnsMinus1()+1; j++ {
      for p=0; p < this.m_pcPic.GetPicSym().GetNumRowsMinus1()+1; p++ {
      	a := (p+1)*int(this.m_pcPic.GetPicSym().GetFrameHeightInCU())/(this.m_pcPic.GetPicSym().GetNumRowsMinus1()+1);
      	b := (p*int(this.m_pcPic.GetPicSym().GetFrameHeightInCU()))/(this.m_pcPic.GetPicSym().GetNumRowsMinus1()+1)
        this.m_pcPic.GetPicSym().GetTComTile(uint( p * (this.m_pcPic.GetPicSym().GetNumColumnsMinus1()+1) + j )).SetTileHeight( uint(a-b) );
      }
    }
  }else{
    //set the width for each tile
    for j=0; j < pcSlice.GetPPS().GetNumRowsMinus1()+1; j++ {
      uiCummulativeTileWidth = 0;
      for i=0; i < pcSlice.GetPPS().GetNumColumnsMinus1(); i++ {
        this.m_pcPic.GetPicSym().GetTComTile(uint(j * (pcSlice.GetPPS().GetNumColumnsMinus1()+1) + i)).SetTileWidth( pcSlice.GetPPS().GetColumnWidth(uint(i)) );
        uiCummulativeTileWidth += pcSlice.GetPPS().GetColumnWidth(uint(i));
      }
      this.m_pcPic.GetPicSym().GetTComTile(uint(j * (pcSlice.GetPPS().GetNumColumnsMinus1()+1) + i)).SetTileWidth( this.m_pcPic.GetPicSym().GetFrameWidthInCU()-uiCummulativeTileWidth );
    }

    //set the height for each tile
    for j=0; j < pcSlice.GetPPS().GetNumColumnsMinus1()+1; j++ {
      uiCummulativeTileHeight = 0;
      for i=0; i < pcSlice.GetPPS().GetNumRowsMinus1(); i++ {
        this.m_pcPic.GetPicSym().GetTComTile(uint(i * (pcSlice.GetPPS().GetNumColumnsMinus1()+1) + j)).SetTileHeight( pcSlice.GetPPS().GetRowHeight(uint(i)) );
        uiCummulativeTileHeight += pcSlice.GetPPS().GetRowHeight(uint(i));
      }
      this.m_pcPic.GetPicSym().GetTComTile(uint(i * (pcSlice.GetPPS().GetNumColumnsMinus1()+1) + j)).SetTileHeight( this.m_pcPic.GetPicSym().GetFrameHeightInCU()-uiCummulativeTileHeight );
    }
  }

  this.m_pcPic.GetPicSym().XInitTiles();

  //generate the Coding Order Map and Inverse Coding Order Map
  uiEncCUAddr := 0;
  for i=0; i<int(this.m_pcPic.GetPicSym().GetNumberOfCUsInFrame()); i++ {
    this.m_pcPic.GetPicSym().SetCUOrderMap(i, uiEncCUAddr);
    this.m_pcPic.GetPicSym().SetInverseCUOrderMap(uiEncCUAddr, i);

    uiEncCUAddr = int(this.m_pcPic.GetPicSym().XCalculateNxtCUAddr(uint(uiEncCUAddr)));
  }
  this.m_pcPic.GetPicSym().SetCUOrderMap(int(this.m_pcPic.GetPicSym().GetNumberOfCUsInFrame()), int(this.m_pcPic.GetPicSym().GetNumberOfCUsInFrame()));
  this.m_pcPic.GetPicSym().SetInverseCUOrderMap(int(this.m_pcPic.GetPicSym().GetNumberOfCUsInFrame()), int(this.m_pcPic.GetPicSym().GetNumberOfCUsInFrame()));

  //convert the start and end CU addresses of the slice and dependent slice into encoding order
  pcSlice.SetDependentSliceCurStartCUAddr( this.m_pcPic.GetPicSym().GetPicSCUEncOrder(pcSlice.GetDependentSliceCurStartCUAddr()) );
  pcSlice.SetDependentSliceCurEndCUAddr( this.m_pcPic.GetPicSym().GetPicSCUEncOrder(pcSlice.GetDependentSliceCurEndCUAddr()) );
  if pcSlice.IsNextSlice(){
    pcSlice.SetSliceCurStartCUAddr(this.m_pcPic.GetPicSym().GetPicSCUEncOrder(pcSlice.GetSliceCurStartCUAddr()));
    pcSlice.SetSliceCurEndCUAddr(this.m_pcPic.GetPicSym().GetPicSCUEncOrder(pcSlice.GetSliceCurEndCUAddr()));
  }

  if this.m_bFirstSliceInPicture {
    if this.m_pcPic.GetNumAllocatedSlice() != 1 {
      this.m_pcPic.ClearSliceBuffer();
    }
  }else{
    this.m_pcPic.AllocateNewSlice();
  }
  //assert(pcPic.GetNumAllocatedSlice() == (this.m_uiSliceIdx + 1));
  this.m_apcSlicePilot = this.m_pcPic.GetPicSym().GetSlice(this.m_uiSliceIdx);
  //fmt.Printf("%v\n", this.m_apcSlicePilot)
  this.m_pcPic.GetPicSym().SetSlice(pcSlice, this.m_uiSliceIdx);

  this.m_pcPic.SetTLayer(nalu.GetTemporalId());

  if bNextSlice {
    pcSlice.CheckCRA(pcSlice.GetRPS(), &this.m_pocCRA, &this.m_prevRAPisBLA, this.m_pcListPic);
    // Set reference list
    pcSlice.SetRefPicList( this.m_pcListPic );

    // For generalized B
    // note: maybe not existed case (always L0 is copied to L1 if L1 is empty)
    if pcSlice.IsInterB() && pcSlice.GetNumRefIdx(TLibCommon.REF_PIC_LIST_1) == 0 {
      iNumRefIdx := pcSlice.GetNumRefIdx(TLibCommon.REF_PIC_LIST_0);
      pcSlice.SetNumRefIdx        ( TLibCommon.REF_PIC_LIST_1, iNumRefIdx );

      for iRefIdx := 0; iRefIdx < iNumRefIdx; iRefIdx++ {
        pcSlice.SetRefPic(pcSlice.GetRefPic(TLibCommon.REF_PIC_LIST_0, iRefIdx), TLibCommon.REF_PIC_LIST_1, iRefIdx);
      }
    }
    if !pcSlice.IsIntra() {
      bLowDelay := true;
      iCurrPOC  := pcSlice.GetPOC();
      iRefIdx := 0;

      for iRefIdx = 0; iRefIdx < pcSlice.GetNumRefIdx(TLibCommon.REF_PIC_LIST_0) && bLowDelay; iRefIdx++ {
        if int(pcSlice.GetRefPic(TLibCommon.REF_PIC_LIST_0, iRefIdx).GetPOC()) > iCurrPOC {
          bLowDelay = false;
        }
      }
      if pcSlice.IsInterB() {
        for iRefIdx = 0; iRefIdx < pcSlice.GetNumRefIdx(TLibCommon.REF_PIC_LIST_1) && bLowDelay; iRefIdx++ {
          if int(pcSlice.GetRefPic(TLibCommon.REF_PIC_LIST_1, iRefIdx).GetPOC()) > iCurrPOC {
            bLowDelay = false;
          }
        }
      }

      pcSlice.SetCheckLDC(bLowDelay);
    }

    //---------------
    pcSlice.SetRefPOCList();
    pcSlice.SetNoBackPredFlag( false );
    if pcSlice.GetSliceType() == TLibCommon.B_SLICE {
      if pcSlice.GetNumRefIdx(TLibCommon.RefPicList( 0 ) ) == pcSlice.GetNumRefIdx(TLibCommon.RefPicList( 1 ) ) {
        pcSlice.SetNoBackPredFlag( true );
        for i=0; i < pcSlice.GetNumRefIdx(TLibCommon.RefPicList( 1 ) ); i++ {
          if pcSlice.GetRefPOC(TLibCommon.RefPicList(1), i) != pcSlice.GetRefPOC(TLibCommon.RefPicList(0), i) {
            pcSlice.SetNoBackPredFlag( false );
            break;
          }
        }
      }
    }
  }

  this.m_pcPic.SetCurrSliceIdx(this.m_uiSliceIdx);
  if pcSlice.GetSPS().GetScalingListFlag() {
    pcSlice.SetScalingList ( pcSlice.GetSPS().GetScalingList()  );
    if pcSlice.GetPPS().GetScalingListPresentFlag() {
      pcSlice.SetScalingList ( pcSlice.GetPPS().GetScalingList()  );
    }
    pcSlice.GetScalingList().SetUseTransformSkip(pcSlice.GetPPS().GetUseTransformSkip());
    if !pcSlice.GetPPS().GetScalingListPresentFlag() && !pcSlice.GetSPS().GetScalingListPresentFlag() {
      pcSlice.SetDefaultScalingList();
    }
    this.m_cTrQuant.SetScalingListDec(pcSlice.GetScalingList());
    this.m_cTrQuant.SetUseScalingList(true);
  }else{
    this.m_cTrQuant.SetFlatScalingList();
    this.m_cTrQuant.SetUseScalingList(false);
  }

  //  Decode a picture
  this.m_cGopDecoder.DecompressSlice(nalu.m_Bitstream, this.m_pcPic, this.m_pTraceFile);

  this.m_bFirstSliceInPicture = false;
  this.m_uiSliceIdx++;

  return false;
}
func (this *TDecTop) xDecodeVPS() {
  vps := TLibCommon.NewTComVPS();

  this.m_cEntropyDecoder.DecodeVPS( vps );
  this.m_parameterSetManagerDecoder.SetVPS(vps);
}
func (this *TDecTop) xDecodeSPS() {
  sps := TLibCommon.NewTComSPS();
  this.m_cEntropyDecoder.DecodeSPS( sps );
  this.m_parameterSetManagerDecoder.SetSPS(sps);
}
func (this *TDecTop) xDecodePPS() {
  pps := TLibCommon.NewTComPPS();
  this.m_cEntropyDecoder.DecodePPS( pps, this.m_parameterSetManagerDecoder );
  this.m_parameterSetManagerDecoder.SetPPS( pps );

//#if DEPENDENT_SLICES
//#if REMOVE_ENTROPY_SLICES
  if pps.GetDependentSliceEnabledFlag() {
//#else
//  if( pps.GetDependentSliceEnabledFlag() && (!pps.GetEntropySliceEnabledFlag()) )
//#endif
	var NumCtx int;
	if pps.GetEntropyCodingSyncEnabledFlag(){
    	NumCtx = 2;
    }else{
    	NumCtx = 1;
    }
    this.m_cSliceDecoder.InitCtxMem(uint(NumCtx));
    for st := 0; st < NumCtx; st++ {
      ctx := NewTDecSbac();
      ctx.Init( this.m_cBinCabac );
      this.m_cSliceDecoder.SetCtxMem( ctx, st );
    }
  }
//#endif
}

//#if SUFFIX_SEI_NUT_DECODED_HASH_SEI
func (this *TDecTop) xDecodeSEI(bs *TLibCommon.TComInputBitstream, nalUnitType TLibCommon.NalUnitType) {
  if this.m_SEIs == nil{
//#if SUFFIX_SEI_NUT_DECODED_HASH_SEI
    if (nalUnitType == TLibCommon.NAL_UNIT_SEI_SUFFIX) && (this.m_pcPic.GetSEIs()!=nil) {
      this.m_SEIs = this.m_pcPic.GetSEIs();          // If suffix SEI and SEI already present, use already existing SEI structure
    }else{
      this.m_SEIs = TLibCommon.NewSEImessages();
    }
  }else{
    //assert(nalUnitType != NAL_UNIT_SEI_SUFFIX);
  }
/*#else
  {
    m_SEIs = new SEImessages;
  }
#endif*/
  this.m_SEIs.SetSPS(this.m_parameterSetManagerDecoder.GetSPS(0));
//#if SUFFIX_SEI_NUT_DECODED_HASH_SEI
  this.m_cSeiReader.ParseSEImessage( bs, this.m_SEIs, nalUnitType );
  if nalUnitType == TLibCommon.NAL_UNIT_SEI_SUFFIX {
    if this.m_pcPic.GetSEIs()==nil{
      this.m_pcPic.SetSEIs(this.m_SEIs); // Only suffix SEI present and new object created; update picture SEI variable
    }
    this.m_SEIs = nil;  // SEI structure already updated using this pointer; not required now.
  }
//#else
//  m_seiReader.parseSEImessage( bs, *m_SEIs );
//#endif
}
