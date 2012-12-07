package TLibCommon

import ()

// ====================================================================================================================
// Constants
// ====================================================================================================================

const MAX_NUM_CTX_MOD = 512 ///< maximum number of supported contexts

const NUM_SPLIT_FLAG_CTX = 3 ///< number of context models for split flag
const NUM_SKIP_FLAG_CTX = 3  ///< number of context models for skip flag

const NUM_MERGE_FLAG_EXT_CTX = 1 ///< number of context models for merge flag of merge extended
const NUM_MERGE_IDX_EXT_CTX = 1  ///< number of context models for merge index of merge extended

const NUM_PART_SIZE_CTX = 4 ///< number of context models for partition size
const NUM_CU_AMP_CTX = 1    ///< number of context models for partition size (AMP)
const NUM_PRED_MODE_CTX = 1 ///< number of context models for prediction mode

const NUM_ADI_CTX = 1 ///< number of context models for intra prediction

const NUM_CHROMA_PRED_CTX = 2 ///< number of context models for intra prediction (chroma)
const NUM_INTER_DIR_CTX = 5   ///< number of context models for inter prediction direction
const NUM_MV_RES_CTX = 2      ///< number of context models for motion vector difference

const NUM_REF_NO_CTX = 2            ///< number of context models for reference index
const NUM_TRANS_SUBDIV_FLAG_CTX = 3 ///< number of context models for transform subdivision flags
const NUM_QT_CBF_CTX = 5            ///< number of context models for QT CBF
const NUM_QT_ROOT_CBF_CTX = 1       ///< number of context models for QT ROOT CBF
const NUM_DELTA_QP_CTX = 3          ///< number of context models for dQP

const NUM_SIG_CG_FLAG_CTX = 2 ///< number of context models for MULTI_LEVEL_SIGNIFICANCE

const NUM_SIG_FLAG_CTX = 42        ///< number of context models for sig flag
const NUM_SIG_FLAG_CTX_LUMA = 27   ///< number of context models for luma sig flag
const NUM_SIG_FLAG_CTX_CHROMA = 15 ///< number of context models for chroma sig flag

const NUM_CTX_LAST_FLAG_XY = 15 ///< number of context models for last coefficient position

const NUM_ONE_FLAG_CTX = 24       ///< number of context models for greater than 1 flag
const NUM_ONE_FLAG_CTX_LUMA = 16  ///< number of context models for greater than 1 flag of luma
const NUM_ONE_FLAG_CTX_CHROMA = 8 ///< number of context models for greater than 1 flag of chroma
const NUM_ABS_FLAG_CTX = 6        ///< number of context models for greater than 2 flag
const NUM_ABS_FLAG_CTX_LUMA = 4   ///< number of context models for greater than 2 flag of luma
const NUM_ABS_FLAG_CTX_CHROMA = 2 ///< number of context models for greater than 2 flag of chroma

const NUM_MVP_IDX_CTX = 2 ///< number of context models for MVP index

const NUM_SAO_MERGE_FLAG_CTX = 1 ///< number of context models for SAO merge flags
const NUM_SAO_TYPE_IDX_CTX = 1   ///< number of context models for SAO type index

const NUM_TRANSFORMSKIP_FLAG_CTX = 1 ///< number of context models for transform skipping 
const NUM_CU_TRANSQUANT_BYPASS_FLAG_CTX = 1
const CNU = 154 ///< dummy initialization value for unused context models 'Context model Not Used'
