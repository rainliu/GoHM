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
    "io"
    "log"
    "os"
)

/// YUV file I/O class
type TVideoIOYuv struct {
    //private:
    m_cHandle        *os.File ///< file handle
    m_fileBitDepthY  int      ///< bitdepth of input/output video file luma component
    m_fileBitDepthC  int      ///< bitdepth of input/output video file chroma component
    m_bitDepthShiftY int      ///< number of bits to increase or decrease luma by before/after write/read
    m_bitDepthShiftC int      ///< number of bits to increase or decrease chroma by before/after write/read
    m_eof            bool
}

func NewTVideoIOYuv() *TVideoIOYuv {
    return &TVideoIOYuv{}
}

///< open or create file
/**
 * Open file for reading/writing Y'CbCr frames.
 *
 * Frames read/written have bitdepth fileBitDepth, and are automatically
 * formatted as 8 or 16 bit word values (see TVideoIOYuv::write()).
 *
 * Image data read or written is converted to/from internalBitDepth
 * (See scalePlane(), TVideoIOYuv::read() and TVideoIOYuv::write() for
 * further details).
 *
 * \param pchFile          file name string
 * \param bWriteMode       file open mode: true=read, false=write
 * \param fileBitDepthY     bit-depth of input/output file data (luma component).
 * \param fileBitDepthC     bit-depth of input/output file data (chroma components).
 * \param internalBitDepthY bit-depth to scale image data to/from when reading/writing (luma component).
 * \param internalBitDepthC bit-depth to scale image data to/from when reading/writing (chroma components).
 */
func (this *TVideoIOYuv) Open(pchFile string, bWriteMode bool, fileBitDepthY, fileBitDepthC, internalBitDepthY, internalBitDepthC int) (err error) {
    this.m_bitDepthShiftY = internalBitDepthY - fileBitDepthY
    this.m_bitDepthShiftC = internalBitDepthC - fileBitDepthC
    this.m_fileBitDepthY = fileBitDepthY
    this.m_fileBitDepthC = fileBitDepthC

    if bWriteMode {
        if this.m_cHandle, err = os.Create(pchFile); err != nil {
            log.Fatal(err)
            return err
        }
    } else {
        if this.m_cHandle, err = os.Open(pchFile); err != nil {
            log.Fatal(err)
            return err
        }
    }

    return nil
}

///< close file
func (this *TVideoIOYuv) Close() {
    if this.m_cHandle != nil {
        this.m_cHandle.Close()
    }
}

/**
 * Skip numFrames in input.
 *
 * This function correctly handles cases where the input file is not
 * seekable, by consuming bytes.
 */
func (this *TVideoIOYuv) SkipFrames(numFrames, width, height uint) (err error) {
    if numFrames == 0 {
        return nil
    }

    var wordsize uint
    if this.m_fileBitDepthY > 8 || this.m_fileBitDepthC > 8 {
        wordsize = 2
    } else {
        wordsize = 1
    }

    framesize := wordsize * width * height * 3 / 2
    offset := framesize * numFrames

    /* attempt to seek */
    //if (!!m_cHandle.seekg(offset, ios::cur))
    //  return; /* success */
    //m_cHandle.clear();
    _, err = this.m_cHandle.Seek(int64(offset), 1)
    return err

    /* fall back to consuming the input
       buf := make([]byte, 512)
       offset_mod_bufsize := offset % len(buf);
       for i := 0; i < offset - offset_mod_bufsize; i += len(buf) {
         this.m_cHandle.Read(buf);
       }
       this.m_cHandle.Read(buf, offset_mod_bufsize);*/
}

func (this *TVideoIOYuv) IsEof() bool {
    return this.m_eof
}

///< read  one YUV frame with padding parameter
/**
 * Read one Y'CbCr frame, performing any required input scaling to change
 * from the bitdepth of the input file to the internal bit-depth.
 *
 * If a bit-depth reduction is required, and internalBitdepth >= 8, then
 * the input file is assumed to be ITU-R BT.601/709 compliant, and the
 * resulting data is clipped to the appropriate legal range, as if the
 * file had been provided at the lower-bitdepth compliant to Rec601/709.
 *
 * @param pPicYuv      input picture YUV buffer class pointer
 * @param aiPad        source padding size, aiPad[0] = horizontal, aiPad[1] = vertical
 * @return true for success, false in case of error
 */
