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

/// general YUV buffer class
type TComYuv struct {
    //private:

    // ------------------------------------------------------------------------------------------------------------------
    //  YUV buffer
    // ------------------------------------------------------------------------------------------------------------------

    m_apiBufY []Pel
    m_apiBufU []Pel
    m_apiBufV []Pel

    // ------------------------------------------------------------------------------------------------------------------
    //  Parameter for general YUV buffer usage
    // ------------------------------------------------------------------------------------------------------------------

    m_iWidth   uint
    m_iHeight  uint
    m_iCWidth  uint
    m_iCHeight uint
}

func (this *TComYuv) getAddrOffset2(uiPartUnitIdx, width uint) int {
    blkX := G_auiRasterToPelX[G_auiZscanToRaster[uiPartUnitIdx]]
    blkY := G_auiRasterToPelY[G_auiZscanToRaster[uiPartUnitIdx]]

    return int(blkX + blkY*width)
}

func (this *TComYuv) getAddrOffset3(iTransUnitIdx, iBlkSize, width uint) int {
    blkX := (iTransUnitIdx * iBlkSize) & (width - 1)
    blkY := (iTransUnitIdx * iBlkSize) & (^(width - 1))

    return int(blkX + blkY*iBlkSize)
}

//public:

func NewTComYuv() *TComYuv {
    return &TComYuv{}
}

// ------------------------------------------------------------------------------------------------------------------
//  Memory management
// ------------------------------------------------------------------------------------------------------------------

func (this *TComYuv) Create(iWidth, iHeight uint) { ///< Create  YUV buffer
    // memory allocation
    this.m_apiBufY = make([]Pel, iWidth*iHeight)
    this.m_apiBufU = make([]Pel, iWidth*iHeight>>2)
    this.m_apiBufV = make([]Pel, iWidth*iHeight>>2)

    // set width and height
    this.m_iWidth = iWidth
    this.m_iHeight = iHeight
    this.m_iCWidth = iWidth >> 1
    this.m_iCHeight = iHeight >> 1
}

func (this *TComYuv) Destroy() { ///< Destroy YUV buffer
    //do nothing
}

func (this *TComYuv) Clear() { ///< clear   YUV buffer
    var x, y uint

    for y = 0; y < this.m_iHeight; y++ {
        for x = 0; x < this.m_iWidth; x++ {
            this.m_apiBufY[y*this.m_iWidth+x] = 0
        }
    }
    for y = 0; y < this.m_iCHeight; y++ {
        for x = 0; x < this.m_iCWidth; x++ {
            this.m_apiBufU[y*this.m_iCWidth+x] = 0
            this.m_apiBufV[y*this.m_iCWidth+x] = 0
        }
    }
}

// ------------------------------------------------------------------------------------------------------------------
//  Copy, load, store YUV buffer
// ------------------------------------------------------------------------------------------------------------------

//  Copy YUV buffer to picture buffer
func (this *TComYuv) CopyToPicYuv(pcPicYuvDst *TComPicYuv, iCuAddr, uiAbsZorderIdx, uiPartDepth, uiPartIdx uint) {
    this.CopyToPicLuma(pcPicYuvDst, iCuAddr, uiAbsZorderIdx, uiPartDepth, uiPartIdx)
    this.CopyToPicChroma(pcPicYuvDst, iCuAddr, uiAbsZorderIdx, uiPartDepth, uiPartIdx)
}
func (this *TComYuv) CopyToPicLuma(pcPicYuvDst *TComPicYuv, iCuAddr, uiAbsZorderIdx, uiPartDepth, uiPartIdx uint) {
    var x, y, iWidth, iHeight int

    iWidth = int(this.m_iWidth >> uiPartDepth)
    iHeight = int(this.m_iHeight >> uiPartDepth)

    pSrc := this.GetLumaAddr2(uiPartIdx, uint(iWidth))
    pDst := pcPicYuvDst.GetLumaAddr2(int(iCuAddr), int(uiAbsZorderIdx))

    iSrcStride := int(this.GetStride())
    iDstStride := pcPicYuvDst.GetStride()

    for y = 0; y < iHeight; y++ {
        for x = 0; x < iWidth; x++ {
            pDst[y*iDstStride+x] = pSrc[y*iSrcStride+x]
        }
    }
}
func (this *TComYuv) CopyToPicChroma(pcPicYuvDst *TComPicYuv, iCuAddr, uiAbsZorderIdx, uiPartDepth, uiPartIdx uint) {
    var x, y, iWidth, iHeight int

    iWidth = int(this.m_iCWidth >> uiPartDepth)
    iHeight = int(this.m_iCHeight >> uiPartDepth)

    pSrcU := this.GetCbAddr2(uiPartIdx, uint(iWidth))
    pSrcV := this.GetCrAddr2(uiPartIdx, uint(iWidth))
    pDstU := pcPicYuvDst.GetCbAddr2(int(iCuAddr), int(uiAbsZorderIdx))
    pDstV := pcPicYuvDst.GetCrAddr2(int(iCuAddr), int(uiAbsZorderIdx))

    iSrcStride := int(this.GetCStride())
    iDstStride := pcPicYuvDst.GetCStride()
    for y = 0; y < iHeight; y++ {
        for x = 0; x < iWidth; x++ {
            pDstU[y*iDstStride+x] = pSrcU[y*iSrcStride+x]
            pDstV[y*iDstStride+x] = pSrcV[y*iSrcStride+x]
        }
    }
}

