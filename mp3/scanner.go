package mp3

import (
	"bytes"
	"io"
	"os"
)

type Scanner struct {
	f io.ReadSeeker

	// The version and layer we are looking for
	version MPEGVersion
	layer   MPEGLayer

	buf                 []byte
	FrameCount, curSize int
	vbrCounter          uint64
	curPos, absPos      int64
	eofPos              int64

	Info *MP3Info
}

func NewScanner(f io.ReadSeeker, version MPEGVersion, layer MPEGLayer) (*Scanner, error) {
	// Seek the end to get the length of the content
	eof, err := f.Seek(0, os.SEEK_END)
	if err != nil {
		return nil, err
	}
	_, err = f.Seek(0, os.SEEK_SET)
	if err != nil {
		return nil, err
	}
	return &Scanner{
		f:       f,
		version: version,
		layer:   layer,
		buf:     make([]byte, 4096),
		Info:    &MP3Info{},
		eofPos:  eof,
	}, nil
}

func NewMP3Scanner(f io.ReadSeeker) (*Scanner, error) {
	return NewScanner(f, MPEG1, LAYER3)
}

// The error returned from here does not indicate that the scan has failed,
// since it could just be eof.
func (s *Scanner) seekTo(pos int64) error {
	var err error
	if _, err = s.f.Seek(pos, os.SEEK_SET); err != nil {
		return err
	}
	s.absPos = pos

	s.curPos = 0
	s.curSize, err = s.f.Read(s.buf)
	if err != nil {
		return err
	}
	s.absPos += int64(s.curSize)
	return err
}

func (s *Scanner) frameDataAt(buf []byte) (int64, *AudioFrame, bool) {
	seekAmount, frame := parseAudioFrame(buf)
	b := !(frame == nil || frame.Version != s.version || frame.Layer != s.layer)
	if !b {
		return 0, nil, false
	}
	return seekAmount, frame, true
}

func (s *Scanner) NextFrame() (*AudioFrame, uint64, error) {
	var err error
	var ok, done bool
	var seekAmount int64
	var frame *AudioFrame
	var framePos int64
	var garbage int64

	for !done {

		// Make sure the buffer has enough unlooked-at data in it, if not move
		// the unlooked at data to the start of the buffer and read more in
		if s.curPos+4 >= int64(s.curSize) {
			rem := int64(s.curSize) - s.curPos
			if rem > 0 {
				copy(s.buf[:rem], s.buf[s.curPos:])
			} else {
				rem = 0
			}

			s.curPos = 0
			s.curSize, err = s.f.Read(s.buf[rem:])
			if err != nil {
				return nil, 0, err
			}
			//NOTE: order matters here, don't add the rem value to curSize
			// before updating absPos, otherwise we will be off by rem bytes
			s.absPos += int64(s.curSize)
			s.curSize += int(rem)
		}

		cur := s.buf[s.curPos : s.curPos+4]
		curAbsPos := s.absPos - int64(s.curSize) + s.curPos

		switch {
		//potentially an audio frame
		case cur[0] == 0xFF && cur[1]&0xE0 == 0xE0:
			if seekAmount, frame, ok = s.frameDataAt(cur); !ok {
				break
			}

			// Where to come back to if the next frame isn't valid
			returnPos := curAbsPos + 1
			framePos = curAbsPos
			// Seek to the byte right before the next frame
			nextFramePos := framePos + seekAmount - 1

			// We seek to where the next frame should be and check that it's
			// there. If it's not we seek back to the return position. seekTo
			// handles reading into buf and setting curSize/curPos and all that.
			nextFrameReal := false
			if err = s.seekTo(nextFramePos); err == nil {
				// Seeking put us right up to the end of the file (minus a possible
				// ID3v1 tag); there are no more frames, but this one is good
				left := s.eofPos - (s.absPos - int64(s.curSize)) - 1
				if left == 0 || left == int64(ID3v1Size) {
					nextFrameReal = true
				} else if _, _, ok = s.frameDataAt(s.buf[1:]); ok {
					nextFrameReal = true
				}
			}
			if !nextFrameReal {
				if err = s.seekTo(returnPos); err != nil && s.curSize == 0 {
					return nil, 0, err
				}
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

			done = true

		//potentially ID3v2 tags
		case bytes.Equal(cur[0:3], ID3v2Header):
			seekAmount, s.Info.ID3v2 = parseID3v2Tag(s.buf[s.curPos+3 : s.curPos+10])
			if err = s.seekTo(curAbsPos + seekAmount - 1); err != nil {
				return nil, 0, err
			}

		case bytes.Equal(cur[0:3], ID3v1Header):
			if (s.eofPos - curAbsPos) != int64(ID3v1Size) {
				break
			}
			s.Info.HasID3v1 = true
			if err = s.seekTo(curAbsPos + int64(ID3v1Size) - 1); err != nil {
				return nil, 0, err
			}

		case s.FrameCount == 0:
			garbage += 1
		}

		s.curPos++
	}

	//if this was the first frame found
	//set the start garbage counter
	if s.FrameCount == 1 {
		s.Info.StartGarbage = garbage
	}

	return frame, uint64(framePos), nil
}
