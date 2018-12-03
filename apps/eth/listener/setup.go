package listener

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strconv"
)

// Config maps the information from the configuration file
type Config struct {
	Network           string
	Address           string
	DatabasePort      uint
	TokenListenerPort uint
}

func ReadBlock() (*big.Int, error) {
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

func WriteBlock(number string) error {
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
