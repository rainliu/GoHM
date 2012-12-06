package TLibDecoder

import (
	"gohm/TLibCommon"
)


type TDecBinCabac struct{//: public TDecBinIf
//private:
  m_pcTComBitstream		*TLibCommon.TComInputBitstream;
  m_uiRange		uint;
  m_uiValue		uint;
  m_bitsNeeded	int;
}
/*
public:
  TDecBinCABAC ();
  virtual ~TDecBinCABAC();
  
  Void  init              ( TComInputBitstream* pcTComBitstream );
  Void  uninit            ();
  
  Void  start             ();
  Void  finish            ();
  Void  flush             ();
  
  Void  decodeBin         ( UInt& ruiBin, ContextModel& rcCtxModel );
  Void  decodeBinEP       ( UInt& ruiBin                           );
  Void  decodeBinsEP      ( UInt& ruiBin, Int numBins              );
  Void  decodeBinTrm      ( UInt& ruiBin                           );
  
  Void  resetBac          ();
#if !REMOVE_BURST_IPCM
  Void  decodeNumSubseqIPCM( Int& numSubseqIPCM ) ;
#endif
  Void  decodePCMAlignBits();
  Void  xReadPCMCode      ( UInt uiLength, UInt& ruiCode );
  
  Void  copyState         ( TDecBinIf* pcTDecBinIf );
  TDecBinCABAC* getTDecBinCABAC()  { return this; }


};*/