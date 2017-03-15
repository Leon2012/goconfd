package util

import (
	"testing"
)

func TestIsExist(t *testing.T) {
	fileName := "/go/src/github.com/Leon2012/goconfd/README11.md"
	exist := IsExist(fileName)
	if exist {
		t.Log("exist")
	} else {
		t.Log("no exist")
	}
}

func TestGetFileList(t *testing.T) {
	dir := "/Users/pengleon/Downloads/goconfd"
	ext := "php"
	files, err := GetFileList(dir, ext)
	if err != nil {
		t.Error(err)
	}
	t.Log(files)
}

func TestGetName(t *testing.T) {
	fileName := "/Users/pengleon/Downloads/goconfd/646576656c6f702e61637469766974792e64616368752e61637439392e6964.php"
	newFileName := GetName(fileName)
	t.Log(newFileName)
}
