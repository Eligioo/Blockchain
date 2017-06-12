package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"
)

//Header for a block
type Header struct {
	BlockNumber  int
	Nonce        int
	Hash         string
	PreviousHash string
	CreatedAt    time.Time
}

//Transaction in the body
type Transaction struct {
	Sender   string
	Reciever string
	Amount   int
}

//Block itself
type Block struct {
	Header Header
	Body   []Transaction
}

//Chain blocks
type Chain struct {
	Blocks []Block
}

//LastBlock returns the last block of a chain
func (c Chain) LastBlock() *Block {
	return &c.Blocks[len(c.Blocks)-1]
}

func generateBlock(chain *Chain) {
	if len(chain.Blocks) == 0 {
		Genesis := Block{}
		fmt.Println(Genesis)
		Genesis.Header.CreatedAt = time.Now()
		Genesis.Header.BlockNumber = 0
		for index := 0; index < 64; index++ {
			Genesis.Header.PreviousHash += "0"
		}
		Genesis.Header.Hash = findHash(&Genesis)
		chain.Blocks = append(chain.Blocks, Genesis)
	} else {
		newBlock := Block{}
		newBlock.Header.CreatedAt = time.Now()
		lastBlock := chain.Blocks[len(chain.Blocks)-1]
		newBlock.Header.PreviousHash = lastBlock.Header.Hash
		newBlock.Header.BlockNumber = lastBlock.Header.BlockNumber + 1
		newBlock.Header.Hash = findHash(&newBlock)
		chain.Blocks = append(chain.Blocks, newBlock)
	}
}

func findHash(block *Block) string {
	byteChar := fmt.Sprintf("%v", block)
	h := sha256.Sum256([]byte(byteChar))
	stringToHash := fmt.Sprintf("%x", h)
	zeroPrefix := "00000"
	for stringToHash[0:len(zeroPrefix)] != zeroPrefix {
		block.Header.Nonce++
		newByteChar := fmt.Sprintf("%v", block)
		newH := &h
		*newH = sha256.Sum256([]byte(newByteChar))
		newStringToHash := &stringToHash
		*newStringToHash = fmt.Sprintf("%x", h)
		block.Header.Hash = stringToHash
	}
	findingHashDuration := time.Now().Sub(block.Header.CreatedAt)
	durationInSeconds := strconv.FormatFloat(findingHashDuration.Seconds(), 'g', -1, 64)
	fmt.Println("Found hash for block: " + strconv.Itoa(block.Header.BlockNumber) + ". It took " + durationInSeconds + " seconds. " + stringToHash)
	return stringToHash
}

//MakeTransaction adds a transaction to a block, most likely the last block in the chain
func MakeTransaction(s, r string, a int) Transaction {
	return Transaction{Sender: s, Reciever: r, Amount: a}
}

//AddTransaction adds a transaction to the last block in the chain
func AddTransaction(chain *Chain, s, r string, a int) {
	newTransaction := MakeTransaction(s, r, a)
	lastBlock := chain.LastBlock()
	lastBlock.Body = append(lastBlock.Body, newTransaction)
}

func main() {
	chain := Chain{}
	generateBlock(&chain)
	AddTransaction(&chain, "daan", "stefan", 12)
	AddTransaction(&chain, "daan", "rick", 22)
	generateBlock(&chain)
	AddTransaction(&chain, "stefan", "daan", 100)
	AddTransaction(&chain, "stefan", "rick", 50)
	generateBlock(&chain)
	json, _ := json.Marshal(chain)
	ioutil.WriteFile("chain.json", json, 0644)
}
