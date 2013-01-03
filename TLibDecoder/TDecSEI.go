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
    "gohm/TLibCommon"
)

type TDecSeiReader struct {
    SyntaxElementParser
}

func NewTDecSeiReader() *TDecSeiReader {
    return &TDecSeiReader{}
}

//#if SUFFIX_SEI_NUT_DECODED_HASH_SEI
func (this *TDecSeiReader) ParseSEImessage(bs *TLibCommon.TComInputBitstream, seis *TLibCommon.SEImessages, nalUnitType TLibCommon.NalUnitType) {
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
