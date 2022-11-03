package models

import "github.com/wangdayong228/cns-backend/models/enums"

type User struct {
	BaseModel
	ApiKey     string               `json:"api_key"`
	Permission enums.UserPermission `json:"permission"`
}
