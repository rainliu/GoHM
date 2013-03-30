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

package TAppDecoder

import (
    "fmt"
    "os"
    //"errors"
    "container/list"
    "gohm/TLibCommon"
    "gohm/TLibDecoder"
)

type TAppDecTop struct {
    TAppDecCfg

    m_cTDecTop              *TLibDecoder.TDecTop
    m_cTVideoIOYuvReconFile *TLibCommon.TVideoIOYuv

    m_abDecFlag       [TLibCommon.MAX_GOP]bool
    m_iPOCLastDisplay int
}

func NewTAppDecTop() *TAppDecTop {
    pAppDecTop := &TAppDecTop{}
    //::memset (m_abDecFlag, 0, sizeof (m_abDecFlag));//memset 0 by Go
    pAppDecTop.m_iPOCLastDisplay = -TLibCommon.MAX_INT
    pAppDecTop.TAppDecCfg.m_targetDecLayerIdSet = list.New()
    pAppDecTop.m_cTDecTop = TLibDecoder.NewTDecTop()
    pAppDecTop.m_cTVideoIOYuvReconFile = TLibCommon.NewTVideoIOYuv()

    return pAppDecTop
}

func (this *TAppDecTop) Create() {
    //do nothing
}

func (this *TAppDecTop) Destroy() {
    //do nothing
}

func (this *TAppDecTop) Decode() (err error) {
	bSkipPictureForBLA := false;
    var poc int
    var pcListPic *list.List           // = NULL;
    var nalUnit, oldNalUnit *list.List //vector<uint8_t>
    var nalu TLibDecoder.InputNALUnit

    bitstreamFile, err := os.Open(this.m_pchBitstreamFile)
    if err != nil {
        fmt.Printf("\nfailed to open bitstream file `%s' for reading\n", this.m_pchBitstreamFile)
        return err
    }
    defer bitstreamFile.Close()

    bytestream := TLibDecoder.NewInputByteStream(bitstreamFile)

    // create & initialize internal classes
    this.xCreateDecLib()
    this.xInitDecLib()
    this.m_iPOCLastDisplay += this.m_iSkipFrame // set the last displayed POC correctly for skip forward.

    // main decoder loop
    recon_opened := false // reconstruction file not yet opened. (must be performed after SPS is seen)
    eof := false
    bNewPicture := false
    iDecodedFrameNum := 0
    for !eof || bNewPicture { // (!!bitstreamFile)
        /* location serves to work around a design fault in the decoder, whereby
         * the process of reading a new slice that is the first slice of a new frame
         * requires the TDecTop::decode() method to be called again with the same
         * nal unit. */
        //streampos location = bitstreamFile.tellg();
        var stats TLibDecoder.AnnexBStats // stats = AnnexBStats();
        bPreviousPictureDecoded := false

        if !bNewPicture {
            nalUnit = list.New() //vector<uint8_t>
            eof,_ = bytestream.ByteStreamNALUnit(nalUnit, &stats)
        } else {
            nalUnit = oldNalUnit
        }

        // call actual decoding function
        if nalUnit.Len() == 0 {
            /* this can happen if the following occur:
             *  - empty input file
             *  - two back-to-back start_code_prefixes
             *  - start_code_prefix immediately followed by EOF
             */
            fmt.Printf("Warning: Attempt to decode an empty NAL unit\n")
            break
        } else {
            //fmt.Printf("NalUnit Len=%d\n", nalUnit.Len())
            oldNalUnit = nalu.Read(nalUnit)

            //fmt.Printf("Type=%d\n", nalu.GetNalUnitType())

            if (this.m_iMaxTemporalLayer >= 0 && int(nalu.GetTemporalId()) > this.m_iMaxTemporalLayer) ||
                !this.IsNaluWithinTargetDecLayerIdSet(&nalu) {
                if bPreviousPictureDecoded {
                    bNewPicture = true
                    bPreviousPictureDecoded = false
                } else {
                    bNewPicture = false
                }
            } else {
                bNewPicture = this.m_cTDecTop.Decode(&nalu, &this.m_iSkipFrame, &this.m_iPOCLastDisplay, &bSkipPictureForBLA, !bNewPicture)
                bPreviousPictureDecoded = true
            }
        }
        if bNewPicture || eof {
            pcListPic = this.m_cTDecTop.ExecuteLoopFilters(&poc, bSkipPictureForBLA)
        }

        if pcListPic != nil {
            if this.m_pchReconFile != "" && !recon_opened {
                if this.m_outputBitDepthY == 0 {
                    this.m_outputBitDepthY = TLibCommon.G_bitDepthY
                }
                if this.m_outputBitDepthC == 0 {
                    this.m_outputBitDepthC = TLibCommon.G_bitDepthC
                }

                this.m_cTVideoIOYuvReconFile.Open(this.m_pchReconFile, true, this.m_outputBitDepthY, this.m_outputBitDepthC, TLibCommon.G_bitDepthY, TLibCommon.G_bitDepthC) // write mode
                recon_opened = true
            }
            //fmt.Printf("bNewPicture=%d m_nalUnitType=%d ", TLibCommon.B2U(bNewPicture), nalu.GetNalUnitType());
            if bNewPicture &&
                (nalu.GetNalUnitType() == TLibCommon.NAL_UNIT_CODED_SLICE_IDR ||
                    nalu.GetNalUnitType() == TLibCommon.NAL_UNIT_CODED_SLICE_IDR_N_LP ||
                    nalu.GetNalUnitType() == TLibCommon.NAL_UNIT_CODED_SLICE_BLA_N_LP ||
                    nalu.GetNalUnitType() == TLibCommon.NAL_UNIT_CODED_SLICE_BLANT ||
                    nalu.GetNalUnitType() == TLibCommon.NAL_UNIT_CODED_SLICE_BLA) {
                this.xFlushOutput(pcListPic)
            }
            // write reconstruction to file
            if bNewPicture {
                this.xWriteOutput(pcListPic, nalu.GetTemporalId())

                iDecodedFrameNum++
                if iDecodedFrameNum >= this.m_iFrameNum && this.m_iFrameNum > 0 {
                    break
                }
            }
        }
    }

    this.xFlushOutput(pcListPic)
    // delete buffers
    this.m_cTDecTop.DeletePicBuffer()

    // destroy internal classes
    this.xDestroyDecLib()

    return nil
}

