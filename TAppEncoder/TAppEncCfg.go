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

package TAppEncoder

import (
    "bufio"
    "errors"
    "fmt"
    "gohm/TLibCommon"
    "gohm/TLibEncoder"
    "io"
    "log"
    "os"
    "strconv"
    "strings"
)

type Option interface {
    GetName() string
    Parse(arg string) error
    SetDefault()
}

/** Type specific option storage */
//template<typename T>
type OptionBase struct {
    opt_name          string
    opt_desc          string
    opt_storage       interface{}
    opt_default_value interface{}
}

func (this *OptionBase) GetName() string {
    return this.opt_name
}

type OptionString struct {
    OptionBase
}

func NewOptionString(name string, storage interface{}, default_value interface{}, desc string) *OptionString {
    return &OptionString{OptionBase{name, desc, storage, default_value}}
}

/* Generic parsing */
//template<typename T>
func (this *OptionString) Parse(arg string) error {
    //this.opt_storage = arg;
    var storage *string
    storage = this.opt_storage.(*string)

    *storage = arg

    return nil
}

func (this *OptionString) SetDefault() {
    var storage *string
    var default_value string
    storage = this.opt_storage.(*string)
    default_value = this.opt_default_value.(string)
    *storage = default_value
}

type OptionInt struct {
    OptionBase
}

func NewOptionInt(name string, storage interface{}, default_value interface{}, desc string) *OptionInt {
    return &OptionInt{OptionBase{name, desc, storage, default_value}}
}

/* Generic parsing */
//template<typename T>
func (this *OptionInt) Parse(arg string) error {
    argint, err := strconv.Atoi(arg)
    if err != nil {
        return err
    }

    //this.opt_storage = argint;
    var storage *int
    storage = this.opt_storage.(*int)

    *storage = argint

    return nil
}

func (this *OptionInt) SetDefault() {
    var storage *int
    var default_value int
    storage = this.opt_storage.(*int)
    default_value = this.opt_default_value.(int)
    *storage = default_value
}

type OptionUInt struct {
    OptionBase
}

func NewOptionUInt(name string, storage interface{}, default_value interface{}, desc string) *OptionUInt {
    return &OptionUInt{OptionBase{name, desc, storage, default_value}}
}

/* Generic parsing */
//template<typename T>
func (this *OptionUInt) Parse(arg string) error {
    argint, err := strconv.Atoi(arg)
    if err != nil {
        return err
    }

    //this.opt_storage = uint(argint);
    var storage *uint
    storage = this.opt_storage.(*uint)

    *storage = uint(argint)

    return nil
}

func (this *OptionUInt) SetDefault() {
    var storage *uint
    var default_value uint
    storage = this.opt_storage.(*uint)
    default_value = this.opt_default_value.(uint)
    *storage = default_value
}

type OptionBool struct {
    OptionBase
}

func NewOptionBool(name string, storage interface{}, default_value interface{}, desc string) *OptionBool {
    return &OptionBool{OptionBase{name, desc, storage, default_value}}
}

/* Generic parsing */
//template<typename T>
func (this *OptionBool) Parse(arg string) error {
    argint, err := strconv.Atoi(arg)
    if err != nil {
        return err
    }

    //this.opt_storage = argint!=0;
    var storage *bool
    storage = this.opt_storage.(*bool)

    *storage = argint != 0

    return nil
}

func (this *OptionBool) SetDefault() {
    var storage *bool
    var default_value bool
    storage = this.opt_storage.(*bool)
    default_value = this.opt_default_value.(bool)
    *storage = default_value
}

type OptionFloat struct {
    OptionBase
}

func NewOptionFloat(name string, storage interface{}, default_value interface{}, desc string) *OptionFloat {
    return &OptionFloat{OptionBase{name, desc, storage, default_value}}
}

/* Generic parsing */
//template<typename T>
func (this *OptionFloat) Parse(arg string) error {
    argfloat, err := strconv.ParseFloat(arg, 64)
    if err != nil {
        return err
    }

    //this.opt_storage = float64(argfloat);
    var storage *float64
    storage = this.opt_storage.(*float64)

    *storage = float64(argfloat)

    return nil
}

func (this *OptionFloat) SetDefault() {
    var storage *float64
    var default_value float64
    storage = this.opt_storage.(*float64)
    default_value = this.opt_default_value.(float64)
    *storage = default_value
}

type OptionGOPEntry struct {
    OptionBase
}

func NewOptionGOPEntry(name string, storage interface{}, default_value interface{}, desc string) *OptionGOPEntry {
    return &OptionGOPEntry{OptionBase{name, desc, storage, default_value}}
}

/* Generic parsing */
//template<typename T>
func (this *OptionGOPEntry) Parse(arg string) error {
    arglist := strings.Fields(arg)
    entry := TLibEncoder.NewGOPEntry()

    in := 0
    entry.SetSliceType(arglist[in])

    in++
    poc, _ := strconv.Atoi(arglist[in])
    entry.SetPOC(poc)

    in++
    qpoffset, _ := strconv.Atoi(arglist[in])
    entry.SetQPOffset(qpoffset)

    in++
    qpfactor, _ := strconv.ParseFloat(arglist[in], 64)
    entry.SetQPFactor(qpfactor)

    //#if VARYING_DBL_PARAMS
    in++
    tcOffsetDiv2, _ := strconv.Atoi(arglist[in])
    entry.SetTcOffsetDiv2(tcOffsetDiv2)

    in++
    betaOffsetDiv2, _ := strconv.Atoi(arglist[in])
    entry.SetBetaOffsetDiv2(betaOffsetDiv2)

    //#endif
    in++
    temporalId, _ := strconv.Atoi(arglist[in])
    entry.SetTemporalId(temporalId)

    in++
    numRefPicsActive, _ := strconv.Atoi(arglist[in])
    entry.SetNumRefPicsActive(numRefPicsActive)

    in++
    numRefPics, _ := strconv.Atoi(arglist[in])
    entry.SetNumRefPics(numRefPics)
    for i := 0; i < entry.GetNumRefPics(); i++ {
        in++
        referencePics, _ := strconv.Atoi(arglist[in])
        entry.SetReferencePics(i, referencePics)
    }

    in++
    interRPSPrediction, _ := strconv.Atoi(arglist[in])
    entry.SetInterRPSPrediction(interRPSPrediction)
    //#if AUTO_INTER_RPS
    if entry.GetInterRPSPrediction() == 1 {
        in++
        deltaRPS, _ := strconv.Atoi(arglist[in])
        entry.SetDeltaRPS(deltaRPS)

        in++
        numRefIdc, _ := strconv.Atoi(arglist[in])
        entry.SetNumRefIdc(numRefIdc)
        for i := 0; i < entry.GetNumRefIdc(); i++ {
            in++
            refIdc, _ := strconv.Atoi(arglist[in])
            entry.SetRefIdc(i, refIdc)
        }
    } else if entry.GetInterRPSPrediction() == 2 {
        in++
        deltaRPS, _ := strconv.Atoi(arglist[in])
        entry.SetDeltaRPS(deltaRPS)
    }
    /*#else
      if (entry.m_interRPSPrediction)
      {
        in>>entry.m_deltaRPS;
        in>>entry.m_numRefIdc;
        for ( Int i = 0; i < entry.m_numRefIdc; i++ )
        {
          in>>entry.m_refIdc[i];
        }
      }
    #endif*/

    //this.opt_storage = entry;
    var storage **TLibEncoder.GOPEntry
    storage = this.opt_storage.(**TLibEncoder.GOPEntry)

    *storage = entry

    return nil
}

func (this *OptionGOPEntry) SetDefault() {
    var storage **TLibEncoder.GOPEntry
    var default_value *TLibEncoder.GOPEntry
    storage = this.opt_storage.(**TLibEncoder.GOPEntry)
    default_value = this.opt_default_value.(*TLibEncoder.GOPEntry)
    *storage = default_value
}

type Options struct {
    opt_map map[string]Option //std::list<Names*>
}

func NewOptions() *Options {
    return &Options{opt_map: make(map[string]Option)}
}

func (opts *Options) AddOption(opt Option) {
    opts.opt_map[opt.GetName()] = opt
}

func (opts *Options) DoHelp(columns uint) {
}

/* for all options in opts, set their storage to their specified default value */
func (opts *Options) SetDefaults() {
    for name, _ := range opts.opt_map {
        //fmt.Printf("%s\n",name);
        opts.opt_map[name].SetDefault()
    }
}

func (opts *Options) StorePair(name string, value string) bool {
    //Options::NamesMap::iterator opt_it;
    opt, found := opts.opt_map[name]

    if !found {
        // not found
        fmt.Printf("Unknown option: '%s' (value:`%s')\n", name, value)
        return false
    }

    opt.Parse(value)

    return true
}

func (opts *Options) ScanLine(line string) {
    /* strip any leading whitespace */
    line = strings.TrimSpace(line)
    if line == "" {
        /* blank line */
        return
    }
    if line[0:1] == "#" {
        /* comment line */
        return
    }
    commentPos := strings.Index(line, "#")
    if commentPos != -1 {
        line = line[0:commentPos]
    }
    commaPos := strings.Index(line, ":")
    if commaPos == -1 {
        // error: badly formatted line
        return
    }
    name := strings.TrimSpace(line[0:commaPos])
    value := strings.TrimSpace(line[commaPos+1:])

    //fmt.Printf("%s : %s\n", name, value);

    /* store the value in option */
    opts.StorePair(name, value)
}
func (opts *Options) ScanFile(in io.Reader) (err error) {
    var line string
    reader := bufio.NewReader(in)
    eof := false

    line, err = reader.ReadString('\n')
    if err == io.EOF {
        err = nil
        eof = true
    } else if err != nil {
        return err
    }

    for !eof {
        opts.ScanLine(line)

        line, err = reader.ReadString('\n')
        if err == io.EOF {
            err = nil
            eof = true
        } else if err != nil {
            return err
        }
    }

    return nil
}

func (opts *Options) ParseConfigFile(filename string) {
    cfgstream, err := os.Open(filename)
    if err != nil {
        log.Fatal(err)
    }
    defer cfgstream.Close()

    opts.ScanFile(cfgstream)
}

// ====================================================================================================================
// Class definition
// ====================================================================================================================

