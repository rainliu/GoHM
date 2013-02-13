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

package TLibEncoder

import (
    "container/list"
    "fmt"
    "gohm/TLibCommon"
    "math"
)

const g_RCInvalidQPValue = -999
const g_RCSmoothWindowSize = 40
const g_RCMaxPicListSize = 32
const g_RCWeightPicTargetBitInGOP = float64(0.9)
const g_RCWeightPicRargetBitInBuffer = float64(1.0) - g_RCWeightPicTargetBitInGOP

type TRCLCU struct {
    m_actualBits    int
    m_QP            int // QP of skip mode is set to g_RCInvalidQPValue
    m_targetBits    int
    m_lambda        float64
    m_MAD           float64
    m_numberOfPixel int
}

type TRCParameter struct {
    m_alpha float64
    m_beta  float64
}

type TEncRCSeq struct {
    m_totalFrames   int
    m_targetRate    int
    m_frameRate     int
    m_GOPSize       int
    m_picWidth      int
    m_picHeight     int
    m_LCUWidth      int
    m_LCUHeight     int
    m_numberOfLevel int
    m_averageBits   int

    m_numberOfPixel int
    m_targetBits    int64
    m_numberOfLCU   int
    m_bitsRatio     []int
    m_GOPID2Level   []int
    m_picPara       []TRCParameter
    m_LCUPara       [][]TRCParameter

    m_framesLeft          int
    m_bitsLeft            int64
    m_seqTargetBpp        float64
    m_alphaUpdate         float64
    m_betaUpdate          float64
    m_useLCUSeparateModel bool
}

func NewTEncRCSeq() *TEncRCSeq {
    return &TEncRCSeq{}
}

func (this *TEncRCSeq) create(totalFrames, targetBitrate, frameRate, GOPSize, picWidth, picHeight, LCUWidth, LCUHeight, numberOfLevel int, useLCUSeparateModel bool) {
    this.m_totalFrames = totalFrames
    this.m_targetRate = targetBitrate
    this.m_frameRate = frameRate
    this.m_GOPSize = GOPSize
    this.m_picWidth = picWidth
    this.m_picHeight = picHeight
    this.m_LCUWidth = LCUWidth
    this.m_LCUHeight = LCUHeight
    this.m_numberOfLevel = numberOfLevel
    this.m_useLCUSeparateModel = useLCUSeparateModel

    this.m_numberOfPixel = this.m_picWidth * this.m_picHeight
    this.m_targetBits = int64(this.m_totalFrames) * int64(this.m_targetRate) / int64(this.m_frameRate)
    this.m_seqTargetBpp = float64(this.m_targetRate) / float64(this.m_frameRate) / float64(this.m_numberOfPixel)
    if this.m_seqTargetBpp < 0.03 {
        this.m_alphaUpdate = 0.01
        this.m_betaUpdate = 0.005
    } else if this.m_seqTargetBpp < 0.08 {
        this.m_alphaUpdate = 0.05
        this.m_betaUpdate = 0.025
    } else {
        this.m_alphaUpdate = 0.1
        this.m_betaUpdate = 0.05
    }
    this.m_averageBits = int(this.m_targetBits / int64(totalFrames))

    var picWidthInBU, picHeightInBU int
    if (this.m_picWidth % this.m_LCUWidth) == 0 {
        picWidthInBU = this.m_picWidth / this.m_LCUWidth
    } else {
        picWidthInBU = this.m_picWidth/this.m_LCUWidth + 1
    }
    if (this.m_picHeight % this.m_LCUHeight) == 0 {
        picHeightInBU = this.m_picHeight / this.m_LCUHeight
    } else {
        picHeightInBU = this.m_picHeight/this.m_LCUHeight + 1
    }

    this.m_numberOfLCU = picWidthInBU * picHeightInBU

    this.m_bitsRatio = make([]int, this.m_GOPSize)
    for i := 0; i < this.m_GOPSize; i++ {
        this.m_bitsRatio[i] = 1
    }

    this.m_GOPID2Level = make([]int, this.m_GOPSize)
    for i := 0; i < this.m_GOPSize; i++ {
        this.m_GOPID2Level[i] = 1
    }

    this.m_picPara = make([]TRCParameter, this.m_numberOfLevel)
    for i := 0; i < this.m_numberOfLevel; i++ {
        this.m_picPara[i].m_alpha = 0.0
        this.m_picPara[i].m_beta = 0.0
    }

    if this.m_useLCUSeparateModel {
        this.m_LCUPara = make([][]TRCParameter, this.m_numberOfLevel)
        for i := 0; i < this.m_numberOfLevel; i++ {
            this.m_LCUPara[i] = make([]TRCParameter, this.m_numberOfLCU)
            for j := 0; j < this.m_numberOfLCU; j++ {
                this.m_LCUPara[i][j].m_alpha = 0.0
                this.m_LCUPara[i][j].m_beta = 0.0
            }
        }
    }

    this.m_framesLeft = this.m_totalFrames
    this.m_bitsLeft = this.m_targetBits
}
func (this *TEncRCSeq) destroy() {}
func (this *TEncRCSeq) initBitsRatio(bitsRatio []int) {
    for i := 0; i < this.m_GOPSize; i++ {
        this.m_bitsRatio[i] = bitsRatio[i]
    }
}

func (this *TEncRCSeq) initGOPID2Level(GOPID2Level []int) {
    for i := 0; i < this.m_GOPSize; i++ {
        this.m_GOPID2Level[i] = GOPID2Level[i]
    }
}

