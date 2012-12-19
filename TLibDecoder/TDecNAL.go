package TLibDecoder

import (
	"fmt"
	"errors"
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
  /* perform anti-emulation prevention */
  pcBitstream := TLibCommon.NewTComInputBitstream(nil);
//#if HM9_NALU_TYPES
  firstByte := nalUnitBuf.Front().Value.(byte)
  this.convertPayloadToRBSP(nalUnitBuf, pcBitstream, (firstByte & 64) == 0);
//#else
//  convertPayloadToRBSP(nalUnitBuf, pcBitstream);
//#endif
  
  this.m_Bitstream = TLibCommon.NewTComInputBitstream(nalUnitBuf);

  this.readNalUnitHeader();
}

func (this *InputNALUnit) GetBitstream() *TLibCommon.TComInputBitstream {
    return this.m_Bitstream
}

func (this *InputNALUnit) SetBitstream(bitstream *TLibCommon.TComInputBitstream) {
    this.m_Bitstream = bitstream
}

//#if HM9_NALU_TYPES
func (this *InputNALUnit) convertPayloadToRBSP(nalUnitBuf *list.List, pcBitstream *TLibCommon.TComInputBitstream, isVclNalUnit bool){
//#else
//static void convertPayloadToRBSP(vector<uint8_t>& nalUnitBuf, TComInputBitstream *pcBitstream)
//#endif
  zeroCount := 0;
  it_write := list.New();
  for e := nalUnitBuf.Front(); e != nil; e = e.Next() {
  	//assert(zeroCount < 2 || *it_read >= 0x03);
	it_read := e.Value.(byte)
	if zeroCount == 2 && it_read == 0x03 {
	  zeroCount = 0;	
	  
	  e = e.Next();
      if e == nil{
        break;
      }else{
      	it_read = e.Value.(byte)
      }
    }
    
    if it_read == 0x00 {
    	zeroCount++;
    }else{
    	zeroCount = 0;
    }
    it_write.PushBack(it_read);
  }
 
  //assert(zeroCount == 0);
  
//#if HM9_NALU_TYPES
  if isVclNalUnit{
    // Remove cabac_zero_word from payload if present
    n := 0;
    
    e := it_write.Back();
    it_read := e.Value.(byte)
    for it_read == 0x00 {
      it_write.Remove(e);
      e = e.Prev();
      it_read = e.Value.(byte)
      n++;
    }
    
    if n > 0 {
      fmt.Printf("\nDetected %d instances of cabac_zero_word", n/2);      
    }
  }
//#endif

  nalUnitBuf.Init();// = .resize(it_write - nalUnitBuf.begin());
  for e := it_write.Front(); e != nil; e = e.Next() {
  	it_read := e.Value.(byte)
  	nalUnitBuf.PushBack(it_read)
  }
}

func (this *InputNALUnit) readNalUnitHeader() error{
  bs := this.m_Bitstream;

  forbidden_zero_bit := bs.ReadBits(1);           // forbidden_zero_bit
  if forbidden_zero_bit != 0{
  	return errors.New("forbidden_zero_bit!=0");
  }
  
  this.SetNalUnitType (TLibCommon.NalUnitType(bs.ReadBits(6)));  // nal_unit_type
  
  this.SetReservedZero6Bits (bs.ReadBits(6));       // nuh_reserved_zero_6bits
  if this.GetReservedZero6Bits() != 0{
  	return errors.New("m_reservedZero6Bits!=0");
  }
  
  this.SetTemporalId (bs.ReadBits(3) - 1);             // nuh_temporal_id_plus1

  if this.GetTemporalId()!=0 {
    if !( this.GetNalUnitType() != TLibCommon.NAL_UNIT_CODED_SLICE_BLA  		&&
          this.GetNalUnitType() != TLibCommon.NAL_UNIT_CODED_SLICE_BLANT		&&
          this.GetNalUnitType() != TLibCommon.NAL_UNIT_CODED_SLICE_BLA_N_LP	&&
          this.GetNalUnitType() != TLibCommon.NAL_UNIT_CODED_SLICE_IDR			&&
          this.GetNalUnitType() != TLibCommon.NAL_UNIT_CODED_SLICE_IDR_N_LP	&&
          this.GetNalUnitType() != TLibCommon.NAL_UNIT_CODED_SLICE_CRA			&&
          this.GetNalUnitType() != TLibCommon.NAL_UNIT_VPS						&&
          this.GetNalUnitType() != TLibCommon.NAL_UNIT_SPS						&&
          this.GetNalUnitType() != TLibCommon.NAL_UNIT_EOS						&&
          this.GetNalUnitType() != TLibCommon.NAL_UNIT_EOB ){
         return errors.New("Wrong this.GetNalUnitType() in readNalUnitHeader")
    }
  }else{
    if !( this.GetNalUnitType() != TLibCommon.NAL_UNIT_CODED_SLICE_TLA 		&& 
          this.GetNalUnitType() != TLibCommon.NAL_UNIT_CODED_SLICE_TSA_N		&& 
          this.GetNalUnitType() != TLibCommon.NAL_UNIT_CODED_SLICE_STSA_R	&& 
          this.GetNalUnitType() != TLibCommon.NAL_UNIT_CODED_SLICE_STSA_N ){
         return errors.New("Wrong this.GetNalUnitType() in readNalUnitHeader")
    }
  }
  
  return nil
}