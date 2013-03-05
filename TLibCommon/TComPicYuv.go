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
    "errors"
    "os"
)

/// picture YUV buffer class
type TComPicYuv struct {
    // ------------------------------------------------------------------------------------------------
    //  YUV buffer
    // ------------------------------------------------------------------------------------------------
    m_apiPicBufY []Pel ///< Buffer (including margin)
    m_apiPicBufU []Pel
    m_apiPicBufV []Pel

    m_piPicOrgY []Pel ///< m_apiPicBufY + m_iMarginLuma*getStride() + m_iMarginLuma
    m_piPicOrgU []Pel
    m_piPicOrgV []Pel

    // ------------------------------------------------------------------------------------------------
    //  Parameter for general YUV buffer usage
    // ------------------------------------------------------------------------------------------------

    m_iPicWidth  int ///< Width of picture
    m_iPicHeight int ///< Height of picture

    m_iCuWidth  int ///< Width of Coding Unit (CU)
    m_iCuHeight int ///< Height of Coding Unit (CU)
    m_cuOffsetY []int
    m_cuOffsetC []int
    m_buOffsetY []int
    m_buOffsetC []int

    m_iLumaMarginX   int
    m_iLumaMarginY   int
    m_iChromaMarginX int
    m_iChromaMarginY int

    m_bIsBorderExtended bool
}

//public:
func NewTComPicYuv() *TComPicYuv {
    return &TComPicYuv{}
}

// ------------------------------------------------------------------------------------------------
//  Memory management
// ------------------------------------------------------------------------------------------------
func (this *TComPicYuv) Create(iPicWidth, iPicHeight int, uiMaxCUWidth, uiMaxCUHeight, uiMaxCUDepth uint) {
    this.m_iPicWidth = iPicWidth
    this.m_iPicHeight = iPicHeight

    // --> After config finished!
    this.m_iCuWidth = int(uiMaxCUWidth)
    this.m_iCuHeight = int(uiMaxCUHeight)

    numCuInWidth := this.m_iPicWidth / this.m_iCuWidth
    if (this.m_iPicWidth % this.m_iCuWidth) != 0 {
        numCuInWidth += 1
    }
    numCuInHeight := this.m_iPicHeight / this.m_iCuHeight
    if (this.m_iPicHeight % this.m_iCuHeight) != 0 {
        numCuInHeight += 1
    }

    this.m_iLumaMarginX = int(uiMaxCUWidth) + 16  // for 16-byte alignment
    this.m_iLumaMarginY = int(uiMaxCUHeight) + 16 // margin for 8-tap filter and infinite padding

    this.m_iChromaMarginX = this.m_iLumaMarginX >> 1
    this.m_iChromaMarginY = this.m_iLumaMarginY >> 1

    this.m_apiPicBufY = make([]Pel, (this.m_iPicWidth+(this.m_iLumaMarginX<<1))*(this.m_iPicHeight+(this.m_iLumaMarginY<<1)))
    this.m_apiPicBufU = make([]Pel, ((this.m_iPicWidth>>1)+(this.m_iChromaMarginX<<1))*((this.m_iPicHeight>>1)+(this.m_iChromaMarginY<<1)))
    this.m_apiPicBufV = make([]Pel, ((this.m_iPicWidth>>1)+(this.m_iChromaMarginX<<1))*((this.m_iPicHeight>>1)+(this.m_iChromaMarginY<<1)))

    this.m_piPicOrgY = this.m_apiPicBufY[this.m_iLumaMarginY*this.GetStride()+this.m_iLumaMarginX:]
    this.m_piPicOrgU = this.m_apiPicBufU[this.m_iChromaMarginY*this.GetCStride()+this.m_iChromaMarginX:]
    this.m_piPicOrgV = this.m_apiPicBufV[this.m_iChromaMarginY*this.GetCStride()+this.m_iChromaMarginX:]

    this.m_bIsBorderExtended = false

    this.m_cuOffsetY = make([]int, numCuInWidth*numCuInHeight)
    this.m_cuOffsetC = make([]int, numCuInWidth*numCuInHeight)
    for cuRow := 0; cuRow < numCuInHeight; cuRow++ {
        for cuCol := 0; cuCol < numCuInWidth; cuCol++ {
            this.m_cuOffsetY[cuRow*numCuInWidth+cuCol] = this.GetStride()*cuRow*this.m_iCuHeight + cuCol*this.m_iCuWidth
            this.m_cuOffsetC[cuRow*numCuInWidth+cuCol] = this.GetCStride()*cuRow*(this.m_iCuHeight/2) + cuCol*(this.m_iCuWidth/2)
        }
    }

    this.m_buOffsetY = make([]int, 1<<(2*uiMaxCUDepth))
    this.m_buOffsetC = make([]int, 1<<(2*uiMaxCUDepth))
    for buRow := 0; buRow < (1 << uiMaxCUDepth); buRow++ {
        for buCol := 0; buCol < (1 << uiMaxCUDepth); buCol++ {
            this.m_buOffsetY[(buRow<<uiMaxCUDepth)+buCol] = this.GetStride()*buRow*int(uiMaxCUHeight>>uiMaxCUDepth) + buCol*int(uiMaxCUWidth>>uiMaxCUDepth)
            this.m_buOffsetC[(buRow<<uiMaxCUDepth)+buCol] = this.GetCStride()*buRow*int((uiMaxCUHeight/2)>>uiMaxCUDepth) + buCol*int((uiMaxCUWidth/2)>>uiMaxCUDepth)
        }
    }
    return
}

