package cmhp_vfs

import "time"

type FileRecord struct {
	Name    string
	Hash    string
	Offset  uint64
	Size    uint64
	Created time.Time
}

func ReadHeader() {

}
