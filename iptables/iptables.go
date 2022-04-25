package iptables

import (
	"fmt"
	"net"

	"github.com/coreos/go-iptables/iptables"
)

type Iptables struct {
	ipv4 *iptables.IPTables
}

type IP struct {
	ip   string
	port string
}

func (i *Iptables) getOutputRuleArrFromSrcAndDst(srcIp, srcPort, dstIp, dstPort string) []string {
	if srcIp == "" || dstIp == "" {
		return nil
	}
	currentDstPort := "80"
	if dstPort != "" {
		currentDstPort = dstPort
	}

	var ruleArr []string
	if srcPort != "" {
		ruleArr = []string{"-p", "tcp", "-d", srcIp, "--dport", srcPort, "-j", "DNAT", "--to-destination", dstIp + ":" + currentDstPort}
	} else {
		ruleArr = []string{"-p", "tcp", "-d", srcIp, "-j", "DNAT", "--to-destination", dstIp + ":" + currentDstPort}
	}
	return ruleArr
}

func (i *Iptables) DeleteIptablesOutputRule(srcIp, srcPort, dstIp, dstPort string) error {
	ruleArr := i.getOutputRuleArrFromSrcAndDst(srcIp, srcPort, dstIp, dstPort)
	if ruleArr == nil {
		return fmt.Errorf("get output rule arr error")
	}
	err := i.ipv4.Delete("nat", "OUTPUT", ruleArr...)
	if err != nil {
		return err
	}
	return nil
}

func (i *Iptables) SetIptablesOutputRule(srcIp, srcPort, dstIp, dstPort string) (func() error, error) {
	ruleArr := i.getOutputRuleArrFromSrcAndDst(srcIp, srcPort, dstIp, dstPort)
	if ruleArr == nil {
		return nil, fmt.Errorf("get output rule arr error")
	}

	err := i.ipv4.Append("nat", "OUTPUT", ruleArr...)
	if err != nil {
		return nil, err
	}

	delFunc := func() error {
		return i.DeleteIptablesOutputRule(srcIp, srcPort, dstIp, dstPort)
	}
	return delFunc, nil
}

/**
 * src 可以带 port 也可以不带
 * 不带 port 表示任何发往 src 的东西不管啥端口都往 dst 上怼
 * dst 必须带个 port, 如果没传的话默认给个 80
 */
func (i *Iptables) SetIpToIp(src, dst IP) (func() error, error) {
	srcIp := net.ParseIP(src.ip)
	if srcIp == nil {
		return nil, fmt.Errorf("src ip is not a valid ip")
	}

	dstIp := net.ParseIP(dst.ip)
	if dstIp == nil {
		return nil, fmt.Errorf("dst ip is not a valid ip")
	}

	delFunc, err := i.SetIptablesOutputRule(src.ip, src.port, dst.ip, dst.port)

	if err != nil {
		return nil, err
	}
	return delFunc, nil
}

func GetNewIptablesClient() (*Iptables, error) {
	iClient, err := iptables.NewWithProtocol(iptables.ProtocolIPv4)
	if err != nil {
		return nil, err
	}
	return &Iptables{
		ipv4: iClient,
	}, nil
}
