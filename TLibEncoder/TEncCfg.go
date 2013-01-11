package TLibEncoder

import (
	"fmt"
	"gohm/TLibCommon"
)

type GOPEntry struct{
  m_POC				int;
  m_QPOffset		int;
  m_QPFactor		float64;
//#if VARYING_DBL_PARAMS
  m_tcOffsetDiv2	int;
  m_betaOffsetDiv2	int;
//#endif
  m_temporalId		int;
  m_refPic			bool;
  m_numRefPicsActive	int;
  m_sliceType			string;
  m_numRefPics			int;
  m_referencePics	[TLibCommon.MAX_NUM_REF_PICS]int;
  m_usedByCurrPic	[TLibCommon.MAX_NUM_REF_PICS]bool;
//#if AUTO_INTER_RPS
  m_interRPSPrediction	int;
/*#else
  Bool m_interRPSPrediction;
#endif*/
  m_deltaRPS	int;
  m_numRefIdc	int;
  m_refIdc	[TLibCommon.MAX_NUM_REF_PICS+1]int;
}

func NewGOPEntry() *GOPEntry{
	gop := &GOPEntry{ m_POC:-1,
					  m_QPOffset:0,
					  m_QPFactor:0,
					//#if VARYING_DBL_PARAMS
					  m_tcOffsetDiv2:0,
					  m_betaOffsetDiv2:0,
					//#endif
					  m_temporalId:0,
					  m_refPic:false,
					  m_numRefPicsActive:0,
					  m_sliceType:"P",
					  m_numRefPics:0,
					  m_interRPSPrediction:0,
					  m_deltaRPS:0,
					  m_numRefIdc:0};
	
	return gop;
}

func (this *GOPEntry) GetPOC() int{
	return this.m_POC;
}

func (this *GOPEntry) SetPOC(poc int){
	this.m_POC = poc;
}

func (this *GOPEntry) GetQPOffset() int{
	return this.m_QPOffset;
}

func (this *GOPEntry) SetQPOffset(QPOffset int){
	this.m_QPOffset = QPOffset;
}

func (this *GOPEntry) GetQPFactor() float64{
	return this.m_QPFactor;
}

func (this *GOPEntry) SetQPFactor(QPFactor float64){
	this.m_QPFactor = QPFactor;
}

func (this *GOPEntry) GetBetaOffsetDiv2() int {
	return this.m_betaOffsetDiv2;
}

func (this *GOPEntry) GetTcOffsetDiv2() int{
	return this.m_tcOffsetDiv2;
}

func (this *GOPEntry) SetBetaOffsetDiv2(betaOffsetDiv2 int) {
	this.m_betaOffsetDiv2 = betaOffsetDiv2;
}

func (this *GOPEntry) SetTcOffsetDiv2(tcOffsetDiv2 int){
	this.m_tcOffsetDiv2 = tcOffsetDiv2;
}

func (this *GOPEntry) GetNumRefPicsActive() int{
	return this.m_numRefPicsActive;
}

func (this *GOPEntry) SetNumRefPicsActive(numRefPicsActive int){
	this.m_numRefPicsActive = numRefPicsActive;
}

func (this *GOPEntry) GetTemporalId() int{
	return this.m_temporalId;
}
func (this *GOPEntry) SetTemporalId(temporalId int){
	this.m_temporalId=temporalId;
}

func (this *GOPEntry) SetNumRefIdc(numRefIdc int){
	this.m_numRefIdc = numRefIdc;
}
func (this *GOPEntry) GetNumRefIdc() int{
	return this.m_numRefIdc;
}

func (this *GOPEntry) GetNumRefPics() int{
	return this.m_numRefPics;
}
func (this *GOPEntry) SetNumRefPics(numRefPics int){
	this.m_numRefPics = numRefPics;
}

func (this *GOPEntry) GetReferencePics(i int) int{
	return this.m_referencePics[i];
}

func (this *GOPEntry) SetReferencePics(i int, value int){
    this.m_referencePics[i] = value;
}

func (this *GOPEntry) SetRefPic(refPic bool) {
	this.m_refPic = refPic;
}

func (this *GOPEntry) GetRefPic() bool {
	return this.m_refPic;
}

func (this *GOPEntry) SetUsedByCurrPic(i int, b bool){
	this.m_usedByCurrPic[i] = b;
}

func (this *GOPEntry) GetUsedByCurrPic(i int) bool{
	return this.m_usedByCurrPic[i];
}

func (this *GOPEntry) SetInterRPSPrediction(interRPSPrediction int) {
	this.m_interRPSPrediction = interRPSPrediction;
}

func (this *GOPEntry) GetInterRPSPrediction() int {
	return this.m_interRPSPrediction;
}

func (this *GOPEntry) SetRefIdc(i int, refIdc int){
	this.m_refIdc[i] = refIdc;
}

func (this *GOPEntry) GetRefIdc(i int) int{
	return this.m_refIdc[i];
}

func (this *GOPEntry) SetDeltaRPS(deltaRPS int) {
	this.m_deltaRPS = deltaRPS;
}

func (this *GOPEntry) GetDeltaRPS() int {
	return this.m_deltaRPS;
}

func (this *GOPEntry) GetSliceType() string{
	return this.m_sliceType;
}

func (this *GOPEntry) SetSliceType(sliceType string){
	this.m_sliceType = sliceType;
}


// ====================================================================================================================
// Class definition
// ====================================================================================================================

