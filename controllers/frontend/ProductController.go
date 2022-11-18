package frontend

import (
	"LeastMall/common"
	"LeastMall/models"
	"math"
	"strconv"
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
