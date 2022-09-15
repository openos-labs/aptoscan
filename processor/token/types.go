package token

import (
	"crypto/sha256"
	"encoding/hex"
)

type TokenTransferEvent struct {
	Id        string
	TokenId   string
	From      string
	To        string
	Amount    int64
	Version   int64
	Timestamp int64
}

func (t *TokenTransferEvent) GetId() string {
	if t.Id != "" {
		return t.Id
	}
	data := []byte(t.TokenId + t.From + t.To)
	hash := sha256.Sum256(data)
	id := hex.EncodeToString(hash[:])
	t.Id = id
	return id
}
