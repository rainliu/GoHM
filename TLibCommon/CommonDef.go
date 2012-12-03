package TLibCommon

import (

)

const NV_VERSION = "9.0.0"


// ====================================================================================================================
// Common constants
// ====================================================================================================================

const _SUMMARY_OUT_   =            0           ///< print-out PSNR results of all slices to summary.txt
const _SUMMARY_PIC_   =            0           ///< print-out PSNR results for each slice type to summary.txt

const MAX_GOP         =            64          ///< max. value of hierarchical GOP size

const MAX_NUM_REF_PICS=            16          ///< max. number of pictures used for reference
const MAX_NUM_REF     =            16          ///< max. number of entries in picture reference list
const MAX_NUM_REF_LC  =            MAX_NUM_REF_PICS  // TODO: remove this macro definition (leftover from combined list concept)

//const MAX_UINT        =            0xFFFFFFFFU ///< max. value of unsigned 32-bit integer
//const MAX_INT         =            2147483647  ///< max. value of signed 32-bit integer
//const MAX_INT64       =            0x7FFFFFFFFFFFFFFFLL  ///< max. value of signed 64-bit integer
//const MAX_DOUBLE      =            1.7e+308    ///< max. value of double-type value

const MIN_QP          =            0
const MAX_QP          =            51

const NOT_VALID       =            -1


// ====================================================================================================================
// Coding tool configuration
// ====================================================================================================================

// AMVP: advanced motion vector prediction
const AMVP_MAX_NUM_CANDS     =     2           ///< max number of final candidates
const AMVP_MAX_NUM_CANDS_MEM =     3           ///< max number of candidates
// MERGE
const MRG_MAX_NUM_CANDS      =     5

// Reference memory management
const DYN_REF_FREE           =     0           ///< dynamic free of reference memories

// Explicit temporal layer QP offset
const MAX_TLAYER             =     8           ///< max number of temporal layer
const HB_LAMBDA_FOR_LDC      =     1           ///< use of B-style lambda for non-key pictures in low-delay mode

// Fast estimation of generalized B in low-delay mode
const GPB_SIMPLE             =     1           ///< Simple GPB mode
//#if     GPB_SIMPLE
const GPB_SIMPLE_UNI         =     1           ///< Simple mode for uni-direction
//#endif

// Fast ME using smoother MV assumption
const FASTME_SMOOTHER_MV     =     1           ///< reduce ME time using faster option

// Adaptive search range depending on POC difference
const ADAPT_SR_SCALE         =     1           ///< division factor for adaptive search range

const CLIP_TO_709_RANGE      =     0

// Early-skip threshold (encoder)
const EARLY_SKIP_THRES       =     1.50        ///< if RD < thres*avg[BestSkipRD]


const MAX_CHROMA_FORMAT_IDC  =    3

// TODO: Existing names used for the different NAL unit types can be altered to better reflect the names in the spec.
//       However, the names in the spec are not yet stable at this point. Once the names are stable, a cleanup 
//       effort can be done without use of macros to alter the names used to indicate the different NAL unit types.
type NalUnitType uint8

