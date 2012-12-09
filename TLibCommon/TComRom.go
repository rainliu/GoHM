package TLibCommon

import ()

// ====================================================================================================================
// Macros
// ====================================================================================================================

const MAX_CU_DEPTH = 7                    // log2(LCUSize)
const MAX_CU_SIZE = (1 << (MAX_CU_DEPTH)) // maximum allowable size of CU
const MIN_PU_SIZE = 4
const MAX_NUM_SPU_W = (MAX_CU_SIZE / MIN_PU_SIZE) // maximum number of SPU in horizontal line

// ====================================================================================================================
// Initialize / destroy functions
// ====================================================================================================================

func InitROM() {
}
func DestroyROM() {
}
func InitSigLastScan(pBuffD, pBuffH, pBuffV *uint, iWidth, iHeight, iDepth int) {
}

var G_uiMaxCUWidth uint = MAX_CU_SIZE
var G_uiMaxCUHeight uint = MAX_CU_SIZE
var G_uiMaxCUDepth uint = MAX_CU_DEPTH
var G_uiAddCUDepth uint = 0
var G_auiZscanToRaster [MAX_NUM_SPU_W * MAX_NUM_SPU_W]uint // = { 0, };
var G_auiRasterToZscan [MAX_NUM_SPU_W * MAX_NUM_SPU_W]uint // = { 0, };
var G_auiRasterToPelX [MAX_NUM_SPU_W * MAX_NUM_SPU_W]uint  // = { 0, };
var G_auiRasterToPelY [MAX_NUM_SPU_W * MAX_NUM_SPU_W]uint  // = { 0, };

var G_bitDepthY int = 8
var G_bitDepthC int = 8

const SCALING_LIST_NUM = 6         ///< list number for quantization matrix
const SCALING_LIST_NUM_32x32 = 2   ///< list number for quantization matrix 32x32
const SCALING_LIST_REM_NUM = 6     ///< remainder of QP/6
const SCALING_LIST_START_VALUE = 8 ///< start value for dpcm mode
const MAX_MATRIX_COEF_NUM = 64     ///< max coefficient number for quantization matrix
const MAX_MATRIX_SIZE_NUM = 8      ///< max size number for quantization matrix
const SCALING_LIST_DC = 16         ///< default DC value
const (                            //enum ScalingListSize
    SCALING_LIST_4x4 = iota
    SCALING_LIST_8x8
    SCALING_LIST_16x16
    SCALING_LIST_32x32
    SCALING_LIST_SIZE_NUM
)

/*
var MatrixType [4][6]string[20] = {{"INTRA4X4_LUMA",
  "INTRA4X4_CHROMAU",
  "INTRA4X4_CHROMAV",
  "INTER4X4_LUMA",
  "INTER4X4_CHROMAU",
  "INTER4X4_CHROMAV"
  },
  {
  "INTRA8X8_LUMA",
  "INTRA8X8_CHROMAU", 
  "INTRA8X8_CHROMAV", 
  "INTER8X8_LUMA",
  "INTER8X8_CHROMAU", 
  "INTER8X8_CHROMAV"  
  },
  {
  "INTRA16X16_LUMA",
  "INTRA16X16_CHROMAU", 
  "INTRA16X16_CHROMAV", 
  "INTER16X16_LUMA",
  "INTER16X16_CHROMAU", 
  "INTER16X16_CHROMAV"  
  },
  {
  "INTRA32X32_LUMA",
  "INTER32X32_LUMA",
  },
};
static const Char MatrixType_DC[4][12][22] =
{
  {
  },
  {
  },
  {
  "INTRA16X16_LUMA_DC",
  "INTRA16X16_CHROMAU_DC", 
  "INTRA16X16_CHROMAV_DC", 
  "INTER16X16_LUMA_DC",
  "INTER16X16_CHROMAU_DC", 
  "INTER16X16_CHROMAV_DC"  
  },
  {
  "INTRA32X32_LUMA_DC",
  "INTER32X32_LUMA_DC",
  },
};
*/

var G_quantTSDefault4x4 [16]int /* = {
  16,16,16,16,
  16,16,16,16,
  16,16,16,16,
  16,16,16,16
};*/

var G_quantIntraDefault8x8 [64]int /* ={
  16,16,16,16,17,18,21,24,
  16,16,16,16,17,19,22,25,
  16,16,17,18,20,22,25,29,
  16,16,18,21,24,27,31,36,
  17,17,20,24,30,35,41,47,
  18,19,22,27,35,44,54,65,
  21,22,25,31,41,54,70,88,
  24,25,29,36,47,65,88,115
};*/

var G_quantInterDefault8x8 [64]int /* = {
  16,16,16,16,17,18,20,24,
  16,16,16,17,18,20,24,25,
  16,16,17,18,20,24,25,28,
  16,17,18,20,24,25,28,33,
  17,18,20,24,25,28,33,41,
  18,20,24,25,28,33,41,54,
  20,24,25,28,33,41,54,71,
  24,25,28,33,41,54,71,91
};*/
var G_scalingListSize [4]uint                    // = {16,64,256,1024}; 
var G_scalingListSizeX [4]uint                   // = { 4, 8, 16,  32};
var G_scalingListNum [SCALING_LIST_SIZE_NUM]uint // = {6,6,6,2};
var G_eTTable [4]int                             // = {0,3,1,2};
