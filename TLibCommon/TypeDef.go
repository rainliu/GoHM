package TLibCommon

import (

)

//! \ingroup TLibCommon
//! \{
const SAVE_BITS_REFPICLIST_MOD_FLAG=               1  ///< K0224 Proposal#1: Send ref_pic_list_modification_flag_lX only when NumPocTotalCurr is greater than 1.

const USE_PIC_CHROMA_QP_OFFSETS_IN_DEBLOCKING=     1  ///< K0220: Use picture-based chroma QP offsets in deblocking filter.

const REMOVE_BURST_IPCM=                1  /// Ticket763
const REMOVE_ENTROPY_SLICES= 1

const DEPENDENT_SLICE_SEGMENT_FLAGS=   1   ///< K0184: Move dependent_slice_enabled_flag after seq_parameter_set_id in PPS.
                                            ///< Move dependent_slice_flag between pic_parameter_set_id and slice_address.
const SPS_INTER_REF_SET_PRED=      1   ///< K0136: Not send inter_ref_pic_set_prediction_flag for index 0
const HM9_NALU_TYPES= 1

const STRONG_INTRA_SMOOTHING=           1  ///< Enables Bilinear interploation of reference samples instead of 121 filter in intra prediction when reference samples are flat.

const RESTRICT_INTRA_BOUNDARY_SMOOTHING=    1  ///< K0380, K0186 
const LINEBUF_CLEANUP =              1 ///< K0101
const MERGE_CLEANUP_AND_K0197 =    1  //<Code cleanup and K0197: removal of indirect use of A1 and B1 in merging candidate list construction.
const RPL_INIT_FIX = 1 ///< K0255 2nd part (editorial)

const MAX_CPB_CNT       =               32  ///< Upper bound of (cpb_cnt_minus1 + 1)
const MAX_NUM_LAYER_IDS =               64

const FLAT_4x4_DSL = 1 ///< Use flat 4x4 default scaling list (see notes on K0203)

const RDOQ_TRANSFORMSKIP   =       1   // Enable RDOQ for transform skip (see noted on K0245)

const COEF_REMAIN_BIN_REDUCTION  =      3 ///< indicates the level at which the VLC 
                                           ///< transitions from Golomb-Rice to TU+EG(k)

const CU_DQP_TU_CMAX= 5                   ///< max number bins for truncated unary
const CU_DQP_EG_k= 0                      ///< expgolomb order

const SBH_THRESHOLD=                    4  ///< I0156: value of the fixed SBH controlling threshold
  
const SEQUENCE_LEVEL_LOSSLESS=           0  ///< H0530: used only for sequence or frame-level lossless coding

const DISABLING_CLIP_FOR_BIPREDME=         1  ///< Ticket #175
  
const C1FLAG_NUMBER=               8 // maximum number of largerThan1 flag coded in one chunk :  16 in HM5
const C2FLAG_NUMBER=               1 // maximum number of largerThan2 flag coded in one chunk:  16 in HM5 

const REMOVE_SAO_LCU_ENC_CONSTRAINTS_3= 1  ///< disable the encoder constraint that conditionally disable SAO for chroma for entire slice in interleaved mode

const SAO_SKIP_RIGHT =                  1  ///< H1101: disallow using unavailable pixel during RDO

const SAO_ENCODING_CHOICE =             1  ///< I0184: picture early termination
//#if SAO_ENCODING_CHOICE
const SAO_ENCODING_RATE  =              0.75
const SAO_ENCODING_CHOICE_CHROMA =      1 ///< J0044: picture early termination Luma and Chroma are handled separatenly
//#if SAO_ENCODING_CHOICE_CHROMA
const SAO_ENCODING_RATE_CHROMA =        0.5
const SAO_ENCODING_CHOICE_CHROMA_BF =   1 ///  K0156: Bug fix for SAO selection consistency
//#endif
//#endif


