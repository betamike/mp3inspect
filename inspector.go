package main

import (
    "mp3inspect/mp3"
    "flag"
    "fmt"
//    "runtime/pprof"
//    "os"
    "log"
)

func main() {
    flag.Parse()
    path := flag.Arg(0)

    if path == "" {
        log.Fatal("Missing file path!")
    }

    /*
    f, err := os.Create("cpu.out")
    if err != nil {
        log.Fatal(err)
    }
    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()
    */

    info, err := mp3.InspectFile(path)
    if err != nil {
        fmt.Printf("MP3 inspection failed: %s", err)
        return
    }

    fmt.Printf("%+v\n", info)
}
