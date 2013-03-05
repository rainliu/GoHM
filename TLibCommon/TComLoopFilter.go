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

package TLibCommon

import (
//"fmt"
)

const DEBLOCK_SMALLEST_BLOCK = 8

// ====================================================================================================================
// Constants
// ====================================================================================================================

const EDGE_VER = 0
const EDGE_HOR = 1
const DEFAULT_INTRA_TC_OFFSET = 2 ///< Default intra TC offset

// ====================================================================================================================
// Tables
// ====================================================================================================================

var tctable_8x8 = [54]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 3, 3, 3, 3, 4, 4, 4, 5, 5, 6, 6, 7, 8, 9, 10, 11, 13, 14, 16, 18, 20, 22, 24}

var betatable_8x8 = [52]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 20, 22, 24, 26, 28, 30, 32, 34, 36, 38, 40, 42, 44, 46, 48, 50, 52, 54, 56, 58, 60, 62, 64}

// ====================================================================================================================
// Class definition
// ====================================================================================================================

/// deblocking filter class
type TComLoopFilter struct {
    //private:
    m_uiNumPartitions uint
    m_aapucBS         [2][]byte ///< Bs for [Ver/Hor][Y/U/V][Blk_Idx]
    m_aapbEdgeFilter  [2][]bool
    m_stLFCUParam     LFCUParam ///< status structure

    m_bLFCrossTileBoundary bool
}

func NewTComLoopFilter() *TComLoopFilter {
    this := &TComLoopFilter{m_uiNumPartitions: 0, m_bLFCrossTileBoundary: true}

    for uiDir := 0; uiDir < 2; uiDir++ {
        this.m_aapucBS[uiDir] = nil
        this.m_aapbEdgeFilter[uiDir] = nil
    }

    return this
}

func (this *TComLoopFilter) Create(uiMaxCUDepth uint) {
    this.Destroy()
    this.m_uiNumPartitions = 1 << (uiMaxCUDepth << 1)
    for uiDir := 0; uiDir < 2; uiDir++ {
        this.m_aapucBS[uiDir] = make([]byte, this.m_uiNumPartitions)
        this.m_aapbEdgeFilter[uiDir] = make([]bool, this.m_uiNumPartitions)
    }
}
func (this *TComLoopFilter) Destroy() {
    /*for uiDir := 0; uiDir < 2; uiDir++ {
      if this.m_aapucBS!=nil {
        //delete [] this.m_aapucBS       [uiDir];
        this.m_aapucBS [uiDir] = nil;
      }
      for uiPlane := 0; uiPlane < 3; uiPlane++ {
        if this.m_aapbEdgeFilter[uiDir][uiPlane]!=nil {
          //delete [] this.m_aapbEdgeFilter[uiDir][uiPlane];
          this.m_aapbEdgeFilter[uiDir][uiPlane] = nil;
        }
      }
    }*/
}

func (this *TComLoopFilter) QpUV(iQpY int) int {
    if iQpY < 0 {
        return iQpY
    } else if iQpY > 57 {
        return iQpY - 6
    }

    return int(G_aucChromaScale[iQpY])
}

/// set configuration
func (this *TComLoopFilter) SetCfg(bLFCrossTileBoundary bool) {
    this.m_bLFCrossTileBoundary = bLFCrossTileBoundary
}

/// picture-level deblocking filter
func (this *TComLoopFilter) LoopFilterPic(pcPic *TComPic) {
    // Horizontal filtering
    for uiCUAddr := uint(0); uiCUAddr < pcPic.GetNumCUsInFrame(); uiCUAddr++ {
        pcCU := pcPic.GetCU(uiCUAddr)

        for i := uint(0); i < this.m_uiNumPartitions; i++ {
            this.m_aapucBS[EDGE_VER][i] = 0            //, sizeof( UChar ) *  );
            this.m_aapbEdgeFilter[EDGE_VER][i] = false //, sizeof( Bool  ) * this.m_uiNumPartitions );
        }

        // CU-based deblocking
        this.xDeblockCU(pcCU, 0, 0, EDGE_VER)
    }

    // Vertical filtering
    for uiCUAddr := uint(0); uiCUAddr < pcPic.GetNumCUsInFrame(); uiCUAddr++ {
        pcCU := pcPic.GetCU(uiCUAddr)

        for i := uint(0); i < this.m_uiNumPartitions; i++ {
            this.m_aapucBS[EDGE_HOR][i] = 0            //, sizeof( UChar ) * this.m_uiNumPartitions );
            this.m_aapbEdgeFilter[EDGE_HOR][i] = false //, sizeof( Bool  ) * this.m_uiNumPartitions );
        }

        // CU-based deblocking
        this.xDeblockCU(pcCU, 0, 0, EDGE_HOR)
    }
}

