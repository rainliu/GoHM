package TLibDecoder

import (
	"io"
    "gohm/TLibCommon"
)

// ====================================================================================================================
// Class definition
// ====================================================================================================================

/// entropy decoder pure class
type TDecEntropyIf interface {
    //public:
      //  Virtual list for SBAC/CAVLC
    ResetEntropy          ( pcSlice *TLibCommon.TComSlice);
    SetBitstream          ( p *TLibCommon.TComInputBitstream);
    SetTraceFile 		  ( traceFile io.Writer);
    SetSliceTrace 		  ( bSliceTrace bool);

    ParseVPS                  ( pcVPS *TLibCommon.TComVPS );
    ParseSPS                  ( pcSPS *TLibCommon.TComSPS );
    ParsePPS                  ( pcPPS *TLibCommon.TComPPS );

    ParseSliceHeader          ( rpcSlice *TLibCommon.TComSlice, parameterSetManager *TLibCommon.ParameterSetManager);

    ParseTerminatingBit       ( ruilsLast *uint );

    ParseMVPIdx        ( riMVPIdx *int );


    ParseSkipFlag      ( pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint );
    ParseCUTransquantBypassFlag( pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint );
    ParseSplitFlag     ( pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint );
    ParseMergeFlag     ( pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth, uiPUIdx uint );
    ParseMergeIndex    ( pcCU *TLibCommon.TComDataCU, ruiMergeIndex *uint, uiAbsPartIdx, uiDepth uint );
    ParsePartSize      ( pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint );
    ParsePredMode      ( pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint );

    ParseIntraDirLumaAng( pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint );

    ParseIntraDirChroma( pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint );

    ParseInterDir      ( pcCU *TLibCommon.TComDataCU, ruiInterDir *uint, uiAbsPartIdx, uiDepth uint );
    ParseRefFrmIdx     ( pcCU *TLibCommon.TComDataCU, riRefFrmIdx *int, uiAbsPartIdx, uiDepth uint, eRefList TLibCommon.RefPicList );
    ParseMvd           ( pcCU *TLibCommon.TComDataCU, uiAbsPartAddr, uiPartIdx, uiDepth uint, eRefList TLibCommon.RefPicList );

    ParseTransformSubdivFlag( ruiSubdivFlag *uint,  uiLog2TransformBlockSize uint);
    ParseQtCbf         ( pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint,  eType TLibCommon.TextType,  uiTrDepth,  uiDepth uint );
    ParseQtRootCbf     ( pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint, uiQtRootCbf *uint );

    ParseDeltaQP       ( pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint ) ;

    ParseIPCMInfo     ( pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint);

    ParseCoeffNxN( pcCU *TLibCommon.TComDataCU, pcCoef []TLibCommon.TCoeff, uiAbsPartIdx, uiWidth, uiHeight, uiDepth uint,  eTType TLibCommon.TextType);
    ParseTransformSkipFlags ( pcCU *TLibCommon.TComDataCU,  uiAbsPartIdx, width,  height, uiDepth uint,  eTType TLibCommon.TextType);
    UpdateContextTables(  eSliceType TLibCommon.SliceType, iQp int ) ;
}


/// entropy decoder class
type TDecEntropy struct {
    //private:
    m_pcEntropyDecoderIf	TDecEntropyIf;
    m_pcPrediction      *TLibCommon.TComPrediction
    m_uiBakAbsPartIdx   uint
    m_uiBakChromaOffset uint
    m_bakAbsPartIdxCU   uint
}

func NewTDecEntropy() *TDecEntropy{
	return &TDecEntropy{};
}

func (this *TDecEntropy) Init(p *TLibCommon.TComPrediction) {
    this.m_pcPrediction = p
}

func (this *TDecEntropy) DecodePUWise       ( pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint, pcSubCU *TLibCommon.TComDataCU){
}
func (this *TDecEntropy) DecodeInterDirPU   ( pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth, uiPartIdx uint ){
}

func (this *TDecEntropy) DecodeRefFrmIdxPU  ( pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth, uiPartIdx uint,  eRefList TLibCommon.RefPicList){
}
func (this *TDecEntropy) DecodeMvdPU        ( pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth, uiPartIdx uint,   eRefList TLibCommon.RefPicList){
}
func (this *TDecEntropy) DecodeMVPIdxPU     ( pcSubCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth, uiPartIdx uint, eRefList TLibCommon.RefPicList){
}

