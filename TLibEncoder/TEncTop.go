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
	"io"
	"os"
    "container/list"
    "fmt"
    "gohm/TLibCommon"
)

// ====================================================================================================================
// Class definition
// ====================================================================================================================
//to be delete
type TEncSampleAdaptiveOffset struct {
    TLibCommon.TComSampleAdaptiveOffset
}

func NewTEncSampleAdaptiveOffset() *TEncSampleAdaptiveOffset{
	return &TEncSampleAdaptiveOffset{};
}

type TEncTop struct {
    m_pcEncCfg *TEncCfg

    // picture
    m_iPOCLast         int        ///< time index (POC)
    m_iNumPicRcvd      int        ///< number of received pictures
    m_uiNumAllPicCoded uint       ///< number of coded pictures
    m_cListPic         *list.List ///< dynamic list of pictures

    // encoder search
    m_cSearch        *TEncSearch  ///< encoder search class
    m_pcEntropyCoder *TEncEntropy ///< entropy encoder
    m_pcCavlcCoder   *TEncCavlc   ///< CAVLC encoder
    // coding tool
    m_cTrQuant         *TLibCommon.TComTrQuant    ///< transform & quantization class
    m_cLoopFilter      *TLibCommon.TComLoopFilter ///< deblocking filter class
    m_cEncSAO          *TEncSampleAdaptiveOffset  ///< sample adaptive offset class
    m_cEntropyCoder    *TEncEntropy               ///< entropy encoder
    m_cCavlcCoder      *TEncCavlc                 ///< CAVLC encoder
    m_cSbacCoder       *TEncSbac                  ///< SBAC encoder
    m_cBinCoderCABAC   *TEncBinCABAC              ///< bin coder CABAC
    m_pcSbacCoders     []*TEncSbac                ///< SBAC encoders (to encode substreams )
    m_pcBinCoderCABACs []*TEncBinCABAC            ///< bin coders CABAC (one per substream)

    // processing unit
    m_cGOPEncoder   *TEncGOP   ///< GOP encoder
    m_cSliceEncoder *TEncSlice ///< slice encoder
    m_cCuEncoder    *TEncCu    ///< CU encoder
    // SPS
    m_cSPS *TLibCommon.TComSPS ///< SPS
    m_cPPS *TLibCommon.TComPPS ///< PPS
    // RD cost computation
    m_cBitCounter      *TLibCommon.TComBitCounter ///< bit counter for RD optimization
    m_cRdCost          *TEncRdCost                ///< RD cost computation class
    m_pppcRDSbacCoder  [][]*TEncSbac              ///< temporal storage for RD computation
    m_cRDGoOnSbacCoder *TEncSbac                  ///< going on SBAC model for RD stage
    //#if FAST_BIT_EST
    m_pppcBinCoderCABAC    [][]*TEncBinCABACCounter ///< temporal CABAC state storage for RD computation
    m_cRDGoOnBinCoderCABAC *TEncBinCABACCounter   ///< going on bin coder CABAC for RD stage
    //#else
    //  TEncBinCABAC***         m_pppcBinCoderCABAC;            ///< temporal CABAC state storage for RD computation
    //  TEncBinCABAC            m_cRDGoOnBinCoderCABAC;         ///< going on bin coder CABAC for RD stage
    //#endif
    m_iNumSubstreams         int                          ///< # of top-level elements allocated.
    m_pcBitCounters          []*TLibCommon.TComBitCounter ///< bit counters for RD optimization per substream
    m_pcRdCosts              []*TEncRdCost                ///< RD cost computation class per substream
    m_ppppcRDSbacCoders      [][][]*TEncSbac              ///< temporal storage for RD computation per substream
    m_pcRDGoOnSbacCoders     []*TEncSbac                  ///< going on SBAC model for RD stage per substream
    m_ppppcBinCodersCABAC    [][][]*TEncBinCABAC          ///< temporal CABAC state storage for RD computation per substream
    m_pcRDGoOnBinCodersCABAC []*TEncBinCABAC              ///< going on bin coder CABAC for RD stage per substream

    // quality control
    //m_cPreanalyzer	*TLibCommon.TComPic;                 ///< image characteristics analyzer for TM5-step3-like adaptive QP

    m_scalingList *TLibCommon.TComScalingList ///< quantization matrix information
    m_cRateCtrl   *TEncRateCtrl               ///< Rate control class
    
    m_pTraceFile   *os.File
}

func NewTEncTop() *TEncTop {
    return &TEncTop{m_iPOCLast: -1,
    				m_cListPic:list.New(),
    				m_cSearch:NewTEncSearch(),
    				m_cTrQuant:TLibCommon.NewTComTrQuant(),
    				m_cLoopFilter:TLibCommon.NewTComLoopFilter(),
    				m_cEncSAO:NewTEncSampleAdaptiveOffset(),
    				m_cEntropyCoder:NewTEncEntropy(),
    				m_cCavlcCoder:NewTEncCavlc(),
    				m_cSbacCoder:NewTEncSbac(),
    				m_cBinCoderCABAC:NewTEncBinCABAC(),
    				m_cGOPEncoder:NewTEncGOP(),
    				m_cSliceEncoder:NewTEncSlice(),
    				m_cCuEncoder:NewTEncCu(),
    				m_cSPS:TLibCommon.NewTComSPS(),
    				m_cPPS:TLibCommon.NewTComPPS(),
    				m_cRDGoOnBinCoderCABAC:NewTEncBinCABACCounter(),
    				m_cBitCounter:TLibCommon.NewTComBitCounter(),
    				m_cRdCost:NewTEncRdCost(),
    				m_cRDGoOnSbacCoder:NewTEncSbac(),
    				m_scalingList:TLibCommon.NewTComScalingList(),
    				m_cRateCtrl:NewTEncRateCtrl(),
    				}
}

