package main

import (
	"fmt"

	"blockchain_go/internal/blockchain"
)

func main() {
	bc := blockchain.CreateBlockchain(2)
	fmt.Println("Blockchain created successfully!")

	bc.AddBlock("Alice", "Bob", 2)
	bc.AddBlock("John", "Bob", 3)

	if bc.IsValid() {
		fmt.Println("Added blocks:")
		for i := range bc.Chain[1:] {
			fmt.Printf("%s - Block[%d]\n", bc.Chain[i+1].Hash, i+1)
		}
	} else {
		fmt.Println("Block verification failed. Please re-create the chain.")
	}
}
