package main

import (
	"btpn-syariah-final-project/controllers"
	"btpn-syariah-final-project/database"
	"btpn-syariah-final-project/middlewares"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/olahol/go-imageupload"
)

func init() {
	database.LoadEnvVariables()
	database.ConnectToDb()
	database.SyncDatabase()
}

var currentImage *imageupload.Image

func main() {
	r := gin.Default()

	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.GET("/validate", middlewares.RequireAuth, controllers.Validate)

	// r.GET("/image", func (c *gin.Context)  {
	// 	if currentImage == nil {
	// 		c.AbortWithStatus(http.StatusNotFound)
	// 		return
	// 	}
	// 	currentImage.Write(c.Writer)
	// })

	// r.GET("/thumbnail", func(c *gin.Context) {
	// 	if currentImage == nil {
	// 		c.AbortWithStatus(http.StatusNotFound)
	// 		return
	// 	}

	// 	t, err := imageupload.ThumbnailJPEG(currentImage, 300, 300, 80)

	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	t.Write(c.Writer)
	// })

	// r.POST("/upload", func(c *gin.Context) {
	// 	img, err := imageupload.Process(c.Request, "image")

	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	currentImage = img
	// })

	r = gin.Default()

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

	r.Run()
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