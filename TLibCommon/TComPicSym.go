package TLibCommon

import ()

// ====================================================================================================================
// Class definition
// ====================================================================================================================

type TComTile struct {
    //private:
    m_uiTileWidth         uint
    m_uiTileHeight        uint
    m_uiRightEdgePosInCU  uint
    m_uiBottomEdgePosInCU uint
    m_uiFirstCUAddr       uint
}


//public:  
func NewTComTile() *TComTile{
	return &TComTile{}
}
 
func (this *TComTile)  SetTileWidth         ( i uint )            { 
	this.m_uiTileWidth = i; 
}
func (this *TComTile)  GetTileWidth         ()  uint              { 
	return this.m_uiTileWidth; 
}
func (this *TComTile)  SetTileHeight        ( i uint )            { 
	this.m_uiTileHeight = i; 
}
func (this *TComTile)  GetTileHeight        ()  uint              { 
	return this.m_uiTileHeight; 
}
func (this *TComTile)  SetRightEdgePosInCU  ( i uint )            { 
	this.m_uiRightEdgePosInCU = i; 
}
func (this *TComTile)  GetRightEdgePosInCU  ()  uint              { 
	return this.m_uiRightEdgePosInCU; 
}
func (this *TComTile)  SetBottomEdgePosInCU ( i uint )            { 
	this.m_uiBottomEdgePosInCU = i; 
}
func (this *TComTile)  GetBottomEdgePosInCU ()  uint              { 
	return this.m_uiBottomEdgePosInCU; 
}
func (this *TComTile)  SetFirstCUAddr       ( i uint )            { 
	this.m_uiFirstCUAddr = i; 
}
func (this *TComTile)  GetFirstCUAddr       ()  uint              { 
	return this.m_uiFirstCUAddr; 
}


/// picture symbol class
type TComPicSym struct {
    //private:
    m_uiWidthInCU  uint
    m_uiHeightInCU uint

    m_uiMaxCUWidth  uint
    m_uiMaxCUHeight uint
    m_uiMinCUWidth  uint
    m_uiMinCUHeight uint

    m_uhTotalDepth      byte ///< max. depth
    m_uiNumPartitions   uint
    m_uiNumPartInWidth  uint
    m_uiNumPartInHeight uint
    m_uiNumCUsInFrame   uint

    m_apcTComSlice        []*TComSlice
    m_uiNumAllocatedSlice uint
    m_apcTComDataCU		[]*TComDataCU;        ///< array of CU data

    m_iTileBoundaryIndependenceIdr int
    m_iNumColumnsMinus1            int
    m_iNumRowsMinus1               int
    m_apcTComTile                  []*TComTile
    m_puiCUOrderMap                []uint //the map of LCU raster scan address relative to LCU encoding order 
    m_puiTileIdxMap                []uint //the map of the tile index relative to LCU raster scan address 
    m_puiInverseCUOrderMap         []uint

    m_saoParam	*SAOParam;
}

/*
public:
  Void        create  ( Int iPicWidth, Int iPicHeight, UInt uiMaxWidth, UInt uiMaxHeight, UInt uiMaxDepth );
  Void        destroy ();

  TComPicSym  ();*/
func (this *TComPicSym) GetSlice(i uint) *TComSlice {
    return this.m_apcTComSlice[i]
}


func (this *TComPicSym)  GetFrameWidthInCU()       uint{ 
	return this.m_uiWidthInCU;                 
}
func (this *TComPicSym)  GetFrameHeightInCU()      uint{ 
	return this.m_uiHeightInCU;                
}
func (this *TComPicSym)  GetMinCUWidth()           uint{ 
	return this.m_uiMinCUWidth;                
}
func (this *TComPicSym)  GetMinCUHeight()          uint{ 
	return this.m_uiMinCUHeight;               
}
func (this *TComPicSym)  GetNumberOfCUsInFrame()   uint{ 
	return this.m_uiNumCUsInFrame;  
}
func (this *TComPicSym)  GetCU( uiCUAddr uint )  *TComDataCU{ 
	return this.m_apcTComDataCU[uiCUAddr];     
}

