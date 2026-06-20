package helpers

import (
	"authorization-microservice/constants"
	"authorization-microservice/models"
	"fmt"
	"math/rand"
	"net/http"
	"strings"

	"github.com/AndreyLebedev1998/auth-grpc"
	"github.com/golang-jwt/jwt/v5"
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

/* func GetParamUserForAuth(paramUser models.Entrance) []any {
	var params []any
	if paramUser.Name != "" {
		params = append(params, paramUser.Name)
	}
	if paramUser.Email != "" {
		params = append(params, paramUser.Email)
	}
	if paramUser.Phone != "" {
		params = append(params, paramUser.Phone)
	}

	if paramUser.Email == "" && paramUser.Phone == "" && paramUser.Name == "" {
		return nil
	}

	return params
} */

func GetUserIDFromToken(r *http.Request) (int, error) {
	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" {
		return 0, fmt.Errorf("no token")
	}

	if len(tokenStr) > 7 && tokenStr[:7] == "Bearer " {
		tokenStr = tokenStr[7:]
	}

	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return constants.JwtKey, nil
	})

	if err != nil || !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	return claims.UserID, nil
}