func (this *TEncTop) Create(pchTraceFile string) {
	if pchTraceFile != "" {
        this.m_pTraceFile, _ = os.Create(pchTraceFile)
    } else {
        this.m_pTraceFile = nil
    }
    
    // initialize global variables
    TLibCommon.InitROM()

	this.m_cRDGoOnSbacCoder.init( this.m_cRDGoOnBinCoderCABAC );
	this.m_cRDGoOnBinCoderCABAC.SetSbac(this.m_cRDGoOnSbacCoder);
//#if FAST_BIT_EST
  	TLibCommon.ContextModel_BuildNextStateTable();
//#endif
    // create processing unit classes
    this.m_cGOPEncoder.create()
    this.m_cSliceEncoder.create(this.GetEncCfg().GetSourceWidth(), this.GetEncCfg().GetSourceHeight(), this.GetEncCfg().GetMaxCUWidth(), this.GetEncCfg().GetMaxCUHeight(), byte(this.GetEncCfg().GetMaxCUDepth()))
    this.m_cCuEncoder.create(byte(this.GetEncCfg().GetMaxCUDepth()), this.GetEncCfg().GetMaxCUWidth(), this.GetEncCfg().GetMaxCUHeight())
    if this.GetEncCfg().m_bUseSAO {
        fmt.Printf("not support SAO\n")
        /*
           this.m_cEncSAO.setSaoLcuBoundary(this.GetEncCfg().GetSaoLcuBoundary());
           this.m_cEncSAO.setSaoLcuBasedOptimization(this.GetEncCfg().GetSaoLcuBasedOptimization());
           this.m_cEncSAO.setMaxNumOffsetsPerPic(this.GetEncCfg().GetMaxNumOffsetsPerPic());
           this.m_cEncSAO.create( this.GetEncCfg().GetSourceWidth(), this.GetEncCfg().GetSourceHeight(), TLibCommon.G_uiMaxCUWidth, TLibCommon.G_uiMaxCUHeight );
           this.m_cEncSAO.createEncBuffer();
        */
    }
    //#if ADAPTIVE_QP_SELECTION
    if this.GetEncCfg().GetUseAdaptQpSelect() {
        this.m_cTrQuant.InitSliceQpDelta()
    }
    //#endif
    this.m_cLoopFilter.Create(this.GetEncCfg().GetMaxCUDepth())

    //#if RATE_CONTROL_LAMBDA_DOMAIN
    if this.GetEncCfg().m_RCEnableRateControl {
        this.m_cRateCtrl.init(this.GetEncCfg().m_framesToBeEncoded, this.GetEncCfg().m_RCTargetBitrate, this.GetEncCfg().m_iFrameRate,
            this.GetEncCfg().m_iGOPSize, this.GetEncCfg().m_iSourceWidth, this.GetEncCfg().m_iSourceHeight,
            int(this.GetEncCfg().GetMaxCUWidth()), int(this.GetEncCfg().GetMaxCUHeight()), this.GetEncCfg().m_RCKeepHierarchicalBit,
            this.GetEncCfg().m_RCUseLCUSeparateModel, this.GetEncCfg().m_GOPList)
    }
    //#else
    //  this.m_cRateCtrl.create(this.GetEncCfg().GetIntraPeriod(), this.GetEncCfg().GetGOPSize(), this.GetEncCfg().GetFrameRate(), this.GetEncCfg().GetTarthis.GetEncCfg().GetBitrate(), this.GetEncCfg().GetQP(), this.GetEncCfg().GetNumLCUInUnit(), this.GetEncCfg().GetSourceWidth(), this.GetEncCfg().GetSourceHeight(), TLibCommon.G_uiMaxCUWidth, TLibCommon.G_uiMaxCUHeight);
    //#endif
    // if SBAC-based RD optimization is used
    if this.GetEncCfg().m_bUseSBACRD {
        this.m_pppcRDSbacCoder = make([][]*TEncSbac, this.GetEncCfg().GetMaxCUDepth()+1)
        //#if FAST_BIT_EST
        this.m_pppcBinCoderCABAC = make([][]*TEncBinCABACCounter, this.GetEncCfg().GetMaxCUDepth()+1)
        //#else
        //    this.m_pppcBinCoderCABAC = new TEncBinCABAC** [TLibCommon.G_uiMaxCUDepth+1];
        //#endif

        for iDepth := 0; iDepth < int(this.GetEncCfg().GetMaxCUDepth()+1); iDepth++ {
            this.m_pppcRDSbacCoder[iDepth] = make([]*TEncSbac, TLibCommon.CI_NUM)
            //#if FAST_BIT_EST
            this.m_pppcBinCoderCABAC[iDepth] = make([]*TEncBinCABACCounter, TLibCommon.CI_NUM)
            //#else
            //      this.m_pppcBinCoderCABAC[iDepth] = new TEncBinCABAC* [CI_NUM];
            //#endif

            for iCIIdx := 0; iCIIdx < TLibCommon.CI_NUM; iCIIdx++ {
                this.m_pppcRDSbacCoder[iDepth][iCIIdx] = NewTEncSbac()
                //#if FAST_BIT_EST
                this.m_pppcBinCoderCABAC[iDepth][iCIIdx] = NewTEncBinCABACCounter()
                //#else
                //        this.m_pppcBinCoderCABAC [iDepth][iCIIdx] = new TEncBinCABAC;
                //#endif
                this.m_pppcRDSbacCoder[iDepth][iCIIdx].init(this.m_pppcBinCoderCABAC[iDepth][iCIIdx])
                this.m_pppcBinCoderCABAC[iDepth][iCIIdx].SetSbac(this.m_pppcRDSbacCoder[iDepth][iCIIdx])
            }
        }
    }
}

func (this *TEncTop) Destroy() {
	if this.m_pTraceFile != nil {
        this.m_pTraceFile.Close()
    }
    
    // destroy processing unit classes
    this.m_cGOPEncoder.destroy()
    this.m_cSliceEncoder.destroy()
    this.m_cCuEncoder.destroy()
    if this.m_cSPS.GetUseSAO() {
        fmt.Printf("not support SAO\n")
        /*
           this.m_cEncSAO.destroy();
           this.m_cEncSAO.destroyEncBuffer();
        */
    }
    this.m_cLoopFilter.Destroy()
    this.m_cRateCtrl.destroy()
    // SBAC RD
    if this.GetEncCfg().m_bUseSBACRD {
        var iDepth int
        for iDepth = 0; iDepth < int(this.GetEncCfg().GetMaxCUDepth()+1); iDepth++ {
            for iCIIdx := 0; iCIIdx < TLibCommon.CI_NUM; iCIIdx++ {
                //delete this.m_pppcRDSbacCoder[iDepth][iCIIdx];
                //delete this.m_pppcBinCoderCABAC[iDepth][iCIIdx];
            }
        }

        for iDepth = 0; iDepth < int(this.GetEncCfg().GetMaxCUDepth()+1); iDepth++ {
            //delete [] this.m_pppcRDSbacCoder[iDepth];
            //delete [] this.m_pppcBinCoderCABAC[iDepth];
        }

        //delete [] this.m_pppcRDSbacCoder;
        //delete [] this.m_pppcBinCoderCABAC;

        for ui := 0; ui < this.m_iNumSubstreams; ui++ {
            for iDepth = 0; iDepth < int(this.GetEncCfg().GetMaxCUDepth()+1); iDepth++ {
                for iCIIdx := 0; iCIIdx < TLibCommon.CI_NUM; iCIIdx++ {
                    //delete this.m_ppppcRDSbacCoders  [ui][iDepth][iCIIdx];
                    //delete this.m_ppppcBinCodersCABAC[ui][iDepth][iCIIdx];
                }
            }

            for iDepth = 0; iDepth < int(this.GetEncCfg().GetMaxCUDepth()+1); iDepth++ {
                //delete [] this.m_ppppcRDSbacCoders  [ui][iDepth];
                //delete [] this.m_ppppcBinCodersCABAC[ui][iDepth];
            }
            //delete[] this.m_ppppcRDSbacCoders  [ui];
            //delete[] this.m_ppppcBinCodersCABAC[ui];
        }
        //delete[] this.m_ppppcRDSbacCoders;
        //delete[] this.m_ppppcBinCodersCABAC;
    }
    //delete[] this.m_pcSbacCoders;
    //delete[] this.m_pcBinCoderCABACs;
    //delete[] this.m_pcRDGoOnSbacCoders;
    //delete[] this.m_pcRDGoOnBinCodersCABAC;
    //delete[] this.m_pcBitCounters;
    //delete[] this.m_pcRdCosts;

    // destroy ROM
    TLibCommon.DestroyROM()

    return
}

