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
	"gohm/TLibCommon"
)

// ====================================================================================================================
// Class definition
// ====================================================================================================================

/// Unit block for storing image characteristics
type TEncQPAdaptationUnit struct{
  m_dActivity	float64;
}


func NewTEncQPAdaptationUnit() *TEncQPAdaptationUnit{
	return &TEncQPAdaptationUnit{};
}

func (this *TEncQPAdaptationUnit) setActivity( d float64) { this.m_dActivity = d; }
func (this *TEncQPAdaptationUnit) getActivity() float64   { return this.m_dActivity; }


/// Local image characteristics for CUs on a specific depth
type TEncPicQPAdaptationLayer struct{
  m_uiAQPartWidth			uint;
  m_uiAQPartHeight			uint;
  m_uiNumAQPartInWidth		uint;
  m_uiNumAQPartInHeight		uint;
  m_acTEncAQU				[]TEncQPAdaptationUnit;
  m_dAvgActivity			float64;
}

func NewTEncPicQPAdaptationLayer() *TEncPicQPAdaptationLayer{
	return &TEncPicQPAdaptationLayer{};
}

func (this *TEncPicQPAdaptationLayer)  create( iWidth, iHeight int, uiAQPartWidth, uiAQPartHeight uint){
  this.m_uiAQPartWidth = uiAQPartWidth;
  this.m_uiAQPartHeight = uiAQPartHeight;
  this.m_uiNumAQPartInWidth  = (uint(iWidth)  + this.m_uiAQPartWidth -1) / this.m_uiAQPartWidth;
  this.m_uiNumAQPartInHeight = (uint(iHeight) + this.m_uiAQPartHeight-1) / this.m_uiAQPartHeight;
  this.m_acTEncAQU = make([]TEncQPAdaptationUnit, this.m_uiNumAQPartInWidth * this.m_uiNumAQPartInHeight );
}
func (this *TEncPicQPAdaptationLayer)  destroy(){
	//do nothing
}

func (this *TEncPicQPAdaptationLayer)  getAQPartWidth()        uint{ return this.m_uiAQPartWidth;       }
func (this *TEncPicQPAdaptationLayer)  getAQPartHeight()       uint{ return this.m_uiAQPartHeight;      }
func (this *TEncPicQPAdaptationLayer)  getNumAQPartInWidth()   uint{ return this.m_uiNumAQPartInWidth;  }
func (this *TEncPicQPAdaptationLayer)  getNumAQPartInHeight()  uint{ return this.m_uiNumAQPartInHeight; }
func (this *TEncPicQPAdaptationLayer)  getAQPartStride()       uint{ return this.m_uiNumAQPartInWidth;  }
func (this *TEncPicQPAdaptationLayer)  getQPAdaptationUnit()   []TEncQPAdaptationUnit{ return this.m_acTEncAQU;           }
func (this *TEncPicQPAdaptationLayer)  getAvgActivity()        float64{ return this.m_dAvgActivity;        }
func (this *TEncPicQPAdaptationLayer)  setAvgActivity( d float64)  { this.m_dAvgActivity = d; }


/// Picture class including local image characteristics information for QP adaptation
type TEncPic struct{
  TLibCommon.TComPic
  m_acAQLayer		[]TEncPicQPAdaptationLayer;
  m_uiMaxAQDepth	uint;
}
func NewTEncPic() *TEncPic{
	return &TEncPic{};
}

func (this *TEncPic)  create( iWidth, iHeight int, uiMaxWidth, uiMaxHeight, uiMaxDepth, uiMaxAQDepth uint,   
                        	  picCroppingWindow *TLibCommon.CroppingWindow, numReorderPics []int, bIsVirtual bool ){//= false ){
  this.TComPic.Create( iWidth, iHeight, uiMaxWidth, uiMaxHeight, uiMaxDepth, picCroppingWindow, numReorderPics, bIsVirtual );
  this.m_uiMaxAQDepth = uiMaxAQDepth;
  if uiMaxAQDepth > 0 {
    this.m_acAQLayer = make([]TEncPicQPAdaptationLayer, this.m_uiMaxAQDepth ); 
    for d := uint(0); d < this.m_uiMaxAQDepth; d++ {
      this.m_acAQLayer[d].create( iWidth, iHeight, uiMaxWidth>>d, uiMaxHeight>>d );
    }
  }	
}
func (this *TEncPic)  destroy(){
	//do nothing
}

