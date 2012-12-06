package TLibCommon

import (
	"container/list"
)

// ====================================================================================================================
// Class definition
// ====================================================================================================================

/// picture class (symbol + YUV buffers)
type TComPic struct{
//private:
  m_uiTLayer			uint;               //  Temporal layer
  m_bUsedByCurr		bool;            //  Used by current picture
  m_bIsLongTerm		bool;            //  IS long term picture
  m_bIsUsedAsLongTerm	bool;      //  long term picture is used as reference before
  m_apcPicSym			*TComPicSym;              //  Symbol
  
  m_apcPicYuv	[2]*TComPicYuv;           //  Texture,  0:org / 1:rec
  
  m_pcPicYuvPred	*TComPicYuv;           //  Prediction
  m_pcPicYuvResi	*TComPicYuv;           //  Residual
  m_bReconstructed	bool;
  m_bNeededForOutput	bool;
  m_uiCurrSliceIdx	uint;         // Index of current slice
  m_pSliceSUMap		*int;
  m_pbValidSlice		*bool;
  m_sliceGranularityForNDBFilter	int;
  m_bIndependentSliceBoundaryForNDBFilter	bool;
  m_bIndependentTileBoundaryForNDBFilter	bool;
  m_pNDBFilterYuvTmp	*TComPicYuv;    //!< temporary picture buffer when non-cross slice/tile boundary in-loop filtering is enabled
  m_bCheckLTMSB		bool;
  m_vSliceCUDataLink	*list.List//std::vector<std::vector<TComDataCU*> > ;

  //SEImessages* m_SEIs; ///< Any SEI messages that have been received.  If !NULL we own the object.
}
//public:
func NewTComPic() *TComPic{
	return &TComPic{}
}

  
func (this *TComPic) Create( iWidth, iHeight int, uiMaxWidth, uiMaxHeight, uiMaxDepth uint, 
						croppingWindow *CroppingWindow, numReorderPics []int, bIsVirtual bool ){
}
func (this *TComPic) Destroy(){
}
/*  
  UInt          getTLayer()                { return m_uiTLayer;   }
  Void          setTLayer( UInt uiTLayer ) { m_uiTLayer = uiTLayer; }

  Bool          getUsedByCurr()             { return m_bUsedByCurr; }
  Void          setUsedByCurr( Bool bUsed ) { m_bUsedByCurr = bUsed; }
  Bool          getIsLongTerm()             { return m_bIsLongTerm; }
  Void          setIsLongTerm( Bool lt ) { m_bIsLongTerm = lt; }
  Bool          getIsUsedAsLongTerm()          { return m_bIsUsedAsLongTerm; }
  Void          setIsUsedAsLongTerm( Bool lt ) { m_bIsUsedAsLongTerm = lt; }
  Void          setCheckLTMSBPresent     (Bool b ) {m_bCheckLTMSB=b;}
  Bool          getCheckLTMSBPresent     () { return m_bCheckLTMSB;}
 */
func (this *TComPic) GetPicSym() *TComPicSym          { 
	return  this.m_apcPicSym;    
}
 
