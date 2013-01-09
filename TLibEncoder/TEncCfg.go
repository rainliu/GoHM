package TLibEncoder

import (
	"gohm/TLibCommon"
)

type GOPEntry struct{
  m_POC				int;
  m_QPOffset		int;
  m_QPFactor		float64;
//#if VARYING_DBL_PARAMS
  m_tcOffsetDiv2	int;
  m_betaOffsetDiv2	int;
//#endif
  m_temporalId		int;
  m_refPic			bool;
  m_numRefPicsActive	int;
  m_sliceType			string;
  m_numRefPics			int;
  m_referencePics	[TLibCommon.MAX_NUM_REF_PICS]int;
  m_usedByCurrPic	[TLibCommon.MAX_NUM_REF_PICS]bool;
//#if AUTO_INTER_RPS
  m_interRPSPrediction	int;
/*#else
  Bool m_interRPSPrediction;
#endif*/
  m_deltaRPS	int;
  m_numRefIdc	int;
  m_refIdc	[TLibCommon.MAX_NUM_REF_PICS+1]int;
}

func NewGOPEntry() *GOPEntry{
	gop := &GOPEntry{ m_POC:-1,
					  m_QPOffset:0,
					  m_QPFactor:0,
					//#if VARYING_DBL_PARAMS
					  m_tcOffsetDiv2:0,
					  m_betaOffsetDiv2:0,
					//#endif
					  m_temporalId:0,
					  m_refPic:false,
					  m_numRefPicsActive:0,
					  m_sliceType:"P",
					  m_numRefPics:0,
					  m_interRPSPrediction:0,
					  m_deltaRPS:0,
					  m_numRefIdc:0};
	
	return gop;
}

func (this *GOPEntry) GetPOC() int{
	return this.m_POC;
}

func (this *GOPEntry) SetPOC(poc int){
	this.m_POC = poc;
}

func (this *GOPEntry) GetQPOffset() int{
	return this.m_QPOffset;
}

func (this *GOPEntry) SetQPOffset(QPOffset int){
	this.m_QPOffset = QPOffset;
}

func (this *GOPEntry) GetQPFactor() float64{
	return this.m_QPFactor;
}

func (this *GOPEntry) SetQPFactor(QPFactor float64){
	this.m_QPFactor = QPFactor;
}

func (this *GOPEntry) GetBetaOffsetDiv2() int {
	return this.m_betaOffsetDiv2;
}

func (this *GOPEntry) GetTcOffsetDiv2() int{
	return this.m_tcOffsetDiv2;
}

func (this *GOPEntry) SetBetaOffsetDiv2(betaOffsetDiv2 int) {
	this.m_betaOffsetDiv2 = betaOffsetDiv2;
}

func (this *GOPEntry) SetTcOffsetDiv2(tcOffsetDiv2 int){
	this.m_tcOffsetDiv2 = tcOffsetDiv2;
}

func (this *GOPEntry) GetNumRefPicsActive() int{
	return this.m_numRefPicsActive;
}

func (this *GOPEntry) SetNumRefPicsActive(numRefPicsActive int){
	this.m_numRefPicsActive = numRefPicsActive;
}

func (this *GOPEntry) GetTemporalId() int{
	return this.m_temporalId;
}
func (this *GOPEntry) SetTemporalId(temporalId int){
	this.m_temporalId=temporalId;
}

func (this *GOPEntry) SetNumRefIdc(numRefIdc int){
	this.m_numRefIdc = numRefIdc;
}
func (this *GOPEntry) GetNumRefIdc() int{
	return this.m_numRefIdc;
}

func (this *GOPEntry) GetNumRefPics() int{
	return this.m_numRefPics;
}
func (this *GOPEntry) SetNumRefPics(numRefPics int){
	this.m_numRefPics = numRefPics;
}

func (this *GOPEntry) GetReferencePics(i int) int{
	return this.m_referencePics[i];
}

func (this *GOPEntry) SetReferencePics(i int, value int){
    this.m_referencePics[i] = value;
}

func (this *GOPEntry) SetRefPic(refPic bool) {
	this.m_refPic = refPic;
}

func (this *GOPEntry) GetRefPic() bool {
	return this.m_refPic;
}

func (this *GOPEntry) SetUsedByCurrPic(i int, b bool){
	this.m_usedByCurrPic[i] = b;
}

func (this *GOPEntry) GetUsedByCurrPic(i int) bool{
	return this.m_usedByCurrPic[i];
}

func (this *GOPEntry) SetInterRPSPrediction(interRPSPrediction int) {
	this.m_interRPSPrediction = interRPSPrediction;
}

func (this *GOPEntry) GetInterRPSPrediction() int {
	return this.m_interRPSPrediction;
}

func (this *GOPEntry) SetRefIdc(i int, refIdc int){
	this.m_refIdc[i] = refIdc;
}

func (this *GOPEntry) GetRefIdc(i int) int{
	return this.m_refIdc[i];
}

func (this *GOPEntry) SetDeltaRPS(deltaRPS int) {
	this.m_deltaRPS = deltaRPS;
}

func (this *GOPEntry) GetDeltaRPS() int {
	return this.m_deltaRPS;
}

func (this *GOPEntry) GetSliceType() string{
	return this.m_sliceType;
}

func (this *GOPEntry) SetSliceType(sliceType string){
	this.m_sliceType = sliceType;
}