func (this *TEncPic)  getAQLayer( uiDepth uint) *TEncPicQPAdaptationLayer  { return &this.m_acAQLayer[uiDepth]; }

func (this *TEncPic)  getMaxAQDepth()     uint        { return this.m_uiMaxAQDepth;        }

func (this *TEncPic)  xPreanalyze(){
  pcPicYuv := this.GetPicYuvOrg();
  iWidth := pcPicYuv.GetWidth();
  iHeight := pcPicYuv.GetHeight();
  iStride := pcPicYuv.GetStride();

  for d := uint(0); d < this.getMaxAQDepth(); d++ {
    pLineY := pcPicYuv.GetLumaAddr();
    pcAQLayer := this.getAQLayer(d);
    uiAQPartWidth := pcAQLayer.getAQPartWidth();
    uiAQPartHeight := pcAQLayer.getAQPartHeight();
    pcAQU := pcAQLayer.getQPAdaptationUnit();

    dSumAct := float64(0.0);
    for y := 0; y < iHeight; y += int(uiAQPartHeight) {
      uiCurrAQPartHeight := uint(TLibCommon.MIN(int(uiAQPartHeight), int(iHeight-y)).(int));
      for x := 0; x < iWidth; x += int(uiAQPartWidth) {
        uiCurrAQPartWidth := uint(TLibCommon.MIN(int(uiAQPartWidth), int(iWidth-x)).(int));
        pBlkY := pLineY[x:];
        var uiSum	=[4]uint{0, 0, 0, 0};
        var uiSumSq=[4]uint{0, 0, 0, 0};
        uiNumPixInAQPart := uint(0);
        by := uint(0);
        for  ; by < uiCurrAQPartHeight>>1; by++ {
          bx := uint(0);
          for  ; bx < uiCurrAQPartWidth>>1; bx++{
            uiSum  [0] += uint(pBlkY[bx]);
            uiSumSq[0] += uint(pBlkY[bx]) * uint(pBlkY[bx]);
            uiNumPixInAQPart++ ;
          }
          for  ; bx < uiCurrAQPartWidth; bx++ {
            uiSum  [1] += uint(pBlkY[bx]);
            uiSumSq[1] += uint(pBlkY[bx]) * uint(pBlkY[bx]);
            uiNumPixInAQPart++;
          }
          pBlkY = pBlkY[iStride:]//+= iStride;
        }
        for  ; by < uiCurrAQPartHeight; by++ {
          bx := uint(0);
          for  ; bx < uiCurrAQPartWidth>>1; bx++ {
            uiSum  [2] += uint(pBlkY[bx]);
            uiSumSq[2] += uint(pBlkY[bx]) * uint(pBlkY[bx]);
            uiNumPixInAQPart++;
          }
          for ; bx < uiCurrAQPartWidth; bx++ {
            uiSum  [3] += uint(pBlkY[bx]);
            uiSumSq[3] += uint(pBlkY[bx]) * uint(pBlkY[bx]);
            uiNumPixInAQPart++;
          }
          pBlkY = pBlkY[iStride:]//+= iStride;
        }

        dMinVar := float64(TLibCommon.MAX_DOUBLE);
        for i:=int(0); i<4; i++ {
          dAverage := float64(uiSum[i]) / float64(uiNumPixInAQPart);
          dVariance := float64(uiSumSq[i]) / float64(uiNumPixInAQPart) - dAverage * dAverage;
          if dMinVar > dVariance {
          	dMinVar = dVariance;
          }
        }
        dActivity := 1.0 + dMinVar;
        pcAQU[0].setActivity( dActivity );
        dSumAct += dActivity;
        
        pcAQU = pcAQU[1:]//++
      }
      pLineY =pLineY[uint(iStride) * uiCurrAQPartHeight:]//+= iStride * uiCurrAQPartHeight;
    }

    dAvgAct := dSumAct / float64(pcAQLayer.getNumAQPartInWidth() * pcAQLayer.getNumAQPartInHeight());
    pcAQLayer.setAvgActivity( dAvgAct );
  }
}