/// CU-level deblocking function
func (this *TComLoopFilter) xDeblockCU(pcCU *TComDataCU, uiAbsZorderIdx, uiDepth uint, Edge int) {
    if pcCU.GetPic() == nil || pcCU.GetPartitionSize1(uiAbsZorderIdx) == SIZE_NONE {
        return
    }
    pcPic := pcCU.GetPic()
    uiCurNumParts := pcPic.GetNumPartInCU() >> (uiDepth << 1)
    uiQNumParts := uiCurNumParts >> 2

    if uint(pcCU.GetDepth1(uiAbsZorderIdx)) > uiDepth {
        for uiPartIdx := 0; uiPartIdx < 4; uiPartIdx++ {
            uiLPelX := pcCU.GetCUPelX() + G_auiRasterToPelX[G_auiZscanToRaster[uiAbsZorderIdx]]
            uiTPelY := pcCU.GetCUPelY() + G_auiRasterToPelY[G_auiZscanToRaster[uiAbsZorderIdx]]
            if (uiLPelX < pcCU.GetSlice().GetSPS().GetPicWidthInLumaSamples()) && (uiTPelY < pcCU.GetSlice().GetSPS().GetPicHeightInLumaSamples()) {
                this.xDeblockCU(pcCU, uiAbsZorderIdx, uiDepth+1, Edge)
            }
            uiAbsZorderIdx += uiQNumParts
        }
        return
    }

    this.xSetLoopfilterParam(pcCU, uiAbsZorderIdx)

    this.xSetEdgefilterTU(pcCU, uiAbsZorderIdx, uiAbsZorderIdx, uiDepth)
    this.xSetEdgefilterPU(pcCU, uiAbsZorderIdx)

    iDir := Edge
    for uiPartIdx := uiAbsZorderIdx; uiPartIdx < uiAbsZorderIdx+uiCurNumParts; uiPartIdx++ {
        var uiBSCheck bool
        if (pcCU.GetSlice().GetSPS().GetMaxCUWidth() >> pcCU.GetSlice().GetSPS().GetMaxCUDepth()) == 4 {
            uiBSCheck = (iDir == EDGE_VER && uiPartIdx%2 == 0) || (iDir == EDGE_HOR && (uiPartIdx-((uiPartIdx>>2)<<2))/2 == 0)
        } else {
            uiBSCheck = true
        }

        if this.m_aapbEdgeFilter[iDir][uiPartIdx] && uiBSCheck {
            this.xGetBoundaryStrengthSingle(pcCU, iDir, uiPartIdx)
        }
    }

    uiPelsInPart := pcCU.GetSlice().GetSPS().GetMaxCUWidth() >> pcCU.GetSlice().GetSPS().GetMaxCUDepth()
    var PartIdxIncr uint

    if DEBLOCK_SMALLEST_BLOCK/uiPelsInPart != 0 {
        PartIdxIncr = DEBLOCK_SMALLEST_BLOCK / uiPelsInPart
    } else {
        PartIdxIncr = 1
    }

    //fmt.Printf("PartIdxIncr=%d ", PartIdxIncr);

    uiSizeInPU := pcPic.GetNumPartInWidth() >> (uiDepth)

    for iEdge := uint(0); iEdge < uiSizeInPU; iEdge += PartIdxIncr {
        this.xEdgeFilterLuma(pcCU, uiAbsZorderIdx, uiDepth, iDir, int(iEdge))
        if (uiPelsInPart > DEBLOCK_SMALLEST_BLOCK) || (iEdge%((DEBLOCK_SMALLEST_BLOCK<<1)/uiPelsInPart)) == 0 {
            this.xEdgeFilterChroma(pcCU, uiAbsZorderIdx, uiDepth, iDir, int(iEdge))
        }
    }
}

// set / get functions
func (this *TComLoopFilter) xSetLoopfilterParam(pcCU *TComDataCU, uiAbsZorderIdx uint) {
    uiX := pcCU.GetCUPelX() + G_auiRasterToPelX[G_auiZscanToRaster[uiAbsZorderIdx]]
    uiY := pcCU.GetCUPelY() + G_auiRasterToPelY[G_auiZscanToRaster[uiAbsZorderIdx]]

    var pcTempCU *TComDataCU
    var uiTempPartIdx uint

    this.m_stLFCUParam.bInternalEdge = !pcCU.GetSlice().GetDeblockingFilterDisable()

    if (uiX == 0) || pcCU.GetSlice().GetDeblockingFilterDisable() {
        this.m_stLFCUParam.bLeftEdge = false
    } else {
        this.m_stLFCUParam.bLeftEdge = true
    }
    if this.m_stLFCUParam.bLeftEdge {
        pcTempCU = pcCU.GetPULeft(&uiTempPartIdx, uiAbsZorderIdx, !pcCU.GetSlice().GetLFCrossSliceBoundaryFlag(), !this.m_bLFCrossTileBoundary)
        if pcTempCU != nil {
            this.m_stLFCUParam.bLeftEdge = true
        } else {
            this.m_stLFCUParam.bLeftEdge = false
        }
    }

    if (uiY == 0) || pcCU.GetSlice().GetDeblockingFilterDisable() {
        this.m_stLFCUParam.bTopEdge = false
    } else {
        this.m_stLFCUParam.bTopEdge = true
    }
    if this.m_stLFCUParam.bTopEdge {
        pcTempCU = pcCU.GetPUAbove(&uiTempPartIdx, uiAbsZorderIdx, !pcCU.GetSlice().GetLFCrossSliceBoundaryFlag(), false, !this.m_bLFCrossTileBoundary)

        if pcTempCU != nil {
            this.m_stLFCUParam.bTopEdge = true
        } else {
            this.m_stLFCUParam.bTopEdge = false
        }
    }
}