func (this *TEncRCSeq) initPicPara(picPara []TRCParameter) { // NULL to initial with default value
    //assert( this.m_picPara != NULL );

    if picPara == nil {
        for i := 0; i < this.m_numberOfLevel; i++ {
            this.m_picPara[i].m_alpha = 3.2003
            this.m_picPara[i].m_beta = -1.367
        }
    } else {
        for i := 0; i < this.m_numberOfLevel; i++ {
            this.m_picPara[i] = picPara[i]
        }
    }
}

func (this *TEncRCSeq) initLCUPara(LCUPara [][]TRCParameter) { // NULL to initial with default value
    if this.m_LCUPara == nil {
        return
    }
    if LCUPara == nil {
        for i := 0; i < this.m_numberOfLevel; i++ {
            for j := 0; j < this.m_numberOfLCU; j++ {
                this.m_LCUPara[i][j].m_alpha = 3.2003
                this.m_LCUPara[i][j].m_beta = -1.367
            }
        }
    } else {
        for i := 0; i < this.m_numberOfLevel; i++ {
            for j := 0; j < this.m_numberOfLCU; j++ {
                this.m_LCUPara[i][j] = LCUPara[i][j]
            }
        }
    }
}

func (this *TEncRCSeq) updateAfterPic(bits int64) {
    this.m_bitsLeft -= bits
    this.m_framesLeft--
}

func (this *TEncRCSeq) getRefineBitsForIntra(orgBits int) int {
    bpp := (float64(orgBits)) / float64(this.m_picHeight) / float64(this.m_picHeight)
    if bpp > 0.2 {
        return orgBits * 5
    }
    if bpp > 0.1 {
        return orgBits * 7
    }
    return orgBits * 10
}

func (this *TEncRCSeq) getTotalFrames() int   { return this.m_totalFrames }
func (this *TEncRCSeq) getTargetRate() int    { return this.m_targetRate }
func (this *TEncRCSeq) getFrameRate() int     { return this.m_frameRate }
func (this *TEncRCSeq) getGOPSize() int       { return this.m_GOPSize }
func (this *TEncRCSeq) getPicWidth() int      { return this.m_picWidth }
func (this *TEncRCSeq) getPicHeight() int     { return this.m_picHeight }
func (this *TEncRCSeq) getLCUWidth() int      { return this.m_LCUWidth }
func (this *TEncRCSeq) getLCUHeight() int     { return this.m_LCUHeight }
func (this *TEncRCSeq) getNumberOfLevel() int { return this.m_numberOfLevel }
func (this *TEncRCSeq) getAverageBits() int   { return this.m_averageBits }
func (this *TEncRCSeq) getLeftAverageBits() int {
    return int(this.m_bitsLeft / int64(this.m_framesLeft))
}
func (this *TEncRCSeq) getUseLCUSeparateModel() bool { return this.m_useLCUSeparateModel }

func (this *TEncRCSeq) getNumPixel() int                         { return this.m_numberOfPixel }
func (this *TEncRCSeq) getTargetBits() int64                     { return this.m_targetBits }
func (this *TEncRCSeq) getNumberOfLCU() int                      { return this.m_numberOfLCU }
func (this *TEncRCSeq) getBitRatio() []int                       { return this.m_bitsRatio }
func (this *TEncRCSeq) getBitRatio1(idx int) int                 { return this.m_bitsRatio[idx] }
func (this *TEncRCSeq) getGOPID2Level() []int                    { return this.m_GOPID2Level }
func (this *TEncRCSeq) getGOPID2Level1(ID int) int               { return this.m_GOPID2Level[ID] }
func (this *TEncRCSeq) getPicPara() []TRCParameter               { return this.m_picPara }
func (this *TEncRCSeq) getPicPara1(level int) TRCParameter       { return this.m_picPara[level] }
func (this *TEncRCSeq) setPicPara2(level int, para TRCParameter) { this.m_picPara[level] = para }
func (this *TEncRCSeq) getLCUPara() [][]TRCParameter             { return this.m_LCUPara }
func (this *TEncRCSeq) getLCUPara1(level int) []TRCParameter     { return this.m_LCUPara[level] }
func (this *TEncRCSeq) getLCUPara2(level, LCUIdx int) TRCParameter {
    return this.getLCUPara1(level)[LCUIdx]
}
func (this *TEncRCSeq) setLCUPara3(level, LCUIdx int, para TRCParameter) {
    this.m_LCUPara[level][LCUIdx] = para
}

func (this *TEncRCSeq) getFramesLeft() int { return this.m_framesLeft }
func (this *TEncRCSeq) getBitsLeft() int64 { return this.m_bitsLeft }

func (this *TEncRCSeq) getSeqBpp() float64      { return this.m_seqTargetBpp }
func (this *TEncRCSeq) getAlphaUpdate() float64 { return this.m_alphaUpdate }
func (this *TEncRCSeq) getBetaUpdate() float64  { return this.m_betaUpdate }

type TEncRCGOP struct {
    m_encRCSeq          *TEncRCSeq
    m_picTargetBitInGOP []int
    m_numPic            int
    m_targetBits        int
    m_picLeft           int
    m_bitsLeft          int
}

func NewTEncRCGOP() *TEncRCGOP {
    return &TEncRCGOP{}
}

