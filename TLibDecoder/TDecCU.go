package TLibDecoder

import (
    "gohm/TLibCommon"
)

// ====================================================================================================================
// Class definition
// ====================================================================================================================

/// CU decoder class
type TDecCu struct {
    //private:
    m_uiMaxDepth 	uint 			  ///< max. number of depth
    m_ppcYuvResi 	[]*TLibCommon.TComYuv;       ///< array of residual buffer
    m_ppcYuvReco	[]*TLibCommon.TComYuv;       ///< array of prediction & reconstruction buffer
    m_ppcCU		    []*TLibCommon.TComDataCU;    ///< CU data array

    // access channel
    m_pcTrQuant        *TLibCommon.TComTrQuant
    m_pcPrediction     *TLibCommon.TComPrediction
    m_pcEntropyDecoder *TDecEntropy

    m_bDecodeDQP bool
}

func NewTDecCu() *TDecCu{
	return &TDecCu{};
}

  /// initialize access channels
func (this *TDecCu) Init  ( pcEntropyDecoder *TDecEntropy, pcTrQuant *TLibCommon.TComTrQuant, pcPrediction *TLibCommon.TComPrediction){
  this.m_pcEntropyDecoder  = pcEntropyDecoder;
  this.m_pcTrQuant         = pcTrQuant;
  this.m_pcPrediction      = pcPrediction;
}

  /// create internal buffers
func (this *TDecCu) Create  (  uiMaxDepth,  uiMaxWidth,  uiMaxHeight uint){
  this.m_uiMaxDepth = uiMaxDepth+1;

  this.m_ppcYuvResi = make([]*TLibCommon.TComYuv, 	  this.m_uiMaxDepth-1);
  this.m_ppcYuvReco = make([]*TLibCommon.TComYuv,	  this.m_uiMaxDepth-1);
  this.m_ppcCU      = make([]*TLibCommon.TComDataCU, this.m_uiMaxDepth-1);

  var uiNumPartitions uint;
  for ui := uint(0); ui < this.m_uiMaxDepth-1; ui++ {
    uiNumPartitions = 1<<( ( this.m_uiMaxDepth - ui - 1 )<<1 );
    uiWidth  := uiMaxWidth  >> ui;
    uiHeight := uiMaxHeight >> ui;

    this.m_ppcYuvResi[ui] = TLibCommon.NewTComYuv();
    this.m_ppcYuvResi[ui].Create( uiWidth, uiHeight );
    this.m_ppcYuvReco[ui] = TLibCommon.NewTComYuv();
    this.m_ppcYuvReco[ui].Create( uiWidth, uiHeight );
    this.m_ppcCU     [ui] = TLibCommon.NewTComDataCU();
    this.m_ppcCU     [ui].Create( uiNumPartitions, uiWidth, uiHeight, true, int(uiMaxWidth >> (this.m_uiMaxDepth - 1)), false );
  }

  this.m_bDecodeDQP = false;

  // initialize partition order.
  piTmp := uint(0);
  TLibCommon.InitZscanToRaster(int(this.m_uiMaxDepth), 1, 0, TLibCommon.G_auiZscanToRaster[:], &piTmp);
  TLibCommon.InitRasterToZscan( uiMaxWidth, uiMaxHeight, this.m_uiMaxDepth );

  // initialize conversion matrix from partition index to pel
  TLibCommon.InitRasterToPelXY( uiMaxWidth, uiMaxHeight, this.m_uiMaxDepth );
//#if !LINEBUF_CLEANUP
//  initMotionReferIdx ( uiMaxWidth, uiMaxHeight, this.m_uiMaxDepth );
//#endif
}

