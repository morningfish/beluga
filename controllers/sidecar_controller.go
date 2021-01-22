package controllers

import (
	"fmt"
	belugav1 "github.com/morningfish/beluga/api/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"reflect"
)

const (
	SidecarName         = "beluga-sidecar"
	SidecarImage        = "shadowsocks/shadowsocks-libev:v3.3.4"
	SidecarPort         = 9123
	SidecarPassword     = "beluga"
	SidecarVolumeName   = "beluga-sidecar-config"
	SidecarVolumePath   = "/ss-server-config.json"
	SidecarConfigMapKey = "ss-server-config.json"
)

var (
	SidecarConfigMapData = fmt.Sprintf(`{
  "server":"0.0.0.0",
  "server_port":%d,
  "password":%s,
  "timeout":600,
  "method":"aes-256-cfb",
  "workers": 1
}`, SidecarPort, SidecarPassword)
)

// InjectSidecar: inject sidecar instance
func InjectSidecar(instance *belugav1.Beluga) error {
	injectDeployment(instance)
	injectService(instance)
	return nil
}
func injectDeployment(instance *belugav1.Beluga) {
	sidecarContainer := newSidecarContainer()
	instance.Spec.DeploymentSpec.Template.Spec.Containers = append(instance.Spec.DeploymentSpec.Template.Spec.Containers, *sidecarContainer)
}
func newSidecarContainer() *corev1.Container {
	return &corev1.Container{
		Name:            SidecarName,
		Image:           SidecarImage,
		ImagePullPolicy: corev1.PullIfNotPresent,
		Ports: []corev1.ContainerPort{
			corev1.ContainerPort{
				Name:          SidecarName,
				ContainerPort: SidecarPort,
				Protocol:      corev1.ProtocolTCP,
			},
		},
		Command: []string{"/bin/sh"},
		Args:    []string{"-c", fmt.Sprintf("ss-server -c %s", SidecarVolumePath)},
		VolumeMounts: []corev1.VolumeMount{
			corev1.VolumeMount{
				Name:      SidecarVolumeName,
				MountPath: SidecarVolumePath,
				SubPath:   SidecarConfigMapKey,
			},
		},
	}
}
func injectService(instance *belugav1.Beluga) {
	servicePort := newServicePort(instance)
	if reflect.DeepEqual(instance.Spec.ServiceSpec, corev1.ServiceSpec{}) {
		instance.Spec.ServiceSpec = corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				*servicePort,
				corev1.ServicePort{
					Name:     fmt.Sprintf("%s80", instance.Name),
					Protocol: corev1.ProtocolTCP,
					Port:     80,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 80,
					},
				},
				corev1.ServicePort{
					Name:     fmt.Sprintf("%s443", instance.Name),
					Protocol: corev1.ProtocolTCP,
					Port:     443,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 443,
					},
				},
			},
			Selector: instance.Labels,
			Type:     corev1.ServiceTypeNodePort,
		}
	} else {
		instance.Spec.ServiceSpec.Ports = append(instance.Spec.ServiceSpec.Ports, *servicePort)
		instance.Spec.ServiceSpec.Type = corev1.ServiceTypeNodePort
	}
}

func newServicePort(instance *belugav1.Beluga) *corev1.ServicePort {
	return &corev1.ServicePort{
		Name:     SidecarName,
		Protocol: corev1.ProtocolTCP,
		Port:     SidecarPort,
		TargetPort: intstr.IntOrString{
			Type:   intstr.Int,
			IntVal: SidecarPort,
		},
		NodePort: instance.Spec.SidecarPort,
	}
}
