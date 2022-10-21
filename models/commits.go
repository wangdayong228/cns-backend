package models

import "github.com/wangdayong228/cns-backend/models/enums"

type Commit struct {
	BaseModel
	CommitCore
}

type CommitCore struct {
	CommitArgs
	CommitHash string           `gorm:"type:varchar(255);uniqueIndex" json:"commit_hash" binding:"required"`
	OrderState enums.OrderState `json:"order_state,omitempty"`
}

type CommitArgs struct {
	Name          string   `gorm:"type:varchar(255)" json:"name"`
	Owner         string   `gorm:"type:varchar(255);index" json:"owner"` //base32地址或hex地址
	Duration      uint     `json:"duration"`
	Secret        string   `gorm:"type:varchar(255)" json:"secret"` //32字节
	Resolver      string   `gorm:"type:varchar(255)" json:"resolver"`
	Data          []string `gorm:"type:varchar(255)" json:"data"` //32字节
	ReverseRecord bool     `json:"reverse_record"`
	Fuses         uint     `json:"fuses"`
	WrapperExpiry uint     `json:"wrapper_expiry"`
}

func FindCommit(commitHash string) (*Commit, error) {
	c := &Commit{}
	c.CommitHash = commitHash
	return c, GetDB().Where(c).First(c).Error
}

func FindCommits(condition *Commit, offset int, limit int) ([]*Commit, error) {
	commits := []*Commit{}
	return commits, GetDB().Where(condition).Find(&commits).Offset(offset).Limit(limit).Error
}