//  Copy YUV buffer from picture buffer
func (this *TComYuv) CopyFromPicYuv(pcPicYuvSrc *TComPicYuv, iCuAddr, uiAbsZorderIdx uint) {
    this.CopyFromPicLuma(pcPicYuvSrc, iCuAddr, uiAbsZorderIdx)
    this.CopyFromPicChroma(pcPicYuvSrc, iCuAddr, uiAbsZorderIdx)
}
func (this *TComYuv) CopyFromPicLuma(pcPicYuvSrc *TComPicYuv, iCuAddr, uiAbsZorderIdx uint) {
    var x, y uint
    pDst := this.m_apiBufY
    pSrc := pcPicYuvSrc.GetLumaAddr2(int(iCuAddr), int(uiAbsZorderIdx))

    iDstStride := this.GetStride()
    iSrcStride := uint(pcPicYuvSrc.GetStride())
    for y = 0; y < this.m_iHeight; y++ {
        for x = 0; x < this.m_iWidth; x++ {
            pDst[y*iDstStride+x] = pSrc[y*iSrcStride+x]
        }
    }
}
func (this *TComYuv) CopyFromPicChroma(pcPicYuvSrc *TComPicYuv, iCuAddr, uiAbsZorderIdx uint) {
    var x, y uint

    pDstU := this.m_apiBufU
    pDstV := this.m_apiBufV
    pSrcU := pcPicYuvSrc.GetCbAddr2(int(iCuAddr), int(uiAbsZorderIdx))
    pSrcV := pcPicYuvSrc.GetCrAddr2(int(iCuAddr), int(uiAbsZorderIdx))

    iDstStride := this.GetCStride()
    iSrcStride := uint(pcPicYuvSrc.GetCStride())
    for y = 0; y < this.m_iCHeight; y++ {
        for x = 0; x < this.m_iCWidth; x++ {
            pDstU[y*iDstStride+x] = pSrcU[y*iSrcStride+x]
            pDstV[y*iDstStride+x] = pSrcV[y*iSrcStride+x]
        }
    }
}

//  Copy Small YUV buffer to the part of other Big YUV buffer
func (this *TComYuv) CopyToPartYuv(pcYuvDst *TComYuv, uiDstPartIdx uint) {
    this.CopyToPartLuma(pcYuvDst, uiDstPartIdx)
    this.CopyToPartChroma(pcYuvDst, uiDstPartIdx)
}
func (this *TComYuv) CopyToPartLuma(pcYuvDst *TComYuv, uiDstPartIdx uint) {
    var x, y uint

    pSrc := this.m_apiBufY
    pDst := pcYuvDst.GetLumaAddr1(uiDstPartIdx)

    iSrcStride := this.GetStride()
    iDstStride := pcYuvDst.GetStride()
    for y = 0; y < this.m_iHeight; y++ {
        for x = 0; x < this.m_iWidth; x++ {
            pDst[y*iDstStride+x] = pSrc[y*iSrcStride+x]
        }
    }
}
func (this *TComYuv) CopyToPartChroma(pcYuvDst *TComYuv, uiDstPartIdx uint) {
    var x, y uint

    pSrcU := this.m_apiBufU
    pSrcV := this.m_apiBufV
    pDstU := pcYuvDst.GetCbAddr1(uiDstPartIdx)
    pDstV := pcYuvDst.GetCrAddr1(uiDstPartIdx)

    iSrcStride := this.GetCStride()
    iDstStride := pcYuvDst.GetCStride()
    for y = 0; y < this.m_iCHeight; y++ {
        for x = 0; x < this.m_iCWidth; x++ {
            pDstU[y*iDstStride+x] = pSrcU[y*iSrcStride+x]
            pDstV[y*iDstStride+x] = pSrcV[y*iSrcStride+x]
        }
    }
}

