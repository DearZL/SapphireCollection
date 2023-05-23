package service

import (
	"P/enum"
	"P/model"
	"P/repository"
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"gorm.io/gorm/clause"
	"log"
	"net/url"
	"time"
)

type OrderService struct {
	OrderRepo     repository.OrderRepoInterface
	CommodityRepo repository.CommodityRepoInterface
	PaySrv        *PayService
}

type OrderServiceInterface interface {
	StartOrderTimer(order *model.Order) error
	CreateOrder(order *model.Order, com *model.Commodity) (*model.Order, error)
	PayOrder(order *model.Order) (*url.URL, error)
	FindOrder(order *model.Order) error
	FindOrderWithCom(order *model.Order) error
	DropOrder(order *model.Order) error
}

// StartOrderTimer 超时取消订单
func (srv *OrderService) StartOrderTimer(order *model.Order) error {
	timeout := time.Duration(viper.GetInt64("order.timeout"))
	timer := time.NewTimer(timeout*time.Minute - time.Since(order.CreatedAt) + 4*time.Second)
	<-timer.C
	err := srv.FindOrderWithCom(order)
	if err != nil {
		log.Println(err)
		return err
	}
	status, err := srv.PaySrv.FindPayStatus(order)
	if err != nil {
		log.Println(err)
		return err
	}
	//只有当订单状态为未付款且支付宝查询状态为空(没有点击支付)或订单超时关闭时调用DropOrder
	if order.Status == enum.OrderStatusUnpaid && (status.Content.TradeStatus == "TRADE_CLOSED" || status.Content.TradeStatus == "") {
		err = srv.DropOrder(order)
		if err != nil {
			log.Println(err)
			return err
		}
		log.Println(order.OrderNum, "超时未支付,已取消!")
	}
	return nil
}

func (srv *OrderService) CreateOrder(order *model.Order, com *model.Commodity) (*model.Order, error) {
	fmt.Println(srv.PaySrv)
	//接收被锁定的商品
	var csEnd []*model.Commodity
	//获取数据库连接指针并开启事务
	tx := srv.OrderRepo.GetDB().Begin()
	//锁定及查询商品
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("name = ? AND order_num = ?", com.Name, "").Order("RAND()").Limit(order.CommodityAmount).Find(&csEnd).Error
	//锁定错误返回
	if err != nil {
		tx.Rollback()
		return nil, errors.New("商品锁定失败")
	}
	//如果被锁定的数量与请求商品数量不符返回错误
	if len(csEnd) != order.CommodityAmount {
		tx.Rollback()
		return nil, errors.New("商品锁定失败!您选择的一件或多件商品已无货,请重新下单")
	}
	order.OrderNum = uuid.NewV4().String()
	//将订单状态设为1(已下单)
	order.Status = enum.OrderStatusUnpaid
	//计算总金额,并修改商品所属订单号
	order.OrderAmount = 0
	for _, c := range csEnd {
		order.OrderAmount = order.OrderAmount + c.Price
		c.OrderNum = order.OrderNum
	}
	order.Commodities = csEnd
	//传入事务指针执行添加订单操作
	err = srv.OrderRepo.AddOrder(order, tx)
	//错误返回
	if err != nil {
		tx.Rollback()
		log.Println("rollBack")
		return nil, errors.New("创建订单失败,请重试")
	}
	//提交事务
	err = tx.Commit().Error
	//错误返回
	if err != nil {
		tx.Rollback()
		return nil, errors.New("创建订单失败,请重试")
	}
	log.Println("commit")
	//启动定时器
	go func() {
		err := srv.StartOrderTimer(order)
		if err != nil {

		}
	}()
	return order, nil
}

func (srv *OrderService) PayOrder(order *model.Order) (*url.URL, error) {
	payUrl, err := srv.PaySrv.Pay(order)
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
	//锁定及查询订单
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Preload("Commodities").
		Where("order_num=?", order.OrderNum).
		Find(o).Error
	if err != nil {
		tx.Rollback()
		log.Println("rollBack")
		return errors.New("取消订单失败!请重试")
	}
	//如果查到的订单为默认值即空结果,返回错误
	if o.ID == 0 {
		tx.Rollback()
		log.Println("rollBack")
		return errors.New("订单不存在!请重试")
	}
	if o.Status == enum.OrderStatusPaid {
		tx.Rollback()
		log.Println("rollBack")
		return errors.New("订单已付款,如要取消,请申请退款")
	}
	if o.Status == enum.OrderStatusCancelled {
		tx.Rollback()
		log.Println("rollBack")
		return errors.New("订单处于未活跃状态")
	}
	refund, err := srv.PaySrv.PayTimeOut(o)
	if err != nil {
		tx.Rollback()
		log.Println("rollBack")
		log.Println(refund.Content)
		return errors.New("支付订单取消失败")
	}
	log.Println(refund.Content)
	//将订单状态设为0(已取消)
	err = srv.OrderRepo.UpdateOrderStatus(o, enum.OrderStatusCancelled, tx)
	if err != nil {
		tx.Rollback()
		log.Println("rollBack")
		return errors.New("取消订单失败!请重试")
	}
	//传入商品,空字符串事务指针将商品order_id和user_id字段更新为空(释放商品)
	err = srv.CommodityRepo.UpdateCommodities(o.Commodities, &model.Order{}, tx)
	if err != nil {
		tx.Rollback()
		log.Println("rollBack")
		return errors.New("取消订单失败!请重试")
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