/// encoder configuration class
type TEncCfg struct{
//protected:
  //==== File I/O ========
  m_iFrameRate	int;
  m_FrameSkip		uint;
  m_iSourceWidth	int;
  m_iSourceHeight	int;
  m_croppingMode	int;
  m_picCroppingWindow	*TLibCommon.CroppingWindow;
  m_iFrameToBeEncoded	int;
  m_adLambdaModifier	[ TLibCommon.MAX_TLAYER ]float64;

  /* profile & level */
  m_profile		TLibCommon.PROFILE;
  m_levelTier	TLibCommon.TIER;
  m_level		TLibCommon.LEVEL;

  //====== Coding Structure ========
  m_uiIntraPeriod	uint;
  m_uiDecodingRefreshType	uint;            ///< the type of decoding refresh employed for the random access.
  m_iGOPSize	int;
  m_GOPList	[TLibCommon.MAX_GOP]*GOPEntry;
  m_extraRPSs	int;
  m_maxDecPicBuffering	[TLibCommon.MAX_TLAYER]int;
  m_numReorderPics		[TLibCommon.MAX_TLAYER]int;
  
  m_iQP	int;                              //  if (AdaptiveQP == OFF)
  
  m_aiPad	[2]int;
  

  m_iMaxRefPicNum	int;                     ///< this is used to mimic the sliding mechanism used by the decoder
                                       // TODO: We need to have a common sliding mechanism used by both the encoder and decoder

  m_maxTempLayer	int;                      ///< Max temporal layer
  m_useAMP	bool;
  //======= Transform =============
  m_uiQuadtreeTULog2MaxSize	uint;
  m_uiQuadtreeTULog2MinSize	uint;
  m_uiQuadtreeTUMaxDepthInter	uint;
  m_uiQuadtreeTUMaxDepthIntra	uint;
  
  //====== Loop/Deblock Filter ========
  m_bLoopFilterDisable	bool;
  m_loopFilterOffsetInPPS	bool;
  m_loopFilterBetaOffsetDiv2	int;
  m_loopFilterTcOffsetDiv2	int;
  m_DeblockingFilterControlPresent	bool;
  m_bUseSAO	bool;
  m_maxNumOffsetsPerPic	int;
  m_saoLcuBoundary	bool;
  m_saoLcuBasedOptimization	bool;

  //====== Lossless ========
  m_useLossless	bool;
  //====== Motion search ========
  m_iFastSearch	int;                      //  0:Full search  1:Diamond  2:PMVFAST
  m_iSearchRange	int;                     //  0:Full frame
  m_bipredSearchRange int;

  //====== Quality control ========
  m_iMaxDeltaQP	int;                      //  Max. absolute delta QP (1:default)
  m_iMaxCuDQPDepth	int;                   //  Max. depth for a minimum CuDQP (0:default)

  m_chromaCbQpOffset	int;                 //  Chroma Cb QP Offset (0:default)
  m_chromaCrQpOffset	int;                 //  Chroma Cr Qp Offset (0:default)

//#if ADAPTIVE_QP_SELECTION
  m_bUseAdaptQpSelect	bool;
//#endif

  m_bUseAdaptiveQP	bool;
  m_iQPAdaptationRange	int;
  
  //====== Tool list ========
  m_bUseSBACRD	bool;
  m_bUseASR	bool;
  m_bUseHADME	bool;
  m_bUseLComb	bool;
  m_useRDOQ	bool;
//#if RDOQ_TRANSFORMSKIP
  m_useRDOQTS	bool;
//#endif
  m_bUseFastEnc	bool;
  m_bUseEarlyCU	bool;
  m_useFastDecisionForMerge	bool;
  m_bUseCbfFastMode	bool;
  m_useEarlySkipDetection	bool;
  m_useTransformSkip	bool;
  m_useTransformSkipFast	bool;
  m_aidQP	[]int;
  m_uiDeltaQpRD	uint;
  
  m_bUseConstrainedIntraPred	bool;
  m_usePCM	bool;
  m_pcmLog2MaxSize	uint;
  m_uiPCMLog2MinSize	uint;
  //====== Slice ========
  m_iSliceMode	int;
  m_iSliceArgument	int;	 
  //====== Dependent Slice ========
  m_iDependentSliceMode	int;
  m_iDependentSliceArgument	int;
//#if DEPENDENT_SLICES
  m_entropySliceEnabledFlag	bool;
//#endif
  m_bLFCrossSliceBoundaryFlag	bool;

  m_bPCMInputBitDepthFlag	bool;
  m_uiPCMBitDepthLuma	uint;
  m_uiPCMBitDepthChroma	uint;
  m_bPCMFilterDisableFlag	bool;
  m_loopFilterAcrossTilesEnabledFlag	bool;
  m_iUniformSpacingIdr	int;
  m_iNumColumnsMinus1	int;
  m_puiColumnWidth	[]uint;
  m_iNumRowsMinus1	int;
  m_puiRowHeight	[]uint;

  m_iWaveFrontSynchro	int;
  m_iWaveFrontSubstreams	int;

  m_decodedPictureHashSEIEnabled	int;              ///< Checksum(3)/CRC(2)/MD5(1)/disable(0) acting on decoded picture hash SEI message
  m_bufferingPeriodSEIEnabled	int;
  m_pictureTimingSEIEnabled	int;
  m_recoveryPointSEIEnabled	int;
//#if SEI_DISPLAY_ORIENTATION
  m_displayOrientationSEIAngle	int;
//#endif
//#if SEI_TEMPORAL_LEVEL0_INDEX
  m_temporalLevel0IndexSEIEnabled	int;
//#endif
  //====== Weighted Prediction ========
  m_bUseWeightPred	bool;       //< Use of Weighting Prediction (P_SLICE)
  m_useWeightedBiPred	bool;    //< Use of Bi-directional Weighting Prediction (B_SLICE)
  m_log2ParallelMergeLevelMinus2	uint;       ///< Parallel merge estimation region
  m_maxNumMergeCand	uint;                    ///< Maximum number of merge candidates
  m_useScalingListId	int;            ///< Using quantization matrix i.e. 0=off, 1=default, 2=file.
  m_scalingListFile	string;          ///< quantization matrix file name
  m_TMVPModeId	int;
  m_signHideFlag	int;
//#if RATE_CONTROL_LAMBDA_DOMAIN
  m_RCEnableRateControl	bool;
  m_RCTargetBitrate	int;
  m_RCKeepHierarchicalBit	bool;
  m_RCLCULevelRC	bool;
  m_RCUseLCUSeparateModel	bool;
  m_RCInitialQP	int;
  m_RCForceIntraQP	bool;
/*#else
  Bool      m_enableRateCtrl;                                ///< Flag for using rate control algorithm
  Int       m_targetBitrate;                                 ///< target bitrate
  Int       m_numLCUInUnit;                                  ///< Total number of LCUs in a frame should be divided by the NumLCUInUnit
#endif*/
  m_TransquantBypassEnableFlag	bool;                     ///< transquant_bypass_enable_flag setting in PPS.
  m_CUTransquantBypassFlagValue	bool;                    ///< if transquant_bypass_enable_flag, the fixed value to use for the per-CU cu_transquant_bypass_flag.
  m_cVPS			*TLibCommon.TComVPS;
  m_recalculateQPAccordingToLambda	bool;                 ///< recalculate QP value according to the lambda value
  m_activeParameterSetsSEIEnabled	int;                  ///< enable active parameter set SEI message 
  m_vuiParametersPresentFlag	bool;                       ///< enable generation of VUI parameters
  m_aspectRatioInfoPresentFlag	bool;                     ///< Signals whether aspect_ratio_idc is present
  m_aspectRatioIdc	int;                                 ///< aspect_ratio_idc
  m_sarWidth	int;                                       ///< horizontal size of the sample aspect ratio
  m_sarHeight	int;                                      ///< vertical size of the sample aspect ratio
  m_overscanInfoPresentFlag	bool;                        ///< Signals whether overscan_appropriate_flag is present
  m_overscanAppropriateFlag	bool;                        ///< Indicates whether cropped decoded pictures are suitable for display using overscan
  m_videoSignalTypePresentFlag	bool;                     ///< Signals whether video_format, video_full_range_flag, and colour_description_present_flag are present
  m_videoFormat	int;                                    ///< Indicates representation of pictures
  m_videoFullRangeFlag	bool;                             ///< Indicates the black level and range of luma and chroma signals
  m_colourDescriptionPresentFlag	bool;                   ///< Signals whether colour_primaries, transfer_characteristics and matrix_coefficients are present
  m_colourPrimaries	int;                                ///< Indicates chromaticity coordinates of the source primaries
  m_transferCharacteristics	int;                        ///< Indicates the opto-electronic transfer characteristics of the source
  m_matrixCoefficients	int;                             ///< Describes the matrix coefficients used in deriving luma and chroma from RGB primaries
  m_chromaLocInfoPresentFlag	bool;                       ///< Signals whether chroma_sample_loc_type_top_field and chroma_sample_loc_type_bottom_field are present
  m_chromaSampleLocTypeTopField	int;                    ///< Specifies the location of chroma samples for top field
  m_chromaSampleLocTypeBottomField	int;                 ///< Specifies the location of chroma samples for bottom field
  m_neutralChromaIndicationFlag	bool;                    ///< Indicates that the value of all decoded chroma samples is equal to 1<<(BitDepthCr-1)
  m_bitstreamRestrictionFlag	bool;                       ///< Signals whether bitstream restriction parameters are present
  m_tilesFixedStructureFlag	bool;                        ///< Indicates that each active picture parameter set has the same values of the syntax elements related to tiles
  m_motionVectorsOverPicBoundariesFlag	bool;             ///< Indicates that no samples outside the picture boundaries are used for inter prediction
//#if MIN_SPATIAL_SEGMENTATION
  m_minSpatialSegmentationIdc	int;                      ///< Indicates the maximum size of the spatial segments in the pictures in the coded video sequence
//#endif
  m_maxBytesPerPicDenom	int;                            ///< Indicates a number of bytes not exceeded by the sum of the sizes of the VCL NAL units associated with any coded picture
  m_maxBitsPerMinCuDenom	int;                           ///< Indicates an upper bound for the number of bits of coding_unit() data
  m_log2MaxMvLengthHorizontal	int;                      ///< Indicate the maximum absolute value of a decoded horizontal MV component in quarter-pel luma units
  m_log2MaxMvLengthVertical	int;                        ///< Indicate the maximum absolute value of a decoded vertical MV component in quarter-pel luma units

//#if STRONG_INTRA_SMOOTHING
  m_useStrongIntraSmoothing	bool;                        ///< enable the use of strong intra smoothing (bi_linear interpolation) for 32x32 blocks when reference samples are flat.
//#endif
}