//  Copy the part of Big YUV buffer to other Small YUV buffer
func (this *TComYuv) CopyPartToYuv(pcYuvDst *TComYuv, uiSrcPartIdx uint) {
    this.CopyPartToLuma(pcYuvDst, uiSrcPartIdx)
    this.CopyPartToChroma(pcYuvDst, uiSrcPartIdx)
}
func (this *TComYuv) CopyPartToLuma(pcYuvDst *TComYuv, uiSrcPartIdx uint) {
    var x, y uint

    pSrc := this.GetLumaAddr1(uiSrcPartIdx)
    pDst := pcYuvDst.GetLumaAddr1(0)

    iSrcStride := this.GetStride()
    iDstStride := pcYuvDst.GetStride()

    uiHeight := pcYuvDst.GetHeight()
    uiWidth := pcYuvDst.GetWidth()

    for y = 0; y < uiHeight; y++ {
        for x = 0; x < uiWidth; x++ {
            pDst[y*iDstStride+x] = pSrc[y*iSrcStride+x]
        }
    }
}
func (this *TComYuv) CopyPartToChroma(pcYuvDst *TComYuv, uiSrcPartIdx uint) {
    var x, y uint

    pSrcU := this.GetCbAddr1(uiSrcPartIdx)
    pSrcV := this.GetCrAddr1(uiSrcPartIdx)
    pDstU := pcYuvDst.GetCbAddr1(0)
    pDstV := pcYuvDst.GetCrAddr1(0)

    iSrcStride := this.GetCStride()
    iDstStride := pcYuvDst.GetCStride()

    uiCHeight := pcYuvDst.GetCHeight()
    uiCWidth := pcYuvDst.GetCWidth()

    for y = 0; y < uiCHeight; y++ {
        for x = 0; x < uiCWidth; x++ {
            pDstU[y*iDstStride+x] = pSrcU[y*iSrcStride+x]
            pDstV[y*iDstStride+x] = pSrcV[y*iSrcStride+x]
        }
    }
}

