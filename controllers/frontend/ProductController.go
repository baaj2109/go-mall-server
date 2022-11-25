package frontend

import (
	"LeastMall/common"
	"LeastMall/models"
	"math"
	"strconv"
	"strings"
)

type ProductController struct {
	BaseController
}

func (c *ProductController) CategoryList() {
	/// 呼叫公共方法
	c.BaseInit()

	id := c.Ctx.Input.Param(":id")
	cateId, _ := strconv.Atoi(id)
	curretProductCate := models.ProductCate{}
	subProductCate := []models.ProductCate{}
	models.DB.Where("id=?", cateId).Find(&curretProductCate)

	/// 當前頁面
	page, _ := c.GetInt("page")
	if page == 0 {
		page = 1
	}
	/// 每一頁顯示數量
	pageSize := 5

	var tempSlice []int
	if curretProductCate.Pid == 0 { /// 主要分類
		/// 次要分類
		models.DB.Where("pid=?", curretProductCate.Id).Find(&subProductCate)
		for i := 0; i < len(subProductCate); i++ {
			tempSlice = append(tempSlice, subProductCate[i].Id)
		}
	} else {
		/// 獲取次要分類對應的同級分類
		models.DB.Where("pid=?", curretProductCate.Pid).Find(&subProductCate)
	}
	tempSlice = append(tempSlice, cateId)
	where := "cate_id in (?)"
	product := []models.Product{}
	models.DB.Where(where, tempSlice).Select("id,title,price,product_img,sub_title").Offset((page - 1) * pageSize).Limit(pageSize).Order("sort desc").Find(&product)
	/// 查詢 product 裡的數量
	var count int
	models.DB.Where(where, tempSlice).Table("product").Count(&count)

	c.Data["productList"] = product
	c.Data["subProductCate"] = subProductCate
	c.Data["curretProductCate"] = curretProductCate
	c.Data["totalPages"] = math.Ceil(float64(count) / float64(pageSize))
	c.Data["page"] = page

	/// 指定分類模板
	tpl := curretProductCate.Template
	if tpl == "" {
		tpl = "frontend/product/list.html"
	}

	c.TplName = tpl
}

func (c *ProductController) Collect() {
	productId, err := c.GetInt("product_id")
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"msg":     "參數錯誤",
		}
		c.ServeJSON()
		return
	}
	user := models.User{}
	ok := models.Cookie.Get(c.Ctx, "userinfo", &user)
	if ok != true {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"msg":     "請先登入",
		}
		c.ServeJSON()
		return
	}
	isExist := models.DB.First(&user)
	if isExist.RowsAffected == 0 {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"msg":     "用戶錯誤",
		}
		c.ServeJSON()
		return
	}

	goodCollect := models.ProductCollect{}
	isExist = models.DB.Where("user_id=? AND product_id=?", user.Id, productId).First(&goodCollect)
	if isExist.RowsAffected == 0 {
		goodCollect.UserId = user.Id
		goodCollect.ProductId = productId
		goodCollect.AddTime = common.FormatDay()
		models.DB.Create(&goodCollect)
		c.Data["json"] = map[string]interface{}{
			"success": true,
			"msg":     "收藏成功",
		}
		c.ServeJSON()
	} else {
		models.DB.Delete(&goodCollect)
		c.Data["json"] = map[string]interface{}{
			"success": true,
			"msg":     "取消收藏",
		}
		c.ServeJSON()
	}
}

func (c *ProductController) ProductItem() {
	c.BaseInit()

	id := c.Ctx.Input.Param(":id")
	/// get product information
	product := models.Product{}
	models.DB.Where("id=?", id).Find(&product)
	c.Data["product"] = product

	/// get relation product
	relationProduct := []models.Product{}
	product.RelationProduct = strings.ReplaceAll(product.RelationProduct, "，", ",")
	relationIds := strings.Split(product.RelationProduct, ",")
	models.DB.Where("id in (?)", relationIds).Select("id,title,price,product_version").Find(&relationProduct)
	c.Data["relationProduct"] = relationProduct

	/// get relation gift product
	productGift := []models.Product{}
	product.ProductGift = strings.ReplaceAll(product.ProductGift, "，", ",")
	giftIds := strings.Split(product.ProductGift, ",")
	models.DB.Where("id in (?)", giftIds).Select("id,title,price,product_img").Find(&productGift)
	c.Data["productGift"] = productGift

	/// get product color
	productColor := []models.ProductColor{}
	product.ProductColor = strings.ReplaceAll(product.ProductColor, "，", ",")
	colorIds := strings.Split(product.ProductColor, ",")
	models.DB.Where("id in (?)", colorIds).Find(&productColor)
	c.Data["productColor"] = productColor

	/// get product fitting
	productFitting := []models.Product{}
	product.ProductFitting = strings.ReplaceAll(product.ProductFitting, "，", ",")
	fittingIds := strings.Split(product.ProductFitting, ",")
	models.DB.Where("id in (?)", fittingIds).Select("id,title,price,product_img").Find(&productFitting)
	c.Data["productFitting"] = productFitting

	/// get product image
	productImage := []models.ProductImage{}
	models.DB.Where("product_id=?", product.Id).Find(&productImage)
	c.Data["productImage"] = productImage

	/// get product attribute
	productAttr := []models.ProductAttr{}
	models.DB.Where("product_id=?", product.Id).Find(&productAttr)
	c.Data["productAttr"] = productAttr
	c.TplName = "frontend/product/item.html"
}

func (c *ProductController) GetImgList() {
	colorId, err1 := c.GetInt("color_id")
	productId, err2 := c.GetInt("product_id")
	/// get product image
	productImage := []models.ProductImage{}
	err3 := models.DB.Where("color_id=? AND product_id=?", colorId, productId).Find(&productImage).Error

	if err1 != nil || err2 != nil || err3 != nil {
		c.Data["json"] = map[string]interface{}{
			"result":  "失败",
			"success": false,
		}
		c.ServeJSON()
	} else {
		if len(productImage) == 0 {
			models.DB.Where("product_id=?", productId).Find(&productImage)
		}
		c.Data["json"] = map[string]interface{}{
			"result":  productImage,
			"success": true,
		}
		c.ServeJSON()
	}
}
