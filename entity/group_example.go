package entity

type GroupExample struct {
	ID         int    `gorm:"primaryKey;column:id" json:"id"`
	Name       string `gorm:"column:name" json:"name"`
	Department string `gorm:"column:department" json:"department"`
}

func (GroupExample) TableName() string {
	return "group_example"
}
