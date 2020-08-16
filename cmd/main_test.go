package main

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
	difficulty = 1
	miners = 1
}

func TestCreateBlockchain(t *testing.T) {
	// Arrange
	now, _ := time.Parse(time.RFC3339, "2020-08-10T00:00:00.000Z")

	// Act
	result := createBlockchain(now)

	// Assert
	expectedGenesisBlock := Block{now.Format(time.RFC3339), genesisBlockData, strings.Repeat("0", 64), "0f723d11ee6c1bb2609ce20084c1c4e062e015a5d039e3ad7cf2c7eac5d3024c", 5}
	assert.Equal(t, expectedGenesisBlock, result[0])
}

func TestCreateBlock(t *testing.T) {
	// Arrange
	now := time.Now()

	// Act
	result, _ := createBlock(now, "Data", Block{Hash: "PreviousHash"})

	// Assert
	expected := Block{now.Format(time.RFC3339), "Data", "PreviousHash", "", 0}
	assert.Equal(t, expected, result)
}

func TestCreateBlockWithNoData(t *testing.T) {
	// Act
	_, err := createBlock(time.Now(), "", Block{})

	// Assert
	assert.NotEqual(t, nil, err)
}

func TestAddBlock(t *testing.T) {
	// Arrange
	blockchain := createBlockchain(time.Now())
	block, _ := createBlock(time.Now(), "block", getLatestBlock(blockchain))

	// Act
	err := addBlock(&blockchain, block)

	// Assert
	assert.Equal(t, nil, err)
	assert.NotEqual(t, "", getLatestBlock(blockchain).Hash)
	assert.Equal(t, 2, len(blockchain))
}

func TestAddBlockWithInvalidBLock(t *testing.T) {
	// Arrange
	blockchain := createBlockchain(time.Now())
	block, _ := createBlock(time.Now(), "detached block", Block{})

	// Act
	err := addBlock(&blockchain, block)

	// Assert
	assert.NotEqual(t, nil, err)
	assert.Equal(t, 1, len(blockchain))
}

func TestCalculateHash(t *testing.T) {
	// Arrange
	block := Block{Data: "Data", Timestamp: "Timestamp", PreviousHash: "PreviousHash", Nonce: 123}

	// Act
	result := calculateHash(block)

	// Assert
	const EXPECTED = "ecfd4574e62af7d96875cf1a96b59f784aed4937387f5a62cf8691739401ec37"
	assert.Equal(t, EXPECTED, result)
}

func TestCalculateHashWithSameInput(t *testing.T) {
	// Arrange
	block1 := Block{Data: "Data", Timestamp: "Timestamp", PreviousHash: "PreviousHash", Nonce: 123}
	block2 := Block{Data: "Data", Timestamp: "Timestamp", PreviousHash: "PreviousHash", Nonce: 123}

	// Act
	result1 := calculateHash(block1)
	result2 := calculateHash(block2)

	// Assert
	assert.Equal(t, result1, result2)
}

func TestCalculateHashWithDifferentInput(t *testing.T) {
	// Arrange
	block := Block{Data: "Data", Timestamp: "Timestamp", PreviousHash: "PreviousHash", Nonce: 123}
	result1 := calculateHash(block)

	// Act
	block.Nonce = 124
	result2 := calculateHash(block)

	// Assert
	assert.NotEqual(t, result1, result2)
}

func TestIsBlockchainValid(t *testing.T) {
	// Arrange
	blockchain := createBlockchain(time.Now())
	block, _ := createBlock(time.Now(), "block", getLatestBlock(blockchain))
	addBlock(&blockchain, block)

	// Act
	result := isBlockchainValid(blockchain)

	// Assert
	assert.Equal(t, true, result)
}

func TestIsBlockchainValidWithGenesisBlockOnly(t *testing.T) {
	// Arrange
	blockchain := createBlockchain(time.Now())

	// Act
	result := isBlockchainValid(blockchain)

	// Assert
	assert.Equal(t, true, result)
}

func TestIsBlockchainValidWithTamperedGenesisBlock(t *testing.T) {
	// Arrange
	blockchain := createBlockchain(time.Now())
	blockchain[0].Data = "TAMPERED"

	// Act
	result := isBlockchainValid(blockchain)

	// Assert
	assert.Equal(t, false, result)
}

func TestIsBlockchainValidWithTamperedBlock(t *testing.T) {
	// Arrange
	blockchain := createBlockchain(time.Now())
	block, _ := createBlock(time.Now(), "block", getLatestBlock(blockchain))
	addBlock(&blockchain, block)
	blockchain[len(blockchain)-1].Data = "TAMPERED"

	// Act
	result := isBlockchainValid(blockchain)

	// Assert
	assert.Equal(t, false, result)
}

func TestIsBlockchainValidWithEmptyBlock(t *testing.T) {
	// Act
	result := isBlockchainValid([]Block{})

	// Assert
	assert.Equal(t, false, result)
}
