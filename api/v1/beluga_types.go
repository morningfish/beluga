/*
Copyright 2021 morningfish.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// BelugaSpec defines the desired state of Beluga
type CallBackHost struct {
	Hostname    string `json:"hostname,omitempty"`
	Path        string `json:"path,omitempty"`
	Rewrite     bool   `json:"rewrite,omitempty"`
	RewritePath string `json:"rewritePath,omitempty"`
}
type BelugaSpec struct {
	DeploymentSpec appsv1.DeploymentSpec         `json:"deploymentSpec"`
	ServiceSpec    corev1.ServiceSpec            `json:"serviceSpec,omitempty"`
	IngressSpec    networkingv1beta1.IngressSpec `json:"ingressSpec,omitempty"`
	AddBindHost    []string                      `json:"addBindHost,omitempty"`
	DelBindHost    []string                      `json:"delBindHost,omitempty"`
	SidecarPort    int32                         `json:"sidecarPort,omitempty"`
	UniqueName     string                        `json:"uniqueName"`
}

// BelugaStatus defines the observed state of Beluga
type BelugaStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// Beluga is the Schema for the belugas API
type Beluga struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BelugaSpec   `json:"spec,omitempty"`
	Status BelugaStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// BelugaList contains a list of Beluga
type BelugaList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Beluga `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Beluga{}, &BelugaList{})
}
