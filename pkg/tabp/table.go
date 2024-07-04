package tabp

import (
	"strings"

	"golang.org/x/exp/slices"
)

// Table is a map based implementation of Table.
type Table struct {
	kv     map[Value]TableEntry
	keys   []Value
	seqLen int
}

// TableEntry define an entry in a Table.
type TableEntry struct {
	index int
	Key   any
	Value any
}

func (mt *Table) set(k any, v TableEntry) {
	if mt.kv == nil {
		mt.kv = make(map[Value]TableEntry)
	}

	mt.kv[k] = v
}

// Set inserts/updates value associated to given key. If provided value is nil,
// key is deleted.
func (mt *Table) Set(k, v Value) {
	entry, entryExists := mt.kv[k]

	// Delete entry.
	if v == nil && entryExists {
		// Sync seqLen.
		if k, isInt := k.(int); isInt {
			if k < mt.seqLen {
				mt.seqLen = k
			}
		}

		mt.keys = slices.Delete(mt.keys, entry.index, entry.index)
		delete(mt.kv, k)
		return
	}

	// Update entry.
	if entryExists {
		entry.Value = v
		mt.kv[k] = entry
	} else { // Insert entry.
		// Sync seqLen.
		if k, isInt := k.(int); isInt {
			if k == mt.seqLen {
				mt.seqLen++
			}
		}

		// Sync sequence length if this set value fill a hole.
		for {
			next := mt.Get(mt.seqLen)
			if next == nil {
				break
			}

			mt.seqLen++
		}

		mt.set(k, TableEntry{
			index: len(mt.keys),
			Key:   k,
			Value: v,
		})
		mt.keys = append(mt.keys, k)
	}
}

// Get returns value associated with given key.
func (mt *Table) Get(k Value) Value {
	if mt.kv == nil {
		return nil
	}

	return mt.kv[k].Value
}

// Append implements Table.
func (mt *Table) Append(v Value) int {
	// Set value.
	mt.Set(mt.seqLen, v)

	return mt.seqLen
}

// SequenceLen implements Table.
func (mt *Table) SequenceLen() int {
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

		// Key is an actual key.
		if intK, isInt := k.(int); !isInt || i != intK {
			result.WriteString(Sexpr(k))
			result.WriteString(": ")
		}

		result.WriteString(Sexpr(entry.Value))
		result.WriteRune(' ')
	}

	result.WriteRune(')')

	return result.String()
}
