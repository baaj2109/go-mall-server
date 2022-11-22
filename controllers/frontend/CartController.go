package frontend

import (
	"LeastMall/models"
	"strconv"
)

/*
	在使用者點擊 加入購物車 後，判斷購物車中有沒有資料
	如果有當前商品資料，則會將購物車商品數量＋１
	如果沒有資料，則把當前資料寫入Ｃｏｏｋｉｅｓ
*/

type CartController struct {
	BaseController
}

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
	models.Cookie.Set(c.Ctx, "carList", cartList)
	c.Data["product"] = product
	c.TplName = "frontend/cart/addcart_success.html"
}