//  Copy YUV partition buffer to other YUV partition buffer
func (this *TComYuv) CopyPartToPartYuv(pcYuvDst *TComYuv, uiPartIdx, uiWidth, uiHeight uint) {
    this.CopyPartToPartLuma(pcYuvDst, uiPartIdx, uiWidth, uiHeight)
    this.CopyPartToPartChroma(pcYuvDst, uiPartIdx, uiWidth>>1, uiHeight>>1)
}
func (this *TComYuv) CopyPartToPartLuma(pcYuvDst *TComYuv, uiPartIdx, uiWidth, uiHeight uint) {
    var x, y uint

    pSrc := this.GetLumaAddr1(uiPartIdx)
    pDst := pcYuvDst.GetLumaAddr1(uiPartIdx)
    //if pSrc == pDst {
    //th not a good idea
    //th best would be to fix the caller
    //  return ;
    //}

    iSrcStride := this.GetStride()
    iDstStride := pcYuvDst.GetStride()
    for y = 0; y < uiHeight; y++ {
        for x = 0; x < uiWidth; x++ {
            pDst[y*iDstStride+x] = pSrc[y*iSrcStride+x]
        }
    }
}
func (this *TComYuv) CopyPartToPartChroma(pcYuvDst *TComYuv, uiPartIdx, uiCWidth, uiCHeight uint) {
    var x, y uint

    pSrcU := this.GetCbAddr1(uiPartIdx)
    pSrcV := this.GetCrAddr1(uiPartIdx)
    pDstU := pcYuvDst.GetCbAddr1(uiPartIdx)
    pDstV := pcYuvDst.GetCrAddr1(uiPartIdx)

    //if( pSrcU == pDstU && pSrcV == pDstV)
    //{
    //th not a good idea
    //th best would be to fix the caller
    //  return ;
    //}

    iSrcStride := this.GetCStride()
    iDstStride := pcYuvDst.GetCStride()
    for y = 0; y < uiCHeight; y++ {
        for x = 0; x < uiCWidth; x++ {
            pDstU[y*iDstStride+x] = pSrcU[y*iSrcStride+x]
            pDstV[y*iDstStride+x] = pSrcV[y*iSrcStride+x]
        }
    }
}
func (this *TComYuv) CopyPartToPartChroma2(pcYuvDst *TComYuv, uiPartIdx, uiWidth, uiHeight, chromaId uint) {
    var x, y uint
    if chromaId == 0 {
        pSrcU := this.GetCbAddr1(uiPartIdx)
        pDstU := pcYuvDst.GetCbAddr1(uiPartIdx)
        //if( pSrcU == pDstU)
        //{
        //  return ;
        //}
        iSrcStride := this.GetCStride()
        iDstStride := pcYuvDst.GetCStride()
        for y = 0; y < uiHeight; y++ {
            for x = 0; x < uiWidth; x++ {
                pDstU[y*iDstStride+x] = pSrcU[y*iSrcStride+x]
            }
        }
    } else if chromaId == 1 {
        pSrcV := this.GetCrAddr1(uiPartIdx)
        pDstV := pcYuvDst.GetCrAddr1(uiPartIdx)
        //if( pSrcV == pDstV)
        //{
        //  return ;
        //}
        iSrcStride := this.GetCStride()
        iDstStride := pcYuvDst.GetCStride()
        for y = 0; y < uiHeight; y++ {
            for x = 0; x < uiWidth; x++ {
                pDstV[y*iDstStride+x] = pSrcV[y*iSrcStride+x]
            }
        }
    } else {
        pSrcU := this.GetCbAddr1(uiPartIdx)
        pSrcV := this.GetCrAddr1(uiPartIdx)
        pDstU := pcYuvDst.GetCbAddr1(uiPartIdx)
        pDstV := pcYuvDst.GetCrAddr1(uiPartIdx)

        //if( pSrcU == pDstU && pSrcV == pDstV)
        //{
        //th not a good idea
        //th best would be to fix the caller
        //  return ;
        //}

        iSrcStride := this.GetCStride()
        iDstStride := pcYuvDst.GetCStride()
        for y = 0; y < uiHeight; y++ {
            for x = 0; x < uiWidth; x++ {
                pDstU[y*iDstStride+x] = pSrcU[y*iSrcStride+x]
                pDstV[y*iDstStride+x] = pSrcV[y*iSrcStride+x]
            }
        }
    }
}

// ------------------------------------------------------------------------------------------------------------------
//  Algebraic operation for YUV buffer
// ------------------------------------------------------------------------------------------------------------------

//  Clip(pcYuvSrc0 + pcYuvSrc1) -> m_apiBuf
func (this *TComYuv) AddClip(pcYuvSrc0 *TComYuv, pcYuvSrc1 *TComYuv, uiTrUnitIdx, uiPartSize uint) {
    this.AddClipLuma(pcYuvSrc0, pcYuvSrc1, uiTrUnitIdx, uiPartSize)
    this.AddClipChroma(pcYuvSrc0, pcYuvSrc1, uiTrUnitIdx, uiPartSize>>1)
}
func (this *TComYuv) AddClipLuma(pcYuvSrc0 *TComYuv, pcYuvSrc1 *TComYuv, uiTrUnitIdx, uiPartSize uint) {
    var x, y uint

    pSrc0 := pcYuvSrc0.GetLumaAddr2(uiTrUnitIdx, uiPartSize)
    pSrc1 := pcYuvSrc1.GetLumaAddr2(uiTrUnitIdx, uiPartSize)
    pDst := this.GetLumaAddr2(uiTrUnitIdx, uiPartSize)

    iSrc0Stride := pcYuvSrc0.GetStride()
    iSrc1Stride := pcYuvSrc1.GetStride()
    iDstStride := this.GetStride()
    for y = 0; y < uiPartSize; y++ {
        for x = 0; x < uiPartSize; x++ {
            pDst[y*iDstStride+x] = ClipY(pSrc0[y*iSrc0Stride+x] + pSrc1[y*iSrc1Stride+x])
        }
    }
}
func (this *TComYuv) AddClipChroma(pcYuvSrc0 *TComYuv, pcYuvSrc1 *TComYuv, uiTrUnitIdx, uiPartSize uint) {
    var x, y uint

    pSrcU0 := pcYuvSrc0.GetCbAddr2(uiTrUnitIdx, uiPartSize)
    pSrcU1 := pcYuvSrc1.GetCbAddr2(uiTrUnitIdx, uiPartSize)
    pSrcV0 := pcYuvSrc0.GetCrAddr2(uiTrUnitIdx, uiPartSize)
    pSrcV1 := pcYuvSrc1.GetCrAddr2(uiTrUnitIdx, uiPartSize)
    pDstU := this.GetCbAddr2(uiTrUnitIdx, uiPartSize)
    pDstV := this.GetCrAddr2(uiTrUnitIdx, uiPartSize)

    iSrc0Stride := pcYuvSrc0.GetCStride()
    iSrc1Stride := pcYuvSrc1.GetCStride()
    iDstStride := this.GetCStride()
    for y = 0; y < uiPartSize; y++ {
        for x = 0; x < uiPartSize; x++ {
            pDstU[y*iDstStride+x] = ClipC(pSrcU0[y*iSrc0Stride+x] + pSrcU1[y*iSrc1Stride+x])
            pDstV[y*iDstStride+x] = ClipC(pSrcV0[y*iSrc0Stride+x] + pSrcV1[y*iSrc1Stride+x])
        }
    }
}

