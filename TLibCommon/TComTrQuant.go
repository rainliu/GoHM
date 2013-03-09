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
    //"fmt"
    "math"
)

// ====================================================================================================================
// Constants
// ====================================================================================================================

const QP_BITS = 15
const RDOQ_CHROMA = true

// ====================================================================================================================
// Type definition
// ====================================================================================================================

type EstBitsSbacStruct struct {
    SignificantCoeffGroupBits [NUM_SIG_CG_FLAG_CTX][2]int
    SignificantBits           [NUM_SIG_FLAG_CTX][2]int
    LastXBits                 [32]int
    LastYBits                 [32]int
    GreaterOneBits            [NUM_ONE_FLAG_CTX][2]int
    LevelAbsBits              [NUM_ABS_FLAG_CTX][2]int

    BlockCbpBits     [3 * NUM_QT_CBF_CTX][2]int
    BlockRootCbpBits [4][2]int
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

type coeffGroupRDStats struct {
    iNNZbeforePos0       int
    d64CodedLevelandDist float64 // distortion and level cost only
    d64UncodedDist       float64 // all zero coded block distortion
    d64SigCost           float64
    d64SigCost_0         float64
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

    m_cQP *QpParam
    //#if RDOQ_CHROMA_LAMBDA
    m_dLambdaLuma   float64
    m_dLambdaChroma float64
    //#endif
    m_dLambda      float64
    m_uiRDOQOffset uint
    m_uiMaxTrSize  uint
    m_bEnc         bool
    m_useRDOQ      bool
    m_useRDOQTS    bool
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
    pTrQuant := &TComTrQuant{m_cQP:NewQpParam()}
    pTrQuant.m_cQP.Clear()

    // allocate temporary buffers
    pTrQuant.m_plTempCoeff = make([]int, MAX_CU_SIZE*MAX_CU_SIZE)

    // allocate bit estimation class  (for RDOQ)
    pTrQuant.m_pcEstBitsSbac = &EstBitsSbacStruct{}
    pTrQuant.InitScalingList()

    return pTrQuant
}

func (this *TComTrQuant) GetEstBitsSbac() *EstBitsSbacStruct {
    return this.m_pcEstBitsSbac
}

// initialize class
func (this *TComTrQuant) Init(uiMaxTrSize uint, bUseRDOQ bool,
    bUseRDOQTS bool,
    bEnc bool, useTransformSkipFast bool,
    //#if ADAPTIVE_QP_SELECTION
    bUseAdaptQpSelect bool) {
    //#endif
    this.m_uiMaxTrSize = uiMaxTrSize
    this.m_bEnc = bEnc
    this.m_useRDOQ = bUseRDOQ
    this.m_useRDOQTS = bUseRDOQTS
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
    
    //fmt.Printf("Enter transformNxN with (%d,%d,%d,%d,%d)\n",uiWidth,uiHeight,eTType,uiAbsPartIdx,B2U(useTransformSkip));
    
    if pcCU.GetCUTransquantBypass1(uiAbsPartIdx) {
        *uiAbsSum = 0
        for k := uint(0); k < uiHeight; k++ {
            for j := uint(0); j < uiWidth; j++ {
                rpcCoeff[k*uiWidth+j] = TCoeff(pcResidual[k*uiStride+j])
                *uiAbsSum += uint(ABS(pcResidual[k*uiStride+j]).(Pel))
            }
        }
        //fmt.Printf("Exit transformNxN\n");
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
       
    //fmt.Printf("Exit transformNxN\n");    
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
       
       /*if G_uiPicNo==70 && uiWidth==4 && pcCoeff[0]==-1 {
	      for  k := uint(0); k<uiHeight; k++ {
	        for  j:= uint(0); j<uiWidth; j++ {
	          fmt.Printf("%8d ", this.m_plTempCoeff[k*uiWidth+j]);
	        }
	        fmt.Printf("\n");
	      }
	    }*/
    
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
    //fmt.Printf("Enter calcPatternSigCtx: %d,%d,%d,%d ",posXCG,posYCG,width,height);
    
    if width == 4 && height == 4 {
        return -1
    }
    sigRight := uint8(0)
    sigLower := uint8(0)

    width >>= 2
    height >>= 2
    if int(posXCG) < width-1 {
        sigRight = B2U(sigCoeffGroupFlag[posYCG*uint(width)+posXCG+1] != 0)
    	//fmt.Printf("%d:%d ",posYCG * uint(width) + posXCG + 1, sigCoeffGroupFlag[ posYCG * uint(width) + posXCG + 1 ]);
    }
    if int(posYCG) < height-1 {
        sigLower = B2U(sigCoeffGroupFlag[(posYCG+1)*uint(width)+posXCG] != 0)
    	//fmt.Printf("%d:%d ",(posYCG  + 1 ) * uint(width) + posXCG, sigCoeffGroupFlag[ (posYCG  + 1 ) * uint(width) + posXCG ]);
    }
    //fmt.Printf("Exit calcPatternSigCtx: %d\n",sigRight + (sigLower<<1));
    return int(sigRight + (sigLower << 1))
}

func GetSigCtxInc(patternSigCtx int,
    scanIdx uint,
    posX int,
    posY int,
    log2BlockSize int,
    textureType TextType) int {
    
    //fmt.Printf("Enter getSigCtxInc: %d,%d,%d,%d,%d,%d",patternSigCtx,scanIdx,posX,posY,log2BlockSize,textureType);
  
  
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
    var nsqth, nsqtw int

    if height < width {
        nsqth = 4
    } else {
        nsqth = 1 //height ratio for NSQT
    }
    if width < height {
        nsqtw = 4
    } else {
        nsqtw = 1 //width ratio for NSQT
    }
    for j := int(0); j < int(height); j++ {
        for i := int(0); i < int(width); i++ {
            quantcoeff[j*int(width)+i] = quantScales / coeff[sizuNum*(j*nsqth/int(ratio))+i*nsqtw/int(ratio)]
        }
    }
    if ratio > 1 {
        quantcoeff[0] = quantScales / int(dc)
    }
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

func (this *TComTrQuant) xTrMxN(bitDepth int, block []int16, coeff []int16, iWidth, iHeight int, uiMode uint) {
    shift_1st := int(G_aucConvertToBit[iWidth]) + 1 + bitDepth - 8 // log2(iWidth) - 1 + g_bitDepth - 8
    shift_2nd := int(G_aucConvertToBit[iHeight]) + 8               // log2(iHeight) + 6

    var tmp [64 * 64]int16

    if iWidth == 4 && iHeight == 4 {
        if uiMode != REG_DCT {
            this.fastForwardDst(block, tmp[:], uint(shift_1st)) // Forward DST BY FAST ALGORITHM, block input, tmp output
            this.fastForwardDst(tmp[:], coeff, uint(shift_2nd)) // Forward DST BY FAST ALGORITHM, tmp input, coeff output
        } else {
            this.partialButterfly4(block, tmp[:], uint(shift_1st), iHeight)
            this.partialButterfly4(tmp[:], coeff, uint(shift_2nd), iWidth)
        }
    } else if iWidth == 8 && iHeight == 8 {
        this.partialButterfly8(block, tmp[:], uint(shift_1st), iHeight)
        this.partialButterfly8(tmp[:], coeff, uint(shift_2nd), iWidth)
    } else if iWidth == 16 && iHeight == 16 {
        this.partialButterfly16(block, tmp[:], uint(shift_1st), iHeight)
        this.partialButterfly16(tmp[:], coeff, uint(shift_2nd), iWidth)
    } else if iWidth == 32 && iHeight == 32 {
        this.partialButterfly32(block, tmp[:], uint(shift_1st), iHeight)
        this.partialButterfly32(tmp[:], coeff, uint(shift_2nd), iWidth)
    }
}

//#endif
//private:
// forward Transform
func (this *TComTrQuant) xT(bitDepth int, uiMode uint, piBlkResi []Pel, uiStride uint, psCoeff []int, iWidth, iHeight int) {
    /*#if MATRIX_MULT
      Int iSize = iWidth;
      xTr(bitDepth, piBlkResi,psCoeff,uiStride,(UInt)iSize,uiMode);
    #else*/
    var i, j int
    {
        var block [64 * 64]int16
        var coeff [64 * 64]int16
        {
            for j = 0; j < iHeight; j++ {
                for i = 0; i < iWidth; i++ {
                    block[j*iWidth+i] = int16(piBlkResi[j*int(uiStride)+i])
                    //fmt.Printf("%d ", block[j*iWidth+i]);
                }
                //fmt.Printf("\n");
                //memcpy( block + j * iWidth, piBlkResi + j * uiStride, iWidth * sizeof( Short ) );
            }
            //fmt.Printf("\n");
        }
        this.xTrMxN(bitDepth, block[:], coeff[:], iWidth, iHeight, uiMode)
        for j = 0; j < iHeight*iWidth; j++ {
            psCoeff[j] = int(coeff[j])
            //fmt.Printf("%d ",coeff[ j ]);
        }
        //fmt.Printf("\n");
        return
    }
    //#endif
}

// skipping Transform
func (this *TComTrQuant) xTransformSkip(bitDepth int, piBlkResi []Pel, uiStride uint, psCoeff []int, width, height int) {
    //assert( width == height );
    uiLog2TrSize := uint(G_aucConvertToBit[width] + 2)
    shift := MAX_TR_DYNAMIC_RANGE - bitDepth - int(uiLog2TrSize)
    var transformSkipShift uint
    var j, k int
    if shift >= 0 {
        transformSkipShift = uint(shift)
        for j = 0; j < height; j++ {
            for k = 0; k < width; k++ {
                psCoeff[j*height+k] = int(piBlkResi[j*int(uiStride)+k]) << transformSkipShift
            }
        }
    } else {
        //The case when uiBitDepth > 13
        var offset int
        transformSkipShift = uint(-shift)
        offset = (1 << (transformSkipShift - 1))
        for j = 0; j < height; j++ {
            for k = 0; k < width; k++ {
                psCoeff[j*height+k] = int(int(piBlkResi[j*int(uiStride)+k])+offset) >> transformSkipShift
            }
        }
    }
}

func (this *TComTrQuant) signBitHidingHDQ(pQCoef []TCoeff, pCoef []int, scan []uint, deltaU []int, width, height int) {
    lastCG := -1
    absSum := 0
    var n int

    for subSet := (width*height - 1) >> LOG2_SCAN_SET_SIZE; subSet >= 0; subSet-- {
        subPos := subSet << LOG2_SCAN_SET_SIZE
        firstNZPosInCG := SCAN_SET_SIZE
        lastNZPosInCG := -1
        absSum = 0

        for n = SCAN_SET_SIZE - 1; n >= 0; n-- {
            if pQCoef[scan[n+subPos]] != 0 {
                lastNZPosInCG = n
                break
            }
        }

        for n = 0; n < SCAN_SET_SIZE; n++ {
            if pQCoef[scan[n+subPos]] != 0 {
                firstNZPosInCG = n
                break
            }
        }

        for n = firstNZPosInCG; n <= lastNZPosInCG; n++ {
            absSum += int(pQCoef[scan[n+subPos]])
        }

        if lastNZPosInCG >= 0 && lastCG == -1 {
            lastCG = 1
        }

        if lastNZPosInCG-firstNZPosInCG >= SBH_THRESHOLD {
            var signbit uint
            if pQCoef[scan[subPos+firstNZPosInCG]] > 0 {
                signbit = 0
            } else {
                signbit = 1
            }
            if signbit != uint(absSum&0x1) { //compare signbit with sum_parity
                minCostInc := MAX_INT
                minPos := -1
                finalChange := 0
                curCost := MAX_INT
                curChange := 0

                if lastCG == 1 {
                    n = lastNZPosInCG
                } else {
                    n = SCAN_SET_SIZE - 1
                }
                for ; n >= 0; n-- {
                    blkPos := scan[n+subPos]
                    if pQCoef[blkPos] != 0 {
                        if deltaU[blkPos] > 0 {
                            curCost = -deltaU[blkPos]
                            curChange = 1
                        } else {
                            //curChange =-1;
                            if n == firstNZPosInCG && ABS(pQCoef[blkPos]) == 1 {
                                curCost = MAX_INT
                            } else {
                                curCost = deltaU[blkPos]
                                curChange = -1
                            }
                        }
                    } else {
                        if n < firstNZPosInCG {
                            var thisSignBit uint
                            if pCoef[blkPos] >= 0 {
                                thisSignBit = 0
                            } else {
                                thisSignBit = 1
                            }
                            if thisSignBit != signbit {
                                curCost = MAX_INT
                            } else {
                                curCost = -(deltaU[blkPos])
                                curChange = 1
                            }
                        } else {
                            curCost = -(deltaU[blkPos])
                            curChange = 1
                        }
                    }

                    if curCost < minCostInc {
                        minCostInc = curCost
                        finalChange = curChange
                        minPos = int(blkPos)
                    }
                }   //CG loop

                if pQCoef[minPos] == 32767 || pQCoef[minPos] == -32768 {
                    finalChange = -1
                }

                if pCoef[minPos] >= 0 {
                    pQCoef[minPos] += TCoeff(finalChange)
                } else {
                    pQCoef[minPos] -= TCoeff(finalChange)
                }
            }   // Hide
        }
        if lastCG == 1 {
            lastCG = 0
        }
    }   // TU loop

    return
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
    piCoef := pSrc
    piQCoef := pDes
    //#if ADAPTIVE_QP_SELECTION
    piArlCCoef := pArlDes
    //#endif
    iAdd := 0

    //#if RDOQ_TRANSFORMSKIP
    var useRDOQ bool
    if pcCU.GetTransformSkip2(uiAbsPartIdx, eTType) {
        useRDOQ = this.m_useRDOQTS
    } else {
        useRDOQ = this.m_useRDOQ
    }
    if useRDOQ && (eTType == TEXT_LUMA || RDOQ_CHROMA) {
        //#if ADAPTIVE_QP_SELECTION
        this.xRateDistOptQuant(pcCU, piCoef, pDes, pArlDes, uint(iWidth), uint(iHeight), uiAcSum, eTType, uiAbsPartIdx)
        //#else
        //    xRateDistOptQuant( pcCU, piCoef, pDes, iWidth, iHeight, uiAcSum, eTType, uiAbsPartIdx );
        //#endif
    } else {
        log2BlockSize := G_aucConvertToBit[iWidth] + 2

        scanIdx := pcCU.GetCoefScanIdx(uiAbsPartIdx, uint(iWidth), eTType == TEXT_LUMA, pcCU.IsIntra(uiAbsPartIdx))
        scan := G_auiSigLastScan[scanIdx][log2BlockSize-1]

        var deltaU [32 * 32]int

        //#if ADAPTIVE_QP_SELECTION
        var cQpBase QpParam
        iQpBase := pcCU.GetSlice().GetSliceQpBase()

        var qpScaled, qpBDOffset int
        if eTType == TEXT_LUMA {
            qpBDOffset = pcCU.GetSlice().GetSPS().GetQpBDOffsetY()
        } else {
            qpBDOffset = pcCU.GetSlice().GetSPS().GetQpBDOffsetC()
        }

        if eTType == TEXT_LUMA {
            qpScaled = iQpBase + qpBDOffset
        } else {
            var chromaQPOffset int
            if eTType == TEXT_CHROMA_U {
                chromaQPOffset = pcCU.GetSlice().GetPPS().GetChromaCbQpOffset() + pcCU.GetSlice().GetSliceQpDeltaCb()
            } else {
                chromaQPOffset = pcCU.GetSlice().GetPPS().GetChromaCrQpOffset() + pcCU.GetSlice().GetSliceQpDeltaCr()
            }
            iQpBase = iQpBase + chromaQPOffset

            qpScaled = CLIP3(-qpBDOffset, 57, iQpBase).(int)

            if qpScaled < 0 {
                qpScaled = qpScaled + qpBDOffset
            } else {
                qpScaled = int(G_aucChromaScale[qpScaled]) + qpBDOffset
            }
        }
        cQpBase.SetQpParam(qpScaled)
        //#endif

        uiLog2TrSize := uint(G_aucConvertToBit[iWidth]) + 2
        var scalingListType int
        if pcCU.IsIntra(uiAbsPartIdx) {
            scalingListType = 0 + G_eTTable[eTType]
        } else {
            scalingListType = 3 + G_eTTable[eTType]
        }
        //assert(scalingListType < 6);
        var piQuantCoeff []int
        piQuantCoeff = this.GetQuantCoeff(uint(scalingListType), uint(this.m_cQP.m_iRem), uiLog2TrSize-2)

        var uiBitDepth uint
        if eTType == TEXT_LUMA {
            uiBitDepth = uint(G_bitDepthY)
        } else {
            uiBitDepth = uint(G_bitDepthC)
        }
        iTransformShift := MAX_TR_DYNAMIC_RANGE - int(uiBitDepth) - int(uiLog2TrSize) // Represents scaling through forward transform

        iQBits := QUANT_SHIFT + this.m_cQP.m_iPer + iTransformShift // Right shift of non-RDOQ quantizer;  level = (coeff*uiQ + offset)>>q_bits

        if pcCU.GetSlice().GetSliceType() == I_SLICE {
            iAdd = 171 << uint(iQBits-9)
        } else {
            iAdd = 85 << uint(iQBits-9)
        }

        //#if ADAPTIVE_QP_SELECTION
        iQBits = QUANT_SHIFT + cQpBase.m_iPer + iTransformShift
        if pcCU.GetSlice().GetSliceType() == I_SLICE {
            iAdd = 171 << uint(iQBits-9)
        } else {
            iAdd = 85 << uint(iQBits-9)
        }
        iQBitsC := QUANT_SHIFT + cQpBase.m_iPer + iTransformShift - ARL_C_PRECISION
        iAddC := 1 << uint(iQBitsC-1)
        //#endif

        qBits8 := iQBits - 8
        for n := 0; n < iWidth*iHeight; n++ {
            var iLevel, iSign int
            uiBlockPos := uint(n)
            iLevel = piCoef[uiBlockPos]
            if iLevel < 0 {
                iSign = -1
            } else {
                iSign = 1
            }

            //#if ADAPTIVE_QP_SELECTION
            tmpLevel := int64(ABS(iLevel).(int)) * int64(piQuantCoeff[uiBlockPos])
            if this.m_bUseAdaptQpSelect {
                piArlCCoef[uiBlockPos] = TCoeff((tmpLevel + int64(iAddC)) >> uint(iQBitsC))
            }
            iLevel = int((tmpLevel + int64(iAdd)) >> uint(iQBits))
            deltaU[uiBlockPos] = int((tmpLevel - int64(iLevel<<uint(iQBits))) >> uint(qBits8))
            //#else
            //      iLevel = ((Int64)abs(iLevel) * piQuantCoeff[uiBlockPos] + iAdd ) >> iQBits;
            //      deltaU[uiBlockPos] = (Int)( ((Int64)abs(piCoef[uiBlockPos]) * piQuantCoeff[uiBlockPos] - (iLevel<<iQBits) )>> qBits8 );
            //#endif
            *uiAcSum += uint(iLevel)
            iLevel *= iSign
            piQCoef[uiBlockPos] = TCoeff(CLIP3(-32768, 32767, iLevel).(int))
        }   // for n
        if pcCU.GetSlice().GetPPS().GetSignHideFlag() {
            if *uiAcSum >= 2 {
                this.signBitHidingHDQ(piQCoef, piCoef, scan, deltaU[:], iWidth, iHeight)
            }
        }
    }   //if RDOQ
    //return;

}

// RDOQ functions

func (this *TComTrQuant) xRateDistOptQuant(pcCU *TComDataCU,
    plSrcCoeff []int,
    piDstCoeff []TCoeff,
    //#if ADAPTIVE_QP_SELECTION
    piArlDstCoeff []TCoeff,
    //#endif
    uiWidth uint,
    uiHeight uint,
    uiAbsSum *uint,
    eTType TextType,
    uiAbsPartIdx uint) {
    
    iQBits := this.m_cQP.m_iBits
    dTemp := float64(0)
    uiLog2TrSize := uint(G_aucConvertToBit[uiWidth]) + 2
    uiQ := uint(G_quantScales[this.m_cQP.m_iRem])

	//fmt.Printf("Enter xRateDistOptQuant with (%d,%d,%d,%d)\n", uiWidth,uiHeight,eTType,uiAbsPartIdx);
    
    var uiBitDepth uint
    if eTType == TEXT_LUMA {
        uiBitDepth = uint(G_bitDepthY)
    } else {
        uiBitDepth = uint(G_bitDepthC)
    }
    iTransformShift := MAX_TR_DYNAMIC_RANGE - int(uiBitDepth) - int(uiLog2TrSize) // Represents scaling through forward transform
    uiGoRiceParam := uint(0)
    d64BlockUncodedCost := float64(0)
    uiLog2BlkSize := G_aucConvertToBit[uiWidth] + 2
    uiMaxNumCoeff := uiWidth * uiHeight
    var scalingListType int
    if pcCU.IsIntra(uiAbsPartIdx) {
        scalingListType = 0 + G_eTTable[eTType]
    } else {
        scalingListType = 3 + G_eTTable[eTType]
    }
    //assert(scalingListType < 6);

    iQBits = QUANT_SHIFT + this.m_cQP.m_iPer + iTransformShift // Right shift of non-RDOQ quantizer;  level = (coeff*uiQ + offset)>>q_bits
    dErrScale := float64(0)
    pdErrScaleOrg := this.GetErrScaleCoeff(uint(scalingListType), uiLog2TrSize-2, uint(this.m_cQP.m_iRem))
    piQCoefOrg := this.GetQuantCoeff(uint(scalingListType), uint(this.m_cQP.m_iRem), uiLog2TrSize-2)
    piQCoef := piQCoefOrg
    pdErrScale := pdErrScaleOrg
    //#if ADAPTIVE_QP_SELECTION
    iQBitsC := iQBits - ARL_C_PRECISION
    iAddC := 1 << uint(iQBitsC-1)
    //#endif
    uiScanIdx := pcCU.GetCoefScanIdx(uiAbsPartIdx, uiWidth, eTType == TEXT_LUMA, pcCU.IsIntra(uiAbsPartIdx))

    //#if ADAPTIVE_QP_SELECTION
    for i := uint(0); i < uiMaxNumCoeff; i++ {
        piArlDstCoeff[i] = 0 //, sizeof(Int) *  uiMaxNumCoeff);
    }
    //#endif

    var pdCostCoeff [32 * 32]float64
    var pdCostSig [32 * 32]float64
    var pdCostCoeff0 [32 * 32]float64
    //::memset( pdCostCoeff, 0, sizeof(Double) *  uiMaxNumCoeff );
    //::memset( pdCostSig,   0, sizeof(Double) *  uiMaxNumCoeff );
    var rateIncUp [32 * 32]int
    var rateIncDown [32 * 32]int
    var sigRateDelta [32 * 32]int
    var deltaU [32 * 32]int
    //::memset( rateIncUp,    0, sizeof(Int) *  uiMaxNumCoeff );
    //::memset( rateIncDown,  0, sizeof(Int) *  uiMaxNumCoeff );
    //::memset( sigRateDelta, 0, sizeof(Int) *  uiMaxNumCoeff );
    //::memset( deltaU,       0, sizeof(Int) *  uiMaxNumCoeff );

    var scanCG []uint
    if uiLog2BlkSize > 3 {
        scanCG = G_auiSigLastScan[uiScanIdx][uiLog2BlkSize-2-1]
    } else {
        scanCG = G_auiSigLastScan[uiScanIdx][0]
    }
    if uiLog2BlkSize == 3 {
        scanCG = G_sigLastScan8x8[uiScanIdx][:]
    } else if uiLog2BlkSize == 5 {
        scanCG = G_sigLastScanCG32x32[:]
    }

    uiCGSize := (1 << MLS_CG_SIZE) // 16
    var pdCostCoeffGroupSig [MLS_GRP_NUM]float64
    var uiSigCoeffGroupFlag [MLS_GRP_NUM]uint
    uiNumBlkSide := uiWidth / MLS_CG_SIZE
    iCGLastScanPos := -1

    uiCtxSet := uint(0)
    c1 := 1
    c2 := 0
    d64BaseCost := float64(0)
    iLastScanPos := -1
    dTemp = dErrScale

    c1Idx := uint(0)
    c2Idx := uint(0)
    var baseLevel int

    scan := G_auiSigLastScan[uiScanIdx][uiLog2BlkSize-1]

    //::memset( pdCostCoeffGroupSig,   0, sizeof(Double) * MLS_GRP_NUM );
    //::memset( uiSigCoeffGroupFlag,   0, sizeof(UInt) * MLS_GRP_NUM );

    uiCGNum := uiWidth * uiHeight >> MLS_CG_SIZE
    var iScanPos int
    

    for iCGScanPos := int(uiCGNum) - 1; iCGScanPos >= 0; iCGScanPos-- {
        uiCGBlkPos := scanCG[iCGScanPos]
        uiCGPosY := uiCGBlkPos / uiNumBlkSide
        uiCGPosX := uiCGBlkPos - (uiCGPosY * uiNumBlkSide)
        var rdStats coeffGroupRDStats
        //::memset( &rdStats, 0, sizeof (coeffGroupRDStats));

        patternSigCtx := CalcPatternSigCtx(uiSigCoeffGroupFlag[:], uiCGPosX, uiCGPosY, int(uiWidth), int(uiHeight))
        for iScanPosinCG := int(uiCGSize) - 1; iScanPosinCG >= 0; iScanPosinCG-- {
            iScanPos = iCGScanPos*uiCGSize + iScanPosinCG
            //===== quantization =====
            uiBlkPos := scan[iScanPos]
            // set coeff
            uiQ = uint(piQCoef[uiBlkPos])
            dTemp = pdErrScale[uiBlkPos]
            lLevelDouble := plSrcCoeff[uiBlkPos]
            lLevelDouble = int(MIN(int64(ABS(int(lLevelDouble)).(int))*int64(uiQ), int64(MAX_INT-(1<<uint(iQBits-1)))).(int64))
            //#if ADAPTIVE_QP_SELECTION
            if this.m_bUseAdaptQpSelect {
                piArlDstCoeff[uiBlkPos] = TCoeff((lLevelDouble + iAddC) >> uint(iQBitsC))
            }
            //#endif
            uiMaxAbsLevel := (lLevelDouble + (1 << uint(iQBits-1))) >> uint(iQBits)
			
			//fmt.Printf("%d ", uiMaxAbsLevel);
			
            dErr := float64(lLevelDouble)
            pdCostCoeff0[iScanPos] = dErr * dErr * dTemp
            d64BlockUncodedCost += pdCostCoeff0[iScanPos]
            piDstCoeff[uiBlkPos] = TCoeff(uiMaxAbsLevel)

            if uiMaxAbsLevel > 0 && iLastScanPos < 0 {
                iLastScanPos = iScanPos
                if iScanPos < SCAN_SET_SIZE || eTType != TEXT_LUMA {
                    uiCtxSet = 0
                } else {
                    uiCtxSet = 2
                }
                iCGLastScanPos = iCGScanPos
            }

            if iLastScanPos >= 0 {
                //===== coefficient level estimation =====
                var uiLevel uint
                uiOneCtx := 4*uiCtxSet + uint(c1)
                uiAbsCtx := uiCtxSet + uint(c2)

                if iScanPos == iLastScanPos {
                    uiLevel = this.xGetCodedLevel(&pdCostCoeff[iScanPos], &pdCostCoeff0[iScanPos], &pdCostSig[iScanPos],
                        lLevelDouble, uint(uiMaxAbsLevel), 0, uint16(uiOneCtx), uint16(uiAbsCtx), uint16(uiGoRiceParam),
                        c1Idx, c2Idx, iQBits, dTemp, true)
                    
                    //fmt.Printf("==%d ", uiLevel);
                } else {
                    uiPosY := uint(uiBlkPos) >> uint(uiLog2BlkSize)
                    uiPosX := uiBlkPos - (uiPosY << uint(uiLog2BlkSize))
                    uiCtxSig := GetSigCtxInc(patternSigCtx, uiScanIdx, int(uiPosX), int(uiPosY), int(uiLog2BlkSize), eTType)
                    
                    //fmt.Printf("uiPosX=%d, uiPosY=%d, uiCtxSig=%d\n", uiPosX, uiPosY, uiCtxSig);
                    
                    uiLevel = this.xGetCodedLevel(&pdCostCoeff[iScanPos], &pdCostCoeff0[iScanPos], &pdCostSig[iScanPos],
                        lLevelDouble, uint(uiMaxAbsLevel), uint16(uiCtxSig), uint16(uiOneCtx), uint16(uiAbsCtx), uint16(uiGoRiceParam),
                        c1Idx, c2Idx, iQBits, dTemp, false)
                    sigRateDelta[uiBlkPos] = this.m_pcEstBitsSbac.SignificantBits[uiCtxSig][1] - this.m_pcEstBitsSbac.SignificantBits[uiCtxSig][0]
                
                	//fmt.Printf("!=%d ", uiLevel);
                }
                deltaU[uiBlkPos] = (lLevelDouble - (int(uiLevel) << uint(iQBits))) >> uint(iQBits-8)
                
                //fmt.Printf("%d,%d,%d ",iScanPos, iLastScanPos, uiLevel);
                
                if uiLevel > 0 {
                    rateNow := this.xGetICRate(uiLevel, uint16(uiOneCtx), uint16(uiAbsCtx), uint16(uiGoRiceParam), c1Idx, c2Idx)
                    rateIncUp[uiBlkPos] = this.xGetICRate(uiLevel+1, uint16(uiOneCtx), uint16(uiAbsCtx), uint16(uiGoRiceParam), c1Idx, c2Idx) - rateNow
                    rateIncDown[uiBlkPos] = this.xGetICRate(uiLevel-1, uint16(uiOneCtx), uint16(uiAbsCtx), uint16(uiGoRiceParam), c1Idx, c2Idx) - rateNow
                	//fmt.Printf("%d:[%d,%d,%d] ", uiLevel, this.xGetICRate(uiLevel+1, uint16(uiOneCtx), uint16(uiAbsCtx), uint16(uiGoRiceParam), c1Idx, c2Idx), 
                	//	rateNow, this.xGetICRate(uiLevel-1, uint16(uiOneCtx), uint16(uiAbsCtx), uint16(uiGoRiceParam), c1Idx, c2Idx));
                } else { // uiLevel == 0
                    rateIncUp[uiBlkPos] = this.m_pcEstBitsSbac.GreaterOneBits[uiOneCtx][0]
                    //fmt.Printf("[%d] ", rateIncUp   [ uiBlkPos ]);
                }
                
                piDstCoeff[uiBlkPos] = TCoeff(uiLevel)
                d64BaseCost += pdCostCoeff[iScanPos]

                if c1Idx < C1FLAG_NUMBER {
                    baseLevel = 2 + int(B2U(c2Idx < C2FLAG_NUMBER))
                } else {
                    baseLevel = 1
                }

                if uiLevel >= uint(baseLevel) {
                    if uiLevel > 3*(1<<uiGoRiceParam) {
                        uiGoRiceParam = MIN(uiGoRiceParam+1, uint(4)).(uint)
                    }
                }
                if uiLevel >= 1 {
                    c1Idx++
                }

                //===== update bin model =====
                if uiLevel > 1 {
                    c1 = 0
                    c2 += int(B2U(c2 < 2))
                    c2Idx++
                } else if (c1 < 3) && (c1 > 0) && uiLevel != 0 {
                    c1++
                }

                //===== context set update =====
                if (iScanPos%SCAN_SET_SIZE == 0) && (iScanPos > 0) {
                    c2 = 0
                    uiGoRiceParam = 0

                    c1Idx = 0
                    c2Idx = 0
                    if iScanPos == SCAN_SET_SIZE || eTType != TEXT_LUMA {
                        uiCtxSet = 0
                    } else {
                        uiCtxSet = 2
                    }

                    if c1 == 0 {
                        uiCtxSet++
                    }
                    c1 = 1
                }
            } else {
                d64BaseCost += pdCostCoeff0[iScanPos]
            }
            rdStats.d64SigCost += pdCostSig[iScanPos]
            if iScanPosinCG == 0 {
                rdStats.d64SigCost_0 = pdCostSig[iScanPos]
            }
            if piDstCoeff[uiBlkPos] != 0 {
            	//fmt.Printf("uiBlkPos=%d, piDstCoeff[ uiBlkPos ]=%d\n",uiBlkPos, piDstCoeff[ uiBlkPos ] );
        
                uiSigCoeffGroupFlag[uiCGBlkPos] = 1
                rdStats.d64CodedLevelandDist += pdCostCoeff[iScanPos] - pdCostSig[iScanPos]
                rdStats.d64UncodedDist += pdCostCoeff0[iScanPos]
                if iScanPosinCG != 0 {
                    rdStats.iNNZbeforePos0++
                }
            }
        }   //end for (iScanPosinCG)

        if iCGLastScanPos >= 0 {
            if iCGScanPos != 0 {
                if uiSigCoeffGroupFlag[uiCGBlkPos] == 0 {
                    uiCtxSig := GetSigCoeffGroupCtxInc(uiSigCoeffGroupFlag[:], uiCGPosX, uiCGPosY, int(uiWidth), int(uiHeight))
                    d64BaseCost += this.xGetRateSigCoeffGroup(0, uint16(uiCtxSig)) - rdStats.d64SigCost

                    pdCostCoeffGroupSig[iCGScanPos] = this.xGetRateSigCoeffGroup(0, uint16(uiCtxSig))
                } else {
                    if iCGScanPos < iCGLastScanPos { //skip the last coefficient group, which will be handled together with last position below.
                        if rdStats.iNNZbeforePos0 == 0 {
                            d64BaseCost -= rdStats.d64SigCost_0
                            rdStats.d64SigCost -= rdStats.d64SigCost_0
                        }
                        // rd-cost if SigCoeffGroupFlag = 0, initialization
                        d64CostZeroCG := d64BaseCost

                        // add SigCoeffGroupFlag cost to total cost
                        uiCtxSig := GetSigCoeffGroupCtxInc(uiSigCoeffGroupFlag[:], uiCGPosX, uiCGPosY, int(uiWidth), int(uiHeight))
                        if iCGScanPos < iCGLastScanPos {
                            d64BaseCost += this.xGetRateSigCoeffGroup(1, uint16(uiCtxSig))
                            d64CostZeroCG += this.xGetRateSigCoeffGroup(0, uint16(uiCtxSig))
                            pdCostCoeffGroupSig[iCGScanPos] = this.xGetRateSigCoeffGroup(1, uint16(uiCtxSig))
                        }
						//fmt.Printf("d64CostZeroCG%f =%f %f %f\n",d64CostZeroCG, rdStats.d64UncodedDist, rdStats.d64CodedLevelandDist, rdStats.d64SigCost);
                        // try to convert the current coeff group from non-zero to all-zero
                        d64CostZeroCG += rdStats.d64UncodedDist       // distortion for resetting non-zero levels to zero levels
                        d64CostZeroCG -= rdStats.d64CodedLevelandDist // distortion and level cost for keeping all non-zero levels
                        d64CostZeroCG -= rdStats.d64SigCost           // sig cost for all coeffs, including zero levels and non-zerl levels

                        // if we can save cost, change this block to all-zero block
                        if d64CostZeroCG < d64BaseCost {
                        	//fmt.Printf("d64CostZeroCG%f<d64BaseCost%f,[ uiCGBlkPos =%d]\n",d64CostZeroCG, d64BaseCost, uiCGBlkPos);
                            uiSigCoeffGroupFlag[uiCGBlkPos] = 0
                            d64BaseCost = d64CostZeroCG
                            if iCGScanPos < iCGLastScanPos {
                                pdCostCoeffGroupSig[iCGScanPos] = this.xGetRateSigCoeffGroup(0, uint16(uiCtxSig))
                            }
                            // reset coeffs to 0 in this block
                            for iScanPosinCG := int(uiCGSize) - 1; iScanPosinCG >= 0; iScanPosinCG-- {
                                iScanPos = iCGScanPos*uiCGSize + iScanPosinCG
                                uiBlkPos := scan[iScanPos]

                                if piDstCoeff[uiBlkPos] != 0 {
                                    piDstCoeff[uiBlkPos] = 0
                                    pdCostCoeff[iScanPos] = pdCostCoeff0[iScanPos]
                                    pdCostSig[iScanPos] = 0
                                }
                            }
                        }   // end if ( d64CostAllZeros < d64BaseCost )
                    }
                }   // end if if (uiSigCoeffGroupFlag[ uiCGBlkPos ] == 0)
            } else {
            	//fmt.Printf("uiCGBlkPos=%d\n",uiCGBlkPos);
                uiSigCoeffGroupFlag[uiCGBlkPos] = 1
            }
        }
    }   //end for (iCGScanPos)
	//fmt.Printf("\n");
    //===== estimate last position =====
    if iLastScanPos < 0 {
        return
    }

    d64BestCost := float64(0)
    ui16CtxCbf := uint(0)
    iBestLastIdxP1 := 0
    if !pcCU.IsIntra(uiAbsPartIdx) && eTType == TEXT_LUMA && pcCU.GetTransformIdx1(uiAbsPartIdx) == 0 {
        ui16CtxCbf = 0
        d64BestCost = d64BlockUncodedCost + this.xGetICost(float64(this.m_pcEstBitsSbac.BlockRootCbpBits[ui16CtxCbf][0]))
        d64BaseCost += this.xGetICost(float64(this.m_pcEstBitsSbac.BlockRootCbpBits[ui16CtxCbf][1]))
    } else {
        ui16CtxCbf = pcCU.GetCtxQtCbf(eTType, uint(pcCU.GetTransformIdx1(uiAbsPartIdx)))
        if eTType != 0 {
            ui16CtxCbf = TEXT_CHROMA*NUM_QT_CBF_CTX + ui16CtxCbf
        } else {
            ui16CtxCbf = uint(eTType)*NUM_QT_CBF_CTX + ui16CtxCbf
        }
        d64BestCost = d64BlockUncodedCost + this.xGetICost(float64(this.m_pcEstBitsSbac.BlockCbpBits[ui16CtxCbf][0]))
        d64BaseCost += this.xGetICost(float64(this.m_pcEstBitsSbac.BlockCbpBits[ui16CtxCbf][1]))
    }

    bFoundLast := false
    for iCGScanPos := iCGLastScanPos; iCGScanPos >= 0; iCGScanPos-- {
        uiCGBlkPos := scanCG[iCGScanPos]

        d64BaseCost -= pdCostCoeffGroupSig[iCGScanPos]
        if uiSigCoeffGroupFlag[uiCGBlkPos] != 0 {
            for iScanPosinCG := int(uiCGSize) - 1; iScanPosinCG >= 0; iScanPosinCG-- {
                iScanPos = iCGScanPos*uiCGSize + iScanPosinCG
                if iScanPos > iLastScanPos {
                    continue
                }
                uiBlkPos := scan[iScanPos]

                if piDstCoeff[uiBlkPos] != 0 {
                    uiPosY := uiBlkPos >> uint(uiLog2BlkSize)
                    uiPosX := uiBlkPos - (uiPosY << uint(uiLog2BlkSize))

                    var d64CostLast float64
                    if uiScanIdx == SCAN_VER {
                        d64CostLast = this.xGetRateLast(uiPosY, uiPosX)
                    } else {
                        d64CostLast = this.xGetRateLast(uiPosX, uiPosY)
                    }
                    totalCost := d64BaseCost + d64CostLast - pdCostSig[iScanPos]

                    if totalCost < d64BestCost {
                        iBestLastIdxP1 = iScanPos + 1
                        d64BestCost = totalCost
                    }
                    if piDstCoeff[uiBlkPos] > 1 {
                        bFoundLast = true
                        break
                    }
                    d64BaseCost -= pdCostCoeff[iScanPos]
                    d64BaseCost += pdCostCoeff0[iScanPos]
                } else {
                    d64BaseCost -= pdCostSig[iScanPos]
                }
            }   //end for
            if bFoundLast {
                break
            }
        }   // end if (uiSigCoeffGroupFlag[ uiCGBlkPos ])
    }   // end for

    for scanPos := 0; scanPos < iBestLastIdxP1; scanPos++ {
        blkPos := scan[scanPos]
        level := piDstCoeff[blkPos]
        *uiAbsSum += uint(level)
        if plSrcCoeff[blkPos] < 0 {
            piDstCoeff[blkPos] = -level
        } else {
            piDstCoeff[blkPos] = level
        }
        //fmt.Printf("%d ", piDstCoeff[blkPos]);
    }
    //fmt.Printf("\n");

    //===== clean uncoded coefficients =====
    for scanPos := iBestLastIdxP1; scanPos <= iLastScanPos; scanPos++ {
        piDstCoeff[scan[scanPos]] = 0
    }

    if pcCU.GetSlice().GetPPS().GetSignHideFlag() && *uiAbsSum >= 2 {
        a := G_invQuantScales[this.m_cQP.m_iRem] * G_invQuantScales[this.m_cQP.m_iRem] * (1 << uint(2*this.m_cQP.m_iPer))
        b := 1 << DISTORTION_PRECISION_ADJUSTMENT(2*(uiBitDepth-8)).(uint)
        rdFactor := int64(float64(a) / this.m_dLambda / 16.0 / float64(b) + 0.5)
        lastCG := -1
        absSum := 0
        var n int
        
		//fmt.Printf("rdFactor %d = a %d/%f/%d\n", rdFactor, G_invQuantScales[this.m_cQP.m_iRem] * G_invQuantScales[this.m_cQP.m_iRem] * (1<<uint(2*this.m_cQP.m_iPer)), this.m_dLambda, (1<<DISTORTION_PRECISION_ADJUSTMENT(2*(uiBitDepth-8)).(uint)));
		
        for subSet := int(uiWidth*uiHeight-1) >> LOG2_SCAN_SET_SIZE; subSet >= 0; subSet-- {
            subPos := subSet << LOG2_SCAN_SET_SIZE
            firstNZPosInCG := SCAN_SET_SIZE
            lastNZPosInCG := -1
            absSum = 0

            for n = SCAN_SET_SIZE - 1; n >= 0; n-- {
                if piDstCoeff[scan[n+subPos]] != 0 {
                    lastNZPosInCG = n
                    break
                }
            }

            for n = 0; n < SCAN_SET_SIZE; n++ {
                if piDstCoeff[scan[n+subPos]] != 0 {
                    firstNZPosInCG = n
                    break
                }
            }

            for n = firstNZPosInCG; n <= lastNZPosInCG; n++ {
                absSum += int(piDstCoeff[scan[n+subPos]])
            }

            if lastNZPosInCG >= 0 && lastCG == -1 {
                lastCG = 1
            }

            if lastNZPosInCG-firstNZPosInCG >= SBH_THRESHOLD {
                var signbit uint
                if piDstCoeff[scan[subPos+firstNZPosInCG]] > 0 {
                    signbit = 0
                } else {
                    signbit = 1
                }
                if signbit != uint(absSum&0x1) { // hide but need tune
                    // calculate the cost
                    minCostInc := int64(MAX_INT64)
                    curCost := int64(MAX_INT64)
                    minPos := -1
                    finalChange := 0
                    curChange := 0

                    if lastCG == 1 {
                        n = lastNZPosInCG
                    } else {
                        n = SCAN_SET_SIZE - 1
                    }
                    for ; n >= 0; n-- {
                        uiBlkPos := scan[n+subPos]
                        
                        //fmt.Printf("uiBlkPos=%d, piDstCoeff[ uiBlkPos ]=%d ",uiBlkPos, piDstCoeff[ uiBlkPos ]);
            
                        if piDstCoeff[uiBlkPos] != 0 {
                            costUp := rdFactor*int64(-deltaU[uiBlkPos]) + int64(rateIncUp[uiBlkPos])
                            var costDown int64
                            if ABS(piDstCoeff[uiBlkPos]).(TCoeff) == 1 {
                                costDown = rdFactor*int64(deltaU[uiBlkPos]) + int64(rateIncDown[uiBlkPos])-((1<<15)+int64(sigRateDelta[uiBlkPos]))
                            } else {
                                costDown = rdFactor*int64(deltaU[uiBlkPos]) + int64(rateIncDown[uiBlkPos])
                            }
							
							//fmt.Printf("(%d,%d,%d, %d,%d,%d,%d) ", costUp, costDown, rdFactor, piDstCoeff[uiBlkPos], deltaU[uiBlkPos], rateIncDown[uiBlkPos], sigRateDelta[uiBlkPos]);
							
                            if lastCG == 1 && lastNZPosInCG == n && ABS(piDstCoeff[uiBlkPos]).(TCoeff) == 1 {
                                costDown -= (4 << 15)
                            }

                            if costUp < costDown {
                                curCost = costUp
                                curChange = 1
                            } else {
                                curChange = -1
                                if n == firstNZPosInCG && ABS(piDstCoeff[uiBlkPos]).(TCoeff) == 1 {
                                    curCost = int64(MAX_INT64)
                                } else {
                                    curCost = costDown
                                }
                            }
                            //fmt.Printf("curCost1=%d \n", curCost);
                        } else {
                            curCost = rdFactor*int64(-(ABS(deltaU[uiBlkPos]).(int))) + int64(1<<15) + int64(rateIncUp[uiBlkPos]+sigRateDelta[uiBlkPos])
                            curChange = 1

                            if n < firstNZPosInCG {
                                var thissignbit uint
                                if plSrcCoeff[uiBlkPos] >= 0 {
                                    thissignbit = 0
                                } else {
                                    thissignbit = 1
                                }
                                if thissignbit != signbit {
                                    curCost = int64(MAX_INT64)
                                }
                            }
                            
                            //fmt.Printf("curCost0=%d \n", curCost);
                        }

                        if curCost < minCostInc {
                        	//fmt.Printf("curCost=%d minCostInc=%d finalChange=%d minPos=%d",curCost, minCostInc, finalChange, minPos);
                            minCostInc = curCost
                            finalChange = curChange
                            minPos = int(uiBlkPos)
                        }
                    }

                    if piQCoef[minPos] == 32767 || piQCoef[minPos] == -32768 {
                        finalChange = -1
                    }
					//fmt.Printf("(%d,%d,%d) ",plSrcCoeff[minPos],piDstCoeff[minPos],finalChange);
                    if plSrcCoeff[minPos] >= 0 {
                        piDstCoeff[minPos] += TCoeff(finalChange)
                    } else {
                        piDstCoeff[minPos] -= TCoeff(finalChange)
                    }
                }
            }

            if lastCG == 1 {
                lastCG = 0
            }
        }
    }
    //fmt.Printf("\n");
    //fmt.Printf("Exit xRateDistOptQuant\n");
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
	
	//fmt.Printf("xGetICRate=%d,%d,%d,%d,%d,%d\n",uiAbsLevel,ui16CtxNumOne,ui16CtxNumAbs,ui16AbsGoRice,c1Idx,c2Idx);
	
	
    if uiAbsLevel >= baseLevel {
        uiSymbol := uint(uiAbsLevel - baseLevel)
        uiMaxVlc := uint(G_auiGoRiceRange[ui16AbsGoRice])
        bExpGolomb := (uiSymbol > uiMaxVlc)

        if bExpGolomb {
            uiAbsLevel = uiSymbol - uiMaxVlc
            iEGS := 1
            for uiMax := uint(2); uiAbsLevel >= uiMax; uiMax <<= 1 {
                iEGS += 2
            }
            iRate += (iEGS << 15)
            uiSymbol = MIN(uiSymbol, (uiMaxVlc + 1)).(uint)
        }

        ui16PrefLen := uint16(uiSymbol>>ui16AbsGoRice) + 1
        ui16NumBins := uint16(MIN(ui16PrefLen, uint16(G_auiGoRicePrefixLen[ui16AbsGoRice])).(uint16)) + ui16AbsGoRice
		
        iRate += (int(ui16NumBins) << 15)
		//fmt.Printf("ui16NumBins=%d, iRate=%d ",ui16NumBins,iRate);
		
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
        	//fmt.Printf("ui16CtxNumOne=%d, iRate=%d, sbac=%d ",ui16CtxNumOne,iRate,
           	//			this.m_pcEstBitsSbac.GreaterOneBits[ ui16CtxNumOne ][ 0 ]);
        
            iRate += this.m_pcEstBitsSbac.GreaterOneBits[ui16CtxNumOne][0]
        } else if uiAbsLevel == 2 {
	        //fmt.Printf("ui16CtxNumOne=(%d,%d), iRate=%d, sbac=(%d,%d)",ui16CtxNumOne,ui16CtxNumAbs,iRate,
	        //   this.m_pcEstBitsSbac.GreaterOneBits[ui16CtxNumOne][1],
	        //   this.m_pcEstBitsSbac.LevelAbsBits[ui16CtxNumAbs][0]);
           
            iRate += this.m_pcEstBitsSbac.GreaterOneBits[ui16CtxNumOne][1]
            iRate += this.m_pcEstBitsSbac.LevelAbsBits[ui16CtxNumAbs][0]
        } else {
            //assert(0);
        }
    }
    return iRate
}
func (this *TComTrQuant) xGetRateLast(uiPosX, uiPosY uint) float64 {
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

    if this.GetUseScalingList() {
        iShift += 4

        piDequantCoef := this.GetDequantCoeff(uint(scalingListType), uint(this.m_cQP.m_iRem), uint(uiLog2TrSize-2))
			
		/*if G_uiPicNo==70 && iWidth==4 && pSrc[0]==-1 {
      		fmt.Printf("%8d %d %d %d\n", piDequantCoef[0], this.m_cQP.m_iRem, iShift, this.m_cQP.m_iPer);
    	}*/
    
        if iShift > this.m_cQP.m_iPer {
            iAdd = 1 << uint(iShift-this.m_cQP.m_iPer-1)
            for n := 0; n < iWidth*iHeight; n++ {
                clipQCoef = CLIP3(TCoeff(-32768), TCoeff(32767), piQCoef[n]).(TCoeff)
                iCoeffQ = ((int(clipQCoef) * piDequantCoef[n]) + iAdd) >> uint(iShift-this.m_cQP.m_iPer)
                piCoef[n] = CLIP3(-32768, 32767, iCoeffQ).(int)
            }
        } else {
            for n := 0; n < iWidth*iHeight; n++ {
                clipQCoef = CLIP3(TCoeff(-32768), TCoeff(32767), piQCoef[n]).(TCoeff)
                /*if G_uiPicNo==70 && iWidth==4 && pSrc[0]==-1 {
                	fmt.Printf("%d ", clipQCoef);
                }*/
                iCoeffQ = CLIP3(int(-32768), int(32767), int(clipQCoef)*piDequantCoef[n]).(int)
                /*if G_uiPicNo==70 && iWidth==4 && pSrc[0]==-1 {
                	fmt.Printf("(%d %d %d %d)", iCoeffQ, clipQCoef, piDequantCoef[n], int(clipQCoef) * piDequantCoef[n]);
                }*/
                piCoef[n] = CLIP3(-32768, 32767, iCoeffQ<<uint(this.m_cQP.m_iPer-iShift)).(int)
            	/*if G_uiPicNo==70 && iWidth==4 && pSrc[0]==-1 {
                	fmt.Printf("%d ", piCoef[n]);
                }*/
            }
        }
        /*if G_uiPicNo==70 && iWidth==4 && pSrc[0]==-1 {
	        fmt.Printf("\n");
	    }*/
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

func (this *TComTrQuant) partialButterfly4(src []int16, dst []int16, shift uint, line int) {
    var j int
    var E, O [2]int
    add := 1 << (shift - 1)

    for j = 0; j < line; j++ {
        /* E and O */
        E[0] = int(src[0]) + int(src[3])
        O[0] = int(src[0]) - int(src[3])
        E[1] = int(src[1]) + int(src[2])
        O[1] = int(src[1]) - int(src[2])

        dst[0] = int16((int(G_aiT4[0][0])*E[0] + int(G_aiT4[0][1])*E[1] + add) >> shift)
        dst[2*line] = int16((int(G_aiT4[2][0])*E[0] + int(G_aiT4[2][1])*E[1] + add) >> shift)
        dst[line] = int16((int(G_aiT4[1][0])*O[0] + int(G_aiT4[1][1])*O[1] + add) >> shift)
        dst[3*line] = int16((int(G_aiT4[3][0])*O[0] + int(G_aiT4[3][1])*O[1] + add) >> shift)

        src = src[4:]
        dst = dst[1:]
    }
}

func (this *TComTrQuant) fastForwardDst(block []int16, coeff []int16, shift uint) { // input block, output coeff
    var i int
    var c [4]int
    rnd_factor := 1 << (shift - 1)
    for i = 0; i < 4; i++ {
        // Intermediate Variables
        c[0] = int(block[4*i+0]) + int(block[4*i+3])
        c[1] = int(block[4*i+1]) + int(block[4*i+3])
        c[2] = int(block[4*i+0]) - int(block[4*i+1])
        c[3] = 74 * int(block[4*i+2])

        coeff[i] = int16((29*c[0] + 55*c[1] + c[3] + rnd_factor) >> shift)
        coeff[4+i] = int16((74*(int(block[4*i+0])+int(block[4*i+1])-int(block[4*i+3])) + rnd_factor) >> shift)
        coeff[8+i] = int16((29*c[2] + 55*c[0] - c[3] + rnd_factor) >> shift)
        coeff[12+i] = int16((55*c[2] - 29*c[1] + c[3] + rnd_factor) >> shift)
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

        block[4*i+0] = int16(CLIP3(int(-32768), int(32767), int((29*c[0]+55*c[1]+c[3]+rnd_factor)>>shift)).(int))
        block[4*i+1] = int16(CLIP3(int(-32768), int(32767), int((55*c[2]-29*c[1]+c[3]+rnd_factor)>>shift)).(int))
        block[4*i+2] = int16(CLIP3(int(-32768), int(32767), int((74*(int(tmp[i])-int(tmp[8+i])+int(tmp[12+i]))+rnd_factor)>>shift)).(int))
        block[4*i+3] = int16(CLIP3(int(-32768), int(32767), int((55*c[0]+29*c[2]-c[3]+rnd_factor)>>shift)).(int))
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
        dst[0+j*4] = int16(CLIP3(int(-32768), int(32767), int((E[0]+O[0]+add)>>shift)).(int))
        dst[1+j*4] = int16(CLIP3(int(-32768), int(32767), int((E[1]+O[1]+add)>>shift)).(int))
        dst[2+j*4] = int16(CLIP3(int(-32768), int(32767), int((E[1]-O[1]+add)>>shift)).(int))
        dst[3+j*4] = int16(CLIP3(int(-32768), int(32767), int((E[0]-O[0]+add)>>shift)).(int))

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
    add := int(1) << (shift - 1)

    for j = 0; j < line; j++ {
        /* E and O*/
        for k = 0; k < 16; k++ {
            E[k] = int(src[k+j*32]) + int(src[31-k+j*32])
            O[k] = int(src[k+j*32]) - int(src[31-k+j*32])
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
