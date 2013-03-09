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
    "container/list"
    "fmt"
)

// ====================================================================================================================
// Class definition
// ====================================================================================================================

/// Reference Picture Set class
type TComReferencePictureSet struct {
    //private:
    m_numberOfPictures         int
    m_numberOfNegativePictures int
    m_numberOfPositivePictures int
    m_numberOfLongtermPictures int
    m_deltaPOC                 [MAX_NUM_REF_PICS]int
    m_POC                      [MAX_NUM_REF_PICS]int
    m_used                     [MAX_NUM_REF_PICS]bool
    m_interRPSPrediction       bool
    m_deltaRIdxMinus1          int
    m_deltaRPS                 int
    m_numRefIdc                int
    m_refIdc                   [MAX_NUM_REF_PICS + 1]int
    m_bCheckLTMSB              [MAX_NUM_REF_PICS]bool
    m_pocLSBLT                 [MAX_NUM_REF_PICS]int
    m_deltaPOCMSBCycleLT       [MAX_NUM_REF_PICS]int
    m_deltaPocMSBPresentFlag   [MAX_NUM_REF_PICS]bool
}

//public:
func NewTComReferencePictureSet() *TComReferencePictureSet {
    return &TComReferencePictureSet{} //init by go
}

func (this *TComReferencePictureSet) GetPocLSBLT(i int) int {
    return this.m_pocLSBLT[i]
}
func (this *TComReferencePictureSet) SetPocLSBLT(i, x int) {
    this.m_pocLSBLT[i] = x
}
func (this *TComReferencePictureSet) GetDeltaPocMSBCycleLT(i int) int {
    return this.m_deltaPOCMSBCycleLT[i]
}
func (this *TComReferencePictureSet) SetDeltaPocMSBCycleLT(i, x int) {
    this.m_deltaPOCMSBCycleLT[i] = x
}
func (this *TComReferencePictureSet) GetDeltaPocMSBPresentFlag(i int) bool {
    return this.m_deltaPocMSBPresentFlag[i]
}
func (this *TComReferencePictureSet) SetDeltaPocMSBPresentFlag(i int, x bool) {
    this.m_deltaPocMSBPresentFlag[i] = x
}
func (this *TComReferencePictureSet) SetUsed(bufferNum int, used bool) {
    this.m_used[bufferNum] = used
}
func (this *TComReferencePictureSet) SetDeltaPOC(bufferNum, deltaPOC int) {
    this.m_deltaPOC[bufferNum] = deltaPOC
}
func (this *TComReferencePictureSet) SetPOC(bufferNum, POC int) {
    this.m_POC[bufferNum] = POC
}
func (this *TComReferencePictureSet) SetNumberOfPictures(numberOfPictures int) {
    this.m_numberOfPictures = numberOfPictures
}
func (this *TComReferencePictureSet) SetCheckLTMSBPresent(bufferNum int, b bool) {
    this.m_bCheckLTMSB[bufferNum] = b
}
func (this *TComReferencePictureSet) GetCheckLTMSBPresent(bufferNum int) bool {
    return this.m_bCheckLTMSB[bufferNum]
}

func (this *TComReferencePictureSet) GetUsed(bufferNum int) bool {
    return this.m_used[bufferNum]
}
func (this *TComReferencePictureSet) GetDeltaPOC(bufferNum int) int {
    return this.m_deltaPOC[bufferNum]
}
func (this *TComReferencePictureSet) GetPOC(bufferNum int) int {
    return this.m_POC[bufferNum]
}
func (this *TComReferencePictureSet) GetNumberOfPictures() int {
    return this.m_numberOfPictures
}

func (this *TComReferencePictureSet) SetNumberOfNegativePictures(number int) {
    this.m_numberOfNegativePictures = number
}
func (this *TComReferencePictureSet) GetNumberOfNegativePictures() int {
    return this.m_numberOfNegativePictures
}
func (this *TComReferencePictureSet) SetNumberOfPositivePictures(number int) {
    this.m_numberOfPositivePictures = number
}
func (this *TComReferencePictureSet) GetNumberOfPositivePictures() int {
    return this.m_numberOfPositivePictures
}
func (this *TComReferencePictureSet) SetNumberOfLongtermPictures(number int) {
    this.m_numberOfLongtermPictures = number
}
func (this *TComReferencePictureSet) GetNumberOfLongtermPictures() int {
    return this.m_numberOfLongtermPictures
}

func (this *TComReferencePictureSet) SetInterRPSPrediction(flag bool) {
    this.m_interRPSPrediction = flag
}
func (this *TComReferencePictureSet) GetInterRPSPrediction() bool {
    return this.m_interRPSPrediction
}
func (this *TComReferencePictureSet) SetDeltaRIdxMinus1(x int) {
    this.m_deltaRIdxMinus1 = x
}
func (this *TComReferencePictureSet) GetDeltaRIdxMinus1() int {
    return this.m_deltaRIdxMinus1
}
func (this *TComReferencePictureSet) SetDeltaRPS(x int) {
    this.m_deltaRPS = x
}
func (this *TComReferencePictureSet) GetDeltaRPS() int {
    return this.m_deltaRPS
}
func (this *TComReferencePictureSet) SetNumRefIdc(x int) {
    this.m_numRefIdc = x
}
func (this *TComReferencePictureSet) GetNumRefIdc() int {
    return this.m_numRefIdc
}

func (this *TComReferencePictureSet) SetRefIdc(bufferNum, refIdc int) {
    this.m_refIdc[bufferNum] = refIdc
}
func (this *TComReferencePictureSet) GetRefIdc(bufferNum int) int {
    return this.m_refIdc[bufferNum]
}

func (this *TComReferencePictureSet) SortDeltaPOC() {
    // sort in increasing order (smallest first)
    for j := 1; j < this.GetNumberOfPictures(); j++ {
        deltaPOC := this.GetDeltaPOC(j)
        used := this.GetUsed(j)
        for k := j - 1; k >= 0; k-- {
            temp := this.GetDeltaPOC(k)
            if deltaPOC < temp {
                this.SetDeltaPOC(k+1, temp)
                this.SetUsed(k+1, this.GetUsed(k))
                this.SetDeltaPOC(k, deltaPOC)
                this.SetUsed(k, used)
            }
        }
    }
    // flip the negative values to largest first
    numNegPics := this.GetNumberOfNegativePictures()
    k := numNegPics - 1
    for j := 0; j < numNegPics>>1; j++ {
        deltaPOC := this.GetDeltaPOC(j)
        used := this.GetUsed(j)
        this.SetDeltaPOC(j, this.GetDeltaPOC(k))
        this.SetUsed(j, this.GetUsed(k))
        this.SetDeltaPOC(k, deltaPOC)
        this.SetUsed(k, used)
        k--
    }
}
func (this *TComReferencePictureSet) PrintDeltaPOC() {
    fmt.Printf("DeltaPOC = { ")
    for j := 0; j < this.GetNumberOfPictures(); j++ {
        if this.GetUsed(j) {
            fmt.Printf("%d%s ", this.GetDeltaPOC(j), "*")
        } else {
            fmt.Printf("%d%s ", this.GetDeltaPOC(j), "")
        }
    }
    if this.GetInterRPSPrediction() {
        fmt.Printf("}, RefIdc = { ")
        for j := 0; j < this.GetNumRefIdc(); j++ {
            fmt.Printf("%d ", this.GetRefIdc(j))
        }
    }
    fmt.Printf("}\n")
}

//};

/// Reference Picture Set set class
type TComRPSList struct {
    //private:
    m_numberOfReferencePictureSets int
    m_referencePictureSets         []TComReferencePictureSet
}

//public:
func NewTComRPSList() *TComRPSList {
    return &TComRPSList{}
}

func (this *TComRPSList) Create(numberOfReferencePictureSets int) {
    this.m_numberOfReferencePictureSets = numberOfReferencePictureSets
    this.m_referencePictureSets = make([]TComReferencePictureSet, numberOfReferencePictureSets)
}
func (this *TComRPSList) Destroy() {
    //if this.m_referencePictureSets
    //{
    //  delete [] m_referencePictureSets;
    //}
    this.m_numberOfReferencePictureSets = 0
    this.m_referencePictureSets = nil
}

func (this *TComRPSList) GetReferencePictureSet(referencePictureSetNum int) *TComReferencePictureSet {
    return &this.m_referencePictureSets[referencePictureSetNum]
}
func (this *TComRPSList) GetNumberOfReferencePictureSets() int {
    return this.m_numberOfReferencePictureSets
}
func (this *TComRPSList) SetNumberOfReferencePictureSets(numberOfReferencePictureSets int) {
    this.m_numberOfReferencePictureSets = numberOfReferencePictureSets
}

/// SCALING_LIST class
type TComScalingList struct {
    m_scalingListDC               [SCALING_LIST_SIZE_NUM][SCALING_LIST_NUM]int   //!< the DC value of the matrix coefficient for 16x16
    m_useDefaultScalingMatrixFlag [SCALING_LIST_SIZE_NUM][SCALING_LIST_NUM]bool  //!< UseDefaultScalingMatrixFlag
    m_refMatrixId                 [SCALING_LIST_SIZE_NUM][SCALING_LIST_NUM]uint  //!< RefMatrixID
    m_scalingListPresentFlag      bool                                           //!< flag for using default matrix
    m_predMatrixId                [SCALING_LIST_SIZE_NUM][SCALING_LIST_NUM]uint  //!< reference list index
    m_scalingListCoef             [SCALING_LIST_SIZE_NUM][SCALING_LIST_NUM][]int //!< quantization matrix
    m_useTransformSkip            bool
}

//public:
func NewTComScalingList() *TComScalingList {
    sl := &TComScalingList{m_useTransformSkip: false}
    sl.Init()

    return sl
}

func (this *TComScalingList) SetScalingListPresentFlag(b bool) {
    this.m_scalingListPresentFlag = b
}
func (this *TComScalingList) GetScalingListPresentFlag() bool {
    return this.m_scalingListPresentFlag
}
func (this *TComScalingList) GetUseTransformSkip() bool {
    return this.m_useTransformSkip
}
func (this *TComScalingList) SetUseTransformSkip(b bool) {
    this.m_useTransformSkip = b
}
func (this *TComScalingList) GetScalingListAddress(sizeId, listId uint) []int {
    return this.m_scalingListCoef[sizeId][listId][:]
}                                                                      //!< get matrix coefficient
func (this *TComScalingList) CheckPredMode(sizeId, listId uint) bool { //encoder func
    /*
       for predListIdx := int(listId) ; predListIdx >= 0; predListIdx-- {
         if( !memcmp(getScalingListAddress(sizeId,listId),((listId == predListIdx) ?
           getScalingListDefaultAddress(sizeId, predListIdx): getScalingListAddress(sizeId, predListIdx)),sizeof(Int)*min(MAX_MATRIX_COEF_NUM,(Int)g_scalingListSize[sizeId])) // check value of matrix
          && ((sizeId < SCALING_LIST_16x16) || (getScalingListDC(sizeId,listId) == getScalingListDC(sizeId,predListIdx)))) // check DC value
         {
           setRefMatrixId(sizeId, listId, predListIdx);
           return false;
         }
       }*/
    return true
}
func (this *TComScalingList) SetRefMatrixId(sizeId, listId, u uint) {
    this.m_refMatrixId[sizeId][listId] = u
}   //!< set reference matrix ID
func (this *TComScalingList) GetRefMatrixId(sizeId, listId uint) uint {
    return this.m_refMatrixId[sizeId][listId]
}   //!< get reference matrix ID
func (this *TComScalingList) GetScalingListDefaultAddress(sizeId, listId uint) []int {
    var src []int
    switch sizeId {
    case SCALING_LIST_4x4:
        src = G_quantTSDefault4x4[:]
    case SCALING_LIST_8x8:
        if listId < 3 {
            src = G_quantIntraDefault8x8[:]
        } else {
            src = G_quantInterDefault8x8[:]
        }
    case SCALING_LIST_16x16:
        if listId < 3 {
            src = G_quantIntraDefault8x8[:]
        } else {
            src = G_quantInterDefault8x8[:]
        }
    case SCALING_LIST_32x32:
        if listId < 1 {
            src = G_quantIntraDefault8x8[:]
        } else {
            src = G_quantInterDefault8x8[:]
        }
    default:
        src = nil
    }
    return src
}   //!< get default matrix coefficient
func (this *TComScalingList) ProcessDefaultMarix(sizeId, listId uint) {
    for i := 0; i < MIN(MAX_MATRIX_COEF_NUM, int(G_scalingListSize[sizeId])).(int); i++ {
        this.GetScalingListAddress(sizeId, listId)[i] = this.GetScalingListDefaultAddress(sizeId, listId)[i]
    }
    //::memcpy(getScalingListAddress(sizeId, listId),getScalingListDefaultAddress(sizeId,listId),sizeof(Int)*);

    this.SetScalingListDC(sizeId, listId, SCALING_LIST_DC)
}
func (this *TComScalingList) SetScalingListDC(sizeId, listId, u uint) {
    this.m_scalingListDC[sizeId][listId] = int(u)
}   //!< set DC value

func (this *TComScalingList) GetScalingListDC(sizeId, listId uint) int {
    return this.m_scalingListDC[sizeId][listId]
}   //!< get DC value
func (this *TComScalingList) CheckDcOfMatrix() {
    for sizeId := uint(0); sizeId < SCALING_LIST_SIZE_NUM; sizeId++ {
        for listId := uint(0); listId < G_scalingListNum[sizeId]; listId++ {
            //check default matrix?
            if this.GetScalingListDC(sizeId, listId) == 0 {
                this.ProcessDefaultMarix(sizeId, listId)
            }
        }
    }
}
func (this *TComScalingList) ProcessRefMatrix(sizeId, listId, refListId uint) {
    if listId == refListId {
        for i := 0; i < MIN(MAX_MATRIX_COEF_NUM, int(G_scalingListSize[sizeId])).(int); i++ {
            this.GetScalingListAddress(sizeId, listId)[i] = this.GetScalingListDefaultAddress(sizeId, refListId)[i]
        }
    } else {
        for i := 0; i < MIN(MAX_MATRIX_COEF_NUM, int(G_scalingListSize[sizeId])).(int); i++ {
            this.GetScalingListAddress(sizeId, listId)[i] = this.GetScalingListAddress(sizeId, refListId)[i]
        }
    }
}
func (this *TComScalingList) XParseScalingList(pchFile string) bool { //Encoder func
    /*FILE *fp;
      Char line[1024];
      UInt sizeIdc,listIdc;
      UInt i,size = 0;
      Int *src=0,data;
      Char *ret;
      UInt  retval;

      if((fp = fopen(pchFile,"r")) == (FILE*)NULL)
      {
        printf("can't open file %s :: set Default Matrix\n",pchFile);
        return true;
      }

      for(sizeIdc = 0; sizeIdc < SCALING_LIST_SIZE_NUM; sizeIdc++)
      {
        size = min(MAX_MATRIX_COEF_NUM,(Int)g_scalingListSize[sizeIdc]);
        for(listIdc = 0; listIdc < g_scalingListNum[sizeIdc]; listIdc++)
        {
          src = getScalingListAddress(sizeIdc, listIdc);

          fseek(fp,0,0);
          do
          {
            ret = fgets(line, 1024, fp);
            if ((ret==NULL)||(strstr(line, MatrixType[sizeIdc][listIdc])==NULL && feof(fp)))
            {
              printf("Error: can't read Matrix :: set Default Matrix\n");
              return true;
            }
          }
          while (strstr(line, MatrixType[sizeIdc][listIdc]) == NULL);
          for (i=0; i<size; i++)
          {
            retval = fscanf(fp, "%d,", &data);
            if (retval!=1)
            {
              printf("Error: can't read Matrix :: set Default Matrix\n");
              return true;
            }
            src[i] = data;
          }
          //set DC value for default matrix check
          setScalingListDC(sizeIdc,listIdc,src[0]);

          if(sizeIdc > SCALING_LIST_8x8)
          {
            fseek(fp,0,0);
            do
            {
              ret = fgets(line, 1024, fp);
              if ((ret==NULL)||(strstr(line, MatrixType_DC[sizeIdc][listIdc])==NULL && feof(fp)))
              {
                printf("Error: can't read DC :: set Default Matrix\n");
                return true;
              }
            }
            while (strstr(line, MatrixType_DC[sizeIdc][listIdc]) == NULL);
            retval = fscanf(fp, "%d,", &data);
            if (retval!=1)
            {
              printf("Error: can't read Matrix :: set Default Matrix\n");
              return true;
            }
            //overwrite DC value when size of matrix is larger than 16x16
            setScalingListDC(sizeIdc,listIdc,data);
          }
        }
      }
      fclose(fp);
    */
    return false
}

//private:
func (this *TComScalingList) Init() {
    for sizeId := uint(0); sizeId < SCALING_LIST_SIZE_NUM; sizeId++ {
        for listId := uint(0); listId < G_scalingListNum[sizeId]; listId++ {
            var l int
            if G_scalingListSize[sizeId] > MAX_MATRIX_COEF_NUM {
                l = MAX_MATRIX_COEF_NUM
            } else {
                l = int(G_scalingListSize[sizeId])
            }
            this.m_scalingListCoef[sizeId][listId] = make([]int, l)
        }
    }
    this.m_scalingListCoef[SCALING_LIST_32x32][3] = this.m_scalingListCoef[SCALING_LIST_32x32][1] // copy address for 32x32
}
func (this *TComScalingList) Destroy() {
    /*for sizeId := 0; sizeId < SCALING_LIST_SIZE_NUM; sizeId++ {
      for listId := 0; listId < g_scalingListNum[sizeId]; listId++ {
        if m_scalingListCoef[sizeId][listId]) delete [] m_scalingListCoef[sizeId][listId];
      }
    }*/
}

//!< transform skipping flag for setting default scaling matrix for 4x4

type ProfileTierLevel struct {
    m_profileSpace             int
    m_tierFlag                 bool
    m_profileIdc               int
    m_profileCompatibilityFlag [32]bool
    m_levelIdc                 int
    m_progressiveSourceFlag    bool
    m_interlacedSourceFlag     bool
    m_nonPackedConstraintFlag  bool
    m_frameOnlyConstraintFlag  bool
}

//public:
func NewProfileTierLevel() *ProfileTierLevel {
    return &ProfileTierLevel{} //init zero by go
}

func (this *ProfileTierLevel) GetProfileSpace() int {
    return this.m_profileSpace
}
func (this *ProfileTierLevel) SetProfileSpace(x int) {
    this.m_profileSpace = x
}

func (this *ProfileTierLevel) GetTierFlag() bool {
    return this.m_tierFlag
}
func (this *ProfileTierLevel) SetTierFlag(x bool) {
    this.m_tierFlag = x
}

func (this *ProfileTierLevel) GetProfileIdc() int {
    return this.m_profileIdc
}
func (this *ProfileTierLevel) SetProfileIdc(x int) {
    this.m_profileIdc = x
}

func (this *ProfileTierLevel) GetProfileCompatibilityFlag(i int) bool {
    return this.m_profileCompatibilityFlag[i]
}
func (this *ProfileTierLevel) SetProfileCompatibilityFlag(i int, x bool) {
    this.m_profileCompatibilityFlag[i] = x
}

func (this *ProfileTierLevel) GetLevelIdc() int {
    return this.m_levelIdc
}
func (this *ProfileTierLevel) SetLevelIdc(x int) {
    this.m_levelIdc = x
}

func (this *ProfileTierLevel) GetProgressiveSourceFlag() bool  { return this.m_progressiveSourceFlag }
func (this *ProfileTierLevel) SetProgressiveSourceFlag(b bool) { this.m_progressiveSourceFlag = b }

func (this *ProfileTierLevel) GetInterlacedSourceFlag() bool  { return this.m_interlacedSourceFlag }
func (this *ProfileTierLevel) SetInterlacedSourceFlag(b bool) { this.m_interlacedSourceFlag = b }

func (this *ProfileTierLevel) GetNonPackedConstraintFlag() bool  { return this.m_nonPackedConstraintFlag }
func (this *ProfileTierLevel) SetNonPackedConstraintFlag(b bool) { this.m_nonPackedConstraintFlag = b }

func (this *ProfileTierLevel) GetFrameOnlyConstraintFlag() bool  { return this.m_frameOnlyConstraintFlag }
func (this *ProfileTierLevel) SetFrameOnlyConstraintFlag(b bool) { this.m_frameOnlyConstraintFlag = b }

type TComPTL struct {
    m_generalPTL                 ProfileTierLevel
    m_subLayerPTL                [6]ProfileTierLevel // max. value of max_sub_layers_minus1 is 6
    m_subLayerProfilePresentFlag [6]bool
    m_subLayerLevelPresentFlag   [6]bool
}

//public:
func NewTComPTL() *TComPTL {
    return &TComPTL{} //init zero by go
}
func (this *TComPTL) GetSubLayerProfilePresentFlag(i int) bool {
    return this.m_subLayerProfilePresentFlag[i]
}
func (this *TComPTL) SetSubLayerProfilePresentFlag(i int, x bool) {
    this.m_subLayerProfilePresentFlag[i] = x
}

