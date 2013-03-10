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

//! \{

// ====================================================================================================================
// Macros
// ====================================================================================================================

const MAX_CU_DEPTH = 7                    // log2(LCUSize)
const MAX_CU_SIZE = (1 << (MAX_CU_DEPTH)) // maximum allowable size of CU
const MIN_PU_SIZE = 4
const MAX_NUM_SPU_W = (MAX_CU_SIZE / MIN_PU_SIZE) // maximum number of SPU in horizontal line

const MAX_TS_WIDTH = 4
const MAX_TS_HEIGHT = 4

const QUANT_IQUANT_SHIFT = 20   // Q(QP%6) * IQ(QP%6) = 2^20
const QUANT_SHIFT = 14          // Q(4) = 2^14
const SCALE_BITS = 15           // Inherited from TMuC, pressumably for fractional bit estimates in RDOQ
const MAX_TR_DYNAMIC_RANGE = 15 // Maximum transform dynamic range (excluding sign bit)

const SHIFT_INV_1ST = 7  // Shift after first inverse transform stage
const SHIFT_INV_2ND = 12 // Shift after second inverse transform stage

const TRACE_FRAME = 0x0001
const TRACE_SLICE = 0x0002
const TRACE_LCU = 0x0004
const TRACE_CU = 0x0008
const TRACE_PU = 0x0010
const TRACE_TU = 0x0020
const TRACE_COEF = 0x0040
const TRACE_RESI = 0x0080
const TRACE_PRED = 0x0100
const TRACE_RECON = 0x0200

const TRACE_LEVEL = 0x0003 //0x0FFF //0x007F     
const TRACE_CABAC = true

const COUNTER_START = 1
const COUNTER_END = 0 //( UInt64(1) << 63 )

const SCALING_LIST_NUM = 6         ///< list number for quantization matrix
const SCALING_LIST_NUM_32x32 = 2   ///< list number for quantization matrix 32x32
const SCALING_LIST_REM_NUM = 6     ///< remainder of QP/6
const SCALING_LIST_START_VALUE = 8 ///< start value for dpcm mode
const MAX_MATRIX_COEF_NUM = 64     ///< max coefficient number for quantization matrix
const MAX_MATRIX_SIZE_NUM = 8      ///< max size number for quantization matrix
const SCALING_LIST_DC = 16         ///< default DC value

type ScalingListSize uint8

//enum ScalingListSize
const (
    SCALING_LIST_4x4 = iota
    SCALING_LIST_8x8
    SCALING_LIST_16x16
    SCALING_LIST_32x32
    SCALING_LIST_SIZE_NUM
)

var MatrixType = [4][6]string{
    {
        "INTRA4X4_LUMA",
        "INTRA4X4_CHROMAU",
        "INTRA4X4_CHROMAV",
        "INTER4X4_LUMA",
        "INTER4X4_CHROMAU",
        "INTER4X4_CHROMAV"},
    {
        "INTRA8X8_LUMA",
        "INTRA8X8_CHROMAU",
        "INTRA8X8_CHROMAV",
        "INTER8X8_LUMA",
        "INTER8X8_CHROMAU",
        "INTER8X8_CHROMAV"},
    {
        "INTRA16X16_LUMA",
        "INTRA16X16_CHROMAU",
        "INTRA16X16_CHROMAV",
        "INTER16X16_LUMA",
        "INTER16X16_CHROMAU",
        "INTER16X16_CHROMAV"},
    {
        "INTRA32X32_LUMA",
        "INTER32X32_LUMA"},
}
var MatrixType_DC = [4][12]string{
    {},
    {},
    {
        "INTRA16X16_LUMA_DC",
        "INTRA16X16_CHROMAU_DC",
        "INTRA16X16_CHROMAV_DC",
        "INTER16X16_LUMA_DC",
        "INTER16X16_CHROMAU_DC",
        "INTER16X16_CHROMAV_DC"},
    {
        "INTRA32X32_LUMA_DC",
        "INTER32X32_LUMA_DC"},
}

