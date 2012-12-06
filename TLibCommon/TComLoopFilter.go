package TLibCommon

import (

)


const DEBLOCK_SMALLEST_BLOCK = 8

// ====================================================================================================================
// Class definition
// ====================================================================================================================

/// deblocking filter class
type TComLoopFilter struct{
//private:
  m_disableDeblockingFilterFlag	bool;
  m_betaOffsetDiv2				int;
  m_tcOffsetDiv2					int;

  m_uiNumPartitions				uint;
  m_aapucBS						[2]*byte;              ///< Bs for [Ver/Hor][Y/U/V][Blk_Idx]
  m_aapbEdgeFilter			[2][3]*bool;
  m_stLFCUParam					LFCUParam;                  ///< status structure
  
  m_bLFCrossTileBoundary			bool;
}
/*
protected:
  /// CU-level deblocking function
  Void xDeblockCU                 ( TComDataCU* pcCU, UInt uiAbsZorderIdx, UInt uiDepth, Int Edge );

  // set / get functions
  Void xSetLoopfilterParam        ( TComDataCU* pcCU, UInt uiAbsZorderIdx );
  // filtering functions
  Void xSetEdgefilterTU           ( TComDataCU* pcCU, UInt absTUPartIdx, UInt uiAbsZorderIdx, UInt uiDepth );
  Void xSetEdgefilterPU           ( TComDataCU* pcCU, UInt uiAbsZorderIdx );
  Void xGetBoundaryStrengthSingle ( TComDataCU* pcCU, UInt uiAbsZorderIdx, Int iDir, UInt uiPartIdx );
  UInt xCalcBsIdx                 ( TComDataCU* pcCU, UInt uiAbsZorderIdx, Int iDir, Int iEdgeIdx, Int iBaseUnitIdx )
  {
    TComPic* const pcPic = pcCU->getPic();
    const UInt uiLCUWidthInBaseUnits = pcPic->getNumPartInWidth();
    if( iDir == 0 )
    {
      return g_auiRasterToZscan[g_auiZscanToRaster[uiAbsZorderIdx] + iBaseUnitIdx * uiLCUWidthInBaseUnits + iEdgeIdx ];
    }
    else
    {
      return g_auiRasterToZscan[g_auiZscanToRaster[uiAbsZorderIdx] + iEdgeIdx * uiLCUWidthInBaseUnits + iBaseUnitIdx ];
    }
  } 
  
  Void xSetEdgefilterMultiple( TComDataCU* pcCU, UInt uiAbsZorderIdx, UInt uiDepth, Int iDir, Int iEdgeIdx, Bool bValue ,UInt uiWidthInBaseUnits = 0, UInt uiHeightInBaseUnits = 0 );
  
  Void xEdgeFilterLuma            ( TComDataCU* pcCU, UInt uiAbsZorderIdx, UInt uiDepth, Int iDir, Int iEdge );
  Void xEdgeFilterChroma          ( TComDataCU* pcCU, UInt uiAbsZorderIdx, UInt uiDepth, Int iDir, Int iEdge );
  
  __inline Void xPelFilterLuma( Pel* piSrc, Int iOffset, Int d, Int beta, Int tc, Bool sw, Bool bPartPNoFilter, Bool bPartQNoFilter, Int iThrCut, Bool bFilterSecondP, Bool bFilterSecondQ);
  __inline Void xPelFilterChroma( Pel* piSrc, Int iOffset, Int tc, Bool bPartPNoFilter, Bool bPartQNoFilter);
  

  __inline Bool xUseStrongFiltering( Int offset, Int d, Int beta, Int tc, Pel* piSrc);
  __inline Int xCalcDP( Pel* piSrc, Int iOffset);
  __inline Int xCalcDQ( Pel* piSrc, Int iOffset);
  
public:
  TComLoopFilter();
  virtual ~TComLoopFilter();
*/  
func (this *TComLoopFilter)  Create                    ( uiMaxCUDepth uint ){
}
func (this *TComLoopFilter)  Destroy                   (){
}
/*  
  /// set configuration
  Void setCfg( Bool bLFCrossTileBoundary );
  
  /// picture-level deblocking filter
  Void loopFilterPic( TComPic* pcPic );
};*/