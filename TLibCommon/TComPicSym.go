/* The copyright in this software is being made available under the BSD
 * License, included below. This software may be subject to other third party
 * and contributor rights, including patent rights, and no such rights are
 * granted under this license.
 *
 * Copyright (c) 2012-2013, H265.net
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 *  * Redistributions of source code must retain the above copyright notice,
 *    this list of conditions and the following disclaimer.
 *  * Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *  * Neither the name of the H265.net nor the names of its contributors may
 *    be used to endorse or promote products derived from this software without
 *    specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS
 * BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF
 * THE POSSIBILITY OF SUCH DAMAGE.
 */

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
func NewTComTile() *TComTile {
    return &TComTile{}
}

func (this *TComTile) SetTileWidth(i uint) {
    this.m_uiTileWidth = i
}
func (this *TComTile) GetTileWidth() uint {
    return this.m_uiTileWidth
}
func (this *TComTile) SetTileHeight(i uint) {
    this.m_uiTileHeight = i
}
func (this *TComTile) GetTileHeight() uint {
    return this.m_uiTileHeight
}
func (this *TComTile) SetRightEdgePosInCU(i uint) {
    this.m_uiRightEdgePosInCU = i
}
func (this *TComTile) GetRightEdgePosInCU() uint {
    return this.m_uiRightEdgePosInCU
}
func (this *TComTile) SetBottomEdgePosInCU(i uint) {
    this.m_uiBottomEdgePosInCU = i
}
func (this *TComTile) GetBottomEdgePosInCU() uint {
    return this.m_uiBottomEdgePosInCU
}
func (this *TComTile) SetFirstCUAddr(i uint) {
    this.m_uiFirstCUAddr = i
}
func (this *TComTile) GetFirstCUAddr() uint {
    return this.m_uiFirstCUAddr
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
    m_apcTComDataCU       []*TComDataCU ///< array of CU data

    m_iTileBoundaryIndependenceIdr int
    m_iNumColumnsMinus1            int
    m_iNumRowsMinus1               int
    m_apcTComTile                  []*TComTile
    m_puiCUOrderMap                []uint //the map of LCU raster scan address relative to LCU encoding order
    m_puiTileIdxMap                []uint //the map of the tile index relative to LCU raster scan address
    m_puiInverseCUOrderMap         []uint

    m_saoParam *SAOParam
}

func NewTComPicSym() *TComPicSym {
    return &TComPicSym{}
}

