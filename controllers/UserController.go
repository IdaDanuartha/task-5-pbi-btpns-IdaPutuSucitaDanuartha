package controllers

import (
	"btpn-syariah-final-project/database"
	"btpn-syariah-final-project/models"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {	
	var body struct {
		models.User
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
      "error": "Failed to hash password",
    })

		return
	}

	user := models.User{
    Username: body.Username,
    Email: body.Email,
    Password: string(hash),
  }

	result := database.DB.Create(&user)

	if result.Error!= nil {
		c.JSON(http.StatusBadRequest, gin.H{
      "error": "Failed to create user",
    })

		return
	}

	c.JSON(http.StatusCreated, gin.H{
    "message": "User created successfully",
  })
}

func Login(c *gin.Context) {
	var body struct {
		models.User
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	var user models.User
	database.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})

		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})

		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})

		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600 * 24 * 30, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "User login successfully",
	})
}

func UpdateUser(c *gin.Context) {
	// Mendapatkan ID dari parameter URL
	idParam := c.Param("userId")

	// Konversi ID dari string ke int
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Mencari record berdasarkan ID
	var record models.User
	result := database.DB.First(&record, id)

	// Memeriksa apakah record ditemukan atau tidak
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Mendapatkan data baru dari body request
	var updatedData models.User
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data format"})
		return
	}

	// Memperbarui data pengguna
	database.DB.Model(&record).Updates(updatedData)

	c.JSON(http.StatusOK, gin.H{"message": "Record updated successfully"})
}

func DeleteUser(c *gin.Context) {
	// Mendapatkan ID dari parameter URL
	idParam := c.Param("userId")

	// Konversi ID dari string ke int
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Mencari record berdasarkan ID
	var record models.User
	result := database.DB.First(&record, id)

	// Memeriksa apakah record ditemukan atau tidak
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Menghapus record dari database
	database.DB.Delete(&record)

	c.JSON(http.StatusOK, gin.H{"message": "Record deleted successfully"})	
}