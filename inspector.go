package main

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/betamike/mp3inspect/mp3"
)

func main() {
	var pretty bool
	var skip bool

	flag.BoolVar(&skip, "s", false, "skip printing the output headers")
	flag.BoolVar(&pretty, "p", false, "print in a pretty human readable format")

	flag.Usage = usage
	flag.Parse()

	if len(flag.Args()) < 1 {
		fatal("missing path to file")
	}

	info, err := mp3.InspectFile(flag.Args()[0])
	if err != nil {
		fatal("MP3 inspection failed: %s", err)
	}
	if pretty {
		prettyPrintInfo(info)
	} else {
		printInfo(info, !skip)
	}
}

func fatal(format string, params ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", params...)
	os.Stderr.Sync()
	os.Exit(-1)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: mp3inspect [flags] <path>\n")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "mp3inspect a simple utility for determining information about an mp3")
}

func printInfo(info *mp3.MP3Info, header bool) {
	tabWriter := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0)
	if header {
		fmt.Fprintln(tabWriter, "num frames\tbr(bps)\tvbr\tsr(hz)\tID3v1\tID3v2\tv2 size(b)\tbad bytes\tbad frames")
	}
	fmt.Fprintf(tabWriter, "%d\t%d\t%t\t%d\t%t\t%t\t", info.FrameCount, info.Bitrate, info.IsVBR, info.Samplerate, info.HasID3v1, (info.ID3v2 != nil))
	if info.ID3v2 == nil {
		fmt.Fprint(tabWriter, "0\t")
	} else {
		fmt.Fprintf(tabWriter, "%d\t", info.ID3v2.Size)
	}
	fmt.Fprintf(tabWriter, "%d\t", info.StartGarbage)
	if info.FoundMPEG2 || info.FoundMPEG25 ||
		info.FoundLayer2 || info.FoundLayer1 {
		fmt.Fprint(tabWriter, "1\t")
	} else {
		fmt.Fprint(tabWriter, "0\t")
	}
	fmt.Fprint(tabWriter, "\n")
	tabWriter.Flush()
}

func prettyPrintInfo(info *mp3.MP3Info) {
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
	tabWriter := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0)
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
