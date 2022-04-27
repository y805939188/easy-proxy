package command

import (
	"easy-proxy/certificate"
	myIpt "easy-proxy/iptables"
	"easy-proxy/tools"
	"encoding/json"
	"fmt"
)

type DelProxy struct {
	iptClient *myIpt.Iptables
}

func (p *DelProxy) removeIptablesRuls(rules []*Iptables) error {
	if rules == nil {
		return nil
	}

	var err error
	for _, rule := range rules {
		// 尽量删, 某条规则删除错了也尽量尝试删下一条
		err = p.iptClient.DeleteIptablesOutputRule(rule.SourceIp, rule.SourcePort, rule.TargetIp, rule.TargetPort)
		if err != nil {
			fmt.Println("删除规则失败, err: ", err.Error())
		}
	}

	return err
}

func (p *DelProxy) killPIDs(pids []int) error {
	if pids == nil {
		return nil
	}

	var err error
	for _, pid := range pids {
		err = tools.KillByPID(pid)
		if err != nil {
			fmt.Println(fmt.Sprintf("终止进程 %d 失败, err: %s", pid, err.Error()))
		}
	}
	return err
}

func (p *DelProxy) removeCerts(certs []*Cert) error {
	if certs == nil {
		return nil
	}

	var err error
	for _, cert := range certs {
		certPath := cert.CertPath
		keyPath := cert.KeyPath
		err = VerifyAddrValid(certPath, keyPath)
		if err != nil {
			return err
		}

		fileName, _ := tools.GetFileNameAndExt(certPath)
		err = certificate.RemoveCertificateFromSystemByCertName(fileName)
		if err != nil {
			fmt.Println(fmt.Sprintf("删除 %s 失败, err: %s", fileName+".crt", err.Error()))
		}

		for _, file := range []string{certPath, keyPath} {
			err = tools.DeleteFile(file)
			if err != nil {
				fmt.Println(fmt.Sprintf("删除证书文件 %s 失败, err: %s", file, err.Error()))
			}
		}
	}

	return err
}

func (p *DelProxy) removeProxyInfo(fileName string) error {
	if fileName == "" {
		return nil
	}

	return tools.DeleteFile(fileName)
}

func (p *DelProxy) DeleteOneProxy(id string) error {
	if id == "" {
		return nil
	}
	filePath, err := tools.GetTmpProxyInfoFilePath(id)
	if err != nil {
		return err
	}
	if !tools.FileIsExisted(filePath) {
		return fmt.Errorf("未找到有效的 proxy 信息文件")
	}
	content, err := tools.ReadFile(filePath)
	if err != nil {
		return err
	}
	var proxyInfo Rule
	err = json.Unmarshal(content, &proxyInfo)
	if err != nil {
		return err
	}

	err = p.removeIptablesRuls(proxyInfo.Iptables)
	if err != nil {
		return err
	}

	err = p.killPIDs(proxyInfo.PIDs)
	if err != nil {
		return err
	}

	err = p.removeCerts(proxyInfo.Certs)
	if err != nil {
		return err
	}

	return p.removeProxyInfo(filePath)
}

func (p *DelProxy) DeleteProxys(ids ...string) error {
	if ids == nil {
		return nil
	}
	for _, id := range ids {
		err := p.DeleteOneProxy(id)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetDelProxy() (*DelProxy, error) {
	iptClient, err := myIpt.GetNewIptablesClient()
	if err != nil {
		return nil, err
	}
	p := &DelProxy{
		iptClient: iptClient,
	}
	return p, nil
}
