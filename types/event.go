//package types
//
//import (
//	"encoding/binary"
//	"encoding/hex"
//	"encoding/json"
//	"strconv"
//)
//
//type Event struct {
//	Key                    string
//	SequenceNumber         int64
//	CreationNumber         int64
//	AccountAddress         string
//	TransactionVersion     int64
//	TransactionBlockHeight int64
//	Type                   string
//	Data                   []byte
//	//todo: add time
//}
//
//func (e *Event) FromEvent(event APIEvent, txVersion, txBlockHeight int64) error {
//	seqNum, err := strconv.ParseInt(event.SequenceNumber, 10, 64)
//	if err != nil {
//		return err
//	}
//	data, err := json.Marshal(&event.Data)
//	if err != nil {
//		return err
//	}
//	key := getKeyFromStr(event.Key)
//	e.Key = event.Key
//	e.CreationNumber = int64(key.CreationNumber)
//	e.AccountAddress = key.AccountAddress
//	e.TransactionVersion = txVersion
//	e.TransactionBlockHeight = txBlockHeight
//	e.SequenceNumber = seqNum
//	e.Type = event.Type
//	e.Data = data
//	return nil
//}
//
//func getKeyFromStr(key string) EventKey {
//	addr := key[18:]
//	numByte, _ := hex.DecodeString(key[2:18])
//	num := binary.LittleEndian.Uint64(numByte)
//	return EventKey{
//		AccountAddress: addr,
//		CreationNumber: num,
//	}
//}
//
//type EventKey struct {
//	AccountAddress string
//	CreationNumber uint64
//}
