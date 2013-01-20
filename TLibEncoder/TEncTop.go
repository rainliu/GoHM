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
	"container/list"
	"gohm/TLibCommon"
)

// ====================================================================================================================
// Class definition
// ====================================================================================================================

type TEncTop struct{
  m_pcEncCfg	*TEncCfg
	
  // picture
  m_iPOCLast		int;                     ///< time index (POC)
  m_iNumPicRcvd		int;                  ///< number of received pictures
  m_uiNumAllPicCoded	uint;             ///< number of coded pictures
  m_cListPic	*list.List;                     ///< dynamic list of pictures
  
  // encoder search
  m_cSearch			*TEncSearch;                      ///< encoder search class
  m_pcEntropyCoder	*TEncEntropy;                     ///< entropy encoder 
  m_pcCavlcCoder	*TEncCavlc;                       ///< CAVLC encoder  
  // coding tool
  m_cTrQuant		*TLibCommon.TComTrQuant;                     ///< transform & quantization class
  m_cLoopFilter		*TLibCommon.TComLoopFilter;                  ///< deblocking filter class
  m_cEncSAO		*TEncSampleAdaptiveOffset;                     ///< sample adaptive offset class
  m_cEntropyCoder	*TEncEntropy;                ///< entropy encoder
  m_cCavlcCoder		*TEncCavlc;                  ///< CAVLC encoder
  m_cSbacCoder		*TEncSbac;                   ///< SBAC encoder
  m_cBinCoderCABAC	*TEncBinCABAC;               ///< bin coder CABAC
  m_pcSbacCoders		[]*TEncSbac;                 ///< SBAC encoders (to encode substreams )
  m_pcBinCoderCABACs	[]*TEncBinCABAC;             ///< bin coders CABAC (one per substream)
  
  // processing unit
  m_cGOPEncoder		*TEncGOP;                  ///< GOP encoder
  m_cSliceEncoder	*TEncSlice;                ///< slice encoder
  m_cCuEncoder		*TEncCu;                   ///< CU encoder
  // SPS
  m_cSPS			*TLibCommon.TComSPS;                         ///< SPS
  m_cPPS			*TLibCommon.TComPPS;                         ///< PPS
  // RD cost computation
  m_cBitCounter		*TLibCommon.TComBitCounter;                  ///< bit counter for RD optimization
  m_cRdCost			*TEncRdCost;                      ///< RD cost computation class
  m_pppcRDSbacCoder	[][]*TEncSbac;              ///< temporal storage for RD computation
  m_cRDGoOnSbacCoder	*TEncSbac;             ///< going on SBAC model for RD stage
//#if FAST_BIT_EST
  m_pppcBinCoderCABAC	[][]*TEncBinCABACCounter;            ///< temporal CABAC state storage for RD computation
  m_cRDGoOnBinCoderCABAC	*TEncBinCABACCounter;         ///< going on bin coder CABAC for RD stage
//#else
//  TEncBinCABAC***         m_pppcBinCoderCABAC;            ///< temporal CABAC state storage for RD computation
//  TEncBinCABAC            m_cRDGoOnBinCoderCABAC;         ///< going on bin coder CABAC for RD stage
//#endif
  m_iNumSubstreams	int;                ///< # of top-level elements allocated.
  m_pcBitCounters	*TLibCommon.TComBitCounter;                 ///< bit counters for RD optimization per substream
  m_pcRdCosts		*TEncRdCost;                     ///< RD cost computation class per substream
  m_ppppcRDSbacCoders	[][][][]TEncSbac;             ///< temporal storage for RD computation per substream
  m_pcRDGoOnSbacCoders	*TEncSbac;            ///< going on SBAC model for RD stage per substream
  m_ppppcBinCodersCABAC	[][][][]TEncBinCABAC;           ///< temporal CABAC state storage for RD computation per substream
  m_pcRDGoOnBinCodersCABAC	*TEncBinCABAC;        ///< going on bin coder CABAC for RD stage per substream

  // quality control
  m_cPreanalyzer	*TLibCommon.TComPic;                 ///< image characteristics analyzer for TM5-step3-like adaptive QP

  m_scalingList		*TLibCommon.TComScalingList;                 ///< quantization matrix information
  m_cRateCtrl		*TEncRateCtrl;                    ///< Rate control class 
}

func NewTEncTop() *TEncTop{
	return &TEncTop{m_iPOCLast:-1}
}
 
func (this *TEncTop)      Create          (){
}
func (this *TEncTop)      Destroy         (){}
func (this *TEncTop)      Init            (){
}
func (this *TEncTop)      DeletePicBuffer (){}

func (this *TEncTop)      CreateWPPCoders( iNumSubstreams int){
}
  
  // -------------------------------------------------------------------------------------------------------------------
  // member access functions
  // -------------------------------------------------------------------------------------------------------------------
