package cmhp_cdb

import (
	"fmt"
	"github.com/maldan/go-cmhp/cmhp_file"
	"sync"
	"time"
)

type Chunk[T any] struct {
	sync.RWMutex
	IsLoad    bool
	IsInit    bool
	IsChanged bool
	List      []T
	Name      string
	Index     int
}

func (c *Chunk[T]) Save() {
	c.Lock()
	defer c.Unlock()

	c.SaveWithoutLock()
}

func (c *Chunk[T]) SaveWithoutLock() {
	// Write to disk
	t := time.Now()
	err := cmhp_file.Write(c.Name+"/chunk_"+fmt.Sprintf("%v", c.Index)+".json", &c.List)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Save chunk [%v:%v] - %v records | %v\n", c.Name, c.Index, len(c.List), time.Since(t))
	c.IsChanged = false
}

func (c *Chunk[T]) SaveIfChanged() {
	if c.IsChanged {
		c.Save()
	}
}

func (c *Chunk[T]) Load() int {
	c.Lock()
	defer c.Unlock()

	t := time.Now()
	chunk, err := cmhp_file.ReadGenericJSON[[]T](c.Name + "/chunk_" + fmt.Sprintf("%v", c.Index) + ".json")
	if err != nil {
		c.List = make([]T, 0)
		c.IsInit = true
		fmt.Printf("Load chunk [%v:%v] - empty\n", c.Name, c.Index)
		return 0
	}

	c.List = chunk
	c.IsLoad = true
	c.IsInit = true
	fmt.Printf("Load chunk [%v:%v] - %v records | %v\n", c.Name, c.Index, len(chunk), time.Since(t))
	return len(chunk)
}
