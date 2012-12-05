package TLibCommon

import (

)


const NTAPS_LUMA        = 8 ///< Number of taps for luma
const NTAPS_CHROMA      = 4 ///< Number of taps for chroma
const IF_INTERNAL_PREC  = 14 ///< Number of bits for internal precision
const IF_FILTER_PREC    = 6 ///< Log2 of sum of filter taps
const IF_INTERNAL_OFFS  = (1<<(IF_INTERNAL_PREC-1)) ///< Offset used internally


/**
 * \brief Interpolation filter class
 */
type TComInterpolationFilter struct{
  m_lumaFilter		[4][NTAPS_LUMA]int16;     ///< Luma filter taps
  m_chromaFilter	[8][NTAPS_CHROMA]int16; ///< Chroma filter taps
}
  
/*  
  static Void filterCopy(Int bitDepth, const Pel *src, Int srcStride, Short *dst, Int dstStride, Int width, Int height, Bool isFirst, Bool isLast);
  
  template<Int N, Bool isVertical, Bool isFirst, Bool isLast>
  static Void filter(Int bitDepth, Pel const *src, Int srcStride, Short *dst, Int dstStride, Int width, Int height, Short const *coeff);

  template<Int N>
  static Void filterHor(Int bitDepth, Pel *src, Int srcStride, Short *dst, Int dstStride, Int width, Int height,               Bool isLast, Short const *coeff);
  template<Int N>
  static Void filterVer(Int bitDepth, Pel *src, Int srcStride, Short *dst, Int dstStride, Int width, Int height, Bool isFirst, Bool isLast, Short const *coeff);

public:
  TComInterpolationFilter() {}
  ~TComInterpolationFilter() {}

  Void filterHorLuma  (Pel *src, Int srcStride, Short *dst, Int dstStride, Int width, Int height, Int frac,               Bool isLast );
  Void filterVerLuma  (Pel *src, Int srcStride, Short *dst, Int dstStride, Int width, Int height, Int frac, Bool isFirst, Bool isLast );
  Void filterHorChroma(Pel *src, Int srcStride, Short *dst, Int dstStride, Int width, Int height, Int frac,               Bool isLast );
  Void filterVerChroma(Pel *src, Int srcStride, Short *dst, Int dstStride, Int width, Int height, Int frac, Bool isFirst, Bool isLast );
  */