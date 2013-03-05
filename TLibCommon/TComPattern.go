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
//"fmt"
)

// ====================================================================================================================
// Class definition
// ====================================================================================================================

//class TComDataCU;

/// neighbouring pixel access class for one component
type TComPatternParam struct {
    //private:
    m_iOffsetLeft     int
    m_iOffsetAbove    int
    m_piPatternOrigin []Pel

    //public:
    m_iROIWidth      int
    m_iROIHeight     int
    m_iPatternStride int
    m_iOffset		 int
}

/// return starting position of buffer
func (this *TComPatternParam) GetPatternOrigin() []Pel {
    return this.m_piPatternOrigin
}

/// return starting position of ROI (ROI = &pattern[AboveOffset][LeftOffset])
func (this *TComPatternParam) GetROIOrigin() []Pel {
    return this.m_piPatternOrigin[this.m_iPatternStride*this.m_iOffsetAbove+this.m_iOffsetLeft:]
}

/// set parameters from Pel buffer for accessing neighbouring pixels
func (this *TComPatternParam) SetPatternParamPel(piTexture []Pel, iRoiWidth, iRoiHeight, iStride, iOffset, iOffsetLeft, iOffsetAbove int) {
    this.m_piPatternOrigin = piTexture
    this.m_iROIWidth = iRoiWidth
    this.m_iROIHeight = iRoiHeight
    this.m_iPatternStride = iStride
    this.m_iOffset = iOffset
    this.m_iOffsetLeft = iOffsetLeft
    this.m_iOffsetAbove = iOffsetAbove
}

/// set parameters of one color component from CU data for accessing neighbouring pixels
func (this *TComPatternParam) SetPatternParamCU(pcCU *TComDataCU, iComp byte, iRoiWidth, iRoiHeight int, iOffsetLeft, iOffsetAbove int, uiAbsPartIdx uint) {
    this.m_iOffsetLeft = iOffsetLeft
    this.m_iOffsetAbove = iOffsetAbove

    this.m_iROIWidth = iRoiWidth
    this.m_iROIHeight = iRoiHeight

    uiAbsZorderIdx := pcCU.GetZorderIdxInCU() + uiAbsPartIdx

    if iComp == 0 {
        this.m_iPatternStride = pcCU.GetPic().GetStride()
        offsetY := pcCU.GetPic().GetPicYuvRec().m_cuOffsetY[pcCU.GetAddr()] + pcCU.GetPic().GetPicYuvRec().m_buOffsetY[G_auiZscanToRaster[uiAbsZorderIdx]]
        this.m_piPatternOrigin = pcCU.GetPic().GetPicYuvRec().GetLumaAddr()[offsetY-this.m_iOffsetAbove*this.m_iPatternStride-this.m_iOffsetLeft:]
    } else {
        this.m_iPatternStride = pcCU.GetPic().GetCStride()
        if iComp == 1 {
            offsetU := pcCU.GetPic().GetPicYuvRec().m_cuOffsetC[pcCU.GetAddr()] + pcCU.GetPic().GetPicYuvRec().m_buOffsetC[G_auiZscanToRaster[uiAbsZorderIdx]]
            this.m_piPatternOrigin = pcCU.GetPic().GetPicYuvRec().GetCbAddr()[offsetU-this.m_iOffsetAbove*this.m_iPatternStride-this.m_iOffsetLeft:]
        } else {
            offsetV := pcCU.GetPic().GetPicYuvRec().m_cuOffsetC[pcCU.GetAddr()] + pcCU.GetPic().GetPicYuvRec().m_buOffsetC[G_auiZscanToRaster[uiAbsZorderIdx]]
            this.m_piPatternOrigin = pcCU.GetPic().GetPicYuvRec().GetCrAddr()[offsetV-this.m_iOffsetAbove*this.m_iPatternStride-this.m_iOffsetLeft:]
        }
    }
}

// ====================================================================================================================
// Tables
// ====================================================================================================================

var m_aucIntraFilter = [5]byte{
    10, //4x4
    7,  //8x8
    1,  //16x16
    0,  //32x32
    10, //64x64
}

/// neighbouring pixel access class for all components
type TComPattern struct {
    //private:
    m_cPatternY  TComPatternParam
    m_cPatternCb TComPatternParam
    m_cPatternCr TComPatternParam

    //m_aucIntraFilter	[5]byte;
}

func NewTComPattern() *TComPattern {
    return &TComPattern{}
}

