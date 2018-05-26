### Eth-VM
This is a single golang EVM, it's based on go-ethereum release1.8
All dependencies is maintained by `godep`.

### Run an example

1. Compile contract
```
cd  $GOPATH/src/github.com/duanbing/go-evm/example/event
solcjs --abi --bin coin.sol
```
Here we get code in xxx.bin, and abi in xx.abi

2. Run

```
godep go run main.go 
```

### Usage

1. Prepare

* Build sol,  get xxx.bin and xxx.abi 
* Make an evm instance, the most important struct is StateDB. 

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

2. Executing a contract

* Create an contract, get the contract code and contractAddr,  you can get code via contractAddr in the MPT tree.

```
contractRef := vm.AccountRef(testAddress)
contractCode, contractAddr, gasLeftover, vmerr := evm.Create(contractRef, data, statedb.GetBalance(testAddress).Uint64(), big.NewInt(0))
```

* Encode the input ,  refer to https://solidity.readthedocs.io/en/develop/abi-spec.html#argument-encoding
```
input, err = abiObj.Pack("send", toAddress, big.NewInt(19))
```

* Execute the evm.Call

```
outputs, gasLeftover, vmerr = evm.Call(senderAcc, testAddress, input, statedb.GetBalance(testAddress).Uint64(), big.NewInt(0))
```

3. Get Logs

* All logs is stored in statdb, you can call GetLogs or Logs to get all the logs

```
logs := statedb.Logs() // logi instruction. logi store i+1 fields, the first one is stored in log.Data, and the last i field stores the parameter of logi 
```
### Embedding

This is the most exciting part!  if you want to combine UTXO or other ledgers tech with EVM, just implement `core/interface.go`! the MPT tree from Ethereum takes care of the smart contract, and your blockchain takes care of the transaction!
