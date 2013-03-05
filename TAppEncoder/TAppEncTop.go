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
    "container/list"
    "fmt"
    "gohm/TLibCommon"
    "gohm/TLibEncoder"
    "io"
    "os"
)

// ====================================================================================================================
// Class definition
// ====================================================================================================================

/// encoder application class
type TAppEncTop struct {
    TAppEncCfg

    // class interface
    m_cTEncTop              *TLibEncoder.TEncTop    ///< encoder class
    m_cTVideoIOYuvInputFile *TLibCommon.TVideoIOYuv ///< input YUV file
    m_cTVideoIOYuvReconFile *TLibCommon.TVideoIOYuv ///< output reconstruction file

    m_cListPicYuvRec *list.List ///< list of reconstruction YUV files TComList<TComPicYuv*>

    m_iFrameRcvd int ///< number of received frames

    m_essentialBytes uint
    m_totalBytes     uint
}

func NewTAppEncTop() *TAppEncTop {
    return &TAppEncTop{m_cTEncTop: TLibEncoder.NewTEncTop(),
        m_cTVideoIOYuvInputFile: TLibCommon.NewTVideoIOYuv(),
        m_cTVideoIOYuvReconFile: TLibCommon.NewTVideoIOYuv(),
        m_cListPicYuvRec:        list.New(),
        m_iFrameRcvd:            0,
        m_totalBytes:            0,
        m_essentialBytes:        0}
}

func (this *TAppEncTop) Encode() (err error) { ///< main encoding function
    bitstreamFile, err := os.Create(this.m_pchBitstreamFile)
    if err != nil {
        fmt.Printf("\nfailed to open bitstream file `%s' for writing\n", this.m_pchBitstreamFile)
        return err
    }
    defer bitstreamFile.Close()

    pcPicYuvOrg := TLibCommon.NewTComPicYuv()
    //var pcPicYuvRec *TLibCommon.TComPicYuv;

    // initialize internal class & member variables
    this.xInitLibCfg()
    this.xCreateLib()
    this.xInitLib()

    // main encoder loop
    iNumEncoded := 0
    bEos := false

    outputAccessUnits := TLibEncoder.NewAccessUnits() ///< list of access units to write out.  is populated by the encoding process

    // allocate original YUV buffer
    pcPicYuvOrg.Create(this.m_iSourceWidth, this.m_iSourceHeight, this.m_uiMaxCUWidth, this.m_uiMaxCUHeight, this.m_uiMaxCUDepth)

    for !bEos {
        // get buffers
        this.xGetBuffer() //pcPicYuvRec =

        // read input YUV file
        this.m_cTVideoIOYuvInputFile.Read(pcPicYuvOrg, this.m_aiPad[:])

        // increase number of received frames
        this.m_iFrameRcvd++

        bEos = (this.m_iFrameRcvd == this.m_framesToBeEncoded)

        flush := false
        // if end of file (which is only detected on a read failure) flush the encoder of any queued pictures
        if this.m_cTVideoIOYuvInputFile.IsEof() {
            flush = true
            bEos = true
            this.m_iFrameRcvd--
            this.m_cTEncTop.GetEncCfg().SetFramesToBeEncoded(this.m_iFrameRcvd)
        }

        // call encoding function for one frame
        if flush {
            this.m_cTEncTop.Encode(bEos, nil, this.m_cListPicYuvRec, outputAccessUnits, &iNumEncoded)
        } else {
            this.m_cTEncTop.Encode(bEos, pcPicYuvOrg, this.m_cListPicYuvRec, outputAccessUnits, &iNumEncoded)
        }
        // write bistream to file if necessary
        if iNumEncoded > 0 {
            this.xWriteOutput(bitstreamFile, iNumEncoded, outputAccessUnits)
            outputAccessUnits.Init()
        }
    }

    this.m_cTEncTop.PrintSummary()

    // delete original YUV buffer
    pcPicYuvOrg.Destroy()
    //delete pcPicYuvOrg;
    pcPicYuvOrg = nil

    // delete used buffers in encoder class
    this.m_cTEncTop.DeletePicBuffer()

    // delete buffers & classes
    this.xDeleteBuffer()
    this.xDestroyLib()

    this.printRateSummary()

    return nil
}

