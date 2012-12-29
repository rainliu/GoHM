package TLibCommon

import (
	//"fmt"
)

const NTAPS_LUMA = 8                                   ///< Number of taps for luma
const NTAPS_CHROMA = 4                                 ///< Number of taps for chroma
const IF_INTERNAL_PREC = 14                            ///< Number of bits for internal precision
const IF_FILTER_PREC = 6                               ///< Log2 of sum of filter taps
const IF_INTERNAL_OFFS = (1 << (IF_INTERNAL_PREC - 1)) ///< Offset used internally


var m_lumaFilter	=[4][NTAPS_LUMA]Pel{
  {  0, 0,   0, 64,  0,   0, 0,  0 },
  { -1, 4, -10, 58, 17,  -5, 1,  0 },
  { -1, 4, -11, 40, 40, -11, 4, -1 },
  {  0, 1,  -5, 17, 58, -10, 4, -1 },
};

var m_chromaFilter =[8][NTAPS_CHROMA]Pel{
  {  0, 64,  0,  0 },
  { -2, 58, 10, -2 },
  { -4, 54, 16, -2 },
  { -6, 46, 28, -4 },
  { -4, 36, 36, -4 },
  { -4, 28, 46, -6 },
  { -2, 16, 54, -4 },
  { -2, 10, 58, -2 },
};
/**
 * \brief Interpolation filter class
 */
type TComInterpolationFilter struct {
    //m_lumaFilter   [4][NTAPS_LUMA]int16   ///< Luma filter taps
    //m_chromaFilter [8][NTAPS_CHROMA]int16 ///< Chroma filter taps
}

func NewTComInterpolationFilter() *TComInterpolationFilter{
	return &TComInterpolationFilter{};
}

func (this *TComInterpolationFilter) FilterHorLuma  (src []Pel, srcStride int, dst []Pel, dstStride, width, height, frac int,            isLast bool){
  //assert(frac >= 0 && frac < 4);
  
  if frac == 0 {
    this.filterCopy(G_bitDepthY, src, srcStride, dst, dstStride, width, height, true, isLast );
  }else{
    this.filterHor(NTAPS_LUMA, G_bitDepthY, src, srcStride, dst, dstStride, width, height, isLast, m_lumaFilter[frac][:]);
  }
}
func (this *TComInterpolationFilter) FilterVerLuma  (src []Pel, srcStride int, dst []Pel, dstStride, width, height, frac int,  isFirst,  isLast bool){
  //assert(frac >= 0 && frac < 4);
  
  if frac == 0 {
    this.filterCopy(G_bitDepthY, src, srcStride, dst, dstStride, width, height, isFirst, isLast );
  }else{
    this.filterVer(NTAPS_LUMA, G_bitDepthY, src, srcStride, dst, dstStride, width, height, isFirst, isLast, m_lumaFilter[frac][:]);
  }
}
func (this *TComInterpolationFilter) FilterHorChroma(src []Pel, srcStride int, dst []Pel, dstStride, width, height, frac int,            isLast bool){
  //assert(frac >= 0 && frac < 8);
  
  if frac == 0 {
    this.filterCopy(G_bitDepthC, src, srcStride, dst, dstStride, width, height, true, isLast );
  }else{
    this.filterHor(NTAPS_CHROMA, G_bitDepthC, src, srcStride, dst, dstStride, width, height, isLast, m_chromaFilter[frac][:]);
  }
}
func (this *TComInterpolationFilter) FilterVerChroma(src []Pel, srcStride int, dst []Pel, dstStride, width, height, frac int,  isFirst,  isLast bool){
  //assert(frac >= 0 && frac < 8);
  
  if frac == 0 {
    this.filterCopy(G_bitDepthC, src, srcStride, dst, dstStride, width, height, isFirst, isLast );
  }else{
    this.filterVer(NTAPS_CHROMA, G_bitDepthC, src, srcStride, dst, dstStride, width, height, isFirst, isLast, m_chromaFilter[frac][:]);
  }
}

func (this *TComInterpolationFilter) filterCopy( bitDepth int, src []Pel, srcStride int, dst []Pel, dstStride, width, height int, isFirst, isLast bool){
  var row, col int; 
  
  if isFirst == isLast {
    for row = 0; row < height; row++ {
      for col = 0; col < width; col++ {
        dst[row*dstStride+col] = src[row*srcStride+col];
      }
      
      //src += srcStride;
      //dst += dstStride;
    }              
  }else if isFirst {
    shift := uint(IF_INTERNAL_PREC - bitDepth);
    
    for row = 0; row < height; row++ {
      for col = 0; col < width; col++ {
        val := src[row*srcStride+col] << shift;
        dst[row*dstStride+col] = val - Pel(IF_INTERNAL_OFFS);
      }
      
      //src += srcStride;
      //dst += dstStride;
    }          
  }else{
    shift := uint(IF_INTERNAL_PREC - bitDepth);
    offset := (IF_INTERNAL_OFFS);
    if shift!=0 {
    	offset += (1 << (shift - 1));
    }else{
    	offset += 0;
    }
    maxVal := (1 << uint(bitDepth)) - 1;
    minVal := 0;
    for row = 0; row < height; row++ {
      for col = 0; col < width; col++ {
        val := int(src[ row*srcStride+col ]);
        val = ( val + offset ) >> shift;
        if val < minVal {
         val = minVal;
        }
        if val > maxVal {
         val = maxVal;
        }
        dst[row*dstStride+col] = Pel(val);
      }
      
      //src += srcStride;
      //dst += dstStride;
    }              
  }
}

