package jwttoken

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var JwtKey = []byte("your_secret_key")

type Claims struct {
	PhoneNumber string `json:"phone_number"`
	jwt.StandardClaims
}

func AuthenticateJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if len(tokenString) < 7 {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthenticated"})
			c.Abort()
			return
		}
		tokenString = tokenString[7:]

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthenticated"})
			c.Abort()
			return
		}
		c.Set("phone_number", claims.PhoneNumber)
		c.Next()
	}
}
