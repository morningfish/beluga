package controllers

import (
	"context"
	"fmt"
	"github.com/morningfish/beluga/api"
	"github.com/morningfish/beluga/api/config"
	v1 "github.com/morningfish/beluga/api/v1"
	"github.com/morningfish/beluga/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"sort"
)

var (
	NodeIp        string
	SubServerList []api.SubServer
)

func registerSubServer() {
	SubServerList = []api.SubServer{
		v1beta1.NewClashSubServer(),
		v1.NewClashSubServer(),
		v1.NewSSSubServer(),
	}
}

func GenerateSubData() error {
	registerSubServer()
	nodeAddress, err := getNodeIp()
	if err != nil {
		return err
	}
	belugaList := &v1.BelugaList{}
	err = BelugaReconcile.List(context.TODO(), belugaList)
	if err != nil {
		return err
	}
	var (
		proxyGroupProxies []string
	)
	for _, beluga := range belugaList.Items {
		for _, subServer := range SubServerList {
			subServer.AddSubServerProxy(v1.ClashSubServerProxy{
				Name:     beluga.Spec.UniqueName,
				Type:     config.BelugaProxyType,
				Server:   nodeAddress,
				Port:     beluga.Spec.SidecarPort,
				Cipher:   config.BelugaProxyCipher,
				Password: config.BelugaPassword,
			})
		}
		proxyGroupProxies = append(proxyGroupProxies, beluga.Spec.UniqueName)
	}
	sort.Strings(proxyGroupProxies)
	for _, subServer := range SubServerList {
		subServer.AddSubServerProxyGroup(v1.ClashSubServerProxyGroup{
			Name:    "Proxy",
			Type:    "select",
			Proxies: proxyGroupProxies,
		})
	}
	for _, subServer := range SubServerList {
		err = subServer.ToFile()
		if err != nil {
			return err
		}
	}
	return nil
}

func getNodeIp() (string, error) {
	if NodeIp != "" {
		return NodeIp, nil
	}
	err := InitNodeIP()
	return NodeIp, err
}
func InitNodeIP() error {
	nodeList := &corev1.NodeList{}
	err := BelugaReconcile.List(context.TODO(), nodeList)
	if err != nil {
		return err
	}
	nodes := nodeList.Items
	for _, node := range nodes {
		for _, conditions := range node.Status.Conditions {
			if conditions.Type == corev1.NodeReady {
				if conditions.Status == corev1.ConditionTrue {
					for _, address := range node.Status.Addresses {
						if address.Type == corev1.NodeInternalIP {
							NodeIp = address.Address
							return nil
						}
					}
				}
			}
		}
	}
	return fmt.Errorf("no available node")
}
