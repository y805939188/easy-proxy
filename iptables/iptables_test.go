package iptables_test

import (
	"easy-proxy/iptables"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MyTestData struct {
	Name string
}

type MyTestResponse struct {
	Code int
	Data MyTestData
}

func helloServer(w http.ResponseWriter, req *http.Request) {
	data := MyTestData{Name: "ding"}
	res := new(MyTestResponse)
	res.Data = data
	resJson, _ := json.Marshal(res)
	io.WriteString(w, string(resJson))
}

func createTestServer(addr *string) (func() error, error) {
	srv := http.Server{
		Addr:    *addr,
		Handler: http.HandlerFunc(helloServer),
	}
	go srv.ListenAndServe()

	closeFunc := func() error {
		err := srv.Close()
		if err != nil {
			return err
		}
		return nil
	}
	return closeFunc, nil
}

func TestIptables(t *testing.T) {
	test := assert.New(t)

	iclient, err := iptables.GetNewIptablesClient()
	if err != nil {
		fmt.Println("创建 iptables 服务发生了错误: ", err.Error())
		fmt.Println("关闭服务失败: ", err.Error())
		return
	}

	// err = iclient.DeleteIptablesOutputRule("1.2.3.4", "", "127.0.0.1", "13190")
	// if err != nil {
	// 	fmt.Println("删除规则失败: ", err.Error())
	// 	return
	// }
	// return

	var testRoute = "127.0.0.1:13190"
	close, err := createTestServer(&testRoute)
	if err != nil {
		fmt.Println("创建测试用服务发生错误: ", err.Error())
		return
	}

	deleteRuleFunc, err := iclient.SetIptablesOutputRule("1.2.3.4", "", "127.0.0.1", "13190")
	if err != nil {
		fmt.Println("设置 iptables rule 发生错误: ", err.Error())
		err = close()
		fmt.Println("关闭服务失败: ", err.Error())
		return
	}

	resp, err := http.Get("http://1.2.3.4")
	if err != nil {
		fmt.Println("向 13190 发送请求失败: ", err.Error())
		err = close()
		fmt.Println("关闭服务失败: ", err.Error())
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取 13190 的 body 失败: ", err.Error())
		err = close()
		fmt.Println("关闭服务失败: ", err.Error())
		return
	}
	var myResponseBody MyTestResponse
	json.Unmarshal(body, &myResponseBody)
	test.Equal(myResponseBody.Data.Name, "ding")

	err = deleteRuleFunc()
	if err != nil {
		fmt.Println("删除 iptables rule 发生错误: ", err.Error())
		err = close()
		fmt.Println("关闭服务失败: ", err.Error())
		return
	}
}
