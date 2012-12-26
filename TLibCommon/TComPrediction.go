package TLibCommon

import ()

// ====================================================================================================================
// Class definition
// ====================================================================================================================

// ====================================================================================================================
// Class definition
// ====================================================================================================================
/// weighting prediction class
type TComWeightPrediction struct {
    m_wp0 [3]wpScalingParam
    m_wp1 [3]wpScalingParam
}

func NewTComWeightPrediction() *TComWeightPrediction {
    return &TComWeightPrediction{}
}


func (this *TComWeightPrediction)  weightBidirY( w0 int, P0 Pel, w1 int, P1 Pel, round, shift, offset int) Pel{
  return ClipY( Pel(( (w0*(int(P0) + IF_INTERNAL_OFFS) + w1*(int(P1) + IF_INTERNAL_OFFS) + round + (offset << uint(shift-1))) >> uint(shift) )) );
}
func (this *TComWeightPrediction)  weightBidirC( w0 int, P0 Pel, w1 int, P1 Pel, round, shift, offset int) Pel{
  return ClipC( Pel(( (w0*(int(P0) + IF_INTERNAL_OFFS) + w1*(int(P1) + IF_INTERNAL_OFFS) + round + (offset << uint(shift-1))) >> uint(shift) )) );
}

func (this *TComWeightPrediction)  weightUnidirY( w0 int, P0 Pel, round, shift, offset int) Pel{
  return ClipY( Pel(( (w0*(int(P0) + IF_INTERNAL_OFFS) + round) >> uint(shift) ) + offset) );
}
func (this *TComWeightPrediction)  weightUnidirC( w0 int, P0 Pel, round, shift, offset int) Pel{
  return ClipC( Pel(( (w0*(int(P0) + IF_INTERNAL_OFFS) + round) >> uint(shift) ) + offset) );
}

func (this *TComWeightPrediction)  GetWpScaling(pcCU *TComDataCU, iRefIdx0, iRefIdx1 int, wp0, wp1 []wpScalingParam){
  pcSlice       := pcCU.GetSlice();
  pps           := pcCU.GetSlice().GetPPS();
  wpBiPred 		:= pps.GetWPBiPred();
  var pwp []wpScalingParam;
  bBiDir        := (iRefIdx0>=0 && iRefIdx1>=0);
  bUniDir       := !bBiDir;

  if bUniDir || wpBiPred { // explicit --------------------
    if iRefIdx0 >= 0 {
      wp0 = pcSlice.GetWpScaling(REF_PIC_LIST_0, iRefIdx0);
    }
    if iRefIdx1 >= 0 {
      wp1 = pcSlice.GetWpScaling(REF_PIC_LIST_1, iRefIdx1);
    }
  }else{
    //assert(0);
  }

  if iRefIdx0 < 0 {
    wp0 = nil;
  }
  if iRefIdx1 < 0 {
    wp1 = nil;
  }

  if bBiDir { // Bi-Dir case
    for yuv:=0 ; yuv<3 ; yuv++ {
      var bitDepth int;
      if yuv!=0 {
      	bitDepth = G_bitDepthC;
      }else{
        bitDepth = G_bitDepthY;
      }
      wp0[yuv].w      = wp0[yuv].iWeight;
      wp0[yuv].o      = wp0[yuv].iOffset * (1 << uint(bitDepth-8));
      wp1[yuv].w      = wp1[yuv].iWeight;
      wp1[yuv].o      = wp1[yuv].iOffset * (1 << uint(bitDepth-8));
      wp0[yuv].offset = wp0[yuv].o + wp1[yuv].o;
      wp0[yuv].shift  = int(wp0[yuv].uiLog2WeightDenom) + 1;
      wp0[yuv].round  = (1 << wp0[yuv].uiLog2WeightDenom);
      wp1[yuv].offset = wp0[yuv].offset;
      wp1[yuv].shift  = wp0[yuv].shift;
      wp1[yuv].round  = wp0[yuv].round;
    }
  }else{  // Unidir
    if iRefIdx0>=0 {
    	pwp = wp0 ;
    }else{
    	pwp = wp1 ;
    }
    for yuv:=0 ; yuv<3 ; yuv++ {
      var bitDepth int;
      if yuv!=0 {
      	bitDepth = G_bitDepthC;
      }else{
        bitDepth = G_bitDepthY;
      }
      pwp[yuv].w      = pwp[yuv].iWeight;
      pwp[yuv].offset = pwp[yuv].iOffset * (1 << uint(bitDepth-8));
      pwp[yuv].shift  = int(pwp[yuv].uiLog2WeightDenom);
      if pwp[yuv].uiLog2WeightDenom>=1 {
      	pwp[yuv].round  = (1 << (pwp[yuv].uiLog2WeightDenom-1)) ;
      }else{
      	pwp[yuv].round  = (0);
      }
    }
  }
}

