package oss

import (
	"time"
	"tool/global/variable"
)

// 上传文件到 OSS
// filename: 文件名
// filePath: 文件路径
func UploadFile(filename, filePath string) (string, error) {

	client, _ := getClient("Aliyun")

	//获取上传目录
	dir := variable.ConfigYml.GetString("Oss.Aliyun.Dir")

	// 生成文件名 年月日小时分秒_文件名
	ossFileName := time.Now().Format("20060102150405") + "_" + filename

	// 上传文件到 OSS
	fileName, err := client.UploadFileFromPath(dir+"/"+ossFileName, filePath)
	if err != nil {
		return "", err
	}

	return fileName, nil
}
