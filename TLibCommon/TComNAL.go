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
//"container/list"
)

/**
 * Represents a single NALunit header and the associated RBSPayload
 */
type NALUnit struct {
    m_nalUnitType       NalUnitType ///< nal_unit_type
    m_temporalId        uint        ///< temporal_id
    m_reservedZero6Bits uint        ///< reserved_zero_6bits
}

/** construct an NALunit structure with given header values. */
func NewNALUnit(nalUnitType NalUnitType, temporalId, reservedZero6Bits uint) *NALUnit {
    return &NALUnit{nalUnitType, temporalId, reservedZero6Bits}
}

/** returns true if the NALunit is a slice NALunit */
func (this *NALUnit) IsSlice() bool {
    return this.m_nalUnitType == NAL_UNIT_CODED_SLICE_TRAIL_R ||
        this.m_nalUnitType == NAL_UNIT_CODED_SLICE_TRAIL_N ||
        this.m_nalUnitType == NAL_UNIT_CODED_SLICE_TLA ||
        this.m_nalUnitType == NAL_UNIT_CODED_SLICE_TSA_N ||
        this.m_nalUnitType == NAL_UNIT_CODED_SLICE_STSA_R ||
        this.m_nalUnitType == NAL_UNIT_CODED_SLICE_STSA_N ||
        this.m_nalUnitType == NAL_UNIT_CODED_SLICE_BLA ||
        this.m_nalUnitType == NAL_UNIT_CODED_SLICE_BLANT ||
        this.m_nalUnitType == NAL_UNIT_CODED_SLICE_BLA_N_LP ||
        this.m_nalUnitType == NAL_UNIT_CODED_SLICE_IDR ||
        this.m_nalUnitType == NAL_UNIT_CODED_SLICE_IDR_N_LP ||
        this.m_nalUnitType == NAL_UNIT_CODED_SLICE_CRA ||
        this.m_nalUnitType == NAL_UNIT_CODED_SLICE_RADL_N ||
        this.m_nalUnitType == NAL_UNIT_CODED_SLICE_DLP ||
        this.m_nalUnitType == NAL_UNIT_CODED_SLICE_RASL_N ||
        this.m_nalUnitType == NAL_UNIT_CODED_SLICE_TFD
}

func (this *NALUnit) IsSei() bool {
    return this.m_nalUnitType == NAL_UNIT_SEI ||
        this.m_nalUnitType == NAL_UNIT_SEI_SUFFIX
}

func (this *NALUnit) IsVcl() bool {
    return (this.m_nalUnitType < 32)
}

func (this *NALUnit) GetReservedZero6Bits() uint {
    return this.m_reservedZero6Bits
}
func (this *NALUnit) GetTemporalId() uint {
    return this.m_temporalId
}
func (this *NALUnit) SetTemporalId(temporalId uint) {
    this.m_temporalId = temporalId
}

func (this *NALUnit) SetReservedZero6Bits(reservedZero6Bits uint) {
    this.m_reservedZero6Bits = reservedZero6Bits
}

func (this *NALUnit) GetNalUnitType() NalUnitType {
    return this.m_nalUnitType
}

func (this *NALUnit) SetNalUnitType(nalUnitType NalUnitType) {
    this.m_nalUnitType = nalUnitType
}
