package TLibDecoder

import (
	"container/list"
    //"gohm/TLibCommon"
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
    CTXMem		*list.List;//std::vector<TDecSbac*> 
    //#endif
}

//public:
func NewTDecSlice() *TDecSlice {
    return &TDecSlice{}
}

func (this *TDecSlice) Init(pcEntropyDecoder *TDecEntropy, pcMbDecoder *TDecCu) {
}

//func (this *TDecSlice) Create            ( TComSlice* pcSlice, Int iWidth, Int iHeight, UInt uiMaxWidth, UInt uiMaxHeight, UInt uiMaxDepth );
func (this *TDecSlice) Destroy() {
}

//func (this *TDecSlice) DecompressSlice   ( TComInputBitstream* pcBitstream, TComInputBitstream** ppcSubstreams,   TComPic*& rpcPic, TDecSbac* pcSbacDecoder, TDecSbac* pcSbacDecoders );

//#if DEPENDENT_SLICES
func (this *TDecSlice)  InitCtxMem(  i uint){
}
func (this *TDecSlice)  SetCtxMem( sb *TDecSbac, b int )   { 
	//this.CTXMem[b] = sb; 
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
