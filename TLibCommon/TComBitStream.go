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
    "container/list"
    "errors"
)

// ====================================================================================================================
// Class definition
// ====================================================================================================================

/// pure virtual class for basic bit handling
type TComBitIf interface {
    //public:
    WriteAlignOne()
    WriteAlignZero()
    Write(uiBits, uiNumberOfBits uint)
    ResetBits()
    GetNumberOfWrittenBits() uint
}

/// class for counting bits
type TComBitCounter struct {
    m_uiBitCounter uint
}

func NewTComBitCounter() *TComBitCounter {
    return &TComBitCounter{}
}

func (this *TComBitCounter) WriteAlignOne()  {}
func (this *TComBitCounter) WriteAlignZero() {}
func (this *TComBitCounter) Write(uiBits, uiNumberOfBits uint) {
    this.m_uiBitCounter += uiNumberOfBits
}
func (this *TComBitCounter) ResetBits()                   { this.m_uiBitCounter = 0 }
func (this *TComBitCounter) GetNumberOfWrittenBits() uint { return this.m_uiBitCounter }

/**
 * Model of a writable bitstream that accumulates bits to produce a
 * bytestream.
 */
type TComOutputBitstream struct { //public TComBitIf
    /**
     * FIFO for storage of bytes.  Use:
     *  - fifo.push_back(x) to append words
     *  - fifo.clear() to empty the FIFO
     *  - &fifo.front() to get a pointer to the data array.
     *    NB, this pointer is only valid until the next push_back()/clear()
     */
    m_fifo *list.List //std::vector<uint8_t> *m_fifo;

    m_num_held_bits uint /// number of bits not flushed to bytestream.
    m_held_bits     byte /// the bits held and not flushed to bytestream.
    /// this value is always msb-aligned, bigendian.
}

//public:
// create / destroy
func NewTComOutputBitstream() *TComOutputBitstream {
    return &TComOutputBitstream{list.New(), 0, 0}
}

// interface for encoding
/**
 * append uiNumberOfBits least significant bits of uiBits to
 * the current bitstream
 */
func (this *TComOutputBitstream) Write(uiBits, uiNumberOfBits uint) {
    //assert( uiNumberOfBits <= 32 );

    /* any modulo 8 remainder of num_total_bits cannot be written this time,
     * and will be held until next time. */
    num_total_bits := uiNumberOfBits + this.m_num_held_bits
    next_num_held_bits := num_total_bits % 8

    /* form a byte aligned word (write_bits), by concatenating any held bits
     * with the new bits, discarding the bits that will form the next_held_bits.
     * eg: H = held bits, V = n new bits        /---- next_held_bits
     * len(H)=7, len(V)=1: ... ---- HHHH HHHV . 0000 0000, next_num_held_bits=0
     * len(H)=7, len(V)=2: ... ---- HHHH HHHV . V000 0000, next_num_held_bits=1
     * if total_bits < 8, the value of v_ is not used */
    next_held_bits := byte(uiBits << (8 - next_num_held_bits))

    if (num_total_bits >> 3) == 0 {
        /* insufficient bits accumulated to write out, append new_held_bits to
         * current held_bits */
        /* NB, this requires that v only contains 0 in bit positions {31..n} */
        this.m_held_bits |= next_held_bits
        this.m_num_held_bits = next_num_held_bits
        return
    }

    /* topword serves to justify held_bits to align with the msb of uiBits */
    topword := (uiNumberOfBits - uint(next_num_held_bits)) & (^(uint(1 << 3) -1));
    write_bits := (uint(this.m_held_bits) << topword) | (uiBits >> uint(next_num_held_bits))

    switch num_total_bits >> 3 {
    case 4:
        this.m_fifo.PushBack(byte(write_bits >> 24))
        fallthrough
    case 3:
        this.m_fifo.PushBack(byte(write_bits >> 16))
        fallthrough
    case 2:
        this.m_fifo.PushBack(byte(write_bits >> 8))
        fallthrough
    case 1:
        this.m_fifo.PushBack(byte(write_bits))
    }

    this.m_held_bits = next_held_bits
    this.m_num_held_bits = next_num_held_bits
}

/** insert one bits until the bitstream is byte-aligned */
func (this *TComOutputBitstream) WriteAlignOne() {
    num_bits := uint(this.GetNumBitsUntilByteAligned())
    this.Write((1<<num_bits)-1, num_bits)
    return
}

/** insert zero bits until the bitstream is byte-aligned */
func (this *TComOutputBitstream) WriteAlignZero() {
    if 0 == this.m_num_held_bits {
        return
    }
    this.m_fifo.PushBack(byte(this.m_held_bits))
    this.m_held_bits = 0
    this.m_num_held_bits = 0
}

