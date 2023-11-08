package controllers

import (
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