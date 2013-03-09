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

import (
    "container/list"
)

// ====================================================================================================================
// Class definition
// ====================================================================================================================

/// Unit block for storing image characteristics
type TEncQPAdaptationUnit struct {
    m_dActivity float64
}

func NewTEncQPAdaptationUnit() *TEncQPAdaptationUnit {
    return &TEncQPAdaptationUnit{}
}

func (this *TEncQPAdaptationUnit) SetActivity(d float64) { this.m_dActivity = d }
func (this *TEncQPAdaptationUnit) GetActivity() float64  { return this.m_dActivity }

/// Local image characteristics for CUs on a specific depth
type TEncPicQPAdaptationLayer struct {
    m_uiAQPartWidth       uint
    m_uiAQPartHeight      uint
    m_uiNumAQPartInWidth  uint
    m_uiNumAQPartInHeight uint
    m_acTEncAQU           []TEncQPAdaptationUnit
    m_dAvgActivity        float64
}

func NewTEncPicQPAdaptationLayer() *TEncPicQPAdaptationLayer {
    return &TEncPicQPAdaptationLayer{}
}

func (this *TEncPicQPAdaptationLayer) create(iWidth, iHeight int, uiAQPartWidth, uiAQPartHeight uint) {
    this.m_uiAQPartWidth = uiAQPartWidth
    this.m_uiAQPartHeight = uiAQPartHeight
    this.m_uiNumAQPartInWidth = (uint(iWidth) + this.m_uiAQPartWidth - 1) / this.m_uiAQPartWidth
    this.m_uiNumAQPartInHeight = (uint(iHeight) + this.m_uiAQPartHeight - 1) / this.m_uiAQPartHeight
    this.m_acTEncAQU = make([]TEncQPAdaptationUnit, this.m_uiNumAQPartInWidth*this.m_uiNumAQPartInHeight)
}
func (this *TEncPicQPAdaptationLayer) destroy() {
    //do nothing
}

func (this *TEncPicQPAdaptationLayer) GetAQPartWidth() uint       { return this.m_uiAQPartWidth }
func (this *TEncPicQPAdaptationLayer) GetAQPartHeight() uint      { return this.m_uiAQPartHeight }
func (this *TEncPicQPAdaptationLayer) GetNumAQPartInWidth() uint  { return this.m_uiNumAQPartInWidth }
func (this *TEncPicQPAdaptationLayer) GetNumAQPartInHeight() uint { return this.m_uiNumAQPartInHeight }
func (this *TEncPicQPAdaptationLayer) GetAQPartStride() uint      { return this.m_uiNumAQPartInWidth }
func (this *TEncPicQPAdaptationLayer) GetQPAdaptationUnit() []TEncQPAdaptationUnit {
    return this.m_acTEncAQU
}
func (this *TEncPicQPAdaptationLayer) GetAvgActivity() float64  { return this.m_dAvgActivity }
func (this *TEncPicQPAdaptationLayer) SetAvgActivity(d float64) { this.m_dAvgActivity = d }

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
    m_pSliceSUMap                           []int
    m_pbValidSlice                          []bool
    m_sliceGranularityForNDBFilter          int
    m_bIndependentSliceBoundaryForNDBFilter bool
    m_bIndependentTileBoundaryForNDBFilter  bool
    m_pNDBFilterYuvTmp                      *TComPicYuv //!< temporary picture buffer when non-cross slice/tile boundary in-loop filtering is enabled
    m_bCheckLTMSB                           bool

    m_numReorderPics       [MAX_TLAYER]int
    m_conformanceWindow    *Window
    m_defaultDisplayWindow *Window

    m_vSliceCUDataLink map[int]*list.List //std::vector<std::vector<TComDataCU*> > ;

    m_SEIs *SEImessages ///< Any SEI messages that have been received.  If !NULL we own the object.

    /// Picture class including local image characteristics information for QP adaptation
    m_acAQLayer    []TEncPicQPAdaptationLayer
    m_uiMaxAQDepth uint
}

//public:
func NewTComPic() *TComPic {
    return &TComPic{}
}

