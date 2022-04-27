package certificate_test

import (
	"crypto/x509"
	"easy-proxy/certificate"
	"easy-proxy/consts"
	"easy-proxy/tools"
	"fmt"
	"testing"
)

func TestCertificate(t *testing.T) {
	test := "127.0.0.1"
	keyPath, certPath, erro := certificate.GenCertificate(test, "./")
	if erro != nil {
		panic("生成证书出错, err: " + erro.Error())
	}

	erro = certificate.SetCertificateToSystemByCertPath(certPath)
	if erro != nil {
		panic("将证书添加到系统受信任失败, err: " + erro.Error())
	}

	return

	testDomain := "www.ding-test.com"

	keyPath, certPath, err := certificate.GenCertificate(testDomain, "./")
	if err != nil {
		panic("生成证书出错, err: " + err.Error())
	}
	exist := tools.FileIsExisted(keyPath)
	if !exist {
		panic("证书私钥文件不存在")
	}
	exist = tools.FileIsExisted(certPath)
	if !exist {
		panic("证书不存在")
	}

	defer func() {
		fmt.Println("即将退出并删除证书")
		tools.DeleteFile(keyPath)
		tools.DeleteFile(certPath)
	}()

	err = certificate.SetCertificateToSystemByCertPath(certPath)
	if err != nil {
		panic("将证书添加到系统受信任失败, err: " + err.Error())
	}

	currentFileName := certificate.GetSystemCertNameFromPath(certPath)
	systemCertPath := consts.UbuntuCaCertificatesPath + "/" + currentFileName

	exist = tools.FileIsExisted(systemCertPath)
	if !exist {
		panic("将证书添加到 " + consts.UbuntuCaCertificatesPath + " 失败")
	}

	systemCertPool, err := x509.SystemCertPool()
	if err != nil {
		fmt.Println("通过 x509 获取系统证书池失败, err: ", err.Error())
		err = certificate.RemoveCertificateFromSystemByCertName(currentFileName)
		if err != nil {
			panic("从系统中移除证书失败, err: " + err.Error())
		}
		return
	}

	cert, err := certificate.ParseLocalCertStringToCertificate(certPath)
	if err != nil {
		fmt.Println("将文本证书转为结构体失败, err: ", err.Error())
		err = certificate.RemoveCertificateFromSystemByCertName(currentFileName)
		if err != nil {
			panic("从系统中移除证书失败, err: " + err.Error())
		}
		return
	}

	opt := x509.VerifyOptions{
		DNSName: testDomain,
		Roots:   systemCertPool,
	}
	_, err = cert.Verify(opt)
	if err != nil {
		fmt.Println("校验证书失败, err: ", err.Error())
		err = certificate.RemoveCertificateFromSystemByCertName(currentFileName)
		if err != nil {
			panic("从系统中移除证书失败, err: " + err.Error())
		}
		return
	}
	fmt.Println("成功将证书添加为系统受信任")
	err = certificate.RemoveCertificateFromSystemByCertName(currentFileName)
	if err != nil {
		panic("从系统中移除证书失败, err: " + err.Error())
	}
}
