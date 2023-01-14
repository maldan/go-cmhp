package cmhp_cdb

import (
	"fmt"
	"github.com/maldan/go-cmhp/cmhp_file"
	"github.com/maldan/go-cmhp/cmhp_slice"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"time"
)

type Chunk[T idAble] struct {
	sync.RWMutex
	IsLoad       bool
	IsInit       bool
	IsChanged    bool
	List         []T
	Name         string
	Id           int
	IndexList    []string
	indexStorage map[string][]*T
	ShowLogs     bool
}

func (c *Chunk[T]) BuildIndexMap() {
	c.indexStorage = make(map[string][]*T)

	// Build index
	for i := 0; i < len(c.List); i++ {
		c.AddToIndex(&c.List[i])
	}
}

func (c *Chunk[T]) AddToIndex(ref *T) {
	c.Lock()
	defer c.Unlock()

	for _, index := range c.IndexList {
		f := reflect.ValueOf(ref).Elem().FieldByName(index)
		mapIndex := reflect.ValueOf(f).Interface()
		strIndex := fmt.Sprintf("%v:%v", index, mapIndex)
		c.indexStorage[strIndex] = append(c.indexStorage[strIndex], ref)
	}
}

func (c *Chunk[T]) Save() {
	c.Lock()
	defer c.Unlock()
	c.SaveWithoutLock()
}

func (c *Chunk[T]) SaveWithoutLock() {
	// Write to disk
	t := time.Now()
	err := cmhp_file.Write(c.Name+"/chunk_"+fmt.Sprintf("%v", c.Id)+".json.tmp", &c.List)
	if err != nil {
		panic(err)
	}

	// Delete old
	cmhp_file.DeleteFile(c.Name + "/chunk_" + fmt.Sprintf("%v", c.Id) + ".json")

	// Replace
	err = cmhp_file.Rename(c.Name+"/chunk_"+fmt.Sprintf("%v", c.Id)+".json.tmp", c.Name+"/chunk_"+fmt.Sprintf("%v", c.Id)+".json")
	if err != nil {
		panic(err)
	}

	if c.ShowLogs {
		name := c.Name
		wd, _ := os.Getwd()
		name, _ = filepath.Rel(wd, c.Name)
		fmt.Printf("Save chunk [%v:%v] - %v records | %v\n", name, c.Id, len(c.List), time.Since(t))
	}
	c.IsChanged = false
}

func (c *Chunk[T]) SaveIfChanged() {
	if c.IsChanged {
		c.Save()
	}
}

func (c *Chunk[T]) Load() int {
	c.Lock()
	defer c.BuildIndexMap()
	defer c.Unlock()

	t := time.Now()
	chunk, err := cmhp_file.ReadGenericJSON[[]T](c.Name + "/chunk_" + fmt.Sprintf("%v", c.Id) + ".json")
	if err != nil {
		c.List = make([]T, 0)
		c.IsInit = true
		if c.ShowLogs {
			fmt.Printf("Load chunk [%v:%v] - empty\n", c.Name, c.Id)
		}
		return 0
	}

	c.List = chunk
	c.IsLoad = true
	c.IsInit = true
	if c.ShowLogs {
		fmt.Printf("Load chunk [%v:%v] - %v records | %v\n", c.Name, c.Id, len(chunk), time.Since(t))
	}
	return len(chunk)
}

// Find value in chunk by [cond]
func (c *Chunk[T]) Find(cond func(v T) bool) (T, bool) {
	c.RLock()
	defer c.RUnlock()

	for _, item := range c.List {
		if cond(item) {
			return item, true
		}
	}
	return *new(T), false
}

// Delete values in chunk by condition [where]
func (c *Chunk[T]) Delete(where func(v T) bool) {
	c.Lock()
	defer c.Unlock()

	// Filter values
	lenWas := len(c.List)
	c.List = cmhp_slice.Filter(c.List, func(i T) bool {
		return !where(i)
	})

	// Elements was deletes
	if lenWas != len(c.List) {
		c.IsChanged = true
		// c.SaveWithoutLock()
	}
}

func (c *Chunk[T]) FindByIndex(indexName string, indexValue any) (T, bool) {
	c.RLock()
	defer c.RUnlock()

	strIndex := fmt.Sprintf("%v:%v", indexName, indexValue)
	for _, val := range c.indexStorage[strIndex] {
		return *val, true
	}
	return *new(T), false
}

func (c *Chunk[T]) FindAllByIndex(indexName string, indexValue any) []T {
	c.RLock()
	defer c.RUnlock()

	strIndex := fmt.Sprintf("%v:%v", indexName, indexValue)
	out := make([]T, 0)
	for _, val := range c.indexStorage[strIndex] {
		out = append(out, *val)
	}
	return out
}

// Replace value in chunk [toHash] by condition [where]
/*func (c *Chunk[T]) Replace(val T, where func(v T) bool) bool {
	c.Lock()
	defer c.Unlock()

	for i := 0; i < len(c.List); i++ {
		if where(c.List[i]) {
			c.List[i] = val
			c.IsChanged = true
			c.SaveWithoutLock()
			return true
		}
	}
	return false
}*/

// Add value to chunk [toHash] and save it
/*func (c *Chunk[T]) Add(v T) {
	//hash := m.Hash(v.GetId(), m.Size)
	//c.FastAdd(v)
	//c.Save()
	c.Lock()
	defer c.Unlock()

	c.List = append(c.List, v)
	c.IsChanged = true
}*/

// FastAdd value to chunk [toHash] without save
/*func (c *Chunk[T]) FastAdd(v T) {
	c.Lock()
	defer c.Unlock()

	c.List = append(c.List, v)
	c.IsChanged = true

	// Build index
	for _, index := range m.IndexList {
		m.AddIndex(index, &v)
	}
}*/

// Contains value in chunk by [cond]
func (c *Chunk[T]) Contains(cond func(v T) bool) bool {
	c.RLock()
	defer c.RUnlock()

	for _, item := range c.List {
		if cond(item) {
			return true
		}
	}
	return false
}

// AddOrReplace value to chunk and save it
/*func (c *Chunk[T]) AddOrReplace(v T, where func(v T) bool) {
	if !c.Replace(v, where) {
		c.Add(v)
	}
}
*/
