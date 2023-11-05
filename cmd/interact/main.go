package main

import (
	"context"
	"crypto/ecdsa"
	storage "go-ethereum-example/gen"
	"log"
	"math/big"
	"os"

	_ "github.com/joho/godotenv/autoload"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// Connect to Ethereum client with RPC endpoint
	client, err := ethclient.Dial(os.Getenv("RPC_ENDPOINT"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	log.Println("Successfully connected to Ethereum client")

	contractAddress := common.HexToAddress("bE7F4aC08B6B58fD4d7085a9AE1811EF1eae1EB4")

	// Create an instance of the contract, specifying its address
	storageInstance, err := storage.NewStorage(contractAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	// Parse wallet private key
	privateKey := parsePrivateKey()
	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	// Get nonce, gas price and chain ID
	nonce, err := client.PendingNonceAt(context.Background(), address)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Suggested gas price: %s", gasPrice)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Chain ID: %d", chainID)

	// Create an transactor with the private key, chain ID and nonce
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatal(err)
	}

	auth.GasPrice = gasPrice
	auth.GasLimit = 3000000
	auth.Nonce = big.NewInt(int64(nonce))

	// Call the contract method (state-changing)
	tx, err := storageInstance.Store(auth, big.NewInt(20))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Transaction hash: %s", tx.Hash().Hex())

	// Wait for the transaction to be mined
	receipt, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Transaction receipt: %s", receipt.TxHash.Hex())

	// Call the contract method (read-only)
	retrieved, err := storageInstance.Retrieve(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Retrieved value: %d", retrieved)
}

func parsePrivateKey() *ecdsa.PrivateKey {
	rawPrivateKey := os.Getenv("PRIVATE_KEY")

	// Parse the private key
	privateKey, err := crypto.HexToECDSA(rawPrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	return privateKey
}
