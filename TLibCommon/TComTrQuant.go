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

import (
    "fmt"
    "math"
)

// ====================================================================================================================
// Constants
// ====================================================================================================================

const QP_BITS = 15

// ====================================================================================================================
// Type definition
// ====================================================================================================================

type EstBitsSbacStruct struct {
    SignificantCoeffGroupBits [NUM_SIG_CG_FLAG_CTX][2]int
    SignificantBits           [NUM_SIG_FLAG_CTX][2]int
    LastXBits                 [32]int
    LastYBits                 [32]int
    GreaterOneBits          [NUM_ONE_FLAG_CTX][2]int
    LevelAbsBits            [NUM_ABS_FLAG_CTX][2]int

    BlockCbpBits     [3 * NUM_QT_CBF_CTX][2]int
    BlockRootCbpBits [4][2]int
    scanZigzag       [2]int ///< flag for zigzag scan
    scanNonZigzag    [2]int ///< flag for non zigzag scan
}

// ====================================================================================================================
// Class definition
// ====================================================================================================================

/// QP class
type QpParam struct {
    m_iQP   int
    m_iPer  int
    m_iRem  int
    m_iBits int
}

//public:

func NewQpParam() *QpParam {
    return &QpParam{}
}

func (this *QpParam) SetQpParam(qpScaled int) {
    this.m_iQP = qpScaled
    this.m_iPer = qpScaled / 6
    this.m_iRem = qpScaled % 6
    this.m_iBits = QP_BITS + this.m_iPer
}

func (this *QpParam) Clear() {
    this.m_iQP = 0
    this.m_iPer = 0
    this.m_iRem = 0
    this.m_iBits = 0
}

func (this *QpParam) GetPer() int {
    return this.m_iPer
}
func (this *QpParam) GetRem() int {
    return this.m_iRem
}
func (this *QpParam) GetBits() int {
    return this.m_iBits
}

func (this *QpParam) GetQp() int {
    return this.m_iQP
}

/// transform and quantization class
type TComTrQuant struct {
    //protected:
    //#if ADAPTIVE_QP_SELECTION
    m_qpDelta       [MAX_QP + 1]int
    m_sliceNsamples [LEVEL_RANGE + 1]int
    m_sliceSumC     [LEVEL_RANGE + 1]float64
    //#endif
    m_plTempCoeff []int

    m_cQP QpParam
    //#if RDOQ_CHROMA_LAMBDA
    m_dLambdaLuma   float64
    m_dLambdaChroma float64
    //#endif
    m_dLambda      float64
    m_uiRDOQOffset uint
    m_uiMaxTrSize  uint
    m_bEnc         bool
    m_useRDOQ      bool
    //#if RDOQ_TRANSFORMSKIP
    m_useRDOQTS bool
    //#endif
    //#if ADAPTIVE_QP_SELECTION
    m_bUseAdaptQpSelect bool
    //#endif
    m_useTransformSkipFast   bool
    m_scalingListEnabledFlag bool
    m_quantCoef              [SCALING_LIST_SIZE_NUM][SCALING_LIST_NUM][SCALING_LIST_REM_NUM][]int     ///< array of quantization matrix coefficient 4x4
    m_dequantCoef            [SCALING_LIST_SIZE_NUM][SCALING_LIST_NUM][SCALING_LIST_REM_NUM][]int     ///< array of dequantization matrix coefficient 4x4
    m_errScale               [SCALING_LIST_SIZE_NUM][SCALING_LIST_NUM][SCALING_LIST_REM_NUM][]float64 ///< array of quantization matrix coefficient 4x4
    m_pcEstBitsSbac          *EstBitsSbacStruct
}

func NewTComTrQuant() *TComTrQuant {
    pTrQuant := &TComTrQuant{}
    pTrQuant.m_cQP.Clear()

    // allocate temporary buffers
    pTrQuant.m_plTempCoeff = make([]int, MAX_CU_SIZE*MAX_CU_SIZE)

    // allocate bit estimation class  (for RDOQ)
    pTrQuant.m_pcEstBitsSbac = &EstBitsSbacStruct{}
    pTrQuant.InitScalingList()

    return pTrQuant
}

func (this *TComTrQuant) GetEstBitsSbac() *EstBitsSbacStruct{
    return this.m_pcEstBitsSbac
}

// initialize class
func (this *TComTrQuant) Init(uiMaxWidth, uiMaxHeight, uiMaxTrSize uint, iSymbolMode int,
    aTable4 []uint, aTable8 []uint, aTableLastPosVlcIndex []uint, bUseRDOQ bool,
    //#if RDOQ_TRANSFORMSKIP
    bUseRDOQTS bool,
    //#endif
    bEnc bool, useTransformSkipFast bool,
    //#if ADAPTIVE_QP_SELECTION
    bUseAdaptQpSelect bool) {
    //#endif
    this.m_uiMaxTrSize = uiMaxTrSize
    this.m_bEnc = bEnc
    this.m_useRDOQ = bUseRDOQ
    //#if RDOQ_TRANSFORMSKIP
    this.m_useRDOQTS = bUseRDOQTS
    //#endif
    //#if ADAPTIVE_QP_SELECTION
    this.m_bUseAdaptQpSelect = bUseAdaptQpSelect
    //#endif
    this.m_useTransformSkipFast = useTransformSkipFast
}

// transform & inverse transform functions
func (this *TComTrQuant) TransformNxN(pcCU *TComDataCU,
    pcResidual []Pel,
    uiStride uint,
    rpcCoeff []TCoeff,
    //#if ADAPTIVE_QP_SELECTION
    rpcArlCoeff []TCoeff,
    //#endif
    uiWidth uint,
    uiHeight uint,
    uiAbsSum *uint,
    eTType TextType,
    uiAbsPartIdx uint,
    useTransformSkip bool) {
    if pcCU.GetCUTransquantBypass1(uiAbsPartIdx) {
        *uiAbsSum = 0
        for k := uint(0); k < uiHeight; k++ {
            for j := uint(0); j < uiWidth; j++ {
                rpcCoeff[k*uiWidth+j] = TCoeff(pcResidual[k*uiStride+j])
                *uiAbsSum += uint(ABS(pcResidual[k*uiStride+j]).(Pel))
            }
        }
        return
    }
    var uiMode uint //luma intra pred
    if eTType == TEXT_LUMA && pcCU.GetPredictionMode1(uiAbsPartIdx) == MODE_INTRA {
        uiMode = uint(pcCU.GetLumaIntraDir1(uiAbsPartIdx))
    } else {
        uiMode = REG_DCT
    }

    *uiAbsSum = 0
    //assert( (pcCU.GetSlice().GetSPS().GetMaxTrSize() >= uiWidth) );
    var bitDepth int
    if eTType == TEXT_LUMA {
        bitDepth = G_bitDepthY
    } else {
        bitDepth = G_bitDepthC
    }

    if useTransformSkip {
        this.xTransformSkip(bitDepth, pcResidual, uiStride, this.m_plTempCoeff, int(uiWidth), int(uiHeight))
    } else {
        this.xT(bitDepth, uiMode, pcResidual, uiStride, this.m_plTempCoeff, int(uiWidth), int(uiHeight))
    }
    this.xQuant(pcCU, this.m_plTempCoeff, rpcCoeff,
        //#if ADAPTIVE_QP_SELECTION
        rpcArlCoeff,
        //#endif
        int(uiWidth), int(uiHeight), uiAbsSum, eTType, uiAbsPartIdx)
}

