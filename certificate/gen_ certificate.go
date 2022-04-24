package certificate

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"easy-proxy/consts"
	"easy-proxy/tools"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func SetCertificateToSystemByCertPath(certPath string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	fileName := filepath.Base(certPath)
	ext := path.Ext(fileName)
	currentFileName := strings.Replace(fileName, ext, ".crt", 1)
	err = tools.FileCopy(certPath, consts.UbuntuCaCertificatesPath+"/"+currentFileName)
	if err != nil {
		return err
	}

	dir := cwd + "/certificate" + "/update-ca-certificates"

	out, _ := tools.Bash("chmod u+x " + dir)
	fmt.Println(out)

	out, _ = tools.Bash(dir)
	fmt.Println(out)

	return nil
}

func GenCertificate(domain string, path string) (string, string, error) {
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
	serialNumber, _ := rand.Int(rand.Reader, max)
	subject := pkix.Name{
		CommonName: domain,
	}

	pk, _ := rsa.GenerateKey(rand.Reader, 2048)

	// 要用下边几行生成 SubjectKeyId 和 AuthorityKeyId
	// SubjectKeyId 和 AuthorityKeyId 一样
	spkiASN1, err := x509.MarshalPKIXPublicKey(&pk.PublicKey)
	if err != nil {
		return "", "", err
	}

	var spki struct {
		Algorithm        pkix.AlgorithmIdentifier
		SubjectPublicKey asn1.BitString
	}
	_, err = asn1.Unmarshal(spkiASN1, &spki)
	skid := sha1.Sum(spki.SubjectPublicKey.Bytes)

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject:      subject,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(365 * 24 * time.Hour),
		// 这俩必须有, 一个 "颁发机构秘钥标识符", 一个 "证书使用者秘钥标识符"
		// 很多客户端的握手也会检查这俩东西
		SubjectKeyId:   skid[:],
		AuthorityKeyId: skid[:],
		// 其他的参数不能瞎加
		// 加了有可能变成 OV SSL 证书
		// OV SSL 证书校验比较严格
		// DV SSL 证书相对来说松快一点
	}

	derBytes, _ := x509.CreateCertificate(rand.Reader, &template, &template, &pk.PublicKey, pk)

	certPath := currentPath + "/" + key + "cert" + ".pem"
	certOut, _ := os.Create(certPath)
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()

	keyPath := currentPath + "/" + key + "key" + ".pem"
	keyOut, _ := os.Create(keyPath)
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)})
	keyOut.Close()

	return keyPath, certPath, nil
}

// func GenCertificate(domain string, path string, index string) (string, string, error) {

// 	if domain == "" {
// 		return "", "", fmt.Errorf("domain is required")
// 	}

// 	domain = strings.Replace(domain, "http://", "", 1)
// 	domain = strings.Replace(domain, "https://", "", 1)

// 	reg, err := regexp.Compile("\\.")
// 	if err != nil {
// 		return "", "", err
// 	}

// 	key := reg.ReplaceAllString(domain, "_")

// 	currentPath := "./"
// 	if path != "" {
// 		if string(path[len(path)-1]) == "/" {
// 			currentPath = path[0 : len(path)-1]
// 		} else {
// 			currentPath = path
// 		}
// 	}

// 	if key == "" {
// 		return "", "", fmt.Errorf("key is error")
// 	}

// 	max := new(big.Int).Lsh(big.NewInt(1), 128)
// 	serialNumber, err := rand.Int(rand.Reader, max)
// 	if err != nil {
// 		return "", "", err
// 	}

// 	// 定义：引用IETF的安全领域的公钥基础实施（PKIX）工作组的标准实例化内容
// 	subject := pkix.Name{
// 		Organization: []string{domain},
// 		// OrganizationalUnit: []string{"client"},
// 		// CommonName:         "",
// 	}

// 	// 设置 SSL证书的属性用途
// 	certificate509 := x509.Certificate{
// 		SerialNumber: serialNumber,
// 		Subject:      subject,
// 		NotBefore:    time.Now(),
// 		NotAfter:     time.Now().Add(100 * 24 * time.Hour),
// 		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
// 		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
// 		// IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
// 	}

// 	// 生成指定位数密匙
// 	pk, err := rsa.GenerateKey(rand.Reader, 1024)
// 	if err != nil {
// 		return "", "", err
// 	}

// 	// 生成 SSL公匙
// 	derBytes, err := x509.CreateCertificate(rand.Reader, &certificate509, &certificate509, &pk.PublicKey, pk)
// 	if err != nil {
// 		return "", "", err
// 	}

// 	certPath := currentPath + "/" + key + "cert" + index + ".pem"
// 	certOut, err := os.Create(certPath)
// 	if err != nil {
// 		return "", "", err
// 	}
// 	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
// 	certOut.Close()

// 	keyPath := currentPath + "/" + key + "key" + index + ".pem"
// 	// 生成 SSL私匙
// 	keyOut, err := os.Create(keyPath)
// 	if err != nil {
// 		return "", "", err
// 	}
// 	pem.Encode(keyOut, &pem.Block{Type: "RAS PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)})
// 	keyOut.Close()
// 	return keyPath, certPath, nil
// }
