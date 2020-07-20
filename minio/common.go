package minio

import (
	"github.com/minio/minio-go/v6"
	"github.com/sirupsen/logrus"
	"path/filepath"
)

type MinioConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
}

type FileInfo struct {
	BucketName string
	ObjectName string
	FileName   string
	FilePath   string
	Unzip      bool
}

func CreateMinioClient(config *MinioConfig) (*minio.Client, error) {
	return minio.New(config.Endpoint, config.AccessKeyID, config.SecretAccessKey, config.UseSSL)
}

func FGetObject(client *minio.Client, fileInfo *FileInfo) (string, error) {
	filePath := filepath.Join("/tmp/", fileInfo.FileName)
	logrus.Debug(`common: file path is '`, filePath, `'`)
	return filePath, client.FGetObject(fileInfo.BucketName, fileInfo.ObjectName, filePath, minio.GetObjectOptions{})
}
