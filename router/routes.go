package router

import (
	"btpn-syariah-final-project/controllers"
	"btpn-syariah-final-project/middlewares"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.GET("/validate", middlewares.RequireAuth, controllers.Validate)

	// Upload file
	r.POST("/upload", func(c *gin.Context) {
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
	})

	// Mengambil daftar file yang sudah diunggah
	r.GET("/files", func(c *gin.Context) {
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
	})

	// Hapus file
	r.DELETE("/files/:filename", func(c *gin.Context) {
		filename := c.Param("filename")
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
	})

	return r
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