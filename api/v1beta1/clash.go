package v1beta1

type Rule string

type SubServer struct {
	Port               int                    `yaml:"port"`
	SocksPort          int                    `yaml:"socks-port"`
	AllowLan           bool                   `yaml:"allow-lan"`
	Mode               string                 `yaml:"mode"`
	LogLevel           string                 `yaml:"log-level"`
	ExternalController string                 `yaml:"external-controller"`
	Experimental       map[string]interface{} `yaml:"experimental"`
	Proxy              []SubServerProxy       `yaml:"Proxy"`
	ProxyGroup         []SubServerProxyGroup  `yaml:"Proxy Group"`
	Rule               []Rule                 `yaml:"Rule"`
}

type SubServerProxy struct {
	Name     string `yaml:"name"`
	Type     string `yaml:"type"`
	Server   string `yaml:"server"`
	Port     int32  `yaml:"port"`
	Cipher   string `yaml:"cipher"`
	Password string `yaml:"password"`
}

type SubServerProxyList []SubServerProxy

type SubServerProxyGroup struct {
	Name    string   `yaml:"name"`
	Type    string   `yaml:"type"`
	Proxies []string `yaml:"proxies"`
}

type SubServerProxyGroupList []SubServerProxyGroup
