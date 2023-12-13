package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	token "go-ethereum-example/gen"
	"math/big"
	"os"

	_ "github.com/joho/godotenv/autoload"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
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

	// Change these addresses to match your contract!
	contractAddress := common.HexToAddress("0x5FbDB2315678afecb367f032d93F642f64180aa3")
	toAddress := common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")

	// Create an instance of the contract, specifying its address
	tokenInstance, err := token.NewToken(contractAddress, client)
	handleError(err)

	// Parse wallet private key
	privateKey := mustParsePrivateKey()
	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	// Get nonce, gas price and chain ID
	nonce, err := client.PendingNonceAt(context.Background(), address)
	handleError(err)

	gasPrice, err := client.SuggestGasPrice(context.Background())
	handleError(err)

	fmt.Printf("Suggested gas price: %s\n", gasPrice)

	chainID, err := client.NetworkID(context.Background())
	handleError(err)

	fmt.Printf("Chain ID: %d\n", chainID)

	// Create an transactor with the private key, chain ID and nonce
	signer, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	handleError(err)

	signer.GasPrice = gasPrice
	signer.GasLimit = 3000000
	signer.Nonce = big.NewInt(int64(nonce))

	// Call transfer method (state-changing)
	tx, err := tokenInstance.Transfer(signer, toAddress, big.NewInt(1000000))
	handleError(err)

	fmt.Printf("Transaction hash: %s\n", tx.Hash().Hex())

	// Wait for the transaction to be mined
	receipt, err := bind.WaitMined(context.Background(), client, tx)
	handleError(err)

	fmt.Printf("Transaction receipt status %d\n", receipt.Status)

	// If the transaction was reverted by the EVM, we can see the reason
	if receipt.Status == 0 {
		msg := ethereum.CallMsg{
			From:     address,
			To:       tx.To(),
			Gas:      tx.Gas(),
			GasPrice: tx.GasPrice(),
			Value:    tx.Value(),
			Data:     tx.Data(),
		}

		_, err = client.CallContract(ctx, msg, nil)
		fmt.Printf("Transaction reverted: %v\n", err)
		return
	}

	// Extract the transfer event from the receipt
	var transferred *token.TokenTransfer

	for _, log := range receipt.Logs {
		transferred, err = tokenInstance.ParseTransfer(*log)
		if err == nil {
			break
		}
	}

	if transferred != nil {
		fmt.Printf("Transferred %d tokens from %s to %s\n", transferred.Value, transferred.From.Hex(), transferred.To.Hex())
	}

	// Call the contract method (read-only)
	toBalance, err := tokenInstance.BalanceOf(&bind.CallOpts{Context: ctx}, toAddress)
	handleError(err)

	fmt.Printf("To balance: %d\n", toBalance)
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