func (this *TComPTL) GetSubLayerLevelPresentFlag(i int) bool {
    return this.m_subLayerLevelPresentFlag[i]
}
func (this *TComPTL) SetSubLayerLevelPresentFlag(i int, x bool) {
    this.m_subLayerLevelPresentFlag[i] = x
}

func (this *TComPTL) GetGeneralPTL() *ProfileTierLevel {
    return &this.m_generalPTL
}
func (this *TComPTL) GetSubLayerPTL(i int) *ProfileTierLevel {
    return &this.m_subLayerPTL[i]
}

type HrdSubLayerInfo struct {
    fixedPicRateFlag          bool
    fixedPicRateWithinCvsFlag bool
    picDurationInTcMinus1     uint
    lowDelayHrdFlag           bool
    cpbCntMinus1              uint
    bitRateValueMinus1        [MAX_CPB_CNT][2]uint
    cpbSizeValue              [MAX_CPB_CNT][2]uint
    ducpbSizeValue            [MAX_CPB_CNT][2]uint
    cbrFlag                   [MAX_CPB_CNT][2]bool
    duBitRateValue            [MAX_CPB_CNT][2]uint
}

type TComHRD struct {
    m_nalHrdParametersPresentFlag       bool
    m_vclHrdParametersPresentFlag       bool
    m_subPicCpbParamsPresentFlag        bool
    m_tickDivisorMinus2                 uint
    m_duCpbRemovalDelayLengthMinus1     uint
    m_subPicCpbParamsInPicTimingSEIFlag bool
    m_dpbOutputDelayDuLengthMinus1      uint

    m_bitRateScale                       uint
    m_cpbSizeScale                       uint
    m_ducpbSizeScale                     uint
    m_initialCpbRemovalDelayLengthMinus1 uint
    m_cpbRemovalDelayLengthMinus1        uint
    m_dpbOutputDelayLengthMinus1         uint
    m_numDU                              uint
    m_HRD                                [MAX_TLAYER]HrdSubLayerInfo
}

func NewTComHRD() *TComHRD {
    return &TComHRD{}
}

func (this *TComHRD) SetNalHrdParametersPresentFlag(flag bool) {
    this.m_nalHrdParametersPresentFlag = flag
}
func (this *TComHRD) GetNalHrdParametersPresentFlag() bool {
    return this.m_nalHrdParametersPresentFlag
}

func (this *TComHRD) SetVclHrdParametersPresentFlag(flag bool) {
    this.m_vclHrdParametersPresentFlag = flag
}
func (this *TComHRD) GetVclHrdParametersPresentFlag() bool {
    return this.m_vclHrdParametersPresentFlag
}

func (this *TComHRD) SetSubPicCpbParamsPresentFlag(flag bool) {
    this.m_subPicCpbParamsPresentFlag = flag
}
func (this *TComHRD) GetSubPicCpbParamsPresentFlag() bool {
    return this.m_subPicCpbParamsPresentFlag
}

func (this *TComHRD) SetTickDivisorMinus2(value uint) {
    this.m_tickDivisorMinus2 = value
}
func (this *TComHRD) GetTickDivisorMinus2() uint { return this.m_tickDivisorMinus2 }

func (this *TComHRD) SetDuCpbRemovalDelayLengthMinus1(value uint) {
    this.m_duCpbRemovalDelayLengthMinus1 = value
}
func (this *TComHRD) GetDuCpbRemovalDelayLengthMinus1() uint {
    return this.m_duCpbRemovalDelayLengthMinus1
}

func (this *TComHRD) SetSubPicCpbParamsInPicTimingSEIFlag(flag bool) {
    this.m_subPicCpbParamsInPicTimingSEIFlag = flag
}
func (this *TComHRD) GetSubPicCpbParamsInPicTimingSEIFlag() bool {
    return this.m_subPicCpbParamsInPicTimingSEIFlag
}

func (this *TComHRD) SetDpbOutputDelayDuLengthMinus1(value uint) {
    this.m_dpbOutputDelayDuLengthMinus1 = value
}
func (this *TComHRD) GetDpbOutputDelayDuLengthMinus1() uint {
    return this.m_dpbOutputDelayDuLengthMinus1
}

func (this *TComHRD) SetBitRateScale(value uint) { this.m_bitRateScale = value }
func (this *TComHRD) GetBitRateScale() uint      { return this.m_bitRateScale }

func (this *TComHRD) SetCpbSizeScale(value uint) { this.m_cpbSizeScale = value }
func (this *TComHRD) GetCpbSizeScale() uint      { return this.m_cpbSizeScale }
func (this *TComHRD) SetDuCpbSizeScale(value uint) {
    this.m_ducpbSizeScale = value
}
func (this *TComHRD) GetDuCpbSizeScale() uint { return this.m_ducpbSizeScale }

func (this *TComHRD) SetInitialCpbRemovalDelayLengthMinus1(value uint) {
    this.m_initialCpbRemovalDelayLengthMinus1 = value
}
func (this *TComHRD) GetInitialCpbRemovalDelayLengthMinus1() uint {
    return this.m_initialCpbRemovalDelayLengthMinus1
}

func (this *TComHRD) SetCpbRemovalDelayLengthMinus1(value uint) {
    this.m_cpbRemovalDelayLengthMinus1 = value
}
func (this *TComHRD) GetCpbRemovalDelayLengthMinus1() uint {
    return this.m_cpbRemovalDelayLengthMinus1
}

func (this *TComHRD) SetDpbOutputDelayLengthMinus1(value uint) {
    this.m_dpbOutputDelayLengthMinus1 = value
}
func (this *TComHRD) GetDpbOutputDelayLengthMinus1() uint {
    return this.m_dpbOutputDelayLengthMinus1
}

func (this *TComHRD) SetFixedPicRateFlag(layer int, flag bool) {
    this.m_HRD[layer].fixedPicRateFlag = flag
}
func (this *TComHRD) GetFixedPicRateFlag(layer int) bool {
    return this.m_HRD[layer].fixedPicRateFlag
}

func (this *TComHRD) SetFixedPicRateWithinCvsFlag(layer int, flag bool) {
    this.m_HRD[layer].fixedPicRateWithinCvsFlag = flag
}
func (this *TComHRD) GetFixedPicRateWithinCvsFlag(layer int) bool {
    return this.m_HRD[layer].fixedPicRateWithinCvsFlag
}

func (this *TComHRD) SetPicDurationInTcMinus1(layer int, value uint) {
    this.m_HRD[layer].picDurationInTcMinus1 = value
}
func (this *TComHRD) GetPicDurationInTcMinus1(layer int) uint {
    return this.m_HRD[layer].picDurationInTcMinus1
}

func (this *TComHRD) SetLowDelayHrdFlag(layer int, flag bool) {
    this.m_HRD[layer].lowDelayHrdFlag = flag
}
func (this *TComHRD) GetLowDelayHrdFlag(layer int) bool {
    return this.m_HRD[layer].lowDelayHrdFlag
}

func (this *TComHRD) SetCpbCntMinus1(layer int, value uint) {
    this.m_HRD[layer].cpbCntMinus1 = value
}
func (this *TComHRD) GetCpbCntMinus1(layer int) uint {
    return this.m_HRD[layer].cpbCntMinus1
}

func (this *TComHRD) SetBitRateValueMinus1(layer int, cpbcnt, nalOrVcl int, value uint) {
    this.m_HRD[layer].bitRateValueMinus1[cpbcnt][nalOrVcl] = value
}
func (this *TComHRD) GetBitRateValueMinus1(layer int, cpbcnt, nalOrVcl int) uint {
    return this.m_HRD[layer].bitRateValueMinus1[cpbcnt][nalOrVcl]
}

func (this *TComHRD) SetCpbSizeValueMinus1(layer int, cpbcnt, nalOrVcl int, value uint) {
    this.m_HRD[layer].cpbSizeValue[cpbcnt][nalOrVcl] = value
}
func (this *TComHRD) GetCpbSizeValueMinus1(layer int, cpbcnt, nalOrVcl int) uint {
    return this.m_HRD[layer].cpbSizeValue[cpbcnt][nalOrVcl]
}
func (this *TComHRD) SetDuCpbSizeValueMinus1(layer int, cpbcnt, nalOrVcl int, value uint) {
    this.m_HRD[layer].ducpbSizeValue[cpbcnt][nalOrVcl] = value
}
func (this *TComHRD) GetDuCpbSizeValueMinus1(layer int, cpbcnt, nalOrVcl int) uint {
    return this.m_HRD[layer].ducpbSizeValue[cpbcnt][nalOrVcl]
}

func (this *TComHRD) SetDuBitRateValueMinus1(layer int, cpbcnt, nalOrVcl int, value uint) {
    this.m_HRD[layer].duBitRateValue[cpbcnt][nalOrVcl] = value
}
func (this *TComHRD) GetDuBitRateValueMinus1(layer int, cpbcnt, nalOrVcl int) uint {
    return this.m_HRD[layer].duBitRateValue[cpbcnt][nalOrVcl]
}

func (this *TComHRD) SetCbrFlag(layer int, cpbcnt, nalOrVcl int, value bool) {
    this.m_HRD[layer].cbrFlag[cpbcnt][nalOrVcl] = value
}
func (this *TComHRD) GetCbrFlag(layer int, cpbcnt, nalOrVcl int) bool {
    return this.m_HRD[layer].cbrFlag[cpbcnt][nalOrVcl]
}

func (this *TComHRD) SetNumDU(value uint) { this.m_numDU = value }
func (this *TComHRD) GetNumDU() uint      { return this.m_numDU }

func (this *TComHRD) GetCpbDpbDelaysPresentFlag() bool {
    return this.GetNalHrdParametersPresentFlag() || this.GetVclHrdParametersPresentFlag()
}

type TimingInfo struct {
    m_timingInfoPresentFlag       bool
    m_numUnitsInTick              uint
    m_timeScale                   uint
    m_pocProportionalToTimingFlag bool
    m_numTicksPocDiffOneMinus1    int
}

func NewTimingInfo() *TimingInfo {
    return &TimingInfo{m_timingInfoPresentFlag: false,
        m_numUnitsInTick:              1001,
        m_timeScale:                   60000,
        m_pocProportionalToTimingFlag: false,
        m_numTicksPocDiffOneMinus1:    0}
}

func (this *TimingInfo) SetTimingInfoPresentFlag(flag bool) {
    this.m_timingInfoPresentFlag = flag
}
func (this *TimingInfo) GetTimingInfoPresentFlag() bool {
    return this.m_timingInfoPresentFlag
}

func (this *TimingInfo) SetNumUnitsInTick(value uint) {
    this.m_numUnitsInTick = value
}
func (this *TimingInfo) GetNumUnitsInTick() uint { return this.m_numUnitsInTick }

func (this *TimingInfo) SetTimeScale(value uint) { this.m_timeScale = value }
func (this *TimingInfo) GetTimeScale() uint      { return this.m_timeScale }

func (this *TimingInfo) GetPocProportionalToTimingFlag() bool {
    return this.m_pocProportionalToTimingFlag
}
func (this *TimingInfo) SetPocProportionalToTimingFlag(x bool) {
    this.m_pocProportionalToTimingFlag = x
}

func (this *TimingInfo) GetNumTicksPocDiffOneMinus1() int {
    return this.m_numTicksPocDiffOneMinus1
}
func (this *TimingInfo) SetNumTicksPocDiffOneMinus1(x int) {
    this.m_numTicksPocDiffOneMinus1 = x
}

type TComVPS struct {
    //private:
    m_VPSId                  int
    m_uiMaxTLayers           uint
    m_uiMaxLayers            uint
    m_bTemporalIdNestingFlag bool

    m_numReorderPics       [MAX_TLAYER]uint
    m_uiMaxDecPicBuffering [MAX_TLAYER]uint
    m_uiMaxLatencyIncrease [MAX_TLAYER]uint

    m_numHrdParameters          uint
    m_maxNuhReservedZeroLayerId uint
    m_hrdParameters             []TComHRD
    m_hrdOpSetIdx               []uint
    m_cprmsPresentFlag          []bool
    m_numOpSets                 uint

    m_layerIdIncludedFlag [MAX_VPS_OP_SETS_PLUS1][MAX_VPS_NUH_RESERVED_ZERO_LAYER_ID_PLUS1]bool

    m_pcPTL TComPTL

    m_timingInfo TimingInfo
}

//public:

func NewTComVPS() *TComVPS {
    return &TComVPS{}
}

func (this *TComVPS) CreateHrdParamBuffer() {
    this.m_hrdParameters = make([]TComHRD, this.GetNumHrdParameters())
    this.m_hrdOpSetIdx = make([]uint, this.GetNumHrdParameters())
    this.m_cprmsPresentFlag = make([]bool, this.GetNumHrdParameters())
}

func (this *TComVPS) GetHrdParameters(i uint) *TComHRD     { return &this.m_hrdParameters[i] }
func (this *TComVPS) GetHrdOpSetIdx(i uint) uint           { return this.m_hrdOpSetIdx[i] }
func (this *TComVPS) SetHrdOpSetIdx(val, i uint)           { this.m_hrdOpSetIdx[i] = val }
func (this *TComVPS) GetCprmsPresentFlag(i uint) bool      { return this.m_cprmsPresentFlag[i] }
func (this *TComVPS) SetCprmsPresentFlag(val bool, i uint) { this.m_cprmsPresentFlag[i] = val }

func (this *TComVPS) GetVPSId() int {
    return this.m_VPSId
}
func (this *TComVPS) SetVPSId(i int) {
    this.m_VPSId = i
}

func (this *TComVPS) GetMaxTLayers() uint {
    return this.m_uiMaxTLayers
}
func (this *TComVPS) SetMaxTLayers(t uint) {
    this.m_uiMaxTLayers = t
}

func (this *TComVPS) GetMaxLayers() uint {
    return this.m_uiMaxLayers
}
func (this *TComVPS) SetMaxLayers(l uint) {
    this.m_uiMaxLayers = l
}

func (this *TComVPS) GetTemporalNestingFlag() bool {
    return this.m_bTemporalIdNestingFlag
}
func (this *TComVPS) SetTemporalNestingFlag(t bool) {
    this.m_bTemporalIdNestingFlag = t
}

func (this *TComVPS) SetNumReorderPics(v, tLayer uint) {
    this.m_numReorderPics[tLayer] = v
}
func (this *TComVPS) GetNumReorderPics(tLayer uint) uint {
    return this.m_numReorderPics[tLayer]
}

func (this *TComVPS) SetMaxDecPicBuffering(v, tLayer uint) {
    this.m_uiMaxDecPicBuffering[tLayer] = v
}
func (this *TComVPS) GetMaxDecPicBuffering(tLayer uint) uint {
    return this.m_uiMaxDecPicBuffering[tLayer]
}

func (this *TComVPS) SetMaxLatencyIncrease(v, tLayer uint) {
    this.m_uiMaxLatencyIncrease[tLayer] = v
}
func (this *TComVPS) GetMaxLatencyIncrease(tLayer uint) uint {
    return this.m_uiMaxLatencyIncrease[tLayer]
}

func (this *TComVPS) GetNumHrdParameters() uint {
    return this.m_numHrdParameters
}
func (this *TComVPS) SetNumHrdParameters(v uint) {
    this.m_numHrdParameters = v
}

func (this *TComVPS) GetMaxNuhReservedZeroLayerId() uint {
    return this.m_maxNuhReservedZeroLayerId
}
func (this *TComVPS) SetMaxNuhReservedZeroLayerId(v uint) {
    this.m_maxNuhReservedZeroLayerId = v
}
func (this *TComVPS) GetMaxOpSets() uint  { return this.m_numOpSets }
func (this *TComVPS) SetMaxOpSets(v uint) { this.m_numOpSets = v }
func (this *TComVPS) GetLayerIdIncludedFlag(opIdx, id uint) bool {
    return this.m_layerIdIncludedFlag[opIdx][id]
}
func (this *TComVPS) SetLayerIdIncludedFlag(v bool, opIdx, id uint) {
    this.m_layerIdIncludedFlag[opIdx][id] = v
}

func (this *TComVPS) GetPTL() *TComPTL {
    return &this.m_pcPTL
}

func (this *TComVPS) GetTimingInfo() *TimingInfo { return &this.m_timingInfo }

type Window struct {
    //private:
    m_enabledFlag     bool
    m_winLeftOffset   int
    m_winRightOffset  int
    m_winTopOffset    int
    m_winBottomOffset int
}

func NewWindow() *Window {
    return &Window{}
}

func (this *Window) GetWindowEnabledFlag() bool {
    return this.m_enabledFlag
}
func (this *Window) SetWindowEnabledFlag(enabledFlag bool) {
    this.m_enabledFlag = enabledFlag
}
func (this *Window) ResetWindow() {
    this.m_enabledFlag = false
    this.m_winLeftOffset = 0
    this.m_winRightOffset = 0
    this.m_winTopOffset = 0
    this.m_winBottomOffset = 0
}
func (this *Window) GetWindowLeftOffset() int {
    if this.m_enabledFlag {
        return this.m_winLeftOffset
    }
    return 0
}
func (this *Window) SetWindowLeftOffset(val int) {
    this.m_winLeftOffset = val
    this.m_enabledFlag = true
}
func (this *Window) GetWindowRightOffset() int {
    if this.m_enabledFlag {
        return this.m_winRightOffset
    }
    return 0
}
func (this *Window) SetWindowRightOffset(val int) {
    this.m_winRightOffset = val
    this.m_enabledFlag = true
}
func (this *Window) GetWindowTopOffset() int {
    if this.m_enabledFlag {
        return this.m_winTopOffset
    }
    return 0
}
func (this *Window) SetWindowTopOffset(val int) {
    this.m_winTopOffset = val
    this.m_enabledFlag = true
}
func (this *Window) GetWindowBottomOffset() int {
    if this.m_enabledFlag {
        return this.m_winBottomOffset
    }
    return 0
}
func (this *Window) SetWindowBottomOffset(val int) {
    this.m_winBottomOffset = val
    this.m_enabledFlag = true
}
func (this *Window) SetWindow(offsetLeft, offsetLRight, offsetLTop, offsetLBottom int) {
    this.m_enabledFlag = true
    this.m_winLeftOffset = offsetLeft
    this.m_winRightOffset = offsetLRight
    this.m_winTopOffset = offsetLTop
    this.m_winBottomOffset = offsetLBottom
}

type TComVUI struct {
    //private:
    m_aspectRatioInfoPresentFlag     bool
    m_aspectRatioIdc                 int
    m_sarWidth                       int
    m_sarHeight                      int
    m_overscanInfoPresentFlag        bool
    m_overscanAppropriateFlag        bool
    m_videoSignalTypePresentFlag     bool
    m_videoFormat                    int
    m_videoFullRangeFlag             bool
    m_colourDescriptionPresentFlag   bool
    m_colourPrimaries                int
    m_transferCharacteristics        int
    m_matrixCoefficients             int
    m_chromaLocInfoPresentFlag       bool
    m_chromaSampleLocTypeTopField    int
    m_chromaSampleLocTypeBottomField int
    m_neutralChromaIndicationFlag    bool
    m_fieldSeqFlag                   bool

    m_defaultDisplayWindow               Window
    m_frameFieldInfoPresentFlag          bool
    m_hrdParametersPresentFlag           bool
    m_bitstreamRestrictionFlag           bool
    m_tilesFixedStructureFlag            bool
    m_motionVectorsOverPicBoundariesFlag bool
    m_restrictedRefPicListsFlag          bool
    m_minSpatialSegmentationIdc          int
    m_maxBytesPerPicDenom                int
    m_maxBitsPerMinCuDenom               int
    m_log2MaxMvLengthHorizontal          int
    m_log2MaxMvLengthVertical            int
    m_hrdParameters                      TComHRD
    m_timingInfo                         TimingInfo
}

//public:
func NewTComVUI() *TComVUI {
    return &TComVUI{
        m_aspectRatioInfoPresentFlag:         false,
        m_aspectRatioIdc:                     0,
        m_sarWidth:                           0,
        m_sarHeight:                          0,
        m_overscanInfoPresentFlag:            false,
        m_overscanAppropriateFlag:            false,
        m_videoSignalTypePresentFlag:         false,
        m_videoFormat:                        5,
        m_videoFullRangeFlag:                 false,
        m_colourDescriptionPresentFlag:       false,
        m_colourPrimaries:                    2,
        m_transferCharacteristics:            2,
        m_matrixCoefficients:                 2,
        m_chromaLocInfoPresentFlag:           false,
        m_chromaSampleLocTypeTopField:        0,
        m_chromaSampleLocTypeBottomField:     0,
        m_neutralChromaIndicationFlag:        false,
        m_fieldSeqFlag:                       false,
        m_frameFieldInfoPresentFlag:          false,
        m_hrdParametersPresentFlag:           false,
        m_bitstreamRestrictionFlag:           false,
        m_tilesFixedStructureFlag:            false,
        m_motionVectorsOverPicBoundariesFlag: true,
        m_restrictedRefPicListsFlag:          true,
        m_minSpatialSegmentationIdc:          0,
        m_maxBytesPerPicDenom:                2,
        m_maxBitsPerMinCuDenom:               1,
        m_log2MaxMvLengthHorizontal:          15,
        m_log2MaxMvLengthVertical:            15,
    }
}

