package frontend

import "LeastMall/models"

type AddressController struct {
	BaseController
}

func (c *AddressController) AddAddress() {
	user := models.User{}
	models.Cookie.Get(c.Ctx, "userinfo", &user)
	name := c.GetString("name")
	phone := c.GetString("phone")
	address := c.GetString("address")
	zipcode := c.GetString("zipcode")
	var addressCount int
	models.DB.Where("uid=?", user.Id).Table("address").Count(&addressCount)
	if addressCount > 10 {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "增加地址失敗 地址數量超過 10 ",
		}
		c.ServeJSON()
		return
	}
	models.DB.Table("address").Where("uid=?", user.Id).Updates(map[string]interface{}{"default_address": 0})
	addressResult := models.Address{
		Uid:            user.Id,
		Name:           name,
		Phone:          phone,
		Address:        address,
		Zipcode:        zipcode,
		DefaultAddress: 1,
	}
	models.DB.Create(&addressResult)
	allAddressResult := []models.Address{}
	models.DB.Where("uid=?", user.Id).Find(&allAddressResult)
	c.Data["json"] = map[string]interface{}{
		"success": true,
		"result":  allAddressResult,
	}
	c.ServeJSON()
}

func (c *AddressController) GetOneAddressList() {
	addressId, err := c.GetInt("address_id")
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "傳入参數錯誤",
		}
		c.ServeJSON()
		return
	}
	address := models.Address{}
	models.DB.Where("id=?", addressId).Find(&address)
	c.Data["json"] = map[string]interface{}{
		"success": true,
		"result":  address,
	}
	c.ServeJSON()
}

func (c *AddressController) GoEditAddressList() {
	user := models.User{}
	models.Cookie.Get(c.Ctx, "userinfo", &user)
	addressId, err := c.GetInt("address_id")
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "傳入参數錯誤",
		}
		c.ServeJSON()
		return
	}
	name := c.GetString("name")
	phone := c.GetString("phone")
	address := c.GetString("address")
	zipcode := c.GetString("zipcode")
	models.DB.Table("address").Where("uid=?", user.Id).Updates(map[string]interface{}{"default_address": 0})
	addressModel := models.Address{}
	models.DB.Where("id=?", addressId).Find(&addressModel)
	addressModel.Name = name
	addressModel.Phone = phone
	addressModel.Address = address
	addressModel.Zipcode = zipcode
	addressModel.DefaultAddress = 1
	models.DB.Save(&addressModel)
	/// get current all address
	allAddressResult := []models.Address{}
	models.DB.Where("uid=?", user.Id).Order("default_address desc").Find(&allAddressResult)

	c.Data["json"] = map[string]interface{}{
		"success": true,
		"result":  allAddressResult,
	}
	c.ServeJSON()

}

func (c *AddressController) ChangeDefaultAddress() {
	user := models.User{}
	models.Cookie.Get(c.Ctx, "userinfo", &user)
	addressId, err := c.GetInt("address_id")
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "傳入参數錯誤",
		}
		c.ServeJSON()
		return
	}
	models.DB.Table("address").Where("uid=?", user.Id).Updates(map[string]interface{}{"default_address": 0})
	models.DB.Table("address").Where("id=?", addressId).Updates(map[string]interface{}{"default_address": 1})
	c.Data["json"] = map[string]interface{}{
		"success": true,
		"result":  "更新默認收貨地址成功",
	}
	c.ServeJSON()
}
