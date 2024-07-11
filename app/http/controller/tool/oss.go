package tool

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"tool/pkg/oss"

	"github.com/gin-gonic/gin"
)

func Upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload error"})
		return
	}

	// 保存文件到临时目录
	tempFilePath := filepath.Join(os.TempDir(), header.Filename)
	out, err := os.Create(tempFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving file"})
		return
	}
	defer out.Close()

	fileContent, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading file"})
		return
	}

	_, err = out.Write(fileContent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error writing file"})
		return
	}

	ossConfig := oss.OssConfig{
		Endpoint:        "",
		AccessKeyID:     "",
		AccessKeySecret: "",
		BucketName:      "",
	}

	ossClient, _ := oss.NewOSSClient(ossConfig)

	// 上传文件到 OSS
	fileName, err := ossClient.UploadFileFromPath(header.Filename, tempFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 删除临时文件
	os.Remove(tempFilePath)

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "file_name": fileName})
}
