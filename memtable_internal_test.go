package cascade

import (
	"fmt"
	"testing"

	sm "github.com/egregors/sortedmap"
)

func TestMemtableInternal(t *testing.T) {
	m := sm.NewFromMap(map[string]int{
		"Bob":   31,
		"Alice": 26,
		"Eve":   84,
	}, func(i, j sm.KV[string, int]) bool {
		return i.Key < j.Key
	})
	fmt.Println(m.Collect())
}

func TestSortedMapInternal(t *testing.T) {

	m := sm.New[map[string]KVEntry](
		func(i, j sm.KV[string, KVEntry]) bool {
			return i.Key < j.Key
		})

	m.Insert("a_hello", KVEntry{Key: "hello", Value: "world", IsTombstone: false})
	fmt.Println(m.Collect())

	m.Insert("c_hello", GenerateUpsert("c_hello", "lol"))
	m.Insert("b_hello", GenerateUpsert("b_hello", "world1"))
	m.Insert("c_hello", GenerateDelete("c_hello"))

	fmt.Println(m.CollectAll())

}
