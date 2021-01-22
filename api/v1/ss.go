package v1

import (
	"encoding/base64"
	"fmt"
	"github.com/morningfish/beluga/api"
	"github.com/morningfish/beluga/api/config"
	"io/ioutil"
)

const SSFile = config.SubFilePath + "/beluga.ss"

type SSServer struct {
	ProxyServer []SSSubServerProxy
}

type SSSubServerProxy struct {
	Name     string `yaml:"name"`
	Type     string `yaml:"type"`
	Server   string `yaml:"server"`
	Port     int32  `yaml:"port"`
	Cipher   string `yaml:"cipher"`
	Password string `yaml:"password"`
}

func NewSSSubServer() api.SubServer {
	return &SSServer{}
}

func (s *SSServer) AddSubServerProxy(subServerProxy interface{}) {
	s.ProxyServer = append(s.ProxyServer, subServerProxy.(SSSubServerProxy))
}
func (s *SSServer) AddSubServerProxyGroup(subServerProxyGroup interface{}) {
}

func (s *SSServer) ToFile() error {
	data := ""
	for _, p := range s.ProxyServer {
		data += p.ToSS() + "\n"
	}
	encoded := base64.StdEncoding.EncodeToString([]byte(data))
	return ioutil.WriteFile(SSFile, []byte(encoded), 0644)
}

func (s *SSSubServerProxy) ToSS() string {
	data := fmt.Sprintf("%s:%s@%s:%d#%s", s.Cipher, s.Password, s.Server, s.Port, s.Name)
	encoded := base64.StdEncoding.EncodeToString([]byte(data))
	url := fmt.Sprintf("ss://%s", encoded)
	return url
}
