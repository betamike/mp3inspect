package mp3

import (
    "testing"
)

func TestFileInspection(t *testing.T) {
/*
    info, err := InspectFile("sky.mp3")
    if err != nil {
        t.Fatal(err)
    }

    if info.Bitrate != 320000 {
        t.Fatalf("Found bitrate of %d should be 320000\n", info.Bitrate)
    }
    if info.IsVBR == true {
        t.Fatal("VBR reported but file is CBR")
    }
    if info.HasXingHeader == true {
        t.Fatal("Xing header found but file does not have one")
    }
    if info.HasVBRiHeader == true {
        t.Fatal("VBRi header found but file does not have one")
    }
    if info.Samplerate != 44100 {
        t.Fatalf("Found samplerate of %d should be 44100\n", info.Bitrate)
    }
    //TODO: add duration and metadata type checks
    */
}

func BenchmarkFileInspection(b *testing.B) {
    _, err := InspectFile("sky.mp3")
    if err != nil {
        b.Fatal(err)
    }
}
