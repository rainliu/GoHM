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
	//"time"
	"math"
	"container/list"
	"gohm/TLibCommon"
)

// ====================================================================================================================
// Enumeration
// ====================================================================================================================
type PROCESSING_STATE uint8
const (
  EXECUTE_INLOOPFILTER = iota
  ENCODE_SLICE
)

type SCALING_LIST_PARAMETER uint8
const (
  SCALING_LIST_OFF	= iota
  SCALING_LIST_DEFAULT
  SCALING_LIST_FILE_READ
)

// ====================================================================================================================
// Class definition
// ====================================================================================================================

/// GOP encoder class
type TEncGOP struct{
  //  Data
  m_bLongtermTestPictureHasBeenCoded	bool;
  m_bLongtermTestPictureHasBeenCoded2	bool;
  m_numLongTermRefPicSPS	uint;
  m_ltRefPicPocLsbSps	[33]uint;
  m_ltRefPicUsedByCurrPicFlag	[33]uint;
  m_iLastIDR	int;
  m_iGopSize	int;
  m_iNumPicCoded	int;
  m_bFirst	bool;
  
  //  Access channel
  m_pcEncTop	*TEncTop;
  m_pcCfg		*TEncCfg;
  m_pcSliceEncoder	*TEncSlice;
  m_pcListPic	*list.List;
  
  m_pcEntropyCoder	*TEncEntropy;
  m_pcCavlcCoder	*TEncCavlc;
  m_pcSbacCoder		*TEncSbac;
  m_pcBinCABAC		*TEncBinCABAC;
  m_pcLoopFilter	*TLibCommon.TComLoopFilter;

  // m_seiWriter		SEIWriter;
  
  //--Adaptive Loop filter
  m_pcSAO	*TEncSampleAdaptiveOffset;
  m_pcBitCounter	*TLibCommon.TComBitCounter;
  m_pcRateCtrl	*TEncRateCtrl;
  // indicate sequence first
  m_bSeqFirst	bool;
  
  // clean decoding refresh
  m_bRefreshPending	bool;
  m_pocCRA	int;
  m_storedStartCUAddrForEncodingSlice	map[int]int;  //vector<int>
  m_storedStartCUAddrForEncodingDependentSlice	map[int]int;

  m_vRVM_RP	map[int]int;
  m_lastBPSEI	uint;
  m_totalCoded	uint;
  m_cpbRemovalDelay	uint;
//#if SEI_TEMPORAL_LEVEL0_INDEX
  m_tl0Idx	uint;
  m_rapIdx	uint;
//#endif
 
  m_gcAnalyzeAll	TEncAnalyze;
  m_gcAnalyzeI	TEncAnalyze;
  m_gcAnalyzeP	TEncAnalyze;
  m_gcAnalyzeB	TEncAnalyze;
}

func NewTEncGOP() *TEncGOP{
	return &TEncGOP{m_bFirst:true, m_bSeqFirst:true}
}
 
func (this *TEncGOP)  create      (  iWidth,  iHeight int,  uiMaxCUWidth,  uiMaxCUHeight uint){}
func (this *TEncGOP)  destroy     () {}
  
func (this *TEncGOP)  init        ( pcTEncTop *TEncTop){
  //fmt.Printf("not init yet in TEncGop\n");
  
  this.m_pcEncTop     		  = pcTEncTop;
  this.m_pcCfg                = pcTEncTop.GetEncCfg();
  this.m_pcSliceEncoder       = pcTEncTop.getSliceEncoder();
  this.m_pcListPic            = pcTEncTop.getListPic();
  
  this.m_pcEntropyCoder       = pcTEncTop.getEntropyCoder();
  this.m_pcCavlcCoder         = pcTEncTop.getCavlcCoder();
  this.m_pcSbacCoder          = pcTEncTop.getSbacCoder();
  this.m_pcBinCABAC           = pcTEncTop.getBinCABAC();
  this.m_pcLoopFilter         = pcTEncTop.getLoopFilter();
  this.m_pcBitCounter         = pcTEncTop.getBitCounter();
  
  //--Adaptive Loop filter
  this.m_pcSAO                = pcTEncTop.getSAO();
  this.m_pcRateCtrl           = pcTEncTop.getRateCtrl();
  
  this.m_lastBPSEI          = 0;
  this.m_totalCoded         = 0;
}

