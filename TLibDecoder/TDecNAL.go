package TLibDecoder

import (
	"container/list"
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

func (this *InputNALUnit) Read(nalUnitBuf *list.List) {
}

func (this *InputNALUnit) GetBitstream() *TLibCommon.TComInputBitstream {
    return this.m_Bitstream
}

func (this *InputNALUnit) SetBitstream(bitstream *TLibCommon.TComInputBitstream) {
    this.m_Bitstream = bitstream
}