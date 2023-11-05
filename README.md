# go-ethereum example

## Generate go code from solidity file

### Generate abi and binary from solidity file

```bash
$ solc --bin --abi contracts/Storage.sol -o build
Compiler run successful. Artifact(s) can be found in directory "build".
```

### Generate go file from abi

```bash
$ mkdir gen
$ abigen --bin=build/Storage.bin --abi=build/Storage.abi --pkg=storage --out=gen/storage.go
```

### Update go.mod

```bash
$ go mod tidy
```

## Deploy contract to local testnet

### Create .env file in root directory

```
PRIVATE_KEY=<YOUR_PRIVATE_KEY>
RPC_ENDPOINT=<YOUR_RPC_ENDPOINT>
```

### Deploy contract

```bash
$ go run ./cmd/deploy/
2023/11/05 20:15:33 Successfully connected to Ethereum client
2023/11/05 20:15:33 Deploying contract from address 0x965B0E63e00E7805569ee3B428Cf96330DFc57EF
2023/11/05 20:15:34 Suggested gas price: 2039359022
2023/11/05 20:15:34 Chain ID: 80001
2023/11/05 20:15:34 Contract deployed! Transaction hash: 0x5b2cba1c0022c76809edb01e9555d40a597fe12d1d0c6d4f0bbd280a4e859a6b
2023/11/05 20:15:34 Contract address: 0xbE7F4aC08B6B58fD4d7085a9AE1811EF1eae1EB4
```