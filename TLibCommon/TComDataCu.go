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
    "container/list"
)

// ====================================================================================================================
// Non-deblocking in-loop filter processing block data structure
// ====================================================================================================================

/// Non-deblocking filter processing block border tag
type NDBFBlockBorderTag uint8

const ( //enum NDBFBlockBorderTag
    SGU_L = iota
    SGU_R
    SGU_T
    SGU_B
    SGU_TL
    SGU_TR
    SGU_BL
    SGU_BR
    NUM_SGU_BORDER
)

/// Non-deblocking filter processing block information
type NDBFBlockInfo struct {
    tileID              int                  //!< tile ID
    sliceID             int                  //!< slice ID
    startSU             uint                 //!< starting SU z-scan address in LCU
    endSU               uint                 //!< ending SU z-scan address in LCU
    widthSU             uint                 //!< number of SUs in width
    heightSU            uint                 //!< number of SUs in height
    posX                uint                 //!< top-left X coordinate in picture
    posY                uint                 //!< top-left Y coordinate in picture
    width               uint                 //!< number of pixels in width
    height              uint                 //!< number of pixels in height
    isBorderAvailable   [NUM_SGU_BORDER]bool //!< the border availabilities
    allBordersAvailable bool

    //NDBFBlockInfo():tileID(0), sliceID(0), startSU(0), endSU(0) {} //!< constructor
    //const NDBFBlockInfo& operator= (const NDBFBlockInfo& src);  //!< "=" operator
}

// ====================================================================================================================
// Class definition
// ====================================================================================================================

/// CU data structure class
type TComDataCU struct {
    //private:

    // -------------------------------------------------------------------------------------------------------------------
    // class pointers
    // -------------------------------------------------------------------------------------------------------------------

    m_pcPic     *TComPic     ///< picture class pointer
    m_pcSlice   *TComSlice   ///< slice header pointer
    m_pcPattern *TComPattern ///< neighbour access class pointer

    // -------------------------------------------------------------------------------------------------------------------
    // CU description
    // -------------------------------------------------------------------------------------------------------------------

    m_uiCUAddr       uint   ///< CU address in a slice
    m_uiAbsIdxInLCU  uint   ///< absolute address in a CU. It's Z scan order
    m_uiCUPelX       uint   ///< CU position in a pixel (X)
    m_uiCUPelY       uint   ///< CU position in a pixel (Y)
    m_uiNumPartition uint   ///< total number of minimum partitions in a CU
    m_puhWidth       []byte ///< array of widths
    m_puhHeight      []byte ///< array of heights
    m_puhDepth       []byte ///< array of depths
    m_unitSize       int    ///< size of a "minimum partition"

    // -------------------------------------------------------------------------------------------------------------------
    // CU data
    // -------------------------------------------------------------------------------------------------------------------
    m_skipFlag           []bool           ///< array of skip flags
    m_pePartSize         []PartSize       ///< array of partition sizes
    m_pePredMode         []PredMode       ///< array of prediction modes
    m_CUTransquantBypass []bool           ///< array of cu_transquant_bypass flags
    m_phQP               []int8           ///< array of QP values
    m_puhTrIdx           []byte           ///< array of transform indices
    m_puhTransformSkip   [3][]bool        ///< array of transform skipping flags
    m_puhCbf             [3][]byte        ///< array of coded block flags (CBF)
    m_acCUMvField        [2]TComCUMvField ///< array of motion vectors
    m_pcTrCoeffY         []TCoeff         ///< transformed coefficient buffer (Y)
    m_pcTrCoeffCb        []TCoeff         ///< transformed coefficient buffer (Cb)
    m_pcTrCoeffCr        []TCoeff         ///< transformed coefficient buffer (Cr)
    //#if ADAPTIVE_QP_SELECTION
    m_pcArlCoeffY                 []TCoeff ///< ARL coefficient buffer (Y)
    m_pcArlCoeffCb                []TCoeff ///< ARL coefficient buffer (Cb)
    m_pcArlCoeffCr                []TCoeff ///< ARL coefficient buffer (Cr)
    m_ArlCoeffIsAliasedAllocation bool     ///< ARL coefficient buffer is an alias of the global buffer and must not be free()'d

    m_pcGlbArlCoeffY  []TCoeff ///< ARL coefficient buffer (Y)
    m_pcGlbArlCoeffCb []TCoeff ///< ARL coefficient buffer (Cb)
    m_pcGlbArlCoeffCr []TCoeff ///< ARL coefficient buffer (Cr)
    //#endif

    m_pcIPCMSampleY  []Pel ///< PCM sample buffer (Y)
    m_pcIPCMSampleCb []Pel ///< PCM sample buffer (Cb)
    m_pcIPCMSampleCr []Pel ///< PCM sample buffer (Cr)

    m_piSliceSUMap    []int ///< pointer of slice ID map
    m_iSliceSUMapAddr int
    m_vNDFBlock       *list.List

    // -------------------------------------------------------------------------------------------------------------------
    // neighbour access variables
    // -------------------------------------------------------------------------------------------------------------------

    m_pcCUAboveLeft  *TComDataCU    ///< pointer of above-left CU
    m_pcCUAboveRight *TComDataCU    ///< pointer of above-right CU
    m_pcCUAbove      *TComDataCU    ///< pointer of above CU
    m_pcCULeft       *TComDataCU    ///< pointer of left CU
    m_apcCUColocated [2]*TComDataCU ///< pointer of temporally colocated CU's for both directions
    m_cMvFieldA      TComMvField    ///< motion vector of position A
    m_cMvFieldB      TComMvField    ///< motion vector of position B
    m_cMvFieldC      TComMvField    ///< motion vector of position C
    m_cMvPred        TComMv         ///< motion vector predictor

    // -------------------------------------------------------------------------------------------------------------------
    // coding tool information
    // -------------------------------------------------------------------------------------------------------------------

    m_pbMergeFlag   []bool ///< array of merge flags
    m_puhMergeIndex []byte ///< array of merge candidate indices
    //#if AMP_MRG
    m_bIsMergeAMP bool
    //#endif
    m_puhLumaIntraDir   []byte    ///< array of intra directions (luma)
    m_puhChromaIntraDir []byte    ///< array of intra directions (chroma)
    m_puhInterDir       []byte    ///< array of inter directions
    m_apiMVPIdx         [2][]int8 ///< array of motion vector predictor candidates
    m_apiMVPNum         [2][]int8 ///< array of number of possible motion vectors predictors
    m_pbIPCMFlag        []bool    ///< array of intra_pcm flags

    // -------------------------------------------------------------------------------------------------------------------
    // misc. variables
    // -------------------------------------------------------------------------------------------------------------------

    m_bDecSubCu           bool    ///< indicates decoder-mode
    m_dTotalCost          float64 ///< sum of partition RD costs
    m_uiTotalDistortion   uint    ///< sum of partition distortion
    m_uiTotalBits         uint    ///< sum of partition bits
    m_uiTotalBins         uint    ///< sum of partition bins
    m_sliceStartCU        []uint  ///< Start CU address of current slice
    m_sliceSegmentStartCU []uint  ///< Start CU address of current slice
    m_codedQP             int8
}

//public:

func NewTComDataCU() *TComDataCU {
    return &TComDataCU{}
}

// -------------------------------------------------------------------------------------------------------------------
// create / destroy / initialize / copy
// -------------------------------------------------------------------------------------------------------------------

func (this *TComDataCU) Create(uiNumPartition, uiWidth, uiHeight uint, bDecSubCu bool, unitSize int,
    //#if ADAPTIVE_QP_SELECTION
    bGlobalRMARLBuffer bool) {
    //#endif
    this.m_bDecSubCu = bDecSubCu

    this.m_pcPic = nil
    this.m_pcSlice = nil
    this.m_uiNumPartition = uiNumPartition
    this.m_unitSize = unitSize

    if !bDecSubCu {
        this.m_phQP = make([]int8, uiNumPartition)
        this.m_puhDepth = make([]byte, uiNumPartition)
        this.m_puhWidth = make([]byte, uiNumPartition)
        this.m_puhHeight = make([]byte, uiNumPartition)

        this.m_skipFlag = make([]bool, uiNumPartition)

        this.m_pePartSize = make([]PartSize, uiNumPartition)
        for i := uint(0); i < uiNumPartition; i++ {
            this.m_pePartSize[i] = SIZE_NONE
        }

        this.m_pePredMode = make([]PredMode, uiNumPartition)
        this.m_CUTransquantBypass = make([]bool, uiNumPartition)
        this.m_pbMergeFlag = make([]bool, uiNumPartition)
        this.m_puhMergeIndex = make([]byte, uiNumPartition)
        this.m_puhLumaIntraDir = make([]byte, uiNumPartition)
        this.m_puhChromaIntraDir = make([]byte, uiNumPartition)
        this.m_puhInterDir = make([]byte, uiNumPartition)

        this.m_puhTrIdx = make([]byte, uiNumPartition)
        this.m_puhTransformSkip[0] = make([]bool, uiNumPartition)
        this.m_puhTransformSkip[1] = make([]bool, uiNumPartition)
        this.m_puhTransformSkip[2] = make([]bool, uiNumPartition)

        this.m_puhCbf[0] = make([]byte, uiNumPartition)
        this.m_puhCbf[1] = make([]byte, uiNumPartition)
        this.m_puhCbf[2] = make([]byte, uiNumPartition)

        this.m_apiMVPIdx[0] = make([]int8, uiNumPartition)
        this.m_apiMVPIdx[1] = make([]int8, uiNumPartition)
        this.m_apiMVPNum[0] = make([]int8, uiNumPartition)
        this.m_apiMVPNum[1] = make([]int8, uiNumPartition)
        for i := uint(0); i < uiNumPartition; i++ {
            this.m_apiMVPIdx[0][i] = -1
            this.m_apiMVPIdx[1][i] = -1
        }
        this.m_pcTrCoeffY = make([]TCoeff, uiWidth*uiHeight)
        this.m_pcTrCoeffCb = make([]TCoeff, uiWidth*uiHeight/4)
        this.m_pcTrCoeffCr = make([]TCoeff, uiWidth*uiHeight/4)

        //#if ADAPTIVE_QP_SELECTION
        if bGlobalRMARLBuffer {
            if this.m_pcGlbArlCoeffY == nil {
                this.m_pcGlbArlCoeffY = make([]TCoeff, uiWidth*uiHeight)
                this.m_pcGlbArlCoeffCb = make([]TCoeff, uiWidth*uiHeight/4)
                this.m_pcGlbArlCoeffCr = make([]TCoeff, uiWidth*uiHeight/4)
            }
            this.m_pcArlCoeffY = this.m_pcGlbArlCoeffY
            this.m_pcArlCoeffCb = this.m_pcGlbArlCoeffCb
            this.m_pcArlCoeffCr = this.m_pcGlbArlCoeffCr
            this.m_ArlCoeffIsAliasedAllocation = true
        } else {
            this.m_pcArlCoeffY = make([]TCoeff, uiWidth*uiHeight)
            this.m_pcArlCoeffCb = make([]TCoeff, uiWidth*uiHeight/4)
            this.m_pcArlCoeffCr = make([]TCoeff, uiWidth*uiHeight/4)
        }
        //#endif

        this.m_pbIPCMFlag = make([]bool, uiNumPartition)
        this.m_pcIPCMSampleY = make([]Pel, uiWidth*uiHeight)
        this.m_pcIPCMSampleCb = make([]Pel, uiWidth*uiHeight/4)
        this.m_pcIPCMSampleCr = make([]Pel, uiWidth*uiHeight/4)

        this.m_acCUMvField[0].Create(uiNumPartition)
        this.m_acCUMvField[1].Create(uiNumPartition)

    } else {
        this.m_acCUMvField[0].SetNumPartition(uiNumPartition)
        this.m_acCUMvField[1].SetNumPartition(uiNumPartition)
    }

    this.m_sliceStartCU = make([]uint, uiNumPartition)
    this.m_sliceSegmentStartCU = make([]uint, uiNumPartition)

    // create pattern memory
    this.m_pcPattern = NewTComPattern()

    // create motion vector fields

    this.m_pcCUAboveLeft = nil
    this.m_pcCUAboveRight = nil
    this.m_pcCUAbove = nil
    this.m_pcCULeft = nil

    this.m_apcCUColocated[0] = nil
    this.m_apcCUColocated[1] = nil

    this.m_vNDFBlock = list.New()
}
func (this *TComDataCU) Destroy() {
    this.m_pcPic = nil
    this.m_pcSlice = nil

    if this.m_pcPattern != nil {
        this.m_pcPattern = nil
    }

    // encoder-side buffer free
    if !this.m_bDecSubCu {
        this.m_phQP = nil
        this.m_puhDepth = nil
        this.m_puhWidth = nil
        this.m_puhHeight = nil

        this.m_skipFlag = nil

        this.m_pePartSize = nil
        this.m_pePredMode = nil
        this.m_CUTransquantBypass = nil
        this.m_puhCbf[0] = nil
        this.m_puhCbf[1] = nil
        this.m_puhCbf[2] = nil
        this.m_puhInterDir = nil
        this.m_pbMergeFlag = nil
        this.m_puhMergeIndex = nil
        this.m_puhLumaIntraDir = nil
        this.m_puhChromaIntraDir = nil
        this.m_puhTrIdx = nil
        this.m_puhTransformSkip[0] = nil
        this.m_puhTransformSkip[1] = nil
        this.m_puhTransformSkip[2] = nil
        this.m_pcTrCoeffY = nil
        this.m_pcTrCoeffCb = nil
        this.m_pcTrCoeffCr = nil
        //#if ADAPTIVE_QP_SELECTION
        if !this.m_ArlCoeffIsAliasedAllocation {
            this.m_pcArlCoeffY = nil
            this.m_pcArlCoeffCb = nil
            this.m_pcArlCoeffCr = nil
        }
        this.m_pcGlbArlCoeffY = nil
        this.m_pcGlbArlCoeffCb = nil
        this.m_pcGlbArlCoeffCr = nil
        //#endi
        this.m_pbIPCMFlag = nil
        this.m_pcIPCMSampleY = nil
        this.m_pcIPCMSampleCb = nil
        this.m_pcIPCMSampleCr = nil
        this.m_apiMVPIdx[0] = nil
        this.m_apiMVPIdx[1] = nil
        this.m_apiMVPNum[0] = nil
        this.m_apiMVPNum[1] = nil

        this.m_acCUMvField[0].Destroy()
        this.m_acCUMvField[1].Destroy()
    }

    this.m_pcCUAboveLeft = nil
    this.m_pcCUAboveRight = nil
    this.m_pcCUAbove = nil
    this.m_pcCULeft = nil

    this.m_apcCUColocated[0] = nil
    this.m_apcCUColocated[1] = nil

    this.m_sliceStartCU = nil
    this.m_sliceSegmentStartCU = nil
}

func (this *TComDataCU) InitCU(pcPic *TComPic, iCUAddr uint) {
    var i int

    this.m_pcPic = pcPic
    this.m_pcSlice = pcPic.GetSlice(pcPic.GetCurrSliceIdx())
    this.m_uiCUAddr = iCUAddr
    this.m_uiCUPelX = (iCUAddr % pcPic.GetFrameWidthInCU()) * this.GetSlice().GetSPS().GetMaxCUWidth()
    this.m_uiCUPelY = (iCUAddr / pcPic.GetFrameWidthInCU()) * this.GetSlice().GetSPS().GetMaxCUHeight()
    this.m_uiAbsIdxInLCU = 0
    this.m_dTotalCost = MAX_DOUBLE
    this.m_uiTotalDistortion = 0
    this.m_uiTotalBits = 0
    this.m_uiTotalBins = 0
    this.m_uiNumPartition = pcPic.GetNumPartInCU()

    for i = 0; i < int(pcPic.GetNumPartInCU()); i++ {
        if pcPic.GetPicSym().GetInverseCUOrderMap(int(iCUAddr))*pcPic.GetNumPartInCU()+uint(i) >= this.GetSlice().GetSliceCurStartCUAddr() {
            this.m_sliceStartCU[i] = this.GetSlice().GetSliceCurStartCUAddr()
        } else {
            this.m_sliceStartCU[i] = pcPic.GetCU(this.GetAddr()).m_sliceStartCU[i]
        }
    }
    for i = 0; i < int(pcPic.GetNumPartInCU()); i++ {
        if pcPic.GetPicSym().GetInverseCUOrderMap(int(iCUAddr))*pcPic.GetNumPartInCU()+uint(i) >= this.GetSlice().GetSliceSegmentCurStartCUAddr() {
            this.m_sliceSegmentStartCU[i] = this.GetSlice().GetSliceSegmentCurStartCUAddr()
        } else {
            this.m_sliceSegmentStartCU[i] = pcPic.GetCU(this.GetAddr()).m_sliceSegmentStartCU[i]
        }
    }

    partStartIdx := int(this.GetSlice().GetSliceSegmentCurStartCUAddr()) - int(pcPic.GetPicSym().GetInverseCUOrderMap(int(iCUAddr))*pcPic.GetNumPartInCU())
    //fmt.Printf("partStartIdx:%d-%d*%d\n", this.GetSlice().GetSliceSegmentCurStartCUAddr(),pcPic.GetPicSym().GetInverseCUOrderMap(int(iCUAddr)), pcPic.GetNumPartInCU());

    var ui uint
    var numElements int
    if partStartIdx < int(this.m_uiNumPartition) {
        numElements = partStartIdx
    } else {
        numElements = int(this.m_uiNumPartition)
    }

    for ui = 0; int(ui) < numElements; ui++ {
        pcFrom := pcPic.GetCU(this.GetAddr())
        this.m_skipFlag[ui] = pcFrom.GetSkipFlag1(ui)
        this.m_pePartSize[ui] = pcFrom.GetPartitionSize1(ui)
        this.m_pePredMode[ui] = pcFrom.GetPredictionMode1(ui)
        this.m_CUTransquantBypass[ui] = pcFrom.GetCUTransquantBypass1(ui)
        this.m_puhDepth[ui] = pcFrom.GetDepth1(ui)
        this.m_puhWidth[ui] = pcFrom.GetWidth1(ui)
        this.m_puhHeight[ui] = pcFrom.GetHeight1(ui)
        this.m_puhTrIdx[ui] = pcFrom.GetTransformIdx1(ui)
        this.m_puhTransformSkip[0][ui] = pcFrom.GetTransformSkip2(ui, TEXT_LUMA)
        this.m_puhTransformSkip[1][ui] = pcFrom.GetTransformSkip2(ui, TEXT_CHROMA_U)
        this.m_puhTransformSkip[2][ui] = pcFrom.GetTransformSkip2(ui, TEXT_CHROMA_V)
        this.m_apiMVPIdx[0][ui] = pcFrom.m_apiMVPIdx[0][ui]

        this.m_apiMVPIdx[1][ui] = pcFrom.m_apiMVPIdx[1][ui]
        this.m_apiMVPNum[0][ui] = pcFrom.m_apiMVPNum[0][ui]
        this.m_apiMVPNum[1][ui] = pcFrom.m_apiMVPNum[1][ui]
        this.m_phQP[ui] = pcFrom.m_phQP[ui]
        this.m_pbMergeFlag[ui] = pcFrom.m_pbMergeFlag[ui]
        this.m_puhMergeIndex[ui] = pcFrom.m_puhMergeIndex[ui]
        this.m_puhLumaIntraDir[ui] = pcFrom.m_puhLumaIntraDir[ui]
        this.m_puhChromaIntraDir[ui] = pcFrom.m_puhChromaIntraDir[ui]
        this.m_puhInterDir[ui] = pcFrom.m_puhInterDir[ui]
        this.m_puhCbf[0][ui] = pcFrom.m_puhCbf[0][ui]
        this.m_puhCbf[1][ui] = pcFrom.m_puhCbf[1][ui]
        this.m_puhCbf[2][ui] = pcFrom.m_puhCbf[2][ui]
        this.m_pbIPCMFlag[ui] = pcFrom.m_pbIPCMFlag[ui]
    }

    var firstElement int
    if partStartIdx > 0 {
        firstElement = partStartIdx
    } else {
        firstElement = 0
    }
    numElements = int(this.m_uiNumPartition) - firstElement

    if numElements > 0 {
        for i := int(0); i < numElements; i++ {
            this.m_skipFlag[firstElement+i] = false
            this.m_pePartSize[firstElement+i] = SIZE_NONE
            this.m_pePredMode[firstElement+i] = MODE_NONE
            this.m_CUTransquantBypass[firstElement+i] = false
            this.m_puhDepth[firstElement+i] = 0
            this.m_puhTrIdx[firstElement+i] = 0
            this.m_puhTransformSkip[0][firstElement+i] = false
            this.m_puhTransformSkip[1][firstElement+i] = false
            this.m_puhTransformSkip[2][firstElement+i] = false
            this.m_puhWidth[firstElement+i] = byte(this.GetSlice().GetSPS().GetMaxCUWidth())
            this.m_puhHeight[firstElement+i] = byte(this.GetSlice().GetSPS().GetMaxCUHeight())
            this.m_apiMVPIdx[0][firstElement+i] = -1
            this.m_apiMVPIdx[1][firstElement+i] = -1
            this.m_apiMVPNum[0][firstElement+i] = -1
            this.m_apiMVPNum[1][firstElement+i] = -1
            this.m_phQP[firstElement+i] = int8(this.GetSlice().GetSliceQp())
            this.m_pbMergeFlag[firstElement+i] = false
            this.m_puhMergeIndex[firstElement+i] = 0
            this.m_puhLumaIntraDir[firstElement+i] = DC_IDX
            this.m_puhChromaIntraDir[firstElement+i] = 0
            this.m_puhInterDir[firstElement+i] = 0
            this.m_puhCbf[0][firstElement+i] = 0
            this.m_puhCbf[1][firstElement+i] = 0
            this.m_puhCbf[2][firstElement+i] = 0
            this.m_pbIPCMFlag[firstElement+i] = false
        }
    }

    uiTmp := this.GetSlice().GetSPS().GetMaxCUWidth() * this.GetSlice().GetSPS().GetMaxCUHeight()
    if 0 >= partStartIdx {
        this.m_acCUMvField[0].ClearMvField()
        this.m_acCUMvField[1].ClearMvField()
        //memSet( this.m_pcTrCoeffY , 0, sizeof( TCoeff ) * uiTmp );
        //#if ADAPTIVE_QP_SELECTION
        //memSet( this.m_pcArlCoeffY , 0, sizeof( Int ) * uiTmp );
        //#endif
        //memSet( this.m_pcIPCMSampleY , 0, sizeof( Pel ) * uiTmp );
        uiTmp >>= 2
        //memSet( this.m_pcTrCoeffCb, 0, sizeof( TCoeff ) * uiTmp );
        //memSet( this.m_pcTrCoeffCr, 0, sizeof( TCoeff ) * uiTmp );
        //#if ADAPTIVE_QP_SELECTION
        //memSet( this.m_pcArlCoeffCb, 0, sizeof( Int ) * uiTmp );
        //memSet( this.m_pcArlCoeffCr, 0, sizeof( Int ) * uiTmp );
        //#endif
        //memSet( this.m_pcIPCMSampleCb , 0, sizeof( Pel ) * uiTmp );
        //memSet( this.m_pcIPCMSampleCr , 0, sizeof( Pel ) * uiTmp );
    } else {
        pcFrom := pcPic.GetCU(this.GetAddr())
        this.m_acCUMvField[0].CopyFrom(&pcFrom.m_acCUMvField[0], this.m_uiNumPartition, 0)
        this.m_acCUMvField[1].CopyFrom(&pcFrom.m_acCUMvField[1], this.m_uiNumPartition, 0)
        for i := uint(0); i < uiTmp; i++ {
            this.m_pcTrCoeffY[i] = pcFrom.m_pcTrCoeffY[i]
            //#if ADAPTIVE_QP_SELECTION
            this.m_pcArlCoeffY[i] = pcFrom.m_pcArlCoeffY[i]
            //#endif
            this.m_pcIPCMSampleY[i] = pcFrom.m_pcIPCMSampleY[i]
        }
        for i := uint(0); i < (uiTmp >> 2); i++ {
            this.m_pcTrCoeffCb[i] = pcFrom.m_pcTrCoeffCb[i]
            this.m_pcTrCoeffCr[i] = pcFrom.m_pcTrCoeffCr[i]
            //#if ADAPTIVE_QP_SELECTION
            this.m_pcArlCoeffCb[i] = pcFrom.m_pcArlCoeffCb[i]
            this.m_pcArlCoeffCr[i] = pcFrom.m_pcArlCoeffCr[i]
            //#endif
            this.m_pcIPCMSampleCb[i] = pcFrom.m_pcIPCMSampleCb[i]
            this.m_pcIPCMSampleCr[i] = pcFrom.m_pcIPCMSampleCr[i]
        }
    }

    // Setting neighbor CU
    this.m_pcCULeft = nil
    this.m_pcCUAbove = nil
    this.m_pcCUAboveLeft = nil
    this.m_pcCUAboveRight = nil

    this.m_apcCUColocated[0] = nil
    this.m_apcCUColocated[1] = nil

    uiWidthInCU := pcPic.GetFrameWidthInCU()
    if this.m_uiCUAddr%uiWidthInCU != 0 {
        this.m_pcCULeft = pcPic.GetCU(this.m_uiCUAddr - 1)
    }

    if this.m_uiCUAddr/uiWidthInCU != 0 {
        this.m_pcCUAbove = pcPic.GetCU(this.m_uiCUAddr - uiWidthInCU)
    }

    if this.m_pcCULeft != nil && this.m_pcCUAbove != nil {
        this.m_pcCUAboveLeft = pcPic.GetCU(this.m_uiCUAddr - uiWidthInCU - 1)
    }

    if this.m_pcCUAbove != nil && ((this.m_uiCUAddr % uiWidthInCU) < (uiWidthInCU - 1)) {
        this.m_pcCUAboveRight = pcPic.GetCU(this.m_uiCUAddr - uiWidthInCU + 1)
    }

    if this.GetSlice().GetNumRefIdx(REF_PIC_LIST_0) > 0 {
        this.m_apcCUColocated[0] = this.GetSlice().GetRefPic(REF_PIC_LIST_0, 0).GetCU(this.m_uiCUAddr)
    }

    if this.GetSlice().GetNumRefIdx(REF_PIC_LIST_1) > 0 {
        this.m_apcCUColocated[1] = this.GetSlice().GetRefPic(REF_PIC_LIST_1, 0).GetCU(this.m_uiCUAddr)
    }
}

//protected:

/// add possible motion vector predictor candidates
func (this *TComDataCU) xAddMVPCand(pInfo *AMVPInfo, eRefPicList RefPicList, iRefIdx int, uiPartUnitIdx uint, eDir MVP_DIR) bool {
    var pcTmpCU *TComDataCU
    var uiIdx uint
    switch eDir {
    case MD_LEFT:
        pcTmpCU = this.GetPULeft(&uiIdx, uiPartUnitIdx, true, true)
    case MD_ABOVE:
        pcTmpCU = this.GetPUAbove(&uiIdx, uiPartUnitIdx, true, false, true)
    case MD_ABOVE_RIGHT:
        pcTmpCU = this.GetPUAboveRight(&uiIdx, uiPartUnitIdx, true)
    case MD_BELOW_LEFT:
        pcTmpCU = this.GetPUBelowLeft(&uiIdx, uiPartUnitIdx, true)
    case MD_ABOVE_LEFT:
        pcTmpCU = this.GetPUAboveLeft(&uiIdx, uiPartUnitIdx, true)
    default:
    }

    if pcTmpCU == nil {
        return false
    }

    if pcTmpCU.GetCUMvField(eRefPicList).GetRefIdx(int(uiIdx)) >= 0 && int(this.m_pcSlice.GetRefPic(eRefPicList, iRefIdx).GetPOC()) == pcTmpCU.GetSlice().GetRefPOC(eRefPicList, int(pcTmpCU.GetCUMvField(eRefPicList).GetRefIdx(int(uiIdx)))) {
        cMvPred := pcTmpCU.GetCUMvField(eRefPicList).GetMv(int(uiIdx))
        //fmt.Printf("xAddMVPCand1:(%d,%d) ", cMvPred.GetHor(), cMvPred.GetVer());
        pInfo.MvCand[pInfo.IN] = cMvPred
        pInfo.IN++
        return true
    }

    eRefPicList2nd := REF_PIC_LIST_0
    if eRefPicList == REF_PIC_LIST_0 {
        eRefPicList2nd = REF_PIC_LIST_1
    } else if eRefPicList == REF_PIC_LIST_1 {
        eRefPicList2nd = REF_PIC_LIST_0
    }

    iCurrRefPOC := this.m_pcSlice.GetRefPic(eRefPicList, iRefIdx).GetPOC()
    var iNeibRefPOC int

    if pcTmpCU.GetCUMvField(RefPicList(eRefPicList2nd)).GetRefIdx(int(uiIdx)) >= 0 {
        iNeibRefPOC = pcTmpCU.GetSlice().GetRefPOC(RefPicList(eRefPicList2nd), int(pcTmpCU.GetCUMvField(RefPicList(eRefPicList2nd)).GetRefIdx(int(uiIdx))))
        if iNeibRefPOC == int(iCurrRefPOC) { // Same Reference Frame But Diff List//
            cMvPred := pcTmpCU.GetCUMvField(RefPicList(eRefPicList2nd)).GetMv(int(uiIdx))
            //fmt.Printf("xAddMVPCand2:(%d,%d) ", cMvPred.GetHor(), cMvPred.GetVer());
            pInfo.MvCand[pInfo.IN] = cMvPred
            pInfo.IN++
            return true
        }
    }
    return false
}
func (this *TComDataCU) xAddMVPCandOrder(pInfo *AMVPInfo, eRefPicList RefPicList, iRefIdx int, uiPartUnitIdx uint, eDir MVP_DIR) bool {
    var pcTmpCU *TComDataCU
    var uiIdx uint
    switch eDir {
    case MD_LEFT:
        pcTmpCU = this.GetPULeft(&uiIdx, uiPartUnitIdx, true, true)
    case MD_ABOVE:
        pcTmpCU = this.GetPUAbove(&uiIdx, uiPartUnitIdx, true, false, true)
    case MD_ABOVE_RIGHT:
        pcTmpCU = this.GetPUAboveRight(&uiIdx, uiPartUnitIdx, true)
    case MD_BELOW_LEFT:
        pcTmpCU = this.GetPUBelowLeft(&uiIdx, uiPartUnitIdx, true)
    case MD_ABOVE_LEFT:
        pcTmpCU = this.GetPUAboveLeft(&uiIdx, uiPartUnitIdx, true)
    default:
    }

    if pcTmpCU == nil {
        return false
    }

    eRefPicList2nd := REF_PIC_LIST_0
    if eRefPicList == REF_PIC_LIST_0 {
        eRefPicList2nd = REF_PIC_LIST_1
    } else if eRefPicList == REF_PIC_LIST_1 {
        eRefPicList2nd = REF_PIC_LIST_0
    }

    iCurrPOC := this.m_pcSlice.GetPOC()
    iCurrRefPOC := this.m_pcSlice.GetRefPic(eRefPicList, iRefIdx).GetPOC()
    iNeibPOC := iCurrPOC
    var iNeibRefPOC int

    bIsCurrRefLongTerm := this.m_pcSlice.GetRefPic(eRefPicList, iRefIdx).GetIsLongTerm()
    bIsNeibRefLongTerm := false
    //---------------  V1 (END) ------------------//
    if pcTmpCU.GetCUMvField(eRefPicList).GetRefIdx(int(uiIdx)) >= 0 {
        iNeibRefPOC = pcTmpCU.GetSlice().GetRefPOC(eRefPicList, int(pcTmpCU.GetCUMvField(eRefPicList).GetRefIdx(int(uiIdx))))
        cMvPred := pcTmpCU.GetCUMvField(eRefPicList).GetMv(int(uiIdx))
        var rcMv TComMv

        bIsNeibRefLongTerm = pcTmpCU.GetSlice().GetRefPic(eRefPicList, int(pcTmpCU.GetCUMvField(eRefPicList).GetRefIdx(int(uiIdx)))).GetIsLongTerm()
        if bIsCurrRefLongTerm == bIsNeibRefLongTerm {
            if bIsCurrRefLongTerm || bIsNeibRefLongTerm {
                rcMv = cMvPred
            } else {
                iScale := this.xGetDistScaleFactor(iCurrPOC, int(iCurrRefPOC), iNeibPOC, iNeibRefPOC)
                if iScale == 4096 {
                    rcMv = cMvPred
                } else {
                    //fmt.Printf("xAddMVPCandOrder1:(%d,%d,%d,%d)=%d",iCurrPOC, int(iCurrRefPOC), iNeibPOC, iNeibRefPOC,iScale);
                    rcMv = cMvPred.ScaleMv(iScale)
                }
            }
            //fmt.Printf("xAddMVPCandOrder1:(%d,%d) ", rcMv.GetHor(), rcMv.GetVer());
            pInfo.MvCand[pInfo.IN] = rcMv
            pInfo.IN++
            return true
        }
    }
    //---------------------- V2(END) --------------------//
    if pcTmpCU.GetCUMvField(RefPicList(eRefPicList2nd)).GetRefIdx(int(uiIdx)) >= 0 {
        iNeibRefPOC = pcTmpCU.GetSlice().GetRefPOC(RefPicList(eRefPicList2nd), int(pcTmpCU.GetCUMvField(RefPicList(eRefPicList2nd)).GetRefIdx(int(uiIdx))))
        cMvPred := pcTmpCU.GetCUMvField(RefPicList(eRefPicList2nd)).GetMv(int(uiIdx))
        var rcMv TComMv

        bIsNeibRefLongTerm = pcTmpCU.GetSlice().GetRefPic(RefPicList(eRefPicList2nd), int(pcTmpCU.GetCUMvField(RefPicList(eRefPicList2nd)).GetRefIdx(int(uiIdx)))).GetIsLongTerm()
        if bIsCurrRefLongTerm == bIsNeibRefLongTerm {
            if bIsCurrRefLongTerm || bIsNeibRefLongTerm {
                rcMv = cMvPred
            } else {
                iScale := this.xGetDistScaleFactor(iCurrPOC, int(iCurrRefPOC), iNeibPOC, iNeibRefPOC)
                if iScale == 4096 {
                    rcMv = cMvPred
                } else {
                    //fmt.Printf("xAddMVPCandOrder2:%d",iScale);
                    rcMv = cMvPred.ScaleMv(iScale)
                }
            }
            //fmt.Printf("xAddMVPCandOrder2:(%d,%d) ", rcMv.GetHor(), rcMv.GetVer());
            pInfo.MvCand[pInfo.IN] = rcMv
            pInfo.IN++
            return true
        }
    }
    //---------------------- V3(END) --------------------//
    return false
}

