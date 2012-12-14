package TLibCommon

import (

)


// ====================================================================================================================
// Type definition
// ====================================================================================================================

/// parameters for AMVP
type AMVPInfo struct{
  m_acMvCand	[ AMVP_MAX_NUM_CANDS_MEM ]TComMv;  ///< array of motion vector predictor candidates
  iN	int;                                ///< number of motion vector predictor candidates
} ;

// ====================================================================================================================
// Class definition
// ====================================================================================================================

/// class for motion vector with reference index
type TComMvField struct{
//private:
  m_acMv		TComMv;
  m_iRefIdx		int;
}
/*
public:
  TComMvField() : m_iRefIdx( NOT_VALID ) {}
  
  Void setMvField( TComMv const & cMv, Int iRefIdx )
  {
    m_acMv    = cMv;
    m_iRefIdx = iRefIdx;
  }
  
  Void setRefIdx( Int refIdx ) { m_iRefIdx = refIdx; }
  
  TComMv const & getMv() const { return  m_acMv; }
  TComMv       & getMv()       { return  m_acMv; }
  
  Int getRefIdx() const { return  m_iRefIdx;       }
  Int getHor   () const { return  m_acMv.getHor(); }
  Int getVer   () const { return  m_acMv.getVer(); }
};*/

/// class for motion information in one CU
type TComCUMvField struct{
//private:
  m_pcMv		*TComMv;
  m_pcMvd		*TComMv;
  m_piRefIdx	*int8;
  m_uiNumPartition	uint;
  m_cAMVPInfo	AMVPInfo;
}
/*    
  template <typename T>
  Void setAll( T *p, T const & val, PartSize eCUMode, Int iPartAddr, UInt uiDepth, Int iPartIdx );

public:
  TComCUMvField() : m_pcMv(NULL), m_pcMvd(NULL), m_piRefIdx(NULL), m_uiNumPartition(0) {}
  ~TComCUMvField() {}
*/
  // ------------------------------------------------------------------------------------------------------------------
  // create / destroy
  // ------------------------------------------------------------------------------------------------------------------
  
func (this *TComCUMvField)  Create(  uiNumPartition uint){
}
func (this *TComCUMvField)  Destroy(){
}
  
  // ------------------------------------------------------------------------------------------------------------------
  // clear / copy
  // ------------------------------------------------------------------------------------------------------------------
/*
  Void    clearMvField();
  
  Void    copyFrom( TComCUMvField const * pcCUMvFieldSrc, Int iNumPartSrc, Int iPartAddrDst );
  Void    copyTo  ( TComCUMvField* pcCUMvFieldDst, Int iPartAddrDst ) const;
  Void    copyTo  ( TComCUMvField* pcCUMvFieldDst, Int iPartAddrDst, UInt uiOffset, UInt uiNumPart ) const;
  
  // ------------------------------------------------------------------------------------------------------------------
  // get
  // ------------------------------------------------------------------------------------------------------------------
  
  TComMv const & getMv    ( Int iIdx ) const { return  m_pcMv    [iIdx]; }
  TComMv const & getMvd   ( Int iIdx ) const { return  m_pcMvd   [iIdx]; }
  Int            getRefIdx( Int iIdx ) const { return  m_piRefIdx[iIdx]; }
  
  AMVPInfo* getAMVPInfo () { return &m_cAMVPInfo; }
  
  // ------------------------------------------------------------------------------------------------------------------
  // set
  // ------------------------------------------------------------------------------------------------------------------
  
  Void    setAllMv     ( TComMv const & rcMv,         PartSize eCUMode, Int iPartAddr, UInt uiDepth, Int iPartIdx=0 );
  Void    setAllMvd    ( TComMv const & rcMvd,        PartSize eCUMode, Int iPartAddr, UInt uiDepth, Int iPartIdx=0 );
  Void    setAllRefIdx ( Int iRefIdx,                 PartSize eMbMode, Int iPartAddr, UInt uiDepth, Int iPartIdx=0 );
  Void    setAllMvField( TComMvField const & mvField, PartSize eMbMode, Int iPartAddr, UInt uiDepth, Int iPartIdx=0 );
*/
func (this *TComCUMvField) SetNumPartition( uiNumPart uint){
    this.m_uiNumPartition = uiNumPart;
}
/*  
func (this *TComCUMvField) linkToWithOffset( src *TComCUMvField,  offset int){
    this.m_pcMv     = src->m_pcMv     + offset;
    this.m_pcMvd    = src->m_pcMvd    + offset;
    this.m_piRefIdx = src->m_piRefIdx + offset;
}
  
func (this *TComCUMvField) compress(Char* pePredMode, Int scale); 
*/
