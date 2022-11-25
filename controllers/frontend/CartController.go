package frontend

import (
	"LeastMall/models"
	"strconv"
)

type CartController struct {
	BaseController
}

func (c *CartController) Get() {
	c.BaseInit()
	cartList := []models.Cart{}
	models.Cookie.Get(c.Ctx, "cartList", &cartList)

	var allPrice float64
	/// 計算總價
	for i := 0; i < len(cartList); i++ {
		if cartList[i].Checked {
			allPrice += cartList[i].Price * float64(cartList[i].Num)
		}
	}
	c.Data["cartList"] = cartList
	c.Data["allPrice"] = allPrice
	c.TplName = "frontend/cart/cart.html"
}

/*
在使用者點擊 加入購物車 後，判斷購物車中有沒有資料
如果有當前商品資料，則會將購物車商品數量＋１
如果沒有資料，則把當前資料寫入Ｃｏｏｋｉｅｓ
*/
func (c *CartController) AddCart() {
	c.BaseInit()
	colorId, colorErr := c.GetInt("color_id")
	productId, productErr := c.GetInt("product_id")

	product := models.Product{}
	productColor := models.ProductColor{}

	canFindProduct := models.DB.Where("id=?", productId).Find(&product).Error
	canFindProductColor := models.DB.Where("id=?", colorId).Find(&productColor).Error

	if colorErr != nil ||
		productErr != nil ||
		canFindProduct != nil ||
		canFindProductColor != nil {
		c.Ctx.Redirect(302, "/item_"+strconv.Itoa(product.Id)+".html")
	}

	/// 新增到購物車的商品資訊
	currentProduct := models.Cart{
		Id:           productId,
		Title:        product.Title,
		Price:        product.Price,
		GoodsVersion: product.ProductVersion,
		Num:          1,
		ProductGift:  product.ProductGift,
		ProductColor: productColor.ColorName,
		ProductImg:   product.ProductImg,
		ProductAttr:  "",
		Checked:      true,
	}
	/// 判斷購物車有沒有資料 (cookie)
	cartList := []models.Cart{}
	models.Cookie.Get(c.Ctx, "carList", &cartList)
	/// 購物車有商品
	if len(cartList) > 0 {
		/// 判斷是否有相同商品
		if models.IsCartHasData(cartList, currentProduct) {
			for i := 0; i < len(cartList); i++ {
				if cartList[i].Id == currentProduct.Id &&
					cartList[i].ProductColor == currentProduct.ProductColor &&
					cartList[i].ProductAttr == currentProduct.ProductAttr {
					cartList[i].Num = cartList[i].Num + 1
				}
			}
		} else {
			cartList = append(cartList, currentProduct)
		}
		// models.Cookie.Set(c.Ctx, "carList", cartList)
	} else { /// 沒有商品在購物車中
		cartList = append(cartList, currentProduct)
	}
	models.Cookie.Set(c.Ctx, "cartList", cartList)
	c.Data["product"] = product
	c.TplName = "frontend/cart/addcart_success.html"
}

func (c *CartController) IncCart() {
	var flag bool
	var allPrice float64
	var currentAllPrice float64
	var num int

	productId, _ := c.GetInt("product_id")
	productColor := c.GetString("product_color")
	productAttr := ""

	cartList := []models.Cart{}
	models.Cookie.Get(c.Ctx, "cartList", &cartList)
	for i := 0; i < len(cartList); i++ {
		if cartList[i].Id == productId &&
			cartList[i].ProductColor == productColor &&
			cartList[i].ProductAttr == productAttr {
			cartList[i].Num = cartList[i].Num + 1
			flag = true
			num = cartList[i].Num
			currentAllPrice = cartList[i].Price * float64(cartList[i].Num)
		}
		if cartList[i].Checked {
			allPrice += cartList[i].Price * float64(cartList[i].Num)
		}
	}

	if flag {
		models.Cookie.Set(c.Ctx, "cartList", cartList)
		c.Data["json"] = map[string]interface{}{
			"success":         true,
			"message":         "修改數量成功",
			"allPrice":        allPrice,
			"currentAllPrice": currentAllPrice,
			"num":             num,
		}

	} else {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "傳入参數錯誤",
		}
	}
	c.ServeJSON()
}