/// destroy internal buffers
func (this *TDecCu) Destroy() {
  for ui := uint(0); ui < this.m_uiMaxDepth-1; ui++ {
    this.m_ppcYuvResi[ui].Destroy();
    //delete m_ppcYuvResi[ui];
    this.m_ppcYuvResi[ui] = nil;
    this.m_ppcYuvReco[ui].Destroy();
    //delete m_ppcYuvReco[ui];
    this.m_ppcYuvReco[ui] = nil;
    this.m_ppcCU     [ui].Destroy();
    //delete m_ppcCU     [ui];
    this.m_ppcCU     [ui] = nil;
  }

  //delete [] m_ppcYuvResi;
  this.m_ppcYuvResi = nil;
  //delete [] m_ppcYuvReco;
  this.m_ppcYuvReco = nil;
  //delete [] m_ppcCU     ;
  this.m_ppcCU      = nil;
}


  /// decode CU information
func (this *TDecCu)  DecodeCU                ( pcCU *TLibCommon.TComDataCU, ruiIsLast *uint){
  if pcCU.GetSlice().GetPPS().GetUseDQP() {
    this.SetdQPFlag(true);
  }

//#if !REMOVE_BURST_IPCM
//  pcCU.SetNumSucIPCM(0);
//#endif

  // start from the top level CU
  this.xDecodeCU( pcCU, 0, 0, ruiIsLast);
}

  /// reconstruct CU information
func (this *TDecCu)  DecompressCU            ( pcCU *TLibCommon.TComDataCU ){
	this.xDecompressCU( pcCU, pcCU, 0,  0 );
}