func (this *TComInterpolationFilter) filter(N int, isVertical, isFirst, isLast bool, bitDepth int, src2 []Pel, srcStride int, dst []Pel, dstStride, width, height int, coeff []Pel){
  var row, col int;
  
  var c	[8]Pel;
  c[0] = coeff[0];
  c[1] = coeff[1];
  if N >= 4 {
    c[2] = coeff[2];
    c[3] = coeff[3];
  }
  if N >= 6 {
    c[4] = coeff[4];
    c[5] = coeff[5];
  }
  if N == 8 {
    c[6] = coeff[6];
    c[7] = coeff[7];
  }
  
  var cStride int;
  if isVertical {
   	cStride = srcStride;
  }else{
  	cStride = 1;
  }

  src := src2[-( N/2 - 1 ) * cStride:];
  //src := src2[( N/2 - 1 ) * cStride-( N/2 - 1 ) * cStride:];
  //fmt.Printf("How to handle src -= ( N/2 - 1 ) * cStride?\n");
  
  var offset int;
  var maxVal int;
  headRoom := IF_INTERNAL_PREC - bitDepth;
  shift := IF_FILTER_PREC;
  if isLast {
    offset = 1 << uint(shift - 1);
    
    if isFirst {
    	shift += 0;
    	offset += 0;
   	}else{
    	shift += headRoom;
    	offset += IF_INTERNAL_OFFS << IF_FILTER_PREC;
   	}
    maxVal = (1 << uint(bitDepth)) - 1;
  }else{
  	if isFirst {
    	shift -= headRoom;
    	offset = -IF_INTERNAL_OFFS << uint(shift);
    }else{
    	shift -=  0;
    	offset =  0;
    }
    maxVal = 0;
  }
  
  for row = 0; row < height; row++ {
    for col = 0; col < width; col++ {
      var sum int;
      
      sum  = int(src[ row*srcStride+col + 0 * cStride] * c[0]);
      sum += int(src[ row*srcStride+col + 1 * cStride] * c[1]);
      if N >= 4 {
        sum += int(src[ row*srcStride+col + 2 * cStride] * c[2]);
        sum += int(src[ row*srcStride+col + 3 * cStride] * c[3]);
      }
      if N >= 6 {
        sum += int(src[ row*srcStride+col + 4 * cStride] * c[4]);
        sum += int(src[ row*srcStride+col + 5 * cStride] * c[5]);
      }
      if N == 8 {
        sum += int(src[ row*srcStride+col + 6 * cStride] * c[6]);
        sum += int(src[ row*srcStride+col + 7 * cStride] * c[7]);        
      }
      
      val := ( sum + offset ) >> uint(shift);
      if isLast {
        if val < 0 {
        	val = 0;
        }else if val > maxVal { 
        	val = maxVal;
        }        
      }
      dst[row*dstStride+col] = Pel(val);
    }
    
    //src += srcStride;
    //dst += dstStride;
  }    
}

func (this *TComInterpolationFilter) filterHor(N int, bitDepth int, src []Pel, srcStride int, dst []Pel, dstStride, width, height int,          isLast bool, coeff []Pel){
  if isLast {
    this.filter(N, false, true, true,  bitDepth, src, srcStride, dst, dstStride, width, height, coeff);
  }else{
    this.filter(N, false, true, false, bitDepth, src, srcStride, dst, dstStride, width, height, coeff);
  }
}

func (this *TComInterpolationFilter) filterVer(N int, bitDepth int, src []Pel, srcStride int, dst []Pel, dstStride, width, height int, isFirst, isLast bool, coeff []Pel){
  if isFirst && isLast {
    this.filter(N, true, true, true,   bitDepth, src, srcStride, dst, dstStride, width, height, coeff);
  }else if isFirst && !isLast {
    this.filter(N, true, true, false,  bitDepth, src, srcStride, dst, dstStride, width, height, coeff);
  }else if !isFirst && isLast {
    this.filter(N, true, false, true,  bitDepth, src, srcStride, dst, dstStride, width, height, coeff);
  }else{
    this.filter(N, true, false, false, bitDepth, src, srcStride, dst, dstStride, width, height, coeff);
  }    
}