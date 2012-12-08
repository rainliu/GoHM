package TLibDecoder

import (
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
    m_parameterSetManagerDecoder ParameterSetManagerDecoder // storage for parameter sets 
    m_apcSlicePilot              *TLibCommon.TComSlice

    m_SEIs *TLibCommon.SEImessages ///< "all" SEI messages.  If not NULL, we own the object.

    // functional classes

    m_cPrediction     TLibCommon.TComPrediction
    m_cTrQuant        TLibCommon.TComTrQuant
    m_cGopDecoder     TDecGop
    m_cSliceDecoder   TDecSlice
    m_cCuDecoder      TDecCu
    m_cEntropyDecoder TDecEntropy
    m_cCavlcDecoder   TDecCavlc
    m_cSbacDecoder    TDecSbac
    m_cBinCabac       TDecBinCabac
    m_cSeiReader      TDecSeiReader
    m_cLoopFilter     TLibCommon.TComLoopFilter
    m_cSAO            TLibCommon.TComSampleAdaptiveOffset

    m_pcPic                 *TLibCommon.TComPic
    m_uiSliceIdx            uint
    m_prevPOC               int
    m_bFirstSliceInPicture  bool
    m_bFirstSliceInSequence bool

    //static
    warningMessage bool
}

//public:
func NewTDecTop() *TDecTop {
    pcListPic := list.New()

    return &TDecTop{m_pcPic: nil,
        m_iGopSize:      0,
        m_bGopSizeSet:   false,
        m_iMaxRefPicNum: 0,
        //#if ENC_DEC_TRACE
        //  g_hTrace = fopen( "TraceDec.txt", "wb" );
        //  g_bJustDoIt = g_bEncDecTraceDisable;
        //  g_nSymbolCounter = 0;
        //#endif
        m_bRefreshPending:       false,
        m_pocCRA:                0,
        m_prevRAPisBLA:          false,
        m_pocRandomAccess:       TLibCommon.MAX_INT,
        m_prevPOC:               TLibCommon.MAX_INT,
        m_bFirstSliceInPicture:  true,
        m_bFirstSliceInSequence: true,
        m_pcListPic:             pcListPic,
        warningMessage:          false}
}

func (this *TDecTop) Create() {
    this.m_cGopDecoder.Create()
    this.m_apcSlicePilot = TLibCommon.NewTComSlice()
    this.m_uiSliceIdx = 0
}
func (this *TDecTop) Destroy() {
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
    this.m_cGopDecoder.Init(&this.m_cEntropyDecoder, &this.m_cSbacDecoder, &this.m_cBinCabac, &this.m_cCavlcDecoder, &this.m_cSliceDecoder, &this.m_cLoopFilter, &this.m_cSAO)
    this.m_cSliceDecoder.Init(&this.m_cEntropyDecoder, &this.m_cCuDecoder)
    this.m_cEntropyDecoder.Init(&this.m_cPrediction)
}

func (this *TDecTop) Decode(nalu *InputNALUnit, iSkipFrame *int, iPOCLastDisplay *int) bool {
    // Initialize entropy decoder
    this.m_cEntropyDecoder.SetEntropyDecoder (&this.m_cCavlcDecoder);
    this.m_cEntropyDecoder.SetBitstream(nalu.GetBitstream())

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

func (this *TDecTop) ExecuteLoopFilters(poc *int, rpcListPic *list.List, iSkipFrame *int, iPOCLastDisplay *int) {
    if this.m_pcPic == nil {
        /* nothing to deblock */
        return
    }

    pcPic := this.m_pcPic

    // Execute Deblock + Cleanup
    this.m_cGopDecoder.FilterPicture(pcPic)

    TLibCommon.SortPicList(this.m_pcListPic)// sorting for application output
    *poc = pcPic.GetSlice(this.m_uiSliceIdx - 1).GetPOC()
    rpcListPic = this.m_pcListPic
    this.m_cCuDecoder.Destroy()
    this.m_bFirstSliceInPicture = true

    return
}

//protected:
func (this *TDecTop) xGetNewPicBuffer(pcSlice *TLibCommon.TComSlice, rpcPic *TLibCommon.TComPic) {
    var numReorderPics [TLibCommon.MAX_TLAYER]int
    picCroppingWindow := pcSlice.GetSPS().GetPicCroppingWindow()

    for temporalLayer := uint(0); temporalLayer < TLibCommon.MAX_TLAYER; temporalLayer++ {
        numReorderPics[temporalLayer] = pcSlice.GetSPS().GetNumReorderPics(temporalLayer)
    }

    this.xUpdateGopSize(pcSlice)

    this.m_iMaxRefPicNum = int(pcSlice.GetSPS().GetMaxDecPicBuffering(pcSlice.GetTLayer())) + pcSlice.GetSPS().GetNumReorderPics(pcSlice.GetTLayer()) + 1 // +1 to have space for the picture currently being decoded
    if this.m_pcListPic.Len() < this.m_iMaxRefPicNum {
        rpcPic := TLibCommon.NewTComPic()

        rpcPic.Create(int(pcSlice.GetSPS().GetPicWidthInLumaSamples()), int(pcSlice.GetSPS().GetPicHeightInLumaSamples()),
            TLibCommon.G_uiMaxCUWidth, TLibCommon.G_uiMaxCUHeight, TLibCommon.G_uiMaxCUDepth,
            picCroppingWindow, numReorderPics[:], true)
        rpcPic.GetPicSym().AllocSaoParam(&this.m_cSAO)
        this.m_pcListPic.PushBack(rpcPic)

        return
    }

    bBufferIsAvailable := false;
	for e := this.m_pcListPic.Front(); e != nil; e = e.Next() {
		// do something with e.Value
		rpcPic := e.Value.(*TLibCommon.TComPic)
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
    rpcPic.GetPicSym().AllocSaoParam(&this.m_cSAO)
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
    this.m_parameterSetManagerDecoder.ApplyPrefetchedPS()

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
    return true
}
func (this *TDecTop) xDecodeVPS() {
  vps := TLibCommon.NewTComVPS();
  
  this.m_cEntropyDecoder.DecodeVPS( vps );
  this.m_parameterSetManagerDecoder.SetPrefetchedVPS(vps); 
}
func (this *TDecTop) xDecodeSPS() {
  sps := TLibCommon.NewTComSPS();
  this.m_cEntropyDecoder.DecodeSPS( sps );
  this.m_parameterSetManagerDecoder.SetPrefetchedSPS(sps);
}
func (this *TDecTop) xDecodePPS() {
  pps := TLibCommon.NewTComPPS();
  this.m_cEntropyDecoder.DecodePPS( pps, &this.m_parameterSetManagerDecoder );
  this.m_parameterSetManagerDecoder.SetPrefetchedPPS( pps );

//#if DEPENDENT_SLICES
//#if REMOVE_ENTROPY_SLICES
  if pps.GetDependentSliceEnabledFlag() {
//#else
//  if( pps->getDependentSliceEnabledFlag() && (!pps->getEntropySliceEnabledFlag()) )
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
      ctx.Init( &this.m_cBinCabac );
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

