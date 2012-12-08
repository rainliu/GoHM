package TLibDecoder

import (
    "gohm/TLibCommon"
)

// ====================================================================================================================
// Class definition
// ====================================================================================================================

type TDecBinIf interface {
    /*
    public:
      virtual Void  init              ( TComInputBitstream* pcTComBitstream )     = 0;
      virtual Void  uninit            ()                                          = 0;

      virtual Void  start             ()                                          = 0;
      virtual Void  finish            ()                                          = 0;
      virtual Void  flush            ()                                           = 0;

      virtual Void  decodeBin         ( UInt& ruiBin, ContextModel& rcCtxModel )  = 0;
      virtual Void  decodeBinEP       ( UInt& ruiBin                           )  = 0;
      virtual Void  decodeBinsEP      ( UInt& ruiBins, Int numBins             )  = 0;
      virtual Void  decodeBinTrm      ( UInt& ruiBin                           )  = 0;

      virtual Void  resetBac          ()                                          = 0;
    #if !REMOVE_BURST_IPCM
      virtual Void  decodeNumSubseqIPCM( Int& numSubseqIPCM )                  = 0;
    #endif
      virtual Void  decodePCMAlignBits()                                          = 0;
      virtual Void  xReadPCMCode      ( UInt uiLength, UInt& ruiCode)              = 0;

      virtual ~TDecBinIf() {}

      virtual Void  copyState         ( TDecBinIf* pcTDecBinIf )                  = 0;
      virtual TDecBinCABAC*   getTDecBinCABAC   ()  { return 0; }
    */
}

type TDecBinCabac struct { //: public TDecBinIf
    //private:
    m_pcTComBitstream *TLibCommon.TComInputBitstream
    m_uiRange         uint
    m_uiValue         uint
    m_bitsNeeded      int
}

/*
public:
  TDecBinCABAC ();
  virtual ~TDecBinCABAC();

  Void  init              ( TComInputBitstream* pcTComBitstream );
  Void  uninit            ();

  Void  start             ();
  Void  finish            ();
  Void  flush             ();

  Void  decodeBin         ( UInt& ruiBin, ContextModel& rcCtxModel );
  Void  decodeBinEP       ( UInt& ruiBin                           );
  Void  decodeBinsEP      ( UInt& ruiBin, Int numBins              );
  Void  decodeBinTrm      ( UInt& ruiBin                           );

  Void  resetBac          ();
#if !REMOVE_BURST_IPCM
  Void  decodeNumSubseqIPCM( Int& numSubseqIPCM ) ;
#endif
  Void  decodePCMAlignBits();
  Void  xReadPCMCode      ( UInt uiLength, UInt& ruiCode );

  Void  copyState         ( TDecBinIf* pcTDecBinIf );
  TDecBinCABAC* getTDecBinCABAC()  { return this; }


};*/

//class SEImessages;

/// SBAC decoder class
type TDecSbac struct { //: public TDecEntropyIf
    //private:
    m_pcBitstream *TLibCommon.TComInputBitstream
    m_pcTDecBinIf		TDecBinIf;

    //private:
    m_uiLastDQpNonZero uint
    m_uiLastQp         uint

    m_contextModels             [TLibCommon.MAX_NUM_CTX_MOD]TLibCommon.ContextModel
    m_numContextModels          int
    m_cCUSplitFlagSCModel       TLibCommon.ContextModel3DBuffer
    m_cCUSkipFlagSCModel        TLibCommon.ContextModel3DBuffer
    m_cCUMergeFlagExtSCModel    TLibCommon.ContextModel3DBuffer
    m_cCUMergeIdxExtSCModel     TLibCommon.ContextModel3DBuffer
    m_cCUPartSizeSCModel        TLibCommon.ContextModel3DBuffer
    m_cCUPredModeSCModel        TLibCommon.ContextModel3DBuffer
    m_cCUIntraPredSCModel       TLibCommon.ContextModel3DBuffer
    m_cCUChromaPredSCModel      TLibCommon.ContextModel3DBuffer
    m_cCUDeltaQpSCModel         TLibCommon.ContextModel3DBuffer
    m_cCUInterDirSCModel        TLibCommon.ContextModel3DBuffer
    m_cCURefPicSCModel          TLibCommon.ContextModel3DBuffer
    m_cCUMvdSCModel             TLibCommon.ContextModel3DBuffer
    m_cCUQtCbfSCModel           TLibCommon.ContextModel3DBuffer
    m_cCUTransSubdivFlagSCModel TLibCommon.ContextModel3DBuffer
    m_cCUQtRootCbfSCModel       TLibCommon.ContextModel3DBuffer

    m_cCUSigCoeffGroupSCModel TLibCommon.ContextModel3DBuffer
    m_cCUSigSCModel           TLibCommon.ContextModel3DBuffer
    m_cCuCtxLastX             TLibCommon.ContextModel3DBuffer
    m_cCuCtxLastY             TLibCommon.ContextModel3DBuffer
    m_cCUOneSCModel           TLibCommon.ContextModel3DBuffer
    m_cCUAbsSCModel           TLibCommon.ContextModel3DBuffer

    m_cMVPIdxSCModel TLibCommon.ContextModel3DBuffer

    m_cCUAMPSCModel                 TLibCommon.ContextModel3DBuffer
    m_cSaoMergeSCModel              TLibCommon.ContextModel3DBuffer
    m_cSaoTypeIdxSCModel            TLibCommon.ContextModel3DBuffer
    m_cTransformSkipSCModel         TLibCommon.ContextModel3DBuffer
    m_CUTransquantBypassFlagSCModel TLibCommon.ContextModel3DBuffer
}


