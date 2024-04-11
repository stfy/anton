package aggregate

import (
	"context"

	"github.com/uptrace/bun/extra/bunbig"

	"github.com/tonindexer/anton/addr"
)

type AccountsReq struct {
	Address *addr.Address

	MinterAddress *addr.Address // NFT or FT minter

	Limit int `form:"limit"`
}

type AccountsRes struct {
	// Address statistics
	TransactionsCount   int `json:"transactions_count,omitempty"`
	OwnedNFTItems       int `json:"owned_nft_items,omitempty"`
	OwnedNFTCollections int `json:"owned_nft_collections,omitempty"`
	OwnedJettonWallets  int `json:"owned_jetton_wallets,omitempty"`

	// NFT minter
	Items       int `json:"items,omitempty"`
	OwnersCount int `json:"owners_count,omitempty"`
	OwnedItems  []*struct {
		OwnerAddress *addr.Address `ch:"type:String" json:"owner_address"`
		ItemsCount   int           `json:"items_count"`
	} `json:"owned_items,omitempty"`
	UniqueOwners []*struct {
		ItemAddress *addr.Address `ch:"type:String" json:"item_address"`
		OwnersCount int           `json:"owners_count"`
	} `json:"unique_owners,omitempty"`

	// FT minter
	Wallets      int         `json:"wallets,omitempty"`
	TotalSupply  *bunbig.Int `json:"total_supply,omitempty"`
	OwnedBalance []*struct {
		WalletAddress *addr.Address `ch:"item_address,type:String" json:"wallet_address"`
		OwnerAddress  *addr.Address `ch:"type:String" json:"owner_address"`
		Balance       *bunbig.Int   `ch:"type:UInt256" json:"balance"`
	} `json:"owned_balance,omitempty"`
}

type AccountRepository interface {
	AggregateAccounts(ctx context.Context, req *AccountsReq) (*AccountsRes, error)
}
