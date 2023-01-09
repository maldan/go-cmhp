package cmhp_cdb

import (
	"fmt"
	"github.com/maldan/go-cmhp/cmhp_file"
	"github.com/maldan/go-cmhp/cmhp_slice"
	"sync"
)

type ChunkMaster[T any] struct {
	sync.Mutex
	Name          string
	Size          int
	AutoIncrement int
	ChunkList     []Chunk[T]
}

type ChunkMasterAutoIncrement struct {
	Counter int `json:"counter"`
}

func (m *ChunkMaster[T]) LoadAll() {
	m.Lock()
	defer m.Unlock()

	if m.Name == "" {
		panic("Chunk name not specified")
	}

	// Read chunk info
	info, _ := cmhp_file.ReadGenericJSON[ChunkMasterAutoIncrement](m.Name + "/counter.json")
	m.AutoIncrement = info.Counter
	fmt.Printf("Prepare chunk [%v] - %v\n", m.Name, info)

	// Init chunks
	m.ChunkList = make([]Chunk[T], m.Size)
	loadTotal := 0
	for i := 0; i < m.Size; i++ {
		m.ChunkList[i].Name = m.Name
		m.ChunkList[i].Index = i
		loadTotal += m.ChunkList[i].Load()
	}

	fmt.Printf("Load chunk [%v] - %v total\n", m.Name, loadTotal)
}

func (m *ChunkMaster[T]) SaveAll() {
	for i := 0; i < m.Size; i++ {
		m.ChunkList[i].Save()
	}
}

func (m *ChunkMaster[T]) SaveChanged() {
	for i := 0; i < m.Size; i++ {
		m.ChunkList[i].SaveIfChanged()
	}
}

func (m *ChunkMaster[T]) SaveChunk(toHash any) {
	hash := m.StrHash(fmt.Sprintf("%v", toHash), m.Size)
	m.ChunkList[hash].Save()
}

// ForEach go over each value in all chunks
func (m *ChunkMaster[T]) ForEach(fn func(item T) bool) {
	status := true
	for i := 0; i < m.Size; i++ {
		m.ChunkList[i].RLock()
		for j := 0; j < len(m.ChunkList[i].List); j++ {
			if !fn(m.ChunkList[i].List[j]) {
				status = false
				break
			}
		}
		m.ChunkList[i].RUnlock()
		if !status {
			break
		}
	}
}

// Find value in chunk [toHash] by [cond]
func (m *ChunkMaster[T]) Find(toHash any, cond func(v T) bool) (T, bool) {
	hash := m.StrHash(fmt.Sprintf("%v", toHash), m.Size)
	m.ChunkList[hash].RLock()
	defer m.ChunkList[hash].RUnlock()

	for _, item := range m.ChunkList[hash].List {
		if cond(item) {
			return item, true
		}
	}
	return *new(T), false
}

// Contains value in chunk [toHash] by [cond]
func (m *ChunkMaster[T]) Contains(toHash any, cond func(v T) bool) bool {
	hash := m.StrHash(fmt.Sprintf("%v", toHash), m.Size)
	m.ChunkList[hash].RLock()
	defer m.ChunkList[hash].RUnlock()

	for _, item := range m.ChunkList[hash].List {
		if cond(item) {
			return true
		}
	}
	return false
}

// AddOrReplace value to chunk [toHash] and save it
func (m *ChunkMaster[T]) AddOrReplace(v T, toHash any, where func(v T) bool) {
	if !m.Replace(v, toHash, where) {
		m.Add(v, toHash)
	}
}

// Add value to chunk [toHash] and save it
func (m *ChunkMaster[T]) Add(v T, toHash any) {
	hash := m.StrHash(fmt.Sprintf("%v", toHash), m.Size)
	m.ChunkList[hash].Lock()
	m.ChunkList[hash].List = append(m.ChunkList[hash].List, v)
	m.ChunkList[hash].IsChanged = true
	m.ChunkList[hash].SaveWithoutLock()
	m.ChunkList[hash].Unlock()
}

