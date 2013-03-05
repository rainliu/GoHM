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
    "fmt"
    "gohm/TLibCommon"
)

type GOPEntry struct {
    m_POC      int
    m_QPOffset int
    m_QPFactor float64

    m_tcOffsetDiv2   int
    m_betaOffsetDiv2 int

    m_temporalId       int
    m_refPic           bool
    m_numRefPicsActive int
    m_sliceType        string
    m_numRefPics       int
    m_referencePics    [TLibCommon.MAX_NUM_REF_PICS]int
    m_usedByCurrPic    [TLibCommon.MAX_NUM_REF_PICS]bool
    //#if AUTO_INTER_RPS
    m_interRPSPrediction int
    /*#else
      Bool m_interRPSPrediction;
    #endif*/
    m_deltaRPS  int
    m_numRefIdc int
    m_refIdc    [TLibCommon.MAX_NUM_REF_PICS + 1]int
}

func NewGOPEntry() *GOPEntry {
    gop := &GOPEntry{m_POC: -1,
        m_QPOffset:           0,
        m_QPFactor:           0,
        m_tcOffsetDiv2:       0,
        m_betaOffsetDiv2:     0,
        m_temporalId:         0,
        m_refPic:             false,
        m_numRefPicsActive:   0,
        m_sliceType:          "P",
        m_numRefPics:         0,
        m_interRPSPrediction: 0,
        m_deltaRPS:           0,
        m_numRefIdc:          0}

    return gop
}

func (this *GOPEntry) GetPOC() int {
    return this.m_POC
}

func (this *GOPEntry) SetPOC(poc int) {
    this.m_POC = poc
}

func (this *GOPEntry) GetQPOffset() int {
    return this.m_QPOffset
}

func (this *GOPEntry) SetQPOffset(QPOffset int) {
    this.m_QPOffset = QPOffset
}

func (this *GOPEntry) GetQPFactor() float64 {
    return this.m_QPFactor
}

func (this *GOPEntry) SetQPFactor(QPFactor float64) {
    this.m_QPFactor = QPFactor
}

func (this *GOPEntry) GetBetaOffsetDiv2() int {
    return this.m_betaOffsetDiv2
}

func (this *GOPEntry) GetTcOffsetDiv2() int {
    return this.m_tcOffsetDiv2
}

func (this *GOPEntry) SetBetaOffsetDiv2(betaOffsetDiv2 int) {
    this.m_betaOffsetDiv2 = betaOffsetDiv2
}

func (this *GOPEntry) SetTcOffsetDiv2(tcOffsetDiv2 int) {
    this.m_tcOffsetDiv2 = tcOffsetDiv2
}

func (this *GOPEntry) GetNumRefPicsActive() int {
    return this.m_numRefPicsActive
}

func (this *GOPEntry) SetNumRefPicsActive(numRefPicsActive int) {
    this.m_numRefPicsActive = numRefPicsActive
}

func (this *GOPEntry) GetTemporalId() int {
    return this.m_temporalId
}
func (this *GOPEntry) SetTemporalId(temporalId int) {
    this.m_temporalId = temporalId
}

func (this *GOPEntry) SetNumRefIdc(numRefIdc int) {
    this.m_numRefIdc = numRefIdc
}
func (this *GOPEntry) GetNumRefIdc() int {
    return this.m_numRefIdc
}

func (this *GOPEntry) GetNumRefPics() int {
    return this.m_numRefPics
}
func (this *GOPEntry) SetNumRefPics(numRefPics int) {
    this.m_numRefPics = numRefPics
}

func (this *GOPEntry) GetReferencePics(i int) int {
    return this.m_referencePics[i]
}

func (this *GOPEntry) SetReferencePics(i int, value int) {
    this.m_referencePics[i] = value
}

func (this *GOPEntry) SetRefPic(refPic bool) {
    this.m_refPic = refPic
}

func (this *GOPEntry) GetRefPic() bool {
    return this.m_refPic
}

func (this *GOPEntry) SetUsedByCurrPic(i int, b bool) {
    this.m_usedByCurrPic[i] = b
}

func (this *GOPEntry) GetUsedByCurrPic(i int) bool {
    return this.m_usedByCurrPic[i]
}

func (this *GOPEntry) SetInterRPSPrediction(interRPSPrediction int) {
    this.m_interRPSPrediction = interRPSPrediction
}

func (this *GOPEntry) GetInterRPSPrediction() int {
    return this.m_interRPSPrediction
}

func (this *GOPEntry) SetRefIdc(i int, refIdc int) {
    this.m_refIdc[i] = refIdc
}

func (this *GOPEntry) GetRefIdc(i int) int {
    return this.m_refIdc[i]
}

func (this *GOPEntry) SetDeltaRPS(deltaRPS int) {
    this.m_deltaRPS = deltaRPS
}

func (this *GOPEntry) GetDeltaRPS() int {
    return this.m_deltaRPS
}

func (this *GOPEntry) GetSliceType() string {
    return this.m_sliceType
}

func (this *GOPEntry) SetSliceType(sliceType string) {
    this.m_sliceType = sliceType
}

// ====================================================================================================================
// Class definition
// ====================================================================================================================

