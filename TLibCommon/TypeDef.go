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

//! \ingroup TLibCommon
//! \{
const L0363_DU_BIT_RATE = 1             ///< L0363: add bit_rate_du_value_minus1 to HRD parameters
const L0328_SPLICING = 1                ///< L0328: splicing support in HRD
const L0044_DU_DPB_OUTPUT_DELAY_HRD = 1 ///< L0044: Include dpb_output_delay_du_length_minus1 in hrd_parameters(), dpb_output_du_delay in
///<        picture timing SEI and DU information SEI
const L0045_PERSISTENCE_FLAGS = 1           ///< L0045: Replace "repetition_period" syntax elements in SEI with "persistence_flag"
const L0045_NON_NESTED_SEI_RESTRICTIONS = 1 ///< L0045; Include restriction on the order of APS and non-nested BP, PT and DU info SEI messages
const L0044_CPB_DPB_DELAY_OFFSET = 1        ///< L0044: Include syntax elements cpb_delay_offset and dpb_delay_offset in the BP SEI message
const L0047_APS_FLAGS = 1                   ///< L0047: Include full_random_access_flag and no_param_set_update_flag in the active parameter set SEI message
const L0043_TIMING_INFO = 1                 ///< L0043: Timing information is signalled in VUI outside hrd_parameters()
const L0046_RENAME_PROG_SRC_IDC = 1         ///< L0046: Rename progressive_source_idc to source_scan_type
const L0045_CONDITION_SIGNALLING = 1        ///< L0045: Condition the signaling of some syntax elements in picture timing SEI message
const L0043_MSS_IDC = 1
const L0116_ENTRY_POINT = 1
const L0363_MORE_BITS = 1
const L0363_MVP_POC = 1
const L0363_BYTE_ALIGN = 1
const L0363_SEI_ALLOW_SUFFIX = 1
const L0323_LIMIT_DEFAULT_LIST_SIZE = 1
const L0046_CONSTRAINT_FLAGS = 1
const L0255_MOVE_PPS_FLAGS = 1 ///< move some flags to earlier positions in the PPS
const L0444_FPA_TYPE = 1       ///< allow only FPA types 3, 4 and 5
const L0372 = 1
const SIGNAL_BITRATE_PICRATE_IN_VPS = 0 ///< K0125: Signal bit_rate and pic_rate in VPS
const L0232_RD_PENALTY = 1              ///< L0232: RD-penalty for 32x32 TU for intra in non-intra slices

const MAX_VPS_NUM_HRD_PARAMETERS = 1
const MAX_VPS_OP_SETS_PLUS1 = 1024
const MAX_VPS_NUH_RESERVED_ZERO_LAYER_ID_PLUS1 = 1

const RATE_CONTROL_LAMBDA_DOMAIN = 1 ///< JCTVC-K0103, rate control by R-lambda model
const L0033_RC_BUGFIX = 1            ///< JCTVC-L0033, bug fix for R-lambda model based rate control

const MAX_CPB_CNT = 32 ///< Upper bound of (cpb_cnt_minus1 + 1)
const MAX_NUM_LAYER_IDS = 64

const COEF_REMAIN_BIN_REDUCTION = 3 ///< indicates the level at which the VLC
///< transitions from Golomb-Rice to TU+EG(k)

const CU_DQP_TU_CMAX = 5 ///< max number bins for truncated unary
const CU_DQP_EG_k = 0    ///< expgolomb order

const SBH_THRESHOLD = 4 ///< I0156: value of the fixed SBH controlling threshold

const SEQUENCE_LEVEL_LOSSLESS = 0 ///< H0530: used only for sequence or frame-level lossless coding

const DISABLING_CLIP_FOR_BIPREDME = 1 ///< Ticket #175

const C1FLAG_NUMBER = 8 // maximum number of largerThan1 flag coded in one chunk :  16 in HM5
const C2FLAG_NUMBER = 1 // maximum number of largerThan2 flag coded in one chunk:  16 in HM5

