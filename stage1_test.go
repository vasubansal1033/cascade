package cascade_test

// Stage 1: Setup
// Goal: basic CRUD against the memtable only — no disk involved.
// Run with: go test -run TestStage1 ./...

import (
	"testing"

	cascade "github.com/anirudhRowjee/cascade"
)

func newTestEngine(t *testing.T) *cascade.Engine {
	t.Helper()
	dir := t.TempDir()
	e, err := cascade.NewEngine(dir)
	if err != nil {
		t.Fatalf("NewEngine: %v", err)
	}
	return e
}

func TestState1_Sanity(t *testing.T) {
	_, err := cascade.NewEngine(t.TempDir())
	if err != nil {
		t.Error("Sanity test is failing!!")
	}
}

func TestStage1_UpsertAndGet(t *testing.T) {
	e := newTestEngine(t)
	if err := e.Upsert("hello", []byte("world")); err != nil {
		t.Fatalf("Upsert: %v", err)
	}
	got, err := e.Get("hello")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if string(got) != "world" {
		t.Fatalf("want %q, got %q", "world", got)
	}
}

func TestStage1_Overwrite(t *testing.T) {
	e := newTestEngine(t)
	_ = e.Upsert("k", []byte("v1"))
	_ = e.Upsert("k", []byte("v2"))
	got, err := e.Get("k")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if string(got) != "v2" {
		t.Fatalf("want %q, got %q", "v2", got)
	}
}

func TestStage1_DeleteProducesTombstone(t *testing.T) {
	e := newTestEngine(t)
	_ = e.Upsert("k", []byte("v"))
	if err := e.Delete("k"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := e.Get("k")
	if err != cascade.ErrNotFound {
		t.Fatalf("want ErrNotFound after delete, got %v", err)
	}
}

func TestStage1_GetMissingKey(t *testing.T) {
	e := newTestEngine(t)
	_, err := e.Get("nonexistent")
	if err != cascade.ErrNotFound {
		t.Fatalf("want ErrNotFound, got %v", err)
	}
}