// ====================================================================================================================
// Initialize / destroy functions
// ====================================================================================================================

//! \ingroup TLibCommon
//! \{

// initialize ROM variables
func InitROM() {
    var i, c int

    // g_aucConvertToBit[ x ]: log2(x/4), if x=4 -> 0, x=8 -> 1, x=16 -> 2, ...
    //  ::memset( G_aucConvertToBit,   -1, sizeof( g_aucConvertToBit ) );
    for i = 0; i < MAX_CU_SIZE+1; i++ {
        G_aucConvertToBit[i] = -1
    }

    c = 0
    for i = 4; i < MAX_CU_SIZE; i *= 2 {
        G_aucConvertToBit[i] = int8(c)
        c++
    }
    G_aucConvertToBit[i] = int8(c)

    // G_auiFrameScanXY[ G_aucConvertToBit[ transformSize ] ]: zigzag scan array for transformSize
    c = 2
    for i = 0; i < MAX_CU_DEPTH; i++ {
        G_auiSigLastScan[0][i] = make([]uint, c*c)
        G_auiSigLastScan[1][i] = make([]uint, c*c)
        G_auiSigLastScan[2][i] = make([]uint, c*c)

        InitSigLastScan(G_auiSigLastScan[0][i], G_auiSigLastScan[1][i], G_auiSigLastScan[2][i], int(c), int(c))

        c <<= 1
    }
}

func DestroyROM() {
    /*var i int;

      for i=0; i<MAX_CU_DEPTH; i++ {
        delete[] G_auiSigLastScan[0][i];
        delete[] G_auiSigLastScan[1][i];
        delete[] G_auiSigLastScan[2][i];
      }*/
}

// ====================================================================================================================
// Data structure related table & variable
// ====================================================================================================================

//var G_uiMaxCUWidth = uint(MAX_CU_SIZE)
//var G_uiMaxCUHeight = uint(MAX_CU_SIZE)
//var G_uiMaxCUDepth = uint(MAX_CU_DEPTH)
//var G_uiAddCUDepth = uint(0)

var G_auiZscanToRaster = [MAX_NUM_SPU_W * MAX_NUM_SPU_W]uint{0}
var G_auiRasterToZscan = [MAX_NUM_SPU_W * MAX_NUM_SPU_W]uint{0}
var G_auiRasterToPelX = [MAX_NUM_SPU_W * MAX_NUM_SPU_W]uint{0}
var G_auiRasterToPelY = [MAX_NUM_SPU_W * MAX_NUM_SPU_W]uint{0}

//#if !LINEBUF_CLEANUP
//UInt G_motionRefer   [ MAX_NUM_SPU_W*MAX_NUM_SPU_W ] = { 0, };
//#endif

var G_auiPUOffset = [8]uint{0, 8, 4, 4, 2, 10, 1, 5}

func InitZscanToRaster(iMaxDepth, iDepth int, uiStartVal uint, rpuiCurrIdx []uint, rpIdx *uint) {
    iStride := uint(1) << uint(iMaxDepth-1)

    if iDepth == iMaxDepth {
        rpuiCurrIdx[*rpIdx] = uiStartVal
        (*rpIdx)++
        //rpuiCurrIdx++;
    } else {
        iStep := iStride >> uint(iDepth)
        InitZscanToRaster(iMaxDepth, iDepth+1, uiStartVal, rpuiCurrIdx, rpIdx)
        InitZscanToRaster(iMaxDepth, iDepth+1, uiStartVal+iStep, rpuiCurrIdx, rpIdx)
        InitZscanToRaster(iMaxDepth, iDepth+1, uiStartVal+iStep*iStride, rpuiCurrIdx, rpIdx)
        InitZscanToRaster(iMaxDepth, iDepth+1, uiStartVal+iStep*iStride+iStep, rpuiCurrIdx, rpIdx)
    }
}

