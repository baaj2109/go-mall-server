package models

import (
	_ "github.com/jinzhu/gorm"
)

type Auth struct {
	Id          int
	ModuleName  string
	ActionName  string
	Type        int // 1.模塊 2.菜單 3. 動作
	Url         string
	ModuleId    int
	Description string
	Status      int
	AddTime     int
	Sort        int
	AuthItem    []Auth `gorm:"foreignkey:ModuleId;association_foreignkey:Id"`
	Checked     bool   `gorm:"-"` // 忽略本字段
}

func (Auth) TableName() string {
	return "auth"
}