func (this *TComTrQuant) InvtransformNxN(transQuantBypass bool, eText TextType, uiMode uint, rpcResidual []Pel, uiStride uint, pcCoeff []TCoeff, uiWidth, uiHeight uint, scalingListType int, useTransformSkip bool) {
    if transQuantBypass {
        for k := uint(0); k < uiHeight; k++ {
            for j := uint(0); j < uiWidth; j++ {
                rpcResidual[k*uiStride+j] = Pel(pcCoeff[k*uiWidth+j])
            }
        }
    } else {
        var bitDepth int
        if eText == TEXT_LUMA {
            bitDepth = G_bitDepthY
        } else {
            bitDepth = G_bitDepthC
        }
        this.xDeQuant(bitDepth, pcCoeff, this.m_plTempCoeff, int(uiWidth), int(uiHeight), scalingListType)
        if useTransformSkip == true {
            this.xITransformSkip(bitDepth, this.m_plTempCoeff, rpcResidual, uiStride, int(uiWidth), int(uiHeight))
        } else {
            this.xIT(bitDepth, uiMode, this.m_plTempCoeff, rpcResidual, uiStride, int(uiWidth), int(uiHeight))
        }
    }
}

// Misc functions
func (this *TComTrQuant) SetQPforQuant(qpy int, eTxtType TextType, qpBdOffset, chromaQPOffset int) {
    var qpScaled int

    if eTxtType == TEXT_LUMA {
        qpScaled = qpy + qpBdOffset
    } else {
        qpScaled = CLIP3(-qpBdOffset, 57, qpy+chromaQPOffset).(int)

        if qpScaled < 0 {
            qpScaled = qpScaled + qpBdOffset
        } else {
            qpScaled = int(G_aucChromaScale[qpScaled]) + qpBdOffset
        }
    }
    this.m_cQP.SetQpParam(qpScaled)
}

//#if RDOQ_CHROMA_LAMBDA
func (this *TComTrQuant) SetLambda(dLambdaLuma, dLambdaChroma float64) {
    this.m_dLambdaLuma = dLambdaLuma
    this.m_dLambdaChroma = dLambdaChroma
}
func (this *TComTrQuant) SelectLambda(eTType TextType) {
    if eTType == TEXT_LUMA {
        this.m_dLambda = this.m_dLambdaLuma
    } else {
        this.m_dLambda = this.m_dLambdaChroma
    }
}

//#else
//  Void setLambda(Double dLambda) { m_dLambda = dLambda;}
//#endif
func (this *TComTrQuant) SetRDOQOffset(uiRDOQOffset uint) {
    this.m_uiRDOQOffset = uiRDOQOffset
}

func CalcPatternSigCtx(sigCoeffGroupFlag []uint, posXCG, posYCG uint, width, height int) int {
    if width == 4 && height == 4 {
        return -1
    }
    sigRight := uint8(0)
    sigLower := uint8(0)

    width >>= 2
    height >>= 2
    if int(posXCG) < width-1 {
        sigRight = B2U(sigCoeffGroupFlag[posYCG*uint(width)+posXCG+1] != 0)
    }
    if int(posYCG) < height-1 {
        sigLower = B2U(sigCoeffGroupFlag[(posYCG+1)*uint(width)+posXCG] != 0)
    }
    return int(sigRight + (sigLower << 1))
}

func GetSigCtxInc(patternSigCtx int,
    scanIdx uint,
    posX int,
    posY int,
    log2BlockSize int,
    width int,
    height int,
    textureType TextType) int {
    var ctxIndMap = [16]int{0, 1, 4, 5,
        2, 3, 4, 5,
        6, 6, 8, 8,
        7, 7, 8, 8}

    if posX+posY == 0 {
        return 0
    }

    if log2BlockSize == 2 {
        return ctxIndMap[4*posY+posX]
    }

    var offset int
    if log2BlockSize == 3 {
        if scanIdx == SCAN_DIAG {
            offset = 9
        } else {
            offset = 15
        }
    } else {
        if textureType == TEXT_LUMA {
            offset = 21
        } else {
            offset = 12
        }
    }

    posXinSubset := posX - ((posX >> 2) << 2)
    posYinSubset := posY - ((posY >> 2) << 2)
    cnt := int(0)
    if patternSigCtx == 0 {
        if posXinSubset+posYinSubset <= 2 {
            if posXinSubset+posYinSubset == 0 {
                cnt = 2
            } else {
                cnt = 1
            }
        } else {
            cnt = 0
        }
    } else if patternSigCtx == 1 {
        if posYinSubset <= 1 {
            if posYinSubset == 0 {
                cnt = 2
            } else {
                cnt = 1
            }
        } else {
            cnt = 0
        }
    } else if patternSigCtx == 2 {
        if posXinSubset <= 1 {
            if posXinSubset == 0 {
                cnt = 2
            } else {
                cnt = 1
            }
        } else {
            cnt = 0
        }
    } else {
        cnt = 2
    }

    if textureType == TEXT_LUMA && ((posX>>2)+(posY>>2)) > 0 {
        return 3 + offset + cnt
    }

    return 0 + offset + cnt
}
func GetSigCoeffGroupCtxInc(uiSigCoeffGroupFlag []uint,
    uiCGPosX uint,
    uiCGPosY uint,
    scanIdx uint,
    width, height int) uint {
    uiRight := uint8(0)
    uiLower := uint8(0)

    width >>= 2
    height >>= 2
    if int(uiCGPosX) < width-1 {
        uiRight = B2U(uiSigCoeffGroupFlag[uiCGPosY*uint(width)+uiCGPosX+1] != 0)
    }
    if int(uiCGPosY) < height-1 {
        uiLower = B2U(uiSigCoeffGroupFlag[(uiCGPosY+1)*uint(width)+uiCGPosX] != 0)
    }
    return uint(B2U(uiRight != 0 || uiLower != 0))
}
func (this *TComTrQuant) InitScalingList() {
    for sizeId := uint(0); sizeId < SCALING_LIST_SIZE_NUM; sizeId++ {
        for listId := uint(0); listId < G_scalingListNum[sizeId]; listId++ {
            for qp := uint(0); qp < SCALING_LIST_REM_NUM; qp++ {
                this.m_quantCoef[sizeId][listId][qp] = make([]int, G_scalingListSize[sizeId])
                this.m_dequantCoef[sizeId][listId][qp] = make([]int, G_scalingListSize[sizeId])
                this.m_errScale[sizeId][listId][qp] = make([]float64, G_scalingListSize[sizeId])
            }
        }
    }
    // alias list [1] as [3].
    for qp := uint(0); qp < SCALING_LIST_REM_NUM; qp++ {
        this.m_quantCoef[SCALING_LIST_32x32][3][qp] = this.m_quantCoef[SCALING_LIST_32x32][1][qp]
        this.m_dequantCoef[SCALING_LIST_32x32][3][qp] = this.m_dequantCoef[SCALING_LIST_32x32][1][qp]
        this.m_errScale[SCALING_LIST_32x32][3][qp] = this.m_errScale[SCALING_LIST_32x32][1][qp]
    }
}
func (this *TComTrQuant) DestroyScalingList() {
    for sizeId := uint(0); sizeId < SCALING_LIST_SIZE_NUM; sizeId++ {
        for listId := uint(0); listId < G_scalingListNum[sizeId]; listId++ {
            for qp := uint(0); qp < SCALING_LIST_REM_NUM; qp++ {
                //if(m_quantCoef   [sizeId][listId][qp]) delete [] m_quantCoef   [sizeId][listId][qp];
                //if(m_dequantCoef [sizeId][listId][qp]) delete [] m_dequantCoef [sizeId][listId][qp];
                //if(m_errScale    [sizeId][listId][qp]) delete [] m_errScale    [sizeId][listId][qp];
            }
        }
    }
}

