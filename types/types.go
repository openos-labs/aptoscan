package types

import (
	"github.com/the729/lcs"
)

const LedgerInfoKey = "ledger_info"
const MaxVersionKey = "processor:%s;chain_id_%d;max_version"

type Value interface {
	isValue()
}

type Null struct{}
type Bool bool
type Number uint64
type String string
type Array []Value
type Object map[string]Value

func (Null) isValue()   {}
func (Bool) isValue()   {}
func (Number) isValue() {}
func (String) isValue() {}
func (Array) isValue()  {}
func (Object) isValue() {}

var _ = lcs.RegisterEnum(
	(*Value)(nil),
	(*Null)(nil),
	Bool(false),
	Number(0),
	String(""),
	Array(nil),
	Object(nil),
)
