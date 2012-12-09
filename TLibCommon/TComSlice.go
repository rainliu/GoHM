package TLibCommon

import (
	"fmt"
    "container/list"
)

// ====================================================================================================================
// Constants
// ====================================================================================================================

/// max number of supported APS in software
const MAX_NUM_SUPPORTED_APS = 1

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
    return &TComReferencePictureSet{}
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
	this.m_used[bufferNum] = used;
}
func (this *TComReferencePictureSet) SetDeltaPOC(bufferNum, deltaPOC int) {
	this.m_deltaPOC[bufferNum] = deltaPOC;
}
func (this *TComReferencePictureSet) SetPOC(bufferNum, POC int) {
	this.m_POC[bufferNum] = POC;
}
func (this *TComReferencePictureSet) SetNumberOfPictures(numberOfPictures int) {
	this.m_numberOfPictures = numberOfPictures;
}
func (this *TComReferencePictureSet) SetCheckLTMSBPresent(bufferNum int, b bool) {
  	this.m_bCheckLTMSB[bufferNum] = b;
}
func (this *TComReferencePictureSet) GetCheckLTMSBPresent(bufferNum int) bool {
    return this.m_bCheckLTMSB[bufferNum];
}

func (this *TComReferencePictureSet) GetUsed(bufferNum int) bool {
    return this.m_used[bufferNum];
}
func (this *TComReferencePictureSet) GetDeltaPOC(bufferNum int) int {
    return this.m_deltaPOC[bufferNum];
}
func (this *TComReferencePictureSet) GetPOC(bufferNum int) int {
    return this.m_POC[bufferNum];
}
func (this *TComReferencePictureSet) GetNumberOfPictures() int {
    return this.m_numberOfPictures;
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
  for j:=1; j < this.GetNumberOfPictures(); j++ { 
    deltaPOC := this.GetDeltaPOC(j);
    used := this.GetUsed(j);
    for k:=j-1; k >= 0; k-- {
      temp := this.GetDeltaPOC(k);
      if deltaPOC < temp {
        this.SetDeltaPOC(k+1, temp);
        this.SetUsed(k+1, this.GetUsed(k));
        this.SetDeltaPOC(k, deltaPOC);
        this.SetUsed(k, used);
      }
    }
  }
  // flip the negative values to largest first
  numNegPics := this.GetNumberOfNegativePictures();
  k:=numNegPics-1;
  for j:=0; j < numNegPics>>1; j++ { 
    deltaPOC := this.GetDeltaPOC(j);
    used := this.GetUsed(j);
    this.SetDeltaPOC(j, this.GetDeltaPOC(k));
    this.SetUsed(j, this.GetUsed(k));
    this.SetDeltaPOC(k, deltaPOC);
    this.SetUsed(k, used);
    k--;
  }
}
func (this *TComReferencePictureSet) PrintDeltaPOC() {
  fmt.Printf("DeltaPOC = { ");
  for j:=0; j < this.GetNumberOfPictures(); j++ {
  	if this.GetUsed(j) {
    	fmt.Printf("%d%s ", this.GetDeltaPOC(j), "*");
  	}else{
  		fmt.Printf("%d%s ", this.GetDeltaPOC(j), "");
  	}
  } 
  if this.GetInterRPSPrediction() {
    fmt.Printf("}, RefIdc = { ");
    for j:=0; j < this.GetNumRefIdc(); j++ {
      fmt.Printf("%d ", this.GetRefIdc(j));
    } 
  }
  fmt.Printf("}\n");
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
    m_scalingListCoef             [][SCALING_LIST_SIZE_NUM][SCALING_LIST_NUM]int //!< quantization matrix
    m_useTransformSkip            bool
}

//public:
func NewTComScalingList() *TComScalingList {
    return &TComScalingList{}
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
}   //!< get matrix coefficient
func (this *TComScalingList) CheckPredMode(sizeId, listId uint) bool {
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
        //#if FLAT_4x4_DSL
        src = G_quantTSDefault4x4[:]
        /*#else
              if( m_useTransformSkip )
              {
                src = g_quantTSDefault4x4;
              }
              else
              {
                src = (listId<3) ? g_quantIntraDefault4x4 : g_quantInterDefault4x4;
              }
        #endif*/
        //break;
    case SCALING_LIST_8x8:
        if listId < 3 {
            src = G_quantIntraDefault8x8[:]
        } else {
            src = G_quantInterDefault8x8[:]
        }
        //src = (listId<3) ? g_quantIntraDefault8x8 : g_quantInterDefault8x8;
        //break;
    case SCALING_LIST_16x16:
        if listId < 3 {
            src = G_quantIntraDefault8x8[:]
        } else {
            src = G_quantInterDefault8x8[:]
        }
        //src = (listId<3) ? g_quantIntraDefault8x8 : g_quantInterDefault8x8;
        //break;
    case SCALING_LIST_32x32:
        if listId < 1 {
            src = G_quantIntraDefault8x8[:]
        } else {
            src = G_quantInterDefault8x8[:]
        }
        //src = (listId<1) ? g_quantIntraDefault8x8 : g_quantInterDefault8x8;
        //break;
    default:
        //  assert(0);
        src = nil //NULL;
        //break;
    }
    return src
}   //!< get default matrix coefficient
func (this *TComScalingList) ProcessDefaultMarix(sizeId, listId uint) {
}
func (this *TComScalingList) SetScalingListDC(sizeId, listId, u uint) {
    this.m_scalingListDC[sizeId][listId] = int(u)
}   //!< set DC value

