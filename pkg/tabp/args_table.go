package tabp

// ArgsTable wraps a table to provide easier retrieval of macros/functions args.
type ArgsTable struct {
	tab      ReadOnlyTable
	seqStart int
}

// NewArgsTable returns a new ArgsTable.
func NewArgsTable(tab ReadOnlyTable) ArgsTable {
	return ArgsTable{
		tab: tab,
		// Skip function/macro name.
		seqStart: 0,
	}
}

func (at *ArgsTable) consumeArg(name Symbol) Value {
	v := at.tab.Get(name)
	if v == nil {
		at.seqStart++
		v = at.tab.Get(at.seqStart)
	}

	return v
}

// ToSExpr implements SExpr.
func (at ArgsTable) ToSExpr() string {
	return at.tab.ToSExpr()
}