const MAX_NUM_VPS =                16
const MAX_NUM_SPS =                16
const MAX_NUM_PPS =                64



const WEIGHTED_CHROMA_DISTORTION=  1   ///< F386: weighting of chroma for RDO
const RDOQ_CHROMA_LAMBDA=          1   ///< F386: weighting of chroma for RDOQ
const SAO_CHROMA_LAMBDA =          1   ///< F386: weighting of chroma for SAO

const MIN_SCAN_POS_CROSS =         4

const FAST_BIT_EST=                1   ///< G763: Table-based bit estimation for CABAC

const MLS_GRP_NUM  =                       64     ///< G644 : Max number of coefficient groups, max(16, 64)
const MLS_CG_SIZE  =                       4      ///< G644 : Coefficient group size of 4x4

const ADAPTIVE_QP_SELECTION  =             1      ///< G382: Adaptive reconstruction levels, non-normative part for adaptive QP selection
//#if ADAPTIVE_QP_SELECTION
const ARL_C_PRECISION     =                7      ///< G382: 7-bit arithmetic precision
const LEVEL_RANGE         =                30     ///< G382: max coefficient level in statistics collection
//#endif

const NS_HAD             =                  0

const K0251              =               1           ///< explicitly signal slice_temporal_mvp_enable_flag in non-IDR I Slices

const HHI_RQT_INTRA_SPEEDUP  =           1           ///< tests one best mode with full rqt
const HHI_RQT_INTRA_SPEEDUP_MOD =        0           ///< tests two best modes with full rqt

//#if HHI_RQT_INTRA_SPEEDUP_MOD && !HHI_RQT_INTRA_SPEEDUP
//#error
//#endif

const VERBOSE_RATE= 0 ///< Print additional rate information in encoder

const AMVP_DECIMATION_FACTOR=            4

const SCAN_SET_SIZE =                    16
const LOG2_SCAN_SET_SIZE=                4

const FAST_UDI_MAX_RDMODE_NUM  =             35          ///< maximum number of RD comparison in fast-UDI estimation loop 

const ZERO_MVD_EST  =                        0           ///< Zero Mvd Estimation in normal mode

const NUM_INTRA_MODE= 36
//#if !REMOVE_LM_CHROMA
const LM_CHROMA_IDX = 35
//#endif

const WRITE_BACK   =                   1           ///< Enable/disable the encoder to replace the deltaPOC and Used by current from the config file with the values derived by the refIdc parameter.
const AUTO_INTER_RPS  =                1           ///< Enable/disable the automatic generation of refIdc from the deltaPOC and Used by current from the config file.
const PRINT_RPS_INFO =                 0           ///< Enable/disable the printing of bits used to send the RPS.
                                                    // using one nearest frame as reference frame, and the other frames are high quality (POC%4==0) frames (1+X)
                                                    // this should be done with encoder only decision
                                                    // but because of the absence of reference frame management, the related code was hard coded currently

const RVM_VCEGAM10_M= 4

const PLANAR_IDX =            0
const VER_IDX    =            26                    // index for intra VERTICAL   mode
const HOR_IDX    =            10                    // index for intra HORIZONTAL mode
const DC_IDX     =            1                     // index for intra DC mode
const NUM_CHROMA_MODE=        5                     // total number of chroma modes
const DM_CHROMA_IDX  =        36                    // chroma mode index for derived from luma intra mode


const FAST_UDI_USE_MPM= 1

const RDO_WITHOUT_DQP_BITS=              0           ///< Disable counting dQP bits in RDO-based mode decision

const FULL_NBIT= 0 ///< When enabled, compute costs using full sample bitdepth.  When disabled, compute costs as if it is 8-bit source video.
//#if FULL_NBIT
//# define DISTORTION_PRECISION_ADJUSTMENT(x) 0
//#else
//#define DISTORTION_PRECISION_ADJUSTMENT(x) (x)
//#endif


