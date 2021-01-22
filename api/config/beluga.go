package config

const (
	BelugaProxyType   = "ss"
	BelugaProxyCipher = "aes-256-cfb"
	BelugaPassword    = "Beluga"
	SubFilePath       = "/tmp/subfiles"
	SubFile           = SubFilePath + "/beluga.yaml"
)

var (
	RuleFile               string
	Rules                  []string
	BindHostFile           string
	BindHost               []string
	SubServerListenAddress string
)
