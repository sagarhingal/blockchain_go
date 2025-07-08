package blockchain

import (
	"testing"

	"blockchain_go/internal/contracts"
)

func TestAddBlockAndValidate(t *testing.T) {
	bc := CreateBlockchain(1)
	bc.AddBlock("Alice", "Bob", 1)
	bc.AddBlock("Bob", "Charlie", 2)
	if !bc.IsValid() {
		t.Fatal("expected chain to be valid")
	}
}

func TestAddContractBlock(t *testing.T) {
	bc := CreateBlockchain(1)
	bc.AddContractBlock(&contracts.TokenTransfer{From: "A", To: "B", Amount: 1})
	if !bc.IsValid() {
		t.Fatal("expected chain with contract to be valid")
	}
}