func (this *TComDataCU) DeriveRightBottomIdx(uiPartIdx uint, ruiPartIdxRB *uint) {
    *ruiPartIdxRB = G_auiRasterToZscan[G_auiZscanToRaster[this.m_uiAbsIdxInLCU]+(((uint(this.m_puhHeight[0])/this.m_pcPic.GetMinCUHeight())>>1)-1)*this.m_pcPic.GetNumPartInWidth()+uint(this.m_puhWidth[0])/this.m_pcPic.GetMinCUWidth()-1]

    switch this.m_pePartSize[0] {
    case SIZE_2Nx2N:
        *ruiPartIdxRB += this.m_uiNumPartition >> 1
        //break;
    case SIZE_2NxN:
        if uiPartIdx == 0 {
            *ruiPartIdxRB += 0
        } else {
            *ruiPartIdxRB += this.m_uiNumPartition >> 1
        }
        //break;
    case SIZE_Nx2N:
        if uiPartIdx == 0 {
            *ruiPartIdxRB += this.m_uiNumPartition >> 2
        } else {
            *ruiPartIdxRB += this.m_uiNumPartition >> 1
        }
        //break;
    case SIZE_NxN:
        *ruiPartIdxRB += (this.m_uiNumPartition >> 2) * (uiPartIdx - 1)
        //break;
    case SIZE_2NxnU:
        if uiPartIdx == 0 {
            *ruiPartIdxRB -= this.m_uiNumPartition >> 3
        } else {
            *ruiPartIdxRB += this.m_uiNumPartition >> 1
        }
        //break;
    case SIZE_2NxnD:
        if uiPartIdx == 0 {
            *ruiPartIdxRB += (this.m_uiNumPartition >> 2) + (this.m_uiNumPartition >> 3)
        } else {
            *ruiPartIdxRB += this.m_uiNumPartition >> 1
        }
        //break;
    case SIZE_nLx2N:
        if uiPartIdx == 0 {
            *ruiPartIdxRB += (this.m_uiNumPartition >> 3) + (this.m_uiNumPartition >> 4)
        } else {
            *ruiPartIdxRB += this.m_uiNumPartition >> 1
        }
        //break;
    case SIZE_nRx2N:
        if uiPartIdx == 0 {
            *ruiPartIdxRB += (this.m_uiNumPartition >> 2) + (this.m_uiNumPartition >> 3) + (this.m_uiNumPartition >> 4)
        } else {
            *ruiPartIdxRB += this.m_uiNumPartition >> 1
        }
        //break;
    default:
        //assert (0);
        //break;
    }
}
func (this *TComDataCU) xGetColMVP(eRefPicList RefPicList, uiCUAddr, uiPartUnitIdx int, rcMv *TComMv, riRefIdx *int) bool {
    uiAbsPartAddr := uiPartUnitIdx

    var eColRefPicList RefPicList
    var iColPOC, iColRefPOC, iCurrPOC, iCurrRefPOC, iScale int
    var cColMv TComMv

    // use coldir.
    var pColPic *TComPic
    if this.GetSlice().IsInterB() {
        pColPic = this.GetSlice().GetRefPic(RefPicList(1-this.GetSlice().GetColFromL0Flag()), int(this.GetSlice().GetColRefIdx()))
    } else {
        pColPic = this.GetSlice().GetRefPic(RefPicList(0), int(this.GetSlice().GetColRefIdx()))
    }

    pColCU := pColPic.GetCU(uint(uiCUAddr))
    if pColCU.GetPic() == nil || pColCU.GetPartitionSize1(uint(uiPartUnitIdx)) == SIZE_NONE {
        return false
    }
    iCurrPOC = this.m_pcSlice.GetPOC()
    iCurrRefPOC = int(this.m_pcSlice.GetRefPic(eRefPicList, *riRefIdx).GetPOC())
    iColPOC = pColCU.GetSlice().GetPOC()

    if pColCU.IsIntra(uint(uiAbsPartAddr)) {
        return false
    }
    if this.GetSlice().GetCheckLDC() {
        eColRefPicList = eRefPicList
    } else {
        eColRefPicList = RefPicList(this.GetSlice().GetColFromL0Flag())
    }

    iColRefIdx := pColCU.GetCUMvField(RefPicList(eColRefPicList)).GetRefIdx(uiAbsPartAddr)

    if iColRefIdx < 0 {
        eColRefPicList = RefPicList(1 - eColRefPicList)
        iColRefIdx = pColCU.GetCUMvField(RefPicList(eColRefPicList)).GetRefIdx(uiAbsPartAddr)

        if iColRefIdx < 0 {
            return false
        }
    }

    // Scale the vector.
    iColRefPOC = pColCU.GetSlice().GetRefPOC(eColRefPicList, int(iColRefIdx))
    cColMv = pColCU.GetCUMvField(eColRefPicList).GetMv(uiAbsPartAddr)
    //fmt.Printf("(%d,%d):%d,%d=(%d,%d)\n",pColCU.GetCUPelX(),pColCU.GetCUPelY(),eColRefPicList,uiAbsPartAddr, cColMv.GetHor(), cColMv.GetVer());

    iCurrRefPOC = int(this.m_pcSlice.GetRefPic(eRefPicList, *riRefIdx).GetPOC())
    bIsCurrRefLongTerm := this.m_pcSlice.GetRefPic(eRefPicList, *riRefIdx).GetIsLongTerm()
    bIsColRefLongTerm := pColCU.GetSlice().GetIsUsedAsLongTerm(int(eColRefPicList), int(iColRefIdx))

    if bIsCurrRefLongTerm != bIsColRefLongTerm {
        return false
    }

    if bIsCurrRefLongTerm || bIsColRefLongTerm {
        *rcMv = cColMv
    } else {
        iScale = this.xGetDistScaleFactor(iCurrPOC, iCurrRefPOC, iColPOC, iColRefPOC)
        //fmt.Printf("xGetColMVP: %d,%d,%d,%d=%d",iCurrPOC, iCurrRefPOC, iColPOC, iColRefPOC, iScale);
        if iScale == 4096 {
            *rcMv = cColMv
        } else {
            //fmt.Printf("xGetColMVP:%d",iScale);
            *rcMv = cColMv.ScaleMv(iScale)
        }
    }
    return true
}

/// compute required bits to encode MVD (used in AMVP)
func (this *TComDataCU) xGetMvdBits(cMvd TComMv) uint {
    return (this.xGetComponentBits(int(cMvd.GetHor())) + this.xGetComponentBits(int(cMvd.GetVer())))
}
func (this *TComDataCU) xGetComponentBits(iVal int) uint {
    uiLength := uint(1)
    var uiTemp int
    if iVal <= 0 {
        uiTemp = -(iVal << 1) + 1
    } else {
        uiTemp = (iVal << 1)
    }

    //assert ( uiTemp );

    for 1 != uiTemp {
        uiTemp >>= 1
        uiLength += 2
    }

    return uiLength
}

/// compute scaling factor from POC difference
func (this *TComDataCU) xGetDistScaleFactor(iCurrPOC, iCurrRefPOC, iColPOC, iColRefPOC int) int {
    iDiffPocD := iColPOC - iColRefPOC
    iDiffPocB := iCurrPOC - iCurrRefPOC

    if iDiffPocD != iDiffPocB {
        var iTDB, iTDD, iX, iScale int
        iTDB = iDiffPocB
        if iTDB < -128 {
            iTDB = -128
        } else if iTDB > 127 {
            iTDB = 127
        }

        iTDD = iDiffPocD
        if iTDD < -128 {
            iTDD = -128
        } else if iTDD > 127 {
            iTDD = 127
        }

        if iTDD < 0 {
            iX = (0x4000 - (iTDD / 2)) / iTDD
        } else {
            iX = (0x4000 + (iTDD / 2)) / iTDD
        }

        iScale = (iTDB*iX + 32) >> 6
        if iScale < -4096 {
            iScale = -4096
        } else if iScale > 4095 {
            iScale = 4095
        }
        return iScale
    }

    return 4096
}

func (this *TComDataCU) xDeriveCenterIdx(uiPartIdx uint, ruiPartIdxCenter *uint) {
    var uiPartAddr uint
    var iPartWidth, iPartHeight int
    this.GetPartIndexAndSize(uiPartIdx, &uiPartAddr, &iPartWidth, &iPartHeight)

    *ruiPartIdxCenter = this.m_uiAbsIdxInLCU + uiPartAddr // partition origin.
    *ruiPartIdxCenter = G_auiRasterToZscan[int(G_auiZscanToRaster[*ruiPartIdxCenter])+(iPartHeight/int(this.m_pcPic.GetMinCUHeight()))/2*int(this.m_pcPic.GetNumPartInWidth())+(iPartWidth/int(this.m_pcPic.GetMinCUWidth()))/2]
}

func (this *TComDataCU) InitEstData(uiDepth uint, qp int) {
    this.m_dTotalCost = MAX_DOUBLE
    this.m_uiTotalDistortion = 0
    this.m_uiTotalBits = 0
    this.m_uiTotalBins = 0

    uhWidth := this.GetSlice().GetSPS().GetMaxCUWidth() >> uiDepth
    uhHeight := this.GetSlice().GetSPS().GetMaxCUHeight() >> uiDepth

    for ui := uint(0); ui < this.m_uiNumPartition; ui++ {
        if this.GetPic().GetPicSym().GetInverseCUOrderMap(int(this.GetAddr()))*this.m_pcPic.GetNumPartInCU()+this.m_uiAbsIdxInLCU+ui >= this.GetSlice().GetSliceSegmentCurStartCUAddr() {
            this.m_apiMVPIdx[0][ui] = -1
            this.m_apiMVPIdx[1][ui] = -1
            this.m_apiMVPNum[0][ui] = -1
            this.m_apiMVPNum[1][ui] = -1
            this.m_puhDepth[ui] = byte(uiDepth)
            this.m_puhWidth[ui] = byte(uhWidth)
            this.m_puhHeight[ui] = byte(uhHeight)
            this.m_puhTrIdx[ui] = 0
            this.m_puhTransformSkip[0][ui] = false
            this.m_puhTransformSkip[1][ui] = false
            this.m_puhTransformSkip[2][ui] = false
            this.m_skipFlag[ui] = false
            this.m_pePartSize[ui] = SIZE_NONE
            this.m_pePredMode[ui] = MODE_NONE
            this.m_CUTransquantBypass[ui] = false
            this.m_pbIPCMFlag[ui] = false
            this.m_phQP[ui] = int8(qp)
            this.m_pbMergeFlag[ui] = false
            this.m_puhMergeIndex[ui] = 0
            this.m_puhLumaIntraDir[ui] = DC_IDX
            this.m_puhChromaIntraDir[ui] = 0
            this.m_puhInterDir[ui] = 0
            this.m_puhCbf[0][ui] = 0
            this.m_puhCbf[1][ui] = 0
            this.m_puhCbf[2][ui] = 0
        }
    }

    uiTmp := uhWidth * uhHeight

    if this.GetPic().GetPicSym().GetInverseCUOrderMap(int(this.GetAddr()))*this.m_pcPic.GetNumPartInCU()+this.m_uiAbsIdxInLCU >= this.GetSlice().GetSliceSegmentCurStartCUAddr() {
        this.m_acCUMvField[0].ClearMvField()
        this.m_acCUMvField[1].ClearMvField()
        uiTmp = uhWidth * uhHeight
        for i := uint(0); i < uiTmp; i++ {
            this.m_pcTrCoeffY[i] = 0 //memset( this.m_pcTrCoeffY,    0, uiTmp * sizeof( *this.m_pcTrCoeffY    ) );
            //#if ADAPTIVE_QP_SELECTION
            this.m_pcArlCoeffY[i] = 0 //memset( this.m_pcArlCoeffY ,  0, uiTmp * sizeof( *this.m_pcArlCoeffY   ) );
            //#endif
            this.m_pcIPCMSampleY[i] = 0 //memset( this.m_pcIPCMSampleY, 0, uiTmp * sizeof( *this.m_pcIPCMSampleY ) );
        }
        uiTmp >>= 2
        for i := uint(0); i < uiTmp; i++ {
            this.m_pcTrCoeffCb[i] = 0 //memset( this.m_pcTrCoeffCb,    0, uiTmp * sizeof( *this.m_pcTrCoeffCb    ) );
            this.m_pcTrCoeffCr[i] = 0 //memset( this.m_pcTrCoeffCr,    0, uiTmp * sizeof( *this.m_pcTrCoeffCr    ) );
            //#if ADAPTIVE_QP_SELECTION
            this.m_pcArlCoeffCb[i] = 0 //memset( this.m_pcArlCoeffCb,   0, uiTmp * sizeof( *this.m_pcArlCoeffCb   ) );
            this.m_pcArlCoeffCr[i] = 0 //memset( this.m_pcArlCoeffCr,   0, uiTmp * sizeof( *this.m_pcArlCoeffCr   ) );
            //#endif
            this.m_pcIPCMSampleCb[i] = 0 //memset( this.m_pcIPCMSampleCb, 0, uiTmp * sizeof( *this.m_pcIPCMSampleCb ) );
            this.m_pcIPCMSampleCr[i] = 0 //memset( this.m_pcIPCMSampleCr, 0, uiTmp * sizeof( *this.m_pcIPCMSampleCr ) );
        }
    }
}
func (this *TComDataCU) InitSubCU(pcCU *TComDataCU, uiPartUnitIdx, uiDepth uint, qp int) {
    //assert( uiPartUnitIdx<4 );

    uiPartOffset := (pcCU.GetTotalNumPart() >> 2) * uiPartUnitIdx

    this.m_pcPic = pcCU.GetPic()
    this.m_pcSlice = this.m_pcPic.GetSlice(this.m_pcPic.GetCurrSliceIdx())
    this.m_uiCUAddr = pcCU.GetAddr()
    this.m_uiAbsIdxInLCU = pcCU.GetZorderIdxInCU() + uiPartOffset

    this.m_uiCUPelX = pcCU.GetCUPelX() + (this.GetSlice().GetSPS().GetMaxCUWidth()>>uiDepth)*(uiPartUnitIdx&1)
    this.m_uiCUPelY = pcCU.GetCUPelY() + (this.GetSlice().GetSPS().GetMaxCUHeight()>>uiDepth)*(uiPartUnitIdx>>1)

    this.m_dTotalCost = MAX_DOUBLE
    this.m_uiTotalDistortion = 0
    this.m_uiTotalBits = 0
    this.m_uiTotalBins = 0
    this.m_uiNumPartition = pcCU.GetTotalNumPart() >> 2

    //iSizeInUchar :=  this.m_uiNumPartition;
    //iSizeInBool  := sizeof( Bool   ) * this.m_uiNumPartition;
    //sizeInChar   := sizeof( Char   ) * this.m_uiNumPartition;
    for i := uint(0); i < this.m_uiNumPartition; i++ {
        this.m_phQP[i] = int8(qp)
        this.m_pbMergeFlag[i] = false
        this.m_puhMergeIndex[i] = 0
        this.m_puhLumaIntraDir[i] = DC_IDX
        this.m_puhChromaIntraDir[i] = 0
        this.m_puhInterDir[i] = 0
        this.m_puhTrIdx[i] = 0
        this.m_puhTransformSkip[0][i] = false
        this.m_puhTransformSkip[1][i] = false
        this.m_puhTransformSkip[2][i] = false
        this.m_puhCbf[0][i] = 0
        this.m_puhCbf[1][i] = 0
        this.m_puhCbf[2][i] = 0
        this.m_puhDepth[i] = byte(uiDepth)
    }

    uhWidth := this.GetSlice().GetSPS().GetMaxCUWidth() >> uiDepth
    uhHeight := this.GetSlice().GetSPS().GetMaxCUHeight() >> uiDepth
    for i := uint(0); i < this.m_uiNumPartition; i++ {
        this.m_puhWidth[i] = byte(uhWidth)
        this.m_puhHeight[i] = byte(uhHeight)
        this.m_pbIPCMFlag[i] = false
    }

    for ui := uint(0); ui < this.m_uiNumPartition; ui++ {
        this.m_skipFlag[ui] = false
        this.m_pePartSize[ui] = SIZE_NONE
        this.m_pePredMode[ui] = MODE_NONE
        this.m_CUTransquantBypass[ui] = false
        this.m_apiMVPIdx[0][ui] = -1
        this.m_apiMVPIdx[1][ui] = -1
        this.m_apiMVPNum[0][ui] = -1
        this.m_apiMVPNum[1][ui] = -1
        if this.m_pcPic.GetPicSym().GetInverseCUOrderMap(int(this.GetAddr()))*this.m_pcPic.GetNumPartInCU()+this.m_uiAbsIdxInLCU+ui < this.GetSlice().GetSliceSegmentCurStartCUAddr() {
            this.m_apiMVPIdx[0][ui] = pcCU.m_apiMVPIdx[0][uiPartOffset+ui]
            this.m_apiMVPIdx[1][ui] = pcCU.m_apiMVPIdx[1][uiPartOffset+ui]

            this.m_apiMVPNum[0][ui] = pcCU.m_apiMVPNum[0][uiPartOffset+ui]

            this.m_apiMVPNum[1][ui] = pcCU.m_apiMVPNum[1][uiPartOffset+ui]

            this.m_puhDepth[ui] = pcCU.GetDepth1(uiPartOffset + ui)
            this.m_puhWidth[ui] = pcCU.GetWidth1(uiPartOffset + ui)
            this.m_puhHeight[ui] = pcCU.GetHeight1(uiPartOffset + ui)
            this.m_puhTrIdx[ui] = pcCU.GetTransformIdx1(uiPartOffset + ui)
            this.m_puhTransformSkip[0][ui] = pcCU.GetTransformSkip2(uiPartOffset+ui, TEXT_LUMA)
            this.m_puhTransformSkip[1][ui] = pcCU.GetTransformSkip2(uiPartOffset+ui, TEXT_CHROMA_U)
            this.m_puhTransformSkip[2][ui] = pcCU.GetTransformSkip2(uiPartOffset+ui, TEXT_CHROMA_V)
            this.m_skipFlag[ui] = pcCU.GetSkipFlag1(uiPartOffset + ui)
            this.m_pePartSize[ui] = pcCU.GetPartitionSize1(uiPartOffset + ui)
            this.m_pePredMode[ui] = pcCU.GetPredictionMode1(uiPartOffset + ui)
            this.m_CUTransquantBypass[ui] = pcCU.GetCUTransquantBypass1(uiPartOffset + ui)
            this.m_pbIPCMFlag[ui] = pcCU.m_pbIPCMFlag[uiPartOffset+ui]
            this.m_phQP[ui] = pcCU.m_phQP[uiPartOffset+ui]
            this.m_pbMergeFlag[ui] = pcCU.m_pbMergeFlag[uiPartOffset+ui]
            this.m_puhMergeIndex[ui] = pcCU.m_puhMergeIndex[uiPartOffset+ui]
            this.m_puhLumaIntraDir[ui] = pcCU.m_puhLumaIntraDir[uiPartOffset+ui]
            this.m_puhChromaIntraDir[ui] = pcCU.m_puhChromaIntraDir[uiPartOffset+ui]
            this.m_puhInterDir[ui] = pcCU.m_puhInterDir[uiPartOffset+ui]
            this.m_puhCbf[0][ui] = pcCU.m_puhCbf[0][uiPartOffset+ui]
            this.m_puhCbf[1][ui] = pcCU.m_puhCbf[1][uiPartOffset+ui]
            this.m_puhCbf[2][ui] = pcCU.m_puhCbf[2][uiPartOffset+ui]

        }
    }

    uiTmp := uhWidth * uhHeight
    for i := uint(0); i < uiTmp; i++ {
        this.m_pcTrCoeffY[i] = 0 //memset( , sizeof(TCoeff)*uiTmp );
        //#if ADAPTIVE_QP_SELECTION
        this.m_pcArlCoeffY[i] = 0 //  memset( , sizeof(Int)*uiTmp );
        //#endif
        this.m_pcIPCMSampleY[i] = 0 //memset( , sizeof( Pel ) * uiTmp );
    }
    uiTmp >>= 2
    for i := uint(0); i < uiTmp; i++ {
        this.m_pcTrCoeffCb[i] = 0 //,memset( this.m_pcTrCoeffCb, 0, sizeof(TCoeff)*uiTmp );
        this.m_pcTrCoeffCr[i] = 0 //memset( this.m_pcTrCoeffCr, 0, sizeof(TCoeff)*uiTmp );
        //#if ADAPTIVE_QP_SELECTION
        this.m_pcArlCoeffCb[i] = 0 //memset( this.m_pcArlCoeffCb, 0, sizeof(Int)*uiTmp );
        this.m_pcArlCoeffCr[i] = 0 //memset( this.m_pcArlCoeffCr, 0, sizeof(Int)*uiTmp );
        //#endif
        this.m_pcIPCMSampleCb[i] = 0 //memset( this.m_pcIPCMSampleCb , 0, sizeof( Pel ) * uiTmp );
        this.m_pcIPCMSampleCr[i] = 0 //memset( this.m_pcIPCMSampleCr , 0, sizeof( Pel ) * uiTmp );
    }
    this.m_acCUMvField[0].ClearMvField()
    this.m_acCUMvField[1].ClearMvField()

    if this.m_pcPic.GetPicSym().GetInverseCUOrderMap(int(this.GetAddr()))*this.m_pcPic.GetNumPartInCU()+this.m_uiAbsIdxInLCU < this.GetSlice().GetSliceSegmentCurStartCUAddr() {
        // Part of this CU contains data from an older slice. Now copy in that data.
        uiMaxCuWidth := pcCU.GetSlice().GetSPS().GetMaxCUWidth()
        uiMaxCuHeight := pcCU.GetSlice().GetSPS().GetMaxCUHeight()
        bigCU := this.GetPic().GetCU(this.GetAddr())
        minui := int(uiPartOffset)
        minui = -minui
        pcCU.m_acCUMvField[0].CopyTo4(&this.m_acCUMvField[0], minui, uiPartOffset, this.m_uiNumPartition)
        pcCU.m_acCUMvField[1].CopyTo4(&this.m_acCUMvField[1], minui, uiPartOffset, this.m_uiNumPartition)
        uiCoffOffset := uiMaxCuWidth * uiMaxCuHeight * this.m_uiAbsIdxInLCU / pcCU.GetPic().GetNumPartInCU()
        uiTmp = uhWidth * uhHeight
        for i := uint(0); i < uiTmp; i++ {
            this.m_pcTrCoeffY[i] = bigCU.m_pcTrCoeffY[uiCoffOffset+i]
            //#if ADAPTIVE_QP_SELECTION
            this.m_pcArlCoeffY[i] = bigCU.m_pcArlCoeffY[uiCoffOffset+i]
            //#endif
            this.m_pcIPCMSampleY[i] = bigCU.m_pcIPCMSampleY[uiCoffOffset+i]
        }
        uiTmp >>= 2
        uiCoffOffset >>= 2
        for i := uint(0); i < uiTmp; i++ {
            this.m_pcTrCoeffCr[i] = bigCU.m_pcTrCoeffCr[uiCoffOffset+i]
            this.m_pcTrCoeffCb[i] = bigCU.m_pcTrCoeffCb[uiCoffOffset+i]
            //#if ADAPTIVE_QP_SELECTION
            this.m_pcArlCoeffCr[i] = bigCU.m_pcArlCoeffCr[uiCoffOffset+i]
            this.m_pcArlCoeffCb[i] = bigCU.m_pcArlCoeffCb[uiCoffOffset+i]
            //#endif
            this.m_pcIPCMSampleCb[i] = bigCU.m_pcIPCMSampleCb[uiCoffOffset+i]
            this.m_pcIPCMSampleCr[i] = bigCU.m_pcIPCMSampleCr[uiCoffOffset+i]
        }
    }

    this.m_pcCULeft = pcCU.GetCULeft()
    this.m_pcCUAbove = pcCU.GetCUAbove()
    this.m_pcCUAboveLeft = pcCU.GetCUAboveLeft()
    this.m_pcCUAboveRight = pcCU.GetCUAboveRight()

    this.m_apcCUColocated[0] = pcCU.GetCUColocated(REF_PIC_LIST_0)
    this.m_apcCUColocated[1] = pcCU.GetCUColocated(REF_PIC_LIST_1)
    for i := uint(0); i < this.m_uiNumPartition; i++ {
        this.m_sliceStartCU[i] = pcCU.m_sliceStartCU[i+uiPartOffset]               //,sizeof(UInt)*this.m_uiNumPartition);
        this.m_sliceSegmentStartCU[i] = pcCU.m_sliceSegmentStartCU[i+uiPartOffset] //,sizeof(UInt)*this.m_uiNumPartition);
    }
}
func (this *TComDataCU) SetOutsideCUPart(uiAbsPartIdx, uiDepth uint) {
    uiNumPartition := this.m_uiNumPartition >> (uiDepth << 1)
    //uiSizeInUchar  := sizeof( UChar  ) * uiNumPartition;

    uhWidth := this.GetSlice().GetSPS().GetMaxCUWidth() >> uiDepth
    uhHeight := this.GetSlice().GetSPS().GetMaxCUHeight() >> uiDepth
    for i := uint(0); i < uiNumPartition; i++ {
        this.m_puhDepth[i+uiAbsPartIdx] = byte(uiDepth)
        this.m_puhWidth[i+uiAbsPartIdx] = byte(uhWidth)
        this.m_puhHeight[i+uiAbsPartIdx] = byte(uhHeight)
    }
}