//public:
func NewTEncCfg() *TEncCfg {
	return &TEncCfg{};  
}

  
func (this *TEncCfg)  setProfile(profile TLibCommon.PROFILE) { this.m_profile = profile; }
func (this *TEncCfg)  setLevel(tier TLibCommon.TIER, level TLibCommon.LEVEL) { this.m_levelTier = tier; this.m_level = level; }

func (this *TEncCfg)  setFrameRate                    ( i int )      { this.m_iFrameRate = i; }
func (this *TEncCfg)  setFrameSkip                    ( i uint )     { this.m_FrameSkip = i; }
func (this *TEncCfg)  setSourceWidth                  ( i int )      { this.m_iSourceWidth = i; }
func (this *TEncCfg)  setSourceHeight                 ( i int )      { this.m_iSourceHeight = i; }

func (this *TEncCfg)  getPicCroppingWindow() *TLibCommon.CroppingWindow         { return this.m_picCroppingWindow; }
func (this *TEncCfg)  setPicCroppingWindow ( cropLeft,  cropRight,  cropTop,  cropBottom int) { this.m_picCroppingWindow.SetPicCropping (cropLeft, cropRight, cropTop, cropBottom); }

func (this *TEncCfg)  setFrameToBeEncoded             ( i int )      { this.m_iFrameToBeEncoded = i; }
  
  //====== Coding Structure ========
func (this *TEncCfg)  setIntraPeriod                  ( i int )      { this.m_uiIntraPeriod = uint(i); }
func (this *TEncCfg)  setDecodingRefreshType          ( i int )      { this.m_uiDecodingRefreshType = uint(i); }
func (this *TEncCfg)  setGOPSize                      ( i int )      { this.m_iGOPSize = i; }
func (this *TEncCfg)  setGopList                      ( GOPList []*GOPEntry) {  
	for i := 0; i < TLibCommon.MAX_GOP; i++ {
	 	this.m_GOPList[i] = GOPList[i]; 
	}
}
func (this *TEncCfg)  setExtraRPSs                    ( i int )      { this.m_extraRPSs = i; }
func (this *TEncCfg)  getGOPEntry                     ( i int ) *GOPEntry     { return this.m_GOPList[i]; }
func (this *TEncCfg)  setMaxDecPicBuffering           ( u, tlayer uint ) { this.m_maxDecPicBuffering[tlayer] = int(u);    }
func (this *TEncCfg)  setNumReorderPics               ( i int, tlayer uint ) { this.m_numReorderPics[tlayer] = i;    }
  
func (this *TEncCfg)  setQP                           ( i int)      { this.m_iQP = i; }
  
func (this *TEncCfg)  setPad                          ( iPad  []int                 )      { 
	for i := 0; i < 2; i++ {
		this.m_aiPad[i] = iPad[i]; 
	}
}
func (this *TEncCfg)  getMaxRefPicNum                 ()   int                           { return this.m_iMaxRefPicNum;           }
func (this *TEncCfg)  setMaxRefPicNum                 ( iMaxRefPicNum int )           { this.m_iMaxRefPicNum = iMaxRefPicNum;  }

func (this *TEncCfg)  getMaxTempLayer                 ()    int                          { return this.m_maxTempLayer;              } 
func (this *TEncCfg)  setMaxTempLayer                 ( maxTempLayer int)            { this.m_maxTempLayer = maxTempLayer;      }
  //======== Transform =============
func (this *TEncCfg)  setQuadtreeTULog2MaxSize        ( u uint )      { this.m_uiQuadtreeTULog2MaxSize = u; }
func (this *TEncCfg)  setQuadtreeTULog2MinSize        ( u uint )      { this.m_uiQuadtreeTULog2MinSize = u; }
func (this *TEncCfg)  setQuadtreeTUMaxDepthInter      ( u uint )      { this.m_uiQuadtreeTUMaxDepthInter = u; }
func (this *TEncCfg)  setQuadtreeTUMaxDepthIntra      ( u uint )      { this.m_uiQuadtreeTUMaxDepthIntra = u; }
  
func (this *TEncCfg)  setUseAMP( b bool ) { this.m_useAMP = b; }
  
  //====== Loop/Deblock Filter ========
func (this *TEncCfg)  setLoopFilterDisable            ( b bool )      { this.m_bLoopFilterDisable       = b; }
func (this *TEncCfg)  setLoopFilterOffsetInPPS        ( b bool )      { this.m_loopFilterOffsetInPPS      = b; }
func (this *TEncCfg)  setLoopFilterBetaOffset         ( i int)      { this.m_loopFilterBetaOffsetDiv2  = i; }
func (this *TEncCfg)  setLoopFilterTcOffset           ( i int )      { this.m_loopFilterTcOffsetDiv2    = i; }
func (this *TEncCfg)  setDeblockingFilterControlPresent ( b bool ) { this.m_DeblockingFilterControlPresent = b; }

  //====== Motion search ========