func (this *TComTrQuant) SetErrScaleCoeff(list, size, qp uint) {
    uiLog2TrSize := int(G_aucConvertToBit[G_scalingListSizeX[size]]) + 2
    var bitDepth int
    if size < SCALING_LIST_32x32 && list != 0 && list != 3 {
        bitDepth = G_bitDepthC
    } else {
        bitDepth = G_bitDepthY
    }
    iTransformShift := MAX_TR_DYNAMIC_RANGE - bitDepth - uiLog2TrSize // Represents scaling through forward transform

    uiMaxNumCoeff := G_scalingListSize[size]
    piQuantcoeff := this.GetQuantCoeff(list, qp, size)
    pdErrScale := this.GetErrScaleCoeff(list, size, qp)

    dErrScale := float64(1 << SCALE_BITS)                                // Compensate for scaling of bitcount in Lagrange cost function
    dErrScale = dErrScale * math.Pow(2.0, -2.0*float64(iTransformShift)) // Compensate for scaling through forward transform
    for i := uint(0); i < uiMaxNumCoeff; i++ {
        a := 1 << uint(2*(bitDepth-8))
        pdErrScale[i] = dErrScale / float64(piQuantcoeff[i]) / float64(piQuantcoeff[i]) / float64(a) //DISTORTION_PRECISION_ADJUSTMENT
    }
}
func (this *TComTrQuant) GetErrScaleCoeff(list, size, qp uint) []float64 {
    return this.m_errScale[size][list][qp]
}   //!< get Error Scale Coefficent
func (this *TComTrQuant) GetQuantCoeff(list, qp, size uint) []int {
    return this.m_quantCoef[size][list][qp]
}   //!< get Quant Coefficent
func (this *TComTrQuant) GetDequantCoeff(list, qp, size uint) []int {
    return this.m_dequantCoef[size][list][qp]
}   //!< get DeQuant Coefficent
func (this *TComTrQuant) SetUseScalingList(bUseScalingList bool) {
    this.m_scalingListEnabledFlag = bUseScalingList
}
func (this *TComTrQuant) GetUseScalingList() bool {
    return this.m_scalingListEnabledFlag
}
func (this *TComTrQuant) SetFlatScalingList() {
    var size, list, qp uint

    for size = 0; size < SCALING_LIST_SIZE_NUM; size++ {
        for list = 0; list < G_scalingListNum[size]; list++ {
            for qp = 0; qp < SCALING_LIST_REM_NUM; qp++ {
                this.xSetFlatScalingList(list, size, qp)
                this.SetErrScaleCoeff(list, size, qp)
            }
        }
    }
}
func (this *TComTrQuant) xSetFlatScalingList(list, size, qp uint) {
    var i, num uint
    num = G_scalingListSize[size]
    var quantcoeff []int
    var dequantcoeff []int
    quantScales := G_quantScales[qp]
    invQuantScales := G_invQuantScales[qp] << 4

    quantcoeff = this.GetQuantCoeff(list, qp, size)
    dequantcoeff = this.GetDequantCoeff(list, qp, size)

    for i = 0; i < num; i++ {
        quantcoeff[i] = quantScales
        dequantcoeff[i] = invQuantScales
    }
}

func (this *TComTrQuant) xSetScalingListEnc(scalingList *TComScalingList, listId, sizeId, qp uint) {
    width := G_scalingListSizeX[sizeId]
    height := G_scalingListSizeX[sizeId]
    ratio := int(G_scalingListSizeX[sizeId]) / MIN(MAX_MATRIX_SIZE_NUM, int(G_scalingListSizeX[sizeId])).(int)
    var quantcoeff []int
    coeff := scalingList.GetScalingListAddress(sizeId, listId)
    quantcoeff = this.GetQuantCoeff(listId, qp, sizeId)

    this.ProcessScalingListEnc(coeff, quantcoeff, G_quantScales[qp]<<4, height, width, uint(ratio), int(MIN(MAX_MATRIX_SIZE_NUM, int(G_scalingListSizeX[sizeId])).(int)), uint(scalingList.GetScalingListDC(sizeId, listId)))
}
func (this *TComTrQuant) xSetScalingListDec(scalingList *TComScalingList, listId, sizeId, qp uint) {
    width := G_scalingListSizeX[sizeId]
    height := G_scalingListSizeX[sizeId]
    ratio := int(G_scalingListSizeX[sizeId]) / MIN(MAX_MATRIX_SIZE_NUM, int(G_scalingListSizeX[sizeId])).(int)
    var dequantcoeff []int
    coeff := scalingList.GetScalingListAddress(sizeId, listId)

    dequantcoeff = this.GetDequantCoeff(listId, qp, sizeId)
    this.ProcessScalingListDec(coeff, dequantcoeff, G_invQuantScales[qp], height, width, uint(ratio), MIN(MAX_MATRIX_SIZE_NUM, int(G_scalingListSizeX[sizeId])).(int), uint(scalingList.GetScalingListDC(sizeId, listId)))
}

func (this *TComTrQuant) SetScalingList(scalingList *TComScalingList) {
    var size, list, qp uint

    for size = 0; size < SCALING_LIST_SIZE_NUM; size++ {
        for list = 0; list < G_scalingListNum[size]; list++ {
            for qp = 0; qp < SCALING_LIST_REM_NUM; qp++ {
                this.xSetScalingListEnc(scalingList, list, size, qp)
                this.xSetScalingListDec(scalingList, list, size, qp)
                this.SetErrScaleCoeff(list, size, qp)
            }
        }
    }
}
func (this *TComTrQuant) SetScalingListDec(scalingList *TComScalingList) {
    var size, list, qp uint

    for size = 0; size < SCALING_LIST_SIZE_NUM; size++ {
        for list = 0; list < G_scalingListNum[size]; list++ {
            for qp = 0; qp < SCALING_LIST_REM_NUM; qp++ {
                this.xSetScalingListDec(scalingList, list, size, qp)
            }
        }
    }
}

func (this *TComTrQuant) ProcessScalingListEnc(coeff []int, quantcoeff []int, quantScales int, height, width, ratio uint, sizuNum int, dc uint) {
    fmt.Printf("ProcessScalingListEnc Empty Func\n")
}
func (this *TComTrQuant) ProcessScalingListDec(coeff []int, dequantcoeff []int, invQuantScales int, height, width, ratio uint, sizuNum int, dc uint) {
    for j := uint(0); j < height; j++ {
        for i := uint(0); i < width; i++ {
            dequantcoeff[j*width+i] = invQuantScales * coeff[uint(sizuNum)*(j/ratio)+i/ratio]
        }
    }
    if ratio > 1 {
        dequantcoeff[0] = invQuantScales * int(dc)
    }
}

