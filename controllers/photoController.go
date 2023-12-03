package controllers

import (
	"btpn-syariah-final-project/database"
	"btpn-syariah-final-project/models"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"btpn-syariah-final-project/helpers"

	"github.com/gin-gonic/gin"
)

func GetPhoto(c *gin.Context) {
	user, err := helpers.GetUserLogin(c)
	
	var photo models.Photo
	database.DB.First(&photo, "user_id = ?", user.ID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "photo not found",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"userPhoto": &photo,
	})

	// files, err := getUploadedFiles()
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": err.Error(),
	// 	})
	// 	return
	// }
}

func StorePhoto(c *gin.Context) {
	user, err := helpers.GetUserLogin(c)

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var body struct {
		models.Photo
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	currentTime := time.Now()

	formattedTime := currentTime.Format("20060102150405")	

	photo := models.Photo{
    UserID: user.ID,
    Title: c.Request.FormValue("title"),
    Caption: c.Request.FormValue("caption"),
    PhotoUrl: "uploads/" + formattedTime + "_" + file.Filename,
  }

	result := database.DB.Create(&photo)

	if result.Error!= nil {
		c.JSON(http.StatusBadRequest, gin.H{
      "error": "Failed to store photo",
    })

		return
	}

	// Simpan file ke server
	err = c.SaveUploadedFile(file, "uploads/" + formattedTime + "_" + file.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Foto profil berhasil ditambahkan"),
	})
}

func UpdatePhoto(c *gin.Context) {
	// Mendapatkan ID dari parameter URL
	idParam := c.Param("photoId")

	file, err := c.FormFile("file")

	currentTime := time.Now()

	formattedTime := currentTime.Format("20060102150405")	

	// Konversi ID dari string ke int
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Mencari photo berdasarkan ID
	var photo models.Photo
	result := database.DB.First(&photo, id)
	
	var photoUrl string
	if file == nil {
		photoUrl = photo.PhotoUrl
	} else {
		photoUrl = "uploads/" + formattedTime + "_" + file.Filename
	}

	// Memeriksa apakah photo ditemukan atau tidak
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Photo user not found"})
		return
	}

	// Memperbarui data pengguna	

	photo.Title = c.Request.FormValue("title")
	photo.Caption = c.Request.FormValue("caption")
	photo.PhotoUrl = photoUrl
	
	database.DB.Save(&photo)	

	// Simpan file ke server
	err = c.SaveUploadedFile(file, "uploads/" + formattedTime + "_" + file.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Photo user updated successfully"})
}

func DeletePhoto(c *gin.Context) {
	// Mendapatkan ID dari parameter URL
	idParam := c.Param("photoId")

	// Konversi ID dari string ke int
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Mencari photo berdasarkan ID
	var photo models.Photo
	result := database.DB.First(&photo, id)

	// Memeriksa apakah record ditemukan atau tidak
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Photo user not found"})
		return
	}

	// Menghapus record dari database
	database.DB.Delete(&photo)

	err = deleteFile(photo.PhotoUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	 
	c.JSON(http.StatusOK, gin.H{"message": "Photo user deleted successfully"})	
}

// Fungsi untuk menghapus file
func deleteFile(filename string) error {
	err := os.Remove(filename)
	return err
}