func (this *TComWeightPrediction)  AddWeightBi( pcYuvSrc0, pcYuvSrc1 *TComYuv, iPartUnitIdx uint, iWidth, iHeight int, wp0, wp1 []wpScalingParam, rpcYuvDst *TComYuv, bRound bool ){
  var x, y int;

  pSrcY0  := pcYuvSrc0.GetLumaAddr1( iPartUnitIdx );
  pSrcU0  := pcYuvSrc0.GetCbAddr1  ( iPartUnitIdx );
  pSrcV0  := pcYuvSrc0.GetCrAddr1  ( iPartUnitIdx );
  
  pSrcY1  := pcYuvSrc1.GetLumaAddr1( iPartUnitIdx );
  pSrcU1  := pcYuvSrc1.GetCbAddr1  ( iPartUnitIdx );
  pSrcV1  := pcYuvSrc1.GetCrAddr1  ( iPartUnitIdx );
  
  pDstY   := rpcYuvDst.GetLumaAddr1( iPartUnitIdx );
  pDstU   := rpcYuvDst.GetCbAddr1  ( iPartUnitIdx );
  pDstV   := rpcYuvDst.GetCrAddr1  ( iPartUnitIdx );
  
  // Luma : --------------------------------------------
  w0       := wp0[0].w;
  offset   := wp0[0].offset;
  shiftNum := IF_INTERNAL_PREC - G_bitDepthY;
  shift    := wp0[0].shift + shiftNum;
  round    := 0;
  if shift!=0{
  	round  = (1<<uint(shift-1)) * int(B2U(bRound));
  }
  w1       := wp1[0].w;

  iSrc0Stride := int(pcYuvSrc0.GetStride());
  iSrc1Stride := int(pcYuvSrc1.GetStride());
  iDstStride  := int(rpcYuvDst.GetStride());
  for y = iHeight-1; y >= 0; y-- {
    for x = iWidth-1; x >= 0; x-=4{
      // note: luma min width is 4
      pDstY[y*iDstStride+x-0] = this.weightBidirY(w0,pSrcY0[y*iSrc0Stride+x-0], w1,pSrcY1[y*iSrc1Stride+x-0], round, shift, offset); //x--;
      pDstY[y*iDstStride+x-1] = this.weightBidirY(w0,pSrcY0[y*iSrc0Stride+x-1], w1,pSrcY1[y*iSrc1Stride+x-1], round, shift, offset); //x--;
      pDstY[y*iDstStride+x-2] = this.weightBidirY(w0,pSrcY0[y*iSrc0Stride+x-2], w1,pSrcY1[y*iSrc1Stride+x-2], round, shift, offset); //x--;
      pDstY[y*iDstStride+x-3] = this.weightBidirY(w0,pSrcY0[y*iSrc0Stride+x-3], w1,pSrcY1[y*iSrc1Stride+x-3], round, shift, offset); //x--;
    }
    //pSrcY0 += iSrc0Stride;
    //pSrcY1 += iSrc1Stride;
    //pDstY  += iDstStride;
  }

  
  // Chroma U : --------------------------------------------
  w0      = wp0[1].w;
  offset  = wp0[1].offset;
  shiftNum = IF_INTERNAL_PREC - G_bitDepthC;
  shift   = wp0[1].shift + shiftNum;
  if shift!=0{
  	round   = (1<<uint(shift-1));
  }else{
  	round   = 0;
  }
  w1      = wp1[1].w;

  iSrc0Stride = int(pcYuvSrc0.GetCStride());
  iSrc1Stride = int(pcYuvSrc1.GetCStride());
  iDstStride  = int(rpcYuvDst.GetCStride());
  
  iWidth  >>=1;
  iHeight >>=1;
  
  for y = iHeight-1; y >= 0; y-- {
    for x = iWidth-1; x >= 0; x-=2{
      // note: chroma min width is 2
      pDstU[y*iDstStride+x-0] = this.weightBidirC(w0,pSrcU0[y*iSrc0Stride+x-0], w1,pSrcU1[y*iSrc1Stride+x-0], round, shift, offset); //x--;
      pDstU[y*iDstStride+x-1] = this.weightBidirC(w0,pSrcU0[y*iSrc0Stride+x-1], w1,pSrcU1[y*iSrc1Stride+x-1], round, shift, offset); //x--;
    }
    //pSrcU0 += iSrc0Stride;
    //pSrcU1 += iSrc1Stride;
    //pDstU  += iDstStride;
  }

  // Chroma V : --------------------------------------------
  w0      = wp0[2].w;
  offset  = wp0[2].offset;
  shift   = wp0[2].shift + shiftNum;
  if shift!=0{
  	round   = (1<<uint(shift-1));
  }else{
  	round   = 0;
  }
  w1      = wp1[2].w;

  for y = iHeight-1; y >= 0; y-- {
    for x = iWidth-1; x >= 0; x-=2{
      // note: chroma min width is 2
      pDstV[y*iDstStride+x-0] = this.weightBidirC(w0,pSrcV0[y*iSrc0Stride+x-0], w1,pSrcV1[y*iSrc1Stride+x-0], round, shift, offset); //x--;
      pDstV[y*iDstStride+x-1] = this.weightBidirC(w0,pSrcV0[y*iSrc0Stride+x-1], w1,pSrcV1[y*iSrc1Stride+x-1], round, shift, offset); //x--;
    }
    //pSrcV0 += iSrc0Stride;
    //pSrcV1 += iSrc1Stride;
    //pDstV  += iDstStride;
  }
}
func (this *TComWeightPrediction)  AddWeightUni( pcYuvSrc0 *TComYuv, iPartUnitIdx uint, iWidth, iHeight int, wp0 []wpScalingParam, rpcYuvDst *TComYuv){
  var x, y int;
  
  pSrcY0  := pcYuvSrc0.GetLumaAddr1( iPartUnitIdx );
  pSrcU0  := pcYuvSrc0.GetCbAddr1  ( iPartUnitIdx );
  pSrcV0  := pcYuvSrc0.GetCrAddr1  ( iPartUnitIdx );
  
  pDstY   := rpcYuvDst.GetLumaAddr1( iPartUnitIdx );
  pDstU   := rpcYuvDst.GetCbAddr1  ( iPartUnitIdx );
  pDstV   := rpcYuvDst.GetCrAddr1  ( iPartUnitIdx );
  
  // Luma : --------------------------------------------
  w0       := wp0[0].w;
  offset   := wp0[0].offset;
  shiftNum := IF_INTERNAL_PREC - G_bitDepthY;
  shift    := wp0[0].shift + shiftNum;
  round    := 0;
  if shift!=0{
  	round   = (1<<uint(shift-1));
  }
  
  iSrc0Stride := int(pcYuvSrc0.GetStride());
  iDstStride  := int(rpcYuvDst.GetStride());
  
  for y = iHeight-1; y >= 0; y-- {
    for x = iWidth-1; x >= 0; x-=4{
      // note: luma min width is 4
      pDstY[y*iDstStride+x-0] = this.weightUnidirY(w0,pSrcY0[y*iSrc0Stride+x-0], round, shift, offset); //x--;
      pDstY[y*iDstStride+x-1] = this.weightUnidirY(w0,pSrcY0[y*iSrc0Stride+x-1], round, shift, offset); //x--;
      pDstY[y*iDstStride+x-2] = this.weightUnidirY(w0,pSrcY0[y*iSrc0Stride+x-2], round, shift, offset); //x--;
      pDstY[y*iDstStride+x-3] = this.weightUnidirY(w0,pSrcY0[y*iSrc0Stride+x-3], round, shift, offset); //x--;
    }
    //pSrcY0 += iSrc0Stride;
    //pDstY  += iDstStride;
  }
  
  // Chroma U : --------------------------------------------
  w0      = wp0[1].w;
  offset  = wp0[1].offset;
  shiftNum = IF_INTERNAL_PREC - G_bitDepthC;
  shift   = wp0[1].shift + shiftNum;
  if shift!=0{
  	round   = (1<<uint(shift-1));
  }else{
  	round   = 0;
  }

  iSrc0Stride = int(pcYuvSrc0.GetCStride());
  iDstStride  = int(rpcYuvDst.GetCStride());
  
  iWidth  >>=1;
  iHeight >>=1;
  
  for y = iHeight-1; y >= 0; y-- {
    for x = iWidth-1; x >= 0; x-=2{
      // note: chroma min width is 2
      pDstU[y*iDstStride+x-0] = this.weightUnidirC(w0,pSrcU0[y*iSrc0Stride+x-0], round, shift, offset); //x--;
      pDstU[y*iDstStride+x-1] = this.weightUnidirC(w0,pSrcU0[y*iSrc0Stride+x-1], round, shift, offset); //x--;
    }
    //pSrcU0 += iSrc0Stride;
    //pDstU  += iDstStride;
  }

  // Chroma V : --------------------------------------------
  w0      = wp0[2].w;
  offset  = wp0[2].offset;
  shift   = wp0[2].shift + shiftNum;
  if shift!=0{
  	round   = (1<<uint(shift-1));
  }else{
  	round   = 0;
  }

  for y = iHeight-1; y >= 0; y-- {
    for x = iWidth-1; x >= 0; x-=2{
      // note: chroma min width is 2
      pDstV[y*iDstStride+x-0] = this.weightUnidirC(w0,pSrcV0[y*iSrc0Stride+x-0], round, shift, offset); //x--;
      pDstV[y*iDstStride+x-1] = this.weightUnidirC(w0,pSrcV0[y*iSrc0Stride+x-1], round, shift, offset); //x--;
    }
    //pSrcV0 += iSrc0Stride;
    //pDstV  += iDstStride;
  }
}