func (this *TComPicYuv) Destroy() {
    //do nothing due to Garbage Collection of GO
}

func (this *TComPicYuv) GetCuOffsetY() []int{
	return this.m_cuOffsetY;
}
func (this *TComPicYuv) GetBuOffsetY() []int{
    return this.m_buOffsetY;
}

func (this *TComPicYuv) CreateLuma(iPicWidth, iPicHeight int, uiMaxCUWidth, uiMaxCUHeight, uiMaxCUDepth uint) {
    this.m_iPicWidth = iPicWidth
    this.m_iPicHeight = iPicHeight

    // --> After config finished!
    this.m_iCuWidth = int(uiMaxCUWidth)
    this.m_iCuHeight = int(uiMaxCUHeight)

    numCuInWidth := this.m_iPicWidth / this.m_iCuWidth
    if (this.m_iPicWidth % this.m_iCuWidth) != 0 {
        numCuInWidth += 1
    }
    numCuInHeight := this.m_iPicHeight / this.m_iCuHeight
    if (this.m_iPicHeight % this.m_iCuHeight) != 0 {
        numCuInHeight += 1
    }

    this.m_iLumaMarginX = int(uiMaxCUWidth) + 16  // for 16-byte alignment
    this.m_iLumaMarginY = int(uiMaxCUHeight) + 16 // margin for 8-tap filter and infinite padding

    this.m_apiPicBufY = make([]Pel, (this.m_iPicWidth+(this.m_iLumaMarginX<<1))*(this.m_iPicHeight+(this.m_iLumaMarginY<<1)))
    this.m_piPicOrgY = this.m_apiPicBufY[this.m_iLumaMarginY*this.GetStride()+this.m_iLumaMarginX:]

    this.m_cuOffsetY = make([]int, numCuInWidth*numCuInHeight)
    for cuRow := 0; cuRow < numCuInHeight; cuRow++ {
        for cuCol := 0; cuCol < numCuInWidth; cuCol++ {
            this.m_cuOffsetY[cuRow*numCuInWidth+cuCol] = this.GetStride()*cuRow*this.m_iCuHeight + cuCol*this.m_iCuWidth
        }
    }

    this.m_buOffsetY = make([]int, 1<<(2*uiMaxCUDepth))
    for buRow := 0; buRow < (1 << uiMaxCUDepth); buRow++ {
        for buCol := 0; buCol < (1 << uiMaxCUDepth); buCol++ {
            this.m_buOffsetY[(buRow<<uiMaxCUDepth)+buCol] = this.GetStride()*buRow*int(uiMaxCUHeight>>uiMaxCUDepth) + buCol*int(uiMaxCUWidth>>uiMaxCUDepth)
        }
    }
    return
}

func (this *TComPicYuv) DestroyLuma() {
    //do nothing
}

// ------------------------------------------------------------------------------------------------
//  Get information of picture
// ------------------------------------------------------------------------------------------------

func (this *TComPicYuv) GetWidth() int {
    return this.m_iPicWidth
}

func (this *TComPicYuv) GetHeight() int {
    return this.m_iPicHeight
}

func (this *TComPicYuv) GetStride() int {
    return (this.m_iPicWidth) + (this.m_iLumaMarginX << 1)
}

func (this *TComPicYuv) GetCStride() int {
    return (this.m_iPicWidth >> 1) + (this.m_iChromaMarginX << 1)
}

func (this *TComPicYuv) GetLumaMarginX() int {
    return this.m_iLumaMarginX
}

func (this *TComPicYuv) GetChromaMarginX() int {
    return this.m_iChromaMarginX
}

func (this *TComPicYuv) GetLumaMarginY() int {
    return this.m_iLumaMarginY
}

func (this *TComPicYuv) GetChromaMarginY() int {
    return this.m_iChromaMarginY
}

// ------------------------------------------------------------------------------------------------
//  Access function for picture buffer
// ------------------------------------------------------------------------------------------------

