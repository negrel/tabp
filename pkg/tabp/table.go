package tabp

import (
	"strings"

	"golang.org/x/exp/slices"
)

// Table is a map based implementation of Table.
type Table struct {
	kv  map[Value]TableEntry
	seq []Value
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
	if kIsInt && arrayK >= 0 && arrayK <= len(mt.seq) {
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
		for i := k + 1; i < len(mt.seq); i++ {
			mt.mapSetEntry(TableEntry{i, mt.seq[i]})
		}

		// Delete k.
		mt.seq = mt.seq[:k]
		return
	}

	// Insert.
	if k == len(mt.seq) {
		mt.seq = append(mt.seq, v)
		mt.fixSeqHole()
		return
	}

	// Update value in array.
	mt.seq[k] = v
}

// 0 <= k <= len(mt.array)
func (mt *Table) arrayGet(k int) Value {
	return mt.seq[k]
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
	mt.mapSetEntry(entry)
}

func (mt *Table) mapSetEntry(entry TableEntry) {
	if mt.kv == nil {
		mt.kv = make(map[Value]TableEntry)
	}

	mt.kv[entry.Key] = entry
}

func (mt *Table) mapGet(k Value) Value {
	v, _ := mt.mapGetEntry(k)
	return v.Value
}

func (mt *Table) mapGetEntry(k Value) (TableEntry, bool) {
	if mt.kv == nil {
		return TableEntry{}, false
	}

	entry, ok := mt.kv[k]
	return entry, ok
}

// Get returns value associated with given key. A nil value is returned if key
// is not found.
func (mt *Table) Get(k Value) Value {
	arrayK, kIsInt := k.(int)
	if kIsInt && arrayK < len(mt.seq) {
		return mt.arrayGet(arrayK)
	}

	return mt.mapGet(k)
}

// Append adds given value at the end of table's sequence. A sequence starts
// with key tab[0] and ends when tab[k] is nil.
func (mt *Table) Append(v Value) int {
	mt.arraySet(len(mt.seq), v)

	return len(mt.seq)
}

// Insert inserts value v at index k (must be an integer). If index k already holds
// an entry, it is inserted at k+1.
func (mt *Table) Insert(k Value, values ...Value) Value {
	if len(values) == 0 {
		return nil
	}

	i, isInt := k.(int)
	if !isInt {
		return Error("failed to insert: index is not an integer")
	}

	mt.insert(i, values...)
	return nil
}

func (mt *Table) insert(startK int, values ...Value) {
	for i, value := range values {
		k := startK + i
		if value == nil {
			continue
		}

		if k >= 0 && k <= len(mt.seq) {
			mt.arrayInsert(k, value)
		} else {
			mt.mapInsert(k, value)
		}

		mt.fixSeqHole()
	}
}

func (mt *Table) arrayInsert(k int, v Value) {
	mt.seq = slices.Insert(mt.seq, k, v)
}

func (mt *Table) mapInsert(k int, v Value) {
	entry, ok := mt.mapGetEntry(k)

	// Entry exists, move next entry first.
	if ok {
		mt.mapInsert(k+1, entry.Value)
	}

	// Insert entry.
	mt.mapSetEntry(TableEntry{
		Key:   k,
		Value: v,
	})
}

// SeqLen returns length of table sequence.
func (mt *Table) SeqLen() int {
	return len(mt.seq)
}

// Len returns number of entries in table.
func (mt *Table) Len() int {
	return len(mt.seq) + len(mt.kv)
}

// Keys returns all keys within table.
func (mt *Table) Keys() []Value {
	keys := make([]Value, 0, len(mt.seq)+len(mt.kv))

	for i := range mt.seq {
		keys[i] = i
	}

	i := 0
	for k := range mt.kv {
		keys[len(mt.seq)+i] = k
		i++
	}

	return keys
}

// Values returns all values within table.
func (mt *Table) Values() []Value {
	values := make([]Value, 0, mt.Len())

	for i, value := range mt.seq {
		values[i] = value
	}

	i := 0
	for _, entry := range mt.kv {
		values[len(mt.seq)+i] = entry.Value
		i++
	}

	return values
}

// Entries returns all entries of table.
func (mt *Table) Entries() []TableEntry {
	entries := make([]TableEntry, 0, len(mt.seq)+len(mt.kv))

	for i, value := range mt.seq {
		entries[i] = TableEntry{i, value}
	}

	i := 0
	for _, entry := range mt.kv {
		entries[len(mt.seq)+i] = entry
		i++
	}

	return entries
}

// Map maps all entries of table using returned value from the given function.
func (mt *Table) Map(fn func(k, v Value) Value) {
	for k := 0; k < len(mt.seq); k++ {
		v := mt.seq[k]
		v = fn(k, v)
		mt.arraySet(k, v)
	}

	for k, entry := range mt.kv {
		entry.Value = fn(k, entry.Value)
		mt.mapSet(entry.Key, entry.Value)
	}
}

// ForEach iterate over all entries until there is no more left or given function
// return true.
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

	totalKeys := len(mt.seq) + len(mt.kv)

	for k, value := range mt.seq {
		result.WriteString(Sexpr(value))
		if k < totalKeys-1 {
			result.WriteRune(' ')
		}
	}

	i := len(mt.seq)
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

func (mt *Table) fixSeqHole() {
	for {
		entry, ok := mt.mapGetEntry(len(mt.seq))
		if !ok {
			break
		}

		delete(mt.kv, len(mt.seq))
		mt.seq = append(mt.seq, entry.Value)
	}
}
