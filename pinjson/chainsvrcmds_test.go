// Copyright (c) 2014 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package pinjson_test

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/nyodeco/pind/pinjson"
	"github.com/nyodeco/pind/wire"
)

// TestChainSvrCmds tests all of the chain server commands marshal and unmarshal
// into valid results include handling of optional fields being omitted in the
// marshalled command, while optional fields with defaults have the default
// assigned on unmarshalled commands.
func TestChainSvrCmds(t *testing.T) {
	t.Parallel()

	testID := int(1)
	tests := []struct {
		name         string
		newCmd       func() (interface{}, error)
		staticCmd    func() interface{}
		marshalled   string
		unmarshalled interface{}
	}{
		{
			name: "addnode",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("addnode", "127.0.0.1", pinjson.ANRemove)
			},
			staticCmd: func() interface{} {
				return pinjson.NewAddNodeCmd("127.0.0.1", pinjson.ANRemove)
			},
			marshalled:   `{"jsonrpc":"1.0","method":"addnode","params":["127.0.0.1","remove"],"id":1}`,
			unmarshalled: &pinjson.AddNodeCmd{Addr: "127.0.0.1", SubCmd: pinjson.ANRemove},
		},
		{
			name: "createrawtransaction",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("createrawtransaction", `[{"txid":"123","vout":1}]`,
					`{"456":0.0123}`)
			},
			staticCmd: func() interface{} {
				txInputs := []pinjson.TransactionInput{
					{Txid: "123", Vout: 1},
				}
				amounts := map[string]float64{"456": .0123}
				return pinjson.NewCreateRawTransactionCmd(txInputs, amounts, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"createrawtransaction","params":[[{"txid":"123","vout":1}],{"456":0.0123}],"id":1}`,
			unmarshalled: &pinjson.CreateRawTransactionCmd{
				Inputs:  []pinjson.TransactionInput{{Txid: "123", Vout: 1}},
				Amounts: map[string]float64{"456": .0123},
			},
		},
		{
			name: "createrawtransaction - no inputs",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("createrawtransaction", `[]`, `{"456":0.0123}`)
			},
			staticCmd: func() interface{} {
				amounts := map[string]float64{"456": .0123}
				return pinjson.NewCreateRawTransactionCmd(nil, amounts, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"createrawtransaction","params":[[],{"456":0.0123}],"id":1}`,
			unmarshalled: &pinjson.CreateRawTransactionCmd{
				Inputs:  []pinjson.TransactionInput{},
				Amounts: map[string]float64{"456": .0123},
			},
		},
		{
			name: "createrawtransaction optional",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("createrawtransaction", `[{"txid":"123","vout":1}]`,
					`{"456":0.0123}`, int64(12312333333))
			},
			staticCmd: func() interface{} {
				txInputs := []pinjson.TransactionInput{
					{Txid: "123", Vout: 1},
				}
				amounts := map[string]float64{"456": .0123}
				return pinjson.NewCreateRawTransactionCmd(txInputs, amounts, pinjson.Int64(12312333333))
			},
			marshalled: `{"jsonrpc":"1.0","method":"createrawtransaction","params":[[{"txid":"123","vout":1}],{"456":0.0123},12312333333],"id":1}`,
			unmarshalled: &pinjson.CreateRawTransactionCmd{
				Inputs:   []pinjson.TransactionInput{{Txid: "123", Vout: 1}},
				Amounts:  map[string]float64{"456": .0123},
				LockTime: pinjson.Int64(12312333333),
			},
		},
		{
			name: "fundrawtransaction - empty opts",
			newCmd: func() (i interface{}, e error) {
				return pinjson.NewCmd("fundrawtransaction", "deadbeef", "{}")
			},
			staticCmd: func() interface{} {
				deadbeef, err := hex.DecodeString("deadbeef")
				if err != nil {
					panic(err)
				}
				return pinjson.NewFundRawTransactionCmd(deadbeef, pinjson.FundRawTransactionOpts{}, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"fundrawtransaction","params":["deadbeef",{}],"id":1}`,
			unmarshalled: &pinjson.FundRawTransactionCmd{
				HexTx:     "deadbeef",
				Options:   pinjson.FundRawTransactionOpts{},
				IsWitness: nil,
			},
		},
		{
			name: "fundrawtransaction - full opts",
			newCmd: func() (i interface{}, e error) {
				return pinjson.NewCmd("fundrawtransaction", "deadbeef", `{"changeAddress":"bcrt1qeeuctq9wutlcl5zatge7rjgx0k45228cxez655","changePosition":1,"change_type":"legacy","includeWatching":true,"lockUnspents":true,"feeRate":0.7,"subtractFeeFromOutputs":[0],"replaceable":true,"conf_target":8,"estimate_mode":"ECONOMICAL"}`)
			},
			staticCmd: func() interface{} {
				deadbeef, err := hex.DecodeString("deadbeef")
				if err != nil {
					panic(err)
				}
				changeAddress := "bcrt1qeeuctq9wutlcl5zatge7rjgx0k45228cxez655"
				change := 1
				changeType := pinjson.ChangeTypeLegacy
				watching := true
				lockUnspents := true
				feeRate := 0.7
				replaceable := true
				confTarget := 8

				return pinjson.NewFundRawTransactionCmd(deadbeef, pinjson.FundRawTransactionOpts{
					ChangeAddress:          &changeAddress,
					ChangePosition:         &change,
					ChangeType:             &changeType,
					IncludeWatching:        &watching,
					LockUnspents:           &lockUnspents,
					FeeRate:                &feeRate,
					SubtractFeeFromOutputs: []int{0},
					Replaceable:            &replaceable,
					ConfTarget:             &confTarget,
					EstimateMode:           &pinjson.EstimateModeEconomical,
				}, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"fundrawtransaction","params":["deadbeef",{"changeAddress":"bcrt1qeeuctq9wutlcl5zatge7rjgx0k45228cxez655","changePosition":1,"change_type":"legacy","includeWatching":true,"lockUnspents":true,"feeRate":0.7,"subtractFeeFromOutputs":[0],"replaceable":true,"conf_target":8,"estimate_mode":"ECONOMICAL"}],"id":1}`,
			unmarshalled: func() interface{} {
				changeAddress := "bcrt1qeeuctq9wutlcl5zatge7rjgx0k45228cxez655"
				change := 1
				changeType := pinjson.ChangeTypeLegacy
				watching := true
				lockUnspents := true
				feeRate := 0.7
				replaceable := true
				confTarget := 8
				return &pinjson.FundRawTransactionCmd{
					HexTx: "deadbeef",
					Options: pinjson.FundRawTransactionOpts{
						ChangeAddress:          &changeAddress,
						ChangePosition:         &change,
						ChangeType:             &changeType,
						IncludeWatching:        &watching,
						LockUnspents:           &lockUnspents,
						FeeRate:                &feeRate,
						SubtractFeeFromOutputs: []int{0},
						Replaceable:            &replaceable,
						ConfTarget:             &confTarget,
						EstimateMode:           &pinjson.EstimateModeEconomical,
					},
					IsWitness: nil,
				}
			}(),
		},
		{
			name: "fundrawtransaction - iswitness",
			newCmd: func() (i interface{}, e error) {
				return pinjson.NewCmd("fundrawtransaction", "deadbeef", "{}", true)
			},
			staticCmd: func() interface{} {
				deadbeef, err := hex.DecodeString("deadbeef")
				if err != nil {
					panic(err)
				}
				t := true
				return pinjson.NewFundRawTransactionCmd(deadbeef, pinjson.FundRawTransactionOpts{}, &t)
			},
			marshalled: `{"jsonrpc":"1.0","method":"fundrawtransaction","params":["deadbeef",{},true],"id":1}`,
			unmarshalled: &pinjson.FundRawTransactionCmd{
				HexTx:   "deadbeef",
				Options: pinjson.FundRawTransactionOpts{},
				IsWitness: func() *bool {
					t := true
					return &t
				}(),
			},
		},
		{
			name: "decoderawtransaction",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("decoderawtransaction", "123")
			},
			staticCmd: func() interface{} {
				return pinjson.NewDecodeRawTransactionCmd("123")
			},
			marshalled:   `{"jsonrpc":"1.0","method":"decoderawtransaction","params":["123"],"id":1}`,
			unmarshalled: &pinjson.DecodeRawTransactionCmd{HexTx: "123"},
		},
		{
			name: "decodescript",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("decodescript", "00")
			},
			staticCmd: func() interface{} {
				return pinjson.NewDecodeScriptCmd("00")
			},
			marshalled:   `{"jsonrpc":"1.0","method":"decodescript","params":["00"],"id":1}`,
			unmarshalled: &pinjson.DecodeScriptCmd{HexScript: "00"},
		},
		{
			name: "deriveaddresses no range",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("deriveaddresses", "00")
			},
			staticCmd: func() interface{} {
				return pinjson.NewDeriveAddressesCmd("00", nil)
			},
			marshalled:   `{"jsonrpc":"1.0","method":"deriveaddresses","params":["00"],"id":1}`,
			unmarshalled: &pinjson.DeriveAddressesCmd{Descriptor: "00"},
		},
		{
			name: "deriveaddresses int range",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd(
					"deriveaddresses", "00", pinjson.DescriptorRange{Value: 2})
			},
			staticCmd: func() interface{} {
				return pinjson.NewDeriveAddressesCmd(
					"00", &pinjson.DescriptorRange{Value: 2})
			},
			marshalled: `{"jsonrpc":"1.0","method":"deriveaddresses","params":["00",2],"id":1}`,
			unmarshalled: &pinjson.DeriveAddressesCmd{
				Descriptor: "00",
				Range:      &pinjson.DescriptorRange{Value: 2},
			},
		},
		{
			name: "deriveaddresses slice range",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd(
					"deriveaddresses", "00",
					pinjson.DescriptorRange{Value: []int{0, 2}},
				)
			},
			staticCmd: func() interface{} {
				return pinjson.NewDeriveAddressesCmd(
					"00", &pinjson.DescriptorRange{Value: []int{0, 2}})
			},
			marshalled: `{"jsonrpc":"1.0","method":"deriveaddresses","params":["00",[0,2]],"id":1}`,
			unmarshalled: &pinjson.DeriveAddressesCmd{
				Descriptor: "00",
				Range:      &pinjson.DescriptorRange{Value: []int{0, 2}},
			},
		},
		{
			name: "getaddednodeinfo",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getaddednodeinfo", true)
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetAddedNodeInfoCmd(true, nil)
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getaddednodeinfo","params":[true],"id":1}`,
			unmarshalled: &pinjson.GetAddedNodeInfoCmd{DNS: true, Node: nil},
		},
		{
			name: "getaddednodeinfo optional",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getaddednodeinfo", true, "127.0.0.1")
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetAddedNodeInfoCmd(true, pinjson.String("127.0.0.1"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getaddednodeinfo","params":[true,"127.0.0.1"],"id":1}`,
			unmarshalled: &pinjson.GetAddedNodeInfoCmd{
				DNS:  true,
				Node: pinjson.String("127.0.0.1"),
			},
		},
		{
			name: "getbestblockhash",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getbestblockhash")
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetBestBlockHashCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getbestblockhash","params":[],"id":1}`,
			unmarshalled: &pinjson.GetBestBlockHashCmd{},
		},
		{
			name: "getblock",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getblock", "123", pinjson.Int(0))
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetBlockCmd("123", pinjson.Int(0))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblock","params":["123",0],"id":1}`,
			unmarshalled: &pinjson.GetBlockCmd{
				Hash:      "123",
				Verbosity: pinjson.Int(0),
			},
		},
		{
			name: "getblock default verbosity",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getblock", "123")
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetBlockCmd("123", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblock","params":["123"],"id":1}`,
			unmarshalled: &pinjson.GetBlockCmd{
				Hash:      "123",
				Verbosity: pinjson.Int(1),
			},
		},
		{
			name: "getblock required optional1",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getblock", "123", pinjson.Int(1))
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetBlockCmd("123", pinjson.Int(1))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblock","params":["123",1],"id":1}`,
			unmarshalled: &pinjson.GetBlockCmd{
				Hash:      "123",
				Verbosity: pinjson.Int(1),
			},
		},
		{
			name: "getblock required optional2",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getblock", "123", pinjson.Int(2))
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetBlockCmd("123", pinjson.Int(2))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblock","params":["123",2],"id":1}`,
			unmarshalled: &pinjson.GetBlockCmd{
				Hash:      "123",
				Verbosity: pinjson.Int(2),
			},
		},
		{
			name: "getblockchaininfo",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getblockchaininfo")
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetBlockChainInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getblockchaininfo","params":[],"id":1}`,
			unmarshalled: &pinjson.GetBlockChainInfoCmd{},
		},
		{
			name: "getblockcount",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getblockcount")
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetBlockCountCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getblockcount","params":[],"id":1}`,
			unmarshalled: &pinjson.GetBlockCountCmd{},
		},
		{
			name: "getblockfilter",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getblockfilter", "0000afaf")
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetBlockFilterCmd("0000afaf", nil)
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getblockfilter","params":["0000afaf"],"id":1}`,
			unmarshalled: &pinjson.GetBlockFilterCmd{"0000afaf", nil},
		},
		{
			name: "getblockfilter optional filtertype",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getblockfilter", "0000afaf", "basic")
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetBlockFilterCmd("0000afaf", pinjson.NewFilterTypeName(pinjson.FilterTypeBasic))
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getblockfilter","params":["0000afaf","basic"],"id":1}`,
			unmarshalled: &pinjson.GetBlockFilterCmd{"0000afaf", pinjson.NewFilterTypeName(pinjson.FilterTypeBasic)},
		},
		{
			name: "getblockhash",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getblockhash", 123)
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetBlockHashCmd(123)
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getblockhash","params":[123],"id":1}`,
			unmarshalled: &pinjson.GetBlockHashCmd{Index: 123},
		},
		{
			name: "getblockheader",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getblockheader", "123")
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetBlockHeaderCmd("123", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblockheader","params":["123"],"id":1}`,
			unmarshalled: &pinjson.GetBlockHeaderCmd{
				Hash:    "123",
				Verbose: pinjson.Bool(true),
			},
		},
		{
			name: "getblockstats height",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getblockstats", pinjson.HashOrHeight{Value: 123})
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetBlockStatsCmd(pinjson.HashOrHeight{Value: 123}, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblockstats","params":[123],"id":1}`,
			unmarshalled: &pinjson.GetBlockStatsCmd{
				HashOrHeight: pinjson.HashOrHeight{Value: 123},
			},
		},
		{
			name: "getblockstats hash",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getblockstats", pinjson.HashOrHeight{Value: "deadbeef"})
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetBlockStatsCmd(pinjson.HashOrHeight{Value: "deadbeef"}, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblockstats","params":["deadbeef"],"id":1}`,
			unmarshalled: &pinjson.GetBlockStatsCmd{
				HashOrHeight: pinjson.HashOrHeight{Value: "deadbeef"},
			},
		},
		{
			name: "getblockstats height optional stats",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getblockstats", pinjson.HashOrHeight{Value: 123}, []string{"avgfee", "maxfee"})
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetBlockStatsCmd(pinjson.HashOrHeight{Value: 123}, &[]string{"avgfee", "maxfee"})
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblockstats","params":[123,["avgfee","maxfee"]],"id":1}`,
			unmarshalled: &pinjson.GetBlockStatsCmd{
				HashOrHeight: pinjson.HashOrHeight{Value: 123},
				Stats:        &[]string{"avgfee", "maxfee"},
			},
		},
		{
			name: "getblockstats hash optional stats",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getblockstats", pinjson.HashOrHeight{Value: "deadbeef"}, []string{"avgfee", "maxfee"})
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetBlockStatsCmd(pinjson.HashOrHeight{Value: "deadbeef"}, &[]string{"avgfee", "maxfee"})
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblockstats","params":["deadbeef",["avgfee","maxfee"]],"id":1}`,
			unmarshalled: &pinjson.GetBlockStatsCmd{
				HashOrHeight: pinjson.HashOrHeight{Value: "deadbeef"},
				Stats:        &[]string{"avgfee", "maxfee"},
			},
		},
		{
			name: "getblocktemplate",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getblocktemplate")
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetBlockTemplateCmd(nil)
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getblocktemplate","params":[],"id":1}`,
			unmarshalled: &pinjson.GetBlockTemplateCmd{Request: nil},
		},
		{
			name: "getblocktemplate optional - template request",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getblocktemplate", `{"mode":"template","capabilities":["longpoll","coinbasetxn"]}`)
			},
			staticCmd: func() interface{} {
				template := pinjson.TemplateRequest{
					Mode:         "template",
					Capabilities: []string{"longpoll", "coinbasetxn"},
				}
				return pinjson.NewGetBlockTemplateCmd(&template)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblocktemplate","params":[{"mode":"template","capabilities":["longpoll","coinbasetxn"]}],"id":1}`,
			unmarshalled: &pinjson.GetBlockTemplateCmd{
				Request: &pinjson.TemplateRequest{
					Mode:         "template",
					Capabilities: []string{"longpoll", "coinbasetxn"},
				},
			},
		},
		{
			name: "getblocktemplate optional - template request with tweaks",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getblocktemplate", `{"mode":"template","capabilities":["longpoll","coinbasetxn"],"sigoplimit":500,"sizelimit":100000000,"maxversion":2}`)
			},
			staticCmd: func() interface{} {
				template := pinjson.TemplateRequest{
					Mode:         "template",
					Capabilities: []string{"longpoll", "coinbasetxn"},
					SigOpLimit:   500,
					SizeLimit:    100000000,
					MaxVersion:   2,
				}
				return pinjson.NewGetBlockTemplateCmd(&template)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblocktemplate","params":[{"mode":"template","capabilities":["longpoll","coinbasetxn"],"sigoplimit":500,"sizelimit":100000000,"maxversion":2}],"id":1}`,
			unmarshalled: &pinjson.GetBlockTemplateCmd{
				Request: &pinjson.TemplateRequest{
					Mode:         "template",
					Capabilities: []string{"longpoll", "coinbasetxn"},
					SigOpLimit:   int64(500),
					SizeLimit:    int64(100000000),
					MaxVersion:   2,
				},
			},
		},
		{
			name: "getblocktemplate optional - template request with tweaks 2",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getblocktemplate", `{"mode":"template","capabilities":["longpoll","coinbasetxn"],"sigoplimit":true,"sizelimit":100000000,"maxversion":2}`)
			},
			staticCmd: func() interface{} {
				template := pinjson.TemplateRequest{
					Mode:         "template",
					Capabilities: []string{"longpoll", "coinbasetxn"},
					SigOpLimit:   true,
					SizeLimit:    100000000,
					MaxVersion:   2,
				}
				return pinjson.NewGetBlockTemplateCmd(&template)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblocktemplate","params":[{"mode":"template","capabilities":["longpoll","coinbasetxn"],"sigoplimit":true,"sizelimit":100000000,"maxversion":2}],"id":1}`,
			unmarshalled: &pinjson.GetBlockTemplateCmd{
				Request: &pinjson.TemplateRequest{
					Mode:         "template",
					Capabilities: []string{"longpoll", "coinbasetxn"},
					SigOpLimit:   true,
					SizeLimit:    int64(100000000),
					MaxVersion:   2,
				},
			},
		},
		{
			name: "getcfilter",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getcfilter", "123",
					wire.GCSFilterRegular)
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetCFilterCmd("123",
					wire.GCSFilterRegular)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getcfilter","params":["123",0],"id":1}`,
			unmarshalled: &pinjson.GetCFilterCmd{
				Hash:       "123",
				FilterType: wire.GCSFilterRegular,
			},
		},
		{
			name: "getcfilterheader",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getcfilterheader", "123",
					wire.GCSFilterRegular)
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetCFilterHeaderCmd("123",
					wire.GCSFilterRegular)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getcfilterheader","params":["123",0],"id":1}`,
			unmarshalled: &pinjson.GetCFilterHeaderCmd{
				Hash:       "123",
				FilterType: wire.GCSFilterRegular,
			},
		},
		{
			name: "getchaintips",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getchaintips")
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetChainTipsCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getchaintips","params":[],"id":1}`,
			unmarshalled: &pinjson.GetChainTipsCmd{},
		},
		{
			name: "getchaintxstats",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getchaintxstats")
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetChainTxStatsCmd(nil, nil)
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getchaintxstats","params":[],"id":1}`,
			unmarshalled: &pinjson.GetChainTxStatsCmd{},
		},
		{
			name: "getchaintxstats optional nblocks",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getchaintxstats", pinjson.Int32(1000))
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetChainTxStatsCmd(pinjson.Int32(1000), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getchaintxstats","params":[1000],"id":1}`,
			unmarshalled: &pinjson.GetChainTxStatsCmd{
				NBlocks: pinjson.Int32(1000),
			},
		},
		{
			name: "getchaintxstats optional nblocks and blockhash",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getchaintxstats", pinjson.Int32(1000), pinjson.String("0000afaf"))
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetChainTxStatsCmd(pinjson.Int32(1000), pinjson.String("0000afaf"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getchaintxstats","params":[1000,"0000afaf"],"id":1}`,
			unmarshalled: &pinjson.GetChainTxStatsCmd{
				NBlocks:   pinjson.Int32(1000),
				BlockHash: pinjson.String("0000afaf"),
			},
		},
		{
			name: "getconnectioncount",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getconnectioncount")
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetConnectionCountCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getconnectioncount","params":[],"id":1}`,
			unmarshalled: &pinjson.GetConnectionCountCmd{},
		},
		{
			name: "getdifficulty",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getdifficulty")
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetDifficultyCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getdifficulty","params":[],"id":1}`,
			unmarshalled: &pinjson.GetDifficultyCmd{},
		},
		{
			name: "getgenerate",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getgenerate")
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetGenerateCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getgenerate","params":[],"id":1}`,
			unmarshalled: &pinjson.GetGenerateCmd{},
		},
		{
			name: "gethashespersec",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("gethashespersec")
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetHashesPerSecCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"gethashespersec","params":[],"id":1}`,
			unmarshalled: &pinjson.GetHashesPerSecCmd{},
		},
		{
			name: "getinfo",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getinfo")
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getinfo","params":[],"id":1}`,
			unmarshalled: &pinjson.GetInfoCmd{},
		},
		{
			name: "getmempoolentry",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getmempoolentry", "txhash")
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetMempoolEntryCmd("txhash")
			},
			marshalled: `{"jsonrpc":"1.0","method":"getmempoolentry","params":["txhash"],"id":1}`,
			unmarshalled: &pinjson.GetMempoolEntryCmd{
				TxID: "txhash",
			},
		},
		{
			name: "getmempoolinfo",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getmempoolinfo")
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetMempoolInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getmempoolinfo","params":[],"id":1}`,
			unmarshalled: &pinjson.GetMempoolInfoCmd{},
		},
		{
			name: "getmininginfo",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getmininginfo")
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetMiningInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getmininginfo","params":[],"id":1}`,
			unmarshalled: &pinjson.GetMiningInfoCmd{},
		},
		{
			name: "getnetworkinfo",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getnetworkinfo")
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetNetworkInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getnetworkinfo","params":[],"id":1}`,
			unmarshalled: &pinjson.GetNetworkInfoCmd{},
		},
		{
			name: "getnettotals",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getnettotals")
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetNetTotalsCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getnettotals","params":[],"id":1}`,
			unmarshalled: &pinjson.GetNetTotalsCmd{},
		},
		{
			name: "getnetworkhashps",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getnetworkhashps")
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetNetworkHashPSCmd(nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getnetworkhashps","params":[],"id":1}`,
			unmarshalled: &pinjson.GetNetworkHashPSCmd{
				Blocks: pinjson.Int(120),
				Height: pinjson.Int(-1),
			},
		},
		{
			name: "getnetworkhashps optional1",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getnetworkhashps", 200)
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetNetworkHashPSCmd(pinjson.Int(200), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getnetworkhashps","params":[200],"id":1}`,
			unmarshalled: &pinjson.GetNetworkHashPSCmd{
				Blocks: pinjson.Int(200),
				Height: pinjson.Int(-1),
			},
		},
		{
			name: "getnetworkhashps optional2",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getnetworkhashps", 200, 123)
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetNetworkHashPSCmd(pinjson.Int(200), pinjson.Int(123))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getnetworkhashps","params":[200,123],"id":1}`,
			unmarshalled: &pinjson.GetNetworkHashPSCmd{
				Blocks: pinjson.Int(200),
				Height: pinjson.Int(123),
			},
		},
		{
			name: "getnodeaddresses",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getnodeaddresses")
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetNodeAddressesCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getnodeaddresses","params":[],"id":1}`,
			unmarshalled: &pinjson.GetNodeAddressesCmd{
				Count: pinjson.Int32(1),
			},
		},
		{
			name: "getnodeaddresses optional",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getnodeaddresses", 10)
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetNodeAddressesCmd(pinjson.Int32(10))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getnodeaddresses","params":[10],"id":1}`,
			unmarshalled: &pinjson.GetNodeAddressesCmd{
				Count: pinjson.Int32(10),
			},
		},
		{
			name: "getpeerinfo",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getpeerinfo")
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetPeerInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getpeerinfo","params":[],"id":1}`,
			unmarshalled: &pinjson.GetPeerInfoCmd{},
		},
		{
			name: "getrawmempool",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getrawmempool")
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetRawMempoolCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawmempool","params":[],"id":1}`,
			unmarshalled: &pinjson.GetRawMempoolCmd{
				Verbose: pinjson.Bool(false),
			},
		},
		{
			name: "getrawmempool optional",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getrawmempool", false)
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetRawMempoolCmd(pinjson.Bool(false))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawmempool","params":[false],"id":1}`,
			unmarshalled: &pinjson.GetRawMempoolCmd{
				Verbose: pinjson.Bool(false),
			},
		},
		{
			name: "getrawtransaction",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getrawtransaction", "123")
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetRawTransactionCmd("123", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawtransaction","params":["123"],"id":1}`,
			unmarshalled: &pinjson.GetRawTransactionCmd{
				Txid:    "123",
				Verbose: pinjson.Int(0),
			},
		},
		{
			name: "getrawtransaction optional",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getrawtransaction", "123", 1)
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetRawTransactionCmd("123", pinjson.Int(1))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawtransaction","params":["123",1],"id":1}`,
			unmarshalled: &pinjson.GetRawTransactionCmd{
				Txid:    "123",
				Verbose: pinjson.Int(1),
			},
		},
		{
			name: "gettxout",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("gettxout", "123", 1)
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetTxOutCmd("123", 1, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"gettxout","params":["123",1],"id":1}`,
			unmarshalled: &pinjson.GetTxOutCmd{
				Txid:           "123",
				Vout:           1,
				IncludeMempool: pinjson.Bool(true),
			},
		},
		{
			name: "gettxout optional",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("gettxout", "123", 1, true)
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetTxOutCmd("123", 1, pinjson.Bool(true))
			},
			marshalled: `{"jsonrpc":"1.0","method":"gettxout","params":["123",1,true],"id":1}`,
			unmarshalled: &pinjson.GetTxOutCmd{
				Txid:           "123",
				Vout:           1,
				IncludeMempool: pinjson.Bool(true),
			},
		},
		{
			name: "gettxoutproof",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("gettxoutproof", []string{"123", "456"})
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetTxOutProofCmd([]string{"123", "456"}, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"gettxoutproof","params":[["123","456"]],"id":1}`,
			unmarshalled: &pinjson.GetTxOutProofCmd{
				TxIDs: []string{"123", "456"},
			},
		},
		{
			name: "gettxoutproof optional",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("gettxoutproof", []string{"123", "456"},
					pinjson.String("000000000000034a7dedef4a161fa058a2d67a173a90155f3a2fe6fc132e0ebf"))
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetTxOutProofCmd([]string{"123", "456"},
					pinjson.String("000000000000034a7dedef4a161fa058a2d67a173a90155f3a2fe6fc132e0ebf"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"gettxoutproof","params":[["123","456"],` +
				`"000000000000034a7dedef4a161fa058a2d67a173a90155f3a2fe6fc132e0ebf"],"id":1}`,
			unmarshalled: &pinjson.GetTxOutProofCmd{
				TxIDs:     []string{"123", "456"},
				BlockHash: pinjson.String("000000000000034a7dedef4a161fa058a2d67a173a90155f3a2fe6fc132e0ebf"),
			},
		},
		{
			name: "gettxoutsetinfo",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("gettxoutsetinfo")
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetTxOutSetInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"gettxoutsetinfo","params":[],"id":1}`,
			unmarshalled: &pinjson.GetTxOutSetInfoCmd{},
		},
		{
			name: "getwork",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getwork")
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetWorkCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getwork","params":[],"id":1}`,
			unmarshalled: &pinjson.GetWorkCmd{
				Data: nil,
			},
		},
		{
			name: "getwork optional",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getwork", "00112233")
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetWorkCmd(pinjson.String("00112233"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getwork","params":["00112233"],"id":1}`,
			unmarshalled: &pinjson.GetWorkCmd{
				Data: pinjson.String("00112233"),
			},
		},
		{
			name: "help",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("help")
			},
			staticCmd: func() interface{} {
				return pinjson.NewHelpCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"help","params":[],"id":1}`,
			unmarshalled: &pinjson.HelpCmd{
				Command: nil,
			},
		},
		{
			name: "help optional",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("help", "getblock")
			},
			staticCmd: func() interface{} {
				return pinjson.NewHelpCmd(pinjson.String("getblock"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"help","params":["getblock"],"id":1}`,
			unmarshalled: &pinjson.HelpCmd{
				Command: pinjson.String("getblock"),
			},
		},
		{
			name: "invalidateblock",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("invalidateblock", "123")
			},
			staticCmd: func() interface{} {
				return pinjson.NewInvalidateBlockCmd("123")
			},
			marshalled: `{"jsonrpc":"1.0","method":"invalidateblock","params":["123"],"id":1}`,
			unmarshalled: &pinjson.InvalidateBlockCmd{
				BlockHash: "123",
			},
		},
		{
			name: "ping",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("ping")
			},
			staticCmd: func() interface{} {
				return pinjson.NewPingCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"ping","params":[],"id":1}`,
			unmarshalled: &pinjson.PingCmd{},
		},
		{
			name: "preciousblock",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("preciousblock", "0123")
			},
			staticCmd: func() interface{} {
				return pinjson.NewPreciousBlockCmd("0123")
			},
			marshalled: `{"jsonrpc":"1.0","method":"preciousblock","params":["0123"],"id":1}`,
			unmarshalled: &pinjson.PreciousBlockCmd{
				BlockHash: "0123",
			},
		},
		{
			name: "reconsiderblock",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("reconsiderblock", "123")
			},
			staticCmd: func() interface{} {
				return pinjson.NewReconsiderBlockCmd("123")
			},
			marshalled: `{"jsonrpc":"1.0","method":"reconsiderblock","params":["123"],"id":1}`,
			unmarshalled: &pinjson.ReconsiderBlockCmd{
				BlockHash: "123",
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("searchrawtransactions", "1Address")
			},
			staticCmd: func() interface{} {
				return pinjson.NewSearchRawTransactionsCmd("1Address", nil, nil, nil, nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address"],"id":1}`,
			unmarshalled: &pinjson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     pinjson.Int(1),
				Skip:        pinjson.Int(0),
				Count:       pinjson.Int(100),
				VinExtra:    pinjson.Int(0),
				Reverse:     pinjson.Bool(false),
				FilterAddrs: nil,
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("searchrawtransactions", "1Address", 0)
			},
			staticCmd: func() interface{} {
				return pinjson.NewSearchRawTransactionsCmd("1Address",
					pinjson.Int(0), nil, nil, nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0],"id":1}`,
			unmarshalled: &pinjson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     pinjson.Int(0),
				Skip:        pinjson.Int(0),
				Count:       pinjson.Int(100),
				VinExtra:    pinjson.Int(0),
				Reverse:     pinjson.Bool(false),
				FilterAddrs: nil,
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("searchrawtransactions", "1Address", 0, 5)
			},
			staticCmd: func() interface{} {
				return pinjson.NewSearchRawTransactionsCmd("1Address",
					pinjson.Int(0), pinjson.Int(5), nil, nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0,5],"id":1}`,
			unmarshalled: &pinjson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     pinjson.Int(0),
				Skip:        pinjson.Int(5),
				Count:       pinjson.Int(100),
				VinExtra:    pinjson.Int(0),
				Reverse:     pinjson.Bool(false),
				FilterAddrs: nil,
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("searchrawtransactions", "1Address", 0, 5, 10)
			},
			staticCmd: func() interface{} {
				return pinjson.NewSearchRawTransactionsCmd("1Address",
					pinjson.Int(0), pinjson.Int(5), pinjson.Int(10), nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0,5,10],"id":1}`,
			unmarshalled: &pinjson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     pinjson.Int(0),
				Skip:        pinjson.Int(5),
				Count:       pinjson.Int(10),
				VinExtra:    pinjson.Int(0),
				Reverse:     pinjson.Bool(false),
				FilterAddrs: nil,
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("searchrawtransactions", "1Address", 0, 5, 10, 1)
			},
			staticCmd: func() interface{} {
				return pinjson.NewSearchRawTransactionsCmd("1Address",
					pinjson.Int(0), pinjson.Int(5), pinjson.Int(10), pinjson.Int(1), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0,5,10,1],"id":1}`,
			unmarshalled: &pinjson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     pinjson.Int(0),
				Skip:        pinjson.Int(5),
				Count:       pinjson.Int(10),
				VinExtra:    pinjson.Int(1),
				Reverse:     pinjson.Bool(false),
				FilterAddrs: nil,
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("searchrawtransactions", "1Address", 0, 5, 10, 1, true)
			},
			staticCmd: func() interface{} {
				return pinjson.NewSearchRawTransactionsCmd("1Address",
					pinjson.Int(0), pinjson.Int(5), pinjson.Int(10), pinjson.Int(1), pinjson.Bool(true), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0,5,10,1,true],"id":1}`,
			unmarshalled: &pinjson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     pinjson.Int(0),
				Skip:        pinjson.Int(5),
				Count:       pinjson.Int(10),
				VinExtra:    pinjson.Int(1),
				Reverse:     pinjson.Bool(true),
				FilterAddrs: nil,
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("searchrawtransactions", "1Address", 0, 5, 10, 1, true, []string{"1Address"})
			},
			staticCmd: func() interface{} {
				return pinjson.NewSearchRawTransactionsCmd("1Address",
					pinjson.Int(0), pinjson.Int(5), pinjson.Int(10), pinjson.Int(1), pinjson.Bool(true), &[]string{"1Address"})
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0,5,10,1,true,["1Address"]],"id":1}`,
			unmarshalled: &pinjson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     pinjson.Int(0),
				Skip:        pinjson.Int(5),
				Count:       pinjson.Int(10),
				VinExtra:    pinjson.Int(1),
				Reverse:     pinjson.Bool(true),
				FilterAddrs: &[]string{"1Address"},
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("searchrawtransactions", "1Address", 0, 5, 10, "null", true, []string{"1Address"})
			},
			staticCmd: func() interface{} {
				return pinjson.NewSearchRawTransactionsCmd("1Address",
					pinjson.Int(0), pinjson.Int(5), pinjson.Int(10), nil, pinjson.Bool(true), &[]string{"1Address"})
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0,5,10,null,true,["1Address"]],"id":1}`,
			unmarshalled: &pinjson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     pinjson.Int(0),
				Skip:        pinjson.Int(5),
				Count:       pinjson.Int(10),
				VinExtra:    nil,
				Reverse:     pinjson.Bool(true),
				FilterAddrs: &[]string{"1Address"},
			},
		},
		{
			name: "sendrawtransaction",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("sendrawtransaction", "1122", &pinjson.AllowHighFeesOrMaxFeeRate{})
			},
			staticCmd: func() interface{} {
				return pinjson.NewSendRawTransactionCmd("1122", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendrawtransaction","params":["1122",false],"id":1}`,
			unmarshalled: &pinjson.SendRawTransactionCmd{
				HexTx: "1122",
				FeeSetting: &pinjson.AllowHighFeesOrMaxFeeRate{
					Value: pinjson.Bool(false),
				},
			},
		},
		{
			name: "sendrawtransaction optional",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("sendrawtransaction", "1122", &pinjson.AllowHighFeesOrMaxFeeRate{Value: pinjson.Bool(false)})
			},
			staticCmd: func() interface{} {
				return pinjson.NewSendRawTransactionCmd("1122", pinjson.Bool(false))
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendrawtransaction","params":["1122",false],"id":1}`,
			unmarshalled: &pinjson.SendRawTransactionCmd{
				HexTx: "1122",
				FeeSetting: &pinjson.AllowHighFeesOrMaxFeeRate{
					Value: pinjson.Bool(false),
				},
			},
		},
		{
			name: "sendrawtransaction optional, bitcoind >= 0.19.0",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("sendrawtransaction", "1122", &pinjson.AllowHighFeesOrMaxFeeRate{Value: pinjson.Int32(1234)})
			},
			staticCmd: func() interface{} {
				return pinjson.NewBitcoindSendRawTransactionCmd("1122", 1234)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendrawtransaction","params":["1122",1234],"id":1}`,
			unmarshalled: &pinjson.SendRawTransactionCmd{
				HexTx: "1122",
				FeeSetting: &pinjson.AllowHighFeesOrMaxFeeRate{
					Value: pinjson.Int32(1234),
				},
			},
		},
		{
			name: "setgenerate",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("setgenerate", true)
			},
			staticCmd: func() interface{} {
				return pinjson.NewSetGenerateCmd(true, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"setgenerate","params":[true],"id":1}`,
			unmarshalled: &pinjson.SetGenerateCmd{
				Generate:     true,
				GenProcLimit: pinjson.Int(-1),
			},
		},
		{
			name: "setgenerate optional",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("setgenerate", true, 6)
			},
			staticCmd: func() interface{} {
				return pinjson.NewSetGenerateCmd(true, pinjson.Int(6))
			},
			marshalled: `{"jsonrpc":"1.0","method":"setgenerate","params":[true,6],"id":1}`,
			unmarshalled: &pinjson.SetGenerateCmd{
				Generate:     true,
				GenProcLimit: pinjson.Int(6),
			},
		},
		{
			name: "signmessagewithprivkey",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("signmessagewithprivkey", "5Hue", "Hey")
			},
			staticCmd: func() interface{} {
				return pinjson.NewSignMessageWithPrivKey("5Hue", "Hey")
			},
			marshalled: `{"jsonrpc":"1.0","method":"signmessagewithprivkey","params":["5Hue","Hey"],"id":1}`,
			unmarshalled: &pinjson.SignMessageWithPrivKeyCmd{
				PrivKey: "5Hue",
				Message: "Hey",
			},
		},
		{
			name: "stop",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("stop")
			},
			staticCmd: func() interface{} {
				return pinjson.NewStopCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"stop","params":[],"id":1}`,
			unmarshalled: &pinjson.StopCmd{},
		},
		{
			name: "submitblock",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("submitblock", "112233")
			},
			staticCmd: func() interface{} {
				return pinjson.NewSubmitBlockCmd("112233", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"submitblock","params":["112233"],"id":1}`,
			unmarshalled: &pinjson.SubmitBlockCmd{
				HexBlock: "112233",
				Options:  nil,
			},
		},
		{
			name: "submitblock optional",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("submitblock", "112233", `{"workid":"12345"}`)
			},
			staticCmd: func() interface{} {
				options := pinjson.SubmitBlockOptions{
					WorkID: "12345",
				}
				return pinjson.NewSubmitBlockCmd("112233", &options)
			},
			marshalled: `{"jsonrpc":"1.0","method":"submitblock","params":["112233",{"workid":"12345"}],"id":1}`,
			unmarshalled: &pinjson.SubmitBlockCmd{
				HexBlock: "112233",
				Options: &pinjson.SubmitBlockOptions{
					WorkID: "12345",
				},
			},
		},
		{
			name: "uptime",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("uptime")
			},
			staticCmd: func() interface{} {
				return pinjson.NewUptimeCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"uptime","params":[],"id":1}`,
			unmarshalled: &pinjson.UptimeCmd{},
		},
		{
			name: "validateaddress",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("validateaddress", "1Address")
			},
			staticCmd: func() interface{} {
				return pinjson.NewValidateAddressCmd("1Address")
			},
			marshalled: `{"jsonrpc":"1.0","method":"validateaddress","params":["1Address"],"id":1}`,
			unmarshalled: &pinjson.ValidateAddressCmd{
				Address: "1Address",
			},
		},
		{
			name: "verifychain",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("verifychain")
			},
			staticCmd: func() interface{} {
				return pinjson.NewVerifyChainCmd(nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"verifychain","params":[],"id":1}`,
			unmarshalled: &pinjson.VerifyChainCmd{
				CheckLevel: pinjson.Int32(3),
				CheckDepth: pinjson.Int32(288),
			},
		},
		{
			name: "verifychain optional1",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("verifychain", 2)
			},
			staticCmd: func() interface{} {
				return pinjson.NewVerifyChainCmd(pinjson.Int32(2), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"verifychain","params":[2],"id":1}`,
			unmarshalled: &pinjson.VerifyChainCmd{
				CheckLevel: pinjson.Int32(2),
				CheckDepth: pinjson.Int32(288),
			},
		},
		{
			name: "verifychain optional2",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("verifychain", 2, 500)
			},
			staticCmd: func() interface{} {
				return pinjson.NewVerifyChainCmd(pinjson.Int32(2), pinjson.Int32(500))
			},
			marshalled: `{"jsonrpc":"1.0","method":"verifychain","params":[2,500],"id":1}`,
			unmarshalled: &pinjson.VerifyChainCmd{
				CheckLevel: pinjson.Int32(2),
				CheckDepth: pinjson.Int32(500),
			},
		},
		{
			name: "verifymessage",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("verifymessage", "1Address", "301234", "test")
			},
			staticCmd: func() interface{} {
				return pinjson.NewVerifyMessageCmd("1Address", "301234", "test")
			},
			marshalled: `{"jsonrpc":"1.0","method":"verifymessage","params":["1Address","301234","test"],"id":1}`,
			unmarshalled: &pinjson.VerifyMessageCmd{
				Address:   "1Address",
				Signature: "301234",
				Message:   "test",
			},
		},
		{
			name: "verifytxoutproof",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("verifytxoutproof", "test")
			},
			staticCmd: func() interface{} {
				return pinjson.NewVerifyTxOutProofCmd("test")
			},
			marshalled: `{"jsonrpc":"1.0","method":"verifytxoutproof","params":["test"],"id":1}`,
			unmarshalled: &pinjson.VerifyTxOutProofCmd{
				Proof: "test",
			},
		},
		{
			name: "getdescriptorinfo",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getdescriptorinfo", "123")
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetDescriptorInfoCmd("123")
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getdescriptorinfo","params":["123"],"id":1}`,
			unmarshalled: &pinjson.GetDescriptorInfoCmd{Descriptor: "123"},
		},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Marshal the command as created by the new static command
		// creation function.
		marshalled, err := pinjson.MarshalCmd(pinjson.RpcVersion1, testID, test.staticCmd())
		if err != nil {
			t.Errorf("MarshalCmd #%d (%s) unexpected error: %v", i,
				test.name, err)
			continue
		}

		if !bytes.Equal(marshalled, []byte(test.marshalled)) {
			t.Errorf("Test #%d (%s) unexpected marshalled data - "+
				"got %s, want %s", i, test.name, marshalled,
				test.marshalled)
			t.Errorf("\n%s\n%s", marshalled, test.marshalled)
			continue
		}

		// Ensure the command is created without error via the generic
		// new command creation function.
		cmd, err := test.newCmd()
		if err != nil {
			t.Errorf("Test #%d (%s) unexpected NewCmd error: %v ",
				i, test.name, err)
		}

		// Marshal the command as created by the generic new command
		// creation function.
		marshalled, err = pinjson.MarshalCmd(pinjson.RpcVersion1, testID, cmd)
		if err != nil {
			t.Errorf("MarshalCmd #%d (%s) unexpected error: %v", i,
				test.name, err)
			continue
		}

		if !bytes.Equal(marshalled, []byte(test.marshalled)) {
			t.Errorf("Test #%d (%s) unexpected marshalled data - "+
				"got %s, want %s", i, test.name, marshalled,
				test.marshalled)
			continue
		}

		var request pinjson.Request
		if err := json.Unmarshal(marshalled, &request); err != nil {
			t.Errorf("Test #%d (%s) unexpected error while "+
				"unmarshalling JSON-RPC request: %v", i,
				test.name, err)
			continue
		}

		cmd, err = pinjson.UnmarshalCmd(&request)
		if err != nil {
			t.Errorf("UnmarshalCmd #%d (%s) unexpected error: %v", i,
				test.name, err)
			continue
		}

		if !reflect.DeepEqual(cmd, test.unmarshalled) {
			t.Errorf("Test #%d (%s) unexpected unmarshalled command "+
				"- got %s, want %s", i, test.name,
				fmt.Sprintf("(%T) %+[1]v", cmd),
				fmt.Sprintf("(%T) %+[1]v\n", test.unmarshalled))
			continue
		}
	}
}

