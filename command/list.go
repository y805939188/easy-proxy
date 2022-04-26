package command

import (
	"easy-proxy/tools"
	"encoding/json"
	"fmt"
	"path"
	"strings"
)

type ListProxy struct{}

func (p *ListProxy) List() error {
	rootPath, err := tools.GetTmpProxyInfoPath()
	if err != nil {
		return err
	}
	files, err := tools.GetAllFile(rootPath)
	if err != nil {
		return err
	}

	res := ""
	for _, file := range files {
		fileName := path.Base(file)
		fileExt := path.Ext(file)
		fileName = strings.Replace(fileName, fileExt, "", 1)
		content, err := tools.ReadFile(file)
		if err != nil {
			return err
		}
		var proxyInfo Rule
		err = json.Unmarshal(content, &proxyInfo)
		if err != nil {
			return err
		}
		identifier := proxyInfo.Identifier
		if identifier.SourcePort != "" && identifier.TargetPort != "" {
			res += fmt.Sprintf(
				"%s      %s:%s      to      %s:%s",
				fileName,
				identifier.SourceAddress,
				identifier.SourcePort,
				identifier.TargetAddress,
				identifier.TargetPort,
			) + "\r\n"
		} else if identifier.SourcePort == "" && identifier.TargetPort != "" {
			res += fmt.Sprintf(
				"%s      %s      to      %s:%s",
				fileName,
				identifier.SourceAddress,
				identifier.TargetAddress,
				identifier.TargetPort,
			) + "\r\n"
		} else {
			res += fmt.Sprintf(
				"%s      %s:%s      to      %s",
				fileName,
				identifier.SourceAddress,
				identifier.SourcePort,
				identifier.TargetAddress,
			) + "\r\n"
		}
	}

	fmt.Println(res)
	return nil
}

func GetListProxy() (*ListProxy, error) {
	p := &ListProxy{}
	return p, nil
}
