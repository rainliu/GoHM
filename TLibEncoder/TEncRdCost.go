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
    "gohm/TLibCommon"
    "math"
)

// ====================================================================================================================
// Type definition
// ====================================================================================================================
// for function pointer
type FpDistFunc func(*DistParam) uint

// ====================================================================================================================
// Class definition
// ====================================================================================================================

/// distortion parameter class
type DistParam struct {
    pOrg       []TLibCommon.Pel
    pCur       []TLibCommon.Pel
    iStrideOrg int
    iStrideCur int
    iRows      int
    iCols      int
    iStep      int
    DistFunc   FpDistFunc
    bitDepth   int

    bApplyWeight bool                        // whether weithed prediction is used or not
    wpCur        []TLibCommon.WpScalingParam // weithed prediction scaling parameters for current ref
    uiComp       uint                        // uiComp = 0 (luma Y), 1 (chroma U), 2 (chroma V)

    //#if NS_HAD
    //              bUseNSHAD	bool;
    //#endif

    // (vertical) subsampling shift (for reducing complexity)
    // - 0 = no subsampling, 1 = even rows, 2 = every 4th, etc.
    iSubShift int
}

func NewDistParam() *DistParam {
    return &DistParam{pOrg: nil,
        pCur:       nil,
        iStrideOrg: 0,
        iStrideCur: 0,
        iRows:      0,
        iCols:      0,
        iStep:      1,
        //DistFunc = NULL;
        iSubShift: 0,
        bitDepth:  0}
    //#if NS_HAD
    //    bUseNSHAD : false
    //#endif
}

/// RD cost computation class, with Weighted Prediction
type TEncRdCostWeightPrediction struct {
    m_w0       int
    m_w1       int // current wp scaling values
    m_shift    uint
    m_offset   int
    m_round    int
    m_xSetDone bool
}

func NewTEncRdCostWeightPrediction() *TEncRdCostWeightPrediction {
    return &TEncRdCostWeightPrediction{}
}

/*func (this *TEncRdCostWeightPrediction) DISTORTION_PRECISION_ADJUSTMENT(x int) uint {
    return uint(x)
}*/

func (this *TEncRdCostWeightPrediction) xSetWPscale(w0, w1, shift, offset, round int) {
    this.m_w0 = w0
    this.m_w1 = w1
    this.m_shift = uint(shift)
    this.m_offset = offset
    this.m_round = round

    this.m_xSetDone = true
}

func (this *TEncRdCostWeightPrediction) xGetSSEw(pcDtParam *DistParam) uint {
    piOrg := pcDtParam.pOrg
    piCur := pcDtParam.pCur
    var pred TLibCommon.Pel
    iRows := pcDtParam.iRows
    iCols := pcDtParam.iCols
    iStrideOrg := pcDtParam.iStrideOrg
    iStrideCur := pcDtParam.iStrideCur

    //assert( pcDtParam.iSubShift == 0 );

    uiComp := pcDtParam.uiComp
    //assert(uiComp<3);
    wpCur := pcDtParam.wpCur[uiComp]
    w0 := wpCur.W
    offset := wpCur.Offset
    shift := uint(wpCur.Shift)
    round := wpCur.Round

    uiSum := uint(0) 
    uiShift := TLibCommon.DISTORTION_PRECISION_ADJUSTMENT(uint(pcDtParam.bitDepth - 8) << 1).(uint)

    var iTemp int

    for ; iRows != 0; iRows-- {
        for n := 0; n < iCols; n++ {
            pred = TLibCommon.Pel(((w0*int(piCur[n]) + round) >> shift) + offset)

            iTemp = int(piOrg[n] - pred)
            uiSum += uint(iTemp*iTemp) >> uiShift
        }
        piOrg = piOrg[iStrideOrg:]
        piCur = piCur[iStrideCur:]
    }

    pcDtParam.uiComp = 255 // reset for DEBUG (assert test)

    return (uiSum)
}
func (this *TEncRdCostWeightPrediction) xGetSADw(pcDtParam *DistParam) uint {
    var pred TLibCommon.Pel
    piOrg := pcDtParam.pOrg
    piCur := pcDtParam.pCur
    iRows := pcDtParam.iRows
    iCols := pcDtParam.iCols
    iStrideCur := pcDtParam.iStrideCur
    iStrideOrg := pcDtParam.iStrideOrg

    uiComp := pcDtParam.uiComp
    //assert(uiComp<3);
    wpCur := (pcDtParam.wpCur[uiComp])
    w0 := wpCur.W
    offset := wpCur.Offset
    shift := uint(wpCur.Shift)
    round := wpCur.Round

    uiSum := uint(0)

    for ; iRows != 0; iRows-- {
        for n := 0; n < iCols; n++ {
            pred = TLibCommon.Pel(((w0*int(piCur[n]) + round) >> shift) + offset)

            uiSum += uint(TLibCommon.ABS(int(piOrg[n] - pred)).(int))
        }
        piOrg = piOrg[iStrideOrg:]
        piCur = piCur[iStrideCur:]
    }

    pcDtParam.uiComp = 255 // reset for DEBUG (assert test)

    return uiSum >> TLibCommon.DISTORTION_PRECISION_ADJUSTMENT(uint(pcDtParam.bitDepth-8)).(uint)
}
func (this *TEncRdCostWeightPrediction) xGetHADs4w(pcDtParam *DistParam) uint {
    piOrg := pcDtParam.pOrg
    piCur := pcDtParam.pCur
    iRows := pcDtParam.iRows
    iStrideCur := pcDtParam.iStrideCur
    iStrideOrg := pcDtParam.iStrideOrg
    iStep := pcDtParam.iStep
    var y int
    iOffsetOrg := iStrideOrg << 2
    iOffsetCur := iStrideCur << 2

    uiSum := uint(0)

    for y = 0; y < iRows; y += 4 {
        uiSum += this.xCalcHADs4x4w(piOrg, piCur, iStrideOrg, iStrideCur, iStep)
        piOrg = piOrg[iOffsetOrg:]
        piCur = piCur[iOffsetCur:]
    }

    return uiSum >> TLibCommon.DISTORTION_PRECISION_ADJUSTMENT(uint(pcDtParam.bitDepth-8)).(uint)
}
func (this *TEncRdCostWeightPrediction) xGetHADs8w(pcDtParam *DistParam) uint {
    piOrg := pcDtParam.pOrg
    piCur := pcDtParam.pCur
    iRows := pcDtParam.iRows
    iStrideCur := pcDtParam.iStrideCur
    iStrideOrg := pcDtParam.iStrideOrg
    iStep := pcDtParam.iStep
    var y int

    uiSum := uint(0)

    if iRows == 4 {
        uiSum += this.xCalcHADs4x4w(piOrg[0:], piCur, iStrideOrg, iStrideCur, iStep)
        uiSum += this.xCalcHADs4x4w(piOrg[4:], piCur[4*iStep:], iStrideOrg, iStrideCur, iStep)
    } else {
        iOffsetOrg := iStrideOrg << 3
        iOffsetCur := iStrideCur << 3
        for y = 0; y < iRows; y += 8 {
            uiSum += this.xCalcHADs8x8w(piOrg, piCur, iStrideOrg, iStrideCur, iStep)
            piOrg = piOrg[iOffsetOrg:]
            piCur = piCur[iOffsetCur:]
        }
    }

    return uiSum >> TLibCommon.DISTORTION_PRECISION_ADJUSTMENT(uint(pcDtParam.bitDepth-8)).(uint)
}
func (this *TEncRdCostWeightPrediction) xGetHADsw(pcDtParam *DistParam) uint {
    piOrg := pcDtParam.pOrg
    piCur := pcDtParam.pCur
    iRows := pcDtParam.iRows
    iCols := pcDtParam.iCols
    iStrideCur := pcDtParam.iStrideCur
    iStrideOrg := pcDtParam.iStrideOrg
    iStep := pcDtParam.iStep

    var x, y int

    uiComp := pcDtParam.uiComp
    //assert(uiComp<3);
    wpCur := (pcDtParam.wpCur[uiComp])

    this.xSetWPscale(wpCur.W, 0, wpCur.Shift, wpCur.Offset, wpCur.Round)

    uiSum := uint(0)

    if (iRows%8 == 0) && (iCols%8 == 0) {
        iOffsetOrg := iStrideOrg << 3
        iOffsetCur := iStrideCur << 3
        for y = 0; y < iRows; y += 8 {
            for x = 0; x < iCols; x += 8 {
                uiSum += this.xCalcHADs8x8w(piOrg[x:], piCur[x*iStep:], iStrideOrg, iStrideCur, iStep)
            }
            piOrg = piOrg[iOffsetOrg:]
            piCur = piCur[iOffsetCur:]
        }
    } else if (iRows%4 == 0) && (iCols%4 == 0) {
        iOffsetOrg := iStrideOrg << 2
        iOffsetCur := iStrideCur << 2

        for y = 0; y < iRows; y += 4 {
            for x = 0; x < iCols; x += 4 {
                uiSum += this.xCalcHADs4x4w(piOrg[x:], piCur[x*iStep:], iStrideOrg, iStrideCur, iStep)
            }
            piOrg = piOrg[iOffsetOrg:]
            piCur = piCur[iOffsetCur:]
        }
    } else {
        for y = 0; y < iRows; y += 2 {
            for x = 0; x < iCols; x += 2 {
                uiSum += this.xCalcHADs2x2w(piOrg[x:], piCur[x*iStep:], iStrideOrg, iStrideCur, iStep)
            }
            piOrg = piOrg[iStrideOrg:]
            piCur = piCur[iStrideCur:]
        }
    }

    this.m_xSetDone = false

    return uiSum >> TLibCommon.DISTORTION_PRECISION_ADJUSTMENT(uint(pcDtParam.bitDepth-8)).(uint)
}
func (this *TEncRdCostWeightPrediction) xCalcHADs2x2w(piOrg []TLibCommon.Pel, piCur []TLibCommon.Pel, iStrideOrg, iStrideCur, iStep int) uint {
    var satd int
    var diff, m [4]int

    //assert( m_xSetDone );
    var pred TLibCommon.Pel

    pred = TLibCommon.Pel(((this.m_w0*int(piCur[0*iStep]) + this.m_round) >> this.m_shift) + this.m_offset)
    diff[0] = int(piOrg[0] - pred)
    pred = TLibCommon.Pel(((this.m_w0*int(piCur[1*iStep]) + this.m_round) >> this.m_shift) + this.m_offset)
    diff[1] = int(piOrg[1] - pred)
    pred = TLibCommon.Pel(((this.m_w0*int(piCur[0*iStep+iStrideCur]) + this.m_round) >> this.m_shift) + this.m_offset)
    diff[2] = int(piOrg[iStrideOrg] - pred)
    pred = TLibCommon.Pel(((this.m_w0*int(piCur[1*iStep+iStrideCur]) + this.m_round) >> this.m_shift) + this.m_offset)
    diff[3] = int(piOrg[iStrideOrg+1] - pred)

    m[0] = diff[0] + diff[2]
    m[1] = diff[1] + diff[3]
    m[2] = diff[0] - diff[2]
    m[3] = diff[1] - diff[3]

    satd += TLibCommon.ABS(m[0] + m[1]).(int)
    satd += TLibCommon.ABS(m[0] - m[1]).(int)
    satd += TLibCommon.ABS(m[2] + m[3]).(int)
    satd += TLibCommon.ABS(m[2] - m[3]).(int)

    return uint(satd)
}