func (this *TAppEncTop) GetTEncTop() *TLibEncoder.TEncTop {
    return this.m_cTEncTop
}   ///< return encoder class pointer reference

//protected:
// initialization
func (this *TAppEncTop) xCreateLib() { ///< create files & encoder class
    // Video I/O
    this.m_cTVideoIOYuvInputFile.Open(this.m_pchInputFile, false, this.m_inputBitDepthY, this.m_inputBitDepthC, this.m_internalBitDepthY, this.m_internalBitDepthC) // read  mode
    this.m_cTVideoIOYuvInputFile.SkipFrames(this.m_FrameSkip, uint(this.m_iSourceWidth-this.m_aiPad[0]), uint(this.m_iSourceHeight-this.m_aiPad[1]))

    if this.m_pchReconFile != "" {
        this.m_cTVideoIOYuvReconFile.Open(this.m_pchReconFile, true, this.m_outputBitDepthY, this.m_outputBitDepthC, this.m_internalBitDepthY, this.m_internalBitDepthC) // write mode
    }
    // Neo Decoder
    this.m_cTEncTop.Create(this.m_pchTraceFile)
}

func (this *TAppEncTop) xInitLibCfg() { ///< initialize internal variables
    vps := TLibCommon.NewTComVPS()

    vps.SetMaxTLayers(uint(this.m_maxTempLayer))
    if this.m_maxTempLayer == 1 {
        vps.SetTemporalNestingFlag(true)
    }
    vps.SetMaxLayers(1)
    for i := 0; i < TLibCommon.MAX_TLAYER; i++ {
        vps.SetNumReorderPics(uint(this.m_numReorderPics[i]), uint(i))
        vps.SetMaxDecPicBuffering(uint(this.m_maxDecPicBuffering[i]), uint(i))
    }

    pcEncCfg := TLibEncoder.NewTEncCfg()

    pcEncCfg.SetVPS(vps)

    pcEncCfg.SetProfile(TLibCommon.PROFILE(this.m_profile))
    pcEncCfg.SetLevel(TLibCommon.TIER(this.m_levelTier), TLibCommon.LEVEL(this.m_level))
    //#if L0046_CONSTRAINT_FLAGS
    pcEncCfg.SetProgressiveSourceFlag(this.m_progressiveSourceFlag)
    pcEncCfg.SetInterlacedSourceFlag(this.m_interlacedSourceFlag)
    pcEncCfg.SetNonPackedConstraintFlag(this.m_nonPackedConstraintFlag)
    pcEncCfg.SetFrameOnlyConstraintFlag(this.m_frameOnlyConstraintFlag)
    //#endif

    pcEncCfg.SetFrameRate(this.m_iFrameRate)
    pcEncCfg.SetFrameSkip(this.m_FrameSkip)
    pcEncCfg.SetSourceWidth(this.m_iSourceWidth)
    pcEncCfg.SetSourceHeight(this.m_iSourceHeight)
    pcEncCfg.SetConformanceWindow(this.m_confLeft, this.m_confRight, this.m_confTop, this.m_confBottom)
    pcEncCfg.SetFramesToBeEncoded(this.m_framesToBeEncoded)

	  // coding unit (CU) definition
	pcEncCfg.SetMaxCUWidth                   ( this.m_uiMaxCUWidth);
	pcEncCfg.SetMaxCUHeight                  ( this.m_uiMaxCUHeight);
	pcEncCfg.SetMaxCUDepth                   ( this.m_uiMaxCUDepth);
	pcEncCfg.SetAddCUDepth                   ( this.m_uiAddCUDepth);
  
    //====== Coding Structure ========
    pcEncCfg.SetIntraPeriod(this.m_iIntraPeriod)
    pcEncCfg.SetDecodingRefreshType(this.m_iDecodingRefreshType)
    pcEncCfg.SetGOPSize(this.m_iGOPSize)
    pcEncCfg.SetGopList(this.m_GOPList[:])
    pcEncCfg.SetExtraRPSs(this.m_extraRPSs)
    for i := 0; i < TLibCommon.MAX_TLAYER; i++ {
        pcEncCfg.SetNumReorderPics(this.m_numReorderPics[i], uint(i))
        pcEncCfg.SetMaxDecPicBuffering(uint(this.m_maxDecPicBuffering[i]), uint(i))
    }
    for uiLoop := 0; uiLoop < TLibCommon.MAX_TLAYER; uiLoop++ {
        pcEncCfg.SetLambdaModifier(uint(uiLoop), this.m_adLambdaModifier[uiLoop])
    }
    pcEncCfg.SetQP(this.m_iQP)

    pcEncCfg.SetPad(this.m_aiPad[:])

    pcEncCfg.SetMaxTempLayer(this.m_maxTempLayer)
    pcEncCfg.SetUseAMP(this.m_enableAMP)

    //===== Slice ========

    //====== Loop/Deblock Filter ========
    pcEncCfg.SetLoopFilterDisable(this.m_bLoopFilterDisable)
    pcEncCfg.SetLoopFilterOffsetInPPS(this.m_loopFilterOffsetInPPS)
    pcEncCfg.SetLoopFilterBetaOffset(this.m_loopFilterBetaOffsetDiv2)
    pcEncCfg.SetLoopFilterTcOffset(this.m_loopFilterTcOffsetDiv2)
    pcEncCfg.SetDeblockingFilterControlPresent(this.m_DeblockingFilterControlPresent)

    //====== Motion search ========
    pcEncCfg.SetFastSearch(this.m_iFastSearch)
    pcEncCfg.SetSearchRange(this.m_iSearchRange)
    pcEncCfg.SetBipredSearchRange(this.m_bipredSearchRange)

    //====== Quality control ========
    pcEncCfg.SetMaxDeltaQP(this.m_iMaxDeltaQP)
    pcEncCfg.SetMaxCuDQPDepth(this.m_iMaxCuDQPDepth)

    pcEncCfg.SetChromaCbQpOffset(this.m_cbQpOffset)
    pcEncCfg.SetChromaCrQpOffset(this.m_crQpOffset)

    //#if ADAPTIVE_QP_SELECTION
    pcEncCfg.SetUseAdaptQpSelect(this.m_bUseAdaptQpSelect)
    //#endif

    var lowestQP int
    lowestQP = -6 * (TLibCommon.G_bitDepthY - 8) // XXX: check

    if (this.m_iMaxDeltaQP == 0) && (this.m_iQP == lowestQP) && (this.m_useLossless == true) {
        this.m_bUseAdaptiveQP = false
    }
    pcEncCfg.SetUseAdaptiveQP(this.m_bUseAdaptiveQP)
    pcEncCfg.SetQPAdaptationRange(this.m_iQPAdaptationRange)

    //====== Tool list ========
    pcEncCfg.SetUseSBACRD(this.m_bUseSBACRD)
    pcEncCfg.SetDeltaQpRD(this.m_uiDeltaQpRD)
    pcEncCfg.SetUseASR(this.m_bUseASR)
    pcEncCfg.SetUseHADME(this.m_bUseHADME)
    pcEncCfg.SetUseLossless(this.m_useLossless)
    pcEncCfg.SetUseLComb(this.m_bUseLComb)
    pcEncCfg.SetdQPs(this.m_aidQP)
    pcEncCfg.SetUseRDOQ(this.m_useRDOQ)
    pcEncCfg.SetUseRDOQTS(this.m_useRDOQTS)
    //#if L0232_RD_PENALTY
    pcEncCfg.SetRDpenalty(this.m_rdPenalty)
    //#endif
    pcEncCfg.SetQuadtreeTULog2MaxSize(this.m_uiQuadtreeTULog2MaxSize)
    pcEncCfg.SetQuadtreeTULog2MinSize(this.m_uiQuadtreeTULog2MinSize)
    pcEncCfg.SetQuadtreeTUMaxDepthInter(this.m_uiQuadtreeTUMaxDepthInter)
    pcEncCfg.SetQuadtreeTUMaxDepthIntra(this.m_uiQuadtreeTUMaxDepthIntra)
    pcEncCfg.SetUseFastEnc(this.m_bUseFastEnc)
    pcEncCfg.SetUseEarlyCU(this.m_bUseEarlyCU)
    pcEncCfg.SetUseFastDecisionForMerge(this.m_useFastDecisionForMerge)
    pcEncCfg.SetUseCbfFastMode(this.m_bUseCbfFastMode)
    pcEncCfg.SetUseEarlySkipDetection(this.m_useEarlySkipDetection)

    pcEncCfg.SetUseTransformSkip(this.m_useTransformSkip)
    pcEncCfg.SetUseTransformSkipFast(this.m_useTransformSkipFast)
    pcEncCfg.SetUseConstrainedIntraPred(this.m_bUseConstrainedIntraPred)
    pcEncCfg.SetPCMLog2MinSize(this.m_uiPCMLog2MinSize)
    pcEncCfg.SetUsePCM(this.m_usePCM)
    pcEncCfg.SetPCMLog2MaxSize(this.m_pcmLog2MaxSize)
    pcEncCfg.SetMaxNumMergeCand(this.m_maxNumMergeCand)

    //====== Weighted Prediction ========
    pcEncCfg.SetUseWP(this.m_useWeightedPred)
    pcEncCfg.SetWPBiPred(this.m_useWeightedBiPred)
    //====== Parallel Merge Estimation ========
    pcEncCfg.SetLog2ParallelMergeLevelMinus2(this.m_log2ParallelMergeLevel - 2)

    //====== Slice ========
    pcEncCfg.SetSliceMode(this.m_sliceMode)
    pcEncCfg.SetSliceArgument(this.m_sliceArgument)

    //====== Dependent Slice ========
    pcEncCfg.SetSliceSegmentMode(this.m_sliceSegmentMode)
    pcEncCfg.SetSliceSegmentArgument(this.m_sliceSegmentArgument)

    iNumPartInCU := 1 << (this.m_uiMaxCUDepth << 1)
    if this.m_sliceSegmentMode == TLibCommon.FIXED_NUMBER_OF_LCU {
        pcEncCfg.SetSliceSegmentArgument(this.m_sliceSegmentArgument * iNumPartInCU)
    }
    if this.m_sliceMode == TLibCommon.FIXED_NUMBER_OF_LCU {
        pcEncCfg.SetSliceArgument(this.m_sliceArgument * iNumPartInCU)
    }
    if this.m_sliceMode == TLibCommon.FIXED_NUMBER_OF_TILES {
        pcEncCfg.SetSliceArgument(this.m_sliceArgument)
    }

    if this.m_sliceMode == 0 {
        this.m_bLFCrossSliceBoundaryFlag = true
    }
    pcEncCfg.SetLFCrossSliceBoundaryFlag(this.m_bLFCrossSliceBoundaryFlag)
    pcEncCfg.SetUseSAO(this.m_bUseSAO)
    pcEncCfg.SetMaxNumOffsetsPerPic(this.m_maxNumOffsetsPerPic)

    pcEncCfg.SetSaoLcuBoundary(this.m_saoLcuBoundary)
    pcEncCfg.SetSaoLcuBasedOptimization(this.m_saoLcuBasedOptimization)
    pcEncCfg.SetPCMInputBitDepthFlag(this.m_bPCMInputBitDepthFlag)
    pcEncCfg.SetPCMFilterDisableFlag(this.m_bPCMFilterDisableFlag)

    pcEncCfg.SetDecodedPictureHashSEIEnabled(this.m_decodedPictureHashSEIEnabled)
    pcEncCfg.SetRecoveryPointSEIEnabled(this.m_recoveryPointSEIEnabled)
    pcEncCfg.SetBufferingPeriodSEIEnabled(this.m_bufferingPeriodSEIEnabled)
    pcEncCfg.SetPictureTimingSEIEnabled(this.m_pictureTimingSEIEnabled)
    pcEncCfg.SetFramePackingArrangementSEIEnabled(this.m_framePackingSEIEnabled)
    pcEncCfg.SetFramePackingArrangementSEIType(this.m_framePackingSEIType)
    pcEncCfg.SetFramePackingArrangementSEIId(this.m_framePackingSEIId)
    pcEncCfg.SetFramePackingArrangementSEIQuincunx(this.m_framePackingSEIQuincunx)
    pcEncCfg.SetFramePackingArrangementSEIInterpretation(this.m_framePackingSEIInterpretation)
    pcEncCfg.SetDisplayOrientationSEIAngle(this.m_displayOrientationSEIAngle)
    pcEncCfg.SetTemporalLevel0IndexSEIEnabled(this.m_temporalLevel0IndexSEIEnabled)
    pcEncCfg.SetUniformSpacingIdr(this.m_iUniformSpacingIdr)
    pcEncCfg.SetNumColumnsMinus1(this.m_iNumColumnsMinus1)
    pcEncCfg.SetNumRowsMinus1(this.m_iNumRowsMinus1)
    if this.m_iUniformSpacingIdr == 0 {
        pcEncCfg.SetColumnWidth(this.m_pColumnWidth)
        pcEncCfg.SetRowHeight(this.m_pRowHeight)
    }
    pcEncCfg.XCheckGSParameters()
    uiTilesCount := (this.m_iNumRowsMinus1 + 1) * (this.m_iNumColumnsMinus1 + 1)
    if uiTilesCount == 1 {
        this.m_bLFCrossTileBoundaryFlag = true
    }
    pcEncCfg.SetLFCrossTileBoundaryFlag(this.m_bLFCrossTileBoundaryFlag)
    pcEncCfg.SetWaveFrontSynchro(this.m_iWaveFrontSynchro)
    pcEncCfg.SetWaveFrontSubstreams(this.m_iWaveFrontSubstreams)
    pcEncCfg.SetTMVPModeId(this.m_TMVPModeId)
    pcEncCfg.SetUseScalingListId(this.m_useScalingListId)
    pcEncCfg.SetScalingListFile(this.m_scalingListFile)
    pcEncCfg.SetSignHideFlag(this.m_signHideFlag)
    //#if RATE_CONTROL_LAMBDA_DOMAIN
    pcEncCfg.SetUseRateCtrl(this.m_RCEnableRateControl)
    pcEncCfg.SetTargetBitrate(this.m_RCTargetBitrate)
    pcEncCfg.SetKeepHierBit(this.m_RCKeepHierarchicalBit)
    pcEncCfg.SetLCULevelRC(this.m_RCLCULevelRC)
    pcEncCfg.SetUseLCUSeparateModel(this.m_RCUseLCUSeparateModel)
    pcEncCfg.SetInitialQP(this.m_RCInitialQP)
    pcEncCfg.SetForceIntraQP(this.m_RCForceIntraQP)
    //#else
    //  pcEncCfg.SetUseRateCtrl     ( this.m_enableRateCtrl);
    //  pcEncCfg.SetTargetBitrate   ( this.m_targetBitrate);
    //  pcEncCfg.SetNumLCUInUnit    ( this.m_numLCUInUnit);
    //#endif
    pcEncCfg.SetTransquantBypassEnableFlag(this.m_TransquantBypassEnableFlag)
    pcEncCfg.SetCUTransquantBypassFlagValue(this.m_CUTransquantBypassFlagValue)
    pcEncCfg.SetUseRecalculateQPAccordingToLambda(this.m_recalculateQPAccordingToLambda)
    pcEncCfg.SetUseStrongIntraSmoothing(this.m_useStrongIntraSmoothing)
    pcEncCfg.SetActiveParameterSetsSEIEnabled(this.m_activeParameterSetsSEIEnabled)
    pcEncCfg.SetVuiParametersPresentFlag(this.m_vuiParametersPresentFlag)
    pcEncCfg.SetAspectRatioIdc(this.m_aspectRatioIdc)
    pcEncCfg.SetSarWidth(this.m_sarWidth)
    pcEncCfg.SetSarHeight(this.m_sarHeight)
    pcEncCfg.SetOverscanInfoPresentFlag(this.m_overscanInfoPresentFlag)
    pcEncCfg.SetOverscanAppropriateFlag(this.m_overscanAppropriateFlag)
    pcEncCfg.SetVideoSignalTypePresentFlag(this.m_videoSignalTypePresentFlag)
    pcEncCfg.SetVideoFormat(this.m_videoFormat)
    pcEncCfg.SetVideoFullRangeFlag(this.m_videoFullRangeFlag)
    pcEncCfg.SetColourDescriptionPresentFlag(this.m_colourDescriptionPresentFlag)
    pcEncCfg.SetColourPrimaries(this.m_colourPrimaries)
    pcEncCfg.SetTransferCharacteristics(this.m_transferCharacteristics)
    pcEncCfg.SetMatrixCoefficients(this.m_matrixCoefficients)
    pcEncCfg.SetChromaLocInfoPresentFlag(this.m_chromaLocInfoPresentFlag)
    pcEncCfg.SetChromaSampleLocTypeTopField(this.m_chromaSampleLocTypeTopField)
    pcEncCfg.SetChromaSampleLocTypeBottomField(this.m_chromaSampleLocTypeBottomField)
    pcEncCfg.SetNeutralChromaIndicationFlag(this.m_neutralChromaIndicationFlag)
    pcEncCfg.SetDefaultDisplayWindow(this.m_defDispWinLeftOffset, this.m_defDispWinRightOffset, this.m_defDispWinTopOffset, this.m_defDispWinBottomOffset)
    pcEncCfg.SetFrameFieldInfoPresentFlag(this.m_frameFieldInfoPresentFlag)
    pcEncCfg.SetPocProportionalToTimingFlag(this.m_pocProportionalToTimingFlag)
    pcEncCfg.SetNumTicksPocDiffOneMinus1(this.m_numTicksPocDiffOneMinus1)
    pcEncCfg.SetBitstreamRestrictionFlag(this.m_bitstreamRestrictionFlag)
    pcEncCfg.SetTilesFixedStructureFlag(this.m_tilesFixedStructureFlag)
    pcEncCfg.SetMotionVectorsOverPicBoundariesFlag(this.m_motionVectorsOverPicBoundariesFlag)
    pcEncCfg.SetMinSpatialSegmentationIdc(this.m_minSpatialSegmentationIdc)
    pcEncCfg.SetMaxBytesPerPicDenom(this.m_maxBytesPerPicDenom)
    pcEncCfg.SetMaxBitsPerMinCuDenom(this.m_maxBitsPerMinCuDenom)
    pcEncCfg.SetLog2MaxMvLengthHorizontal(this.m_log2MaxMvLengthHorizontal)
    pcEncCfg.SetLog2MaxMvLengthVertical(this.m_log2MaxMvLengthVertical)

    this.m_cTEncTop.SetEncCfg(pcEncCfg)
}