func InitRasterToZscan(uiMaxCUWidth, uiMaxCUHeight, uiMaxDepth uint) {
    uiMinCUWidth := uiMaxCUWidth >> (uiMaxDepth - 1)
    uiMinCUHeight := uiMaxCUHeight >> (uiMaxDepth - 1)

    uiNumPartInWidth := uint(uiMaxCUWidth / uiMinCUWidth)
    uiNumPartInHeight := uint(uiMaxCUHeight / uiMinCUHeight)

    for i := uint(0); i < uiNumPartInWidth*uiNumPartInHeight; i++ {
        G_auiRasterToZscan[G_auiZscanToRaster[i]] = i
    }
}

/*
#if !LINEBUF_CLEANUP
Void initMotionReferIdx ( UInt uiMaxCUWidth, UInt uiMaxCUHeight, UInt uiMaxDepth )
{
  Int  minSUWidth  = (Int)uiMaxCUWidth  >> ( (Int)uiMaxDepth - 1 );
  Int  minSUHeight = (Int)uiMaxCUHeight >> ( (Int)uiMaxDepth - 1 );

  Int  numPartInWidth  = (Int)uiMaxCUWidth  / (Int)minSUWidth;
  Int  numPartInHeight = (Int)uiMaxCUHeight / (Int)minSUHeight;

  for ( Int i = 0; i < numPartInWidth*numPartInHeight; i++ )
  {
    G_motionRefer[i] = i;
  }

  UInt maxCUDepth = G_uiMaxCUDepth - ( G_uiAddCUDepth - 1);
  Int  minCUWidth  = (Int)uiMaxCUWidth  >> ( (Int)maxCUDepth - 1);

  if(!(minCUWidth == 8 && minSUWidth == 4)) //check if Minimum PU width == 4
  {
    return;
  }

  Int compressionNum = 2;

  for ( Int i = numPartInWidth*(numPartInHeight-1); i < numPartInWidth*numPartInHeight; i += compressionNum*2)
  {
    for ( Int j = 1; j < compressionNum; j++ )
    {
      G_motionRefer[G_auiRasterToZscan[i+j]] = G_auiRasterToZscan[i];
    }
  }

  for ( Int i = numPartInWidth*(numPartInHeight-1)+compressionNum*2-1; i < numPartInWidth*numPartInHeight; i += compressionNum*2)
  {
    for ( Int j = 1; j < compressionNum; j++ )
    {
      G_motionRefer[G_auiRasterToZscan[i-j]] = G_auiRasterToZscan[i];
    }
  }
}

#endif
*/

func InitRasterToPelXY(uiMaxCUWidth, uiMaxCUHeight, uiMaxDepth uint) {
    var i, j uint

    //UInt* uiTempX = &G_auiRasterToPelX[0];
    //UInt* uiTempY = &G_auiRasterToPelY[0];

    uiMinCUWidth := uiMaxCUWidth >> (uiMaxDepth - 1)
    uiMinCUHeight := uiMaxCUHeight >> (uiMaxDepth - 1)

    uiNumPartInWidth := uiMaxCUWidth / uiMinCUWidth
    uiNumPartInHeight := uiMaxCUHeight / uiMinCUHeight

    G_auiRasterToPelX[0] = 0 //uiTempX++;
    for i = 1; i < uiNumPartInWidth; i++ {
        G_auiRasterToPelX[i] = G_auiRasterToPelX[i-1] + uiMinCUWidth //uiTempX++;
    }
    for j = 1; j < uiNumPartInHeight; j++ {
        for i = 0; i < uiNumPartInWidth; i++ {
            G_auiRasterToPelX[j*uiNumPartInWidth+i] = G_auiRasterToPelX[i]
        }
        //memcpy(uiTempX, uiTempX-uiNumPartInWidth, sizeof(UInt)*uiNumPartInWidth);
        //uiTempX += uiNumPartInWidth;
    }

    for i = 1; i < uiNumPartInWidth*uiNumPartInHeight; i++ {
        G_auiRasterToPelY[i] = (i / uiNumPartInWidth) * uiMinCUWidth
    }
}

