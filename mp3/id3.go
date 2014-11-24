package mp3

var ID3v1Header = []byte{'T', 'A', 'G'}
var ID3v2Header = []byte{'I', 'D', '3'}

var ID3v1Size = 128

type ID3v2Tag struct {
	Version  uint8
	Revision uint8
	Flags    uint16
	Size     uint32
}

func parseID3v2Tag(buf []byte) (int64, *ID3v2Tag) {
	id3 := &ID3v2Tag{}
	id3.Version = buf[0]
	id3.Revision = buf[1]
	id3.Flags = uint16(buf[2])

	//unsynchsafe the size
	id3.Size = uint32(buf[3])<<21 | uint32(buf[4])<<14 | uint32(buf[5])<<7 | uint32(buf[6]) + uint32(10)
	//footer present
	if id3.Flags&0x10 == 0x10 {
		id3.Size += 10
	}

	return int64(id3.Size), id3
}
