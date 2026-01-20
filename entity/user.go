package entity

type User struct {
	ID   string `gorm:"type:varchar(32);primaryKey" json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (User) TableName() string {
	return "user"
}