// ROI & pattern information, (ROI = &pattern[AboveOffset][LeftOffset])
func (this *TComPattern) GetROIY() []Pel {
    return this.m_cPatternY.GetROIOrigin()
}
func (this *TComPattern) GetROIYWidth() int {
    return this.m_cPatternY.m_iROIWidth
}
func (this *TComPattern) GetROIYHeight() int {
    return this.m_cPatternY.m_iROIHeight
}
func (this *TComPattern) GetPatternLStride() int {
    return this.m_cPatternY.m_iPatternStride
}
func (this *TComPattern) GetPatternLOffset() int {
    return this.m_cPatternY.m_iOffset
}
// access functions of ADI buffers
func (this *TComPattern) GetAdiOrgBuf(iCuWidth, iCuHeight int, piAdiBuf []Pel) []Pel {
    return piAdiBuf
}
func (this *TComPattern) GetAdiCbBuf(iCuWidth, iCuHeight int, piAdiBuf []Pel) []Pel {
    return piAdiBuf
}
func (this *TComPattern) GetAdiCrBuf(iCuWidth, iCuHeight int, piAdiBuf []Pel) []Pel {
    return piAdiBuf[(iCuWidth*2+1)*(iCuHeight*2+1):]
}

func (this *TComPattern) GetPredictorPtr(uiDirMode, log2BlkSize uint, piAdiBuf []Pel) []Pel {
    var piSrc []Pel
    //assert(log2BlkSize >= 2 && log2BlkSize < 7);
    diff := MIN(ABS(int(uiDirMode)-HOR_IDX).(int), ABS(int(uiDirMode)-VER_IDX).(int)).(int)
    ucFiltIdx := diff > int(m_aucIntraFilter[log2BlkSize-2])

    if uiDirMode == DC_IDX {
        ucFiltIdx = false //no smoothing for DC or LM chroma
    }

    //assert( ucFiltIdx <= 1 );

    width := 1 << log2BlkSize
    height := 1 << log2BlkSize

    piSrc = this.GetAdiOrgBuf(width, height, piAdiBuf)

    if ucFiltIdx {
        //fmt.Printf("hit\n");
        return piSrc[(2*width+1)*(2*height+1):]
    }

    return piSrc
}

// -------------------------------------------------------------------------------------------------------------------
// initialization functions
// -------------------------------------------------------------------------------------------------------------------
/// set parameters from Pel buffers for accessing neighbouring pixels
func (this *TComPattern) InitPattern(piY []Pel, piCb []Pel, piCr []Pel,
    iRoiWidth, iRoiHeight, iStride, iOffsetY, iOffsetCb, iOffsetCr, iOffsetLeft, iOffsetAbove int) {
    this.m_cPatternY.SetPatternParamPel (piY,  iRoiWidth,    iRoiHeight,    iStride,    iOffsetY,  iOffsetLeft,    iOffsetAbove)
    this.m_cPatternCb.SetPatternParamPel(piCb, iRoiWidth>>1, iRoiHeight>>1, iStride>>1, iOffsetCb, iOffsetLeft>>1, iOffsetAbove>>1)
    this.m_cPatternCr.SetPatternParamPel(piCr, iRoiWidth>>1, iRoiHeight>>1, iStride>>1, iOffsetCr, iOffsetLeft>>1, iOffsetAbove>>1)

    return
}

/// set parameters from CU data for accessing neighbouring pixels
func (this *TComPattern) InitPattern3(pcCU *TComDataCU, uiPartDepth, uiAbsPartIdx uint) {
    uiOffsetLeft := 0
    uiOffsetAbove := 0

    uiWidth := uint(pcCU.GetWidth1(0)) >> uiPartDepth
    uiHeight := uint(pcCU.GetHeight1(0)) >> uiPartDepth

    uiAbsZorderIdx := pcCU.GetZorderIdxInCU() + uiAbsPartIdx
    uiCurrPicPelX := pcCU.GetCUPelX() + G_auiRasterToPelX[G_auiZscanToRaster[uiAbsZorderIdx]]
    uiCurrPicPelY := pcCU.GetCUPelY() + G_auiRasterToPelY[G_auiZscanToRaster[uiAbsZorderIdx]]

    if uiCurrPicPelX != 0 {
        uiOffsetLeft = 1
    }

    if uiCurrPicPelY != 0 {
        uiOffsetAbove = 1
    }

    this.m_cPatternY.SetPatternParamCU(pcCU, 0, int(uiWidth), int(uiHeight), uiOffsetLeft, uiOffsetAbove, uiAbsPartIdx)
    this.m_cPatternCb.SetPatternParamCU(pcCU, 1, int(uiWidth)>>1, int(uiHeight)>>1, uiOffsetLeft, uiOffsetAbove, uiAbsPartIdx)
    this.m_cPatternCr.SetPatternParamCU(pcCU, 2, int(uiWidth)>>1, int(uiHeight)>>1, uiOffsetLeft, uiOffsetAbove, uiAbsPartIdx)
}

