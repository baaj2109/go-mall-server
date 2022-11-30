package backend

import (
	"LeastMall/common"
	"LeastMall/models"
	"os"
	"strconv"

	"github.com/astaxie/beego"
)

type BannerController struct {
	BaseController
}

func (c *BannerController) Get() {
	banner := []models.Banner{}
	models.DB.Find(&banner)
	c.Data["bannerList"] = banner
	c.TplName = "backend/banner/index.html"
}

func (c *BannerController) Add() {
	c.TplName = "backend/banner/add.html"
}

func (c *BannerController) GoAdd() {
	bannerType, err1 := c.GetInt("banner_type")
	title := c.GetString("title")
	link := c.GetString("link")
	sort, err2 := c.GetInt("sort")
	status, err3 := c.GetInt("status")
	if err1 != nil || err3 != nil {
		c.Error("非法請求", "/banner")
		return
	}
	if err2 != nil {
		c.Error("排序表單里面输入的數據不合法", "/banner/add")
		return
	}
	bannerImgSrc, err4 := c.UploadImg("banner_img")
	if err4 == nil {
		banner := models.Banner{
			Title:      title,
			BannerType: bannerType,
			BannerImg:  bannerImgSrc,
			Link:       link,
			Sort:       sort,
			Status:     status,
			AddTime:    int(common.GetUnix()),
		}
		models.DB.Create(&banner)
		c.Success("增加輪播圖成功", "/banner")
	} else {
		c.Error("增加輪播圖失敗", "/banner/add")
		return
	}
}

func (c *BannerController) Edit() {
	id, err := c.GetInt("id")
	if err != nil {
		c.Error("非法請求", "/banner")
		return
	}
	banner := models.Banner{Id: id}
	models.DB.Find(&banner)
	c.Data["banner"] = banner
	c.TplName = "backend/banner/edit.html"
}

func (c *BannerController) GoEdit() {
	id, err := c.GetInt("id")
	bannerType, err1 := c.GetInt("banner_type")
	title := c.GetString("title")
	link := c.GetString("link")
	sort, err2 := c.GetInt("sort")
	status, err3 := c.GetInt("status")
	if err != nil || err1 != nil || err3 != nil {
		c.Error("非法請求", "/banner")
		return
	}
	if err2 != nil {
		c.Error("排序表單里面输入的數據不合法", "/banner/edit?id="+strconv.Itoa(id))
		return
	}
	bannerImgSrc, _ := c.UploadImg("banner_img")
	banner := models.Banner{Id: id}
	models.DB.Find(&banner)
	banner.Title = title
	banner.BannerType = bannerType
	banner.Link = link
	banner.Sort = sort
	banner.Status = status
	if bannerImgSrc != "" {
		banner.BannerImg = bannerImgSrc
	}
	err5 := models.DB.Save(&banner).Error
	if err5 != nil {
		c.Error("修改輪播圖失敗", "/banner/edit?id="+strconv.Itoa(id))
		return
	}
	c.Success("修改輪播圖成功", "/banner")
}

func (c *BannerController) Delete() {
	id, err := c.GetInt("id")
	if err != nil {
		c.Error("傳入参數錯誤", "/banner")
		return
	}
	banner := models.Banner{Id: id}
	models.DB.Find(&banner)
	// address := "D:/gowork/src/gitee.com/shirdonl/LeastMall/" + banner.BannerImg
	address := "/Users/baaj2109/go/src/LeastMall/" + banner.BannerImg
	test := os.Remove(address)
	if test != nil {
		beego.Error(test)
		c.Error("刪除本機上圖片錯誤", "/banner")
		return
	}
	models.DB.Delete(&banner)
	c.Success("刪除輪播圖成功", "/banner")
}
