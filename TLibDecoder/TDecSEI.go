package TLibDecoder

import (
	"gohm/TLibCommon"
)

type TDecSeiReader struct {
    SyntaxElementParser
}

func NewTDecSeiReader() *TDecSeiReader{
	return &TDecSeiReader{}
}

//#if SUFFIX_SEI_NUT_DECODED_HASH_SEI
func (this *TDecSeiReader)  ParseSEImessage(bs *TLibCommon.TComInputBitstream, seis *TLibCommon.SEImessages, nalUnitType TLibCommon.NalUnitType){
}
/*
//#else
//  Void parseSEImessage(TComInputBitstream* bs, SEImessages& seis);
//#endif
//protected:
//#if SUFFIX_SEI_NUT_DECODED_HASH_SEI
  Void xReadSEImessage                (SEImessages& seis, const NalUnitType nalUnitType);
//#else
//  Void xReadSEImessage                (SEImessages& seis);
//#endif
  Void xParseSEIuserDataUnregistered  (SEIuserDataUnregistered &sei, UInt payloadSize);
  Void xParseSEIActiveParameterSets   (SEIActiveParameterSets  &sei, UInt payloadSize);
  Void xParseSEIDecodedPictureHash    (SEIDecodedPictureHash& sei, UInt payloadSize);
  Void xParseSEIBufferingPeriod       (SEIBufferingPeriod& sei, UInt payloadSize);
  Void xParseSEIPictureTiming         (SEIPictureTiming& sei, UInt payloadSize);
  Void xParseSEIRecoveryPoint         (SEIRecoveryPoint& sei, UInt payloadSize);
//#if SEI_DISPLAY_ORIENTATION
  Void xParseSEIDisplayOrientation    (SEIDisplayOrientation &sei, UInt payloadSize);
//#endif
//#if SEI_TEMPORAL_LEVEL0_INDEX
  Void xParseSEITemporalLevel0Index   (SEITemporalLevel0Index &sei, UInt payloadSize);
//#endif
  Void xParseByteAlign();
//};
*/