func (this *TEncRdCostWeightPrediction) xCalcHADs4x4w(piOrg []TLibCommon.Pel, piCur []TLibCommon.Pel, iStrideOrg, iStrideCur, iStep int) uint {
    var k, satd int
    var diff, m, d [16]int

    //assert( this.m_xSetDone );
    var pred TLibCommon.Pel

    for k = 0; k < 16; k += 4 {
        pred = TLibCommon.Pel(((this.m_w0*int(piCur[0*iStep]) + this.m_round) >> this.m_shift) + this.m_offset)
        diff[k+0] = int(piOrg[0] - pred)
        pred = TLibCommon.Pel(((this.m_w0*int(piCur[1*iStep]) + this.m_round) >> this.m_shift) + this.m_offset)
        diff[k+1] = int(piOrg[1] - pred)
        pred = TLibCommon.Pel(((this.m_w0*int(piCur[2*iStep]) + this.m_round) >> this.m_shift) + this.m_offset)
        diff[k+2] = int(piOrg[2] - pred)
        pred = TLibCommon.Pel(((this.m_w0*int(piCur[3*iStep]) + this.m_round) >> this.m_shift) + this.m_offset)
        diff[k+3] = int(piOrg[3] - pred)

        piCur = piCur[iStrideCur:]
        piOrg = piOrg[iStrideOrg:]
    }

    /*===== hadamard transform =====*/
    m[0] = diff[0] + diff[12]
    m[1] = diff[1] + diff[13]
    m[2] = diff[2] + diff[14]
    m[3] = diff[3] + diff[15]
    m[4] = diff[4] + diff[8]
    m[5] = diff[5] + diff[9]
    m[6] = diff[6] + diff[10]
    m[7] = diff[7] + diff[11]
    m[8] = diff[4] - diff[8]
    m[9] = diff[5] - diff[9]
    m[10] = diff[6] - diff[10]
    m[11] = diff[7] - diff[11]
    m[12] = diff[0] - diff[12]
    m[13] = diff[1] - diff[13]
    m[14] = diff[2] - diff[14]
    m[15] = diff[3] - diff[15]

    d[0] = m[0] + m[4]
    d[1] = m[1] + m[5]
    d[2] = m[2] + m[6]
    d[3] = m[3] + m[7]
    d[4] = m[8] + m[12]
    d[5] = m[9] + m[13]
    d[6] = m[10] + m[14]
    d[7] = m[11] + m[15]
    d[8] = m[0] - m[4]
    d[9] = m[1] - m[5]
    d[10] = m[2] - m[6]
    d[11] = m[3] - m[7]
    d[12] = m[12] - m[8]
    d[13] = m[13] - m[9]
    d[14] = m[14] - m[10]
    d[15] = m[15] - m[11]

    m[0] = d[0] + d[3]
    m[1] = d[1] + d[2]
    m[2] = d[1] - d[2]
    m[3] = d[0] - d[3]
    m[4] = d[4] + d[7]
    m[5] = d[5] + d[6]
    m[6] = d[5] - d[6]
    m[7] = d[4] - d[7]
    m[8] = d[8] + d[11]
    m[9] = d[9] + d[10]
    m[10] = d[9] - d[10]
    m[11] = d[8] - d[11]
    m[12] = d[12] + d[15]
    m[13] = d[13] + d[14]
    m[14] = d[13] - d[14]
    m[15] = d[12] - d[15]

    d[0] = m[0] + m[1]
    d[1] = m[0] - m[1]
    d[2] = m[2] + m[3]
    d[3] = m[3] - m[2]
    d[4] = m[4] + m[5]
    d[5] = m[4] - m[5]
    d[6] = m[6] + m[7]
    d[7] = m[7] - m[6]
    d[8] = m[8] + m[9]
    d[9] = m[8] - m[9]
    d[10] = m[10] + m[11]
    d[11] = m[11] - m[10]
    d[12] = m[12] + m[13]
    d[13] = m[12] - m[13]
    d[14] = m[14] + m[15]
    d[15] = m[15] - m[14]

    for k = 0; k < 16; k++ {
        satd += TLibCommon.ABS(d[k]).(int)
    }
    satd = ((satd + 1) >> 1)

    return uint(satd)
}

func (this *TEncRdCostWeightPrediction) xCalcHADs8x8w(piOrg []TLibCommon.Pel, piCur []TLibCommon.Pel, iStrideOrg, iStrideCur, iStep int) uint {
    var k, i, j, jj, sad int
    var diff [64]int
    var m1, m2, m3 [8][8]int
    iStep2 := iStep << 1
    iStep3 := iStep2 + iStep
    iStep4 := iStep3 + iStep
    iStep5 := iStep4 + iStep
    iStep6 := iStep5 + iStep
    iStep7 := iStep6 + iStep

    //assert( m_xSetDone );
    var pred TLibCommon.Pel

    for k = 0; k < 64; k += 8 {
        pred = TLibCommon.Pel(((this.m_w0*int(piCur[0]) + this.m_round) >> this.m_shift) + this.m_offset)
        diff[k+0] = int(piOrg[0] - pred)
        pred = TLibCommon.Pel(((this.m_w0*int(piCur[iStep]) + this.m_round) >> this.m_shift) + this.m_offset)
        diff[k+1] = int(piOrg[1] - pred)
        pred = TLibCommon.Pel(((this.m_w0*int(piCur[iStep2]) + this.m_round) >> this.m_shift) + this.m_offset)
        diff[k+2] = int(piOrg[2] - pred)
        pred = TLibCommon.Pel(((this.m_w0*int(piCur[iStep3]) + this.m_round) >> this.m_shift) + this.m_offset)
        diff[k+3] = int(piOrg[3] - pred)
        pred = TLibCommon.Pel(((this.m_w0*int(piCur[iStep4]) + this.m_round) >> this.m_shift) + this.m_offset)
        diff[k+4] = int(piOrg[4] - pred)
        pred = TLibCommon.Pel(((this.m_w0*int(piCur[iStep5]) + this.m_round) >> this.m_shift) + this.m_offset)
        diff[k+5] = int(piOrg[5] - pred)
        pred = TLibCommon.Pel(((this.m_w0*int(piCur[iStep6]) + this.m_round) >> this.m_shift) + this.m_offset)
        diff[k+6] = int(piOrg[6] - pred)
        pred = TLibCommon.Pel(((this.m_w0*int(piCur[iStep7]) + this.m_round) >> this.m_shift) + this.m_offset)
        diff[k+7] = int(piOrg[7] - pred)

        piCur = piCur[iStrideCur:]
        piOrg = piOrg[iStrideOrg:]
    }

    //horizontal
    for j = 0; j < 8; j++ {
        jj = j << 3
        m2[j][0] = diff[jj] + diff[jj+4]
        m2[j][1] = diff[jj+1] + diff[jj+5]
        m2[j][2] = diff[jj+2] + diff[jj+6]
        m2[j][3] = diff[jj+3] + diff[jj+7]
        m2[j][4] = diff[jj] - diff[jj+4]
        m2[j][5] = diff[jj+1] - diff[jj+5]
        m2[j][6] = diff[jj+2] - diff[jj+6]
        m2[j][7] = diff[jj+3] - diff[jj+7]

        m1[j][0] = m2[j][0] + m2[j][2]
        m1[j][1] = m2[j][1] + m2[j][3]
        m1[j][2] = m2[j][0] - m2[j][2]
        m1[j][3] = m2[j][1] - m2[j][3]
        m1[j][4] = m2[j][4] + m2[j][6]
        m1[j][5] = m2[j][5] + m2[j][7]
        m1[j][6] = m2[j][4] - m2[j][6]
        m1[j][7] = m2[j][5] - m2[j][7]

        m2[j][0] = m1[j][0] + m1[j][1]
        m2[j][1] = m1[j][0] - m1[j][1]
        m2[j][2] = m1[j][2] + m1[j][3]
        m2[j][3] = m1[j][2] - m1[j][3]
        m2[j][4] = m1[j][4] + m1[j][5]
        m2[j][5] = m1[j][4] - m1[j][5]
        m2[j][6] = m1[j][6] + m1[j][7]
        m2[j][7] = m1[j][6] - m1[j][7]
    }

    //vertical
    for i = 0; i < 8; i++ {
        m3[0][i] = m2[0][i] + m2[4][i]
        m3[1][i] = m2[1][i] + m2[5][i]
        m3[2][i] = m2[2][i] + m2[6][i]
        m3[3][i] = m2[3][i] + m2[7][i]
        m3[4][i] = m2[0][i] - m2[4][i]
        m3[5][i] = m2[1][i] - m2[5][i]
        m3[6][i] = m2[2][i] - m2[6][i]
        m3[7][i] = m2[3][i] - m2[7][i]

        m1[0][i] = m3[0][i] + m3[2][i]
        m1[1][i] = m3[1][i] + m3[3][i]
        m1[2][i] = m3[0][i] - m3[2][i]
        m1[3][i] = m3[1][i] - m3[3][i]
        m1[4][i] = m3[4][i] + m3[6][i]
        m1[5][i] = m3[5][i] + m3[7][i]
        m1[6][i] = m3[4][i] - m3[6][i]
        m1[7][i] = m3[5][i] - m3[7][i]

        m2[0][i] = m1[0][i] + m1[1][i]
        m2[1][i] = m1[0][i] - m1[1][i]
        m2[2][i] = m1[2][i] + m1[3][i]
        m2[3][i] = m1[2][i] - m1[3][i]
        m2[4][i] = m1[4][i] + m1[5][i]
        m2[5][i] = m1[4][i] - m1[5][i]
        m2[6][i] = m1[6][i] + m1[7][i]
        m2[7][i] = m1[6][i] - m1[7][i]
    }

    for j = 0; j < 8; j++ {
        for i = 0; i < 8; i++ {
            sad += TLibCommon.ABS(m2[j][i]).(int)
        }
    }

    sad = ((sad + 2) >> 2)

    return uint(sad)
}

/// RD cost computation class
type TEncRdCost struct {
    TEncRdCostWeightPrediction
    // for distortion
    m_iBlkWidth  int
    m_iBlkHeight int

    //#if AMP_SAD
    m_afpDistortFunc [64]FpDistFunc // [eDFunc]
    //#else
    //  FpDistFunc              this.m_afpDistortFunc[33]; // [eDFunc]
    //#endif

    //#if WEIGHTED_CHROMA_DISTORTION
    m_chromaDistortionWeight float64
    //#endif
    m_dLambda           float64
    m_sqrtLambda        float64
    m_uiLambdaMotionSAD uint
    m_uiLambdaMotionSSE uint
    m_dFrameLambda      float64

    // for motion cost
    //#if FIX203
    m_mvPredictor TLibCommon.TComMv
    /*#else
      UInt*                   this.m_puiComponentCostOriginP;
      UInt*                   this.m_puiComponentCost;
      UInt*                   this.m_puiVerCost;
      UInt*                   this.m_puiHorCost;
    #endif*/
    m_uiCost     uint
    m_iCostScale int
    /*#if !FIX203
      Int                     this.m_iSearchLimit;
    #endif*/
}

