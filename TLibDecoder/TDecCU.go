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
    m_uiMaxDepth uint ///< max. number of depth
    //TComYuv**           m_ppcYuvResi;       ///< array of residual buffer
    //TComYuv**           m_ppcYuvReco;       ///< array of prediction & reconstruction buffer
    //TComDataCU**        m_ppcCU;            ///< CU data array

    // access channel
    m_pcTrQuant        *TLibCommon.TComTrQuant
    m_pcPrediction     *TLibCommon.TComPrediction
    m_pcEntropyDecoder *TDecEntropy

    m_bDecodeDQP bool
}

/*
public:
  TDecCu();
  virtual ~TDecCu();

  /// initialize access channels
  Void  init                    ( TDecEntropy* pcEntropyDecoder, TComTrQuant* pcTrQuant, TComPrediction* pcPrediction );

  /// create internal buffers
  Void  create                  ( UInt uiMaxDepth, UInt uiMaxWidth, UInt uiMaxHeight );
*/
/// destroy internal buffers
func (this *TDecCu) Destroy() {
}

/*  
  /// decode CU information
  Void  decodeCU                ( TComDataCU* pcCU, UInt& ruiIsLast );

  /// reconstruct CU information
  Void  decompressCU            ( TComDataCU* pcCU );

protected:

  Void xDecodeCU                ( TComDataCU* pcCU,                       UInt uiAbsPartIdx, UInt uiDepth, UInt &ruiIsLast);
  Void xFinishDecodeCU          ( TComDataCU* pcCU,                       UInt uiAbsPartIdx, UInt uiDepth, UInt &ruiIsLast);
  Bool xDecodeSliceEnd          ( TComDataCU* pcCU,                       UInt uiAbsPartIdx, UInt uiDepth);
  Void xDecompressCU            ( TComDataCU* pcCU, TComDataCU* pcCUCur,  UInt uiAbsPartIdx, UInt uiDepth );

  Void xReconInter              ( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiDepth );

  Void  xReconIntraQT           ( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiDepth );
  Void  xIntraRecLumaBlk        ( TComDataCU* pcCU, UInt uiTrDepth, UInt uiAbsPartIdx, TComYuv* pcRecoYuv, TComYuv* pcPredYuv, TComYuv* pcResiYuv );
  Void  xIntraRecChromaBlk      ( TComDataCU* pcCU, UInt uiTrDepth, UInt uiAbsPartIdx, TComYuv* pcRecoYuv, TComYuv* pcPredYuv, TComYuv* pcResiYuv, UInt uiChromaId );

  Void  xReconPCM               ( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiDepth );

  Void xDecodeInterTexture      ( TComDataCU* pcCU, UInt uiAbsPartIdx, UInt uiDepth );
  Void xDecodePCMTexture        ( TComDataCU* pcCU, UInt uiPartIdx, Pel *piPCM, Pel* piReco, UInt uiStride, UInt uiWidth, UInt uiHeight, TextType ttText);

  Void xCopyToPic               ( TComDataCU* pcCU, TComPic* pcPic, UInt uiZorderIdx, UInt uiDepth );

  Void  xIntraLumaRecQT         ( TComDataCU* pcCU, UInt uiTrDepth, UInt uiAbsPartIdx, TComYuv* pcRecoYuv, TComYuv* pcPredYuv, TComYuv* pcResiYuv );
  Void  xIntraChromaRecQT       ( TComDataCU* pcCU, UInt uiTrDepth, UInt uiAbsPartIdx, TComYuv* pcRecoYuv, TComYuv* pcPredYuv, TComYuv* pcResiYuv );

  Bool getdQPFlag               ()                        { return m_bDecodeDQP;        }
  Void setdQPFlag               ( Bool b )                { m_bDecodeDQP = b;           }
  Void xFillPCMBuffer           (TComDataCU* pCU, UInt absPartIdx, UInt depth);
};*/
