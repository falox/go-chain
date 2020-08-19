[![Build Status](https://travis-ci.com/falox/go-chain.svg?branch=master)](https://travis-ci.com/falox/go-chain)

# go-chain, a blockchain implementation in Go

A [blockchain](https://en.wikipedia.org/wiki/Blockchain) is a list of records, called *blocks*, that are linked using cryptography. Each block contains a cryptographic hash of the previous block, a timestamp, and transaction data.

*go-chain* implements some typical blockchain data structures and algorithms in Go:

- Blocks
- Proof-of-Work (PoW)
- Blockchain validation (prevents tampering)
- Concurrent mining

The code is inspired by [Savjee's blockchain Javascript implementation](https://github.com/Savjee/SavjeeCoin).

> NOTE: The code is for education purposes only. This is not a complete and secure implementation.

## Instructions

Download the source code:

```bash
git clone https://github.com/falox/go-chain.git
```

You can run the program with:

```bash
go run ./...
```

You can run the unit tests with:

```bash
go test -v ./...
```

### Mining parameters

You can change the difficulty of the proof-of-work algorithm by editing the `difficulty` global variable:

```go
var difficulty = 2
```

The higher the difficulty is, the more time-consuming the algorithm will be. Minimum is 1. 

Depending on the hardware, you can speed up the mining by enabling concurrency. Set the `miners` global variable to a value greater then 1:

```go
var miners = 3 // will spawn 3 threads
```

## Next Steps

- Extend `Block.Data` to support *transactions* (probably using a [Merkle tree](https://en.wikipedia.org/wiki/Merkle_tree))
- Add wallet generation (private/public key) and the ability to sign a transaction
