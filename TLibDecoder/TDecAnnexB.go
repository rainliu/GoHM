package TLibDecoder

import (
	"io"
	"os"
	"fmt"
	"container/list"
)


type InputByteStream struct{
//private:
  m_NumFutureBytes		uint; /* number of valid bytes in m_FutureBytes */
  m_FutureBytes	uint32; /* bytes that have been peeked */
  m_Input		*os.File; /* Input stream to read from */
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
func NewInputByteStream(istream *os.File) *InputByteStream{
	return &InputByteStream{m_NumFutureBytes:0, m_FutureBytes:0, m_Input:istream}
}

  /**
   * Reset the internal state.  Must be called if input stream is
   * modified externally to this class
   */
func (this *InputByteStream) Reset(){
    this.m_NumFutureBytes = 0;
    this.m_FutureBytes = 0;
}

  /**
   * returns true if an EOF will be encountered within the next
   * n bytes.
   */
func (this *InputByteStream) EofBeforeNBytes(n uint) bool{
    if n > 4 {
    	fmt.Printf("n must be smaller or equal to 4\n")
    	return false
    }
    
    if this.m_NumFutureBytes >= n{
      return false;
    }

    n -= this.m_NumFutureBytes;
	
	buf := make([]byte, 1)
	
    for i := uint(0); i < n; i++ {
    	_, err := this.m_Input.Read(buf)
    	if err == io.EOF{
    		return true
    	}
        this.m_FutureBytes = (this.m_FutureBytes << 8) | uint32(buf[0]);
        this.m_NumFutureBytes++;
    }
    
    return false;
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
func (this *InputByteStream) PeekBytes(n uint) uint32{
    this.EofBeforeNBytes(n);
    return this.m_FutureBytes >> uint32(8*(this.m_NumFutureBytes - n));
  }

  /**
   * consume and return one byte from the input.
   *
   * If bytestream is already at EOF prior to a call to readByte(),
   * an exception std::ios_base::failure is thrown.
   */
func (this *InputByteStream) ReadByte() byte {
	wanted_byte := make([]byte, 1)
	
    if this.m_NumFutureBytes==0 {
    	this.m_Input.Read(wanted_byte)
      	return wanted_byte[0];
    }
    this.m_NumFutureBytes--;
    wanted_byte[0] = byte(this.m_FutureBytes >> uint32(8*this.m_NumFutureBytes));
    this.m_FutureBytes &= ^(0xff << uint32(8*this.m_NumFutureBytes));
    return wanted_byte[0];
  }

  /**
   * consume and return n bytes from the input.  n bytes from
   * bytestream are interpreted as bigendian when assembling
   * the return value.
   */
func (this *InputByteStream) ReadBytes(n uint) uint32 {
    var val uint32 = 0;
    for i := uint(0); i < n; i++ {
      val = (val << 8) | uint32(this.ReadByte());
    }
    
    return val;
  }

func (this *InputByteStream) byteStreamNALUnit(nalUnit *list.List, stats *AnnexBStats ) (err error) {
	return nil
}

func (this *InputByteStream) ByteStreamNALUnit(nalUnit *list.List, stats *AnnexBStats) bool {
  eof := false;
 
  if err:= this.byteStreamNALUnit(nalUnit, stats); err!=nil {
    eof = true;
  }

  stats.m_numBytesInNALUnit = uint(nalUnit.Len());
  return eof;	
}
/**
 * Statistics associated with AnnexB bytestreams
 */
type AnnexBStats struct{
  m_numLeadingZero8BitsBytes	uint;
  m_numZeroByteBytes			uint;
  m_numStartCodePrefixBytes	uint;
  m_numBytesInNALUnit			uint;
  m_numTrailingZero8BitsBytes	uint;
}

func (this *AnnexBStats)  Add(rhs *AnnexBStats) {
    this.m_numLeadingZero8BitsBytes += rhs.m_numLeadingZero8BitsBytes;
    this.m_numZeroByteBytes += rhs.m_numZeroByteBytes;
    this.m_numStartCodePrefixBytes += rhs.m_numStartCodePrefixBytes;
    this.m_numBytesInNALUnit += rhs.m_numBytesInNALUnit;
    this.m_numTrailingZero8BitsBytes += rhs.m_numTrailingZero8BitsBytes;
  }