/** this function should never be called */
func (this *TComOutputBitstream) ResetBits() {
    //do nothing
}

// utility functions

/**
 * Return a pointer to the start of the byte-stream buffer.
 * Pointer is valid until the next write/flush/reset call.
 * NB, data is arranged such that subsequent bytes in the
 * bytestream are stored in ascending addresses.
 */
//func (this *TComOutputBitstream) GetByteStream() *byte {
//    return nil;//(Char*) &this.m_fifo.Front();
//}

/**
 * Return the number of valid bytes available from  getByteStream()
 */
func (this *TComOutputBitstream) GetByteStreamLength() uint {
    return uint(this.m_fifo.Len())
}

/**
 * Reset all internal state.
 */
func (this *TComOutputBitstream) Clear() {
    this.m_fifo.Init()
    this.m_held_bits = 0
    this.m_num_held_bits = 0
}

/**
 * returns the number of bits that need to be written to
 * achieve byte alignment.
 */
func (this *TComOutputBitstream) GetNumBitsUntilByteAligned() int {
    return int((8 - this.m_num_held_bits) & 0x7)
}

/**
 * Return the number of bits that have been written since the last clear()
 */
func (this *TComOutputBitstream) GetNumberOfWrittenBits() uint {
    //return uint(m_fifo->size()) * 8 + this.m_num_held_bits;
    return uint(this.m_fifo.Len())*8 + this.m_num_held_bits
}

func (this *TComOutputBitstream) InsertAt(src *TComOutputBitstream, pos uint) {
    //src_bits := src.GetNumberOfWrittenBits();
    //assert(0 == src_bits % 8);

    i := uint(0)
    for e := this.m_fifo.Front(); e != nil; e = e.Next() {
        if i == pos {
            for f := src.m_fifo.Front(); f != nil; f = f.Next() {
                v := f.Value.(byte)
                this.m_fifo.InsertBefore(v, e)
            }
            break
        }

        i++
    }
    //vector<uint8_t>::iterator at = this->m_fifo->begin() + pos;
    //this->m_fifo->insert(at, src.m_fifo->begin(), src.m_fifo->end());
}

/**
 * Return a reference to the internal fifo
 */
func (this *TComOutputBitstream) GetFIFO() *list.List {
    return this.m_fifo
}

func (this *TComOutputBitstream) GetHeldBits() byte {
    return this.m_held_bits
}

func (this *TComOutputBitstream) Copy(src *TComOutputBitstream) {
    //vector<uint8_t>::iterator at = this->m_fifo->begin();
    //this->m_fifo->insert(at, src.m_fifo->begin(), src.m_fifo->end());

    e := this.m_fifo.Front()
    for f := src.m_fifo.Front(); f != nil; f = f.Next() {
        v := f.Value.(byte)
        this.m_fifo.InsertBefore(v, e)
    }

    this.m_num_held_bits = src.m_num_held_bits
    this.m_held_bits = src.m_held_bits
}

func (this *TComOutputBitstream) AddSubstream(pcSubstream *TComOutputBitstream) {
    uiNumBits := pcSubstream.GetNumberOfWrittenBits()

    rbsp := pcSubstream.GetFIFO()
    for e := rbsp.Front(); e != nil; e = e.Next() {
        v := e.Value.(byte)
        this.Write(uint(v), 8)
    }
    if uiNumBits&0x7 != 0 {
        this.Write(uint(pcSubstream.GetHeldBits()>>(8-(uiNumBits&0x7))), uiNumBits&0x7)
    }
}

func (this *TComOutputBitstream) WriteByteAlignment() {
    this.Write(1, 1)
    this.WriteAlignZero()
}

/**
 * Write rbsp_trailing_bits to bs causing it to become byte-aligned
 */
func (this *TComOutputBitstream) WriteRBSPTrailingBits() {
    this.Write(1, 1)
    this.WriteAlignZero()
}

/**
 * Model of an input bitstream that extracts bits from a predefined
 * bytestream.
 */
type TComInputBitstream struct {
    m_fifo *list.List //std::vector<uint8_t> *m_fifo; /// FIFO for storage of complete bytes

    m_emulationPreventionByteLocation map[uint]uint; //  std::vector<UInt>

    //protected:
    m_fifo_idx uint /// Read index into m_fifo
    m_fifo_ptr *list.Element

    m_num_held_bits uint
    m_held_bits     byte
    m_numBitsRead   uint
}