func NewTEncRdCost() *TEncRdCost {
    return &TEncRdCost{}
}

func (this *TEncRdCost) calcRdCost(uiBits, uiDistortion uint, bFlag bool, eDFunc TLibCommon.DFunc) float64 {
    dRdCost := float64(0.0)
    dLambda := float64(0.0)

    switch eDFunc {
    case TLibCommon.DF_SSE:
        //assert(0);
    case TLibCommon.DF_SAD:
        dLambda = float64(this.m_uiLambdaMotionSAD)

    case TLibCommon.DF_DEFAULT:
        dLambda = this.m_dLambda

    case TLibCommon.DF_SSE_FRAME:
        dLambda = this.m_dFrameLambda

    default:
        //assert (0);

    }

    if bFlag {
        // Intra8x8, Intra4x4 Block only...
        //#if SEQUENCE_LEVEL_LOSSLESS
        //    dRdCost = (Double)(uiBits);
        //#else
        dRdCost = ((float64(uiDistortion)) + (float64(uiBits) * dLambda))
        //#endif
    } else {
        if eDFunc == TLibCommon.DF_SAD {
            dRdCost = (float64(uiDistortion) + float64(int(float64(uiBits)*dLambda+.5)>>16))
            dRdCost = float64(uint(math.Floor(dRdCost)))
        } else {
            //#if SEQUENCE_LEVEL_LOSSLESS
            //      dRdCost = (Double)(uiBits);
            //#else
            dRdCost = (float64(uiDistortion) + float64(int(float64(uiBits)*dLambda+.5)))
            dRdCost = float64(uint(math.Floor(dRdCost)))
            //#endif
        }
    }

    return dRdCost
}
func (this *TEncRdCost) calcRdCost64(uiBits, uiDistortion uint64, bFlag bool, eDFunc TLibCommon.DFunc) float64 {
    dRdCost := float64(0.0)
    dLambda := float64(0.0)

    switch eDFunc {
    case TLibCommon.DF_SSE:
        //      assert(0);
    case TLibCommon.DF_SAD:
        dLambda = float64(this.m_uiLambdaMotionSAD)
    case TLibCommon.DF_DEFAULT:
        dLambda = this.m_dLambda
    case TLibCommon.DF_SSE_FRAME:
        dLambda = this.m_dFrameLambda
    default:
        //      assert (0);
    }

    if bFlag {
        // Intra8x8, Intra4x4 Block only...
        //#if SEQUENCE_LEVEL_LOSSLESS
        //    dRdCost = (Double)(uiBits);
        //#else
        dRdCost = (float64(int64(uiDistortion)) + (float64(uiBits) * dLambda))
        //#endif
    } else {
        if eDFunc == TLibCommon.DF_SAD {
            dRdCost = (float64(int64(uiDistortion)) + float64(int(float64(uiBits)*dLambda+.5)>>16))
            dRdCost = float64(uint(math.Floor(dRdCost)))
        } else {
            //#if SEQUENCE_LEVEL_LOSSLESS
            //      dRdCost = (Double)(uiBits);
            //#else
            dRdCost = (float64(int64(uiDistortion)) + float64(int(float64(uiBits)*dLambda+.5)))
            dRdCost = float64(uint(math.Floor(dRdCost)))
            //#endif
        }
    }

    return dRdCost
}

//#if WEIGHTED_CHROMA_DISTORTION
func (this *TEncRdCost) setChromaDistortionWeight(chromaDistortionWeight float64) {
    this.m_chromaDistortionWeight = chromaDistortionWeight
}

//#endif
func (this *TEncRdCost) setLambda(dLambda float64) {
    this.m_dLambda = dLambda
    this.m_sqrtLambda = math.Sqrt(this.m_dLambda)
    this.m_uiLambdaMotionSAD = uint(math.Floor(65536.0 * this.m_sqrtLambda))
    this.m_uiLambdaMotionSSE = uint(math.Floor(65536.0 * this.m_dLambda))
}
func (this *TEncRdCost) setFrameLambda(dLambda float64) { this.m_dFrameLambda = dLambda }

func (this *TEncRdCost) getSqrtLambda() float64 { return this.m_sqrtLambda }

//#if RATE_CONTROL_LAMBDA_DOMAIN
func (this *TEncRdCost) getLambda() float64 { return this.m_dLambda }

//#endif

// Distortion Functions
func (this *TEncRdCost) init() {
    this.m_afpDistortFunc[0] = nil // for TLibCommon.DF_DEFAULT

    this.m_afpDistortFunc[1] = func(d *DistParam) uint { return this.xGetSSE(d) }
    this.m_afpDistortFunc[2] = func(d *DistParam) uint { return this.xGetSSE4(d) }
    this.m_afpDistortFunc[3] = func(d *DistParam) uint { return this.xGetSSE8(d) }
    this.m_afpDistortFunc[4] = func(d *DistParam) uint { return this.xGetSSE16(d) }
    this.m_afpDistortFunc[5] = func(d *DistParam) uint { return this.xGetSSE32(d) }
    this.m_afpDistortFunc[6] = func(d *DistParam) uint { return this.xGetSSE64(d) }
    this.m_afpDistortFunc[7] = func(d *DistParam) uint { return this.xGetSSE16N(d) }

    this.m_afpDistortFunc[8] = func(d *DistParam) uint { return this.xGetSAD(d) }
    this.m_afpDistortFunc[9] = func(d *DistParam) uint { return this.xGetSAD4(d) }
    this.m_afpDistortFunc[10] = func(d *DistParam) uint { return this.xGetSAD8(d) }
    this.m_afpDistortFunc[11] = func(d *DistParam) uint { return this.xGetSAD16(d) }
    this.m_afpDistortFunc[12] = func(d *DistParam) uint { return this.xGetSAD32(d) }
    this.m_afpDistortFunc[13] = func(d *DistParam) uint { return this.xGetSAD64(d) }
    this.m_afpDistortFunc[14] = func(d *DistParam) uint { return this.xGetSAD16N(d) }

    this.m_afpDistortFunc[15] = func(d *DistParam) uint { return this.xGetSAD(d) }
    this.m_afpDistortFunc[16] = func(d *DistParam) uint { return this.xGetSAD4(d) }
    this.m_afpDistortFunc[17] = func(d *DistParam) uint { return this.xGetSAD8(d) }
    this.m_afpDistortFunc[18] = func(d *DistParam) uint { return this.xGetSAD16(d) }
    this.m_afpDistortFunc[19] = func(d *DistParam) uint { return this.xGetSAD32(d) }
    this.m_afpDistortFunc[20] = func(d *DistParam) uint { return this.xGetSAD64(d) }
    this.m_afpDistortFunc[21] = func(d *DistParam) uint { return this.xGetSAD16N(d) }

    //#if AMP_SAD
    this.m_afpDistortFunc[43] = func(d *DistParam) uint { return this.xGetSAD12(d) }
    this.m_afpDistortFunc[44] = func(d *DistParam) uint { return this.xGetSAD24(d) }
    this.m_afpDistortFunc[45] = func(d *DistParam) uint { return this.xGetSAD48(d) }

    this.m_afpDistortFunc[46] = func(d *DistParam) uint { return this.xGetSAD12(d) }
    this.m_afpDistortFunc[47] = func(d *DistParam) uint { return this.xGetSAD24(d) }
    this.m_afpDistortFunc[48] = func(d *DistParam) uint { return this.xGetSAD48(d) }
    //#endif
    this.m_afpDistortFunc[22] = func(d *DistParam) uint { return this.xGetHADs(d) }
    this.m_afpDistortFunc[23] = func(d *DistParam) uint { return this.xGetHADs(d) }
    this.m_afpDistortFunc[24] = func(d *DistParam) uint { return this.xGetHADs(d) }
    this.m_afpDistortFunc[25] = func(d *DistParam) uint { return this.xGetHADs(d) }
    this.m_afpDistortFunc[26] = func(d *DistParam) uint { return this.xGetHADs(d) }
    this.m_afpDistortFunc[27] = func(d *DistParam) uint { return this.xGetHADs(d) }
    this.m_afpDistortFunc[28] = func(d *DistParam) uint { return this.xGetHADs(d) }

    /*#if !FIX203
      this.m_puiComponentCostOriginP = NULL;
      this.m_puiComponentCost        = NULL;
      this.m_puiVerCost              = NULL;
      this.m_puiHorCost              = NULL;
    #endif*/
    this.m_uiCost = 0
    this.m_iCostScale = 0
    /*#if !FIX203
      this.m_iSearchLimit            = 0xdeaddead;
    #endif*/
}

func (this *TEncRdCost) setDistParam1(uiBlkWidth, uiBlkHeight uint, eDFunc TLibCommon.DFunc, rcDistParam *DistParam) {
    // set Block Width / Height
    rcDistParam.iCols = int(uiBlkWidth)
    rcDistParam.iRows = int(uiBlkHeight)
    rcDistParam.DistFunc = this.m_afpDistortFunc[int(eDFunc)+int(TLibCommon.G_aucConvertToBit[rcDistParam.iCols])+1]

    // initialize
    rcDistParam.iSubShift = 0
}

func (this *TEncRdCost) setDistParam2(pcPatternKey *TLibCommon.TComPattern, piRefY []TLibCommon.Pel, iRefStride int, rcDistParam *DistParam) {
    // set Original & Curr Pointer / Stride
    rcDistParam.pOrg = pcPatternKey.GetROIY()
    rcDistParam.pCur = piRefY

    rcDistParam.iStrideOrg = pcPatternKey.GetPatternLStride()
    rcDistParam.iStrideCur = iRefStride

    // set Block Width / Height
    rcDistParam.iCols = pcPatternKey.GetROIYWidth()
    rcDistParam.iRows = pcPatternKey.GetROIYHeight()
    rcDistParam.DistFunc = this.m_afpDistortFunc[int(TLibCommon.DF_SAD)+int(TLibCommon.G_aucConvertToBit[rcDistParam.iCols])+1]

    //#if AMP_SAD
    if rcDistParam.iCols == 12 {
        rcDistParam.DistFunc = this.m_afpDistortFunc[43]
    } else if rcDistParam.iCols == 24 {
        rcDistParam.DistFunc = this.m_afpDistortFunc[44]
    } else if rcDistParam.iCols == 48 {
        rcDistParam.DistFunc = this.m_afpDistortFunc[45]
    }
    //#endif

    // initialize
    rcDistParam.iSubShift = 0
}