func (this *TEncTop) Init() {
    //var aTable4, aTable8,aTableLastPosVlcIndex []uint;

    // initialize SPS
    this.xInitSPS()

    /* set the VPS profile information */
    *this.GetEncCfg().m_cVPS.GetPTL() = *this.m_cSPS.GetPTL()
    //#if L0043_TIMING_INFO
    this.GetEncCfg().m_cVPS.GetTimingInfo().SetTimingInfoPresentFlag(false)
    //#endif
    // initialize PPS
    this.m_cPPS.SetSPS(this.m_cSPS)
    this.xInitPPS()
    this.xInitRPS()

    this.xInitPPSforTiles()

    // initialize processing unit classes
    this.m_cGOPEncoder.init(this)
    this.m_cSliceEncoder.init(this)
    this.m_cCuEncoder.init(this)

    // initialize transform & quantization class
    this.m_pcCavlcCoder = this.getCavlcCoder()

    this.m_cTrQuant.Init(1<<this.GetEncCfg().m_uiQuadtreeTULog2MaxSize,
        this.GetEncCfg().m_useRDOQ,
        this.GetEncCfg().m_useRDOQTS,
        true,
        this.GetEncCfg().m_useTransformSkipFast,
        //#if ADAPTIVE_QP_SELECTION
        this.GetEncCfg().m_bUseAdaptQpSelect)
    //#endif

    // initialize encoder search class
    this.m_cSearch.init(this.GetEncCfg(), this.m_cTrQuant, this.GetEncCfg().m_iSearchRange, this.GetEncCfg().m_bipredSearchRange,
        this.GetEncCfg().m_iFastSearch, 0, this.m_cEntropyCoder, this.m_cRdCost, this.getRDSbacCoder(), this.getRDGoOnSbacCoder())

    this.GetEncCfg().m_iMaxRefPicNum = 0
}

func (this *TEncTop) DeletePicBuffer() {
    for iterPic := this.m_cListPic.Front(); iterPic != nil; iterPic = iterPic.Next() {
        pcPic := iterPic.Value.(*TLibCommon.TComPic)
        pcPic.Destroy()
    }

    this.m_cListPic.Init()
}

func (this *TEncTop) CreateWPPCoders(iNumSubstreams int) {
    if this.m_pcSbacCoders != nil {
        return // already generated.
    }

    this.m_iNumSubstreams = iNumSubstreams
    this.m_pcSbacCoders = make([]*TEncSbac, iNumSubstreams)
    this.m_pcBinCoderCABACs = make([]*TEncBinCABAC, iNumSubstreams)
    this.m_pcRDGoOnSbacCoders = make([]*TEncSbac, iNumSubstreams)
    this.m_pcRDGoOnBinCodersCABAC = make([]*TEncBinCABAC, iNumSubstreams)
    this.m_pcBitCounters = make([]*TLibCommon.TComBitCounter, iNumSubstreams)
    this.m_pcRdCosts = make([]*TEncRdCost, iNumSubstreams)

    for ui := 0; ui < iNumSubstreams; ui++ {
        this.m_pcSbacCoders[ui] = NewTEncSbac()
        this.m_pcBinCoderCABACs[ui] = NewTEncBinCABAC()
        this.m_pcRDGoOnSbacCoders[ui] = NewTEncSbac()
        this.m_pcRDGoOnBinCodersCABAC[ui] = NewTEncBinCABAC()
        this.m_pcBitCounters[ui] = TLibCommon.NewTComBitCounter()
        this.m_pcRdCosts[ui] = NewTEncRdCost()

        this.m_pcRDGoOnSbacCoders[ui].init(this.m_pcRDGoOnBinCodersCABAC[ui])
        this.m_pcSbacCoders[ui].init(this.m_pcBinCoderCABACs[ui])
    }
    if this.GetEncCfg().m_bUseSBACRD {
        this.m_ppppcRDSbacCoders = make([][][]*TEncSbac, iNumSubstreams)
        this.m_ppppcBinCodersCABAC = make([][][]*TEncBinCABAC, iNumSubstreams)
        for ui := 0; ui < iNumSubstreams; ui++ {
            this.m_ppppcRDSbacCoders[ui] = make([][]*TEncSbac, this.GetEncCfg().GetMaxCUDepth()+1)
            this.m_ppppcBinCodersCABAC[ui] = make([][]*TEncBinCABAC, this.GetEncCfg().GetMaxCUDepth()+1)

            for iDepth := 0; iDepth < int(this.GetEncCfg().GetMaxCUDepth()+1); iDepth++ {
                this.m_ppppcRDSbacCoders[ui][iDepth] = make([]*TEncSbac, TLibCommon.CI_NUM)
                this.m_ppppcBinCodersCABAC[ui][iDepth] = make([]*TEncBinCABAC, TLibCommon.CI_NUM)

                for iCIIdx := 0; iCIIdx < TLibCommon.CI_NUM; iCIIdx++ {
                    this.m_ppppcRDSbacCoders[ui][iDepth][iCIIdx] = NewTEncSbac()
                    this.m_ppppcBinCodersCABAC[ui][iDepth][iCIIdx] = NewTEncBinCABAC()
                    this.m_ppppcRDSbacCoders[ui][iDepth][iCIIdx].init(this.m_ppppcBinCodersCABAC[ui][iDepth][iCIIdx])
                }
            }
        }
    }
}

