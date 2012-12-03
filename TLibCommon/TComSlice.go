package TLibCommon

import (

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
type TComReferencePictureSet struct{
//private:
  m_numberOfPictures			int;
  m_numberOfNegativePictures	int;
  m_numberOfPositivePictures	int;
  m_numberOfLongtermPictures	int;
  m_deltaPOC	[MAX_NUM_REF_PICS]int;
  m_POC		[MAX_NUM_REF_PICS]int;
  m_used		[MAX_NUM_REF_PICS]bool;
  m_interRPSPrediction			bool;
  m_deltaRIdxMinus1			int;   
  m_deltaRPS					int; 
  m_numRefIdc					int; 
  m_refIdc	[MAX_NUM_REF_PICS+1]int;
  m_bCheckLTMSB	[MAX_NUM_REF_PICS]bool;
  m_pocLSBLT		[MAX_NUM_REF_PICS]int;
  m_deltaPOCMSBCycleLT	[MAX_NUM_REF_PICS]int;
  m_deltaPocMSBPresentFlag	[MAX_NUM_REF_PICS]bool;
}

//public:
func  NewTComReferencePictureSet() *TComReferencePictureSet{
	return &TComReferencePictureSet{}
}

func  (this *TComReferencePictureSet) GetPocLSBLT(i int) int { 
	return this.m_pocLSBLT[i]; 
}
func  (this *TComReferencePictureSet) SetPocLSBLT(i, x int) { 
	this.m_pocLSBLT[i] = x; 
}
func  (this *TComReferencePictureSet) GetDeltaPocMSBCycleLT(i int) int { 
	return this.m_deltaPOCMSBCycleLT[i]; 
}
func  (this *TComReferencePictureSet) SetDeltaPocMSBCycleLT(i, x int) { 
	this.m_deltaPOCMSBCycleLT[i] = x; 
}
func  (this *TComReferencePictureSet) GetDeltaPocMSBPresentFlag(i int)  bool       { 
	return this.m_deltaPocMSBPresentFlag[i]; 
}
func  (this *TComReferencePictureSet) SetDeltaPocMSBPresentFlag(i int, x bool) { 
	this.m_deltaPocMSBPresentFlag[i] = x;    
}
func  (this *TComReferencePictureSet) SetUsed(bufferNum int, used bool){
}
func  (this *TComReferencePictureSet) SetDeltaPOC(bufferNum, deltaPOC int){
}
func  (this *TComReferencePictureSet) SetPOC(bufferNum, deltaPOC int){
}
func  (this *TComReferencePictureSet) SetNumberOfPictures(numberOfPictures int){
}
func  (this *TComReferencePictureSet) SetCheckLTMSBPresent2(bufferNum int, b int){
}
func  (this *TComReferencePictureSet) SetCheckLTMSBPresent1(bufferNum int) bool{
	return true
}

func  (this *TComReferencePictureSet)   GetUsed(bufferNum int)int{
	return 0
}
func  (this *TComReferencePictureSet)   GetDeltaPOC(bufferNum int)int{
	return 0
}
func  (this *TComReferencePictureSet)   GetPOC(bufferNum int) int {
	return 0
}
func  (this *TComReferencePictureSet)   GetNumberOfPictures() int {
	return 0
}

func  (this *TComReferencePictureSet)   SetNumberOfNegativePictures(number int)  { 
	this.m_numberOfNegativePictures = number; 
}
func  (this *TComReferencePictureSet)   GetNumberOfNegativePictures()   int         { 
	return this.m_numberOfNegativePictures; 
}
func  (this *TComReferencePictureSet)   SetNumberOfPositivePictures(number int)  { 
	this.m_numberOfPositivePictures = number; 
}
func  (this *TComReferencePictureSet)   GetNumberOfPositivePictures()   int         { 
	return this.m_numberOfPositivePictures; 
}
func  (this *TComReferencePictureSet)   SetNumberOfLongtermPictures(number int)  { 
	this.m_numberOfLongtermPictures = number; 
}
func  (this *TComReferencePictureSet)   GetNumberOfLongtermPictures()   int         { 
	return this.m_numberOfLongtermPictures; 
}

func  (this *TComReferencePictureSet)   SetInterRPSPrediction(flag bool)         { 
	this.m_interRPSPrediction = flag; 
}
func  (this *TComReferencePictureSet)   GetInterRPSPrediction()     bool             { 
	return this.m_interRPSPrediction; 
}
func  (this *TComReferencePictureSet)   SetDeltaRIdxMinus1(x int)                { 
	this.m_deltaRIdxMinus1 = x; 
}
func  (this *TComReferencePictureSet)   GetDeltaRIdxMinus1()        int              { 
	return this.m_deltaRIdxMinus1; 
}
func  (this *TComReferencePictureSet)   SetDeltaRPS(x int)                       { 
	this.m_deltaRPS = x; 
}
func  (this *TComReferencePictureSet)   GetDeltaRPS()               int              { 
	return this.m_deltaRPS; 
}
func  (this *TComReferencePictureSet)   SetNumRefIdc(x int)                      { 
	this.m_numRefIdc = x; 
}
func  (this *TComReferencePictureSet)   GetNumRefIdc()              int             { 
	return this.m_numRefIdc; 
}

func  (this *TComReferencePictureSet)   SetRefIdc(bufferNum, refIdc int) {
	this.m_refIdc[bufferNum] = refIdc;
}
func  (this *TComReferencePictureSet)   GetRefIdc(bufferNum int) int{
	return this.m_refIdc[bufferNum];	
}

func  (this *TComReferencePictureSet)   SortDeltaPOC(){
}
func  (this *TComReferencePictureSet)   PrintDeltaPOC(){
}
//};

/// Reference Picture Set set class
type TComRPSList struct{
//private:
  m_numberOfReferencePictureSets int;
  m_referencePictureSets []TComReferencePictureSet;
}
  
//public:
func NewTComRPSList() *TComRPSList{
	return &TComRPSList{}
}
 
func (this *TComRPSList) Create  (numberOfReferencePictureSets int){
  this.m_numberOfReferencePictureSets = numberOfReferencePictureSets;
  this.m_referencePictureSets = make([]TComReferencePictureSet, numberOfReferencePictureSets);
}
func (this *TComRPSList) Destroy (){
}


func (this *TComRPSList) GetReferencePictureSet(referencePictureSetNum int) *  TComReferencePictureSet{
	return &this.m_referencePictureSets[referencePictureSetNum];
}
func (this *TComRPSList) GetNumberOfReferencePictureSets() int {
	return this.m_numberOfReferencePictureSets;
}
func (this *TComRPSList) SetNumberOfReferencePictureSets(numberOfReferencePictureSets int){
	this.m_numberOfReferencePictureSets = numberOfReferencePictureSets;
}

/// SCALING_LIST class
type TComScalingList struct{
  m_scalingListDC               [SCALING_LIST_SIZE_NUM][SCALING_LIST_NUM]int; //!< the DC value of the matrix coefficient for 16x16
  m_useDefaultScalingMatrixFlag [SCALING_LIST_SIZE_NUM][SCALING_LIST_NUM]bool; //!< UseDefaultScalingMatrixFlag
  m_refMatrixId                 [SCALING_LIST_SIZE_NUM][SCALING_LIST_NUM]uint; //!< RefMatrixID
  m_scalingListPresentFlag												 bool;                                                //!< flag for using default matrix
  m_predMatrixId                [SCALING_LIST_SIZE_NUM][SCALING_LIST_NUM]uint; //!< reference list index
  m_scalingListCoef           [][SCALING_LIST_SIZE_NUM][SCALING_LIST_NUM]int; //!< quantization matrix
  m_useTransformSkip													 bool;  
}

//public:
func NewTComScalingList() *TComScalingList{
	return &TComScalingList{};
}

func (this *TComScalingList) SetScalingListPresentFlag (b bool) { 
	this.m_scalingListPresentFlag = b;    
}
func (this *TComScalingList) GetScalingListPresentFlag ()  bool { 
	return this.m_scalingListPresentFlag; 
}
func (this *TComScalingList) GetUseTransformSkip    () bool  { 
	return this.m_useTransformSkip; 
}      
func (this *TComScalingList) SetUseTransformSkip    (b bool)  { 
	this.m_useTransformSkip = b;    
}
func (this *TComScalingList) GetScalingListAddress          (sizeId, listId uint)  []int         { 
	return this.m_scalingListCoef[sizeId][listId][:]; 
} //!< get matrix coefficient
func (this *TComScalingList) CheckPredMode                  (sizeId, listId uint) bool{
	return true;
}
func (this *TComScalingList) SetRefMatrixId                 (sizeId, listId, u uint)    { 
	this.m_refMatrixId[sizeId][listId] = u;    
}     //!< set reference matrix ID
func (this *TComScalingList) GetRefMatrixId                 (sizeId, listId uint)   uint        { 
	return this.m_refMatrixId[sizeId][listId]; 
}     //!< get reference matrix ID
func (this *TComScalingList) GetScalingListDefaultAddress   (sizeId, listId uint) []int{
  var src []int;
  switch sizeId {
    case SCALING_LIST_4x4:
//#if FLAT_4x4_DSL
      src = g_quantTSDefault4x4[:];
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
    	if listId<3{
    		src = g_quantIntraDefault8x8[:]
    	}else{
    		src = g_quantInterDefault8x8[:]
    	}
      //src = (listId<3) ? g_quantIntraDefault8x8 : g_quantInterDefault8x8;
      //break;
    case SCALING_LIST_16x16:
    	if listId<3{
    		src = g_quantIntraDefault8x8[:]
    	}else{
    		src = g_quantInterDefault8x8[:]
    	}
      //src = (listId<3) ? g_quantIntraDefault8x8 : g_quantInterDefault8x8;
      //break;
    case SCALING_LIST_32x32:
      	if listId<1{
    		src = g_quantIntraDefault8x8[:]
    	}else{
    		src = g_quantInterDefault8x8[:]
    	}
      //src = (listId<1) ? g_quantIntraDefault8x8 : g_quantInterDefault8x8;
      //break;
    default:
    //  assert(0);
      src = nil;//NULL;
      //break;
  }
  return src;
}                                                        //!< get default matrix coefficient
func (this *TComScalingList) ProcessDefaultMarix            (sizeId, listId uint){
}
func (this *TComScalingList) SetScalingListDC               (sizeId, listId, u uint)   { 
	this.m_scalingListDC[sizeId][listId] = int(u); 
}      //!< set DC value

func (this *TComScalingList) GetScalingListDC               (sizeId, listId uint)  int         { 
	return this.m_scalingListDC[sizeId][listId]; 
}   //!< get DC value
func (this *TComScalingList) CheckDcOfMatrix                (){
}
func (this *TComScalingList) ProcessRefMatrix               (sizeId, listId, refListId uint){
}
func (this *TComScalingList) XParseScalingList              (pchFile string) bool {
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
  return false;
}

//private:
func (this *TComScalingList) init                    (){
}
func (this *TComScalingList) destroy                 (){
}
                                                    //!< transform skipping flag for setting default scaling matrix for 4x4


type ProfileTierLevel struct{
  m_profileSpace	int;
  m_tierFlag		bool;
  m_profileIdc		int;
  m_profileCompatibilityFlag	[32]bool;
  m_levelIdc		int;
}
//public:
func NewProfileTierLevel() *ProfileTierLevel{
	return &ProfileTierLevel{};
}

func (this *ProfileTierLevel)  GetProfileSpace() int   { 
	return this.m_profileSpace; 
}
func (this *ProfileTierLevel)  SetProfileSpace(x int)    { 
	this.m_profileSpace = x; 
}

func (this *ProfileTierLevel)  GetTierFlag()     bool   { 
	return this.m_tierFlag; 
}
func (this *ProfileTierLevel)  SetTierFlag(x bool)       { 
	this.m_tierFlag = x; 
}

func (this *ProfileTierLevel)  GetProfileIdc()   int   { 
	return this.m_profileIdc; 
}
func (this *ProfileTierLevel)  SetProfileIdc(x int)      { 
	this.m_profileIdc = x; 
}

func (this *ProfileTierLevel)  GetProfileCompatibilityFlag(i int) bool    { 
	return this.m_profileCompatibilityFlag[i]; 
}
func (this *ProfileTierLevel)  SetProfileCompatibilityFlag(i int, x bool)  { 
	this.m_profileCompatibilityFlag[i] = x; 
}

func (this *ProfileTierLevel)  GetLevelIdc()   int   { 
	return this.m_levelIdc; 
}
func (this *ProfileTierLevel)  SetLevelIdc(x int)      { 
	this.m_levelIdc = x; 
}



type TComPTL struct{
  m_generalPTL ProfileTierLevel;
  m_subLayerPTL	[6]ProfileTierLevel ;      // max. value of max_sub_layers_minus1 is 6
  m_subLayerProfilePresentFlag	[6]bool;
  m_subLayerLevelPresentFlag	[6]bool;
}

//public:
func NewTComPTL() *TComPTL{
	return &TComPTL{}
}
func (this *TComPTL)  GetSubLayerProfilePresentFlag(i int) bool { 
	return this.m_subLayerProfilePresentFlag[i]; 
}
func (this *TComPTL)  SetSubLayerProfilePresentFlag(i int, x bool) { 
	this.m_subLayerProfilePresentFlag[i] = x; 
}
  
func (this *TComPTL)  GetSubLayerLevelPresentFlag(i int) bool { 
	return this.m_subLayerLevelPresentFlag[i]; 
}
func (this *TComPTL)  SetSubLayerLevelPresentFlag(i int, x bool) { 
	this.m_subLayerLevelPresentFlag[i] = x; 
}

func (this *TComPTL)  GetGeneralPTL() *ProfileTierLevel  { 
	return &this.m_generalPTL; 
}
func (this *TComPTL)  GetSubLayerPTL(i int) *ProfileTierLevel  { 
	return &this.m_subLayerPTL[i]; 
}
/// VPS class

type TComVPS struct{
//private:
  m_VPSId			int;
  m_uiMaxTLayers	uint;
  m_uiMaxLayers		uint;
  m_bTemporalIdNestingFlag	bool;
  
  m_numReorderPics			[MAX_TLAYER]uint;
  m_uiMaxDecPicBuffering	[MAX_TLAYER]uint; 
  m_uiMaxLatencyIncrease	[MAX_TLAYER]uint;
  m_pcPTL 			TComPTL;
}
//public:

func NewTComVPS() *TComVPS{
	return &TComVPS{};
}
  
func (this *TComVPS)  GetVPSId       ()     int              { 
	return this.m_VPSId;          
}
func (this *TComVPS)  SetVPSId       (i int)              { 
	this.m_VPSId = i;             
}

func (this *TComVPS)  GetMaxTLayers  ()     uint              { 
	return this.m_uiMaxTLayers;   
}
func (this *TComVPS)  SetMaxTLayers  (t uint)             { 
	this.m_uiMaxTLayers = t; 
}
  
func (this *TComVPS)  GetMaxLayers   ()     uint              { 
	return this.m_uiMaxLayers;   
}
func (this *TComVPS)  SetMaxLayers   (l uint)             { 
	this.m_uiMaxLayers = l; 
}

func (this *TComVPS)  GetTemporalNestingFlag   () bool         { 
	return this.m_bTemporalIdNestingFlag;   
}
func (this *TComVPS)  SetTemporalNestingFlag   (t bool)   { 
	this.m_bTemporalIdNestingFlag = t; 
}
  
func (this *TComVPS)  SetNumReorderPics(v, tLayer uint)                { 
	this.m_numReorderPics[tLayer] = v;    
}
func (this *TComVPS)  GetNumReorderPics(tLayer uint)  uint                      { 
	return this.m_numReorderPics[tLayer]; 
}
  
func (this *TComVPS)  SetMaxDecPicBuffering(v, tLayer uint)            { 
	this.m_uiMaxDecPicBuffering[tLayer] = v;    
}
func (this *TComVPS)  GetMaxDecPicBuffering(tLayer uint)  uint                  { 
	return this.m_uiMaxDecPicBuffering[tLayer]; 
}
  
func (this *TComVPS)  SetMaxLatencyIncrease(v, tLayer uint)            { 
	this.m_uiMaxLatencyIncrease[tLayer] = v;    
}
func (this *TComVPS)  GetMaxLatencyIncrease(tLayer uint)  uint                  { 
	return this.m_uiMaxLatencyIncrease[tLayer]; 
}
func (this *TComVPS)  GetPTL() *TComPTL { 
	return &this.m_pcPTL; 
}


type HrdSubLayerInfo struct{
  fixedPicRateFlag			bool;
  picDurationInTcMinus1		uint;
  lowDelayHrdFlag			bool;
  cpbCntMinus1				uint;
  bitRateValueMinus1	[MAX_CPB_CNT][2]uint;
  cpbSizeValue     	 	[MAX_CPB_CNT][2]uint;
  cbrFlag           	[MAX_CPB_CNT][2]bool;
};

type TComVUI struct{
//private:
  m_aspectRatioInfoPresentFlag		bool;
  m_aspectRatioIdc					int;
  m_sarWidth						int;
  m_sarHeight						int;
  m_overscanInfoPresentFlag		bool;
  m_overscanAppropriateFlag		bool;
  m_videoSignalTypePresentFlag		bool;
  m_videoFormat					int;
  m_videoFullRangeFlag				bool;
  m_colourDescriptionPresentFlag	bool;
  m_colourPrimaries				int;
  m_transferCharacteristics		int;
  m_matrixCoefficients				int;
  m_chromaLocInfoPresentFlag		bool;
  m_chromaSampleLocTypeTopField	int;
  m_chromaSampleLocTypeBottomField	int;
  m_neutralChromaIndicationFlag	bool;
  m_fieldSeqFlag					bool;
  m_hrdParametersPresentFlag		bool;
  m_bitstreamRestrictionFlag		bool;
  m_tilesFixedStructureFlag		bool;
  m_motionVectorsOverPicBoundariesFlag	bool;
  m_maxBytesPerPicDenom			int;
  m_maxBitsPerMinCuDenom			int;
  m_log2MaxMvLengthHorizontal		int;
  m_log2MaxMvLengthVertical		int;
  m_timingInfoPresentFlag			bool;
  m_numUnitsInTick					uint;
  m_timeScale						uint;
  m_nalHrdParametersPresentFlag	bool;
  m_vclHrdParametersPresentFlag	bool;
  m_subPicCpbParamsPresentFlag		bool;
  m_tickDivisorMinus2				uint;
  m_duCpbRemovalDelayLengthMinus1	uint;
  m_bitRateScale					uint;
  m_cpbSizeScale					uint;
  m_initialCpbRemovalDelayLengthMinus1	uint;
  m_cpbRemovalDelayLengthMinus1	uint;
  m_dpbOutputDelayLengthMinus1		uint;
  m_numDU							uint;
  m_HRD		[MAX_TLAYER]HrdSubLayerInfo;
}
//public:
func NewTComVUI() *TComVUI{
	return &TComVUI{
    m_aspectRatioInfoPresentFlag:false,				
    m_aspectRatioIdc:0,
    m_sarWidth:0,
    m_sarHeight:0,
    m_overscanInfoPresentFlag:false,
    m_overscanAppropriateFlag:false,
    m_videoSignalTypePresentFlag:false,
    m_videoFormat:5,
    m_videoFullRangeFlag:false,
    m_colourDescriptionPresentFlag:false,
    m_colourPrimaries:2,
    m_transferCharacteristics:2,
    m_matrixCoefficients:2,
    m_chromaLocInfoPresentFlag:false,
    m_chromaSampleLocTypeTopField:0,
    m_chromaSampleLocTypeBottomField:0,
    m_neutralChromaIndicationFlag:false,
    m_fieldSeqFlag:false,
    m_hrdParametersPresentFlag:false,
    m_bitstreamRestrictionFlag:false,
    m_tilesFixedStructureFlag:false,
    m_motionVectorsOverPicBoundariesFlag:true,
    m_maxBytesPerPicDenom:2,
    m_maxBitsPerMinCuDenom:1,
    m_log2MaxMvLengthHorizontal:15,
    m_log2MaxMvLengthVertical:15,
    m_timingInfoPresentFlag:false,
    m_numUnitsInTick:1001,
    m_timeScale:60000,
    m_nalHrdParametersPresentFlag:false,
    m_vclHrdParametersPresentFlag:false,
    m_subPicCpbParamsPresentFlag:false,
    m_tickDivisorMinus2:0,
    m_duCpbRemovalDelayLengthMinus1:0,
    m_bitRateScale:0,
    m_cpbSizeScale:0,
    m_initialCpbRemovalDelayLengthMinus1:0,
    m_cpbRemovalDelayLengthMinus1:0,
    m_dpbOutputDelayLengthMinus1:0,
  }
}

func (this *TComVUI)   GetAspectRatioInfoPresentFlag() bool  { 
	return this.m_aspectRatioInfoPresentFlag; 
}
func (this *TComVUI)   SetAspectRatioInfoPresentFlag(i bool) { 
	this.m_aspectRatioInfoPresentFlag = i; 
}

func (this *TComVUI)   GetAspectRatioIdc() int  { 
	return this.m_aspectRatioIdc; 
}
func (this *TComVUI)   SetAspectRatioIdc(i int) { 
	this.m_aspectRatioIdc = i; 
}

func (this *TComVUI)   GetSarWidth() int  { 
	return this.m_sarWidth; 
}
func (this *TComVUI)   SetSarWidth(i int) { 
	this.m_sarWidth = i; 
}

func (this *TComVUI)   GetSarHeight() int  { 
	return this.m_sarHeight; 
}
func (this *TComVUI)   SetSarHeight(i int) { 
	this.m_sarHeight = i; 
}

func (this *TComVUI)   GetOverscanInfoPresentFlag() bool  { 
	return this.m_overscanInfoPresentFlag; 
}
func (this *TComVUI)   SetOverscanInfoPresentFlag(i bool) { 
	this.m_overscanInfoPresentFlag = i; 
}

func (this *TComVUI)   GetOverscanAppropriateFlag() bool  { 
	return this.m_overscanAppropriateFlag; 
}
func (this *TComVUI)   SetOverscanAppropriateFlag(i bool) { 
	this.m_overscanAppropriateFlag = i; 
}

func (this *TComVUI)   GetVideoSignalTypePresentFlag() bool  { 
	return this.m_videoSignalTypePresentFlag; 
}
func (this *TComVUI)   SetVideoSignalTypePresentFlag(i bool) { 
	this.m_videoSignalTypePresentFlag = i; 
}

func (this *TComVUI)   GetVideoFormat() int  { 
	return this.m_videoFormat; 
}
func (this *TComVUI)   SetVideoFormat(i int) { 
	this.m_videoFormat = i; 
}

func (this *TComVUI)   GetVideoFullRangeFlag() bool  { 
	return this.m_videoFullRangeFlag; 
}
func (this *TComVUI)   SetVideoFullRangeFlag(i bool) { 
	this.m_videoFullRangeFlag = i; 
}

func (this *TComVUI)   GetColourDescriptionPresentFlag() bool  { 
	return this.m_colourDescriptionPresentFlag; 
}
func (this *TComVUI)   SetColourDescriptionPresentFlag(i bool) { 
	this.m_colourDescriptionPresentFlag = i; 
}

func (this *TComVUI)   GetColourPrimaries() int  { 
	return this.m_colourPrimaries; 
}
func (this *TComVUI)   SetColourPrimaries(i int) { 
	this.m_colourPrimaries = i; 
}

func (this *TComVUI)   GetTransferCharacteristics() int  { 
	return this.m_transferCharacteristics; 
}
func (this *TComVUI)   SetTransferCharacteristics(i int) { 
	this.m_transferCharacteristics = i; 
}

func (this *TComVUI)   GetMatrixCoefficients() int  { 
	return this.m_matrixCoefficients; 
}
func (this *TComVUI)   SetMatrixCoefficients(i int) { 
	this.m_matrixCoefficients = i; 
}

func (this *TComVUI)   GetChromaLocInfoPresentFlag() bool  { 
	return this.m_chromaLocInfoPresentFlag; 
}
func (this *TComVUI)   SetChromaLocInfoPresentFlag(i bool) { 
	this.m_chromaLocInfoPresentFlag = i; 
}

func (this *TComVUI)   GetChromaSampleLocTypeTopField() int  { 
	return this.m_chromaSampleLocTypeTopField; 
}
func (this *TComVUI)   SetChromaSampleLocTypeTopField(i int) { 
	this.m_chromaSampleLocTypeTopField = i; 
}

func (this *TComVUI)   GetChromaSampleLocTypeBottomField() int  { 
	return this.m_chromaSampleLocTypeBottomField; 
}
func (this *TComVUI)   SetChromaSampleLocTypeBottomField(i int) { 
	this.m_chromaSampleLocTypeBottomField = i; 
}

func (this *TComVUI)   GetNeutralChromaIndicationFlag() bool  { 
	return this.m_neutralChromaIndicationFlag; 
}
func (this *TComVUI)   SetNeutralChromaIndicationFlag(i bool) { 
	this.m_neutralChromaIndicationFlag = i; 
}

func (this *TComVUI)   GetFieldSeqFlag() bool  { 
	return this.m_fieldSeqFlag; 
}
func (this *TComVUI)   SetFieldSeqFlag(i bool) { 
	this.m_fieldSeqFlag = i; 
}

func (this *TComVUI)   GetHrdParametersPresentFlag() bool  { 
	return this.m_hrdParametersPresentFlag; 
}
func (this *TComVUI)   SetHrdParametersPresentFlag(i bool) { 
	this.m_hrdParametersPresentFlag = i; 
}

func (this *TComVUI)   GetBitstreamRestrictionFlag() bool  { 
	return this.m_bitstreamRestrictionFlag; 
}
func (this *TComVUI)   SetBitstreamRestrictionFlag(i bool) { 
	this.m_bitstreamRestrictionFlag = i; 
}

func (this *TComVUI)   GetTilesFixedStructureFlag() bool  { 
	return this.m_tilesFixedStructureFlag; 
}
func (this *TComVUI)   SetTilesFixedStructureFlag(i bool) { 
	this.m_tilesFixedStructureFlag = i; 
}

func (this *TComVUI)   GetMotionVectorsOverPicBoundariesFlag() bool  { 
	return this.m_motionVectorsOverPicBoundariesFlag; 
}
func (this *TComVUI)   SetMotionVectorsOverPicBoundariesFlag(i bool) { 
	this.m_motionVectorsOverPicBoundariesFlag = i; 
}

func (this *TComVUI)   GetMaxBytesPerPicDenom() int  { 
	return this.m_maxBytesPerPicDenom; 
}
func (this *TComVUI)   SetMaxBytesPerPicDenom(i int) { 
	this.m_maxBytesPerPicDenom = i; 
}

func (this *TComVUI)   GetMaxBitsPerMinCuDenom() int  { 
	return this.m_maxBitsPerMinCuDenom; 
}
func (this *TComVUI)   SetMaxBitsPerMinCuDenom(i int) { 
	this.m_maxBitsPerMinCuDenom = i; 
}

func (this *TComVUI)   GetLog2MaxMvLengthHorizontal() int  { 
	return this.m_log2MaxMvLengthHorizontal; 
}
func (this *TComVUI)   SetLog2MaxMvLengthHorizontal(i int) { 
	this.m_log2MaxMvLengthHorizontal = i; 
}

func (this *TComVUI)   GetLog2MaxMvLengthVertical() int  {
	return this.m_log2MaxMvLengthVertical; 
}
func (this *TComVUI)   SetLog2MaxMvLengthVertical(i int) { 
	this.m_log2MaxMvLengthVertical = i; 
}

func (this *TComVUI)   SetTimingInfoPresentFlag             ( flag bool )  { 
	this.m_timingInfoPresentFlag = flag;               
}
func (this *TComVUI)   GetTimingInfoPresentFlag             ( )            bool{ 
	return this.m_timingInfoPresentFlag;               
}

func (this *TComVUI)   SetNumUnitsInTick                    ( value uint ) { 
	this.m_numUnitsInTick = value;                     
}
func (this *TComVUI)   GetNumUnitsInTick                    ( )            uint{ 
	return this.m_numUnitsInTick;                      
}

func (this *TComVUI)   SetTimeScale                         ( value uint ) { 
	this.m_timeScale = value;                          
}
func (this *TComVUI)   GetTimeScale                         ( )            uint{ 
	return this.m_timeScale;                           
}

func (this *TComVUI)   SetNalHrdParametersPresentFlag       ( flag bool )  { 
	this.m_nalHrdParametersPresentFlag = flag;         
}
func (this *TComVUI)   GetNalHrdParametersPresentFlag       ( )            bool{ 
	return this.m_nalHrdParametersPresentFlag;         
}

func (this *TComVUI)   SetVclHrdParametersPresentFlag       ( flag bool )  { 
	this.m_vclHrdParametersPresentFlag = flag;         
}
func (this *TComVUI)   GetVclHrdParametersPresentFlag       ( )            bool{ 
	return this.m_vclHrdParametersPresentFlag;         
}

func (this *TComVUI)   SetSubPicCpbParamsPresentFlag        ( flag bool )  { 
	this.m_subPicCpbParamsPresentFlag = flag;          
}
func (this *TComVUI)   GetSubPicCpbParamsPresentFlag        ( )            bool{ 
	return this.m_subPicCpbParamsPresentFlag;          
}
  
func (this *TComVUI)   SetTickDivisorMinus2                 ( value uint ) { 
	this.m_tickDivisorMinus2 = value;                  
}
func (this *TComVUI)   GetTickDivisorMinus2                 ( )            uint{ 
	return this.m_tickDivisorMinus2;                   
}

func (this *TComVUI)   SetDuCpbRemovalDelayLengthMinus1     ( value uint ) { 
	this.m_duCpbRemovalDelayLengthMinus1 = value;      
}
func (this *TComVUI)   GetDuCpbRemovalDelayLengthMinus1     ( )            uint{ 
	return this.m_duCpbRemovalDelayLengthMinus1;      
}

func (this *TComVUI)   SetBitRateScale                      ( value uint ) { 
	this.m_bitRateScale = value;                       
}
func (this *TComVUI)   GetBitRateScale                      ( )            uint{ 
	return this.m_bitRateScale;                        
}

func (this *TComVUI)   SetCpbSizeScale                      ( value uint ) { 
	this.m_cpbSizeScale = value;                       
}
func (this *TComVUI)   GetCpbSizeScale                      ( )            uint{ 
	return this.m_cpbSizeScale;                        
}

func (this *TComVUI)   SetInitialCpbRemovalDelayLengthMinus1( value uint ) { 
	this.m_initialCpbRemovalDelayLengthMinus1 = value; 
}
func (this *TComVUI)   GetInitialCpbRemovalDelayLengthMinus1( )            uint{ 
	return this.m_initialCpbRemovalDelayLengthMinus1;  
}

func (this *TComVUI)   SetCpbRemovalDelayLengthMinus1       ( value uint ) { 
	this.m_cpbRemovalDelayLengthMinus1 = value;        
}
func (this *TComVUI)   GetCpbRemovalDelayLengthMinus1       ( )            uint{ 
	return this.m_cpbRemovalDelayLengthMinus1;         
}

func (this *TComVUI)   SetDpbOutputDelayLengthMinus1        ( value uint ) { 
	this.m_dpbOutputDelayLengthMinus1 = value;         
}
func (this *TComVUI)   GetDpbOutputDelayLengthMinus1        ( )            uint{ 
	return this.m_dpbOutputDelayLengthMinus1;          
}

func (this *TComVUI)   SetFixedPicRateFlag       ( layer int, flag bool )  { 
	this.m_HRD[layer].fixedPicRateFlag = flag;         
}
func (this *TComVUI)   GetFixedPicRateFlag       ( layer int            )  bool{ 
	return this.m_HRD[layer].fixedPicRateFlag;         
}

func (this *TComVUI)   SetPicDurationInTcMinus1  ( layer int, value uint ) { 
	this.m_HRD[layer].picDurationInTcMinus1 = value;   
}
func (this *TComVUI)   GetPicDurationInTcMinus1  ( layer int             ) uint{ 
	return this.m_HRD[layer].picDurationInTcMinus1;    
}

func (this *TComVUI)   SetLowDelayHrdFlag        ( layer int, flag bool )  { 
	this.m_HRD[layer].lowDelayHrdFlag = flag;          
}
func (this *TComVUI)   GetLowDelayHrdFlag        ( layer int            )  bool{ 
	return this.m_HRD[layer].lowDelayHrdFlag;          
}

func (this *TComVUI)   SetCpbCntMinus1           ( layer int, value uint ) { 
	this.m_HRD[layer].cpbCntMinus1 = value; 
}
func (this *TComVUI)   GetCpbCntMinus1           ( layer int            )  uint{ 
	return this.m_HRD[layer].cpbCntMinus1; 
}

func (this *TComVUI)   SetBitRateValueMinus1     ( layer, cpbcnt, nalOrVcl int, value uint ) { 
	this.m_HRD[layer].bitRateValueMinus1[cpbcnt][nalOrVcl] = value; 
}
func (this *TComVUI)   GetBitRateValueMinus1     ( layer, cpbcnt, nalOrVcl int             ) uint{ 
	return this.m_HRD[layer].bitRateValueMinus1[cpbcnt][nalOrVcl];  
}

func (this *TComVUI)   SetCpbSizeValueMinus1     ( layer, cpbcnt, nalOrVcl int, value uint ) {
 	this.m_HRD[layer].cpbSizeValue[cpbcnt][nalOrVcl] = value;       
}
func (this *TComVUI)   GetCpbSizeValueMinus1     ( layer, cpbcnt, nalOrVcl int            )  uint{ 
	return this.m_HRD[layer].cpbSizeValue[cpbcnt][nalOrVcl];        
}

func (this *TComVUI)   SetCbrFlag                ( layer, cpbcnt, nalOrVcl int, value bool ) { 
	this.m_HRD[layer].cbrFlag[cpbcnt][nalOrVcl] = value;            
}
func (this *TComVUI)   GetCbrFlag                ( layer, cpbcnt, nalOrVcl int             ) bool{ 
	return this.m_HRD[layer].cbrFlag[cpbcnt][nalOrVcl];             
}


func (this *TComVUI)   SetNumDU                              ( value uint ) { 
	this.m_numDU = value;                            
}
func (this *TComVUI)   GetNumDU                              ( )            uint{ 
	return this.m_numDU;          
}

/// SPS class
type TComSPS struct{
//private:
  m_SPSId				int;
  m_VPSId				int;
  m_chromaFormatIdc		int;

  m_uiMaxTLayers		uint;           // maximum number of temporal layers

  	// Structure
  m_picWidthInLumaSamples		uint;
  m_picHeightInLumaSamples		uint;
  m_picCroppingFlag				bool;
  m_picCropLeftOffset			int;
  m_picCropRightOffset			int;
  m_picCropTopOffset			int;
  m_picCropBottomOffset			int;
  m_uiMaxCUWidth				uint;
  m_uiMaxCUHeight				uint;
  m_uiMaxCUDepth				uint;
  m_uiMinTrDepth				uint;
  m_uiMaxTrDepth				uint;
  m_RPSList						TComRPSList;
  m_bLongTermRefsPresent		bool;
  m_TMVPFlagsPresent			bool;
  m_numReorderPics	[MAX_TLAYER]int;
  
  	// Tool list
  m_uiQuadtreeTULog2MaxSize		uint;
  m_uiQuadtreeTULog2MinSize		uint;
  m_uiQuadtreeTUMaxDepthInter	uint;
  m_uiQuadtreeTUMaxDepthIntra	uint;
  m_usePCM						bool;
  m_pcmLog2MaxSize				uint;
  m_uiPCMLog2MinSize			uint;
  m_useAMP						bool;

  m_bUseLComb					bool;
  
  m_restrictedRefPicListsFlag		bool;
  m_listsModificationPresentFlag	bool;

  	// Parameter
  m_bitDepthY					int;
  m_bitDepthC					int;
  m_qpBDOffsetY					int;
  m_qpBDOffsetC					int;

  m_useLossless					bool;

  m_uiPCMBitDepthLuma			uint;
  m_uiPCMBitDepthChroma			uint;
  m_bPCMFilterDisableFlag		bool;

  m_uiBitsForPOC				uint;
  m_numLongTermRefPicSPS		uint;
  m_ltRefPicPocLsbSps		[33]uint;
  m_usedByCurrPicLtSPSFlag	[33]bool;
  	// Max physical transform size
  m_uiMaxTrSize					uint;
  
  m_iAMPAcc			[MAX_CU_DEPTH]int;
  m_bUseSAO						bool; 

  m_bTemporalIdNestingFlag		bool; // temporal_id_nesting_flag

  m_scalingListEnabledFlag		bool;
  m_scalingListPresentFlag		bool;
  m_scalingList		*TComScalingList;   //!< ScalingList class pointer
  m_uiMaxDecPicBuffering	[MAX_TLAYER]uint; 
  m_uiMaxLatencyIncrease	[MAX_TLAYER]uint;

  m_useDF						bool;
//NTRA_SMOOTHING
  m_useStrongIntraSmoothing		bool; 
//

  m_vuiParametersPresentFlag	bool;
  m_vuiParameters				TComVUI;

  m_cropUnitX	[MAX_CHROMA_FORMAT_IDC+1]int;
  m_cropUnitY	[MAX_CHROMA_FORMAT_IDC+1]int;
  m_pcPTL						TComPTL;
 }
 
//public:
func NewTComSPS() *TComSPS{
	return &TComSPS{}
}

func (this *TComSPS)  GetVPSId       () int     { 
	return this.m_VPSId;          
}
func (this *TComSPS)  SetVPSId       (i int)    { 
	this.m_VPSId = i;             
}
func (this *TComSPS)  GetSPSId       () int     { 
	return this.m_SPSId;          
}
func (this *TComSPS)  SetSPSId       (i int)    { 
	this.m_SPSId = i;             
}
func (this *TComSPS)  GetChromaFormatIdc () int     { 
	return this.m_chromaFormatIdc;       
}
func (this *TComSPS)  SetChromaFormatIdc (i int)    { 
	this.m_chromaFormatIdc = i;          
}

func (this *TComSPS)  GetCropUnitX (chromaFormatIdc int) int { 
	//assert (chromaFormatIdc > 0 && chromaFormatIdc <= MAX_CHROMA_FORMAT_IDC); 
	return this.m_cropUnitX[chromaFormatIdc];      
}
func (this *TComSPS)  GetCropUnitY (chromaFormatIdc int) int { 
	//assert (chromaFormatIdc > 0 && chromaFormatIdc <= MAX_CHROMA_FORMAT_IDC); 
	return this.m_cropUnitY[chromaFormatIdc];      
}
  
  // structure
func (this *TComSPS)  SetPicWidthInLumaSamples       ( u uint ) { 
	this.m_picWidthInLumaSamples = u;        
}
func (this *TComSPS)  GetPicWidthInLumaSamples       ()  uint   { 
	return  this.m_picWidthInLumaSamples;    
}
func (this *TComSPS)  SetPicHeightInLumaSamples      ( u uint ) { 
	this.m_picHeightInLumaSamples = u;       
}
func (this *TComSPS)  GetPicHeightInLumaSamples      ()  uint   { 
	return  this.m_picHeightInLumaSamples;   
}

func (this *TComSPS)  GetPicCroppingFlag() bool           { 
	return this.m_picCroppingFlag; 
}
func (this *TComSPS)  SetPicCroppingFlag(val bool)        { 
	this.m_picCroppingFlag = val; 
}
func (this *TComSPS)  GetPicCropLeftOffset() int          { 
	return this.m_picCropLeftOffset; 
}
func (this *TComSPS)  SetPicCropLeftOffset(val int)       { 
	this.m_picCropLeftOffset = val; 
}
func (this *TComSPS)  GetPicCropRightOffset() int         { 
	return this.m_picCropRightOffset; 
}
func (this *TComSPS)  SetPicCropRightOffset(val int)      { 
	this.m_picCropRightOffset = val; 
}
func (this *TComSPS)  GetPicCropTopOffset() int           { 
	return this.m_picCropTopOffset; 
}
func (this *TComSPS)  SetPicCropTopOffset(val int)        { 
	this.m_picCropTopOffset = val; 
}
func (this *TComSPS)  GetPicCropBottomOffset() int        { 
	return this.m_picCropBottomOffset; 
}
func (this *TComSPS)  SetPicCropBottomOffset(val int)     { 
	this.m_picCropBottomOffset = val; 
}
func (this *TComSPS)  GetNumLongTermRefPicSPS() uint     { 
	return this.m_numLongTermRefPicSPS; 
}
func (this *TComSPS)  SetNumLongTermRefPicSPS(val uint)  { 
	this.m_numLongTermRefPicSPS = val; 
}

func (this *TComSPS)  GetLtRefPicPocLsbSps(index uint) uint     { 
	return this.m_ltRefPicPocLsbSps[index]; 
}
func (this *TComSPS)  SetLtRefPicPocLsbSps(index, val uint)     { 
	this.m_ltRefPicPocLsbSps[index] = val; 
}

func (this *TComSPS)  GetUsedByCurrPicLtSPSFlag(i int) bool      {
	return this.m_usedByCurrPicLtSPSFlag[i];
}
func (this *TComSPS)  SetUsedByCurrPicLtSPSFlag(i int, x bool)   { 
	this.m_usedByCurrPicLtSPSFlag[i] = x;
}
func (this *TComSPS)  SetMaxCUWidth  ( u uint ) { 
	this.m_uiMaxCUWidth = u;      
}
func (this *TComSPS)  GetMaxCUWidth  ()  uint   { 
	return  this.m_uiMaxCUWidth;  
}
func (this *TComSPS)  SetMaxCUHeight ( u uint ) { 
	this.m_uiMaxCUHeight = u;     
}
func (this *TComSPS)  GetMaxCUHeight ()  uint   { 
	return  this.m_uiMaxCUHeight; 
}
func (this *TComSPS)  SetMaxCUDepth  ( u uint)  { 
	this.m_uiMaxCUDepth = u;      
}
func (this *TComSPS)  GetMaxCUDepth  ()  uint   { 
	return  this.m_uiMaxCUDepth;  
}
func (this *TComSPS)  SetUsePCM      ( b bool ) { 
	this.m_usePCM = b;           
}
func (this *TComSPS)  GetUsePCM      ()  bool   { 
	return this.m_usePCM;        
}
func (this *TComSPS)  SetPCMLog2MaxSize  ( u uint ) { 
	this.m_pcmLog2MaxSize = u;      
}
func (this *TComSPS)  GetPCMLog2MaxSize  ()  uint   { 
	return  this.m_pcmLog2MaxSize;  
}
func (this *TComSPS)  SetPCMLog2MinSize  ( u uint ) { 
	this.m_uiPCMLog2MinSize = u;      
}
func (this *TComSPS)  GetPCMLog2MinSize  ()  uint   { 
	return  this.m_uiPCMLog2MinSize;  
}
func (this *TComSPS)  SetBitsForPOC  ( u uint ) 	 { 
	this.m_uiBitsForPOC = u;      
}
func (this *TComSPS)  GetBitsForPOC  ()  uint   	 { 
	return this.m_uiBitsForPOC;   
}
func (this *TComSPS)  GetUseAMP()  bool	{ 
	return this.m_useAMP; 
}
func (this *TComSPS)  SetUseAMP( b bool ) 	{ 
	this.m_useAMP = b; 
}
func (this *TComSPS)  SetMinTrDepth  ( u uint ) { 
	this.m_uiMinTrDepth = u;      
}
func (this *TComSPS)  GetMinTrDepth  ()  uint   { 
	return  this.m_uiMinTrDepth;  
}
func (this *TComSPS)  SetMaxTrDepth  ( u uint ) { 
	this.m_uiMaxTrDepth = u;      
}
func (this *TComSPS)  GetMaxTrDepth  ()  uint   { 
	return  this.m_uiMaxTrDepth;  
}
func (this *TComSPS)  SetQuadtreeTULog2MaxSize( u uint ) { 
	this.m_uiQuadtreeTULog2MaxSize = u;    
}
func (this *TComSPS)  GetQuadtreeTULog2MaxSize()  uint   { 
	return this.m_uiQuadtreeTULog2MaxSize; 
}
func (this *TComSPS)  SetQuadtreeTULog2MinSize( u uint ) { 
	this.m_uiQuadtreeTULog2MinSize = u;    
}
func (this *TComSPS)  GetQuadtreeTULog2MinSize()  uint   { 
	return this.m_uiQuadtreeTULog2MinSize; 
}
func (this *TComSPS)  SetQuadtreeTUMaxDepthInter( u uint ) { 
	this.m_uiQuadtreeTUMaxDepthInter = u;    
}
func (this *TComSPS)  SetQuadtreeTUMaxDepthIntra( u uint ) { 
	this.m_uiQuadtreeTUMaxDepthIntra = u;    
}
func (this *TComSPS)  GetQuadtreeTUMaxDepthInter()  uint   { 
	return this.m_uiQuadtreeTUMaxDepthInter; 
}
func (this *TComSPS)  GetQuadtreeTUMaxDepthIntra()  uint   { 
	return this.m_uiQuadtreeTUMaxDepthIntra; 
}
func (this *TComSPS)  SetNumReorderPics(i int, tlayer uint)              { 
	this.m_numReorderPics[tlayer] = i;    
}
func (this *TComSPS)  GetNumReorderPics(tlayer uint)  int               { 
	return this.m_numReorderPics[tlayer]; 
}
func (this *TComSPS)  CreateRPSList( numRPS int ){
}
func (this *TComSPS)  GetRPSList() *TComRPSList                  { 
	return &this.m_RPSList;          
}
func (this *TComSPS)  GetLongTermRefsPresent() bool        { 
	return this.m_bLongTermRefsPresent; 
}
func (this *TComSPS)  SetLongTermRefsPresent(b bool)       { 
	this.m_bLongTermRefsPresent=b;      
}
func (this *TComSPS)  GetTMVPFlagsPresent() bool           { 
	return this.m_TMVPFlagsPresent; 
}
func (this *TComSPS)  SetTMVPFlagsPresent(b bool)          { 
	this.m_TMVPFlagsPresent=b;      
}  
  // physical transform
func (this *TComSPS)  SetMaxTrSize   ( u uint ) { 
	this.m_uiMaxTrSize = u;       
}
func (this *TComSPS)  GetMaxTrSize   ()  uint   { 
	return  this.m_uiMaxTrSize;   
}
  
  // Tool list
func (this *TComSPS)  SetUseLComb    (b bool)   { 
	this.m_bUseLComb = b;         
}
func (this *TComSPS)  GetUseLComb    () bool    { 
	return this.m_bUseLComb;      
}

func (this *TComSPS)  GetUseLossless () bool    { 
	return this.m_useLossless; 
}
func (this *TComSPS)  SetUseLossless (b bool )  { 
	this.m_useLossless  = b; 
}
  
func (this *TComSPS)  GetRestrictedRefPicListsFlag    ()  bool    { 
	return this.m_restrictedRefPicListsFlag;   
}
func (this *TComSPS)  SetRestrictedRefPicListsFlag    ( b bool )  { 
	this.m_restrictedRefPicListsFlag = b;      
}
func (this *TComSPS)  GetListsModificationPresentFlag ()  bool    { 	
	return this.m_listsModificationPresentFlag; 
}
func (this *TComSPS)  SetListsModificationPresentFlag ( b bool )  { 
	this.m_listsModificationPresentFlag = b;    
}

  // AMP accuracy
func (this *TComSPS)  GetAMPAcc   ( uiDepth uint ) int 	   { 
	return this.m_iAMPAcc[uiDepth]; 
}
func (this *TComSPS)  SetAMPAcc   ( uiDepth uint, iAccu int) { 
	//assert( uiDepth < g_uiMaxCUDepth);  
	this.m_iAMPAcc[uiDepth] = iAccu; 
}

  // Bit-depth
func (this *TComSPS)  GetBitDepthY() int  { 
	return this.m_bitDepthY; 
}
func (this *TComSPS)  SetBitDepthY(u int) { 
	this.m_bitDepthY = u; 
}
func (this *TComSPS)  GetBitDepthC() int  { 
	return this.m_bitDepthC; 
}
func (this *TComSPS)  SetBitDepthC(u int) { 
	this.m_bitDepthC = u; 
}
func (this *TComSPS)  GetQpBDOffsetY  () int         { 
	return this.m_qpBDOffsetY;   
}
func (this *TComSPS)  SetQpBDOffsetY  ( value int  ) { 
	this.m_qpBDOffsetY = value;  
}
func (this *TComSPS)  GetQpBDOffsetC  () int         { 
	return this.m_qpBDOffsetC;   
}
func (this *TComSPS)  SetQpBDOffsetC  ( value int  ) { 
	this.m_qpBDOffsetC = value;  
}
func (this *TComSPS)  SetUseSAO                  (bVal bool)  {
	this.m_bUseSAO = bVal;
}
func (this *TComSPS)  GetUseSAO                  ()    bool   {
	return this.m_bUseSAO;
}

func (this *TComSPS)  GetMaxTLayers() uint                      { 
	return this.m_uiMaxTLayers; 
}
func (this *TComSPS)  SetMaxTLayers( uiMaxTLayers uint )        { 
	//assert( uiMaxTLayers <= MAX_TLAYER ); 
	this.m_uiMaxTLayers = uiMaxTLayers; 
}

func (this *TComSPS)  GetTemporalIdNestingFlag() bool           { 
	return this.m_bTemporalIdNestingFlag; 
}
func (this *TComSPS)  SetTemporalIdNestingFlag( bValue bool )   { 
	this.m_bTemporalIdNestingFlag = bValue; 
}
func (this *TComSPS)  GetPCMBitDepthLuma     ()  uint   { 
	return this.m_uiPCMBitDepthLuma;     
}
func (this *TComSPS)  SetPCMBitDepthLuma     ( u uint ) { 
	this.m_uiPCMBitDepthLuma = u;        
}
func (this *TComSPS)  GetPCMBitDepthChroma   ()  uint   { 
	return this.m_uiPCMBitDepthChroma;   
}
func (this *TComSPS)  SetPCMBitDepthChroma   ( u uint ) { 
	this.m_uiPCMBitDepthChroma = u;      
}
func (this *TComSPS)  SetPCMFilterDisableFlag     ( bValue bool )    { 
	this.m_bPCMFilterDisableFlag = bValue; 
}
func (this *TComSPS)  GetPCMFilterDisableFlag     ()       bool      { 
	return this.m_bPCMFilterDisableFlag;   
} 

func (this *TComSPS)  GetScalingListFlag       ()  bool   { 
	return this.m_scalingListEnabledFlag;     
}
func (this *TComSPS)  SetScalingListFlag       ( b bool ) { 
	this.m_scalingListEnabledFlag  = b;       
}
func (this *TComSPS)  GetScalingListPresentFlag()  bool   { 
	return this.m_scalingListPresentFlag;     
}
func (this *TComSPS)  SetScalingListPresentFlag( b bool ) { 
	this.m_scalingListPresentFlag  = b;       
}
func (this *TComSPS)  SetScalingList      ( scalingList *TComScalingList){
	this.m_scalingList = scalingList;
}
func (this *TComSPS)  GetScalingList () *TComScalingList      { 
	return this.m_scalingList; 
}               //!< get ScalingList class pointer in SPS
func (this *TComSPS)  GetMaxDecPicBuffering  (tlayer uint) uint           { 
	return this.m_uiMaxDecPicBuffering[tlayer]; 
}
func (this *TComSPS)  SetMaxDecPicBuffering  (ui, tlayer uint)			   { 
	this.m_uiMaxDecPicBuffering[tlayer] = ui;   
}
func (this *TComSPS)  GetMaxLatencyIncrease  (tlayer uint) uint           { 
	return this.m_uiMaxLatencyIncrease[tlayer];   
}
func (this *TComSPS)  SetMaxLatencyIncrease  (ui, tlayer uint)			   { 
	this.m_uiMaxLatencyIncrease[tlayer] = ui;     
}

//#if STRONG_INTRA_SMOOTHING
func (this *TComSPS)  SetUseStrongIntraSmoothing (bVal bool)  {
	this.m_useStrongIntraSmoothing = bVal;
}
func (this *TComSPS)  GetUseStrongIntraSmoothing ()    bool   {
	return this.m_useStrongIntraSmoothing;
}
//#endif

func (this *TComSPS)  GetVuiParametersPresentFlag() bool  { 
	return this.m_vuiParametersPresentFlag; 
}
func (this *TComSPS)  SetVuiParametersPresentFlag(b bool) { 
	this.m_vuiParametersPresentFlag = b; 
}
func (this *TComSPS)  GetVuiParameters() *TComVUI	   { 
	return &this.m_vuiParameters; 
}
func (this *TComSPS)  SetHrdParameters(frameRate, numDU, bitRate uint, randomAccess bool){
}

func (this *TComSPS)  GetPTL() *TComPTL     { 
	return &this.m_pcPTL; 
}
//};

/// Reference Picture Lists class
type TComRefPicListModification struct{
//private:
  m_bRefPicListModificationFlagL0	bool;  
  m_bRefPicListModificationFlagL1	bool;  
  m_RefPicSetIdxL0				[32]uint;
  m_RefPicSetIdxL1				[32]uint;
}
    
//public:
func NewTComRefPicListModification() *TComRefPicListModification{
	return &TComRefPicListModification{};
}
 
func (this *TComRefPicListModification)  Create                    (){
}
func (this *TComRefPicListModification)  Destroy                   (){
}

func (this *TComRefPicListModification)  GetRefPicListModificationFlagL0()    bool  { 
	return this.m_bRefPicListModificationFlagL0; 
}
func (this *TComRefPicListModification)  SetRefPicListModificationFlagL0(flag bool) { 
	this.m_bRefPicListModificationFlagL0 = flag; 
}
func (this *TComRefPicListModification)  GetRefPicListModificationFlagL1() 	  bool  { 
	return this.m_bRefPicListModificationFlagL1; 
}
func (this *TComRefPicListModification)  SetRefPicListModificationFlagL1(flag bool) { 
	this.m_bRefPicListModificationFlagL1 = flag; 
}
func (this *TComRefPicListModification)  SetRefPicSetIdxL0(idx, refPicSetIdx uint)  { 
	this.m_RefPicSetIdxL0[idx] = refPicSetIdx; 
}
func (this *TComRefPicListModification)  GetRefPicSetIdxL0(idx uint) 		 uint   { 
	return this.m_RefPicSetIdxL0[idx]; 
}
func (this *TComRefPicListModification)  SetRefPicSetIdxL1(idx, refPicSetIdx uint)  { 
	this.m_RefPicSetIdxL1[idx] = refPicSetIdx; 
}
func (this *TComRefPicListModification)  GetRefPicSetIdxL1(idx uint) 		 uint   { 
	return this.m_RefPicSetIdxL1[idx]; 
}


/// PPS class
type TComPPS struct{
//private:
  m_PPSId	int;                    // pic_parameter_set_id
  m_SPSId	int;                    // seq_parameter_set_id
  m_picInitQPMinus26 int;
  m_useDQP	bool;
  m_bConstrainedIntraPred	bool;    // constrained_intra_pred_flag
  m_bSliceChromaQpFlag		bool;       // slicelevel_chroma_qp_flag

  	// access channel
  m_pcSPS			*TComSPS;
  m_uiMaxCuDQPDepth	uint;
  m_uiMinCuDQPSize	uint;

  m_chromaCbQpOffset	int;
  m_chromaCrQpOffset	int;

  m_numRefIdxL0DefaultActive	uint;
  m_numRefIdxL1DefaultActive	uint;

  m_bUseWeightPred			bool;           // Use of Weighting Prediction (P_SLICE)
  m_useWeightedBiPred		bool;        // Use of Weighting Bi-Prediction (B_SLICE)
  m_OutputFlagPresentFlag	bool;   // Indicates the presence of output_flag in slice header

  m_TransquantBypassEnableFlag	bool; // Indicates presence of cu_transquant_bypass_flag in CUs.
  m_useTransformSkip			bool;
  m_dependentSliceEnabledFlag	bool;     //!< Indicates the presence of dependent slices
  m_tilesEnabledFlag			bool;              //!< Indicates the presence of tiles
  m_entropyCodingSyncEnabledFlag	bool;  //!< Indicates the presence of wavefronts
//#if !REMOVE_ENTROPY_SLICES
//  Bool        m_entropySliceEnabledFlag;       //!< Indicates the presence of entropy slices
//#endif
  
  m_loopFilterAcrossTilesEnabledFlag	bool;
  m_uniformSpacingFlag					bool;
  m_iNumColumnsMinus1					int;
  m_puiColumnWidth					   []uint;
  m_iNumRowsMinus1						int;
  m_puiRowHeight					   []uint;

  m_iNumSubstreams						int;

  m_signHideFlag						int;

  m_cabacInitPresentFlag				bool;
  m_encCABACTableIdx					uint;           // Used to transmit table selection across slices

  m_sliceHeaderExtensionPresentFlag	bool;
  m_loopFilterAcrossSlicesEnabledFlag		bool;
  m_deblockingFilterControlPresentFlag		bool;
  m_deblockingFilterOverrideEnabledFlag	bool;
  m_picDisableDeblockingFilterFlag			bool;
  m_deblockingFilterBetaOffsetDiv2		int;    //< beta offset for deblocking filter
  m_deblockingFilterTcOffsetDiv2		int;      //< tc offset for deblocking filter
  m_scalingListPresentFlag				bool;
  m_scalingList	*TComScalingList;   //!< ScalingList class pointer
  m_log2ParallelMergeLevelMinus2			uint;
}

//public:
func NewTComPPS() *TComPPS{
	return &TComPPS{}
}
  
func (this *TComPPS)  GetPPSId () int  { 
	return this.m_PPSId; 
}
func (this *TComPPS)  SetPPSId (i int) { 
	this.m_PPSId = i; 
}
func (this *TComPPS)  GetSPSId () int  { 
	return this.m_SPSId; 
}
func (this *TComPPS)  SetSPSId (i int) { 
	this.m_SPSId = i; 
}
  
func (this *TComPPS)  GetPicInitQPMinus26 ()  int    { 
	return  this.m_picInitQPMinus26; 
}
func (this *TComPPS)  SetPicInitQPMinus26 ( i int )  { 
	this.m_picInitQPMinus26 = i;     
}
func (this *TComPPS)  GetUseDQP ()  bool             { 
	return this.m_useDQP;        
}
func (this *TComPPS)  SetUseDQP ( b bool )           { 
	this.m_useDQP   = b;         
}
func (this *TComPPS)  GetConstrainedIntraPred ()  bool   { 
	return  this.m_bConstrainedIntraPred; 
}
func (this *TComPPS)  SetConstrainedIntraPred ( b bool ) { 
	this.m_bConstrainedIntraPred = b;     
}
func (this *TComPPS)  GetSliceChromaQpFlag ()  bool   { 
	return  this.m_bSliceChromaQpFlag; 
}
func (this *TComPPS)  SetSliceChromaQpFlag ( b bool ) { 
	this.m_bSliceChromaQpFlag = b;     
}

func (this *TComPPS)  SetSPS              ( pcSPS *TComSPS) { 
	this.m_pcSPS = pcSPS; 
}
func (this *TComPPS)  GetSPS              ()      *TComSPS  { 
	return this.m_pcSPS;          
}
func (this *TComPPS)  SetMaxCuDQPDepth    ( u uint ) { 
	this.m_uiMaxCuDQPDepth = u;   
}
func (this *TComPPS)  GetMaxCuDQPDepth    ()  uint   { 
	return this.m_uiMaxCuDQPDepth;
}
func (this *TComPPS)  SetMinCuDQPSize     ( u uint ) { 
	this.m_uiMinCuDQPSize = u;    
}
func (this *TComPPS)  GetMinCuDQPSize     ()  uint   { 
	return this.m_uiMinCuDQPSize; 
}

func (this *TComPPS)  SetChromaCbQpOffset( i int ) { 
	this.m_chromaCbQpOffset = i;    
}
func (this *TComPPS)  GetChromaCbQpOffset()  int   { 
	return this.m_chromaCbQpOffset; 
}
func (this *TComPPS)  SetChromaCrQpOffset( i int ) { 
	this.m_chromaCrQpOffset = i;    
}
func (this *TComPPS)  GetChromaCrQpOffset()  int   { 
	return this.m_chromaCrQpOffset; 
}

func (this *TComPPS)  SetNumRefIdxL0DefaultActive(ui uint)    { 
	this.m_numRefIdxL0DefaultActive=ui;     
}
func (this *TComPPS)  GetNumRefIdxL0DefaultActive()  uint     { 
	return this.m_numRefIdxL0DefaultActive; 
}
func (this *TComPPS)  SetNumRefIdxL1DefaultActive(ui uint)    { 
	this.m_numRefIdxL1DefaultActive=ui;     
}
func (this *TComPPS)  GetNumRefIdxL1DefaultActive()  uint     { 
	return this.m_numRefIdxL1DefaultActive; 
}

func (this *TComPPS)  GetUseWP                     ()  bool    { 
	return this.m_bUseWeightPred;  
}
func (this *TComPPS)  GetWPBiPred                  ()  bool    { 
	return this.m_useWeightedBiPred;     
}
func (this *TComPPS)  SetUseWP                     ( b bool )  { 
	this.m_bUseWeightPred = b;     
}
func (this *TComPPS)  SetWPBiPred                  ( b bool )  { 
	this.m_useWeightedBiPred = b;  
}
func (this *TComPPS)  SetOutputFlagPresentFlag( b bool )  { 
	this.m_OutputFlagPresentFlag = b;    
}
func (this *TComPPS)  GetOutputFlagPresentFlag()  bool    { 
	return this.m_OutputFlagPresentFlag; 
}
func (this *TComPPS)  SetTransquantBypassEnableFlag( b bool ) { 
	this.m_TransquantBypassEnableFlag = b; 
}
func (this *TComPPS)  GetTransquantBypassEnableFlag()  bool   { 
	return this.m_TransquantBypassEnableFlag; 
}

func (this *TComPPS)  GetUseTransformSkip       ()  bool   { 
	return this.m_useTransformSkip;     
}
func (this *TComPPS)  SetUseTransformSkip       ( b bool ) { 
	this.m_useTransformSkip  = b;       
}

func (this *TComPPS)  SetLoopFilterAcrossTilesEnabledFlag  (b bool)    { 
	this.m_loopFilterAcrossTilesEnabledFlag = b; 
}
func (this *TComPPS)  GetLoopFilterAcrossTilesEnabledFlag  () bool     { 
	return this.m_loopFilterAcrossTilesEnabledFlag;  
}
func (this *TComPPS)  GetDependentSliceEnabledFlag()   bool            { 
	return this.m_dependentSliceEnabledFlag; 
}
func (this *TComPPS)  SetDependentSliceEnabledFlag(val bool)           { 
	this.m_dependentSliceEnabledFlag = val; 
}
func (this *TComPPS)  GetTilesEnabledFlag()            bool            { 
	return this.m_tilesEnabledFlag; 
}
func (this *TComPPS)  SetTilesEnabledFlag(val bool)                    { 
	this.m_tilesEnabledFlag = val; 
}
func (this *TComPPS)  GetEntropyCodingSyncEnabledFlag() bool           { 
	return this.m_entropyCodingSyncEnabledFlag; 
}
func (this *TComPPS)  SetEntropyCodingSyncEnabledFlag(val bool)        { 
	this.m_entropyCodingSyncEnabledFlag = val; 
}
/*#if !REMOVE_ENTROPY_SLICES
  Bool    GetEntropySliceEnabledFlag() const               { return this.m_entropySliceEnabledFlag; }
  Void    SetEntropySliceEnabledFlag(Bool val)             { this.m_entropySliceEnabledFlag = val; }
#endif*/
func (this *TComPPS)  SetUniformSpacingFlag            ( b bool )          { 
	this.m_uniformSpacingFlag = b; 
}
func (this *TComPPS)  GetUniformSpacingFlag            ()  bool            { 
	return this.m_uniformSpacingFlag; 
}
func (this *TComPPS)  SetNumColumnsMinus1              ( i int )           { 
	this.m_iNumColumnsMinus1 = i; 
}
func (this *TComPPS)  GetNumColumnsMinus1              ()  int             { 
	return this.m_iNumColumnsMinus1; 
}
func (this *TComPPS)  SetColumnWidth ( columnWidth []uint){
	if this.m_uniformSpacingFlag == false && this.m_iNumColumnsMinus1 > 0 {
      this.m_puiColumnWidth = make([]uint, this.m_iNumColumnsMinus1);
      for i:=0; i<this.m_iNumColumnsMinus1; i++ {
        this.m_puiColumnWidth[i] = columnWidth[i];
      }
    }
}
func (this *TComPPS)  GetColumnWidth  (columnIdx uint) uint { 
	return this.m_puiColumnWidth[columnIdx]; 
}
func (this *TComPPS)  SetNumRowsMinus1(i int)               { 
	this.m_iNumRowsMinus1 = i; 
}
func (this *TComPPS)  GetNumRowsMinus1() int                { 
	return this.m_iNumRowsMinus1; 
}
func (this *TComPPS)  SetRowHeight    ( rowHeight []uint )   {
    if this.m_uniformSpacingFlag == false && this.m_iNumRowsMinus1 > 0 {
      this.m_puiRowHeight = make([]uint, this.m_iNumRowsMinus1);
      for i:=0; i<this.m_iNumRowsMinus1; i++ {
        this.m_puiRowHeight[i] = rowHeight[i];
      }
    }
  }
func (this *TComPPS)  GetRowHeight    (rowIdx uint) uint    { 
	return this.m_puiRowHeight[rowIdx]; 
}
func (this *TComPPS)  SetNumSubstreams(iNumSubstreams int)               { 
	this.m_iNumSubstreams = iNumSubstreams; 
}
func (this *TComPPS)  GetNumSubstreams()              int                { 
	return this.m_iNumSubstreams; 
}

func (this *TComPPS)  SetSignHideFlag( signHideFlag int )  { 
	this.m_signHideFlag = signHideFlag; 
}
func (this *TComPPS)  GetSignHideFlag()             int    { 
	return this.m_signHideFlag; 
}

func (this *TComPPS)  SetCabacInitPresentFlag( flag bool )     { 
	this.m_cabacInitPresentFlag = flag;    
}
func (this *TComPPS)  SetEncCABACTableIdx( idx uint )           { 
	this.m_encCABACTableIdx = idx;         
}
func (this *TComPPS)  GetCabacInitPresentFlag()     bool       { 
	return this.m_cabacInitPresentFlag;    
}
func (this *TComPPS)  GetEncCABACTableIdx()         uint       { 
	return this.m_encCABACTableIdx;        
}
func (this *TComPPS)  SetDeblockingFilterControlPresentFlag( val bool )  { 
	this.m_deblockingFilterControlPresentFlag = val; 
}
func (this *TComPPS)  GetDeblockingFilterControlPresentFlag()    bool    { 
	return this.m_deblockingFilterControlPresentFlag; 
}
func (this *TComPPS)  SetDeblockingFilterOverrideEnabledFlag(val bool )  { 
	this.m_deblockingFilterOverrideEnabledFlag = val; 
}
func (this *TComPPS)  GetDeblockingFilterOverrideEnabledFlag()   bool    { 
	return this.m_deblockingFilterOverrideEnabledFlag; 
}
func (this *TComPPS)  SetPicDisableDeblockingFilterFlag(val bool)        { 
	this.m_picDisableDeblockingFilterFlag = val; 
}       //!< Set offSet for deblocking filter disabled
func (this *TComPPS)  GetPicDisableDeblockingFilterFlag()   bool         { 
	return this.m_picDisableDeblockingFilterFlag; 
}      //!< Get offset for deblocking filter disabled
func (this *TComPPS)  SetDeblockingFilterBetaOffsetDiv2(val int)         { 
	this.m_deblockingFilterBetaOffsetDiv2 = val; 
}       //!< set beta offset for deblocking filter
func (this *TComPPS)  GetDeblockingFilterBetaOffsetDiv2()   int          { 
	return this.m_deblockingFilterBetaOffsetDiv2; 
}      //!< Get beta offset for deblocking filter
func (this *TComPPS)  SetDeblockingFilterTcOffsetDiv2(val   int)         { 
	this.m_deblockingFilterTcOffsetDiv2 = val; 
}               //!< set tc offset for deblocking filter
func (this *TComPPS)  GetDeblockingFilterTcOffsetDiv2()     int          { 
	return this.m_deblockingFilterTcOffsetDiv2; 
}              //!< Get tc offset for deblocking filter
func (this *TComPPS)  GetScalingListPresentFlag()         	 bool		  { 
	return this.m_scalingListPresentFlag;     
}
func (this *TComPPS)  SetScalingListPresentFlag( b bool ) 				  { 
	this.m_scalingListPresentFlag  = b;       }

func (this *TComPPS)  SetScalingList( scalingList *TComScalingList) 	  {
	this.m_scalingList = scalingList;
}
func (this *TComPPS)  GetScalingList () *TComScalingList         { 
	return this.m_scalingList; 
}         //!< Get ScalingList class pointer in PPS
func (this *TComPPS)  GetLog2ParallelMergeLevelMinus2      ()   uint         { 
	return this.m_log2ParallelMergeLevelMinus2; 
}
func (this *TComPPS)  SetLog2ParallelMergeLevelMinus2      (mrgLevel uint)   { 
	this.m_log2ParallelMergeLevelMinus2 = mrgLevel; 
}
func (this *TComPPS)  SetLoopFilterAcrossSlicesEnabledFlag ( bValue bool )    { 
	this.m_loopFilterAcrossSlicesEnabledFlag = bValue; 
}
func (this *TComPPS)  GetLoopFilterAcrossSlicesEnabledFlag ()       bool      { 
	return this.m_loopFilterAcrossSlicesEnabledFlag;   
} 
func (this *TComPPS)  GetSliceHeaderExtensionPresentFlag   ()            bool      { 
	return this.m_sliceHeaderExtensionPresentFlag;
}
func (this *TComPPS)  SetSliceHeaderExtensionPresentFlag   (val bool)              { 
	this.m_sliceHeaderExtensionPresentFlag = val; 
}


type wpScalingParam struct {
  // Explicit weighted prediction parameters parsed in slice header,
  // or Implicit weighted prediction parameters (8 bits depth values).
  bPresentFlag			bool;
  uiLog2WeightDenom		uint;
  iWeight				int;
  iOffset				int;

  // Weighted prediction scaling values built from above parameters (bitdepth scaled):
  w, o, offset, shift, round int;
};

type wpACDCParam struct {
  iAC int64;
  iDC	int64;
};

/// slice header class
type TComSlice struct{
//private:
  //  Bitstream writing
  m_saoEnabledFlag			bool;
  m_saoEnabledFlagChroma	bool;      ///< SAO Cb&Cr enabled flag
  m_iPPSId	int;               ///< picture parameter set ID
  m_PicOutputFlag	bool;        ///< pic_output_flag 
  m_iPOC	int;
  m_iLastIDR	int;
  m_prevPOC	int;
  m_pcRPS	*TComReferencePictureSet;
  m_LocalRPS *TComReferencePictureSet;
  m_iBDidx	int; 
  m_iCombinationBDidx	int;
  m_bCombineWithReferenceFlag	bool;
  m_RefPicListModification	TComRefPicListModification;
  m_eNalUnitType	NalUnitType;         ///< Nal unit type for the slice
  m_eSliceType		SliceType;
  m_iSliceQp		int;
  m_dependentSliceFlag	bool;
//#if ADAPTIVE_QP_SELECTION
  m_iSliceQpBase	int;
//#endif
  m_deblockingFilterDisable	bool;
  m_deblockingFilterOverrideFlag	bool;      //< offsets for deblocking filter inherit from PPS
  m_deblockingFilterBetaOffsetDiv2	int;    //< beta offset for deblocking filter
  m_deblockingFilterTcOffsetDiv2	int;      //< tc offset for deblocking filter
  
  m_aiNumRefIdx   [3]int;    //  for multiple reference of current slice

  m_iRefIdxOfLC		 [2][MAX_NUM_REF_LC]int;
  m_eListIdFromIdxOfLC	[MAX_NUM_REF_LC]int;
  m_iRefIdxFromIdxOfLC	[MAX_NUM_REF_LC]int;
  m_iRefIdxOfL1FromRefIdxOfL0	[MAX_NUM_REF_LC]int;
  m_iRefIdxOfL0FromRefIdxOfL1	[MAX_NUM_REF_LC]int;
  m_bRefPicListModificationFlagLC	bool;
  m_bRefPicListCombinationFlag		bool;

  m_bCheckLDC	bool;

  //  Data
  m_iSliceQpDelta	int;
  m_iSliceQpDeltaCb	int;
  m_iSliceQpDeltaCr	int;
  m_apcRefPicList [][2][MAX_NUM_REF+1]TComPic;
  m_aiRefPOCList    [2][MAX_NUM_REF+1]int;
  m_iDepth	int;
  
  // referenced slice?
  m_bRefenced	bool;
  
  // access channel
  m_pcVPS	*TComVPS;
  m_pcSPS	*TComSPS;
  m_pcPPS	*TComPPS;
  m_pcPic	*TComPic;
//#if ADAPTIVE_QP_SELECTION
  m_pcTrQuant	*TComTrQuant;
//#endif  
  m_colFromL0Flag	uint;  // collocated picture from List0 flag
  
  m_colRefIdx	uint;
  m_maxNumMergeCand	uint;


//#if SAO_CHROMA_LAMBDA
  m_dLambdaLuma		float64;
  m_dLambdaChroma	float64;
//#else
//  Double      m_dLambda;
//#endif

  m_abEqualRef  [2][MAX_NUM_REF][MAX_NUM_REF]bool;
  
  m_bNoBackPredFlag	bool;
  m_uiTLayer	uint;
  m_bTLayerSwitchingFlag	bool;

  m_uiSliceMode	uint;
  m_uiSliceArgument	uint;
  m_uiSliceCurStartCUAddr	uint;
  m_uiSliceCurEndCUAddr	uint;
  m_uiSliceIdx	uint;
  m_uiDependentSliceMode	uint;
  m_uiDependentSliceArgument	uint;
  m_uiDependentSliceCurStartCUAddr	uint;
  m_uiDependentSliceCurEndCUAddr	uint;
  m_bNextSlice	bool;
  m_bNextDependentSlice	bool;
  m_uiSliceBits	uint;
  m_uiDependentSliceCounter	uint;
  m_bFinalized	bool;

  m_weightPredTable	[2][MAX_NUM_REF][3]wpScalingParam; // [REF_PIC_LIST_0 or REF_PIC_LIST_1][refIdx][0:Y, 1:U, 2:V]
  m_weightACDCParam	[3]wpACDCParam;                 // [0:Y, 1:U, 2:V]

  m_tileByteLocation	*list.List;
  m_uiTileOffstForMultES	uint;

  m_puiSubstreamSizes		*uint;
  m_scalingList	*TComScalingList;                 //!< pointer of quantization matrix
  m_cabacInitFlag	bool; 

  m_bLMvdL1Zero	bool;
  m_numEntryPointOffsets	int;
  m_temporalLayerNonReferenceFlag	bool;
  m_LFCrossSliceBoundaryFlag	bool;

  m_enableTMVPFlag	bool;
}

//public:
func NewTComSlice() *TComSlice{
	return &TComSlice{}
}

func (this *TComSlice)  initSlice       (){
}

func (this *TComSlice)  setVPS          ( pcVPS *TComVPS) { 
	m_pcVPS = pcVPS; 
}
func (this *TComSlice)  getVPS          () *TComVPS { 
	return m_pcVPS; 
}
func (this *TComSlice)  setSPS          ( pcSPS *TComSPS) { 
	m_pcSPS = pcSPS; 
}
func (this *TComSlice)  getSPS          () *TComVPS { 
	return m_pcSPS; 
}
  
func (this *TComSlice)  setPPS          ( pcPPS *TComPPS)         { 
	//assert(pcPPS!=NULL); 
	m_pcPPS = pcPPS; 
	m_iPPSId = pcPPS.getPPSId(); 
}
func (this *TComSlice)  getPPS          () *TComVPS { 
	return m_pcPPS; 
}

//#if ADAPTIVE_QP_SELECTION
func (this *TComSlice)  setTrQuant          ( pcTrQuant *TComTrQuant) { 
	m_pcTrQuant = pcTrQuant; 
}
func (this *TComSlice)  getTrQuant          () *TComTrQuant { 
	return m_pcTrQuant; 
}
//#endif

func (this *TComSlice)  setPPSId        ( PPSId int )         { 
	m_iPPSId = PPSId; 
}
func (this *TComSlice)  getPPSId        () int { 
	return m_iPPSId; 
}
func (this *TComSlice)  setPicOutputFlag( b bool)         { 
	m_PicOutputFlag = b;    
}
func (this *TComSlice)  getPicOutputFlag() bool                { 
	return m_PicOutputFlag; 
}
func (this *TComSlice)  setSaoEnabledFlag(s bool) {
	m_saoEnabledFlag =s; 
}
func (this *TComSlice)  getSaoEnabledFlag() bool { 
	return m_saoEnabledFlag; 
}
func (this *TComSlice)  setSaoEnabledFlagChroma(s bool) {
	m_saoEnabledFlagChroma = s; 
}       //!< set SAO Cb&Cr enabled flag
func (this *TComSlice)  getSaoEnabledFlagChroma() bool { 
	return m_saoEnabledFlagChroma; 
}        //!< get SAO Cb&Cr enabled flag
func (this *TComSlice)  setRPS          ( pcRPS *TComReferencePictureSet) { 
	m_pcRPS = pcRPS; 
}
func (this *TComSlice)  getRPS          () *TComReferencePictureSet{ 
	return m_pcRPS; 
}
func (this *TComSlice)  getLocalRPS     () *TComReferencePictureSet{ 
	return &m_LocalRPS; 
}

func (this *TComSlice)  setRPSidx          ( iBDidx int ) { 
	m_iBDidx = iBDidx; 
}
func (this *TComSlice)  getRPSidx          () int{ 
	return m_iBDidx; 
}
func (this *TComSlice)  setCombinationBDidx          ( iCombinationBDidx int ) {
	m_iCombinationBDidx = iCombinationBDidx; 
}
func (this *TComSlice)  getCombinationBDidx          () int { 
	return m_iCombinationBDidx; 
}
func (this *TComSlice)  setCombineWithReferenceFlag          ( bCombineWithReferenceFlag bool) { 
	m_bCombineWithReferenceFlag = bCombineWithReferenceFlag; 
}
func (this *TComSlice)  getCombineWithReferenceFlag          () bool { 
	return m_bCombineWithReferenceFlag; 
}
func (this *TComSlice)  getPrevPOC      ()                     int     { 
	return  m_prevPOC;      
}
func (this *TComSlice)  getRefPicListModification() *TComRefPicListModification { 
	return &m_RefPicListModification; 
}
func (this *TComSlice)  setLastIDR(iIDRPOC int)                       { 
	m_iLastIDR = iIDRPOC; 
}
func (this *TComSlice)  getLastIDR()            int                      { 
	return m_iLastIDR; 
}
func (this *TComSlice)  getSliceType    ()      SliceType                 { 
	return  m_eSliceType;         
}
func (this *TComSlice)  getPOC          ()      int                    { 
	return  m_iPOC;           
}
func (this *TComSlice)  getSliceQp      ()      int                    { 
	return  m_iSliceQp;           
}
func (this *TComSlice)  getDependentSliceFlag() bool               { 
	return m_dependentSliceFlag; 
}
func (this *TComSlice)  setDependentSliceFlag(val bool)             { 
	m_dependentSliceFlag = val; 
}
//#if ADAPTIVE_QP_SELECTION
func (this *TComSlice)  getSliceQpBase  ()       int                   { 
	return  m_iSliceQpBase;       
}
//#endif
func (this *TComSlice)  getSliceQpDelta ()       int                   { 
	return  m_iSliceQpDelta;      
}
func (this *TComSlice)  getSliceQpDeltaCb ()     int                  { 
	return  m_iSliceQpDeltaCb;      
}
func (this *TComSlice)  getSliceQpDeltaCr ()     int                  { 
	return  m_iSliceQpDeltaCr;      
}
func (this *TComSlice)  getDeblockingFilterDisable()  bool              { 
	return  m_deblockingFilterDisable; 
}
func (this *TComSlice)  getDeblockingFilterOverrideFlag() bool          { 
	return  m_deblockingFilterOverrideFlag; 
}
func (this *TComSlice)  getDeblockingFilterBetaOffsetDiv2()  int       { 
	return  m_deblockingFilterBetaOffsetDiv2; 
}
func (this *TComSlice)  getDeblockingFilterTcOffsetDiv2()    int       { 
	return  m_deblockingFilterTcOffsetDiv2; 
}

func (this *TComSlice)  getNumRefIdx        ( e RefPicList )   int             { 
	return  m_aiNumRefIdx[e];             
}
func (this *TComSlice)  getPic              ()                *TComPic              { 
	return  m_pcPic;                      
}
func (this *TComSlice)  getRefPic           ( e RefPicList, iRefIdx int ) *TComPic   { 
	return  m_apcRefPicList[e][iRefIdx];  
}
func (this *TComSlice)  getRefPOC           ( e RefPicList, iRefIdx int )  int  { 
	return  m_aiRefPOCList[e][iRefIdx];   
}
func (this *TComSlice)  getDepth            ()               int               { 
	return  m_iDepth;                     
}
func (this *TComSlice)  getColFromL0Flag    ()               uint               { 
	return  m_colFromL0Flag;              
}
func (this *TComSlice)  getColRefIdx        ()               uint               { 
	return  m_colRefIdx;                  
}
func (this *TComSlice)  checkColRefIdx      ( curSliceIdx uint, pic *TComPic){
}
func (this *TComSlice)  getCheckLDC     ()               bool                   { 
	return m_bCheckLDC; 
}
func (this *TComSlice)  getMvdL1ZeroFlag ()              bool                    { 
	return m_bLMvdL1Zero;    
}
func (this *TComSlice)  getNumRpsCurrTempList()			int{
}
func (this *TComSlice)  getRefIdxOfLC       ( e RefPicList, iRefIdx int ) int    { 
	return m_iRefIdxOfLC[e][iRefIdx];           
}
func (this *TComSlice)  getListIdFromIdxOfLC(iRefIdx int)       int            { 
	return m_eListIdFromIdxOfLC[iRefIdx];       
}
func (this *TComSlice)  getRefIdxFromIdxOfLC(iRefIdx int)       int            { 
	return m_iRefIdxFromIdxOfLC[iRefIdx];       
}

func (this *TComSlice)  getRefIdxOfL0FromRefIdxOfL1(iRefIdx int) int           { 
	return m_iRefIdxOfL0FromRefIdxOfL1[iRefIdx];
}
func (this *TComSlice)  getRefIdxOfL1FromRefIdxOfL0(iRefIdx int) int           { 
	return m_iRefIdxOfL1FromRefIdxOfL0[iRefIdx];
}
func (this *TComSlice)  getRefPicListModificationFlagLC()        bool           {
	return m_bRefPicListModificationFlagLC;
}
func (this *TComSlice)  setRefPicListModificationFlagLC(bflag bool)         {
	m_bRefPicListModificationFlagLC=bflag;
}     
func (this *TComSlice)  getRefPicListCombinationFlag()          bool            {
	return m_bRefPicListCombinationFlag;
}
func (this *TComSlice)  setRefPicListCombinationFlag(bflag bool)            {
	m_bRefPicListCombinationFlag=bflag;
}   
func (this *TComSlice)  setReferenced(b bool)                               { 
	m_bRefenced = b; 
}
func (this *TComSlice)  isReferenced()                      bool                { 
	return m_bRefenced; 
}
func (this *TComSlice)  setPOC              ( i int )                       { 
	m_iPOC              = i; 
	if getTLayer()==0 {
		m_prevPOC=i; 
	}
}
func (this *TComSlice)  setNalUnitType      ( e NalUnitType)               { 
	m_eNalUnitType      = e;      
}
func (this *TComSlice)  getNalUnitType    ()            NalUnitType                  { 
	return m_eNalUnitType;        
}
func (this *TComSlice)  getRapPicFlag       ()	bool{
}
func (this *TComSlice)  getIdrPicFlag       ()    bool                          { 
	return getNalUnitType() == NAL_UNIT_CODED_SLICE_IDR || getNalUnitType() == NAL_UNIT_CODED_SLICE_IDR_N_LP; 
}
func (this *TComSlice)  checkCRA(pReferencePictureSet *TComReferencePictureSet, pocCRA *int, prevRAPisBLA *bool, rcListPic *list.List){
}
func (this *TComSlice)  decodingRefreshMarking(pocCRA *int, bRefreshPending *bool, rcListPic *list.List){
}
func (this *TComSlice)  setSliceType        ( e SliceType)                 { 
	m_eSliceType        = e;      
}
func (this *TComSlice)  setSliceQp          ( i int)                       { 
	m_iSliceQp          = i;      
}
//#if ADAPTIVE_QP_SELECTION
func (this *TComSlice)  setSliceQpBase      ( i int)                       { 
	m_iSliceQpBase      = i;      
}
//#endif
func (this *TComSlice)  setSliceQpDelta     ( i int )                       { 
	m_iSliceQpDelta     = i;      
}
func (this *TComSlice)  setSliceQpDeltaCb   ( i int )                       { 
	m_iSliceQpDeltaCb   = i;      
}
func (this *TComSlice)  setSliceQpDeltaCr   ( i int )                       { 
	m_iSliceQpDeltaCr   = i;     
}
func (this *TComSlice)  setDeblockingFilterDisable( b bool)                { 
	m_deblockingFilterDisable= b;      
}
func (this *TComSlice)  setDeblockingFilterOverrideFlag( b bool)           { 
	m_deblockingFilterOverrideFlag = b; 
}
func (this *TComSlice)  setDeblockingFilterBetaOffsetDiv2( i int)          { 
	m_deblockingFilterBetaOffsetDiv2 = i; 
}
func (this *TComSlice)  setDeblockingFilterTcOffsetDiv2( i int)            { 
	m_deblockingFilterTcOffsetDiv2 = i; 
}
  
func (this *TComSlice)  setRefPic           ( p *TComPic, e RefPicList, iRefIdx int ) { 
	m_apcRefPicList[e][iRefIdx] = p; 
}
func (this *TComSlice)  setRefPOC           ( i int, e RefPicList, iRefIdx int ) { 
	m_aiRefPOCList[e][iRefIdx] = i; 
}
func (this *TComSlice)  setNumRefIdx        ( e RefPicList, i int )         { 
	m_aiNumRefIdx[e]    = i;      
}
func (this *TComSlice)  setPic              ( p *TComPic )                  { 
	m_pcPic             = p;      
}
func (this *TComSlice)  setDepth            ( iDepth int)                  { 
	m_iDepth            = iDepth; 
}
  
func (this *TComSlice)  setRefPicList       ( rcListPic *list.List){
}
func (this *TComSlice)  setRefPOCList       (){
}
func (this *TComSlice)  setColFromL0Flag    ( colFromL0 uint ) { 
	m_colFromL0Flag = colFromL0; 
}
func (this *TComSlice)  setColRefIdx        ( refIdx uint) { 
	m_colRefIdx = refIdx; 
}
func (this *TComSlice)  setCheckLDC         ( b bool)                      { 
	m_bCheckLDC = b; 
}
func (this *TComSlice)  setMvdL1ZeroFlag    ( b bool)                       { 
	m_bLMvdL1Zero = b; 
}

func (this *TComSlice)  isIntra         ()     bool                     { 
	return  m_eSliceType == I_SLICE;  
}
func (this *TComSlice)  isInterB        ()     bool                     { 
	return  m_eSliceType == B_SLICE;  
}
func (this *TComSlice)  isInterP        ()     bool                     { 
	return  m_eSliceType == P_SLICE;  
}
  
//#if SAO_CHROMA_LAMBDA  
func (this *TComSlice)  setLambda( d, e float64 ) { 
	m_dLambdaLuma = d; 
	m_dLambdaChroma = e;
}
func (this *TComSlice)  getLambdaLuma()		float64 { 
	return m_dLambdaLuma;        
}
func (this *TComSlice)  getLambdaChroma() 	float64 { 
	return m_dLambdaChroma;        
}
//#else
//  Void      setLambda( Double d ) { m_dLambda = d; }
//  Double    getLambda() { return m_dLambda;        }
//#endif
  
func (this *TComSlice)  initEqualRef(){
}
func (this *TComSlice)  isEqualRef  ( e RefPicList, iRefIdx1 int, iRefIdx2 int ) bool{
    if iRefIdx1 < 0 || iRefIdx2 < 0 {
    	return false;
    }
    
    return m_abEqualRef[e][iRefIdx1][iRefIdx2];
}
  
func (this *TComSlice)  setEqualRef( e RefPicList, iRefIdx1 int, iRefIdx2 int, b bool){
    m_abEqualRef[e][iRefIdx1][iRefIdx2] = b;
    m_abEqualRef[e][iRefIdx2][iRefIdx1] = b;
}
  
func (this *TComSlice)  sortPicList         ( rcListPic *list.List){
}
  
func (this *TComSlice)  getNoBackPredFlag() bool{ 
	return m_bNoBackPredFlag; 
}
func (this *TComSlice)  setNoBackPredFlag( b bool ) { 
	m_bNoBackPredFlag = b; 
}
func (this *TComSlice)  generateCombinedList       (){
}

func (this *TComSlice)  getTLayer             ()       uint                     { 
	return m_uiTLayer;                      
}
func (this *TComSlice)  setTLayer             ( uiTLayer uint)             { 
	m_uiTLayer = uiTLayer;                  
}

func (this *TComSlice)  setTLayerInfo( uiTLayer uint ){
}
func (this *TComSlice)  decodingMarking( rcListPic *list.List, iGOPSIze int, iMaxRefPicNum *int ){
}
func (this *TComSlice)  applyReferencePictureSet( rcListPic *list.List, RPSList *TComReferencePictureSet){
}
func (this *TComSlice)  isTemporalLayerSwitchingPoint( rcListPic *list.List, RPSList *TComReferencePictureSet) bool{
}
func (this *TComSlice)  isStepwiseTemporalLayerSwitchingPointCandidate( rcListPic *list.List, RPSList *TComReferencePictureSet) bool{
}
func (this *TComSlice)  checkThatAllRefPicsAreAvailable( rcListPic *list.List, pReferencePictureSet *TComReferencePictureSet, printErrors bool, pocRandomAccess int) int{
}
func (this *TComSlice)  createExplicitReferencePictureSetFromReference( rcListPic *list.List, pReferencePictureSet *TComReferencePictureSet){
}

func (this *TComSlice)  setMaxNumMergeCand               (val uint)         { 
	m_maxNumMergeCand = val;                    
}
func (this *TComSlice)  getMaxNumMergeCand               ()      uint            { 
	return m_maxNumMergeCand;                   
}

func (this *TComSlice)  setSliceMode                     (uiMode uint)     { 
	m_uiSliceMode = uiMode;                     
}
func (this *TComSlice)  getSliceMode                     ()     uint             { 
	return m_uiSliceMode;                       
}
func (this *TComSlice)  setSliceArgument                 (uiArgument uint) { 
	m_uiSliceArgument = uiArgument;             
}
func (this *TComSlice)  getSliceArgument                 ()     uint             { 
	return m_uiSliceArgument;                   
}
func (this *TComSlice)  setSliceCurStartCUAddr           (uiAddr uint)     { m_uiSliceCurStartCUAddr = uiAddr;           }
func (this *TComSlice)  getSliceCurStartCUAddr           ()     uint             { return m_uiSliceCurStartCUAddr;             }
func (this *TComSlice)  setSliceCurEndCUAddr             (uiAddr uint)     { m_uiSliceCurEndCUAddr = uiAddr;             }
func (this *TComSlice)  getSliceCurEndCUAddr             ()     uint             { return m_uiSliceCurEndCUAddr;               }
func (this *TComSlice)  setSliceIdx                      (i uint)           { m_uiSliceIdx = i;                           }
func (this *TComSlice)  getSliceIdx                      ()     uint             { return  m_uiSliceIdx;                       }
func (this *TComSlice)  copySliceInfo                    (pcSliceSrc *TComSlice);
func (this *TComSlice)  setDependentSliceMode              ( uiMode uint )     { m_uiDependentSliceMode = uiMode;              }
func (this *TComSlice)  getDependentSliceMode              ()   uint               { return m_uiDependentSliceMode;                }
func (this *TComSlice)  setDependentSliceArgument          ( uiArgument uint ) { m_uiDependentSliceArgument = uiArgument;      }
func (this *TComSlice)  getDependentSliceArgument          ()   uint               { return m_uiDependentSliceArgument;            }
func (this *TComSlice)  setDependentSliceCurStartCUAddr    ( uiAddr uint )     { m_uiDependentSliceCurStartCUAddr = uiAddr;    }
func (this *TComSlice)  getDependentSliceCurStartCUAddr    ()   uint               { return m_uiDependentSliceCurStartCUAddr;      }
func (this *TComSlice)  setDependentSliceCurEndCUAddr      ( uiAddr uint )     { m_uiDependentSliceCurEndCUAddr = uiAddr;      }
func (this *TComSlice)  getDependentSliceCurEndCUAddr      ()   uint               { return m_uiDependentSliceCurEndCUAddr;        }
func (this *TComSlice)  setNextSlice                     ( b bool )          { m_bNextSlice = b;                           }
func (this *TComSlice)  isNextSlice                      ()     bool             { return m_bNextSlice;                        }
func (this *TComSlice)  setNextDependentSlice              ( b bool )          { m_bNextDependentSlice = b;                    }
func (this *TComSlice)  isNextDependentSlice               ()   bool               { return m_bNextDependentSlice;                 }
func (this *TComSlice)  setSliceBits                     ( uiVal uint )      { m_uiSliceBits = uiVal;                      }
func (this *TComSlice)  getSliceBits                     ()     uint             { return m_uiSliceBits;                       }  
func (this *TComSlice)  setDependentSliceCounter           (  uiVal uint)      { m_uiDependentSliceCounter = uiVal;            }
func (this *TComSlice)  getDependentSliceCounter           ()   uint               { return m_uiDependentSliceCounter;             }
func (this *TComSlice)  setFinalized                     ( uiVal bool)      { m_bFinalized = uiVal;                       }
func (this *TComSlice)  getFinalized                     ()     bool             { return m_bFinalized;                        }
func (this *TComSlice)  setWpScaling    ( wp [2][MAX_NUM_REF][3]wpScalingParam ) { memcpy(m_weightPredTable, wp, sizeof(wpScalingParam)*2*MAX_NUM_REF*3); }
func (this *TComSlice)  getWpScaling    ( e RefPicList, iRefIdx int, wp *wpScalingParam);

func (this *TComSlice)  resetWpScaling   (  wp [2][MAX_NUM_REF][3]wpScalingParam);
func (this *TComSlice)  initWpScaling1   (  wp [2][MAX_NUM_REF][3]wpScalingParam);
func (this *TComSlice)  initWpScaling   (){
}
func (this *TComSlice)  applyWP   () bool { 
	return( (m_eSliceType==P_SLICE && m_pcPPS->getUseWP()) || (m_eSliceType==B_SLICE && m_pcPPS->getWPBiPred()) ); 
}

func (this *TComSlice)  setWpAcDcParam  (  wp [3]wpACDCParam ) { memcpy(m_weightACDCParam, wp, sizeof(wpACDCParam)*3); }
func (this *TComSlice)  getWpAcDcParam  (  wp *wpACDCParam );
func (this *TComSlice)  initWpAcDcParam ();
  
func (this *TComSlice)  setTileLocationCount             ( cnt uint)               { return m_tileByteLocation.resize(cnt);    }
func (this *TComSlice)  getTileLocationCount             ()    uint                     { return (UInt) m_tileByteLocation.size();  }
func (this *TComSlice)  setTileLocation                  ( idx int, location uint ) { assert (idx<m_tileByteLocation.size());
                                                                m_tileByteLocation[idx] = location;       }
func (this *TComSlice)  addTileLocation                  ( location uint )          { m_tileByteLocation.push_back(location);   }
func (this *TComSlice)  getTileLocation                  ( idx int)   uint             { return m_tileByteLocation[idx];           }

func (this *TComSlice)  setTileOffstForMultES            ( uiOffset uint)      { m_uiTileOffstForMultES = uiOffset;        }
func (this *TComSlice)  getTileOffstForMultES            ()    uint                { return m_uiTileOffstForMultES;            }
func (this *TComSlice)  allocSubstreamSizes              ( uiNumSubstreams uint);
func (this *TComSlice)  getSubstreamSizes               ()    *uint              { return m_puiSubstreamSizes; }
func (this *TComSlice)  setScalingList              	( scalingList *TComScalingList ) { m_scalingList = scalingList; }
func (this *TComSlice)  getScalingList 					()        *TComScalingList                        { return m_scalingList; }
func (this *TComSlice)  setDefaultScalingList       ();
func (this *TComSlice)  checkDefaultScalingList     () bool{
}
func (this *TComSlice)  setCabacInitFlag  ( val bool ) { m_cabacInitFlag = val;      }  //!< set CABAC initial flag 
func (this *TComSlice)  getCabacInitFlag  ()     bool      { return m_cabacInitFlag;     }  //!< get CABAC initial flag 
func (this *TComSlice)  setNumEntryPointOffsets(val int)  { m_numEntryPointOffsets = val;     }
func (this *TComSlice)  getNumEntryPointOffsets()   int      { return m_numEntryPointOffsets;    }
func (this *TComSlice)  getTemporalLayerNonReferenceFlag()  bool     { return m_temporalLayerNonReferenceFlag;}
func (this *TComSlice)  setTemporalLayerNonReferenceFlag(x bool) { m_temporalLayerNonReferenceFlag = x;}
func (this *TComSlice)  setLFCrossSliceBoundaryFlag     ( val bool )    { m_LFCrossSliceBoundaryFlag = val; }
func (this *TComSlice)  getLFCrossSliceBoundaryFlag     ()     bool           { return m_LFCrossSliceBoundaryFlag;} 

func (this *TComSlice)  setEnableTMVPFlag     ( b bool)    { m_enableTMVPFlag = b; }
func (this *TComSlice)  getEnableTMVPFlag     ()      bool        { return m_enableTMVPFlag;}

//protected:
func (this *TComSlice)  xGetRefPic  (rcListPic *list.List, poc int) *TComPic{
}
func (this *TComSlice)  xGetLongTermRefPic  (rcListPic *list.List, poc int) *TComPic{
}
//};// END CLASS DEFINITION TComSlice

/*
template <class T> class ParameterSetMap
{
public:
  ParameterSetMap(Int maxId)
  :m_maxId (maxId)
  {}

  ~ParameterSetMap()
  {
    for (typename std::map<Int,T *>::iterator i = m_paramsetMap.begin(); i!= m_paramsetMap.end(); i++)
    {
      delete (*i).second;
    }
  }

  Void storePS(Int psId, T *ps)
  {
    assert ( psId < m_maxId );
    if ( m_paramsetMap.find(psId) != m_paramsetMap.end() )
    {
      delete m_paramsetMap[psId];
    }
    m_paramsetMap[psId] = ps; 
  }

  Void mergePSList(ParameterSetMap<T> &rPsList)
  {
    for (typename std::map<Int,T *>::iterator i = rPsList.m_paramsetMap.begin(); i!= rPsList.m_paramsetMap.end(); i++)
    {
      storePS(i->first, i->second);
    }
    rPsList.m_paramsetMap.clear();
  }


  T* getPS(Int psId)
  {
    return ( m_paramsetMap.find(psId) == m_paramsetMap.end() ) ? NULL : m_paramsetMap[psId];
  }

  T* getFirstPS()
  {
    return (m_paramsetMap.begin() == m_paramsetMap.end() ) ? NULL : m_paramsetMap.begin()->second;
  }

private:
  std::map<Int,T *> m_paramsetMap;
  Int               m_maxId;
};

class ParameterSetManager
{
public:
  ParameterSetManager();
  virtual ~ParameterSetManager();

  //! store sequence parameter set and take ownership of it 
  Void storeVPS(TComVPS *vps) { m_vpsMap.storePS( vps->getVPSId(), vps); };
  //! get pointer to existing video parameter set  
  TComVPS* getVPS(Int vpsId)  { return m_vpsMap.getPS(vpsId); };
  TComVPS* getFirstVPS()      { return m_vpsMap.getFirstPS(); };
  
  //! store sequence parameter set and take ownership of it 
  Void storeSPS(TComSPS *sps) { m_spsMap.storePS( sps->getSPSId(), sps); };
  //! get pointer to existing sequence parameter set  
  TComSPS* getSPS(Int spsId)  { return m_spsMap.getPS(spsId); };
  TComSPS* getFirstSPS()      { return m_spsMap.getFirstPS(); };

  //! store picture parameter set and take ownership of it 
  Void storePPS(TComPPS *pps) { m_ppsMap.storePS( pps->getPPSId(), pps); };
  //! get pointer to existing picture parameter set  
  TComPPS* getPPS(Int ppsId)  { return m_ppsMap.getPS(ppsId); };
  TComPPS* getFirstPPS()      { return m_ppsMap.getFirstPS(); };

protected:
  
  ParameterSetMap<TComVPS> m_vpsMap;
  ParameterSetMap<TComSPS> m_spsMap; 
  ParameterSetMap<TComPPS> m_ppsMap;
};
*/