//public:
/**
 * Create a new bitstream reader object that reads from #buf#.  Ownership
 * of #buf# remains with the callee, although the constructed object
 * will hold a reference to #buf#
 */

func NewTComInputBitstream(buf *list.List) *TComInputBitstream { // std::vector<uint8_t>* buf);
    if buf!=nil {
        return &TComInputBitstream{buf, make(map[uint]uint), 0, buf.Front(), 0, 0, 0}
    }

    return &TComInputBitstream{nil, make(map[uint]uint), 0, nil, 0, 0, 0}
}

// interface for decoding
func (this *TComInputBitstream) PseudoRead(uiNumberOfBits uint, ruiBits *uint) {
    saved_num_held_bits := this.m_num_held_bits
    saved_held_bits := this.m_held_bits
    saved_fifo_idx := this.m_fifo_idx
    saved_fifo_ptr := this.m_fifo_ptr

    var num_bits_to_read uint
    if uiNumberOfBits < this.GetNumBitsLeft() {
        num_bits_to_read = uiNumberOfBits
    } else {
        num_bits_to_read = this.GetNumBitsLeft()
    }
    this.Read(num_bits_to_read, ruiBits)
    *ruiBits <<= uint(uiNumberOfBits - num_bits_to_read)

	this.m_fifo_ptr = saved_fifo_ptr
    this.m_fifo_idx = saved_fifo_idx
    this.m_held_bits = saved_held_bits
    this.m_num_held_bits = saved_num_held_bits
}
func (this *TComInputBitstream) Read(uiNumberOfBits uint, ruiBits *uint) (err error) {
    //assert( uiNumberOfBits <= 32 );

    this.m_numBitsRead += uiNumberOfBits

    /* NB, bits are extracted from the MSB of each byte. */
    retval := uint(0)
    if uiNumberOfBits <= this.m_num_held_bits {
        /* n=1, len(H)=7:   -VHH HHHH, shift_down=6, mask=0xfe
         * n=3, len(H)=7:   -VVV HHHH, shift_down=4, mask=0xf8
         */
        retval = uint(this.m_held_bits) >> (uint(this.m_num_held_bits) - uiNumberOfBits)
        retval &= ^(0xff << uiNumberOfBits)
        this.m_num_held_bits -= uiNumberOfBits
        *ruiBits = retval
        return
    }

    /* all num_held_bits will go into retval
     *   => need to mask leftover bits from previous extractions
     *   => align retval with top of extracted word */
    /* n=5, len(H)=3: ---- -VVV, mask=0x07, shift_up=5-3=2,
     * n=9, len(H)=3: ---- -VVV, mask=0x07, shift_up=9-3=6 */
    uiNumberOfBits -= this.m_num_held_bits
    retval = uint(this.m_held_bits) & ^(0xff << uint(this.m_num_held_bits))
    retval <<= uiNumberOfBits

    /* number of whole bytes that need to be loaded to form retval */
    /* n=32, len(H)=0, load 4bytes, shift_down=0
     * n=32, len(H)=1, load 4bytes, shift_down=1
     * n=31, len(H)=1, load 4bytes, shift_down=1+1
     * n=8,  len(H)=0, load 1byte,  shift_down=0
     * n=8,  len(H)=3, load 1byte,  shift_down=3
     * n=5,  len(H)=1, load 1byte,  shift_down=1+3
     */
    aligned_word := uint(0)
    num_bytes_to_load := (uiNumberOfBits - 1) >> 3
    //assert(m_fifo_idx + num_bytes_to_load < m_fifo->size());
    if this.m_fifo_idx+num_bytes_to_load >= uint(this.m_fifo.Len()) {
        err = errors.New("this.m_fifo_idx >= this.m_fifo.Len()")
        return err
    }

    switch num_bytes_to_load {
    case 3:
        aligned_word = this.ReadByte() << 24
        fallthrough
    case 2:
        aligned_word |= this.ReadByte() << 16
        fallthrough
    case 1:
        aligned_word |= this.ReadByte() << 8
        fallthrough
    case 0:
        aligned_word |= this.ReadByte()
    }

    /* resolve remainder bits */
    next_num_held_bits := (32 - uiNumberOfBits) % 8

    /* copy required part of aligned_word into retval */
    retval |= aligned_word >> next_num_held_bits

    /* store held bits */
    this.m_num_held_bits = next_num_held_bits
    this.m_held_bits = byte(aligned_word)

    *ruiBits = uint(retval)

    return nil
}
func (this *TComInputBitstream) ReadByte1(ruiBits *uint) (err error) {
    //assert(m_fifo_idx < m_fifo->size());
    if this.m_fifo_idx >= uint(this.m_fifo.Len()) {
        err = errors.New("this.m_fifo_idx >= this.m_fifo.Len()")
        return err
    }

    /*idx := uint(0)
    elm := this.m_fifo.Front()
    for ; elm != nil; elm = elm.Next() {
        if this.m_fifo_idx == idx {
            break
        }

        idx++
    }*/

    *ruiBits = uint(this.m_fifo_ptr.Value.(byte));//uint(elm.Value.(byte)) //this.m_fifo.(*m_fifo)[m_fifo_idx++];
    this.m_fifo_ptr=this.m_fifo_ptr.Next()
    this.m_fifo_idx++

    return nil
}

