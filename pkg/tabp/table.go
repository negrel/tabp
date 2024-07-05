package tabp

import (
	"strings"
)

// Table is a map based implementation of Table.
type Table struct {
	kv    map[Value]TableEntry
	array []Value
}

// TableEntry define an entry in a Table.
type TableEntry struct {
	Key   any
	Value any
}

// Set inserts/updates value associated to given key. If provided value is nil,
// entry is deleted.
func (mt *Table) Set(k, v Value) {
	arrayK, kIsInt := k.(int)
	if kIsInt && arrayK <= len(mt.array) {
		mt.arraySet(arrayK, v)
		return
	}

	mt.mapSet(k, v)
}

// 0 <= k <= len(mt.array)
func (mt *Table) arraySet(k int, v Value) {
	// Delete.
	if v == nil {
		// migrate element after k to map.
		for i := k + 1; i < len(mt.array); i++ {
			mt.mapSetEntry(i, TableEntry{i, mt.array[i]})
		}

		// Delete k.
		mt.array = mt.array[:k]
		return
	}

	// Insert.
	if k == len(mt.array) {
		mt.array = append(mt.array, v)
		mt.fillSeqHole()
		return
	}

	// Update value in array.
	mt.array[k] = v
}

// 0 <= k <= len(mt.array)
func (mt *Table) arrayGet(k int) Value {
	return mt.array[k]
}

func (mt *Table) mapSet(k Value, v Value) {
	// Delete.
	if v == nil {
		delete(mt.kv, k)
		return
	}

	entry, entryExists := mt.kv[k]
	entry.Value = v

	// Insert.
	if !entryExists {
		entry.Key = k
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
	if kIsInt && arrayK < len(mt.array) {
		return mt.arrayGet(arrayK)
	}

	return mt.mapGet(k)
}

// Append adds given value at the end of table's sequence. A sequence starts
// with key tab[0] and ends when tab[k] is nil.
func (mt *Table) Append(v Value) int {
	mt.arraySet(len(mt.array), v)

	return len(mt.array)
}

func (mt *Table) Insert(k Value, v Value) Value {
	i, isInt := k.(int)
	if !isInt {
		return Error("failed to insert: index is not an integer")
	}

	mt.insert(i, v)
	return nil
}

func (mt *Table) insert(k int, v Value) {
	head := mt.array[:k]
	tail := mt.array[k:]
	mt.array = append(head, TableEntry{
		Key:   k,
		Value: v,
	})
	mt.array = append(mt.array, tail)
}

// SeqLen returns length of table sequence.
func (mt *Table) SeqLen() int {
	return len(mt.array)
}

// Len returns number of entries in table.
func (mt *Table) Len() int {
	return len(mt.array) + len(mt.kv)
}

// Keys implements Table.
func (mt *Table) Keys() []Value {
	keys := make([]Value, 0, len(mt.array)+len(mt.kv))

	for i := range mt.array {
		keys[i] = i
	}

	i := 0
	for k := range mt.kv {
		keys[len(mt.array)+i] = k
		i++
	}

	return keys
}

// Values implements Table.
func (mt *Table) Values() []Value {
	values := make([]Value, 0, mt.Len())

	for i, value := range mt.array {
		values[i] = value
	}

	i := 0
	for _, entry := range mt.kv {
		values[len(mt.array)+i] = entry.Value
		i++
	}

	return values
}

// Entries implements Table.
func (mt *Table) Entries() []TableEntry {
	entries := make([]TableEntry, 0, len(mt.array)+len(mt.kv))

	for i, value := range mt.array {
		entries[i] = TableEntry{i, value}
	}

	i := 0
	for _, entry := range mt.kv {
		entries[len(mt.array)+i] = entry
		i++
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

	totalKeys := len(mt.array) + len(mt.kv)

	for k, value := range mt.array {
		result.WriteString(Sexpr(value))
		if k < totalKeys-1 {
			result.WriteRune(' ')
		}
	}

	i := len(mt.array)
	for _, entry := range mt.kv {
		result.WriteString(Sexpr(entry.Key))
		result.WriteString(": ")
		result.WriteString(Sexpr(entry.Value))
		if i < totalKeys-1 {
			result.WriteRune(' ')
		}
		i++
	}

	result.WriteRune(')')

	return result.String()
}

func (mt *Table) fillSeqHole() {
	for {
		entry := mt.mapGetEntry(len(mt.array))
		if entry.Value == nil {
			break
		}

		delete(mt.kv, len(mt.array))
		mt.array = append(mt.array, entry.Value)
	}
}
