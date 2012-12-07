package TLibDecoder

import (
    "gohm/TLibCommon"
)

// ====================================================================================================================
// Class definition
// ====================================================================================================================

/*
#if ENC_DEC_TRACE

#define READ_CODE(length, code, name)     xReadCodeTr ( length, code, name )
#define READ_UVLC(        code, name)     xReadUvlcTr (         code, name )
#define READ_SVLC(        code, name)     xReadSvlcTr (         code, name )
#define READ_FLAG(        code, name)     xReadFlagTr (         code, name )

#else

#define READ_CODE(length, code, name)     xReadCode ( length, code )
#define READ_UVLC(        code, name)     xReadUvlc (         code )
#define READ_SVLC(        code, name)     xReadSvlc (         code )
#define READ_FLAG(        code, name)     xReadFlag (         code )

#endif
*/
//! \ingroup TLibDecoder
//! \{

// ====================================================================================================================
// Class definition
// ====================================================================================================================

type SyntaxElementParser struct {
    //protected:
    m_pcBitstream *TLibCommon.TComInputBitstream
}

/*
  SyntaxElementParser()
  : m_pcBitstream (NULL)
  {};
  virtual ~SyntaxElementParser() {};

  Void  xReadCode    ( UInt   length, UInt& val );
  Void  xReadUvlc    ( UInt&  val );
  Void  xReadSvlc    ( Int&   val );
  Void  xReadFlag    ( UInt&  val );
#if ENC_DEC_TRACE
  Void  xReadCodeTr  (UInt  length, UInt& rValue, const Char *pSymbolName);
  Void  xReadUvlcTr  (              UInt& rValue, const Char *pSymbolName);
  Void  xReadSvlcTr  (               Int& rValue, const Char *pSymbolName);
  Void  xReadFlagTr  (              UInt& rValue, const Char *pSymbolName);
#endif
public:
  Void  setBitstream ( TComInputBitstream* p )   { m_pcBitstream = p; }
  TComInputBitstream* getBitstream() { return m_pcBitstream; }
};*/

//class SEImessages;

/// CAVLC decoder class
type TDecCavlc struct {
    SyntaxElementParser //, public TDecEntropyIf
}

/*
public:
  TDecCavlc();
  virtual ~TDecCavlc();

protected:
  Void  xReadEpExGolomb       ( UInt& ruiSymbol, UInt uiCount );
  Void  xReadExGolombLevel    ( UInt& ruiSymbol );
  Void  xReadUnaryMaxSymbol   ( UInt& ruiSymbol, UInt uiMaxSymbol );

  Void  xReadPCMAlignZero     ();

  UInt  xGetBit             ();

  void  parseShortTermRefPicSet            (TComSPS* pcSPS, TComReferencePictureSet* pcRPS, Int idx);
private:

public:

  /// rest entropy coder by intial QP and IDC in CABAC
  Void  resetEntropy        ( TComSlice* pcSlice  )     { assert(0); };
  Void  setBitstream        ( TComInputBitstream* p )   { m_pcBitstream = p; }
  Void  parseTransformSubdivFlag( UInt& ruiSubdivFlag, UInt uiLog2TransformBlockSize );
  Void  parseQtCbf          ( TComDataCU* pcCU, UInt uiAbsPartIdx, TextType eType, UInt uiTrDepth, UInt uiDepth );
  Void  parseQtRootCbf      ( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiDepth, UInt& uiQtRootCbf );
  Void  parseVPS            ( TComVPS* pcVPS );
  Void  parseSPS            ( TComSPS* pcSPS );
  Void  parsePPS            ( TComPPS* pcPPS);
  Void  parseVUI            ( TComVUI* pcVUI, TComSPS* pcSPS );
  Void  parseSEI(SEImessages&);
  Void  parsePTL            ( TComPTL *rpcPTL, Bool profilePresentFlag, Int maxNumSubLayersMinus1 );
  Void  parseProfileTier    (ProfileTierLevel *ptl);
#if SIGNAL_BITRATE_PICRATE_IN_VPS
  Void  parseBitratePicRateInfo(TComBitRatePicRateInfo *info, Int tempLevelLow, Int tempLevelHigh);
#endif
  Void  parseSliceHeader    ( TComSlice*& rpcSlice, ParameterSetManagerDecoder *parameterSetManager);
  Void  parseTerminatingBit ( UInt& ruiBit );

  Void  parseMVPIdx         ( Int& riMVPIdx );

  Void  parseSkipFlag       ( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiDepth );
  Void  parseCUTransquantBypassFlag( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiDepth );
  Void parseMergeFlag       ( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiDepth, UInt uiPUIdx );
  Void parseMergeIndex      ( TComDataCU* pcCU, UInt& ruiMergeIndex, UInt uiAbsPartIdx, UInt uiDepth );
  Void parseSplitFlag       ( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiDepth );
  Void parsePartSize        ( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiDepth );
  Void parsePredMode        ( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiDepth );

  Void parseIntraDirLumaAng ( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiDepth );

  Void parseIntraDirChroma  ( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiDepth );

  Void parseInterDir        ( TComDataCU* pcCU, UInt& ruiInterDir, UInt uiAbsPartIdx, UInt uiDepth );
  Void parseRefFrmIdx       ( TComDataCU* pcCU, Int& riRefFrmIdx,  UInt uiAbsPartIdx, UInt uiDepth, RefPicList eRefList );
  Void parseMvd             ( TComDataCU* pcCU, UInt uiAbsPartAddr,UInt uiPartIdx,    UInt uiDepth, RefPicList eRefList );

  Void parseDeltaQP         ( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiDepth );
  Void parseCoeffNxN        ( TComDataCU* pcCU, TCoeff* pcCoef, UInt uiAbsPartIdx, UInt uiWidth, UInt uiHeight, UInt uiDepth, TextType eTType );
  Void parseTransformSkipFlags ( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt width, UInt height, UInt uiDepth, TextType eTType);

  Void parseIPCMInfo        ( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiDepth);

  Void updateContextTables  ( SliceType eSliceType, Int iQp ) { return; }

  Void xParsePredWeightTable ( TComSlice* pcSlice );
  Void  parseScalingList               ( TComScalingList* scalingList );
  Void xDecodeScalingList    ( TComScalingList *scalingList, UInt sizeId, UInt listId);
protected:
  Bool  xMoreRbspData();
};*/