//#if ADAPTIVE_QP_SELECTION
func (this *TComTrQuant) InitSliceQpDelta() {
    for qp := 0; qp <= MAX_QP; qp++ {
        if qp < 17 {
            this.m_qpDelta[qp] = 0
        } else {
            this.m_qpDelta[qp] = 1
        }
    }
}
func (this *TComTrQuant) StoreSliceQpNext(pcSlice *TComSlice) {
    qpBase := pcSlice.GetSliceQpBase()
    sliceQpused := pcSlice.GetSliceQp()
    var sliceQpnext int
    var alpha float64
    if qpBase < 17 {
        alpha = 0.5
    } else {
        alpha = 1
    }

    cnt := 0
    for u := 1; u <= LEVEL_RANGE; u++ {
        cnt += this.m_sliceNsamples[u]
    }

    if !this.m_useRDOQ {
        sliceQpused = qpBase
        alpha = 0.5
    }

    if cnt > 120 {
        sum := float64(0)
        k := 0
        for u := 1; u < LEVEL_RANGE; u++ {
            sum += float64(u) * this.m_sliceSumC[u]
            k += u * u * this.m_sliceNsamples[u]
        }

        var v int
        var q [MAX_QP + 1]float64
        for v = 0; v <= MAX_QP; v++ {
            q[v] = float64(G_invQuantScales[v%6]*(1<<uint(v/6))) / 64
        }

        qnext := sum / float64(k) * q[sliceQpused] / (1 << ARL_C_PRECISION)

        for v = 0; v < MAX_QP; v++ {
            if qnext < alpha*q[v]+(1-alpha)*q[v+1] {
                break
            }
        }
        sliceQpnext = CLIP3(sliceQpused-3, sliceQpused+3, v).(int)
    } else {
        sliceQpnext = sliceQpused
    }

    this.m_qpDelta[qpBase] = sliceQpnext - qpBase
}
func (this *TComTrQuant) ClearSliceARLCnt() {
    for i := 0; i < LEVEL_RANGE+1; i++ {
        this.m_sliceSumC[i] = 0
        this.m_sliceNsamples[i] = 0
    }
    //memset(m_sliceSumC, 0, sizeof(Double)*(LEVEL_RANGE+1));
    //memset(m_sliceNsamples, 0, sizeof(Int)*(LEVEL_RANGE+1));
}
func (this *TComTrQuant) GetQpDelta(qp int) int {
    return this.m_qpDelta[qp]
}
func (this *TComTrQuant) GetSliceNSamples() []int {
    return this.m_sliceNsamples[:]
}
func (this *TComTrQuant) GetSliceSumC() []float64 {
    return this.m_sliceSumC[:]
}

//#endif
//private:
// forward Transform
func (this *TComTrQuant) xT(bitDepth int, uiMode uint, pResidual []Pel, uiStride uint, plCoeff []int, iWidth, iHeight int) {
    fmt.Printf("xT Empty Func\n")
}

// skipping Transform
func (this *TComTrQuant) xTransformSkip(bitDepth int, piBlkResi []Pel, uiStride uint, psCoeff []int, width, height int) {
    fmt.Printf("xTransformSkip Empty Func\n")
}

func (this *TComTrQuant) signBitHidingHDQ(pcCU *TComDataCU, pQCoef []TCoeff, pCoef []TCoeff, scan *uint, deltaU *int, width, height int) {
    fmt.Printf("signBitHidingHDQ Empty Func\n")
}

// quantization
func (this *TComTrQuant) xQuant(pcCU *TComDataCU,
    pSrc []int,
    pDes []TCoeff,
    //#if ADAPTIVE_QP_SELECTION
    pArlDes []TCoeff,
    //#endif
    iWidth int,
    iHeight int,
    uiAcSum *uint,
    eTType TextType,
    uiAbsPartIdx uint) {
    fmt.Printf("xQuant Empty Func\n")
}

// RDOQ functions

func (this *TComTrQuant) xRateDistOptQuant(pcCU *TComDataCU,
    plSrcCoeff []int,
    piDstCoeff []TCoeff,
    //#if ADAPTIVE_QP_SELECTION
    piArlDstCoeff []int,
    //#endif
    uiWidth uint,
    uiHeight uint,
    uiAbsSum *uint,
    eTType TextType,
    uiAbsPartIdx uint) {
    fmt.Printf("xRateDistOptQuant Empty Func\n")
}
func (this *TComTrQuant) xGetCodedLevel(rd64CodedCost *float64,
    rd64CodedCost0 *float64,
    rd64CodedCostSig *float64,
    lLevelDouble int,
    uiMaxAbsLevel uint,
    ui16CtxNumSig uint16,
    ui16CtxNumOne uint16,
    ui16CtxNumAbs uint16,
    ui16AbsGoRice uint16,
    c1Idx uint,
    c2Idx uint,
    iQBits int,
    dTemp float64,
    bLast bool) uint {
    dCurrCostSig := float64(0)
    uiBestAbsLevel := uint(0)

    if !bLast && uiMaxAbsLevel < 3 {
        *rd64CodedCostSig = this.xGetRateSigCoef(0, ui16CtxNumSig)
        *rd64CodedCost = *rd64CodedCost0 + *rd64CodedCostSig
        if uiMaxAbsLevel == 0 {
            return uiBestAbsLevel
        }
    } else {
        *rd64CodedCost = MAX_DOUBLE
    }

    if !bLast {
        dCurrCostSig = this.xGetRateSigCoef(1, ui16CtxNumSig)
    }

    var uiMinAbsLevel uint
    if uiMaxAbsLevel > 1 {
        uiMinAbsLevel = uiMaxAbsLevel - 1
    } else {
        uiMinAbsLevel = 1
    }

    for uiAbsLevel := int(uiMaxAbsLevel); uiAbsLevel >= int(uiMinAbsLevel); uiAbsLevel-- {
        dErr := float64(lLevelDouble - (uiAbsLevel << uint(iQBits)))
        dCurrCost := dErr*dErr*dTemp + this.xGetICRateCost(uint(uiAbsLevel), ui16CtxNumOne, ui16CtxNumAbs, ui16AbsGoRice, c1Idx, c2Idx)
        dCurrCost += dCurrCostSig

        if dCurrCost < *rd64CodedCost {
            uiBestAbsLevel = uint(uiAbsLevel)
            *rd64CodedCost = dCurrCost
            *rd64CodedCostSig = dCurrCostSig
        }
    }

    return uiBestAbsLevel
}
func (this *TComTrQuant) xGetICRateCost(uiAbsLevel uint,
    ui16CtxNumOne uint16,
    ui16CtxNumAbs uint16,
    ui16AbsGoRice uint16,
    c1Idx uint,
    c2Idx uint) float64 {
    iRate := this.xGetIEPRate()
    var baseLevel uint
    if c1Idx < C1FLAG_NUMBER {
        baseLevel = 2 + uint(B2U(c2Idx < C2FLAG_NUMBER))
    } else {
        baseLevel = 1
    }

    if uiAbsLevel >= baseLevel {
        symbol := uiAbsLevel - baseLevel
        var length uint
        if symbol < (COEF_REMAIN_BIN_REDUCTION << ui16AbsGoRice) {
            length = symbol >> ui16AbsGoRice
            iRate += float64((length + 1 + uint(ui16AbsGoRice)) << 15)
        } else {
            length = uint(ui16AbsGoRice)
            symbol = symbol - (COEF_REMAIN_BIN_REDUCTION << ui16AbsGoRice)
            for symbol >= (1 << length) {
                symbol -= (1 << (length))
                length++
            }
            iRate += float64((COEF_REMAIN_BIN_REDUCTION + length + 1 - uint(ui16AbsGoRice) + length) << 15)
        }
        if c1Idx < C1FLAG_NUMBER {
            iRate += float64(this.m_pcEstBitsSbac.GreaterOneBits[ui16CtxNumOne][1])

            if c2Idx < C2FLAG_NUMBER {
                iRate += float64(this.m_pcEstBitsSbac.LevelAbsBits[ui16CtxNumAbs][1])
            }
        }
    } else {
        if uiAbsLevel == 1 {
            iRate += float64(this.m_pcEstBitsSbac.GreaterOneBits[ui16CtxNumOne][0])
        } else if uiAbsLevel == 2 {
            iRate += float64(this.m_pcEstBitsSbac.GreaterOneBits[ui16CtxNumOne][1])
            iRate += float64(this.m_pcEstBitsSbac.LevelAbsBits[ui16CtxNumAbs][0])
        } else {
            //assert (0);
        }
    }
    return this.xGetICost(iRate)
}
func (this *TComTrQuant) xGetICRate(uiAbsLevel uint,
    ui16CtxNumOne uint16,
    ui16CtxNumAbs uint16,
    ui16AbsGoRice uint16,
    c1Idx uint,
    c2Idx uint) int {
    iRate := 0
    var baseLevel uint
    if c1Idx < C1FLAG_NUMBER {
        baseLevel = 2 + uint(B2U(c2Idx < C2FLAG_NUMBER))
    } else {
        baseLevel = 1
    }

    if uiAbsLevel >= baseLevel {
        uiSymbol := uiAbsLevel - baseLevel
        uiMaxVlc := G_auiGoRiceRange[ui16AbsGoRice]
        bExpGolomb := (uiSymbol > uiMaxVlc)

        if bExpGolomb {
            uiAbsLevel = uiSymbol - uiMaxVlc
            iEGS := 1
            for uiMax := uint(2); uiAbsLevel >= uiMax; uiMax <<= 1 {
                iEGS += 2
            }
            iRate += iEGS << 15
            uiSymbol = MIN(uiSymbol, (uiMaxVlc + 1)).(uint)
        }

        ui16PrefLen := uint16(uiSymbol>>ui16AbsGoRice) + 1
        ui16NumBins := uint16(MIN(ui16PrefLen, G_auiGoRicePrefixLen[ui16AbsGoRice]).(uint)) + ui16AbsGoRice

        iRate += int(ui16NumBins << 15)

        if c1Idx < C1FLAG_NUMBER {
            iRate += this.m_pcEstBitsSbac.GreaterOneBits[ui16CtxNumOne][1]

            if c2Idx < C2FLAG_NUMBER {
                iRate += this.m_pcEstBitsSbac.LevelAbsBits[ui16CtxNumAbs][1]
            }
        }
    } else {
        if uiAbsLevel == 0 {
            return 0
        } else if uiAbsLevel == 1 {
            iRate += this.m_pcEstBitsSbac.GreaterOneBits[ui16CtxNumOne][0]
        } else if uiAbsLevel == 2 {
            iRate += this.m_pcEstBitsSbac.GreaterOneBits[ui16CtxNumOne][1]
            iRate += this.m_pcEstBitsSbac.LevelAbsBits[ui16CtxNumAbs][0]
        } else {
            //assert(0);
        }
    }
    return iRate
}
func (this *TComTrQuant) xGetRateLast(uiPosX, uiPosY, uiBlkWdth uint) float64 {
    uiCtxX := G_uiGroupIdx[uiPosX]
    uiCtxY := G_uiGroupIdx[uiPosY]
    uiCost := float64(this.m_pcEstBitsSbac.LastXBits[uiCtxX] + this.m_pcEstBitsSbac.LastYBits[uiCtxY])
    if uiCtxX > 3 {
        uiCost += float64(this.xGetIEPRate() * float64((uiCtxX-2)>>1))
    }

    if uiCtxY > 3 {
        uiCost += float64(this.xGetIEPRate() * float64((uiCtxY-2)>>1))
    }
    return this.xGetICost(uiCost)
}
func (this *TComTrQuant) xGetRateSigCoeffGroup(uiSignificanceCoeffGroup, ui16CtxNumSig uint16) float64 {
    return float64(this.xGetICost(float64(this.m_pcEstBitsSbac.SignificantCoeffGroupBits[ui16CtxNumSig][uiSignificanceCoeffGroup])))
}
func (this *TComTrQuant) xGetRateSigCoef(uiSignificance, ui16CtxNumSig uint16) float64 {
    return float64(this.xGetICost(float64(this.m_pcEstBitsSbac.SignificantBits[ui16CtxNumSig][uiSignificance])))
}
func (this *TComTrQuant) xGetICost(dRate float64) float64 {
    return this.m_dLambda * dRate
}
func (this *TComTrQuant) xGetIEPRate() float64 {
    return 32768
}

