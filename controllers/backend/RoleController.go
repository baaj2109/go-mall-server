package backend

import (
	"LeastMall/common"
	"LeastMall/models"
	"strconv"
	"strings"
)

type RoleController struct {
	BaseController
}

func (c *RoleController) Get() {
	role := []models.Role{}
	models.DB.Find(&role)
	c.Data["rolelist"] = role
	c.TplName = "backend/role/index.html"
}

func (c *RoleController) Add() {
	c.TplName = "backend/role/add.html"
}

func (c *RoleController) GoAdd() {
	title := strings.Trim(c.GetString("title"), "")
	description := strings.Trim(c.GetString("description"), "")
	if title == "" {
		c.Error("title can't be null", "/role/add")
		return
	}
	roleList := []models.Role{}
	models.DB.Where("title=?", title).Find(&roleList)
	if len(roleList) != 0 {
		c.Error("部門已存在！", "/role/add")
		return
	}
	role := models.Role{}
	role.Title = title
	role.Description = description
	role.Status = 1
	role.AddTime = int(common.GetUnix())
	err := models.DB.Create(&role).Error
	if err != nil {
		c.Error("增加部門失敗", "/role/add")
	} else {
		c.Success("增加部門成功", "/role")
	}
}

func (c *RoleController) Edit() {
	id, err := c.GetInt("id")
	if err != nil {
		c.Error("傳入参數錯誤", "/role")
		return
	}
	role := models.Role{Id: id}
	models.DB.Find(&role)
	c.Data["role"] = role
	c.TplName = "backend/role/edit.html"
}

func (c *RoleController) GoEdit() {
	title := strings.Trim(c.GetString("title"), "")
	description := strings.Trim(c.GetString("description"), "")
	id, err := c.GetInt("id")
	if err != nil {
		c.Error("傳入参數錯誤", "/role")
		return
	}
	role := models.Role{Id: id}
	models.DB.Find(&role)
	role.Title = title
	role.Description = description
	err2 := models.DB.Save(&role).Error
	if err2 != nil {
		c.Error("修改部門失敗", "/role/edit?id="+strconv.Itoa(id))
	} else {
		c.Success("修改部門成功", "/role")
	}
}

func (c *RoleController) Delete() {
	id, err := c.GetInt("id")
	if err != nil {
		c.Error("傳入参數錯誤", "/role")
		return
	}
	role := models.Role{Id: id}
	administrator := []models.Administrator{}
	roleAuth := models.RoleAuth{}
	models.DB.Where("role_id=?", id).Delete(&roleAuth)
	models.DB.Preload("Role").Where("role_id=?", id).Find(&administrator)
	if len(administrator) > 0 {
		c.Error("仍有未處理的員工，無法刪除部門", "/role")
		return
	}
	models.DB.Delete(&role)
	c.Success("刪除部門成功", "/role")
}

func (c *RoleController) Auth() {
	roleId, err := c.GetInt("id")
	if err != nil {
		c.Error("傳入参數錯誤", "/role")
		return
	}
	auth := []models.Auth{}
	models.DB.Preload("AuthItem").Where("module_id=0").Find(&auth)
	/// get auth of current department, and put authItem in a map
	roleAuth := []models.RoleAuth{}
	models.DB.Where("role_id=?", roleId).Find(&roleAuth)
	roleAuthMap := make(map[int]int)
	for _, v := range roleAuth {
		roleAuthMap[v.AuthId] = v.AuthId
	}
	for i := 0; i < len(auth); i++ {
		if _, ok := roleAuthMap[auth[i].Id]; ok {
			auth[i].Checked = true
		}
		for j := 0; j < len(auth[i].AuthItem); j++ {
			if _, ok := roleAuthMap[auth[i].AuthItem[j].Id]; ok {
				auth[i].AuthItem[j].Checked = true
			}
		}
	}

	c.Data["authList"] = auth
	c.Data["roleId"] = roleId
	c.TplName = "backend/role/auth.html"
}

func (c *RoleController) GoAuth() {
	roleId, err := c.GetInt("id")
	if err != nil {
		c.Error("傳入参數錯誤", "/role")
		return
	}
	authNode := c.GetStrings("auth_node")
	roleAuth := models.RoleAuth{}
	models.DB.Where("role_id=?", roleId).Delete(&roleAuth)
	for _, v := range authNode {
		authId, _ := strconv.Atoi(v)
		roleAuth.AuthId = authId
		roleAuth.RoleId = roleId
		models.DB.Create(&roleAuth)
	}
	c.Success("授權成功", "/role/auth?id="+strconv.Itoa(roleId))
}