//  Access starting position of picture buffer with margin
func (this *TComPicYuv) GetBufY() []Pel {
    return this.m_apiPicBufY
}

func (this *TComPicYuv) GetBufU() []Pel {
    return this.m_apiPicBufU
}

func (this *TComPicYuv) GetBufV() []Pel {
    return this.m_apiPicBufV
}

//  Access starting position of original picture
func (this *TComPicYuv) GetLumaAddr() []Pel {
    return this.m_piPicOrgY
}

func (this *TComPicYuv) GetCbAddr() []Pel {
    return this.m_piPicOrgU
}

func (this *TComPicYuv) GetCrAddr() []Pel {
    return this.m_piPicOrgV
}

//  Access starting position of original picture for specific coding unit (CU) or partition unit (PU)

func (this *TComPicYuv) GetLumaAddr1(iCuAddr int) []Pel {
    return this.m_piPicOrgY[this.m_cuOffsetY[iCuAddr]:]
}

func (this *TComPicYuv) GetCbAddr1(iCuAddr int) []Pel {
    return this.m_piPicOrgU[this.m_cuOffsetC[iCuAddr]:]
}

func (this *TComPicYuv) GetCrAddr1(iCuAddr int) []Pel {
    return this.m_piPicOrgV[this.m_cuOffsetC[iCuAddr]:]
}

func (this *TComPicYuv) GetLumaAddr2(iCuAddr, uiAbsZorderIdx int) []Pel {
    return this.m_piPicOrgY[this.m_cuOffsetY[iCuAddr]+this.m_buOffsetY[G_auiZscanToRaster[uiAbsZorderIdx]]:]
}

func (this *TComPicYuv) GetCbAddr2(iCuAddr, uiAbsZorderIdx int) []Pel {
    return this.m_piPicOrgU[this.m_cuOffsetC[iCuAddr]+this.m_buOffsetC[G_auiZscanToRaster[uiAbsZorderIdx]]:]
}

func (this *TComPicYuv) GetCrAddr2(iCuAddr, uiAbsZorderIdx int) []Pel {
    return this.m_piPicOrgV[this.m_cuOffsetC[iCuAddr]+this.m_buOffsetC[G_auiZscanToRaster[uiAbsZorderIdx]]:]
}

// ------------------------------------------------------------------------------------------------
//  Miscellaneous
// ------------------------------------------------------------------------------------------------

//  Copy function to picture
func (this *TComPicYuv) CopyToPic(pcPicYuvDst *TComPicYuv) (err error) {
    if this.m_iPicWidth != pcPicYuvDst.GetWidth() ||
        this.m_iPicHeight != pcPicYuvDst.GetHeight() {
        err = errors.New("this.m_iPicWidth  != pcPicYuvDst.GetWidth() || this.m_iPicHeight != pcPicYuvDst.GetHeight()")
        return err
    }

    this.CopyToPicLuma(pcPicYuvDst)
    this.CopyToPicCb(pcPicYuvDst)
    this.CopyToPicCr(pcPicYuvDst)

    return nil
}

func (this *TComPicYuv) CopyToPicLuma(pcPicYuvDst *TComPicYuv) (err error) {
    if this.m_iPicWidth != pcPicYuvDst.GetWidth() ||
        this.m_iPicHeight != pcPicYuvDst.GetHeight() {
        err = errors.New("this.m_iPicWidth  != pcPicYuvDst.GetWidth() || this.m_iPicHeight != pcPicYuvDst.GetHeight()")
        return err
    }

    pcPicYuvDstY := pcPicYuvDst.GetBufY()
    for k := 0; k < (this.m_iPicHeight+(this.m_iLumaMarginY<<1))*(this.m_iPicWidth+(this.m_iLumaMarginX<<1)); k++ {
        pcPicYuvDstY[k] = this.m_apiPicBufY[k]
    }

    return nil
}

func (this *TComPicYuv) CopyToPicCb(pcPicYuvDst *TComPicYuv) (err error) {
    if this.m_iPicWidth != pcPicYuvDst.GetWidth() ||
        this.m_iPicHeight != pcPicYuvDst.GetHeight() {
        err = errors.New("this.m_iPicWidth  != pcPicYuvDst.GetWidth() || this.m_iPicHeight != pcPicYuvDst.GetHeight()")
        return err
    }

    pcPicYuvDstU := pcPicYuvDst.GetBufU()
    for k := 0; k < ((this.m_iPicWidth>>1)+(this.m_iChromaMarginX<<1))*((this.m_iPicHeight>>1)+(this.m_iChromaMarginY<<1)); k++ {
        pcPicYuvDstU[k] = this.m_apiPicBufU[k]
    }

    return nil
}

