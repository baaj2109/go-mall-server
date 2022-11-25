package frontend

import (
	"LeastMall/models"
	"fmt"
	"net/url"
	"strings"

	"github.com/astaxie/beego"
	"github.com/jinzhu/gorm"
)

type BaseController struct {
	beego.Controller
}

func (c *BaseController) BaseInit() {
	/// get top menu list
	topMenu := []models.Menu{}
	hasTopMenu := models.CacheDB.Get(
		"topMenu", &topMenu,
	)
	if hasTopMenu {
		c.Data["topMenuList"] = topMenu
	} else {
		models.DB.Where("status = 1 AND position = 1").Order("sort desc").Find(&topMenu)
		c.Data["topMenuList"] = topMenu
		models.CacheDB.Set("topMenu", topMenu)
	}

	/// load left side category (pre-load)
	productCate := []models.ProductCate{}
	hasProductCate := models.CacheDB.Get(
		"productCate", &productCate,
	)
	if hasProductCate {
		c.Data["productCate"] = productCate
	} else {
		models.DB.Preload("ProductCateItem",
			func(db *gorm.DB) *gorm.DB {
				return db.Where("product_cate.status=1").
					Order("product_cate.sort DESC")
			}).Where("pid=0 AND status=1").Order("sort desc", true).Find(&productCate)
		c.Data["productCateList"] = productCate
		models.CacheDB.Set("productCate", productCate)
	}

	/// load middle menu
	middleMenu := []models.Menu{}
	hasMiddleMenu := models.CacheDB.Get("middleMenu", &middleMenu)
	if hasMiddleMenu {
		c.Data["middleMenuList"] = middleMenu
	} else {
		models.DB.Where("status=1 AND position=2").Order("sort desc").
			Find(&middleMenu)

		for i := 0; i < len(middleMenu); i++ {
			/// get associate product
			middleMenu[i].Relation = strings.ReplaceAll(middleMenu[i].Relation, "，", ",")
			relation := strings.Split(middleMenu[i].Relation, ",")
			product := []models.Product{}
			models.DB.Where("id in (?)", relation).Limit(6).Order("sort ASC").
				Select("id,title,product_img,price").Find(&product)
			middleMenu[i].ProductItem = product
		}
		c.Data["middleMenuList"] = middleMenu
		models.CacheDB.Set("middleMenu", middleMenu)
	}

	/// 判斷使用者是否登入
	user := models.User{}
	models.Cookie.Get(c.Ctx, "userinfo", &user)
	if len(user.Phone) == 4 {
		str := fmt.Sprintf(`<ul>
			<li class="userinfo">
				<a href="#">%v</a>
				<i class="i"></i>
				<ol>
					<li><a href="/user">個人中心</a></li>
					<li><a href="#">我的收藏</a></li>
					<li><a href="/auth/loginOut">退出</a></li>
				</ol>
			</li>
		</ul> `, user.Phone)
		c.Data["userinfo"] = str
	} else {
		str := fmt.Sprint(`<ul>
			<li><a href="/auth/login" target="_self">登入</a></li>
			<li>|</li>
			<li><a href="/auth/registerStep1" target="_self" >註冊</a></li>
			</ul>`)
		c.Data["userinfo"] = str
	}

	urlPath, _ := url.Parse(c.Ctx.Request.URL.String())
	c.Data["pathname"] = urlPath.Path
}
