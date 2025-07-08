package contracts

// Contract represents a simple interface for smart contracts
// executed on blocks.
type Contract interface {
	Execute() error
}

// TokenTransfer is a basic implementation that could perform
// token transfer logic in a real blockchain. Here it's a stub
// used for demonstration purposes.
type TokenTransfer struct {
	From   string
	To     string
	Amount float64
}

func (t *TokenTransfer) Execute() error {
	// In a real blockchain this would update balances.
	return nil
}