func (this *TDecCu)  xDecodeCU               ( pcCU *TLibCommon.TComDataCU,  uiAbsPartIdx,  uiDepth uint,  ruiIsLast *uint){
  pcPic := pcCU.GetPic();
  uiCurNumParts := pcPic.GetNumPartInCU() >> (uiDepth<<1);
  uiQNumParts   := uiCurNumParts>>2;
  
  bBoundary := false;
  uiLPelX   := pcCU.GetCUPelX() + TLibCommon.G_auiRasterToPelX[ TLibCommon.G_auiZscanToRaster[uiAbsPartIdx] ];
  uiRPelX   := uiLPelX + (TLibCommon.G_uiMaxCUWidth>>uiDepth)  - 1;
  uiTPelY   := pcCU.GetCUPelY() + TLibCommon.G_auiRasterToPelY[ TLibCommon.G_auiZscanToRaster[uiAbsPartIdx] ];
  uiBPelY   := uiTPelY + (TLibCommon.G_uiMaxCUHeight>>uiDepth) - 1;
  
  pcSlice := pcCU.GetPic().GetSlice(pcCU.GetPic().GetCurrSliceIdx());
  bStartInCU := pcCU.GetSCUAddr()+uiAbsPartIdx+uiCurNumParts>pcSlice.GetDependentSliceCurStartCUAddr()&&pcCU.GetSCUAddr()+uiAbsPartIdx<pcSlice.GetDependentSliceCurStartCUAddr();
  if (!bStartInCU) && ( uiRPelX < pcSlice.GetSPS().GetPicWidthInLumaSamples() ) && ( uiBPelY < pcSlice.GetSPS().GetPicHeightInLumaSamples() ) {
/*#if !REMOVE_BURST_IPCM
    if(pcCU.GetNumSucIPCM() == 0)
    {
      m_pcEntropyDecoder->decodeSplitFlag( pcCU, uiAbsPartIdx, uiDepth );
    }
    else
    {
      pcCU.SetDepthSubParts( uiDepth, uiAbsPartIdx );
    }
#else*/
    this.m_pcEntropyDecoder.DecodeSplitFlag( pcCU, uiAbsPartIdx, uiDepth );
//#endif
  }else{
    bBoundary = true;
  }
  
  if ( ( uiDepth < uint(pcCU.GetDepth1( uiAbsPartIdx )) ) && ( uiDepth < TLibCommon.G_uiMaxCUDepth - TLibCommon.G_uiAddCUDepth ) ) || bBoundary {
    uiIdx := uiAbsPartIdx;
    if (TLibCommon.G_uiMaxCUWidth>>uiDepth) == pcCU.GetSlice().GetPPS().GetMinCuDQPSize() && pcCU.GetSlice().GetPPS().GetUseDQP() {
      this.SetdQPFlag(true);
      pcCU.SetQPSubParts( int(pcCU.GetRefQP(uiAbsPartIdx)), uiAbsPartIdx, uiDepth ); // set QP to default QP
    }

    for uiPartUnitIdx := uint(0); uiPartUnitIdx < 4; uiPartUnitIdx++ {
      uiLPelX   = pcCU.GetCUPelX() + TLibCommon.G_auiRasterToPelX[ TLibCommon.G_auiZscanToRaster[uiIdx] ];
      uiTPelY   = pcCU.GetCUPelY() + TLibCommon.G_auiRasterToPelY[ TLibCommon.G_auiZscanToRaster[uiIdx] ];
      
      bSubInSlice := pcCU.GetSCUAddr()+uiIdx+uiQNumParts>pcSlice.GetDependentSliceCurStartCUAddr();
      if bSubInSlice {
        if ( uiLPelX < pcCU.GetSlice().GetSPS().GetPicWidthInLumaSamples() ) && ( uiTPelY < pcCU.GetSlice().GetSPS().GetPicHeightInLumaSamples() )  {
          this.xDecodeCU( pcCU, uiIdx, uiDepth+1, ruiIsLast );
        }else{
          pcCU.SetOutsideCUPart( uiIdx, uiDepth+1 );
        }
      }
      if *ruiIsLast!=0 {
        break;
      }
      
      uiIdx += uiQNumParts;
    }
    if (TLibCommon.G_uiMaxCUWidth>>uiDepth) == pcCU.GetSlice().GetPPS().GetMinCuDQPSize() && pcCU.GetSlice().GetPPS().GetUseDQP() {
      if this.GetdQPFlag() {
        var  uiQPSrcPartIdx uint;
        if pcPic.GetCU( pcCU.GetAddr() ).GetDependentSliceStartCU(uiAbsPartIdx) != pcSlice.GetDependentSliceCurStartCUAddr() {
          uiQPSrcPartIdx = pcSlice.GetDependentSliceCurStartCUAddr() % pcPic.GetNumPartInCU();
        }else{
          uiQPSrcPartIdx = uiAbsPartIdx;
        }
        pcCU.SetQPSubParts( int(pcCU.GetRefQP( uiQPSrcPartIdx )), uiAbsPartIdx, uiDepth ); // set QP to default QP
      }
    }
    return;
  }
  
  if (TLibCommon.G_uiMaxCUWidth>>uiDepth) >= pcCU.GetSlice().GetPPS().GetMinCuDQPSize() && pcCU.GetSlice().GetPPS().GetUseDQP() {
    this.SetdQPFlag(true);
    pcCU.SetQPSubParts( int(pcCU.GetRefQP(uiAbsPartIdx)), uiAbsPartIdx, uiDepth ); // set QP to default QP
  }

//#if !REMOVE_BURST_IPCM
//  if (pcCU.GetSlice().GetPPS().GetTransquantBypassEnableFlag() && pcCU.GetNumSucIPCM() == 0 )
//#else
  if pcCU.GetSlice().GetPPS().GetTransquantBypassEnableFlag() {
//#endif
    this.m_pcEntropyDecoder.DecodeCUTransquantBypassFlag( pcCU, uiAbsPartIdx, uiDepth );
  }
  
  // decode CU mode and the partition size
//#if !REMOVE_BURST_IPCM
//  if( !pcCU.GetSlice()->isIntra() && pcCU.GetNumSucIPCM() == 0 )
//#else
  if !pcCU.GetSlice().IsIntra() {
//#endif
    this.m_pcEntropyDecoder.DecodeSkipFlag( pcCU, uiAbsPartIdx, uiDepth );
  }
 
  if pcCU.IsSkipped(uiAbsPartIdx) {
    this.m_ppcCU[uiDepth].CopyInterPredInfoFrom( pcCU, uiAbsPartIdx, TLibCommon.REF_PIC_LIST_0 );
    this.m_ppcCU[uiDepth].CopyInterPredInfoFrom( pcCU, uiAbsPartIdx, TLibCommon.REF_PIC_LIST_1 );
    var cMvFieldNeighbours		[TLibCommon.MRG_MAX_NUM_CANDS << 1]TLibCommon.TComMvField; // double length for mv of both lists
    var uhInterDirNeighbours	[TLibCommon.MRG_MAX_NUM_CANDS]byte;
    numValidMergeCand := 0;
    for ui := uint(0); ui < this.m_ppcCU[uiDepth].GetSlice().GetMaxNumMergeCand(); ui++ {
      uhInterDirNeighbours[ui] = 0;
    }
    this.m_pcEntropyDecoder.DecodeMergeIndex( pcCU, 0, uiAbsPartIdx, TLibCommon.SIZE_2Nx2N, uhInterDirNeighbours[:], cMvFieldNeighbours[:], uiDepth );
    uiMergeIndex := pcCU.GetMergeIndex1(uiAbsPartIdx);
    this.m_ppcCU[uiDepth].GetInterMergeCandidates( 0, 0, cMvFieldNeighbours[:], uhInterDirNeighbours[:], &numValidMergeCand, int(uiMergeIndex) );
    pcCU.SetInterDirSubParts( uint(uhInterDirNeighbours[uiMergeIndex]), uiAbsPartIdx, 0, uiDepth );

    cTmpMv := TLibCommon.NewTComMv( 0, 0 );
    for uiRefListIdx := 0; uiRefListIdx < 2; uiRefListIdx++ {        
      if pcCU.GetSlice().GetNumRefIdx( TLibCommon.RefPicList( uiRefListIdx ) ) > 0 {
        pcCU.SetMVPIdxSubParts( 0, TLibCommon.RefPicList( uiRefListIdx ), uiAbsPartIdx, 0, uiDepth);
        pcCU.SetMVPNumSubParts( 0, TLibCommon.RefPicList( uiRefListIdx ), uiAbsPartIdx, 0, uiDepth);
        pcCU.GetCUMvField( TLibCommon.RefPicList( uiRefListIdx ) ).SetAllMvd( cTmpMv, TLibCommon.SIZE_2Nx2N, int(uiAbsPartIdx), uiDepth, 0 );
        pcCU.GetCUMvField( TLibCommon.RefPicList( uiRefListIdx ) ).SetAllMvField( &cMvFieldNeighbours[ 2*int(uiMergeIndex) + uiRefListIdx ], TLibCommon.SIZE_2Nx2N, int(uiAbsPartIdx), uiDepth, 0 );
      }
    }
    this.xFinishDecodeCU( pcCU, uiAbsPartIdx, uiDepth, ruiIsLast );
    return;
  }

/*#if !REMOVE_BURST_IPCM
  if( pcCU.GetNumSucIPCM() == 0 ) 
  {
    m_pcEntropyDecoder->decodePredMode( pcCU, uiAbsPartIdx, uiDepth );
    m_pcEntropyDecoder->decodePartSize( pcCU, uiAbsPartIdx, uiDepth );
  }
  else
  {
    pcCU.SetPredModeSubParts( MODE_INTRA, uiAbsPartIdx, uiDepth );
    pcCU.SetPartSizeSubParts( SIZE_2Nx2N, uiAbsPartIdx, uiDepth );
    pcCU.SetSizeSubParts( g_uiMaxCUWidth>>uiDepth, g_uiMaxCUHeight>>uiDepth, uiAbsPartIdx, uiDepth ); 
    pcCU.SetTrIdxSubParts( 0, uiAbsPartIdx, uiDepth );
  }
#else*/
  this.m_pcEntropyDecoder.DecodePredMode( pcCU, uiAbsPartIdx, uiDepth );
  this.m_pcEntropyDecoder.DecodePartSize( pcCU, uiAbsPartIdx, uiDepth );
//#endif

  if pcCU.IsIntra( uiAbsPartIdx ) && pcCU.GetPartitionSize1( uiAbsPartIdx ) == TLibCommon.SIZE_2Nx2N  {
    this.m_pcEntropyDecoder.DecodeIPCMInfo( pcCU, uiAbsPartIdx, uiDepth );

    if pcCU.GetIPCMFlag1(uiAbsPartIdx) {
      this.xFinishDecodeCU( pcCU, uiAbsPartIdx, uiDepth, ruiIsLast );
      return;
    }
  }

  uiCurrWidth      := pcCU.GetWidth1 ( uiAbsPartIdx );
  uiCurrHeight     := pcCU.GetHeight1( uiAbsPartIdx );
  
  // prediction mode ( Intra : direction mode, Inter : Mv, reference idx )
  this.m_pcEntropyDecoder.DecodePredInfo( pcCU, uiAbsPartIdx, uiDepth, this.m_ppcCU[uiDepth]);
  
  // Coefficient decoding
  bCodeDQP := this.GetdQPFlag();
  this.m_pcEntropyDecoder.DecodeCoeff( pcCU, uiAbsPartIdx, uiDepth, uint(uiCurrWidth), uint(uiCurrHeight), &bCodeDQP );
  this.SetdQPFlag( bCodeDQP );
  this.xFinishDecodeCU( pcCU, uiAbsPartIdx, uiDepth, ruiIsLast );
}