/// encoder configuration class
type TEncCfg struct {
    //protected:
    //==== File I/O ========
    m_iFrameRate        int
    m_FrameSkip         uint
    m_iSourceWidth      int
    m_iSourceHeight     int
    m_conformanceMode   int
    m_conformanceWindow *TLibCommon.Window
    m_framesToBeEncoded int
    m_adLambdaModifier  [TLibCommon.MAX_TLAYER]float64

  // coding unit (CU) definition
        m_uiMaxCUWidth	uint;                                   ///< max. CU width in pixel
        m_uiMaxCUHeight	uint;                                  ///< max. CU height in pixel
        m_uiMaxCUDepth	uint;                                   ///< max. CU depth
        m_uiAddCUDepth	uint;
  
    /* profile & level */
    m_profile   TLibCommon.PROFILE
    m_levelTier TLibCommon.TIER
    m_level     TLibCommon.LEVEL

    //#if L0046_CONSTRAINT_FLAGS
    m_progressiveSourceFlag   bool
    m_interlacedSourceFlag    bool
    m_nonPackedConstraintFlag bool
    m_frameOnlyConstraintFlag bool
    //#endif

    //====== Coding Structure ========
    m_uiIntraPeriod         uint
    m_uiDecodingRefreshType uint ///< the type of decoding refresh employed for the random access.
    m_iGOPSize              int
    m_GOPList               [TLibCommon.MAX_GOP]*GOPEntry
    m_extraRPSs             int
    m_maxDecPicBuffering    [TLibCommon.MAX_TLAYER]int
    m_numReorderPics        [TLibCommon.MAX_TLAYER]int

    m_iQP int //  if (AdaptiveQP == OFF)

    m_aiPad [2]int

    m_iMaxRefPicNum int ///< this is used to mimic the sliding mechanism used by the decoder
    // TODO: We need to have a common sliding mechanism used by both the encoder and decoder

    m_maxTempLayer int ///< Max temporal layer
    m_useAMP       bool
    //======= Transform =============
    m_uiQuadtreeTULog2MaxSize   uint
    m_uiQuadtreeTULog2MinSize   uint
    m_uiQuadtreeTUMaxDepthInter uint
    m_uiQuadtreeTUMaxDepthIntra uint

    //====== Loop/Deblock Filter ========
    m_bLoopFilterDisable             bool
    m_loopFilterOffsetInPPS          bool
    m_loopFilterBetaOffsetDiv2       int
    m_loopFilterTcOffsetDiv2         int
    m_DeblockingFilterControlPresent bool
    m_bUseSAO                        bool
    m_maxNumOffsetsPerPic            int
    m_saoLcuBoundary                 bool
    m_saoLcuBasedOptimization        bool

    //====== Lossless ========
    m_useLossless bool
    //====== Motion search ========
    m_iFastSearch       int //  0:Full search  1:Diamond  2:PMVFAST
    m_iSearchRange      int //  0:Full frame
    m_bipredSearchRange int

    //====== Quality control ========
    m_iMaxDeltaQP    int //  Max. absolute delta QP (1:default)
    m_iMaxCuDQPDepth int //  Max. depth for a minimum CuDQP (0:default)

    m_chromaCbQpOffset int //  Chroma Cb QP Offset (0:default)
    m_chromaCrQpOffset int //  Chroma Cr Qp Offset (0:default)

    //#if ADAPTIVE_QP_SELECTION
    m_bUseAdaptQpSelect bool
    //#endif

    m_bUseAdaptiveQP     bool
    m_iQPAdaptationRange int

    //====== Tool list ========
    m_bUseSBACRD bool
    m_bUseASR    bool
    m_bUseHADME  bool
    m_bUseLComb  bool
    m_useRDOQ    bool
    m_useRDOQTS  bool
    //#if L0232_RD_PENALTY
    m_rdPenalty uint
    //#endif
    m_bUseFastEnc             bool
    m_bUseEarlyCU             bool
    m_useFastDecisionForMerge bool
    m_bUseCbfFastMode         bool
    m_useEarlySkipDetection   bool
    m_useTransformSkip        bool
    m_useTransformSkipFast    bool
    m_aidQP                   []int
    m_uiDeltaQpRD             uint

    m_bUseConstrainedIntraPred bool
    m_usePCM                   bool
    m_pcmLog2MaxSize           uint
    m_uiPCMLog2MinSize         uint
    //====== Slice ========
    m_sliceMode     int
    m_sliceArgument int
    //====== Dependent Slice ========
    m_sliceSegmentMode          int
    m_sliceSegmentArgument      int
    m_bLFCrossSliceBoundaryFlag bool

    m_bPCMInputBitDepthFlag            bool
    m_uiPCMBitDepthLuma                uint
    m_uiPCMBitDepthChroma              uint
    m_bPCMFilterDisableFlag            bool
    m_loopFilterAcrossTilesEnabledFlag bool
    m_iUniformSpacingIdr               int
    m_iNumColumnsMinus1                int
    m_puiColumnWidth                   []int
    m_iNumRowsMinus1                   int
    m_puiRowHeight                     []int

    m_iWaveFrontSynchro    int
    m_iWaveFrontSubstreams int

    m_decodedPictureHashSEIEnabled      int ///< Checksum(3)/CRC(2)/MD5(1)/disable(0) acting on decoded picture hash SEI message
    m_bufferingPeriodSEIEnabled         int
    m_pictureTimingSEIEnabled           int
    m_recoveryPointSEIEnabled           int
    m_framePackingSEIEnabled            int
    m_framePackingSEIType               int
    m_framePackingSEIId                 int
    m_framePackingSEIQuincunx           int
    m_framePackingSEIInterpretation     int
    m_displayOrientationSEIAngle        int
    m_temporalLevel0IndexSEIEnabled     int
    m_gradualDecodingRefreshInfoEnabled int
    m_decodingUnitInfoSEIEnabled        int
    //====== Weighted Prediction ========
    m_useWeightedPred              bool   //< Use of Weighting Prediction (P_SLICE)
    m_useWeightedBiPred            bool   //< Use of Bi-directional Weighting Prediction (B_SLICE)
    m_log2ParallelMergeLevelMinus2 uint   ///< Parallel merge estimation region
    m_maxNumMergeCand              uint   ///< Maximum number of merge candidates
    m_useScalingListId             int    ///< Using quantization matrix i.e. 0=off, 1=default, 2=file.
    m_scalingListFile              string ///< quantization matrix file name
    m_TMVPModeId                   int
    m_signHideFlag                 int
    //#if RATE_CONTROL_LAMBDA_DOMAIN
    m_RCEnableRateControl   bool
    m_RCTargetBitrate       int
    m_RCKeepHierarchicalBit bool
    m_RCLCULevelRC          bool
    m_RCUseLCUSeparateModel bool
    m_RCInitialQP           int
    m_RCForceIntraQP        bool
    /*#else
      Bool      m_enableRateCtrl;                                ///< Flag for using rate control algorithm
      Int       m_targetBitrate;                                 ///< target bitrate
      Int       m_numLCUInUnit;                                  ///< Total number of LCUs in a frame should be divided by the NumLCUInUnit
    #endif*/
    m_TransquantBypassEnableFlag         bool ///< transquant_bypass_enable_flag setting in PPS.
    m_CUTransquantBypassFlagValue        bool ///< if transquant_bypass_enable_flag, the fixed value to use for the per-CU cu_transquant_bypass_flag.
    m_cVPS                               *TLibCommon.TComVPS
    m_recalculateQPAccordingToLambda     bool               ///< recalculate QP value according to the lambda value
    m_activeParameterSetsSEIEnabled      int                ///< enable active parameter set SEI message
    m_vuiParametersPresentFlag           bool               ///< enable generation of VUI parameters
    m_aspectRatioInfoPresentFlag         bool               ///< Signals whether aspect_ratio_idc is present
    m_aspectRatioIdc                     int                ///< aspect_ratio_idc
    m_sarWidth                           int                ///< horizontal size of the sample aspect ratio
    m_sarHeight                          int                ///< vertical size of the sample aspect ratio
    m_overscanInfoPresentFlag            bool               ///< Signals whether overscan_appropriate_flag is present
    m_overscanAppropriateFlag            bool               ///< Indicates whether conformant decoded pictures are suitable for display using overscan
    m_videoSignalTypePresentFlag         bool               ///< Signals whether video_format, video_full_range_flag, and colour_description_present_flag are present
    m_videoFormat                        int                ///< Indicates representation of pictures
    m_videoFullRangeFlag                 bool               ///< Indicates the black level and range of luma and chroma signals
    m_colourDescriptionPresentFlag       bool               ///< Signals whether colour_primaries, transfer_characteristics and matrix_coefficients are present
    m_colourPrimaries                    int                ///< Indicates chromaticity coordinates of the source primaries
    m_transferCharacteristics            int                ///< Indicates the opto-electronic transfer characteristics of the source
    m_matrixCoefficients                 int                ///< Describes the matrix coefficients used in deriving luma and chroma from RGB primaries
    m_chromaLocInfoPresentFlag           bool               ///< Signals whether chroma_sample_loc_type_top_field and chroma_sample_loc_type_bottom_field are present
    m_chromaSampleLocTypeTopField        int                ///< Specifies the location of chroma samples for top field
    m_chromaSampleLocTypeBottomField     int                ///< Specifies the location of chroma samples for bottom field
    m_neutralChromaIndicationFlag        bool               ///< Indicates that the value of all decoded chroma samples is equal to 1<<(BitDepthCr-1)
    m_defaultDisplayWindow               *TLibCommon.Window ///< Represents the default display window parameters
    m_frameFieldInfoPresentFlag          bool               ///< Indicates that pic_struct and other field coding related values are present in picture timing SEI messages
    m_pocProportionalToTimingFlag        bool               ///< Indicates that the POC value is proportional to the output time w.r.t. first picture in CVS
    m_numTicksPocDiffOneMinus1           int                ///< Number of ticks minus 1 that for a POC difference of one
    m_bitstreamRestrictionFlag           bool               ///< Signals whether bitstream restriction parameters are present
    m_tilesFixedStructureFlag            bool               ///< Indicates that each active picture parameter set has the same values of the syntax elements related to tiles
    m_motionVectorsOverPicBoundariesFlag bool               ///< Indicates that no samples outside the picture boundaries are used for inter prediction
    m_minSpatialSegmentationIdc          int                ///< Indicates the maximum size of the spatial segments in the pictures in the coded video sequence
    m_maxBytesPerPicDenom                int                ///< Indicates a number of bytes not exceeded by the sum of the sizes of the VCL NAL units associated with any coded picture
    m_maxBitsPerMinCuDenom               int                ///< Indicates an upper bound for the number of bits of coding_unit() data
    m_log2MaxMvLengthHorizontal          int                ///< Indicate the maximum absolute value of a decoded horizontal MV component in quarter-pel luma units
    m_log2MaxMvLengthVertical            int                ///< Indicate the maximum absolute value of a decoded vertical MV component in quarter-pel luma units
    m_useStrongIntraSmoothing            bool               ///< enable the use of strong intra smoothing (bi_linear interpolation) for 32x32 blocks when reference samples are flat.
}

