package backend

import (
	"LeastMall/common"
	"LeastMall/models"
	"strconv"
	"strings"
)

type ProductTypeAttrController struct {
	BaseController
}

func (c *ProductTypeAttrController) Get() {

	cateId, err1 := c.GetInt("cate_id")
	if err1 != nil {
		c.Error("非法請求", "/productType")
	}
	//獲取當前的類型
	productType := models.ProductType{Id: cateId}
	models.DB.Find(&productType)
	c.Data["productType"] = productType

	//查詢當前類型下面的商品類型属性
	productTypeAttr := []models.ProductTypeAttribute{}
	models.DB.Where("cate_id=?", cateId).Find(&productTypeAttr)
	c.Data["productTypeAttrList"] = productTypeAttr

	c.TplName = "backend/productTypeAttribute/index.html"
}

func (c *ProductTypeAttrController) Add() {

	cateId, err1 := c.GetInt("cate_id")
	if err1 != nil {
		c.Error("非法請求", "/productType")
	}

	productType := []models.ProductType{}
	models.DB.Find(&productType)
	c.Data["productTypeList"] = productType
	c.Data["cateId"] = cateId
	c.TplName = "backend/productTypeAttribute/add.html"
}

func (c *ProductTypeAttrController) GoAdd() {

	title := c.GetString("title")
	cateId, err1 := c.GetInt("cate_id")
	attrType, err2 := c.GetInt("attr_type")
	attrValue := c.GetString("attr_value")
	sort, err4 := c.GetInt("sort")
	if err1 != nil || err2 != nil {
		c.Error("非法請求", "/productType")
		return
	}
	if strings.Trim(title, " ") == "" {
		c.Error("商品類型属性名稱不能為空", "/productTypeAttribute/add?cate_id="+strconv.Itoa(cateId))
		return
	}
	if err4 != nil {
		c.Error("排序值錯誤", "/productTypeAttribute/add?cate_id="+strconv.Itoa(cateId))
		return
	}
	productTypeAttr := models.ProductTypeAttribute{
		Title:     title,
		CateId:    cateId,
		AttrType:  attrType,
		AttrValue: attrValue,
		Status:    1,
		AddTime:   int(common.GetUnix()),
		Sort:      sort,
	}
	models.DB.Create(&productTypeAttr)
	c.Success("增加成功", "/productTypeAttribute?cate_id="+strconv.Itoa(cateId))

}

func (c *ProductTypeAttrController) Edit() {
	id, err1 := c.GetInt("id")
	if err1 != nil {
		c.Error("非法請求", "/goodType")
		return
	}
	productTypeAttr := models.ProductTypeAttribute{Id: id}
	models.DB.Find(&productTypeAttr)
	c.Data["productTypeAttr"] = productTypeAttr
	productType := []models.ProductType{}
	models.DB.Find(&productType)
	c.Data["productTypeList"] = productType
	c.TplName = "backend/productTypeAttribute/edit.html"
}

func (c *ProductTypeAttrController) GoEdit() {
	id, err := c.GetInt("id")
	title := c.GetString("title")
	cateId, err1 := c.GetInt("cate_id")
	attrType, err2 := c.GetInt("attr_type")
	attrValue := c.GetString("attr_value")
	sort, err4 := c.GetInt("sort")
	if err != nil || err1 != nil || err2 != nil {
		c.Error("非法請求", "/productTypeAttribute")
		return
	}
	if strings.Trim(title, " ") == "" {
		c.Error("商品類型属性名稱不能為空", "/productTypeAttribute/edit?cate_id="+strconv.Itoa(id))
		return
	}
	if err4 != nil {
		c.Error("排序值錯誤", "/productTypeAttribute/edit?cate_id="+strconv.Itoa(id))
		return
	}
	productTypeAttr := models.ProductTypeAttribute{Id: id}
	models.DB.Find(&productTypeAttr)
	productTypeAttr.Title = title
	productTypeAttr.CateId = cateId
	productTypeAttr.AttrType = attrType
	productTypeAttr.AttrValue = attrValue
	productTypeAttr.Sort = sort
	err3 := models.DB.Save(&productTypeAttr).Error
	if err3 != nil {
		c.Error("修改數據失敗", "/productTypeAttribute/edit?cate_id="+strconv.Itoa(id))
	}
	c.Success("修改數據成功", "/productTypeAttribute?cate_id="+strconv.Itoa(cateId))
}
func (c *ProductTypeAttrController) Delete() {
	id, err := c.GetInt("id")
	cateId, err1 := c.GetInt("cate_id")
	if err != nil {
		c.Error("傳入参數錯誤", "/productTypeAttribute?cate_id="+strconv.Itoa(cateId))
		return
	}
	if err1 != nil {
		c.Error("非法請求", "/productType")
	}
	productTypeAttr := models.ProductTypeAttribute{Id: id}
	models.DB.Delete(&productTypeAttr)
	c.Success("刪除數據成功", "/productTypeAttribute?cate_id="+strconv.Itoa(cateId))
}
