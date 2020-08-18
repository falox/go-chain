package main

import (
	"fmt"
	"time"
)

func dumpBlockchain(blockchain []Block) {
	for i, block := range blockchain {
		fmt.Printf("Block %d: \"%s\" (%s)\n", i+1, block.Data, block.Hash)
	}

	if isBlockchainValid(blockchain) {
		fmt.Printf("Blockchain is VALID\n")
	} else {
		fmt.Printf("Blockchain is NOT VALID\n")
	}
}

func main() {
	blockchain := createBlockchain(time.Now())

	block, _ := createBlock(time.Now(), "first block", getLatestBlock(blockchain))
	addBlock(&blockchain, block)

	block, _ = createBlock(time.Now(), "second block", getLatestBlock(blockchain))
	addBlock(&blockchain, block)

	dumpBlockchain(blockchain)
}