//public:
func NewTEncCfg() *TEncCfg {
    return &TEncCfg{m_conformanceWindow:TLibCommon.NewWindow(),
    				m_defaultDisplayWindow:TLibCommon.NewWindow()}
}

func (this *TEncCfg) SetProfile(profile TLibCommon.PROFILE) { this.m_profile = profile }
func (this *TEncCfg) SetLevel(tier TLibCommon.TIER, level TLibCommon.LEVEL) {
    this.m_levelTier = tier
    this.m_level = level
}

func (this *TEncCfg) SetFrameRate(i int)    { this.m_iFrameRate = i }
func (this *TEncCfg) SetFrameSkip(i uint)   { this.m_FrameSkip = i }
func (this *TEncCfg) SetSourceWidth(i int)  { this.m_iSourceWidth = i }
func (this *TEncCfg) SetSourceHeight(i int) { this.m_iSourceHeight = i }

func (this *TEncCfg) GetConformanceWindow() *TLibCommon.Window { return this.m_conformanceWindow }
func (this *TEncCfg) SetConformanceWindow(confLeft, confRight, confTop, confBottom int) {
    this.m_conformanceWindow.SetWindow(confLeft, confRight, confTop, confBottom)
}

func (this *TEncCfg) SetFramesToBeEncoded(i int) { this.m_framesToBeEncoded = i }

  // coding unit (CU) definition
func (this *TEncCfg) SetMaxCUWidth                   ( i uint)        { this.m_uiMaxCUWidth = i;}
func (this *TEncCfg) SetMaxCUHeight                  ( i uint)        { this.m_uiMaxCUHeight = i;}
func (this *TEncCfg) SetMaxCUDepth                   ( i uint)        { this.m_uiMaxCUDepth = i;}
func (this *TEncCfg) SetAddCUDepth                   ( i uint)        { this.m_uiAddCUDepth = i;}
func (this *TEncCfg) GetMaxCUWidth                   () uint       { return this.m_uiMaxCUWidth ;}
func (this *TEncCfg) GetMaxCUHeight                  () uint       { return this.m_uiMaxCUHeight;}
func (this *TEncCfg) GetMaxCUDepth                   () uint       { return this.m_uiMaxCUDepth ;}
func (this *TEncCfg) GetAddCUDepth                   () uint       { return this.m_uiAddCUDepth ;}
  
//====== Coding Structure ========
func (this *TEncCfg) SetIntraPeriod(i int) { this.m_uiIntraPeriod = uint(i) }
func (this *TEncCfg) SetDecodingRefreshType(i int) {
    this.m_uiDecodingRefreshType = uint(i)
}
func (this *TEncCfg) SetGOPSize(i int) { this.m_iGOPSize = i }
func (this *TEncCfg) SetGopList(GOPList []*GOPEntry) {
    for i := 0; i < TLibCommon.MAX_GOP; i++ {
        this.m_GOPList[i] = GOPList[i]
    }
}
func (this *TEncCfg) SetExtraRPSs(i int)          { this.m_extraRPSs = i }
func (this *TEncCfg) GetGOPEntry(i int) *GOPEntry { return this.m_GOPList[i] }
func (this *TEncCfg) SetMaxDecPicBuffering(u, tlayer uint) {
    this.m_maxDecPicBuffering[tlayer] = int(u)
}
func (this *TEncCfg) SetNumReorderPics(i int, tlayer uint) {
    this.m_numReorderPics[tlayer] = i
}

func (this *TEncCfg) SetQP(i int) { this.m_iQP = i }

func (this *TEncCfg) SetPad(iPad []int) {
    for i := 0; i < 2; i++ {
        this.m_aiPad[i] = iPad[i]
    }
}
func (this *TEncCfg) GetMaxRefPicNum() int { return this.m_iMaxRefPicNum }
func (this *TEncCfg) SetMaxRefPicNum(iMaxRefPicNum int) {
    this.m_iMaxRefPicNum = iMaxRefPicNum
}

func (this *TEncCfg) GetMaxTempLayer() int { return this.m_maxTempLayer }
func (this *TEncCfg) SetMaxTempLayer(maxTempLayer int) {
    this.m_maxTempLayer = maxTempLayer
}

//======== Transform =============
func (this *TEncCfg) SetQuadtreeTULog2MaxSize(u uint)   { this.m_uiQuadtreeTULog2MaxSize = u }
func (this *TEncCfg) SetQuadtreeTULog2MinSize(u uint)   { this.m_uiQuadtreeTULog2MinSize = u }
func (this *TEncCfg) SetQuadtreeTUMaxDepthInter(u uint) { this.m_uiQuadtreeTUMaxDepthInter = u }
func (this *TEncCfg) SetQuadtreeTUMaxDepthIntra(u uint) { this.m_uiQuadtreeTUMaxDepthIntra = u }