func (this *TComDataCU) CopySubCU(pcCU *TComDataCU, uiAbsPartIdx, uiDepth uint) {
    uiPart := uiAbsPartIdx

    this.m_pcPic = pcCU.GetPic()
    this.m_pcSlice = pcCU.GetSlice()
    this.m_uiCUAddr = pcCU.GetAddr()
    this.m_uiAbsIdxInLCU = uiAbsPartIdx

    this.m_uiCUPelX = pcCU.GetCUPelX() + G_auiRasterToPelX[G_auiZscanToRaster[uiAbsPartIdx]]
    this.m_uiCUPelY = pcCU.GetCUPelY() + G_auiRasterToPelY[G_auiZscanToRaster[uiAbsPartIdx]]

    uiWidth := this.GetSlice().GetSPS().GetMaxCUWidth() >> uiDepth
    uiHeight := this.GetSlice().GetSPS().GetMaxCUHeight() >> uiDepth

    this.m_skipFlag = pcCU.GetSkipFlag()[uiPart:]

    this.m_phQP = pcCU.GetQP()[uiPart:]
    this.m_pePartSize = pcCU.GetPartitionSize()[uiPart:]
    this.m_pePredMode = pcCU.GetPredictionMode()[uiPart:]
    this.m_CUTransquantBypass = pcCU.GetCUTransquantBypass()[uiPart:]

    this.m_pbMergeFlag = pcCU.GetMergeFlag()[uiPart:]
    this.m_puhMergeIndex = pcCU.GetMergeIndex()[uiPart:]

    this.m_puhLumaIntraDir = pcCU.GetLumaIntraDir()[uiPart:]
    this.m_puhChromaIntraDir = pcCU.GetChromaIntraDir()[uiPart:]
    this.m_puhInterDir = pcCU.GetInterDir()[uiPart:]
    this.m_puhTrIdx = pcCU.GetTransformIdx()[uiPart:]
    this.m_puhTransformSkip[0] = pcCU.GetTransformSkip1(TEXT_LUMA)[uiPart:]
    this.m_puhTransformSkip[1] = pcCU.GetTransformSkip1(TEXT_CHROMA_U)[uiPart:]
    this.m_puhTransformSkip[2] = pcCU.GetTransformSkip1(TEXT_CHROMA_V)[uiPart:]

    this.m_puhCbf[0] = pcCU.GetCbf1(TEXT_LUMA)[uiPart:]
    this.m_puhCbf[1] = pcCU.GetCbf1(TEXT_CHROMA_U)[uiPart:]
    this.m_puhCbf[2] = pcCU.GetCbf1(TEXT_CHROMA_V)[uiPart:]

    this.m_puhDepth = pcCU.GetDepth()[uiPart:]
    this.m_puhWidth = pcCU.GetWidth()[uiPart:]
    this.m_puhHeight = pcCU.GetHeight()[uiPart:]

    this.m_apiMVPIdx[0] = pcCU.GetMVPIdx1(REF_PIC_LIST_0)[uiPart:]
    this.m_apiMVPIdx[1] = pcCU.GetMVPIdx1(REF_PIC_LIST_1)[uiPart:]
    this.m_apiMVPNum[0] = pcCU.GetMVPNum1(REF_PIC_LIST_0)[uiPart:]
    this.m_apiMVPNum[1] = pcCU.GetMVPNum1(REF_PIC_LIST_1)[uiPart:]

    this.m_pbIPCMFlag = pcCU.GetIPCMFlag()[uiPart:]

    this.m_pcCUAboveLeft = pcCU.GetCUAboveLeft()
    this.m_pcCUAboveRight = pcCU.GetCUAboveRight()
    this.m_pcCUAbove = pcCU.GetCUAbove()
    this.m_pcCULeft = pcCU.GetCULeft()

    this.m_apcCUColocated[0] = pcCU.GetCUColocated(REF_PIC_LIST_0)
    this.m_apcCUColocated[1] = pcCU.GetCUColocated(REF_PIC_LIST_1)

    uiTmp := uiWidth * uiHeight
    uiMaxCuWidth := pcCU.GetSlice().GetSPS().GetMaxCUWidth()
    uiMaxCuHeight := pcCU.GetSlice().GetSPS().GetMaxCUHeight()

    uiCoffOffset := uiMaxCuWidth * uiMaxCuHeight * uiAbsPartIdx / pcCU.GetPic().GetNumPartInCU()

    this.m_pcTrCoeffY = pcCU.GetCoeffY()[uiCoffOffset:]
    //#if ADAPTIVE_QP_SELECTION
    this.m_pcArlCoeffY = pcCU.GetArlCoeffY()[uiCoffOffset:]
    //#endif
    this.m_pcIPCMSampleY = pcCU.GetPCMSampleY()[uiCoffOffset:]

    uiTmp >>= 2
    uiCoffOffset >>= 2
    this.m_pcTrCoeffCb = pcCU.GetCoeffCb()[uiCoffOffset:]
    this.m_pcTrCoeffCr = pcCU.GetCoeffCr()[uiCoffOffset:]
    //#if ADAPTIVE_QP_SELECTION
    this.m_pcArlCoeffCb = pcCU.GetArlCoeffCb()[uiCoffOffset:]
    this.m_pcArlCoeffCr = pcCU.GetArlCoeffCr()[uiCoffOffset:]
    //#endif
    this.m_pcIPCMSampleCb = pcCU.GetPCMSampleCb()[uiCoffOffset:]
    this.m_pcIPCMSampleCr = pcCU.GetPCMSampleCr()[uiCoffOffset:]

    this.m_acCUMvField[0].LinkToWithOffset(pcCU.GetCUMvField(REF_PIC_LIST_0), int(uiPart))
    this.m_acCUMvField[1].LinkToWithOffset(pcCU.GetCUMvField(REF_PIC_LIST_1), int(uiPart))

    for i := uint(0); i < this.m_uiNumPartition; i++ {
        this.m_sliceStartCU[i] = pcCU.m_sliceStartCU[i+uiPart]               //,sizeof(UInt)*this.m_uiNumPartition);
        this.m_sliceSegmentStartCU[i] = pcCU.m_sliceSegmentStartCU[i+uiPart] //,sizeof(UInt)*this.m_uiNumPartition);
    }
}

func (this *TComDataCU) CopyInterPredInfoFrom(pcCU *TComDataCU, uiAbsPartIdx uint, eRefPicList RefPicList) {
    this.m_pcPic = pcCU.GetPic()
    this.m_pcSlice = pcCU.GetSlice()
    this.m_uiCUAddr = pcCU.GetAddr()
    this.m_uiAbsIdxInLCU = uiAbsPartIdx

    iRastPartIdx := G_auiZscanToRaster[uiAbsPartIdx]
    this.m_uiCUPelX = pcCU.GetCUPelX() + this.m_pcPic.GetMinCUWidth()*(iRastPartIdx%this.m_pcPic.GetNumPartInWidth())
    this.m_uiCUPelY = pcCU.GetCUPelY() + this.m_pcPic.GetMinCUHeight()*(iRastPartIdx/this.m_pcPic.GetNumPartInWidth())

    this.m_pcCUAboveLeft = pcCU.GetCUAboveLeft()
    this.m_pcCUAboveRight = pcCU.GetCUAboveRight()
    this.m_pcCUAbove = pcCU.GetCUAbove()
    this.m_pcCULeft = pcCU.GetCULeft()

    this.m_apcCUColocated[0] = pcCU.GetCUColocated(REF_PIC_LIST_0)
    this.m_apcCUColocated[1] = pcCU.GetCUColocated(REF_PIC_LIST_1)

    this.m_skipFlag = pcCU.GetSkipFlag()[uiAbsPartIdx:]

    this.m_pePartSize = pcCU.GetPartitionSize()[uiAbsPartIdx:]
    this.m_pePredMode = pcCU.GetPredictionMode()[uiAbsPartIdx:]
    this.m_CUTransquantBypass = pcCU.GetCUTransquantBypass()[uiAbsPartIdx:]
    this.m_puhInterDir = pcCU.GetInterDir()[uiAbsPartIdx:]

    this.m_puhDepth = pcCU.GetDepth()[uiAbsPartIdx:]
    this.m_puhWidth = pcCU.GetWidth()[uiAbsPartIdx:]
    this.m_puhHeight = pcCU.GetHeight()[uiAbsPartIdx:]

    this.m_pbMergeFlag = pcCU.GetMergeFlag()[uiAbsPartIdx:]
    this.m_puhMergeIndex = pcCU.GetMergeIndex()[uiAbsPartIdx:]

    this.m_apiMVPIdx[eRefPicList] = pcCU.GetMVPIdx1(eRefPicList)[uiAbsPartIdx:]
    this.m_apiMVPNum[eRefPicList] = pcCU.GetMVPNum1(eRefPicList)[uiAbsPartIdx:]

    this.m_acCUMvField[eRefPicList].LinkToWithOffset(pcCU.GetCUMvField(eRefPicList), int(uiAbsPartIdx))

    for i := uint(0); i < this.m_uiNumPartition; i++ {
        this.m_sliceStartCU[i] = pcCU.m_sliceStartCU[i+uiAbsPartIdx]               //,sizeof(UInt)*this.m_uiNumPartition);
        this.m_sliceSegmentStartCU[i] = pcCU.m_sliceSegmentStartCU[i+uiAbsPartIdx] //,sizeof(UInt)*this.m_uiNumPartition);
    }
}
func (this *TComDataCU) CopyPartFrom(pcCU *TComDataCU, uiPartUnitIdx, uiDepth uint) {
    //assert( uiPartUnitIdx<4 );

    this.m_dTotalCost += pcCU.GetTotalCost()
    this.m_uiTotalDistortion += pcCU.GetTotalDistortion()
    this.m_uiTotalBits += pcCU.GetTotalBits()

    uiOffset := pcCU.GetTotalNumPart() * uiPartUnitIdx

    uiNumPartition := pcCU.GetTotalNumPart()
    //iSizeInUchar   := sizeof( UChar ) * uiNumPartition;
    //iSizeInBool    := sizeof( Bool  ) * uiNumPartition;
    //sizeInChar     := sizeof( Char ) * uiNumPartition;
    for i := uint(0); i < uiNumPartition; i++ {
        this.m_skipFlag[uiOffset+i] = pcCU.GetSkipFlag1(i)                                //,       sizeof( *this.m_skipFlag )   * uiNumPartition );
        this.m_phQP[uiOffset+i] = pcCU.GetQP1(i)                                          //,             sizeInChar                        );
        this.m_pePartSize[uiOffset+i] = pcCU.GetPartitionSize1(i)                         //,  sizeof( *this.m_pePartSize ) * uiNumPartition );
        this.m_pePredMode[uiOffset+i] = pcCU.GetPredictionMode1(i)                        //, sizeof( *this.m_pePredMode ) * uiNumPartition );
        this.m_CUTransquantBypass[uiOffset+i] = pcCU.GetCUTransquantBypass1(i)            //, sizeof( *this.m_CUTransquantBypass ) * uiNumPartition );
        this.m_pbMergeFlag[uiOffset+i] = pcCU.GetMergeFlag1(i)                            //,         iSizeInBool  );
        this.m_puhMergeIndex[uiOffset+i] = pcCU.GetMergeIndex1(i)                         //,        iSizeInUchar );
        this.m_puhLumaIntraDir[uiOffset+i] = pcCU.GetLumaIntraDir1(i)                     //,      iSizeInUchar );
        this.m_puhChromaIntraDir[uiOffset+i] = pcCU.GetChromaIntraDir1(i)                 //,    iSizeInUchar );
        this.m_puhInterDir[uiOffset+i] = pcCU.GetInterDir1(i)                             //,          iSizeInUchar );
        this.m_puhTrIdx[uiOffset+i] = pcCU.GetTransformIdx1(i)                            //,      iSizeInUchar );
        this.m_puhTransformSkip[0][uiOffset+i] = pcCU.GetTransformSkip2(i, TEXT_LUMA)     //,     iSizeInUchar );
        this.m_puhTransformSkip[1][uiOffset+i] = pcCU.GetTransformSkip2(i, TEXT_CHROMA_U) //, iSizeInUchar );
        this.m_puhTransformSkip[2][uiOffset+i] = pcCU.GetTransformSkip2(i, TEXT_CHROMA_V) //, iSizeInUchar );

        this.m_puhCbf[0][uiOffset+i] = pcCU.GetCbf2(i, TEXT_LUMA)     //    , iSizeInUchar );
        this.m_puhCbf[1][uiOffset+i] = pcCU.GetCbf2(i, TEXT_CHROMA_U) //, iSizeInUchar );
        this.m_puhCbf[2][uiOffset+i] = pcCU.GetCbf2(i, TEXT_CHROMA_V) //, iSizeInUchar );

        this.m_puhDepth[uiOffset+i] = pcCU.GetDepth1(i)   //,  iSizeInUchar );
        this.m_puhWidth[uiOffset+i] = pcCU.GetWidth1(i)   //,  iSizeInUchar );
        this.m_puhHeight[uiOffset+i] = pcCU.GetHeight1(i) //, iSizeInUchar );

        this.m_apiMVPIdx[0][uiOffset+i] = pcCU.GetMVPIdx2(REF_PIC_LIST_0, i) //, iSizeInUchar );
        this.m_apiMVPIdx[1][uiOffset+i] = pcCU.GetMVPIdx2(REF_PIC_LIST_1, i) //, iSizeInUchar );
        this.m_apiMVPNum[0][uiOffset+i] = pcCU.GetMVPNum2(REF_PIC_LIST_0, i) //, iSizeInUchar );
        this.m_apiMVPNum[1][uiOffset+i] = pcCU.GetMVPNum2(REF_PIC_LIST_1, i) //, iSizeInUchar );

        this.m_pbIPCMFlag[uiOffset+i] = pcCU.GetIPCMFlag1(i) //, iSizeInBool );
    }

    this.m_pcCUAboveLeft = pcCU.GetCUAboveLeft()
    this.m_pcCUAboveRight = pcCU.GetCUAboveRight()
    this.m_pcCUAbove = pcCU.GetCUAbove()
    this.m_pcCULeft = pcCU.GetCULeft()

    this.m_apcCUColocated[0] = pcCU.GetCUColocated(REF_PIC_LIST_0)
    this.m_apcCUColocated[1] = pcCU.GetCUColocated(REF_PIC_LIST_1)

    this.m_acCUMvField[0].CopyFrom(pcCU.GetCUMvField(REF_PIC_LIST_0), pcCU.GetTotalNumPart(), int(uiOffset))
    this.m_acCUMvField[1].CopyFrom(pcCU.GetCUMvField(REF_PIC_LIST_1), pcCU.GetTotalNumPart(), int(uiOffset))

    uiTmp := this.GetSlice().GetSPS().GetMaxCUWidth() * this.GetSlice().GetSPS().GetMaxCUHeight() >> (uiDepth << 1)
    uiTmp2 := uiPartUnitIdx * uiTmp
    for i := uint(0); i < uiTmp; i++ {
        this.m_pcTrCoeffY[i+uiTmp2] = pcCU.GetCoeffY()[i]
        //#if ADAPTIVE_QP_SELECTION
        this.m_pcArlCoeffY[i+uiTmp2] = pcCU.GetArlCoeffY()[i]
        //#endif
        this.m_pcIPCMSampleY[i+uiTmp2] = pcCU.GetPCMSampleY()[i]
    }

    uiTmp >>= 2
    uiTmp2 >>= 2
    for i := uint(0); i < uiTmp; i++ {
        this.m_pcTrCoeffCb[i+uiTmp2] = pcCU.GetCoeffCb()[i] //, sizeof(TCoeff)*uiTmp );
        this.m_pcTrCoeffCr[i+uiTmp2] = pcCU.GetCoeffCr()[i] //, sizeof(TCoeff)*uiTmp );
        //#if ADAPTIVE_QP_SELECTION
        this.m_pcArlCoeffCb[i+uiTmp2] = pcCU.GetArlCoeffCb()[i] //, sizeof(Int)*uiTmp );
        this.m_pcArlCoeffCr[i+uiTmp2] = pcCU.GetArlCoeffCr()[i] //, sizeof(Int)*uiTmp );
        //#endif
        this.m_pcIPCMSampleCb[i+uiTmp2] = pcCU.GetPCMSampleCb()[i] //, sizeof(Pel) * uiTmp );
        this.m_pcIPCMSampleCr[i+uiTmp2] = pcCU.GetPCMSampleCr()[i] //, sizeof(Pel) * uiTmp );
    }
    this.m_uiTotalBins += pcCU.GetTotalBins()

    for i := uint(0); i < uiNumPartition; i++ {
        this.m_sliceStartCU[i+uiOffset] = pcCU.m_sliceStartCU[i]               //,sizeof(UInt)*this.m_uiNumPartition);
        this.m_sliceSegmentStartCU[i+uiOffset] = pcCU.m_sliceSegmentStartCU[i] //,sizeof(UInt)*this.m_uiNumPartition);
    }
}

func (this *TComDataCU) CopyToPic1(uhDepth uint) {
    rpcCU := this.m_pcPic.GetCU(this.m_uiCUAddr)

    rpcCU.SetTotalCost(this.m_dTotalCost)
    rpcCU.SetTotalDistortion(this.m_uiTotalDistortion)
    rpcCU.SetTotalBits(this.m_uiTotalBits)

    //iSizeInUchar  := sizeof( UChar ) * this.m_uiNumPartition;
    //iSizeInBool   := sizeof( Bool  ) * this.m_uiNumPartition;
    //sizeInChar    := sizeof( Char  ) * this.m_uiNumPartition;

    for i := uint(0); i < this.m_uiNumPartition; i++ {
        rpcCU.GetSkipFlag()[i+this.m_uiAbsIdxInLCU] = this.m_skipFlag[i]                               //, sizeof( *this.m_skipFlag ) * this.m_uiNumPartition );
        rpcCU.GetQP()[i+this.m_uiAbsIdxInLCU] = this.m_phQP[i]                                         //, sizeInChar  );
        rpcCU.GetPartitionSize()[i+this.m_uiAbsIdxInLCU] = this.m_pePartSize[i]                        //, sizeof( *this.m_pePartSize ) * this.m_uiNumPartition );
        rpcCU.GetPredictionMode()[i+this.m_uiAbsIdxInLCU] = this.m_pePredMode[i]                       //, sizeof( *this.m_pePredMode ) * this.m_uiNumPartition );
        rpcCU.GetCUTransquantBypass()[i+this.m_uiAbsIdxInLCU] = this.m_CUTransquantBypass[i]           //, sizeof( *this.m_CUTransquantBypass ) * this.m_uiNumPartition );
        rpcCU.GetMergeFlag()[i+this.m_uiAbsIdxInLCU] = this.m_pbMergeFlag[i]                           //,         iSizeInBool  );
        rpcCU.GetMergeIndex()[i+this.m_uiAbsIdxInLCU] = this.m_puhMergeIndex[i]                        //,       iSizeInUchar );
        rpcCU.GetLumaIntraDir()[i+this.m_uiAbsIdxInLCU] = this.m_puhLumaIntraDir[i]                    //,     iSizeInUchar );
        rpcCU.GetChromaIntraDir()[i+this.m_uiAbsIdxInLCU] = this.m_puhChromaIntraDir[i]                //,   iSizeInUchar );
        rpcCU.GetInterDir()[i+this.m_uiAbsIdxInLCU] = this.m_puhInterDir[i]                            //,         iSizeInUchar );
        rpcCU.GetTransformIdx()[i+this.m_uiAbsIdxInLCU] = this.m_puhTrIdx[i]                           //,            iSizeInUchar );
        rpcCU.GetTransformSkip1(TEXT_LUMA)[i+this.m_uiAbsIdxInLCU] = this.m_puhTransformSkip[0][i]     //, iSizeInUchar );
        rpcCU.GetTransformSkip1(TEXT_CHROMA_U)[i+this.m_uiAbsIdxInLCU] = this.m_puhTransformSkip[1][i] //, iSizeInUchar );
        rpcCU.GetTransformSkip1(TEXT_CHROMA_V)[i+this.m_uiAbsIdxInLCU] = this.m_puhTransformSkip[2][i] //, iSizeInUchar );

        rpcCU.GetCbf1(TEXT_LUMA)[i+this.m_uiAbsIdxInLCU] = this.m_puhCbf[0][i]     //, iSizeInUchar );
        rpcCU.GetCbf1(TEXT_CHROMA_U)[i+this.m_uiAbsIdxInLCU] = this.m_puhCbf[1][i] //, iSizeInUchar );
        rpcCU.GetCbf1(TEXT_CHROMA_V)[i+this.m_uiAbsIdxInLCU] = this.m_puhCbf[2][i] //, iSizeInUchar );

        rpcCU.GetDepth()[i+this.m_uiAbsIdxInLCU] = this.m_puhDepth[i]   //,  iSizeInUchar );
        rpcCU.GetWidth()[i+this.m_uiAbsIdxInLCU] = this.m_puhWidth[i]   //,  iSizeInUchar );
        rpcCU.GetHeight()[i+this.m_uiAbsIdxInLCU] = this.m_puhHeight[i] //, iSizeInUchar );

        rpcCU.GetMVPIdx1(REF_PIC_LIST_0)[i+this.m_uiAbsIdxInLCU] = this.m_apiMVPIdx[0][i] //, iSizeInUchar );
        rpcCU.GetMVPIdx1(REF_PIC_LIST_1)[i+this.m_uiAbsIdxInLCU] = this.m_apiMVPIdx[1][i] //, iSizeInUchar );
        rpcCU.GetMVPNum1(REF_PIC_LIST_0)[i+this.m_uiAbsIdxInLCU] = this.m_apiMVPNum[0][i] //, iSizeInUchar );
        rpcCU.GetMVPNum1(REF_PIC_LIST_1)[i+this.m_uiAbsIdxInLCU] = this.m_apiMVPNum[1][i] //, iSizeInUchar );

        rpcCU.GetIPCMFlag()[i+this.m_uiAbsIdxInLCU] = this.m_pbIPCMFlag[i] //,         iSizeInBool  );
    }

    this.m_acCUMvField[0].CopyTo2(rpcCU.GetCUMvField(REF_PIC_LIST_0), int(this.m_uiAbsIdxInLCU))
    this.m_acCUMvField[1].CopyTo2(rpcCU.GetCUMvField(REF_PIC_LIST_1), int(this.m_uiAbsIdxInLCU))

    uiTmp := (this.GetSlice().GetSPS().GetMaxCUWidth() * this.GetSlice().GetSPS().GetMaxCUHeight()) >> (uhDepth << 1)
    uiTmp2 := this.m_uiAbsIdxInLCU * this.m_pcPic.GetMinCUWidth() * this.m_pcPic.GetMinCUHeight()
    for i := uint(0); i < uiTmp; i++ {
        rpcCU.GetCoeffY()[i+uiTmp2] = this.m_pcTrCoeffY[i] //,  sizeof(TCoeff)*uiTmp  );
        //#if ADAPTIVE_QP_SELECTION
        rpcCU.GetArlCoeffY()[i+uiTmp2] = this.m_pcArlCoeffY[i] //,  sizeof(Int)*uiTmp  );
        //#endif
        rpcCU.GetPCMSampleY()[i+uiTmp2] = this.m_pcIPCMSampleY[i] //, sizeof(Pel)*uiTmp );
    }

    uiTmp >>= 2
    uiTmp2 >>= 2
    for i := uint(0); i < uiTmp; i++ {
        rpcCU.GetCoeffCb()[i+uiTmp2] = this.m_pcTrCoeffCb[i] //, sizeof(TCoeff)*uiTmp  );
        rpcCU.GetCoeffCr()[i+uiTmp2] = this.m_pcTrCoeffCr[i] //, sizeof(TCoeff)*uiTmp  );
        //#if ADAPTIVE_QP_SELECTION
        rpcCU.GetArlCoeffCb()[i+uiTmp2] = this.m_pcArlCoeffCb[i] //, sizeof(Int)*uiTmp  );
        rpcCU.GetArlCoeffCr()[i+uiTmp2] = this.m_pcArlCoeffCr[i] //, sizeof(Int)*uiTmp  );
        //#endif
        rpcCU.GetPCMSampleCb()[i+uiTmp2] = this.m_pcIPCMSampleCb[i] //, sizeof( Pel ) * uiTmp );
        rpcCU.GetPCMSampleCr()[i+uiTmp2] = this.m_pcIPCMSampleCr[i] //, sizeof( Pel ) * uiTmp );
    }

    rpcCU.SetTotalBins(this.m_uiTotalBins)

    for i := uint(0); i < this.m_uiNumPartition; i++ {
        rpcCU.m_sliceStartCU[i+this.m_uiAbsIdxInLCU] = this.m_sliceStartCU[i]               //,sizeof(UInt)*this.m_uiNumPartition);
        rpcCU.m_sliceSegmentStartCU[i+this.m_uiAbsIdxInLCU] = this.m_sliceSegmentStartCU[i] //,sizeof(UInt)*this.m_uiNumPartition);
    }
}
func (this *TComDataCU) CopyToPic3(uhDepth, uiPartIdx, uiPartDepth uint) {
    rpcCU := this.m_pcPic.GetCU(this.m_uiCUAddr)
    uiQNumPart := this.m_uiNumPartition >> (uiPartDepth << 1)

    uiPartStart := uiPartIdx * uiQNumPart
    uiPartOffset := this.m_uiAbsIdxInLCU + uiPartStart

    rpcCU.SetTotalCost(this.m_dTotalCost)
    rpcCU.SetTotalDistortion(this.m_uiTotalDistortion)
    rpcCU.SetTotalBits(this.m_uiTotalBits)

    //iSizeInUchar  := sizeof( UChar  ) * uiQNumPart;
    //iSizeInBool   := sizeof( Bool   ) * uiQNumPart;
    //sizeInChar    := sizeof( Char  ) * uiQNumPart;
    for i := uint(0); i < uiQNumPart; i++ {
        rpcCU.GetSkipFlag()[i+uiPartOffset] = this.m_skipFlag[i]                               //,   sizeof( *this.m_skipFlag )   * uiQNumPart );
        rpcCU.GetQP()[i+uiPartOffset] = this.m_phQP[i]                                         //, sizeInChar );
        rpcCU.GetPartitionSize()[i+uiPartOffset] = this.m_pePartSize[i]                        //, sizeof( *this.m_pePartSize ) * uiQNumPart );
        rpcCU.GetPredictionMode()[i+uiPartOffset] = this.m_pePredMode[i]                       //, sizeof( *this.m_pePredMode ) * uiQNumPart );
        rpcCU.GetCUTransquantBypass()[i+uiPartOffset] = this.m_CUTransquantBypass[i]           //, sizeof( *this.m_CUTransquantBypass ) * uiQNumPart );
        rpcCU.GetMergeFlag()[i+uiPartOffset] = this.m_pbMergeFlag[i]                           //,         iSizeInBool  );
        rpcCU.GetMergeIndex()[i+uiPartOffset] = this.m_puhMergeIndex[i]                        //,       iSizeInUchar );
        rpcCU.GetLumaIntraDir()[i+uiPartOffset] = this.m_puhLumaIntraDir[i]                    //,     iSizeInUchar );
        rpcCU.GetChromaIntraDir()[i+uiPartOffset] = this.m_puhChromaIntraDir[i]                //,   iSizeInUchar );
        rpcCU.GetInterDir()[i+uiPartOffset] = this.m_puhInterDir[i]                            //,         iSizeInUchar );
        rpcCU.GetTransformIdx()[i+uiPartOffset] = this.m_puhTrIdx[i]                           //,            iSizeInUchar );
        rpcCU.GetTransformSkip1(TEXT_LUMA)[i+uiPartOffset] = this.m_puhTransformSkip[0][i]     //, iSizeInUchar );
        rpcCU.GetTransformSkip1(TEXT_CHROMA_U)[i+uiPartOffset] = this.m_puhTransformSkip[1][i] //, iSizeInUchar );
        rpcCU.GetTransformSkip1(TEXT_CHROMA_V)[i+uiPartOffset] = this.m_puhTransformSkip[2][i] //, iSizeInUchar );
        rpcCU.GetCbf1(TEXT_LUMA)[i+uiPartOffset] = this.m_puhCbf[0][i]                         //, iSizeInUchar );
        rpcCU.GetCbf1(TEXT_CHROMA_U)[i+uiPartOffset] = this.m_puhCbf[1][i]                     //, iSizeInUchar );
        rpcCU.GetCbf1(TEXT_CHROMA_V)[i+uiPartOffset] = this.m_puhCbf[2][i]                     //, iSizeInUchar );

        rpcCU.GetDepth()[i+uiPartOffset] = this.m_puhDepth[i]   //,  iSizeInUchar );
        rpcCU.GetWidth()[i+uiPartOffset] = this.m_puhWidth[i]   //,  iSizeInUchar );
        rpcCU.GetHeight()[i+uiPartOffset] = this.m_puhHeight[i] //, iSizeInUchar );

        rpcCU.GetMVPIdx1(REF_PIC_LIST_0)[i+uiPartOffset] = this.m_apiMVPIdx[0][i] //, iSizeInUchar );
        rpcCU.GetMVPIdx1(REF_PIC_LIST_1)[i+uiPartOffset] = this.m_apiMVPIdx[1][i] //, iSizeInUchar );
        rpcCU.GetMVPNum1(REF_PIC_LIST_0)[i+uiPartOffset] = this.m_apiMVPNum[0][i] //, iSizeInUchar );
        rpcCU.GetMVPNum1(REF_PIC_LIST_1)[i+uiPartOffset] = this.m_apiMVPNum[1][i] //, iSizeInUchar );
        rpcCU.GetIPCMFlag()[i+uiPartOffset] = this.m_pbIPCMFlag[i]                //,         iSizeInBool  );
    }
    this.m_acCUMvField[0].CopyTo4(rpcCU.GetCUMvField(REF_PIC_LIST_0), int(this.m_uiAbsIdxInLCU), uiPartStart, uiQNumPart)
    this.m_acCUMvField[1].CopyTo4(rpcCU.GetCUMvField(REF_PIC_LIST_1), int(this.m_uiAbsIdxInLCU), uiPartStart, uiQNumPart)

    uiTmp := (this.GetSlice().GetSPS().GetMaxCUWidth() * this.GetSlice().GetSPS().GetMaxCUHeight()) >> ((uhDepth + uiPartDepth) << 1)
    uiTmp2 := uiPartOffset * this.m_pcPic.GetMinCUWidth() * this.m_pcPic.GetMinCUHeight()
    for i := uint(0); i < uiTmp; i++ {
        rpcCU.GetCoeffY()[i+uiTmp2] = this.m_pcTrCoeffY[i] //,  sizeof(TCoeff)*uiTmp  );
        //#if ADAPTIVE_QP_SELECTION
        rpcCU.GetArlCoeffY()[i+uiTmp2] = this.m_pcArlCoeffY[i] //,  sizeof(Int)*uiTmp  );
        //#endif
        rpcCU.GetPCMSampleY()[i+uiTmp2] = this.m_pcIPCMSampleY[i] //, sizeof( Pel ) * uiTmp );
    }

    uiTmp >>= 2
    uiTmp2 >>= 2
    for i := uint(0); i < uiTmp; i++ {
        rpcCU.GetCoeffCb()[i+uiTmp2] = this.m_pcTrCoeffCb[i] //, sizeof(TCoeff)*uiTmp  );
        rpcCU.GetCoeffCr()[i+uiTmp2] = this.m_pcTrCoeffCr[i] //, sizeof(TCoeff)*uiTmp  );
        //#if ADAPTIVE_QP_SELECTION
        rpcCU.GetArlCoeffCb()[i+uiTmp2] = this.m_pcArlCoeffCb[i] //, sizeof(Int)*uiTmp  );
        rpcCU.GetArlCoeffCr()[i+uiTmp2] = this.m_pcArlCoeffCr[i] //, sizeof(Int)*uiTmp  );
        //#endif
        rpcCU.GetPCMSampleCb()[i+uiTmp2] = this.m_pcIPCMSampleCb[i] //, sizeof( Pel ) * uiTmp );
        rpcCU.GetPCMSampleCr()[i+uiTmp2] = this.m_pcIPCMSampleCr[i] //, sizeof( Pel ) * uiTmp );
    }

    rpcCU.SetTotalBins(this.m_uiTotalBins)

    for i := uint(0); i < uiQNumPart; i++ {
        rpcCU.m_sliceStartCU[i+uiPartOffset] = this.m_sliceStartCU[i]               //,sizeof(UInt)*this.m_uiNumPartition);
        rpcCU.m_sliceSegmentStartCU[i+uiPartOffset] = this.m_sliceSegmentStartCU[i] //,sizeof(UInt)*this.m_uiNumPartition);
    }
}

// -------------------------------------------------------------------------------------------------------------------
// member functions for CU description
// -------------------------------------------------------------------------------------------------------------------

func (this *TComDataCU) GetPic() *TComPic {
    return this.m_pcPic
}