func (this *TComWeightPrediction)  xWeightedPredictionUni( pcCU *TComDataCU, pcYuvSrc *TComYuv, uiPartAddr uint, iWidth, iHeight int, eRefPicList RefPicList, rpcYuvPred *TComYuv, iPartIdx, iRefIdx int ){
  var  pwp, pwpTmp []wpScalingParam;
  if iRefIdx < 0 {
    iRefIdx = int(pcCU.GetCUMvField( eRefPicList ).GetRefIdx( int(uiPartAddr) ));
  }
  //assert (iRefIdx >= 0);

  if eRefPicList == REF_PIC_LIST_0 {
    this.GetWpScaling(pcCU, iRefIdx, -1, pwp, pwpTmp);
  }else{
    this.GetWpScaling(pcCU, -1, iRefIdx, pwpTmp, pwp);
  }
  this.AddWeightUni( pcYuvSrc, uiPartAddr, iWidth, iHeight, pwp, rpcYuvPred );
}
func (this *TComWeightPrediction)  xWeightedPredictionBi( pcCU *TComDataCU, pcYuvSrc0, pcYuvSrc1 *TComYuv, iRefIdx0, iRefIdx1 int, uiPartIdx uint, iWidth, iHeight int, rpcYuvDst *TComYuv){
  var pwp0, pwp1 []wpScalingParam;
  //pps := pcCU.GetSlice().GetPPS();
  //assert( pps.GetWPBiPred());

  this.GetWpScaling(pcCU, iRefIdx0, iRefIdx1, pwp0, pwp1);

  if iRefIdx0 >= 0 && iRefIdx1 >= 0 {
    this.AddWeightBi(pcYuvSrc0, pcYuvSrc1, uiPartIdx, iWidth, iHeight, pwp0, pwp1, rpcYuvDst, true );
  }else if iRefIdx0 >= 0 && iRefIdx1 <  0 {
    this.AddWeightUni( pcYuvSrc0, uiPartIdx, iWidth, iHeight, pwp0, rpcYuvDst );
  }else if iRefIdx0 <  0 && iRefIdx1 >= 0 {
    this.AddWeightUni( pcYuvSrc1, uiPartIdx, iWidth, iHeight, pwp1, rpcYuvDst );
  }else{
    //assert (0);
  }
}

/// prediction class
type TComPrediction struct {
    TComWeightPrediction
    //protected:
    m_piYuvExt     []int
    m_iYuvExtStride int
    m_iYuvExtHeight int

    m_acYuvPred        [2]TComYuv
    m_cYuvPredTemp     TComYuv
    m_filteredBlock    [4][4]TComYuv
    m_filteredBlockTmp [4]TComYuv

    m_if TComInterpolationFilter

    m_pLumaRecBuffer []Pel ///< array for downsampled reconstructed luma sample
    m_iLumaRecStride int  ///< stride of #m_pLumaRecBuffer array
}


func NewTComPrediction() *TComPrediction{
	return &TComPrediction{m_iLumaRecStride:0};
}

