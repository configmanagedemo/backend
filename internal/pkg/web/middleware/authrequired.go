package middleware

import (
	"fmt"
	"main/internal/pkg/cache"
	"main/internal/pkg/e"
	"main/internal/pkg/web/service"

	"github.com/gin-gonic/gin"
)

// AuthRequired 检查登录权限
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取
		var (
			token string
			err   error
		)
		token, err = c.Cookie("token")
		if err != nil {
			token = c.GetHeader("x-token")
		}

		if token == "" {
			service.ResponseError(c, e.NotLogin, "")
			return
		}

		if val, err := cache.Get(token); err != nil {
			if err != cache.Nil {
				panic(err.Error())
			}
			service.ResponseError(c, e.NotLogin, "")
			return
		} else {
			fmt.Println(val)
			c.Set("uid", val)
		}
		c.Next()
	}
}