const REMOVE_SAO_LCU_ENC_CONSTRAINTS_3 = 1 ///< disable the encoder constraint that conditionally disable SAO for chroma for entire slice in interleaved mode

const REMOVE_SINGLE_SEI_EXTENSION_FLAGS = 1 ///< remove display orientation SEI extension flag (there is a generic SEI extension mechanism now)

const SAO_ENCODING_CHOICE = 1 ///< I0184: picture early termination
//#if SAO_ENCODING_CHOICE
const SAO_ENCODING_RATE = 0.75
const SAO_ENCODING_CHOICE_CHROMA = 1 ///< J0044: picture early termination Luma and Chroma are handled separately
//#if SAO_ENCODING_CHOICE_CHROMA
const SAO_ENCODING_RATE_CHROMA = 0.5

//#endif
//#endif

const MAX_NUM_VPS = 16
const MAX_NUM_SPS = 16
const MAX_NUM_PPS = 64

const WEIGHTED_CHROMA_DISTORTION = 1 ///< F386: weighting of chroma for RDO
const RDOQ_CHROMA_LAMBDA = 1         ///< F386: weighting of chroma for RDOQ
const SAO_CHROMA_LAMBDA = 1          ///< F386: weighting of chroma for SAO

const MIN_SCAN_POS_CROSS = 4

const FAST_BIT_EST = 1 ///< G763: Table-based bit estimation for CABAC

const MLS_GRP_NUM = 64 ///< G644 : Max number of coefficient groups, max(16, 64)
const MLS_CG_SIZE = 4  ///< G644 : Coefficient group size of 4x4

const ADAPTIVE_QP_SELECTION = 1 ///< G382: Adaptive reconstruction levels, non-normative part for adaptive QP selection
//#if ADAPTIVE_QP_SELECTION
const ARL_C_PRECISION = 7 ///< G382: 7-bit arithmetic precision
const LEVEL_RANGE = 30    ///< G382: max coefficient level in statistics collection
//#endif

const NS_HAD = 0

const HHI_RQT_INTRA_SPEEDUP = 1     ///< tests one best mode with full rqt
const HHI_RQT_INTRA_SPEEDUP_MOD = 0 ///< tests two best modes with full rqt

//#if HHI_RQT_INTRA_SPEEDUP_MOD && !HHI_RQT_INTRA_SPEEDUP
//#error
//#endif

const VERBOSE_RATE = 0 ///< Print additional rate information in encoder

const AMVP_DECIMATION_FACTOR = 4

const SCAN_SET_SIZE = 16
const LOG2_SCAN_SET_SIZE = 4

const FAST_UDI_MAX_RDMODE_NUM = 35 ///< maximum number of RD comparison in fast-UDI estimation loop

const ZERO_MVD_EST = 0 ///< Zero Mvd Estimation in normal mode

const NUM_INTRA_MODE = 36

//#if !REMOVE_LM_CHROMA
const LM_CHROMA_IDX = 35

//#endif

const WRITE_BACK = 1     ///< Enable/disable the encoder to replace the deltaPOC and Used by current from the config file with the values derived by the refIdc parameter.
const AUTO_INTER_RPS = 1 ///< Enable/disable the automatic generation of refIdc from the deltaPOC and Used by current from the config file.
const PRINT_RPS_INFO = 0 ///< Enable/disable the printing of bits used to send the RPS.
// using one nearest frame as reference frame, and the other frames are high quality (POC%4==0) frames (1+X)
// this should be done with encoder only decision
// but because of the absence of reference frame management, the related code was hard coded currently

const RVM_VCEGAM10_M = 4

const PLANAR_IDX = 0
const VER_IDX = 26        // index for intra VERTICAL   mode
const HOR_IDX = 10        // index for intra HORIZONTAL mode
const DC_IDX = 1          // index for intra DC mode
const NUM_CHROMA_MODE = 5 // total number of chroma modes
const DM_CHROMA_IDX = 36  // chroma mode index for derived from luma intra mode

