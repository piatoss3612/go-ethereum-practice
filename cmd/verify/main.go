package main

import (
	"bytes"
	"context"
	"net/http"
	"net/url"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

var etherscanURL = "https://api.etherscan.io/api"

func main() {
	client := http.DefaultClient

	data := url.Values{
		"apiKey": []string{os.Getenv("ETHERSCAN_API_KEY")},
		"module": []string{"contract"},
		"action": []string{"verifysourcecode"},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, etherscanURL, bytes.NewBufferString(data.Encode()))
	handleErr(err)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	handleErr(err)

	defer resp.Body.Close()
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
