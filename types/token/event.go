package token

import (
	"apotscan/types"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/the729/lcs"
	"strconv"
)

const (
	TypeWithdrawEvent = "0x3::token::WithdrawEvent"
	TypeDepositEvent  = "0x3::token::DepositEvent"

	TypeCreateTokenDataEvent        = "0x3::token::CreateTokenDataEvent"
	TypeCollectionCreationEvent     = "0x3::token::CreateCollectionEvent"
	TypeBurnTokenEvent              = "0x3::token::BurnTokenEvent"
	TypeMutateTokenPropertyMapEvent = "0x3::token::MutateTokenPropertyMapEvent"
	TypeMintTokenEvent              = "0x3::token::MintTokenEvent"

	TypeTokenSwapEvent    = "0x3::token_coin_swap::TokenSwapEvent"
	TypeTokenListingEvent = "0x3::token_coin_swap::TokenListingEvent"

	TypeTokenOfferEvent       = "0x3::token_transfers::TokenOfferEvent"
	TypeTokenClaimEvent       = "0x3::token_transfers::TokenClaimEvent"
	TypeTokenCancelOfferEvent = "0x3::token_transfers::TokenCancelOfferEvent"
)

type TokenEvent struct {
	Key            string `json:"key"`
	SequenceNumber string `json:"sequence_number"`
	Type           string `json:"type"`
	TokenEventData TokenEventData
}

type TokenEventData interface {
	EventType() string
	//ToStruct()
}

type TokenOfferEvent struct {
	ToAddress string  `json:"to_address"`
	TokenId   TokenId `json:"token_id"`
	Amount    uint64  `json:"amount"`
}

func (TokenOfferEvent) EventType() string {
	return TypeTokenOfferEvent
}

type TokenClaimEvent struct {
	ToAddress string  `json:"to_address"`
	TokenId   TokenId `json:"token_id"`
	Amount    uint64  `json:"amount"`
}

func (TokenClaimEvent) EventType() string {
	return TypeTokenClaimEvent
}

type TokenCancelOfferEvent struct {
	ToAddress string  `json:"to_address"`
	TokenId   TokenId `json:"token_id"`
	Amount    uint64  `json:"amount"`
}

func (TokenCancelOfferEvent) EventType() string {
	return TypeTokenCancelOfferEvent
}

type TokenListingEvent struct {
	TokenId         TokenId        `json:"token_id"`
	Amount          uint64         `json:"amount"`
	MinPrice        uint64         `json:"min_price"`
	LockedUntilSecs uint64         `json:"locked_until_secs"`
	CoinTypeInfo    types.TypeInfo `json:"coin_type_info"`
}

func (TokenListingEvent) EventType() string {
	return TypeTokenListingEvent
}

type TokenSwapEvent struct {
	TokenId      TokenId        `json:"token_id"`
	TokenBuyer   string         `json:"token_buyer"`
	TokenAmount  int64          `json:"token_amount"`
	CoinAmount   int64          `json:"coin_amount"`
	CoinTypeInfo types.TypeInfo `json:"coin_type_info"`
}

func (TokenSwapEvent) EventType() string {
	return TypeTokenSwapEvent
}

type WithdrawEvent struct {
	Amount int64   `json:"amount"`
	Id     TokenId `json:"id"`
}

func (WithdrawEvent) EventType() string {
	return TypeWithdrawEvent
}

type DepositEvent struct {
	Amount uint64  `json:"amount"`
	Id     TokenId `json:"id"`
}

func (DepositEvent) EventType() string {
	return TypeDepositEvent
}

type BurnTokenEvent struct {
	Amount uint64  `json:"amount"`
	Id     TokenId `json:"id"`
}

func (BurnTokenEvent) EventType() string {
	return TypeBurnTokenEvent
}

type MintTokenEvent struct {
	Amount uint64  `json:"amount"`
	Id     TokenId `json:"id"`
}

func (MintTokenEvent) EventType() string {
	return TypeMintTokenEvent
}