// dequantization
func (this *TComTrQuant) xDeQuant(bitDepth int, pSrc []TCoeff, pDes []int, iWidth, iHeight, scalingListType int) {
    piQCoef := pSrc
    piCoef := pDes

    if iWidth > int(this.m_uiMaxTrSize) {
        iWidth = int(this.m_uiMaxTrSize)
        iHeight = int(this.m_uiMaxTrSize)
    }

    var iShift, iAdd, iCoeffQ int
    uiLog2TrSize := int(G_aucConvertToBit[iWidth]) + 2

    iTransformShift := MAX_TR_DYNAMIC_RANGE - bitDepth - uiLog2TrSize

    iShift = QUANT_IQUANT_SHIFT - QUANT_SHIFT - iTransformShift

    var clipQCoef TCoeff
    bitRange := MIN(15, int(12+uiLog2TrSize+bitDepth-this.m_cQP.m_iPer)).(int)
    levelLimit := 1 << uint(bitRange)

    if this.GetUseScalingList() {
        iShift += 4
        if iShift > this.m_cQP.m_iPer {
            iAdd = 1 << uint(iShift-this.m_cQP.m_iPer-1)
        } else {
            iAdd = 0
        }
        piDequantCoef := this.GetDequantCoeff(uint(scalingListType), uint(this.m_cQP.m_iRem), uint(uiLog2TrSize-2))

        if iShift > this.m_cQP.m_iPer {
            for n := 0; n < iWidth*iHeight; n++ {
                clipQCoef = CLIP3(TCoeff(-32768), TCoeff(32767), piQCoef[n]).(TCoeff)
                iCoeffQ = ((int(clipQCoef) * piDequantCoef[n]) + iAdd) >> uint(iShift-this.m_cQP.m_iPer)
                piCoef[n] = CLIP3(-32768, 32767, iCoeffQ).(int)
            }
        } else {
            for n := 0; n < iWidth*iHeight; n++ {
                clipQCoef = CLIP3(TCoeff(-levelLimit), TCoeff(levelLimit-1), piQCoef[n]).(TCoeff)
                iCoeffQ = (int(clipQCoef) * piDequantCoef[n]) << uint(this.m_cQP.m_iPer-iShift)
                piCoef[n] = CLIP3(-32768, 32767, iCoeffQ).(int)
            }
        }
    } else {
        iAdd = 1 << uint(iShift-1)
        scale := G_invQuantScales[this.m_cQP.m_iRem] << uint(this.m_cQP.m_iPer)

        for n := 0; n < iWidth*iHeight; n++ {
            clipQCoef = CLIP3(TCoeff(-32768), TCoeff(32767), piQCoef[n]).(TCoeff)
            iCoeffQ = (int(clipQCoef)*scale + iAdd) >> uint(iShift)
            piCoef[n] = CLIP3(-32768, 32767, iCoeffQ).(int)
        }
    }
}

func (this *TComTrQuant) fastInverseDst(tmp []int16, block []int16, shift uint) { // input tmp, output block
    var i int
    var c [4]int
    rnd_factor := 1 << (shift - 1)
    for i = 0; i < 4; i++ {
        // Intermediate Variables
        c[0] = int(tmp[i]) + int(tmp[8+i])
        c[1] = int(tmp[8+i]) + int(tmp[12+i])
        c[2] = int(tmp[i]) - int(tmp[12+i])
        c[3] = 74 * int(tmp[4+i])

        block[4*i+0] = int16(CLIP3(-32768, 32767, (29*c[0]+55*c[1]+c[3]+rnd_factor)>>shift).(int))
        block[4*i+1] = int16(CLIP3(-32768, 32767, (55*c[2]-29*c[1]+c[3]+rnd_factor)>>shift).(int))
        block[4*i+2] = int16(CLIP3(-32768, 32767, (74*(int(tmp[i])-int(tmp[8+i])+int(tmp[12+i]))+rnd_factor)>>shift).(int))
        block[4*i+3] = int16(CLIP3(-32768, 32767, (55*c[0]+29*c[2]-c[3]+rnd_factor)>>shift).(int))
    }
}