func (this *TComVUI) GetAspectRatioInfoPresentFlag() bool {
    return this.m_aspectRatioInfoPresentFlag
}
func (this *TComVUI) SetAspectRatioInfoPresentFlag(i bool) {
    this.m_aspectRatioInfoPresentFlag = i
}

func (this *TComVUI) GetAspectRatioIdc() int {
    return this.m_aspectRatioIdc
}
func (this *TComVUI) SetAspectRatioIdc(i int) {
    this.m_aspectRatioIdc = i
}

func (this *TComVUI) GetSarWidth() int {
    return this.m_sarWidth
}
func (this *TComVUI) SetSarWidth(i int) {
    this.m_sarWidth = i
}

func (this *TComVUI) GetSarHeight() int {
    return this.m_sarHeight
}
func (this *TComVUI) SetSarHeight(i int) {
    this.m_sarHeight = i
}

func (this *TComVUI) GetOverscanInfoPresentFlag() bool {
    return this.m_overscanInfoPresentFlag
}
func (this *TComVUI) SetOverscanInfoPresentFlag(i bool) {
    this.m_overscanInfoPresentFlag = i
}

func (this *TComVUI) GetOverscanAppropriateFlag() bool {
    return this.m_overscanAppropriateFlag
}
func (this *TComVUI) SetOverscanAppropriateFlag(i bool) {
    this.m_overscanAppropriateFlag = i
}

func (this *TComVUI) GetVideoSignalTypePresentFlag() bool {
    return this.m_videoSignalTypePresentFlag
}
func (this *TComVUI) SetVideoSignalTypePresentFlag(i bool) {
    this.m_videoSignalTypePresentFlag = i
}

func (this *TComVUI) GetVideoFormat() int {
    return this.m_videoFormat
}
func (this *TComVUI) SetVideoFormat(i int) {
    this.m_videoFormat = i
}

func (this *TComVUI) GetVideoFullRangeFlag() bool {
    return this.m_videoFullRangeFlag
}
func (this *TComVUI) SetVideoFullRangeFlag(i bool) {
    this.m_videoFullRangeFlag = i
}

func (this *TComVUI) GetColourDescriptionPresentFlag() bool {
    return this.m_colourDescriptionPresentFlag
}
func (this *TComVUI) SetColourDescriptionPresentFlag(i bool) {
    this.m_colourDescriptionPresentFlag = i
}

func (this *TComVUI) GetColourPrimaries() int {
    return this.m_colourPrimaries
}
func (this *TComVUI) SetColourPrimaries(i int) {
    this.m_colourPrimaries = i
}

func (this *TComVUI) GetTransferCharacteristics() int {
    return this.m_transferCharacteristics
}
func (this *TComVUI) SetTransferCharacteristics(i int) {
    this.m_transferCharacteristics = i
}

func (this *TComVUI) GetMatrixCoefficients() int {
    return this.m_matrixCoefficients
}
func (this *TComVUI) SetMatrixCoefficients(i int) {
    this.m_matrixCoefficients = i
}

func (this *TComVUI) GetChromaLocInfoPresentFlag() bool {
    return this.m_chromaLocInfoPresentFlag
}
func (this *TComVUI) SetChromaLocInfoPresentFlag(i bool) {
    this.m_chromaLocInfoPresentFlag = i
}

func (this *TComVUI) GetChromaSampleLocTypeTopField() int {
    return this.m_chromaSampleLocTypeTopField
}
func (this *TComVUI) SetChromaSampleLocTypeTopField(i int) {
    this.m_chromaSampleLocTypeTopField = i
}

func (this *TComVUI) GetChromaSampleLocTypeBottomField() int {
    return this.m_chromaSampleLocTypeBottomField
}
func (this *TComVUI) SetChromaSampleLocTypeBottomField(i int) {
    this.m_chromaSampleLocTypeBottomField = i
}

func (this *TComVUI) GetNeutralChromaIndicationFlag() bool {
    return this.m_neutralChromaIndicationFlag
}
func (this *TComVUI) SetNeutralChromaIndicationFlag(i bool) {
    this.m_neutralChromaIndicationFlag = i
}

func (this *TComVUI) GetFieldSeqFlag() bool {
    return this.m_fieldSeqFlag
}
func (this *TComVUI) SetFieldSeqFlag(i bool) {
    this.m_fieldSeqFlag = i
}

func (this *TComVUI) GetFrameFieldInfoPresentFlag() bool  { return this.m_frameFieldInfoPresentFlag }
func (this *TComVUI) SetFrameFieldInfoPresentFlag(i bool) { this.m_frameFieldInfoPresentFlag = i }

func (this *TComVUI) GetDefaultDisplayWindow() *Window { return &this.m_defaultDisplayWindow }
func (this *TComVUI) SetDefaultDisplayWindow(defaultDisplayWindow *Window) {
    this.m_defaultDisplayWindow = *defaultDisplayWindow
}

func (this *TComVUI) GetHrdParametersPresentFlag() bool {
    return this.m_hrdParametersPresentFlag
}
func (this *TComVUI) SetHrdParametersPresentFlag(i bool) {
    this.m_hrdParametersPresentFlag = i
}

func (this *TComVUI) GetBitstreamRestrictionFlag() bool {
    return this.m_bitstreamRestrictionFlag
}
func (this *TComVUI) SetBitstreamRestrictionFlag(i bool) {
    this.m_bitstreamRestrictionFlag = i
}

func (this *TComVUI) GetTilesFixedStructureFlag() bool {
    return this.m_tilesFixedStructureFlag
}
func (this *TComVUI) SetTilesFixedStructureFlag(i bool) {
    this.m_tilesFixedStructureFlag = i
}

func (this *TComVUI) GetMotionVectorsOverPicBoundariesFlag() bool {
    return this.m_motionVectorsOverPicBoundariesFlag
}
func (this *TComVUI) SetMotionVectorsOverPicBoundariesFlag(i bool) {
    this.m_motionVectorsOverPicBoundariesFlag = i
}

func (this *TComVUI) GetRestrictedRefPicListsFlag() bool {
    return this.m_restrictedRefPicListsFlag
}
func (this *TComVUI) SetRestrictedRefPicListsFlag(b bool) {
    this.m_restrictedRefPicListsFlag = b
}

func (this *TComVUI) GetMinSpatialSegmentationIdc() int {
    return this.m_minSpatialSegmentationIdc
}
func (this *TComVUI) SetMinSpatialSegmentationIdc(i int) {
    this.m_minSpatialSegmentationIdc = i
}

func (this *TComVUI) GetMaxBytesPerPicDenom() int {
    return this.m_maxBytesPerPicDenom
}
func (this *TComVUI) SetMaxBytesPerPicDenom(i int) {
    this.m_maxBytesPerPicDenom = i
}

func (this *TComVUI) GetMaxBitsPerMinCuDenom() int {
    return this.m_maxBitsPerMinCuDenom
}
func (this *TComVUI) SetMaxBitsPerMinCuDenom(i int) {
    this.m_maxBitsPerMinCuDenom = i
}

func (this *TComVUI) GetLog2MaxMvLengthHorizontal() int {
    return this.m_log2MaxMvLengthHorizontal
}
func (this *TComVUI) SetLog2MaxMvLengthHorizontal(i int) {
    this.m_log2MaxMvLengthHorizontal = i
}

func (this *TComVUI) GetLog2MaxMvLengthVertical() int {
    return this.m_log2MaxMvLengthVertical
}
func (this *TComVUI) SetLog2MaxMvLengthVertical(i int) {
    this.m_log2MaxMvLengthVertical = i
}
func (this *TComVUI) GetHrdParameters() *TComHRD { return &this.m_hrdParameters }

func (this *TComVUI) GetTimingInfo() *TimingInfo { return &this.m_timingInfo }

/// SPS class
type TComSPS struct {
    //private:
    m_SPSId           int
    m_VPSId           int
    m_chromaFormatIdc int

    m_uiMaxTLayers uint // maximum number of temporal layers

    // Structure
    m_picWidthInLumaSamples  uint
    m_picHeightInLumaSamples uint

    m_conformanceWindow *Window

    m_uiMaxCUWidth         uint
    m_uiMaxCUHeight        uint
    m_uiMaxCUDepth         uint
    m_uiAddCUDepth		   uint
    m_uiMinTrDepth         uint
    m_uiMaxTrDepth         uint
    m_RPSList              *TComRPSList
    m_bLongTermRefsPresent bool
    m_TMVPFlagsPresent     bool
    m_numReorderPics       [MAX_TLAYER]int

    // Tool list
    m_uiQuadtreeTULog2MaxSize   uint
    m_uiQuadtreeTULog2MinSize   uint
    m_uiQuadtreeTUMaxDepthInter uint
    m_uiQuadtreeTUMaxDepthIntra uint
    m_usePCM                    bool
    m_pcmLog2MaxSize            uint
    m_uiPCMLog2MinSize          uint
    m_useAMP                    bool

    m_bUseLComb bool

    // Parameter
    m_bitDepthY   int
    m_bitDepthC   int
    m_qpBDOffsetY int
    m_qpBDOffsetC int

    m_useLossless bool

    m_uiPCMBitDepthLuma     uint
    m_uiPCMBitDepthChroma   uint
    m_bPCMFilterDisableFlag bool

    m_uiBitsForPOC           uint
    m_numLongTermRefPicSPS   uint
    m_ltRefPicPocLsbSps      [33]uint
    m_usedByCurrPicLtSPSFlag [33]bool
    // Max physical transform size
    m_uiMaxTrSize uint

    m_iAMPAcc [MAX_CU_DEPTH]int
    m_bUseSAO bool

    m_bTemporalIdNestingFlag bool // temporal_id_nesting_flag

    m_scalingListEnabledFlag bool
    m_scalingListPresentFlag bool
    m_scalingList            *TComScalingList //!< ScalingList class pointer
    m_uiMaxDecPicBuffering   [MAX_TLAYER]uint
    m_uiMaxLatencyIncrease   [MAX_TLAYER]uint

    m_useDF bool

    m_useStrongIntraSmoothing bool

    m_vuiParametersPresentFlag bool
    m_vuiParameters            TComVUI

    m_winUnitX [MAX_CHROMA_FORMAT_IDC + 1]int
    m_winUnitY [MAX_CHROMA_FORMAT_IDC + 1]int
    m_pcPTL    TComPTL
}

//public:
func NewTComSPS() *TComSPS {
    sps := &TComSPS{m_SPSId: 0,
        m_VPSId:           0,
        m_chromaFormatIdc: CHROMA_420,
        m_uiMaxTLayers:    1,
        // Structure
        m_picWidthInLumaSamples:     352,
        m_picHeightInLumaSamples:    288,
        m_uiMaxCUWidth:              32,
        m_uiMaxCUHeight:             32,
        m_uiMaxCUDepth:              3,
        m_uiMinTrDepth:              0,
        m_uiMaxTrDepth:              1,
        m_RPSList:					 NewTComRPSList(),
        m_bLongTermRefsPresent:      false,
        m_uiQuadtreeTULog2MaxSize:   0,
        m_uiQuadtreeTULog2MinSize:   0,
        m_uiQuadtreeTUMaxDepthInter: 0,
        m_uiQuadtreeTUMaxDepthIntra: 0,
        m_usePCM:                    false,
        m_pcmLog2MaxSize:            5,
        m_uiPCMLog2MinSize:          7,
        m_bUseLComb:                 false,
        m_bitDepthY:                 8,
        m_bitDepthC:                 8,
        m_qpBDOffsetY:               0,
        m_qpBDOffsetC:               0,
        m_useLossless:               false,
        m_uiPCMBitDepthLuma:         8,
        m_uiPCMBitDepthChroma:       8,
        m_bPCMFilterDisableFlag:     false,
        m_uiBitsForPOC:              8,
        m_numLongTermRefPicSPS:      0,
        m_uiMaxTrSize:               32,
        m_bUseSAO:                   false,
        m_bTemporalIdNestingFlag:    false,
        m_scalingListEnabledFlag:    false,
        m_useStrongIntraSmoothing:   false,
        m_vuiParametersPresentFlag:  false}

    for i := 0; i < MAX_TLAYER; i++ {
        sps.m_uiMaxLatencyIncrease[i] = 0
        sps.m_uiMaxDecPicBuffering[i] = 0
        sps.m_numReorderPics[i] = 0
    }
    sps.m_scalingList = NewTComScalingList()
    for i := 0; i < 33; i++ {
        sps.m_ltRefPicPocLsbSps[i] = 0
        sps.m_usedByCurrPicLtSPSFlag[i] = false
    }

    sps.m_conformanceWindow = NewWindow()
    sps.m_winUnitX[0] = 1
    sps.m_winUnitX[1] = 2
    sps.m_winUnitX[2] = 2
    sps.m_winUnitX[3] = 1
    sps.m_winUnitY[0] = 1
    sps.m_winUnitY[1] = 2
    sps.m_winUnitY[2] = 1
    sps.m_winUnitY[3] = 1

    return sps
}

func (this *TComSPS) GetVPSId() int {
    return this.m_VPSId
}
func (this *TComSPS) SetVPSId(i int) {
    this.m_VPSId = i
}
func (this *TComSPS) GetSPSId() int {
    return this.m_SPSId
}
func (this *TComSPS) SetSPSId(i int) {
    this.m_SPSId = i
}
func (this *TComSPS) GetChromaFormatIdc() int {
    return this.m_chromaFormatIdc
}
func (this *TComSPS) SetChromaFormatIdc(i int) {
    this.m_chromaFormatIdc = i
}

func (this *TComSPS) GetWinUnitX(chromaFormatIdc int) int {
    //assert (chromaFormatIdc > 0 && chromaFormatIdc <= MAX_CHROMA_FORMAT_IDC);
    return this.m_winUnitX[chromaFormatIdc]
}
func (this *TComSPS) GetWinUnitY(chromaFormatIdc int) int {
    //assert (chromaFormatIdc > 0 && chromaFormatIdc <= MAX_CHROMA_FORMAT_IDC);
    return this.m_winUnitY[chromaFormatIdc]
}

// structure
func (this *TComSPS) SetPicWidthInLumaSamples(u uint) {
    this.m_picWidthInLumaSamples = u
}
func (this *TComSPS) GetPicWidthInLumaSamples() uint {
    return this.m_picWidthInLumaSamples
}
func (this *TComSPS) SetPicHeightInLumaSamples(u uint) {
    this.m_picHeightInLumaSamples = u
}
func (this *TComSPS) GetPicHeightInLumaSamples() uint {
    return this.m_picHeightInLumaSamples
}

func (this *TComSPS) GetConformanceWindow() *Window {
    return this.m_conformanceWindow
}
func (this *TComSPS) SetConformanceWindow(conformanceWindow *Window) {
    this.m_conformanceWindow = conformanceWindow
}

func (this *TComSPS) GetNumLongTermRefPicSPS() uint {
    return this.m_numLongTermRefPicSPS
}
func (this *TComSPS) SetNumLongTermRefPicSPS(val uint) {
    this.m_numLongTermRefPicSPS = val
}

func (this *TComSPS) GetLtRefPicPocLsbSps(index uint) uint {
    return this.m_ltRefPicPocLsbSps[index]
}
func (this *TComSPS) SetLtRefPicPocLsbSps(index, val uint) {
    this.m_ltRefPicPocLsbSps[index] = val
}

func (this *TComSPS) GetUsedByCurrPicLtSPSFlag(i int) bool {
    return this.m_usedByCurrPicLtSPSFlag[i]
}
func (this *TComSPS) SetUsedByCurrPicLtSPSFlag(i int, x bool) {
    this.m_usedByCurrPicLtSPSFlag[i] = x
}
func (this *TComSPS) SetMaxCUWidth(u uint) {
    this.m_uiMaxCUWidth = u
}
func (this *TComSPS) GetMaxCUWidth() uint {
    return this.m_uiMaxCUWidth
}
func (this *TComSPS) SetMaxCUHeight(u uint) {
    this.m_uiMaxCUHeight = u
}
func (this *TComSPS) GetMaxCUHeight() uint {
    return this.m_uiMaxCUHeight
}
func (this *TComSPS) SetMaxCUDepth(u uint) {
    this.m_uiMaxCUDepth = u
}
func (this *TComSPS) GetMaxCUDepth() uint {
    return this.m_uiMaxCUDepth
}
func (this *TComSPS) SetAddCUDepth  ( u uint) { 
	this.m_uiAddCUDepth = u;      
}
func (this *TComSPS) GetAddCUDepth  ()  uint { 
	return  this.m_uiAddCUDepth;  
}
func (this *TComSPS) SetUsePCM(b bool) {
    this.m_usePCM = b
}
func (this *TComSPS) GetUsePCM() bool {
    return this.m_usePCM
}
func (this *TComSPS) SetPCMLog2MaxSize(u uint) {
    this.m_pcmLog2MaxSize = u
}
func (this *TComSPS) GetPCMLog2MaxSize() uint {
    return this.m_pcmLog2MaxSize
}
func (this *TComSPS) SetPCMLog2MinSize(u uint) {
    this.m_uiPCMLog2MinSize = u
}
func (this *TComSPS) GetPCMLog2MinSize() uint {
    return this.m_uiPCMLog2MinSize
}
func (this *TComSPS) SetBitsForPOC(u uint) {
    this.m_uiBitsForPOC = u
}
func (this *TComSPS) GetBitsForPOC() uint {
    return this.m_uiBitsForPOC
}
func (this *TComSPS) GetUseAMP() bool {
    return this.m_useAMP
}
func (this *TComSPS) SetUseAMP(b bool) {
    this.m_useAMP = b
}
func (this *TComSPS) SetMinTrDepth(u uint) {
    this.m_uiMinTrDepth = u
}
func (this *TComSPS) GetMinTrDepth() uint {
    return this.m_uiMinTrDepth
}
func (this *TComSPS) SetMaxTrDepth(u uint) {
    this.m_uiMaxTrDepth = u
}
func (this *TComSPS) GetMaxTrDepth() uint {
    return this.m_uiMaxTrDepth
}
func (this *TComSPS) SetQuadtreeTULog2MaxSize(u uint) {
    this.m_uiQuadtreeTULog2MaxSize = u
}
func (this *TComSPS) GetQuadtreeTULog2MaxSize() uint {
    return this.m_uiQuadtreeTULog2MaxSize
}
func (this *TComSPS) SetQuadtreeTULog2MinSize(u uint) {
    this.m_uiQuadtreeTULog2MinSize = u
}
func (this *TComSPS) GetQuadtreeTULog2MinSize() uint {
    return this.m_uiQuadtreeTULog2MinSize
}
func (this *TComSPS) SetQuadtreeTUMaxDepthInter(u uint) {
    this.m_uiQuadtreeTUMaxDepthInter = u
}
func (this *TComSPS) SetQuadtreeTUMaxDepthIntra(u uint) {
    this.m_uiQuadtreeTUMaxDepthIntra = u
}
func (this *TComSPS) GetQuadtreeTUMaxDepthInter() uint {
    return this.m_uiQuadtreeTUMaxDepthInter
}
func (this *TComSPS) GetQuadtreeTUMaxDepthIntra() uint {
    return this.m_uiQuadtreeTUMaxDepthIntra
}
func (this *TComSPS) SetNumReorderPics(i int, tlayer uint) {
    this.m_numReorderPics[tlayer] = i
}
func (this *TComSPS) GetNumReorderPics(tlayer uint) int {
    return this.m_numReorderPics[tlayer]
}
func (this *TComSPS) CreateRPSList(numRPS int) {
    this.m_RPSList.Destroy()
    this.m_RPSList.Create(numRPS)
}
func (this *TComSPS) GetRPSList() *TComRPSList {
    return this.m_RPSList
}
func (this *TComSPS) GetLongTermRefsPresent() bool {
    return this.m_bLongTermRefsPresent
}
func (this *TComSPS) SetLongTermRefsPresent(b bool) {
    this.m_bLongTermRefsPresent = b
}
func (this *TComSPS) GetTMVPFlagsPresent() bool {
    return this.m_TMVPFlagsPresent
}
func (this *TComSPS) SetTMVPFlagsPresent(b bool) {
    this.m_TMVPFlagsPresent = b
}

// physical transform
func (this *TComSPS) SetMaxTrSize(u uint) {
    this.m_uiMaxTrSize = u
}
func (this *TComSPS) GetMaxTrSize() uint {
    return this.m_uiMaxTrSize
}

// Tool list
func (this *TComSPS) SetUseLComb(b bool) {
    this.m_bUseLComb = b
}
func (this *TComSPS) GetUseLComb() bool {
    return this.m_bUseLComb
}

func (this *TComSPS) GetUseLossless() bool {
    return this.m_useLossless
}
func (this *TComSPS) SetUseLossless(b bool) {
    this.m_useLossless = b
}

// AMP accuracy
func (this *TComSPS) GetAMPAcc(uiDepth uint) int {
    return this.m_iAMPAcc[uiDepth]
}
func (this *TComSPS) SetAMPAcc(uiDepth uint, iAccu int) {
    //assert( uiDepth < g_uiMaxCUDepth);
    this.m_iAMPAcc[uiDepth] = iAccu
}

