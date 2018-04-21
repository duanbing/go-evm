### Eth-VM
This is a single golang EVM.

### Usage for example

##### Step 1: Compile contract
```
cd  $GOPATH/src/baidu.com/evm/example/event
solcjs --abi --bin coin.sol
```
Here we get code in xxx.bin, and abi in xx.abi

##### Run

```
go run mainXX.go 
```


