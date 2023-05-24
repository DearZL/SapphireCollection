package service

import (
	"P/model"
	"P/repository"
	"errors"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

type BlockServiceInterface interface {
	GetBlockChain(chain *model.Blockchain) error
	AddBlock(cs *model.Commodities, chainId string, tx1 ...*gorm.DB) error
}

type BlockService struct {
	BlockRepo repository.BlockRepoInterface
}

func (srv *BlockService) GetBlockChain(chain *model.Blockchain) error {
	return srv.BlockRepo.GetBlockChain(chain)
}

func (srv *BlockService) AddBlock(cs *model.Commodities, chainId string, tx0 ...*gorm.DB) error {
	location := "./file/commodity/"
	latestBlock := &model.Block{}
	err := srv.BlockRepo.GetLastBlock(latestBlock)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	chain := &model.Blockchain{Id: chainId}
	chain.Blocks = append(chain.Blocks, latestBlock)
	for i, com := range cs.Commodities {
		imagePath := location + com.Image
		imageData, err := os.ReadFile(imagePath)
		if err != nil {
			log.Println("第", i, "个商品图片FileName:", imagePath, "读取失败,添加商品失败，请重试")
			return err
		}
		chain.AddBlock(string(imageData))
		time.Sleep(1000 * time.Microsecond)
		com.Hash = chain.Blocks[len(chain.Blocks)-1].Hash

	}
	var tx *gorm.DB
	if len(tx0) != 0 {
		tx = tx0[0]
	} else {
		tx = srv.BlockRepo.GetDB().Begin()
	}
	err = srv.BlockRepo.AddBlock(chain, tx)
	if err != nil {
		tx.Rollback()
		log.Println("rollback")
		return err
	}
	result := tx.Model(latestBlock).
		Where("hash=? AND version = ?",
			latestBlock.Hash, latestBlock.Version).
		Update("version", latestBlock.Version+1)
	if result.RowsAffected == 0 {
		tx.Rollback()
		log.Println("rollback")
		return errors.New("事务执行期间已有其他事务提交,请放弃或重试")
	}
	if len(tx0) == 0 {
		err = tx.Commit().Error
		if err != nil {
			tx.Rollback()
			log.Println("rollback")
			return err
		}
	}
	return nil
}