// -------------------------------------------------------------------------------------------------------------------
// member access functions
// -------------------------------------------------------------------------------------------------------------------
func (this *TEncTop) xGetNewPicBuffer() *TLibCommon.TComPic { ///< get picture buffer which will be processed
    var rpcPic *TLibCommon.TComPic

    TLibCommon.SortPicList(this.m_cListPic)

    if this.m_cListPic.Len() >= (this.GetEncCfg().m_iGOPSize + this.GetEncCfg().GetMaxDecPicBuffering(TLibCommon.MAX_TLAYER-1) + 2) {

        //Int iSize = Int( this.m_cListPic.size() );
        for iterPic := this.m_cListPic.Front(); iterPic != nil; iterPic = iterPic.Next() {
            rpcPic = iterPic.Value.(*TLibCommon.TComPic)

            if rpcPic.GetSlice(0).IsReferenced() == false {
                break
            }
        }
    } else {
        if this.GetEncCfg().GetUseAdaptiveQP() {
            pcEPic := TLibCommon.NewTComPic()
            pcEPic.Create(this.GetEncCfg().m_iSourceWidth, this.GetEncCfg().m_iSourceHeight,
                this.GetEncCfg().GetMaxCUWidth(), this.GetEncCfg().GetMaxCUHeight(), this.GetEncCfg().GetMaxCUDepth(), this.m_cPPS.GetMaxCuDQPDepth()+1,
                this.GetEncCfg().m_conformanceWindow, this.GetEncCfg().m_defaultDisplayWindow, this.GetEncCfg().m_numReorderPics[:], false)
            rpcPic = pcEPic
        } else {
            rpcPic = TLibCommon.NewTComPic()
            rpcPic.Create(this.GetEncCfg().m_iSourceWidth, this.GetEncCfg().m_iSourceHeight,
                this.GetEncCfg().GetMaxCUWidth(), this.GetEncCfg().GetMaxCUHeight(), this.GetEncCfg().GetMaxCUDepth(), 0,
                this.GetEncCfg().m_conformanceWindow, this.GetEncCfg().m_defaultDisplayWindow, this.GetEncCfg().m_numReorderPics[:], false)
        }
        if this.GetEncCfg().GetUseSAO() {
            fmt.Printf("not support SAO\n")
            //rpcPic.GetPicSym().AllocSaoParam(this.m_cEncSAO);
        }
        this.m_cListPic.PushBack(rpcPic)
    }
    rpcPic.SetReconMark(false)

    this.m_iPOCLast++
    this.m_iNumPicRcvd++

    rpcPic.GetSlice(0).SetPOC(this.m_iPOCLast)
    // mark it should be extended
    rpcPic.GetPicYuvRec().SetBorderExtension(false)

    return rpcPic
}

