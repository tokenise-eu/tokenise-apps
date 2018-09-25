package main

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// Event maps all the needed information for a contract event
type Event struct {
	Signature []byte
	Hash      common.Hash
	Hex       string
}

func newEvent(name string) *Event {
	var e Event
	e.Signature = []byte(name)
	e.setHash(name)
	e.setHex(e.Hash)
	return &e
}

func (e *Event) setHash(name string) {
	byteName := []byte(name)
	e.Hash = crypto.Keccak256Hash(byteName)
}

func (e *Event) setHex(hash common.Hash) {
	e.Hex = hash.Hex()
}