// Bit-depth
func (this *TComSPS) GetBitDepthY() int {
    return this.m_bitDepthY
}
func (this *TComSPS) SetBitDepthY(u int) {
    this.m_bitDepthY = u
}
func (this *TComSPS) GetBitDepthC() int {
    return this.m_bitDepthC
}
func (this *TComSPS) SetBitDepthC(u int) {
    this.m_bitDepthC = u
}
func (this *TComSPS) GetQpBDOffsetY() int {
    return this.m_qpBDOffsetY
}
func (this *TComSPS) SetQpBDOffsetY(value int) {
    this.m_qpBDOffsetY = value
}
func (this *TComSPS) GetQpBDOffsetC() int {
    return this.m_qpBDOffsetC
}
func (this *TComSPS) SetQpBDOffsetC(value int) {
    this.m_qpBDOffsetC = value
}
func (this *TComSPS) SetUseSAO(bVal bool) {
    this.m_bUseSAO = bVal
}
func (this *TComSPS) GetUseSAO() bool {
    return this.m_bUseSAO
}

func (this *TComSPS) GetMaxTLayers() uint {
    return this.m_uiMaxTLayers
}
func (this *TComSPS) SetMaxTLayers(uiMaxTLayers uint) {
    //assert( uiMaxTLayers <= MAX_TLAYER );
    this.m_uiMaxTLayers = uiMaxTLayers
}

func (this *TComSPS) GetTemporalIdNestingFlag() bool {
    return this.m_bTemporalIdNestingFlag
}
func (this *TComSPS) SetTemporalIdNestingFlag(bValue bool) {
    this.m_bTemporalIdNestingFlag = bValue
}
func (this *TComSPS) GetPCMBitDepthLuma() uint {
    return this.m_uiPCMBitDepthLuma
}
func (this *TComSPS) SetPCMBitDepthLuma(u uint) {
    this.m_uiPCMBitDepthLuma = u
}
func (this *TComSPS) GetPCMBitDepthChroma() uint {
    return this.m_uiPCMBitDepthChroma
}
func (this *TComSPS) SetPCMBitDepthChroma(u uint) {
    this.m_uiPCMBitDepthChroma = u
}
func (this *TComSPS) SetPCMFilterDisableFlag(bValue bool) {
    this.m_bPCMFilterDisableFlag = bValue
}
func (this *TComSPS) GetPCMFilterDisableFlag() bool {
    return this.m_bPCMFilterDisableFlag
}

func (this *TComSPS) GetScalingListFlag() bool {
    return this.m_scalingListEnabledFlag
}
func (this *TComSPS) SetScalingListFlag(b bool) {
    this.m_scalingListEnabledFlag = b
}
func (this *TComSPS) GetScalingListPresentFlag() bool {
    return this.m_scalingListPresentFlag
}
func (this *TComSPS) SetScalingListPresentFlag(b bool) {
    this.m_scalingListPresentFlag = b
}
func (this *TComSPS) SetScalingList(scalingList *TComScalingList) {
    this.m_scalingList = scalingList
}
func (this *TComSPS) GetScalingList() *TComScalingList {
    return this.m_scalingList
}   //!< get ScalingList class pointer in SPS
func (this *TComSPS) GetMaxDecPicBuffering(tlayer uint) uint {
    return this.m_uiMaxDecPicBuffering[tlayer]
}
func (this *TComSPS) SetMaxDecPicBuffering(ui, tlayer uint) {
    this.m_uiMaxDecPicBuffering[tlayer] = ui
}
func (this *TComSPS) GetMaxLatencyIncrease(tlayer uint) uint {
    return this.m_uiMaxLatencyIncrease[tlayer]
}
func (this *TComSPS) SetMaxLatencyIncrease(ui, tlayer uint) {
    this.m_uiMaxLatencyIncrease[tlayer] = ui
}
func (this *TComSPS) SetUseStrongIntraSmoothing(bVal bool) {
    this.m_useStrongIntraSmoothing = bVal
}
func (this *TComSPS) GetUseStrongIntraSmoothing() bool {
    return this.m_useStrongIntraSmoothing
}

func (this *TComSPS) GetVuiParametersPresentFlag() bool {
    return this.m_vuiParametersPresentFlag
}
func (this *TComSPS) SetVuiParametersPresentFlag(b bool) {
    this.m_vuiParametersPresentFlag = b
}
func (this *TComSPS) GetVuiParameters() *TComVUI {
    return &this.m_vuiParameters
}
func (this *TComSPS) SetHrdParameters(frameRate, numDU, bitRate uint, randomAccess bool) {
    if !this.GetVuiParametersPresentFlag() {
        return
    }

    vui := this.GetVuiParameters()
    hrd := vui.GetHrdParameters()

    timingInfo := vui.GetTimingInfo()
    timingInfo.SetTimingInfoPresentFlag(true)

    switch frameRate {
    case 24:
        timingInfo.SetNumUnitsInTick(1125000)
        timingInfo.SetTimeScale(27000000)

    case 25:
        timingInfo.SetNumUnitsInTick(1080000)
        timingInfo.SetTimeScale(27000000)

    case 30:
        timingInfo.SetNumUnitsInTick(900900)
        timingInfo.SetTimeScale(27000000)

    case 50:
        timingInfo.SetNumUnitsInTick(540000)
        timingInfo.SetTimeScale(27000000)

    case 60:
        timingInfo.SetNumUnitsInTick(450450)
        timingInfo.SetTimeScale(27000000)

    default:
        timingInfo.SetNumUnitsInTick(1001)
        timingInfo.SetTimeScale(60000)

    }

    rateCnt := (bitRate > 0)
    hrd.SetNalHrdParametersPresentFlag(rateCnt)
    hrd.SetVclHrdParametersPresentFlag(rateCnt)

    hrd.SetSubPicCpbParamsPresentFlag((numDU > 1))

    if hrd.GetSubPicCpbParamsPresentFlag() {
        hrd.SetTickDivisorMinus2(100 - 2)       //
        hrd.SetDuCpbRemovalDelayLengthMinus1(7) // 8-bit precision ( plus 1 for last DU in AU )
        hrd.SetSubPicCpbParamsInPicTimingSEIFlag(true)
        hrd.SetDpbOutputDelayDuLengthMinus1(5 + 7) // With sub-clock tick factor of 100, at least 7 bits to have the same value as AU dpb delay
    } else {
        hrd.SetSubPicCpbParamsInPicTimingSEIFlag(false)
    }

    hrd.SetBitRateScale(4)   // in units of 2~( 6 + 4 ) = 1,024 bps
    hrd.SetCpbSizeScale(6)   // in units of 2~( 4 + 4 ) = 1,024 bit
    hrd.SetDuCpbSizeScale(6) // in units of 2~( 4 + 4 ) = 1,024 bit

    hrd.SetInitialCpbRemovalDelayLengthMinus1(15) // assuming 0.5 sec, log2( 90,000 * 0.5 ) = 16-bit
    if randomAccess {
        hrd.SetCpbRemovalDelayLengthMinus1(5) // 32 = 2^5 (plus 1)
        hrd.SetDpbOutputDelayLengthMinus1(5)  // 32 + 3 = 2^6
    } else {
        hrd.SetCpbRemovalDelayLengthMinus1(9) // max. 2^10
        hrd.SetDpbOutputDelayLengthMinus1(9)  // max. 2^10
    }

    /*
       Note: only the case of "vps_max_temporal_layers_minus1 = 0" is supported.
    */
    var i, j int
    var birateValue, cpbSizeValue uint
    var ducpbSizeValue uint
    duBitRateValue := uint(0)

    for i = 0; i < MAX_TLAYER; i++ {
        hrd.SetFixedPicRateFlag(i, true)
        hrd.SetPicDurationInTcMinus1(i, 0)
        hrd.SetLowDelayHrdFlag(i, false)
        hrd.SetCpbCntMinus1(i, 0)

        birateValue = bitRate
        cpbSizeValue = bitRate // 1 second

        ducpbSizeValue = bitRate / numDU
        duBitRateValue = bitRate

        for j = 0; j < int(hrd.GetCpbCntMinus1(i)+1); j++ {
            hrd.SetBitRateValueMinus1(i, j, 0, (birateValue - 1))
            hrd.SetCpbSizeValueMinus1(i, j, 0, (cpbSizeValue - 1))
            hrd.SetDuCpbSizeValueMinus1(i, j, 0, (ducpbSizeValue - 1))
            hrd.SetCbrFlag(i, j, 0, (j == 0))

            hrd.SetBitRateValueMinus1(i, j, 1, (birateValue - 1))
            hrd.SetCpbSizeValueMinus1(i, j, 1, (cpbSizeValue - 1))
            hrd.SetDuCpbSizeValueMinus1(i, j, 1, (ducpbSizeValue - 1))
            hrd.SetDuBitRateValueMinus1(i, j, 1, (duBitRateValue - 1))
            hrd.SetCbrFlag(i, j, 1, (j == 0))
        }
    }
}

func (this *TComSPS) GetPTL() *TComPTL {
    return &this.m_pcPTL
}

//};

/// Reference Picture Lists class
type TComRefPicListModification struct {
    //private:
    m_bRefPicListModificationFlagL0 bool
    m_bRefPicListModificationFlagL1 bool
    m_RefPicSetIdxL0                [32]uint
    m_RefPicSetIdxL1                [32]uint
}

//public:
func NewTComRefPicListModification() *TComRefPicListModification {
    return &TComRefPicListModification{} //init zero by go
}

func (this *TComRefPicListModification) Create() {
    //do nothing
}
func (this *TComRefPicListModification) Destroy() {
    //do nothing
}

func (this *TComRefPicListModification) GetRefPicListModificationFlagL0() bool {
    return this.m_bRefPicListModificationFlagL0
}
func (this *TComRefPicListModification) SetRefPicListModificationFlagL0(flag bool) {
    this.m_bRefPicListModificationFlagL0 = flag
}
func (this *TComRefPicListModification) GetRefPicListModificationFlagL1() bool {
    return this.m_bRefPicListModificationFlagL1
}
func (this *TComRefPicListModification) SetRefPicListModificationFlagL1(flag bool) {
    this.m_bRefPicListModificationFlagL1 = flag
}
func (this *TComRefPicListModification) SetRefPicSetIdxL0(idx, refPicSetIdx uint) {
    this.m_RefPicSetIdxL0[idx] = refPicSetIdx
}
func (this *TComRefPicListModification) GetRefPicSetIdxL0(idx uint) uint {
    return this.m_RefPicSetIdxL0[idx]
}
func (this *TComRefPicListModification) SetRefPicSetIdxL1(idx, refPicSetIdx uint) {
    this.m_RefPicSetIdxL1[idx] = refPicSetIdx
}
func (this *TComRefPicListModification) GetRefPicSetIdxL1(idx uint) uint {
    return this.m_RefPicSetIdxL1[idx]
}

/// PPS class
type TComPPS struct {
    //private:
    m_PPSId                 int // pic_parameter_set_id
    m_SPSId                 int // seq_parameter_set_id
    m_picInitQPMinus26      int
    m_useDQP                bool
    m_bConstrainedIntraPred bool // constrained_intra_pred_flag
    m_bSliceChromaQpFlag    bool // slicelevel_chroma_qp_flag

    // access channel
    m_pcSPS           *TComSPS
    m_uiMaxCuDQPDepth uint
    m_uiMinCuDQPSize  uint

    m_chromaCbQpOffset int
    m_chromaCrQpOffset int

    m_numRefIdxL0DefaultActive uint
    m_numRefIdxL1DefaultActive uint

    m_bUseWeightPred        bool // Use of Weighting Prediction (P_SLICE)
    m_useWeightedBiPred     bool // Use of Weighting Bi-Prediction (B_SLICE)
    m_OutputFlagPresentFlag bool // Indicates the presence of output_flag in slice header

    m_TransquantBypassEnableFlag        bool // Indicates presence of cu_transquant_bypass_flag in CUs.
    m_useTransformSkip                  bool
    m_dependentSliceSegmentsEnabledFlag bool //!< Indicates the presence of dependent slices
    m_tilesEnabledFlag                  bool //!< Indicates the presence of tiles
    m_entropyCodingSyncEnabledFlag      bool //!< Indicates the presence of wavefronts

    m_loopFilterAcrossTilesEnabledFlag bool
    m_uniformSpacingFlag               bool
    m_iNumColumnsMinus1                int
    m_puiColumnWidth                   []int
    m_iNumRowsMinus1                   int
    m_puiRowHeight                     []int

    m_iNumSubstreams int

    m_signHideFlag bool

    m_cabacInitPresentFlag bool
    m_encCABACTableIdx     uint // Used to transmit table selection across slices

    m_sliceHeaderExtensionPresentFlag     bool
    m_loopFilterAcrossSlicesEnabledFlag   bool
    m_deblockingFilterControlPresentFlag  bool
    m_deblockingFilterOverrideEnabledFlag bool
    m_picDisableDeblockingFilterFlag      bool
    m_deblockingFilterBetaOffsetDiv2      int //< beta offset for deblocking filter
    m_deblockingFilterTcOffsetDiv2        int //< tc offset for deblocking filter
    m_scalingListPresentFlag              bool
    m_scalingList                         *TComScalingList //!< ScalingList class pointer

    m_listsModificationPresentFlag bool
    m_log2ParallelMergeLevelMinus2 uint
    m_numExtraSliceHeaderBits      int
}

//public:
func NewTComPPS() *TComPPS {
    pps := &TComPPS{m_PPSId: 0,
        m_SPSId:                             0,
        m_picInitQPMinus26:                  0,
        m_useDQP:                            false,
        m_bConstrainedIntraPred:             false,
        m_bSliceChromaQpFlag:                false,
        m_pcSPS:                             nil,
        m_uiMaxCuDQPDepth:                   0,
        m_uiMinCuDQPSize:                    0,
        m_chromaCbQpOffset:                  0,
        m_chromaCrQpOffset:                  0,
        m_numRefIdxL0DefaultActive:          1,
        m_numRefIdxL1DefaultActive:          1,
        m_TransquantBypassEnableFlag:        false,
        m_useTransformSkip:                  false,
        m_dependentSliceSegmentsEnabledFlag: false,
        m_tilesEnabledFlag:                  false,
        m_entropyCodingSyncEnabledFlag:      false,
        m_loopFilterAcrossTilesEnabledFlag:  true,
        m_uniformSpacingFlag:                false,
        m_iNumColumnsMinus1:                 0,
        m_puiColumnWidth:                    nil,
        m_iNumRowsMinus1:                    0,
        m_puiRowHeight:                      nil,
        m_iNumSubstreams:                    1,
        m_signHideFlag:                      false,
        m_cabacInitPresentFlag:              false,
        m_encCABACTableIdx:                  I_SLICE,
        m_sliceHeaderExtensionPresentFlag:   false,
        m_loopFilterAcrossSlicesEnabledFlag: false,
        m_listsModificationPresentFlag:      false,
        m_numExtraSliceHeaderBits:           0}

    pps.m_scalingList = NewTComScalingList()

    return pps
}

func (this *TComPPS) GetPPSId() int {
    return this.m_PPSId
}
func (this *TComPPS) SetPPSId(i int) {
    this.m_PPSId = i
}
func (this *TComPPS) GetSPSId() int {
    return this.m_SPSId
}
func (this *TComPPS) SetSPSId(i int) {
    this.m_SPSId = i
}

func (this *TComPPS) GetPicInitQPMinus26() int {
    return this.m_picInitQPMinus26
}
func (this *TComPPS) SetPicInitQPMinus26(i int) {
    this.m_picInitQPMinus26 = i
}
func (this *TComPPS) GetUseDQP() bool {
    return this.m_useDQP
}
func (this *TComPPS) SetUseDQP(b bool) {
    this.m_useDQP = b
}
func (this *TComPPS) GetConstrainedIntraPred() bool {
    return this.m_bConstrainedIntraPred
}
func (this *TComPPS) SetConstrainedIntraPred(b bool) {
    this.m_bConstrainedIntraPred = b
}
func (this *TComPPS) GetSliceChromaQpFlag() bool {
    return this.m_bSliceChromaQpFlag
}
func (this *TComPPS) SetSliceChromaQpFlag(b bool) {
    this.m_bSliceChromaQpFlag = b
}

func (this *TComPPS) SetSPS(pcSPS *TComSPS) {
    this.m_pcSPS = pcSPS
}
func (this *TComPPS) GetSPS() *TComSPS {
    return this.m_pcSPS
}
func (this *TComPPS) SetMaxCuDQPDepth(u uint) {
    this.m_uiMaxCuDQPDepth = u
}
func (this *TComPPS) GetMaxCuDQPDepth() uint {
    return this.m_uiMaxCuDQPDepth
}
func (this *TComPPS) SetMinCuDQPSize(u uint) {
    this.m_uiMinCuDQPSize = u
}
func (this *TComPPS) GetMinCuDQPSize() uint {
    return this.m_uiMinCuDQPSize
}

func (this *TComPPS) SetChromaCbQpOffset(i int) {
    this.m_chromaCbQpOffset = i
}
func (this *TComPPS) GetChromaCbQpOffset() int {
    return this.m_chromaCbQpOffset
}
func (this *TComPPS) SetChromaCrQpOffset(i int) {
    this.m_chromaCrQpOffset = i
}
func (this *TComPPS) GetChromaCrQpOffset() int {
    return this.m_chromaCrQpOffset
}

func (this *TComPPS) SetNumRefIdxL0DefaultActive(ui uint) {
    this.m_numRefIdxL0DefaultActive = ui
}
func (this *TComPPS) GetNumRefIdxL0DefaultActive() uint {
    return this.m_numRefIdxL0DefaultActive
}
func (this *TComPPS) SetNumRefIdxL1DefaultActive(ui uint) {
    this.m_numRefIdxL1DefaultActive = ui
}
func (this *TComPPS) GetNumRefIdxL1DefaultActive() uint {
    return this.m_numRefIdxL1DefaultActive
}

func (this *TComPPS) GetUseWP() bool {
    return this.m_bUseWeightPred
}
func (this *TComPPS) GetWPBiPred() bool {
    return this.m_useWeightedBiPred
}
func (this *TComPPS) SetUseWP(b bool) {
    this.m_bUseWeightPred = b
}
func (this *TComPPS) SetWPBiPred(b bool) {
    this.m_useWeightedBiPred = b
}
func (this *TComPPS) SetOutputFlagPresentFlag(b bool) {
    this.m_OutputFlagPresentFlag = b
}
func (this *TComPPS) GetOutputFlagPresentFlag() bool {
    return this.m_OutputFlagPresentFlag
}
func (this *TComPPS) SetTransquantBypassEnableFlag(b bool) {
    this.m_TransquantBypassEnableFlag = b
}
func (this *TComPPS) GetTransquantBypassEnableFlag() bool {
    return this.m_TransquantBypassEnableFlag
}

func (this *TComPPS) GetUseTransformSkip() bool {
    return this.m_useTransformSkip
}
func (this *TComPPS) SetUseTransformSkip(b bool) {
    this.m_useTransformSkip = b
}

func (this *TComPPS) SetLoopFilterAcrossTilesEnabledFlag(b bool) {
    this.m_loopFilterAcrossTilesEnabledFlag = b
}
func (this *TComPPS) GetLoopFilterAcrossTilesEnabledFlag() bool {
    return this.m_loopFilterAcrossTilesEnabledFlag
}
func (this *TComPPS) GetDependentSliceSegmentsEnabledFlag() bool {
    return this.m_dependentSliceSegmentsEnabledFlag
}
func (this *TComPPS) SetDependentSliceSegmentsEnabledFlag(val bool) {
    this.m_dependentSliceSegmentsEnabledFlag = val
}
func (this *TComPPS) GetTilesEnabledFlag() bool {
    return this.m_tilesEnabledFlag
}
func (this *TComPPS) SetTilesEnabledFlag(val bool) {
    this.m_tilesEnabledFlag = val
}
func (this *TComPPS) GetEntropyCodingSyncEnabledFlag() bool {
    return this.m_entropyCodingSyncEnabledFlag
}
func (this *TComPPS) SetEntropyCodingSyncEnabledFlag(val bool) {
    this.m_entropyCodingSyncEnabledFlag = val
}
func (this *TComPPS) SetUniformSpacingFlag(b bool) {
    this.m_uniformSpacingFlag = b
}
func (this *TComPPS) GetUniformSpacingFlag() bool {
    return this.m_uniformSpacingFlag
}
func (this *TComPPS) SetNumColumnsMinus1(i int) {
    this.m_iNumColumnsMinus1 = i
}
func (this *TComPPS) GetNumColumnsMinus1() int {
    return this.m_iNumColumnsMinus1
}
func (this *TComPPS) SetColumnWidth(columnWidth []int) {
    if this.m_uniformSpacingFlag == false && this.m_iNumColumnsMinus1 > 0 {
        this.m_puiColumnWidth = make([]int, this.m_iNumColumnsMinus1)
        for i := 0; i < this.m_iNumColumnsMinus1; i++ {
            this.m_puiColumnWidth[i] = columnWidth[i]
        }
    }
}
func (this *TComPPS) GetColumnWidth(columnIdx int) int {
    return this.m_puiColumnWidth[columnIdx]
}
func (this *TComPPS) SetNumRowsMinus1(i int) {
    this.m_iNumRowsMinus1 = i
}
func (this *TComPPS) GetNumRowsMinus1() int {
    return this.m_iNumRowsMinus1
}
func (this *TComPPS) SetRowHeight(rowHeight []int) {
    if this.m_uniformSpacingFlag == false && this.m_iNumRowsMinus1 > 0 {
        this.m_puiRowHeight = make([]int, this.m_iNumRowsMinus1)
        for i := 0; i < this.m_iNumRowsMinus1; i++ {
            this.m_puiRowHeight[i] = rowHeight[i]
        }
    }
}
func (this *TComPPS) GetRowHeight(rowIdx int) int {
    return this.m_puiRowHeight[rowIdx]
}
func (this *TComPPS) SetNumSubstreams(iNumSubstreams int) {
    this.m_iNumSubstreams = iNumSubstreams
}
func (this *TComPPS) GetNumSubstreams() int {
    return this.m_iNumSubstreams
}

