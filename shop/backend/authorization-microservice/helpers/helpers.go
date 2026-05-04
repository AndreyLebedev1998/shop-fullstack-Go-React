package helpers

import (
	"authorization-microservice/models"
	"math/rand"
	"net/http"
	"strings"

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

func SixRandomNumbers() string {
	sliceInt := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}
	var result string
	for i := 0; i < 6; i++ {
		result += string(sliceInt[rand.Intn(10)])
	}
	return result
}

func GetClientIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")

	if ip != "" {
		// может быть список: "ip1, ip2"
		return strings.Split(ip, ",")[0]
	}

	ip = r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}

	// fallback
	return r.RemoteAddr
}

func GetParamUserForAuth(paramUser models.Entrance) any {
	if paramUser.Email != "" {
		return paramUser.Email
	} else if paramUser.Phone != "" {
		return paramUser.Phone
	} else if paramUser.Name != "" {
		return paramUser.Name
	} else {
		return nil
	}
}
