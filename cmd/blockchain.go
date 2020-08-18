package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
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

// Difficulty sets how difficult mining is. (1 = easiest)
var difficulty = 1

// Miners sets the number of concurrent mining threads. (1 = no concurrency)
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
	// A valid blockchain contains at least one block
	if len(blockchain) == 0 {
		return false
	}

	// Validate the genesis (first) block
	if len(blockchain) > 0 && blockchain[0].Data != genesisBlockData {
		return false
	}

	// Validate next blocks
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
	// A new block must be "mined" before to be added to the blockchain. Mining is a proof-of-work (https://en.wikipedia.org/wiki/Proof_of_work)
	minedBlock := mineBlock(block)

	candidateBlockchain := append(*blockchain, minedBlock)
	if isBlockchainValid(candidateBlockchain) == false {
		// If the mined block is inconsistent, the blockchain is not updated
		return errors.New("Invalid block (Previous hash is inconsistent). Blockchain was not updated")
	}

	*blockchain = candidateBlockchain
	return
}

func mineBlock(block Block) Block {
	// Mining is guessing the Block.Nonce until a Block.Hash that matches the targetHashPrefix is found

	if miners <= 1 {
		// difficulty defines the number of 0s leading the hash. The higher difficulty, the more time-consuming
		targetHashPrefix := strings.Repeat("0", difficulty)

		for {
			block.Hash = calculateHash(block)

			if strings.HasPrefix(block.Hash, targetHashPrefix) {
				return block
			}

			block.Nonce++
		}
	} else {
		// Depending on the hardware, concurrency could speed up mining
		result := make(chan Block, miners)
		var wg sync.WaitGroup
		var stop uint32 = 0

		for i := 0; i < miners; i++ {
			wg.Add(1)
			go concurrentMineBlock(block, miners, result, &stop, &wg)
			block.Nonce++ // Every miner starts with a different nonce
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