func (this *TEncRCGOP) create(encRCSeq *TEncRCSeq, numPic int) {
    targetBits := this.xEstGOPTargetBits(encRCSeq, numPic)

    this.m_picTargetBitInGOP = make([]int, numPic)
    var i, totalPicRatio, currPicRatio int
    for i = 0; i < numPic; i++ {
        totalPicRatio += encRCSeq.getBitRatio1(i)
    }
    for i = 0; i < numPic; i++ {
        currPicRatio = encRCSeq.getBitRatio1(i)
        this.m_picTargetBitInGOP[i] = targetBits * currPicRatio / totalPicRatio
    }

    this.m_encRCSeq = encRCSeq
    this.m_numPic = numPic
    this.m_targetBits = targetBits
    this.m_picLeft = this.m_numPic
    this.m_bitsLeft = this.m_targetBits
}
func (this *TEncRCGOP) destroy() {}
func (this *TEncRCGOP) updateAfterPicture(bitsCost int) {
    this.m_bitsLeft -= bitsCost
    this.m_picLeft--
}

func (this *TEncRCGOP) xEstGOPTargetBits(encRCSeq *TEncRCSeq, GOPSize int) int {
    realInfluencePicture := TLibCommon.MIN(g_RCSmoothWindowSize, encRCSeq.getFramesLeft()).(int)
    averageTargetBitsPerPic := int(encRCSeq.getTargetBits() / int64(encRCSeq.getTotalFrames()))
    currentTargetBitsPerPic := int((int64(encRCSeq.getBitsLeft()) - int64(averageTargetBitsPerPic)*(int64(encRCSeq.getFramesLeft())-int64(realInfluencePicture))) / int64(realInfluencePicture))
    targetBits := currentTargetBitsPerPic * GOPSize

    if targetBits < 200 {
        targetBits = 200 // at least allocate 200 bits for one GOP
    }

    return targetBits
}

func (this *TEncRCGOP) getEncRCSeq() *TEncRCSeq     { return this.m_encRCSeq }
func (this *TEncRCGOP) getNumPic() int              { return this.m_numPic }
func (this *TEncRCGOP) getTargetBits() int          { return this.m_targetBits }
func (this *TEncRCGOP) getPicLeft() int             { return this.m_picLeft }
func (this *TEncRCGOP) getBitsLeft() int            { return this.m_bitsLeft }
func (this *TEncRCGOP) getTargetBitInGOP(i int) int { return this.m_picTargetBitInGOP[i] }

type TEncRCPic struct {
    m_encRCSeq *TEncRCSeq
    m_encRCGOP *TEncRCGOP

    m_frameLevel    int
    m_numberOfPixel int
    m_numberOfLCU   int
    m_targetBits    int
    m_estHeaderBits int
    m_estPicQP      int
    m_estPicLambda  float64

    m_LCULeft    int
    m_bitsLeft   int
    m_pixelsLeft int

    m_LCUs                []TRCLCU
    m_picActualHeaderBits int // only SH and potential APS
    m_totalMAD            float64
    m_picActualBits       int // the whole picture, including header
    m_picQP               int // in integer form
    m_picLambda           float64
    m_lastPicture         *TEncRCPic
}

func NewTEncRCPic() *TEncRCPic {
    return &TEncRCPic{}
}

func (this *TEncRCPic) create(encRCSeq *TEncRCSeq, encRCGOP *TEncRCGOP, frameLevel int, listPreviousPictures *list.List) {
    this.m_encRCSeq = encRCSeq
    this.m_encRCGOP = encRCGOP

    targetBits := this.xEstPicTargetBits(encRCSeq, encRCGOP)
    estHeaderBits := this.xEstPicHeaderBits(listPreviousPictures, frameLevel)

    if targetBits < estHeaderBits+100 {
        targetBits = estHeaderBits + 100 // at least allocate 100 bits for picture data
    }

    this.m_frameLevel = frameLevel
    this.m_numberOfPixel = encRCSeq.getNumPixel()
    this.m_numberOfLCU = encRCSeq.getNumberOfLCU()
    this.m_estPicLambda = 100.0
    this.m_targetBits = targetBits
    this.m_estHeaderBits = estHeaderBits
    this.m_bitsLeft = this.m_targetBits
    picWidth := encRCSeq.getPicWidth()
    picHeight := encRCSeq.getPicHeight()
    LCUWidth := encRCSeq.getLCUWidth()
    LCUHeight := encRCSeq.getLCUHeight()
    var picWidthInLCU, picHeightInLCU int
    if (picWidth % LCUWidth) == 0 {
        picWidthInLCU = picWidth / LCUWidth
    } else {
        picWidthInLCU = picWidth/LCUWidth + 1
    }
    if (picHeight % LCUHeight) == 0 {
        picHeightInLCU = picHeight / LCUHeight
    } else {
        picHeightInLCU = picHeight/LCUHeight + 1
    }

    this.m_LCULeft = this.m_numberOfLCU
    this.m_bitsLeft -= this.m_estHeaderBits
    this.m_pixelsLeft = this.m_numberOfPixel

    this.m_LCUs = make([]TRCLCU, this.m_numberOfLCU)
    var i, j, LCUIdx int
    for i = 0; i < picWidthInLCU; i++ {
        for j = 0; j < picHeightInLCU; j++ {
            LCUIdx = j*picWidthInLCU + i
            this.m_LCUs[LCUIdx].m_actualBits = 0
            this.m_LCUs[LCUIdx].m_QP = 0
            this.m_LCUs[LCUIdx].m_lambda = 0.0
            this.m_LCUs[LCUIdx].m_targetBits = 0
            this.m_LCUs[LCUIdx].m_MAD = 0.0
            var currWidth, currHeight int
            if i == picWidthInLCU-1 {
                currWidth = picWidth - LCUWidth*(picWidthInLCU-1)
            } else {
                currWidth = LCUWidth
            }
            if j == picHeightInLCU-1 {
                currHeight = picHeight - LCUHeight*(picHeightInLCU-1)
            } else {
                currHeight = LCUHeight
            }
            this.m_LCUs[LCUIdx].m_numberOfPixel = currWidth * currHeight
        }
    }
    this.m_picActualHeaderBits = 0
    this.m_totalMAD = 0.0
    this.m_picActualBits = 0
    this.m_picQP = 0
    this.m_picLambda = 0.0

    this.m_lastPicture = nil
    //list<TEncRCPic*>::reverse_iterator it;
    for it := listPreviousPictures.Front(); it != nil; it = it.Next() {
        v := it.Value.(*TEncRCPic)
        if v.getFrameLevel() == this.m_frameLevel {
            this.m_lastPicture = v
            break
        }
    }
}

