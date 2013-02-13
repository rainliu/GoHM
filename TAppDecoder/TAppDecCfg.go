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

package TAppDecoder

import (
    "container/list"
    "errors"
    "strconv"
)

type TAppDecCfg struct {
    m_pchBitstreamFile             string     ///< input bitstream file name
    m_pchReconFile                 string     ///< output reconstruction file name
    m_pchTraceFile                 string     ///< trace file name
    m_iFrameNum                    int        ///< output frame number
    m_iSkipFrame                   int        ///< counter for frames prior to the random access point to skip
    m_outputBitDepthY              int        ///< bit depth used for writing output (luma)
    m_outputBitDepthC              int        ///< bit depth used for writing output (chroma)t
    m_iMaxTemporalLayer            int        ///< maximum temporal layer to be decoded
    m_decodedPictureHashSEIEnabled int        ///< Checksum(3)/CRC(2)/MD5(1)/disable(0) acting on decoded picture hash SEI message
    m_targetDecLayerIdSet          *list.List ///< set of LayerIds to be included in the sub-bitstream extraction process.
    m_respectDefDispWindow         int        ///< Only output content inside the default display window
}

func NewTAppDecCfg() *TAppDecCfg {
    return &TAppDecCfg{}
}

func (this *TAppDecCfg) ParseCfg(argc int, argv []string) (err error) { ///< initialize option class from configuration
    if argc <= 3 {
        err = errors.New("Too few arguments for HM Decoder")
        return err
    }

    this.m_iFrameNum = -1
    this.m_iSkipFrame = 0
    this.m_outputBitDepthY = 0
    this.m_outputBitDepthC = 0
    this.m_iMaxTemporalLayer = -1
    this.m_decodedPictureHashSEIEnabled = 1
    //this.m_targetDecLayerIdSet = 0

    if argc >= 4 {
        this.m_pchBitstreamFile = argv[2]
        this.m_pchReconFile = argv[3]
    }

    if argc >= 5 {
        this.m_iFrameNum, err = strconv.Atoi(argv[4])
        if err != nil {
            return err
        }
    }

    if argc >= 6 {
        this.m_pchTraceFile = argv[5]
    }

    return nil
}