func (this *TEncCfg) SetUseAMP(b bool) { this.m_useAMP = b }

//====== Loop/Deblock Filter ========
func (this *TEncCfg) SetLoopFilterDisable(b bool)     { this.m_bLoopFilterDisable = b }
func (this *TEncCfg) SetLoopFilterOffsetInPPS(b bool) { this.m_loopFilterOffsetInPPS = b }
func (this *TEncCfg) SetLoopFilterBetaOffset(i int)   { this.m_loopFilterBetaOffsetDiv2 = i }
func (this *TEncCfg) SetLoopFilterTcOffset(i int)     { this.m_loopFilterTcOffsetDiv2 = i }
func (this *TEncCfg) SetDeblockingFilterControlPresent(b bool) {
    this.m_DeblockingFilterControlPresent = b
}

//====== Motion search ========
func (this *TEncCfg) SetFastSearch(i int)        { this.m_iFastSearch = i }
func (this *TEncCfg) SetSearchRange(i int)       { this.m_iSearchRange = i }
func (this *TEncCfg) SetBipredSearchRange(i int) { this.m_bipredSearchRange = i }

//====== Quality control ========
func (this *TEncCfg) SetMaxDeltaQP(i int)    { this.m_iMaxDeltaQP = i }
func (this *TEncCfg) SetMaxCuDQPDepth(i int) { this.m_iMaxCuDQPDepth = i }

func (this *TEncCfg) SetChromaCbQpOffset(i int) { this.m_chromaCbQpOffset = i }
func (this *TEncCfg) SetChromaCrQpOffset(i int) { this.m_chromaCrQpOffset = i }

//#if ADAPTIVE_QP_SELECTION
func (this *TEncCfg) SetUseAdaptQpSelect(i bool) { this.m_bUseAdaptQpSelect = i }
func (this *TEncCfg) GetUseAdaptQpSelect() bool  { return this.m_bUseAdaptQpSelect }

//#endif

func (this *TEncCfg) SetUseAdaptiveQP(b bool)    { this.m_bUseAdaptiveQP = b }
func (this *TEncCfg) SetQPAdaptationRange(i int) { this.m_iQPAdaptationRange = i }

//====== Lossless ========
func (this *TEncCfg) SetUseLossless(b bool) { this.m_useLossless = b }

//====== Sequence ========
func (this *TEncCfg) GetFrameRate() int         { return this.m_iFrameRate }
func (this *TEncCfg) GetFrameSkip() uint        { return this.m_FrameSkip }
func (this *TEncCfg) GetSourceWidth() int       { return this.m_iSourceWidth }
func (this *TEncCfg) GetSourceHeight() int      { return this.m_iSourceHeight }
func (this *TEncCfg) GetFramesToBeEncoded() int { return this.m_framesToBeEncoded }
func (this *TEncCfg) SetLambdaModifier(uiIndex uint, dValue float64) {
    this.m_adLambdaModifier[uiIndex] = dValue
}
func (this *TEncCfg) GetLambdaModifier(uiIndex uint) float64 {
    return this.m_adLambdaModifier[uiIndex]
}

//==== Coding Structure ========
func (this *TEncCfg) GetIntraPeriod() uint         { return this.m_uiIntraPeriod }
func (this *TEncCfg) GetDecodingRefreshType() uint { return this.m_uiDecodingRefreshType }
func (this *TEncCfg) GetGOPSize() int              { return this.m_iGOPSize }
func (this *TEncCfg) GetMaxDecPicBuffering(tlayer uint) int {
    return this.m_maxDecPicBuffering[tlayer]
}
func (this *TEncCfg) GetNumReorderPics(tlayer uint) int {
    return this.m_numReorderPics[tlayer]
}
func (this *TEncCfg) GetQP() int { return this.m_iQP }

func (this *TEncCfg) GetPad(i int) int { return this.m_aiPad[i] }

//======== Transform =============
func (this *TEncCfg) GetQuadtreeTULog2MaxSize() uint {
    return this.m_uiQuadtreeTULog2MaxSize
}
func (this *TEncCfg) GetQuadtreeTULog2MinSize() uint {
    return this.m_uiQuadtreeTULog2MinSize
}
func (this *TEncCfg) GetQuadtreeTUMaxDepthInter() uint {
    return this.m_uiQuadtreeTUMaxDepthInter
}
func (this *TEncCfg) GetQuadtreeTUMaxDepthIntra() uint {
    return this.m_uiQuadtreeTUMaxDepthIntra
}

//==== Loop/Deblock Filter ========
func (this *TEncCfg) GetLoopFilterDisable() bool { return this.m_bLoopFilterDisable }
func (this *TEncCfg) GetLoopFilterOffsetInPPS() bool {
    return this.m_loopFilterOffsetInPPS
}
func (this *TEncCfg) GetLoopFilterBetaOffset() int {
    return this.m_loopFilterBetaOffsetDiv2
}
func (this *TEncCfg) GetLoopFilterTcOffset() int {
    return this.m_loopFilterTcOffsetDiv2
}
func (this *TEncCfg) GetDeblockingFilterControlPresent() bool {
    return this.m_DeblockingFilterControlPresent
}

//==== Motion search ========
func (this *TEncCfg) GetFastSearch() int  { return this.m_iFastSearch }
func (this *TEncCfg) GetSearchRange() int { return this.m_iSearchRange }

//==== Quality control ========
func (this *TEncCfg) GetMaxDeltaQP() int        { return this.m_iMaxDeltaQP }
func (this *TEncCfg) GetMaxCuDQPDepth() int     { return this.m_iMaxCuDQPDepth }
func (this *TEncCfg) GetUseAdaptiveQP() bool    { return this.m_bUseAdaptiveQP }
func (this *TEncCfg) GetQPAdaptationRange() int { return this.m_iQPAdaptationRange }

//====== Lossless ========
func (this *TEncCfg) GetUseLossless() bool { return this.m_useLossless }

//==== Tool list ========
func (this *TEncCfg) SetUseSBACRD(b bool) { this.m_bUseSBACRD = b }
func (this *TEncCfg) SetUseASR(b bool)    { this.m_bUseASR = b }
func (this *TEncCfg) SetUseHADME(b bool)  { this.m_bUseHADME = b }
func (this *TEncCfg) SetUseLComb(b bool)  { this.m_bUseLComb = b }
func (this *TEncCfg) SetUseRDOQ(b bool)   { this.m_useRDOQ = b }
func (this *TEncCfg) SetUseRDOQTS(b bool) { this.m_useRDOQTS = b }

//#if L0232_RD_PENALTY
func (this *TEncCfg) SetRDpenalty(b uint) { this.m_rdPenalty = b }

