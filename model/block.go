package model

import (
	"gorm.io/gorm"
	"math/big"
)

type Block struct {
	gorm.Model
	Timestamp     int64  `json:"timestamp"`
	Data          []byte `json:"data" gorm:"-"`
	PrevBlockHash string `json:"prevBlockHash"`
	Hash          string `json:"hash" gorm:"unique"`
	Nonce         int    `json:"nonce"`
	Version       uint   `json:"version" gorm:"default:0"`
	ChainId       string `json:"chainId" gorm:"default:0"`
}

type Blockchain struct {
	Id     string
	Blocks []*Block
}

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}
