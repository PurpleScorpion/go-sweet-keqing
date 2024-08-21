package keqing

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func MkdirAll(path string) {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		panic(err)
	}
}

/*
文件存在判断
*/
func FileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

/*
获取一个目录下的所有文件路径
*/
func ReadDirFiles(dir string) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}

/*
获取指定路径下md5值和文件名称
*/
func FileMD5(filePath string) (string, error) {
	if !FileExists(filePath) {
		return "", errors.New("file not exists")
	}
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 创建MD5哈希对象
	hash := md5.New()

	// 将文件内容读取到哈希对象中
	_, err = io.Copy(hash, file)
	if err != nil {
		return "", err
	}

	// 计算MD5值
	md5sum := hash.Sum(nil)

	// 将MD5值转换为十六进制字符串
	md5String := fmt.Sprintf("%x", md5sum)

	return md5String, nil
}

func SaveFile(path, data string) error {
	if err := os.WriteFile(path, []byte(data), 0644); err != nil {
		return err
	}
	return nil
}

/*
保存文件, 使用字符串路径保存
*/
func CopyFile(filePath, savePath string) error {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	dst, err := os.Create(savePath)
	defer dst.Close()
	if err != nil {
		return err
	}
	_, err = io.Copy(dst, file)
	if err != nil {
		return err
	}
	return nil
}

/*
保存文件,直接使用io流保存
*/
func CopyFile2IO(dst io.Writer, src io.Reader) error {
	if dst == nil || src == nil {
		return errors.New("dst or src is nil")
	}
	_, err := io.Copy(dst, src)
	if err != nil {
		return err
	}
	return nil
}

/*
删除指定路径的文件
*/
func RemoveFile(path string) error {
	err := os.Remove(path)
	if err != nil {
		return fmt.Errorf("cannot delete file: %w", err)
	}
	return nil
}

/*
删除整个目录及其所有内容
*/
func RemoveDir(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		return fmt.Errorf("Cannot delete directory: %w", err)
	}
	return nil
}

func ReNameFile(oldPath, newPath string) error {
	err := os.Rename(oldPath, newPath)
	if err != nil {
		return fmt.Errorf("Unable to rename file: %w", err)
	}
	return nil
}
func HomePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return home
}

func FileSize(path string) int64 {
	if !FileExists(path) {
		return 0
	}
	file, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return file.Size()
}
