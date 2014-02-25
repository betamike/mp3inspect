package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/betamike/mp3inspect/mp3"
)

func main() {
	flag.Parse()
	path := flag.Arg(0)

	if path == "" {
		log.Fatal("Missing file path!")
	}

	info, err := mp3.InspectFile(path)
	if err != nil {
		fmt.Printf("MP3 inspection failed: %s", err)
		return
	}

	printInfo(info)
}

func printInfo(info *mp3.MP3Info) {
	if info.FoundMPEG1 && info.FoundLayer3 {
		fmt.Printf("MPEG v1 Layer III (Inspected %d frames)\n", info.FrameCount)
	}
	if info.FoundMPEG2 || info.FoundMPEG25 ||
		info.FoundLayer2 || info.FoundLayer1 {
		fmt.Printf("(Found non-mp3 frame versions)\n")
	}
	brType := "CBR"
	if info.IsVBR {
		brType = "VBR"
	}
	fmt.Printf("Bitrate:    %d %s\n", info.Bitrate, brType)
	fmt.Printf("Samplerate: %d\n", info.Samplerate)
	if info.HasID3v1 {
		fmt.Printf("ID3v1:      Found\n")
	}
	if info.ID3v2 != nil {
		fmt.Printf("ID3v2:      Found (%db)\n", info.ID3v2.Size)
	}
}