func (this *TEncGOP)  compressGOP ( iPOCLast, iNumPicRcvd int, rcListPic, rcListPicYuvRecOut *list.List, accessUnitsInGOP *list.List ){
	fmt.Printf("not compressGOP yet in TEncGop\n");
/*
  var        pcPic *TLibCommon.TComPic;
  var    pcPicYuvRecOut	*TLibCommon.TComPicYuv;
  var      pcSlice	*TLibCommon.TComSlice;
  var 	pcBitstreamRedirect	*TLibCommon.TComOutputBitstream;
  pcBitstreamRedirect = TLibCommon.NewTComOutputBitstream();
  //AccessUnit::iterator  itLocationToPushSliceHeaderNALU; // used to store location where NALU containing slice header is to be inserted
                   uiOneBitstreamPerSliceLength := uint(0);
  var pcSbacCoders *TEncSbac;
  var pcSubstreamsOut *TLibCommon.TComOutputBitstream;

  this.xInitGOP( iPOCLast, iNumPicRcvd, rcListPic, rcListPicYuvRecOut );

  this.m_iNumPicCoded = 0;
  //var pictureTimingSEI SEIPictureTiming;
  var accumBitsDU,accumNalsDU  *uint;
  for iGOPid:=0; iGOPid < this.m_iGopSize; iGOPid++ {
    uiColDir := uint(1);
    //-- For time output for each slice
    iBeforeTime := time.Now()

    //select uiColDir
    iCloseLeft:=1
    iCloseRight:=-1;
    for i := 0; i<this.m_pcCfg.GetGOPEntry(iGOPid).m_numRefPics; i++ {
      iRef := this.m_pcCfg.GetGOPEntry(iGOPid).m_referencePics[i];
      if iRef>0&&(iRef<iCloseRight||iCloseRight==-1) {
        iCloseRight=iRef;
      }else if iRef<0&&(iRef>iCloseLeft||iCloseLeft==1) {
        iCloseLeft=iRef;
      }
    }
    if iCloseRight>-1 {
      iCloseRight=iCloseRight+this.m_pcCfg.GetGOPEntry(iGOPid).m_POC-1;
    }
    if iCloseLeft<1 {
      iCloseLeft=iCloseLeft+this.m_pcCfg.GetGOPEntry(iGOPid).m_POC-1;
      for iCloseLeft<0 {
        iCloseLeft+=this.m_iGopSize;
      }
    }
    iLeftQP:=0;
    iRightQP:=0;
    for i:=0; i<this.m_iGopSize; i++ {
      if this.m_pcCfg.GetGOPEntry(i).m_POC==(iCloseLeft%this.m_iGopSize)+1 {
        iLeftQP= this.m_pcCfg.GetGOPEntry(i).m_QPOffset;
      }
      if this.m_pcCfg.GetGOPEntry(i).m_POC==(iCloseRight%this.m_iGopSize)+1 {
        iRightQP=this.m_pcCfg.GetGOPEntry(i).m_QPOffset;
      }
    }
    if iCloseRight>-1&&iRightQP<iLeftQP {
      uiColDir=0;
    }

    /////////////////////////////////////////////////////////////////////////////////////////////////// Initial to start encoding
    pocCurr := iPOCLast -iNumPicRcvd+ this.m_pcCfg.GetGOPEntry(iGOPid).m_POC;
    iTimeOffset := this.m_pcCfg.GetGOPEntry(iGOPid).m_POC;
    if iPOCLast == 0 {
      pocCurr=0;
      iTimeOffset = 1;
    }
    if pocCurr>=this.m_pcCfg.GetFrameToBeEncoded() {
      continue;
    }

    if this.getNalUnitType(pocCurr) == TLibCommon.NAL_UNIT_CODED_SLICE_IDR || this.getNalUnitType(pocCurr) == TLibCommon.NAL_UNIT_CODED_SLICE_IDR_N_LP {
      this.m_iLastIDR = pocCurr;
    }        
    // start a new access unit: create an entry in the list of output access units
    accessUnitsInGOP.PushBack(AccessUnit());
    accessUnit := accessUnitsInGOP.Back();
    this.xGetBuffer( rcListPic, rcListPicYuvRecOut, iNumPicRcvd, iTimeOffset, pcPic, pcPicYuvRecOut, pocCurr );

    //  Slice data initialization
    pcPic.clearSliceBuffer();
    //assert(pcPic.GetNumAllocatedSlice() == 1);
    this.m_pcSliceEncoder.setSliceIdx(0);
    pcPic.SetCurrSliceIdx(0);

    this.m_pcSliceEncoder.initEncSlice ( pcPic, iPOCLast, pocCurr, iNumPicRcvd, iGOPid, pcSlice, this.m_pcEncTop.getSPS(), this.m_pcEncTop.getPPS() );
    pcSlice.setLastIDR(this.m_iLastIDR);
    pcSlice.setSliceIdx(0);
    //set default slice level flag to the same as SPS level flag
    pcSlice.setLFCrossSliceBoundaryFlag(  pcSlice.GetPPS().getLoopFilterAcrossSlicesEnabledFlag()  );
    pcSlice.setScalingList ( this.m_pcEncTop.getScalingList()  );
    pcSlice.GetScalingList().setUseTransformSkip(this.m_pcEncTop.getPPS().getUseTransformSkip());
    if this.m_pcEncTop.getUseScalingListId() == SCALING_LIST_OFF {
      this.m_pcEncTop.getTrQuant().setFlatScalingList();
      this.m_pcEncTop.getTrQuant().setUseScalingList(false);
      this.m_pcEncTop.getSPS().setScalingListPresentFlag(false);
      this.m_pcEncTop.getPPS().setScalingListPresentFlag(false);
    }else if this.m_pcEncTop.getUseScalingListId() == SCALING_LIST_DEFAULT {
      pcSlice.setDefaultScalingList ();
      this.m_pcEncTop.getSPS().setScalingListPresentFlag(false);
      this.m_pcEncTop.getPPS().setScalingListPresentFlag(false);
      this.m_pcEncTop.getTrQuant().setScalingList(pcSlice.GetScalingList());
      this.m_pcEncTop.getTrQuant().setUseScalingList(true);
    }else if this.m_pcEncTop.getUseScalingListId() == SCALING_LIST_FILE_READ {
      if pcSlice.GetScalingList().xParseScalingList(this.m_pcCfg.GetScalingListFile()) {
        pcSlice.setDefaultScalingList ();
      }
      pcSlice.GetScalingList().checkDcOfMatrix();
      this.m_pcEncTop.getSPS().setScalingListPresentFlag(pcSlice.checkDefaultScalingList());
      this.m_pcEncTop.getPPS().setScalingListPresentFlag(false);
      this.m_pcEncTop.getTrQuant().setScalingList(pcSlice.GetScalingList());
      this.m_pcEncTop.getTrQuant().setUseScalingList(true);
    }else{
      fmt.Printf("error : ScalingList == %d no support\n",this.m_pcEncTop.getUseScalingListId());
      //assert(0);
    }

    if pcSlice.GetSliceType()==TLibCommon.B_SLICE&&this.m_pcCfg.GetGOPEntry(iGOPid).m_sliceType=="P"{
      pcSlice.setSliceType(TLibCommon.P_SLICE);
    }
    // Set the nal unit type
    pcSlice.setNalUnitType(getNalUnitType(pocCurr));
    if pcSlice.GetNalUnitType()==TLibCommon.NAL_UNIT_CODED_SLICE_TRAIL_R {
      if pcSlice.GetTemporalLayerNonReferenceFlag() {
        pcSlice.setNalUnitType(TLibCommon.NAL_UNIT_CODED_SLICE_TRAIL_N);
      }
    }

    // Do decoding refresh marking if any 
    pcSlice.decodingRefreshMarking(this.m_pocCRA, this.m_bRefreshPending, rcListPic);
    this.m_pcEncTop.selectReferencePictureSet(pcSlice, pocCurr, iGOPid,rcListPic);
    pcSlice.GetRPS().setNumberOfLongtermPictures(0);

    if pcSlice.checkThatAllRefPicsAreAvailable(rcListPic, pcSlice.GetRPS(), false) != 0 {
      pcSlice.createExplicitReferencePictureSetFromReference(rcListPic, pcSlice.GetRPS());
    }
    pcSlice.applyReferencePictureSet(rcListPic, pcSlice.GetRPS());

    if pcSlice.GetTLayer() > 0 {
      if pcSlice.isTemporalLayerSwitchingPoint(rcListPic, pcSlice.GetRPS()) || pcSlice.GetSPS().getTemporalIdNestingFlag() {
        if pcSlice.GetTemporalLayerNonReferenceFlag() {
          pcSlice.setNalUnitType(TLibCommon.NAL_UNIT_CODED_SLICE_TSA_N);
        }else{
          pcSlice.setNalUnitType(TLibCommon.NAL_UNIT_CODED_SLICE_TLA);
        }
      }else if pcSlice.isStepwiseTemporalLayerSwitchingPointCandidate(rcListPic, pcSlice.GetRPS()) {
        isSTSA:=true;
        for ii:=iGOPid+1;(ii<this.m_pcCfg.GetGOPSize() && isSTSA==true);ii++ {
          lTid := this.m_pcCfg.GetGOPEntry(ii).m_temporalId;
          if lTid==pcSlice.GetTLayer() {
            nRPS := pcSlice.GetSPS().getRPSList().getReferencePictureSet(ii);
            for jj:=0;jj<nRPS.getNumberOfPictures();jj++ {
              if nRPS.getUsed(jj) {
                tPoc:=this.m_pcCfg.GetGOPEntry(ii).m_POC+nRPS.getDeltaPOC(jj);
                kk:=0;
                for kk=0;kk<this.m_pcCfg.GetGOPSize();kk++ {
                  if this.m_pcCfg.GetGOPEntry(kk).m_POC==tPoc{
                    break;
                  }
                }
                tTid:=this.m_pcCfg.GetGOPEntry(kk).m_temporalId;
                if tTid >= pcSlice.GetTLayer() {
                  isSTSA=false;
                  break;
                }
              }
            }
          }
        }
        if isSTSA==true { 
          if pcSlice.GetTemporalLayerNonReferenceFlag() {
            pcSlice.setNalUnitType(TLibCommon.NAL_UNIT_CODED_SLICE_STSA_N);
          }else{
            pcSlice.setNalUnitType(TLibCommon.NAL_UNIT_CODED_SLICE_STSA_R);
          }
        }
      }
    }
    arrangeLongtermPicturesInRPS(pcSlice, rcListPic);
    TLibCommon.TComRefPicListModification* refPicListModification = pcSlice.GetRefPicListModification();
    refPicListModification.setRefPicListModificationFlagL0(0);
    refPicListModification.setRefPicListModificationFlagL1(0);
    pcSlice.setNumRefIdx(TLibCommon.REF_PIC_LIST_0,min(this.m_pcCfg.GetGOPEntry(iGOPid).m_numRefPicsActive,pcSlice.GetRPS().getNumberOfPictures()));
    pcSlice.setNumRefIdx(TLibCommon.REF_PIC_LIST_1,min(this.m_pcCfg.GetGOPEntry(iGOPid).m_numRefPicsActive,pcSlice.GetRPS().getNumberOfPictures()));

//#if ADAPTIVE_QP_SELECTION
    pcSlice.setTrQuant( this.m_pcEncTop.getTrQuant() );
//#endif      

    //  Set reference list
    pcSlice.setRefPicList ( rcListPic );

    //  Slice info. refinement
    if  (pcSlice.GetSliceType() == TLibCommon.B_SLICE) && (pcSlice.GetNumRefIdx(TLibCommon.REF_PIC_LIST_1) == 0) {
      pcSlice.setSliceType ( TLibCommon.P_SLICE );
    }

    if pcSlice.GetSliceType() != TLibCommon.B_SLICE || !pcSlice.GetSPS().getUseLComb() {
      pcSlice.setNumRefIdx(TLibCommon.REF_PIC_LIST_C, 0);
      pcSlice.setRefPicListCombinationFlag(false);
      pcSlice.setRefPicListModificationFlagLC(false);
    }else{
      pcSlice.setRefPicListCombinationFlag(pcSlice.GetSPS().getUseLComb());
      pcSlice.setNumRefIdx(TLibCommon.REF_PIC_LIST_C, pcSlice.GetNumRefIdx(TLibCommon.REF_PIC_LIST_0));
    }

    if pcSlice.GetSliceType() == TLibCommon.B_SLICE {
      pcSlice.setColFromL0Flag(1-uiColDir);
      bLowDelay := true;
       iCurrPOC  := pcSlice.GetPOC();
      iRefIdx := 0;

      for iRefIdx = 0; iRefIdx < pcSlice.GetNumRefIdx(TLibCommon.REF_PIC_LIST_0) && bLowDelay; iRefIdx++ {
        if pcSlice.GetRefPic(TLibCommon.REF_PIC_LIST_0, iRefIdx).getPOC() > iCurrPOC {
          bLowDelay = false;
        }
      }
      for iRefIdx = 0; iRefIdx < pcSlice.GetNumRefIdx(TLibCommon.REF_PIC_LIST_1) && bLowDelay; iRefIdx++ {
        if pcSlice.GetRefPic(TLibCommon.REF_PIC_LIST_1, iRefIdx).getPOC() > iCurrPOC {
          bLowDelay = false;
        }
      }

      pcSlice.SetCheckLDC(bLowDelay);  
    }

    uiColDir = 1-uiColDir;

    //-------------------------------------------------------------
    pcSlice.setRefPOCList();

    pcSlice.setNoBackPredFlag( false );
    if  pcSlice.GetSliceType() == TLibCommon.B_SLICE && !pcSlice.GetRefPicListCombinationFlag() {
      if  pcSlice.GetNumRefIdx(TLibCommon.B_SLICE( 0 ) ) == pcSlice.GetNumRefIdx(TLibCommon.B_SLICE( 1 ) ) {
        pcSlice.setNoBackPredFlag( true );
        var i int;
        for i=0; i < pcSlice.GetNumRefIdx(TLibCommon.B_SLICE( 1 ) ); i++ {
          if  pcSlice.GetRefPOC(TLibCommon.B_SLICE(1), i) != pcSlice.GetRefPOC(TLibCommon.B_SLICE(0), i) {
            pcSlice.setNoBackPredFlag( false );
            break;
          }
        }
      }
    }

    if pcSlice.GetNoBackPredFlag() {
      pcSlice.setNumRefIdx(TLibCommon.REF_PIC_LIST_C, 0);
    }
    pcSlice.generateCombinedList();

    if this.m_pcEncTop.getTMVPModeId() == 2 {
      if iGOPid == 0{ // first picture in SOP (i.e. forward B)
        pcSlice.setEnableTMVPFlag(0);
      }else{
        // Note: pcSlice.GetColFromL0Flag() is assumed to be always 0 and getcolRefIdx() is always 0.
        pcSlice.setEnableTMVPFlag(1);
      }
      pcSlice.GetSPS().setTMVPFlagsPresent(1);
    }else if this.m_pcEncTop.getTMVPModeId() == 1 {
      pcSlice.GetSPS().setTMVPFlagsPresent(1);
      pcSlice.setEnableTMVPFlag(1);
    }else{
      pcSlice.GetSPS().setTMVPFlagsPresent(0);
      pcSlice.setEnableTMVPFlag(0);
    }
    /////////////////////////////////////////////////////////////////////////////////////////////////// Compress a slice
    //  Slice compression
    if this.m_pcCfg.GetUseASR() {
      this.m_pcSliceEncoder.setSearchRange(pcSlice);
    }

     bGPBcheck:=false;
    if  pcSlice.GetSliceType() == TLibCommon.B_SLICE {
      if  pcSlice.GetNumRefIdx(TLibCommon.B_SLICE( 0 ) ) == pcSlice.GetNumRefIdx(TLibCommon.B_SLICE( 1 ) ) {
        bGPBcheck=true;
        var i int;
        for i=0; i < pcSlice.GetNumRefIdx(TLibCommon.B_SLICE( 1 ) ); i++ {
          if  pcSlice.GetRefPOC(TLibCommon.B_SLICE(1), i) != pcSlice.GetRefPOC(TLibCommon.B_SLICE(0), i) {
            bGPBcheck=false;
            break;
          }
        }
      }
    }
    if bGPBcheck {
      pcSlice.setMvdL1ZeroFlag(true);
    }else{
      pcSlice.setMvdL1ZeroFlag(false);
    }
    pcPic.GetSlice(pcSlice.GetSliceIdx()).setMvdL1ZeroFlag(pcSlice.GetMvdL1ZeroFlag());

//#if RATE_CONTROL_LAMBDA_DOMAIN
     sliceQP              := pcSlice.GetSliceQp();
     lambda            := float64(0.0);
     actualHeadBits       := 0;
     actualTotalBits      := 0;
     estimatedBits        := 0;
     tmpBitsBeforeWriting := 0;
    if  this.m_pcCfg.GetUseRateCtrl() {
       frameLevel := this.m_pcRateCtrl.getRCSeq().getGOPID2Level( iGOPid );
      if  pcPic.GetSlice(0).getSliceType() == TLibCommon.I_SLICE {
        frameLevel = 0;
      }
      this.m_pcRateCtrl.initRCPic( frameLevel );
      estimatedBits = this.m_pcRateCtrl.getRCPic().getTargetBits();

      if  ( pcSlice.GetPOC() == 0 && this.m_pcCfg.GetInitialQP() > 0 ) || ( frameLevel == 0 && this.m_pcCfg.GetForceIntraQP() ) { // QP is specified
        sliceQP              = this.m_pcCfg.GetInitialQP();
            NumberBFrames := ( this.m_pcCfg.GetGOPSize() - 1 );
         dLambda_scale := 1.0 - TLibCommon.CLIP3( 0.0, 0.5, 0.05*float64(NumberBFrames) ).(float64);
         dQPFactor     := 0.57*dLambda_scale;
            SHIFT_QP      := 12;
            bitdepth_luma_qp_scale := 0;
         qp_temp := float64( sliceQP )+ bitdepth_luma_qp_scale - TLibCommon.SHIFT_QP;
        lambda = dQPFactor*pow( 2.0, qp_temp/3.0 );
      }else if frameLevel == 0 {  // intra case, but use the model
        if  this.m_pcCfg.GetIntraPeriod() != 1 {  // do not refine allocated bits for all intra case
          bits := this.m_pcRateCtrl.getRCSeq().getLeftAverageBits();
          bits = this.m_pcRateCtrl.getRCSeq().getRefineBitsForIntra( bits );
          if bits < 200 {
            bits = 200;
          }
          this.m_pcRateCtrl.getRCPic().setTargetBits( bits );
        }

        listPreviousPicture := this.m_pcRateCtrl.getPicList();
        lambda  = this.m_pcRateCtrl.getRCPic().estimatePicLambda( listPreviousPicture );
        sliceQP = this.m_pcRateCtrl.getRCPic().estimatePicQP( lambda, listPreviousPicture );
      }else{    // normal case
        listPreviousPicture := this.m_pcRateCtrl.getPicList();
        lambda  = this.m_pcRateCtrl.getRCPic().estimatePicLambda( listPreviousPicture );
        sliceQP = this.m_pcRateCtrl.getRCPic().estimatePicQP( lambda, listPreviousPicture );
      }

      sliceQP = Clip3( -pcSlice.GetSPS().getQpBDOffsetY(), MAX_QP, sliceQP );
      this.m_pcRateCtrl.getRCPic().setPicEstQP( sliceQP );

      this.m_pcSliceEncoder.resetQP( pcPic, sliceQP, lambda );
    }
//#endif

     uiNumSlices := uint(1);

     uiInternalAddress := pcPic.GetNumPartInCU()-4;
     uiExternalAddress := pcPic.GetPicSym().getNumberOfCUsInFrame()-1;
     uiPosX := ( uiExternalAddress % pcPic.GetFrameWidthInCU() ) * TLibCommon.G_uiMaxCUWidth+ TLibCommon.G_auiRasterToPelX[ TLibCommon.G_auiZscanToRaster[uiInternalAddress] ];
     uiPosY := ( uiExternalAddress / pcPic.GetFrameWidthInCU() ) * TLibCommon.G_uiMaxCUHeight+ TLibCommon.G_auiRasterToPelY[ TLibCommon.G_auiZscanToRaster[uiInternalAddress] ];
     uiWidth := pcSlice.GetSPS().getPicWidthInLumaSamples();
     uiHeight := pcSlice.GetSPS().getPicHeightInLumaSamples();
    for uiPosX>=uiWidth||uiPosY>=uiHeight {
      uiInternalAddress--;
      uiPosX = ( uiExternalAddress % pcPic.GetFrameWidthInCU() ) * TLibCommon.G_uiMaxCUWidth+ TLibCommon.G_auiRasterToPelX[ TLibCommon.G_auiZscanToRaster[uiInternalAddress] ];
      uiPosY = ( uiExternalAddress / pcPic.GetFrameWidthInCU() ) * TLibCommon.G_uiMaxCUHeight+ TLibCommon.G_auiRasterToPelY[ TLibCommon.G_auiZscanToRaster[uiInternalAddress] ];
    }
    uiInternalAddress++;
    if uiInternalAddress==pcPic.GetNumPartInCU() {
      uiInternalAddress = 0;
      uiExternalAddress++;
    }
     uiRealEndAddress := uiExternalAddress*pcPic.GetNumPartInCU()+uiInternalAddress;

    var uiCummulativeTileWidth, uiCummulativeTileHeight uint;
    var  p, j int;
    var uiEncCUAddr uint;

    //set NumColumnsMinus1 and NumRowsMinus1
    pcPic.GetPicSym().setNumColumnsMinus1( pcSlice.GetPPS().getNumColumnsMinus1() );
    pcPic.GetPicSym().setNumRowsMinus1( pcSlice.GetPPS().getNumRowsMinus1() );

    //create the TLibCommon.TComTileArray
    pcPic.GetPicSym().xCreateTComTileArray();

    if pcSlice.GetPPS().getUniformSpacingFlag() == 1 {
      //set the width for each tile
      for j=0; j < pcPic.GetPicSym().getNumRowsMinus1()+1; j++ {
        for p=0; p < pcPic.GetPicSym().getNumColumnsMinus1()+1; p++ {
          pcPic.GetPicSym().getTComTile( j * (pcPic.GetPicSym().getNumColumnsMinus1()+1) + p ).setTileWidth( (p+1)*pcPic.GetPicSym().getFrameWidthInCU()/(pcPic.GetPicSym().getNumColumnsMinus1()+1)- 
            						(p*pcPic.GetPicSym().getFrameWidthInCU())/(pcPic.GetPicSym().getNumColumnsMinus1()+1) );
        }
      }

      //set the height for each tile
      for j=0; j < pcPic.GetPicSym().getNumColumnsMinus1()+1; j++ {
        for p=0; p < pcPic.GetPicSym().getNumRowsMinus1()+1; p++ {
          pcPic.GetPicSym().getTComTile( p * (pcPic.GetPicSym().getNumColumnsMinus1()+1) + j ).setTileHeight( (p+1)*pcPic.GetPicSym().getFrameHeightInCU()/(pcPic.GetPicSym().getNumRowsMinus1()+1)- 
             (p*pcPic.GetPicSym().getFrameHeightInCU())/(pcPic.GetPicSym().getNumRowsMinus1()+1) );   
        }
      }
    }else{
      //set the width for each tile
      for j=0; j < pcPic.GetPicSym().getNumRowsMinus1()+1; j++ {
        uiCummulativeTileWidth = 0;
        for p=0; p < pcPic.GetPicSym().getNumColumnsMinus1(); p++ {
          pcPic.GetPicSym().getTComTile( j * (pcPic.GetPicSym().getNumColumnsMinus1()+1) + p ).setTileWidth( pcSlice.GetPPS().getColumnWidth(p) );
          uiCummulativeTileWidth += pcSlice.GetPPS().getColumnWidth(p);
        }
        pcPic.GetPicSym().getTComTile(j * (pcPic.GetPicSym().getNumColumnsMinus1()+1) + p).setTileWidth( pcPic.GetPicSym().getFrameWidthInCU()-uiCummulativeTileWidth );
      }

      //set the height for each tile
      for j=0; j < pcPic.GetPicSym().getNumColumnsMinus1()+1; j++ {
        uiCummulativeTileHeight = 0;
        for p=0; p < pcPic.GetPicSym().getNumRowsMinus1(); p++ {
          pcPic.GetPicSym().getTComTile( p * (pcPic.GetPicSym().getNumColumnsMinus1()+1) + j ).setTileHeight( pcSlice.GetPPS().getRowHeight(p) );
          uiCummulativeTileHeight += pcSlice.GetPPS().getRowHeight(p);
        }
        pcPic.GetPicSym().getTComTile(p * (pcPic.GetPicSym().getNumColumnsMinus1()+1) + j).setTileHeight( pcPic.GetPicSym().getFrameHeightInCU()-uiCummulativeTileHeight );
      }
    }
    //intialize each tile of the current picture
    pcPic.GetPicSym().xInitTiles();

    // Allocate some coders, now we know how many tiles there are.
     iNumSubstreams := pcSlice.GetPPS().getNumSubstreams();

    //generate the Coding Order Map and Inverse Coding Order Map
     uiEncCUAddr=0
    for p=0; p<pcPic.GetPicSym().getNumberOfCUsInFrame(); p++ {
      pcPic.GetPicSym().setCUOrderMap(p, uiEncCUAddr);
      pcPic.GetPicSym().setInverseCUOrderMap(uiEncCUAddr, p);
      uiEncCUAddr = pcPic.GetPicSym().xCalculateNxtCUAddr(uiEncCUAddr)
    }
    pcPic.GetPicSym().setCUOrderMap(pcPic.GetPicSym().getNumberOfCUsInFrame(), pcPic.GetPicSym().getNumberOfCUsInFrame());    
    pcPic.GetPicSym().setInverseCUOrderMap(pcPic.GetPicSym().getNumberOfCUsInFrame(), pcPic.GetPicSym().getNumberOfCUsInFrame());

    // Allocate some coders, now we know how many tiles there are.
    this.m_pcEncTop.createWPPCoders(iNumSubstreams);
    pcSbacCoders = this.m_pcEncTop.getSbacCoders();
    pcSubstreamsOut = TLibCommon.NewTComOutputBitstream(iNumSubstreams);

     uiStartCUAddrSliceIdx := uint(0); // used to index "m_uiStoredStartCUAddrForEncodingSlice" containing locations of slice boundaries
     uiStartCUAddrSlice    := uint(0); // used to keep track of current slice's starting CU addr.
    pcSlice.setSliceCurStartCUAddr( uiStartCUAddrSlice ); // Setting "start CU addr" for current slice
    this.m_storedStartCUAddrForEncodingSlice.clear();

     uiStartCUAddrDependentSliceIdx := uint(0); // used to index "m_uiStoredStartCUAddrForEntropyEncodingSlice" containing locations of slice boundaries
     uiStartCUAddrDependentSlice    := uint(0); // used to keep track of current Dependent slice's starting CU addr.
    pcSlice.setDependentSliceCurStartCUAddr( uiStartCUAddrDependentSlice ); // Setting "start CU addr" for current Dependent slice

    this.m_storedStartCUAddrForEncodingDependentSlice.clear();
     uiNextCUAddr := uint(0);
    this.m_storedStartCUAddrForEncodingSlice.PushBack (uiNextCUAddr);
    uiStartCUAddrSliceIdx++;
    this.m_storedStartCUAddrForEncodingDependentSlice.PushBack(uiNextCUAddr);
    uiStartCUAddrDependentSliceIdx++;

    for uiNextCUAddr<uiRealEndAddress { // determine slice boundaries
      pcSlice.setNextSlice       ( false );
      pcSlice.setNextDependentSlice( false );
      assert(pcPic.GetNumAllocatedSlice() == uiStartCUAddrSliceIdx);
      this.m_pcSliceEncoder.precompressSlice( pcPic );
      this.m_pcSliceEncoder.compressSlice   ( pcPic );

      bNoBinBitConstraintViolated := (!pcSlice.isNextSlice() && !pcSlice.isNextDependentSlice());
      if pcSlice.isNextSlice() || (bNoBinBitConstraintViolated && this.m_pcCfg.GetSliceMode()==TLibCommon.AD_HOC_SLICES_FIXED_NUMBER_OF_LCU_IN_SLICE) {
        uiStartCUAddrSlice = pcSlice.GetSliceCurEndCUAddr();
        // Reconstruction slice
        this.m_storedStartCUAddrForEncodingSlice.push_back(uiStartCUAddrSlice);
        uiStartCUAddrSliceIdx++;
        // Dependent slice
        if uiStartCUAddrDependentSliceIdx>0 && this.m_storedStartCUAddrForEncodingDependentSlice[uiStartCUAddrDependentSliceIdx-1] != uiStartCUAddrSlice {
          this.m_storedStartCUAddrForEncodingDependentSlice.push_back(uiStartCUAddrSlice);
          uiStartCUAddrDependentSliceIdx++;
        }

        if uiStartCUAddrSlice < uiRealEndAddress {
          pcPic.allocateNewSlice();          
          pcPic.SetCurrSliceIdx                  ( uiStartCUAddrSliceIdx-1 );
          this.m_pcSliceEncoder.setSliceIdx           ( uiStartCUAddrSliceIdx-1 );
          pcSlice = pcPic.GetSlice               ( uiStartCUAddrSliceIdx-1 );
          pcSlice.copySliceInfo                  ( pcPic.GetSlice(0)      );
          pcSlice.setSliceIdx                    ( uiStartCUAddrSliceIdx-1 );
          pcSlice.setSliceCurStartCUAddr         ( uiStartCUAddrSlice      );
          pcSlice.setDependentSliceCurStartCUAddr  ( uiStartCUAddrSlice      );
          pcSlice.setSliceBits(0);
          uiNumSlices ++;
        }
      }else if pcSlice.isNextDependentSlice() || (bNoBinBitConstraintViolated && this.m_pcCfg.GetDependentSliceMode()==TLibCommon.SHARP_FIXED_NUMBER_OF_LCU_IN_DEPENDENT_SLICE) {
        uiStartCUAddrDependentSlice                                                     = pcSlice.GetDependentSliceCurEndCUAddr();
        this.m_storedStartCUAddrForEncodingDependentSlice.push_back(uiStartCUAddrDependentSlice);
        uiStartCUAddrDependentSliceIdx++;
        pcSlice.setDependentSliceCurStartCUAddr( uiStartCUAddrDependentSlice );
      }else{
        uiStartCUAddrSlice                                                            = pcSlice.GetSliceCurEndCUAddr();
        uiStartCUAddrDependentSlice                                                     = pcSlice.GetDependentSliceCurEndCUAddr();
      }        

	  if uiStartCUAddrSlice > uiStartCUAddrDependentSlice{
      	uiNextCUAddr = uiStartCUAddrSlice ;
      }else{
      	uiNextCUAddr = uiStartCUAddrDependentSlice;
      }
    }
    this.m_storedStartCUAddrForEncodingSlice.push_back( pcSlice.GetSliceCurEndCUAddr());
    uiStartCUAddrSliceIdx++;
    this.m_storedStartCUAddrForEncodingDependentSlice.push_back(pcSlice.GetSliceCurEndCUAddr());
    uiStartCUAddrDependentSliceIdx++;

    pcSlice = pcPic.GetSlice(0);

    // SAO parameter estimation using non-deblocked pixels for LCU bottom and right boundary areas
    if this.m_pcCfg.GetSaoLcuBasedOptimization() && this.m_pcCfg.GetSaoLcuBoundary() {
      this.m_pcSAO.resetStats();
      this.m_pcSAO.calcSaoStatsCu_BeforeDblk( pcPic );
    }

    //-- Loop filter
    bLFCrossTileBoundary := pcSlice.GetPPS().getLoopFilterAcrossTilesEnabledFlag();
    this.m_pcLoopFilter.setCfg(bLFCrossTileBoundary);
    this.m_pcLoopFilter.loopFilterPic( pcPic );

    pcSlice = pcPic.GetSlice(0);
    if pcSlice.GetSPS().getUseSAO() {
      var LFCrossSliceBoundaryFlag map[int]bool;
      for s:=0; s< uiNumSlices; s++ {
      	if uiNumSlices==1 {
        	LFCrossSliceBoundaryFlag[s] =true;//:pcPic.GetSlice(s).getLFCrossSliceBoundaryFlag()) );
        }else{
        	LFCrossSliceBoundaryFlag[s] =pcPic.GetSlice(s).getLFCrossSliceBoundaryFlag();
        }
      }
      this.m_storedStartCUAddrForEncodingSlice.resize(uiNumSlices+1);
      pcPic.createNonDBFilterInfo(this.m_storedStartCUAddrForEncodingSlice, 0, &LFCrossSliceBoundaryFlag ,pcPic.GetPicSym().getNumTiles() ,bLFCrossTileBoundary);
    }


    pcSlice = pcPic.GetSlice(0);

    if pcSlice.GetSPS().getUseSAO() {
      this.m_pcSAO.createPicSaoInfo(pcPic, uiNumSlices);
    }

    /////////////////////////////////////////////////////////////////////////////////////////////////// File writing
    // Set entropy coder
    this.m_pcEntropyCoder.setEntropyCoder   ( this.m_pcCavlcCoder, pcSlice );

    // write various header sets. 
    if  this.m_bSeqFirst {
      nalu = TLibCommon.NewOutputNALUnit(TLibCommon.NAL_UNIT_VPS);
      this.m_pcEntropyCoder.setBitstream(&nalu.m_Bitstream);
      this.m_pcEntropyCoder.encodeVPS(this.m_pcEncTop.getVPS());
      writeRBSPTrailingBits(nalu.m_Bitstream);
      accessUnit.push_back(NewNALUnitEBSP(nalu));
//#if RATE_CONTROL_LAMBDA_DOMAIN
      actualTotalBits += uint(accessUnit.back().m_nalUnitData.str().size()) * 8;
//#endif

      nalu = NALUnit(TLibCommon.NAL_UNIT_SPS);
      this.m_pcEntropyCoder.setBitstream(&nalu.m_Bitstream);
      if this.m_bSeqFirst {
        pcSlice.GetSPS().setNumLongTermRefPicSPS(this.m_numLongTermRefPicSPS);
        for k := 0; k < this.m_numLongTermRefPicSPS; k++ {
          pcSlice.GetSPS().setLtRefPicPocLsbSps(k, this.m_ltRefPicPocLsbSps[k]);
          pcSlice.GetSPS().setUsedByCurrPicLtSPSFlag(k, this.m_ltRefPicUsedByCurrPicFlag[k]);
        }
      }
      if this.m_pcCfg.GetPictureTimingSEIEnabled()  {
         maxCU := this.m_pcCfg.GetSliceArgument() >> ( pcSlice.GetSPS().getMaxCUDepth() << 1);
        var numDU uint;
        if this.m_pcCfg.GetSliceMode() == 1 {
        	numDU = ( pcPic.GetNumCUsInFrame() / maxCU );
        }else{
        	numDU = ( 0 );
        }
        
        if pcPic.GetNumCUsInFrame() % maxCU != 0  {
          numDU ++;
        }
        pcSlice.GetSPS().getVuiParameters().setNumDU( numDU );
        pcSlice.GetSPS().setHrdParameters( this.m_pcCfg.GetFrameRate(), numDU, this.m_pcCfg.GetTargetBitrate(), ( this.m_pcCfg.GetIntraPeriod() > 0 ) );
      }
      if this.m_pcCfg.GetBufferingPeriodSEIEnabled() || this.m_pcCfg.GetPictureTimingSEIEnabled() {
        pcSlice.GetSPS().getVuiParameters().setHrdParametersPresentFlag( true );
      }
      this.m_pcEntropyCoder.encodeSPS(pcSlice.GetSPS());
      writeRBSPTrailingBits(nalu.m_Bitstream);
      accessUnit.push_back(NewNALUnitEBSP(nalu));
//#if RATE_CONTROL_LAMBDA_DOMAIN
      actualTotalBits += uint(accessUnit.back().m_nalUnitData.str().size()) * 8;
//#endif

      nalu = NALUnit(TLibCommon.NAL_UNIT_PPS);
      this.m_pcEntropyCoder.setBitstream(&nalu.m_Bitstream);
      this.m_pcEntropyCoder.encodePPS(pcSlice.GetPPS());
      writeRBSPTrailingBits(nalu.m_Bitstream);
      accessUnit.push_back(NewNALUnitEBSP(nalu));
//#if RATE_CONTROL_LAMBDA_DOMAIN
      actualTotalBits += unt(accessUnit.back().m_nalUnitData.str().size()) * 8;
//#endif

      if this.m_pcCfg.GetActiveParameterSetsSEIEnabled() {
        var sei_active_parameter_sets SEIActiveParameterSets; 
        sei_active_parameter_sets.activeVPSId = this.m_pcCfg.GetVPS().getVPSId(); 
        if this.m_pcCfg.GetActiveParameterSetsSEIEnabled()==2 {
        	sei_active_parameter_sets.activeSPSIdPresentFlag =  0;
        }else{
        	sei_active_parameter_sets.activeSPSIdPresentFlag =  0;
        }
        
        if sei_active_parameter_sets.activeSPSIdPresentFlag {
          sei_active_parameter_sets.activeSeqParamSetId = pcSlice.GetSPS().getSPSId(); 
        }
//#if !HLS_REMOVE_ACTIVE_PARAM_SET_SEI_EXT_FLAG
//        sei_active_parameter_sets.activeParamSetSEIExtensionFlag = 0;
//#endif // HLS_REMOVE_ACTIVE_PARAM_SET_SEI_EXT_FLAG

        nalu = NALUnit(TLibCommon.NAL_UNIT_SEI); 
        this.m_pcEntropyCoder.setBitstream(&nalu.m_Bitstream);
        this.m_seiWriter.writeSEImessage(nalu.m_Bitstream, sei_active_parameter_sets); 
        writeRBSPTrailingBits(nalu.m_Bitstream);
        accessUnit.push_back(NewNALUnitEBSP(nalu));
      }
//#if SEI_DISPLAY_ORIENTATION
      if this.m_pcCfg.GetDisplayOrientationSEIAngle() {
        var sei_display_orientation SEIDisplayOrientation;
        sei_display_orientation.cancelFlag = false;
        sei_display_orientation.horFlip = false;
        sei_display_orientation.verFlip = false;
        sei_display_orientation.anticlockwiseRotation = this.m_pcCfg.GetDisplayOrientationSEIAngle();

        nalu = NALUnit(TLibCommon.NAL_UNIT_SEI); 
        this.m_pcEntropyCoder.setBitstream(&nalu.m_Bitstream);
        this.m_seiWriter.writeSEImessage(nalu.m_Bitstream, sei_display_orientation); 
        writeRBSPTrailingBits(nalu.m_Bitstream);
        accessUnit.push_back(NewNALUnitEBSP(nalu));
      }
//#endif

      this.m_bSeqFirst = false;
    }

    if  ( this.m_pcCfg.GetPictureTimingSEIEnabled() ) &&
        ( pcSlice.GetSPS().getVuiParametersPresentFlag() ) && 
        ( ( pcSlice.GetSPS().getVuiParameters().getNalHrdParametersPresentFlag() ) || 
        ( pcSlice.GetSPS().getVuiParameters().getVclHrdParametersPresentFlag() ) ) {
      if pcSlice.GetSPS().getVuiParameters().getSubPicCpbParamsPresentFlag() {
        numDU := pcSlice.GetSPS().getVuiParameters().getNumDU();
        pictureTimingSEI.m_numDecodingUnitsMinus1     = ( numDU - 1 );
        pictureTimingSEI.m_duCommonCpbRemovalDelayFlag = 0;

        if pictureTimingSEI.m_numNalusInDuMinus1 == nil {
          pictureTimingSEI.m_numNalusInDuMinus1       = make([]uint, numDU );
        }
        if pictureTimingSEI.m_duCpbRemovalDelayMinus1  == nil {
          pictureTimingSEI.m_duCpbRemovalDelayMinus1  = make([]uint, numDU );
        }
        if accumBitsDU == nil {
          accumBitsDU                                  = make([]uint, numDU );
        }
        if accumNalsDU == nil {
          accumNalsDU                                  = make([]uint, numDU );
        }
      }
      pictureTimingSEI.m_auCpbRemovalDelay = this.m_totalCoded - this.m_lastBPSEI;
      pictureTimingSEI.m_picDpbOutputDelay = pcSlice.GetSPS().getNumReorderPics(0) + pcSlice.GetPOC() - this.m_totalCoded;
    }
    if  ( this.m_pcCfg.GetBufferingPeriodSEIEnabled() ) && ( pcSlice.GetSliceType() == TLibCommon.I_SLICE ) &&
        ( pcSlice.GetSPS().getVuiParametersPresentFlag() ) && 
        ( ( pcSlice.GetSPS().getVuiParameters().getNalHrdParametersPresentFlag() ) || 
        ( pcSlice.GetSPS().getVuiParameters().getVclHrdParametersPresentFlag() ) ) {
       nalu = NewOutputNALUnit(TLibCommon.NAL_UNIT_SEI);
      this.m_pcEntropyCoder.setEntropyCoder(this.m_pcCavlcCoder, pcSlice);
      this.m_pcEntropyCoder.setBitstream(&nalu.m_Bitstream);

      var sei_buffering_period SEIBufferingPeriod;
      
      uiInitialCpbRemovalDelay := uint(90000/2);                      // 0.5 sec
      sei_buffering_period.m_initialCpbRemovalDelay      [0][0]     = uiInitialCpbRemovalDelay;
      sei_buffering_period.m_initialCpbRemovalDelayOffset[0][0]     = uiInitialCpbRemovalDelay;
      sei_buffering_period.m_initialCpbRemovalDelay      [0][1]     = uiInitialCpbRemovalDelay;
      sei_buffering_period.m_initialCpbRemovalDelayOffset[0][1]     = uiInitialCpbRemovalDelay;

      dTmp := float64(pcSlice.GetSPS().getVuiParameters().getNumUnitsInTick()) / float64(pcSlice.GetSPS().getVuiParameters().getTimeScale());

      uiTmp := uint( dTmp * 90000.0 ); 
      uiInitialCpbRemovalDelay -= uiTmp;
      uiInitialCpbRemovalDelay -= uiTmp / ( pcSlice.GetSPS().getVuiParameters().getTickDivisorMinus2() + 2 );
      sei_buffering_period.m_initialAltCpbRemovalDelay      [0][0]  = uiInitialCpbRemovalDelay;
      sei_buffering_period.m_initialAltCpbRemovalDelayOffset[0][0]  = uiInitialCpbRemovalDelay;
      sei_buffering_period.m_initialAltCpbRemovalDelay      [0][1]  = uiInitialCpbRemovalDelay;
      sei_buffering_period.m_initialAltCpbRemovalDelayOffset[0][1]  = uiInitialCpbRemovalDelay;

      sei_buffering_period.m_altCpbParamsPresentFlag              = 0;
      sei_buffering_period.m_sps                                  = pcSlice.GetSPS();

      this.m_seiWriter.writeSEImessage( nalu.m_Bitstream, sei_buffering_period );
      writeRBSPTrailingBits(nalu.m_Bitstream);
      accessUnit.push_back(NewNALUnitEBSP(nalu));

      this.m_lastBPSEI = this.m_totalCoded;
      this.m_cpbRemovalDelay = 0;
    }
    this.m_cpbRemovalDelay ++;
    if ( this.m_pcEncTop.getRecoveryPointSEIEnabled() ) && ( pcSlice.GetSliceType() == TLibCommon.I_SLICE ) {
      // Recovery point SEI
       nalu := NewOutputNALUnit(TLibCommon.NAL_UNIT_SEI);
      this.m_pcEntropyCoder.setEntropyCoder(this.m_pcCavlcCoder, pcSlice);
      this.m_pcEntropyCoder.setBitstream(&nalu.m_Bitstream);

      var sei_recovery_point SEIRecoveryPoint;
      sei_recovery_point.m_recoveryPocCnt    = 0;
      sei_recovery_point.m_exactMatchingFlag =  pcSlice.GetPOC() == 0 ;
      sei_recovery_point.m_brokenLinkFlag    = false;

      this.m_seiWriter.writeSEImessage( nalu.m_Bitstream, sei_recovery_point );
      writeRBSPTrailingBits(nalu.m_Bitstream);
      accessUnit.push_back(NewNALUnitEBSP(nalu));
    }
    // use the main bitstream buffer for storing the marshalled picture 
    this.m_pcEntropyCoder.setBitstream(NULL);

    uiStartCUAddrSliceIdx = 0;
    uiStartCUAddrSlice    = 0; 

    uiStartCUAddrDependentSliceIdx = 0;
    uiStartCUAddrDependentSlice    = 0; 
    uiNextCUAddr                 = 0;
    pcSlice = pcPic.GetSlice(uiStartCUAddrSliceIdx);

    var processingState int;
    if pcSlice.GetSPS().getUseSAO() {
    	processingState = (TLibCommon.EXECUTE_INLOOPFILTER);
    }else{
    	processingState = (TLibCommon.ENCODE_SLICE);
    }
    
     skippedSlice:=false;
    for uiNextCUAddr < uiRealEndAddress {// Iterate over all slices
      switch processingState {
      	case TLibCommon.ENCODE_SLICE:
          pcSlice.setNextSlice       ( false );
          pcSlice.setNextDependentSlice( false );
          if uiNextCUAddr == this.m_storedStartCUAddrForEncodingSlice[uiStartCUAddrSliceIdx] {
            pcSlice = pcPic.GetSlice(uiStartCUAddrSliceIdx);
            if uiStartCUAddrSliceIdx > 0 && pcSlice.GetSliceType()!= TLibCommon.I_SLICE {
              pcSlice.checkColRefIdx(uiStartCUAddrSliceIdx, pcPic);
            }
            pcPic.SetCurrSliceIdx(uiStartCUAddrSliceIdx);
            this.m_pcSliceEncoder.setSliceIdx(uiStartCUAddrSliceIdx);
            assert(uiStartCUAddrSliceIdx == pcSlice.GetSliceIdx());
            // Reconstruction slice
            pcSlice.setSliceCurStartCUAddr( uiNextCUAddr );  // to be used in encodeSlice() + context restriction
            pcSlice.setSliceCurEndCUAddr  ( this.m_storedStartCUAddrForEncodingSlice[uiStartCUAddrSliceIdx+1 ] );
            // Dependent slice
            pcSlice.setDependentSliceCurStartCUAddr( uiNextCUAddr );  // to be used in encodeSlice() + context restriction
            pcSlice.setDependentSliceCurEndCUAddr  ( this.m_storedStartCUAddrForEncodingDependentSlice[uiStartCUAddrDependentSliceIdx+1 ] );

            pcSlice.setNextSlice       ( true );

            uiStartCUAddrSliceIdx++;
            uiStartCUAddrDependentSliceIdx++;
          }else if uiNextCUAddr == this.m_storedStartCUAddrForEncodingDependentSlice[uiStartCUAddrDependentSliceIdx] {
            // Dependent slice
            pcSlice.setDependentSliceCurStartCUAddr( uiNextCUAddr );  // to be used in encodeSlice() + context restriction
            pcSlice.setDependentSliceCurEndCUAddr  ( this.m_storedStartCUAddrForEncodingDependentSlice[uiStartCUAddrDependentSliceIdx+1 ] );

            pcSlice.setNextDependentSlice( true );

            uiStartCUAddrDependentSliceIdx++;
          }

          pcSlice.setRPS(pcPic.GetSlice(0).getRPS());
          pcSlice.setRPSidx(pcPic.GetSlice(0).getRPSidx());
          var uiDummyStartCUAddr, uiDummyBoundingCUAddr uint;
          this.m_pcSliceEncoder.xDetermineStartAndBoundingCUAddr(uiDummyStartCUAddr,uiDummyBoundingCUAddr,pcPic,true);

          uiInternalAddress = pcPic.GetPicSym().getPicSCUAddr(pcSlice.GetDependentSliceCurEndCUAddr()-1) % pcPic.GetNumPartInCU();
          uiExternalAddress = pcPic.GetPicSym().getPicSCUAddr(pcSlice.GetDependentSliceCurEndCUAddr()-1) / pcPic.GetNumPartInCU();
          uiPosX = ( uiExternalAddress % pcPic.GetFrameWidthInCU() ) * TLibCommon.G_uiMaxCUWidth+ TLibCommon.G_auiRasterToPelX[ TLibCommon.G_auiZscanToRaster[uiInternalAddress] ];
          uiPosY = ( uiExternalAddress / pcPic.GetFrameWidthInCU() ) * TLibCommon.G_uiMaxCUHeight+ TLibCommon.G_auiRasterToPelY[ TLibCommon.G_auiZscanToRaster[uiInternalAddress] ];
          uiWidth = pcSlice.GetSPS().getPicWidthInLumaSamples();
          uiHeight = pcSlice.GetSPS().getPicHeightInLumaSamples();
          while(uiPosX>=uiWidth||uiPosY>=uiHeight)
          {
            uiInternalAddress--;
            uiPosX = ( uiExternalAddress % pcPic.GetFrameWidthInCU() ) * TLibCommon.G_uiMaxCUWidth+ TLibCommon.G_auiRasterToPelX[ TLibCommon.G_auiZscanToRaster[uiInternalAddress] ];
            uiPosY = ( uiExternalAddress / pcPic.GetFrameWidthInCU() ) * TLibCommon.G_uiMaxCUHeight+ TLibCommon.G_auiRasterToPelY[ TLibCommon.G_auiZscanToRaster[uiInternalAddress] ];
          }
          uiInternalAddress++;
          if uiInternalAddress==pcPic.GetNumPartInCU() {
            uiInternalAddress = 0;
            uiExternalAddress = pcPic.GetPicSym().getCUOrderMap(pcPic.GetPicSym().getInverseCUOrderMap(uiExternalAddress)+1);
          }
          uiEndAddress := pcPic.GetPicSym().getPicSCUEncOrder(uiExternalAddress*pcPic.GetNumPartInCU()+uiInternalAddress);
          if uiEndAddress<=pcSlice.GetDependentSliceCurStartCUAddr() {
            var uiBoundingAddrSlice, uiBoundingAddrDependentSlice uint;
            uiBoundingAddrSlice          = this.m_storedStartCUAddrForEncodingSlice[uiStartCUAddrSliceIdx];          
            uiBoundingAddrDependentSlice = this.m_storedStartCUAddrForEncodingDependentSlice[uiStartCUAddrDependentSliceIdx];          
            uiNextCUAddr               = min(uiBoundingAddrSlice, uiBoundingAddrDependentSlice);
            if pcSlice.isNextSlice() {
              skippedSlice=true;
            }
            continue;
          }
          if skippedSlice {
            pcSlice.setNextSlice       ( true );
            pcSlice.setNextDependentSlice( false );
          }
          skippedSlice=false;
          pcSlice.allocSubstreamSizes( iNumSubstreams );
          for ui := 0 ; ui < iNumSubstreams; ui++ {
            pcSubstreamsOut[ui].clear();
          }

          this.m_pcEntropyCoder.setEntropyCoder   ( this.m_pcCavlcCoder, pcSlice );
          this.m_pcEntropyCoder.resetEntropy      ();
          // start slice NALunit 
          nalu := NewOutputNALUnit( pcSlice.GetNalUnitType(), pcSlice.GetTLayer() );
          bDependentSlice := (!pcSlice.isNextSlice());
          if !bDependentSlice {
            uiOneBitstreamPerSliceLength = 0; // start of a new slice
          }
          this.m_pcEntropyCoder.setBitstream(&nalu.m_Bitstream);
//#if RATE_CONTROL_LAMBDA_DOMAIN
          tmpBitsBeforeWriting = this.m_pcEntropyCoder.getNumberOfWrittenBits();
//#endif
          this.m_pcEntropyCoder.encodeSliceHeader(pcSlice);
//#if RATE_CONTROL_LAMBDA_DOMAIN
          actualHeadBits += ( this.m_pcEntropyCoder.getNumberOfWrittenBits() - tmpBitsBeforeWriting );
//#endif

          // is it needed?
          {
            if !bDependentSlice{
              pcBitstreamRedirect.writeAlignOne();
            }else{
              // We've not completed our slice header info yet, do the alignment later.
            }
            this.m_pcSbacCoder.init( this.m_pcBinCABAC );
            this.m_pcEntropyCoder.setEntropyCoder ( this.m_pcSbacCoder, pcSlice );
            this.m_pcEntropyCoder.resetEntropy    ();
            for ui = 0 ; ui < pcSlice.GetPPS().getNumSubstreams() ; ui++ {
              this.m_pcEntropyCoder.setEntropyCoder ( &pcSbacCoders[ui], pcSlice );
              this.m_pcEntropyCoder.resetEntropy    ();
            }
          }

          if pcSlice.isNextSlice() {
            // set entropy coder for writing
            this.m_pcSbacCoder.init( this.m_pcBinCABAC );
            {
              for ui = 0 ; ui < pcSlice.GetPPS().getNumSubstreams() ; ui++ {
                this.m_pcEntropyCoder.setEntropyCoder ( &pcSbacCoders[ui], pcSlice );
                this.m_pcEntropyCoder.resetEntropy    ();
              }
              pcSbacCoders[0].load(this.m_pcSbacCoder);
              this.m_pcEntropyCoder.setEntropyCoder ( &pcSbacCoders[0], pcSlice );  //ALF is written in substream #0 with CABAC coder #0 (see ALF param encoding below)
            }
            m_pcEntropyCoder.resetEntropy    ();
            // File writing
            if !bDependentSlice {
              this.m_pcEntropyCoder.setBitstream(pcBitstreamRedirect);
            }else
            {
              this.m_pcEntropyCoder.setBitstream(&nalu.m_Bitstream);
            }
            // for now, override the TILES_DECODER setting in order to write substreams.
            this.m_pcEntropyCoder.setBitstream    ( &pcSubstreamsOut[0] );

          }
          pcSlice.setFinalized(true);

          this.m_pcSbacCoder.load( &pcSbacCoders[0] );

          pcSlice.setTileOffstForMultES( uiOneBitstreamPerSliceLength );
          if !bDependentSlice {
            pcSlice.setTileLocationCount ( 0 );
            this.m_pcSliceEncoder.encodeSlice(pcPic, pcBitstreamRedirect, pcSubstreamsOut); // redirect is only used for CAVLC tile position info.
          }else{
            this.m_pcSliceEncoder.encodeSlice(pcPic, &nalu.m_Bitstream, pcSubstreamsOut); // nalu.m_Bitstream is only used for CAVLC tile position info.
          }

          {
            // Construct the final bitstream by flushing and concatenating substreams.
            // The final bitstream is either nalu.m_Bitstream or pcBitstreamRedirect;
            puiSubstreamSizes := pcSlice.GetSubstreamSizes();
            uiTotalCodedSize := 0; // for padding calcs.
            uiNumSubstreamsPerTile := iNumSubstreams;
            if iNumSubstreams > 1 {
              uiNumSubstreamsPerTile /= pcPic.GetPicSym().getNumTiles();
            }
            for ui = 0 ; ui < iNumSubstreams; ui++ {
              // Flush all substreams -- this includes empty ones.
              // Terminating bit and flush.
              this.m_pcEntropyCoder.setEntropyCoder   ( &pcSbacCoders[ui], pcSlice );
              this.m_pcEntropyCoder.setBitstream      (  &pcSubstreamsOut[ui] );
              this.m_pcEntropyCoder.encodeTerminatingBit( 1 );
              this.m_pcEntropyCoder.encodeSliceFinish();

              pcSubstreamsOut[ui].writeByteAlignment();   // Byte-alignment in slice_data() at end of sub-stream
              // Byte alignment is necessary between tiles when tiles are independent.
              uiTotalCodedSize += pcSubstreamsOut[ui].getNumberOfWrittenBits();

              bNextSubstreamInNewTile := ((ui+1) < iNumSubstreams)&& ((ui+1)%uiNumSubstreamsPerTile == 0);
              if bNextSubstreamInNewTile {
                pcSlice.setTileLocation(ui/uiNumSubstreamsPerTile, pcSlice.GetTileOffstForMultES()+(uiTotalCodedSize>>3));
              }
              if ui+1 < pcSlice.GetPPS().getNumSubstreams() {
                puiSubstreamSizes[ui] = pcSubstreamsOut[ui].getNumberOfWrittenBits();
              }
            }

            // Complete the slice header info.
            this.m_pcEntropyCoder.setEntropyCoder   ( this.m_pcCavlcCoder, pcSlice );
            this.m_pcEntropyCoder.setBitstream(&nalu.m_Bitstream);
            this.m_pcEntropyCoder.encodeTilesWPPEntryPoint( pcSlice );

            // Substreams...
            pcOut := pcBitstreamRedirect;
           offs := 0;
           nss := pcSlice.GetPPS().getNumSubstreams();
          if pcSlice.GetPPS().getEntropyCodingSyncEnabledFlag() {
            // 1st line present for WPP.
//#if DEPENDENT_SLICES
            offs = pcSlice.GetDependentSliceCurStartCUAddr()/pcSlice.GetPic().getNumPartInCU()/pcSlice.GetPic().getFrameWidthInCU();
//#else
//            offs = pcSlice.GetSliceCurStartCUAddr()/pcSlice.GetPic().getNumPartInCU()/pcSlice.GetPic().getFrameWidthInCU();
//#endif
            nss  = pcSlice.GetNumEntryPointOffsets()+1;
          }
          	for ui = 0 ; ui < nss; ui++ {
            	pcOut.addSubstream(&pcSubstreamsOut[ui+offs]);
            }
          }

          var uiBoundingAddrSlice, uiBoundingAddrDependentSlice	uint;
          uiBoundingAddrSlice          = this.m_storedStartCUAddrForEncodingSlice[uiStartCUAddrSliceIdx];          
          uiBoundingAddrDependentSlice = this.m_storedStartCUAddrForEncodingDependentSlice[uiStartCUAddrDependentSliceIdx];          
          uiNextCUAddr               = TLibCommon.MIN(uiBoundingAddrSlice, uiBoundingAddrDependentSlice);
          // If current NALU is the first NALU of slice (containing slice header) and more NALUs exist (due to multiple dependent slices) then buffer it.
          // If current NALU is the last NALU of slice and a NALU was buffered, then (a) Write current NALU (b) Update an write buffered NALU at approproate location in NALU list.
          bNALUAlignedWrittenToList    := false; // used to ensure current NALU is not written more than once to the NALU list.
          xWriteTileLocationToSliceHeader(nalu, pcBitstreamRedirect, pcSlice);
          accessUnit.push_back(NewNALUnitEBSP(nalu));
//#if RATE_CONTROL_LAMBDA_DOMAIN
          actualTotalBits += UInt(accessUnit.back().m_nalUnitData.str().size()) * 8;
//#endif
          bNALUAlignedWrittenToList = true; 
          uiOneBitstreamPerSliceLength += nalu.m_Bitstream.getNumberOfWrittenBits(); // length of bitstream after byte-alignment

          if !bNALUAlignedWrittenToList {
            nalu.m_Bitstream.writeAlignZero();
            
            accessUnit.push_back(NewNALUnitEBSP(nalu));
            uiOneBitstreamPerSliceLength += nalu.m_Bitstream.getNumberOfWrittenBits() + 24; // length of bitstream after byte-alignment + 3 byte startcode 0x000001
          }

          if  ( this.m_pcCfg.GetPictureTimingSEIEnabled() ) &&
              ( pcSlice.GetSPS().getVuiParametersPresentFlag() ) && 
              ( ( pcSlice.GetSPS().getVuiParameters().getNalHrdParametersPresentFlag() ) || 
              ( pcSlice.GetSPS().getVuiParameters().getVclHrdParametersPresentFlag() ) ) &&
              ( pcSlice.GetSPS().getVuiParameters().getSubPicCpbParamsPresentFlag() ) {
             numRBSPBytes := uint(0);
            for it := accessUnit.Front(); it != nil; it=it.Next() {
               //v:= it.Value.();
               numRBSPBytes_nal := uint((*it).m_nalUnitData.str().size());
//#if HM9_NALU_TYPES
              if (*it).m_nalUnitType != TLibCommon.NAL_UNIT_SEI && (*it).m_nalUnitType != TLibCommon.NAL_UNIT_SEI_SUFFIX {
//#else
//              if ((*it).m_nalUnitType != TLibCommon.NAL_UNIT_SEI)
//#endif
              
                numRBSPBytes += numRBSPBytes_nal;
              }
            }
            accumBitsDU[ pcSlice.GetSliceIdx() ] = ( numRBSPBytes << 3 );
            accumNalsDU[ pcSlice.GetSliceIdx() ] = uint(accessUnit.size());
          }
          processingState = TLibCommon.ENCODE_SLICE;

        case TLibCommon.EXECUTE_INLOOPFILTER:
          
            // set entropy coder for RD
            this.m_pcEntropyCoder.setEntropyCoder ( this.m_pcSbacCoder, pcSlice );
            if  pcSlice.GetSPS().getUseSAO() {
              this.m_pcEntropyCoder.resetEntropy();
              this.m_pcEntropyCoder.setBitstream( this.m_pcBitCounter );
              this.m_pcSAO.startSaoEnc(pcPic, this.m_pcEntropyCoder, this.m_pcEncTop.getRDSbacCoder(), m_pcEncTop.getRDGoOnSbacCoder());
              cSaoParam := pcSlice.GetPic().getPicSym().getSaoParam();

//#if SAO_CHROMA_LAMBDA
//#if SAO_ENCODING_CHOICE
              this.m_pcSAO.SAOProcess(&cSaoParam, pcPic.GetSlice(0).getLambdaLuma(), pcPic.GetSlice(0).getLambdaChroma(), pcPic.GetSlice(0).getDepth());
//#else
//              this.m_pcSAO.SAOProcess(&cSaoParam, pcPic.GetSlice(0).getLambdaLuma(), pcPic.GetSlice(0).getLambdaChroma());
//#endif
//#else
//              this.m_pcSAO.SAOProcess(&cSaoParam, pcPic.GetSlice(0).getLambda());
//#endif
              this.m_pcSAO.endSaoEnc();
              this.m_pcSAO.PCMLFDisableProcess(pcPic);
            }
//#if SAO_RDO
            this.m_pcEntropyCoder.setEntropyCoder ( this.m_pcCavlcCoder, pcSlice );
//#endif
            processingState = TLibCommon.ENCODE_SLICE;

            for s:=0; s< uiNumSlices; s++ {
              if pcSlice.GetSPS().getUseSAO() {
                pcPic.GetSlice(s).setSaoEnabledFlag((pcSlice.GetPic().getPicSym().getSaoParam().bSaoFlag[0]==1));
              }
            }
        default:
            fmt.Printf("Not a supported encoding state\n");
            
        }
      } // end iteration over slices

      if pcSlice.GetSPS().getUseSAO() {
        if pcSlice.GetSPS().getUseSAO() {
          this.m_pcSAO.destroyPicSaoInfo();
        }
        pcPic.destroyNonDBFilterInfo();
      }

      pcPic.compressMotion(); 
      
      //-- For time output for each slice
      dEncTime := time.Now().Sub(iBeforeTime) //Double dEncTime = (Double)(clock()-iBeforeTime) / CLOCKS_PER_SEC;

      var digestStr string;
      if this.m_pcCfg.GetDecodedPictureHashSEIEnabled() {
        //calculate MD5sum for entire reconstructed picture 
        
        var sei_recon_picture_digest	SEIDecodedPictureHash;
        if this.m_pcCfg.GetDecodedPictureHashSEIEnabled() == 1 {
          sei_recon_picture_digest.method = SEIDecodedPictureHash::MD5;
          calcMD5(*pcPic.GetPicYuvRec(), sei_recon_picture_digest.digest);
          digestStr = digestToString(sei_recon_picture_digest.digest, 16);
        }else if(this.m_pcCfg.GetDecodedPictureHashSEIEnabled() == 2{
          sei_recon_picture_digest.method = SEIDecodedPictureHash::CRC;
          calcCRC(*pcPic.GetPicYuvRec(), sei_recon_picture_digest.digest);
          digestStr = digestToString(sei_recon_picture_digest.digest, 2);
        }
        else if(this.m_pcCfg.GetDecodedPictureHashSEIEnabled() == 3)
        {
          sei_recon_picture_digest.method = SEIDecodedPictureHash::CHECKSUM;
          calcChecksum(*pcPic.GetPicYuvRec(), sei_recon_picture_digest.digest);
          digestStr = digestToString(sei_recon_picture_digest.digest, 4);
        }
#if SUFFIX_SEI_NUT_DECODED_HASH_SEI
        OutputNALUnit nalu(TLibCommon.NAL_UNIT_SEI_SUFFIX, pcSlice.GetTLayer());
#else
        OutputNALUnit nalu(TLibCommon.NAL_UNIT_SEI, pcSlice.GetTLayer());
#endif

        //write the SEI messages 
        this.m_pcEntropyCoder.setEntropyCoder(this.m_pcCavlcCoder, pcSlice);
        this.m_seiWriter.writeSEImessage(nalu.m_Bitstream, sei_recon_picture_digest);
        writeRBSPTrailingBits(nalu.m_Bitstream);

#if SUFFIX_SEI_NUT_DECODED_HASH_SEI
        accessUnit.insert(accessUnit.end(), new NALUnitEBSP(nalu));
#else
        // insert the SEI message NALUnit before any Slice NALUnits 
        AccessUnit::iterator it = find_if(accessUnit.begin(), accessUnit.end(), mem_fun(&NALUnit::isSlice));
        accessUnit.insert(it, new NALUnitEBSP(nalu));
#endif
		
      }
//#if SEI_TEMPORAL_LEVEL0_INDEX
      if this.m_pcCfg.GetTemporalLevel0IndexSEIEnabled() {
        var sei_temporal_level0_index	SEITemporalLevel0Index;
        if pcSlice.GetRapPicFlag() {
          this.m_tl0Idx = 0;
          this.m_rapIdx = (this.m_rapIdx + 1) & 0xFF;
        }else{
          if pcSlice.GetTLayer()!=0 {
          	this.m_tl0Idx = (this.m_tl0Idx + 0) & 0xFF;
          }else{
          	this.m_tl0Idx = (this.m_tl0Idx + 1) & 0xFF;
          }
        }
        sei_temporal_level0_index.tl0Idx = this.m_tl0Idx;
        sei_temporal_level0_index.rapIdx = this.m_rapIdx;

         nalu := NewOutputNALUnit(TLibCommon.NAL_UNIT_SEI); 

        // write the SEI messages 
        this.m_pcEntropyCoder.setEntropyCoder(this.m_pcCavlcCoder, pcSlice);
        this.m_seiWriter.writeSEImessage(nalu.m_Bitstream, sei_temporal_level0_index);
        writeRBSPTrailingBits(nalu.m_Bitstream);

        // insert the SEI message NALUnit before any Slice NALUnits 
        //???it := find_if(accessUnit.begin(), accessUnit.end(), (&NALUnit::isSlice));
        accessUnit.insert(it, NewNALUnitEBSP(nalu));
      }
//#endif

      this.xCalculateAddPSNR( pcPic, pcPic.GetPicYuvRec(), accessUnit, dEncTime );

      if digestStr {
        if this.m_pcCfg.GetDecodedPictureHashSEIEnabled() == 1 {
          fmt.Printf(" [MD5:%s]", digestStr);
        }else if this.m_pcCfg.GetDecodedPictureHashSEIEnabled() == 2 {
          fmt.Printf(" [CRC:%s]", digestStr);
        }else if this.m_pcCfg.GetDecodedPictureHashSEIEnabled() == 3 {
          fmt.Printf(" [Checksum:%s]", digestStr);
        }
      }
//#if RATE_CONTROL_LAMBDA_DOMAIN
      if  this.m_pcCfg.GetUseRateCtrl() {
         effectivePercentage := this.m_pcRateCtrl.getRCPic().getEffectivePercentage();
         avgQP     := this.m_pcRateCtrl.getRCPic().calAverageQP();
         avgLambda := this.m_pcRateCtrl.getRCPic().calAverageLambda();
        if  avgLambda < 0.0 {
          avgLambda = lambda;
        }
        this.m_pcRateCtrl.getRCPic().updateAfterPicture( actualHeadBits, actualTotalBits, avgQP, avgLambda, effectivePercentage );
        this.m_pcRateCtrl.getRCPic().addToPictureLsit( this.m_pcRateCtrl.getPicList() );

        this.m_pcRateCtrl.getRCSeq().updateAfterPic( actualTotalBits );
        if  pcSlice.GetSliceType() != TLibCommon.I_SLICE {
          this.m_pcRateCtrl.getRCGOP().updateAfterPicture( actualTotalBits );
        }else{    // for intra picture, the estimated bits are used to update the current status in the GOP
          this.m_pcRateCtrl.getRCGOP().updateAfterPicture( estimatedBits );
        }
      }
//#else
//      if(this.m_pcCfg.GetUseRateCtrl())
//      {
//        UInt  frameBits = this.m_vRVM_RP[this.m_vRVM_RP.size()-1];
//        this.m_pcRateCtrl.updataRCFrameStatus((Int)frameBits, pcSlice.GetSliceType());
//      }
//#endif
      if  ( this.m_pcCfg.GetPictureTimingSEIEnabled() ) &&
          ( pcSlice.GetSPS().getVuiParametersPresentFlag() ) && 
          ( ( pcSlice.GetSPS().getVuiParameters().getNalHrdParametersPresentFlag() ) || 
          ( pcSlice.GetSPS().getVuiParameters().getVclHrdParametersPresentFlag() ) ) {
         nalu := NewOutputNALUnit(TLibCommon.NAL_UNIT_SEI, pcSlice.GetTLayer());
        vui := pcSlice.GetSPS().getVuiParameters();

        if vui.getSubPicCpbParamsPresentFlag() {
          var i	int;
          var ui64Tmp uint64;
          var uiTmp, uiPrev, uiCurr uint;

          uiPrev = 0;
          for i = 0; i < ( pictureTimingSEI.m_numDecodingUnitsMinus1 + 1 ); i ++ {
          	if i == 0 {
            	pictureTimingSEI.m_numNalusInDuMinus1[ i ]       =  ( accumNalsDU[ i ] ) ;
            }else{
            	pictureTimingSEI.m_numNalusInDuMinus1[ i ]       =  ( accumNalsDU[ i ] - accumNalsDU[ i - 1] - 1 );
            }
            ui64Tmp = ( ( ( accumBitsDU[ pictureTimingSEI.m_numDecodingUnitsMinus1 ] - accumBitsDU[ i ] ) * ( vui.getTimeScale() / vui.getNumUnitsInTick() ) * ( vui.getTickDivisorMinus2() + 2 ) ) /
                         ( this.m_pcCfg.GetTargetBitrate() << 10 ) );

            uiTmp = uint(ui64Tmp);
            if uiTmp >= ( vui.getTickDivisorMinus2() + 2 ) {
                  uiCurr = 0;
            }else{                                                     
            	uiCurr = ( vui.getTickDivisorMinus2() + 2 ) - uiTmp;
			}
				
            if i == pictureTimingSEI.m_numDecodingUnitsMinus1 {
             uiCurr = vui.getTickDivisorMinus2() + 2;
            }
            if uiCurr <= uiPrev {
             uiCurr = uiPrev + 1;
			}
            pictureTimingSEI.m_duCpbRemovalDelayMinus1[ i ]              = (uiCurr - uiPrev) - 1;
            uiPrev = uiCurr;
          }
        }
        this.m_pcEntropyCoder.setEntropyCoder(m_pcCavlcCoder, pcSlice);
        pictureTimingSEI.m_sps = pcSlice.GetSPS();
        this.m_seiWriter.writeSEImessage(nalu.m_Bitstream, pictureTimingSEI);
        writeRBSPTrailingBits(nalu.m_Bitstream);

        //??? AccessUnit::iterator it = find_if(accessUnit.begin(), accessUnit.end(), mem_fun(&NALUnit::isSlice));
        accessUnit.insert(it, NewNALUnitEBSP(nalu));
      }
      pcPic.GetPicYuvRec().copyToPic(pcPicYuvRecOut);

      pcPic.SetReconMark   ( true );
      this.m_bFirst = false;
      this.m_iNumPicCoded++;
      this.m_totalCoded ++;
      //logging: insert a newline at end of picture period 
      fmt.Printf("\n");
      fflush(stdout);

      //delete[] pcSubstreamsOut;
  }
*/
  
/*#if !RATE_CONTROL_LAMBDA_DOMAIN
  if(m_pcCfg.GetUseRateCtrl())
  {
    this.m_pcRateCtrl.updateRCGOPStatus();
  }
#endif
  delete pcBitstreamRedirect;

  if( accumBitsDU != NULL) delete accumBitsDU;
  if( accumNalsDU != NULL) delete accumNalsDU;

  assert ( this.m_iNumPicCoded == iNumPicRcvd );
  */
}