func (this *TComInputBitstream) ReadOutTrailingBits() {
    uiBits := uint(0)

    for (this.GetNumBitsLeft() > 0) && (this.GetNumBitsUntilByteAligned() != 0) {
        this.Read(1, &uiBits)
    }
}
func (this *TComInputBitstream) GetHeldBits() byte {
    return this.m_held_bits
}
func (this *TComInputBitstream) Copy(src *TComOutputBitstream) {
}
func (this *TComInputBitstream) GetByteLocation() uint {
    return this.m_fifo_idx
}

// Peek at bits in word-storage. Used in determining if we have completed reading of current bitstream and therefore slice in LCEC.
func (this *TComInputBitstream) PeekBits(uiBits uint) uint {
    var tmp uint
    this.PseudoRead(uiBits, &tmp)
    return tmp
}

// utility functions
func (this *TComInputBitstream) ReadBits(numberOfBits uint) uint {
    var tmp uint
    this.Read(numberOfBits, &tmp)
    return tmp
}
func (this *TComInputBitstream) ReadByte() uint {
    var tmp uint
    this.ReadByte1(&tmp)
    return tmp
}
func (this *TComInputBitstream) GetNumBitsUntilByteAligned() uint {
    return this.m_num_held_bits & (0x7)
}

func (this *TComInputBitstream) GetNumBitsLeft() uint {
    return 8*(uint(this.m_fifo.Len())-this.m_fifo_idx) + this.m_num_held_bits
}

func (this *TComInputBitstream) ExtractSubstream(uiNumBits uint) *TComInputBitstream { // Read the nominated number of bits, and return as a bitstream.
    uiNumBytes := uiNumBits / 8

    buf := list.New() //std::vector<uint8_t>* buf = new std::vector<uint8_t>;
    var uiByte, ui uint
    for ui = 0; ui < uiNumBytes; ui++ {
        this.Read(8, &uiByte)
        buf.PushBack(byte(uiByte))
    }

    if uiNumBits&0x7 != 0 {
        uiByte = 0
        this.Read(uiNumBits&0x7, &uiByte)
        uiByte <<= 8 - (uiNumBits & 0x7)
        buf.PushBack(byte(uiByte))
    }
    return &TComInputBitstream{buf, make(map[uint]uint), 0, buf.Front(), 0, 0, 0}
}
func (this *TComInputBitstream) DeleteFifo() { // Delete internal fifo of bitstream.
    //delete m_fifo;
    //m_fifo = NULL;
}
func (this *TComInputBitstream) GetNumBitsRead() uint {
    return this.m_numBitsRead
}
func (this *TComInputBitstream) ReadByteAlignment() {
    var code uint
    this.Read(1, &code)
    //assert(code == 1);

    numBits := this.GetNumBitsUntilByteAligned()
    if numBits != 0 {
        //assert(numBits <= getNumBitsLeft());
        this.Read(numBits, &code)
        //assert(code == 0);
    }
}

func (this *TComInputBitstream) PushEmulationPreventionByteLocation  ( pos uint)           { this.m_emulationPreventionByteLocation[uint(len(this.m_emulationPreventionByteLocation))] = pos; }
func (this *TComInputBitstream) NumEmulationPreventionBytesRead      () uint               { return uint(len(this.m_emulationPreventionByteLocation));     }
func (this *TComInputBitstream) GetEmulationPreventionByteLocation   () map[uint]uint      { return this.m_emulationPreventionByteLocation;          }
func (this *TComInputBitstream) GetEmulationPreventionByteLocationIdx( idx uint)  uint     { return this.m_emulationPreventionByteLocation[ idx ];   }
func (this *TComInputBitstream) ClearEmulationPreventionByteLocation ()                    { this.m_emulationPreventionByteLocation = make(map[uint]uint); }
func (this *TComInputBitstream) SetEmulationPreventionByteLocation   ( vec map[uint]uint)  { this.m_emulationPreventionByteLocation = vec;           }


//! \}