func (this *TComPic) Create(iWidth, iHeight int, uiMaxWidth, uiMaxHeight, uiMaxDepth, uiMaxAQDepth uint,
    conformanceWindow, defaultDisplayWindow *Window, numReorderPics []int, bIsVirtual bool) {
    this.m_apcPicSym = NewTComPicSym()
    this.m_apcPicSym.Create(iWidth, iHeight, uiMaxWidth, uiMaxHeight, uiMaxDepth)
    if !bIsVirtual {
        this.m_apcPicYuv[0] = NewTComPicYuv()
        this.m_apcPicYuv[0].Create(iWidth, iHeight, uiMaxWidth, uiMaxHeight, uiMaxDepth)
    }
    this.m_apcPicYuv[1] = NewTComPicYuv()
    this.m_apcPicYuv[1].Create(iWidth, iHeight, uiMaxWidth, uiMaxHeight, uiMaxDepth)

    /* there are no SEI messages associated with this picture initially */
    this.m_SEIs = nil
    this.m_bUsedByCurr = false

    /* store cropping parameters with picture */
    this.m_conformanceWindow = conformanceWindow
    this.m_defaultDisplayWindow = defaultDisplayWindow

    /* store number of reorder pics with picture */
    for i := 0; i < MAX_TLAYER; i++ {
        this.m_numReorderPics[i] = numReorderPics[i]
    }
    //memcpy(m_numReorderPics, numReorderPics, MAX_TLAYER*sizeof(Int));

    this.m_uiMaxAQDepth = uiMaxAQDepth
    if uiMaxAQDepth > 0 {
        this.m_acAQLayer = make([]TEncPicQPAdaptationLayer, this.m_uiMaxAQDepth)
        for d := uint(0); d < this.m_uiMaxAQDepth; d++ {
            this.m_acAQLayer[d].create(iWidth, iHeight, uiMaxWidth>>d, uiMaxHeight>>d)
        }
    }
}

func (this *TComPic) GetAQLayer(uiDepth uint) *TEncPicQPAdaptationLayer {
    return &this.m_acAQLayer[uiDepth]
}

func (this *TComPic) GetMaxAQDepth() uint { return this.m_uiMaxAQDepth }

func (this *TComPic) XPreanalyze() {
    pcPicYuv := this.GetPicYuvOrg()
    iWidth := pcPicYuv.GetWidth()
    iHeight := pcPicYuv.GetHeight()
    iStride := pcPicYuv.GetStride()

    for d := uint(0); d < this.GetMaxAQDepth(); d++ {
        pLineY := pcPicYuv.GetLumaAddr()
        pcAQLayer := this.GetAQLayer(d)
        uiAQPartWidth := pcAQLayer.GetAQPartWidth()
        uiAQPartHeight := pcAQLayer.GetAQPartHeight()
        pcAQU := pcAQLayer.GetQPAdaptationUnit()

        dSumAct := float64(0.0)
        for y := 0; y < iHeight; y += int(uiAQPartHeight) {
            uiCurrAQPartHeight := uint(MIN(int(uiAQPartHeight), int(iHeight-y)).(int))
            for x := 0; x < iWidth; x += int(uiAQPartWidth) {
                uiCurrAQPartWidth := uint(MIN(int(uiAQPartWidth), int(iWidth-x)).(int))
                pBlkY := pLineY[x:]
                var uiSum = [4]uint{0, 0, 0, 0}
                var uiSumSq = [4]uint{0, 0, 0, 0}
                uiNumPixInAQPart := uint(0)
                by := uint(0)
                for ; by < uiCurrAQPartHeight>>1; by++ {
                    bx := uint(0)
                    for ; bx < uiCurrAQPartWidth>>1; bx++ {
                        uiSum[0] += uint(pBlkY[bx])
                        uiSumSq[0] += uint(pBlkY[bx]) * uint(pBlkY[bx])
                        uiNumPixInAQPart++
                    }
                    for ; bx < uiCurrAQPartWidth; bx++ {
                        uiSum[1] += uint(pBlkY[bx])
                        uiSumSq[1] += uint(pBlkY[bx]) * uint(pBlkY[bx])
                        uiNumPixInAQPart++
                    }
                    pBlkY = pBlkY[iStride:] //+= iStride;
                }
                for ; by < uiCurrAQPartHeight; by++ {
                    bx := uint(0)
                    for ; bx < uiCurrAQPartWidth>>1; bx++ {
                        uiSum[2] += uint(pBlkY[bx])
                        uiSumSq[2] += uint(pBlkY[bx]) * uint(pBlkY[bx])
                        uiNumPixInAQPart++
                    }
                    for ; bx < uiCurrAQPartWidth; bx++ {
                        uiSum[3] += uint(pBlkY[bx])
                        uiSumSq[3] += uint(pBlkY[bx]) * uint(pBlkY[bx])
                        uiNumPixInAQPart++
                    }
                    pBlkY = pBlkY[iStride:] //+= iStride;
                }

                dMinVar := float64(MAX_DOUBLE)
                for i := int(0); i < 4; i++ {
                    dAverage := float64(uiSum[i]) / float64(uiNumPixInAQPart)
                    dVariance := float64(uiSumSq[i])/float64(uiNumPixInAQPart) - dAverage*dAverage
                    if dMinVar > dVariance {
                        dMinVar = dVariance
                    }
                }
                dActivity := 1.0 + dMinVar
                pcAQU[0].SetActivity(dActivity)
                dSumAct += dActivity

                pcAQU = pcAQU[1:] //++
            }
            pLineY = pLineY[uint(iStride)*uiCurrAQPartHeight:] //+= iStride * uiCurrAQPartHeight;
        }

        dAvgAct := dSumAct / float64(pcAQLayer.GetNumAQPartInWidth()*pcAQLayer.GetNumAQPartInHeight())
        pcAQLayer.SetAvgActivity(dAvgAct)
    }
}

