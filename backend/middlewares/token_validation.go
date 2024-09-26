package middlewares

import (
	"os"
	"strings"

	"github.com/chekuthankl13/sparrow_dine/helpers"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func ValidateToken() gin.HandlerFunc {
	key := os.Getenv("SECRET_KEY")
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			helpers.UnauthorizedResponse(c, "Unauthorized")
			return
		}
		// fmt.Println(strings.Split(tokenString, " ")[1])
		token, err := jwt.Parse(strings.Split(tokenString, " ")[1], func(t *jwt.Token) (interface{}, error) {
			return []byte(key), nil
		})

		if err != nil {
			helpers.UnauthorizedResponse(c, err.Error())
			return
		}

		if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Next()
		} else {
			helpers.UnauthorizedResponse(c, "Invalid token")
			return
		}

	}
}
