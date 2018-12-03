package listener

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/tokenise-eu/tokenise-apps/apps/eth/token"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Listen will start listening on the ethereum blockchain for events
// emitted by the contract.
func Listen(contract *token.Token, client *ethclient.Client, c chan string) error {

	watchOpts := bind.WatchOpts{
		Start:   nil,
		Context: nil,
	}

	lockChan := make(chan *token.TokenLock)
	lockSub, err := contract.TokenFilterer.WatchLock(&watchOpts, lockChan, []common.Address{})
	if err != nil {
		return err
	}

	freezeChan := make(chan *token.TokenFreeze)
	freezeSub, err := contract.TokenFilterer.WatchFreeze(&watchOpts, freezeChan)
	if err != nil {
		return err
	}

	migrateChan := make(chan *token.TokenMigrate)
	migrateSub, err := contract.TokenFilterer.WatchMigrate(&watchOpts, migrateChan)
	if err != nil {
		return err
	}

	transferChan := make(chan *token.TokenTransfer)
	transferSub, err := contract.TokenFilterer.WatchTransfer(&watchOpts, transferChan, []common.Address{}, []common.Address{})
	if err != nil {
		return err
	}

	addChan := make(chan *token.TokenVerifiedAddressAdded)
	addSub, err := contract.TokenFilterer.WatchVerifiedAddressAdded(&watchOpts, addChan, []common.Address{}, []common.Address{})
	if err != nil {
		return err
	}

	updateChan := make(chan *token.TokenVerifiedAddressUpdated)
	updateSub, err := contract.TokenFilterer.WatchVerifiedAddressUpdated(&watchOpts, updateChan, []common.Address{}, []common.Address{})
	if err != nil {
		return err
	}

	removeChan := make(chan *token.TokenVerifiedAddressRemoved)
	removeSub, err := contract.TokenFilterer.WatchVerifiedAddressRemoved(&watchOpts, removeChan, []common.Address{}, []common.Address{})
	if err != nil {
		return err
	}

	fmt.Println("Listener set up complete")
	for {
		select {
		case err := <-lockSub.Err():
			return err
		case err := <-freezeSub.Err():
			return err
		case err := <-migrateSub.Err():
			return err
		case err := <-transferSub.Err():
			return err
		case err := <-addSub.Err():
			return err
		case err := <-updateSub.Err():
			return err
		case err := <-removeSub.Err():
			return err
		case lock := <-lockChan:
			/*if err := db.AddTx(lock.Addr.Bytes()); err != nil {
				return err
			}*/
			fmt.Println(lock)
		case freeze := <-freezeChan:
			fmt.Println(freeze)
		case migrate := <-migrateChan:
			fmt.Println(migrate)
		case transfer := <-transferChan:
			/*if err := db.AddTx(transfer.From.String(), transfer.To.String(), transfer.Value.String()); err != nil {
				return err
			}*/
			fmt.Println(transfer)
		case add := <-addChan:
			/*if err := db.AddUser(add.Hash[:]); err != nil {
				return err
			}*/
			fmt.Println(add)
		case update := <-updateChan:
			fmt.Println(update)
		case remove := <-removeChan:
			fmt.Println(remove)
		case <-c:
			return nil
		}
	}
}