func (this *TComPPS) SetSignHideFlag(signHideFlag bool) {
    this.m_signHideFlag = signHideFlag
}
func (this *TComPPS) GetSignHideFlag() bool {
    return this.m_signHideFlag
}

func (this *TComPPS) SetCabacInitPresentFlag(flag bool) {
    this.m_cabacInitPresentFlag = flag
}
func (this *TComPPS) SetEncCABACTableIdx(idx uint) {
    this.m_encCABACTableIdx = idx
}
func (this *TComPPS) GetCabacInitPresentFlag() bool {
    return this.m_cabacInitPresentFlag
}
func (this *TComPPS) GetEncCABACTableIdx() uint {
    return this.m_encCABACTableIdx
}
func (this *TComPPS) SetDeblockingFilterControlPresentFlag(val bool) {
    this.m_deblockingFilterControlPresentFlag = val
}
func (this *TComPPS) GetDeblockingFilterControlPresentFlag() bool {
    return this.m_deblockingFilterControlPresentFlag
}
func (this *TComPPS) SetDeblockingFilterOverrideEnabledFlag(val bool) {
    this.m_deblockingFilterOverrideEnabledFlag = val
}
func (this *TComPPS) GetDeblockingFilterOverrideEnabledFlag() bool {
    return this.m_deblockingFilterOverrideEnabledFlag
}
func (this *TComPPS) SetPicDisableDeblockingFilterFlag(val bool) {
    this.m_picDisableDeblockingFilterFlag = val
}   //!< Set offSet for deblocking filter disabled
func (this *TComPPS) GetPicDisableDeblockingFilterFlag() bool {
    return this.m_picDisableDeblockingFilterFlag
}   //!< Get offset for deblocking filter disabled
func (this *TComPPS) SetDeblockingFilterBetaOffsetDiv2(val int) {
    this.m_deblockingFilterBetaOffsetDiv2 = val
}   //!< set beta offset for deblocking filter
func (this *TComPPS) GetDeblockingFilterBetaOffsetDiv2() int {
    return this.m_deblockingFilterBetaOffsetDiv2
}   //!< Get beta offset for deblocking filter
func (this *TComPPS) SetDeblockingFilterTcOffsetDiv2(val int) {
    this.m_deblockingFilterTcOffsetDiv2 = val
}   //!< set tc offset for deblocking filter
func (this *TComPPS) GetDeblockingFilterTcOffsetDiv2() int {
    return this.m_deblockingFilterTcOffsetDiv2
}   //!< Get tc offset for deblocking filter
func (this *TComPPS) GetScalingListPresentFlag() bool {
    return this.m_scalingListPresentFlag
}
func (this *TComPPS) SetScalingListPresentFlag(b bool) {
    this.m_scalingListPresentFlag = b
}

func (this *TComPPS) SetScalingList(scalingList *TComScalingList) {
    this.m_scalingList = scalingList
}
func (this *TComPPS) GetScalingList() *TComScalingList {
    return this.m_scalingList
}   //!< Get ScalingList class pointer in PPS
func (this *TComPPS) GetListsModificationPresentFlag() bool {
    return this.m_listsModificationPresentFlag
}
func (this *TComPPS) SetListsModificationPresentFlag(b bool) {
    this.m_listsModificationPresentFlag = b
}
func (this *TComPPS) GetLog2ParallelMergeLevelMinus2() uint {
    return this.m_log2ParallelMergeLevelMinus2
}
func (this *TComPPS) SetLog2ParallelMergeLevelMinus2(mrgLevel uint) {
    this.m_log2ParallelMergeLevelMinus2 = mrgLevel
}
func (this *TComPPS) GetNumExtraSliceHeaderBits() int {
    return this.m_numExtraSliceHeaderBits
}
func (this *TComPPS) SetNumExtraSliceHeaderBits(i int) {
    this.m_numExtraSliceHeaderBits = i
}
func (this *TComPPS) SetLoopFilterAcrossSlicesEnabledFlag(bValue bool) {
    this.m_loopFilterAcrossSlicesEnabledFlag = bValue
}
func (this *TComPPS) GetLoopFilterAcrossSlicesEnabledFlag() bool {
    return this.m_loopFilterAcrossSlicesEnabledFlag
}
func (this *TComPPS) GetSliceHeaderExtensionPresentFlag() bool {
    return this.m_sliceHeaderExtensionPresentFlag
}
func (this *TComPPS) SetSliceHeaderExtensionPresentFlag(val bool) {
    this.m_sliceHeaderExtensionPresentFlag = val
}

type WpScalingParam struct {
    // Explicit weighted prediction parameters parsed in slice header,
    // or Implicit weighted prediction parameters (8 bits depth values).
    bPresentFlag      bool
    uiLog2WeightDenom uint
    iWeight           int
    iOffset           int

    // Weighted prediction scaling values built from above parameters (bitdepth scaled):
    W, O, Offset, Shift, Round int
}

func (this *WpScalingParam) SetPresentFlag(bPresentFlag bool) {
    this.bPresentFlag = bPresentFlag
}
func (this *WpScalingParam) SetLog2WeightDenom(uiLog2WeightDenom uint) {
    this.uiLog2WeightDenom = uiLog2WeightDenom
}
func (this *WpScalingParam) SetWeight(iWeight int) {
    this.iWeight = iWeight
}
func (this *WpScalingParam) SetOffset(iOffset int) {
    this.iOffset = iOffset
}
func (this *WpScalingParam) GetPresentFlag() bool {
    return this.bPresentFlag
}
func (this *WpScalingParam) GetLog2WeightDenom() uint {
    return this.uiLog2WeightDenom
}
func (this *WpScalingParam) GetWeight() int {
    return this.iWeight
}
func (this *WpScalingParam) GetOffset() int {
    return this.iOffset
}

type wpACDCParam struct {
    iAC int64
    iDC int64
}

var TComSlice_prevPOC   int

/// slice header class
type TComSlice struct {
    //private:
    //  Bitstream writing
    m_saoEnabledFlag       bool
    m_saoEnabledFlagChroma bool ///< SAO Cb&Cr enabled flag
    m_iPPSId               int  ///< picture parameter set ID
    m_PicOutputFlag        bool ///< pic_output_flag
    m_iPOC                 int
    m_iLastIDR             int
    //m_prevPOC              int
    m_pcRPS                *TComReferencePictureSet
    m_LocalRPS             TComReferencePictureSet
    m_iBDidx               int

    m_RefPicListModification    TComRefPicListModification
    m_eNalUnitType              NalUnitType ///< Nal unit type for the slice
    m_eSliceType                SliceType
    m_iSliceQp                  int
    m_dependentSliceSegmentFlag bool
    //#if ADAPTIVE_QP_SELECTION
    m_iSliceQpBase int
    //#endif
    m_deblockingFilterDisable        bool
    m_deblockingFilterOverrideFlag   bool //< offsets for deblocking filter inherit from PPS
    m_deblockingFilterBetaOffsetDiv2 int  //< beta offset for deblocking filter
    m_deblockingFilterTcOffsetDiv2   int  //< tc offset for deblocking filter

    m_aiNumRefIdx [3]int //  for multiple reference of current slice

    m_iRefIdxOfLC                   [2][MAX_NUM_REF_LC]int
    m_eListIdFromIdxOfLC            [MAX_NUM_REF_LC]int
    m_iRefIdxFromIdxOfLC            [MAX_NUM_REF_LC]int
    m_iRefIdxOfL1FromRefIdxOfL0     [MAX_NUM_REF_LC]int
    m_iRefIdxOfL0FromRefIdxOfL1     [MAX_NUM_REF_LC]int
    m_bRefPicListModificationFlagLC bool
    m_bRefPicListCombinationFlag    bool

    m_bCheckLDC bool

    //  Data
    m_iSliceQpDelta     int
    m_iSliceQpDeltaCb   int
    m_iSliceQpDeltaCr   int
    m_apcRefPicList     [2][MAX_NUM_REF + 1]*TComPic
    m_aiRefPOCList      [2][MAX_NUM_REF + 1]int
    m_bIsUsedAsLongTerm [2][MAX_NUM_REF + 1]bool
    m_iDepth            int

    // referenced slice?
    m_bRefenced bool

    // access channel
    m_pcVPS *TComVPS
    m_pcSPS *TComSPS
    m_pcPPS *TComPPS
    m_pcPic *TComPic
    //#if ADAPTIVE_QP_SELECTION
    m_pcTrQuant *TComTrQuant
    //#endif
    m_colFromL0Flag uint // collocated picture from List0 flag

    m_colRefIdx       uint
    m_maxNumMergeCand uint

    //#if SAO_CHROMA_LAMBDA
    m_dLambdaLuma   float64
    m_dLambdaChroma float64
    //#else
    //  Double      m_dLambda;
    //#endif

    m_abEqualRef [2][MAX_NUM_REF][MAX_NUM_REF]bool

    m_bNoBackPredFlag      bool
    m_uiTLayer             uint
    m_bTLayerSwitchingFlag bool

    m_sliceMode                  uint
    m_sliceArgument              uint
    m_sliceCurStartCUAddr        uint
    m_sliceCurEndCUAddr          uint
    m_sliceIdx                   uint
    m_sliceSegmentMode           uint
    m_sliceSegmentArgument       uint
    m_sliceSegmentCurStartCUAddr uint
    m_sliceSegmentCurEndCUAddr   uint
    m_nextSlice                  bool
    m_nextSliceSegment           bool
    m_sliceBits                  uint
    m_sliceSegmentBits           uint
    m_bFinalized                 bool

    m_weightPredTable [2][MAX_NUM_REF][3]WpScalingParam // [REF_PIC_LIST_0 or REF_PIC_LIST_1][refIdx][0:Y, 1:U, 2:V]
    m_weightACDCParam [3]wpACDCParam                    // [0:Y, 1:U, 2:V]

    m_tileByteLocation     map[int]uint //*list.List
    m_uiTileOffstForMultES uint

    m_puiSubstreamSizes []uint
    m_scalingList       *TComScalingList //!< pointer of quantization matrix
    m_cabacInitFlag     bool

    m_bLMvdL1Zero                   bool
    m_numEntryPointOffsets          int
    m_temporalLayerNonReferenceFlag bool
    m_LFCrossSliceBoundaryFlag      bool

    m_enableTMVPFlag bool
}

//public:
func NewTComSlice() *TComSlice {
    pSlice := &TComSlice{m_iPPSId: -1,
        m_iPOC:                      0,
        m_iLastIDR:                  0,
        m_eNalUnitType:              NAL_UNIT_CODED_SLICE_IDR,
        m_eSliceType:                I_SLICE,
        m_iSliceQp:                  0,
        m_dependentSliceSegmentFlag: false,
        //#if ADAPTIVE_QP_SELECTION
        m_iSliceQpBase: 0,
        //#endif
        m_deblockingFilterDisable:        false,
        m_deblockingFilterOverrideFlag:   false,
        m_deblockingFilterBetaOffsetDiv2: 0,
        m_deblockingFilterTcOffsetDiv2:   0,
        m_bRefPicListModificationFlagLC:  false,
        m_bRefPicListCombinationFlag:     false,
        m_bCheckLDC:                      false,
        m_iSliceQpDelta:                  0,
        m_iSliceQpDeltaCb:                0,
        m_iSliceQpDeltaCr:                0,
        m_iDepth:                         0,
        m_bRefenced:                      false,
        m_pcSPS:                          nil,
        m_pcPPS:                          nil,
        m_pcPic:                          nil,
        m_colFromL0Flag:                  1,
        m_colRefIdx:                      0,
        //#if SAO_CHROMA_LAMBDA
        m_dLambdaLuma:   0.0,
        m_dLambdaChroma: 0.0,
        //#else
        //, m_dLambda                       ( 0.0 )
        //#endif
        m_bNoBackPredFlag:            false,
        m_uiTLayer:                   0,
        m_bTLayerSwitchingFlag:       false,
        m_sliceMode:                  0,
        m_sliceArgument:              0,
        m_sliceCurStartCUAddr:        0,
        m_sliceCurEndCUAddr:          0,
        m_sliceIdx:                   0,
        m_sliceSegmentMode:           0,
        m_sliceSegmentArgument:       0,
        m_sliceSegmentCurStartCUAddr: 0,
        m_sliceSegmentCurEndCUAddr:   0,
        m_nextSlice:                  false,
        m_nextSliceSegment:           false,
        m_sliceBits:                  0,
        m_sliceSegmentBits:           0,
        m_bFinalized:                 false,
        m_uiTileOffstForMultES:       0,
        //m_puiSubstreamSizes             : NULL ,
        m_cabacInitFlag:                 false,
        m_bLMvdL1Zero:                   false,
        m_numEntryPointOffsets:          0,
        m_temporalLayerNonReferenceFlag: false,
        m_enableTMVPFlag:                true}

    pSlice.m_aiNumRefIdx[0] = 0
    pSlice.m_aiNumRefIdx[1] = 0
    pSlice.m_aiNumRefIdx[2] = 0

    pSlice.InitEqualRef()

    for iNumCount := 0; iNumCount < MAX_NUM_REF_LC; iNumCount++ {
        pSlice.m_iRefIdxOfLC[REF_PIC_LIST_0][iNumCount] = -1
        pSlice.m_iRefIdxOfLC[REF_PIC_LIST_1][iNumCount] = -1
        pSlice.m_eListIdFromIdxOfLC[iNumCount] = 0
        pSlice.m_iRefIdxFromIdxOfLC[iNumCount] = 0
        pSlice.m_iRefIdxOfL0FromRefIdxOfL1[iNumCount] = -1
        pSlice.m_iRefIdxOfL1FromRefIdxOfL0[iNumCount] = -1
    }
    for iNumCount := 0; iNumCount < MAX_NUM_REF; iNumCount++ {
        pSlice.m_apcRefPicList[0][iNumCount] = nil
        pSlice.m_apcRefPicList[1][iNumCount] = nil
        pSlice.m_aiRefPOCList[0][iNumCount] = 0
        pSlice.m_aiRefPOCList[1][iNumCount] = 0
    }
    pSlice.ResetWpScaling(pSlice.m_weightPredTable)
    pSlice.InitWpAcDcParam()
    pSlice.m_saoEnabledFlag = false

    return pSlice
}

func (this *TComSlice) InitSlice() {
    this.m_aiNumRefIdx[0] = 0
    this.m_aiNumRefIdx[1] = 0

    this.m_colFromL0Flag = 1

    this.m_colRefIdx = 0
    this.InitEqualRef()
    this.m_bNoBackPredFlag = false
    this.m_bRefPicListCombinationFlag = false
    this.m_bRefPicListModificationFlagLC = false
    this.m_bCheckLDC = false
    this.m_iSliceQpDeltaCb = 0
    this.m_iSliceQpDeltaCr = 0

    this.m_aiNumRefIdx[REF_PIC_LIST_C] = 0

    this.m_maxNumMergeCand = MRG_MAX_NUM_CANDS

    this.m_bFinalized = false

    this.m_tileByteLocation = make(map[int]uint) //list.New();
    this.m_cabacInitFlag = false
    this.m_numEntryPointOffsets = 0
    this.m_enableTMVPFlag = true
}

func (this *TComSlice) SetVPS(pcVPS *TComVPS) {
    this.m_pcVPS = pcVPS
}
func (this *TComSlice) GetVPS() *TComVPS {
    return this.m_pcVPS
}
func (this *TComSlice) SetSPS(pcSPS *TComSPS) {
    this.m_pcSPS = pcSPS
}
func (this *TComSlice) GetSPS() *TComSPS {
    return this.m_pcSPS
}

func (this *TComSlice) SetPPS(pcPPS *TComPPS) {
    //assert(pcPPS!=NULL);
    this.m_pcPPS = pcPPS
    this.m_iPPSId = pcPPS.GetPPSId()
}
func (this *TComSlice) GetPPS() *TComPPS {
    return this.m_pcPPS
}

//#if ADAPTIVE_QP_SELECTION
func (this *TComSlice) SetTrQuant(pcTrQuant *TComTrQuant) {
    this.m_pcTrQuant = pcTrQuant
}
func (this *TComSlice) GetTrQuant() *TComTrQuant {
    return this.m_pcTrQuant
}

//#endif

func (this *TComSlice) SetPPSId(PPSId int) {
    this.m_iPPSId = PPSId
}
func (this *TComSlice) GetPPSId() int {
    return this.m_iPPSId
}
func (this *TComSlice) SetPicOutputFlag(b bool) {
    this.m_PicOutputFlag = b
}
func (this *TComSlice) GetPicOutputFlag() bool {
    return this.m_PicOutputFlag
}
func (this *TComSlice) SetSaoEnabledFlag(s bool) {
    this.m_saoEnabledFlag = s
}
func (this *TComSlice) GetSaoEnabledFlag() bool {
    return this.m_saoEnabledFlag
}
func (this *TComSlice) SetSaoEnabledFlagChroma(s bool) {
    this.m_saoEnabledFlagChroma = s
}   //!< Set SAO Cb&Cr enabled flag
func (this *TComSlice) GetSaoEnabledFlagChroma() bool {
    return this.m_saoEnabledFlagChroma
}   //!< Get SAO Cb&Cr enabled flag
func (this *TComSlice) SetRPS(pcRPS *TComReferencePictureSet) {
    this.m_pcRPS = pcRPS
}
func (this *TComSlice) GetRPS() *TComReferencePictureSet {
    return this.m_pcRPS
}
func (this *TComSlice) GetLocalRPS() *TComReferencePictureSet {
    return &this.m_LocalRPS
}

func (this *TComSlice) SetRPSidx(iBDidx int) {
    this.m_iBDidx = iBDidx
}
func (this *TComSlice) GetRPSidx() int {
    return this.m_iBDidx
}
func (this *TComSlice) GetPrevPOC() int {
    return TComSlice_prevPOC;//this.m_prevPOC
}
func (this *TComSlice) GetRefPicListModification() *TComRefPicListModification {
    return &this.m_RefPicListModification
}
func (this *TComSlice) SetLastIDR(iIDRPOC int) {
    this.m_iLastIDR = iIDRPOC
}
func (this *TComSlice) GetLastIDR() int {
    return this.m_iLastIDR
}
func (this *TComSlice) GetSliceType() SliceType {
    return this.m_eSliceType
}
func (this *TComSlice) GetPOC() int {
    return this.m_iPOC
}
func (this *TComSlice) GetSliceQp() int {
    return this.m_iSliceQp
}
func (this *TComSlice) GetDependentSliceSegmentFlag() bool {
    return this.m_dependentSliceSegmentFlag
}
func (this *TComSlice) SetDependentSliceSegmentFlag(val bool) {
    this.m_dependentSliceSegmentFlag = val
}

//#if ADAPTIVE_QP_SELECTION
func (this *TComSlice) GetSliceQpBase() int {
    return this.m_iSliceQpBase
}

//#endif
func (this *TComSlice) GetSliceQpDelta() int {
    return this.m_iSliceQpDelta
}
func (this *TComSlice) GetSliceQpDeltaCb() int {
    return this.m_iSliceQpDeltaCb
}
func (this *TComSlice) GetSliceQpDeltaCr() int {
    return this.m_iSliceQpDeltaCr
}
func (this *TComSlice) GetDeblockingFilterDisable() bool {
    return this.m_deblockingFilterDisable
}
func (this *TComSlice) GetDeblockingFilterOverrideFlag() bool {
    return this.m_deblockingFilterOverrideFlag
}
func (this *TComSlice) GetDeblockingFilterBetaOffsetDiv2() int {
    return this.m_deblockingFilterBetaOffsetDiv2
}
func (this *TComSlice) GetDeblockingFilterTcOffsetDiv2() int {
    return this.m_deblockingFilterTcOffsetDiv2
}

