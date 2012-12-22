package main

import (
    "fmt"
    "gohm/TAppDecoder"
    "gohm/TLibCommon"
    "log"
    "os"
    "time"
)

func main() {
    fmt.Printf("GoHM Software: Decoder Version [%s]\n", TLibCommon.NV_VERSION)

    cTAppDecTop := TAppDecoder.NewTAppDecTop()

    // create application decoder class
    cTAppDecTop.Create()

    if true {
        var args =[]string{"gohm.exe", "test.bin", "test.yuv", "1", "trace.txt"};

        // parse configuration
        if err := cTAppDecTop.ParseCfg(len(args), args); err != nil {
            log.Fatal(err)
            return
        }
    }else{
    // parse configuration
    if err := cTAppDecTop.ParseCfg(len(os.Args), os.Args); err != nil {
        log.Fatal(err)
        return
    }
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