/// set luma parameters from CU data for accessing ADI data
func (this *TComPattern) InitAdiPattern(pcCU *TComDataCU, uiZorderIdxInPart, uiPartDepth uint, piAdiBuf []Pel, iOrgBufStride, iOrgBufHeight int,
    bAbove, bLeft *bool, bLMmode bool) {
    var piRoiOrigin []Pel
    var piAdiTemp []Pel
    uiCuWidth := uint(pcCU.GetWidth1(0)) >> uiPartDepth
    uiCuHeight := uint(pcCU.GetHeight1(0)) >> uiPartDepth
    uiCuWidth2 := uiCuWidth << 1
    uiCuHeight2 := uiCuHeight << 1
    var uiWidth, uiHeight uint
    iPicStride := pcCU.GetPic().GetStride()
    iUnitSize := 0
    iNumUnitsInCu := 0
    iTotalUnits := 0
    var bNeighborFlags [4*MAX_NUM_SPU_W + 1]bool
    iNumIntraNeighbor := 0

    var uiPartIdxLT, uiPartIdxRT, uiPartIdxLB uint

    pcCU.DeriveLeftRightTopIdxAdi(&uiPartIdxLT, &uiPartIdxRT, uiZorderIdxInPart, uiPartDepth)
    pcCU.DeriveLeftBottomIdxAdi(&uiPartIdxLB, uiZorderIdxInPart, uiPartDepth)

    iUnitSize = int(pcCU.GetSlice().GetSPS().GetMaxCUWidth() >> pcCU.GetSlice().GetSPS().GetMaxCUDepth())
    iNumUnitsInCu = int(uiCuWidth) / iUnitSize
    iTotalUnits = (iNumUnitsInCu << 2) + 1

    bNeighborFlags[iNumUnitsInCu*2] = this.IsAboveLeftAvailable(pcCU, uiPartIdxLT)
    iNumIntraNeighbor += int(B2U(bNeighborFlags[iNumUnitsInCu*2]))
    iNumIntraNeighbor += this.IsAboveAvailable(pcCU, uiPartIdxLT, uiPartIdxRT, bNeighborFlags[(iNumUnitsInCu*2)+1:])
    iNumIntraNeighbor += this.IsAboveRightAvailable(pcCU, uiPartIdxLT, uiPartIdxRT, bNeighborFlags[(iNumUnitsInCu*3)+1:])
    iNumIntraNeighbor += this.IsLeftAvailable(pcCU, uiPartIdxLT, uiPartIdxLB, bNeighborFlags[iNumUnitsInCu:(iNumUnitsInCu*2)], iNumUnitsInCu)
    iNumIntraNeighbor += this.IsBelowLeftAvailable(pcCU, uiPartIdxLT, uiPartIdxLB, bNeighborFlags[0:iNumUnitsInCu], iNumUnitsInCu)

    *bAbove = true
    *bLeft = true

    uiWidth = uint(uiCuWidth2) + 1
    uiHeight = uint(uiCuHeight2) + 1

    if ((uiWidth << 2) > uint(iOrgBufStride)) || ((uiHeight << 2) > uint(iOrgBufHeight)) {
        return
    }

    pPicYuvRec := pcCU.GetPic().GetPicYuvRec()
    offsetY := pPicYuvRec.m_cuOffsetY[pcCU.GetAddr()] + pPicYuvRec.m_buOffsetY[G_auiZscanToRaster[pcCU.GetZorderIdxInCU()+uiZorderIdxInPart]]
    //fmt.Printf("offsetY=%d\n", offsetY);
    piRoiOrigin = pPicYuvRec.GetBufY()[pPicYuvRec.m_iLumaMarginY*pPicYuvRec.GetStride()+pPicYuvRec.m_iLumaMarginX+offsetY-iPicStride-1:]
    piAdiTemp = piAdiBuf

    this.FillReferenceSamples(G_bitDepthY, piRoiOrigin, piAdiTemp, bNeighborFlags[:], iNumIntraNeighbor, iUnitSize, iNumUnitsInCu, iTotalUnits, uint(uiCuWidth), uint(uiCuHeight), uiWidth, uiHeight, iPicStride, bLMmode)

    var i int
    // generate filtered intra prediction samples
    iBufSize := uiCuHeight2 + uiCuWidth2 + 1 // left and left above border + above and above right border + top left corner = length of 3. filter buffer

    uiWH := uiWidth * uiHeight // number of elements in one buffer

    piFilteredBuf1 := piAdiBuf[uiWH:]       // 1. filter buffer
    piFilteredBuf2 := piFilteredBuf1[uiWH:] // 2. filter buffer
    piFilterBuf := piFilteredBuf2[uiWH:]    // buffer for 2. filtering (sequential)
    piFilterBufN := piFilterBuf[iBufSize:]  // buffer for 1. filtering (sequential)

    l := 0
    // left border from bottom to top
    for i = 0; i < int(uiCuHeight2); i++ {
        piFilterBuf[l] = piAdiTemp[uiWidth*(uiCuHeight2-uint(i))]
        l++
        //fmt.Printf("%x ", piAdiTemp[uiWidth * (uiCuHeight2 - uint(i))]);
    }
    // top left corner
    piFilterBuf[l] = piAdiTemp[0]
    l++
    //fmt.Printf("%x ", piAdiTemp[0]);
    // above border from left to right
    for i = 0; i < int(uiCuWidth2); i++ {
        piFilterBuf[l] = piAdiTemp[1+i]
        l++
        //fmt.Printf("%x ", piAdiTemp[1+i]);
    }
    //fmt.Printf("\n");

    if pcCU.GetSlice().GetSPS().GetUseStrongIntraSmoothing() {
        blkSize := uint(32)
        bottomLeft := piFilterBuf[0]
        topLeft := piFilterBuf[uiCuHeight2]
        topRight := piFilterBuf[iBufSize-1]
        threshold := Pel(1 << uint(G_bitDepthY-5))
        bilinearLeft := ABS(bottomLeft+topLeft-2*piFilterBuf[uiCuHeight]).(Pel) < threshold
        bilinearAbove := ABS(topLeft+topRight-2*piFilterBuf[uiCuHeight2+uiCuHeight]).(Pel) < threshold

        if uiCuWidth >= blkSize && (bilinearLeft && bilinearAbove) {
            shift := G_aucConvertToBit[uiCuWidth] + 3 // log2(uiCuHeight2)
            piFilterBufN[0] = piFilterBuf[0]
            piFilterBufN[uiCuHeight2] = piFilterBuf[uiCuHeight2]
            piFilterBufN[iBufSize-1] = piFilterBuf[iBufSize-1]
            for i = 1; i < int(uiCuHeight2); i++ {
                piFilterBufN[i] = Pel(((int(uiCuHeight2)-i)*int(bottomLeft) + i*int(topLeft) + int(uiCuHeight)) >> uint(shift))
            }

            for i = 1; i < int(uiCuWidth2); i++ {
                piFilterBufN[int(uiCuHeight2)+i] = Pel(((int(uiCuWidth2)-i)*int(topLeft) + i*int(topRight) + int(uiCuWidth)) >> uint(shift))
            }
        } else {
            // 1. filtering with [1 2 1]
            piFilterBufN[0] = piFilterBuf[0]
            piFilterBufN[iBufSize-1] = piFilterBuf[iBufSize-1]
            for i = 1; i < int(iBufSize)-1; i++ {
                piFilterBufN[i] = (piFilterBuf[i-1] + 2*piFilterBuf[i] + piFilterBuf[i+1] + 2) >> 2
            }
        }
    } else {
        // 1. filtering with [1 2 1]
        piFilterBufN[0] = piFilterBuf[0]
        piFilterBufN[iBufSize-1] = piFilterBuf[iBufSize-1]
        for i = 1; i < int(iBufSize)-1; i++ {
            piFilterBufN[i] = (piFilterBuf[i-1] + 2*piFilterBuf[i] + piFilterBuf[i+1] + 2) >> 2
        }
    }

    //fmt.Printf("1: ");
    // fill 1. filter buffer with filtered values
    l = 0
    for i = 0; i < int(uiCuHeight2); i++ {
        piFilteredBuf1[uiWidth*(uiCuHeight2-uint(i))] = piFilterBufN[l]
        l++
        //fmt.Printf("%x ", piFilteredBuf1[uiWidth * (uiCuHeight2 - uint(i))]);
    }
    piFilteredBuf1[0] = piFilterBufN[l]
    //fmt.Printf("%x ", piFilteredBuf1[0]);
    l++
    for i = 0; i < int(uiCuWidth2); i++ {
        piFilteredBuf1[1+i] = piFilterBufN[l]
        l++
        //fmt.Printf("%x ", piFilteredBuf1[1 + i]);
    }
    //fmt.Printf("\n");
}

