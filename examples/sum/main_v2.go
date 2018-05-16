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
	//testHash    = common.StringToHash("duanbing")
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
	binFileName := "./sum_sol_sum.bin"
	abiFileName := "./sum_sol_sum.abi"
	data := loadBin(binFileName)
	msg := ec.NewMessage(testAddress, &toAddress, nonce, amount, gasLimit, big.NewInt(1), data, false)
	header := types.Header{
		// ParentHash: common.Hash{},
		// UncleHash:  common.Hash{},
		Coinbase: coinbase,
		//	Root:        common.Hash{},
		//	TxHash:      common.Hash{},
		//	ReceiptHash: common.Hash{},
		//	Bloom:      types.BytesToBloom([]byte("duanbing")),
		Difficulty: big.NewInt(1),
		Number:     big.NewInt(1),
		GasLimit:   gasLimit,
		GasUsed:    big.NewInt(1),
		Time:       big.NewInt(time.Now().Unix()),
		Extra:      nil,
		//MixDigest:  testHash,
		//Nonce:      types.EncodeNonce(1),
	}
	cc := ChainContext{}
	ctx := ec.NewEVMContext(msg, &header, cc, &testAddress)
	mdb, err := ethdb.NewMemDatabase()
	must(err)
	db := state.NewDatabase(mdb)
	statedb, err := state.New(common.Hash{}, db)
	//set balance
	statedb.GetOrNewStateObject(testAddress)
	statedb.GetOrNewStateObject(toAddress)
	statedb.AddBalance(testAddress, big.NewInt(1e18))
	testBalance := statedb.GetBalance(testAddress)
	fmt.Println("init testBalance =", testBalance)
	must(err)

	//	config := params.TestnetChainConfig
	config := params.AllProtocolChanges
	logConfig := vm.LogConfig{}
	structLogger := vm.NewStructLogger(&logConfig)
	vmConfig := vm.Config{Debug: true, Tracer: structLogger, DisableGasMetering: false /*, JumpTable: vm.NewByzantiumInstructionSet()*/}

	evm := vm.NewEVM(ctx, statedb, config, vmConfig)
	contractRef := vm.AccountRef(testAddress)
	contractCode, _, gasLeftover, vmerr := evm.Create(contractRef, data, statedb.GetBalance(testAddress).Uint64(), big.NewInt(0))
	must(vmerr)
	statedb.SetBalance(testAddress, big.NewInt(0).SetUint64(gasLeftover))
	testBalance = statedb.GetBalance(testAddress)
	fmt.Println("after create contract, testBalance =", testBalance)
	// set input ,  formatted accocding to https://solidity.readthedocs.io/en/develop/abi-spec.html
	//find methods := "multiply(uint)"
	abiObj := loadAbi(abiFileName)
	method := abiObj.Methods["multiply"]
	//make params := "0xa"
	pm := abi.U256(big.NewInt(10))
	//concat method and params
	//inputstr := hexutil.Encode(method.Id()) + pm[2:]
	input := append(method.Id(), pm...)
	//fmt.Println(hexutil.Encode(input))
	fmt.Println("begin to exec contract")
	statedb.SetCode(testAddress, contractCode)
	outputs, gasLeftover, vmerr := evm.Call(contractRef, testAddress, input, statedb.GetBalance(testAddress).Uint64(), big.NewInt(0))
	must(vmerr)

	statedb.SetBalance(testAddress, big.NewInt(0).SetUint64(gasLeftover))
	testBalance = statedb.GetBalance(testAddress)
	fmt.Println("after call contract, testBalance =", testBalance)
	for _, op := range method.Outputs {
		switch op.Type.String() {
		case "uint256":
			fmt.Printf("Output name=%s, value=%d\n", op.Name, big.NewInt(0).SetBytes(outputs))

		default:
			fmt.Println(op.Name, op.Type.String())
		}
	}

}

type ChainContext struct{}

func (cc ChainContext) GetHeader(hash common.Hash, number uint64) *types.Header {
	fmt.Println("(cc ChainContext) GetHeader(hash common.Hash, number uint64)")
	return nil
	//return &header
}
