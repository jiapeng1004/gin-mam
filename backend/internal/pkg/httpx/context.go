package httpx

import "github.com/gin-gonic/gin"

const ctxKeyUserID = "userId"
const ctxKeyUsername = "username"

func GetUserID(c *gin.Context) string {
	v, ok := c.Get(ctxKeyUserID)
	if !ok {
		return ""
	}
	s, _ := v.(string)
	return s
}

func GetUsername(c *gin.Context) string {
	v, ok := c.Get(ctxKeyUsername)
	if !ok {
		return ""
	}
	s, _ := v.(string)
	return s
}