const FAST_UDI_USE_MPM = 1

const RDO_WITHOUT_DQP_BITS = 0 ///< Disable counting dQP bits in RDO-based mode decision

const FULL_NBIT = 0 ///< When enabled, compute costs using full sample bitdepth.  When disabled, compute costs as if it is 8-bit source video.

func DISTORTION_PRECISION_ADJUSTMENT(x interface{}) interface{} {
    //#if FULL_NBIT
    //# define DISTORTION_PRECISION_ADJUSTMENT(x) 0
    //#else
    return x // DISTORTION_PRECISION_ADJUSTMENT(x) (x)
    //#endif
}

const LOG2_MAX_NUM_COLUMNS_MINUS1 = 7
const LOG2_MAX_NUM_ROWS_MINUS1 = 7
const LOG2_MAX_COLUMN_WIDTH = 13
const LOG2_MAX_ROW_HEIGHT = 13

const MATRIX_MULT = 0 // Brute force matrix multiplication instead of partial butterfly

const REG_DCT = 65535

const AMP_SAD = 1         ///< dedicated SAD functions for AMP
const AMP_ENC_SPEEDUP = 1 ///< encoder only speed-up by AMP mode skipping
//#if AMP_ENC_SPEEDUP
const AMP_MRG = 1 ///< encoder only force merge for AMP partition (no motion search for AMP)
//#endif

const SCALING_LIST_OUTPUT_RESULT = 0 //JCTVC-G880/JCTVC-G1016 quantization matrices

const CABAC_INIT_PRESENT_FLAG = 1

// ====================================================================================================================
// Basic type redefinition
// ====================================================================================================================
/*
typedef       void                Void;
typedef       bool                Bool;

typedef       char                Char;
typedef       unsigned char       UChar;
typedef       short               Short;
typedef       unsigned short      UShort;
typedef       int                 Int;
typedef       unsigned int        UInt;
typedef       double              Double;
typedef       float               Float;

// ====================================================================================================================
// 64-bit integer type
// ====================================================================================================================
typedef       long long           Int64;
typedef       unsigned long long  UInt64;
*/
// ====================================================================================================================
// Type definition
// ====================================================================================================================

type Pxl byte     ///< 8-bit pixel type
type Pel int16    ///< 16-bit pixel type
type TCoeff int32 ///< transform coefficient

/// parameters for adaptive loop filter
//class TComPicSym;

// Slice / Slice segment encoding modes
type SliceConstraint uint8

const ( //SliceConstraint
    NO_SLICES             = 0 ///< don't use slices / slice segments
    FIXED_NUMBER_OF_LCU   = 1 ///< Limit maximum number of largest coding tree blocks in a slice / slice segments
    FIXED_NUMBER_OF_BYTES = 2 ///< Limit maximum number of bytes in a slice / slice segment
    FIXED_NUMBER_OF_TILES = 3 ///< slices / slice segments span an integer number of tiles
)

const NUM_DOWN_PART = 4

type SAOTypeLen uint8

const ( //enum SAOTypeLen
    SAO_EO_LEN         = 4
    SAO_BO_LEN         = 4
    SAO_MAX_BO_CLASSES = 32
)

type SAOType uint8

const ( //enum SAOType
    SAO_EO_0 = iota
    SAO_EO_1
    SAO_EO_2
    SAO_EO_3
    SAO_BO
    MAX_NUM_SAO_TYPE
)

type SAOQTPart struct {
    iBestType  int
    iLength    int
    subTypeIdx int ///< indicates EO class or BO band position
    iOffset    [4]int
    StartCUX   int
    StartCUY   int
    EndCUX     int
    EndCUY     int

    PartIdx   int
    PartLevel int
    PartCol   int
    PartRow   int

    DownPartsIdx [NUM_DOWN_PART]int
    UpPartIdx    int

    bSplit bool

    //---- encoder only start -----//
    bProcessed bool
    dMinCost   float64
    iMinDist   int
    iMinRate   int
    //---- encoder only end -----//
}