func (this *TComPic) Destroy() {
}

func (this *TComPic) GetTLayer() uint {
    return this.m_uiTLayer
}
func (this *TComPic) SetTLayer(uiTLayer uint) {
    this.m_uiTLayer = uiTLayer
}

func (this *TComPic) GetUsedByCurr() bool {
    return this.m_bUsedByCurr
}
func (this *TComPic) SetUsedByCurr(bUsed bool) {
    this.m_bUsedByCurr = bUsed
}
func (this *TComPic) GetIsLongTerm() bool {
    return this.m_bIsLongTerm
}
func (this *TComPic) SetIsLongTerm(lt bool) {
    this.m_bIsLongTerm = lt
}
func (this *TComPic) SetCheckLTMSBPresent(b bool) {
    this.m_bCheckLTMSB = b
}
func (this *TComPic) GetCheckLTMSBPresent() bool {
    return this.m_bCheckLTMSB
}

func (this *TComPic) GetPicSym() *TComPicSym {
    return this.m_apcPicSym
}

func (this *TComPic) GetSlice(i uint) *TComSlice {
    return this.m_apcPicSym.GetSlice(i)
}

func (this *TComPic) GetPOC() int {
    return this.m_apcPicSym.GetSlice(this.m_uiCurrSliceIdx).GetPOC()
}
func (this *TComPic) GetCU(uiCUAddr uint) *TComDataCU {
    return this.m_apcPicSym.GetCU(uiCUAddr)
}

func (this *TComPic) GetPicYuvOrg() *TComPicYuv {
    return this.m_apcPicYuv[0]
}
func (this *TComPic) GetPicYuvRec() *TComPicYuv {
    return this.m_apcPicYuv[1]
}

func (this *TComPic) GetPicYuvPred() *TComPicYuv {
    return this.m_pcPicYuvPred
}
func (this *TComPic) GetPicYuvResi() *TComPicYuv {
    return this.m_pcPicYuvResi
}
func (this *TComPic) SetPicYuvPred(pcPicYuv *TComPicYuv) {
    this.m_pcPicYuvPred = pcPicYuv
}
func (this *TComPic) SetPicYuvResi(pcPicYuv *TComPicYuv) {
    this.m_pcPicYuvResi = pcPicYuv
}

