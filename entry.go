package cascade

import "fmt"

type KVEntry struct {
	Key         string
	Value       string
	IsTombstone bool
}

func GenerateDelete(key string) KVEntry {
	return KVEntry{
		Key:         key,
		Value:       "",
		IsTombstone: true,
	}
}

func GenerateUpsert(key, value string) KVEntry {
	return KVEntry{
		Key:         key,
		Value:       value,
		IsTombstone: false,
	}
}

// Generates a delete for key `key-000...<keyNum>`
func GenerateNumberedDelete(keyNum int64) KVEntry {
	key := fmt.Sprintf("key-%010d", keyNum)
	return KVEntry{
		Key:         key,
		Value:       "",
		IsTombstone: true,
	}
}

// Generates an upsert for key `key-000...<keyNum>`
func GenerateNumberedUpsert(keyNum int64, value string) KVEntry {
	key := fmt.Sprintf("key-%010d", keyNum)
	return KVEntry{
		Key:         key,
		Value:       value,
		IsTombstone: false,
	}
}