//#endif
func (this *TEncCfg) SetUseFastEnc(b bool)              { this.m_bUseFastEnc = b }
func (this *TEncCfg) SetUseEarlyCU(b bool)              { this.m_bUseEarlyCU = b }
func (this *TEncCfg) SetUseFastDecisionForMerge(b bool) { this.m_useFastDecisionForMerge = b }
func (this *TEncCfg) SetUseCbfFastMode(b bool)          { this.m_bUseCbfFastMode = b }
func (this *TEncCfg) SetUseEarlySkipDetection(b bool)   { this.m_useEarlySkipDetection = b }
func (this *TEncCfg) SetUseConstrainedIntraPred(b bool) { this.m_bUseConstrainedIntraPred = b }
func (this *TEncCfg) SetPCMInputBitDepthFlag(b bool)    { this.m_bPCMInputBitDepthFlag = b }
func (this *TEncCfg) SetPCMFilterDisableFlag(b bool)    { this.m_bPCMFilterDisableFlag = b }
func (this *TEncCfg) SetUsePCM(b bool)                  { this.m_usePCM = b }
func (this *TEncCfg) SetPCMLog2MaxSize(u uint)          { this.m_pcmLog2MaxSize = u }
func (this *TEncCfg) SetPCMLog2MinSize(u uint)          { this.m_uiPCMLog2MinSize = u }
func (this *TEncCfg) SetdQPs(p []int)                   { this.m_aidQP = p }
func (this *TEncCfg) SetDeltaQpRD(u uint)               { this.m_uiDeltaQpRD = u }
func (this *TEncCfg) GetUseSBACRD() bool                { return this.m_bUseSBACRD }
func (this *TEncCfg) GetUseASR() bool                   { return this.m_bUseASR }
func (this *TEncCfg) GetUseHADME() bool                 { return this.m_bUseHADME }
func (this *TEncCfg) GetUseLComb() bool                 { return this.m_bUseLComb }
func (this *TEncCfg) GetUseRDOQ() bool                  { return this.m_useRDOQ }
func (this *TEncCfg) GetUseRDOQTS() bool                { return this.m_useRDOQTS }

//#if L0232_RD_PENALTY
func (this *TEncCfg) GetRDpenalty() uint { return this.m_rdPenalty }

//#endif
func (this *TEncCfg) GetUseFastEnc() bool              { return this.m_bUseFastEnc }
func (this *TEncCfg) GetUseEarlyCU() bool              { return this.m_bUseEarlyCU }
func (this *TEncCfg) GetUseFastDecisionForMerge() bool { return this.m_useFastDecisionForMerge }
func (this *TEncCfg) GetUseCbfFastMode() bool          { return this.m_bUseCbfFastMode }
func (this *TEncCfg) GetUseEarlySkipDetection() bool   { return this.m_useEarlySkipDetection }
func (this *TEncCfg) GetUseConstrainedIntraPred() bool { return this.m_bUseConstrainedIntraPred }
func (this *TEncCfg) GetPCMInputBitDepthFlag() bool    { return this.m_bPCMInputBitDepthFlag }
func (this *TEncCfg) GetPCMFilterDisableFlag() bool    { return this.m_bPCMFilterDisableFlag }
func (this *TEncCfg) GetUsePCM() bool                  { return this.m_usePCM }
func (this *TEncCfg) GetPCMLog2MaxSize() uint          { return this.m_pcmLog2MaxSize }
func (this *TEncCfg) GetPCMLog2MinSize() uint          { return this.m_uiPCMLog2MinSize }

func (this *TEncCfg) GetUseTransformSkip() bool  { return this.m_useTransformSkip }
func (this *TEncCfg) SetUseTransformSkip(b bool) { this.m_useTransformSkip = b }
func (this *TEncCfg) GetUseTransformSkipFast() bool {
    return this.m_useTransformSkipFast
}
func (this *TEncCfg) SetUseTransformSkipFast(b bool) { this.m_useTransformSkipFast = b }
func (this *TEncCfg) GetdQPs() []int                 { return this.m_aidQP }
func (this *TEncCfg) GetDeltaQpRD() uint             { return this.m_uiDeltaQpRD }

//====== Slice ========
func (this *TEncCfg) SetSliceMode(i int)     { this.m_sliceMode = i }
func (this *TEncCfg) SetSliceArgument(i int) { this.m_sliceArgument = i }
func (this *TEncCfg) GetSliceMode() int      { return this.m_sliceMode }
func (this *TEncCfg) GetSliceArgument() int  { return this.m_sliceArgument }

//====== Dependent Slice ========
func (this *TEncCfg) SetSliceSegmentMode(i int)     { this.m_sliceSegmentMode = i }
func (this *TEncCfg) SetSliceSegmentArgument(i int) { this.m_sliceSegmentArgument = i }
func (this *TEncCfg) GetSliceSegmentMode() int      { return this.m_sliceSegmentMode }
func (this *TEncCfg) GetSliceSegmentArgument() int  { return this.m_sliceSegmentArgument }

func (this *TEncCfg) SetLFCrossSliceBoundaryFlag(bValue bool) {
    this.m_bLFCrossSliceBoundaryFlag = bValue
}
func (this *TEncCfg) GetLFCrossSliceBoundaryFlag() bool {
    return this.m_bLFCrossSliceBoundaryFlag
}

func (this *TEncCfg) SetUseSAO(bVal bool) { this.m_bUseSAO = bVal }
func (this *TEncCfg) GetUseSAO() bool     { return this.m_bUseSAO }
func (this *TEncCfg) SetMaxNumOffsetsPerPic(iVal int) {
    this.m_maxNumOffsetsPerPic = iVal
}
func (this *TEncCfg) GetMaxNumOffsetsPerPic() int {
    return this.m_maxNumOffsetsPerPic
}
func (this *TEncCfg) SetSaoLcuBoundary(val bool) { this.m_saoLcuBoundary = val }
func (this *TEncCfg) GetSaoLcuBoundary() bool    { return this.m_saoLcuBoundary }
func (this *TEncCfg) SetSaoLcuBasedOptimization(val bool) {
    this.m_saoLcuBasedOptimization = val
}
func (this *TEncCfg) GetSaoLcuBasedOptimization() bool {
    return this.m_saoLcuBasedOptimization
}
func (this *TEncCfg) SetLFCrossTileBoundaryFlag(val bool) {
    this.m_loopFilterAcrossTilesEnabledFlag = val
}
func (this *TEncCfg) GetLFCrossTileBoundaryFlag() bool {
    return this.m_loopFilterAcrossTilesEnabledFlag
}
func (this *TEncCfg) SetUniformSpacingIdr(i int) { this.m_iUniformSpacingIdr = i }
func (this *TEncCfg) GetUniformSpacingIdr() int  { return this.m_iUniformSpacingIdr }
func (this *TEncCfg) SetNumColumnsMinus1(i int)  { this.m_iNumColumnsMinus1 = i }
func (this *TEncCfg) GetNumColumnsMinus1() int   { return this.m_iNumColumnsMinus1 }

