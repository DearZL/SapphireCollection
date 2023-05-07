package model

import (
	"gorm.io/gorm"
	"math/big"
)

type Block struct {
	gorm.Model
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

type Blockchain struct {
	Blocks []*Block
}

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}
