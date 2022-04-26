package main_test

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"testing"
)

func TestMain(m *testing.M) {
	cmd := exec.Command(
		"./easy-proxy-service",
		"start",
		"-ip",
		"127.0.0.1",
		"-port",
		"13191",
		"-user-ip",
		"1.2.3.4",
		"-user-port",
		"6666",
		"-cert",
		"/root/tmp-easy-proxy/tmp-ca/www_baidu_comcert.pem",
		"-key",
		"/root/tmp-easy-proxy/tmp-ca/www_baidu_comkey.pem",
	)
	cmd.Stdout = os.Stdout
	err := cmd.Start()
	if err != nil {
		panic(fmt.Sprintf("发生了错误, err: %s", err.Error()))
	}
	fmt.Println("这里的 pid 号是: ", cmd.Process.Pid)
	err = syscall.Kill(cmd.Process.Pid, syscall.SIGINT)
	if err != nil {
		fmt.Println("杀死进程失败")
	} else {
		fmt.Println("杀死进程成功")
	}
}