func (this *TEncCfg) SetColumnWidth(columnWidth []int) {
    if this.m_iUniformSpacingIdr == 0 && this.m_iNumColumnsMinus1 > 0 {
        var m_iWidthInCU int
        if this.m_iSourceWidth%int(this.m_uiMaxCUWidth) != 0 {
            m_iWidthInCU = this.m_iSourceWidth/int(this.m_uiMaxCUWidth) + 1
        } else {
            m_iWidthInCU = this.m_iSourceWidth/int(this.m_uiMaxCUWidth)
        }
        this.m_puiColumnWidth = make([]int, this.m_iNumColumnsMinus1)

        for i := 0; i < this.m_iNumColumnsMinus1; i++ {
            this.m_puiColumnWidth[i] = columnWidth[i]
            fmt.Printf("col: this.m_iWidthInCU= %4d i=%4d width= %4d\n", m_iWidthInCU, i, this.m_puiColumnWidth[i]) //AFU
        }
    }
}
func (this *TEncCfg) GetColumnWidth(columnidx int) int {
    return this.m_puiColumnWidth[columnidx]
}
func (this *TEncCfg) SetNumRowsMinus1(i int) { this.m_iNumRowsMinus1 = i }
func (this *TEncCfg) GetNumRowsMinus1() int {
    return this.m_iNumRowsMinus1
}

func (this *TEncCfg) SetRowHeight(rowHeight []int) {
    if this.m_iUniformSpacingIdr == 0 && this.m_iNumRowsMinus1 > 0 {
        var m_iHeightInCU int
        if this.m_iSourceHeight%int(this.m_uiMaxCUHeight) != 0 {
            m_iHeightInCU = this.m_iSourceHeight/int(this.m_uiMaxCUHeight) + 1
        } else {
            m_iHeightInCU = this.m_iSourceHeight/int(this.m_uiMaxCUHeight)
        }
        this.m_puiRowHeight = make([]int, this.m_iNumRowsMinus1)

        for i := 0; i < this.m_iNumRowsMinus1; i++ {
            this.m_puiRowHeight[i] = rowHeight[i]
            fmt.Printf("row: this.m_iHeightInCU=%4d i=%4d height=%4d\n", m_iHeightInCU, i, this.m_puiRowHeight[i]) //AFU
        }
    }
}
func (this *TEncCfg) GetRowHeight(rowIdx int) int {
    return this.m_puiRowHeight[rowIdx]
}
func (this *TEncCfg) XCheckGSParameters() {
}
func (this *TEncCfg) SetWaveFrontSynchro(iWaveFrontSynchro int) {
    this.m_iWaveFrontSynchro = iWaveFrontSynchro
}
func (this *TEncCfg) GetWaveFrontsynchro() int {
    return this.m_iWaveFrontSynchro
}
func (this *TEncCfg) SetWaveFrontSubstreams(iWaveFrontSubstreams int) {
    this.m_iWaveFrontSubstreams = iWaveFrontSubstreams
}
func (this *TEncCfg) GetWaveFrontSubstreams() int {
    return this.m_iWaveFrontSubstreams
}
func (this *TEncCfg) SetDecodedPictureHashSEIEnabled(b int) { this.m_decodedPictureHashSEIEnabled = b }
func (this *TEncCfg) GetDecodedPictureHashSEIEnabled() int {
    return this.m_decodedPictureHashSEIEnabled
}
func (this *TEncCfg) SetBufferingPeriodSEIEnabled(b int) { this.m_bufferingPeriodSEIEnabled = b }
func (this *TEncCfg) GetBufferingPeriodSEIEnabled() int {
    return this.m_bufferingPeriodSEIEnabled
}
func (this *TEncCfg) SetPictureTimingSEIEnabled(b int) { this.m_pictureTimingSEIEnabled = b }
func (this *TEncCfg) GetPictureTimingSEIEnabled() int {
    return this.m_pictureTimingSEIEnabled
}
func (this *TEncCfg) SetRecoveryPointSEIEnabled(b int) { this.m_recoveryPointSEIEnabled = b }
func (this *TEncCfg) GetRecoveryPointSEIEnabled() int {
    return this.m_recoveryPointSEIEnabled
}
func (this *TEncCfg) SetFramePackingArrangementSEIEnabled(b int) { this.m_framePackingSEIEnabled = b }
func (this *TEncCfg) GetFramePackingArrangementSEIEnabled() int  { return this.m_framePackingSEIEnabled }
func (this *TEncCfg) SetFramePackingArrangementSEIType(b int)    { this.m_framePackingSEIType = b }
func (this *TEncCfg) GetFramePackingArrangementSEIType() int     { return this.m_framePackingSEIType }
func (this *TEncCfg) SetFramePackingArrangementSEIId(b int)      { this.m_framePackingSEIId = b }
func (this *TEncCfg) GetFramePackingArrangementSEIId() int       { return this.m_framePackingSEIId }
func (this *TEncCfg) SetFramePackingArrangementSEIQuincunx(b int) {
    this.m_framePackingSEIQuincunx = b
}
func (this *TEncCfg) GetFramePackingArrangementSEIQuincunx() int {
    return this.m_framePackingSEIQuincunx
}
func (this *TEncCfg) SetFramePackingArrangementSEIInterpretation(b int) {
    this.m_framePackingSEIInterpretation = b
}
func (this *TEncCfg) GetFramePackingArrangementSEIInterpretation() int {
    return this.m_framePackingSEIInterpretation
}
func (this *TEncCfg) SetDisplayOrientationSEIAngle(b int) { this.m_displayOrientationSEIAngle = b }
func (this *TEncCfg) GetDisplayOrientationSEIAngle() int {
    return this.m_displayOrientationSEIAngle
}
func (this *TEncCfg) SetTemporalLevel0IndexSEIEnabled(b int) {
    this.m_temporalLevel0IndexSEIEnabled = b
}
func (this *TEncCfg) GetTemporalLevel0IndexSEIEnabled() int {
    return this.m_temporalLevel0IndexSEIEnabled
}
func (this *TEncCfg) SetGradualDecodingRefreshInfoEnabled(b int) {
    this.m_gradualDecodingRefreshInfoEnabled = b
}
func (this *TEncCfg) GetGradualDecodingRefreshInfoEnabled() int {
    return this.m_gradualDecodingRefreshInfoEnabled
}
func (this *TEncCfg) SetDecodingUnitInfoSEIEnabled(b int) { this.m_decodingUnitInfoSEIEnabled = b }
func (this *TEncCfg) GetDecodingUnitInfoSEIEnabled() int {
    return this.m_decodingUnitInfoSEIEnabled
}
func (this *TEncCfg) SetUseWP(b bool)    { this.m_useWeightedPred = b }
func (this *TEncCfg) SetWPBiPred(b bool) { this.m_useWeightedBiPred = b }
func (this *TEncCfg) GetUseWP() bool     { return this.m_useWeightedPred }
func (this *TEncCfg) GetWPBiPred() bool  { return this.m_useWeightedBiPred }
func (this *TEncCfg) SetLog2ParallelMergeLevelMinus2(u uint) {
    this.m_log2ParallelMergeLevelMinus2 = u
}
func (this *TEncCfg) GetLog2ParallelMergeLevelMinus2() uint {
    return this.m_log2ParallelMergeLevelMinus2
}
func (this *TEncCfg) SetMaxNumMergeCand(u uint)        { this.m_maxNumMergeCand = u }
func (this *TEncCfg) GetMaxNumMergeCand() uint         { return this.m_maxNumMergeCand }
func (this *TEncCfg) SetUseScalingListId(u int)        { this.m_useScalingListId = u }
func (this *TEncCfg) GetUseScalingListId() int         { return this.m_useScalingListId }
func (this *TEncCfg) SetScalingListFile(pch string)    { this.m_scalingListFile = pch }
func (this *TEncCfg) GetScalingListFile() string       { return this.m_scalingListFile }
func (this *TEncCfg) SetTMVPModeId(u int)              { this.m_TMVPModeId = u }
func (this *TEncCfg) GetTMVPModeId() int               { return this.m_TMVPModeId }
func (this *TEncCfg) SetSignHideFlag(signHideFlag int) { this.m_signHideFlag = signHideFlag }
func (this *TEncCfg) GetSignHideFlag() int             { return this.m_signHideFlag }

