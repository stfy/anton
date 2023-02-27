package core

import (
	"context"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/extra/bunbig"
	"github.com/uptrace/go-clickhouse/ch"
	"github.com/xssnick/tonutils-go/tlb"

	"github.com/iam047801/tonidx/abi"
	"github.com/iam047801/tonidx/internal/addr"
)

type AccountStatus string

const (
	Uninit   = AccountStatus(tlb.AccountStatusUninit)
	Active   = AccountStatus(tlb.AccountStatusActive)
	Frozen   = AccountStatus(tlb.AccountStatusFrozen)
	NonExist = AccountStatus(tlb.AccountStatusNonExist)
)

type AccountState struct {
	ch.CHModel    `ch:"account_states,partition:types,is_active,status"`
	bun.BaseModel `bun:"table:account_states"`

	Latest bool `ch:"-" json:"latest"`

	Address  addr.Address  `ch:"type:String,pk" bun:"type:bytea,pk,notnull" json:"address"`
	IsActive bool          `json:"is_active"`
	Status   AccountStatus `ch:",lc" bun:"type:account_status" json:"status"` // TODO: enum
	Balance  uint64        `json:"balance"`

	LastTxLT   uint64 `ch:",pk" bun:",pk,notnull" json:"last_tx_lt"`
	LastTxHash []byte `ch:",pk" bun:"type:bytea,unique" json:"last_tx_hash"`

	StateData *AccountData `ch:"-" bun:"rel:belongs-to,join:address=address,join:last_tx_lt=last_tx_lt" json:"state_data,omitempty"`

	StateHash []byte `bun:"type:bytea" json:"state_hash,omitempty"` // only if account is frozen
	Code      []byte `bun:"type:bytea" json:"code,omitempty"`
	CodeHash  []byte `bun:"type:bytea" json:"code_hash,omitempty"`
	Data      []byte `bun:"type:bytea" json:"data,omitempty"`
	DataHash  []byte `bun:"type:bytea" json:"data_hash,omitempty"`

	// TODO: do we need it?
	Depth uint64 `json:"depth"`
	Tick  bool   `json:"tick"`
	Tock  bool   `json:"tock"`

	// TODO: list all get method hashes
	Types []string `ch:",lc" bun:",array" json:"types,omitempty"` // TODO: ContractType here, go-ch bug
}

type NFTCollectionData struct {
	NextItemIndex uint64 `ch:"type:UInt64" json:"next_item_index,omitempty"`
	// OwnerAddress  Address
}

type NFTRoyaltyData struct {
	RoyaltyAddress *addr.Address `ch:"type:String" bun:"type:bytea" json:"royalty_address,omitempty"`
	RoyaltyFactor  uint16        `ch:"type:UInt16" json:"royalty_factor,omitempty"`
	RoyaltyBase    uint16        `ch:"type:UInt16" json:"royalty_base,omitempty"`
}

type NFTContentData struct {
	ContentURI         string `ch:"type:String" json:"content_uri,omitempty"`
	ContentName        string `ch:"type:String" json:"content_name,omitempty"`
	ContentDescription string `ch:"type:String" json:"content_description,omitempty"`
	ContentImage       string `ch:"type:String" json:"content_image,omitempty"`
	ContentImageData   []byte `ch:"type:String" json:"content_image_data,omitempty"`
}

type NFTItemData struct {
	Initialized       bool          `ch:"type:Bool" json:"initialized,omitempty"`
	ItemIndex         uint64        `ch:"type:UInt64" json:"item_index,omitempty"`
	CollectionAddress *addr.Address `ch:"type:String" bun:"type:bytea" json:"collection_address,omitempty"`
	EditorAddress     *addr.Address `ch:"type:String" bun:"type:bytea" json:"editor_address,omitempty"`
	// OwnerAddress      Address
}

type FTMasterData struct {
	TotalSupply  *bunbig.Int   `ch:"type:UInt64" bun:"type:decimal" json:"total_supply,omitempty"` // TODO: pointer here, bun bug
	Mintable     bool          `json:"mintable,omitempty"`
	AdminAddress *addr.Address `ch:"type:String" bun:"type:bytea" json:"admin_addr,omitempty"`
	// Content     nft.ContentAny
	// WalletCode  *cell.Cell
}

type FTWalletData struct {
	Balance *bunbig.Int `ch:"UInt64" bun:"type:decimal" json:"balance"`
	// OwnerAddress  Address
	MasterAddress *addr.Address `ch:"type:String" bun:"type:bytea" json:"master_address"`
	// WalletCode  *cell.Cell
}

type AccountData struct {
	ch.CHModel    `ch:"account_data,partition:types"`
	bun.BaseModel `bun:"table:account_data"`

	Address    addr.Address `ch:"type:String,pk" bun:"type:bytea,pk,notnull" json:"address"`
	LastTxLT   uint64       `ch:",pk" bun:",pk,notnull" json:"last_tx_lt"`
	LastTxHash []byte       `ch:",pk" bun:"type:bytea,notnull,unique" json:"last_tx_hash"`

	Types []string `ch:",lc" bun:",array" json:"types"` // TODO: ContractType here, ch bug

	OwnerAddress *addr.Address `ch:"type:String" bun:"type:bytea" json:"owner_address,omitempty"` // universal column for many contracts

	NFTCollectionData
	NFTRoyaltyData
	NFTContentData
	NFTItemData

	FTMasterData
	FTWalletData
}

type AccountStateFilter struct {
	Address     *addr.Address
	LatestState bool

	// contract data filter
	WithData          bool
	ContractTypes     []abi.ContractName
	OwnerAddress      *addr.Address
	CollectionAddress *addr.Address
}

type AccountRepository interface {
	AddAccountStates(ctx context.Context, tx bun.Tx, states []*AccountState) error
	AddAccountData(ctx context.Context, tx bun.Tx, data []*AccountData) error
	GetAccountStates(ctx context.Context, filter *AccountStateFilter, offset, limit int) ([]*AccountState, error)
}
