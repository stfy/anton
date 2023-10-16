package known_test

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/tonindexer/anton/abi"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/tvm/cell"
	"math/big"
	"testing"
)

type AmmData struct {
	Balance               tlb.Coins         `tlb:"." json:"balance"`
	VaultAddress          *address.Address  `tlb:"addr" json:"vault_address"`
	AssetId               uint16            `tlb:"## 16" json:"asset_id"`
	CloseOnly             bool              `tlb:"bool" json:"close_only"`
	Paused                bool              `tlb:"bool" json:"paused"`
	OracleLastPrice       tlb.Coins         `tlb:"." json:"oracle_last_price"`
	OracleLastSpread      tlb.Coins         `tlb:"." json:"oracle_last_spread"`
	OracleLastTimestamp   uint32            `tlb:"## 32" json:"oracle_last_timestamp"`
	OracleMaxDeviation    tlb.Coins         `tlb:"." json:"oracle_max_deviation"`
	OracleValidityPeriod  uint32            `tlb:"## 32" json:"oracle_validity_period"`
	OraclePublicKeysCount uint32            `tlb:"## 4" json:"oracle_public_keys_count"`
	AmmState              *AmmState         `tlb:"^" json:"amm_state"`
	Settings              *ExchangeSettings `tlb:"^" json:"settings"`
}

type AmmState struct {
	QuoteAssetReserve                    tlb.Coins `tlb:"." json:"quote_asset_reserve"`
	BaseAssetReserve                     tlb.Coins `tlb:"." json:"base_asset_reserve"`
	QuoteAssetWeight                     uint64    `tlb:"## 64" json:"quote_asset_reserve_weight"`
	TotalLongPositionSize                tlb.Coins `tlb:"." json:"total_long_position_size"`
	TotalShortPositionSize               tlb.Coins `tlb:"." json:"total_short_position_size"`
	OpenInterestLong                     tlb.Coins `tlb:"." json:"open_interest_long"`
	OpenInterestShort                    tlb.Coins `tlb:"." json:"open_interest_short"`
	LatestLongCumulativePremiumFraction  int64     `tlb:"## 64" json:"latest_long_cumulative_premium_fraction"`
	LatestShortCumulativePremiumFraction int64     `tlb:"## 64" json:"latest_short_cumulative_premium_fraction"`
	NextFundingBlockTimestamp            uint32    `tlb:"## 32" json:"next_funding_block_timestamp"`
}

type ExchangeSettings struct {
	Fee                           uint32    `tlb:"## 32" json:"fee"`
	RolloverFee                   uint32    `tlb:"## 32" json:"rollover_fee"`
	FundingPeriod                 uint32    `tlb:"## 32" json:"funding_period"`
	InitMarginRatio               uint32    `tlb:"## 32" json:"init_margin_ratio"`
	MaintenanceMarginRatio        uint32    `tlb:"## 32" json:"maintenance_margin_ratio"`
	LiquidationFeeRatio           uint32    `tlb:"## 32" json:"liquidation_fee_ratio"`
	PartialLiquidationRatio       uint32    `tlb:"## 32" json:"partial_liquidation_ratio"`
	SpreadLimit                   uint32    `tlb:"## 32" json:"spread_limit"`
	MaxPriceImpact                uint32    `tlb:"## 32" json:"max_price_impact"`
	MaxPriceSpread                uint32    `tlb:"## 32" json:"max_price_spread"`
	MaxOpenNotional               tlb.Coins `tlb:"." json:"max_open_notional"`
	FeeToStakersPercent           uint32    `tlb:"## 32" json:"fee_to_stakers_percent"`
	FundingMode                   uint8     `tlb:"## 2" json:"funding_mode"`
	MinPartialLiquidationNotional tlb.Coins `tlb:"." json:"min_partial_liquidation_notional"`
	MinLeverage                   uint32    `tlb:"## 32" json:"min_leverage"`
}

type PositionManagerData struct {
	TraderAddress *address.Address `tlb:"addr" json:"trader_address"`
	VaultAddress  *address.Address `tlb:"addr" json:"vault_address"`
	AmmAddress    *address.Address `tlb:"addr" json:"amm_address"`
	Long          *cell.Cell       `tlb:"maybe ^" json:"long"`
	Short         *cell.Cell       `tlb:"maybe ^" json:"short"`
	ReferralData  *cell.Cell       `tlb:"maybe ^" json:"referral_data"`
	LimitOrders   *cell.Dictionary `tlb:"dict" json:"limit_orders"`
	OrdersBitset  uint32           `tlb:"## 8" json:"limit_orders_bitset"`
}