/// encoder configuration class
type TAppEncCfg struct {
    //protected:
    // file I/O
    m_pchInputFile     string                         ///< source file name
    m_pchBitstreamFile string                         ///< output bitstream file
    m_pchReconFile     string                         ///< output reconstruction file
    m_pchTraceFile     string     					  ///< trace file name
    m_adLambdaModifier [TLibCommon.MAX_TLAYER]float64 ///< Lambda modifier array for each temporal layer
    // source specification
    m_iFrameRate        int  ///< source frame-rates (Hz)
    m_FrameSkip         uint ///< number of skipped frames from the beginning
    m_iSourceWidth      int  ///< source width in pixel
    m_iSourceHeight     int  ///< source height in pixel
    m_conformanceMode   int
    m_confLeft          int
    m_confRight         int
    m_confTop           int
    m_confBottom        int
    m_framesToBeEncoded int    ///< number of encoded frames
    m_aiPad             [2]int ///< number of padded pixels for width and height

    // profile/level
    m_profile   int //TLibCommon.PROFILE;
    m_levelTier int
    TLibCommon.TIER
    m_level int
    TLibCommon.LEVEL
    //#if L0046_CONSTRAINT_FLAGS
    m_progressiveSourceFlag   bool
    m_interlacedSourceFlag    bool
    m_nonPackedConstraintFlag bool
    m_frameOnlyConstraintFlag bool
    //#endif

    // coding structure
    m_iIntraPeriod         int                                       ///< period of I-slice (random access period)
    m_iDecodingRefreshType int                                       ///< random access type
    m_iGOPSize             int                                       ///< GOP size of hierarchical structure
    m_extraRPSs            int                                       ///< extra RPSs added to handle CRA
    m_GOPList              [TLibCommon.MAX_GOP]*TLibEncoder.GOPEntry ///< the coding structure entries from the config file
    m_numReorderPics       [TLibCommon.MAX_TLAYER]int                ///< total number of reorder pictures
    m_maxDecPicBuffering   [TLibCommon.MAX_TLAYER]int                ///< total number of reference pictures needed for decoding
    m_bUseLComb            bool                                      ///< flag for using combined reference list for uni-prediction in B-slices (JCTVC-D421)
    m_useTransformSkip     bool                                      ///< flag for enabling intra transform skipping
    m_useTransformSkipFast bool                                      ///< flag for enabling fast intra transform skipping
    m_enableAMP            bool
    // coding quality
    m_fQP            float64 ///< QP value of key-picture (floating point)
    m_iQP            int     ///< QP value of key-picture (integer)
    m_pchdQPFile     string  ///< QP offset for each slice (initialized from external file)
    m_aidQP          []int   ///< array of slice QP values
    m_iMaxDeltaQP    int     ///< max. |delta QP|
    m_uiDeltaQpRD    uint    ///< dQP range for multi-pass slice QP optimization
    m_iMaxCuDQPDepth int     ///< Max. depth for a minimum CuDQPSize (0:default)

    m_cbQpOffset int ///< Chroma Cb QP Offset (0:default)
    m_crQpOffset int ///< Chroma Cr QP Offset (0:default)

    //#if ADAPTIVE_QP_SELECTION
    m_bUseAdaptQpSelect bool
    //#endif

    m_bUseAdaptiveQP     bool ///< Flag for enabling QP adaptation based on a psycho-visual model
    m_iQPAdaptationRange int  ///< dQP range by QP adaptation

    m_maxTempLayer int ///< Max temporal layer

    // coding unit (CU) definition
    m_uiMaxCUWidth  uint ///< max. CU width in pixel
    m_uiMaxCUHeight uint ///< max. CU height in pixel
    m_uiMaxCUDepth  uint ///< max. CU depth
    m_uiAddCUDepth	uint; 

    // transfom unit (TU) definition
    m_uiQuadtreeTULog2MaxSize uint
    m_uiQuadtreeTULog2MinSize uint

    m_uiQuadtreeTUMaxDepthInter uint
    m_uiQuadtreeTUMaxDepthIntra uint

    // coding tools (bit-depth)
    m_inputBitDepthY    int ///< bit-depth of input file (luma component)
    m_inputBitDepthC    int ///< bit-depth of input file (chroma component)
    m_outputBitDepthY   int ///< bit-depth of output file (luma component)
    m_outputBitDepthC   int ///< bit-depth of output file (chroma component)
    m_internalBitDepthY int ///< bit-depth codec operates at in luma (input/output files will be converted)
    m_internalBitDepthC int ///< bit-depth codec operates at in chroma (input/output files will be converted)

    // coding tools (PCM bit-depth)
    m_bPCMInputBitDepthFlag bool ///< 0: PCM bit-depth is internal bit-depth. 1: PCM bit-depth is input bit-depth.

    // coding tool (lossless)
    m_useLossless             bool ///< flag for using lossless coding
    m_bUseSAO                 bool
    m_maxNumOffsetsPerPic     int  ///< SAO maximun number of offset per picture
    m_saoLcuBoundary          bool ///< SAO parameter estimation using non-deblocked pixels for LCU bottom and right boundary areas
    m_saoLcuBasedOptimization bool ///< SAO LCU-based optimization
    // coding tools (loop filter)
    m_bLoopFilterDisable             bool ///< flag for using deblocking filter
    m_loopFilterOffsetInPPS          bool ///< offset for deblocking filter in 0 = slice header, 1 = PPS
    m_loopFilterBetaOffsetDiv2       int  ///< beta offset for deblocking filter
    m_loopFilterTcOffsetDiv2         int  ///< tc offset for deblocking filter
    m_DeblockingFilterControlPresent bool ///< deblocking filter control present flag in PPS

    // coding tools (PCM)
    m_usePCM                bool ///< flag for using IPCM
    m_pcmLog2MaxSize        uint ///< log2 of maximum PCM block size
    m_uiPCMLog2MinSize      uint ///< log2 of minimum PCM block size
    m_bPCMFilterDisableFlag bool ///< PCM filter disable flag

    // coding tools (encoder-only parameters)
    m_bUseSBACRD bool ///< flag for using RD optimization based on SBAC
    m_bUseASR    bool ///< flag for using adaptive motion search range
    m_bUseHADME  bool ///< flag for using HAD in sub-pel ME
    m_useRDOQ    bool ///< flag for using RD optimized quantization
    m_useRDOQTS  bool ///< flag for using RD optimized quantization for transform skip
    //#if L0232_RD_PENALTY
    m_rdPenalty uint ///< RD-penalty for 32x32 TU for intra in non-intra slices (0: no RD-penalty, 1: RD-penalty, 2: maximum RD-penalty)
    //#endif
    m_iFastSearch             int  ///< ME mode, 0 = full, 1 = diamond, 2 = PMVFAST
    m_iSearchRange            int  ///< ME search range
    m_bipredSearchRange       int  ///< ME search range for bipred refinement
    m_bUseFastEnc             bool ///< flag for using fast encoder setting
    m_bUseEarlyCU             bool ///< flag for using Early CU setting
    m_useFastDecisionForMerge bool ///< flag for using Fast Decision Merge RD-Cost
    m_bUseCbfFastMode         bool ///< flag for using Cbf Fast PU Mode Decision
    m_useEarlySkipDetection   bool ///< flag for using Early SKIP Detection
    m_sliceMode               int  ///< 0: no slice limits, 1 : max number of CTBs per slice, 2: max number of bytes per slice,
    ///< 3: max number of tiles per slice
    m_sliceArgument    int ///< argument according to selected slice mode
    m_sliceSegmentMode int ///< 0: no slice segment limits, 1 : max number of CTBs per slice segment, 2: max number of bytes per slice segment,
    ///< 3: max number of tiles per slice segment
    m_sliceSegmentArgument int ///< argument according to selected slice segment mode

    m_bLFCrossSliceBoundaryFlag bool ///< 1: filter across slice boundaries 0: do not filter across slice boundaries
    m_bLFCrossTileBoundaryFlag  bool ///< 1: filter across tile boundaries  0: do not filter across tile boundaries
    m_iUniformSpacingIdr        int
    m_iNumColumnsMinus1         int
    m_pchColumnWidth            []byte
    m_iNumRowsMinus1            int
    m_pchRowHeight              []byte
    m_pColumnWidth              []int
    m_pRowHeight                []int
    m_iWaveFrontSynchro         int //< 0: no WPP. >= 1: WPP is enabled, the "Top right" from which inheritance occurs is this LCU offset in the line above the current.
    m_iWaveFrontSubstreams      int //< If iWaveFrontSynchro, this is the number of substreams per frame (dependent tiles) or per tile (independent tiles).

    m_bUseConstrainedIntraPred bool ///< flag for using constrained intra prediction

    m_decodedPictureHashSEIEnabled      int ///< Checksum(3)/CRC(2)/MD5(1)/disable(0) acting on decoded picture hash SEI message
    m_recoveryPointSEIEnabled           int
    m_bufferingPeriodSEIEnabled         int
    m_pictureTimingSEIEnabled           int
    m_framePackingSEIEnabled            int
    m_framePackingSEIType               int
    m_framePackingSEIId                 int
    m_framePackingSEIQuincunx           int
    m_framePackingSEIInterpretation     int
    m_displayOrientationSEIAngle        int
    m_temporalLevel0IndexSEIEnabled     int
    m_gradualDecodingRefreshInfoEnabled int
    m_decodingUnitInfoSEIEnabled        int

    // weighted prediction
    m_useWeightedPred   bool ///< Use of explicit Weighting Prediction for P_SLICE
    m_useWeightedBiPred bool ///< Use of Bi-Directional Weighting Prediction (B_SLICE)

    m_log2ParallelMergeLevel uint ///< Parallel merge estimation region
    m_maxNumMergeCand        uint ///< Max number of merge candidates

    m_TMVPModeId   int
    m_signHideFlag int
    //#if RATE_CONTROL_LAMBDA_DOMAIN
    m_RCEnableRateControl   bool ///< enable rate control or not
    m_RCTargetBitrate       int  ///< target bitrate when rate control is enabled
    m_RCKeepHierarchicalBit bool ///< whether keeping hierarchical bit allocation structure or not
    m_RCLCULevelRC          bool ///< true: LCU level rate control; false: picture level rate control
    m_RCUseLCUSeparateModel bool ///< use separate R-lambda model at LCU level
    m_RCInitialQP           int  ///< inital QP for rate control
    m_RCForceIntraQP        bool ///< force all intra picture to use initial QP or not
    /*#else
      Bool      m_enableRateCtrl;                                   ///< Flag for using rate control algorithm
      Int       m_targetBitrate;                                 ///< target bitrate
      Int       m_numLCUInUnit;                                  ///< Total number of LCUs in a frame should be completely divided by the NumLCUInUnit
    #endif*/
    m_useScalingListId int    ///< using quantization matrix
    m_scalingListFile  string ///< quantization matrix file name

    m_TransquantBypassEnableFlag  bool ///< transquant_bypass_enable_flag setting in PPS.
    m_CUTransquantBypassFlagValue bool ///< if transquant_bypass_enable_flag, the fixed value to use for the per-CU cu_transquant_bypass_flag.

    m_recalculateQPAccordingToLambda bool ///< recalculate QP value according to the lambda value
    m_useStrongIntraSmoothing        bool ///< enable strong intra smoothing for 32x32 blocks where the reference samples are flat
    m_activeParameterSetsSEIEnabled  int

    m_vuiParametersPresentFlag           bool ///< enable generation of VUI parameters
    m_aspectRatioInfoPresentFlag         bool ///< Signals whether aspect_ratio_idc is present
    m_aspectRatioIdc                     int  ///< aspect_ratio_idc
    m_sarWidth                           int  ///< horizontal size of the sample aspect ratio
    m_sarHeight                          int  ///< vertical size of the sample aspect ratio
    m_overscanInfoPresentFlag            bool ///< Signals whether overscan_appropriate_flag is present
    m_overscanAppropriateFlag            bool ///< Indicates whether cropped decoded pictures are suitable for display using overscan
    m_videoSignalTypePresentFlag         bool ///< Signals whether video_format, video_full_range_flag, and colour_description_present_flag are present
    m_videoFormat                        int  ///< Indicates representation of pictures
    m_videoFullRangeFlag                 bool ///< Indicates the black level and range of luma and chroma signals
    m_colourDescriptionPresentFlag       bool ///< Signals whether colour_primaries, transfer_characteristics and matrix_coefficients are present
    m_colourPrimaries                    int  ///< Indicates chromaticity coordinates of the source primaries
    m_transferCharacteristics            int  ///< Indicates the opto-electronic transfer characteristics of the source
    m_matrixCoefficients                 int  ///< Describes the matrix coefficients used in deriving luma and chroma from RGB primaries
    m_chromaLocInfoPresentFlag           bool ///< Signals whether chroma_sample_loc_type_top_field and chroma_sample_loc_type_bottom_field are present
    m_chromaSampleLocTypeTopField        int  ///< Specifies the location of chroma samples for top field
    m_chromaSampleLocTypeBottomField     int  ///< Specifies the location of chroma samples for bottom field
    m_neutralChromaIndicationFlag        bool ///< Indicates that the value of all decoded chroma samples is equal to 1<<(BitDepthCr-1)
    m_defaultDisplayWindowFlag           bool ///< Indicates the presence of the default window parameters
    m_defDispWinLeftOffset               int  ///< Specifies the left offset from the conformance window of the default window
    m_defDispWinRightOffset              int  ///< Specifies the right offset from the conformance window of the default window
    m_defDispWinTopOffset                int  ///< Specifies the top offset from the conformance window of the default window
    m_defDispWinBottomOffset             int  ///< Specifies the bottom offset from the conformance window of the default window
    m_frameFieldInfoPresentFlag          bool ///< Indicates that pic_struct values are present in picture timing SEI messages
    m_pocProportionalToTimingFlag        bool ///< Indicates that the POC value is proportional to the output time w.r.t. first picture in CVS
    m_numTicksPocDiffOneMinus1           int  ///< Number of ticks minus 1 that for a POC difference of one
    m_bitstreamRestrictionFlag           bool ///< Signals whether bitstream restriction parameters are present
    m_tilesFixedStructureFlag            bool ///< Indicates that each active picture parameter set has the same values of the syntax elements related to tiles
    m_motionVectorsOverPicBoundariesFlag bool ///< Indicates that no samples outside the picture boundaries are used for inter prediction
    m_minSpatialSegmentationIdc          int  ///< Indicates the maximum size of the spatial segments in the pictures in the coded video sequence
    m_maxBytesPerPicDenom                int  ///< Indicates a number of bytes not exceeded by the sum of the sizes of the VCL NAL units associated with any coded picture
    m_maxBitsPerMinCuDenom               int  ///< Indicates an upper bound for the number of bits of codinTLibCommon.G_unit() data
    m_log2MaxMvLengthHorizontal          int  ///< Indicate the maximum absolute value of a decoded horizontal MV component in quarter-pel luma units
    m_log2MaxMvLengthVertical            int  ///< Indicate the maximum absolute value of a decoded vertical MV component in quarter-pel luma units
}

func NewTAppEncCfg() *TAppEncCfg {
    return &TAppEncCfg{}
}