//#if NS_HAD
func (this *TEncRdCost) setDistParam3(pcPatternKey *TLibCommon.TComPattern, piRefY []TLibCommon.Pel, iRefStride, iStep int, rcDistParam *DistParam, bHADME bool) {
    // set Original & Curr Pointer / Stride
    rcDistParam.pOrg = pcPatternKey.GetROIY()
    rcDistParam.pCur = piRefY

    rcDistParam.iStrideOrg = pcPatternKey.GetPatternLStride()
    rcDistParam.iStrideCur = iRefStride * iStep

    // set Step for interpolated buffer
    rcDistParam.iStep = iStep

    // set Block Width / Height
    rcDistParam.iCols = pcPatternKey.GetROIYWidth()
    rcDistParam.iRows = pcPatternKey.GetROIYHeight()
    //#if NS_HAD
    //  rcDistParam.bUseNSHAD = bUseNSHAD;
    //#endif

    // set distortion function
    if !bHADME {
        rcDistParam.DistFunc = this.m_afpDistortFunc[int(TLibCommon.DF_SADS)+int(TLibCommon.G_aucConvertToBit[rcDistParam.iCols])+1]
        //#if AMP_SAD
        if rcDistParam.iCols == 12 {
            rcDistParam.DistFunc = this.m_afpDistortFunc[46]
        } else if rcDistParam.iCols == 24 {
            rcDistParam.DistFunc = this.m_afpDistortFunc[47]
        } else if rcDistParam.iCols == 48 {
            rcDistParam.DistFunc = this.m_afpDistortFunc[48]
        }
        //#endif
    } else {
        rcDistParam.DistFunc = this.m_afpDistortFunc[int(TLibCommon.DF_HADS)+int(TLibCommon.G_aucConvertToBit[rcDistParam.iCols])+1]
    }

    // initialize
    rcDistParam.iSubShift = 0
}

func (this *TEncRdCost) setDistParam4(rcDP *DistParam, bitDepth int, p1 []TLibCommon.Pel, iStride1 int, p2 []TLibCommon.Pel, iStride2, iWidth, iHeight int, bHadamard bool) {
    rcDP.pOrg = p1
    rcDP.pCur = p2
    rcDP.iStrideOrg = iStride1
    rcDP.iStrideCur = iStride2
    rcDP.iCols = iWidth
    rcDP.iRows = iHeight
    rcDP.iStep = 1
    rcDP.iSubShift = 0
    rcDP.bitDepth = bitDepth
    if bHadamard {
        rcDP.DistFunc = this.m_afpDistortFunc[int(TLibCommon.DF_HADS)+int(TLibCommon.G_aucConvertToBit[iWidth])+1]
    } else {
        rcDP.DistFunc = this.m_afpDistortFunc[int(TLibCommon.DF_SADS)+int(TLibCommon.G_aucConvertToBit[iWidth])+1]
    }
    //#if NS_HAD
    //  rcDP.bUseNSHAD  = bUseNSHAD;
    //#endif
}

//#else
//  Void    setDistParam( TComPattern* pcPatternKey, TLibCommon.Pel* piRefY, Int iRefStride, Int iStep, DistParam& rcDistParam, Bool bHADME=false );
//  Void    setDistParam( DistParam& rcDP, Int bitDepth, TLibCommon.Pel* p1, Int iStride1, TLibCommon.Pel* p2, Int iStride2, Int iWidth, Int iHeight, Bool bHadamard = false );
//#endif

func (this *TEncRdCost) calcHAD(bitDepth int, pi0 []TLibCommon.Pel, iStride0 int, pi1 []TLibCommon.Pel, iStride1, iWidth, iHeight int) uint {
    uiSum := uint(0)
    var x, y int

    if ((iWidth % 8) == 0) && ((iHeight % 8) == 0) {
        for y = 0; y < iHeight; y += 8 {
            for x = 0; x < iWidth; x += 8 {
                uiSum += this.xCalcHADs8x8(pi0[x:], pi1[x:], iStride0, iStride1, 1)
            }
            pi0 = pi0[iStride0*8:]
            pi1 = pi1[iStride1*8:]
        }
    } else if ((iWidth % 4) == 0) && ((iHeight % 4) == 0) {
        for y = 0; y < iHeight; y += 4 {
            for x = 0; x < iWidth; x += 4 {
                uiSum += this.xCalcHADs4x4(pi0[x:], pi1[x:], iStride0, iStride1, 1)
            }
            pi0 = pi0[iStride0*4:]
            pi1 = pi1[iStride1*4:]
        }
    } else {
        for y = 0; y < iHeight; y += 2 {
            for x = 0; x < iWidth; x += 2 {
                uiSum += this.xCalcHADs8x8(pi0[x:], pi1[x:], iStride0, iStride1, 1)
            }
            pi0 = pi0[iStride0*2:]
            pi1 = pi1[iStride1*2:]
        }
    }

    return uiSum >> TLibCommon.DISTORTION_PRECISION_ADJUSTMENT(uint(bitDepth-8)).(uint)
}

// for motion cost
//#if !FIX203
//  Void    initRateDistortionModel( Int iSubTLibCommon.PelSearchLimit );
//  Void    xUninit();
//#endif
func (this *TEncRdCost) xGetComponentBits(iVal int) uint {
    uiLength := uint(1)
    var uiTemp uint
    if iVal <= 0 {
        uiTemp = uint(-iVal<<1) + 1
    } else {
        uiTemp = uint(iVal << 1)
    }

    //assert ( uiTemp );

    for 1 != uiTemp {
        uiTemp >>= 1
        uiLength += 2
    }

    return uiLength
}

func (this *TEncRdCost) getMotionCost(bSad bool, iAdd int) {
    if bSad {
        this.m_uiCost = this.m_uiLambdaMotionSAD + uint(iAdd)
    } else {
        this.m_uiCost = this.m_uiLambdaMotionSSE + uint(iAdd)
    }
}
func (this *TEncRdCost) setPredictor(rcMv *TLibCommon.TComMv) {
    //#if FIX203
    this.m_mvPredictor = *rcMv
    //#else
    //    this.m_puiHorCost = this.m_puiComponentCost - rcMv.getHor();
    //    this.m_puiVerCost = this.m_puiComponentCost - rcMv.getVer();
    //#endif
}
func (this *TEncRdCost) setCostScale(iCostScale int) { this.m_iCostScale = iCostScale }
func (this *TEncRdCost) getCost2(x, y int) uint {
    //#if FIX203
    return this.m_uiCost * this.getBits(x, y) >> 16
    //#else
    //    return (( this.m_uiCost * (this.m_puiHorCost[ x * (1<<this.m_iCostScale) ] + this.m_puiVerCost[ y * (1<<this.m_iCostScale) ]) ) >> 16);
    //#endif
}
func (this *TEncRdCost) getCost1(b uint) uint { return (this.m_uiCost * b) >> 16 }
func (this *TEncRdCost) getBits(x, y int) uint {
    //#if FIX203
    return this.xGetComponentBits((x<<uint(this.m_iCostScale))-int(this.m_mvPredictor.GetHor())) +
        this.xGetComponentBits((y<<uint(this.m_iCostScale))-int(this.m_mvPredictor.GetVer()))
    //#else
    //    return this.m_puiHorCost[ x * (1<<this.m_iCostScale)] + this.m_puiVerCost[ y * (1<<this.m_iCostScale) ];
    //#endif
}

//#if WEIGHTED_CHROMA_DISTORTION
func (this *TEncRdCost) getDistPart(bitDepth int, piCur []TLibCommon.Pel, iCurStride int, piOrg []TLibCommon.Pel, iOrgStride int, uiBlkWidth, uiBlkHeight uint, bWeighted bool, eDFunc TLibCommon.DFunc) uint {
    var cDtParam DistParam
    this.setDistParam1(uiBlkWidth, uiBlkHeight, eDFunc, &cDtParam)
    cDtParam.pOrg = piOrg
    cDtParam.pCur = piCur
    cDtParam.iStrideOrg = iOrgStride
    cDtParam.iStrideCur = iCurStride
    cDtParam.iStep = 1

    cDtParam.bApplyWeight = false
    cDtParam.uiComp = 255 // just for assert: to be sure it was set before use, since only values 0,1 or 2 are allowed.
    cDtParam.bitDepth = bitDepth

    //#if WEIGHTED_CHROMA_DISTORTION
    if bWeighted {
        return uint(this.m_chromaDistortionWeight * float64(cDtParam.DistFunc(&cDtParam)))
    }
    return cDtParam.DistFunc(&cDtParam)

    //#else
    //  return cDtParam.DistFunc( &cDtParam );
    //#endif
}

//#else
//  UInt   getDistPart(Int bitDepth, TLibCommon.Pel* piCur, Int iCurStride,  TLibCommon.Pel* piOrg, Int iOrgStride, UInt uiBlkWidth, UInt uiBlkHeight, TLibCommon.DFunc eDFunc = TLibCommon.DF_SSE );
//#endif

//#if RATE_CONTROL_LAMBDA_DOMAIN
func (this *TEncRdCost) getSADPart(bitDepth int, pelCur []TLibCommon.Pel, curStride int, pelOrg []TLibCommon.Pel, orgStride int, width, height int) uint {
    SAD := uint(0)
    shift := TLibCommon.DISTORTION_PRECISION_ADJUSTMENT(uint(bitDepth - 8)).(uint)
    for i := 0; i < height; i++ {
        for j := 0; j < width; j++ {
            SAD += uint(TLibCommon.ABS(int(pelCur[j]-pelOrg[j])).(int)) >> shift
        }
        pelCur = pelCur[curStride:]
        pelOrg = pelOrg[orgStride:]
    }
    return SAD
}

//#endif

func (this *TEncRdCost) xGetSSE(pcDtParam *DistParam) uint {
    if pcDtParam.bApplyWeight {
        return this.xGetSSEw(pcDtParam)
    }
    piOrg := pcDtParam.pOrg
    piCur := pcDtParam.pCur
    iRows := pcDtParam.iRows
    iCols := pcDtParam.iCols
    iStrideOrg := pcDtParam.iStrideOrg
    iStrideCur := pcDtParam.iStrideCur

    uiSum := uint(0)
    uiShift := TLibCommon.DISTORTION_PRECISION_ADJUSTMENT(uint(pcDtParam.bitDepth - 8) << 1).(uint)

    var iTemp int

    for ; iRows != 0; iRows-- {
        for n := 0; n < iCols; n++ {
            iTemp = int(piOrg[n] - piCur[n])
            uiSum += uint(iTemp*iTemp) >> uiShift
        }
        piOrg = piOrg[iStrideOrg:]
        piCur = piCur[iStrideCur:]
    }

    return uiSum
}

