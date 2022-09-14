package types

import (
	"fmt"
	"github.com/the729/lcs"
)

const LedgerInfoKey = "ledger_info"
const MaxVersionKey = "processor:%s;chain_id_%d;max_version"

type Value interface {
	Marshal() ([]byte, error)
}

type Null struct{}
type Bool bool
type Number uint64
type String string
type Array []Value
type Object map[string]Value

func (n *Null) Marshal() ([]byte, error) {
	return lcs.Marshal(n)
}
func (b *Bool) Marshal() ([]byte, error) {
	return lcs.Marshal(b)
}
func (n *Number) Marshal() ([]byte, error) {
	return lcs.Marshal(n)
}
func (s *String) Marshal() ([]byte, error) {
	return lcs.Marshal(s)
}
func (a *Array) Marshal() ([]byte, error) {
	return lcs.Marshal(a)
}
func (o *Object) Marshal() ([]byte, error) {
	return lcs.Marshal(o)
}

var _ = lcs.RegisterEnum(
	(*Value)(nil),
	(*Null)(nil),
	Bool(false),
	Number(0),
	String(""),
	Array(nil),
	Object(nil),
)

type TypeInfo struct {
	AccountAddress string `json:"account_address"`
	ModuleName     string `json:"module_name"`
	StructName     string `json:"struct_name"`
}

func (t TypeInfo) ToString() string {
	return fmt.Sprintf("%s::%s::%s", t.AccountAddress, t.ModuleName, t.StructName)

}