//  pcYuvSrc0 - pcYuvSrc1 -> m_apiBuf
func (this *TComYuv) Subtract(pcYuvSrc0 *TComYuv, pcYuvSrc1 *TComYuv, uiTrUnitIdx, uiPartSize uint) {
    this.SubtractLuma(pcYuvSrc0, pcYuvSrc1, uiTrUnitIdx, uiPartSize)
    this.SubtractChroma(pcYuvSrc0, pcYuvSrc1, uiTrUnitIdx, uiPartSize>>1)
}
func (this *TComYuv) SubtractLuma(pcYuvSrc0 *TComYuv, pcYuvSrc1 *TComYuv, uiTrUnitIdx, uiPartSize uint) {
    var x, y uint

    pSrc0 := pcYuvSrc0.GetLumaAddr2(uiTrUnitIdx, uiPartSize)
    pSrc1 := pcYuvSrc1.GetLumaAddr2(uiTrUnitIdx, uiPartSize)
    pDst := this.GetLumaAddr2(uiTrUnitIdx, uiPartSize)

    iSrc0Stride := pcYuvSrc0.GetStride()
    iSrc1Stride := pcYuvSrc1.GetStride()
    iDstStride := this.GetStride()
    for y = 0; y < uiPartSize; y++ {
        for x = 0; x < uiPartSize; x++ {
            pDst[y*iDstStride+x] = pSrc0[y*iSrc0Stride+x] - pSrc1[y*iSrc1Stride+x]
        }
    }
}
func (this *TComYuv) SubtractChroma(pcYuvSrc0 *TComYuv, pcYuvSrc1 *TComYuv, uiTrUnitIdx, uiPartSize uint) {
    var x, y uint

    pSrcU0 := pcYuvSrc0.GetCbAddr2(uiTrUnitIdx, uiPartSize)
    pSrcU1 := pcYuvSrc1.GetCbAddr2(uiTrUnitIdx, uiPartSize)
    pSrcV0 := pcYuvSrc0.GetCrAddr2(uiTrUnitIdx, uiPartSize)
    pSrcV1 := pcYuvSrc1.GetCrAddr2(uiTrUnitIdx, uiPartSize)
    pDstU := this.GetCbAddr2(uiTrUnitIdx, uiPartSize)
    pDstV := this.GetCrAddr2(uiTrUnitIdx, uiPartSize)

    iSrc0Stride := pcYuvSrc0.GetCStride()
    iSrc1Stride := pcYuvSrc1.GetCStride()
    iDstStride := this.GetCStride()
    for y = 0; y < uiPartSize; y++ {
        for x = 0; x < uiPartSize; x++ {
            pDstU[y*iDstStride+x] = pSrcU0[y*iSrc0Stride+x] - pSrcU1[y*iSrc1Stride+x]
            pDstV[y*iDstStride+x] = pSrcV0[y*iSrc0Stride+x] - pSrcV1[y*iSrc1Stride+x]
        }
    }
}

