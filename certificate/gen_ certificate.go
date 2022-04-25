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
	"io/ioutil"
	"math/big"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func ParseLocalCertStringToCertificate(certPath string) (*x509.Certificate, error) {
	currentCert, err := ioutil.ReadFile(certPath)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode([]byte(currentCert))
	if block == nil {
		return nil, fmt.Errorf("将证书文本转为 block 失败")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	return cert, nil
}

func RemoveCertificateFromSystemByCertName(certName string) error {
	if certName == "" {
		return nil
	}
	currentCertName := strings.Replace(certName, ".pem", "", 1)
	currentCertName = strings.Replace(currentCertName, ".crt", "", 1)
	currentCertName += ".crt"
	certPath := consts.UbuntuCaCertificatesPath + "/" + currentCertName
	err := tools.DeleteFile(certPath)
	if err != nil {
		return err
	}
	err = ExecCertificateUpdateBash()
	if err != nil {
		return err
	}
	return nil
}

func ExecCertificateUpdateBash() error {
	// cwd, err := os.Getwd()
	// if err != nil {
	// 	return err
	// }
	// file := cwd + "/certificate" + "/update-ca-certificates"
	file := "./update-ca-certificates"
	out, code := tools.Bash("chmod u+x " + file)
	if code != 0 {
		return fmt.Errorf(out)
	}
	fmt.Println(out)
	out, code = tools.Bash(file)
	if code != 0 {
		return fmt.Errorf(out)
	}
	fmt.Println(out)
	return nil
}

func GetSystemCertNameFromPath(certPath string) string {
	fileName := filepath.Base(certPath)
	ext := path.Ext(fileName)
	currentFileName := strings.Replace(fileName, ext, ".crt", 1)
	return currentFileName
}

func SetCertificateToSystemByCertPath(certPath string) error {
	currentFileName := GetSystemCertNameFromPath(certPath)
	err := tools.FileCopy(certPath, consts.UbuntuCaCertificatesPath+"/"+currentFileName)
	if err != nil {
		return err
	}

	err = ExecCertificateUpdateBash()
	if err != nil {
		return err
	}
	return nil
}

// https://serverfault.com/questions/18364/can-i-use-the-same-wildcard-certification-for-domain-com-and-domain-com
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
