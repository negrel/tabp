package tabp

import (
	"strings"

	"golang.org/x/exp/slices"
)

// Table is a map based implementation of Table.
type Table struct {
	kv     map[Value]TableEntry
	keys   []Value
	array  []TableEntry
	seqLen int
}

// TableEntry define an entry in a Table.
type TableEntry struct {
	keyIndex int
	Key      any
	Value    any
}

// Set inserts/updates value associated to given key. If provided value is nil,
// entry is deleted.
func (mt *Table) Set(k, v Value) {
	arrayK, kIsInt := k.(int)
	if kIsInt && arrayK <= mt.seqLen {
		mt.arraySet(arrayK, v)
		return
	}

	mt.mapSet(k, v)
}

// 0 <= k <= mt.seqLen
func (mt *Table) arraySet(k int, v Value) {
	// Delete.
	if v == nil {
		// migrate element after k to map.
		for i := k + 1; i < mt.seqLen; i++ {
			mt.mapSetEntry(i, mt.array[i])
		}

		// Delete k.
		mt.seqLen = k
		index := mt.array[k].keyIndex
		mt.array = mt.array[:k]
		mt.keys = slices.Delete(mt.keys, index, index+1)
		return
	}

	entry := TableEntry{
		keyIndex: len(mt.keys),
		Key:      k,
		Value:    v,
	}

	// Grow if needed. (Insert)
	if len(mt.array)-1 <= k {
		mt.array = append(mt.array, entry)
		mt.keys = append(mt.keys, k) // Add key on insert.
		mt.seqLen++
		mt.fillSeqHole()
		return
	}

	// Update value in array.
	mt.array[k] = entry
}

func (mt *Table) arrayGet(k int) Value {
	// Grow if needed.
	if mt.seqLen < k {
		return nil
	}

	return mt.array[k].Value
}

func (mt *Table) mapSet(k Value, v Value) {
	// Delete.
	if v == nil {
		entry, entryExists := mt.kv[k]
		if entryExists {
			mt.keys = slices.Delete(mt.keys, entry.keyIndex, entry.keyIndex)
			delete(mt.kv, k)
		}
		return
	}

	entry, entryExists := mt.kv[k]
	entry.Value = v

	// Insert.
	if !entryExists {
		entry.keyIndex = len(mt.keys)
		entry.Key = k
		mt.keys = append(mt.keys, k)
	}

	// Update.
	mt.mapSetEntry(k, entry)
}

func (mt *Table) mapSetEntry(k Value, entry TableEntry) {
	if mt.kv == nil {
		mt.kv = make(map[Value]TableEntry)
	}

	mt.kv[k] = entry
}

func (mt *Table) mapGet(k Value) Value {
	return mt.mapGetEntry(k).Value
}

func (mt *Table) mapGetEntry(k Value) TableEntry {
	if mt.kv == nil {
		return TableEntry{}
	}

	return mt.kv[k]
}

// Get returns value associated with given key. A nil value is returned if key
// is not found.
func (mt *Table) Get(k Value) Value {
	arrayK, kIsInt := k.(int)
	if kIsInt && arrayK < mt.seqLen {
		return mt.arrayGet(arrayK)
	}

	return mt.mapGet(k)
}

// Append adds given value at the end of table's sequence. A sequence starts
// with key tab[0] and ends when tab[k] is nil.
func (mt *Table) Append(v Value) int {
	// Set value.
	mt.Set(mt.seqLen, v)

	return mt.seqLen
}

// SeqLen returns length of table sequence.
func (mt *Table) SeqLen() int {
	return mt.seqLen
}

// Keys implements Table.
func (mt *Table) Keys() []Value {
	return mt.keys
}

// Values implements Table.
func (mt *Table) Values() []Value {
	var result []Value

	for _, k := range mt.keys {
		result = append(result, mt.kv[k].Value)
	}

	return result
}

// Entries implements Table.
func (mt *Table) Entries() []TableEntry {
	var entries []TableEntry
	for _, k := range mt.keys {
		entries = append(entries, mt.kv[k])
	}

	return entries
}

// Map implements Table.
func (mt *Table) Map(fn func(k, v Value) Value) {
	for k, v := range mt.kv {
		v.Value = fn(k, v.Value)
		mt.kv[k] = v
	}
}

// ForEach implements Table.
func (mt *Table) ForEach(fn func(k, v Value) bool) {
	for k, v := range mt.kv {
		stop := fn(k, v.Value)
		if stop {
			break
		}
	}
}

// ToSExpr implements SExpr.
func (mt *Table) ToSExpr() string {
	var result strings.Builder
	result.WriteRune('(')

	// Iterate over keys in insertion order.
	for i, k := range mt.keys {
		entry := mt.kv[k]

		intK, kIsInt := k.(int)

		// Key is an actual key.
		if !kIsInt || i != intK {
			result.WriteString(Sexpr(k))
			result.WriteString(": ")
		}

		// entry is stored in array.
		if kIsInt && intK < mt.seqLen {
			entry = mt.array[intK]
		}

		result.WriteString(Sexpr(entry.Value))
		if i < len(mt.keys)-1 {
			result.WriteRune(' ')
		}
	}

	result.WriteRune(')')

	return result.String()
}

func (mt *Table) fillSeqHole() {
	for {
		entry := mt.mapGetEntry(mt.seqLen)
		if entry.Value == nil {
			break
		}

		mt.kv[mt.seqLen] = TableEntry{}
		mt.array = append(mt.array, entry)
		mt.seqLen++
	}
}