type PositionData struct {
	Size                         *big.Int  `tlb:"## 128" json:"size"`
	Direction                    uint8     `tlb:"## 1" json:"direction"`
	Margin                       tlb.Coins `tlb:"." json:"margin"`
	OpenNotional                 tlb.Coins `tlb:"." json:"open_notional"`
	LastUpdatedCumulativePremium int64     `tlb:"## 64" json:"last_updated_cumulative_premium"`
	Fee                          uint64    `tlb:"## 32" json:"fee"`
	Discount                     uint64    `tlb:"## 32" json:"discount"`
	Rebate                       uint64    `tlb:"## 32" json:"rebate"`
	LastUpdatedTimestamp         uint64    `tlb:"## 32" json:"last_updated_timestamp"`
}

type PositionRecord struct {
	Locked          bool             `tlb:"bool" json:"locked"`
	RedirectAddress *address.Address `tlb:"addr" json:"redirect_address"`
	OrdersBitset    uint32           `tlb:"## 8" json:"orders_bitset"`
	Orders          *abi.Orders      `tlb:"." json:"orders"`
	LockedTimestamp uint32           `tlb:"## 32" json:"locked_timestamp"`
	Position        *PositionData    `tlb:"^" json:"state"`
}

type ReferralData struct {
	Address  *address.Address `tlb:"addr" json:"address"`
	Discount uint64           `tlb:"## 32" json:"discount"`
	Rebate   uint64           `tlb:"## 32" json:"rebate"`
}

//{
//        "name": "get_manager_data",
//        "return_values": [
//          {
//            "name": "trader_address",
//            "stack_type": "slice",
//            "format": "addr"
//          },
//          {
//            "name": "vault_address",
//            "stack_type": "slice",
//            "format": "addr"
//          },
//          {
//            "name": "amm_address",
//            "stack_type": "slice",
//            "format": "addr"
//          },
//          {
//            "name": "long",
//            "stack_type": "cell",
//            "format": "position_record"
//          },
//          {
//            "name": "short",
//            "stack_type": "cell",
//            "format": "position_record"
//          },
//          {
//            "name": "orders_dict",
//            "stack_type": "cell",
//            "format": "struct",
//            "struct_fields": [
//              {
//                "name": "orders",
//                "tlb_type": ".",
//                "format": "orders"
//              }
//            ]
//          },
//          {
//            "name": "referral_data",
//            "stack_type": "cell",
//            "format": "referral_data"
//          },
//          {
//            "name": "orders_bitset",
//            "stack_type": "int"
//          }
//        ]
//      },

func Test_AmmStateDesc(t *testing.T) {
	//showDesc := func(anyStruct any) {
	//	desc, err := abi.NewStructDesc(anyStruct)
	//	require.Nil(t, err)
	//
	//	res, err := json.Marshal(desc)
	//	require.Nil(t, err)
	//
	//	fmt.Println(string(res))
	//}
	//
	//structs := []any{
	//	//(*AmmState)(nil),
	//	//(*ExchangeSettings)(nil),
	//	//(*AmmData)(nil),
	//	//
	//	(*PositionManagerData)(nil),
	//}
	//
	//for _, str := range structs {
	//	showDesc(str)
	//}

	pm := PositionManagerData{}
	pmBoc, err := base64.StdEncoding.DecodeString("te6cckEBBAEAmAACy4ADTyzEMY0WQdBG+5/0DyGiK4mpiAxTrk/S7oRHky4abTAAkBq1xiMkjUf+ihMTkGPPkcuJKG4Zq//w/CZdH9ICW44AfrdvctSRVoEVKAmIlM7VXnbnIXSX43ZXR+0iwvc1cWsf9AECABEAAAAAAAAAACABAdwDADkmUuYR4oukO3QAAAAAASoF8gAAIBfXhAIC+vCAQGSFr0Q=")
	require.Nil(t, err)

	c, err := cell.FromBOC(pmBoc)
	require.Nil(t, err)

	err = tlb.LoadFromCell(&pm, c.BeginParse())

	require.Nil(t, err)

	res, err := json.Marshal(pm)

	fmt.Println(err)

	fmt.Println(string(res))
}