func (this *TAppEncCfg) Create() { ///< create option handling class
    //do nothing
}
func (this *TAppEncCfg) Destroy() { ///< destroy option handling class
    //do nothing
}
func (this *TAppEncCfg) ParseCfg(argc int, argv []string) error { ///< parse configuration file to fill member variables
    //do_help := false;

    var cfg_InputFile string
    var cfg_BitstreamFile string
    var cfg_ReconFile string
    var cfg_dQPFile string
    var cfg_ColumnWidth string
    var cfg_RowHeight string
    var cfg_ScalingListFile string

    opts := NewOptions()

    //("help", do_help, false, "this help text")
    //("c", po::parseConfigFile, "configuration file name")

    // File, I/O and source parameters
    opts.AddOption(NewOptionString("InputFile", &cfg_InputFile, string(""), "Original YUV input file name"))
    opts.AddOption(NewOptionString("BitstreamFile", &cfg_BitstreamFile, string(""), "Bitstream output file name"))
    opts.AddOption(NewOptionString("ReconFile", &cfg_ReconFile, string(""), "Reconstructed YUV output file name"))
    opts.AddOption(NewOptionInt("SourceWidth", &this.m_iSourceWidth, 0, "Source picture width"))
    opts.AddOption(NewOptionInt("SourceHeight", &this.m_iSourceHeight, 0, "Source picture height"))
    opts.AddOption(NewOptionInt("InputBitDepth", &this.m_inputBitDepthY, 8, "Bit-depth of input file"))
    opts.AddOption(NewOptionInt("OutputBitDepth", &this.m_outputBitDepthY, 0, "Bit-depth of output file (default:InternalBitDepth)"))
    opts.AddOption(NewOptionInt("InternalBitDepth", &this.m_internalBitDepthY, 0, "Bit-depth the codec operates at. (default:InputBitDepth) If different to InputBitDepth, source data will be converted"))
    opts.AddOption(NewOptionInt("InputBitDepthC", &this.m_inputBitDepthC, 0, "As per InputBitDepth but for chroma component. (default:InputBitDepth)"))
    opts.AddOption(NewOptionInt("OutputBitDepthC", &this.m_outputBitDepthC, 0, "As per OutputBitDepth but for chroma component. (default:InternalBitDepthC)"))
    opts.AddOption(NewOptionInt("InternalBitDepthC", &this.m_internalBitDepthC, 0, "As per InternalBitDepth but for chroma component. (default:IntrenalBitDepth)"))
    opts.AddOption(NewOptionInt("ConformanceMode", &this.m_conformanceMode, 0, "Window conformance mode (0: no window, 1:automatic padding, 2:padding, 3:conformance"))
    opts.AddOption(NewOptionInt("HorizontalPadding", &this.m_aiPad[0], 0, "Horizontal source padding for conformance window mode 2"))
    opts.AddOption(NewOptionInt("VerticalPadding", &this.m_aiPad[1], 0, "Vertical source padding for conformance window mode 2"))
    opts.AddOption(NewOptionInt("ConfLeft", &this.m_confLeft, 0, "Left cropping for conformance window mode 3"))
    opts.AddOption(NewOptionInt("ConfRight", &this.m_confRight, 0, "Right cropping for conformance window mode 3"))
    opts.AddOption(NewOptionInt("ConfTop", &this.m_confTop, 0, "Top cropping for conformance window mode 3"))
    opts.AddOption(NewOptionInt("ConfBottom", &this.m_confBottom, 0, "Bottom cropping for conformance window mode 3"))
    opts.AddOption(NewOptionInt("FrameRate", &this.m_iFrameRate, 0, "Frame rate"))
    opts.AddOption(NewOptionUInt("FrameSkip", &this.m_FrameSkip, uint(0), "Number of frames to skip at start of input YUV"))
    opts.AddOption(NewOptionInt("FramesToBeEncoded", &this.m_framesToBeEncoded, 0, "Number of frames to be encoded (default=all)"))

    // Profile and level
    opts.AddOption(NewOptionInt("Profile", &this.m_profile, 0, "Profile to be used when encoding (Incomplete)"))
    opts.AddOption(NewOptionInt("Level", &this.m_level, 0, "Level limit to be used, eg 5.1 (Incomplete)"))
    opts.AddOption(NewOptionInt("Tier", &this.m_levelTier, 0, "Tier to use for interpretation of --Level"))

    //#if L0046_CONSTRAINT_FLAGS
    opts.AddOption(NewOptionBool("ProgressiveSource", &this.m_progressiveSourceFlag, false, "Indicate that source is progressive"))
    opts.AddOption(NewOptionBool("InterlacedSource", &this.m_interlacedSourceFlag, false, "Indicate that source is interlaced"))
    opts.AddOption(NewOptionBool("NonPackedSource", &this.m_nonPackedConstraintFlag, false, "Indicate that source does not contain frame packing"))
    opts.AddOption(NewOptionBool("FrameOnly", &this.m_frameOnlyConstraintFlag, false, "Indicate that the bitstream contains only frames"))
    //#endif

    // Unit definition parameters
    opts.AddOption(NewOptionUInt("MaxCUWidth", &this.m_uiMaxCUWidth, uint(64), ""))
    opts.AddOption(NewOptionUInt("MaxCUHeight", &this.m_uiMaxCUHeight, uint(64), ""))
    // todo: remove defaults from MaxCUSize
    opts.AddOption(NewOptionUInt("MaxCUSize", &this.m_uiMaxCUWidth, uint(64), "Maximum CU size"))
    opts.AddOption(NewOptionUInt("MaxCUSize", &this.m_uiMaxCUHeight, uint(64), "Maximum CU size"))
    opts.AddOption(NewOptionUInt("MaxPartitionDepth", &this.m_uiMaxCUDepth, uint(4), "CU depth"))
    opts.AddOption(NewOptionUInt("QuadtreeTULog2MaxSize", &this.m_uiQuadtreeTULog2MaxSize, uint(6), "Maximum TU size in logarithm base 2"))
    opts.AddOption(NewOptionUInt("QuadtreeTULog2MinSize", &this.m_uiQuadtreeTULog2MinSize, uint(2), "Minimum TU size in logarithm base 2"))
    opts.AddOption(NewOptionUInt("QuadtreeTUMaxDepthIntra", &this.m_uiQuadtreeTUMaxDepthIntra, uint(1), "Depth of TU tree for intra CUs"))
    opts.AddOption(NewOptionUInt("QuadtreeTUMaxDepthInter", &this.m_uiQuadtreeTUMaxDepthInter, uint(2), "Depth of TU tree for inter CUs"))

    // Coding structure paramters
    opts.AddOption(NewOptionInt("IntraPeriod", &this.m_iIntraPeriod, -1, "Intra period in frames, (-1: only first frame)"))
    opts.AddOption(NewOptionInt("DecodingRefreshType", &this.m_iDecodingRefreshType, 0, "Intra refresh type (0:none 1:CRA 2:IDR)"))
    opts.AddOption(NewOptionInt("GOPSize", &this.m_iGOPSize, 1, "GOP size of temporal structure"))
    opts.AddOption(NewOptionBool("ListCombination", &this.m_bUseLComb, true, "Combined reference list for uni-prediction estimation in B-slices"))

    // motion options
    opts.AddOption(NewOptionInt("FastSearch", &this.m_iFastSearch, 1, "0:Full search  1:Diamond  2:PMVFAST"))
    opts.AddOption(NewOptionInt("SearchRange", &this.m_iSearchRange, 96, "Motion search range"))
    opts.AddOption(NewOptionInt("BipredSearchRange", &this.m_bipredSearchRange, 4, "Motion search range for bipred refinement"))
    opts.AddOption(NewOptionBool("HadamardME", &this.m_bUseHADME, true, "Hadamard ME for fractional-pel"))
    opts.AddOption(NewOptionBool("ASR", &this.m_bUseASR, false, "Adaptive motion search range"))

    // Mode decision parameters
    opts.AddOption(NewOptionFloat("LambdaModifier0", &this.m_adLambdaModifier[0], 1.0, "Lambda modifier for temporal layer 0"))
    opts.AddOption(NewOptionFloat("LambdaModifier1", &this.m_adLambdaModifier[1], 1.0, "Lambda modifier for temporal layer 1"))
    opts.AddOption(NewOptionFloat("LambdaModifier2", &this.m_adLambdaModifier[2], 1.0, "Lambda modifier for temporal layer 2"))
    opts.AddOption(NewOptionFloat("LambdaModifier3", &this.m_adLambdaModifier[3], 1.0, "Lambda modifier for temporal layer 3"))
    opts.AddOption(NewOptionFloat("LambdaModifier4", &this.m_adLambdaModifier[4], 1.0, "Lambda modifier for temporal layer 4"))
    opts.AddOption(NewOptionFloat("LambdaModifier5", &this.m_adLambdaModifier[5], 1.0, "Lambda modifier for temporal layer 5"))
    opts.AddOption(NewOptionFloat("LambdaModifier6", &this.m_adLambdaModifier[6], 1.0, "Lambda modifier for temporal layer 6"))
    opts.AddOption(NewOptionFloat("LambdaModifier7", &this.m_adLambdaModifier[7], 1.0, "Lambda modifier for temporal layer 7"))

    /* Quantization parameters */
    opts.AddOption(NewOptionFloat("QP", &this.m_fQP, 30.0, "Qp value, if value is float, QP is switched once during encoding"))
    opts.AddOption(NewOptionUInt("DeltaQpRD", &this.m_uiDeltaQpRD, uint(0), "max dQp offset for slice"))
    opts.AddOption(NewOptionInt("MaxDeltaQP", &this.m_iMaxDeltaQP, 0, "max dQp offset for block"))
    opts.AddOption(NewOptionInt("MaxCuDQPDepth", &this.m_iMaxCuDQPDepth, 0, "max depth for a minimum CuDQP"))
    opts.AddOption(NewOptionInt("CbQpOffset", &this.m_cbQpOffset, 0, "Chroma Cb QP Offset"))
    opts.AddOption(NewOptionInt("CrQpOffset", &this.m_crQpOffset, 0, "Chroma Cr QP Offset"))

    //#if ADAPTIVE_QP_SELECTION
    opts.AddOption(NewOptionBool("AdaptiveQpSelection", &this.m_bUseAdaptQpSelect, false, "AdaptiveQpSelection"))
    //#endif

    opts.AddOption(NewOptionBool("AdaptiveQP", &this.m_bUseAdaptiveQP, false, "QP adaptation based on a psycho-visual model"))
    opts.AddOption(NewOptionInt("MaxQPAdaptationRange", &this.m_iQPAdaptationRange, 6, "QP adaptation range"))
    opts.AddOption(NewOptionString("dQPFile", &cfg_dQPFile, string(""), "dQP file name"))
    opts.AddOption(NewOptionBool("RDOQ", &this.m_useRDOQ, true, ""))
    opts.AddOption(NewOptionBool("RDOQTS", &this.m_useRDOQTS, true, ""))
    //#if L0232_RD_PENALTY
    opts.AddOption(NewOptionUInt("RDpenalty", &this.m_rdPenalty, uint(0), "RD-penalty for 32x32 TU for intra in non-intra slices. 0:disbaled  1:RD-penalty  2:maximum RD-penalty"))
    //#endif
    // Entropy coding parameters
    opts.AddOption(NewOptionBool("SBACRD", &this.m_bUseSBACRD, true, "SBAC based RD estimation"))

    // Deblocking filter parameters
    opts.AddOption(NewOptionBool("LoopFilterDisable", &this.m_bLoopFilterDisable, false, ""))
    opts.AddOption(NewOptionBool("LoopFilterOffsetInPPS", &this.m_loopFilterOffsetInPPS, false, ""))
    opts.AddOption(NewOptionInt("LoopFilterBetaOffset_div2", &this.m_loopFilterBetaOffsetDiv2, 0, ""))
    opts.AddOption(NewOptionInt("LoopFilterTcOffset_div2", &this.m_loopFilterTcOffsetDiv2, 0, ""))
    opts.AddOption(NewOptionBool("DeblockingFilterControlPresent", &this.m_DeblockingFilterControlPresent, false, ""))

    // Coding tools
    opts.AddOption(NewOptionBool("AMP", &this.m_enableAMP, true, "Enable asymmetric motion partitions"))
    opts.AddOption(NewOptionBool("TransformSkip", &this.m_useTransformSkip, false, "Intra transform skipping"))
    opts.AddOption(NewOptionBool("TransformSkipFast", &this.m_useTransformSkipFast, false, "Fast intra transform skipping"))
    opts.AddOption(NewOptionBool("SAO", &this.m_bUseSAO, true, "Enable Sample Adaptive Offset"))
    opts.AddOption(NewOptionInt("MaxNumOffsetsPerPic", &this.m_maxNumOffsetsPerPic, 2048, "Max number of SAO offset per picture (Default: 2048)"))
    opts.AddOption(NewOptionBool("SAOLcuBoundary", &this.m_saoLcuBoundary, false, "0: right/bottom LCU boundary areas skipped from SAO parameter estimation, 1: non-deblocked pixels are used for those areas"))
    opts.AddOption(NewOptionBool("SAOLcuBasedOptimization", &this.m_saoLcuBasedOptimization, true, "0: SAO picture-based optimization, 1: SAO LCU-based optimization "))
    opts.AddOption(NewOptionInt("SliceMode", &this.m_sliceMode, 0, "0: Disable all Recon slice limits, 1: Enforce max # of LCUs, 2: Enforce max # of bytes"))
    opts.AddOption(NewOptionInt("SliceArgument", &this.m_sliceArgument, 0, "if SliceMode==1 SliceArgument represents max # of LCUs. if SliceMode==2 SliceArgument represents max # of bytes."))
    opts.AddOption(NewOptionInt("SliceSegmentMode", &this.m_sliceSegmentMode, 0, "0: Disable all dependent slice limits, 1: Enforce max # of LCUs, 2: Enforce constraint based dependent slices"))
    opts.AddOption(NewOptionInt("SliceSegmentArgument", &this.m_sliceSegmentArgument, 0, "if DependentSliceMode==1 SliceArgument represents max # of LCUs. if DependentSliceMode==2 DependentSliceArgument represents max # of bins."))
    opts.AddOption(NewOptionBool("LFCrossSliceBoundaryFlag", &this.m_bLFCrossSliceBoundaryFlag, true, "LFCrossSliceBoundaryFlag"))
    opts.AddOption(NewOptionBool("ConstrainedIntraPred", &this.m_bUseConstrainedIntraPred, false, "Constrained Intra Prediction"))
    opts.AddOption(NewOptionBool("PCMEnabledFlag", &this.m_usePCM, false, ""))
    opts.AddOption(NewOptionUInt("PCMLog2MaxSize", &this.m_pcmLog2MaxSize, uint(5), ""))
    opts.AddOption(NewOptionUInt("PCMLog2MinSize", &this.m_uiPCMLog2MinSize, uint(3), ""))

    opts.AddOption(NewOptionBool("PCMInputBitDepthFlag", &this.m_bPCMInputBitDepthFlag, true, ""))
    opts.AddOption(NewOptionBool("PCMFilterDisableFlag", &this.m_bPCMFilterDisableFlag, false, ""))
    opts.AddOption(NewOptionBool("LosslessCuEnabled", &this.m_useLossless, false, ""))
    opts.AddOption(NewOptionBool("weighted_pred_flag", &this.m_useWeightedPred, false, "weighted prediction flag (P-Slices)"))
    opts.AddOption(NewOptionBool("weighted_bipred_flag", &this.m_useWeightedBiPred, false, "weighted bipred flag (B-Slices)"))
    opts.AddOption(NewOptionUInt("Log2ParallelMergeLevel", &this.m_log2ParallelMergeLevel, uint(2), "Parallel merge estimation region"))
    opts.AddOption(NewOptionInt("UniformSpacingIdc", &this.m_iUniformSpacingIdr, 0, "Indicates if the column and row boundaries are distributed uniformly"))
    opts.AddOption(NewOptionInt("NumTileColumnsMinus1", &this.m_iNumColumnsMinus1, 0, "Number of columns in a picture minus 1"))
    opts.AddOption(NewOptionString("ColumnWidthArray", &cfg_ColumnWidth, string(""), "Array containing ColumnWidth values in units of LCU"))
    opts.AddOption(NewOptionInt("NumTileRowsMinus1", &this.m_iNumRowsMinus1, 0, "Number of rows in a picture minus 1"))
    opts.AddOption(NewOptionString("RowHeightArray", &cfg_RowHeight, string(""), "Array containing RowHeight values in units of LCU"))
    opts.AddOption(NewOptionBool("LFCrossTileBoundaryFlag", &this.m_bLFCrossTileBoundaryFlag, true, "1: cross-tile-boundary loop filtering. 0:non-cross-tile-boundary loop filtering"))
    opts.AddOption(NewOptionInt("WaveFrontSynchro", &this.m_iWaveFrontSynchro, 0, "0: no synchro; 1 synchro with TR; 2 TRR etc"))
    opts.AddOption(NewOptionInt("ScalingList", &this.m_useScalingListId, 0, "0: no scaling list, 1: default scaling lists, 2: scaling lists specified in ScalingListFile"))
    opts.AddOption(NewOptionString("ScalingListFile", &cfg_ScalingListFile, string(""), "Scaling list file name"))
    opts.AddOption(NewOptionInt("SignHideFlag,-SBH", &this.m_signHideFlag, 1, ""))
    opts.AddOption(NewOptionUInt("MaxNumMergeCand", &this.m_maxNumMergeCand, uint(5), "Maximum number of merge candidates"))

    /* Misc. */
    opts.AddOption(NewOptionInt("SEIDecodedPictureHash", &this.m_decodedPictureHashSEIEnabled, 0, "Control generation of decode picture hash SEI messages\n"))
    opts.AddOption(NewOptionInt("SEIpictureDigest", &this.m_decodedPictureHashSEIEnabled, 0, "deprecated alias for SEIDecodedPictureHash"))
    opts.AddOption(NewOptionInt("TMVPMode", &this.m_TMVPModeId, 1, "TMVP mode 0: TMVP disable for all slices. 1: TMVP enable for all slices (default) 2: TMVP enable for certain slices only"))
    opts.AddOption(NewOptionBool("FEN", &this.m_bUseFastEnc, false, "fast encoder setting"))
    opts.AddOption(NewOptionBool("ECU", &this.m_bUseEarlyCU, false, "Early CU setting"))
    opts.AddOption(NewOptionBool("FDM", &this.m_useFastDecisionForMerge, true, "Fast decision for Merge RD Cost"))
    opts.AddOption(NewOptionBool("CFM", &this.m_bUseCbfFastMode, false, "Cbf fast mode setting"))
    opts.AddOption(NewOptionBool("ESD", &this.m_useEarlySkipDetection, false, "Early SKIP detection setting"))
    //#if RATE_CONTROL_LAMBDA_DOMAIN
    opts.AddOption(NewOptionBool("RateControl", &this.m_RCEnableRateControl, false, "Rate control: enable rate control"))
    opts.AddOption(NewOptionInt("TargetBitrate", &this.m_RCTargetBitrate, 0, "Rate control: target bitrate"))
    opts.AddOption(NewOptionBool("KeepHierarchicalBit", &this.m_RCKeepHierarchicalBit, false, "Rate control: keep hierarchical bit allocation in rate control algorithm"))
    opts.AddOption(NewOptionBool("LCULevelRateControl", &this.m_RCLCULevelRC, true, "Rate control: true: LCU level RC; false: picture level RC"))
    opts.AddOption(NewOptionBool("RCLCUSeparateModel", &this.m_RCUseLCUSeparateModel, true, "Rate control: use LCU level separate R-lambda model"))
    opts.AddOption(NewOptionInt("InitialQP", &this.m_RCInitialQP, 0, "Rate control: initial QP"))
    opts.AddOption(NewOptionBool("RCForceIntraQP", &this.m_RCForceIntraQP, false, "Rate control: force intra QP to be equal to initial QP"))
    /*#else
      ("RateCtrl,-rc", this.m_enableRateCtrl, false, "Rate control on/off")
      ("TargetBitrate,-tbr", this.m_targetBitrate, 0, "Input target bitrate")
      ("NumLCUInUnit,-nu", this.m_numLCUInUnit, 0, "Number of LCUs in an Unit")
    #endif*/

    opts.AddOption(NewOptionBool("TransquantBypassEnableFlag", &this.m_TransquantBypassEnableFlag, false, "transquant_bypass_enable_flag indicator in PPS"))
    opts.AddOption(NewOptionBool("CUTransquantBypassFlagValue", &this.m_CUTransquantBypassFlagValue, false, "Fixed cu_transquant_bypass_flag value, when transquant_bypass_enable_flag is enabled"))
    opts.AddOption(NewOptionBool("RecalculateQPAccordingToLambda", &this.m_recalculateQPAccordingToLambda, false, "Recalculate QP values according to lambda values. Do not suggest to be enabled in all intra case"))
    opts.AddOption(NewOptionBool("StrongIntraSmoothing", &this.m_useStrongIntraSmoothing, true, "Enable strong intra smoothing for 32x32 blocks"))
    opts.AddOption(NewOptionInt("SEIActiveParameterSets", &this.m_activeParameterSetsSEIEnabled, 0, "Control generation of active parameter sets SEI messages"))
    opts.AddOption(NewOptionBool("VuiParametersPresent", &this.m_vuiParametersPresentFlag, false, "Enable generation of vui_parameters()"))
    opts.AddOption(NewOptionBool("AspectRatioInfoPresent", &this.m_aspectRatioInfoPresentFlag, false, "Signals whether aspect_ratio_idc is present"))
    opts.AddOption(NewOptionInt("AspectRatioIdc", &this.m_aspectRatioIdc, 0, "aspect_ratio_idc"))
    opts.AddOption(NewOptionInt("SarWidth", &this.m_sarWidth, 0, "horizontal size of the sample aspect ratio"))
    opts.AddOption(NewOptionInt("SarHeight", &this.m_sarHeight, 0, "vertical size of the sample aspect ratio"))
    opts.AddOption(NewOptionBool("OverscanInfoPresent", &this.m_overscanInfoPresentFlag, false, "Indicates whether conformant decoded pictures are suitable for display using overscan\n"))
    opts.AddOption(NewOptionBool("OverscanAppropriate", &this.m_overscanAppropriateFlag, false, "Indicates whether conformant decoded pictures are suitable for display using overscan\n"))
    opts.AddOption(NewOptionBool("VideoSignalTypePresent", &this.m_videoSignalTypePresentFlag, false, "Signals whether video_format, video_full_range_flag, and colour_description_present_flag are present"))
    opts.AddOption(NewOptionInt("VideoFormat", &this.m_videoFormat, 5, "Indicates representation of pictures"))
    opts.AddOption(NewOptionBool("VideoFullRange", &this.m_videoFullRangeFlag, false, "Indicates the black level and range of luma and chroma signals"))
    opts.AddOption(NewOptionBool("ColourDescriptionPresent", &this.m_colourDescriptionPresentFlag, false, "Signals whether colour_primaries, transfer_characteristics and matrix_coefficients are present"))
    opts.AddOption(NewOptionInt("ColourPrimaries", &this.m_colourPrimaries, 2, "Indicates chromaticity coordinates of the source primaries"))
    opts.AddOption(NewOptionInt("TransferCharateristics", &this.m_transferCharacteristics, 2, "Indicates the opto-electronic transfer characteristics of the source"))
    opts.AddOption(NewOptionInt("MatrixCoefficients", &this.m_matrixCoefficients, 2, "Describes the matrix coefficients used in deriving luma and chroma from RGB primaries"))
    opts.AddOption(NewOptionBool("ChromaLocInfoPresent", &this.m_chromaLocInfoPresentFlag, false, "Signals whether chroma_sample_loc_type_top_field and chroma_sample_loc_type_bottothis.m_field are present"))
    opts.AddOption(NewOptionInt("ChromaSampleLocTypeTopField", &this.m_chromaSampleLocTypeTopField, 0, "Specifies the location of chroma samples for top field"))
    opts.AddOption(NewOptionInt("ChromaSampleLocTypeBottomField", &this.m_chromaSampleLocTypeBottomField, 0, "Specifies the location of chroma samples for bottom field"))
    opts.AddOption(NewOptionBool("NeutralChromaIndication", &this.m_neutralChromaIndicationFlag, false, "Indicates that the value of all decoded chroma samples is equal to 1<<(BitDepthCr-1)"))
    opts.AddOption(NewOptionBool("DefaultDisplayWindowFlag", &this.m_defaultDisplayWindowFlag, false, "Indicates the presence of the Default Window parameters"))
    opts.AddOption(NewOptionInt("DefDispWinLeftOffset", &this.m_defDispWinLeftOffset, 0, "Specifies the left offset of the default display window from the conformance window"))
    opts.AddOption(NewOptionInt("DefDispWinRightOffset", &this.m_defDispWinRightOffset, 0, "Specifies the right offset of the default display window from the conformance window"))
    opts.AddOption(NewOptionInt("DefDispWinTopOffset", &this.m_defDispWinTopOffset, 0, "Specifies the top offset of the default display window from the conformance window"))
    opts.AddOption(NewOptionInt("DefDispWinBottomOffset", &this.m_defDispWinBottomOffset, 0, "Specifies the bottom offset of the default display window from the conformance window"))
    opts.AddOption(NewOptionBool("FrameFieldInfoPresentFlag", &this.m_frameFieldInfoPresentFlag, false, "Indicates that pic_struct and field coding related values are present in picture timing SEI messages"))
    opts.AddOption(NewOptionBool("PocProportionalToTimingFlag", &this.m_pocProportionalToTimingFlag, false, "Indicates that the POC value is proportional to the output time w.r.t. first picture in CVS"))
    opts.AddOption(NewOptionInt("NumTicksPocDiffOneMinus1", &this.m_numTicksPocDiffOneMinus1, 0, "Number of ticks minus 1 that for a POC difference of one"))
    opts.AddOption(NewOptionBool("BitstreamRestriction", &this.m_bitstreamRestrictionFlag, false, "Signals whether bitstream restriction parameters are present"))
    opts.AddOption(NewOptionBool("TilesFixedStructure", &this.m_tilesFixedStructureFlag, false, "Indicates that each active picture parameter set has the same values of the syntax elements related to tiles"))
    opts.AddOption(NewOptionBool("MotionVectorsOverPicBoundaries", &this.m_motionVectorsOverPicBoundariesFlag, false, "Indicates that no samples outside the picture boundaries are used for inter prediction"))
    opts.AddOption(NewOptionInt("MaxBytesPerPicDenom", &this.m_maxBytesPerPicDenom, 2, "Indicates a number of bytes not exceeded by the sum of the sizes of the VCL NAL units associated with any coded picture"))
    opts.AddOption(NewOptionInt("MaxBitsPerMinCuDenom", &this.m_maxBitsPerMinCuDenom, 1, "Indicates an upper bound for the number of bits of coding_unit() data"))
    opts.AddOption(NewOptionInt("Log2MaxMvLengthHorizontal", &this.m_log2MaxMvLengthHorizontal, 15, "Indicate the maximum absolute value of a decoded horizontal MV component in quarter-pel luma units"))
    opts.AddOption(NewOptionInt("Log2MaxMvLengthVertical", &this.m_log2MaxMvLengthVertical, 15, "Indicate the maximum absolute value of a decoded vertical MV component in quarter-pel luma units"))
    opts.AddOption(NewOptionInt("SEIRecoveryPoint", &this.m_recoveryPointSEIEnabled, 0, "Control generation of recovery point SEI messages"))
    opts.AddOption(NewOptionInt("SEIBufferingPeriod", &this.m_bufferingPeriodSEIEnabled, 0, "Control generation of buffering period SEI messages"))
    opts.AddOption(NewOptionInt("SEIPictureTiming", &this.m_pictureTimingSEIEnabled, 0, "Control generation of picture timing SEI messages"))
    opts.AddOption(NewOptionInt("SEIFramePacking", &this.m_framePackingSEIEnabled, 0, "Control generation of frame packing SEI messages"))
    opts.AddOption(NewOptionInt("SEIFramePackingType", &this.m_framePackingSEIType, 0, "Define frame packing arrangement\n"))
    opts.AddOption(NewOptionInt("SEIFramePackingId", &this.m_framePackingSEIId, 0, "Id of frame packing SEI message for a given session"))
    opts.AddOption(NewOptionInt("SEIFramePackingQuincunx", &this.m_framePackingSEIQuincunx, 0, "Indicate the presence of a Quincunx type video frame"))
    opts.AddOption(NewOptionInt("SEIFramePackingInterpretation", &this.m_framePackingSEIInterpretation, 0, "Indicate the interpretation of the frame pair\n"))
    opts.AddOption(NewOptionInt("SEIDisplayOrientation", &this.m_displayOrientationSEIAngle, 0, "Control generation of display orientation SEI messages"))
    opts.AddOption(NewOptionInt("SEITemporalLevel0Index", &this.m_temporalLevel0IndexSEIEnabled, 0, "Control generation of temporal level 0 index SEI messages"))
    opts.AddOption(NewOptionInt("SEIGradualDecodingRefreshInfo", &this.m_gradualDecodingRefreshInfoEnabled, 0, "Control generation of gradual decoding refresh information SEI message"))
    opts.AddOption(NewOptionInt("SEIDecodingUnitInfo", &this.m_decodingUnitInfoSEIEnabled, 0, "Control generation of decoding unit information SEI message."))

    var emptyGOPEntry *TLibEncoder.GOPEntry
    emptyGOPEntry = nil
    for i := 1; i < TLibCommon.MAX_GOP+1; i++ {
        cOSS := fmt.Sprintf("Frame%d", i)
        //fmt.Printf("%s===\n", cOSS);
        opts.AddOption(NewOptionGOPEntry(cOSS, &this.m_GOPList[i-1], emptyGOPEntry, "GOPEntry"))
    }
    opts.SetDefaults()

    /*if (argc == 1)
      {
        // argc == 1: no options have been specified
        po::doHelp(cout, opts);
        return false;
      }*/

    // parse cfg file
    opts.ParseConfigFile(argv[2])
    
    if argc >= 4 {
        this.m_pchTraceFile = argv[3]
    }

    /*
     * Set any derived parameters
     */
    /* convert std::string to c string for compatability */
    this.m_pchInputFile = string(cfg_InputFile)
    this.m_pchBitstreamFile = string(cfg_BitstreamFile)
    this.m_pchReconFile = string(cfg_ReconFile)
    this.m_pchdQPFile = string(cfg_dQPFile)

    var err error

    //Char* pColumnWidth = cfg_ColumnWidth.empty() ? NULL: strdup(cfg_ColumnWidth.c_str());
    //Char* pRowHeight = cfg_RowHeight.empty() ? NULL : strdup(cfg_RowHeight.c_str());
    if this.m_iUniformSpacingIdr == 0 && this.m_iNumColumnsMinus1 > 0 {
        //char *columnWidth;
        //i:=0;
        this.m_pColumnWidth = make([]int, this.m_iNumColumnsMinus1)
        columnWidth := strings.Fields(string(cfg_ColumnWidth)) //strtok(pColumnWidth, " ,-");
        if len(columnWidth) != this.m_iNumColumnsMinus1 {
            return errors.New("The number of columns whose width are defined is not equal to the allowed number of columns.\n")
            //return false;
        } else {
            for i := 0; i < this.m_iNumColumnsMinus1; i++ {
                this.m_pColumnWidth[i], err = strconv.Atoi(columnWidth[i])
                if err != nil {
                    return errors.New("Can't convert columnWidth[i] to integer\n")
                    //return false;
                }
            }
        }
        /*while(columnWidth!=NULL)
          {
            if( i>=m_iNumColumnsMinus1 )
            {
              printf( "The number of columns whose width are defined is larger than the allowed number of columns.\n" );
              exit( EXIT_FAILURE );
            }
            *( m_pColumnWidth + i ) = atoi( columnWidth );
            columnWidth = strtok(NULL, " ,-");
            i++;
          }
          if( i<m_iNumColumnsMinus1 )
          {
            printf( "The width of some columns is not defined.\n" );
            exit( EXIT_FAILURE );
          }*/
    } else {
        this.m_pColumnWidth = nil
    }

    if this.m_iUniformSpacingIdr == 0 && this.m_iNumRowsMinus1 > 0 {
        //char *rowHeight;
        //int  i=0;
        this.m_pRowHeight = make([]int, this.m_iNumRowsMinus1)
        rowHeight := strings.Fields(string(cfg_RowHeight)) //strtok(pRowHeight, " ,-");
        if len(rowHeight) != this.m_iNumRowsMinus1 {
            return errors.New("The number of rows whose width are defined is not equal to the allowed number of rows.\n")
            //return false;
        } else {
            for i := 0; i < this.m_iNumRowsMinus1; i++ {
                this.m_pRowHeight[i], err = strconv.Atoi(rowHeight[i])
                if err != nil {
                    return errors.New("Can't convert rowHeight[i] to integer\n")
                    //return false;
                }
            }
        }
        /*while(rowHeight!=NULL)
          {
            if( i>=m_iNumRowsMinus1 )
            {
              printf( "The number of rows whose height are defined is larger than the allowed number of rows.\n" );
              exit( EXIT_FAILURE );
            }
            *( m_pRowHeight + i ) = atoi( rowHeight );
            rowHeight = strtok(NULL, " ,-");
            i++;
          }
          if( i<m_iNumRowsMinus1 )
          {
            printf( "The height of some rows is not defined.\n" );
            exit( EXIT_FAILURE );
         }*/
    } else {
        this.m_pRowHeight = nil
    }

    this.m_scalingListFile = string(cfg_ScalingListFile) //.empty() ? NULL : strdup(cfg_ScalingListFile.c_str());

    /* rules for input, output and internal bitdepths as per help text */
    if this.m_internalBitDepthY == 0 {
        this.m_internalBitDepthY = this.m_inputBitDepthY
    }
    if this.m_internalBitDepthC == 0 {
        this.m_internalBitDepthC = this.m_internalBitDepthY
    }
    if this.m_inputBitDepthC == 0 {
        this.m_inputBitDepthC = this.m_inputBitDepthY
    }
    if this.m_outputBitDepthY == 0 {
        this.m_outputBitDepthY = this.m_internalBitDepthY
    }
    if this.m_outputBitDepthC == 0 {
        this.m_outputBitDepthC = this.m_internalBitDepthC
    }

    sps := TLibCommon.NewTComSPS()

    // TODO:ChromaFmt assumes 4:2:0 below
    switch this.m_conformanceMode {
    case 0:
        // no cropping or padding
        this.m_confLeft = 0
        this.m_confRight = 0
        this.m_confTop = 0
        this.m_confBottom = 0
        this.m_aiPad[1] = 0
        this.m_aiPad[0] = 0
    case 1:
        // automatic padding to minimum CU size
        minCuSize := int(this.m_uiMaxCUHeight >> (this.m_uiMaxCUDepth - 1))
        if this.m_iSourceWidth%minCuSize != 0 {
            this.m_aiPad[0] = ((this.m_iSourceWidth/minCuSize)+1)*minCuSize - this.m_iSourceWidth
            this.m_confRight = ((this.m_iSourceWidth/minCuSize)+1)*minCuSize - this.m_iSourceWidth
            this.m_iSourceWidth += this.m_confRight
        }
        if this.m_iSourceHeight%minCuSize != 0 {
            this.m_aiPad[1] = ((this.m_iSourceHeight/minCuSize)+1)*minCuSize - this.m_iSourceHeight
            this.m_confBottom = ((this.m_iSourceHeight/minCuSize)+1)*minCuSize - this.m_iSourceHeight
            this.m_iSourceHeight += this.m_confBottom
        }
        if this.m_aiPad[0]%sps.GetWinUnitX(TLibCommon.CHROMA_420) != 0 {
            return errors.New("Error: picture width is not an integer multiple of the specified chroma subsampling\n")
            //return false;
            //exit(EXIT_FAILURE);
        }
        if this.m_aiPad[1]%sps.GetWinUnitY(TLibCommon.CHROMA_420) != 0 {
            return errors.New("Error: picture height is not an integer multiple of the specified chroma subsampling\n")
            //return false;
            //exit(EXIT_FAILURE);
        }
    case 2:
        //padding
        this.m_iSourceWidth += this.m_aiPad[0]
        this.m_iSourceHeight += this.m_aiPad[1]
        this.m_confRight = this.m_aiPad[0]
        this.m_confBottom = this.m_aiPad[1]
    case 3:
        // cropping
        if (this.m_confLeft == 0) && (this.m_confRight == 0) && (this.m_confTop == 0) && (this.m_confBottom == 0) {
            fmt.Printf("Warning: Cropping enabled, but all cropping parameters set to zero\n")
        }
        if (this.m_aiPad[1] != 0) || (this.m_aiPad[0] != 0) {
            fmt.Printf("Warning: Cropping enabled, padding parameters will be ignored\n")
        }
        this.m_aiPad[1] = 0
        this.m_aiPad[0] = 0
    }

    // allocate slice-based dQP values
    this.m_aidQP = make([]int, this.m_framesToBeEncoded+this.m_iGOPSize+1)
    //::memset( this.m_aidQP, 0, sizeof(Int)*( this.m_iFrameToBeEncoded + this.m_iGOPSize + 1 ) );

    // handling of floating-point QP values
    // if QP is not integer, sequence is split into two sections having QP and QP+1
    this.m_iQP = int(this.m_fQP)
    if float64(this.m_iQP) < this.m_fQP {
        iSwitchPOC := int(float64(this.m_framesToBeEncoded) - (this.m_fQP-float64(this.m_iQP))*float64(this.m_framesToBeEncoded) + 0.5)

        iSwitchPOC = int(float64(iSwitchPOC)/float64(this.m_iGOPSize)+0.5) * this.m_iGOPSize
        for i := iSwitchPOC; i < this.m_framesToBeEncoded+this.m_iGOPSize+1; i++ {
            this.m_aidQP[i] = 1
        }
    }

    // reading external dQP description from file
    /*if this.m_pchdQPFile!="" {
      	fpt, err := os.Open(this.m_pchdQPFile);
    	if err!=nil {
    		log.Fatal(err)
    	}
    	defer fpt.Close()

        var iValue int;
        iPOC := 0;

        reader := bufio.NewReader(traceFile)
    	eof := false;
        for iPOC < this.m_framesToBeEncoded {
    	    line, err = reader.ReadString('\n')
    	    if err == io.EOF {
    	        break;
    	    } else if err != nil {
    	        return err
    	    }

            if ( fscanf(fpt, "%d", &iValue ) == EOF ) break;
            this.m_aidQP[ iPOC ] = iValue;
            iPOC++;
        }
      }*/

    if this.m_iWaveFrontSynchro != 0 {
        this.m_iWaveFrontSubstreams = int((uint(this.m_iSourceHeight) + this.m_uiMaxCUHeight - 1) / this.m_uiMaxCUHeight)
    } else {
        this.m_iWaveFrontSubstreams = 1
    }

    // check validity of input parameters
    if (this.xCheckParameter()==false) {
    	return errors.New("Error: invalidity of input parameters\n") 
    }

    // set global varibles
    this.xSetGlobal()

    // print-out parameters
    this.xPrintParameter()

    return nil
}