func (this *TComPrediction) InitTempBuff(){
  if this.m_piYuvExt == nil {
    extWidth  := G_uiMaxCUWidth + 16; 
    extHeight := G_uiMaxCUHeight + 1;
    var i, j int;
    for i = 0; i < 4; i++ {
      this.m_filteredBlockTmp[i].Create(extWidth, extHeight + 7);
      for j = 0; j < 4; j++{
        this.m_filteredBlock[i][j].Create(extWidth, extHeight);
      }
    }
    this.m_iYuvExtHeight = int((G_uiMaxCUHeight + 2) << 4);
    this.m_iYuvExtStride = int((G_uiMaxCUWidth  + 8) << 4);
    this.m_piYuvExt = make([]int, this.m_iYuvExtStride * this.m_iYuvExtHeight );

    // new structure
    this.m_acYuvPred[0].Create( G_uiMaxCUWidth, G_uiMaxCUHeight );
    this.m_acYuvPred[1].Create( G_uiMaxCUWidth, G_uiMaxCUHeight );

    this.m_cYuvPredTemp.Create( G_uiMaxCUWidth, G_uiMaxCUHeight );
  }

  if this.m_iLumaRecStride != int(G_uiMaxCUWidth>>1) + 1 {
    this.m_iLumaRecStride =  int(G_uiMaxCUWidth>>1) + 1;
    this.m_pLumaRecBuffer = make([]Pel, this.m_iLumaRecStride * this.m_iLumaRecStride );
  }
}


func (this *TComPrediction)  xPredIntraAng            ( bitDepth int, pSrc2 []int, srcStride int, rpDst []Pel, dstStride int, 
													    width, height, dirMode uint, blkAboveAvailable, blkLeftAvailable, bFilter bool){
  var k,l int;
  blkSize := int(width);
  //pSrc    := pSrc2[srcStride+1:];
  pDst    := rpDst;

  // Map the mode index to main prediction direction and angle
  //assert( dirMode > 0 ); //no planar
  modeDC        := dirMode < 2;
  modeHor       := !modeDC && (dirMode < 18);
  modeVer       := !modeDC && !modeHor;
  var intraPredAngle int;
  if modeVer {
  	intraPredAngle = int(dirMode) - VER_IDX;
  }else if modeHor {
  	intraPredAngle = -(int(dirMode) - HOR_IDX) 
  }else{
  	intraPredAngle = 0;
  }
  absAng         := ABS(intraPredAngle).(int);
  var signAng int;
  if intraPredAngle < 0 {
  	signAng = -1;
  }else{
  	signAng = 1;
  }
  // Set bitshifts and scale the angle parameter to block size
  var angTable =[9]int{0,    2,    5,   9,  13,  17,  21,  26,  32};
  var invAngTable=[9]int{0, 4096, 1638, 910, 630, 482, 390, 315, 256}; // (256 * 32) / Angle
  invAngle       := invAngTable[absAng];
  absAng             = angTable[absAng];
  intraPredAngle     = signAng * absAng;

  // Do the DC prediction
  if modeDC {
    dcval := this.PredIntraGetPredValDC(pSrc2, srcStride, width, height, blkAboveAvailable, blkLeftAvailable);

    for k=0;k<blkSize;k++ {
      for l=0;l<blkSize;l++ {
        pDst[k*dstStride+l] = dcval;
      }
    }
  }else {// Do angular predictions
    var refMain []Pel;
    var refSide []Pel;
    var refAbove	[2*MAX_CU_SIZE+1]Pel;
    var refLeft		[2*MAX_CU_SIZE+1]Pel;

    // Initialise the Main and Left reference array.
    if intraPredAngle < 0 {
      for k=0;k<blkSize+1;k++ {
        refAbove[k+blkSize-1] = Pel(pSrc2[srcStride+1+k-srcStride-1]);//Pel(pSrc[k-srcStride-1]);
      }
      for k=0;k<blkSize+1;k++ {
        refLeft[k+blkSize-1]  = Pel(pSrc2[srcStride+1+(k-1)*srcStride-1]);//Pel(pSrc[(k-1)*srcStride-1]);
      }
      if modeVer{
      	refMain = refAbove [ (blkSize-1):];
      	refSide = refLeft  [ (blkSize-1):];
	  }else{
	  	refMain = refLeft  [ (blkSize-1):];
      	refSide = refAbove [ (blkSize-1):];
	  }
      // Extend the Main reference to the left.
      invAngleSum := 128;       // rounding for (shift by 8)
      for k=-1; k>blkSize*intraPredAngle>>5; k-- {
        invAngleSum += invAngle;
        refMain[k] = refSide[invAngleSum>>8];
      }
    }else{
      for k=0;k<2*blkSize+1;k++{
        refAbove[k] = Pel(pSrc2[srcStride+1+k-srcStride-1]);//Pel(pSrc[k-srcStride-1]);
      }
      for k=0;k<2*blkSize+1;k++{
        refLeft[k]  = Pel(pSrc2[srcStride+1+(k-1)*srcStride-1]);//Pel(pSrc[(k-1)*srcStride-1]);
      }
      
      if modeVer{
      	refMain = refAbove [:];
      	refSide = refLeft  [:];
	  }else{
	  	refMain = refLeft  [:];
      	refSide = refAbove [:];
	  }

    }

    if intraPredAngle == 0{
      for k=0;k<blkSize;k++ {
        for l=0;l<blkSize;l++{
          pDst[k*dstStride+l] = refMain[l+1];
        }
      }

      if bFilter {
        for k=0;k<blkSize;k++ {
          pDst[k*dstStride] = CLIP3(0, (1<<uint(bitDepth))-1, pDst[k*dstStride] + (( refSide[k+1] - refSide[0] ) >> 1) ).(Pel);
        }
      }
    }else{
      deltaPos:=0;
      var deltaInt,deltaFract,refMainIndex int;

      for k=0;k<blkSize;k++ {
        deltaPos += intraPredAngle;
        deltaInt   = deltaPos >> 5;
        deltaFract = deltaPos & (32 - 1);

        if deltaFract!=0{
          // Do linear filtering
          for l=0;l<blkSize;l++ {
            refMainIndex        = l+deltaInt+1;
            pDst[k*dstStride+l] = Pel( ((32-deltaFract)*int(refMain[refMainIndex])+deltaFract*int(refMain[refMainIndex+1])+16) >> 5 );
          }
        }else{
          // Just copy the integer samples
          for l=0;l<blkSize;l++ {
            pDst[k*dstStride+l] = refMain[l+deltaInt+1];
          }
        }
      }
    }

    // Flip the block if this is the horizontal mode
    if modeHor{
      var tmp Pel;
      for k=0;k<blkSize-1;k++ {
        for l=k+1;l<blkSize;l++ {
          tmp                 = pDst[k*dstStride+l];
          pDst[k*dstStride+l] = pDst[l*dstStride+k];
          pDst[l*dstStride+k] = tmp;
        }
      }
    }
  }
}
func (this *TComPrediction)  xPredIntraPlanar         ( pSrc []int, srcStride int, rpDst []Pel, dstStride int, width,  height uint){
  //assert(width == height);

  var k, l, bottomLeft, topRight int;
  var horPred int;
  var leftColumn, topRow, bottomRow, rightColumn	[MAX_CU_SIZE]int;
  blkSize := int(width);
  offset2D := width;
  shift1D := uint(G_aucConvertToBit[ width ]) + 2;
  shift2D := shift1D + 1;

  // Get left and above reference column and row
  for k=0;k<blkSize+1;k++ {
    topRow[k] = pSrc[k-srcStride];
    leftColumn[k] = pSrc[k*srcStride-1];
  }

  // Prepare intermediate variables used in interpolation
  bottomLeft = leftColumn[blkSize];
  topRight   = topRow[blkSize];
  for k=0;k<blkSize;k++ {
    bottomRow[k]   = bottomLeft - topRow[k];
    rightColumn[k] = topRight   - leftColumn[k];
    topRow[k]      <<= shift1D;
    leftColumn[k]  <<= shift1D;
  }

  // Generate prediction signal
  for k=0;k<blkSize;k++ {
    horPred = leftColumn[k] + int(offset2D);
    for l=0;l<blkSize;l++ {
      horPred += rightColumn[k];
      topRow[l] += bottomRow[l];
      rpDst[k*dstStride+l] = Pel( (horPred + topRow[l]) >> shift2D );
    }
  }
}

  // motion compensation functions
