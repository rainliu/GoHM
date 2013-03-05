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

package TLibEncoder

import (
	//"os"
	//"fmt"
    "gohm/TLibCommon"
    "math"
)

//#if FASTME_SMOOTHER_MV
const FIRSTSEARCHSTOP = true

//#else
//#define FIRSTSEARCHSTOP     0
//#endif

//#define TZ_SEARCH_CONFIGURATION
const iRaster = 5 /* TZ soll von aussen ?ergeben werden */
const bTestOtherPredictedMV = false
const bTestZeroVector = true
const bTestZeroVectorStart = false
const bTestZeroVectorStop = false
const bFirstSearchDiamond = true /* 1 = xTZ8PointDiamondSearch   0 = xTZ8PointSquareSearch */
const bFirstSearchStop = FIRSTSEARCHSTOP
const uiFirstSearchRounds = 3 /* first search stop X rounds after best match (must be >=1) */
const bEnableRasterSearch = true
const bAlwaysRasterSearch = false      /* ===== 1: BETTER but factor 2 slower ===== */
const bRasterRefinementEnable = false  /* enable either raster refinement or star refinement */
const bRasterRefinementDiamond = false /* 1 = xTZ8PointDiamondSearch   0 = xTZ8PointSquareSearch */
const bStarRefinementEnable = true     /* enable either star refinement or raster refinement */
const bStarRefinementDiamond = true    /* 1 = xTZ8PointDiamondSearch   0 = xTZ8PointSquareSearch */
const bStarRefinementStop = false
const uiStarRefinementRounds = 2 /* star refinement stop X rounds after best match (must be >=1) */

var s_acMvRefineH = [9]*TLibCommon.TComMv{TLibCommon.NewTComMv(0, 0), // 0
    TLibCommon.NewTComMv(0, -1),  // 1
    TLibCommon.NewTComMv(0, 1),   // 2
    TLibCommon.NewTComMv(-1, 0),  // 3
    TLibCommon.NewTComMv(1, 0),   // 4
    TLibCommon.NewTComMv(-1, -1), // 5
    TLibCommon.NewTComMv(1, -1),  // 6
    TLibCommon.NewTComMv(-1, 1),  // 7
    TLibCommon.NewTComMv(1, 1),   // 8
}

var s_acMvRefineQ = [9]*TLibCommon.TComMv{TLibCommon.NewTComMv(0, 0), // 0
    TLibCommon.NewTComMv(0, -1),  // 1
    TLibCommon.NewTComMv(0, 1),   // 2
    TLibCommon.NewTComMv(-1, -1), // 5
    TLibCommon.NewTComMv(1, -1),  // 6
    TLibCommon.NewTComMv(-1, 0),  // 3
    TLibCommon.NewTComMv(1, 0),   // 4
    TLibCommon.NewTComMv(-1, 1),  // 7
    TLibCommon.NewTComMv(1, 1),   // 8
}

var s_auiDFilter = [9]uint{0, 1, 0,
    2, 3, 2,
    0, 1, 0,
}

type IntTZSearchStruct struct {
    piRefY         []TLibCommon.Pel
    iOffset		   int
    iYStride       int
    iBestX         int
    iBestY         int
    uiBestRound    uint
    uiBestDistance uint
    uiBestSad      uint
    ucPointNr      byte
}

// ====================================================================================================================
// Class definition
// ====================================================================================================================

/// encoder search class
type TEncSearch struct {
    TLibCommon.TComPrediction

    m_ppcQTTempCoeffY  [][]TLibCommon.TCoeff
    m_ppcQTTempCoeffCb [][]TLibCommon.TCoeff
    m_ppcQTTempCoeffCr [][]TLibCommon.TCoeff
    m_pcQTTempCoeffY   []TLibCommon.TCoeff
    m_pcQTTempCoeffCb  []TLibCommon.TCoeff
    m_pcQTTempCoeffCr  []TLibCommon.TCoeff
    //#if ADAPTIVE_QP_SELECTION
    m_ppcQTTempArlCoeffY  [][]TLibCommon.TCoeff //int
    m_ppcQTTempArlCoeffCb [][]TLibCommon.TCoeff //int
    m_ppcQTTempArlCoeffCr [][]TLibCommon.TCoeff //int
    m_pcQTTempArlCoeffY   []TLibCommon.TCoeff   //int
    m_pcQTTempArlCoeffCb  []TLibCommon.TCoeff   //int
    m_pcQTTempArlCoeffCr  []TLibCommon.TCoeff   //int
    //#endif
    m_puhQTTempTrIdx []byte
    m_puhQTTempCbf   [3][]byte

    m_pcQTTempTComYuv              []TLibCommon.TComYuv
    m_tmpYuvPred                   *TLibCommon.TComYuv // To be used in xGetInterPredictionError() to avoid constant memory allocation/deallocation
    m_pSharedPredTransformSkip     [3][]TLibCommon.Pel
    m_pcQTTempTUCoeffY             []TLibCommon.TCoeff
    m_pcQTTempTUCoeffCb            []TLibCommon.TCoeff
    m_pcQTTempTUCoeffCr            []TLibCommon.TCoeff
    m_puhQTTempTransformSkipFlag   [3][]bool
    m_pcQTTempTransformSkipTComYuv TLibCommon.TComYuv
    //#if ADAPTIVE_QP_SELECTION
    m_ppcQTTempTUArlCoeffY  []TLibCommon.TCoeff
    m_ppcQTTempTUArlCoeffCb []TLibCommon.TCoeff
    m_ppcQTTempTUArlCoeffCr []TLibCommon.TCoeff
    //#endif
    //protected:
    // interface to option
    m_pcEncCfg *TEncCfg

    // interface to classes
    m_pcTrQuant      *TLibCommon.TComTrQuant
    m_pcRdCost       *TEncRdCost
    m_pcEntropyCoder *TEncEntropy

    // ME parameters
    m_iSearchRange      int
    m_bipredSearchRange int // Search range for bi-prediction
    m_iFastSearch       int
    m_aaiAdaptSR        [2][33]int
    m_cSrchRngLT        TLibCommon.TComMv
    m_cSrchRngRB        TLibCommon.TComMv
    m_acMvPredictors    [3]TLibCommon.TComMv

    // RD computation
    m_pppcRDSbacCoder   [][]*TEncSbac
    m_pcRDGoOnSbacCoder *TEncSbac
    m_bUseSBACRD        bool
    m_cDistParam        DistParam

    // Misc.
    m_pTempPel    []TLibCommon.Pel
    m_puiDFilter  []uint
    m_iMaxDeltaQP int

    // AMVP cost computation
    // UInt            m_auiMVPIdxCost[AMVP_MAX_NUM_CANDS+1][AMVP_MAX_NUM_CANDS];
    m_auiMVPIdxCost [TLibCommon.AMVP_MAX_NUM_CANDS + 1][TLibCommon.AMVP_MAX_NUM_CANDS + 1]uint //th array bounds
}

func NewTEncSearch() *TEncSearch {
    this := &TEncSearch{}
    this.m_ppcQTTempCoeffY = nil
    this.m_ppcQTTempCoeffCb = nil
    this.m_ppcQTTempCoeffCr = nil
    this.m_pcQTTempCoeffY = nil
    this.m_pcQTTempCoeffCb = nil
    this.m_pcQTTempCoeffCr = nil
    //#if ADAPTIVE_QP_SELECTION
    this.m_ppcQTTempArlCoeffY = nil
    this.m_ppcQTTempArlCoeffCb = nil
    this.m_ppcQTTempArlCoeffCr = nil
    this.m_pcQTTempArlCoeffY = nil
    this.m_pcQTTempArlCoeffCb = nil
    this.m_pcQTTempArlCoeffCr = nil
    //#endif
    this.m_puhQTTempTrIdx = nil
    this.m_puhQTTempCbf[0] = nil
    this.m_puhQTTempCbf[1] = nil
    this.m_puhQTTempCbf[2] = nil
    this.m_pcQTTempTComYuv = nil
    this.m_pcEncCfg = nil
    this.m_pcEntropyCoder = nil
    this.m_pTempPel = nil
    this.m_pSharedPredTransformSkip[0] = nil
    this.m_pSharedPredTransformSkip[1] = nil
    this.m_pSharedPredTransformSkip[2] = nil
    this.m_pcQTTempTUCoeffY = nil
    this.m_pcQTTempTUCoeffCb = nil
    this.m_pcQTTempTUCoeffCr = nil
    //#if ADAPTIVE_QP_SELECTION
    this.m_ppcQTTempTUArlCoeffY = nil
    this.m_ppcQTTempTUArlCoeffCb = nil
    this.m_ppcQTTempTUArlCoeffCr = nil
    //#endif
    this.m_puhQTTempTransformSkipFlag[0] = nil
    this.m_puhQTTempTransformSkipFlag[1] = nil
    this.m_puhQTTempTransformSkipFlag[2] = nil

    this.setWpScalingDistParam(nil, -1, TLibCommon.REF_PIC_LIST_X)

    return this
}

func (this *TEncSearch) init(pcEncCfg *TEncCfg,
    pcTrQuant *TLibCommon.TComTrQuant,
    iSearchRange int,
    bipredSearchRange int,
    iFastSearch int,
    iMaxDeltaQP int,
    pcEntropyCoder *TEncEntropy,
    pcRdCost *TEncRdCost,
    pppcRDSbacCoder [][]*TEncSbac,
    pcRDGoOnSbacCoder *TEncSbac) {
    this.m_pcEncCfg = pcEncCfg
    this.m_pcTrQuant = pcTrQuant
    this.m_iSearchRange = iSearchRange
    this.m_bipredSearchRange = bipredSearchRange
    this.m_iFastSearch = iFastSearch
    this.m_iMaxDeltaQP = iMaxDeltaQP
    this.m_pcEntropyCoder = pcEntropyCoder
    this.m_pcRdCost = pcRdCost

    this.m_pppcRDSbacCoder = pppcRDSbacCoder
    this.m_pcRDGoOnSbacCoder = pcRDGoOnSbacCoder

    this.m_bUseSBACRD = pppcRDSbacCoder != nil // ? true : false;

    for iDir := 0; iDir < 2; iDir++ {
        for iRefIdx := 0; iRefIdx < 33; iRefIdx++ {
            this.m_aaiAdaptSR[iDir][iRefIdx] = iSearchRange
        }
    }

    this.m_puiDFilter = s_auiDFilter[4:]

    // initialize motion cost
    //#if !FIX203
    //  this.m_pcRdCost->initRateDistortionModel( this.m_iSearchRange << 2 );
    //#endif

    for iNum := 0; iNum < TLibCommon.AMVP_MAX_NUM_CANDS+1; iNum++ {
        for iIdx := 0; iIdx < TLibCommon.AMVP_MAX_NUM_CANDS; iIdx++ {
            if iIdx < iNum {
                this.m_auiMVPIdxCost[iIdx][iNum] = this.xGetMvpIdxBits(iIdx, iNum)
            } else {
                this.m_auiMVPIdxCost[iIdx][iNum] = uint(TLibCommon.MAX_INT)
            }
        }
    }

    this.InitTempBuff(pcEncCfg.GetMaxCUWidth(), pcEncCfg.GetMaxCUHeight())

    this.m_pTempPel = make([]TLibCommon.Pel, pcEncCfg.GetMaxCUWidth()*pcEncCfg.GetMaxCUHeight())

    uiNumLayersToAllocate := pcEncCfg.GetQuadtreeTULog2MaxSize() - pcEncCfg.GetQuadtreeTULog2MinSize() + 1
    this.m_ppcQTTempCoeffY = make([][]TLibCommon.TCoeff, uiNumLayersToAllocate)
    this.m_ppcQTTempCoeffCb = make([][]TLibCommon.TCoeff, uiNumLayersToAllocate)
    this.m_ppcQTTempCoeffCr = make([][]TLibCommon.TCoeff, uiNumLayersToAllocate)
    this.m_pcQTTempCoeffY  = make([]TLibCommon.TCoeff, pcEncCfg.GetMaxCUWidth()*pcEncCfg.GetMaxCUHeight())
    this.m_pcQTTempCoeffCb = make([]TLibCommon.TCoeff, pcEncCfg.GetMaxCUWidth()*pcEncCfg.GetMaxCUHeight()>>2)
    this.m_pcQTTempCoeffCr = make([]TLibCommon.TCoeff, pcEncCfg.GetMaxCUWidth()*pcEncCfg.GetMaxCUHeight()>>2)
    //#if ADAPTIVE_QP_SELECTION
    this.m_ppcQTTempArlCoeffY = make([][]TLibCommon.TCoeff, uiNumLayersToAllocate)
    this.m_ppcQTTempArlCoeffCb = make([][]TLibCommon.TCoeff, uiNumLayersToAllocate)
    this.m_ppcQTTempArlCoeffCr = make([][]TLibCommon.TCoeff, uiNumLayersToAllocate)
    this.m_pcQTTempArlCoeffY  = make([]TLibCommon.TCoeff, pcEncCfg.GetMaxCUWidth()*pcEncCfg.GetMaxCUHeight())
    this.m_pcQTTempArlCoeffCb = make([]TLibCommon.TCoeff, pcEncCfg.GetMaxCUWidth()*pcEncCfg.GetMaxCUHeight()>>2)
    this.m_pcQTTempArlCoeffCr = make([]TLibCommon.TCoeff, pcEncCfg.GetMaxCUWidth()*pcEncCfg.GetMaxCUHeight()>>2)
    //#endif

    uiNumPartitions := 1 << (pcEncCfg.GetMaxCUDepth() << 1)
    this.m_puhQTTempTrIdx = make([]byte, uiNumPartitions)
    this.m_puhQTTempCbf[0] = make([]byte, uiNumPartitions)
    this.m_puhQTTempCbf[1] = make([]byte, uiNumPartitions)
    this.m_puhQTTempCbf[2] = make([]byte, uiNumPartitions)
    this.m_pcQTTempTComYuv = make([]TLibCommon.TComYuv, uiNumLayersToAllocate)
    for ui := uint(0); ui < uiNumLayersToAllocate; ui++ {
        this.m_ppcQTTempCoeffY[ui] = make([]TLibCommon.TCoeff, pcEncCfg.GetMaxCUWidth()*pcEncCfg.GetMaxCUHeight())
        this.m_ppcQTTempCoeffCb[ui] = make([]TLibCommon.TCoeff, pcEncCfg.GetMaxCUWidth()*pcEncCfg.GetMaxCUHeight()>>2)
        this.m_ppcQTTempCoeffCr[ui] = make([]TLibCommon.TCoeff, pcEncCfg.GetMaxCUWidth()*pcEncCfg.GetMaxCUHeight()>>2)
        //#if ADAPTIVE_QP_SELECTION
        this.m_ppcQTTempArlCoeffY[ui] = make([]TLibCommon.TCoeff, pcEncCfg.GetMaxCUWidth()*pcEncCfg.GetMaxCUHeight())
        this.m_ppcQTTempArlCoeffCb[ui] = make([]TLibCommon.TCoeff, pcEncCfg.GetMaxCUWidth()*pcEncCfg.GetMaxCUHeight()>>2)
        this.m_ppcQTTempArlCoeffCr[ui] = make([]TLibCommon.TCoeff, pcEncCfg.GetMaxCUWidth()*pcEncCfg.GetMaxCUHeight()>>2)
        //#endif
        this.m_pcQTTempTComYuv[ui].Create(pcEncCfg.GetMaxCUWidth(), pcEncCfg.GetMaxCUHeight())
    }
    this.m_pSharedPredTransformSkip[0] = make([]TLibCommon.Pel, TLibCommon.MAX_TS_WIDTH*TLibCommon.MAX_TS_HEIGHT)
    this.m_pSharedPredTransformSkip[1] = make([]TLibCommon.Pel, TLibCommon.MAX_TS_WIDTH*TLibCommon.MAX_TS_HEIGHT)
    this.m_pSharedPredTransformSkip[2] = make([]TLibCommon.Pel, TLibCommon.MAX_TS_WIDTH*TLibCommon.MAX_TS_HEIGHT)
    this.m_pcQTTempTUCoeffY = make([]TLibCommon.TCoeff, TLibCommon.MAX_TS_WIDTH*TLibCommon.MAX_TS_HEIGHT)
    this.m_pcQTTempTUCoeffCb = make([]TLibCommon.TCoeff, TLibCommon.MAX_TS_WIDTH*TLibCommon.MAX_TS_HEIGHT)
    this.m_pcQTTempTUCoeffCr = make([]TLibCommon.TCoeff, TLibCommon.MAX_TS_WIDTH*TLibCommon.MAX_TS_HEIGHT)
    //#if ADAPTIVE_QP_SELECTION
    this.m_ppcQTTempTUArlCoeffY = make([]TLibCommon.TCoeff, TLibCommon.MAX_TS_WIDTH*TLibCommon.MAX_TS_HEIGHT)
    this.m_ppcQTTempTUArlCoeffCb = make([]TLibCommon.TCoeff, TLibCommon.MAX_TS_WIDTH*TLibCommon.MAX_TS_HEIGHT)
    this.m_ppcQTTempTUArlCoeffCr = make([]TLibCommon.TCoeff, TLibCommon.MAX_TS_WIDTH*TLibCommon.MAX_TS_HEIGHT)
    //#endif
    this.m_pcQTTempTransformSkipTComYuv.Create(pcEncCfg.GetMaxCUWidth(), pcEncCfg.GetMaxCUHeight())

    this.m_puhQTTempTransformSkipFlag[0] = make([]bool, uiNumPartitions)
    this.m_puhQTTempTransformSkipFlag[1] = make([]bool, uiNumPartitions)
    this.m_puhQTTempTransformSkipFlag[2] = make([]bool, uiNumPartitions)
    this.m_tmpYuvPred = TLibCommon.NewTComYuv()
    this.m_tmpYuvPred.Create(TLibCommon.MAX_CU_SIZE, TLibCommon.MAX_CU_SIZE)
}

// sub-functions for ME
func (this *TEncSearch) xTZSearchHelp(pcPatternKey *TLibCommon.TComPattern, rcStruct *IntTZSearchStruct, iSearchX, iSearchY int, ucPointNr byte, uiDistance uint ) {
    var uiSad uint

    var piRefSrch []TLibCommon.Pel

    piRefSrch = rcStruct.piRefY[rcStruct.iOffset + iSearchY*rcStruct.iYStride+iSearchX:]

    //-- jclee for using the SAD function pointer
    this.m_pcRdCost.setDistParam2(pcPatternKey, piRefSrch, rcStruct.iYStride, &this.m_cDistParam)

    // fast encoder decision: use subsampled SAD when rows > 8 for integer ME
    if this.m_pcEncCfg.GetUseFastEnc() {
        if this.m_cDistParam.iRows > 8 {
            this.m_cDistParam.iSubShift = 1
        }
    }

    this.setDistParamComp(0) // Y component

    // distortion
    this.m_cDistParam.bitDepth = TLibCommon.G_bitDepthY
    uiSad = this.m_cDistParam.DistFunc(&this.m_cDistParam)

    // motion cost
    uiSad += this.m_pcRdCost.getCost2(iSearchX, iSearchY)

    if uiSad < rcStruct.uiBestSad {
        rcStruct.uiBestSad = uiSad
        rcStruct.iBestX = iSearchX
        rcStruct.iBestY = iSearchY
        rcStruct.uiBestDistance = uiDistance
        rcStruct.uiBestRound = 0
        rcStruct.ucPointNr = ucPointNr
    }
}
func (this *TEncSearch) xTZ2PointSearch(pcPatternKey *TLibCommon.TComPattern, rcStruct *IntTZSearchStruct, pcMvSrchRngLT *TLibCommon.TComMv, pcMvSrchRngRB *TLibCommon.TComMv) {
    iSrchRngHorLeft := int(pcMvSrchRngLT.GetHor())
    iSrchRngHorRight := int(pcMvSrchRngRB.GetHor())
    iSrchRngVerTop := int(pcMvSrchRngLT.GetVer())
    iSrchRngVerBottom := int(pcMvSrchRngRB.GetVer())

    // 2 point search,                   //   1 2 3
    // check only the 2 untested points  //   4 0 5
    // around the start point            //   6 7 8
    iStartX := rcStruct.iBestX
    iStartY := rcStruct.iBestY
    switch rcStruct.ucPointNr {
    case 1:
        if (iStartX - 1) >= iSrchRngHorLeft {
            this.xTZSearchHelp(pcPatternKey, rcStruct, iStartX-1, iStartY, 0, 2)
        }
        if (iStartY - 1) >= iSrchRngVerTop {
            this.xTZSearchHelp(pcPatternKey, rcStruct, iStartX, iStartY-1, 0, 2)
        }
    case 2:
        if (iStartY - 1) >= iSrchRngVerTop {
            if (iStartX - 1) >= iSrchRngHorLeft {
                this.xTZSearchHelp(pcPatternKey, rcStruct, iStartX-1, iStartY-1, 0, 2)
            }
            if (iStartX + 1) <= iSrchRngHorRight {
                this.xTZSearchHelp(pcPatternKey, rcStruct, iStartX+1, iStartY-1, 0, 2)
            }
        }
    case 3:
        if (iStartY - 1) >= iSrchRngVerTop {
            this.xTZSearchHelp(pcPatternKey, rcStruct, iStartX, iStartY-1, 0, 2)
        }
        if (iStartX + 1) <= iSrchRngHorRight {
            this.xTZSearchHelp(pcPatternKey, rcStruct, iStartX+1, iStartY, 0, 2)
        }
    case 4:
        if (iStartX - 1) >= iSrchRngHorLeft {
            if (iStartY + 1) <= iSrchRngVerBottom {
                this.xTZSearchHelp(pcPatternKey, rcStruct, iStartX-1, iStartY+1, 0, 2)
            }
            if (iStartY - 1) >= iSrchRngVerTop {
                this.xTZSearchHelp(pcPatternKey, rcStruct, iStartX-1, iStartY-1, 0, 2)
            }
        }
    case 5:
        if (iStartX + 1) <= iSrchRngHorRight {
            if (iStartY - 1) >= iSrchRngVerTop {
                this.xTZSearchHelp(pcPatternKey, rcStruct, iStartX+1, iStartY-1, 0, 2)
            }
            if (iStartY + 1) <= iSrchRngVerBottom {
                this.xTZSearchHelp(pcPatternKey, rcStruct, iStartX+1, iStartY+1, 0, 2)
            }
        }
    case 6:
        if (iStartX - 1) >= iSrchRngHorLeft {
            this.xTZSearchHelp(pcPatternKey, rcStruct, iStartX-1, iStartY, 0, 2)
        }
        if (iStartY + 1) <= iSrchRngVerBottom {
            this.xTZSearchHelp(pcPatternKey, rcStruct, iStartX, iStartY+1, 0, 2)
        }
    case 7:
        if (iStartY + 1) <= iSrchRngVerBottom {
            if (iStartX - 1) >= iSrchRngHorLeft {
                this.xTZSearchHelp(pcPatternKey, rcStruct, iStartX-1, iStartY+1, 0, 2)
            }
            if (iStartX + 1) <= iSrchRngHorRight {
                this.xTZSearchHelp(pcPatternKey, rcStruct, iStartX+1, iStartY+1, 0, 2)
            }
        }
    case 8:
        if (iStartX + 1) <= iSrchRngHorRight {
            this.xTZSearchHelp(pcPatternKey, rcStruct, iStartX+1, iStartY, 0, 2)
        }
        if (iStartY + 1) <= iSrchRngVerBottom {
            this.xTZSearchHelp(pcPatternKey, rcStruct, iStartX, iStartY+1, 0, 2)
        }
    default:
        {
        }
    }   // switch( rcStruct.ucPointNr )
}
func (this *TEncSearch) xTZ8PointSquareSearch(pcPatternKey *TLibCommon.TComPattern, rcStruct *IntTZSearchStruct, pcMvSrchRngLT *TLibCommon.TComMv, pcMvSrchRngRB *TLibCommon.TComMv, iStartX, iStartY, iDist int) {
    iSrchRngHorLeft := int(pcMvSrchRngLT.GetHor())
    iSrchRngHorRight := int(pcMvSrchRngRB.GetHor())
    iSrchRngVerTop := int(pcMvSrchRngLT.GetVer())
    iSrchRngVerBottom := int(pcMvSrchRngRB.GetVer())

    // 8 point search,                   //   1 2 3
    // search around the start point     //   4 0 5
    // with the required  distance       //   6 7 8
    //assert( iDist != 0 );
    iTop := iStartY - iDist
    iBottom := iStartY + iDist
    iLeft := iStartX - iDist
    iRight := iStartX + iDist
    rcStruct.uiBestRound += 1

    if iTop >= iSrchRngVerTop { // check top
        if iLeft >= iSrchRngHorLeft { // check top left
            this.xTZSearchHelp(pcPatternKey, rcStruct, iLeft, iTop, 1, uint(iDist))
        }
        // top middle
        this.xTZSearchHelp(pcPatternKey, rcStruct, iStartX, iTop, 2, uint(iDist))

        if iRight <= iSrchRngHorRight { // check top right
            this.xTZSearchHelp(pcPatternKey, rcStruct, iRight, iTop, 3, uint(iDist))
        }
    }   // check top
    if iLeft >= iSrchRngHorLeft { // check middle left
        this.xTZSearchHelp(pcPatternKey, rcStruct, iLeft, iStartY, 4, uint(iDist))
    }
    if iRight <= iSrchRngHorRight { // check middle right
        this.xTZSearchHelp(pcPatternKey, rcStruct, iRight, iStartY, 5, uint(iDist))
    }
    if iBottom <= iSrchRngVerBottom { // check bottom
        if iLeft >= iSrchRngHorLeft { // check bottom left
            this.xTZSearchHelp(pcPatternKey, rcStruct, iLeft, iBottom, 6, uint(iDist))
        }
        // check bottom middle
        this.xTZSearchHelp(pcPatternKey, rcStruct, iStartX, iBottom, 7, uint(iDist))

        if iRight <= iSrchRngHorRight { // check bottom right
            this.xTZSearchHelp(pcPatternKey, rcStruct, iRight, iBottom, 8, uint(iDist))
        }
    }   // check bottom
}
func (this *TEncSearch) xTZ8PointDiamondSearch(pcPatternKey *TLibCommon.TComPattern, rcStruct *IntTZSearchStruct, pcMvSrchRngLT *TLibCommon.TComMv, pcMvSrchRngRB *TLibCommon.TComMv, iStartX, iStartY, iDist int) {
    iSrchRngHorLeft := int(pcMvSrchRngLT.GetHor())
    iSrchRngHorRight := int(pcMvSrchRngRB.GetHor())
    iSrchRngVerTop := int(pcMvSrchRngLT.GetVer())
    iSrchRngVerBottom := int(pcMvSrchRngRB.GetVer())

    // 8 point search,                   //   1 2 3
    // search around the start point     //   4 0 5
    // with the required  distance       //   6 7 8
    //assert ( iDist != 0 );
    iTop := iStartY - iDist
    iBottom := iStartY + iDist
    iLeft := iStartX - iDist
    iRight := iStartX + iDist
    rcStruct.uiBestRound += 1

    if iDist == 1 { // iDist == 1
        if iTop >= iSrchRngVerTop { // check top
            this.xTZSearchHelp(pcPatternKey, rcStruct, iStartX, iTop, 2, uint(iDist))
        }
        if iLeft >= iSrchRngHorLeft { // check middle left
            this.xTZSearchHelp(pcPatternKey, rcStruct, iLeft, iStartY, 4, uint(iDist))
        }
        if iRight <= iSrchRngHorRight { // check middle right
            this.xTZSearchHelp(pcPatternKey, rcStruct, iRight, iStartY, 5, uint(iDist))
        }
        if iBottom <= iSrchRngVerBottom { // check bottom
            this.xTZSearchHelp(pcPatternKey, rcStruct, iStartX, iBottom, 7, uint(iDist))
        }
    } else { // if (uint(iDist) != 1)
        if iDist <= 8 {
            iTop_2 := iStartY - (iDist >> 1)
            iBottom_2 := iStartY + (iDist >> 1)
            iLeft_2 := iStartX - (iDist >> 1)
            iRight_2 := iStartX + (iDist >> 1)

            if iTop >= iSrchRngVerTop && iLeft >= iSrchRngHorLeft &&
                iRight <= iSrchRngHorRight && iBottom <= iSrchRngVerBottom { // check border
                this.xTZSearchHelp(pcPatternKey, rcStruct, iStartX, iTop, 2, uint(iDist))
                this.xTZSearchHelp(pcPatternKey, rcStruct, iLeft_2, iTop_2, 1, uint(iDist)>>1)
                this.xTZSearchHelp(pcPatternKey, rcStruct, iRight_2, iTop_2, 3, uint(iDist)>>1)
                this.xTZSearchHelp(pcPatternKey, rcStruct, iLeft, iStartY, 4, uint(iDist))
                this.xTZSearchHelp(pcPatternKey, rcStruct, iRight, iStartY, 5, uint(iDist))
                this.xTZSearchHelp(pcPatternKey, rcStruct, iLeft_2, iBottom_2, 6, uint(iDist)>>1)
                this.xTZSearchHelp(pcPatternKey, rcStruct, iRight_2, iBottom_2, 8, uint(iDist)>>1)
                this.xTZSearchHelp(pcPatternKey, rcStruct, iStartX, iBottom, 7, uint(iDist))
            } else { // check border
                if iTop >= iSrchRngVerTop { // check top
                    this.xTZSearchHelp(pcPatternKey, rcStruct, iStartX, iTop, 2, uint(iDist))
                }
                if iTop_2 >= iSrchRngVerTop { // check half top
                    if iLeft_2 >= iSrchRngHorLeft { // check half left
                        this.xTZSearchHelp(pcPatternKey, rcStruct, iLeft_2, iTop_2, 1, (uint(iDist) >> 1))
                    }
                    if iRight_2 <= iSrchRngHorRight { // check half right
                        this.xTZSearchHelp(pcPatternKey, rcStruct, iRight_2, iTop_2, 3, (uint(iDist) >> 1))
                    }
                }   // check half top
                if iLeft >= iSrchRngHorLeft { // check left
                    this.xTZSearchHelp(pcPatternKey, rcStruct, iLeft, iStartY, 4, uint(iDist))
                }
                if iRight <= iSrchRngHorRight { // check right
                    this.xTZSearchHelp(pcPatternKey, rcStruct, iRight, iStartY, 5, uint(iDist))
                }
                if iBottom_2 <= iSrchRngVerBottom { // check half bottom
                    if iLeft_2 >= iSrchRngHorLeft { // check half left
                        this.xTZSearchHelp(pcPatternKey, rcStruct, iLeft_2, iBottom_2, 6, (uint(iDist) >> 1))
                    }
                    if iRight_2 <= iSrchRngHorRight { // check half right
                        this.xTZSearchHelp(pcPatternKey, rcStruct, iRight_2, iBottom_2, 8, (uint(iDist) >> 1))
                    }
                }   // check half bottom
                if iBottom <= iSrchRngVerBottom { // check bottom
                    this.xTZSearchHelp(pcPatternKey, rcStruct, iStartX, iBottom, 7, uint(iDist))
                }
            }   // check border
        } else { // uint(iDist) > 8
            if iTop >= iSrchRngVerTop && iLeft >= iSrchRngHorLeft &&
                iRight <= iSrchRngHorRight && iBottom <= iSrchRngVerBottom { // check border
                this.xTZSearchHelp(pcPatternKey, rcStruct, iStartX, iTop, 0, uint(iDist))
                this.xTZSearchHelp(pcPatternKey, rcStruct, iLeft, iStartY, 0, uint(iDist))
                this.xTZSearchHelp(pcPatternKey, rcStruct, iRight, iStartY, 0, uint(iDist))
                this.xTZSearchHelp(pcPatternKey, rcStruct, iStartX, iBottom, 0, uint(iDist))
                for index := 1; index < 4; index++ {
                    iPosYT := iTop + ((iDist >> 2) * index)
                    iPosYB := iBottom - ((iDist >> 2) * index)
                    iPosXL := iStartX - ((iDist >> 2) * index)
                    iPosXR := iStartX + ((iDist >> 2) * index)
                    this.xTZSearchHelp(pcPatternKey, rcStruct, iPosXL, iPosYT, 0, uint(iDist))
                    this.xTZSearchHelp(pcPatternKey, rcStruct, iPosXR, iPosYT, 0, uint(iDist))
                    this.xTZSearchHelp(pcPatternKey, rcStruct, iPosXL, iPosYB, 0, uint(iDist))
                    this.xTZSearchHelp(pcPatternKey, rcStruct, iPosXR, iPosYB, 0, uint(iDist))
                }
            } else { // check border
                if iTop >= iSrchRngVerTop { // check top
                    this.xTZSearchHelp(pcPatternKey, rcStruct, iStartX, iTop, 0, uint(iDist))
                }
                if iLeft >= iSrchRngHorLeft { // check left
                    this.xTZSearchHelp(pcPatternKey, rcStruct, iLeft, iStartY, 0, uint(iDist))
                }
                if iRight <= iSrchRngHorRight { // check right
                    this.xTZSearchHelp(pcPatternKey, rcStruct, iRight, iStartY, 0, uint(iDist))
                }
                if iBottom <= iSrchRngVerBottom { // check bottom
                    this.xTZSearchHelp(pcPatternKey, rcStruct, iStartX, iBottom, 0, uint(iDist))
                }
                for index := 1; index < 4; index++ {
                    iPosYT := iTop + ((iDist >> 2) * index)
                    iPosYB := iBottom - ((iDist >> 2) * index)
                    iPosXL := iStartX - ((iDist >> 2) * index)
                    iPosXR := iStartX + ((iDist >> 2) * index)

                    if iPosYT >= iSrchRngVerTop { // check top
                        if iPosXL >= iSrchRngHorLeft { // check left
                            this.xTZSearchHelp(pcPatternKey, rcStruct, iPosXL, iPosYT, 0, uint(iDist))
                        }
                        if iPosXR <= iSrchRngHorRight { // check right
                            this.xTZSearchHelp(pcPatternKey, rcStruct, iPosXR, iPosYT, 0, uint(iDist))
                        }
                    }   // check top
                    if iPosYB <= iSrchRngVerBottom { // check bottom
                        if iPosXL >= iSrchRngHorLeft { // check left
                            this.xTZSearchHelp(pcPatternKey, rcStruct, iPosXL, iPosYB, 0, uint(iDist))
                        }
                        if iPosXR <= iSrchRngHorRight { // check right
                            this.xTZSearchHelp(pcPatternKey, rcStruct, iPosXR, iPosYB, 0, uint(iDist))
                        }
                    }   // check bottom
                }   // for ...
            }   // check border
        }   // iDist <= 8
    }   // iDist == 1
}

/// sub-function for motion vector refinement used in fractional-pel accuracy
func (this *TEncSearch) xPatternRefinement(pcPatternKey *TLibCommon.TComPattern,
    baseRefMv TLibCommon.TComMv,
    iFrac int, rcMvFrac *TLibCommon.TComMv) uint {
    var uiDist uint
    uiDistBest := uint(TLibCommon.MAX_UINT)
    uiDirecBest := uint(0)

    var piRefPos []TLibCommon.Pel
    iRefStride := int(this.GetFilteredBlock(0, 0).GetStride())
    //#if NS_HAD
    //  this.m_pcRdCost.setDistParam3( pcPatternKey, this.GetFilteredBlock(0,0).GetLumaAddr(), iRefStride, 1, &this.m_cDistParam, this.m_pcEncCfg.GetUseHADME(), this.m_pcEncCfg.GetUseNSQT() );
    //#else
    this.m_pcRdCost.setDistParam3(pcPatternKey, this.GetFilteredBlock(0, 0).GetLumaAddr(), iRefStride, 1, &this.m_cDistParam, this.m_pcEncCfg.GetUseHADME())
    //#endif

    var pcMvRefine []*TLibCommon.TComMv
    var cMvTest TLibCommon.TComMv
    if iFrac == 2 {
        pcMvRefine = s_acMvRefineH[:]
    } else {
        pcMvRefine = s_acMvRefineQ[:]
    }

    for i := uint(0); i < 9; i++ {
        cMvTest = *pcMvRefine[i]
        cMvTest.AddMv(baseRefMv)

        horVal := int(cMvTest.GetHor()) * iFrac
        verVal := int(cMvTest.GetVer()) * iFrac
        piRefPos = this.GetFilteredBlock(verVal&3, horVal&3).GetLumaAddr()
        if horVal == 2 && (verVal&1) == 0 {
            piRefPos = piRefPos[1:]
        }
        if (horVal&1) == 0 && verVal == 2 {
            piRefPos = piRefPos[iRefStride:]
        }
        cMvTest = *pcMvRefine[i]
        cMvTest.AddMv(*rcMvFrac)

        this.setDistParamComp(0) // Y component

        this.m_cDistParam.pCur = piRefPos
        this.m_cDistParam.bitDepth = TLibCommon.G_bitDepthY
        uiDist = this.m_cDistParam.DistFunc(&this.m_cDistParam)
        uiDist += this.m_pcRdCost.getCost2(int(cMvTest.GetHor()), int(cMvTest.GetVer()))
		//fmt.Printf("pCur[0]=%d, pOrg[0]=%d\n",this.m_cDistParam.pCur[0],this.m_cDistParam.pOrg[0]);
        if uiDist < uiDistBest {
            uiDistBest = uiDist
            uiDirecBest = i
        }
    }

    *rcMvFrac = *pcMvRefine[uiDirecBest]

    return uiDistBest
}

func (this *TEncSearch) xGetInterPredictionError(pcCU *TLibCommon.TComDataCU, pcYuvOrg *TLibCommon.TComYuv, iPartIdx int, ruiErr *uint, Hadamard bool) {
    this.MotionCompensation(pcCU, this.m_tmpYuvPred, TLibCommon.REF_PIC_LIST_X, iPartIdx)

    uiAbsPartIdx := uint(0)
    iWidth := 0
    iHeight := 0
    pcCU.GetPartIndexAndSize(uint(iPartIdx), &uiAbsPartIdx, &iWidth, &iHeight)

    var cDistParam DistParam

    cDistParam.bApplyWeight = false

    this.m_pcRdCost.setDistParam4(&cDistParam, TLibCommon.G_bitDepthY,
        pcYuvOrg.GetLumaAddr1(uiAbsPartIdx), int(pcYuvOrg.GetStride()),
        this.m_tmpYuvPred.GetLumaAddr1(uiAbsPartIdx), int(this.m_tmpYuvPred.GetStride()),
        //#if NS_HAD
        //                            iWidth, iHeight, this.m_pcEncCfg.GetUseHADME(), this.m_pcEncCfg.GetUseNSQT() );
        //#else
        iWidth, iHeight, this.m_pcEncCfg.GetUseHADME())
    //#endif
    *ruiErr = cDistParam.DistFunc(&cDistParam)
}

func (this *TEncSearch) preestChromaPredMode(pcCU *TLibCommon.TComDataCU,
    pcOrgYuv *TLibCommon.TComYuv,
    pcPredYuv *TLibCommon.TComYuv,
    uiChromaId uint) {
    uiWidth := uint(pcCU.GetWidth1(0)) >> 1
    uiHeight := uint(pcCU.GetHeight1(0)) >> 1
    uiStride := pcOrgYuv.GetCStride()
    piOrgU := pcOrgYuv.GetCbAddr1(0)
    piOrgV := pcOrgYuv.GetCrAddr1(0)
    piPredU := pcPredYuv.GetCbAddr1(0)
    piPredV := pcPredYuv.GetCrAddr1(0)

    //===== init pattern =====
    bAboveAvail := false
    bLeftAvail := false
    pcCU.GetPattern().InitPattern3(pcCU, 0, 0)
    pcCU.GetPattern().InitAdiPatternChroma(pcCU, 0, 0, this.GetYuvExt(), this.GetYuvExtStride(), this.GetYuvExtHeight(), &bAboveAvail, &bLeftAvail, uiChromaId)
    pPatChromaU := pcCU.GetPattern().GetAdiCbBuf(int(uiWidth), int(uiHeight), this.GetYuvExt())
    pPatChromaV := pcCU.GetPattern().GetAdiCrBuf(int(uiWidth), int(uiHeight), this.GetYuvExt())

    //===== get best prediction modes (using SAD) =====
    uiMinMode := uint(0)
    uiMaxMode := uint(4)
    uiBestMode := uint(TLibCommon.MAX_UINT)
    uiMinSAD := uint(TLibCommon.MAX_UINT)
    for uiMode := uiMinMode; uiMode < uiMaxMode; uiMode++ {
        //--- get prediction ---
        this.PredIntraChromaAng(pPatChromaU, uiMode, piPredU, uiStride, int(uiWidth), int(uiHeight), bAboveAvail, bLeftAvail)
        this.PredIntraChromaAng(pPatChromaV, uiMode, piPredV, uiStride, int(uiWidth), int(uiHeight), bAboveAvail, bLeftAvail)

        //--- get SAD ---
        uiSAD := this.m_pcRdCost.calcHAD(TLibCommon.G_bitDepthC, piOrgU, int(uiStride), piPredU, int(uiStride), int(uiWidth), int(uiHeight))
        uiSAD += this.m_pcRdCost.calcHAD(TLibCommon.G_bitDepthC, piOrgV, int(uiStride), piPredV, int(uiStride), int(uiWidth), int(uiHeight))
        //--- check ---
        if uiSAD < uiMinSAD {
            uiMinSAD = uiSAD
            uiBestMode = uiMode
        }
    }

    //===== set chroma pred mode =====
    pcCU.SetChromIntraDirSubParts(uiBestMode, 0, uint(pcCU.GetDepth1(0)))
}

func (this *TEncSearch) estIntraPredQT(pcCU *TLibCommon.TComDataCU,
    pcOrgYuv *TLibCommon.TComYuv,
    pcPredYuv *TLibCommon.TComYuv,
    pcResiYuv *TLibCommon.TComYuv,
    pcRecoYuv *TLibCommon.TComYuv,
    ruiDistC *uint,
    bLumaOnly bool) {
    uiDepth := uint(pcCU.GetDepth1(0))
    uiNumPU := uint(pcCU.GetNumPartInter())
    var uiInitTrDepth uint
    if pcCU.GetPartitionSize1(0) == TLibCommon.SIZE_2Nx2N {
        uiInitTrDepth = 0
    } else {
        uiInitTrDepth = 1
    }
    uiWidth := uint(pcCU.GetWidth1(0)) >> uiInitTrDepth
    uiHeight := uint(pcCU.GetHeight1(0)) >> uiInitTrDepth
    uiQNumParts := pcCU.GetTotalNumPart() >> 2
    uiWidthBit := pcCU.GetIntraSizeIdx(0)
    uiOverallDistY := uint(0)
    uiOverallDistC := uint(0)
    var CandNum uint
    var CandCostList [TLibCommon.FAST_UDI_MAX_RDMODE_NUM]float64

	//fmt.Printf("Enter estIntraPredQT\n")
	
    //===== set QP and clear Cbf =====
    if pcCU.GetSlice().GetPPS().GetUseDQP() {
        pcCU.SetQPSubParts(int(pcCU.GetQP1(0)), 0, uiDepth)
    } else {
        pcCU.SetQPSubParts(pcCU.GetSlice().GetSliceQp(), 0, uiDepth)
    }

    //===== loop over partitions =====
    uiPartOffset := uint(0)
    for uiPU := uint(0); uiPU < uiNumPU; uiPU++ {
        //===== init pattern for luma prediction =====
        bAboveAvail := false
        bLeftAvail := false
        pcCU.GetPattern().InitPattern3(pcCU, uiInitTrDepth, uiPartOffset)
        pcCU.GetPattern().InitAdiPattern(pcCU, uiPartOffset, uiInitTrDepth, this.GetYuvExt(), this.GetYuvExtStride(), this.GetYuvExtHeight(), &bAboveAvail, &bLeftAvail, false)

        //===== determine set of modes to be tested (using prediction signal only) =====
        numModesAvailable := 35 //total number of Intra modes
        piOrg := pcOrgYuv.GetLumaAddr2(uiPU, uiWidth)
        piPred := pcPredYuv.GetLumaAddr2(uiPU, uiWidth)
        uiStride := pcPredYuv.GetStride()
        var uiRdModeList [TLibCommon.FAST_UDI_MAX_RDMODE_NUM]uint
        numModesForFullRD := int(TLibCommon.G_aucIntraModeNumFast[uiWidthBit])

        doFastSearch := (numModesForFullRD != numModesAvailable)
        if doFastSearch {
            //assert(numModesForFullRD < numModesAvailable);

            for i := 0; i < numModesForFullRD; i++ {
                CandCostList[i] = TLibCommon.MAX_DOUBLE
            }
            CandNum = 0

            for modeIdx := 0; modeIdx < numModesAvailable; modeIdx++ {
                uiMode := uint(modeIdx)

                this.PredIntraLumaAng(pcCU.GetPattern(), uiMode, piPred, uiStride, int(uiWidth), int(uiHeight), bAboveAvail, bLeftAvail)
				//fmt.Printf("piOrg=%d, piPred=%d\n", piOrg[0], piPred[0]);
				
                // use hadamard transform here
                uiSad := this.m_pcRdCost.calcHAD(TLibCommon.G_bitDepthY, piOrg, int(uiStride), piPred, int(uiStride), int(uiWidth), int(uiHeight))
				//fmt.Printf("uiSad=%d\n", uiSad);
				
                iModeBits := this.xModeBitsIntra(pcCU, uiMode, uiPU, uiPartOffset, uiDepth, uiInitTrDepth)
                cost := float64(uiSad) + float64(iModeBits)*this.m_pcRdCost.getSqrtLambda()

                CandNum += this.xUpdateCandList(uiMode, cost, uint(numModesForFullRD), uiRdModeList[:], CandCostList[:])
            }

            //#if FAST_UDI_USE_MPM
            var uiPreds = [3]int{-1, -1, -1}
            iMode := -1
            numCand := pcCU.GetIntraDirLumaPredictor(uiPartOffset, uiPreds[:], &iMode)
            if iMode >= 0 {
                numCand = iMode
            }

            for j := 0; j < numCand; j++ {
                mostProbableModeIncluded := false
                mostProbableMode := uiPreds[j]

                for i := 0; i < numModesForFullRD; i++ {
                    mostProbableModeIncluded = mostProbableModeIncluded || (mostProbableMode == int(uiRdModeList[i]))
                }
                if !mostProbableModeIncluded {        
                    uiRdModeList[numModesForFullRD] = uint(mostProbableMode)
                    numModesForFullRD++
                    
                    //fmt.Printf("hit !mostProbableModeIncluded %d =%d\n",TLibCommon.B2U(mostProbableModeIncluded), numModesForFullRD);
                }
            }
            //#endif // FAST_UDI_USE_MPM
        } else {
            for i := 0; i < numModesForFullRD; i++ {
                uiRdModeList[i] = uint(i)
            }
        }

        //===== check modes (using r-d costs) =====
        //#if HHI_RQT_INTRA_SPEEDUP_MOD
        //uiSecondBestMode := TLibCommon.MAX_UINT
        //dSecondBestPUCost := TLibCommon.MAX_DOUBLE
        //#endif

        uiBestPUMode := uint(0)
        uiBestPUDistY := uint(0)
        uiBestPUDistC := uint(0)
        dBestPUCost := float64(TLibCommon.MAX_DOUBLE)
        for uiMode := 0; uiMode < numModesForFullRD; uiMode++ {
        	//fmt.Printf("uiMode=%d\n", uiMode);
            // set luma prediction mode
            uiOrgMode := uiRdModeList[uiMode]

            pcCU.SetLumaIntraDirSubParts(uiOrgMode, uiPartOffset, uiDepth+uiInitTrDepth)

            // set context models
            if this.m_bUseSBACRD {
                this.m_pcRDGoOnSbacCoder.load(this.m_pppcRDSbacCoder[uiDepth][TLibCommon.CI_CURR_BEST])
            }

            // determine residual for partition
            uiPUDistY := uint(0)
            uiPUDistC := uint(0)
            dPUCost := float64(0.0)
            //#if HHI_RQT_INTRA_SPEEDUP
            this.xRecurIntraCodingQT(pcCU, uiInitTrDepth, uiPartOffset, bLumaOnly, pcOrgYuv, pcPredYuv, pcResiYuv, &uiPUDistY, &uiPUDistC, true, &dPUCost)
            //#else
            //      xRecurIntraCodingQT( pcCU, uiInitTrDepth, uiPartOffset, bLumaOnly, pcOrgYuv, pcPredYuv, pcResiYuv, uiPUDistY, uiPUDistC, dPUCost );
            //#endif
			//fmt.Printf("uiOrgMode=%d, uiBestPUMode=%d, dBestPUCost=%f, dPUCost=%f\n",uiOrgMode, uiBestPUMode, dBestPUCost, dPUCost);
            // check r-d cost
            if dPUCost < dBestPUCost {
                //#if HHI_RQT_INTRA_SPEEDUP_MOD
                //uiSecondBestMode = uiBestPUMode
                //dSecondBestPUCost = dBestPUCost
                //#endif
                uiBestPUMode = uiOrgMode
                uiBestPUDistY = uiPUDistY
                uiBestPUDistC = uiPUDistC
                dBestPUCost = dPUCost

                this.xSetIntraResultQT(pcCU, uiInitTrDepth, uiPartOffset, bLumaOnly, pcRecoYuv)

                uiQPartNum := pcCU.GetPic().GetNumPartInCU() >> ((uint(pcCU.GetDepth1(0)) + uiInitTrDepth) << 1)

                for i := uint(0); i < uiQPartNum; i++ {
                    this.m_puhQTTempTrIdx[i] = pcCU.GetTransformIdx()[i+uiPartOffset]                                          //, uiQPartNum * sizeof( byte ) );
                    this.m_puhQTTempCbf[0][i] = pcCU.GetCbf1(TLibCommon.TEXT_LUMA)[i+uiPartOffset]                             //, uiQPartNum * sizeof( byte ) );
                    this.m_puhQTTempCbf[1][i] = pcCU.GetCbf1(TLibCommon.TEXT_CHROMA_U)[i+uiPartOffset]                         //, uiQPartNum * sizeof( byte ) );
                    this.m_puhQTTempCbf[2][i] = pcCU.GetCbf1(TLibCommon.TEXT_CHROMA_V)[i+uiPartOffset]                         //, uiQPartNum * sizeof( byte ) );
                    this.m_puhQTTempTransformSkipFlag[0][i] = pcCU.GetTransformSkip1(TLibCommon.TEXT_LUMA)[i+uiPartOffset]     //, uiQPartNum * sizeof( byte ) );
                    this.m_puhQTTempTransformSkipFlag[1][i] = pcCU.GetTransformSkip1(TLibCommon.TEXT_CHROMA_U)[i+uiPartOffset] //, uiQPartNum * sizeof( byte ) );
                    this.m_puhQTTempTransformSkipFlag[2][i] = pcCU.GetTransformSkip1(TLibCommon.TEXT_CHROMA_V)[i+uiPartOffset] //, uiQPartNum * sizeof( byte ) );
                }
            //#if HHI_RQT_INTRA_SPEEDUP_MOD
            //} else if dPUCost < dSecondBestPUCost {
            //    uiSecondBestMode = uiOrgMode
            //    dSecondBestPUCost = dPUCost
            }
            //#endif
        }   // Mode loop

        //#if HHI_RQT_INTRA_SPEEDUP
        //#if HHI_RQT_INTRA_SPEEDUP_MOD
        //for ui := 0; ui < 2; ui++ {
        //#endif
        for ui := 0; ui < 1; ui++ {
            //#if HHI_RQT_INTRA_SPEEDUP_MOD
            /*var uiOrgMode uint
            if ui != 0 {
                uiOrgMode = uiSecondBestMode
            } else {
                uiOrgMode = uiBestPUMode
            }
            if uiOrgMode == TLibCommon.MAX_UINT {
                break
            }
            #else*/
            uiOrgMode := uiBestPUMode
            //#endif
            pcCU.SetLumaIntraDirSubParts(uiOrgMode, uiPartOffset, uiDepth+uiInitTrDepth)

            // set context models
            if this.m_bUseSBACRD {
                this.m_pcRDGoOnSbacCoder.load(this.m_pppcRDSbacCoder[uiDepth][TLibCommon.CI_CURR_BEST])
            }

            // determine residual for partition
            uiPUDistY := uint(0)
            uiPUDistC := uint(0)
            dPUCost := float64(0.0)
            //fmt.Printf("uiBestPUMode=%d\n", uiBestPUMode);
            this.xRecurIntraCodingQT(pcCU, uiInitTrDepth, uiPartOffset, bLumaOnly, pcOrgYuv, pcPredYuv, pcResiYuv, &uiPUDistY, &uiPUDistC, false, &dPUCost)

            // check r-d cost
            if dPUCost < dBestPUCost {
                uiBestPUMode = uiOrgMode
                uiBestPUDistY = uiPUDistY
                uiBestPUDistC = uiPUDistC
                dBestPUCost = dPUCost

                this.xSetIntraResultQT(pcCU, uiInitTrDepth, uiPartOffset, bLumaOnly, pcRecoYuv)

                uiQPartNum := pcCU.GetPic().GetNumPartInCU() >> ((uint(pcCU.GetDepth1(0)) + uiInitTrDepth) << 1)
                for i := uint(0); i < uiQPartNum; i++ {
                    this.m_puhQTTempTrIdx[i] = pcCU.GetTransformIdx()[i+uiPartOffset]                                          //, uiQPartNum * sizeof( byte ) );
                    this.m_puhQTTempCbf[0][i] = pcCU.GetCbf1(TLibCommon.TEXT_LUMA)[i+uiPartOffset]                             //, uiQPartNum * sizeof( byte ) );
                    this.m_puhQTTempCbf[1][i] = pcCU.GetCbf1(TLibCommon.TEXT_CHROMA_U)[i+uiPartOffset]                         //, uiQPartNum * sizeof( byte ) );
                    this.m_puhQTTempCbf[2][i] = pcCU.GetCbf1(TLibCommon.TEXT_CHROMA_V)[i+uiPartOffset]                         //, uiQPartNum * sizeof( byte ) );
                    this.m_puhQTTempTransformSkipFlag[0][i] = pcCU.GetTransformSkip1(TLibCommon.TEXT_LUMA)[i+uiPartOffset]     //, uiQPartNum * sizeof( byte ) );
                    this.m_puhQTTempTransformSkipFlag[1][i] = pcCU.GetTransformSkip1(TLibCommon.TEXT_CHROMA_U)[i+uiPartOffset] //, uiQPartNum * sizeof( byte ) );
                    this.m_puhQTTempTransformSkipFlag[2][i] = pcCU.GetTransformSkip1(TLibCommon.TEXT_CHROMA_V)[i+uiPartOffset] //, uiQPartNum * sizeof( byte ) );
                }
            }
        }   // Mode loop
        //#endif

        //--- update overall distortion ---
        uiOverallDistY += uiBestPUDistY
        uiOverallDistC += uiBestPUDistC

        //--- update transform index and cbf ---
        uiQPartNum := pcCU.GetPic().GetNumPartInCU() >> ((uint(pcCU.GetDepth1(0)) + uiInitTrDepth) << 1)
        for i := uint(0); i < uiQPartNum; i++ {
            pcCU.GetTransformIdx()[i+uiPartOffset] = this.m_puhQTTempTrIdx[i]                                          // uiQPartNum * sizeof( byte ) );
            pcCU.GetCbf1(TLibCommon.TEXT_LUMA)[i+uiPartOffset] = this.m_puhQTTempCbf[0][i]                             // uiQPartNum * sizeof( byte ) );
            pcCU.GetCbf1(TLibCommon.TEXT_CHROMA_U)[i+uiPartOffset] = this.m_puhQTTempCbf[1][i]                         // uiQPartNum * sizeof( byte ) );
            pcCU.GetCbf1(TLibCommon.TEXT_CHROMA_V)[i+uiPartOffset] = this.m_puhQTTempCbf[2][i]                         //, uiQPartNum * sizeof( byte ) );
            pcCU.GetTransformSkip1(TLibCommon.TEXT_LUMA)[i+uiPartOffset] = this.m_puhQTTempTransformSkipFlag[0][i]     //, uiQPartNum * sizeof( byte ) );
            pcCU.GetTransformSkip1(TLibCommon.TEXT_CHROMA_U)[i+uiPartOffset] = this.m_puhQTTempTransformSkipFlag[1][i] //, uiQPartNum * sizeof( byte ) );
            pcCU.GetTransformSkip1(TLibCommon.TEXT_CHROMA_V)[i+uiPartOffset] = this.m_puhQTTempTransformSkipFlag[2][i] //, uiQPartNum * sizeof( byte ) );
        }
        //--- set reconstruction for next intra prediction blocks ---
        if uiPU != uiNumPU-1 {
            bSkipChroma := false
            bChromaSame := false
            uiLog2TrSize := TLibCommon.G_aucConvertToBit[pcCU.GetSlice().GetSPS().GetMaxCUWidth()>>(uint(pcCU.GetDepth1(0))+uiInitTrDepth)] + 2
            if !bLumaOnly && uiLog2TrSize == 2 {
                //assert( uiInitTrDepth  > 0 );
                bSkipChroma = (uiPU != 0)
                bChromaSame = true
            }

            uiCompWidth := uint(pcCU.GetWidth1(0)) >> uiInitTrDepth
            uiCompHeight := uint(pcCU.GetHeight1(0)) >> uiInitTrDepth
            uiZOrder := pcCU.GetZorderIdxInCU() + uiPartOffset
            piDes := pcCU.GetPic().GetPicYuvRec().GetLumaAddr2(int(pcCU.GetAddr()), int(uiZOrder))
            uiDesStride := uint(pcCU.GetPic().GetPicYuvRec().GetStride())
            piSrc := pcRecoYuv.GetLumaAddr1(uiPartOffset)
            uiSrcStride := uint(pcRecoYuv.GetStride())
            for uiY := uint(0); uiY < uiCompHeight; uiY++ {
                for uiX := uint(0); uiX < uiCompWidth; uiX++ {
                    piDes[uiY*uiDesStride+uiX] = piSrc[uiY*uiSrcStride+uiX]
                }
                //piSrc = piSrc[uiSrcStride:]
                //piDes = piDes[uiDesStride:]
            }
            if !bLumaOnly && !bSkipChroma {
                if !bChromaSame {
                    uiCompWidth >>= 1
                    uiCompHeight >>= 1
                }
                piDes = pcCU.GetPic().GetPicYuvRec().GetCbAddr2(int(pcCU.GetAddr()), int(uiZOrder))
                uiDesStride = uint(pcCU.GetPic().GetPicYuvRec().GetCStride())
                piSrc = pcRecoYuv.GetCbAddr1(uiPartOffset)
                uiSrcStride = uint(pcRecoYuv.GetCStride())
                for uiY := uint(0); uiY < uiCompHeight; uiY++ {
                    for uiX := uint(0); uiX < uiCompWidth; uiX++ {
                        piDes[uiY*uiDesStride+uiX] = piSrc[uiY*uiSrcStride+uiX]
                    }
                    //piSrc = piSrc[uiSrcStride:]
                    //piDes = piDes[uiDesStride:]
                }
                piDes = pcCU.GetPic().GetPicYuvRec().GetCrAddr2(int(pcCU.GetAddr()), int(uiZOrder))
                piSrc = pcRecoYuv.GetCrAddr1(uiPartOffset)
                for uiY := uint(0); uiY < uiCompHeight; uiY++ {
                    for uiX := uint(0); uiX < uiCompWidth; uiX++ {
                        piDes[uiY*uiDesStride+uiX] = piSrc[uiY*uiSrcStride+uiX]
                    }
                    //piSrc = piSrc[uiSrcStride:]
                    //piDes = piDes[uiDesStride:]
                }
            }
        }

        //=== update PU data ====
        pcCU.SetLumaIntraDirSubParts(uiBestPUMode, uiPartOffset, uiDepth+uiInitTrDepth)
        pcCU.CopyToPic3(uiDepth, uiPU, uiInitTrDepth)

        uiPartOffset += uiQNumParts
    }   // PU loop

    if uiNumPU > 1 {
        // set Cbf for all blocks
        uiCombCbfY := byte(0)
        uiCombCbfU := byte(0)
        uiCombCbfV := byte(0)
        uiPartIdx := uint(0)
        for uiPart := uint(0); uiPart < 4; uiPart++ {
            uiCombCbfY = uiCombCbfY | pcCU.GetCbf3(uiPartIdx, TLibCommon.TEXT_LUMA, 1)
            uiCombCbfU = uiCombCbfU | pcCU.GetCbf3(uiPartIdx, TLibCommon.TEXT_CHROMA_U, 1)
            uiCombCbfV = uiCombCbfV | pcCU.GetCbf3(uiPartIdx, TLibCommon.TEXT_CHROMA_V, 1)

            uiPartIdx += uiQNumParts
        }
        for uiOffs := uint(0); uiOffs < 4*uiQNumParts; uiOffs++ {
            pcCU.GetCbf1(TLibCommon.TEXT_LUMA)[uiOffs] |= uiCombCbfY
            pcCU.GetCbf1(TLibCommon.TEXT_CHROMA_U)[uiOffs] |= uiCombCbfU
            pcCU.GetCbf1(TLibCommon.TEXT_CHROMA_V)[uiOffs] |= uiCombCbfV
        }
    }

    //===== reset context models =====
    if this.m_bUseSBACRD {
        this.m_pcRDGoOnSbacCoder.load(this.m_pppcRDSbacCoder[uiDepth][TLibCommon.CI_CURR_BEST])
    }

    //===== set distortion (rate and r-d costs are determined later) =====
    *ruiDistC = uiOverallDistC
    pcCU.SetTotalDistortion(uiOverallDistY + uiOverallDistC)
    
    //fmt.Printf("Exit estIntraPredQT\n")
}

func (this *TEncSearch) estIntraPredChromaQT(pcCU *TLibCommon.TComDataCU,
    pcOrgYuv *TLibCommon.TComYuv,
    pcPredYuv *TLibCommon.TComYuv,
    pcResiYuv *TLibCommon.TComYuv,
    pcRecoYuv *TLibCommon.TComYuv,
    uiPreCalcDistC uint) {
    uiDepth := uint(pcCU.GetDepth1(0))
    uiBestMode := uint(0)
    uiBestDist := uint(0)
    dBestCost := float64(TLibCommon.MAX_DOUBLE)

    //----- init mode list -----
    uiMinMode := uint(0)
    var uiModeList [TLibCommon.NUM_CHROMA_MODE]uint
    pcCU.GetAllowedChromaDir(0, uiModeList[:])
    uiMaxMode := uint(TLibCommon.NUM_CHROMA_MODE)
	
    //----- check chroma modes -----
    for uiMode := uiMinMode; uiMode < uiMaxMode; uiMode++ {
        //----- restore context models -----
        if this.m_bUseSBACRD {
            this.m_pcRDGoOnSbacCoder.load(this.m_pppcRDSbacCoder[uiDepth][TLibCommon.CI_CURR_BEST])
        }

        //----- chroma coding -----
        uiDist := uint(0)
        pcCU.SetChromIntraDirSubParts(uiModeList[uiMode], 0, uiDepth)
        this.xRecurIntraChromaCodingQT(pcCU, 0, 0, pcOrgYuv, pcPredYuv, pcResiYuv, &uiDist)
        if this.m_bUseSBACRD && pcCU.GetSlice().GetPPS().GetUseTransformSkip() {
            this.m_pcRDGoOnSbacCoder.load(this.m_pppcRDSbacCoder[uiDepth][TLibCommon.CI_CURR_BEST])
        }
        uiBits := this.xGetIntraBitsQT(pcCU, 0, 0, false, true, false)
        dCost := this.m_pcRdCost.calcRdCost(uiBits, uiDist, false, TLibCommon.DF_DEFAULT)

        //----- compare -----
        if dCost < dBestCost {
            dBestCost = dCost
            uiBestDist = uiDist
            uiBestMode = uiModeList[uiMode]
            uiQPN := pcCU.GetPic().GetNumPartInCU() >> (uiDepth << 1)
            this.xSetIntraResultChromaQT(pcCU, 0, 0, pcRecoYuv)

            for i := uint(0); i < uiQPN; i++ {
                this.m_puhQTTempCbf[1][i] = pcCU.GetCbf1(TLibCommon.TEXT_CHROMA_U)[i]                         //, uiQPN * sizeof( byte ) );
                this.m_puhQTTempCbf[2][i] = pcCU.GetCbf1(TLibCommon.TEXT_CHROMA_V)[i]                         //, uiQPN * sizeof( byte ) );
                this.m_puhQTTempTransformSkipFlag[1][i] = pcCU.GetTransformSkip1(TLibCommon.TEXT_CHROMA_U)[i] //, uiQPN * sizeof( byte ) );
                this.m_puhQTTempTransformSkipFlag[2][i] = pcCU.GetTransformSkip1(TLibCommon.TEXT_CHROMA_V)[i] //, uiQPN * sizeof( byte ) );
            }
        }
    }

    //----- set data -----
    uiQPN := pcCU.GetPic().GetNumPartInCU() >> (uiDepth << 1)
    for i := uint(0); i < uiQPN; i++ {
        pcCU.GetCbf1(TLibCommon.TEXT_CHROMA_U)[i] = this.m_puhQTTempCbf[1][i]                         //, uiQPN * sizeof( byte ) );
        pcCU.GetCbf1(TLibCommon.TEXT_CHROMA_V)[i] = this.m_puhQTTempCbf[2][i]                         //, uiQPN * sizeof( byte ) );
        pcCU.GetTransformSkip1(TLibCommon.TEXT_CHROMA_U)[i] = this.m_puhQTTempTransformSkipFlag[1][i] //, uiQPN * sizeof( byte ) );
        pcCU.GetTransformSkip1(TLibCommon.TEXT_CHROMA_V)[i] = this.m_puhQTTempTransformSkipFlag[2][i] //, uiQPN * sizeof( byte ) );
    }
    pcCU.SetChromIntraDirSubParts(uiBestMode, 0, uiDepth)
    pcCU.SetTotalDistortion(pcCU.GetTotalDistortion() + uiBestDist - uiPreCalcDistC)

    //----- restore context models -----
    if this.m_bUseSBACRD {
        this.m_pcRDGoOnSbacCoder.load(this.m_pppcRDSbacCoder[uiDepth][TLibCommon.CI_CURR_BEST])
    }
}

/// encoder estimation - inter prediction (non-skip)
func (this *TEncSearch) predInterSearch(pcCU *TLibCommon.TComDataCU,
    pcOrgYuv *TLibCommon.TComYuv,
    rpcPredYuv *TLibCommon.TComYuv,
    rpcResiYuv *TLibCommon.TComYuv,
    rpcRecoYuv *TLibCommon.TComYuv,
    bUseRes bool, //= false
    //#if AMP_MRG
    bUseMRG bool) { //= false
    //#endif

    this.GetYuvPred(0).Clear()
    this.GetYuvPred(1).Clear()
    this.GetYuvPredTemp().Clear()
    rpcPredYuv.Clear()

    if !bUseRes {
        rpcResiYuv.Clear()
    }

    rpcRecoYuv.Clear()

    var cMvZero, TempMv TLibCommon.TComMv //kolya cMvSrchRngLT, cMvSrchRngRB,

    var cMv [2]TLibCommon.TComMv
    var cMvBi [2]TLibCommon.TComMv
    var cMvTemp [2][33]TLibCommon.TComMv

    iNumPart := int(pcCU.GetNumPartInter())
    var iNumPredDir int
    if pcCU.GetSlice().IsInterP() {
        iNumPredDir = 1
    } else {
        iNumPredDir = 2
    }

    var cMvPred [2][33]TLibCommon.TComMv

    var cMvPredBi [2][33]TLibCommon.TComMv
    var aaiMvpIdxBi [2][33]int

    var aaiMvpIdx [2][33]int
    var aaiMvpNum [2][33]int

    var aacAMVPInfo [2][33]TLibCommon.AMVPInfo

    var iRefIdx = [2]int{0, 0} //If un-initialized, may cause SEGV in bi-directional prediction iterative stage.
    var iRefIdxBi [2]int

    var uiPartAddr uint
    var iRoiWidth, iRoiHeight int

    var uiMbBits = []uint{1, 1, 0}

    uiLastMode := uint(0)
    var iRefStart, iRefEnd int

    ePartSize := pcCU.GetPartitionSize1(0)

    bestBiPRefIdxL1 := 0
    bestBiPMvpL1 := 0
    biPDistTemp := uint(TLibCommon.MAX_INT)

    /*#if ZERO_MVD_EST
      var           aiZeroMvdMvpIdx =[2]int{-1, -1};
      var           aiZeroMvdRefIdx =[2]int{0, 0};
      iZeroMvdDir := -1;
    //#endif*/

    var cMvFieldNeighbours [TLibCommon.MRG_MAX_NUM_CANDS << 1]TLibCommon.TComMvField // double length for mv of both lists
    var uhInterDirNeighbours [TLibCommon.MRG_MAX_NUM_CANDS]byte
    numValidMergeCand := 0

    for iPartIdx := 0; iPartIdx < iNumPart; iPartIdx++ {
        var uiCost = [2]uint{TLibCommon.MAX_UINT, TLibCommon.MAX_UINT}
        uiCostBi := uint(TLibCommon.MAX_UINT)
        var uiCostTemp uint

        var uiBits [3]uint
        var uiBitsTemp uint
        /*#if ZERO_MVD_EST
                         uiZeroMvdCost := uint(TLibCommon.MAX_UINT);
            var          uiZeroMvdCostTemp, uiZeroMvdBitsTemp uint;
                         uiZeroMvdDistTemp := uint(TLibCommon.MAX_UINT);
            var          auiZeroMvdBits [3]uint;
        //#endif*/
        bestBiPDist := uint(TLibCommon.MAX_INT)

        var uiCostTempL0 [TLibCommon.MAX_NUM_REF]uint
        for iNumRef := 0; iNumRef < TLibCommon.MAX_NUM_REF; iNumRef++ {
            uiCostTempL0[iNumRef] = TLibCommon.MAX_UINT
        }
        var uiBitsTempL0 [TLibCommon.MAX_NUM_REF]uint

        this.xGetBlkBits(ePartSize, pcCU.GetSlice().IsInterP(), iPartIdx, uiLastMode, uiMbBits[:])

        pcCU.GetPartIndexAndSize(uint(iPartIdx), &uiPartAddr, &iRoiWidth, &iRoiHeight)

        //#if AMP_MRG
        bTestNormalMC := true

        if bUseMRG && pcCU.GetWidth1(0) > 8 && iNumPart == 2 {
            bTestNormalMC = false
        }

        if bTestNormalMC {
            //#endif

            //  Uni-directional prediction
            for iRefList := 0; iRefList < iNumPredDir; iRefList++ {
                var eRefPicList TLibCommon.RefPicList
                if iRefList != 0 {
                    eRefPicList = TLibCommon.REF_PIC_LIST_1
                } else {
                    eRefPicList = TLibCommon.REF_PIC_LIST_0
                }

                for iRefIdxTemp := 0; iRefIdxTemp < pcCU.GetSlice().GetNumRefIdx(eRefPicList); iRefIdxTemp++ {
                    uiBitsTemp = uiMbBits[iRefList]
                    //fmt.Printf("2.0.0:uiBitsTemp=%d uiMbBits[%d]=[%d,%d,%d]\n",uiBitsTemp,iRefList,uiMbBits[0],uiMbBits[1],uiMbBits[2]);
                    if pcCU.GetSlice().GetNumRefIdx(eRefPicList) > 1 {
                        uiBitsTemp += uint(iRefIdxTemp + 1)
                        if iRefIdxTemp == pcCU.GetSlice().GetNumRefIdx(eRefPicList)-1 {
                            uiBitsTemp--
                        }
                    }
                    //fmt.Printf("2.0.1:uiBitsTemp=%d\n",uiBitsTemp);
                    //#if ZERO_MVD_EST
                    //        this.xEstimateMvPredAMVP( pcCU, pcOrgYuv, uint(iPartIdx), eRefPicList, iRefIdxTemp, cMvPred[iRefList][iRefIdxTemp][:], false, &biPDistTemp, &uiZeroMvdDistTemp);
                    //#else
                    this.xEstimateMvPredAMVP(pcCU, pcOrgYuv, uint(iPartIdx), eRefPicList, iRefIdxTemp, &cMvPred[iRefList][iRefIdxTemp], false, &biPDistTemp)
                    //#endif
                    aaiMvpIdx[iRefList][iRefIdxTemp] = int(pcCU.GetMVPIdx2(eRefPicList, uiPartAddr))
                    aaiMvpNum[iRefList][iRefIdxTemp] = int(pcCU.GetMVPNum2(eRefPicList, uiPartAddr))

                    if pcCU.GetSlice().GetMvdL1ZeroFlag() && iRefList == 1 && biPDistTemp < bestBiPDist {
                        bestBiPDist = biPDistTemp
                        bestBiPMvpL1 = aaiMvpIdx[iRefList][iRefIdxTemp]
                        bestBiPRefIdxL1 = iRefIdxTemp
                    }
					//fmt.Printf("2.0:uiBitsTemp=%d\n",uiBitsTemp);  
                    uiBitsTemp += this.m_auiMVPIdxCost[aaiMvpIdx[iRefList][iRefIdxTemp]][TLibCommon.AMVP_MAX_NUM_CANDS]
                    //fmt.Printf("2.1:uiBitsTemp=%d\n",uiBitsTemp);  
                    /*#if ZERO_MVD_EST
                            if (iRefList != 1 || !pcCU.GetSlice().GetNoBackPredFlag()) &&
                               (pcCU.GetSlice().GetNumRefIdx(TLibCommon.REF_PIC_LIST_C) <= 0 || pcCU.GetSlice().GetRefIdxOfLC(eRefPicList, iRefIdxTemp)>=0) {
                              uiZeroMvdBitsTemp = uiBitsTemp;
                              uiZeroMvdBitsTemp += 2; //zero mvd bits

                              this.m_pcRdCost.GetMotionCost( 1, 0 );
                              uiZeroMvdCostTemp = uiZeroMvdDistTemp + this.m_pcRdCost.GetCost(uiZeroMvdBitsTemp);

                              if uiZeroMvdCostTemp < uiZeroMvdCost {
                                uiZeroMvdCost = uiZeroMvdCostTemp;
                                iZeroMvdDir = iRefList + 1;
                                aiZeroMvdRefIdx[iRefList] = iRefIdxTemp;
                                aiZeroMvdMvpIdx[iRefList] = aaiMvpIdx[iRefList][iRefIdxTemp];
                                auiZeroMvdBits[iRefList] = uiZeroMvdBitsTemp;
                              }
                            }
                    #endif*/

                    //#if GPB_SIMPLE_UNI
                    if pcCU.GetSlice().GetNumRefIdx(TLibCommon.REF_PIC_LIST_C) > 0 {
                        if iRefList != 0 && (pcCU.GetSlice().GetNoBackPredFlag() || (pcCU.GetSlice().GetNumRefIdx(TLibCommon.REF_PIC_LIST_C) > 0 && !pcCU.GetSlice().GetNoBackPredFlag() && pcCU.GetSlice().GetRefIdxOfL0FromRefIdxOfL1(iRefIdxTemp) >= 0)) {
                            if pcCU.GetSlice().GetNoBackPredFlag() {
                                cMvTemp[1][iRefIdxTemp] = cMvTemp[0][iRefIdxTemp]
                                uiCostTemp = uiCostTempL0[iRefIdxTemp]
                                /*first subtract the bit-rate part of the cost of the other list*/
                                uiCostTemp -= this.m_pcRdCost.getCost1(uiBitsTempL0[iRefIdxTemp])
                            } else {
                                cMvTemp[1][iRefIdxTemp] = cMvTemp[0][pcCU.GetSlice().GetRefIdxOfL0FromRefIdxOfL1(iRefIdxTemp)]
                                uiCostTemp = uiCostTempL0[pcCU.GetSlice().GetRefIdxOfL0FromRefIdxOfL1(iRefIdxTemp)]
                                /*first subtract the bit-rate part of the cost of the other list*/
                                uiCostTemp -= this.m_pcRdCost.getCost1(uiBitsTempL0[pcCU.GetSlice().GetRefIdxOfL0FromRefIdxOfL1(iRefIdxTemp)])
                            }
                            /*correct the bit-rate part of the current ref*/
                            this.m_pcRdCost.setPredictor(&cMvPred[iRefList][iRefIdxTemp])
                            uiBitsTemp += this.m_pcRdCost.getBits(int(cMvTemp[1][iRefIdxTemp].GetHor()), int(cMvTemp[1][iRefIdxTemp].GetVer()))
                            /*calculate the correct cost*/
                            uiCostTemp += this.m_pcRdCost.getCost1(uiBitsTemp)
                        } else {
                            this.xMotionEstimation(pcCU, pcOrgYuv, iPartIdx, eRefPicList, &cMvPred[iRefList][iRefIdxTemp], iRefIdxTemp, &cMvTemp[iRefList][iRefIdxTemp], &uiBitsTemp, &uiCostTemp, false)
                        }
                    } else {
                        if iRefList != 0 && pcCU.GetSlice().GetNoBackPredFlag() {
                            uiCostTemp = TLibCommon.MAX_UINT
                            cMvTemp[1][iRefIdxTemp] = cMvTemp[0][iRefIdxTemp]
                        } else {
                        	//fmt.Printf("2:uiBitsTemp=%d\n",uiBitsTemp);
                            this.xMotionEstimation(pcCU, pcOrgYuv, iPartIdx, eRefPicList, &cMvPred[iRefList][iRefIdxTemp], iRefIdxTemp, &cMvTemp[iRefList][iRefIdxTemp], &uiBitsTemp, &uiCostTemp, false)
                        }
                    }
                    //#else
                    //        xMotionEstimation ( pcCU, pcOrgYuv, iPartIdx, eRefPicList, &cMvPred[iRefList][iRefIdxTemp], iRefIdxTemp, cMvTemp[iRefList][iRefIdxTemp], uiBitsTemp, uiCostTemp );
                    //#endif
                    this.xCopyAMVPInfo(pcCU.GetCUMvField(eRefPicList).GetAMVPInfo(), &aacAMVPInfo[iRefList][iRefIdxTemp]) // must always be done ( also when AMVP_MODE = AM_NONE )
                    this.xCheckBestMVP(pcCU, eRefPicList, cMvTemp[iRefList][iRefIdxTemp], &cMvPred[iRefList][iRefIdxTemp], &aaiMvpIdx[iRefList][iRefIdxTemp], &uiBitsTemp, &uiCostTemp)

                    if pcCU.GetSlice().GetNumRefIdx(TLibCommon.REF_PIC_LIST_C) > 0 && !pcCU.GetSlice().GetNoBackPredFlag() {
                        if iRefList == TLibCommon.REF_PIC_LIST_0 {
                            uiCostTempL0[iRefIdxTemp] = uiCostTemp
                            uiBitsTempL0[iRefIdxTemp] = uiBitsTemp
                            if pcCU.GetSlice().GetRefIdxOfLC(TLibCommon.REF_PIC_LIST_0, iRefIdxTemp) < 0 {
                                uiCostTemp = uint(TLibCommon.MAX_UINT)
                            }
                        } else {
                            if pcCU.GetSlice().GetRefIdxOfLC(TLibCommon.REF_PIC_LIST_1, iRefIdxTemp) < 0 {
                                uiCostTemp = uint(TLibCommon.MAX_UINT)
                            }
                        }
                    }

                    if (iRefList == 0 && uiCostTemp < uiCost[iRefList]) ||
                        (iRefList == 1 && pcCU.GetSlice().GetNoBackPredFlag() && iRefIdxTemp == iRefIdx[0]) ||
                        (iRefList == 1 && (pcCU.GetSlice().GetNumRefIdx(TLibCommon.REF_PIC_LIST_C) > 0) && (iRefIdxTemp == 0 || iRefIdxTemp == iRefIdx[0]) && !pcCU.GetSlice().GetNoBackPredFlag() && (iRefIdxTemp == pcCU.GetSlice().GetRefIdxOfL0FromRefIdxOfL1(iRefIdxTemp))) ||
                        (iRefList == 1 && !pcCU.GetSlice().GetNoBackPredFlag() && uiCostTemp < uiCost[iRefList]) {
                        uiCost[iRefList] = uiCostTemp
                        uiBits[iRefList] = uiBitsTemp // storing for bi-prediction

                        // set motion
                        cMv[iRefList] = cMvTemp[iRefList][iRefIdxTemp]
                        iRefIdx[iRefList] = iRefIdxTemp
                        pcCU.GetCUMvField(eRefPicList).SetAllMv(cMv[iRefList], ePartSize, int(uiPartAddr), 0, iPartIdx)
                        pcCU.GetCUMvField(eRefPicList).SetAllRefIdx(int8(iRefIdx[iRefList]), ePartSize, int(uiPartAddr), 0, iPartIdx)

                        if !pcCU.GetSlice().GetMvdL1ZeroFlag() {
                            // storing list 1 prediction signal for iterative bi-directional prediction
                            if eRefPicList == TLibCommon.REF_PIC_LIST_1 {
                                pcYuvPred := this.GetYuvPred(iRefList)
                                this.MotionCompensation(pcCU, pcYuvPred, eRefPicList, iPartIdx)
                            }
                            if (pcCU.GetSlice().GetNoBackPredFlag() || (pcCU.GetSlice().GetNumRefIdx(TLibCommon.REF_PIC_LIST_C) > 0 && pcCU.GetSlice().GetRefIdxOfL0FromRefIdxOfL1(0) == 0)) && eRefPicList == TLibCommon.REF_PIC_LIST_0 {
                                pcYuvPred := this.GetYuvPred(iRefList)
                                this.MotionCompensation(pcCU, pcYuvPred, eRefPicList, iPartIdx)
                            }
                        }
                    }
                }
            }
            //  Bi-directional prediction
            if (pcCU.GetSlice().IsInterB()) && (pcCU.IsBipredRestriction(uint(iPartIdx)) == false) {

                cMvBi[0] = cMv[0]
                cMvBi[1] = cMv[1]
                iRefIdxBi[0] = iRefIdx[0]
                iRefIdxBi[1] = iRefIdx[1]

                cMvPredBi = cMvPred
                aaiMvpIdxBi = aaiMvpIdx
                //::memcpy(cMvPredBi, cMvPred, sizeof(cMvPred));
                //::memcpy(aaiMvpIdxBi, aaiMvpIdx, sizeof(aaiMvpIdx));

                var uiMotBits [2]uint

                if pcCU.GetSlice().GetMvdL1ZeroFlag() {
                    this.xCopyAMVPInfo(&aacAMVPInfo[1][bestBiPRefIdxL1], pcCU.GetCUMvField(TLibCommon.REF_PIC_LIST_1).GetAMVPInfo())
                    pcCU.SetMVPIdxSubParts(bestBiPMvpL1, TLibCommon.REF_PIC_LIST_1, uiPartAddr, uint(iPartIdx), uint(pcCU.GetDepth1(uiPartAddr)))
                    aaiMvpIdxBi[1][bestBiPRefIdxL1] = bestBiPMvpL1
                    cMvPredBi[1][bestBiPRefIdxL1] = pcCU.GetCUMvField(TLibCommon.REF_PIC_LIST_1).GetAMVPInfo().MvCand[bestBiPMvpL1]

                    cMvBi[1] = cMvPredBi[1][bestBiPRefIdxL1]
                    iRefIdxBi[1] = bestBiPRefIdxL1
                    pcCU.GetCUMvField(TLibCommon.REF_PIC_LIST_1).SetAllMv(cMvBi[1], ePartSize, int(uiPartAddr), 0, iPartIdx)
                    pcCU.GetCUMvField(TLibCommon.REF_PIC_LIST_1).SetAllRefIdx(int8(iRefIdxBi[1]), ePartSize, int(uiPartAddr), 0, iPartIdx)
                    pcYuvPred := this.GetYuvPred(1)
                    this.MotionCompensation(pcCU, pcYuvPred, TLibCommon.REF_PIC_LIST_1, iPartIdx)

                    uiMotBits[0] = uiBits[0] - uiMbBits[0]
                    uiMotBits[1] = uiMbBits[1]

                    if pcCU.GetSlice().GetNumRefIdx(TLibCommon.REF_PIC_LIST_1) > 1 {
                        uiMotBits[1] += uint(bestBiPRefIdxL1) + 1
                        if bestBiPRefIdxL1 == pcCU.GetSlice().GetNumRefIdx(TLibCommon.REF_PIC_LIST_1)-1 {
                            uiMotBits[1]--
                        }
                    }

                    uiMotBits[1] += this.m_auiMVPIdxCost[aaiMvpIdxBi[1][bestBiPRefIdxL1]][TLibCommon.AMVP_MAX_NUM_CANDS]

                    uiBits[2] = uiMbBits[2] + uiMotBits[0] + uiMotBits[1]

                    cMvTemp[1][bestBiPRefIdxL1] = cMvBi[1]
                } else {
                    uiMotBits[0] = uiBits[0] - uiMbBits[0]
                    uiMotBits[1] = uiBits[1] - uiMbBits[1]
                    uiBits[2] = uiMbBits[2] + uiMotBits[0] + uiMotBits[1]
                }

                // 4-times iteration (default)
                iNumIter := 4

                // fast encoder setting: only one iteration
                if this.m_pcEncCfg.GetUseFastEnc() || pcCU.GetSlice().GetMvdL1ZeroFlag() {
                    iNumIter = 1
                }

                for iIter := 0; iIter < iNumIter; iIter++ {
                    iRefList := iIter % 2
                    if this.m_pcEncCfg.GetUseFastEnc() && (pcCU.GetSlice().GetNoBackPredFlag() || (pcCU.GetSlice().GetNumRefIdx(TLibCommon.REF_PIC_LIST_C) > 0 && pcCU.GetSlice().GetRefIdxOfL0FromRefIdxOfL1(0) == 0)) {
                        iRefList = 1
                    }
                    var eRefPicList TLibCommon.RefPicList
                    if iRefList != 0 {
                        eRefPicList = TLibCommon.REF_PIC_LIST_1
                    } else {
                        eRefPicList = TLibCommon.REF_PIC_LIST_0
                    }

                    if pcCU.GetSlice().GetMvdL1ZeroFlag() {
                        iRefList = 0
                        eRefPicList = TLibCommon.REF_PIC_LIST_0
                    }

                    bChanged := false

                    iRefStart = 0
                    iRefEnd = pcCU.GetSlice().GetNumRefIdx(eRefPicList) - 1

                    for iRefIdxTemp := iRefStart; iRefIdxTemp <= iRefEnd; iRefIdxTemp++ {
                        uiBitsTemp = uiMbBits[2] + uiMotBits[1-iRefList]
                        if pcCU.GetSlice().GetNumRefIdx(eRefPicList) > 1 {
                            uiBitsTemp += uint(iRefIdxTemp) + 1
                            if iRefIdxTemp == pcCU.GetSlice().GetNumRefIdx(eRefPicList)-1 {
                                uiBitsTemp--
                            }
                        }
                        uiBitsTemp += this.m_auiMVPIdxCost[aaiMvpIdxBi[iRefList][iRefIdxTemp]][TLibCommon.AMVP_MAX_NUM_CANDS]
                        // call ME
                        this.xMotionEstimation(pcCU, pcOrgYuv, iPartIdx, eRefPicList, &cMvPredBi[iRefList][iRefIdxTemp], iRefIdxTemp, &cMvTemp[iRefList][iRefIdxTemp], &uiBitsTemp, &uiCostTemp, true)
                        this.xCopyAMVPInfo(&aacAMVPInfo[iRefList][iRefIdxTemp], pcCU.GetCUMvField(eRefPicList).GetAMVPInfo())
                        this.xCheckBestMVP(pcCU, eRefPicList, cMvTemp[iRefList][iRefIdxTemp], &cMvPredBi[iRefList][iRefIdxTemp], &aaiMvpIdxBi[iRefList][iRefIdxTemp], &uiBitsTemp, &uiCostTemp)

                        if uiCostTemp < uiCostBi {
                            bChanged = true

                            cMvBi[iRefList] = cMvTemp[iRefList][iRefIdxTemp]
                            iRefIdxBi[iRefList] = iRefIdxTemp

                            uiCostBi = uiCostTemp
                            uiMotBits[iRefList] = uiBitsTemp - uiMbBits[2] - uiMotBits[1-iRefList]
                            uiBits[2] = uiBitsTemp

                            if iNumIter != 1 {
                                //  Set motion
                                pcCU.GetCUMvField(eRefPicList).SetAllMv(cMvBi[iRefList], ePartSize, int(uiPartAddr), 0, iPartIdx)
                                pcCU.GetCUMvField(eRefPicList).SetAllRefIdx(int8(iRefIdxBi[iRefList]), ePartSize, int(uiPartAddr), 0, iPartIdx)

                                pcYuvPred := this.GetYuvPred(iRefList)
                                this.MotionCompensation(pcCU, pcYuvPred, eRefPicList, iPartIdx)
                            }
                        }
                    }   // for loop-iRefIdxTemp

                    if !bChanged {
                        if uiCostBi <= uiCost[0] && uiCostBi <= uiCost[1] {
                            this.xCopyAMVPInfo(&aacAMVPInfo[0][iRefIdxBi[0]], pcCU.GetCUMvField(TLibCommon.REF_PIC_LIST_0).GetAMVPInfo())
                            this.xCheckBestMVP(pcCU, TLibCommon.REF_PIC_LIST_0, cMvBi[0], &cMvPredBi[0][iRefIdxBi[0]], &aaiMvpIdxBi[0][iRefIdxBi[0]], &uiBits[2], &uiCostBi)
                            if !pcCU.GetSlice().GetMvdL1ZeroFlag() {
                                this.xCopyAMVPInfo(&aacAMVPInfo[1][iRefIdxBi[1]], pcCU.GetCUMvField(TLibCommon.REF_PIC_LIST_1).GetAMVPInfo())
                                this.xCheckBestMVP(pcCU, TLibCommon.REF_PIC_LIST_1, cMvBi[1], &cMvPredBi[1][iRefIdxBi[1]], &aaiMvpIdxBi[1][iRefIdxBi[1]], &uiBits[2], &uiCostBi)
                            }
                        }
                        break
                    }
                }   // for loop-iter
            }   // if (B_SLICE)
            /*#if ZERO_MVD_EST
                  if (pcCU.GetSlice().IsInterB()) && (pcCU.IsBipredRestriction(iPartIdx) == false) {
                    this.m_pcRdCost.GetMotionCost( 1, 0 );

                    for iL0RefIdxTemp := 0; iL0RefIdxTemp <= pcCU.GetSlice().GetNumRefIdx(TLibCommon.REF_PIC_LIST_0)-1; iL0RefIdxTemp++ {
              	      for iL1RefIdxTemp := 0; iL1RefIdxTemp <= pcCU.GetSlice().GetNumRefIdx(TLibCommon.REF_PIC_LIST_1)-1; iL1RefIdxTemp++ {
              	        uiRefIdxBitsTemp := uint(0);
              	        if pcCU.GetSlice().GetNumRefIdx(TLibCommon.REF_PIC_LIST_0) > 1 {
              	          uiRefIdxBitsTemp += iL0RefIdxTemp+1;
              	          if iL0RefIdxTemp == pcCU.GetSlice().GetNumRefIdx(TLibCommon.REF_PIC_LIST_0)-1 {
              	          	 uiRefIdxBitsTemp--;
              	          }
              	        }
              	        if pcCU.GetSlice().GetNumRefIdx(TLibCommon.REF_PIC_LIST_1) > 1 {
              	          uiRefIdxBitsTemp += iL1RefIdxTemp+1;
              	          if iL1RefIdxTemp == pcCU.GetSlice().GetNumRefIdx(TLibCommon.REF_PIC_LIST_1)-1 {
              	           uiRefIdxBitsTemp--;
              	          }
              	        }

              	        iL0MVPIdx := 0;
              	        iL1MVPIdx := 0;

              	        for iL0MVPIdx = 0; iL0MVPIdx < aaiMvpNum[0][iL0RefIdxTemp]; iL0MVPIdx++ {
              	          for iL1MVPIdx = 0; iL1MVPIdx < aaiMvpNum[1][iL1RefIdxTemp]; iL1MVPIdx++ {
              	            uiZeroMvdBitsTemp = uiRefIdxBitsTemp;
              	            uiZeroMvdBitsTemp += uiMbBits[2];
              	            uiZeroMvdBitsTemp += this.m_auiMVPIdxCost[iL0MVPIdx][aaiMvpNum[0][iL0RefIdxTemp]] + this.m_auiMVPIdxCost[iL1MVPIdx][aaiMvpNum[1][iL1RefIdxTemp]];
              	            uiZeroMvdBitsTemp += 4; //zero mvd for both directions
              	            pcCU.GetCUMvField( TLibCommon.REF_PIC_LIST_0 ).SetAllMvField( aacAMVPInfo[0][iL0RefIdxTemp].this.MvCand[iL0MVPIdx], iL0RefIdxTemp, ePartSize, uiPartAddr, iPartIdx, 0 );
              	            pcCU.GetCUMvField( TLibCommon.REF_PIC_LIST_1 ).SetAllMvField( aacAMVPInfo[1][iL1RefIdxTemp].this.MvCand[iL1MVPIdx], iL1RefIdxTemp, ePartSize, uiPartAddr, iPartIdx, 0 );

              	            this.xGetInterPredictionError( pcCU, pcOrgYuv, iPartIdx, uiZeroMvdDistTemp, this.m_pcEncCfg.GetUseHADME() );
              	            uiZeroMvdCostTemp = uiZeroMvdDistTemp + this.m_pcRdCost.GetCost( uiZeroMvdBitsTemp );
              	            if uiZeroMvdCostTemp < uiZeroMvdCost {
              	              uiZeroMvdCost = uiZeroMvdCostTemp;
              	              iZeroMvdDir = 3;
              	              aiZeroMvdMvpIdx[0] = iL0MVPIdx;
              	              aiZeroMvdMvpIdx[1] = iL1MVPIdx;
              	              aiZeroMvdRefIdx[0] = iL0RefIdxTemp;
              	              aiZeroMvdRefIdx[1] = iL1RefIdxTemp;
              	              auiZeroMvdBits[2] = uiZeroMvdBitsTemp;
              	            }
              	          }
              	        }
              	      }
              	   }
                  }
              //#endif*/

            //#if AMP_MRG
        }   //end if bTestNormalMC
        //#endif
        //  Clear Motion Field
        pcCU.GetCUMvField(TLibCommon.REF_PIC_LIST_0).SetAllMvField(TLibCommon.NewTComMvField(), ePartSize, int(uiPartAddr), 0, iPartIdx)
        pcCU.GetCUMvField(TLibCommon.REF_PIC_LIST_1).SetAllMvField(TLibCommon.NewTComMvField(), ePartSize, int(uiPartAddr), 0, iPartIdx)
        pcCU.GetCUMvField(TLibCommon.REF_PIC_LIST_0).SetAllMvd(cMvZero, ePartSize, int(uiPartAddr), 0, iPartIdx)
        pcCU.GetCUMvField(TLibCommon.REF_PIC_LIST_1).SetAllMvd(cMvZero, ePartSize, int(uiPartAddr), 0, iPartIdx)

        pcCU.SetMVPIdxSubParts(-1, TLibCommon.REF_PIC_LIST_0, uiPartAddr, uint(iPartIdx), uint(pcCU.GetDepth1(uiPartAddr)))
        pcCU.SetMVPNumSubParts(-1, TLibCommon.REF_PIC_LIST_0, uiPartAddr, uint(iPartIdx), uint(pcCU.GetDepth1(uiPartAddr)))
        pcCU.SetMVPIdxSubParts(-1, TLibCommon.REF_PIC_LIST_1, uiPartAddr, uint(iPartIdx), uint(pcCU.GetDepth1(uiPartAddr)))
        pcCU.SetMVPNumSubParts(-1, TLibCommon.REF_PIC_LIST_1, uiPartAddr, uint(iPartIdx), uint(pcCU.GetDepth1(uiPartAddr)))

        uiMEBits := uint(0)
        // Set Motion Field_
        if pcCU.GetSlice().GetNoBackPredFlag() || (pcCU.GetSlice().GetNumRefIdx(TLibCommon.REF_PIC_LIST_C) > 0 && pcCU.GetSlice().GetRefIdxOfL0FromRefIdxOfL1(0) == 0) {
            uiCost[1] = uint(TLibCommon.MAX_UINT)
        }
        //#if AMP_MRG
        if bTestNormalMC {
            //#endif
            /*#if ZERO_MVD_EST
                if uiZeroMvdCost <= uiCostBi && uiZeroMvdCost <= uiCost[0] && uiZeroMvdCost <= uiCost[1] {
                  if iZeroMvdDir == 3 {
                    uiLastMode = 2;

                    pcCU.GetCUMvField(TLibCommon.REF_PIC_LIST_0).SetAllMvField( aacAMVPInfo[0][aiZeroMvdRefIdx[0]].this.MvCand[aiZeroMvdMvpIdx[0]], aiZeroMvdRefIdx[0], ePartSize, uiPartAddr, iPartIdx, 0 );
                    pcCU.GetCUMvField(TLibCommon.REF_PIC_LIST_1).SetAllMvField( aacAMVPInfo[1][aiZeroMvdRefIdx[1]].this.MvCand[aiZeroMvdMvpIdx[1]], aiZeroMvdRefIdx[1], ePartSize, uiPartAddr, iPartIdx, 0 );

                    pcCU.SetInterDirSubParts( 3, uiPartAddr, iPartIdx, pcCU.GetDepth(0) );

                    pcCU.SetMVPIdxSubParts( aiZeroMvdMvpIdx[0], TLibCommon.REF_PIC_LIST_0, uiPartAddr, iPartIdx, pcCU.GetDepth(uiPartAddr));
                    pcCU.SetMVPNumSubParts( aaiMvpNum[0][aiZeroMvdRefIdx[0]], TLibCommon.REF_PIC_LIST_0, uiPartAddr, iPartIdx, pcCU.GetDepth(uiPartAddr));
                    pcCU.SetMVPIdxSubParts( aiZeroMvdMvpIdx[1], TLibCommon.REF_PIC_LIST_1, uiPartAddr, iPartIdx, pcCU.GetDepth(uiPartAddr));
                    pcCU.SetMVPNumSubParts( aaiMvpNum[1][aiZeroMvdRefIdx[1]], TLibCommon.REF_PIC_LIST_1, uiPartAddr, iPartIdx, pcCU.GetDepth(uiPartAddr));
                    uiMEBits = auiZeroMvdBits[2];
                  }else if iZeroMvdDir == 1 {
                    uiLastMode = 0;

                    pcCU.GetCUMvField(TLibCommon.REF_PIC_LIST_0).SetAllMvField( aacAMVPInfo[0][aiZeroMvdRefIdx[0]].this.MvCand[aiZeroMvdMvpIdx[0]], aiZeroMvdRefIdx[0], ePartSize, uiPartAddr, iPartIdx, 0 );

                    pcCU.SetInterDirSubParts( 1, uiPartAddr, iPartIdx, pcCU.GetDepth(0) );

                    pcCU.SetMVPIdxSubParts( aiZeroMvdMvpIdx[0], TLibCommon.REF_PIC_LIST_0, uiPartAddr, iPartIdx, pcCU.GetDepth(uiPartAddr));
                    pcCU.SetMVPNumSubParts( aaiMvpNum[0][aiZeroMvdRefIdx[0]], TLibCommon.REF_PIC_LIST_0, uiPartAddr, iPartIdx, pcCU.GetDepth(uiPartAddr));
                    uiMEBits = auiZeroMvdBits[0];
                  }else if iZeroMvdDir == 2 {
                    uiLastMode = 1;

                    pcCU.GetCUMvField(TLibCommon.REF_PIC_LIST_1).SetAllMvField( aacAMVPInfo[1][aiZeroMvdRefIdx[1]].this.MvCand[aiZeroMvdMvpIdx[1]], aiZeroMvdRefIdx[1], ePartSize, uiPartAddr, iPartIdx, 0 );

                    pcCU.SetInterDirSubParts( 2, uiPartAddr, iPartIdx, pcCU.GetDepth(0) );

                    pcCU.SetMVPIdxSubParts( aiZeroMvdMvpIdx[1], TLibCommon.REF_PIC_LIST_1, uiPartAddr, iPartIdx, pcCU.GetDepth(uiPartAddr));
                    pcCU.SetMVPNumSubParts( aaiMvpNum[1][aiZeroMvdRefIdx[1]], TLibCommon.REF_PIC_LIST_1, uiPartAddr, iPartIdx, pcCU.GetDepth(uiPartAddr));
                    uiMEBits = auiZeroMvdBits[1];
                  }else{
                    //assert(0);
                  }
                }else{
            //#endif*/
            {
                if uiCostBi <= uiCost[0] && uiCostBi <= uiCost[1] {
                    uiLastMode = 2
                    {
                        pcCU.GetCUMvField(TLibCommon.REF_PIC_LIST_0).SetAllMv(cMvBi[0], ePartSize, int(uiPartAddr), 0, iPartIdx)
                        pcCU.GetCUMvField(TLibCommon.REF_PIC_LIST_0).SetAllRefIdx(int8(iRefIdxBi[0]), ePartSize, int(uiPartAddr), 0, iPartIdx)
                        pcCU.GetCUMvField(TLibCommon.REF_PIC_LIST_1).SetAllMv(cMvBi[1], ePartSize, int(uiPartAddr), 0, iPartIdx)
                        pcCU.GetCUMvField(TLibCommon.REF_PIC_LIST_1).SetAllRefIdx(int8(iRefIdxBi[1]), ePartSize, int(uiPartAddr), 0, iPartIdx)
                    }
                    {
                        TempMv = TLibCommon.SubMvs(cMvBi[0], cMvPredBi[0][iRefIdxBi[0]])
                        pcCU.GetCUMvField(TLibCommon.REF_PIC_LIST_0).SetAllMvd(TempMv, ePartSize, int(uiPartAddr), 0, iPartIdx)
                    }
                    {
                        TempMv = TLibCommon.SubMvs(cMvBi[1], cMvPredBi[1][iRefIdxBi[1]])
                        pcCU.GetCUMvField(TLibCommon.REF_PIC_LIST_1).SetAllMvd(TempMv, ePartSize, int(uiPartAddr), 0, iPartIdx)
                    }

                    pcCU.SetInterDirSubParts(3, uiPartAddr, uint(iPartIdx), uint(pcCU.GetDepth1(0)))

                    pcCU.SetMVPIdxSubParts(aaiMvpIdxBi[0][iRefIdxBi[0]], TLibCommon.REF_PIC_LIST_0, uiPartAddr, uint(iPartIdx), uint(pcCU.GetDepth1(uiPartAddr)))
                    pcCU.SetMVPNumSubParts(aaiMvpNum[0][iRefIdxBi[0]], TLibCommon.REF_PIC_LIST_0, uiPartAddr, uint(iPartIdx), uint(pcCU.GetDepth1(uiPartAddr)))
                    pcCU.SetMVPIdxSubParts(aaiMvpIdxBi[1][iRefIdxBi[1]], TLibCommon.REF_PIC_LIST_1, uiPartAddr, uint(iPartIdx), uint(pcCU.GetDepth1(uiPartAddr)))
                    pcCU.SetMVPNumSubParts(aaiMvpNum[1][iRefIdxBi[1]], TLibCommon.REF_PIC_LIST_1, uiPartAddr, uint(iPartIdx), uint(pcCU.GetDepth1(uiPartAddr)))

                    uiMEBits = uiBits[2]
                } else if uiCost[0] <= uiCost[1] {
                    uiLastMode = 0
                    pcCU.GetCUMvField(TLibCommon.REF_PIC_LIST_0).SetAllMv(cMv[0], ePartSize, int(uiPartAddr), 0, iPartIdx)
                    pcCU.GetCUMvField(TLibCommon.REF_PIC_LIST_0).SetAllRefIdx(int8(iRefIdx[0]), ePartSize, int(uiPartAddr), 0, iPartIdx)
                    {
                        TempMv = TLibCommon.SubMvs(cMv[0], cMvPred[0][iRefIdx[0]])
                        pcCU.GetCUMvField(TLibCommon.REF_PIC_LIST_0).SetAllMvd(TempMv, ePartSize, int(uiPartAddr), 0, iPartIdx)
                    }
                    pcCU.SetInterDirSubParts(1, uiPartAddr, uint(iPartIdx), uint(pcCU.GetDepth1(0)))

                    pcCU.SetMVPIdxSubParts(aaiMvpIdx[0][iRefIdx[0]], TLibCommon.REF_PIC_LIST_0, uiPartAddr, uint(iPartIdx), uint(pcCU.GetDepth1(uiPartAddr)))
                    pcCU.SetMVPNumSubParts(aaiMvpNum[0][iRefIdx[0]], TLibCommon.REF_PIC_LIST_0, uiPartAddr, uint(iPartIdx), uint(pcCU.GetDepth1(uiPartAddr)))

                    uiMEBits = uiBits[0]
                } else {
                    uiLastMode = 1
                    pcCU.GetCUMvField(TLibCommon.REF_PIC_LIST_1).SetAllMv(cMv[1], ePartSize, int(uiPartAddr), 0, iPartIdx)
                    pcCU.GetCUMvField(TLibCommon.REF_PIC_LIST_1).SetAllRefIdx(int8(iRefIdx[1]), ePartSize, int(uiPartAddr), 0, iPartIdx)
                    {
                        TempMv = TLibCommon.SubMvs(cMv[1], cMvPred[1][iRefIdx[1]])
                        pcCU.GetCUMvField(TLibCommon.REF_PIC_LIST_1).SetAllMvd(TempMv, ePartSize, int(uiPartAddr), 0, iPartIdx)
                    }
                    pcCU.SetInterDirSubParts(2, uiPartAddr, uint(iPartIdx), uint(pcCU.GetDepth1(0)))

                    pcCU.SetMVPIdxSubParts(aaiMvpIdx[1][iRefIdx[1]], TLibCommon.REF_PIC_LIST_1, uiPartAddr, uint(iPartIdx), uint(pcCU.GetDepth1(uiPartAddr)))
                    pcCU.SetMVPNumSubParts(aaiMvpNum[1][iRefIdx[1]], TLibCommon.REF_PIC_LIST_1, uiPartAddr, uint(iPartIdx), uint(pcCU.GetDepth1(uiPartAddr)))

                    uiMEBits = uiBits[1]
                }
            }
            //#if AMP_MRG
        }   // end if bTestNormalMC
        //#endif

        if pcCU.GetPartitionSize1(uiPartAddr) != TLibCommon.SIZE_2Nx2N {
            uiMRGInterDir := uint(0)
            var cMRGMvField [2]TLibCommon.TComMvField
            uiMRGIndex := uint(0)

            uiMEInterDir := uint(0)
            var cMEMvField [2]TLibCommon.TComMvField

            this.m_pcRdCost.getMotionCost(true, 0)
            //#if AMP_MRG
            // calculate ME cost
            uiMEError := uint(TLibCommon.MAX_UINT)
            uiMECost := uint(TLibCommon.MAX_UINT)

            if bTestNormalMC {
                this.xGetInterPredictionError(pcCU, pcOrgYuv, iPartIdx, &uiMEError, this.m_pcEncCfg.GetUseHADME())
                uiMECost = uiMEError + this.m_pcRdCost.getCost1(uiMEBits)
            }
            //#else
            // calculate ME cost
            //      UInt uiMEError = MAX_UINT;
            //      xGetInterPredictionError( pcCU, pcOrgYuv, iPartIdx, uiMEError, this.m_pcEncCfg.GetUseHADME() );
            //      UInt uiMECost = uiMEError + this.m_pcRdCost.GetCost( uiMEBits );
            //#endif
            // save ME result.
            uiMEInterDir = uint(pcCU.GetInterDir1(uiPartAddr))
            pcCU.GetMvField(pcCU, uiPartAddr, TLibCommon.REF_PIC_LIST_0, &cMEMvField[0])
            pcCU.GetMvField(pcCU, uiPartAddr, TLibCommon.REF_PIC_LIST_1, &cMEMvField[1])

            // find Merge result
            uiMRGCost := uint(TLibCommon.MAX_UINT)
            this.xMergeEstimation(pcCU, pcOrgYuv, iPartIdx, &uiMRGInterDir, cMRGMvField[:], &uiMRGIndex, &uiMRGCost, cMvFieldNeighbours[:], uhInterDirNeighbours[:], &numValidMergeCand)
            if uiMRGCost < uiMECost {
                // set Merge result
                pcCU.SetMergeFlagSubParts(true, uiPartAddr, uint(iPartIdx), uint(pcCU.GetDepth1(uiPartAddr)))
                pcCU.SetMergeIndexSubParts(uiMRGIndex, uiPartAddr, uint(iPartIdx), uint(pcCU.GetDepth1(uiPartAddr)))
                pcCU.SetInterDirSubParts(uiMRGInterDir, uiPartAddr, uint(iPartIdx), uint(pcCU.GetDepth1(uiPartAddr)))

                pcCU.GetCUMvField(TLibCommon.REF_PIC_LIST_0).SetAllMvField(&cMRGMvField[0], ePartSize, int(uiPartAddr), 0, iPartIdx)
                pcCU.GetCUMvField(TLibCommon.REF_PIC_LIST_1).SetAllMvField(&cMRGMvField[1], ePartSize, int(uiPartAddr), 0, iPartIdx)

                pcCU.GetCUMvField(TLibCommon.REF_PIC_LIST_0).SetAllMvd(cMvZero, ePartSize, int(uiPartAddr), 0, iPartIdx)
                pcCU.GetCUMvField(TLibCommon.REF_PIC_LIST_1).SetAllMvd(cMvZero, ePartSize, int(uiPartAddr), 0, iPartIdx)

                pcCU.SetMVPIdxSubParts(-1, TLibCommon.REF_PIC_LIST_0, uiPartAddr, uint(iPartIdx), uint(pcCU.GetDepth1(uiPartAddr)))
                pcCU.SetMVPNumSubParts(-1, TLibCommon.REF_PIC_LIST_0, uiPartAddr, uint(iPartIdx), uint(pcCU.GetDepth1(uiPartAddr)))
                pcCU.SetMVPIdxSubParts(-1, TLibCommon.REF_PIC_LIST_1, uiPartAddr, uint(iPartIdx), uint(pcCU.GetDepth1(uiPartAddr)))
                pcCU.SetMVPNumSubParts(-1, TLibCommon.REF_PIC_LIST_1, uiPartAddr, uint(iPartIdx), uint(pcCU.GetDepth1(uiPartAddr)))
            } else {
                // set ME result
                pcCU.SetMergeFlagSubParts(false, uiPartAddr, uint(iPartIdx), uint(pcCU.GetDepth1(uiPartAddr)))
                pcCU.SetInterDirSubParts(uiMEInterDir, uiPartAddr, uint(iPartIdx), uint(pcCU.GetDepth1(uiPartAddr)))

                pcCU.GetCUMvField(TLibCommon.REF_PIC_LIST_0).SetAllMvField(&cMEMvField[0], ePartSize, int(uiPartAddr), 0, iPartIdx)
                pcCU.GetCUMvField(TLibCommon.REF_PIC_LIST_1).SetAllMvField(&cMEMvField[1], ePartSize, int(uiPartAddr), 0, iPartIdx)
            }
        }

        //  MC
        this.MotionCompensation(pcCU, rpcPredYuv, TLibCommon.REF_PIC_LIST_X, iPartIdx)

    }   //  end of for ( int iPartIdx = 0; iPartIdx < iNumPart; iPartIdx++ )

    this.setWpScalingDistParam(pcCU, -1, TLibCommon.REF_PIC_LIST_X)
	
    return
}

/// encode residual and compute rd-cost for inter mode
func (this *TEncSearch) encodeResAndCalcRdInterCU(pcCU *TLibCommon.TComDataCU,
    pcYuvOrg *TLibCommon.TComYuv,
    pcYuvPred *TLibCommon.TComYuv,
    rpcYuvResi *TLibCommon.TComYuv,
    rpcYuvResiBest *TLibCommon.TComYuv,
    rpcYuvRec *TLibCommon.TComYuv,
    bSkipRes bool) {
    if pcCU.IsIntra(0) {
        return
    }

    bHighPass := pcCU.GetSlice().GetDepth() != 0 // ? true : false;
    uiBits := uint(0)
    uiBitsBest := uint(0)
    uiDistortion := uint(0)
    uiDistortionBest := uint(0)

    uiWidth := uint(pcCU.GetWidth1(0))
    uiHeight := uint(pcCU.GetHeight1(0))

    //  No residual coding : SKIP mode
    if bSkipRes {
        pcCU.SetSkipFlagSubParts(true, 0, uint(pcCU.GetDepth1(0)))

        rpcYuvResi.Clear()

        pcYuvPred.CopyToPartYuv(rpcYuvRec, 0)

        //#if WEIGHTED_CHROMA_DISTORTION
        uiDistortion = this.m_pcRdCost.getDistPart(TLibCommon.G_bitDepthY, rpcYuvRec.GetLumaAddr(), int(rpcYuvRec.GetStride()), pcYuvOrg.GetLumaAddr(), int(pcYuvOrg.GetStride()), uiWidth, uiHeight, TLibCommon.TEXT_LUMA, TLibCommon.DF_SSE) +
            this.m_pcRdCost.getDistPart(TLibCommon.G_bitDepthC, rpcYuvRec.GetCbAddr(), int(rpcYuvRec.GetCStride()), pcYuvOrg.GetCbAddr(), int(pcYuvOrg.GetCStride()), uiWidth>>1, uiHeight>>1, TLibCommon.TEXT_CHROMA_U, TLibCommon.DF_SSE) +
            this.m_pcRdCost.getDistPart(TLibCommon.G_bitDepthC, rpcYuvRec.GetCrAddr(), int(rpcYuvRec.GetCStride()), pcYuvOrg.GetCrAddr(), int(pcYuvOrg.GetCStride()), uiWidth>>1, uiHeight>>1, TLibCommon.TEXT_CHROMA_V, TLibCommon.DF_SSE)
        //#else
        //    uiDistortion = this.m_pcRdCost.GetDistPart(TLibCommon.G_bitDepthY, rpcYuvRec.GetLumaAddr(), rpcYuvRec.GetStride(),  pcYuvOrg.GetLumaAddr(), pcYuvOrg.GetStride(),  uiWidth,      uiHeight      )
        //    + this.m_pcRdCost.GetDistPart(TLibCommon.G_bitDepthC, rpcYuvRec.GetCbAddr(),   rpcYuvRec.GetCStride(), pcYuvOrg.GetCbAddr(),   pcYuvOrg.GetCStride(), uiWidth >> 1, uiHeight >> 1 )
        //    + this.m_pcRdCost.GetDistPart(TLibCommon.G_bitDepthC, rpcYuvRec.GetCrAddr(),   rpcYuvRec.GetCStride(), pcYuvOrg.GetCrAddr(),   pcYuvOrg.GetCStride(), uiWidth >> 1, uiHeight >> 1 );
        //#endif

        if this.m_bUseSBACRD {
            this.m_pcRDGoOnSbacCoder.load(this.m_pppcRDSbacCoder[pcCU.GetDepth1(0)][TLibCommon.CI_CURR_BEST])
        }

        this.m_pcEntropyCoder.resetBits()
        if pcCU.GetSlice().GetPPS().GetTransquantBypassEnableFlag() {
            this.m_pcEntropyCoder.encodeCUTransquantBypassFlag(pcCU, 0, true)
        }
        this.m_pcEntropyCoder.encodeSkipFlag(pcCU, 0, true)
        this.m_pcEntropyCoder.encodeMergeIndex(pcCU, 0, true)

        uiBits = this.m_pcEntropyCoder.getNumberOfWrittenBits()
        pcCU.SetTotalBits(uiBits)
        pcCU.SetTotalDistortion(uiDistortion)
        pcCU.SetTotalCost(this.m_pcRdCost.calcRdCost(uiBits, uiDistortion, false, TLibCommon.DF_DEFAULT))

        if this.m_bUseSBACRD {
            this.m_pcRDGoOnSbacCoder.store(this.m_pppcRDSbacCoder[pcCU.GetDepth1(0)][TLibCommon.CI_TEMP_BEST])
        }
        pcCU.SetCbfSubParts(0, 0, 0, 0, uint(pcCU.GetDepth1(0)))
        pcCU.SetTrIdxSubParts(0, 0, uint(pcCU.GetDepth1(0)))

        return
    }

    //  Residual coding.
    var qp, qpBest, qpMin, qpMax int
    var dCost float64
    dCostBest := float64(TLibCommon.MAX_DOUBLE)

    uiTrLevel := uint(0)
    if uint(pcCU.GetWidth1(0)) > pcCU.GetSlice().GetSPS().GetMaxTrSize() {
        for uint(pcCU.GetWidth1(0)) > (pcCU.GetSlice().GetSPS().GetMaxTrSize() << uiTrLevel) {
            uiTrLevel++
        }
    }
    uiMaxTrMode := pcCU.GetSlice().GetSPS().GetMaxTrDepth() + uiTrLevel

    for (uiWidth >> uiMaxTrMode) < (this.m_pcEncCfg.GetMaxCUWidth()>>this.m_pcEncCfg.GetMaxCUDepth()) {
        uiMaxTrMode--
    }

    if bHighPass {
        qpMin = TLibCommon.CLIP3(int(-pcCU.GetSlice().GetSPS().GetQpBDOffsetY()), int(TLibCommon.MAX_QP), int(int(pcCU.GetQP1(0))-this.m_iMaxDeltaQP)).(int)
        qpMax = TLibCommon.CLIP3(int(-pcCU.GetSlice().GetSPS().GetQpBDOffsetY()), int(TLibCommon.MAX_QP), int(int(pcCU.GetQP1(0))+this.m_iMaxDeltaQP)).(int)
    } else {
        qpMin = int(pcCU.GetQP1(0))
        qpMax = int(pcCU.GetQP1(0))
    }

    rpcYuvResi.Subtract(pcYuvOrg, pcYuvPred, 0, uiWidth)

    for qp = qpMin; qp <= qpMax; qp++ {
        dCost = 0.
        uiBits = 0
        uiDistortion = 0
        if this.m_bUseSBACRD {
            this.m_pcRDGoOnSbacCoder.load(this.m_pppcRDSbacCoder[pcCU.GetDepth1(0)][TLibCommon.CI_CURR_BEST])
        }

        uiZeroDistortion := uint(0)
        this.xEstimateResidualQT(pcCU, 0, 0, 0, rpcYuvResi, uint(pcCU.GetDepth1(0)), &dCost, &uiBits, &uiDistortion, &uiZeroDistortion)

        this.m_pcEntropyCoder.resetBits()
        this.m_pcEntropyCoder.encodeQtRootCbfZero(pcCU)
        zeroResiBits := uint(this.m_pcEntropyCoder.getNumberOfWrittenBits())
        dZeroCost := this.m_pcRdCost.calcRdCost(zeroResiBits, uiZeroDistortion, false, TLibCommon.DF_DEFAULT)
        if pcCU.IsLosslessCoded(0) {
            dZeroCost = dCost + 1
        }
        if dZeroCost < dCost {
            dCost = dZeroCost
            uiBits = 0
            uiDistortion = uiZeroDistortion

            uiQPartNum := pcCU.GetPic().GetNumPartInCU() >> (pcCU.GetDepth1(0) << 1)

            for i := uint(0); i < uiQPartNum; i++ {
                pcCU.GetTransformIdx()[i] = 0                 //, uiQPartNum * sizeof(byte) );
                pcCU.GetCbf1(TLibCommon.TEXT_LUMA)[i] = 0     //, uiQPartNum * sizeof(byte) );
                pcCU.GetCbf1(TLibCommon.TEXT_CHROMA_U)[i] = 0 //, uiQPartNum * sizeof(byte) );
                pcCU.GetCbf1(TLibCommon.TEXT_CHROMA_V)[i] = 0 //, uiQPartNum * sizeof(byte) );
            }
            for i := uint(0); i < uiWidth*uiHeight; i++ {
                pcCU.GetCoeffY()[i] = 0 //, uiWidth * uiHeight * sizeof( TLibCommon.TCoeff )      );
            }
            for i := uint(0); i < (uiWidth*uiHeight)>>2; i++ {
                pcCU.GetCoeffCb()[i] = 0 //, uiWidth * uiHeight * sizeof( TLibCommon.TCoeff ) >> 2 );
                pcCU.GetCoeffCr()[i] = 0 //, uiWidth * uiHeight * sizeof( TLibCommon.TCoeff ) >> 2 );
            }
            pcCU.SetTransformSkipSubParts5(false, false, false, 0, uint(pcCU.GetDepth1(0)))
        } else {
            this.xSetResidualQTData(pcCU, 0, 0, 0, nil, uint(pcCU.GetDepth1(0)), false)
        }

        if this.m_bUseSBACRD {
            this.m_pcRDGoOnSbacCoder.load(this.m_pppcRDSbacCoder[pcCU.GetDepth1(0)][TLibCommon.CI_CURR_BEST])
        }
        /*#if 0 // check
            {
              this.m_pcEntropyCoder->resetBits();
              this.m_pcEntropyCoder->encodeCoeff( pcCU, 0, pcCU.GetDepth1(0), pcCU.GetWidth1(0), pcCU.GetHeight(0) );
              const UInt uiBitsForCoeff = this.m_pcEntropyCoder.getNumberOfWrittenBits();
              if( this.m_bUseSBACRD )
              {
                this.m_pcRDGoOnSbacCoder->load( this.m_pppcRDSbacCoder[pcCU.GetDepth1(0)][TLibCommon.CI_CURR_BEST] );
              }
              if( uiBitsForCoeff != uiBits )
                assert( 0 );
            }
        #endif*/
        uiBits = 0
        {
            var pDummy *TLibCommon.TComYuv
            this.xAddSymbolBitsInter(pcCU, 0, 0, &uiBits, pDummy, nil, pDummy)
        }

        dExactCost := this.m_pcRdCost.calcRdCost(uiBits, uiDistortion, false, TLibCommon.DF_DEFAULT)
        dCost = dExactCost

        if dCost < dCostBest {
            if !pcCU.GetQtRootCbf(0) {
                rpcYuvResiBest.Clear()
            } else {
                this.xSetResidualQTData(pcCU, 0, 0, 0, rpcYuvResiBest, uint(pcCU.GetDepth1(0)), true)
            }

            if qpMin != qpMax && qp != qpMax {
                uiQPartNum := pcCU.GetPic().GetNumPartInCU() >> (uint(pcCU.GetDepth1(0)) << 1)
                for i := uint(0); i < uiQPartNum; i++ {
                    this.m_puhQTTempTrIdx[i] = pcCU.GetTransformIdx()[i]                                          //,        uiQPartNum * sizeof(byte) );
                    this.m_puhQTTempCbf[0][i] = pcCU.GetCbf1(TLibCommon.TEXT_LUMA)[i]                             //,     uiQPartNum * sizeof(byte) );
                    this.m_puhQTTempCbf[1][i] = pcCU.GetCbf1(TLibCommon.TEXT_CHROMA_U)[i]                         //, uiQPartNum * sizeof(byte) );
                    this.m_puhQTTempCbf[2][i] = pcCU.GetCbf1(TLibCommon.TEXT_CHROMA_V)[i]                         //, uiQPartNum * sizeof(byte) );
                    this.m_puhQTTempTransformSkipFlag[0][i] = pcCU.GetTransformSkip1(TLibCommon.TEXT_LUMA)[i]     //,     uiQPartNum * sizeof( byte ) );
                    this.m_puhQTTempTransformSkipFlag[1][i] = pcCU.GetTransformSkip1(TLibCommon.TEXT_CHROMA_U)[i] //, uiQPartNum * sizeof( byte ) );
                    this.m_puhQTTempTransformSkipFlag[2][i] = pcCU.GetTransformSkip1(TLibCommon.TEXT_CHROMA_V)[i] //, uiQPartNum * sizeof( byte ) );
                }
                for i := uint(0); i < uiWidth*uiHeight; i++ {
                    this.m_pcQTTempCoeffY[i] = pcCU.GetCoeffY()[i]       //,  uiWidth * uiHeight * sizeof( TLibCommon.TCoeff )      );
                    this.m_pcQTTempArlCoeffY[i] = pcCU.GetArlCoeffY()[i] //,  uiWidth * uiHeight * sizeof( int )      );
                }
                for i := uint(0); i < (uiWidth*uiHeight)>>2; i++ {
                    this.m_pcQTTempCoeffCb[i] = pcCU.GetCoeffCb()[i] //, uiWidth * uiHeight * sizeof( TLibCommon.TCoeff ) >> 2 );
                    this.m_pcQTTempCoeffCr[i] = pcCU.GetCoeffCr()[i] //, uiWidth * uiHeight * sizeof( TLibCommon.TCoeff ) >> 2 );

                    this.m_pcQTTempArlCoeffCb[i] = pcCU.GetArlCoeffCb()[i] //, uiWidth * uiHeight * sizeof( int ) >> 2 );
                    this.m_pcQTTempArlCoeffCr[i] = pcCU.GetArlCoeffCr()[i] //, uiWidth * uiHeight * sizeof( int ) >> 2 );
                }
                //#if ADAPTIVE_QP_SELECTION

                //#endif

            }
            uiBitsBest = uiBits
            uiDistortionBest = uiDistortion
            dCostBest = dCost
            qpBest = qp
            if this.m_bUseSBACRD {
                this.m_pcRDGoOnSbacCoder.store(this.m_pppcRDSbacCoder[uint(pcCU.GetDepth1(0))][TLibCommon.CI_TEMP_BEST])
            }
        }
    }

    //assert ( dCostBest != MAX_DOUBLE );

    if qpMin != qpMax && qpBest != qpMax {
        if this.m_bUseSBACRD {
            //assert( 0 ); // check
            this.m_pcRDGoOnSbacCoder.load(this.m_pppcRDSbacCoder[uint(pcCU.GetDepth1(0))][TLibCommon.CI_TEMP_BEST])
        }
        // copy best cbf and trIdx to pcCU
        uiQPartNum := pcCU.GetPic().GetNumPartInCU() >> (uint(pcCU.GetDepth1(0)) << 1)
        for i := uint(0); i < uiQPartNum; i++ {
            pcCU.GetTransformIdx()[i] = this.m_puhQTTempTrIdx[i]                                          //,  uiQPartNum * sizeof(byte) );
            pcCU.GetCbf1(TLibCommon.TEXT_LUMA)[i] = this.m_puhQTTempCbf[0][i]                             //, uiQPartNum * sizeof(byte) );
            pcCU.GetCbf1(TLibCommon.TEXT_CHROMA_U)[i] = this.m_puhQTTempCbf[1][i]                         //, uiQPartNum * sizeof(byte) );
            pcCU.GetCbf1(TLibCommon.TEXT_CHROMA_V)[i] = this.m_puhQTTempCbf[2][i]                         //, uiQPartNum * sizeof(byte) );
            pcCU.GetTransformSkip1(TLibCommon.TEXT_LUMA)[i] = this.m_puhQTTempTransformSkipFlag[0][i]     //, uiQPartNum * sizeof( byte ) );
            pcCU.GetTransformSkip1(TLibCommon.TEXT_CHROMA_U)[i] = this.m_puhQTTempTransformSkipFlag[1][i] //, uiQPartNum * sizeof( byte ) );
            pcCU.GetTransformSkip1(TLibCommon.TEXT_CHROMA_V)[i] = this.m_puhQTTempTransformSkipFlag[2][i] //, uiQPartNum * sizeof( byte ) );
        }
        for i := uint(0); i < uiWidth*uiHeight; i++ {
            pcCU.GetCoeffY()[i] = this.m_pcQTTempCoeffY[i]       //,  uiWidth * uiHeight * sizeof( TLibCommon.TCoeff )      );
            pcCU.GetArlCoeffY()[i] = this.m_pcQTTempArlCoeffY[i] //,  uiWidth * uiHeight * sizeof( int )      );
        }
        for i := uint(0); i < (uiWidth*uiHeight)>>2; i++ {
            pcCU.GetCoeffCb()[i] = this.m_pcQTTempCoeffCb[i] //, uiWidth * uiHeight * sizeof( TLibCommon.TCoeff ) >> 2 );
            pcCU.GetCoeffCr()[i] = this.m_pcQTTempCoeffCr[i] //, uiWidth * uiHeight * sizeof( TLibCommon.TCoeff ) >> 2 );
            //#if ADAPTIVE_QP_SELECTION
            pcCU.GetArlCoeffCb()[i] = this.m_pcQTTempArlCoeffCb[i] //, uiWidth * uiHeight * sizeof( int ) >> 2 );
            pcCU.GetArlCoeffCr()[i] = this.m_pcQTTempArlCoeffCr[i] //, uiWidth * uiHeight * sizeof( int ) >> 2 );
            //#endif
        }

    }
    rpcYuvRec.AddClip(pcYuvPred, rpcYuvResiBest, 0, uiWidth)

    // update with clipped distortion and cost (qp estimation loop uses unclipped values)
    //#if WEIGHTED_CHROMA_DISTORTION
    uiDistortionBest = this.m_pcRdCost.getDistPart(TLibCommon.G_bitDepthY, rpcYuvRec.GetLumaAddr(), int(rpcYuvRec.GetStride()), pcYuvOrg.GetLumaAddr(), int(pcYuvOrg.GetStride()), uiWidth, uiHeight, TLibCommon.TEXT_LUMA, TLibCommon.DF_SSE) +
        this.m_pcRdCost.getDistPart(TLibCommon.G_bitDepthC, rpcYuvRec.GetCbAddr(), int(rpcYuvRec.GetCStride()), pcYuvOrg.GetCbAddr(), int(pcYuvOrg.GetCStride()), uiWidth>>1, uiHeight>>1, TLibCommon.TEXT_CHROMA_U, TLibCommon.DF_SSE) +
        this.m_pcRdCost.getDistPart(TLibCommon.G_bitDepthC, rpcYuvRec.GetCrAddr(), int(rpcYuvRec.GetCStride()), pcYuvOrg.GetCrAddr(), int(pcYuvOrg.GetCStride()), uiWidth>>1, uiHeight>>1, TLibCommon.TEXT_CHROMA_V, TLibCommon.DF_SSE)
    /*#else
      uiDistortionBest = this.m_pcRdCost.getDistPart(TLibCommon.G_bitDepthY, rpcYuvRec.GetLumaAddr(), rpcYuvRec.GetStride(),  pcYuvOrg.GetLumaAddr(), pcYuvOrg.GetStride(),  uiWidth,      uiHeight      )
      + this.m_pcRdCost.getDistPart(TLibCommon.G_bitDepthC, rpcYuvRec.GetCbAddr(),   rpcYuvRec.GetCStride(), pcYuvOrg.GetCbAddr(),   pcYuvOrg.GetCStride(), uiWidth >> 1, uiHeight >> 1 )
      + this.m_pcRdCost.getDistPart(TLibCommon.G_bitDepthC, rpcYuvRec.GetCrAddr(),   rpcYuvRec.GetCStride(), pcYuvOrg.GetCrAddr(),   pcYuvOrg.GetCStride(), uiWidth >> 1, uiHeight >> 1 );
    #endif*/
    dCostBest = this.m_pcRdCost.calcRdCost(uiBitsBest, uiDistortionBest, false, TLibCommon.DF_DEFAULT)

    pcCU.SetTotalBits(uiBitsBest)
    pcCU.SetTotalDistortion(uiDistortionBest)
    pcCU.SetTotalCost(dCostBest)

    if pcCU.IsSkipped(0) {
        pcCU.SetCbfSubParts(0, 0, 0, 0, uint(pcCU.GetDepth1(0)))
    }

    pcCU.SetQPSubParts(qpBest, 0, uint(pcCU.GetDepth1(0)))
}

/// set ME search range
func (this *TEncSearch) setAdaptiveSearchRange(iDir, iRefIdx, iSearchRange int) {
    this.m_aaiAdaptSR[iDir][iRefIdx] = iSearchRange
}

func (this *TEncSearch) xEncPCM(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint, piOrg []TLibCommon.Pel, piPCM []TLibCommon.Pel, piPred []TLibCommon.Pel, piResi []TLibCommon.Pel, piReco []TLibCommon.Pel, uiStride, uiWidth, uiHeight uint, eText TLibCommon.TextType) {
    var uiX, uiY, uiReconStride uint
    pOrg := piOrg
    pPCM := piPCM
    pPred := piPred
    pResi := piResi
    pReco := piReco
    var pRecoPic []TLibCommon.Pel
    var shiftPcm uint

    if eText == TLibCommon.TEXT_LUMA {
        uiReconStride = uint(pcCU.GetPic().GetPicYuvRec().GetStride())
        pRecoPic = pcCU.GetPic().GetPicYuvRec().GetLumaAddr2(int(pcCU.GetAddr()), int(pcCU.GetZorderIdxInCU()+uiAbsPartIdx))
        shiftPcm = uint(TLibCommon.G_bitDepthY) - pcCU.GetSlice().GetSPS().GetPCMBitDepthLuma()
    } else {
        uiReconStride = uint(pcCU.GetPic().GetPicYuvRec().GetCStride())

        if eText == TLibCommon.TEXT_CHROMA_U {
            pRecoPic = pcCU.GetPic().GetPicYuvRec().GetCbAddr2(int(pcCU.GetAddr()), int(pcCU.GetZorderIdxInCU()+uiAbsPartIdx))
        } else {
            pRecoPic = pcCU.GetPic().GetPicYuvRec().GetCrAddr2(int(pcCU.GetAddr()), int(pcCU.GetZorderIdxInCU()+uiAbsPartIdx))
        }
        shiftPcm = uint(TLibCommon.G_bitDepthC) - pcCU.GetSlice().GetSPS().GetPCMBitDepthChroma()
    }

    // Reset pred and residual
    for uiY = 0; uiY < uiHeight; uiY++ {
        for uiX = 0; uiX < uiWidth; uiX++ {
            pPred[uiY*uiStride+uiX] = 0
            pResi[uiY*uiStride+uiX] = 0
        }
        //pPred = pPred[uiStride:]
        //pResi = pResi[uiStride:]
    }

    // Encode
    for uiY = 0; uiY < uiHeight; uiY++ {
        for uiX = 0; uiX < uiWidth; uiX++ {
            pPCM[uiY*uiWidth+uiX] = pOrg[uiY*uiStride+uiX] >> shiftPcm
        }
        //pPCM = pPCM[uiWidth:]
        //pOrg = pOrg[uiStride:]
    }

    //pPCM = piPCM

    // Reconstruction
    for uiY = 0; uiY < uiHeight; uiY++ {
        for uiX = 0; uiX < uiWidth; uiX++ {
            pReco[uiY*uiStride+uiX] = pPCM[uiY*uiWidth+uiX] << shiftPcm
            pRecoPic[uiY*uiReconStride+uiX] = pReco[uiY*uiStride+uiX]
        }
        //pPCM = pPCM[uiWidth:]
        //pReco = pReco[uiStride:]
        //pRecoPic = pRecoPic[uiReconStride:]
    }
}

func (this *TEncSearch) IPCMSearch(pcCU *TLibCommon.TComDataCU, pcOrgYuv *TLibCommon.TComYuv, rpcPredYuv *TLibCommon.TComYuv, rpcResiYuv *TLibCommon.TComYuv, rpcRecoYuv *TLibCommon.TComYuv) {
    uiDepth := uint(pcCU.GetDepth1(0))
    uiWidth := uint(pcCU.GetWidth1(0))
    uiHeight := uint(pcCU.GetHeight1(0))
    uiStride := rpcPredYuv.GetStride()
    uiStrideC := rpcPredYuv.GetCStride()
    uiWidthC := uiWidth >> 1
    uiHeightC := uiHeight >> 1
    uiDistortion := uint(0)
    var uiBits uint

    var dCost float64

    var pOrig, pResi, pReco, pPred, pPCM []TLibCommon.Pel

    uiAbsPartIdx := uint(0)

    uiMinCoeffSize := pcCU.GetPic().GetMinCUWidth() * pcCU.GetPic().GetMinCUHeight()
    uiLumaOffset := uiMinCoeffSize * uiAbsPartIdx
    uiChromaOffset := uiLumaOffset >> 2

    // Luminance
    pOrig = pcOrgYuv.GetLumaAddr2(0, uiWidth)
    pResi = rpcResiYuv.GetLumaAddr2(0, uiWidth)
    pPred = rpcPredYuv.GetLumaAddr2(0, uiWidth)
    pReco = rpcRecoYuv.GetLumaAddr2(0, uiWidth)
    pPCM = pcCU.GetPCMSampleY()[uiLumaOffset:]

    this.xEncPCM(pcCU, 0, pOrig, pPCM, pPred, pResi, pReco, uiStride, uiWidth, uiHeight, TLibCommon.TEXT_LUMA)

    // Chroma U
    pOrig = pcOrgYuv.GetCbAddr()
    pResi = rpcResiYuv.GetCbAddr()
    pPred = rpcPredYuv.GetCbAddr()
    pReco = rpcRecoYuv.GetCbAddr()
    pPCM = pcCU.GetPCMSampleCb()[uiChromaOffset:]

    this.xEncPCM(pcCU, 0, pOrig, pPCM, pPred, pResi, pReco, uiStrideC, uiWidthC, uiHeightC, TLibCommon.TEXT_CHROMA_U)

    // Chroma V
    pOrig = pcOrgYuv.GetCrAddr()
    pResi = rpcResiYuv.GetCrAddr()
    pPred = rpcPredYuv.GetCrAddr()
    pReco = rpcRecoYuv.GetCrAddr()
    pPCM = pcCU.GetPCMSampleCr()[uiChromaOffset:]

    this.xEncPCM(pcCU, 0, pOrig, pPCM, pPred, pResi, pReco, uiStrideC, uiWidthC, uiHeightC, TLibCommon.TEXT_CHROMA_V)

    this.m_pcEntropyCoder.resetBits()
    this.xEncIntraHeader(pcCU, uiDepth, uiAbsPartIdx, true, false)
    uiBits = this.m_pcEntropyCoder.getNumberOfWrittenBits()

    dCost = this.m_pcRdCost.calcRdCost(uiBits, uiDistortion, false, TLibCommon.DF_DEFAULT)

    if this.m_bUseSBACRD {
        this.m_pcRDGoOnSbacCoder.load(this.m_pppcRDSbacCoder[uiDepth][TLibCommon.CI_CURR_BEST])
    }

    pcCU.SetTotalBits(uiBits)
    pcCU.SetTotalCost(dCost)
    pcCU.SetTotalDistortion(uiDistortion)

    pcCU.CopyToPic3(uiDepth, 0, 0)
}

// -------------------------------------------------------------------------------------------------------------------
// Intra search
// -------------------------------------------------------------------------------------------------------------------

func (this *TEncSearch) xEncSubdivCbfQT(pcCU *TLibCommon.TComDataCU,
    uiTrDepth uint,
    uiAbsPartIdx uint,
    bLuma bool,
    bChroma bool) {
    //fmt.Printf("Enter xEncSubdivCbfQT\n")
    uiFullDepth := uint(pcCU.GetDepth1(0)) + uiTrDepth
    uiTrMode := uint(pcCU.GetTransformIdx1(uiAbsPartIdx))
    uiSubdiv := uint(TLibCommon.B2U(uiTrMode > uiTrDepth))
    uiLog2TrafoSize := uint(TLibCommon.G_aucConvertToBit[pcCU.GetSlice().GetSPS().GetMaxCUWidth()]) + 2 - uiFullDepth

    {
        if pcCU.GetPredictionMode1(0) == TLibCommon.MODE_INTRA && pcCU.GetPartitionSize1(0) == TLibCommon.SIZE_NxN && uiTrDepth == 0 {
            //assert(uiSubdiv)
        } else if uiLog2TrafoSize > pcCU.GetSlice().GetSPS().GetQuadtreeTULog2MaxSize() {
            //assert( uiSubdiv );
        } else if uiLog2TrafoSize == pcCU.GetSlice().GetSPS().GetQuadtreeTULog2MinSize() {
            //assert( !uiSubdiv );
        } else if uiLog2TrafoSize == pcCU.GetQuadtreeTULog2MinSizeInCU(uiAbsPartIdx) {
            //assert( !uiSubdiv );
        } else {
            //assert( uiLog2TrafoSize > pcCU.GetQuadtreeTULog2MinSizeInCU(uiAbsPartIdx) );
            if bLuma {
                this.m_pcEntropyCoder.encodeTransformSubdivFlag(uiSubdiv, 5-uiLog2TrafoSize)
            }
        }
    }

    if bChroma {
        if uiLog2TrafoSize > 2 {
            if uiTrDepth == 0 || pcCU.GetCbf3(uiAbsPartIdx, TLibCommon.TEXT_CHROMA_U, uiTrDepth-1) != 0 {
                this.m_pcEntropyCoder.encodeQtCbf(pcCU, uiAbsPartIdx, TLibCommon.TEXT_CHROMA_U, uiTrDepth)
            }
            if uiTrDepth == 0 || pcCU.GetCbf3(uiAbsPartIdx, TLibCommon.TEXT_CHROMA_V, uiTrDepth-1) != 0 {
                this.m_pcEntropyCoder.encodeQtCbf(pcCU, uiAbsPartIdx, TLibCommon.TEXT_CHROMA_V, uiTrDepth)
            }
        }
    }

    if uiSubdiv != 0 {
        uiQPartNum := pcCU.GetPic().GetNumPartInCU() >> ((uiFullDepth + 1) << 1)
        for uiPart := uint(0); uiPart < 4; uiPart++ {
            this.xEncSubdivCbfQT(pcCU, uiTrDepth+1, uiAbsPartIdx+uiPart*uiQPartNum, bLuma, bChroma)
        }
        //fmt.Printf("Exit xEncSubdivCbfQT\n");
        return
    }

    {
        //===== Cbfs =====
        if bLuma {
        	//fmt.Printf("%v\n", this.m_pcEntropyCoder.m_pcEntropyCoderIf)
            this.m_pcEntropyCoder.encodeQtCbf(pcCU, uiAbsPartIdx, TLibCommon.TEXT_LUMA, uiTrMode)
        }
    }
    //fmt.Printf("Exit xEncSubdivCbfQT\n");
}

func (this *TEncSearch) xEncCoeffQT(pcCU *TLibCommon.TComDataCU,
    uiTrDepth uint,
    uiAbsPartIdx uint,
    eTextType TLibCommon.TextType,
    bRealCoeff bool) {
    uiFullDepth := uint(pcCU.GetDepth1(0)) + uiTrDepth
    uiTrMode := uint(pcCU.GetTransformIdx1(uiAbsPartIdx))
    uiSubdiv := uint(TLibCommon.B2U(uiTrMode > uiTrDepth))
    uiLog2TrafoSize := uint(TLibCommon.G_aucConvertToBit[pcCU.GetSlice().GetSPS().GetMaxCUWidth()]) + 2 - uiFullDepth
    uiChroma := uint(TLibCommon.B2U(eTextType != TLibCommon.TEXT_LUMA))
	
	//fmt.Printf("Enter xEncCoeffQT (%d,%d,%d,%d)\n",uiTrDepth,uiAbsPartIdx,eTextType,TLibCommon.B2U(bRealCoeff));
  
	
    if uiSubdiv != 0 {
        uiQPartNum := pcCU.GetPic().GetNumPartInCU() >> ((uiFullDepth + 1) << 1)
        for uiPart := uint(0); uiPart < 4; uiPart++ {
            this.xEncCoeffQT(pcCU, uiTrDepth+1, uiAbsPartIdx+uiPart*uiQPartNum, eTextType, bRealCoeff)
        }
        //fmt.Printf("Exit xEncCoeffQT\n");
        return
    }

    if eTextType != TLibCommon.TEXT_LUMA && uiLog2TrafoSize == 2 {
        //assert( uiTrDepth > 0 );
        uiTrDepth--
        uiQPDiv := pcCU.GetPic().GetNumPartInCU() >> ((uint(pcCU.GetDepth1(0)) + uiTrDepth) << 1)
        bFirstQ := ((uiAbsPartIdx % uiQPDiv) == 0)
        if !bFirstQ {
        	//fmt.Printf("Exit xEncCoeffQT\n");
            return
        }
    }

    //===== coefficients =====
    uiWidth := uint(uint(pcCU.GetWidth1(0))) >> (uiTrDepth + uiChroma)
    uiHeight := uint(pcCU.GetHeight1(0)) >> (uiTrDepth + uiChroma)
    uiCoeffOffset := (pcCU.GetPic().GetMinCUWidth() * pcCU.GetPic().GetMinCUHeight() * uiAbsPartIdx) >> (uiChroma << 1)
    uiQTLayer := pcCU.GetSlice().GetSPS().GetQuadtreeTULog2MaxSize() - uiLog2TrafoSize
    var pcCoeff []TLibCommon.TCoeff
    switch eTextType {
    case TLibCommon.TEXT_LUMA:
        if bRealCoeff {
            pcCoeff = pcCU.GetCoeffY()
        } else {
            pcCoeff = this.m_ppcQTTempCoeffY[uiQTLayer]
        }
    case TLibCommon.TEXT_CHROMA_U:
        if bRealCoeff {
            pcCoeff = pcCU.GetCoeffCb()
        } else {
            pcCoeff = this.m_ppcQTTempCoeffCb[uiQTLayer]
        }
    case TLibCommon.TEXT_CHROMA_V:
        if bRealCoeff {
            pcCoeff = pcCU.GetCoeffCr()
        } else {
            pcCoeff = this.m_ppcQTTempCoeffCr[uiQTLayer]
        }
        //default:            assert(0);
    }
    pcCoeff = pcCoeff[uiCoeffOffset:]

    this.m_pcEntropyCoder.encodeCoeffNxN(pcCU, pcCoeff, uiAbsPartIdx, uiWidth, uiHeight, uiFullDepth, eTextType)
	
	//fmt.Printf("Exit xEncCoeffQT with uiBits=%d\n", this.m_pcEntropyCoder.getNumberOfWrittenBits());
}

func (this *TEncSearch) xEncIntraHeader(pcCU *TLibCommon.TComDataCU,
    uiTrDepth uint,
    uiAbsPartIdx uint,
    bLuma bool,
    bChroma bool) {
    if bLuma {
        // CU header
        if uiAbsPartIdx == 0 {
            if !pcCU.GetSlice().IsIntra() {
                if pcCU.GetSlice().GetPPS().GetTransquantBypassEnableFlag() {
                    this.m_pcEntropyCoder.encodeCUTransquantBypassFlag(pcCU, 0, true)
                }
                this.m_pcEntropyCoder.encodeSkipFlag(pcCU, 0, true)
                this.m_pcEntropyCoder.encodePredMode(pcCU, 0, true)
            }

            this.m_pcEntropyCoder.encodePartSize(pcCU, 0, uint(pcCU.GetDepth1(0)), true)

            if pcCU.IsIntra(0) && pcCU.GetPartitionSize1(0) == TLibCommon.SIZE_2Nx2N {
                this.m_pcEntropyCoder.encodeIPCMInfo(pcCU, 0, true)

                if pcCU.GetIPCMFlag1(0) {
                    return
                }
            }
        }
        // luma prediction mode
        if pcCU.GetPartitionSize1(0) == TLibCommon.SIZE_2Nx2N {
            if uiAbsPartIdx == 0 {
                this.m_pcEntropyCoder.encodeIntraDirModeLuma(pcCU, 0, false)
            }
        } else {
            uiQNumParts := pcCU.GetTotalNumPart() >> 2
            if uiTrDepth == 0 {
                //assert( uiAbsPartIdx == 0 );
                for uiPart := uint(0); uiPart < 4; uiPart++ {
                    this.m_pcEntropyCoder.encodeIntraDirModeLuma(pcCU, uiPart*uiQNumParts, false)
                }
            } else if (uiAbsPartIdx % uiQNumParts) == 0 {
                this.m_pcEntropyCoder.encodeIntraDirModeLuma(pcCU, uiAbsPartIdx, false)
            }
        }
    }
    if bChroma {
        // chroma prediction mode
        if uiAbsPartIdx == 0 {
            this.m_pcEntropyCoder.encodeIntraDirModeChroma(pcCU, 0, true)
        }
    }
}

func (this *TEncSearch) xGetIntraBitsQT(pcCU *TLibCommon.TComDataCU,
    uiTrDepth uint,
    uiAbsPartIdx uint,
    bLuma bool,
    bChroma bool,
    bRealCoeff bool) uint {
    //fmt.Printf("Enter xGetIntraBitsQT\n");
    this.m_pcEntropyCoder.resetBits()
    //fmt.Printf("uiBits0=%d ", this.m_pcEntropyCoder.getNumberOfWrittenBits());
    this.xEncIntraHeader(pcCU, uiTrDepth, uiAbsPartIdx, bLuma, bChroma)
    //fmt.Printf("uiBits1=%d ", this.m_pcEntropyCoder.getNumberOfWrittenBits());
    this.xEncSubdivCbfQT(pcCU, uiTrDepth, uiAbsPartIdx, bLuma, bChroma)
	//fmt.Printf("uiBits2=%d ", this.m_pcEntropyCoder.getNumberOfWrittenBits());
	
    if bLuma {
        this.xEncCoeffQT(pcCU, uiTrDepth, uiAbsPartIdx, TLibCommon.TEXT_LUMA, bRealCoeff)
    }
    if bChroma {
        this.xEncCoeffQT(pcCU, uiTrDepth, uiAbsPartIdx, TLibCommon.TEXT_CHROMA_U, bRealCoeff)
        this.xEncCoeffQT(pcCU, uiTrDepth, uiAbsPartIdx, TLibCommon.TEXT_CHROMA_V, bRealCoeff)
    }
    uiBits := this.m_pcEntropyCoder.getNumberOfWrittenBits()
    //fmt.Printf("Exit xGetIntraBitsQT\n");
    
    return uiBits
}

func (this *TEncSearch) xGetIntraBitsQTChroma(pcCU *TLibCommon.TComDataCU,
    uiTrDepth uint,
    uiAbsPartIdx uint,
    uiChromaId uint,
    bRealCoeff bool) uint {
    this.m_pcEntropyCoder.resetBits()
    if uiChromaId == TLibCommon.TEXT_CHROMA_U {
        this.xEncCoeffQT(pcCU, uiTrDepth, uiAbsPartIdx, TLibCommon.TEXT_CHROMA_U, bRealCoeff)
    } else if uiChromaId == TLibCommon.TEXT_CHROMA_V {
        this.xEncCoeffQT(pcCU, uiTrDepth, uiAbsPartIdx, TLibCommon.TEXT_CHROMA_V, bRealCoeff)
    }

    uiBits := this.m_pcEntropyCoder.getNumberOfWrittenBits()
    return uiBits
}

func (this *TEncSearch) xIntraCodingLumaBlk(pcCU *TLibCommon.TComDataCU,
    uiTrDepth uint,
    uiAbsPartIdx uint,
    pcOrgYuv *TLibCommon.TComYuv,
    pcPredYuv *TLibCommon.TComYuv,
    pcResiYuv *TLibCommon.TComYuv,
    ruiDist *uint,
    default0Save1Load2 int) {
    //fmt.Printf("Enter xIntraCodingLumaBlk\n");
    
    uiLumaPredMode := uint(pcCU.GetLumaIntraDir1(uiAbsPartIdx))
    uiFullDepth := uint(pcCU.GetDepth1(0)) + uiTrDepth
    uiWidth := uint(uint(pcCU.GetWidth1(0))) >> uiTrDepth
    uiHeight := uint(pcCU.GetHeight1(0)) >> uiTrDepth
    uiStride := uint(pcOrgYuv.GetStride())
    piOrg := pcOrgYuv.GetLumaAddr1(uiAbsPartIdx)
    piPred := pcPredYuv.GetLumaAddr1(uiAbsPartIdx)
    piResi := pcResiYuv.GetLumaAddr1(uiAbsPartIdx)
    piReco := pcPredYuv.GetLumaAddr1(uiAbsPartIdx)

    uiLog2TrSize := uint(TLibCommon.G_aucConvertToBit[pcCU.GetSlice().GetSPS().GetMaxCUWidth()>>uiFullDepth]) + 2
    uiQTLayer := pcCU.GetSlice().GetSPS().GetQuadtreeTULog2MaxSize() - uiLog2TrSize
    uiNumCoeffPerInc := pcCU.GetSlice().GetSPS().GetMaxCUWidth() * pcCU.GetSlice().GetSPS().GetMaxCUHeight() >> (pcCU.GetSlice().GetSPS().GetMaxCUDepth() << 1)
    pcCoeff := this.m_ppcQTTempCoeffY[uiQTLayer][uiNumCoeffPerInc*uiAbsPartIdx:]
    //#if ADAPTIVE_QP_SELECTION
    pcArlCoeff := this.m_ppcQTTempArlCoeffY[uiQTLayer][uiNumCoeffPerInc*uiAbsPartIdx:]
    //#endif
    piRecQt := this.m_pcQTTempTComYuv[uiQTLayer].GetLumaAddr1(uiAbsPartIdx)
    uiRecQtStride := uint(this.m_pcQTTempTComYuv[uiQTLayer].GetStride())

    uiZOrder := pcCU.GetZorderIdxInCU() + uiAbsPartIdx
    piRecIPred := pcCU.GetPic().GetPicYuvRec().GetLumaAddr2(int(pcCU.GetAddr()), int(uiZOrder))
    uiRecIPredStride := uint(pcCU.GetPic().GetPicYuvRec().GetStride())
    useTransformSkip := pcCU.GetTransformSkip2(uiAbsPartIdx, TLibCommon.TEXT_LUMA)
    //===== init availability pattern =====
    bAboveAvail := false
    bLeftAvail := false
    if default0Save1Load2 != 2 {
        pcCU.GetPattern().InitPattern3(pcCU, uiTrDepth, uiAbsPartIdx)
        pcCU.GetPattern().InitAdiPattern(pcCU, uiAbsPartIdx, uiTrDepth, this.GetYuvExt(), this.GetYuvExtStride(), this.GetYuvExtHeight(), &bAboveAvail, &bLeftAvail, false)
        //===== get prediction signal =====
        this.PredIntraLumaAng(pcCU.GetPattern(), uiLumaPredMode, piPred, uiStride, int(uiWidth), int(uiHeight), bAboveAvail, bLeftAvail)
        // save prediction
        if default0Save1Load2 == 1 {
            pPred := piPred
            pPredBuf := this.m_pSharedPredTransformSkip[0]
            k := 0
            for uiY := uint(0); uiY < uiHeight; uiY++ {
                for uiX := uint(0); uiX < uiWidth; uiX++ {
                    pPredBuf[k] = pPred[uiY*uiStride+uiX]
                    k++
                }
                //pPred = pPred[uiStride:]
            }
        }
    } else {
        // load prediction
        pPred := piPred
        pPredBuf := this.m_pSharedPredTransformSkip[0]
        k := 0
        for uiY := uint(0); uiY < uiHeight; uiY++ {
            for uiX := uint(0); uiX < uiWidth; uiX++ {
                pPred[uiY*uiStride+uiX] = pPredBuf[k]
                k++
            }
            //pPred = pPred[uiStride:]
        }
    }
    //===== get residual signal =====
    {
        // get residual
        pOrg := piOrg
        pPred := piPred
        pResi := piResi
        for uiY := uint(0); uiY < uiHeight; uiY++ {
            for uiX := uint(0); uiX < uiWidth; uiX++ {
                pResi[uiY*uiStride+uiX] = pOrg[uiY*uiStride+uiX] - pPred[uiY*uiStride+uiX]
            	//fmt.Printf("%d-%d ", pOrg[ uiY*uiStride+uiX ], pPred[ uiY*uiStride+uiX ]);
            }
            //fmt.Printf("\n");
            //pOrg = pOrg[uiStride:]
            //pResi = pResi[uiStride:]
            //pPred = pPred[uiStride:]
        }
        //fmt.Printf("\n");
    }

    //===== transform and quantization =====
    //--- init rate estimation arrays for RDOQ ---
    var rdoqflag bool
    if useTransformSkip {
        rdoqflag = this.m_pcEncCfg.GetUseRDOQTS()
    } else {
        rdoqflag = this.m_pcEncCfg.GetUseRDOQ()
    }
    if rdoqflag {
        this.m_pcEntropyCoder.estimateBit(this.m_pcTrQuant.GetEstBitsSbac(), int(uiWidth), int(uiWidth), TLibCommon.TEXT_LUMA)
    }
    //--- transform and quantization ---
    uiAbsSum := uint(0)
    pcCU.SetTrIdxSubParts(uiTrDepth, uiAbsPartIdx, uiFullDepth)

    this.m_pcTrQuant.SetQPforQuant(int(pcCU.GetQP1(0)), TLibCommon.TEXT_LUMA, pcCU.GetSlice().GetSPS().GetQpBDOffsetY(), 0)

    //#if RDOQ_CHROMA_LAMBDA
    this.m_pcTrQuant.SelectLambda(TLibCommon.TEXT_LUMA)
    //#endif

    this.m_pcTrQuant.TransformNxN(pcCU, piResi, uiStride, pcCoeff,
        //#if ADAPTIVE_QP_SELECTION
        pcArlCoeff,
        //#endif
        uiWidth, uiHeight, &uiAbsSum, TLibCommon.TEXT_LUMA, uiAbsPartIdx, useTransformSkip)

    //--- set coded block flag ---
    pcCU.SetCbfSubParts4(byte(TLibCommon.B2U(uiAbsSum != 0))<<uiTrDepth, TLibCommon.TEXT_LUMA, uiAbsPartIdx, uiFullDepth)
    //--- inverse transform ---
    if uiAbsSum != 0 {
        scalingListType := 0 + int(TLibCommon.G_eTTable[TLibCommon.TEXT_LUMA])
        //assert(scalingListType < 6);
        
        /*#ifdef ENC_DEC_TRACE*/
	    this.m_pcEntropyCoder.m_pcEntropyCoderIf.XTraceCoefHeader(TLibCommon.TRACE_COEF)
	
	    for k := uint(0); k < uiHeight; k++ {
	        this.m_pcEntropyCoder.m_pcEntropyCoderIf.XReadCeofTr(pcCoeff[k*uiWidth:], uiWidth, TLibCommon.TRACE_COEF)
	    }
	    /*#endif*/
        
        this.m_pcTrQuant.InvtransformNxN(pcCU.GetCUTransquantBypass1(uiAbsPartIdx), TLibCommon.TEXT_LUMA, uint(pcCU.GetLumaIntraDir1(uiAbsPartIdx)), piResi, uiStride, pcCoeff, uiWidth, uiHeight, scalingListType, useTransformSkip)
    	
    	/*#ifdef ENC_DEC_TRACE*/
	    this.m_pcEntropyCoder.m_pcEntropyCoderIf.XTraceResiHeader(TLibCommon.TRACE_RESI)
	
	    for k := uint(0); k < uiHeight; k++ {
	        this.m_pcEntropyCoder.m_pcEntropyCoderIf.XReadResiTr(piResi[k*uiStride:], uiWidth, TLibCommon.TRACE_RESI)
	    }
	    /*#endif*/
    } else {
        pResi := piResi
        for i := uint(0); i < uiWidth*uiHeight; i++ {
            pcCoeff[i] = 0 //, sizeof( TLibCommon.TCoeff ) * uiWidth * uiHeight );
        }
        for uiY := uint(0); uiY < uiHeight; uiY++ {
            for uiX := uint(0); uiX < uiWidth; uiX++ {
                pResi[uiY*uiStride+uiX] = 0 //, sizeof( TLibCommon.Pel ) * uiWidth );
            }
            //pResi = pResi[uiStride:]
        }
    }

    //===== reconstruction =====
    {
        pPred := piPred
        pResi := piResi
        pReco := piReco
        pRecQt := piRecQt
        pRecIPred := piRecIPred
        for uiY := uint(0); uiY < uiHeight; uiY++ {
            for uiX := uint(0); uiX < uiWidth; uiX++ {
            	//fmt.Printf("(%d, %d) ", pPred[ uiY*uiStride+uiX ], pResi[ uiY*uiStride+uiX ]);
                pReco[uiY*uiStride+uiX] = TLibCommon.ClipY(pPred[uiY*uiStride+uiX] + pResi[uiY*uiStride+uiX])
                pRecQt[uiY*uiRecQtStride+uiX] = pReco[uiY*uiStride+uiX]
                pRecIPred[uiY*uiRecIPredStride+uiX] = pReco[uiY*uiStride+uiX]            	
            }
            //fmt.Printf("\n");
            //pPred = pPred[uiStride:]
            //pResi = pResi[uiStride:]
            //pReco = pReco[uiStride:]
            //pRecQt = pRecQt[uiRecQtStride:]
            //pRecIPred = pRecIPred[uiRecIPredStride:]
        }
        //fmt.Printf("\n");
    }

    //===== update distortion =====
    *ruiDist += this.m_pcRdCost.getDistPart(TLibCommon.G_bitDepthY, piReco, int(uiStride), piOrg, int(uiStride), uiWidth, uiHeight, TLibCommon.TEXT_LUMA, TLibCommon.DF_SSE)
	
	//fmt.Printf("Exit xIntraCodingLumaBlk\n");
}

func (this *TEncSearch) xIntraCodingChromaBlk(pcCU *TLibCommon.TComDataCU,
    uiTrDepth uint,
    uiAbsPartIdx uint,
    pcOrgYuv *TLibCommon.TComYuv,
    pcPredYuv *TLibCommon.TComYuv,
    pcResiYuv *TLibCommon.TComYuv,
    ruiDist *uint,
    uiChromaId uint,
    default0Save1Load2 int) {
    //fmt.Printf("Enter xIntraCodingChromaBlk with (%d,%d,%d)\n", uiTrDepth, uiAbsPartIdx, uiChromaId);
    
    uiOrgTrDepth := uiTrDepth
    uiFullDepth := uint(pcCU.GetDepth1(0)) + uiTrDepth
    uiLog2TrSize := uint(TLibCommon.G_aucConvertToBit[pcCU.GetSlice().GetSPS().GetMaxCUWidth()>>uiFullDepth]) + 2
    if uiLog2TrSize == 2 {
        //assert( uiTrDepth > 0 );
        uiTrDepth--
        uiQPDiv := pcCU.GetPic().GetNumPartInCU() >> ((uint(pcCU.GetDepth1(0)) + uiTrDepth) << 1)
        bFirstQ := ((uiAbsPartIdx % uiQPDiv) == 0)
        //fmt.Printf("bFirstQ%d=%d %d\n", bFirstQ, uiAbsPartIdx, uiQPDiv);
        if !bFirstQ {
        	//fmt.Printf("Exit xIntraCodingChromaBlk\n");
            return
        }
    }

    var eText TLibCommon.TextType
    if uiChromaId > 0 {
        eText = TLibCommon.TEXT_CHROMA_V
    } else {
        eText = TLibCommon.TEXT_CHROMA_U
    }
    uiChromaPredMode := uint(pcCU.GetChromaIntraDir1(uiAbsPartIdx))
    uiWidth := uint(uint(pcCU.GetWidth1(0))) >> (uiTrDepth + 1)
    uiHeight := uint(pcCU.GetHeight1(0)) >> (uiTrDepth + 1)
    uiStride := uint(pcOrgYuv.GetCStride())
    var piOrg, piPred, piResi, piReco []TLibCommon.Pel
    if uiChromaId > 0 {
        piOrg = pcOrgYuv.GetCrAddr1(uiAbsPartIdx)
        piPred = pcPredYuv.GetCrAddr1(uiAbsPartIdx)
        piResi = pcResiYuv.GetCrAddr1(uiAbsPartIdx)
        piReco = pcPredYuv.GetCrAddr1(uiAbsPartIdx)
    } else {
        piOrg = pcOrgYuv.GetCbAddr1(uiAbsPartIdx)
        piPred = pcPredYuv.GetCbAddr1(uiAbsPartIdx)
        piResi = pcResiYuv.GetCbAddr1(uiAbsPartIdx)
        piReco = pcPredYuv.GetCbAddr1(uiAbsPartIdx)
    }

    uiQTLayer := pcCU.GetSlice().GetSPS().GetQuadtreeTULog2MaxSize() - uiLog2TrSize
    uiNumCoeffPerInc := (pcCU.GetSlice().GetSPS().GetMaxCUWidth() * pcCU.GetSlice().GetSPS().GetMaxCUHeight() >> (pcCU.GetSlice().GetSPS().GetMaxCUDepth() << 1)) >> 2
    uiZOrder := pcCU.GetZorderIdxInCU() + uiAbsPartIdx
    uiRecQtStride := this.m_pcQTTempTComYuv[uiQTLayer].GetCStride()

    var pcCoeff []TLibCommon.TCoeff
    var pcArlCoeff []TLibCommon.TCoeff
    var piRecQt, piRecIPred []TLibCommon.Pel
    if uiChromaId > 0 {
        pcCoeff = this.m_ppcQTTempCoeffCr[uiQTLayer][uiNumCoeffPerInc*uiAbsPartIdx:]
        pcArlCoeff = this.m_ppcQTTempArlCoeffCr[uiQTLayer][uiNumCoeffPerInc*uiAbsPartIdx:]
        piRecQt = this.m_pcQTTempTComYuv[uiQTLayer].GetCrAddr1(uiAbsPartIdx)
        piRecIPred = pcCU.GetPic().GetPicYuvRec().GetCrAddr2(int(pcCU.GetAddr()), int(uiZOrder))
    } else {
        pcCoeff = this.m_ppcQTTempCoeffCb[uiQTLayer][uiNumCoeffPerInc*uiAbsPartIdx:]
        pcArlCoeff = this.m_ppcQTTempArlCoeffCb[uiQTLayer][uiNumCoeffPerInc*uiAbsPartIdx:]
        piRecQt = this.m_pcQTTempTComYuv[uiQTLayer].GetCbAddr1(uiAbsPartIdx)
        piRecIPred = pcCU.GetPic().GetPicYuvRec().GetCbAddr2(int(pcCU.GetAddr()), int(uiZOrder))
    }
    //#if ADAPTIVE_QP_SELECTION
    //#endif

    uiRecIPredStride := uint(pcCU.GetPic().GetPicYuvRec().GetCStride())
    useTransformSkipChroma := pcCU.GetTransformSkip2(uiAbsPartIdx, eText)
    //===== update chroma mode =====
    if uiChromaPredMode == TLibCommon.DM_CHROMA_IDX {
        uiChromaPredMode = uint(pcCU.GetLumaIntraDir1(0))
    }

    //===== init availability pattern =====
    bAboveAvail := false
    bLeftAvail := false
    if default0Save1Load2 != 2 {
        pcCU.GetPattern().InitPattern3(pcCU, uiTrDepth, uiAbsPartIdx)
        pcCU.GetPattern().InitAdiPatternChroma(pcCU, uiAbsPartIdx, uiTrDepth, this.GetYuvExt(), this.GetYuvExtStride(), this.GetYuvExtHeight(), &bAboveAvail, &bLeftAvail, uiChromaId)
        var pPatChroma []TLibCommon.Pel
        if uiChromaId > 0 {
            pPatChroma = pcCU.GetPattern().GetAdiCrBuf(int(uiWidth), int(uiHeight), this.GetYuvExt())
        } else {
            pPatChroma = pcCU.GetPattern().GetAdiCbBuf(int(uiWidth), int(uiHeight), this.GetYuvExt())
        }
        //===== get prediction signal =====
        {
            this.PredIntraChromaAng(pPatChroma, uiChromaPredMode, piPred, uiStride, int(uiWidth), int(uiHeight), bAboveAvail, bLeftAvail)
        }
        // save prediction
        if default0Save1Load2 == 1 {
            pPred := piPred
            pPredBuf := this.m_pSharedPredTransformSkip[1+uiChromaId]
            k := 0
            for uiY := uint(0); uiY < uiHeight; uiY++ {
                for uiX := uint(0); uiX < uiWidth; uiX++ {
                    pPredBuf[k] = pPred[uiY*uiStride+uiX]
                    k++
                }
                //pPred = pPred[uiStride:]
            }
        }
    } else {
        // load prediction
        pPred := piPred
        pPredBuf := this.m_pSharedPredTransformSkip[1+uiChromaId]
        k := 0
        for uiY := uint(0); uiY < uiHeight; uiY++ {
            for uiX := uint(0); uiX < uiWidth; uiX++ {
                pPred[uiY*uiStride+uiX] = pPredBuf[k]
                k++
            }
            //pPred = pPred[uiStride:]
        }
    }
    //===== get residual signal =====
    {
        // get residual
        pOrg := piOrg
        pPred := piPred
        pResi := piResi
        for uiY := uint(0); uiY < uiHeight; uiY++ {
            for uiX := uint(0); uiX < uiWidth; uiX++ {
                pResi[uiY*uiStride+uiX] = pOrg[uiY*uiStride+uiX] - pPred[uiY*uiStride+uiX]
            }
            //pOrg = pOrg[uiStride:]
            //pResi = pResi[uiStride:]
            //pPred = pPred[uiStride:]
        }
    }

    //===== transform and quantization =====
    {
        //--- init rate estimation arrays for RDOQ ---
        var rdoqflag bool
        if useTransformSkipChroma {
            rdoqflag = this.m_pcEncCfg.GetUseRDOQTS()
        } else {
            rdoqflag = this.m_pcEncCfg.GetUseRDOQ()
        }
        if rdoqflag {
            this.m_pcEntropyCoder.estimateBit(this.m_pcTrQuant.GetEstBitsSbac(), int(uiWidth), int(uiWidth), eText)
        }
        //--- transform and quantization ---
        uiAbsSum := uint(0)

        var curChromaQpOffset int
        if eText == TLibCommon.TEXT_CHROMA_U {
            curChromaQpOffset = pcCU.GetSlice().GetPPS().GetChromaCbQpOffset() + pcCU.GetSlice().GetSliceQpDeltaCb()
        } else {
            curChromaQpOffset = pcCU.GetSlice().GetPPS().GetChromaCrQpOffset() + pcCU.GetSlice().GetSliceQpDeltaCr()
        }
        this.m_pcTrQuant.SetQPforQuant(int(pcCU.GetQP1(0)), TLibCommon.TEXT_CHROMA, pcCU.GetSlice().GetSPS().GetQpBDOffsetC(), curChromaQpOffset)

        //#if RDOQ_CHROMA_LAMBDA
        this.m_pcTrQuant.SelectLambda(TLibCommon.TEXT_CHROMA)
        //#endif
        this.m_pcTrQuant.TransformNxN(pcCU, piResi, uiStride, pcCoeff,
            //#if ADAPTIVE_QP_SELECTION
            pcArlCoeff,
            //#endif
            uiWidth, uiHeight, &uiAbsSum, eText, uiAbsPartIdx, useTransformSkipChroma)
        //--- set coded block flag ---
        pcCU.SetCbfSubParts4(byte(TLibCommon.B2U(uiAbsSum != 0))<<uiOrgTrDepth, eText, uiAbsPartIdx, uint(pcCU.GetDepth1(0))+uiTrDepth)
        //--- inverse transform ---
        if uiAbsSum != 0 {
            scalingListType := 0 + int(TLibCommon.G_eTTable[eText])
            //assert(scalingListType < 6);
            /*#ifdef ENC_DEC_TRACE*/
		    this.m_pcEntropyCoder.m_pcEntropyCoderIf.XTraceCoefHeader(TLibCommon.TRACE_COEF)
		
		    for k := uint(0); k < uiHeight; k++ {
		        this.m_pcEntropyCoder.m_pcEntropyCoderIf.XReadCeofTr(pcCoeff[k*uiWidth:], uiWidth, TLibCommon.TRACE_COEF)
		    }
		    /*#endif*/
	    
            this.m_pcTrQuant.InvtransformNxN(pcCU.GetCUTransquantBypass1(uiAbsPartIdx), TLibCommon.TEXT_CHROMA, TLibCommon.REG_DCT, piResi, uiStride, pcCoeff, uiWidth, uiHeight, scalingListType, useTransformSkipChroma)
        
        	/*#ifdef ENC_DEC_TRACE*/
		    this.m_pcEntropyCoder.m_pcEntropyCoderIf.XTraceResiHeader(TLibCommon.TRACE_RESI)
		
		    for k := uint(0); k < uiHeight; k++ {
		        this.m_pcEntropyCoder.m_pcEntropyCoderIf.XReadResiTr(piResi[k*uiStride:], uiWidth, TLibCommon.TRACE_RESI)
		    }
		    /*#endif*/
        } else {
            pResi := piResi
            for i := uint(0); i < uiWidth*uiHeight; i++ {
                pcCoeff[i] = 0 //, sizeof( TLibCommon.TCoeff ) * uiWidth * uiHeight );
            }
            for uiY := uint(0); uiY < uiHeight; uiY++ {
                for uiX := uint(0); uiX < uiWidth; uiX++ {
                    pResi[uiY*uiStride+uiX] = 0 //, sizeof( TLibCommon.Pel ) * uiWidth );
                }
                //pResi = pResi[uiStride:]
            }
        }
    }

    //===== reconstruction =====
    {
        pPred := piPred
        pResi := piResi
        pReco := piReco
        pRecQt := piRecQt
        pRecIPred := piRecIPred
        for uiY := uint(0); uiY < uiHeight; uiY++ {
            for uiX := uint(0); uiX < uiWidth; uiX++ {
                pReco[uiY*uiStride+uiX] = TLibCommon.ClipC(pPred[uiY*uiStride+uiX] + pResi[uiY*uiStride+uiX])
                pRecQt[uiY*uiRecQtStride+uiX] = pReco[uiY*uiStride+uiX]
                pRecIPred[uiY*uiRecIPredStride+uiX] = pReco[uiY*uiStride+uiX]
            }
            //pPred = pPred[uiStride:]
            //pResi = pResi[uiStride:]
            //pReco = pReco[uiStride:]
            //pRecQt = pRecQt[uiRecQtStride:]
            //pRecIPred = pRecIPred[uiRecIPredStride:]
        }
    }

    //===== update distortion =====
    //#if WEIGHTED_CHROMA_DISTORTION
    *ruiDist += this.m_pcRdCost.getDistPart(TLibCommon.G_bitDepthC, piReco, int(uiStride), piOrg, int(uiStride), uiWidth, uiHeight, eText, TLibCommon.DF_SSE)
    //#else
    //  ruiDist += m_pcRdCost->getDistPart(g_bitDepthC, piReco, uiStride, piOrg, uiStride, uiWidth, uiHeight );
    //#endif
    
    //fmt.Printf("Exit xIntraCodingChromaBlk\n");
}

func (this *TEncSearch) xRecurIntraCodingQT(pcCU *TLibCommon.TComDataCU,
    uiTrDepth uint,
    uiAbsPartIdx uint,
    bLumaOnly bool,
    pcOrgYuv *TLibCommon.TComYuv,
    pcPredYuv *TLibCommon.TComYuv,
    pcResiYuv *TLibCommon.TComYuv,
    ruiDistY *uint,
    ruiDistC *uint,
    //#if HHI_RQT_INTRA_SPEEDUP
    bCheckFirst bool,
    //#endif
    dRDCost *float64) {
    
    //fmt.Printf("Enter xRecurIntraCodingQT with uiTrDepth=%d uiAbsPartIdx=%d\n", uiTrDepth, uiAbsPartIdx);
    
    uiFullDepth := uint(pcCU.GetDepth1(0)) + uiTrDepth
    uiLog2TrSize := uint(TLibCommon.G_aucConvertToBit[pcCU.GetSlice().GetSPS().GetMaxCUWidth()>>uiFullDepth]) + 2
    bCheckFull := (uiLog2TrSize <= pcCU.GetSlice().GetSPS().GetQuadtreeTULog2MaxSize())
    bCheckSplit := (uiLog2TrSize > pcCU.GetQuadtreeTULog2MinSizeInCU(uiAbsPartIdx))
	
    //#if HHI_RQT_INTRA_SPEEDUP
    //#if L0232_RD_PENALTY
    maxTuSize := int(pcCU.GetSlice().GetSPS().GetQuadtreeTULog2MaxSize())
    isIntraSlice := (pcCU.GetSlice().GetSliceType() == TLibCommon.I_SLICE)
    // don't check split if TU size is less or equal to max TU size
    noSplitIntraMaxTuSize := bCheckFull
    if this.m_pcEncCfg.GetRDpenalty() != 0 && !isIntraSlice {
        // in addition don't check split if TU size is less or equal to 16x16 TU size for non-intra slice
        noSplitIntraMaxTuSize = (int(uiLog2TrSize) <= TLibCommon.MIN(maxTuSize, 4).(int))

        // if maximum RD-penalty don't check TU size 32x32
        if this.m_pcEncCfg.GetRDpenalty() == 2 {
            bCheckFull = (int(uiLog2TrSize) <= TLibCommon.MIN(maxTuSize, 4).(int))
        }
    }
    if bCheckFirst && noSplitIntraMaxTuSize {
        //#else
        //    if bCheckFirst && bCheckFull {
        //#endif
        bCheckSplit = false
    }
    //#endif
    dSingleCost := float64(TLibCommon.MAX_DOUBLE)
    uiSingleDistY := uint(0)
    uiSingleDistC := uint(0)
    uiSingleCbfY := uint(0)
    uiSingleCbfU := uint(0)
    uiSingleCbfV := uint(0)
    checkTransformSkip := pcCU.GetSlice().GetPPS().GetUseTransformSkip()
    widthTransformSkip := uint(uint(pcCU.GetWidth1(0))) >> uiTrDepth
    heightTransformSkip := uint(pcCU.GetHeight1(0)) >> uiTrDepth
    bestModeId := 0
    var bestModeIdUV = [2]int{0, 0}
    checkTransformSkip = checkTransformSkip && (widthTransformSkip == 4 && heightTransformSkip == 4)
    checkTransformSkip = checkTransformSkip && (!pcCU.GetCUTransquantBypass1(0))
    checkTransformSkip = checkTransformSkip && (!((pcCU.GetQP1(0) == 0) && (pcCU.GetSlice().GetSPS().GetUseLossless())))
    if this.m_pcEncCfg.GetUseTransformSkipFast() {
        checkTransformSkip = checkTransformSkip && (pcCU.GetPartitionSize1(uiAbsPartIdx) == TLibCommon.SIZE_NxN)
    }
    if bCheckFull {
        if checkTransformSkip == true {
            //----- store original entropy coding status -----
            if this.m_bUseSBACRD {
                this.m_pcRDGoOnSbacCoder.store(this.m_pppcRDSbacCoder[uiFullDepth][TLibCommon.CI_QT_TRAFO_ROOT])
            }
            singleDistYTmp := uint(0)
            singleDistCTmp := uint(0)
            singleCbfYTmp := uint(0)
            singleCbfUTmp := uint(0)
            singleCbfVTmp := uint(0)
            singleCostTmp := float64(0)
            default0Save1Load2 := 0
            firstCheckId := 0

            uiQPDiv := pcCU.GetPic().GetNumPartInCU() >> ((uint(pcCU.GetDepth1(0)) + (uiTrDepth - 1)) << 1)
            bFirstQ := ((uiAbsPartIdx % uiQPDiv) == 0)

            for modeId := firstCheckId; modeId < 2; modeId++ {
                singleDistYTmp = 0
                singleDistCTmp = 0
                pcCU.SetTransformSkipSubParts4(modeId != 0, TLibCommon.TEXT_LUMA, uiAbsPartIdx, uiFullDepth)
                if modeId == firstCheckId {
                    default0Save1Load2 = 1
                } else {
                    default0Save1Load2 = 2
                }
                //----- code luma block with given intra prediction mode and store Cbf-----
                this.xIntraCodingLumaBlk(pcCU, uiTrDepth, uiAbsPartIdx, pcOrgYuv, pcPredYuv, pcResiYuv, &singleDistYTmp, default0Save1Load2)
                singleCbfYTmp = uint(pcCU.GetCbf3(uiAbsPartIdx, TLibCommon.TEXT_LUMA, uiTrDepth))
                //----- code chroma blocks with given intra prediction mode and store Cbf-----
                if !bLumaOnly {
                    if bFirstQ {
                        pcCU.SetTransformSkipSubParts4(modeId != 0, TLibCommon.TEXT_CHROMA_U, uiAbsPartIdx, uiFullDepth)
                        pcCU.SetTransformSkipSubParts4(modeId != 0, TLibCommon.TEXT_CHROMA_V, uiAbsPartIdx, uiFullDepth)
                    }
                    this.xIntraCodingChromaBlk(pcCU, uiTrDepth, uiAbsPartIdx, pcOrgYuv, pcPredYuv, pcResiYuv, &singleDistCTmp, 0, default0Save1Load2)
                    this.xIntraCodingChromaBlk(pcCU, uiTrDepth, uiAbsPartIdx, pcOrgYuv, pcPredYuv, pcResiYuv, &singleDistCTmp, 1, default0Save1Load2)
                    singleCbfUTmp = uint(pcCU.GetCbf3(uiAbsPartIdx, TLibCommon.TEXT_CHROMA_U, uiTrDepth))
                    singleCbfVTmp = uint(pcCU.GetCbf3(uiAbsPartIdx, TLibCommon.TEXT_CHROMA_V, uiTrDepth))
                }
                //----- determine rate and r-d cost -----
                if modeId == 1 && singleCbfYTmp == 0 {
                    //In order not to code TS flag when cbf is zero, the case for TS with cbf being zero is forbidden.
                    singleCostTmp = float64(TLibCommon.MAX_DOUBLE)
                } else {
                    uiSingleBits := this.xGetIntraBitsQT(pcCU, uiTrDepth, uiAbsPartIdx, true, !bLumaOnly, false)
                    //#if L0232_RD_PENALTY
                    if this.m_pcEncCfg.GetRDpenalty() != 0 && (uiLog2TrSize == 5) && !isIntraSlice {
                        uiSingleBits = uiSingleBits * 4
                    }
                    //#endif
                    singleCostTmp = this.m_pcRdCost.calcRdCost(uiSingleBits, singleDistYTmp+singleDistCTmp, false, TLibCommon.DF_DEFAULT)
                	//fmt.Printf("singleCostTmp %f = uiSingleBits %d, singleDistYTmp %d + singleDistCTmp %d\n",singleCostTmp, uiSingleBits, singleDistYTmp, singleDistCTmp);
                }

                if singleCostTmp < dSingleCost {
                    dSingleCost = singleCostTmp
                    uiSingleDistY = singleDistYTmp
                    uiSingleDistC = singleDistCTmp
                    uiSingleCbfY = singleCbfYTmp
                    uiSingleCbfU = singleCbfUTmp
                    uiSingleCbfV = singleCbfVTmp
                    bestModeId = modeId
                    if bestModeId == firstCheckId {
                        this.xStoreIntraResultQT(pcCU, uiTrDepth, uiAbsPartIdx, bLumaOnly)
                        if this.m_bUseSBACRD {
                            this.m_pcRDGoOnSbacCoder.store(this.m_pppcRDSbacCoder[uiFullDepth][TLibCommon.CI_TEMP_BEST])
                        }
                    }
                }
                if modeId == firstCheckId {
                    this.m_pcRDGoOnSbacCoder.load(this.m_pppcRDSbacCoder[uiFullDepth][TLibCommon.CI_QT_TRAFO_ROOT])
                }
            }

            pcCU.SetTransformSkipSubParts4(bestModeId != 0, TLibCommon.TEXT_LUMA, uiAbsPartIdx, uiFullDepth)

            if bestModeId == firstCheckId {
                this.xLoadIntraResultQT(pcCU, uiTrDepth, uiAbsPartIdx, bLumaOnly)
                pcCU.SetCbfSubParts4(byte(uiSingleCbfY<<uiTrDepth), TLibCommon.TEXT_LUMA, uiAbsPartIdx, uiFullDepth)
                if !bLumaOnly {
                    if bFirstQ {
                        pcCU.SetCbfSubParts4(byte(uiSingleCbfU<<uiTrDepth), TLibCommon.TEXT_CHROMA_U, uiAbsPartIdx, uint(pcCU.GetDepth1(0))+uiTrDepth-1)
                        pcCU.SetCbfSubParts4(byte(uiSingleCbfV<<uiTrDepth), TLibCommon.TEXT_CHROMA_V, uiAbsPartIdx, uint(pcCU.GetDepth1(0))+uiTrDepth-1)
                    }
                }
                if this.m_bUseSBACRD {
                    this.m_pcRDGoOnSbacCoder.load(this.m_pppcRDSbacCoder[uiFullDepth][TLibCommon.CI_TEMP_BEST])
                }
            }

            if !bLumaOnly {
                bestModeIdUV[0] = bestModeId
                bestModeIdUV[1] = bestModeId
                if bFirstQ && bestModeId == 1 {
                    //In order not to code TS flag when cbf is zero, the case for TS with cbf being zero is forbidden.
                    if uiSingleCbfU == 0 {
                        pcCU.SetTransformSkipSubParts4(false, TLibCommon.TEXT_CHROMA_U, uiAbsPartIdx, uiFullDepth)
                        bestModeIdUV[0] = 0
                    }
                    if uiSingleCbfV == 0 {
                        pcCU.SetTransformSkipSubParts4(false, TLibCommon.TEXT_CHROMA_V, uiAbsPartIdx, uiFullDepth)
                        bestModeIdUV[1] = 0
                    }
                }
            }
        } else {
            pcCU.SetTransformSkipSubParts4(false, TLibCommon.TEXT_LUMA, uiAbsPartIdx, uiFullDepth)
            //----- store original entropy coding status -----
            if this.m_bUseSBACRD && bCheckSplit {
                this.m_pcRDGoOnSbacCoder.store(this.m_pppcRDSbacCoder[uiFullDepth][TLibCommon.CI_QT_TRAFO_ROOT])
            }
            //----- code luma block with given intra prediction mode and store Cbf-----
            dSingleCost = 0.0
            this.xIntraCodingLumaBlk(pcCU, uiTrDepth, uiAbsPartIdx, pcOrgYuv, pcPredYuv, pcResiYuv, &uiSingleDistY, 0)
            //fmt.Printf("uiSingleDistY=%d\n", uiSingleDistY);
            if bCheckSplit {
                uiSingleCbfY = uint(pcCU.GetCbf3(uiAbsPartIdx, TLibCommon.TEXT_LUMA, uiTrDepth))
            }
            //----- code chroma blocks with given intra prediction mode and store Cbf-----
            if !bLumaOnly {
                pcCU.SetTransformSkipSubParts4(false, TLibCommon.TEXT_CHROMA_U, uiAbsPartIdx, uiFullDepth)
                pcCU.SetTransformSkipSubParts4(false, TLibCommon.TEXT_CHROMA_V, uiAbsPartIdx, uiFullDepth)
                this.xIntraCodingChromaBlk(pcCU, uiTrDepth, uiAbsPartIdx, pcOrgYuv, pcPredYuv, pcResiYuv, &uiSingleDistC, 0, 0)
                this.xIntraCodingChromaBlk(pcCU, uiTrDepth, uiAbsPartIdx, pcOrgYuv, pcPredYuv, pcResiYuv, &uiSingleDistC, 1, 0)
                if bCheckSplit {
                    uiSingleCbfU = uint(pcCU.GetCbf3(uiAbsPartIdx, TLibCommon.TEXT_CHROMA_U, uiTrDepth))
                    uiSingleCbfV = uint(pcCU.GetCbf3(uiAbsPartIdx, TLibCommon.TEXT_CHROMA_V, uiTrDepth))
                }
            }
            //----- determine rate and r-d cost -----
            uiSingleBits := this.xGetIntraBitsQT(pcCU, uiTrDepth, uiAbsPartIdx, true, !bLumaOnly, false)
//#if L0232_RD_PENALTY
		    if this.m_pcEncCfg.GetRDpenalty()!=0 && (uiLog2TrSize==5) && !isIntraSlice {
		        uiSingleBits=uiSingleBits*4; 
		    }
//#endif
            dSingleCost = this.m_pcRdCost.calcRdCost(uiSingleBits, uiSingleDistY+uiSingleDistC, false, TLibCommon.DF_DEFAULT)
        	//fmt.Printf("dSingleCost %f = uiSingleBits %d, uiSingleDistY %d + uiSingleDistC %d\n",dSingleCost, uiSingleBits, uiSingleDistY, uiSingleDistC);
        }
    }

    if bCheckSplit {
        //----- store full entropy coding status, load original entropy coding status -----
        if this.m_bUseSBACRD {
            if bCheckFull {
                this.m_pcRDGoOnSbacCoder.store(this.m_pppcRDSbacCoder[uiFullDepth][TLibCommon.CI_QT_TRAFO_TEST])
                this.m_pcRDGoOnSbacCoder.load(this.m_pppcRDSbacCoder[uiFullDepth][TLibCommon.CI_QT_TRAFO_ROOT])
            } else {
                this.m_pcRDGoOnSbacCoder.store(this.m_pppcRDSbacCoder[uiFullDepth][TLibCommon.CI_QT_TRAFO_ROOT])
            }
        }
        //----- code splitted block -----
        dSplitCost := float64(0.0)
        uiSplitDistY := uint(0)
        uiSplitDistC := uint(0)
        uiQPartsDiv := pcCU.GetPic().GetNumPartInCU() >> ((uiFullDepth + 1) << 1)
        uiAbsPartIdxSub := uiAbsPartIdx

        uiSplitCbfY := uint(0)
        uiSplitCbfU := uint(0)
        uiSplitCbfV := uint(0)

        for uiPart := uint(0); uiPart < 4; uiPart++ {
            //#if HHI_RQT_INTRA_SPEEDUP
            this.xRecurIntraCodingQT(pcCU, uiTrDepth+1, uiAbsPartIdxSub, bLumaOnly, pcOrgYuv, pcPredYuv, pcResiYuv, &uiSplitDistY, &uiSplitDistC, bCheckFirst, &dSplitCost)
            //#else
            //      xRecurIntraCodingQT( pcCU, uiTrDepth + 1, uiAbsPartIdxSub, bLumaOnly, pcOrgYuv, pcPredYuv, pcResiYuv, uiSplitDistY, uiSplitDistC, dSplitCost );
            //#endif

            uiSplitCbfY |= uint(pcCU.GetCbf3(uiAbsPartIdxSub, TLibCommon.TEXT_LUMA, uiTrDepth+1))
            if !bLumaOnly {
                uiSplitCbfU |= uint(pcCU.GetCbf3(uiAbsPartIdxSub, TLibCommon.TEXT_CHROMA_U, uiTrDepth+1))
                uiSplitCbfV |= uint(pcCU.GetCbf3(uiAbsPartIdxSub, TLibCommon.TEXT_CHROMA_V, uiTrDepth+1))
            }

            uiAbsPartIdxSub += uiQPartsDiv
        }

        for uiOffs := uint(0); uiOffs < 4*uiQPartsDiv; uiOffs++ {
            pcCU.GetCbf1(TLibCommon.TEXT_LUMA)[uiAbsPartIdx+uiOffs] |= byte(uiSplitCbfY << uiTrDepth)
        }
        if !bLumaOnly {
            for uiOffs := uint(0); uiOffs < 4*uiQPartsDiv; uiOffs++ {
                pcCU.GetCbf1(TLibCommon.TEXT_CHROMA_U)[uiAbsPartIdx+uiOffs] |= byte(uiSplitCbfU << uiTrDepth)
                pcCU.GetCbf1(TLibCommon.TEXT_CHROMA_V)[uiAbsPartIdx+uiOffs] |= byte(uiSplitCbfV << uiTrDepth)
            }
        }
        //----- restore context states -----
        if this.m_bUseSBACRD {
            this.m_pcRDGoOnSbacCoder.load(this.m_pppcRDSbacCoder[uiFullDepth][TLibCommon.CI_QT_TRAFO_ROOT])
        }
        //----- determine rate and r-d cost -----
        uiSplitBits := this.xGetIntraBitsQT(pcCU, uiTrDepth, uiAbsPartIdx, true, !bLumaOnly, false)
        dSplitCost = this.m_pcRdCost.calcRdCost(uiSplitBits, uiSplitDistY+uiSplitDistC, false, TLibCommon.DF_DEFAULT)
		//fmt.Printf("dSplitCost %f dSingleCost %f = uiSplitBits %d, uiSplitDistY %d + uiSplitDistC %d\n",dSplitCost, dSingleCost, uiSplitBits, uiSplitDistY, uiSplitDistC);
        //===== compare and set best =====
        if dSplitCost < dSingleCost {
            //--- update cost ---
            *ruiDistY += uiSplitDistY
            *ruiDistC += uiSplitDistC
            *dRDCost += dSplitCost
            //fmt.Printf("Exit xRecurIntraCodingQT with uiTrDepth=%d uiAbsPartIdx=%d in dSplitCost %f< dSingleCost %f\n", uiTrDepth, uiAbsPartIdx, dSplitCost, dSingleCost);
            return
        }
        //----- set entropy coding status -----
        if this.m_bUseSBACRD {
            this.m_pcRDGoOnSbacCoder.load(this.m_pppcRDSbacCoder[uiFullDepth][TLibCommon.CI_QT_TRAFO_TEST])
        }

        //--- set transform index and Cbf values ---
        pcCU.SetTrIdxSubParts(uiTrDepth, uiAbsPartIdx, uiFullDepth)
        pcCU.SetCbfSubParts4(byte(uiSingleCbfY<<uiTrDepth), TLibCommon.TEXT_LUMA, uiAbsPartIdx, uiFullDepth)
        pcCU.SetTransformSkipSubParts4(bestModeId != 0, TLibCommon.TEXT_LUMA, uiAbsPartIdx, uiFullDepth)
        if !bLumaOnly {
            pcCU.SetCbfSubParts4(byte(uiSingleCbfU<<uiTrDepth), TLibCommon.TEXT_CHROMA_U, uiAbsPartIdx, uiFullDepth)
            pcCU.SetCbfSubParts4(byte(uiSingleCbfV<<uiTrDepth), TLibCommon.TEXT_CHROMA_V, uiAbsPartIdx, uiFullDepth)
            pcCU.SetTransformSkipSubParts4(bestModeIdUV[0] != 0, TLibCommon.TEXT_CHROMA_U, uiAbsPartIdx, uiFullDepth)
            pcCU.SetTransformSkipSubParts4(bestModeIdUV[1] != 0, TLibCommon.TEXT_CHROMA_V, uiAbsPartIdx, uiFullDepth)
        }

        //--- set reconstruction for next intra prediction blocks ---
        uiWidth := uint(uint(pcCU.GetWidth1(0))) >> uiTrDepth
        uiHeight := uint(pcCU.GetHeight1(0)) >> uiTrDepth
        uiQTLayer := pcCU.GetSlice().GetSPS().GetQuadtreeTULog2MaxSize() - uiLog2TrSize
        uiZOrder := pcCU.GetZorderIdxInCU() + uiAbsPartIdx
        piSrc := this.m_pcQTTempTComYuv[uiQTLayer].GetLumaAddr1(uiAbsPartIdx)
        uiSrcStride := uint(this.m_pcQTTempTComYuv[uiQTLayer].GetStride())
        piDes := pcCU.GetPic().GetPicYuvRec().GetLumaAddr2(int(pcCU.GetAddr()), int(uiZOrder))
        uiDesStride := uint(pcCU.GetPic().GetPicYuvRec().GetStride())
        for uiY := uint(0); uiY < uiHeight; uiY++ {
        	uiYSrc := uiY*uiSrcStride;
            uiYDes := uiY*uiDesStride;
            for uiX := uint(0); uiX < uiWidth; uiX++ {
                piDes[uiYDes+uiX] = piSrc[uiYSrc+uiX]
            }
            //piSrc = piSrc[uiSrcStride:]
            //piDes = piDes[uiDesStride:]
        }
        if !bLumaOnly {
            uiWidth >>= 1
            uiHeight >>= 1
            piSrc = this.m_pcQTTempTComYuv[uiQTLayer].GetCbAddr1(uiAbsPartIdx)
            uiSrcStride = uint(this.m_pcQTTempTComYuv[uiQTLayer].GetCStride())
            piDes = pcCU.GetPic().GetPicYuvRec().GetCbAddr2(int(pcCU.GetAddr()), int(uiZOrder))
            uiDesStride = uint(pcCU.GetPic().GetPicYuvRec().GetCStride())
            for uiY := uint(0); uiY < uiHeight; uiY++ {
            	uiYSrc := uiY*uiSrcStride;
            	uiYDes := uiY*uiDesStride;
                for uiX := uint(0); uiX < uiWidth; uiX++ {
                    piDes[uiYDes+uiX] = piSrc[uiYSrc+uiX]
                }
                //piSrc = piSrc[uiSrcStride:]
                //piDes = piDes[uiDesStride:]
            }
            piSrc = this.m_pcQTTempTComYuv[uiQTLayer].GetCrAddr1(uiAbsPartIdx)
            piDes = pcCU.GetPic().GetPicYuvRec().GetCrAddr2(int(pcCU.GetAddr()), int(uiZOrder))
            for uiY := uint(0); uiY < uiHeight; uiY++ {
            	uiYSrc := uiY*uiSrcStride;
            	uiYDes := uiY*uiDesStride;
                for uiX := uint(0); uiX < uiWidth; uiX++ {
                    piDes[uiYDes+uiX] = piSrc[uiYSrc+uiX]
                }
                //piSrc = piSrc[uiSrcStride:]
                //piDes = piDes[uiDesStride:]
            }
        }
    }
    *ruiDistY += uiSingleDistY
    *ruiDistC += uiSingleDistC
    *dRDCost += dSingleCost
    
    //fmt.Printf("Exit xRecurIntraCodingQT with uiTrDepth=%d uiAbsPartIdx=%d ruiDistY=%d\n", uiTrDepth, uiAbsPartIdx, *ruiDistY);
}

func (this *TEncSearch) xSetIntraResultQT(pcCU *TLibCommon.TComDataCU,
    uiTrDepth uint,
    uiAbsPartIdx uint,
    bLumaOnly bool,
    pcRecoYuv *TLibCommon.TComYuv) {
    uiFullDepth := uint(pcCU.GetDepth1(0)) + uiTrDepth
    uiTrMode := uint(pcCU.GetTransformIdx1(uiAbsPartIdx))
    if uiTrMode == uiTrDepth {
        uiLog2TrSize := uint(TLibCommon.G_aucConvertToBit[pcCU.GetSlice().GetSPS().GetMaxCUWidth()>>uiFullDepth]) + 2
        uiQTLayer := pcCU.GetSlice().GetSPS().GetQuadtreeTULog2MaxSize() - uiLog2TrSize

        bSkipChroma := false
        bChromaSame := false
        if !bLumaOnly && uiLog2TrSize == 2 {
            //assert( uiTrDepth > 0 );
            uiQPDiv := pcCU.GetPic().GetNumPartInCU() >> ((uint(pcCU.GetDepth1(0)) + uiTrDepth - 1) << 1)
            bSkipChroma = ((uiAbsPartIdx % uiQPDiv) != 0)
            bChromaSame = true
        }

        //===== copy transform coefficients =====
        uiNumCoeffY := (pcCU.GetSlice().GetSPS().GetMaxCUWidth() * pcCU.GetSlice().GetSPS().GetMaxCUHeight()) >> (uiFullDepth << 1)
        uiNumCoeffIncY := (pcCU.GetSlice().GetSPS().GetMaxCUWidth() * pcCU.GetSlice().GetSPS().GetMaxCUHeight()) >> (pcCU.GetSlice().GetSPS().GetMaxCUDepth() << 1)
        pcCoeffSrcY := this.m_ppcQTTempCoeffY[uiQTLayer][(uiNumCoeffIncY * uiAbsPartIdx):]
        pcCoeffDstY := pcCU.GetCoeffY()[(uiNumCoeffIncY * uiAbsPartIdx):]
        for i := uint(0); i < uiNumCoeffY; i++ {
            pcCoeffDstY[i] = pcCoeffSrcY[i] //, sizeof( TLibCommon.TCoeff ) * uiNumCoeffY );
        }

        //#if ADAPTIVE_QP_SELECTION
        pcArlCoeffSrcY := this.m_ppcQTTempArlCoeffY[uiQTLayer][(uiNumCoeffIncY * uiAbsPartIdx):]
        pcArlCoeffDstY := pcCU.GetArlCoeffY()[(uiNumCoeffIncY * uiAbsPartIdx):]
        for i := uint(0); i < uiNumCoeffY; i++ {
            pcArlCoeffDstY[i] = pcArlCoeffSrcY[i] //, sizeof( int ) * uiNumCoeffY );
        }
        //#endif
        if !bLumaOnly && !bSkipChroma {
            var uiNumCoeffC uint
            if bChromaSame {
                uiNumCoeffC = uiNumCoeffY
            } else {
                uiNumCoeffC = uiNumCoeffY >> 2
            }
            uiNumCoeffIncC := uiNumCoeffIncY >> 2
            pcCoeffSrcU := this.m_ppcQTTempCoeffCb[uiQTLayer][(uiNumCoeffIncC * uiAbsPartIdx):]
            pcCoeffSrcV := this.m_ppcQTTempCoeffCr[uiQTLayer][(uiNumCoeffIncC * uiAbsPartIdx):]
            pcCoeffDstU := pcCU.GetCoeffCb()[(uiNumCoeffIncC * uiAbsPartIdx):]
            pcCoeffDstV := pcCU.GetCoeffCr()[(uiNumCoeffIncC * uiAbsPartIdx):]
            for i := uint(0); i < uiNumCoeffC; i++ {
                pcCoeffDstU[i] = pcCoeffSrcU[i] //, sizeof( TLibCommon.TCoeff ) * uiNumCoeffC );
                pcCoeffDstV[i] = pcCoeffSrcV[i] //, sizeof( TLibCommon.TCoeff ) * uiNumCoeffC );
            }
            //#if ADAPTIVE_QP_SELECTION
            pcArlCoeffSrcU := this.m_ppcQTTempArlCoeffCb[uiQTLayer][(uiNumCoeffIncC * uiAbsPartIdx):]
            pcArlCoeffSrcV := this.m_ppcQTTempArlCoeffCr[uiQTLayer][(uiNumCoeffIncC * uiAbsPartIdx):]
            pcArlCoeffDstU := pcCU.GetArlCoeffCb()[(uiNumCoeffIncC * uiAbsPartIdx):]
            pcArlCoeffDstV := pcCU.GetArlCoeffCr()[(uiNumCoeffIncC * uiAbsPartIdx):]
            for i := uint(0); i < uiNumCoeffC; i++ {
                pcArlCoeffDstU[i] = pcArlCoeffSrcU[i] //, sizeof( int ) * uiNumCoeffC );
                pcArlCoeffDstV[i] = pcArlCoeffSrcV[i] //, sizeof( int ) * uiNumCoeffC );
            }
            //#endif
        }

        //===== copy reconstruction =====
        this.m_pcQTTempTComYuv[uiQTLayer].CopyPartToPartLuma(pcRecoYuv, uiAbsPartIdx, 1<<uiLog2TrSize, 1<<uiLog2TrSize)
        if !bLumaOnly && !bSkipChroma {
            var uiLog2TrSizeChroma uint
            if bChromaSame {
                uiLog2TrSizeChroma = uiLog2TrSize
            } else {
                uiLog2TrSizeChroma = uiLog2TrSize - 1
            }
            this.m_pcQTTempTComYuv[uiQTLayer].CopyPartToPartChroma(pcRecoYuv, uiAbsPartIdx, 1<<uiLog2TrSizeChroma, 1<<uiLog2TrSizeChroma)
        }
    } else {
        uiNumQPart := pcCU.GetPic().GetNumPartInCU() >> ((uiFullDepth + 1) << 1)
        for uiPart := uint(0); uiPart < 4; uiPart++ {
            this.xSetIntraResultQT(pcCU, uiTrDepth+1, uiAbsPartIdx+uiPart*uiNumQPart, bLumaOnly, pcRecoYuv)
        }
    }
}

func (this *TEncSearch) xRecurIntraChromaCodingQT(pcCU *TLibCommon.TComDataCU,
    uiTrDepth uint,
    uiAbsPartIdx uint,
    pcOrgYuv *TLibCommon.TComYuv,
    pcPredYuv *TLibCommon.TComYuv,
    pcResiYuv *TLibCommon.TComYuv,
    ruiDist *uint) {
    uiFullDepth := uint(pcCU.GetDepth1(0)) + uiTrDepth
    uiTrMode := uint(pcCU.GetTransformIdx1(uiAbsPartIdx))
    if uiTrMode == uiTrDepth {
        checkTransformSkip := pcCU.GetSlice().GetPPS().GetUseTransformSkip()
        uiLog2TrSize := uint(TLibCommon.G_aucConvertToBit[pcCU.GetSlice().GetSPS().GetMaxCUWidth()>>uiFullDepth]) + 2

        actualTrDepth := uiTrDepth
        if uiLog2TrSize == 2 {
            //assert( uiTrDepth > 0 );
            actualTrDepth--
            uiQPDiv := pcCU.GetPic().GetNumPartInCU() >> ((uint(pcCU.GetDepth1(0)) + actualTrDepth) << 1)
            bFirstQ := ((uiAbsPartIdx % uiQPDiv) == 0)
            if !bFirstQ {
                return
            }
        }

        checkTransformSkip = checkTransformSkip && (uiLog2TrSize <= 3)
        if this.m_pcEncCfg.GetUseTransformSkipFast() {
            checkTransformSkip = checkTransformSkip && (uiLog2TrSize < 3)
            if checkTransformSkip {
                nbLumaSkip := 0
                for absPartIdxSub := uiAbsPartIdx; absPartIdxSub < uiAbsPartIdx+4; absPartIdxSub++ {
                    nbLumaSkip += int(TLibCommon.B2U(pcCU.GetTransformSkip2(absPartIdxSub, TLibCommon.TEXT_LUMA)))
                }
                checkTransformSkip = checkTransformSkip && (nbLumaSkip > 0)
            }
        }

        if checkTransformSkip {
            //use RDO to decide whether Cr/Cb takes TS
            if this.m_bUseSBACRD {
                this.m_pcRDGoOnSbacCoder.store(this.m_pppcRDSbacCoder[uiFullDepth][TLibCommon.CI_QT_TRAFO_ROOT])
            }

            for chromaId := 0; chromaId < 2; chromaId++ {
                dSingleCost := float64(TLibCommon.MAX_DOUBLE)
                bestModeId := uint(0)
                singleDistC := uint(0)
                singleCbfC := uint(0)
                singleDistCTmp := uint(0)
                singleCostTmp := float64(0)
                singleCbfCTmp := uint(0)

                default0Save1Load2 := 0
                firstCheckId := 0

                for chromaModeId := firstCheckId; chromaModeId < 2; chromaModeId++ {
                    pcCU.SetTransformSkipSubParts4(chromaModeId != 0, TLibCommon.TextType(chromaId+2), uiAbsPartIdx, uint(pcCU.GetDepth1(0))+actualTrDepth)
                    if chromaModeId == firstCheckId {
                        default0Save1Load2 = 1
                    } else {
                        default0Save1Load2 = 2
                    }
                    singleDistCTmp = 0
                    this.xIntraCodingChromaBlk(pcCU, uiTrDepth, uiAbsPartIdx, pcOrgYuv, pcPredYuv, pcResiYuv, &singleDistCTmp, uint(chromaId), default0Save1Load2)
                    singleCbfCTmp = uint(pcCU.GetCbf3(uiAbsPartIdx, TLibCommon.TextType(chromaId+2), uiTrDepth))

                    if chromaModeId == 1 && singleCbfCTmp == 0 {
                        //In order not to code TS flag when cbf is zero, the case for TS with cbf being zero is forbidden.
                        singleCostTmp = float64(TLibCommon.MAX_DOUBLE)
                    } else {
                        bitsTmp := this.xGetIntraBitsQTChroma(pcCU, uiTrDepth, uiAbsPartIdx, uint(chromaId+2), false)
                        singleCostTmp = this.m_pcRdCost.calcRdCost(bitsTmp, singleDistCTmp, false, TLibCommon.DF_DEFAULT)
                    }

                    if singleCostTmp < dSingleCost {
                        dSingleCost = singleCostTmp
                        singleDistC = singleDistCTmp
                        bestModeId = uint(chromaModeId)
                        singleCbfC = singleCbfCTmp

                        if bestModeId == uint(firstCheckId) {
                            this.xStoreIntraResultChromaQT(pcCU, uiTrDepth, uiAbsPartIdx, uint(chromaId))
                            if this.m_bUseSBACRD {
                                this.m_pcRDGoOnSbacCoder.store(this.m_pppcRDSbacCoder[uiFullDepth][TLibCommon.CI_TEMP_BEST])
                            }
                        }
                    }
                    if chromaModeId == firstCheckId {
                        this.m_pcRDGoOnSbacCoder.load(this.m_pppcRDSbacCoder[uiFullDepth][TLibCommon.CI_QT_TRAFO_ROOT])
                    }
                }

                if bestModeId == uint(firstCheckId) {
                    this.xLoadIntraResultChromaQT(pcCU, uiTrDepth, uiAbsPartIdx, uint(chromaId))
                    pcCU.SetCbfSubParts4(byte(singleCbfC<<uiTrDepth), TLibCommon.TextType(chromaId+2), uiAbsPartIdx, uint(pcCU.GetDepth1(0))+actualTrDepth)
                    if this.m_bUseSBACRD {
                        this.m_pcRDGoOnSbacCoder.load(this.m_pppcRDSbacCoder[uiFullDepth][TLibCommon.CI_TEMP_BEST])
                    }
                }
                pcCU.SetTransformSkipSubParts4(bestModeId != 0, TLibCommon.TextType(chromaId+2), uiAbsPartIdx, uint(pcCU.GetDepth1(0))+actualTrDepth)
                *ruiDist += singleDistC

                if chromaId == 0 {
                    if this.m_bUseSBACRD {
                        this.m_pcRDGoOnSbacCoder.store(this.m_pppcRDSbacCoder[uiFullDepth][TLibCommon.CI_QT_TRAFO_ROOT])
                    }
                }
            }
        } else {
            pcCU.SetTransformSkipSubParts4(false, TLibCommon.TEXT_CHROMA_U, uiAbsPartIdx, uint(pcCU.GetDepth1(0))+actualTrDepth)
            pcCU.SetTransformSkipSubParts4(false, TLibCommon.TEXT_CHROMA_V, uiAbsPartIdx, uint(pcCU.GetDepth1(0))+actualTrDepth)
            this.xIntraCodingChromaBlk(pcCU, uiTrDepth, uiAbsPartIdx, pcOrgYuv, pcPredYuv, pcResiYuv, ruiDist, 0, 0)
            this.xIntraCodingChromaBlk(pcCU, uiTrDepth, uiAbsPartIdx, pcOrgYuv, pcPredYuv, pcResiYuv, ruiDist, 1, 0)
        }
    } else {
        uiSplitCbfU := uint(0)
        uiSplitCbfV := uint(0)
        uiQPartsDiv := pcCU.GetPic().GetNumPartInCU() >> ((uiFullDepth + 1) << 1)
        uiAbsPartIdxSub := uiAbsPartIdx
        for uiPart := uint(0); uiPart < 4; uiPart++ {
            this.xRecurIntraChromaCodingQT(pcCU, uiTrDepth+1, uiAbsPartIdxSub, pcOrgYuv, pcPredYuv, pcResiYuv, ruiDist)
            uiSplitCbfU |= uint(pcCU.GetCbf3(uiAbsPartIdxSub, TLibCommon.TEXT_CHROMA_U, uiTrDepth+1))
            uiSplitCbfV |= uint(pcCU.GetCbf3(uiAbsPartIdxSub, TLibCommon.TEXT_CHROMA_V, uiTrDepth+1))
            uiAbsPartIdxSub += uiQPartsDiv
        }
        for uiOffs := uint(0); uiOffs < 4*uiQPartsDiv; uiOffs++ {
            pcCU.GetCbf1(TLibCommon.TEXT_CHROMA_U)[uiAbsPartIdx+uiOffs] |= byte(uiSplitCbfU << uiTrDepth)
            pcCU.GetCbf1(TLibCommon.TEXT_CHROMA_V)[uiAbsPartIdx+uiOffs] |= byte(uiSplitCbfV << uiTrDepth)
        }
    }
}

func (this *TEncSearch) xSetIntraResultChromaQT(pcCU *TLibCommon.TComDataCU,
    uiTrDepth uint,
    uiAbsPartIdx uint,
    pcRecoYuv *TLibCommon.TComYuv) {
    uiFullDepth := uint(pcCU.GetDepth1(0)) + uiTrDepth
    uiTrMode := uint(pcCU.GetTransformIdx1(uiAbsPartIdx))
    if uiTrMode == uiTrDepth {
        uiLog2TrSize := uint(TLibCommon.G_aucConvertToBit[pcCU.GetSlice().GetSPS().GetMaxCUWidth()>>uiFullDepth]) + 2
        uiQTLayer := pcCU.GetSlice().GetSPS().GetQuadtreeTULog2MaxSize() - uiLog2TrSize

        bChromaSame := false
        if uiLog2TrSize == 2 {
            //assert( uiTrDepth > 0 );
            uiQPDiv := pcCU.GetPic().GetNumPartInCU() >> ((uint(pcCU.GetDepth1(0)) + uiTrDepth - 1) << 1)
            if (uiAbsPartIdx % uiQPDiv) != 0 {
                return
            }
            bChromaSame = true
        }

        //===== copy transform coefficients =====
        uiNumCoeffC := (pcCU.GetSlice().GetSPS().GetMaxCUWidth() * pcCU.GetSlice().GetSPS().GetMaxCUHeight()) >> (uiFullDepth << 1)
        if !bChromaSame {
            uiNumCoeffC >>= 2
        }
        uiNumCoeffIncC := (pcCU.GetSlice().GetSPS().GetMaxCUWidth() * pcCU.GetSlice().GetSPS().GetMaxCUHeight()) >> ((pcCU.GetSlice().GetSPS().GetMaxCUDepth() << 1) + 2)
        pcCoeffSrcU := this.m_ppcQTTempCoeffCb[uiQTLayer][(uiNumCoeffIncC * uiAbsPartIdx):]
        pcCoeffSrcV := this.m_ppcQTTempCoeffCr[uiQTLayer][(uiNumCoeffIncC * uiAbsPartIdx):]
        pcCoeffDstU := pcCU.GetCoeffCb()[(uiNumCoeffIncC * uiAbsPartIdx):]
        pcCoeffDstV := pcCU.GetCoeffCr()[(uiNumCoeffIncC * uiAbsPartIdx):]
        for i := uint(0); i < uiNumCoeffC; i++ {
            pcCoeffDstU[i] = pcCoeffSrcU[i] //, sizeof( TLibCommon.TCoeff ) * uiNumCoeffC );
            pcCoeffDstV[i] = pcCoeffSrcV[i] //, sizeof( TLibCommon.TCoeff ) * uiNumCoeffC );
        }
        //#if ADAPTIVE_QP_SELECTION
        pcArlCoeffSrcU := this.m_ppcQTTempArlCoeffCb[uiQTLayer][(uiNumCoeffIncC * uiAbsPartIdx):]
        pcArlCoeffSrcV := this.m_ppcQTTempArlCoeffCr[uiQTLayer][(uiNumCoeffIncC * uiAbsPartIdx):]
        pcArlCoeffDstU := pcCU.GetArlCoeffCb()[(uiNumCoeffIncC * uiAbsPartIdx):]
        pcArlCoeffDstV := pcCU.GetArlCoeffCr()[(uiNumCoeffIncC * uiAbsPartIdx):]
        for i := uint(0); i < uiNumCoeffC; i++ {
            pcArlCoeffDstU[i] = pcArlCoeffSrcU[i] //, sizeof( int ) * uiNumCoeffC );
            pcArlCoeffDstV[i] = pcArlCoeffSrcV[i] //, sizeof( int ) * uiNumCoeffC );
        }
        //#endif

        //===== copy reconstruction =====
        var uiLog2TrSizeChroma uint
        if bChromaSame {
            uiLog2TrSizeChroma = uiLog2TrSize
        } else {
            uiLog2TrSizeChroma = uiLog2TrSize - 1
        }
        this.m_pcQTTempTComYuv[uiQTLayer].CopyPartToPartChroma(pcRecoYuv, uiAbsPartIdx, 1<<uiLog2TrSizeChroma, 1<<uiLog2TrSizeChroma)
    } else {
        uiNumQPart := pcCU.GetPic().GetNumPartInCU() >> ((uiFullDepth + 1) << 1)
        for uiPart := uint(0); uiPart < 4; uiPart++ {
            this.xSetIntraResultChromaQT(pcCU, uiTrDepth+1, uiAbsPartIdx+uiPart*uiNumQPart, pcRecoYuv)
        }
    }
}

func (this *TEncSearch) xStoreIntraResultQT(pcCU *TLibCommon.TComDataCU,
    uiTrDepth uint,
    uiAbsPartIdx uint,
    bLumaOnly bool) {
    uiFullDepth := uint(pcCU.GetDepth1(0)) + uiTrDepth
    //uiTrMode := uint(pcCU.GetTransformIdx1(uiAbsPartIdx))
    //assert(  uiTrMode == uiTrDepth );
    uiLog2TrSize := uint(TLibCommon.G_aucConvertToBit[pcCU.GetSlice().GetSPS().GetMaxCUWidth()>>uiFullDepth]) + 2
    uiQTLayer := pcCU.GetSlice().GetSPS().GetQuadtreeTULog2MaxSize() - uiLog2TrSize

    bSkipChroma := false
    bChromaSame := false
    if !bLumaOnly && uiLog2TrSize == 2 {
        //assert( uiTrDepth > 0 );
        uiQPDiv := pcCU.GetPic().GetNumPartInCU() >> ((uint(pcCU.GetDepth1(0)) + uiTrDepth - 1) << 1)
        bSkipChroma = ((uiAbsPartIdx % uiQPDiv) != 0)
        bChromaSame = true
    }

    //===== copy transform coefficients =====
    uiNumCoeffY := (pcCU.GetSlice().GetSPS().GetMaxCUWidth() * pcCU.GetSlice().GetSPS().GetMaxCUHeight()) >> (uiFullDepth << 1)
    uiNumCoeffIncY := (pcCU.GetSlice().GetSPS().GetMaxCUWidth() * pcCU.GetSlice().GetSPS().GetMaxCUHeight()) >> (pcCU.GetSlice().GetSPS().GetMaxCUDepth() << 1)
    pcCoeffSrcY := this.m_ppcQTTempCoeffY[uiQTLayer][(uiNumCoeffIncY * uiAbsPartIdx):]
    pcCoeffDstY := this.m_pcQTTempTUCoeffY[:]

    for i := uint(0); i < uiNumCoeffY; i++ {
        pcCoeffDstY[i] = pcCoeffSrcY[i] //, sizeof( TLibCommon.TCoeff ) * uiNumCoeffY );
    }
    //#if ADAPTIVE_QP_SELECTION
    pcArlCoeffSrcY := this.m_ppcQTTempArlCoeffY[uiQTLayer][(uiNumCoeffIncY * uiAbsPartIdx):]
    pcArlCoeffDstY := this.m_ppcQTTempTUArlCoeffY[:]
    for i := uint(0); i < uiNumCoeffY; i++ {
        pcArlCoeffDstY[i] = pcArlCoeffSrcY[i] //, sizeof( int ) * uiNumCoeffY );
    }
    //#endif
    if !bLumaOnly && !bSkipChroma {
        var uiNumCoeffC uint
        if bChromaSame {
            uiNumCoeffC = uiNumCoeffY
        } else {
            uiNumCoeffC = uiNumCoeffY >> 2
        }
        uiNumCoeffIncC := uiNumCoeffIncY >> 2
        pcCoeffSrcU := this.m_ppcQTTempCoeffCb[uiQTLayer][(uiNumCoeffIncC * uiAbsPartIdx):]
        pcCoeffSrcV := this.m_ppcQTTempCoeffCr[uiQTLayer][(uiNumCoeffIncC * uiAbsPartIdx):]
        pcCoeffDstU := this.m_pcQTTempTUCoeffCb[:]
        pcCoeffDstV := this.m_pcQTTempTUCoeffCr[:]
        for i := uint(0); i < uiNumCoeffC; i++ {
            pcCoeffDstU[i] = pcCoeffSrcU[i] //, sizeof( TLibCommon.TCoeff ) * uiNumCoeffC );
            pcCoeffDstV[i] = pcCoeffSrcV[i] //, sizeof( TLibCommon.TCoeff ) * uiNumCoeffC );
        }
        //#if ADAPTIVE_QP_SELECTION
        pcArlCoeffSrcU := this.m_ppcQTTempArlCoeffCb[uiQTLayer][(uiNumCoeffIncC * uiAbsPartIdx):]
        pcArlCoeffSrcV := this.m_ppcQTTempArlCoeffCr[uiQTLayer][(uiNumCoeffIncC * uiAbsPartIdx):]
        pcArlCoeffDstU := this.m_ppcQTTempTUArlCoeffCb[:]
        pcArlCoeffDstV := this.m_ppcQTTempTUArlCoeffCr[:]
        for i := uint(0); i < uiNumCoeffC; i++ {
            pcArlCoeffDstU[i] = pcArlCoeffSrcU[i] //, sizeof( int ) * uiNumCoeffC );
            pcArlCoeffDstV[i] = pcArlCoeffSrcV[i] //, sizeof( int ) * uiNumCoeffC );
        }
        //#endif
    }

    //===== copy reconstruction =====
    this.m_pcQTTempTComYuv[uiQTLayer].CopyPartToPartLuma(&this.m_pcQTTempTransformSkipTComYuv, uiAbsPartIdx, 1<<uiLog2TrSize, 1<<uiLog2TrSize)

    if !bLumaOnly && !bSkipChroma {
        var uiLog2TrSizeChroma uint
        if bChromaSame {
            uiLog2TrSizeChroma = uiLog2TrSize
        } else {
            uiLog2TrSizeChroma = uiLog2TrSize - 1
        }
        this.m_pcQTTempTComYuv[uiQTLayer].CopyPartToPartChroma(&this.m_pcQTTempTransformSkipTComYuv, uiAbsPartIdx, 1<<uiLog2TrSizeChroma, 1<<uiLog2TrSizeChroma)
    }
}
func (this *TEncSearch) xLoadIntraResultQT(pcCU *TLibCommon.TComDataCU,
    uiTrDepth uint,
    uiAbsPartIdx uint,
    bLumaOnly bool) {
    uiFullDepth := uint(pcCU.GetDepth1(0)) + uiTrDepth
    //uiTrMode := uint(pcCU.GetTransformIdx1(uiAbsPartIdx))
    //assert(  uiTrMode == uiTrDepth );
    uiLog2TrSize := uint(TLibCommon.G_aucConvertToBit[pcCU.GetSlice().GetSPS().GetMaxCUWidth()>>uiFullDepth]) + 2
    uiQTLayer := pcCU.GetSlice().GetSPS().GetQuadtreeTULog2MaxSize() - uiLog2TrSize

    bSkipChroma := false
    bChromaSame := false
    if !bLumaOnly && uiLog2TrSize == 2 {
        //assert( uiTrDepth > 0 );
        uiQPDiv := pcCU.GetPic().GetNumPartInCU() >> ((uint(pcCU.GetDepth1(0)) + uiTrDepth - 1) << 1)
        bSkipChroma = ((uiAbsPartIdx % uiQPDiv) != 0)
        bChromaSame = true
    }

    //===== copy transform coefficients =====
    uiNumCoeffY := (pcCU.GetSlice().GetSPS().GetMaxCUWidth() * pcCU.GetSlice().GetSPS().GetMaxCUHeight()) >> (uiFullDepth << 1)
    uiNumCoeffIncY := (pcCU.GetSlice().GetSPS().GetMaxCUWidth() * pcCU.GetSlice().GetSPS().GetMaxCUHeight()) >> (pcCU.GetSlice().GetSPS().GetMaxCUDepth() << 1)
    pcCoeffDstY := this.m_ppcQTTempCoeffY[uiQTLayer][(uiNumCoeffIncY * uiAbsPartIdx):]
    pcCoeffSrcY := this.m_pcQTTempTUCoeffY[:]

    for i := uint(0); i < uiNumCoeffY; i++ {
        pcCoeffDstY[i] = pcCoeffSrcY[i] //, sizeof( TLibCommon.TCoeff ) * uiNumCoeffY );
    }
    //#if ADAPTIVE_QP_SELECTION
    pcArlCoeffDstY := this.m_ppcQTTempArlCoeffY[uiQTLayer][(uiNumCoeffIncY * uiAbsPartIdx):]
    pcArlCoeffSrcY := this.m_ppcQTTempTUArlCoeffY[:]
    for i := uint(0); i < uiNumCoeffY; i++ {
        pcArlCoeffDstY[i] = pcArlCoeffSrcY[i] //, sizeof( int ) * uiNumCoeffY );
    }
    //#endif
    if !bLumaOnly && !bSkipChroma {
        var uiNumCoeffC uint
        if bChromaSame {
            uiNumCoeffC = uiNumCoeffY
        } else {
            uiNumCoeffC = uiNumCoeffY >> 2
        }
        uiNumCoeffIncC := uiNumCoeffIncY >> 2
        pcCoeffDstU := this.m_ppcQTTempCoeffCb[uiQTLayer][(uiNumCoeffIncC * uiAbsPartIdx):]
        pcCoeffDstV := this.m_ppcQTTempCoeffCr[uiQTLayer][(uiNumCoeffIncC * uiAbsPartIdx):]
        pcCoeffSrcU := this.m_pcQTTempTUCoeffCb[:]
        pcCoeffSrcV := this.m_pcQTTempTUCoeffCr[:]
        for i := uint(0); i < uiNumCoeffC; i++ {
            pcCoeffDstU[i] = pcCoeffSrcU[i] //, sizeof( TLibCommon.TCoeff ) * uiNumCoeffC );
            pcCoeffDstV[i] = pcCoeffSrcV[i] //, sizeof( TLibCommon.TCoeff ) * uiNumCoeffC );
        }
        //#if ADAPTIVE_QP_SELECTION
        pcArlCoeffDstU := this.m_ppcQTTempArlCoeffCb[uiQTLayer][(uiNumCoeffIncC * uiAbsPartIdx):]
        pcArlCoeffDstV := this.m_ppcQTTempArlCoeffCr[uiQTLayer][(uiNumCoeffIncC * uiAbsPartIdx):]
        pcArlCoeffSrcU := this.m_ppcQTTempTUArlCoeffCb[:]
        pcArlCoeffSrcV := this.m_ppcQTTempTUArlCoeffCr[:]
        for i := uint(0); i < uiNumCoeffC; i++ {
            pcArlCoeffDstU[i] = pcArlCoeffSrcU[i] // sizeof( int ) * uiNumCoeffC );
            pcArlCoeffDstV[i] = pcArlCoeffSrcV[i] // sizeof( int ) * uiNumCoeffC );
        }
        //#endif
    }

    //===== copy reconstruction =====
    this.m_pcQTTempTransformSkipTComYuv.CopyPartToPartLuma(&this.m_pcQTTempTComYuv[uiQTLayer], uiAbsPartIdx, 1<<uiLog2TrSize, 1<<uiLog2TrSize)

    if !bLumaOnly && !bSkipChroma {
        var uiLog2TrSizeChroma uint
        if bChromaSame {
            uiLog2TrSizeChroma = uiLog2TrSize
        } else {
            uiLog2TrSizeChroma = uiLog2TrSize - 1
        }
        this.m_pcQTTempTransformSkipTComYuv.CopyPartToPartChroma(&this.m_pcQTTempTComYuv[uiQTLayer], uiAbsPartIdx, 1<<uiLog2TrSizeChroma, 1<<uiLog2TrSizeChroma)
    }

    uiZOrder := pcCU.GetZorderIdxInCU() + uiAbsPartIdx
    piRecIPred := pcCU.GetPic().GetPicYuvRec().GetLumaAddr2(int(pcCU.GetAddr()), int(uiZOrder))
    uiRecIPredStride := uint(pcCU.GetPic().GetPicYuvRec().GetStride())
    piRecQt := this.m_pcQTTempTComYuv[uiQTLayer].GetLumaAddr1(uiAbsPartIdx)
    uiRecQtStride := uint(this.m_pcQTTempTComYuv[uiQTLayer].GetStride())
    uiWidth := uint(uint(pcCU.GetWidth1(0))) >> uiTrDepth
    uiHeight := uint(pcCU.GetHeight1(0)) >> uiTrDepth
    pRecQt := piRecQt
    pRecIPred := piRecIPred
    for uiY := uint(0); uiY < uiHeight; uiY++ {
        for uiX := uint(0); uiX < uiWidth; uiX++ {
            pRecIPred[uiY*uiRecIPredStride+uiX] = pRecQt[uiY*uiRecQtStride+uiX]
        }
        //pRecQt = pRecQt[uiRecQtStride:]
        //pRecIPred = pRecIPred[uiRecIPredStride:]
    }

    if !bLumaOnly && !bSkipChroma {
        piRecIPred = pcCU.GetPic().GetPicYuvRec().GetCbAddr2(int(pcCU.GetAddr()), int(uiZOrder))
        piRecQt = this.m_pcQTTempTComYuv[uiQTLayer].GetCbAddr1(uiAbsPartIdx)
        pRecQt = piRecQt
        pRecIPred = piRecIPred
        for uiY := uint(0); uiY < uiHeight; uiY++ {
            for uiX := uint(0); uiX < uiWidth; uiX++ {
                pRecIPred[uiY*uiRecIPredStride+uiX] = pRecQt[uiY*uiRecQtStride+uiX]
            }
            //pRecQt = pRecQt[uiRecQtStride:]
            //pRecIPred = pRecIPred[uiRecIPredStride:]
        }

        piRecIPred = pcCU.GetPic().GetPicYuvRec().GetCrAddr2(int(pcCU.GetAddr()), int(uiZOrder))
        piRecQt = this.m_pcQTTempTComYuv[uiQTLayer].GetCrAddr1(uiAbsPartIdx)
        pRecQt = piRecQt
        pRecIPred = piRecIPred
        for uiY := uint(0); uiY < uiHeight; uiY++ {
            for uiX := uint(0); uiX < uiWidth; uiX++ {
                pRecIPred[uiY*uiRecIPredStride+uiX] = pRecQt[uiY*uiRecQtStride+uiX]
            }
            //pRecQt = pRecQt[uiRecQtStride:]
            //pRecIPred = pRecIPred[uiRecIPredStride:]
        }
    }
}

func (this *TEncSearch) xStoreIntraResultChromaQT(pcCU *TLibCommon.TComDataCU,
    uiTrDepth uint,
    uiAbsPartIdx uint,
    stateU0V1Both2 uint) {
    uiFullDepth := uint(pcCU.GetDepth1(0)) + uiTrDepth
    uiTrMode := uint(pcCU.GetTransformIdx1(uiAbsPartIdx))
    if uiTrMode == uiTrDepth {
        uiLog2TrSize := uint(TLibCommon.G_aucConvertToBit[pcCU.GetSlice().GetSPS().GetMaxCUWidth()>>uiFullDepth]) + 2
        uiQTLayer := pcCU.GetSlice().GetSPS().GetQuadtreeTULog2MaxSize() - uiLog2TrSize

        bChromaSame := false
        if uiLog2TrSize == 2 {
            //assert( uiTrDepth > 0 );
            uiTrDepth--
            uiQPDiv := pcCU.GetPic().GetNumPartInCU() >> ((uint(pcCU.GetDepth1(0)) + uiTrDepth) << 1)
            if (uiAbsPartIdx % uiQPDiv) != 0 {
                return
            }
            bChromaSame = true
        }

        //===== copy transform coefficients =====
        uiNumCoeffC := (pcCU.GetSlice().GetSPS().GetMaxCUWidth() * pcCU.GetSlice().GetSPS().GetMaxCUHeight()) >> (uiFullDepth << 1)
        if !bChromaSame {
            uiNumCoeffC >>= 2
        }
        uiNumCoeffIncC := (pcCU.GetSlice().GetSPS().GetMaxCUWidth() * pcCU.GetSlice().GetSPS().GetMaxCUHeight()) >> ((pcCU.GetSlice().GetSPS().GetMaxCUDepth() << 1) + 2)
        if stateU0V1Both2 == 0 || stateU0V1Both2 == 2 {
            pcCoeffSrcU := this.m_ppcQTTempCoeffCb[uiQTLayer][(uiNumCoeffIncC * uiAbsPartIdx):]
            pcCoeffDstU := this.m_pcQTTempTUCoeffCb[:]
            for i := uint(0); i < uiNumCoeffC; i++ {
                pcCoeffDstU[i] = pcCoeffSrcU[i] //, sizeof( TLibCommon.TCoeff ) * uiNumCoeffC );
            }

            //#if ADAPTIVE_QP_SELECTION
            pcArlCoeffSrcU := this.m_ppcQTTempArlCoeffCb[uiQTLayer][(uiNumCoeffIncC * uiAbsPartIdx):]
            pcArlCoeffDstU := this.m_ppcQTTempTUArlCoeffCb[:]
            for i := uint(0); i < uiNumCoeffC; i++ {
                pcArlCoeffDstU[i] = pcArlCoeffSrcU[i] //, sizeof( int ) * uiNumCoeffC );
            }
            //#endif
        }
        if stateU0V1Both2 == 1 || stateU0V1Both2 == 2 {
            pcCoeffSrcV := this.m_ppcQTTempCoeffCr[uiQTLayer][(uiNumCoeffIncC * uiAbsPartIdx):]
            pcCoeffDstV := this.m_pcQTTempTUCoeffCr[:]
            for i := uint(0); i < uiNumCoeffC; i++ {
                pcCoeffDstV[i] = pcCoeffSrcV[i] //, sizeof( TLibCommon.TCoeff ) * uiNumCoeffC );
            }
            //#if ADAPTIVE_QP_SELECTION
            pcArlCoeffSrcV := this.m_ppcQTTempArlCoeffCr[uiQTLayer][(uiNumCoeffIncC * uiAbsPartIdx):]
            pcArlCoeffDstV := this.m_ppcQTTempTUArlCoeffCr[:]
            for i := uint(0); i < uiNumCoeffC; i++ {
                pcArlCoeffDstV[i] = pcArlCoeffSrcV[i] //, sizeof( int ) * uiNumCoeffC );
            }
            //#endif
        }

        //===== copy reconstruction =====
        var uiLog2TrSizeChroma uint
        if bChromaSame {
            uiLog2TrSizeChroma = uiLog2TrSize
        } else {
            uiLog2TrSizeChroma = uiLog2TrSize - 1
        }
        this.m_pcQTTempTComYuv[uiQTLayer].CopyPartToPartChroma2(&this.m_pcQTTempTransformSkipTComYuv, uiAbsPartIdx, 1<<uiLog2TrSizeChroma, 1<<uiLog2TrSizeChroma, stateU0V1Both2)
    }
}

func (this *TEncSearch) xLoadIntraResultChromaQT(pcCU *TLibCommon.TComDataCU,
    uiTrDepth uint,
    uiAbsPartIdx uint,
    stateU0V1Both2 uint) {
    uiFullDepth := uint(pcCU.GetDepth1(0)) + uiTrDepth
    uiTrMode := uint(pcCU.GetTransformIdx1(uiAbsPartIdx))
    if uiTrMode == uiTrDepth {
        uiLog2TrSize := uint(TLibCommon.G_aucConvertToBit[pcCU.GetSlice().GetSPS().GetMaxCUWidth()>>uiFullDepth]) + 2
        uiQTLayer := pcCU.GetSlice().GetSPS().GetQuadtreeTULog2MaxSize() - uiLog2TrSize

        bChromaSame := false
        if uiLog2TrSize == 2 {
            //assert( uiTrDepth > 0 );
            uiTrDepth--
            uiQPDiv := pcCU.GetPic().GetNumPartInCU() >> ((uint(pcCU.GetDepth1(0)) + uiTrDepth) << 1)
            if (uiAbsPartIdx % uiQPDiv) != 0 {
                return
            }
            bChromaSame = true
        }

        //===== copy transform coefficients =====
        uiNumCoeffC := (pcCU.GetSlice().GetSPS().GetMaxCUWidth() * pcCU.GetSlice().GetSPS().GetMaxCUHeight()) >> (uiFullDepth << 1)
        if !bChromaSame {
            uiNumCoeffC >>= 2
        }
        uiNumCoeffIncC := (pcCU.GetSlice().GetSPS().GetMaxCUWidth() * pcCU.GetSlice().GetSPS().GetMaxCUHeight()) >> ((pcCU.GetSlice().GetSPS().GetMaxCUDepth() << 1) + 2)

        if stateU0V1Both2 == 0 || stateU0V1Both2 == 2 {
            pcCoeffDstU := this.m_ppcQTTempCoeffCb[uiQTLayer][(uiNumCoeffIncC * uiAbsPartIdx):]
            pcCoeffSrcU := this.m_pcQTTempTUCoeffCb[:]
            for i := uint(0); i < uiNumCoeffC; i++ {
                pcCoeffDstU[i] = pcCoeffSrcU[i] //, sizeof( TLibCommon.TCoeff ) * uiNumCoeffC );
            }
            //#if ADAPTIVE_QP_SELECTION
            pcArlCoeffDstU := this.m_ppcQTTempArlCoeffCb[uiQTLayer][(uiNumCoeffIncC * uiAbsPartIdx):]
            pcArlCoeffSrcU := this.m_ppcQTTempTUArlCoeffCb[:]
            for i := uint(0); i < uiNumCoeffC; i++ {
                pcArlCoeffDstU[i] = pcArlCoeffSrcU[i] //, sizeof( int ) * uiNumCoeffC );
            }
            //#endif
        }
        if stateU0V1Both2 == 1 || stateU0V1Both2 == 2 {
            pcCoeffDstV := this.m_ppcQTTempCoeffCr[uiQTLayer][(uiNumCoeffIncC * uiAbsPartIdx):]
            pcCoeffSrcV := this.m_pcQTTempTUCoeffCr[:]
            for i := uint(0); i < uiNumCoeffC; i++ {
                pcCoeffDstV[i] = pcCoeffSrcV[i] //, sizeof( TLibCommon.TCoeff ) * uiNumCoeffC );
            }
            //#if ADAPTIVE_QP_SELECTION
            pcArlCoeffDstV := this.m_ppcQTTempArlCoeffCr[uiQTLayer][(uiNumCoeffIncC * uiAbsPartIdx):]
            pcArlCoeffSrcV := this.m_ppcQTTempTUArlCoeffCr[:]
            for i := uint(0); i < uiNumCoeffC; i++ {
                pcArlCoeffDstV[i] = pcArlCoeffSrcV[i] //, sizeof( int ) * uiNumCoeffC );
            }
            //#endif
        }

        //===== copy reconstruction =====
        var uiLog2TrSizeChroma uint
        if bChromaSame {
            uiLog2TrSizeChroma = uiLog2TrSize
        } else {
            uiLog2TrSizeChroma = uiLog2TrSize - 1
        }
        this.m_pcQTTempTransformSkipTComYuv.CopyPartToPartChroma2(&this.m_pcQTTempTComYuv[uiQTLayer], uiAbsPartIdx, 1<<uiLog2TrSizeChroma, 1<<uiLog2TrSizeChroma, stateU0V1Both2)

        uiZOrder := pcCU.GetZorderIdxInCU() + uiAbsPartIdx
        uiWidth := uint(uint(pcCU.GetWidth1(0))) >> (uiTrDepth + 1)
        uiHeight := uint(pcCU.GetHeight1(0)) >> (uiTrDepth + 1)
        uiRecQtStride := uint(this.m_pcQTTempTComYuv[uiQTLayer].GetCStride())
        uiRecIPredStride := uint(pcCU.GetPic().GetPicYuvRec().GetCStride())

        if stateU0V1Both2 == 0 || stateU0V1Both2 == 2 {
            piRecIPred := pcCU.GetPic().GetPicYuvRec().GetCbAddr2(int(pcCU.GetAddr()), int(uiZOrder))
            piRecQt := this.m_pcQTTempTComYuv[uiQTLayer].GetCbAddr1(uiAbsPartIdx)
            pRecQt := piRecQt
            pRecIPred := piRecIPred
            for uiY := uint(0); uiY < uiHeight; uiY++ {
                for uiX := uint(0); uiX < uiWidth; uiX++ {
                    pRecIPred[uiY*uiRecIPredStride+uiX] = pRecQt[uiY*uiRecQtStride+uiX]
                }
                //pRecQt = pRecQt[uiRecQtStride:]
                //pRecIPred = pRecIPred[uiRecIPredStride:]
            }
        }
        if stateU0V1Both2 == 1 || stateU0V1Both2 == 2 {
            piRecIPred := pcCU.GetPic().GetPicYuvRec().GetCrAddr2(int(pcCU.GetAddr()), int(uiZOrder))
            piRecQt := this.m_pcQTTempTComYuv[uiQTLayer].GetCrAddr1(uiAbsPartIdx)
            pRecQt := piRecQt
            pRecIPred := piRecIPred
            for uiY := uint(0); uiY < uiHeight; uiY++ {
                for uiX := uint(0); uiX < uiWidth; uiX++ {
                    pRecIPred[uiY*uiRecIPredStride+uiX] = pRecQt[uiY*uiRecQtStride+uiX]
                }
                //pRecQt = pRecQt[uiRecQtStride:]
                //pRecIPred = pRecIPred[uiRecIPredStride:]
            }
        }
    }
}

// -------------------------------------------------------------------------------------------------------------------
// Inter search (AMP)
// -------------------------------------------------------------------------------------------------------------------

func (this *TEncSearch) xEstimateMvPredAMVP(pcCU *TLibCommon.TComDataCU,
    pcOrgYuv *TLibCommon.TComYuv,
    uiPartIdx uint,
    eRefPicList TLibCommon.RefPicList,
    iRefIdx int,
    rcMvPred *TLibCommon.TComMv,
    bFilled bool,
    puiDistBiP *uint) {
    //#if ZERO_MVD_EST
    //       puiDist *uint
    //#endif
    pcAMVPInfo := pcCU.GetCUMvField(eRefPicList).GetAMVPInfo()

    var cBestMv TLibCommon.TComMv //, cZeroMv, cMvPred

    iBestIdx := 0

    uiBestCost := uint(TLibCommon.MAX_INT)
    uiPartAddr := uint(0)
    var iRoiWidth, iRoiHeight, i int

    pcCU.GetPartIndexAndSize(uiPartIdx, &uiPartAddr, &iRoiWidth, &iRoiHeight)
    // Fill the MV Candidates
    if !bFilled {
        pcCU.FillMvpCand(uiPartIdx, uiPartAddr, eRefPicList, iRefIdx, pcAMVPInfo)
    }

    // initialize Mvp index & Mvp
    iBestIdx = 0
    cBestMv = pcAMVPInfo.MvCand[0]
    //#if !ZERO_MVD_EST
    if pcAMVPInfo.IN <= 1 {
        *rcMvPred = cBestMv

        pcCU.SetMVPIdxSubParts(iBestIdx, eRefPicList, uiPartAddr, uiPartIdx, uint(pcCU.GetDepth1(uiPartAddr)))
        pcCU.SetMVPNumSubParts(pcAMVPInfo.IN, eRefPicList, uiPartAddr, uiPartIdx, uint(pcCU.GetDepth1(uiPartAddr)))

        if pcCU.GetSlice().GetMvdL1ZeroFlag() && eRefPicList == TLibCommon.REF_PIC_LIST_1 {
            //#if ZERO_MVD_EST
            //      (*puiDistBiP) = xGetTemplateCost( pcCU, uiPartIdx, uiPartAddr, pcOrgYuv, &this.m_cYuvPredTemp, rcMvPred, 0, TLibCommon.AMVP_MAX_NUM_CANDS, eRefPicList, iRefIdx, iRoiWidth, iRoiHeight, uiDist );
            //#else
            (*puiDistBiP) = this.xGetTemplateCost(pcCU, uiPartIdx, uiPartAddr, pcOrgYuv, this.GetYuvPredTemp(), *rcMvPred, 0, TLibCommon.AMVP_MAX_NUM_CANDS, eRefPicList, iRefIdx, iRoiWidth, iRoiHeight)
            //#endif
        }
        return
    }
    //#endif
    if bFilled {
        //assert(pcCU.GetMVPIdx(eRefPicList,uiPartAddr) >= 0);
        *rcMvPred = pcAMVPInfo.MvCand[pcCU.GetMVPIdx2(eRefPicList, uiPartAddr)]
        return
    }

    this.GetYuvPredTemp().Clear()
    //#if ZERO_MVD_EST
    //  UInt uiDist;
    //#endif
    //-- Check Minimum Cost.
    for i = 0; i < pcAMVPInfo.IN; i++ {
        var uiTmpCost uint
        //#if ZERO_MVD_EST
        //    uiTmpCost = xGetTemplateCost( pcCU, uiPartIdx, uiPartAddr, pcOrgYuv, this.GetYuvPredTemp(), pcAMVPInfo->this.MvCand[i], i, TLibCommon.AMVP_MAX_NUM_CANDS, eRefPicList, iRefIdx, iRoiWidth, iRoiHeight, uiDist );
        //#else
        uiTmpCost = this.xGetTemplateCost(pcCU, uiPartIdx, uiPartAddr, pcOrgYuv, this.GetYuvPredTemp(), pcAMVPInfo.MvCand[i], i, TLibCommon.AMVP_MAX_NUM_CANDS, eRefPicList, iRefIdx, iRoiWidth, iRoiHeight)
        //#endif
        if uiBestCost > uiTmpCost {
            uiBestCost = uiTmpCost
            cBestMv = pcAMVPInfo.MvCand[i]
            iBestIdx = i
            (*puiDistBiP) = uiTmpCost
            //#if ZERO_MVD_EST
            //      (*puiDist) = uiDist;
            //#endif
        }
    }

    this.GetYuvPredTemp().Clear()

    // Setting Best MVP
    *rcMvPred = cBestMv
    pcCU.SetMVPIdxSubParts(iBestIdx, eRefPicList, uiPartAddr, uiPartIdx, uint(pcCU.GetDepth1(uiPartAddr)))
    pcCU.SetMVPNumSubParts(pcAMVPInfo.IN, eRefPicList, uiPartAddr, uiPartIdx, uint(pcCU.GetDepth1(uiPartAddr)))
    return
}

func (this *TEncSearch) xCheckBestMVP(pcCU *TLibCommon.TComDataCU,
    eRefPicList TLibCommon.RefPicList,
    cMv TLibCommon.TComMv,
    rcMvPred *TLibCommon.TComMv,
    riMVPIdx *int,
    ruiBits *uint,
    ruiCost *uint) {
    pcAMVPInfo := pcCU.GetCUMvField(eRefPicList).GetAMVPInfo()

    //assert(pcAMVPInfo->this.MvCand[riMVPIdx] == rcMvPred);

    if pcAMVPInfo.IN < 2 {
        return
    }

    this.m_pcRdCost.getMotionCost(true, 0)
    this.m_pcRdCost.setCostScale(0)

    iBestMVPIdx := *riMVPIdx

    this.m_pcRdCost.setPredictor(rcMvPred)
    iOrgMvBits := this.m_pcRdCost.getBits(int(cMv.GetHor()), int(cMv.GetVer()))
    iOrgMvBits += this.m_auiMVPIdxCost[*riMVPIdx][TLibCommon.AMVP_MAX_NUM_CANDS]
    iBestMvBits := iOrgMvBits

    for iMVPIdx := 0; iMVPIdx < pcAMVPInfo.IN; iMVPIdx++ {
        if iMVPIdx == *riMVPIdx {
            continue
        }

        this.m_pcRdCost.setPredictor(&pcAMVPInfo.MvCand[iMVPIdx])

        iMvBits := this.m_pcRdCost.getBits(int(cMv.GetHor()), int(cMv.GetVer()))
        iMvBits += this.m_auiMVPIdxCost[iMVPIdx][TLibCommon.AMVP_MAX_NUM_CANDS]

        if iMvBits < iBestMvBits {
            iBestMvBits = iMvBits
            iBestMVPIdx = iMVPIdx
        }
    }

    if iBestMVPIdx != *riMVPIdx { //if changed
        *rcMvPred = pcAMVPInfo.MvCand[iBestMVPIdx]

        *riMVPIdx = iBestMVPIdx
        uiOrgBits := *ruiBits
        *ruiBits = uiOrgBits - iOrgMvBits + iBestMvBits
        *ruiCost = (*ruiCost - this.m_pcRdCost.getCost1(uiOrgBits)) + this.m_pcRdCost.getCost1(*ruiBits)
    }
}

func (this *TEncSearch) xGetTemplateCost(pcCU *TLibCommon.TComDataCU,
    uiPartIdx uint,
    uiPartAddr uint,
    pcOrgYuv *TLibCommon.TComYuv,
    pcTemplateCand *TLibCommon.TComYuv,
    cMvCand TLibCommon.TComMv,
    iMVPIdx int,
    iMVPNum int,
    eRefPicList TLibCommon.RefPicList,
    iRefIdx int,
    iSizeX int,
    iSizeY int) uint {
    //#if ZERO_MVD_EST
    //UInt&       ruiDist
    //#endif

    uiCost := uint(TLibCommon.MAX_INT)

    pcPicYuvRef := pcCU.GetSlice().GetRefPic(eRefPicList, iRefIdx).GetPicYuvRec()

    pcCU.ClipMv(&cMvCand)

    // prediction pattern
    if pcCU.GetSlice().GetPPS().GetUseWP() && pcCU.GetSlice().GetSliceType() == TLibCommon.P_SLICE {
        this.XPredInterLumaBlk(pcCU, pcPicYuvRef, uiPartAddr, &cMvCand, iSizeX, iSizeY, pcTemplateCand, true)
    } else {
        this.XPredInterLumaBlk(pcCU, pcPicYuvRef, uiPartAddr, &cMvCand, iSizeX, iSizeY, pcTemplateCand, false)
    }

    if pcCU.GetSlice().GetPPS().GetUseWP() && pcCU.GetSlice().GetSliceType() == TLibCommon.P_SLICE {
        this.XWeightedPredictionUni(pcCU, pcTemplateCand, uiPartAddr, iSizeX, iSizeY, eRefPicList, pcTemplateCand, iRefIdx)
    }

    // calc distortion

    //#if WEIGHTED_CHROMA_DISTORTION
    uiCost = this.m_pcRdCost.getDistPart(TLibCommon.G_bitDepthY, pcTemplateCand.GetLumaAddr1(uiPartAddr), int(pcTemplateCand.GetStride()), pcOrgYuv.GetLumaAddr1(uiPartAddr), int(pcOrgYuv.GetStride()), uint(iSizeX), uint(iSizeY), TLibCommon.TEXT_LUMA, TLibCommon.DF_SAD)
    //#else
    //   uiCost = this.m_pcRdCost.getDistPart(TLibCommon.G_bitDepthY, pcTemplateCand.GetLumaAddr(uiPartAddr), pcTemplateCand.GetStride(), pcOrgYuv.GetLumaAddr(uiPartAddr), pcOrgYuv.GetStride(), iSizeX, iSizeY, DF_SAD );
    //#endif
    uiCost = uint(this.m_pcRdCost.calcRdCost(this.m_auiMVPIdxCost[iMVPIdx][iMVPNum], uiCost, false, TLibCommon.DF_SAD))

    return uiCost
}

func (this *TEncSearch) xCopyAMVPInfo(pSrc, pDst *TLibCommon.AMVPInfo) {
    pDst.IN = pSrc.IN
    for i := 0; i < pSrc.IN; i++ {
        pDst.MvCand[i] = pSrc.MvCand[i]
    }
}

func (this *TEncSearch) xGetMvpIdxBits(iIdx, iNum int) uint {
    //assert(iIdx >= 0 && iNum >= 0 && iIdx < iNum);

    if iNum == 1 {
        return 0
    }

    uiLength := uint(1)
    iTemp := iIdx
    if iTemp == 0 {
        return uiLength
    }

    bCodeLast := (iNum-1 > iTemp)

    uiLength += uint(iTemp - 1)

    if bCodeLast {
        uiLength++
    }

    return uiLength
}

func (this *TEncSearch) xGetBlkBits(eCUMode TLibCommon.PartSize, bPSlice bool, iPartIdx int, uiLastMode uint, uiBlkBit []uint) {
    if eCUMode == TLibCommon.SIZE_2Nx2N {
        if !bPSlice {
            uiBlkBit[0] = 3
        } else {
            uiBlkBit[0] = 1
        }
        uiBlkBit[1] = 3
        uiBlkBit[2] = 5
    } else if (eCUMode == TLibCommon.SIZE_2NxN || eCUMode == TLibCommon.SIZE_2NxnU) || eCUMode == TLibCommon.SIZE_2NxnD {
        var aauiMbBits = [2][3][3]uint{{{0, 0, 3}, {0, 0, 0}, {0, 0, 0}}, {{5, 7, 7}, {7, 5, 7}, {9 - 3, 9 - 3, 9 - 3}}}
        if bPSlice {
            uiBlkBit[0] = 3
            uiBlkBit[1] = 0
            uiBlkBit[2] = 0
        } else {
            for i := uint(0); i < 3; i++ {
                uiBlkBit[i] = aauiMbBits[iPartIdx][uiLastMode][i] //, 3*sizeof(UInt) );
            }
        }
    } else if (eCUMode == TLibCommon.SIZE_Nx2N || eCUMode == TLibCommon.SIZE_nLx2N) || eCUMode == TLibCommon.SIZE_nRx2N {
        var aauiMbBits = [2][3][3]uint{{{0, 2, 3}, {0, 0, 0}, {0, 0, 0}}, {{5, 7, 7}, {7 - 2, 7 - 2, 9 - 2}, {9 - 3, 9 - 3, 9 - 3}}}
        if bPSlice {
            uiBlkBit[0] = 3
            uiBlkBit[1] = 0
            uiBlkBit[2] = 0
        } else {
            for i := uint(0); i < 3; i++ {
                uiBlkBit[i] = aauiMbBits[iPartIdx][uiLastMode][i] //, 3*sizeof(UInt) );
            }
        }
    } else if eCUMode == TLibCommon.SIZE_NxN {
        if !bPSlice {
            uiBlkBit[0] = 3
        } else {
            uiBlkBit[0] = 1
        }
        uiBlkBit[1] = 3
        uiBlkBit[2] = 5
    } else {
        println("Wrong!")
        //assert( 0 );
    }
}

func (this *TEncSearch) xMergeEstimation(pcCU *TLibCommon.TComDataCU,
    pcYuvOrg *TLibCommon.TComYuv,
    iPUIdx int,
    uiInterDir *uint,
    pacMvField []TLibCommon.TComMvField,
    uiMergeIndex *uint,
    ruiCost *uint,
    cMvFieldNeighbours []TLibCommon.TComMvField,
    uhInterDirNeighbours []byte,
    numValidMergeCand *int) {
    uiAbsPartIdx := uint(0)
    iWidth := 0
    iHeight := 0

    pcCU.GetPartIndexAndSize(uint(iPUIdx), &uiAbsPartIdx, &iWidth, &iHeight)
    uiDepth := uint(pcCU.GetDepth1(uiAbsPartIdx))
    partSize := pcCU.GetPartitionSize1(0)
    if pcCU.GetSlice().GetPPS().GetLog2ParallelMergeLevelMinus2() != 0 && partSize != TLibCommon.SIZE_2Nx2N && uint(pcCU.GetWidth1(0)) <= 8 {
        pcCU.SetPartSizeSubParts(TLibCommon.SIZE_2Nx2N, 0, uiDepth)
        if iPUIdx == 0 {
            pcCU.GetInterMergeCandidates(0, 0, cMvFieldNeighbours, uhInterDirNeighbours, numValidMergeCand, -1)
        }
        pcCU.SetPartSizeSubParts(partSize, 0, uiDepth)
    } else {
        pcCU.GetInterMergeCandidates(uiAbsPartIdx, uint(iPUIdx), cMvFieldNeighbours, uhInterDirNeighbours, numValidMergeCand, -1)
    }
    this.xRestrictBipredMergeCand(pcCU, uint(iPUIdx), cMvFieldNeighbours[:], uhInterDirNeighbours[:], *numValidMergeCand)

    *ruiCost = uint(TLibCommon.MAX_UINT)
    for uiMergeCand := uint(0); uiMergeCand < uint(*numValidMergeCand); uiMergeCand++ {
        uiCostCand := TLibCommon.MAX_UINT
        uiBitsCand := uint(0)

        ePartSize := pcCU.GetPartitionSize1(0)

        pcCU.GetCUMvField(TLibCommon.REF_PIC_LIST_0).SetAllMvField(&cMvFieldNeighbours[0+2*uiMergeCand], ePartSize, int(uiAbsPartIdx), 0, iPUIdx)
        pcCU.GetCUMvField(TLibCommon.REF_PIC_LIST_1).SetAllMvField(&cMvFieldNeighbours[1+2*uiMergeCand], ePartSize, int(uiAbsPartIdx), 0, iPUIdx)

        this.xGetInterPredictionError(pcCU, pcYuvOrg, iPUIdx, &uiCostCand, this.m_pcEncCfg.GetUseHADME())
        uiBitsCand = uiMergeCand + 1
        if uiMergeCand == this.m_pcEncCfg.GetMaxNumMergeCand()-1 {
            uiBitsCand--
        }
        uiCostCand = uiCostCand + this.m_pcRdCost.getCost1(uiBitsCand)
        if uiCostCand < *ruiCost {
            *ruiCost = uiCostCand
            pacMvField[0] = cMvFieldNeighbours[0+2*uiMergeCand]
            pacMvField[1] = cMvFieldNeighbours[1+2*uiMergeCand]
            *uiInterDir = uint(uhInterDirNeighbours[uiMergeCand])
            *uiMergeIndex = uiMergeCand
        }
    }

}

func (this *TEncSearch) xRestrictBipredMergeCand(pcCU *TLibCommon.TComDataCU,
    puIdx uint,
    mvFieldNeighbours []TLibCommon.TComMvField,
    interDirNeighbours []byte,
    numValidMergeCand int) {
    if pcCU.IsBipredRestriction(puIdx) {
        for mergeCand := 0; mergeCand < numValidMergeCand; mergeCand++ {
            if interDirNeighbours[mergeCand] == 3 {
                interDirNeighbours[mergeCand] = 1
                mvFieldNeighbours[(mergeCand<<1)+1].SetMvField(*TLibCommon.NewTComMv(0, 0), -1)
            }
        }
    }
}

// -------------------------------------------------------------------------------------------------------------------
// motion estimation
// -------------------------------------------------------------------------------------------------------------------

func (this *TEncSearch) xMotionEstimation(pcCU *TLibCommon.TComDataCU,
    pcYuvOrg *TLibCommon.TComYuv,
    iPartIdx int,
    eRefPicList TLibCommon.RefPicList,
    pcMvPred *TLibCommon.TComMv,
    iRefIdxPred int,
    rcMv *TLibCommon.TComMv,
    ruiBits *uint,
    ruiCost *uint,
    bBi bool) {
    
    //fmt.Printf("Enter xMotionEstimation with iPartIdx%d,eRefPicList%d,iRefIdxPred%d,ruiBits%d,ruiCost%d,bBi%d\n", iPartIdx,eRefPicList,iRefIdxPred,*ruiBits,*ruiCost,TLibCommon.B2U(bBi));
    
    var uiPartAddr uint
    var iRoiWidth int
    var iRoiHeight int

    var cMvHalf, cMvQter, cMvSrchRngLT, cMvSrchRngRB TLibCommon.TComMv

    pcYuv := pcYuvOrg
    this.m_iSearchRange = this.m_aaiAdaptSR[eRefPicList][iRefIdxPred]

    var iSrchRng int
    if bBi {
        iSrchRng = this.m_bipredSearchRange
    } else {
        iSrchRng = this.m_iSearchRange
    }
    pcPatternKey := pcCU.GetPattern()
    fWeight := float64(1.0)

    pcCU.GetPartIndexAndSize(uint(iPartIdx), &uiPartAddr, &iRoiWidth, &iRoiHeight)

    if bBi {
        pcYuvOther := this.GetYuvPred(1 - int(eRefPicList))
        pcYuv = this.GetYuvPredTemp()

        pcYuvOrg.CopyPartToPartYuv(pcYuv, uiPartAddr, uint(iRoiWidth), uint(iRoiHeight))

        pcYuv.RemoveHighFreq(pcYuvOther, uiPartAddr, uint(iRoiWidth), uint(iRoiHeight))

        fWeight = 0.5
    }

    //  Search key pattern initialization
    pcPatternKey.InitPattern(pcYuv.GetLumaAddr1(uiPartAddr),
        pcYuv.GetCbAddr1(uiPartAddr),
        pcYuv.GetCrAddr1(uiPartAddr),
        iRoiWidth,
        iRoiHeight,
        int(pcYuv.GetStride()), 
        0, 0, 0,
        0, 0)

    pPicYuvRec := pcCU.GetSlice().GetRefPic(eRefPicList, iRefIdxPred).GetPicYuvRec()
    piRefY := pPicYuvRec.GetBufY();//pcCU.GetSlice().GetRefPic(eRefPicList, iRefIdxPred).GetPicYuvRec().GetLumaAddr2(int(pcCU.GetAddr()), int(pcCU.GetZorderIdxInCU()+uiPartAddr))
    iOffset := pPicYuvRec.GetLumaMarginY()*pPicYuvRec.GetStride()+pPicYuvRec.GetLumaMarginX() + 
    		   pPicYuvRec.GetCuOffsetY()[pcCU.GetAddr()] + pPicYuvRec.GetBuOffsetY()[TLibCommon.G_auiZscanToRaster[pcCU.GetZorderIdxInCU()+uiPartAddr]]
    iRefStride := pPicYuvRec.GetStride()

    cMvPred := *pcMvPred

    if bBi {
        this.xSetSearchRange(pcCU, rcMv, iSrchRng, &cMvSrchRngLT, &cMvSrchRngRB)
    } else {
        this.xSetSearchRange(pcCU, &cMvPred, iSrchRng, &cMvSrchRngLT, &cMvSrchRngRB)
    }

    this.m_pcRdCost.getMotionCost(true, 0)

    this.m_pcRdCost.setPredictor(pcMvPred)
    this.m_pcRdCost.setCostScale(2)

    this.setWpScalingDistParam(pcCU, iRefIdxPred, eRefPicList)
    //  Do integer search
    if this.m_iFastSearch == 0 || bBi {
    	this.xPatternSearch	   (pcCU, pcPatternKey, piRefY, iOffset, iRefStride, &cMvSrchRngLT, &cMvSrchRngRB, rcMv, ruiCost)
    } else {
        *rcMv = *pcMvPred
        this.xPatternSearchFast(pcCU, pcPatternKey, piRefY, iOffset, iRefStride, &cMvSrchRngLT, &cMvSrchRngRB, rcMv, ruiCost)
    }

    this.m_pcRdCost.getMotionCost(true, 0)
    this.m_pcRdCost.setCostScale(1)

	this.xPatternSearchFracDIF(pcCU, pcPatternKey, piRefY, iOffset, iRefStride, rcMv, &cMvHalf, &cMvQter, ruiCost, bBi)
	
	//fmt.Printf("rcMv=(%d,%d), cMvHalf=(%d,%d), cMvQter=(%d,%d), ruiCost=%d\n",rcMv.GetHor(),rcMv.GetVer(),cMvHalf.GetHor(),cMvHalf.GetVer(),cMvQter.GetHor(),cMvQter.GetVer(), *ruiCost);
	  
    this.m_pcRdCost.setCostScale(0)
    rcMv.ShiftMv(2) // <<= 2    
    cMvHalf.ShiftMv(1)//(cMvHalf <<= 1);    
    rcMv.AddMv(cMvHalf)     
    rcMv.AddMv(cMvQter)
    
    uiMvBits := this.m_pcRdCost.getBits(int(rcMv.GetHor()), int(rcMv.GetVer()))
	//fmt.Printf("uiMvBits=%d\n",uiMvBits);

    *ruiBits += uiMvBits
    *ruiCost = uint(math.Floor(fWeight*(float64(*ruiCost)-float64(this.m_pcRdCost.getCost1(uiMvBits)))) + float64(this.m_pcRdCost.getCost1(*ruiBits)))
	
	//os.Exit(-1);
	//fmt.Printf("Exit xMotionEstimation with mv_x%d,mv_y%d,ruiBits%d,ruiCost%d\n", rcMv.GetHor(),rcMv.GetVer(),*ruiBits,*ruiCost);
}

func (this *TEncSearch) xTZSearch(pcCU *TLibCommon.TComDataCU,
    pcPatternKey *TLibCommon.TComPattern,
    piRefY []TLibCommon.Pel,
    iOffset int,
    iRefStride int,
    pcMvSrchRngLT *TLibCommon.TComMv,
    pcMvSrchRngRB *TLibCommon.TComMv,
    rcMv *TLibCommon.TComMv,
    ruiSAD *uint) {
    
    
    iSrchRngHorLeft := int(pcMvSrchRngLT.GetHor())
    iSrchRngHorRight := int(pcMvSrchRngRB.GetHor())
    iSrchRngVerTop := int(pcMvSrchRngLT.GetVer())
    iSrchRngVerBottom := int(pcMvSrchRngRB.GetVer())

    //TZ_SEARCH_CONFIGURATION

    uiSearchRange := this.m_iSearchRange
    //fmt.Printf("xTZSearch rcMV0=(%d,%d)\n", rcMv.GetHor(), rcMv.GetVer());
    pcCU.ClipMv(rcMv)
    //fmt.Printf("xTZSearch rcMV1=(%d,%d)\n", rcMv.GetHor(), rcMv.GetVer());
    rcMv.Set(rcMv.GetHor()>>2, rcMv.GetVer()>>2) // >>= 2
    //fmt.Printf("xTZSearch rcMV2=(%d,%d)\n", rcMv.GetHor(), rcMv.GetVer());
    // init TZSearchStruct
    var cStruct IntTZSearchStruct
    cStruct.iYStride = iRefStride
    cStruct.piRefY = piRefY
    cStruct.iOffset = iOffset
    cStruct.uiBestSad = uint(TLibCommon.MAX_UINT)
	
    // set rcMv (Median predictor) as start point and as best point
    this.xTZSearchHelp(pcPatternKey, &cStruct, int(rcMv.GetHor()), int(rcMv.GetVer()), 0, 0)

    // test whether one of PRED_A, PRED_B, PRED_C MV is better start point than Median predictor
    if bTestOtherPredictedMV {
        for index := 0; index < 3; index++ {
            cMv := this.m_acMvPredictors[index]
            pcCU.ClipMv(&cMv)
            cMv.Set(cMv.GetHor()/4, cMv.GetVer()/4) // >>= 2
            //            cMv >>= 2
            this.xTZSearchHelp(pcPatternKey, &cStruct, int(cMv.GetHor()), int(cMv.GetVer()), 0, 0)
        }
    }

    // test whether zero Mv is better start point than Median predictor
    if bTestZeroVector {
        this.xTZSearchHelp(pcPatternKey, &cStruct, 0, 0, 0, 0)
    }

    // start search
    iDist := 0
    iStartX := cStruct.iBestX
    iStartY := cStruct.iBestY

    // first search
    for iDist = 1; iDist <= int(uiSearchRange); iDist *= 2 {
        if bFirstSearchDiamond {
            this.xTZ8PointDiamondSearch(pcPatternKey, &cStruct, pcMvSrchRngLT, pcMvSrchRngRB, iStartX, iStartY, iDist)
        } else {
            this.xTZ8PointSquareSearch(pcPatternKey, &cStruct, pcMvSrchRngLT, pcMvSrchRngRB, iStartX, iStartY, iDist)
        }

        if bFirstSearchStop && (cStruct.uiBestRound >= uiFirstSearchRounds) { // stop criterion
            break
        }
    }

    // test whether zero Mv is a better start point than Median predictor
    if bTestZeroVectorStart && ((cStruct.iBestX != 0) || (cStruct.iBestY != 0)) {
        this.xTZSearchHelp(pcPatternKey, &cStruct, 0, 0, 0, 0)
        if (cStruct.iBestX == 0) && (cStruct.iBestY == 0) {
            // test its neighborhood
            for iDist = 1; iDist <= int(uiSearchRange); iDist *= 2 {
                this.xTZ8PointDiamondSearch(pcPatternKey, &cStruct, pcMvSrchRngLT, pcMvSrchRngRB, 0, 0, iDist)
                if bTestZeroVectorStop && (cStruct.uiBestRound > 0) { // stop criterion
                    break
                }
            }
        }
    }

    // calculate only 2 missing points instead 8 points if cStruct.uiBestDistance == 1
    if cStruct.uiBestDistance == 1 {
        cStruct.uiBestDistance = 0
        this.xTZ2PointSearch(pcPatternKey, &cStruct, pcMvSrchRngLT, pcMvSrchRngRB)
    }

    // raster search if distance is too big
    if bEnableRasterSearch && (((int)(cStruct.uiBestDistance) > iRaster) || bAlwaysRasterSearch) {
        cStruct.uiBestDistance = iRaster
        for iStartY = iSrchRngVerTop; iStartY <= iSrchRngVerBottom; iStartY += iRaster {
            for iStartX = iSrchRngHorLeft; iStartX <= iSrchRngHorRight; iStartX += iRaster {
                this.xTZSearchHelp(pcPatternKey, &cStruct, iStartX, iStartY, 0, iRaster)
            }
        }
    }

    // raster refinement
    if bRasterRefinementEnable && cStruct.uiBestDistance > 0 {
        for cStruct.uiBestDistance > 0 {
            iStartX = cStruct.iBestX
            iStartY = cStruct.iBestY
            if cStruct.uiBestDistance > 1 {
                cStruct.uiBestDistance >>= 1
                iDist = int(cStruct.uiBestDistance)
                if bRasterRefinementDiamond {
                    this.xTZ8PointDiamondSearch(pcPatternKey, &cStruct, pcMvSrchRngLT, pcMvSrchRngRB, iStartX, iStartY, iDist)
                } else {
                    this.xTZ8PointSquareSearch(pcPatternKey, &cStruct, pcMvSrchRngLT, pcMvSrchRngRB, iStartX, iStartY, iDist)
                }
            }

            // calculate only 2 missing points instead 8 points if cStruct.uiBestDistance == 1
            if cStruct.uiBestDistance == 1 {
                cStruct.uiBestDistance = 0
                if cStruct.ucPointNr != 0 {
                    this.xTZ2PointSearch(pcPatternKey, &cStruct, pcMvSrchRngLT, pcMvSrchRngRB)
                }
            }
        }
    }

    // start refinement
    if bStarRefinementEnable && cStruct.uiBestDistance > 0 {
        for cStruct.uiBestDistance > 0 {
            iStartX = cStruct.iBestX
            iStartY = cStruct.iBestY
            cStruct.uiBestDistance = 0
            cStruct.ucPointNr = 0
            for iDist = 1; iDist < int(uiSearchRange)+1; iDist *= 2 {
                if bStarRefinementDiamond {
                    this.xTZ8PointDiamondSearch(pcPatternKey, &cStruct, pcMvSrchRngLT, pcMvSrchRngRB, iStartX, iStartY, iDist)
                } else {
                    this.xTZ8PointSquareSearch(pcPatternKey, &cStruct, pcMvSrchRngLT, pcMvSrchRngRB, iStartX, iStartY, iDist)
                }
                if bStarRefinementStop && (cStruct.uiBestRound >= uiStarRefinementRounds) { // stop criterion
                    break
                }
            }

            // calculate only 2 missing points instead 8 points if cStrukt.uiBestDistance == 1
            if cStruct.uiBestDistance == 1 {
                cStruct.uiBestDistance = 0
                if cStruct.ucPointNr != 0 {
                    this.xTZ2PointSearch(pcPatternKey, &cStruct, pcMvSrchRngLT, pcMvSrchRngRB)
                }
            }
        }
    }

    // write out best match
    rcMv.Set(int16(cStruct.iBestX), int16(cStruct.iBestY))
    *ruiSAD = cStruct.uiBestSad - this.m_pcRdCost.getCost2(cStruct.iBestX, cStruct.iBestY)
}

func (this *TEncSearch) xSetSearchRange(pcCU *TLibCommon.TComDataCU,
    cMvPred *TLibCommon.TComMv,
    iSrchRng int,
    rcMvSrchRngLT *TLibCommon.TComMv,
    rcMvSrchRngRB *TLibCommon.TComMv) {
    iMvShift := uint(2)
    cTmpMvPred := cMvPred
    pcCU.ClipMv(cTmpMvPred)

    rcMvSrchRngLT.SetHor(cTmpMvPred.GetHor() - int16(iSrchRng<<iMvShift))
    rcMvSrchRngLT.SetVer(cTmpMvPred.GetVer() - int16(iSrchRng<<iMvShift))

    rcMvSrchRngRB.SetHor(cTmpMvPred.GetHor() + int16(iSrchRng<<iMvShift))
    rcMvSrchRngRB.SetVer(cTmpMvPred.GetVer() + int16(iSrchRng<<iMvShift))
    pcCU.ClipMv(rcMvSrchRngLT)
    pcCU.ClipMv(rcMvSrchRngRB)

    rcMvSrchRngLT.Set(rcMvSrchRngLT.GetHor()>>iMvShift, rcMvSrchRngLT.GetVer()>>iMvShift)
    rcMvSrchRngRB.Set(rcMvSrchRngRB.GetHor()>>iMvShift, rcMvSrchRngRB.GetVer()>>iMvShift)
    //rcMvSrchRngRB >>= iMvShift
}

func (this *TEncSearch) xPatternSearchFast(pcCU *TLibCommon.TComDataCU,
    pcPatternKey *TLibCommon.TComPattern,
    piRefY []TLibCommon.Pel,
    iOffset int,
    iRefStride int,
    pcMvSrchRngLT *TLibCommon.TComMv,
    pcMvSrchRngRB *TLibCommon.TComMv,
    rcMv *TLibCommon.TComMv,
    ruiSAD *uint) {

    this.m_acMvPredictors[0] = pcCU.GetMvPredLeft()
    this.m_acMvPredictors[1] = pcCU.GetMvPredAbove()
    this.m_acMvPredictors[2] = pcCU.GetMvPredAboveRight()

    switch this.m_iFastSearch {
    case 1:
        this.xTZSearch(pcCU, pcPatternKey, piRefY, iOffset, iRefStride, pcMvSrchRngLT, pcMvSrchRngRB, rcMv, ruiSAD)
    default:
    }
}

func (this *TEncSearch) xPatternSearch(pcCU *TLibCommon.TComDataCU,
	pcPatternKey *TLibCommon.TComPattern,
    piRefY []TLibCommon.Pel,
    iOffset int,
    iRefStride int,
    pcMvSrchRngLT *TLibCommon.TComMv,
    pcMvSrchRngRB *TLibCommon.TComMv,
    rcMv *TLibCommon.TComMv,
    ruiSAD *uint) {
    iSrchRngHorLeft := int16(pcMvSrchRngLT.GetHor())
    iSrchRngHorRight := int16(pcMvSrchRngRB.GetHor())
    iSrchRngVerTop := int16(pcMvSrchRngLT.GetVer())
    iSrchRngVerBottom := int16(pcMvSrchRngRB.GetVer())

    var uiSad uint
    uiSadBest := uint(TLibCommon.MAX_UINT)
    iBestX := int16(0)
    iBestY := int16(0)

    var piRefSrch []TLibCommon.Pel

    //-- jclee for using the SAD function pointer
    this.m_pcRdCost.setDistParam2(pcPatternKey, piRefY, iRefStride, &this.m_cDistParam)

    // fast encoder decision: use subsampled SAD for integer ME
    if this.m_pcEncCfg.GetUseFastEnc() {
        if this.m_cDistParam.iRows > 8 {
            this.m_cDistParam.iSubShift = 1
        }
    }

    //piRefY = piRefY[(iSrchRngVerTop * iRefStride):]
    for mv_y := iSrchRngVerTop; mv_y <= iSrchRngVerBottom; mv_y++ {
        for mv_x := iSrchRngHorLeft; mv_x <= iSrchRngHorRight; mv_x++ {
        	candMV := TLibCommon.NewTComMv(mv_x, mv_y)
        	pcCU.ClipMv(candMV)
        	
            //  find min. distortion position
            piRefSrch = piRefY[iOffset + int(candMV.GetVer()) * iRefStride + int(candMV.GetHor()):]
            this.m_cDistParam.pCur = piRefSrch

            this.setDistParamComp(0)

            this.m_cDistParam.bitDepth = TLibCommon.G_bitDepthY
            uiSad = this.m_cDistParam.DistFunc(&this.m_cDistParam)

            // motion cost
            uiSad += this.m_pcRdCost.getCost2(int(candMV.GetHor()), int(candMV.GetVer()))

            if uiSad < uiSadBest {
                uiSadBest = uiSad
                iBestX = candMV.GetHor()
                iBestY = candMV.GetVer()
            }
        }
        //piRefY = piRefY[iRefStride:]
    }

    rcMv.Set(int16(iBestX), int16(iBestY))

    *ruiSAD = uiSadBest - this.m_pcRdCost.getCost2(int(iBestX), int(iBestY))
    return
}

func (this *TEncSearch) xPatternSearchFracDIF(pcCU *TLibCommon.TComDataCU,
    pcPatternKey *TLibCommon.TComPattern,
    piRefY []TLibCommon.Pel,
    iOffset int,
    iRefStride int,
    pcMvInt *TLibCommon.TComMv,
    rcMvHalf *TLibCommon.TComMv,
    rcMvQter *TLibCommon.TComMv,
    ruiCost *uint,
    biPred bool) {
    
    //fmt.Printf("ruiCost=%d, pcMvInt=(%d,%d)\n", *ruiCost, pcMvInt.GetHor(), pcMvInt.GetVer());
    
    //  Reference pattern initialization (integer scale)
    var cPatternRoi TLibCommon.TComPattern
    //iOffset := int(pcMvInt.GetHor()) + int(pcMvInt.GetVer())*iRefStride
    iOffset += int(pcMvInt.GetHor()) + int(pcMvInt.GetVer())*iRefStride
    cPatternRoi.InitPattern(piRefY, nil, nil,
        pcPatternKey.GetROIYWidth(),
        pcPatternKey.GetROIYHeight(),
        iRefStride,
        iOffset, 0, 0,
        0, 0)

    //  Half-pel refinement
    this.xExtDIFUpSamplingH(&cPatternRoi, biPred)

    *rcMvHalf = *pcMvInt
    rcMvHalf.ShiftMv(1) // <<= 1 // for mv-cost
    baseRefMv := TLibCommon.NewTComMv(0, 0)
    cPatternRoi.InitPattern(piRefY[iOffset:], nil, nil,
        pcPatternKey.GetROIYWidth(),
        pcPatternKey.GetROIYHeight(),
        iRefStride,
        iOffset, 0, 0,
        0, 0)
    *ruiCost = this.xPatternRefinement(pcPatternKey, *baseRefMv, 2, rcMvHalf)
	//fmt.Printf("ruiCost=%d, rcMvHalf=(%d,%d)\n", *ruiCost, rcMvHalf.GetHor(), rcMvHalf.GetVer());
    
    this.m_pcRdCost.setCostScale(0)

	cPatternRoi.InitPattern(piRefY, nil, nil,
        pcPatternKey.GetROIYWidth(),
        pcPatternKey.GetROIYHeight(),
        iRefStride,
        iOffset, 0, 0,
        0, 0)
    this.xExtDIFUpSamplingQ(&cPatternRoi, *rcMvHalf, biPred)
    *baseRefMv = *rcMvHalf
    baseRefMv.ShiftMv(1) // <<= 1

    *rcMvQter = *pcMvInt
    rcMvQter.ShiftMv(1) // <<= 1 // for mv-cost
    rcMvQter.AddMv(*rcMvHalf)
    rcMvQter.ShiftMv(1) // <<= 1
    cPatternRoi.InitPattern(piRefY[iOffset:], nil, nil,
        pcPatternKey.GetROIYWidth(),
        pcPatternKey.GetROIYHeight(),
        iRefStride,
        iOffset, 0, 0,
        0, 0)
    *ruiCost = this.xPatternRefinement(pcPatternKey, *baseRefMv, 1, rcMvQter)
    
    //fmt.Printf("ruiCost=%d, rcMvQter=(%d,%d)\n", *ruiCost, rcMvQter.GetHor(), rcMvQter.GetVer());
}

func (this *TEncSearch) xExtDIFUpSamplingH(pattern *TLibCommon.TComPattern, biPred bool) {
    width := pattern.GetROIYWidth()
    height := pattern.GetROIYHeight()
    srcStride := pattern.GetPatternLStride()
    srcOffset := pattern.GetPatternLOffset()

    intStride := int(this.GetFilteredBlockTmp(0).GetStride())
    dstStride := int(this.GetFilteredBlock(0, 0).GetStride())
    var intPtr, dstPtr []TLibCommon.Pel
    filterSize := TLibCommon.NTAPS_LUMA
    halfFilterSize := (filterSize >> 1)
    srcPtr := pattern.GetROIY()[srcOffset-halfFilterSize*srcStride-1-(halfFilterSize-1)*1:]

    this.GetIf().FilterHorLuma(srcPtr, srcStride, this.GetFilteredBlockTmp(0).GetLumaAddr(), intStride, width+1, height+filterSize, 0, false)
    this.GetIf().FilterHorLuma(srcPtr, srcStride, this.GetFilteredBlockTmp(2).GetLumaAddr(), intStride, width+1, height+filterSize, 2, false)
  
    intPtr = this.GetFilteredBlockTmp(0).GetLumaAddr();//[halfFilterSize*intStride+1:]
    dstPtr = this.GetFilteredBlock(0, 0).GetLumaAddr()
    this.GetIf().FilterVerLuma(intPtr[-(halfFilterSize-1)*intStride+halfFilterSize*intStride+1:], intStride, dstPtr, dstStride, width+0, height+0, 0, false, true)

    intPtr = this.GetFilteredBlockTmp(0).GetLumaAddr();//[(halfFilterSize-1)*intStride+1:]
    dstPtr = this.GetFilteredBlock(2, 0).GetLumaAddr()
    this.GetIf().FilterVerLuma(intPtr[-(halfFilterSize-1)*intStride+(halfFilterSize-1)*intStride+1:], intStride, dstPtr, dstStride, width+0, height+1, 2, false, true)

    intPtr = this.GetFilteredBlockTmp(2).GetLumaAddr();//[halfFilterSize*intStride:]
    dstPtr = this.GetFilteredBlock(0, 2).GetLumaAddr()
    this.GetIf().FilterVerLuma(intPtr[-(halfFilterSize-1)*intStride+halfFilterSize*intStride:], intStride, dstPtr, dstStride, width+1, height+0, 0, false, true)

    intPtr = this.GetFilteredBlockTmp(2).GetLumaAddr();//[(halfFilterSize-1)*intStride:]
    dstPtr = this.GetFilteredBlock(2, 2).GetLumaAddr()
    this.GetIf().FilterVerLuma(intPtr[-(halfFilterSize-1)*intStride+(halfFilterSize-1)*intStride:], intStride, dstPtr, dstStride, width+1, height+1, 2, false, true)
}

func (this *TEncSearch) xExtDIFUpSamplingQ(pattern *TLibCommon.TComPattern, halfPelRef TLibCommon.TComMv, biPred bool) {
    width := pattern.GetROIYWidth()
    height := pattern.GetROIYHeight()
    srcStride := pattern.GetPatternLStride()
    srcOffset := pattern.GetPatternLOffset()

    var srcPtr []TLibCommon.Pel
    intStride := int(this.GetFilteredBlockTmp(0).GetStride())
    dstStride := int(this.GetFilteredBlock(0, 0).GetStride())
    var intPtr, dstPtr []TLibCommon.Pel
    filterSize := TLibCommon.NTAPS_LUMA
    halfFilterSize := (filterSize >> 1)

    var extHeight int
    if halfPelRef.GetVer() == 0 {
        extHeight = int(height + filterSize)
    } else {
        extHeight = int(height + filterSize - 1)
    }

    // Horizontal filter 1/4
    srcPtr = pattern.GetROIY()[srcOffset-halfFilterSize*srcStride-1-(halfFilterSize-1)*1:]
    intPtr = this.GetFilteredBlockTmp(1).GetLumaAddr()
    if halfPelRef.GetVer() > 0 {
        srcPtr = srcPtr[srcStride:]
    }
    if halfPelRef.GetHor() >= 0 {
        srcPtr = srcPtr[1:]
    }
    this.GetIf().FilterHorLuma(srcPtr, srcStride, intPtr, int(intStride), width, extHeight, 1, false)

    // Horizontal filter 3/4
    srcPtr = pattern.GetROIY()[srcOffset-halfFilterSize*srcStride-1-(halfFilterSize-1)*1:]
    intPtr = this.GetFilteredBlockTmp(3).GetLumaAddr()
    if halfPelRef.GetVer() > 0 {
        srcPtr = srcPtr[srcStride:]
    }
    if halfPelRef.GetHor() > 0 {
        srcPtr = srcPtr[1:]
    }
    this.GetIf().FilterHorLuma(srcPtr, srcStride, intPtr, int(intStride), width, extHeight, 3, false)

    // Generate @ 1,1
    intPtr = this.GetFilteredBlockTmp(1).GetLumaAddr();//[(halfFilterSize-1)*int(intStride):]
    dstPtr = this.GetFilteredBlock(1, 1).GetLumaAddr()
    if halfPelRef.GetVer() == 0 {
        intPtr = intPtr[intStride:]
    }
    this.GetIf().FilterVerLuma(intPtr[-(halfFilterSize-1)*intStride+(halfFilterSize-1)*intStride:], int(intStride), dstPtr, int(dstStride), width, height, 1, false, true)

    // Generate @ 3,1
    intPtr = this.GetFilteredBlockTmp(1).GetLumaAddr();//[(halfFilterSize-1)*int(intStride):]
    dstPtr = this.GetFilteredBlock(3, 1).GetLumaAddr()
    this.GetIf().FilterVerLuma(intPtr[-(halfFilterSize-1)*intStride+(halfFilterSize-1)*intStride:], int(intStride), dstPtr, int(dstStride), width, height, 3, false, true)

    if halfPelRef.GetVer() != 0 {
        // Generate @ 2,1
        intPtr = this.GetFilteredBlockTmp(1).GetLumaAddr();//[(halfFilterSize-1)*int(intStride):]
        dstPtr = this.GetFilteredBlock(2, 1).GetLumaAddr()
        if halfPelRef.GetVer() == 0 {
            intPtr = intPtr[intStride:]
        }
        this.GetIf().FilterVerLuma(intPtr[-(halfFilterSize-1)*intStride+(halfFilterSize-1)*intStride:], int(intStride), dstPtr, int(dstStride), width, height, 2, false, true)

        // Generate @ 2,3
        intPtr = this.GetFilteredBlockTmp(3).GetLumaAddr();//[(halfFilterSize-1)*intStride:]
        dstPtr = this.GetFilteredBlock(2, 3).GetLumaAddr()
        if halfPelRef.GetVer() == 0 {
            intPtr = intPtr[intStride:]
        }
        this.GetIf().FilterVerLuma(intPtr[-(halfFilterSize-1)*intStride+(halfFilterSize-1)*intStride:], intStride, dstPtr, dstStride, width, height, 2, false, true)
    } else {
        // Generate @ 0,1
        intPtr = this.GetFilteredBlockTmp(1).GetLumaAddr();//[halfFilterSize*intStride:]
        dstPtr = this.GetFilteredBlock(0, 1).GetLumaAddr()
        this.GetIf().FilterVerLuma(intPtr[-(halfFilterSize-1)*intStride+halfFilterSize*intStride:], intStride, dstPtr, dstStride, width, height, 0, false, true)

        // Generate @ 0,3
        intPtr = this.GetFilteredBlockTmp(3).GetLumaAddr();//[halfFilterSize*intStride:]
        dstPtr = this.GetFilteredBlock(0, 3).GetLumaAddr()
        this.GetIf().FilterVerLuma(intPtr[-(halfFilterSize-1)*intStride+halfFilterSize*intStride:], intStride, dstPtr, dstStride, width, height, 0, false, true)
    }

    if halfPelRef.GetHor() != 0 {
        // Generate @ 1,2
        intPtr = this.GetFilteredBlockTmp(2).GetLumaAddr();//[(halfFilterSize-1)*intStride:]
        dstPtr = this.GetFilteredBlock(1, 2).GetLumaAddr()
        if halfPelRef.GetHor() > 0 {
            intPtr = intPtr[1:]
        }
        if halfPelRef.GetVer() >= 0 {
            intPtr = intPtr[intStride:]
        }
        this.GetIf().FilterVerLuma(intPtr[-(halfFilterSize-1)*intStride+(halfFilterSize-1)*intStride:], intStride, dstPtr, dstStride, width, height, 1, false, true)

        // Generate @ 3,2
        intPtr = this.GetFilteredBlockTmp(2).GetLumaAddr();//[(halfFilterSize-1)*intStride:]
        dstPtr = this.GetFilteredBlock(3, 2).GetLumaAddr()
        if halfPelRef.GetHor() > 0 {
            intPtr = intPtr[1:]
        }
        if halfPelRef.GetVer() > 0 {
            intPtr = intPtr[intStride:]
        }
        this.GetIf().FilterVerLuma(intPtr[-(halfFilterSize-1)*intStride+(halfFilterSize-1)*intStride:], intStride, dstPtr, dstStride, width, height, 3, false, true)
    } else {
        // Generate @ 1,0
        intPtr = this.GetFilteredBlockTmp(0).GetLumaAddr();//[(halfFilterSize-1)*intStride+1:]
        dstPtr = this.GetFilteredBlock(1, 0).GetLumaAddr()
        if halfPelRef.GetVer() >= 0 {
            intPtr = intPtr[intStride:]
        }
        this.GetIf().FilterVerLuma(intPtr[-(halfFilterSize-1)*intStride+(halfFilterSize-1)*intStride+1:], intStride, dstPtr, dstStride, width, height, 1, false, true)

        // Generate @ 3,0
        intPtr = this.GetFilteredBlockTmp(0).GetLumaAddr();//[(halfFilterSize-1)*intStride+1:]
        dstPtr = this.GetFilteredBlock(3, 0).GetLumaAddr()
        if halfPelRef.GetVer() > 0 {
            intPtr = intPtr[intStride:]
        }
        this.GetIf().FilterVerLuma(intPtr[-(halfFilterSize-1)*intStride+(halfFilterSize-1)*intStride+1:], intStride, dstPtr, dstStride, width, height, 3, false, true)
    }

    // Generate @ 1,3
    intPtr = this.GetFilteredBlockTmp(3).GetLumaAddr();//[(halfFilterSize-1)*intStride:]
    dstPtr = this.GetFilteredBlock(1, 3).GetLumaAddr()
    if halfPelRef.GetVer() == 0 {
        intPtr = intPtr[intStride:]
    }
    this.GetIf().FilterVerLuma(intPtr[-(halfFilterSize-1)*intStride+(halfFilterSize-1)*intStride:], intStride, dstPtr, dstStride, width, height, 1, false, true)

    // Generate @ 3,3
    intPtr = this.GetFilteredBlockTmp(3).GetLumaAddr();//[(halfFilterSize-1)*intStride:]
    dstPtr = this.GetFilteredBlock(3, 3).GetLumaAddr()
    this.GetIf().FilterVerLuma(intPtr[-(halfFilterSize-1)*intStride+(halfFilterSize-1)*intStride:], intStride, dstPtr, dstStride, width, height, 3, false, true)
}

// -------------------------------------------------------------------------------------------------------------------
// T & Q & Q-1 & T-1
// -------------------------------------------------------------------------------------------------------------------

func (this *TEncSearch) xEncodeResidualQT(pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint, bSubdivAndCbf bool, eType TLibCommon.TextType) {
    //assert( pcCU.GetDepth1( 0 ) == pcCU.GetDepth1( uiAbsPartIdx ) );
    uiCurrTrMode := uiDepth - uint(pcCU.GetDepth1(0))
    uiTrMode := uint(pcCU.GetTransformIdx1(uiAbsPartIdx))

    bSubdiv := uiCurrTrMode != uiTrMode

    uiLog2TrSize := uint(TLibCommon.G_aucConvertToBit[pcCU.GetSlice().GetSPS().GetMaxCUWidth()>>uiDepth]) + 2

    if bSubdivAndCbf && uiLog2TrSize <= pcCU.GetSlice().GetSPS().GetQuadtreeTULog2MaxSize() && uiLog2TrSize > pcCU.GetQuadtreeTULog2MinSizeInCU(uiAbsPartIdx) {
        this.m_pcEntropyCoder.encodeTransformSubdivFlag(uint(TLibCommon.B2U(bSubdiv)), 5-uiLog2TrSize)
    }

    {
        //assert( pcCU.GetPredictionMode(uiAbsPartIdx) != MODE_INTRA );
        if bSubdivAndCbf {
            bFirstCbfOfCU := uiCurrTrMode == 0
            if bFirstCbfOfCU || uiLog2TrSize > 2 {
                if bFirstCbfOfCU || pcCU.GetCbf3(uiAbsPartIdx, TLibCommon.TEXT_CHROMA_U, uiCurrTrMode-1) != 0 {
                    this.m_pcEntropyCoder.encodeQtCbf(pcCU, uiAbsPartIdx, TLibCommon.TEXT_CHROMA_U, uiCurrTrMode)
                }
                if bFirstCbfOfCU || pcCU.GetCbf3(uiAbsPartIdx, TLibCommon.TEXT_CHROMA_V, uiCurrTrMode-1) != 0 {
                    this.m_pcEntropyCoder.encodeQtCbf(pcCU, uiAbsPartIdx, TLibCommon.TEXT_CHROMA_V, uiCurrTrMode)
                }
            } else if uiLog2TrSize == 2 {
                //assert( pcCU.GetCbf( uiAbsPartIdx, TLibCommon.TEXT_CHROMA_U, uiCurrTrMode ) == pcCU.GetCbf( uiAbsPartIdx, TLibCommon.TEXT_CHROMA_U, uiCurrTrMode - 1 ) );
                //assert( pcCU.GetCbf( uiAbsPartIdx, TLibCommon.TEXT_CHROMA_V, uiCurrTrMode ) == pcCU.GetCbf( uiAbsPartIdx, TLibCommon.TEXT_CHROMA_V, uiCurrTrMode - 1 ) );
            }
        }
    }

    if !bSubdiv {
        uiNumCoeffPerAbsPartIdxIncrement := pcCU.GetSlice().GetSPS().GetMaxCUWidth() * pcCU.GetSlice().GetSPS().GetMaxCUHeight() >> (pcCU.GetSlice().GetSPS().GetMaxCUDepth() << 1)
        //assert( 16 == uiNumCoeffPerAbsPartIdxIncrement ); // check
        uiQTTempAccessLayer := pcCU.GetSlice().GetSPS().GetQuadtreeTULog2MaxSize() - uiLog2TrSize
        pcCoeffCurrY := this.m_ppcQTTempCoeffY[uiQTTempAccessLayer][uiNumCoeffPerAbsPartIdxIncrement*uiAbsPartIdx:]
        pcCoeffCurrU := this.m_ppcQTTempCoeffCb[uiQTTempAccessLayer][(uiNumCoeffPerAbsPartIdxIncrement * uiAbsPartIdx >> 2):]
        pcCoeffCurrV := this.m_ppcQTTempCoeffCr[uiQTTempAccessLayer][(uiNumCoeffPerAbsPartIdxIncrement * uiAbsPartIdx >> 2):]

        bCodeChroma := true
        uiTrModeC := uiTrMode
        uiLog2TrSizeC := uiLog2TrSize - 1
        if uiLog2TrSize == 2 {
            uiLog2TrSizeC++
            uiTrModeC--
            uiQPDiv := pcCU.GetPic().GetNumPartInCU() >> ((uint(pcCU.GetDepth1(0)) + uiTrModeC) << 1)
            bCodeChroma = ((uiAbsPartIdx % uiQPDiv) == 0)
        }

        if bSubdivAndCbf {
            this.m_pcEntropyCoder.encodeQtCbf(pcCU, uiAbsPartIdx, TLibCommon.TEXT_LUMA, uiTrMode)
        } else {
            if eType == TLibCommon.TEXT_LUMA && pcCU.GetCbf3(uiAbsPartIdx, TLibCommon.TEXT_LUMA, uiTrMode) != 0 {
                trWidth := uint(1) << uiLog2TrSize
                trHeight := uint(1) << uiLog2TrSize
                this.m_pcEntropyCoder.encodeCoeffNxN(pcCU, pcCoeffCurrY, uiAbsPartIdx, trWidth, trHeight, uiDepth, TLibCommon.TEXT_LUMA)
            }
            if bCodeChroma {
                trWidth := uint(1) << uiLog2TrSizeC
                trHeight := uint(1) << uiLog2TrSizeC
                if eType == TLibCommon.TEXT_CHROMA_U && pcCU.GetCbf3(uiAbsPartIdx, TLibCommon.TEXT_CHROMA_U, uiTrMode) != 0 {
                    this.m_pcEntropyCoder.encodeCoeffNxN(pcCU, pcCoeffCurrU, uiAbsPartIdx, trWidth, trHeight, uiDepth, TLibCommon.TEXT_CHROMA_U)
                }
                if eType == TLibCommon.TEXT_CHROMA_V && pcCU.GetCbf3(uiAbsPartIdx, TLibCommon.TEXT_CHROMA_V, uiTrMode) != 0 {
                    this.m_pcEntropyCoder.encodeCoeffNxN(pcCU, pcCoeffCurrV, uiAbsPartIdx, trWidth, trHeight, uiDepth, TLibCommon.TEXT_CHROMA_V)
                }
            }
        }
    } else {
        if bSubdivAndCbf || pcCU.GetCbf3(uiAbsPartIdx, eType, uiCurrTrMode) != 0 {
            uiQPartNumSubdiv := pcCU.GetPic().GetNumPartInCU() >> ((uiDepth + 1) << 1)
            for ui := uint(0); ui < 4; ui++ {
                this.xEncodeResidualQT(pcCU, uiAbsPartIdx+ui*uiQPartNumSubdiv, uiDepth+1, bSubdivAndCbf, eType)
            }
        }
    }
}

func (this *TEncSearch) xEstimateResidualQT(pcCU *TLibCommon.TComDataCU, uiQuadrant, uiAbsPartIdx, absTUPartIdx uint, pcResi *TLibCommon.TComYuv, uiDepth uint, rdCost *float64, ruiBits, ruiDist, puiZeroDist *uint) {
    uiTrMode := uiDepth - uint(pcCU.GetDepth1(0))

    //assert( pcCU.GetDepth1( 0 ) == pcCU.GetDepth1( uiAbsPartIdx ) );
    uiLog2TrSize := uint(TLibCommon.G_aucConvertToBit[pcCU.GetSlice().GetSPS().GetMaxCUWidth()>>uiDepth] + 2)

    SplitFlag := ((pcCU.GetSlice().GetSPS().GetQuadtreeTUMaxDepthInter() == 1) && pcCU.GetPredictionMode1(uiAbsPartIdx) == TLibCommon.MODE_INTER && (pcCU.GetPartitionSize1(uiAbsPartIdx) != TLibCommon.SIZE_2Nx2N))
    var bCheckFull bool
    if SplitFlag && uiDepth == uint(pcCU.GetDepth1(uiAbsPartIdx)) && (uiLog2TrSize > pcCU.GetQuadtreeTULog2MinSizeInCU(uiAbsPartIdx)) {
        bCheckFull = false
    } else {
        bCheckFull = (uiLog2TrSize <= pcCU.GetSlice().GetSPS().GetQuadtreeTULog2MaxSize())
    }
    bCheckSplit := uiLog2TrSize > pcCU.GetQuadtreeTULog2MinSizeInCU(uiAbsPartIdx)

    //assert( bCheckFull || bCheckSplit );

    bCodeChroma := true
    uiTrModeC := uiTrMode
    uiLog2TrSizeC := uiLog2TrSize - 1
    if uiLog2TrSize == 2 {
        uiLog2TrSizeC++
        uiTrModeC--
        uiQPDiv := pcCU.GetPic().GetNumPartInCU() >> ((uint(pcCU.GetDepth1(0)) + uiTrModeC) << 1)
        bCodeChroma = ((uiAbsPartIdx % uiQPDiv) == 0)
    }

    uiSetCbf := 1 << uiTrMode
    // code full block
    dSingleCost := TLibCommon.MAX_DOUBLE
    uiSingleBits := uint(0)
    uiSingleDist := uint(0)
    uiAbsSumY := uint(0)
    uiAbsSumU := uint(0)
    uiAbsSumV := uint(0)
    var uiBestTransformMode [3]uint // = {0};

    if this.m_bUseSBACRD {
        this.m_pcRDGoOnSbacCoder.store(this.m_pppcRDSbacCoder[uiDepth][TLibCommon.CI_QT_TRAFO_ROOT])
    }

    if bCheckFull {
        uiNumCoeffPerAbsPartIdxIncrement := pcCU.GetSlice().GetSPS().GetMaxCUWidth() * pcCU.GetSlice().GetSPS().GetMaxCUHeight() >> (pcCU.GetSlice().GetSPS().GetMaxCUDepth() << 1)
        uiQTTempAccessLayer := pcCU.GetSlice().GetSPS().GetQuadtreeTULog2MaxSize() - uiLog2TrSize
        pcCoeffCurrY := this.m_ppcQTTempCoeffY[uiQTTempAccessLayer][uiNumCoeffPerAbsPartIdxIncrement*uiAbsPartIdx:]
        pcCoeffCurrU := this.m_ppcQTTempCoeffCb[uiQTTempAccessLayer][(uiNumCoeffPerAbsPartIdxIncrement * uiAbsPartIdx >> 2):]
        pcCoeffCurrV := this.m_ppcQTTempCoeffCr[uiQTTempAccessLayer][(uiNumCoeffPerAbsPartIdxIncrement * uiAbsPartIdx >> 2):]
        //#if ADAPTIVE_QP_SELECTION
        pcArlCoeffCurrY := this.m_ppcQTTempArlCoeffY[uiQTTempAccessLayer][uiNumCoeffPerAbsPartIdxIncrement*uiAbsPartIdx:]
        pcArlCoeffCurrU := this.m_ppcQTTempArlCoeffCb[uiQTTempAccessLayer][(uiNumCoeffPerAbsPartIdxIncrement * uiAbsPartIdx >> 2):]
        pcArlCoeffCurrV := this.m_ppcQTTempArlCoeffCr[uiQTTempAccessLayer][(uiNumCoeffPerAbsPartIdxIncrement * uiAbsPartIdx >> 2):]
        //#endif

        trWidth := uint(0)
        trHeight := uint(0)
        trWidthC := uint(0)
        trHeightC := uint(0)
        absTUPartIdxC := uiAbsPartIdx

        trHeight = 1 << uiLog2TrSize
        trWidth = trHeight
        trHeightC = 1 << uiLog2TrSizeC
        trWidthC = trHeightC
        pcCU.SetTrIdxSubParts(uiDepth-uint(pcCU.GetDepth1(0)), uiAbsPartIdx, uiDepth)
        minCostY := TLibCommon.MAX_DOUBLE
        minCostU := TLibCommon.MAX_DOUBLE
        minCostV := TLibCommon.MAX_DOUBLE
        checkTransformSkipY := pcCU.GetSlice().GetPPS().GetUseTransformSkip() && trWidth == 4 && trHeight == 4
        checkTransformSkipUV := pcCU.GetSlice().GetPPS().GetUseTransformSkip() && trWidthC == 4 && trHeightC == 4

        checkTransformSkipY = checkTransformSkipY && (!pcCU.IsLosslessCoded(0))
        checkTransformSkipUV = checkTransformSkipUV && (!pcCU.IsLosslessCoded(0))

        pcCU.SetTransformSkipSubParts4(false, TLibCommon.TEXT_LUMA, uiAbsPartIdx, uiDepth)
        if bCodeChroma {
            pcCU.SetTransformSkipSubParts4(false, TLibCommon.TEXT_CHROMA_U, uiAbsPartIdx, uint(pcCU.GetDepth1(0))+uiTrModeC)
            pcCU.SetTransformSkipSubParts4(false, TLibCommon.TEXT_CHROMA_V, uiAbsPartIdx, uint(pcCU.GetDepth1(0))+uiTrModeC)
        }

        if this.m_pcEncCfg.GetUseRDOQ() {
            this.m_pcEntropyCoder.estimateBit(this.m_pcTrQuant.GetEstBitsSbac(), int(trWidth), int(trHeight), TLibCommon.TEXT_LUMA)
        }

        this.m_pcTrQuant.SetQPforQuant(int(pcCU.GetQP1(0)), TLibCommon.TEXT_LUMA, pcCU.GetSlice().GetSPS().GetQpBDOffsetY(), 0)

        //#if RDOQ_CHROMA_LAMBDA
        this.m_pcTrQuant.SelectLambda(TLibCommon.TEXT_LUMA)
        //#endif
        this.m_pcTrQuant.TransformNxN(pcCU, pcResi.GetLumaAddr1(absTUPartIdx), pcResi.GetStride(), pcCoeffCurrY,
            //#if ADAPTIVE_QP_SELECTION
            pcArlCoeffCurrY,
            //#endif
            uint(trWidth), uint(trHeight), &uiAbsSumY, TLibCommon.TEXT_LUMA, uiAbsPartIdx, false)

        if uiAbsSumY != 0 {
            pcCU.SetCbfSubParts4(byte(uiSetCbf), TLibCommon.TEXT_LUMA, uiAbsPartIdx, uiDepth)
        } else {
            pcCU.SetCbfSubParts4(0, TLibCommon.TEXT_LUMA, uiAbsPartIdx, uiDepth)
        }

        if bCodeChroma {
            if this.m_pcEncCfg.GetUseRDOQ() {
                this.m_pcEntropyCoder.estimateBit(this.m_pcTrQuant.GetEstBitsSbac(), int(trWidthC), int(trHeightC), TLibCommon.TEXT_CHROMA)
            }

            curChromaQpOffset := pcCU.GetSlice().GetPPS().GetChromaCbQpOffset() + pcCU.GetSlice().GetSliceQpDeltaCb()
            this.m_pcTrQuant.SetQPforQuant(int(pcCU.GetQP1(0)), TLibCommon.TEXT_CHROMA, pcCU.GetSlice().GetSPS().GetQpBDOffsetC(), curChromaQpOffset)

            //#if RDOQ_CHROMA_LAMBDA
            this.m_pcTrQuant.SelectLambda(TLibCommon.TEXT_CHROMA)
            //#endif

            this.m_pcTrQuant.TransformNxN(pcCU, pcResi.GetCbAddr1(absTUPartIdxC), pcResi.GetCStride(), pcCoeffCurrU,
                //#if ADAPTIVE_QP_SELECTION
                pcArlCoeffCurrU,
                //#endif
                uint(trWidthC), uint(trHeightC), &uiAbsSumU, TLibCommon.TEXT_CHROMA_U, uiAbsPartIdx, false)

            curChromaQpOffset = pcCU.GetSlice().GetPPS().GetChromaCrQpOffset() + pcCU.GetSlice().GetSliceQpDeltaCr()
            this.m_pcTrQuant.SetQPforQuant(int(pcCU.GetQP1(0)), TLibCommon.TEXT_CHROMA, pcCU.GetSlice().GetSPS().GetQpBDOffsetC(), curChromaQpOffset)
            this.m_pcTrQuant.TransformNxN(pcCU, pcResi.GetCrAddr1(absTUPartIdxC), pcResi.GetCStride(), pcCoeffCurrV,
                //#if ADAPTIVE_QP_SELECTION
                pcArlCoeffCurrV,
                //#endif
                uint(trWidthC), uint(trHeightC), &uiAbsSumV, TLibCommon.TEXT_CHROMA_V, uiAbsPartIdx, false)
            if uiAbsSumU != 0 {
                pcCU.SetCbfSubParts4(byte(uiSetCbf), TLibCommon.TEXT_CHROMA_U, uiAbsPartIdx, uint(pcCU.GetDepth1(0))+uiTrModeC)
            } else {
                pcCU.SetCbfSubParts4(0, TLibCommon.TEXT_CHROMA_U, uiAbsPartIdx, uint(pcCU.GetDepth1(0))+uiTrModeC)
            }
            if uiAbsSumV != 0 {
                pcCU.SetCbfSubParts4(byte(uiSetCbf), TLibCommon.TEXT_CHROMA_V, uiAbsPartIdx, uint(pcCU.GetDepth1(0))+uiTrModeC)
            } else {
                pcCU.SetCbfSubParts4(0, TLibCommon.TEXT_CHROMA_V, uiAbsPartIdx, uint(pcCU.GetDepth1(0))+uiTrModeC)
            }
        }

        this.m_pcEntropyCoder.resetBits()

        this.m_pcEntropyCoder.encodeQtCbf(pcCU, uiAbsPartIdx, TLibCommon.TEXT_LUMA, uiTrMode)

        this.m_pcEntropyCoder.encodeCoeffNxN(pcCU, pcCoeffCurrY, uiAbsPartIdx, uint(trWidth), uint(trHeight), uiDepth, TLibCommon.TEXT_LUMA)
        uiSingleBitsY := this.m_pcEntropyCoder.getNumberOfWrittenBits()

        uiSingleBitsU := uint(0)
        uiSingleBitsV := uint(0)
        if bCodeChroma {
            this.m_pcEntropyCoder.encodeQtCbf(pcCU, uiAbsPartIdx, TLibCommon.TEXT_CHROMA_U, uiTrMode)

            this.m_pcEntropyCoder.encodeCoeffNxN(pcCU, pcCoeffCurrU, uiAbsPartIdx, uint(trWidthC), uint(trHeightC), uiDepth, TLibCommon.TEXT_CHROMA_U)
            uiSingleBitsU = this.m_pcEntropyCoder.getNumberOfWrittenBits() - uiSingleBitsY

            this.m_pcEntropyCoder.encodeQtCbf(pcCU, uiAbsPartIdx, TLibCommon.TEXT_CHROMA_V, uiTrMode)

            this.m_pcEntropyCoder.encodeCoeffNxN(pcCU, pcCoeffCurrV, uiAbsPartIdx, uint(trWidthC), uint(trHeightC), uiDepth, TLibCommon.TEXT_CHROMA_V)
            uiSingleBitsV = this.m_pcEntropyCoder.getNumberOfWrittenBits() - (uiSingleBitsY + uiSingleBitsU)
        }

        uiNumSamplesLuma := 1 << (uiLog2TrSize << 1)
        uiNumSamplesChro := 1 << (uiLog2TrSizeC << 1)

        for i := int(0); i < uiNumSamplesLuma; i++ {
            this.m_pTempPel[i] = 0 //, sizeof( TLibCommon.Pel ) * uiNumSamplesLuma ); // not necessary needed for inside of recursion (only at the beginning)
        }
        uiDistY := this.m_pcRdCost.getDistPart(TLibCommon.G_bitDepthY, this.m_pTempPel, int(trWidth), pcResi.GetLumaAddr1(absTUPartIdx), int(pcResi.GetStride()), uint(trWidth), uint(trHeight), TLibCommon.TEXT_LUMA, TLibCommon.DF_SSE) // initialized with zero residual destortion

        if puiZeroDist != nil {
            *puiZeroDist += uiDistY
        }
        if uiAbsSumY != 0 {
            pcResiCurrY := this.m_pcQTTempTComYuv[uiQTTempAccessLayer].GetLumaAddr1(absTUPartIdx)

            this.m_pcTrQuant.SetQPforQuant(int(pcCU.GetQP1(0)), TLibCommon.TEXT_LUMA, pcCU.GetSlice().GetSPS().GetQpBDOffsetY(), 0)

            scalingListType := 3 + TLibCommon.G_eTTable[int(TLibCommon.TEXT_LUMA)]
            //assert(scalingListType < 6);
            this.m_pcTrQuant.InvtransformNxN(pcCU.GetCUTransquantBypass1(uiAbsPartIdx), TLibCommon.TEXT_LUMA, TLibCommon.REG_DCT, pcResiCurrY, this.m_pcQTTempTComYuv[uiQTTempAccessLayer].GetStride(), pcCoeffCurrY, uint(trWidth), uint(trHeight), scalingListType, false) //this is for inter mode only

            uiNonzeroDistY := this.m_pcRdCost.getDistPart(TLibCommon.G_bitDepthY, this.m_pcQTTempTComYuv[uiQTTempAccessLayer].GetLumaAddr1(absTUPartIdx), int(this.m_pcQTTempTComYuv[uiQTTempAccessLayer].GetStride()),
                pcResi.GetLumaAddr1(absTUPartIdx), int(pcResi.GetStride()), trWidth, trHeight, TLibCommon.TEXT_LUMA, TLibCommon.DF_SSE)
            if pcCU.IsLosslessCoded(0) {
                uiDistY = uiNonzeroDistY
            } else {
                singleCostY := this.m_pcRdCost.calcRdCost(uiSingleBitsY, uiNonzeroDistY, false, TLibCommon.DF_DEFAULT)
                this.m_pcEntropyCoder.resetBits()
                this.m_pcEntropyCoder.encodeQtCbfZero(pcCU, TLibCommon.TEXT_LUMA, uiTrMode)
                uiNullBitsY := this.m_pcEntropyCoder.getNumberOfWrittenBits()
                nullCostY := this.m_pcRdCost.calcRdCost(uiNullBitsY, uiDistY, false, TLibCommon.DF_DEFAULT)
                if nullCostY < singleCostY {
                    uiAbsSumY = 0
                    for i := int(0); i < uiNumSamplesLuma; i++ {
                        pcCoeffCurrY[i] = 0 //, sizeof( TLibCommon.TCoeff ) * uiNumSamplesLuma );
                    }
                    if checkTransformSkipY {
                        minCostY = nullCostY
                    }
                } else {
                    uiDistY = uiNonzeroDistY
                    if checkTransformSkipY {
                        minCostY = singleCostY
                    }
                }
            }
        } else if checkTransformSkipY {
            this.m_pcEntropyCoder.resetBits()
            this.m_pcEntropyCoder.encodeQtCbfZero(pcCU, TLibCommon.TEXT_LUMA, uiTrMode)
            uiNullBitsY := this.m_pcEntropyCoder.getNumberOfWrittenBits()
            minCostY = this.m_pcRdCost.calcRdCost(uiNullBitsY, uiDistY, false, TLibCommon.DF_DEFAULT)
        }

        if uiAbsSumY == 0 {
            pcPtr := this.m_pcQTTempTComYuv[uiQTTempAccessLayer].GetLumaAddr1(absTUPartIdx)
            uiStride := this.m_pcQTTempTComYuv[uiQTTempAccessLayer].GetStride()
            for uiY := uint(0); uiY < trHeight; uiY++ {
                for uiX := uint(0); uiX < trWidth; uiX++ {
                    pcPtr[uiY*uiStride+uiX] = 0 //, sizeof( TLibCommon.Pel ) * trWidth );
                }

                //pcPtr = pcPtr[uiStride:]
            }
        }

        uiDistU := uint(0)
        uiDistV := uint(0)
        if bCodeChroma {
            uiDistU = this.m_pcRdCost.getDistPart(TLibCommon.G_bitDepthC, this.m_pTempPel, int(trWidthC), pcResi.GetCbAddr1(absTUPartIdxC), int(pcResi.GetCStride()), trWidthC, trHeightC, TLibCommon.TEXT_CHROMA_U, TLibCommon.DF_SSE)
            //#if WEIGHTED_CHROMA_DISTORTION

            //#endif
            // initialized with zero residual destortion
            if puiZeroDist != nil {
                *puiZeroDist += uiDistU
            }
            if uiAbsSumU != 0 {
                pcResiCurrU := this.m_pcQTTempTComYuv[uiQTTempAccessLayer].GetCbAddr1(absTUPartIdxC)

                curChromaQpOffset := pcCU.GetSlice().GetPPS().GetChromaCbQpOffset() + pcCU.GetSlice().GetSliceQpDeltaCb()
                this.m_pcTrQuant.SetQPforQuant(int(pcCU.GetQP1(0)), TLibCommon.TEXT_CHROMA, pcCU.GetSlice().GetSPS().GetQpBDOffsetC(), curChromaQpOffset)

                scalingListType := 3 + TLibCommon.G_eTTable[int(TLibCommon.TEXT_CHROMA_U)]
                //assert(scalingListType < 6);
                this.m_pcTrQuant.InvtransformNxN(pcCU.GetCUTransquantBypass1(uiAbsPartIdx), TLibCommon.TEXT_CHROMA, TLibCommon.REG_DCT, pcResiCurrU, this.m_pcQTTempTComYuv[uiQTTempAccessLayer].GetCStride(), pcCoeffCurrU, trWidthC, trHeightC, scalingListType, false)

                uiNonzeroDistU := this.m_pcRdCost.getDistPart(TLibCommon.G_bitDepthC, this.m_pcQTTempTComYuv[uiQTTempAccessLayer].GetCbAddr1(absTUPartIdxC), int(this.m_pcQTTempTComYuv[uiQTTempAccessLayer].GetCStride()),
                    pcResi.GetCbAddr1(absTUPartIdxC), int(pcResi.GetCStride()), trWidthC, trHeightC, TLibCommon.TEXT_CHROMA_U, TLibCommon.DF_SSE)
                //#if WEIGHTED_CHROMA_DISTORTION

                //#endif

                if pcCU.IsLosslessCoded(0) {
                    uiDistU = uiNonzeroDistU
                } else {
                    dSingleCostU := this.m_pcRdCost.calcRdCost(uiSingleBitsU, uiNonzeroDistU, false, TLibCommon.DF_DEFAULT)
                    this.m_pcEntropyCoder.resetBits()
                    this.m_pcEntropyCoder.encodeQtCbfZero(pcCU, TLibCommon.TEXT_CHROMA_U, uiTrMode)
                    uiNullBitsU := this.m_pcEntropyCoder.getNumberOfWrittenBits()
                    dNullCostU := this.m_pcRdCost.calcRdCost(uiNullBitsU, uiDistU, false, TLibCommon.DF_DEFAULT)
                    if dNullCostU < dSingleCostU {
                        uiAbsSumU = 0
                        for i := int(0); i < uiNumSamplesChro; i++ {
                            pcCoeffCurrU[i] = 0 //, sizeof( TLibCommon.TCoeff ) *  );
                        }
                        if checkTransformSkipUV {
                            minCostU = dNullCostU
                        }
                    } else {
                        uiDistU = uiNonzeroDistU
                        if checkTransformSkipUV {
                            minCostU = dSingleCostU
                        }
                    }
                }
            } else if checkTransformSkipUV {
                this.m_pcEntropyCoder.resetBits()
                this.m_pcEntropyCoder.encodeQtCbfZero(pcCU, TLibCommon.TEXT_CHROMA_U, uiTrModeC)
                uiNullBitsU := this.m_pcEntropyCoder.getNumberOfWrittenBits()
                minCostU = this.m_pcRdCost.calcRdCost(uiNullBitsU, uiDistU, false, TLibCommon.DF_DEFAULT)
            }
            if uiAbsSumU == 0 {
                pcPtr := this.m_pcQTTempTComYuv[uiQTTempAccessLayer].GetCbAddr1(absTUPartIdxC)
                uiStride := this.m_pcQTTempTComYuv[uiQTTempAccessLayer].GetCStride()
                for uiY := uint(0); uiY < trHeightC; uiY++ {
                    for uiX := uint(0); uiX < trWidthC; uiX++ {
                        pcPtr[uiY*uiStride+uiX] = 0 //, sizeof(TLibCommon.Pel) *  );
                    }
                    //pcPtr = pcPtr[uiStride:]
                }
            }

            uiDistV = this.m_pcRdCost.getDistPart(TLibCommon.G_bitDepthC, this.m_pTempPel, int(trWidthC), pcResi.GetCrAddr1(absTUPartIdxC), int(pcResi.GetCStride()), trWidthC, trHeightC, TLibCommon.TEXT_CHROMA_V, TLibCommon.DF_SSE)
            //#if WEIGHTED_CHROMA_DISTORTION

            //#endif
            // initialized with zero residual destortion
            if puiZeroDist != nil {
                *puiZeroDist += uiDistV
            }
            if uiAbsSumV != 0 {
                pcResiCurrV := this.m_pcQTTempTComYuv[uiQTTempAccessLayer].GetCrAddr1(absTUPartIdxC)
                curChromaQpOffset := pcCU.GetSlice().GetPPS().GetChromaCrQpOffset() + pcCU.GetSlice().GetSliceQpDeltaCr()
                this.m_pcTrQuant.SetQPforQuant(int(pcCU.GetQP1(0)), TLibCommon.TEXT_CHROMA, pcCU.GetSlice().GetSPS().GetQpBDOffsetC(), curChromaQpOffset)

                scalingListType := 3 + TLibCommon.G_eTTable[int(TLibCommon.TEXT_CHROMA_V)]
                //assert(scalingListType < 6);
                this.m_pcTrQuant.InvtransformNxN(pcCU.GetCUTransquantBypass1(uiAbsPartIdx), TLibCommon.TEXT_CHROMA, TLibCommon.REG_DCT, pcResiCurrV, this.m_pcQTTempTComYuv[uiQTTempAccessLayer].GetCStride(), pcCoeffCurrV, trWidthC, trHeightC, scalingListType, false)

                uiNonzeroDistV := this.m_pcRdCost.getDistPart(TLibCommon.G_bitDepthC, this.m_pcQTTempTComYuv[uiQTTempAccessLayer].GetCrAddr1(absTUPartIdxC), int(this.m_pcQTTempTComYuv[uiQTTempAccessLayer].GetCStride()),
                    pcResi.GetCrAddr1(absTUPartIdxC), int(pcResi.GetCStride()), trWidthC, trHeightC, TLibCommon.TEXT_CHROMA_V, TLibCommon.DF_SSE)
                //#if WEIGHTED_CHROMA_DISTORTION

                //#endif

                if pcCU.IsLosslessCoded(0) {
                    uiDistV = uiNonzeroDistV
                } else {
                    dSingleCostV := this.m_pcRdCost.calcRdCost(uiSingleBitsV, uiNonzeroDistV, false, TLibCommon.DF_DEFAULT)
                    this.m_pcEntropyCoder.resetBits()
                    this.m_pcEntropyCoder.encodeQtCbfZero(pcCU, TLibCommon.TEXT_CHROMA_V, uiTrMode)
                    uiNullBitsV := this.m_pcEntropyCoder.getNumberOfWrittenBits()
                    dNullCostV := this.m_pcRdCost.calcRdCost(uiNullBitsV, uiDistV, false, TLibCommon.DF_DEFAULT)
                    if dNullCostV < dSingleCostV {
                        uiAbsSumV = 0
                        for i := int(0); i < uiNumSamplesChro; i++ {
                            pcCoeffCurrV[i] = 0 //, sizeof( TLibCommon.TCoeff ) * uiNumSamplesChro );
                        }
                        if checkTransformSkipUV {
                            minCostV = dNullCostV
                        }
                    } else {
                        uiDistV = uiNonzeroDistV
                        if checkTransformSkipUV {
                            minCostV = dSingleCostV
                        }
                    }
                }
            } else if checkTransformSkipUV {
                this.m_pcEntropyCoder.resetBits()
                this.m_pcEntropyCoder.encodeQtCbfZero(pcCU, TLibCommon.TEXT_CHROMA_V, uiTrModeC)
                uiNullBitsV := this.m_pcEntropyCoder.getNumberOfWrittenBits()
                minCostV = this.m_pcRdCost.calcRdCost(uiNullBitsV, uiDistV, false, TLibCommon.DF_DEFAULT)
            }
            if uiAbsSumV == 0 {
                pcPtr := this.m_pcQTTempTComYuv[uiQTTempAccessLayer].GetCrAddr1(absTUPartIdxC)
                uiStride := this.m_pcQTTempTComYuv[uiQTTempAccessLayer].GetCStride()
                for uiY := uint(0); uiY < trHeightC; uiY++ {
                    for uiX := uint(0); uiX < trWidthC; uiX++ {
                        pcPtr[uiY*uiStride+uiX] = 0 //, sizeof(TLibCommon.Pel) * trWidthC );
                    }
                    //pcPtr = pcPtr[uiStride:]
                }
            }
        }
        if uiAbsSumY != 0 {
            pcCU.SetCbfSubParts4(byte(uiSetCbf), TLibCommon.TEXT_LUMA, uiAbsPartIdx, uiDepth)
        } else {
            pcCU.SetCbfSubParts4(0, TLibCommon.TEXT_LUMA, uiAbsPartIdx, uiDepth)
        }
        if bCodeChroma {
            if uiAbsSumU != 0 {
                pcCU.SetCbfSubParts4(byte(uiSetCbf), TLibCommon.TEXT_CHROMA_U, uiAbsPartIdx, uint(pcCU.GetDepth1(0))+uiTrModeC)
            } else {
                pcCU.SetCbfSubParts4(0, TLibCommon.TEXT_CHROMA_U, uiAbsPartIdx, uint(pcCU.GetDepth1(0))+uiTrModeC)
            }
            if uiAbsSumV != 0 {
                pcCU.SetCbfSubParts4(byte(uiSetCbf), TLibCommon.TEXT_CHROMA_V, uiAbsPartIdx, uint(pcCU.GetDepth1(0))+uiTrModeC)
            } else {
                pcCU.SetCbfSubParts4(0, TLibCommon.TEXT_CHROMA_V, uiAbsPartIdx, uint(pcCU.GetDepth1(0))+uiTrModeC)
            }
        }

        if checkTransformSkipY {
            var uiNonzeroDistY, uiAbsSumTransformSkipY uint
            var dSingleCostY float64

            pcResiCurrY := this.m_pcQTTempTComYuv[uiQTTempAccessLayer].GetLumaAddr1(absTUPartIdx)
            resiYStride := this.m_pcQTTempTComYuv[uiQTTempAccessLayer].GetStride()

            var bestCoeffY [32 * 32]TLibCommon.TCoeff
            for i := int(0); i < uiNumSamplesLuma; i++ {
                bestCoeffY[i] = pcCoeffCurrY[i] //, sizeof(TLibCommon.TCoeff) * uiNumSamplesLuma );
            }

            //#if ADAPTIVE_QP_SELECTION
            var bestArlCoeffY [32 * 32]TLibCommon.TCoeff
            for i := int(0); i < uiNumSamplesLuma; i++ {
                bestArlCoeffY[i] = pcArlCoeffCurrY[i] //, sizeof(TLibCommon.TCoeff) * uiNumSamplesLuma );
            }
            //#endif

            var bestResiY [32 * 32]TLibCommon.Pel
            for i := uint(0); i < trHeight; i++ {
                for j := uint(0); j < trWidth; j++ {
                    bestResiY[i*trWidth+j] = pcResiCurrY[i*resiYStride+j] //, sizeof(TLibCommon.Pel) * trWidth );
                }
            }

            if this.m_bUseSBACRD {
                this.m_pcRDGoOnSbacCoder.load(this.m_pppcRDSbacCoder[uiDepth][TLibCommon.CI_QT_TRAFO_ROOT])
            }

            pcCU.SetTransformSkipSubParts4(true, TLibCommon.TEXT_LUMA, uiAbsPartIdx, uiDepth)

            if this.m_pcEncCfg.GetUseRDOQTS() {
                this.m_pcEntropyCoder.estimateBit(this.m_pcTrQuant.GetEstBitsSbac(), int(trWidth), int(trHeight), TLibCommon.TEXT_LUMA)
            }

            this.m_pcTrQuant.SetQPforQuant(int(pcCU.GetQP1(0)), TLibCommon.TEXT_LUMA, pcCU.GetSlice().GetSPS().GetQpBDOffsetY(), 0)

            //#if RDOQ_CHROMA_LAMBDA
            this.m_pcTrQuant.SelectLambda(TLibCommon.TEXT_LUMA)
            //#endif
            this.m_pcTrQuant.TransformNxN(pcCU, pcResi.GetLumaAddr1(absTUPartIdx), pcResi.GetStride(), pcCoeffCurrY,
                //#if ADAPTIVE_QP_SELECTION
                pcArlCoeffCurrY,
                //#endif
                trWidth, trHeight, &uiAbsSumTransformSkipY, TLibCommon.TEXT_LUMA, uiAbsPartIdx, true)
            if uiAbsSumTransformSkipY != 0 {
                pcCU.SetCbfSubParts4(byte(uiSetCbf), TLibCommon.TEXT_LUMA, uiAbsPartIdx, uiDepth)
            } else {
                pcCU.SetCbfSubParts4(0, TLibCommon.TEXT_LUMA, uiAbsPartIdx, uiDepth)
            }

            if uiAbsSumTransformSkipY != 0 {
                this.m_pcEntropyCoder.resetBits()
                this.m_pcEntropyCoder.encodeQtCbf(pcCU, uiAbsPartIdx, TLibCommon.TEXT_LUMA, uiTrMode)
                this.m_pcEntropyCoder.encodeCoeffNxN(pcCU, pcCoeffCurrY, uiAbsPartIdx, trWidth, trHeight, uiDepth, TLibCommon.TEXT_LUMA)
                uiTsSingleBitsY := this.m_pcEntropyCoder.getNumberOfWrittenBits()

                this.m_pcTrQuant.SetQPforQuant(int(pcCU.GetQP1(0)), TLibCommon.TEXT_LUMA, pcCU.GetSlice().GetSPS().GetQpBDOffsetY(), 0)

                scalingListType := 3 + TLibCommon.G_eTTable[int(TLibCommon.TEXT_LUMA)]
                //assert(scalingListType < 6);

                this.m_pcTrQuant.InvtransformNxN(pcCU.GetCUTransquantBypass1(uiAbsPartIdx), TLibCommon.TEXT_LUMA, TLibCommon.REG_DCT, pcResiCurrY, this.m_pcQTTempTComYuv[uiQTTempAccessLayer].GetStride(), pcCoeffCurrY, trWidth, trHeight, scalingListType, true)

                uiNonzeroDistY = this.m_pcRdCost.getDistPart(TLibCommon.G_bitDepthY, this.m_pcQTTempTComYuv[uiQTTempAccessLayer].GetLumaAddr1(absTUPartIdx), int(this.m_pcQTTempTComYuv[uiQTTempAccessLayer].GetStride()),
                    pcResi.GetLumaAddr1(absTUPartIdx), int(pcResi.GetStride()), trWidth, trHeight, TLibCommon.TEXT_LUMA, TLibCommon.DF_SSE)

                dSingleCostY = this.m_pcRdCost.calcRdCost(uiTsSingleBitsY, uiNonzeroDistY, false, TLibCommon.DF_DEFAULT)
            }

            if uiAbsSumTransformSkipY == 0 || minCostY < dSingleCostY {
                pcCU.SetTransformSkipSubParts4(false, TLibCommon.TEXT_LUMA, uiAbsPartIdx, uiDepth)
                for i := int(0); i < uiNumSamplesLuma; i++ {
                    pcCoeffCurrY[i] = bestCoeffY[i] //, sizeof(TLibCommon.TCoeff) * uiNumSamplesLuma );
                    //#if ADAPTIVE_QP_SELECTION
                    pcArlCoeffCurrY[i] = bestArlCoeffY[i] //, sizeof(TLibCommon.TCoeff) * uiNumSamplesLuma );
                    //#endif
                }
                for i := uint(0); i < trHeight; i++ {
                    for j := uint(0); j < trWidth; j++ {
                        pcResiCurrY[i*resiYStride+j] = bestResiY[i*trWidth+j] //, sizeof(TLibCommon.Pel) * trWidth );
                    }
                }
            } else {
                uiDistY = uiNonzeroDistY
                uiAbsSumY = uiAbsSumTransformSkipY
                uiBestTransformMode[0] = 1
            }

            if uiAbsSumY != 0 {
                pcCU.SetCbfSubParts4(byte(uiSetCbf), TLibCommon.TEXT_LUMA, uiAbsPartIdx, uiDepth)
            } else {
                pcCU.SetCbfSubParts4(0, TLibCommon.TEXT_LUMA, uiAbsPartIdx, uiDepth)
            }
        }

        if bCodeChroma && checkTransformSkipUV {
            var uiNonzeroDistU, uiNonzeroDistV, uiAbsSumTransformSkipU, uiAbsSumTransformSkipV uint
            var dSingleCostU, dSingleCostV float64

            pcResiCurrU := this.m_pcQTTempTComYuv[uiQTTempAccessLayer].GetCbAddr1(absTUPartIdxC)
            pcResiCurrV := this.m_pcQTTempTComYuv[uiQTTempAccessLayer].GetCrAddr1(absTUPartIdxC)
            resiCStride := this.m_pcQTTempTComYuv[uiQTTempAccessLayer].GetCStride()

            var bestCoeffU, bestCoeffV [32 * 32]TLibCommon.TCoeff
            var bestArlCoeffU, bestArlCoeffV [32 * 32]TLibCommon.TCoeff
            for i := int(0); i < uiNumSamplesChro; i++ {
                bestCoeffU[i] = pcCoeffCurrU[i] //, sizeof(TLibCommon.TCoeff) * uiNumSamplesChro );
                bestCoeffV[i] = pcCoeffCurrV[i] //, sizeof(TLibCommon.TCoeff) * uiNumSamplesChro );

                //#if ADAPTIVE_QP_SELECTION
                bestArlCoeffU[i] = pcArlCoeffCurrU[i] //, sizeof(TLibCommon.TCoeff) * uiNumSamplesChro );
                bestArlCoeffV[i] = pcArlCoeffCurrV[i] //, sizeof(TLibCommon.TCoeff) * uiNumSamplesChro );
                //#endif
            }

            var bestResiU, bestResiV [32 * 32]TLibCommon.Pel
            for i := uint(0); i < trHeightC; i++ {
                for j := uint(0); j < trWidthC; j++ {
                    bestResiU[i*trWidthC+j] = pcResiCurrU[i*resiCStride+j] //, sizeof(TLibCommon.Pel) * trWidthC );
                    bestResiV[i*trWidthC+j] = pcResiCurrV[i*resiCStride+j] //, sizeof(TLibCommon.Pel) * trWidthC );
                }
            }

            if this.m_bUseSBACRD {
                this.m_pcRDGoOnSbacCoder.load(this.m_pppcRDSbacCoder[uiDepth][TLibCommon.CI_QT_TRAFO_ROOT])
            }

            pcCU.SetTransformSkipSubParts4(true, TLibCommon.TEXT_CHROMA_U, uiAbsPartIdx, uint(pcCU.GetDepth1(0))+uiTrModeC)
            pcCU.SetTransformSkipSubParts4(true, TLibCommon.TEXT_CHROMA_V, uiAbsPartIdx, uint(pcCU.GetDepth1(0))+uiTrModeC)

            if this.m_pcEncCfg.GetUseRDOQTS() {
                this.m_pcEntropyCoder.estimateBit(this.m_pcTrQuant.GetEstBitsSbac(), int(trWidthC), int(trHeightC), TLibCommon.TEXT_CHROMA)
            }

            curChromaQpOffset := pcCU.GetSlice().GetPPS().GetChromaCbQpOffset() + pcCU.GetSlice().GetSliceQpDeltaCb()
            this.m_pcTrQuant.SetQPforQuant(int(pcCU.GetQP1(0)), TLibCommon.TEXT_CHROMA, pcCU.GetSlice().GetSPS().GetQpBDOffsetC(), curChromaQpOffset)

            //#if RDOQ_CHROMA_LAMBDA
            this.m_pcTrQuant.SelectLambda(TLibCommon.TEXT_CHROMA)
            //#endif

            this.m_pcTrQuant.TransformNxN(pcCU, pcResi.GetCbAddr1(absTUPartIdxC), pcResi.GetCStride(), pcCoeffCurrU,
                //#if ADAPTIVE_QP_SELECTION
                pcArlCoeffCurrU,
                //#endif
                trWidthC, trHeightC, &uiAbsSumTransformSkipU, TLibCommon.TEXT_CHROMA_U, uiAbsPartIdx, true)
            curChromaQpOffset = pcCU.GetSlice().GetPPS().GetChromaCrQpOffset() + pcCU.GetSlice().GetSliceQpDeltaCr()
            this.m_pcTrQuant.SetQPforQuant(int(pcCU.GetQP1(0)), TLibCommon.TEXT_CHROMA, pcCU.GetSlice().GetSPS().GetQpBDOffsetC(), curChromaQpOffset)
            this.m_pcTrQuant.TransformNxN(pcCU, pcResi.GetCrAddr1(absTUPartIdxC), pcResi.GetCStride(), pcCoeffCurrV,
                //#if ADAPTIVE_QP_SELECTION
                pcArlCoeffCurrV,
                //#endif
                trWidthC, trHeightC, &uiAbsSumTransformSkipV, TLibCommon.TEXT_CHROMA_V, uiAbsPartIdx, true)

            if uiAbsSumTransformSkipU != 0 {
                pcCU.SetCbfSubParts4(byte(uiSetCbf), TLibCommon.TEXT_CHROMA_U, uiAbsPartIdx, uint(pcCU.GetDepth1(0))+uiTrModeC)
            } else {
                pcCU.SetCbfSubParts4(0, TLibCommon.TEXT_CHROMA_U, uiAbsPartIdx, uint(pcCU.GetDepth1(0))+uiTrModeC)
            }
            if uiAbsSumTransformSkipV != 0 {
                pcCU.SetCbfSubParts4(byte(uiSetCbf), TLibCommon.TEXT_CHROMA_V, uiAbsPartIdx, uint(pcCU.GetDepth1(0))+uiTrModeC)
            } else {
                pcCU.SetCbfSubParts4(0, TLibCommon.TEXT_CHROMA_V, uiAbsPartIdx, uint(pcCU.GetDepth1(0))+uiTrModeC)
            }
            this.m_pcEntropyCoder.resetBits()
            uiSingleBitsU = 0
            uiSingleBitsV = 0

            if uiAbsSumTransformSkipU != 0 {
                this.m_pcEntropyCoder.encodeQtCbf(pcCU, uiAbsPartIdx, TLibCommon.TEXT_CHROMA_U, uiTrMode)
                this.m_pcEntropyCoder.encodeCoeffNxN(pcCU, pcCoeffCurrU, uiAbsPartIdx, trWidthC, trHeightC, uiDepth, TLibCommon.TEXT_CHROMA_U)
                uiSingleBitsU = this.m_pcEntropyCoder.getNumberOfWrittenBits()

                curChromaQpOffset = pcCU.GetSlice().GetPPS().GetChromaCbQpOffset() + pcCU.GetSlice().GetSliceQpDeltaCb()
                this.m_pcTrQuant.SetQPforQuant(int(pcCU.GetQP1(0)), TLibCommon.TEXT_CHROMA, pcCU.GetSlice().GetSPS().GetQpBDOffsetC(), curChromaQpOffset)

                scalingListType := 3 + TLibCommon.G_eTTable[int(TLibCommon.TEXT_CHROMA_U)]
                //assert(scalingListType < 6);

                this.m_pcTrQuant.InvtransformNxN(pcCU.GetCUTransquantBypass1(uiAbsPartIdx), TLibCommon.TEXT_CHROMA, TLibCommon.REG_DCT, pcResiCurrU, this.m_pcQTTempTComYuv[uiQTTempAccessLayer].GetCStride(), pcCoeffCurrU, trWidthC, trHeightC, scalingListType, true)

                uiNonzeroDistU = this.m_pcRdCost.getDistPart(TLibCommon.G_bitDepthC, this.m_pcQTTempTComYuv[uiQTTempAccessLayer].GetCbAddr1(absTUPartIdxC), int(this.m_pcQTTempTComYuv[uiQTTempAccessLayer].GetCStride()),
                    pcResi.GetCbAddr1(absTUPartIdxC), int(pcResi.GetCStride()), trWidthC, trHeightC, TLibCommon.TEXT_CHROMA_U, TLibCommon.DF_SSE)
                //#if WEIGHTED_CHROMA_DISTORTION

                //#endif

                dSingleCostU = this.m_pcRdCost.calcRdCost(uiSingleBitsU, uiNonzeroDistU, false, TLibCommon.DF_DEFAULT)
            }

            if uiAbsSumTransformSkipU == 0 || minCostU < dSingleCostU {
                pcCU.SetTransformSkipSubParts4(false, TLibCommon.TEXT_CHROMA_U, uiAbsPartIdx, uint(pcCU.GetDepth1(0))+uiTrModeC)

                for i := int(0); i < uiNumSamplesChro; i++ {
                    pcCoeffCurrU[i] = bestCoeffU[i] //, sizeof (TLibCommon.TCoeff) * uiNumSamplesChro );
                    //#if ADAPTIVE_QP_SELECTION
                    pcArlCoeffCurrU[i] = bestArlCoeffU[i] //, sizeof (TLibCommon.TCoeff) * uiNumSamplesChro );
                    //#endif
                }
                for i := uint(0); i < trHeightC; i++ {
                    for j := uint(0); j < trWidthC; j++ {
                        pcResiCurrU[i*resiCStride+j] = bestResiU[i*trWidthC+j] //, sizeof(TLibCommon.Pel) * trWidthC );
                    }
                }
            } else {
                uiDistU = uiNonzeroDistU
                uiAbsSumU = uiAbsSumTransformSkipU
                uiBestTransformMode[1] = 1
            }

            if uiAbsSumTransformSkipV != 0 {
                this.m_pcEntropyCoder.encodeQtCbf(pcCU, uiAbsPartIdx, TLibCommon.TEXT_CHROMA_V, uiTrMode)
                this.m_pcEntropyCoder.encodeCoeffNxN(pcCU, pcCoeffCurrV, uiAbsPartIdx, trWidthC, trHeightC, uiDepth, TLibCommon.TEXT_CHROMA_V)
                uiSingleBitsV = this.m_pcEntropyCoder.getNumberOfWrittenBits() - uiSingleBitsU

                curChromaQpOffset = pcCU.GetSlice().GetPPS().GetChromaCrQpOffset() + pcCU.GetSlice().GetSliceQpDeltaCr()
                this.m_pcTrQuant.SetQPforQuant(int(pcCU.GetQP1(0)), TLibCommon.TEXT_CHROMA, pcCU.GetSlice().GetSPS().GetQpBDOffsetC(), curChromaQpOffset)

                scalingListType := 3 + TLibCommon.G_eTTable[int(TLibCommon.TEXT_CHROMA_V)]
                //assert(scalingListType < 6);

                this.m_pcTrQuant.InvtransformNxN(pcCU.GetCUTransquantBypass1(uiAbsPartIdx), TLibCommon.TEXT_CHROMA, TLibCommon.REG_DCT, pcResiCurrV, this.m_pcQTTempTComYuv[uiQTTempAccessLayer].GetCStride(), pcCoeffCurrV, trWidthC, trHeightC, scalingListType, true)

                uiNonzeroDistV = this.m_pcRdCost.getDistPart(TLibCommon.G_bitDepthC, this.m_pcQTTempTComYuv[uiQTTempAccessLayer].GetCrAddr1(absTUPartIdxC), int(this.m_pcQTTempTComYuv[uiQTTempAccessLayer].GetCStride()),
                    pcResi.GetCrAddr1(absTUPartIdxC), int(pcResi.GetCStride()), trWidthC, trHeightC, TLibCommon.TEXT_CHROMA_V, TLibCommon.DF_SSE)
                //#if WEIGHTED_CHROMA_DISTORTION

                //#endif

                dSingleCostV = this.m_pcRdCost.calcRdCost(uiSingleBitsV, uiNonzeroDistV, false, TLibCommon.DF_DEFAULT)
            }

            if uiAbsSumTransformSkipV == 0 || minCostV < dSingleCostV {
                pcCU.SetTransformSkipSubParts4(false, TLibCommon.TEXT_CHROMA_V, uiAbsPartIdx, uint(pcCU.GetDepth1(0))+uiTrModeC)

                for i := int(0); i < uiNumSamplesChro; i++ {
                    pcCoeffCurrV[i] = bestCoeffV[i] //, sizeof(TLibCommon.TCoeff) * uiNumSamplesChro );
                    //#if ADAPTIVE_QP_SELECTION
                    pcArlCoeffCurrV[i] = bestArlCoeffV[i] //, sizeof(TLibCommon.TCoeff) * uiNumSamplesChro );
                    //#endif
                }
                for i := uint(0); i < trHeightC; i++ {
                    for j := uint(0); j < trWidthC; j++ {
                        pcResiCurrV[i*resiCStride+j] = bestResiV[i*trWidthC+j] //, sizeof(TLibCommon.Pel) * trWidthC );
                    }
                }
            } else {
                uiDistV = uiNonzeroDistV
                uiAbsSumV = uiAbsSumTransformSkipV
                uiBestTransformMode[2] = 1
            }

            if uiAbsSumU != 0 {
                pcCU.SetCbfSubParts4(byte(uiSetCbf), TLibCommon.TEXT_CHROMA_U, uiAbsPartIdx, uint(pcCU.GetDepth1(0))+uiTrModeC)
            } else {
                pcCU.SetCbfSubParts4(0, TLibCommon.TEXT_CHROMA_U, uiAbsPartIdx, uint(pcCU.GetDepth1(0))+uiTrModeC)
            }
            if uiAbsSumV != 0 {
                pcCU.SetCbfSubParts4(byte(uiSetCbf), TLibCommon.TEXT_CHROMA_V, uiAbsPartIdx, uint(pcCU.GetDepth1(0))+uiTrModeC)
            } else {
                pcCU.SetCbfSubParts4(0, TLibCommon.TEXT_CHROMA_V, uiAbsPartIdx, uint(pcCU.GetDepth1(0))+uiTrModeC)
            }
        }

        if this.m_bUseSBACRD {
            this.m_pcRDGoOnSbacCoder.load(this.m_pppcRDSbacCoder[uiDepth][TLibCommon.CI_QT_TRAFO_ROOT])
        }

        this.m_pcEntropyCoder.resetBits()

        {
            if uiLog2TrSize > pcCU.GetQuadtreeTULog2MinSizeInCU(uiAbsPartIdx) {
                this.m_pcEntropyCoder.encodeTransformSubdivFlag(0, 5-uiLog2TrSize)
            }
        }

        {
            if bCodeChroma {
                this.m_pcEntropyCoder.encodeQtCbf(pcCU, uiAbsPartIdx, TLibCommon.TEXT_CHROMA_U, uiTrMode)
                this.m_pcEntropyCoder.encodeQtCbf(pcCU, uiAbsPartIdx, TLibCommon.TEXT_CHROMA_V, uiTrMode)
            }

            this.m_pcEntropyCoder.encodeQtCbf(pcCU, uiAbsPartIdx, TLibCommon.TEXT_LUMA, uiTrMode)
        }

        this.m_pcEntropyCoder.encodeCoeffNxN(pcCU, pcCoeffCurrY, uiAbsPartIdx, trWidth, trHeight, uiDepth, TLibCommon.TEXT_LUMA)

        if bCodeChroma {
            this.m_pcEntropyCoder.encodeCoeffNxN(pcCU, pcCoeffCurrU, uiAbsPartIdx, trWidthC, trHeightC, uiDepth, TLibCommon.TEXT_CHROMA_U)
            this.m_pcEntropyCoder.encodeCoeffNxN(pcCU, pcCoeffCurrV, uiAbsPartIdx, trWidthC, trHeightC, uiDepth, TLibCommon.TEXT_CHROMA_V)
        }

        uiSingleBits = this.m_pcEntropyCoder.getNumberOfWrittenBits()

        uiSingleDist = uiDistY + uiDistU + uiDistV
        dSingleCost = this.m_pcRdCost.calcRdCost(uiSingleBits, uiSingleDist, false, TLibCommon.DF_DEFAULT)
    }

    // code sub-blocks
    if bCheckSplit {
        if this.m_bUseSBACRD && bCheckFull {
            this.m_pcRDGoOnSbacCoder.store(this.m_pppcRDSbacCoder[uiDepth][TLibCommon.CI_QT_TRAFO_TEST])
            this.m_pcRDGoOnSbacCoder.load(this.m_pppcRDSbacCoder[uiDepth][TLibCommon.CI_QT_TRAFO_ROOT])
        }
        uiSubdivDist := uint(0)
        uiSubdivBits := uint(0)
        dSubdivCost := float64(0.0)

        uiQPartNumSubdiv := pcCU.GetPic().GetNumPartInCU() >> ((uiDepth + 1) << 1)
        for ui := uint(0); ui < 4; ui++ {
            nsAddr := uiAbsPartIdx + ui*uiQPartNumSubdiv
            if bCheckFull {
                this.xEstimateResidualQT(pcCU, ui, uiAbsPartIdx+ui*uiQPartNumSubdiv, nsAddr, pcResi, uiDepth+1, &dSubdivCost, &uiSubdivBits, &uiSubdivDist, nil)
            } else {
                this.xEstimateResidualQT(pcCU, ui, uiAbsPartIdx+ui*uiQPartNumSubdiv, nsAddr, pcResi, uiDepth+1, &dSubdivCost, &uiSubdivBits, &uiSubdivDist, puiZeroDist)
            }
        }

        uiYCbf := uint(0)
        uiUCbf := uint(0)
        uiVCbf := uint(0)
        for ui := uint(0); ui < 4; ui++ {
            uiYCbf |= uint(pcCU.GetCbf3(uiAbsPartIdx+ui*uiQPartNumSubdiv, TLibCommon.TEXT_LUMA, uiTrMode+1))
            uiUCbf |= uint(pcCU.GetCbf3(uiAbsPartIdx+ui*uiQPartNumSubdiv, TLibCommon.TEXT_CHROMA_U, uiTrMode+1))
            uiVCbf |= uint(pcCU.GetCbf3(uiAbsPartIdx+ui*uiQPartNumSubdiv, TLibCommon.TEXT_CHROMA_V, uiTrMode+1))
        }
        for ui := uint(0); ui < 4*uiQPartNumSubdiv; ui++ {
            pcCU.GetCbf1(TLibCommon.TEXT_LUMA)[uiAbsPartIdx+ui] |= byte(uiYCbf << uiTrMode)
            pcCU.GetCbf1(TLibCommon.TEXT_CHROMA_U)[uiAbsPartIdx+ui] |= byte(uiUCbf << uiTrMode)
            pcCU.GetCbf1(TLibCommon.TEXT_CHROMA_V)[uiAbsPartIdx+ui] |= byte(uiVCbf << uiTrMode)
        }

        if this.m_bUseSBACRD {
            this.m_pcRDGoOnSbacCoder.load(this.m_pppcRDSbacCoder[uiDepth][TLibCommon.CI_QT_TRAFO_ROOT])
        }
        this.m_pcEntropyCoder.resetBits()

        {
            this.xEncodeResidualQT(pcCU, uiAbsPartIdx, uiDepth, true, TLibCommon.TEXT_LUMA)
            this.xEncodeResidualQT(pcCU, uiAbsPartIdx, uiDepth, false, TLibCommon.TEXT_LUMA)
            this.xEncodeResidualQT(pcCU, uiAbsPartIdx, uiDepth, false, TLibCommon.TEXT_CHROMA_U)
            this.xEncodeResidualQT(pcCU, uiAbsPartIdx, uiDepth, false, TLibCommon.TEXT_CHROMA_V)
        }

        uiSubdivBits = this.m_pcEntropyCoder.getNumberOfWrittenBits()
        dSubdivCost = this.m_pcRdCost.calcRdCost(uiSubdivBits, uiSubdivDist, false, TLibCommon.DF_DEFAULT)

        if uiYCbf != 0 || uiUCbf != 0 || uiVCbf != 0 || !bCheckFull {
            if dSubdivCost < dSingleCost {
                *rdCost += dSubdivCost
                *ruiBits += uiSubdivBits
                *ruiDist += uiSubdivDist
                return
            }
        }
        pcCU.SetTransformSkipSubParts4(uiBestTransformMode[0] != 0, TLibCommon.TEXT_LUMA, uiAbsPartIdx, uiDepth)
        if bCodeChroma {
            pcCU.SetTransformSkipSubParts4(uiBestTransformMode[1] != 0, TLibCommon.TEXT_CHROMA_U, uiAbsPartIdx, uint(pcCU.GetDepth1(0))+uiTrModeC)
            pcCU.SetTransformSkipSubParts4(uiBestTransformMode[2] != 0, TLibCommon.TEXT_CHROMA_V, uiAbsPartIdx, uint(pcCU.GetDepth1(0))+uiTrModeC)
        }
        //assert( bCheckFull );
        if this.m_bUseSBACRD {
            this.m_pcRDGoOnSbacCoder.load(this.m_pppcRDSbacCoder[uiDepth][TLibCommon.CI_QT_TRAFO_TEST])
        }
    }
    *rdCost += dSingleCost
    *ruiBits += uiSingleBits
    *ruiDist += uiSingleDist

    pcCU.SetTrIdxSubParts(uiTrMode, uiAbsPartIdx, uiDepth)

    if uiAbsSumY != 0 {
        pcCU.SetCbfSubParts4(byte(uiSetCbf), TLibCommon.TEXT_LUMA, uiAbsPartIdx, uiDepth)
    } else {
        pcCU.SetCbfSubParts4(0, TLibCommon.TEXT_LUMA, uiAbsPartIdx, uiDepth)
    }
    if bCodeChroma {
        if uiAbsSumU != 0 {
            pcCU.SetCbfSubParts4(byte(uiSetCbf), TLibCommon.TEXT_CHROMA_U, uiAbsPartIdx, uint(pcCU.GetDepth1(0))+uiTrModeC)
        } else {
            pcCU.SetCbfSubParts4(0, TLibCommon.TEXT_CHROMA_U, uiAbsPartIdx, uint(pcCU.GetDepth1(0))+uiTrModeC)
        }
        if uiAbsSumV != 0 {
            pcCU.SetCbfSubParts4(byte(uiSetCbf), TLibCommon.TEXT_CHROMA_V, uiAbsPartIdx, uint(pcCU.GetDepth1(0))+uiTrModeC)
        } else {
            pcCU.SetCbfSubParts4(0, TLibCommon.TEXT_CHROMA_V, uiAbsPartIdx, uint(pcCU.GetDepth1(0))+uiTrModeC)
        }
    }
}

func (this *TEncSearch) xSetResidualQTData(pcCU *TLibCommon.TComDataCU, uiQuadrant, uiAbsPartIdx, absTUPartIdx uint, pcResi *TLibCommon.TComYuv, uiDepth uint, bSpatial bool) {
    //assert( pcCU.GetDepth1( 0 ) == pcCU.GetDepth1( uiAbsPartIdx ) );
    uiCurrTrMode := uiDepth - uint(pcCU.GetDepth1(0))
    uiTrMode := uint(pcCU.GetTransformIdx1(uiAbsPartIdx))

    if uiCurrTrMode == uiTrMode {
        uiLog2TrSize := uint(TLibCommon.G_aucConvertToBit[pcCU.GetSlice().GetSPS().GetMaxCUWidth()>>uiDepth] + 2)
        uiQTTempAccessLayer := pcCU.GetSlice().GetSPS().GetQuadtreeTULog2MaxSize() - uiLog2TrSize

        bCodeChroma := true
        uiTrModeC := uiTrMode
        uiLog2TrSizeC := uiLog2TrSize - 1
        if uiLog2TrSize == 2 {
            uiLog2TrSizeC++
            uiTrModeC--
            uiQPDiv := pcCU.GetPic().GetNumPartInCU() >> ((uint(pcCU.GetDepth1(0)) + uiTrModeC) << 1)
            bCodeChroma = ((uiAbsPartIdx % uiQPDiv) == 0)
        }

        if bSpatial {
            trWidth := 1 << uiLog2TrSize
            trHeight := 1 << uiLog2TrSize
            this.m_pcQTTempTComYuv[uiQTTempAccessLayer].CopyPartToPartLuma(pcResi, absTUPartIdx, uint(trWidth), uint(trHeight))

            if bCodeChroma {

                this.m_pcQTTempTComYuv[uiQTTempAccessLayer].CopyPartToPartChroma(pcResi, uiAbsPartIdx, 1<<uiLog2TrSizeC, 1<<uiLog2TrSizeC)

            }
        } else {
            uiNumCoeffPerAbsPartIdxIncrement := pcCU.GetSlice().GetSPS().GetMaxCUWidth() * pcCU.GetSlice().GetSPS().GetMaxCUHeight() >> (pcCU.GetSlice().GetSPS().GetMaxCUDepth() << 1)
            uiNumCoeffY := (1 << (uiLog2TrSize << 1))
            pcCoeffSrcY := this.m_ppcQTTempCoeffY[uiQTTempAccessLayer][uiNumCoeffPerAbsPartIdxIncrement*uiAbsPartIdx:]
            pcCoeffDstY := pcCU.GetCoeffY()[uiNumCoeffPerAbsPartIdxIncrement*uiAbsPartIdx:]
            for i := int(0); i < uiNumCoeffY; i++ {
                pcCoeffDstY[i] = pcCoeffSrcY[i] //, sizeof( TLibCommon.TCoeff ) * uiNumCoeffY );
            }
            //#if ADAPTIVE_QP_SELECTION
            pcArlCoeffSrcY := this.m_ppcQTTempArlCoeffY[uiQTTempAccessLayer][uiNumCoeffPerAbsPartIdxIncrement*uiAbsPartIdx:]
            pcArlCoeffDstY := pcCU.GetArlCoeffY()[uiNumCoeffPerAbsPartIdxIncrement*uiAbsPartIdx:]
            for i := int(0); i < uiNumCoeffY; i++ {
                pcArlCoeffDstY[i] = pcArlCoeffSrcY[i] //, sizeof( int ) * uiNumCoeffY );
            }
            //#endif
            if bCodeChroma {
                uiNumCoeffC := (1 << (uiLog2TrSizeC << 1))
                pcCoeffSrcU := this.m_ppcQTTempCoeffCb[uiQTTempAccessLayer][(uiNumCoeffPerAbsPartIdxIncrement * uiAbsPartIdx >> 2):]
                pcCoeffSrcV := this.m_ppcQTTempCoeffCr[uiQTTempAccessLayer][(uiNumCoeffPerAbsPartIdxIncrement * uiAbsPartIdx >> 2):]
                pcCoeffDstU := pcCU.GetCoeffCb()[(uiNumCoeffPerAbsPartIdxIncrement * uiAbsPartIdx >> 2):]
                pcCoeffDstV := pcCU.GetCoeffCr()[(uiNumCoeffPerAbsPartIdxIncrement * uiAbsPartIdx >> 2):]
                for i := int(0); i < uiNumCoeffC; i++ {
                    pcCoeffDstU[i] = pcCoeffSrcU[i] //, sizeof( TLibCommon.TCoeff ) * uiNumCoeffC );
                    pcCoeffDstV[i] = pcCoeffSrcV[i] //, sizeof( TLibCommon.TCoeff ) * uiNumCoeffC );
                }
                //#if ADAPTIVE_QP_SELECTION
                pcArlCoeffSrcU := this.m_ppcQTTempArlCoeffCb[uiQTTempAccessLayer][(uiNumCoeffPerAbsPartIdxIncrement * uiAbsPartIdx >> 2):]
                pcArlCoeffSrcV := this.m_ppcQTTempArlCoeffCr[uiQTTempAccessLayer][(uiNumCoeffPerAbsPartIdxIncrement * uiAbsPartIdx >> 2):]
                pcArlCoeffDstU := pcCU.GetArlCoeffCb()[(uiNumCoeffPerAbsPartIdxIncrement * uiAbsPartIdx >> 2):]
                pcArlCoeffDstV := pcCU.GetArlCoeffCr()[(uiNumCoeffPerAbsPartIdxIncrement * uiAbsPartIdx >> 2):]
                for i := int(0); i < uiNumCoeffC; i++ {
                    pcArlCoeffDstU[i] = pcArlCoeffSrcU[i] //, sizeof( int ) * uiNumCoeffC );
                    pcArlCoeffDstV[i] = pcArlCoeffSrcV[i] //, sizeof( int ) * uiNumCoeffC );
                }
                //#endif
            }
        }
    } else {
        uiQPartNumSubdiv := pcCU.GetPic().GetNumPartInCU() >> ((uiDepth + 1) << 1)
        for ui := uint(0); ui < 4; ui++ {
            nsAddr := uiAbsPartIdx + ui*uiQPartNumSubdiv
            this.xSetResidualQTData(pcCU, ui, uiAbsPartIdx+ui*uiQPartNumSubdiv, nsAddr, pcResi, uiDepth+1, bSpatial)
        }
    }
}

func (this *TEncSearch) xModeBitsIntra(pcCU *TLibCommon.TComDataCU, uiMode, uiPU, uiPartOffset, uiDepth, uiInitTrDepth uint) uint {
    if this.m_bUseSBACRD {
        // Reload only contexts required for coding intra mode information
        this.m_pcRDGoOnSbacCoder.loadIntraDirModeLuma(this.m_pppcRDSbacCoder[uiDepth][TLibCommon.CI_CURR_BEST])
    }

    pcCU.SetLumaIntraDirSubParts(uiMode, uiPartOffset, uiDepth+uiInitTrDepth)

    this.m_pcEntropyCoder.resetBits()
    this.m_pcEntropyCoder.encodeIntraDirModeLuma(pcCU, uiPartOffset, false)

    return this.m_pcEntropyCoder.getNumberOfWrittenBits()
}

func (this *TEncSearch) xUpdateCandList(uiMode uint, uiCost float64, uiFastCandNum uint, CandModeList []uint, CandCostList []float64) uint {
    var i uint
    shift := uint(0)

    for shift < uiFastCandNum && uiCost < CandCostList[uiFastCandNum-1-shift] {
        shift++
    }

    if shift != 0 {
        for i = 1; i < shift; i++ {
            CandModeList[uiFastCandNum-i] = CandModeList[uiFastCandNum-1-i]
            CandCostList[uiFastCandNum-i] = CandCostList[uiFastCandNum-1-i]
        }
        CandModeList[uiFastCandNum-shift] = uiMode
        CandCostList[uiFastCandNum-shift] = uiCost
        return 1
    }

    return 0
}

// -------------------------------------------------------------------------------------------------------------------
// compute symbol bits
// -------------------------------------------------------------------------------------------------------------------

func (this *TEncSearch) xAddSymbolBitsInter(pcCU *TLibCommon.TComDataCU,
    uiQp uint,
    uiTrMode uint,
    ruiBits *uint,
    rpcYuvRec *TLibCommon.TComYuv,
    pcYuvPred *TLibCommon.TComYuv,
    rpcYuvResi *TLibCommon.TComYuv) {
    if pcCU.GetMergeFlag1(0) && pcCU.GetPartitionSize1(0) == TLibCommon.SIZE_2Nx2N && !pcCU.GetQtRootCbf(0) {
        pcCU.SetSkipFlagSubParts(true, 0, uint(pcCU.GetDepth1(0)))

        this.m_pcEntropyCoder.resetBits()
        if pcCU.GetSlice().GetPPS().GetTransquantBypassEnableFlag() {
            this.m_pcEntropyCoder.encodeCUTransquantBypassFlag(pcCU, 0, true)
        }
        this.m_pcEntropyCoder.encodeSkipFlag(pcCU, 0, true)
        this.m_pcEntropyCoder.encodeMergeIndex(pcCU, 0, true)
        *ruiBits += this.m_pcEntropyCoder.getNumberOfWrittenBits()
    } else {
        this.m_pcEntropyCoder.resetBits()
        if pcCU.GetSlice().GetPPS().GetTransquantBypassEnableFlag() {
            this.m_pcEntropyCoder.encodeCUTransquantBypassFlag(pcCU, 0, true)
        }
        this.m_pcEntropyCoder.encodeSkipFlag(pcCU, 0, true)
        this.m_pcEntropyCoder.encodePredMode(pcCU, 0, true)
        this.m_pcEntropyCoder.encodePartSize(pcCU, 0, uint(pcCU.GetDepth1(0)), true)
        this.m_pcEntropyCoder.encodePredInfo(pcCU, 0, true)
        bDummy := false
        this.m_pcEntropyCoder.encodeCoeff(pcCU, 0, uint(pcCU.GetDepth1(0)), uint(uint(pcCU.GetWidth1(0))), uint(pcCU.GetHeight1(0)), &bDummy)

        *ruiBits += this.m_pcEntropyCoder.getNumberOfWrittenBits()
    }
}

func (this *TEncSearch) setWpScalingDistParam(pcCU *TLibCommon.TComDataCU, iRefIdx int, eRefPicListCur TLibCommon.RefPicList) {
    if iRefIdx < 0 {
        this.m_cDistParam.bApplyWeight = false
        return
    }

    pcSlice := pcCU.GetSlice()
    pps := pcCU.GetSlice().GetPPS()
    var wp0, wp1 []TLibCommon.WpScalingParam
    this.m_cDistParam.bApplyWeight = (pcSlice.GetSliceType() == TLibCommon.P_SLICE && pps.GetUseWP()) || (pcSlice.GetSliceType() == TLibCommon.B_SLICE && pps.GetWPBiPred())
    if !this.m_cDistParam.bApplyWeight {
        return
    }
    var iRefIdx0, iRefIdx1 int

    if eRefPicListCur == TLibCommon.REF_PIC_LIST_0 {
        iRefIdx0 = iRefIdx
    } else {
        iRefIdx0 = (-1)
    }
    if eRefPicListCur == TLibCommon.REF_PIC_LIST_1 {
        iRefIdx1 = iRefIdx
    } else {
        iRefIdx1 = (-1)
    }

    wp0, wp1 = this.GetWpScaling(pcCU, iRefIdx0, iRefIdx1)

    if iRefIdx0 < 0 {
        wp0 = nil
    }
    if iRefIdx1 < 0 {
        wp1 = nil
    }

    this.m_cDistParam.wpCur = nil

    if eRefPicListCur == TLibCommon.REF_PIC_LIST_0 {
        this.m_cDistParam.wpCur = wp0
    } else {
        this.m_cDistParam.wpCur = wp1
    }
}

func (this *TEncSearch) setDistParamComp(uiComp uint) { this.m_cDistParam.uiComp = uiComp }
