# Blockchain

The `internal/blockchain` package implements a minimal proof‑of‑work blockchain.
Each block stores arbitrary data, the hash of the previous block and a nonce
found during mining. Blocks are verified by checking the hash chain and mining
rules.
