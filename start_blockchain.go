package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

//	Block : This interface contains all the data
//	related to a transaction.
type Block struct {
	data         map[string]interface{}
	hash         string
	previousHash string
	timestamp    time.Time
	pow          int
}

// Blockchain : This interface contains all the blocks
// and its data.
type Blockchain struct {
	genesisBlock Block
	chain        []Block
	difficulty   int
}

// calculateHash : This method takes the block data as an
// input and generates a hash using SHA256 algorithm.
func (b Block) calculateHash() string {
	data, _ := json.Marshal(b.data)
	blockData := b.previousHash + string(data) + b.timestamp.String() + strconv.Itoa(b.pow)
	blockHash := sha256.Sum256([]byte(blockData))
	return fmt.Sprintf("%x", blockHash)
}

// mine : This method finds a valid hash by
// incrementing the pow and calculating the block hash.
func (b *Block) mine(difficulty int) {
	for !strings.HasPrefix(b.hash, strings.Repeat("0", difficulty)) {
		b.pow++
		b.hash = b.calculateHash()
	}
}

// CreateBlockchain : This method creates the genesis block
// and initializes the chain for the first time.
func CreateBlockchain(difficulty int) Blockchain {

	// Create the genesis block with default values.
	genesisBlock := Block{
		hash:      "0",
		timestamp: time.Now(),
	}

	// Return the chain with the genesis (first) block.
	return Blockchain{
		genesisBlock,
		[]Block{genesisBlock},
		difficulty,
	}
}

// addBlock : This method is used to add a new block to the chain
// using the input data parameters for the transaction.
func (b *Blockchain) addBlock(from, to string, amount float64) {

	// Create the block data interface with the input parameters.
	blockData := map[string]interface{}{
		"from":   from,
		"to":     to,
		"amount": amount,
	}

	// Get the last block details.
	lastBlock := b.chain[len(b.chain)-1]

	// Create the new block with lastBlock as an input.
	newBlock := Block{
		data:         blockData,
		previousHash: lastBlock.hash,
		timestamp:    time.Now(),
	}

	// Generate the block hash using the mine function and
	// provided difficulty parameter from the Blockchain.
	newBlock.mine(b.difficulty)

	// Append the newly created block to the existing chain.
	b.chain = append(b.chain, newBlock)
}

// isValid : This method is used to validate the blocks in the chain
// by verifying the hash of the current and previous block.
func (b Blockchain) isValid() bool {

	// Traverese through the chain excluding the genesis block.
	for i := range b.chain[1:] {
		previousBlock := b.chain[i]
		currentBlock := b.chain[i+1]

		// Check the validation condition
		if currentBlock.hash != currentBlock.calculateHash() || currentBlock.previousHash != previousBlock.hash {
			return false
		}
	}
	return true
}

func main() {

	// Create a new blockchain instance with mining difficulty of 2.
	blockchain := CreateBlockchain(2)
	fmt.Println("Blockchain created successfully!")

	// Record transactions on the blockchain for the Alice, Bob and John.
	blockchain.addBlock("Alice", "Bob", 2)
	blockchain.addBlock("John", "Bob", 3)

	// Check if the blockchain is valid
	if blockchain.isValid() {
		fmt.Println("Added blocks: ")
		for i := range blockchain.chain[1:] {
			fmt.Println(blockchain.chain[i+1].hash, " - Block[", i+1, "]")
		}
	} else {
		fmt.Println("Block verfication failed. Please re-create the chain.")
	}
}
