package TLibCommon

import (

)

const NV_VERSION = "9.0.0"

const MAX_GOP = 64

func Clip3( minVal, maxVal, a Pel) Pel  { 
	if a < minVal {
		a = minVal
	}else if a > maxVal {
		a = maxVal
	}

	return a; 
}  ///< general min/max clip