//  (pcYuvSrc0 + pcYuvSrc1)/2 for YUV partition
func (this *TComYuv) AddAvg(pcYuvSrc0 *TComYuv, pcYuvSrc1 *TComYuv, iPartUnitIdx, iWidth, iHeight uint) {
    var x, y uint

    pSrcY0 := pcYuvSrc0.GetLumaAddr1(iPartUnitIdx)
    pSrcU0 := pcYuvSrc0.GetCbAddr1(iPartUnitIdx)
    pSrcV0 := pcYuvSrc0.GetCrAddr1(iPartUnitIdx)

    pSrcY1 := pcYuvSrc1.GetLumaAddr1(iPartUnitIdx)
    pSrcU1 := pcYuvSrc1.GetCbAddr1(iPartUnitIdx)
    pSrcV1 := pcYuvSrc1.GetCrAddr1(iPartUnitIdx)

    pDstY := this.GetLumaAddr1(iPartUnitIdx)
    pDstU := this.GetCbAddr1(iPartUnitIdx)
    pDstV := this.GetCrAddr1(iPartUnitIdx)

    iSrc0Stride := pcYuvSrc0.GetStride()
    iSrc1Stride := pcYuvSrc1.GetStride()
    iDstStride := this.GetStride()

    shiftNum := uint(IF_INTERNAL_PREC + 1 - G_bitDepthY)
    offset := int((1 << (shiftNum - 1)) + 2*IF_INTERNAL_OFFS)

    for y = 0; y < iHeight; y++ {
        for x = 0; x < iWidth; x += 4 {
            pDstY[y*iDstStride+x+0] = ClipY(Pel((int(pSrcY0[y*iSrc0Stride+x+0]) + int(pSrcY1[y*iSrc1Stride+x+0]) + offset) >> shiftNum))
            pDstY[y*iDstStride+x+1] = ClipY(Pel((int(pSrcY0[y*iSrc0Stride+x+1]) + int(pSrcY1[y*iSrc1Stride+x+1]) + offset) >> shiftNum))
            pDstY[y*iDstStride+x+2] = ClipY(Pel((int(pSrcY0[y*iSrc0Stride+x+2]) + int(pSrcY1[y*iSrc1Stride+x+2]) + offset) >> shiftNum))
            pDstY[y*iDstStride+x+3] = ClipY(Pel((int(pSrcY0[y*iSrc0Stride+x+3]) + int(pSrcY1[y*iSrc1Stride+x+3]) + offset) >> shiftNum))
        }
    }

    shiftNum = uint(IF_INTERNAL_PREC + 1 - G_bitDepthC)
    offset = int((1 << (shiftNum - 1)) + 2*IF_INTERNAL_OFFS)

    iSrc0Stride = pcYuvSrc0.GetCStride()
    iSrc1Stride = pcYuvSrc1.GetCStride()
    iDstStride = this.GetCStride()

    iWidth >>= 1
    iHeight >>= 1

    for y = 0; y < iHeight; y++ {
        for x = 0; x < iWidth; x += 2 {
            // note: chroma min width is 2
            pDstU[y*iDstStride+x+0] = ClipC(Pel((int(pSrcU0[y*iSrc0Stride+x+0]) + int(pSrcU1[y*iSrc1Stride+x+0]) + offset) >> shiftNum))
            pDstV[y*iDstStride+x+0] = ClipC(Pel((int(pSrcV0[y*iSrc0Stride+x+0]) + int(pSrcV1[y*iSrc1Stride+x+0]) + offset) >> shiftNum))
            pDstU[y*iDstStride+x+1] = ClipC(Pel((int(pSrcU0[y*iSrc0Stride+x+1]) + int(pSrcU1[y*iSrc1Stride+x+1]) + offset) >> shiftNum))
            pDstV[y*iDstStride+x+1] = ClipC(Pel((int(pSrcV0[y*iSrc0Stride+x+1]) + int(pSrcV1[y*iSrc1Stride+x+1]) + offset) >> shiftNum))
        }
    }
}