func (this *TComPicSym)  SetSlice(p *TComSlice, i uint) { 
	this.m_apcTComSlice[i] = p;           
}
func (this *TComPicSym)  GetNumAllocatedSlice()    uint{ 
	return this.m_uiNumAllocatedSlice;         
}
func (this *TComPicSym)  AllocateNewSlice(){
}
func (this *TComPicSym)  ClearSliceBuffer(){
}
func (this *TComPicSym)  GetNumPartition()         uint{ 
	return this.m_uiNumPartitions;             
}
func (this *TComPicSym)  GetNumPartInWidth()       uint{ 
	return this.m_uiNumPartInWidth;            
}
func (this *TComPicSym)  GetNumPartInHeight()      uint{ 
	return this.m_uiNumPartInHeight;           
}
func (this *TComPicSym)  SetNumColumnsMinus1( i int )                          { 
	this.m_iNumColumnsMinus1 = i; 
}
func (this *TComPicSym)  GetNumColumnsMinus1()     int                            { 
	return this.m_iNumColumnsMinus1; 
}  
func (this *TComPicSym)  SetNumRowsMinus1( i int )                             { 
	this.m_iNumRowsMinus1 = i; 
}
func (this *TComPicSym)  GetNumRowsMinus1()        int                            { 
	return this.m_iNumRowsMinus1; 
}
func (this *TComPicSym)  GetNumTiles()             int                            { 
	return (this.m_iNumRowsMinus1+1)*(this.m_iNumColumnsMinus1+1); 
}
func (this *TComPicSym)  GetTComTile  ( tileIdx uint ) *TComTile                         { 
	return this.m_apcTComTile[tileIdx]; 
}
func (this *TComPicSym)  SetCUOrderMap( encCUOrder, cuAddr int )           { 
	this.m_puiCUOrderMap[encCUOrder] = uint(cuAddr); 
}
func (this *TComPicSym)  GetCUOrderMap( encCUOrder int )     uint                  { 
	if encCUOrder>=int(this.m_uiNumCUsInFrame) {
		return this.m_puiCUOrderMap[this.m_uiNumCUsInFrame];
	}
	
	return this.m_puiCUOrderMap[encCUOrder]; 
}
func (this *TComPicSym)  GetTileIdxMap( i int )              uint                  { 
	return this.m_puiTileIdxMap[i]; 
}
func (this *TComPicSym)  SetInverseCUOrderMap( cuAddr, encCUOrder int )    { 
	this.m_puiInverseCUOrderMap[cuAddr] = uint(encCUOrder); 
}
func (this *TComPicSym)  GetInverseCUOrderMap( cuAddr int )  uint                   { 
	if cuAddr>=int(this.m_uiNumCUsInFrame) {
		return this.m_puiInverseCUOrderMap[this.m_uiNumCUsInFrame]
	}
	
	return this.m_puiInverseCUOrderMap [cuAddr]; 
}
func (this *TComPicSym)  GetPicSCUEncOrder( SCUAddr uint )   uint{
	return this.GetInverseCUOrderMap(int(SCUAddr/this.m_uiNumPartitions))*this.m_uiNumPartitions + SCUAddr%this.m_uiNumPartitions; 
}
func (this *TComPicSym)  GetPicSCUAddr( SCUEncOrder uint )  uint{
  return this.GetCUOrderMap(int(SCUEncOrder/this.m_uiNumPartitions))*this.m_uiNumPartitions + SCUEncOrder%this.m_uiNumPartitions;
}
func (this *TComPicSym)  xCreateTComTileArray(){
  /*this.m_apcTComTile = NewTComTile*[(m_iNumColumnsMinus1+1)*(m_iNumRowsMinus1+1)];
  for( UInt i=0; i<(m_iNumColumnsMinus1+1)*(m_iNumRowsMinus1+1); i++ )
  {
    m_apcTComTile[i] = new TComTile;
  }*/
}
func (this *TComPicSym)  xInitTiles(){
}
func (this *TComPicSym)  xCalculateNxtCUAddr( uiCurrCUAddr uint ) uint{
  var  uiNxtCUAddr, uiTileIdx uint;
  
  //get the tile index for the current LCU
  uiTileIdx = this.GetTileIdxMap(int(uiCurrCUAddr));

  //get the raster scan address for the next LCU
  if uiCurrCUAddr % this.m_uiWidthInCU == this.GetTComTile(uiTileIdx).GetRightEdgePosInCU() && 
     uiCurrCUAddr / this.m_uiWidthInCU == this.GetTComTile(uiTileIdx).GetBottomEdgePosInCU()  {
  //the current LCU is the last LCU of the tile
    if int(uiTileIdx) == (this.m_iNumColumnsMinus1+1)*(this.m_iNumRowsMinus1+1)-1 {
      uiNxtCUAddr = this.m_uiNumCUsInFrame;
    }else{
      uiNxtCUAddr = this.GetTComTile(uiTileIdx+1).GetFirstCUAddr();
    }
  }else{ //the current LCU is not the last LCU of the tile
    if uiCurrCUAddr % this.m_uiWidthInCU == this.GetTComTile(uiTileIdx).GetRightEdgePosInCU() {  //the current LCU is on the rightmost edge of the tile
      uiNxtCUAddr = uiCurrCUAddr + this.m_uiWidthInCU - this.GetTComTile(uiTileIdx).GetTileWidth() + 1;
    }else{
      uiNxtCUAddr = uiCurrCUAddr + 1;
    }
  }

  return uiNxtCUAddr;
}

func (this *TComPicSym) AllocSaoParam(sao *TComSampleAdaptiveOffset) {
}
  
func (this *TComPicSym) SetSaoParam() *SAOParam { 
	return this.m_saoParam; 
}