// TestChainSvrCmdErrors ensures any errors that occur in the command during
// custom mashal and unmarshal are as expected.
func TestChainSvrCmdErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		result     interface{}
		marshalled string
		err        error
	}{
		{
			name:       "template request with invalid type",
			result:     &pinjson.TemplateRequest{},
			marshalled: `{"mode":1}`,
			err:        &json.UnmarshalTypeError{},
		},
		{
			name:       "invalid template request sigoplimit field",
			result:     &pinjson.TemplateRequest{},
			marshalled: `{"sigoplimit":"invalid"}`,
			err:        pinjson.Error{ErrorCode: pinjson.ErrInvalidType},
		},
		{
			name:       "invalid template request sizelimit field",
			result:     &pinjson.TemplateRequest{},
			marshalled: `{"sizelimit":"invalid"}`,
			err:        pinjson.Error{ErrorCode: pinjson.ErrInvalidType},
		},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		err := json.Unmarshal([]byte(test.marshalled), &test.result)
		if reflect.TypeOf(err) != reflect.TypeOf(test.err) {
			t.Errorf("Test #%d (%s) wrong error - got %T (%v), "+
				"want %T", i, test.name, err, err, test.err)
			continue
		}

		if terr, ok := test.err.(pinjson.Error); ok {
			gotErrorCode := err.(pinjson.Error).ErrorCode
			if gotErrorCode != terr.ErrorCode {
				t.Errorf("Test #%d (%s) mismatched error code "+
					"- got %v (%v), want %v", i, test.name,
					gotErrorCode, terr, terr.ErrorCode)
				continue
			}
		}
	}
}