func (this *TComDataCU) GetSlice() *TComSlice {
    return this.m_pcSlice
}
func (this *TComDataCU) GetAddr() uint {
    return this.m_uiCUAddr
}
func (this *TComDataCU) GetZorderIdxInCU() uint {
    return this.m_uiAbsIdxInLCU
}
func (this *TComDataCU) GetSCUAddr() uint {
    return this.GetPic().GetPicSym().GetInverseCUOrderMap(int(this.m_uiCUAddr))*(1<<(this.m_pcSlice.GetSPS().GetMaxCUDepth()<<1)) + this.m_uiAbsIdxInLCU
}
func (this *TComDataCU) GetCUPelX() uint {
    return this.m_uiCUPelX
}
func (this *TComDataCU) GetCUPelY() uint {
    return this.m_uiCUPelY
}
func (this *TComDataCU) GetPattern() *TComPattern {
    return this.m_pcPattern
}

func (this *TComDataCU) GetDepth() []byte {
    return this.m_puhDepth
}
func (this *TComDataCU) GetDepth1(uiIdx uint) byte {
    return this.m_puhDepth[uiIdx]
}
func (this *TComDataCU) SetDepth(uiIdx uint, uh byte) {
    this.m_puhDepth[uiIdx] = uh
}

func (this *TComDataCU) SetDepthSubParts(uiDepth, uiAbsPartIdx uint) {
    uiCurrPartNumb := this.m_pcPic.GetNumPartInCU() >> (uiDepth << 1)
    for i := uint(0); i < uiCurrPartNumb; i++ {
        this.m_puhDepth[i+uiAbsPartIdx] = byte(uiDepth)
    }
}

// -------------------------------------------------------------------------------------------------------------------
// member functions for CU data
// -------------------------------------------------------------------------------------------------------------------

func (this *TComDataCU) GetPartitionSize() []PartSize {
    return this.m_pePartSize
}
func (this *TComDataCU) GetPartitionSize1(uiIdx uint) PartSize {
    return PartSize(this.m_pePartSize[uiIdx])
}
func (this *TComDataCU) SetPartitionSize(uiIdx uint, uh PartSize) {
    this.m_pePartSize[uiIdx] = uh
}
func (this *TComDataCU) SetPartSizeSubParts(eMode PartSize, uiAbsPartIdx, uiDepth uint) {
    //assert( sizeof( *m_pePartSize) == 1 );
    for i := uint(0); i < this.m_pcPic.GetNumPartInCU()>>(2*uiDepth); i++ {
        this.m_pePartSize[i+uiAbsPartIdx] = eMode
    }
    //memset( m_pePartSize + uiAbsPartIdx, eMode, m_pcPic.GetNumPartInCU() >> ( 2 * uiDepth ) );
}
func (this *TComDataCU) SetCUTransquantBypassSubParts(flag bool, uiAbsPartIdx, uiDepth uint) {
    for i := uint(0); i < this.m_pcPic.GetNumPartInCU()>>(2*uiDepth); i++ {
        this.m_CUTransquantBypass[i+uiAbsPartIdx] = flag
    }
    //memset( m_CUTransquantBypass + uiAbsPartIdx, flag, m_pcPic.GetNumPartInCU() >> ( 2 * uiDepth ) );
}

func (this *TComDataCU) GetSkipFlag() []bool {
    return this.m_skipFlag
}
func (this *TComDataCU) GetSkipFlag1(idx uint) bool {
    return this.m_skipFlag[idx]
}
func (this *TComDataCU) SetSkipFlag(idx uint, skip bool) {
    this.m_skipFlag[idx] = skip
}
func (this *TComDataCU) SetSkipFlagSubParts(skip bool, absPartIdx, depth uint) {
    //assert( sizeof( *m_skipFlag) == 1 );
    for i := uint(0); i < this.m_pcPic.GetNumPartInCU()>>(2*depth); i++ {
        this.m_skipFlag[i+absPartIdx] = skip
    }
    //memset( m_skipFlag + absPartIdx, skip, m_pcPic.GetNumPartInCU() >> ( 2 * depth ) );
}

func (this *TComDataCU) GetPredictionMode() []PredMode {
    return this.m_pePredMode
}
func (this *TComDataCU) GetPredictionMode1(uiIdx uint) PredMode {
    return PredMode(this.m_pePredMode[uiIdx])
}
func (this *TComDataCU) GetCUTransquantBypass() []bool {
    return this.m_CUTransquantBypass
}
func (this *TComDataCU) GetCUTransquantBypass1(uiIdx uint) bool {
    return this.m_CUTransquantBypass[uiIdx]
}
func (this *TComDataCU) SetPredictionMode(uiIdx uint, uh PredMode) {
    this.m_pePredMode[uiIdx] = uh
}
func (this *TComDataCU) SetPredModeSubParts(eMode PredMode, uiAbsPartIdx, uiDepth uint) {
    //assert( sizeof( *m_pePredMode) == 1 );
    for i := uint(0); i < this.m_pcPic.GetNumPartInCU()>>(2*uiDepth); i++ {
        this.m_pePredMode[i+uiAbsPartIdx] = eMode
    }
    //memset( m_pePredMode + uiAbsPartIdx, eMode, m_pcPic.GetNumPartInCU() >> ( 2 * uiDepth ) );
}

func (this *TComDataCU) GetWidth() []byte {
    return this.m_puhWidth
}
func (this *TComDataCU) GetWidth1(uiIdx uint) byte {
    return this.m_puhWidth[uiIdx]
}

func (this *TComDataCU) SetWidth(uiIdx uint, uh byte) {
    this.m_puhWidth[uiIdx] = uh
}

func (this *TComDataCU) GetHeight() []byte {
    return this.m_puhHeight
}
func (this *TComDataCU) GetHeight1(uiIdx uint) byte {
    return this.m_puhHeight[uiIdx]
}
func (this *TComDataCU) SetHeight(uiIdx uint, uh byte) {
    this.m_puhHeight[uiIdx] = uh
}

func (this *TComDataCU) SetSizeSubParts(uiWidth, uiHeight, uiAbsPartIdx, uiDepth uint) {
    uiCurrPartNumb := this.m_pcPic.GetNumPartInCU() >> (uiDepth << 1)
    for i := uint(0); i < uiCurrPartNumb; i++ {
        this.m_puhWidth[i+uiAbsPartIdx] = byte(uiWidth)
        this.m_puhHeight[i+uiAbsPartIdx] = byte(uiHeight)
    }
    //memset( m_puhWidth  + uiAbsPartIdx, uiWidth,  sizeof(UChar)*uiCurrPartNumb );
    //memset( m_puhHeight + uiAbsPartIdx, uiHeight, sizeof(UChar)*uiCurrPartNumb );
}
func (this *TComDataCU) GetQP() []int8 {
    return this.m_phQP
}
func (this *TComDataCU) GetQP1(uiIdx uint) int8 {
    return this.m_phQP[uiIdx]
}
func (this *TComDataCU) SetQP(uiIdx int, value int8) {
    this.m_phQP[uiIdx] = value
}
func (this *TComDataCU) SetQPSubParts(qp int, uiAbsPartIdx, uiDepth uint) {
    uiCurrPartNumb := this.m_pcPic.GetNumPartInCU() >> (uiDepth << 1)
    pcSlice := this.GetPic().GetSlice(this.GetPic().GetCurrSliceIdx())

    for uiSCUIdx := uiAbsPartIdx; uiSCUIdx < uiAbsPartIdx+uiCurrPartNumb; uiSCUIdx++ {
        if this.m_pcPic.GetCU(this.GetAddr()).GetSliceSegmentStartCU(uiSCUIdx+this.GetZorderIdxInCU()) == pcSlice.GetSliceSegmentCurStartCUAddr() {
            this.m_phQP[uiSCUIdx] = int8(qp)
        }
    }
}
func (this *TComDataCU) GetLastValidPartIdx(iAbsPartIdx int) int {
    iLastValidPartIdx := iAbsPartIdx - 1
    for iLastValidPartIdx >= 0 && this.GetPredictionMode1(uint(iLastValidPartIdx)) == MODE_NONE {
        uiDepth := this.GetDepth1(uint(iLastValidPartIdx))
        iLastValidPartIdx -= int(this.m_uiNumPartition >> (uiDepth << 1))
    }
    return iLastValidPartIdx
}
func (this *TComDataCU) GetLastCodedQP(uiAbsPartIdx uint) int8 {
    var uiQUPartIdxMask uint
    uiQUPartIdxMask = ^((1 << ((this.GetSlice().GetSPS().GetMaxCUDepth() - this.GetSlice().GetPPS().GetMaxCuDQPDepth()) << 1)) - 1)
    iLastValidPartIdx := this.GetLastValidPartIdx(int(uiAbsPartIdx & uiQUPartIdxMask))

    if uiAbsPartIdx < this.m_uiNumPartition && (this.GetSCUAddr()+uint(iLastValidPartIdx) < this.GetSliceStartCU(this.m_uiAbsIdxInLCU+uiAbsPartIdx)) {
        return int8(this.GetSlice().GetSliceQp())
    } else if iLastValidPartIdx >= 0 {
        return this.GetQP1(uint(iLastValidPartIdx))
    } else if this.GetZorderIdxInCU() > 0 {
        return this.GetPic().GetCU(this.GetAddr()).GetLastCodedQP(this.GetZorderIdxInCU())
    } else if this.GetPic().GetPicSym().GetInverseCUOrderMap(int(this.GetAddr())) > 0 &&
        this.GetPic().GetPicSym().GetTileIdxMap(int(this.GetAddr())) == this.GetPic().GetPicSym().GetTileIdxMap(int(this.GetPic().GetPicSym().GetCUOrderMap(int(this.GetPic().GetPicSym().GetInverseCUOrderMap(int(this.GetAddr())))-1))) &&
        !(this.GetSlice().GetPPS().GetEntropyCodingSyncEnabledFlag() && this.GetAddr()%this.GetPic().GetFrameWidthInCU() == 0) {
        return this.GetPic().GetCU(this.GetPic().GetPicSym().GetCUOrderMap(int(this.GetPic().GetPicSym().GetInverseCUOrderMap(int(this.GetAddr())) - 1))).GetLastCodedQP(this.GetPic().GetNumPartInCU())
    }

    return int8(this.GetSlice().GetSliceQp())
}
func (this *TComDataCU) SetQPSubCUs(qp int, pcCU *TComDataCU, absPartIdx, depth uint, foundNonZeroCbf *bool) {
    currPartNumb := this.m_pcPic.GetNumPartInCU() >> (depth << 1)
    currPartNumQ := currPartNumb >> 2

    if !*foundNonZeroCbf {
        if uint(pcCU.GetDepth1(absPartIdx)) > depth {
            for partUnitIdx := uint(0); partUnitIdx < 4; partUnitIdx++ {
                pcCU.SetQPSubCUs(qp, pcCU, absPartIdx+partUnitIdx*currPartNumQ, depth+1, foundNonZeroCbf)
            }
        } else {
            if pcCU.GetCbf2(absPartIdx, TEXT_LUMA) != 0 || pcCU.GetCbf2(absPartIdx, TEXT_CHROMA_U) != 0 || pcCU.GetCbf2(absPartIdx, TEXT_CHROMA_V) != 0 {
                *foundNonZeroCbf = true
            } else {
                this.SetQPSubParts(qp, absPartIdx, depth)
            }
        }
    }
}
func (this *TComDataCU) SetCodedQP(qp int8) {
    this.m_codedQP = qp
}
func (this *TComDataCU) GetCodedQP() int8 {
    return this.m_codedQP
}

func (this *TComDataCU) IsLosslessCoded(absPartIdx uint) bool {
    return (this.GetSlice().GetPPS().GetTransquantBypassEnableFlag() && this.GetCUTransquantBypass1(absPartIdx))
}

func (this *TComDataCU) GetTransformIdx() []byte {
    return this.m_puhTrIdx
}
func (this *TComDataCU) GetTransformIdx1(uiIdx uint) byte {
    return this.m_puhTrIdx[uiIdx]
}
func (this *TComDataCU) SetTrIdxSubParts(uiTrIdx, uiAbsPartIdx, uiDepth uint) {
    uiCurrPartNumb := this.m_pcPic.GetNumPartInCU() >> (uiDepth << 1)
    for i := uint(0); i < uiCurrPartNumb; i++ {
        this.m_puhTrIdx[i+uiAbsPartIdx] = byte(uiTrIdx)
    }
    //memset( m_puhTrIdx + uiAbsPartIdx, uiTrIdx, sizeof(UChar)*uiCurrPartNumb );
}

func (this *TComDataCU) GetTransformSkip1(eType TextType) []bool {
    return this.m_puhTransformSkip[G_aucConvertTxtTypeToIdx[eType]]
}
func (this *TComDataCU) GetTransformSkip2(uiIdx uint, eType TextType) bool {
    return this.m_puhTransformSkip[G_aucConvertTxtTypeToIdx[eType]][uiIdx]
}
func (this *TComDataCU) SetTransformSkipSubParts4(useTransformSkip bool, eType TextType, uiAbsPartIdx, uiDepth uint) {
    uiCurrPartNumb := this.m_pcPic.GetNumPartInCU() >> (uiDepth << 1)
    for i := uint(0); i < uiCurrPartNumb; i++ {
        this.m_puhTransformSkip[G_aucConvertTxtTypeToIdx[eType]][i+uiAbsPartIdx] = useTransformSkip
    }
    //memset( m_puhTransformSkip[g_aucConvertTxtTypeToIdx[eType]] + uiAbsPartIdx, useTransformSkip, sizeof( UChar ) * uiCurrPartNumb );
}
func (this *TComDataCU) SetTransformSkipSubParts5(useTransformSkipY, useTransformSkipU, useTransformSkipV bool, uiAbsPartIdx, uiDepth uint) {
    uiCurrPartNumb := this.m_pcPic.GetNumPartInCU() >> (uiDepth << 1)
    for i := uint(0); i < uiCurrPartNumb; i++ {
        this.m_puhTransformSkip[0][i+uiAbsPartIdx] = (useTransformSkipY) //, sizeof( UChar ) * uiCurrPartNumb );
        this.m_puhTransformSkip[1][i+uiAbsPartIdx] = (useTransformSkipU) //, sizeof( UChar ) * uiCurrPartNumb );
        this.m_puhTransformSkip[2][i+uiAbsPartIdx] = (useTransformSkipV) //, sizeof( UChar ) * uiCurrPartNumb );
    }

}

func (this *TComDataCU) GetQuadtreeTULog2MinSizeInCU(absPartIdx uint) uint {
    log2CbSize := uint(G_aucConvertToBit[this.GetWidth1(absPartIdx)] + 2)
    partSize := this.GetPartitionSize1(absPartIdx)
    var quadtreeTUMaxDepth uint
    if this.GetPredictionMode1(absPartIdx) == MODE_INTRA {
        quadtreeTUMaxDepth = this.m_pcSlice.GetSPS().GetQuadtreeTUMaxDepthIntra()
    } else {
        quadtreeTUMaxDepth = this.m_pcSlice.GetSPS().GetQuadtreeTUMaxDepthInter()
    }
    var intraSplitFlag uint
    if this.GetPredictionMode1(absPartIdx) == MODE_INTRA && partSize == SIZE_NxN {
        intraSplitFlag = 1
    } else {
        intraSplitFlag = 0
    }

    interSplitFlag := uint(B2U((quadtreeTUMaxDepth == 1) && (this.GetPredictionMode1(absPartIdx) == MODE_INTER) && (partSize != SIZE_2Nx2N)))

    log2MinTUSizeInCU := uint(0)
    if log2CbSize < (uint(this.m_pcSlice.GetSPS().GetQuadtreeTULog2MinSize()) + quadtreeTUMaxDepth - 1 + interSplitFlag + intraSplitFlag) {
        // when fully making use of signaled TUMaxDepth + inter/intraSplitFlag, resulting luma TB size is < QuadtreeTULog2MinSize
        log2MinTUSizeInCU = this.m_pcSlice.GetSPS().GetQuadtreeTULog2MinSize()
    } else {
        // when fully making use of signaled TUMaxDepth + inter/intraSplitFlag, resulting luma TB size is still >= QuadtreeTULog2MinSize
        log2MinTUSizeInCU = log2CbSize - (quadtreeTUMaxDepth - 1 + interSplitFlag + intraSplitFlag) // stop when trafoDepth == hierarchy_depth = splitFlag
        if log2MinTUSizeInCU > this.m_pcSlice.GetSPS().GetQuadtreeTULog2MaxSize() {
            // when fully making use of signaled TUMaxDepth + inter/intraSplitFlag, resulting luma TB size is still > QuadtreeTULog2MaxSize
            log2MinTUSizeInCU = this.m_pcSlice.GetSPS().GetQuadtreeTULog2MaxSize()
        }
    }
    return log2MinTUSizeInCU
}

func (this *TComDataCU) GetCUMvField(e RefPicList) *TComCUMvField {
    return &this.m_acCUMvField[e]
}

func (this *TComDataCU) GetCoeffY() []TCoeff {
    return this.m_pcTrCoeffY
}
func (this *TComDataCU) GetCoeffCb() []TCoeff {
    return this.m_pcTrCoeffCb
}
func (this *TComDataCU) GetCoeffCr() []TCoeff {
    return this.m_pcTrCoeffCr
}

//#if ADAPTIVE_QP_SELECTION
func (this *TComDataCU) GetArlCoeffY() []TCoeff {
    return this.m_pcArlCoeffY
}
func (this *TComDataCU) GetArlCoeffCb() []TCoeff {
    return this.m_pcArlCoeffCb
}
func (this *TComDataCU) GetArlCoeffCr() []TCoeff {
    return this.m_pcArlCoeffCr
}

//#endif

func (this *TComDataCU) GetPCMSampleY() []Pel {
    return this.m_pcIPCMSampleY
}
func (this *TComDataCU) GetPCMSampleCb() []Pel {
    return this.m_pcIPCMSampleCb
}
func (this *TComDataCU) GetPCMSampleCr() []Pel {
    return this.m_pcIPCMSampleCr
}

func (this *TComDataCU) GetCbf2(uiIdx uint, eType TextType) byte {
    return this.m_puhCbf[G_aucConvertTxtTypeToIdx[eType]][uiIdx]
}
func (this *TComDataCU) GetCbf1(eType TextType) []byte {
    return this.m_puhCbf[G_aucConvertTxtTypeToIdx[eType]]
}
func (this *TComDataCU) GetCbf3(uiIdx uint, eType TextType, uiTrDepth uint) byte {
    return ((this.GetCbf2(uiIdx, eType) >> uiTrDepth) & 0x1)
}
func (this *TComDataCU) SetCbf(uiIdx uint, eType TextType, uh byte) {
    this.m_puhCbf[G_aucConvertTxtTypeToIdx[eType]][uiIdx] = uh
}
func (this *TComDataCU) ClearCbf(uiIdx uint, eType TextType, uiNumParts uint) {
    for i := uint(0); i < uiNumParts; i++ {
        this.m_puhCbf[G_aucConvertTxtTypeToIdx[eType]][i+uiIdx] = 0
    }
    //::memset( &m_puhCbf[g_aucConvertTxtTypeToIdx[eType]][uiIdx], 0, sizeof(UChar)*uiNumParts);
}
func (this *TComDataCU) GetQtRootCbf(uiIdx uint) bool {
    return this.GetCbf3(uiIdx, TEXT_LUMA, 0) != 0 || this.GetCbf3(uiIdx, TEXT_CHROMA_U, 0) != 0 || this.GetCbf3(uiIdx, TEXT_CHROMA_V, 0) != 0
}

func (this *TComDataCU) SetCbfSubParts(uiCbfY, uiCbfU, uiCbfV byte, uiAbsPartIdx, uiDepth uint) {
    uiCurrPartNumb := this.m_pcPic.GetNumPartInCU() >> (uiDepth << 1)
    for i := uint(0); i < uiCurrPartNumb; i++ {
        this.m_puhCbf[0][i+uiAbsPartIdx] = byte(uiCbfY) //, sizeof( UChar ) * uiCurrPartNumb );
        this.m_puhCbf[1][i+uiAbsPartIdx] = byte(uiCbfU) //, sizeof( UChar ) * uiCurrPartNumb );
        this.m_puhCbf[2][i+uiAbsPartIdx] = byte(uiCbfV) //, sizeof( UChar ) * uiCurrPartNumb );
    }
}
func (this *TComDataCU) SetCbfSubParts4(uiCbf byte, eTType TextType, uiAbsPartIdx, uiDepth uint) {
    uiCurrPartNumb := this.m_pcPic.GetNumPartInCU() >> (uiDepth << 1)
    for i := uint(0); i < uiCurrPartNumb; i++ {
        this.m_puhCbf[G_aucConvertTxtTypeToIdx[eTType]][i+uiAbsPartIdx] = byte(uiCbf) //, sizeof( UChar ) * uiCurrPartNumb );
    }
}
func (this *TComDataCU) SetCbfSubParts5(uiCbf byte, eTType TextType, uiAbsPartIdx, uiPartIdx, uiDepth uint) {
    this.SetSubPartByte(uiCbf, this.m_puhCbf[G_aucConvertTxtTypeToIdx[eTType]], uiAbsPartIdx, uiDepth, uiPartIdx)
}

// -------------------------------------------------------------------------------------------------------------------
// member functions for coding tool information
// -------------------------------------------------------------------------------------------------------------------

func (this *TComDataCU) GetMergeFlag() []bool {
    return this.m_pbMergeFlag
}
func (this *TComDataCU) GetMergeFlag1(uiIdx uint) bool {
    return this.m_pbMergeFlag[uiIdx]
}
func (this *TComDataCU) SetMergeFlag(uiIdx uint, b bool) {
    this.m_pbMergeFlag[uiIdx] = b
}
func (this *TComDataCU) SetMergeFlagSubParts(bMergeFlag bool, uiAbsPartIdx, uiPartIdx, uiDepth uint) {
    this.SetSubPartBool(bMergeFlag, this.m_pbMergeFlag, uiAbsPartIdx, uiDepth, uiPartIdx)
}

func (this *TComDataCU) GetMergeIndex() []byte {
    return this.m_puhMergeIndex
}
func (this *TComDataCU) GetMergeIndex1(uiIdx uint) byte {
    return this.m_puhMergeIndex[uiIdx]
}
func (this *TComDataCU) SetMergeIndex(uiIdx uint, uiMergeIndex byte) {
    this.m_puhMergeIndex[uiIdx] = uiMergeIndex
}
func (this *TComDataCU) SetMergeIndexSubParts(uiMergeIndex, uiAbsPartIdx, uiPartIdx, uiDepth uint) {
    this.SetSubPartByte(byte(uiMergeIndex), this.m_puhMergeIndex, uiAbsPartIdx, uiDepth, uiPartIdx)
}

//  template <typename T>
func (this *TComDataCU) SetSubPartByte(uiParameter byte, puhBaseLCU []byte, uiCUAddr, uiCUDepth, uiPUIdx uint) {
    //assert( sizeof(T) == 1 ); // Using memset() works only for types of size 1

    uiCurrPartNumQ := (this.m_pcPic.GetNumPartInCU() >> (2 * uiCUDepth)) >> 2
    switch this.m_pePartSize[uiCUAddr] {
    case SIZE_2Nx2N:
        for i := uint(0); i < 4*uiCurrPartNumQ; i++ {
            puhBaseLCU[i+uiCUAddr] = uiParameter //, 4 * uiCurrPartNumQ );
        }
    case SIZE_2NxN:
        for i := uint(0); i < 2*uiCurrPartNumQ; i++ {
            puhBaseLCU[i+uiCUAddr] = uiParameter //, 2 * uiCurrPartNumQ );
        }
    case SIZE_Nx2N:
        for i := uint(0); i < uiCurrPartNumQ; i++ {
            puhBaseLCU[i+uiCUAddr] = uiParameter                  //, uiCurrPartNumQ );
            puhBaseLCU[i+uiCUAddr+2*uiCurrPartNumQ] = uiParameter //, uiCurrPartNumQ );
        }
    case SIZE_NxN:
        for i := uint(0); i < uiCurrPartNumQ; i++ {
            puhBaseLCU[i+uiCUAddr] = uiParameter //, uiCurrPartNumQ );
        }
    case SIZE_2NxnU:
        if uiPUIdx == 0 {
            for i := uint(0); i < (uiCurrPartNumQ >> 1); i++ {
                puhBaseLCU[i+uiCUAddr] = uiParameter                //, (uiCurrPartNumQ >> 1) );
                puhBaseLCU[i+uiCUAddr+uiCurrPartNumQ] = uiParameter //, (uiCurrPartNumQ >> 1) );
            }
        } else if uiPUIdx == 1 {
            for i := uint(0); i < (uiCurrPartNumQ >> 1); i++ {
                puhBaseLCU[i+uiCUAddr] = uiParameter //, (uiCurrPartNumQ >> 1) );
            }
            for i := uint(0); i < ((uiCurrPartNumQ >> 1) + (uiCurrPartNumQ << 1)); i++ {
                puhBaseLCU[i+uiCUAddr+uiCurrPartNumQ] = uiParameter //, ((uiCurrPartNumQ >> 1) + (uiCurrPartNumQ << 1)) );
            }
        } else {
            //assert(0);
        }

    case SIZE_2NxnD:
        if uiPUIdx == 0 {
            for i := uint(0); i < ((uiCurrPartNumQ << 1) + (uiCurrPartNumQ >> 1)); i++ {
                puhBaseLCU[i+uiCUAddr] = uiParameter //, ((uiCurrPartNumQ << 1) + (uiCurrPartNumQ >> 1)) );
            }
            for i := uint(0); i < (uiCurrPartNumQ >> 1); i++ {
                puhBaseLCU[i+uiCUAddr+(uiCurrPartNumQ<<1)+uiCurrPartNumQ] = uiParameter //, (uiCurrPartNumQ >> 1) );
            }
        } else if uiPUIdx == 1 {
            for i := uint(0); i < (uiCurrPartNumQ >> 1); i++ {
                puhBaseLCU[i+uiCUAddr] = uiParameter                //, (uiCurrPartNumQ >> 1) );
                puhBaseLCU[i+uiCUAddr+uiCurrPartNumQ] = uiParameter //, (uiCurrPartNumQ >> 1) );
            }
        } else {
            //assert(0);
        }

    case SIZE_nLx2N:
        if uiPUIdx == 0 {
            for i := uint(0); i < (uiCurrPartNumQ >> 2); i++ {
                puhBaseLCU[i+uiCUAddr] = uiParameter                                         //, (uiCurrPartNumQ >> 2) );
                puhBaseLCU[i+uiCUAddr+(uiCurrPartNumQ>>1)] = uiParameter                     //, (uiCurrPartNumQ >> 2) );
                puhBaseLCU[i+uiCUAddr+(uiCurrPartNumQ<<1)] = uiParameter                     //, (uiCurrPartNumQ >> 2) );
                puhBaseLCU[i+uiCUAddr+(uiCurrPartNumQ<<1)+(uiCurrPartNumQ>>1)] = uiParameter //, (uiCurrPartNumQ >> 2) );
            }
        } else if uiPUIdx == 1 {
            for i := uint(0); i < (uiCurrPartNumQ >> 2); i++ {
                puhBaseLCU[i+uiCUAddr] = uiParameter                     //, (uiCurrPartNumQ >> 2) );
                puhBaseLCU[i+uiCUAddr+(uiCurrPartNumQ<<1)] = uiParameter //, (uiCurrPartNumQ >> 2) );

            }
            for i := uint(0); i < (uiCurrPartNumQ + (uiCurrPartNumQ >> 2)); i++ {
                puhBaseLCU[i+uiCUAddr+(uiCurrPartNumQ>>1)] = uiParameter                     //, (uiCurrPartNumQ + (uiCurrPartNumQ >> 2)) );
                puhBaseLCU[i+uiCUAddr+(uiCurrPartNumQ<<1)+(uiCurrPartNumQ>>1)] = uiParameter //, (uiCurrPartNumQ + (uiCurrPartNumQ >> 2)) );
            }
        } else {
            //assert(0);
        }

    case SIZE_nRx2N:
        if uiPUIdx == 0 {
            for i := uint(0); i < (uiCurrPartNumQ >> 2); i++ {
                puhBaseLCU[i+uiCUAddr+uiCurrPartNumQ+(uiCurrPartNumQ>>1)] = uiParameter                     //, (uiCurrPartNumQ >> 2) );
                puhBaseLCU[i+uiCUAddr+(uiCurrPartNumQ<<1)+uiCurrPartNumQ+(uiCurrPartNumQ>>1)] = uiParameter //, (uiCurrPartNumQ >> 2) );
            }
            for i := uint(0); i < (uiCurrPartNumQ + (uiCurrPartNumQ >> 2)); i++ {
                puhBaseLCU[i+uiCUAddr] = uiParameter                     //, (uiCurrPartNumQ + (uiCurrPartNumQ >> 2)) );
                puhBaseLCU[i+uiCUAddr+(uiCurrPartNumQ<<1)] = uiParameter //, (uiCurrPartNumQ + (uiCurrPartNumQ >> 2)) );
            }
        } else if uiPUIdx == 1 {
            for i := uint(0); i < (uiCurrPartNumQ >> 2); i++ {
                puhBaseLCU[i+uiCUAddr] = uiParameter                                         // (uiCurrPartNumQ >> 2) );
                puhBaseLCU[i+uiCUAddr+(uiCurrPartNumQ>>1)] = uiParameter                     // (uiCurrPartNumQ >> 2) );
                puhBaseLCU[i+uiCUAddr+(uiCurrPartNumQ<<1)] = uiParameter                     // (uiCurrPartNumQ >> 2) );
                puhBaseLCU[i+uiCUAddr+(uiCurrPartNumQ<<1)+(uiCurrPartNumQ>>1)] = uiParameter // (uiCurrPartNumQ >> 2) );
            }
        } else {
            //assert(0);
        }

    default:
        //assert( 0 );
    }
}

