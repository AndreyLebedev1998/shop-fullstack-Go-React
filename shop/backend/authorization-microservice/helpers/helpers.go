package helpers

import (
	"math/rand"

	"github.com/AndreyLebedev1998/auth-grpc"
)

func RandomString() string {
	var str = "absdefghijklmnopABCDEFGHIJKLMNOP1234567890"
	var result string
	var strLength int = 42
	for i := 0; i < 42; i++ {
		result += string(str[rand.Intn(strLength)])
	}
	return result
}

func GetParamUser(paramUser *auth.ParamUser) any {
	if paramUser.Id != 0 {
		return paramUser.Id
	} else if paramUser.Email != "" {
		return paramUser.Email
	} else if paramUser.Phone != "" {
		return paramUser.Phone
	} else {
		return nil
	}
}
