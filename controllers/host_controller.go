package controllers

import (
	"fmt"
	"github.com/morningfish/beluga/api/config"
	belugav1 "github.com/morningfish/beluga/api/v1"
	"github.com/morningfish/beluga/tools"
	"io/ioutil"
	corev1 "k8s.io/api/core/v1"
	"strings"
)

func InitBindHost() error {
	if tools.Exists(config.BindHostFile) {
		b, err := ioutil.ReadFile(config.BindHostFile)
		if err != nil {
			return err
		}
		config.BindHost = strings.Split(string(b), ",")
		err = InitRule()
		if err != nil {
			return err
		}
		return nil
	} else {
		return fmt.Errorf("bind host file not exist")
	}
}

// 附加host
func InjectHost(instance *belugav1.Beluga) error {
	hostMap := make(map[string]int)
	for _, host := range config.BindHost {
		if _, ok := hostMap[host]; !ok {
			hostMap[host] = 1
		}
	}
	for _, addHost := range instance.Spec.AddBindHost {
		if _, ok := hostMap[addHost]; !ok {
			hostMap[addHost] = 1
		}
	}
	for _, delHost := range instance.Spec.DelBindHost {
		if _, ok := hostMap[delHost]; ok {
			delete(hostMap, delHost)
		}
	}
	hostAlias := corev1.HostAlias{
		IP:        "127.0.0.1",
		Hostnames: []string{},
	}
	for host, _ := range hostMap {
		hostAlias.Hostnames = append(hostAlias.Hostnames, host)
	}
	instance.Spec.DeploymentSpec.Template.Spec.HostAliases = append(instance.Spec.DeploymentSpec.Template.Spec.HostAliases, hostAlias)
	return nil
}