// filtering functions
func (this *TComLoopFilter) xSetEdgefilterTU(pcCU *TComDataCU, absTUPartIdx, uiAbsZorderIdx, uiDepth uint) {
    if uint(pcCU.GetTransformIdx1(uiAbsZorderIdx)+pcCU.GetDepth1(uiAbsZorderIdx)) > uiDepth {
        uiCurNumParts := pcCU.GetPic().GetNumPartInCU() >> (uiDepth << 1)
        uiQNumParts := uiCurNumParts >> 2
        for uiPartIdx := uint(0); uiPartIdx < 4; uiPartIdx++ {
            nsAddr := uiAbsZorderIdx
            this.xSetEdgefilterTU(pcCU, nsAddr, uiAbsZorderIdx, uiDepth+1)
            uiAbsZorderIdx += uiQNumParts
        }
        return
    }

    trWidth := uint(pcCU.GetWidth1(uiAbsZorderIdx) >> pcCU.GetTransformIdx1(uiAbsZorderIdx))
    trHeight := uint(pcCU.GetHeight1(uiAbsZorderIdx) >> pcCU.GetTransformIdx1(uiAbsZorderIdx))

    uiWidthInBaseUnits := trWidth / (pcCU.GetSlice().GetSPS().GetMaxCUWidth() >> pcCU.GetSlice().GetSPS().GetMaxCUDepth())
    uiHeightInBaseUnits := trHeight / (pcCU.GetSlice().GetSPS().GetMaxCUWidth() >> pcCU.GetSlice().GetSPS().GetMaxCUDepth())

    this.xSetEdgefilterMultiple(pcCU, absTUPartIdx, uiDepth, EDGE_VER, 0, this.m_stLFCUParam.bInternalEdge, uiWidthInBaseUnits, uiHeightInBaseUnits)
    this.xSetEdgefilterMultiple(pcCU, absTUPartIdx, uiDepth, EDGE_HOR, 0, this.m_stLFCUParam.bInternalEdge, uiWidthInBaseUnits, uiHeightInBaseUnits)
}
func (this *TComLoopFilter) xSetEdgefilterPU(pcCU *TComDataCU, uiAbsZorderIdx uint) {
    uiDepth := uint(pcCU.GetDepth1(uiAbsZorderIdx))
    uiWidthInBaseUnits := pcCU.GetPic().GetNumPartInWidth() >> uiDepth
    uiHeightInBaseUnits := pcCU.GetPic().GetNumPartInHeight() >> uiDepth
    uiHWidthInBaseUnits := uiWidthInBaseUnits >> 1
    uiHHeightInBaseUnits := uiHeightInBaseUnits >> 1
    uiQWidthInBaseUnits := uiWidthInBaseUnits >> 2
    uiQHeightInBaseUnits := uiHeightInBaseUnits >> 2

    this.xSetEdgefilterMultiple(pcCU, uiAbsZorderIdx, uiDepth, EDGE_VER, 0, this.m_stLFCUParam.bLeftEdge, 0, 0)
    this.xSetEdgefilterMultiple(pcCU, uiAbsZorderIdx, uiDepth, EDGE_HOR, 0, this.m_stLFCUParam.bTopEdge, 0, 0)

    switch pcCU.GetPartitionSize1(uiAbsZorderIdx) {
    case SIZE_2Nx2N:
        //break;
    case SIZE_2NxN:
        this.xSetEdgefilterMultiple(pcCU, uiAbsZorderIdx, uiDepth, EDGE_HOR, int(uiHHeightInBaseUnits), this.m_stLFCUParam.bInternalEdge, 0, 0)
    case SIZE_Nx2N:
        this.xSetEdgefilterMultiple(pcCU, uiAbsZorderIdx, uiDepth, EDGE_VER, int(uiHWidthInBaseUnits), this.m_stLFCUParam.bInternalEdge, 0, 0)
    case SIZE_NxN:
        this.xSetEdgefilterMultiple(pcCU, uiAbsZorderIdx, uiDepth, EDGE_VER, int(uiHWidthInBaseUnits), this.m_stLFCUParam.bInternalEdge, 0, 0)
        this.xSetEdgefilterMultiple(pcCU, uiAbsZorderIdx, uiDepth, EDGE_HOR, int(uiHHeightInBaseUnits), this.m_stLFCUParam.bInternalEdge, 0, 0)
    case SIZE_2NxnU:
        this.xSetEdgefilterMultiple(pcCU, uiAbsZorderIdx, uiDepth, EDGE_HOR, int(uiQHeightInBaseUnits), this.m_stLFCUParam.bInternalEdge, 0, 0)
    case SIZE_2NxnD:
        this.xSetEdgefilterMultiple(pcCU, uiAbsZorderIdx, uiDepth, EDGE_HOR, int(uiHeightInBaseUnits-uiQHeightInBaseUnits), this.m_stLFCUParam.bInternalEdge, 0, 0)
    case SIZE_nLx2N:
        this.xSetEdgefilterMultiple(pcCU, uiAbsZorderIdx, uiDepth, EDGE_VER, int(uiQWidthInBaseUnits), this.m_stLFCUParam.bInternalEdge, 0, 0)
    case SIZE_nRx2N:
        this.xSetEdgefilterMultiple(pcCU, uiAbsZorderIdx, uiDepth, EDGE_VER, int(uiWidthInBaseUnits-uiQWidthInBaseUnits), this.m_stLFCUParam.bInternalEdge, 0, 0)
    default:
        //break;
    }
}
func (this *TComLoopFilter) xGetBoundaryStrengthSingle(pcCU *TComDataCU, iDir int, uiAbsPartIdx uint) {
    pcSlice := pcCU.GetSlice()

    uiPartQ := uiAbsPartIdx
    pcCUQ := pcCU

    var uiPartP uint
    var pcCUP *TComDataCU
    uiBs := uint(0)

    //-- Calculate Block Index
    if iDir == EDGE_VER {
        pcCUP = pcCUQ.GetPULeft(&uiPartP, uiPartQ, !pcCU.GetSlice().GetLFCrossSliceBoundaryFlag(), !this.m_bLFCrossTileBoundary)
        //fmt.Printf("V:%d\n",uiPartP);
    } else { // (iDir == EDGE_HOR)
        pcCUP = pcCUQ.GetPUAbove(&uiPartP, uiPartQ, !pcCU.GetSlice().GetLFCrossSliceBoundaryFlag(), false, !this.m_bLFCrossTileBoundary)
    }

    //-- Set BS for Intra MB : BS = 4 or 3
    if pcCUP.IsIntra(uiPartP) || pcCUQ.IsIntra(uiPartQ) {
        uiBs = 2
    }

    //-- Set BS for not Intra MB : BS = 2 or 1 or 0
    if !pcCUP.IsIntra(uiPartP) && !pcCUQ.IsIntra(uiPartQ) {
        nsPartQ := uiPartQ
        nsPartP := uiPartP
        //fmt.Printf("1:m_aapucBS[%d][%d]=%d\n",iDir,uiAbsPartIdx,this.m_aapucBS[iDir][uiAbsPartIdx]);
        if this.m_aapucBS[iDir][uiAbsPartIdx] != 0 && (pcCUQ.GetCbf3(nsPartQ, TEXT_LUMA, uint(pcCUQ.GetTransformIdx1(nsPartQ))) != 0 || pcCUP.GetCbf3(nsPartP, TEXT_LUMA, uint(pcCUP.GetTransformIdx1(nsPartP))) != 0) {
            uiBs = 1
            //fmt.Printf("1.1:m_aapucBS[%d][%d]=%d\n",iDir,uiAbsPartIdx,this.m_aapucBS[iDir][uiAbsPartIdx]);
        } else {
            if iDir == EDGE_HOR {
                pcCUP = pcCUQ.GetPUAbove(&uiPartP, uiPartQ, !pcCU.GetSlice().GetLFCrossSliceBoundaryFlag(), false, !this.m_bLFCrossTileBoundary)
            }
            if pcSlice.IsInterB() || pcCUP.GetSlice().IsInterB() {
                //fmt.Printf("1.2:m_aapucBS[%d][%d]=%d\n",iDir,uiAbsPartIdx,this.m_aapucBS[iDir][uiAbsPartIdx]);
                var iRefIdx int
                var piRefP0, piRefP1, piRefQ0, piRefQ1 *TComPic
                iRefIdx = int(pcCUP.GetCUMvField(REF_PIC_LIST_0).GetRefIdx(int(uiPartP)))
                if iRefIdx < 0 {
                    piRefP0 = nil
                } else {
                    piRefP0 = pcCUP.GetSlice().GetRefPic(REF_PIC_LIST_0, iRefIdx)
                }
                iRefIdx = int(pcCUP.GetCUMvField(REF_PIC_LIST_1).GetRefIdx(int(uiPartP)))
                if iRefIdx < 0 {
                    piRefP1 = nil
                } else {
                    piRefP1 = pcCUP.GetSlice().GetRefPic(REF_PIC_LIST_1, iRefIdx)
                }
                iRefIdx = int(pcCUQ.GetCUMvField(REF_PIC_LIST_0).GetRefIdx(int(uiPartQ)))
                if iRefIdx < 0 {
                    piRefQ0 = nil
                } else {
                    piRefQ0 = pcSlice.GetRefPic(REF_PIC_LIST_0, iRefIdx)
                }
                iRefIdx = int(pcCUQ.GetCUMvField(REF_PIC_LIST_1).GetRefIdx(int(uiPartQ)))
                if iRefIdx < 0 {
                    piRefQ1 = nil
                } else {
                    piRefQ1 = pcSlice.GetRefPic(REF_PIC_LIST_1, iRefIdx)
                }

                pcMvP0 := pcCUP.GetCUMvField(REF_PIC_LIST_0).GetMv(int(uiPartP))
                pcMvP1 := pcCUP.GetCUMvField(REF_PIC_LIST_1).GetMv(int(uiPartP))
                pcMvQ0 := pcCUQ.GetCUMvField(REF_PIC_LIST_0).GetMv(int(uiPartQ))
                pcMvQ1 := pcCUQ.GetCUMvField(REF_PIC_LIST_1).GetMv(int(uiPartQ))

                if piRefP0 == nil {
                    pcMvP0.SetZero()
                }
                if piRefP1 == nil {
                    pcMvP1.SetZero()
                }
                if piRefQ0 == nil {
                    pcMvQ0.SetZero()
                }
                if piRefQ1 == nil {
                    pcMvQ1.SetZero()
                }

                if ((piRefP0 == piRefQ0) && (piRefP1 == piRefQ1)) || ((piRefP0 == piRefQ1) && (piRefP1 == piRefQ0)) {
                    //fmt.Printf("1.3:m_aapucBS[%d][%d]=%d\n",iDir,uiAbsPartIdx,this.m_aapucBS[iDir][uiAbsPartIdx]);
                    uiBs = 0
                    if piRefP0 != piRefP1 { // Different L0 & L1
                        if piRefP0 == piRefQ0 {
                            if (ABS(pcMvQ0.GetHor()-pcMvP0.GetHor()).(int16) >= 4) ||
                                (ABS(pcMvQ0.GetVer()-pcMvP0.GetVer()).(int16) >= 4) ||
                                (ABS(pcMvQ1.GetHor()-pcMvP1.GetHor()).(int16) >= 4) ||
                                (ABS(pcMvQ1.GetVer()-pcMvP1.GetVer()).(int16) >= 4) {
                                uiBs = 1
                            } else {
                                uiBs = 0
                            }
                        } else {
                            if (ABS(pcMvQ1.GetHor()-pcMvP0.GetHor()).(int16) >= 4) ||
                                (ABS(pcMvQ1.GetVer()-pcMvP0.GetVer()).(int16) >= 4) ||
                                (ABS(pcMvQ0.GetHor()-pcMvP1.GetHor()).(int16) >= 4) ||
                                (ABS(pcMvQ0.GetVer()-pcMvP1.GetVer()).(int16) >= 4) {
                                uiBs = 1
                            } else {
                                uiBs = 0
                            }
                        }
                    } else { // Same L0 & L1
                        if ((ABS(pcMvQ0.GetHor()-pcMvP0.GetHor()).(int16) >= 4) ||
                            (ABS(pcMvQ0.GetVer()-pcMvP0.GetVer()).(int16) >= 4) ||
                            (ABS(pcMvQ1.GetHor()-pcMvP1.GetHor()).(int16) >= 4) ||
                            (ABS(pcMvQ1.GetVer()-pcMvP1.GetVer()).(int16) >= 4)) &&
                            ((ABS(pcMvQ1.GetHor()-pcMvP0.GetHor()).(int16) >= 4) ||
                                (ABS(pcMvQ1.GetVer()-pcMvP0.GetVer()).(int16) >= 4) ||
                                (ABS(pcMvQ0.GetHor()-pcMvP1.GetHor()).(int16) >= 4) ||
                                (ABS(pcMvQ0.GetVer()-pcMvP1.GetVer()).(int16) >= 4)) {
                            uiBs = 1
                        } else {
                            uiBs = 0
                        }
                    }
                } else { // for all different Ref_Idx
                    uiBs = 1
                }
            } else { // pcSlice->isInterP()
                var iRefIdx int
                var piRefP0, piRefQ0 *TComPic
                iRefIdx = int(pcCUP.GetCUMvField(REF_PIC_LIST_0).GetRefIdx(int(uiPartP)))
                if iRefIdx < 0 {
                    piRefP0 = nil
                } else {
                    piRefP0 = pcCUP.GetSlice().GetRefPic(REF_PIC_LIST_0, iRefIdx)
                }
                iRefIdx = int(pcCUQ.GetCUMvField(REF_PIC_LIST_0).GetRefIdx(int(uiPartQ)))
                if iRefIdx < 0 {
                    piRefQ0 = nil
                } else {
                    piRefQ0 = pcSlice.GetRefPic(REF_PIC_LIST_0, iRefIdx)
                }
                pcMvP0 := pcCUP.GetCUMvField(REF_PIC_LIST_0).GetMv(int(uiPartP))
                pcMvQ0 := pcCUQ.GetCUMvField(REF_PIC_LIST_0).GetMv(int(uiPartQ))
                //fmt.Printf("p(%d,%d),q(%d,%d)\n", pcMvP0.GetAbsHor(), pcMvP0.GetAbsVer(), pcMvQ0.GetAbsHor(), pcMvQ0.GetAbsVer());

                if piRefP0 == nil {
                    pcMvP0.SetZero()
                }
                if piRefQ0 == nil {
                    pcMvQ0.SetZero()
                }

                if (piRefP0 != piRefQ0) ||
                    (ABS(pcMvQ0.GetHor()-pcMvP0.GetHor()).(int16) >= 4) ||
                    (ABS(pcMvQ0.GetVer()-pcMvP0.GetVer()).(int16) >= 4) {
                    uiBs = 1
                } else {
                    uiBs = 0
                }
                //fmt.Printf("(%d,%d):%d | %d>=4 | %d>=4\n",uiPartP,uiPartQ,B2U(piRefP0!=piRefQ0),pcMvP0.GetAbsHor(),pcMvP0.GetAbsVer());
            }
        }   // enf of "if( one of BCBP == 0 )"
    }   // enf of "if( not Intra )"

    this.m_aapucBS[iDir][uiAbsPartIdx] = byte(uiBs)
    //fmt.Printf("2:m_aapucBS[%d][%d]=%d\n",iDir,uiAbsPartIdx,this.m_aapucBS[iDir][uiAbsPartIdx]);
}

