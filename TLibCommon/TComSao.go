package TLibCommon

import ()

// ====================================================================================================================
// Constants
// ====================================================================================================================
const SAO_MAX_DEPTH = 4
const SAO_BO_BITS = 5
const LUMA_GROUP_NUM = (1 << SAO_BO_BITS)
const MAX_NUM_SAO_OFFSETS = 4
const MAX_NUM_SAO_CLASS = 33

// ====================================================================================================================
// Class definition
// ====================================================================================================================

/// Sample Adaptive Offset class
type TComSampleAdaptiveOffset struct {
    //protected:
    m_pcPic *TComPic

    m_uiMaxDepth         uint
    m_aiNumCulPartsLevel [5]uint
    m_auiEoTable         [9]uint
    m_iOffsetBo          *int
    m_iChromaOffsetBo    *int
    m_iOffsetEo          [LUMA_GROUP_NUM]int

    m_iPicWidth           int
    m_iPicHeight          int
    m_uiMaxSplitLevel     uint
    m_uiMaxCUWidth        uint
    m_uiMaxCUHeight       uint
    m_iNumCuInWidth       int
    m_iNumCuInHeight      int
    m_iNumTotalParts      int
    m_iNumClass           [MAX_NUM_SAO_TYPE]int
    m_eSliceType          SliceType
    m_iPicNalReferenceIdc int

    m_uiSaoBitIncreaseY uint
    m_uiSaoBitIncreaseC uint //for chroma
    m_uiQP              uint

    m_pClipTable           *Pel
    m_pClipTableBase       *Pel
    m_lumaTableBo          *Pel
    m_pChromaClipTable     *Pel
    m_pChromaClipTableBase *Pel
    m_chromaTableBo        *Pel
    m_iUpBuff1             *int
    m_iUpBuff2             *int
    m_iUpBufft             *int
    m_ipSwap               *int
    m_bUseNIF              bool        //!< true for performing non-cross slice boundary ALF
    m_uiNumSlicesInPic     uint        //!< number of slices in picture
    m_iSGDepth             int         //!< slice granularity depth
    m_pcYuvTmp             *TComPicYuv //!< temporary picture buffer pointer when non-across slice/tile boundary SAO is enabled

    m_pTmpU1                  *Pel
    m_pTmpU2                  *Pel
    m_pTmpL1                  *Pel
    m_pTmpL2                  *Pel
    m_iLcuPartIdx             *int
    m_maxNumOffsetsPerPic     int
    m_saoLcuBoundary          bool
    m_saoLcuBasedOptimization bool
}


func (this *TComSampleAdaptiveOffset) xPCMRestoration        (pcPic *TComPic){
}
func (this *TComSampleAdaptiveOffset) xPCMCURestoration      (pcCU *TComDataCU,  uiAbsZorderIdx,  uiDepth uint){
}
func (this *TComSampleAdaptiveOffset) xPCMSampleRestoration  (pcCU *TComDataCU,  uiAbsZorderIdx,  uiDepth uint,  ttText TextType){
}
//public:
func NewTComSampleAdaptiveOffset() *TComSampleAdaptiveOffset{
	return &TComSampleAdaptiveOffset{}
}

func (this *TComSampleAdaptiveOffset) Create(uiSourceWidth, uiSourceHeight, uiMaxCUWidth, uiMaxCUHeight uint) {
}
func (this *TComSampleAdaptiveOffset) Destroy() {
}


func (this *TComSampleAdaptiveOffset) ConvertLevelRowCol2Idx( level,  row,  col int) int{
	return 0;
}

func (this *TComSampleAdaptiveOffset) InitSAOParam   (pcSaoParam *SAOParam,  iPartLevel,  iPartRow,  iPartCol,  iParentPartIdx,  StartCUX,  EndCUX,  StartCUY,  EndCUY,  iYCbCr int){
}
func (this *TComSampleAdaptiveOffset) AllocSaoParam  (pcSaoParam *SAOParam){
}
func (this *TComSampleAdaptiveOffset) ResetSAOParam  (pcSaoParam *SAOParam){
}
func (this *TComSampleAdaptiveOffset) FreeSaoParam   (pcSaoParam *SAOParam){
}

func (this *TComSampleAdaptiveOffset) SAOProcess(pcSaoParam *SAOParam){
}
func (this *TComSampleAdaptiveOffset) ProcessSaoCu( iAddr,  iSaoType,  iYCbCr int){
}
func (this *TComSampleAdaptiveOffset) GetPicYuvAddr(pcPicYuv *TComPicYuv,  iYCbCr, iAddr int) *Pel{
	return nil
}

func (this *TComSampleAdaptiveOffset) ProcessSaoCuOrg( iAddr,  iPartIdx,  iYCbCr int){  //!< LCU-basd SAO process without slice granularity
}
func (this *TComSampleAdaptiveOffset) CreatePicSaoInfo(pcPic *TComPic,  numSlicesInPic int){
}
func (this *TComSampleAdaptiveOffset) DestroyPicSaoInfo(){
}
func (this *TComSampleAdaptiveOffset) ProcessSaoBlock(pDec *Pel, pRest *Pel,  stride,  iSaoType int,  width,  height uint, pbBorderAvail *bool,  iYCbCr int){
}

func (this *TComSampleAdaptiveOffset) ResetLcuPart(saoLcuParam *SaoLcuParam){
}
func (this *TComSampleAdaptiveOffset) ConvertQT2SaoUnit(saoParam *SAOParam,  partIdx uint, yCbCr int){
}
func (this *TComSampleAdaptiveOffset) ConvertOnePart2SaoUnit(saoParam *SAOParam,  partIdx uint,  yCbCr int){
}
func (this *TComSampleAdaptiveOffset) ProcessSaoUnitAll(saoLcuParam *SaoLcuParam, oneUnitFlag bool,  yCbCr int){
}
func (this *TComSampleAdaptiveOffset) SetSaoLcuBoundary ( bVal bool)  {
	this.m_saoLcuBoundary = bVal;
}
func (this *TComSampleAdaptiveOffset) GetSaoLcuBoundary ()    bool       {
	return this.m_saoLcuBoundary;
}
func (this *TComSampleAdaptiveOffset) SetSaoLcuBasedOptimization ( bVal bool)  {
	this.m_saoLcuBasedOptimization = bVal;
}
func (this *TComSampleAdaptiveOffset) GetSaoLcuBasedOptimization ()  bool         {
	return this.m_saoLcuBasedOptimization;
}
func (this *TComSampleAdaptiveOffset) ResetSaoUnit(saoUnit *SaoLcuParam){
}
func (this *TComSampleAdaptiveOffset) CopySaoUnit(saoUnitDst *SaoLcuParam, saoUnitSrc *SaoLcuParam){
}
func (this *TComSampleAdaptiveOffset) PCMLFDisableProcess    ( pcPic *TComPic ){                       ///< interface function for ALF process
}
