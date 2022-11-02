package enums

import "errors"

type OrderState uint

const (
	ORDER_STATE_INIT = iota + 1
	ORDER_STATE_MADE
	ORDER_STATE_SUCCESS
)

var (
	orderStateValue2StrMap map[OrderState]string
	orderStateStr2ValueMap map[string]OrderState
)

var (
	ErrUnkownOrderState = errors.New("unknown order state")
)

func init() {
	orderStateValue2StrMap = map[OrderState]string{
		ORDER_STATE_INIT:    "init",
		ORDER_STATE_MADE:    "made",
		ORDER_STATE_SUCCESS: "success",
	}

	orderStateStr2ValueMap = make(map[string]OrderState)
	for k, v := range orderStateValue2StrMap {
		orderStateStr2ValueMap[v] = k
	}
}

func (t *OrderState) String() string {
	v, ok := orderStateValue2StrMap[*t]
	if ok {
		return v
	}
	return "unknown"
}

func ParseOrderState(str string) (*OrderState, bool) {
	v, ok := orderStateStr2ValueMap[str]
	return &v, ok
}
