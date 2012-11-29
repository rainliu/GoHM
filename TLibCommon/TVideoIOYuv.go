package TLibCommon

import (
	"os"
	"log"
)

/// YUV file I/O class
type TVideoIOYuv struct{
//private:
  m_cHandle			*os.File;                                      ///< file handle
  m_fileBitDepthY	int; ///< bitdepth of input/output video file luma component
  m_fileBitDepthC 	int; ///< bitdepth of input/output video file chroma component
  m_bitDepthShiftY 	int;  ///< number of bits to increase or decrease luma by before/after write/read
  m_bitDepthShiftC 	int;  ///< number of bits to increase or decrease chroma by before/after write/read
}

func NewTVideoIOYuv() (*TVideoIOYuv){
	return &TVideoIOYuv{};
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
func (this *TVideoIOYuv) Open(pchFile string, bWriteMode bool, fileBitDepthY, fileBitDepthC, internalBitDepthY, internalBitDepthC int) (err error){ 
  this.m_bitDepthShiftY = internalBitDepthY - fileBitDepthY;
  this.m_bitDepthShiftC = internalBitDepthC - fileBitDepthC;
  this.m_fileBitDepthY = fileBitDepthY;
  this.m_fileBitDepthC = fileBitDepthC;

  if bWriteMode {
    if this.m_cHandle, err = os.Create( pchFile ); err!=nil{
		log.Fatal(err)
		return err
	}
  }else{
    if this.m_cHandle, err = os.Open( pchFile ); err!=nil{
		log.Fatal(err)
		return err
	}
  }
  
  return nil;
}
 
///< close file 
func (this *TVideoIOYuv) Close (){                                           
	this.m_cHandle.Close();
}
 

/**
 * Skip numFrames in input.
 *
 * This function correctly handles cases where the input file is not
 * seekable, by consuming bytes.
 */ 
func (this *TVideoIOYuv) SkipFrames(numFrames, width, height uint) (err error){
  if numFrames==0 {
    return nil;
  }
  
  var wordsize uint
  if (this.m_fileBitDepthY > 8 || this.m_fileBitDepthC > 8){
  	wordsize = 2;
  }else{
  	wordsize = 1;
  }
  
  framesize := wordsize * width * height * 3 / 2;
  offset := framesize * numFrames;

  /* attempt to seek */
  //if (!!m_cHandle.seekg(offset, ios::cur))
  //  return; /* success */
  //m_cHandle.clear();
  _, err = this.m_cHandle.Seek(int64(offset), 1)
  return err;
  
  /* fall back to consuming the input 
  buf := make([]byte, 512)
  offset_mod_bufsize := offset % len(buf);
  for i := 0; i < offset - offset_mod_bufsize; i += len(buf) {
    this.m_cHandle.Read(buf);
  }
  this.m_cHandle.Read(buf, offset_mod_bufsize);*/
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
func (this *TVideoIOYuv) Read ( pPicYuv *TComPicYuv, aiPad []int ) bool {   
  // check end-of-file
  if this.IsEof() {
  	return false;
  }
  
  iStride := uint(pPicYuv.GetStride());
  
  // compute actual YUV width & height excluding padding size
  pad_h 		:= uint(aiPad[0]);
  pad_v 		:= uint(aiPad[1]);
  width_full 	:= uint(pPicYuv.GetWidth());
  height_full 	:= uint(pPicYuv.GetHeight());
  width  		:= uint(width_full - pad_h);
  height 		:= uint(height_full - pad_v);
  is16bit 		:= this.m_fileBitDepthY > 8 || this.m_fileBitDepthC > 8;

  desired_bitdepthY := uint(this.m_fileBitDepthY + this.m_bitDepthShiftY);
  desired_bitdepthC := uint(this.m_fileBitDepthC + this.m_bitDepthShiftC);
  minvalY := Pel(0);
  minvalC := Pel(0);
  maxvalY := Pel((1 << desired_bitdepthY) - 1);
  maxvalC := Pel((1 << desired_bitdepthC) - 1);

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
  if !this.readPlane(pPicYuv.GetLumaAddr(), this.m_cHandle, is16bit, iStride, width, height, pad_h, pad_v){
    return false;   
  }
  this.scalePlane(pPicYuv.GetLumaAddr(), iStride, width_full, height_full, this.m_bitDepthShiftY, minvalY, maxvalY);

  iStride >>= 1;
  width_full >>= 1;
  height_full >>= 1;
  width >>= 1;
  height >>= 1;
  pad_h >>= 1;
  pad_v >>= 1;

  if !this.readPlane(pPicYuv.GetCbAddr(), this.m_cHandle, is16bit, iStride, width, height, pad_h, pad_v){
    return false;
  }
  this.scalePlane(pPicYuv.GetCbAddr(), iStride, width_full, height_full, this.m_bitDepthShiftC, minvalC, maxvalC);

  if !this.readPlane(pPicYuv.GetCrAddr(), this.m_cHandle, is16bit, iStride, width, height, pad_h, pad_v){
    return false;
  }
  this.scalePlane(pPicYuv.GetCrAddr(), iStride, width_full, height_full, this.m_bitDepthShiftC, minvalC, maxvalC);

  return true;
}

func (this *TVideoIOYuv) Write( pPicYuv *TComPicYuv, cropLeft, cropRight, cropTop, cropBottom int ) bool {

	return true
}
  
///< check for end-of-file  
func (this *TVideoIOYuv) IsEof () bool {                                           

	return true
}

///< check for failure
func (this *TVideoIOYuv) IsFail() bool {                                           

	return true
}

func (this *TVideoIOYuv) readPlane(dst *Pel, fd *os.File, is16bit bool, stride, width, height, pad_x, pad_y uint) bool{
	return true
}

func (this *TVideoIOYuv) scalePlane(img *Pel, stride, width, height uint, shiftbits int, minval, maxval Pel){
}