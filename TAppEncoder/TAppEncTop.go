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
	"io"
	"container/list"
	"gohm/TLibCommon"
	"gohm/TLibEncoder"
)

// ====================================================================================================================
// Class definition
// ====================================================================================================================

/// encoder application class
type TAppEncTop struct{
	TAppEncCfg

  // class interface
  m_cTEncTop				*TLibEncoder.TEncTop;                    ///< encoder class
  m_cTVideoIOYuvInputFile	*TLibCommon.TVideoIOYuv;       ///< input YUV file
  m_cTVideoIOYuvReconFile	*TLibCommon.TVideoIOYuv;       ///< output reconstruction file
  
  m_cListPicYuvRec			*list.List;              ///< list of reconstruction YUV files TComList<TComPicYuv*>      
  
  m_iFrameRcvd				int;                  ///< number of received frames
  
  m_essentialBytes			uint;
  m_totalBytes				uint;
}

func NewTAppEncTop() *TAppEncTop{
	return &TAppEncTop{};
}
  
func (this *TAppEncTop) Encode      (){                               ///< main encoding function
}
/*
func (this *TAppEncTop) GetTEncTop  () *TEncTop{ 
	return  &m_cTEncTop; 
}      ///< return encoder class pointer reference
*/
//protected:
  // initialization
func (this *TAppEncTop)  xCreateLib        (){                               ///< create files & encoder class
}
func (this *TAppEncTop)  xInitLibCfg       (){                               ///< initialize internal variables
}
func (this *TAppEncTop)  xInitLib          (){                               ///< initialize encoder class
}
func (this *TAppEncTop)  xDestroyLib       (){                               ///< destroy encoder class
}
  
  /// obtain required buffers
func (this *TAppEncTop)  xGetBuffer(rpcPicYuvRec *TLibCommon.TComPicYuv){
}
  
  /// delete allocated buffers
func (this *TAppEncTop)  xDeleteBuffer     (){
}
  
  // file I/O
func (this *TAppEncTop)  xWriteOutput(bitstreamFile io.Writer, iNumEncoded int, accessUnits *list.List) { //const std::list<AccessUnit>& ///< write bitstream to file
}

//func (this *TAppEncTop)  rateStatsAccum(au *AccessUnit, stats *list.List){//const std::vector<UInt>
//}

func (this *TAppEncTop)  printRateSummary(){
}


