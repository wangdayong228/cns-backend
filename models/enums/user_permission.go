package enums

import (
	"errors"
	"strings"
)

type UserPermission uint

const (
	USER_PERMIT_REGISTER = iota + 1
	USER_PERMIT_RENEW    = USER_PERMIT_REGISTER << 1
)

var (
	userPermissionValue2StrMap map[UserPermission]string
	userPermissionStr2ValueMap map[string]UserPermission
)

var (
	ErrUnkownUserPermission = errors.New("unknown user permission")
)

func init() {
	userPermissionValue2StrMap = map[UserPermission]string{
		USER_PERMIT_REGISTER: "REGISTER",
		USER_PERMIT_RENEW:    "RENEW",
	}

	userPermissionStr2ValueMap = make(map[string]UserPermission)
	for k, v := range userPermissionValue2StrMap {
		userPermissionStr2ValueMap[v] = k
	}
}

func (t *UserPermission) HasPermission(perm UserPermission) bool {
	return *t&perm == 1
}

func (t *UserPermission) String() string {
	v, ok := userPermissionValue2StrMap[*t]
	if ok {
		return v
	}

	var permits []UserPermission

	if t.HasPermission(USER_PERMIT_REGISTER) {
		permits = append(permits, USER_PERMIT_REGISTER)
	}
	if t.HasPermission(USER_PERMIT_RENEW) {
		permits = append(permits, USER_PERMIT_RENEW)
	}

	permitStrs := []string{}
	for _, p := range permits {
		permitStrs = append(permitStrs, p.String())
	}
	return strings.Join(permitStrs, ",")
}

func ParseUserPermission(str string) ([]UserPermission, bool) {
	v, ok := userPermissionStr2ValueMap[str]
	if ok {
		return []UserPermission{v}, true
	}

	permitStrs := strings.Split(str, ",")

	var permits []UserPermission
	for _, p := range permitStrs {
		permit, ok := ParseUserPermission(p)
		if !ok {
			return nil, false
		}
		permits = append(permits, permit...)
	}
	return permits, true
}
