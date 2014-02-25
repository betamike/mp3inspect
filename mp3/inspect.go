package mp3

import (
    "os"
    "io"
    "bytes"
//    "fmt"
)

func InspectFile(path string) (*MP3Info, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var mp3info MP3Info

    buf := make([]byte, 4096)
    framesSeen := 0
    var vbrBRCounter uint64

    var r, seek int64
    var n int
    var frame *AudioFrame
    var cur []byte

    for {
        if seek > 0 {
            if (r + seek) > int64(n) {
                //remove already buffered data from seek offset
                seek -= (int64(n) - r)
                _, err := file.Seek(seek, os.SEEK_CUR)
                if err != nil {
                    return nil, err
                }

                r = 0
                n, err = file.Read(buf)
                if n < 0 || err != nil {
                    if err == io.EOF {
                        break
                    }
                    return nil, err
                }
            } else {
                r += seek
            }
        }
        seek = 0

        if (r + 10) > int64(n) {
            rem := int64(n) - r
            if rem > 0 {
                copy(buf[:rem], buf[r:])
            } else {
                rem = 0
            }

            r = 0
            n, err = file.Read(buf[rem:])
            if n < 0 || err != nil {
                if err == io.EOF {
                    break
                }
                return nil, err
            }
        }

        //fmt.Printf("%d : %d == %d\n", r, 4, r+4)
        cur = buf[r:r+4]
        r += 4

        switch {
        //potentially an audio frame
        case cur[0] == 0xFF && cur[1] & 0xE0 == 0xE0:
            seek, frame = parseAudioFrame(cur)
            if frame == nil {
                //fmt.Printf("Bad potential frame\n")
                seek = 0
                r -= 3
                break
            }

            framesSeen++
            if mp3info.Bitrate != frame.Bitrate {
                if mp3info.Bitrate > 0 {
                    mp3info.IsVBR = true
                } else {
                    mp3info.Bitrate = frame.Bitrate
                }
            }
            vbrBRCounter += frame.Bitrate

            if mp3info.Samplerate != frame.Samplerate {
                mp3info.Samplerate = frame.Samplerate
            }

        //potentially ID3v1 tags
        case bytes.Equal(cur[0:3], ID3v1Header):
            seek = 127

        //potentially ID3v2 tags
        case bytes.Equal(cur[0:3], ID3v2Header):
            seek, mp3info.ID3v2 = parseID3v2Tag(buf[r-1:r+7])
            r += 6

        //potentially APE tags
        case bytes.Equal(buf, APEHeader):
        default:
            r -= 3
        }
    }
    //fmt.Printf("Seen: %d\n", framesSeen)
    return &mp3info, nil
}
