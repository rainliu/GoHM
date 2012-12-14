package TLibCommon

import (

)


// ====================================================================================================================
// Class definition
// ====================================================================================================================

//class TComDataCU;

/// neighbouring pixel access class for one component
type TComPatternParam struct{
//private:
  m_iOffsetLeft		int;
  m_iOffsetRight		int;
  m_iOffsetAbove		int;
  m_iOffsetBottom		int;
  m_piPatternOrigin	*Pel;
  
//public:
  m_iROIWidth			int;
  m_iROIHeight		int;
  m_iPatternStride	int;
}
/*
  
  /// return starting position of buffer
  Pel*  getPatternOrigin()        { return  m_piPatternOrigin; }
  
  /// return starting position of ROI (ROI = &pattern[AboveOffset][LeftOffset])
  __inline Pel*  getROIOrigin()
  {
    return  m_piPatternOrigin + m_iPatternStride * m_iOffsetAbove + m_iOffsetLeft;
  }
  
  /// set parameters from Pel buffer for accessing neighbouring pixels
  Void setPatternParamPel ( Pel*        piTexture,
                           Int         iRoiWidth,
                           Int         iRoiHeight,
                           Int         iStride,
                           Int         iOffsetLeft,
                           Int         iOffsetRight,
                           Int         iOffsetAbove,
                           Int         iOffsetBottom );
  
  /// set parameters of one color component from CU data for accessing neighbouring pixels
  Void setPatternParamCU  ( TComDataCU* pcCU,
                           UChar       iComp,
                           UChar       iRoiWidth,
                           UChar       iRoiHeight,
                           Int         iOffsetLeft,
                           Int         iOffsetRight,
                           Int         iOffsetAbove,
                           Int         iOffsetBottom,
                           UInt        uiPartDepth,
                           UInt        uiAbsZorderIdx );
};
*/
/// neighbouring pixel access class for all components
type TComPattern struct{
//private:
  m_cPatternY	TComPatternParam;
  m_cPatternCb	TComPatternParam;
  m_cPatternCr	TComPatternParam;
  
  m_aucIntraFilter	[5]byte;
}

func NewTComPattern() *TComPattern{
	return &TComPattern{};
}
/*
public:
  
  // ROI & pattern information, (ROI = &pattern[AboveOffset][LeftOffset])
  Pel*  getROIY()                 { return m_cPatternY.getROIOrigin();    }
  Int   getROIYWidth()            { return m_cPatternY.m_iROIWidth;       }
  Int   getROIYHeight()           { return m_cPatternY.m_iROIHeight;      }
  Int   getPatternLStride()       { return m_cPatternY.m_iPatternStride;  }

  // access functions of ADI buffers
  Int*  getAdiOrgBuf              ( Int iCuWidth, Int iCuHeight, Int* piAdiBuf );
  Int*  getAdiCbBuf               ( Int iCuWidth, Int iCuHeight, Int* piAdiBuf );
  Int*  getAdiCrBuf               ( Int iCuWidth, Int iCuHeight, Int* piAdiBuf );
  
  Int*  getPredictorPtr           ( UInt uiDirMode, UInt uiWidthBits, Int* piAdiBuf );
  // -------------------------------------------------------------------------------------------------------------------
  // initialization functions
  // -------------------------------------------------------------------------------------------------------------------
  
  /// set parameters from Pel buffers for accessing neighbouring pixels
  Void initPattern            ( Pel*        piY,
                               Pel*        piCb,
                               Pel*        piCr,
                               Int         iRoiWidth,
                               Int         iRoiHeight,
                               Int         iStride,
                               Int         iOffsetLeft,
                               Int         iOffsetRight,
                               Int         iOffsetAbove,
                               Int         iOffsetBottom );
  
  /// set parameters from CU data for accessing neighbouring pixels
  Void  initPattern           ( TComDataCU* pcCU,
                               UInt        uiPartDepth,
                               UInt        uiAbsPartIdx );
  
  /// set luma parameters from CU data for accessing ADI data
  Void  initAdiPattern        ( TComDataCU* pcCU,
                               UInt        uiZorderIdxInPart,
                               UInt        uiPartDepth,
                               Int*        piAdiBuf,
                               Int         iOrgBufStride,
                               Int         iOrgBufHeight,
                               Bool&       bAbove,
                               Bool&       bLeft
                              ,Bool        bLMmode = false // using for LM chroma or not
                               );
  
  /// set chroma parameters from CU data for accessing ADI data
  Void  initAdiPatternChroma  ( TComDataCU* pcCU,
                               UInt        uiZorderIdxInPart,
                               UInt        uiPartDepth,
                               Int*        piAdiBuf,
                               Int         iOrgBufStride,
                               Int         iOrgBufHeight,
                               Bool&       bAbove,
                               Bool&       bLeft );

private:

  /// padding of unavailable reference samples for intra prediction
  Void  fillReferenceSamples        (Int bitDepth, TComDataCU* pcCU, Pel* piRoiOrigin, Int* piAdiTemp, Bool* bNeighborFlags, Int iNumIntraNeighbor, Int iUnitSize, Int iNumUnitsInCu, Int iTotalUnits, UInt uiCuWidth, UInt uiCuHeight, UInt uiWidth, UInt uiHeight, Int iPicStride, Bool bLMmode = false);
  

  /// constrained intra prediction
  Bool  isAboveLeftAvailable  ( TComDataCU* pcCU, UInt uiPartIdxLT );
  Int   isAboveAvailable      ( TComDataCU* pcCU, UInt uiPartIdxLT, UInt uiPartIdxRT, Bool* bValidFlags );
  Int   isLeftAvailable       ( TComDataCU* pcCU, UInt uiPartIdxLT, UInt uiPartIdxLB, Bool* bValidFlags );
  Int   isAboveRightAvailable ( TComDataCU* pcCU, UInt uiPartIdxLT, UInt uiPartIdxRT, Bool* bValidFlags );
  Int   isBelowLeftAvailable  ( TComDataCU* pcCU, UInt uiPartIdxLT, UInt uiPartIdxLB, Bool* bValidFlags );

};*/