func (this *TComPic) GetNumCUsInFrame() uint {
    return this.m_apcPicSym.GetNumberOfCUsInFrame()
}
func (this *TComPic) GetNumPartInWidth() uint {
    return this.m_apcPicSym.GetNumPartInWidth()
}
func (this *TComPic) GetNumPartInHeight() uint {
    return this.m_apcPicSym.GetNumPartInHeight()
}
func (this *TComPic) GetNumPartInCU() uint {
    return this.m_apcPicSym.GetNumPartition()
}
func (this *TComPic) GetFrameWidthInCU() uint {
    return this.m_apcPicSym.GetFrameWidthInCU()
}
func (this *TComPic) GetFrameHeightInCU() uint {
    return this.m_apcPicSym.GetFrameHeightInCU()
}
func (this *TComPic) GetMinCUWidth() uint {
    return this.m_apcPicSym.GetMinCUWidth()
}
func (this *TComPic) GetMinCUHeight() uint {
    return this.m_apcPicSym.GetMinCUHeight()
}

func (this *TComPic) GetParPelX(uhPartIdx byte) uint {
    return this.GetParPelX(uhPartIdx)
}
func (this *TComPic) GetParPelY(uhPartIdx byte) uint {
    return this.GetParPelX(uhPartIdx)
}

func (this *TComPic) GetStride() int {
    return this.m_apcPicYuv[1].GetStride()
}
func (this *TComPic) GetCStride() int {
    return this.m_apcPicYuv[1].GetCStride()
}

func (this *TComPic) SetReconMark(b bool) {
    this.m_bReconstructed = b
}
func (this *TComPic) GetReconMark() bool {
    return this.m_bReconstructed
}
func (this *TComPic) SetOutputMark(b bool) {
    this.m_bNeededForOutput = b
}
func (this *TComPic) GetOutputMark() bool {
    return this.m_bNeededForOutput
}
func (this *TComPic) SetNumReorderPics(i int, tlayer uint) {
    this.m_numReorderPics[tlayer] = i
}
func (this *TComPic) GetNumReorderPics(tlayer uint) int {
    return this.m_numReorderPics[tlayer]
}

func (this *TComPic) CompressMotion() {
    pPicSym := this.GetPicSym()
    for uiCUAddr := uint(0); uiCUAddr < pPicSym.GetFrameHeightInCU()*pPicSym.GetFrameWidthInCU(); uiCUAddr++ {
        pcCU := pPicSym.GetCU(uiCUAddr)
        pcCU.CompressMV()
    }
}
func (this *TComPic) GetCurrSliceIdx() uint {
    return this.m_uiCurrSliceIdx
}
func (this *TComPic) SetCurrSliceIdx(i uint) {
    this.m_uiCurrSliceIdx = i
}
func (this *TComPic) GetNumAllocatedSlice() uint {
    return this.m_apcPicSym.GetNumAllocatedSlice()
}
func (this *TComPic) AllocateNewSlice() {
    this.m_apcPicSym.AllocateNewSlice()
}
func (this *TComPic) ClearSliceBuffer() {
    this.m_apcPicSym.ClearSliceBuffer()
}

func (this *TComPic) GetConformanceWindow() *Window { return this.m_conformanceWindow }
func (this *TComPic) GetDefDisplayWindow() *Window  { return this.m_defaultDisplayWindow }

