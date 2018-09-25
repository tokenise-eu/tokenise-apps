package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

func main() {
	var currentBlock uint64

	e := configure()
	check(e)

	clientURL := "wss://" + network + ".infura.io/ws"
	client, err := ethclient.Dial(clientURL)
	check(err)

	newestBlock, err := client.HeaderByNumber(context.Background(), nil)
	check(err)

	lastBlock, err := readBlock()
	check(err)

	contractAddress := common.HexToAddress(address)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
		FromBlock: lastBlock,
		ToBlock: newestBlock.Number,
	}

	// Signature declaration
	freeze := newEvent("Freeze()")
	deployed := newEvent("Deployed(address)")
	ready := newEvent("Ready()")
	migrate := newEvent("Migrate()")
	lock := newEvent("Lock(address)")

	logs := make(chan types.Log)

	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	check(err)

	fmt.Println("Subscribed")

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:
			if vLog.BlockNumber != currentBlock {
				currentBlock = vLog.BlockNumber
				e := writeBlock(string(vLog.BlockNumber))
				check(e)
			}

			eventSignature := vLog.Topics[0].Hex()

			switch eventSignature {
			case freeze.Hex:
				// Signify
			case deployed.Hex:
				fmt.Println(vLog.Topics[1].Hex())
				fmt.Println(tokenListenerPort)
			case ready.Hex:
				fmt.Println(databasePort)
			case migrate.Hex:
				// Signify shutdown
			case lock.Hex:
				lockedAddress := vLog.Topics[1].Hex()
				fmt.Println(lockedAddress)
				// Signify
			default:
				fmt.Println("Emitted event was not recognized: ", eventSignature)
			}
		}
	}
}

func check(err error) {
	if err != nil {
		log.Fatal("Error: ", err)
	}
}
