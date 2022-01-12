// +kubebuilder:object:generate=true
package v1beta1

import (
	hypv1alpha1 "github.com/openshift/hypershift/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

type ConditionType string

const (
	PlatformBeingConfigured         = "PlatformInfrastructureBeingConfigured"
	PlatformConfiguredAsExpected    = "PlatformInfrastructureConfiguredAsExpected"
	PlatfromDestroy                 = "PlatformInfrastructureDestroy"
	PlatformMisConfiguredReason     = "PlatformInfrastructureMisconfigured"
	PlatformIAMBeingConfigured      = "PlatformIAMBeingConfigured"
	PlatformIAMConfiguredAsExpected = "PlatformIAMConfiguredAsExpected"
	PlatformIAMRemove               = "PlatformIAMRemove"
	PlatformIAMMisConfiguredReason  = "PlatformIAMMisconfigured"

	// PlatformConfigured indicates (if status is true) that the
	// platform configuration specified for the platform provider has been deployed
	PlatformConfigured    ConditionType = "PlatformInfrastructureConfigured"
	PlatformIAMConfigured ConditionType = "PlatformIAMConfigured"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ClusterCuratorSpec defines the desired state of ClusterCurator
type ClusterCuratorSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// This is the desired curation that will occur
	// +kubebuilder:validation:Enum={install,scale,upgrade,destroy,delete-cluster-namespace}
	DesiredCuration string `json:"desiredCuration,omitempty"`

	// Points to the Cloud Provider or Ansible Provider secret, format: namespace/secretName
	ProviderCredentialPath string `json:"providerCredentialPath,omitempty"`

	// During an install curation run these Pre/Post hooks
	Install Hooks `json:"install,omitempty"`

	// During an scale curation run these Pre/Post hooks
	Scale Hooks `json:"scale,omitempty"`

	// During an destroy curation run these **Pre hook ONLY**
	Destroy Hooks `json:"destroy,omitempty"`

	// During an upgrade curation run these
	Upgrade UpgradeHooks `json:"upgrade,omitempty"`

	// Kubernetes job resource created for curation of a cluster
	CuratingJob string `json:"curatorJob,omitempty"`

	// Hypershift cluster definition, used to instantiated a HostedCluster and NodePools
	// +optional
	Hypershift *HypershiftSpec `json:"hypershift,omitempty"`
}

//Hypershift Specification
type HypershiftSpec struct {
	// Infrastructure instructions and pointers so either ClusterCurator generates what is needed or
	// skips it when the user provides the infrastructure values
	// +immutable
	Infrastructure InfraSpec `json:"infrastructure"`

	// Infrastructure ID, this is used to tag resources in the Cloud Provider, it will be generated
	// if not provided
	// +immutable
	// +optional
	InfraID string `json:"infra-id,omitempty"`

	// HostedCluster that will be applied to the ManagementCluster by ACM, if omitted, it will be generated
	// +optional
	HostedClusterSpec *hypv1alpha1.HostedClusterSpec `json:"hostedCluster,omitempty"`

	// NodePools is an array of NodePool resources that will be applied to the ManagementCluster by ACM,
	// if omitted, a default NodePool will be generated
	// +optional
	NodePools []*hypv1alpha1.NodePoolSpec `json:"nodePools,omitempty"`
}

type InfraSpec struct {
	// Configure the infrastructure using the provided CloudProvider, or user provided
	// +immutable
	Configure bool `json:"configure"`

	// Platform has infrastructure related information for the deployment, this is from the HostedCluster CRD
	// +optional
	Platform hypv1alpha1.PlatformSpec `json:"platform,omitempty"`

	// CloudProvider secret, contains the Cloud credenetial, Pull Secret and Base Domain
	// +immutable
	CloudProvider corev1.LocalObjectReference `json:"cloudProvider"`
}

type Hook struct {
	// Name of the Ansible Template to run in Tower as a job
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// Ansible job extra_vars is passed to the Ansible job at execution time
	// and is a known Ansible entity.
	// +kubebuilder:pruning:PreserveUnknownFields
	ExtraVars *runtime.RawExtension `json:"extra_vars,omitempty"`
}

type Hooks struct {

	// TowerAuthSecret is ansible secret used in template to run in tower
	// +kubebuilder:validation:Required
	TowerAuthSecret string `json:"towerAuthSecret,omitempty"`

	// Jobs to run before the cluster deployment
	Prehook []Hook `json:"prehook,omitempty"`

	// Jobs to run after the cluster import
	Posthook []Hook `json:"posthook,omitempty"`

	// When provided, this is a Job specification and overrides the default flow
	// +kubebuilder:pruning:PreserveUnknownFields
	OverrideJob *runtime.RawExtension `json:"overrideJob,omitempty"`
}

type UpgradeHooks struct {

	// TowerAuthSecret is ansible secret used in template to run in tower
	// +kubebuilder:validation:Required
	TowerAuthSecret string `json:"towerAuthSecret,omitempty"`

	// DesiredUpdate indicates the desired value of
	// the cluster version. Setting this value will trigger an upgrade (if
	// the current version does not match the desired version).
	// +optional
	DesiredUpdate string `json:"desiredUpdate,omitempty"`

	// Channel is an identifier for explicitly requesting that a non-default
	// set of updates be applied to this cluster. The default channel will be
	// contain stable updates that are appropriate for production clusters.
	// +optional
	Channel string `json:"channel,omitempty"`

	// Upstream may be used to specify the preferred update server. By default
	// it will use the appropriate update server for the cluster and region.
	// +optional
	Upstream string `json:"upstream,omitempty"`

	// Jobs to run before the cluster upgrade
	Prehook []Hook `json:"prehook,omitempty"`

	// Jobs to run after the cluster upgrade
	Posthook []Hook `json:"posthook,omitempty"`

	// When provided, this is a Job specification and overrides the default flow
	// +kubebuilder:pruning:PreserveUnknownFields
	OverrideJob *runtime.RawExtension `json:"overrideJob,omitempty"`
}

// ClusterCuratorStatus defines the observed state of ClusterCurator work
type ClusterCuratorStatus struct {
	// Track the conditions for each step in the desired curation that is being
	// executed as a job
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=clustercurator,shortName=cc;ccs,scope=Namespaced
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Reason",type="string",JSONPath=".status.conditions[?(@.type==\"PlatformInfrastructureConfigured\")].reason",description="Reason"
// +kubebuilder:printcolumn:name="IAM Ready",type="string",JSONPath=".status.conditions[?(@.type==\"PlatformIAMConfigured\")].status",description="Configured"
// +kubebuilder:printcolumn:name="Reason",type="string",JSONPath=".status.conditions[?(@.type==\"PlatformIAMConfigured\")].reason",description="Reason"

// +kubebuilder:printcolumn:name="INFRA Ready",type="string",JSONPath=".status.conditions[?(@.type==\"PlatformInfrastructureConfigured\")].status",description="Configured"
// ClusterCurator is the Schema for the clustercurators API
// This kind allows for prehook and posthook jobs to be executed prior to Hive provisioning
// and import of a cluster.
type ClusterCurator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterCuratorSpec   `json:"spec,omitempty"`
	Status ClusterCuratorStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ClusterCuratorList contains a list of ClusterCurator
type ClusterCuratorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterCurator `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ClusterCurator{}, &ClusterCuratorList{})
}
