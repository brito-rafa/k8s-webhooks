/*


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

package v2beta2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// RockBandSpec defines the desired state of RockBand
type RockBandSpec struct {
	// +kubebuilder:validation:Optional
	Genre string `json:"genre"`
	// +kubebuilder:validation:Optional
	NumberComponents int32 `json:"numberComponents"`
	// +kubebuilder:validation:Optional
	LeadSinger string `json:"leadSinger"`
	// +kubebuilder:validation:Optional
	LeadGuitar string `json:"leadGuitar"`
	// +kubebuilder:validation:Optional
	Drummer string `json:"drummer"`
}

// RockBandStatus defines the observed state of RockBand
type RockBandStatus struct {
	LastPlayed string `json:"lastPlayed"`
}

// +kubebuilder:validation:Optional
// +kubebuilder:resource:shortName={"rb"}
// +kubebuilder:printcolumn:name="Genre",type=string,JSONPath=`.spec.genre`
// +kubebuilder:printcolumn:name="Number_Components",type=integer,JSONPath=`.spec.numberComponents`
// +kubebuilder:printcolumn:name="Lead_Singer",type=string,JSONPath=`.spec.leadSinger`

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion

// RockBand is the Schema for the rockbands API
type RockBand struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RockBandSpec   `json:"spec,omitempty"`
	Status RockBandStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// RockBandList contains a list of RockBand
type RockBandList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RockBand `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RockBand{}, &RockBandList{})
}
