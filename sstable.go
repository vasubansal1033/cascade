package cascade

import "errors"

var ErrKeyNotFound = errors.New("key not found")

type SSTable struct {
	Path string
}

func WriteSSTable(path string, entries []KVEntry) (*SSTable, error) { return nil, nil }

func (s *SSTable) Get(key string, counter *IOCounter) (KVEntry, bool, error) {
	return KVEntry{}, false, nil
}

func (s *SSTable) Scan(counter *IOCounter) ([]KVEntry, error) { return nil, nil }
