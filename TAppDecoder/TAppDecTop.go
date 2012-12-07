package TAppDecoder

import (
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

func (this *TAppDecTop) Destroy() {
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

func (this *TAppDecTop) Decode() {
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
	    if pcPic.GetOutputMark() && pcPic.GetPOC() > this.m_iPOCLastDisplay {
	        not_displayed++;
	    }
	}

	for e := pcListPic.Front(); e != nil; e = e.Next() {
	    pcPic := e.Value.(*TLibCommon.TComPic)
	    if pcPic.GetOutputMark() && (not_displayed >  pcPic.GetNumReorderPics(tId) && pcPic.GetPOC() > this.m_iPOCLastDisplay){
            // write to file
            not_displayed--;
            if this.m_pchReconFile != "" {
              crop := pcPic.GetCroppingWindow();
              this.m_cTVideoIOYuvReconFile.Write( pcPic.GetPicYuvRec(), crop.GetPicCropLeftOffset(), crop.GetPicCropRightOffset(), crop.GetPicCropTopOffset(), crop.GetPicCropBottomOffset() );
            }

            // update POC of display order
            this.m_iPOCLastDisplay = pcPic.GetPOC();

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

func (this *TAppDecTop) xFlushOutput() {
}

func (this *TAppDecTop) isNaluWithinTargetDecLayerIdSet() {
}