func (this *TDecCu)  xFinishDecodeCU         ( pcCU *TLibCommon.TComDataCU,  uiAbsPartIdx,  uiDepth uint,  ruiIsLast *uint){
  if  pcCU.GetSlice().GetPPS().GetUseDQP() {
  	if this.GetdQPFlag() {
    	pcCU.SetQPSubParts( int(pcCU.GetRefQP(uiAbsPartIdx)), uiAbsPartIdx, uiDepth ); // set QP
    }else{
    	pcCU.SetQPSubParts( int(pcCU.GetCodedQP()), uiAbsPartIdx, uiDepth ); // set QP
    }
  }

/*#if !REMOVE_BURST_IPCM
  if( pcCU.GetNumSucIPCM() > 0 )
  {
    ruiIsLast = 0;
    return;
  }
#endif*/

  *ruiIsLast = uint(TLibCommon.B2U(this.xDecodeSliceEnd( pcCU, uiAbsPartIdx, uiDepth)));
}

func (this *TDecCu)  xDecodeSliceEnd         ( pcCU *TLibCommon.TComDataCU,                        uiAbsPartIdx,  uiDepth uint) bool{
  var uiIsLast uint;
  pcPic := pcCU.GetPic();
  pcSlice := pcPic.GetSlice(pcPic.GetCurrSliceIdx());
  uiCurNumParts := pcPic.GetNumPartInCU() >> (uiDepth<<1);
  uiWidth := pcSlice.GetSPS().GetPicWidthInLumaSamples();
  uiHeight := pcSlice.GetSPS().GetPicHeightInLumaSamples();
  uiGranularityWidth := TLibCommon.G_uiMaxCUWidth;
  uiPosX := pcCU.GetCUPelX() + TLibCommon.G_auiRasterToPelX[ TLibCommon.G_auiZscanToRaster[uiAbsPartIdx] ];
  uiPosY := pcCU.GetCUPelY() + TLibCommon.G_auiRasterToPelY[ TLibCommon.G_auiZscanToRaster[uiAbsPartIdx] ];

  if ((uiPosX+uint(pcCU.GetWidth1 (uiAbsPartIdx)))%uiGranularityWidth==0||(uiPosX+uint(pcCU.GetWidth1 (uiAbsPartIdx))==uiWidth)) &&
     ((uiPosY+uint(pcCU.GetHeight1(uiAbsPartIdx)))%uiGranularityWidth==0||(uiPosY+uint(pcCU.GetHeight1(uiAbsPartIdx))==uiHeight)) {
    this.m_pcEntropyDecoder.DecodeTerminatingBit( &uiIsLast );
  }else{
    uiIsLast=0;
  }

  if uiIsLast!=0 {
    if pcSlice.IsNextDependentSlice()&&!pcSlice.IsNextSlice() {
      pcSlice.SetDependentSliceCurEndCUAddr(pcCU.GetSCUAddr()+uiAbsPartIdx+uiCurNumParts);
    }else{
      pcSlice.SetSliceCurEndCUAddr(pcCU.GetSCUAddr()+uiAbsPartIdx+uiCurNumParts);
      pcSlice.SetDependentSliceCurEndCUAddr(pcCU.GetSCUAddr()+uiAbsPartIdx+uiCurNumParts);
    }
  }

  return uiIsLast>0;
}
func (this *TDecCu)  xDecompressCU           ( pcCU *TLibCommon.TComDataCU, pcCUCur *TLibCommon.TComDataCU,   uiAbsPartIdx,  uiDepth uint){
}