func (this *TAppEncTop) xInitLib() { ///< initialize encoder class
    this.m_cTEncTop.Init()
}

func (this *TAppEncTop) xDestroyLib() { ///< destroy encoder class
    // Video I/O
    this.m_cTVideoIOYuvInputFile.Close()
    this.m_cTVideoIOYuvReconFile.Close()

    // Neo Decoder
    this.m_cTEncTop.Destroy()
}

/// obtain required buffers
func (this *TAppEncTop) xGetBuffer() *TLibCommon.TComPicYuv {
    //assert( this.m_iGOPSize > 0 );
    var rpcPicYuvRec *TLibCommon.TComPicYuv

    // org. buffer
    if this.m_cListPicYuvRec.Len() == this.m_iGOPSize {
        e := this.m_cListPicYuvRec.Front()
        rpcPicYuvRec = e.Value.(*TLibCommon.TComPicYuv)
        this.m_cListPicYuvRec.Remove(e)
    } else {
        rpcPicYuvRec = TLibCommon.NewTComPicYuv()

        rpcPicYuvRec.Create(this.m_iSourceWidth, this.m_iSourceHeight, this.m_uiMaxCUWidth, this.m_uiMaxCUHeight, this.m_uiMaxCUDepth)
    }
    this.m_cListPicYuvRec.PushBack(rpcPicYuvRec)

    return rpcPicYuvRec
}