func (this *TEncRCPic) destroy() {}

func (this *TEncRCPic) estimatePicLambda(listPreviousPictures *list.List) float64 {
    alpha := this.m_encRCSeq.getPicPara1(this.m_frameLevel).m_alpha
    beta := this.m_encRCSeq.getPicPara1(this.m_frameLevel).m_beta
    bpp := float64(this.m_targetBits) / float64(this.m_numberOfPixel)
    estLambda := alpha * math.Pow(bpp, beta)
    lastLevelLambda := float64(-1.0)
    lastPicLambda := float64(-1.0)
    lastValidLambda := float64(-1.0)
    //list<TEncRCPic*>::iterator it;
    for it := listPreviousPictures.Front(); it != nil; it = it.Next() {
        v := it.Value.(*TEncRCPic)
        if v.getFrameLevel() == this.m_frameLevel {
            lastLevelLambda = v.getPicActualLambda()
        }
        lastPicLambda = v.getPicActualLambda()

        if lastPicLambda > 0.0 {
            lastValidLambda = lastPicLambda
        }
    }

    if lastLevelLambda > 0.0 {
        lastLevelLambda = TLibCommon.CLIP3(float64(0.1), float64(10000.0), lastLevelLambda).(float64)
        estLambda = TLibCommon.CLIP3(lastLevelLambda*math.Pow(2.0, -3.0/3.0), lastLevelLambda*math.Pow(2.0, 3.0/3.0), estLambda).(float64)
    }

    if lastPicLambda > 0.0 {
        lastPicLambda = TLibCommon.CLIP3(0.1, 2000.0, lastPicLambda).(float64)
        estLambda = TLibCommon.CLIP3(lastPicLambda*math.Pow(2.0, -10.0/3.0), lastPicLambda*math.Pow(2.0, 10.0/3.0), estLambda).(float64)
    } else if lastValidLambda > 0.0 {
        lastValidLambda = TLibCommon.CLIP3(0.1, 2000.0, lastValidLambda).(float64)
        estLambda = TLibCommon.CLIP3(lastValidLambda*math.Pow(2.0, -10.0/3.0), lastValidLambda*math.Pow(2.0, 10.0/3.0), estLambda).(float64)
    } else {
        estLambda = TLibCommon.CLIP3(0.1, 10000.0, estLambda).(float64)
    }

    if estLambda < 0.1 {
        estLambda = 0.1
    }

    this.m_estPicLambda = estLambda
    return estLambda
}

func (this *TEncRCPic) estimatePicQP(lambda float64, listPreviousPictures *list.List) int {
    QP := int(4.2005*math.Log(lambda) + 13.7122 + 0.5)

    lastLevelQP := g_RCInvalidQPValue
    lastPicQP := g_RCInvalidQPValue
    lastValidQP := g_RCInvalidQPValue
    //list<TEncRCPic*>::iterator it;
    for it := listPreviousPictures.Front(); it != nil; it = it.Next() {
        v := it.Value.(*TEncRCPic)
        if v.getFrameLevel() == this.m_frameLevel {
            lastLevelQP = v.getPicActualQP()
        }
        lastPicQP = v.getPicActualQP()
        if lastPicQP > g_RCInvalidQPValue {
            lastValidQP = lastPicQP
        }
    }

    if lastLevelQP > g_RCInvalidQPValue {
        QP = TLibCommon.CLIP3(lastLevelQP-3, lastLevelQP+3, QP).(int)
    }

    if lastPicQP > g_RCInvalidQPValue {
        QP = TLibCommon.CLIP3(lastPicQP-10, lastPicQP+10, QP).(int)
    } else if lastValidQP > g_RCInvalidQPValue {
        QP = TLibCommon.CLIP3(lastValidQP-10, lastValidQP+10, QP).(int)
    }

    return QP
}

func (this *TEncRCPic) getLCUTargetBpp() float64 {
    LCUIdx := this.getLCUCoded()
    bpp := float64(-1.0)
    avgBits := 0
    totalMAD := float64(-1.0)
    MAD := float64(-1.0)

    if this.m_lastPicture == nil {
        avgBits = int(this.m_bitsLeft / this.m_LCULeft)
    } else {
        MAD = this.m_lastPicture.getLCU(LCUIdx).m_MAD
        totalMAD = this.m_lastPicture.getTotalMAD()
        for i := 0; i < LCUIdx; i++ {
            totalMAD -= this.m_lastPicture.getLCU(i).m_MAD
        }

        if totalMAD > 0.1 {
            avgBits = int(float64(this.m_bitsLeft) * MAD / totalMAD)
        } else {
            avgBits = int(this.m_bitsLeft / this.m_LCULeft)
        }
    }

    //#if L0033_RC_BUGFIX
    if avgBits < 1 {
        avgBits = 1
    }
    //#else
    //  if  avgBits < 5 {
    //    avgBits = 5;
    //  }
    //#endif

    bpp = float64(avgBits) / float64(this.m_LCUs[LCUIdx].m_numberOfPixel)
    this.m_LCUs[LCUIdx].m_targetBits = avgBits

    return bpp
}