type SaoLcuParam struct {
    MergeUpFlag   bool
    MergeLeftFlag bool
    TypeIdx       int
    SubTypeIdx    int ///< indicates EO class or BO band position
    Offset        [4]int
    PartIdx       int
    PartIdxTmp    int
    Length        int
}

type SAOParam struct {
    SaoFlag       [2]bool
    SaoPart       [3][]SAOQTPart
    MaxSplitLevel int
    OneUnitFlag   [3]bool
    SaoLcuParam   [3][]SaoLcuParam
    NumCuInHeight int
    NumCuInWidth  int
}

/// parameters for deblocking filter
type LFCUParam struct {
    bInternalEdge bool ///< indicates internal edge
    bLeftEdge     bool ///< indicates left edge
    bTopEdge      bool ///< indicates top edge
}

// ====================================================================================================================
// Enumeration
// ====================================================================================================================

/// supported slice type
type SliceType uint8

const ( //enum SliceType
    B_SLICE = iota
    P_SLICE
    I_SLICE
)

/// chroma formats (according to semantics of chroma_format_idc)
type ChromaFormat uint8

const ( //enum ChromaFormat
    CHROMA_400 = 0
    CHROMA_420 = 1
    CHROMA_422 = 2
    CHROMA_444 = 3
)

/// supported partition shape
type PartSize uint8

const ( //enum PartSize
    SIZE_2Nx2N = iota ///< symmetric motion partition,  2Nx2N
    SIZE_2NxN         ///< symmetric motion partition,  2Nx N
    SIZE_Nx2N         ///< symmetric motion partition,   Nx2N
    SIZE_NxN          ///< symmetric motion partition,   Nx N
    SIZE_2NxnU        ///< asymmetric motion partition, 2Nx( N/2) + 2Nx(3N/2)
    SIZE_2NxnD        ///< asymmetric motion partition, 2Nx(3N/2) + 2Nx( N/2)
    SIZE_nLx2N        ///< asymmetric motion partition, ( N/2)x2N + (3N/2)x2N
    SIZE_nRx2N        ///< asymmetric motion partition, (3N/2)x2N + ( N/2)x2N
    SIZE_NONE  = 15
)

/// supported prediction type
type PredMode uint8

const ( //enum PredMode
    MODE_INTER = 0 ///< inter-prediction mode
    MODE_INTRA = 1 ///< intra-prediction mode
    MODE_NONE  = 15
)

/// texture component type
type TextType uint8

const ( //enum TextType
    TEXT_LUMA     = iota ///< luma
    TEXT_CHROMA          ///< chroma (U+V)
    TEXT_CHROMA_U        ///< chroma U
    TEXT_CHROMA_V        ///< chroma V
    TEXT_ALL             ///< Y+U+V
    TEXT_NONE     = 15
)

/// reference list index
type RefPicList uint8

const ( //enum RefPicList
    REF_PIC_LIST_0 = 0   ///< reference list 0
    REF_PIC_LIST_1 = 1   ///< reference list 1
    REF_PIC_LIST_C = 2   ///< combined reference list for uni-prediction in B-Slices
    REF_PIC_LIST_X = 100 ///< special mark
)

/// distortion function index
type DFunc uint8