func (this *TComDataCU) SetSubPartInt8(uiParameter int8, puhBaseLCU []int8, uiCUAddr, uiCUDepth, uiPUIdx uint) {
    //assert( sizeof(T) == 1 ); // Using memset() works only for types of size 1

    uiCurrPartNumQ := (this.m_pcPic.GetNumPartInCU() >> (2 * uiCUDepth)) >> 2
    switch this.m_pePartSize[uiCUAddr] {
    case SIZE_2Nx2N:
        for i := uint(0); i < 4*uiCurrPartNumQ; i++ {
            puhBaseLCU[i+uiCUAddr] = uiParameter //, 4 * uiCurrPartNumQ );
        }
    case SIZE_2NxN:
        for i := uint(0); i < 2*uiCurrPartNumQ; i++ {
            puhBaseLCU[i+uiCUAddr] = uiParameter //, 2 * uiCurrPartNumQ );
        }
    case SIZE_Nx2N:
        for i := uint(0); i < uiCurrPartNumQ; i++ {
            puhBaseLCU[i+uiCUAddr] = uiParameter                  //, uiCurrPartNumQ );
            puhBaseLCU[i+uiCUAddr+2*uiCurrPartNumQ] = uiParameter //, uiCurrPartNumQ );
        }
    case SIZE_NxN:
        for i := uint(0); i < uiCurrPartNumQ; i++ {
            puhBaseLCU[i+uiCUAddr] = uiParameter //, uiCurrPartNumQ );
        }
    case SIZE_2NxnU:
        if uiPUIdx == 0 {
            for i := uint(0); i < (uiCurrPartNumQ >> 1); i++ {
                puhBaseLCU[i+uiCUAddr] = uiParameter                //, (uiCurrPartNumQ >> 1) );
                puhBaseLCU[i+uiCUAddr+uiCurrPartNumQ] = uiParameter //, (uiCurrPartNumQ >> 1) );
            }
        } else if uiPUIdx == 1 {
            for i := uint(0); i < (uiCurrPartNumQ >> 1); i++ {
                puhBaseLCU[i+uiCUAddr] = uiParameter //, (uiCurrPartNumQ >> 1) );
            }
            for i := uint(0); i < ((uiCurrPartNumQ >> 1) + (uiCurrPartNumQ << 1)); i++ {
                puhBaseLCU[i+uiCUAddr+uiCurrPartNumQ] = uiParameter //, ((uiCurrPartNumQ >> 1) + (uiCurrPartNumQ << 1)) );
            }
        } else {
            //assert(0);
        }

    case SIZE_2NxnD:
        if uiPUIdx == 0 {
            for i := uint(0); i < ((uiCurrPartNumQ << 1) + (uiCurrPartNumQ >> 1)); i++ {
                puhBaseLCU[i+uiCUAddr] = uiParameter //, ((uiCurrPartNumQ << 1) + (uiCurrPartNumQ >> 1)) );
            }
            for i := uint(0); i < (uiCurrPartNumQ >> 1); i++ {
                puhBaseLCU[i+uiCUAddr+(uiCurrPartNumQ<<1)+uiCurrPartNumQ] = uiParameter //, (uiCurrPartNumQ >> 1) );
            }
        } else if uiPUIdx == 1 {
            for i := uint(0); i < (uiCurrPartNumQ >> 1); i++ {
                puhBaseLCU[i+uiCUAddr] = uiParameter                //, (uiCurrPartNumQ >> 1) );
                puhBaseLCU[i+uiCUAddr+uiCurrPartNumQ] = uiParameter //, (uiCurrPartNumQ >> 1) );
            }
        } else {
            //assert(0);
        }

    case SIZE_nLx2N:
        if uiPUIdx == 0 {
            for i := uint(0); i < (uiCurrPartNumQ >> 2); i++ {
                puhBaseLCU[i+uiCUAddr] = uiParameter                                         //, (uiCurrPartNumQ >> 2) );
                puhBaseLCU[i+uiCUAddr+(uiCurrPartNumQ>>1)] = uiParameter                     //, (uiCurrPartNumQ >> 2) );
                puhBaseLCU[i+uiCUAddr+(uiCurrPartNumQ<<1)] = uiParameter                     //, (uiCurrPartNumQ >> 2) );
                puhBaseLCU[i+uiCUAddr+(uiCurrPartNumQ<<1)+(uiCurrPartNumQ>>1)] = uiParameter //, (uiCurrPartNumQ >> 2) );
            }
        } else if uiPUIdx == 1 {
            for i := uint(0); i < (uiCurrPartNumQ >> 2); i++ {
                puhBaseLCU[i+uiCUAddr] = uiParameter                     //, (uiCurrPartNumQ >> 2) );
                puhBaseLCU[i+uiCUAddr+(uiCurrPartNumQ<<1)] = uiParameter //, (uiCurrPartNumQ >> 2) );

            }
            for i := uint(0); i < (uiCurrPartNumQ + (uiCurrPartNumQ >> 2)); i++ {
                puhBaseLCU[i+uiCUAddr+(uiCurrPartNumQ>>1)] = uiParameter                     //, (uiCurrPartNumQ + (uiCurrPartNumQ >> 2)) );
                puhBaseLCU[i+uiCUAddr+(uiCurrPartNumQ<<1)+(uiCurrPartNumQ>>1)] = uiParameter //, (uiCurrPartNumQ + (uiCurrPartNumQ >> 2)) );
            }
        } else {
            //assert(0);
        }

    case SIZE_nRx2N:
        if uiPUIdx == 0 {
            for i := uint(0); i < (uiCurrPartNumQ >> 2); i++ {
                puhBaseLCU[i+uiCUAddr+uiCurrPartNumQ+(uiCurrPartNumQ>>1)] = uiParameter                     //, (uiCurrPartNumQ >> 2) );
                puhBaseLCU[i+uiCUAddr+(uiCurrPartNumQ<<1)+uiCurrPartNumQ+(uiCurrPartNumQ>>1)] = uiParameter //, (uiCurrPartNumQ >> 2) );
            }
            for i := uint(0); i < (uiCurrPartNumQ + (uiCurrPartNumQ >> 2)); i++ {
                puhBaseLCU[i+uiCUAddr] = uiParameter                     //, (uiCurrPartNumQ + (uiCurrPartNumQ >> 2)) );
                puhBaseLCU[i+uiCUAddr+(uiCurrPartNumQ<<1)] = uiParameter //, (uiCurrPartNumQ + (uiCurrPartNumQ >> 2)) );
            }
        } else if uiPUIdx == 1 {
            for i := uint(0); i < (uiCurrPartNumQ >> 2); i++ {
                puhBaseLCU[i+uiCUAddr] = uiParameter                                         // (uiCurrPartNumQ >> 2) );
                puhBaseLCU[i+uiCUAddr+(uiCurrPartNumQ>>1)] = uiParameter                     // (uiCurrPartNumQ >> 2) );
                puhBaseLCU[i+uiCUAddr+(uiCurrPartNumQ<<1)] = uiParameter                     // (uiCurrPartNumQ >> 2) );
                puhBaseLCU[i+uiCUAddr+(uiCurrPartNumQ<<1)+(uiCurrPartNumQ>>1)] = uiParameter // (uiCurrPartNumQ >> 2) );
            }
        } else {
            //assert(0);
        }

    default:
        //assert( 0 );
    }
}
func (this *TComDataCU) SetSubPartBool(uiParameter bool, puhBaseLCU []bool, uiCUAddr, uiCUDepth, uiPUIdx uint) {
    //assert( sizeof(T) == 1 ); // Using memset() works only for types of size 1

    uiCurrPartNumQ := (this.m_pcPic.GetNumPartInCU() >> (2 * uiCUDepth)) >> 2
    switch this.m_pePartSize[uiCUAddr] {
    case SIZE_2Nx2N:
        for i := uint(0); i < 4*uiCurrPartNumQ; i++ {
            puhBaseLCU[i+uiCUAddr] = uiParameter //, 4 * uiCurrPartNumQ );
        }
    case SIZE_2NxN:
        for i := uint(0); i < 2*uiCurrPartNumQ; i++ {
            puhBaseLCU[i+uiCUAddr] = uiParameter //, 2 * uiCurrPartNumQ );
        }
    case SIZE_Nx2N:
        for i := uint(0); i < uiCurrPartNumQ; i++ {
            puhBaseLCU[i+uiCUAddr] = uiParameter                  //, uiCurrPartNumQ );
            puhBaseLCU[i+uiCUAddr+2*uiCurrPartNumQ] = uiParameter //, uiCurrPartNumQ );
        }
    case SIZE_NxN:
        for i := uint(0); i < uiCurrPartNumQ; i++ {
            puhBaseLCU[i+uiCUAddr] = uiParameter //, uiCurrPartNumQ );
        }
    case SIZE_2NxnU:
        if uiPUIdx == 0 {
            for i := uint(0); i < (uiCurrPartNumQ >> 1); i++ {
                puhBaseLCU[i+uiCUAddr] = uiParameter                //, (uiCurrPartNumQ >> 1) );
                puhBaseLCU[i+uiCUAddr+uiCurrPartNumQ] = uiParameter //, (uiCurrPartNumQ >> 1) );
            }
        } else if uiPUIdx == 1 {
            for i := uint(0); i < (uiCurrPartNumQ >> 1); i++ {
                puhBaseLCU[i+uiCUAddr] = uiParameter //, (uiCurrPartNumQ >> 1) );
            }
            for i := uint(0); i < ((uiCurrPartNumQ >> 1) + (uiCurrPartNumQ << 1)); i++ {
                puhBaseLCU[i+uiCUAddr+uiCurrPartNumQ] = uiParameter //, ((uiCurrPartNumQ >> 1) + (uiCurrPartNumQ << 1)) );
            }
        } else {
            //assert(0);
        }

    case SIZE_2NxnD:
        if uiPUIdx == 0 {
            for i := uint(0); i < ((uiCurrPartNumQ << 1) + (uiCurrPartNumQ >> 1)); i++ {
                puhBaseLCU[i+uiCUAddr] = uiParameter //, ((uiCurrPartNumQ << 1) + (uiCurrPartNumQ >> 1)) );
            }
            for i := uint(0); i < (uiCurrPartNumQ >> 1); i++ {
                puhBaseLCU[i+uiCUAddr+(uiCurrPartNumQ<<1)+uiCurrPartNumQ] = uiParameter //, (uiCurrPartNumQ >> 1) );
            }
        } else if uiPUIdx == 1 {
            for i := uint(0); i < (uiCurrPartNumQ >> 1); i++ {
                puhBaseLCU[i+uiCUAddr] = uiParameter                //, (uiCurrPartNumQ >> 1) );
                puhBaseLCU[i+uiCUAddr+uiCurrPartNumQ] = uiParameter //, (uiCurrPartNumQ >> 1) );
            }
        } else {
            //assert(0);
        }

    case SIZE_nLx2N:
        if uiPUIdx == 0 {
            for i := uint(0); i < (uiCurrPartNumQ >> 2); i++ {
                puhBaseLCU[i+uiCUAddr] = uiParameter                                         //, (uiCurrPartNumQ >> 2) );
                puhBaseLCU[i+uiCUAddr+(uiCurrPartNumQ>>1)] = uiParameter                     //, (uiCurrPartNumQ >> 2) );
                puhBaseLCU[i+uiCUAddr+(uiCurrPartNumQ<<1)] = uiParameter                     //, (uiCurrPartNumQ >> 2) );
                puhBaseLCU[i+uiCUAddr+(uiCurrPartNumQ<<1)+(uiCurrPartNumQ>>1)] = uiParameter //, (uiCurrPartNumQ >> 2) );
            }
        } else if uiPUIdx == 1 {
            for i := uint(0); i < (uiCurrPartNumQ >> 2); i++ {
                puhBaseLCU[i+uiCUAddr] = uiParameter                     //, (uiCurrPartNumQ >> 2) );
                puhBaseLCU[i+uiCUAddr+(uiCurrPartNumQ<<1)] = uiParameter //, (uiCurrPartNumQ >> 2) );

            }
            for i := uint(0); i < (uiCurrPartNumQ + (uiCurrPartNumQ >> 2)); i++ {
                puhBaseLCU[i+uiCUAddr+(uiCurrPartNumQ>>1)] = uiParameter                     //, (uiCurrPartNumQ + (uiCurrPartNumQ >> 2)) );
                puhBaseLCU[i+uiCUAddr+(uiCurrPartNumQ<<1)+(uiCurrPartNumQ>>1)] = uiParameter //, (uiCurrPartNumQ + (uiCurrPartNumQ >> 2)) );
            }
        } else {
            //assert(0);
        }

    case SIZE_nRx2N:
        if uiPUIdx == 0 {
            for i := uint(0); i < (uiCurrPartNumQ >> 2); i++ {
                puhBaseLCU[i+uiCUAddr+uiCurrPartNumQ+(uiCurrPartNumQ>>1)] = uiParameter                     //, (uiCurrPartNumQ >> 2) );
                puhBaseLCU[i+uiCUAddr+(uiCurrPartNumQ<<1)+uiCurrPartNumQ+(uiCurrPartNumQ>>1)] = uiParameter //, (uiCurrPartNumQ >> 2) );
            }
            for i := uint(0); i < (uiCurrPartNumQ + (uiCurrPartNumQ >> 2)); i++ {
                puhBaseLCU[i+uiCUAddr] = uiParameter                     //, (uiCurrPartNumQ + (uiCurrPartNumQ >> 2)) );
                puhBaseLCU[i+uiCUAddr+(uiCurrPartNumQ<<1)] = uiParameter //, (uiCurrPartNumQ + (uiCurrPartNumQ >> 2)) );
            }
        } else if uiPUIdx == 1 {
            for i := uint(0); i < (uiCurrPartNumQ >> 2); i++ {
                puhBaseLCU[i+uiCUAddr] = uiParameter                                         // (uiCurrPartNumQ >> 2) );
                puhBaseLCU[i+uiCUAddr+(uiCurrPartNumQ>>1)] = uiParameter                     // (uiCurrPartNumQ >> 2) );
                puhBaseLCU[i+uiCUAddr+(uiCurrPartNumQ<<1)] = uiParameter                     // (uiCurrPartNumQ >> 2) );
                puhBaseLCU[i+uiCUAddr+(uiCurrPartNumQ<<1)+(uiCurrPartNumQ>>1)] = uiParameter // (uiCurrPartNumQ >> 2) );
            }
        } else {
            //assert(0);
        }

    default:
        //assert( 0 );
    }
}

//#if AMP_MRG
func (this *TComDataCU) SetMergeAMP(b bool) {
    this.m_bIsMergeAMP = b
}
func (this *TComDataCU) GetMergeAMP() bool {
    return this.m_bIsMergeAMP
}

//#endif

func (this *TComDataCU) GetLumaIntraDir() []byte {
    return this.m_puhLumaIntraDir
}
func (this *TComDataCU) GetLumaIntraDir1(uiIdx uint) byte {
    return this.m_puhLumaIntraDir[uiIdx]
}
func (this *TComDataCU) SetLumaIntraDir(uiIdx uint, uh byte) {
    this.m_puhLumaIntraDir[uiIdx] = uh
}
func (this *TComDataCU) SetLumaIntraDirSubParts(uiDir, uiAbsPartIdx, uiDepth uint) {
    uiCurrPartNumb := this.m_pcPic.GetNumPartInCU() >> (uiDepth << 1)
    for i := uint(0); i < uiCurrPartNumb; i++ {
        this.m_puhLumaIntraDir[i+uiAbsPartIdx] = byte(uiDir)
    }
    //memset( m_puhLumaIntraDir + uiAbsPartIdx, uiDir, sizeof(UChar)*uiCurrPartNumb );
}

func (this *TComDataCU) GetChromaIntraDir() []byte {
    return this.m_puhChromaIntraDir
}
func (this *TComDataCU) GetChromaIntraDir1(uiIdx uint) byte {
    return this.m_puhChromaIntraDir[uiIdx]
}
func (this *TComDataCU) SetChromaIntraDir(uiIdx uint, uh byte) {
    this.m_puhChromaIntraDir[uiIdx] = uh
}
func (this *TComDataCU) SetChromIntraDirSubParts(uiDir, uiAbsPartIdx, uiDepth uint) {
    uiCurrPartNumb := this.m_pcPic.GetNumPartInCU() >> (uiDepth << 1)
    for i := uint(0); i < uiCurrPartNumb; i++ {
        this.m_puhChromaIntraDir[i+uiAbsPartIdx] = byte(uiDir)
    }
    //memset( m_puhChromaIntraDir + uiAbsPartIdx, uiDir, sizeof(UChar)*uiCurrPartNumb );
}

func (this *TComDataCU) GetInterDir() []byte {
    return this.m_puhInterDir
}
func (this *TComDataCU) GetInterDir1(uiIdx uint) byte {
    return this.m_puhInterDir[uiIdx]
}
func (this *TComDataCU) SetInterDir(uiIdx uint, uh byte) {
    this.m_puhInterDir[uiIdx] = uh
}
func (this *TComDataCU) SetInterDirSubParts(uiDir, uiAbsPartIdx, uiPartIdx, uiDepth uint) {
    this.SetSubPartByte(byte(uiDir), this.m_puhInterDir, uiAbsPartIdx, uiDepth, uiPartIdx)
}
func (this *TComDataCU) GetIPCMFlag() []bool {
    return this.m_pbIPCMFlag
}
func (this *TComDataCU) GetIPCMFlag1(uiIdx uint) bool {
    return this.m_pbIPCMFlag[uiIdx]
}
func (this *TComDataCU) SetIPCMFlag(uiIdx uint, b bool) {
    this.m_pbIPCMFlag[uiIdx] = b
}
func (this *TComDataCU) SetIPCMFlagSubParts(bIpcmFlag bool, uiAbsPartIdx, uiDepth uint) {
    uiCurrPartNumb := this.m_pcPic.GetNumPartInCU() >> (uiDepth << 1)
    for i := uint(0); i < uiCurrPartNumb; i++ {
        this.m_pbIPCMFlag[i+uiAbsPartIdx] = bIpcmFlag
    }
    //memset(m_pbIPCMFlag + uiAbsPartIdx, bIpcmFlag, sizeof(Bool)*uiCurrPartNumb );
}

/// Get slice ID for SU
func (this *TComDataCU) GetSUSliceID(uiIdx uint) int {
    return this.m_piSliceSUMap[uiIdx]
}

/// Get the pointer of slice ID map
func (this *TComDataCU) GetSliceSUMap() ([]int, int) {
    return this.m_piSliceSUMap, this.m_iSliceSUMapAddr
}

/// Set the pointer of slice ID map
func (this *TComDataCU) SetSliceSUMap(pi []int, addr int) {
    this.m_piSliceSUMap = pi
    this.m_iSliceSUMapAddr = addr
}

func (this *TComDataCU) GetNDBFilterBlocks() *list.List {
    return this.m_vNDFBlock
}
func (this *TComDataCU) SetNDBFilterBlockBorderAvailability(numLCUInPicWidth, numLCUInPicHeight, numSUInLCUWidth, numSUInLCUHeight, picWidth, picHeight uint,
    LFCrossSliceBoundary map[int]bool, bTopTileBoundary, bDownTileBoundary, bLeftTileBoundary, bRightTileBoundary, bIndependentTileBoundaryEnabled bool) {
    numSUInLCU := int(numSUInLCUWidth * numSUInLCUHeight)
    //pSliceIDMapLCU := this.m_piSliceSUMap;
    pSliceIDMapLCU := this.m_iSliceSUMapAddr
    onlyOneSliceInPic := len(LFCrossSliceBoundary) == 1
    var uiLPelX, uiTPelY, width, height uint
    var bPicRBoundary, bPicBBoundary, bPicTBoundary, bPicLBoundary bool
    bLCURBoundary := false
    bLCUBBoundary := false
    bLCUTBoundary := false
    bLCULBoundary := false
    var pbAvailBorder []bool
    var pbAvail *bool
    var rTLSU, rBRSU, widthSU, heightSU, zRefSU uint
    var pRefID *int
    rTRefSU := 0
    rBRefSU := 0
    rLRefSU := 0
    rRRefSU := 0
    var pRefMapLCU, pRRefMapLCU, pLRefMapLCU, pTRefMapLCU, pBRefMapLCU int
    var sliceID int
    //numSGU := uint(this.m_vNDFBlock.Len());

    //for i:=uint(0); i< numSGU; i++ {
    for e := this.m_vNDFBlock.Front(); e != nil; e = e.Next() {
        rSGU := e.Value.(*NDBFBlockInfo) //this.m_vNDFBlock[i];

        sliceID = rSGU.sliceID
        uiLPelX = rSGU.posX
        uiTPelY = rSGU.posY
        width = rSGU.width
        height = rSGU.height
        rTLSU = G_auiZscanToRaster[rSGU.startSU]
        rBRSU = G_auiZscanToRaster[rSGU.endSU]
        widthSU = rSGU.widthSU
        heightSU = rSGU.heightSU

        pbAvailBorder = rSGU.isBorderAvailable[:]

        bPicTBoundary = uiTPelY == 0
        bPicLBoundary = uiLPelX == 0
        bPicRBoundary = !(uiLPelX+width < picWidth)
        bPicBBoundary = !(uiTPelY+height < picHeight)

        bLCULBoundary = rTLSU%numSUInLCUWidth == 0
        bLCURBoundary = (rTLSU+widthSU)%numSUInLCUWidth == 0
        bLCUTBoundary = (rTLSU / numSUInLCUWidth) == 0
        bLCUBBoundary = (rBRSU / numSUInLCUWidth) == (numSUInLCUHeight - 1)

        //       SGU_L
        pbAvail = &(pbAvailBorder[SGU_L])
        if bPicLBoundary {
            *pbAvail = false
        } else if onlyOneSliceInPic {
            *pbAvail = true
        } else {
            //      bLCULBoundary = (rTLSU % uiNumSUInLCUWidth == 0)?(true):(false);
            if bLCULBoundary {
                rLRefSU = int(rTLSU + numSUInLCUWidth - 1)
                zRefSU = G_auiRasterToZscan[rLRefSU]
                pRefMapLCU = (pSliceIDMapLCU - numSUInLCU)
                pLRefMapLCU = (pSliceIDMapLCU - numSUInLCU)
            } else {
                zRefSU = G_auiRasterToZscan[rTLSU-1]
                pRefMapLCU = pSliceIDMapLCU
            }
            pRefID = &this.m_piSliceSUMap[pRefMapLCU+int(zRefSU)]
            if *pRefID == sliceID {
                *pbAvail = true
            } else if *pRefID > sliceID {
                *pbAvail = (LFCrossSliceBoundary[*pRefID])
            } else {
                *pbAvail = (LFCrossSliceBoundary[sliceID])
            }
        }

        //       SGU_R
        pbAvail = &(pbAvailBorder[SGU_R])
        if bPicRBoundary {
            *pbAvail = false
        } else if onlyOneSliceInPic {
            *pbAvail = true
        } else {
            //       bLCURBoundary = ( (rTLSU+ uiWidthSU) % uiNumSUInLCUWidth == 0)?(true):(false);
            if bLCURBoundary {
                rRRefSU = int(rTLSU + widthSU - numSUInLCUWidth)
                zRefSU = G_auiRasterToZscan[rRRefSU]
                pRefMapLCU = (pSliceIDMapLCU + numSUInLCU)
                pRRefMapLCU = (pSliceIDMapLCU + numSUInLCU)
            } else {
                zRefSU = G_auiRasterToZscan[rTLSU+widthSU]
                pRefMapLCU = pSliceIDMapLCU
            }
            pRefID = &this.m_piSliceSUMap[pRefMapLCU+int(zRefSU)]
            if *pRefID == sliceID {
                *pbAvail = true
            } else if *pRefID > sliceID {
                *pbAvail = (LFCrossSliceBoundary[*pRefID])
            } else {
                *pbAvail = (LFCrossSliceBoundary[sliceID])
            }
        }

        //       SGU_T
        pbAvail = &(pbAvailBorder[SGU_T])
        if bPicTBoundary {
            *pbAvail = false
        } else if onlyOneSliceInPic {
            *pbAvail = true
        } else {
            //      bLCUTBoundary = ( (UInt)(rTLSU / uiNumSUInLCUWidth)== 0)?(true):(false);
            if bLCUTBoundary {
                rTRefSU = numSUInLCU - int(numSUInLCUWidth-rTLSU)
                zRefSU = G_auiRasterToZscan[rTRefSU]
                pRefMapLCU = (pSliceIDMapLCU - int(numLCUInPicWidth)*numSUInLCU)
                pTRefMapLCU = (pSliceIDMapLCU - int(numLCUInPicWidth)*numSUInLCU)
            } else {
                zRefSU = G_auiRasterToZscan[rTLSU-numSUInLCUWidth]
                pRefMapLCU = pSliceIDMapLCU
            }
            pRefID = &this.m_piSliceSUMap[pRefMapLCU+int(zRefSU)]
            if *pRefID == sliceID {
                *pbAvail = (true)
            } else if *pRefID > sliceID {
                *pbAvail = (LFCrossSliceBoundary[*pRefID])
            } else {
                *pbAvail = (LFCrossSliceBoundary[sliceID])
            }
        }

        //       SGU_B
        pbAvail = &(pbAvailBorder[SGU_B])
        if bPicBBoundary {
            *pbAvail = false
        } else if onlyOneSliceInPic {
            *pbAvail = true
        } else {
            //      bLCUBBoundary = ( (UInt)(rBRSU / uiNumSUInLCUWidth) == (uiNumSUInLCUHeight-1) )?(true):(false);
            if bLCUBBoundary {
                rBRefSU = int(rTLSU) % int(numSUInLCUWidth)
                zRefSU = G_auiRasterToZscan[rBRefSU]
                pRefMapLCU = (pSliceIDMapLCU + int(numLCUInPicWidth)*numSUInLCU)
                pBRefMapLCU = (pSliceIDMapLCU + int(numLCUInPicWidth)*numSUInLCU)
            } else {
                zRefSU = G_auiRasterToZscan[rTLSU+(heightSU*numSUInLCUWidth)]
                pRefMapLCU = pSliceIDMapLCU
            }
            pRefID = &this.m_piSliceSUMap[pRefMapLCU+int(zRefSU)]
            if *pRefID == sliceID {
                *pbAvail = (true)
            } else if *pRefID > sliceID {
                *pbAvail = (LFCrossSliceBoundary[*pRefID])
            } else {
                *pbAvail = (LFCrossSliceBoundary[sliceID])
            }
        }

        //       SGU_TL
        pbAvail = &(pbAvailBorder[SGU_TL])
        if bPicTBoundary || bPicLBoundary {
            *pbAvail = false
        } else if onlyOneSliceInPic {
            *pbAvail = true
        } else {
            if bLCUTBoundary && bLCULBoundary {
                zRefSU = uint(numSUInLCU) - 1
                pRefMapLCU = pSliceIDMapLCU - (int(numLCUInPicWidth+1) * numSUInLCU)
            } else if bLCUTBoundary {
                zRefSU = G_auiRasterToZscan[rTRefSU-1]
                pRefMapLCU = pTRefMapLCU
            } else if bLCULBoundary {
                zRefSU = G_auiRasterToZscan[rLRefSU-int(numSUInLCUWidth)]
                pRefMapLCU = pLRefMapLCU
            } else { //inside LCU
                zRefSU = G_auiRasterToZscan[int(rTLSU)-int(numSUInLCUWidth)-1]
                pRefMapLCU = pSliceIDMapLCU
            }
            pRefID = &this.m_piSliceSUMap[pRefMapLCU+int(zRefSU)]
            if *pRefID == sliceID {
                *pbAvail = (true)
            } else if *pRefID > sliceID {
                *pbAvail = (LFCrossSliceBoundary[*pRefID])
            } else {
                *pbAvail = (LFCrossSliceBoundary[sliceID])
            }
        }

        //       SGU_TR
        pbAvail = &(pbAvailBorder[SGU_TR])
        if bPicTBoundary || bPicRBoundary {
            *pbAvail = false
        } else if onlyOneSliceInPic {
            *pbAvail = true
        } else {
            if bLCUTBoundary && bLCURBoundary {
                zRefSU = G_auiRasterToZscan[numSUInLCU-int(numSUInLCUWidth)]
                pRefMapLCU = pSliceIDMapLCU - (int(numLCUInPicWidth-1) * numSUInLCU)
            } else if bLCUTBoundary {
                zRefSU = G_auiRasterToZscan[rTRefSU+int(widthSU)]
                pRefMapLCU = pTRefMapLCU
            } else if bLCURBoundary {
                zRefSU = G_auiRasterToZscan[rRRefSU-int(numSUInLCUWidth)]
                pRefMapLCU = pRRefMapLCU
            } else { //inside LCU
                zRefSU = G_auiRasterToZscan[rTLSU-numSUInLCUWidth+widthSU]
                pRefMapLCU = pSliceIDMapLCU
            }
            pRefID = &this.m_piSliceSUMap[pRefMapLCU+int(zRefSU)]
            if *pRefID == sliceID {
                *pbAvail = (true)
            } else if *pRefID > sliceID {
                *pbAvail = (LFCrossSliceBoundary[*pRefID])
            } else {
                *pbAvail = (LFCrossSliceBoundary[sliceID])
            }
        }

        //       SGU_BL
        pbAvail = &(pbAvailBorder[SGU_BL])
        if bPicBBoundary || bPicLBoundary {
            *pbAvail = false
        } else if onlyOneSliceInPic {
            *pbAvail = true
        } else {
            if bLCUBBoundary && bLCULBoundary {
                zRefSU = G_auiRasterToZscan[numSUInLCUWidth-1]
                pRefMapLCU = pSliceIDMapLCU + (int(numLCUInPicWidth-1) * numSUInLCU)
            } else if bLCUBBoundary {
                zRefSU = G_auiRasterToZscan[rBRefSU-1]
                pRefMapLCU = pBRefMapLCU
            } else if bLCULBoundary {
                zRefSU = G_auiRasterToZscan[rLRefSU+int(heightSU*numSUInLCUWidth)]
                pRefMapLCU = pLRefMapLCU
            } else { //inside LCU
                zRefSU = G_auiRasterToZscan[int(rTLSU)+int(heightSU*numSUInLCUWidth)-1]
                pRefMapLCU = pSliceIDMapLCU
            }
            pRefID = &this.m_piSliceSUMap[pRefMapLCU+int(zRefSU)]
            if *pRefID == sliceID {
                *pbAvail = (true)
            } else if *pRefID > sliceID {
                *pbAvail = (LFCrossSliceBoundary[*pRefID])
            } else {
                *pbAvail = (LFCrossSliceBoundary[sliceID])
            }
        }

        //       SGU_BR
        pbAvail = &(pbAvailBorder[SGU_BR])
        if bPicBBoundary || bPicRBoundary {
            *pbAvail = false
        } else if onlyOneSliceInPic {
            *pbAvail = true
        } else {
            if bLCUBBoundary && bLCURBoundary {
                zRefSU = 0
                pRefMapLCU = pSliceIDMapLCU + (int(numLCUInPicWidth+1) * numSUInLCU)
            } else if bLCUBBoundary {
                zRefSU = G_auiRasterToZscan[rBRefSU+int(widthSU)]
                pRefMapLCU = pBRefMapLCU
            } else if bLCURBoundary {
                zRefSU = G_auiRasterToZscan[rRRefSU+int(heightSU*numSUInLCUWidth)]
                pRefMapLCU = pRRefMapLCU
            } else { //inside LCU
                zRefSU = G_auiRasterToZscan[int(rTLSU)+int(heightSU*numSUInLCUWidth)+int(widthSU)]
                pRefMapLCU = pSliceIDMapLCU
            }
            pRefID = &this.m_piSliceSUMap[pRefMapLCU+int(zRefSU)]
            if *pRefID == sliceID {
                *pbAvail = (true)
            } else if *pRefID > sliceID {
                *pbAvail = (LFCrossSliceBoundary[*pRefID])
            } else {
                *pbAvail = (LFCrossSliceBoundary[sliceID])
            }
        }

        if bIndependentTileBoundaryEnabled {
            //left LCU boundary
            if !bPicLBoundary && bLCULBoundary {
                if bLeftTileBoundary {
                    pbAvailBorder[SGU_L] = false
                    pbAvailBorder[SGU_TL] = false
                    pbAvailBorder[SGU_BL] = false
                }
            }
            //right LCU boundary
            if !bPicRBoundary && bLCURBoundary {
                if bRightTileBoundary {
                    pbAvailBorder[SGU_R] = false
                    pbAvailBorder[SGU_TR] = false
                    pbAvailBorder[SGU_BR] = false
                }
            }
            //top LCU boundary
            if !bPicTBoundary && bLCUTBoundary {
                if bTopTileBoundary {
                    pbAvailBorder[SGU_T] = false
                    pbAvailBorder[SGU_TL] = false
                    pbAvailBorder[SGU_TR] = false
                }
            }
            //down LCU boundary
            if !bPicBBoundary && bLCUBBoundary {
                if bDownTileBoundary {
                    pbAvailBorder[SGU_B] = false
                    pbAvailBorder[SGU_BL] = false
                    pbAvailBorder[SGU_BR] = false
                }
            }
        }
        rSGU.allBordersAvailable = true
        for b := 0; b < NUM_SGU_BORDER; b++ {
            if pbAvailBorder[b] == false {
                rSGU.allBordersAvailable = false
                break
            }
        }
    }
}