func (this *TComScalingList) GetScalingListDC(sizeId, listId uint) int {
    return this.m_scalingListDC[sizeId][listId]
}   //!< get DC value
func (this *TComScalingList) CheckDcOfMatrix() {
}
func (this *TComScalingList) ProcessRefMatrix(sizeId, listId, refListId uint) {
}
func (this *TComScalingList) XParseScalingList(pchFile string) bool {
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
func (this *TComScalingList) init() {
}
func (this *TComScalingList) destroy() {
}

//!< transform skipping flag for setting default scaling matrix for 4x4

type ProfileTierLevel struct {
    m_profileSpace             int
    m_tierFlag                 bool
    m_profileIdc               int
    m_profileCompatibilityFlag [32]bool
    m_levelIdc                 int
}

//public:
func NewProfileTierLevel() *ProfileTierLevel {
    return &ProfileTierLevel{}
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

type TComPTL struct {
    m_generalPTL                 ProfileTierLevel
    m_subLayerPTL                [6]ProfileTierLevel // max. value of max_sub_layers_minus1 is 6
    m_subLayerProfilePresentFlag [6]bool
    m_subLayerLevelPresentFlag   [6]bool
}

//public:
func NewTComPTL() *TComPTL {
    return &TComPTL{}
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

/// VPS class

//#if SIGNAL_BITRATE_PICRATE_IN_VPS
type TComBitRatePicRateInfo struct{
  m_bitRateInfoPresentFlag	[MAX_TLAYER]bool;
  m_picRateInfoPresentFlag	[MAX_TLAYER]bool;
  m_avgBitRate				[MAX_TLAYER]int;
  m_maxBitRate				[MAX_TLAYER]int;
  m_constantPicRateIdc		[MAX_TLAYER]int;
  m_avgPicRate				[MAX_TLAYER]int;
}

func NewTComBitRatePicRateInfo() *TComBitRatePicRateInfo{
	return &TComBitRatePicRateInfo{}
}

func (this *TComBitRatePicRateInfo)  GetBitRateInfoPresentFlag(i int)   bool  {
	return this.m_bitRateInfoPresentFlag[i];
}
func (this *TComBitRatePicRateInfo)  SetBitRateInfoPresentFlag(i int, x bool) {
	this.m_bitRateInfoPresentFlag[i] = x;
}

func (this *TComBitRatePicRateInfo)  GetPicRateInfoPresentFlag(i int) 	bool  {
	return this.m_picRateInfoPresentFlag[i];
}
func (this *TComBitRatePicRateInfo)  SetPicRateInfoPresentFlag(i int, x bool) {
	this.m_picRateInfoPresentFlag[i] = x;
}

func (this *TComBitRatePicRateInfo)  GetAvgBitRate(i int) int {
	return this.m_avgBitRate[i];
}
func (this *TComBitRatePicRateInfo)  SetAvgBitRate(i, x int)  {
	this.m_avgBitRate[i] = x;
}

func (this *TComBitRatePicRateInfo)  GetMaxBitRate(i int) int {
	return this.m_maxBitRate[i];
}
func (this *TComBitRatePicRateInfo)  SetMaxBitRate(i, x int)  {
	this.m_maxBitRate[i] = x;
}

func (this *TComBitRatePicRateInfo)  GetConstantPicRateIdc(i int) int {
	return this.m_constantPicRateIdc[i];
}
func (this *TComBitRatePicRateInfo)  SetConstantPicRateIdc(i, x int)  {
	this.m_constantPicRateIdc[i] = x;
}

func (this *TComBitRatePicRateInfo)  GetAvgPicRate(i int) int {
	return this.m_avgPicRate[i];
}
func (this *TComBitRatePicRateInfo)  SetAvgPicRate(i, x int)  {
	this.m_avgPicRate[i] = x;
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
   

//#if VPS_OPERATING_POINT
  	m_numHrdParameters	uint;
  	m_maxNuhReservedZeroLayerId	uint;
  	m_opLayerIdIncludedFlag	[MAX_VPS_NUM_HRD_PARAMETERS_ALLOWED_PLUS1][MAX_VPS_NUH_RESERVED_ZERO_LAYER_ID_PLUS1]bool;
//#endif    
    
    m_pcPTL                TComPTL
    
//#if SIGNAL_BITRATE_PICRATE_IN_VPS
  	m_bitRatePicRateInfo	TComBitRatePicRateInfo;
//#endif    
}

//public:

func NewTComVPS() *TComVPS {
    return &TComVPS{}
}

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
//#if VPS_OPERATING_POINT
func (this *TComVPS)  GetNumHrdParameters()  uint                               { 
	return this.m_numHrdParameters; 
}
func (this *TComVPS)  SetNumHrdParameters(v uint)                           { 
	this.m_numHrdParameters = v;    
}

func (this *TComVPS)  GetMaxNuhReservedZeroLayerId() uint                       { 
	return this.m_maxNuhReservedZeroLayerId; 
}
func (this *TComVPS)  SetMaxNuhReservedZeroLayerId(v uint)                  { 
	this.m_maxNuhReservedZeroLayerId = v;    
}

func (this *TComVPS)  GetOpLayerIdIncludedFlag( opIdx, id uint)  bool       { 
	return this.m_opLayerIdIncludedFlag[opIdx][id]; 
}
func (this *TComVPS)  SetOpLayerIdIncludedFlag( v bool,  opIdx,  id uint) { 
	this.m_opLayerIdIncludedFlag[opIdx][id] = v;    
}
//#endif
func (this *TComVPS) GetPTL() *TComPTL {
    return &this.m_pcPTL
}
//#if SIGNAL_BITRATE_PICRATE_IN_VPS
func (this *TComVPS)  GetBitratePicrateInfo() *TComBitRatePicRateInfo{ 
	return &this.m_bitRatePicRateInfo; 
}
//#endif

type HrdSubLayerInfo struct {
    fixedPicRateFlag      bool
    picDurationInTcMinus1 uint
    lowDelayHrdFlag       bool
    cpbCntMinus1          uint
    bitRateValueMinus1    [MAX_CPB_CNT][2]uint
    cpbSizeValue          [MAX_CPB_CNT][2]uint
    cbrFlag               [MAX_CPB_CNT][2]bool
}

type TComVUI struct {
    //private:
    m_aspectRatioInfoPresentFlag         bool
    m_aspectRatioIdc                     int
    m_sarWidth                           int
    m_sarHeight                          int
    m_overscanInfoPresentFlag            bool
    m_overscanAppropriateFlag            bool
    m_videoSignalTypePresentFlag         bool
    m_videoFormat                        int
    m_videoFullRangeFlag                 bool
    m_colourDescriptionPresentFlag       bool
    m_colourPrimaries                    int
    m_transferCharacteristics            int
    m_matrixCoefficients                 int
    m_chromaLocInfoPresentFlag           bool
    m_chromaSampleLocTypeTopField        int
    m_chromaSampleLocTypeBottomField     int
    m_neutralChromaIndicationFlag        bool
    m_fieldSeqFlag                       bool
    m_hrdParametersPresentFlag           bool
    m_bitstreamRestrictionFlag           bool
    m_tilesFixedStructureFlag            bool
    m_motionVectorsOverPicBoundariesFlag bool
    m_maxBytesPerPicDenom                int
    m_maxBitsPerMinCuDenom               int
    m_log2MaxMvLengthHorizontal          int
    m_log2MaxMvLengthVertical            int
    m_timingInfoPresentFlag              bool
    m_numUnitsInTick                     uint
    m_timeScale                          uint
    m_nalHrdParametersPresentFlag        bool
    m_vclHrdParametersPresentFlag        bool
    m_subPicCpbParamsPresentFlag         bool
    m_tickDivisorMinus2                  uint
    m_duCpbRemovalDelayLengthMinus1      uint
    m_bitRateScale                       uint
    m_cpbSizeScale                       uint
    m_initialCpbRemovalDelayLengthMinus1 uint
    m_cpbRemovalDelayLengthMinus1        uint
    m_dpbOutputDelayLengthMinus1         uint
    m_numDU                              uint
    m_HRD                                [MAX_TLAYER]HrdSubLayerInfo
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
        m_hrdParametersPresentFlag:           false,
        m_bitstreamRestrictionFlag:           false,
        m_tilesFixedStructureFlag:            false,
        m_motionVectorsOverPicBoundariesFlag: true,
        m_maxBytesPerPicDenom:                2,
        m_maxBitsPerMinCuDenom:               1,
        m_log2MaxMvLengthHorizontal:          15,
        m_log2MaxMvLengthVertical:            15,
        m_timingInfoPresentFlag:              false,
        m_numUnitsInTick:                     1001,
        m_timeScale:                          60000,
        m_nalHrdParametersPresentFlag:        false,
        m_vclHrdParametersPresentFlag:        false,
        m_subPicCpbParamsPresentFlag:         false,
        m_tickDivisorMinus2:                  0,
        m_duCpbRemovalDelayLengthMinus1:      0,
        m_bitRateScale:                       0,
        m_cpbSizeScale:                       0,
        m_initialCpbRemovalDelayLengthMinus1: 0,
        m_cpbRemovalDelayLengthMinus1:        0,
        m_dpbOutputDelayLengthMinus1:         0,
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

func (this *TComVUI) SetTimingInfoPresentFlag(flag bool) {
    this.m_timingInfoPresentFlag = flag
}
func (this *TComVUI) GetTimingInfoPresentFlag() bool {
    return this.m_timingInfoPresentFlag
}

func (this *TComVUI) SetNumUnitsInTick(value uint) {
    this.m_numUnitsInTick = value
}
func (this *TComVUI) GetNumUnitsInTick() uint {
    return this.m_numUnitsInTick
}

func (this *TComVUI) SetTimeScale(value uint) {
    this.m_timeScale = value
}
func (this *TComVUI) GetTimeScale() uint {
    return this.m_timeScale
}

func (this *TComVUI) SetNalHrdParametersPresentFlag(flag bool) {
    this.m_nalHrdParametersPresentFlag = flag
}
func (this *TComVUI) GetNalHrdParametersPresentFlag() bool {
    return this.m_nalHrdParametersPresentFlag
}

func (this *TComVUI) SetVclHrdParametersPresentFlag(flag bool) {
    this.m_vclHrdParametersPresentFlag = flag
}
func (this *TComVUI) GetVclHrdParametersPresentFlag() bool {
    return this.m_vclHrdParametersPresentFlag
}

func (this *TComVUI) SetSubPicCpbParamsPresentFlag(flag bool) {
    this.m_subPicCpbParamsPresentFlag = flag
}
func (this *TComVUI) GetSubPicCpbParamsPresentFlag() bool {
    return this.m_subPicCpbParamsPresentFlag
}

func (this *TComVUI) SetTickDivisorMinus2(value uint) {
    this.m_tickDivisorMinus2 = value
}
func (this *TComVUI) GetTickDivisorMinus2() uint {
    return this.m_tickDivisorMinus2
}

func (this *TComVUI) SetDuCpbRemovalDelayLengthMinus1(value uint) {
    this.m_duCpbRemovalDelayLengthMinus1 = value
}
func (this *TComVUI) GetDuCpbRemovalDelayLengthMinus1() uint {
    return this.m_duCpbRemovalDelayLengthMinus1
}

func (this *TComVUI) SetBitRateScale(value uint) {
    this.m_bitRateScale = value
}
func (this *TComVUI) GetBitRateScale() uint {
    return this.m_bitRateScale
}

func (this *TComVUI) SetCpbSizeScale(value uint) {
    this.m_cpbSizeScale = value
}
func (this *TComVUI) GetCpbSizeScale() uint {
    return this.m_cpbSizeScale
}

func (this *TComVUI) SetInitialCpbRemovalDelayLengthMinus1(value uint) {
    this.m_initialCpbRemovalDelayLengthMinus1 = value
}
func (this *TComVUI) GetInitialCpbRemovalDelayLengthMinus1() uint {
    return this.m_initialCpbRemovalDelayLengthMinus1
}

func (this *TComVUI) SetCpbRemovalDelayLengthMinus1(value uint) {
    this.m_cpbRemovalDelayLengthMinus1 = value
}
func (this *TComVUI) GetCpbRemovalDelayLengthMinus1() uint {
    return this.m_cpbRemovalDelayLengthMinus1
}

func (this *TComVUI) SetDpbOutputDelayLengthMinus1(value uint) {
    this.m_dpbOutputDelayLengthMinus1 = value
}
func (this *TComVUI) GetDpbOutputDelayLengthMinus1() uint {
    return this.m_dpbOutputDelayLengthMinus1
}

func (this *TComVUI) SetFixedPicRateFlag(layer int, flag bool) {
    this.m_HRD[layer].fixedPicRateFlag = flag
}
func (this *TComVUI) GetFixedPicRateFlag(layer int) bool {
    return this.m_HRD[layer].fixedPicRateFlag
}

func (this *TComVUI) SetPicDurationInTcMinus1(layer int, value uint) {
    this.m_HRD[layer].picDurationInTcMinus1 = value
}
func (this *TComVUI) GetPicDurationInTcMinus1(layer int) uint {
    return this.m_HRD[layer].picDurationInTcMinus1
}

func (this *TComVUI) SetLowDelayHrdFlag(layer int, flag bool) {
    this.m_HRD[layer].lowDelayHrdFlag = flag
}
func (this *TComVUI) GetLowDelayHrdFlag(layer int) bool {
    return this.m_HRD[layer].lowDelayHrdFlag
}

func (this *TComVUI) SetCpbCntMinus1(layer int, value uint) {
    this.m_HRD[layer].cpbCntMinus1 = value
}
func (this *TComVUI) GetCpbCntMinus1(layer int) uint {
    return this.m_HRD[layer].cpbCntMinus1
}

func (this *TComVUI) SetBitRateValueMinus1(layer, cpbcnt, nalOrVcl int, value uint) {
    this.m_HRD[layer].bitRateValueMinus1[cpbcnt][nalOrVcl] = value
}
func (this *TComVUI) GetBitRateValueMinus1(layer, cpbcnt, nalOrVcl int) uint {
    return this.m_HRD[layer].bitRateValueMinus1[cpbcnt][nalOrVcl]
}

func (this *TComVUI) SetCpbSizeValueMinus1(layer, cpbcnt, nalOrVcl int, value uint) {
    this.m_HRD[layer].cpbSizeValue[cpbcnt][nalOrVcl] = value
}
func (this *TComVUI) GetCpbSizeValueMinus1(layer, cpbcnt, nalOrVcl int) uint {
    return this.m_HRD[layer].cpbSizeValue[cpbcnt][nalOrVcl]
}

func (this *TComVUI) SetCbrFlag(layer, cpbcnt, nalOrVcl int, value bool) {
    this.m_HRD[layer].cbrFlag[cpbcnt][nalOrVcl] = value
}
func (this *TComVUI) GetCbrFlag(layer, cpbcnt, nalOrVcl int) bool {
    return this.m_HRD[layer].cbrFlag[cpbcnt][nalOrVcl]
}

func (this *TComVUI) SetNumDU(value uint) {
    this.m_numDU = value
}
func (this *TComVUI) GetNumDU() uint {
    return this.m_numDU
}

type CroppingWindow struct {
    //private:
    m_picCroppingFlag     bool
    m_picCropLeftOffset   int
    m_picCropRightOffset  int
    m_picCropTopOffset    int
    m_picCropBottomOffset int
}

func NewCroppingWindow() *CroppingWindow{
	return &CroppingWindow{}
}

func (this *CroppingWindow) GetPicCroppingFlag() bool {
    return this.m_picCroppingFlag
}
func (this *CroppingWindow) SetPicCroppingFlag(val bool) {
    this.m_picCroppingFlag = val
}
func (this *CroppingWindow) GetPicCropLeftOffset() int {
    return this.m_picCropLeftOffset
}
func (this *CroppingWindow) SetPicCropLeftOffset(val int) {
    this.m_picCropLeftOffset = val
}
func (this *CroppingWindow) GetPicCropRightOffset() int {
    return this.m_picCropRightOffset
}
func (this *CroppingWindow) SetPicCropRightOffset(val int) {
    this.m_picCropRightOffset = val
}
func (this *CroppingWindow) GetPicCropTopOffset() int {
    return this.m_picCropTopOffset
}
func (this *CroppingWindow) SetPicCropTopOffset(val int) {
    this.m_picCropTopOffset = val
}
func (this *CroppingWindow) GetPicCropBottomOffset() int {
    return this.m_picCropBottomOffset
}
func (this *CroppingWindow) SetPicCropBottomOffset(val int) {
    this.m_picCropBottomOffset = val
}

func (this *CroppingWindow) ResetCropping() {
    this.m_picCroppingFlag = false
    this.m_picCropLeftOffset = 0
    this.m_picCropRightOffset = 0
    this.m_picCropTopOffset = 0
    this.m_picCropBottomOffset = 0
}

func (this *CroppingWindow) SetPicCropping(cropLeft, cropRight, cropTop, cropBottom int) {
    this.m_picCroppingFlag = true
    this.m_picCropLeftOffset = cropLeft
    this.m_picCropRightOffset = cropRight
    this.m_picCropTopOffset = cropTop
    this.m_picCropBottomOffset = cropBottom
}


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

    m_picCroppingWindow *CroppingWindow

    m_uiMaxCUWidth         uint
    m_uiMaxCUHeight        uint
    m_uiMaxCUDepth         uint
    m_uiMinTrDepth         uint
    m_uiMaxTrDepth         uint
    m_RPSList              TComRPSList
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

    m_restrictedRefPicListsFlag    bool
    m_listsModificationPresentFlag bool

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
    //NTRA_SMOOTHING
    m_useStrongIntraSmoothing bool
    //

    m_vuiParametersPresentFlag bool
    m_vuiParameters            TComVUI

    m_cropUnitX [MAX_CHROMA_FORMAT_IDC + 1]int
    m_cropUnitY [MAX_CHROMA_FORMAT_IDC + 1]int
    m_pcPTL     TComPTL
}

//public:
func NewTComSPS() *TComSPS {
	sps := &TComSPS{};
	sps.m_picCroppingWindow = NewCroppingWindow();
	sps.m_cropUnitX[0]=1;
	sps.m_cropUnitX[1]=2;
	sps.m_cropUnitX[2]=2;
	sps.m_cropUnitX[3]=1;
	sps.m_cropUnitY[0]=1;
	sps.m_cropUnitY[1]=2;
	sps.m_cropUnitY[2]=1;
	sps.m_cropUnitY[3]=1;
	
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

func (this *TComSPS) GetCropUnitX(chromaFormatIdc int) int {
    //assert (chromaFormatIdc > 0 && chromaFormatIdc <= MAX_CHROMA_FORMAT_IDC); 
    return this.m_cropUnitX[chromaFormatIdc]
}
func (this *TComSPS) GetCropUnitY(chromaFormatIdc int) int {
    //assert (chromaFormatIdc > 0 && chromaFormatIdc <= MAX_CHROMA_FORMAT_IDC); 
    return this.m_cropUnitY[chromaFormatIdc]
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

func (this *TComSPS) GetPicCroppingWindow() *CroppingWindow {
    return this.m_picCroppingWindow
}
func (this *TComSPS) SetPicCroppingWindow(croppingWindow *CroppingWindow) {
    this.m_picCroppingWindow = croppingWindow
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
  	this.m_RPSList.Destroy();
  	this.m_RPSList.Create(numRPS);
}
func (this *TComSPS) GetRPSList() *TComRPSList {
    return &this.m_RPSList
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

func (this *TComSPS) GetRestrictedRefPicListsFlag() bool {
    return this.m_restrictedRefPicListsFlag
}
func (this *TComSPS) SetRestrictedRefPicListsFlag(b bool) {
    this.m_restrictedRefPicListsFlag = b
}
func (this *TComSPS) GetListsModificationPresentFlag() bool {
    return this.m_listsModificationPresentFlag
}
func (this *TComSPS) SetListsModificationPresentFlag(b bool) {
    this.m_listsModificationPresentFlag = b
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

//#if STRONG_INTRA_SMOOTHING
func (this *TComSPS) SetUseStrongIntraSmoothing(bVal bool) {
    this.m_useStrongIntraSmoothing = bVal
}
func (this *TComSPS) GetUseStrongIntraSmoothing() bool {
    return this.m_useStrongIntraSmoothing
}

//#endif

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
    return &TComRefPicListModification{}
}

func (this *TComRefPicListModification) Create() {
}
func (this *TComRefPicListModification) Destroy() {
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

    m_TransquantBypassEnableFlag   bool // Indicates presence of cu_transquant_bypass_flag in CUs.
    m_useTransformSkip             bool
    m_dependentSliceEnabledFlag    bool //!< Indicates the presence of dependent slices
    m_tilesEnabledFlag             bool //!< Indicates the presence of tiles
    m_entropyCodingSyncEnabledFlag bool //!< Indicates the presence of wavefronts
    //#if !REMOVE_ENTROPY_SLICES
    //  Bool        m_entropySliceEnabledFlag;       //!< Indicates the presence of entropy slices
    //#endif

    m_loopFilterAcrossTilesEnabledFlag bool
    m_uniformSpacingFlag               bool
    m_iNumColumnsMinus1                int
    m_puiColumnWidth                   []uint
    m_iNumRowsMinus1                   int
    m_puiRowHeight                     []uint

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
    
//#if HLS_MOVE_SPS_PICLIST_FLAGS
  	m_listsModificationPresentFlag		 bool;
//#endif /* HLS_MOVE_SPS_PICLIST_FLAGS */
  	m_log2ParallelMergeLevelMinus2        uint
//#if HLS_EXTRA_SLICE_HEADER_BITS
  	m_numExtraSliceHeaderBits			int;
//#endif /* HLS_EXTRA_SLICE_HEADER_BITS */    
}

//public:
func NewTComPPS() *TComPPS {
    return &TComPPS{}
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
func (this *TComPPS) GetDependentSliceEnabledFlag() bool {
    return this.m_dependentSliceEnabledFlag
}
func (this *TComPPS) SetDependentSliceEnabledFlag(val bool) {
    this.m_dependentSliceEnabledFlag = val
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

/*#if !REMOVE_ENTROPY_SLICES
  Bool    GetEntropySliceEnabledFlag() const               { return this.m_entropySliceEnabledFlag; }
  Void    SetEntropySliceEnabledFlag(Bool val)             { this.m_entropySliceEnabledFlag = val; }
#endif*/
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
func (this *TComPPS) SetColumnWidth(columnWidth []uint) {
    if this.m_uniformSpacingFlag == false && this.m_iNumColumnsMinus1 > 0 {
        this.m_puiColumnWidth = make([]uint, this.m_iNumColumnsMinus1)
        for i := 0; i < this.m_iNumColumnsMinus1; i++ {
            this.m_puiColumnWidth[i] = columnWidth[i]
        }
    }
}
func (this *TComPPS) GetColumnWidth(columnIdx uint) uint {
    return this.m_puiColumnWidth[columnIdx]
}
func (this *TComPPS) SetNumRowsMinus1(i int) {
    this.m_iNumRowsMinus1 = i
}
func (this *TComPPS) GetNumRowsMinus1() int {
    return this.m_iNumRowsMinus1
}
func (this *TComPPS) SetRowHeight(rowHeight []uint) {
    if this.m_uniformSpacingFlag == false && this.m_iNumRowsMinus1 > 0 {
        this.m_puiRowHeight = make([]uint, this.m_iNumRowsMinus1)
        for i := 0; i < this.m_iNumRowsMinus1; i++ {
            this.m_puiRowHeight[i] = rowHeight[i]
        }
    }
}
func (this *TComPPS) GetRowHeight(rowIdx uint) uint {
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
//#if HLS_MOVE_SPS_PICLIST_FLAGS
func (this *TComPPS)  GetListsModificationPresentFlag ()  bool   { 
	return this.m_listsModificationPresentFlag; 
}
func (this *TComPPS)  SetListsModificationPresentFlag ( b bool)  { 
	this.m_listsModificationPresentFlag = b;    
}
//#endif /* HLS_MOVE_SPS_PICLIST_FLAGS */
func (this *TComPPS) GetLog2ParallelMergeLevelMinus2() uint {
    return this.m_log2ParallelMergeLevelMinus2
}
func (this *TComPPS) SetLog2ParallelMergeLevelMinus2(mrgLevel uint) {
    this.m_log2ParallelMergeLevelMinus2 = mrgLevel
}
//#if HLS_EXTRA_SLICE_HEADER_BITS
func (this *TComPPS)  GetNumExtraSliceHeaderBits()  int  { 
	return this.m_numExtraSliceHeaderBits; 
}
func (this *TComPPS)  SetNumExtraSliceHeaderBits(i int) { 
	this.m_numExtraSliceHeaderBits = i; 
}
//#endif /* HLS_EXTRA_SLICE_HEADER_BITS */

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

type wpScalingParam struct {
    // Explicit weighted prediction parameters parsed in slice header,
    // or Implicit weighted prediction parameters (8 bits depth values).
    bPresentFlag      bool
    uiLog2WeightDenom uint
    iWeight           int
    iOffset           int

    // Weighted prediction scaling values built from above parameters (bitdepth scaled):
    w, o, offset, shift, round int
}

type wpACDCParam struct {
    iAC int64
    iDC int64
}

/// slice header class
type TComSlice struct {
    //private:
    //  Bitstream writing
    m_saoEnabledFlag            bool
    m_saoEnabledFlagChroma      bool ///< SAO Cb&Cr enabled flag
    m_iPPSId                    int  ///< picture parameter set ID
    m_PicOutputFlag             bool ///< pic_output_flag 
    m_iPOC                      int
    m_iLastIDR                  int
    m_prevPOC                   int
    m_pcRPS                     *TComReferencePictureSet
    m_LocalRPS                  TComReferencePictureSet
    m_iBDidx                    int
    m_iCombinationBDidx         int
    m_bCombineWithReferenceFlag bool
    m_RefPicListModification    TComRefPicListModification
    m_eNalUnitType              NalUnitType ///< Nal unit type for the slice
    m_eSliceType                SliceType
    m_iSliceQp                  int
    m_dependentSliceFlag        bool
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
    m_iSliceQpDelta   int
    m_iSliceQpDeltaCb int
    m_iSliceQpDeltaCr int
    m_apcRefPicList   [2][MAX_NUM_REF + 1]*TComPic
    m_aiRefPOCList    [2][MAX_NUM_REF + 1]int
    m_iDepth          int

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

    m_uiSliceMode                    uint
    m_uiSliceArgument                uint
    m_uiSliceCurStartCUAddr          uint
    m_uiSliceCurEndCUAddr            uint
    m_uiSliceIdx                     uint
    m_uiDependentSliceMode           uint
    m_uiDependentSliceArgument       uint
    m_uiDependentSliceCurStartCUAddr uint
    m_uiDependentSliceCurEndCUAddr   uint
    m_bNextSlice                     bool
    m_bNextDependentSlice            bool
    m_uiSliceBits                    uint
    m_uiDependentSliceCounter        uint
    m_bFinalized                     bool

    m_weightPredTable [2][MAX_NUM_REF][3]wpScalingParam // [REF_PIC_LIST_0 or REF_PIC_LIST_1][refIdx][0:Y, 1:U, 2:V]
    m_weightACDCParam [3]wpACDCParam                    // [0:Y, 1:U, 2:V]

    m_tileByteLocation     *list.List
    m_uiTileOffstForMultES uint

    m_puiSubstreamSizes *uint
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
    return &TComSlice{}
}

func (this *TComSlice) InitSlice() {
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
func (this *TComSlice) SetCombinationBDidx(iCombinationBDidx int) {
    this.m_iCombinationBDidx = iCombinationBDidx
}
func (this *TComSlice) GetCombinationBDidx() int {
    return this.m_iCombinationBDidx
}
func (this *TComSlice) SetCombineWithReferenceFlag(bCombineWithReferenceFlag bool) {
    this.m_bCombineWithReferenceFlag = bCombineWithReferenceFlag
}
func (this *TComSlice) GetCombineWithReferenceFlag() bool {
    return this.m_bCombineWithReferenceFlag
}
func (this *TComSlice) GetPrevPOC() int {
    return this.m_prevPOC
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
func (this *TComSlice) GetDependentSliceFlag() bool {
    return this.m_dependentSliceFlag
}
func (this *TComSlice) SetDependentSliceFlag(val bool) {
    this.m_dependentSliceFlag = val
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
}
func (this *TComSlice) GetCheckLDC() bool {
    return this.m_bCheckLDC
}
func (this *TComSlice) GetMvdL1ZeroFlag() bool {
    return this.m_bLMvdL1Zero
}
func (this *TComSlice) GetNumRpsCurrTempList() int {
    return 0
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
        this.m_prevPOC = i
    }
}
func (this *TComSlice) SetNalUnitType(e NalUnitType) {
    this.m_eNalUnitType = e
}
func (this *TComSlice) GetNalUnitType() NalUnitType {
    return this.m_eNalUnitType
}
func (this *TComSlice) GetRapPicFlag() bool {
    return true
}
func (this *TComSlice) GetIdrPicFlag() bool {
    return this.GetNalUnitType() == NAL_UNIT_CODED_SLICE_IDR || this.GetNalUnitType() == NAL_UNIT_CODED_SLICE_IDR_N_LP
}
func (this *TComSlice) CheckCRA(pReferencePictureSet *TComReferencePictureSet, pocCRA *int, prevRAPisBLA *bool, rcListPic *list.List) {
}
func (this *TComSlice) DecodingRefreshMarking(pocCRA *int, bRefreshPending *bool, rcListPic *list.List) {
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
}
func (this *TComSlice) SetRefPOCList() {
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
}

func (this *TComSlice) GetNoBackPredFlag() bool {
    return this.m_bNoBackPredFlag
}
func (this *TComSlice) SetNoBackPredFlag(b bool) {
    this.m_bNoBackPredFlag = b
}
func (this *TComSlice) GenerateCombinedList() {
}

func (this *TComSlice) GetTLayer() uint {
    return this.m_uiTLayer
}
func (this *TComSlice) SetTLayer(uiTLayer uint) {
    this.m_uiTLayer = uiTLayer
}

func (this *TComSlice) SetTLayerInfo(uiTLayer uint) {
}
func (this *TComSlice) DecodingMarking(rcListPic *list.List, iGOPSIze int, iMaxRefPicNum *int) {
}
func (this *TComSlice) ApplyReferencePictureSet(rcListPic *list.List, RPSList *TComReferencePictureSet) {
}
func (this *TComSlice) IsTemporalLayerSwitchingPoint(rcListPic *list.List, RPSList *TComReferencePictureSet) bool {
    return true
}
func (this *TComSlice) IsStepwiseTemporalLayerSwitchingPointCandidate(rcListPic *list.List, RPSList *TComReferencePictureSet) bool {
    return true
}
func (this *TComSlice) CheckThatAllRefPicsAreAvailable(rcListPic *list.List, pReferencePictureSet *TComReferencePictureSet, printErrors bool, pocRandomAccess int) int {
    return 0
}
func (this *TComSlice) CreateExplicitReferencePictureSetFromReference(rcListPic *list.List, pReferencePictureSet *TComReferencePictureSet) {
}

func (this *TComSlice) SetMaxNumMergeCand(val uint) {
    this.m_maxNumMergeCand = val
}
func (this *TComSlice) GetMaxNumMergeCand() uint {
    return this.m_maxNumMergeCand
}

func (this *TComSlice) SetSliceMode(uiMode uint) {
    this.m_uiSliceMode = uiMode
}
func (this *TComSlice) GetSliceMode() uint {
    return this.m_uiSliceMode
}
func (this *TComSlice) SetSliceArgument(uiArgument uint) {
    this.m_uiSliceArgument = uiArgument
}
func (this *TComSlice) GetSliceArgument() uint {
    return this.m_uiSliceArgument
}
func (this *TComSlice) SetSliceCurStartCUAddr(uiAddr uint) {
    this.m_uiSliceCurStartCUAddr = uiAddr
}
func (this *TComSlice) GetSliceCurStartCUAddr() uint {
    return this.m_uiSliceCurStartCUAddr
}
func (this *TComSlice) SetSliceCurEndCUAddr(uiAddr uint) {
    this.m_uiSliceCurEndCUAddr = uiAddr
}
func (this *TComSlice) GetSliceCurEndCUAddr() uint {
    return this.m_uiSliceCurEndCUAddr
}
func (this *TComSlice) SetSliceIdx(i uint) {
    this.m_uiSliceIdx = i
}
func (this *TComSlice) GetSliceIdx() uint {
    return this.m_uiSliceIdx
}
func (this *TComSlice) CopySliceInfo(pcSliceSrc *TComSlice) {
}
func (this *TComSlice) SetDependentSliceMode(uiMode uint) {
    this.m_uiDependentSliceMode = uiMode
}
func (this *TComSlice) GetDependentSliceMode() uint {
    return this.m_uiDependentSliceMode
}
func (this *TComSlice) SetDependentSliceArgument(uiArgument uint) {
    this.m_uiDependentSliceArgument = uiArgument
}
func (this *TComSlice) GetDependentSliceArgument() uint {
    return this.m_uiDependentSliceArgument
}
func (this *TComSlice) SetDependentSliceCurStartCUAddr(uiAddr uint) {
    this.m_uiDependentSliceCurStartCUAddr = uiAddr
}
func (this *TComSlice) GetDependentSliceCurStartCUAddr() uint {
    return this.m_uiDependentSliceCurStartCUAddr
}
func (this *TComSlice) SetDependentSliceCurEndCUAddr(uiAddr uint) {
    this.m_uiDependentSliceCurEndCUAddr = uiAddr
}
func (this *TComSlice) GetDependentSliceCurEndCUAddr() uint {
    return this.m_uiDependentSliceCurEndCUAddr
}
func (this *TComSlice) SetNextSlice(b bool) {
    this.m_bNextSlice = b
}
func (this *TComSlice) IsNextSlice() bool {
    return this.m_bNextSlice
}
func (this *TComSlice) SetNextDependentSlice(b bool) {
    this.m_bNextDependentSlice = b
}
func (this *TComSlice) IsNextDependentSlice() bool {
    return this.m_bNextDependentSlice
}
func (this *TComSlice) SetSliceBits(uiVal uint) {
    this.m_uiSliceBits = uiVal
}
func (this *TComSlice) GetSliceBits() uint {
    return this.m_uiSliceBits
}
func (this *TComSlice) SetDependentSliceCounter(uiVal uint) {
    this.m_uiDependentSliceCounter = uiVal
}
func (this *TComSlice) GetDependentSliceCounter() uint {
    return this.m_uiDependentSliceCounter
}
func (this *TComSlice) SetFinalized(uiVal bool) {
    this.m_bFinalized = uiVal
}
func (this *TComSlice) GetFinalized() bool {
    return this.m_bFinalized
}
func (this *TComSlice) SetWpScaling(wp [2][MAX_NUM_REF][3]wpScalingParam) {
    //memcpy(this.m_weightPredTable, wp, sizeof(wpScalingParam)*2*MAX_NUM_REF*3); 
    this.m_weightPredTable = wp
}
func (this *TComSlice) GetWpScaling(e RefPicList, iRefIdx int, wp *wpScalingParam) {
}

func (this *TComSlice) ResetWpScaling(wp [2][MAX_NUM_REF][3]wpScalingParam) {
}
func (this *TComSlice) InitWpScaling1(wp [2][MAX_NUM_REF][3]wpScalingParam) {
}
func (this *TComSlice) InitWpScaling() {
}
func (this *TComSlice) ApplyWP() bool {
    return ((this.m_eSliceType == P_SLICE && this.m_pcPPS.GetUseWP()) || (this.m_eSliceType == B_SLICE && this.m_pcPPS.GetWPBiPred()))
}

func (this *TComSlice) SetWpAcDcParam(wp [3]wpACDCParam) {
    //memcpy(this.m_weightACDCParam, wp, sizeof(wpACDCParam)*3); 
    this.m_weightACDCParam = wp
}
func (this *TComSlice) GetWpAcDcParam(wp *wpACDCParam) {
}
func (this *TComSlice) InitWpAcDcParam() {
}

func (this *TComSlice) SetTileLocationCount(cnt uint) {
    //	return this.m_tileByteLocation.Resize(cnt);    
}
func (this *TComSlice) GetTileLocationCount() uint {
    return uint(this.m_tileByteLocation.Len())
}
func (this *TComSlice) SetTileLocation(idx int, location uint) {
    //assert (idx<this.m_tileByteLocation.size());
    //this.m_tileByteLocation[idx] = location;       
}
func (this *TComSlice) AddTileLocation(location uint) {
    this.m_tileByteLocation.PushBack(location)
}
func (this *TComSlice) GetTileLocation(idx int) uint {
    return 0 //this.m_tileByteLocation[idx];          
}

func (this *TComSlice) SetTileOffstForMultES(uiOffset uint) {
    this.m_uiTileOffstForMultES = uiOffset
}
func (this *TComSlice) GetTileOffstForMultES() uint {
    return this.m_uiTileOffstForMultES
}
func (this *TComSlice) AllocSubstreamSizes(uiNumSubstreams uint) {
}
func (this *TComSlice) GetSubstreamSizes() *uint {
    return this.m_puiSubstreamSizes
}
func (this *TComSlice) SetScalingList(scalingList *TComScalingList) {
    this.m_scalingList = scalingList
}
func (this *TComSlice) GetScalingList() *TComScalingList {
    return this.m_scalingList
}
func (this *TComSlice) SetDefaultScalingList() {
}
func (this *TComSlice) CheckDefaultScalingList() bool {
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
    return nil
}
func (this *TComSlice) xGetLongTermRefPic(rcListPic *list.List, poc int) *TComPic {
    return nil
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
}

//public:
func NewParameterSetManager() *ParameterSetManager {
    return &ParameterSetManager{make(map[int]*TComVPS), make(map[int]*TComSPS), make(map[int]*TComPPS)}
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
//func (this *ParameterSetManager)  TComPPS* getFirstPPS()      { return m_ppsMap.getFirstPS(); };
