package helpers

import (
	"btpn-syariah-final-project/database"
	"btpn-syariah-final-project/models"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func GetUserLogin(c *gin.Context) (*models.User, error) {
	tokenString, err := c.Cookie("Authorization")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil || !token.Valid {
		c.AbortWithStatus(http.StatusUnauthorized)
		return nil, fmt.Errorf("Invalid or expired token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// check the exp
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return nil, fmt.Errorf("Token has expired")
		}

		// find the user with token sub
		var user models.User
		if err := database.DB.First(&user, claims["sub"]).Error; err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return nil, fmt.Errorf("User not found")
		}

		return &user, nil
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
		return nil, fmt.Errorf("Invalid token claims")
	}
}