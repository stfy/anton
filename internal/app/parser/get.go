package parser

import (
	"context"

	"github.com/pkg/errors"
	"github.com/uptrace/bun/extra/bunbig"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton/nft"

	"github.com/iam047801/tonidx/abi"
	"github.com/iam047801/tonidx/internal/addr"
	"github.com/iam047801/tonidx/internal/core"
)

func mapContentDataNFT(ret *core.AccountData, c nft.ContentAny) {
	switch content := c.(type) {
	case *nft.ContentSemichain: // TODO: remove this (?)
		ret.ContentURI = content.URI
		ret.ContentName = content.Name
		ret.ContentDescription = content.Description
		ret.ContentImage = content.Image
		ret.ContentImageData = content.ImageData

	case *nft.ContentOnchain:
		ret.ContentName = content.Name
		ret.ContentDescription = content.Description
		ret.ContentImage = content.Image
		ret.ContentImageData = content.ImageData

	case *nft.ContentOffchain:
		ret.ContentURI = content.URI
	}
}

func mapCollectionDataNFT(ret *core.AccountData, data *abi.NFTCollectionData) {
	ret.NextItemIndex = bunbig.FromMathBig(data.NextItemIndex)
	ret.OwnerAddress, _ = new(addr.Address).FromTU(data.OwnerAddress)
	mapContentDataNFT(ret, data.Content)
}

func mapRoyaltyDataNFT(ret *core.AccountData, params *abi.NFTRoyaltyData) {
	var err error
	ret.RoyaltyAddress, err = new(addr.Address).FromTU(params.Address)
	if err != nil {
		ret.Errors = append(ret.Errors, errors.Wrap(err, "nft royalty address from TU").Error())
	}
	ret.RoyaltyBase = params.Base
	ret.RoyaltyFactor = params.Factor
}

func mapItemDataNFT(ret *core.AccountData, data *abi.NFTItemData) {
	var err error
	ret.Initialized = data.Initialized
	ret.ItemIndex = bunbig.FromMathBig(data.Index)
	ret.CollectionAddress, err = new(addr.Address).FromTU(data.CollectionAddress)
	if err != nil {
		ret.Errors = append(ret.Errors, errors.Wrap(err, "nft item collection_address from TU").Error())
	}
	ret.OwnerAddress, _ = new(addr.Address).FromTU(data.OwnerAddress)
	mapContentDataNFT(ret, data.Content)
}

func mapEditorDataNFT(ret *core.AccountData, data *abi.NFTEditableData) {
	ret.EditorAddress, _ = new(addr.Address).FromTU(data.Editor)
}

func (s *Service) getAccountDataNFT(ctx context.Context, b *tlb.BlockInfo, a *address.Address, types []abi.ContractName, ret *core.AccountData) (ok bool) {
	var unknown int

	for _, t := range types {
		switch t {
		case abi.NFTCollection:
			data, err := abi.GetNFTCollectionData(ctx, s.api, b, a)
			if err != nil {
				ret.Errors = append(ret.Errors, errors.Wrap(err, "get nft collection data").Error())
				continue
			}
			mapCollectionDataNFT(ret, data)

		case abi.NFTRoyalty:
			data, err := abi.GetNFTRoyaltyData(ctx, s.api, b, a)
			if err != nil {
				ret.Errors = append(ret.Errors, errors.Wrap(err, "get nft royalty data").Error())
				continue
			}
			mapRoyaltyDataNFT(ret, data)

		case abi.NFTItem:
			data, err := abi.GetNFTItemData(ctx, s.api, b, a)
			if err != nil {
				ret.Errors = append(ret.Errors, errors.Wrap(err, "get nft item data").Error())
				continue
			}
			mapItemDataNFT(ret, data)

		case abi.NFTEditable:
			data, err := abi.GetNFTEditableData(ctx, s.api, b, a)
			if err != nil {
				ret.Errors = append(ret.Errors, errors.Wrap(err, "get nft editable data").Error())
				continue
			}
			mapEditorDataNFT(ret, data)

		default:
			unknown++
		}
	}

	return unknown != len(types)
}

func (s *Service) getAccountDataFT(ctx context.Context, b *tlb.BlockInfo, a *address.Address, types []abi.ContractName, ret *core.AccountData) (ok bool) {
	var unknown int

	for _, t := range types {
		switch t {
		case abi.JettonMinter:
			data, err := abi.GetJettonData(ctx, s.api, b, a)
			if err != nil {
				ret.Errors = append(ret.Errors, errors.Wrap(err, "get jetton minter data").Error())
				continue
			}
			if data.TotalSupply != nil {
				ret.TotalSupply = bunbig.FromMathBig(data.TotalSupply)
			}
			ret.Mintable = data.Mintable
			ret.AdminAddress, err = new(addr.Address).FromTU(data.AdminAddr)
			if err != nil {
				ret.Errors = append(ret.Errors, errors.Wrap(err, "jetton minter admin addr from TU").Error())
			}
			mapContentDataNFT(ret, data.Content)

		case abi.JettonWallet:
			data, err := abi.GetJettonWalletData(ctx, s.api, b, a)
			if err != nil {
				ret.Errors = append(ret.Errors, errors.Wrap(err, "get jetton wallet data").Error())
				continue
			}
			if data.Balance != nil {
				ret.JettonBalance = bunbig.FromMathBig(data.Balance)
			}
			ret.OwnerAddress, err = new(addr.Address).FromTU(data.OwnerAddress)
			if err != nil {
				ret.Errors = append(ret.Errors, errors.Wrap(err, "jetton wallet owner addr from TU").Error())
			}
			ret.MasterAddress, err = new(addr.Address).FromTU(data.MasterAddress)
			if err != nil {
				ret.Errors = append(ret.Errors, errors.Wrap(err, "jetton wallet master addr from TU").Error())
			}

		default:
			unknown++
		}
	}

	return unknown != len(types)
}