package TLibCommon

import (
    "container/list"
)

// ====================================================================================================================
// Class definition
// ====================================================================================================================

/// picture class (symbol + YUV buffers)
type TComPic struct {
    //private:
    m_uiTLayer          uint        //  Temporal layer
    m_bUsedByCurr       bool        //  Used by current picture
    m_bIsLongTerm       bool        //  IS long term picture
    m_bIsUsedAsLongTerm bool        //  long term picture is used as reference before
    m_apcPicSym         *TComPicSym //  Symbol

    m_apcPicYuv [2]*TComPicYuv //  Texture,  0:org / 1:rec

    m_pcPicYuvPred                          *TComPicYuv //  Prediction
    m_pcPicYuvResi                          *TComPicYuv //  Residual
    m_bReconstructed                        bool
    m_bNeededForOutput                      bool
    m_uiCurrSliceIdx                        uint // Index of current slice
    m_pSliceSUMap                           *int
    m_pbValidSlice                          *bool
    m_sliceGranularityForNDBFilter          int
    m_bIndependentSliceBoundaryForNDBFilter bool
    m_bIndependentTileBoundaryForNDBFilter  bool
    m_pNDBFilterYuvTmp                      *TComPicYuv //!< temporary picture buffer when non-cross slice/tile boundary in-loop filtering is enabled
    m_bCheckLTMSB                           bool
    m_vSliceCUDataLink                      *list.List //std::vector<std::vector<TComDataCU*> > ;

    //SEImessages* m_SEIs; ///< Any SEI messages that have been received.  If !NULL we own the object.
}

//public:
func NewTComPic() *TComPic {
    return &TComPic{}
}

func (this *TComPic) Create(iWidth, iHeight int, uiMaxWidth, uiMaxHeight, uiMaxDepth uint,
    croppingWindow *CroppingWindow, numReorderPics []int, bIsVirtual bool) {
}
func (this *TComPic) Destroy() {
}

func (this *TComPic) GetTLayer() uint {
    return this.m_uiTLayer
}
func (this *TComPic) SetTLayer(uiTLayer uint) {
    this.m_uiTLayer = uiTLayer
}

func (this *TComPic) GetUsedByCurr() bool         {
    return this.m_bUsedByCurr
}
func (this *TComPic) SetUsedByCurr(bUsed bool)    {
    this.m_bUsedByCurr = bUsed
}
func (this *TComPic) GetIsLongTerm() bool         {
    return this.m_bIsLongTerm
}
func (this *TComPic) SetIsLongTerm(lt bool)       {
    this.m_bIsLongTerm = lt
}
func (this *TComPic) GetIsUsedAsLongTerm() bool   {
    return this.m_bIsUsedAsLongTerm
}
func (this *TComPic) SetIsUsedAsLongTerm(lt bool) {
    this.m_bIsUsedAsLongTerm = lt
}
func (this *TComPic) SetCheckLTMSBPresent(b bool) {
    this.m_bCheckLTMSB = b
}
func (this *TComPic) GetCheckLTMSBPresent() bool  {
    return m_bCheckLTMSB
}

func (this *TComPic) GetPicSym() *TComPicSym {
    return this.m_apcPicSym
}

func (this *TComPic) GetSlice(i uint) *TComSlice {
    return this.m_apcPicSym.GetSlice(i)
}


func (this *TComPic)  Int           getPOC()        uint    { return  m_apcPicSym->getSlice(m_uiCurrSliceIdx)->getPOC();  }
func (this *TComPic)  TComDataCU*&  getCU( uiCUAddr uint ) *TComDataCU { return  m_apcPicSym->getCU( uiCUAddr ); }

func (this *TComPic)  TComPicYuv*   getPicYuvOrg()        { return  m_apcPicYuv[0]; }
func (this *TComPic)  TComPicYuv*   getPicYuvRec()        { return  m_apcPicYuv[1]; }

func (this *TComPic)  TComPicYuv*   getPicYuvPred()       { return  m_pcPicYuvPred; }
func (this *TComPic)  TComPicYuv*   getPicYuvResi()       { return  m_pcPicYuvResi; }
func (this *TComPic)  Void          setPicYuvPred( TComPicYuv* pcPicYuv )       { m_pcPicYuvPred = pcPicYuv; }
func (this *TComPic)  Void          setPicYuvResi( TComPicYuv* pcPicYuv )       { m_pcPicYuvResi = pcPicYuv; }