var G_quantScales = [6]int{26214, 23302, 20560, 18396, 16384, 14564}

var G_invQuantScales = [6]int{40, 45, 51, 57, 64, 72}

var G_aiT4 = [4][4]int16{
    {64, 64, 64, 64},
    {83, 36, -36, -83},
    {64, -64, -64, 64},
    {36, -83, 83, -36}}

var G_aiT8 = [8][8]int16{
    {64, 64, 64, 64, 64, 64, 64, 64},
    {89, 75, 50, 18, -18, -50, -75, -89},
    {83, 36, -36, -83, -83, -36, 36, 83},
    {75, -18, -89, -50, 50, 89, 18, -75},
    {64, -64, -64, 64, 64, -64, -64, 64},
    {50, -89, 18, 75, -75, -18, 89, -50},
    {36, -83, 83, -36, -36, 83, -83, 36},
    {18, -50, 75, -89, 89, -75, 50, -18}}

var G_aiT16 = [16][16]int16{
    {64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64},
    {90, 87, 80, 70, 57, 43, 25, 9, -9, -25, -43, -57, -70, -80, -87, -90},
    {89, 75, 50, 18, -18, -50, -75, -89, -89, -75, -50, -18, 18, 50, 75, 89},
    {87, 57, 9, -43, -80, -90, -70, -25, 25, 70, 90, 80, 43, -9, -57, -87},
    {83, 36, -36, -83, -83, -36, 36, 83, 83, 36, -36, -83, -83, -36, 36, 83},
    {80, 9, -70, -87, -25, 57, 90, 43, -43, -90, -57, 25, 87, 70, -9, -80},
    {75, -18, -89, -50, 50, 89, 18, -75, -75, 18, 89, 50, -50, -89, -18, 75},
    {70, -43, -87, 9, 90, 25, -80, -57, 57, 80, -25, -90, -9, 87, 43, -70},
    {64, -64, -64, 64, 64, -64, -64, 64, 64, -64, -64, 64, 64, -64, -64, 64},
    {57, -80, -25, 90, -9, -87, 43, 70, -70, -43, 87, 9, -90, 25, 80, -57},
    {50, -89, 18, 75, -75, -18, 89, -50, -50, 89, -18, -75, 75, 18, -89, 50},
    {43, -90, 57, 25, -87, 70, 9, -80, 80, -9, -70, 87, -25, -57, 90, -43},
    {36, -83, 83, -36, -36, 83, -83, 36, 36, -83, 83, -36, -36, 83, -83, 36},
    {25, -70, 90, -80, 43, 9, -57, 87, -87, 57, -9, -43, 80, -90, 70, -25},
    {18, -50, 75, -89, 89, -75, 50, -18, -18, 50, -75, 89, -89, 75, -50, 18},
    {9, -25, 43, -57, 70, -80, 87, -90, 90, -87, 80, -70, 57, -43, 25, -9}}

