package v1alpha1

import (
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type VaultList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Vault `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Vault struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              VaultSpec   `json:"spec"`
	Status            VaultStatus `json:"status,omitempty"`
}

type VaultSpec struct {
	Size           int32                  `json:"size"`
	Image          string                 `json:"image"`
	Config         map[string]interface{} `json:"config"`
	ExternalConfig map[string]interface{} `json:"externalConfig"`
	UnsealConfig   UnsealConfig           `json:"unsealConfig"`
}

func (spec *VaultSpec) ConfigJSON() string {
	config, _ := json.Marshal(spec.Config)
	return string(config)
}

func (spec *VaultSpec) ExternalConfigJSON() string {
	config, _ := json.Marshal(spec.ExternalConfig)
	return string(config)
}

type VaultStatus struct {
	Nodes []string `json:"nodes"`
}

type UnsealConfig struct {
	Kubernetes *KubernetesUnsealConfig `json:"kubernetes"`
}

func (usc *UnsealConfig) ToArgs(vault *Vault) []string {
	if usc.Kubernetes != nil {
		secretNamespace := vault.Namespace
		if usc.Kubernetes.SecretNamespace != "" {
			secretNamespace = usc.Kubernetes.SecretNamespace
		}
		secretName := vault.Name + "-unseal-keys"
		if usc.Kubernetes.SecretName != "" {
			secretName = usc.Kubernetes.SecretName
		}
		return []string{"--mode", "k8s", "--k8s-secret-namespace", secretNamespace, "--k8s-secret-name", secretName}
	}
	return []string{}
}

type KubernetesUnsealConfig struct {
	SecretNamespace string `json:"secretNamespace"`
	SecretName      string `json:"secretName"`
}
