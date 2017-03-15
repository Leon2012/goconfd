package util

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func IsWritable(folder string) (err error) {
	fileInfo, err := os.Stat(folder)
	if err != nil {
		return err
	}
	if !fileInfo.IsDir() {
		return errors.New("Not a valid folder!")
	}
	perm := fileInfo.Mode().Perm()
	if 0200&perm != 0 {
		return nil
	}
	return errors.New("Not writable!")
}

func IsExist(name string) bool {
	_, err := os.Stat(name)
	fmt.Println(err)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func GetFileList(dir string, ext string) ([]string, error) {
	files := []string{}
	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, f := range fs {
		if !f.IsDir() {
			fileName := f.Name()
			if strings.HasSuffix(fileName, ext) {
				fileName = filepath.Join(dir, fileName)
				files = append(files, fileName)
			}
		}
	}
	return files, nil
}

func GetFileContent(file string) ([]byte, error) {
	return ioutil.ReadFile(file)
}

func GetName(fileName string) string {
	pathSeparatorIdx := strings.LastIndex(fileName, "/")
	if pathSeparatorIdx == -1 {
		return fileName
	}
	baseFileName := fileName[(pathSeparatorIdx + 1):len(fileName)]
	pointSeparatorIdx := strings.LastIndex(baseFileName, ".")
	if pointSeparatorIdx == -1 {
		return baseFileName
	}
	name := baseFileName[0:pointSeparatorIdx]
	return name
}
