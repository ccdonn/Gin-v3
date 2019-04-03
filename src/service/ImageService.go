package service

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"../constant"
	ApiErr "../error"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetImage : get image from file
func GetImage(c *gin.Context) {
	name := c.Param("name")

	// add more checker for security issue
	if _, err := os.Stat(constant.ImagePath + name); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, ApiErr.ErrNotFound)
		return
	}

	c.File(constant.ImagePath + name)
	return
}

// UploadImage : user upload images to server
func UploadImage(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["upload"]

	uuid, err := uuid.NewUUID()
	if err != nil {
		c.AbortWithStatusJSON(500, ApiErr.ErrUUIDGen)
		return
	}

	fmt.Println(uuid, files, files[0].Filename, files[0].Header, files[0].Size)
	fmt.Println(strconv.FormatInt(files[0].Size, 10))

	file, err := files[0].Open()
	if err != nil {
		c.AbortWithStatusJSON(500, ApiErr.ErrFileOpen)
		return
	}

	data := make([]byte, files[0].Size)
	count, err := file.Read(data)
	if err != nil {
		c.AbortWithStatusJSON(500, ApiErr.ErrFileRead)
		return
	}

	if count <= 0 {
		c.AbortWithStatusJSON(500, ApiErr.ErrFileSize)
		return
	}

	// check folder before write file

	err = ioutil.WriteFile(constant.ImagePath+uuid.String()+".jpg", data, 0644)
	if err != nil {
		c.AbortWithStatusJSON(500, ApiErr.ErrFileWrite)
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
	})

	return
}
