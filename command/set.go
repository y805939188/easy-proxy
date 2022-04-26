package command

import (
	"easy-proxy/certificate"
	myIpt "easy-proxy/iptables"
	myNet "easy-proxy/net"
	"easy-proxy/tools"
	"encoding/json"
	"fmt"
)

type Cert struct {
	CertPath string
	KeyPath  string
}

type Iptables struct {
	SourceIp   string
	SourcePort string
	TargetIp   string
	TargetPort string
}

type Rule struct {
	Certs      []*Cert
	Ports      []string
	Iptables   []*Iptables
	PIDs       []int
	Identifier *Identifier
}

type Identifier struct {
	SourceAddress string
	SourcePort    string
	TargetAddress string
	TargetPort    string
}

type SetProxy struct {
	iptClient         *myIpt.Iptables
	delFunc           []*func() error
	currentIdentifier string
	rules             map[string]*Rule
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

func createCaAndAddToSystemp(commonName string) (string, string, error) {
	tmpPath, err := tools.GetTmpCaPath()
	if err != nil {
		return "", "", err
	}
	var keyPath string
	var certPath string
	// if domain != "" {
	// 	keyPath, certPath, err = certificate.GenCertificate(domain, tmpPath)
	// } else {
	// 	keyPath, certPath, err = certificate.GenCertificate(ip, tmpPath)
	// }
	keyPath, certPath, err = certificate.GenCertificate(commonName, tmpPath)
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

func (p *SetProxy) setCaInfo(cert, key string) {
	if rule, ok := p.rules[p.currentIdentifier]; ok {
		// 记录下这俩证书相关的位置以便日后清除
		if rule.Certs == nil {
			rule.Certs = []*Cert{}
		}
		rule.Certs = append(rule.Certs, &Cert{
			CertPath: cert,
			KeyPath:  key,
		})
	}
}

func (p *SetProxy) setPort(port string) {
	if rule, ok := p.rules[p.currentIdentifier]; ok {
		// 记录下 port 以便日后 kill
		if rule.Ports == nil {
			rule.Ports = []string{}
		}
		rule.Ports = append(rule.Ports, port)
	}
}

func (p *SetProxy) setPID(pid int) {
	if rule, ok := p.rules[p.currentIdentifier]; ok {
		// 记录下 port 以便日后 kill
		if rule.PIDs == nil {
			rule.PIDs = []int{}
		}
		rule.PIDs = append(rule.PIDs, pid)
	}
}

func (p *SetProxy) setIptablesRule(srcIp, srcPort, targetIp, targetPort string) {
	if rule, ok := p.rules[p.currentIdentifier]; ok {
		// 记录下 iptables 规则以便日后 -D
		if rule.Iptables == nil {
			rule.Iptables = []*Iptables{}
		}
		rule.Iptables = append(rule.Iptables, &Iptables{
			SourceIp:   srcIp,
			SourcePort: srcPort,
			TargetIp:   targetIp,
			TargetPort: targetPort,
		})
	}
}

func (p *SetProxy) proxyHttpsIpToIp(srcIps []string, srcPort, targetIp, targetPort string, domain string) error {
	if srcIps == nil {
		return nil
	}

	// 如果没有 domain 的话说明就是单纯地代理 ip
	// 要给每个 ip 都做个证书以及本地的代理服务
	// 用户不感知这一层
	if domain == "" {
		for _, srcIp := range srcIps {
			// 创建 ca 证书并将其干到系统受信任里
			keyPath, certPath, err := createCaAndAddToSystemp(srcIp)
			if err != nil {
				return err
			}
			p.setCaInfo(certPath, keyPath)

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
			p.setIptablesRule(srcIp, srcPort, "127.0.0.1", port)

			// 创建一个本地的 https://127.0.0.1:port 的服务用来代理
			pid, err := createLocalHttpsMiddleService(port, targetIp, targetPort, keyPath, certPath)
			if err != nil {
				return err
			}
			p.setPort(port)
			if pid > 0 {
				p.setPID(pid)
			}
		}
	} else {
		// 如果有 domain 的话
		// 此时可能一个域名会对应多个 ip 地址
		// 但是多个 ip 地址可以映射到同一个本地代理服务
		keyPath, certPath, err := createCaAndAddToSystemp(domain)
		if err != nil {
			return err
		}
		p.setCaInfo(certPath, keyPath)

		port, err := myNet.GetAvailablePort()
		if err != nil {
			return err
		}

		// 把每个 ip 都映射到同一个本地代理服务
		for _, srcIp := range srcIps {
			_, err = p.iptClient.SetIpToIp(myIpt.IP{IP: srcIp, Port: srcPort}, myIpt.IP{IP: "127.0.0.1", Port: port})
			if err != nil {
				return err
			}
			p.setIptablesRule(srcIp, srcPort, "127.0.0.1", port)
		}

		pid, err := createLocalHttpsMiddleService(port, targetIp, targetPort, keyPath, certPath)
		if err != nil {
			return err
		}
		p.setPort(port)
		if pid > 0 {
			p.setPID(pid)
		}
	}

	return nil
}

func (p *SetProxy) proxyHttpIpToIp(srcIps []string, srcPort, targetIp, targetPort string) error {
	if srcIps == nil {
		return nil
	}

	for _, srcIp := range srcIps {
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
	}
	return nil
}

func (p *SetProxy) proxyIpsToIp(isHttps bool, srcIps []string, srcPort, targetIp, targetPort, domain string) error {
	if srcIps == nil {
		return nil
	}

	if isHttps {
		if srcPort == "" {
			srcPort = "443"
		}
		if targetPort == "" {
			targetPort = srcPort
		}
		// 如果是 https 的代理的话需要中间加一层
		err := p.proxyHttpsIpToIp(srcIps, srcPort, targetIp, targetPort, domain)
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
		err := p.proxyHttpIpToIp(srcIps, srcPort, targetIp, targetPort)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *SetProxy) writeProxyInfoToDisk(id string) error {
	if rule, ok := p.rules[id]; ok {
		res, err := json.Marshal(rule)
		if err != nil {
			return err
		}
		filePath, err := tools.GetTmpProxyInfoFilePath(id)
		if err != nil {
			return err
		}
		return tools.CreateFile(filePath, []byte(res), 0766)
	}
	return nil
}

func (p *SetProxy) checkProxyInfoExist(id string) bool {
	proxyInfoName, err := tools.GetTmpProxyInfoFilePath(id)
	if err != nil {
		return false
	}
	if tools.FileIsExisted(proxyInfoName) {
		return true
	}
	return false
}

func (p *SetProxy) TransformIdentifierToStr(id *Identifier) string {
	return fmt.Sprintf(
		"%s-%s-%s-%s",
		id.SourceAddress,
		id.SourcePort,
		id.TargetAddress,
		id.TargetPort,
	)
}

func (p *SetProxy) ProxyIpToIp(isHttps bool, srcIp, srcPort, targetIp, targetPort string) error {
	err := VerifyAddrValid(srcIp, srcPort, targetIp, targetPort)
	if err != nil {
		return err
	}
	id := &Identifier{
		SourceAddress: srcIp,
		SourcePort:    srcPort,
		TargetAddress: targetIp,
		TargetPort:    targetPort,
	}

	idStr := tools.Sha1(p.TransformIdentifierToStr(id))
	if p.checkProxyInfoExist(idStr) {
		return fmt.Errorf("该规则已经存在")
	}

	p.rules[idStr] = &Rule{
		Identifier: id,
	}
	originId := p.currentIdentifier
	p.currentIdentifier = idStr
	err = p.proxyIpsToIp(isHttps, []string{srcIp}, srcPort, targetIp, targetPort, "")
	p.currentIdentifier = originId
	if err != nil {
		return err
	}
	p.writeProxyInfoToDisk(idStr)
	return nil
}

func (p *SetProxy) ProxyDomainToIp(isHttps bool, srcDomain, srcPort, targetIp, targetPort string) error {
	err := VerifyAddrValid(srcDomain, targetIp)
	if err != nil {
		return err
	}

	// 先根据域名去解析对应的 ip list
	ipList, err := myNet.DnsAnalyzer(srcDomain)
	if err != nil {
		return err
	}

	id := &Identifier{
		SourceAddress: srcDomain,
		SourcePort:    srcPort,
		TargetAddress: targetIp,
		TargetPort:    targetPort,
	}
	idStr := tools.Sha1(p.TransformIdentifierToStr(id))
	if p.checkProxyInfoExist(idStr) {
		return fmt.Errorf("该规则已经存在")
	}

	p.rules[idStr] = &Rule{
		Identifier: id,
	}
	originId := p.currentIdentifier
	p.currentIdentifier = idStr
	err = p.proxyIpsToIp(isHttps, ipList, srcPort, targetIp, targetPort, srcDomain)
	p.currentIdentifier = originId
	if err != nil {
		return err
	}
	p.writeProxyInfoToDisk(idStr)
	return nil
}

func GetSetProxy() (*SetProxy, error) {
	iptClient, err := myIpt.GetNewIptablesClient()
	if err != nil {
		return nil, err
	}
	p := &SetProxy{
		iptClient: iptClient,
		rules:     make(map[string]*Rule),
	}
	return p, nil
}
