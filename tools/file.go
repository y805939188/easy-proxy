package tools

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func GetFileNameAndExt(filePath string) (string, string) {
	fileName := filepath.Base(filePath)
	ext := path.Ext(fileName)
	fileName = strings.Replace(fileName, ext, "", 1)
	return fileName, ext
}

func ReadFile(filePath string) ([]byte, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func GetAllFile(pathname string) ([]string, error) {
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		return nil, err
	}
	s := []string{}

	for _, fi := range rd {
		if !fi.IsDir() {
			fullName := pathname + "/" + fi.Name()
			s = append(s, fullName)
		}
	}
	return s, nil
}

func CreateFile(filename string, data []byte, perm os.FileMode) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}

func CreateDir(dirName string) error {
	err := os.MkdirAll(dirName, 0766)
	if err != nil {
		return err
	}
	return nil
}

func FileIsExisted(filename string) bool {
	existed := true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		existed = false
	}
	return existed
}

func DeleteFile(path string) error {
	if path == "" {
		return nil
	}
	if FileIsExisted(path) {
		err := os.Remove(path)
		if err != nil {
			return err
		}
	}
	return nil
}

func FileCopy(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}
