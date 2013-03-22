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

package TLibDecoder

import (
    "container/list"
    "errors"
    "fmt"
    "gohm/TLibCommon"
)

/**
 * A convenience wrapper to NALUnit that also provides a
 * bitstream object.
 */
type InputNALUnit struct {
    TLibCommon.NALUnit
    m_Bitstream *TLibCommon.TComInputBitstream
}

func NewInputNALUnit() *InputNALUnit {
    return &InputNALUnit{}
}

func (this *InputNALUnit) Read(nalUnitBuf *list.List) *list.List {
    /* perform anti-emulation prevention */
    pcBitstream := TLibCommon.NewTComInputBitstream(nil)

    firstByte := nalUnitBuf.Front().Value.(byte)
    oldNalUnitBuf := this.convertPayloadToRBSP(nalUnitBuf, pcBitstream, (firstByte&64) == 0)

    this.m_Bitstream = TLibCommon.NewTComInputBitstream(nalUnitBuf)

    this.m_Bitstream.SetEmulationPreventionByteLocation(pcBitstream.GetEmulationPreventionByteLocation());

    this.readNalUnitHeader()

    return oldNalUnitBuf
}

func (this *InputNALUnit) GetBitstream() *TLibCommon.TComInputBitstream {
    return this.m_Bitstream
}

func (this *InputNALUnit) SetBitstream(bitstream *TLibCommon.TComInputBitstream) {
    this.m_Bitstream = bitstream
}

func (this *InputNALUnit) convertPayloadToRBSP(nalUnitBuf *list.List, pcBitstream *TLibCommon.TComInputBitstream, isVclNalUnit bool) *list.List {
    zeroCount := 0
    it_write := list.New()
    oldBuf := list.New()

    pos := uint(0);
    pcBitstream.ClearEmulationPreventionByteLocation();

    for e := nalUnitBuf.Front(); e != nil; e, pos = e.Next(), pos+1 {
        //assert(zeroCount < 2 || *it_read >= 0x03);
        it_read := e.Value.(byte)
        oldBuf.PushBack(it_read)
        if zeroCount == 2 && it_read == 0x03 {
            pcBitstream.PushEmulationPreventionByteLocation( pos );
            pos++;

            zeroCount = 0

            e = e.Next()
            if e == nil {
                break
            } else {
                it_read = e.Value.(byte)
                oldBuf.PushBack(it_read)
            }
        }

        if it_read == 0x00 {
            zeroCount++
        } else {
            zeroCount = 0
        }
        it_write.PushBack(it_read)
    }

    //assert(zeroCount == 0);
    if isVclNalUnit {
        // Remove cabac_zero_word from payload if present
        n := 0

        e := it_write.Back()
        it_read := e.Value.(byte)
        for it_read == 0x00 {
            it_write.Remove(e)
            n++
            e = it_write.Back()
            if e!=nil{
            	it_read = e.Value.(byte)
            }else{
            	break;
            }
        }

        if n > 0 {
            fmt.Printf("\nDetected %d instances of cabac_zero_word", n/2)
        }
    }

    nalUnitBuf.Init() // = .resize(it_write - nalUnitBuf.begin());
    for e := it_write.Front(); e != nil; e = e.Next() {
        it_read := e.Value.(byte)
        nalUnitBuf.PushBack(it_read)
    }

    return oldBuf
}

func (this *InputNALUnit) readNalUnitHeader() error {
    bs := this.m_Bitstream

    forbidden_zero_bit := bs.ReadBits(1) // forbidden_zero_bit
    if forbidden_zero_bit != 0 {
        return errors.New("forbidden_zero_bit!=0")
    }

    this.SetNalUnitType(TLibCommon.NalUnitType(bs.ReadBits(6))) // nal_unit_type

    this.SetReservedZero6Bits(bs.ReadBits(6)) // nuh_reserved_zero_6bits
    if this.GetReservedZero6Bits() != 0 {
        return errors.New("m_reservedZero6Bits!=0")
    }

    this.SetTemporalId(bs.ReadBits(3) - 1) // nuh_temporal_id_plus1

    if this.GetTemporalId() != 0 {
        if !(this.GetNalUnitType() != TLibCommon.NAL_UNIT_CODED_SLICE_BLA &&
            this.GetNalUnitType() != TLibCommon.NAL_UNIT_CODED_SLICE_BLANT &&
            this.GetNalUnitType() != TLibCommon.NAL_UNIT_CODED_SLICE_BLA_N_LP &&
            this.GetNalUnitType() != TLibCommon.NAL_UNIT_CODED_SLICE_IDR &&
            this.GetNalUnitType() != TLibCommon.NAL_UNIT_CODED_SLICE_IDR_N_LP &&
            this.GetNalUnitType() != TLibCommon.NAL_UNIT_CODED_SLICE_CRA &&
            this.GetNalUnitType() != TLibCommon.NAL_UNIT_VPS &&
            this.GetNalUnitType() != TLibCommon.NAL_UNIT_SPS &&
            this.GetNalUnitType() != TLibCommon.NAL_UNIT_EOS &&
            this.GetNalUnitType() != TLibCommon.NAL_UNIT_EOB) {
            return errors.New("Wrong this.GetNalUnitType() in readNalUnitHeader")
        }
    } else {
        if !(this.GetNalUnitType() != TLibCommon.NAL_UNIT_CODED_SLICE_TLA &&
            this.GetNalUnitType() != TLibCommon.NAL_UNIT_CODED_SLICE_TSA_N &&
            this.GetNalUnitType() != TLibCommon.NAL_UNIT_CODED_SLICE_STSA_R &&
            this.GetNalUnitType() != TLibCommon.NAL_UNIT_CODED_SLICE_STSA_N) {
            return errors.New("Wrong this.GetNalUnitType() in readNalUnitHeader")
        }
    }

    return nil
}
