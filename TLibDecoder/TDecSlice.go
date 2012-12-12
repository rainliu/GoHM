package TLibDecoder

import (
	//"container/list"
    "gohm/TLibCommon"
)

/// slice decoder class
type TDecSlice struct {
    //private:
    // access channel
    m_pcEntropyDecoder	*TDecEntropy;
    m_pcCuDecoder		*TDecCu;
    m_uiCurrSliceIdx uint

    m_pcBufferSbacDecoders		*TDecSbac;   ///< line to store temporary contexts, one per column of tiles.
    m_pcBufferBinCABACs			*TDecBinCabac;
    m_pcBufferLowLatSbacDecoders	*TDecSbac;   ///< dependent tiles: line to store temporary contexts, one per column of tiles.
    m_pcBufferLowLatBinCABACs		*TDecBinCabac;
    //#if DEPENDENT_SLICES
    CTXMem							map[int]*TDecSbac;//*list.List;//std::vector<TDecSbac*> 
    //#endif
}

//public:
func NewTDecSlice() *TDecSlice {
    return &TDecSlice{CTXMem:make(map[int]*TDecSbac)}
}

func (this *TDecSlice) Init(pcEntropyDecoder *TDecEntropy, pcCuDecoder *TDecCu) {
  this.m_pcEntropyDecoder  = pcEntropyDecoder;
  this.m_pcCuDecoder       = pcCuDecoder;
}

func (this *TDecSlice) Create( pcSlice *TLibCommon.TComSlice, iWidth, iHeight int, uiMaxWidth, uiMaxHeight, uiMaxDepth uint){
	//do nothing
}
func (this *TDecSlice) Destroy() {
	//do nothing
}