func (this *TEncRdCost) xGetSSE4(pcDtParam *DistParam) uint {
    if pcDtParam.bApplyWeight {
        //assert( pcDtParam.iCols == 4 );
        return this.xGetSSEw(pcDtParam)
    }
    piOrg := pcDtParam.pOrg
    piCur := pcDtParam.pCur
    iRows := pcDtParam.iRows
    iStrideOrg := pcDtParam.iStrideOrg
    iStrideCur := pcDtParam.iStrideCur

    uiSum := uint(0)
    uiShift := TLibCommon.DISTORTION_PRECISION_ADJUSTMENT(uint(pcDtParam.bitDepth - 8) << 1).(uint)

    var iTemp int

    for ; iRows != 0; iRows-- {
        iTemp = int(piOrg[0] - piCur[0])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[1] - piCur[1])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[2] - piCur[2])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[3] - piCur[3])
        uiSum += uint(iTemp*iTemp) >> uiShift

        piOrg = piOrg[iStrideOrg:]
        piCur = piCur[iStrideCur:]
    }

    return (uiSum)
}

func (this *TEncRdCost) xGetSSE8(pcDtParam *DistParam) uint {
    if pcDtParam.bApplyWeight {
        //assert( pcDtParam.iCols == 8 );
        return this.xGetSSEw(pcDtParam)
    }
    piOrg := pcDtParam.pOrg
    piCur := pcDtParam.pCur
    iRows := pcDtParam.iRows
    iStrideOrg := pcDtParam.iStrideOrg
    iStrideCur := pcDtParam.iStrideCur

    uiSum := uint(0)
    uiShift := TLibCommon.DISTORTION_PRECISION_ADJUSTMENT(uint(pcDtParam.bitDepth - 8) << 1).(uint) 

    var iTemp int

    for ; iRows != 0; iRows-- {
        iTemp = int(piOrg[0] - piCur[0])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[1] - piCur[1])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[2] - piCur[2])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[3] - piCur[3])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[4] - piCur[4])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[5] - piCur[5])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[6] - piCur[6])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[7] - piCur[7])
        uiSum += uint(iTemp*iTemp) >> uiShift

        piOrg = piOrg[iStrideOrg:]
        piCur = piCur[iStrideCur:]
    }

    return (uiSum)
}
func (this *TEncRdCost) xGetSSE16(pcDtParam *DistParam) uint {
    if pcDtParam.bApplyWeight {
        //assert( pcDtParam.iCols == 16 );
        return this.xGetSSEw(pcDtParam)
    }
    piOrg := pcDtParam.pOrg
    piCur := pcDtParam.pCur
    iRows := pcDtParam.iRows
    iStrideOrg := pcDtParam.iStrideOrg
    iStrideCur := pcDtParam.iStrideCur

    uiSum := uint(0)
    uiShift := TLibCommon.DISTORTION_PRECISION_ADJUSTMENT(uint(pcDtParam.bitDepth - 8) << 1).(uint)

    var iTemp int

    for ; iRows != 0; iRows-- {
        iTemp = int(piOrg[0] - piCur[0])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[1] - piCur[1])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[2] - piCur[2])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[3] - piCur[3])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[4] - piCur[4])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[5] - piCur[5])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[6] - piCur[6])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[7] - piCur[7])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[8] - piCur[8])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[9] - piCur[9])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[10] - piCur[10])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[11] - piCur[11])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[12] - piCur[12])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[13] - piCur[13])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[14] - piCur[14])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[15] - piCur[15])
        uiSum += uint(iTemp*iTemp) >> uiShift

        piOrg = piOrg[iStrideOrg:]
        piCur = piCur[iStrideCur:]
    }

    return (uiSum)
}

func (this *TEncRdCost) xGetSSE32(pcDtParam *DistParam) uint {
    if pcDtParam.bApplyWeight {
        //assert( pcDtParam.iCols == 32 );
        return this.xGetSSEw(pcDtParam)
    }
    piOrg := pcDtParam.pOrg
    piCur := pcDtParam.pCur
    iRows := pcDtParam.iRows
    iStrideOrg := pcDtParam.iStrideOrg
    iStrideCur := pcDtParam.iStrideCur

    uiSum := uint(0)
    uiShift := TLibCommon.DISTORTION_PRECISION_ADJUSTMENT(uint(pcDtParam.bitDepth - 8) << 1).(uint)
    var iTemp int

    for ; iRows != 0; iRows-- {
        iTemp = int(piOrg[0] - piCur[0])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[1] - piCur[1])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[2] - piCur[2])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[3] - piCur[3])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[4] - piCur[4])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[5] - piCur[5])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[6] - piCur[6])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[7] - piCur[7])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[8] - piCur[8])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[9] - piCur[9])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[10] - piCur[10])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[11] - piCur[11])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[12] - piCur[12])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[13] - piCur[13])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[14] - piCur[14])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[15] - piCur[15])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[16] - piCur[16])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[17] - piCur[17])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[18] - piCur[18])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[19] - piCur[19])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[20] - piCur[20])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[21] - piCur[21])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[22] - piCur[22])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[23] - piCur[23])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[24] - piCur[24])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[25] - piCur[25])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[26] - piCur[26])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[27] - piCur[27])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[28] - piCur[28])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[29] - piCur[29])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[30] - piCur[30])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[31] - piCur[31])
        uiSum += uint(iTemp*iTemp) >> uiShift

        piOrg = piOrg[iStrideOrg:]
        piCur = piCur[iStrideCur:]
    }

    return (uiSum)
}

func (this *TEncRdCost) xGetSSE64(pcDtParam *DistParam) uint {
    if pcDtParam.bApplyWeight {
        //assert( pcDtParam.iCols == 64 );
        return this.xGetSSEw(pcDtParam)
    }
    piOrg := pcDtParam.pOrg
    piCur := pcDtParam.pCur
    iRows := pcDtParam.iRows
    iStrideOrg := pcDtParam.iStrideOrg
    iStrideCur := pcDtParam.iStrideCur

    uiSum := uint(0)
    uiShift := TLibCommon.DISTORTION_PRECISION_ADJUSTMENT(uint(pcDtParam.bitDepth - 8) << 1).(uint)
    var iTemp int

    for ; iRows != 0; iRows-- {
        iTemp = int(piOrg[0] - piCur[0])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[1] - piCur[1])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[2] - piCur[2])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[3] - piCur[3])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[4] - piCur[4])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[5] - piCur[5])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[6] - piCur[6])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[7] - piCur[7])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[8] - piCur[8])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[9] - piCur[9])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[10] - piCur[10])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[11] - piCur[11])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[12] - piCur[12])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[13] - piCur[13])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[14] - piCur[14])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[15] - piCur[15])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[16] - piCur[16])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[17] - piCur[17])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[18] - piCur[18])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[19] - piCur[19])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[20] - piCur[20])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[21] - piCur[21])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[22] - piCur[22])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[23] - piCur[23])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[24] - piCur[24])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[25] - piCur[25])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[26] - piCur[26])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[27] - piCur[27])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[28] - piCur[28])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[29] - piCur[29])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[30] - piCur[30])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[31] - piCur[31])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[32] - piCur[32])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[33] - piCur[33])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[34] - piCur[34])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[35] - piCur[35])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[36] - piCur[36])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[37] - piCur[37])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[38] - piCur[38])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[39] - piCur[39])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[40] - piCur[40])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[41] - piCur[41])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[42] - piCur[42])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[43] - piCur[43])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[44] - piCur[44])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[45] - piCur[45])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[46] - piCur[46])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[47] - piCur[47])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[48] - piCur[48])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[49] - piCur[49])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[50] - piCur[50])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[51] - piCur[51])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[52] - piCur[52])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[53] - piCur[53])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[54] - piCur[54])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[55] - piCur[55])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[56] - piCur[56])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[57] - piCur[57])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[58] - piCur[58])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[59] - piCur[59])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[60] - piCur[60])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[61] - piCur[61])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[62] - piCur[62])
        uiSum += uint(iTemp*iTemp) >> uiShift
        iTemp = int(piOrg[63] - piCur[63])
        uiSum += uint(iTemp*iTemp) >> uiShift

        piOrg = piOrg[iStrideOrg:]
        piCur = piCur[iStrideCur:]
    }

    return (uiSum)
}

func (this *TEncRdCost) xGetSSE16N(pcDtParam *DistParam) uint {
    if pcDtParam.bApplyWeight {
        return this.xGetSSEw(pcDtParam)
    }
    piOrg := pcDtParam.pOrg
    piCur := pcDtParam.pCur
    iRows := pcDtParam.iRows
    iCols := pcDtParam.iCols
    iStrideOrg := pcDtParam.iStrideOrg
    iStrideCur := pcDtParam.iStrideCur

    uiSum := uint(0)
    uiShift := TLibCommon.DISTORTION_PRECISION_ADJUSTMENT(uint(pcDtParam.bitDepth - 8) << 1).(uint)
    var iTemp int

    for ; iRows != 0; iRows-- {
        for n := 0; n < iCols; n += 16 {
            iTemp = int(piOrg[n+0] - piCur[n+0])
            uiSum += uint(iTemp*iTemp) >> uiShift
            iTemp = int(piOrg[n+1] - piCur[n+1])
            uiSum += uint(iTemp*iTemp) >> uiShift
            iTemp = int(piOrg[n+2] - piCur[n+2])
            uiSum += uint(iTemp*iTemp) >> uiShift
            iTemp = int(piOrg[n+3] - piCur[n+3])
            uiSum += uint(iTemp*iTemp) >> uiShift
            iTemp = int(piOrg[n+4] - piCur[n+4])
            uiSum += uint(iTemp*iTemp) >> uiShift
            iTemp = int(piOrg[n+5] - piCur[n+5])
            uiSum += uint(iTemp*iTemp) >> uiShift
            iTemp = int(piOrg[n+6] - piCur[n+6])
            uiSum += uint(iTemp*iTemp) >> uiShift
            iTemp = int(piOrg[n+7] - piCur[n+7])
            uiSum += uint(iTemp*iTemp) >> uiShift
            iTemp = int(piOrg[n+8] - piCur[n+8])
            uiSum += uint(iTemp*iTemp) >> uiShift
            iTemp = int(piOrg[n+9] - piCur[n+9])
            uiSum += uint(iTemp*iTemp) >> uiShift
            iTemp = int(piOrg[n+10] - piCur[n+10])
            uiSum += uint(iTemp*iTemp) >> uiShift
            iTemp = int(piOrg[n+11] - piCur[n+11])
            uiSum += uint(iTemp*iTemp) >> uiShift
            iTemp = int(piOrg[n+12] - piCur[n+12])
            uiSum += uint(iTemp*iTemp) >> uiShift
            iTemp = int(piOrg[n+13] - piCur[n+13])
            uiSum += uint(iTemp*iTemp) >> uiShift
            iTemp = int(piOrg[n+14] - piCur[n+14])
            uiSum += uint(iTemp*iTemp) >> uiShift
            iTemp = int(piOrg[n+15] - piCur[n+15])
            uiSum += uint(iTemp*iTemp) >> uiShift

        }
        piOrg = piOrg[iStrideOrg:]
        piCur = piCur[iStrideCur:]
    }

    return (uiSum)
}

