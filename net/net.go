package net

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strings"
)

type IUrl struct {
	Scheam   string
	IsHttps  bool
	Address  string
	Port     string
	IsDomain bool
}

func IsValidIpv4(ip string) bool {
	address := net.ParseIP(ip)
	if address == nil {
		return false
	}
	return true
}

func IsValidAddr(addr string) (*IUrl, error) {
	if 0 == len(addr) {
		return nil, fmt.Errorf("url 不能为空")
	}
	parse, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	if "http" != parse.Scheme && "https" != parse.Scheme {
		return nil, fmt.Errorf("url 必须包含 http 或 https 域名")
	}
	re := regexp.MustCompile(`^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$|^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)+([A-Za-z]|[A-Za-z][A-Za-z0-9\-]*[A-Za-z0-9])$`)

	u := strings.Split(parse.Host, ":")
	var resUrl = &IUrl{
		Scheam:   parse.Scheme,
		IsDomain: false,
		Port:     "",
	}
	if len(u) == 1 || len(u) == 2 {
		resUrl.Address = u[0]
	}
	isIp := IsValidIpv4(u[0])
	if !isIp {
		resUrl.IsDomain = true
	}
	if parse.Scheme == "https" {
		resUrl.IsHttps = true
	}
	if len(u) == 2 && len(u[1]) != 0 {
		result := re.FindAllStringSubmatch(u[0], -1)
		if result == nil {
			return nil, errors.New("url 无效: " + u[0])
		}
		resUrl.Port = u[1]
		return resUrl, nil
	}

	result := re.FindAllStringSubmatch(parse.Host, -1)
	if result == nil {
		return nil, errors.New("url 无效: " + parse.Host)
	}

	return resUrl, nil
}
