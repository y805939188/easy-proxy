package tools_test

import (
	"easy-proxy/tools"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestPath(t *testing.T) {
	fmt.Println(filepath.Abs(filepath.Dir(os.Args[0])))
	return
	p, err := tools.GetTmpCaPath()
	if err != nil {
		panic("获取 tmp ca path 失败, err: " + err.Error())
	}
	fmt.Println("ca 临时路径是: ", p)
}