var G_aiT32 = [32][32]int16{
    {64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64},
    {90, 90, 88, 85, 82, 78, 73, 67, 61, 54, 46, 38, 31, 22, 13, 4, -4, -13, -22, -31, -38, -46, -54, -61, -67, -73, -78, -82, -85, -88, -90, -90},
    {90, 87, 80, 70, 57, 43, 25, 9, -9, -25, -43, -57, -70, -80, -87, -90, -90, -87, -80, -70, -57, -43, -25, -9, 9, 25, 43, 57, 70, 80, 87, 90},
    {90, 82, 67, 46, 22, -4, -31, -54, -73, -85, -90, -88, -78, -61, -38, -13, 13, 38, 61, 78, 88, 90, 85, 73, 54, 31, 4, -22, -46, -67, -82, -90},
    {89, 75, 50, 18, -18, -50, -75, -89, -89, -75, -50, -18, 18, 50, 75, 89, 89, 75, 50, 18, -18, -50, -75, -89, -89, -75, -50, -18, 18, 50, 75, 89},
    {88, 67, 31, -13, -54, -82, -90, -78, -46, -4, 38, 73, 90, 85, 61, 22, -22, -61, -85, -90, -73, -38, 4, 46, 78, 90, 82, 54, 13, -31, -67, -88},
    {87, 57, 9, -43, -80, -90, -70, -25, 25, 70, 90, 80, 43, -9, -57, -87, -87, -57, -9, 43, 80, 90, 70, 25, -25, -70, -90, -80, -43, 9, 57, 87},
    {85, 46, -13, -67, -90, -73, -22, 38, 82, 88, 54, -4, -61, -90, -78, -31, 31, 78, 90, 61, 4, -54, -88, -82, -38, 22, 73, 90, 67, 13, -46, -85},
    {83, 36, -36, -83, -83, -36, 36, 83, 83, 36, -36, -83, -83, -36, 36, 83, 83, 36, -36, -83, -83, -36, 36, 83, 83, 36, -36, -83, -83, -36, 36, 83},
    {82, 22, -54, -90, -61, 13, 78, 85, 31, -46, -90, -67, 4, 73, 88, 38, -38, -88, -73, -4, 67, 90, 46, -31, -85, -78, -13, 61, 90, 54, -22, -82},
    {80, 9, -70, -87, -25, 57, 90, 43, -43, -90, -57, 25, 87, 70, -9, -80, -80, -9, 70, 87, 25, -57, -90, -43, 43, 90, 57, -25, -87, -70, 9, 80},
    {78, -4, -82, -73, 13, 85, 67, -22, -88, -61, 31, 90, 54, -38, -90, -46, 46, 90, 38, -54, -90, -31, 61, 88, 22, -67, -85, -13, 73, 82, 4, -78},
    {75, -18, -89, -50, 50, 89, 18, -75, -75, 18, 89, 50, -50, -89, -18, 75, 75, -18, -89, -50, 50, 89, 18, -75, -75, 18, 89, 50, -50, -89, -18, 75},
    {73, -31, -90, -22, 78, 67, -38, -90, -13, 82, 61, -46, -88, -4, 85, 54, -54, -85, 4, 88, 46, -61, -82, 13, 90, 38, -67, -78, 22, 90, 31, -73},
    {70, -43, -87, 9, 90, 25, -80, -57, 57, 80, -25, -90, -9, 87, 43, -70, -70, 43, 87, -9, -90, -25, 80, 57, -57, -80, 25, 90, 9, -87, -43, 70},
    {67, -54, -78, 38, 85, -22, -90, 4, 90, 13, -88, -31, 82, 46, -73, -61, 61, 73, -46, -82, 31, 88, -13, -90, -4, 90, 22, -85, -38, 78, 54, -67},
    {64, -64, -64, 64, 64, -64, -64, 64, 64, -64, -64, 64, 64, -64, -64, 64, 64, -64, -64, 64, 64, -64, -64, 64, 64, -64, -64, 64, 64, -64, -64, 64},
    {61, -73, -46, 82, 31, -88, -13, 90, -4, -90, 22, 85, -38, -78, 54, 67, -67, -54, 78, 38, -85, -22, 90, 4, -90, 13, 88, -31, -82, 46, 73, -61},
    {57, -80, -25, 90, -9, -87, 43, 70, -70, -43, 87, 9, -90, 25, 80, -57, -57, 80, 25, -90, 9, 87, -43, -70, 70, 43, -87, -9, 90, -25, -80, 57},
    {54, -85, -4, 88, -46, -61, 82, 13, -90, 38, 67, -78, -22, 90, -31, -73, 73, 31, -90, 22, 78, -67, -38, 90, -13, -82, 61, 46, -88, 4, 85, -54},
    {50, -89, 18, 75, -75, -18, 89, -50, -50, 89, -18, -75, 75, 18, -89, 50, 50, -89, 18, 75, -75, -18, 89, -50, -50, 89, -18, -75, 75, 18, -89, 50},
    {46, -90, 38, 54, -90, 31, 61, -88, 22, 67, -85, 13, 73, -82, 4, 78, -78, -4, 82, -73, -13, 85, -67, -22, 88, -61, -31, 90, -54, -38, 90, -46},
    {43, -90, 57, 25, -87, 70, 9, -80, 80, -9, -70, 87, -25, -57, 90, -43, -43, 90, -57, -25, 87, -70, -9, 80, -80, 9, 70, -87, 25, 57, -90, 43},
    {38, -88, 73, -4, -67, 90, -46, -31, 85, -78, 13, 61, -90, 54, 22, -82, 82, -22, -54, 90, -61, -13, 78, -85, 31, 46, -90, 67, 4, -73, 88, -38},
    {36, -83, 83, -36, -36, 83, -83, 36, 36, -83, 83, -36, -36, 83, -83, 36, 36, -83, 83, -36, -36, 83, -83, 36, 36, -83, 83, -36, -36, 83, -83, 36},
    {31, -78, 90, -61, 4, 54, -88, 82, -38, -22, 73, -90, 67, -13, -46, 85, -85, 46, 13, -67, 90, -73, 22, 38, -82, 88, -54, -4, 61, -90, 78, -31},
    {25, -70, 90, -80, 43, 9, -57, 87, -87, 57, -9, -43, 80, -90, 70, -25, -25, 70, -90, 80, -43, -9, 57, -87, 87, -57, 9, 43, -80, 90, -70, 25},
    {22, -61, 85, -90, 73, -38, -4, 46, -78, 90, -82, 54, -13, -31, 67, -88, 88, -67, 31, 13, -54, 82, -90, 78, -46, 4, 38, -73, 90, -85, 61, -22},
    {18, -50, 75, -89, 89, -75, 50, -18, -18, 50, -75, 89, -89, 75, -50, 18, 18, -50, 75, -89, 89, -75, 50, -18, -18, 50, -75, 89, -89, 75, -50, 18},
    {13, -38, 61, -78, 88, -90, 85, -73, 54, -31, 4, 22, -46, 67, -82, 90, -90, 82, -67, 46, -22, -4, 31, -54, 73, -85, 90, -88, 78, -61, 38, -13},
    {9, -25, 43, -57, 70, -80, 87, -90, 90, -87, 80, -70, 57, -43, 25, -9, -9, 25, -43, 57, -70, 80, -87, 90, -90, 87, -80, 70, -57, 43, -25, 9},
    {4, -13, 22, -31, 38, -46, 54, -61, 67, -73, 78, -82, 85, -88, 90, -90, 90, -90, 88, -85, 82, -78, 73, -67, 61, -54, 46, -38, 31, -22, 13, -4}}

