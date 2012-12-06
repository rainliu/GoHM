package TLibDecoder

import (
	"gohm/TLibCommon"
)

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

type SyntaxElementParser struct{
//protected:
  m_pcBitstream 	*TLibCommon.TComInputBitstream;
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