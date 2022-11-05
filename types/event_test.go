package types

import (
	"encoding/binary"
	"encoding/hex"
	"testing"
)

func TestName(t *testing.T) {
	key := "0x02000000000000000ef46a87d9d404e4dfdb2b5b73e6fc0d80b832cc80e5efd294a9f2a469430e91"
	seq := key[2:18]
	add := key[18:]

	//t.Logf("%d,%v\n", num, err)
	mid, _ := hex.DecodeString(seq)
	num := binary.LittleEndian.Uint64(mid)
	t.Log(num)
	t.Log(seq)
	t.Log(add)

}