func (this *TDecSlice) DecompressSlice   ( pcBitstream *TLibCommon.TComInputBitstream, ppcSubstreams []*TLibCommon.TComInputBitstream,   
										   rpcPic *TLibCommon.TComPic, pcSbacDecoder *TDecSbac, pcSbacDecoders []TDecSbac){
  /*TComDataCU* pcCU;
  UInt        uiIsLast = 0;
  Int   iStartCUEncOrder = max(rpcPic->getSlice(rpcPic->getCurrSliceIdx())->getSliceCurStartCUAddr()/rpcPic->getNumPartInCU(), rpcPic->getSlice(rpcPic->getCurrSliceIdx())->getDependentSliceCurStartCUAddr()/rpcPic->getNumPartInCU());
  Int   iStartCUAddr = rpcPic->getPicSym()->getCUOrderMap(iStartCUEncOrder);

  // decoder don't need prediction & residual frame buffer
  rpcPic->setPicYuvPred( 0 );
  rpcPic->setPicYuvResi( 0 );
  
#if ENC_DEC_TRACE
  g_bJustDoIt = g_bEncDecTraceEnable;
#endif
  DTRACE_CABAC_VL( g_nSymbolCounter++ );
  DTRACE_CABAC_T( "\tPOC: " );
  DTRACE_CABAC_V( rpcPic->getPOC() );
  DTRACE_CABAC_T( "\n" );

#if ENC_DEC_TRACE
  g_bJustDoIt = g_bEncDecTraceDisable;
#endif

  UInt uiTilesAcross   = rpcPic->getPicSym()->getNumColumnsMinus1()+1;
  TComSlice*  pcSlice = rpcPic->getSlice(rpcPic->getCurrSliceIdx());
  Int  iNumSubstreams = pcSlice->getPPS()->getNumSubstreams();

  // delete decoders if already allocated in previous slice
  if (m_pcBufferSbacDecoders)
  {
    delete [] m_pcBufferSbacDecoders;
  }
  if (m_pcBufferBinCABACs) 
  {
    delete [] m_pcBufferBinCABACs;
  }
  // allocate new decoders based on tile numbaer
  m_pcBufferSbacDecoders = new TDecSbac    [uiTilesAcross];  
  m_pcBufferBinCABACs    = new TDecBinCABAC[uiTilesAcross];
  for (UInt ui = 0; ui < uiTilesAcross; ui++)
  {
    m_pcBufferSbacDecoders[ui].init(&m_pcBufferBinCABACs[ui]);
  }
  //save init. state
  for (UInt ui = 0; ui < uiTilesAcross; ui++)
  {
    m_pcBufferSbacDecoders[ui].load(pcSbacDecoder);
  }

  // free memory if already allocated in previous call
  if (m_pcBufferLowLatSbacDecoders)
  {
    delete [] m_pcBufferLowLatSbacDecoders;
  }
  if (m_pcBufferLowLatBinCABACs)
  {
    delete [] m_pcBufferLowLatBinCABACs;
  }
  m_pcBufferLowLatSbacDecoders = new TDecSbac    [uiTilesAcross];  
  m_pcBufferLowLatBinCABACs    = new TDecBinCABAC[uiTilesAcross];
  for (UInt ui = 0; ui < uiTilesAcross; ui++)
  {
    m_pcBufferLowLatSbacDecoders[ui].init(&m_pcBufferLowLatBinCABACs[ui]);
  }
  //save init. state
  for (UInt ui = 0; ui < uiTilesAcross; ui++)
  {
    m_pcBufferLowLatSbacDecoders[ui].load(pcSbacDecoder);
  }

  UInt uiWidthInLCUs  = rpcPic->getPicSym()->getFrameWidthInCU();
  //UInt uiHeightInLCUs = rpcPic->getPicSym()->getFrameHeightInCU();
  UInt uiCol=0, uiLin=0, uiSubStrm=0;

  UInt uiTileCol;
  UInt uiTileStartLCU;
  UInt uiTileLCUX;
  Int iNumSubstreamsPerTile = 1; // if independent.
#if DEPENDENT_SLICES
  Bool bAllowDependence = false;
#if REMOVE_ENTROPY_SLICES
  if( rpcPic->getSlice(rpcPic->getCurrSliceIdx())->getPPS()->getDependentSliceEnabledFlag() )
#else
  if( rpcPic->getSlice(rpcPic->getCurrSliceIdx())->getPPS()->getDependentSliceEnabledFlag()&& (!rpcPic->getSlice(rpcPic->getCurrSliceIdx())->getPPS()->getEntropySliceEnabledFlag()) )
#endif
  {
    bAllowDependence = true;
  }
  if( bAllowDependence )
  {
    if( !rpcPic->getSlice(rpcPic->getCurrSliceIdx())->isNextSlice() )
    {
      uiTileCol = 0;
      if(pcSlice->getPPS()->getEntropyCodingSyncEnabledFlag())
      {
        m_pcBufferSbacDecoders[uiTileCol].loadContexts( CTXMem[1]  );//2.LCU
      }
      pcSbacDecoder->loadContexts(CTXMem[0] ); //end of depSlice-1
      pcSbacDecoders[uiSubStrm].loadContexts(pcSbacDecoder);
    }
    else
    {
      if(pcSlice->getPPS()->getEntropyCodingSyncEnabledFlag())
      {
        CTXMem[1]->loadContexts(pcSbacDecoder);
      }
      CTXMem[0]->loadContexts(pcSbacDecoder);
    }
  }
#endif
  for( Int iCUAddr = iStartCUAddr; !uiIsLast && iCUAddr < rpcPic->getNumCUsInFrame(); iCUAddr = rpcPic->getPicSym()->xCalculateNxtCUAddr(iCUAddr) )
  {
    pcCU = rpcPic->getCU( iCUAddr );
    pcCU->initCU( rpcPic, iCUAddr );

#ifdef ENC_DEC_TRACE
    xTraceLCUHeader(pcCU, TRACE_LCU);
    xReadAeTr (iCUAddr, "lcu_address", TRACE_LCU);
    xReadAeTr (rpcPic->getPicSym()->getTileIdxMap(iCUAddr), "tile_id", TRACE_LCU);
#endif

    uiTileCol = rpcPic->getPicSym()->getTileIdxMap(iCUAddr) % (rpcPic->getPicSym()->getNumColumnsMinus1()+1); // what column of tiles are we in?
    uiTileStartLCU = rpcPic->getPicSym()->getTComTile(rpcPic->getPicSym()->getTileIdxMap(iCUAddr))->getFirstCUAddr();
    uiTileLCUX = uiTileStartLCU % uiWidthInLCUs;
    uiCol     = iCUAddr % uiWidthInLCUs;
    // The 'line' is now relative to the 1st line in the slice, not the 1st line in the picture.
    uiLin     = (iCUAddr/uiWidthInLCUs)-(iStartCUAddr/uiWidthInLCUs);
    // inherit from TR if necessary, select substream to use.
#if DEPENDENT_SLICES
    if( (pcSlice->getPPS()->getNumSubstreams() > 1) || ( bAllowDependence  && (uiCol == uiTileLCUX)&&(pcSlice->getPPS()->getEntropyCodingSyncEnabledFlag()) ))
#else
    if( pcSlice->getPPS()->getNumSubstreams() > 1 )
#endif
    {
      // independent tiles => substreams are "per tile".  iNumSubstreams has already been multiplied.
      iNumSubstreamsPerTile = iNumSubstreams/rpcPic->getPicSym()->getNumTiles();
      uiSubStrm = rpcPic->getPicSym()->getTileIdxMap(iCUAddr)*iNumSubstreamsPerTile
                  + uiLin%iNumSubstreamsPerTile;
      m_pcEntropyDecoder->setBitstream( ppcSubstreams[uiSubStrm] );
      // Synchronize cabac probabilities with upper-right LCU if it's available and we're at the start of a line.
#if DEPENDENT_SLICES
      if (((pcSlice->getPPS()->getNumSubstreams() > 1) || bAllowDependence ) && (uiCol == uiTileLCUX)&&(pcSlice->getPPS()->getEntropyCodingSyncEnabledFlag()))
#else
      if (pcSlice->getPPS()->getNumSubstreams() > 1 && uiCol == uiTileLCUX)
#endif
      {
        // We'll sync if the TR is available.
        TComDataCU *pcCUUp = pcCU->getCUAbove();
        UInt uiWidthInCU = rpcPic->getFrameWidthInCU();
        TComDataCU *pcCUTR = NULL;
        if ( pcCUUp && ((iCUAddr%uiWidthInCU+1) < uiWidthInCU)  )
        {
          pcCUTR = rpcPic->getCU( iCUAddr - uiWidthInCU + 1 );
        }
        UInt uiMaxParts = 1<<(pcSlice->getSPS()->getMaxCUDepth()<<1);

        if ( (true && //bEnforceSliceRestriction
             ((pcCUTR==NULL) || (pcCUTR->getSlice()==NULL) || 
             ((pcCUTR->getSCUAddr()+uiMaxParts-1) < pcSlice->getSliceCurStartCUAddr()) ||
             ((rpcPic->getPicSym()->getTileIdxMap( pcCUTR->getAddr() ) != rpcPic->getPicSym()->getTileIdxMap(iCUAddr)))
             ))||
             (true && //bEnforceDependentSliceRestriction
             ((pcCUTR==NULL) || (pcCUTR->getSlice()==NULL) || 
             ((pcCUTR->getSCUAddr()+uiMaxParts-1) < pcSlice->getDependentSliceCurStartCUAddr()) ||
             ((rpcPic->getPicSym()->getTileIdxMap( pcCUTR->getAddr() ) != rpcPic->getPicSym()->getTileIdxMap(iCUAddr)))
             ))
           )
        {
#if DEPENDENT_SLICES
          if( (iCUAddr!=0) && pcCUTR && ((pcCUTR->getSCUAddr()+uiMaxParts-1) >= pcSlice->getSliceCurStartCUAddr()) && bAllowDependence)
          {
             pcSbacDecoders[uiSubStrm].loadContexts( &m_pcBufferSbacDecoders[uiTileCol] ); 
          }
#endif
          // TR not available.
        }
        else
        {
          // TR is available, we use it.
          pcSbacDecoders[uiSubStrm].loadContexts( &m_pcBufferSbacDecoders[uiTileCol] );
        }
      }
      pcSbacDecoder->load(&pcSbacDecoders[uiSubStrm]);  //this load is used to simplify the code (avoid to change all the call to pcSbacDecoders)
    }
    else if ( pcSlice->getPPS()->getNumSubstreams() <= 1 )
    {
      // Set variables to appropriate values to avoid later code change.
      iNumSubstreamsPerTile = 1;
    }

    if ( (iCUAddr == rpcPic->getPicSym()->getTComTile(rpcPic->getPicSym()->getTileIdxMap(iCUAddr))->getFirstCUAddr()) && // 1st in tile.
         (iCUAddr!=0) && (iCUAddr!=rpcPic->getPicSym()->getPicSCUAddr(rpcPic->getSlice(rpcPic->getCurrSliceIdx())->getSliceCurStartCUAddr())/rpcPic->getNumPartInCU())
#if DEPENDENT_SLICES
         && (iCUAddr!=rpcPic->getPicSym()->getPicSCUAddr(rpcPic->getSlice(rpcPic->getCurrSliceIdx())->getDependentSliceCurStartCUAddr())/rpcPic->getNumPartInCU())
#endif
         ) // !1st in frame && !1st in slice
    {
      if (pcSlice->getPPS()->getNumSubstreams() > 1)
      {
        // We're crossing into another tile, tiles are independent.
        // When tiles are independent, we have "substreams per tile".  Each substream has already been terminated, and we no longer
        // have to perform it here.
        // For TILES_DECODER, there can be a header at the start of the 1st substream in a tile.  These are read when the substreams
        // are extracted, not here.
      }
      else
      {
        SliceType sliceType  = pcSlice->getSliceType();
        if (pcSlice->getCabacInitFlag())
        {
          switch (sliceType)
          {
          case P_SLICE:           // change initialization table to B_SLICE intialization
            sliceType = B_SLICE; 
            break;
          case B_SLICE:           // change initialization table to P_SLICE intialization
            sliceType = P_SLICE; 
            break;
          default     :           // should not occur
            assert(0);
          }
        }
        m_pcEntropyDecoder->updateContextTables( sliceType, pcSlice->getSliceQp() );
      }
      
    }

#if ENC_DEC_TRACE
    g_bJustDoIt = g_bEncDecTraceEnable;
#endif
    if ( pcSlice->getSPS()->getUseSAO() && (pcSlice->getSaoEnabledFlag()||pcSlice->getSaoEnabledFlagChroma()) )
    {
      SAOParam *saoParam = rpcPic->getPicSym()->getSaoParam();
      saoParam->bSaoFlag[0] = pcSlice->getSaoEnabledFlag();
      if (iCUAddr == iStartCUAddr)
      {
        saoParam->bSaoFlag[1] = pcSlice->getSaoEnabledFlagChroma();
      }
      Int numCuInWidth     = saoParam->numCuInWidth;
      Int cuAddrInSlice = iCUAddr - rpcPic->getPicSym()->getCUOrderMap(pcSlice->getSliceCurStartCUAddr()/rpcPic->getNumPartInCU());
      Int cuAddrUpInSlice  = cuAddrInSlice - numCuInWidth;
      Int rx = iCUAddr % numCuInWidth;
      Int ry = iCUAddr / numCuInWidth;
      Int allowMergeLeft = 1;
      Int allowMergeUp   = 1;
      if (rx!=0)
      {
        if (rpcPic->getPicSym()->getTileIdxMap(iCUAddr-1) != rpcPic->getPicSym()->getTileIdxMap(iCUAddr))
        {
          allowMergeLeft = 0;
        }
      }
      if (ry!=0)
      {
        if (rpcPic->getPicSym()->getTileIdxMap(iCUAddr-numCuInWidth) != rpcPic->getPicSym()->getTileIdxMap(iCUAddr))
        {
          allowMergeUp = 0;
        }
      }
      pcSbacDecoder->parseSaoOneLcuInterleaving(rx, ry, saoParam,pcCU, cuAddrInSlice, cuAddrUpInSlice, allowMergeLeft, allowMergeUp);
    }
    m_pcCuDecoder->decodeCU     ( pcCU, uiIsLast );
    m_pcCuDecoder->decompressCU ( pcCU );
    
#if ENC_DEC_TRACE
    g_bJustDoIt = g_bEncDecTraceDisable;
#endif
    pcSbacDecoders[uiSubStrm].load(pcSbacDecoder);

    //Store probabilities of second LCU in line into buffer
#if DEPENDENT_SLICES
    if ( (uiCol == uiTileLCUX+1)&& (bAllowDependence || (pcSlice->getPPS()->getNumSubstreams() > 1)) && (pcSlice->getPPS()->getEntropyCodingSyncEnabledFlag()) )
#else
    if (pcSlice->getPPS()->getNumSubstreams() > 1 && (uiCol == uiTileLCUX+1))
#endif
    {
      m_pcBufferSbacDecoders[uiTileCol].loadContexts( &pcSbacDecoders[uiSubStrm] );
    }
#if DEPENDENT_SLICES
    if( uiIsLast && bAllowDependence )
    {
      if (pcSlice->getPPS()->getEntropyCodingSyncEnabledFlag())
       {
         CTXMem[1]->loadContexts( &m_pcBufferSbacDecoders[uiTileCol] );//ctx 2.LCU
       }
      CTXMem[0]->loadContexts( pcSbacDecoder );//ctx end of dep.slice
      return;
    }
#endif
  }*/
}

