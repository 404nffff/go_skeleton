package oss

import (
	"tool/app/global/variable"
	"tool/pkg/oss"
)

var client *oss.OSSClient

func init() {
	ossConfig := oss.OssConfig{
		Endpoint:        variable.ConfigYml.GetString("Oss.Endpoint"),
		AccessKeyID:     variable.ConfigYml.GetString("Oss.AccessKeyId"),
		AccessKeySecret: variable.ConfigYml.GetString("Oss.AccessKeySecret"),
		BucketName:      variable.ConfigYml.GetString("Oss.BucketName"),
	}

	client, _ = oss.NewOSSClient(ossConfig)
}

// 上传文件到 OSS
func uploadFile(filename, tempFilePath string) (string, error) {

	// 上传文件到 OSS
	fileName, err := client.UploadFileFromPath(filename, tempFilePath)
	if err != nil {
		return "", err
	}

	return fileName, nil
}
