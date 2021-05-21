package jwt

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"easycache/pkg/define"
	"easycache/pkg/util"
)

// JWT is jwt middleware
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}

		code = define.SUCCESS
		token := c.Request.Header.Get("token")
		if token == "" {
			code = define.INVALID_PARAMS
		} else {
			_, err := util.ParseToken(token)
			if err != nil {
				switch err.(*jwt.ValidationError).Errors {
				case jwt.ValidationErrorExpired:
					code = define.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
				default:
					code = define.ERROR_AUTH_CHECK_TOKEN_FAIL
				}
			}
		}

		if code != define.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  define.GetMsg(code),
				"data": data,
			})

			c.Abort()
			return
		}

		c.Next()
	}
}
