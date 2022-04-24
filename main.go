package main

import (
	"easy-proxy/certificate"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
)

// https://serverfault.com/questions/18364/can-i-use-the-same-wildcard-certification-for-domain-com-and-domain-com

func dnsAnalyzer(domain string) ([]string, error) {
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

func main() {

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