func (this *TEncRCPic) getLCUEstLambda(bpp float64) float64 {
    LCUIdx := this.getLCUCoded()
    var alpha, beta float64
    if this.m_encRCSeq.getUseLCUSeparateModel() {
        alpha = this.m_encRCSeq.getLCUPara2(this.m_frameLevel, LCUIdx).m_alpha
        beta = this.m_encRCSeq.getLCUPara2(this.m_frameLevel, LCUIdx).m_beta
    } else {
        alpha = this.m_encRCSeq.getPicPara1(this.m_frameLevel).m_alpha
        beta = this.m_encRCSeq.getPicPara1(this.m_frameLevel).m_beta
    }

    estLambda := alpha * math.Pow(bpp, beta)
    //for Lambda clip, picture level clip
    clipPicLambda := this.m_estPicLambda

    //for Lambda clip, LCU level clip
    clipNeighbourLambda := float64(-1.0)
    for i := LCUIdx - 1; i >= 0; i-- {
        if this.m_LCUs[i].m_lambda > 0 {
            clipNeighbourLambda = this.m_LCUs[i].m_lambda
            break
        }
    }

    if clipNeighbourLambda > 0.0 {
        estLambda = TLibCommon.CLIP3(clipNeighbourLambda*math.Pow(2.0, -1.0/3.0), clipNeighbourLambda*math.Pow(2.0, 1.0/3.0), estLambda).(float64)
    }

    if clipPicLambda > 0.0 {
        estLambda = TLibCommon.CLIP3(clipPicLambda*math.Pow(2.0, -2.0/3.0), clipPicLambda*math.Pow(2.0, 2.0/3.0), estLambda).(float64)
    } else {
        estLambda = TLibCommon.CLIP3(10.0, 1000.0, estLambda).(float64)
    }

    if estLambda < 0.1 {
        estLambda = 0.1
    }

    return estLambda
}

func (this *TEncRCPic) getLCUEstQP(lambda float64, clipPicQP int) int {
    LCUIdx := int(this.getLCUCoded())
    estQP := int(4.2005*math.Log(lambda) + 13.7122 + 0.5)

    //for Lambda clip, LCU level clip
    clipNeighbourQP := g_RCInvalidQPValue
    //#if L0033_RC_BUGFIX
    for i := LCUIdx - 1; i >= 0; i-- {
        //#else
        //  for  i:=LCUIdx; i>=0; i-- {
        //#endif
        if (this.getLCU(i)).m_QP > g_RCInvalidQPValue {
            clipNeighbourQP = this.getLCU(i).m_QP
            break
        }
    }

    if clipNeighbourQP > g_RCInvalidQPValue {
        estQP = TLibCommon.CLIP3(clipNeighbourQP-1, clipNeighbourQP+1, estQP).(int)
    }

    estQP = TLibCommon.CLIP3(clipPicQP-2, clipPicQP+2, estQP).(int)

    return estQP
}

func (this *TEncRCPic) updateAfterLCU(LCUIdx, bits, QP int, lambda float64, updateLCUParameter bool) {
    this.m_LCUs[LCUIdx].m_actualBits = bits
    this.m_LCUs[LCUIdx].m_QP = QP
    this.m_LCUs[LCUIdx].m_lambda = lambda

    this.m_LCULeft--
    this.m_bitsLeft -= bits
    this.m_pixelsLeft -= this.m_LCUs[LCUIdx].m_numberOfPixel

    if !updateLCUParameter {
        return
    }

    if !this.m_encRCSeq.getUseLCUSeparateModel() {
        return
    }

    alpha := this.m_encRCSeq.getLCUPara2(this.m_frameLevel, LCUIdx).m_alpha
    beta := this.m_encRCSeq.getLCUPara2(this.m_frameLevel, LCUIdx).m_beta

    LCUActualBits := this.m_LCUs[LCUIdx].m_actualBits
    LCUTotalPixels := this.m_LCUs[LCUIdx].m_numberOfPixel
    bpp := float64(LCUActualBits) / float64(LCUTotalPixels)
    calLambda := alpha * math.Pow(bpp, beta)
    inputLambda := this.m_LCUs[LCUIdx].m_lambda

    if inputLambda < 0.01 || calLambda < 0.01 || bpp < 0.0001 {
        alpha *= (1.0 - this.m_encRCSeq.getAlphaUpdate()/2.0)
        beta *= (1.0 - this.m_encRCSeq.getBetaUpdate()/2.0)

        alpha = TLibCommon.CLIP3(0.05, 20.0, alpha).(float64)
        beta = TLibCommon.CLIP3(-3.0, -0.1, beta).(float64)

        var rcPara TRCParameter
        rcPara.m_alpha = alpha
        rcPara.m_beta = beta
        this.m_encRCSeq.setLCUPara3(this.m_frameLevel, LCUIdx, rcPara)

        return
    }

    calLambda = TLibCommon.CLIP3(inputLambda/10.0, inputLambda*10.0, calLambda).(float64)
    alpha += this.m_encRCSeq.getAlphaUpdate() * (math.Log(inputLambda) - math.Log(calLambda)) * alpha
    lnbpp := math.Log(bpp)
    lnbpp = TLibCommon.CLIP3(-5.0, 1.0, lnbpp).(float64)
    beta += this.m_encRCSeq.getBetaUpdate() * (math.Log(inputLambda) - math.Log(calLambda)) * lnbpp

    alpha = TLibCommon.CLIP3(0.05, 20.0, alpha).(float64)
    beta = TLibCommon.CLIP3(-3.0, -0.1, beta).(float64)
    var rcPara TRCParameter
    rcPara.m_alpha = alpha
    rcPara.m_beta = beta
    this.m_encRCSeq.setLCUPara3(this.m_frameLevel, LCUIdx, rcPara)
}

