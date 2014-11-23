package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

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
	if info.StartGarbage != 0 {
		fmt.Printf("(found %d bytes of non-audio, non-tag data before first valid frame)\n", info.StartGarbage)
	}
	if info.FoundMPEG2 || info.FoundMPEG25 ||
		info.FoundLayer2 || info.FoundLayer1 {
		fmt.Printf("(Found non-mp3 frame versions)\n")
	}

	tabWriter := tabwriter.NewWriter(os.Stdout, 0, 4, 0, '\t', 0)
	brType := "CBR"
	if info.IsVBR {
		brType = "VBR"
	}
	fmt.Fprintf(tabWriter, "Bitrate:\t%d %s\n", info.Bitrate, brType)
	fmt.Fprintf(tabWriter, "Samplerate:\t%d\n", info.Samplerate)
	if info.HasID3v1 {
		fmt.Fprintf(tabWriter, "ID3v1:\tFound\n")
	}
	if info.ID3v2 != nil {
		fmt.Fprintf(tabWriter, "ID3v2:\tFound (%db)\n", info.ID3v2.Size)
	}

	tabWriter.Flush()
}