func (this *TComSlice) GetNumRefIdx(e RefPicList) int {
    return this.m_aiNumRefIdx[e]
}
func (this *TComSlice) GetPic() *TComPic {
    return this.m_pcPic
}
func (this *TComSlice) GetRefPic(e RefPicList, iRefIdx int) *TComPic {
    return this.m_apcRefPicList[e][iRefIdx]
}
func (this *TComSlice) GetRefPOC(e RefPicList, iRefIdx int) int {
    return this.m_aiRefPOCList[e][iRefIdx]
}
func (this *TComSlice) GetDepth() int {
    return this.m_iDepth
}
func (this *TComSlice) GetColFromL0Flag() uint {
    return this.m_colFromL0Flag
}
func (this *TComSlice) GetColRefIdx() uint {
    return this.m_colRefIdx
}
func (this *TComSlice) CheckColRefIdx(curSliceIdx uint, pic *TComPic) {
    var i int
    curSlice := pic.GetSlice(curSliceIdx)
    currColRefPOC := curSlice.GetRefPOC(RefPicList(1-curSlice.GetColFromL0Flag()), int(curSlice.GetColRefIdx()))
    var preSlice *TComSlice
    var preColRefPOC int
    for i = int(curSliceIdx) - 1; i >= 0; i-- {
        preSlice = pic.GetSlice(uint(i))
        if preSlice.GetSliceType() != I_SLICE {
            preColRefPOC = preSlice.GetRefPOC(RefPicList(1-preSlice.GetColFromL0Flag()), int(preSlice.GetColRefIdx()))
            if currColRefPOC != preColRefPOC {
                fmt.Printf("Collocated_ref_idx shall always be the same for all slices of a coded picture!\n")
                return //exit(EXIT_FAILURE);
            } else {
                break
            }
        }
    }
}
func (this *TComSlice) GetIsUsedAsLongTerm(i, j int) bool { return this.m_bIsUsedAsLongTerm[i][j] }

func (this *TComSlice) GetCheckLDC() bool {
    return this.m_bCheckLDC
}
func (this *TComSlice) GetMvdL1ZeroFlag() bool {
    return this.m_bLMvdL1Zero
}
func (this *TComSlice) GetNumRpsCurrTempList() int {
    numRpsCurrTempList := 0

    if this.m_eSliceType == I_SLICE {
        return 0
    }
    for i := 0; i < this.m_pcRPS.GetNumberOfNegativePictures()+this.m_pcRPS.GetNumberOfPositivePictures()+this.m_pcRPS.GetNumberOfLongtermPictures(); i++ {
        if this.m_pcRPS.GetUsed(i) {
            numRpsCurrTempList++
        }
    }
    return numRpsCurrTempList
}
func (this *TComSlice) GetRefIdxOfLC(e RefPicList, iRefIdx int) int {
    return this.m_iRefIdxOfLC[e][iRefIdx]
}
func (this *TComSlice) GetListIdFromIdxOfLC(iRefIdx int) int {
    return this.m_eListIdFromIdxOfLC[iRefIdx]
}
func (this *TComSlice) GetRefIdxFromIdxOfLC(iRefIdx int) int {
    return this.m_iRefIdxFromIdxOfLC[iRefIdx]
}

func (this *TComSlice) GetRefIdxOfL0FromRefIdxOfL1(iRefIdx int) int {
    return this.m_iRefIdxOfL0FromRefIdxOfL1[iRefIdx]
}
func (this *TComSlice) GetRefIdxOfL1FromRefIdxOfL0(iRefIdx int) int {
    return this.m_iRefIdxOfL1FromRefIdxOfL0[iRefIdx]
}
func (this *TComSlice) GetRefPicListModificationFlagLC() bool {
    return this.m_bRefPicListModificationFlagLC
}
func (this *TComSlice) SetRefPicListModificationFlagLC(bflag bool) {
    this.m_bRefPicListModificationFlagLC = bflag
}
func (this *TComSlice) GetRefPicListCombinationFlag() bool {
    return this.m_bRefPicListCombinationFlag
}
func (this *TComSlice) SetRefPicListCombinationFlag(bflag bool) {
    this.m_bRefPicListCombinationFlag = bflag
}
func (this *TComSlice) SetReferenced(b bool) {
    this.m_bRefenced = b
}
func (this *TComSlice) IsReferenced() bool {
    return this.m_bRefenced
}
func (this *TComSlice) SetPOC(i int) {
    this.m_iPOC = i
    if this.GetTLayer() == 0 {
    	//this.m_prevPOC = i
    	TComSlice_prevPOC = i
    }
}
func (this *TComSlice) SetNalUnitType(e NalUnitType) {
    this.m_eNalUnitType = e
}
func (this *TComSlice) GetNalUnitType() NalUnitType {
    return this.m_eNalUnitType
}
func (this *TComSlice) GetRapPicFlag() bool {
    return this.GetNalUnitType() == NAL_UNIT_CODED_SLICE_IDR ||
        this.GetNalUnitType() == NAL_UNIT_CODED_SLICE_IDR_N_LP ||
        this.GetNalUnitType() == NAL_UNIT_CODED_SLICE_BLA_N_LP ||
        this.GetNalUnitType() == NAL_UNIT_CODED_SLICE_BLANT ||
        this.GetNalUnitType() == NAL_UNIT_CODED_SLICE_BLA ||
        this.GetNalUnitType() == NAL_UNIT_CODED_SLICE_CRA
}
func (this *TComSlice) GetIdrPicFlag() bool {
    return this.GetNalUnitType() == NAL_UNIT_CODED_SLICE_IDR || this.GetNalUnitType() == NAL_UNIT_CODED_SLICE_IDR_N_LP
}
func (this *TComSlice) CheckCRA(pReferencePictureSet *TComReferencePictureSet, pocCRA *int, prevRAPisBLA *bool) {
    for i := int(0); i < pReferencePictureSet.GetNumberOfNegativePictures()+pReferencePictureSet.GetNumberOfPositivePictures(); i++ {
        if uint(*pocCRA) < MAX_UINT && this.GetPOC() > *pocCRA {
            //assert(getPOC()+pReferencePictureSet.GetDeltaPOC(i) >= pocCRA);
        }
    }
    for i := int(pReferencePictureSet.GetNumberOfNegativePictures() + pReferencePictureSet.GetNumberOfPositivePictures()); i < pReferencePictureSet.GetNumberOfPictures(); i++ {
        if uint(*pocCRA) < MAX_UINT && this.GetPOC() > *pocCRA {
            //assert(pReferencePictureSet.GetPOC(i) >= pocCRA);
        }
    }
    if this.GetNalUnitType() == NAL_UNIT_CODED_SLICE_IDR || this.GetNalUnitType() == NAL_UNIT_CODED_SLICE_IDR_N_LP { // IDR picture found
        *pocCRA = this.GetPOC()
        *prevRAPisBLA = false
    } else if this.GetNalUnitType() == NAL_UNIT_CODED_SLICE_CRA { // CRA picture found
        *pocCRA = this.GetPOC()
        *prevRAPisBLA = false
    } else if this.GetNalUnitType() == NAL_UNIT_CODED_SLICE_BLA ||
        this.GetNalUnitType() == NAL_UNIT_CODED_SLICE_BLANT ||
        this.GetNalUnitType() == NAL_UNIT_CODED_SLICE_BLA_N_LP { // BLA picture found
        *pocCRA = this.GetPOC()
        *prevRAPisBLA = true
    }
}
func (this *TComSlice) DecodingRefreshMarking(pocCRA *int, bRefreshPending *bool, rcListPic *list.List) {
    var rpcPic *TComPic
    pocCurr := this.GetPOC()

    if this.GetNalUnitType() == NAL_UNIT_CODED_SLICE_BLA ||
        this.GetNalUnitType() == NAL_UNIT_CODED_SLICE_BLANT ||
        this.GetNalUnitType() == NAL_UNIT_CODED_SLICE_BLA_N_LP ||
        this.GetNalUnitType() == NAL_UNIT_CODED_SLICE_IDR ||
        this.GetNalUnitType() == NAL_UNIT_CODED_SLICE_IDR_N_LP { // IDR or BLA picture
        // mark all pictures as not used for reference
        iterPic := rcListPic.Front() // TComList<TComPic*>::iterator
        for iterPic != nil {
            rpcPic = iterPic.Value.(*TComPic)
            rpcPic.SetCurrSliceIdx(0)
            if int(rpcPic.GetPOC()) != pocCurr {
                rpcPic.GetSlice(0).SetReferenced(false)
            }
            iterPic = iterPic.Next()
        }
        if this.GetNalUnitType() == NAL_UNIT_CODED_SLICE_BLA ||
            this.GetNalUnitType() == NAL_UNIT_CODED_SLICE_BLANT ||
            this.GetNalUnitType() == NAL_UNIT_CODED_SLICE_BLA_N_LP {
            *pocCRA = pocCurr
        }
    } else { // CRA or No DR
        if *bRefreshPending == true && pocCurr > *pocCRA { // CRA reference marking pending
            iterPic := rcListPic.Front() // TComList<TComPic*>::iterator
            for iterPic != nil {
                rpcPic = iterPic.Value.(*TComPic)
                if int(rpcPic.GetPOC()) != pocCurr && int(rpcPic.GetPOC()) != *pocCRA {
                    rpcPic.GetSlice(0).SetReferenced(false)
                }
                iterPic = iterPic.Next()
            }
            *bRefreshPending = false
        }
        if this.GetNalUnitType() == NAL_UNIT_CODED_SLICE_CRA { // CRA picture found
            *bRefreshPending = true
            *pocCRA = pocCurr
        }
    }
}
func (this *TComSlice) SetSliceType(e SliceType) {
    this.m_eSliceType = e
}
func (this *TComSlice) SetSliceQp(i int) {
    this.m_iSliceQp = i
}

//#if ADAPTIVE_QP_SELECTION
func (this *TComSlice) SetSliceQpBase(i int) {
    this.m_iSliceQpBase = i
}

//#endif
func (this *TComSlice) SetSliceQpDelta(i int) {
    this.m_iSliceQpDelta = i
}
func (this *TComSlice) SetSliceQpDeltaCb(i int) {
    this.m_iSliceQpDeltaCb = i
}
func (this *TComSlice) SetSliceQpDeltaCr(i int) {
    this.m_iSliceQpDeltaCr = i
}
func (this *TComSlice) SetDeblockingFilterDisable(b bool) {
    this.m_deblockingFilterDisable = b
}
func (this *TComSlice) SetDeblockingFilterOverrideFlag(b bool) {
    this.m_deblockingFilterOverrideFlag = b
}
func (this *TComSlice) SetDeblockingFilterBetaOffsetDiv2(i int) {
    this.m_deblockingFilterBetaOffsetDiv2 = i
}
func (this *TComSlice) SetDeblockingFilterTcOffsetDiv2(i int) {
    this.m_deblockingFilterTcOffsetDiv2 = i
}

func (this *TComSlice) SetRefPic(p *TComPic, e RefPicList, iRefIdx int) {
    this.m_apcRefPicList[e][iRefIdx] = p
}
func (this *TComSlice) SetRefPOC(i int, e RefPicList, iRefIdx int) {
    this.m_aiRefPOCList[e][iRefIdx] = i
}
func (this *TComSlice) SetNumRefIdx(e RefPicList, i int) {
    this.m_aiNumRefIdx[e] = i
}
func (this *TComSlice) SetPic(p *TComPic) {
    this.m_pcPic = p
}
func (this *TComSlice) SetDepth(iDepth int) {
    this.m_iDepth = iDepth
}

func (this *TComSlice) SetRefPicList(rcListPic *list.List) {
    if this.m_eSliceType == I_SLICE {
        for j := 0; j < 2; j++ {
            for i := 0; i < MAX_NUM_REF+1; i++ {
                this.m_apcRefPicList[j][i] = nil
            }
        }
        for i := 0; i < 3; i++ {
            this.m_aiNumRefIdx[i] = 0
        }
        //::memset( m_apcRefPicList, 0, sizeof (m_apcRefPicList));
        //::memset( m_aiNumRefIdx,   0, sizeof ( m_aiNumRefIdx ));
        return
    }

    this.m_aiNumRefIdx[0] = this.GetNumRefIdx(REF_PIC_LIST_0)
    this.m_aiNumRefIdx[1] = this.GetNumRefIdx(REF_PIC_LIST_1)

    var pcRefPic *TComPic
    var RefPicSetStCurr0 [16]*TComPic
    var RefPicSetStCurr1 [16]*TComPic
    var RefPicSetLtCurr [16]*TComPic
    NumPocStCurr0 := 0
    NumPocStCurr1 := 0
    NumPocLtCurr := 0

    var i int
    for i = 0; i < this.m_pcRPS.GetNumberOfNegativePictures(); i++ {
        if this.m_pcRPS.GetUsed(i) {
            pcRefPic = this.xGetRefPic(rcListPic, this.GetPOC()+this.m_pcRPS.GetDeltaPOC(i))
            pcRefPic.SetIsLongTerm(false)
            pcRefPic.GetPicYuvRec().ExtendPicBorder()
            RefPicSetStCurr0[NumPocStCurr0] = pcRefPic
            NumPocStCurr0++
            pcRefPic.SetCheckLTMSBPresent(false)
        }
    }
    for ; i < this.m_pcRPS.GetNumberOfNegativePictures()+this.m_pcRPS.GetNumberOfPositivePictures(); i++ {
        if this.m_pcRPS.GetUsed(i) {
            pcRefPic = this.xGetRefPic(rcListPic, this.GetPOC()+this.m_pcRPS.GetDeltaPOC(i))
            pcRefPic.SetIsLongTerm(false)
            pcRefPic.GetPicYuvRec().ExtendPicBorder()
            RefPicSetStCurr1[NumPocStCurr1] = pcRefPic
            NumPocStCurr1++
            pcRefPic.SetCheckLTMSBPresent(false)
        }
    }
    for i = this.m_pcRPS.GetNumberOfNegativePictures() + this.m_pcRPS.GetNumberOfPositivePictures() + this.m_pcRPS.GetNumberOfLongtermPictures() - 1; i > this.m_pcRPS.GetNumberOfNegativePictures()+this.m_pcRPS.GetNumberOfPositivePictures()-1; i-- {
        if this.m_pcRPS.GetUsed(i) {
            pcRefPic = this.xGetLongTermRefPic(rcListPic, this.m_pcRPS.GetPOC(i), this.m_pcRPS.GetCheckLTMSBPresent(i))
            pcRefPic.SetIsLongTerm(true)
            pcRefPic.GetPicYuvRec().ExtendPicBorder()
            RefPicSetLtCurr[NumPocLtCurr] = pcRefPic
            NumPocLtCurr++
        }
        if pcRefPic == nil {
            pcRefPic = this.xGetLongTermRefPic(rcListPic, this.m_pcRPS.GetPOC(i), this.m_pcRPS.GetCheckLTMSBPresent(i))
        }
        pcRefPic.SetCheckLTMSBPresent(this.m_pcRPS.GetCheckLTMSBPresent(i))
    }

    // ref_pic_list_init
    var rpsCurrList0 [MAX_NUM_REF + 1]*TComPic
    var rpsCurrList1 [MAX_NUM_REF + 1]*TComPic
    numPocTotalCurr := NumPocStCurr0 + NumPocStCurr1 + NumPocLtCurr

    //{
    cIdx := 0
    for i = 0; i < NumPocStCurr0; i++ {
        rpsCurrList0[cIdx] = RefPicSetStCurr0[i]
        cIdx++
    }
    for i = 0; i < NumPocStCurr1; i++ {
        rpsCurrList0[cIdx] = RefPicSetStCurr1[i]
        cIdx++
    }
    for i = 0; i < NumPocLtCurr; i++ {
        rpsCurrList0[cIdx] = RefPicSetLtCurr[i]
        cIdx++
    }

    if this.m_eSliceType == B_SLICE {
        cIdx := 0
        for i = 0; i < NumPocStCurr1; i++ {
            rpsCurrList1[cIdx] = RefPicSetStCurr1[i]
            cIdx++
        }
        for i = 0; i < NumPocStCurr0; i++ {
            rpsCurrList1[cIdx] = RefPicSetStCurr0[i]
            cIdx++
        }
        for i = 0; i < NumPocLtCurr; i++ {
            rpsCurrList1[cIdx] = RefPicSetLtCurr[i]
            cIdx++
        }
    }

    for j := 0; j < 2; j++ {
        for i := 0; i < MAX_NUM_REF+1; i++ {
            this.m_bIsUsedAsLongTerm[j][i] = false //, [2][MAX_NUM_REF+1]boolsizeof(m_bIsUsedAsLongTerm));
        }
    }

    for rIdx := 0; rIdx <= (this.m_aiNumRefIdx[0] - 1); rIdx++ {
        if this.m_RefPicListModification.GetRefPicListModificationFlagL0() {
            this.m_apcRefPicList[0][rIdx] = rpsCurrList0[this.m_RefPicListModification.GetRefPicSetIdxL0(uint(rIdx))]
            this.m_bIsUsedAsLongTerm[0][rIdx] = (this.m_RefPicListModification.GetRefPicSetIdxL0(uint(rIdx)) >= uint(NumPocStCurr0+NumPocStCurr1))
        } else {
            this.m_apcRefPicList[0][rIdx] = rpsCurrList0[rIdx%numPocTotalCurr]
            this.m_bIsUsedAsLongTerm[0][rIdx] = ((rIdx % numPocTotalCurr) >= (NumPocStCurr0 + NumPocStCurr1))
        }
    }
    if this.m_eSliceType == P_SLICE {
        this.m_aiNumRefIdx[1] = 0
        //::memset( m_apcRefPicList[1], 0, sizeof(m_apcRefPicList[1]));
    } else {
        for rIdx := 0; rIdx <= (this.m_aiNumRefIdx[1] - 1); rIdx++ {
            if this.m_RefPicListModification.GetRefPicListModificationFlagL1() {
                this.m_apcRefPicList[1][rIdx] = rpsCurrList1[this.m_RefPicListModification.GetRefPicSetIdxL1(uint(rIdx))]
                this.m_bIsUsedAsLongTerm[1][rIdx] = (this.m_RefPicListModification.GetRefPicSetIdxL1(uint(rIdx)) >= uint(NumPocStCurr0+NumPocStCurr1))
            } else {
                this.m_apcRefPicList[1][rIdx] = rpsCurrList1[rIdx%numPocTotalCurr]
                this.m_bIsUsedAsLongTerm[1][rIdx] = ((rIdx % numPocTotalCurr) >= (NumPocStCurr0 + NumPocStCurr1))
            }
        }
    }
}
func (this *TComSlice) SetRefPOCList() {
    for iDir := 0; iDir < 2; iDir++ {
        for iNumRefIdx := 0; iNumRefIdx < this.m_aiNumRefIdx[iDir]; iNumRefIdx++ {
            this.m_aiRefPOCList[iDir][iNumRefIdx] = int(this.m_apcRefPicList[iDir][iNumRefIdx].GetPOC())
        }
    }
}
func (this *TComSlice) SetColFromL0Flag(colFromL0 uint) {
    this.m_colFromL0Flag = colFromL0
}
func (this *TComSlice) SetColRefIdx(refIdx uint) {
    this.m_colRefIdx = refIdx
}
func (this *TComSlice) SetCheckLDC(b bool) {
    this.m_bCheckLDC = b
}
func (this *TComSlice) SetMvdL1ZeroFlag(b bool) {
    this.m_bLMvdL1Zero = b
}

func (this *TComSlice) IsIntra() bool {
    return this.m_eSliceType == I_SLICE
}
func (this *TComSlice) IsInterB() bool {
    return this.m_eSliceType == B_SLICE
}
func (this *TComSlice) IsInterP() bool {
    return this.m_eSliceType == P_SLICE
}

//#if SAO_CHROMA_LAMBDA
func (this *TComSlice) SetLambda(d, e float64) {
    this.m_dLambdaLuma = d
    this.m_dLambdaChroma = e
}
func (this *TComSlice) GetLambdaLuma() float64 {
    return this.m_dLambdaLuma
}
func (this *TComSlice) GetLambdaChroma() float64 {
    return this.m_dLambdaChroma
}

//#else
//  Void      SetLambda( Double d ) { this.m_dLambda = d; }
//  Double    GetLambda() { return this.m_dLambda;        }
//#endif

func (this *TComSlice) InitEqualRef() {
    for iDir := int(0); iDir < 2; iDir++ {
        for iRefIdx1 := int(0); iRefIdx1 < MAX_NUM_REF; iRefIdx1++ {
            for iRefIdx2 := iRefIdx1; iRefIdx2 < MAX_NUM_REF; iRefIdx2++ {
                if iRefIdx1 == iRefIdx2 {
                    this.m_abEqualRef[iDir][iRefIdx1][iRefIdx2] = true
                    this.m_abEqualRef[iDir][iRefIdx2][iRefIdx1] = true
                } else {
                    this.m_abEqualRef[iDir][iRefIdx1][iRefIdx2] = false
                    this.m_abEqualRef[iDir][iRefIdx2][iRefIdx1] = false
                }
            }
        }
    }
}
func (this *TComSlice) IsEqualRef(e RefPicList, iRefIdx1 int, iRefIdx2 int) bool {
    if iRefIdx1 < 0 || iRefIdx2 < 0 {
        return false
    }

    return this.m_abEqualRef[e][iRefIdx1][iRefIdx2]
}

func (this *TComSlice) SetEqualRef(e RefPicList, iRefIdx1 int, iRefIdx2 int, b bool) {
    this.m_abEqualRef[e][iRefIdx1][iRefIdx2] = b
    this.m_abEqualRef[e][iRefIdx2][iRefIdx1] = b
}

