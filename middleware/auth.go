package middleware

import (
	"dbmsbackend/util"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Auth(config *util.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")

		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				util.GeneralAPIResponse{
					Status:  http.StatusUnauthorized,
					Message: "missing jwt token",
				},
			)
			return
		}
		token := strings.Split(auth, "Bearer ")[1]

		claims, err := util.ValidateToken(config, token)

		if err, ok := err.(*util.ErrUnauthorized); ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				util.GeneralAPIResponse{
					Status:  http.StatusUnauthorized,
					Message: err.Error(),
				},
			)
			return
		}

		c.Set("userID", claims.ID)
		c.Next()
	}
}