func (this *TEncTop) xInitSPS() { ///< initialize SPS from encoder options
    profileTierLevel := this.m_cSPS.GetPTL().GetGeneralPTL()
    profileTierLevel.SetLevelIdc(int(this.GetEncCfg().m_level))
    profileTierLevel.SetTierFlag(this.GetEncCfg().m_levelTier != 0)
    profileTierLevel.SetProfileIdc(int(this.GetEncCfg().m_profile))
    profileTierLevel.SetProfileCompatibilityFlag(int(this.GetEncCfg().m_profile), true)
    //#if L0046_CONSTRAINT_FLAGS
    profileTierLevel.SetProgressiveSourceFlag(this.GetEncCfg().m_progressiveSourceFlag)
    profileTierLevel.SetInterlacedSourceFlag(this.GetEncCfg().m_interlacedSourceFlag)
    profileTierLevel.SetNonPackedConstraintFlag(this.GetEncCfg().m_nonPackedConstraintFlag)
    profileTierLevel.SetFrameOnlyConstraintFlag(this.GetEncCfg().m_frameOnlyConstraintFlag)
    //#endif

    if this.GetEncCfg().m_profile == TLibCommon.PROFILE_MAIN10 && TLibCommon.G_bitDepthY == 8 && TLibCommon.G_bitDepthC == 8 {
        /* The above constraint is equal to Profile::MAIN */
        profileTierLevel.SetProfileCompatibilityFlag(TLibCommon.PROFILE_MAIN, true)
    }
    if this.GetEncCfg().m_profile == TLibCommon.PROFILE_MAIN {
        /* A Profile::MAIN10 decoder can always decode Profile::MAIN */
        profileTierLevel.SetProfileCompatibilityFlag(TLibCommon.PROFILE_MAIN10, true)
    }
    /* XXX: should Main be marked as compatible with still picture? */
    /* XXX: may be a good idea to refactor the above into a function
     * that chooses the actual compatibility based upon options */

    this.m_cSPS.SetPicWidthInLumaSamples(uint(this.GetEncCfg().m_iSourceWidth))
    this.m_cSPS.SetPicHeightInLumaSamples(uint(this.GetEncCfg().m_iSourceHeight))
    this.m_cSPS.SetConformanceWindow(this.GetEncCfg().m_conformanceWindow)
    this.m_cSPS.SetMaxCUWidth(this.GetEncCfg().GetMaxCUWidth())
    this.m_cSPS.SetMaxCUHeight(this.GetEncCfg().GetMaxCUHeight())
    this.m_cSPS.SetMaxCUDepth(this.GetEncCfg().GetMaxCUDepth())
    this.m_cSPS.SetAddCUDepth(this.GetEncCfg().GetAddCUDepth());
    this.m_cSPS.SetMinTrDepth(0)
    this.m_cSPS.SetMaxTrDepth(1)

    this.m_cSPS.SetPCMLog2MinSize(this.GetEncCfg().m_uiPCMLog2MinSize)
    this.m_cSPS.SetUsePCM(this.GetEncCfg().m_usePCM)
    this.m_cSPS.SetPCMLog2MaxSize(this.GetEncCfg().m_pcmLog2MaxSize)

    this.m_cSPS.SetQuadtreeTULog2MaxSize(this.GetEncCfg().m_uiQuadtreeTULog2MaxSize)
    this.m_cSPS.SetQuadtreeTULog2MinSize(this.GetEncCfg().m_uiQuadtreeTULog2MinSize)
    this.m_cSPS.SetQuadtreeTUMaxDepthInter(this.GetEncCfg().m_uiQuadtreeTUMaxDepthInter)
    this.m_cSPS.SetQuadtreeTUMaxDepthIntra(this.GetEncCfg().m_uiQuadtreeTUMaxDepthIntra)

    this.m_cSPS.SetTMVPFlagsPresent(false)
    this.m_cSPS.SetUseLossless(this.GetEncCfg().m_useLossless)

    this.m_cSPS.SetMaxTrSize(1 << this.GetEncCfg().m_uiQuadtreeTULog2MaxSize)

    this.m_cSPS.SetUseLComb(this.GetEncCfg().m_bUseLComb)

    var i uint

    for i = 0; i < this.GetEncCfg().GetMaxCUDepth()-this.GetEncCfg().GetAddCUDepth(); i++ {
        this.m_cSPS.SetAMPAcc(i, int(TLibCommon.B2U(this.GetEncCfg().m_useAMP)))
        //this.m_cSPS.setAMPAcc( i, 1 );
    }

    this.m_cSPS.SetUseAMP(this.GetEncCfg().m_useAMP)

    for i = this.GetEncCfg().GetMaxCUDepth() - this.GetEncCfg().GetAddCUDepth(); i < this.GetEncCfg().GetMaxCUDepth(); i++ {
        this.m_cSPS.SetAMPAcc(i, 0)
    }

    this.m_cSPS.SetBitDepthY(TLibCommon.G_bitDepthY)
    this.m_cSPS.SetBitDepthC(TLibCommon.G_bitDepthC)

    this.m_cSPS.SetQpBDOffsetY(6 * (TLibCommon.G_bitDepthY - 8))
    this.m_cSPS.SetQpBDOffsetC(6 * (TLibCommon.G_bitDepthC - 8))

    this.m_cSPS.SetUseSAO(this.GetEncCfg().m_bUseSAO)

    this.m_cSPS.SetMaxTLayers(uint(this.GetEncCfg().m_maxTempLayer))
    this.m_cSPS.SetTemporalIdNestingFlag((this.GetEncCfg().m_maxTempLayer == 1))
    for i = 0; i < this.m_cSPS.GetMaxTLayers(); i++ {
        this.m_cSPS.SetMaxDecPicBuffering(uint(this.GetEncCfg().m_maxDecPicBuffering[i]), i)
        this.m_cSPS.SetNumReorderPics(this.GetEncCfg().m_numReorderPics[i], i)
    }
    this.m_cSPS.SetPCMBitDepthLuma(uint(TLibCommon.G_uiPCMBitDepthLuma))
    this.m_cSPS.SetPCMBitDepthChroma(uint(TLibCommon.G_uiPCMBitDepthChroma))
    this.m_cSPS.SetPCMFilterDisableFlag(this.GetEncCfg().m_bPCMFilterDisableFlag)
    this.m_cSPS.SetScalingListFlag(this.GetEncCfg().m_useScalingListId != 0)
    this.m_cSPS.SetUseStrongIntraSmoothing(this.GetEncCfg().m_useStrongIntraSmoothing)

    this.m_cSPS.SetVuiParametersPresentFlag(this.GetEncCfg().GetVuiParametersPresentFlag())
    if this.m_cSPS.GetVuiParametersPresentFlag() {
        pcVUI := this.m_cSPS.GetVuiParameters()
        pcVUI.SetAspectRatioInfoPresentFlag(this.GetEncCfg().GetAspectRatioIdc() != -1)
        pcVUI.SetAspectRatioIdc(this.GetEncCfg().GetAspectRatioIdc())
        pcVUI.SetSarWidth(this.GetEncCfg().GetSarWidth())
        pcVUI.SetSarHeight(this.GetEncCfg().GetSarHeight())
        pcVUI.SetOverscanInfoPresentFlag(this.GetEncCfg().GetOverscanInfoPresentFlag())
        pcVUI.SetOverscanAppropriateFlag(this.GetEncCfg().GetOverscanAppropriateFlag())
        pcVUI.SetVideoSignalTypePresentFlag(this.GetEncCfg().GetVideoSignalTypePresentFlag())
        pcVUI.SetVideoFormat(this.GetEncCfg().GetVideoFormat())
        pcVUI.SetVideoFullRangeFlag(this.GetEncCfg().GetVideoFullRangeFlag())
        pcVUI.SetColourDescriptionPresentFlag(this.GetEncCfg().GetColourDescriptionPresentFlag())
        pcVUI.SetColourPrimaries(this.GetEncCfg().GetColourPrimaries())
        pcVUI.SetTransferCharacteristics(this.GetEncCfg().GetTransferCharacteristics())
        pcVUI.SetMatrixCoefficients(this.GetEncCfg().GetMatrixCoefficients())
        pcVUI.SetChromaLocInfoPresentFlag(this.GetEncCfg().GetChromaLocInfoPresentFlag())
        pcVUI.SetChromaSampleLocTypeTopField(this.GetEncCfg().GetChromaSampleLocTypeTopField())
        pcVUI.SetChromaSampleLocTypeBottomField(this.GetEncCfg().GetChromaSampleLocTypeBottomField())
        pcVUI.SetNeutralChromaIndicationFlag(this.GetEncCfg().GetNeutralChromaIndicationFlag())
        pcVUI.SetDefaultDisplayWindow(this.GetEncCfg().GetDefaultDisplayWindow())
        pcVUI.SetFrameFieldInfoPresentFlag(this.GetEncCfg().GetFrameFieldInfoPresentFlag())
        pcVUI.SetFieldSeqFlag(false)
        pcVUI.SetHrdParametersPresentFlag(false)
        //#if L0043_TIMING_INFO
        pcVUI.GetTimingInfo().SetPocProportionalToTimingFlag(this.GetEncCfg().GetPocProportionalToTimingFlag())
        pcVUI.GetTimingInfo().SetNumTicksPocDiffOneMinus1(this.GetEncCfg().GetNumTicksPocDiffOneMinus1())
        //#else
        //    pcVUI.SetPocProportionalToTimingFlag(this.GetEncCfg().GetPocProportionalToTimingFlag());
        //    pcVUI.SetNumTicksPocDiffOneMinus1   (this.GetEncCfg().GetNumTicksPocDiffOneMinus1()   );
        //#endif
        pcVUI.SetBitstreamRestrictionFlag(this.GetEncCfg().GetBitstreamRestrictionFlag())
        pcVUI.SetTilesFixedStructureFlag(this.GetEncCfg().GetTilesFixedStructureFlag())
        pcVUI.SetMotionVectorsOverPicBoundariesFlag(this.GetEncCfg().GetMotionVectorsOverPicBoundariesFlag())
        pcVUI.SetMinSpatialSegmentationIdc(this.GetEncCfg().GetMinSpatialSegmentationIdc())
        pcVUI.SetMaxBytesPerPicDenom(this.GetEncCfg().GetMaxBytesPerPicDenom())
        pcVUI.SetMaxBitsPerMinCuDenom(this.GetEncCfg().GetMaxBitsPerMinCuDenom())
        pcVUI.SetLog2MaxMvLengthHorizontal(this.GetEncCfg().GetLog2MaxMvLengthHorizontal())
        pcVUI.SetLog2MaxMvLengthVertical(this.GetEncCfg().GetLog2MaxMvLengthVertical())
    }
}