/// delete allocated buffers
func (this *TAppEncTop) xDeleteBuffer() {

    //iSize = Int( this.m_cListPicYuvRec.size() );

    for iterPicYuvRec := this.m_cListPicYuvRec.Front(); iterPicYuvRec != nil; iterPicYuvRec = iterPicYuvRec.Next() {
        pcPicYuvRec := iterPicYuvRec.Value.(*TLibCommon.TComPicYuv)
        pcPicYuvRec.Destroy()
        //delete pcPicYuvRec; pcPicYuvRec = NULL;
    }

    this.m_cListPicYuvRec.Init()
}

// file I/O
func (this *TAppEncTop) xWriteOutput(bitstreamFile io.Writer, iNumEncoded int, accessUnits *TLibEncoder.AccessUnits) { //const std::list<AccessUnit>& ///< write bitstream to file
    var i int

    iterPicYuvRec := this.m_cListPicYuvRec.Back() //TComList<TComPicYuv*>::iterator
    iterBitstream := accessUnits.Front()          //list<AccessUnit>::const_iterator

    for i = 0; i < iNumEncoded-1; i++ {
        iterPicYuvRec = iterPicYuvRec.Prev()
    }

    for i = 0; i < iNumEncoded; i++ {
        pcPicYuvRec := iterPicYuvRec.Value.(*TLibCommon.TComPicYuv)
        if this.m_pchReconFile != "" {
            this.m_cTVideoIOYuvReconFile.Write(pcPicYuvRec, this.m_confLeft, this.m_confRight, this.m_confTop, this.m_confBottom)
        }

        au := iterBitstream.Value.(*TLibEncoder.AccessUnit)
        stats := TLibEncoder.WriteAnnexB(bitstreamFile, au) //const vector<UInt>&
        this.rateStatsAccum(au, stats)

        iterPicYuvRec = iterPicYuvRec.Next()
        iterBitstream = iterBitstream.Next()
    }
}

