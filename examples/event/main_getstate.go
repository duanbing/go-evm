/***************************************************************************
 *
 * Copyright (c) 2017 Baidu.com, Inc. All Rights Reserved
 * @author duanbing(duanbing@baidu.com)
 *
 **************************************************************************/

/**
 * @filename main.go
 * @desc
 * @create time 2018-04-19 15:49:26
**/
package main

import (
	"fmt"
	ec "github.com/duanbing/go-evm/core"
	"github.com/duanbing/go-evm/state"
	"github.com/duanbing/go-evm/types"
	"github.com/duanbing/go-evm/vm"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
	"io/ioutil"
	"math/big"
	"os"
	"time"
)

var (
	testHash    = common.StringToHash("duanbing")
	testAddress = common.StringToAddress("duanbing")
	toAddress   = common.StringToAddress("andone")
	amount      = big.NewInt(1)
	nonce       = uint64(0)
	gasLimit    = big.NewInt(100000)
	coinbase    = common.HexToAddress("0x0000000000000000000000000000000000000000")
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}
func loadBin(filename string) []byte {
	code, err := ioutil.ReadFile(filename)
	must(err)
	return hexutil.MustDecode("0x" + string(code))
}
func loadAbi(filename string) abi.ABI {
	abiFile, err := os.Open(filename)
	must(err)
	defer abiFile.Close()
	abiObj, err := abi.JSON(abiFile)
	must(err)
	return abiObj
}

func main() {
	binFileName := "./coin_sol_Coin.bin"
	abiFileName := "./coin_sol_Coin.abi"
	data := loadBin(binFileName)
	msg := ec.NewMessage(testAddress, &toAddress, nonce, amount, gasLimit, big.NewInt(1), data, false)
	cc := ChainContext{}
	ctx := ec.NewEVMContext(msg, cc.GetHeader(testHash, 0), cc, &testAddress)
	mdb, err := ethdb.NewLDBDatabase("/tmp/a.txt", 100, 100)
	must(err)
	db := state.NewDatabase(mdb)
	statedb, err := state.New(common.HexToHash("0x2b92c9cabb651051e874e27ef8ec9b7a168249e5b7881d5902177010527b8ad2"), db)
	must(err)
	//set balance
	statedb.GetOrNewStateObject(testAddress)
	statedb.GetOrNewStateObject(toAddress)
	testBalance := statedb.GetBalance(testAddress)
	fmt.Println("init testBalance =", testBalance)

	//	config := params.TestnetChainConfig
	config := params.MainnetChainConfig
	logConfig := vm.LogConfig{}
	structLogger := vm.NewStructLogger(&logConfig)
	vmConfig := vm.Config{Debug: true, Tracer: structLogger /*, JumpTable: vm.NewByzantiumInstructionSet()*/}

	evm := vm.NewEVM(ctx, statedb, config, vmConfig)
	contractRef := vm.AccountRef(testAddress)
	contractCode, contractAddr, gasLeftover, vmerr := evm.Create(contractRef, data, statedb.GetBalance(testAddress).Uint64(), big.NewInt(0))
	must(vmerr)
	statedb.SetBalance(testAddress, big.NewInt(0).SetUint64(gasLeftover))
	testBalance = statedb.GetBalance(testAddress)
	fmt.Println("after create contract, testBalance =", testBalance)
	abiObj := loadAbi(abiFileName)
	// get minter
	input, err := abiObj.Pack("minter")
	must(err)
	evm.StateDB.SetCode(testAddress, contractCode)
	outputs, gasLeftover, vmerr := evm.Call(contractRef, testAddress, input, statedb.GetBalance(testAddress).Uint64(), big.NewInt(0))
	must(vmerr)
	sender := outputs
	senderAcc := vm.AccountRef(common.BytesToAddress(sender))

	// get balance
	input, err = abiObj.Pack("balances", toAddress)
	must(err)
	outputs, gasLeftover, vmerr = evm.Call(contractRef, testAddress, input, statedb.GetBalance(testAddress).Uint64(), big.NewInt(0))
	must(vmerr)
	Print(outputs, "balances")

	// get balance
	input, err = abiObj.Pack("balances", common.BytesToAddress(sender))
	must(err)
	outputs, gasLeftover, vmerr = evm.Call(contractRef, testAddress, input, statedb.GetBalance(testAddress).Uint64(), big.NewInt(0))
	must(vmerr)
	Print(outputs, "balances")
	//send 2
	input, err = abiObj.Pack("send", toAddress, big.NewInt(19))
	must(err)
	outputs, gasLeftover, vmerr = evm.Call(senderAcc, testAddress, input, statedb.GetBalance(testAddress).Uint64(), big.NewInt(0))
	must(vmerr)

	// get balance
	input, err = abiObj.Pack("balances", toAddress)
	must(err)
	outputs, gasLeftover, vmerr = evm.Call(contractRef, testAddress, input, statedb.GetBalance(testAddress).Uint64(), big.NewInt(0))
	must(vmerr)
	Print(outputs, "balances")

	// get balance
	input, err = abiObj.Pack("balances", common.BytesToAddress(sender))
	must(err)
	outputs, gasLeftover, vmerr = evm.Call(contractRef, testAddress, input, statedb.GetBalance(testAddress).Uint64(), big.NewInt(0))
	must(vmerr)
	Print(outputs, "balances")

	// get event
	logs := statedb.Logs()

	for _, log := range logs {
		fmt.Printf("%#v\n", log)
		for _, topic := range log.Topics {
			fmt.Printf("topic: %#v\n", topic)
		}
		fmt.Printf("data: %#v\n", log.Data)
	}

	idx := 1
	getstateFunc := func(key, value common.Hash) bool {
		fmt.Printf("------------- idx=%d, key=%v,value=%v\n", idx, key, value)
		idx += 1
		return true
	}
	statedb.ForEachStorage(contractAddr, getstateFunc)
	root, err := statedb.Commit(true)
	must(err)
	fmt.Println("Root hash", root.Hex())
}

func Print(outputs []byte, name string) {
	fmt.Printf("method=%s, output=%#v\n", name, outputs)
}

type ChainContext struct{}

func (cc ChainContext) GetHeader(hash common.Hash, number uint64) *types.Header {

	return &types.Header{
		// ParentHash: common.Hash{},
		// UncleHash:  common.Hash{},
		Coinbase: coinbase,
		//	Root:        common.Hash{},
		//	TxHash:      common.Hash{},
		//	ReceiptHash: common.Hash{},
		//	Bloom:      types.BytesToBloom([]byte("duanbing")),
		Difficulty: big.NewInt(1),
		Number:     big.NewInt(1),
		GasLimit:   1000000,
		GasUsed:    1,
		Time:       big.NewInt(time.Now().Unix()),
		Extra:      nil,
		//MixDigest:  testHash,
		//Nonce:      types.EncodeNonce(1),
	}
}
