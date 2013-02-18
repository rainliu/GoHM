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
)

// ====================================================================================================================
// Class definition
// ====================================================================================================================
// Move to TEncPic struct
/*
/// Source picture analyzer class
type TEncPreanalyzer struct{
public:
  TEncPreanalyzer();
  virtual ~TEncPreanalyzer();

  Void xPreanalyze( TEncPic* pcPic );
};
*/
// ====================================================================================================================
// Class definition
// ====================================================================================================================

/// encoder analyzer class
type TEncAnalyze struct{
  m_dPSNRSumY		float64;
  m_dPSNRSumU		float64;
  m_dPSNRSumV		float64;
  m_dAddBits		float64;
  m_uiNumPic		uint;
  m_dFrmRate		float64; //--CFG_KDY
}
  
func NewTEncAnalyze() *TEncAnalyze { 
	return &TEncAnalyze{};
	//m_dPSNRSumY = m_dPSNRSumU = m_dPSNRSumV = m_dAddBits = m_uiNumPic = 0;  
}
  
func (this *TEncAnalyze)  addResult(  psnrY,  psnrU,  psnrV,  bits float64) {
    this.m_dPSNRSumY += psnrY;
    this.m_dPSNRSumU += psnrU;
    this.m_dPSNRSumV += psnrV;
    this.m_dAddBits  += bits;
    
    this.m_uiNumPic++;
}
  
func (this *TEncAnalyze)  getPsnrY() float64 { return  this.m_dPSNRSumY;  }
func (this *TEncAnalyze)  getPsnrU() float64 { return  this.m_dPSNRSumU;  }
func (this *TEncAnalyze)  getPsnrV() float64 { return  this.m_dPSNRSumV;  }
func (this *TEncAnalyze)  getBits()  float64 { return  this.m_dAddBits;   }
func (this *TEncAnalyze)  getNumPic() uint   { return  this.m_uiNumPic;   }
  
func (this *TEncAnalyze)  setFrmRate  ( dFrameRate float64) { this.m_dFrmRate = dFrameRate; } //--CFG_KDY
func (this *TEncAnalyze)  clear() { 
	this.m_dPSNRSumY = 0;
	this.m_dPSNRSumU = 0;
	this.m_dPSNRSumV = 0;
	this.m_dAddBits  = 0;
	this.m_uiNumPic = 0;  
}
func (this *TEncAnalyze)  printOut ( cDelim string){
    dFps     :=   this.m_dFrmRate; //--CFG_KDY
    dScale   := dFps / 1000 / float64(this.m_uiNumPic);
    
    fmt.Printf( "\tTotal Frames |     Bitrate      Y-PSNR      U-PSNR      V-PSNR \n" );
    //printf( "\t------------ "  " ----------"   " -------- "  " -------- "  " --------\n" );
    fmt.Printf( "\t %8d    %s          %12.4f      %8.4f     %8.4f      %8.4f\n",
           this.getNumPic(), cDelim,
           this.getBits() * dScale,
           this.getPsnrY() / float64(this.getNumPic()),
           this.getPsnrU() / float64(this.getNumPic()),
           this.getPsnrV() / float64(this.getNumPic()) );
  }
  
func (this *TEncAnalyze)  printSummaryOut (){
    /*FILE* pFile = fopen ("summaryTotal.txt", "at");
    Double dFps     =   this.m_dFrmRate; //--CFG_KDY
    Double dScale   = dFps / 1000 / (Double)this.m_uiNumPic;
    
    fprintf(pFile, "%f\t %f\t %f\t %f\n", getBits() * dScale,
            getPsnrY() / (Double)getNumPic(),
            getPsnrU() / (Double)getNumPic(),
            getPsnrV() / (Double)getNumPic() );
    fclose(pFile);*/
}
  
func (this *TEncAnalyze)  printSummary(ch string) {
    /*FILE* pFile = NULL;
    
    switch( ch ) 
    {
      case 'I':
        pFile = fopen ("summary_I.txt", "at");
        break;
      case 'P':
        pFile = fopen ("summary_P.txt", "at");
        break;
      case 'B':
        pFile = fopen ("summary_B.txt", "at");
        break;
      default:
        assert(0);
        return;
        break;
    }
    
    Double dFps     =   this.m_dFrmRate; //--CFG_KDY
    Double dScale   = dFps / 1000 / (Double)this.m_uiNumPic;
    
    fprintf(pFile, "%f\t %f\t %f\t %f\n",
            getBits() * dScale,
            getPsnrY() / (Double)getNumPic(),
            getPsnrU() / (Double)getNumPic(),
            getPsnrV() / (Double)getNumPic() );
    
    fclose(pFile);
  }*/
};