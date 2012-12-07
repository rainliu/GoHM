package TAppDecoder

import (
	"fmt"
	"os"
	//"errors"
    "container/list"
    "gohm/TLibCommon"
    "gohm/TLibDecoder"
)

type TAppDecTop struct {
    TAppDecCfg

    m_cTDecTop              TLibDecoder.TDecTop
    m_cTVideoIOYuvReconFile TLibCommon.TVideoIOYuv

    m_abDecFlag       [TLibCommon.MAX_GOP]bool
    m_iPOCLastDisplay int
}

func NewTAppDecTop() *TAppDecTop {
    pAppDecTop := &TAppDecTop{}
    //::memset (m_abDecFlag, 0, sizeof (m_abDecFlag));//memset 0 by Go
    pAppDecTop.m_iPOCLastDisplay = -TLibCommon.MAX_GOP

    return pAppDecTop
}

func (this *TAppDecTop) Create() {
    //do nothing
}

func (this *TAppDecTop) Destroy(){
  /*if this.m_pchBitstreamFile !=
  {
    free (m_pchBitstreamFile);
    m_pchBitstreamFile = NULL;
  }
  if (m_pchReconFile)
  {
    free (m_pchReconFile);
    m_pchReconFile = NULL;
  }*/
}

func (this *TAppDecTop) Decode() (err error){
  //var poc int;
  var pcListPic *list.List;// = NULL;

  bitstreamFile, err := os.Open(this.m_pchBitstreamFile);
  if err != nil{
  	fmt.Printf("\nfailed to open bitstream file `%s' for reading\n", this.m_pchBitstreamFile);
  	return err
  }
  defer bitstreamFile.Close();

  bytestream := TLibDecoder.NewInputByteStream (bitstreamFile);

  // create & initialize internal classes
  this.xCreateDecLib();
  this.xInitDecLib  ();
  this.m_iPOCLastDisplay += this.m_iSkipFrame;      // set the last displayed POC correctly for skip forward.

  // main decoder loop
  recon_opened := false; // reconstruction file not yet opened. (must be performed after SPS is seen)

  for {// (!!bitstreamFile)
    /* location serves to work around a design fault in the decoder, whereby
     * the process of reading a new slice that is the first slice of a new frame
     * requires the TDecTop::decode() method to be called again with the same
     * nal unit. */
    //streampos location = bitstreamFile.tellg();
    var stats TLibDecoder.AnnexBStats;// stats = AnnexBStats();
    bPreviousPictureDecoded := false;

    nalUnit	:= list.New(); //vector<uint8_t> 
    var	nalu TLibCommon.InputNALUnit;
    bytestream.ByteStreamNALUnit(nalUnit, &stats);

    // call actual decoding function
    bNewPicture := false;
    if nalUnit.Len()==0 {
      /* this can happen if the following occur:
       *  - empty input file
       *  - two back-to-back start_code_prefixes
       *  - start_code_prefix immediately followed by EOF
       */
      fmt.Printf("Warning: Attempt to decode an empty NAL unit\n");
    }else{
      nalu.Read(nalUnit);
      if (this.m_iMaxTemporalLayer >= 0 && int(nalu.GetTemporalId()) > this.m_iMaxTemporalLayer) || 
      	 !this.IsNaluWithinTargetDecLayerIdSet(&nalu) {
        if bPreviousPictureDecoded {
          bNewPicture = true;
          bPreviousPictureDecoded = false;
        }else{
          bNewPicture = false;
        }
      }else{
        bNewPicture = this.m_cTDecTop.Decode(&nalu, &this.m_iSkipFrame, &this.m_iPOCLastDisplay);
        if bNewPicture {
          //bitstreamFile.clear();
          /* location points to the current nalunit payload[1] due to the
           * need for the annexB parser to read three extra bytes.
           * [1] except for the first NAL unit in the file
           *     (but bNewPicture doesn't happen then) */
          //bitstreamFile.seekg(location-streamoff(3));
          //bytestream.reset();
/*#if ENC_DEC_TRACE
          g_bSliceTrace = false;
#endif*/
        }
/*#if ENC_DEC_TRACE
        else
        {
          g_bSliceTrace = true;
        }
#endif*/
        bPreviousPictureDecoded = true; 
      }
    }
    //if bNewPicture || !bitstreamFile {
    //  this.m_cTDecTop.executeLoopFilters(poc, pcListPic, m_iSkipFrame, m_iPOCLastDisplay);
    //}

    if pcListPic != nil {
      if this.m_pchReconFile!="" && !recon_opened {
        if  this.m_outputBitDepthY==0 { 
        	this.m_outputBitDepthY = TLibCommon.G_bitDepthY; 
        }
        if this.m_outputBitDepthC==0 { 
        	this.m_outputBitDepthC = TLibCommon.G_bitDepthC; 
        }

        this.m_cTVideoIOYuvReconFile.Open( this.m_pchReconFile, true, this.m_outputBitDepthY, this.m_outputBitDepthC, TLibCommon.G_bitDepthY, TLibCommon.G_bitDepthC ); // write mode
        recon_opened = true;
      }
      if  bNewPicture && 
         ( nalu.GetNalUnitType() == TLibCommon.NAL_UNIT_CODED_SLICE_IDR		||
           nalu.GetNalUnitType() == TLibCommon.NAL_UNIT_CODED_SLICE_IDR_N_LP	||
           nalu.GetNalUnitType() == TLibCommon.NAL_UNIT_CODED_SLICE_BLA_N_LP	||
           nalu.GetNalUnitType() == TLibCommon.NAL_UNIT_CODED_SLICE_BLANT		||
           nalu.GetNalUnitType() == TLibCommon.NAL_UNIT_CODED_SLICE_BLA ) {
        this.xFlushOutput( pcListPic );
      }
      // write reconstruction to file
      if bNewPicture {
        this.xWriteOutput( pcListPic, nalu.GetTemporalId() );
      }
    }
  }
  
  this.xFlushOutput( pcListPic );
  // delete buffers
  this.m_cTDecTop.DeletePicBuffer();
  
  // destroy internal classes
  this.xDestroyDecLib();
  
  return nil
}

