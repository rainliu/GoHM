package TLibCommon

import (

)


const NTAPS_LUMA        = 8 ///< Number of taps for luma
const NTAPS_CHROMA      = 4 ///< Number of taps for chroma
const IF_INTERNAL_PREC  = 14 ///< Number of bits for internal precision
const IF_FILTER_PREC    = 6 ///< Log2 of sum of filter taps
const IF_INTERNAL_OFFS  = (1<<(IF_INTERNAL_PREC-1)) ///< Offset used internally