func (this *TEncCfg)  setFastSearch                   ( i int )      { this.m_iFastSearch = i; }
func (this *TEncCfg)  setSearchRange                  ( i int )      { this.m_iSearchRange = i; }
func (this *TEncCfg)  setBipredSearchRange            ( i int )      { this.m_bipredSearchRange = i; }

  //====== Quality control ========
func (this *TEncCfg)  setMaxDeltaQP                   ( i int )      { this.m_iMaxDeltaQP = i; }
func (this *TEncCfg)  setMaxCuDQPDepth                ( i int )      { this.m_iMaxCuDQPDepth = i; }

func (this *TEncCfg)  setChromaCbQpOffset             ( i int )      { this.m_chromaCbQpOffset = i; }
func (this *TEncCfg)  setChromaCrQpOffset             ( i int )      { this.m_chromaCrQpOffset = i; }

//#if ADAPTIVE_QP_SELECTION
func (this *TEncCfg)  setUseAdaptQpSelect             (i bool) { this.m_bUseAdaptQpSelect    = i; }
func (this *TEncCfg)  getUseAdaptQpSelect             ()  bool         { return   this.m_bUseAdaptQpSelect; }
//#endif

func (this *TEncCfg)  setUseAdaptiveQP                ( b bool )      { this.m_bUseAdaptiveQP = b; }
func (this *TEncCfg)  setQPAdaptationRange            ( i int )      { this.m_iQPAdaptationRange = i; }
  
  //====== Lossless ========
func (this *TEncCfg)  setUseLossless                  ( b bool )        { this.m_useLossless = b;  }
  //====== Sequence ========
func (this *TEncCfg)  getFrameRate                    ()  int     { return  this.m_iFrameRate; }
func (this *TEncCfg)  getFrameSkip                    ()  uint    { return  this.m_FrameSkip; }
func (this *TEncCfg)  getSourceWidth                  ()  int     { return  this.m_iSourceWidth; }
func (this *TEncCfg)  getSourceHeight                 ()  int     { return  this.m_iSourceHeight; }
func (this *TEncCfg)  getFrameToBeEncoded             ()  int     { return  this.m_iFrameToBeEncoded; }
func (this *TEncCfg)  setLambdaModifier                    ( uiIndex uint,  dValue float64) { this.m_adLambdaModifier[ uiIndex ] = dValue; }
func (this *TEncCfg)  getLambdaModifier                  ( uiIndex uint ) float64 { return this.m_adLambdaModifier[ uiIndex ]; }

  //==== Coding Structure ========
func (this *TEncCfg)  getIntraPeriod                  ()   uint   { return  this.m_uiIntraPeriod; }
func (this *TEncCfg)  getDecodingRefreshType          ()   uint   { return  this.m_uiDecodingRefreshType; }
func (this *TEncCfg)  getGOPSize                      ()   int    { return  this.m_iGOPSize; }
func (this *TEncCfg)  getMaxDecPicBuffering           ( tlayer uint) int { return this.m_maxDecPicBuffering[tlayer]; }
func (this *TEncCfg)  getNumReorderPics               ( tlayer uint) int { return this.m_numReorderPics[tlayer]; }
func (this *TEncCfg)  getQP                           ()   int    { return  this.m_iQP; }
  
func (this *TEncCfg)  getPad                          ( i int ) int     {  return  this.m_aiPad[i]; }
  
  //======== Transform =============
func (this *TEncCfg)  getQuadtreeTULog2MaxSize        ()      uint { return this.m_uiQuadtreeTULog2MaxSize; }
func (this *TEncCfg)  getQuadtreeTULog2MinSize        ()      uint { return this.m_uiQuadtreeTULog2MinSize; }
func (this *TEncCfg)  getQuadtreeTUMaxDepthInter      ()      uint { return this.m_uiQuadtreeTUMaxDepthInter; }
func (this *TEncCfg)  getQuadtreeTUMaxDepthIntra      ()      uint { return this.m_uiQuadtreeTUMaxDepthIntra; }
 
  //==== Loop/Deblock Filter ========
func (this *TEncCfg)  getLoopFilterDisable            ()      bool{ return  this.m_bLoopFilterDisable;       }
func (this *TEncCfg)  getLoopFilterOffsetInPPS        ()      bool{ return this.m_loopFilterOffsetInPPS; }
func (this *TEncCfg)  getLoopFilterBetaOffset         ()      int { return this.m_loopFilterBetaOffsetDiv2; }
func (this *TEncCfg)  getLoopFilterTcOffset           ()      int { return this.m_loopFilterTcOffsetDiv2; }
func (this *TEncCfg)  getDeblockingFilterControlPresent()  	bool{ return  this.m_DeblockingFilterControlPresent; }

  //==== Motion search ========
func (this *TEncCfg)  getFastSearch                   ()      int{ return  this.m_iFastSearch; }
func (this *TEncCfg)  getSearchRange                  ()      int{ return  this.m_iSearchRange; }

  //==== Quality control ========
func (this *TEncCfg)  getMaxDeltaQP                   ()      int{ return  this.m_iMaxDeltaQP; }
func (this *TEncCfg)  getMaxCuDQPDepth                ()      int{ return  this.m_iMaxCuDQPDepth; }
func (this *TEncCfg)  getUseAdaptiveQP                ()      bool{ return  this.m_bUseAdaptiveQP; }
func (this *TEncCfg)  getQPAdaptationRange            ()      int{ return  this.m_iQPAdaptationRange; }
  //====== Lossless ========
func (this *TEncCfg)  getUseLossless                  ()      bool{ return  this.m_useLossless;  }
  
  //==== Tool list ========
