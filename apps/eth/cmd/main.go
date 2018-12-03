package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/tokenise-eu/tokenise-apps/apps/eth/token"

	"github.com/tokenise-eu/tokenise-apps/apps/eth/listener"
	"github.com/tokenise-eu/tokenise-apps/apps/eth/migrator"
)

var mainnet = flag.Bool("mainnet", false, "Deploy the contract and listen "+
	"on the Ethereum mainnet (default)")
var ropsten = flag.Bool("ropsten", false, "Deploy the contract and listen "+
	"on the Ropsten testnet")
var rinkeby = flag.Bool("rinkeby", false, "Deploy the contract and listen "+
	"on the Rinkeby testnet")
var addr = flag.String("contract", "0x0", "Specify a contract listen to.")

var key = `{"address":"20c62701345b727ef8a4cbb61ddc9764cd9d2634","crypto":{"cipher":"aes-128-ctr","ciphertext":"8de8ecc697cf199fc9331cf311fd33ab84b1713cdf3fc1fbe37a5c881541fd21","cipherparams":{"iv":"d1c7f2550163ab9df89015e671088391"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"335c48f497bcf282e35cba90df1ce950ca0cae305d52f2b84aa12bf0d36ba62d"},"mac":"40cff15f203211262ad1fd540b64d03b852b1f8c97afb978590821415ed03526"},"id":"000c9aa5-5f68-45ae-a2e1-1af7bc955bc4","version":3}`

func main() {
	// Parse flags first
	flag.Parse()
	network := "mainnet"
	switch {
	case *ropsten:
		network = "ropsten"
	case *rinkeby:
		network = "rinkeby"
	default:
		break
	}

	// Connect and deploy
	clientURL := "wss://" + network + ".infura.io/ws"
	client, err := ethclient.Dial(clientURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	// Get token contract instance
	var address common.Address
	var contract *token.Token
	if *addr == "0x0" {
		auth, err := bind.NewTransactor(strings.NewReader(key), "")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		address, _, contract, err = token.DeployToken(auth, client, "Test", "TST")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		fmt.Printf("Contract deployed at %v on the %v network\n", address.String(), network)
		if err := migrator.Populate(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
	} else {
		contract, err = token.NewToken(common.HexToAddress(*addr), client)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
	}

	// Set up listener
	var c chan string
	go func() {
		if err := listener.Listen(contract, client, c); err != nil {
			fmt.Fprintf(os.Stderr, "listener error: %v", err)
			return
		}
	}()

	// Listen for ctrl+c
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// Event handler
	for {
		select {
		case <-quit:
			// c <- "stop"
			return
		case message := <-c:
			switch message {
			case "freeze":

			case "lock":

			case "migrate":
				c <- "stop"
				if err := migrator.PackUp(); err != nil {
					fmt.Fprintf(os.Stderr, "migration error during packup: %v", err)
				}

				return
			}
		}
	}
}
