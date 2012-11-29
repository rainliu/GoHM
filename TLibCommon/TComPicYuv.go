package TLibCommon

import (

)

/// picture YUV buffer class
type TComPicYuv struct{
  // ------------------------------------------------------------------------------------------------
  //  YUV buffer
  // ------------------------------------------------------------------------------------------------
  m_apiPicBufY *Pel;           ///< Buffer (including margin)
  m_apiPicBufU *Pel;
  m_apiPicBufV *Pel;
  
  m_piPicOrgY *Pel;            ///< m_apiPicBufY + m_iMarginLuma*getStride() + m_iMarginLuma
  m_piPicOrgU *Pel;
  m_piPicOrgV *Pel;
  
  // ------------------------------------------------------------------------------------------------
  //  Parameter for general YUV buffer usage
  // ------------------------------------------------------------------------------------------------
  
  m_iPicWidth		int;            ///< Width of picture
  m_iPicHeight		int;            ///< Height of picture
  
  m_iCuWidth		int;             ///< Width of Coding Unit (CU)
  m_iCuHeight		int;             ///< Height of Coding Unit (CU)
  m_cuOffsetY		*int;
  m_cuOffsetC		*int;
  m_buOffsetY		*int;
  m_buOffsetC		*int;
  
  m_iLumaMarginX	int;
  m_iLumaMarginY	int;
  m_iChromaMarginX	int;
  m_iChromaMarginY	int;
  
  m_bIsBorderExtended	bool;
}

 
//protected:
func (this *TComPicYuv) xExtendPicCompBorder (piTxt *Pel, iStride, iWidth, iHeight, iMarginX, iMarginY int){
}
  
//public:
func NewTComPicYuv() (*TComPicYuv){
	return &TComPicYuv{}
}
 
  // ------------------------------------------------------------------------------------------------
  //  Memory management
  // ------------------------------------------------------------------------------------------------
func (this *TComPicYuv) Create      ( iPicWidth, iPicHeight int, uiMaxCUWidth, uiMaxCUHeight, uiMaxCUDepth uint ){
}

func (this *TComPicYuv) Destroy     (){
}
  
func (this *TComPicYuv) CreateLuma  ( iPicWidth, iPicHeight int, uiMaxCUWidth, uiMaxCUHeight, uiMaxCUDepth uint ){
}

func (this *TComPicYuv)	DestroyLuma (){
}
  
  // ------------------------------------------------------------------------------------------------
  //  Get information of picture
  // ------------------------------------------------------------------------------------------------
  
func (this *TComPicYuv)	GetWidth    () int    { 
	return  this.m_iPicWidth;    
}

func (this *TComPicYuv) GetHeight   () int    { 
	return  this.m_iPicHeight;   
}
  
func (this *TComPicYuv) GetStride   () int    { 
	return (this.m_iPicWidth     ) + (this.m_iLumaMarginX  <<1); 
}

func (this *TComPicYuv) GetCStride  () int    { 
	return (this.m_iPicWidth >> 1) + (this.m_iChromaMarginX<<1); 
}
  
func (this *TComPicYuv) GetLumaMargin   () int{ 
	return this.m_iLumaMarginX;  
}

func (this *TComPicYuv) GetChromaMargin () int{ 
	return this.m_iChromaMarginX;
}
  
  // ------------------------------------------------------------------------------------------------
  //  Access function for picture buffer
  // ------------------------------------------------------------------------------------------------
  
  //  Access starting position of picture buffer with margin
func (this *TComPicYuv)   GetBufY     ()  *Pel   { 
	return  this.m_apiPicBufY;   
}

func (this *TComPicYuv)   GetBufU     ()  *Pel   { 
	return  this.m_apiPicBufU;   
}

func (this *TComPicYuv)   GetBufV     ()  *Pel   { 
	return  this.m_apiPicBufV;   
}
  
  //  Access starting position of original picture
func (this *TComPicYuv)   GetLumaAddr ()  *Pel   { 
	return  this.m_piPicOrgY;    
}

func (this *TComPicYuv)   GetCbAddr   ()  *Pel   { 
	return  this.m_piPicOrgU;    
}

func (this *TComPicYuv)   GetCrAddr   ()  *Pel   { 
	return  this.m_piPicOrgV;    
}
 
  //  Access starting position of original picture for specific coding unit (CU) or partition unit (PU)
/* 
func (this *TComPicYuv)   GetLumaAddr1 ( iCuAddr int ) *Pel{ 
	return this.m_piPicOrgY + this.m_cuOffsetY[ iCuAddr ]; 
}

func (this *TComPicYuv)   GetCbAddr1   ( iCuAddr int ) *Pel{ 
	return this.m_piPicOrgU + this.m_cuOffsetC[ iCuAddr ]; 
}

func (this *TComPicYuv)   GetCrAddr1   ( iCuAddr int ) *Pel{ 
	return this.m_piPicOrgV + this.m_cuOffsetC[ iCuAddr ]; 
}

func (this *TComPicYuv)   GetLumaAddr2 ( iCuAddr, uiAbsZorderIdx int ) *Pel{ 
	return this.m_piPicOrgY + this.m_cuOffsetY[iCuAddr] + this.m_buOffsetY[g_auiZscanToRaster[uiAbsZorderIdx]]; 
}

func (this *TComPicYuv)   GetCbAddr2   ( iCuAddr, uiAbsZorderIdx int ) *Pel{ 
	return this.m_piPicOrgU + this.m_cuOffsetC[iCuAddr] + this.m_buOffsetC[g_auiZscanToRaster[uiAbsZorderIdx]]; 
}

func (this *TComPicYuv)   GetCrAddr2   ( iCuAddr, uiAbsZorderIdx int ) *Pel{ 
	return this.m_piPicOrgV + this.m_cuOffsetC[iCuAddr] + this.m_buOffsetC[g_auiZscanToRaster[uiAbsZorderIdx]]; 
}
*/  
  // ------------------------------------------------------------------------------------------------
  //  Miscellaneous
  // ------------------------------------------------------------------------------------------------
  
  //  Copy function to picture
func (this *TComPicYuv)   CopyToPic       ( pcPicYuvDst *TComPicYuv ){
}

func (this *TComPicYuv)   CopyToPicLuma   ( pcPicYuvDst *TComPicYuv ){
}

func (this *TComPicYuv)   CopyToPicCb     ( pcPicYuvDst *TComPicYuv ){
}

func (this *TComPicYuv)   CopyToPicCr     ( pcPicYuvDst *TComPicYuv ){
}
  
  //  Extend function of picture buffer
func (this *TComPicYuv)   ExtendPicBorder      (){
}
  
  //  Dump picture
func (this *TComPicYuv)   Dump (pFileName string, bAdd bool){
}
  
  // Set border extension flag
func (this *TComPicYuv)   SetBorderExtension(bIsBorderExtended bool) { 
	this.m_bIsBorderExtended = bIsBorderExtended; 
}