func (this *TComTrQuant) partialButterflyInverse4(src []int16, dst []int16, shift uint, line int) {
    var j int
    var E [2]int
    var O [2]int
    add := 1 << (shift - 1)

    for j = 0; j < line; j++ {
        /* Utilizing symmetry properties to the maximum to minimize the number of multiplications */
        O[0] = int(G_aiT4[1][0])*int(src[line+j]) + int(G_aiT4[3][0])*int(src[3*line+j])
        O[1] = int(G_aiT4[1][1])*int(src[line+j]) + int(G_aiT4[3][1])*int(src[3*line+j])
        E[0] = int(G_aiT4[0][0])*int(src[0+j]) + int(G_aiT4[2][0])*int(src[2*line+j])
        E[1] = int(G_aiT4[0][1])*int(src[0+j]) + int(G_aiT4[2][1])*int(src[2*line+j])

        /* Combining even and odd terms at each hierarchy levels to calculate the final spatial domain vector */
        dst[0+j*4] = int16(CLIP3(-32768, 32767, (E[0]+O[0]+add)>>shift).(int))
        dst[1+j*4] = int16(CLIP3(-32768, 32767, (E[1]+O[1]+add)>>shift).(int))
        dst[2+j*4] = int16(CLIP3(-32768, 32767, (E[1]-O[1]+add)>>shift).(int))
        dst[3+j*4] = int16(CLIP3(-32768, 32767, (E[0]-O[0]+add)>>shift).(int))

        //src ++;
        //dst += 4;
    }
}

func (this *TComTrQuant) partialButterfly8(src []int16, dst []int16, shift uint, line int) {
    var j, k int
    var E [4]int
    var O [4]int
    var EE [2]int
    var EO [2]int
    add := 1 << (shift - 1)

    for j = 0; j < line; j++ {
        /* E and O*/
        for k = 0; k < 4; k++ {
            E[k] = int(src[k+j*8]) + int(src[7-k+j*8])
            O[k] = int(src[k+j*8]) - int(src[7-k+j*8])
        }
        /* EE and EO */
        EE[0] = E[0] + E[3]
        EO[0] = E[0] - E[3]
        EE[1] = E[1] + E[2]
        EO[1] = E[1] - E[2]

        dst[0+j] = int16((int(G_aiT8[0][0])*EE[0] + int(G_aiT8[0][1])*EE[1] + add) >> shift)
        dst[4*line+j] = int16((int(G_aiT8[4][0])*EE[0] + int(G_aiT8[4][1])*EE[1] + add) >> shift)
        dst[2*line+j] = int16((int(G_aiT8[2][0])*EO[0] + int(G_aiT8[2][1])*EO[1] + add) >> shift)
        dst[6*line+j] = int16((int(G_aiT8[6][0])*EO[0] + int(G_aiT8[6][1])*EO[1] + add) >> shift)

        dst[line+j] = int16((int(G_aiT8[1][0])*O[0] + int(G_aiT8[1][1])*O[1] + int(G_aiT8[1][2])*O[2] + int(G_aiT8[1][3])*O[3] + add) >> shift)
        dst[3*line+j] = int16((int(G_aiT8[3][0])*O[0] + int(G_aiT8[3][1])*O[1] + int(G_aiT8[3][2])*O[2] + int(G_aiT8[3][3])*O[3] + add) >> shift)
        dst[5*line+j] = int16((int(G_aiT8[5][0])*O[0] + int(G_aiT8[5][1])*O[1] + int(G_aiT8[5][2])*O[2] + int(G_aiT8[5][3])*O[3] + add) >> shift)
        dst[7*line+j] = int16((int(G_aiT8[7][0])*O[0] + int(G_aiT8[7][1])*O[1] + int(G_aiT8[7][2])*O[2] + int(G_aiT8[7][3])*O[3] + add) >> shift)

        //src += 8;
        //dst ++;
    }
}

func (this *TComTrQuant) partialButterflyInverse8(src []int16, dst []int16, shift uint, line int) {
    var j, k int
    var E [4]int
    var O [4]int
    var EE [2]int
    var EO [2]int
    add := 1 << (shift - 1)

    for j = 0; j < line; j++ {
        /* Utilizing symmetry properties to the maximum to minimize the number of multiplications */
        for k = 0; k < 4; k++ {
            O[k] = int(G_aiT8[1][k])*int(src[line+j]) + int(G_aiT8[3][k])*int(src[3*line+j]) +
                int(G_aiT8[5][k])*int(src[5*line+j]) + int(G_aiT8[7][k])*int(src[7*line+j])
        }

        EO[0] = int(G_aiT8[2][0])*int(src[2*line+j]) + int(G_aiT8[6][0])*int(src[6*line+j])
        EO[1] = int(G_aiT8[2][1])*int(src[2*line+j]) + int(G_aiT8[6][1])*int(src[6*line+j])
        EE[0] = int(G_aiT8[0][0])*int(src[0+j]) + int(G_aiT8[4][0])*int(src[4*line+j])
        EE[1] = int(G_aiT8[0][1])*int(src[0+j]) + int(G_aiT8[4][1])*int(src[4*line+j])

        /* Combining even and odd terms at each hierarchy levels to calculate the final spatial domain vector */
        E[0] = EE[0] + EO[0]
        E[3] = EE[0] - EO[0]
        E[1] = EE[1] + EO[1]
        E[2] = EE[1] - EO[1]
        for k = 0; k < 4; k++ {
            dst[k+j*8] = int16(CLIP3(-32768, 32767, (E[k]+O[k]+add)>>shift).(int))
            dst[k+4+j*8] = int16(CLIP3(-32768, 32767, (E[3-k]-O[3-k]+add)>>shift).(int))
        }
        //src ++;
        //dst += 8;
    }
}

func (this *TComTrQuant) partialButterfly16(src []int16, dst []int16, shift uint, line int) {
    var j, k int
    var E [8]int
    var O [8]int
    var EE [4]int
    var EO [4]int
    var EEE [2]int
    var EEO [2]int
    add := 1 << (shift - 1)

    for j = 0; j < line; j++ {
        /* E and O*/
        for k = 0; k < 8; k++ {
            E[k] = int(src[k+j*16]) + int(src[15-k+j*16])
            O[k] = int(src[k+j*16]) - int(src[15-k+j*16])
        }
        /* EE and EO */
        for k = 0; k < 4; k++ {
            EE[k] = E[k] + E[7-k]
            EO[k] = E[k] - E[7-k]
        }
        /* EEE and EEO */
        EEE[0] = EE[0] + EE[3]
        EEO[0] = EE[0] - EE[3]
        EEE[1] = EE[1] + EE[2]
        EEO[1] = EE[1] - EE[2]

        dst[0+j] = int16((int(G_aiT16[0][0])*EEE[0] + int(G_aiT16[0][1])*EEE[1] + add) >> shift)
        dst[8*line+j] = int16((int(G_aiT16[8][0])*EEE[0] + int(G_aiT16[8][1])*EEE[1] + add) >> shift)
        dst[4*line+j] = int16((int(G_aiT16[4][0])*EEO[0] + int(G_aiT16[4][1])*EEO[1] + add) >> shift)
        dst[12*line+j] = int16((int(G_aiT16[12][0])*EEO[0] + int(G_aiT16[12][1])*EEO[1] + add) >> shift)

        for k = 2; k < 16; k += 4 {
            dst[k*line+j] = int16((int(G_aiT16[k][0])*EO[0] + int(G_aiT16[k][1])*EO[1] + int(G_aiT16[k][2])*EO[2] + int(G_aiT16[k][3])*EO[3] + add) >> shift)
        }

        for k = 1; k < 16; k += 2 {
            dst[k*line+j] = int16((int(G_aiT16[k][0])*O[0] + int(G_aiT16[k][1])*O[1] + int(G_aiT16[k][2])*O[2] + int(G_aiT16[k][3])*O[3] +
                int(G_aiT16[k][4])*O[4] + int(G_aiT16[k][5])*O[5] + int(G_aiT16[k][6])*O[6] + int(G_aiT16[k][7])*O[7] + add) >> shift)
        }

        //src += 16;
        //dst ++;
    }
}

