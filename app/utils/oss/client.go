package oss

import (
	"sync"
	"tool/app/global/variable"
	"tool/pkg/oss"
)

var (
	clients = make(map[string]*oss.OSSClient)
	lock    sync.Mutex
)

// 初始化 OSS 客户端
func getClient(name string) (*oss.OSSClient, error) {
	lock.Lock()
	defer lock.Unlock()

	if client, ok := clients[name]; ok {
		return client, nil
	}

	//oss 前缀
	ossPrefix := "Oss." + name + "."

	ossConfig := oss.OssConfig{
		Endpoint:         variable.ConfigYml.GetString(ossPrefix + "Endpoint"),
		AccessKeyID:      variable.ConfigYml.GetString(ossPrefix + "AccessKeyId"),
		AccessKeySecret:  variable.ConfigYml.GetString(ossPrefix + "AccessKeySecret"),
		BucketName:       variable.ConfigYml.GetString(ossPrefix + "BucketName"),
		ConnectTimeout:   variable.ConfigYml.GetInt(ossPrefix + "ConnectTimeout"),
		ReadWriteTimeout: variable.ConfigYml.GetInt(ossPrefix + "ReadWriteTimeout"),
	}

	client, err := oss.NewOSSClient(ossConfig)
	if err != nil {
		return nil, err
	}

	clients[name] = client
	return client, nil
}
