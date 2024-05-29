package controllers

import (
	"context"
	"fmt"
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
	log.Info(fmt.Sprintf("currentTime %d", currentHour))

	if currentHour >= startTime && currentHour <= endTime {
		err = scaleDeployment(test, r, ctx, replicas)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
}

func scaleDeployment(test *myworkspotinv1alpha1.Test, r *TestReconciler, ctx context.Context, replicas int32) error {
	for _, deploy := range test.Spec.Deployment {
		dep := &appsv1.Deployment{}
		err := r.Get(ctx, types.NamespacedName{
			Namespace: deploy.Namespace,
			Name:      deploy.Name,
		}, dep)
		if err != nil {
			return err
		}

		if *dep.Spec.Replicas != replicas {
			dep.Spec.Replicas = &replicas
			err = r.Update(ctx, dep)
			if err != nil {
				test.Status.Status = myworkspotinv1alpha1.FAILED
				return err
			}
		}
	}

	test.Status.Status = myworkspotinv1alpha1.SUCCESS
	err := r.Status().Update(ctx, test)
	if err != nil {
		return err
	}

	return nil
}

func (r *TestReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&myworkspotinv1alpha1.Test{}).
		Complete(r)
}