func (this *TComPicYuv) CopyToPicCr(pcPicYuvDst *TComPicYuv) (err error) {
    if this.m_iPicWidth != pcPicYuvDst.GetWidth() ||
        this.m_iPicHeight != pcPicYuvDst.GetHeight() {
        err = errors.New("this.m_iPicWidth  != pcPicYuvDst.GetWidth() || this.m_iPicHeight != pcPicYuvDst.GetHeight()")
        return err
    }

    pcPicYuvDstV := pcPicYuvDst.GetBufV()
    for k := 0; k < ((this.m_iPicWidth>>1)+(this.m_iChromaMarginX<<1))*((this.m_iPicHeight>>1)+(this.m_iChromaMarginY<<1)); k++ {
        pcPicYuvDstV[k] = this.m_apiPicBufV[k]
    }

    return nil
}

// Set border extension flag
func (this *TComPicYuv) SetBorderExtension(bIsBorderExtended bool) {
    this.m_bIsBorderExtended = bIsBorderExtended
}

//  Extend function of picture buffer
func (this *TComPicYuv) ExtendPicBorder() {
    if this.m_bIsBorderExtended {
        return
    }

    this.xExtendPicCompBorder(this.GetBufY(), this.GetLumaAddr(), this.GetStride(), this.GetWidth(), this.GetHeight(), this.m_iLumaMarginX, this.m_iLumaMarginY)
    this.xExtendPicCompBorder(this.GetBufU(), this.GetCbAddr(), this.GetCStride(), this.GetWidth()>>1, this.GetHeight()>>1, this.m_iChromaMarginX, this.m_iChromaMarginY)
    this.xExtendPicCompBorder(this.GetBufV(), this.GetCrAddr(), this.GetCStride(), this.GetWidth()>>1, this.GetHeight()>>1, this.m_iChromaMarginX, this.m_iChromaMarginY)

    this.m_bIsBorderExtended = true
}

//  Dump picture
func (this *TComPicYuv) Dump(pFileName string, bAdd bool) (err error) {
    var pFile *os.File

    if !bAdd {
        pFile, err = os.Create(pFileName)
    } else {
        pFile, err = os.OpenFile(pFileName, os.O_APPEND, os.ModeAppend)
    }
    if err != nil {
        return err
    }
    defer pFile.Close()

    var offset Pel

    shift := uint(G_bitDepthY - 8)
    if shift > 0 {
        offset = 1 << (shift - 1)
    } else {
        offset = 0
    }

    var x, y int

    piY := this.GetLumaAddr()
    piCb := this.GetCbAddr()
    piCr := this.GetCrAddr()
    iStride := this.GetStride()
    iStrideC := this.GetCStride()
    uy := make([]byte, this.m_iPicWidth)
    uc := make([]byte, this.m_iPicWidth>>1)

    for y = 0; y < this.m_iPicHeight; y++ {
        for x = 0; x < this.m_iPicWidth; x++ {
            uy[x] = byte(CLIP3(0, 255, (piY[y*iStride+x]+offset)>>shift).(Pel))
        }
        pFile.Write(uy)
    }

    shift = uint(G_bitDepthC - 8)
    if shift > 0 {
        offset = 1 << (shift - 1)
    } else {
        offset = 0
    }

    for y = 0; y < this.m_iPicHeight>>1; y++ {
        for x = 0; x < this.m_iPicWidth>>1; x++ {
            uc[x] = byte(CLIP3(0, 255, (piCb[y*iStrideC+x]+offset)>>shift).(Pel))
        }
        pFile.Write(uc)
    }

    for y = 0; y < this.m_iPicHeight>>1; y++ {
        for x = 0; x < this.m_iPicWidth>>1; x++ {
            uc[x] = byte(CLIP3(0, 255, (piCr[y*iStrideC+x]+offset)>>shift).(Pel))
        }
        pFile.Write(uc)
    }

    return nil
}

//protected:
func (this *TComPicYuv) xExtendPicCompBorder(pi []Pel, piTxt []Pel, iStride, iWidth, iHeight, iMarginX, iMarginY int) {
    var x, y int

    for y = 0; y < iHeight; y++ {
        for x = 0; x < iMarginX; x++ {
            pi[(y+iMarginY)*iStride-iMarginX+x+iMarginX] = piTxt[y*iStride+0]
            pi[(y+iMarginY)*iStride+iWidth+x+iMarginX] = piTxt[y*iStride+iWidth-1]
        }
    }

    for y = 0; y < iMarginY; y++ {
        for x = 0; x < iWidth+(iMarginX<<1); x++ {
            pi[y*iStride+x] = pi[iMarginY*iStride+x]
            pi[(y+iHeight+iMarginY)*iStride+x] = pi[(iHeight-1+iMarginY)*iStride+x]
        }
    }
}
