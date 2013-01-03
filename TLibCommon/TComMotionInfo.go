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

import ()

// ====================================================================================================================
// Type definition
// ====================================================================================================================

/// parameters for AMVP
type AMVPInfo struct {
    MvCand [AMVP_MAX_NUM_CANDS_MEM]TComMv ///< array of motion vector predictor candidates
    IN     int                            ///< number of motion vector predictor candidates
}

// ====================================================================================================================
// Class definition
// ====================================================================================================================

/// class for motion vector with reference index
type TComMvField struct {
    //private:
    m_acMv    TComMv
    m_iRefIdx int8
}

func NewTComMvField() *TComMvField {
    return &TComMvField{m_iRefIdx: NOT_VALID}
}

func (this *TComMvField) SetMvField(cMv TComMv, iRefIdx int8) {
    this.m_acMv.SetHor(cMv.GetHor())
    this.m_acMv.SetVer(cMv.GetVer())

    this.m_iRefIdx = iRefIdx
}

func (this *TComMvField) SetRefIdx(refIdx int8) {
    this.m_iRefIdx = refIdx
}

func (this *TComMvField) GetMv() TComMv {
    return this.m_acMv
}
func (this *TComMvField) GetRefIdx() int8 {
    return this.m_iRefIdx
}
func (this *TComMvField) GetHor() int16 {
    return this.m_acMv.GetHor()
}
func (this *TComMvField) GetVer() int16 {
    return this.m_acMv.GetVer()
}

/// class for motion information in one CU
type TComCUMvField struct {
    //private:
    m_pcMv           []TComMv
    m_pcMvd          []TComMv
    m_piRefIdx       []int8
    m_uiNumPartition uint
    m_cAMVPInfo      AMVPInfo
}

/*
  template <typename T>
  Void setAll( T *p, T const & val, PartSize eCUMode, Int iPartAddr, UInt uiDepth, Int iPartIdx );
*/
//public:
func NewTComCUMvField() *TComCUMvField {
    return &TComCUMvField{}
}

// ------------------------------------------------------------------------------------------------------------------
// create / destroy
// ------------------------------------------------------------------------------------------------------------------

func (this *TComCUMvField) Create(uiNumPartition uint) {
    this.m_pcMv = make([]TComMv, uiNumPartition)
    this.m_pcMvd = make([]TComMv, uiNumPartition)
    this.m_piRefIdx = make([]int8, uiNumPartition)

    this.m_uiNumPartition = uiNumPartition
}

func (this *TComCUMvField) Destroy() {
    this.m_pcMv = nil
    this.m_pcMvd = nil
    this.m_piRefIdx = nil

    this.m_uiNumPartition = 0
}

// ------------------------------------------------------------------------------------------------------------------
// clear / copy
// ------------------------------------------------------------------------------------------------------------------

func (this *TComCUMvField) ClearMvField() {
    for i := uint(0); i < this.m_uiNumPartition; i++ {
        this.m_pcMv[i].SetZero()
        this.m_pcMvd[i].SetZero()
    }
    //assert( sizeof( *m_piRefIdx ) == 1 );
    for i := uint(0); i < this.m_uiNumPartition; i++ {
        this.m_piRefIdx[i] = NOT_VALID
    }
    //memset( m_piRefIdx, NOT_VALID, this.m_uiNumPartition * sizeof( *this.m_piRefIdx ) );
}

