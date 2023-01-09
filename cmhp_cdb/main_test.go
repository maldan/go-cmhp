package cmhp_cdb_test

import (
	"github.com/maldan/go-cmhp/cmhp_cdb"
	"testing"
)

type Test struct {
	Id int    `json:"id"`
	A  string `json:"a"`
}

func TestM(t *testing.T) {
	testChunk := cmhp_cdb.ChunkMaster[Test]{Size: 10, Name: "../a/gas"}
	testChunk.LoadAll()

	// Add
	testChunk.Add(Test{Id: 1}, 1)
	testChunk.Add(Test{Id: 2}, 2)
	testChunk.Add(Test{Id: 3}, 3)
	if testChunk.TotalElements() != 3 {
		t.Errorf("Fuck you")
	}

	// Delete
	testChunk.Delete(1, func(t Test) bool {
		return t.Id == 1
	})
	if testChunk.TotalElements() != 2 {
		t.Errorf("Fuck you")
	}

	// Test if changed
	if !testChunk.ChunkList[9].IsChanged {
		t.Errorf("Chunk must be changed")
	}
	if testChunk.ChunkList[3].IsChanged {
		t.Errorf("Chunk must not be changed")
	}

	// All
	if len(testChunk.All()) != 2 {
		t.Errorf("Fuck you")
	}
}

func TestAdd(t *testing.T) {
	testChunk := cmhp_cdb.ChunkMaster[Test]{Size: 10, Name: "../a/gas"}
	testChunk.LoadAll()

	// Add
	for i := 0; i < 10000; i++ {
		testChunk.Add(Test{Id: i}, i)
	}
	if testChunk.TotalElements() != 10000 {
		t.Errorf("Element amount not match")
	}

	// Test if changed
	if !testChunk.ChunkList[9].IsChanged {
		t.Errorf("Chunk must be changed")
	}
}

func TestDelete(t *testing.T) {
	testChunk := cmhp_cdb.ChunkMaster[Test]{Size: 10, Name: "../a/gas"}
	testChunk.LoadAll()

	// Add
	for i := 0; i < 10000; i++ {
		testChunk.Add(Test{Id: i}, i)
	}

	// Delete
	testChunk.Delete(1, func(t Test) bool {
		return t.Id == 1
	})
	if testChunk.TotalElements() != 10000-1 {
		t.Errorf("Delete not working")
	}

	// Delete in all
	testChunk.Delete(555, func(t Test) bool {
		return t.Id == 555
	})
	if testChunk.TotalElements() != 10000-2 {
		t.Errorf("Delete not working")
	}

	// Find
	_, ok := testChunk.Find(1, func(t Test) bool {
		return t.Id == 1
	})
	if ok {
		t.Errorf("Delete not working")
	}
	_, ok = testChunk.Find(555, func(t Test) bool {
		return t.Id == 555
	})
	if ok {
		t.Errorf("Delete not working")
	}
	_, ok = testChunk.FindInAll(func(t Test) bool {
		return t.Id == 1
	})
	if ok {
		t.Errorf("Delete not working")
	}
	_, ok = testChunk.FindInAll(func(t Test) bool {
		return t.Id == 555
	})
	if ok {
		t.Errorf("Delete not working")
	}
}

func TestUpdate(t *testing.T) {
	testChunk := cmhp_cdb.ChunkMaster[Test]{Size: 10, Name: "../a/gas"}
	testChunk.LoadAll()

	// Add
	testChunk.Add(Test{Id: 1}, 1)
	testChunk.Add(Test{Id: 2}, 2)
	testChunk.Add(Test{Id: 3}, 3)
	if testChunk.TotalElements() != 3 {
		t.Errorf("Fuck you")
	}

	// Update
	testChunk.Replace(Test{A: "gas"}, 1, func(t Test) bool {
		return t.Id == 1
	})
	if testChunk.TotalElements() != 3 {
		t.Errorf("Fuck you")
	}

	// Test if changed
	if testChunk.ChunkList[9].List[0].A != "gas" {
		t.Errorf("Update not working")
	}
}

func TestFind(t *testing.T) {
	testChunk := cmhp_cdb.ChunkMaster[Test]{Size: 10, Name: "../a/gas"}
	testChunk.LoadAll()

	// Add
	for i := 0; i < 1000; i++ {
		testChunk.Add(Test{Id: i}, i)
	}
	if testChunk.TotalElements() != 1000 {
		t.Errorf("Fuck you")
	}

	// Find in chunk
	v, ok := testChunk.Find(432, func(t Test) bool {
		return t.Id == 432
	})
	if !ok {
		t.Errorf("Value not found")
	}

	// Test if changed
	if v.Id != 432 {
		t.Errorf("Find not working")
	}

	// Find in all
	v, ok = testChunk.FindInAll(func(t Test) bool {
		return t.Id == 768
	})
	if !ok {
		t.Errorf("Value not found")
	}

	// Test if changed
	if v.Id != 768 {
		t.Errorf("Find not working")
	}
}

func TestFilter(t *testing.T) {
	testChunk := cmhp_cdb.ChunkMaster[Test]{Size: 10, Name: "../a/gas"}
	testChunk.LoadAll()

	// Add
	for i := 0; i < 1000; i++ {
		testChunk.Add(Test{Id: i}, i)
	}
	if testChunk.TotalElements() != 1000 {
		t.Errorf("Fuck you")
	}

	// Find in chunk
	list := testChunk.FilterInAll(func(t Test) bool {
		return t.Id >= 500
	})
	if len(list) != 500 {
		t.Errorf("Fuck you")
	}
}