func (this *TDecEntropy) SetEntropyDecoder(p TDecEntropyIf) {
	this.m_pcEntropyDecoderIf = p;
}
func (this *TDecEntropy) SetBitstream(p *TLibCommon.TComInputBitstream) {
   this.m_pcEntropyDecoderIf.SetBitstream(p);
}
func (this *TDecEntropy) SetTraceFile( traceFile io.Writer){
   this.m_pcEntropyDecoderIf.SetTraceFile(traceFile);
}
func (this *TDecEntropy) SetSliceTrace( bSliceTrace bool){
   this.m_pcEntropyDecoderIf.SetSliceTrace(bSliceTrace);
}
func (this *TDecEntropy)   ResetEntropy                ( p 		*TLibCommon.TComSlice)           {
	this.m_pcEntropyDecoderIf.ResetEntropy(p);
}
func (this *TDecEntropy)   DecodeVPS                   ( pcVPS 	*TLibCommon.TComVPS) {
	this.m_pcEntropyDecoderIf.ParseVPS(pcVPS);
}
func (this *TDecEntropy)   DecodeSPS                   ( pcSPS  *TLibCommon.TComSPS)    {
	this.m_pcEntropyDecoderIf.ParseSPS(pcSPS);
}
func (this *TDecEntropy)   DecodePPS                   ( pcPPS	*TLibCommon.TComPPS, parameterSet *TLibCommon.ParameterSetManager )    {
	this.m_pcEntropyDecoderIf.ParsePPS(pcPPS);
}
func (this *TDecEntropy)   DecodeSliceHeader           ( rpcSlice	*TLibCommon.TComSlice, parameterSetManager	*TLibCommon.ParameterSetManager)  {
	this.m_pcEntropyDecoderIf.ParseSliceHeader(rpcSlice, parameterSetManager);
}

func (this *TDecEntropy)   DecodeTerminatingBit        ( ruiIsLast *uint )       {
	this.m_pcEntropyDecoderIf.ParseTerminatingBit(ruiIsLast);
}

func (this *TDecEntropy)   GetEntropyDecoder() TDecEntropyIf {
	return this.m_pcEntropyDecoderIf;
}

//public:
func (this *TDecEntropy)   DecodeSplitFlag         		( pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint ){
}
func (this *TDecEntropy)   DecodeSkipFlag          		( pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint ){
}
func (this *TDecEntropy)   DecodeCUTransquantBypassFlag	( pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint ){
}
func (this *TDecEntropy)   DecodeMergeFlag         		( pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth, uiPUIdx uint ){
}
func (this *TDecEntropy)   DecodeMergeIndex        ( pcSubCU *TLibCommon.TComDataCU, uiPartIdx, uiPartAddr uint, eCUMode TLibCommon.PartSize, puhInterDirNeighbours *byte, pcMvFieldNeighbours *TLibCommon.TComMvField, uiDepth uint){
}
func (this *TDecEntropy)   DecodePredMode          ( pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint ){
}
func (this *TDecEntropy)   DecodePartSize          ( pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint ){
}

func (this *TDecEntropy)   DecodeIPCMInfo          ( pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint ){
}

func (this *TDecEntropy)   DecodePredInfo          ( pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint, pcSubCU *TLibCommon.TComDataCU){
}

func (this *TDecEntropy)   DecodeIntraDirModeLuma  ( pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint ){
}
func (this *TDecEntropy)   DecodeIntraDirModeChroma( pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth uint ){
}

func (this *TDecEntropy)   DecodeQP                ( pcCU *TLibCommon.TComDataCU, uiAbsPartIdx uint){
}

func (this *TDecEntropy)   UpdateContextTables     ( eSliceType TLibCommon.SliceType, iQp int ) {
	this.m_pcEntropyDecoderIf.UpdateContextTables( eSliceType, iQp );
}
func (this *TDecEntropy)   DecodeCoeff             ( pcCU *TLibCommon.TComDataCU, uiAbsPartIdx, uiDepth, uiWidth, uiHeight uint, bCodeDQP *bool){
}


//private:
func (this *TDecEntropy)   xDecodeTransform        ( pcCU *TLibCommon.TComDataCU, offsetLuma, offsetChroma, uiAbsPartIdx, absTUPartIdx, uiDepth, width, height, uiTrIdx, uiInnerQuadIdx uint, bCodeDQP *bool ){
}

