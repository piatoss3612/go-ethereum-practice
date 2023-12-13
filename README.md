# go-ethereum practice: deploy ERC20 token contract

## Table of Contents

- [1. Generate Go code from solidity file](#1-generate-go-code-from-solidity-file)
    - [Generate abi and binary from solidity file](#generate-abi-and-binary-from-solidity-file)
    - [Install solc](#install-solc)
    - [Install geth](#install-geth)
    - [Build abigen](#build-abigen)
    - [Install node modules](#install-node-modules)
    - [Generate go file from abi](#generate-go-file-from-abi)
    - [Update go module](#update-go-module)
- [2. Run local test network](#2-run-local-test-network)
- [3. Deploy contract to local test network](#3-deploy-contract-to-local-test-network)
    - [Create .env file in root directory](#create-env-file-in-root-directory)
    - [Deploy contract](#deploy-contract)
- [4. Interact with contract](#4-interact-with-contract)
    - [Set contract address and to address for transfer](#set-contract-address-and-to-address-for-transfer)
    - [Interact with contract](#interact-with-contract)
- [5. Subscribe to events](#5-subscribe-to-events)
    - [Subscribe to events](#subscribe-to-events)
    - [Trigger event](#trigger-event)
    - [Output](#output)
- [6. Verify contract](#6-verify-contract)
    - [Generate metadata from solidity file](#generate-metadata-from-solidity-file)
    - [Generate standard json input file from metadata](#generate-standard-json-input-file-from-metadata)
    - [Test standard json input file](#test-standard-json-input-file)
    - [Update .env file](#update-env-file)
    - [Check solc version](#check-solc-version)
    - [Verify contract](#verify-contract)
    - [Check on polygonscan (mumbai testnet)](#check-on-polygonscan-mumbai-testnet)

## 1. Generate Go code from solidity file

### Install solc

- [solc installation](https://docs.soliditylang.org/en/latest/installing-solidity.html)

### Install geth

- [geth installation](https://geth.ethereum.org/docs/install-and-build/installing-geth)

### Build abigen

```bash
$ cd $GOPATH/src/github.com/ethereum/go-ethereum
$ go build ./cmd/abigen
```

### Install node modules

```bash
$ npm install
```

### Generate abi and binary from solidity file

```bash
$ solc @openzeppelin/=$(pwd)/node_modules/@openzeppelin/ --optimize --abi --bin --pretty-json contracts/MyToken.sol -o build --overwrite
Compiler run successful. Artifact(s) can be found in directory "build".
```

> binary file is required to deploy contract, if you don't want to deploy contract, you can ignore --bin option
> --overwrite option is required to overwrite existing files (if needed)
> --pretty-json is optional to generate pretty json file

### Generate go file from abi

```bash
$ mkdir -p gen
$ abigen --bin=build/MyToken.bin --abi=build/MyToken.abi --pkg=token --out=gen/token.go
```

### Update go module

```bash
$ go mod tidy
```

## 2. Run local test network

- Installation of Foundry required to run anvil (https://book.getfoundry.sh/getting-started/installation)
- or you can use ganache-cli

```bash
$ anvil
```

## 3. Deploy contract to local test network

### Create .env file in root directory

```.env
PRIVATE_KEY=<FIRST_ACCOUNT_FROM_ANVIL> # Should not be 0x prefixed
RPC_ENDPOINT=http://localhost:8545
RPC_WS_ENDPOINT=ws://localhost:8545
```

### Deploy contract

```bash
$ go run ./cmd/deploy/
Successfully connected to Ethereum client
Deploying contract from address 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
Suggested gas price: 2000000000
Chain ID: 31337
Transaction hash: 0x9b3b8b92cf9370e0a51c4f1f35726385839ef4aba748c42f22d7327f00cca5ad
Contract deployed! Contract address: 0x5FbDB2315678afecb367f032d93F642f64180aa3
```

## 4. Interact with contract

### Set contract address and to address for transfer

```go
// Contract address you deployed and to address from anvil accounts
contractAddress := common.HexToAddress("0x5FbDB2315678afecb367f032d93F642f64180aa3")
toAddress := common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")
```

### Interact with contract

```bash
$ go run ./cmd/interact/
2023/12/12 13:38:36 Successfully connected to Ethereum client
Suggested gas price: 1879547559
Chain ID: 31337
Transaction hash: 0x8336dceaef677f437377c6c90caf4ba15c3851c8f386c060e386fab5f90df64c
Transaction receipt status 1
Transferred 1000000 tokens from 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266 to 0x70997970C51812dc3A010C7d01b50e0d17dc79C8
To balance: 1000000
```

## 5. Subscribe to events

### Subscribe to events

```bash
$ go run ./cmd/subscribe/
Successfully connected to Ethereum client
Successfully subscribed to Transfer events
```

### Trigger event

```bash
$ go run ./cmd/interact/
```

### Output

```bash
Transfer event received: from=0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266 to=0x70997970C51812dc3A010C7d01b50e0d17dc79C8 value=1000000
```

## 6. Verify contract

### Generate metadata from solidity file

```bash
$ solc @openzeppelin/=$(pwd)/node_modules/@openzeppelin/ --optimize --metadata --metadata-literal contracts/MyToken.sol -o build
```

### Generate standard json input file from metadata
    
```bash
$ go run ./cmd/input
```

### Test standard json input file
    
```bash
$ solc --standard-json ./verify/MyToken_input.json
```

### Update .env file

```.env
ETHERSCAN_API_KEY=<YOUR_ETHERSCAN_API_KEY>
```

### Check solc version

```bash
$ solc --version
solc, the solidity compiler commandline interface
Version: 0.8.22+commit.4fc1097e.Linux.g++
```

### Verify contract

```bash
$ go run ./cmd/verify
{"status":"1","message":"OK","result":"vykmzujkyimxbzxn5cek1iyfmv8hj1hf2xdwxw4ephyk8maeyb"}
```

### Check on polygonscan (mumbai testnet)

- https://mumbai.polygonscan.com/address/0x7fc3c9ae336291ec87296bb10d4b03f7d23357e4#code