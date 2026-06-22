package cascade_test

// Stage 2: Flush
// Goal: persist the memtable to an SSTable on disk and read it back.
// Run with: go test -run TestStage2 ./...

import (
	"errors"
	"io/fs"
	"os"
	"testing"

	cascade "github.com/anirudhRowjee/cascade"
)

func TestStage2_HomeDirIsCreated(t *testing.T) {
	testPath := "./cascade_data"
	cascade.NewEngine(testPath)

	info, err := os.Stat(testPath)
	if err != nil {
		if !info.IsDir() {
			t.Fatal("File exists but not as directory")
		}
		if errors.Is(err, fs.ErrNotExist) {
			t.Fatal("Directory does not exist")
		}
		if errors.Is(err, fs.ErrPermission) {
			t.Fatal("Permissions not present on directory")
		}
	}

}

// Flushing a non-empty memtable must create at least one SSTable file on disk.
func TestStage2_FlushCreatesSSTableFile(t *testing.T) {
	e := newTestEngine(t)

	_ = e.Upsert("a", []byte("1"))
	_ = e.Upsert("b", []byte("2"))

	if err := e.Flush(); err != nil {
		t.Fatalf("Flush: %v", err)
	}

	// TODO: assert at least one SSTable file exists under the engine's L0 directory.
}

// Each flush must produce a distinct SSTable; n flushes → n L0 files.
func TestStage2_MultipleFlushesProduceMultipleSSTables(t *testing.T) {
	e := newTestEngine(t)

	batches := []struct{ key, val string }{
		{"a", "1"},
		{"b", "2"},
		{"c", "3"},
	}
	for _, b := range batches {
		_ = e.Upsert(b.key, []byte(b.val))
		if err := e.Flush(); err != nil {
			t.Fatalf("Flush: %v", err)
		}
	}

	// TODO: assert the L0 SSTable count equals len(batches).
}

// A tombstone written before a flush must suppress the key on subsequent reads.
func TestStage2_TombstonePersistedToSSTable(t *testing.T) {
	e := newTestEngine(t)

	_ = e.Upsert("gone", []byte("here"))
	_ = e.Delete("gone")
	if err := e.Flush(); err != nil {
		t.Fatalf("Flush: %v", err)
	}

	_, err := e.Get("gone")
	if err != cascade.ErrNotFound {
		t.Fatalf("want ErrNotFound for deleted key after flush, got %v", err)
	}
}

// When the same key appears in multiple L0 SSTables, the newest value must win.
// This validates that the read path scans L0 from most-recently-written to oldest.
func TestStage2_ReadSearchesNewestFirst(t *testing.T) {
	e := newTestEngine(t)

	_ = e.Upsert("k", []byte("old"))
	_ = e.Flush()

	_ = e.Upsert("k", []byte("new"))
	_ = e.Flush()

	got, err := e.Get("k")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if string(got) != "new" {
		t.Fatalf("want %q (newest value), got %q", "new", got)
	}
}

// A key written to the memtable and then flushed must be readable afterward.
func TestStage2_GetAfterFlush(t *testing.T) {
	e := newTestEngine(t)

	_ = e.Upsert("hello", []byte("world"))
	if err := e.Flush(); err != nil {
		t.Fatalf("Flush: %v", err)
	}

	got, err := e.Get("hello")
	if err != nil {
		t.Fatalf("Get after flush: %v", err)
	}
	if string(got) != "world" {
		t.Fatalf("want %q, got %q", "world", got)
	}
}
