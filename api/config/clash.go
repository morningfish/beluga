package config

const (
	Port               = 7890
	SocksPort          = 7891
	Mode               = "Rule"
	AllowLan           = true
	LogLevel           = "info"
	ExternalController = "0.0.0.0:9090"
	BelugaProxyType    = "ss"
	BelugaProxyCipher  = "aes-256-cfb"
	SubFilePath        = "/tmp/subfiles"
	SubFile            = SubFilePath + "/beluga.yaml"
)
