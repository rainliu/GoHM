package TLibCommon

import ()

// ====================================================================================================================
// Class definition
// ====================================================================================================================

/// context model 3D buffer class
type ContextModel3DBuffer struct {
    //protected:
    m_contextModel *ContextModel ///< array of context models
    m_sizeX        uint          ///< X size of 3D buffer
    m_sizeXY       uint          ///< X times Y size of 3D buffer
    m_sizeXYZ      uint          ///< total size of 3D buffer
}

/* 
public:
  ContextModel3DBuffer  ( UInt uiSizeZ, UInt uiSizeY, UInt uiSizeX, ContextModel *basePtr, Int &count );
  ~ContextModel3DBuffer () {}

  // access functions
  ContextModel& get( UInt uiZ, UInt uiY, UInt uiX )
  {
    return  m_contextModel[ uiZ * m_sizeXY + uiY * m_sizeX + uiX ];
  }
  ContextModel* get( UInt uiZ, UInt uiY )
  {
    return &m_contextModel[ uiZ * m_sizeXY + uiY * m_sizeX ];
  }
  ContextModel* get( UInt uiZ )
  {
    return &m_contextModel[ uiZ * m_sizeXY ];
  }

  // initialization & copy functions
  Void initBuffer( SliceType eSliceType, Int iQp, UChar* ctxModel );          ///< initialize 3D buffer by slice type & QP

  UInt calcCost( SliceType sliceType, Int qp, UChar* ctxModel );      ///< determine cost of choosing a probability table based on current probabilities
  // copy from another buffer
  // \param src buffer to copy from

  Void copyFrom( ContextModel3DBuffer* src )
  {
    assert( m_sizeXYZ == src->m_sizeXYZ );
    ::memcpy( m_contextModel, src->m_contextModel, sizeof(ContextModel) * m_sizeXYZ );
  }
};*/