func (this *TAppDecTop) xCreateDecLib() {
    //create decoder class
    this.m_cTDecTop.Create()
}

func (this *TAppDecTop) xDestroyDecLib() {
    if this.m_pchReconFile != "" {
        this.m_cTVideoIOYuvReconFile.Close()
    }

    //destroy decoder class
    this.m_cTDecTop.Destroy()
}

func (this *TAppDecTop) xInitDecLib() {
    //initialize decoder class
    this.m_cTDecTop.Init()
    this.m_cTDecTop.SetDecodedPictureHashSEIEnabled(this.m_decodedPictureHashSEIEnabled)

}

func (this *TAppDecTop) xWriteOutput(pcListPic *list.List, tId uint) {
    not_displayed := 0;

	for e := pcListPic.Front(); e != nil; e = e.Next() {
	    pcPic := e.Value.(*TLibCommon.TComPic)
	    if pcPic.GetOutputMark() && int(pcPic.GetPOC()) > this.m_iPOCLastDisplay {
	        not_displayed++;
	    }
	}

	for e := pcListPic.Front(); e != nil; e = e.Next() {
	    pcPic := e.Value.(*TLibCommon.TComPic)
	    if pcPic.GetOutputMark() && (not_displayed >  pcPic.GetNumReorderPics(tId) && int(pcPic.GetPOC()) > this.m_iPOCLastDisplay){
            // write to file
            not_displayed--;
            if this.m_pchReconFile != "" {
              crop := pcPic.GetCroppingWindow();
              this.m_cTVideoIOYuvReconFile.Write( pcPic.GetPicYuvRec(), crop.GetPicCropLeftOffset(), crop.GetPicCropRightOffset(), crop.GetPicCropTopOffset(), crop.GetPicCropBottomOffset() );
            }

            // update POC of display order
            this.m_iPOCLastDisplay = int(pcPic.GetPOC());

            // erase non-referenced picture in the reference picture list after display
            if !pcPic.GetSlice(0).IsReferenced() && pcPic.GetReconMark() == true {
//#   if !D   YN_REF_FREE
              pcPic.SetReconMark(false);

              // mark it should be extended later
              pcPic.GetPicYuvRec().SetBorderExtension( false );

//#   else
//              pcPic->destroy();
//              pcListPic->erase( iterPic );
//              iterPic = pcListPic->begin(); // to the beginning, non-efficient way, have to be revised!
//              continue;
//#   endif
            }
            pcPic.SetOutputMark(false);
       }
	}

}

func (this *TAppDecTop) xFlushOutput( pcListPic *list.List ) {
  if pcListPic==nil {
    return;
  } 
  
  for e := pcListPic.Front(); e != nil; e = e.Next() {
	pcPic := e.Value.(*TLibCommon.TComPic)
	if pcPic.GetOutputMark() {
      // write to file
      if this.m_pchReconFile !="" {
        crop := pcPic.GetCroppingWindow();
        this.m_cTVideoIOYuvReconFile.Write( pcPic.GetPicYuvRec(), crop.GetPicCropLeftOffset(), crop.GetPicCropRightOffset(), crop.GetPicCropTopOffset(), crop.GetPicCropBottomOffset() );
      }
      
      // update POC of display order
      this.m_iPOCLastDisplay = int(pcPic.GetPOC());
      
      // erase non-referenced picture in the reference picture list after display
      if !pcPic.GetSlice(0).IsReferenced() && pcPic.GetReconMark() == true {
//#if !DYN_REF_FREE
        pcPic.SetReconMark(false);
        
        // mark it should be extended later
        pcPic.GetPicYuvRec().SetBorderExtension( false );
        
//#else
//        pcPic->destroy();
//        pcListPic->erase( iterPic );
//        iterPic = pcListPic->begin(); // to the beginning, non-efficient way, have to be revised!
//        continue;
//#endif
      }
      pcPic.SetOutputMark(false);
    }
//#if !DYN_REF_FREE
    if pcPic !=nil {
      pcPic.Destroy();
      //delete pcPic;
      pcPic = nil;
    }
//#endif
	pcListPic.Remove(e)    
  }
  
  //pcListPic.Clear();
  this.m_iPOCLastDisplay = - TLibCommon.MAX_INT;
}

func (this *TAppDecTop) IsNaluWithinTargetDecLayerIdSet( nalu *TLibCommon.InputNALUnit ) bool {
  if this.m_targetDecLayerIdSet.Len() == 0 { // By default, the set is empty, meaning all LayerIds are allowed
    return true;
  }
  for e := this.m_targetDecLayerIdSet.Front(); e != nil; e = e.Next() {
	it := e.Value.(int)
    if int(nalu.GetReservedZero6Bits()) == it {
    	return true;
    }
  }
  return false;
}
