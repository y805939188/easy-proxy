package main

import (
	"easy-proxy/certificate"
	myNet "easy-proxy/net"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/y805939188/dcommand"
)

type IUrl struct {
	Scheam   string
	IsHttps  bool
	Address  string
	Port     string
	IsDomain bool
}

func main() {
	cmd := &dcommand.DCommand{}
	cmd.Command("easy-proxy").
		Operator("set").
		Flag(&dcommand.FlagInfo{Name: "source", Short: "s"}).
		Flag(&dcommand.FlagInfo{Name: "target", Short: "t"}).
		Operator("del").
		Flag(&dcommand.FlagInfo{Name: "source", Short: "s"}).
		Flag(&dcommand.FlagInfo{Name: "target", Short: "t"}).
		Operator("fresh").
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
				return fmt.Errorf("easy-proxy 命令需要至少一个 operator: set | del | fresh")
			}
			cmd := fc.GetCommandIfExist(command)
			op := fc.GetOperatorIfExistByCommand(currentOperator, cmd)
			sourceFlag := fc.GetFlagIfExistInOperatorByOperator("source", true, op)
			targetFlag := fc.GetFlagIfExistInOperatorByOperator("target", true, op)
			switch currentOperator {
			case "set":
				fallthrough
			case "del":
				if !sourceFlag.Passed || !targetFlag.Passed {
					return fmt.Errorf(fmt.Sprintf("%s 指令一定需要 %s & %s 参数", currentOperator, "-s", "-t"))
				}
				if len(sourceFlag.Params) != 1 || len(targetFlag.Params) != 1 {
					return fmt.Errorf(fmt.Sprintf("-s & -t 参数需要 1 个参数"))
				}
				sourceHost, err := myNet.IsValidAddr(sourceFlag.Params[0])
				if err != nil {
					return err
				}
				targetHost, err := myNet.IsValidAddr(targetFlag.Params[0])
				if err != nil {
					return err
				}
				fmt.Println("这里的 sourceHost 是: ", sourceHost)
				fmt.Println("这里的 targetHost 是: ", targetHost)
				break
			case "fresh":
				fmt.Println("走到了 fresh")
				break
			default:
				return fmt.Errorf("无效的操作符")
			}
			return nil
		})

	testCmd := "easy-proxy " + strings.Join(os.Args[1:], " ")
	fmt.Println("这里的 cmd 是: ", testCmd)
	err := cmd.ExecuteStr(testCmd)
	if err != nil {
		fmt.Println(err.Error())
	}
	return

	keyPath, certPath, err := certificate.GenCertificate("www.baidu.com", "./")
	if err != nil {
		fmt.Println("这里的 err 是: ", err.Error())
		return
	}

	err = certificate.SetCertificateToSystemByCertPath(certPath)
	if err != nil {
		fmt.Println("这里的 err 是: ", err.Error())
	}

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("ding test")
		io.WriteString(w, "Hello, TLS!\n")
	})

	// One can use generate_cert.go in crypto/tls to generate cert.pem and key.pem.
	log.Printf("About to listen on 13191. Go to https://127.0.0.1:13191/")
	err = http.ListenAndServeTLS(":13191", certPath, keyPath, nil)
	log.Fatal(err)

}

// openssl req -x509 -newkey rsa:2048 -sha256 -nodes -keyout key1.pem -out cert1.pem -subj "/CN=www.baidu.com" -days 1
