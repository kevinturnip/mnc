package auth

import (
	"mnc/db"
	jwttoken "mnc/jwt"
	"mnc/model"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var loginData struct {
		PhoneNumber string `json:"phone_number"`
		Pin         string `json:"pin"`
	}
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	var user model.User
	if err := db.DB.Where("phone_number = ? AND pin = ?", loginData.PhoneNumber, loginData.Pin).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Phone Number and PIN doesn’t match."})
		return
	}

	// Create JWT Token
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &jwttoken.Claims{
		PhoneNumber: loginData.PhoneNumber,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwttoken.JwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not login"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "SUCCESS",
		"result": gin.H{
			"access_token": tokenString,
		},
	})
}
