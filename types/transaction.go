//package types
//
//import (
//	aptos "github.com/portto/aptos-go-sdk/client"
//	"gorm.io/gorm"
//	"strconv"
//	"time"
//)
//
////todo: 数据库部分
//type Transaction struct {
//	Sender         string      `json:"sender"`
//	SequenceNumber string      `json:"sequence_number"`
//	Payload        JSONPayload `json:"payload"`
//
//	Type      string     `json:"type"`
//	Timestamp string     `json:"timestamp"`
//	Events    []APIEvent `json:"events"`
//	Version   int64      `json:"version"`
//	Hash      string     `json:"hash"`
//	Success   bool       `json:"success"`
//	Changes   []Change   `json:"changes"`
//}
//
//func (transaction *Transaction) FromAptos(tx aptos.TransactionResp) error {
//	code := Code{
//		Bytecode: tx.Payload.Code.Bytecode,
//		ABI:      tx.Payload.Code.ABI,
//	}
//
//	var modules []Code
//	for _, module := range tx.Payload.Modules {
//		modules = append(modules, Code{
//			Bytecode: module.Bytecode,
//			ABI:      module.ABI,
//		})
//	}
//
//	payload := JSONPayload{
//		Type:          tx.Payload.Type,
//		TypeArguments: tx.Payload.TypeArguments,
//		Arguments:     tx.Payload.Arguments,
//		Code:          code,
//		Modules:       modules,
//		Function:      tx.Payload.Function,
//	}
//
//	var events []APIEvent
//	for _, event := range tx.Events {
//		events = append(events, APIEvent{
//			Key:            event.Key,
//			SequenceNumber: event.SequenceNumber,
//			Type:           event.Type,
//			Data:           event.Data,
//		})
//	}
//
//	var changes []Change
//	for _, change := range changes {
//		changes = append(changes, Change{
//			Type:         change.Type,
//			StateKeyHash: change.StateKeyHash,
//			Address:      change.Address,
//			Module:       change.Module,
//			Resource:     change.Resource,
//			Data: struct {
//				Handle   string                 `json:"handle"`
//				Key      string                 `json:"key"`
//				Value    string                 `json:"value"`
//				Bytecode string                 `json:"bytecode"`
//				ABI      interface{}            `json:"abi"`
//				Type     string                 `json:"type"`
//				Data     map[string]interface{} `json:"data"`
//			}{
//				Handle:   change.Data.Handle,
//				Key:      change.Data.Key,
//				Value:    change.Data.Value,
//				Bytecode: change.Data.Bytecode,
//				ABI:      change.Data.ABI,
//				Type:     change.Data.Type,
//				Data:     change.Data.Data,
//			},
//		})
//	}
//	version, err := strconv.ParseInt(tx.Version, 10, 64)
//	if err != nil {
//		return err
//	}
//	*transaction = Transaction{
//		Sender:         tx.Sender,
//		SequenceNumber: tx.SequenceNumber,
//		Payload:        payload,
//		Type:           tx.Type,
//		Timestamp:      tx.Timestamp,
//		Events:         events,
//		Version:        version,
//		Hash:           tx.Hash,
//		Success:        tx.Success,
//		Changes:        changes,
//	}
//	return nil
//}
//
//type JSONPayload struct {
//	Type          string   `json:"type"`
//	TypeArguments []string `json:"type_arguments"`
//	Arguments     []string `json:"arguments"`
//
//	// ScriptPayload
//	Code Code `json:"code,omitempty"`
//	// ModuleBundlePayload
//	Modules []Code `json:"modules,omitempty"`
//	// EntryFunctionPayload
//	Function string `json:"function,omitempty"`
//}
//
//type Code struct {
//	Bytecode string      `json:"bytecode"`
//	ABI      interface{} `json:"abi,omitempty"`
//}
//
//type Change struct {
//	Type         string `json:"type"`
//	StateKeyHash string `json:"state_key_hash"`
//	Address      string `json:"address"`
//	Module       string `json:"module"`
//	Resource     string `json:"resource"`
//	Data         struct {
//		Handle   string                 `json:"handle"`
//		Key      string                 `json:"key"`
//		Value    string                 `json:"value"`
//		Bytecode string                 `json:"bytecode"`
//		ABI      interface{}            `json:"abi"`
//		Type     string                 `json:"type"`
//		Data     map[string]interface{} `json:"data"`
//	} `json:"data"`
//}
//
//type APIEvent struct {
//	Key            string                 `json:"key"`
//	SequenceNumber string                 `json:"sequence_number"`
//	Type           string                 `json:"type"`
//	Data           map[string]interface{} `json:"data"`
//}
//
//type TransactionInDB struct {
//	Type                string
//	Payload             []byte
//	Version             int64
//	Hash                string
//	StateRootHash       string
//	EventRootHash       string
//	GasUsed             int64
//	Success             bool
//	VMStatus            string
//	AccumulatorRootHash string
//
//	CreatedAt *time.Time `gorm:"autoCreateTime"`
//	UpdatedAt *time.Time `gorm:"autoUpdateTime;not null"`
//}
//
//func (TransactionInDB) TableName() string {
//	return "transactions"
//}
//
//func AutoCreateTransactionTable(db *gorm.DB) error {
//	err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").AutoMigrate(&TransactionInDB{})
//	if err != nil {
//		return err
//	}
//	return nil
//}