func (this *TComTrQuant) partialButterflyInverse16(src []int16, dst []int16, shift uint, line int) {
    var j, k int
    var E [8]int
    var O [8]int
    var EE [4]int
    var EO [4]int
    var EEE [2]int
    var EEO [2]int
    add := 1 << (shift - 1)

    for j = 0; j < line; j++ {
        /* Utilizing symmetry properties to the maximum to minimize the number of multiplications */
        for k = 0; k < 8; k++ {
            O[k] = int(G_aiT16[1][k])*int(src[line+j]) + int(G_aiT16[3][k])*int(src[3*line+j]) +
                int(G_aiT16[5][k])*int(src[5*line+j]) + int(G_aiT16[7][k])*int(src[7*line+j]) +
                int(G_aiT16[9][k])*int(src[9*line+j]) + int(G_aiT16[11][k])*int(src[11*line+j]) +
                int(G_aiT16[13][k])*int(src[13*line+j]) + int(G_aiT16[15][k])*int(src[15*line+j])
        }
        for k = 0; k < 4; k++ {
            EO[k] = int(G_aiT16[2][k])*int(src[2*line+j]) + int(G_aiT16[6][k])*int(src[6*line+j]) +
                int(G_aiT16[10][k])*int(src[10*line+j]) + int(G_aiT16[14][k])*int(src[14*line+j])
        }
        EEO[0] = int(G_aiT16[4][0])*int(src[4*line+j]) + int(G_aiT16[12][0])*int(src[12*line+j])
        EEE[0] = int(G_aiT16[0][0])*int(src[0+j]) + int(G_aiT16[8][0])*int(src[8*line+j])
        EEO[1] = int(G_aiT16[4][1])*int(src[4*line+j]) + int(G_aiT16[12][1])*int(src[12*line+j])
        EEE[1] = int(G_aiT16[0][1])*int(src[0+j]) + int(G_aiT16[8][1])*int(src[8*line+j])

        /* Combining even and odd terms at each hierarchy levels to calculate the final spatial domain vector */
        for k = 0; k < 2; k++ {
            EE[k] = EEE[k] + EEO[k]
            EE[k+2] = EEE[1-k] - EEO[1-k]
        }
        for k = 0; k < 4; k++ {
            E[k] = EE[k] + EO[k]
            E[k+4] = EE[3-k] - EO[3-k]
        }
        for k = 0; k < 8; k++ {
            dst[k+j*16] = int16(CLIP3(-32768, 32767, (E[k]+O[k]+add)>>shift).(int))
            dst[k+8+j*16] = int16(CLIP3(-32768, 32767, (E[7-k]-O[7-k]+add)>>shift).(int))
        }
        //src ++;
        //dst += 16;
    }
}

func (this *TComTrQuant) partialButterfly32(src []int16, dst []int16, shift uint, line int) {
    var j, k int
    var E [16]int
    var O [16]int
    var EE [8]int
    var EO [8]int
    var EEE [4]int
    var EEO [4]int
    var EEEE [2]int
    var EEEO [2]int
    add := 1 << (shift - 1)

    for j = 0; j < line; j++ {
        /* E and O*/
        for k = 0; k < 16; k++ {
            E[k] = int(src[k]) + int(src[31-k+j*32])
            O[k] = int(src[k]) - int(src[31-k+j*32])
        }
        /* EE and EO */
        for k = 0; k < 8; k++ {
            EE[k] = E[k] + E[15-k]
            EO[k] = E[k] - E[15-k]
        }
        /* EEE and EEO */
        for k = 0; k < 4; k++ {
            EEE[k] = EE[k] + EE[7-k]
            EEO[k] = EE[k] - EE[7-k]
        }
        /* EEEE and EEEO */
        EEEE[0] = EEE[0] + EEE[3]
        EEEO[0] = EEE[0] - EEE[3]
        EEEE[1] = EEE[1] + EEE[2]
        EEEO[1] = EEE[1] - EEE[2]

        dst[0+j] = int16((int(G_aiT32[0][0])*EEEE[0] + int(G_aiT32[0][1])*EEEE[1] + add) >> shift)
        dst[16*line+j] = int16((int(G_aiT32[16][0])*EEEE[0] + int(G_aiT32[16][1])*EEEE[1] + add) >> shift)
        dst[8*line+j] = int16((int(G_aiT32[8][0])*EEEO[0] + int(G_aiT32[8][1])*EEEO[1] + add) >> shift)
        dst[24*line+j] = int16((int(G_aiT32[24][0])*EEEO[0] + int(G_aiT32[24][1])*EEEO[1] + add) >> shift)
        for k = 4; k < 32; k += 8 {
            dst[k*line+j] = int16((int(G_aiT32[k][0])*EEO[0] + int(G_aiT32[k][1])*EEO[1] + int(G_aiT32[k][2])*EEO[2] + int(G_aiT32[k][3])*EEO[3] + add) >> shift)
        }
        for k = 2; k < 32; k += 4 {
            dst[k*line+j] = int16((int(G_aiT32[k][0])*EO[0] + int(G_aiT32[k][1])*EO[1] + int(G_aiT32[k][2])*EO[2] + int(G_aiT32[k][3])*EO[3] +
                int(G_aiT32[k][4])*EO[4] + int(G_aiT32[k][5])*EO[5] + int(G_aiT32[k][6])*EO[6] + int(G_aiT32[k][7])*EO[7] + add) >> shift)
        }
        for k = 1; k < 32; k += 2 {
            dst[k*line+j] = int16((int(G_aiT32[k][0])*O[0] + int(G_aiT32[k][1])*O[1] + int(G_aiT32[k][2])*O[2] + int(G_aiT32[k][3])*O[3] +
                int(G_aiT32[k][4])*O[4] + int(G_aiT32[k][5])*O[5] + int(G_aiT32[k][6])*O[6] + int(G_aiT32[k][7])*O[7] +
                int(G_aiT32[k][8])*O[8] + int(G_aiT32[k][9])*O[9] + int(G_aiT32[k][10])*O[10] + int(G_aiT32[k][11])*O[11] +
                int(G_aiT32[k][12])*O[12] + int(G_aiT32[k][13])*O[13] + int(G_aiT32[k][14])*O[14] + int(G_aiT32[k][15])*O[15] + add) >> shift)
        }
        //src += 32;
        //dst ++;
    }
}