func (c *CartController) DecCart() {
	var flag bool
	var allPrice float64
	var currentAllPrice float64
	var num int
	productId, _ := c.GetInt("product_id")
	productColor := c.GetString("product_color")
	productAttr := ""
	cartList := []models.Cart{}
	models.Cookie.Get(c.Ctx, "cartList", &cartList)
	for i := 0; i < len(cartList); i++ {
		if cartList[i].Id == productId &&
			cartList[i].ProductColor == productColor &&
			cartList[i].ProductAttr == productAttr {
			/// main logic
			cartList[i].Num = cartList[i].Num - 1
			flag = true
			num = cartList[i].Num
			currentAllPrice = cartList[i].Price * float64(cartList[i].Num)
		}
		if cartList[i].Checked {
			allPrice += cartList[i].Price * float64(cartList[i].Num)
		}
	}
	if flag {
		models.Cookie.Set(c.Ctx, "cartList", cartList)
		c.Data["json"] = map[string]interface{}{
			"success":         true,
			"message":         "修改數量成功",
			"allPrice":        allPrice,
			"currentAllPrice": currentAllPrice,
			"num":             num,
		}

	} else {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "傳入参數錯誤",
		}
	}
	c.ServeJSON()
}

func (c *CartController) DelCart() {
	productId, _ := c.GetInt("product_id")
	productColor := c.GetString("product_color")
	productAttr := ""

	cartList := []models.Cart{}
	models.Cookie.Get(c.Ctx, "cartList", &cartList)
	for i := 0; i < len(cartList); i++ {
		if cartList[i].Id == productId &&
			cartList[i].ProductColor == productColor &&
			cartList[i].ProductAttr == productAttr {
			cartList = append(cartList[:i], cartList[(i+1):]...)
		}
	}
	models.Cookie.Set(c.Ctx, "cartList", cartList)

	c.Redirect("/cart", 302)
}

func (c *CartController) ChangeAllCart() {
	flag, _ := c.GetInt("flag")
	var allPrice float64
	cartList := []models.Cart{}
	models.Cookie.Get(c.Ctx, "cartList", &cartList)
	for i := 0; i < len(cartList); i++ {
		if flag == 1 {
			cartList[i].Checked = true
		} else {
			cartList[i].Checked = false
		}
		/// calaulate total price
		if cartList[i].Checked {
			allPrice += cartList[i].Price * float64(cartList[i].Num)
		}
	}
	models.Cookie.Set(c.Ctx, "cartList", cartList)

	c.Data["json"] = map[string]interface{}{
		"success":  true,
		"allPrice": allPrice,
	}
	c.ServeJSON()
}

func (c *CartController) ChangeOneCart() {
	var flag bool
	var allPrice float64
	productId, _ := c.GetInt("product_id")
	productColor := c.GetString("product_color")
	productAttr := ""

	cartList := []models.Cart{}
	models.Cookie.Get(c.Ctx, "cartList", &cartList)

	for i := 0; i < len(cartList); i++ {
		if cartList[i].Id == productId &&
			cartList[i].ProductColor == productColor &&
			cartList[i].ProductAttr == productAttr {
			cartList[i].Checked = !cartList[i].Checked
			flag = true
		}
		if cartList[i].Checked {
			allPrice += cartList[i].Price * float64(cartList[i].Num)
		}
	}

	if flag {
		models.Cookie.Set(c.Ctx, "cartList", cartList)
		c.Data["json"] = map[string]interface{}{
			"success":  true,
			"message":  "修改狀態成功",
			"allPrice": allPrice,
		}

	} else {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "傳入参數錯誤",
		}
	}
	c.ServeJSON()
}
