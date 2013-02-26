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
// Class definition
// ====================================================================================================================

/// basic motion vector class
type TComMv struct {
    //private:
    m_iHor int16 ///< horizontal component of motion vector
    m_iVer int16 ///< vertical component of motion vector
}

func NewTComMv(iHor, iVer int16) *TComMv {
    return &TComMv{iHor, iVer}
}

// ------------------------------------------------------------------------------------------------------------------
// constructors
// ------------------------------------------------------------------------------------------------------------------

// ------------------------------------------------------------------------------------------------------------------
// set
// ------------------------------------------------------------------------------------------------------------------

func (this *TComMv) Set(iHor, iVer int16) {
    this.m_iHor = iHor
    this.m_iVer = iVer
}
func (this *TComMv) SetHor(i int16) {
    this.m_iHor = i
}
func (this *TComMv) SetVer(i int16) {
    this.m_iVer = i
}
func (this *TComMv) SetZero() {
    this.m_iHor = 0
    this.m_iVer = 0
}

// ------------------------------------------------------------------------------------------------------------------
// get
// ------------------------------------------------------------------------------------------------------------------

func (this *TComMv) GetHor() int16 {
    return this.m_iHor
}
func (this *TComMv) GetVer() int16 {
    return this.m_iVer
}
func (this *TComMv) GetAbsHor() int16 {
    if this.m_iHor < 0 {
        return -this.m_iHor
    }

    return this.m_iHor
}
func (this *TComMv) GetAbsVer() int16 {
    if this.m_iVer < 0 {
        return -this.m_iVer
    }

    return this.m_iVer
}

// ------------------------------------------------------------------------------------------------------------------
// operations
// ------------------------------------------------------------------------------------------------------------------

func (this *TComMv) AddMv(rcMv TComMv) {
    this.m_iHor += rcMv.m_iHor
    this.m_iVer += rcMv.m_iVer
}

func (this *TComMv) SubMv(rcMv TComMv) {
    this.m_iHor -= rcMv.m_iHor
    this.m_iVer -= rcMv.m_iVer
}

func (this *TComMv) RightShift(i uint) {
    this.m_iHor >>= i
    this.m_iVer >>= i
}

func (this *TComMv) LeftShift(i uint) {
    this.m_iHor <<= i
    this.m_iVer <<= i
}

func SubMvs(lcMv, rcMv TComMv) TComMv {
    return TComMv{lcMv.m_iHor - rcMv.m_iHor, lcMv.m_iVer - rcMv.m_iVer}
}

func AddMvs(lcMv, rcMv TComMv) TComMv {
    return TComMv{lcMv.m_iHor + rcMv.m_iHor, lcMv.m_iVer + rcMv.m_iVer}
}

func EqualMvs(lcMv, rcMv TComMv) bool {
    return (lcMv.m_iHor == rcMv.m_iHor && lcMv.m_iVer == rcMv.m_iVer)
}

func NotEqualMvs(lcMv, rcMv TComMv) bool {
    return (lcMv.m_iHor != rcMv.m_iHor || lcMv.m_iVer != rcMv.m_iVer)
}

func (this *TComMv) ScaleMv(iScale int) TComMv{
    mvx := (iScale*int(this.GetHor()) + 127 + int(B2U(iScale*int(this.GetHor()) < 0))) >> 8

    if mvx < -32768 {
        mvx = -32768
    } else if mvx > 32767 {
        mvx = 32767
    }

    mvy := (iScale*int(this.GetVer()) + 127 + int(B2U(iScale*int(this.GetVer()) < 0))) >> 8

    if mvy < -32768 {
        mvy = -32768
    } else if mvy > 32767 {
        mvy = 32767
    }

    return TComMv{m_iHor: int16(mvx), m_iVer: int16(mvy)}
}


func (this *TComMv) ShiftMv(uiShift uint){
    this.m_iHor <<= uiShift;
	this.m_iVer <<= uiShift;
}
