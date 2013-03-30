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
    "io"
)

type InputByteStream struct {
    //private:
    m_NumFutureBytes uint     /* number of valid bytes in m_FutureBytes */
    m_FutureBytes    uint32   /* bytes that have been peeked */
    m_Input          io.Reader /* Input stream to read from */
}

//public:
/**
 * Create a bytestream reader that will extract bytes from
 * istream.
 *
 * NB, it isn't safe to access istream while in use by a
 * InputByteStream.
 *
 * Side-effects: the exception mask of istream is set to eofbit
 */
func NewInputByteStream(istream io.Reader) *InputByteStream {
    return &InputByteStream{m_NumFutureBytes: 0, m_FutureBytes: 0, m_Input: istream}
}

/**
 * Reset the internal state.  Must be called if input stream is
 * modified externally to this class
 */
func (this *InputByteStream) Reset() {
    this.m_NumFutureBytes = 0
    this.m_FutureBytes = 0
}

/**
 * returns true if an EOF will be encountered within the next
 * n bytes.
 */
func (this *InputByteStream) EofBeforeNBytes(n uint) bool {
    if n > 4 {
        fmt.Printf("n must be smaller or equal to 4\n")
        return false
    }

    if this.m_NumFutureBytes >= n {
        return false
    }

    n -= this.m_NumFutureBytes

    buf := make([]byte, 1)

    for i := uint(0); i < n; i++ {
        _, err := this.m_Input.Read(buf)
        if err == io.EOF {
            return true
        }
        this.m_FutureBytes = (this.m_FutureBytes << 8) | uint32(buf[0])
        this.m_NumFutureBytes++
    }

    return false
}

/**
 * return the next n bytes in the stream without advancing
 * the stream pointer.
 *
 * Returns: an unsigned integer representing an n byte bigendian
 * word.
 *
 * If an attempt is made to read past EOF, an n-byte word is
 * returned, but the portion that required input bytes beyond EOF
 * is undefined.
 *
 */
func (this *InputByteStream) PeekBytes(n uint) uint32 {
    this.EofBeforeNBytes(n)
    return this.m_FutureBytes >> uint32(8*(this.m_NumFutureBytes-n))
}

/**
 * consume and return one byte from the input.
 *
 * If bytestream is already at EOF prior to a call to readByte(),
 * an exception std::ios_base::failure is thrown.
 */
func (this *InputByteStream) ReadByte() (byte, error) {
    wanted_byte := make([]byte, 1)

    if this.m_NumFutureBytes == 0 {
        _, err := this.m_Input.Read(wanted_byte)
        return wanted_byte[0], err
    }
    this.m_NumFutureBytes--
    wanted_byte[0] = byte(this.m_FutureBytes >> uint32(8*this.m_NumFutureBytes))
    this.m_FutureBytes &= ^(0xff << uint32(8*this.m_NumFutureBytes))
    return wanted_byte[0], nil
}

/**
 * consume and return n bytes from the input.  n bytes from
 * bytestream are interpreted as bigendian when assembling
 * the return value.
 */
func (this *InputByteStream) ReadBytes(n uint) (uint32, error) {
    var val uint32 = 0
    for i := uint(0); i < n; i++ {
        b, err := this.ReadByte()
        if err != nil {
            return val, err
        } else {
            val = (val << 8) | uint32(b)
        }
    }

    return val, nil
}