func (this *TEncGOP)  xWriteTileLocationToSliceHeader ( rNalu *OutputNALUnit, rpcBitstreamRedirect *TLibCommon.TComOutputBitstream, rpcSlice *TLibCommon.TComSlice){
  // Byte-align
  rNalu.m_Bitstream.WriteByteAlignment();   // Slice header byte-alignment

  // Perform bitstream concatenation
  if rpcBitstreamRedirect.GetNumberOfWrittenBits() > 0{
    uiBitCount := rpcBitstreamRedirect.GetNumberOfWrittenBits();
    if rpcBitstreamRedirect.GetByteStreamLength()>0 {
      pucStart  :=  rpcBitstreamRedirect.GetFIFO().Front();
      uiWriteByteCount := uint(0);
      for uiWriteByteCount < (uiBitCount >> 3) {
        v := pucStart.Value.(byte)
        uiBits := uint(v);
        rNalu.m_Bitstream.Write(uiBits, 8);
        pucStart = pucStart.Next();
        uiWriteByteCount++;
      }
    }
    uiBitsHeld := uint(uiBitCount & 0x07);
    for  uiIdx:=uint(0); uiIdx < uiBitsHeld; uiIdx++ {
      rNalu.m_Bitstream.Write(uint((rpcBitstreamRedirect.GetHeldBits() & (1 << (7-uiIdx))) >> (7-uiIdx)), 1);
    }          
  }

  this.m_pcEntropyCoder.setBitstream(rNalu.m_Bitstream);

  //delete rpcBitstreamRedirect;
  rpcBitstreamRedirect =  TLibCommon.NewTComOutputBitstream();
}

  
func (this *TEncGOP)  getGOPSize()    int      { return  this.m_iGopSize;  }
  
