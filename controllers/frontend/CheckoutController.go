package frontend

import (
	"LeastMall/common"
	"LeastMall/models"
	"fmt"
	"strconv"
)

type CheckoutController struct {
	BaseController
}

func (c *CheckoutController) Checkout() {

	c.BaseInit()
	/// get product which want to buy
	cartList := []models.Cart{}  /// all product in cart list
	orderList := []models.Cart{} /// product want to but in cart list
	models.Cookie.Get(c.Ctx, "cartList", &cartList)

	var allPrice float64
	/// measure current order total price
	for i := 0; i < len(cartList); i++ {
		if cartList[i].Checked {
			allPrice += cartList[i].Price * float64(cartList[i].Num)
			orderList = append(orderList, cartList[i])
		}
	}
	/// no order return to index
	if len(orderList) == 0 {
		c.Redirect("/", 302)
		return
	}

	c.Data["orderList"] = orderList
	c.Data["allPrice"] = allPrice

	/// get address
	user := models.User{}
	models.Cookie.Get(c.Ctx, "userinfo", &user)
	addressList := []models.Address{}
	models.DB.Where("uid=?", user.Id).Order("default_address desc").Find(&addressList)
	c.Data["addressList"] = addressList

	/// create sign to prevent duplicate order
	orderSign := common.Md5(common.GetRandomNum())
	c.SetSession("orderSign", orderSign)
	c.Data["orderSign"] = orderSign

	c.TplName = "frontend/buy/checkout.html"
}

/*
提交訂單
1. 收貨地址資訊
2. 購買商品資訊
3. 訂單訊息放進訂單表 商品訊息放進商品表
4. 刪除購物車選中商品
*/
func (c *CheckoutController) GoOrder() {
	/// check order sign
	orderSign := c.GetString("orderSign")
	sessionOrderSign := c.GetSession("orderSign")
	if sessionOrderSign != orderSign {
		c.Redirect("/", 302)
		return
	}
	c.DelSession("orderSign")

	/// get address
	user := models.User{}
	models.Cookie.Get(c.Ctx, "userinfo", &user)

	addressResult := []models.Address{}
	models.DB.Where("uid=? AND default_address=1", user.Id).Find(&addressResult)

	if len(addressResult) > 0 {

		/// get product information , order list is product wanna buy
		cartList := []models.Cart{}
		orderList := []models.Cart{}
		models.Cookie.Get(c.Ctx, "cartList", &cartList)
		var allPrice float64
		for i := 0; i < len(cartList); i++ {
			if cartList[i].Checked {
				allPrice += cartList[i].Price * float64(cartList[i].Num)
				orderList = append(orderList, cartList[i])
			}
		}

		/// set order into order table in db, product into product table in db
		order := models.Order{
			OrderId:     common.GenerateOrderId(),
			Uid:         user.Id,
			AllPrice:    allPrice,
			Phone:       addressResult[0].Phone,
			Name:        addressResult[0].Name,
			Address:     addressResult[0].Address,
			Zipcode:     addressResult[0].Zipcode,
			PayStatus:   0,
			PayType:     0,
			OrderStatus: 0,
			AddTime:     int(common.GetUnix()),
		}
		err := models.DB.Create(&order).Error
		if err == nil {
			for i := 0; i < len(orderList); i++ {
				orderItem := models.OrderItem{
					OrderId:      order.Id,
					Uid:          user.Id,
					ProductTitle: orderList[i].Title,
					ProductId:    orderList[i].Id,
					ProductImg:   orderList[i].ProductImg,
					ProductPrice: orderList[i].Price,
					ProductNum:   orderList[i].Num,
					GoodsVersion: orderList[i].GoodsVersion,
					GoodsColor:   orderList[i].ProductColor,
					AddTime:      int(common.GetUnix()),
				}
				err := models.DB.Create(&orderItem).Error
				if err != nil {
					fmt.Println(err)
				}
			}
			/// delete product in cart
			noSelectedCartList := []models.Cart{}
			for i := 0; i < len(cartList); i++ {
				if !cartList[i].Checked {
					noSelectedCartList = append(noSelectedCartList, cartList[i])
				}
			}
			models.Cookie.Set(c.Ctx, "cartList", noSelectedCartList)
			c.Redirect("/buy/confirm?id="+strconv.Itoa(order.Id), 302)

		} else {
			c.Redirect("/", 302)
		}
	} else {
		c.Redirect("/", 302)
	}
}

// / 下墊單
func (c *CheckoutController) Confirm() {
	c.BaseInit()
	id, err := c.GetInt("id")
	if err != nil {
		c.Redirect("/", 302)
		return
	}
	/// get user
	user := models.User{}
	models.Cookie.Get(c.Ctx, "userinfo", &user)

	/// get order from db
	order := models.Order{}
	models.DB.Where("id=?", id).Find(&order)
	c.Data["order"] = order
	/// check current order id is from current user
	if user.Id != order.Uid {
		c.Redirect("/", 302)
		return
	}

	/// get product infromation under current order
	orderItem := []models.OrderItem{}
	models.DB.Where("order_id=?", id).Find(&orderItem)
	c.Data["orderItem"] = orderItem
	c.TplName = "frontend/buy/confirm.html"
}

// / order status
func (c *CheckoutController) OrderPayStatus() {
	/// get order id
	id, err := c.GetInt("id")
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "傳入參數錯誤",
		}
		c.ServeJSON()
		return
	}
	/// get order from order id
	order := models.Order{}
	models.DB.Where("id=?", id).Find(&order)

	user := models.User{}
	models.Cookie.Get(c.Ctx, "userinfo", &user)

	/// check order user id equal user id
	if user.Id != order.Uid {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "傳入參數錯誤",
		}
		c.ServeJSON()
		return
	}

	/// check order payment status
	if order.PayStatus == 1 && order.OrderStatus == 1 {
		c.Data["json"] = map[string]interface{}{
			"success": true,
			"message": "已支付",
		}
		c.ServeJSON()

	} else {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "未支付",
		}
		c.ServeJSON()
	}
}