func (this *TComPic)  UInt          getNumCUsInFrame()      { return m_apcPicSym->getNumberOfCUsInFrame(); }
func (this *TComPic)  UInt          getNumPartInWidth()     { return m_apcPicSym->getNumPartInWidth();     }
func (this *TComPic)  UInt          getNumPartInHeight()    { return m_apcPicSym->getNumPartInHeight();    }
func (this *TComPic)  UInt          getNumPartInCU()        { return m_apcPicSym->getNumPartition();       }
func (this *TComPic)  UInt          getFrameWidthInCU()     { return m_apcPicSym->getFrameWidthInCU();     }
func (this *TComPic)  UInt          getFrameHeightInCU()    { return m_apcPicSym->getFrameHeightInCU();    }
func (this *TComPic)  UInt          getMinCUWidth()         { return m_apcPicSym->getMinCUWidth();         }
func (this *TComPic)  UInt          getMinCUHeight()        { return m_apcPicSym->getMinCUHeight();        }

func (this *TComPic)  UInt          getParPelX(UChar uhPartIdx) { return getParPelX(uhPartIdx); }
func (this *TComPic)  UInt          getParPelY(UChar uhPartIdx) { return getParPelX(uhPartIdx); }

func (this *TComPic)  Int           getStride()           { return m_apcPicYuv[1]->getStride(); }
func (this *TComPic)  Int           getCStride()          { return m_apcPicYuv[1]->getCStride(); }

func (this *TComPic)  Void          setReconMark (Bool b) { m_bReconstructed = b;     }
func (this *TComPic)  Bool          getReconMark ()       { return m_bReconstructed;  }
func (this *TComPic)  Void          setOutputMark (Bool b) { m_bNeededForOutput = b;     }
func (this *TComPic)  Bool          getOutputMark ()       { return m_bNeededForOutput;  }

func (this *TComPic)  Void          compressMotion();
func (this *TComPic)  UInt          getCurrSliceIdx()            { return m_uiCurrSliceIdx;                }
func (this *TComPic)  Void          setCurrSliceIdx(UInt i)      { m_uiCurrSliceIdx = i;                   }
func (this *TComPic)  UInt          getNumAllocatedSlice()       {return m_apcPicSym->getNumAllocatedSlice();}
func (this *TComPic)  Void          allocateNewSlice()           {m_apcPicSym->allocateNewSlice();         }
func (this *TComPic)  Void          clearSliceBuffer()           {m_apcPicSym->clearSliceBuffer();         }

func (this *TComPic)  Void          createNonDBFilterInfo   (std::vector<Int> sliceStartAddress, Int sliceGranularityDepth
                                        ,std::vector<Bool>* LFCrossSliceBoundary
                                        ,Int  numTiles = 1
                                        ,Bool bNDBFilterCrossTileBoundary = true);
func (this *TComPic)  Void          createNonDBFilterInfoLCU(Int tileID, Int sliceID, TComDataCU* pcCU, UInt startSU, UInt endSU, Int sliceGranularyDepth, UInt picWidth, UInt picHeight);
func (this *TComPic)  Void          destroyNonDBFilterInfo();

func (this *TComPic)  Bool          getValidSlice                                  (Int sliceID)  {return m_pbValidSlice[sliceID];}
func (this *TComPic)  Bool          getIndependentSliceBoundaryForNDBFilter        ()             {return m_bIndependentSliceBoundaryForNDBFilter;}
func (this *TComPic)  Bool          getIndependentTileBoundaryForNDBFilter         ()             {return m_bIndependentTileBoundaryForNDBFilter; }
func (this *TComPic)  TComPicYuv*   getYuvPicBufferForIndependentBoundaryProcessing()             {return m_pNDBFilterYuvTmp;}
func (this *TComPic)  std::vector<TComDataCU*>& getOneSliceCUDataForNDBFilter      (Int sliceID) { return m_vSliceCUDataLink[sliceID];}

  // transfer ownership of seis to this picture
func (this *TComPic)  void setSEIs(SEImessages* seis) { m_SEIs = seis; }

  //return the current list of SEI messages associated with this picture.
  //Pointer is valid until this->destroy() is called
func (this *TComPic)  SEImessages* getSEIs() { return m_SEIs; }

  //return the current list of SEI messages associated with this picture.
  // Pointer is valid until this->destroy() is called
func (this *TComPic)  const SEImessages* getSEIs() const { return m_SEIs; }
