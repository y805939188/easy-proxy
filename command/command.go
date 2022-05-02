package command

import (
	"fmt"

	myNet "easy-proxy/net"
)

func ExecSetCommand(source, target *myNet.IUrl) error {
	fmt.Println("开始设置代理规则......")
	p, err := GetSetProxy()
	if err != nil {
		return err
	}
	if source.IsDomain {
		err = p.ProxyDomainToIp(source.IsHttps, source.Address, source.Port, target.Address, target.Port)
	} else {
		err = p.ProxyIpToIp(source.IsHttps, source.Address, source.Port, target.Address, target.Port)
	}
	if err != nil {
		_err := p.Fresh()
		if _err != nil {
			fmt.Println(err)
		}
		return err
	}
	fmt.Println("完成!")
	return nil
}

func ExecDelCommand(ids ...string) error {
	fmt.Println("正在删除规则......")
	p, err := GetDelProxy()
	if err != nil {
		return err
	}
	err = p.DeleteProxys(ids...)
	if err != nil {
		return err
	}
	fmt.Println("完成!")
	return nil
}

func ExecFreshCommand() error {
	fmt.Println("正在清空规则......")
	p, err := GetFreshProxy()
	if err != nil {
		return err
	}
	err = p.FreshAll()
	if err != nil {
		return err
	}
	fmt.Println("完成!")
	return nil
}

func ExecListCommand() error {
	p, err := GetListProxy()
	if err != nil {
		return err
	}
	err = p.List()
	if err != nil {
		return err
	}
	return nil
}
