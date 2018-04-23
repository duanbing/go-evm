### Eth-VM
This is a single golang EVM.

### Example

1. Compile contract
```
cd  $GOPATH/src/baidu.com/evm/example/event
solcjs --abi --bin coin.sol
```
Here we get code in xxx.bin, and abi in xx.abi

2. Run

```
go run mainXX.go 
```

### Usage

0. Prepare

* build sol,  get xxx.bin and xxx.abi 
* make an evm instance, the most important struct is StateDB. 

```

msg := ec.NewMessage(testAddress, &toAddress, nonce, amount, gasLimit, big.NewInt(1), data, false)
header := types.Header{
    // ParentHash: common.Hash{},
    // UncleHash:  common.Hash{},
    Coinbase: coinbase,
    //  Root:        common.Hash{},
    //  TxHash:      common.Hash{},
    //  ReceiptHash: common.Hash{},
    //  Bloom:      types.BytesToBloom([]byte("duanbing")),
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

//  config := params.TestnetChainConfig
config := params.AllProtocolChanges
logConfig := vm.LogConfig{}
structLogger := vm.NewStructLogger(&logConfig)
vmConfig := vm.Config{Debug: true, Tracer: structLogger, DisableGasMetering: false /*, JumpTable: vm.NewByzantiumInstructionSet()*/}

evm := vm.NewEVM(ctx, statedb, config, vmConfig)
```

1. Executing a contract

* create an contract, get the contract code

```
contractRef := vm.AccountRef(testAddress)
contractCode, _, gasLeftover, vmerr := evm.Create(contractRef, data, statedb.GetBalance(testAddress).Uint64(), big.NewInt(0))
```

* encode the input ,  refer to https://solidity.readthedocs.io/en/develop/abi-spec.html#argument-encoding
```
method = abiObj.Methods["mint"]
input = append(method.Id(), sender...)
pm = abi.U256(big.NewInt(1000000))
input = append(input, pm...)
```

* execute the evm.Call

```
outputs, gasLeftover, vmerr = evm.Call(senderAcc, testAddress, input, statedb.GetBalance(testAddress).Uint64(), big.NewInt(0))
```

2. Get Logs

* all logs is stored in statdb, you can call GetLogs or Logs to get all the logs

```
logs := statedb.Logs() // logi instruction. logi store i+1 fields, the first one is stored in log.Data, and the last i field stores the parameter of logi 
```