//#if RATE_CONTROL_LAMBDA_DOMAIN
func (this *TEncCfg) GetUseRateCtrl() bool          { return this.m_RCEnableRateControl }
func (this *TEncCfg) SetUseRateCtrl(b bool)         { this.m_RCEnableRateControl = b }
func (this *TEncCfg) GetTargetBitrate() int         { return this.m_RCTargetBitrate }
func (this *TEncCfg) SetTargetBitrate(bitrate int)  { this.m_RCTargetBitrate = bitrate }
func (this *TEncCfg) GetKeepHierBit() bool          { return this.m_RCKeepHierarchicalBit }
func (this *TEncCfg) SetKeepHierBit(b bool)         { this.m_RCKeepHierarchicalBit = b }
func (this *TEncCfg) GetLCULevelRC() bool           { return this.m_RCLCULevelRC }
func (this *TEncCfg) SetLCULevelRC(b bool)          { this.m_RCLCULevelRC = b }
func (this *TEncCfg) GetUseLCUSeparateModel() bool  { return this.m_RCUseLCUSeparateModel }
func (this *TEncCfg) SetUseLCUSeparateModel(b bool) { this.m_RCUseLCUSeparateModel = b }
func (this *TEncCfg) GetInitialQP() int             { return this.m_RCInitialQP }
func (this *TEncCfg) SetInitialQP(QP int)           { this.m_RCInitialQP = QP }
func (this *TEncCfg) GetForceIntraQP() bool         { return this.m_RCForceIntraQP }
func (this *TEncCfg) SetForceIntraQP(b bool)        { this.m_RCForceIntraQP = b }