func (this *TComCUMvField) CopyFrom(pcCUMvFieldSrc *TComCUMvField, uiNumPartSrc uint, iPartAddrDst int) {
    //Int iSizeInTComMv := sizeof( TComMv ) * iNumPartSrc;
    for i := 0; i < int(uiNumPartSrc); i++ {
        this.m_pcMv[i+iPartAddrDst] = pcCUMvFieldSrc.m_pcMv[i]
        this.m_pcMvd[i+iPartAddrDst] = pcCUMvFieldSrc.m_pcMvd[i]
        this.m_piRefIdx[i+iPartAddrDst] = pcCUMvFieldSrc.m_piRefIdx[i]
    }
    //memcpy( m_pcMv     + iPartAddrDst, pcCUMvFieldSrc->m_pcMv,     iSizeInTComMv );
    //memcpy( m_pcMvd    + iPartAddrDst, pcCUMvFieldSrc->m_pcMvd,    iSizeInTComMv );
    //memcpy( m_piRefIdx + iPartAddrDst, pcCUMvFieldSrc->m_piRefIdx, sizeof( *m_piRefIdx ) * iNumPartSrc );
}
func (this *TComCUMvField) CopyTo2(pcCUMvFieldDst *TComCUMvField, iPartAddrDst int) {
    this.CopyTo4(pcCUMvFieldDst, iPartAddrDst, 0, this.m_uiNumPartition)
}
func (this *TComCUMvField) CopyTo4(pcCUMvFieldDst *TComCUMvField, iPartAddrDst int, uiOffset, uiNumPart uint) {
    //Int iSizeInTComMv = sizeof( TComMv ) * uiNumPart;
    iOffset := uiOffset + uint(iPartAddrDst)
    for i := uint(0); i < uiNumPart; i++ {
        pcCUMvFieldDst.m_pcMv[i+iOffset] = this.m_pcMv[i+uiOffset]
        pcCUMvFieldDst.m_pcMvd[i+iOffset] = this.m_pcMvd[i+uiOffset]
        pcCUMvFieldDst.m_piRefIdx[i+iOffset] = this.m_piRefIdx[i+uiOffset]
    }
    //memcpy( pcCUMvFieldDst->m_pcMv     + iOffset, m_pcMv     + uiOffset, iSizeInTComMv );
    //memcpy( pcCUMvFieldDst->m_pcMvd    + iOffset, m_pcMvd    + uiOffset, iSizeInTComMv );
    //memcpy( pcCUMvFieldDst->m_piRefIdx + iOffset, m_piRefIdx + uiOffset, sizeof( *m_piRefIdx ) * uiNumPart );
}

// ------------------------------------------------------------------------------------------------------------------
// get
// ------------------------------------------------------------------------------------------------------------------
func (this *TComCUMvField) GetMvs(offset int) []TComMv {
    return this.m_pcMv[offset:]
}
func (this *TComCUMvField) GetMvds(offset int) []TComMv {
    return this.m_pcMvd[offset:]
}
func (this *TComCUMvField) GetRefIdxs(offset int) []int8 {
    return this.m_piRefIdx[offset:]
}

func (this *TComCUMvField) GetMv(iIdx int) TComMv {
    return this.m_pcMv[iIdx]
}
func (this *TComCUMvField) GetMvd(iIdx int) TComMv {
    return this.m_pcMvd[iIdx]
}
func (this *TComCUMvField) GetRefIdx(iIdx int) int8 {
    return this.m_piRefIdx[iIdx]
}

func (this *TComCUMvField) GetAMVPInfo() *AMVPInfo {
    return &this.m_cAMVPInfo
}

// --------------------------------------------------------------------------------------------------------------------
// Set
// --------------------------------------------------------------------------------------------------------------------