func (this *TEncCfg)  setUseSBACRD                    ( b bool )     { this.m_bUseSBACRD  = b; }
func (this *TEncCfg)  setUseASR                       ( b bool )     { this.m_bUseASR     = b; }
func (this *TEncCfg)  setUseHADME                     ( b bool )     { this.m_bUseHADME   = b; }
func (this *TEncCfg)  setUseLComb                     ( b bool )     { this.m_bUseLComb   = b; }
func (this *TEncCfg)  setUseRDOQ                      ( b bool )     { this.m_useRDOQ    = b; }
//#if RDOQ_TRANSFORMSKIP
func (this *TEncCfg)  setUseRDOQTS                    ( b bool )     { this.m_useRDOQTS  = b; }
//#endif
func (this *TEncCfg)  setUseFastEnc                   ( b bool )     { this.m_bUseFastEnc = b; }
func (this *TEncCfg)  setUseEarlyCU                   ( b bool )     { this.m_bUseEarlyCU = b; }
func (this *TEncCfg)  setUseFastDecisionForMerge      ( b bool )     { this.m_useFastDecisionForMerge = b; }
func (this *TEncCfg)  setUseCbfFastMode            ( b bool )     { this.m_bUseCbfFastMode = b; }
func (this *TEncCfg)  setUseEarlySkipDetection        ( b bool )     { this.m_useEarlySkipDetection = b; }
func (this *TEncCfg)  setUseConstrainedIntraPred      ( b bool )     { this.m_bUseConstrainedIntraPred = b; }
func (this *TEncCfg)  setPCMInputBitDepthFlag         ( b bool )     { this.m_bPCMInputBitDepthFlag = b; }
func (this *TEncCfg)  setPCMFilterDisableFlag         ( b bool )     {  this.m_bPCMFilterDisableFlag = b; }
func (this *TEncCfg)  setUsePCM                       ( b bool )     {  this.m_usePCM = b;               }
func (this *TEncCfg)  setPCMLog2MaxSize               ( u uint )      { this.m_pcmLog2MaxSize = u;      }
func (this *TEncCfg)  setPCMLog2MinSize               ( u uint )     { this.m_uiPCMLog2MinSize = u;      }
func (this *TEncCfg)  setdQPs                         ( p []int )     { this.m_aidQP       = p; }
func (this *TEncCfg)  setDeltaQpRD                    ( u uint )     {this.m_uiDeltaQpRD  = u; }
func (this *TEncCfg)  getUseSBACRD                    ()   bool   { return this.m_bUseSBACRD;  }
func (this *TEncCfg)  getUseASR                       ()   bool   { return this.m_bUseASR;     }
func (this *TEncCfg)  getUseHADME                     ()   bool   { return this.m_bUseHADME;   }
func (this *TEncCfg)  getUseLComb                     ()   bool   { return this.m_bUseLComb;   }
func (this *TEncCfg)  getUseRDOQ                      ()   bool   { return this.m_useRDOQ;    }
//#if RDOQ_TRANSFORMSKIP
func (this *TEncCfg)  getUseRDOQTS                    ()   bool   { return this.m_useRDOQTS;  }
//#endif
func (this *TEncCfg)  getUseFastEnc                   ()  bool   { return this.m_bUseFastEnc; }
func (this *TEncCfg)  getUseEarlyCU                   ()  bool   { return this.m_bUseEarlyCU; }
func (this *TEncCfg)  getUseFastDecisionForMerge      ()  bool   { return this.m_useFastDecisionForMerge; }
func (this *TEncCfg)  getUseCbfFastMode               ()  bool   { return this.m_bUseCbfFastMode; }
func (this *TEncCfg)  getUseEarlySkipDetection        ()  bool   { return this.m_useEarlySkipDetection; }
func (this *TEncCfg)  getUseConstrainedIntraPred      ()  bool   { return this.m_bUseConstrainedIntraPred; }
func (this *TEncCfg)  getPCMInputBitDepthFlag         ()  bool   { return this.m_bPCMInputBitDepthFlag;   }
func (this *TEncCfg)  getPCMFilterDisableFlag         ()  bool   { return this.m_bPCMFilterDisableFlag;   } 
func (this *TEncCfg)  getUsePCM                       ()  bool   { return this.m_usePCM;                 }
func (this *TEncCfg)  getPCMLog2MaxSize               ()  uint   { return this.m_pcmLog2MaxSize;  }
func (this *TEncCfg)  getPCMLog2MinSize               ()  uint   { return  this.m_uiPCMLog2MinSize;  }

func (this *TEncCfg)  getUseTransformSkip                  ()  bool    { return this.m_useTransformSkip;        }
func (this *TEncCfg)  setUseTransformSkip                  ( b bool ) { this.m_useTransformSkip  = b;       }
func (this *TEncCfg)  getUseTransformSkipFast              ()  bool    { return this.m_useTransformSkipFast;    }
func (this *TEncCfg)  setUseTransformSkipFast              ( b bool ) { this.m_useTransformSkipFast  = b;   }
func (this *TEncCfg)  getdQPs                         ()  []int    { return this.m_aidQP;       }
func (this *TEncCfg)  getDeltaQpRD                    ()  uint    { return this.m_uiDeltaQpRD; }

  //====== Slice ========
func (this *TEncCfg)  setSliceMode                   ( i int )       { this.m_iSliceMode = i;              }
func (this *TEncCfg)  setSliceArgument               ( i int )       { this.m_iSliceArgument = i;          }
func (this *TEncCfg)  getSliceMode                   ()      int        { return this.m_iSliceMode;           }
func (this *TEncCfg)  getSliceArgument               ()      int        { return this.m_iSliceArgument;       }
  //====== Dependent Slice ========
func (this *TEncCfg)  setDependentSliceMode            ( i int )      { this.m_iDependentSliceMode = i;       }
func (this *TEncCfg)  setDependentSliceArgument        ( i int )      { this.m_iDependentSliceArgument = i;   }
func (this *TEncCfg)  getDependentSliceMode            ()    int          { return this.m_iDependentSliceMode;    }
func (this *TEncCfg)  getDependentSliceArgument        ()    int          { return this.m_iDependentSliceArgument;}
//#if DEPENDENT_SLICES && !REMOVE_ENTROPY_SLICES
func (this *TEncCfg)  setEntropySliceEnabledFlag       ( b bool )     { this.m_entropySliceEnabledFlag = b;    }
func (this *TEncCfg)  getEntropySliceEnabledFlag       ()    bool          { return this.m_entropySliceEnabledFlag; }
//#endif
func (this *TEncCfg)  setLFCrossSliceBoundaryFlag     ( bValue bool  )    { this.m_bLFCrossSliceBoundaryFlag = bValue; }
func (this *TEncCfg)  getLFCrossSliceBoundaryFlag     ()   bool                 { return this.m_bLFCrossSliceBoundaryFlag;   }

