package compress

import (
	"archive/tar"
	"archive/zip"
	"compress/bzip2"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
)

func Unzip(zipFile, destDir string) error {
	zipReader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	for _, f := range zipReader.File {
		fpath := filepath.Join(destDir, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}

			inFile, err := f.Open()
			if err != nil {
				return err
			}
			defer inFile.Close()

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, inFile)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func UnTar(tarFile, destDir string) error {
	err := unTarFile(tarFile, destDir)
	if err == nil {
		return nil
	}
	logrus.WithError(err).Info("compress: untar file failed")
	err = unTarGz(tarFile, destDir)
	if err == nil {
		return nil
	}
	logrus.WithError(err).Info("compress: untar gzip failed")
	err = unTarBzip(tarFile, destDir)
	if err == nil {
		return nil
	}
	logrus.WithError(err).Info("compress: untar bzip failed")
	err = unTarFlate(tarFile, destDir)
	if err == nil {
		return nil
	}
	logrus.WithError(err).Info("compress: untar flate failed")
	err = unTarZlib(tarFile, destDir)
	if err != nil {
		logrus.WithError(err).Info("compress: untar zlib failed")
	}
	return err
}

func unTarFile(tarFile, destDir string) error {
	// 打开准备解压的 tar 包
	fileReader, err := os.Open(tarFile)
	if err != nil {
		return err
	}
	defer fileReader.Close()

	return unTar(fileReader, destDir)
}

func unTarBzip(tarFile, destDir string) error {
	// 打开准备解压的 tar 包
	fileReader, err := os.Open(tarFile)
	if err != nil {
		return err
	}
	defer fileReader.Close()

	//将打开的文件先解压
	bzip2Reader := bzip2.NewReader(fileReader)

	return unTar(bzip2Reader, destDir)
}

func unTarFlate(tarFile, destDir string) error {
	// 打开准备解压的 tar 包
	fileReader, err := os.Open(tarFile)
	if err != nil {
		return err
	}
	defer fileReader.Close()

	//将打开的文件先解压
	flateReader := flate.NewReader(fileReader)
	defer flateReader.Close()

	return unTar(flateReader, destDir)
}

func unTarGz(tarFile, destDir string) error {
	// 打开准备解压的 tar 包
	fileReader, err := os.Open(tarFile)
	if err != nil {
		return err
	}
	defer fileReader.Close()

	//将打开的文件先解压
	gzReader, err := gzip.NewReader(fileReader)
	if err != nil {
		logrus.WithError(err).Info("common: not a gzip type")
		return err
	}
	defer gzReader.Close()

	return unTar(gzReader, destDir)
}

func unTarZlib(tarFile, destDir string) error {
	// 打开准备解压的 tar 包
	fileReader, err := os.Open(tarFile)
	if err != nil {
		return err
	}
	defer fileReader.Close()

	//将打开的文件先解压
	zlibReader, err := zlib.NewReader(fileReader)
	if err != nil {
		logrus.WithError(err).Info("common: not a gzip type")
		return err
	}
	defer zlibReader.Close()

	return unTar(zlibReader, destDir)
}

func unTar(fileReader io.Reader, destDir string) error {
	// 通过 gzReader 创建 tar.Reader
	tarReader := tar.NewReader(fileReader)

	// 现在已经获得了 tar.Reader 结构了，只需要循环里面的数据写入文件就可以了
	for {
		hdr, err := tarReader.Next()

		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case hdr == nil:
			continue
		}

		// 处理下保存路径，将要保存的目录加上 header 中的 Name
		// 这个变量保存的有可能是目录，有可能是文件，所以就叫 FileDir 了……
		dstFileDir := filepath.Join(destDir, hdr.Name)

		// 根据 header 的 Typeflag 字段，判断文件的类型
		switch hdr.Typeflag {
		case tar.TypeDir: // 如果是目录时候，创建目录
			// 判断下目录是否存在，不存在就创建
			if b := existDir(dstFileDir); !b {
				// 使用 MkdirAll 不使用 Mkdir ，就类似 Linux 终端下的 mkdir -p，
				// 可以递归创建每一级目录
				if err := os.MkdirAll(dstFileDir, 0775); err != nil {
					return err
				}
			}
		case tar.TypeReg: // 如果是文件就写入到磁盘
			// 创建一个可以读写的文件，权限就使用 header 中记录的权限
			// 因为操作系统的 FileMode 是 int32 类型的，hdr 中的是 int64，所以转换下
			file, err := os.OpenFile(dstFileDir, os.O_CREATE|os.O_RDWR, os.FileMode(hdr.Mode))
			if err != nil {
				return err
			}
			n, err := io.Copy(file, tarReader)
			if err != nil {
				return err
			}
			// 将解压结果输出显示
			logrus.Info("compress: uncompress success： ", "destDir", " handle ", n, " chars")

			// 不要忘记关闭打开的文件，因为它是在 for 循环中，不能使用 defer
			// 如果想使用 defer 就放在一个单独的函数中
			file.Close()
		}
	}
}

// 判断目录是否存在
func existDir(dirname string) bool {
	fi, err := os.Stat(dirname)
	return (err == nil || os.IsExist(err)) && fi.IsDir()
}
