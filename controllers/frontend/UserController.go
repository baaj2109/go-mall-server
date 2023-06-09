package frontend

import (
	"LeastMall/models"
	"math"
	"strconv"
	"time"
)

type UserController struct {
	BaseController
}

func (c *UserController) Get() {
	c.BaseInit()
	user := models.User{}
	models.Cookie.Get(c.Ctx, "userinfo", &user)
	c.Data["user"] = user
	time := time.Now().Hour()
	if time >= 12 && time <= 18 {
		c.Data["Hello"] = "下午好"
	} else if time >= 6 && time < 12 {
		c.Data["Hello"] = "上午好"
	} else {
		c.Data["Hello"] = "晚上好"
	}
	order := []models.Order{}
	models.DB.Where("uid=?", user.Id).Find(&order)
	var wait_pay int
	var wait_rec int
	for i := 0; i < len(order); i++ {
		if order[i].PayStatus == 0 {
			wait_pay += 1
		}
		if order[i].OrderStatus >= 2 && order[i].OrderStatus < 4 {
			wait_rec += 1
		}
	}
	c.Data["wait_pay"] = wait_pay
	c.Data["wait_rec"] = wait_rec
	c.TplName = "frontend/user/welcome.html"
}

// 展示訂單清單
/*
	1. 獲取使用者資訊
	2. 獲取使用者訂單資訊並分頁
	3. 獲取搜尋關鍵字
	4. 獲取篩選條件
	5. 計算總數量
*/
func (c *UserController) OrderList() {
	c.BaseInit()
	/// 獲取當前使用者
	user := models.User{}
	models.Cookie.Get(c.Ctx, "userinfo", &user)
	/// 獲取使用者的訂單訊息
	page, _ := c.GetInt("page")
	if page == 0 {
		page = 1
	}
	pageSize := 2
	/// 獲取搜尋關鍵字
	where := "uid=?"
	keywords := c.GetString("keywords")
	if keywords != "" {
		orderitem := []models.OrderItem{}
		models.DB.Where("product_title like ?", "%"+keywords+"%").Find(&orderitem)
		var str string
		for i := 0; i < len(orderitem); i++ {
			if i == 0 {
				str += strconv.Itoa(orderitem[i].OrderId)
			} else {
				str += "," + strconv.Itoa(orderitem[i].OrderId)
			}
		}
		where += " AND id in (" + str + ")"
	}
	/// 獲取篩選條件
	orderStatus, err := c.GetInt("order_status")
	if err == nil {
		where += " AND order_status=" + strconv.Itoa(orderStatus)
		c.Data["orderStatus"] = orderStatus
	} else {
		c.Data["orderStatus"] = "nil"
	}

	/// 訂單總數量
	var count int
	models.DB.Where(where, user.Id).Table("order").Count(&count)
	order := []models.Order{}
	models.DB.Where(where, user.Id).Offset((page - 1) * pageSize).Limit(pageSize).Preload("OrderItem").Order("add_time desc").Find(&order)

	c.Data["order"] = order
	c.Data["totalPages"] = math.Ceil(float64(count) / float64(pageSize))
	c.Data["page"] = page
	c.Data["keywords"] = keywords
	c.TplName = "frontend/user/order.html"
}

// 訂單詳情
// 根據訂單編號展示訂單詳情
func (c *UserController) OrderInfo() {
	c.BaseInit()
	id, _ := c.GetInt("id")
	user := models.User{}
	models.Cookie.Get(c.Ctx, "userinfo", &user)
	order := models.Order{}
	models.DB.Where("id=? AND uid=?", id, user.Id).Preload("OrderItem").Find(&order)
	c.Data["order"] = order
	if order.OrderId == "" {
		c.Redirect("/", 302)
	}
	c.TplName = "frontend/user/order_info.html"
}