func (this *TEncGOP)  getListPic()   *list.List   { return this.m_pcListPic; }
  
func (this *TEncGOP)  printOutSummary      (  uiNumAllPicCoded uint){
  //assert (uiNumAllPicCoded == this.m_gcAnalyzeAll.getNumPic());
  
  //--CFG_KDY
  this.m_gcAnalyzeAll.setFrmRate( float64(this.m_pcCfg.GetFrameRate()) );
  this.m_gcAnalyzeI.setFrmRate( float64(this.m_pcCfg.GetFrameRate()) );
  this.m_gcAnalyzeP.setFrmRate( float64(this.m_pcCfg.GetFrameRate()) );
  this.m_gcAnalyzeB.setFrmRate( float64(this.m_pcCfg.GetFrameRate()) );
  
  //-- all
  fmt.Printf( "\n\nSUMMARY --------------------------------------------------------\n" );
  this.m_gcAnalyzeAll.printOut("a");
  
  fmt.Printf( "\n\nI Slices--------------------------------------------------------\n" );
  this.m_gcAnalyzeI.printOut("i");
  
  fmt.Printf( "\n\nP Slices--------------------------------------------------------\n" );
  this.m_gcAnalyzeP.printOut("p");
  
  fmt.Printf( "\n\nB Slices--------------------------------------------------------\n" );
  this.m_gcAnalyzeB.printOut("b");
  
//#if _SUMMARY_OUT_
//  this.m_gcAnalyzeAll.printSummaryOut();
//#endif
/*#if _SUMMARY_PIC_
  this.m_gcAnalyzeI.printSummary("I");
  this.m_gcAnalyzeP.printSummary("P");
  this.m_gcAnalyzeB.printSummary("B");
#endif
*/
  fmt.Printf("\nRVM: %.3lf\n" , this.xCalculateRVM());
}

