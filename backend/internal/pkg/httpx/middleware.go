package httpx

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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