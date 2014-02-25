package mp3

import (
	"bytes"
	"os"
)

type Scanner struct {
	f *os.File

	buf                 []byte
	FrameCount, curSize int
	vbrCounter          uint64
	r, seek             int64

	Info *MP3Info
}

func NewScanner(path string) (*Scanner, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return &Scanner{f: f, buf: make([]byte, 4096), Info: &MP3Info{}}, nil
}

func (s *Scanner) NextFrame() (*AudioFrame, error) {
	var frame *AudioFrame
	var err error
	for frame == nil {
		if s.seek > 0 {
			if (s.r + s.seek) > int64(s.curSize) {
				//remove already buffered data from seek offset
				s.seek -= (int64(s.curSize) - s.r)
				_, err = s.f.Seek(s.seek, os.SEEK_CUR)
				if err != nil {
					return nil, err
				}

				s.r = 0
				s.curSize, err = s.f.Read(s.buf)
				if s.curSize < 0 || err != nil {
					return nil, err
				}
			} else {
				s.r += s.seek
			}
		}
		s.seek = 0

		if (s.r + 10) > int64(s.curSize) {
			rem := int64(s.curSize) - s.r
			if rem > 0 {
				copy(s.buf[:rem], s.buf[s.r:])
			} else {
				rem = 0
			}

			s.r = 0
			s.curSize, err = s.f.Read(s.buf[rem:])
			if s.curSize < 0 || err != nil {
				return nil, err
			}
		}

		//fmt.Printf("%d : %d == %d\n", r, 4, r+4)
		cur := s.buf[s.r : s.r+4]
		s.r += 4

		switch {
		//potentially an audio frame
		case cur[0] == 0xFF && cur[1]&0xE0 == 0xE0:
			s.seek, frame = parseAudioFrame(cur)
			if frame == nil {
				//fmt.Printf("Bad potential frame\n")
				s.seek = 0
				s.r -= 3
				break
			}

			s.FrameCount++
			if s.Info.Bitrate != frame.Bitrate {
				if s.Info.Bitrate > 0 {
					s.Info.IsVBR = true
				} else {
					s.Info.Bitrate = frame.Bitrate
				}
			}
			s.vbrCounter += frame.Bitrate

			if s.Info.Samplerate != frame.Samplerate {
				s.Info.Samplerate = frame.Samplerate
			}

			switch frame.Version {
			case MPEG1:
				s.Info.FoundMPEG1 = true
			case MPEG2:
				s.Info.FoundMPEG2 = true
			case MPEG25:
				s.Info.FoundMPEG25 = true
			}

			switch frame.Layer {
			case LAYER1:
				s.Info.FoundLayer1 = true
			case LAYER2:
				s.Info.FoundLayer2 = true
			case LAYER3:
				s.Info.FoundLayer3 = true
			}

		//potentially ID3v1 tags
		case bytes.Equal(cur[0:3], ID3v1Header):
			s.seek = 127

		//potentially ID3v2 tags
		case bytes.Equal(cur[0:3], ID3v2Header):
			s.seek, s.Info.ID3v2 = parseID3v2Tag(s.buf[s.r-1 : s.r+7])
			s.r += 6

		//potentially APE tags
		case bytes.Equal(s.buf, APEHeader):
		default:
			s.r -= 3
		}
	}

	return frame, nil
}

func (s *Scanner) Close() error {
	return s.f.Close()
}
