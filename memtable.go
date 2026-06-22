package cascade

import "github.com/egregors/sortedmap"

const MemtableMaxBytes = 1024 // 1 KB

type Memtable struct {
	data sortedmap.SortedMap[map[string]KVEntry, string, KVEntry]
	size int
	/*
		What data do we need to solve the following problems?
		1. How do we keep track of how full the memtable is?
	*/
}

func NewMemtable() *Memtable { return nil }

func (m *Memtable) Get(key string) (KVEntry, bool) { return KVEntry{}, false }
func (m *Memtable) Set(key string, entry KVEntry)  {}
func (m *Memtable) IsFull() bool                   { return false }
func (m *Memtable) Size() int                      { return 0 }
func (m *Memtable) SortedEntries() []KVEntry       { return nil }