func (this *TComPrediction)  xPredInterUni            ( pcCU *TComDataCU, uiPartAddr uint, iWidth, iHeight int, eRefPicList RefPicList, rpcYuvPred *TComYuv, iPartIdx int, bi bool ){
  iRefIdx     := pcCU.GetCUMvField( eRefPicList ).GetRefIdx( int(uiPartAddr) );           //assert (iRefIdx >= 0);
  cMv         := pcCU.GetCUMvField( eRefPicList ).GetMv( int(uiPartAddr) );
  pcCU.ClipMv(cMv);
  this.xPredInterLumaBlk  ( pcCU, pcCU.GetSlice().GetRefPic( eRefPicList, int(iRefIdx) ).GetPicYuvRec(), uiPartAddr, cMv, iWidth, iHeight, rpcYuvPred, bi );
  this.xPredInterChromaBlk( pcCU, pcCU.GetSlice().GetRefPic( eRefPicList, int(iRefIdx) ).GetPicYuvRec(), uiPartAddr, cMv, iWidth, iHeight, rpcYuvPred, bi );
}
func (this *TComPrediction)  xPredInterBi             ( pcCU *TComDataCU, uiPartAddr uint, iWidth, iHeight int, rpcYuvPred *TComYuv, iPartIdx int){
  var pcMbYuv *TComYuv;
  var iRefIdx = [2]int{-1, -1};

  for iRefList := 0; iRefList < 2; iRefList++ {
    var eRefPicList RefPicList;
    if iRefList!=0 {
    	eRefPicList = REF_PIC_LIST_1 ;
    }else{
    	eRefPicList = REF_PIC_LIST_0;
    }
    
    iRefIdx[iRefList] = int(pcCU.GetCUMvField( eRefPicList ).GetRefIdx( int(uiPartAddr) ));

    if iRefIdx[iRefList] < 0 {
      continue;
    }

    //assert( iRefIdx[iRefList] < pcCU.GetSlice().GetNumRefIdx(eRefPicList) );

    pcMbYuv = &this.m_acYuvPred[iRefList];
    if pcCU.GetCUMvField( REF_PIC_LIST_0 ).GetRefIdx( int(uiPartAddr) ) >= 0 && pcCU.GetCUMvField( REF_PIC_LIST_1 ).GetRefIdx( int(uiPartAddr) ) >= 0 {
      this.xPredInterUni ( pcCU, uiPartAddr, iWidth, iHeight, eRefPicList, pcMbYuv, iPartIdx, true );
    }else{
      if   ( pcCU.GetSlice().GetPPS().GetUseWP()    && pcCU.GetSlice().GetSliceType() == P_SLICE ) || 
           ( pcCU.GetSlice().GetPPS().GetWPBiPred() && pcCU.GetSlice().GetSliceType() == B_SLICE ) {
        this.xPredInterUni ( pcCU, uiPartAddr, iWidth, iHeight, eRefPicList, pcMbYuv, iPartIdx, true );
      }else{
        this.xPredInterUni ( pcCU, uiPartAddr, iWidth, iHeight, eRefPicList, pcMbYuv, iPartIdx, false );
      }
    }
  }

  if pcCU.GetSlice().GetPPS().GetWPBiPred() && pcCU.GetSlice().GetSliceType() == B_SLICE {
    this.xWeightedPredictionBi( pcCU, &this.m_acYuvPred[0], &this.m_acYuvPred[1], iRefIdx[0], iRefIdx[1], uiPartAddr, iWidth, iHeight, rpcYuvPred );
  }else if pcCU.GetSlice().GetPPS().GetUseWP() && pcCU.GetSlice().GetSliceType() == P_SLICE {
    this.xWeightedPredictionUni( pcCU, &this.m_acYuvPred[0], uiPartAddr, iWidth, iHeight, REF_PIC_LIST_0, rpcYuvPred, iPartIdx, -1 ); 
  }else{
    this.xWeightedAverage( pcCU, &this.m_acYuvPred[0], &this.m_acYuvPred[1], iRefIdx[0], iRefIdx[1], uiPartAddr, uint(iWidth), uint(iHeight), rpcYuvPred );
  }
}
func (this *TComPrediction)  xPredInterLumaBlk  ( cu *TComDataCU, refPic *TComPicYuv, partAddr uint, mv *TComMv, width, height int, dstPic *TComYuv, bi bool){
  refStride := refPic.GetStride();  
  refOffset := int( mv.GetHor() >> 2 ) + int( mv.GetVer() >> 2 ) * refStride;
  ref      := refPic.GetLumaAddr2( int(cu.GetAddr()), int(cu.GetZorderIdxInCU() + partAddr) )[ refOffset:];
  
  dstStride := int(dstPic.GetStride());
  dst       := dstPic.GetLumaAddr1( partAddr );
  
  xFrac := int(mv.GetHor() & 0x3);
  yFrac := int(mv.GetVer() & 0x3);

  if yFrac == 0 {
    this.m_if.FilterHorLuma( ref, refStride, dst, dstStride, width, height, xFrac,       !bi );
  }else if xFrac == 0 {
    this.m_if.FilterVerLuma( ref, refStride, dst, dstStride, width, height, yFrac, true, !bi );
  }else {
    tmpStride := int(this.m_filteredBlockTmp[0].GetStride());
    tmp    := this.m_filteredBlockTmp[0].GetLumaAddr();

    filterSize := NTAPS_LUMA;
    halfFilterSize := ( filterSize >> 1 );

    this.m_if.FilterHorLuma(ref [- (halfFilterSize-1)*refStride:], refStride, tmp, tmpStride, width, height+filterSize-1, xFrac, false     );
    this.m_if.FilterVerLuma(tmp [ (halfFilterSize-1)*tmpStride:], tmpStride, dst, dstStride, width, height,              yFrac, false, !bi);    
  }
}
func (this *TComPrediction)  xPredInterChromaBlk( cu *TComDataCU, refPic *TComPicYuv, partAddr uint, mv *TComMv, width, height int, dstPic *TComYuv, bi bool){
  refStride  := refPic.GetCStride();
  dstStride  := int(dstPic.GetCStride());
  
  refOffset  := int(mv.GetHor() >> 3) + int(mv.GetVer() >> 3) * refStride;
  
  refCb     := refPic.GetCbAddr2( int(cu.GetAddr()), int(cu.GetZorderIdxInCU() + partAddr) ) [refOffset:];
  refCr     := refPic.GetCrAddr2( int(cu.GetAddr()), int(cu.GetZorderIdxInCU() + partAddr) ) [refOffset:];
  
  dstCb := dstPic.GetCbAddr1( partAddr );
  dstCr := dstPic.GetCrAddr1( partAddr );
  
  xFrac  := int(mv.GetHor() & 0x7);
  yFrac  := int(mv.GetVer() & 0x7);
  cxWidth  := width  >> 1;
  cxHeight := height >> 1;
  
  extStride := int(this.m_filteredBlockTmp[0].GetStride());
  extY      := this.m_filteredBlockTmp[0].GetLumaAddr();
  
  filterSize := NTAPS_CHROMA;
  
  halfFilterSize := (filterSize>>1);
  
  if yFrac == 0 {
    this.m_if.FilterHorChroma(refCb, refStride, dstCb,  dstStride, cxWidth, cxHeight, xFrac, !bi);    
    this.m_if.FilterHorChroma(refCr, refStride, dstCr,  dstStride, cxWidth, cxHeight, xFrac, !bi);    
  }else if xFrac == 0 {
    this.m_if.FilterVerChroma(refCb, refStride, dstCb, dstStride, cxWidth, cxHeight, yFrac, true, !bi);    
    this.m_if.FilterVerChroma(refCr, refStride, dstCr, dstStride, cxWidth, cxHeight, yFrac, true, !bi);    
  }else{
    this.m_if.FilterHorChroma(refCb [- (halfFilterSize-1)*refStride:], refStride, extY,  extStride, cxWidth, cxHeight+filterSize-1, xFrac, false);
    this.m_if.FilterVerChroma(extY  [  (halfFilterSize-1)*extStride:], extStride, dstCb, dstStride, cxWidth, cxHeight  , yFrac, false, !bi);
    
    this.m_if.FilterHorChroma(refCr [- (halfFilterSize-1)*refStride:], refStride, extY,  extStride, cxWidth, cxHeight+filterSize-1, xFrac, false);
    this.m_if.FilterVerChroma(extY  [  (halfFilterSize-1)*extStride:], extStride, dstCr, dstStride, cxWidth, cxHeight  , yFrac, false, !bi);    
  }
}
func (this *TComPrediction)  xWeightedAverage  ( pcCU *TComDataCU, pcYuvSrc0 *TComYuv, pcYuvSrc1 *TComYuv, iRefIdx0, iRefIdx1 int, uiPartIdx uint, iWidth, iHeight uint, rpcYuvDst *TComYuv){
  if iRefIdx0 >= 0 && iRefIdx1 >= 0 {
    rpcYuvDst.AddAvg( pcYuvSrc0, pcYuvSrc1, uiPartIdx, iWidth, iHeight );
  }else if iRefIdx0 >= 0 && iRefIdx1 <  0 {
    pcYuvSrc0.CopyPartToPartYuv( rpcYuvDst, uiPartIdx, iWidth, iHeight );
  }else if iRefIdx0 <  0 && iRefIdx1 >= 0 {
    pcYuvSrc1.CopyPartToPartYuv( rpcYuvDst, uiPartIdx, iWidth, iHeight );
  }
}

