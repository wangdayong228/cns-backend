package models

import (
	"github.com/wangdayong228/cns-backend/models/enums"
)

type User struct {
	BaseModel
	Name       string               `gorm:"type:varchar(255)" json:"name"`
	ApiKeyHash string               `gorm:"type:varchar(255)" json:"api_key"`
	Permission enums.UserPermission `gorm:"uint" json:"permission"`
}

func GetAllUsers() ([]*User, error) {
	users := []*User{}
	if err := GetDB().Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
