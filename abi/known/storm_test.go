package known_test

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
	"github.com/tonindexer/anton/abi"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tvm/cell"
	"os"
	"testing"
)

func Test_ReferralItem(t *testing.T) {
	var (
		interfaces []*abi.InterfaceDesc
		i          *abi.InterfaceDesc
	)

	j, err := os.ReadFile("storm.json")
	require.Nil(t, err)

	err = json.Unmarshal(j, &interfaces)
	require.Nil(t, err)

	for _, i = range interfaces {
		if i.Name == "referral" {
			err := abi.RegisterDefinitions(i.Definitions)
			require.Nil(t, err)
			break
		}
	}

	var testCases = []*struct {
		name           string
		boc            string
		code           string
		data           string
		addr           string
		expected       string
		method         string
		expectedMethod []any
	}{
		{
			addr:           "EQB3N9xY1ZgMk1b506q-7jcvSHr9_N_NqBVcygdgetUrJ-Jk",
			code:           "te6cckECFAEABAEAART/APSkE/S88sgLAQIBYgIDAgLLBAUCAVgGBwTZ0IMcA8kAB0NMDAXGwjzozMYAg1yHTHwGCEPXU61K6jyaAQNchAfpAMPpEMQH6ADDIAfoC+EUSgwf0SzD4Zds8cPhDgEDbPOBb4PpA+kAx+gAxcdch+gAx+gAwc6m0ANs8MALTHyGCEMtN3Dy6hENCQgAEdPpEMMAA8uLGgEduPz9s8yMn4QfhC+ENVAoCQEnu66ts8MPhG0NMf0x8w+ET4RUMwgJBI6PMjFsIjL4QgHHBfLhNfpA0wH6QNT0BDAE+GMC+GQB+GYB+GXbPIIImJaAcPsCcAGDBts84CGCEITc7Xq64wIhghB+Q/X8uhENCgsATu1E0NM/+kAC+GH4YiDXSY4T+kDTAfQE10wD+GMB+GT4Zfhmf+AwcAO2MTIzAYIJycOAufLRNyH6RAHy0S/4RVIQgwf0Dm+hjphbghDdWqtzyMsfWM8WAc8WcPhCWIBA2zzhMwH6APpAMAP6ADDIAqD6AvhFEoMH9Esw+GVwAYBA2zzbPBANEQT+jqAQI18DMoIK+vCAufLROPhGghAS2WAdyMsfzHBZgEDbPOAhghBplhWNuo9SMWwiMvhCAccF8uE19AT4RVgggwf0hm+lkJ9SBIMH9DYwUSGDB/R8b6XoXwP4ZfpA+gD6QDAC+kTIMgL6AvhFEoMH9EsB+GXy4TlwAYBA2zzbPBANEQwEtuAhghCHtsQ4uo85MWwiMvhCAccF8uE19AT6QDD4RVgggwf0hm+lkJ9SBIMH9DYwUSGDB/R8b6XoXwP4ZXABgELbPNs84CGCEFVtFiC64wI0A9M/JIIQX8w9FLoNEQ4PAChwgBDIywVQA88WUAP6AstqyQH7AAK2MTIzAYIQC+vCALny0Tr4QxLHBfLhLfpAMCD6RAHy0S/4RVIQgwf0Dm+h8uE7yHD6AvhFQTCDB/RLMPhl+gAw+EGCEPXU61LIyx/LPwH6AvhDzxZwWYBA2zzbPBARBP6P+TRa+ETAAfLRM/hDFMcF8uEt+kAh8BD6QHHXIfoABoIImJaAoSGUUxWgod4i1wsBwwAgkgahkTbiIML/8uE0IY6ZghAFE42RyPhDzxZQCM8WcSUESRNUR6DbPJI2MOIDjpAi8BATghDVMnbbUARtcds8kmwx4vhj2zzgMDIzExMREgAwcIAYyMsFUATPFlAE+gISy2oBzxfJAfsAADb4RvhF+ET4QcjLP/hCzxb4Q88WywH0AMzJ7VQBTAGCEC/LJqK6jpZwghCLdxc1+EHIy//4Qs8WQTCAQNs84FuED/LwEwBMcIAQyMsFUAfPFlAF+gIVy2oSyx/LPyJus5RYzxcBkTLiAckB+wDcOop6",
			data:           "te6cckEBAwEAfgAClwAAAAAAAAREgAoYScXCRmSg88WyTKkXBoMpMxXm7gJ+UZuUNbCF1P7EUAAV5PCRxhIAbzrz16mcCx0Rd95pKRQ8vTfJND7jJz5lvMABAgBDoB5TovRBPFeYylItsAV008R5Zs3C9fWsPwimaHc3SnWcgQAQAvrwgAX14QBgZAgh",
			name:           "mint_referral",
			method:         "get_referral_data",
			boc:            `te6cckEBAwEAegACj8tN3DyAAK8nhI4wkAN51569TOBY6Iu+80lIoeXpvkmh9xk58y3kAAV5PCRxhIAbzrz16mcCx0Rd95pKRQ8vTfJND7jJz5lvwAECABAC+vCABfXhAABDoB5TovRBPFeYylItsAV008R5Zs3C9fWsPwimaHc3SnWcgbh75XU=`,
			expected:       `{"nft_owner":"EQAFeTwkcYSAG8689epnAsdEXfeaSkUPL03yTQ-4yc-Zb4mt","ref_type":0,"redirect_address":"EQAFeTwkcYSAG8689epnAsdEXfeaSkUPL03yTQ-4yc-Zb4mt","parameters":{"discount":50000000,"rebate":100000000},"balances_dict":{"EQDynReiCeK8xlKRbYArpp4jyzZuF6-tYfhFM0O5ulOs5H0L":0}}`,
			expectedMethod: []any{},
		},
	}

	for _, test := range testCases {
		j := loadOperation(t, i, test.name, test.boc)

		require.Equal(t, test.expected, j)

		ret := execGetMethod(t, i, address.MustParseAddr(test.addr), test.method, test.code, test.data)

		spew.Dump(ret)
	}
}