func (this *TEncRdCost) xGetSAD(pcDtParam *DistParam) uint {
    if pcDtParam.bApplyWeight {
        return this.xGetSADw(pcDtParam)
    }
    piOrg := pcDtParam.pOrg
    piCur := pcDtParam.pCur
    iRows := pcDtParam.iRows
    iCols := pcDtParam.iCols
    iStrideCur := pcDtParam.iStrideCur
    iStrideOrg := pcDtParam.iStrideOrg

    uiSum := uint(0)

    for ; iRows != 0; iRows-- {
        for n := 0; n < iCols; n++ {
            uiSum += uint(TLibCommon.ABS(int(piOrg[n] - piCur[n])).(int))
        }
        piOrg = piOrg[iStrideOrg:]
        piCur = piCur[iStrideCur:]
    }

    return uiSum >> TLibCommon.DISTORTION_PRECISION_ADJUSTMENT(uint(pcDtParam.bitDepth-8)).(uint)
}

func (this *TEncRdCost) xGetSAD4(pcDtParam *DistParam) uint {
    if pcDtParam.bApplyWeight {
        return this.xGetSADw(pcDtParam)
    }
    piOrg := pcDtParam.pOrg
    piCur := pcDtParam.pCur
    iRows := pcDtParam.iRows
    iSubShift := pcDtParam.iSubShift
    iSubStep := (1 << uint(iSubShift))
    iStrideCur := pcDtParam.iStrideCur * iSubStep
    iStrideOrg := pcDtParam.iStrideOrg * iSubStep

    uiSum := uint(0)

    for ; iRows != 0; iRows -= iSubStep {
        uiSum += uint(TLibCommon.ABS(int(piOrg[0] - piCur[0])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[1] - piCur[1])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[2] - piCur[2])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[3] - piCur[3])).(int))

        piOrg = piOrg[iStrideOrg:]
        piCur = piCur[iStrideCur:]
    }

    uiSum <<= uint(iSubShift)
    return uiSum >> TLibCommon.DISTORTION_PRECISION_ADJUSTMENT(uint(pcDtParam.bitDepth-8)).(uint)
}

func (this *TEncRdCost) xGetSAD8(pcDtParam *DistParam) uint {
    if pcDtParam.bApplyWeight {
        return this.xGetSADw(pcDtParam)
    }
    piOrg := pcDtParam.pOrg
    piCur := pcDtParam.pCur
    iRows := pcDtParam.iRows
    iSubShift := pcDtParam.iSubShift
    iSubStep := (1 << uint(iSubShift))
    iStrideCur := pcDtParam.iStrideCur * iSubStep
    iStrideOrg := pcDtParam.iStrideOrg * iSubStep

    uiSum := uint(0)

    for ; iRows != 0; iRows -= iSubStep {
        uiSum += uint(TLibCommon.ABS(int(piOrg[0] - piCur[0])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[1] - piCur[1])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[2] - piCur[2])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[3] - piCur[3])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[4] - piCur[4])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[5] - piCur[5])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[6] - piCur[6])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[7] - piCur[7])).(int))

        piOrg = piOrg[iStrideOrg:]
        piCur = piCur[iStrideCur:]
    }

    uiSum <<= uint(iSubShift)
    return uiSum >> TLibCommon.DISTORTION_PRECISION_ADJUSTMENT(uint(pcDtParam.bitDepth-8)).(uint)
}

func (this *TEncRdCost) xGetSAD16(pcDtParam *DistParam) uint {
    if pcDtParam.bApplyWeight {
        return this.xGetSADw(pcDtParam)
    }
    piOrg := pcDtParam.pOrg
    piCur := pcDtParam.pCur
    iRows := pcDtParam.iRows
    iSubShift := pcDtParam.iSubShift
    iSubStep := (1 << uint(iSubShift))
    iStrideCur := pcDtParam.iStrideCur * iSubStep
    iStrideOrg := pcDtParam.iStrideOrg * iSubStep

    uiSum := uint(0)

    for ; iRows != 0; iRows -= iSubStep {
        uiSum += uint(TLibCommon.ABS(int(piOrg[0] - piCur[0])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[1] - piCur[1])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[2] - piCur[2])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[3] - piCur[3])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[4] - piCur[4])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[5] - piCur[5])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[6] - piCur[6])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[7] - piCur[7])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[8] - piCur[8])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[9] - piCur[9])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[10] - piCur[10])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[11] - piCur[11])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[12] - piCur[12])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[13] - piCur[13])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[14] - piCur[14])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[15] - piCur[15])).(int))

        piOrg = piOrg[iStrideOrg:]
        piCur = piCur[iStrideCur:]
    }

    uiSum <<= uint(iSubShift)
    return uiSum >> TLibCommon.DISTORTION_PRECISION_ADJUSTMENT(uint(pcDtParam.bitDepth-8)).(uint)
}

func (this *TEncRdCost) xGetSAD32(pcDtParam *DistParam) uint {
    if pcDtParam.bApplyWeight {
        return this.xGetSADw(pcDtParam)
    }
    piOrg := pcDtParam.pOrg
    piCur := pcDtParam.pCur
    iRows := pcDtParam.iRows
    iSubShift := pcDtParam.iSubShift
    iSubStep := (1 << uint(iSubShift))
    iStrideCur := pcDtParam.iStrideCur * iSubStep
    iStrideOrg := pcDtParam.iStrideOrg * iSubStep

    uiSum := uint(0)

    for ; iRows != 0; iRows -= iSubStep {
        uiSum += uint(TLibCommon.ABS(int(piOrg[0] - piCur[0])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[1] - piCur[1])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[2] - piCur[2])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[3] - piCur[3])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[4] - piCur[4])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[5] - piCur[5])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[6] - piCur[6])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[7] - piCur[7])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[8] - piCur[8])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[9] - piCur[9])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[10] - piCur[10])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[11] - piCur[11])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[12] - piCur[12])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[13] - piCur[13])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[14] - piCur[14])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[15] - piCur[15])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[16] - piCur[16])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[17] - piCur[17])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[18] - piCur[18])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[19] - piCur[19])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[20] - piCur[20])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[21] - piCur[21])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[22] - piCur[22])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[23] - piCur[23])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[24] - piCur[24])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[25] - piCur[25])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[26] - piCur[26])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[27] - piCur[27])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[28] - piCur[28])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[29] - piCur[29])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[30] - piCur[30])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[31] - piCur[31])).(int))

        piOrg = piOrg[iStrideOrg:]
        piCur = piCur[iStrideCur:]
    }

    uiSum <<= uint(iSubShift)
    return uiSum >> TLibCommon.DISTORTION_PRECISION_ADJUSTMENT(uint(pcDtParam.bitDepth-8)).(uint)
}

func (this *TEncRdCost) xGetSAD64(pcDtParam *DistParam) uint {
    if pcDtParam.bApplyWeight {
        return this.xGetSADw(pcDtParam)
    }
    piOrg := pcDtParam.pOrg
    piCur := pcDtParam.pCur
    iRows := pcDtParam.iRows
    iSubShift := pcDtParam.iSubShift
    iSubStep := (1 << uint(iSubShift))
    iStrideCur := pcDtParam.iStrideCur * iSubStep
    iStrideOrg := pcDtParam.iStrideOrg * iSubStep

    uiSum := uint(0)

    for ; iRows != 0; iRows -= iSubStep {
        uiSum += uint(TLibCommon.ABS(int(piOrg[0] - piCur[0])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[1] - piCur[1])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[2] - piCur[2])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[3] - piCur[3])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[4] - piCur[4])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[5] - piCur[5])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[6] - piCur[6])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[7] - piCur[7])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[8] - piCur[8])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[9] - piCur[9])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[10] - piCur[10])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[11] - piCur[11])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[12] - piCur[12])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[13] - piCur[13])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[14] - piCur[14])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[15] - piCur[15])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[16] - piCur[16])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[17] - piCur[17])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[18] - piCur[18])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[19] - piCur[19])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[20] - piCur[20])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[21] - piCur[21])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[22] - piCur[22])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[23] - piCur[23])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[24] - piCur[24])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[25] - piCur[25])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[26] - piCur[26])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[27] - piCur[27])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[28] - piCur[28])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[29] - piCur[29])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[30] - piCur[30])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[31] - piCur[31])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[32] - piCur[32])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[33] - piCur[33])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[34] - piCur[34])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[35] - piCur[35])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[36] - piCur[36])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[37] - piCur[37])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[38] - piCur[38])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[39] - piCur[39])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[40] - piCur[40])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[41] - piCur[41])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[42] - piCur[42])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[43] - piCur[43])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[44] - piCur[44])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[45] - piCur[45])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[46] - piCur[46])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[47] - piCur[47])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[48] - piCur[48])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[49] - piCur[49])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[50] - piCur[50])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[51] - piCur[51])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[52] - piCur[52])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[53] - piCur[53])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[54] - piCur[54])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[55] - piCur[55])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[56] - piCur[56])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[57] - piCur[57])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[58] - piCur[58])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[59] - piCur[59])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[60] - piCur[60])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[61] - piCur[61])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[62] - piCur[62])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[63] - piCur[63])).(int))

        piOrg = piOrg[iStrideOrg:]
        piCur = piCur[iStrideCur:]
    }

    uiSum <<= uint(iSubShift)
    return uiSum >> TLibCommon.DISTORTION_PRECISION_ADJUSTMENT(uint(pcDtParam.bitDepth-8)).(uint)
}

func (this *TEncRdCost) xGetSAD16N(pcDtParam *DistParam) uint {
    piOrg := pcDtParam.pOrg
    piCur := pcDtParam.pCur
    iRows := pcDtParam.iRows
    iCols := pcDtParam.iCols
    iSubShift := pcDtParam.iSubShift
    iSubStep := (1 << uint(iSubShift))
    iStrideCur := pcDtParam.iStrideCur * iSubStep
    iStrideOrg := pcDtParam.iStrideOrg * iSubStep

    uiSum := uint(0)

    for ; iRows != 0; iRows -= iSubStep {
        for n := 0; n < iCols; n += 16 {
            uiSum += uint(TLibCommon.ABS(int(piOrg[n+0] - piCur[n+0])).(int))
            uiSum += uint(TLibCommon.ABS(int(piOrg[n+1] - piCur[n+1])).(int))
            uiSum += uint(TLibCommon.ABS(int(piOrg[n+2] - piCur[n+2])).(int))
            uiSum += uint(TLibCommon.ABS(int(piOrg[n+3] - piCur[n+3])).(int))
            uiSum += uint(TLibCommon.ABS(int(piOrg[n+4] - piCur[n+4])).(int))
            uiSum += uint(TLibCommon.ABS(int(piOrg[n+5] - piCur[n+5])).(int))
            uiSum += uint(TLibCommon.ABS(int(piOrg[n+6] - piCur[n+6])).(int))
            uiSum += uint(TLibCommon.ABS(int(piOrg[n+7] - piCur[n+7])).(int))
            uiSum += uint(TLibCommon.ABS(int(piOrg[n+8] - piCur[n+8])).(int))
            uiSum += uint(TLibCommon.ABS(int(piOrg[n+9] - piCur[n+9])).(int))
            uiSum += uint(TLibCommon.ABS(int(piOrg[n+10] - piCur[n+10])).(int))
            uiSum += uint(TLibCommon.ABS(int(piOrg[n+11] - piCur[n+11])).(int))
            uiSum += uint(TLibCommon.ABS(int(piOrg[n+12] - piCur[n+12])).(int))
            uiSum += uint(TLibCommon.ABS(int(piOrg[n+13] - piCur[n+13])).(int))
            uiSum += uint(TLibCommon.ABS(int(piOrg[n+14] - piCur[n+14])).(int))
            uiSum += uint(TLibCommon.ABS(int(piOrg[n+15] - piCur[n+15])).(int))
        }
        piOrg = piOrg[iStrideOrg:]
        piCur = piCur[iStrideCur:]
    }

    uiSum <<= uint(iSubShift)
    return uiSum >> TLibCommon.DISTORTION_PRECISION_ADJUSTMENT(uint(pcDtParam.bitDepth-8)).(uint)
}