// internal member functions
func (this *TAppEncCfg) xSetGlobal() { ///< set global variables
    // set max CU width & height
    //TLibCommon.G_uiMaxCUWidth = this.m_uiMaxCUWidth
    //TLibCommon.G_uiMaxCUHeight = this.m_uiMaxCUHeight

    // compute actual CU depth with respect to config depth and max transform size
    this.m_uiAddCUDepth = 0
    for (this.m_uiMaxCUWidth >> this.m_uiMaxCUDepth) > (1 << (this.m_uiQuadtreeTULog2MinSize + this.m_uiAddCUDepth)) {
        this.m_uiAddCUDepth++
    }

    this.m_uiMaxCUDepth += this.m_uiAddCUDepth
    this.m_uiAddCUDepth++
    //TLibCommon.G_uiMaxCUDepth = this.m_uiMaxCUDepth
    //TLibCommon.G_uiAddCUDepth = this.m_uiAddCUDepth;

    // set internal bit-depth and constants
    TLibCommon.G_bitDepthY = this.m_internalBitDepthY
    TLibCommon.G_bitDepthC = this.m_internalBitDepthC

    if this.m_bPCMInputBitDepthFlag {
        TLibCommon.G_uiPCMBitDepthLuma = this.m_inputBitDepthY
        TLibCommon.G_uiPCMBitDepthChroma = this.m_inputBitDepthC
    } else {
        TLibCommon.G_uiPCMBitDepthLuma = this.m_internalBitDepthY
        TLibCommon.G_uiPCMBitDepthChroma = this.m_internalBitDepthC
    }
}

