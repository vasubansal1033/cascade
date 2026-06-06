package cascade

type Entry struct {
	Value     []byte
	Tombstone bool
}

type KVEntry struct {
	Key   string
	Entry Entry
}
