package oss

// 上传文件到 OSS
func UploadFile(filename, tempFilePath string) (string, error) {

	client, _ := getClient("Aliyun")

	// 上传文件到 OSS
	fileName, err := client.UploadFileFromPath(filename, tempFilePath)
	if err != nil {
		return "", err
	}

	return fileName, nil
}
