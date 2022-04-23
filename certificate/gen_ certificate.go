package certificate

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"regexp"
	"strings"
	"time"
	// openssl "github.com/spacemonkeygo/openssl"
)

func SetCertificateToSystem(certPath string) error {
	rootCAs, err := x509.SystemCertPool()
	if err != nil {
		return err
	}

	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	certs, err := ioutil.ReadFile(certPath)
	if err != nil {
		return err
	}

	if ok := rootCAs.AppendCertsFromPEM(certs); !ok {
		return fmt.Errorf("add cert to system error")
	}

	return nil
}

func GenCertificate(domain string, path string, index string) (string, string, error) {

	if domain == "" {
		return "", "", fmt.Errorf("domain is required")
	}

	domain = strings.Replace(domain, "http://", "", 1)
	domain = strings.Replace(domain, "https://", "", 1)

	reg, err := regexp.Compile("\\.")
	if err != nil {
		return "", "", err
	}

	key := reg.ReplaceAllString(domain, "_")

	currentPath := "./"
	if path != "" {
		if string(path[len(path)-1]) == "/" {
			currentPath = path[0 : len(path)-1]
		} else {
			currentPath = path
		}
	}

	if key == "" {
		return "", "", fmt.Errorf("key is error")
	}

	max := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", "", err
	}

	// 定义：引用IETF的安全领域的公钥基础实施（PKIX）工作组的标准实例化内容
	subject := pkix.Name{
		Organization: []string{domain},
		// OrganizationalUnit: []string{"client"},
		// CommonName:         "",
	}

	// 设置 SSL证书的属性用途
	certificate509 := x509.Certificate{
		SerialNumber: serialNumber,
		Subject:      subject,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(100 * 24 * time.Hour),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		// IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
	}

	// 生成指定位数密匙
	pk, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return "", "", err
	}

	// 生成 SSL公匙
	derBytes, err := x509.CreateCertificate(rand.Reader, &certificate509, &certificate509, &pk.PublicKey, pk)
	if err != nil {
		return "", "", err
	}

	certPath := currentPath + "/" + key + "cert" + index + ".pem"
	certOut, err := os.Create(certPath)
	if err != nil {
		return "", "", err
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()

	keyPath := currentPath + "/" + key + "key" + index + ".pem"
	// 生成 SSL私匙
	keyOut, err := os.Create(keyPath)
	if err != nil {
		return "", "", err
	}
	pem.Encode(keyOut, &pem.Block{Type: "RAS PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)})
	keyOut.Close()
	return keyPath, certPath, nil
}