func (this *TEncGOP)  preLoopFilterPicAll  ( pcPic *TLibCommon.TComPic, ruiDist *uint64, ruiBits *uint64 ){
  pcSlice := pcPic.GetSlice(pcPic.GetCurrSliceIdx());
  bCalcDist := false;
//#if VARYING_DBL_PARAMS
  this.m_pcLoopFilter.SetCfg(this.m_pcCfg.GetLFCrossTileBoundaryFlag());
//#else
//  this.m_pcLoopFilter.setCfg(m_pcCfg.GetLFCrossTileBoundaryFlag());
//#endif
  this.m_pcLoopFilter.LoopFilterPic( pcPic );
  
  this.m_pcEntropyCoder.setEntropyCoder ( this.m_pcEncTop.getRDGoOnSbacCoder(), pcSlice );
  this.m_pcEntropyCoder.resetEntropy    ();
  this.m_pcEntropyCoder.setBitstream    ( this.m_pcBitCounter );
  pcSlice = pcPic.GetSlice(0);
  if pcSlice.GetSPS().GetUseSAO() {
    var LFCrossSliceBoundaryFlag	map[int]bool//(1, true); //std::vector<Bool>
    var sliceStartAddress map[int]int; //std::vector<Int>
    sliceStartAddress[0] = 0;
    sliceStartAddress[1] = int(pcPic.GetNumCUsInFrame()* pcPic.GetNumPartInCU());
    pcPic.CreateNonDBFilterInfo(sliceStartAddress, 0, LFCrossSliceBoundaryFlag, 1, true);
  }
  
  if pcSlice.GetSPS().GetUseSAO() {
    pcPic.DestroyNonDBFilterInfo();
  }
  
  this.m_pcEntropyCoder.resetEntropy    ();
  *ruiBits += uint64(this.m_pcEntropyCoder.getNumberOfWrittenBits());
  
  if !bCalcDist {
    *ruiDist = this.xFindDistortionFrame(pcPic.GetPicYuvOrg(), pcPic.GetPicYuvRec());
  }
}
  
