package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	token "go-ethereum-example/gen"
	"math/big"
	"os"

	_ "github.com/joho/godotenv/autoload"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Connect to Ethereum client with RPC endpoint
	client, err := ethclient.DialContext(ctx, os.Getenv("RPC_ENDPOINT"))
	handleError(err)

	defer client.Close()

	fmt.Println("Successfully connected to Ethereum client")

	// Parse wallet private key
	privateKey := mustParsePrivateKey()

	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	fmt.Printf("Deploying contract from address %s\n", address.Hex())

	// Get nonce, gas price and chain ID from the Ethereum client
	nonce, err := client.PendingNonceAt(ctx, address)
	handleError(err)

	gasPrice, err := client.SuggestGasPrice(ctx)
	handleError(err)

	fmt.Printf("Suggested gas price: %s\n", gasPrice)

	chainID, err := client.NetworkID(ctx)
	handleError(err)

	fmt.Printf("Chain ID: %d\n", chainID)

	// Create an signer with the private key, chain ID and nonce
	signer, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	handleError(err)

	signer.GasPrice = gasPrice
	signer.GasLimit = 3000000
	signer.Nonce = big.NewInt(int64(nonce))

	// Deploy the contract with initial supply of 1,000,000 tokens
	initialSupply := big.NewInt(0).Mul(big.NewInt(1000000), big.NewInt(0).Exp(big.NewInt(10), big.NewInt(18), nil))

	_, tx, _, err := token.DeployToken(signer, client, initialSupply)
	handleError(err)

	fmt.Printf("Transaction hash: %s\n", tx.Hash().Hex())

	// Wait for the deployment to be mined
	contractAddress, err := bind.WaitDeployed(ctx, client, tx)
	handleError(err)

	fmt.Printf("Contract deployed! Contract address: %s\n", contractAddress.Hex())
}

func mustParsePrivateKey() *ecdsa.PrivateKey {
	rawPrivateKey := os.Getenv("PRIVATE_KEY")

	// Parse the private key
	privateKey, err := crypto.HexToECDSA(rawPrivateKey)
	handleError(err)

	return privateKey
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