func (this *TComLoopFilter) xCalcBsIdx(pcCU *TComDataCU, uiAbsZorderIdx uint, iDir, iEdgeIdx, iBaseUnitIdx int) uint {
    pcPic := pcCU.GetPic()
    uiLCUWidthInBaseUnits := pcPic.GetNumPartInWidth()
    if iDir == 0 {
        return G_auiRasterToZscan[G_auiZscanToRaster[uiAbsZorderIdx]+uint(iBaseUnitIdx)*uiLCUWidthInBaseUnits+uint(iEdgeIdx)]
    }

    return G_auiRasterToZscan[G_auiZscanToRaster[uiAbsZorderIdx]+uint(iEdgeIdx)*uiLCUWidthInBaseUnits+uint(iBaseUnitIdx)]
}

func (this *TComLoopFilter) xSetEdgefilterMultiple(pcCU *TComDataCU, uiScanIdx, uiDepth uint, iDir, iEdgeIdx int, bValue bool, uiWidthInBaseUnits, uiHeightInBaseUnits uint) {
    if uiWidthInBaseUnits == 0 {
        uiWidthInBaseUnits = pcCU.GetPic().GetNumPartInWidth() >> uiDepth
    }
    if uiHeightInBaseUnits == 0 {
        uiHeightInBaseUnits = pcCU.GetPic().GetNumPartInHeight() >> uiDepth
    }
    var uiNumElem uint
    if iDir == 0 {
        uiNumElem = uiHeightInBaseUnits
    } else {
        uiNumElem = uiWidthInBaseUnits
    }
    //assert( uiNumElem > 0 );
    //assert( uiWidthInBaseUnits > 0 );
    //assert( uiHeightInBaseUnits > 0 );
    for ui := uint(0); ui < uiNumElem; ui++ {
        uiBsIdx := this.xCalcBsIdx(pcCU, uiScanIdx, iDir, iEdgeIdx, int(ui))
        this.m_aapbEdgeFilter[iDir][uiBsIdx] = bValue
        if iEdgeIdx == 0 {
            this.m_aapucBS[iDir][uiBsIdx] = B2U(bValue)
        }
    }
}