func (this *TEncTop) xInitPPS() { ///< initialize PPS from encoder options
    this.m_cPPS.SetConstrainedIntraPred(this.GetEncCfg().m_bUseConstrainedIntraPred)
    bUseDQP := (this.GetEncCfg().GetMaxCuDQPDepth() > 0)

    lowestQP := -this.m_cSPS.GetQpBDOffsetY()

    if this.GetEncCfg().GetUseLossless() {
        if (this.GetEncCfg().GetMaxCuDQPDepth() == 0) && (this.GetEncCfg().GetMaxDeltaQP() == 0) && (this.GetEncCfg().GetQP() == lowestQP) {
            bUseDQP = false
        } else {
            bUseDQP = true
        }
    } else {
        if bUseDQP == false {
            if (this.GetEncCfg().GetMaxDeltaQP() != 0) || this.GetEncCfg().GetUseAdaptiveQP() {
                bUseDQP = true
            }
        }
    }

    if bUseDQP {
        this.m_cPPS.SetUseDQP(true)
        this.m_cPPS.SetMaxCuDQPDepth(uint(this.GetEncCfg().m_iMaxCuDQPDepth))
        this.m_cPPS.SetMinCuDQPSize(this.m_cPPS.GetSPS().GetMaxCUWidth() >> (this.m_cPPS.GetMaxCuDQPDepth()))
    } else {
        this.m_cPPS.SetUseDQP(false)
        this.m_cPPS.SetMaxCuDQPDepth(0)
        this.m_cPPS.SetMinCuDQPSize(this.m_cPPS.GetSPS().GetMaxCUWidth() >> (this.m_cPPS.GetMaxCuDQPDepth()))
    }

    //#if RATE_CONTROL_LAMBDA_DOMAIN
    if this.GetEncCfg().m_RCEnableRateControl {
        this.m_cPPS.SetUseDQP(true)
        this.m_cPPS.SetMaxCuDQPDepth(0)
        this.m_cPPS.SetMinCuDQPSize(this.m_cPPS.GetSPS().GetMaxCUWidth() >> (this.m_cPPS.GetMaxCuDQPDepth()))
    }
    //#endif

    this.m_cPPS.SetChromaCbQpOffset(this.GetEncCfg().m_chromaCbQpOffset)
    this.m_cPPS.SetChromaCrQpOffset(this.GetEncCfg().m_chromaCrQpOffset)

    this.m_cPPS.SetNumSubstreams(this.GetEncCfg().m_iWaveFrontSubstreams)
    this.m_cPPS.SetEntropyCodingSyncEnabledFlag(this.GetEncCfg().m_iWaveFrontSynchro > 0)
    this.m_cPPS.SetTilesEnabledFlag((this.GetEncCfg().m_iNumColumnsMinus1 > 0 || this.GetEncCfg().m_iNumRowsMinus1 > 0))
    this.m_cPPS.SetUseWP(this.GetEncCfg().m_useWeightedPred)
    this.m_cPPS.SetWPBiPred(this.GetEncCfg().m_useWeightedBiPred)
    this.m_cPPS.SetOutputFlagPresentFlag(false)
    this.m_cPPS.SetSignHideFlag(this.GetEncCfg().GetSignHideFlag() != 0)
    this.m_cPPS.SetDeblockingFilterControlPresentFlag(this.GetEncCfg().m_DeblockingFilterControlPresent)
    this.m_cPPS.SetLog2ParallelMergeLevelMinus2(this.GetEncCfg().m_log2ParallelMergeLevelMinus2)
    this.m_cPPS.SetCabacInitPresentFlag(TLibCommon.CABAC_INIT_PRESENT_FLAG != 0)
    this.m_cPPS.SetLoopFilterAcrossSlicesEnabledFlag(this.GetEncCfg().m_bLFCrossSliceBoundaryFlag)
    var histogram [TLibCommon.MAX_NUM_REF + 1]int
    for i := 0; i <= TLibCommon.MAX_NUM_REF; i++ {
        histogram[i] = 0
    }
    for i := 0; i < this.GetEncCfg().GetGOPSize(); i++ {
        //assert(this.GetEncCfg().GetGOPEntry(i).this.m_numRefPicsActive >= 0 && this.GetEncCfg().GetGOPEntry(i).this.m_numRefPicsActive <= MAX_NUM_REF);
        histogram[this.GetEncCfg().GetGOPEntry(i).m_numRefPicsActive]++
    }
    maxHist := -1
    bestPos := 0
    for i := 0; i <= TLibCommon.MAX_NUM_REF; i++ {
        if histogram[i] > maxHist {
            maxHist = histogram[i]
            bestPos = i
        }
    }
    this.m_cPPS.SetNumRefIdxL0DefaultActive(uint(bestPos))
    this.m_cPPS.SetNumRefIdxL1DefaultActive(uint(bestPos))
    this.m_cPPS.SetTransquantBypassEnableFlag(this.GetEncCfg().GetTransquantBypassEnableFlag())
    this.m_cPPS.SetUseTransformSkip(this.GetEncCfg().m_useTransformSkip)
    if this.GetEncCfg().m_sliceSegmentMode != 0 {
        this.m_cPPS.SetDependentSliceSegmentsEnabledFlag(true)
    }
    if this.m_cPPS.GetDependentSliceSegmentsEnabledFlag() {
        var NumCtx uint
        if this.m_cPPS.GetEntropyCodingSyncEnabledFlag() {
            NumCtx = 2
        } else {
            NumCtx = 1
        }
        this.m_cSliceEncoder.initCtxMem(NumCtx)
        for st := uint(0); st < NumCtx; st++ {
            //TEncSbac* ctx = NULL;
            ctx := NewTEncSbac()
            ctx.init(this.m_cBinCoderCABAC)
            this.m_cSliceEncoder.setCtxMem(ctx, int(st))
        }
    }
}

func (this *TEncTop) xInitPPSforTiles() {
    this.m_cPPS.SetUniformSpacingFlag(this.GetEncCfg().m_iUniformSpacingIdr != 0)
    this.m_cPPS.SetNumColumnsMinus1(this.GetEncCfg().m_iNumColumnsMinus1)
    this.m_cPPS.SetNumRowsMinus1(this.GetEncCfg().m_iNumRowsMinus1)
    if this.GetEncCfg().m_iUniformSpacingIdr == 0 {
        this.m_cPPS.SetColumnWidth(this.GetEncCfg().m_puiColumnWidth)
        this.m_cPPS.SetRowHeight(this.GetEncCfg().m_puiRowHeight)
    }
    this.m_cPPS.SetLoopFilterAcrossTilesEnabledFlag(this.GetEncCfg().m_loopFilterAcrossTilesEnabledFlag)

    // # substreams is "per tile" when tiles are independent.
    if this.GetEncCfg().m_iWaveFrontSynchro != 0 {
        this.m_cPPS.SetNumSubstreams(this.GetEncCfg().m_iWaveFrontSubstreams * (this.GetEncCfg().m_iNumColumnsMinus1 + 1))
    }
}