var G_aucChromaScale = [58]uint8{
    0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
    17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 29, 30, 31, 32,
    33, 33, 34, 34, 35, 35, 36, 36, 37, 37, 38, 39, 40, 41, 42, 43, 44,
    45, 46, 47, 48, 49, 50, 51}

// Mode-Dependent DCT/DST
var G_as_DST_MAT_4 = [4][4]int16{
    {29, 55, 74, 84},
    {74, 74, 0, -74},
    {84, -29, -74, 55},
    {55, -84, 74, -29},
}

// ====================================================================================================================
// ADI
// ====================================================================================================================

//#if FAST_UDI_USE_MPM
var G_aucIntraModeNumFast = [7]uint8{
    3,  //   2x2
    8,  //   4x4
    8,  //   8x8
    3,  //  16x16
    3,  //  32x32
    3,  //  64x64
    3}  // 128x128

/*#else // FAST_UDI_USE_MPM
const UChar G_aucIntraModeNumFast[7] =
{
  3,  //   2x2
  9,  //   4x4
  9,  //   8x8
  4,  //  16x16   33
  4,  //  32x32   33
  5,  //  64x64   33
  4   // 128x128  33
};
#endif // FAST_UDI_USE_MPM
*/
// chroma

var G_aucConvertTxtTypeToIdx = [4]uint8{0, 1, 1, 2}

