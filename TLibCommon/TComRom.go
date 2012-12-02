package TLibCommon

import (

)

// ====================================================================================================================
// Macros
// ====================================================================================================================

const     MAX_CU_DEPTH    =        7                           // log2(LCUSize)
const     MAX_CU_SIZE     =        (1<<(MAX_CU_DEPTH))         // maximum allowable size of CU
const     MIN_PU_SIZE     =        4
const     MAX_NUM_SPU_W   =        (MAX_CU_SIZE/MIN_PU_SIZE)   // maximum number of SPU in horizontal line


var g_uiMaxCUWidth  uint = MAX_CU_SIZE;
var g_uiMaxCUHeight uint = MAX_CU_SIZE;
var g_uiMaxCUDepth  uint = MAX_CU_DEPTH;
var g_uiAddCUDepth  uint = 0;
var g_auiZscanToRaster [ MAX_NUM_SPU_W*MAX_NUM_SPU_W ]uint;// = { 0, };
var g_auiRasterToZscan [ MAX_NUM_SPU_W*MAX_NUM_SPU_W ]uint;// = { 0, };
var g_auiRasterToPelX  [ MAX_NUM_SPU_W*MAX_NUM_SPU_W ]uint;// = { 0, };
var g_auiRasterToPelY  [ MAX_NUM_SPU_W*MAX_NUM_SPU_W ]uint;// = { 0, };

var g_bitDepthY      int = 8;
var g_bitDepthC 	 int = 8;