package service

import (
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetImage : get image from file
func GetImage(c *gin.Context) {
	name := c.Param("name")
	// add more checker for security issue
	c.File(name)
	return
}

// UploadImage : user upload images to server
func UploadImage(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["upload"]

	uuid, err := uuid.NewUUID()
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"status":       "failure",
			"errorCode":    10008000,
			"errorMessage": "uuid generate fail",
		})
		return
	}

	fmt.Println(uuid, files, files[0].Filename, files[0].Header, files[0].Size)
	fmt.Println(strconv.FormatInt(files[0].Size, 10))

	file, err := files[0].Open()
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"status":       "failure",
			"errorCode":    10008001,
			"errorMessage": "file upload fail",
		})
		return
	}

	data := make([]byte, files[0].Size)
	count, err := file.Read(data)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"status":       "failure",
			"errorCode":    10008002,
			"errorMessage": "file upload fail",
		})
		return
	}

	if count <= 0 {
		c.AbortWithStatusJSON(500, gin.H{
			"status":       "failure",
			"errorCode":    10008003,
			"errorMessage": "file upload fail(file size?)",
		})
		return
	}

	// check folder before write file

	err = ioutil.WriteFile("../../images/"+uuid.String()+".jpg", data, 0644)
	if err != nil {
		panic(err)
	}

	c.JSON(200, gin.H{
		"status": "success",
	})

	return
}
