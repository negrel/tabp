//go:build goexperiment.rangefunc

package tabp

import (
	"iter"
	"strings"
)

// ReadOnlyTable define any table like object with Table read only methods.
type ReadOnlyTable interface {
	SExpr

	Has(k Value) bool
	Get(k Value) Value
	Seq() []Value
	SeqLen() int
	KVsLen() int
	Len() int
	IterSeq() iter.Seq2[int, Value]
	IterKVs() iter.Seq2[Value, Value]
	Iter() iter.Seq2[Value, Value]
}

// Table is a datastructure that acts as a map and a vector/slice at the same
// time. All values are stored in map except values that are part of the
// sequence. Entries with an integer key 'n' are stored in the slice (and part
// of the sequence) if for i from 0 to n tab.Get(i) is not nil.
type Table struct {
	kv  map[any]TableEntry
	seq []Value
}

// TableEntry define an entry in a Table.
type TableEntry struct {
	Key   Value
	Value Value
}

// Set inserts/updates value associated to given key. If provided value is nil,
// entry is deleted.
func (mt *Table) Set(k, v Value) {
	if k.Type == IntValueType {
		k := k.AsInt()
		if k >= 0 && k <= len(mt.seq) {
			mt.arraySet(k, v)
			return
		}
	}

	mt.mapSet(k, v)
}

// 0 <= k <= len(mt.array)
func (mt *Table) arraySet(k int, v Value) {
	// Delete.
	if v.Type == NilValueType {
		// migrate element after k to map.
		for i := k + 1; i < len(mt.seq); i++ {
			mt.mapSetEntry(TableEntry{IntValue(i), mt.seq[i]})
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
	if v.Type == NilValueType {
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
		mt.kv = make(map[any]TableEntry)
	}

	mt.kv[entry.Key.AsAny()] = entry
}

func (mt *Table) mapGet(k Value) Value {
	return mt.mapGetEntry(k).Value
}

func (mt *Table) mapGetEntry(k Value) TableEntry {
	if mt.kv == nil {
		return TableEntry{}
	}

	return mt.kv[k.AsAny()]
}

// Get returns value associated with given key. A nil value is returned if key
// is not found.
func (mt *Table) Get(k Value) Value {
	if k.Type == IntValueType {
		k := k.AsInt()
		if k >= 0 && k < len(mt.seq) {
			return mt.arrayGet(k)
		}
	}

	return mt.mapGet(k)
}

// Has returns whether table contains given key.
func (mt *Table) Has(k Value) bool {
	return (k.Type == IntValueType && mt.arrayHas(k.AsInt())) || mt.mapHas(k)
}

func (mt *Table) arrayHas(k int) bool {
	return k >= 0 && k < len(mt.seq)
}

func (mt *Table) mapHas(k Value) bool {
	return mt.mapGet(k).Type != NilValueType
}

// Append adds given value at the end of table's sequence. A sequence starts
// with key tab[0] and ends when tab[k] is nil.
func (mt *Table) Append(v Value) int {
	mt.arraySet(len(mt.seq), v)

	return len(mt.seq)
}

// Insert inserts value v at index k (must be an integer). If index k already holds
// an entry, it is inserted at k+1.
func (mt *Table) Insert(k int, values ...Value) {
	// TODO improve algo efficiency.
	if len(values) == 0 {
		return
	}

	mt.insert(k, values...)
}

func (mt *Table) insert(startK int, values ...Value) {
	for i, value := range values {
		if value.Type == NilValueType {
			continue
		}
		k := startK + i

		if k >= 0 && k < len(mt.seq) {
			mt.arrayInsert(k, value)
		} else {
			mt.mapInsert(k, value)
		}

		mt.fixSeqHole()
	}
}

func (mt *Table) arrayInsert(k int, v Value) {
	mt.copyEntryTo(k, 1)
	mt.arraySet(k, v)
}

func (mt *Table) mapInsert(k int, v Value) {
	// Entry exists, move next entry first.
	if mt.Has(IntValue(k)) {
		mt.copyEntryTo(k, 1)
	}

	// Insert entry.
	mt.mapSetEntry(TableEntry{
		Key:   IntValue(k),
		Value: v,
	})
}

func (mt *Table) copyEntryTo(k, delta int) {
	dstK := k + delta
	copyFromSeq := k >= 0 && k < len(mt.seq)
	copyToSeq := dstK >= 0 && dstK < len(mt.seq)
	appendToSeq := copyToSeq && dstK == len(mt.seq)

	if copyFromSeq {
		if copyToSeq { // Copy value within sequence.
			if !appendToSeq {
				mt.copyEntryTo(dstK, delta)
			}
			mt.arraySet(dstK, mt.arrayGet(k))
		} else { // Copy value from sequence to map.
			// Copy destination map value if any.
			if mt.mapGet(IntValue(dstK)).Type != NilValueType {
				mt.copyEntryTo(dstK, delta)
			}
			mt.mapSet(IntValue(dstK), mt.arrayGet(k))
		}
	} else {
		if copyToSeq { // Copy value from map to sequence.
			if !appendToSeq {
				mt.copyEntryTo(dstK, delta)
			}
			mt.arraySet(dstK, mt.mapGet(IntValue(k)))
		} else { // Copy value from map to map.
			// Copy destination map value if any.
			if mt.mapGet(IntValue(dstK)).Type != NilValueType {
				mt.copyEntryTo(dstK, delta)
			}
			mt.mapSet(IntValue(dstK), mt.mapGet(IntValue(k)))
		}
	}
}

// SeqLen returns length of table sequence.
func (mt *Table) SeqLen() int {
	return len(mt.seq)
}

// Seq returns table sequence.
func (mt *Table) Seq() []Value {
	return mt.seq
}

// IterSeq returns an iter.Seq over table sequence.
func (mt *Table) IterSeq() iter.Seq2[int, Value] {
	return func(yield func(i int, v Value) bool) {
		for i := 0; i < len(mt.seq); i++ {
			if !yield(i, mt.seq[i]) {
				return
			}
		}
	}
}

// IterKVs returns an iter.Seq over table keys and values, sequence excluded.
func (mt *Table) IterKVs() iter.Seq2[Value, Value] {
	return func(yield func(k, v Value) bool) {
		for _, v := range mt.kv {
			if !yield(v.Key, v.Value) {
				break
			}
		}
	}
}

// Iterreturns an iter.Seq over table sequence, keys and values.
func (mt *Table) Iter() iter.Seq2[Value, Value] {
	return func(yield func(k, v Value) bool) {
		for i, v := range mt.IterSeq() {
			if !yield(IntValue(i), v) {
				return
			}
		}

		for k, v := range mt.IterKVs() {
			if !yield(k, v) {
				return
			}
		}
	}
}

// Len returns number of entries in table.
func (mt *Table) KVsLen() int {
	return len(mt.kv)
}

// Len returns number of entries in table.
func (mt *Table) Len() int {
	return len(mt.seq) + len(mt.kv)
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
		entry := mt.mapGetEntry(IntValue(len(mt.seq)))
		if entry.Value.Type == NilValueType {
			break
		}

		delete(mt.kv, IntValue(len(mt.seq)))
		mt.seq = append(mt.seq, entry.Value)
	}
}
