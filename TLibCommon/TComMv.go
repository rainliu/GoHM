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

func NewTComMv(iHor, iVer int16)*TComMv {
    return &TComMv{iHor, iVer};
}
  // ------------------------------------------------------------------------------------------------------------------
  // constructors
  // ------------------------------------------------------------------------------------------------------------------

  // ------------------------------------------------------------------------------------------------------------------
  // set
  // ------------------------------------------------------------------------------------------------------------------

func (this *TComMv)  Set       (  iHor,  iVer int16)     {
	this.m_iHor = iHor;
	this.m_iVer = iVer;
}
func (this *TComMv)   SetHor    (  i int16)                   {
	this.m_iHor = i;
}
func (this *TComMv)   SetVer    (  i int16)                   {
	this.m_iVer = i;
}
func (this *TComMv)   SetZero   ()                            {
	this.m_iHor = 0;
	this.m_iVer = 0;
}

  // ------------------------------------------------------------------------------------------------------------------
  // get
  // ------------------------------------------------------------------------------------------------------------------

func (this *TComMv)     GetHor    () int16 {
	return this.m_iHor;
}
func (this *TComMv)     GetVer    () int16 {
	return this.m_iVer;
}
func (this *TComMv)     GetAbsHor () int16 {
	if this.m_iHor<0 {
		return -this.m_iHor;
	}

	return this.m_iHor;
}
func (this *TComMv)     GetAbsVer () int16 {
	if this.m_iVer<0 {
		return -this.m_iVer;
	}

	return this.m_iVer;
}

  // ------------------------------------------------------------------------------------------------------------------
  // operations
  // ------------------------------------------------------------------------------------------------------------------
 /*
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
  */
func (this *TComMv) ScaleMv(  iScale int) *TComMv{
    mvx := (iScale * int(this.GetHor()) + 127 + int(B2U(iScale * int(this.GetHor()) < 0))) >> 8;

    if mvx < -32768{
    	mvx = -32768;
    }else if mvx > 32767{
    	mvx = 32767;
    }

    mvy := (iScale * int(this.GetVer()) + 127 + int(B2U(iScale * int(this.GetVer()) < 0))) >> 8;

    if mvy < -32768{
    	mvy = -32768;
    }else if mvy > 32767{
    	mvy = 32767;
    }

    return &TComMv{ m_iHor:int16(mvx), m_iVer:int16(mvy) };
}
