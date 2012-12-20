package TLibCommon

import (

)


// ====================================================================================================================
// Type definition
// ====================================================================================================================

/// parameters for AMVP
type AMVPInfo struct{
  MvCand	[ AMVP_MAX_NUM_CANDS_MEM ]TComMv;  ///< array of motion vector predictor candidates
  IN	int;                                ///< number of motion vector predictor candidates
} ;

// ====================================================================================================================
// Class definition
// ====================================================================================================================

/// class for motion vector with reference index
type TComMvField struct{
//private:
  m_acMv		TComMv;
  m_iRefIdx		int8;
}

func NewTComMvField() *TComMvField{ 
	return &TComMvField{m_iRefIdx: NOT_VALID}
}
  
func (this *TComMvField) SetMvField( cMv *TComMv,  iRefIdx int8){
    this.m_acMv.SetHor( cMv.GetHor() );
    this.m_acMv.SetVer( cMv.GetVer() );
    
    this.m_iRefIdx = iRefIdx;
}
  
func (this *TComMvField) SetRefIdx( refIdx int8) { 
	this.m_iRefIdx = refIdx; 
}
  
func (this *TComMvField)  GetMv() *TComMv { 
	return  &this.m_acMv; 
}  
func (this *TComMvField)  GetRefIdx() int8 { 
	return  this.m_iRefIdx;       
}
func (this *TComMvField)  GetHor   () int16 { 
	return  this.m_acMv.GetHor(); 
}
func (this *TComMvField)  GetVer   () int16 { 
	return  this.m_acMv.GetVer(); 
}


/// class for motion information in one CU
type TComCUMvField struct{
//private:
  m_pcMv		[]TComMv;
  m_pcMvd		[]TComMv;
  m_piRefIdx	[]int8;
  m_uiNumPartition	uint;
  m_cAMVPInfo	AMVPInfo;
}
/*    
  template <typename T>
  Void setAll( T *p, T const & val, PartSize eCUMode, Int iPartAddr, UInt uiDepth, Int iPartIdx );
*/
//public:
func NewTComCUMvField()*TComCUMvField{ 
	return &TComCUMvField{};
}
  
  // ------------------------------------------------------------------------------------------------------------------
  // create / destroy
  // ------------------------------------------------------------------------------------------------------------------
  
func (this *TComCUMvField)  Create(  uiNumPartition uint){
  this.m_pcMv     = make([]TComMv, uiNumPartition);
  this.m_pcMvd    = make([]TComMv, uiNumPartition);
  this.m_piRefIdx = make([]int8,   uiNumPartition);
  
  this.m_uiNumPartition = uiNumPartition;
}

func (this *TComCUMvField)  Destroy(){
  this.m_pcMv     = nil;
  this.m_pcMvd    = nil;
  this.m_piRefIdx = nil;
  
  this.m_uiNumPartition = 0;
}
  
  // ------------------------------------------------------------------------------------------------------------------
  // clear / copy
  // ------------------------------------------------------------------------------------------------------------------

func (this *TComCUMvField)  ClearMvField(){
}
  
func (this *TComCUMvField)  CopyFrom( pcCUMvFieldSrc *TComCUMvField,  uiNumPartSrc uint,  iPartAddrDst int){
}
func (this *TComCUMvField)  CopyTo2 ( pcCUMvFieldDst *TComCUMvField,  iPartAddrDst int){
}
func (this *TComCUMvField)  CopyTo4 ( pcCUMvFieldDst *TComCUMvField,  iPartAddrDst int,  uiOffset,  uiNumPart uint){
}
  // ------------------------------------------------------------------------------------------------------------------
  // get
  // ------------------------------------------------------------------------------------------------------------------
func (this *TComCUMvField)  GetMvs     ( offset int) []TComMv { 
	return  this.m_pcMv[offset:]; 
}
func (this *TComCUMvField)  GetMvds    ( offset int) []TComMv { 
	return  this.m_pcMvd[offset:]; 
}
func (this *TComCUMvField)  GetRefIdxs( offset int) []int8 { 
	return  this.m_piRefIdx[offset:]; 
}
  
    
func (this *TComCUMvField)  GetMv     ( iIdx int) *TComMv { 
	return  &this.m_pcMv    [iIdx]; 
}
func (this *TComCUMvField)  GetMvd    ( iIdx int) *TComMv { 
	return  &this.m_pcMvd   [iIdx]; 
}
func (this *TComCUMvField)  GetRefIdx( iIdx int) int8 { 
	return  this.m_piRefIdx[iIdx]; 
}
  
func (this *TComCUMvField)  GetAMVPInfo () *AMVPInfo{ 
	return &this.m_cAMVPInfo; 
}
  
  // ------------------------------------------------------------------------------------------------------------------
  // set
  // ------------------------------------------------------------------------------------------------------------------
  
func (this *TComCUMvField)  SetAllMv     ( rcMv  *TComMv,  		 eCUMode PartSize,  iPartAddr int,  uiDepth uint,  iPartIdx int ){
}
func (this *TComCUMvField)  SetAllMvd    ( rcMvd *TComMv,  	     eCUMode PartSize,  iPartAddr int,  uiDepth uint,  iPartIdx int ){
}
func (this *TComCUMvField)  SetAllRefIdx ( iRefIdx int,          eCUMode PartSize,  iPartAddr int,  uiDepth uint,  iPartIdx int ){
}
func (this *TComCUMvField)  SetAllMvField( mvField *TComMvField, eCUMode PartSize,  iPartAddr int,  uiDepth uint,  iPartIdx int ){
}

func (this *TComCUMvField) SetNumPartition( uiNumPart uint){
    this.m_uiNumPartition = uiNumPart;
}
 
func (this *TComCUMvField) LinkToWithOffset( src *TComCUMvField,  offset int){
    this.m_pcMv     = src.GetMvs(offset);
    this.m_pcMvd    = src.GetMvds(offset);
    this.m_piRefIdx = src.GetRefIdxs(offset);
}
  
func (this *TComCUMvField) Compress(pePredMode []PredMode, scale int){ 
}

