package TAppDecoder

import (
	"gohm/TLibCommon"
)

type TAppDecTop struct{
	TAppDecCfg
	
	//m_cTDecTop	TDecTop
	m_cTVideoIOYuvReconFile	TLibCommon.TVideoIOYuv
	
	m_abDecFlag	[TLibCommon.MAX_GOP]bool
	m_iPOCLastDisplay	int
}

func NewTAppDecTop() (*TAppDecTop){
	pAppDecTop := &TAppDecTop{}
	pAppDecTop.m_iPOCLastDisplay = -TLibCommon.MAX_GOP
	
	return pAppDecTop
}

func (this *TAppDecTop) Create(){
}

func (this *TAppDecTop) Destroy(){
}

func (this *TAppDecTop) Decode(){
}

func (this *TAppDecTop) xCreateDecLib(){
	//create decoder class
  	//this.m_cTDecTop.create();
}

func (this *TAppDecTop) xDestroyDecLib(){
   	if this.m_pchReconFile != "" {
    	this.m_cTVideoIOYuvReconFile.Close();
  	}
  
	//destroy decoder class
  	//this.m_cTDecTop.destroy();
}

func (this *TAppDecTop) xInitDecLib(){
	//initialize decoder class
  	//this.m_cTDecTop.init();
  	//this.m_cTDecTop.setDecodedPictureHashSEIEnabled(this.m_decodedPictureHashSEIEnabled);
	
}

func (this *TAppDecTop) xWriteOutput(){
}

func (this *TAppDecTop) xFlushOutput(){
}

func (this *TAppDecTop) isNaluWithinTargetDecLayerIdSet(){
}