type CreateTokenDataEventRaw struct {
}

type CreateTokenDataEvent struct {
	Id                       TokenDataId `json:"id"`
	Description              string      `json:"description"`
	Maximum                  uint64      `json:"maximum"`
	Uri                      string      `json:"uri"`
	RoyaltyPayeeAddress      string      `json:"royalty_payee_address"`
	RoyaltyPointsDenominator string      `json:"royalty_points_denominator"`
	RoyaltyPointsNumerator   int64       `json:"royal_points_numerator"`
	Name                     string      `json:"name"`
	MutabilityConfig         types.Value `json:"mutability_config"`
	PropertyKeys             types.Value `json:"property_keys"`
	PropertyValues           types.Value `json:"property_values"`
	PropertyTypes            types.Value `json:"property_types"`
}

func (CreateTokenDataEvent) EventType() string {
	return TypeCreateTokenDataEvent
}

type CollectionCreationEventRaw struct {
	Creator        string `json:"creator"`
	CollectionName string `json:"collection_name"`
	Uri            string `json:"uri"`
	Description    string `json:"description"`
	Maximum        string `json:"maximum"`
}

func (c CollectionCreationEventRaw) GetEvent() CollectionCreationEvent {
	maximum, _ := strconv.ParseUint(c.Maximum, 10, 64)
	return CollectionCreationEvent{
		Creator:        c.Creator,
		CollectionName: c.CollectionName,
		Uri:            c.Uri,
		Description:    c.Description,
		Maximum:        maximum,
	}
}

type CollectionCreationEvent struct {
	Creator        string `json:"creator"`
	CollectionName string `json:"collection_name"`
	Uri            string `json:"uri"`
	Description    string `json:"description"`
	Maximum        uint64 `json:"maximum"`
}

func (CollectionCreationEvent) EventType() string {
	return TypeCollectionCreationEvent
}

type MutateTokenPropertyMapEvent struct {
	OldId  TokenId     `json:"old_id"`
	NewID  TokenId     `json:"new_id"`
	Keys   types.Value `json:"keys"`
	Values types.Value `json:"values"`
	Types  types.Value `json:"types"`
}

func (MutateTokenPropertyMapEvent) EventType() string {
	return TypeMutateTokenPropertyMapEvent
}

type TokenId struct {
	TokenDataId     TokenDataId `json:"token_data_id"`
	PropertyVersion uint64      `json:"property_version"`
}

func (t TokenId) ToString() string {
	data, _ := json.Marshal(t)
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

type TokenDataId struct {
	Creator    string `json:"creator"`
	Collection string `json:"collection"`
	Name       string `json:"name"`
}

func (t TokenDataId) ToString() string {
	data, _ := json.Marshal(t)
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

var _ = lcs.RegisterEnum(
	(*TokenEventData)(nil),
	(*WithdrawEvent)(nil),
	(*DepositEvent)(nil),
	(*MintTokenEvent)(nil),
	(*CreateTokenDataEvent)(nil),
	(*CollectionCreationEvent)(nil),
	(*BurnTokenEvent)(nil),
	(*MutateTokenPropertyMapEvent)(nil),
	(*TokenListingEvent)(nil),
	(*TokenSwapEvent)(nil),
	(*TokenOfferEvent)(nil),
	(*TokenClaimEvent)(nil),
	(*TokenCancelOfferEvent)(nil),
)

//var TokenEventSet mapset.Set
//
//func Init() {
//	TokenEventSet = mapset.NewSet()
//	TokenEventSet.Add(TypeWithdrawEvent)
//	TokenEventSet.Add(TypeDepositEvent)
//	TokenEventSet.Add(TypeCreateTokenDataEvent)
//	TokenEventSet.Add(TypeCollectionCreationEvent)
//	TokenEventSet.Add(TypeBurnTokenEvent)
//	TokenEventSet.Add(TypeMutateTokenPropertyMapEvent)
//	TokenEventSet.Add(TypeMintTokenEvent)
//}