func (this *TComPrediction)  xDCPredFiltering( pSrc2 []int, iSrcStride int, rpDst []Pel, iDstStride, iWidth, iHeight int){
  pSrc := pSrc2[iSrcStride+1:]; //ptrSrc[sw+1:]
  pDst := rpDst;
  var x, y, iDstStride2, iSrcStride2 int;

  // boundary pixels processing
  //pDst[0] = Pel((pSrc[-iSrcStride] + pSrc[-1] + 2 * int(pDst[0]) + 2) >> 2);
  pDst[0] = Pel((pSrc2[iSrcStride+1-iSrcStride] + pSrc2[iSrcStride+1-1] + 2 * int(pDst[0]) + 2) >> 2);


  for x = 1; x < iWidth; x++ {
    //pDst[x] = Pel((pSrc[x - iSrcStride] +  3 * int(pDst[x]) + 2) >> 2);
    pDst[x] = Pel((pSrc2[x +iSrcStride+1- iSrcStride] +  3 * int(pDst[x]) + 2) >> 2);
  }

  iDstStride2 = iDstStride;
  iSrcStride2 = iSrcStride-1;
  for y = 1; y < iHeight; y++ {
    pDst[iDstStride2] = Pel((pSrc[iSrcStride2] + 3 * int(pDst[iDstStride2]) + 2) >> 2);
    iDstStride2+=iDstStride; 
    iSrcStride2+=iSrcStride;
  }

  return;
}
func (this *TComPrediction)  xCheckIdenticalMotion    ( pcCU *TComDataCU, PartAddr uint) bool{
  if pcCU.GetSlice().IsInterB() && !pcCU.GetSlice().GetPPS().GetWPBiPred() {
    if pcCU.GetCUMvField(REF_PIC_LIST_0).GetRefIdx(int(PartAddr)) >= 0 && pcCU.GetCUMvField(REF_PIC_LIST_1).GetRefIdx(int(PartAddr)) >= 0 {
      RefPOCL0 := pcCU.GetSlice().GetRefPic(REF_PIC_LIST_0, int(pcCU.GetCUMvField(REF_PIC_LIST_0).GetRefIdx(int(PartAddr)))).GetPOC();
      RefPOCL1 := pcCU.GetSlice().GetRefPic(REF_PIC_LIST_1, int(pcCU.GetCUMvField(REF_PIC_LIST_1).GetRefIdx(int(PartAddr)))).GetPOC();
      if RefPOCL0 == RefPOCL1 && pcCU.GetCUMvField(REF_PIC_LIST_0).GetMv(int(PartAddr)) == pcCU.GetCUMvField(REF_PIC_LIST_1).GetMv(int(PartAddr)) {
        return true;
      }
    }
  }
  return false;
}

  // inter
