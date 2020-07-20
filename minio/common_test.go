package minio

import "testing"

func TestName(t *testing.T) {
	minioConfig := &MinioConfig{
		Endpoint:        "minio.dev.segma.tech",
		AccessKeyID:     "AKIAIOSFODNN7EXAMPLESCDI",
		SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
		UseSSL:          false,
	}
	_, err := CreateMinioClient(minioConfig)
	t.Log(err)
}
