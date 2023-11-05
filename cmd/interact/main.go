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
	client, err := ethclient.Dial(os.Getenv("RPC_ENDPOINT"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	log.Println("Successfully connected to Ethereum client")

	contractAddress := common.HexToAddress("bE7F4aC08B6B58fD4d7085a9AE1811EF1eae1EB4")

	storageInstance, err := storage.NewStorage(contractAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	privateKey := parsePrivateKey()
	address := crypto.PubkeyToAddress(privateKey.PublicKey)

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

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatal(err)
	}

	auth.GasPrice = gasPrice
	auth.GasLimit = 3000000
	auth.Nonce = big.NewInt(int64(nonce))

	tx, err := storageInstance.Store(auth, big.NewInt(20))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Transaction hash: %s", tx.Hash().Hex())

	receipt, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Transaction receipt: %s", receipt.TxHash.Hex())

	retrieved, err := storageInstance.Retrieve(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Retrieved value: %d", retrieved)
}

func parsePrivateKey() *ecdsa.PrivateKey {
	rawPrivateKey := os.Getenv("PRIVATE_KEY")

	privateKey, err := crypto.HexToECDSA(rawPrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	return privateKey
}