func (this *TComPrediction) MotionCompensation         ( pcCU *TComDataCU, pcYuvPred *TComYuv,  eRefPicList RefPicList,  iPartIdx int ){
  var        iWidth, iHeight int;
  var        uiPartAddr uint;

  if iPartIdx >= 0 {
    pcCU.GetPartIndexAndSize( uint(iPartIdx), &uiPartAddr, &iWidth, &iHeight );
    if eRefPicList != REF_PIC_LIST_X {
      if pcCU.GetSlice().GetPPS().GetUseWP() {
        this.xPredInterUni (pcCU, uiPartAddr, iWidth, iHeight, eRefPicList, pcYuvPred, iPartIdx, true );
      }else{
        this.xPredInterUni (pcCU, uiPartAddr, iWidth, iHeight, eRefPicList, pcYuvPred, iPartIdx, false );
      }
      if pcCU.GetSlice().GetPPS().GetUseWP() {
        this.xWeightedPredictionUni( pcCU, pcYuvPred, uiPartAddr, iWidth, iHeight, eRefPicList, pcYuvPred, iPartIdx, -1 );
      }
    }else{
      if this.xCheckIdenticalMotion( pcCU, uiPartAddr ) {
        this.xPredInterUni (pcCU, uiPartAddr, iWidth, iHeight, REF_PIC_LIST_0, pcYuvPred, iPartIdx, false );
      }else{
        this.xPredInterBi  (pcCU, uiPartAddr, iWidth, iHeight, pcYuvPred, iPartIdx );
      }
    }
    return;
  }

  for iPartIdx = 0; iPartIdx < int(pcCU.GetNumPartInter()); iPartIdx++ {
    pcCU.GetPartIndexAndSize( uint(iPartIdx), &uiPartAddr, &iWidth, &iHeight );

    if eRefPicList != REF_PIC_LIST_X {
      if pcCU.GetSlice().GetPPS().GetUseWP() {
        this.xPredInterUni (pcCU, uiPartAddr, iWidth, iHeight, eRefPicList, pcYuvPred, iPartIdx, true );
      }else{
        this.xPredInterUni (pcCU, uiPartAddr, iWidth, iHeight, eRefPicList, pcYuvPred, iPartIdx, false );
      }
      if pcCU.GetSlice().GetPPS().GetUseWP() {
        this.xWeightedPredictionUni( pcCU, pcYuvPred, uiPartAddr, iWidth, iHeight, eRefPicList, pcYuvPred, iPartIdx, -1 );
      }
    }else{
      if this.xCheckIdenticalMotion( pcCU, uiPartAddr ) {
        this.xPredInterUni (pcCU, uiPartAddr, iWidth, iHeight, REF_PIC_LIST_0, pcYuvPred, iPartIdx, false );
      }else{
        this.xPredInterBi  (pcCU, uiPartAddr, iWidth, iHeight, pcYuvPred, iPartIdx );
      }
    }
  }
  return;
}

  // motion vector prediction
