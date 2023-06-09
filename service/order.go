package service

import (
	"P/enum"
	"P/model"
	"P/repository"
	"errors"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"net/url"
	"time"
)

type OrderService struct {
	OrderRepo     repository.OrderRepoInterface
	CommodityRepo repository.CommodityRepoInterface
	PaySrvI       PayServiceInterface
}

type OrderServiceInterface interface {
	StartOrderTimer(order *model.Order) error
	CreateOrder(order *model.Order) error
	PayOrder(order *model.Order) (*url.URL, error)
	FindOrder(order *model.Order) error
	FindOrderWithCom(order *model.Order) error
	DropOrder(order *model.Order) error
}

// StartOrderTimer 超时取消订单
func (srv *OrderService) StartOrderTimer(order *model.Order) error {
	timeout := time.Duration(viper.GetInt64("order.timeout"))
	timer := time.NewTimer(timeout*time.Minute - time.Since(order.CreatedAt))
	<-timer.C
	err := srv.OrderRepo.FindOrderByOrderNum(order)
	if err != nil {
		log.Println(err)
		return err
	}
	status, err := srv.PaySrvI.FindPayStatus(order)
	if err != nil {
		log.Println(err)
		return err
	}
	//只有当订单状态为未付款且支付宝查询状态为空(没有点击支付)或订单超时关闭时调用DropOrder
	if order.Status == enum.OrderStatusUnpaid && (status.Content.TradeStatus != "TRADE_FINISHED" && status.Content.TradeStatus != "TRADE_SUCCESS") {
		err = srv.DropOrder(order)
		if err != nil {
			log.Println(err)
			return err
		}
		log.Println(order.OrderNum, "超时未支付,已取消!")
	}
	return nil
}

func (srv *OrderService) CreateOrder(order *model.Order) error {
	//钱包类型订单处理方式
	if order.OrderType == enum.OrderTypeWallet {
		//补充订单信息
		order.OrderNum = uuid.NewV4().String()
		//将订单状态设为1(已下单)
		order.Status = enum.OrderStatusUnpaid
		err := srv.OrderRepo.AddOrder(order)
		//错误返回
		if err != nil {
			log.Println(err)
			return errors.New("创建订单失败,请重试")
		}
		//启动定时器
		go func() {
			_ = srv.StartOrderTimer(order)
		}()
		return nil
	}
	//商品类型订单处理方式
	//接收被锁定的商品
	var csEnd []*model.Commodity
	var csTmp []*model.Commodity
	//获取数据库连接指针并开启事务
	tx := srv.OrderRepo.GetDB().Begin()
	//锁定及查询商品
	for _, c := range order.Commodities {
		result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("name = ? AND order_num = ? And status = ?", c.Name, "", enum.CommodityActive).
			Order("RAND()").Limit(c.Amount).
			Find(&csTmp)
		//锁定错误返回
		if result.Error != nil {
			tx.Rollback()
			log.Println("rollBack")
			return errors.New("商品锁定失败")
		}
		//如果被锁定的数量与订单商品数量不符返回错误
		if result.RowsAffected != int64(c.Amount) {
			tx.Rollback()
			log.Println("rollBack")
			return errors.New("商品锁定失败!您选择的一件或多件商品已无货,请重新下单")
		}
		csEnd = append(csEnd, csTmp...)
	}
	//补充订单信息
	order.OrderNum = uuid.NewV4().String()
	//将订单状态设为1(已下单)
	order.Status = enum.OrderStatusUnpaid
	order.SellerId = csEnd[0].UserId
	//计算总金额,并修改商品所属订单号
	order.OrderAmount = 0

	for _, c := range csEnd {
		//计算订单总金额
		order.OrderAmount = order.OrderAmount + c.Price
		c.OrderNum = order.OrderNum
		//商品个数
		c.Amount = 1
	}
	//将锁定的商品切片赋给订单中商品字段
	order.Commodities = csEnd
	//传入事务指针执行添加订单操作
	err := srv.OrderRepo.AddOrder(order, tx)
	//错误返回
	if err != nil {
		tx.Rollback()
		log.Println("rollBack")
		return errors.New("创建订单失败,请重试")
	}
	//提交事务
	err = tx.Commit().Error
	//错误返回
	if err != nil {
		tx.Rollback()
		return errors.New("创建订单失败,请重试")
	}
	log.Println("commit")
	//启动定时器
	go func() {
		_ = srv.StartOrderTimer(order)
	}()
	return nil
}

func (srv *OrderService) PayOrder(order *model.Order) (*url.URL, error) {
	payUrl, err := srv.PaySrvI.Pay(order)
	if err != nil {
		return nil, err
	}
	return payUrl, nil
}

func (srv *OrderService) FindOrder(order *model.Order) error {
	return srv.OrderRepo.FindOrderByOrderNum(order)
}

func (srv *OrderService) FindOrderWithCom(order *model.Order) error {
	return srv.OrderRepo.FindOrderWithComByOrderNum(order)
}

// DropOrder 释放订单及库存
func (srv *OrderService) DropOrder(order *model.Order) error {
	//取数据库指针
	tx := srv.OrderRepo.GetDB()
	//开启事务
	tx = tx.Begin()
	//接收查询到的订单信息
	o := &model.Order{}
	var result *gorm.DB
	//锁定及查询订单
	//钱包订单处理方式
	if order.OrderType == enum.OrderTypeWallet {
		result = tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("order_num=?", order.OrderNum).
			Find(o)
	} else {
		//商品订单处理方式
		result = tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Preload("Commodities").
			Where("order_num=?", order.OrderNum).
			Find(o)
	}
	if result.Error != nil {
		tx.Rollback()
		log.Println("rollBack")
		return errors.New("取消订单失败!请重试")
	}
	//如果查到的订单记录为空,返回错误
	if result.RowsAffected == 0 {
		tx.Rollback()
		log.Println("rollBack")
		return errors.New("订单不存在!请重试")
	}
	if o.Status == enum.OrderStatusPaid {
		tx.Rollback()
		log.Println("rollBack")
		return errors.New("订单已付款,如要取消,请申请退款")
	}
	if o.Status == enum.OrderStatusFinish {
		tx.Rollback()
		log.Println("rollBack")
		return errors.New("订单已完成,无法取消")
	}
	if o.Status == enum.OrderStatusCancelled {
		tx.Rollback()
		log.Println("rollBack")
		return errors.New("订单处于未活跃状态")
	}
	tradeClose, err := srv.PaySrvI.ClosePay(o)
	log.Println(tradeClose.Content)
	if err != nil {
		tx.Rollback()
		log.Println("rollBack")
		return errors.New("支付订单取消失败")
	}
	//将订单状态设为0(已取消)
	err = srv.OrderRepo.UpdateOrderStatus(o, enum.OrderStatusCancelled, tx)
	if err != nil {
		tx.Rollback()
		log.Println("rollBack")
		return errors.New("取消订单失败!请重试")
	}
	//对于商品订单:传入商品,空字符串事务指针将商品order_id和user_id字段更新为空(释放商品)
	if o.OrderType == enum.OrderTypeCommodity {
		err = srv.CommodityRepo.UpdateCommoditiesOrderNumUserId(o.Commodities, &model.Order{}, tx)
		if err != nil {
			tx.Rollback()
			log.Println("rollBack")
			return errors.New("取消订单失败!请重试")
		}
	}
	//提交事务
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		log.Println("rollBack")
		return errors.New("取消订单失败!请重试")
	}
	log.Println("commit")
	return nil
}