// -------------------------------------------------------------------------------------------------------------------
// member functions for accessing partition information
// -------------------------------------------------------------------------------------------------------------------
func (this *TComDataCU) GetPartIndexAndSizePos(uiPartIdx uint, ruiPartAddr *uint, riWidth, riHeight, rPosX, rPosY *int) {
    switch this.m_pePartSize[0] {
    case SIZE_2NxN:
        *riWidth = int(this.GetWidth1(0))
        *riHeight = int(this.GetHeight1(0)) >> 1
        if uiPartIdx == 0 {
            *ruiPartAddr = 0
            *rPosY = 0
        } else {
            *ruiPartAddr = this.m_uiNumPartition >> 1
            *rPosY = *riHeight
        }
        *rPosX = 0

    case SIZE_Nx2N:
        *riWidth = int(this.GetWidth1(0)) >> 1
        *riHeight = int(this.GetHeight1(0))
        if uiPartIdx == 0 {
            *ruiPartAddr = 0
            *rPosX = 0
        } else {
            *ruiPartAddr = this.m_uiNumPartition >> 2
            *rPosX = *riWidth
        }
        *rPosY = 0

    case SIZE_NxN:
        *riWidth = int(this.GetWidth1(0)) >> 1
        *riHeight = int(this.GetHeight1(0)) >> 1
        *ruiPartAddr = (this.m_uiNumPartition >> 2) * uiPartIdx
        if uiPartIdx%2 != 0 {
            *rPosX = *riWidth
        } else {
            *rPosX = 0
        }
        if uiPartIdx/2 != 0 {
            *rPosY = *riHeight
        } else {
            *rPosY = 0
        }

    case SIZE_2NxnU:
        *riWidth = int(this.GetWidth1(0))
        *rPosX = 0
        if uiPartIdx == 0 {
            *riHeight = int(this.GetHeight1(0)) >> 2
            *ruiPartAddr = 0
            *rPosY = 0
        } else {
            *riHeight = int(this.GetHeight1(0)>>2) + int(this.GetHeight1(0)>>1)
            *ruiPartAddr = this.m_uiNumPartition >> 3
            *rPosY = int(this.GetHeight1(0) >> 2)
        }

    case SIZE_2NxnD:
        *riWidth = int(this.GetWidth1(0))
        *rPosX = 0
        if uiPartIdx == 0 {
            *riHeight = int(this.GetHeight1(0)>>2) + int(this.GetHeight1(0)>>1)
            *ruiPartAddr = 0
            *rPosY = 0
        } else {
            *riHeight = int(this.GetHeight1(0) >> 2)
            *ruiPartAddr = (this.m_uiNumPartition >> 1) + (this.m_uiNumPartition >> 3)
            *rPosY = int(this.GetHeight1(0)>>2) + int(this.GetHeight1(0)>>1)
        }

    case SIZE_nLx2N:
        *riHeight = int(this.GetHeight1(0))
        *rPosY = 0
        if uiPartIdx == 0 {
            *riWidth = int(this.GetWidth1(0) >> 2)
            *ruiPartAddr = 0
            *rPosX = 0
        } else {
            *riWidth = int(this.GetWidth1(0)>>2) + int(this.GetWidth1(0)>>1)
            *ruiPartAddr = this.m_uiNumPartition >> 4
            *rPosX = int(this.GetWidth1(0) >> 2)
        }

    case SIZE_nRx2N:
        *riHeight = int(this.GetHeight1(0))
        *rPosY = 0
        if uiPartIdx == 0 {
            *riWidth = int(this.GetWidth1(0)>>2) + int(this.GetWidth1(0)>>1)
            *ruiPartAddr = 0
            *rPosX = 0
        } else {
            *riWidth = int(this.GetWidth1(0) >> 2)
            *ruiPartAddr = (this.m_uiNumPartition >> 2) + (this.m_uiNumPartition >> 4)
            *rPosX = int(this.GetWidth1(0)>>2) + int(this.GetWidth1(0)>>1)
        }

    default:
        //assert ( this.m_pePartSize[0] == SIZE_2Nx2N );
        *riWidth = int(this.GetWidth1(0))
        *riHeight = int(this.GetHeight1(0))
        *ruiPartAddr = 0
        *rPosX = 0
        *rPosY = 0
    }
}

func (this *TComDataCU) GetPartIndexAndSize(uiPartIdx uint, ruiPartAddr *uint, riWidth, riHeight *int) {
    switch this.m_pePartSize[0] {
    case SIZE_2NxN:
        *riWidth = int(this.GetWidth1(0))
        *riHeight = int(this.GetHeight1(0)) >> 1
        if uiPartIdx == 0 {
            *ruiPartAddr = 0
        } else {
            *ruiPartAddr = this.m_uiNumPartition >> 1
        }
    case SIZE_Nx2N:
        *riWidth = int(this.GetWidth1(0)) >> 1
        *riHeight = int(this.GetHeight1(0))
        if uiPartIdx == 0 {
            *ruiPartAddr = 0
        } else {
            *ruiPartAddr = this.m_uiNumPartition >> 2
        }
    case SIZE_NxN:
        *riWidth = int(this.GetWidth1(0)) >> 1
        *riHeight = int(this.GetHeight1(0)) >> 1
        *ruiPartAddr = (this.m_uiNumPartition >> 2) * uiPartIdx

    case SIZE_2NxnU:
        *riWidth = int(this.GetWidth1(0))
        if uiPartIdx == 0 {
            *riHeight = int(this.GetHeight1(0)) >> 2
            *ruiPartAddr = 0
        } else {
            *riHeight = int(this.GetHeight1(0)>>2) + int(this.GetHeight1(0)>>1)
            *ruiPartAddr = this.m_uiNumPartition >> 3
        }
    case SIZE_2NxnD:
        *riWidth = int(this.GetWidth1(0))
        if uiPartIdx == 0 {
            *riHeight = int(this.GetHeight1(0)>>2) + int(this.GetHeight1(0)>>1)
            *ruiPartAddr = 0
        } else {
            *riHeight = int(this.GetHeight1(0)) >> 2
            *ruiPartAddr = (this.m_uiNumPartition >> 1) + (this.m_uiNumPartition >> 3)
        }

    case SIZE_nLx2N:
        *riHeight = int(this.GetHeight1(0))
        if uiPartIdx == 0 {
            *riWidth = int(this.GetWidth1(0)) >> 2
            *ruiPartAddr = 0
        } else {
            *riWidth = int(this.GetWidth1(0)>>2) + int(this.GetWidth1(0)>>1)
            *ruiPartAddr = this.m_uiNumPartition >> 4
        }

    case SIZE_nRx2N:
        *riHeight = int(this.GetHeight1(0))
        if uiPartIdx == 0 {
            *riWidth = int(this.GetWidth1(0)>>2) + int(this.GetWidth1(0)>>1)
            *ruiPartAddr = 0
        } else {
            *riWidth = int(this.GetWidth1(0) >> 2)
            *ruiPartAddr = (this.m_uiNumPartition >> 2) + (this.m_uiNumPartition >> 4)
        }

    default:
        //assert ( m_pePartSize[0] == SIZE_2Nx2N );
        *riWidth = int(this.GetWidth1(0))
        *riHeight = int(this.GetHeight1(0))
        *ruiPartAddr = 0

    }
}
func (this *TComDataCU) GetNumPartInter() byte {
    iNumPart := byte(0)

    switch this.m_pePartSize[0] {
    case SIZE_2Nx2N:
        iNumPart = 1
    case SIZE_2NxN:
        iNumPart = 2
    case SIZE_Nx2N:
        iNumPart = 2
    case SIZE_NxN:
        iNumPart = 4
    case SIZE_2NxnU:
        iNumPart = 2
    case SIZE_2NxnD:
        iNumPart = 2
    case SIZE_nLx2N:
        iNumPart = 2
    case SIZE_nRx2N:
        iNumPart = 2
        //default:            assert (0);  ;
    }

    return iNumPart
}
func (this *TComDataCU) IsFirstAbsZorderIdxInDepth(uiAbsPartIdx, uiDepth uint) bool {
    uiPartNumb := this.m_pcPic.GetNumPartInCU() >> (uiDepth << 1)
    return (((this.m_uiAbsIdxInLCU + uiAbsPartIdx) % uiPartNumb) == 0)
}

// -------------------------------------------------------------------------------------------------------------------
// member functions for motion vector
// -------------------------------------------------------------------------------------------------------------------

func (this *TComDataCU) GetMvField(pcCU *TComDataCU, uiAbsPartIdx uint, eRefPicList RefPicList, rcMvField *TComMvField) {
    if pcCU == nil { // OUT OF BOUNDARY
        cZeroMv := NewTComMv(0, 0)
        rcMvField.SetMvField(*cZeroMv, NOT_VALID)
        return
    }

    pcCUMvField := pcCU.GetCUMvField(eRefPicList)
    rcMvField.SetMvField(pcCUMvField.GetMv(int(uiAbsPartIdx)), pcCUMvField.GetRefIdx(int(uiAbsPartIdx)))
}

func (this *TComDataCU) FillMvpCand(uiPartIdx, uiPartAddr uint, eRefPicList RefPicList, iRefIdx int, pInfo *AMVPInfo) {
    //var cMvPred TComMv;
    bAddedSmvp := false

    pInfo.IN = 0
    if iRefIdx < 0 {
        return
    }

    //-- Get Spatial MV
    var uiPartIdxLT, uiPartIdxRT, uiPartIdxLB uint
    uiNumPartInCUWidth := this.m_pcPic.GetNumPartInWidth()
    bAdded := false

    this.DeriveLeftRightTopIdx(uiPartIdx, &uiPartIdxLT, &uiPartIdxRT)
    this.DeriveLeftBottomIdx(uiPartIdx, &uiPartIdxLB)

    var tmpCU *TComDataCU
    var idx uint
    tmpCU = this.GetPUBelowLeft(&idx, uiPartIdxLB, true)
    bAddedSmvp = (tmpCU != nil) && (tmpCU.GetPredictionMode1(idx) != MODE_INTRA)

    if !bAddedSmvp {
        tmpCU = this.GetPULeft(&idx, uiPartIdxLB, true, true)
        bAddedSmvp = (tmpCU != nil) && (tmpCU.GetPredictionMode1(idx) != MODE_INTRA)
    }

    // Left predictor search
    bAdded = this.xAddMVPCand(pInfo, eRefPicList, iRefIdx, uiPartIdxLB, MD_BELOW_LEFT)
    if !bAdded {
        bAdded = this.xAddMVPCand(pInfo, eRefPicList, iRefIdx, uiPartIdxLB, MD_LEFT)
    }

    if !bAdded {
        bAdded = this.xAddMVPCandOrder(pInfo, eRefPicList, iRefIdx, uiPartIdxLB, MD_BELOW_LEFT)
        if !bAdded {
            bAdded = this.xAddMVPCandOrder(pInfo, eRefPicList, iRefIdx, uiPartIdxLB, MD_LEFT)
        }
    }
    // Above predictor search
    bAdded = this.xAddMVPCand(pInfo, eRefPicList, iRefIdx, uiPartIdxRT, MD_ABOVE_RIGHT)

    if !bAdded {
        bAdded = this.xAddMVPCand(pInfo, eRefPicList, iRefIdx, uiPartIdxRT, MD_ABOVE)
    }

    if !bAdded {
        bAdded = this.xAddMVPCand(pInfo, eRefPicList, iRefIdx, uiPartIdxLT, MD_ABOVE_LEFT)
    }
    bAdded = bAddedSmvp
    if pInfo.IN == 2 {
        bAdded = true
    }

    if !bAdded {
        bAdded = this.xAddMVPCandOrder(pInfo, eRefPicList, iRefIdx, uiPartIdxRT, MD_ABOVE_RIGHT)
        if !bAdded {
            bAdded = this.xAddMVPCandOrder(pInfo, eRefPicList, iRefIdx, uiPartIdxRT, MD_ABOVE)
        }

        if !bAdded {
            bAdded = this.xAddMVPCandOrder(pInfo, eRefPicList, iRefIdx, uiPartIdxLT, MD_ABOVE_LEFT)
        }
    }

    if pInfo.IN == 2 {
        if pInfo.MvCand[0] == pInfo.MvCand[1] {
            pInfo.IN = 1
        }
    }

    //fmt.Printf("iN=%d\n",pInfo.IN);
    if this.GetSlice().GetEnableTMVPFlag() {
        // Get Temporal Motion Predictor
        iRefIdx_Col := iRefIdx
        var cColMv TComMv
        var uiPartIdxRB, uiAbsPartIdx, uiAbsPartAddr uint
        uiLCUIdx := int(this.GetAddr())

        this.DeriveRightBottomIdx(uiPartIdx, &uiPartIdxRB)
        uiAbsPartAddr = this.m_uiAbsIdxInLCU + uiPartAddr

        //----  co-located RightBottom Temporal Predictor (H) ---//
        uiAbsPartIdx = G_auiZscanToRaster[uiPartIdxRB]
        if (this.m_pcPic.GetCU(this.m_uiCUAddr).GetCUPelX() + G_auiRasterToPelX[uiAbsPartIdx] + this.m_pcPic.GetMinCUWidth()) >= this.m_pcSlice.GetSPS().GetPicWidthInLumaSamples() { // image boundary check
            uiLCUIdx = -1
        } else if (this.m_pcPic.GetCU(this.m_uiCUAddr).GetCUPelY() + G_auiRasterToPelY[uiAbsPartIdx] + this.m_pcPic.GetMinCUHeight()) >= this.m_pcSlice.GetSPS().GetPicHeightInLumaSamples() {
            uiLCUIdx = -1
        } else {
            if (uiAbsPartIdx%uiNumPartInCUWidth < uiNumPartInCUWidth-1) && // is not at the last column of LCU
                (uiAbsPartIdx/uiNumPartInCUWidth < this.m_pcPic.GetNumPartInHeight()-1) { // is not at the last row    of LCU
                uiAbsPartAddr = G_auiRasterToZscan[uiAbsPartIdx+uiNumPartInCUWidth+1]
                uiLCUIdx = int(this.GetAddr())
            } else if uiAbsPartIdx%uiNumPartInCUWidth < uiNumPartInCUWidth-1 { // is not at the last column of LCU But is last row of LCU
                uiAbsPartAddr = G_auiRasterToZscan[(uiAbsPartIdx+uiNumPartInCUWidth+1)%this.m_pcPic.GetNumPartInCU()]
                uiLCUIdx = -1
            } else if uiAbsPartIdx/uiNumPartInCUWidth < this.m_pcPic.GetNumPartInHeight()-1 { // is not at the last row of LCU But is last column of LCU
                uiAbsPartAddr = G_auiRasterToZscan[uiAbsPartIdx+1]
                uiLCUIdx = int(this.GetAddr()) + 1
            } else { //is the right bottom corner of LCU
                uiAbsPartAddr = 0
                uiLCUIdx = -1
            }
        }
        if uiLCUIdx >= 0 && this.xGetColMVP(eRefPicList, uiLCUIdx, int(uiAbsPartAddr), &cColMv, &iRefIdx_Col) {
            //fmt.Printf("cColMv1=(%d,%d) ", cColMv.GetHor(), cColMv.GetVer());
            pInfo.MvCand[pInfo.IN] = cColMv
            pInfo.IN++
        } else {
            var uiPartIdxCenter uint
            uiCurLCUIdx := int(this.GetAddr())
            this.xDeriveCenterIdx(uiPartIdx, &uiPartIdxCenter)
            if this.xGetColMVP(eRefPicList, uiCurLCUIdx, int(uiPartIdxCenter), &cColMv, &iRefIdx_Col) {
                //fmt.Printf("cColMv2=(%d,%d) ", cColMv.GetHor(), cColMv.GetVer());
                pInfo.MvCand[pInfo.IN] = cColMv
                pInfo.IN++
            }
        }
        //----  co-located RightBottom Temporal Predictor  ---//
    }

    if pInfo.IN > AMVP_MAX_NUM_CANDS {
        pInfo.IN = AMVP_MAX_NUM_CANDS
    }
    for pInfo.IN < AMVP_MAX_NUM_CANDS {
        pInfo.MvCand[pInfo.IN].Set(0, 0)
        pInfo.IN++
    }

    /*for i:=0; i<pInfo.IN; i++ {
        fmt.Printf("(%d,%d) ",pInfo.MvCand[i].GetHor(),pInfo.MvCand[i].GetVer());
      }
      fmt.Printf("\n");*/

    return
}
func (this *TComDataCU) IsDiffMER(xN, yN, xP, yP int) bool {
    plevel := this.GetSlice().GetPPS().GetLog2ParallelMergeLevelMinus2() + 2
    if (xN >> plevel) != (xP >> plevel) {
        return true
    }

    if (yN >> plevel) != (yP >> plevel) {
        return true
    }

    return false
}
func (this *TComDataCU) GetPartPosition(partIdx uint, xP, yP, nPSW, nPSH *int) {
    col := int(this.m_uiCUPelX)
    row := int(this.m_uiCUPelY)

    switch this.m_pePartSize[0] {
    case SIZE_2NxN:
        *nPSW = int(this.GetWidth1(0))
        *nPSH = int(this.GetHeight1(0)) >> 1
        *xP = col
        if partIdx == 0 {
            *yP = row
        } else {
            *yP = row + *nPSH
        }
    case SIZE_Nx2N:
        *nPSW = int(this.GetWidth1(0)) >> 1
        *nPSH = int(this.GetHeight1(0))
        if partIdx == 0 {
            *xP = col
        } else {
            *xP = col + *nPSW
        }
        *yP = row
    case SIZE_NxN:
        *nPSW = int(this.GetWidth1(0)) >> 1
        *nPSH = int(this.GetHeight1(0)) >> 1
        *xP = col + int(partIdx&0x1)*(*nPSW)
        *yP = row + int(partIdx>>1)*(*nPSH)

    case SIZE_2NxnU:
        *nPSW = int(this.GetWidth1(0))
        *xP = col
        if partIdx == 0 {
            *nPSH = int(this.GetHeight1(0)) >> 2
            *yP = row
        } else {
            *nPSH = int(this.GetHeight1(0)>>2) + int(this.GetHeight1(0)>>1)
            *yP = row + int(this.GetHeight1(0)) - (*nPSH)
        }

    case SIZE_2NxnD:
        *nPSW = int(this.GetWidth1(0))
        *xP = col
        if partIdx == 0 {
            *nPSH = int(this.GetHeight1(0)>>2) + int(this.GetHeight1(0)>>1)
            *yP = row
        } else {
            *nPSH = int(this.GetHeight1(0)) >> 2
            *yP = row + int(this.GetHeight1(0)) - (*nPSH)
        }
    case SIZE_nLx2N:
        if partIdx == 0 {
            *nPSW = int(this.GetWidth1(0)) >> 2
            *xP = col
        } else {
            *nPSW = int(this.GetWidth1(0)>>2) + int(this.GetWidth1(0)>>1)
            *xP = col + int(this.GetWidth1(0)) - (*nPSW)
        }
        *nPSH = int(this.GetHeight1(0))
        *yP = row

    case SIZE_nRx2N:
        if partIdx == 0 {
            *nPSW = int(this.GetWidth1(0)>>2) + int(this.GetWidth1(0)>>1)
            *xP = col
        } else {
            *nPSW = int(this.GetWidth1(0) >> 2)
            *xP = col + int(this.GetWidth1(0)) - (*nPSW)
        }
        *nPSH = int(this.GetHeight1(0))
        *yP = row

    default:
        //assert ( m_pePartSize[0] == SIZE_2Nx2N );
        *nPSW = int(this.GetWidth1(0))
        *nPSH = int(this.GetHeight1(0))
        *xP = col
        *yP = row
    }
}
func (this *TComDataCU) SetMVPIdx(eRefPicList RefPicList, uiIdx uint, iMVPIdx int8) {
    this.m_apiMVPIdx[eRefPicList][uiIdx] = iMVPIdx
}
func (this *TComDataCU) GetMVPIdx2(eRefPicList RefPicList, uiIdx uint) int8 {
    return this.m_apiMVPIdx[eRefPicList][uiIdx]
}
func (this *TComDataCU) GetMVPIdx1(eRefPicList RefPicList) []int8 {
    return this.m_apiMVPIdx[eRefPicList]
}

func (this *TComDataCU) SetMVPNum(eRefPicList RefPicList, uiIdx uint, iMVPNum int8) {
    this.m_apiMVPNum[eRefPicList][uiIdx] = iMVPNum
}
func (this *TComDataCU) GetMVPNum2(eRefPicList RefPicList, uiIdx uint) int8 {
    return this.m_apiMVPNum[eRefPicList][uiIdx]
}
func (this *TComDataCU) GetMVPNum1(eRefPicList RefPicList) []int8 {
    return this.m_apiMVPNum[eRefPicList]
}

func (this *TComDataCU) SetMVPIdxSubParts(iMVPIdx int, eRefPicList RefPicList, uiAbsPartIdx, uiPartIdx, uiDepth uint) {
    this.SetSubPartInt8(int8(iMVPIdx), this.m_apiMVPIdx[eRefPicList], uiAbsPartIdx, uiDepth, uiPartIdx)
}
func (this *TComDataCU) SetMVPNumSubParts(iMVPNum int, eRefPicList RefPicList, uiAbsPartIdx, uiPartIdx, uiDepth uint) {
    this.SetSubPartInt8(int8(iMVPNum), this.m_apiMVPNum[eRefPicList], uiAbsPartIdx, uiDepth, uiPartIdx)
}

func (this *TComDataCU) ClipMv(rcMv *TComMv) {
    iMvShift := uint(2)
    iOffset := 8
    iHorMax := int16(this.m_pcSlice.GetSPS().GetPicWidthInLumaSamples()+uint(iOffset)-this.m_uiCUPelX-1) << iMvShift
    iHorMin := int16(-int(this.GetSlice().GetSPS().GetMaxCUWidth())-iOffset-int(this.m_uiCUPelX)+1) << iMvShift

    iVerMax := int16(this.m_pcSlice.GetSPS().GetPicHeightInLumaSamples()+uint(iOffset)-this.m_uiCUPelY-1) << iMvShift
    iVerMin := int16(-int(this.GetSlice().GetSPS().GetMaxCUHeight())-iOffset-int(this.m_uiCUPelY)+1) << iMvShift

    rcMv.SetHor(MIN(iHorMax, MAX(iHorMin, rcMv.GetHor()).(int16)).(int16))
    rcMv.SetVer(MIN(iVerMax, MAX(iVerMin, rcMv.GetVer()).(int16)).(int16))
}
func (this *TComDataCU) GetMvPredLeft() TComMv {
    return this.m_cMvFieldA.GetMv()
}
func (this *TComDataCU) GetMvPredAbove() TComMv {
    return this.m_cMvFieldB.GetMv()
}
func (this *TComDataCU) GetMvPredAboveRight() TComMv {
    return this.m_cMvFieldC.GetMv()
}

func (this *TComDataCU) CompressMV() {
    scaleFactor := 4 * AMVP_DECIMATION_FACTOR / this.m_unitSize
    if scaleFactor > 0 {
        this.m_acCUMvField[0].Compress(this.m_pePredMode, scaleFactor)
        this.m_acCUMvField[1].Compress(this.m_pePredMode, scaleFactor)
    }
}

// -------------------------------------------------------------------------------------------------------------------
// utility functions for neighbouring information
// -------------------------------------------------------------------------------------------------------------------

func (this *TComDataCU) GetCULeft() *TComDataCU {
    return this.m_pcCULeft
}
func (this *TComDataCU) GetCUAbove() *TComDataCU {
    return this.m_pcCUAbove
}
func (this *TComDataCU) GetCUAboveLeft() *TComDataCU {
    return this.m_pcCUAboveLeft
}
func (this *TComDataCU) GetCUAboveRight() *TComDataCU {
    return this.m_pcCUAboveRight
}

func (this *TComDataCU) GetCUColocated(eRefPicList RefPicList) *TComDataCU {
    return this.m_apcCUColocated[eRefPicList]
}

func (this *TComDataCU) GetPULeft(uiLPartUnitIdx *uint,
    uiCurrPartUnitIdx uint,
    bEnforceSliceRestriction bool,
    bEnforceTileRestriction bool) *TComDataCU {
    uiAbsPartIdx := G_auiZscanToRaster[uiCurrPartUnitIdx]
    uiAbsZorderCUIdx := G_auiZscanToRaster[this.m_uiAbsIdxInLCU]
    uiNumPartInCUWidth := this.m_pcPic.GetNumPartInWidth()
    //fmt.Printf("uiNumPartInCUWidth=%d\n %v", uiNumPartInCUWidth, IsZeroCol( int(uiAbsPartIdx), int(uiNumPartInCUWidth) ));

    if !IsZeroCol(int(uiAbsPartIdx), int(uiNumPartInCUWidth)) {
        *uiLPartUnitIdx = G_auiRasterToZscan[uiAbsPartIdx-1]
        if IsEqualCol(int(uiAbsPartIdx), int(uiAbsZorderCUIdx), int(uiNumPartInCUWidth)) {
            return this.m_pcPic.GetCU(this.GetAddr())
        } else {
            *uiLPartUnitIdx -= this.m_uiAbsIdxInLCU
            return this
        }
    }

    *uiLPartUnitIdx = G_auiRasterToZscan[uiAbsPartIdx+uiNumPartInCUWidth-1]

    if (bEnforceSliceRestriction && (this.m_pcCULeft == nil || this.m_pcCULeft.GetSlice() == nil || this.m_pcCULeft.GetSCUAddr()+(*uiLPartUnitIdx) < this.m_pcPic.GetCU(this.GetAddr()).GetSliceStartCU(uiCurrPartUnitIdx))) ||
        (bEnforceTileRestriction && (this.m_pcCULeft == nil || this.m_pcCULeft.GetSlice() == nil || (this.m_pcPic.GetPicSym().GetTileIdxMap(int(this.m_pcCULeft.GetAddr())) != this.m_pcPic.GetPicSym().GetTileIdxMap(int(this.GetAddr()))))) {
        return nil
    }
    return this.m_pcCULeft
}

