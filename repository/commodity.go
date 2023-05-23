package repository

import (
	"P/model"
	"gorm.io/gorm"
)

type CommodityRepository struct {
	DB *gorm.DB
}
type CommodityRepoInterface interface {
	AddCommodities(cs []*model.Commodity, tx ...*gorm.DB) error
	FindComSByOrderNum(num string, tx ...*gorm.DB) ([]*model.Commodity, error)
	FindComSByUserId(userId string) ([]*model.Commodity, error)
	UpdateCommodities(cs []*model.Commodity, order *model.Order, tx ...*gorm.DB) error
	DelCommodities(cs []*model.Commodity, tx ...*gorm.DB) error
}

func (repo *CommodityRepository) AddCommodities(cs []*model.Commodity, tx ...*gorm.DB) error {
	//将指针复制给新变量以便使用tx时不影响原数据库指针
	db := repo.DB
	//判断tx是否传入
	if len(tx) != 0 {
		//若传入则将变量的值替换为tx
		db = tx[0]
	}
	//执行方法
	err := db.Create(cs).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *CommodityRepository) FindComSByOrderNum(num string, tx ...*gorm.DB) ([]*model.Commodity, error) {
	db := repo.DB
	if len(tx) != 0 {
		db = tx[0]
	}
	var comS []*model.Commodity
	err := db.Where("order_num=?", num).Find(&comS).Error
	if err != nil {
		return nil, err
	}
	return comS, nil
}

func (repo *CommodityRepository) FindComSByUserId(userId string) ([]*model.Commodity, error) {
	var comS []*model.Commodity
	err := repo.DB.Where("user_id=?", userId).Find(&comS).Error
	if err != nil {
		return nil, err
	}
	return comS, nil
}

func (repo *CommodityRepository) UpdateCommodities(cs []*model.Commodity, order *model.Order, tx ...*gorm.DB) error {
	db := repo.DB
	if len(tx) != 0 {
		db = tx[0]
	}
	var hashes []string
	for _, c := range cs {
		hashes = append(hashes, c.Hash)
	}
	m := map[string]interface{}{"order_num": order.OrderNum, "user_id": order.BuyerId}
	err := db.Model(&model.Commodity{}).Where("hash IN (?)", hashes).Updates(m).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *CommodityRepository) DelCommodities(cs []*model.Commodity, tx ...*gorm.DB) error {
	db := repo.DB
	if len(tx) != 0 {
		db = tx[0]
	}
	err := db.Delete(cs).Error
	if err != nil {
		return err
	}
	return nil
}
