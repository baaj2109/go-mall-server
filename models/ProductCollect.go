package models

type ProductCollect struct {
	Id        int
	UserId    int
	ProductId int
	AddTime   string
}

func (ProductCollect) TableName() string {
	return "product_collect"
}