//#if AMP_SAD
func (this *TEncRdCost) xGetSAD12(pcDtParam *DistParam) uint {
    if pcDtParam.bApplyWeight {
        return this.xGetSADw(pcDtParam)
    }
    piOrg := pcDtParam.pOrg
    piCur := pcDtParam.pCur
    iRows := pcDtParam.iRows
    iSubShift := pcDtParam.iSubShift
    iSubStep := (1 << uint(iSubShift))
    iStrideCur := pcDtParam.iStrideCur * iSubStep
    iStrideOrg := pcDtParam.iStrideOrg * iSubStep

    uiSum := uint(0)

    for ; iRows != 0; iRows -= iSubStep {
        uiSum += uint(TLibCommon.ABS(int(piOrg[0] - piCur[0])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[1] - piCur[1])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[2] - piCur[2])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[3] - piCur[3])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[4] - piCur[4])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[5] - piCur[5])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[6] - piCur[6])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[7] - piCur[7])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[8] - piCur[8])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[9] - piCur[9])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[10] - piCur[10])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[11] - piCur[11])).(int))

        piOrg = piOrg[iStrideOrg:]
        piCur = piCur[iStrideCur:]
    }

    uiSum <<= uint(iSubShift)
    return uiSum >> TLibCommon.DISTORTION_PRECISION_ADJUSTMENT(uint(pcDtParam.bitDepth-8)).(uint)
}

func (this *TEncRdCost) xGetSAD24(pcDtParam *DistParam) uint {
    if pcDtParam.bApplyWeight {
        return this.xGetSADw(pcDtParam)
    }
    piOrg := pcDtParam.pOrg
    piCur := pcDtParam.pCur
    iRows := pcDtParam.iRows
    iSubShift := pcDtParam.iSubShift
    iSubStep := (1 << uint(iSubShift))
    iStrideCur := pcDtParam.iStrideCur * iSubStep
    iStrideOrg := pcDtParam.iStrideOrg * iSubStep

    uiSum := uint(0)

    for ; iRows != 0; iRows -= iSubStep {
        uiSum += uint(TLibCommon.ABS(int(piOrg[0] - piCur[0])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[1] - piCur[1])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[2] - piCur[2])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[3] - piCur[3])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[4] - piCur[4])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[5] - piCur[5])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[6] - piCur[6])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[7] - piCur[7])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[8] - piCur[8])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[9] - piCur[9])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[10] - piCur[10])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[11] - piCur[11])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[12] - piCur[12])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[13] - piCur[13])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[14] - piCur[14])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[15] - piCur[15])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[16] - piCur[16])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[17] - piCur[17])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[18] - piCur[18])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[19] - piCur[19])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[20] - piCur[20])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[21] - piCur[21])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[22] - piCur[22])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[23] - piCur[23])).(int))

        piOrg = piOrg[iStrideOrg:]
        piCur = piCur[iStrideCur:]
    }

    uiSum <<= uint(iSubShift)
    return uiSum >> TLibCommon.DISTORTION_PRECISION_ADJUSTMENT(uint(pcDtParam.bitDepth-8)).(uint)
}

func (this *TEncRdCost) xGetSAD48(pcDtParam *DistParam) uint {
    if pcDtParam.bApplyWeight {
        return this.xGetSADw(pcDtParam)
    }
    piOrg := pcDtParam.pOrg
    piCur := pcDtParam.pCur
    iRows := pcDtParam.iRows
    iSubShift := pcDtParam.iSubShift
    iSubStep := (1 << uint(iSubShift))
    iStrideCur := pcDtParam.iStrideCur * iSubStep
    iStrideOrg := pcDtParam.iStrideOrg * iSubStep

    uiSum := uint(0)

    for ; iRows != 0; iRows -= iSubStep {
        uiSum += uint(TLibCommon.ABS(int(piOrg[0] - piCur[0])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[1] - piCur[1])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[2] - piCur[2])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[3] - piCur[3])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[4] - piCur[4])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[5] - piCur[5])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[6] - piCur[6])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[7] - piCur[7])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[8] - piCur[8])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[9] - piCur[9])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[10] - piCur[10])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[11] - piCur[11])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[12] - piCur[12])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[13] - piCur[13])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[14] - piCur[14])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[15] - piCur[15])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[16] - piCur[16])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[17] - piCur[17])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[18] - piCur[18])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[19] - piCur[19])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[20] - piCur[20])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[21] - piCur[21])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[22] - piCur[22])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[23] - piCur[23])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[24] - piCur[24])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[25] - piCur[25])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[26] - piCur[26])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[27] - piCur[27])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[28] - piCur[28])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[29] - piCur[29])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[30] - piCur[30])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[31] - piCur[31])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[32] - piCur[32])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[33] - piCur[33])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[34] - piCur[34])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[35] - piCur[35])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[36] - piCur[36])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[37] - piCur[37])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[38] - piCur[38])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[39] - piCur[39])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[40] - piCur[40])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[41] - piCur[41])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[42] - piCur[42])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[43] - piCur[43])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[44] - piCur[44])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[45] - piCur[45])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[46] - piCur[46])).(int))
        uiSum += uint(TLibCommon.ABS(int(piOrg[47] - piCur[47])).(int))

        piOrg = piOrg[iStrideOrg:]
        piCur = piCur[iStrideCur:]
    }

    uiSum <<= uint(iSubShift)
    return uiSum >> TLibCommon.DISTORTION_PRECISION_ADJUSTMENT(uint(pcDtParam.bitDepth-8)).(uint)
}

//#endif

func (this *TEncRdCost) xGetHADs4(pcDtParam *DistParam) uint {
    if pcDtParam.bApplyWeight {
        return this.xGetHADs4w(pcDtParam)
    }
    piOrg := pcDtParam.pOrg
    piCur := pcDtParam.pCur
    iRows := pcDtParam.iRows
    iStrideCur := pcDtParam.iStrideCur
    iStrideOrg := pcDtParam.iStrideOrg
    iStep := pcDtParam.iStep
    var y int
    iOffsetOrg := iStrideOrg << 2
    iOffsetCur := iStrideCur << 2

    uiSum := uint(0)

    for y = 0; y < iRows; y += 4 {
        uiSum += this.xCalcHADs4x4(piOrg, piCur, iStrideOrg, iStrideCur, iStep)
        piOrg = piOrg[iOffsetOrg:]
        piCur = piCur[iOffsetCur:]
    }

    return uiSum >> TLibCommon.DISTORTION_PRECISION_ADJUSTMENT(uint(pcDtParam.bitDepth-8)).(uint)
}

func (this *TEncRdCost) xGetHADs8(pcDtParam *DistParam) uint {
    if pcDtParam.bApplyWeight {
        return this.xGetHADs8w(pcDtParam)
    }
    piOrg := pcDtParam.pOrg
    piCur := pcDtParam.pCur
    iRows := pcDtParam.iRows
    iStrideCur := pcDtParam.iStrideCur
    iStrideOrg := pcDtParam.iStrideOrg
    iStep := pcDtParam.iStep
    var y int

    uiSum := uint(0)

    if iRows == 4 {
        uiSum += this.xCalcHADs4x4(piOrg[0:], piCur, iStrideOrg, iStrideCur, iStep)
        uiSum += this.xCalcHADs4x4(piOrg[4:], piCur[4*iStep:], iStrideOrg, iStrideCur, iStep)
    } else {
        iOffsetOrg := iStrideOrg << 3
        iOffsetCur := iStrideCur << 3
        for y = 0; y < iRows; y += 8 {
            uiSum += this.xCalcHADs8x8(piOrg, piCur, iStrideOrg, iStrideCur, iStep)
            piOrg = piOrg[iOffsetOrg:]
            piCur = piCur[iOffsetCur:]
        }
    }

    return uiSum >> TLibCommon.DISTORTION_PRECISION_ADJUSTMENT(uint(pcDtParam.bitDepth-8)).(uint)
}

func (this *TEncRdCost) xGetHADs(pcDtParam *DistParam) uint {
    if pcDtParam.bApplyWeight {
        return this.xGetHADsw(pcDtParam)
    }
    piOrg := pcDtParam.pOrg
    piCur := pcDtParam.pCur
    iRows := pcDtParam.iRows
    iCols := pcDtParam.iCols
    iStrideCur := pcDtParam.iStrideCur
    iStrideOrg := pcDtParam.iStrideOrg
    iStep := pcDtParam.iStep

    var x, y int

    uiSum := uint(0)

    //#if NS_HAD
    //  if( ( ( iRows % 8 == 0) && (iCols % 8 == 0) && ( iRows == iCols ) ) || ( ( iRows % 8 == 0 ) && (iCols % 8 == 0) && !pcDtParam.bUseNSHAD ) )
    //#else
    if (iRows%8 == 0) && (iCols%8 == 0) {
        //#endif
        iOffsetOrg := iStrideOrg << 3
        iOffsetCur := iStrideCur << 3
        for y = 0; y < iRows; y += 8 {
            for x = 0; x < iCols; x += 8 {
                uiSum += this.xCalcHADs8x8(piOrg[x:], piCur[x*iStep:], iStrideOrg, iStrideCur, iStep)
            }
            piOrg = piOrg[iOffsetOrg:]
            piCur = piCur[iOffsetCur:]
        }
        /*#if NS_HAD
          else if ( ( iCols > 8 ) && ( iCols > iRows ) && pcDtParam.bUseNSHAD )
          {
            Int  iOffsetOrg = iStrideOrg<<2;
            Int  iOffsetCur = iStrideCur<<2;
            for ( y=0; y<iRows; y+= 4 )
            {
              for ( x=0; x<iCols; x+= 16 )
              {
                uiSum += xCalcHADs16x4( piOrg[x:], piCur[x*iStep:], iStrideOrg, iStrideCur, iStep );
              }
              piOrg =piOrg[ iOffsetOrg:];
              piCur =piCur[ iOffsetCur:];
            }
          }
          else if ( ( iRows > 8 ) && ( iCols < iRows ) && pcDtParam.bUseNSHAD )
          {
            Int  iOffsetOrg = iStrideOrg<<4;
            Int  iOffsetCur = iStrideCur<<4;
            for ( y=0; y<iRows; y+= 16 )
            {
              for ( x=0; x<iCols; x+= 4 )
              {
                uiSum += xCalcHADs4x16( piOrg[x:], piCur[x*iStep:], iStrideOrg, iStrideCur, iStep );
              }
              piOrg =piOrg[ iOffsetOrg:];
              piCur =piCur[ iOffsetCur:];
            }
          }
        #endif*/
    } else if (iRows%4 == 0) && (iCols%4 == 0) {
        iOffsetOrg := iStrideOrg << 2
        iOffsetCur := iStrideCur << 2

        for y = 0; y < iRows; y += 4 {
            for x = 0; x < iCols; x += 4 {
                uiSum += this.xCalcHADs4x4(piOrg[x:], piCur[x*iStep:], iStrideOrg, iStrideCur, iStep)
            }
            piOrg = piOrg[iOffsetOrg:]
            piCur = piCur[iOffsetCur:]
        }
    } else if (iRows%2 == 0) && (iCols%2 == 0) {
        iOffsetOrg := iStrideOrg << 1
        iOffsetCur := iStrideCur << 1
        for y = 0; y < iRows; y += 2 {
            for x = 0; x < iCols; x += 2 {
                uiSum += this.xCalcHADs2x2(piOrg[x:], piCur[x*iStep:], iStrideOrg, iStrideCur, iStep)
            }
            piOrg = piOrg[iOffsetOrg:]
            piCur = piCur[iOffsetCur:]
        }
    } else {
        //    assert(false);
    }

    return uiSum >> TLibCommon.DISTORTION_PRECISION_ADJUSTMENT(uint(pcDtParam.bitDepth-8)).(uint)
}