func (this *TEncRCPic) updateAfterPicture(actualHeaderBits, actualTotalBits int, averageQP, averageLambda, effectivePercentage float64) {
    this.m_picActualHeaderBits = actualHeaderBits
    this.m_picActualBits = actualTotalBits
    if averageQP > 0.0 {
        this.m_picQP = int(averageQP + 0.5)
    } else {
        this.m_picQP = g_RCInvalidQPValue
    }
    this.m_picLambda = averageLambda
    for i := 0; i < this.m_numberOfLCU; i++ {
        this.m_totalMAD += this.m_LCUs[i].m_MAD
    }

    alpha := this.m_encRCSeq.getPicPara1(this.m_frameLevel).m_alpha
    beta := this.m_encRCSeq.getPicPara1(this.m_frameLevel).m_beta

    // update parameters
    picActualBits := float64(this.m_picActualBits)
    picActualBpp := picActualBits / float64(this.m_numberOfPixel)
    calLambda := alpha * math.Pow(picActualBpp, beta)
    inputLambda := this.m_picLambda

    if inputLambda < 0.01 || calLambda < 0.01 || picActualBpp < 0.0001 || effectivePercentage < 0.05 {
        alpha *= (1.0 - this.m_encRCSeq.getAlphaUpdate()/2.0)
        beta *= (1.0 - this.m_encRCSeq.getBetaUpdate()/2.0)

        alpha = TLibCommon.CLIP3(0.05, 20.0, alpha).(float64)
        beta = TLibCommon.CLIP3(-3.0, -0.1, beta).(float64)
        var rcPara TRCParameter
        rcPara.m_alpha = alpha
        rcPara.m_beta = beta
        this.m_encRCSeq.setPicPara2(this.m_frameLevel, rcPara)

        return
    }

    calLambda = TLibCommon.CLIP3(inputLambda/10.0, inputLambda*10.0, calLambda).(float64)
    alpha += this.m_encRCSeq.getAlphaUpdate() * (math.Log(inputLambda) - math.Log(calLambda)) * alpha
    lnbpp := math.Log(picActualBpp)
    lnbpp = TLibCommon.CLIP3(-5.0, 1.0, lnbpp).(float64)
    beta += this.m_encRCSeq.getBetaUpdate() * (math.Log(inputLambda) - math.Log(calLambda)) * lnbpp

    alpha = TLibCommon.CLIP3(0.05, 20.0, alpha).(float64)
    beta = TLibCommon.CLIP3(-3.0, -0.1, beta).(float64)

    var rcPara TRCParameter
    rcPara.m_alpha = alpha
    rcPara.m_beta = beta

    this.m_encRCSeq.setPicPara2(this.m_frameLevel, rcPara)
}

func (this *TEncRCPic) addToPictureLsit(listPreviousPictures *list.List) {
    if listPreviousPictures.Len() > g_RCMaxPicListSize {
        p := listPreviousPictures.Front() //.Value.(*TEncRCPic);
        listPreviousPictures.Remove(p)    //pop_front();
        //p.destroy();
        //delete p;
    }

    listPreviousPictures.PushBack(this)
}

func (this *TEncRCPic) getEffectivePercentage() float64 {
    effectivePiexels := 0
    totalPixels := 0

    for i := 0; i < this.m_numberOfLCU; i++ {
        totalPixels += this.m_LCUs[i].m_numberOfPixel
        if this.m_LCUs[i].m_QP > 0 {
            effectivePiexels += this.m_LCUs[i].m_numberOfPixel
        }
    }

    effectivePixelPercentage := float64(effectivePiexels) / float64(totalPixels)
    return effectivePixelPercentage
}

func (this *TEncRCPic) calAverageQP() float64 {
    totalQPs := 0
    numTotalLCUs := 0

    var i int
    for i = 0; i < this.m_numberOfLCU; i++ {
        if this.m_LCUs[i].m_QP > 0 {
            totalQPs += this.m_LCUs[i].m_QP
            numTotalLCUs++
        }
    }

    avgQP := float64(0.0)
    if numTotalLCUs == 0 {
        avgQP = g_RCInvalidQPValue
    } else {
        avgQP = float64(totalQPs) / float64(numTotalLCUs)
    }
    return avgQP
}

func (this *TEncRCPic) calAverageLambda() float64 {
    totalLambdas := float64(0.0)
    numTotalLCUs := 0

    var i int
    for i = 0; i < this.m_numberOfLCU; i++ {
        if this.m_LCUs[i].m_lambda > 0.01 {
            totalLambdas += math.Log(this.m_LCUs[i].m_lambda)
            numTotalLCUs++
        }
    }

    var avgLambda float64
    if numTotalLCUs == 0 {
        avgLambda = -1.0
    } else {
        avgLambda = math.Pow(2.7183, totalLambdas/float64(numTotalLCUs))
    }
    return avgLambda
}