func (this *TAppDecTop) xCreateDecLib() {
    //create decoder class
    this.m_cTDecTop.Create(this.m_pchTraceFile)
}

func (this *TAppDecTop) xDestroyDecLib() {
    if this.m_pchReconFile != "" {
        this.m_cTVideoIOYuvReconFile.Close()
    }

    //destroy decoder class
    this.m_cTDecTop.Destroy()
}

func (this *TAppDecTop) xInitDecLib() {
    //initialize decoder class
    this.m_cTDecTop.Init()
    this.m_cTDecTop.SetDecodedPictureHashSEIEnabled(this.m_decodedPictureHashSEIEnabled)

}

func (this *TAppDecTop) xWriteOutput(pcListPic *list.List, tId uint) {
    not_displayed := 0
	maxNumReorderPics := 0;
	
    for e := pcListPic.Front(); e != nil; e = e.Next() {
        pcPic := e.Value.(*TLibCommon.TComPic)
        if pcPic.GetOutputMark() && int(pcPic.GetPOC()) > this.m_iPOCLastDisplay {
        	for i:=uint(0); i<TLibCommon.MAX_TLAYER; i++ {
		        if pcPic.GetNumReorderPics(i)>maxNumReorderPics {
		          maxNumReorderPics = pcPic.GetNumReorderPics(i);
		        }
		    }
            not_displayed++
        }
    }

    for e := pcListPic.Front(); e != nil; e = e.Next() {
        pcPic := e.Value.(*TLibCommon.TComPic)
        //fmt.Printf("tId=%d, %v, %d, %d, %d, %d\n", tId, pcPic.GetOutputMark(), not_displayed, pcPic.GetNumReorderPics(tId), int(pcPic.GetPOC()), this.m_iPOCLastDisplay);
        if pcPic.GetOutputMark() && (not_displayed > maxNumReorderPics && int(pcPic.GetPOC()) > this.m_iPOCLastDisplay) {
            // write to file
            not_displayed--
            if this.m_pchReconFile != "" {
                conf := pcPic.GetConformanceWindow()
                var defDisp *TLibCommon.Window
                if this.m_respectDefDispWindow != 0 {
                    defDisp = pcPic.GetDefDisplayWindow()
                } else {
                    defDisp = TLibCommon.NewWindow()
                }
				
				/*fmt.Printf("(%d %d) (%d %d) (%d %d) (%d %d)\n", conf.GetWindowLeftOffset(), defDisp.GetWindowLeftOffset(),
               conf.GetWindowRightOffset(), defDisp.GetWindowRightOffset(),
               conf.GetWindowTopOffset(), defDisp.GetWindowTopOffset(),
               conf.GetWindowBottomOffset(), defDisp.GetWindowBottomOffset());*/
                //fmt.Printf(" [xWriteOutput POC %4d] ", pcPic.GetPOC());
                
                this.m_cTVideoIOYuvReconFile.Write(pcPic.GetPicYuvRec(),
                    conf.GetWindowLeftOffset()+defDisp.GetWindowLeftOffset(),
                    conf.GetWindowRightOffset()+defDisp.GetWindowRightOffset(),
                    conf.GetWindowTopOffset()+defDisp.GetWindowTopOffset(),
                    conf.GetWindowBottomOffset()+defDisp.GetWindowBottomOffset())
            }

            // update POC of display order
            this.m_iPOCLastDisplay = int(pcPic.GetPOC())
            //fmt.Printf("m_iPOCLastDisplay=%d\n",this.m_iPOCLastDisplay);

            // erase non-referenced picture in the reference picture list after display
            if !pcPic.GetSlice(0).IsReferenced() && pcPic.GetReconMark() == true {
                //#   if !D   YN_REF_FREE
                pcPic.SetReconMark(false)

                // mark it should be extended later
                pcPic.GetPicYuvRec().SetBorderExtension(false)

                //#   else
                //              pcPic->destroy();
                //              pcListPic->erase( iterPic );
                //              iterPic = pcListPic->begin(); // to the beginning, non-efficient way, have to be revised!
                //              continue;
                //#   endif
            }
            pcPic.SetOutputMark(false)
        }
    }

}