// ====================================================================================================================
// Bit-depth
// ====================================================================================================================

var G_bitDepthY = int(8)
var G_bitDepthC = int(8)

var G_uiPCMBitDepthLuma = int(8)   // PCM bit-depth
var G_uiPCMBitDepthChroma = int(8) // PCM bit-depth

// ====================================================================================================================
// Misc.
// ====================================================================================================================

var G_aucConvertToBit [MAX_CU_SIZE + 1]int8

/*
#if ENC_DEC_TRACE
FILE*  G_hTrace = NULL;
const Bool G_bEncDecTraceEnable  = true;
const Bool G_bEncDecTraceDisable = false;
Bool   G_HLSTraceEnable = true;
Bool   G_bJustDoIt = false;
UInt64 G_nSymbolCounter = 0;
Bool   G_bSliceTrace = true;
#endif
*/
var G_uiPicNo = uint(0);
// ====================================================================================================================
// Scanning order & context model mapping
// ====================================================================================================================

// scanning order table
var G_auiSigLastScan [3][MAX_CU_DEPTH][]uint

var G_sigLastScan8x8 = [3][4]uint{
    {0, 2, 1, 3},
    {0, 1, 2, 3},
    {0, 2, 1, 3}}
var G_sigLastScanCG32x32 [64]uint

var G_uiMinInGroup = [10]uint{0, 1, 2, 3, 4, 6, 8, 12, 16, 24}
var G_uiGroupIdx = [32]uint{0, 1, 2, 3, 4, 4, 5, 5, 6, 6, 6, 6, 7, 7, 7, 7, 8, 8, 8, 8, 8, 8, 8, 8, 9, 9, 9, 9, 9, 9, 9, 9}

// Rice parameters for absolute transform levels
var G_auiGoRiceRange = [5]uint{7, 14, 26, 46, 78}

var G_auiGoRicePrefixLen = [5]uint{8, 7, 6, 5, 4}

