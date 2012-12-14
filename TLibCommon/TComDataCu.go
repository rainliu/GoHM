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
  m_pePartSize			[]int8;         ///< array of partition sizes
  m_pePredMode			[]int8;         ///< array of prediction modes
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
  
func (this *TComDataCU) Create(  uiNumPartition,  uiWidth,  uiHeight uint, bDecSubCu bool,  unitSize int,
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

    this.m_pePartSize         = make([]int8,    uiNumPartition);
    for i:=uint(0); i<uiNumPartition; i++{
    	this.m_pePartSize[i] = SIZE_NONE;
    }
    
    this.m_pePredMode         = make([]int8,    uiNumPartition);
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
func (this *TComDataCU) 		 Destroy(){
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

func (this *TComDataCU)          InitCU                ( pcPic *TComPic,  uiCUAddr uint){
}
func (this *TComDataCU)          InitEstData           (  uiDepth uint,  qp int){
}
func (this *TComDataCU)          InitSubCU             ( pcCU *TComDataCU,  uiPartUnitIdx,  uiDepth uint,  qp int){
}
func (this *TComDataCU)          SetOutsideCUPart      (  uiAbsPartIdx,  uiDepth uint){
}

func (this *TComDataCU)          CopySubCU             ( pcCU *TComDataCU,  uiPartUnitIdx,  uiDepth uint){
}
func (this *TComDataCU)          CopyInterPredInfoFrom ( pcCU *TComDataCU,  uiAbsPartIdx uint,  eRefPicList RefPicList){
}
func (this *TComDataCU)          CopyPartFrom          ( pcCU *TComDataCU,  uiPartUnitIdx,  uiDepth uint){
}
  
func (this *TComDataCU)          CopyToPic1             (  uiDepth uint){
}
func (this *TComDataCU)          CopyToPic3             (  uiDepth,  uiPartIdx,  uiPartDepth uint){
}
  
  // -------------------------------------------------------------------------------------------------------------------
  // member functions for CU description
  // -------------------------------------------------------------------------------------------------------------------

func (this *TComDataCU) GetPic                ()             *TComPic                  { 
	return this.m_pcPic;           
}
 
func (this *TComDataCU)	GetSlice              ()           *TComSlice                   { 
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
/*  
func (this *TComDataCU)  Char*         getPartitionSize      ()                        { return m_pePartSize;        }
func (this *TComDataCU)  PartSize      getPartitionSize      ( UInt uiIdx )            { return static_cast<PartSize>( m_pePartSize[uiIdx] ); }
func (this *TComDataCU)  Void          setPartitionSize      ( UInt uiIdx, PartSize uh){ m_pePartSize[uiIdx] = uh;   }
func (this *TComDataCU)  Void          setPartSizeSubParts   ( PartSize eMode, UInt uiAbsPartIdx, UInt uiDepth );
func (this *TComDataCU)  Void          setCUTransquantBypassSubParts( Bool flag, UInt uiAbsPartIdx, UInt uiDepth );
  
func (this *TComDataCU)  Bool*        getSkipFlag            ()                        { return m_skipFlag;          }
func (this *TComDataCU)  Bool         getSkipFlag            (UInt idx)                { return m_skipFlag[idx];     }
func (this *TComDataCU)  Void         setSkipFlag           ( UInt idx, Bool skip)     { m_skipFlag[idx] = skip;   }
func (this *TComDataCU)  Void         setSkipFlagSubParts   ( Bool skip, UInt absPartIdx, UInt depth );

func (this *TComDataCU)  Char*         getPredictionMode     ()                        { return m_pePredMode;        }
func (this *TComDataCU)  PredMode      getPredictionMode     ( UInt uiIdx )            { return static_cast<PredMode>( m_pePredMode[uiIdx] ); }
func (this *TComDataCU)  Bool*         getCUTransquantBypass ()                        { return m_CUTransquantBypass;        }
func (this *TComDataCU)  Bool          getCUTransquantBypass( UInt uiIdx )             { return m_CUTransquantBypass[uiIdx]; }
func (this *TComDataCU)  Void          setPredictionMode     ( UInt uiIdx, PredMode uh){ m_pePredMode[uiIdx] = uh;   }
func (this *TComDataCU)  Void          setPredModeSubParts   ( PredMode eMode, UInt uiAbsPartIdx, UInt uiDepth );
*/  
func (this *TComDataCU)  GetWidth              () []byte                       { 
	return this.m_puhWidth;          
}
func (this *TComDataCU)  GetWidth1             (  uiIdx uint) byte           { 
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
/*  
  Void          setSizeSubParts       ( UInt uiWidth, UInt uiHeight, UInt uiAbsPartIdx, UInt uiDepth );
  
  Char*         getQP                 ()                        { return m_phQP;              }
  Char          getQP                 ( UInt uiIdx )            { return m_phQP[uiIdx];       }
  Void          setQP                 ( UInt uiIdx, Char value ){ m_phQP[uiIdx] =  value;     }
  Void          setQPSubParts         ( Int qp,   UInt uiAbsPartIdx, UInt uiDepth );
  Int           getLastValidPartIdx   ( Int iAbsPartIdx );
  Char          getLastCodedQP        ( UInt uiAbsPartIdx );
  Void          setQPSubCUs           ( Int qp, TComDataCU* pcCU, UInt absPartIdx, UInt depth, Bool &foundNonZeroCbf );
  Void          setCodedQP            ( Char qp )               { m_codedQP = qp;             }
  Char          getCodedQP            ()                        { return m_codedQP;           }

  Bool          isLosslessCoded(UInt absPartIdx);
  
  UChar*        getTransformIdx       ()                        { return m_puhTrIdx;          }
  UChar         getTransformIdx       ( UInt uiIdx )            { return m_puhTrIdx[uiIdx];   }
  Void          setTrIdxSubParts      ( UInt uiTrIdx, UInt uiAbsPartIdx, UInt uiDepth );

  UChar*        getTransformSkip      ( TextType eType)    { return m_puhTransformSkip[g_aucConvertTxtTypeToIdx[eType]];}
  UChar         getTransformSkip      ( UInt uiIdx,TextType eType)    { return m_puhTransformSkip[g_aucConvertTxtTypeToIdx[eType]][uiIdx];}
  Void          setTransformSkipSubParts  ( UInt useTransformSkip, TextType eType, UInt uiAbsPartIdx, UInt uiDepth); 
  Void          setTransformSkipSubParts  ( UInt useTransformSkipY, UInt useTransformSkipU, UInt useTransformSkipV, UInt uiAbsPartIdx, UInt uiDepth );

  UInt          getQuadtreeTULog2MinSizeInCU( UInt absPartIdx );
  
  TComCUMvField* getCUMvField         ( RefPicList e )          { return  &m_acCUMvField[e];  }
  
  TCoeff*&      getCoeffY             ()                        { return m_pcTrCoeffY;        }
  TCoeff*&      getCoeffCb            ()                        { return m_pcTrCoeffCb;       }
  TCoeff*&      getCoeffCr            ()                        { return m_pcTrCoeffCr;       }
#if ADAPTIVE_QP_SELECTION
  Int*&         getArlCoeffY          ()                        { return m_pcArlCoeffY;       }
  Int*&         getArlCoeffCb         ()                        { return m_pcArlCoeffCb;      }
  Int*&         getArlCoeffCr         ()                        { return m_pcArlCoeffCr;      }
#endif
  
  Pel*&         getPCMSampleY         ()                        { return m_pcIPCMSampleY;     }
  Pel*&         getPCMSampleCb        ()                        { return m_pcIPCMSampleCb;    }
  Pel*&         getPCMSampleCr        ()                        { return m_pcIPCMSampleCr;    }

  UChar         getCbf    ( UInt uiIdx, TextType eType )                  { return m_puhCbf[g_aucConvertTxtTypeToIdx[eType]][uiIdx];  }
  UChar*        getCbf    ( TextType eType )                              { return m_puhCbf[g_aucConvertTxtTypeToIdx[eType]];         }
  UChar         getCbf    ( UInt uiIdx, TextType eType, UInt uiTrDepth )  { return ( ( getCbf( uiIdx, eType ) >> uiTrDepth ) & 0x1 ); }
  Void          setCbf    ( UInt uiIdx, TextType eType, UChar uh )        { m_puhCbf[g_aucConvertTxtTypeToIdx[eType]][uiIdx] = uh;    }
  Void          clearCbf  ( UInt uiIdx, TextType eType, UInt uiNumParts );
  UChar         getQtRootCbf          ( UInt uiIdx )                      { return getCbf( uiIdx, TEXT_LUMA, 0 ) || getCbf( uiIdx, TEXT_CHROMA_U, 0 ) || getCbf( uiIdx, TEXT_CHROMA_V, 0 ); }
  
  Void          setCbfSubParts        ( UInt uiCbfY, UInt uiCbfU, UInt uiCbfV, UInt uiAbsPartIdx, UInt uiDepth          );
  Void          setCbfSubParts        ( UInt uiCbf, TextType eTType, UInt uiAbsPartIdx, UInt uiDepth                    );
  Void          setCbfSubParts        ( UInt uiCbf, TextType eTType, UInt uiAbsPartIdx, UInt uiPartIdx, UInt uiDepth    );
  
  // -------------------------------------------------------------------------------------------------------------------
  // member functions for coding tool information
  // -------------------------------------------------------------------------------------------------------------------
  
  Bool*         getMergeFlag          ()                        { return m_pbMergeFlag;               }
  Bool          getMergeFlag          ( UInt uiIdx )            { return m_pbMergeFlag[uiIdx];        }
  Void          setMergeFlag          ( UInt uiIdx, Bool b )    { m_pbMergeFlag[uiIdx] = b;           }
  Void          setMergeFlagSubParts  ( Bool bMergeFlag, UInt uiAbsPartIdx, UInt uiPartIdx, UInt uiDepth );

  UChar*        getMergeIndex         ()                        { return m_puhMergeIndex;                         }
  UChar         getMergeIndex         ( UInt uiIdx )            { return m_puhMergeIndex[uiIdx];                  }
  Void          setMergeIndex         ( UInt uiIdx, UInt uiMergeIndex ) { m_puhMergeIndex[uiIdx] = uiMergeIndex;  }
  Void          setMergeIndexSubParts ( UInt uiMergeIndex, UInt uiAbsPartIdx, UInt uiPartIdx, UInt uiDepth );
  template <typename T>
  Void          setSubPart            ( T bParameter, T* pbBaseLCU, UInt uiCUAddr, UInt uiCUDepth, UInt uiPUIdx );

#if AMP_MRG
  Void          setMergeAMP( Bool b )      { m_bIsMergeAMP = b; }
  Bool          getMergeAMP( )             { return m_bIsMergeAMP; }
#endif

  UChar*        getLumaIntraDir       ()                        { return m_puhLumaIntraDir;           }
  UChar         getLumaIntraDir       ( UInt uiIdx )            { return m_puhLumaIntraDir[uiIdx];    }
  Void          setLumaIntraDir       ( UInt uiIdx, UChar  uh ) { m_puhLumaIntraDir[uiIdx] = uh;      }
  Void          setLumaIntraDirSubParts( UInt uiDir,  UInt uiAbsPartIdx, UInt uiDepth );
  
  UChar*        getChromaIntraDir     ()                        { return m_puhChromaIntraDir;         }
  UChar         getChromaIntraDir     ( UInt uiIdx )            { return m_puhChromaIntraDir[uiIdx];  }
  Void          setChromaIntraDir     ( UInt uiIdx, UChar  uh ) { m_puhChromaIntraDir[uiIdx] = uh;    }
  Void          setChromIntraDirSubParts( UInt uiDir,  UInt uiAbsPartIdx, UInt uiDepth );
  
  UChar*        getInterDir           ()                        { return m_puhInterDir;               }
  UChar         getInterDir           ( UInt uiIdx )            { return m_puhInterDir[uiIdx];        }
  Void          setInterDir           ( UInt uiIdx, UChar  uh ) { m_puhInterDir[uiIdx] = uh;          }
  Void          setInterDirSubParts   ( UInt uiDir,  UInt uiAbsPartIdx, UInt uiPartIdx, UInt uiDepth );
  Bool*         getIPCMFlag           ()                        { return m_pbIPCMFlag;               }
  Bool          getIPCMFlag           (UInt uiIdx )             { return m_pbIPCMFlag[uiIdx];        }
  Void          setIPCMFlag           (UInt uiIdx, Bool b )     { m_pbIPCMFlag[uiIdx] = b;           }
  Void          setIPCMFlagSubParts   (Bool bIpcmFlag, UInt uiAbsPartIdx, UInt uiDepth);

#if !REMOVE_BURST_IPCM
  Int           getNumSucIPCM         ()                        { return m_numSucIPCM;             }
  Void          setNumSucIPCM         ( Int num )               { m_numSucIPCM = num;              }
  Bool          getLastCUSucIPCMFlag  ()                        { return m_lastCUSucIPCMFlag;        }
  Void          setLastCUSucIPCMFlag  ( Bool flg )              { m_lastCUSucIPCMFlag = flg;         }
#endif

  /// get slice ID for SU
  Int           getSUSliceID          (UInt uiIdx)              {return m_piSliceSUMap[uiIdx];      } 

  /// get the pointer of slice ID map
  Int*          getSliceSUMap         ()                        {return m_piSliceSUMap;             }

  /// set the pointer of slice ID map
  Void          setSliceSUMap         (Int *pi)                 {m_piSliceSUMap = pi;               }

  std::vector<NDBFBlockInfo>* getNDBFilterBlocks()      {return &m_vNDFBlock;}
  Void setNDBFilterBlockBorderAvailability(UInt numLCUInPicWidth, UInt numLCUInPicHeight, UInt numSUInLCUWidth, UInt numSUInLCUHeight, UInt picWidth, UInt picHeight
                                          ,std::vector<Bool>& LFCrossSliceBoundary
                                          ,Bool bTopTileBoundary, Bool bDownTileBoundary, Bool bLeftTileBoundary, Bool bRightTileBoundary
                                          ,Bool bIndependentTileBoundaryEnabled );
  // -------------------------------------------------------------------------------------------------------------------
  // member functions for accessing partition information
  // -------------------------------------------------------------------------------------------------------------------
  void          getPartIndexAndSizePos( UInt uiPartIdx, UInt& ruiPartAddr, Int& riWidth, Int& riHeight, Int& rPosX, Int& rPosY );
  Void          getPartIndexAndSize   ( UInt uiPartIdx, UInt& ruiPartAddr, Int& riWidth, Int& riHeight );
  UChar         getNumPartInter       ();
  Bool          isFirstAbsZorderIdxInDepth (UInt uiAbsPartIdx, UInt uiDepth);
  
  // -------------------------------------------------------------------------------------------------------------------
  // member functions for motion vector
  // -------------------------------------------------------------------------------------------------------------------
  
  Void          getMvField            ( TComDataCU* pcCU, UInt uiAbsPartIdx, RefPicList eRefPicList, TComMvField& rcMvField );
  
  Void          fillMvpCand           ( UInt uiPartIdx, UInt uiPartAddr, RefPicList eRefPicList, Int iRefIdx, AMVPInfo* pInfo );
  Bool          isDiffMER             ( Int xN, Int yN, Int xP, Int yP);
  Void          getPartPosition       ( UInt partIdx, Int& xP, Int& yP, Int& nPSW, Int& nPSH);
  Void          setMVPIdx             ( RefPicList eRefPicList, UInt uiIdx, Int iMVPIdx)  { m_apiMVPIdx[eRefPicList][uiIdx] = iMVPIdx;  }
  Int           getMVPIdx             ( RefPicList eRefPicList, UInt uiIdx)               { return m_apiMVPIdx[eRefPicList][uiIdx];     }
  Char*         getMVPIdx             ( RefPicList eRefPicList )                          { return m_apiMVPIdx[eRefPicList];            }

  Void          setMVPNum             ( RefPicList eRefPicList, UInt uiIdx, Int iMVPNum ) { m_apiMVPNum[eRefPicList][uiIdx] = iMVPNum;  }
  Int           getMVPNum             ( RefPicList eRefPicList, UInt uiIdx )              { return m_apiMVPNum[eRefPicList][uiIdx];     }
  Char*         getMVPNum             ( RefPicList eRefPicList )                          { return m_apiMVPNum[eRefPicList];            }
  
  Void          setMVPIdxSubParts     ( Int iMVPIdx, RefPicList eRefPicList, UInt uiAbsPartIdx, UInt uiPartIdx, UInt uiDepth );
  Void          setMVPNumSubParts     ( Int iMVPNum, RefPicList eRefPicList, UInt uiAbsPartIdx, UInt uiPartIdx, UInt uiDepth );
  
  Void          clipMv                ( TComMv&     rcMv     );
  Void          getMvPredLeft         ( TComMv&     rcMvPred )   { rcMvPred = m_cMvFieldA.getMv(); }
  Void          getMvPredAbove        ( TComMv&     rcMvPred )   { rcMvPred = m_cMvFieldB.getMv(); }
  Void          getMvPredAboveRight   ( TComMv&     rcMvPred )   { rcMvPred = m_cMvFieldC.getMv(); }
  
  Void          compressMV            ();
 */ 
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

/*
  TComDataCU*   getPULeft                   ( UInt&  uiLPartUnitIdx, 
                                              UInt uiCurrPartUnitIdx, 
                                              Bool bEnforceSliceRestriction=true, 
                                              Bool bEnforceDependentSliceRestriction=true,
                                              Bool bEnforceTileRestriction=true );
#if !LINEBUF_CLEANUP
  TComDataCU*   getPUAbove                  ( UInt&  uiAPartUnitIdx, 
                                              UInt uiCurrPartUnitIdx, 
                                              Bool bEnforceSliceRestriction=true, 
                                              Bool bEnforceDependentSliceRestriction=true, 
                                              Bool MotionDataCompresssion = false,
                                              Bool planarAtLCUBoundary = false,
                                              Bool bEnforceTileRestriction=true );
  TComDataCU*   getPUAboveLeft              ( UInt&  uiALPartUnitIdx, UInt uiCurrPartUnitIdx, Bool bEnforceSliceRestriction=true, Bool bEnforceDependentSliceRestriction=true, Bool MotionDataCompresssion = false );
  TComDataCU*   getPUAboveRight             ( UInt&  uiARPartUnitIdx, UInt uiCurrPartUnitIdx, Bool bEnforceSliceRestriction=true, Bool bEnforceDependentSliceRestriction=true, Bool MotionDataCompresssion = false );
#else
  TComDataCU*   getPUAbove                  ( UInt&  uiAPartUnitIdx, 
                                              UInt uiCurrPartUnitIdx, 
                                              Bool bEnforceSliceRestriction=true, 
                                              Bool bEnforceDependentSliceRestriction=true, 
                                              Bool planarAtLCUBoundary = false,
                                              Bool bEnforceTileRestriction=true );
  TComDataCU*   getPUAboveLeft              ( UInt&  uiALPartUnitIdx, UInt uiCurrPartUnitIdx, Bool bEnforceSliceRestriction=true, Bool bEnforceDependentSliceRestriction=true );
  TComDataCU*   getPUAboveRight             ( UInt&  uiARPartUnitIdx, UInt uiCurrPartUnitIdx, Bool bEnforceSliceRestriction=true, Bool bEnforceDependentSliceRestriction=true );
#endif
  TComDataCU*   getPUBelowLeft              ( UInt&  uiBLPartUnitIdx, UInt uiCurrPartUnitIdx, Bool bEnforceSliceRestriction=true, Bool bEnforceDependentSliceRestriction=true );

  TComDataCU*   getQpMinCuLeft              ( UInt&  uiLPartUnitIdx , UInt uiCurrAbsIdxInLCU );
  TComDataCU*   getQpMinCuAbove             ( UInt&  aPartUnitIdx , UInt currAbsIdxInLCU );
  Char          getRefQP                    ( UInt   uiCurrAbsIdxInLCU                       );

  TComDataCU*   getPUAboveRightAdi          ( UInt&  uiARPartUnitIdx, UInt uiCurrPartUnitIdx, UInt uiPartUnitOffset = 1, Bool bEnforceSliceRestriction=true, Bool bEnforceDependentSliceRestriction=true );
  TComDataCU*   getPUBelowLeftAdi           ( UInt&  uiBLPartUnitIdx, UInt uiCurrPartUnitIdx, UInt uiPartUnitOffset = 1, Bool bEnforceSliceRestriction=true, Bool bEnforceDependentSliceRestriction=true );
  
  Void          deriveLeftRightTopIdx       ( UInt uiPartIdx, UInt& ruiPartIdxLT, UInt& ruiPartIdxRT );
  Void          deriveLeftBottomIdx         ( UInt uiPartIdx, UInt& ruiPartIdxLB );
  
  Void          deriveLeftRightTopIdxAdi    ( UInt& ruiPartIdxLT, UInt& ruiPartIdxRT, UInt uiPartOffset, UInt uiPartDepth );
  Void          deriveLeftBottomIdxAdi      ( UInt& ruiPartIdxLB, UInt  uiPartOffset, UInt uiPartDepth );
  
  Bool          hasEqualMotion              ( UInt uiAbsPartIdx, TComDataCU* pcCandCU, UInt uiCandAbsPartIdx );
  Void          getInterMergeCandidates       ( UInt uiAbsPartIdx, UInt uiPUIdx, TComMvField* pcMFieldNeighbours, UChar* puhInterDirNeighbours, Int& numValidMergeCand, Int mrgCandIdx = -1 );
  Void          deriveLeftRightTopIdxGeneral  ( UInt uiAbsPartIdx, UInt uiPartIdx, UInt& ruiPartIdxLT, UInt& ruiPartIdxRT );
  Void          deriveLeftBottomIdxGeneral    ( UInt uiAbsPartIdx, UInt uiPartIdx, UInt& ruiPartIdxLB );
  
  
  // -------------------------------------------------------------------------------------------------------------------
  // member functions for modes
  // -------------------------------------------------------------------------------------------------------------------
  
  Bool          isIntra   ( UInt uiPartIdx )  { return m_pePredMode[ uiPartIdx ] == MODE_INTRA; }
  Bool          isSkipped ( UInt uiPartIdx );                                                     ///< SKIP (no residual)
  Bool          isBipredRestriction( UInt puIdx );

  // -------------------------------------------------------------------------------------------------------------------
  // member functions for symbol prediction (most probable / mode conversion)
  // -------------------------------------------------------------------------------------------------------------------
  
  UInt          getIntraSizeIdx                 ( UInt uiAbsPartIdx                                       );
  
  Void          getAllowedChromaDir             ( UInt uiAbsPartIdx, UInt* uiModeList );
  Int           getIntraDirLumaPredictor        ( UInt uiAbsPartIdx, Int* uiIntraDirPred, Int* piMode = nil );
  
  // -------------------------------------------------------------------------------------------------------------------
  // member functions for SBAC context
  // -------------------------------------------------------------------------------------------------------------------
  
  UInt          getCtxSplitFlag                 ( UInt   uiAbsPartIdx, UInt uiDepth                   );
  UInt          getCtxQtCbf                     ( TextType eType, UInt uiTrDepth );

  UInt          getCtxSkipFlag                  ( UInt   uiAbsPartIdx                                 );
  UInt          getCtxInterDir                  ( UInt   uiAbsPartIdx                                 );
  
  UInt          getSliceStartCU         ( UInt pos )                  { return m_uiSliceStartCU[pos-m_uiAbsIdxInLCU];                                                                                          }
  UInt          getDependentSliceStartCU  ( UInt pos )                  { return m_uiDependentSliceStartCU[pos-m_uiAbsIdxInLCU];                                                                                   }
  UInt&         getTotalBins            ()                            { return m_uiTotalBins;                                                                                                  }
  // -------------------------------------------------------------------------------------------------------------------
  // member functions for RD cost storage
  // -------------------------------------------------------------------------------------------------------------------
  
  Double&       getTotalCost()                  { return m_dTotalCost;        }
  UInt&         getTotalDistortion()            { return m_uiTotalDistortion; }
  UInt&         getTotalBits()                  { return m_uiTotalBits;       }
  UInt&         getTotalNumPart()               { return m_uiNumPartition;    }

  UInt          getCoefScanIdx(UInt uiAbsPartIdx, UInt uiWidth, Bool bIsLuma, Bool bIsIntra);

};
*/
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