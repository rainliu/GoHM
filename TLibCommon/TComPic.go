package TLibCommon

import (
    "container/list"
)

// ====================================================================================================================
// Class definition
// ====================================================================================================================

/// picture class (symbol + YUV buffers)
type TComPic struct {
    //private:
    m_uiTLayer          uint        //  Temporal layer
    m_bUsedByCurr       bool        //  Used by current picture
    m_bIsLongTerm       bool        //  IS long term picture
    m_bIsUsedAsLongTerm bool        //  long term picture is used as reference before
    m_apcPicSym         *TComPicSym //  Symbol

    m_apcPicYuv [2]*TComPicYuv //  Texture,  0:org / 1:rec

    m_pcPicYuvPred                          *TComPicYuv //  Prediction
    m_pcPicYuvResi                          *TComPicYuv //  Residual
    m_bReconstructed                        bool
    m_bNeededForOutput                      bool
    m_uiCurrSliceIdx                        uint // Index of current slice
    m_pSliceSUMap                           []int
    m_pbValidSlice                          []bool
    m_sliceGranularityForNDBFilter          int
    m_bIndependentSliceBoundaryForNDBFilter bool
    m_bIndependentTileBoundaryForNDBFilter  bool
    m_pNDBFilterYuvTmp                      *TComPicYuv //!< temporary picture buffer when non-cross slice/tile boundary in-loop filtering is enabled
    m_bCheckLTMSB                           bool
    m_numReorderPics				[MAX_TLAYER]int;
  	m_croppingWindow				*CroppingWindow;
  
    m_vSliceCUDataLink                      *list.List //std::vector<std::vector<TComDataCU*> > ;

    m_SEIs		*SEImessages; ///< Any SEI messages that have been received.  If !NULL we own the object.
}

//public:
func NewTComPic() *TComPic {
    return &TComPic{}
}

func (this *TComPic) Create(iWidth, iHeight int, uiMaxWidth, uiMaxHeight, uiMaxDepth uint,
    croppingWindow *CroppingWindow, numReorderPics []int, bIsVirtual bool) {
  this.m_apcPicSym = NewTComPicSym();  
  this.m_apcPicSym.Create( iWidth, iHeight, uiMaxWidth, uiMaxHeight, uiMaxDepth );
  if !bIsVirtual {
    this.m_apcPicYuv[0] = NewTComPicYuv();  
    this.m_apcPicYuv[0].Create( iWidth, iHeight, uiMaxWidth, uiMaxHeight, uiMaxDepth );
  }
  this.m_apcPicYuv[1] = NewTComPicYuv();  
  this.m_apcPicYuv[1].Create( iWidth, iHeight, uiMaxWidth, uiMaxHeight, uiMaxDepth );
  
  /* there are no SEI messages associated with this picture initially */
  this.m_SEIs = nil;
  this.m_bUsedByCurr = false;

  /* store cropping parameters with picture */
  this.m_croppingWindow = croppingWindow;

  /* store number of reorder pics with picture */
  //memcpy(m_numReorderPics, numReorderPics, MAX_TLAYER*sizeof(Int));

  return;
}
func (this *TComPic) Destroy() {
}

func (this *TComPic) GetTLayer() uint {
    return this.m_uiTLayer
}
func (this *TComPic) SetTLayer(uiTLayer uint) {
    this.m_uiTLayer = uiTLayer
}

func (this *TComPic) GetUsedByCurr() bool         {
    return this.m_bUsedByCurr
}
func (this *TComPic) SetUsedByCurr(bUsed bool)    {
    this.m_bUsedByCurr = bUsed
}
func (this *TComPic) GetIsLongTerm() bool         {
    return this.m_bIsLongTerm
}
func (this *TComPic) SetIsLongTerm(lt bool)       {
    this.m_bIsLongTerm = lt
}
func (this *TComPic) GetIsUsedAsLongTerm() bool   {
    return this.m_bIsUsedAsLongTerm
}
func (this *TComPic) SetIsUsedAsLongTerm(lt bool) {
    this.m_bIsUsedAsLongTerm = lt
}
func (this *TComPic) SetCheckLTMSBPresent(b bool) {
    this.m_bCheckLTMSB = b
}
func (this *TComPic) GetCheckLTMSBPresent() bool  {
    return this.m_bCheckLTMSB
}

func (this *TComPic) GetPicSym() *TComPicSym {
    return this.m_apcPicSym
}

func (this *TComPic) GetSlice(i uint) *TComSlice {
    return this.m_apcPicSym.GetSlice(i)
}


func (this *TComPic)  GetPOC()        uint    { 
	return  uint(this.m_apcPicSym.GetSlice(this.m_uiCurrSliceIdx).GetPOC());  
}
func (this *TComPic)  GetCU( uiCUAddr uint ) *TComDataCU { 
	return  this.m_apcPicSym.GetCU( uiCUAddr ); 
}

func (this *TComPic)  GetPicYuvOrg()  *TComPicYuv      { 
	return  this.m_apcPicYuv[0]; 
}
func (this *TComPic)  GetPicYuvRec()  *TComPicYuv      { 
	return  this.m_apcPicYuv[1]; 
}