func NewTDecSbac() *TDecSbac{
	return &TDecSbac{}
}

  

func (this *TDecSbac) Init ( p TDecBinIf)    { 
	this.m_pcTDecBinIf = p; 
}
//  Void  uninit                    (              )    { m_pcTDecBinIf = 0; }
/*
  Void load                          ( TDecSbac* pScr );
  Void loadContexts                  ( TDecSbac* pScr );
  Void xCopyFrom           ( TDecSbac* pSrc );
  Void xCopyContextsFrom       ( TDecSbac* pSrc );

  Void  resetEntropy (TComSlice* pSlice );
  Void  setBitstream              ( TComInputBitstream* p  ) { m_pcBitstream = p; m_pcTDecBinIf->init( p ); }
  Void  parseVPS                  ( TComVPS* pcVPS )  {}
  Void  parseSPS                  ( TComSPS* pcSPS         ) {}
  Void  parsePPS                  ( TComPPS* pcPPS         ) {}

  Void  parseSliceHeader          ( TComSlice*& rpcSlice, ParameterSetManagerDecoder *parameterSetManager) {}
  Void  parseTerminatingBit       ( UInt& ruiBit );
  Void  parseMVPIdx               ( Int& riMVPIdx          );
  Void  parseSaoMaxUvlc           ( UInt& val, UInt maxSymbol );
  Void  parseSaoMerge         ( UInt&  ruiVal   );
  Void  parseSaoTypeIdx           ( UInt&  ruiVal  );
  Void  parseSaoUflc              ( UInt uiLength, UInt& ruiVal     );
  Void  parseSaoOneLcuInterleaving(Int rx, Int ry, SAOParam* pSaoParam, TComDataCU* pcCU, Int iCUAddrInSlice, Int iCUAddrUpInSlice, Int allowMergeLeft, Int allowMergeUp);
  Void  parseSaoOffset            (SaoLcuParam* psSaoLcuParam, UInt compIdx);
private:
  Void  xReadUnarySymbol    ( UInt& ruiSymbol, ContextModel* pcSCModel, Int iOffset );
  Void  xReadUnaryMaxSymbol ( UInt& ruiSymbol, ContextModel* pcSCModel, Int iOffset, UInt uiMaxSymbol );
  Void  xReadEpExGolomb     ( UInt& ruiSymbol, UInt uiCount );
  Void  xReadCoefRemainExGolomb ( UInt &rSymbol, UInt &rParam );


public:

  Void parseSkipFlag      ( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiDepth );
  Void parseCUTransquantBypassFlag( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiDepth );
  Void parseSplitFlag     ( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiDepth );
  Void parseMergeFlag     ( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiDepth, UInt uiPUIdx );
  Void parseMergeIndex    ( TComDataCU* pcCU, UInt& ruiMergeIndex, UInt uiAbsPartIdx, UInt uiDepth );
  Void parsePartSize      ( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiDepth );
  Void parsePredMode      ( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiDepth );

  Void parseIntraDirLumaAng( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiDepth );

  Void parseIntraDirChroma( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiDepth );

  Void parseInterDir      ( TComDataCU* pcCU, UInt& ruiInterDir, UInt uiAbsPartIdx, UInt uiDepth );
  Void parseRefFrmIdx     ( TComDataCU* pcCU, Int& riRefFrmIdx, UInt uiAbsPartIdx, UInt uiDepth, RefPicList eRefList );
  Void parseMvd           ( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiPartIdx, UInt uiDepth, RefPicList eRefList );

  Void parseTransformSubdivFlag( UInt& ruiSubdivFlag, UInt uiLog2TransformBlockSize );
  Void parseQtCbf         ( TComDataCU* pcCU, UInt uiAbsPartIdx, TextType eType, UInt uiTrDepth, UInt uiDepth );
  Void parseQtRootCbf     ( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiDepth, UInt& uiQtRootCbf );

  Void parseDeltaQP       ( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiDepth );

  Void parseIPCMInfo      ( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiDepth);

  Void parseLastSignificantXY( UInt& uiPosLastX, UInt& uiPosLastY, Int width, Int height, TextType eTType, UInt uiScanIdx );
  Void parseCoeffNxN      ( TComDataCU* pcCU, TCoeff* pcCoef, UInt uiAbsPartIdx, UInt uiWidth, UInt uiHeight, UInt uiDepth, TextType eTType );
  Void parseTransformSkipFlags ( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt width, UInt height, UInt uiDepth, TextType eTType);

  Void updateContextTables( SliceType eSliceType, Int iQp );

  Void  parseScalingList ( TComScalingList* scalingList ) {}

};*/