func (this *TEncGOP)  getSliceEncoder()  *TEncSlice { return this.m_pcSliceEncoder; }

func (this *TEncGOP)  getNalUnitType( pocCurr int ) TLibCommon.NalUnitType{
  if pocCurr == 0 {
    return TLibCommon.NAL_UNIT_CODED_SLICE_IDR;
  }
  if pocCurr % int(this.m_pcCfg.GetIntraPeriod()) == 0 {
    if this.m_pcCfg.GetDecodingRefreshType() == 1 {
      return TLibCommon.NAL_UNIT_CODED_SLICE_CRA;
    }else if this.m_pcCfg.GetDecodingRefreshType() == 2 {
      return TLibCommon.NAL_UNIT_CODED_SLICE_IDR;
    }
  }
  if this.m_pocCRA > 0 {
    if pocCurr< this.m_pocCRA {
      // All leading pictures are being marked as TFD pictures here since current encoder uses all 
      // reference pictures while encoding leading pictures. An encoder can ensure that a leading 
      // picture can be still decodable when random accessing to a CRA/CRANT/BLA/BLANT picture by 
      // controlling the reference pictures used for encoding that leading picture. Such a leading 
      // picture need not be marked as a TFD picture.
      return TLibCommon.NAL_UNIT_CODED_SLICE_TFD;
    }
  }
  return TLibCommon.NAL_UNIT_CODED_SLICE_TRAIL_R;
}

