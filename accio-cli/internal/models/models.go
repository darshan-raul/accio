package models

type CloudProvider string

const (
	AWS   CloudProvider = "aws"
	GCP   CloudProvider = "gcp"
	Azure CloudProvider = "azure"
)

type CapabilityType string

const (
	CapabilityCompute    CapabilityType = "compute"
	CapabilityKubernetes CapabilityType = "kubernetes"
	CapabilityDatabase   CapabilityType = "database"
	CapabilityStorage    CapabilityType = "object-storage"
	CapabilityNetworking CapabilityType = "networking"
)

type InfraIntent struct {
	APIVersion string     `yaml:"apiVersion" json:"apiVersion"`
	Kind       string     `yaml:"kind" json:"kind"`
	Metadata   Metadata   `yaml:"metadata" json:"metadata"`
	Spec       IntentSpec `yaml:"spec" json:"spec"`
}

type Metadata struct {
	Name string `yaml:"name" json:"name"`
}

type IntentSpec struct {
	Cloud        CloudProvider    `yaml:"cloud" json:"cloud"`
	Environment  string           `yaml:"environment" json:"environment"`
	Capabilities []CapabilitySpec `yaml:"capabilities" json:"capabilities"`
}

type CapabilitySpec struct {
	Type   CapabilityType `yaml:"type" json:"type"`
	Size   string         `yaml:"size" json:"size"`
	Region string         `yaml:"region" json:"region"`
	// Add more fields as needed for the wizard
}

type StackStatus string

const (
	StatusReady       StackStatus = "Ready"
	StatusReconciling StackStatus = "Reconciling"
	StatusFailed      StackStatus = "Failed"
)

type Stack struct {
	Name      string      `json:"name"`
	Cloud     string      `json:"cloud"`
	Region    string      `json:"region"`
	Status    StackStatus `json:"status"`
	Resources []Resource  `json:"resources"`
}

type Resource struct {
	Name   string `json:"name"`
	Kind   string `json:"kind"`
	Status string `json:"status"`
}
