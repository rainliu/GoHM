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
}

/** insert one bits until the bitstream is byte-aligned */
func (this *TComOutputBitstream) WriteAlignOne() {
}

/** insert zero bits until the bitstream is byte-aligned */
func (this *TComOutputBitstream) WriteAlignZero() {
}

/** this function should never be called */
func (this *TComOutputBitstream) ResetBits() {
}

// utility functions

/**
 * Return a pointer to the start of the byte-stream buffer.
 * Pointer is valid until the next write/flush/reset call.
 * NB, data is arranged such that subsequent bytes in the
 * bytestream are stored in ascending addresses.
 */
func (this *TComOutputBitstream) GetByteStream() *byte {
    return nil
}

/**
 * Return the number of valid bytes available from  getByteStream()
 */
func (this *TComOutputBitstream) GetByteStreamLength() uint {
    return 0
}

/**
 * Reset all internal state.
 */
func (this *TComOutputBitstream) Clear() {
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
}

func (this *TComOutputBitstream) AddSubstream(pcSubstream *TComOutputBitstream) {
}

func (this *TComOutputBitstream) WriteByteAlignment() {
}

/**
 * Model of an input bitstream that extracts bits from a predefined
 * bytestream.
 */
type TComInputBitstream struct {
    m_fifo *list.List //std::vector<uint8_t> *m_fifo; /// FIFO for storage of complete bytes

    //protected:
    m_fifo_idx uint /// Read index into m_fifo

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
    return &TComInputBitstream{buf, 0, 0, 0, 0}
}

// interface for decoding
func (this *TComInputBitstream) PseudoRead(uiNumberOfBits uint, ruiBits *uint) {
    saved_num_held_bits := this.m_num_held_bits
    saved_held_bits := this.m_held_bits
    saved_fifo_idx := this.m_fifo_idx

    var num_bits_to_read uint
    if uiNumberOfBits < this.GetNumBitsLeft() {
        num_bits_to_read = uiNumberOfBits
    } else {
        num_bits_to_read = this.GetNumBitsLeft()
    }
    this.Read(num_bits_to_read, ruiBits)
    *ruiBits <<= uint(uiNumberOfBits - num_bits_to_read)

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

    idx := uint(0)
    elm := this.m_fifo.Front()
    for ; elm != nil; elm = elm.Next() {
        if this.m_fifo_idx == idx {
            break
        }

        idx++
    }

    *ruiBits = uint(elm.Value.(byte)) //this.m_fifo.(*m_fifo)[m_fifo_idx++];
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
        buf.PushBack(uiByte)
    }

    if uiNumBits&0x7 != 0 {
        uiByte = 0
        this.Read(uiNumBits&0x7, &uiByte)
        uiByte <<= 8 - (uiNumBits & 0x7)
        buf.PushBack(uiByte)
    }
    return &TComInputBitstream{buf, 0, 0, 0, 0}
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

//! \}