func (this *TComPrediction) GetMvPredAMVP              ( pcCU *TComDataCU,  uiPartIdx,  uiPartAddr uint,  eRefPicList RefPicList,  iRefIdx int, rcMvPred *TComMv ){
  pcAMVPInfo := pcCU.GetCUMvField(eRefPicList).GetAMVPInfo();
  if pcAMVPInfo.IN <= 1 {
    rcMvPred = &pcAMVPInfo.MvCand[0];

    pcCU.SetMVPIdxSubParts( 0, eRefPicList, uiPartAddr, uiPartIdx, uint(pcCU.GetDepth1(uiPartAddr)));
    pcCU.SetMVPNumSubParts( pcAMVPInfo.IN, eRefPicList, uiPartAddr, uiPartIdx, uint(pcCU.GetDepth1(uiPartAddr)));
    return;
  }

  //assert(pcCU.GetMVPIdx(eRefPicList,uiPartAddr) >= 0);
  rcMvPred = &pcAMVPInfo.MvCand[pcCU.GetMVPIdx2(eRefPicList,uiPartAddr)];
  return;
}

  // Angular Intra
func (this *TComPrediction) PredIntraLumaAng           ( pcTComPattern *TComPattern,  uiDirMode uint, piPred []Pel,  uiStride uint,  iWidth,  iHeight int,  pcCU *TComDataCU,  bAbove,  bLeft bool){
  pDst := piPred;
  var ptrSrc []int;

  //assert( G_aucConvertToBit[ iWidth ] >= 0 ); //   4x  4
  //assert( G_aucConvertToBit[ iWidth ] <= 5 ); // 128x128
  //assert( iWidth == iHeight  );

  ptrSrc = pcTComPattern.GetPredictorPtr( uiDirMode, uint(G_aucConvertToBit[ iWidth ]) + 2, this.m_piYuvExt );

  // get starting pixel in block
  sw := 2 * iWidth + 1;

  // Create the prediction
  if uiDirMode == PLANAR_IDX {
    this.xPredIntraPlanar( ptrSrc[sw+1:], sw, pDst, int(uiStride), uint(iWidth), uint(iHeight) );
  }else{
//#if RESTRICT_INTRA_BOUNDARY_SMOOTHING
    if (iWidth > 16) || (iHeight > 16) {
      //this.xPredIntraAng(G_bitDepthY, ptrSrc[sw+1:], sw, pDst, int(uiStride), uint(iWidth), uint(iHeight), uiDirMode, bAbove, bLeft, false );
	  this.xPredIntraAng(G_bitDepthY, ptrSrc, sw, pDst, int(uiStride), uint(iWidth), uint(iHeight), uiDirMode, bAbove, bLeft, false );
    }else{
//#endif
      //this.xPredIntraAng(G_bitDepthY, ptrSrc[sw+1:], sw, pDst, int(uiStride), uint(iWidth), uint(iHeight), uiDirMode, bAbove, bLeft, true );
      this.xPredIntraAng(G_bitDepthY, ptrSrc, sw, pDst, int(uiStride), uint(iWidth), uint(iHeight), uiDirMode, bAbove, bLeft, true );
      		
      if (uiDirMode == DC_IDX ) && bAbove && bLeft {
        //this.xDCPredFiltering( ptrSrc[sw+1:], sw, pDst, int(uiStride), iWidth, iHeight);
        this.xDCPredFiltering( ptrSrc, sw, pDst, int(uiStride), iWidth, iHeight);
      }
//#if RESTRICT_INTRA_BOUNDARY_SMOOTHING
    }
//#endif
  }
}
func (this *TComPrediction) PredIntraChromaAng         ( pcTComPattern *TComPattern, piSrc []int,  uiDirMode uint, piPred []Pel,  uiStride uint,  iWidth,  iHeight int, pcCU *TComDataCU,  bAbove,  bLeft bool){
  pDst := piPred;
  ptrSrc := piSrc;

  // get starting pixel in block
  sw := 2 * iWidth + 1;

  if uiDirMode == PLANAR_IDX {
    this.xPredIntraPlanar( ptrSrc[sw+1:], sw, pDst, int(uiStride), uint(iWidth), uint(iHeight) );
  }else{
    // Create the prediction
    //this.xPredIntraAng(G_bitDepthC, ptrSrc[sw+1:], sw, pDst, int(uiStride), uint(iWidth), uint(iHeight), uiDirMode, bAbove, bLeft, false );
	this.xPredIntraAng(G_bitDepthC, ptrSrc, sw, pDst, int(uiStride), uint(iWidth), uint(iHeight), uiDirMode, bAbove, bLeft, false );
  }
}

func (this *TComPrediction) PredIntraGetPredValDC      ( pSrc2 []int,  iSrcStride int,  iWidth,  iHeight uint,  bAbove,  bLeft bool) Pel{
  var iInd int;
  iSum := 0;
  var pDcVal Pel;

  if bAbove{
    for iInd = 0;iInd < int(iWidth);iInd++ {
      iSum += pSrc2[iSrcStride+1+iInd-iSrcStride];
    }
  }
  if bLeft {
    for iInd = 0;iInd < int(iHeight);iInd++ {
      iSum += pSrc2[iSrcStride+1+iInd*iSrcStride-1];
    }
  }

  if bAbove && bLeft {
    pDcVal = Pel((uint(iSum) + iWidth) / (iWidth + iHeight));
  }else if bAbove{
    pDcVal = Pel((uint(iSum) + iWidth/2) / iWidth);
  }else if bLeft{
    pDcVal = Pel((uint(iSum )+ iHeight/2) / iHeight);
  }else{
    pDcVal = Pel(pSrc2[iSrcStride+1-1]); // Default DC value already calculated and placed in the prediction array if no neighbors are available
  }
  
  return pDcVal;
}

func (this *TComPrediction) GetPredicBuf()         []int   { 
	return this.m_piYuvExt;      
}
func (this *TComPrediction) GetPredicBufWidth()     int   { 
	return this.m_iYuvExtStride; 
}
func (this *TComPrediction) GetPredicBufHeight()    int   { 
	return this.m_iYuvExtHeight; 
}


