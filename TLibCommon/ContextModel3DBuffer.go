/* The copyright in this software is being made available under the BSD
 * License, included below. This software may be subject to other third party
 * and contributor rights, including patent rights, and no such rights are
 * granted under this license.
 *
 * Copyright (c) 2012-2013, H265.net
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 *  * Redistributions of source code must retain the above copyright notice,
 *    this list of conditions and the following disclaimer.
 *  * Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *  * Neither the name of the H265.net nor the names of its contributors may
 *    be used to endorse or promote products derived from this software without
 *    specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS
 * BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF
 * THE POSSIBILITY OF SUCH DAMAGE.
 */

package TLibCommon

import ()

// ====================================================================================================================
// Class definition
// ====================================================================================================================

/// context model 3D buffer class
type ContextModel3DBuffer struct {
    //protected:
    m_contextModel []ContextModel ///< array of context models
    m_sizeX        uint           ///< X size of 3D buffer
    m_sizeXY       uint           ///< X times Y size of 3D buffer
    m_sizeXYZ      uint           ///< total size of 3D buffer
}

func NewContextModel3DBuffer(uiSizeZ, uiSizeY, uiSizeX uint, basePtr []ContextModel, count *int) *ContextModel3DBuffer {
    pContextModel3DBuffer := &ContextModel3DBuffer{m_sizeX: uiSizeX,
        m_sizeXY:       uiSizeX * uiSizeY,
        m_sizeXYZ:      uiSizeX * uiSizeY * uiSizeZ,
        m_contextModel: basePtr}
    *count += int(uiSizeX * uiSizeY * uiSizeZ)

    return pContextModel3DBuffer
}

// access functions
func (this *ContextModel3DBuffer) Get3(uiZ, uiY, uiX uint) *ContextModel {
    return &this.m_contextModel[uiZ*this.m_sizeXY+uiY*this.m_sizeX+uiX]
}

func (this *ContextModel3DBuffer) Get2(uiZ, uiY uint) []ContextModel {
    return this.m_contextModel[uiZ*this.m_sizeXY+uiY*this.m_sizeX:]
}

func (this *ContextModel3DBuffer) Get1(uiZ uint) []ContextModel {
    return this.m_contextModel[uiZ*this.m_sizeXY:]
}

// initialization & copy functions
func (this *ContextModel3DBuffer) InitBuffer(sliceType SliceType, qp int, ctxModel []byte) { ///< initialize 3D buffer by slice type & QP
    ctxModel = ctxModel[uint(sliceType)*this.m_sizeXYZ:]
    for n := uint(0); n < this.m_sizeXYZ; n++ {
        this.m_contextModel[n].Init(qp, int(ctxModel[n]))
        this.m_contextModel[n].SetBinsCoded(0)
    }
}

var aStateToProbLPS = []float64{0.50000000, 0.47460857, 0.45050660, 0.42762859, 0.40591239, 0.38529900, 0.36573242, 0.34715948, 0.32952974, 0.31279528, 0.29691064, 0.28183267, 0.26752040, 0.25393496, 0.24103941, 0.22879875, 0.21717969, 0.20615069, 0.19568177, 0.18574449, 0.17631186, 0.16735824, 0.15885931, 0.15079198, 0.14313433, 0.13586556, 0.12896592, 0.12241667, 0.11620000, 0.11029903, 0.10469773, 0.09938088, 0.09433404, 0.08954349, 0.08499621, 0.08067986, 0.07658271, 0.07269362, 0.06900203, 0.06549791, 0.06217174, 0.05901448, 0.05601756, 0.05317283, 0.05047256, 0.04790942, 0.04547644, 0.04316702, 0.04097487, 0.03889405, 0.03691890, 0.03504406, 0.03326442, 0.03157516, 0.02997168, 0.02844963, 0.02700488, 0.02563349, 0.02433175, 0.02309612, 0.02192323, 0.02080991, 0.01975312, 0.01875000}

func (this *ContextModel3DBuffer) CalcCost(sliceType SliceType, qp int, ctxModel []byte) uint { ///< determine cost of choosing a probability table based on current probabilities
    cost := uint(0)
	ctxModel = ctxModel[uint(sliceType)*this.m_sizeXYZ:]
    for n := uint(0); n < this.m_sizeXYZ; n++ {
        tmpContextModel := NewContextModel()
        tmpContextModel.Init(qp, int(ctxModel[n]))

        // Map the 64 CABAC states to their corresponding probability values

        probLPS := aStateToProbLPS[this.m_contextModel[n].GetState()]
        var prob0, prob1 float64
        if this.m_contextModel[n].GetMps() == 1 {
            prob0 = probLPS
            prob1 = 1.0 - prob0
        } else {
            prob1 = probLPS
            prob0 = 1.0 - prob1
        }

        if this.m_contextModel[n].GetBinsCoded() > 0 {
            cost += uint(prob0*float64(tmpContextModel.GetEntropyBits(0)) + prob1*float64(tmpContextModel.GetEntropyBits(1)))
        }
    }

    return cost
}

// copy from another buffer
// \param src buffer to copy from

func (this *ContextModel3DBuffer) CopyFrom(src *ContextModel3DBuffer) {
    //    assert( m_sizeXYZ == src->m_sizeXYZ );
    //    ::memcpy( m_contextModel, src->m_contextModel, sizeof(ContextModel) * m_sizeXYZ );
    for i := uint(0); i < this.m_sizeXYZ; i++ {
        this.m_contextModel[i] = src.m_contextModel[i]
    }
}
