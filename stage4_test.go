package cascade_test

// Stage 4: Checkpointing and Recovery
// Goal: persist engine state to disk with Sync; reconstruct it with Restart.
// Run with: go test -run TestStage4 ./...

import (
	"testing"

	cascade "github.com/anirudhRowjee/cascade"
)

// Data flushed to L0 before a Sync must be readable after Restart.
func TestStage4_RecoverRestoresFlushedData(t *testing.T) {
	e := newTestEngine(t)

	_ = e.Upsert("persistent", []byte("yes"))
	_ = e.Flush()
	if err := e.Sync(); err != nil {
		t.Fatalf("Sync: %v", err)
	}

	if err := e.Restart(); err != nil {
		t.Fatalf("Restart: %v", err)
	}

	got, err := e.Get("persistent")
	if err != nil {
		t.Fatalf("Get after restart: %v", err)
	}
	if string(got) != "yes" {
		t.Fatalf("want %q, got %q", "yes", got)
	}
}

// Data written to the active memtable but not yet flushed must survive a
// Sync/Restart cycle because Checkpoint snapshots the memtable.
func TestStage4_RecoverRestoresUnflushedMemtable(t *testing.T) {
	e := newTestEngine(t)

	_ = e.Upsert("unflushed", []byte("memtable"))
	if err := e.Sync(); err != nil {
		t.Fatalf("Sync: %v", err)
	}

	if err := e.Restart(); err != nil {
		t.Fatalf("Restart: %v", err)
	}

	got, err := e.Get("unflushed")
	if err != nil {
		t.Fatalf("Get after restart: %v", err)
	}
	if string(got) != "memtable" {
		t.Fatalf("want %q, got %q", "memtable", got)
	}
}

// A key deleted before Sync must remain deleted after Restart.
func TestStage4_RecoverPreservesTombstone(t *testing.T) {
	e := newTestEngine(t)

	_ = e.Upsert("gone", []byte("here"))
	_ = e.Flush()
	_ = e.Delete("gone")
	if err := e.Sync(); err != nil {
		t.Fatalf("Sync: %v", err)
	}

	if err := e.Restart(); err != nil {
		t.Fatalf("Restart: %v", err)
	}

	_, err := e.Get("gone")
	if err != cascade.ErrNotFound {
		t.Fatalf("want ErrNotFound after restart, got %v", err)
	}
}

// Restart on an engine that has never called Sync must succeed without error
// and leave the engine in a clean empty state.
func TestStage4_RestartWithNoCheckpointIsNoop(t *testing.T) {
	e := newTestEngine(t)

	if err := e.Restart(); err != nil {
		t.Fatalf("Restart with no checkpoint: %v", err)
	}
	_, err := e.Get("anything")
	if err != cascade.ErrNotFound {
		t.Fatalf("want ErrNotFound in empty restarted engine, got %v", err)
	}
}