const AD_HOC_SLICES_FIXED_NUMBER_OF_LCU_IN_SLICE =     1          ///< OPTION IDENTIFIER. mode==1 -> Limit maximum number of largest coding tree blocks in a slice
const AD_HOC_SLICES_FIXED_NUMBER_OF_BYTES_IN_SLICE=    2          ///< OPTION IDENTIFIER. mode==2 -> Limit maximum number of bins/bits in a slice
const AD_HOC_SLICES_FIXED_NUMBER_OF_TILES_IN_SLICE=    3

const DEPENDENT_SLICES  =     1 ///< JCTVC-I0229
// Dependent slice options
const SHARP_FIXED_NUMBER_OF_LCU_IN_DEPENDENT_SLICE =           1          ///< OPTION IDENTIFIER. Limit maximum number of largest coding tree blocks in an dependent slice
const SHARP_MULTIPLE_CONSTRAINT_BASED_DEPENDENT_SLICE =        2          ///< OPTION IDENTIFIER. Limit maximum number of bins/bits in an dependent slice
//#if DEPENDENT_SLICES
const FIXED_NUMBER_OF_TILES_IN_DEPENDENT_SLICE =         3 // JCTVC-I0229
//#endif

const LOG2_MAX_NUM_COLUMNS_MINUS1 =       7
const LOG2_MAX_NUM_ROWS_MINUS1    =       7
const LOG2_MAX_COLUMN_WIDTH       =       13
const LOG2_MAX_ROW_HEIGHT         =       13

const MATRIX_MULT                  =           0   // Brute force matrix multiplication instead of partial butterfly

const REG_DCT = 65535

const AMP_SAD   =                            1           ///< dedicated SAD functions for AMP
const AMP_ENC_SPEEDUP =                      1           ///< encoder only speed-up by AMP mode skipping
//#if AMP_ENC_SPEEDUP
const AMP_MRG  =                             1           ///< encoder only force merge for AMP partition (no motion search for AMP)
//#endif

const SCALING_LIST_OUTPUT_RESULT =   0 //JCTVC-G880/JCTVC-G1016 quantization matrices

const CABAC_INIT_PRESENT_FLAG  =   1

// ====================================================================================================================
// VPS constants
// ====================================================================================================================
const MAX_LAYER_NUM =                     10

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

type       Pxl			   	byte;        ///< 8-bit pixel type
type       Pel				int16;       ///< 16-bit pixel type
type       TCoeff			int32;       ///< transform coefficient


/// parameters for adaptive loop filter
//class TComPicSym;

const NUM_DOWN_PART = 4

type SAOTypeLen	uint8
const (//enum SAOTypeLen
  SAO_EO_LEN    = 4 
  SAO_BO_LEN    = 4
  SAO_MAX_BO_CLASSES = 32
)

type SAOType 	uint8
const (//enum SAOType
  SAO_EO_0 = iota
  SAO_EO_1
  SAO_EO_2 
  SAO_EO_3
  SAO_BO
  MAX_NUM_SAO_TYPE
)
/*
typedef struct _SaoQTPart
{
  Int         iBestType;
  Int         iLength;
  Int         subTypeIdx ;                 ///< indicates EO class or BO band position
  Int         iOffset[4];
  Int         StartCUX;
  Int         StartCUY;
  Int         EndCUX;
  Int         EndCUY;

  Int         PartIdx;
  Int         PartLevel;
  Int         PartCol;
  Int         PartRow;

  Int         DownPartsIdx[NUM_DOWN_PART];
  Int         UpPartIdx;

  Bool        bSplit;

  //---- encoder only start -----//
  Bool        bProcessed;
  Double      dMinCost;
  Int64       iMinDist;
  Int         iMinRate;
  //---- encoder only end -----//
} SAOQTPart;

typedef struct _SaoLcuParam
{
  Bool       mergeUpFlag;
  Bool       mergeLeftFlag;
  Int        typeIdx;
  Int        subTypeIdx;                  ///< indicates EO class or BO band position
  Int        offset[4];
  Int        partIdx;
  Int        partIdxTmp;
  Int        length;
} SaoLcuParam;

struct SAOParam
{
  Bool       bSaoFlag[2];
  SAOQTPart* psSaoPart[3];
  Int        iMaxSplitLevel;
  Bool         oneUnitFlag[3];
  SaoLcuParam* saoLcuParam[3];
  Int          numCuInHeight;
  Int          numCuInWidth;
  ~SAOParam();
};

/// parameters for deblocking filter
typedef struct _LFCUParam
{
  Bool bInternalEdge;                     ///< indicates internal edge
  Bool bLeftEdge;                         ///< indicates left edge
  Bool bTopEdge;                          ///< indicates top edge
} LFCUParam;
*/
// ====================================================================================================================
// Enumeration
// ====================================================================================================================

