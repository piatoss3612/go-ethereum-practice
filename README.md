# go-ethereum example

## Generate abi 

## Generate abi and binary from solidity file

```bash
$ solc --bin --abi contracts/Storage.sol -o build
Compiler run successful. Artifact(s) can be found in directory "build".
```

## Generate go file from abi

```bash
$ mkdir gen
$ abigen --bin=build/Storage.bin --abi=build/Storage.abi --pkg=storage --out=gen/storage.go
```

## Update go.mod

```bash
$ go mod tidy
```