const ( //enum DFunc
    DF_DEFAULT = 0
    DF_SSE     = 1 ///< general size SSE
    DF_SSE4    = 2 ///<   4xM SSE
    DF_SSE8    = 3 ///<   8xM SSE
    DF_SSE16   = 4 ///<  16xM SSE
    DF_SSE32   = 5 ///<  32xM SSE
    DF_SSE64   = 6 ///<  64xM SSE
    DF_SSE16N  = 7 ///< 16NxM SSE

    DF_SAD    = 8  ///< general size SAD
    DF_SAD4   = 9  ///<   4xM SAD
    DF_SAD8   = 10 ///<   8xM SAD
    DF_SAD16  = 11 ///<  16xM SAD
    DF_SAD32  = 12 ///<  32xM SAD
    DF_SAD64  = 13 ///<  64xM SAD
    DF_SAD16N = 14 ///< 16NxM SAD

    DF_SADS    = 15 ///< general size SAD with step
    DF_SADS4   = 16 ///<   4xM SAD with step
    DF_SADS8   = 17 ///<   8xM SAD with step
    DF_SADS16  = 18 ///<  16xM SAD with step
    DF_SADS32  = 19 ///<  32xM SAD with step
    DF_SADS64  = 20 ///<  64xM SAD with step
    DF_SADS16N = 21 ///< 16NxM SAD with step

    DF_HADS    = 22 ///< general size Hadamard with step
    DF_HADS4   = 23 ///<   4xM HAD with step
    DF_HADS8   = 24 ///<   8xM HAD with step
    DF_HADS16  = 25 ///<  16xM HAD with step
    DF_HADS32  = 26 ///<  32xM HAD with step
    DF_HADS64  = 27 ///<  64xM HAD with step
    DF_HADS16N = 28 ///< 16NxM HAD with step

    //#if AMP_SAD
    DF_SAD12 = 43
    DF_SAD24 = 44
    DF_SAD48 = 45

    DF_SADS12 = 46
    DF_SADS24 = 47
    DF_SADS48 = 48

    DF_SSE_FRAME = 50 ///< Frame-based SSE
//#else
//  DF_SSE_FRAME = 33     ///< Frame-based SSE
//#endif
)

/// index for SBAC based RD optimization
type CI_IDX uint8

const ( //enum CI_IDX
    CI_CURR_BEST    = iota ///< best mode index
    CI_NEXT_BEST           ///< next best index
    CI_TEMP_BEST           ///< temporal index
    CI_CHROMA_INTRA        ///< chroma intra index
    CI_QT_TRAFO_TEST
    CI_QT_TRAFO_ROOT
    CI_NUM ///< total number
)

/// motion vector predictor direction used in AMVP
type MVP_DIR uint8

const ( //enum MVP_DIR
    MD_LEFT        = iota ///< MVP of left block
    MD_ABOVE              ///< MVP of above block
    MD_ABOVE_RIGHT        ///< MVP of above right block
    MD_BELOW_LEFT         ///< MVP of below left block
    MD_ABOVE_LEFT         ///< MVP of above left block
)

/// coefficient scanning type used in ACS
type COEFF_SCAN_TYPE uint8

const ( //enum COEFF_SCAN_TYPE
    SCAN_DIAG = iota ///< up-right diagonal scan
    SCAN_HOR         ///< horizontal first scan
    SCAN_VER         ///< vertical first scan
)

//namespace Profile
//{
type PROFILE uint8

const ( //enum Name
    PROFILE_NONE             = 0
    PROFILE_MAIN             = 1
    PROFILE_MAIN10           = 2
    PROFILE_MAINSTILLPICTURE = 3
)

//}

//namespace Level
//{
type TIER uint8

const ( //enum Tier
    TIER_MAIN = 0
    TIER_HIGH = 1
)

type LEVEL uint8

const ( //enum Name
    LEVELNONE = 0
    LEVEL1    = 30
    LEVEL2    = 60
    LEVEL2_1  = 63
    LEVEL3    = 90
    LEVEL3_1  = 93
    LEVEL4    = 120
    LEVEL4_1  = 123
    LEVEL5    = 150
    LEVEL5_1  = 153
    LEVEL5_2  = 156
    LEVEL6    = 180
    LEVEL6_1  = 183
    LEVEL6_2  = 186
)

//}
//! \}
