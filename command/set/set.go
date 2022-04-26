package command

import (
	"easy-proxy/certificate"
	myIpt "easy-proxy/iptables"
	myNet "easy-proxy/net"
	"easy-proxy/tools"
	"fmt"
)

type Proxy struct {
	iptClient     *myIpt.Iptables
	delFunc       []*func() error
	domain        string
	subServicePid []int
}

func VerifyAddrValid(addrs ...string) error {
	// TODO: 验证可以更严格一些
	for _, addr := range addrs {
		if addr == "" {
			return fmt.Errorf(fmt.Sprintf("addr %s 不是一个有效的地址", addr))
		}
	}
	return nil
}

func createCaAndAddToSystemp(domain, ip string) (string, string, error) {
	tmpPath, err := tools.GetTmpCaPath()
	if err != nil {
		return "", "", err
	}
	var keyPath string
	var certPath string
	if domain != "" {
		keyPath, certPath, err = certificate.GenCertificate(domain, tmpPath)
	} else {
		keyPath, certPath, err = certificate.GenCertificate(ip, tmpPath)
	}
	if err != nil {
		return "", "", err
	}
	err = certificate.SetCertificateToSystemByCertPath(certPath)
	if err != nil {
		return "", "", err
	}

	return keyPath, certPath, nil
}

func _createService(userTargetIp, userTargetPort string, port string, certPath, keyPath string) (int, error) {
	svcBinaryPath, err := tools.GetTmpLocalServicePath()
	if err != nil {
		return -1, err
	}
	pid, err := tools.ExecCommand(
		svcBinaryPath,
		"start",
		"-ip",
		"127.0.0.1",
		"-port",
		port,
		"-user-ip",
		userTargetIp,
		"-user-port",
		userTargetPort,
		"-cert",
		certPath,
		"-key",
		keyPath,
	)
	if err != nil {
		return pid, err
	}
	return pid, nil
}

func createLocalHttpsMiddleService(localPort, userTargetIp, tuserTargetPort, keyPath, certPath string) (int, error) {
	err := VerifyAddrValid(keyPath, certPath)
	if err != nil {
		return -1, err
	}

	pid, err := _createService(userTargetIp, tuserTargetPort, localPort, certPath, keyPath)
	return pid, nil
}

func (p *Proxy) ProxyHttpsIpToIp(srcIp, srcPort, targetIp, targetPort string) error {
	err := VerifyAddrValid(srcIp, srcPort, targetIp, targetPort)
	if err != nil {
		return err
	}

	// 创建 ca 证书并将其干到系统受信任里
	keyPath, certPath, err := createCaAndAddToSystemp(p.domain, srcIp)
	if err != nil {
		return err
	}

	// 获取一个有效的端口号
	port, err := myNet.GetAvailablePort()
	if err != nil {
		return err
	}

	// 设置 iptables 规则, 把要被代理的 ip 以及 port 给代理到 127.0.0.1:port
	_, err = p.iptClient.SetIpToIp(myIpt.IP{IP: srcIp, Port: srcPort}, myIpt.IP{IP: "127.0.0.1", Port: port})
	if err != nil {
		return err
	}
	// 创建一个本地的 https://127.0.0.1:port 的服务用来代理
	// targetIp 和 targetPort 表示用户真正想要打到的自己的服务那里
	// 用户不感知中间这层 https://127.0.0.1:port
	pid, err := createLocalHttpsMiddleService(port, targetIp, targetPort, keyPath, certPath)
	if err != nil {
		return err
	}
	if pid > 0 {
		p.subServicePid = append(p.subServicePid, pid)
	}

	return nil
}

func (p *Proxy) ProxyHttpIpToIp(srcIp, srcPort, targetIp, targetPort string) error {
	err := VerifyAddrValid(srcIp, srcPort, targetIp, targetPort)
	if err != nil {
		return err
	}

	del, err := p.iptClient.SetIpToIp(myIpt.IP{IP: srcIp, Port: srcPort}, myIpt.IP{IP: targetIp, Port: targetPort})
	if err != nil {
		return err
	}
	if p.delFunc == nil {
		p.delFunc = [](*func() error){
			&del,
		}
	} else {
		p.delFunc = append(p.delFunc, &del)
	}
	return nil
}

func (p *Proxy) ProxyIpToIp(isHttps bool, srcIp, srcPort, targetIp, targetPort string) error {
	err := VerifyAddrValid(srcIp, targetIp)
	if err != nil {
		return err
	}

	if isHttps {
		if srcPort == "" {
			srcPort = "443"
		}
		if targetPort == "" {
			targetPort = srcPort
		}
		// 如果是 https 的代理的话需要中间加一层
		err = p.ProxyHttpsIpToIp(srcIp, srcPort, targetIp, targetPort)
		if err != nil {
			return err
		}
	} else {
		if srcPort == "" {
			srcPort = "80"
		}
		if targetPort == "" {
			targetPort = srcPort
		}
		// 如果是 http 的代理直接走 iptables 规则就 ok
		err = p.ProxyHttpIpToIp(srcIp, srcPort, targetIp, targetPort)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Proxy) ProxyDomainToIp(isHttps bool, srcDomain, srcPort, targetIp, targetPort string) error {
	err := VerifyAddrValid(srcDomain, targetIp)
	if err != nil {
		return err
	}

	// 先根据域名去解析对应的 ip list
	ipList, err := myNet.DnsAnalyzer(srcDomain)
	if err != nil {
		return err
	}

	p.domain = srcDomain

	for _, ip := range ipList {
		// 对域名对应的 ip 做代理
		err = p.ProxyIpToIp(isHttps, ip, srcPort, targetIp, targetPort)
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateProxy() (*Proxy, error) {
	iptClient, err := myIpt.GetNewIptablesClient()
	if err != nil {
		return nil, err
	}
	p := &Proxy{
		iptClient:     iptClient,
		subServicePid: []int{},
	}
	return p, nil
}