func (this *TAppEncTop) rateStatsAccum(au *TLibEncoder.AccessUnit, annexBsizes *list.List) { //const std::vector<UInt>
    it_stats := annexBsizes.Front()
    for it_au := au.Front(); it_au != nil; it_au = it_au.Next() {
        nalu := it_au.Value.(*TLibEncoder.NALUnitEBSP)
        stats := it_stats.Value.(uint)
        switch nalu.GetNalUnitType() {
        case TLibCommon.NAL_UNIT_CODED_SLICE_TRAIL_R:
            fallthrough
        case TLibCommon.NAL_UNIT_CODED_SLICE_TRAIL_N:
            fallthrough
        case TLibCommon.NAL_UNIT_CODED_SLICE_TLA:
            fallthrough
        case TLibCommon.NAL_UNIT_CODED_SLICE_TSA_N:
            fallthrough
        case TLibCommon.NAL_UNIT_CODED_SLICE_STSA_R:
            fallthrough
        case TLibCommon.NAL_UNIT_CODED_SLICE_STSA_N:
            fallthrough
        case TLibCommon.NAL_UNIT_CODED_SLICE_BLA:
            fallthrough
        case TLibCommon.NAL_UNIT_CODED_SLICE_BLANT:
            fallthrough
        case TLibCommon.NAL_UNIT_CODED_SLICE_BLA_N_LP:
            fallthrough
        case TLibCommon.NAL_UNIT_CODED_SLICE_IDR:
            fallthrough
        case TLibCommon.NAL_UNIT_CODED_SLICE_IDR_N_LP:
            fallthrough
        case TLibCommon.NAL_UNIT_CODED_SLICE_CRA:
            fallthrough
        case TLibCommon.NAL_UNIT_CODED_SLICE_RADL_N:
            fallthrough
        case TLibCommon.NAL_UNIT_CODED_SLICE_DLP:
            fallthrough
        case TLibCommon.NAL_UNIT_CODED_SLICE_RASL_N:
            fallthrough
        case TLibCommon.NAL_UNIT_CODED_SLICE_TFD:
            fallthrough
        case TLibCommon.NAL_UNIT_VPS:
            fallthrough
        case TLibCommon.NAL_UNIT_SPS:
            fallthrough
        case TLibCommon.NAL_UNIT_PPS:
            this.m_essentialBytes += stats
        default:
        }

        this.m_totalBytes += stats

        it_stats = it_stats.Next()
    }
}

func (this *TAppEncTop) printRateSummary() {
    time := float64(this.m_iFrameRcvd) / float64(this.m_iFrameRate)
    fmt.Printf("Bytes written to file: %d (%.3f kbps)\n", this.m_totalBytes, 0.008*float64(this.m_totalBytes)/time)
    //#if VERBOSE_RATE
    fmt.Printf("Bytes for SPS/PPS/Slice (Incl. Annex B): %d (%.3f kbps)\n", this.m_essentialBytes, 0.008*float64(this.m_essentialBytes)/time)
    //#endif
}