//   Remove High frequency
func (this *TComYuv) RemoveHighFreq(pcYuvSrc *TComYuv, uiPartIdx, uiWidth, uiHeight uint) {
    var x, y uint

    pSrc := pcYuvSrc.GetLumaAddr1(uiPartIdx)
    pSrcU := pcYuvSrc.GetCbAddr1(uiPartIdx)
    pSrcV := pcYuvSrc.GetCrAddr1(uiPartIdx)

    pDst := this.GetLumaAddr1(uiPartIdx)
    pDstU := this.GetCbAddr1(uiPartIdx)
    pDstV := this.GetCrAddr1(uiPartIdx)

    iSrcStride := pcYuvSrc.GetStride()
    iDstStride := this.GetStride()

    for y = 0; y < uiHeight; y++ {
        for x = 0; x < uiWidth; x++ {
            //#if DISABLING_CLIP_FOR_BIPREDME
            pDst[y*iDstStride+x] = (pDst[y*iDstStride+x] << 1) - pSrc[y*iSrcStride+x]
            //#else
            //      pDst[x ] = Clip( (pDst[x ]<<1) - pSrc[x ] );
            //#endif
        }
    }

    iSrcStride = pcYuvSrc.GetCStride()
    iDstStride = this.GetCStride()

    uiHeight >>= 1
    uiWidth >>= 1

    for y = 0; y < uiHeight; y++ {
        for x = 0; x < uiWidth; x++ {
            //#if DISABLING_CLIP_FOR_BIPREDME
            pDstU[y*iDstStride+x] = (pDstU[y*iDstStride+x] << 1) - pSrcU[y*iSrcStride+x]
            pDstV[y*iDstStride+x] = (pDstV[y*iDstStride+x] << 1) - pSrcV[y*iSrcStride+x]
            //#else
            //      pDstU[x ] = Clip( (pDstU[x ]<<1) - pSrcU[x ] );
            //      pDstV[x ] = Clip( (pDstV[x ]<<1) - pSrcV[x ] );
            //#endif
        }
    }
}

// ------------------------------------------------------------------------------------------------------------------
//  Access function for YUV buffer
// ------------------------------------------------------------------------------------------------------------------

//  Access starting position of YUV buffer
func (this *TComYuv) GetLumaAddr() []Pel {
    return this.m_apiBufY
}
func (this *TComYuv) GetCbAddr() []Pel {
    return this.m_apiBufU
}
func (this *TComYuv) GetCrAddr() []Pel {
    return this.m_apiBufV
}

//  Access starting position of YUV partition unit buffer
func (this *TComYuv) GetLumaAddr1(iPartUnitIdx uint) []Pel {
    return this.m_apiBufY[this.getAddrOffset2(iPartUnitIdx, this.m_iWidth):]
}
func (this *TComYuv) GetCbAddr1(iPartUnitIdx uint) []Pel {
    return this.m_apiBufU[this.getAddrOffset2(iPartUnitIdx, this.m_iCWidth)>>1:]
}
func (this *TComYuv) GetCrAddr1(iPartUnitIdx uint) []Pel {
    return this.m_apiBufV[this.getAddrOffset2(iPartUnitIdx, this.m_iCWidth)>>1:]
}

//  Access starting position of YUV transform unit buffer
func (this *TComYuv) GetLumaAddr2(iTransUnitIdx, iBlkSize uint) []Pel {
    return this.m_apiBufY[this.getAddrOffset3(iTransUnitIdx, iBlkSize, this.m_iWidth):]
}
func (this *TComYuv) GetCbAddr2(iTransUnitIdx, iBlkSize uint) []Pel {
    return this.m_apiBufU[this.getAddrOffset3(iTransUnitIdx, iBlkSize, this.m_iCWidth):]
}
func (this *TComYuv) GetCrAddr2(iTransUnitIdx, iBlkSize uint) []Pel {
    return this.m_apiBufV[this.getAddrOffset3(iTransUnitIdx, iBlkSize, this.m_iCWidth):]
}

//  Get stride value of YUV buffer
func (this *TComYuv) GetStride() uint {
    return this.m_iWidth
}
func (this *TComYuv) GetCStride() uint {
    return this.m_iCWidth
}
func (this *TComYuv) GetHeight() uint {
    return this.m_iHeight
}
func (this *TComYuv) GetWidth() uint {
    return this.m_iWidth
}
func (this *TComYuv) GetCHeight() uint {
    return this.m_iCHeight
}
func (this *TComYuv) GetCWidth() uint {
    return this.m_iCWidth
}

//};// END CLASS DEFINITION TComYuv
