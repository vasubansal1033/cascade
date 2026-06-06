package cascade_test

// Stage 3: Reads Across Levels
// Goal: implement the read fan-out path (memtable → L0 → L1 → L2) and instrument
// it with the IO counter. At this stage all data lives in L0; cross-level IO costs
// become observable in stage 6 after compaction is wired in.
// Run with: go test -run TestStage3 ./...

import (
	"fmt"
	"testing"
)

// Memtable hits must not increment the IO counter — they are in-memory reads.
func TestStage3_ReadFromMemtableDoesNotIncrementCounter(t *testing.T) {
	e := newTestEngine(t)

	_ = e.Upsert("k", []byte("v"))
	e.ResetIOCount()

	if _, err := e.Get("k"); err != nil {
		t.Fatalf("Get: %v", err)
	}
	if e.IOCount() != 0 {
		t.Fatalf("expected 0 IOs for memtable hit, got %d", e.IOCount())
	}
}

// Reading a key from a single L0 SSTable must cost exactly 1 IO.
func TestStage3_ReadFromL0IncrementsCounter(t *testing.T) {
	e := newTestEngine(t)

	_ = e.Upsert("k", []byte("v"))
	_ = e.Flush()
	e.ResetIOCount()

	if _, err := e.Get("k"); err != nil {
		t.Fatalf("Get: %v", err)
	}
	if e.IOCount() != 1 {
		t.Fatalf("expected exactly 1 IO for single L0 SSTable hit, got %d", e.IOCount())
	}
}

// Each L0 SSTable consulted during a search must add exactly 1 to the IO counter.
// A key buried in the Nth L0 file (oldest) costs exactly N IOs; the newest costs 1.
func TestStage3_MultipleL0SSTablesAccumulateIOs(t *testing.T) {
	e := newTestEngine(t)

	const n = 5
	for i := range n {
		_ = e.Upsert(fmt.Sprintf("key-%02d", i), []byte("v"))
		_ = e.Flush()
	}

	// Newest key is in the first SSTable searched — costs 1 IO.
	e.ResetIOCount()
	if _, err := e.Get(fmt.Sprintf("key-%02d", n-1)); err != nil {
		t.Fatalf("Get newest key: %v", err)
	}
	if e.IOCount() != 1 {
		t.Fatalf("newest key: expected 1 IO, got %d", e.IOCount())
	}

	// Oldest key is in the last SSTable searched — costs exactly n IOs.
	e.ResetIOCount()
	if _, err := e.Get("key-00"); err != nil {
		t.Fatalf("Get oldest key: %v", err)
	}
	if e.IOCount() != n {
		t.Fatalf("oldest key: expected %d IOs, got %d", n, e.IOCount())
	}
}

// The search must stop as soon as it encounters a tombstone — a later (older)
// version of the same key must not be returned even if it exists.
func TestStage3_SearchStopsAtTombstone(t *testing.T) {
	e := newTestEngine(t)

	// Oldest SSTable: key has a live value.
	_ = e.Upsert("k", []byte("old"))
	_ = e.Flush()

	// Newest SSTable: key is deleted.
	_ = e.Delete("k")
	_ = e.Flush()

	e.ResetIOCount()
	_, err := e.Get("k")
	if err == nil {
		t.Fatal("expected ErrNotFound after tombstone, but got a value")
	}
	// Tombstone is in the newest (first-searched) SSTable — must cost exactly 1 IO.
	if e.IOCount() != 1 {
		t.Fatalf("expected 1 IO (stopped at tombstone), got %d", e.IOCount())
	}
}