func (this *TEncCfg)  setUseSAO                  (bVal bool)     {this.m_bUseSAO = bVal;}
func (this *TEncCfg)  getUseSAO                  ()    bool          {return this.m_bUseSAO;}
func (this *TEncCfg)  setMaxNumOffsetsPerPic                   (iVal int)            { this.m_maxNumOffsetsPerPic = iVal; }
func (this *TEncCfg)  getMaxNumOffsetsPerPic                   ()   int                 { return this.m_maxNumOffsetsPerPic; }
func (this *TEncCfg)  setSaoLcuBoundary              ( val bool)      { this.m_saoLcuBoundary = val; }
func (this *TEncCfg)  getSaoLcuBoundary              ()       bool       { return this.m_saoLcuBoundary; }
func (this *TEncCfg)  setSaoLcuBasedOptimization               (val bool)            { this.m_saoLcuBasedOptimization = val; }
func (this *TEncCfg)  getSaoLcuBasedOptimization               ()    bool                { return this.m_saoLcuBasedOptimization; }
func (this *TEncCfg)  setLFCrossTileBoundaryFlag               ( val bool  )       { this.m_loopFilterAcrossTilesEnabledFlag = val; }
func (this *TEncCfg)  getLFCrossTileBoundaryFlag               ()    bool                { return this.m_loopFilterAcrossTilesEnabledFlag;   }
func (this *TEncCfg)  setUniformSpacingIdr           ( i int )           { this.m_iUniformSpacingIdr = i; }
func (this *TEncCfg)  getUniformSpacingIdr           ()   int               { return this.m_iUniformSpacingIdr; }
func (this *TEncCfg)  setNumColumnsMinus1            ( i int )           { this.m_iNumColumnsMinus1 = i; }
func (this *TEncCfg)  getNumColumnsMinus1            ()   int               { return this.m_iNumColumnsMinus1; }
//#if MIN_SPATIAL_SEGMENTATION
func (this *TEncCfg)  setColumnWidth ( columnWidth []uint ) {
    if this.m_iUniformSpacingIdr == 0 && this.m_iNumColumnsMinus1 > 0 {
      var m_iWidthInCU int
      if this.m_iSourceWidth%int(TLibCommon.G_uiMaxCUWidth)!=0 {
      	m_iWidthInCU = this.m_iSourceWidth/int(TLibCommon.G_uiMaxCUWidth) + 1;
      }else{
      	m_iWidthInCU = this.m_iSourceWidth/int(TLibCommon.G_uiMaxCUWidth);
      }
      this.m_puiColumnWidth = make([]uint, this.m_iNumColumnsMinus1);

      for i :=0; i<this.m_iNumColumnsMinus1; i++ {
        this.m_puiColumnWidth[i] = columnWidth[i];
        fmt.Printf("col: this.m_iWidthInCU= %4d i=%4d width= %4d\n",m_iWidthInCU,i,this.m_puiColumnWidth[i]); //AFU
      }
    }
  }
/*#else
  Void  setColumnWidth ( Char* str )
  {
    Char *columnWidth;
    i int=0;
    Int  this.m_iWidthInCU = ( this.m_iSourceWidth%g_uiMaxCUWidth ) ? this.m_iSourceWidth/g_uiMaxCUWidth + 1 : this.m_iSourceWidth/g_uiMaxCUWidth;

    if( this.m_iUniformSpacingIdr == 0 && this.m_iNumColumnsMinus1 > 0 )
    {
      this.m_puiColumnWidth = new UInt[this.m_iNumColumnsMinus1];

      columnWidth = strtok(str, " ,-");
      while(columnWidth!=NULL)
      {
        if( i>=this.m_iNumColumnsMinus1 )
        {
          printf( "The number of columns whose width are defined is larger than the allowed number of columns.\n" );
          exit( EXIT_FAILURE );
        }
        *( this.m_puiColumnWidth + i ) = atoi( columnWidth );
        printf("col: this.m_iWidthInCU= %4d i=%4d width= %4d\n",this.m_iWidthInCU,i,this.m_puiColumnWidth[i]); //AFU
        columnWidth = strtok(NULL, " ,-");
        i++;
      }
      if( i<this.m_iNumColumnsMinus1 )
      {
        printf( "The width of some columns is not defined.\n" );
        exit( EXIT_FAILURE );
      }
    }
  }
#endif*/
func (this *TEncCfg)  getColumnWidth                 ( columnidx uint ) uint { return this.m_puiColumnWidth [columnidx]; }
func (this *TEncCfg)  setNumRowsMinus1               ( i int )           	   { this.m_iNumRowsMinus1 = i; }
func (this *TEncCfg)  getNumRowsMinus1               ()                 int  { return this.m_iNumRowsMinus1; }
//#if MIN_SPATIAL_SEGMENTATION
func (this *TEncCfg)  setRowHeight (rowHeight []uint) {
    if this.m_iUniformSpacingIdr == 0 && this.m_iNumRowsMinus1 > 0 {
      var m_iHeightInCU int
      if this.m_iSourceHeight%int(TLibCommon.G_uiMaxCUHeight)!=0 {
      	m_iHeightInCU = this.m_iSourceHeight/int(TLibCommon.G_uiMaxCUHeight) + 1;
      }else{
      	m_iHeightInCU = this.m_iSourceHeight/int(TLibCommon.G_uiMaxCUHeight);
      }
      this.m_puiRowHeight = make([]uint, this.m_iNumRowsMinus1 );

      for i:=0; i<this.m_iNumRowsMinus1; i++ {
        this.m_puiRowHeight[i] = rowHeight[i];
        fmt.Printf("row: this.m_iHeightInCU=%4d i=%4d height=%4d\n",m_iHeightInCU,i,this.m_puiRowHeight[i]); //AFU
      }
    }
  }