func InitSigLastScan(pBuffD, pBuffH, pBuffV []uint, iWidth, iHeight int) {
    uiNumScanPos := iWidth * iWidth
    uiNextScanPos := 0

    if iWidth < 16 {
        pBuffTemp := pBuffD
        if iWidth == 8 {
            pBuffTemp = G_sigLastScanCG32x32[:]
        }
        for uiScanLine := 0; uiNextScanPos < uiNumScanPos; uiScanLine++ {
            iPrimDim := uiScanLine
            iScndDim := 0
            for iPrimDim >= iWidth {
                iScndDim++
                iPrimDim--
            }
            for iPrimDim >= 0 && iScndDim < iWidth {
                pBuffTemp[uiNextScanPos] = uint(iPrimDim*iWidth + iScndDim)
                uiNextScanPos++
                iScndDim++
                iPrimDim--
            }
        }
    }

    if iWidth > 4 {
        uiNumBlkSide := uint(iWidth) >> 2
        uiNumBlks := uiNumBlkSide * uiNumBlkSide
        log2Blk := G_aucConvertToBit[uiNumBlkSide] + 1

        for uiBlk := uint(0); uiBlk < uiNumBlks; uiBlk++ {
            uiNextScanPos = 0
            initBlkPos := G_auiSigLastScan[SCAN_DIAG][log2Blk][uiBlk]
            if iWidth == 32 {
                initBlkPos = G_sigLastScanCG32x32[uiBlk]
            }
            offsetY := initBlkPos / uiNumBlkSide
            offsetX := initBlkPos - offsetY*uiNumBlkSide
            offsetD := 4 * (offsetX + offsetY*uint(iWidth))
            offsetScan := 16 * uiBlk
            for uiScanLine := 0; uiNextScanPos < 16; uiScanLine++ {
                iPrimDim := int(uiScanLine)
                iScndDim := 0
                for iPrimDim >= 4 {
                    iScndDim++
                    iPrimDim--
                }
                for iPrimDim >= 0 && iScndDim < 4 {
                    pBuffD[uint(uiNextScanPos)+offsetScan] = uint(iPrimDim*iWidth+iScndDim) + offsetD
                    uiNextScanPos++
                    iScndDim++
                    iPrimDim--
                }
            }
        }
    }

    uiCnt := 0
    if iWidth > 2 {
        numBlkSide := iWidth >> 2
        for blkY := 0; blkY < numBlkSide; blkY++ {
            for blkX := 0; blkX < numBlkSide; blkX++ {
                offset := blkY*4*iWidth + blkX*4
                for y := 0; y < 4; y++ {
                    for x := 0; x < 4; x++ {
                        pBuffH[uiCnt] = uint(y*iWidth + x + offset)
                        uiCnt++
                    }
                }
            }
        }

        uiCnt = 0
        for blkX := 0; blkX < numBlkSide; blkX++ {
            for blkY := 0; blkY < numBlkSide; blkY++ {
                offset := blkY*4*iWidth + blkX*4
                for x := 0; x < 4; x++ {
                    for y := 0; y < 4; y++ {
                        pBuffV[uiCnt] = uint(y*iWidth + x + offset)
                        uiCnt++
                    }
                }
            }
        }
    } else {
        for iY := 0; iY < iHeight; iY++ {
            for iX := 0; iX < iWidth; iX++ {
                pBuffH[uiCnt] = uint(iY*iWidth + iX)
                uiCnt++
            }
        }

        uiCnt = 0
        for iX := 0; iX < iWidth; iX++ {
            for iY := 0; iY < iHeight; iY++ {
                pBuffV[uiCnt] = uint(iY*iWidth + iX)
                uiCnt++
            }
        }
    }
}

var G_quantTSDefault4x4 = [16]int{
    16, 16, 16, 16,
    16, 16, 16, 16,
    16, 16, 16, 16,
    16, 16, 16, 16}

var G_quantIntraDefault8x8 = [64]int{
    16, 16, 16, 16, 17, 18, 21, 24,
    16, 16, 16, 16, 17, 19, 22, 25,
    16, 16, 17, 18, 20, 22, 25, 29,
    16, 16, 18, 21, 24, 27, 31, 36,
    17, 17, 20, 24, 30, 35, 41, 47,
    18, 19, 22, 27, 35, 44, 54, 65,
    21, 22, 25, 31, 41, 54, 70, 88,
    24, 25, 29, 36, 47, 65, 88, 115}

var G_quantInterDefault8x8 = [64]int{
    16, 16, 16, 16, 17, 18, 20, 24,
    16, 16, 16, 17, 18, 20, 24, 25,
    16, 16, 17, 18, 20, 24, 25, 28,
    16, 17, 18, 20, 24, 25, 28, 33,
    17, 18, 20, 24, 25, 28, 33, 41,
    18, 20, 24, 25, 28, 33, 41, 54,
    20, 24, 25, 28, 33, 41, 54, 71,
    24, 25, 28, 33, 41, 54, 71, 91}

var G_scalingListSize = [4]uint{16, 64, 256, 1024}
var G_scalingListSizeX = [4]uint{4, 8, 16, 32}
var G_scalingListNum = [SCALING_LIST_SIZE_NUM]uint{6, 6, 6, 2}
var G_eTTable = [4]int{0, 3, 1, 2}

//! \}