/// supported slice type
type SliceType	uint8
const (//enum SliceType
  B_SLICE = iota
  P_SLICE
  I_SLICE
)

/// chroma formats (according to semantics of chroma_format_idc)
type ChromaFormat	uint8
const (//enum ChromaFormat
  CHROMA_400  = 0
  CHROMA_420  = 1
  CHROMA_422  = 2
  CHROMA_444  = 3
)

/// supported partition shape
type PartSize	uint8
const (//enum PartSize
  SIZE_2Nx2N = iota    ///< symmetric motion partition,  2Nx2N
  SIZE_2NxN             ///< symmetric motion partition,  2Nx N
  SIZE_Nx2N             ///< symmetric motion partition,   Nx2N
  SIZE_NxN              ///< symmetric motion partition,   Nx N
  SIZE_2NxnU            ///< asymmetric motion partition, 2Nx( N/2) + 2Nx(3N/2)
  SIZE_2NxnD            ///< asymmetric motion partition, 2Nx(3N/2) + 2Nx( N/2)
  SIZE_nLx2N            ///< asymmetric motion partition, ( N/2)x2N + (3N/2)x2N
  SIZE_nRx2N            ///< asymmetric motion partition, (3N/2)x2N + ( N/2)x2N
  SIZE_NONE = 15
)

/// supported prediction type
type PredMode	uint8
const (//enum PredMode
  MODE_INTER = 0           ///< inter-prediction mode
  MODE_INTRA = 1           ///< intra-prediction mode
  MODE_NONE = 15
)

/// texture component type
type TextType uint8
const (//enum TextType
  TEXT_LUMA = iota            ///< luma
  TEXT_CHROMA          ///< chroma (U+V)
  TEXT_CHROMA_U        ///< chroma U
  TEXT_CHROMA_V        ///< chroma V
  TEXT_ALL             ///< Y+U+V
  TEXT_NONE = 15
)

/// reference list index
type RefPicList	uint8
const (//enum RefPicList
  REF_PIC_LIST_0 = 0   ///< reference list 0
  REF_PIC_LIST_1 = 1   ///< reference list 1
  REF_PIC_LIST_C = 2   ///< combined reference list for uni-prediction in B-Slices
  REF_PIC_LIST_X = 100  ///< special mark
)