/// set chroma parameters from CU data for accessing ADI data
func (this *TComPattern) InitAdiPatternChroma(pcCU *TComDataCU, uiZorderIdxInPart, uiPartDepth uint, piAdiBuf []Pel, iOrgBufStride, iOrgBufHeight int,
    bAbove, bLeft *bool, uiChromaId uint) {
    var piRoiOrigin []Pel
    var piAdiTemp []Pel
    uiCuWidth := uint(pcCU.GetWidth1(0)) >> uiPartDepth
    uiCuHeight := uint(pcCU.GetHeight1(0)) >> uiPartDepth
    var uiWidth, uiHeight uint
    iPicStride := pcCU.GetPic().GetCStride()

    iUnitSize := 0
    iNumUnitsInCu := 0
    iTotalUnits := 0
    var bNeighborFlags [4*MAX_NUM_SPU_W + 1]bool
    iNumIntraNeighbor := 0

    var uiPartIdxLT, uiPartIdxRT, uiPartIdxLB uint

    pcCU.DeriveLeftRightTopIdxAdi(&uiPartIdxLT, &uiPartIdxRT, uiZorderIdxInPart, uiPartDepth)
    pcCU.DeriveLeftBottomIdxAdi(&uiPartIdxLB, uiZorderIdxInPart, uiPartDepth)

    iUnitSize = int(pcCU.GetSlice().GetSPS().GetMaxCUWidth()>>pcCU.GetSlice().GetSPS().GetMaxCUDepth()) >> 1 // for chroma
    iNumUnitsInCu = (int(uiCuWidth) / iUnitSize) >> 1    // for chroma
    iTotalUnits = (iNumUnitsInCu << 2) + 1

    bNeighborFlags[iNumUnitsInCu*2] = this.IsAboveLeftAvailable(pcCU, uiPartIdxLT)
    iNumIntraNeighbor += int(B2U(bNeighborFlags[iNumUnitsInCu*2]))
    iNumIntraNeighbor += this.IsAboveAvailable(pcCU, uiPartIdxLT, uiPartIdxRT, bNeighborFlags[(iNumUnitsInCu*2)+1:])
    iNumIntraNeighbor += this.IsAboveRightAvailable(pcCU, uiPartIdxLT, uiPartIdxRT, bNeighborFlags[(iNumUnitsInCu*3)+1:])
    iNumIntraNeighbor += this.IsLeftAvailable(pcCU, uiPartIdxLT, uiPartIdxLB, bNeighborFlags[iNumUnitsInCu:(iNumUnitsInCu*2)], iNumUnitsInCu)
    iNumIntraNeighbor += this.IsBelowLeftAvailable(pcCU, uiPartIdxLT, uiPartIdxLB, bNeighborFlags[0:iNumUnitsInCu], iNumUnitsInCu)

    *bAbove = true
    *bLeft = true

    uiCuWidth = uiCuWidth >> 1   // for chroma
    uiCuHeight = uiCuHeight >> 1 // for chroma

    uiWidth = uiCuWidth*2 + 1
    uiHeight = uiCuHeight*2 + 1

    if (4*uiWidth > uint(iOrgBufStride)) || (4*uiHeight > uint(iOrgBufHeight)) {
        return
    }

    pPicYuvRec := pcCU.GetPic().GetPicYuvRec()

    if uiChromaId == 0 {
        // get Cb pattern
        offsetU := pPicYuvRec.m_cuOffsetC[pcCU.GetAddr()] + pPicYuvRec.m_buOffsetC[G_auiZscanToRaster[pcCU.GetZorderIdxInCU()+uiZorderIdxInPart]]
        piRoiOrigin = pPicYuvRec.GetBufU()[pPicYuvRec.m_iChromaMarginY*pPicYuvRec.GetCStride()+pPicYuvRec.m_iChromaMarginX+offsetU-iPicStride-1:]
        piAdiTemp = piAdiBuf

        this.FillReferenceSamples(G_bitDepthC, piRoiOrigin, piAdiTemp, bNeighborFlags[:], iNumIntraNeighbor, iUnitSize, iNumUnitsInCu, iTotalUnits, uiCuWidth, uiCuHeight, uiWidth, uiHeight, iPicStride, false)
    } else {
        // get Cr pattern
        offsetV := pPicYuvRec.m_cuOffsetC[pcCU.GetAddr()] + pPicYuvRec.m_buOffsetC[G_auiZscanToRaster[pcCU.GetZorderIdxInCU()+uiZorderIdxInPart]]
        piRoiOrigin = pPicYuvRec.GetBufV()[pPicYuvRec.m_iChromaMarginY*pPicYuvRec.GetCStride()+pPicYuvRec.m_iChromaMarginX+offsetV-iPicStride-1:]
        piAdiTemp = piAdiBuf[uiWidth*uiHeight:]

        this.FillReferenceSamples(G_bitDepthC, piRoiOrigin, piAdiTemp, bNeighborFlags[:], iNumIntraNeighbor, iUnitSize, iNumUnitsInCu, iTotalUnits, uiCuWidth, uiCuHeight, uiWidth, uiHeight, iPicStride, false)
    }
}

