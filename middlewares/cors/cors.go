package cors

import (
	"github.com/gin-gonic/gin"
)

func CorsSeting() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "PUT,DELETE,POST,OPTIONS,GET,PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With,X-PAGE-CURRENT,Order-Attr,Order-Type,X-PAGE-SIZE,Sid,UserId")

		// ie 禁用缓存
		c.Writer.Header().Set("expries", "-1")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Pragma", "no-cache")
		if "OPTIONS" == c.Request.Method {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}

	}
}
