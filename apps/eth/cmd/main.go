package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/tokenise-eu/tokenise-apps/apps/eth/db"
	"github.com/tokenise-eu/tokenise-apps/apps/eth/deployer"
	"github.com/tokenise-eu/tokenise-apps/apps/eth/listener"
	"github.com/tokenise-eu/tokenise-apps/apps/eth/migrator"
)

var mainnet = flag.Bool("mainnet", false, "Deploy the contract and listen "+
	"on the Ethereum mainnet (default)")
var ropsten = flag.Bool("ropsten", false, "Deploy the contract and listen "+
	"on the Ropsten testnet")
var rinkeby = flag.Bool("rinkeby", false, "Deploy the contract and listen "+
	"on the Rinkeby testnet")
var addr = flag.String("contract", "0x0", "Specify a contract listen to, "+
	"instead of deploying a new one")

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

	if *addr == "0x0" {
		var err error
		*addr, err = deployer.Deploy(network)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error deploying contract: %v", err)
			return
		}
	}

	var c chan string

	go func() {
		if err := listener.Listen(*addr, network, c); err != nil {
			fmt.Fprintf(os.Stderr, "listener error: %v", err)
			return
		}
	}()

	// Listen for ctrl+c
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT)

	for {
		select {
		case <-quit:
			c <- "stop"
			if err := db.Stop(); err != nil {
				fmt.Fprintf(os.Stderr, "error shutting down database: %v", err)
			}

			return
		case message := <-c:
			switch message {
			case "freeze":

			case "lock":

			case "ready":
				if err := db.Start(); err != nil {
					fmt.Fprintf(os.Stderr, "error starting database: %v", err)
					c <- "stop"
					return
				}
			case "migrate":
				c <- "stop"
				db.Stop()
				if err := migrator.PackUp(); err != nil {
					fmt.Fprintf(os.Stderr, "migration error during packup: %v", err)
				}

				return
			}
		}
	}
}
