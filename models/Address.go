package models

//	type Address struct {
//		Id      int    `json:"id"`
//		Uid     int    `json:"uid"`
//		Phone   string `json:"phone"`
//		Name    string `json:"name"`
//		Zipcode string `json:"zipcode"`
//		AddTime string `json:"add_time"`
//	}
type Address struct {
	Id             int    `json:"id"`
	Uid            int    `json:"uid"`
	Phone          string `json:"phone"`
	Name           string `json:"name"`
	Address        string `json:"address"`
	Zipcode        string `json:"zipcode"`
	DefaultAddress int    `json:"default_address"`
	AddTime        int    `json:"add_time"`
}

func (Address) TableName() string {
	return "address"
}