func (this *TComLoopFilter) xEdgeFilterLuma(pcCU *TComDataCU, uiAbsZorderIdx, uiDepth uint, iDir, iEdge int) {
    pcPicYuvRec := pcCU.GetPic().GetPicYuvRec()

    //piSrc    := pcPicYuvRec.GetLumaAddr2( int(pcCU.GetAddr()), int(uiAbsZorderIdx) );
    //piTmpSrc := piSrc;
    //iTmpSrcOffset := 0;

    piSrc := pcPicYuvRec.GetBufY()
    offsetY := pcPicYuvRec.m_cuOffsetY[pcCU.GetAddr()] + pcPicYuvRec.m_buOffsetY[G_auiZscanToRaster[uiAbsZorderIdx]]
    iTmpSrcOffset := pcPicYuvRec.m_iLumaMarginY*pcPicYuvRec.GetStride() + pcPicYuvRec.m_iLumaMarginX + offsetY

    //fmt.Printf("(%d,%d)--%x--\n", pcCU.GetAddr(), uiAbsZorderIdx, piSrc[iTmpSrcOffset]);

    iStride := pcPicYuvRec.GetStride()
    iQP := 0
    iQP_P := 0
    iQP_Q := 0
    uiNumParts := pcCU.GetPic().GetNumPartInWidth() >> uiDepth

    uiPelsInPart := pcCU.GetSlice().GetSPS().GetMaxCUWidth() >> pcCU.GetSlice().GetSPS().GetMaxCUDepth()
    uiBsAbsIdx := uint(0)
    uiBs := uint(0)
    var iOffset, iSrcStep int

    bPCMFilter := pcCU.GetSlice().GetSPS().GetUsePCM() && pcCU.GetSlice().GetSPS().GetPCMFilterDisableFlag()
    bPartPNoFilter := false
    bPartQNoFilter := false
    uiPartPIdx := uint(0)
    uiPartQIdx := uint(0)
    pcCUP := pcCU
    pcCUQ := pcCU
    betaOffsetDiv2 := pcCUQ.GetSlice().GetDeblockingFilterBetaOffsetDiv2()
    tcOffsetDiv2 := pcCUQ.GetSlice().GetDeblockingFilterTcOffsetDiv2()

    if iDir == EDGE_VER {
        iOffset = 1
        iSrcStep = iStride
        iTmpSrcOffset += iEdge * int(uiPelsInPart) //piTmpSrc = piTmpSrc[ iEdge*int(uiPelsInPart):];
    } else { // (iDir == EDGE_HOR)
        iOffset = iStride
        iSrcStep = 1
        iTmpSrcOffset += iEdge * int(uiPelsInPart) * iStride //piTmpSrc = piTmpSrc[ iEdge*int(uiPelsInPart)*iStride:];
    }

    for iIdx := uint(0); iIdx < uiNumParts; iIdx++ {
        uiBsAbsIdx = this.xCalcBsIdx(pcCU, uiAbsZorderIdx, iDir, iEdge, int(iIdx))
        uiBs = uint(this.m_aapucBS[iDir][uiBsAbsIdx])
        //fmt.Printf("iIdx=%d,uiNumParts=%d,iDir=%d,uiBsAbsIdx=%d,uiBs=%d\n", iIdx, uiNumParts,iDir,uiBsAbsIdx,uiBs);
        if uiBs != 0 {
            iQP_Q = int(pcCU.GetQP1(uiBsAbsIdx))
            uiPartQIdx = uiBsAbsIdx
            // Derive neighboring PU index
            if iDir == EDGE_VER {
                pcCUP = pcCUQ.GetPULeft(&uiPartPIdx, uiPartQIdx, !pcCU.GetSlice().GetLFCrossSliceBoundaryFlag(), !this.m_bLFCrossTileBoundary)
            } else { // (iDir == EDGE_HOR)
                pcCUP = pcCUQ.GetPUAbove(&uiPartPIdx, uiPartQIdx, !pcCU.GetSlice().GetLFCrossSliceBoundaryFlag(), false, !this.m_bLFCrossTileBoundary)
            }

            iQP_P = int(pcCUP.GetQP1(uiPartPIdx))
            iQP = (iQP_P + iQP_Q + 1) >> 1
            iBitdepthScale := 1 << uint(G_bitDepthY-8)

            iIndexTC := CLIP3(0, MAX_QP+DEFAULT_INTRA_TC_OFFSET, int(iQP+DEFAULT_INTRA_TC_OFFSET*int(uiBs-1)+(tcOffsetDiv2<<1))).(int)
            iIndexB := CLIP3(0, MAX_QP, iQP+(betaOffsetDiv2<<1)).(int)

            iTc := int(tctable_8x8[iIndexTC]) * iBitdepthScale
            iBeta := int(betatable_8x8[iIndexB]) * iBitdepthScale
            iSideThreshold := (iBeta + (iBeta >> 1)) >> 3
            iThrCut := iTc * 10

            var uiBlocksInPart uint
            if uiPelsInPart/4 != 0 {
                uiBlocksInPart = uiPelsInPart / 4
            } else {
                uiBlocksInPart = 1
            }
            //fmt.Printf("\nuiBlocksInPart=%d: ",uiBlocksInPart);
            for iBlkIdx := uint(0); iBlkIdx < uiBlocksInPart; iBlkIdx++ {
                dp0 := this.xCalcDP(piSrc[iTmpSrcOffset+iSrcStep*int(iIdx*uiPelsInPart+iBlkIdx*4+0)-iOffset*3:], iOffset)
                dq0 := this.xCalcDQ(piSrc[iTmpSrcOffset+iSrcStep*int(iIdx*uiPelsInPart+iBlkIdx*4+0):], iOffset)
                dp3 := this.xCalcDP(piSrc[iTmpSrcOffset+iSrcStep*int(iIdx*uiPelsInPart+iBlkIdx*4+3)-iOffset*3:], iOffset)
                dq3 := this.xCalcDQ(piSrc[iTmpSrcOffset+iSrcStep*int(iIdx*uiPelsInPart+iBlkIdx*4+3):], iOffset)
                d0 := dp0 + dq0
                d3 := dp3 + dq3

                dp := dp0 + dp3
                dq := dq0 + dq3
                d := d0 + d3

                if bPCMFilter || pcCU.GetSlice().GetPPS().GetTransquantBypassEnableFlag() {
                    // Check if each of PUs is I_PCM with LF disabling
                    bPartPNoFilter = (bPCMFilter && pcCUP.GetIPCMFlag1(uiPartPIdx))
                    bPartQNoFilter = (bPCMFilter && pcCUQ.GetIPCMFlag1(uiPartQIdx))

                    // check if each of PUs is lossless coded
                    bPartPNoFilter = bPartPNoFilter || (pcCUP.IsLosslessCoded(uiPartPIdx))
                    bPartQNoFilter = bPartQNoFilter || (pcCUQ.IsLosslessCoded(uiPartQIdx))
                }
                //fmt.Printf("%d<%d: ",d, iBeta);
                if int(d) < iBeta {
                    bFilterP := (int(dp) < iSideThreshold)
                    bFilterQ := (int(dq) < iSideThreshold)

                    sw := this.xUseStrongFiltering(iOffset, 2*int(d0), iBeta, iTc, piSrc[iTmpSrcOffset+iSrcStep*int(iIdx*uiPelsInPart+iBlkIdx*4+0)-iOffset*4:]) &&
                        this.xUseStrongFiltering(iOffset, 2*int(d3), iBeta, iTc, piSrc[iTmpSrcOffset+iSrcStep*int(iIdx*uiPelsInPart+iBlkIdx*4+3)-iOffset*4:])

                    for i := uint(0); i < DEBLOCK_SMALLEST_BLOCK/2; i++ {
                        this.xPelFilterLuma(piSrc[iTmpSrcOffset+iSrcStep*int(iIdx*uiPelsInPart+iBlkIdx*4+i)-iOffset*4:], iOffset, iTc, sw, bPartPNoFilter, bPartQNoFilter, iThrCut, bFilterP, bFilterQ)
                    }
                }
            }
        }
    }
}
func (this *TComLoopFilter) xEdgeFilterChroma(pcCU *TComDataCU, uiAbsZorderIdx, uiDepth uint, iDir, iEdge int) {
    pcPicYuvRec := pcCU.GetPic().GetPicYuvRec()
    iStride := pcPicYuvRec.GetCStride()
    //piSrcCb     := pcPicYuvRec.GetCbAddr2( int(pcCU.GetAddr()), int(uiAbsZorderIdx) );
    //piSrcCr     := pcPicYuvRec.GetCrAddr2( int(pcCU.GetAddr()), int(uiAbsZorderIdx) );

    piSrcCb := pcPicYuvRec.GetBufU()
    piSrcCr := pcPicYuvRec.GetBufV()

    offsetChroma := pcPicYuvRec.m_cuOffsetC[pcCU.GetAddr()] + pcPicYuvRec.m_buOffsetC[G_auiZscanToRaster[uiAbsZorderIdx]]
    piTmpSrcOffsetChroma := pcPicYuvRec.m_iChromaMarginY*pcPicYuvRec.GetCStride() + pcPicYuvRec.m_iChromaMarginX + offsetChroma

    iQP := 0
    iQP_P := 0
    iQP_Q := 0

    uiPelsInPartChroma := pcCU.GetSlice().GetSPS().GetMaxCUWidth() >> (pcCU.GetSlice().GetSPS().GetMaxCUDepth() + 1)

    var iOffset, iSrcStep int

    uiLCUWidthInBaseUnits := pcCU.GetPic().GetNumPartInWidth()

    bPCMFilter := pcCU.GetSlice().GetSPS().GetUsePCM() && pcCU.GetSlice().GetSPS().GetPCMFilterDisableFlag()
    bPartPNoFilter := false
    bPartQNoFilter := false
    var uiPartPIdx, uiPartQIdx uint
    var pcCUP *TComDataCU
    pcCUQ := pcCU
    tcOffsetDiv2 := pcCU.GetSlice().GetDeblockingFilterTcOffsetDiv2()

    // Vertical Position
    uiEdgeNumInLCUVert := G_auiZscanToRaster[uiAbsZorderIdx]%uiLCUWidthInBaseUnits + uint(iEdge)
    uiEdgeNumInLCUHor := G_auiZscanToRaster[uiAbsZorderIdx]/uiLCUWidthInBaseUnits + uint(iEdge)

    if (uiPelsInPartChroma < DEBLOCK_SMALLEST_BLOCK) && (((uiEdgeNumInLCUVert%(DEBLOCK_SMALLEST_BLOCK/uiPelsInPartChroma) != 0) && (iDir == 0)) || ((uiEdgeNumInLCUHor%(DEBLOCK_SMALLEST_BLOCK/uiPelsInPartChroma) != 0) && iDir != 0)) {
        return
    }

    uiNumParts := pcCU.GetPic().GetNumPartInWidth() >> uiDepth

    var uiBsAbsIdx uint
    var ucBs byte

    //piTmpSrcCb := piSrcCb;
    //piTmpSrcCr := piSrcCr;

    if iDir == EDGE_VER {
        iOffset = 1
        iSrcStep = iStride
        piTmpSrcOffsetChroma += iEdge * int(uiPelsInPartChroma) //piTmpSrcCb =piTmpSrcCb[ uint(iEdge)*uiPelsInPartChroma:];
        //piTmpSrcCr =piTmpSrcCr[ uint(iEdge)*uiPelsInPartChroma:];
    } else { // (iDir == EDGE_HOR)
        iOffset = iStride
        iSrcStep = 1
        piTmpSrcOffsetChroma += iEdge * iStride * int(uiPelsInPartChroma) //piTmpSrcCb =piTmpSrcCb[ iEdge*iStride*int(uiPelsInPartChroma):];
        //piTmpSrcCr =piTmpSrcCr[ iEdge*iStride*int(uiPelsInPartChroma):];
    }

    for iIdx := uint(0); iIdx < uiNumParts; iIdx++ {
        ucBs = 0

        uiBsAbsIdx = this.xCalcBsIdx(pcCU, uiAbsZorderIdx, iDir, iEdge, int(iIdx))
        ucBs = this.m_aapucBS[iDir][uiBsAbsIdx]

        if ucBs > 1 {
            iQP_Q = int(pcCU.GetQP1(uiBsAbsIdx))
            uiPartQIdx = uiBsAbsIdx
            // Derive neighboring PU index
            if iDir == EDGE_VER {
                pcCUP = pcCUQ.GetPULeft(&uiPartPIdx, uiPartQIdx, !pcCU.GetSlice().GetLFCrossSliceBoundaryFlag(), !this.m_bLFCrossTileBoundary)
            } else { // (iDir == EDGE_HOR)
                pcCUP = pcCUQ.GetPUAbove(&uiPartPIdx, uiPartQIdx, !pcCU.GetSlice().GetLFCrossSliceBoundaryFlag(), false, !this.m_bLFCrossTileBoundary)
            }

            iQP_P = int(pcCUP.GetQP1(uiPartPIdx))

            if bPCMFilter || pcCU.GetSlice().GetPPS().GetTransquantBypassEnableFlag() {
                // Check if each of PUs is I_PCM with LF disabling
                bPartPNoFilter = (bPCMFilter && pcCUP.GetIPCMFlag1(uiPartPIdx))
                bPartQNoFilter = (bPCMFilter && pcCUQ.GetIPCMFlag1(uiPartQIdx))

                // check if each of PUs is lossless coded
                bPartPNoFilter = bPartPNoFilter || (pcCUP.IsLosslessCoded(uiPartPIdx))
                bPartQNoFilter = bPartQNoFilter || (pcCUQ.IsLosslessCoded(uiPartQIdx))
            }

            for chromaIdx := 0; chromaIdx < 2; chromaIdx++ {
                var chromaQPOffset int
                var piTmpSrcChroma []Pel

                if chromaIdx == 0 {
                    chromaQPOffset = pcCU.GetSlice().GetPPS().GetChromaCbQpOffset()
                    piTmpSrcChroma = piSrcCb //piTmpSrcCb;
                } else {
                    chromaQPOffset = pcCU.GetSlice().GetPPS().GetChromaCrQpOffset()
                    piTmpSrcChroma = piSrcCr //piTmpSrcCr;
                }

                iQP = this.QpUV(((iQP_P + iQP_Q + 1) >> 1) + chromaQPOffset)
                iBitdepthScale := 1 << uint(G_bitDepthC-8)

                iIndexTC := CLIP3(0, MAX_QP+DEFAULT_INTRA_TC_OFFSET, iQP+DEFAULT_INTRA_TC_OFFSET*int(ucBs-1)+(tcOffsetDiv2<<1)).(int)
                iTc := int(tctable_8x8[iIndexTC]) * iBitdepthScale

                for uiStep := uint(0); uiStep < uiPelsInPartChroma; uiStep++ {
                    this.xPelFilterChroma(piTmpSrcChroma[piTmpSrcOffsetChroma+iSrcStep*int(uiStep+iIdx*uiPelsInPartChroma)-iOffset*2:], iOffset, iTc, bPartPNoFilter, bPartQNoFilter)
                }
            }
        }
    }
}