//template <typename T>
func (this *TComCUMvField) SetAll(p []TComMv, val TComMv, eCUMode PartSize, iPartAddr int, uiDepth uint, iPartIdx int) {
    var i uint
    p = p[iPartAddr:]
    numElements := this.m_uiNumPartition >> (2 * uiDepth)

    switch eCUMode {
    case SIZE_2Nx2N:
        for i = 0; i < numElements; i++ {
            p[i] = val
        }
    case SIZE_2NxN:
        numElements >>= 1
        for i = 0; i < numElements; i++ {
            p[i] = val
        }
    case SIZE_Nx2N:
        numElements >>= 2
        for i = 0; i < numElements; i++ {
            p[i] = val
            p[i+2*numElements] = val
        }
    case SIZE_NxN:
        numElements >>= 2
        for i = 0; i < numElements; i++ {
            p[i] = val
        }
    case SIZE_2NxnU:
        iCurrPartNumQ := numElements >> 2
        if iPartIdx == 0 {
            pT := p
            pT2 := p[iCurrPartNumQ:]
            for i = 0; i < (iCurrPartNumQ >> 1); i++ {
                pT[i] = val
                pT2[i] = val
            }
        } else {
            pT := p
            for i = 0; i < (iCurrPartNumQ >> 1); i++ {
                pT[i] = val
            }

            pT = p[iCurrPartNumQ:]
            for i = 0; i < ((iCurrPartNumQ >> 1) + (iCurrPartNumQ << 1)); i++ {
                pT[i] = val
            }
        }
    case SIZE_2NxnD:
        iCurrPartNumQ := numElements >> 2
        if iPartIdx == 0 {
            pT := p
            for i = 0; i < ((iCurrPartNumQ >> 1) + (iCurrPartNumQ << 1)); i++ {
                pT[i] = val
            }
            pT = p[(numElements - iCurrPartNumQ):]
            for i = 0; i < (iCurrPartNumQ >> 1); i++ {
                pT[i] = val
            }
        } else {
            pT := p
            pT2 := p[iCurrPartNumQ:]
            for i = 0; i < (iCurrPartNumQ >> 1); i++ {
                pT[i] = val
                pT2[i] = val
            }
        }
    case SIZE_nLx2N:
        iCurrPartNumQ := numElements >> 2
        if iPartIdx == 0 {
            pT := p
            pT2 := p[(iCurrPartNumQ << 1):]
            pT3 := p[(iCurrPartNumQ >> 1):]
            pT4 := p[(iCurrPartNumQ<<1)+(iCurrPartNumQ>>1):]

            for i = 0; i < (iCurrPartNumQ >> 2); i++ {
                pT[i] = val
                pT2[i] = val
                pT3[i] = val
                pT4[i] = val
            }
        } else {
            pT := p
            pT2 := p[(iCurrPartNumQ << 1):]
            for i = 0; i < (iCurrPartNumQ >> 2); i++ {
                pT[i] = val
                pT2[i] = val
            }

            pT = p[(iCurrPartNumQ >> 1):]
            pT2 = p[(iCurrPartNumQ<<1)+(iCurrPartNumQ>>1):]
            for i = 0; i < ((iCurrPartNumQ >> 2) + iCurrPartNumQ); i++ {
                pT[i] = val
                pT2[i] = val
            }
        }
    case SIZE_nRx2N:
        iCurrPartNumQ := numElements >> 2
        if iPartIdx == 0 {
            pT := p
            pT2 := p[(iCurrPartNumQ << 1):]
            for i = 0; i < ((iCurrPartNumQ >> 2) + iCurrPartNumQ); i++ {
                pT[i] = val
                pT2[i] = val
            }

            pT = p[iCurrPartNumQ+(iCurrPartNumQ>>1):]
            pT2 = p[numElements-iCurrPartNumQ+(iCurrPartNumQ>>1):]
            for i = 0; i < (iCurrPartNumQ >> 2); i++ {
                pT[i] = val
                pT2[i] = val
            }
        } else {
            pT := p
            pT2 := p[(iCurrPartNumQ >> 1):]
            pT3 := p[(iCurrPartNumQ << 1):]
            pT4 := p[(iCurrPartNumQ<<1)+(iCurrPartNumQ>>1):]
            for i = 0; i < (iCurrPartNumQ >> 2); i++ {
                pT[i] = val
                pT2[i] = val
                pT3[i] = val
                pT4[i] = val
            }
        }
    default:
        //assert(0);
    }
}

func (this *TComCUMvField) SetAll2(p []int8, val int8, eCUMode PartSize, iPartAddr int, uiDepth uint, iPartIdx int) {
    var i uint
    p = p[iPartAddr:]
    numElements := this.m_uiNumPartition >> (2 * uiDepth)

    switch eCUMode {
    case SIZE_2Nx2N:
        for i = 0; i < numElements; i++ {
            p[i] = val
        }
    case SIZE_2NxN:
        numElements >>= 1
        for i = 0; i < numElements; i++ {
            p[i] = val
        }
    case SIZE_Nx2N:
        numElements >>= 2
        for i = 0; i < numElements; i++ {
            p[i] = val
            p[i+2*numElements] = val
        }
    case SIZE_NxN:
        numElements >>= 2
        for i = 0; i < numElements; i++ {
            p[i] = val
        }
    case SIZE_2NxnU:
        iCurrPartNumQ := numElements >> 2
        if iPartIdx == 0 {
            pT := p
            pT2 := p[iCurrPartNumQ:]
            for i = 0; i < (iCurrPartNumQ >> 1); i++ {
                pT[i] = val
                pT2[i] = val
            }
        } else {
            pT := p
            for i = 0; i < (iCurrPartNumQ >> 1); i++ {
                pT[i] = val
            }

            pT = p[iCurrPartNumQ:]
            for i = 0; i < ((iCurrPartNumQ >> 1) + (iCurrPartNumQ << 1)); i++ {
                pT[i] = val
            }
        }
    case SIZE_2NxnD:
        iCurrPartNumQ := numElements >> 2
        if iPartIdx == 0 {
            pT := p
            for i = 0; i < ((iCurrPartNumQ >> 1) + (iCurrPartNumQ << 1)); i++ {
                pT[i] = val
            }
            pT = p[(numElements - iCurrPartNumQ):]
            for i = 0; i < (iCurrPartNumQ >> 1); i++ {
                pT[i] = val
            }
        } else {
            pT := p
            pT2 := p[iCurrPartNumQ:]
            for i = 0; i < (iCurrPartNumQ >> 1); i++ {
                pT[i] = val
                pT2[i] = val
            }
        }
    case SIZE_nLx2N:
        iCurrPartNumQ := numElements >> 2
        if iPartIdx == 0 {
            pT := p
            pT2 := p[(iCurrPartNumQ << 1):]
            pT3 := p[(iCurrPartNumQ >> 1):]
            pT4 := p[(iCurrPartNumQ<<1)+(iCurrPartNumQ>>1):]

            for i = 0; i < (iCurrPartNumQ >> 2); i++ {
                pT[i] = val
                pT2[i] = val
                pT3[i] = val
                pT4[i] = val
            }
        } else {
            pT := p
            pT2 := p[(iCurrPartNumQ << 1):]
            for i = 0; i < (iCurrPartNumQ >> 2); i++ {
                pT[i] = val
                pT2[i] = val
            }

            pT = p[(iCurrPartNumQ >> 1):]
            pT2 = p[(iCurrPartNumQ<<1)+(iCurrPartNumQ>>1):]
            for i = 0; i < ((iCurrPartNumQ >> 2) + iCurrPartNumQ); i++ {
                pT[i] = val
                pT2[i] = val
            }
        }
    case SIZE_nRx2N:
        iCurrPartNumQ := numElements >> 2
        if iPartIdx == 0 {
            pT := p
            pT2 := p[(iCurrPartNumQ << 1):]
            for i = 0; i < ((iCurrPartNumQ >> 2) + iCurrPartNumQ); i++ {
                pT[i] = val
                pT2[i] = val
            }

            pT = p[iCurrPartNumQ+(iCurrPartNumQ>>1):]
            pT2 = p[numElements-iCurrPartNumQ+(iCurrPartNumQ>>1):]
            for i = 0; i < (iCurrPartNumQ >> 2); i++ {
                pT[i] = val
                pT2[i] = val
            }
        } else {
            pT := p
            pT2 := p[(iCurrPartNumQ >> 1):]
            pT3 := p[(iCurrPartNumQ << 1):]
            pT4 := p[(iCurrPartNumQ<<1)+(iCurrPartNumQ>>1):]
            for i = 0; i < (iCurrPartNumQ >> 2); i++ {
                pT[i] = val
                pT2[i] = val
                pT3[i] = val
                pT4[i] = val
            }
        }
    default:
        //assert(0);
    }
}

