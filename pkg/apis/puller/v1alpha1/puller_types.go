package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +kubebuilder:resource:scope="Cluster",singular="puller",path="pullers"
// +kubebuilder:subresource:status
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Puller is the Schema for the fast api
type Puller struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PullerSpec   `json:"spec,omitempty"`
	Status PullerStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PullerList contains a list of Puller
type PullerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Puller `json:"items"`
}

// PullerSpec defines the desired state of Puller
type PullerSpec struct {
	// +kubebuilder:validation:Optional
	Registries []Registry `json:"registries,omitempty"`

	// +kubebuilder:validation:Optional
	NamespaceAffinity *metav1.LabelSelector `json:"namespaceAffinity,omitempty"`
}

type Registry struct {
	// +kubebuilder:validation:Optional
	Server string `json:"server,omitempty"`

	// +kubebuilder:validation:Optional
	Username string `json:"username,omitempty"`

	// +kubebuilder:validation:Optional
	Password string `json:"password,omitempty"`

	// +kubebuilder:validation:Optional
	Email string `json:"email,omitempty"`

	// +kubebuilder:validation:Optional
	Auth string `json:"auth,omitempty"`
}

// PullerStatus defines the observed state of Puller
type PullerStatus struct {
	// +kubebuilder:validation:Optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}