func (this *TVideoIOYuv) Read(pPicYuv *TComPicYuv, aiPad []int) bool {
    // check end-of-file
    if this.IsEof() {
        return false
    }

    iStride := pPicYuv.GetStride()

    // compute actual YUV width & height excluding padding size
    pad_h := aiPad[0]
    pad_v := aiPad[1]
    width_full := pPicYuv.GetWidth()
    height_full := pPicYuv.GetHeight()
    width := width_full - pad_h
    height := height_full - pad_v
    is16bit := this.m_fileBitDepthY > 8 || this.m_fileBitDepthC > 8

    desired_bitdepthY := uint(this.m_fileBitDepthY + this.m_bitDepthShiftY)
    desired_bitdepthC := uint(this.m_fileBitDepthC + this.m_bitDepthShiftC)
    minvalY := Pel(0)
    minvalC := Pel(0)
    maxvalY := Pel((1 << desired_bitdepthY) - 1)
    maxvalC := Pel((1 << desired_bitdepthC) - 1)

    /*
       #if CLIP_TO_709_RANGE
         if (m_bitdepthShiftY < 0 && desired_bitdepthY >= 8)
         {
           // ITU-R BT.709 compliant clipping for converting say 10b to 8b
           minvalY = 1 << (desired_bitdepthY - 8);
           maxvalY = (0xff << (desired_bitdepthY - 8)) -1;
         }
         if (m_bitdepthShiftC < 0 && desired_bitdepthC >= 8)
         {
           // ITU-R BT.709 compliant clipping for converting say 10b to 8b
           minvalC = 1 << (desired_bitdepthC - 8);
           maxvalC = (0xff << (desired_bitdepthC - 8)) -1;
         }
       #endif
    */
    if !this.readPlane(pPicYuv.GetLumaAddr(), this.m_cHandle, is16bit, iStride, width, height, pad_h, pad_v) {
        return false
    }
    this.scalePlane(pPicYuv.GetLumaAddr(), iStride, width_full, height_full, this.m_bitDepthShiftY, minvalY, maxvalY)

    iStride >>= 1
    width_full >>= 1
    height_full >>= 1
    width >>= 1
    height >>= 1
    pad_h >>= 1
    pad_v >>= 1

    if !this.readPlane(pPicYuv.GetCbAddr(), this.m_cHandle, is16bit, iStride, width, height, pad_h, pad_v) {
        return false
    }
    this.scalePlane(pPicYuv.GetCbAddr(), iStride, width_full, height_full, this.m_bitDepthShiftC, minvalC, maxvalC)

    if !this.readPlane(pPicYuv.GetCrAddr(), this.m_cHandle, is16bit, iStride, width, height, pad_h, pad_v) {
        return false
    }
    this.scalePlane(pPicYuv.GetCrAddr(), iStride, width_full, height_full, this.m_bitDepthShiftC, minvalC, maxvalC)

    return true
}

/**
 * Write one Y'CbCr frame. No bit-depth conversion is performed, pcPicYuv is
 * assumed to be at TVideoIO::m_fileBitdepth depth.
 *
 * @param pPicYuv     input picture YUV buffer class pointer
 * @param aiPad       source padding size, aiPad[0] = horizontal, aiPad[1] = vertical
 * @return true for success, false in case of error
 */
func (this *TVideoIOYuv) Write(pPicYuv *TComPicYuv, confLeft, confRight, confTop, confBottom int) bool {
    // compute actual YUV frame size excluding padding size
    iStride := pPicYuv.GetStride()
    width := pPicYuv.GetWidth() - confLeft - confRight
    height := pPicYuv.GetHeight() - confTop - confBottom
    is16bit := this.m_fileBitDepthY > 8 || this.m_fileBitDepthC > 8
    var dstPicYuv *TComPicYuv
    retval := true

    if this.m_bitDepthShiftY != 0 || this.m_bitDepthShiftC != 0 {
        dstPicYuv = NewTComPicYuv()
        dstPicYuv.Create(pPicYuv.GetWidth(), pPicYuv.GetHeight(), 1, 1, 0)
        pPicYuv.CopyToPic(dstPicYuv)

        minvalY := Pel(0)
        minvalC := Pel(0)
        maxvalY := Pel((1 << uint(this.m_fileBitDepthY)) - 1)
        maxvalC := Pel((1 << uint(this.m_fileBitDepthC)) - 1)

        /*#if CLIP_TO_709_RANGE
            if (-m_bitDepthShiftY < 0 && m_fileBitDepthY >= 8)
            {
              // ITU-R BT.709 compliant clipping for converting say 10b to 8b
              minvalY = 1 << (m_fileBitDepthY - 8);
              maxvalY = (0xff << (m_fileBitDepthY - 8)) -1;
            }
            if (-m_bitDepthShiftC < 0 && m_fileBitDepthC >= 8)
            {
              // ITU-R BT.709 compliant clipping for converting say 10b to 8b
              minvalC = 1 << (m_fileBitDepthC - 8);
              maxvalC = (0xff << (m_fileBitDepthC - 8)) -1;
            }
        #endif*/

        this.scalePlane(dstPicYuv.GetLumaAddr(), dstPicYuv.GetStride(), dstPicYuv.GetWidth(), dstPicYuv.GetHeight(), -this.m_bitDepthShiftY, minvalY, maxvalY)
        this.scalePlane(dstPicYuv.GetCbAddr(), dstPicYuv.GetCStride(), dstPicYuv.GetWidth()>>1, dstPicYuv.GetHeight()>>1, -this.m_bitDepthShiftC, minvalC, maxvalC)
        this.scalePlane(dstPicYuv.GetCrAddr(), dstPicYuv.GetCStride(), dstPicYuv.GetWidth()>>1, dstPicYuv.GetHeight()>>1, -this.m_bitDepthShiftC, minvalC, maxvalC)
    } else {
        dstPicYuv = pPicYuv
    }
    // location of upper left pel in a plane
    planeOffset := confLeft + confTop * iStride;

    if !this.writePlane(this.m_cHandle, dstPicYuv.GetLumaAddr()[planeOffset:], is16bit, iStride, width, height) {
        retval = false
        goto exit
    }

    width >>= 1
    height >>= 1
    iStride >>= 1
    confLeft >>= 1
    confRight >>= 1
    confTop >>= 1;
  	confBottom >>= 1;

    planeOffset = confLeft + confTop * iStride;

    if !this.writePlane(this.m_cHandle, dstPicYuv.GetCbAddr()[planeOffset:], is16bit, iStride, width, height) {
        retval = false
        goto exit
    }
    if !this.writePlane(this.m_cHandle, dstPicYuv.GetCrAddr()[planeOffset:], is16bit, iStride, width, height) {
        retval = false
        goto exit
    }

exit:
    if this.m_bitDepthShiftY != 0 || this.m_bitDepthShiftC != 0 {
        dstPicYuv.Destroy()
    }
    return retval
}

