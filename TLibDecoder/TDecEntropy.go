package TLibDecoder

import (
    "gohm/TLibCommon"
)

// ====================================================================================================================
// Class definition
// ====================================================================================================================

/// entropy decoder pure class
type TDecEntropyIf interface {
    /*public:
      //  Virtual list for SBAC/CAVLC
      virtual Void  resetEntropy          ( TComSlice* pcSlice )     = 0;
      virtual Void  setBitstream          ( TComInputBitstream* p )  = 0;

      virtual Void  parseVPS                  ( TComVPS* pcVPS )                       = 0;
      virtual Void  parseSPS                  ( TComSPS* pcSPS )                                      = 0;
      virtual Void  parsePPS                  ( TComPPS* pcPPS )                                      = 0;

      virtual Void parseSliceHeader          ( TComSlice*& rpcSlice, ParameterSetManagerDecoder *parameterSetManager)       = 0;

      virtual Void  parseTerminatingBit       ( UInt& ruilsLast )                                     = 0;

      virtual Void parseMVPIdx        ( Int& riMVPIdx ) = 0;

    public:
      virtual Void parseSkipFlag      ( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiDepth ) = 0;
      virtual Void parseCUTransquantBypassFlag( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiDepth ) = 0;
      virtual Void parseSplitFlag     ( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiDepth ) = 0;
      virtual Void parseMergeFlag     ( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiDepth, UInt uiPUIdx ) = 0;
      virtual Void parseMergeIndex    ( TComDataCU* pcCU, UInt& ruiMergeIndex, UInt uiAbsPartIdx, UInt uiDepth ) = 0;
      virtual Void parsePartSize      ( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiDepth ) = 0;
      virtual Void parsePredMode      ( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiDepth ) = 0;

      virtual Void parseIntraDirLumaAng( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiDepth ) = 0;

      virtual Void parseIntraDirChroma( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiDepth ) = 0;

      virtual Void parseInterDir      ( TComDataCU* pcCU, UInt& ruiInterDir, UInt uiAbsPartIdx, UInt uiDepth ) = 0;
      virtual Void parseRefFrmIdx     ( TComDataCU* pcCU, Int& riRefFrmIdx, UInt uiAbsPartIdx, UInt uiDepth, RefPicList eRefList ) = 0;
      virtual Void parseMvd           ( TComDataCU* pcCU, UInt uiAbsPartAddr, UInt uiPartIdx, UInt uiDepth, RefPicList eRefList ) = 0;

      virtual Void parseTransformSubdivFlag( UInt& ruiSubdivFlag, UInt uiLog2TransformBlockSize ) = 0;
      virtual Void parseQtCbf         ( TComDataCU* pcCU, UInt uiAbsPartIdx, TextType eType, UInt uiTrDepth, UInt uiDepth ) = 0;
      virtual Void parseQtRootCbf     ( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiDepth, UInt& uiQtRootCbf ) = 0;

      virtual Void parseDeltaQP       ( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiDepth ) = 0;

      virtual Void parseIPCMInfo     ( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiDepth) = 0;

      virtual Void parseCoeffNxN( TComDataCU* pcCU, TCoeff* pcCoef, UInt uiAbsPartIdx, UInt uiWidth, UInt uiHeight, UInt uiDepth, TextType eTType ) = 0;
      virtual Void parseTransformSkipFlags ( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt width, UInt height, UInt uiDepth, TextType eTType) = 0;
      virtual Void updateContextTables( SliceType eSliceType, Int iQp ) = 0;

      virtual ~TDecEntropyIf() {}
    */
}

/// entropy decoder class
type TDecEntropy struct {
    //private:
    //TDecEntropyIf*  m_pcEntropyDecoderIf;
    m_pcPrediction      *TLibCommon.TComPrediction
    m_uiBakAbsPartIdx   uint
    m_uiBakChromaOffset uint
    m_bakAbsPartIdxCU   uint
}

