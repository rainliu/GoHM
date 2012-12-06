package TLibCommon

import (

)


// ====================================================================================================================
// Constants
// ====================================================================================================================
const SAO_MAX_DEPTH         =        4
const SAO_BO_BITS           =        5
const LUMA_GROUP_NUM        =        (1<<SAO_BO_BITS)
const MAX_NUM_SAO_OFFSETS   =        4
const MAX_NUM_SAO_CLASS     =        33
// ====================================================================================================================
// Class definition
// ====================================================================================================================

/// Sample Adaptive Offset class
type TComSampleAdaptiveOffset struct{
//protected:
  m_pcPic		*TComPic;

  m_uiMaxDepth			uint;
  m_aiNumCulPartsLevel	[5]uint;
  m_auiEoTable		[9]uint;
  m_iOffsetBo						*int;
  m_iChromaOffsetBo				*int;
  m_iOffsetEo			[LUMA_GROUP_NUM]int;

  m_iPicWidth			int;
  m_iPicHeight			int;
  m_uiMaxSplitLevel	uint;
  m_uiMaxCUWidth		uint;
  m_uiMaxCUHeight		uint;
  m_iNumCuInWidth		int;
  m_iNumCuInHeight		int;
  m_iNumTotalParts		int;
  m_iNumClass	[MAX_NUM_SAO_TYPE]int;
  m_eSliceType	SliceType;
  m_iPicNalReferenceIdc	int;

  m_uiSaoBitIncreaseY	uint;
  m_uiSaoBitIncreaseC	uint;  //for chroma
  m_uiQP				uint;

  m_pClipTable			*Pel;
  m_pClipTableBase		*Pel;
  m_lumaTableBo			*Pel;
  m_pChromaClipTable		*Pel;
  m_pChromaClipTableBase	*Pel;
  m_chromaTableBo		*Pel;
  m_iUpBuff1				*int;
  m_iUpBuff2				*int;
  m_iUpBufft				*int;
  m_ipSwap					*int;
  m_bUseNIF				bool;       //!< true for performing non-cross slice boundary ALF
  m_uiNumSlicesInPic		uint;      //!< number of slices in picture
  m_iSGDepth				int;              //!< slice granularity depth
  m_pcYuvTmp		*TComPicYuv;    //!< temporary picture buffer pointer when non-across slice/tile boundary SAO is enabled

  m_pTmpU1		*Pel;
  m_pTmpU2		*Pel;
  m_pTmpL1		*Pel;
  m_pTmpL2		*Pel;
  m_iLcuPartIdx	*int;
  m_maxNumOffsetsPerPic	int;
  m_saoLcuBoundary		bool;
  m_saoLcuBasedOptimization	bool;
}
/*
  Void xPCMRestoration        (TComPic* pcPic);
  Void xPCMCURestoration      (TComDataCU* pcCU, UInt uiAbsZorderIdx, UInt uiDepth);
  Void xPCMSampleRestoration  (TComDataCU* pcCU, UInt uiAbsZorderIdx, UInt uiDepth, TextType ttText);
public:
  TComSampleAdaptiveOffset         ();
  virtual ~TComSampleAdaptiveOffset();
*/
func (this *TComSampleAdaptiveOffset) Create( uiSourceWidth, uiSourceHeight, uiMaxCUWidth, uiMaxCUHeight uint ){
}
func (this *TComSampleAdaptiveOffset) Destroy (){
}
/*
  Int  convertLevelRowCol2Idx(Int level, Int row, Int col);

  Void initSAOParam   (SAOParam *pcSaoParam, Int iPartLevel, Int iPartRow, Int iPartCol, Int iParentPartIdx, Int StartCUX, Int EndCUX, Int StartCUY, Int EndCUY, Int iYCbCr);
  Void allocSaoParam  (SAOParam* pcSaoParam);
  Void resetSAOParam  (SAOParam *pcSaoParam);
  static Void freeSaoParam   (SAOParam *pcSaoParam);
  
  Void SAOProcess(SAOParam* pcSaoParam);
  Void processSaoCu(Int iAddr, Int iSaoType, Int iYCbCr);
  Pel* getPicYuvAddr(TComPicYuv* pcPicYuv, Int iYCbCr,Int iAddr = 0);

  Void processSaoCuOrg(Int iAddr, Int iPartIdx, Int iYCbCr);  //!< LCU-basd SAO process without slice granularity 
  Void createPicSaoInfo(TComPic* pcPic, Int numSlicesInPic = 1);
  Void destroyPicSaoInfo();
  Void processSaoBlock(Pel* pDec, Pel* pRest, Int stride, Int iSaoType, UInt width, UInt height, Bool* pbBorderAvail, Int iYCbCr);

  Void resetLcuPart(SaoLcuParam* saoLcuParam);
  Void convertQT2SaoUnit(SAOParam* saoParam, UInt partIdx, Int yCbCr);
  Void convertOnePart2SaoUnit(SAOParam *saoParam, UInt partIdx, Int yCbCr);
  Void processSaoUnitAll(SaoLcuParam* saoLcuParam, Bool oneUnitFlag, Int yCbCr);
  Void setSaoLcuBoundary (Bool bVal)  {m_saoLcuBoundary = bVal;}
  Bool getSaoLcuBoundary ()           {return m_saoLcuBoundary;}
  Void setSaoLcuBasedOptimization (Bool bVal)  {m_saoLcuBasedOptimization = bVal;}
  Bool getSaoLcuBasedOptimization ()           {return m_saoLcuBasedOptimization;}
  Void resetSaoUnit(SaoLcuParam* saoUnit);
  Void copySaoUnit(SaoLcuParam* saoUnitDst, SaoLcuParam* saoUnitSrc );
  Void PCMLFDisableProcess    ( TComPic* pcPic);                        ///< interface function for ALF process 
};*/