func (this *TComPicSym) Create(iPicWidth, iPicHeight int, uiMaxWidth, uiMaxHeight, uiMaxDepth uint) {
    this.m_uhTotalDepth = byte(uiMaxDepth)
    this.m_uiNumPartitions = 1 << (this.m_uhTotalDepth << 1)

    this.m_uiMaxCUWidth = uiMaxWidth
    this.m_uiMaxCUHeight = uiMaxHeight

    this.m_uiMinCUWidth = uiMaxWidth >> this.m_uhTotalDepth
    this.m_uiMinCUHeight = uiMaxHeight >> this.m_uhTotalDepth

    this.m_uiNumPartInWidth = this.m_uiMaxCUWidth / this.m_uiMinCUWidth
    this.m_uiNumPartInHeight = this.m_uiMaxCUHeight / this.m_uiMinCUHeight

    if uint(iPicWidth)%this.m_uiMaxCUWidth != 0 {
        this.m_uiWidthInCU = uint(iPicWidth)/this.m_uiMaxCUWidth + 1
    } else {
        this.m_uiWidthInCU = uint(iPicWidth) / this.m_uiMaxCUWidth
    }

    if uint(iPicHeight)%this.m_uiMaxCUHeight != 0 {
        this.m_uiHeightInCU = uint(iPicHeight)/this.m_uiMaxCUHeight + 1
    } else {
        this.m_uiHeightInCU = uint(iPicHeight) / this.m_uiMaxCUHeight
    }

    this.m_uiNumCUsInFrame = this.m_uiWidthInCU * this.m_uiHeightInCU
    this.m_apcTComDataCU = make([]*TComDataCU, this.m_uiNumCUsInFrame)

    /*if this.m_uiNumAllocatedSlice>0 {
      for i:=0; i<m_uiNumAllocatedSlice ; i++ )
      {
        delete m_apcTComSlice[i];
      }
      delete [] m_apcTComSlice;
    }*/
    this.m_apcTComSlice = make([]*TComSlice, this.m_uiNumCUsInFrame*this.m_uiNumPartitions)
    this.m_apcTComSlice[0] = NewTComSlice()
    this.m_uiNumAllocatedSlice = 1
    for i := uint(0); i < this.m_uiNumCUsInFrame; i++ {
        this.m_apcTComDataCU[i] = NewTComDataCU()
        this.m_apcTComDataCU[i].Create(this.m_uiNumPartitions, this.m_uiMaxCUWidth, this.m_uiMaxCUHeight,
            false, int(this.m_uiMaxCUWidth>>this.m_uhTotalDepth),
            //#if ADAPTIVE_QP_SELECTION
            true)
        //#endif

    }

    this.m_puiCUOrderMap = make([]uint, this.m_uiNumCUsInFrame+1)
    this.m_puiTileIdxMap = make([]uint, this.m_uiNumCUsInFrame)
    this.m_puiInverseCUOrderMap = make([]uint, this.m_uiNumCUsInFrame+1)

    for i := uint(0); i < this.m_uiNumCUsInFrame; i++ {
        this.m_puiCUOrderMap[i] = i
        this.m_puiInverseCUOrderMap[i] = i
    }
    this.m_saoParam = nil
}
func (this *TComPicSym) Destroy() {
}

func (this *TComPicSym) GetSlice(i uint) *TComSlice {
    return this.m_apcTComSlice[i]
}

func (this *TComPicSym) GetFrameWidthInCU() uint {
    return this.m_uiWidthInCU
}
func (this *TComPicSym) GetFrameHeightInCU() uint {
    return this.m_uiHeightInCU
}
func (this *TComPicSym) GetMinCUWidth() uint {
    return this.m_uiMinCUWidth
}
func (this *TComPicSym) GetMinCUHeight() uint {
    return this.m_uiMinCUHeight
}
func (this *TComPicSym) GetNumberOfCUsInFrame() uint {
    return this.m_uiNumCUsInFrame
}
func (this *TComPicSym) GetCU(uiCUAddr uint) *TComDataCU {
    return this.m_apcTComDataCU[uiCUAddr]
}