func (this *TComDataCU) GetPUAbove(uiAPartUnitIdx *uint,
    uiCurrPartUnitIdx uint,
    bEnforceSliceRestriction bool,
    planarAtLCUBoundary bool,
    bEnforceTileRestriction bool) *TComDataCU {
    uiAbsPartIdx := int(G_auiZscanToRaster[uiCurrPartUnitIdx])
    uiAbsZorderCUIdx := int(G_auiZscanToRaster[this.m_uiAbsIdxInLCU])
    uiNumPartInCUWidth := int(this.m_pcPic.GetNumPartInWidth())

    if !IsZeroRow(uiAbsPartIdx, uiNumPartInCUWidth) {
        *uiAPartUnitIdx = G_auiRasterToZscan[uiAbsPartIdx-uiNumPartInCUWidth]
        if IsEqualRow(uiAbsPartIdx, uiAbsZorderCUIdx, uiNumPartInCUWidth) {
            return this.m_pcPic.GetCU(this.GetAddr())
        } else {
            *uiAPartUnitIdx -= this.m_uiAbsIdxInLCU
            return this
        }
    }

    if planarAtLCUBoundary {
        return nil
    }

    *uiAPartUnitIdx = G_auiRasterToZscan[uiAbsPartIdx+int(this.m_pcPic.GetNumPartInCU())-uiNumPartInCUWidth]
    /*#if !LINEBUF_CLEANUP
      if(MotionDataCompresssion)
      {
        uiAPartUnitIdx = G_motionRefer[uiAPartUnitIdx];
      }
    #endif*/

    if (bEnforceSliceRestriction && (this.m_pcCUAbove == nil || this.m_pcCUAbove.GetSlice() == nil || this.m_pcCUAbove.GetSCUAddr()+(*uiAPartUnitIdx) < this.m_pcPic.GetCU(this.GetAddr()).GetSliceStartCU(uiCurrPartUnitIdx))) ||
        (bEnforceTileRestriction && (this.m_pcCUAbove == nil || this.m_pcCUAbove.GetSlice() == nil || (this.m_pcPic.GetPicSym().GetTileIdxMap(int(this.m_pcCUAbove.GetAddr())) != this.m_pcPic.GetPicSym().GetTileIdxMap(int(this.GetAddr()))))) {
        return nil
    }
    return this.m_pcCUAbove
}
func (this *TComDataCU) GetPUAboveLeft(uiALPartUnitIdx *uint, uiCurrPartUnitIdx uint, bEnforceSliceRestriction bool) *TComDataCU {
    uiAbsPartIdx := int(G_auiZscanToRaster[uiCurrPartUnitIdx])
    uiAbsZorderCUIdx := int(G_auiZscanToRaster[this.m_uiAbsIdxInLCU])
    uiNumPartInCUWidth := int(this.m_pcPic.GetNumPartInWidth())

    if !IsZeroCol(uiAbsPartIdx, uiNumPartInCUWidth) {
        if !IsZeroRow(uiAbsPartIdx, uiNumPartInCUWidth) {
            *uiALPartUnitIdx = G_auiRasterToZscan[uiAbsPartIdx-uiNumPartInCUWidth-1]
            if IsEqualRowOrCol(uiAbsPartIdx, uiAbsZorderCUIdx, uiNumPartInCUWidth) {
                return this.m_pcPic.GetCU(this.GetAddr())
            } else {
                *uiALPartUnitIdx -= this.m_uiAbsIdxInLCU
                return this
            }
        }
        *uiALPartUnitIdx = G_auiRasterToZscan[uiAbsPartIdx+int(this.GetPic().GetNumPartInCU())-uiNumPartInCUWidth-1]
        if bEnforceSliceRestriction &&
            (this.m_pcCUAbove == nil ||
                this.m_pcCUAbove.GetSlice() == nil ||
                this.m_pcCUAbove.GetSCUAddr()+(*uiALPartUnitIdx) < this.m_pcPic.GetCU(this.GetAddr()).GetSliceStartCU(uiCurrPartUnitIdx) ||
                this.m_pcPic.GetPicSym().GetTileIdxMap(int(this.m_pcCUAbove.GetAddr())) != this.m_pcPic.GetPicSym().GetTileIdxMap(int(this.GetAddr()))) {
            return nil
        }
        return this.m_pcCUAbove
    }

    if !IsZeroRow(uiAbsPartIdx, uiNumPartInCUWidth) {
        *uiALPartUnitIdx = G_auiRasterToZscan[uiAbsPartIdx-1]
        if bEnforceSliceRestriction &&
            (this.m_pcCULeft == nil ||
                this.m_pcCULeft.GetSlice() == nil ||
                this.m_pcCULeft.GetSCUAddr()+(*uiALPartUnitIdx) < this.m_pcPic.GetCU(this.GetAddr()).GetSliceStartCU(uiCurrPartUnitIdx) ||
                this.m_pcPic.GetPicSym().GetTileIdxMap(int(this.m_pcCULeft.GetAddr())) != this.m_pcPic.GetPicSym().GetTileIdxMap(int(this.GetAddr()))) {
            return nil
        }
        return this.m_pcCULeft
    }

    *uiALPartUnitIdx = G_auiRasterToZscan[this.m_pcPic.GetNumPartInCU()-1]
    if bEnforceSliceRestriction &&
        (this.m_pcCUAboveLeft == nil ||
            this.m_pcCUAboveLeft.GetSlice() == nil ||
            this.m_pcCUAboveLeft.GetSCUAddr()+(*uiALPartUnitIdx) < this.m_pcPic.GetCU(this.GetAddr()).GetSliceStartCU(uiCurrPartUnitIdx) ||
            this.m_pcPic.GetPicSym().GetTileIdxMap(int(this.m_pcCUAboveLeft.GetAddr())) != this.m_pcPic.GetPicSym().GetTileIdxMap(int(this.GetAddr()))) {
        return nil
    }
    return this.m_pcCUAboveLeft
}
func (this *TComDataCU) GetPUAboveRight(uiARPartUnitIdx *uint, uiCurrPartUnitIdx uint, bEnforceSliceRestriction bool) *TComDataCU {
    uiAbsPartIdxRT := int(G_auiZscanToRaster[uiCurrPartUnitIdx])
    uiAbsZorderCUIdx := int(G_auiZscanToRaster[this.m_uiAbsIdxInLCU]) + int(this.m_puhWidth[0])/int(this.m_pcPic.GetMinCUWidth()) - 1
    uiNumPartInCUWidth := int(this.m_pcPic.GetNumPartInWidth())

    if (this.m_pcPic.GetCU(this.m_uiCUAddr).GetCUPelX() + G_auiRasterToPelX[uiAbsPartIdxRT] + this.m_pcPic.GetMinCUWidth()) >= this.m_pcSlice.GetSPS().GetPicWidthInLumaSamples() {
        *uiARPartUnitIdx = MAX_UINT
        return nil
    }

    if LessThanCol(uiAbsPartIdxRT, uiNumPartInCUWidth-1, uiNumPartInCUWidth) {
        if !IsZeroRow(uiAbsPartIdxRT, uiNumPartInCUWidth) {
            if uiCurrPartUnitIdx > G_auiRasterToZscan[uiAbsPartIdxRT-uiNumPartInCUWidth+1] {
                *uiARPartUnitIdx = G_auiRasterToZscan[uiAbsPartIdxRT-uiNumPartInCUWidth+1]
                if IsEqualRowOrCol(uiAbsPartIdxRT, uiAbsZorderCUIdx, uiNumPartInCUWidth) {
                    return this.m_pcPic.GetCU(this.GetAddr())
                } else {
                    *uiARPartUnitIdx -= this.m_uiAbsIdxInLCU
                    return this
                }
            }
            *uiARPartUnitIdx = MAX_UINT
            return nil
        }
        *uiARPartUnitIdx = G_auiRasterToZscan[uiAbsPartIdxRT+int(this.m_pcPic.GetNumPartInCU())-uiNumPartInCUWidth+1]
        if bEnforceSliceRestriction &&
            (this.m_pcCUAbove == nil ||
                this.m_pcCUAbove.GetSlice() == nil ||
                this.m_pcCUAbove.GetSCUAddr()+(*uiARPartUnitIdx) < this.m_pcPic.GetCU(this.GetAddr()).GetSliceStartCU(uiCurrPartUnitIdx) ||
                this.m_pcPic.GetPicSym().GetTileIdxMap(int(this.m_pcCUAbove.GetAddr())) != this.m_pcPic.GetPicSym().GetTileIdxMap(int(this.GetAddr()))) {
            return nil
        }
        return this.m_pcCUAbove
    }

    if !IsZeroRow(uiAbsPartIdxRT, uiNumPartInCUWidth) {
        *uiARPartUnitIdx = MAX_UINT
        return nil
    }

    *uiARPartUnitIdx = G_auiRasterToZscan[int(this.m_pcPic.GetNumPartInCU())-uiNumPartInCUWidth]
    if bEnforceSliceRestriction &&
        (this.m_pcCUAboveRight == nil ||
            this.m_pcCUAboveRight.GetSlice() == nil ||
            this.m_pcPic.GetPicSym().GetInverseCUOrderMap(int(this.m_pcCUAboveRight.GetAddr())) > this.m_pcPic.GetPicSym().GetInverseCUOrderMap(int(this.GetAddr())) ||
            this.m_pcCUAboveRight.GetSCUAddr()+(*uiARPartUnitIdx) < this.m_pcPic.GetCU(this.GetAddr()).GetSliceStartCU(uiCurrPartUnitIdx) ||
            this.m_pcPic.GetPicSym().GetTileIdxMap(int(this.m_pcCUAboveRight.GetAddr())) != this.m_pcPic.GetPicSym().GetTileIdxMap(int(this.GetAddr()))) {
        return nil
    }
    return this.m_pcCUAboveRight
}

//#endif
func (this *TComDataCU) GetPUBelowLeft(uiBLPartUnitIdx *uint, uiCurrPartUnitIdx uint, bEnforceSliceRestriction bool) *TComDataCU {
    uiAbsPartIdxLB := int(G_auiZscanToRaster[uiCurrPartUnitIdx])
    uiAbsZorderCUIdxLB := int(G_auiZscanToRaster[this.m_uiAbsIdxInLCU]) + (int(this.m_puhHeight[0])/int(this.m_pcPic.GetMinCUHeight())-1)*int(this.m_pcPic.GetNumPartInWidth())
    uiNumPartInCUWidth := int(this.m_pcPic.GetNumPartInWidth())

    if (this.m_pcPic.GetCU(this.m_uiCUAddr).GetCUPelY() + G_auiRasterToPelY[uiAbsPartIdxLB] + this.m_pcPic.GetMinCUHeight()) >= this.m_pcSlice.GetSPS().GetPicHeightInLumaSamples() {
        *uiBLPartUnitIdx = MAX_UINT
        return nil
    }

    if LessThanRow(uiAbsPartIdxLB, int(this.m_pcPic.GetNumPartInHeight())-1, uiNumPartInCUWidth) {
        if !IsZeroCol(uiAbsPartIdxLB, uiNumPartInCUWidth) {
            if uiCurrPartUnitIdx > G_auiRasterToZscan[uiAbsPartIdxLB+uiNumPartInCUWidth-1] {
                *uiBLPartUnitIdx = G_auiRasterToZscan[uiAbsPartIdxLB+uiNumPartInCUWidth-1]
                if IsEqualRowOrCol(uiAbsPartIdxLB, uiAbsZorderCUIdxLB, uiNumPartInCUWidth) {
                    return this.m_pcPic.GetCU(this.GetAddr())
                } else {
                    *uiBLPartUnitIdx -= this.m_uiAbsIdxInLCU
                    return this
                }
            }
            *uiBLPartUnitIdx = MAX_UINT
            return nil
        }
        *uiBLPartUnitIdx = G_auiRasterToZscan[uiAbsPartIdxLB+uiNumPartInCUWidth*2-1]
        if bEnforceSliceRestriction &&
            (this.m_pcCULeft == nil ||
                this.m_pcCULeft.GetSlice() == nil ||
                this.m_pcCULeft.GetSCUAddr()+(*uiBLPartUnitIdx) < this.m_pcPic.GetCU(this.GetAddr()).GetSliceStartCU(uiCurrPartUnitIdx) ||
                this.m_pcPic.GetPicSym().GetTileIdxMap(int(this.m_pcCULeft.GetAddr())) != this.m_pcPic.GetPicSym().GetTileIdxMap(int(this.GetAddr()))) {
            return nil
        }
        return this.m_pcCULeft
    }

    *uiBLPartUnitIdx = MAX_UINT
    return nil
}

func (this *TComDataCU) GetQpMinCuLeft(uiLPartUnitIdx *uint, uiCurrAbsIdxInLCU uint) *TComDataCU {
    numPartInCUWidth := int(this.m_pcPic.GetNumPartInWidth())
    absZorderQpMinCUIdx := (uiCurrAbsIdxInLCU >> ((this.GetSlice().GetSPS().GetMaxCUDepth() - this.GetSlice().GetPPS().GetMaxCuDQPDepth()) << 1)) << ((this.GetSlice().GetSPS().GetMaxCUDepth() - this.GetSlice().GetPPS().GetMaxCuDQPDepth()) << 1)
    absRorderQpMinCUIdx := int(G_auiZscanToRaster[absZorderQpMinCUIdx])

    // check for left LCU boundary
    if IsZeroCol(absRorderQpMinCUIdx, numPartInCUWidth) {
        return nil
    }

    // this.Get index of left-CU relative to top-left corner of current quantization group
    *uiLPartUnitIdx = G_auiRasterToZscan[absRorderQpMinCUIdx-1]

    // return pointer to current LCU
    return this.m_pcPic.GetCU(this.GetAddr())
}
func (this *TComDataCU) GetQpMinCuAbove(aPartUnitIdx *uint, currAbsIdxInLCU uint) *TComDataCU {
    numPartInCUWidth := int(this.m_pcPic.GetNumPartInWidth())
    absZorderQpMinCUIdx := (currAbsIdxInLCU >> ((this.GetSlice().GetSPS().GetMaxCUDepth() - this.GetSlice().GetPPS().GetMaxCuDQPDepth()) << 1)) << ((this.GetSlice().GetSPS().GetMaxCUDepth() - this.GetSlice().GetPPS().GetMaxCuDQPDepth()) << 1)
    absRorderQpMinCUIdx := int(G_auiZscanToRaster[absZorderQpMinCUIdx])

    // check for top LCU boundary
    if IsZeroRow(absRorderQpMinCUIdx, numPartInCUWidth) {
        return nil
    }

    // this.Get index of top-CU relative to top-left corner of current quantization group
    *aPartUnitIdx = G_auiRasterToZscan[absRorderQpMinCUIdx-numPartInCUWidth]

    // return pointer to current LCU
    return this.m_pcPic.GetCU(this.GetAddr())
}
func (this *TComDataCU) GetRefQP(uiCurrAbsIdxInLCU uint) int8 {
    lPartIdx := uint(0)
    aPartIdx := uint(0)
    cULeft := this.GetQpMinCuLeft(&lPartIdx, this.m_uiAbsIdxInLCU+uiCurrAbsIdxInLCU)
    cUAbove := this.GetQpMinCuAbove(&aPartIdx, this.m_uiAbsIdxInLCU+uiCurrAbsIdxInLCU)

    if cULeft != nil && cUAbove != nil {
        return (cULeft.GetQP1(lPartIdx) + cUAbove.GetQP1(aPartIdx) + 1) >> 1
    } else if cUAbove != nil {
        return (this.GetLastCodedQP(uiCurrAbsIdxInLCU) + cUAbove.GetQP1(aPartIdx) + 1) >> 1
    } else if cULeft != nil {
        return (cULeft.GetQP1(lPartIdx) + this.GetLastCodedQP(uiCurrAbsIdxInLCU) + 1) >> 1
    }

    return (this.GetLastCodedQP(uiCurrAbsIdxInLCU) + this.GetLastCodedQP(uiCurrAbsIdxInLCU) + 1) >> 1
}

func (this *TComDataCU) GetPUAboveRightAdi(uiARPartUnitIdx *uint, uiCurrPartUnitIdx, uiPartUnitOffset uint, bEnforceSliceRestriction bool) *TComDataCU {
    uiAbsPartIdxRT := int(G_auiZscanToRaster[uiCurrPartUnitIdx])
    uiAbsZorderCUIdx := int(G_auiZscanToRaster[this.m_uiAbsIdxInLCU]) + (int(this.m_puhWidth[0]) / int(this.m_pcPic.GetMinCUWidth())) - 1
    uiNumPartInCUWidth := int(this.m_pcPic.GetNumPartInWidth())

    if (this.m_pcPic.GetCU(this.m_uiCUAddr).GetCUPelX() + G_auiRasterToPelX[uiAbsPartIdxRT] + (this.m_pcPic.GetPicSym().GetMinCUHeight() * uiPartUnitOffset)) >= this.m_pcSlice.GetSPS().GetPicWidthInLumaSamples() {
        *uiARPartUnitIdx = MAX_UINT
        return nil
    }

    if LessThanCol(uiAbsPartIdxRT, uiNumPartInCUWidth-int(uiPartUnitOffset), uiNumPartInCUWidth) {
        if !IsZeroRow(uiAbsPartIdxRT, uiNumPartInCUWidth) {
            if uiCurrPartUnitIdx > G_auiRasterToZscan[uiAbsPartIdxRT-uiNumPartInCUWidth+int(uiPartUnitOffset)] {
                *uiARPartUnitIdx = G_auiRasterToZscan[uiAbsPartIdxRT-uiNumPartInCUWidth+int(uiPartUnitOffset)]
                if IsEqualRowOrCol(uiAbsPartIdxRT, uiAbsZorderCUIdx, uiNumPartInCUWidth) {
                    return this.m_pcPic.GetCU(this.GetAddr())
                } else {
                    *uiARPartUnitIdx -= this.m_uiAbsIdxInLCU
                    return this
                }
            }
            *uiARPartUnitIdx = MAX_UINT
            return nil
        }
        *uiARPartUnitIdx = G_auiRasterToZscan[uiAbsPartIdxRT+int(this.m_pcPic.GetNumPartInCU())-uiNumPartInCUWidth+int(uiPartUnitOffset)]
        if bEnforceSliceRestriction &&
            (this.m_pcCUAbove == nil ||
                this.m_pcCUAbove.GetSlice() == nil ||
                this.m_pcCUAbove.GetSCUAddr()+(*uiARPartUnitIdx) < this.m_pcPic.GetCU(this.GetAddr()).GetSliceStartCU(uiCurrPartUnitIdx) ||
                this.m_pcPic.GetPicSym().GetTileIdxMap(int(this.m_pcCUAbove.GetAddr())) != this.m_pcPic.GetPicSym().GetTileIdxMap(int(this.GetAddr()))) {
            return nil
        }
        return this.m_pcCUAbove
    }

    if !IsZeroRow(uiAbsPartIdxRT, uiNumPartInCUWidth) {
        *uiARPartUnitIdx = MAX_UINT
        return nil
    }

    *uiARPartUnitIdx = G_auiRasterToZscan[int(this.m_pcPic.GetNumPartInCU())-uiNumPartInCUWidth+int(uiPartUnitOffset)-1]
    if bEnforceSliceRestriction &&
        (this.m_pcCUAboveRight == nil ||
            this.m_pcCUAboveRight.GetSlice() == nil ||
            this.m_pcPic.GetPicSym().GetInverseCUOrderMap(int(this.m_pcCUAboveRight.GetAddr())) > this.m_pcPic.GetPicSym().GetInverseCUOrderMap(int(this.GetAddr())) ||
            this.m_pcCUAboveRight.GetSCUAddr()+(*uiARPartUnitIdx) < this.m_pcPic.GetCU(this.GetAddr()).GetSliceStartCU(uiCurrPartUnitIdx) ||
            this.m_pcPic.GetPicSym().GetTileIdxMap(int(this.m_pcCUAboveRight.GetAddr())) != this.m_pcPic.GetPicSym().GetTileIdxMap(int(this.GetAddr()))) {
        return nil
    }
    return this.m_pcCUAboveRight
}
func (this *TComDataCU) GetPUBelowLeftAdi(uiBLPartUnitIdx *uint, uiCurrPartUnitIdx, uiPartUnitOffset uint, bEnforceSliceRestriction bool) *TComDataCU {
    uiAbsPartIdxLB := int(G_auiZscanToRaster[uiCurrPartUnitIdx])
    uiAbsZorderCUIdxLB := int(G_auiZscanToRaster[this.m_uiAbsIdxInLCU]) + ((int(this.m_puhHeight[0])/int(this.m_pcPic.GetMinCUHeight()))-1)*int(this.m_pcPic.GetNumPartInWidth())
    uiNumPartInCUWidth := int(this.m_pcPic.GetNumPartInWidth())

    if (this.m_pcPic.GetCU(this.m_uiCUAddr).GetCUPelY() + G_auiRasterToPelY[uiAbsPartIdxLB] + (this.m_pcPic.GetPicSym().GetMinCUHeight() * uiPartUnitOffset)) >= this.m_pcSlice.GetSPS().GetPicHeightInLumaSamples() {
        *uiBLPartUnitIdx = MAX_UINT
        return nil
    }

    if LessThanRow(uiAbsPartIdxLB, int(this.m_pcPic.GetNumPartInHeight()-uiPartUnitOffset), uiNumPartInCUWidth) {
        if !IsZeroCol(uiAbsPartIdxLB, uiNumPartInCUWidth) {
            if uiCurrPartUnitIdx > G_auiRasterToZscan[uiAbsPartIdxLB+int(uiPartUnitOffset)*uiNumPartInCUWidth-1] {
                *uiBLPartUnitIdx = G_auiRasterToZscan[uiAbsPartIdxLB+int(uiPartUnitOffset)*uiNumPartInCUWidth-1]
                if IsEqualRowOrCol(uiAbsPartIdxLB, uiAbsZorderCUIdxLB, uiNumPartInCUWidth) {
                    return this.m_pcPic.GetCU(this.GetAddr())
                } else {
                    *uiBLPartUnitIdx -= this.m_uiAbsIdxInLCU
                    return this
                }
            }
            *uiBLPartUnitIdx = MAX_UINT
            return nil
        }
        *uiBLPartUnitIdx = G_auiRasterToZscan[uiAbsPartIdxLB+(1+int(uiPartUnitOffset))*uiNumPartInCUWidth-1]
        if bEnforceSliceRestriction &&
            (this.m_pcCULeft == nil ||
                this.m_pcCULeft.GetSlice() == nil ||
                this.m_pcCULeft.GetSCUAddr()+(*uiBLPartUnitIdx) < this.m_pcPic.GetCU(this.GetAddr()).GetSliceStartCU(uiCurrPartUnitIdx) ||
                this.m_pcPic.GetPicSym().GetTileIdxMap(int(this.m_pcCULeft.GetAddr())) != this.m_pcPic.GetPicSym().GetTileIdxMap(int(this.GetAddr()))) {
            return nil
        }
        return this.m_pcCULeft
    }

    *uiBLPartUnitIdx = MAX_UINT
    return nil
}

func (this *TComDataCU) DeriveLeftRightTopIdx(uiPartIdx uint, ruiPartIdxLT, ruiPartIdxRT *uint) {
    *ruiPartIdxLT = this.m_uiAbsIdxInLCU
    *ruiPartIdxRT = G_auiRasterToZscan[G_auiZscanToRaster[*ruiPartIdxLT]+uint(this.m_puhWidth[0])/this.m_pcPic.GetMinCUWidth()-1]

    switch this.m_pePartSize[0] {
    case SIZE_2Nx2N:
        fallthrough
    case SIZE_2NxN:
        if uiPartIdx == 0 {
            *ruiPartIdxLT += 0
            *ruiPartIdxRT += 0
        } else {
            *ruiPartIdxLT += this.m_uiNumPartition >> 1
            *ruiPartIdxRT += this.m_uiNumPartition >> 1
        }
    case SIZE_Nx2N:
        if uiPartIdx == 0 {
            *ruiPartIdxLT += 0
        } else {
            *ruiPartIdxLT += this.m_uiNumPartition >> 2
        }
        if uiPartIdx == 1 {
            *ruiPartIdxRT -= 0
        } else {
            *ruiPartIdxRT -= this.m_uiNumPartition >> 2
        }
    case SIZE_NxN:
        *ruiPartIdxLT += (this.m_uiNumPartition >> 2) * uiPartIdx
        *ruiPartIdxRT += (this.m_uiNumPartition >> 2) * (uiPartIdx - 1)

    case SIZE_2NxnU:
        if uiPartIdx == 0 {
            *ruiPartIdxLT += 0
            *ruiPartIdxRT += 0
        } else {
            *ruiPartIdxLT += this.m_uiNumPartition >> 3
            *ruiPartIdxRT += this.m_uiNumPartition >> 3
        }

    case SIZE_2NxnD:
        if uiPartIdx == 0 {
            *ruiPartIdxLT += 0
            *ruiPartIdxRT += 0
        } else {
            *ruiPartIdxLT += (this.m_uiNumPartition >> 1) + (this.m_uiNumPartition >> 3)
            *ruiPartIdxRT += (this.m_uiNumPartition >> 1) + (this.m_uiNumPartition >> 3)
        }
    case SIZE_nLx2N:
        if uiPartIdx == 0 {
            *ruiPartIdxLT += 0
        } else {
            *ruiPartIdxLT += this.m_uiNumPartition >> 4
        }
        if uiPartIdx == 1 {
            *ruiPartIdxRT -= 0
        } else {
            *ruiPartIdxRT -= (this.m_uiNumPartition >> 2) + (this.m_uiNumPartition >> 4)
        }
    case SIZE_nRx2N:
        if uiPartIdx == 0 {
            *ruiPartIdxLT += 0
        } else {
            *ruiPartIdxLT += (this.m_uiNumPartition >> 2) + (this.m_uiNumPartition >> 4)
        }
        if uiPartIdx == 1 {
            *ruiPartIdxRT -= 0
        } else {
            *ruiPartIdxRT -= this.m_uiNumPartition >> 4
        }

    default:

    }
}

func (this *TComDataCU) DeriveLeftRightTopIdxGeneral(uiAbsPartIdx, uiPartIdx uint, ruiPartIdxLT, ruiPartIdxRT *uint) {
    *ruiPartIdxLT = this.m_uiAbsIdxInLCU + uiAbsPartIdx
    uiPUWidth := byte(0)

    switch this.m_pePartSize[uiAbsPartIdx] {
    case SIZE_2Nx2N:
        uiPUWidth = this.m_puhWidth[uiAbsPartIdx]
    case SIZE_2NxN:
        uiPUWidth = this.m_puhWidth[uiAbsPartIdx]
    case SIZE_Nx2N:
        uiPUWidth = this.m_puhWidth[uiAbsPartIdx] >> 1
    case SIZE_NxN:
        uiPUWidth = this.m_puhWidth[uiAbsPartIdx] >> 1
    case SIZE_2NxnU:
        uiPUWidth = this.m_puhWidth[uiAbsPartIdx]
    case SIZE_2NxnD:
        uiPUWidth = this.m_puhWidth[uiAbsPartIdx]
    case SIZE_nLx2N:
        if uiPartIdx == 0 {
            uiPUWidth = this.m_puhWidth[uiAbsPartIdx] >> 2
        } else if uiPartIdx == 1 {
            uiPUWidth = (this.m_puhWidth[uiAbsPartIdx] >> 1) + (this.m_puhWidth[uiAbsPartIdx] >> 2)
        } else {
            //assert(0);
        }

    case SIZE_nRx2N:
        if uiPartIdx == 0 {
            uiPUWidth = (this.m_puhWidth[uiAbsPartIdx] >> 1) + (this.m_puhWidth[uiAbsPartIdx] >> 2)
        } else if uiPartIdx == 1 {
            uiPUWidth = this.m_puhWidth[uiAbsPartIdx] >> 2
        } else {
            //assert(0);
        }

    default:
        //assert (0);
        //break;
    }

    *ruiPartIdxRT = G_auiRasterToZscan[G_auiZscanToRaster[*ruiPartIdxLT]+uint(uiPUWidth)/this.m_pcPic.GetMinCUWidth()-1]
}
func (this *TComDataCU) DeriveLeftBottomIdx(uiPartIdx uint, ruiPartIdxLB *uint) {
    *ruiPartIdxLB = G_auiRasterToZscan[G_auiZscanToRaster[this.m_uiAbsIdxInLCU]+(((uint(this.m_puhHeight[0])/this.m_pcPic.GetMinCUHeight())>>1)-1)*this.m_pcPic.GetNumPartInWidth()]

    switch this.m_pePartSize[0] {
    case SIZE_2Nx2N:
        *ruiPartIdxLB += this.m_uiNumPartition >> 1

    case SIZE_2NxN:
        if uiPartIdx == 0 {
            *ruiPartIdxLB += 0
        } else {
            *ruiPartIdxLB += this.m_uiNumPartition >> 1
        }

    case SIZE_Nx2N:
        if uiPartIdx == 0 {
            *ruiPartIdxLB += this.m_uiNumPartition >> 1
        } else {
            *ruiPartIdxLB += (this.m_uiNumPartition >> 2) * 3
        }
    case SIZE_NxN:
        *ruiPartIdxLB += (this.m_uiNumPartition >> 2) * uiPartIdx

    case SIZE_2NxnU:
        if uiPartIdx == 0 {
            *ruiPartIdxLB -= this.m_uiNumPartition >> 3 //-(int(this.m_uiNumPartition) >> 3);
        } else {
            *ruiPartIdxLB += this.m_uiNumPartition >> 1
        }

    case SIZE_2NxnD:
        if uiPartIdx == 0 {
            *ruiPartIdxLB += (this.m_uiNumPartition >> 2) + (this.m_uiNumPartition >> 3)
        } else {
            *ruiPartIdxLB += this.m_uiNumPartition >> 1
        }

    case SIZE_nLx2N:
        if uiPartIdx == 0 {
            *ruiPartIdxLB += this.m_uiNumPartition >> 1
        } else {
            *ruiPartIdxLB += (this.m_uiNumPartition >> 1) + (this.m_uiNumPartition >> 4)
        }

    case SIZE_nRx2N:
        if uiPartIdx == 0 {
            *ruiPartIdxLB += this.m_uiNumPartition >> 1
        } else {
            *ruiPartIdxLB += (this.m_uiNumPartition >> 1) + (this.m_uiNumPartition >> 2) + (this.m_uiNumPartition >> 4)
        }

    default:

    }
}

func (this *TComDataCU) DeriveLeftBottomIdxGeneral(uiAbsPartIdx, uiPartIdx uint, ruiPartIdxLB *uint) {
    uiPUHeight := byte(0)
    switch this.m_pePartSize[uiAbsPartIdx] {
    case SIZE_2Nx2N:
        uiPUHeight = this.m_puhHeight[uiAbsPartIdx]
    case SIZE_2NxN:
        uiPUHeight = this.m_puhHeight[uiAbsPartIdx] >> 1
    case SIZE_Nx2N:
        uiPUHeight = this.m_puhHeight[uiAbsPartIdx]
    case SIZE_NxN:
        uiPUHeight = this.m_puhHeight[uiAbsPartIdx] >> 1
    case SIZE_2NxnU:
        if uiPartIdx == 0 {
            uiPUHeight = this.m_puhHeight[uiAbsPartIdx] >> 2
        } else if uiPartIdx == 1 {
            uiPUHeight = (this.m_puhHeight[uiAbsPartIdx] >> 1) + (this.m_puhHeight[uiAbsPartIdx] >> 2)
        } else {
            // assert(0);
        }

    case SIZE_2NxnD:
        if uiPartIdx == 0 {
            uiPUHeight = (this.m_puhHeight[uiAbsPartIdx] >> 1) + (this.m_puhHeight[uiAbsPartIdx] >> 2)
        } else if uiPartIdx == 1 {
            uiPUHeight = this.m_puhHeight[uiAbsPartIdx] >> 2
        } else {
            //assert(0);
        }

    case SIZE_nLx2N:
        uiPUHeight = this.m_puhHeight[uiAbsPartIdx]
    case SIZE_nRx2N:
        uiPUHeight = this.m_puhHeight[uiAbsPartIdx]
    default:
        // assert (0);
        // break;
    }

    *ruiPartIdxLB = G_auiRasterToZscan[G_auiZscanToRaster[this.m_uiAbsIdxInLCU+uiAbsPartIdx]+((uint(uiPUHeight)/this.m_pcPic.GetMinCUHeight())-1)*this.m_pcPic.GetNumPartInWidth()]
}

func (this *TComDataCU) DeriveLeftRightTopIdxAdi(ruiPartIdxLT, ruiPartIdxRT *uint, uiPartOffset, uiPartDepth uint) {
    uiNumPartInWidth := (uint(this.m_puhWidth[0]) / this.m_pcPic.GetMinCUWidth()) >> uiPartDepth
    *ruiPartIdxLT = this.m_uiAbsIdxInLCU + uiPartOffset
    *ruiPartIdxRT = G_auiRasterToZscan[G_auiZscanToRaster[*ruiPartIdxLT]+uiNumPartInWidth-1]
}
func (this *TComDataCU) DeriveLeftBottomIdxAdi(ruiPartIdxLB *uint, uiPartOffset, uiPartDepth uint) {
    var uiAbsIdx, uiMinCuWidth, uiWidthInMinCus uint

    uiMinCuWidth = this.GetPic().GetMinCUWidth()
    uiWidthInMinCus = (uint(this.GetWidth1(0)) / uiMinCuWidth) >> uiPartDepth
    uiAbsIdx = this.GetZorderIdxInCU() + uiPartOffset + (this.m_uiNumPartition >> (uiPartDepth << 1)) - 1
    uiAbsIdx = G_auiZscanToRaster[uiAbsIdx] - (uiWidthInMinCus - 1)
    *ruiPartIdxLB = G_auiRasterToZscan[uiAbsIdx]
}

