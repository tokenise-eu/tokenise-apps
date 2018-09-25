package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strconv"
)

var (
	network           string
	address           string
	databasePort      uint
	tokenListenerPort uint
)

// Config maps the information from the configuration file
type Config struct {
	Network           string
	Address           string
	DatabasePort      uint
	TokenListenerPort uint
}

func configure() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	file, err := os.Open(dir + "/config.json")
	if err != nil {
		return err
	}

	defer file.Close()
	decoder := json.NewDecoder(file)
	config := Config{}
	err2 := decoder.Decode(&config)
	if err2 != nil {
		return err2
	}

	switch config.Network {
	case "mainnet", "ropsten", "rinkeby", "kovan":
		break
	default:
		e := errors.New("network not recognized")
		return e
	}

	network = config.Network
	address = config.Address
	databasePort = config.DatabasePort
	tokenListenerPort = config.TokenListenerPort
	return nil
}

func readBlock() (*big.Int, error) {
	txt, err := ioutil.ReadFile("lastblock.txt")
	if err != nil {
		return nil, err
	}

	str := string(txt)
	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return nil, err
	}

	block := big.NewInt(num)
	return block, nil
}

func writeBlock(number string) error {
	file, err := os.Create("lastblock.txt")
	if err != nil {
		return err
	}

	defer file.Close()
	_, err2 := fmt.Fprintf(file, number)
	if err2 != nil {
		return err
	}
	return nil
}