package listener

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Listen will start listening on the ethereum blockchain for events
// emitted by the contract at address `addr`.
func Listen(addr string, network string, c chan string) error {
	var currentBlock uint64

	clientURL := "wss://" + network + ".infura.io/ws"
	client, err := ethclient.Dial(clientURL)
	if err != nil {
		return err
	}

	newestBlock, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return err
	}

	lastBlock, err := readBlock()
	if err != nil {
		return err
	}

	contractAddress := common.HexToAddress(addr)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
		FromBlock: lastBlock,
		ToBlock:   newestBlock.Number,
	}

	// Signature declaration
	freeze := newEvent("Freeze()")
	ready := newEvent("Ready()")
	migrate := newEvent("Migrate()")
	lock := newEvent("Lock(address)")

	logs := make(chan types.Log)

	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		return err
	}

	c <- "subscribed"

	for {
		select {
		case err := <-sub.Err():
			return err
		case vLog := <-logs:
			if vLog.BlockNumber != currentBlock {
				currentBlock = vLog.BlockNumber
				if err := writeBlock(string(vLog.BlockNumber)); err != nil {
					return err
				}
			}

			eventSignature := vLog.Topics[0].Hex()
			switch eventSignature {
			case freeze.Hex:
				c <- "freeze"
			case ready.Hex:
				//
			case migrate.Hex:
				//
			case lock.Hex:
				lockedAddress := vLog.Topics[1].Hex()
				fmt.Println(lockedAddress)
				// Signify
			default:
				fmt.Println("Emitted event was not recognized: ", eventSignature)
			}
		case <-c:
			return nil
		}
	}
}