func (this *TEncRCPic) xEstPicTargetBits(encRCSeq *TEncRCSeq, encRCGOP *TEncRCGOP) int {
    targetBits := 0
    GOPbitsLeft := encRCGOP.getBitsLeft()

    var i int
    currPicPosition := encRCGOP.getNumPic() - encRCGOP.getPicLeft()
    currPicRatio := encRCSeq.getBitRatio1(currPicPosition)
    totalPicRatio := 0
    for i = currPicPosition; i < encRCGOP.getNumPic(); i++ {
        totalPicRatio += encRCSeq.getBitRatio1(i)
    }

    targetBits = int(GOPbitsLeft * currPicRatio / totalPicRatio)

    if targetBits < 100 {
        targetBits = 100 // at least allocate 100 bits for one picture
    }

    if this.m_encRCSeq.getFramesLeft() > 16 {
        targetBits = int(g_RCWeightPicRargetBitInBuffer*float64(targetBits) + g_RCWeightPicTargetBitInGOP*float64(this.m_encRCGOP.getTargetBitInGOP(currPicPosition)))
    }

    return targetBits
}

func (this *TEncRCPic) xEstPicHeaderBits(listPreviousPictures *list.List, frameLevel int) int {
    numPreviousPics := 0
    totalPreviousBits := 0

    //list<TEncRCPic*>::iterator it;
    for it := listPreviousPictures.Front(); it != nil; it = it.Next() {
        v := it.Value.(*TEncRCPic)
        if v.getFrameLevel() == frameLevel {
            totalPreviousBits += v.getPicActualHeaderBits()
            numPreviousPics++
        }
    }

    estHeaderBits := 0
    if numPreviousPics > 0 {
        estHeaderBits = totalPreviousBits / numPreviousPics
    }

    return estHeaderBits
}

func (this *TEncRCPic) getRCSequence() *TEncRCSeq { return this.m_encRCSeq }
func (this *TEncRCPic) getRCGOP() *TEncRCGOP      { return this.m_encRCGOP }

func (this *TEncRCPic) getFrameLevel() int {
    return this.m_frameLevel
}
func (this *TEncRCPic) getNumberOfPixel() int {
    return this.m_numberOfPixel
}
func (this *TEncRCPic) getNumberOfLCU() int {
    return this.m_numberOfLCU
}
func (this *TEncRCPic) getTargetBits() int {
    return this.m_targetBits
}
func (this *TEncRCPic) setTargetBits(bits int) { this.m_targetBits = bits }
func (this *TEncRCPic) getEstHeaderBits() int {
    return this.m_estHeaderBits
}
func (this *TEncRCPic) getLCULeft() int { return this.m_LCULeft }
func (this *TEncRCPic) getBitsLeft() int {
    return this.m_bitsLeft
}
func (this *TEncRCPic) getPixelsLeft() int {
    return this.m_pixelsLeft
}
func (this *TEncRCPic) getBitsCoded() int {
    return this.m_targetBits - this.m_estHeaderBits - this.m_bitsLeft
}
func (this *TEncRCPic) getLCUCoded() int {
    return this.m_numberOfLCU - this.m_LCULeft
}
func (this *TEncRCPic) getLCUs() []TRCLCU { return this.m_LCUs }
func (this *TEncRCPic) getLCU(LCUIdx int) *TRCLCU {
    return &this.m_LCUs[LCUIdx]
}
func (this *TEncRCPic) getPicActualHeaderBits() int {
    return this.m_picActualHeaderBits
}
func (this *TEncRCPic) getTotalMAD() float64 {
    return this.m_totalMAD
}
func (this *TEncRCPic) setTotalMAD(MAD float64) { this.m_totalMAD = MAD }
func (this *TEncRCPic) getPicActualBits() int {
    return this.m_picActualBits
}
func (this *TEncRCPic) getPicActualQP() int { return this.m_picQP }
func (this *TEncRCPic) getPicActualLambda() float64 {
    return this.m_picLambda
}
func (this *TEncRCPic) getPicEstQP() int {
    return this.m_estPicQP
}
func (this *TEncRCPic) setPicEstQP(QP int) { this.m_estPicQP = QP }
func (this *TEncRCPic) getPicEstLambda() float64 {
    return this.m_estPicLambda
}
func (this *TEncRCPic) setPicEstLambda(lambda float64) { this.m_picLambda = lambda }

type TEncRateCtrl struct {
    m_encRCSeq       *TEncRCSeq
    m_encRCGOP       *TEncRCGOP
    m_encRCPic       *TEncRCPic
    m_listRCPictures *list.List //<TEncRCPic*>;
    m_RCQP           int
}

func NewTEncRateCtrl() *TEncRateCtrl {
    return &TEncRateCtrl{}
}