func (this *TEncTop)  xGetNewPicBuffer  ( rpcPic *TLibCommon.TComPic){           ///< get picture buffer which will be processed
}

func (this *TEncTop)  xInitSPS          (){                             ///< initialize SPS from encoder options
}
func (this *TEncTop)  xInitPPS          (){                             ///< initialize PPS from encoder options
}  
func (this *TEncTop)  xInitPPSforTiles  (){
}
func (this *TEncTop)  xInitRPS          (){                             ///< initialize PPS from encoder options
}
func (this *TEncTop)  GetEncCfg             () *TEncCfg		  {return this.m_pcEncCfg; }
func (this *TEncTop)  SetEncCfg             (pcEncCfg *TEncCfg)		  { this.m_pcEncCfg = pcEncCfg; }
func (this *TEncTop)  getListPic            () *list.List     { return  this.m_cListPic;             }
func (this *TEncTop)  getPredSearch         () *TEncSearch             { return  this.m_cSearch;              }
  
func (this *TEncTop)  getTrQuant            () *TLibCommon.TComTrQuant           { return  this.m_cTrQuant;             }
func (this *TEncTop)  getLoopFilter         () *TLibCommon.TComLoopFilter         { return  this.m_cLoopFilter;          }
func (this *TEncTop)  getSAO                () *TEncSampleAdaptiveOffset{ return  this.m_cEncSAO;              }
func (this *TEncTop)  getGOPEncoder         () *TEncGOP               { return  this.m_cGOPEncoder;          }
func (this *TEncTop)  getSliceEncoder       () *TEncSlice              { return  this.m_cSliceEncoder;        }
func (this *TEncTop)  getCuEncoder          () *TEncCu                 { return  this.m_cCuEncoder;           }
func (this *TEncTop)  getEntropyCoder       () *TEncEntropy            { return  this.m_cEntropyCoder;        }
func (this *TEncTop)  getCavlcCoder         () *TEncCavlc              { return  this.m_cCavlcCoder;          }
func (this *TEncTop)  getSbacCoder          () *TEncSbac               { return  this.m_cSbacCoder;           }
func (this *TEncTop)  getBinCABAC           () *TEncBinCABAC           { return  this.m_cBinCoderCABAC;       }
func (this *TEncTop)  getSbacCoders         () []*TEncSbac               { return  this.m_pcSbacCoders;      }
func (this *TEncTop)  getBinCABACs          () []*TEncBinCABAC           { return  this.m_pcBinCoderCABACs;      }
  
func (this *TEncTop)  getBitCounter         () *TLibCommon.TComBitCounter         { return  this.m_cBitCounter;          }
func (this *TEncTop)  getRdCost             () *TEncRdCost             { return  this.m_cRdCost;              }
func (this *TEncTop)  getRDSbacCoder        () [][]*TEncSbac             { return  this.m_pppcRDSbacCoder;       }
func (this *TEncTop)  getRDGoOnSbacCoder    () *TEncSbac               { return  this.m_cRDGoOnSbacCoder;     }
func (this *TEncTop)  getBitCounters        () *TLibCommon.TComBitCounter         { return  this.m_pcBitCounters;         }
func (this *TEncTop)  getRdCosts            () *TEncRdCost             { return  this.m_pcRdCosts;             }
func (this *TEncTop)  getRDSbacCoders       () [][][][]TEncSbac            { return  this.m_ppppcRDSbacCoders;     }
func (this *TEncTop)  getRDGoOnSbacCoders   () *TEncSbac               { return  this.m_pcRDGoOnSbacCoders;   }
func (this *TEncTop)  getRateCtrl           () *TEncRateCtrl           { return this.m_cRateCtrl;             }
func (this *TEncTop)  getSPS                () *TLibCommon.TComSPS                { return  this.m_cSPS;                 }
func (this *TEncTop)  getPPS                () *TLibCommon.TComPPS                { return  this.m_cPPS;                 }
func (this *TEncTop)  selectReferencePictureSet(slice *TLibCommon.TComSlice, POCCurr, GOPid int, listPic *list.List){
}
func (this *TEncTop)  getScalingList        () *TLibCommon.TComScalingList        { return  this.m_scalingList;         }
  // -------------------------------------------------------------------------------------------------------------------
  // encoder function
  // -------------------------------------------------------------------------------------------------------------------

  /// encode several number of pictures until end-of-sequence
func (this *TEncTop) Encode(  bEos bool, pcPicYuvOrg *TLibCommon.TComPicYuv, rcListPicYuvRecOut *list.List, accessUnitsOut *AccessUnits, iNumEncoded *int ){
}  

func (this *TEncTop) PrintSummary() { 
	this.m_cGOPEncoder.printOutSummary (this.m_uiNumAllPicCoded); 
}
