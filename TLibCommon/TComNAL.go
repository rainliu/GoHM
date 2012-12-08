package TLibCommon

import (
    //"container/list"
)

/**
 * Represents a single NALunit header and the associated RBSPayload
 */
type NALUnit struct {
    m_nalUnitType       NalUnitType ///< nal_unit_type
    m_temporalId        uint        ///< temporal_id
    m_reservedZero6Bits uint        ///< reserved_zero_6bits
}

/** construct an NALunit structure with given header values. */
func NewNALUnit(nalUnitType NalUnitType, temporalId, reservedZero6Bits uint) *NALUnit {
    return &NALUnit{nalUnitType, temporalId, reservedZero6Bits}
}

/** returns true if the NALunit is a slice NALunit */
func (this *NALUnit) IsSlice() bool {
    return this.m_nalUnitType == NAL_UNIT_CODED_SLICE_TRAIL_R ||
        this.m_nalUnitType == NAL_UNIT_CODED_SLICE_TRAIL_N ||
        this.m_nalUnitType == NAL_UNIT_CODED_SLICE_TLA ||
        this.m_nalUnitType == NAL_UNIT_CODED_SLICE_TSA_N ||
        this.m_nalUnitType == NAL_UNIT_CODED_SLICE_STSA_R ||
        this.m_nalUnitType == NAL_UNIT_CODED_SLICE_STSA_N ||
        this.m_nalUnitType == NAL_UNIT_CODED_SLICE_BLA ||
        this.m_nalUnitType == NAL_UNIT_CODED_SLICE_BLANT ||
        this.m_nalUnitType == NAL_UNIT_CODED_SLICE_BLA_N_LP ||
        this.m_nalUnitType == NAL_UNIT_CODED_SLICE_IDR ||
        this.m_nalUnitType == NAL_UNIT_CODED_SLICE_IDR_N_LP ||
        this.m_nalUnitType == NAL_UNIT_CODED_SLICE_CRA ||
        this.m_nalUnitType == NAL_UNIT_CODED_SLICE_DLP ||
        this.m_nalUnitType == NAL_UNIT_CODED_SLICE_TFD
}

func (this *NALUnit) GetReservedZero6Bits() uint {
	return this.m_reservedZero6Bits
}
func (this *NALUnit) GetTemporalId() uint{
	return this.m_temporalId;
}
func (this *NALUnit) SetTemporalId(temporalId uint){
	this.m_temporalId = temporalId;
}

func (this *NALUnit) SetReservedZero6Bits(reservedZero6Bits uint) {
	this.m_reservedZero6Bits = reservedZero6Bits
}

func (this *NALUnit) GetNalUnitType() NalUnitType {
    return this.m_nalUnitType
}

func (this *NALUnit) SetNalUnitType(nalUnitType NalUnitType) {
    this.m_nalUnitType = nalUnitType
}