/*#else
func (this *TEncCfg)  GetUseRateCtrl    ()                { return this.m_enableRateCtrl;    }
func (this *TEncCfg)  SetUseRateCtrl    (Bool flag)       { this.m_enableRateCtrl = flag;    }
func (this *TEncCfg)  GetTargetBitrate  ()                { return this.m_targetBitrate;     }
func (this *TEncCfg)  SetTargetBitrate  (Int target)      { this.m_targetBitrate  = target;  }
func (this *TEncCfg)  GetNumLCUInUnit   ()                { return this.m_numLCUInUnit;      }
func (this *TEncCfg)  SetNumLCUInUnit   (Int numLCUs)     { this.m_numLCUInUnit   = numLCUs; }
#endif*/
func (this *TEncCfg) GetTransquantBypassEnableFlag() bool { return this.m_TransquantBypassEnableFlag }
func (this *TEncCfg) SetTransquantBypassEnableFlag(flag bool) {
    this.m_TransquantBypassEnableFlag = flag
}
func (this *TEncCfg) GetCUTransquantBypassFlagValue() bool {
    return this.m_CUTransquantBypassFlagValue
}
func (this *TEncCfg) SetCUTransquantBypassFlagValue(flag bool) {
    this.m_CUTransquantBypassFlagValue = flag
}
func (this *TEncCfg) SetVPS(p *TLibCommon.TComVPS) { this.m_cVPS = p }
func (this *TEncCfg) GetVPS() *TLibCommon.TComVPS  { return this.m_cVPS }
func (this *TEncCfg) SetUseRecalculateQPAccordingToLambda(b bool) {
    this.m_recalculateQPAccordingToLambda = b
}
func (this *TEncCfg) GetUseRecalculateQPAccordingToLambda() bool {
    return this.m_recalculateQPAccordingToLambda
}
func (this *TEncCfg) SetUseStrongIntraSmoothing(b bool) { this.m_useStrongIntraSmoothing = b }
func (this *TEncCfg) GetUseStrongIntraSmoothing() bool  { return this.m_useStrongIntraSmoothing }
func (this *TEncCfg) SetActiveParameterSetsSEIEnabled(b int) {
    this.m_activeParameterSetsSEIEnabled = b
}
func (this *TEncCfg) GetActiveParameterSetsSEIEnabled() int {
    return this.m_activeParameterSetsSEIEnabled
}
func (this *TEncCfg) GetVuiParametersPresentFlag() bool {
    return this.m_vuiParametersPresentFlag
}
func (this *TEncCfg) SetVuiParametersPresentFlag(i bool) { this.m_vuiParametersPresentFlag = i }
func (this *TEncCfg) GetAspectRatioInfoPresentFlag() bool {
    return this.m_aspectRatioInfoPresentFlag
}
func (this *TEncCfg) SetAspectRatioInfoPresentFlag(i bool) { this.m_aspectRatioInfoPresentFlag = i }
func (this *TEncCfg) GetAspectRatioIdc() int               { return this.m_aspectRatioIdc }
func (this *TEncCfg) SetAspectRatioIdc(i int)              { this.m_aspectRatioIdc = i }
func (this *TEncCfg) GetSarWidth() int                     { return this.m_sarWidth }
func (this *TEncCfg) SetSarWidth(i int)                    { this.m_sarWidth = i }
func (this *TEncCfg) GetSarHeight() int                    { return this.m_sarHeight }
func (this *TEncCfg) SetSarHeight(i int)                   { this.m_sarHeight = i }
func (this *TEncCfg) GetOverscanInfoPresentFlag() bool     { return this.m_overscanInfoPresentFlag }
func (this *TEncCfg) SetOverscanInfoPresentFlag(i bool)    { this.m_overscanInfoPresentFlag = i }
func (this *TEncCfg) GetOverscanAppropriateFlag() bool     { return this.m_overscanAppropriateFlag }
func (this *TEncCfg) SetOverscanAppropriateFlag(i bool)    { this.m_overscanAppropriateFlag = i }
func (this *TEncCfg) GetVideoSignalTypePresentFlag() bool {
    return this.m_videoSignalTypePresentFlag
}
func (this *TEncCfg) SetVideoSignalTypePresentFlag(i bool) { this.m_videoSignalTypePresentFlag = i }
func (this *TEncCfg) GetVideoFormat() int                  { return this.m_videoFormat }
func (this *TEncCfg) SetVideoFormat(i int)                 { this.m_videoFormat = i }
func (this *TEncCfg) GetVideoFullRangeFlag() bool          { return this.m_videoFullRangeFlag }
func (this *TEncCfg) SetVideoFullRangeFlag(i bool)         { this.m_videoFullRangeFlag = i }
func (this *TEncCfg) GetColourDescriptionPresentFlag() bool {
    return this.m_colourDescriptionPresentFlag
}
func (this *TEncCfg) SetColourDescriptionPresentFlag(i bool) { this.m_colourDescriptionPresentFlag = i }
func (this *TEncCfg) GetColourPrimaries() int                { return this.m_colourPrimaries }
func (this *TEncCfg) SetColourPrimaries(i int)               { this.m_colourPrimaries = i }
func (this *TEncCfg) GetTransferCharacteristics() int        { return this.m_transferCharacteristics }
func (this *TEncCfg) SetTransferCharacteristics(i int)       { this.m_transferCharacteristics = i }
func (this *TEncCfg) GetMatrixCoefficients() int             { return this.m_matrixCoefficients }
func (this *TEncCfg) SetMatrixCoefficients(i int)            { this.m_matrixCoefficients = i }
func (this *TEncCfg) GetChromaLocInfoPresentFlag() bool {
    return this.m_chromaLocInfoPresentFlag
}
func (this *TEncCfg) SetChromaLocInfoPresentFlag(i bool) { this.m_chromaLocInfoPresentFlag = i }
func (this *TEncCfg) GetChromaSampleLocTypeTopField() int {
    return this.m_chromaSampleLocTypeTopField
}
func (this *TEncCfg) SetChromaSampleLocTypeTopField(i int) { this.m_chromaSampleLocTypeTopField = i }
func (this *TEncCfg) GetChromaSampleLocTypeBottomField() int {
    return this.m_chromaSampleLocTypeBottomField
}
func (this *TEncCfg) SetChromaSampleLocTypeBottomField(i int) {
    this.m_chromaSampleLocTypeBottomField = i
}
func (this *TEncCfg) GetNeutralChromaIndicationFlag() bool {
    return this.m_neutralChromaIndicationFlag
}
func (this *TEncCfg) SetNeutralChromaIndicationFlag(i bool) { this.m_neutralChromaIndicationFlag = i }
func (this *TEncCfg) GetDefaultDisplayWindow() *TLibCommon.Window {
    return this.m_defaultDisplayWindow
}
func (this *TEncCfg) SetDefaultDisplayWindow(offsetLeft, offsetRight, offsetTop, offsetBottom int) {
    this.m_defaultDisplayWindow.SetWindow(offsetLeft, offsetRight, offsetTop, offsetBottom)
}
func (this *TEncCfg) GetFrameFieldInfoPresentFlag() bool    { return this.m_frameFieldInfoPresentFlag }
func (this *TEncCfg) SetFrameFieldInfoPresentFlag(i bool)   { this.m_frameFieldInfoPresentFlag = i }
func (this *TEncCfg) GetPocProportionalToTimingFlag() bool  { return this.m_pocProportionalToTimingFlag }
func (this *TEncCfg) SetPocProportionalToTimingFlag(x bool) { this.m_pocProportionalToTimingFlag = x }
func (this *TEncCfg) GetNumTicksPocDiffOneMinus1() int      { return this.m_numTicksPocDiffOneMinus1 }
func (this *TEncCfg) SetNumTicksPocDiffOneMinus1(x int)     { this.m_numTicksPocDiffOneMinus1 = x }
func (this *TEncCfg) GetBitstreamRestrictionFlag() bool {
    return this.m_bitstreamRestrictionFlag
}
func (this *TEncCfg) SetBitstreamRestrictionFlag(i bool) { this.m_bitstreamRestrictionFlag = i }
func (this *TEncCfg) GetTilesFixedStructureFlag() bool   { return this.m_tilesFixedStructureFlag }
func (this *TEncCfg) SetTilesFixedStructureFlag(i bool)  { this.m_tilesFixedStructureFlag = i }
func (this *TEncCfg) GetMotionVectorsOverPicBoundariesFlag() bool {
    return this.m_motionVectorsOverPicBoundariesFlag
}
func (this *TEncCfg) SetMotionVectorsOverPicBoundariesFlag(i bool) {
    this.m_motionVectorsOverPicBoundariesFlag = i
}
func (this *TEncCfg) GetMinSpatialSegmentationIdc() int {
    return this.m_minSpatialSegmentationIdc
}
func (this *TEncCfg) SetMinSpatialSegmentationIdc(i int) { this.m_minSpatialSegmentationIdc = i }
func (this *TEncCfg) GetMaxBytesPerPicDenom() int        { return this.m_maxBytesPerPicDenom }
func (this *TEncCfg) SetMaxBytesPerPicDenom(i int)       { this.m_maxBytesPerPicDenom = i }
func (this *TEncCfg) GetMaxBitsPerMinCuDenom() int       { return this.m_maxBitsPerMinCuDenom }
func (this *TEncCfg) SetMaxBitsPerMinCuDenom(i int)      { this.m_maxBitsPerMinCuDenom = i }
func (this *TEncCfg) GetLog2MaxMvLengthHorizontal() int {
    return this.m_log2MaxMvLengthHorizontal
}
func (this *TEncCfg) SetLog2MaxMvLengthHorizontal(i int) { this.m_log2MaxMvLengthHorizontal = i }
func (this *TEncCfg) GetLog2MaxMvLengthVertical() int    { return this.m_log2MaxMvLengthVertical }
func (this *TEncCfg) SetLog2MaxMvLengthVertical(i int)   { this.m_log2MaxMvLengthVertical = i }

//#if L0046_CONSTRAINT_FLAGS
func (this *TEncCfg) GetProgressiveSourceFlag() bool  { return this.m_progressiveSourceFlag }
func (this *TEncCfg) SetProgressiveSourceFlag(b bool) { this.m_progressiveSourceFlag = b }

func (this *TEncCfg) GetInterlacedSourceFlag() bool  { return this.m_interlacedSourceFlag }
func (this *TEncCfg) SetInterlacedSourceFlag(b bool) { this.m_interlacedSourceFlag = b }

func (this *TEncCfg) GetNonPackedConstraintFlag() bool  { return this.m_nonPackedConstraintFlag }
func (this *TEncCfg) SetNonPackedConstraintFlag(b bool) { this.m_nonPackedConstraintFlag = b }

func (this *TEncCfg) GetFrameOnlyConstraintFlag() bool  { return this.m_frameOnlyConstraintFlag }
func (this *TEncCfg) SetFrameOnlyConstraintFlag(b bool) { this.m_frameOnlyConstraintFlag = b }

//#endif