//#if DEPENDENT_SLICES
func (this *TDecSlice)  InitCtxMem(  i uint){
  for j := 0; j < len(this.CTXMem); j++ {
    delete (this.CTXMem, j);
  }
  
  this.CTXMem = make(map[int]*TDecSbac, i);
}
func (this *TDecSlice)  SetCtxMem( sb *TDecSbac, b int )   { 
	this.CTXMem[b] = sb; 
}
//#endif
//};
/*
type ParameterSetManagerDecoder struct {
    TLibCommon.ParameterSetManager
    //private:
    //  ParameterSetMap<TComVPS> m_vpsBuffer;
    //  ParameterSetMap<TComSPS> m_spsBuffer; 
    //  ParameterSetMap<TComPPS> m_ppsBuffer;
}


func NewParameterSetManagerDecoder() *ParameterSetManagerDecoder{
	return ParameterSetManagerDecoder{TLibCommon.ParameterSetManager{make(map[int]*TLibCommon.TComVPS), 
																	 make(map[int]*TLibCommon.TComSPS), 
																	 make(map[int]*TLibCommon.TComPPS)}}
}

func (this *ParameterSetManagerDecoder)  SetPrefetchedVPS(vps *TLibCommon.TComVPS)  { 
	this.SetVPS(vps); 
}
func (this *ParameterSetManagerDecoder)  GetPrefetchedVPS  (vpsId int) *TLibCommon.TComVPS {
	return this.GetVPS(vpsId)
}
func (this *ParameterSetManagerDecoder)  SetPrefetchedSPS(sps *TLibCommon.TComSPS)  { 
	this.SetSPS(sps); 
};
func (this *ParameterSetManagerDecoder)  GetPrefetchedSPS  (spsId int) *TLibCommon.TComSPS{
	return this.GetSPS(spsId)
}
func (this *ParameterSetManagerDecoder)  SetPrefetchedPPS(pps *TLibCommon.TComPPS)  { 
	this.SetPPS(pps); 
}
func (this *ParameterSetManagerDecoder)  GetPrefetchedPPS  (ppsId int) *TLibCommon.TComPPS{
	return this.GetPPS(ppsId)
}
func (this *ParameterSetManagerDecoder)  ApplyPrefetchedPS() {
}*/