func (this *TComDataCU) HasEqualMotion(uiAbsPartIdx uint, pcCandCU *TComDataCU, uiCandAbsPartIdx uint) bool {
    if this.GetInterDir1(uiAbsPartIdx) != pcCandCU.GetInterDir1(uiCandAbsPartIdx) {
        return false
    }

    for uiRefListIdx := uint(0); uiRefListIdx < 2; uiRefListIdx++ {
        if (this.GetInterDir1(uiAbsPartIdx) & (1 << uiRefListIdx)) != 0 {
            if this.GetCUMvField(RefPicList(uiRefListIdx)).GetMv(int(uiAbsPartIdx)) != pcCandCU.GetCUMvField(RefPicList(uiRefListIdx)).GetMv(int(uiCandAbsPartIdx)) ||
                this.GetCUMvField(RefPicList(uiRefListIdx)).GetRefIdx(int(uiAbsPartIdx)) != pcCandCU.GetCUMvField(RefPicList(uiRefListIdx)).GetRefIdx(int(uiCandAbsPartIdx)) {
                return false
            }
        }
    }

    return true
}
func (this *TComDataCU) GetInterMergeCandidates(uiAbsPartIdx, uiPUIdx uint, pcMvFieldNeighbours []TComMvField, puhInterDirNeighbours []byte, numValidMergeCand *int, mrgCandIdx int) {
    uiAbsPartAddr := this.m_uiAbsIdxInLCU + uiAbsPartIdx
    var abCandIsInter [MRG_MAX_NUM_CANDS]bool
    for ui := uint(0); ui < this.GetSlice().GetMaxNumMergeCand(); ui++ {
        abCandIsInter[ui] = false
        pcMvFieldNeighbours[(ui << 1)].SetRefIdx(NOT_VALID)
        pcMvFieldNeighbours[(ui<<1)+1].SetRefIdx(NOT_VALID)
    }
    *numValidMergeCand = int(this.GetSlice().GetMaxNumMergeCand())
    // compute the location of the current PU
    var xP, yP, nPSW, nPSH int
    this.GetPartPosition(uiPUIdx, &xP, &yP, &nPSW, &nPSH)

    iCount := int(0)

    var uiPartIdxLT, uiPartIdxRT, uiPartIdxLB uint
    cCurPS := this.GetPartitionSize1(uiAbsPartIdx)
    this.DeriveLeftRightTopIdxGeneral(uiAbsPartIdx, uiPUIdx, &uiPartIdxLT, &uiPartIdxRT)
    this.DeriveLeftBottomIdxGeneral(uiAbsPartIdx, uiPUIdx, &uiPartIdxLB)

    //left
    uiLeftPartIdx := uint(0)
    var pcCULeft *TComDataCU
    pcCULeft = this.GetPULeft(&uiLeftPartIdx, uiPartIdxLB, true, true)
    isAvailableA1 := pcCULeft != nil &&
        pcCULeft.IsDiffMER(xP-1, yP+nPSH-1, xP, yP) &&
        !(uiPUIdx == 1 && (cCurPS == SIZE_Nx2N || cCurPS == SIZE_nLx2N || cCurPS == SIZE_nRx2N)) &&
        !pcCULeft.IsIntra(uiLeftPartIdx)
    if isAvailableA1 {
        abCandIsInter[iCount] = true
        // this.Get Inter Dir
        puhInterDirNeighbours[iCount] = pcCULeft.GetInterDir1(uiLeftPartIdx)
        // this.Get Mv from Left
        pcCULeft.GetMvField(pcCULeft, uiLeftPartIdx, REF_PIC_LIST_0, &pcMvFieldNeighbours[iCount<<1])
        if this.GetSlice().IsInterB() {
            pcCULeft.GetMvField(pcCULeft, uiLeftPartIdx, REF_PIC_LIST_1, &pcMvFieldNeighbours[(iCount<<1)+1])
        }
        if mrgCandIdx == iCount {
            return
        }
        iCount++
    }

    // early termination
    if iCount == int(this.GetSlice().GetMaxNumMergeCand()) {
        return
    }
    // above
    uiAbovePartIdx := uint(0)
    var pcCUAbove *TComDataCU
    pcCUAbove = this.GetPUAbove(&uiAbovePartIdx, uiPartIdxRT, true, false, true)
    isAvailableB1 := pcCUAbove != nil &&
        pcCUAbove.IsDiffMER(xP+nPSW-1, yP-1, xP, yP) &&
        !(uiPUIdx == 1 && (cCurPS == SIZE_2NxN || cCurPS == SIZE_2NxnU || cCurPS == SIZE_2NxnD)) &&
        !pcCUAbove.IsIntra(uiAbovePartIdx)
    if isAvailableB1 && (!isAvailableA1 || !pcCULeft.HasEqualMotion(uiLeftPartIdx, pcCUAbove, uiAbovePartIdx)) {
        abCandIsInter[iCount] = true
        // this.Get Inter Dir
        puhInterDirNeighbours[iCount] = pcCUAbove.GetInterDir1(uiAbovePartIdx)
        // this.Get Mv from Left
        pcCUAbove.GetMvField(pcCUAbove, uiAbovePartIdx, REF_PIC_LIST_0, &pcMvFieldNeighbours[iCount<<1])
        if this.GetSlice().IsInterB() {
            pcCUAbove.GetMvField(pcCUAbove, uiAbovePartIdx, REF_PIC_LIST_1, &pcMvFieldNeighbours[(iCount<<1)+1])
        }
        if mrgCandIdx == iCount {
            return
        }
        iCount++
    }
    // early termination
    if iCount == int(this.GetSlice().GetMaxNumMergeCand()) {
        return
    }

    // above right
    uiAboveRightPartIdx := uint(0)
    var pcCUAboveRight *TComDataCU

    pcCUAboveRight = this.GetPUAboveRight(&uiAboveRightPartIdx, uiPartIdxRT, true)
    isAvailableB0 := pcCUAboveRight != nil &&
        pcCUAboveRight.IsDiffMER(xP+nPSW, yP-1, xP, yP) &&
        !pcCUAboveRight.IsIntra(uiAboveRightPartIdx)
    if isAvailableB0 && (!isAvailableB1 || !pcCUAbove.HasEqualMotion(uiAbovePartIdx, pcCUAboveRight, uiAboveRightPartIdx)) {
        abCandIsInter[iCount] = true
        // this.Get Inter Dir
        puhInterDirNeighbours[iCount] = pcCUAboveRight.GetInterDir1(uiAboveRightPartIdx)
        // this.Get Mv from Left
        pcCUAboveRight.GetMvField(pcCUAboveRight, uiAboveRightPartIdx, REF_PIC_LIST_0, &pcMvFieldNeighbours[iCount<<1])
        if this.GetSlice().IsInterB() {
            pcCUAboveRight.GetMvField(pcCUAboveRight, uiAboveRightPartIdx, REF_PIC_LIST_1, &pcMvFieldNeighbours[(iCount<<1)+1])
        }
        if mrgCandIdx == iCount {
            return
        }
        iCount++
    }
    // early termination
    if iCount == int(this.GetSlice().GetMaxNumMergeCand()) {
        return
    }

    //left bottom
    uiLeftBottomPartIdx := uint(0)
    var pcCULeftBottom *TComDataCU
    pcCULeftBottom = this.GetPUBelowLeft(&uiLeftBottomPartIdx, uiPartIdxLB, true)
    isAvailableA0 := pcCULeftBottom != nil &&
        pcCULeftBottom.IsDiffMER(xP-1, yP+nPSH, xP, yP) &&
        !pcCULeftBottom.IsIntra(uiLeftBottomPartIdx)
    if isAvailableA0 && (!isAvailableA1 || !pcCULeft.HasEqualMotion(uiLeftPartIdx, pcCULeftBottom, uiLeftBottomPartIdx)) {
        abCandIsInter[iCount] = true
        // this.Get Inter Dir
        puhInterDirNeighbours[iCount] = pcCULeftBottom.GetInterDir1(uiLeftBottomPartIdx)
        // this.Get Mv from Left
        pcCULeftBottom.GetMvField(pcCULeftBottom, uiLeftBottomPartIdx, REF_PIC_LIST_0, &pcMvFieldNeighbours[iCount<<1])
        if this.GetSlice().IsInterB() {
            pcCULeftBottom.GetMvField(pcCULeftBottom, uiLeftBottomPartIdx, REF_PIC_LIST_1, &pcMvFieldNeighbours[(iCount<<1)+1])
        }
        if mrgCandIdx == iCount {
            return
        }
        iCount++
    }
    // early termination
    if iCount == int(this.GetSlice().GetMaxNumMergeCand()) {
        return
    }
    // above left
    if iCount < 4 {
        uiAboveLeftPartIdx := uint(0)
        var pcCUAboveLeft *TComDataCU
        pcCUAboveLeft = this.GetPUAboveLeft(&uiAboveLeftPartIdx, uiAbsPartAddr, true)
        isAvailableB2 := pcCUAboveLeft != nil &&
            pcCUAboveLeft.IsDiffMER(xP-1, yP-1, xP, yP) &&
            !pcCUAboveLeft.IsIntra(uiAboveLeftPartIdx)
        if isAvailableB2 && (!isAvailableA1 || !pcCULeft.HasEqualMotion(uiLeftPartIdx, pcCUAboveLeft, uiAboveLeftPartIdx)) &&
            (!isAvailableB1 || !pcCUAbove.HasEqualMotion(uiAbovePartIdx, pcCUAboveLeft, uiAboveLeftPartIdx)) {
            abCandIsInter[iCount] = true
            // this.Get Inter Dir
            puhInterDirNeighbours[iCount] = pcCUAboveLeft.GetInterDir1(uiAboveLeftPartIdx)
            // this.Get Mv from Left
            pcCUAboveLeft.GetMvField(pcCUAboveLeft, uiAboveLeftPartIdx, REF_PIC_LIST_0, &pcMvFieldNeighbours[iCount<<1])
            if this.GetSlice().IsInterB() {
                pcCUAboveLeft.GetMvField(pcCUAboveLeft, uiAboveLeftPartIdx, REF_PIC_LIST_1, &pcMvFieldNeighbours[(iCount<<1)+1])
            }
            if mrgCandIdx == iCount {
                return
            }
            iCount++
        }
    }
    // early termination
    if iCount == int(this.GetSlice().GetMaxNumMergeCand()) {
        return
    }
    if this.GetSlice().GetEnableTMVPFlag() {
        //>> MTK colocated-RightBottom
        var uiPartIdxRB uint
        uiLCUIdx := int(this.GetAddr())

        this.DeriveRightBottomIdx(uiPUIdx, &uiPartIdxRB)

        uiAbsPartIdxTmp := G_auiZscanToRaster[uiPartIdxRB]
        uiNumPartInCUWidth := this.m_pcPic.GetNumPartInWidth()

        var cColMv TComMv
        var iRefIdx int

        if (this.m_pcPic.GetCU(this.m_uiCUAddr).GetCUPelX() + G_auiRasterToPelX[uiAbsPartIdxTmp] + this.m_pcPic.GetMinCUWidth()) >= this.m_pcSlice.GetSPS().GetPicWidthInLumaSamples() { // image boundary check
            uiLCUIdx = -1
        } else if (this.m_pcPic.GetCU(this.m_uiCUAddr).GetCUPelY() + G_auiRasterToPelY[uiAbsPartIdxTmp] + this.m_pcPic.GetMinCUHeight()) >= this.m_pcSlice.GetSPS().GetPicHeightInLumaSamples() {
            uiLCUIdx = -1
        } else {
            if (uiAbsPartIdxTmp%uiNumPartInCUWidth < uiNumPartInCUWidth-1) && // is not at the last column of LCU
                (uiAbsPartIdxTmp/uiNumPartInCUWidth < this.m_pcPic.GetNumPartInHeight()-1) { // is not at the last row    of LCU
                uiAbsPartAddr = G_auiRasterToZscan[uiAbsPartIdxTmp+uiNumPartInCUWidth+1]
                uiLCUIdx = int(this.GetAddr())
            } else if uiAbsPartIdxTmp%uiNumPartInCUWidth < uiNumPartInCUWidth-1 { // is not at the last column of LCU But is last row of LCU
                uiAbsPartAddr = G_auiRasterToZscan[(uiAbsPartIdxTmp+uiNumPartInCUWidth+1)%this.m_pcPic.GetNumPartInCU()]
                uiLCUIdx = -1
            } else if uiAbsPartIdxTmp/uiNumPartInCUWidth < this.m_pcPic.GetNumPartInHeight()-1 { // is not at the last row of LCU But is last column of LCU
                uiAbsPartAddr = G_auiRasterToZscan[uiAbsPartIdxTmp+1]
                uiLCUIdx = int(this.GetAddr()) + 1
            } else { //is the right bottom corner of LCU
                uiAbsPartAddr = 0
                uiLCUIdx = -1
            }
        }
        iRefIdx = 0

        bExistMV := false
        var uiPartIdxCenter uint
        uiCurLCUIdx := this.GetAddr()
        dir := 0
        uiArrayAddr := uint(iCount)
        this.xDeriveCenterIdx(uiPUIdx, &uiPartIdxCenter)
        bExistMV = uiLCUIdx >= 0 && this.xGetColMVP(REF_PIC_LIST_0, uiLCUIdx, int(uiAbsPartAddr), &cColMv, &iRefIdx)
        if bExistMV == false {
            bExistMV = this.xGetColMVP(REF_PIC_LIST_0, int(uiCurLCUIdx), int(uiPartIdxCenter), &cColMv, &iRefIdx)
        }
        if bExistMV {
            dir |= 1
            pcMvFieldNeighbours[2*uiArrayAddr].SetMvField(cColMv, int8(iRefIdx))
        }
        if this.GetSlice().IsInterB() {
            bExistMV = uiLCUIdx >= 0 && this.xGetColMVP(REF_PIC_LIST_1, uiLCUIdx, int(uiAbsPartAddr), &cColMv, &iRefIdx)
            if bExistMV == false {
                bExistMV = this.xGetColMVP(REF_PIC_LIST_1, int(uiCurLCUIdx), int(uiPartIdxCenter), &cColMv, &iRefIdx)
            }
            if bExistMV {
                dir |= 2
                pcMvFieldNeighbours[2*uiArrayAddr+1].SetMvField(cColMv, int8(iRefIdx))
            }
        }

        if dir != 0 {
            puhInterDirNeighbours[uiArrayAddr] = byte(dir)
            abCandIsInter[uiArrayAddr] = true

            if mrgCandIdx == iCount {
                return
            }
            iCount++
        }
    }
    // early termination
    if iCount == int(this.GetSlice().GetMaxNumMergeCand()) {
        return
    }
    uiArrayAddr := iCount
    uiCutoff := uiArrayAddr

    if this.GetSlice().IsInterB() {
        var uiPriorityList0 = [12]uint{0, 1, 0, 2, 1, 2, 0, 3, 1, 3, 2, 3}
        var uiPriorityList1 = [12]uint{1, 0, 2, 0, 2, 1, 3, 0, 3, 1, 3, 2}

        for idx := int(0); idx < uiCutoff*(uiCutoff-1) && uiArrayAddr != int(this.GetSlice().GetMaxNumMergeCand()); idx++ {
            i := uiPriorityList0[idx]
            j := uiPriorityList1[idx]
            if abCandIsInter[i] && abCandIsInter[j] && (puhInterDirNeighbours[i]&0x1) != 0 && (puhInterDirNeighbours[j]&0x2) != 0 {
                abCandIsInter[uiArrayAddr] = true
                puhInterDirNeighbours[uiArrayAddr] = 3

                // this.Get Mv from cand[i] and cand[j]
                pcMvFieldNeighbours[uiArrayAddr<<1].SetMvField(pcMvFieldNeighbours[i<<1].GetMv(), pcMvFieldNeighbours[i<<1].GetRefIdx())
                pcMvFieldNeighbours[(uiArrayAddr<<1)+1].SetMvField(pcMvFieldNeighbours[(j<<1)+1].GetMv(), pcMvFieldNeighbours[(j<<1)+1].GetRefIdx())

                iRefPOCL0 := this.m_pcSlice.GetRefPOC(REF_PIC_LIST_0, int(pcMvFieldNeighbours[(uiArrayAddr<<1)].GetRefIdx()))
                iRefPOCL1 := this.m_pcSlice.GetRefPOC(REF_PIC_LIST_1, int(pcMvFieldNeighbours[(uiArrayAddr<<1)+1].GetRefIdx()))
                if iRefPOCL0 == iRefPOCL1 && pcMvFieldNeighbours[(uiArrayAddr<<1)].GetMv() == pcMvFieldNeighbours[(uiArrayAddr<<1)+1].GetMv() {
                    abCandIsInter[uiArrayAddr] = false
                } else {
                    uiArrayAddr++
                }
            }
        }
    }
    // early termination
    if uiArrayAddr == int(this.GetSlice().GetMaxNumMergeCand()) {
        return
    }

    var iNumRefIdx int
    if this.GetSlice().IsInterB() {
        if this.m_pcSlice.GetNumRefIdx(REF_PIC_LIST_0) < this.m_pcSlice.GetNumRefIdx(REF_PIC_LIST_1) {
            iNumRefIdx = this.m_pcSlice.GetNumRefIdx(REF_PIC_LIST_0)
        } else {
            iNumRefIdx = this.m_pcSlice.GetNumRefIdx(REF_PIC_LIST_1)
        }
    } else {
        iNumRefIdx = this.m_pcSlice.GetNumRefIdx(REF_PIC_LIST_0)
    }

    r := int(0)
    refcnt := int(0)
    for uiArrayAddr < int(this.GetSlice().GetMaxNumMergeCand()) {
        abCandIsInter[uiArrayAddr] = true
        puhInterDirNeighbours[uiArrayAddr] = 1
        pcMvFieldNeighbours[uiArrayAddr<<1].SetMvField(TComMv{0, 0}, int8(r))

        if this.GetSlice().IsInterB() {
            puhInterDirNeighbours[uiArrayAddr] = 3
            pcMvFieldNeighbours[(uiArrayAddr<<1)+1].SetMvField(TComMv{0, 0}, int8(r))
        }
        uiArrayAddr++
        if refcnt == iNumRefIdx-1 {
            r = 0
        } else {
            r++
            refcnt++
        }
    }

    *numValidMergeCand = uiArrayAddr
}

// -------------------------------------------------------------------------------------------------------------------
// member functions for modes
// -------------------------------------------------------------------------------------------------------------------

func (this *TComDataCU) IsIntra(uiPartIdx uint) bool {
    return this.m_pePredMode[uiPartIdx] == MODE_INTRA
}
func (this *TComDataCU) IsSkipped(uiPartIdx uint) bool {
    return this.GetSkipFlag1(uiPartIdx)
}   ///< SKIP (no residual)
func (this *TComDataCU) IsBipredRestriction(puIdx uint) bool {
    width := int(0)
    height := int(0)
    var partAddr uint

    this.GetPartIndexAndSize(puIdx, &partAddr, &width, &height)
    if this.GetWidth1(0) == 8 && (width < 8 || height < 8) {
        return true
    }
    return false
}

// -------------------------------------------------------------------------------------------------------------------
// member functions for symbol prediction (most probable / mode conversion)
// -------------------------------------------------------------------------------------------------------------------

func (this *TComDataCU) GetIntraSizeIdx(uiAbsPartIdx uint) uint {
    var uiShift uint

    //uiShift := ( (m_puhTrIdx[uiAbsPartIdx]==0) && (m_pePartSize[uiAbsPartIdx]==SIZE_NxN) ) ? m_puhTrIdx[uiAbsPartIdx]+1 : m_puhTrIdx[uiAbsPartIdx];
    if this.m_pePartSize[uiAbsPartIdx] == SIZE_NxN {
        uiShift = 1
    } else {
        uiShift = 0
    }

    uiWidth := this.m_puhWidth[uiAbsPartIdx] >> uiShift
    uiCnt := uint(0)
    for uiWidth != 0 {
        uiCnt++
        uiWidth >>= 1
    }
    uiCnt -= 2

    if uiCnt > 6 {
        return 6
    }

    return uiCnt
}

func (this *TComDataCU) GetAllowedChromaDir(uiAbsPartIdx uint, uiModeList []uint) {
    uiModeList[0] = PLANAR_IDX
    uiModeList[1] = VER_IDX
    uiModeList[2] = HOR_IDX
    uiModeList[3] = DC_IDX
    uiModeList[4] = DM_CHROMA_IDX

    uiLumaMode := uint(this.GetLumaIntraDir1(uiAbsPartIdx))

    for i := 0; i < NUM_CHROMA_MODE-1; i++ {
        if uiLumaMode == uiModeList[i] {
            uiModeList[i] = 34 // VER+8 mode
            break
        }
    }
}

func (this *TComDataCU) GetIntraDirLumaPredictor(uiAbsPartIdx uint, uiIntraDirPred []int, piMode *int) int {
    var pcTempCU *TComDataCU
    var uiTempPartIdx uint
    var iLeftIntraDir, iAboveIntraDir int
    uiPredNum := 0

    // Get intra direction of left PU
    pcTempCU = this.GetPULeft(&uiTempPartIdx, this.m_uiAbsIdxInLCU+uiAbsPartIdx, true, true)

    if pcTempCU != nil {
        if pcTempCU.IsIntra(uiTempPartIdx) {
            iLeftIntraDir = int(pcTempCU.GetLumaIntraDir1(uiTempPartIdx))
        } else {
            iLeftIntraDir = DC_IDX
        }
    } else {
        iLeftIntraDir = DC_IDX
    }

    // Get intra direction of above PU
    pcTempCU = this.GetPUAbove(&uiTempPartIdx, this.m_uiAbsIdxInLCU+uiAbsPartIdx, true, true, true)

    if pcTempCU != nil {
        if pcTempCU.IsIntra(uiTempPartIdx) {
            iAboveIntraDir = int(pcTempCU.GetLumaIntraDir1(uiTempPartIdx))
        } else {
            iAboveIntraDir = DC_IDX
        }
    } else {
        iAboveIntraDir = DC_IDX
    }

    uiPredNum = 3
    if iLeftIntraDir == iAboveIntraDir {
        if piMode != nil {
            *piMode = 1
        }

        if iLeftIntraDir > 1 { // angular modes
            uiIntraDirPred[0] = iLeftIntraDir
            uiIntraDirPred[1] = ((iLeftIntraDir + 29) % 32) + 2
            uiIntraDirPred[2] = ((iLeftIntraDir - 1) % 32) + 2
        } else { //non-angular
            uiIntraDirPred[0] = PLANAR_IDX
            uiIntraDirPred[1] = DC_IDX
            uiIntraDirPred[2] = VER_IDX
        }
    } else {
        if piMode != nil {
            *piMode = 2
        }
        uiIntraDirPred[0] = iLeftIntraDir
        uiIntraDirPred[1] = iAboveIntraDir

        if iLeftIntraDir != 0 && iAboveIntraDir != 0 { //both modes are non-planar
            uiIntraDirPred[2] = PLANAR_IDX
        } else {
            if (iLeftIntraDir + iAboveIntraDir) < 2 {
                uiIntraDirPred[2] = VER_IDX
            } else {
                uiIntraDirPred[2] = DC_IDX
            }
        }
    }

    return uiPredNum
}

// -------------------------------------------------------------------------------------------------------------------
// member functions for SBAC context
// -------------------------------------------------------------------------------------------------------------------

func (this *TComDataCU) GetCtxSplitFlag(uiAbsPartIdx, uiDepth uint) uint {
    var pcTempCU *TComDataCU
    var uiTempPartIdx, uiCtx uint

    // Get left split flag
    //#if DEPENDENT_SLICES
    pcTempCU = this.GetPULeft(&uiTempPartIdx, this.m_uiAbsIdxInLCU+uiAbsPartIdx, true, true)

    if pcTempCU != nil {
        if uint(pcTempCU.GetDepth1(uiTempPartIdx)) > uiDepth {
            uiCtx = 1
        } else {
            uiCtx = 0
        }
    } else {
        uiCtx = 0
    }

    // Get above split flag
    pcTempCU = this.GetPUAbove(&uiTempPartIdx, this.m_uiAbsIdxInLCU+uiAbsPartIdx, true, false, true)

    if pcTempCU != nil {
        if uint(pcTempCU.GetDepth1(uiTempPartIdx)) > uiDepth {
            uiCtx += 1
        } else {
            uiCtx += 0
        }
    } else {
        uiCtx += 0
    }

    return uiCtx
}

func (this *TComDataCU) GetCtxQtCbf(eType TextType, uiTrDepth uint) uint {
    if eType != 0 {
        return uiTrDepth
    }

    if uiTrDepth == 0 {
        return 1
    }

    return 0
}

func (this *TComDataCU) GetCtxSkipFlag(uiAbsPartIdx uint) uint {
    var pcTempCU *TComDataCU
    var uiTempPartIdx uint
    uiCtx := uint(0)

    // Get BCBP of left PU
    pcTempCU = this.GetPULeft(&uiTempPartIdx, this.m_uiAbsIdxInLCU+uiAbsPartIdx, true, true)

    if pcTempCU != nil {
        uiCtx = uint(B2U(pcTempCU.IsSkipped(uiTempPartIdx)))
    } else {
        uiCtx = 0
    }

    // Get BCBP of above PU
    pcTempCU = this.GetPUAbove(&uiTempPartIdx, this.m_uiAbsIdxInLCU+uiAbsPartIdx, true, false, true)

    if pcTempCU != nil {
        uiCtx += uint(B2U(pcTempCU.IsSkipped(uiTempPartIdx)))
    } else {
        uiCtx += 0
    }

    return uiCtx
}
func (this *TComDataCU) GetCtxInterDir(uiAbsPartIdx uint) uint {
    return uint(this.GetDepth1(uiAbsPartIdx))
}

func (this *TComDataCU) GetSliceStartCU(pos uint) uint {
    return this.m_sliceStartCU[pos-this.m_uiAbsIdxInLCU]
}
func (this *TComDataCU) GetSliceSegmentStartCU(pos uint) uint {
    return this.m_sliceSegmentStartCU[pos-this.m_uiAbsIdxInLCU]
}
func (this *TComDataCU) GetTotalBins() uint {
    return this.m_uiTotalBins
}
func (this *TComDataCU) SetTotalBins(uiTotalBins uint) {
    this.m_uiTotalBins = uiTotalBins
}

// -------------------------------------------------------------------------------------------------------------------
// member functions for RD cost storage
// -------------------------------------------------------------------------------------------------------------------

func (this *TComDataCU) GetTotalCost() float64 {
    return this.m_dTotalCost
}
func (this *TComDataCU) GetTotalDistortion() uint {
    return this.m_uiTotalDistortion
}
func (this *TComDataCU) GetTotalBits() uint {
    return this.m_uiTotalBits
}
func (this *TComDataCU) GetTotalNumPart() uint {
    return this.m_uiNumPartition
}

func (this *TComDataCU) SetTotalCost(dTotalCost float64) {
    this.m_dTotalCost = dTotalCost
}
func (this *TComDataCU) SetTotalDistortion(uiTotalDistortion uint) {
    this.m_uiTotalDistortion = uiTotalDistortion
}
func (this *TComDataCU) SetTotalBits(uiTotalBits uint) {
    this.m_uiTotalBits = uiTotalBits
}
func (this *TComDataCU) SetTotalNumPart(uiNumPartition uint) {
    this.m_uiNumPartition = uiNumPartition
}

func (this *TComDataCU) GetCoefScanIdx(uiAbsPartIdx, uiWidth uint, bIsLuma, bIsIntra bool) uint {
    var uiCTXIdx, uiScanIdx, uiDirMode uint

    if !bIsIntra {
        uiScanIdx = SCAN_DIAG
        return uiScanIdx
    }

    switch uiWidth {
    case 2:
        uiCTXIdx = 6
    case 4:
        uiCTXIdx = 5
    case 8:
        uiCTXIdx = 4
    case 16:
        uiCTXIdx = 3
    case 32:
        uiCTXIdx = 2
    case 64:
        uiCTXIdx = 1
    default:
        uiCTXIdx = 0
    }

    if bIsLuma {
        uiDirMode = uint(this.GetLumaIntraDir1(uiAbsPartIdx))
        uiScanIdx = SCAN_DIAG
        if uiCTXIdx > 3 && uiCTXIdx < 6 { //if multiple scans supported for transform size
            var a, b int
            if int(uiDirMode)-VER_IDX < 0 {
                a = VER_IDX - int(uiDirMode)
            } else {
                a = int(uiDirMode) - VER_IDX
            }
            if int(uiDirMode)-HOR_IDX < 0 {
                b = HOR_IDX - int(uiDirMode)
            } else {
                b = int(uiDirMode) - HOR_IDX
            }

            if a < 5 {
                uiScanIdx = SCAN_HOR
            } else if b < 5 {
                uiScanIdx = SCAN_VER
            } else {
                uiScanIdx = SCAN_DIAG
            }
        }
    } else {
        uiDirMode = uint(this.GetChromaIntraDir1(uiAbsPartIdx))
        if uiDirMode == DM_CHROMA_IDX {
            // this.Get number of partitions in current CU
            depth := this.GetDepth1(uiAbsPartIdx)
            numParts := this.GetPic().GetNumPartInCU() >> (2 * depth)

            // this.Get luma mode from upper-left corner of current CU
            uiDirMode = uint(this.GetLumaIntraDir1((uiAbsPartIdx / numParts) * numParts))
        }

        uiScanIdx = SCAN_DIAG
        if uiCTXIdx > 4 && uiCTXIdx < 7 { //if multiple scans supported for transform size
            var a, b int
            if int(uiDirMode)-VER_IDX < 0 {
                a = VER_IDX - int(uiDirMode)
            } else {
                a = int(uiDirMode) - VER_IDX
            }
            if int(uiDirMode)-HOR_IDX < 0 {
                b = HOR_IDX - int(uiDirMode)
            } else {
                b = int(uiDirMode) - HOR_IDX
            }

            if a < 5 {
                uiScanIdx = SCAN_HOR
            } else if b < 5 {
                uiScanIdx = SCAN_VER
            } else {
                uiScanIdx = SCAN_DIAG
            }
        }
    }

    return uiScanIdx
}

//namespace RasterAddress
//{
/** Check whether 2 addresses point to the same column
 * \param addrA          First address in raster scan order
 * \param addrB          Second address in raters scan order
 * \param numUnitsPerRow Number of units in a row
 * \return Result of test
 */
func IsEqualCol(addrA, addrB, numUnitsPerRow int) bool {
    // addrA % numUnitsPerRow == addrB % numUnitsPerRow
    return ((addrA ^ addrB) & (numUnitsPerRow - 1)) == 0
}

/** Check whether 2 addresses point to the same row
 * \param addrA          First address in raster scan order
 * \param addrB          Second address in raters scan order
 * \param numUnitsPerRow Number of units in a row
 * \return Result of test
 */
func IsEqualRow(addrA, addrB, numUnitsPerRow int) bool {
    // addrA / numUnitsPerRow == addrB / numUnitsPerRow
    return ((addrA ^ addrB) & (^(numUnitsPerRow - 1))) == 0
}

/** Check whether 2 addresses point to the same row or column
 * \param addrA          First address in raster scan order
 * \param addrB          Second address in raters scan order
 * \param numUnitsPerRow Number of units in a row
 * \return Result of test
 */
func IsEqualRowOrCol(addrA, addrB, numUnitsPerRow int) bool {
    return IsEqualCol(addrA, addrB, numUnitsPerRow) || IsEqualRow(addrA, addrB, numUnitsPerRow)
}

/** Check whether one address points to the first column
 * \param addr           Address in raster scan order
 * \param numUnitsPerRow Number of units in a row
 * \return Result of test
 */
func IsZeroCol(addr, numUnitsPerRow int) bool {
    // addr % numUnitsPerRow == 0
    return (addr & (numUnitsPerRow - 1)) == 0
}

/** Check whether one address points to the first row
 * \param addr           Address in raster scan order
 * \param numUnitsPerRow Number of units in a row
 * \return Result of test
 */
func IsZeroRow(addr, numUnitsPerRow int) bool {
    // addr / numUnitsPerRow == 0
    return (addr & ^(numUnitsPerRow - 1)) == 0
}

/** Check whether one address points to a column whose index is smaller than a given value
 * \param addr           Address in raster scan order
 * \param val            Given column index value
 * \param numUnitsPerRow Number of units in a row
 * \return Result of test
 */
func LessThanCol(addr, val, numUnitsPerRow int) bool {
    // addr % numUnitsPerRow < val
    return (addr & (numUnitsPerRow - 1)) < val
}

/** Check whether one address points to a row whose index is smaller than a given value
 * \param addr           Address in raster scan order
 * \param val            Given row index value
 * \param numUnitsPerRow Number of units in a row
 * \return Result of test
 */
func LessThanRow(addr, val, numUnitsPerRow int) bool {
    // addr / numUnitsPerRow < val
    return addr < val*numUnitsPerRow
}

//};
