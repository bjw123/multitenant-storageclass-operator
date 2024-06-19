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

package controller

import (
	"context"
	v1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	multitenantwrapperv1 "multitenant.storageclass/namespaced-wrapper/api/v1"
)

// NSStorageClassReconciler reconciles a NSStorageClass object
type NSStorageClassReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=multitenant-wrapper.multitenant.storageclass,resources=nsstorageclasses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=multitenant-wrapper.multitenant.storageclass,resources=nsstorageclasses/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=multitenant-wrapper.multitenant.storageclass,resources=nsstorageclasses/finalizers,verbs=update
//permissions to create storage class
//+kubebuilder:rbac:groups=storage.k8s.io/v1,resources=storageclass,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=storage.k8s.io/v1,resources=storageclass/ownerreferences,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the NSStorageClass object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.3/pkg/reconcile
func (r *NSStorageClassReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// check if namespaced storageclass exists
	var namespacedStorageClass multitenantwrapperv1.NSStorageClass
	if err := r.Get(ctx, req.NamespacedName, &namespacedStorageClass); err != nil {
		logger.Error(err, "unable to fetch Namespaced Storage Class")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	// check if  storageClass exists, create if needed
	var sc storagev1.StorageClass
	scName := types.NamespacedName{Name: namespacedStorageClass.Namespace + namespacedStorageClass.Name, Namespace: ""} //storage class is not bound to a namespace
	if err := r.Get(ctx, scName, &sc); err != nil {
		//create storageClass if does not exist
		sc = createStorageClass(namespacedStorageClass)
		err := controllerutil.SetOwnerReference(&namespacedStorageClass, &sc, r.Scheme)
		if err != nil {
			logger.Error(err, "unable to add owner ref to storageClass")
			return ctrl.Result{}, err
		}
		if err = r.Create(ctx, &sc); err != nil {
			logger.Error(err, "unable to create storageClass")
			return ctrl.Result{}, err
		}
	}

	// finalizer logic to clean up storageClass on deletion
	finalizer := "multitenant.storageclass/finalizer"
	if namespacedStorageClass.ObjectMeta.DeletionTimestamp.IsZero() {
		if !controllerutil.ContainsFinalizer(&namespacedStorageClass, finalizer) {
			//not being deleted, lets add the finalizer
			controllerutil.AddFinalizer(&namespacedStorageClass, finalizer)
			if err := r.Update(ctx, &namespacedStorageClass); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		//object is being deleted
		if controllerutil.ContainsFinalizer(&namespacedStorageClass, finalizer) {
			if err := r.Delete(ctx, &sc); err != nil {
				logger.Error(err, "unable to delete storageclass bound to namespaced storage class")
				return ctrl.Result{}, err
			}
		}
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NSStorageClassReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&multitenantwrapperv1.NSStorageClass{}).
		Complete(r)
}

func createStorageClass(nsc multitenantwrapperv1.NSStorageClass) storagev1.StorageClass {
	var pol v1.PersistentVolumeReclaimPolicy
	if nsc.Spec.ReclaimPolicy == "Reclaim" {
		pol = v1.PersistentVolumeReclaimRetain
	} else {
		pol = v1.PersistentVolumeReclaimDelete
	}

	storageClass := storagev1.StorageClass{
		TypeMeta:             metav1.TypeMeta{},
		ObjectMeta:           metav1.ObjectMeta{},
		Provisioner:          nsc.Spec.Provisioner,
		Parameters:           nsc.Spec.Parameters,
		ReclaimPolicy:        &pol,
		MountOptions:         nsc.Spec.MountOptions,
		AllowVolumeExpansion: nsc.Spec.AllowVolumeExpansion,
		VolumeBindingMode:    (*storagev1.VolumeBindingMode)(nsc.Spec.VolumeBindingMode),
		AllowedTopologies:    nil,
	}
	ownerRef := metav1.OwnerReference{Name: nsc.Name, Kind: nsc.Kind, UID: nsc.UID, APIVersion: nsc.APIVersion}
	storageClass.OwnerReferences = []metav1.OwnerReference{ownerRef}
	return storageClass
}
