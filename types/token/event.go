package token

import (
	"apotscan/types"
	"github.com/the729/lcs"
)

const (
	TypeWithdrawEvent               = "0x3::token::WithdrawEvent"
	TypeDepositEvent                = "0x3::token::DepositEvent"
	TypeCreateTokenDataEvent        = "0x3::token::CreateTokenDataEvent"
	TypeCollectionCreationEvent     = "0x3::token::CreateCollectionEvent"
	TypeBurnTokenEvent              = "0x3::token::BurnTokenEvent"
	TypeMutateTokenPropertyMapEvent = "0x3::token::MutateTokenPropertyMapEvent"
	TypeMintTokenEvent              = "0x3::token::MintTokenEvent"
)

type TokenEvent interface {
	EventType() string
	//ToStruct()
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

type CreateTokenDataEvent struct {
	Id                       TokenDataId `json:"id"`
	Description              string      `json:"description"`
	Maximum                  uint64      `json:"maximum"`
	Uri                      string      `json:"uri"`
	RoyaltyPayeeAddress      string      `json:"royalty_payee_address"`
	RoyaltyPointsDenominator string      `json:"royalty_points_denominator"`
	Name                     string      `json:"name"`
	MutabilityConfig         types.Value `json:"mutability_config"`
	PropertyKeys             types.Value `json:"property_keys"`
	PropertyValues           types.Value `json:"property_values"`
	PropertyTypes            types.Value `json:"property_types"`
}

func (CreateTokenDataEvent) EventType() string {
	return TypeCreateTokenDataEvent
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
	ProperTyVersion uint64      `json:"proper_ty_version"`
}

type TokenDataId struct {
	Creator    string `json:"creator"`
	Collection string `json:"collection"`
	Name       string `json:"name"`
}

var _ = lcs.RegisterEnum(
	(*TokenEvent)(nil),
	(*WithdrawEvent)(nil),
	(*DepositEvent)(nil),
	(*MintTokenEvent)(nil),
	(*CreateTokenDataEvent)(nil),
	(*CollectionCreationEvent)(nil),
	(*BurnTokenEvent)(nil),
	(*MutateTokenPropertyMapEvent)(nil),
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
