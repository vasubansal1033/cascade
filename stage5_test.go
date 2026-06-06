package cascade_test

// Stage 5: Compaction in Isolation
// Goal: understand the compaction algorithm as a standalone operation —
// merge two SSTable streams, handle tombstones, and produce a new sorted SSTable.
// Run with: go test -run TestStage5 ./...

import (
	"testing"

	cascade "github.com/anirudhRowjee/cascade"
)

// Merging two SSTables with non-overlapping keys must produce a single
// sorted SSTable containing all keys from both inputs.
func TestStage5_MergeNonOverlappingKeys(t *testing.T) {
	dir := t.TempDir()

	a, err := cascade.WriteSSTable(dir+"/a.sst", []cascade.KVEntry{
		{Key: "apple", Entry: cascade.Entry{Value: []byte("1")}},
		{Key: "cherry", Entry: cascade.Entry{Value: []byte("3")}},
	})
	if err != nil {
		t.Fatalf("WriteSSTable a: %v", err)
	}
	b, err := cascade.WriteSSTable(dir+"/b.sst", []cascade.KVEntry{
		{Key: "banana", Entry: cascade.Entry{Value: []byte("2")}},
		{Key: "date", Entry: cascade.Entry{Value: []byte("4")}},
	})
	if err != nil {
		t.Fatalf("WriteSSTable b: %v", err)
	}

	counter := cascade.NewIOCounter()
	merged, err := cascade.MergeSSTables([]*cascade.SSTable{a}, []*cascade.SSTable{b}, dir+"/merged.sst", counter)
	if err != nil {
		t.Fatalf("MergeSSTables: %v", err)
	}

	entries, err := merged.Scan(counter)
	if err != nil {
		t.Fatalf("Scan: %v", err)
	}
	want := []string{"apple", "banana", "cherry", "date"}
	if len(entries) != len(want) {
		t.Fatalf("want %d entries, got %d", len(want), len(entries))
	}
	for i, e := range entries {
		if e.Key != want[i] {
			t.Errorf("entry[%d]: want key %q, got %q", i, want[i], e.Key)
		}
	}
}

// When both inputs contain the same key, the value from the newer SSTable wins.
func TestStage5_MergeOverlappingKeysNewerWins(t *testing.T) {
	dir := t.TempDir()

	older, _ := cascade.WriteSSTable(dir+"/older.sst", []cascade.KVEntry{
		{Key: "k", Entry: cascade.Entry{Value: []byte("old")}},
	})
	newer, _ := cascade.WriteSSTable(dir+"/newer.sst", []cascade.KVEntry{
		{Key: "k", Entry: cascade.Entry{Value: []byte("new")}},
	})

	counter := cascade.NewIOCounter()
	// newer is passed as the first (higher-priority) argument.
	merged, err := cascade.MergeSSTables([]*cascade.SSTable{newer}, []*cascade.SSTable{older}, dir+"/merged.sst", counter)
	if err != nil {
		t.Fatalf("MergeSSTables: %v", err)
	}

	entry, found, err := merged.Get("k", counter)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if !found {
		t.Fatal("key not found in merged SSTable")
	}
	if string(entry.Value) != "new" {
		t.Fatalf("want %q, got %q", "new", entry.Value)
	}
}

// Tombstones from the newer SSTable must suppress the older live value and
// must not appear in the output — dead keys are fully dropped.
func TestStage5_MergeDropsTombstonesWhenOlderValueExists(t *testing.T) {
	dir := t.TempDir()

	older, _ := cascade.WriteSSTable(dir+"/older.sst", []cascade.KVEntry{
		{Key: "k", Entry: cascade.Entry{Value: []byte("alive")}},
	})
	newer, _ := cascade.WriteSSTable(dir+"/newer.sst", []cascade.KVEntry{
		{Key: "k", Entry: cascade.Entry{Tombstone: true}},
	})

	counter := cascade.NewIOCounter()
	merged, err := cascade.MergeSSTables([]*cascade.SSTable{newer}, []*cascade.SSTable{older}, dir+"/merged.sst", counter)
	if err != nil {
		t.Fatalf("MergeSSTables: %v", err)
	}

	entries, err := merged.Scan(counter)
	if err != nil {
		t.Fatalf("Scan: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected empty merged SSTable after tombstone drop, got %d entries", len(entries))
	}
}

// The output of MergeSSTables must be sorted in ascending key order.
func TestStage5_MergeOutputIsSorted(t *testing.T) {
	dir := t.TempDir()

	a, _ := cascade.WriteSSTable(dir+"/a.sst", []cascade.KVEntry{
		{Key: "c", Entry: cascade.Entry{Value: []byte("3")}},
		{Key: "a", Entry: cascade.Entry{Value: []byte("1")}},
	})
	b, _ := cascade.WriteSSTable(dir+"/b.sst", []cascade.KVEntry{
		{Key: "d", Entry: cascade.Entry{Value: []byte("4")}},
		{Key: "b", Entry: cascade.Entry{Value: []byte("2")}},
	})

	counter := cascade.NewIOCounter()
	merged, err := cascade.MergeSSTables([]*cascade.SSTable{a}, []*cascade.SSTable{b}, dir+"/merged.sst", counter)
	if err != nil {
		t.Fatalf("MergeSSTables: %v", err)
	}

	entries, err := merged.Scan(counter)
	if err != nil {
		t.Fatalf("Scan: %v", err)
	}
	for i := 1; i < len(entries); i++ {
		if entries[i].Key <= entries[i-1].Key {
			t.Errorf("output not sorted at index %d: %q >= %q", i, entries[i-1].Key, entries[i].Key)
		}
	}
}
