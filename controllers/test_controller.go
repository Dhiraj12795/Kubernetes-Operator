/*
Copyright 2024 Dhiraj Bhattad.

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

package controllers

import (
	"context"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	myworkspotinv1alpha1 "github.com/Dhiraj12795/Kubernetes-Operator/api/v1alpha1"
)

// TestReconciler reconciles a Test object
type TestReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=myworkspot.in.myworkspot.in,resources=tests,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=myworkspot.in.myworkspot.in,resources=tests/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=myworkspot.in.myworkspot.in,resources=tests/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO: Modify the Reconcile function to compare the state specified by
// the Test object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *TestReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Reconcile called", "Request.Namespace", req.Namespace, "Request.Name", req.Name)

	test := &myworkspotinv1alpha1.Test{}
	err := r.Get(ctx, req.NamespacedName, test)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	startTime := test.Spec.Start
	endTime := test.Spec.End
	replicas := test.Spec.Replicas

	currentHour := time.Now().UTC().Hour()

	if currentHour >= startTime && currentHour <= endTime {
		for _, deploy := range test.Spec.Deployment {
			deployment := &appsv1.Deployment{}
			err := r.Get(ctx, types.NamespacedName{
				Namespace: deploy.Namespace,
				Name:      deploy.Name,
			}, deployment)
			if err != nil {
				return ctrl.Result{}, err
			}

			if *deployment.Spec.Replicas != replicas {
				deployment.Spec.Replicas = &replicas
				err := r.Update(ctx, deployment)
				if err != nil {
					return ctrl.Result{}, err
				}
			}
		}
	}

	return ctrl.Result{RequestAfter: time.Duration(30 * time.Second)}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TestReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&myworkspotinv1alpha1.Test{}).
		Complete(r)
}