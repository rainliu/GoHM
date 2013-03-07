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

package main

import (
    "fmt"
    "gohm/TAppDecoder"
    "gohm/TAppEncoder"
    "gohm/TLibCommon"
    "log"
    "os"
    "time"
    "runtime/pprof"
)

func Encoder() {
    cTAppEncTop := TAppEncoder.NewTAppEncTop()

    // create application encoder class
    cTAppEncTop.Create()

    // parse configuration
    if err := cTAppEncTop.ParseCfg(len(os.Args), os.Args); err != nil {
        log.Fatal(err)
        return
    }

    // starting time
    lBefore := time.Now()

    // call encoding function
    cTAppEncTop.Encode()

    // ending time
    lAfter := time.Now()

    fmt.Printf("\n\nTotal Encoding Time: %v.\n", lAfter.Sub(lBefore))

    // destroy application encoder class
    cTAppEncTop.Destroy()
}

func Decoder() {
    cTAppDecTop := TAppDecoder.NewTAppDecTop()

    // create application decoder class
    cTAppDecTop.Create()

    // parse configuration
    if err := cTAppDecTop.ParseCfg(len(os.Args), os.Args); err != nil {
        log.Fatal(err)
        return
    }

    // starting time
    lBefore := time.Now()

    // call decoding function
    cTAppDecTop.Decode()

    // ending time
    lAfter := time.Now()

    fmt.Printf("\n\nTotal Decoding Time: %v.\n", lAfter.Sub(lBefore))

    // destroy application decoder class
    cTAppDecTop.Destroy()
}

func main() {
	f, err := os.Create("cpuprofile.prof")
   	if err != nil {
    	log.Fatal(err)
    }
    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()
    
    fmt.Printf("GoHM Software Version [%s]\n", TLibCommon.NV_VERSION)
    if len(os.Args) <= 2 {
        fmt.Printf("Usage: \n")
        fmt.Printf("	HM Encoder: gohm.exe -c encoder.cfg [trace.txt]\n")
        fmt.Printf("	HM Decoder: gohm.exe -d test.bin test.yuv [n trace.txt]\n")
    } else {
        if os.Args[1] == "-c" {
            Encoder()
        } else if os.Args[1] == "-d" {
            Decoder()
        } else {
            fmt.Printf("Unknown argment %s\n", os.Args[1])
        }
    }
    
    m, err := os.Create("memprofile.prof")
    if err != nil {
    	log.Fatal(err)
    }
    pprof.WriteHeapProfile(m)
    m.Close()
}
