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

import ()

/**
 * Abstract class representing an SEI message with lightweight RTTI.
 */
/*
class SEI
{
public:
  enum PayloadType
  {
    BUFFERING_PERIOD       = 0,
    PICTURE_TIMING         = 1,
    USER_DATA_UNREGISTERED = 5,
    RECOVERY_POINT         = 6,
#if SEI_DISPLAY_ORIENTATION
    DISPLAY_ORIENTATION    = 47,
#endif
    ACTIVE_PARAMETER_SETS  = 130,
#if SEI_TEMPORAL_LEVEL0_INDEX
    TEMPORAL_LEVEL0_INDEX  = 132,
#endif
#if SUFFIX_SEI_NUT_DECODED_HASH_SEI
    DECODED_PICTURE_HASH   = 133,
#else
    DECODED_PICTURE_HASH   = 256,
#endif
  };

  SEI() {}
  virtual ~SEI() {}

  virtual PayloadType payloadType() const = 0;
};

class SEIuserDataUnregistered : public SEI
{
public:
  PayloadType payloadType() const { return USER_DATA_UNREGISTERED; }

  SEIuserDataUnregistered()
    : userData(0)
    {}

  virtual ~SEIuserDataUnregistered()
  {
    delete userData;
  }

  UChar uuid_iso_iec_11578[16];
  UInt userDataLength;
  UChar *userData;
};

class SEIDecodedPictureHash : public SEI
{
public:
  PayloadType payloadType() const { return DECODED_PICTURE_HASH; }

  SEIDecodedPictureHash() {}
  virtual ~SEIDecodedPictureHash() {}

  enum Method
  {
    MD5,
    CRC,
    CHECKSUM,
    RESERVED,
  } method;

  UChar digest[3][16];
};

class SEIActiveParameterSets : public SEI
{
public:
  PayloadType payloadType() const { return ACTIVE_PARAMETER_SETS; }

  SEIActiveParameterSets()
    :activeSPSIdPresentFlag(1)
#if !HLS_REMOVE_ACTIVE_PARAM_SET_SEI_EXT_FLAG
    ,activeParamSetSEIExtensionFlag(0)
#endif // HLS_REMOVE_ACTIVE_PARAM_SET_SEI_EXT_FLAG
  {}
  virtual ~SEIActiveParameterSets() {}

  Int activeVPSId;
  Int activeSPSIdPresentFlag;
  Int activeSeqParamSetId;
#if !HLS_REMOVE_ACTIVE_PARAM_SET_SEI_EXT_FLAG
  Int activeParamSetSEIExtensionFlag;
#endif // HLS_REMOVE_ACTIVE_PARAM_SET_SEI_EXT_FLAG
};

class SEIBufferingPeriod : public SEI
{
public:
  PayloadType payloadType() const { return BUFFERING_PERIOD; }

  SEIBufferingPeriod()
  :m_sps (NULL)
  {}
  virtual ~SEIBufferingPeriod() {}

  UInt m_seqParameterSetId;
  Bool m_altCpbParamsPresentFlag;
  UInt m_initialCpbRemovalDelay         [MAX_CPB_CNT][2];
  UInt m_initialCpbRemovalDelayOffset   [MAX_CPB_CNT][2];
  UInt m_initialAltCpbRemovalDelay      [MAX_CPB_CNT][2];
  UInt m_initialAltCpbRemovalDelayOffset[MAX_CPB_CNT][2];
  TComSPS* m_sps;
};
class SEIPictureTiming : public SEI
{
public:
  PayloadType payloadType() const { return PICTURE_TIMING; }

  SEIPictureTiming()
  : m_numNalusInDuMinus1      (NULL)
  , m_duCpbRemovalDelayMinus1 (NULL)
  , m_sps                     (NULL)
  {}
  virtual ~SEIPictureTiming()
  {
    if( m_numNalusInDuMinus1 != NULL )
    {
      delete m_numNalusInDuMinus1;
    }
    if( m_duCpbRemovalDelayMinus1  != NULL )
    {
      delete m_duCpbRemovalDelayMinus1;
    }
  }

  UInt  m_auCpbRemovalDelay;
  UInt  m_picDpbOutputDelay;
  UInt  m_numDecodingUnitsMinus1;
  Bool  m_duCommonCpbRemovalDelayFlag;
  UInt  m_duCommonCpbRemovalDelayMinus1;
  UInt* m_numNalusInDuMinus1;
  UInt* m_duCpbRemovalDelayMinus1;
  TComSPS* m_sps;
};
class SEIRecoveryPoint : public SEI
{
public:
  PayloadType payloadType() const { return RECOVERY_POINT; }

  SEIRecoveryPoint() {}
  virtual ~SEIRecoveryPoint() {}

  Int  m_recoveryPocCnt;
  Bool m_exactMatchingFlag;
  Bool m_brokenLinkFlag;
};
#if SEI_DISPLAY_ORIENTATION
class SEIDisplayOrientation : public SEI
{
public:
  PayloadType payloadType() const { return DISPLAY_ORIENTATION; }

  SEIDisplayOrientation()
    : cancelFlag(true)
    , repetitionPeriod(1)
    , extensionFlag(false)
    {}
  virtual ~SEIDisplayOrientation() {}

  Bool cancelFlag;
  Bool horFlip;
  Bool verFlip;

  UInt anticlockwiseRotation;
  UInt repetitionPeriod;
  Bool extensionFlag;
};
#endif
#if SEI_TEMPORAL_LEVEL0_INDEX
class SEITemporalLevel0Index : public SEI
{
public:
  PayloadType payloadType() const { return TEMPORAL_LEVEL0_INDEX; }

  SEITemporalLevel0Index()
    : tl0Idx(0)
    , rapIdx(0)
    {}
  virtual ~SEITemporalLevel0Index() {}

  UInt tl0Idx;
  UInt rapIdx;
};
#endif
*/
/**
 * A structure to collate all SEI messages.  This ought to be replaced
 * with a list of std::list<SEI*>.  However, since there is only one
 * user of the SEI framework, this will do initially */
type SEImessages struct {
    /*SEIuserDataUnregistered* user_data_unregistered;
      SEIActiveParameterSets* active_parameter_sets;
      SEIDecodedPictureHash* picture_digest;
      SEIBufferingPeriod* buffering_period;
      SEIPictureTiming* picture_timing;*/
    m_pSPS *TComSPS /*
         SEIRecoveryPoint* recovery_point;
       #if SEI_DISPLAY_ORIENTATION
         SEIDisplayOrientation* display_orientation;
       #endif
       #if SEI_TEMPORAL_LEVEL0_INDEX
         SEITemporalLevel0Index* temporal_level0_index;
       #endif
    */
}

func NewSEImessages() *SEImessages {
    return &SEImessages{}
}

func (this *SEImessages) GetSPS() *TComSPS {
    return this.m_pSPS
}

func (this *SEImessages) SetSPS(sps *TComSPS) {
    this.m_pSPS = sps
}