func (this *TComTrQuant) partialButterflyInverse32(src []int16, dst []int16, shift uint, line int) {
    var j, k int
    var E [16]int
    var O [16]int
    var EE [8]int
    var EO [8]int
    var EEE [4]int
    var EEO [4]int
    var EEEE [2]int
    var EEEO [2]int
    add := 1 << (shift - 1)

    for j = 0; j < line; j++ {
        /* Utilizing symmetry properties to the maximum to minimize the number of multiplications */
        for k = 0; k < 16; k++ {
            O[k] = int(G_aiT32[1][k])*int(src[line+j]) + int(G_aiT32[3][k])*int(src[3*line+j]) + int(G_aiT32[5][k])*int(src[5*line+j]) + int(G_aiT32[7][k])*int(src[7*line+j]) +
                int(G_aiT32[9][k])*int(src[9*line+j]) + int(G_aiT32[11][k])*int(src[11*line+j]) + int(G_aiT32[13][k])*int(src[13*line+j]) + int(G_aiT32[15][k])*int(src[15*line+j]) +
                int(G_aiT32[17][k])*int(src[17*line+j]) + int(G_aiT32[19][k])*int(src[19*line+j]) + int(G_aiT32[21][k])*int(src[21*line+j]) + int(G_aiT32[23][k])*int(src[23*line+j]) +
                int(G_aiT32[25][k])*int(src[25*line+j]) + int(G_aiT32[27][k])*int(src[27*line+j]) + int(G_aiT32[29][k])*int(src[29*line+j]) + int(G_aiT32[31][k])*int(src[31*line+j])
        }
        for k = 0; k < 8; k++ {
            EO[k] = int(G_aiT32[2][k])*int(src[2*line+j]) + int(G_aiT32[6][k])*int(src[6*line+j]) +
                int(G_aiT32[10][k])*int(src[10*line+j]) + int(G_aiT32[14][k])*int(src[14*line+j]) +
                int(G_aiT32[18][k])*int(src[18*line+j]) + int(G_aiT32[22][k])*int(src[22*line+j]) +
                int(G_aiT32[26][k])*int(src[26*line+j]) + int(G_aiT32[30][k])*int(src[30*line+j])
        }
        for k = 0; k < 4; k++ {
            EEO[k] = int(G_aiT32[4][k])*int(src[4*line+j]) + int(G_aiT32[12][k])*int(src[12*line+j]) +
                int(G_aiT32[20][k])*int(src[20*line+j]) + int(G_aiT32[28][k])*int(src[28*line+j])
        }
        EEEO[0] = int(G_aiT32[8][0])*int(src[8*line+j]) + int(G_aiT32[24][0])*int(src[24*line+j])
        EEEO[1] = int(G_aiT32[8][1])*int(src[8*line+j]) + int(G_aiT32[24][1])*int(src[24*line+j])
        EEEE[0] = int(G_aiT32[0][0])*int(src[0+j]) + int(G_aiT32[16][0])*int(src[16*line+j])
        EEEE[1] = int(G_aiT32[0][1])*int(src[0+j]) + int(G_aiT32[16][1])*int(src[16*line+j])

        /* Combining even and odd terms at each hierarchy levels to calculate the final spatial domain vector */
        EEE[0] = EEEE[0] + EEEO[0]
        EEE[3] = EEEE[0] - EEEO[0]
        EEE[1] = EEEE[1] + EEEO[1]
        EEE[2] = EEEE[1] - EEEO[1]
        for k = 0; k < 4; k++ {
            EE[k] = EEE[k] + EEO[k]
            EE[k+4] = EEE[3-k] - EEO[3-k]
        }
        for k = 0; k < 8; k++ {
            E[k] = EE[k] + EO[k]
            E[k+8] = EE[7-k] - EO[7-k]
        }
        for k = 0; k < 16; k++ {
            dst[k+j*32] = int16(CLIP3(-32768, 32767, (E[k]+O[k]+add)>>shift).(int))
            dst[k+16+j*32] = int16(CLIP3(-32768, 32767, (E[15-k]-O[15-k]+add)>>shift).(int))
        }
        //src ++;
        //dst += 32;
    }
}

/** MxN inverse transform (2D)
*  \param coeff input data (transform coefficients)
*  \param block output data (residual)
*  \param iWidth input data (width of transform)
*  \param iHeight input data (height of transform)
 */
func (this *TComTrQuant) xITrMxN(bitDepth int, coeff []int16, block []int16, iWidth, iHeight int, uiMode uint) {
    shift_1st := uint(SHIFT_INV_1ST)
    shift_2nd := uint(SHIFT_INV_2ND - (bitDepth - 8))

    var tmp [64 * 64]int16
    if iWidth == 4 && iHeight == 4 {
        if uiMode != REG_DCT {
            this.fastInverseDst(coeff, tmp[:], shift_1st) // Inverse DST by FAST Algorithm, coeff input, tmp output
            this.fastInverseDst(tmp[:], block, shift_2nd) // Inverse DST by FAST Algorithm, tmp input, coeff output
        } else {
            this.partialButterflyInverse4(coeff, tmp[:], shift_1st, iWidth)
            this.partialButterflyInverse4(tmp[:], block, shift_2nd, iHeight)
        }
    } else if iWidth == 8 && iHeight == 8 {
        this.partialButterflyInverse8(coeff, tmp[:], shift_1st, iWidth)
        this.partialButterflyInverse8(tmp[:], block, shift_2nd, iHeight)
    } else if iWidth == 16 && iHeight == 16 {
        this.partialButterflyInverse16(coeff, tmp[:], shift_1st, iWidth)
        this.partialButterflyInverse16(tmp[:], block, shift_2nd, iHeight)
    } else if iWidth == 32 && iHeight == 32 {
        this.partialButterflyInverse32(coeff, tmp[:], shift_1st, iWidth)
        this.partialButterflyInverse32(tmp[:], block, shift_2nd, iHeight)
    }
}

// inverse transform
func (this *TComTrQuant) xIT(bitDepth int, uiMode uint, plCoef []int, pResidual []Pel, uiStride uint, iWidth, iHeight int) {
    /*#if MATRIX_MULT
      Int iSize = iWidth;
      xITr(bitDepth, plCoef,pResidual,uiStride,(UInt)iSize,uiMode);
    #else*/
    var i, j int
    {
        var block [64 * 64]int16
        var coeff [64 * 64]int16
        for j = 0; j < iHeight*iWidth; j++ {
            coeff[j] = int16(plCoef[j])
        }
        this.xITrMxN(bitDepth, coeff[:], block[:], iWidth, iHeight, uiMode)
        {
            for j = 0; j < iHeight; j++ {
                for i = 0; i < iWidth; i++ {
                    pResidual[j*int(uiStride)+i] = Pel(block[j*iWidth+i])
                }
                //memcpy( pResidual + j * uiStride, block + j * iWidth, iWidth * sizeof(Short) );
            }
        }
        return
    }
    //#endif
}

// inverse skipping transform
func (this *TComTrQuant) xITransformSkip(bitDepth int, plCoef []int, pResidual []Pel, uiStride uint, width, height int) {
    //assert( width == height );
    uiLog2TrSize := int(G_aucConvertToBit[width]) + 2
    shift := uint(MAX_TR_DYNAMIC_RANGE - bitDepth - uiLog2TrSize)
    var transformSkipShift uint
    var j, k int
    if shift > 0 {
        var offset int
        transformSkipShift = shift
        offset = (1 << (transformSkipShift - 1))
        for j = 0; j < height; j++ {
            for k = 0; k < width; k++ {
                pResidual[j*int(uiStride)+k] = Pel((plCoef[j*width+k] + offset) >> transformSkipShift)
            }
        }
    } else {
        //The case when uiBitDepth >= 13
        transformSkipShift = -shift
        for j = 0; j < height; j++ {
            for k = 0; k < width; k++ {
                pResidual[j*int(uiStride)+k] = Pel(plCoef[j*width+k] << transformSkipShift)
            }
        }
    }
}
