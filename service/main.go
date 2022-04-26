package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/y805939188/dcommand"
)

func startHttpsService(userTargetIp, userTargetPort string, port string, certPath, keyPath string) error {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("ding test")
		io.WriteString(w, "Hello, TLS!\n")
	})

	err := http.ListenAndServeTLS(":"+port, certPath, keyPath, nil)
	fmt.Println("创建服务发生错误, err: ", err.Error())

	return err
}

func main() {
	cmd := &dcommand.DCommand{}
	cmd.Command("easy-proxy-service").
		Operator("start").
		Flag(&dcommand.FlagInfo{Name: "ip", Short: "ip"}).
		Flag(&dcommand.FlagInfo{Name: "port", Short: "port"}).
		Flag(&dcommand.FlagInfo{Name: "user-ip", Short: "user-ip"}).
		Flag(&dcommand.FlagInfo{Name: "user-port", Short: "user-port"}).
		Flag(&dcommand.FlagInfo{Name: "cert", Short: "cert"}).
		Flag(&dcommand.FlagInfo{Name: "key", Short: "key"}).
		Operator("kill").
		Flag(&dcommand.FlagInfo{Name: "pid", Short: "pid"}).
		Handler(func(command string, fc *dcommand.DCommand) error {
			_cmd := fc.GetCommandIfExist(command)
			currentOperator := ""
			for _, operator := range _cmd.Operators {
				if operator.Passed {
					if currentOperator != "" {
						return fmt.Errorf(fmt.Sprintf("无法同时设置 operator: %s & %s", currentOperator, operator.Name))
					}
					currentOperator = operator.Name
				}
			}
			if currentOperator == "" {
				return fmt.Errorf("easy-proxy-service 命令需要至少一个 operator: start | kill")
			}
			cmd := fc.GetCommandIfExist(command)
			op := fc.GetOperatorIfExistByCommand(currentOperator, cmd)
			switch currentOperator {
			case "start":
				ipFlag := fc.GetFlagIfExistInOperatorByOperator("ip", true, op)
				if !ipFlag.Passed {
					return fmt.Errorf("-ip 参数是必须的")
				}
				portFlag := fc.GetFlagIfExistInOperatorByOperator("port", true, op)
				if !portFlag.Passed {
					return fmt.Errorf("-port 参数是必须的")
				}
				userIpFlag := fc.GetFlagIfExistInOperatorByOperator("user-ip", true, op)
				if !ipFlag.Passed {
					return fmt.Errorf("-user-ip 参数是必须的")
				}
				userPortFlag := fc.GetFlagIfExistInOperatorByOperator("user-port", true, op)
				if !portFlag.Passed {
					return fmt.Errorf("-user-port 参数是必须的")
				}
				certFlag := fc.GetFlagIfExistInOperatorByOperator("cert", true, op)
				keyFlag := fc.GetFlagIfExistInOperatorByOperator("key", true, op)
				if certFlag.Passed || keyFlag.Passed {
					if !certFlag.Passed || !keyFlag.Passed {
						return fmt.Errorf("如何想开启 https 服务的话, -cert 和 -key 参数都是必须的")
					}
				}
				_ = ipFlag.Params[0] // 当前一定是 127.0.0.1
				port := portFlag.Params[0]
				userIp := userIpFlag.Params[0]
				userPort := userPortFlag.Params[0]
				cert := ""
				key := ""
				if certFlag.Passed && keyFlag.Passed {
					cert = certFlag.Params[0]
					key = keyFlag.Params[0]
				}
				if cert != "" && key != "" {
					err := startHttpsService(userIp, userPort, port, cert, key)
					if err != nil {
						return err
					}
				} else {
					fmt.Println("走 http 的服务逻辑, easy-proxy 中暂时不需要用到 http 服务")
					return nil
				}
			case "kill":
				pidFlag := fc.GetFlagIfExistInOperatorByOperator("pid", true, op)
				if !pidFlag.Passed {
					return fmt.Errorf("kill 指令必须传递 -pid 参数")
				}
				pid := ""
				if len(pidFlag.Params) > 0 {
					pid = pidFlag.Params[0]
				}
				if pid == "" {
					return fmt.Errorf("pid 参数不能为空")
				}
				fmt.Println("这里要干掉的 pid 号是: ", pid)
			default:
				return fmt.Errorf("不支持的操作")
			}

			return nil
		})

	testCmd := "easy-proxy-service " + strings.Join(os.Args[1:], " ")
	err := cmd.ExecuteStr(testCmd)
	if err != nil {
		fmt.Println(err.Error())
	}
}
