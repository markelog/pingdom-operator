package v1alpha1

import (
	"github.com/russellcardullo/go-pingdom/pingdom"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// HTTPCheck is representation of the http check
type HTTPCheck pingdom.HttpCheck

// ChecksSpec defines the desired state of Checks
// +k8s:openapi-gen=true
type ChecksSpec struct {
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file

	// You don't have to define this, in fact, it's better to use
	// kubernetes secret for these settings, you can set them up
	// as env variables - "PINGDOM_(USER|PASSWORD|KEY|BASE_URL)"
	User     string `json:"user"`
	Password string `json:"password"`
	Key      string `json:"key"`
	BaseURL  string `json:"base-url"`

	HTTP *HTTPCheck `json:"http"`
}

// ChecksStatus defines the observed state of Checks
// +k8s:openapi-gen=true
type ChecksStatus struct {
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	ID int
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Checks is the Schema for the checks API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=checks,scope=Namespaced
type Checks struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ChecksSpec   `json:"spec,omitempty"`
	Status ChecksStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ChecksList contains a list of Checks
type ChecksList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Checks `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Checks{}, &ChecksList{})
}