const (//enum NalUnitType
//#if HM9_NALU_TYPES
  NAL_UNIT_CODED_SLICE_TRAIL_N = iota   // 0
  NAL_UNIT_CODED_SLICE_TRAIL_R   // 1
  
  NAL_UNIT_CODED_SLICE_TSA_N     // 2
  NAL_UNIT_CODED_SLICE_TLA       // 3   // Current name in the spec: TSA_R
  
  NAL_UNIT_CODED_SLICE_STSA_N    // 4
  NAL_UNIT_CODED_SLICE_STSA_R    // 5

  NAL_UNIT_CODED_SLICE_RADL_N    // 6
  NAL_UNIT_CODED_SLICE_DLP       // 7 // Current name in the spec: RADL_R
  
  NAL_UNIT_CODED_SLICE_RASL_N    // 8
  NAL_UNIT_CODED_SLICE_TFD       // 9 // Current name in the spec: RASL_R

  NAL_UNIT_RESERVED_10
  NAL_UNIT_RESERVED_11
  NAL_UNIT_RESERVED_12
  NAL_UNIT_RESERVED_13
  NAL_UNIT_RESERVED_14
  NAL_UNIT_RESERVED_15

  NAL_UNIT_CODED_SLICE_BLA       // 16   // Current name in the spec: BLA_W_LP
  NAL_UNIT_CODED_SLICE_BLANT     // 17   // Current name in the spec: BLA_W_DLP
  NAL_UNIT_CODED_SLICE_BLA_N_LP  // 18
  NAL_UNIT_CODED_SLICE_IDR       // 19  // Current name in the spec: IDR_W_DLP
  NAL_UNIT_CODED_SLICE_IDR_N_LP  // 20
  NAL_UNIT_CODED_SLICE_CRA       // 21
  NAL_UNIT_RESERVED_22
  NAL_UNIT_RESERVED_23

  NAL_UNIT_RESERVED_24
  NAL_UNIT_RESERVED_25
  NAL_UNIT_RESERVED_26
  NAL_UNIT_RESERVED_27
  NAL_UNIT_RESERVED_28
  NAL_UNIT_RESERVED_29
  NAL_UNIT_RESERVED_30
  NAL_UNIT_RESERVED_31

  NAL_UNIT_VPS                   // 32
  NAL_UNIT_SPS                   // 33
  NAL_UNIT_PPS                   // 34
  NAL_UNIT_ACCESS_UNIT_DELIMITER // 35
  NAL_UNIT_EOS                   // 36
  NAL_UNIT_EOB                   // 37
  NAL_UNIT_FILLER_DATA           // 38
  NAL_UNIT_SEI                   // 39 Prefix SEI
  NAL_UNIT_SEI_SUFFIX            // 40 Suffix SEI

  NAL_UNIT_RESERVED_41
  NAL_UNIT_RESERVED_42
  NAL_UNIT_RESERVED_43
  NAL_UNIT_RESERVED_44
  NAL_UNIT_RESERVED_45
  NAL_UNIT_RESERVED_46
  NAL_UNIT_RESERVED_47
  NAL_UNIT_UNSPECIFIED_48
  NAL_UNIT_UNSPECIFIED_49
  NAL_UNIT_UNSPECIFIED_50
  NAL_UNIT_UNSPECIFIED_51
  NAL_UNIT_UNSPECIFIED_52
  NAL_UNIT_UNSPECIFIED_53
  NAL_UNIT_UNSPECIFIED_54
  NAL_UNIT_UNSPECIFIED_55
  NAL_UNIT_UNSPECIFIED_56
  NAL_UNIT_UNSPECIFIED_57
  NAL_UNIT_UNSPECIFIED_58
  NAL_UNIT_UNSPECIFIED_59
  NAL_UNIT_UNSPECIFIED_60
  NAL_UNIT_UNSPECIFIED_61
  NAL_UNIT_UNSPECIFIED_62
  NAL_UNIT_UNSPECIFIED_63
  NAL_UNIT_INVALID
)

func MAX(a, b Pel) Pel {
	if a < b {
		return b
	}
	
	return a
}

func MIN(a, b Pel) Pel {
	if a > b {
		return b
	}
	
	return a
}

func ClipY(a Pel) Pel { 
	if a < 0 {
		a = 0
	}else if a > (1 << uint(g_bitDepthY))-1 {
		a = (1 << uint(g_bitDepthY))-1
	}

	return a; 
}
func ClipC(a Pel) Pel { 
	if a < 0 {
		a = 0
	}else if a > (1 << uint(g_bitDepthC))-1 {
		a = (1 << uint(g_bitDepthC))-1
	}

	return a; 
}

func Clip3( minVal, maxVal, a Pel) Pel  { 
	if a < minVal {
		a = minVal
	}else if a > maxVal {
		a = maxVal
	}

	return a; 
}  ///< general min/max clip