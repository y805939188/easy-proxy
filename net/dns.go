package net

import "net"

func DnsAnalyzer(domain string) ([]string, error) {
	var res []string
	ns, err := net.LookupHost(domain)
	if err != nil {
		return nil, err
	}

	if len(ns) != 0 {
		for _, n := range ns {
			res = append(res, n)
		}
	}
	return res, nil
}