func (this *TComCUMvField) SetAllMv(mv TComMv, eCUMode PartSize, iPartAddr int, uiDepth uint, iPartIdx int) {
    this.SetAll(this.m_pcMv, mv, eCUMode, iPartAddr, uiDepth, iPartIdx)
}
func (this *TComCUMvField) SetAllMvd(mvd TComMv, eCUMode PartSize, iPartAddr int, uiDepth uint, iPartIdx int) {
    this.SetAll(this.m_pcMvd, mvd, eCUMode, iPartAddr, uiDepth, iPartIdx)
}
func (this *TComCUMvField) SetAllRefIdx(iRefIdx int8, eCUMode PartSize, iPartAddr int, uiDepth uint, iPartIdx int) {
    this.SetAll2(this.m_piRefIdx, iRefIdx, eCUMode, iPartAddr, uiDepth, iPartIdx)
}
func (this *TComCUMvField) SetAllMvField(mvField *TComMvField, eCUMode PartSize, iPartAddr int, uiDepth uint, iPartIdx int) {
    this.SetAllMv(mvField.GetMv(), eCUMode, iPartAddr, uiDepth, iPartIdx)
    this.SetAllRefIdx(mvField.GetRefIdx(), eCUMode, iPartAddr, uiDepth, iPartIdx)
}

func (this *TComCUMvField) SetNumPartition(uiNumPart uint) {
    this.m_uiNumPartition = uiNumPart
}

func (this *TComCUMvField) LinkToWithOffset(src *TComCUMvField, offset int) {
    this.m_pcMv = src.GetMvs(offset)
    this.m_pcMvd = src.GetMvds(offset)
    this.m_piRefIdx = src.GetRefIdxs(offset)
}

func (this *TComCUMvField) Compress(pePredMode []PredMode, scale int) {
    N := scale * scale
    //assert( N > 0 && N <= m_uiNumPartition);

    for uiPartIdx := 0; uiPartIdx < int(this.m_uiNumPartition); uiPartIdx += N {
        //cMv :=NewTComMv(0,0);
        //predMode := MODE_INTRA;
        //iRefIdx := 0;

        cMv := this.m_pcMv[uiPartIdx]
        predMode := pePredMode[uiPartIdx]
        iRefIdx := this.m_piRefIdx[uiPartIdx]
        for i := 0; i < N; i++ {
            this.m_pcMv[uiPartIdx+i] = cMv
            pePredMode[uiPartIdx+i] = predMode
            this.m_piRefIdx[uiPartIdx+i] = iRefIdx
        }
    }
}
