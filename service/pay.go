package service

import (
	"P/enum"
	"P/model"
	"errors"
	"github.com/smartwalle/alipay/v3"
	"github.com/spf13/viper"
	"log"
	"net/url"
	"strconv"
	"time"
)

type PayService struct{}

type PayServiceInterface interface {
	Pay(order *model.Order) (*url.URL, error)
	ClosePay(order *model.Order) (*alipay.TradeCloseRsp, error)
	FindPayStatus(order *model.Order) (*alipay.TradeQueryRsp, error)
}

func (srv *PayService) Pay(order *model.Order) (*url.URL, error) {
	if order.OrderNum == "" {
		return nil, errors.New("参数错误")
	}
	if order.Status == enum.OrderStatusCancelled {
		return nil, errors.New("支付失败,订单已取消")
	}
	if order.Status == enum.OrderStatusPaid {
		return nil, errors.New("支付失败,订单已支付")
	}
	if time.Now().After(order.CreatedAt.Add(time.Duration(viper.GetInt64("order.timeout")) * time.Minute)) {
		return nil, errors.New("订单已过期！")
	}
	appID := viper.GetString("alipay.appId")
	privateKey := viper.GetString("alipay.privateKey")
	aliPublicKey := viper.GetString("alipay.aliPublicKey")
	client, err := alipay.New(appID, privateKey, false)
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("生成支付链接失败")
	}
	err = client.LoadAliPayPublicKey(aliPublicKey)
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("生成支付链接失败")
	}
	var p = alipay.TradePagePay{}                         // page支付方式使用
	p.NotifyURL = "/api/pay/success"                      // 支付结果回调的url，注意内网穿透问题
	p.ReturnURL = "http://localhost:9090/api/pay/success" // 支付成功后倒计时结束跳转的页面
	p.Subject = order.BuyerId + " 的订单"
	p.OutTradeNo = order.OrderNum //传递一个唯一单号
	p.TotalAmount = strconv.FormatFloat(float64(order.OrderAmount), 'f', 2, 64)
	p.ProductCode = "FAST_INSTANT_TRADE_PAY" // page支付必须使用这个配置
	p.TimeExpire = order.CreatedAt.Add(time.Duration(viper.GetInt64("order.timeout")) * time.Minute).Format("2006-01-02 15:04:05")
	payUrl, err := client.TradePagePay(p)
	if err != nil {
		log.Println(err)
		return nil, errors.New("生成支付链接失败")
	}
	return payUrl, nil
}

func (srv *PayService) ClosePay(order *model.Order) (*alipay.TradeCloseRsp, error) {
	if order.OrderNum == "" {
		return nil, errors.New("参数错误")
	}
	if order.Status == enum.OrderStatusPaid {
		return nil, errors.New("支付订单取消失败,订单已支付")
	}
	appID := viper.GetString("alipay.appId")
	privateKey := viper.GetString("alipay.privateKey")
	aliPublicKey := viper.GetString("alipay.aliPublicKey")
	var p = alipay.TradeClose{}
	p.OutTradeNo = order.OrderNum
	client, err := alipay.New(appID, privateKey, false)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	err = client.LoadAliPayPublicKey(aliPublicKey)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	tradeClose, err := client.TradeClose(p)
	if err != nil {
		return nil, err
	}
	return tradeClose, nil
}

func (srv *PayService) FindPayStatus(order *model.Order) (*alipay.TradeQueryRsp, error) {
	appID := viper.GetString("alipay.appId")
	privateKey := viper.GetString("alipay.privateKey")
	aliPublicKey := viper.GetString("alipay.aliPublicKey")
	var param alipay.TradeQuery
	param.OutTradeNo = order.OrderNum
	client, err := alipay.New(appID, privateKey, false)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	err = client.LoadAliPayPublicKey(aliPublicKey)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	result, err := client.TradeQuery(param)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return result, nil
}