/// distortion function index
type Dunc 	uint8
const (//enum DFunc
  DF_DEFAULT  = 0
  DF_SSE      = 1      ///< general size SSE
  DF_SSE4     = 2      ///<   4xM SSE
  DF_SSE8     = 3      ///<   8xM SSE
  DF_SSE16    = 4      ///<  16xM SSE
  DF_SSE32    = 5      ///<  32xM SSE
  DF_SSE64    = 6      ///<  64xM SSE
  DF_SSE16N   = 7      ///< 16NxM SSE
  
  DF_SAD      = 8      ///< general size SAD
  DF_SAD4     = 9      ///<   4xM SAD
  DF_SAD8     = 10     ///<   8xM SAD
  DF_SAD16    = 11     ///<  16xM SAD
  DF_SAD32    = 12     ///<  32xM SAD
  DF_SAD64    = 13     ///<  64xM SAD
  DF_SAD16N   = 14     ///< 16NxM SAD
  
  DF_SADS     = 15     ///< general size SAD with step
  DF_SADS4    = 16     ///<   4xM SAD with step
  DF_SADS8    = 17     ///<   8xM SAD with step
  DF_SADS16   = 18     ///<  16xM SAD with step
  DF_SADS32   = 19     ///<  32xM SAD with step
  DF_SADS64   = 20     ///<  64xM SAD with step
  DF_SADS16N  = 21     ///< 16NxM SAD with step
  
  DF_HADS     = 22     ///< general size Hadamard with step
  DF_HADS4    = 23     ///<   4xM HAD with step
  DF_HADS8    = 24     ///<   8xM HAD with step
  DF_HADS16   = 25     ///<  16xM HAD with step
  DF_HADS32   = 26     ///<  32xM HAD with step
  DF_HADS64   = 27     ///<  64xM HAD with step
  DF_HADS16N  = 28     ///< 16NxM HAD with step
  
//#if AMP_SAD
  DF_SAD12    = 43
  DF_SAD24    = 44
  DF_SAD48    = 45

  DF_SADS12   = 46
  DF_SADS24   = 47
  DF_SADS48   = 48

  DF_SSE_FRAME = 50     ///< Frame-based SSE
//#else
//  DF_SSE_FRAME = 33     ///< Frame-based SSE
//#endif
)

/// index for SBAC based RD optimization
type CI_IDX	uint8
const (//enum CI_IDX
  CI_CURR_BEST = iota     ///< best mode index
  CI_NEXT_BEST         ///< next best index
  CI_TEMP_BEST         ///< temporal index
  CI_CHROMA_INTRA      ///< chroma intra index
  CI_QT_TRAFO_TEST
  CI_QT_TRAFO_ROOT
  CI_NUM               ///< total number
)

/// motion vector predictor direction used in AMVP
type MVP_DIR	uint8
const (//enum MVP_DIR
  MD_LEFT = iota          ///< MVP of left block
  MD_ABOVE             ///< MVP of above block
  MD_ABOVE_RIGHT       ///< MVP of above right block
  MD_BELOW_LEFT        ///< MVP of below left block
  MD_ABOVE_LEFT         ///< MVP of above left block
)

/// coefficient scanning type used in ACS
type COEFF_SCAN_TYPE	uint8
const (//enum COEFF_SCAN_TYPE
  SCAN_ZIGZAG = iota      ///< typical zigzag scan
  SCAN_HOR             ///< horizontal first scan
  SCAN_VER              ///< vertical first scan
  SCAN_DIAG              ///< up-right diagonal scan
)

//namespace Profile
//{
type PROFILE	uint8
  const(//enum Name
    PROFILE_NONE = 0
    PROFILE_MAIN = 1
    PROFILE_MAIN10 = 2
    PROFILE_MAINSTILLPICTURE = 3
  )
//}

//namespace Level
//{
type TIER	uint8
  const (//enum Tier
    TIER_MAIN = 0
    TIER_HIGH = 1
  )

type LEVEL	uint8
  const (//enum Name
    LEVELNONE= 0
    LEVEL1   = 30
    LEVEL2   = 60
    LEVEL2_1 = 63
    LEVEL3   = 90
    LEVEL3_1 = 93
    LEVEL4   = 120
    LEVEL4_1 = 123
    LEVEL5   = 150
    LEVEL5_1 = 153
    LEVEL5_2 = 156
    LEVEL6   = 180
    LEVEL6_1 = 183
    LEVEL6_2 = 186
  )
//}
//! \}
