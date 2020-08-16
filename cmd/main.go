package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const genesisBlockData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks" // https://en.bitcoin.it/wiki/Genesis_block

// Block is the element of the Blockchain.
type Block struct {
	Timestamp    string
	Data         string
	PreviousHash string
	Hash         string
	Nonce        int
}

// How difficult and time-consuming it is to find the right hash for each block. (1 = easiest)
var difficulty = 1

// Defines the number of concurrent mining threads. (1 = no concurrency)
var miners = 1

func createBlockchain(timestamp time.Time) []Block {
	genesisBlock := Block{}
	genesisBlock.Timestamp = timestamp.Format(time.RFC3339)
	genesisBlock.PreviousHash = strings.Repeat("0", 64)
	genesisBlock.Data = genesisBlockData
	genesisBlock.Nonce = 0

	blockchain := make([]Block, 0, 10)

	addBlock(&blockchain, genesisBlock)

	return blockchain
}

func createBlock(timestamp time.Time, data string, previousBlock Block) (block Block, err error) {
	if len(data) == 0 {
		err = errors.New("Cannot create a block with no data")
		return
	}

	block = Block{}
	block.Timestamp = timestamp.Format(time.RFC3339)
	block.PreviousHash = previousBlock.Hash
	block.Data = data
	block.Nonce = 0
	return
}

func isBlockchainValid(blockchain []Block) bool {
	// Empty blockchain is not valid
	if len(blockchain) == 0 {
		return false
	}

	// Validate the genesis block
	if len(blockchain) > 0 && blockchain[0].Data != genesisBlockData {
		return false
	}

	// Validate the subsequent blocks
	for i := 1; i < len(blockchain); i++ {
		currentBlock := blockchain[i]
		previousBlock := blockchain[i-1]

		if currentBlock.Hash != calculateHash(currentBlock) {
			return false
		}

		if currentBlock.PreviousHash != previousBlock.Hash {
			return false
		}
	}

	return true
}

func addBlock(blockchain *[]Block, block Block) (err error) {
	minedBlock := mineBlock(block)

	candidateBlockchain := append(*blockchain, minedBlock)
	if isBlockchainValid(candidateBlockchain) == false {
		return errors.New("Invalid block (Previous hash is inconsistent with the chain). Blockchain was not updated")
	}

	*blockchain = candidateBlockchain
	return
}

func mineBlock(block Block) Block {
	if miners <= 1 {
		targetHashPrefix := strings.Repeat("0", difficulty)

		for {
			block.Hash = calculateHash(block)
			if strings.HasPrefix(block.Hash, targetHashPrefix) {
				return block
			}
			block.Nonce++
		}
	} else {
		result := make(chan Block, miners)
		var wg sync.WaitGroup
		var stop uint32 = 0

		for i := 0; i < miners; i++ {
			wg.Add(1)
			go concurrentMineBlock(block, miners, result, &stop, &wg)
			block.Nonce++ // every miner starts with a different nonce
		}

		wg.Wait()
		close(result)

		return <-result
	}
}

func concurrentMineBlock(block Block, nonceIncrementStep int, result chan Block, stop *uint32, wg *sync.WaitGroup) {
	defer wg.Done()

	isMiningRunning := func() bool {
		return atomic.LoadUint32(stop) == 0
	}

	stopMining := func() {
		atomic.StoreUint32(stop, 1)
	}

	targetHashPrefix := strings.Repeat("0", difficulty)

	for isMiningRunning() {
		block.Hash = calculateHash(block)

		if strings.HasPrefix(block.Hash, targetHashPrefix) {
			result <- block
			stopMining()
			break
		}

		block.Nonce = block.Nonce + nonceIncrementStep
	}
}

func calculateHash(block Block) string {
	data := block.Timestamp + block.Data + block.PreviousHash + strconv.Itoa(block.Nonce)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func getLatestBlock(blockchain []Block) Block {
	// ASSUMPTION: Blockchain has at least 1 (genesis) block
	return blockchain[len(blockchain)-1]
}

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
