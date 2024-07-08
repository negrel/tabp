package tabp

// ArgsTable wraps a table to provide easier retrieval of macros/functions args.
type ArgsTable struct {
	tab      ReadOnlyTable
	seqStart int
}

// NewArgsTable returns a new ArgsTable.
func NewArgsTable(tab ReadOnlyTable) ArgsTable {
	return ArgsTable{
		tab:      tab,
		seqStart: 0,
	}
}

func (at *ArgsTable) consumeArg(name Symbol) Value {
	v := at.tab.Get(SymbolValue(name))
	if v.Type == NilValueType {
		at.seqStart++
		v = at.tab.Get(IntValue(at.seqStart))
	}

	return v
}

// ToSExpr implements SExpr.
func (at ArgsTable) ToSExpr() string {
	return at.tab.ToSExpr()
}
