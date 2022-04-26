package main

import (
	mySvc "easy-proxy/binary_service"
	"easy-proxy/certificate"
	command "easy-proxy/command"
	"os"

	myNet "easy-proxy/net"
	"easy-proxy/tools"
	"fmt"
	"strings"

	"github.com/y805939188/dcommand"
)

func execProxy(operator string, source, target *myNet.IUrl, ids ...string) error {
	if operator == "set" {
		fmt.Println("开始设置代理规则......")
		p, err := command.GetSetProxy()
		if err != nil {
			return err
		}
		if source.IsDomain {
			err = p.ProxyDomainToIp(source.IsHttps, source.Address, source.Port, target.Address, target.Port)
		} else {
			err = p.ProxyIpToIp(source.IsHttps, source.Address, source.Port, target.Address, target.Port)
		}
		if err != nil {
			return err
		}
	} else if operator == "del" {
		fmt.Println("正在删除规则......")
		p, err := command.GetDelProxy()
		if err != nil {
			return nil
		}
		err = p.DeleteProxys(ids...)
		if err != nil {
			return err
		}
	} else if operator == "fresh" {
		fmt.Println("正在清空规则......")
	} else if operator == "list" {
		p, err := command.GetListProxy()
		if err != nil {
			return err
		}
		err = p.List()
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf(fmt.Sprintf("暂不支持的操作: %s", operator))
	}
	fmt.Println("完成!")
	return nil
}

func initFiles() error {
	easyProxyRootPath, err := tools.GetEasyRootPath()
	if err != nil {
		return fmt.Errorf("获取程序根路径失败")
	}
	if !tools.FileIsExisted(easyProxyRootPath) {
		err := tools.CreateDir(easyProxyRootPath)
		if err != nil {
			return fmt.Errorf("创建程序根路径失败")
		}
		tmpCaPath, err := tools.GetTmpCaPath()
		if err != nil {
			return fmt.Errorf("获取程序临时证书存放目录失败")
		}
		err = tools.CreateDir(tmpCaPath)
		if err != nil {
			return fmt.Errorf("创建程序证书临时存储目录失败")
		}
	}
	scriptPath, err := tools.GetUpdateSystemCaScriptPath()
	if err != nil {
		return err
	}
	if !tools.FileIsExisted(scriptPath) {
		err = tools.CreateFile(
			scriptPath,
			[]byte(certificate.UpdateCaCertificatesBashScriptContent),
			0766,
		)
		if err != nil {
			return err
		}
	}

	servicePath, err := tools.GetTmpLocalServicePath()
	if err != nil {
		return err
	}
	if !tools.FileIsExisted(servicePath) {
		// 拷贝启动本地代理服务器的二进制到临时目录下
		// svcRootPath 是这个二进制文件的父级目录
		svcRootPath, err := tools.GetTmpLocalServicePath()
		if err != nil {
			return err
		}
		err = mySvc.RestoreAssets(svcRootPath, "service")
		if err != nil {
			return err
		}
	}

	proxyInfoPath, err := tools.GetTmpProxyInfoPath()
	if err != nil {
		return err
	}
	if !tools.FileIsExisted(proxyInfoPath) {
		err := tools.CreateDir(proxyInfoPath)
		if err != nil {
			return fmt.Errorf("创建 proxy info 目录失败")
		}
	}
	return nil
}

func main() {
	err := initFiles()
	if err != nil {
		fmt.Println("初始化项目路径失败, err: ", err.Error())
		return
	}

	cmd := &dcommand.DCommand{}
	cmd.Command("easy-proxy").
		Operator("set").
		Flag(&dcommand.FlagInfo{Name: "source", Short: "s"}).
		Flag(&dcommand.FlagInfo{Name: "target", Short: "t"}).
		Operator("del").
		Flag(&dcommand.FlagInfo{Name: "id", Short: "id"}).
		Operator("fresh").
		Operator("list").
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
			switch currentOperator {
			case "set":
				sourceFlag := fc.GetFlagIfExistInOperatorByOperator("source", true, op)
				targetFlag := fc.GetFlagIfExistInOperatorByOperator("target", true, op)
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
				if strings.HasPrefix(targetFlag.Params[0], "https://") {
					return fmt.Errorf("暂不支持把请求代理到 https 服务")
				}
				if !strings.HasPrefix(targetFlag.Params[0], "http://") {
					targetFlag.Params[0] = "http://" + targetFlag.Params[0]
				}
				targetHost, err := myNet.IsValidAddr(targetFlag.Params[0])
				if err != nil {
					return err
				}
				if targetHost.IsDomain {
					return fmt.Errorf("target 暂不支持域名")
				}
				err = execProxy(currentOperator, sourceHost, targetHost)
				if err != nil {
					return err
				}
			case "del":
				idFlag := fc.GetFlagIfExistInOperatorByOperator("id", true, op)
				if !idFlag.Passed {
					return fmt.Errorf("del 指令一定需要一个 -id 作为参数")
				}
				if len(idFlag.Params) == 0 {
					return fmt.Errorf("-id 缺少参数")
				}
				err = execProxy(currentOperator, nil, nil, idFlag.Params...)
				if err != nil {
					return err
				}
			case "fresh":
				fallthrough
			case "list":
				err := execProxy(currentOperator, nil, nil)
				if err != nil {
					return err
				}
				break
			default:
				return fmt.Errorf("无效的操作符")
			}
			return nil
		})

	testCmd := "easy-proxy " + strings.Join(os.Args[1:], " ")
	// testCmd := "easy-proxy set -s https://www.baidu.com -t 127.0.0.1:13191"
	// testCmd := "easy-proxy list"
	fmt.Println("这里的 cmd 是: ", testCmd)
	err = cmd.ExecuteStr(testCmd)
	if err != nil {
		fmt.Println(err.Error())
	}
}
