package repository

import (
	"P/model"
	"gorm.io/gorm"
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
	err := db.Order("id desc").First(&latestBlock).Error
	if err != nil {
		return err
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
	//var Blocks []*model.Block
	//err := db.Where("chain_id=?", chain.Id).Find(Blocks).Error
	//if err != nil {
	//	return err
	//}
	//chain.Blocks = Blocks
	//return nil

	db.Create(chain.Blocks)
	return nil
}