func (this *TAppDecTop) xFlushOutput(pcListPic *list.List) {
    if pcListPic == nil {
        return
    }

    //fmt.Printf("list len=%d\n", pcListPic.Len());

    for e := pcListPic.Front(); e != nil; e = e.Next() {
        pcPic := e.Value.(*TLibCommon.TComPic)
        if pcPic.GetOutputMark() {
            // write to file
            if this.m_pchReconFile != "" {
                conf := pcPic.GetConformanceWindow()
                var defDisp *TLibCommon.Window
                if this.m_respectDefDispWindow != 0 {
                    defDisp = pcPic.GetDefDisplayWindow()
                } else {
                    defDisp = TLibCommon.NewWindow()
                }
                
                //fmt.Printf(" [xFlushOutput POC %4d] ", pcPic.GetPOC());

                this.m_cTVideoIOYuvReconFile.Write(pcPic.GetPicYuvRec(),
                    conf.GetWindowLeftOffset()+defDisp.GetWindowLeftOffset(),
                    conf.GetWindowRightOffset()+defDisp.GetWindowRightOffset(),
                    conf.GetWindowTopOffset()+defDisp.GetWindowTopOffset(),
                    conf.GetWindowBottomOffset()+defDisp.GetWindowBottomOffset())
            }

            // update POC of display order
            this.m_iPOCLastDisplay = int(pcPic.GetPOC())
            //fmt.Printf("m_iPOCLastDisplay=%d\n",this.m_iPOCLastDisplay);

            // erase non-referenced picture in the reference picture list after display
            if !pcPic.GetSlice(0).IsReferenced() && pcPic.GetReconMark() == true {
                //#if !DYN_REF_FREE
                pcPic.SetReconMark(false)

                // mark it should be extended later
                pcPic.GetPicYuvRec().SetBorderExtension(false)

                //#else
                //        pcPic->destroy();
                //        pcListPic->erase( iterPic );
                //        iterPic = pcListPic->begin(); // to the beginning, non-efficient way, have to be revised!
                //        continue;
                //#endif
            }
            pcPic.SetOutputMark(false)
        }
        //#if !DYN_REF_FREE
        if pcPic != nil {
            pcPic.Destroy()
            //delete pcPic;
            pcPic = nil
        }
        //#endif
    }

    pcListPic.Init()
    this.m_iPOCLastDisplay = -TLibCommon.MAX_INT
}

func (this *TAppDecTop) IsNaluWithinTargetDecLayerIdSet(nalu *TLibDecoder.InputNALUnit) bool {
    if this.m_targetDecLayerIdSet.Len() == 0 { // By default, the set is empty, meaning all LayerIds are allowed
        return true
    }
    for e := this.m_targetDecLayerIdSet.Front(); e != nil; e = e.Next() {
        it := e.Value.(int)
        if int(nalu.GetReservedZero6Bits()) == it {
            return true
        }
    }
    return false
}
