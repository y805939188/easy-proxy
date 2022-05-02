package consts

type EaxyProxyPath = string

const (
	EasyProxyRootDirName                  EaxyProxyPath = "tmp-easy-proxy"
	TmpCaCertSaveDirName                  EaxyProxyPath = "tmp-ca"
	TmpUpdateCaCertificatesBashScriptName EaxyProxyPath = "update-ca-certificates"
	TmpLocalServiceBinaryName             EaxyProxyPath = "easy-proxy-service"
	TmpProxyInfoJsonDirName               EaxyProxyPath = "tmp-proxy-info"
)

type Ubuntu = string

const (
	UbuntuCaCertificatesPath Ubuntu = "/usr/local/share/ca-certificates"
	UbuntuSystemtRootCaPath  Ubuntu = "/etc/ssl/certs/ca-certificates.crt"
)

// TODO: other platform