func (this *TEncGOP) getLSB( poc,  maxLSB int) int{
  if poc >= 0 {
    return poc % maxLSB;
  }
      
  return (maxLSB - ((-poc) % maxLSB)) % maxLSB;
}

func (this *TEncGOP)  arrangeLongtermPicturesInRPS(pcSlice *TLibCommon.TComSlice, rcListPic *list.List){
  rps := pcSlice.GetRPS();
  if rps.GetNumberOfLongtermPictures()==0 {
    return;
  }

  // Arrange long-term reference pictures in the correct order of LSB and MSB,
  // and assign values for pocLSBLT and MSB present flag
  var longtermPicsPoc, longtermPicsLSB, indices	[TLibCommon.MAX_NUM_REF_PICS]int;
//#if REMOVE_LTRP_LSB_RESTRICTIONS
  var longtermPicsMSB	[TLibCommon.MAX_NUM_REF_PICS]int;
//#endif
  var mSBPresentFlag	[TLibCommon.MAX_NUM_REF_PICS]bool;
//  ::memset(longtermPicsPoc, 0, sizeof(longtermPicsPoc));    // Store POC values of LTRP
//  ::memset(longtermPicsLSB, 0, sizeof(longtermPicsLSB));    // Store POC LSB values of LTRP
//#if REMOVE_LTRP_LSB_RESTRICTIONS
//  ::memset(longtermPicsMSB, 0, sizeof(longtermPicsMSB));    // Store POC LSB values of LTRP
//#endif
//  ::memset(indices        , 0, sizeof(indices));            // Indices to aid in tracking sorted LTRPs
//  ::memset(mSBPresentFlag , 0, sizeof(mSBPresentFlag));     // Indicate if MSB needs to be present

  // Get the long-term reference pictures 
  offset := rps.GetNumberOfNegativePictures() + rps.GetNumberOfPositivePictures();
  var i, j, ctr int;
  maxPicOrderCntLSB := 1 << pcSlice.GetSPS().GetBitsForPOC();
  for i = rps.GetNumberOfPictures() - 1; i >= offset; i-- {
    longtermPicsPoc[ctr] = rps.GetPOC(i);                                  // LTRP POC
    longtermPicsLSB[ctr] = this.getLSB(longtermPicsPoc[ctr], maxPicOrderCntLSB); // LTRP POC LSB
    indices[ctr]      = i; 
//#if REMOVE_LTRP_LSB_RESTRICTIONS
    longtermPicsMSB[ctr] = longtermPicsPoc[ctr] - longtermPicsLSB[ctr];
//#endif
	ctr++
  }
  numLongPics := rps.GetNumberOfLongtermPictures();
  //assert(ctr == numLongPics);

//#if REMOVE_LTRP_LSB_RESTRICTIONS
  // Arrange pictures in decreasing order of MSB; 
  for i = 0; i < numLongPics; i++ {
    for j = 0; j < numLongPics - 1; j++ {
      if longtermPicsMSB[j] < longtermPicsMSB[j+1] {
        var tmp int;
        tmp = longtermPicsPoc[j];
        longtermPicsPoc[j] = longtermPicsPoc[j+1];
        longtermPicsPoc[j+1] = tmp;
        
        tmp = longtermPicsLSB[j];
        longtermPicsLSB[j] = longtermPicsLSB[j+1];
        longtermPicsLSB[j+1] = tmp;
        
        tmp = longtermPicsMSB[j];
        longtermPicsMSB[j] = longtermPicsMSB[j+1];
        longtermPicsMSB[j+1] = tmp;
        
        tmp = indices[j];
        indices[j] = indices[j+1];
        indices[j+1] = tmp;
      }
    }
  }
/*#else
  // Arrange LTR pictures in decreasing order of LSB
  for(i = 0; i < numLongPics; i++)
  {
    for(Int j = 0; j < numLongPics - 1; j++)
    {
      if(longtermPicsLSB[j] < longtermPicsLSB[j+1])
      {
        std::swap(longtermPicsPoc[j], longtermPicsPoc[j+1]);
        std::swap(longtermPicsLSB[j], longtermPicsLSB[j+1]);
        std::swap(indices[j]        , indices[j+1]        );
      }
    }
  }
  // Now for those pictures that have the same LSB, arrange them 
  // in increasing MSB cycle, or equivalently decreasing MSB
  for(i = 0; i < numLongPics;)    // i incremented using j
  {
    Int j = i + 1;
    Int pocLSB = longtermPicsLSB[i];
    for(; j < numLongPics; j++)
    {
      if(pocLSB != longtermPicsLSB[j])
      {
        break;
      }
    }
    // Last index upto which lsb equals pocLSB is j - 1 
    // Now sort based on the MSB values
    Int sta, end;
    for(sta = i; sta < j; sta++)
    {
      for(end = i; end < j - 1; end++)
      {
      // longtermPicsMSB = longtermPicsPoc - longtermPicsLSB
        if(longtermPicsPoc[end] - longtermPicsLSB[end] < longtermPicsPoc[end+1] - longtermPicsLSB[end+1])
        {
          std::swap(longtermPicsPoc[end], longtermPicsPoc[end+1]);
          std::swap(longtermPicsLSB[end], longtermPicsLSB[end+1]);
          std::swap(indices[end]        , indices[end+1]        );
        }
      }
    }
    i = j;
  }
#endif*/

  for i = 0; i < numLongPics; i++ {
    // Check if MSB present flag should be enabled.
    // Check if the buffer contains any pictures that have the same LSB.
    iterPic := rcListPic.Front();  
    var pcPic *TLibCommon.TComPic;
    for iterPic != nil {
      pcPic = iterPic.Value.(*TLibCommon.TComPic);
      if (this.getLSB(int(pcPic.GetPOC()), maxPicOrderCntLSB) == longtermPicsLSB[i])   &&     // Same LSB
                                      (pcPic.GetSlice(0).IsReferenced())     &&    // Reference picture
                                        (int(pcPic.GetPOC()) != longtermPicsPoc[i])    {  // Not the LTRP itself
        mSBPresentFlag[i] = true;
        break;
      }
      iterPic=iterPic.Next();      
    }
  }

  // tempArray for usedByCurr flag
  var tempArray	[TLibCommon.MAX_NUM_REF_PICS]bool; //::memset(tempArray, 0, sizeof(tempArray));
  for i = 0; i < numLongPics; i++ {
    tempArray[i] = rps.GetUsed(indices[i]);
  }
  // Now write the final values;
  ctr = 0;
  currMSB := 0;
  currLSB := 0;
  // currPicPoc = currMSB + currLSB
  currLSB = this.getLSB(pcSlice.GetPOC(), maxPicOrderCntLSB);  
  currMSB = pcSlice.GetPOC() - currLSB;

  for i = int(rps.GetNumberOfPictures()) - 1; i >= offset; i-- {
    rps.SetPOC                   (i, longtermPicsPoc[ctr]);
    rps.SetDeltaPOC              (i, - pcSlice.GetPOC() + longtermPicsPoc[ctr]);
    rps.SetUsed                  (i, tempArray[ctr]);
    rps.SetPocLSBLT              (i, longtermPicsLSB[ctr]);
    rps.SetDeltaPocMSBCycleLT    (i, (currMSB - (longtermPicsPoc[ctr] - longtermPicsLSB[ctr])) / maxPicOrderCntLSB);
    rps.SetDeltaPocMSBPresentFlag(i, mSBPresentFlag[ctr]);     

    //assert(rps.GetDeltaPocMSBCycleLT(i) >= 0);   // Non-negative value
    ctr++
  }
//#if DISALLOW_LTRP_REPETITIONS
  ctr = 1
  for i = rps.GetNumberOfPictures() - 1; i >= offset; i-- {
    for j = rps.GetNumberOfPictures() - 1 - ctr; j >= offset; j-- {
      // Here at the encoder we know that we have set the full POC value for the LTRPs, hence we 
      // don't have to check the MSB present flag values for this constraint.
      //assert( rps.GetPOC(i) != rps.GetPOC(j) ); // If assert fails, LTRP entry repeated in RPS!!!
    }
    ctr++
  }
//#endif
}

