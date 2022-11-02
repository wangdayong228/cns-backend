package models

type User struct {
	BaseModel
	ApiKey     string `json:"api_key"`
	Permission uint   `json:"permission"`
}