func /*(this *TComSlice)*/ SortPicList(rcListPic *list.List) {
    var pcPicExtract *TComPic
    var pcPicInsert *TComPic

    for i := 1; i < rcListPic.Len(); i++ {
        iterPicExtract := rcListPic.Front()
        for j := 0; j < i; j++ {
            iterPicExtract = iterPicExtract.Next()
        }
        pcPicExtract = iterPicExtract.Value.(*TComPic)
        pcPicExtract.SetCurrSliceIdx(0)

        iterPicInsert := rcListPic.Front()
        for iterPicInsert != iterPicExtract {
            pcPicInsert = iterPicInsert.Value.(*TComPic)
            pcPicInsert.SetCurrSliceIdx(0)
            if pcPicInsert.GetPOC() >= pcPicExtract.GetPOC() {
                break
            }

            iterPicInsert = iterPicInsert.Next()
        }

        //  swap iterPicExtract and iterPicInsert, iterPicExtract = curr. / iterPicInsert = insertion position
        rcListPic.InsertBefore(pcPicExtract, iterPicInsert) // (, iterPicExtract, iterPicExtract_1);
        rcListPic.Remove(iterPicExtract)
    }
}

func (this *TComSlice) GetNoBackPredFlag() bool {
    return this.m_bNoBackPredFlag
}
func (this *TComSlice) SetNoBackPredFlag(b bool) {
    this.m_bNoBackPredFlag = b
}
func (this *TComSlice) GenerateCombinedList() {
    if this.m_aiNumRefIdx[REF_PIC_LIST_C] > 0 {
        this.m_aiNumRefIdx[REF_PIC_LIST_C] = 0
        for iNumCount := 0; iNumCount < MAX_NUM_REF_LC; iNumCount++ {
            this.m_iRefIdxOfLC[REF_PIC_LIST_0][iNumCount] = -1
            this.m_iRefIdxOfLC[REF_PIC_LIST_1][iNumCount] = -1
            this.m_eListIdFromIdxOfLC[iNumCount] = 0
            this.m_iRefIdxFromIdxOfLC[iNumCount] = 0
            this.m_iRefIdxOfL0FromRefIdxOfL1[iNumCount] = -1
            this.m_iRefIdxOfL1FromRefIdxOfL0[iNumCount] = -1
        }

        for iNumRefIdx := 0; iNumRefIdx < MAX_NUM_REF; iNumRefIdx++ {
            if iNumRefIdx < this.m_aiNumRefIdx[REF_PIC_LIST_0] {
                bTempRefIdxInL2 := true
                for iRefIdxLC := 0; iRefIdxLC < this.m_aiNumRefIdx[REF_PIC_LIST_C]; iRefIdxLC++ {
                    if this.m_apcRefPicList[REF_PIC_LIST_0][iNumRefIdx].GetPOC() == this.m_apcRefPicList[this.m_eListIdFromIdxOfLC[iRefIdxLC]][this.m_iRefIdxFromIdxOfLC[iRefIdxLC]].GetPOC() {
                        this.m_iRefIdxOfL1FromRefIdxOfL0[iNumRefIdx] = this.m_iRefIdxFromIdxOfLC[iRefIdxLC]
                        this.m_iRefIdxOfL0FromRefIdxOfL1[this.m_iRefIdxFromIdxOfLC[iRefIdxLC]] = iNumRefIdx
                        bTempRefIdxInL2 = false
                        break
                    }
                }

                if bTempRefIdxInL2 == true {
                    this.m_eListIdFromIdxOfLC[this.m_aiNumRefIdx[REF_PIC_LIST_C]] = REF_PIC_LIST_0
                    this.m_iRefIdxFromIdxOfLC[this.m_aiNumRefIdx[REF_PIC_LIST_C]] = iNumRefIdx
                    this.m_iRefIdxOfLC[REF_PIC_LIST_0][iNumRefIdx] = this.m_aiNumRefIdx[REF_PIC_LIST_C]
                    this.m_aiNumRefIdx[REF_PIC_LIST_C]++
                }
            }

            if iNumRefIdx < this.m_aiNumRefIdx[REF_PIC_LIST_1] {
                bTempRefIdxInL2 := true
                for iRefIdxLC := 0; iRefIdxLC < this.m_aiNumRefIdx[REF_PIC_LIST_C]; iRefIdxLC++ {
                    if this.m_apcRefPicList[REF_PIC_LIST_1][iNumRefIdx].GetPOC() == this.m_apcRefPicList[this.m_eListIdFromIdxOfLC[iRefIdxLC]][this.m_iRefIdxFromIdxOfLC[iRefIdxLC]].GetPOC() {
                        this.m_iRefIdxOfL0FromRefIdxOfL1[iNumRefIdx] = this.m_iRefIdxFromIdxOfLC[iRefIdxLC]
                        this.m_iRefIdxOfL1FromRefIdxOfL0[this.m_iRefIdxFromIdxOfLC[iRefIdxLC]] = iNumRefIdx
                        bTempRefIdxInL2 = false
                        break
                    }
                }
                if bTempRefIdxInL2 == true {
                    this.m_eListIdFromIdxOfLC[this.m_aiNumRefIdx[REF_PIC_LIST_C]] = REF_PIC_LIST_1
                    this.m_iRefIdxFromIdxOfLC[this.m_aiNumRefIdx[REF_PIC_LIST_C]] = iNumRefIdx
                    this.m_iRefIdxOfLC[REF_PIC_LIST_1][iNumRefIdx] = this.m_aiNumRefIdx[REF_PIC_LIST_C]
                    this.m_aiNumRefIdx[REF_PIC_LIST_C]++
                }
            }
        }
    }
}

func (this *TComSlice) GetTLayer() uint {
    return this.m_uiTLayer
}
func (this *TComSlice) SetTLayer(uiTLayer uint) {
    this.m_uiTLayer = uiTLayer
}

func (this *TComSlice) SetTLayerInfo(uiTLayer uint) {
    this.m_uiTLayer = uiTLayer
}
func (this *TComSlice) DecodingMarking(rcListPic *list.List, iGOPSIze int, iMaxRefPicNum *int) {
    //do nothing
}
func (this *TComSlice) ApplyReferencePictureSet(rcListPic *list.List, pReferencePictureSet *TComReferencePictureSet) {
    var rpcPic *TComPic
    var i, isReference int

    j := 0
    // loop through all pictures in the reference picture buffer
    iterPic := rcListPic.Front()
    for iterPic != nil {
        j++
        rpcPic = iterPic.Value.(*TComPic)
        iterPic = iterPic.Next()

        isReference = 0
        // loop through all pictures in the Reference Picture Set
        // to see if the picture should be kept as reference picture
        for i = 0; i < pReferencePictureSet.GetNumberOfPositivePictures()+pReferencePictureSet.GetNumberOfNegativePictures(); i++ {
            if !rpcPic.GetIsLongTerm() && rpcPic.GetPicSym().GetSlice(0).GetPOC() == this.GetPOC()+pReferencePictureSet.GetDeltaPOC(i) {
                isReference = 1
                rpcPic.SetUsedByCurr(pReferencePictureSet.GetUsed(i))
                rpcPic.SetIsLongTerm(false)
            }
        }
        for ; i < pReferencePictureSet.GetNumberOfPictures(); i++ {
            if pReferencePictureSet.GetCheckLTMSBPresent(i) == true {
                if rpcPic.GetIsLongTerm() && (rpcPic.GetPicSym().GetSlice(0).GetPOC()) == pReferencePictureSet.GetPOC(i) {
                    isReference = 1
                    rpcPic.SetUsedByCurr(pReferencePictureSet.GetUsed(i))
                }
            } else {
                if rpcPic.GetIsLongTerm() && (rpcPic.GetPicSym().GetSlice(0).GetPOC()%(1<<rpcPic.GetPicSym().GetSlice(0).GetSPS().GetBitsForPOC())) == pReferencePictureSet.GetPOC(i)%(1<<rpcPic.GetPicSym().GetSlice(0).GetSPS().GetBitsForPOC()) {
                    isReference = 1
                    rpcPic.SetUsedByCurr(pReferencePictureSet.GetUsed(i))
                }
            }

        }
        // mark the picture as "unused for reference" if it is not in
        // the Reference Picture Set
        if rpcPic.GetPicSym().GetSlice(0).GetPOC() != this.GetPOC() && isReference == 0 {
            rpcPic.GetSlice(0).SetReferenced(false)
            rpcPic.SetIsLongTerm(false)
        }
        if isReference==1 {
     	 	rpcPic.GetSlice(0).SetReferenced( true );  
    	}
        //check that pictures of higher temporal layers are not used
        //assert(rpcPic.GetSlice( 0 )->isReferenced()==0||rpcPic.GetUsedByCurr()==0||rpcPic.GetTLayer()<=this.GetTLayer());
        //check that pictures of higher or equal temporal layer are not in the RPS if the current picture is a TSA picture
        if this.GetNalUnitType() == NAL_UNIT_CODED_SLICE_TLA || this.GetNalUnitType() == NAL_UNIT_CODED_SLICE_TSA_N {
            //assert(rpcPic.GetSlice( 0 )->isReferenced()==0||rpcPic.GetTLayer()<this.GetTLayer());
        }
        //check that pictures marked as temporal layer non-reference pictures are not used for reference
        if rpcPic.GetPicSym().GetSlice(0).GetPOC() != this.GetPOC() && rpcPic.GetTLayer() == this.GetTLayer() {
            //assert(rpcPic.GetSlice( 0 )->isReferenced()==0||rpcPic.GetUsedByCurr()==0||rpcPic.GetSlice( 0 ).GetTemporalLayerNonReferenceFlag()==false);
        }
    }
}
func (this *TComSlice) IsTemporalLayerSwitchingPoint(rcListPic *list.List) bool {
    var rpcPic *TComPic
    // loop through all pictures in the reference picture buffer
    iterPic := rcListPic.Front()
    for iterPic != nil {
        rpcPic = iterPic.Value.(*TComPic)
        iterPic = iterPic.Next()
        if rpcPic.GetSlice(0).IsReferenced() && int(rpcPic.GetPOC()) != this.GetPOC() {
            if rpcPic.GetTLayer() >= this.GetTLayer() {
                return false
            }
        }
    }
    return true
}
func (this *TComSlice) IsStepwiseTemporalLayerSwitchingPointCandidate(rcListPic *list.List) bool {
    var rpcPic *TComPic

    iterPic := rcListPic.Front()
    for iterPic != nil {
        rpcPic = iterPic.Value.(*TComPic)
        iterPic = iterPic.Next()
        if rpcPic.GetSlice(0).IsReferenced() && (rpcPic.GetUsedByCurr() == true) && int(rpcPic.GetPOC()) != this.GetPOC() {
            if rpcPic.GetTLayer() >= this.GetTLayer() {
                return false
            }
        }
    }
    return true
}
func (this *TComSlice) CheckThatAllRefPicsAreAvailable(rcListPic *list.List, pReferencePictureSet *TComReferencePictureSet, printErrors bool, pocRandomAccess int) int {
    var rpcPic *TComPic
    var i, isAvailable, j int
    atLeastOneLost := 0
    atLeastOneRemoved := 0
    iPocLost := 0

    // loop through all long-term pictures in the Reference Picture Set
    // to see if the picture should be kept as reference picture
    for i = pReferencePictureSet.GetNumberOfNegativePictures() + pReferencePictureSet.GetNumberOfPositivePictures(); i < pReferencePictureSet.GetNumberOfPictures(); i++ {
        j = 0
        isAvailable = 0
        // loop through all pictures in the reference picture buffer
        iterPic := rcListPic.Front()
        for iterPic != nil {
            j++
            rpcPic = iterPic.Value.(*TComPic)
            iterPic = iterPic.Next()
            if pReferencePictureSet.GetCheckLTMSBPresent(i) == true {
                if rpcPic.GetIsLongTerm() && (rpcPic.GetPicSym().GetSlice(0).GetPOC()) == pReferencePictureSet.GetPOC(i) && rpcPic.GetSlice(0).IsReferenced() {
                    isAvailable = 1
                }
            } else {
                if rpcPic.GetIsLongTerm() &&
                    (rpcPic.GetPicSym().GetSlice(0).GetPOC()%(1<<rpcPic.GetPicSym().GetSlice(0).GetSPS().GetBitsForPOC())) == pReferencePictureSet.GetPOC(i)%(1<<rpcPic.GetPicSym().GetSlice(0).GetSPS().GetBitsForPOC()) &&
                    rpcPic.GetSlice(0).IsReferenced() {
                    isAvailable = 1
                }
            }
        }
        // if there was no such long-term check the short terms
        if isAvailable == 0 {
            iterPic = rcListPic.Front()
            for iterPic != nil {
                j++
                rpcPic = iterPic.Value.(*TComPic)
                iterPic = iterPic.Next()
                if (rpcPic.GetPicSym().GetSlice(0).GetPOC()%(1<<rpcPic.GetPicSym().GetSlice(0).GetSPS().GetBitsForPOC())) == (this.GetPOC()+pReferencePictureSet.GetDeltaPOC(i))%(1<<rpcPic.GetPicSym().GetSlice(0).GetSPS().GetBitsForPOC()) && rpcPic.GetSlice(0).IsReferenced() {
                    isAvailable = 1
                    rpcPic.SetIsLongTerm(true)
                    break
                }
            }
        }
        // report that a picture is lost if it is in the Reference Picture Set
        // but not available as reference picture
        if isAvailable == 0 {
            if this.GetPOC()+pReferencePictureSet.GetDeltaPOC(i) >= pocRandomAccess {
                if !pReferencePictureSet.GetUsed(i) {
                    if printErrors {
                        fmt.Printf("\nLong-term reference picture with POC = %3d seems to have been removed or not correctly decoded.", this.GetPOC()+pReferencePictureSet.GetDeltaPOC(i))
                    }
                    atLeastOneRemoved = 1
                } else {
                    if printErrors {
                        fmt.Printf("\nLong-term reference picture with POC = %3d is lost or not correctly decoded!", this.GetPOC()+pReferencePictureSet.GetDeltaPOC(i))
                    }
                    atLeastOneLost = 1
                    iPocLost = this.GetPOC() + pReferencePictureSet.GetDeltaPOC(i)
                }
            }
        }
    }
    // loop through all short-term pictures in the Reference Picture Set
    // to see if the picture should be kept as reference picture
    for i = 0; i < pReferencePictureSet.GetNumberOfNegativePictures()+pReferencePictureSet.GetNumberOfPositivePictures(); i++ {
        j = 0
        isAvailable = 0
        // loop through all pictures in the reference picture buffer
        iterPic := rcListPic.Front()
        for iterPic != nil {
            j++
            rpcPic = iterPic.Value.(*TComPic)
            iterPic = iterPic.Next()
            if !rpcPic.GetIsLongTerm() && rpcPic.GetPicSym().GetSlice(0).GetPOC() == this.GetPOC()+pReferencePictureSet.GetDeltaPOC(i) && rpcPic.GetSlice(0).IsReferenced() {
                isAvailable = 1
            }
        }
        // report that a picture is lost if it is in the Reference Picture Set
        // but not available as reference picture
        if isAvailable == 0 {
            if this.GetPOC()+pReferencePictureSet.GetDeltaPOC(i) >= pocRandomAccess {
                if !pReferencePictureSet.GetUsed(i) {
                    if printErrors {
                        fmt.Printf("\nShort-term reference picture with POC = %3d seems to have been removed or not correctly decoded.", this.GetPOC()+pReferencePictureSet.GetDeltaPOC(i))
                    }
                    atLeastOneRemoved = 1
                } else {
                    if printErrors {
                        fmt.Printf("\nShort-term reference picture with POC = %3d is lost or not correctly decoded!", this.GetPOC()+pReferencePictureSet.GetDeltaPOC(i))
                    }
                    atLeastOneLost = 1
                    iPocLost = this.GetPOC() + pReferencePictureSet.GetDeltaPOC(i)
                }
            }
        }
    }
    if atLeastOneLost != 0 {
        return iPocLost + 1
    }
    if atLeastOneRemoved != 0 {
        return -2
    }

    return 0
}
func (this *TComSlice) CreateExplicitReferencePictureSetFromReference(rcListPic *list.List, pReferencePictureSet *TComReferencePictureSet) {
    var rpcPic *TComPic
    var i, j int
    k := 0
    nrOfNegativePictures := 0
    nrOfPositivePictures := 0
    pcRPS := this.GetLocalRPS()

    // loop through all pictures in the Reference Picture Set
    for i = 0; i < pReferencePictureSet.GetNumberOfPictures(); i++ {
        j = 0
        // loop through all pictures in the reference picture buffer
        iterPic := rcListPic.Front()
        for iterPic != nil {
            j++
            rpcPic = iterPic.Value.(*TComPic)
            iterPic = iterPic.Next()

            if rpcPic.GetPicSym().GetSlice(0).GetPOC() == this.GetPOC()+pReferencePictureSet.GetDeltaPOC(i) && rpcPic.GetSlice(0).IsReferenced() {
                // This picture exists as a reference picture
                // and should be added to the explicit Reference Picture Set
                pcRPS.SetDeltaPOC(k, pReferencePictureSet.GetDeltaPOC(i))
                pcRPS.SetUsed(k, pReferencePictureSet.GetUsed(i))
                if pcRPS.GetDeltaPOC(k) < 0 {
                    nrOfNegativePictures++
                } else {
                    nrOfPositivePictures++
                }
                k++
            }
        }
    }
    pcRPS.SetNumberOfNegativePictures(nrOfNegativePictures)
    pcRPS.SetNumberOfPositivePictures(nrOfPositivePictures)
    pcRPS.SetNumberOfPictures(nrOfNegativePictures + nrOfPositivePictures)
    // This is a simplistic inter rps example. A smarter encoder will look for a better reference RPS to do the
    // inter RPS prediction with.  Here we just use the reference used by pReferencePictureSet.
    // If pReferencePictureSet is not inter_RPS_predicted, then inter_RPS_prediction is for the current RPS also disabled.
    if !pReferencePictureSet.GetInterRPSPrediction() {
        pcRPS.SetInterRPSPrediction(false)
        pcRPS.SetNumRefIdc(0)
    } else {
        rIdx := this.GetRPSidx() - pReferencePictureSet.GetDeltaRIdxMinus1() - 1
        deltaRPS := pReferencePictureSet.GetDeltaRPS()
        pcRefRPS := this.GetSPS().GetRPSList().GetReferencePictureSet(rIdx)
        iRefPics := pcRefRPS.GetNumberOfPictures()
        iNewIdc := 0
        for i = 0; i <= iRefPics; i++ {
            var deltaPOC int
            if i != iRefPics {
                deltaPOC = pcRefRPS.GetDeltaPOC(i) // check if the reference abs POC is >= 0
            } else {
                deltaPOC = 0 // check if the reference abs POC is >= 0
            }
            iRefIdc := 0
            for j = 0; j < pcRPS.GetNumberOfPictures(); j++ { // loop through the  pictures in the new RPS
                if (deltaPOC + deltaRPS) == pcRPS.GetDeltaPOC(j) {
                    if pcRPS.GetUsed(j) {
                        iRefIdc = 1
                    } else {
                        iRefIdc = 2
                    }
                }
            }
            pcRPS.SetRefIdc(i, iRefIdc)
            iNewIdc++
        }
        pcRPS.SetInterRPSPrediction(true)
        pcRPS.SetNumRefIdc(iNewIdc)
        pcRPS.SetDeltaRPS(deltaRPS)
        pcRPS.SetDeltaRIdxMinus1(pReferencePictureSet.GetDeltaRIdxMinus1() + this.GetSPS().GetRPSList().GetNumberOfReferencePictureSets() - this.GetRPSidx())
    }

    this.SetRPS(pcRPS)
    this.SetRPSidx(-1)
}

func (this *TComSlice) SetMaxNumMergeCand(val uint) {
    this.m_maxNumMergeCand = val
}
func (this *TComSlice) GetMaxNumMergeCand() uint {
    return this.m_maxNumMergeCand
}

