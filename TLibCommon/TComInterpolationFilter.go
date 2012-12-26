package TLibCommon

import ()

const NTAPS_LUMA = 8                                   ///< Number of taps for luma
const NTAPS_CHROMA = 4                                 ///< Number of taps for chroma
const IF_INTERNAL_PREC = 14                            ///< Number of bits for internal precision
const IF_FILTER_PREC = 6                               ///< Log2 of sum of filter taps
const IF_INTERNAL_OFFS = (1 << (IF_INTERNAL_PREC - 1)) ///< Offset used internally

/**
 * \brief Interpolation filter class
 */
type TComInterpolationFilter struct {
    m_lumaFilter   [4][NTAPS_LUMA]int16   ///< Luma filter taps
    m_chromaFilter [8][NTAPS_CHROMA]int16 ///< Chroma filter taps
}

func NewTComInterpolationFilter() *TComInterpolationFilter{
	return &TComInterpolationFilter{};
}

func (this *TComInterpolationFilter) FilterHorLuma  (src []Pel, srcStride int, dst []Pel, dstStride, width, height, frac int,            isLast bool){
}
func (this *TComInterpolationFilter) FilterVerLuma  (src []Pel, srcStride int, dst []Pel, dstStride, width, height, frac int,  isFirst,  isLast bool){
}
func (this *TComInterpolationFilter) FilterHorChroma(src []Pel, srcStride int, dst []Pel, dstStride, width, height, frac int,            isLast bool){
}
func (this *TComInterpolationFilter) FilterVerChroma(src []Pel, srcStride int, dst []Pel, dstStride, width, height, frac int,  isFirst,  isLast bool){
}

func (this *TComInterpolationFilter) filterCopy( bitDepth int, src []Pel, srcStride int, dst []Pel, dstStride, width, height int, isFirst, isLast bool){
}

/*
  template<Int N, Bool isVertical, Bool isFirst, Bool isLast>
func (this *TComInterpolationFilter) filter(Int bitDepth, Pel const *src, Int srcStride, Short *dst, Int dstStride, Int width, Int height, Short const *coeff);

  template<Int N>
func (this *TComInterpolationFilter) filterHor(Int bitDepth, Pel *src, Int srcStride, Short *dst, Int dstStride, Int width, Int height,               Bool isLast, Short const *coeff);
  template<Int N>
func (this *TComInterpolationFilter) filterVer(Int bitDepth, Pel *src, Int srcStride, Short *dst, Int dstStride, Int width, Int height, Bool isFirst, Bool isLast, Short const *coeff);
*/