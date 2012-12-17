package TLibCommon

import (
	"container/list"
)


// ====================================================================================================================
// Non-deblocking in-loop filter processing block data structure
// ====================================================================================================================

/// Non-deblocking filter processing block border tag
type NDBFBlockBorderTag	 uint8
const (//enum NDBFBlockBorderTag
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
type NDBFBlockInfo struct{
  tileID		int;   //!< tile ID
  sliceID		int;  //!< slice ID
  startSU		uint;  //!< starting SU z-scan address in LCU
  endSU		uint;    //!< ending SU z-scan address in LCU
  widthSU		uint;  //!< number of SUs in width
  heightSU	uint; //!< number of SUs in height
  posX		uint;     //!< top-left X coordinate in picture
  posY		uint;     //!< top-left Y coordinate in picture
  width		uint;    //!< number of pixels in width
  height		uint;   //!< number of pixels in height
  isBorderAvailable	[NUM_SGU_BORDER]bool;  //!< the border availabilities
  allBordersAvailable				bool;

  //NDBFBlockInfo():tileID(0), sliceID(0), startSU(0), endSU(0) {} //!< constructor
  //const NDBFBlockInfo& operator= (const NDBFBlockInfo& src);  //!< "=" operator
};


// ====================================================================================================================
// Class definition
// ====================================================================================================================

/// CU data structure class
type TComDataCU struct{
//private:
  
  // -------------------------------------------------------------------------------------------------------------------
  // class pointers
  // -------------------------------------------------------------------------------------------------------------------
  
  m_pcPic			*TComPic;              ///< picture class pointer
  m_pcSlice			*TComSlice;            ///< slice header pointer
  m_pcPattern		*TComPattern;          ///< neighbour access class pointer
  
  // -------------------------------------------------------------------------------------------------------------------
  // CU description
  // -------------------------------------------------------------------------------------------------------------------
  
  m_uiCUAddr			uint;           ///< CU address in a slice
  m_uiAbsIdxInLCU		uint;      ///< absolute address in a CU. It's Z scan order
  m_uiCUPelX			uint;           ///< CU position in a pixel (X)
  m_uiCUPelY			uint;           ///< CU position in a pixel (Y)
  m_uiNumPartition		uint;     ///< total number of minimum partitions in a CU
  m_puhWidth			[]byte;           ///< array of widths
  m_puhHeight			[]byte;          ///< array of heights
  m_puhDepth			[]byte;           ///< array of depths
  m_unitSize			int;           ///< size of a "minimum partition"
  
  // -------------------------------------------------------------------------------------------------------------------
  // CU data
  // -------------------------------------------------------------------------------------------------------------------
  m_skipFlag			[]bool;           ///< array of skip flags
  m_pePartSize			[]PartSize;         ///< array of partition sizes
  m_pePredMode			[]PredMode;         ///< array of prediction modes
  m_CUTransquantBypass	[]bool;   ///< array of cu_transquant_bypass flags
  m_phQP				[]int8;               ///< array of QP values
  m_puhTrIdx			[]byte;           ///< array of transform indices
  m_puhTransformSkip  	[3][]byte;///< array of transform skipping flags
  m_puhCbf				[3][]byte;          ///< array of coded block flags (CBF)
  m_acCUMvField			[2]TComCUMvField;     ///< array of motion vectors
  m_pcTrCoeffY			[]TCoeff;         ///< transformed coefficient buffer (Y)
  m_pcTrCoeffCb			[]TCoeff;        ///< transformed coefficient buffer (Cb)
  m_pcTrCoeffCr			[]TCoeff;        ///< transformed coefficient buffer (Cr)
//#if ADAPTIVE_QP_SELECTION
  m_pcArlCoeffY			[]int;        ///< ARL coefficient buffer (Y)
  m_pcArlCoeffCb		[]int;       ///< ARL coefficient buffer (Cb)
  m_pcArlCoeffCr		[]int;       ///< ARL coefficient buffer (Cr)
  m_ArlCoeffIsAliasedAllocation	bool; ///< ARL coefficient buffer is an alias of the global buffer and must not be free()'d

  m_pcGlbArlCoeffY		[]int;     ///< ARL coefficient buffer (Y)
  m_pcGlbArlCoeffCb		[]int;    ///< ARL coefficient buffer (Cb)
  m_pcGlbArlCoeffCr		[]int;    ///< ARL coefficient buffer (Cr)
//#endif
  
  m_pcIPCMSampleY		[]Pel;      ///< PCM sample buffer (Y)
  m_pcIPCMSampleCb		[]Pel;     ///< PCM sample buffer (Cb)
  m_pcIPCMSampleCr		[]Pel;     ///< PCM sample buffer (Cr)

  m_piSliceSUMap		[]int;       ///< pointer of slice ID map
  m_vNDFBlock 			*list.List;

  // -------------------------------------------------------------------------------------------------------------------
  // neighbour access variables
  // -------------------------------------------------------------------------------------------------------------------
  
  m_pcCUAboveLeft			*TComDataCU;      ///< pointer of above-left CU
  m_pcCUAboveRight			*TComDataCU;     ///< pointer of above-right CU
  m_pcCUAbove				*TComDataCU;          ///< pointer of above CU
  m_pcCULeft				*TComDataCU;           ///< pointer of left CU
  m_apcCUColocated		 [2]*TComDataCU;  ///< pointer of temporally colocated CU's for both directions
  m_cMvFieldA				 TComMvField;          ///< motion vector of position A
  m_cMvFieldB				 TComMvField;          ///< motion vector of position B
  m_cMvFieldC				 TComMvField;          ///< motion vector of position C
  m_cMvPred					 TComMv;            ///< motion vector predictor
  
  // -------------------------------------------------------------------------------------------------------------------
  // coding tool information
  // -------------------------------------------------------------------------------------------------------------------
  
  m_pbMergeFlag				[]bool;        ///< array of merge flags
  m_puhMergeIndex			[]byte;      ///< array of merge candidate indices
//#if AMP_MRG
  m_bIsMergeAMP				bool;
//#endif
  m_puhLumaIntraDir			[]byte;    ///< array of intra directions (luma)
  m_puhChromaIntraDir		[]byte;  ///< array of intra directions (chroma)
  m_puhInterDir				[]byte;        ///< array of inter directions
  m_apiMVPIdx				[2][]int8;       ///< array of motion vector predictor candidates
  m_apiMVPNum				[2][]int8;       ///< array of number of possible motion vectors predictors
  m_pbIPCMFlag				[]bool;         ///< array of intra_pcm flags

  m_numSucIPCM				int;         ///< the number of succesive IPCM blocks associated with the current log2CUSize
  m_lastCUSucIPCMFlag		bool;  ///< True indicates that the last CU is IPCM and shares the same root as the current CU.  

  // -------------------------------------------------------------------------------------------------------------------
  // misc. variables
  // -------------------------------------------------------------------------------------------------------------------
  
  m_bDecSubCu			bool;          ///< indicates decoder-mode
  m_dTotalCost			float64;         ///< sum of partition RD costs
  m_uiTotalDistortion	uint;  ///< sum of partition distortion
  m_uiTotalBits			uint;        ///< sum of partition bits
  m_uiTotalBins			uint;       ///< sum of partition bins
  m_uiSliceStartCU			[]uint;    ///< Start CU address of current slice
  m_uiDependentSliceStartCU	[]uint; ///< Start CU address of current slice
  m_codedQP					int8;
}

//protected:
  
  /// add possible motion vector predictor candidates
func (this *TComDataCU)  xAddMVPCand           ( pInfo *AMVPInfo,  eRefPicList RefPicList,  iRefIdx int,  uiPartUnitIdx uint,  eDir MVP_DIR)bool{
	return true;
}
func (this *TComDataCU)  xAddMVPCandOrder      ( pInfo *AMVPInfo,  eRefPicList RefPicList,  iRefIdx int,  uiPartUnitIdx uint,  eDir MVP_DIR)bool{
	return true;
}

func (this *TComDataCU)  deriveRightBottomIdx        (  uiPartIdx uint, ruiPartIdxRB *uint){
}
func (this *TComDataCU)  xGetColMVP(  eRefPicList RefPicList,  uiCUAddr,  uiPartUnitIdx int, rcMv *TComMv, riRefIdx *int)bool{
	return true;
}
  
  /// compute required bits to encode MVD (used in AMVP)
func (this *TComDataCU)  xGetMvdBits           (  cMvd TComMv)uint{
	return 0;
}
func (this *TComDataCU)  xGetComponentBits     (  iVal int)uint{
	return 0;
}
  
  /// compute scaling factor from POC difference
func (this *TComDataCU)  xGetDistScaleFactor   (  iCurrPOC,  iCurrRefPOC,  iColPOC,  iColRefPOC int)int{
	return 0;
}
  
func (this *TComDataCU)  xDeriveCenterIdx(  uiPartIdx uint, ruiPartIdxCenter *uint){
}
func (this *TComDataCU)  xGetCenterCol(  uiPartIdx uint,  eRefPicList RefPicList,  iRefIdx int, pcMv *TComMv)bool{
	return true;
}


//public:

func NewTComDataCU() *TComDataCU{
	return &TComDataCU{}
}
 
  // -------------------------------------------------------------------------------------------------------------------
  // create / destroy / initialize / copy
  // -------------------------------------------------------------------------------------------------------------------
  
func (this *TComDataCU)  Create(  uiNumPartition,  uiWidth,  uiHeight uint, bDecSubCu bool,  unitSize int,
//#if ADAPTIVE_QP_SELECTION
    bGlobalRMARLBuffer bool){
//#endif  
  this.m_bDecSubCu = bDecSubCu;
  
  this.m_pcPic              = nil;
  this.m_pcSlice            = nil;
  this.m_uiNumPartition     = uiNumPartition;
  this.m_unitSize = unitSize;
  
  if !bDecSubCu {
    this.m_phQP               = make([]int8,    uiNumPartition);
    this.m_puhDepth           = make([]byte,    uiNumPartition);
    this.m_puhWidth           = make([]byte,    uiNumPartition);
    this.m_puhHeight          = make([]byte,    uiNumPartition);

    this.m_skipFlag           = make([]bool, 	uiNumPartition);

    this.m_pePartSize         = make([]PartSize,    uiNumPartition);
    for i:=uint(0); i<uiNumPartition; i++{
    	this.m_pePartSize[i] = SIZE_NONE;
    }
    
    this.m_pePredMode         = make([]PredMode,    uiNumPartition);
    this.m_CUTransquantBypass = make([]bool,    uiNumPartition);
    this.m_pbMergeFlag        = make([]bool,    uiNumPartition);
    this.m_puhMergeIndex      = make([]byte,  uiNumPartition);
    this.m_puhLumaIntraDir    = make([]byte,  uiNumPartition);
    this.m_puhChromaIntraDir  = make([]byte,  uiNumPartition);
    this.m_puhInterDir        = make([]byte,  uiNumPartition);
    
    this.m_puhTrIdx           = make([]byte,  uiNumPartition);
    this.m_puhTransformSkip[0] = make([]byte,  uiNumPartition);
    this.m_puhTransformSkip[1] = make([]byte,  uiNumPartition);
    this.m_puhTransformSkip[2] = make([]byte,  uiNumPartition);

    this.m_puhCbf[0]          = make([]byte,  uiNumPartition);
    this.m_puhCbf[1]          = make([]byte,  uiNumPartition);
    this.m_puhCbf[2]          = make([]byte,  uiNumPartition);
    
    this.m_apiMVPIdx[0]       = make([]int8,    uiNumPartition);
    this.m_apiMVPIdx[1]       = make([]int8,    uiNumPartition);
    this.m_apiMVPNum[0]       = make([]int8,    uiNumPartition);
    this.m_apiMVPNum[1]       = make([]int8,    uiNumPartition);
    for i:=uint(0); i<uiNumPartition; i++{
    	this.m_apiMVPIdx[0][i]=-1;
    	this.m_apiMVPIdx[1][i]=-1;
    }
    this.m_pcTrCoeffY         = make([]TCoeff, uiWidth*uiHeight);
    this.m_pcTrCoeffCb        = make([]TCoeff, uiWidth*uiHeight/4);
    this.m_pcTrCoeffCr        = make([]TCoeff, uiWidth*uiHeight/4);

//#if ADAPTIVE_QP_SELECTION    
    if bGlobalRMARLBuffer {
      if this.m_pcGlbArlCoeffY == nil {
        this.m_pcGlbArlCoeffY   = make([]int, uiWidth*uiHeight);
        this.m_pcGlbArlCoeffCb  = make([]int, uiWidth*uiHeight/4);
        this.m_pcGlbArlCoeffCr  = make([]int, uiWidth*uiHeight/4);
      }
      this.m_pcArlCoeffY        = this.m_pcGlbArlCoeffY;
      this.m_pcArlCoeffCb       = this.m_pcGlbArlCoeffCb;
      this.m_pcArlCoeffCr       = this.m_pcGlbArlCoeffCr;
      this.m_ArlCoeffIsAliasedAllocation = true;
    }else{
      this.m_pcArlCoeffY        = make([]int, uiWidth*uiHeight);
      this.m_pcArlCoeffCb       = make([]int, uiWidth*uiHeight/4);
      this.m_pcArlCoeffCr       = make([]int, uiWidth*uiHeight/4);
    }
//#endif
    
    this.m_pbIPCMFlag         =  make([]bool, uiNumPartition);
    this.m_pcIPCMSampleY      =  make([]Pel , uiWidth*uiHeight);
    this.m_pcIPCMSampleCb     =  make([]Pel , uiWidth*uiHeight/4);
    this.m_pcIPCMSampleCr     =  make([]Pel , uiWidth*uiHeight/4);

    this.m_acCUMvField[0].Create( uiNumPartition );
    this.m_acCUMvField[1].Create( uiNumPartition );
    
  }else{
    this.m_acCUMvField[0].SetNumPartition(uiNumPartition );
    this.m_acCUMvField[1].SetNumPartition(uiNumPartition );
  }
  
  this.m_uiSliceStartCU          = make([]uint, uiNumPartition);
  this.m_uiDependentSliceStartCU = make([]uint, uiNumPartition);
  
  // create pattern memory
  this.m_pcPattern            =  NewTComPattern();
  
  // create motion vector fields
  
  this.m_pcCUAboveLeft      = nil;
  this.m_pcCUAboveRight     = nil;
  this.m_pcCUAbove          = nil;
  this.m_pcCULeft           = nil;
  
  this.m_apcCUColocated[0]  = nil;
  this.m_apcCUColocated[1]  = nil;	
}
func (this *TComDataCU)  Destroy(){
  this.m_pcPic              = nil;
  this.m_pcSlice            = nil;
  
  if this.m_pcPattern !=nil{ 
    this.m_pcPattern = nil;
  }
  
  // encoder-side buffer free
  if !this.m_bDecSubCu {
    this.m_phQP              = nil; 
    this.m_puhDepth          = nil; 
    this.m_puhWidth          = nil; 
    this.m_puhHeight         = nil; 

    this.m_skipFlag          = nil; 

    this.m_pePartSize        = nil; 
    this.m_pePredMode        = nil; 
    this.m_CUTransquantBypass = nil;
    this.m_puhCbf[0]         = nil; 
    this.m_puhCbf[1]         = nil; 
    this.m_puhCbf[2]         = nil; 
    this.m_puhInterDir       = nil; 
    this.m_pbMergeFlag       = nil; 
    this.m_puhMergeIndex     = nil; 
    this.m_puhLumaIntraDir   = nil; 
    this.m_puhChromaIntraDir = nil; 
    this.m_puhTrIdx          = nil; 
    this.m_puhTransformSkip[0] = nil; 
    this.m_puhTransformSkip[1] = nil; 
    this.m_puhTransformSkip[2] = nil; 
    this.m_pcTrCoeffY        = nil; 
    this.m_pcTrCoeffCb       = nil; 
    this.m_pcTrCoeffCr       = nil; 
//#if ADAPTIVE_QP_SELECTION
    if !this.m_ArlCoeffIsAliasedAllocation{
      this.m_pcArlCoeffY  = nil;
      this.m_pcArlCoeffCb = nil;
      this.m_pcArlCoeffCr = nil;
    }
    this.m_pcGlbArlCoeffY    = nil;
    this.m_pcGlbArlCoeffCb   = nil;
    this.m_pcGlbArlCoeffCr   = nil;
//#endi
    this.m_pbIPCMFlag        = nil;
    this.m_pcIPCMSampleY     = nil;
    this.m_pcIPCMSampleCb    = nil;
    this.m_pcIPCMSampleCr    = nil;
    this.m_apiMVPIdx[0]      = nil;
    this.m_apiMVPIdx[1]      = nil;
    this.m_apiMVPNum[0]      = nil;
    this.m_apiMVPNum[1]      = nil;
    
    this.m_acCUMvField[0].Destroy();
    this.m_acCUMvField[1].Destroy();
  }
  
  this.m_pcCUAboveLeft       = nil;
  this.m_pcCUAboveRight      = nil;
  this.m_pcCUAbove           = nil;
  this.m_pcCULeft            = nil;
  
  this.m_apcCUColocated[0]   = nil;
  this.m_apcCUColocated[1]   = nil;

  this.m_uiSliceStartCU=nil;
  this.m_uiDependentSliceStartCU=nil;
}

func (this *TComDataCU)  InitCU ( pcPic *TComPic,  iCUAddr uint){
  var i int;
  
  this.m_pcPic              = pcPic;
  this.m_pcSlice            = pcPic.GetSlice(pcPic.GetCurrSliceIdx());
  this.m_uiCUAddr           = iCUAddr;
  this.m_uiCUPelX           = ( iCUAddr % pcPic.GetFrameWidthInCU() ) * G_uiMaxCUWidth;
  this.m_uiCUPelY           = ( iCUAddr / pcPic.GetFrameWidthInCU() ) * G_uiMaxCUHeight;
  this.m_uiAbsIdxInLCU      = 0;
  this.m_dTotalCost         = MAX_DOUBLE;
  this.m_uiTotalDistortion  = 0;
  this.m_uiTotalBits        = 0;
  this.m_uiTotalBins        = 0;
  this.m_uiNumPartition     = pcPic.GetNumPartInCU();
  this.m_numSucIPCM       = 0;
  this.m_lastCUSucIPCMFlag   = false;

  for i=0; i< int(pcPic.GetNumPartInCU()); i++ {
    if pcPic.GetPicSym().GetInverseCUOrderMap(int(iCUAddr))*pcPic.GetNumPartInCU()+uint(i)>=this.GetSlice().GetSliceCurStartCUAddr() {
      this.m_uiSliceStartCU[i]=this.GetSlice().GetSliceCurStartCUAddr();
    }else{
      this.m_uiSliceStartCU[i]=pcPic.GetCU(this.GetAddr()).m_uiSliceStartCU[i];
    }
  }
  for i=0; i< int(pcPic.GetNumPartInCU()); i++ {
    if pcPic.GetPicSym().GetInverseCUOrderMap(int(iCUAddr))*pcPic.GetNumPartInCU()+uint(i)>=this.GetSlice().GetDependentSliceCurStartCUAddr() {
      this.m_uiDependentSliceStartCU[i]=this.GetSlice().GetDependentSliceCurStartCUAddr();
    }else{
      this.m_uiDependentSliceStartCU[i]=pcPic.GetCU(this.GetAddr()).m_uiDependentSliceStartCU[i];
    }
  }

  partStartIdx := this.GetSlice().GetDependentSliceCurStartCUAddr() - pcPic.GetPicSym().GetInverseCUOrderMap(int(iCUAddr)) * pcPic.GetNumPartInCU();

  var ui, numElements uint;
  if partStartIdx < this.m_uiNumPartition {
  	numElements = partStartIdx;
  }else{
  	numElements = this.m_uiNumPartition;
  }
  
  for ui = 0; ui < numElements; ui++ {
    pcFrom := pcPic.GetCU(this.GetAddr());
    this.m_skipFlag[ui]   = pcFrom.GetSkipFlag1(ui);
    this.m_pePartSize[ui] = pcFrom.GetPartitionSize1(ui);
    this.m_pePredMode[ui] = pcFrom.GetPredictionMode1(ui);
    this.m_CUTransquantBypass[ui] = pcFrom.GetCUTransquantBypass1(ui);
    this.m_puhDepth[ui] = pcFrom.GetDepth1(ui);
    this.m_puhWidth  [ui] = pcFrom.GetWidth1(ui);
    this.m_puhHeight [ui] = pcFrom.GetHeight1(ui);
    this.m_puhTrIdx  [ui] = pcFrom.GetTransformIdx1(ui);
    this.m_puhTransformSkip[0][ui] = pcFrom.GetTransformSkip2(ui,TEXT_LUMA);
    this.m_puhTransformSkip[1][ui] = pcFrom.GetTransformSkip2(ui,TEXT_CHROMA_U);
    this.m_puhTransformSkip[2][ui] = pcFrom.GetTransformSkip2(ui,TEXT_CHROMA_V);
    this.m_apiMVPIdx[0][ui] = pcFrom.m_apiMVPIdx[0][ui];;
    this.m_apiMVPIdx[1][ui] = pcFrom.m_apiMVPIdx[1][ui];
    this.m_apiMVPNum[0][ui] = pcFrom.m_apiMVPNum[0][ui];
    this.m_apiMVPNum[1][ui] = pcFrom.m_apiMVPNum[1][ui];
    this.m_phQP[ui]=pcFrom.m_phQP[ui];
    this.m_pbMergeFlag[ui]=pcFrom.m_pbMergeFlag[ui];
    this.m_puhMergeIndex[ui]=pcFrom.m_puhMergeIndex[ui];
    this.m_puhLumaIntraDir[ui]=pcFrom.m_puhLumaIntraDir[ui];
    this.m_puhChromaIntraDir[ui]=pcFrom.m_puhChromaIntraDir[ui];
    this.m_puhInterDir[ui]=pcFrom.m_puhInterDir[ui];
    this.m_puhCbf[0][ui]=pcFrom.m_puhCbf[0][ui];
    this.m_puhCbf[1][ui]=pcFrom.m_puhCbf[1][ui];
    this.m_puhCbf[2][ui]=pcFrom.m_puhCbf[2][ui];
    this.m_pbIPCMFlag[ui] = pcFrom.m_pbIPCMFlag[ui];
  }
  
  var firstElement uint;
  if partStartIdx > 0 {
  	firstElement = partStartIdx;
  }else{
  	firstElement = 0;
  }
  numElements = this.m_uiNumPartition - firstElement;
  
  if numElements > 0 {
  	for i:=uint(0); i<numElements; i++{
     this.m_skipFlag           [ firstElement+i]= false;                    
     this.m_pePartSize         [ firstElement+i]= SIZE_NONE;                
     this.m_pePredMode         [ firstElement+i]= MODE_NONE;                
     this.m_CUTransquantBypass [ firstElement+i]= false;                
     this.m_puhDepth           [ firstElement+i]= 0;                        
     this.m_puhTrIdx           [ firstElement+i]= 0;                        
     this.m_puhTransformSkip[0][ firstElement+i]= 0;                        
     this.m_puhTransformSkip[1][ firstElement+i]= 0;                        
     this.m_puhTransformSkip[2][ firstElement+i]= 0;                        
     this.m_puhWidth           [ firstElement+i]= byte(G_uiMaxCUWidth);          
     this.m_puhHeight          [ firstElement+i]= byte(G_uiMaxCUHeight);          
     this.m_apiMVPIdx[0]       [ firstElement+i]= -1;                       
     this.m_apiMVPIdx[1]       [ firstElement+i]= -1;                       
     this.m_apiMVPNum[0]       [ firstElement+i]= -1;                       
     this.m_apiMVPNum[1]       [ firstElement+i]= -1;                       
     this.m_phQP               [ firstElement+i]= int8(this.GetSlice().GetSliceQp()); 
     this.m_pbMergeFlag        [ firstElement+i]= false;                  
     this.m_puhMergeIndex      [ firstElement+i]= 0;                        
     this.m_puhLumaIntraDir    [ firstElement+i]= DC_IDX;                   
     this.m_puhChromaIntraDir  [ firstElement+i]= 0;                        
     this.m_puhInterDir        [ firstElement+i]= 0;                        
     this.m_puhCbf[0]          [ firstElement+i]= 0;                        
     this.m_puhCbf[1]          [ firstElement+i]= 0;                        
     this.m_puhCbf[2]          [ firstElement+i]= 0;                        
     this.m_pbIPCMFlag         [ firstElement+i]= false;                    
    }
  }
  
  uiTmp := G_uiMaxCUWidth*G_uiMaxCUHeight;
  if 0 >= partStartIdx {
    this.m_acCUMvField[0].ClearMvField();
    this.m_acCUMvField[1].ClearMvField();
    //memSet( this.m_pcTrCoeffY , 0, sizeof( TCoeff ) * uiTmp );
//#if ADAPTIVE_QP_SELECTION
    //memSet( this.m_pcArlCoeffY , 0, sizeof( Int ) * uiTmp );  
//#endif
    //memSet( this.m_pcIPCMSampleY , 0, sizeof( Pel ) * uiTmp );
    uiTmp  >>= 2;
    //memSet( this.m_pcTrCoeffCb, 0, sizeof( TCoeff ) * uiTmp );
    //memSet( this.m_pcTrCoeffCr, 0, sizeof( TCoeff ) * uiTmp );
//#if ADAPTIVE_QP_SELECTION  
    //memSet( this.m_pcArlCoeffCb, 0, sizeof( Int ) * uiTmp );
    //memSet( this.m_pcArlCoeffCr, 0, sizeof( Int ) * uiTmp );
//#endif
    //memSet( this.m_pcIPCMSampleCb , 0, sizeof( Pel ) * uiTmp );
    //memSet( this.m_pcIPCMSampleCr , 0, sizeof( Pel ) * uiTmp );
  }else{
    pcFrom := pcPic.GetCU(this.GetAddr());
    this.m_acCUMvField[0].CopyFrom(&pcFrom.m_acCUMvField[0],this.m_uiNumPartition,0);
    this.m_acCUMvField[1].CopyFrom(&pcFrom.m_acCUMvField[1],this.m_uiNumPartition,0);
    for i:=uint(0); i<uiTmp; i++ {
      this.m_pcTrCoeffY[i]=pcFrom.m_pcTrCoeffY[i];
//#if ADAPTIVE_QP_SELECTION
      this.m_pcArlCoeffY[i]=pcFrom.m_pcArlCoeffY[i];
//#endif
      this.m_pcIPCMSampleY[i]=pcFrom.m_pcIPCMSampleY[i];
    }
    for i:=uint(0); i<(uiTmp>>2); i++ {
      this.m_pcTrCoeffCb[i]=pcFrom.m_pcTrCoeffCb[i];
      this.m_pcTrCoeffCr[i]=pcFrom.m_pcTrCoeffCr[i];
//#if ADAPTIVE_QP_SELECTION
      this.m_pcArlCoeffCb[i]=pcFrom.m_pcArlCoeffCb[i];
      this.m_pcArlCoeffCr[i]=pcFrom.m_pcArlCoeffCr[i];
//#endif
      this.m_pcIPCMSampleCb[i]=pcFrom.m_pcIPCMSampleCb[i];
      this.m_pcIPCMSampleCr[i]=pcFrom.m_pcIPCMSampleCr[i];
    }
  }

  // Setting neighbor CU
  this.m_pcCULeft        = nil;
  this.m_pcCUAbove       = nil;
  this.m_pcCUAboveLeft   = nil;
  this.m_pcCUAboveRight  = nil;

  this.m_apcCUColocated[0] = nil;
  this.m_apcCUColocated[1] = nil;

  uiWidthInCU := pcPic.GetFrameWidthInCU();
  if this.m_uiCUAddr % uiWidthInCU != 0 {
    this.m_pcCULeft = pcPic.GetCU( this.m_uiCUAddr - 1 );
  }

  if this.m_uiCUAddr / uiWidthInCU != 0 {
    this.m_pcCUAbove = pcPic.GetCU( this.m_uiCUAddr - uiWidthInCU );
  }

  if this.m_pcCULeft!=nil && this.m_pcCUAbove!=nil {
    this.m_pcCUAboveLeft = pcPic.GetCU( this.m_uiCUAddr - uiWidthInCU - 1 );
  }

  if this.m_pcCUAbove!=nil && ( (this.m_uiCUAddr%uiWidthInCU) < (uiWidthInCU-1) )  {
    this.m_pcCUAboveRight = pcPic.GetCU( this.m_uiCUAddr - uiWidthInCU + 1 );
  }

  if this.GetSlice().GetNumRefIdx( REF_PIC_LIST_0 ) > 0 {
    this.m_apcCUColocated[0] = this.GetSlice().GetRefPic( REF_PIC_LIST_0, 0).GetCU( this.m_uiCUAddr );
  }

  if this.GetSlice().GetNumRefIdx( REF_PIC_LIST_1 ) > 0 {
    this.m_apcCUColocated[1] = this.GetSlice().GetRefPic( REF_PIC_LIST_1, 0).GetCU( this.m_uiCUAddr );
  }
}
func (this *TComDataCU)  InitEstData           (  uiDepth uint,  qp int){
}
func (this *TComDataCU)  InitSubCU             ( pcCU *TComDataCU,  uiPartUnitIdx,  uiDepth uint,  qp int){
}
func (this *TComDataCU)  SetOutsideCUPart      (  uiAbsPartIdx,  uiDepth uint){
}

func (this *TComDataCU)  CopySubCU             ( pcCU *TComDataCU,  uiPartUnitIdx,  uiDepth uint){
}
func (this *TComDataCU)  CopyInterPredInfoFrom ( pcCU *TComDataCU,  uiAbsPartIdx uint,  eRefPicList RefPicList){
}
func (this *TComDataCU)  CopyPartFrom          ( pcCU *TComDataCU,  uiPartUnitIdx,  uiDepth uint){
}
  
func (this *TComDataCU)  CopyToPic1             (  uiDepth uint){
}
func (this *TComDataCU)  CopyToPic3             (  uiDepth,  uiPartIdx,  uiPartDepth uint){
}
  
  // -------------------------------------------------------------------------------------------------------------------
  // member functions for CU description
  // -------------------------------------------------------------------------------------------------------------------

func (this *TComDataCU)  GetPic                ()             *TComPic                  { 
	return this.m_pcPic;           
}
 
func (this *TComDataCU)  GetSlice              ()           *TComSlice                   { 
	return this.m_pcSlice;         
}
func (this *TComDataCU)  GetAddr               ()   uint                     { 
	return this.m_uiCUAddr;        
}
func (this *TComDataCU)  GetZorderIdxInCU      ()   uint                     { 
	return this.m_uiAbsIdxInLCU; 
}
func (this *TComDataCU)  GetSCUAddr            ()	   uint{
	return 0;
}
func (this *TComDataCU)  GetCUPelX             ()    uint                    { 
	return this.m_uiCUPelX;        
}
func (this *TComDataCU)  GetCUPelY             ()    uint                    { 
	return this.m_uiCUPelY;        }
func (this *TComDataCU)  GetPattern            ()   *TComPattern                     { 
	return this.m_pcPattern;       }
  
func (this *TComDataCU)  GetDepth              ()   []byte                     { 
	return this.m_puhDepth;        }
func (this *TComDataCU)  GetDepth1              (  uiIdx uint) byte           { 
	return this.m_puhDepth[uiIdx]; }
func (this *TComDataCU)  SetDepth              (  uiIdx uint,   uh byte) { 
	this.m_puhDepth[uiIdx] = uh;   }
  
func (this *TComDataCU)  SetDepthSubParts      (  uiDepth,  uiAbsPartIdx uint){
}
 
  // -------------------------------------------------------------------------------------------------------------------
  // member functions for CU data
  // -------------------------------------------------------------------------------------------------------------------
  
func (this *TComDataCU)  GetPartitionSize      ()     []PartSize                   { 
	return this.m_pePartSize;        
	}
func (this *TComDataCU)  GetPartitionSize1      ( uiIdx uint ) PartSize        { 
	return PartSize( this.m_pePartSize[uiIdx] ); 
}
func (this *TComDataCU)  SetPartitionSize      ( uiIdx uint, uh PartSize ){ 
	this.m_pePartSize[uiIdx] = uh;   
}
func (this *TComDataCU)  SetPartSizeSubParts   ( eMode PartSize, uiAbsPartIdx, uiDepth uint ){
}
func (this *TComDataCU)  SetCUTransquantBypassSubParts( flag bool, uiAbsPartIdx, uiDepth uint ){
}
  
func (this *TComDataCU)  GetSkipFlag            ()      []bool                  { 
	return this.m_skipFlag;          
}
func (this *TComDataCU)  GetSkipFlag1            ( idx uint)    bool            { 
	return this.m_skipFlag[idx];     
}
func (this *TComDataCU)  SetSkipFlag           (  idx uint, skip bool)     { 
	this.m_skipFlag[idx] = skip;   
}
func (this *TComDataCU)  SetSkipFlagSubParts   ( skip bool, absPartIdx, depth uint ){
}

func (this *TComDataCU)  GetPredictionMode     ()       []PredMode                 { 
	return this.m_pePredMode;        
}
func (this *TComDataCU)  GetPredictionMode1     ( uiIdx uint ) PredMode           { 
	return PredMode( this.m_pePredMode[uiIdx] ); 
}
func (this *TComDataCU)  GetCUTransquantBypass ()              []bool          { 
	return this.m_CUTransquantBypass;        
}
func (this *TComDataCU)  GetCUTransquantBypass1(  uiIdx uint )     bool        { 
	return this.m_CUTransquantBypass[uiIdx]; 
}
func (this *TComDataCU)  SetPredictionMode     ( uiIdx uint, uh PredMode ){ 
	this.m_pePredMode[uiIdx] = uh;  
}
func (this *TComDataCU)  SetPredModeSubParts   ( eMode PredMode, uiAbsPartIdx, uiDepth uint ){
}
  
func (this *TComDataCU)  GetWidth              () []byte                       { 
	return this.m_puhWidth;          
}
func (this *TComDataCU)  GetWidth1             ( uiIdx uint) byte           { 
	return this.m_puhWidth[uiIdx];   
}

func (this *TComDataCU)  SetWidth              (  uiIdx uint,   uh byte) { 
	this.m_puhWidth[uiIdx] = uh;     
}
  
func (this *TComDataCU)  GetHeight             ()  []byte                      { 
	return this.m_puhHeight;         
}
func (this *TComDataCU)  GetHeight1            (  uiIdx uint) byte           { 
	return this.m_puhHeight[uiIdx];  
}
func (this *TComDataCU)  SetHeight             (  uiIdx uint,   uh byte) { 
	this.m_puhHeight[uiIdx] = uh;    
}

func (this *TComDataCU)  SetSizeSubParts       (  uiWidth,  uiHeight,  uiAbsPartIdx,  uiDepth uint){
}
func (this *TComDataCU)  GetQP                 ()                        []int8{ 
	return this.m_phQP;              
}
func (this *TComDataCU)  GetQP1                (  uiIdx uint)            int8{ 
	return this.m_phQP[uiIdx];       
}
func (this *TComDataCU)  SetQP                 (  uiIdx int,  value int8){ 
	this.m_phQP[uiIdx] =  value;     
}
func (this *TComDataCU)  SetQPSubParts         (  qp int,    uiAbsPartIdx,  uiDepth uint ){
}
func (this *TComDataCU)  GetLastValidPartIdx   (  iAbsPartIdx int) int{
  iLastValidPartIdx := iAbsPartIdx-1;
  for iLastValidPartIdx >= 0 && this.GetPredictionMode1( uint(iLastValidPartIdx) ) == MODE_NONE {
    uiDepth := this.GetDepth1( uint(iLastValidPartIdx) );
    iLastValidPartIdx -= int(this.m_uiNumPartition>>(uiDepth<<1));
  }
  return iLastValidPartIdx;
}
func (this *TComDataCU)  GetLastCodedQP        (  uiAbsPartIdx uint) int8{
  var uiQUPartIdxMask uint;
  uiQUPartIdxMask = ^((1<<((G_uiMaxCUDepth - this.GetSlice().GetPPS().GetMaxCuDQPDepth())<<1))-1);
  iLastValidPartIdx := this.GetLastValidPartIdx( int(uiAbsPartIdx&uiQUPartIdxMask) );
  
  if uiAbsPartIdx < this.m_uiNumPartition && (this.GetSCUAddr()+uint(iLastValidPartIdx) < this.GetSliceStartCU(this.m_uiAbsIdxInLCU+uiAbsPartIdx)) {
    return int8(this.GetSlice().GetSliceQp());
  }else if iLastValidPartIdx >= 0 {
    return this.GetQP1( uint(iLastValidPartIdx) );
  }else if this.GetZorderIdxInCU() > 0 {
      return this.GetPic().GetCU( this.GetAddr() ).GetLastCodedQP( this.GetZorderIdxInCU() );
  }else if this.GetPic().GetPicSym().GetInverseCUOrderMap(int(this.GetAddr())) > 0 &&
           this.GetPic().GetPicSym().GetTileIdxMap(int(this.GetAddr())) == this.GetPic().GetPicSym().GetTileIdxMap(int(this.GetPic().GetPicSym().GetCUOrderMap(int(this.GetPic().GetPicSym().GetInverseCUOrderMap(int(this.GetAddr())))-1))) &&
           !( this.GetSlice().GetPPS().GetEntropyCodingSyncEnabledFlag() && this.GetAddr() % this.GetPic().GetFrameWidthInCU() == 0 ) {
      return this.GetPic().GetCU( this.GetPic().GetPicSym().GetCUOrderMap(int(this.GetPic().GetPicSym().GetInverseCUOrderMap(int(this.GetAddr()))-1)) ).GetLastCodedQP( this.GetPic().GetNumPartInCU() );
  }

  return int8(this.GetSlice().GetSliceQp());
}
func (this *TComDataCU)  SetQPSubCUs           (  qp int, pcCU *TComDataCU,  absPartIdx,  depth uint, foundNonZeroCbf *bool){
}
func (this *TComDataCU)  SetCodedQP            (  qp int8)               { 
	this.m_codedQP = qp;             
}
func (this *TComDataCU)  GetCodedQP            ()                        int8{ 
	return this.m_codedQP;           
}

func (this *TComDataCU)  IsLosslessCoded( absPartIdx uint) bool{
  return (this.GetSlice().GetPPS().GetTransquantBypassEnableFlag() && this.GetCUTransquantBypass1 (absPartIdx));
}
  
func (this *TComDataCU)  GetTransformIdx       ()                  []byte      { 
	return this.m_puhTrIdx;          
}
func (this *TComDataCU)  GetTransformIdx1      (  uiIdx uint)        byte      { 
	return this.m_puhTrIdx[uiIdx];   
}
func (this *TComDataCU)  SetTrIdxSubParts      (  uiTrIdx,  uiAbsPartIdx,  uiDepth uint){
}

func (this *TComDataCU)  GetTransformSkip1      (  eType TextType)    []byte{ 
	return this.m_puhTransformSkip[G_aucConvertTxtTypeToIdx[eType]];
}
func (this *TComDataCU)  GetTransformSkip2      (  uiIdx uint, eType TextType)   byte { 
	return this.m_puhTransformSkip[G_aucConvertTxtTypeToIdx[eType]][uiIdx];
}
func (this *TComDataCU)  SetTransformSkipSubParts4  (  useTransformSkip uint,  eType TextType,  uiAbsPartIdx,  uiDepth uint){ 
}
func (this *TComDataCU)  SetTransformSkipSubParts5  (  useTransformSkipY,  useTransformSkipU,  useTransformSkipV,  uiAbsPartIdx,  uiDepth uint ){
}

func (this *TComDataCU)  GetQuadtreeTULog2MinSizeInCU(  absPartIdx uint) uint{
  log2CbSize := uint(G_aucConvertToBit[this.GetWidth1( absPartIdx )] + 2);
  partSize  := this.GetPartitionSize1( absPartIdx );
  var quadtreeTUMaxDepth uint;
  if this.GetPredictionMode1( absPartIdx ) == MODE_INTRA{
  	quadtreeTUMaxDepth = this.m_pcSlice.GetSPS().GetQuadtreeTUMaxDepthIntra();
  }else{
  	quadtreeTUMaxDepth = this.m_pcSlice.GetSPS().GetQuadtreeTUMaxDepthInter(); 
  }
  var intraSplitFlag uint
  if this.GetPredictionMode1( absPartIdx ) == MODE_INTRA && partSize == SIZE_NxN {
  	intraSplitFlag = 1;
  }else{
  	intraSplitFlag = 0;
  }

  interSplitFlag := uint(B2U((quadtreeTUMaxDepth == 1) && (this.GetPredictionMode1( absPartIdx ) == MODE_INTER) && (partSize != SIZE_2Nx2N) ));
  
  log2MinTUSizeInCU := uint(0);
  if log2CbSize < (uint(this.m_pcSlice.GetSPS().GetQuadtreeTULog2MinSize()) + quadtreeTUMaxDepth - 1 + interSplitFlag + intraSplitFlag) {
    // when fully making use of signaled TUMaxDepth + inter/intraSplitFlag, resulting luma TB size is < QuadtreeTULog2MinSize
    log2MinTUSizeInCU = this.m_pcSlice.GetSPS().GetQuadtreeTULog2MinSize();
  }else{
    // when fully making use of signaled TUMaxDepth + inter/intraSplitFlag, resulting luma TB size is still >= QuadtreeTULog2MinSize
    log2MinTUSizeInCU = log2CbSize - ( quadtreeTUMaxDepth - 1 + interSplitFlag + intraSplitFlag); // stop when trafoDepth == hierarchy_depth = splitFlag
    if log2MinTUSizeInCU > this.m_pcSlice.GetSPS().GetQuadtreeTULog2MaxSize() {
      // when fully making use of signaled TUMaxDepth + inter/intraSplitFlag, resulting luma TB size is still > QuadtreeTULog2MaxSize
      log2MinTUSizeInCU = this.m_pcSlice.GetSPS().GetQuadtreeTULog2MaxSize();
    }  
  }
  return log2MinTUSizeInCU;
}
  
func (this *TComDataCU)  GetCUMvField         (  e RefPicList)      *TComCUMvField    { 
	return  &this.m_acCUMvField[e];  
}
  
func (this *TComDataCU)  GetCoeffY             ()                []TCoeff        { 
	return this.m_pcTrCoeffY;        
}
func (this *TComDataCU)  GetCoeffCb            ()                []TCoeff        { 
	return this.m_pcTrCoeffCb;       
}
func (this *TComDataCU)  GetCoeffCr            ()                []TCoeff        { 
	return this.m_pcTrCoeffCr;       
}
//#if ADAPTIVE_QP_SELECTION
func (this *TComDataCU)  GetArlCoeffY          ()                []int        { 
	return this.m_pcArlCoeffY;       
}
func (this *TComDataCU)  GetArlCoeffCb         ()                []int        { 
	return this.m_pcArlCoeffCb;      
}
func (this *TComDataCU)  GetArlCoeffCr         ()                []int        { 
	return this.m_pcArlCoeffCr;     
}
//#endif
  
func (this *TComDataCU)  GetPCMSampleY         ()                []Pel        { 
	return this.m_pcIPCMSampleY;    
}
func (this *TComDataCU)  GetPCMSampleCb        ()                []Pel        { 
	return this.m_pcIPCMSampleCb;    
}
func (this *TComDataCU)  GetPCMSampleCr        ()                []Pel        { 
	return this.m_pcIPCMSampleCr;    
}

func (this *TComDataCU)  GetCbf2    (  uiIdx uint,  eType TextType)     byte             { 
	return this.m_puhCbf[G_aucConvertTxtTypeToIdx[eType]][uiIdx];  
}
func (this *TComDataCU)  GetCbf1    (  eType TextType)                []byte              { 
	return this.m_puhCbf[G_aucConvertTxtTypeToIdx[eType]];         
}
func (this *TComDataCU)  GetCbf3    (  uiIdx uint,  eType TextType,  uiTrDepth uint) byte  { 
	return ( ( this.GetCbf2( uiIdx, eType ) >> uiTrDepth ) & 0x1 ); 
}
func (this *TComDataCU)  SetCbf    (  uiIdx uint,  eType TextType,  uh byte)        { 
	this.m_puhCbf[G_aucConvertTxtTypeToIdx[eType]][uiIdx] = uh;    
}
func (this *TComDataCU)  ClearCbf  (  uiIdx uint,  eType TextType,  uiNumParts uint){
}
func (this *TComDataCU)  GetQtRootCbf          (  uiIdx uint )            bool          { 
	return this.GetCbf3( uiIdx, TEXT_LUMA, 0 )!=0 || this.GetCbf3( uiIdx, TEXT_CHROMA_U, 0 )!=0 || this.GetCbf3( uiIdx, TEXT_CHROMA_V, 0 )!=0; 
}
  
func (this *TComDataCU)  SetCbfSubParts        (  uiCbfY,  uiCbfU,  uiCbfV,  uiAbsPartIdx,  uiDepth uint         ){
}
func (this *TComDataCU)  SetCbfSubParts4        (  uiCbf uint,  eTType TextType,  uiAbsPartIdx,  uiDepth  uint     ){
}
func (this *TComDataCU)  SetCbfSubParts5        (  uiCbf uint,  eTType TextType,  uiAbsPartIdx,  uiPartIdx,  uiDepth uint   ){
}
  
  // -------------------------------------------------------------------------------------------------------------------
  // member functions for coding tool information
  // -------------------------------------------------------------------------------------------------------------------
  
func (this *TComDataCU)  GetMergeFlag          ()                        []bool{ 
	return this.m_pbMergeFlag;               
}
func (this *TComDataCU)  GetMergeFlag1          (  uiIdx uint)              bool{ 
	return this.m_pbMergeFlag[uiIdx];        
}
func (this *TComDataCU)  SetMergeFlag          (  uiIdx uint,  b bool)    { 
	this.m_pbMergeFlag[uiIdx] = b;           
}
func (this *TComDataCU)  SetMergeFlagSubParts  (  bMergeFlag bool,  uiAbsPartIdx,  uiPartIdx,  uiDepth uint){
}

func (this *TComDataCU)  GetMergeIndex         ()                        []byte{ 
	return this.m_puhMergeIndex;                         
}
func (this *TComDataCU)  GetMergeIndex1         (  uiIdx uint)              byte{ 
	return this.m_puhMergeIndex[uiIdx];                  
}
func (this *TComDataCU)  SetMergeIndex         (  uiIdx uint,  uiMergeIndex byte) { 
	this.m_puhMergeIndex[uiIdx] = uiMergeIndex; 
}
func (this *TComDataCU)  SetMergeIndexSubParts (  uiMergeIndex,  uiAbsPartIdx,  uiPartIdx,  uiDepth uint){
}
//  template <typename T>
//func (this *TComDataCU)  SetSubPart            ( T bParameter, T* pbBaseLCU, UInt uiCUAddr, UInt uiCUDepth, UInt uiPUIdx );

//#if AMP_MRG
func (this *TComDataCU)  SetMergeAMP(  b bool)      { 
	this.m_bIsMergeAMP = b; 
}
func (this *TComDataCU)  GetMergeAMP( )         bool    { 
	return this.m_bIsMergeAMP; 
}
//#endif

func (this *TComDataCU)  GetLumaIntraDir       ()                   []byte    { 
	return this.m_puhLumaIntraDir;           
}
func (this *TComDataCU)  GetLumaIntraDir1       (  uiIdx uint)         byte   { 
	return this.m_puhLumaIntraDir[uiIdx];    
}
func (this *TComDataCU)  SetLumaIntraDir       (  uiIdx uint,  uh byte) { 
	this.m_puhLumaIntraDir[uiIdx] = uh;      
}
func (this *TComDataCU)  SetLumaIntraDirSubParts(  uiDir,  uiAbsPartIdx,  uiDepth uint){
}
  
func (this *TComDataCU)  GetChromaIntraDir     ()                  []byte      { 
	return this.m_puhChromaIntraDir;         
}
func (this *TComDataCU)  GetChromaIntraDir1     (  uiIdx uint)        byte    { 
	return this.m_puhChromaIntraDir[uiIdx];  
}
func (this *TComDataCU)  SetChromaIntraDir     (  uiIdx uint,   uh byte) { 
	this.m_puhChromaIntraDir[uiIdx] = uh;    
}
func (this *TComDataCU)  SetChromIntraDirSubParts(  uiDir,   uiAbsPartIdx,  uiDepth uint){
}
  
func (this *TComDataCU)  GetInterDir           ()                   []byte     { 
	return this.m_puhInterDir;               
}
func (this *TComDataCU)  GetInterDir1          (  uiIdx uint)         byte   { 
	return this.m_puhInterDir[uiIdx];        
}
func (this *TComDataCU)  SetInterDir           (  uiIdx uint,   uh byte) { 
	this.m_puhInterDir[uiIdx] = uh;          
}
func (this *TComDataCU)  SetInterDirSubParts   (  uiDir,   uiAbsPartIdx,  uiPartIdx,  uiDepth uint){
}
func (this *TComDataCU)  GetIPCMFlag           ()                   []bool     { 
	return this.m_pbIPCMFlag;               
}
func (this *TComDataCU)  GetIPCMFlag1          ( uiIdx uint )          bool   { 
	return this.m_pbIPCMFlag[uiIdx];        
}
func (this *TComDataCU)  SetIPCMFlag           ( uiIdx uint,  b bool)     { 
	this.m_pbIPCMFlag[uiIdx] = b;           
}
func (this *TComDataCU)  SetIPCMFlagSubParts   ( bIpcmFlag bool,  uiAbsPartIdx,  uiDepth uint){
}

/*#if !REMOVE_BURST_IPCM
  Int           GetNumSucIPCM         ()                        { return this.m_numSucIPCM;             }
  Void          SetNumSucIPCM         ( Int num )               { this.m_numSucIPCM = num;              }
  Bool          GetLastCUSucIPCMFlag  ()                        { return this.m_lastCUSucIPCMFlag;        }
  Void          SetLastCUSucIPCMFlag  ( Bool flg )              { this.m_lastCUSucIPCMFlag = flg;         }
#endif*/

  /// Get slice ID for SU
func (this *TComDataCU)  GetSUSliceID          ( uiIdx uint)            int  {
	return this.m_piSliceSUMap[uiIdx];      
} 

  /// Get the pointer of slice ID map
func (this *TComDataCU)  GetSliceSUMap         ()                      []int  {
	return this.m_piSliceSUMap;             
}

  /// Set the pointer of slice ID map
func (this *TComDataCU)  SetSliceSUMap         (pi []int)                 {
	this.m_piSliceSUMap = pi;               
}

func (this *TComDataCU)  GetNDBFilterBlocks()     *list.List {
	return this.m_vNDFBlock;
}
func (this *TComDataCU)  SetNDBFilterBlockBorderAvailability( numLCUInPicWidth,  numLCUInPicHeight,  numSUInLCUWidth,  numSUInLCUHeight,  picWidth,  picHeight uint,
                                          LFCrossSliceBoundary *list.List, bTopTileBoundary, bDownTileBoundary, bLeftTileBoundary, bRightTileBoundary, bIndependentTileBoundaryEnabled bool){
}
  // -------------------------------------------------------------------------------------------------------------------
  // member functions for accessing partition information
  // -------------------------------------------------------------------------------------------------------------------
func (this *TComDataCU)  GetPartIndexAndSizePos(  uiPartIdx uint, ruiPartAddr *uint, riWidth, riHeight, rPosX, rPosY *int ){
}
func (this *TComDataCU)  GetPartIndexAndSize   (  uiPartIdx uint, ruiPartAddr *uint, riWidth, riHeight *int ){
}
func (this *TComDataCU)  GetNumPartInter       () byte{
  iNumPart := byte(0);
  
  switch this.m_pePartSize[0] {
    case SIZE_2Nx2N:    iNumPart = 1;
    case SIZE_2NxN:     iNumPart = 2;
    case SIZE_Nx2N:     iNumPart = 2;
    case SIZE_NxN:      iNumPart = 4;
    case SIZE_2NxnU:    iNumPart = 2;
    case SIZE_2NxnD:    iNumPart = 2;
    case SIZE_nLx2N:    iNumPart = 2;
    case SIZE_nRx2N:    iNumPart = 2;
    //default:            assert (0);  ;
  }
  
  return  iNumPart;
}
func (this *TComDataCU)  IsFirstAbsZorderIdxInDepth ( uiAbsPartIdx,  uiDepth uint) bool{
  uiPartNumb := this.m_pcPic.GetNumPartInCU() >> (uiDepth << 1);
  return (((this.m_uiAbsIdxInLCU + uiAbsPartIdx)% uiPartNumb) == 0);
}
  
  // -------------------------------------------------------------------------------------------------------------------
  // member functions for motion vector
  // -------------------------------------------------------------------------------------------------------------------
  
func (this *TComDataCU)  GetMvField            ( pcCU *TComDataCU,  uiAbsPartIdx uint,  eRefPicList RefPicList, rcMvField *TComMvField){
}
  
func (this *TComDataCU)  FillMvpCand           (  uiPartIdx,  uiPartAddr uint,  eRefPicList RefPicList,  iRefIdx int, pInfo *AMVPInfo){
}
func (this *TComDataCU)  IsDiffMER             (  xN,  yN,  xP,  yP int) bool{
  plevel := this.GetSlice().GetPPS().GetLog2ParallelMergeLevelMinus2() + 2;
  if (xN>>plevel)!= (xP>>plevel) {
    return true;
  }
  
  if (yN>>plevel)!= (yP>>plevel) {
    return true;
  }
  
  return false;
}
func (this *TComDataCU)  GetPartPosition       (  partIdx uint, xP, yP, nPSW, nPSH *int){
}
func (this *TComDataCU)  SetMVPIdx             (  eRefPicList RefPicList, uiIdx uint, iMVPIdx int8)  { 
	this.m_apiMVPIdx[eRefPicList][uiIdx] = iMVPIdx;  
}
func (this *TComDataCU)  GetMVPIdx2             ( eRefPicList RefPicList, uiIdx uint)     int8          { 
	return this.m_apiMVPIdx[eRefPicList][uiIdx];     
}
func (this *TComDataCU)  GetMVPIdx1             ( eRefPicList RefPicList)                []int8          { 
	return this.m_apiMVPIdx[eRefPicList];            
}

func (this *TComDataCU)  SetMVPNum             ( eRefPicList RefPicList, uiIdx uint,  iMVPNum int8) { 
	this.m_apiMVPNum[eRefPicList][uiIdx] = iMVPNum;  
}
func (this *TComDataCU)  GetMVPNum2             ( eRefPicList RefPicList, uiIdx uint )     int8         { 
	return this.m_apiMVPNum[eRefPicList][uiIdx];     
}
func (this *TComDataCU)  GetMVPNum1             ( eRefPicList RefPicList)                []int8          { 
	return this.m_apiMVPNum[eRefPicList];            
}
  
func (this *TComDataCU)  SetMVPIdxSubParts     ( iMVPIdx int,  eRefPicList RefPicList,  uiAbsPartIdx,  uiPartIdx,  uiDepth uint){
}
func (this *TComDataCU)  SetMVPNumSubParts     ( iMVPIdx int,  eRefPicList RefPicList,  uiAbsPartIdx,  uiPartIdx,  uiDepth uint ){
}
  
func (this *TComDataCU)  ClipMv                ( rcMv  *TComMv   ){
}
func (this *TComDataCU)  GetMvPredLeft         ( )   *TComMv{ 
	return this.m_cMvFieldA.GetMv(); 
}
func (this *TComDataCU)  GetMvPredAbove        ( )   *TComMv{ 
	return this.m_cMvFieldB.GetMv(); 
}
func (this *TComDataCU)  GetMvPredAboveRight   ( )   *TComMv{ 
	return this.m_cMvFieldC.GetMv(); 
}
  
func (this *TComDataCU)  CompressMV            (){
}

  // -------------------------------------------------------------------------------------------------------------------
  // utility functions for neighbouring information
  // -------------------------------------------------------------------------------------------------------------------

func (this *TComDataCU)   GetCULeft                   () *TComDataCU{ 
	return this.m_pcCULeft;       
}
func (this *TComDataCU)   GetCUAbove                  () *TComDataCU{ 
	return this.m_pcCUAbove;    
}
func (this *TComDataCU)   GetCUAboveLeft              () *TComDataCU{ 
	return this.m_pcCUAboveLeft;  
}
func (this *TComDataCU)   GetCUAboveRight             () *TComDataCU{ 
	return this.m_pcCUAboveRight; 
}
func (this *TComDataCU)   GetCUColocated              (  eRefPicList RefPicList) *TComDataCU{ 
	return this.m_apcCUColocated[eRefPicList]; 
}

func (this *TComDataCU)  GetPULeft          ( uiLPartUnitIdx *uint, 
                                              uiCurrPartUnitIdx uint, 
                                              bEnforceSliceRestriction bool, 
                                              bEnforceDependentSliceRestriction bool,
                                              bEnforceTileRestriction bool ) *TComDataCU{
  uiAbsPartIdx       := G_auiZscanToRaster[uiCurrPartUnitIdx];
  uiAbsZorderCUIdx   := G_auiZscanToRaster[this.m_uiAbsIdxInLCU];
  uiNumPartInCUWidth := this.m_pcPic.GetNumPartInWidth();
  
  if  !IsZeroCol( int(uiAbsPartIdx), int(uiNumPartInCUWidth) ) {
    *uiLPartUnitIdx = G_auiRasterToZscan[ uiAbsPartIdx - 1 ];
    if IsEqualCol( int(uiAbsPartIdx), int(uiAbsZorderCUIdx), int(uiNumPartInCUWidth) ) {
      return this.m_pcPic.GetCU( this.GetAddr() );
    }else{
      *uiLPartUnitIdx -= this.m_uiAbsIdxInLCU;
      return this;
    }
  }
  
  *uiLPartUnitIdx = G_auiRasterToZscan[ uiAbsPartIdx + uiNumPartInCUWidth - 1 ];


  if  (bEnforceSliceRestriction 		 && (this.m_pcCULeft==nil || this.m_pcCULeft.GetSlice()==nil || this.m_pcCULeft.GetSCUAddr()+(*uiLPartUnitIdx) < this.m_pcPic.GetCU( this.GetAddr() ).GetSliceStartCU(uiCurrPartUnitIdx))) ||
      (bEnforceDependentSliceRestriction && (this.m_pcCULeft==nil || this.m_pcCULeft.GetSlice()==nil || this.m_pcCULeft.GetSCUAddr()+(*uiLPartUnitIdx) < this.m_pcPic.GetCU( this.GetAddr() ).GetDependentSliceStartCU(uiCurrPartUnitIdx)))  ||
      (bEnforceTileRestriction 			 && (this.m_pcCULeft==nil || this.m_pcCULeft.GetSlice()==nil || (this.m_pcPic.GetPicSym().GetTileIdxMap( int(this.m_pcCULeft.GetAddr()) ) != this.m_pcPic.GetPicSym().GetTileIdxMap( int(this.GetAddr())))  )  ) {
    return nil;
  }
  return this.m_pcCULeft;                                              
}
/*#if !LINEBUF_CLEANUP
  TComDataCU*   GetPUAbove                  ( UInt&  uiAPartUnitIdx, 
                                              UInt uiCurrPartUnitIdx, 
                                              Bool bEnforceSliceRestriction=true, 
                                              Bool bEnforceDependentSliceRestriction=true, 
                                              Bool MotionDataCompresssion = false,
                                              Bool planarAtLCUBoundary = false,
                                              Bool bEnforceTileRestriction=true );
  TComDataCU*   GetPUAboveLeft              ( UInt&  uiALPartUnitIdx, UInt uiCurrPartUnitIdx, Bool bEnforceSliceRestriction=true, Bool bEnforceDependentSliceRestriction=true, Bool MotionDataCompresssion = false );
  TComDataCU*   GetPUAboveRight             ( UInt&  uiARPartUnitIdx, UInt uiCurrPartUnitIdx, Bool bEnforceSliceRestriction=true, Bool bEnforceDependentSliceRestriction=true, Bool MotionDataCompresssion = false );
#else*/
func (this *TComDataCU)  GetPUAbove                  (  uiAPartUnitIdx *uint, 
						                                uiCurrPartUnitIdx uint, 
						                                bEnforceSliceRestriction bool, 
						                                bEnforceDependentSliceRestriction bool, 
						                                planarAtLCUBoundary bool,
						                                bEnforceTileRestriction bool )*TComDataCU{
  uiAbsPartIdx       := int(G_auiZscanToRaster[uiCurrPartUnitIdx])
  uiAbsZorderCUIdx   := int(G_auiZscanToRaster[this.m_uiAbsIdxInLCU])
  uiNumPartInCUWidth := int(this.m_pcPic.GetNumPartInWidth());
  
  if !IsZeroRow( uiAbsPartIdx, uiNumPartInCUWidth )  {
    *uiAPartUnitIdx = G_auiRasterToZscan[ uiAbsPartIdx - uiNumPartInCUWidth ];
    if IsEqualRow( uiAbsPartIdx, uiAbsZorderCUIdx, uiNumPartInCUWidth ) {
      return this.m_pcPic.GetCU( this.GetAddr() );
    }else{
      *uiAPartUnitIdx -= this.m_uiAbsIdxInLCU;
      return this;
    }
  }

  if planarAtLCUBoundary {
    return nil;
  }
  
  *uiAPartUnitIdx = G_auiRasterToZscan[ uiAbsPartIdx + int(this.m_pcPic.GetNumPartInCU()) - uiNumPartInCUWidth ];
/*#if !LINEBUF_CLEANUP
  if(MotionDataCompresssion)
  {
    uiAPartUnitIdx = G_motionRefer[uiAPartUnitIdx];
  }
#endif*/

  if (bEnforceSliceRestriction 			&& (this.m_pcCUAbove==nil || this.m_pcCUAbove.GetSlice()==nil || this.m_pcCUAbove.GetSCUAddr()+(*uiAPartUnitIdx) < this.m_pcPic.GetCU( this.GetAddr() ).GetSliceStartCU(uiCurrPartUnitIdx))) ||
     (bEnforceDependentSliceRestriction && (this.m_pcCUAbove==nil || this.m_pcCUAbove.GetSlice()==nil || this.m_pcCUAbove.GetSCUAddr()+(*uiAPartUnitIdx) < this.m_pcPic.GetCU( this.GetAddr() ).GetDependentSliceStartCU(uiCurrPartUnitIdx))) ||
     (bEnforceTileRestriction 			&& (this.m_pcCUAbove==nil || this.m_pcCUAbove.GetSlice()==nil || (this.m_pcPic.GetPicSym().GetTileIdxMap( int(this.m_pcCUAbove.GetAddr()) ) != this.m_pcPic.GetPicSym().GetTileIdxMap(int(this.GetAddr()))))) {
    return nil;
  }
  return this.m_pcCUAbove;
}
func (this *TComDataCU)  GetPUAboveLeft              ( uiALPartUnitIdx *uint, uiCurrPartUnitIdx uint, bEnforceSliceRestriction, bEnforceDependentSliceRestriction bool ) *TComDataCU{
  uiAbsPartIdx       := int(G_auiZscanToRaster[uiCurrPartUnitIdx]);
  uiAbsZorderCUIdx   := int(G_auiZscanToRaster[this.m_uiAbsIdxInLCU]);
  uiNumPartInCUWidth := int(this.m_pcPic.GetNumPartInWidth());
  
  if !IsZeroCol( uiAbsPartIdx, uiNumPartInCUWidth ) {
    if !IsZeroRow( uiAbsPartIdx, uiNumPartInCUWidth ) {
      *uiALPartUnitIdx = G_auiRasterToZscan[ uiAbsPartIdx - uiNumPartInCUWidth - 1 ];
      if IsEqualRowOrCol( uiAbsPartIdx, uiAbsZorderCUIdx, uiNumPartInCUWidth ) {
        return this.m_pcPic.GetCU( this.GetAddr() );
      }else{
        *uiALPartUnitIdx -= this.m_uiAbsIdxInLCU;
        return this;
      }
    }
    *uiALPartUnitIdx = G_auiRasterToZscan[ uiAbsPartIdx + int(this.GetPic().GetNumPartInCU()) - uiNumPartInCUWidth - 1 ];
/*#if !LINEBUF_CLEANUP
    if(MotionDataCompresssion)
    {
      uiALPartUnitIdx = g_motionRefer[uiALPartUnitIdx];
    }
#endif*/
    if (bEnforceSliceRestriction          && (this.m_pcCUAbove==nil || this.m_pcCUAbove.GetSlice()==nil || this.m_pcCUAbove.GetSCUAddr()+(*uiALPartUnitIdx) < this.m_pcPic.GetCU( this.GetAddr() ).GetSliceStartCU(uiCurrPartUnitIdx)         ||(this.m_pcPic.GetPicSym().GetTileIdxMap( int(this.m_pcCUAbove.GetAddr()) ) != this.m_pcPic.GetPicSym().GetTileIdxMap(int(this.GetAddr()))) ))||
       (bEnforceDependentSliceRestriction && (this.m_pcCUAbove==nil || this.m_pcCUAbove.GetSlice()==nil || this.m_pcCUAbove.GetSCUAddr()+(*uiALPartUnitIdx) < this.m_pcPic.GetCU( this.GetAddr() ).GetDependentSliceStartCU(uiCurrPartUnitIdx)||(this.m_pcPic.GetPicSym().GetTileIdxMap( int(this.m_pcCUAbove.GetAddr()) ) != this.m_pcPic.GetPicSym().GetTileIdxMap(int(this.GetAddr()))) )) {
      return nil;
    }
    return this.m_pcCUAbove;
  }
  
  if !IsZeroRow( uiAbsPartIdx, uiNumPartInCUWidth ) {
    *uiALPartUnitIdx = G_auiRasterToZscan[ uiAbsPartIdx - 1 ];
    if (bEnforceSliceRestriction && (this.m_pcCULeft==nil || this.m_pcCULeft.GetSlice()==nil || 
        this.m_pcCULeft.GetSCUAddr()+(*uiALPartUnitIdx) < this.m_pcPic.GetCU( this.GetAddr() ).GetSliceStartCU(uiCurrPartUnitIdx)||
       (this.m_pcPic.GetPicSym().GetTileIdxMap( int(this.m_pcCULeft.GetAddr()) ) != this.m_pcPic.GetPicSym().GetTileIdxMap(int(this.GetAddr()))) ))||
       (bEnforceDependentSliceRestriction && (this.m_pcCULeft==nil || this.m_pcCULeft.GetSlice()==nil || 
        this.m_pcCULeft.GetSCUAddr()+(*uiALPartUnitIdx) < this.m_pcPic.GetCU( this.GetAddr() ).GetDependentSliceStartCU(uiCurrPartUnitIdx)||
       (this.m_pcPic.GetPicSym().GetTileIdxMap( int(this.m_pcCULeft.GetAddr()) ) != this.m_pcPic.GetPicSym().GetTileIdxMap(int(this.GetAddr()))) )) {
      return nil;
    }
    return this.m_pcCULeft;
  }
  
  *uiALPartUnitIdx = G_auiRasterToZscan[ this.m_pcPic.GetNumPartInCU() - 1 ];
/*#if !LINEBUF_CLEANUP
  if(MotionDataCompresssion)
  {
    uiALPartUnitIdx = g_motionRefer[uiALPartUnitIdx];
  }
#endif*/
  if (bEnforceSliceRestriction && (this.m_pcCUAboveLeft==nil || this.m_pcCUAboveLeft.GetSlice()==nil || 
      this.m_pcCUAboveLeft.GetSCUAddr()+(*uiALPartUnitIdx) < this.m_pcPic.GetCU( this.GetAddr() ).GetSliceStartCU(uiCurrPartUnitIdx)||
      (this.m_pcPic.GetPicSym().GetTileIdxMap( int(this.m_pcCUAboveLeft.GetAddr()) ) != this.m_pcPic.GetPicSym().GetTileIdxMap(int(this.GetAddr()))) ))||
     (bEnforceDependentSliceRestriction && (this.m_pcCUAboveLeft==nil || this.m_pcCUAboveLeft.GetSlice()==nil || 
      this.m_pcCUAboveLeft.GetSCUAddr()+(*uiALPartUnitIdx) < this.m_pcPic.GetCU( this.GetAddr() ).GetDependentSliceStartCU(uiCurrPartUnitIdx)||
      (this.m_pcPic.GetPicSym().GetTileIdxMap( int(this.m_pcCUAboveLeft.GetAddr()) ) != this.m_pcPic.GetPicSym().GetTileIdxMap(int(this.GetAddr()))) ))  {
    return nil;
  }
  return this.m_pcCUAboveLeft;
}
func (this *TComDataCU)  GetPUAboveRight             ( uiARPartUnitIdx *uint, uiCurrPartUnitIdx uint, bEnforceSliceRestriction, bEnforceDependentSliceRestriction bool ) *TComDataCU{
  uiAbsPartIdxRT     := int(G_auiZscanToRaster[uiCurrPartUnitIdx]);
  uiAbsZorderCUIdx   := int(G_auiZscanToRaster[this.m_uiAbsIdxInLCU ]) + int(this.m_puhWidth[0]) / int(this.m_pcPic.GetMinCUWidth()) - 1;
  uiNumPartInCUWidth := int(this.m_pcPic.GetNumPartInWidth());
  
  if ( this.m_pcPic.GetCU(this.m_uiCUAddr).GetCUPelX() + G_auiRasterToPelX[uiAbsPartIdxRT] + this.m_pcPic.GetMinCUWidth() ) >= this.m_pcSlice.GetSPS().GetPicWidthInLumaSamples() {
    *uiARPartUnitIdx = MAX_UINT;
    return nil;
  }
  
  if LessThanCol( uiAbsPartIdxRT, uiNumPartInCUWidth - 1, uiNumPartInCUWidth ) {
    if !IsZeroRow( uiAbsPartIdxRT, uiNumPartInCUWidth ) {
      if uiCurrPartUnitIdx > G_auiRasterToZscan[ uiAbsPartIdxRT - uiNumPartInCUWidth + 1 ] {
        *uiARPartUnitIdx = G_auiRasterToZscan[ uiAbsPartIdxRT - uiNumPartInCUWidth + 1 ];
        if IsEqualRowOrCol( uiAbsPartIdxRT, uiAbsZorderCUIdx, uiNumPartInCUWidth ) {
          return this.m_pcPic.GetCU( this.GetAddr() );
        }else{
          *uiARPartUnitIdx -= this.m_uiAbsIdxInLCU;
          return this;
        }
      }
      *uiARPartUnitIdx = MAX_UINT;
      return nil;
    }
    *uiARPartUnitIdx = G_auiRasterToZscan[ uiAbsPartIdxRT + int(this.m_pcPic.GetNumPartInCU()) - uiNumPartInCUWidth + 1 ];
/*#if !LINEBUF_CLEANUP
    if(MotionDataCompresssion)
    {
      uiARPartUnitIdx = g_motionRefer[uiARPartUnitIdx];
    }
#endif*/
    if (bEnforceSliceRestriction && (this.m_pcCUAbove==nil || this.m_pcCUAbove.GetSlice()==nil || 
       this.m_pcCUAbove.GetSCUAddr()+(*uiARPartUnitIdx) < this.m_pcPic.GetCU( this.GetAddr() ).GetSliceStartCU(uiCurrPartUnitIdx)||
       (this.m_pcPic.GetPicSym().GetTileIdxMap( int(this.m_pcCUAbove.GetAddr()) ) != this.m_pcPic.GetPicSym().GetTileIdxMap(int(this.GetAddr()))) ))||
       (bEnforceDependentSliceRestriction && (this.m_pcCUAbove==nil || this.m_pcCUAbove.GetSlice()==nil || 
       this.m_pcCUAbove.GetSCUAddr()+(*uiARPartUnitIdx) < this.m_pcPic.GetCU( this.GetAddr() ).GetDependentSliceStartCU(uiCurrPartUnitIdx)||
       (this.m_pcPic.GetPicSym().GetTileIdxMap( int(this.m_pcCUAbove.GetAddr()) ) != this.m_pcPic.GetPicSym().GetTileIdxMap(int(this.GetAddr()))) )) {
      return nil;
    }
    return this.m_pcCUAbove;
  }
  
  if !IsZeroRow( uiAbsPartIdxRT, uiNumPartInCUWidth ) {
    *uiARPartUnitIdx = MAX_UINT;
    return nil;
  }
  
  *uiARPartUnitIdx = G_auiRasterToZscan[ int(this.m_pcPic.GetNumPartInCU()) - uiNumPartInCUWidth ];
/*#if !LINEBUF_CLEANUP
  if(MotionDataCompresssion)
  {
    uiARPartUnitIdx = g_motionRefer[uiARPartUnitIdx];
  }
#endif*/
  if  (bEnforceSliceRestriction && (this.m_pcCUAboveRight==nil || this.m_pcCUAboveRight.GetSlice()==nil ||
       this.m_pcPic.GetPicSym().GetInverseCUOrderMap( int(this.m_pcCUAboveRight.GetAddr())) > this.m_pcPic.GetPicSym().GetInverseCUOrderMap( int(this.GetAddr())) ||
       this.m_pcCUAboveRight.GetSCUAddr()+(*uiARPartUnitIdx) < this.m_pcPic.GetCU( this.GetAddr() ).GetSliceStartCU(uiCurrPartUnitIdx)||
       (this.m_pcPic.GetPicSym().GetTileIdxMap( int(this.m_pcCUAboveRight.GetAddr()) ) != this.m_pcPic.GetPicSym().GetTileIdxMap(int(this.GetAddr()))) ))||
       (bEnforceDependentSliceRestriction && (this.m_pcCUAboveRight==nil || this.m_pcCUAboveRight.GetSlice()==nil || 
       this.m_pcPic.GetPicSym().GetInverseCUOrderMap( int(this.m_pcCUAboveRight.GetAddr())) > this.m_pcPic.GetPicSym().GetInverseCUOrderMap( int(this.GetAddr())) ||
       this.m_pcCUAboveRight.GetSCUAddr()+(*uiARPartUnitIdx) < this.m_pcPic.GetCU( this.GetAddr() ).GetDependentSliceStartCU(uiCurrPartUnitIdx)||
       (this.m_pcPic.GetPicSym().GetTileIdxMap( int(this.m_pcCUAboveRight.GetAddr()) ) != this.m_pcPic.GetPicSym().GetTileIdxMap(int(this.GetAddr()))) )){
    return nil;
  }
  return this.m_pcCUAboveRight;
}
//#endif
func (this *TComDataCU)  GetPUBelowLeft              ( uiBLPartUnitIdx *uint, uiCurrPartUnitIdx uint, bEnforceSliceRestriction, bEnforceDependentSliceRestriction bool) *TComDataCU{
  uiAbsPartIdxLB     := int(G_auiZscanToRaster[uiCurrPartUnitIdx]);
  uiAbsZorderCUIdxLB := int(G_auiZscanToRaster[this.m_uiAbsIdxInLCU ]) + (int(this.m_puhHeight[0]) / int(this.m_pcPic.GetMinCUHeight()) - 1)*int(this.m_pcPic.GetNumPartInWidth());
  uiNumPartInCUWidth := int(this.m_pcPic.GetNumPartInWidth());
  
  if ( this.m_pcPic.GetCU(this.m_uiCUAddr).GetCUPelY() + G_auiRasterToPelY[uiAbsPartIdxLB] + this.m_pcPic.GetMinCUHeight() ) >= this.m_pcSlice.GetSPS().GetPicHeightInLumaSamples() {
    *uiBLPartUnitIdx = MAX_UINT;
    return nil;
  }
  
  if LessThanRow( uiAbsPartIdxLB, int(this.m_pcPic.GetNumPartInHeight()) - 1, uiNumPartInCUWidth ) {
    if !IsZeroCol( uiAbsPartIdxLB, uiNumPartInCUWidth ) {
      if uiCurrPartUnitIdx > G_auiRasterToZscan[ uiAbsPartIdxLB + uiNumPartInCUWidth - 1 ] {
        *uiBLPartUnitIdx = G_auiRasterToZscan[ uiAbsPartIdxLB + uiNumPartInCUWidth - 1 ];
        if IsEqualRowOrCol( uiAbsPartIdxLB, uiAbsZorderCUIdxLB, uiNumPartInCUWidth ) {
          return this.m_pcPic.GetCU( this.GetAddr() );
        }else{
          *uiBLPartUnitIdx -= this.m_uiAbsIdxInLCU;
          return this;
        }
      }
      *uiBLPartUnitIdx = MAX_UINT;
      return nil;
    }
    *uiBLPartUnitIdx = G_auiRasterToZscan[ uiAbsPartIdxLB + uiNumPartInCUWidth*2 - 1 ];
    if (bEnforceSliceRestriction && (this.m_pcCULeft==nil || this.m_pcCULeft.GetSlice()==nil || 
       this.m_pcCULeft.GetSCUAddr()+(*uiBLPartUnitIdx) < this.m_pcPic.GetCU( this.GetAddr() ).GetSliceStartCU(uiCurrPartUnitIdx)||
       (this.m_pcPic.GetPicSym().GetTileIdxMap( int(this.m_pcCULeft.GetAddr()) ) != this.m_pcPic.GetPicSym().GetTileIdxMap(int(this.GetAddr()))) ))||
       (bEnforceDependentSliceRestriction && (this.m_pcCULeft==nil || this.m_pcCULeft.GetSlice()==nil || 
       this.m_pcCULeft.GetSCUAddr()+(*uiBLPartUnitIdx) < this.m_pcPic.GetCU( this.GetAddr() ).GetDependentSliceStartCU(uiCurrPartUnitIdx)||
       (this.m_pcPic.GetPicSym().GetTileIdxMap( int(this.m_pcCULeft.GetAddr()) ) != this.m_pcPic.GetPicSym().GetTileIdxMap(int(this.GetAddr()))) )) {
      return nil;
    }
    return this.m_pcCULeft;
  }
  
  *uiBLPartUnitIdx = MAX_UINT;
  return nil;
}

func (this *TComDataCU)  GetQpMinCuLeft              ( uiLPartUnitIdx *uint, uiCurrAbsIdxInLCU uint) *TComDataCU{
  numPartInCUWidth    := int(this.m_pcPic.GetNumPartInWidth());
  absZorderQpMinCUIdx := (uiCurrAbsIdxInLCU>>((G_uiMaxCUDepth - this.GetSlice().GetPPS().GetMaxCuDQPDepth())<<1))<<((G_uiMaxCUDepth - this.GetSlice().GetPPS().GetMaxCuDQPDepth())<<1);
  absRorderQpMinCUIdx := int(G_auiZscanToRaster[absZorderQpMinCUIdx]);

  // check for left LCU boundary
  if IsZeroCol(absRorderQpMinCUIdx, numPartInCUWidth) {
    return nil;
  }

  // get index of left-CU relative to top-left corner of current quantization group
  *uiLPartUnitIdx = G_auiRasterToZscan[absRorderQpMinCUIdx - 1];

  // return pointer to current LCU
  return this.m_pcPic.GetCU( this.GetAddr() );
}
func (this *TComDataCU)  GetQpMinCuAbove             ( aPartUnitIdx *uint, currAbsIdxInLCU uint) *TComDataCU{
  numPartInCUWidth    := int(this.m_pcPic.GetNumPartInWidth());
  absZorderQpMinCUIdx := (currAbsIdxInLCU>>((G_uiMaxCUDepth - this.GetSlice().GetPPS().GetMaxCuDQPDepth())<<1))<<((G_uiMaxCUDepth - this.GetSlice().GetPPS().GetMaxCuDQPDepth())<<1);
  absRorderQpMinCUIdx := int(G_auiZscanToRaster[absZorderQpMinCUIdx]);

  // check for top LCU boundary
  if IsZeroRow( absRorderQpMinCUIdx, numPartInCUWidth) {
    return nil;
  }

  // get index of top-CU relative to top-left corner of current quantization group
  *aPartUnitIdx = G_auiRasterToZscan[absRorderQpMinCUIdx - numPartInCUWidth];

  // return pointer to current LCU
  return this.m_pcPic.GetCU( this.GetAddr() );
}
func (this *TComDataCU)  GetRefQP                    ( uiCurrAbsIdxInLCU uint) int8{
  lPartIdx := uint(0);
  aPartIdx := uint(0);
  cULeft  := this.GetQpMinCuLeft ( &lPartIdx, this.m_uiAbsIdxInLCU + uiCurrAbsIdxInLCU );
  cUAbove := this.GetQpMinCuAbove( &aPartIdx, this.m_uiAbsIdxInLCU + uiCurrAbsIdxInLCU );
  
  if cULeft!=nil && cUAbove!=nil{
  	return (cULeft.GetQP1( lPartIdx ) + cUAbove.GetQP1( aPartIdx ) + 1) >> 1;
  }else if cUAbove!=nil {
  	return (this.GetLastCodedQP( uiCurrAbsIdxInLCU ) + cUAbove.GetQP1( aPartIdx )  + 1) >> 1;
  }else if cULeft!=nil {
	return (cULeft.GetQP1( lPartIdx ) + this.GetLastCodedQP( uiCurrAbsIdxInLCU ) + 1) >> 1;
  }

  return (this.GetLastCodedQP( uiCurrAbsIdxInLCU ) +  this.GetLastCodedQP( uiCurrAbsIdxInLCU ) + 1) >> 1;
}

func (this *TComDataCU)  GetPUAboveRightAdi          ( uiARPartUnitIdx *uint,  uiCurrPartUnitIdx,  uiPartUnitOffset uint,  bEnforceSliceRestriction,  bEnforceDependentSliceRestriction bool ) *TComDataCU{
  uiAbsPartIdxRT     := int(G_auiZscanToRaster[uiCurrPartUnitIdx]);
  uiAbsZorderCUIdx   := int(G_auiZscanToRaster[ this.m_uiAbsIdxInLCU ]) + (int(this.m_puhWidth[0]) / int(this.m_pcPic.GetMinCUWidth())) - 1;
  uiNumPartInCUWidth := int(this.m_pcPic.GetNumPartInWidth());
  
  if ( this.m_pcPic.GetCU(this.m_uiCUAddr).GetCUPelX() + G_auiRasterToPelX[uiAbsPartIdxRT] + (this.m_pcPic.GetPicSym().GetMinCUHeight() * uiPartUnitOffset)) >= this.m_pcSlice.GetSPS().GetPicWidthInLumaSamples() {
    *uiARPartUnitIdx = MAX_UINT;
    return nil;
  }
  
  if LessThanCol( uiAbsPartIdxRT, uiNumPartInCUWidth - int(uiPartUnitOffset), uiNumPartInCUWidth ) {
    if !IsZeroRow( uiAbsPartIdxRT, uiNumPartInCUWidth ) {
      if uiCurrPartUnitIdx > G_auiRasterToZscan[ uiAbsPartIdxRT - uiNumPartInCUWidth + int(uiPartUnitOffset) ] {
        *uiARPartUnitIdx = G_auiRasterToZscan[ uiAbsPartIdxRT - uiNumPartInCUWidth + int(uiPartUnitOffset) ];
        if IsEqualRowOrCol( uiAbsPartIdxRT, uiAbsZorderCUIdx, uiNumPartInCUWidth ) {
          return this.m_pcPic.GetCU( this.GetAddr() );
        }else{
          *uiARPartUnitIdx -= this.m_uiAbsIdxInLCU;
          return this;
        }
      }
      *uiARPartUnitIdx = MAX_UINT;
      return nil;
    }
    *uiARPartUnitIdx = G_auiRasterToZscan[ uiAbsPartIdxRT + int(this.m_pcPic.GetNumPartInCU()) - uiNumPartInCUWidth + int(uiPartUnitOffset) ];
    if (bEnforceSliceRestriction && (this.m_pcCUAbove==nil || this.m_pcCUAbove.GetSlice()==nil || 
       this.m_pcCUAbove.GetSCUAddr()+(*uiARPartUnitIdx) < this.m_pcPic.GetCU( this.GetAddr() ).GetSliceStartCU(uiCurrPartUnitIdx)||
       (this.m_pcPic.GetPicSym().GetTileIdxMap( int(this.m_pcCUAbove.GetAddr()) ) != this.m_pcPic.GetPicSym().GetTileIdxMap(int(this.GetAddr()))) ))||
       (bEnforceDependentSliceRestriction && (this.m_pcCUAbove==nil || this.m_pcCUAbove.GetSlice()==nil || 
       this.m_pcCUAbove.GetSCUAddr()+(*uiARPartUnitIdx) < this.m_pcPic.GetCU( this.GetAddr() ).GetDependentSliceStartCU(uiCurrPartUnitIdx)||
       (this.m_pcPic.GetPicSym().GetTileIdxMap( int(this.m_pcCUAbove.GetAddr()) ) != this.m_pcPic.GetPicSym().GetTileIdxMap(int(this.GetAddr()))) )) {
      return nil;
    }
    return this.m_pcCUAbove;
  }
  
  if !IsZeroRow( uiAbsPartIdxRT, uiNumPartInCUWidth ) {
    *uiARPartUnitIdx = MAX_UINT;
    return nil;
  }
  
  *uiARPartUnitIdx = G_auiRasterToZscan[ int(this.m_pcPic.GetNumPartInCU()) - uiNumPartInCUWidth + int(uiPartUnitOffset)-1 ];
  if (bEnforceSliceRestriction && (this.m_pcCUAboveRight==nil || this.m_pcCUAboveRight.GetSlice()==nil ||
       this.m_pcPic.GetPicSym().GetInverseCUOrderMap( int(this.m_pcCUAboveRight.GetAddr())) > this.m_pcPic.GetPicSym().GetInverseCUOrderMap( int(this.GetAddr())) ||
       this.m_pcCUAboveRight.GetSCUAddr()+(*uiARPartUnitIdx) < this.m_pcPic.GetCU( this.GetAddr() ).GetSliceStartCU(uiCurrPartUnitIdx)||
       (this.m_pcPic.GetPicSym().GetTileIdxMap( int(this.m_pcCUAboveRight.GetAddr()) ) != this.m_pcPic.GetPicSym().GetTileIdxMap(int(this.GetAddr()))) ))||
       (bEnforceDependentSliceRestriction && (this.m_pcCUAboveRight==nil || this.m_pcCUAboveRight.GetSlice()==nil || 
       this.m_pcPic.GetPicSym().GetInverseCUOrderMap( int(this.m_pcCUAboveRight.GetAddr())) > this.m_pcPic.GetPicSym().GetInverseCUOrderMap( int(this.GetAddr())) ||
       this.m_pcCUAboveRight.GetSCUAddr()+(*uiARPartUnitIdx) < this.m_pcPic.GetCU( this.GetAddr() ).GetDependentSliceStartCU(uiCurrPartUnitIdx)||
       (this.m_pcPic.GetPicSym().GetTileIdxMap( int(this.m_pcCUAboveRight.GetAddr()) ) != this.m_pcPic.GetPicSym().GetTileIdxMap(int(this.GetAddr()))) )) {
    return nil;
  }
  return this.m_pcCUAboveRight;
}
func (this *TComDataCU)  GetPUBelowLeftAdi           ( uiBLPartUnitIdx *uint,  uiCurrPartUnitIdx,  uiPartUnitOffset uint,  bEnforceSliceRestriction,  bEnforceDependentSliceRestriction bool ) *TComDataCU{
  uiAbsPartIdxLB     := int(G_auiZscanToRaster[uiCurrPartUnitIdx]);
  uiAbsZorderCUIdxLB := int(G_auiZscanToRaster[ this.m_uiAbsIdxInLCU ]) + ((int(this.m_puhHeight[0]) / int(this.m_pcPic.GetMinCUHeight())) - 1)*int(this.m_pcPic.GetNumPartInWidth());
  uiNumPartInCUWidth := int(this.m_pcPic.GetNumPartInWidth());
  
  if ( this.m_pcPic.GetCU(this.m_uiCUAddr).GetCUPelY() + G_auiRasterToPelY[uiAbsPartIdxLB] + (this.m_pcPic.GetPicSym().GetMinCUHeight() * uiPartUnitOffset)) >= this.m_pcSlice.GetSPS().GetPicHeightInLumaSamples() {
    *uiBLPartUnitIdx = MAX_UINT;
    return nil;
  }
  
  if LessThanRow( uiAbsPartIdxLB, int(this.m_pcPic.GetNumPartInHeight() - uiPartUnitOffset), uiNumPartInCUWidth ) {
    if !IsZeroCol( uiAbsPartIdxLB, uiNumPartInCUWidth ) {
      if uiCurrPartUnitIdx > G_auiRasterToZscan[ uiAbsPartIdxLB + int(uiPartUnitOffset) * uiNumPartInCUWidth - 1 ] {
        *uiBLPartUnitIdx = G_auiRasterToZscan[ uiAbsPartIdxLB + int(uiPartUnitOffset) * uiNumPartInCUWidth - 1 ];
        if IsEqualRowOrCol( uiAbsPartIdxLB, uiAbsZorderCUIdxLB, uiNumPartInCUWidth ) {
          return this.m_pcPic.GetCU( this.GetAddr() );
        }else{
          *uiBLPartUnitIdx -= this.m_uiAbsIdxInLCU;
          return this;
        }
      }
      *uiBLPartUnitIdx = MAX_UINT;
      return nil;
    }
    *uiBLPartUnitIdx = G_auiRasterToZscan[ uiAbsPartIdxLB + (1+int(uiPartUnitOffset)) * uiNumPartInCUWidth - 1 ];
    if (bEnforceSliceRestriction && (this.m_pcCULeft==nil || this.m_pcCULeft.GetSlice()==nil || 
       this.m_pcCULeft.GetSCUAddr()+(*uiBLPartUnitIdx) < this.m_pcPic.GetCU( this.GetAddr() ).GetSliceStartCU(uiCurrPartUnitIdx)||
       (this.m_pcPic.GetPicSym().GetTileIdxMap( int(this.m_pcCULeft.GetAddr()) ) != this.m_pcPic.GetPicSym().GetTileIdxMap(int(this.GetAddr()))) ))||
       (bEnforceDependentSliceRestriction && (this.m_pcCULeft==nil || this.m_pcCULeft.GetSlice()==nil || 
       this.m_pcCULeft.GetSCUAddr()+(*uiBLPartUnitIdx) < this.m_pcPic.GetCU( this.GetAddr() ).GetDependentSliceStartCU(uiCurrPartUnitIdx)||
       (this.m_pcPic.GetPicSym().GetTileIdxMap( int(this.m_pcCULeft.GetAddr()) ) != this.m_pcPic.GetPicSym().GetTileIdxMap(int(this.GetAddr()))) )) {
      return nil;
    }
    return this.m_pcCULeft;
  }
  
  *uiBLPartUnitIdx = MAX_UINT;
  return nil;
}
  
func (this *TComDataCU)  DeriveLeftRightTopIdx       ( uiPartIdx uint, ruiPartIdxLT, ruiPartIdxRT *uint){
}
func (this *TComDataCU)  DeriveLeftBottomIdx         ( uiPartIdx uint, ruiPartIdxLB *uint){
}
  
func (this *TComDataCU)  DeriveLeftRightTopIdxAdi    ( ruiPartIdxLT, ruiPartIdxRT *uint,  uiPartOffSet,  uiPartDepth uint){
}
func (this *TComDataCU)  DeriveLeftBottomIdxAdi      ( ruiPartIdxLB *uint,   uiPartOffSet,  uiPartDepth uint){
}
  
func (this *TComDataCU)  HasEqualMotion              (  uiAbsPartIdx uint, pcCandCU *TComDataCU,  uiCandAbsPartIdx uint) bool{
  if this.GetInterDir1( uiAbsPartIdx ) != pcCandCU.GetInterDir1( uiCandAbsPartIdx ) {
    return false;
  }

  for uiRefListIdx := uint(0); uiRefListIdx < 2; uiRefListIdx++ {
    if (this.GetInterDir1( uiAbsPartIdx ) & ( 1 << uiRefListIdx ))!=0 {
      if this.GetCUMvField( RefPicList( uiRefListIdx ) ).GetMv    ( int(uiAbsPartIdx) ) != pcCandCU.GetCUMvField( RefPicList( uiRefListIdx ) ).GetMv    ( int(uiCandAbsPartIdx) ) || 
         this.GetCUMvField( RefPicList( uiRefListIdx ) ).GetRefIdx( int(uiAbsPartIdx) ) != pcCandCU.GetCUMvField( RefPicList( uiRefListIdx ) ).GetRefIdx( int(uiCandAbsPartIdx) )  {
        return false;
      }
    }
  }

  return true;
}
func (this *TComDataCU)  GetInterMergeCandidates       ( uiAbsPartIdx,  uiPUIdx uint, pcMFieldNeighbours *TComMvField, puhInterDirNeighbours *byte, numValidMergeCand *int, mrgCandIdx int ){
}
func (this *TComDataCU)  DeriveLeftRightTopIdxGeneral  (  uiAbsPartIdx,  uiPartIdx uint, ruiPartIdxLT, ruiPartIdxRT *uint ){
}
func (this *TComDataCU)  DeriveLeftBottomIdxGeneral    (  uiAbsPartIdx,  uiPartIdx uint, UruiPartIdxLB *uint){
}
  
  
  // -------------------------------------------------------------------------------------------------------------------
  // member functions for modes
  // -------------------------------------------------------------------------------------------------------------------
  
func (this *TComDataCU)  IsIntra   (  uiPartIdx uint)  bool{ 
	return this.m_pePredMode[ uiPartIdx ] == MODE_INTRA; 
}
func (this *TComDataCU)  IsSkipped (  uiPartIdx uint) bool{
	return this.GetSkipFlag1( uiPartIdx );
}                                  ///< SKIP (no residual)
func (this *TComDataCU)  IsBipredRestriction(  puIdx uint) bool{
  width := int(0);
  height := int(0);
  var partAddr uint;

  this.GetPartIndexAndSize( puIdx, &partAddr, &width, &height );
  if this.GetWidth1(0) == 8 && (width < 8 || height < 8) {
    return true;
  }
  return false;
}

  // -------------------------------------------------------------------------------------------------------------------
  // member functions for symbol prediction (most probable / mode conversion)
  // -------------------------------------------------------------------------------------------------------------------
  
func (this *TComDataCU)  GetIntraSizeIdx                 (  uiAbsPartIdx    uint                                   )uint{
  var uiShift uint;
  
  //uiShift := ( (m_puhTrIdx[uiAbsPartIdx]==0) && (m_pePartSize[uiAbsPartIdx]==SIZE_NxN) ) ? m_puhTrIdx[uiAbsPartIdx]+1 : m_puhTrIdx[uiAbsPartIdx];
  if this.m_pePartSize[uiAbsPartIdx]==SIZE_NxN {
  	uiShift = 1;
  }else{
  	uiShift = 0;
  }
  
  uiWidth := this.m_puhWidth[uiAbsPartIdx]>>uiShift;
  uiCnt := uint(0);
  for uiWidth!=0 {
    uiCnt++;
    uiWidth>>=1;
  }
  uiCnt-=2;
  
  if uiCnt > 6 {
  	return 6;
  } 
  
  return uiCnt;
}
  
func (this *TComDataCU)  GetAllowedChromaDir             (  uiAbsPartIdx uint, uiModeList *uint ){
}
func (this *TComDataCU)  GetIntraDirLumaPredictor        (  uiAbsPartIdx uint, uiIntraDirPred *int, piMode *int ) int{
}
  
  // -------------------------------------------------------------------------------------------------------------------
  // member functions for SBAC context
  // -------------------------------------------------------------------------------------------------------------------
  
func (this *TComDataCU)  GetCtxSplitFlag                 (    uiAbsPartIdx,  uiDepth   uint                ) uint{
}
func (this *TComDataCU)  GetCtxQtCbf                     (  eType TextType, uiTrDepth uint ) uint{
}

func (this *TComDataCU)  GetCtxSkipFlag                  (    uiAbsPartIdx  uint                              )uint{
}
func (this *TComDataCU)  GetCtxInterDir                  (    uiAbsPartIdx  uint                               )uint{
}
  
func (this *TComDataCU)  GetSliceStartCU         (  pos uint)                  uint{ 
	return this.m_uiSliceStartCU[pos-this.m_uiAbsIdxInLCU];                                                                                          
}
func (this *TComDataCU)  GetDependentSliceStartCU  (  pos uint)                uint{ 
	return this.m_uiDependentSliceStartCU[pos-this.m_uiAbsIdxInLCU];                                                                                   
}
func (this *TComDataCU)  GetTotalBins            ()                            uint{ 
	return this.m_uiTotalBins;                                                                                                  
}
  // -------------------------------------------------------------------------------------------------------------------
  // member functions for RD cost storage
  // -------------------------------------------------------------------------------------------------------------------
  
func (this *TComDataCU)  GetTotalCost()                 float64 { 
	return this.m_dTotalCost;        
}
func (this *TComDataCU)  GetTotalDistortion()           uint { 
	return this.m_uiTotalDistortion; 
}
func (this *TComDataCU)  GetTotalBits()                 uint { 
	return this.m_uiTotalBits;       
}
func (this *TComDataCU)  GetTotalNumPart()              uint { 
	return this.m_uiNumPartition;    
}

func (this *TComDataCU)  GetCoefScanIdx( uiAbsPartIdx,  uiWidth uint,  bIsLuma,  bIsIntra bool) uint{
}

//namespace RasterAddress
//{
  /** Check whether 2 addresses point to the same column
   * \param addrA          First address in raster scan order
   * \param addrB          Second address in raters scan order
   * \param numUnitsPerRow Number of units in a row
   * \return Result of test
   */
  func IsEqualCol( addrA, addrB, numUnitsPerRow int ) bool{
    // addrA % numUnitsPerRow == addrB % numUnitsPerRow
    return (( addrA ^ addrB ) &  ( numUnitsPerRow - 1 ) ) == 0;
  }
  
  /** Check whether 2 addresses point to the same row
   * \param addrA          First address in raster scan order
   * \param addrB          Second address in raters scan order
   * \param numUnitsPerRow Number of units in a row
   * \return Result of test
   */
  func IsEqualRow( addrA, addrB, numUnitsPerRow int ) bool{
    // addrA / numUnitsPerRow == addrB / numUnitsPerRow
    return (( addrA ^ addrB ) & (^( numUnitsPerRow - 1 )) ) == 0;
  }
  
  /** Check whether 2 addresses point to the same row or column
   * \param addrA          First address in raster scan order
   * \param addrB          Second address in raters scan order
   * \param numUnitsPerRow Number of units in a row
   * \return Result of test
   */
  func IsEqualRowOrCol( addrA, addrB, numUnitsPerRow int ) bool{
    return IsEqualCol( addrA, addrB, numUnitsPerRow ) || IsEqualRow( addrA, addrB, numUnitsPerRow );
  }
  
  /** Check whether one address points to the first column
   * \param addr           Address in raster scan order
   * \param numUnitsPerRow Number of units in a row
   * \return Result of test
   */
  func IsZeroCol( addr, numUnitsPerRow int ) bool{
    // addr % numUnitsPerRow == 0
    return ( addr & ( numUnitsPerRow - 1 ) ) == 0;
  }
  
  /** Check whether one address points to the first row
   * \param addr           Address in raster scan order
   * \param numUnitsPerRow Number of units in a row
   * \return Result of test
   */
  func IsZeroRow( addr, numUnitsPerRow int ) bool{
    // addr / numUnitsPerRow == 0
    return ( addr & ^( numUnitsPerRow - 1 ) ) == 0;
  }
  
  /** Check whether one address points to a column whose index is smaller than a given value
   * \param addr           Address in raster scan order
   * \param val            Given column index value
   * \param numUnitsPerRow Number of units in a row
   * \return Result of test
   */
  func LessThanCol( addr, val, numUnitsPerRow int ) bool{
    // addr % numUnitsPerRow < val
    return ( addr & ( numUnitsPerRow - 1 ) ) < val;
  }
  
  /** Check whether one address points to a row whose index is smaller than a given value
   * \param addr           Address in raster scan order
   * \param val            Given row index value
   * \param numUnitsPerRow Number of units in a row
   * \return Result of test
   */
  func LessThanRow( addr, val, numUnitsPerRow int ) bool{
    // addr / numUnitsPerRow < val
    return addr < val * numUnitsPerRow;
  }
//};