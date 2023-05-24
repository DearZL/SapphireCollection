package repository

import (
	"P/model"
	"errors"
	"gorm.io/gorm"
	"log"
)

type BlockRepoInterface interface {
	GetDB() *gorm.DB
	GetLastBlock(latestBlock *model.Block, tx ...*gorm.DB) error
	AddBlock(chain *model.Blockchain, tx ...*gorm.DB) error
	GetBlockChain(chain *model.Blockchain, tx ...*gorm.DB) error
}
type BlockRepository struct {
	DB *gorm.DB
}

func (repo *BlockRepository) GetDB() *gorm.DB {
	db := repo.DB
	return db
}

func (repo *BlockRepository) GetLastBlock(latestBlock *model.Block, tx ...*gorm.DB) error {
	db := repo.DB
	if len(tx) != 0 {
		db = tx[0]
	}
	result := db.Order("id desc").First(&latestBlock)
	if result.Error != nil {
		log.Println(result.Error)
		return errors.New("查询失败")
	}
	if result.RowsAffected == 0 {
		return errors.New("未找到记录,请先创建创世块")
	}
	return nil
}

func (repo *BlockRepository) AddBlock(chain *model.Blockchain, tx ...*gorm.DB) error {
	db := repo.DB
	if len(tx) != 0 {
		db = tx[0]
	}
	err := db.Create(chain.Blocks[1:]).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *BlockRepository) GetBlockChain(chain *model.Blockchain, tx ...*gorm.DB) error {
	db := repo.DB
	if len(tx) != 0 {
		db = tx[0]
	}
	var blocks []*model.Block
	result := db.Where("chain_id=?", chain.Id).Find(&blocks)
	if result.Error != nil {
		log.Println(result.Error)
		return errors.New("查询失败")
	}
	if result.RowsAffected == 0 {
		return errors.New("未找到记录,请检查参数")
	}
	chain.Blocks = blocks
	return nil
}