func (this *TComPic)  GetPicYuvPred() *TComPicYuv      { 
	return  this.m_pcPicYuvPred; 
}
func (this *TComPic)  GetPicYuvResi() *TComPicYuv      { 
	return  this.m_pcPicYuvResi; 
}
func (this *TComPic)  SetPicYuvPred( pcPicYuv *TComPicYuv )       { 
	this.m_pcPicYuvPred = pcPicYuv; 
}
func (this *TComPic)  SetPicYuvResi( pcPicYuv *TComPicYuv )       { 
	this.m_pcPicYuvResi = pcPicYuv; 
}

func (this *TComPic)  GetNumCUsInFrame()    uint  { 
	return this.m_apcPicSym.GetNumberOfCUsInFrame(); 
}
func (this *TComPic)  GetNumPartInWidth()   uint  { 
	return this.m_apcPicSym.GetNumPartInWidth();     
}
func (this *TComPic)  GetNumPartInHeight()  uint  { 
	return this.m_apcPicSym.GetNumPartInHeight();    
}
func (this *TComPic)  GetNumPartInCU()      uint  { 
	return this.m_apcPicSym.GetNumPartition();       
}
func (this *TComPic)  GetFrameWidthInCU()   uint  { 
	return this.m_apcPicSym.GetFrameWidthInCU();     
}
func (this *TComPic)  GetFrameHeightInCU()  uint  { 
	return this.m_apcPicSym.GetFrameHeightInCU();    
}
func (this *TComPic)  GetMinCUWidth()       uint  { 
	return this.m_apcPicSym.GetMinCUWidth();         
}
func (this *TComPic)  GetMinCUHeight()      uint  { 
	return this.m_apcPicSym.GetMinCUHeight();        
}

func (this *TComPic)  GetParPelX(uhPartIdx	byte) uint{ 
	return this.GetParPelX(uhPartIdx); 
}
func (this *TComPic)  GetParPelY(uhPartIdx	byte) uint{ 
	return this.GetParPelX(uhPartIdx); 
}

func (this *TComPic)  GetStride()     int      { 
	return this.m_apcPicYuv[1].GetStride(); 
}
func (this *TComPic)  GetCStride()    int      { 
	return this.m_apcPicYuv[1].GetCStride(); 
}

func (this *TComPic)  SetReconMark (b bool) { 
	this.m_bReconstructed = b;     
}
func (this *TComPic)  GetReconMark () bool  { 
	return this.m_bReconstructed;  
}
func (this *TComPic)  SetOutputMark (b bool) { 
	this.m_bNeededForOutput = b;     
}
func (this *TComPic)  GetOutputMark () bool  { 
	return this.m_bNeededForOutput;  
}
func (this *TComPic)  SetNumReorderPics(i int, tlayer uint) { 
	this.m_numReorderPics[tlayer] = i;    
}
func (this *TComPic)  GetNumReorderPics(tlayer uint) int       { 
	return this.m_numReorderPics[tlayer]; 
}

func (this *TComPic)  CompressMotion(){
}
func (this *TComPic)  GetCurrSliceIdx() uint       { 
	return this.m_uiCurrSliceIdx;                
}
func (this *TComPic)  SetCurrSliceIdx(i uint)      { 
	this.m_uiCurrSliceIdx = i;                   
}
func (this *TComPic)  GetNumAllocatedSlice() uint  {
	return this.m_apcPicSym.GetNumAllocatedSlice();
}
func (this *TComPic)  AllocateNewSlice()           {
	this.m_apcPicSym.AllocateNewSlice();         
}
func (this *TComPic)  ClearSliceBuffer()           {
	this.m_apcPicSym.ClearSliceBuffer();         
}

func (this *TComPic)  GetCroppingWindow() *CroppingWindow        { 
	return this.m_croppingWindow; 
}

func (this *TComPic)  CreateNonDBFilterInfo   (sliceStartAddress *list.List, sliceGranularityDepth int, LFCrossSliceBoundary *list.List, numTiles int, bNDBFilterCrossTileBoundary bool){
}
func (this *TComPic)  CreateNonDBFilterInfoLCU(tileID, sliceID int, pcCU *TComDataCU, startSU, endSU uint, sliceGranularyDepth int, picWidth, picHeight uint){
}
func (this *TComPic)  DestroyNonDBFilterInfo(){
}

func (this *TComPic)  GetValidSlice                                  (sliceID int) bool {
	return this.m_pbValidSlice[sliceID];
}
func (this *TComPic)  GetIndependentSliceBoundaryForNDBFilter        ()            bool {
	return this.m_bIndependentSliceBoundaryForNDBFilter;
}
func (this *TComPic)  GetIndependentTileBoundaryForNDBFilter         ()            bool {
	return this.m_bIndependentTileBoundaryForNDBFilter; 
}
func (this *TComPic)  GetYuvPicBufferForIndependentBoundaryProcessing() *TComPicYuv            {
	return this.m_pNDBFilterYuvTmp;
}
func (this *TComPic)  GetOneSliceCUDataForNDBFilter      (sliceID int) *list.List{ 
	return nil;//this.m_vSliceCUDataLink[sliceID];
}

  // transfer ownership of seis to this picture
func (this *TComPic)  SetSEIs(seis *SEImessages) { 
	this.m_SEIs = seis; 
}

  //return the current list of SEI messages associated with this picture.
  //Pointer is valid until this.destroy() is called
func (this *TComPic)  GetSEIs() *SEImessages { 
	return this.m_SEIs; 
}
