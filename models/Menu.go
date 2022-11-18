package models

import (
	_ "github.com/jinzhu/gorm"
)

type Menu struct {
	Id          int
	Title       string
	Link        string
	Position    int
	IsOpennew   int
	Relation    string
	Sort        int
	Status      int
	Addtime     int
	ProductItem []Product `gorm:"-"`
}