func (this *TEncTop) xInitRPS() { ///< initialize PPS from encoder options
    var rps *TLibCommon.TComReferencePictureSet

    this.m_cSPS.CreateRPSList(this.GetEncCfg().GetGOPSize() + this.GetEncCfg().m_extraRPSs)
    rpsList := this.m_cSPS.GetRPSList()
    
	//fmt.Printf("getGOPSize()%d+m_extraRPSs%d\n", this.GetEncCfg().GetGOPSize(), this.GetEncCfg().m_extraRPSs);
    for i := 0; i < this.GetEncCfg().GetGOPSize()+this.GetEncCfg().m_extraRPSs; i++ {
        ge := this.GetEncCfg().GetGOPEntry(i)
        rps = rpsList.GetReferencePictureSet(i)
        rps.SetNumberOfPictures(ge.m_numRefPics)
        rps.SetNumRefIdc(ge.m_numRefIdc)
        //fmt.Printf("(%d %d) ", ge.m_numRefPics, ge.m_numRefIdc);
        
        numNeg := 0
        numPos := 0
        for j := 0; j < ge.m_numRefPics; j++ {
            rps.SetDeltaPOC(j, ge.m_referencePics[j])
            rps.SetUsed(j, ge.m_usedByCurrPic[j])
            if ge.m_referencePics[j] > 0 {
                numPos++
            } else {
                numNeg++
            }
        }
        rps.SetNumberOfNegativePictures(numNeg)
        rps.SetNumberOfPositivePictures(numPos)

        // handle inter RPS intialization from the config file.
        //#if AUTO_INTER_RPS
        rps.SetInterRPSPrediction(ge.m_interRPSPrediction > 0) // not very clean, converting anything > 0 to true.
        rps.SetDeltaRIdxMinus1(0)                              // index to the Reference RPS is always the previous one.
        var RPSRef *TLibCommon.TComReferencePictureSet; 
        if (i-1 < 0){
      		//fmt.Printf("Warning: getReferencePictureSet(i-1):i-1<0\n");
      		RPSRef = nil;
   		}else{
        	RPSRef = rpsList.GetReferencePictureSet(i - 1)        // get the reference RPS
		}
		
        if ge.m_interRPSPrediction == 2 { // Automatic generation of the inter RPS idc based on the RIdx provided.
            deltaRPS := this.GetEncCfg().GetGOPEntry(i-1).m_POC - ge.m_POC // the ref POC - current POC
            numRefDeltaPOC := RPSRef.GetNumberOfPictures()

            rps.SetDeltaRPS(deltaRPS)            // set delta RPS
            rps.SetNumRefIdc(numRefDeltaPOC + 1) // set the numRefIdc to the number of pictures in the reference RPS + 1.
            count := 0
            for j := 0; j <= numRefDeltaPOC; j++ { // cycle through pics in reference RPS.
                var RefDeltaPOC int
                if j < numRefDeltaPOC {
                    RefDeltaPOC = RPSRef.GetDeltaPOC(j) // if it is the last decoded picture, set RefDeltaPOC = 0
                } else {
                    RefDeltaPOC = 0
                }

                rps.SetRefIdc(j, 0)
                for k := 0; k < rps.GetNumberOfPictures(); k++ { // cycle through pics in current RPS.
                    if rps.GetDeltaPOC(k) == (RefDeltaPOC + deltaRPS) { // if the current RPS has a same picture as the reference RPS.
                        if rps.GetUsed(k) {
                            rps.SetRefIdc(j, 1)
                        } else {
                            rps.SetRefIdc(j, 2)
                        }
                        count++
                        break
                    }
                }
            }
            if count != rps.GetNumberOfPictures() {
                fmt.Printf("Warning: Unable fully predict all delta POCs using the reference RPS index given in the config file.  Setting Inter RPS to false for this RPS.\n")
                rps.SetInterRPSPrediction(false)
            }
        } else if ge.m_interRPSPrediction == 1 { // inter RPS idc based on the RefIdc values provided in config file.
            rps.SetDeltaRPS(ge.m_deltaRPS)
            rps.SetNumRefIdc(ge.m_numRefIdc)
            for j := 0; j < ge.m_numRefIdc; j++ {
                rps.SetRefIdc(j, ge.m_refIdc[j])
            }
            //#if WRITE_BACK
            // the folowing code overwrite the deltaPOC and Used by current values read from the config file with the ones
            // computed from the RefIdc.  A warning is printed if they are not identical.
            numNeg = 0
            numPos = 0
            var RPSTemp TLibCommon.TComReferencePictureSet // temporary variable

            for j := 0; j < ge.m_numRefIdc; j++ {
                if ge.m_refIdc[j] != 0 {
                    var deltaPOC int
                    if j < RPSRef.GetNumberOfPictures() {
                        deltaPOC = ge.m_deltaRPS + RPSRef.GetDeltaPOC(j)
                    } else {
                        deltaPOC = ge.m_deltaRPS + 0
                    }

                    RPSTemp.SetDeltaPOC((numNeg + numPos), deltaPOC)
                    RPSTemp.SetUsed((numNeg + numPos), ge.m_refIdc[j] == 1)
                    if deltaPOC < 0 {
                        numNeg++
                    } else {
                        numPos++
                    }
                }
            }
            if numNeg != rps.GetNumberOfNegativePictures() {
                fmt.Printf("Warning: number of negative pictures in RPS is different between intra and inter RPS specified in the config file.\n")
                rps.SetNumberOfNegativePictures(numNeg)
                rps.SetNumberOfPositivePictures(numNeg + numPos)
            }
            if numPos != rps.GetNumberOfPositivePictures() {
                fmt.Printf("Warning: number of positive pictures in RPS is different between intra and inter RPS specified in the config file.\n")
                rps.SetNumberOfPositivePictures(numPos)
                rps.SetNumberOfPositivePictures(numNeg + numPos)
            }
            RPSTemp.SetNumberOfPictures(numNeg + numPos)
            RPSTemp.SetNumberOfNegativePictures(numNeg)
            RPSTemp.SortDeltaPOC() // sort the created delta POC before comparing
            // check if Delta POC and Used are the same
            // print warning if they are not.
            for j := 0; j < ge.m_numRefIdc; j++ {
                if RPSTemp.GetDeltaPOC(j) != rps.GetDeltaPOC(j) {
                    fmt.Printf("Warning: delta POC is different between intra RPS and inter RPS specified in the config file.\n")
                    rps.SetDeltaPOC(j, RPSTemp.GetDeltaPOC(j))
                }
                if RPSTemp.GetUsed(j) != rps.GetUsed(j) {
                    fmt.Printf("Warning: Used by Current in RPS is different between intra and inter RPS specified in the config file.\n")
                    rps.SetUsed(j, RPSTemp.GetUsed(j))
                }
            }
            //#endif
        }
        /*#else
            rps.SetInterRPSPrediction(ge.m_interRPSPrediction);
            if (ge.m_interRPSPrediction)
            {
              rps.SetDeltaRIdxMinus1(0);
              rps.SetDeltaRPS(ge.m_deltaRPS);
              rps.SetNumRefIdc(ge.m_numRefIdc);
              for (Int j = 0; j < ge.m_numRefIdc; j++ )
              {
                rps.SetRefIdc(j, ge.m_refIdc[j]);
              }
        #if WRITE_BACK
              // the folowing code overwrite the deltaPOC and Used by current values read from the config file with the ones
              // computed from the RefIdc.  This is not necessary if both are identical. Currently there is no check to see if they are identical.
              numNeg = 0;
              numPos = 0;
              TComReferencePictureSet*     RPSRef = this.m_RPSList.this.GetEncCfg().GetReferencePictureSet(i-1);

              for (Int j = 0; j < ge.m_numRefIdc; j++ )
              {
                if (ge.m_refIdc[j])
                {
                  Int deltaPOC = ge.m_deltaRPS + ((j < RPSRef.this.GetEncCfg().GetNumberOfPictures())? RPSRef.getDeltaPOC(j) : 0);
                  rps.SetDeltaPOC((numNeg+numPos),deltaPOC);
                  rps.SetUsed((numNeg+numPos),ge.m_refIdc[j]==1?1:0);
                  if (deltaPOC<0)
                  {
                    numNeg++;
                  }
                  else
                  {
                    numPos++;
                  }
                }
              }
              rps.SetNumberOfNegativePictures(numNeg);
              rps.SetNumberOfPositivePictures(numPos);
              rps.sortDeltaPOC();
        #endif
            }
        #endif //INTER_RPS_AUTO
        */
    }

}