func (this *TComPicSym) SetSlice(p *TComSlice, i uint) {
    this.m_apcTComSlice[i] = p
}
func (this *TComPicSym) GetNumAllocatedSlice() uint {
    return this.m_uiNumAllocatedSlice
}
func (this *TComPicSym) AllocateNewSlice() {
    this.m_apcTComSlice[this.m_uiNumAllocatedSlice] = NewTComSlice()
    this.m_uiNumAllocatedSlice++
    if this.m_uiNumAllocatedSlice >= 2 {
        this.m_apcTComSlice[this.m_uiNumAllocatedSlice-1].CopySliceInfo(this.m_apcTComSlice[this.m_uiNumAllocatedSlice-2])
        this.m_apcTComSlice[this.m_uiNumAllocatedSlice-1].InitSlice()
    }
}
func (this *TComPicSym) ClearSliceBuffer() {
    for i := uint(1); i < this.m_uiNumAllocatedSlice; i++ {
        this.m_apcTComSlice[i] = nil
    }
    this.m_uiNumAllocatedSlice = 1
}
func (this *TComPicSym) GetNumPartition() uint {
    return this.m_uiNumPartitions
}
func (this *TComPicSym) GetNumPartInWidth() uint {
    return this.m_uiNumPartInWidth
}
func (this *TComPicSym) GetNumPartInHeight() uint {
    return this.m_uiNumPartInHeight
}
func (this *TComPicSym) SetNumColumnsMinus1(i int) {
    this.m_iNumColumnsMinus1 = i
}
func (this *TComPicSym) GetNumColumnsMinus1() int {
    return this.m_iNumColumnsMinus1
}
func (this *TComPicSym) SetNumRowsMinus1(i int) {
    this.m_iNumRowsMinus1 = i
}
func (this *TComPicSym) GetNumRowsMinus1() int {
    return this.m_iNumRowsMinus1
}
func (this *TComPicSym) GetNumTiles() int {
    return (this.m_iNumRowsMinus1 + 1) * (this.m_iNumColumnsMinus1 + 1)
}
func (this *TComPicSym) GetTComTile(tileIdx uint) *TComTile {
    return this.m_apcTComTile[tileIdx]
}
func (this *TComPicSym) SetCUOrderMap(encCUOrder, cuAddr int) {
    this.m_puiCUOrderMap[encCUOrder] = uint(cuAddr)
}
func (this *TComPicSym) GetCUOrderMap(encCUOrder int) uint {
    if encCUOrder >= int(this.m_uiNumCUsInFrame) {
        return this.m_puiCUOrderMap[this.m_uiNumCUsInFrame]
    }

    return this.m_puiCUOrderMap[encCUOrder]
}
func (this *TComPicSym) GetTileIdxMap(i int) uint {
    return this.m_puiTileIdxMap[i]
}
func (this *TComPicSym) SetInverseCUOrderMap(cuAddr, encCUOrder int) {
    this.m_puiInverseCUOrderMap[cuAddr] = uint(encCUOrder)
}
func (this *TComPicSym) GetInverseCUOrderMap(cuAddr int) uint {
    if cuAddr >= int(this.m_uiNumCUsInFrame) {
        return this.m_puiInverseCUOrderMap[this.m_uiNumCUsInFrame]
    }

    return this.m_puiInverseCUOrderMap[cuAddr]
}
func (this *TComPicSym) GetPicSCUEncOrder(SCUAddr uint) uint {
    return this.GetInverseCUOrderMap(int(SCUAddr/this.m_uiNumPartitions))*this.m_uiNumPartitions + SCUAddr%this.m_uiNumPartitions
}
func (this *TComPicSym) GetPicSCUAddr(SCUEncOrder uint) uint {
    return this.GetCUOrderMap(int(SCUEncOrder/this.m_uiNumPartitions))*this.m_uiNumPartitions + SCUEncOrder%this.m_uiNumPartitions
}
func (this *TComPicSym) XCreateTComTileArray() {
    this.m_apcTComTile = make([]*TComTile, (this.m_iNumColumnsMinus1+1)*(this.m_iNumRowsMinus1+1))
    for i := 0; i < (this.m_iNumColumnsMinus1+1)*(this.m_iNumRowsMinus1+1); i++ {
        this.m_apcTComTile[i] = NewTComTile()
    }
}

