package models

// type User struct {
// 	Id       int
// 	Phone    string
// 	Password string
// 	AddTime  int
// 	LastIp   string
// 	Email    string
// 	Status   int
// }

type User struct {
	Id       int
	Email    string
	Password string
	AddTime  int
	LastIP   string
	Status   int
}

func (User) TableName() string {
	return "user"
}