func (this *TComSlice) SetSliceMode(uiMode uint) { this.m_sliceMode = uiMode }
func (this *TComSlice) GetSliceMode() uint       { return this.m_sliceMode }
func (this *TComSlice) SetSliceArgument(uiArgument uint) {
    this.m_sliceArgument = uiArgument
}
func (this *TComSlice) GetSliceArgument() uint { return this.m_sliceArgument }
func (this *TComSlice) SetSliceCurStartCUAddr(uiAddr uint) {
    this.m_sliceCurStartCUAddr = uiAddr
}
func (this *TComSlice) GetSliceCurStartCUAddr() uint { return this.m_sliceCurStartCUAddr }
func (this *TComSlice) SetSliceCurEndCUAddr(uiAddr uint) {
    this.m_sliceCurEndCUAddr = uiAddr
}
func (this *TComSlice) GetSliceCurEndCUAddr() uint { return this.m_sliceCurEndCUAddr }
func (this *TComSlice) SetSliceIdx(i uint)         { this.m_sliceIdx = i }
func (this *TComSlice) GetSliceIdx() uint          { return this.m_sliceIdx }
func (this *TComSlice) CopySliceInfo(pSrc *TComSlice) {
    //assert( pSrc != NULL );

    var i, j, k int

    this.m_iPOC = pSrc.m_iPOC
    this.m_eNalUnitType = pSrc.m_eNalUnitType
    this.m_eSliceType = pSrc.m_eSliceType
    this.m_iSliceQp = pSrc.m_iSliceQp
    //#if ADAPTIVE_QP_SELECTION
    this.m_iSliceQpBase = pSrc.m_iSliceQpBase
    //#endif
    this.m_deblockingFilterDisable = pSrc.m_deblockingFilterDisable
    this.m_deblockingFilterOverrideFlag = pSrc.m_deblockingFilterOverrideFlag
    this.m_deblockingFilterBetaOffsetDiv2 = pSrc.m_deblockingFilterBetaOffsetDiv2
    this.m_deblockingFilterTcOffsetDiv2 = pSrc.m_deblockingFilterTcOffsetDiv2

    for i = 0; i < 3; i++ {
        this.m_aiNumRefIdx[i] = pSrc.m_aiNumRefIdx[i]
    }

    for i = 0; i < 2; i++ {
        for j = 0; j < MAX_NUM_REF_LC; j++ {
            this.m_iRefIdxOfLC[i][j] = pSrc.m_iRefIdxOfLC[i][j]
        }
    }
    for i = 0; i < MAX_NUM_REF_LC; i++ {
        this.m_eListIdFromIdxOfLC[i] = pSrc.m_eListIdFromIdxOfLC[i]
        this.m_iRefIdxFromIdxOfLC[i] = pSrc.m_iRefIdxFromIdxOfLC[i]
        this.m_iRefIdxOfL1FromRefIdxOfL0[i] = pSrc.m_iRefIdxOfL1FromRefIdxOfL0[i]
        this.m_iRefIdxOfL0FromRefIdxOfL1[i] = pSrc.m_iRefIdxOfL0FromRefIdxOfL1[i]
    }
    this.m_bRefPicListModificationFlagLC = pSrc.m_bRefPicListModificationFlagLC
    this.m_bRefPicListCombinationFlag = pSrc.m_bRefPicListCombinationFlag
    this.m_bCheckLDC = pSrc.m_bCheckLDC
    this.m_iSliceQpDelta = pSrc.m_iSliceQpDelta
    this.m_iSliceQpDeltaCb = pSrc.m_iSliceQpDeltaCb
    this.m_iSliceQpDeltaCr = pSrc.m_iSliceQpDeltaCr
    for i = 0; i < 2; i++ {
        for j = 0; j < MAX_NUM_REF; j++ {
            this.m_apcRefPicList[i][j] = pSrc.m_apcRefPicList[i][j]
            this.m_aiRefPOCList[i][j] = pSrc.m_aiRefPOCList[i][j]
        }
    }
    for i = 0; i < 2; i++ {
        for j = 0; j < MAX_NUM_REF+1; j++ {
            this.m_bIsUsedAsLongTerm[i][j] = pSrc.m_bIsUsedAsLongTerm[i][j]
        }
    }
    this.m_iDepth = pSrc.m_iDepth

    // referenced slice
    this.m_bRefenced = pSrc.m_bRefenced

    // access channel
    this.m_pcSPS = pSrc.m_pcSPS
    this.m_pcPPS = pSrc.m_pcPPS
    this.m_pcRPS = pSrc.m_pcRPS
    this.m_iLastIDR = pSrc.m_iLastIDR

    this.m_pcPic = pSrc.m_pcPic

    this.m_colFromL0Flag = pSrc.m_colFromL0Flag
    this.m_colRefIdx = pSrc.m_colRefIdx
    //#if SAO_CHROMA_LAMBDA
    this.m_dLambdaLuma = pSrc.m_dLambdaLuma
    this.m_dLambdaChroma = pSrc.m_dLambdaChroma
    //#else
    //  m_dLambda              = pSrc.m_dLambda;
    //#endif
    for i = 0; i < 2; i++ {
        for j = 0; j < MAX_NUM_REF; j++ {
            for k = 0; k < MAX_NUM_REF; k++ {
                this.m_abEqualRef[i][j][k] = pSrc.m_abEqualRef[i][j][k]
            }
        }
    }

    this.m_bNoBackPredFlag = pSrc.m_bNoBackPredFlag
    this.m_uiTLayer = pSrc.m_uiTLayer
    this.m_bTLayerSwitchingFlag = pSrc.m_bTLayerSwitchingFlag

    this.m_sliceMode = pSrc.m_sliceMode
    this.m_sliceArgument = pSrc.m_sliceArgument
    this.m_sliceCurStartCUAddr = pSrc.m_sliceCurStartCUAddr
    this.m_sliceCurEndCUAddr = pSrc.m_sliceCurEndCUAddr
    this.m_sliceIdx = pSrc.m_sliceIdx
    this.m_sliceSegmentMode = pSrc.m_sliceSegmentMode
    this.m_sliceSegmentArgument = pSrc.m_sliceSegmentArgument
    this.m_sliceSegmentCurStartCUAddr = pSrc.m_sliceSegmentCurStartCUAddr
    this.m_sliceSegmentCurEndCUAddr = pSrc.m_sliceSegmentCurEndCUAddr
    this.m_nextSlice = pSrc.m_nextSlice
    this.m_nextSliceSegment = pSrc.m_nextSliceSegment
    for e := 0; e < 2; e++ {
        for n := 0; n < MAX_NUM_REF; n++ {
            for m := 0; m < 3; m++ {
                this.m_weightPredTable[e][n][m] = pSrc.m_weightPredTable[e][n][m]
            }
            //memcpy(this.m_weightPredTable[e][n], pSrc.m_weightPredTable[e][n], sizeof(WpScalingParam)*3 );
        }
    }
    this.m_saoEnabledFlag = pSrc.m_saoEnabledFlag
    this.m_saoEnabledFlagChroma = pSrc.m_saoEnabledFlagChroma
    this.m_cabacInitFlag = pSrc.m_cabacInitFlag
    this.m_numEntryPointOffsets = pSrc.m_numEntryPointOffsets

    this.m_bLMvdL1Zero = pSrc.m_bLMvdL1Zero
    this.m_LFCrossSliceBoundaryFlag = pSrc.m_LFCrossSliceBoundaryFlag
    this.m_enableTMVPFlag = pSrc.m_enableTMVPFlag
    this.m_maxNumMergeCand = pSrc.m_maxNumMergeCand
}
func (this *TComSlice) SetSliceSegmentMode(uiMode uint) {
    this.m_sliceSegmentMode = uiMode
}
func (this *TComSlice) GetSliceSegmentMode() uint { return this.m_sliceSegmentMode }
func (this *TComSlice) SetSliceSegmentArgument(uiArgument uint) {
    this.m_sliceSegmentArgument = uiArgument
}
func (this *TComSlice) GetSliceSegmentArgument() uint { return this.m_sliceSegmentArgument }
func (this *TComSlice) SetSliceSegmentCurStartCUAddr(uiAddr uint) {
    this.m_sliceSegmentCurStartCUAddr = uiAddr
}
func (this *TComSlice) GetSliceSegmentCurStartCUAddr() uint {
    return this.m_sliceSegmentCurStartCUAddr
}
func (this *TComSlice) SetSliceSegmentCurEndCUAddr(uiAddr uint) {
    this.m_sliceSegmentCurEndCUAddr = uiAddr
}
func (this *TComSlice) GetSliceSegmentCurEndCUAddr() uint {
    return this.m_sliceSegmentCurEndCUAddr
}
func (this *TComSlice) SetNextSlice(b bool)        { this.m_nextSlice = b }
func (this *TComSlice) IsNextSlice() bool          { return this.m_nextSlice }
func (this *TComSlice) SetNextSliceSegment(b bool) { this.m_nextSliceSegment = b }
func (this *TComSlice) IsNextSliceSegment() bool   { return this.m_nextSliceSegment }
func (this *TComSlice) SetSliceBits(uiVal uint)    { this.m_sliceBits = uiVal }
func (this *TComSlice) GetSliceBits() uint         { return this.m_sliceBits }
func (this *TComSlice) SetSliceSegmentBits(uiVal uint) {
    this.m_sliceSegmentBits = uiVal
}
func (this *TComSlice) GetSliceSegmentBits() uint { return this.m_sliceSegmentBits }
func (this *TComSlice) SetFinalized(uiVal bool) {
    this.m_bFinalized = uiVal
}
func (this *TComSlice) GetFinalized() bool {
    return this.m_bFinalized
}
func (this *TComSlice) SetWpScaling(wp [2][MAX_NUM_REF][3]WpScalingParam) {
    //memcpy(this.m_weightPredTable, wp, sizeof(WpScalingParam)*2*MAX_NUM_REF*3);
    this.m_weightPredTable = wp
}
func (this *TComSlice) GetWpScaling(e RefPicList, iRefIdx int) []WpScalingParam {
    return this.m_weightPredTable[e][iRefIdx][:]
}

func (this *TComSlice) ResetWpScaling(wp [2][MAX_NUM_REF][3]WpScalingParam) {
    for e := 0; e < 2; e++ {
        for i := 0; i < MAX_NUM_REF; i++ {
            for yuv := 0; yuv < 3; yuv++ {
                wp[e][i][yuv].bPresentFlag = false
                wp[e][i][yuv].uiLog2WeightDenom = 0
                wp[e][i][yuv].uiLog2WeightDenom = 0
                wp[e][i][yuv].iWeight = 1
                wp[e][i][yuv].iOffset = 0
            }
        }
    }
}
func (this *TComSlice) InitWpScaling1(wp [2][MAX_NUM_REF][3]WpScalingParam) {
    for e := 0; e < 2; e++ {
        for i := 0; i < MAX_NUM_REF; i++ {
            for yuv := 0; yuv < 3; yuv++ {
                if !wp[e][i][yuv].bPresentFlag {
                    // Inferring values not present :
                    wp[e][i][yuv].iWeight = (1 << wp[e][i][yuv].uiLog2WeightDenom)
                    wp[e][i][yuv].iOffset = 0
                }

                wp[e][i][yuv].W = wp[e][i][yuv].iWeight
                var bitDepth uint
                if yuv != 0 {
                    bitDepth = uint(G_bitDepthC)
                } else {
                    bitDepth = uint(G_bitDepthY)
                }
                wp[e][i][yuv].O = wp[e][i][yuv].iOffset << (bitDepth - 8)
                wp[e][i][yuv].Shift = int(wp[e][i][yuv].uiLog2WeightDenom)
                if wp[e][i][yuv].uiLog2WeightDenom >= 1 {
                    wp[e][i][yuv].Round = (1 << (wp[e][i][yuv].uiLog2WeightDenom - 1))
                } else {
                    wp[e][i][yuv].Round = (0)
                }
            }
        }
    }
}
func (this *TComSlice) InitWpScaling() {
    this.InitWpScaling1(this.m_weightPredTable)
}
func (this *TComSlice) ApplyWP() bool {
    return ((this.m_eSliceType == P_SLICE && this.m_pcPPS.GetUseWP()) || (this.m_eSliceType == B_SLICE && this.m_pcPPS.GetWPBiPred()))
}

func (this *TComSlice) SetWpAcDcParam(wp [3]wpACDCParam) {
    //memcpy(this.m_weightACDCParam, wp, sizeof(wpACDCParam)*3);
    this.m_weightACDCParam = wp
}
func (this *TComSlice) GetWpAcDcParam() [3]wpACDCParam {
    return this.m_weightACDCParam
}
func (this *TComSlice) InitWpAcDcParam() {
    for iComp := 0; iComp < 3; iComp++ {
        this.m_weightACDCParam[iComp].iAC = 0
        this.m_weightACDCParam[iComp].iDC = 0
    }
}

func (this *TComSlice) SetTileLocationCount1(cnt uint) {
    this.m_tileByteLocation = make(map[int]uint, cnt)
}
func (this *TComSlice) GetTileLocationCount() uint {
    return uint(len(this.m_tileByteLocation))
}
func (this *TComSlice) SetTileLocation(idx int, location uint) {
    //assert (idx<this.m_tileByteLocation.size());
    this.m_tileByteLocation[idx] = location
}
func (this *TComSlice) AddTileLocation(location uint) {
    this.m_tileByteLocation[len(this.m_tileByteLocation)+1] = location
}
func (this *TComSlice) GetTileLocation(idx int) uint {
    return this.m_tileByteLocation[idx]
}

func (this *TComSlice) SetTileOffstForMultES(uiOffset uint) {
    this.m_uiTileOffstForMultES = uiOffset
}
func (this *TComSlice) GetTileOffstForMultES() uint {
    return this.m_uiTileOffstForMultES
}
func (this *TComSlice) AllocSubstreamSizes(uiNumSubstreams uint) {
    //delete[] m_puiSubstreamSizes;
    if uiNumSubstreams > 0 {
        this.m_puiSubstreamSizes = make([]uint, uiNumSubstreams-1)
    } else {
        this.m_puiSubstreamSizes = make([]uint, 0)
    }
}
func (this *TComSlice) GetSubstreamSizes() []uint {
    return this.m_puiSubstreamSizes
}
func (this *TComSlice) SetScalingList(scalingList *TComScalingList) {
    this.m_scalingList = scalingList
}
func (this *TComSlice) GetScalingList() *TComScalingList {
    return this.m_scalingList
}
func (this *TComSlice) SetDefaultScalingList() {
    for sizeId := 0; sizeId < SCALING_LIST_SIZE_NUM; sizeId++ {
        for listId := uint(0); listId < G_scalingListNum[sizeId]; listId++ {
            this.GetScalingList().ProcessDefaultMarix(uint(sizeId), listId)
        }
    }
}
func (this *TComSlice) CheckDefaultScalingList() bool {
    /* Encoder func

         defaultCounter:=uint(0);

         for sizeId := uint(0); sizeId < SCALING_LIST_SIZE_NUM; sizeId++ {
           for listId:= uint(0);listId<G_scalingListNum[sizeId];listId++ {
           	slDstAddr := this.GetScalingList().GetScalingListAddress(sizeId,listId);
           	slSrcAddr := this.GetScalingList().GetScalingListDefaultAddress(sizeId, listId);

       		sizeof(Int)*min(MAX_MATRIX_COEF_NUM,G_scalingListSize[sizeId])
             if( !memcmp() // check value of matrix
            && ((sizeId < SCALING_LIST_16x16) || (getScalingList().GetScalingListDC(sizeId,listId) == 16))) // check DC value
             {
               defaultCounter++;
             }
           }
         }
         return (defaultCounter == (SCALING_LIST_NUM * SCALING_LIST_SIZE_NUM - 4)) ? false : true; // -4 for 32x32
    */
    return true
}
func (this *TComSlice) SetCabacInitFlag(val bool) {
    this.m_cabacInitFlag = val
}   //!< Set CABAC initial flag
func (this *TComSlice) GetCabacInitFlag() bool {
    return this.m_cabacInitFlag
}   //!< Get CABAC initial flag
func (this *TComSlice) SetNumEntryPointOffsets(val int) {
    this.m_numEntryPointOffsets = val
}
func (this *TComSlice) GetNumEntryPointOffsets() int {
    return this.m_numEntryPointOffsets
}
func (this *TComSlice) GetTemporalLayerNonReferenceFlag() bool {
    return this.m_temporalLayerNonReferenceFlag
}
func (this *TComSlice) SetTemporalLayerNonReferenceFlag(x bool) {
    this.m_temporalLayerNonReferenceFlag = x
}
func (this *TComSlice) SetLFCrossSliceBoundaryFlag(val bool) {
    this.m_LFCrossSliceBoundaryFlag = val
}
func (this *TComSlice) GetLFCrossSliceBoundaryFlag() bool {
    return this.m_LFCrossSliceBoundaryFlag
}

func (this *TComSlice) SetEnableTMVPFlag(b bool) {
    this.m_enableTMVPFlag = b
}
func (this *TComSlice) GetEnableTMVPFlag() bool {
    return this.m_enableTMVPFlag
}

//protected:
func (this *TComSlice) xGetRefPic(rcListPic *list.List, poc int) *TComPic {
    var pcPic *TComPic

    for e := rcListPic.Front(); e != nil; e = e.Next() {
        pcPic = e.Value.(*TComPic)
        if pcPic.GetPOC() == poc {
            break
        }
    }

    return pcPic
}

func (this *TComSlice) xGetLongTermRefPic(rcListPic *list.List, poc int, bCheckLTMSBPresent bool) *TComPic {
    var pcPic *TComPic
    var pcStPic *TComPic
    
    if rcListPic.Front()!=nil {
    	pcStPic = rcListPic.Front().Value.(*TComPic);
    }else{
    	pcStPic = nil;
    }
    
    for e := rcListPic.Front(); e != nil; e = e.Next() {
        pcPic = e.Value.(*TComPic)
        if bCheckLTMSBPresent==true {
	      if pcPic != nil && pcPic.GetPOC() == poc {
	        if pcPic.GetIsLongTerm() {
	          return pcPic;
	        }else {
	          pcStPic = pcPic;
	        }
	        break;
	      }
	    }else{
	        if pcPic != nil && pcPic.GetPOC()%int(1<<this.GetSPS().GetBitsForPOC()) == poc%int(1<<this.GetSPS().GetBitsForPOC()) {
	            if pcPic.GetIsLongTerm() {
	                return pcPic
	            } else {
	                pcStPic = pcPic
	            }
	            break
	        }
	    }
    }

    return pcStPic
}

//};// END CLASS DEFINITION TComSlice

/*
type ParameterSetMap struct{
//private:
  m_maxId	int;
  m_paramsetMap map[int]interface{};
};

//public:
func NewParameterSetMap(maxId int) *ParameterSetMap{
	return &ParameterSetMap{m_maxId:maxId}
}

func (this *ParameterSetMap) StorePS(psId int, ps interface{}){
    //assert ( psId < m_maxId );
    m_paramsetMap[psId] = ps;
}
func (this *ParameterSetMap) MergePSList(rPsList *ParameterSetMap){
    for id, ps := this.m_paramsetMap {
      storePS(i->first, i->second);
    }
}


func (this *ParameterSetMap) GetPS(psId int) interface{}{
	value, ok := m_paramsetMap[psId];
	if ok {
		return value
	}

	return nil
}

 T* getFirstPS()
  {
    return (m_paramsetMap.begin() == m_paramsetMap.end() ) ? NULL : m_paramsetMap.begin()->second;
  }*/

type ParameterSetManager struct {
    m_vpsMap map[int]*TComVPS
    m_spsMap map[int]*TComSPS
    m_ppsMap map[int]*TComPPS

    m_activeVPSId int
    m_activeSPSId int
    m_activePPSId int
}

//public:
func NewParameterSetManager() *ParameterSetManager {
    return &ParameterSetManager{make(map[int]*TComVPS), make(map[int]*TComSPS), make(map[int]*TComPPS), -1, -1, -1}
}

//! store sequence parameter set and take ownership of it
func (this *ParameterSetManager) SetVPS(vps *TComVPS) {
    this.m_vpsMap[vps.GetVPSId()] = vps
}

//! get pointer to existing video parameter set
func (this *ParameterSetManager) GetVPS(vpsId int) *TComVPS {
    return this.m_vpsMap[vpsId]
}

//func (this *ParameterSetManager)  TComVPS* getFirstVPS()      { return m_vpsMap.getFirstPS(); };

//! store sequence parameter set and take ownership of it
func (this *ParameterSetManager) SetSPS(sps *TComSPS) {
    this.m_spsMap[sps.GetSPSId()] = sps
}

//! get pointer to existing sequence parameter set
func (this *ParameterSetManager) GetSPS(spsId int) *TComSPS {
    return this.m_spsMap[spsId]
}

//func (this *ParameterSetManager)  TComSPS* getFirstSPS()      { return m_spsMap.getFirstPS(); };

//! store picture parameter set and take ownership of it
func (this *ParameterSetManager) SetPPS(pps *TComPPS) {
    this.m_ppsMap[pps.GetPPSId()] = pps
}

//! get pointer to existing picture parameter set
func (this *ParameterSetManager) GetPPS(ppsId int) *TComPPS {
    return this.m_ppsMap[ppsId]
}

func (this *ParameterSetManager) ApplyPS() {
}

func (this *ParameterSetManager) GetActiveVPS() *TComVPS { return this.m_vpsMap[this.m_activeVPSId] }
func (this *ParameterSetManager) GetActiveSPS() *TComSPS { return this.m_spsMap[this.m_activeSPSId] }
func (this *ParameterSetManager) GetActivePPS() *TComPPS { return this.m_ppsMap[this.m_activePPSId] }

func (this *ParameterSetManager) ActivatePPS(ppsId int, isIDR bool) bool {
    pps := this.m_ppsMap[ppsId]
    if pps != nil {
        spsId := pps.GetSPSId()
        if !isIDR && (spsId != this.m_activeSPSId) {
            fmt.Printf("Warning: tried to activate PPS referring to a inactive SPS at non-IDR.")
            return false
        }
        sps := this.m_spsMap[spsId]
        if sps != nil {
            vpsId := sps.GetVPSId()
            if !isIDR && (vpsId != this.m_activeVPSId) {
                fmt.Printf("Warning: tried to activate PPS referring to a inactive VPS at non-IDR.")
                return false
            }
            if this.m_vpsMap[vpsId] != nil {
                this.m_activePPSId = ppsId
                this.m_activeVPSId = vpsId
                this.m_activeSPSId = spsId
                return true
            } else {
                fmt.Printf("Warning: tried to activate PPS that refers to a non-existing VPS.")
            }
        } else {
            fmt.Printf("Warning: tried to activate a PPS that refers to a non-existing SPS.")
        }
    } else {
        fmt.Printf("Warning: tried to activate non-existing PPS.")
    }
    return false
}

//func (this *ParameterSetManager)  TComPPS* getFirstPPS()      { return m_ppsMap.getFirstPS(); };