func (this *TAppEncCfg) xConfirmPara(bflag bool, message string) int {
    if !bflag {
        return 0
    }

    fmt.Printf("Error: %s\n", message)
    return 1
}

func (this *TAppEncCfg) xCheckParameter() bool { ///< check validity of configuration values
    if this.m_decodedPictureHashSEIEnabled == 0 {
        fmt.Printf("******************************************************************\n");
    	fmt.Printf("** WARNING: --SEIDecodedPictureHash is now disabled by default. **\n");
    	fmt.Printf("**          Automatic verification of decoded pictures by a     **\n");
    	fmt.Printf("**          decoder requires this option to be enabled.         **\n");
    	fmt.Printf("******************************************************************\n");
    }

    check_failed := 0 /* abort if there is a fatal configuration problem */
    //#define xConfirmPara(a,b) check_failed += confirmPara(a,b)
    // check range of parameters
    check_failed += this.xConfirmPara(this.m_inputBitDepthY < 8, "InputBitDepth must be at least 8")
    check_failed += this.xConfirmPara(this.m_inputBitDepthC < 8, "InputBitDepthC must be at least 8")
    check_failed += this.xConfirmPara(this.m_iFrameRate <= 0, "Frame rate must be more than 1")
    check_failed += this.xConfirmPara(this.m_framesToBeEncoded <= 0, "Total Number Of Frames encoded must be more than 0")
    check_failed += this.xConfirmPara(this.m_iGOPSize < 1, "GOP Size must be greater or equal to 1")
    check_failed += this.xConfirmPara(this.m_iGOPSize > 1 && this.m_iGOPSize%2 != 0, "GOP Size must be a multiple of 2, if GOP Size is greater than 1")
    check_failed += this.xConfirmPara((this.m_iIntraPeriod > 0 && this.m_iIntraPeriod < this.m_iGOPSize) || this.m_iIntraPeriod == 0, "Intra period must be more than GOP size, or -1 , not 0")
    check_failed += this.xConfirmPara(this.m_iDecodingRefreshType < 0 || this.m_iDecodingRefreshType > 2, "Decoding Refresh Type must be equal to 0, 1 or 2")
    check_failed += this.xConfirmPara(this.m_iQP < -6*(this.m_internalBitDepthY-8) || this.m_iQP > 51, "QP exceeds supported range (-QpBDOffsety to 51)")
    check_failed += this.xConfirmPara(this.m_loopFilterBetaOffsetDiv2 < -13 || this.m_loopFilterBetaOffsetDiv2 > 13, "Loop Filter Beta Offset div. 2 exceeds supported range (-13 to 13)")
    check_failed += this.xConfirmPara(this.m_loopFilterTcOffsetDiv2 < -13 || this.m_loopFilterTcOffsetDiv2 > 13, "Loop Filter Tc Offset div. 2 exceeds supported range (-13 to 13)")
    check_failed += this.xConfirmPara(this.m_iFastSearch < 0 || this.m_iFastSearch > 2, "Fast Search Mode is not supported value (0:Full search  1:Diamond  2:PMVFAST)")
    check_failed += this.xConfirmPara(this.m_iSearchRange < 0, "Search Range must be more than 0")
    check_failed += this.xConfirmPara(this.m_bipredSearchRange < 0, "Search Range must be more than 0")
    check_failed += this.xConfirmPara(this.m_iMaxDeltaQP > 7, "Absolute Delta QP exceeds supported range (0 to 7)")
    check_failed += this.xConfirmPara(this.m_iMaxCuDQPDepth > int(this.m_uiMaxCUDepth)-1, "Absolute depth for a minimum CuDQP exceeds maximum coding unit depth")
	check_failed += this.xConfirmPara(this.m_bUseSAO == true, "Current GoHM 10.0 don't support SAO")
    check_failed += this.xConfirmPara(this.m_cbQpOffset < -12, "Min. Chroma Cb QP Offset is -12")
    check_failed += this.xConfirmPara(this.m_cbQpOffset > 12, "Max. Chroma Cb QP Offset is  12")
    check_failed += this.xConfirmPara(this.m_crQpOffset < -12, "Min. Chroma Cr QP Offset is -12")
    check_failed += this.xConfirmPara(this.m_crQpOffset > 12, "Max. Chroma Cr QP Offset is  12")

    check_failed += this.xConfirmPara(this.m_iQPAdaptationRange <= 0, "QP Adaptation Range must be more than 0")
    if this.m_iDecodingRefreshType == 2 {
        check_failed += this.xConfirmPara(this.m_iIntraPeriod > 0 && this.m_iIntraPeriod <= this.m_iGOPSize, "Intra period must be larger than GOP size for periodic IDR pictures")
    }
    check_failed += this.xConfirmPara((this.m_uiMaxCUWidth>>this.m_uiMaxCUDepth) < 4, "Minimum partition width size should be larger than or equal to 8")
    check_failed += this.xConfirmPara((this.m_uiMaxCUHeight>>this.m_uiMaxCUDepth) < 4, "Minimum partition height size should be larger than or equal to 8")
    check_failed += this.xConfirmPara(this.m_uiMaxCUWidth < 16, "Maximum partition width size should be larger than or equal to 16")
    check_failed += this.xConfirmPara(this.m_uiMaxCUHeight < 16, "Maximum partition height size should be larger than or equal to 16")
    check_failed += this.xConfirmPara((uint(this.m_iSourceWidth)%(this.m_uiMaxCUWidth>>(this.m_uiMaxCUDepth-1))) != 0, "Resulting coded frame width must be a multiple of the minimum CU size")
    check_failed += this.xConfirmPara((uint(this.m_iSourceHeight)%(this.m_uiMaxCUHeight>>(this.m_uiMaxCUDepth-1))) != 0, "Resulting coded frame height must be a multiple of the minimum CU size")

    check_failed += this.xConfirmPara(this.m_uiQuadtreeTULog2MinSize < 2, "QuadtreeTULog2MinSize must be 2 or greater.")
    check_failed += this.xConfirmPara(this.m_uiQuadtreeTULog2MaxSize > 5, "QuadtreeTULog2MaxSize must be 5 or smaller.")
    check_failed += this.xConfirmPara((1<<this.m_uiQuadtreeTULog2MaxSize) > this.m_uiMaxCUWidth, "QuadtreeTULog2MaxSize must be log2(maxCUSize) or smaller.")
    check_failed += this.xConfirmPara(this.m_uiQuadtreeTULog2MaxSize < this.m_uiQuadtreeTULog2MinSize, "QuadtreeTULog2MaxSize must be greater than or equal to this.m_uiQuadtreeTULog2MinSize.")
    check_failed += this.xConfirmPara((1<<this.m_uiQuadtreeTULog2MinSize) > (this.m_uiMaxCUWidth>>(this.m_uiMaxCUDepth-1)), "QuadtreeTULog2MinSize must not be greater than minimum CU size")  // HS
    check_failed += this.xConfirmPara((1<<this.m_uiQuadtreeTULog2MinSize) > (this.m_uiMaxCUHeight>>(this.m_uiMaxCUDepth-1)), "QuadtreeTULog2MinSize must not be greater than minimum CU size") // HS
    check_failed += this.xConfirmPara((1<<this.m_uiQuadtreeTULog2MinSize) > (this.m_uiMaxCUWidth>>this.m_uiMaxCUDepth), "Minimum CU width must be greater than minimum transform size.")
    check_failed += this.xConfirmPara((1<<this.m_uiQuadtreeTULog2MinSize) > (this.m_uiMaxCUHeight>>this.m_uiMaxCUDepth), "Minimum CU height must be greater than minimum transform size.")
    check_failed += this.xConfirmPara(this.m_uiQuadtreeTUMaxDepthInter < 1, "QuadtreeTUMaxDepthInter must be greater than or equal to 1")
    check_failed += this.xConfirmPara(this.m_uiMaxCUWidth < ( 1 << (this.m_uiQuadtreeTULog2MinSize + this.m_uiQuadtreeTUMaxDepthInter - 1) ), "QuadtreeTUMaxDepthInter must be less than or equal to the difference between log2(maxCUSize) and QuadtreeTULog2MinSize plus 1")
    check_failed += this.xConfirmPara(this.m_uiQuadtreeTUMaxDepthIntra < 1, "QuadtreeTUMaxDepthIntra must be greater than or equal to 1")
    check_failed += this.xConfirmPara(this.m_uiMaxCUWidth < ( 1 << (this.m_uiQuadtreeTULog2MinSize + this.m_uiQuadtreeTUMaxDepthIntra - 1) ), "QuadtreeTUMaxDepthIntra must be less than or equal to the difference between log2(maxCUSize) and QuadtreeTULog2MinSize plus 1")

    check_failed += this.xConfirmPara(this.m_maxNumMergeCand < 1, "MaxNumMergeCand must be 1 or greater.")
    check_failed += this.xConfirmPara(this.m_maxNumMergeCand > 5, "MaxNumMergeCand must be 5 or smaller.")

    //#if ADAPTIVE_QP_SELECTION
    check_failed += this.xConfirmPara(this.m_bUseAdaptQpSelect == true && this.m_iQP < 0, "AdaptiveQpSelection must be disabled when QP < 0.")
    check_failed += this.xConfirmPara(this.m_bUseAdaptQpSelect == true && (this.m_cbQpOffset != 0 || this.m_crQpOffset != 0), "AdaptiveQpSelection must be disabled when ChromaQpOffset is not equal to 0.")
    //#endif

    if this.m_usePCM {
        check_failed += this.xConfirmPara(this.m_uiPCMLog2MinSize < 3, "PCMLog2MinSize must be 3 or greater.")
        check_failed += this.xConfirmPara(this.m_uiPCMLog2MinSize > 5, "PCMLog2MinSize must be 5 or smaller.")
        check_failed += this.xConfirmPara(this.m_pcmLog2MaxSize > 5, "PCMLog2MaxSize must be 5 or smaller.")
        check_failed += this.xConfirmPara(this.m_pcmLog2MaxSize < this.m_uiPCMLog2MinSize, "PCMLog2MaxSize must be equal to or greater than this.m_uiPCMLog2MinSize.")
    }

    check_failed += this.xConfirmPara(this.m_sliceMode < 0 || this.m_sliceMode > 3, "SliceMode exceeds supported range (0 to 3)")
    if this.m_sliceMode != 0 {
        check_failed += this.xConfirmPara(this.m_sliceArgument < 1, "SliceArgument should be larger than or equal to 1")
    }
    check_failed += this.xConfirmPara(this.m_sliceSegmentMode < 0 || this.m_sliceSegmentMode > 2, "SliceSegmentMode exceeds supported range (0 to 2)")
    if this.m_sliceSegmentMode != 0 {
        check_failed += this.xConfirmPara(this.m_sliceSegmentArgument < 1, "SliceSegmentArgument should be larger than or equal to 1")
    }

    tileFlag := (this.m_iNumColumnsMinus1 > 0 || this.m_iNumRowsMinus1 > 0)
    check_failed += this.xConfirmPara(tileFlag && this.m_iWaveFrontSynchro != 0, "Tile and Wavefront can not be applied together")

    //TODO:ChromaFmt assumes 4:2:0 below
    sps := TLibCommon.NewTComSPS()
    check_failed += this.xConfirmPara(this.m_iSourceWidth%sps.GetWinUnitX(TLibCommon.CHROMA_420) != 0, "Picture width must be an integer multiple of the specified chroma subsampling")
    check_failed += this.xConfirmPara(this.m_iSourceHeight%sps.GetWinUnitY(TLibCommon.CHROMA_420) != 0, "Picture height must be an integer multiple of the specified chroma subsampling")

    check_failed += this.xConfirmPara(this.m_aiPad[0]%sps.GetWinUnitX(TLibCommon.CHROMA_420) != 0, "Horizontal padding must be an integer multiple of the specified chroma subsampling")
    check_failed += this.xConfirmPara(this.m_aiPad[1]%sps.GetWinUnitY(TLibCommon.CHROMA_420) != 0, "Vertical padding must be an integer multiple of the specified chroma subsampling")

    check_failed += this.xConfirmPara(this.m_confLeft%sps.GetWinUnitX(TLibCommon.CHROMA_420) != 0, "Left cropping must be an integer multiple of the specified chroma subsampling")
    check_failed += this.xConfirmPara(this.m_confRight%sps.GetWinUnitX(TLibCommon.CHROMA_420) != 0, "Right cropping must be an integer multiple of the specified chroma subsampling")
    check_failed += this.xConfirmPara(this.m_confTop%sps.GetWinUnitY(TLibCommon.CHROMA_420) != 0, "Top cropping must be an integer multiple of the specified chroma subsampling")
    check_failed += this.xConfirmPara(this.m_confBottom%sps.GetWinUnitY(TLibCommon.CHROMA_420) != 0, "Bottom cropping must be an integer multiple of the specified chroma subsampling")

    // max CU width and height should be power of 2
    ui := this.m_uiMaxCUWidth
    for ui != 0 {
        ui >>= 1
        if (ui & 1) == 1 {
            check_failed += this.xConfirmPara(ui != 1, "Width should be 2^n")
        }
    }
    ui = this.m_uiMaxCUHeight
    for ui != 0 {
        ui >>= 1
        if (ui & 1) == 1 {
            check_failed += this.xConfirmPara(ui != 1, "Height should be 2^n")
        }
    }

    /* if this is an intra-only sequence, ie IntraPeriod=1, don't verify the GOP structure
     * This permits the ability to omit a GOP structure specification */
    if this.m_iIntraPeriod == 1 {
        // && this.m_GOPList[0].GetPOC() == -1 {
        this.m_GOPList[0] = TLibEncoder.NewGOPEntry()
        this.m_GOPList[0].SetQPFactor(1)

        this.m_GOPList[0].SetBetaOffsetDiv2(0)
        this.m_GOPList[0].SetTcOffsetDiv2(0)

        this.m_GOPList[0].SetPOC(1)
        this.m_GOPList[0].SetNumRefPicsActive(4)
    }

    verifiedGOP := false
    errorGOP := false
    checkGOP := 1
    numRefs := 1
    var refList [TLibCommon.MAX_NUM_REF_PICS + 1]int
    refList[0] = 0
    var isOK [TLibCommon.MAX_GOP]bool
    for i := 0; i < TLibCommon.MAX_GOP; i++ {
        isOK[i] = false
    }
    numOK := 0
    check_failed += this.xConfirmPara(this.m_iIntraPeriod >= 0 && (this.m_iIntraPeriod%this.m_iGOPSize != 0), "Intra period must be a multiple of GOPSize, or -1")
	
	//fmt.Printf("#        Type POC QPoffset QPfactor tcOffsetDiv2 betaOffsetDiv2  temporal_id #ref_pics_active #ref_pics reference pictures predict deltaRPS #ref_idcs reference idcs\n");
    for i := 0; i < this.m_iGOPSize; i++ {
        if this.m_GOPList[i].GetPOC() == this.m_iGOPSize {
            check_failed += this.xConfirmPara(this.m_GOPList[i].GetTemporalId() != 0, "The last frame in each GOP must have temporal ID = 0 ")
        }
        
        //fmt.Printf("%s %d %d %f %d %d %d %d %d %d \n",this.m_GOPList[i].GetSliceType(), this.m_GOPList[i].GetPOC(), this.m_GOPList[i].GetQPOffset(),
        //   this.m_GOPList[i].GetQPFactor(), this.m_GOPList[i].GetTcOffsetDiv2(), this.m_GOPList[i].GetBetaOffsetDiv2(), this.m_GOPList[i].GetTemporalId(),
        //   this.m_GOPList[i].GetNumRefPicsActive(), this.m_GOPList[i].GetNumRefPics(), this.m_GOPList[i].GetInterRPSPrediction());
    }

    if (this.m_iIntraPeriod != 1) && !this.m_loopFilterOffsetInPPS && this.m_DeblockingFilterControlPresent && (!this.m_bLoopFilterDisable) {
        for i := 0; i < this.m_iGOPSize; i++ {
            check_failed += this.xConfirmPara((this.m_GOPList[i].GetBetaOffsetDiv2()+this.m_loopFilterBetaOffsetDiv2) < -6 || (this.m_GOPList[i].GetBetaOffsetDiv2()+this.m_loopFilterBetaOffsetDiv2) > 6, "Loop Filter Beta Offset div. 2 for one of the GOP entries exceeds supported range (-6 to 6)")
            check_failed += this.xConfirmPara((this.m_GOPList[i].GetTcOffsetDiv2()+this.m_loopFilterTcOffsetDiv2) < -6 || (this.m_GOPList[i].GetTcOffsetDiv2()+this.m_loopFilterTcOffsetDiv2) > 6, "Loop Filter Tc Offset div. 2 for one of the GOP entries exceeds supported range (-6 to 6)")
        }
    }

    this.m_extraRPSs = 0
    //start looping through frames in coding order until we can verify that the GOP structure is correct.
    for !verifiedGOP && !errorGOP {
        curGOP := (checkGOP - 1) % this.m_iGOPSize
        curPOC := ((checkGOP-1)/this.m_iGOPSize)*this.m_iGOPSize + this.m_GOPList[curGOP].GetPOC()
        //fmt.Printf("curPOC%d=checkGOP%d/m_iGOPSize%d +m_GOPList[%d].m_POC%d/m_numRefPics%d\n", curPOC,checkGOP,this.m_iGOPSize,curGOP,this.m_GOPList[curGOP].GetPOC(),this.m_GOPList[curGOP].GetNumRefPics());
        if this.m_GOPList[curGOP].GetPOC() < 0 {
            fmt.Printf("\nError: found fewer Reference Picture Sets than GOPSize\n")
            errorGOP = true
        } else {
            //check that all reference pictures are available, or have a POC < 0 meaning they might be available in the next GOP.
            beforeI := false
            for i := 0; i < this.m_GOPList[curGOP].GetNumRefPics(); i++ {
                absPOC := curPOC + this.m_GOPList[curGOP].GetReferencePics(i)
                //fmt.Printf("absPOC %d = curPOC %d + m_referencePics %d\n", absPOC, curGOP, this.m_GOPList[curGOP].GetReferencePics(i));
                if absPOC < 0 {
                    beforeI = true
                } else {
                    found := false
                    for j := 0; j < numRefs; j++ {
                        if refList[j] == absPOC {
                            found = true
                            for k := 0; k < this.m_iGOPSize; k++ {
                                if absPOC%this.m_iGOPSize == this.m_GOPList[k].GetPOC()%this.m_iGOPSize {
                                    if this.m_GOPList[k].GetTemporalId() == this.m_GOPList[curGOP].GetTemporalId() {
                                        this.m_GOPList[k].SetRefPic(true)
                                    }
                                    this.m_GOPList[curGOP].SetUsedByCurrPic(i, this.m_GOPList[k].GetTemporalId() <= this.m_GOPList[curGOP].GetTemporalId())
                                }
                            }
                        }
                    }
                    if !found {
                        fmt.Printf("\nError: ref pic %d is not available for GOP frame %d\n", this.m_GOPList[curGOP].GetReferencePics(i), curGOP+1)
                        errorGOP = true
                    }
                }
            }
            if !beforeI && !errorGOP {
                //all ref frames were present
                if !isOK[curGOP] {
                    numOK++
                    isOK[curGOP] = true
                    if numOK == this.m_iGOPSize {
                        verifiedGOP = true
                    }
                }
            } else {
                //create a new GOPEntry for this frame containing all the reference pictures that were available (POC > 0)
                this.m_GOPList[this.m_iGOPSize+this.m_extraRPSs] = TLibEncoder.NewGOPEntry();
                *(this.m_GOPList[this.m_iGOPSize+this.m_extraRPSs]) = *(this.m_GOPList[curGOP])
                newRefs := 0
                for i := 0; i < this.m_GOPList[curGOP].GetNumRefPics(); i++ {
                    absPOC := curPOC + this.m_GOPList[curGOP].GetReferencePics(i)
                    if absPOC >= 0 {
                        this.m_GOPList[this.m_iGOPSize+this.m_extraRPSs].SetReferencePics(newRefs, this.m_GOPList[curGOP].GetReferencePics(i))
                        this.m_GOPList[this.m_iGOPSize+this.m_extraRPSs].SetUsedByCurrPic(newRefs, this.m_GOPList[curGOP].GetUsedByCurrPic(i))
                        newRefs++
                    }
                }
                numPrefRefs := this.m_GOPList[curGOP].GetNumRefPicsActive()

                for offset := -1; offset > -checkGOP; offset-- {
                    //step backwards in coding order and include any extra available pictures we might find useful to replace the ones with POC < 0.
                    offGOP := (checkGOP - 1 + offset) % this.m_iGOPSize
                    offPOC := ((checkGOP-1+offset)/this.m_iGOPSize)*this.m_iGOPSize + this.m_GOPList[offGOP].GetPOC()
                    if offPOC >= 0 && this.m_GOPList[offGOP].GetTemporalId() <= this.m_GOPList[curGOP].GetTemporalId() {
                        newRef := false
                        for i := 0; i < numRefs; i++ {
                            if refList[i] == offPOC {
                                newRef = true
                            }
                        }
                        for i := 0; i < newRefs; i++ {
                            if this.m_GOPList[this.m_iGOPSize+this.m_extraRPSs].GetReferencePics(i) == offPOC-curPOC {
                                newRef = false
                            }
                        }
                        if newRef {
                            insertPoint := newRefs
                            //this picture can be added, find appropriate place in list and insert it.
                            if this.m_GOPList[offGOP].GetTemporalId() == this.m_GOPList[curGOP].GetTemporalId() {
                                this.m_GOPList[offGOP].SetRefPic(true)
                            }
                            for j := 0; j < newRefs; j++ {
                                if this.m_GOPList[this.m_iGOPSize+this.m_extraRPSs].GetReferencePics(j) < offPOC-curPOC || this.m_GOPList[this.m_iGOPSize+this.m_extraRPSs].GetReferencePics(j) > 0 {
                                    insertPoint = j
                                    break
                                }
                            }
                            prev := offPOC - curPOC
                            prevUsed := this.m_GOPList[offGOP].GetTemporalId() <= this.m_GOPList[curGOP].GetTemporalId()
                            for j := insertPoint; j < newRefs+1; j++ {
                                newPrev := this.m_GOPList[this.m_iGOPSize+this.m_extraRPSs].GetReferencePics(j)
                                newUsed := this.m_GOPList[this.m_iGOPSize+this.m_extraRPSs].GetUsedByCurrPic(j)
                                this.m_GOPList[this.m_iGOPSize+this.m_extraRPSs].SetReferencePics(j, prev)
                                this.m_GOPList[this.m_iGOPSize+this.m_extraRPSs].SetUsedByCurrPic(j, prevUsed)
                                prevUsed = newUsed
                                prev = newPrev
                            }
                            newRefs++
                        }
                    }
                    if newRefs >= numPrefRefs {
                        break
                    }
                }
                this.m_GOPList[this.m_iGOPSize+this.m_extraRPSs].SetNumRefPics(newRefs)
                this.m_GOPList[this.m_iGOPSize+this.m_extraRPSs].SetPOC(curPOC)
                //fmt.Printf("m_GOPList[m_iGOPSize%d+m_extraRPSs%d].m_numRefPics%d VS %d\n", this.m_iGOPSize, this.m_extraRPSs, newRefs, this.m_GOPList[curGOP].GetNumRefPics());
                if this.m_extraRPSs == 0 {
                    this.m_GOPList[this.m_iGOPSize+this.m_extraRPSs].SetInterRPSPrediction(0)
                    this.m_GOPList[this.m_iGOPSize+this.m_extraRPSs].SetNumRefIdc(0)
                } else {
                    rIdx := this.m_iGOPSize + this.m_extraRPSs - 1
                    refPOC := this.m_GOPList[rIdx].GetPOC()
                    refPics := this.m_GOPList[rIdx].GetNumRefPics()
                    newIdc := 0
                    for i := 0; i <= refPics; i++ {
                        var deltaPOC int
                        if i != refPics {
                            deltaPOC = this.m_GOPList[rIdx].GetReferencePics(i) // check if the reference abs POC is >= 0
                        } else {
                            deltaPOC = 0 // check if the reference abs POC is >= 0
                        }
                        absPOCref := refPOC + deltaPOC
                        refIdc := 0
                        for j := 0; j < this.m_GOPList[this.m_iGOPSize+this.m_extraRPSs].GetNumRefPics(); j++ {
                            if (absPOCref - curPOC) == this.m_GOPList[this.m_iGOPSize+this.m_extraRPSs].GetReferencePics(j) {
                                if this.m_GOPList[this.m_iGOPSize+this.m_extraRPSs].GetUsedByCurrPic(j) {
                                    refIdc = 1
                                } else {
                                    refIdc = 2
                                }
                            }
                        }
                        this.m_GOPList[this.m_iGOPSize+this.m_extraRPSs].SetRefIdc(newIdc, refIdc)
                        newIdc++
                    }
                    this.m_GOPList[this.m_iGOPSize+this.m_extraRPSs].SetInterRPSPrediction(1)
                    this.m_GOPList[this.m_iGOPSize+this.m_extraRPSs].SetNumRefIdc(newIdc)
                    this.m_GOPList[this.m_iGOPSize+this.m_extraRPSs].SetDeltaRPS(refPOC - this.m_GOPList[this.m_iGOPSize+this.m_extraRPSs].GetPOC())
                }
                curGOP = this.m_iGOPSize + this.m_extraRPSs
                this.m_extraRPSs++
            }
            numRefs = 0
            for i := 0; i < this.m_GOPList[curGOP].GetNumRefPics(); i++ {
                absPOC := curPOC + this.m_GOPList[curGOP].GetReferencePics(i)
                if absPOC >= 0 {
                    refList[numRefs] = absPOC
                    numRefs++
                }
            }
            refList[numRefs] = curPOC
            numRefs++
        }
        checkGOP++
    }
    //fmt.Printf("m_extraRPSs=%d\n", this.m_extraRPSs);
    
    check_failed += this.xConfirmPara(errorGOP, "Invalid GOP structure given")
    this.m_maxTempLayer = 1
    for i := 0; i < this.m_iGOPSize; i++ {
        if this.m_GOPList[i].GetTemporalId() >= this.m_maxTempLayer {
            this.m_maxTempLayer = this.m_GOPList[i].GetTemporalId() + 1
        }
        check_failed += this.xConfirmPara(this.m_GOPList[i].GetSliceType() != "B" && this.m_GOPList[i].GetSliceType() != "P", "Slice type must be equal to B or P")
    }
    for i := 0; i < TLibCommon.MAX_TLAYER; i++ {
        this.m_numReorderPics[i] = 0
        this.m_maxDecPicBuffering[i] = 0
    }
    for i := 0; i < this.m_iGOPSize; i++ {
        if this.m_GOPList[i].GetNumRefPics() > this.m_maxDecPicBuffering[this.m_GOPList[i].GetTemporalId()] {
            this.m_maxDecPicBuffering[this.m_GOPList[i].GetTemporalId()] = this.m_GOPList[i].GetNumRefPics()
        }
        highestDecodingNumberWithLowerPOC := 0
        for j := 0; j < this.m_iGOPSize; j++ {
            if this.m_GOPList[j].GetPOC() <= this.m_GOPList[i].GetPOC() {
                highestDecodingNumberWithLowerPOC = j
            }
        }
        numReorder := 0
        for j := 0; j < highestDecodingNumberWithLowerPOC; j++ {
            if this.m_GOPList[j].GetTemporalId() <= this.m_GOPList[i].GetTemporalId() && this.m_GOPList[j].GetPOC() > this.m_GOPList[i].GetPOC() {
                numReorder++
            }
        }
        if numReorder > this.m_numReorderPics[this.m_GOPList[i].GetTemporalId()] {
            this.m_numReorderPics[this.m_GOPList[i].GetTemporalId()] = numReorder
        }
    }
    for i := 0; i < TLibCommon.MAX_TLAYER-1; i++ {
        // a lower layer can not have higher value of this.m_numReorderPics than a higher layer
        if this.m_numReorderPics[i+1] < this.m_numReorderPics[i] {
            this.m_numReorderPics[i+1] = this.m_numReorderPics[i]
        }
        // the value of nuthis.m_reorder_pics[ i ] shall be in the range of 0 to max_dec_pic_buffering[ i ], inclusive
        if this.m_numReorderPics[i] > this.m_maxDecPicBuffering[i] {
            this.m_maxDecPicBuffering[i] = this.m_numReorderPics[i]
        }
        // a lower layer can not have higher value of this.m_uiMaxDecPicBuffering than a higher layer
        if this.m_maxDecPicBuffering[i+1] < this.m_maxDecPicBuffering[i] {
            this.m_maxDecPicBuffering[i+1] = this.m_maxDecPicBuffering[i]
        }
    }
    // the value of nuthis.m_reorder_pics[ i ] shall be in the range of 0 to max_dec_pic_buffering[ i ], inclusive
    if this.m_numReorderPics[TLibCommon.MAX_TLAYER-1] > this.m_maxDecPicBuffering[TLibCommon.MAX_TLAYER-1] {
        this.m_maxDecPicBuffering[TLibCommon.MAX_TLAYER-1] = this.m_numReorderPics[TLibCommon.MAX_TLAYER-1]
    }

    if this.m_vuiParametersPresentFlag && this.m_bitstreamRestrictionFlag {
        PicSizeInSamplesY := this.m_iSourceWidth * this.m_iSourceHeight
        if tileFlag {
            maxTileWidth := 0
            maxTileHeight := 0
            var widthInCU, heightInCU int
            if this.m_iSourceWidth%int(this.m_uiMaxCUWidth) != 0 {
                widthInCU = this.m_iSourceWidth/int(this.m_uiMaxCUWidth) + 1
            } else {
                widthInCU = this.m_iSourceWidth / int(this.m_uiMaxCUWidth)
            }
            if this.m_iSourceHeight%int(this.m_uiMaxCUHeight) != 0 {
                heightInCU = this.m_iSourceHeight/int(this.m_uiMaxCUHeight) + 1
            } else {
                heightInCU = this.m_iSourceHeight / int(this.m_uiMaxCUHeight)
            }

            if this.m_iUniformSpacingIdr != 0 {
                maxTileWidth = int(this.m_uiMaxCUWidth) * ((widthInCU + this.m_iNumColumnsMinus1) / (this.m_iNumColumnsMinus1 + 1))
                maxTileHeight = int(this.m_uiMaxCUHeight) * ((heightInCU + this.m_iNumRowsMinus1) / (this.m_iNumRowsMinus1 + 1))
                // if only the last tile-row is one treeblock higher than the others
                // the maxTileHeight becomes smaller if the last row of treeblocks has lower height than the others
                if ((heightInCU - 1) % (this.m_iNumRowsMinus1 + 1)) == 0 {
                    maxTileHeight = maxTileHeight - int(this.m_uiMaxCUHeight) + (this.m_iSourceHeight % int(this.m_uiMaxCUHeight))
                }
                // if only the last tile-column is one treeblock wider than the others
                // the maxTileWidth becomes smaller if the last column of treeblocks has lower width than the others
                if ((widthInCU - 1) % (this.m_iNumColumnsMinus1 + 1)) == 0 {
                    maxTileWidth = maxTileWidth - int(this.m_uiMaxCUWidth) + (this.m_iSourceWidth % int(this.m_uiMaxCUWidth))
                }
            } else { // not uniform spacing
                if this.m_iNumColumnsMinus1 < 1 {
                    maxTileWidth = this.m_iSourceWidth
                } else {
                    accColumnWidth := 0
                    for col := 0; col < (this.m_iNumColumnsMinus1); col++ {
                        if int(this.m_pColumnWidth[col]) > maxTileWidth {
                            maxTileWidth = int(this.m_pColumnWidth[col])
                        } else {
                            maxTileWidth = maxTileWidth
                        }
                        accColumnWidth += int(this.m_pColumnWidth[col])
                    }
                    if (widthInCU - accColumnWidth) > maxTileWidth {
                        maxTileWidth = int(this.m_uiMaxCUWidth) * (widthInCU - accColumnWidth)
                    } else {
                        maxTileWidth = int(this.m_uiMaxCUWidth) * maxTileWidth
                    }
                }
                if this.m_iNumRowsMinus1 < 1 {
                    maxTileHeight = this.m_iSourceHeight
                } else {
                    accRowHeight := 0
                    for row := 0; row < (this.m_iNumRowsMinus1); row++ {
                        if int(this.m_pRowHeight[row]) > maxTileHeight {
                            maxTileHeight = int(this.m_pRowHeight[row])
                        } else {
                            maxTileHeight = maxTileHeight
                        }
                        accRowHeight += int(this.m_pRowHeight[row])
                    }
                    if (heightInCU - accRowHeight) > maxTileHeight {
                        maxTileHeight = int(this.m_uiMaxCUHeight) * (heightInCU - accRowHeight)
                    } else {
                        maxTileHeight = int(this.m_uiMaxCUHeight) * maxTileHeight
                    }
                }
            }
            maxSizeInSamplesY := maxTileWidth * maxTileHeight
            this.m_minSpatialSegmentationIdc = 4*PicSizeInSamplesY/maxSizeInSamplesY - 4
        } else if this.m_iWaveFrontSynchro != 0 {
            this.m_minSpatialSegmentationIdc = 4*PicSizeInSamplesY/((2*this.m_iSourceHeight+this.m_iSourceWidth)*int(this.m_uiMaxCUHeight)) - 4
        } else if this.m_sliceMode == 1 {
            this.m_minSpatialSegmentationIdc = 4*PicSizeInSamplesY/(this.m_sliceArgument*int(this.m_uiMaxCUWidth*this.m_uiMaxCUHeight)) - 4
        } else {
            this.m_minSpatialSegmentationIdc = 0
        }
    }

    check_failed += this.xConfirmPara(this.m_bUseLComb == false && this.m_numReorderPics[TLibCommon.MAX_TLAYER-1] != 0, "ListCombination can only be 0 in low delay coding (more precisely when L0 and L1 are identical)") // Note however this is not the full necessary condition as ref_pic_list_combination_flag can only be 0 if L0 == L1.
    check_failed += this.xConfirmPara(this.m_iWaveFrontSynchro < 0, "WaveFrontSynchro cannot be negative")
    check_failed += this.xConfirmPara(this.m_iWaveFrontSubstreams <= 0, "WaveFrontSubstreams must be positive")
    check_failed += this.xConfirmPara(this.m_iWaveFrontSubstreams > 1 && this.m_iWaveFrontSynchro == 0, "Must have WaveFrontSynchro > 0 in order to have WaveFrontSubstreams > 1")

    check_failed += this.xConfirmPara(this.m_decodedPictureHashSEIEnabled < 0 || this.m_decodedPictureHashSEIEnabled > 3, "this hash type is not correct!\n")

    //#if RATE_CONTROL_LAMBDA_DOMAIN
    if this.m_RCEnableRateControl {
        if this.m_RCForceIntraQP {
            if this.m_RCInitialQP == 0 {
                fmt.Printf("\nInitial QP for rate control is not specified. Reset not to use force intra QP!")
                this.m_RCForceIntraQP = false
            }
        }
        check_failed += this.xConfirmPara(this.m_uiDeltaQpRD > 0, "Rate control cannot be used together with slice level multiple-QP optimization!\n")
    }
    /*#else
      if(this.m_enableRateCtrl)
      {
        Int numLCUInWidth  = (this.m_iSourceWidth  / this.m_uiMaxCUWidth) + (( this.m_iSourceWidth  %  this.m_uiMaxCUWidth ) ? 1 : 0);
        Int numLCUInHeight = (this.m_iSourceHeight / this.m_uiMaxCUHeight)+ (( this.m_iSourceHeight %  this.m_uiMaxCUHeight) ? 1 : 0);
        Int numLCUInPic    =  numLCUInWidth * numLCUInHeight;

        check_failed += this.xConfirmPara( (numLCUInPic % this.m_numLCUInUnit) != 0, "total number of LCUs in a frame should be completely divided by NumLCUInUnit" );

        this.m_iMaxDeltaQP       = MAX_DELTA_QP;
        this.m_iMaxCuDQPDepth    = MAX_CUDQP_DEPTH;
      }
    #endif*/

    check_failed += this.xConfirmPara(!this.m_TransquantBypassEnableFlag && this.m_CUTransquantBypassFlagValue, "CUTransquantBypassFlagValue cannot be 1 when TransquantBypassEnableFlag is 0")

    check_failed += this.xConfirmPara(this.m_log2ParallelMergeLevel < 2, "Log2ParallelMergeLevel should be larger than or equal to 2")

    //#if L0444_FPA_TYPE
    if this.m_framePackingSEIEnabled != 0 {
        check_failed += this.xConfirmPara(this.m_framePackingSEIType < 3 || this.m_framePackingSEIType > 5, "SEIFramePackingType must be in rage 3 to 5")
    }
    //#endif

    //#undef xConfirmPara
    if check_failed > 0 {
        return false //exit(EXIT_FAILURE);
    }

    return true
}