// Replace value in chunk [toHash] by condition [where]
func (m *ChunkMaster[T]) Replace(val T, toHash any, where func(v T) bool) bool {
	hash := m.StrHash(fmt.Sprintf("%v", toHash), m.Size)
	m.ChunkList[hash].Lock()
	defer m.ChunkList[hash].Unlock()
	for i := 0; i < len(m.ChunkList[hash].List); i++ {
		if where(m.ChunkList[hash].List[i]) {
			m.ChunkList[hash].List[i] = val
			m.ChunkList[hash].IsChanged = true
			m.ChunkList[hash].SaveWithoutLock()
			return true
		}
	}
	return false
}

// Delete values in chunk [toHash] by condition [where]
func (m *ChunkMaster[T]) Delete(toHash any, where func(v T) bool) {
	hash := m.StrHash(fmt.Sprintf("%v", toHash), m.Size)
	m.ChunkList[hash].Lock()
	defer m.ChunkList[hash].Unlock()

	// Filter values
	lenWas := len(m.ChunkList[hash].List)
	m.ChunkList[hash].List = cmhp_slice.Filter(m.ChunkList[hash].List, func(i T) bool {
		return !where(i)
	})

	// Elements was deletes
	if lenWas != len(m.ChunkList[hash].List) {
		m.ChunkList[hash].IsChanged = true
		m.ChunkList[hash].SaveWithoutLock()
	}
}

// DeleteInAll values in all chunks by condition [where]
func (m *ChunkMaster[T]) DeleteInAll(where func(v T) bool) {
	for i := 0; i < m.Size; i++ {
		m.ChunkList[i].Lock()

		// Filter values
		lenWas := len(m.ChunkList[i].List)
		m.ChunkList[i].List = cmhp_slice.Filter(m.ChunkList[i].List, func(i T) bool {
			return !where(i)
		})

		// Elements was deletes
		if lenWas != len(m.ChunkList[i].List) {
			m.ChunkList[i].IsChanged = true
		}

		m.ChunkList[i].Unlock()
	}
}

func (m *ChunkMaster[T]) FindInAll(fn func(v T) bool) (T, bool) {
	out := *new(T)
	isFound := false
	m.ForEach(func(item T) bool {
		if fn(item) {
			out = item
			isFound = true
			return false
		}
		return true
	})
	return out, isFound
}

func (m *ChunkMaster[T]) ContainsInAll(fn func(v T) bool) bool {
	isFound := false
	m.ForEach(func(item T) bool {
		if fn(item) {
			isFound = true
			return false
		}
		return true
	})
	return isFound
}

func (m *ChunkMaster[T]) FilterInAll(fn func(v T) bool) []T {
	out := make([]T, 0)

	m.ForEach(func(item T) bool {
		if fn(item) {
			out = append(out, item)
		}
		return true
	})
	return out
}

// All Copy all values from list
func (m *ChunkMaster[T]) All() []T {
	out := make([]T, 0)

	m.ForEach(func(item T) bool {
		out = append(out, item)
		return true
	})
	return out
}

// GetId generate thread safe autoincrement id
func (m *ChunkMaster[T]) GetId() int {
	out := 0
	m.Lock()
	m.AutoIncrement += 1
	out = m.AutoIncrement
	info := ChunkMasterAutoIncrement{m.AutoIncrement}
	cmhp_file.Write(m.Name+"/counter.json", &info)
	m.Unlock()

	return out
}

func (m *ChunkMaster[T]) TotalElements() int {
	count := 0
	for i := 0; i < m.Size; i++ {
		m.ChunkList[i].RLock()
		count += len(m.ChunkList[i].List)
		m.ChunkList[i].RUnlock()
	}
	return count
}

func (m *ChunkMaster[T]) StrHash(str string, max int) int {
	hash := 0
	for i := 0; i < len(str); i++ {
		hash += int(str[i])
	}
	return hash % max
}