func (this *TDecCu)  xReconInter             ( pcCU *TLibCommon.TComDataCU, uiAbsPartIdx,  uiDepth uint ){
}

func (this *TDecCu)  xReconIntraQT           ( pcCU *TLibCommon.TComDataCU, uiAbsPartIdx,  uiDepth uint ){
}
func (this *TDecCu)  xIntraRecLumaBlk        ( pcCU *TLibCommon.TComDataCU, uiTrDepth,  uiAbsPartIdx uint, pcRecoYuv *TLibCommon.TComYuv, pcPredYuv *TLibCommon.TComYuv, pcResiYuv *TLibCommon.TComYuv){
}
func (this *TDecCu)  xIntraRecChromaBlk      ( pcCU *TLibCommon.TComDataCU, uiTrDepth,  uiAbsPartIdx uint, pcRecoYuv *TLibCommon.TComYuv, pcPredYuv *TLibCommon.TComYuv, pcResiYuv *TLibCommon.TComYuv,  uiChromaId uint){
}

func (this *TDecCu)  xReconPCM               ( pcCU *TLibCommon.TComDataCU, uiAbsPartIdx,  uiDepth uint ){
}

func (this *TDecCu)  xDecodeInterTexture     ( pcCU *TLibCommon.TComDataCU, uiAbsPartIdx,  uiDepth uint ){
}
func (this *TDecCu)  xDecodePCMTexture       ( pcCU *TLibCommon.TComDataCU,  uiPartIdx uint, piPCM *TLibCommon.Pel, piReco *TLibCommon.Pel,  uiStride,  uiWidth,  uiHeight uint,  ttText TLibCommon.TextType){
}

func (this *TDecCu)  xCopyToPic              ( pcCU *TLibCommon.TComDataCU, pcPic *TLibCommon.TComPic,  uiZorderIdx,  uiDepth uint){
}

func (this *TDecCu)  xIntraLumaRecQT         ( pcCU *TLibCommon.TComDataCU,  uiTrDepth,  uiAbsPartIdx uint, pcRecoYuv *TLibCommon.TComYuv, pcPredYuv *TLibCommon.TComYuv, pcResiYuv *TLibCommon.TComYuv ){
}
func (this *TDecCu)  xIntraChromaRecQT       ( pcCU *TLibCommon.TComDataCU,  uiTrDepth,  uiAbsPartIdx uint, pcRecoYuv *TLibCommon.TComYuv, pcPredYuv *TLibCommon.TComYuv, pcResiYuv *TLibCommon.TComYuv ){
}

func (this *TDecCu) GetdQPFlag               ()         bool               {
	return this.m_bDecodeDQP;
}
func (this *TDecCu) SetdQPFlag               (  b bool)                {
	this.m_bDecodeDQP = b;
}
func (this *TDecCu) xFillPCMBuffer           ( pCU *TLibCommon.TComDataCU,  absPartIdx,  depth uint){
}