func (this *InputByteStream) byteStreamNALUnit(nalUnit *list.List, stats *AnnexBStats) (err error) {
    /* At the beginning of the decoding process, the decoder initialises its
     * current position in the byte stream to the beginning of the byte stream.
     * It then extracts and discards each leading_zero_8bits syntax element (if
     * present), moving the current position in the byte stream forward one
     * byte at a time, until the current position in the byte stream is such
     * that the next four bytes in the bitstream form the four-byte sequence
     * 0x00000001.
     */
    //fmt.Printf("0\n");
    for (this.EofBeforeNBytes(24/8) || this.PeekBytes(24/8) != 0x000001) &&
        (this.EofBeforeNBytes(32/8) || this.PeekBytes(32/8) != 0x00000001) {
        leading_zero_8bits, err := this.ReadByte()
        if leading_zero_8bits != 0 || err != nil {
            err = errors.New("leading_zero_8bits!=0 || err!=nil")
            return err
        }
        stats.NumLeadingZero8BitsBytes++
    }
    //fmt.Printf("1\n");
    /* 1. When the next four bytes in the bitstream form the four-byte sequence
     * 0x00000001, the next byte in the byte stream (which is a zero_byte
     * syntax element) is extracted and discarded and the current position in
     * the byte stream is set equal to the position of the byte following this
     * discarded byte.
     */
    /* NB, the previous step guarantees this will succeed -- if EOF was
     * encountered, an exception will stop execution getting this far */
    if this.PeekBytes(24/8) != 0x000001 {
        zero_byte, err := this.ReadByte()
        if zero_byte != 0 || err != nil {
            err = errors.New("zero_byte!=0 || err!=nil")
            return err
        }
        stats.NumZeroByteBytes++
    }
    //fmt.Printf("2\n");
    /* 2. The next three-byte sequence in the byte stream (which is a
     * start_code_prefix_one_3bytes) is extracted and discarded and the current
     * position in the byte stream is set equal to the position of the byte
     * following this three-byte sequence.
     */
    /* NB, (1) guarantees that the next three bytes are 0x00 00 01 */
    start_code_prefix_one_3bytes, err := this.ReadBytes(24 / 8)
    if start_code_prefix_one_3bytes != 0x000001 || err != nil {
        err = errors.New("start_code_prefix_one_3bytes != 0x000001 || err!=nil")
        return err
    }
    stats.NumStartCodePrefixBytes += 3
    //fmt.Printf("3\n");
    /* 3. NumBytesInNALunit is set equal to the number of bytes starting with
     * the byte at the current position in the byte stream up to and including
     * the last byte that precedes the location of any of the following
     * conditions:
     *   a. A subsequent byte-aligned three-byte sequence equal to 0x000000, or
     *   b. A subsequent byte-aligned three-byte sequence equal to 0x000001, or
     *   c. The end of the byte stream, as determined by unspecified means.
     */
    /* 4. NumBytesInNALunit bytes are removed from the bitstream and the
     * current position in the byte stream is advanced by NumBytesInNALunit
     * bytes. This sequence of bytes is nal_unit( NumBytesInNALunit ) and is
     * decoded using the NAL unit decoding process
     */
    /* NB, (unsigned)x > 2 implies n!=0 && n!=1 */
    for this.EofBeforeNBytes(24/8) || this.PeekBytes(24/8) > 2 {
        b, err := this.ReadByte()
        if err != nil {
            return err
        }
        nalUnit.PushBack(b)
    }
    //fmt.Printf("5\n");
    /* 5. When the current position in the byte stream is:
     *  - not at the end of the byte stream (as determined by unspecified means)
     *  - and the next bytes in the byte stream do not start with a three-byte
     *    sequence equal to 0x000001
     *  - and the next bytes in the byte stream do not start with a four byte
     *    sequence equal to 0x00000001,
     * the decoder extracts and discards each trailing_zero_8bits syntax
     * element, moving the current position in the byte stream forward one byte
     * at a time, until the current position in the byte stream is such that:
     *  - the next bytes in the byte stream form the four-byte sequence
     *    0x00000001 or
     *  - the end of the byte stream has been encountered (as determined by
     *    unspecified means).
     */
    /* NB, (3) guarantees there are at least three bytes available or none */
    for (this.EofBeforeNBytes(24/8) || this.PeekBytes(24/8) != 0x000001) &&
        (this.EofBeforeNBytes(32/8) || this.PeekBytes(32/8) != 0x00000001) {
        trailing_zero_8bits, err := this.ReadByte()
        if trailing_zero_8bits != 0 || err != nil {
            err = errors.New("trailing_zero_8bits!=0 || err!=nil")
            return err
        }
        stats.NumTrailingZero8BitsBytes++
    }

    return nil
}

func (this *InputByteStream) ByteStreamNALUnit(nalUnit *list.List, stats *AnnexBStats) (bool, error) {
    var err error
    eof := false

    if err = this.byteStreamNALUnit(nalUnit, stats); err != nil {
        eof = true
    }

    stats.NumBytesInNALUnit = uint(nalUnit.Len())
    return eof, err
}

func (this *InputByteStream) ByteStreamNALUnits(nalUnits *list.List) error{
	var stats AnnexBStats
	var err error
	eof := false
	
	for !eof {
		nalUnit := list.New()
		eof, err = this.ByteStreamNALUnit(nalUnit, &stats);
		if nalUnit.Len() == 0 {
			return errors.New("Warning: Attempt to decode an empty NAL unit\n");
		}else if err !=nil  {
			return err;
		}else{
			nalUnits.PushBack(nalUnit);
		}
	}
	
	return nil;
}


/**
 * Statistics associated with AnnexB bytestreams
 */
type AnnexBStats struct {
    NumLeadingZero8BitsBytes  uint
    NumZeroByteBytes          uint
    NumStartCodePrefixBytes   uint
    NumBytesInNALUnit         uint
    NumTrailingZero8BitsBytes uint
}

func (this *AnnexBStats) Add(rhs *AnnexBStats) {
    this.NumLeadingZero8BitsBytes += rhs.NumLeadingZero8BitsBytes
    this.NumZeroByteBytes += rhs.NumZeroByteBytes
    this.NumStartCodePrefixBytes += rhs.NumStartCodePrefixBytes
    this.NumBytesInNALUnit += rhs.NumBytesInNALUnit
    this.NumTrailingZero8BitsBytes += rhs.NumTrailingZero8BitsBytes
}
