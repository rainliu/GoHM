package TLibCommon

import (

)


// ====================================================================================================================
// Class definition
// ====================================================================================================================

/// context model class
type ContextModel struct{
//private:
  m_ucState			byte;                                                                  ///< internal state variable
  m_aucNextStateMPS	[ 128 ]byte;
  m_aucNextStateLPS	[ 128 ]byte;
  m_entropyBits	[ 128 ]int;
//#if FAST_BIT_EST
  m_nextState	[128][2]byte;
//#endif
  m_binsCoded	uint;
}
/*
public:
  ContextModel  ()                        { m_ucState = 0; m_binsCoded = 0; }
  ~ContextModel ()                        {}
  
  UChar getState  ()                { return ( m_ucState >> 1 ); }                    ///< get current state
  UChar getMps    ()                { return ( m_ucState  & 1 ); }                    ///< get curret MPS
  Void  setStateAndMps( UChar ucState, UChar ucMPS) { m_ucState = (ucState << 1) + ucMPS; } ///< set state and MPS
  
  Void init ( Int qp, Int initValue );   ///< initialize state with initial probability
  
  Void updateLPS ()
  {
    m_ucState = m_aucNextStateLPS[ m_ucState ];
  }
  
  Void updateMPS ()
  {
    m_ucState = m_aucNextStateMPS[ m_ucState ];
  }
  
  Int getEntropyBits(Short val) { return m_entropyBits[m_ucState ^ val]; }
    
#if FAST_BIT_EST
  Void update( Int binVal )
  {
    m_ucState = m_nextState[m_ucState][binVal];
  }
  static Void buildNextStateTable();
  static Int getEntropyBitsTrm( Int val ) { return m_entropyBits[126 ^ val]; }
#endif
  Void setBinsCoded(UInt val)   { m_binsCoded = val;  }
  UInt getBinsCoded()           { return m_binsCoded;   }
  
};*/