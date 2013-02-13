package TLibEncoder

import (
    //"io"
    "gohm/TLibCommon"
)

var emulation_prevention_three_byte = [1]byte{3}

/**
 * A convenience wrapper to NALUnit that also provides a
 * bitstream object.
 */
type OutputNALUnit struct {
    TLibCommon.NALUnit
    m_Bitstream *TLibCommon.TComOutputBitstream
}

/**
 * construct an OutputNALunit structure with given header values and
 * storage for a bitstream.  Upon construction the NALunit header is
 * written to the bitstream.
 */
func NewOutputNALUnit(nalUnitType TLibCommon.NalUnitType, temporalID, reserved_zero_6bits uint) *OutputNALUnit {
    out := &OutputNALUnit{}

    out.NALUnit.SetNalUnitType(nalUnitType)
    out.NALUnit.SetTemporalId(temporalID)
    out.NALUnit.SetReservedZero6Bits(reserved_zero_6bits)
    out.m_Bitstream = nil


    return out
}

func (this *OutputNALUnit) Copy(src *TLibCommon.NALUnit) {
    this.m_Bitstream.Clear()

    this.NALUnit.SetNalUnitType(src.GetNalUnitType())
    this.NALUnit.SetReservedZero6Bits(src.GetReservedZero6Bits())
    this.NALUnit.SetTemporalId(src.GetTemporalId())
    //static_cast<NALUnit*>(this)->operator=(src);
}

func (this *OutputNALUnit) WriteNalUnitHeader() { // nal_unit_header()
    //bsNALUHeader := TLibCommon.NewTComOutputBitstream();//*TLibCommon.TComOutputBitstream;

    this.m_Bitstream.Write(0, 1)                           // forbidden_zero_bit
    this.m_Bitstream.Write(uint(this.GetNalUnitType()), 6) // nal_unit_type
    this.m_Bitstream.Write(this.GetReservedZero6Bits(), 6) // nuh_reserved_zero_6bits
    this.m_Bitstream.Write(this.GetTemporalId()+1, 3)      // nuh_temporal_id_plus1

    /*var buf [1]byte;
      for e:= bsNALUHeader.GetFIFO().Front(); e!=nil; e=e.Next() {
      	buf[0] = e.Value.(byte);
      	out.Write(buf[:]);
      }	*/
    //out.write(bsNALUHeader.getByteStream(), bsNALUHeader.getByteStreamLength());
}

func (this *OutputNALUnit) Write() {
    this.WriteNalUnitHeader()
    /* write out rsbp_byte's, inserting any required
     * emulation_prevention_three_byte's */
    /* 7.4.1 ...
     * emulation_prevention_three_byte is a byte equal to 0x03. When an
     * emulation_prevention_three_byte is present in the NAL unit, it shall be
     * discarded by the decoding process.
     * The last byte of the NAL unit shall not be equal to 0x00.
     * Within the NAL unit, the following three-byte sequences shall not occur at
     * any byte-aligned position:
     *  - 0x000000
     *  - 0x000001
     *  - 0x000002
     * Within the NAL unit, any four-byte sequence that starts with 0x000003
     * other than the following sequences shall not occur at any byte-aligned
     * position:
     *  - 0x00000300
     *  - 0x00000301
     *  - 0x00000302
     *  - 0x00000303
     */
    rbsp := this.m_Bitstream.GetFIFO()
    var v0, v1 byte
    for it := rbsp.Front(); it != nil; it = it.Next() {
        /* 1) find the next emulated 00 00 {00,01,02,03}
         * 2a) if not found, write all remaining bytes out, stop.
         * 2b) otherwise, write all non-emulated bytes out
         * 3) insert emulation_prevention_three_byte
         */
        found := it
        for {
            /* NB, end()-1, prevents finding a trailing two byte sequence */
            //found = search_n(found, rbsp.end()-1, 2, 0);
            for found != rbsp.Back() {
                v0 = found.Value.(byte)
                if found.Next() != rbsp.Back() {
                    v1 = found.Next().Value.(byte)
                } else {
                    v1 = 0xFF
                }
                found = found.Next()

                if v0 == 0 && v1 == 0 {
                    break
                }
            }

            found = found.Next()

            /* if not found, found == end, otherwise found = second zero byte */
            if found == nil {
                break
            }

            found = found.Next()

            if found.Value.(byte) <= 3 {
                break
            }
        }

        it = found
        if found != nil {
            it = rbsp.InsertBefore(emulation_prevention_three_byte[0], found)
        }
    }

    /*
      var buf [1]byte;
      for e:= rbsp.Front(); e!=nil; e=e.Next() {
      	buf[0] = e.Value.(byte);
      	out.Write(buf[:]);
      }	*/
    //out.write((Char*)&(*rbsp.begin()), rbsp.end() - rbsp.begin());

    /* 7.4.1.1
     * ... when the last byte of the RBSP data is equal to 0x00 (which can
     * only occur when the RBSP ends in a cabac_zero_word), a final byte equal
     * to 0x03 is appended to the end of the data.
     */
    if rbsp.Back().Value.(byte) == 0x00 {
        rbsp.PushBack(emulation_prevention_three_byte[0])
        //out.Write(emulation_prevention_three_byte[:]);
    }
}

/**
 * Copy NALU from naluSrc to naluDest
 */
func (this *OutputNALUnit) CopyNaluData(naluSrc *OutputNALUnit) {
    this.SetNalUnitType(naluSrc.GetNalUnitType())
    this.SetReservedZero6Bits(naluSrc.GetReservedZero6Bits())
    this.SetTemporalId(naluSrc.GetTemporalId())
    this.m_Bitstream = naluSrc.m_Bitstream
}

/**
 * A single NALunit, with complete payload in EBSP format.
 */
/* 
type NALUnitEBSP struct{
  TLibCommon.NALUnit
  //m_nalUnitData io.Writer;
}*/
/**
 * convert the OutputNALUnit #nalu# into EBSP format by writing out
 * the NALUnit header, then the rbsp_bytes including any
 * emulation_prevention_three_byte symbols.
 */
/*
func NewNALUnitEBSP(nalu *OutputNALUnit) *NALUnitEBSP{
  nalu.Write();
}*/
