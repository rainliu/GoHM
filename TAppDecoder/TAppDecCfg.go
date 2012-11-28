package TAppDecoder

import (
	"strconv"
	"errors"
)

type TAppDecCfg struct{
  m_pchBitstreamFile	string;                ///< input bitstream file name
  m_pchReconFile		string;                ///< output reconstruction file name
  m_iFrameNum			int;				   ///< output frame number
  m_iSkipFrame			int;                   ///< counter for frames prior to the random access point to skip
  m_outputBitDepthY		int;                   ///< bit depth used for writing output (luma)
  m_outputBitDepthC		int;                   ///< bit depth used for writing output (chroma)t
  m_iMaxTemporalLayer	int;                   ///< maximum temporal layer to be decoded
  m_decodedPictureHashSEIEnabled	int;       ///< Checksum(3)/CRC(2)/MD5(1)/disable(0) acting on decoded picture hash SEI message
  //m_targetDecLayerIdSet	list.List;             ///< set of LayerIds to be included in the sub-bitstream extraction process.
}  
  
func NewTAppDecCfg() (*TAppDecCfg){
	return &TAppDecCfg{}
}

func (this *TAppDecCfg) ParseCfg(argc int, argv []string) (err error){   ///< initialize option class from configuration
	if argc <= 2 {
		err = errors.New("Too few arguments")
		return err
	}
	
	this.m_iFrameNum = -1
	this.m_iSkipFrame = 0
	this.m_outputBitDepthY = 0
	this.m_outputBitDepthC = 0
	this.m_iMaxTemporalLayer = -1
	this.m_decodedPictureHashSEIEnabled = 1
	//this.m_targetDecLayerIdSet = 0
	
	if argc >= 3{
		this.m_pchBitstreamFile = argv[1]
		this.m_pchReconFile     = argv[2]
	}
	
	if argc >= 4 {
		this.m_iFrameNum, err = strconv.Atoi(argv[3])
		if err!=nil{
			return err
		}
	}
	
	return nil
}