func Test_Vault(t *testing.T) {
	var (
		interfaces []*abi.InterfaceDesc
		i          *abi.InterfaceDesc
	)

	j, err := os.ReadFile("storm.json")
	require.Nil(t, err)

	err = json.Unmarshal(j, &interfaces)
	require.Nil(t, err)

	for _, i = range interfaces {
		if i.Name == "vault" {
			err := abi.RegisterDefinitions(i.Definitions)
			require.Nil(t, err)
			break
		}
	}

	var testCases = []*struct {
		name           string
		boc            string
		code           string
		data           string
		addr           string
		expected       string
		methods        []string
		expectedMethod []any
	}{
		{
			addr:    "EQB3N9xY1ZgMk1b506q-7jcvSHr9_N_NqBVcygdgetUrJ-Jk",
			code:    "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",
			data:    "te6cckEBCQEA9AACgoAHeYCR5Ci4fR4I7farJu8gitLqyNDXPGW4XuA0WVMf+aAAAAAHc1lAAOAiGyYt2AAMPKIKlnEAACOG8m/BAAABAgEAG0EeGjAENaTpAEKbknAIAgFiBgMCASAFBABDIAC+2oHD5GHcK+scyMd5Kega/Y/4qhFseztwQK1pXDWQlABDIAUQzR5NwpfUETjq3mgphPwycXow0X0GPGSTh3KhQJqu9AIBIAgHAEMgAkfRLKzYlIpEhwtkUzEnn+lgaiT/z5FnHs2sA64rqQnUAEMgBTG+KEZbRdYwaOPmbNkQfyLv4YiSCWFGCOUA+YjIbRnEGSDwyw==",
			methods: []string{"get_vault_contract_data", "get_vault_data", "get_buffer_data"},
		},
		{
			addr:    "EQB3N9xY1ZgMk1b506q-7jcvSHr9_N_NqBVcygdgetUrJ-Jk",
			code:    "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",
			data:    "te6cckEBCAEA1gABaIAHeYCR5Ci4fR4I7farJu8gitLqyNDXPGW4XuA0WVMf+aAAAAAHc1lAAAAAAAAAAAAAAAEBAgFiBQICASAEAwBDIAC+2oHD5GHcK+scyMd5Kega/Y/4qhFseztwQK1pXDWQlABDIAUQzR5NwpfUETjq3mgphPwycXow0X0GPGSTh3KhQJqu9AIBIAcGAEMgAkfRLKzYlIpEhwtkUzEnn+lgaiT/z5FnHs2sA64rqQnUAEMgBTG+KEZbRdYwaOPmbNkQfyLv4YiSCWFGCOUA+YjIbRnEU16g3g==",
			methods: []string{"get_vault_contract_data", "get_vault_data"},
		},
	}

	for _, test := range testCases {
		for _, method := range test.methods {
			ret := execGetMethod(t, i, address.MustParseAddr(test.addr), method, test.code, test.data)

			fmt.Println(method)
			fmt.Println(ret)
		}
	}
}

func Test_PositionManager(t *testing.T) {
	var (
		interfaces []*abi.InterfaceDesc
		i          *abi.InterfaceDesc
	)

	j, err := os.ReadFile("storm.json")
	require.Nil(t, err)

	err = json.Unmarshal(j, &interfaces)
	require.Nil(t, err)

	for _, i = range interfaces {
		if i.Name == "position_manager" {
			err := abi.RegisterDefinitions(i.Definitions)
			require.Nil(t, err)
			break
		}
	}

	data := "te6cckEBBgEA7AADy4AZnKSjwXj7b7qulrljp4DKyOPHk3TDrKeywGQaWHD1jPABeEG4/wEK3O8yOGHVsDsHVl9aJuGZ5nMny4gsh3nqnBoAOlwxxRHV3b1P4G/s3Yov3TlxXQFWLhsFWZmjbsVKOV13/AQCAQARAAAAAAAAAAAgAQMf6AMAZ////////////////9+dW4+oEo9EW3qDetzRJoAAAAAAAAAAAAknwAAAAAAAAAAAMsVvq8ABAx/oBQBnAAAAAAAAAAAAAAAAIGKkcSgSj0RbeoN63NEmgAAAAAAAAAAACSfAAAAAAAAAAAAyxW+rwNAApbE="
	boc, err := base64.StdEncoding.DecodeString(data)
	require.Nil(t, err)

	c, err := cell.FromBOC(boc)
	require.Nil(t, err)

	method := getMethodDescByName(i, "get_position_manager_contract_data")

	v, err := method.ReturnValues[0].Fields.FromCell(c)
	require.Nil(t, err)

	spew.Dump(v)
}
