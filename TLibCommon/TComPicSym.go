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

/*
public:  
  TComTile();
  virtual ~TComTile();

  Void      setTileWidth         ( UInt i )            { m_uiTileWidth = i; }
  UInt      getTileWidth         ()                    { return m_uiTileWidth; }
  Void      setTileHeight        ( UInt i )            { m_uiTileHeight = i; }
  UInt      getTileHeight        ()                    { return m_uiTileHeight; }
  Void      setRightEdgePosInCU  ( UInt i )            { m_uiRightEdgePosInCU = i; }
  UInt      getRightEdgePosInCU  ()                    { return m_uiRightEdgePosInCU; }
  Void      setBottomEdgePosInCU ( UInt i )            { m_uiBottomEdgePosInCU = i; }
  UInt      getBottomEdgePosInCU ()                    { return m_uiBottomEdgePosInCU; }
  Void      setFirstCUAddr       ( UInt i )            { m_uiFirstCUAddr = i; }
  UInt      getFirstCUAddr       ()                    { return m_uiFirstCUAddr; }
};
*/
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
    //m_apcTComDataCU		**TComDataCU;        ///< array of CU data

    m_iTileBoundaryIndependenceIdr int
    m_iNumColumnsMinus1            int
    m_iNumRowsMinus1               int
    m_apcTComTile                  **TComTile
    m_puiCUOrderMap                *uint //the map of LCU raster scan address relative to LCU encoding order 
    m_puiTileIdxMap                *uint //the map of the tile index relative to LCU raster scan address 
    m_puiInverseCUOrderMap         *uint

    //m_saoParam	*SAOParam;
}

/*
public:
  Void        create  ( Int iPicWidth, Int iPicHeight, UInt uiMaxWidth, UInt uiMaxHeight, UInt uiMaxDepth );
  Void        destroy ();

  TComPicSym  ();*/
func (this *TComPicSym) GetSlice(i uint) *TComSlice {
    return this.m_apcTComSlice[i]
}

/* 
  UInt        getFrameWidthInCU()       { return m_uiWidthInCU;                 }
  UInt        getFrameHeightInCU()      { return m_uiHeightInCU;                }
  UInt        getMinCUWidth()           { return m_uiMinCUWidth;                }
  UInt        getMinCUHeight()          { return m_uiMinCUHeight;               }
  UInt        getNumberOfCUsInFrame()   { return m_uiNumCUsInFrame;  }
  TComDataCU*&  getCU( UInt uiCUAddr )  { return m_apcTComDataCU[uiCUAddr];     }

  Void        setSlice(TComSlice* p, UInt i) { m_apcTComSlice[i] = p;           }
  UInt        getNumAllocatedSlice()    { return m_uiNumAllocatedSlice;         }
  Void        allocateNewSlice();
  Void        clearSliceBuffer();
  UInt        getNumPartition()         { return m_uiNumPartitions;             }
  UInt        getNumPartInWidth()       { return m_uiNumPartInWidth;            }
  UInt        getNumPartInHeight()      { return m_uiNumPartInHeight;           }
  Void         setNumColumnsMinus1( Int i )                          { m_iNumColumnsMinus1 = i; }
  Int          getNumColumnsMinus1()                                 { return m_iNumColumnsMinus1; }  
  Void         setNumRowsMinus1( Int i )                             { m_iNumRowsMinus1 = i; }
  Int          getNumRowsMinus1()                                    { return m_iNumRowsMinus1; }
  Int          getNumTiles()                                         { return (m_iNumRowsMinus1+1)*(m_iNumColumnsMinus1+1); }
  TComTile*    getTComTile  ( UInt tileIdx )                         { return *(m_apcTComTile + tileIdx); }
  Void         setCUOrderMap( Int encCUOrder, Int cuAddr )           { *(m_puiCUOrderMap + encCUOrder) = cuAddr; }
  UInt         getCUOrderMap( Int encCUOrder )                       { return *(m_puiCUOrderMap + (encCUOrder>=m_uiNumCUsInFrame ? m_uiNumCUsInFrame : encCUOrder)); }
  UInt         getTileIdxMap( Int i )                                { return *(m_puiTileIdxMap + i); }
  Void         setInverseCUOrderMap( Int cuAddr, Int encCUOrder )    { *(m_puiInverseCUOrderMap + cuAddr) = encCUOrder; }
  UInt         getInverseCUOrderMap( Int cuAddr )                    { return *(m_puiInverseCUOrderMap + (cuAddr>=m_uiNumCUsInFrame ? m_uiNumCUsInFrame : cuAddr)); }
  UInt         getPicSCUEncOrder( UInt SCUAddr );
  UInt         getPicSCUAddr( UInt SCUEncOrder );
  Void         xCreateTComTileArray();
  Void         xInitTiles();
  UInt         xCalculateNxtCUAddr( UInt uiCurrCUAddr );
*/
func (this *TComPicSym) AllocSaoParam(sao *TComSampleAdaptiveOffset) {
}

/*  
  SAOParam *getSaoParam() { return m_saoParam; }
};// END CLASS DEFINITION TComPicSym
*/