/// padding of unavailable reference samples for intra prediction
func (this *TComPattern) FillReferenceSamples(bitDepth int, piRoiOrigin2 []Pel, piAdiTemp []Pel, bNeighborFlags []bool,
    iNumIntraNeighbor, iUnitSize, iNumUnitsInCu, iTotalUnits int,
    uiCuWidth, uiCuHeight, uiWidth, uiHeight uint, iPicStride int, bLMmode bool) {
    //piRoiOrigin := piRoiOrigin2[iPicStride+1:];

    var piRoiTemp []Pel
    iDCValue := Pel(1 << uint(bitDepth-1))

    if iNumIntraNeighbor == 0 {
        var i uint
        // Fill border with DC value
        for i = 0; i < (uiWidth); i++ {
            piAdiTemp[i] = iDCValue
        }
        for i = 1; i < (uiHeight); i++ {
            piAdiTemp[i*(uiWidth)] = iDCValue
        }
    } else if iNumIntraNeighbor == iTotalUnits {
        var i uint
        // Fill top-left border with rec. samples
        piRoiTemp = piRoiOrigin2 // - iPicStride - 1;
        piAdiTemp[0] = piRoiTemp[0]

        // Fill left border with rec. samples
        piRoiTemp = piRoiOrigin2[iPicStride:] // - 1;

        /*if bLMmode {
          piRoiTemp --; // move to the second left column
        }*/

        for i = 0; i < uiCuHeight; i++ {
            piAdiTemp[(1+i)*uiWidth] = piRoiTemp[i*uint(iPicStride)]
            //piRoiTemp += iPicStride;
        }

        piRoiTemp = piRoiOrigin2[iPicStride+int(uiCuHeight)*iPicStride:]

        // Fill below left border with rec. samples
        for i = 0; i < uiCuHeight; i++ {
            piAdiTemp[(1+uiCuHeight+i)*uiWidth] = piRoiTemp[i*uint(iPicStride)]
            //piRoiTemp += iPicStride;
        }

        // Fill top border with rec. samples
        piRoiTemp = piRoiOrigin2[1:]
        for i = 0; i < uiCuWidth; i++ {
            piAdiTemp[1+i] = piRoiTemp[i]
        }

        // Fill top right border with rec. samples
        piRoiTemp = piRoiOrigin2[1+uiCuWidth:]
        for i = 0; i < uiCuWidth; i++ {
            piAdiTemp[1+uiCuWidth+i] = piRoiTemp[i]
        }
    } else { // reference samples are partially available
        var i, j int
        iNumUnits2 := iNumUnitsInCu << 1
        iTotalSamples := iTotalUnits * iUnitSize
        var piAdiLine [5 * MAX_CU_SIZE]Pel
        var piAdiLineTemp []Pel
        var pbNeighborFlags []bool
        var iNext, iCurr int
        piRef := Pel(0)

        // Initialize
        for i = 0; i < iTotalSamples; i++ {
            piAdiLine[i] = iDCValue
        }

        // Fill top-left sample
        piRoiTemp = piRoiOrigin2[:]
        piAdiLineTemp = piAdiLine[(iNumUnits2 * iUnitSize):]
        pbNeighborFlags = bNeighborFlags[iNumUnits2:]
        if pbNeighborFlags[0] {
            piAdiLineTemp[0] = piRoiTemp[0]
            for i = 1; i < iUnitSize; i++ {
                piAdiLineTemp[i] = piAdiLineTemp[0]
            }
        }

        // Fill left & below-left samples
        piRoiTemp = piRoiOrigin2[iPicStride:]
        /*if bLMmode {
          piRoiTemp --; // move the second left column
        }*/
        piAdiLineTemp = piAdiLine[:(iNumUnits2 * iUnitSize)]
        pbNeighborFlags = bNeighborFlags[:iNumUnits2]
        for j = 0; j < iNumUnits2; j++ {
            if pbNeighborFlags[iNumUnits2-1-j] {
                for i = 0; i < iUnitSize; i++ {
                    piAdiLineTemp[(iNumUnits2*iUnitSize-1)-i-j*iUnitSize] = piRoiTemp[j*iUnitSize*iPicStride+i*iPicStride]
                }
            }
            //piRoiTemp += iUnitSize*iPicStride;
            //piAdiLineTemp -= iUnitSize;
            //pbNeighborFlags--;
        }

        // Fill above & above-right samples
        piRoiTemp = piRoiOrigin2[1:]
        piAdiLineTemp = piAdiLine[((iNumUnits2 + 1) * iUnitSize):]
        pbNeighborFlags = bNeighborFlags[iNumUnits2+1:]
        for j = 0; j < iNumUnits2; j++ {
            if pbNeighborFlags[j] {
                for i = 0; i < iUnitSize; i++ {
                    piAdiLineTemp[i+j*iUnitSize] = piRoiTemp[i+j*iUnitSize]
                }
            }
            //piRoiTemp += iUnitSize;
            //piAdiLineTemp += iUnitSize;
            //pbNeighborFlags++;
        }

        // Pad reference samples when necessary
        iCurr = 0
        iNext = 1
        piAdiLineTemp = piAdiLine[:]
        for iCurr < iTotalUnits {
            if !bNeighborFlags[iCurr] {
                if iCurr == 0 {
                    for iNext < iTotalUnits && !bNeighborFlags[iNext] {
                        iNext++
                    }
                    piRef = piAdiLine[iNext*iUnitSize]
                    // Pad unavailable samples with new value
                    for iCurr < iNext {
                        for i = 0; i < iUnitSize; i++ {
                            piAdiLineTemp[i] = piRef
                        }
                        piAdiLineTemp = piAdiLineTemp[iUnitSize:]
                        iCurr++
                    }
                } else {
                    piRef = piAdiLine[iCurr*iUnitSize-1]
                    for i = 0; i < iUnitSize; i++ {
                        piAdiLineTemp[i] = piRef
                    }
                    piAdiLineTemp = piAdiLineTemp[iUnitSize:]
                    iCurr++
                }
            } else {
                piAdiLineTemp = piAdiLineTemp[iUnitSize:]
                iCurr++
            }
        }

        // Copy processed samples
        piAdiLineTemp = piAdiLine[int(uiHeight)+iUnitSize-2:]
        for i = 0; i < int(uiWidth); i++ {
            piAdiTemp[i] = piAdiLineTemp[i]
        }
        piAdiLineTemp = piAdiLine[:uiHeight-1]
        for i = 1; i < int(uiHeight); i++ {
            piAdiTemp[i*int(uiWidth)] = piAdiLineTemp[int(uiHeight)-1-i]
        }
    }
}

