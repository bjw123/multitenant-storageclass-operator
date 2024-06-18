/*
Copyright 2024.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// NSStorageClassSpec defines the desired state of NSStorageClass
type NSStorageClassSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Provisioner string `json:"provisioner"`
	//+optional
	Parameters map[string]string `json:"parameters,omitempty"`
	//+optional
	ReclaimPolicy ReclaimPolicy `json:"reclaimPolicy,omitempty"`
	//+optional
	MountOptions []string `json:"mountOptions,omitempty"`
	//+optional
	AllowVolumeExpansion *bool `json:"allowVolumeExpansion,omitempty"`
	//+optional
	VolumeBindingMode *VolumeBindingMode `json:"volumeBindingMode,omitempty"`
}

type ReclaimPolicy string

const (
	Retain ReclaimPolicy = "Retain"
	Delete ReclaimPolicy = "Delete" //default
	//Recycle ReclaimPolicy = "Recycle" //depreciated

)

// VolumeBindingMode indicates how PersistentVolumeClaims should be bound.
type VolumeBindingMode string

const (
	// VolumeBindingImmediate indicates that PersistentVolumeClaims should be
	// immediately provisioned and bound.
	VolumeBindingImmediate VolumeBindingMode = "Immediate"

	// VolumeBindingWaitForFirstConsumer indicates that PersistentVolumeClaims
	// should not be provisioned and bound until the first Pod is created that
	// references the PeristentVolumeClaim.  The volume provisioning and
	// binding will occur during Pod scheduing.
	VolumeBindingWaitForFirstConsumer VolumeBindingMode = "WaitForFirstConsumer"
)

// NSStorageClassStatus defines the observed state of NSStorageClass
type NSStorageClassStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// NSStorageClass is the Schema for the nsstorageclasses API
type NSStorageClass struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              NSStorageClassSpec   `json:"spec,omitempty"`
	Status            NSStorageClassStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// NSStorageClassList contains a list of NSStorageClass
type NSStorageClassList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NSStorageClass `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NSStorageClass{}, &NSStorageClassList{})
}
