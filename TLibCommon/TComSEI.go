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
      SEIPictureTiming* picture_timing;
      TComSPS* m_pSPS;
      SEIRecoveryPoint* recovery_point;
    #if SEI_DISPLAY_ORIENTATION
      SEIDisplayOrientation* display_orientation;
    #endif
    #if SEI_TEMPORAL_LEVEL0_INDEX
      SEITemporalLevel0Index* temporal_level0_index;
    #endif
    */
}

/*
public:
  SEImessages()
    : user_data_unregistered(0)
    , active_parameter_sets(0)
    , picture_digest(0)
    , buffering_period(0)
    , picture_timing(0)
    , recovery_point(0)
#if SEI_DISPLAY_ORIENTATION
    , display_orientation(0)
#endif
#if SEI_TEMPORAL_LEVEL0_INDEX
    , temporal_level0_index(0)
#endif
    {}

  ~SEImessages()
  {
    delete user_data_unregistered;
    delete active_parameter_sets; 
    delete picture_digest;
    delete buffering_period;
    delete picture_timing;
    delete recovery_point;
#if SEI_DISPLAY_ORIENTATION
    delete display_orientation;
#endif
#if SEI_TEMPORAL_LEVEL0_INDEX
    delete temporal_level0_index;
#endif
  }
*/