/// constrained intra prediction
func (this *TComPattern) IsAboveLeftAvailable(pcCU *TComDataCU, uiPartIdxLT uint) bool {
    var bAboveLeftFlag bool
    var uiPartAboveLeft uint
    pcCUAboveLeft := pcCU.GetPUAboveLeft(&uiPartAboveLeft, uiPartIdxLT, true)
    if pcCU.GetSlice().GetPPS().GetConstrainedIntraPred() {
        bAboveLeftFlag = (pcCUAboveLeft != nil && pcCUAboveLeft.GetPredictionMode1(uiPartAboveLeft) == MODE_INTRA)
    } else {
        bAboveLeftFlag = pcCUAboveLeft != nil
    }
    return bAboveLeftFlag
}
func (this *TComPattern) IsAboveAvailable(pcCU *TComDataCU, uiPartIdxLT uint, uiPartIdxRT uint, bValidFlags []bool) int {
    uiRasterPartBegin := G_auiZscanToRaster[uiPartIdxLT]
    uiRasterPartEnd := G_auiZscanToRaster[uiPartIdxRT] + 1
    uiIdxStep := uint(1)
    //Bool *pbValidFlags := bValidFlags;
    iNumIntra := 0

    for uiRasterPart := uiRasterPartBegin; uiRasterPart < uiRasterPartEnd; uiRasterPart += uiIdxStep {
        var uiPartAbove uint
        pcCUAbove := pcCU.GetPUAbove(&uiPartAbove, G_auiRasterToZscan[uiRasterPart], true, false, true)
        if pcCU.GetSlice().GetPPS().GetConstrainedIntraPred() {
            if pcCUAbove != nil && pcCUAbove.GetPredictionMode1(uiPartAbove) == MODE_INTRA {
                iNumIntra++
                bValidFlags[(uiRasterPart-uiRasterPartBegin)/uiIdxStep] = true
            } else {
                bValidFlags[(uiRasterPart-uiRasterPartBegin)/uiIdxStep] = false
            }
        } else {
            if pcCUAbove != nil {
                iNumIntra++
                bValidFlags[(uiRasterPart-uiRasterPartBegin)/uiIdxStep] = true
            } else {
                bValidFlags[(uiRasterPart-uiRasterPartBegin)/uiIdxStep] = false
            }
        }
        //pbValidFlags++;
    }
    return iNumIntra
}
func (this *TComPattern) IsLeftAvailable(pcCU *TComDataCU, uiPartIdxLT uint, uiPartIdxLB uint, bValidFlags []bool, iNumUnitsInCu int) int {
    uiRasterPartBegin := G_auiZscanToRaster[uiPartIdxLT]
    uiRasterPartEnd := G_auiZscanToRaster[uiPartIdxLB] + 1
    uiIdxStep := pcCU.GetPic().GetNumPartInWidth()
    //Bool *pbValidFlags = bValidFlags;
    iNumIntra := 0

    for uiRasterPart := uiRasterPartBegin; uiRasterPart < uiRasterPartEnd; uiRasterPart += uiIdxStep {
        var uiPartLeft uint
        pcCULeft := pcCU.GetPULeft(&uiPartLeft, G_auiRasterToZscan[uiRasterPart], true, true)
        if pcCU.GetSlice().GetPPS().GetConstrainedIntraPred() {
            if pcCULeft != nil && pcCULeft.GetPredictionMode1(uiPartLeft) == MODE_INTRA {
                iNumIntra++
                bValidFlags[iNumUnitsInCu-1-int((uiRasterPart-uiRasterPartBegin)/uiIdxStep)] = true
            } else {
                bValidFlags[iNumUnitsInCu-1-int((uiRasterPart-uiRasterPartBegin)/uiIdxStep)] = false
            }
        } else {
            if pcCULeft != nil {
                iNumIntra++
                bValidFlags[iNumUnitsInCu-1-int((uiRasterPart-uiRasterPartBegin)/uiIdxStep)] = true
            } else {
                bValidFlags[iNumUnitsInCu-1-int((uiRasterPart-uiRasterPartBegin)/uiIdxStep)] = false
            }
        }
        //pbValidFlags--; // opposite direction
    }

    return iNumIntra
}
func (this *TComPattern) IsAboveRightAvailable(pcCU *TComDataCU, uiPartIdxLT uint, uiPartIdxRT uint, bValidFlags []bool) int {
    uiNumUnitsInPU := G_auiZscanToRaster[uiPartIdxRT] - G_auiZscanToRaster[uiPartIdxLT] + 1
    //Bool *pbValidFlags = bValidFlags;
    iNumIntra := 0

    for uiOffset := uint(1); uiOffset <= uiNumUnitsInPU; uiOffset++ {
        var uiPartAboveRight uint
        pcCUAboveRight := pcCU.GetPUAboveRightAdi(&uiPartAboveRight, uiPartIdxRT, uiOffset, true)
        if pcCU.GetSlice().GetPPS().GetConstrainedIntraPred() {
            if pcCUAboveRight != nil && pcCUAboveRight.GetPredictionMode1(uiPartAboveRight) == MODE_INTRA {
                iNumIntra++
                bValidFlags[uiOffset-1] = true
            } else {
                bValidFlags[uiOffset-1] = false
            }
        } else {
            if pcCUAboveRight != nil {
                iNumIntra++
                bValidFlags[uiOffset-1] = true
            } else {
                bValidFlags[uiOffset-1] = false
            }
        }
        //pbValidFlags++;
    }

    return iNumIntra
}
func (this *TComPattern) IsBelowLeftAvailable(pcCU *TComDataCU, uiPartIdxLT uint, uiPartIdxLB uint, bValidFlags []bool, iNumUnitsInCu int) int {
    uiNumUnitsInPU := (G_auiZscanToRaster[uiPartIdxLB]-G_auiZscanToRaster[uiPartIdxLT])/pcCU.GetPic().GetNumPartInWidth() + 1
    //Bool *pbValidFlags = bValidFlags;
    iNumIntra := 0

    for uiOffset := uint(1); uiOffset <= uiNumUnitsInPU; uiOffset++ {
        var uiPartBelowLeft uint
        pcCUBelowLeft := pcCU.GetPUBelowLeftAdi(&uiPartBelowLeft, uiPartIdxLB, uiOffset, true)
        if pcCU.GetSlice().GetPPS().GetConstrainedIntraPred() {
            if pcCUBelowLeft != nil && pcCUBelowLeft.GetPredictionMode1(uiPartBelowLeft) == MODE_INTRA {
                iNumIntra++
                bValidFlags[iNumUnitsInCu-1-int(uiOffset-1)] = true
            } else {
                bValidFlags[iNumUnitsInCu-1-int(uiOffset-1)] = false
            }
        } else {
            if pcCUBelowLeft != nil {
                iNumIntra++
                bValidFlags[iNumUnitsInCu-1-int(uiOffset-1)] = true
            } else {
                bValidFlags[iNumUnitsInCu-1-int(uiOffset-1)] = false
            }
        }
        //pbValidFlags--; // opposite direction
    }

    return iNumIntra
}
