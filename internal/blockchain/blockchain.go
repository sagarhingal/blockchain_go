package blockchain

import (
	"blockchain_go/internal/contracts"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Block represents a single block in the chain. A block may contain
// transaction data or a smart contract execution and is linked to the
// previous block via PrevHash.
type Block struct {
	Data      map[string]interface{}
	Hash      string
	PrevHash  string
	Timestamp time.Time
	Contract  contracts.Contract
	Nonce     int
}

// Blockchain stores an ordered list of blocks beginning with a genesis block.
// Difficulty controls the Proof-of-Work mining target.
type Blockchain struct {
	Chain      []Block
	Difficulty int
}

// calculateHash generates the SHA-256 hash of the block fields.
func (b Block) calculateHash() string {
	data, _ := json.Marshal(b.Data)
	blockData := b.PrevHash + string(data) + b.Timestamp.String() + strconv.Itoa(b.Nonce)
	blockHash := sha256.Sum256([]byte(blockData))
	return fmt.Sprintf("%x", blockHash)
}

// mine repeatedly hashes the block with an incrementing nonce until the
// resulting hash has the required number of leading zeros.
func (b *Block) mine(difficulty int) {
	for !strings.HasPrefix(b.Hash, strings.Repeat("0", difficulty)) {
		b.Nonce++
		b.Hash = b.calculateHash()
	}
}

// CreateBlockchain initializes a new blockchain with a genesis block and
// returns it. The provided difficulty determines how many leading zeros are
// required in a block hash during mining.
func CreateBlockchain(difficulty int) Blockchain {
	// Create the genesis block with default values.
	genesisBlock := Block{
		Hash:      "0",
		Timestamp: time.Now(),
	}

	// Initialize the chain containing the genesis block only.
	return Blockchain{
		Chain:      []Block{genesisBlock},
		Difficulty: difficulty,
	}
}

// AddBlock mines and appends a new block containing the given
// transaction data to the chain.
func (b *Blockchain) AddBlock(from, to string, amount float64) {

	// Create the block Data interface with the input parameters.
	blockData := map[string]interface{}{
		"from":   from,
		"to":     to,
		"amount": amount,
	}

	// Get the last block details.
	lastBlock := b.Chain[len(b.Chain)-1]

	// Create the new block with lastBlock as an input.
	newBlock := Block{
		Data:      blockData,
		PrevHash:  lastBlock.Hash,
		Timestamp: time.Now(),
	}

	// Generate the block hash using the mine function and
	// provided Difficulty parameter from the Blockchain.
	newBlock.mine(b.Difficulty)

	// Append the newly created block to the existing chain.
	b.Chain = append(b.Chain, newBlock)
}

// IsValid verifies the hashes of every block and returns true if the
// entire chain is intact.
func (b Blockchain) IsValid() bool {

	// Traverse through the chain excluding the genesis block.
	for i := range b.Chain[1:] {
		previousBlock := b.Chain[i]
		currentBlock := b.Chain[i+1]

		// Check the validation condition
		if currentBlock.Hash != currentBlock.calculateHash() || currentBlock.PrevHash != previousBlock.Hash {
			return false
		}
	}
	return true
}

// AddContractBlock mines and appends a block that executes the provided
// smart contract.
func (b *Blockchain) AddContractBlock(c contracts.Contract) {
	lastBlock := b.Chain[len(b.Chain)-1]
	newBlock := Block{
		PrevHash:  lastBlock.Hash,
		Timestamp: time.Now(),
		Contract:  c,
	}
	newBlock.mine(b.Difficulty)
	b.Chain = append(b.Chain, newBlock)
}
