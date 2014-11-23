package mp3

import (
	"time"
)

type MPEGVersion int

const (
	MPEG25 MPEGVersion = iota
	MPEG_RESERVED
	MPEG2
	MPEG1
)

type MPEGLayer int

const (
	LAYER_RESERVED MPEGLayer = iota
	LAYER3
	LAYER2
	LAYER1
)

const (
	STEREO         = iota
	JOINT_STEREO   // stereo
	DUAL_CHANNE    // two mono chans
	SINGLE_CHANNEL // mono
)

var ZeroedBitrates = []uint64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

var MPEG1BitrateTable = [][]uint64{
	ZeroedBitrates, //reserved layer value
	{0, 32, 40, 48, 56, 64, 80, 96, 112, 128, 160, 192, 224, 256, 320, 0},     //layer 3
	{0, 32, 48, 56, 64, 80, 96, 112, 128, 160, 192, 224, 256, 320, 384, 0},    //layer 2
	{0, 32, 64, 96, 128, 160, 192, 224, 256, 288, 320, 352, 384, 416, 448, 0}, //layer 1
}
var MPEG2L1Bitrates = []uint64{0, 32, 48, 56, 64, 80, 96, 112, 128, 144, 160, 176, 192, 224, 256, 0}
var MPEG2L2L3Bitrates = []uint64{0, 8, 16, 24, 32, 40, 48, 56, 64, 80, 96, 112, 128, 144, 160, 0}
var MPEG2BitrateTable = [][]uint64{
	ZeroedBitrates,    //reserved layer value
	MPEG2L2L3Bitrates, //layer 3
	MPEG2L2L3Bitrates, //layer 2
	MPEG2L1Bitrates,   //layer 1
}

//Lookup by Version, Layer, and Bitrate Index
var BitrateLookupTable = [][][]uint64{
	MPEG2BitrateTable,                                                //MPEG 2.5
	{ZeroedBitrates, ZeroedBitrates, ZeroedBitrates, ZeroedBitrates}, //reserved mpeg value
	MPEG2BitrateTable,                                                //MPEG 2
	MPEG1BitrateTable,                                                //MPEG 1
}

//Lookup by Version, Samplerate Index
var SamplerateLookupTable = [][]uint64{
	{11025, 12000, 8000, 0},
	{0, 0, 0, 0},
	{22050, 24000, 16000, 0},
	{44100, 48000, 32000, 0},
}

//Lookup by Version, Layer
var SamplesLookupTable = [][]uint64{
	{0, 576, 1152, 384},
	{0, 0, 0, 0},
	{0, 576, 1152, 384},
	{0, 1152, 1152, 384},
}

type ApeTag struct {
	Version  uint32
	Size     uint32
	NumItems uint32
	Flags    uint32
}

type MP3Info struct {
	Bitrate       uint64
	IsVBR         bool
	HasXingHeader bool
	HasVBRiHeader bool

	Samplerate uint64

	Duration time.Duration

	HasID3v1 bool
	ID3v2    *ID3v2Tag
	Ape      *ApeTag

	FoundLayer1 bool
	FoundLayer2 bool
	FoundLayer3 bool

	FoundMPEG1  bool
	FoundMPEG2  bool
	FoundMPEG25 bool

	UnknownBytes uint64

	FrameCount   int
	StartGarbage int64
}

//information about one mp3 audio frame
type AudioFrame struct {
	Version         MPEGVersion
	Layer           MPEGLayer
	CRC             bool
	BitrateIndex    uint8
	SamplerateIndex uint8
	Padding         uint8
	Private         bool
	Mode            uint8
	ModeExt         uint8
	Copyright       bool
	Original        bool
	Emphasis        uint8

	Bitrate    uint64
	Samplerate uint64

	Size uint64
}

func parseAudioFrame(buf []byte) (int64, *AudioFrame) {
	frame := &AudioFrame{}
	frame.Version = MPEGVersion(buf[1] & 0x18 >> 3)
	frame.Layer = MPEGLayer(buf[1] & 0x6 >> 1)
	frame.CRC = (buf[1] & 0x1) != 1
	frame.BitrateIndex = buf[2] & 0xF0 >> 4
	frame.SamplerateIndex = buf[2] & 0xC >> 2
	frame.Padding = buf[2] & 0x2 >> 1
	frame.Private = (buf[2] & 0x1) == 1
	frame.Mode = buf[3] & 0xC0 >> 6
	frame.ModeExt = buf[3] & 0x30 >> 4
	frame.Copyright = (buf[3] & 0x8 >> 3) == 1
	frame.Original = (buf[3] & 0x4 >> 2) == 1
	frame.Emphasis = buf[3] & 0x3

	frame.Bitrate = BitrateLookupTable[frame.Version][frame.Layer][frame.BitrateIndex] * 1000
	frame.Samplerate = SamplerateLookupTable[frame.Version][frame.SamplerateIndex]

	//typically if any of these cases are met, this is just random audio
	//data that kinda looks like an audio header
	if frame.Version == MPEG_RESERVED || frame.Layer == LAYER_RESERVED ||
		frame.Samplerate == 0 || frame.Bitrate == 0 {
		return 0, nil
	}

	samples := SamplesLookupTable[frame.Version][frame.Layer]
	frame.Size = (((samples / 8) * frame.Bitrate) / frame.Samplerate) + uint64(frame.Padding)

	return int64(frame.Size), frame
}