///< check for end-of-file
//func (this *TVideoIOYuv) IsEof () bool {
//	this.m_cHandle.
//	return true
//}

///< check for failure
//func (this *TVideoIOYuv) IsFail() bool {
//	return true
//}

func (this *TVideoIOYuv) readPlane(dst []Pel, fd *os.File, is16bit bool, stride, width, height, pad_x, pad_y int) bool {
    var read_len, x, y int

    if is16bit {
        read_len = width * 2
    } else {
        read_len = width
    }

    buf := make([]byte, read_len)
    for y = 0; y < height; y++ {
        n, err := fd.Read(buf)
        if err == io.EOF || n != read_len {
            this.m_eof = true
            return false
        }

        if !is16bit {
            for x = 0; x < width; x++ {
                dst[y*stride+x] = Pel(buf[x])
            }
        } else {
            for x = 0; x < width; x++ {
                dst[y*stride+x] = (Pel(buf[2*x+1]) << 8) | Pel(buf[2*x])
            }
        }

        for x = width; x < width+pad_x; x++ {
            dst[y*stride+x] = dst[y*stride+width-1]
        }
    }

    for y = height; y < height+pad_y; y++ {
        for x = 0; x < width+pad_x; x++ {
            dst[y*stride+x] = dst[(y-1)*stride+x]
        }
    }

    return true
}

func (this *TVideoIOYuv) writePlane(fd *os.File, src []Pel, is16bit bool, stride, width, height int) bool {
    var write_len, x, y int

    if is16bit {
        write_len = width * 2
    } else {
        write_len = width
    }

    buf := make([]byte, write_len)
    for y = 0; y < height; y++ {
        if !is16bit {
            for x = 0; x < width; x++ {
                buf[x] = byte(src[y*stride+x])
            }
        } else {
            for x = 0; x < width; x++ {
                buf[2*x] = byte(src[y*stride+x] & 0xff)
                buf[2*x+1] = byte((src[y*stride+x] >> 8) & 0xff)
            }
        }

        n, err := fd.Write(buf)
        if err != nil || n != write_len {
            return false
        }
    }

    return true
}

func (this *TVideoIOYuv) scalePlane(img []Pel, stride, width, height, shiftbits int, minval, maxval Pel) {
    if shiftbits == 0 {
        return
    }

    if shiftbits > 0 {
        this.fwdScalePlane(img, stride, width, height, shiftbits)
    } else {
        this.invScalePlane(img, stride, width, height, -shiftbits, minval, maxval)
    }
}

/**
 * Multiply all pixels in img by 2<sup>shiftbits</sup>.
 *
 * @param img        pointer to image to be transformed
 * @param stride     distance between vertically adjacent pixels of img.
 * @param width      width of active area in img.
 * @param height     height of active area in img.
 * @param shiftbits  number of bits to shift
 */
func (this *TVideoIOYuv) fwdScalePlane(img []Pel, stride, width, height, shiftbits int) {
    for y := 0; y < height; y++ {
        for x := 0; x < width; x++ {
            img[y*stride+x] <<= uint(shiftbits)
        }
    }
}

/**
 * Perform division with rounding of all pixels in img by
 * 2<sup>shiftbits</sup>. All pixels are clipped to [minval, maxval]
 *
 * @param img        pointer to image to be transformed
 * @param stride     distance between vertically adjacent pixels of img.
 * @param width      width of active area in img.
 * @param height     height of active area in img.
 * @param shiftbits  number of rounding bits
 * @param minval     minimum clipping value
 * @param maxval     maximum clipping value
 */
func (this *TVideoIOYuv) invScalePlane(img []Pel, stride, width, height, shiftbits int, minval, maxval Pel) {
    offset := Pel(1 << uint(shiftbits-1))

    for y := 0; y < height; y++ {
        for x := 0; x < width; x++ {
            val := (img[y*stride+x] + offset) >> uint(shiftbits)
            img[y*stride+x] = CLIP3(minval, maxval, val).(Pel)
        }
    }
}
