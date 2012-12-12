package TLibCommon

import ()

const DEBLOCK_SMALLEST_BLOCK = 8

// ====================================================================================================================
// Class definition
// ====================================================================================================================

/// deblocking filter class
type TComLoopFilter struct {
    //private:
    m_disableDeblockingFilterFlag bool
    m_betaOffsetDiv2              int
    m_tcOffsetDiv2                int

    m_uiNumPartitions uint
    m_aapucBS         [2]*byte ///< Bs for [Ver/Hor][Y/U/V][Blk_Idx]
    m_aapbEdgeFilter  [2][3]*bool
    m_stLFCUParam     LFCUParam ///< status structure

    m_bLFCrossTileBoundary bool
}

  /// CU-level deblocking function
func (this *TComLoopFilter) xDeblockCU                 ( pcCU *TComDataCU,  uiAbsZorderIdx,  uiDepth uint,  iEdge int ){
}

  // set / get functions
func (this *TComLoopFilter) xSetLoopfilterParam        ( pcCU *TComDataCU,  uiAbsZorderIdx uint ){
}
  // filtering functions
func (this *TComLoopFilter) xSetEdgefilterTU           ( pcCU *TComDataCU,  absTUPartIdx,  uiAbsZorderIdx,  uiDepth uint){
}
func (this *TComLoopFilter) xSetEdgefilterPU           ( pcCU *TComDataCU,  uiAbsZorderIdx uint){
}
func (this *TComLoopFilter) xGetBoundaryStrengthSingle ( pcCU *TComDataCU,  uiAbsZorderIdx uint,  iDir int,  uiPartIdx uint){
}
  
func (this *TComLoopFilter)   xCalcBsIdx           ( pcCU *TComDataCU,  uiAbsZorderIdx uint,  iDir,  iEdgeIdx,  iBaseUnitIdx int) uint{
    pcPic := pcCU.GetPic();
    uiLCUWidthInBaseUnits := pcPic.GetNumPartInWidth();
    if iDir == 0 {
      return G_auiRasterToZscan[G_auiZscanToRaster[uiAbsZorderIdx] + uint(iBaseUnitIdx) * uiLCUWidthInBaseUnits + uint(iEdgeIdx) ];
    }

	return G_auiRasterToZscan[G_auiZscanToRaster[uiAbsZorderIdx] + uint(iEdgeIdx) * uiLCUWidthInBaseUnits + uint(iBaseUnitIdx) ];
}

func (this *TComLoopFilter) xSetEdgefilterMultiple( pcCU *TComDataCU,  uiAbsZorderIdx,  uiDepth uint,  iDir,  iEdge int,  bValue bool, uiWidthInBaseUnits,  uiHeightInBaseUnits uint ){
}

func (this *TComLoopFilter) xEdgeFilterLuma            ( pcCU *TComDataCU,  uiAbsZorderIdx,  uiDepth uint,  iDir,  iEdge int){
}
func (this *TComLoopFilter) xEdgeFilterChroma          ( pcCU *TComDataCU,  uiAbsZorderIdx,  uiDepth uint,  iDir,  iEdge int){
}

func (this *TComLoopFilter) xPelFilterLuma( piSrc *Pel,  iOffset,  d,  beta,  tc int,  sw,  bPartPNoFilter,  bPartQNoFilter bool,  iThrCut int,  bFilterSecondP,  bFilterSecondQ bool){
}
func (this *TComLoopFilter) xPelFilterChroma( piSrc *Pel,  iOffset,  tc int,  bPartPNoFilter,  bPartQNoFilter bool){
}

func (this *TComLoopFilter) xUseStrongFiltering(  offset,  d,  beta,  tc int, piSrc []Pel) bool {
  m4  := piSrc[0];
  m3  := piSrc[-offset];
  m7  := piSrc[ offset*3];
  m0  := piSrc[-offset*4];

  d_strong := int(ABS(m0-m3) + ABS(m7-m4));

  return  (d_strong < (beta>>3)) && (d<(beta>>2)) && ( int(ABS(m3-m4)) < ((tc*5+1)>>1)) ;
	
}
func (this *TComLoopFilter) xCalcDP( piSrc []Pel, iOffset int) Pel{
	return ABS( piSrc[-iOffset*3] - 2*piSrc[-iOffset*2] + piSrc[-iOffset] ) ;
}
func (this *TComLoopFilter) xCalcDQ( piSrc []Pel, iOffset int) Pel{
	return ABS( piSrc[0] - 2*piSrc[iOffset] + piSrc[iOffset*2] );
}


func NewTComLoopFilter() *TComLoopFilter{
	return &TComLoopFilter{};
}

func (this *TComLoopFilter) Create(uiMaxCUDepth uint) {
}
func (this *TComLoopFilter) Destroy() {
}


  /// set configuration
func (this *TComLoopFilter) SetCfg(  bLFCrossTileBoundary bool){
}
  /// picture-level deblocking filter
func (this *TComLoopFilter) LoopFilterPic( pcPic *TComPic){
}