func (this *TAppEncCfg) xPrintParameter() { ///< print configuration values
    fmt.Printf("\n")
    fmt.Printf("Input          File          : %s\n", this.m_pchInputFile)
    fmt.Printf("Bitstream      File          : %s\n", this.m_pchBitstreamFile)
    fmt.Printf("Reconstruction File          : %s\n", this.m_pchReconFile)
    fmt.Printf("Real     Format              : %dx%d %dHz\n", this.m_iSourceWidth-this.m_confLeft-this.m_confRight, this.m_iSourceHeight-this.m_confTop-this.m_confBottom, this.m_iFrameRate)
    fmt.Printf("Internal Format              : %dx%d %dHz\n", this.m_iSourceWidth, this.m_iSourceHeight, this.m_iFrameRate)
    fmt.Printf("Frame index                  : %d - %d (%d frames)\n", this.m_FrameSkip, int(this.m_FrameSkip)+this.m_framesToBeEncoded-1, this.m_framesToBeEncoded)
    fmt.Printf("CU size / depth              : %d / %d\n", this.m_uiMaxCUWidth, this.m_uiMaxCUDepth)
    fmt.Printf("RQT trans. size (min / max)  : %d / %d\n", 1<<this.m_uiQuadtreeTULog2MinSize, 1<<this.m_uiQuadtreeTULog2MaxSize)
    fmt.Printf("Max RQT depth inter          : %d\n", this.m_uiQuadtreeTUMaxDepthInter)
    fmt.Printf("Max RQT depth intra          : %d\n", this.m_uiQuadtreeTUMaxDepthIntra)
    fmt.Printf("Min PCM size                 : %d\n", 1<<this.m_uiPCMLog2MinSize)
    fmt.Printf("Motion search range          : %d\n", this.m_iSearchRange)
    fmt.Printf("Intra period                 : %d\n", this.m_iIntraPeriod)
    fmt.Printf("Decoding refresh type        : %d\n", this.m_iDecodingRefreshType)
    fmt.Printf("QP                           : %5.2f\n", this.m_fQP)
    fmt.Printf("Max dQP signaling depth      : %d\n", this.m_iMaxCuDQPDepth)

    fmt.Printf("Cb QP Offset                 : %d\n", this.m_cbQpOffset)
    fmt.Printf("Cr QP Offset                 : %d\n", this.m_crQpOffset)

    if this.m_bUseAdaptiveQP {
        fmt.Printf("QP adaptation                : %d (range=%d)\n", TLibCommon.B2U(this.m_bUseAdaptiveQP), this.m_iQPAdaptationRange)
    } else {
        fmt.Printf("QP adaptation                : %d (range=%d)\n", TLibCommon.B2U(this.m_bUseAdaptiveQP), 0)
    }
    fmt.Printf("GOP size                     : %d\n", this.m_iGOPSize)
    fmt.Printf("Internal bit depth           : (Y:%d, C:%d)\n", this.m_internalBitDepthY, this.m_internalBitDepthC)
    fmt.Printf("PCM sample bit depth         : (Y:%d, C:%d)\n", TLibCommon.G_uiPCMBitDepthLuma, TLibCommon.G_uiPCMBitDepthChroma)
    //#if RATE_CONTROL_LAMBDA_DOMAIN
    fmt.Printf("RateControl                  : %d\n", TLibCommon.B2U(this.m_RCEnableRateControl))
    if this.m_RCEnableRateControl {
        fmt.Printf("TargetBitrate                : %d\n", this.m_RCTargetBitrate)
        fmt.Printf("KeepHierarchicalBit          : %d\n", this.m_RCKeepHierarchicalBit)
        fmt.Printf("LCULevelRC                   : %d\n", this.m_RCLCULevelRC)
        fmt.Printf("UseLCUSeparateModel          : %d\n", this.m_RCUseLCUSeparateModel)
        fmt.Printf("InitialQP                    : %d\n", this.m_RCInitialQP)
        fmt.Printf("ForceIntraQP                 : %d\n", this.m_RCForceIntraQP)
    }
    /*#else
      fmt.Printf("RateControl                  : %d\n", this.m_enableRateCtrl);
      if(this.m_enableRateCtrl)
      {
        fmt.Printf("TargetBitrate                : %d\n", this.m_targetBitrate);
        fmt.Printf("NumLCUInUnit                 : %d\n", this.m_numLCUInUnit);
      }
    #endif*/
    fmt.Printf("Max Num Merge Candidates     : %d\n", this.m_maxNumMergeCand)
    fmt.Printf("\n")

    fmt.Printf("TOOL CFG: ")
    fmt.Printf("IBD:%d ", TLibCommon.B2U(TLibCommon.G_bitDepthY > this.m_inputBitDepthY || TLibCommon.G_bitDepthC > this.m_inputBitDepthC))
    fmt.Printf("HAD:%d ", TLibCommon.B2U(this.m_bUseHADME))
    fmt.Printf("SRD:%d ", TLibCommon.B2U(this.m_bUseSBACRD))
    fmt.Printf("RDQ:%d ", TLibCommon.B2U(this.m_useRDOQ))
    fmt.Printf("RDQTS:%d ", TLibCommon.B2U(this.m_useRDOQTS))
    //#if L0232_RD_PENALTY
    fmt.Printf("RDpenalty:%d ", this.m_rdPenalty)
    //#endif
    fmt.Printf("SQP:%d ", this.m_uiDeltaQpRD)
    fmt.Printf("ASR:%d ", TLibCommon.B2U(this.m_bUseASR))
    fmt.Printf("LComb:%d ", TLibCommon.B2U(this.m_bUseLComb))
    fmt.Printf("FEN:%d ", TLibCommon.B2U(this.m_bUseFastEnc))
    fmt.Printf("ECU:%d ", TLibCommon.B2U(this.m_bUseEarlyCU))
    fmt.Printf("FDM:%d ", TLibCommon.B2U(this.m_useFastDecisionForMerge))
    fmt.Printf("CFM:%d ", TLibCommon.B2U(this.m_bUseCbfFastMode))
    fmt.Printf("ESD:%d ", TLibCommon.B2U(this.m_useEarlySkipDetection))
    fmt.Printf("RQT:%d ", 1)
    fmt.Printf("TransformSkip:%d ", TLibCommon.B2U(this.m_useTransformSkip))
    fmt.Printf("TransformSkipFast:%d ", TLibCommon.B2U(this.m_useTransformSkipFast))
    fmt.Printf("Slice: M=%d ", this.m_sliceMode)
    if this.m_sliceMode != 0 {
        fmt.Printf("A=%d ", this.m_sliceArgument)
    }
    fmt.Printf("SliceSegment: M=%d ", this.m_sliceSegmentMode)
    if this.m_sliceSegmentMode != 0 {
        fmt.Printf("A=%d ", this.m_sliceSegmentArgument)
    }
    fmt.Printf("CIP:%d ", TLibCommon.B2U(this.m_bUseConstrainedIntraPred))
    fmt.Printf("SAO:%d ", TLibCommon.B2U(this.m_bUseSAO))
    fmt.Printf("PCM:%d ", TLibCommon.B2U(this.m_usePCM && (1<<this.m_uiPCMLog2MinSize) <= this.m_uiMaxCUWidth))
    fmt.Printf("SAOLcuBasedOptimization:%d ", TLibCommon.B2U(this.m_saoLcuBasedOptimization))

    fmt.Printf("LosslessCuEnabled:%d ", TLibCommon.B2U(this.m_useLossless))
    fmt.Printf("WPP:%d ", TLibCommon.B2U(this.m_useWeightedPred))
    fmt.Printf("WPB:%d ", TLibCommon.B2U(this.m_useWeightedBiPred))
    fmt.Printf("PME:%d ", this.m_log2ParallelMergeLevel)
    fmt.Printf(" WaveFrontSynchro:%d WaveFrontSubstreams:%d", this.m_iWaveFrontSynchro, this.m_iWaveFrontSubstreams)
    fmt.Printf(" ScalingList:%d ", this.m_useScalingListId)
    fmt.Printf("TMVPMode:%d ", this.m_TMVPModeId)
    //#if ADAPTIVE_QP_SELECTION
    fmt.Printf("AQpS:%d", TLibCommon.B2U(this.m_bUseAdaptQpSelect))
    //#endif

    fmt.Printf(" SignBitHidingFlag:%d ", this.m_signHideFlag)
    fmt.Printf("RecalQP:%d", TLibCommon.B2U(this.m_recalculateQPAccordingToLambda))
    fmt.Printf("\n\n")

    //fflush(stdout);
}
