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
	rootPaht, err := GetEasyRootPath()
	if err != nil {
		return "", err
	}
	return rootPaht + "/" + consts.TmpLocalServiceBinaryName, nil
}

func GetTmpLocalServicePath() (string, error) {
	rootPaht, err := GetTmpLocalServiceRootPath()
	if err != nil {
		return "", err
	}
	return rootPaht + "/service" + "/" + consts.TmpLocalServiceBinaryName, nil
}
