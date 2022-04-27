package command

import "easy-proxy/tools"

type FreshProxy struct {
	delProxy *DelProxy
}

func (p *FreshProxy) FreshAll() error {
	dirPath, err := tools.GetTmpProxyInfoPath()
	if err != nil {
		return err
	}
	files, err := tools.GetAllFile(dirPath)
	if err != nil {
		return err
	}

	ids := []string{}
	for _, file := range files {
		name, _ := tools.GetFileNameAndExt(file)
		ids = append(ids, name)
	}

	return p.delProxy.DeleteProxys(ids...)
}

func GetFreshProxy() (*FreshProxy, error) {
	delProxy, err := GetDelProxy()
	if err != nil {
		return nil, err
	}
	p := &FreshProxy{
		delProxy: delProxy,
	}
	return p, nil
}
