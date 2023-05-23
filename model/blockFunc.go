package model

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"time"
)

const targetBits = 0

func NewBlockchain(chainId string) *Blockchain {
	return &Blockchain{
		Blocks: []*Block{
			NewGenesisBlock(chainId),
		},
	}
}

func NewBlock(data string, prevBlockHash string, chainId string) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Data:          []byte(data),
		PrevBlockHash: prevBlockHash,
		Hash:          "",
		ChainId:       chainId,
	}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	block.Hash = fmt.Sprintf("%x", hash[:])
	block.Nonce = nonce
	return block
}

func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash, bc.Id)
	bc.Blocks = append(bc.Blocks, newBlock)
}

func NewGenesisBlock(chainId string) *Block {
	return NewBlock("Genesis Block", "", chainId)
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0
	fmt.Println("Mining the block containing")
	for nonce < math.MaxInt64 {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(pow.Target) == -1 {
			break
		} else {
			nonce++
		}
	}
	return nonce, hash[:]
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))
	pow := &ProofOfWork{b, target}
	return pow
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			[]byte(pow.Block.PrevBlockHash),
			pow.Block.Data,
			[]byte(strconv.FormatInt(pow.Block.Timestamp, 10)),
			[]byte(strconv.FormatInt(int64(targetBits), 10)),
			[]byte(strconv.FormatInt(int64(nonce), 10)),
		},
		[]byte{},
	)
	return data
}

func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int
	data := pow.prepareData(pow.Block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])
	isValid := hashInt.Cmp(pow.Target) == -1
	return isValid
}

//func exmple() {
//	bc := NewBlockchain("0")
//	bc.AddBlock("Send 1 BTC to Ivan")
//	bc.AddBlock("Send 2 more BTC to Ivan")
//
//	pow := NewProofOfWork(bc.Blocks[1])
//	pow1 := NewProofOfWork(bc.Blocks[2])
//
//	for _, block := range bc.Blocks {
//		fmt.Printf("Prev. hash: %b\n", block.PrevBlockHash)
//		fmt.Printf("Data: %s\n", block.Data)
//		fmt.Printf("Hash: %b\n", block.Hash)
//		fmt.Printf("Nonce: %d\n", block.Nonce)
//
//		fmt.Println()
//	}
//	fmt.Println(pow.Target)
//	fmt.Println(pow1.Target)
//
//}
