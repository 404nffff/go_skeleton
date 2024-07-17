package oss

import (
	"fmt"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type OSSClient struct {
	Bucket *oss.Bucket
}

// Aliyun:
//     Open: true  # 是否启用 OSS
//     Endpoint: "oss-cn-hangzhou.aliyuncs.com"
//     AccessKeyId: "xxxx"
//     AccessKeySecret: "xxx"
//     BucketName: "xxx"
//     Dir: "uploads"
//     ConnectTimeout: 10 # HTTP超时时间，单位为秒。默认值为10秒，0表示不超时。
//     ReadWriteTimeout: 60 # HTTP读取或写入超时时间，单位为秒。默认值为20秒，0表示不超时。

// oss 配置
type OssConfig struct {
	Endpoint         string // oss endpoint
	AccessKeyID      string // oss access key id
	AccessKeySecret  string // oss access key secret
	BucketName       string // oss bucket name
	Dir              string // 目录
	ConnectTimeout   int    // 超时时间
	ReadWriteTimeout int    // 读写超时时间
}

// 创建 OSS 客户端
func NewOSSClient(config OssConfig) (*OSSClient, error) {

	// 设置HTTP连接超时时间为20秒，HTTP读取或写入超时时间为60秒。
	time := oss.Timeout(int64(config.ConnectTimeout), int64(config.ReadWriteTimeout))

	client, err := oss.New(config.Endpoint, config.AccessKeyID, config.AccessKeySecret, time)

	if err != nil {
		return nil, fmt.Errorf("error creating OSS client: %v", err)
	}

	bucket, err := client.Bucket(config.BucketName)
	if err != nil {
		return nil, fmt.Errorf("error getting bucket: %v", err)
	}

	return &OSSClient{Bucket: bucket}, nil
}

// 上传文件
// fileName: 文件名
// filePath: 文件路径
func (o *OSSClient) UploadFileFromPath(fileName, filePath string) (string, error) {

	err := o.Bucket.PutObjectFromFile(fileName, filePath)
	if err != nil {
		return "", fmt.Errorf("error uploading file: %v", err)
	}

	return fileName, nil
}

// 删除文件
func (o *OSSClient) DeleteFile(fileName string) error {
	err := o.Bucket.DeleteObject(fileName)
	if err != nil {
		return fmt.Errorf("error deleting file: %v", err)
	}
	return nil
}
