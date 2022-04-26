package command

import (
	myIpt "easy-proxy/iptables"
	"easy-proxy/tools"
	"encoding/json"
	"fmt"
)

type DelProxy struct {
	iptClient *myIpt.Iptables
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
	return nil
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