func (this *TEncGOP)  getRateCtrl()    *TEncRateCtrl   { return this.m_pcRateCtrl;  }

func (this *TEncGOP)  xInitGOP          (  iPOCLast,  iNumPicRcvd int, rcListPic *list.List, rcListPicYuvRecOut *list.List){
  //assert( iNumPicRcvd > 0 );
  //  Exception for the first frame
  if iPOCLast == 0 {
    this.m_iGopSize    = 1;
  }else{
    this.m_iGopSize    = this.m_pcCfg.GetGOPSize();
  }
  //assert (m_iGopSize > 0); 

  return;
}

func (this *TEncGOP)  xGetBuffer        (  rcListPic, rcListPicYuvRecOut *list.List,  iNumPicRcvd,  iTimeOffset int, rpcPic *TLibCommon.TComPic,  rpcPicYuvRecOut *TLibCommon.TComPicYuv, pocCurr int){
  var i int;
  //  Rec. output
  iterPicYuvRec := rcListPicYuvRecOut.Back();
  for  i = 0; i < iNumPicRcvd - iTimeOffset + 1; i++  {
    iterPicYuvRec=iterPicYuvRec.Prev();
  }
  
  rpcPicYuvRecOut = iterPicYuvRec.Value.(*TLibCommon.TComPicYuv);
  
  //  Current pic.
  iterPic       := rcListPic.Front();
  for iterPic != nil {
    rpcPic = iterPic.Value.(*TLibCommon.TComPic);
    rpcPic.SetCurrSliceIdx(0);
    if int(rpcPic.GetPOC()) == pocCurr {
      break;
    }
    iterPic=iterPic.Next();
  }
  
  //assert (rpcPic.GetPOC() == pocCurr);
  
  return;
}
  
func (this *TEncGOP)  xCalculateAddPSNR ( pcPic *TLibCommon.TComPic, pcPicD *TLibCommon.TComPicYuv, accessUnit *list.List, dEncTime float64){
  var     x, y int;
   uiSSDY  := uint64(0);
   uiSSDU  := uint64(0);
   uiSSDV  := uint64(0);
  
    dYPSNR  := float64(0.0);
    dUPSNR  := float64(0.0);
    dVPSNR  := float64(0.0);
  
  //===== calculate PSNR =====
    pOrg    := pcPic.GetPicYuvOrg().GetLumaAddr();
    pRec    := pcPicD.GetLumaAddr();
     iStride := pcPicD.GetStride();
  
  var   iWidth,   iHeight int;
  
  iWidth  = pcPicD.GetWidth () - this.m_pcEncTop.m_pcEncCfg.GetPad(0);
  iHeight = pcPicD.GetHeight() - this.m_pcEncTop.m_pcEncCfg.GetPad(1);
  
     iSize   := iWidth*iHeight;
  
  for y = 0; y < iHeight; y++ {
    for x = 0; x < iWidth; x++  {
      iDiff := int( pOrg[x] - pRec[x] );
      uiSSDY   += uint64(iDiff * iDiff);
    }
    pOrg = pOrg[iStride:];
    pRec = pRec[iStride:];
  }
  
  iHeight >>= 1;
  iWidth  >>= 1;
  iStride >>= 1;
  pOrg  = pcPic.GetPicYuvOrg().GetCbAddr();
  pRec  = pcPicD.GetCbAddr();
  
  for y = 0; y < iHeight; y++ {
    for x = 0; x < iWidth; x++ {
      iDiff := int( pOrg[x] - pRec[x] );
      uiSSDU += uint64(iDiff * iDiff);
    }
    pOrg = pOrg[iStride:];
    pRec = pRec[iStride:];
  }
  
  pOrg  = pcPic.GetPicYuvOrg().GetCrAddr();
  pRec  = pcPicD.GetCrAddr();
  
  for y = 0; y < iHeight; y++ {
    for x = 0; x < iWidth; x++ {
       iDiff := int( pOrg[x] - pRec[x] );
      uiSSDV   += uint64(iDiff * iDiff);
    }
    pOrg = pOrg[iStride:];
    pRec = pRec[iStride:];
  }
  
   maxvalY := 255 << uint(TLibCommon.G_bitDepthY-8);
   maxvalC := 255 << uint(TLibCommon.G_bitDepthC-8);
   fRefValueY := float64 (maxvalY * maxvalY * iSize);
   fRefValueC := float64 (maxvalC * maxvalC * iSize) / 4.0;
  if uiSSDY!=0 {
  	dYPSNR            = 10.0 * math.Log10( fRefValueY / float64(uiSSDY)) ;
  }else{
    dYPSNR            = 99.99 ;
  }
  if uiSSDU!=0 {
  	dUPSNR            = 10.0 * math.Log10( fRefValueC / float64(uiSSDU)) ;
  }else{
    dUPSNR            = 99.99 ;
  }
  if uiSSDV!=0 {
  	dVPSNR            = 10.0 * math.Log10( fRefValueC / float64(uiSSDV)) ;
  }else{
    dVPSNR            = 99.99 ;
  }
 
  /* calculate the size of the access unit, excluding:
   *  - any AnnexB contributions (start_code_prefix, zero_byte, etc.,)
   *  - SEI NAL units
   */
  fmt.Printf("not implement yet xCalculateAddPSNR\n");
 
  numRBSPBytes := uint(0);
   /*
  for it := accessUnit.Front(); it != nil; it=it.Next() {
    //v := it.Value.()
    numRBSPBytes_nal := uint((*it).m_nalUnitData.str().size());
//#if VERBOSE_RATE
    printf("*** %6s numBytesInNALunit: %u\n", nalUnitTypeToString((*it).m_nalUnitType), numRBSPBytes_nal);
//#endif
//#if HM9_NALU_TYPES
    if (*it).m_nalUnitType != TLibCommon.NAL_UNIT_SEI && (*it).m_nalUnitType != TLibCommon.NAL_UNIT_SEI_SUFFIX {
//#else
//    if ((*it).m_nalUnitType != TLibCommon.NAL_UNIT_SEI)
//#endif
      numRBSPBytes += numRBSPBytes_nal;
    }
  }
  */

  uibits := numRBSPBytes * 8;
  this.m_vRVM_RP[len(this.m_vRVM_RP)] = int( uibits );

  //===== add PSNR =====
  this.m_gcAnalyzeAll.addResult (dYPSNR, dUPSNR, dVPSNR, float64(uibits));
  pcSlice := pcPic.GetSlice(0);
  if pcSlice.IsIntra() {
    this.m_gcAnalyzeI.addResult (dYPSNR, dUPSNR, dVPSNR, float64(uibits));
  }
  if pcSlice.IsInterP(){
    this.m_gcAnalyzeP.addResult (dYPSNR, dUPSNR, dVPSNR, float64(uibits));
  }
  if pcSlice.IsInterB(){
    this.m_gcAnalyzeB.addResult (dYPSNR, dUPSNR, dVPSNR, float64(uibits));
  }

  var c string;
  if pcSlice.IsIntra() {
  	c = "I";
  }else if pcSlice.IsInterP() {
  	c = "P";
  }else{
  	c = "B";
  }
  if !pcSlice.IsReferenced() {
  	//c += 32;
  	  if pcSlice.IsIntra() {
  		c = "i";
	  }else if pcSlice.IsInterP() {
	  	c = "p";
	  }else{
	  	c = "b";
	  }
  }
//#if ADAPTIVE_QP_SELECTION
  fmt.Printf("POC %4d TId: %1d ( %c-SLICE, nQP %d QP %d ) %10d bits",
         pcSlice.GetPOC(),
         pcSlice.GetTLayer(),
         c,
         pcSlice.GetSliceQpBase(),
         pcSlice.GetSliceQp(),
         uibits );
/*#else
  printf("POC %4d TId: %1d ( %c-SLICE, QP %d ) %10d bits",
         pcSlice.GetPOC()-pcSlice.GetLastIDR(),
         pcSlice.GetTLayer(),
         c,
         pcSlice.GetSliceQp(),
         uibits );
#endif*/

  fmt.Printf(" [Y %6.4lf dB    U %6.4lf dB    V %6.4lf dB]", dYPSNR, dUPSNR, dVPSNR );
  fmt.Printf(" [ET %5.0f ]", dEncTime );
  
  for  iRefList := 0; iRefList < 2; iRefList++ {
    fmt.Printf(" [L%d ", iRefList);
    for iRefIndex := 0; iRefIndex < pcSlice.GetNumRefIdx(TLibCommon.RefPicList(iRefList)); iRefIndex++ {
      fmt.Printf ("%d ", pcSlice.GetRefPOC(TLibCommon.RefPicList(iRefList), iRefIndex)-pcSlice.GetLastIDR());
    }
    fmt.Printf("]");
  }
}
  
func (this *TEncGOP)  xFindDistortionFrame (pcPic0 *TLibCommon.TComPicYuv, pcPic1 *TLibCommon.TComPicYuv) uint64{
  var     x, y int;
  pSrc0   := pcPic0.GetLumaAddr();
  pSrc1   := pcPic1.GetLumaAddr();
  uiShift := 2 * TLibCommon.DISTORTION_PRECISION_ADJUSTMENT(uint(TLibCommon.G_bitDepthY-8)).(uint);
  var   iTemp int;
  
     iStride := pcPic0.GetStride();
     iWidth  := pcPic0.GetWidth();
     iHeight := pcPic0.GetHeight();
  
    uiTotalDiff := uint64(0);
  
  for y = 0; y < iHeight; y++ {
    for x = 0; x < iWidth; x++ {
      iTemp = int(pSrc0[x] - pSrc1[x]); 
      uiTotalDiff += uint64(iTemp*iTemp) >> uiShift;
    }
    pSrc0 = pSrc0[iStride:];
    pSrc1 = pSrc1[iStride:];
  }
  
  uiShift = 2 * TLibCommon.DISTORTION_PRECISION_ADJUSTMENT(uint(TLibCommon.G_bitDepthC-8)).(uint);
  iHeight >>= 1;
  iWidth  >>= 1;
  iStride >>= 1;
  
  pSrc0  = pcPic0.GetCbAddr();
  pSrc1  = pcPic1.GetCbAddr();
  
  for y = 0; y < iHeight; y++ {
    for x = 0; x < iWidth; x++ {
      iTemp = int(pSrc0[x] - pSrc1[x]); 
      uiTotalDiff += uint64(iTemp*iTemp) >> uiShift;
    }
    pSrc0 = pSrc0[iStride:];
    pSrc1 = pSrc1[iStride:];
  }
  
  pSrc0  = pcPic0.GetCrAddr();
  pSrc1  = pcPic1.GetCrAddr();
  
  for y = 0; y < iHeight; y++ {
    for x = 0; x < iWidth; x++ {
      iTemp = int(pSrc0[x] - pSrc1[x]); 
      uiTotalDiff += uint64(iTemp*iTemp) >> uiShift;
    }
    pSrc0 = pSrc0[iStride:];
    pSrc1 = pSrc1[iStride:];
  }
  
  return uiTotalDiff;
}

func (this *TEncGOP)  xCalculateRVM() float64{
  dRVM := float64(0);
  fmt.Printf("not implement yet xCalculateRVM\n");
  /*
  if this.m_pcCfg.GetGOPSize() == 1 && this.m_pcCfg.GetIntraPeriod() != 1 && this.m_pcCfg.GetFrameToBeEncoded() > TLibCommon.RVM_VCEGAM10_M * 2 {
    // calculate RVM only for lowdelay configurations
    std::vector<Double> vRL , vB;
    size_t N = this.m_vRVM_RP.size();
    vRL.resize( N );
    vB.resize( N );
    
    Int i;
    Double dRavg = 0 , dBavg = 0;
    vB[RVM_VCEGAM10_M] = 0;
    for( i = RVM_VCEGAM10_M + 1 ; i < N - RVM_VCEGAM10_M + 1 ; i++ )
    {
      vRL[i] = 0;
      for( Int j = i - RVM_VCEGAM10_M ; j <= i + RVM_VCEGAM10_M - 1 ; j++ )
        vRL[i] += this.m_vRVM_RP[j];
      vRL[i] /= ( 2 * RVM_VCEGAM10_M );
      vB[i] = vB[i-1] + this.m_vRVM_RP[i] - vRL[i];
      dRavg += this.m_vRVM_RP[i];
      dBavg += vB[i];
    }
    
    dRavg /= ( N - 2 * RVM_VCEGAM10_M );
    dBavg /= ( N - 2 * RVM_VCEGAM10_M );
    
    Double dSigamB = 0;
    for( i = RVM_VCEGAM10_M + 1 ; i < N - RVM_VCEGAM10_M + 1 ; i++ )
    {
      Double tmp = vB[i] - dBavg;
      dSigamB += tmp * tmp;
    }
    dSigamB = sqrt( dSigamB / ( N - 2 * RVM_VCEGAM10_M ) );
    
    Double f = sqrt( 12.0 * ( RVM_VCEGAM10_M - 1 ) / ( RVM_VCEGAM10_M + 1 ) );
    
    dRVM = dSigamB / dRavg * f;
  }
  */
  return( dRVM );
}
