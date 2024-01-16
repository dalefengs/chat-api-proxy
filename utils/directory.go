package utils

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"errors"
	"github.com/dalefengs/chat-api-proxy/global"
	"io"
	"os"

	"go.uber.org/zap"
)

//@author: [likfees](https://github.com/dalefengs)
//@function: PathDirExists
//@description: 文件目录是否存在
//@param: path string
//@return: bool, error

func PathDirExists(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err == nil {
		if fi.IsDir() {
			return true, nil
		}
		return false, errors.New("存在同名文件")
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func PathExists(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err == nil {
		if fi.IsDir() {
			return true, nil
		}
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//@author: [likfees](https://github.com/dalefengs)
//@function: CreateDir
//@description: 批量创建文件夹
//@param: dirs ...string
//@return: err error

func CreateDir(dirs ...string) (err error) {
	for _, v := range dirs {
		exist, err := PathDirExists(v)
		if err != nil {
			return err
		}
		if !exist {
			global.Log.Debug("create directory" + v)
			if err := os.MkdirAll(v, os.ModePerm); err != nil {
				global.Log.Error("create directory"+v, zap.Any(" error:", err))
				return err
			}
		}
	}
	return err
}

// ReadGzipFile 读取 gzip 文件内容
func ReadGzipFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	gz, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}
	defer gz.Close()

	var content []byte
	reader := bufio.NewReader(gz)
	for {
		line, err := reader.ReadBytes('\n')
		if line != nil {
			content = append(content, line...)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
	}
	return content, nil
}

// ZipMarshalAndWriteToFile JSON序列化后压缩并写入文件
func ZipMarshalAndWriteToFile(filePath string, data interface{}) error {
	// 先进行gzip压缩
	zippedData, err := MarshalJsonAndGzip(data)
	if err != nil {
		return err
	}
	// 将压缩后的数据写入到指定路径的文件中
	return os.WriteFile(filePath, zippedData, 0644)
}

// ZipAndWriteToFile 压缩并写入文件
func ZipAndWriteToFile(filePath string, data []byte) error {
	// 先进行gzip压缩
	zippedData, err := GzipEncode(data)
	if err != nil {
		return err
	}
	// 将压缩后的数据写入到指定路径的文件中
	return os.WriteFile(filePath, zippedData, 0644)
}

// GzipEncode gzip 压缩
func GzipEncode(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)
	defer gzipWriter.Close()
	_, err := gzipWriter.Write(data)
	if err != nil {
		return nil, errors.New("gzip write error: " + err.Error())
	}

	if err := gzipWriter.Close(); err != nil {
		return nil, errors.New("gzip close error: " + err.Error())
	}
	return buf.Bytes(), nil
}

func GzipDecode(input []byte) ([]byte, error) {
	// 创建一个新的 gzip.Reader
	bytesReader := bytes.NewReader(input)
	gzipReader, err := gzip.NewReader(bytesReader)
	defer gzipReader.Close()
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	// 从 Reader 中读取出数据
	if _, err := buf.ReadFrom(gzipReader); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// MarshalJsonAndGzip 压缩
func MarshalJsonAndGzip(data interface{}) ([]byte, error) {
	marshalData, err := global.Json.Marshal(data)
	if err != nil {
		return nil, err
	}
	gzipData, err := GzipEncode(marshalData)
	if err != nil {
		return nil, err
	}
	return gzipData, err
}
