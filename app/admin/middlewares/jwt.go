package middlewares

import (
	"fmt"
	"go-admin/internal/lib/config"
	"go-admin/pkg/utils/auth"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenKey := "Authorization"
		tokenString := c.Request.Header.Get(tokenKey)

		bearerPrefix := "Bearer "
		tokenString = strings.TrimPrefix(tokenString, bearerPrefix)

		secretKey := config.Settings.AdminServer.Server.JwtSecret
		token, err := jwt.ParseWithClaims(tokenString, &auth.MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusForbidden,
				"message": fmt.Sprintf("token parse error: %v", err),
			})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusForbidden,
				"message": fmt.Sprintf("token valid error: %v", err),
			})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*auth.MyCustomClaims); ok {
			fmt.Printf("claims: %+v\n\n", claims)
			fmt.Printf("%+v\n\n", claims.User)
			// username := claims.Subject
			// username := claims.User["username"].(string)
			userID := claims.User["id"].(float64)
			c.Request.Header.Set("user_id", strconv.Itoa(int(userID)))

			c.Next()
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusForbidden,
				"message": fmt.Sprintf("get claims error: %v", err),
			})
			c.Abort()
			return
		}
	}
}