func (this *TEncTop) selectReferencePictureSet(slice *TLibCommon.TComSlice, POCCurr, GOPid int) {
    slice.SetRPSidx(GOPid)

    for extraNum := this.GetEncCfg().m_iGOPSize; extraNum < this.GetEncCfg().m_extraRPSs+this.GetEncCfg().m_iGOPSize; extraNum++ {
        if this.GetEncCfg().m_uiIntraPeriod > 0 && this.GetEncCfg().GetDecodingRefreshType() > 0 {
            POCIndex := POCCurr % int(this.GetEncCfg().m_uiIntraPeriod)
            if POCIndex == 0 {
                POCIndex = int(this.GetEncCfg().m_uiIntraPeriod)
            }
            if POCIndex == this.GetEncCfg().m_GOPList[extraNum].m_POC {
                slice.SetRPSidx(extraNum)
            }
        } else {
            if POCCurr == this.GetEncCfg().m_GOPList[extraNum].m_POC {
                slice.SetRPSidx(extraNum)
            }
        }
    }

    slice.SetRPS(this.m_cSPS.GetRPSList().GetReferencePictureSet(slice.GetRPSidx()))
    slice.GetRPS().SetNumberOfPictures(slice.GetRPS().GetNumberOfNegativePictures() + slice.GetRPS().GetNumberOfPositivePictures())
}

// -------------------------------------------------------------------------------------------------------------------
// encoder function
// -------------------------------------------------------------------------------------------------------------------

/// encode several number of pictures until end-of-sequence
func (this *TEncTop) Encode(flush bool, pcPicYuvOrg *TLibCommon.TComPicYuv, rcListPicYuvRecOut *list.List, accessUnitsOut *AccessUnits, iNumEncoded *int) {
    if pcPicYuvOrg != nil {
        // get original YUV
        pcPicCurr := this.xGetNewPicBuffer()
        pcPicYuvOrg.CopyToPic(pcPicCurr.GetPicYuvOrg())

        // compute image characteristics
        if this.GetEncCfg().GetUseAdaptiveQP() {
            pcPicCurr.XPreanalyze()
        }
    }

    if this.m_iNumPicRcvd == 0 || (!flush && this.m_iPOCLast != 0 && this.m_iNumPicRcvd != this.GetEncCfg().m_iGOPSize && this.GetEncCfg().m_iGOPSize != 0) {
        *iNumEncoded = 0
        return
    }

    //#if RATE_CONTROL_LAMBDA_DOMAIN
    if this.GetEncCfg().m_RCEnableRateControl {
        this.m_cRateCtrl.initRCGOP(this.m_iNumPicRcvd)
    }
    //#endif

    // compress GOP
    this.m_cGOPEncoder.compressGOP(this.m_iPOCLast, this.m_iNumPicRcvd, this.m_cListPic, rcListPicYuvRecOut, accessUnitsOut)

    //#if RATE_CONTROL_LAMBDA_DOMAIN
    if this.GetEncCfg().m_RCEnableRateControl {
        this.m_cRateCtrl.destroyRCGOP()
    }
    //#endif

    *iNumEncoded = this.m_iNumPicRcvd
    this.m_iNumPicRcvd = 0
    this.m_uiNumAllPicCoded += uint(*iNumEncoded)
}

func (this *TEncTop) PrintSummary() {
    this.m_cGOPEncoder.printOutSummary(this.m_uiNumAllPicCoded)
}

func (this *TEncTop) GetEncCfg() *TEncCfg                       { return this.m_pcEncCfg }
func (this *TEncTop) SetEncCfg(pcEncCfg *TEncCfg)               { this.m_pcEncCfg = pcEncCfg }
func (this *TEncTop) getListPic() *list.List                    { return this.m_cListPic }
func (this *TEncTop) getPredSearch() *TEncSearch                { return this.m_cSearch }
func (this *TEncTop) getTrQuant() *TLibCommon.TComTrQuant       { return this.m_cTrQuant }
func (this *TEncTop) getLoopFilter() *TLibCommon.TComLoopFilter { return this.m_cLoopFilter }
func (this *TEncTop) getSAO() *TEncSampleAdaptiveOffset         { return this.m_cEncSAO }
func (this *TEncTop) getGOPEncoder() *TEncGOP                   { return this.m_cGOPEncoder }
func (this *TEncTop) getSliceEncoder() *TEncSlice               { return this.m_cSliceEncoder }
func (this *TEncTop) getCuEncoder() *TEncCu                     { return this.m_cCuEncoder }
func (this *TEncTop) getEntropyCoder() *TEncEntropy             { return this.m_cEntropyCoder }
func (this *TEncTop) getTraceFile()   io.Writer			    	{ return this.m_pTraceFile }
func (this *TEncTop) getCavlcCoder() *TEncCavlc                 { return this.m_cCavlcCoder }
func (this *TEncTop) getSbacCoder() *TEncSbac                   { return this.m_cSbacCoder }
func (this *TEncTop) getBinCABAC() *TEncBinCABAC                { return this.m_cBinCoderCABAC }
func (this *TEncTop) getSbacCoders() []*TEncSbac                { return this.m_pcSbacCoders }
func (this *TEncTop) getBinCABACs() []*TEncBinCABAC             { return this.m_pcBinCoderCABACs }
func (this *TEncTop) getBitCounter() *TLibCommon.TComBitCounter { return this.m_cBitCounter }
func (this *TEncTop) getRdCost() *TEncRdCost                    { return this.m_cRdCost }
func (this *TEncTop) getRDSbacCoder() [][]*TEncSbac             { return this.m_pppcRDSbacCoder }
func (this *TEncTop) getRDGoOnSbacCoder() *TEncSbac             { return this.m_cRDGoOnSbacCoder }
func (this *TEncTop) getBitCounters() []*TLibCommon.TComBitCounter {
    return this.m_pcBitCounters
}
func (this *TEncTop) getRdCosts() []*TEncRdCost                   { return this.m_pcRdCosts }
func (this *TEncTop) getRDSbacCoders() [][][]*TEncSbac            { return this.m_ppppcRDSbacCoders }
func (this *TEncTop) getRDGoOnSbacCoders() []*TEncSbac            { return this.m_pcRDGoOnSbacCoders }
func (this *TEncTop) getRateCtrl() *TEncRateCtrl                  { return this.m_cRateCtrl }
func (this *TEncTop) getSPS() *TLibCommon.TComSPS                 { return this.m_cSPS }
func (this *TEncTop) getPPS() *TLibCommon.TComPPS                 { return this.m_cPPS }
func (this *TEncTop) getScalingList() *TLibCommon.TComScalingList { return this.m_scalingList }
