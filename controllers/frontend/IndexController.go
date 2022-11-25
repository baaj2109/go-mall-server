package frontend

import (
	"LeastMall/models"
	"fmt"
	"time"
)

type IndexController struct {
	BaseController
}

func (c *IndexController) Get() {
	/// init
	c.BaseInit()
	/// start time
	startTime := time.Now().UnixNano()
	/// get banner
	banner := []models.Banner{}
	hasBanner := models.CacheDB.Get("banner", &banner)
	if hasBanner {
		c.Data["bannerList"] = banner
	} else {
		models.DB.Where("status = 1 AND banner_type = 1").Order("sort desc").Find(&banner)
		c.Data["bannerList"] = banner
		models.CacheDB.Set("banner", banner)
	}

	/// get phone product list
	phoneList := []models.Product{}
	hasPhone := models.CacheDB.Get("phone", &phoneList)
	if hasPhone {
		c.Data["phoneList"] = phoneList
	} else {
		phoneList := models.GetProductByCategory(1, "hot", 8)
		c.Data["phoneList"] = phoneList
		models.CacheDB.Set("phone", phoneList)
	}

	/// get tv product list
	tvList := []models.Product{}
	hasTv := models.CacheDB.Get("tv", &tvList)
	if hasTv {
		c.Data["tvList"] = tvList
	} else {
		tv := models.GetProductByCategory(4, "best", 8)
		c.Data["tvList"] = tv
		models.CacheDB.Set("tv", tv)
	}

	/// finish time
	endTime := time.Now().UnixNano()
	fmt.Println("執行時間", endTime-startTime)
	c.TplName = "frontend/index/index.html"
	// c.TplName = "frontend/auth/register_step1.html"
}
