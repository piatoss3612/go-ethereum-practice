package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	_ "github.com/joho/godotenv/autoload"
)

var etherscanURL = "https://api-testnet.polygonscan.com/api"

func main() {
	// read input file
	f, err := os.Open("verify/MyToken_input.json")
	handleErr(err)

	defer f.Close()

	sourceCodeBytes, err := io.ReadAll(f)
	handleErr(err)

	client := http.DefaultClient

	// generate abi-encoded constructor arguments
	initialSupply := big.NewInt(0).Mul(big.NewInt(1000000), big.NewInt(0).Exp(big.NewInt(10), big.NewInt(18), nil))

	uint256Ty, err := abi.NewType("uint256", "uint256", nil)
	handleErr(err)

	args := abi.Arguments{
		{
			Type: uint256Ty,
		},
	}

	encodedArgsBytes, err := args.Pack(initialSupply)
	handleErr(err)

	encodedArgsHex := fmt.Sprintf("%x", encodedArgsBytes)

	// set url encode data
	data := url.Values{
		"apiKey":                []string{os.Getenv("ETHERSCAN_API_KEY")},
		"module":                []string{"contract"},
		"action":                []string{"verifysourcecode"},
		"sourceCode":            []string{string(sourceCodeBytes)},
		"contractaddress":       []string{"0x7Fc3c9ae336291EC87296bb10D4B03f7d23357e4"},
		"codeformat":            []string{"solidity-standard-json-input"},
		"contractname":          []string{"contracts/MyToken.sol:MyToken"}, // contractfile.sol:contractname format
		"compilerversion":       []string{"v0.8.22+commit.4fc1097e"},
		"optimizationUsed":      []string{"1"},
		"constructorArguements": []string{encodedArgsHex}, // abi-encoded constructor arguments
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// create request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, etherscanURL, bytes.NewBufferString(data.Encode()))
	handleErr(err)

	// set request headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// send request
	resp, err := client.Do(req)
	handleErr(err)

	defer resp.Body.Close()

	// read response body
	body, err := io.ReadAll(resp.Body)
	handleErr(err)

	// success: 200 OK

	// print response body
	fmt.Println(string(body))
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
