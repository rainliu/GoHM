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
    "io"
)

type AccessUnit struct {
    *list.List
}

func NewAccessUnit() *AccessUnit {
    return &AccessUnit{list.New()}
}

type AccessUnits struct {
    *list.List
}

func NewAccessUnits() *AccessUnits {
    return &AccessUnits{list.New()}
}

/**
 * write all NALunits in au to bytestream out in a manner satisfying
 * AnnexB of AVC.  NALunits are written in the order they are found in au.
 * the zero_byte word is appended to:
 *  - the initial startcode in the access unit,
 *  - any SPS/PPS nal units
 */
var start_code_prefix = [4]byte{0, 0, 0, 1}

func WriteAnnexB(out io.Writer, au *AccessUnit) *list.List {
    annexBsizes := list.New()

    for it := au.Front(); it != nil; it = it.Next() {
        nalu := it.Value.(*NALUnitEBSP)
        size := uint(0) /* size of annexB unit in bytes */

        if it == au.Front() || nalu.GetNalUnitType() == TLibCommon.NAL_UNIT_SPS || nalu.GetNalUnitType() == TLibCommon.NAL_UNIT_PPS {
            /* From AVC, When any of the following conditions are fulfilled, the
             * zero_byte syntax element shall be present:
             *  - the nal_unit_type within the nal_unit() is equal to 7 (sequence
             *    parameter set) or 8 (picture parameter set),
             *  - the byte stream NAL unit syntax structure contains the first NAL
             *    unit of an access unit in decoding order, as specified by subclause
             *    7.4.1.2.3.
             */
            out.Write(start_code_prefix[:])
            size += 4
        } else {
            out.Write(start_code_prefix[1:])
            size += 3
        }

        var buf [1]byte
        for e := nalu.m_Bitstream.GetFIFO().Front(); e != nil; e = e.Next() {
            buf[0] = e.Value.(byte)
            out.Write(buf[:])
        }
        //out << nalu.m_nalUnitData.str();
        size += uint(nalu.m_Bitstream.GetFIFO().Len())

        annexBsizes.PushBack(size)
    }

    return annexBsizes
}