func (this *TComPic) GetSlice(i uint) *TComSlice { 
	return  this.m_apcPicSym.GetSlice(i);  
}
 /* 
  Int           getPOC()              { return  m_apcPicSym->getSlice(m_uiCurrSliceIdx)->getPOC();  }
  TComDataCU*&  getCU( UInt uiCUAddr )  { return  m_apcPicSym->getCU( uiCUAddr ); }
  
  TComPicYuv*   getPicYuvOrg()        { return  m_apcPicYuv[0]; }
  TComPicYuv*   getPicYuvRec()        { return  m_apcPicYuv[1]; }
  
  TComPicYuv*   getPicYuvPred()       { return  m_pcPicYuvPred; }
  TComPicYuv*   getPicYuvResi()       { return  m_pcPicYuvResi; }
  Void          setPicYuvPred( TComPicYuv* pcPicYuv )       { m_pcPicYuvPred = pcPicYuv; }
  Void          setPicYuvResi( TComPicYuv* pcPicYuv )       { m_pcPicYuvResi = pcPicYuv; }
  
  UInt          getNumCUsInFrame()      { return m_apcPicSym->getNumberOfCUsInFrame(); }
  UInt          getNumPartInWidth()     { return m_apcPicSym->getNumPartInWidth();     }
  UInt          getNumPartInHeight()    { return m_apcPicSym->getNumPartInHeight();    }
  UInt          getNumPartInCU()        { return m_apcPicSym->getNumPartition();       }
  UInt          getFrameWidthInCU()     { return m_apcPicSym->getFrameWidthInCU();     }
  UInt          getFrameHeightInCU()    { return m_apcPicSym->getFrameHeightInCU();    }
  UInt          getMinCUWidth()         { return m_apcPicSym->getMinCUWidth();         }
  UInt          getMinCUHeight()        { return m_apcPicSym->getMinCUHeight();        }
  
  UInt          getParPelX(UChar uhPartIdx) { return getParPelX(uhPartIdx); }
  UInt          getParPelY(UChar uhPartIdx) { return getParPelX(uhPartIdx); }
  
  Int           getStride()           { return m_apcPicYuv[1]->getStride(); }
  Int           getCStride()          { return m_apcPicYuv[1]->getCStride(); }
  
  Void          setReconMark (Bool b) { m_bReconstructed = b;     }
  Bool          getReconMark ()       { return m_bReconstructed;  }
  Void          setOutputMark (Bool b) { m_bNeededForOutput = b;     }
  Bool          getOutputMark ()       { return m_bNeededForOutput;  }
 
  Void          compressMotion(); 
  UInt          getCurrSliceIdx()            { return m_uiCurrSliceIdx;                }
  Void          setCurrSliceIdx(UInt i)      { m_uiCurrSliceIdx = i;                   }
  UInt          getNumAllocatedSlice()       {return m_apcPicSym->getNumAllocatedSlice();}
  Void          allocateNewSlice()           {m_apcPicSym->allocateNewSlice();         }
  Void          clearSliceBuffer()           {m_apcPicSym->clearSliceBuffer();         }

  Void          createNonDBFilterInfo   (std::vector<Int> sliceStartAddress, Int sliceGranularityDepth
                                        ,std::vector<Bool>* LFCrossSliceBoundary
                                        ,Int  numTiles = 1
                                        ,Bool bNDBFilterCrossTileBoundary = true);
  Void          createNonDBFilterInfoLCU(Int tileID, Int sliceID, TComDataCU* pcCU, UInt startSU, UInt endSU, Int sliceGranularyDepth, UInt picWidth, UInt picHeight);
  Void          destroyNonDBFilterInfo();

  Bool          getValidSlice                                  (Int sliceID)  {return m_pbValidSlice[sliceID];}
  Bool          getIndependentSliceBoundaryForNDBFilter        ()             {return m_bIndependentSliceBoundaryForNDBFilter;}
  Bool          getIndependentTileBoundaryForNDBFilter         ()             {return m_bIndependentTileBoundaryForNDBFilter; }
  TComPicYuv*   getYuvPicBufferForIndependentBoundaryProcessing()             {return m_pNDBFilterYuvTmp;}
  std::vector<TComDataCU*>& getOneSliceCUDataForNDBFilter      (Int sliceID) { return m_vSliceCUDataLink[sliceID];}

  // transfer ownership of seis to this picture 
  void setSEIs(SEImessages* seis) { m_SEIs = seis; }

  //return the current list of SEI messages associated with this picture.
  //Pointer is valid until this->destroy() is called 
  SEImessages* getSEIs() { return m_SEIs; }

  //return the current list of SEI messages associated with this picture.
  // Pointer is valid until this->destroy() is called
  const SEImessages* getSEIs() const { return m_SEIs; }

};*/ 