package cascade

const MemtableMaxBytes = 1024 // 1 KB

type Memtable struct {
	data map[string]Entry
	size int
}

func NewMemtable() *Memtable { return nil }

func (m *Memtable) Get(key string) (Entry, bool)    { return Entry{}, false }
func (m *Memtable) Set(key string, entry Entry)     {}
func (m *Memtable) IsFull() bool                    { return false }
func (m *Memtable) Size() int                       { return 0 }
func (m *Memtable) SortedEntries() []KVEntry        { return nil }