func (this *TEncRateCtrl) init(totalFrames, targetBitrate, frameRate, GOPSize, picWidth, picHeight, LCUWidth, LCUHeight int, keepHierBits, useLCUSeparateModel bool, GOPList [TLibCommon.MAX_GOP]*GOPEntry) {
    isLowdelay := true
    for i := 0; i < GOPSize-1; i++ {
        if GOPList[i].m_POC > GOPList[i+1].m_POC {
            isLowdelay = false
            break
        }
    }

    numberOfLevel := 1
    if keepHierBits {
        numberOfLevel = int(math.Log(float64(GOPSize))/math.Log(2.0)+0.5) + 1
    }
    if !isLowdelay && GOPSize == 8 {
        numberOfLevel = int(math.Log(float64(GOPSize))/math.Log(2.0)+0.5) + 1
    }
    numberOfLevel++ // intra picture
    numberOfLevel++ // non-reference picture

    var bitsRatio []int
    bitsRatio = make([]int, GOPSize)
    for i := 0; i < GOPSize; i++ {
        bitsRatio[i] = 10
        if !GOPList[i].m_refPic {
            bitsRatio[i] = 2
        }
    }
    if keepHierBits {
        bpp := float64(targetBitrate) / float64(frameRate*picWidth*picHeight)
        if GOPSize == 4 && isLowdelay {
            if bpp > 0.2 {
                bitsRatio[0] = 2
                bitsRatio[1] = 3
                bitsRatio[2] = 2
                bitsRatio[3] = 6
            } else if bpp > 0.1 {
                bitsRatio[0] = 2
                bitsRatio[1] = 3
                bitsRatio[2] = 2
                bitsRatio[3] = 10
            } else if bpp > 0.05 {
                bitsRatio[0] = 2
                bitsRatio[1] = 3
                bitsRatio[2] = 2
                bitsRatio[3] = 12
            } else {
                bitsRatio[0] = 2
                bitsRatio[1] = 3
                bitsRatio[2] = 2
                bitsRatio[3] = 14
            }
        } else if GOPSize == 8 && !isLowdelay {
            if bpp > 0.2 {
                bitsRatio[0] = 15
                bitsRatio[1] = 5
                bitsRatio[2] = 4
                bitsRatio[3] = 1
                bitsRatio[4] = 1
                bitsRatio[5] = 4
                bitsRatio[6] = 1
                bitsRatio[7] = 1
            } else if bpp > 0.1 {
                bitsRatio[0] = 20
                bitsRatio[1] = 6
                bitsRatio[2] = 4
                bitsRatio[3] = 1
                bitsRatio[4] = 1
                bitsRatio[5] = 4
                bitsRatio[6] = 1
                bitsRatio[7] = 1
            } else if bpp > 0.05 {
                bitsRatio[0] = 25
                bitsRatio[1] = 7
                bitsRatio[2] = 4
                bitsRatio[3] = 1
                bitsRatio[4] = 1
                bitsRatio[5] = 4
                bitsRatio[6] = 1
                bitsRatio[7] = 1
            } else {
                bitsRatio[0] = 30
                bitsRatio[1] = 8
                bitsRatio[2] = 4
                bitsRatio[3] = 1
                bitsRatio[4] = 1
                bitsRatio[5] = 4
                bitsRatio[6] = 1
                bitsRatio[7] = 1
            }
        } else {
            fmt.Printf("\n hierarchical bit allocation is not support for the specified coding structure currently.")
        }
    }

    GOPID2Level := make([]int, GOPSize)
    for i := 0; i < GOPSize; i++ {
        GOPID2Level[i] = 1
        if !GOPList[i].m_refPic {
            GOPID2Level[i] = 2
        }
    }
    if keepHierBits {
        if GOPSize == 4 && isLowdelay {
            GOPID2Level[0] = 3
            GOPID2Level[1] = 2
            GOPID2Level[2] = 3
            GOPID2Level[3] = 1
        } else if GOPSize == 8 && !isLowdelay {
            GOPID2Level[0] = 1
            GOPID2Level[1] = 2
            GOPID2Level[2] = 3
            GOPID2Level[3] = 4
            GOPID2Level[4] = 4
            GOPID2Level[5] = 3
            GOPID2Level[6] = 4
            GOPID2Level[7] = 4
        }
    }

    if !isLowdelay && GOPSize == 8 {
        GOPID2Level[0] = 1
        GOPID2Level[1] = 2
        GOPID2Level[2] = 3
        GOPID2Level[3] = 4
        GOPID2Level[4] = 4
        GOPID2Level[5] = 3
        GOPID2Level[6] = 4
        GOPID2Level[7] = 4
    }

    this.m_encRCSeq = NewTEncRCSeq()
    this.m_encRCSeq.create(totalFrames, targetBitrate, frameRate, GOPSize, picWidth, picHeight, LCUWidth, LCUHeight, numberOfLevel, useLCUSeparateModel)
    this.m_encRCSeq.initBitsRatio(bitsRatio)
    this.m_encRCSeq.initGOPID2Level(GOPID2Level)
    this.m_encRCSeq.initPicPara(nil)
    if useLCUSeparateModel {
        this.m_encRCSeq.initLCUPara(nil)
    }

    //delete[] bitsRatio;
    //delete[] GOPID2Level;
}

func (this *TEncRateCtrl) destroy() {}
func (this *TEncRateCtrl) initRCPic(frameLevel int) {
    this.m_encRCPic = NewTEncRCPic()
    this.m_encRCPic.create(this.m_encRCSeq, this.m_encRCGOP, frameLevel, this.m_listRCPictures)
}

func (this *TEncRateCtrl) initRCGOP(numberOfPictures int) {
    this.m_encRCGOP = NewTEncRCGOP()
    this.m_encRCGOP.create(this.m_encRCSeq, numberOfPictures)
}

func (this *TEncRateCtrl) destroyRCGOP() {}

func (this *TEncRateCtrl) setRCQP(QP int)         { this.m_RCQP = QP }
func (this *TEncRateCtrl) getRCQP() int           { return this.m_RCQP }
func (this *TEncRateCtrl) getRCSeq() *TEncRCSeq   { return this.m_encRCSeq }
func (this *TEncRateCtrl) getRCGOP() *TEncRCGOP   { return this.m_encRCGOP }
func (this *TEncRateCtrl) getRCPic() *TEncRCPic   { return this.m_encRCPic }
func (this *TEncRateCtrl) getPicList() *list.List { return this.m_listRCPictures }