/*#else
  Void  setRowHeight (Char* str)
  {
    Char *rowHeight;
    i int=0;
    Int  this.m_iHeightInCU = ( this.m_iSourceHeight%g_uiMaxCUHeight ) ? this.m_iSourceHeight/g_uiMaxCUHeight + 1 : this.m_iSourceHeight/g_uiMaxCUHeight;

    if( this.m_iUniformSpacingIdr == 0 && this.m_iNumRowsMinus1 > 0 )
    {
      this.m_puiRowHeight = new UInt[this.m_iNumRowsMinus1];

      rowHeight = strtok(str, " ,-");
      while(rowHeight!=NULL)
      {
        if( i>=this.m_iNumRowsMinus1 )
        {
          printf( "The number of rows whose height are defined is larger than the allowed number of rows.\n" );
          exit( EXIT_FAILURE );
        }
        *( this.m_puiRowHeight + i ) = atoi( rowHeight );
        printf("row: this.m_iHeightInCU=%4d i=%4d height=%4d\n",this.m_iHeightInCU,i,this.m_puiRowHeight[i]); //AFU
        rowHeight = strtok(NULL, " ,-");
        i++;
      }
      if( i<this.m_iNumRowsMinus1 )
      {
        printf( "The height of some rows is not defined.\n" );
        exit( EXIT_FAILURE );
     }
    }
  }
#endif*/
func (this *TEncCfg)  getRowHeight                   ( rowIdx uint )  uint   { return this.m_puiRowHeight [ rowIdx ]; }
func (this *TEncCfg)  xCheckGSParameters() {
}
func (this *TEncCfg)  setWaveFrontSynchro(iWaveFrontSynchro int)       { this.m_iWaveFrontSynchro = iWaveFrontSynchro; }
func (this *TEncCfg)  getWaveFrontsynchro()                           int { return this.m_iWaveFrontSynchro; }
func (this *TEncCfg)  setWaveFrontSubstreams(iWaveFrontSubstreams int) { this.m_iWaveFrontSubstreams = iWaveFrontSubstreams; }
func (this *TEncCfg)  getWaveFrontSubstreams()                        int { return this.m_iWaveFrontSubstreams; }
func (this *TEncCfg)  setDecodedPictureHashSEIEnabled(b int)           { this.m_decodedPictureHashSEIEnabled = b; }
func (this *TEncCfg)  getDecodedPictureHashSEIEnabled()               int { return this.m_decodedPictureHashSEIEnabled; }
func (this *TEncCfg)  setBufferingPeriodSEIEnabled(b int)              { this.m_bufferingPeriodSEIEnabled = b; }
func (this *TEncCfg)  getBufferingPeriodSEIEnabled()                  int { return this.m_bufferingPeriodSEIEnabled; }
func (this *TEncCfg)  setPictureTimingSEIEnabled(b int)                { this.m_pictureTimingSEIEnabled = b; }
func (this *TEncCfg)  getPictureTimingSEIEnabled()                    int { return this.m_pictureTimingSEIEnabled; }
func (this *TEncCfg)  setRecoveryPointSEIEnabled(b int)                { this.m_recoveryPointSEIEnabled = b; }
func (this *TEncCfg)  getRecoveryPointSEIEnabled()                    int { return this.m_recoveryPointSEIEnabled; }
//#if SEI_DISPLAY_ORIENTATION
func (this *TEncCfg)  setDisplayOrientationSEIAngle(b int)             { this.m_displayOrientationSEIAngle = b; }
func (this *TEncCfg)  getDisplayOrientationSEIAngle()                 int { return this.m_displayOrientationSEIAngle; }
//#endif
//#if SEI_TEMPORAL_LEVEL0_INDEX
func (this *TEncCfg)  setTemporalLevel0IndexSEIEnabled(b int)          { this.m_temporalLevel0IndexSEIEnabled = b; }
func (this *TEncCfg)  getTemporalLevel0IndexSEIEnabled()              int { return this.m_temporalLevel0IndexSEIEnabled; }
//#endif
func (this *TEncCfg)  setUseWP               ( b bool )   { this.m_bUseWeightPred    = b;    }
func (this *TEncCfg)  setWPBiPred            ( b bool )    { this.m_useWeightedBiPred = b;    }
func (this *TEncCfg)  getUseWP               () bool           { return this.m_bUseWeightPred;    }
func (this *TEncCfg)  getWPBiPred            () bool           { return this.m_useWeightedBiPred; }
func (this *TEncCfg)  setLog2ParallelMergeLevelMinus2   ( u uint )    { this.m_log2ParallelMergeLevelMinus2       = u;    }
func (this *TEncCfg)  getLog2ParallelMergeLevelMinus2   () uint           { return this.m_log2ParallelMergeLevelMinus2;       }
func (this *TEncCfg)  setMaxNumMergeCand                ( u uint )    { this.m_maxNumMergeCand = u;      }
func (this *TEncCfg)  getMaxNumMergeCand                () uint           { return this.m_maxNumMergeCand;   }
func (this *TEncCfg)  setUseScalingListId    (  u int)    { this.m_useScalingListId       = u;   }
func (this *TEncCfg)  getUseScalingListId    ()  int          { return this.m_useScalingListId;      }
func (this *TEncCfg)  setScalingListFile     ( pch string){ this.m_scalingListFile     = pch; }
func (this *TEncCfg)  getScalingListFile     ()  string          { return this.m_scalingListFile;    }
func (this *TEncCfg)  setTMVPModeId ( u int) { this.m_TMVPModeId = u;    }
func (this *TEncCfg)  getTMVPModeId ()   int      { return this.m_TMVPModeId; }
func (this *TEncCfg)  setSignHideFlag( signHideFlag int) { this.m_signHideFlag = signHideFlag; }
func (this *TEncCfg)  getSignHideFlag()  int                  { return this.m_signHideFlag; }
//#if RATE_CONTROL_LAMBDA_DOMAIN
func (this *TEncCfg)  getUseRateCtrl         ()   bool           { return this.m_RCEnableRateControl;   }
func (this *TEncCfg)  setUseRateCtrl         ( b bool )      { this.m_RCEnableRateControl = b;      }
func (this *TEncCfg)  getTargetBitrate       ()   int           { return this.m_RCTargetBitrate;       }
func (this *TEncCfg)  setTargetBitrate       ( bitrate int ) { this.m_RCTargetBitrate  = bitrate;   }
func (this *TEncCfg)  getKeepHierBit         ()   bool           { return this.m_RCKeepHierarchicalBit; }
func (this *TEncCfg)  setKeepHierBit         ( b bool )      { this.m_RCKeepHierarchicalBit = b;    }
func (this *TEncCfg)  getLCULevelRC          ()   bool           { return this.m_RCLCULevelRC; }
func (this *TEncCfg)  setLCULevelRC          ( b bool )      { this.m_RCLCULevelRC = b; }
func (this *TEncCfg)  getUseLCUSeparateModel ()   bool           { return this.m_RCUseLCUSeparateModel; }
func (this *TEncCfg)  setUseLCUSeparateModel ( b bool )      { this.m_RCUseLCUSeparateModel = b;    }
func (this *TEncCfg)  getInitialQP           ()   int           { return this.m_RCInitialQP;           }
func (this *TEncCfg)  setInitialQP           ( QP int )      { this.m_RCInitialQP = QP;             }
func (this *TEncCfg)  getForceIntraQP        ()   bool           { return this.m_RCForceIntraQP;        }
func (this *TEncCfg)  setForceIntraQP        ( b bool )      { this.m_RCForceIntraQP = b;           }
/*#else
func (this *TEncCfg)  getUseRateCtrl    ()                { return this.m_enableRateCtrl;    }
func (this *TEncCfg)  setUseRateCtrl    (Bool flag)       { this.m_enableRateCtrl = flag;    }
func (this *TEncCfg)  getTargetBitrate  ()                { return this.m_targetBitrate;     }
func (this *TEncCfg)  setTargetBitrate  (Int target)      { this.m_targetBitrate  = target;  }
func (this *TEncCfg)  getNumLCUInUnit   ()                { return this.m_numLCUInUnit;      }
func (this *TEncCfg)  setNumLCUInUnit   (Int numLCUs)     { this.m_numLCUInUnit   = numLCUs; }
#endif*/
func (this *TEncCfg)  getTransquantBypassEnableFlag()   bool        { return this.m_TransquantBypassEnableFlag; }
func (this *TEncCfg)  setTransquantBypassEnableFlag( flag bool)  { this.m_TransquantBypassEnableFlag = flag; }
func (this *TEncCfg)  getCUTransquantBypassFlagValue()   bool       { return this.m_CUTransquantBypassFlagValue; }
func (this *TEncCfg)  setCUTransquantBypassFlagValue( flag bool) { this.m_CUTransquantBypassFlagValue = flag; }
func (this *TEncCfg)  setVPS(p *TLibCommon.TComVPS) { this.m_cVPS = p; }
func (this *TEncCfg)  getVPS() *TLibCommon.TComVPS{ return this.m_cVPS; }
func (this *TEncCfg)  setUseRecalculateQPAccordingToLambda ( b bool ) { this.m_recalculateQPAccordingToLambda = b;    }
func (this *TEncCfg)  getUseRecalculateQPAccordingToLambda ()  bool       { return this.m_recalculateQPAccordingToLambda; }

//#if STRONG_INTRA_SMOOTHING
func (this *TEncCfg)  setUseStrongIntraSmoothing ( b bool ) { this.m_useStrongIntraSmoothing = b;    }
func (this *TEncCfg)  getUseStrongIntraSmoothing ()  bool       { return this.m_useStrongIntraSmoothing; }
//#endif

