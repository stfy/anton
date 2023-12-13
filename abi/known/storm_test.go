package known_test

import (
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
	"github.com/tonindexer/anton/abi"
	"github.com/xssnick/tonutils-go/address"
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