func (this *TComLoopFilter) xPelFilterLuma(piSrc2 []Pel, iOffset, tc int, sw, bPartPNoFilter, bPartQNoFilter bool, iThrCut int, bFilterSecondP, bFilterSecondQ bool) {
    var delta int

    m4 := piSrc2[iOffset*4+0]
    m3 := piSrc2[iOffset*4-iOffset]
    m5 := piSrc2[iOffset*4+iOffset]
    m2 := piSrc2[iOffset*4-iOffset*2]
    m6 := piSrc2[iOffset*4+iOffset*2]
    m1 := piSrc2[iOffset*4-iOffset*3]
    m7 := piSrc2[iOffset*4+iOffset*3]
    m0 := piSrc2[iOffset*4-iOffset*4]

    //fmt.Printf("%x %x %x %x %x %x %x %x\n", m0, m1, m2, m3, m4, m5, m6, m7);

    if sw {
        piSrc2[iOffset*4-iOffset] = CLIP3(m3-2*Pel(tc), m3+2*Pel(tc), ((m1 + 2*m2 + 2*m3 + 2*m4 + m5 + 4) >> 3)).(Pel)
        piSrc2[iOffset*4+0] = CLIP3(m4-2*Pel(tc), m4+2*Pel(tc), ((m2 + 2*m3 + 2*m4 + 2*m5 + m6 + 4) >> 3)).(Pel)
        piSrc2[iOffset*4-iOffset*2] = CLIP3(m2-2*Pel(tc), m2+2*Pel(tc), ((m1 + m2 + m3 + m4 + 2) >> 2)).(Pel)
        piSrc2[iOffset*4+iOffset] = CLIP3(m5-2*Pel(tc), m5+2*Pel(tc), ((m3 + m4 + m5 + m6 + 2) >> 2)).(Pel)
        piSrc2[iOffset*4-iOffset*3] = CLIP3(m1-2*Pel(tc), m1+2*Pel(tc), ((2*m0 + 3*m1 + m2 + m3 + m4 + 4) >> 3)).(Pel)
        piSrc2[iOffset*4+iOffset*2] = CLIP3(m6-2*Pel(tc), m6+2*Pel(tc), ((m3 + m4 + m5 + 3*m6 + 2*m7 + 4) >> 3)).(Pel)
    } else {
        /* Weak filter */
        delta = int(9*(m4-m3)-3*(m5-m2)+8) >> 4

        if ABS(delta).(int) < iThrCut {
            delta = CLIP3(-tc, tc, delta).(int)
            piSrc2[iOffset*4-iOffset] = ClipY((m3 + Pel(delta)))
            piSrc2[iOffset*4+0] = ClipY((m4 - Pel(delta)))

            tc2 := tc >> 1
            if bFilterSecondP {
                delta1 := CLIP3(-tc2, tc2, ((int(((m1+m3+1)>>1)-m2) + delta) >> 1)).(int)
                piSrc2[iOffset*4-iOffset*2] = ClipY((m2 + Pel(delta1)))
            }

            if bFilterSecondQ {
                delta2 := CLIP3(-tc2, tc2, ((int(((m6+m4+1)>>1)-m5) - delta) >> 1)).(int)
                piSrc2[iOffset*4+iOffset] = ClipY((m5 + Pel(delta2)))
            }
        }
    }

    if bPartPNoFilter {
        piSrc2[iOffset*4-iOffset] = m3
        piSrc2[iOffset*4-iOffset*2] = m2
        piSrc2[iOffset*4-iOffset*3] = m1
    }
    if bPartQNoFilter {
        piSrc2[iOffset*4+0] = m4
        piSrc2[iOffset*4+iOffset] = m5
        piSrc2[iOffset*4+iOffset*2] = m6
    }
}
func (this *TComLoopFilter) xPelFilterChroma(piSrc2 []Pel, iOffset, tc int, bPartPNoFilter, bPartQNoFilter bool) {
    var delta int

    m4 := piSrc2[iOffset*2+0]
    m3 := piSrc2[iOffset*2-iOffset]
    m5 := piSrc2[iOffset*2+iOffset]
    m2 := piSrc2[iOffset*2-iOffset*2]

    delta = CLIP3(-tc, tc, int((((m4-m3)<<2)+m2-m5+4)>>3)).(int)
    piSrc2[iOffset*2-iOffset] = ClipC(m3 + Pel(delta))
    piSrc2[iOffset*2+0] = ClipC(m4 - Pel(delta))

    if bPartPNoFilter {
        piSrc2[iOffset*2-iOffset] = m3
    }
    if bPartQNoFilter {
        piSrc2[iOffset*2+0] = m4
    }
}

func (this *TComLoopFilter) xUseStrongFiltering(offset, d, beta, tc int, piSrc2 []Pel) bool {
    m4 := piSrc2[offset*4+0]
    m3 := piSrc2[offset*4-offset]
    m7 := piSrc2[offset*4+offset*3]
    m0 := piSrc2[offset*4-offset*4]

    d_strong := int(ABS(m0-m3).(Pel) + ABS(m7-m4).(Pel))

    return (d_strong < (beta >> 3)) && (d < (beta >> 2)) && (int(ABS(m3-m4).(Pel)) < ((tc*5 + 1) >> 1))

}
func (this *TComLoopFilter) xCalcDP(piSrc2 []Pel, iOffset int) Pel {
    return ABS(piSrc2[iOffset*3-iOffset*3] - 2*piSrc2[iOffset*3-iOffset*2] + piSrc2[iOffset*3-iOffset]).(Pel)
}
func (this *TComLoopFilter) xCalcDQ(piSrc []Pel, iOffset int) Pel {
    return ABS(piSrc[0] - 2*piSrc[iOffset] + piSrc[iOffset*2]).(Pel)
}