func (this *TEncRdCost) xCalcHADs2x2(piOrg []TLibCommon.Pel, piCur []TLibCommon.Pel, iStrideOrg, iStrideCur, iStep int) uint {
    satd := 0
    var diff, m [4]int
    //    assert( iStep == 1 );
    diff[0] = int(piOrg[0] - piCur[0])
    diff[1] = int(piOrg[1] - piCur[1])
    diff[2] = int(piOrg[iStrideOrg] - piCur[0+iStrideCur])
    diff[3] = int(piOrg[iStrideOrg+1] - piCur[1+iStrideCur])
    m[0] = diff[0] + diff[2]
    m[1] = diff[1] + diff[3]
    m[2] = diff[0] - diff[2]
    m[3] = diff[1] - diff[3]

    satd += int(TLibCommon.ABS(m[0] + m[1]).(int))
    satd += int(TLibCommon.ABS(m[0] - m[1]).(int))
    satd += int(TLibCommon.ABS(m[2] + m[3]).(int))
    satd += int(TLibCommon.ABS(m[2] - m[3]).(int))

    return uint(satd)
}

func (this *TEncRdCost) xCalcHADs4x4(piOrg []TLibCommon.Pel, piCur []TLibCommon.Pel, iStrideOrg, iStrideCur, iStep int) uint {
    var k, satd int
    var diff, m, d [16]int

    //assert( iStep == 1 );
    for k = 0; k < 16; k += 4 {
        diff[k+0] = int(piOrg[0] - piCur[0])
        diff[k+1] = int(piOrg[1] - piCur[1])
        diff[k+2] = int(piOrg[2] - piCur[2])
        diff[k+3] = int(piOrg[3] - piCur[3])

        piCur = piCur[iStrideCur:]
        piOrg = piOrg[iStrideOrg:]
    }

    /*===== hadamard transform =====*/
    m[0] = diff[0] + diff[12]
    m[1] = diff[1] + diff[13]
    m[2] = diff[2] + diff[14]
    m[3] = diff[3] + diff[15]
    m[4] = diff[4] + diff[8]
    m[5] = diff[5] + diff[9]
    m[6] = diff[6] + diff[10]
    m[7] = diff[7] + diff[11]
    m[8] = diff[4] - diff[8]
    m[9] = diff[5] - diff[9]
    m[10] = diff[6] - diff[10]
    m[11] = diff[7] - diff[11]
    m[12] = diff[0] - diff[12]
    m[13] = diff[1] - diff[13]
    m[14] = diff[2] - diff[14]
    m[15] = diff[3] - diff[15]

    d[0] = m[0] + m[4]
    d[1] = m[1] + m[5]
    d[2] = m[2] + m[6]
    d[3] = m[3] + m[7]
    d[4] = m[8] + m[12]
    d[5] = m[9] + m[13]
    d[6] = m[10] + m[14]
    d[7] = m[11] + m[15]
    d[8] = m[0] - m[4]
    d[9] = m[1] - m[5]
    d[10] = m[2] - m[6]
    d[11] = m[3] - m[7]
    d[12] = m[12] - m[8]
    d[13] = m[13] - m[9]
    d[14] = m[14] - m[10]
    d[15] = m[15] - m[11]

    m[0] = d[0] + d[3]
    m[1] = d[1] + d[2]
    m[2] = d[1] - d[2]
    m[3] = d[0] - d[3]
    m[4] = d[4] + d[7]
    m[5] = d[5] + d[6]
    m[6] = d[5] - d[6]
    m[7] = d[4] - d[7]
    m[8] = d[8] + d[11]
    m[9] = d[9] + d[10]
    m[10] = d[9] - d[10]
    m[11] = d[8] - d[11]
    m[12] = d[12] + d[15]
    m[13] = d[13] + d[14]
    m[14] = d[13] - d[14]
    m[15] = d[12] - d[15]

    d[0] = m[0] + m[1]
    d[1] = m[0] - m[1]
    d[2] = m[2] + m[3]
    d[3] = m[3] - m[2]
    d[4] = m[4] + m[5]
    d[5] = m[4] - m[5]
    d[6] = m[6] + m[7]
    d[7] = m[7] - m[6]
    d[8] = m[8] + m[9]
    d[9] = m[8] - m[9]
    d[10] = m[10] + m[11]
    d[11] = m[11] - m[10]
    d[12] = m[12] + m[13]
    d[13] = m[12] - m[13]
    d[14] = m[14] + m[15]
    d[15] = m[15] - m[14]

    for k = 0; k < 16; k++ {
        satd += TLibCommon.ABS(d[k]).(int)
    }
    satd = ((satd + 1) >> 1)

    return uint(satd)
}

func (this *TEncRdCost) xCalcHADs8x8(piOrg []TLibCommon.Pel, piCur []TLibCommon.Pel, iStrideOrg, iStrideCur, iStep int) uint {
    var k, i, j, jj, sad int
    var diff [64]int
    var m1, m2, m3 [8][8]int
    //    assert( iStep == 1 );
    for k = 0; k < 64; k += 8 {
        diff[k+0] = int(piOrg[0] - piCur[0])
        diff[k+1] = int(piOrg[1] - piCur[1])
        diff[k+2] = int(piOrg[2] - piCur[2])
        diff[k+3] = int(piOrg[3] - piCur[3])
        diff[k+4] = int(piOrg[4] - piCur[4])
        diff[k+5] = int(piOrg[5] - piCur[5])
        diff[k+6] = int(piOrg[6] - piCur[6])
        diff[k+7] = int(piOrg[7] - piCur[7])

        piCur = piCur[iStrideCur:]
        piOrg = piOrg[iStrideOrg:]
    }

    //horizontal
    for j = 0; j < 8; j++ {
        jj = j << 3
        m2[j][0] = diff[jj] + diff[jj+4]
        m2[j][1] = diff[jj+1] + diff[jj+5]
        m2[j][2] = diff[jj+2] + diff[jj+6]
        m2[j][3] = diff[jj+3] + diff[jj+7]
        m2[j][4] = diff[jj] - diff[jj+4]
        m2[j][5] = diff[jj+1] - diff[jj+5]
        m2[j][6] = diff[jj+2] - diff[jj+6]
        m2[j][7] = diff[jj+3] - diff[jj+7]

        m1[j][0] = m2[j][0] + m2[j][2]
        m1[j][1] = m2[j][1] + m2[j][3]
        m1[j][2] = m2[j][0] - m2[j][2]
        m1[j][3] = m2[j][1] - m2[j][3]
        m1[j][4] = m2[j][4] + m2[j][6]
        m1[j][5] = m2[j][5] + m2[j][7]
        m1[j][6] = m2[j][4] - m2[j][6]
        m1[j][7] = m2[j][5] - m2[j][7]

        m2[j][0] = m1[j][0] + m1[j][1]
        m2[j][1] = m1[j][0] - m1[j][1]
        m2[j][2] = m1[j][2] + m1[j][3]
        m2[j][3] = m1[j][2] - m1[j][3]
        m2[j][4] = m1[j][4] + m1[j][5]
        m2[j][5] = m1[j][4] - m1[j][5]
        m2[j][6] = m1[j][6] + m1[j][7]
        m2[j][7] = m1[j][6] - m1[j][7]
    }

    //vertical
    for i = 0; i < 8; i++ {
        m3[0][i] = m2[0][i] + m2[4][i]
        m3[1][i] = m2[1][i] + m2[5][i]
        m3[2][i] = m2[2][i] + m2[6][i]
        m3[3][i] = m2[3][i] + m2[7][i]
        m3[4][i] = m2[0][i] - m2[4][i]
        m3[5][i] = m2[1][i] - m2[5][i]
        m3[6][i] = m2[2][i] - m2[6][i]
        m3[7][i] = m2[3][i] - m2[7][i]

        m1[0][i] = m3[0][i] + m3[2][i]
        m1[1][i] = m3[1][i] + m3[3][i]
        m1[2][i] = m3[0][i] - m3[2][i]
        m1[3][i] = m3[1][i] - m3[3][i]
        m1[4][i] = m3[4][i] + m3[6][i]
        m1[5][i] = m3[5][i] + m3[7][i]
        m1[6][i] = m3[4][i] - m3[6][i]
        m1[7][i] = m3[5][i] - m3[7][i]

        m2[0][i] = m1[0][i] + m1[1][i]
        m2[1][i] = m1[0][i] - m1[1][i]
        m2[2][i] = m1[2][i] + m1[3][i]
        m2[3][i] = m1[2][i] - m1[3][i]
        m2[4][i] = m1[4][i] + m1[5][i]
        m2[5][i] = m1[4][i] - m1[5][i]
        m2[6][i] = m1[6][i] + m1[7][i]
        m2[7][i] = m1[6][i] - m1[7][i]
    }

    for i = 0; i < 8; i++ {
        for j = 0; j < 8; j++ {
            sad += TLibCommon.ABS(m2[i][j]).(int)
        }
    }

    sad = ((sad + 2) >> 2)

    return uint(sad)
}

//#if NS_HAD
//func (this *TEncRdCost)  xCalcHADs16x4     ( piOrg []TLibCommon.Pel, piCurr []TLibCommon.Pel,  iStrideOrg,  iStrideCur,  iStep int)uint{}
//func (this *TEncRdCost)  xCalcHADs4x16     ( piOrg []TLibCommon.Pel, piCurr []TLibCommon.Pel,  iStrideOrg,  iStrideCur,  iStep int)uint{}
//#endif
