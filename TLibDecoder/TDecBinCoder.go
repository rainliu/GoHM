package TLibDecoder

import (

)

//class TDecBinCABAC;

type TDecBinIf	interface{
/*
public:
  virtual Void  init              ( TComInputBitstream* pcTComBitstream )     = 0;
  virtual Void  uninit            ()                                          = 0;

  virtual Void  start             ()                                          = 0;
  virtual Void  finish            ()                                          = 0;
  virtual Void  flush            ()                                           = 0;

  virtual Void  decodeBin         ( UInt& ruiBin, ContextModel& rcCtxModel )  = 0;
  virtual Void  decodeBinEP       ( UInt& ruiBin                           )  = 0;
  virtual Void  decodeBinsEP      ( UInt& ruiBins, Int numBins             )  = 0;
  virtual Void  decodeBinTrm      ( UInt& ruiBin                           )  = 0;
  
  virtual Void  resetBac          ()                                          = 0;
#if !REMOVE_BURST_IPCM
  virtual Void  decodeNumSubseqIPCM( Int& numSubseqIPCM )                  = 0;
#endif
  virtual Void  decodePCMAlignBits()                                          = 0;
  virtual Void  xReadPCMCode      ( UInt uiLength, UInt& ruiCode)              = 0;

  virtual ~TDecBinIf() {}

  virtual Void  copyState         ( TDecBinIf* pcTDecBinIf )                  = 0;
  virtual TDecBinCABAC*   getTDecBinCABAC   ()  { return 0; }
 */
};