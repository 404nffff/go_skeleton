package tool

import (
	"image"
	"image/jpeg"
	_ "image/png" // 支持PNG格式
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"tool/app/global/variable"
	"tool/app/utils/common"
	"tool/app/utils/oss"

	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
)

func Upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload error"})
		return
	}

	//获取配置 检测文件大小
	if header.Size/1024/1024 > int64(variable.ConfigYml.GetInt("UploadFile.MaxSize")) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds the limit"})
		return
	}

	//检测文件后缀
	ext := filepath.Ext(header.Filename)

	//ext 去除点
	ext = strings.TrimPrefix(ext, ".")

	allowExt := variable.ConfigYml.GetString("UploadFile.AllowExt")

	//分割字符串
	allowExtSlice := strings.Split(allowExt, ",")

	//判断是否存在
	if !common.InArray(ext, allowExtSlice) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File suffix is not allowed"})
		return
	}

	//检测文件类型
	mime := header.Header.Get("Content-Type")

	allowMime := variable.ConfigYml.GetString("UploadFile.AllowMime")

	//分割字符串
	allowMimeSlice := strings.Split(allowMime, ",")
	if !common.InArray(mime, allowMimeSlice) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File type is not allowed"})
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

	// 打开临时文件进行压缩
	tempFile, err := os.Open(tempFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error opening file"})
		return
	}
	defer tempFile.Close()

	img, _, err := image.Decode(tempFile)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error decoding image"})
		return
	}

	// 获取配置的压缩参数
	maxWidth := variable.ConfigYml.GetInt("UploadFile.ResizeWidth")
	maxHeight := variable.ConfigYml.GetInt("UploadFile.ResizeHeight")
	jpegQuality := variable.ConfigYml.GetInt("UploadFile.JPEGQuality")

	// 获取原始图片的宽高
	originalWidth := img.Bounds().Dx()
	originalHeight := img.Bounds().Dy()

	// 计算新的宽高
	newWidth := originalWidth
	newHeight := originalHeight

	if originalWidth > maxWidth || originalHeight > maxHeight {
		aspectRatio := float64(originalWidth) / float64(originalHeight)
		if originalWidth > maxWidth {
			newWidth = maxWidth
			newHeight = int(float64(maxWidth) / aspectRatio)
		}
		if newHeight > maxHeight {
			newHeight = maxHeight
			newWidth = int(float64(maxHeight) * aspectRatio)
		}
	}

	// 压缩图片
	resizedImg := resize.Resize(uint(newWidth), uint(newHeight), img, resize.Lanczos3)

	// 将压缩后的图片保存到新临时文件
	compressedFilePath := filepath.Join(os.TempDir(), "compressed_"+header.Filename)
	outCompressed, err := os.Create(compressedFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating compressed file"})
		return
	}
	defer outCompressed.Close()

	var opt jpeg.Options
	opt.Quality = jpegQuality // 设置图片质量
	err = jpeg.Encode(outCompressed, resizedImg, &opt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error encoding compressed image"})
		return
	}

	// 上传压缩后的文件到 OSS
	fileName, err := oss.UploadFile("compressed_"+header.Filename, compressedFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 删除临时文件
	os.Remove(tempFilePath)
	os.Remove(compressedFilePath)

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded and compressed successfully", "file_name": fileName})
}
