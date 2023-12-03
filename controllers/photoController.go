package controllers

import (
	"btpn-syariah-final-project/database"
	"btpn-syariah-final-project/models"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func GetPhoto(c *gin.Context) {
	files, err := getUploadedFiles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"files": files,
	})
}

func StorePhoto(c *gin.Context) {
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

	photo := models.Photo{
    UserID: 1,
    Title: c.Request.FormValue("title"),
    Caption: c.Request.FormValue("caption"),
    PhotoUrl: "uploads/"+file.Filename,
  }

	result := database.DB.Create(&photo)

	if result.Error!= nil {
		c.JSON(http.StatusBadRequest, gin.H{
      "error": "Failed to store photo",
    })

		return
	}

	// Simpan file ke server
	err = c.SaveUploadedFile(file, "uploads/"+file.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("File %s berhasil diunggah!", file.Filename),
	})
}

func UpdatePhoto(c *gin.Context) {
	// file, err := c.FormFile("file")
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error": err.Error(),
	// 	})
	// 	return
	// }

	// // Simpan file ke server
	// err = c.SaveUploadedFile(file, "uploads/"+file.Filename)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": err.Error(),
	// 	})
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("File berhasil diupdate!"),
	})
}

func DeletePhoto(c *gin.Context) {
	filename := c.Param("photoId")
	err := deleteFile(filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("File %s berhasil dihapus!", filename),
	})
}

// Fungsi untuk mendapatkan daftar file yang sudah diunggah
func getUploadedFiles() ([]string, error) {
	files, err := os.ReadDir("uploads")
	if err != nil {
		return nil, err
	}

	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}

	return fileNames, nil
}

// Fungsi untuk menghapus file
func deleteFile(filename string) error {
	err := os.Remove("uploads/" + filename)
	return err
}