//public:
func (this *TDecEntropy) Init(p *TLibCommon.TComPrediction) {
    this.m_pcPrediction = p
}

/*  Void decodePUWise       ( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiDepth, TComDataCU* pcSubCU );
Void decodeInterDirPU   ( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiDepth, UInt uiPartIdx );
Void decodeRefFrmIdxPU  ( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiDepth, UInt uiPartIdx, RefPicList eRefList );
Void decodeMvdPU        ( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiDepth, UInt uiPartIdx, RefPicList eRefList );
Void decodeMVPIdxPU     ( TComDataCU* pcSubCU, UInt uiPartAddr, UInt uiDepth, UInt uiPartIdx, RefPicList eRefList );
*/
func (this *TDecEntropy) SetEntropyDecoder(p *TDecEntropyIf) {
}
func (this *TDecEntropy) SetBitstream(p *TLibCommon.TComInputBitstream) {
    //this.m_pcEntropyDecoderIf.SetBitstream(p);                    
}

/* 
  Void    resetEntropy                ( TComSlice* p)           { m_pcEntropyDecoderIf->resetEntropy(p);                    }
  Void    decodeVPS                   ( TComVPS* pcVPS ) { m_pcEntropyDecoderIf->parseVPS(pcVPS); }
  Void    decodeSPS                   ( TComSPS* pcSPS     )    { m_pcEntropyDecoderIf->parseSPS(pcSPS);                    }
  Void    decodePPS                   ( TComPPS* pcPPS, ParameterSetManagerDecoder *parameterSet    )    { m_pcEntropyDecoderIf->parsePPS(pcPPS);                    }
  Void    decodeSliceHeader           ( TComSlice*& rpcSlice, ParameterSetManagerDecoder *parameterSetManager)  { m_pcEntropyDecoderIf->parseSliceHeader(rpcSlice, parameterSetManager);         }

  Void    decodeTerminatingBit        ( UInt& ruiIsLast )       { m_pcEntropyDecoderIf->parseTerminatingBit(ruiIsLast);     }

  TDecEntropyIf* getEntropyDecoder() { return m_pcEntropyDecoderIf; }

public:
  Void decodeSplitFlag         ( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiDepth );
  Void decodeSkipFlag          ( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiDepth );
  Void decodeCUTransquantBypassFlag( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiDepth );
  Void decodeMergeFlag         ( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiDepth, UInt uiPUIdx );
  Void decodeMergeIndex        ( TComDataCU* pcSubCU, UInt uiPartIdx, UInt uiPartAddr, PartSize eCUMode, UChar* puhInterDirNeighbours, TComMvField* pcMvFieldNeighbours, UInt uiDepth );
  Void decodePredMode          ( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiDepth );
  Void decodePartSize          ( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiDepth );

  Void decodeIPCMInfo          ( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiDepth );

  Void decodePredInfo          ( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiDepth, TComDataCU* pcSubCU );

  Void decodeIntraDirModeLuma  ( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiDepth );
  Void decodeIntraDirModeChroma( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiDepth );

  Void decodeQP                ( TComDataCU* pcCU, UInt uiAbsPartIdx );

  Void updateContextTables    ( SliceType eSliceType, Int iQp ) { m_pcEntropyDecoderIf->updateContextTables( eSliceType, iQp ); }


private:
  Void xDecodeTransform        ( TComDataCU* pcCU, UInt offsetLuma, UInt offsetChroma, UInt uiAbsPartIdx, UInt absTUPartIdx, UInt uiDepth, UInt width, UInt height, UInt uiTrIdx, UInt uiInnerQuadIdx, Bool& bCodeDQP );

public:
  Void decodeCoeff             ( TComDataCU* pcCU                 , UInt uiAbsPartIdx, UInt uiDepth, UInt uiWidth, UInt uiHeight, Bool& bCodeDQP );

};// END CLASS DEFINITION TDecEntropy
*/