func (this *TComPicSym) XInitTiles() {
    var uiTileIdx, uiColumnIdx, uiRowIdx, uiRightEdgePosInCU, uiBottomEdgePosInCU uint
    var i, j int

    //initialize each tile of the current picture
    for uiRowIdx = 0; uiRowIdx < uint(this.m_iNumRowsMinus1)+1; uiRowIdx++ {
        for uiColumnIdx = 0; uiColumnIdx < uint(this.m_iNumColumnsMinus1)+1; uiColumnIdx++ {
            uiTileIdx = uiRowIdx*uint(this.m_iNumColumnsMinus1+1) + uiColumnIdx

            //initialize the RightEdgePosInCU for each tile
            uiRightEdgePosInCU = 0
            for i = 0; i <= int(uiColumnIdx); i++ {
                uiRightEdgePosInCU += this.GetTComTile(uiRowIdx*uint(this.m_iNumColumnsMinus1+1) + uint(i)).GetTileWidth()
            }
            this.GetTComTile(uiTileIdx).SetRightEdgePosInCU(uiRightEdgePosInCU - 1)

            //initialize the BottomEdgePosInCU for each tile
            uiBottomEdgePosInCU = 0
            for i = 0; i <= int(uiRowIdx); i++ {
                uiBottomEdgePosInCU += this.GetTComTile(uint(i)*uint(this.m_iNumColumnsMinus1+1) + uiColumnIdx).GetTileHeight()
            }
            this.GetTComTile(uiTileIdx).SetBottomEdgePosInCU(uiBottomEdgePosInCU - 1)

            //initialize the FirstCUAddr for each tile
            this.GetTComTile(uiTileIdx).SetFirstCUAddr((this.GetTComTile(uiTileIdx).GetBottomEdgePosInCU()-this.GetTComTile(uiTileIdx).GetTileHeight()+1)*this.m_uiWidthInCU +
                this.GetTComTile(uiTileIdx).GetRightEdgePosInCU() - this.GetTComTile(uiTileIdx).GetTileWidth() + 1)
        }
    }

    //initialize the TileIdxMap
    for i = 0; i < int(this.m_uiNumCUsInFrame); i++ {
        for j = 0; j < int(this.m_iNumColumnsMinus1+1); j++ {
            if uint(i)%this.m_uiWidthInCU <= this.GetTComTile(uint(j)).GetRightEdgePosInCU() {
                uiColumnIdx = uint(j)
                j = this.m_iNumColumnsMinus1 + 1
            }
        }
        for j = 0; j < this.m_iNumRowsMinus1+1; j++ {
            if uint(i)/this.m_uiWidthInCU <= this.GetTComTile(uint(j*(this.m_iNumColumnsMinus1+1))).GetBottomEdgePosInCU() {
                uiRowIdx = uint(j)
                j = this.m_iNumRowsMinus1 + 1
            }
        }
        this.m_puiTileIdxMap[i] = uiRowIdx*uint(this.m_iNumColumnsMinus1+1) + uiColumnIdx
    }
}

func (this *TComPicSym) XCalculateNxtCUAddr(uiCurrCUAddr uint) uint {
    var uiNxtCUAddr, uiTileIdx uint

    //get the tile index for the current LCU
    uiTileIdx = this.GetTileIdxMap(int(uiCurrCUAddr))

    //get the raster scan address for the next LCU
    if uiCurrCUAddr%this.m_uiWidthInCU == this.GetTComTile(uiTileIdx).GetRightEdgePosInCU() &&
        uiCurrCUAddr/this.m_uiWidthInCU == this.GetTComTile(uiTileIdx).GetBottomEdgePosInCU() {
        //the current LCU is the last LCU of the tile
        if int(uiTileIdx) == (this.m_iNumColumnsMinus1+1)*(this.m_iNumRowsMinus1+1)-1 {
            uiNxtCUAddr = this.m_uiNumCUsInFrame
        } else {
            uiNxtCUAddr = this.GetTComTile(uiTileIdx + 1).GetFirstCUAddr()
        }
    } else { //the current LCU is not the last LCU of the tile
        if uiCurrCUAddr%this.m_uiWidthInCU == this.GetTComTile(uiTileIdx).GetRightEdgePosInCU() { //the current LCU is on the rightmost edge of the tile
            uiNxtCUAddr = uiCurrCUAddr + this.m_uiWidthInCU - this.GetTComTile(uiTileIdx).GetTileWidth() + 1
        } else {
            uiNxtCUAddr = uiCurrCUAddr + 1
        }
    }

    return uiNxtCUAddr
}

func (this *TComPicSym) AllocSaoParam(sao *TComSampleAdaptiveOffset) {
    this.m_saoParam = &SAOParam{}
    sao.AllocSaoParam(this.m_saoParam)
}

func (this *TComPicSym) GetSaoParam() *SAOParam {
    return this.m_saoParam
}
