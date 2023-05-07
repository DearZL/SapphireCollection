package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"
	"log"
)

type PayHandler struct {
}

func (h *PayHandler) Pay(c *gin.Context) {
	appID := "2021000122682269"
	privateKey := "MIIEpQIBAAKCAQEAy2DolUYc9Oc29mT86MgQyjQoq5pmYOOXtHCj2cSLigA9pFz9qqP2tBIAUupz9V3G7r1KzDeZM8gZyvvamxkTHqT+Kvaoo1phS7inXoam405kOU4UpOM5eUZQnsV91giRVHha5aO6no3Bn6hxxbOPz5Dtc7MolOBKPIiQvLPa0YzB6upVe3978pezVbPn+3ni5clijs2niJHPP22WpxLdzfifDGhHy0VbxY2EUExe23ny6jxt9hnzmDhsdKJ+lCRlQGBFKlqM3WGTRMSrEJRl4ndWMcYObk6JDhYxEFKbgJA66zkpwQsFDdm55yXvP+52pTqSizt7iCtkCTS/P+EOcwIDAQABAoIBAQC4hFSf0fvFidzI0TjP7WumOIpJnoySDQr/L07I7VP4QV2ruJ6AacATAV3/3CyWiZ1Jzr2E6FB7tWkJS1S7cJVzMRhUXHMFuaMacw6OaTYSdnXhs+Bw9KKZT90nH2Cahi1saMF3JQPUhCIOO2H1j4LDO+bjGMGRyKgxoWlHexnlEOQ5WqDd3nAXSvlp5xwIA3Qo7Q/FRK7g9PGiudsUarqisy+Oxyq+LuuWAry48z0SOxdc/4mNhmfcQypNCXagzzUCeCu8+8iacyQeJYX2qtjATz9Tx3hX3+2A17qYa40diSOH+uD1nLRVDGLs9hbksaAuIdAexXP5FURhHNdcLlM5AoGBAOw5Zu8tAjhD2nr/9bvbQ5NVjQ070x9Qx5Y/3rIbXEDzjq67TveCTrRetKvlnRFSHLdu3Utm4Mm7mQSSkQJZjFlQV4GkprYIAEsDRXO/zOk0bezfMDfHxssQn6N65PO/n+DUjjypZFB6qhaXe7AyhXtsPcovp+Pn0ALYIbYiHYpNAoGBANxnlIAFaTSN41Mr6+E0nAXcRiEstaKSyGqFe4k/KIS6SXucxfhUvHQM7PCzdwRjfPT5kN6fdrndVT3IfHGSVg++vWvPah0h7I/1TMorrOUASe9dgCa1pmMRWnvTxWucJKD8i2W2e2n3WaRzXbZrLVRQE+o5+nSvX5sMFcAWI9u/AoGBAIqExoVt4SVZNJ53xYMY+jFFM2cVM6HjXoYOgeny/U/hAkQX9iBROxGtj0hVZpsniUtPKVjzxNDGvt3djEbSd+hPomCVSmTnoDRcgLd1OxVs9yC1Z7Lt5PZikxnsEKGWNoxCV/3eXsKKi36f6ZnSpk9Pk5QiCdMstd9VGb+RlbzpAoGAI/6qiunXT6TofjnLEQF1haN+tIZHt6A/KN5Z2YU+CccenxhwYGj+SfmebITyp/3Td3KWjTT/v8T82dU3NZkPgwzEhKngC5fxuWT3QIE3gEK20Ge1uRyrarx3yYdBU5yxgrUb0uWlbB3gPvI0WMlSItXdGsCTPaEyfPDRUIiNHKMCgYEAxn3sDXgfOE6FKLNEe7eeszf/PUFKviYOPBU/cw90yBD5ppROSGWnh5gsu8LQqhnIu3EjoycOcyY8y5OKPZhLE5wmUeslNplRUpoF6YE6JtR+PFPk4axGtaIX4mYiUYzn7wyL4lMMEwRK4eto0FQWGS9qdJrAM+C0XSeXGvqqliE="
	aliPublicKey := "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAhlncP9uSU/fXcBPlIfuZFF2pd/QIu0Hp2FQjjSsNYH/6KRaxX74G7ACJEiv+ZZCI8jSMvIDo065XX9cIF8TEsa/2h0vFIAhzx/cSfxDL+PmKuRxvrGpJrfZEuSz5hZaN/XC1IdhA6ni10SdYy8AtUwg6ah1LtVZtqoU/HHf/tPngUqu/Y0Nje7NQQvCHr7bgU8xhYNjpo4dT9lELKkHZ88DaDL7nj6KjcwZPrPahN0CCw1At6ZonC6fPuTi7dwcQpMUWigLOyogtnkSyEyFVKKzOV/2UOTRUJcSRPxeweUkN8ghOFGQPoZL8damg5yLKmSQD95GAmXDWAV0QXp5JmQIDAQAB"

	var client, err = alipay.New(appID, privateKey, false)
	if err != nil {
		log.Println(err.Error())
	}
	err = client.LoadAliPayPublicKey(aliPublicKey)
	if err != nil {
		log.Println(err.Error())
	}

	var p = alipay.TradePagePay{} // page支付方式使用
	p.NotifyURL = ""              // 支付结果回调的url，注意内网穿透问题
	p.ReturnURL = "www.baidu.com" // 支付成功后倒计时结束跳转的页面
	p.Subject = "标题"
	p.OutTradeNo = "sn12389342479" //传递一个唯一单号
	p.TotalAmount = "12.00"
	p.ProductCode = "FAST_INSTANT_TRADE_PAY" // page支付必须使用这个配置

	url, err := client.TradePagePay(p)
	if err != nil {
		panic(err)
	}

	c.String(200, url.String())
}
