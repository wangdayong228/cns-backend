package services

import "errors"

var (
	ErrMakeCommithashFirst       = errors.New("commitment not found, please make commit before make order")
	ErrCommitsUnsubmitOnContract = errors.New("commitment invalid: not submit on contract")
	ErrCommitsExpired            = errors.New("commitment invalid: expired")
	ErrOrderUnexists             = errors.New("order is exists, if need refresh url please invoke API 'RefreshUrl'")
	ErrOrderCompleted            = errors.New("order is completed")
)
