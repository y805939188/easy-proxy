package tools

import (
	"easy-proxy/consts"
	"os/user"
)

func GetUserHome() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}
	return user.HomeDir, nil
}

func GetUpdateSystemCaScriptPath() (string, error) {
	rootPaht, err := GetEasyRootPath()
	if err != nil {
		return "", err
	}
	return rootPaht + "/" + consts.TmpUpdateCaCertificatesBashScriptName, nil
}

func GetEasyRootPath() (string, error) {
	homeDir, err := GetUserHome()
	if err != nil {
		return "", err
	}
	return homeDir + "/" + consts.EasyProxyRootDirName, nil
}

func GetTmpCaPath() (string, error) {
	homeDir, err := GetUserHome()
	if err != nil {
		return "", err
	}
	return homeDir + "/" + consts.EasyProxyRootDirName + "/" + consts.TmpCaCertSaveDirName, nil
}

func GetTmpLocalServiceRootPath() (string, error) {
	rootPath, err := GetEasyRootPath()
	if err != nil {
		return "", err
	}
	return rootPath + "/" + consts.TmpLocalServiceBinaryName, nil
}

func GetTmpLocalServicePath() (string, error) {
	rootPath, err := GetTmpLocalServiceRootPath()
	if err != nil {
		return "", err
	}
	return rootPath + "/service" + "/" + consts.TmpLocalServiceBinaryName, nil
}

func GetTmpProxyInfoPath() (string, error) {
	rootPaht, err := GetEasyRootPath()
	if err != nil {
		return "", err
	}
	return rootPaht + "/" + consts.TmpProxyInfoJsonDirName, nil
}

func GetTmpProxyInfoFilePath(id string) (string, error) {
	proxyInfoPath, err := GetTmpProxyInfoPath()
	if err != nil {
		return "", err
	}
	fileName := proxyInfoPath + "/" + id + ".json"
	return fileName, nil
}
