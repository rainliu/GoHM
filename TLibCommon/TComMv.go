package TLibCommon

import (

)


// ====================================================================================================================
// Class definition
// ====================================================================================================================

/// basic motion vector class
type TComMv struct{
//private:
  m_iHor	int16;     ///< horizontal component of motion vector
  m_iVer	int16;     ///< vertical component of motion vector
}

/*
public:
  
  // ------------------------------------------------------------------------------------------------------------------
  // constructors
  // ------------------------------------------------------------------------------------------------------------------
  
  TComMv() :
  m_iHor(0),
  m_iVer(0)
  {
  }
  
  TComMv( Short iHor, Short iVer ) :
  m_iHor(iHor),
  m_iVer(iVer)
  {
  }
  
  // ------------------------------------------------------------------------------------------------------------------
  // set
  // ------------------------------------------------------------------------------------------------------------------
  
  Void  set       ( Short iHor, Short iVer)     { m_iHor = iHor;  m_iVer = iVer;            }
  Void  setHor    ( Short i )                   { m_iHor = i;                               }
  Void  setVer    ( Short i )                   { m_iVer = i;                               }
  Void  setZero   ()                            { m_iHor = m_iVer = 0;  }
  
  // ------------------------------------------------------------------------------------------------------------------
  // get
  // ------------------------------------------------------------------------------------------------------------------
  
  Int   getHor    () const { return m_iHor;          }
  Int   getVer    () const { return m_iVer;          }
  Int   getAbsHor () const { return abs( m_iHor );   }
  Int   getAbsVer () const { return abs( m_iVer );   }
  
  // ------------------------------------------------------------------------------------------------------------------
  // operations
  // ------------------------------------------------------------------------------------------------------------------
  
  const TComMv& operator += (const TComMv& rcMv)
  {
    m_iHor += rcMv.m_iHor;
    m_iVer += rcMv.m_iVer;
    return  *this;
  }
  
  const TComMv& operator-= (const TComMv& rcMv)
  {
    m_iHor -= rcMv.m_iHor;
    m_iVer -= rcMv.m_iVer;
    return  *this;
  }
  
  const TComMv& operator>>= (const Int i)
  {
    m_iHor >>= i;
    m_iVer >>= i;
    return  *this;
  }
  
  const TComMv& operator<<= (const Int i)
  {
    m_iHor <<= i;
    m_iVer <<= i;
    return  *this;
  }
  
  const TComMv operator - ( const TComMv& rcMv ) const
  {
    return TComMv( m_iHor - rcMv.m_iHor, m_iVer - rcMv.m_iVer );
  }
  
  const TComMv operator + ( const TComMv& rcMv ) const
  {
    return TComMv( m_iHor + rcMv.m_iHor, m_iVer + rcMv.m_iVer );
  }
  
  Bool operator== ( const TComMv& rcMv ) const
  {
    return (m_iHor==rcMv.m_iHor && m_iVer==rcMv.m_iVer);
  }
  
  Bool operator!= ( const TComMv& rcMv ) const
  {
    return (m_iHor!=rcMv.m_iHor || m_iVer!=rcMv.m_iVer);
  }
  
  const TComMv scaleMv( Int iScale ) const
  {
    Int mvx = Clip3( -32768, 32767, (iScale * getHor() + 127 + (iScale * getHor() < 0)) >> 8 );
    Int mvy = Clip3( -32768, 32767, (iScale * getVer() + 127 + (iScale * getVer() < 0)) >> 8 );
    return TComMv( mvx, mvy );
  }
};// END CLASS DEFINITION TComMV
*/
