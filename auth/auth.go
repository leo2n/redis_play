package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/valyala/fasthttp"
	"go_redis/jsonStruct"
	"go_redis/mysql/shop/structure"
	"go_redis/utils"
	"net/http"
	"time"
)

// server side sign token need secret
var secret = "1hXNV1rlgoEoT9U9gWqSmyYS9G1"

// 生成符合要求的JWT token
// 要求如下: 24h后过期
func GenerateToken(user *structure.UserLogin) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // 24 hours expire
	})
	return token.SignedString([]byte(secret))
}

// token middleware
func MiddleAuth(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		// 首先, 验证header中key: authorization的值是否符合要求?
		tokenStr := string(ctx.Request.Header.Peek("Authorization"))
		if tokenStr == "" {
			utils.ResponseWithJson(ctx, 401, jsonStruct.CommonResponse{
				Code: 8401,
				Msg:  "unauthorized",
				Data: nil,
			})
		} else {
			// 验证token是否可以被解析
			token, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (i interface{}, err error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					utils.ResponseWithJson(ctx, http.StatusUnauthorized, jsonStruct.CommonResponse{
						Code: 8401,
						Msg:  "unauthorized",
						Data: nil,
					})
				}
				return []byte(""), nil // default return
			})
			// 验证token是否合法
			if !token.Valid {
				utils.ResponseWithJson(ctx, 401, jsonStruct.CommonResponse{
					Code: 8401,
					Msg:  "unauthorized",
					Data: nil,
				})
			}
			handler(ctx)
		}
	}
}