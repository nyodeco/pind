// Copyright (c) 2014-2020 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package pinjson_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/nyodeco/pind/pinjson"
	"github.com/nyodeco/pinutil"
)

// TestWalletSvrCmds tests all of the wallet server commands marshal and
// unmarshal into valid results include handling of optional fields being
// omitted in the marshalled command, while optional fields with defaults have
// the default assigned on unmarshalled commands.
func TestWalletSvrCmds(t *testing.T) {
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
			name: "addmultisigaddress",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("addmultisigaddress", 2, []string{"031234", "035678"})
			},
			staticCmd: func() interface{} {
				keys := []string{"031234", "035678"}
				return pinjson.NewAddMultisigAddressCmd(2, keys, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"addmultisigaddress","params":[2,["031234","035678"]],"id":1}`,
			unmarshalled: &pinjson.AddMultisigAddressCmd{
				NRequired: 2,
				Keys:      []string{"031234", "035678"},
				Account:   nil,
			},
		},
		{
			name: "addmultisigaddress optional",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("addmultisigaddress", 2, []string{"031234", "035678"}, "test")
			},
			staticCmd: func() interface{} {
				keys := []string{"031234", "035678"}
				return pinjson.NewAddMultisigAddressCmd(2, keys, pinjson.String("test"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"addmultisigaddress","params":[2,["031234","035678"],"test"],"id":1}`,
			unmarshalled: &pinjson.AddMultisigAddressCmd{
				NRequired: 2,
				Keys:      []string{"031234", "035678"},
				Account:   pinjson.String("test"),
			},
		},
		{
			name: "createwallet",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("createwallet", "mywallet", true, true, "secret", true)
			},
			staticCmd: func() interface{} {
				return pinjson.NewCreateWalletCmd("mywallet",
					pinjson.Bool(true), pinjson.Bool(true),
					pinjson.String("secret"), pinjson.Bool(true))
			},
			marshalled: `{"jsonrpc":"1.0","method":"createwallet","params":["mywallet",true,true,"secret",true],"id":1}`,
			unmarshalled: &pinjson.CreateWalletCmd{
				WalletName:         "mywallet",
				DisablePrivateKeys: pinjson.Bool(true),
				Blank:              pinjson.Bool(true),
				Passphrase:         pinjson.String("secret"),
				AvoidReuse:         pinjson.Bool(true),
			},
		},
		{
			name: "createwallet - optional1",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("createwallet", "mywallet")
			},
			staticCmd: func() interface{} {
				return pinjson.NewCreateWalletCmd("mywallet",
					nil, nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"createwallet","params":["mywallet"],"id":1}`,
			unmarshalled: &pinjson.CreateWalletCmd{
				WalletName:         "mywallet",
				DisablePrivateKeys: pinjson.Bool(false),
				Blank:              pinjson.Bool(false),
				Passphrase:         pinjson.String(""),
				AvoidReuse:         pinjson.Bool(false),
			},
		},
		{
			name: "createwallet - optional2",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("createwallet", "mywallet", "null", "null", "secret")
			},
			staticCmd: func() interface{} {
				return pinjson.NewCreateWalletCmd("mywallet",
					nil, nil, pinjson.String("secret"), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"createwallet","params":["mywallet",null,null,"secret"],"id":1}`,
			unmarshalled: &pinjson.CreateWalletCmd{
				WalletName:         "mywallet",
				DisablePrivateKeys: nil,
				Blank:              nil,
				Passphrase:         pinjson.String("secret"),
				AvoidReuse:         pinjson.Bool(false),
			},
		},
		{
			name: "addwitnessaddress",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("addwitnessaddress", "1address")
			},
			staticCmd: func() interface{} {
				return pinjson.NewAddWitnessAddressCmd("1address")
			},
			marshalled: `{"jsonrpc":"1.0","method":"addwitnessaddress","params":["1address"],"id":1}`,
			unmarshalled: &pinjson.AddWitnessAddressCmd{
				Address: "1address",
			},
		},
		{
			name: "backupwallet",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("backupwallet", "backup.dat")
			},
			staticCmd: func() interface{} {
				return pinjson.NewBackupWalletCmd("backup.dat")
			},
			marshalled:   `{"jsonrpc":"1.0","method":"backupwallet","params":["backup.dat"],"id":1}`,
			unmarshalled: &pinjson.BackupWalletCmd{Destination: "backup.dat"},
		},
		{
			name: "loadwallet",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("loadwallet", "wallet.dat")
			},
			staticCmd: func() interface{} {
				return pinjson.NewLoadWalletCmd("wallet.dat")
			},
			marshalled:   `{"jsonrpc":"1.0","method":"loadwallet","params":["wallet.dat"],"id":1}`,
			unmarshalled: &pinjson.LoadWalletCmd{WalletName: "wallet.dat"},
		},
		{
			name: "unloadwallet",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("unloadwallet", "default")
			},
			staticCmd: func() interface{} {
				return pinjson.NewUnloadWalletCmd(pinjson.String("default"))
			},
			marshalled:   `{"jsonrpc":"1.0","method":"unloadwallet","params":["default"],"id":1}`,
			unmarshalled: &pinjson.UnloadWalletCmd{WalletName: pinjson.String("default")},
		},
		{name: "unloadwallet - nil arg",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("unloadwallet")
			},
			staticCmd: func() interface{} {
				return pinjson.NewUnloadWalletCmd(nil)
			},
			marshalled:   `{"jsonrpc":"1.0","method":"unloadwallet","params":[],"id":1}`,
			unmarshalled: &pinjson.UnloadWalletCmd{WalletName: nil},
		},
		{
			name: "createmultisig",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("createmultisig", 2, []string{"031234", "035678"})
			},
			staticCmd: func() interface{} {
				keys := []string{"031234", "035678"}
				return pinjson.NewCreateMultisigCmd(2, keys)
			},
			marshalled: `{"jsonrpc":"1.0","method":"createmultisig","params":[2,["031234","035678"]],"id":1}`,
			unmarshalled: &pinjson.CreateMultisigCmd{
				NRequired: 2,
				Keys:      []string{"031234", "035678"},
			},
		},
		{
			name: "dumpprivkey",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("dumpprivkey", "1Address")
			},
			staticCmd: func() interface{} {
				return pinjson.NewDumpPrivKeyCmd("1Address")
			},
			marshalled: `{"jsonrpc":"1.0","method":"dumpprivkey","params":["1Address"],"id":1}`,
			unmarshalled: &pinjson.DumpPrivKeyCmd{
				Address: "1Address",
			},
		},
		{
			name: "encryptwallet",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("encryptwallet", "pass")
			},
			staticCmd: func() interface{} {
				return pinjson.NewEncryptWalletCmd("pass")
			},
			marshalled: `{"jsonrpc":"1.0","method":"encryptwallet","params":["pass"],"id":1}`,
			unmarshalled: &pinjson.EncryptWalletCmd{
				Passphrase: "pass",
			},
		},
		{
			name: "estimatefee",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("estimatefee", 6)
			},
			staticCmd: func() interface{} {
				return pinjson.NewEstimateFeeCmd(6)
			},
			marshalled: `{"jsonrpc":"1.0","method":"estimatefee","params":[6],"id":1}`,
			unmarshalled: &pinjson.EstimateFeeCmd{
				NumBlocks: 6,
			},
		},
		{
			name: "estimatesmartfee - no mode",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("estimatesmartfee", 6)
			},
			staticCmd: func() interface{} {
				return pinjson.NewEstimateSmartFeeCmd(6, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"estimatesmartfee","params":[6],"id":1}`,
			unmarshalled: &pinjson.EstimateSmartFeeCmd{
				ConfTarget:   6,
				EstimateMode: &pinjson.EstimateModeConservative,
			},
		},
		{
			name: "estimatesmartfee - economical mode",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("estimatesmartfee", 6, pinjson.EstimateModeEconomical)
			},
			staticCmd: func() interface{} {
				return pinjson.NewEstimateSmartFeeCmd(6, &pinjson.EstimateModeEconomical)
			},
			marshalled: `{"jsonrpc":"1.0","method":"estimatesmartfee","params":[6,"ECONOMICAL"],"id":1}`,
			unmarshalled: &pinjson.EstimateSmartFeeCmd{
				ConfTarget:   6,
				EstimateMode: &pinjson.EstimateModeEconomical,
			},
		},
		{
			name: "estimatepriority",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("estimatepriority", 6)
			},
			staticCmd: func() interface{} {
				return pinjson.NewEstimatePriorityCmd(6)
			},
			marshalled: `{"jsonrpc":"1.0","method":"estimatepriority","params":[6],"id":1}`,
			unmarshalled: &pinjson.EstimatePriorityCmd{
				NumBlocks: 6,
			},
		},
		{
			name: "getaccount",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getaccount", "1Address")
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetAccountCmd("1Address")
			},
			marshalled: `{"jsonrpc":"1.0","method":"getaccount","params":["1Address"],"id":1}`,
			unmarshalled: &pinjson.GetAccountCmd{
				Address: "1Address",
			},
		},
		{
			name: "getaccountaddress",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getaccountaddress", "acct")
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetAccountAddressCmd("acct")
			},
			marshalled: `{"jsonrpc":"1.0","method":"getaccountaddress","params":["acct"],"id":1}`,
			unmarshalled: &pinjson.GetAccountAddressCmd{
				Account: "acct",
			},
		},
		{
			name: "getaddressesbyaccount",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getaddressesbyaccount", "acct")
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetAddressesByAccountCmd("acct")
			},
			marshalled: `{"jsonrpc":"1.0","method":"getaddressesbyaccount","params":["acct"],"id":1}`,
			unmarshalled: &pinjson.GetAddressesByAccountCmd{
				Account: "acct",
			},
		},
		{
			name: "getaddressinfo",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getaddressinfo", "1234")
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetAddressInfoCmd("1234")
			},
			marshalled: `{"jsonrpc":"1.0","method":"getaddressinfo","params":["1234"],"id":1}`,
			unmarshalled: &pinjson.GetAddressInfoCmd{
				Address: "1234",
			},
		},
		{
			name: "getbalance",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getbalance")
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetBalanceCmd(nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getbalance","params":[],"id":1}`,
			unmarshalled: &pinjson.GetBalanceCmd{
				Account: nil,
				MinConf: pinjson.Int(1),
			},
		},
		{
			name: "getbalance optional1",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getbalance", "acct")
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetBalanceCmd(pinjson.String("acct"), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getbalance","params":["acct"],"id":1}`,
			unmarshalled: &pinjson.GetBalanceCmd{
				Account: pinjson.String("acct"),
				MinConf: pinjson.Int(1),
			},
		},
		{
			name: "getbalance optional2",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getbalance", "acct", 6)
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetBalanceCmd(pinjson.String("acct"), pinjson.Int(6))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getbalance","params":["acct",6],"id":1}`,
			unmarshalled: &pinjson.GetBalanceCmd{
				Account: pinjson.String("acct"),
				MinConf: pinjson.Int(6),
			},
		},
		{
			name: "getbalances",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getbalances")
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetBalancesCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getbalances","params":[],"id":1}`,
			unmarshalled: &pinjson.GetBalancesCmd{},
		},
		{
			name: "getnewaddress",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getnewaddress")
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetNewAddressCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getnewaddress","params":[],"id":1}`,
			unmarshalled: &pinjson.GetNewAddressCmd{
				Account: nil,
			},
		},
		{
			name: "getnewaddress optional",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getnewaddress", "acct")
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetNewAddressCmd(pinjson.String("acct"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getnewaddress","params":["acct"],"id":1}`,
			unmarshalled: &pinjson.GetNewAddressCmd{
				Account: pinjson.String("acct"),
			},
		},
		{
			name: "getrawchangeaddress",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getrawchangeaddress")
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetRawChangeAddressCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawchangeaddress","params":[],"id":1}`,
			unmarshalled: &pinjson.GetRawChangeAddressCmd{
				Account: nil,
			},
		},
		{
			name: "getrawchangeaddress optional",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getrawchangeaddress", "acct")
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetRawChangeAddressCmd(pinjson.String("acct"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawchangeaddress","params":["acct"],"id":1}`,
			unmarshalled: &pinjson.GetRawChangeAddressCmd{
				Account: pinjson.String("acct"),
			},
		},
		{
			name: "getreceivedbyaccount",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getreceivedbyaccount", "acct")
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetReceivedByAccountCmd("acct", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getreceivedbyaccount","params":["acct"],"id":1}`,
			unmarshalled: &pinjson.GetReceivedByAccountCmd{
				Account: "acct",
				MinConf: pinjson.Int(1),
			},
		},
		{
			name: "getreceivedbyaccount optional",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getreceivedbyaccount", "acct", 6)
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetReceivedByAccountCmd("acct", pinjson.Int(6))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getreceivedbyaccount","params":["acct",6],"id":1}`,
			unmarshalled: &pinjson.GetReceivedByAccountCmd{
				Account: "acct",
				MinConf: pinjson.Int(6),
			},
		},
		{
			name: "getreceivedbyaddress",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getreceivedbyaddress", "1Address")
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetReceivedByAddressCmd("1Address", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getreceivedbyaddress","params":["1Address"],"id":1}`,
			unmarshalled: &pinjson.GetReceivedByAddressCmd{
				Address: "1Address",
				MinConf: pinjson.Int(1),
			},
		},
		{
			name: "getreceivedbyaddress optional",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getreceivedbyaddress", "1Address", 6)
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetReceivedByAddressCmd("1Address", pinjson.Int(6))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getreceivedbyaddress","params":["1Address",6],"id":1}`,
			unmarshalled: &pinjson.GetReceivedByAddressCmd{
				Address: "1Address",
				MinConf: pinjson.Int(6),
			},
		},
		{
			name: "gettransaction",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("gettransaction", "123")
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetTransactionCmd("123", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"gettransaction","params":["123"],"id":1}`,
			unmarshalled: &pinjson.GetTransactionCmd{
				Txid:             "123",
				IncludeWatchOnly: pinjson.Bool(false),
			},
		},
		{
			name: "gettransaction optional",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("gettransaction", "123", true)
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetTransactionCmd("123", pinjson.Bool(true))
			},
			marshalled: `{"jsonrpc":"1.0","method":"gettransaction","params":["123",true],"id":1}`,
			unmarshalled: &pinjson.GetTransactionCmd{
				Txid:             "123",
				IncludeWatchOnly: pinjson.Bool(true),
			},
		},
		{
			name: "getwalletinfo",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("getwalletinfo")
			},
			staticCmd: func() interface{} {
				return pinjson.NewGetWalletInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getwalletinfo","params":[],"id":1}`,
			unmarshalled: &pinjson.GetWalletInfoCmd{},
		},
		{
			name: "importprivkey",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("importprivkey", "abc")
			},
			staticCmd: func() interface{} {
				return pinjson.NewImportPrivKeyCmd("abc", nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"importprivkey","params":["abc"],"id":1}`,
			unmarshalled: &pinjson.ImportPrivKeyCmd{
				PrivKey: "abc",
				Label:   nil,
				Rescan:  pinjson.Bool(true),
			},
		},
		{
			name: "importprivkey optional1",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("importprivkey", "abc", "label")
			},
			staticCmd: func() interface{} {
				return pinjson.NewImportPrivKeyCmd("abc", pinjson.String("label"), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"importprivkey","params":["abc","label"],"id":1}`,
			unmarshalled: &pinjson.ImportPrivKeyCmd{
				PrivKey: "abc",
				Label:   pinjson.String("label"),
				Rescan:  pinjson.Bool(true),
			},
		},
		{
			name: "importprivkey optional2",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("importprivkey", "abc", "label", false)
			},
			staticCmd: func() interface{} {
				return pinjson.NewImportPrivKeyCmd("abc", pinjson.String("label"), pinjson.Bool(false))
			},
			marshalled: `{"jsonrpc":"1.0","method":"importprivkey","params":["abc","label",false],"id":1}`,
			unmarshalled: &pinjson.ImportPrivKeyCmd{
				PrivKey: "abc",
				Label:   pinjson.String("label"),
				Rescan:  pinjson.Bool(false),
			},
		},
		{
			name: "keypoolrefill",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("keypoolrefill")
			},
			staticCmd: func() interface{} {
				return pinjson.NewKeyPoolRefillCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"keypoolrefill","params":[],"id":1}`,
			unmarshalled: &pinjson.KeyPoolRefillCmd{
				NewSize: pinjson.Uint(100),
			},
		},
		{
			name: "keypoolrefill optional",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("keypoolrefill", 200)
			},
			staticCmd: func() interface{} {
				return pinjson.NewKeyPoolRefillCmd(pinjson.Uint(200))
			},
			marshalled: `{"jsonrpc":"1.0","method":"keypoolrefill","params":[200],"id":1}`,
			unmarshalled: &pinjson.KeyPoolRefillCmd{
				NewSize: pinjson.Uint(200),
			},
		},
		{
			name: "listaccounts",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("listaccounts")
			},
			staticCmd: func() interface{} {
				return pinjson.NewListAccountsCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listaccounts","params":[],"id":1}`,
			unmarshalled: &pinjson.ListAccountsCmd{
				MinConf: pinjson.Int(1),
			},
		},
		{
			name: "listaccounts optional",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("listaccounts", 6)
			},
			staticCmd: func() interface{} {
				return pinjson.NewListAccountsCmd(pinjson.Int(6))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listaccounts","params":[6],"id":1}`,
			unmarshalled: &pinjson.ListAccountsCmd{
				MinConf: pinjson.Int(6),
			},
		},
		{
			name: "listaddressgroupings",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("listaddressgroupings")
			},
			staticCmd: func() interface{} {
				return pinjson.NewListAddressGroupingsCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"listaddressgroupings","params":[],"id":1}`,
			unmarshalled: &pinjson.ListAddressGroupingsCmd{},
		},
		{
			name: "listlockunspent",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("listlockunspent")
			},
			staticCmd: func() interface{} {
				return pinjson.NewListLockUnspentCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"listlockunspent","params":[],"id":1}`,
			unmarshalled: &pinjson.ListLockUnspentCmd{},
		},
		{
			name: "listreceivedbyaccount",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("listreceivedbyaccount")
			},
			staticCmd: func() interface{} {
				return pinjson.NewListReceivedByAccountCmd(nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaccount","params":[],"id":1}`,
			unmarshalled: &pinjson.ListReceivedByAccountCmd{
				MinConf:          pinjson.Int(1),
				IncludeEmpty:     pinjson.Bool(false),
				IncludeWatchOnly: pinjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaccount optional1",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("listreceivedbyaccount", 6)
			},
			staticCmd: func() interface{} {
				return pinjson.NewListReceivedByAccountCmd(pinjson.Int(6), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaccount","params":[6],"id":1}`,
			unmarshalled: &pinjson.ListReceivedByAccountCmd{
				MinConf:          pinjson.Int(6),
				IncludeEmpty:     pinjson.Bool(false),
				IncludeWatchOnly: pinjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaccount optional2",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("listreceivedbyaccount", 6, true)
			},
			staticCmd: func() interface{} {
				return pinjson.NewListReceivedByAccountCmd(pinjson.Int(6), pinjson.Bool(true), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaccount","params":[6,true],"id":1}`,
			unmarshalled: &pinjson.ListReceivedByAccountCmd{
				MinConf:          pinjson.Int(6),
				IncludeEmpty:     pinjson.Bool(true),
				IncludeWatchOnly: pinjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaccount optional3",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("listreceivedbyaccount", 6, true, false)
			},
			staticCmd: func() interface{} {
				return pinjson.NewListReceivedByAccountCmd(pinjson.Int(6), pinjson.Bool(true), pinjson.Bool(false))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaccount","params":[6,true,false],"id":1}`,
			unmarshalled: &pinjson.ListReceivedByAccountCmd{
				MinConf:          pinjson.Int(6),
				IncludeEmpty:     pinjson.Bool(true),
				IncludeWatchOnly: pinjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaddress",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("listreceivedbyaddress")
			},
			staticCmd: func() interface{} {
				return pinjson.NewListReceivedByAddressCmd(nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaddress","params":[],"id":1}`,
			unmarshalled: &pinjson.ListReceivedByAddressCmd{
				MinConf:          pinjson.Int(1),
				IncludeEmpty:     pinjson.Bool(false),
				IncludeWatchOnly: pinjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaddress optional1",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("listreceivedbyaddress", 6)
			},
			staticCmd: func() interface{} {
				return pinjson.NewListReceivedByAddressCmd(pinjson.Int(6), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaddress","params":[6],"id":1}`,
			unmarshalled: &pinjson.ListReceivedByAddressCmd{
				MinConf:          pinjson.Int(6),
				IncludeEmpty:     pinjson.Bool(false),
				IncludeWatchOnly: pinjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaddress optional2",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("listreceivedbyaddress", 6, true)
			},
			staticCmd: func() interface{} {
				return pinjson.NewListReceivedByAddressCmd(pinjson.Int(6), pinjson.Bool(true), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaddress","params":[6,true],"id":1}`,
			unmarshalled: &pinjson.ListReceivedByAddressCmd{
				MinConf:          pinjson.Int(6),
				IncludeEmpty:     pinjson.Bool(true),
				IncludeWatchOnly: pinjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaddress optional3",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("listreceivedbyaddress", 6, true, false)
			},
			staticCmd: func() interface{} {
				return pinjson.NewListReceivedByAddressCmd(pinjson.Int(6), pinjson.Bool(true), pinjson.Bool(false))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaddress","params":[6,true,false],"id":1}`,
			unmarshalled: &pinjson.ListReceivedByAddressCmd{
				MinConf:          pinjson.Int(6),
				IncludeEmpty:     pinjson.Bool(true),
				IncludeWatchOnly: pinjson.Bool(false),
			},
		},
		{
			name: "listsinceblock",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("listsinceblock")
			},
			staticCmd: func() interface{} {
				return pinjson.NewListSinceBlockCmd(nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listsinceblock","params":[],"id":1}`,
			unmarshalled: &pinjson.ListSinceBlockCmd{
				BlockHash:           nil,
				TargetConfirmations: pinjson.Int(1),
				IncludeWatchOnly:    pinjson.Bool(false),
			},
		},
		{
			name: "listsinceblock optional1",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("listsinceblock", "123")
			},
			staticCmd: func() interface{} {
				return pinjson.NewListSinceBlockCmd(pinjson.String("123"), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listsinceblock","params":["123"],"id":1}`,
			unmarshalled: &pinjson.ListSinceBlockCmd{
				BlockHash:           pinjson.String("123"),
				TargetConfirmations: pinjson.Int(1),
				IncludeWatchOnly:    pinjson.Bool(false),
			},
		},
		{
			name: "listsinceblock optional2",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("listsinceblock", "123", 6)
			},
			staticCmd: func() interface{} {
				return pinjson.NewListSinceBlockCmd(pinjson.String("123"), pinjson.Int(6), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listsinceblock","params":["123",6],"id":1}`,
			unmarshalled: &pinjson.ListSinceBlockCmd{
				BlockHash:           pinjson.String("123"),
				TargetConfirmations: pinjson.Int(6),
				IncludeWatchOnly:    pinjson.Bool(false),
			},
		},
		{
			name: "listsinceblock optional3",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("listsinceblock", "123", 6, true)
			},
			staticCmd: func() interface{} {
				return pinjson.NewListSinceBlockCmd(pinjson.String("123"), pinjson.Int(6), pinjson.Bool(true))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listsinceblock","params":["123",6,true],"id":1}`,
			unmarshalled: &pinjson.ListSinceBlockCmd{
				BlockHash:           pinjson.String("123"),
				TargetConfirmations: pinjson.Int(6),
				IncludeWatchOnly:    pinjson.Bool(true),
			},
		},
		{
			name: "listsinceblock pad null",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("listsinceblock", "null", 1, false)
			},
			staticCmd: func() interface{} {
				return pinjson.NewListSinceBlockCmd(nil, pinjson.Int(1), pinjson.Bool(false))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listsinceblock","params":[null,1,false],"id":1}`,
			unmarshalled: &pinjson.ListSinceBlockCmd{
				BlockHash:           nil,
				TargetConfirmations: pinjson.Int(1),
				IncludeWatchOnly:    pinjson.Bool(false),
			},
		},
		{
			name: "listtransactions",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("listtransactions")
			},
			staticCmd: func() interface{} {
				return pinjson.NewListTransactionsCmd(nil, nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listtransactions","params":[],"id":1}`,
			unmarshalled: &pinjson.ListTransactionsCmd{
				Account:          nil,
				Count:            pinjson.Int(10),
				From:             pinjson.Int(0),
				IncludeWatchOnly: pinjson.Bool(false),
			},
		},
		{
			name: "listtransactions optional1",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("listtransactions", "acct")
			},
			staticCmd: func() interface{} {
				return pinjson.NewListTransactionsCmd(pinjson.String("acct"), nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listtransactions","params":["acct"],"id":1}`,
			unmarshalled: &pinjson.ListTransactionsCmd{
				Account:          pinjson.String("acct"),
				Count:            pinjson.Int(10),
				From:             pinjson.Int(0),
				IncludeWatchOnly: pinjson.Bool(false),
			},
		},
		{
			name: "listtransactions optional2",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("listtransactions", "acct", 20)
			},
			staticCmd: func() interface{} {
				return pinjson.NewListTransactionsCmd(pinjson.String("acct"), pinjson.Int(20), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listtransactions","params":["acct",20],"id":1}`,
			unmarshalled: &pinjson.ListTransactionsCmd{
				Account:          pinjson.String("acct"),
				Count:            pinjson.Int(20),
				From:             pinjson.Int(0),
				IncludeWatchOnly: pinjson.Bool(false),
			},
		},
		{
			name: "listtransactions optional3",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("listtransactions", "acct", 20, 1)
			},
			staticCmd: func() interface{} {
				return pinjson.NewListTransactionsCmd(pinjson.String("acct"), pinjson.Int(20),
					pinjson.Int(1), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listtransactions","params":["acct",20,1],"id":1}`,
			unmarshalled: &pinjson.ListTransactionsCmd{
				Account:          pinjson.String("acct"),
				Count:            pinjson.Int(20),
				From:             pinjson.Int(1),
				IncludeWatchOnly: pinjson.Bool(false),
			},
		},
		{
			name: "listtransactions optional4",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("listtransactions", "acct", 20, 1, true)
			},
			staticCmd: func() interface{} {
				return pinjson.NewListTransactionsCmd(pinjson.String("acct"), pinjson.Int(20),
					pinjson.Int(1), pinjson.Bool(true))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listtransactions","params":["acct",20,1,true],"id":1}`,
			unmarshalled: &pinjson.ListTransactionsCmd{
				Account:          pinjson.String("acct"),
				Count:            pinjson.Int(20),
				From:             pinjson.Int(1),
				IncludeWatchOnly: pinjson.Bool(true),
			},
		},
		{
			name: "listunspent",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("listunspent")
			},
			staticCmd: func() interface{} {
				return pinjson.NewListUnspentCmd(nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listunspent","params":[],"id":1}`,
			unmarshalled: &pinjson.ListUnspentCmd{
				MinConf:   pinjson.Int(1),
				MaxConf:   pinjson.Int(9999999),
				Addresses: nil,
			},
		},
		{
			name: "listunspent optional1",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("listunspent", 6)
			},
			staticCmd: func() interface{} {
				return pinjson.NewListUnspentCmd(pinjson.Int(6), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listunspent","params":[6],"id":1}`,
			unmarshalled: &pinjson.ListUnspentCmd{
				MinConf:   pinjson.Int(6),
				MaxConf:   pinjson.Int(9999999),
				Addresses: nil,
			},
		},
		{
			name: "listunspent optional2",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("listunspent", 6, 100)
			},
			staticCmd: func() interface{} {
				return pinjson.NewListUnspentCmd(pinjson.Int(6), pinjson.Int(100), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listunspent","params":[6,100],"id":1}`,
			unmarshalled: &pinjson.ListUnspentCmd{
				MinConf:   pinjson.Int(6),
				MaxConf:   pinjson.Int(100),
				Addresses: nil,
			},
		},
		{
			name: "listunspent optional3",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("listunspent", 6, 100, []string{"1Address", "1Address2"})
			},
			staticCmd: func() interface{} {
				return pinjson.NewListUnspentCmd(pinjson.Int(6), pinjson.Int(100),
					&[]string{"1Address", "1Address2"})
			},
			marshalled: `{"jsonrpc":"1.0","method":"listunspent","params":[6,100,["1Address","1Address2"]],"id":1}`,
			unmarshalled: &pinjson.ListUnspentCmd{
				MinConf:   pinjson.Int(6),
				MaxConf:   pinjson.Int(100),
				Addresses: &[]string{"1Address", "1Address2"},
			},
		},
		{
			name: "lockunspent",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("lockunspent", true, `[{"txid":"123","vout":1}]`)
			},
			staticCmd: func() interface{} {
				txInputs := []pinjson.TransactionInput{
					{Txid: "123", Vout: 1},
				}
				return pinjson.NewLockUnspentCmd(true, txInputs)
			},
			marshalled: `{"jsonrpc":"1.0","method":"lockunspent","params":[true,[{"txid":"123","vout":1}]],"id":1}`,
			unmarshalled: &pinjson.LockUnspentCmd{
				Unlock: true,
				Transactions: []pinjson.TransactionInput{
					{Txid: "123", Vout: 1},
				},
			},
		},
		{
			name: "move",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("move", "from", "to", 0.5)
			},
			staticCmd: func() interface{} {
				return pinjson.NewMoveCmd("from", "to", 0.5, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"move","params":["from","to",0.5],"id":1}`,
			unmarshalled: &pinjson.MoveCmd{
				FromAccount: "from",
				ToAccount:   "to",
				Amount:      0.5,
				MinConf:     pinjson.Int(1),
				Comment:     nil,
			},
		},
		{
			name: "move optional1",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("move", "from", "to", 0.5, 6)
			},
			staticCmd: func() interface{} {
				return pinjson.NewMoveCmd("from", "to", 0.5, pinjson.Int(6), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"move","params":["from","to",0.5,6],"id":1}`,
			unmarshalled: &pinjson.MoveCmd{
				FromAccount: "from",
				ToAccount:   "to",
				Amount:      0.5,
				MinConf:     pinjson.Int(6),
				Comment:     nil,
			},
		},
		{
			name: "move optional2",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("move", "from", "to", 0.5, 6, "comment")
			},
			staticCmd: func() interface{} {
				return pinjson.NewMoveCmd("from", "to", 0.5, pinjson.Int(6), pinjson.String("comment"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"move","params":["from","to",0.5,6,"comment"],"id":1}`,
			unmarshalled: &pinjson.MoveCmd{
				FromAccount: "from",
				ToAccount:   "to",
				Amount:      0.5,
				MinConf:     pinjson.Int(6),
				Comment:     pinjson.String("comment"),
			},
		},
		{
			name: "sendfrom",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("sendfrom", "from", "1Address", 0.5)
			},
			staticCmd: func() interface{} {
				return pinjson.NewSendFromCmd("from", "1Address", 0.5, nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendfrom","params":["from","1Address",0.5],"id":1}`,
			unmarshalled: &pinjson.SendFromCmd{
				FromAccount: "from",
				ToAddress:   "1Address",
				Amount:      0.5,
				MinConf:     pinjson.Int(1),
				Comment:     nil,
				CommentTo:   nil,
			},
		},
		{
			name: "sendfrom optional1",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("sendfrom", "from", "1Address", 0.5, 6)
			},
			staticCmd: func() interface{} {
				return pinjson.NewSendFromCmd("from", "1Address", 0.5, pinjson.Int(6), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendfrom","params":["from","1Address",0.5,6],"id":1}`,
			unmarshalled: &pinjson.SendFromCmd{
				FromAccount: "from",
				ToAddress:   "1Address",
				Amount:      0.5,
				MinConf:     pinjson.Int(6),
				Comment:     nil,
				CommentTo:   nil,
			},
		},
		{
			name: "sendfrom optional2",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("sendfrom", "from", "1Address", 0.5, 6, "comment")
			},
			staticCmd: func() interface{} {
				return pinjson.NewSendFromCmd("from", "1Address", 0.5, pinjson.Int(6),
					pinjson.String("comment"), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendfrom","params":["from","1Address",0.5,6,"comment"],"id":1}`,
			unmarshalled: &pinjson.SendFromCmd{
				FromAccount: "from",
				ToAddress:   "1Address",
				Amount:      0.5,
				MinConf:     pinjson.Int(6),
				Comment:     pinjson.String("comment"),
				CommentTo:   nil,
			},
		},
		{
			name: "sendfrom optional3",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("sendfrom", "from", "1Address", 0.5, 6, "comment", "commentto")
			},
			staticCmd: func() interface{} {
				return pinjson.NewSendFromCmd("from", "1Address", 0.5, pinjson.Int(6),
					pinjson.String("comment"), pinjson.String("commentto"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendfrom","params":["from","1Address",0.5,6,"comment","commentto"],"id":1}`,
			unmarshalled: &pinjson.SendFromCmd{
				FromAccount: "from",
				ToAddress:   "1Address",
				Amount:      0.5,
				MinConf:     pinjson.Int(6),
				Comment:     pinjson.String("comment"),
				CommentTo:   pinjson.String("commentto"),
			},
		},
		{
			name: "sendmany",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("sendmany", "from", `{"1Address":0.5}`)
			},
			staticCmd: func() interface{} {
				amounts := map[string]float64{"1Address": 0.5}
				return pinjson.NewSendManyCmd("from", amounts, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendmany","params":["from",{"1Address":0.5}],"id":1}`,
			unmarshalled: &pinjson.SendManyCmd{
				FromAccount: "from",
				Amounts:     map[string]float64{"1Address": 0.5},
				MinConf:     pinjson.Int(1),
				Comment:     nil,
			},
		},
		{
			name: "sendmany optional1",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("sendmany", "from", `{"1Address":0.5}`, 6)
			},
			staticCmd: func() interface{} {
				amounts := map[string]float64{"1Address": 0.5}
				return pinjson.NewSendManyCmd("from", amounts, pinjson.Int(6), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendmany","params":["from",{"1Address":0.5},6],"id":1}`,
			unmarshalled: &pinjson.SendManyCmd{
				FromAccount: "from",
				Amounts:     map[string]float64{"1Address": 0.5},
				MinConf:     pinjson.Int(6),
				Comment:     nil,
			},
		},
		{
			name: "sendmany optional2",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("sendmany", "from", `{"1Address":0.5}`, 6, "comment")
			},
			staticCmd: func() interface{} {
				amounts := map[string]float64{"1Address": 0.5}
				return pinjson.NewSendManyCmd("from", amounts, pinjson.Int(6), pinjson.String("comment"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendmany","params":["from",{"1Address":0.5},6,"comment"],"id":1}`,
			unmarshalled: &pinjson.SendManyCmd{
				FromAccount: "from",
				Amounts:     map[string]float64{"1Address": 0.5},
				MinConf:     pinjson.Int(6),
				Comment:     pinjson.String("comment"),
			},
		},
		{
			name: "sendtoaddress",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("sendtoaddress", "1Address", 0.5)
			},
			staticCmd: func() interface{} {
				return pinjson.NewSendToAddressCmd("1Address", 0.5, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendtoaddress","params":["1Address",0.5],"id":1}`,
			unmarshalled: &pinjson.SendToAddressCmd{
				Address:   "1Address",
				Amount:    0.5,
				Comment:   nil,
				CommentTo: nil,
			},
		},
		{
			name: "sendtoaddress optional1",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("sendtoaddress", "1Address", 0.5, "comment", "commentto")
			},
			staticCmd: func() interface{} {
				return pinjson.NewSendToAddressCmd("1Address", 0.5, pinjson.String("comment"),
					pinjson.String("commentto"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendtoaddress","params":["1Address",0.5,"comment","commentto"],"id":1}`,
			unmarshalled: &pinjson.SendToAddressCmd{
				Address:   "1Address",
				Amount:    0.5,
				Comment:   pinjson.String("comment"),
				CommentTo: pinjson.String("commentto"),
			},
		},
		{
			name: "setaccount",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("setaccount", "1Address", "acct")
			},
			staticCmd: func() interface{} {
				return pinjson.NewSetAccountCmd("1Address", "acct")
			},
			marshalled: `{"jsonrpc":"1.0","method":"setaccount","params":["1Address","acct"],"id":1}`,
			unmarshalled: &pinjson.SetAccountCmd{
				Address: "1Address",
				Account: "acct",
			},
		},
		{
			name: "settxfee",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("settxfee", 0.0001)
			},
			staticCmd: func() interface{} {
				return pinjson.NewSetTxFeeCmd(0.0001)
			},
			marshalled: `{"jsonrpc":"1.0","method":"settxfee","params":[0.0001],"id":1}`,
			unmarshalled: &pinjson.SetTxFeeCmd{
				Amount: 0.0001,
			},
		},
		{
			name: "signmessage",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("signmessage", "1Address", "message")
			},
			staticCmd: func() interface{} {
				return pinjson.NewSignMessageCmd("1Address", "message")
			},
			marshalled: `{"jsonrpc":"1.0","method":"signmessage","params":["1Address","message"],"id":1}`,
			unmarshalled: &pinjson.SignMessageCmd{
				Address: "1Address",
				Message: "message",
			},
		},
		{
			name: "signrawtransaction",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("signrawtransaction", "001122")
			},
			staticCmd: func() interface{} {
				return pinjson.NewSignRawTransactionCmd("001122", nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"signrawtransaction","params":["001122"],"id":1}`,
			unmarshalled: &pinjson.SignRawTransactionCmd{
				RawTx:    "001122",
				Inputs:   nil,
				PrivKeys: nil,
				Flags:    pinjson.String("ALL"),
			},
		},
		{
			name: "signrawtransaction optional1",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("signrawtransaction", "001122", `[{"txid":"123","vout":1,"scriptPubKey":"00","redeemScript":"01"}]`)
			},
			staticCmd: func() interface{} {
				txInputs := []pinjson.RawTxInput{
					{
						Txid:         "123",
						Vout:         1,
						ScriptPubKey: "00",
						RedeemScript: "01",
					},
				}

				return pinjson.NewSignRawTransactionCmd("001122", &txInputs, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"signrawtransaction","params":["001122",[{"txid":"123","vout":1,"scriptPubKey":"00","redeemScript":"01"}]],"id":1}`,
			unmarshalled: &pinjson.SignRawTransactionCmd{
				RawTx: "001122",
				Inputs: &[]pinjson.RawTxInput{
					{
						Txid:         "123",
						Vout:         1,
						ScriptPubKey: "00",
						RedeemScript: "01",
					},
				},
				PrivKeys: nil,
				Flags:    pinjson.String("ALL"),
			},
		},
		{
			name: "signrawtransaction optional2",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("signrawtransaction", "001122", `[]`, `["abc"]`)
			},
			staticCmd: func() interface{} {
				txInputs := []pinjson.RawTxInput{}
				privKeys := []string{"abc"}
				return pinjson.NewSignRawTransactionCmd("001122", &txInputs, &privKeys, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"signrawtransaction","params":["001122",[],["abc"]],"id":1}`,
			unmarshalled: &pinjson.SignRawTransactionCmd{
				RawTx:    "001122",
				Inputs:   &[]pinjson.RawTxInput{},
				PrivKeys: &[]string{"abc"},
				Flags:    pinjson.String("ALL"),
			},
		},
		{
			name: "signrawtransaction optional3",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("signrawtransaction", "001122", `[]`, `[]`, "ALL")
			},
			staticCmd: func() interface{} {
				txInputs := []pinjson.RawTxInput{}
				privKeys := []string{}
				return pinjson.NewSignRawTransactionCmd("001122", &txInputs, &privKeys,
					pinjson.String("ALL"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"signrawtransaction","params":["001122",[],[],"ALL"],"id":1}`,
			unmarshalled: &pinjson.SignRawTransactionCmd{
				RawTx:    "001122",
				Inputs:   &[]pinjson.RawTxInput{},
				PrivKeys: &[]string{},
				Flags:    pinjson.String("ALL"),
			},
		},
		{
			name: "signrawtransactionwithwallet",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("signrawtransactionwithwallet", "001122")
			},
			staticCmd: func() interface{} {
				return pinjson.NewSignRawTransactionWithWalletCmd("001122", nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"signrawtransactionwithwallet","params":["001122"],"id":1}`,
			unmarshalled: &pinjson.SignRawTransactionWithWalletCmd{
				RawTx:       "001122",
				Inputs:      nil,
				SigHashType: pinjson.String("ALL"),
			},
		},
		{
			name: "signrawtransactionwithwallet optional1",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("signrawtransactionwithwallet", "001122", `[{"txid":"123","vout":1,"scriptPubKey":"00","redeemScript":"01","witnessScript":"02","amount":1.5}]`)
			},
			staticCmd: func() interface{} {
				txInputs := []pinjson.RawTxWitnessInput{
					{
						Txid:          "123",
						Vout:          1,
						ScriptPubKey:  "00",
						RedeemScript:  pinjson.String("01"),
						WitnessScript: pinjson.String("02"),
						Amount:        pinjson.Float64(1.5),
					},
				}

				return pinjson.NewSignRawTransactionWithWalletCmd("001122", &txInputs, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"signrawtransactionwithwallet","params":["001122",[{"txid":"123","vout":1,"scriptPubKey":"00","redeemScript":"01","witnessScript":"02","amount":1.5}]],"id":1}`,
			unmarshalled: &pinjson.SignRawTransactionWithWalletCmd{
				RawTx: "001122",
				Inputs: &[]pinjson.RawTxWitnessInput{
					{
						Txid:          "123",
						Vout:          1,
						ScriptPubKey:  "00",
						RedeemScript:  pinjson.String("01"),
						WitnessScript: pinjson.String("02"),
						Amount:        pinjson.Float64(1.5),
					},
				},
				SigHashType: pinjson.String("ALL"),
			},
		},
		{
			name: "signrawtransactionwithwallet optional1 with blank fields in input",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("signrawtransactionwithwallet", "001122", `[{"txid":"123","vout":1,"scriptPubKey":"00","redeemScript":"01"}]`)
			},
			staticCmd: func() interface{} {
				txInputs := []pinjson.RawTxWitnessInput{
					{
						Txid:         "123",
						Vout:         1,
						ScriptPubKey: "00",
						RedeemScript: pinjson.String("01"),
					},
				}

				return pinjson.NewSignRawTransactionWithWalletCmd("001122", &txInputs, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"signrawtransactionwithwallet","params":["001122",[{"txid":"123","vout":1,"scriptPubKey":"00","redeemScript":"01"}]],"id":1}`,
			unmarshalled: &pinjson.SignRawTransactionWithWalletCmd{
				RawTx: "001122",
				Inputs: &[]pinjson.RawTxWitnessInput{
					{
						Txid:         "123",
						Vout:         1,
						ScriptPubKey: "00",
						RedeemScript: pinjson.String("01"),
					},
				},
				SigHashType: pinjson.String("ALL"),
			},
		},
		{
			name: "signrawtransactionwithwallet optional2",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("signrawtransactionwithwallet", "001122", `[]`, "ALL")
			},
			staticCmd: func() interface{} {
				txInputs := []pinjson.RawTxWitnessInput{}
				return pinjson.NewSignRawTransactionWithWalletCmd("001122", &txInputs, pinjson.String("ALL"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"signrawtransactionwithwallet","params":["001122",[],"ALL"],"id":1}`,
			unmarshalled: &pinjson.SignRawTransactionWithWalletCmd{
				RawTx:       "001122",
				Inputs:      &[]pinjson.RawTxWitnessInput{},
				SigHashType: pinjson.String("ALL"),
			},
		},
		{
			name: "walletlock",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("walletlock")
			},
			staticCmd: func() interface{} {
				return pinjson.NewWalletLockCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"walletlock","params":[],"id":1}`,
			unmarshalled: &pinjson.WalletLockCmd{},
		},
		{
			name: "walletpassphrase",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("walletpassphrase", "pass", 60)
			},
			staticCmd: func() interface{} {
				return pinjson.NewWalletPassphraseCmd("pass", 60)
			},
			marshalled: `{"jsonrpc":"1.0","method":"walletpassphrase","params":["pass",60],"id":1}`,
			unmarshalled: &pinjson.WalletPassphraseCmd{
				Passphrase: "pass",
				Timeout:    60,
			},
		},
		{
			name: "walletpassphrasechange",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd("walletpassphrasechange", "old", "new")
			},
			staticCmd: func() interface{} {
				return pinjson.NewWalletPassphraseChangeCmd("old", "new")
			},
			marshalled: `{"jsonrpc":"1.0","method":"walletpassphrasechange","params":["old","new"],"id":1}`,
			unmarshalled: &pinjson.WalletPassphraseChangeCmd{
				OldPassphrase: "old",
				NewPassphrase: "new",
			},
		},
		{
			name: "importmulti with descriptor + options",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd(
					"importmulti",
					// Cannot use a native string, due to special types like timestamp.
					[]pinjson.ImportMultiRequest{
						{Descriptor: pinjson.String("123"), Timestamp: pinjson.TimestampOrNow{Value: 0}},
					},
					`{"rescan": true}`,
				)
			},
			staticCmd: func() interface{} {
				requests := []pinjson.ImportMultiRequest{
					{Descriptor: pinjson.String("123"), Timestamp: pinjson.TimestampOrNow{Value: 0}},
				}
				options := pinjson.ImportMultiOptions{Rescan: true}
				return pinjson.NewImportMultiCmd(requests, &options)
			},
			marshalled: `{"jsonrpc":"1.0","method":"importmulti","params":[[{"desc":"123","timestamp":0}],{"rescan":true}],"id":1}`,
			unmarshalled: &pinjson.ImportMultiCmd{
				Requests: []pinjson.ImportMultiRequest{
					{
						Descriptor: pinjson.String("123"),
						Timestamp:  pinjson.TimestampOrNow{Value: 0},
					},
				},
				Options: &pinjson.ImportMultiOptions{Rescan: true},
			},
		},
		{
			name: "importmulti with descriptor + no options",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd(
					"importmulti",
					// Cannot use a native string, due to special types like timestamp.
					[]pinjson.ImportMultiRequest{
						{
							Descriptor: pinjson.String("123"),
							Timestamp:  pinjson.TimestampOrNow{Value: 0},
							WatchOnly:  pinjson.Bool(false),
							Internal:   pinjson.Bool(true),
							Label:      pinjson.String("aaa"),
							KeyPool:    pinjson.Bool(false),
						},
					},
				)
			},
			staticCmd: func() interface{} {
				requests := []pinjson.ImportMultiRequest{
					{
						Descriptor: pinjson.String("123"),
						Timestamp:  pinjson.TimestampOrNow{Value: 0},
						WatchOnly:  pinjson.Bool(false),
						Internal:   pinjson.Bool(true),
						Label:      pinjson.String("aaa"),
						KeyPool:    pinjson.Bool(false),
					},
				}
				return pinjson.NewImportMultiCmd(requests, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"importmulti","params":[[{"desc":"123","timestamp":0,"internal":true,"watchonly":false,"label":"aaa","keypool":false}]],"id":1}`,
			unmarshalled: &pinjson.ImportMultiCmd{
				Requests: []pinjson.ImportMultiRequest{
					{
						Descriptor: pinjson.String("123"),
						Timestamp:  pinjson.TimestampOrNow{Value: 0},
						WatchOnly:  pinjson.Bool(false),
						Internal:   pinjson.Bool(true),
						Label:      pinjson.String("aaa"),
						KeyPool:    pinjson.Bool(false),
					},
				},
			},
		},
		{
			name: "importmulti with descriptor + string timestamp",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd(
					"importmulti",
					// Cannot use a native string, due to special types like timestamp.
					[]pinjson.ImportMultiRequest{
						{
							Descriptor: pinjson.String("123"),
							Timestamp:  pinjson.TimestampOrNow{Value: "now"},
						},
					},
				)
			},
			staticCmd: func() interface{} {
				requests := []pinjson.ImportMultiRequest{
					{Descriptor: pinjson.String("123"), Timestamp: pinjson.TimestampOrNow{Value: "now"}},
				}
				return pinjson.NewImportMultiCmd(requests, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"importmulti","params":[[{"desc":"123","timestamp":"now"}]],"id":1}`,
			unmarshalled: &pinjson.ImportMultiCmd{
				Requests: []pinjson.ImportMultiRequest{
					{Descriptor: pinjson.String("123"), Timestamp: pinjson.TimestampOrNow{Value: "now"}},
				},
			},
		},
		{
			name: "importmulti with scriptPubKey script",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd(
					"importmulti",
					// Cannot use a native string, due to special types like timestamp and scriptPubKey
					[]pinjson.ImportMultiRequest{
						{
							ScriptPubKey: &pinjson.ScriptPubKey{Value: "script"},
							RedeemScript: pinjson.String("123"),
							Timestamp:    pinjson.TimestampOrNow{Value: 0},
							PubKeys:      &[]string{"aaa"},
						},
					},
				)
			},
			staticCmd: func() interface{} {
				requests := []pinjson.ImportMultiRequest{
					{
						ScriptPubKey: &pinjson.ScriptPubKey{Value: "script"},
						RedeemScript: pinjson.String("123"),
						Timestamp:    pinjson.TimestampOrNow{Value: 0},
						PubKeys:      &[]string{"aaa"},
					},
				}
				return pinjson.NewImportMultiCmd(requests, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"importmulti","params":[[{"scriptPubKey":"script","timestamp":0,"redeemscript":"123","pubkeys":["aaa"]}]],"id":1}`,
			unmarshalled: &pinjson.ImportMultiCmd{
				Requests: []pinjson.ImportMultiRequest{
					{
						ScriptPubKey: &pinjson.ScriptPubKey{Value: "script"},
						RedeemScript: pinjson.String("123"),
						Timestamp:    pinjson.TimestampOrNow{Value: 0},
						PubKeys:      &[]string{"aaa"},
					},
				},
			},
		},
		{
			name: "importmulti with scriptPubKey address",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd(
					"importmulti",
					// Cannot use a native string, due to special types like timestamp and scriptPubKey
					[]pinjson.ImportMultiRequest{
						{
							ScriptPubKey:  &pinjson.ScriptPubKey{Value: pinjson.ScriptPubKeyAddress{Address: "addr"}},
							WitnessScript: pinjson.String("123"),
							Timestamp:     pinjson.TimestampOrNow{Value: 0},
							Keys:          &[]string{"aaa"},
						},
					},
				)
			},
			staticCmd: func() interface{} {
				requests := []pinjson.ImportMultiRequest{
					{
						ScriptPubKey:  &pinjson.ScriptPubKey{Value: pinjson.ScriptPubKeyAddress{Address: "addr"}},
						WitnessScript: pinjson.String("123"),
						Timestamp:     pinjson.TimestampOrNow{Value: 0},
						Keys:          &[]string{"aaa"},
					},
				}
				return pinjson.NewImportMultiCmd(requests, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"importmulti","params":[[{"scriptPubKey":{"address":"addr"},"timestamp":0,"witnessscript":"123","keys":["aaa"]}]],"id":1}`,
			unmarshalled: &pinjson.ImportMultiCmd{
				Requests: []pinjson.ImportMultiRequest{
					{
						ScriptPubKey:  &pinjson.ScriptPubKey{Value: pinjson.ScriptPubKeyAddress{Address: "addr"}},
						WitnessScript: pinjson.String("123"),
						Timestamp:     pinjson.TimestampOrNow{Value: 0},
						Keys:          &[]string{"aaa"},
					},
				},
			},
		},
		{
			name: "importmulti with ranged (int) descriptor",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd(
					"importmulti",
					// Cannot use a native string, due to special types like timestamp.
					[]pinjson.ImportMultiRequest{
						{
							Descriptor: pinjson.String("123"),
							Timestamp:  pinjson.TimestampOrNow{Value: 0},
							Range:      &pinjson.DescriptorRange{Value: 7},
						},
					},
				)
			},
			staticCmd: func() interface{} {
				requests := []pinjson.ImportMultiRequest{
					{
						Descriptor: pinjson.String("123"),
						Timestamp:  pinjson.TimestampOrNow{Value: 0},
						Range:      &pinjson.DescriptorRange{Value: 7},
					},
				}
				return pinjson.NewImportMultiCmd(requests, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"importmulti","params":[[{"desc":"123","timestamp":0,"range":7}]],"id":1}`,
			unmarshalled: &pinjson.ImportMultiCmd{
				Requests: []pinjson.ImportMultiRequest{
					{
						Descriptor: pinjson.String("123"),
						Timestamp:  pinjson.TimestampOrNow{Value: 0},
						Range:      &pinjson.DescriptorRange{Value: 7},
					},
				},
			},
		},
		{
			name: "importmulti with ranged (slice) descriptor",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd(
					"importmulti",
					// Cannot use a native string, due to special types like timestamp.
					[]pinjson.ImportMultiRequest{
						{
							Descriptor: pinjson.String("123"),
							Timestamp:  pinjson.TimestampOrNow{Value: 0},
							Range:      &pinjson.DescriptorRange{Value: []int{1, 7}},
						},
					},
				)
			},
			staticCmd: func() interface{} {
				requests := []pinjson.ImportMultiRequest{
					{
						Descriptor: pinjson.String("123"),
						Timestamp:  pinjson.TimestampOrNow{Value: 0},
						Range:      &pinjson.DescriptorRange{Value: []int{1, 7}},
					},
				}
				return pinjson.NewImportMultiCmd(requests, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"importmulti","params":[[{"desc":"123","timestamp":0,"range":[1,7]}]],"id":1}`,
			unmarshalled: &pinjson.ImportMultiCmd{
				Requests: []pinjson.ImportMultiRequest{
					{
						Descriptor: pinjson.String("123"),
						Timestamp:  pinjson.TimestampOrNow{Value: 0},
						Range:      &pinjson.DescriptorRange{Value: []int{1, 7}},
					},
				},
			},
		},
		{
			name: "walletcreatefundedpsbt",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd(
					"walletcreatefundedpsbt",
					[]pinjson.PsbtInput{
						{
							Txid:     "1234",
							Vout:     0,
							Sequence: 0,
						},
					},
					[]pinjson.PsbtOutput{
						pinjson.NewPsbtOutput("1234", pinutil.Amount(1234)),
						pinjson.NewPsbtDataOutput([]byte{1, 2, 3, 4}),
					},
					pinjson.Uint32(1),
					pinjson.WalletCreateFundedPsbtOpts{},
					pinjson.Bool(true),
				)
			},
			staticCmd: func() interface{} {
				return pinjson.NewWalletCreateFundedPsbtCmd(
					[]pinjson.PsbtInput{
						{
							Txid:     "1234",
							Vout:     0,
							Sequence: 0,
						},
					},
					[]pinjson.PsbtOutput{
						pinjson.NewPsbtOutput("1234", pinutil.Amount(1234)),
						pinjson.NewPsbtDataOutput([]byte{1, 2, 3, 4}),
					},
					pinjson.Uint32(1),
					&pinjson.WalletCreateFundedPsbtOpts{},
					pinjson.Bool(true),
				)
			},
			marshalled: `{"jsonrpc":"1.0","method":"walletcreatefundedpsbt","params":[[{"txid":"1234","vout":0,"sequence":0}],[{"1234":0.00001234},{"data":"01020304"}],1,{},true],"id":1}`,
			unmarshalled: &pinjson.WalletCreateFundedPsbtCmd{
				Inputs: []pinjson.PsbtInput{
					{
						Txid:     "1234",
						Vout:     0,
						Sequence: 0,
					},
				},
				Outputs: []pinjson.PsbtOutput{
					pinjson.NewPsbtOutput("1234", pinutil.Amount(1234)),
					pinjson.NewPsbtDataOutput([]byte{1, 2, 3, 4}),
				},
				Locktime:    pinjson.Uint32(1),
				Options:     &pinjson.WalletCreateFundedPsbtOpts{},
				Bip32Derivs: pinjson.Bool(true),
			},
		},
		{
			name: "walletprocesspsbt",
			newCmd: func() (interface{}, error) {
				return pinjson.NewCmd(
					"walletprocesspsbt", "1234", pinjson.Bool(true), pinjson.String("ALL"), pinjson.Bool(true))
			},
			staticCmd: func() interface{} {
				return pinjson.NewWalletProcessPsbtCmd(
					"1234", pinjson.Bool(true), pinjson.String("ALL"), pinjson.Bool(true))
			},
			marshalled: `{"jsonrpc":"1.0","method":"walletprocesspsbt","params":["1234",true,"ALL",true],"id":1}`,
			unmarshalled: &pinjson.WalletProcessPsbtCmd{
				Psbt:        "1234",
				Sign:        pinjson.Bool(true),
				SighashType: pinjson.String("ALL"),
				Bip32Derivs: pinjson.Bool(true),
			},
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
