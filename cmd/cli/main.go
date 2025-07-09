package main

import (
	"os"

	"blockchain_go/internal/blockchain"
	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	bc := blockchain.CreateBlockchain(2)
	logger.Info().Msg("Blockchain created successfully!")

	bc.AddBlock("Alice", "Bob", 2)
	bc.AddBlock("John", "Bob", 3)

	if bc.IsValid() {
		logger.Info().Msg("Added blocks:")
		for i := range bc.Chain[1:] {
			logger.Info().Msgf("%s - Block[%d]", bc.Chain[i+1].Hash, i+1)
		}
	} else {
		logger.Error().Msg("Block verification failed. Please re-create the chain.")
	}
}
