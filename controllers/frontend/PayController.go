package frontend

import (
	"LeastMall/models"

	"github.com/astaxie/beego"
	"github.com/processout/applepay"
)

type PayController struct {
	BaseController
}

var (
	ap *applepay.Merchant
)

func init() {
	var err error

	applepay.AppleRootCertificatePath = "AppleRootCA-G3.crt"
	ap, err = applepay.New(
		"merchant.com.LeastMall.test",
		applepay.MerchantDisplayName("LeastMall"),
		applepay.MerchantDomainName("applepay.LeastMall.com"),
		applepay.MerchantCertificateLocation(
			"certs/cert-merchant.crt",
			"certs/cert-merchant-key.pem",
		),
		applepay.ProcessingCertificateLocation(
			"certs/cert-processing.crt",
			"certs/cert-processing-key.pem",
		),
	)
	if err != nil {
		panic(err)
	}
	beego.Info("Apple Pay test app starting")
}

func (c *PayController) Applepay() {
	orderId, err := c.GetInt("id")
	if err != nil {
		c.Redirect(c.Ctx.Request.Referer(), 302)
	}

	orderItem := []models.OrderItem{}
	models.DB.Where("order_id=?", orderId).Find(orderItem)

	var totalAmount float64
	for i := 0; i < len(orderItem); i++ {
		totalAmount = totalAmount + orderItem[i].ProductPrice
	}

	payLoad, err := ap.Session(c.Ctx.Request.RequestURI)
	if err != nil {
		beego.Info(err)
		// c.Status(http.StatusInternalServerError)
		return
	}
	str := string(payLoad)
	c.Data["userinfo"] = str
}