func (this *TComPic) CreateNonDBFilterInfo(sliceStartAddress map[int]int, sliceGranularityDepth int, LFCrossSliceBoundary map[int]bool, numTiles int, bNDBFilterCrossTileBoundary bool) {
    maxNumSUInLCU := this.GetNumPartInCU()
    numLCUInPic := this.GetNumCUsInFrame()
    picWidth := this.GetSlice(0).GetSPS().GetPicWidthInLumaSamples()
    picHeight := this.GetSlice(0).GetSPS().GetPicHeightInLumaSamples()
    numLCUsInPicWidth := this.GetFrameWidthInCU()
    numLCUsInPicHeight := this.GetFrameHeightInCU()
    maxNumSUInLCUWidth := this.GetNumPartInWidth()
    maxNumSUInLCUHeight := this.GetNumPartInHeight()
    numSlices := len(sliceStartAddress) - 1
    this.m_bIndependentSliceBoundaryForNDBFilter = false
    if numSlices > 1 {
        for s := 0; s < numSlices; s++ {
            if LFCrossSliceBoundary[s] == false {
                this.m_bIndependentSliceBoundaryForNDBFilter = true
            }
        }
    }
    this.m_sliceGranularityForNDBFilter = sliceGranularityDepth
    if bNDBFilterCrossTileBoundary {
        this.m_bIndependentTileBoundaryForNDBFilter = false
    } else if numTiles > 1 {
        this.m_bIndependentTileBoundaryForNDBFilter = true
    } else {
        this.m_bIndependentTileBoundaryForNDBFilter = false
    }

    this.m_pbValidSlice = make([]bool, numSlices)
    for s := 0; s < numSlices; s++ {
        this.m_pbValidSlice[s] = true
    }
    this.m_pSliceSUMap = make([]int, maxNumSUInLCU*numLCUInPic)

    //initialization
    for i := uint(0); i < maxNumSUInLCU*numLCUInPic; i++ {
        this.m_pSliceSUMap[i] = -1
    }
    for CUAddr := uint(0); CUAddr < numLCUInPic; CUAddr++ {
        pcCU := this.GetCU(CUAddr)
        pcCU.SetSliceSUMap(this.m_pSliceSUMap, int(CUAddr*maxNumSUInLCU))
        pcCU.GetNDBFilterBlocks().Init()
    }
    //this.m_vSliceCUDataLink.clear();
    //this.m_vSliceCUDataLink.resize(numSlices);
    this.m_vSliceCUDataLink = make(map[int]*list.List)

    var startAddr, endAddr, firstCUInStartLCU, startLCU, endLCU, lastCUInEndLCU, uiAddr uint
    var LPelX, TPelY, LCUX, LCUY uint
    var currSU, startSU, endSU uint

    for s := 0; s < numSlices; s++ {
        this.m_vSliceCUDataLink[s] = list.New()

        //1st step: decide the real start address
        startAddr = uint(sliceStartAddress[s])
        endAddr = uint(sliceStartAddress[s+1] - 1)

        startLCU = startAddr / maxNumSUInLCU
        firstCUInStartLCU = startAddr % maxNumSUInLCU

        endLCU = endAddr / maxNumSUInLCU
        lastCUInEndLCU = endAddr % maxNumSUInLCU

        uiAddr = this.m_apcPicSym.GetCUOrderMap(int(startLCU))

        LCUX = this.GetCU(uiAddr).GetCUPelX()
        LCUY = this.GetCU(uiAddr).GetCUPelY()
        LPelX = LCUX + G_auiRasterToPelX[G_auiZscanToRaster[firstCUInStartLCU]]
        TPelY = LCUY + G_auiRasterToPelY[G_auiZscanToRaster[firstCUInStartLCU]]
        currSU = firstCUInStartLCU

        bMoveToNextLCU := false
        bSliceInOneLCU := (startLCU == endLCU)

        for !(LPelX < picWidth) || !(TPelY < picHeight) {
            currSU++

            if bSliceInOneLCU {
                if currSU > lastCUInEndLCU {
                    this.m_pbValidSlice[s] = false
                    break
                }
            }

            if currSU >= maxNumSUInLCU {
                bMoveToNextLCU = true
                break
            }

            LPelX = LCUX + G_auiRasterToPelX[G_auiZscanToRaster[currSU]]
            TPelY = LCUY + G_auiRasterToPelY[G_auiZscanToRaster[currSU]]
        }

        if !this.m_pbValidSlice[s] {
            continue
        }

        if currSU != firstCUInStartLCU {
            if !bMoveToNextLCU {
                firstCUInStartLCU = currSU
            } else {
                startLCU++
                firstCUInStartLCU = 0
                //assert( startLCU < this.GetNumCUsInFrame());
            }
            //assert(startLCU*maxNumSUInLCU + firstCUInStartLCU < endAddr);
        }

        //2nd step: assign NonDBFilterInfo to each processing block
        for i := uint(startLCU); i <= endLCU; i++ {
            if i == startLCU {
                startSU = firstCUInStartLCU
            } else {
                startSU = 0
            }
            if i == endLCU {
                endSU = lastCUInEndLCU
            } else {
                endSU = maxNumSUInLCU - 1
            }

            uiAddr = this.m_apcPicSym.GetCUOrderMap(int(i))
            iTileID := this.m_apcPicSym.GetTileIdxMap(int(uiAddr))

            pcCU := this.GetCU(uiAddr)
            this.m_vSliceCUDataLink[s].PushBack(pcCU)

            this.CreateNonDBFilterInfoLCU(int(iTileID), s, pcCU, startSU, endSU, this.m_sliceGranularityForNDBFilter, picWidth, picHeight)
        }
    }

    //step 3: border availability
    for s := 0; s < numSlices; s++ {
        if !this.m_pbValidSlice[s] {
            continue
        }

        for e := this.m_vSliceCUDataLink[s].Front(); e != nil; e = e.Next() {
            pcCU := e.Value.(*TComDataCU) //this.m_vSliceCUDataLink[s][i];
            uiAddr = pcCU.GetAddr()

            if pcCU.GetPic() == nil {
                continue
            }
            iTileID := this.m_apcPicSym.GetTileIdxMap(int(uiAddr))
            bTopTileBoundary := false
            bDownTileBoundary := false
            bLeftTileBoundary := false
            bRightTileBoundary := false

            if this.m_bIndependentTileBoundaryForNDBFilter {
                //left
                if uiAddr%numLCUsInPicWidth != 0 {
                    bLeftTileBoundary = (this.m_apcPicSym.GetTileIdxMap(int(uiAddr)-1) != iTileID)
                }
                //right
                if (uiAddr % numLCUsInPicWidth) != (numLCUsInPicWidth - 1) {
                    bRightTileBoundary = (this.m_apcPicSym.GetTileIdxMap(int(uiAddr)+1) != iTileID)
                }
                //top
                if uiAddr >= numLCUsInPicWidth {
                    bTopTileBoundary = (this.m_apcPicSym.GetTileIdxMap(int(uiAddr-numLCUsInPicWidth)) != iTileID)
                }
                //down
                if uiAddr+numLCUsInPicWidth < numLCUInPic {
                    bDownTileBoundary = (this.m_apcPicSym.GetTileIdxMap(int(uiAddr+numLCUsInPicWidth)) != iTileID)
                }

            }

            pcCU.SetNDBFilterBlockBorderAvailability(numLCUsInPicWidth, numLCUsInPicHeight, maxNumSUInLCUWidth, maxNumSUInLCUHeight, picWidth, picHeight, LFCrossSliceBoundary,
                bTopTileBoundary, bDownTileBoundary, bLeftTileBoundary, bRightTileBoundary, this.m_bIndependentTileBoundaryForNDBFilter)
        }

    }

    if this.m_bIndependentSliceBoundaryForNDBFilter || this.m_bIndependentTileBoundaryForNDBFilter {
        this.m_pNDBFilterYuvTmp = NewTComPicYuv()
        this.m_pNDBFilterYuvTmp.Create(int(picWidth), int(picHeight), this.GetSlice(0).GetSPS().GetMaxCUWidth(), this.GetSlice(0).GetSPS().GetMaxCUHeight(), this.GetSlice(0).GetSPS().GetMaxCUDepth())
    }
}
func (this *TComPic) CreateNonDBFilterInfoLCU(tileID, sliceID int, pcCU *TComDataCU, startSU, endSU uint, sliceGranularyDepth int, picWidth, picHeight uint) {
    LCUX := pcCU.GetCUPelX()
    LCUY := pcCU.GetCUPelY()
    pCUSliceMap, iCUSliceMapAddr := pcCU.GetSliceSUMap()
    maxNumSUInLCU := this.GetNumPartInCU()
    maxNumSUInSGU := maxNumSUInLCU >> uint(sliceGranularyDepth<<1)
    maxNumSUInLCUWidth := this.GetNumPartInWidth()
    var LPelX, TPelY, currSU uint

    //Get the number of valid NBFilterBLock
    currSU = startSU
    for currSU <= endSU {
        LPelX = LCUX + G_auiRasterToPelX[G_auiZscanToRaster[currSU]]
        TPelY = LCUY + G_auiRasterToPelY[G_auiZscanToRaster[currSU]]

        for !(LPelX < picWidth) || !(TPelY < picHeight) {
            currSU += maxNumSUInSGU
            if currSU >= maxNumSUInLCU || currSU > endSU {
                break
            }
            LPelX = LCUX + G_auiRasterToPelX[G_auiZscanToRaster[currSU]]
            TPelY = LCUY + G_auiRasterToPelY[G_auiZscanToRaster[currSU]]
        }

        if currSU >= maxNumSUInLCU || currSU > endSU {
            break
        }

        NDBFBlock := &NDBFBlockInfo{}

        NDBFBlock.tileID = tileID
        NDBFBlock.sliceID = sliceID
        NDBFBlock.posY = TPelY
        NDBFBlock.posX = LPelX
        NDBFBlock.startSU = currSU

        uiLastValidSU := currSU
        var uiIdx, uiLPelX_su, uiTPelY_su uint
        for uiIdx = currSU; uiIdx < currSU+maxNumSUInSGU; uiIdx++ {
            if uiIdx > endSU {
                break
            }
            uiLPelX_su = LCUX + G_auiRasterToPelX[G_auiZscanToRaster[uiIdx]]
            uiTPelY_su = LCUY + G_auiRasterToPelY[G_auiZscanToRaster[uiIdx]]
            if !(uiLPelX_su < picWidth) || !(uiTPelY_su < picHeight) {
                continue
            }
            pCUSliceMap[iCUSliceMapAddr+int(uiIdx)] = sliceID
            uiLastValidSU = uiIdx
        }
        NDBFBlock.endSU = uiLastValidSU

        rTLSU := G_auiZscanToRaster[NDBFBlock.startSU]
        rBRSU := G_auiZscanToRaster[NDBFBlock.endSU]
        NDBFBlock.widthSU = (rBRSU % maxNumSUInLCUWidth) - (rTLSU % maxNumSUInLCUWidth) + 1
        NDBFBlock.heightSU = uint(rBRSU/maxNumSUInLCUWidth) - uint(rTLSU/maxNumSUInLCUWidth) + 1
        NDBFBlock.width = NDBFBlock.widthSU * this.GetMinCUWidth()
        NDBFBlock.height = NDBFBlock.heightSU * this.GetMinCUHeight()

        pcCU.GetNDBFilterBlocks().PushBack(NDBFBlock)

        currSU += maxNumSUInSGU
    }
}
func (this *TComPic) DestroyNonDBFilterInfo() {
    if this.m_pbValidSlice != nil {
        //delete[] this.m_pbValidSlice;
        this.m_pbValidSlice = nil
    }

    if this.m_pSliceSUMap != nil {
        //delete[] this.m_pSliceSUMap;
        this.m_pSliceSUMap = nil
    }
    for CUAddr := uint(0); CUAddr < this.GetNumCUsInFrame(); CUAddr++ {
        pcCU := this.GetCU(CUAddr)
        pcCU.GetNDBFilterBlocks().Init()
    }

    if this.m_bIndependentSliceBoundaryForNDBFilter || this.m_bIndependentTileBoundaryForNDBFilter {
        this.m_pNDBFilterYuvTmp.Destroy()
        //delete this.m_pNDBFilterYuvTmp;
        this.m_pNDBFilterYuvTmp = nil
    }
}

func (this *TComPic) GetValidSlice(sliceID int) bool {
    return this.m_pbValidSlice[sliceID]
}
func (this *TComPic) GetIndependentSliceBoundaryForNDBFilter() bool {
    return this.m_bIndependentSliceBoundaryForNDBFilter
}
func (this *TComPic) GetIndependentTileBoundaryForNDBFilter() bool {
    return this.m_bIndependentTileBoundaryForNDBFilter
}
func (this *TComPic) GetYuvPicBufferForIndependentBoundaryProcessing() *TComPicYuv {
    return this.m_pNDBFilterYuvTmp
}

/*func (this *TComPic)  GetOneSliceCUDataForNDBFilter      (sliceID int) *list.List{
	return nil;//this.m_vSliceCUDataLink[sliceID];
}*/

// transfer ownership of seis to this picture
func (this *TComPic) SetSEIs(seis *SEImessages) {
    this.m_SEIs = seis
}

//return the current list of SEI messages associated with this picture.
//Pointer is valid until this.destroy() is called
func (this *TComPic) GetSEIs() *SEImessages {
    return this.m_SEIs
}
