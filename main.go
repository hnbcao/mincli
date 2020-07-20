package main

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"go-minio/compress"
	"go-minio/minio"
	"os"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	app := cli.NewApp()
	app.Name = "minio client"
	app.Usage = "download file from minio"
	app.Action = run
	app.Version = "v0.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "endpoint",
			Usage:  "Minio Endpoint",
			EnvVar: "MINIO_ENDPOINT,ENDPOINT",
		},
		cli.StringFlag{
			Name:   "accessKeyID",
			Usage:  "Minio accessKeyID",
			EnvVar: "MINIO_ACCESS_KEY_ID,ACCESS_KEY_ID",
		},
		cli.StringFlag{
			Name:   "secretAccessKey",
			Usage:  "Minio secretAccessKey",
			EnvVar: "MINIO_SECRET_ACCESS_KEY,SECRET_ACCESS_KEY",
		},
		cli.BoolFlag{
			Name:   "useSSL",
			Usage:  "Minio useSSL",
			EnvVar: "MINIO_USE_SSL,USE_SSL",
		},
		cli.StringFlag{
			Name:   "bucketName",
			Usage:  "Minio bucketName",
			EnvVar: "MINIO_BUCKET_NAME,BUCKET_NAME",
		},
		cli.StringFlag{
			Name:   "objectName",
			Usage:  "Minio objectName",
			EnvVar: "MINIO_OBJECT_NAME,OBJECT_NAME",
		},
		cli.StringFlag{
			Name:   "fileName",
			Usage:  "Minio fileName",
			EnvVar: "MINIO_FILE_NAME,FILE_NAME",
		},
		cli.StringFlag{
			Name:   "filePath",
			Usage:  "Minio filePath",
			EnvVar: "MINIO_FILE_PATH,FILE_PATH",
			Value:  "/data/file/",
		},
		cli.BoolFlag{
			Name:   "unzip",
			Usage:  "if or not unzip",
			EnvVar: "MINIO_FILE_UNZIP,FILE_UNZIP",
		},
	}
	if err := app.Run(os.Args); err != nil {
		logrus.WithError(err).Fatal("application execute failed")
	}

}

func run(c *cli.Context) {
	minioConfig := &minio.MinioConfig{
		Endpoint:        c.String("endpoint"),
		AccessKeyID:     c.String("accessKeyID"),
		SecretAccessKey: c.String("secretAccessKey"),
		UseSSL:          c.Bool("useSSL"),
	}

	minioClient, err := minio.CreateMinioClient(minioConfig)

	if err != nil {
		logrus.WithError(err).Fatalln("main: init minio client error")
	}
	fileInfo := &minio.FileInfo{
		BucketName: c.String("bucketName"),
		ObjectName: c.String("objectName"),
		FileName:   c.String("fileName"),
		FilePath:   c.String("filePath"),
		Unzip:      c.Bool("unzip"),
	}

	zipPath, err := minio.FGetObject(minioClient, fileInfo)

	if err != nil {
		logrus.WithError(err).Fatalln("main: download file failed")
	}

	if fileInfo.Unzip {
		err = compress.Unzip(zipPath, fileInfo.FilePath)

		if err != nil {
			logrus.WithError(err).Info("main: unzip file failed")
			err = compress.UnTar(zipPath, fileInfo.FilePath)
			if err != nil {
				logrus.WithError(err).Fatalln("main: untar file failed")
			}
		}
	}

	logrus.Printf("main: download file success")
}
