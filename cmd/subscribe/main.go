package main

import (
	"context"
	"fmt"
	token "go-ethereum-example/gen"

	_ "github.com/joho/godotenv/autoload"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Connect to Ethereum client with RPC endpoint
	client, err := ethclient.DialContext(ctx, "ws://localhost:8545")
	handleError(err)

	defer client.Close()

	fmt.Println("Successfully connected to Ethereum client")

	// Change these addresses to match your contract!
	contractAddress := common.HexToAddress("0x5FbDB2315678afecb367f032d93F642f64180aa3")

	// Create an instance of the contract, specifying its address
	tokenInstance, err := token.NewToken(contractAddress, client)
	handleError(err)

	transferChan := make(chan *token.TokenTransfer)

	// Subscribe to Transfer events
	sub, err := tokenInstance.WatchTransfer(&bind.WatchOpts{}, transferChan, nil, nil)
	handleError(err)

	defer sub.Unsubscribe()

	fmt.Println("Successfully subscribed to Transfer events")

	for {
		select {
		case err := <-sub.Err():
			handleError(err)
		case transfer := <-transferChan:
			fmt.Printf("Transfer event received: from=%s to=%s value=%d\n", transfer.From.Hex(), transfer.To.Hex(), transfer.Value)
		}
	}
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