func (this *TEncCfg)  setActiveParameterSetsSEIEnabled ( b int )  { this.m_activeParameterSetsSEIEnabled = b; }  
func (this *TEncCfg)  getActiveParameterSetsSEIEnabled ()  int       { return this.m_activeParameterSetsSEIEnabled; }
func (this *TEncCfg)  getVuiParametersPresentFlag()        bool         { return this.m_vuiParametersPresentFlag; }
func (this *TEncCfg)  setVuiParametersPresentFlag(i bool)           { this.m_vuiParametersPresentFlag = i; }
func (this *TEncCfg)  getAspectRatioInfoPresentFlag()      bool         { return this.m_aspectRatioInfoPresentFlag; }
func (this *TEncCfg)  setAspectRatioInfoPresentFlag(i bool)         { this.m_aspectRatioInfoPresentFlag = i; }
func (this *TEncCfg)  getAspectRatioIdc()                  int         { return this.m_aspectRatioIdc; }
func (this *TEncCfg)  setAspectRatioIdc(i int)                      { this.m_aspectRatioIdc = i; }
func (this *TEncCfg)  getSarWidth()                        int         { return this.m_sarWidth; }
func (this *TEncCfg)  setSarWidth(i int)                            { this.m_sarWidth = i; }
func (this *TEncCfg)  getSarHeight()                       int         { return this.m_sarHeight; }
func (this *TEncCfg)  setSarHeight(i int)                           { this.m_sarHeight = i; }
func (this *TEncCfg)  getOverscanInfoPresentFlag()         bool         { return this.m_overscanInfoPresentFlag; }
func (this *TEncCfg)  setOverscanInfoPresentFlag(i bool)            { this.m_overscanInfoPresentFlag = i; }
func (this *TEncCfg)  getOverscanAppropriateFlag()         bool         { return this.m_overscanAppropriateFlag; }
func (this *TEncCfg)  setOverscanAppropriateFlag(i bool)            { this.m_overscanAppropriateFlag = i; }
func (this *TEncCfg)  getVideoSignalTypePresentFlag()      bool         { return this.m_videoSignalTypePresentFlag; }
func (this *TEncCfg)  setVideoSignalTypePresentFlag(i bool)         { this.m_videoSignalTypePresentFlag = i; }
func (this *TEncCfg)  getVideoFormat()                     int         { return this.m_videoFormat; }
func (this *TEncCfg)  setVideoFormat(i int)                         { this.m_videoFormat = i; }
func (this *TEncCfg)  getVideoFullRangeFlag()              bool         { return this.m_videoFullRangeFlag; }
func (this *TEncCfg)  setVideoFullRangeFlag(i bool)                 { this.m_videoFullRangeFlag = i; }
func (this *TEncCfg)  getColourDescriptionPresentFlag()    bool         { return this.m_colourDescriptionPresentFlag; }
func (this *TEncCfg)  setColourDescriptionPresentFlag(i bool)       { this.m_colourDescriptionPresentFlag = i; }
func (this *TEncCfg)  getColourPrimaries()                 int         { return this.m_colourPrimaries; }
func (this *TEncCfg)  setColourPrimaries(i int)                     { this.m_colourPrimaries = i; }
func (this *TEncCfg)  getTransferCharacteristics()         int         { return this.m_transferCharacteristics; }
func (this *TEncCfg)  setTransferCharacteristics(i int)             { this.m_transferCharacteristics = i; }
func (this *TEncCfg)  getMatrixCoefficients()              int         { return this.m_matrixCoefficients; }
func (this *TEncCfg)  setMatrixCoefficients(i int)                  { this.m_matrixCoefficients = i; }
func (this *TEncCfg)  getChromaLocInfoPresentFlag()        bool         { return this.m_chromaLocInfoPresentFlag; }
func (this *TEncCfg)  setChromaLocInfoPresentFlag(i bool)           { this.m_chromaLocInfoPresentFlag = i; }
func (this *TEncCfg)  getChromaSampleLocTypeTopField()     int         { return this.m_chromaSampleLocTypeTopField; }
func (this *TEncCfg)  setChromaSampleLocTypeTopField(i int)         { this.m_chromaSampleLocTypeTopField = i; }
func (this *TEncCfg)  getChromaSampleLocTypeBottomField()  int        { return this.m_chromaSampleLocTypeBottomField; }
func (this *TEncCfg)  setChromaSampleLocTypeBottomField(i int)      { this.m_chromaSampleLocTypeBottomField = i; }
func (this *TEncCfg)  getNeutralChromaIndicationFlag()     bool         { return this.m_neutralChromaIndicationFlag; }
func (this *TEncCfg)  setNeutralChromaIndicationFlag(i bool)        { this.m_neutralChromaIndicationFlag = i; }
func (this *TEncCfg)  getBitstreamRestrictionFlag()        bool         { return this.m_bitstreamRestrictionFlag; }
func (this *TEncCfg)  setBitstreamRestrictionFlag(i bool)           { this.m_bitstreamRestrictionFlag = i; }
func (this *TEncCfg)  getTilesFixedStructureFlag()         bool         { return this.m_tilesFixedStructureFlag; }
func (this *TEncCfg)  setTilesFixedStructureFlag(i bool)            { this.m_tilesFixedStructureFlag = i; }
func (this *TEncCfg)  getMotionVectorsOverPicBoundariesFlag() bool      { return this.m_motionVectorsOverPicBoundariesFlag; }
func (this *TEncCfg)  setMotionVectorsOverPicBoundariesFlag(i bool) { this.m_motionVectorsOverPicBoundariesFlag = i; }
//#if MIN_SPATIAL_SEGMENTATION
func (this *TEncCfg)  getMinSpatialSegmentationIdc()       int         { return this.m_minSpatialSegmentationIdc; }
func (this *TEncCfg)  setMinSpatialSegmentationIdc(i int)           { this.m_minSpatialSegmentationIdc = i; }
//#endif
func (this *TEncCfg)  getMaxBytesPerPicDenom()             int         { return this.m_maxBytesPerPicDenom; }
func (this *TEncCfg)  setMaxBytesPerPicDenom(i int)                 { this.m_maxBytesPerPicDenom = i; }
func (this *TEncCfg)  getMaxBitsPerMinCuDenom()            int          { return this.m_maxBitsPerMinCuDenom; }
func (this *TEncCfg)  setMaxBitsPerMinCuDenom(i int)                { this.m_maxBitsPerMinCuDenom = i; }
func (this *TEncCfg)  getLog2MaxMvLengthHorizontal()       int         { return this.m_log2MaxMvLengthHorizontal; }
func (this *TEncCfg)  setLog2MaxMvLengthHorizontal(i int)           { this.m_log2MaxMvLengthHorizontal = i; }
func (this *TEncCfg)  getLog2MaxMvLengthVertical()         int         { return this.m_log2MaxMvLengthVertical; }
func (this *TEncCfg)  setLog2MaxMvLengthVertical(i int)             { this.m_log2MaxMvLengthVertical = i; }
