package main

import (
	// "easy-proxy/certificate"

	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"math/big"
	"net"
	"net/http"
	"os"
	"time"

	openssl "github.com/spacemonkeygo/openssl"
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

func tempFn() {
	key, err := openssl.GenerateRSAKey(768)
	if err != nil {
		fmt.Println(1, err)
	}
	info := &openssl.CertificateInfo{
		Serial:  big.NewInt(int64(1)),
		Issued:  0,
		Expires: 24 * time.Hour,
		// Country:      "US",
		Organization: "www.baidu.com",
		CommonName:   "localhost",
	}
	cert, err := openssl.NewCertificate(info, key)
	if err != nil {
		fmt.Println(2, err)
	}
	fmt.Println(cert)
}

// func main() {
// 	// fmt.Println("进来了")
// 	// // fmt.Println(dnsAnalyzer("www.baidu.com"))
// 	// http.HandleFunc("/ding1", func(w http.ResponseWriter, r *http.Request) {
// 	// 	fmt.Println("来了!!!")
// 	// 	w.Write([]byte("hello, world !"))
// 	// })

// 	// log.Fatal(http.ListenAndServeTLS(":13191", "./cert6.pem", "./key6.pem", nil))
// 	tempFn()
// 	return
// 	index := "3"
// 	// keyPath, certPath, err := certificate.GenCertificate("www.baidu.com", "./", index)
// 	// if err != nil {
// 	// 	fmt.Println("这里的 error 是: ", err)
// 	// }
// 	// fmt.Println(keyPath)
// 	// fmt.Println(certPath)
// 	// // err = certificate.SetCertificateToSystem(certPath)
// 	// // if err != nil {
// 	// // 	fmt.Println("这里的 222 error 是: ", err)
// 	// // }
// 	http.HandleFunc("/ding1", func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Println("来了!!!")
// 		w.Write([]byte("hello, world !"))
// 	})
// 	log.Fatal(http.ListenAndServeTLS(":13191", "./www_baidu_comcert"+index+".pem", "./www_baidu_comkey"+index+".pem", nil))
// 	// reg, _ := regexp.Compile("\\.")

// 	// key := reg.ReplaceAllString("www.baidu.com", "_")

// 	// fmt.Println("这里的 key 是: ", key)

// }

// func gen() {
// 	// max := new(big.Int).Lsh(big.NewInt(1), 128)
// 	// // serialNumber, _ := rand.Int(rand.Reader, max)
// 	// subject := pkix.Name{
// 	// 	// Organization:       []string{"baidu.com"},
// 	// 	// OrganizationalUnit: []string{"/CN"},
// 	// 	// CommonName:         "GO Web",
// 	// }

// 	rootTemplate := x509.Certificate{
// 		// SerialNumber: serialNumber,
// 		// SerialNumber: big.NewInt(1),
// 		// Subject:      subject,
// 		Subject: pkix.Name{
// 			// Organization: []string{"www.baidu.com"},
// 			CommonName: "www.baidu.com",
// 		},
// 		NotBefore:   time.Now(),
// 		NotAfter:    time.Now().Add(365 * time.Hour),
// 		KeyUsage:    x509.KeyUsageKeyEncipherment,
// 		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
// 		// IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
// 	}
// 	pk, _ := rsa.GenerateKey(rand.Reader, 2048)
// 	derBytes, _ := x509.CreateCertificate(rand.Reader, &rootTemplate, &rootTemplate, &pk.PublicKey, pk)

// 	certOut, _ := os.Create("cert9.pem")
// 	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

// 	certOut.Close()

// 	keyOut, _ := os.Create("key9.pem")
// 	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)})
// 	keyOut.Close()

// }

func gen() {
	max := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, _ := rand.Int(rand.Reader, max)
	subject := pkix.Name{
		// Organization:       []string{"My Company"},
		// OrganizationalUnit: []string{"Person"},
		CommonName: "www.baidu.com",
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject:      subject,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(365 * 24 * time.Hour),
		// IsCA:         true,
		// SubjectKeyId:   []byte("ding1"),
		// AuthorityKeyId: []byte("ding1"),
		SubjectKeyId:   []byte{201, 32, 246, 127, 177, 93, 206, 156, 86, 226, 91, 218, 23, 165, 148, 22, 57, 155, 65, 147},
		AuthorityKeyId: []byte{201, 32, 246, 127, 177, 93, 206, 156, 86, 226, 91, 218, 23, 165, 148, 22, 57, 155, 65, 147},
		// KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		// ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		// IPAddresses:                 []net.IP{net.ParseIP("127.0.0.1")},
		// PermittedDNSDomainsCritical: false,
	}

	pk, _ := rsa.GenerateKey(rand.Reader, 2048)

	derBytes, _ := x509.CreateCertificate(rand.Reader, &template, &template, &pk.PublicKey, pk)

	certOut, _ := os.Create("cert-7.pem")
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()

	keyOut, _ := os.Create("key-7.pem")
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)})
	keyOut.Close()
}

func main() {
	// fmt.Println([]byte("ding1"))
	// return
	// gen()
	// return
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Hello, TLS!\n")
	})

	// One can use generate_cert.go in crypto/tls to generate cert.pem and key.pem.
	log.Printf("About to listen on 13191. Go to https://127.0.0.1:13191/")
	err := http.ListenAndServeTLS(":13191", "cert-7.pem", "key-7.pem", nil)
	log.Fatal(err)

}

// openssl req -x509 -newkey rsa:2048 -sha256 -nodes -keyout key1.pem -out cert1.pem -subj "/CN=www.baidu.com" -days 1
