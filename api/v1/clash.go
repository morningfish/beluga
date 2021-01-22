package v1

import (
	"github.com/morningfish/beluga/api"
	"github.com/morningfish/beluga/api/config"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

const ClashV1File = config.SubFilePath + "/beluga_v1.yaml"

// clash server 类型结构体
type ClashSubServer struct {
	Port               int                        `yaml:"port"`
	SocksPort          int                        `yaml:"socks-port"`
	AllowLan           bool                       `yaml:"allow-lan"`
	Mode               string                     `yaml:"mode"`
	LogLevel           string                     `yaml:"log-level"`
	ExternalController string                     `yaml:"external-controller"`
	Experimental       map[string]interface{}     `yaml:"experimental"`
	Proxy              []ClashSubServerProxy      `yaml:"proxies"`
	ProxyGroup         []ClashSubServerProxyGroup `yaml:"proxy-groups"`
	Rule               []string                   `yaml:"rules"`
}

type ClashSubServerProxy struct {
	Name     string `yaml:"name"`
	Type     string `yaml:"type"`
	Server   string `yaml:"server"`
	Port     int32  `yaml:"port"`
	Cipher   string `yaml:"cipher"`
	Password string `yaml:"password"`
}

type ClashSubServerProxyGroup struct {
	Name    string   `yaml:"name"`
	Type    string   `yaml:"type"`
	Proxies []string `yaml:"proxies"`
}

func NewClashSubServer() api.SubServer {
	return &ClashSubServer{
		Port:               config.ClashPort,
		SocksPort:          config.ClashSocksPort,
		AllowLan:           config.ClashAllowLan,
		Mode:               config.ClashMode,
		LogLevel:           config.ClashLogLevel,
		ExternalController: config.ClashExternalController,
		Experimental: map[string]interface{}{
			"ignore-resolve-fail": true,
		},
		Rule: config.Rules,
	}
}

func (s *ClashSubServer) AddSubServerProxy(subServerProxy interface{}) {
	s.Proxy = append(s.Proxy, subServerProxy.(ClashSubServerProxy))
}
func (s *ClashSubServer) AddSubServerProxyGroup(subServerProxyGroup interface{}) {
	s.ProxyGroup = append(s.ProxyGroup, subServerProxyGroup.(ClashSubServerProxyGroup))
}

func (s *ClashSubServer) ToFile() error {
	data, err := yaml.Marshal(s)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(ClashV1File, data, 0644)
}
