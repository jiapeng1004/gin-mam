package httpx

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const ErrCodeUnauthorized = 40100

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) == 0 || c.Writer.Written() {
			return
		}
		err := c.Errors.Last().Err
		if biz, ok := err.(*BizError); ok {
			Fail(c, biz.HTTPStatus, biz.ErrCode, biz.ErrMsg)
			return
		}
		Fail(c, http.StatusInternalServerError, 50000, "内部错误")
	}
}

func JWTMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			Fail(c, http.StatusUnauthorized, ErrCodeUnauthorized, "未授权")
			c.Abort()
			return
		}
		tokenStr := strings.TrimSpace(strings.TrimPrefix(auth, "Bearer "))
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			if t.Method != jwt.SigningMethodHS256 {
				return nil, jwt.ErrTokenSignatureInvalid
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			Fail(c, http.StatusUnauthorized, ErrCodeUnauthorized, "未授权")
			c.Abort()
			return
		}
		c.Next()
	}
}