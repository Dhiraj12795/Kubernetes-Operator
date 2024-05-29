// api/v1alpha1/test_types.go

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TestSpec defines the desired state of Test
type TestSpec struct {
	Start      int      `json:"start"`
	End        int      `json:"end"`
	Replicas   int32    `json:"replicas"`
	Deployment []Deploy `json:"deployment"`
}

// Deploy defines the deployment details
type Deploy struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

// TestStatus defines the observed state of Test
type TestStatus struct {
	Status string `json:"status,omitempty"`
}

// Test is the Schema for the tests API
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type Test struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TestSpec   `json:"spec,omitempty"`
	Status TestStatus `json:"status,omitempty"`
}

// TestList contains a list of Test
// +kubebuilder:object:root=true
type TestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Test `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Test{}, &TestList{})
}
