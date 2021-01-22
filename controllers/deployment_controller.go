package controllers

import (
	"context"
	belugav1 "github.com/morningfish/beluga/api/v1"
	appv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// 配置 Deployment
func (r *BelugaReconciler) reconcileDeployment(instance *belugav1.Beluga) (reconcile.Result, error) {
	r.Log.Info("Reconciling Beluga Deployment")
	var err error
	// Define a new Deployment object
	deployment := r.newDeploymentForCR(instance)
	// Set Beluga instance as the owner and controller
	if err = controllerutil.SetControllerReference(instance, deployment, r.Scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this deployment already exists
	found := &appv1.Deployment{}
	err = r.Get(context.TODO(), types.NamespacedName{Name: deployment.Name, Namespace: deployment.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		r.Log.Info("Creating a new deployment", "deployment.Namespace", deployment.Namespace, "deployment.Name", deployment.Name)
		err = r.Create(context.TODO(), deployment)
		if err != nil {
			return reconcile.Result{}, err
		}
	} else if err != nil {
		return reconcile.Result{}, err
	} else {
		err = r.Update(context.TODO(), deployment)
		if err != nil {
			return reconcile.Result{}, err
		}
		r.Log.Info("Skip reconcile: deployment already exists", "deployment.Namespace", found.Namespace, "deployment.Name", found.Name)
	}

	return reconcile.Result{}, nil
}

// 返回一个 pod 地址
func (r *BelugaReconciler) newDeploymentForCR(instance *belugav1.Beluga) *appv1.Deployment {

	return &appv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
			Labels:    instance.Labels,
		},
		Spec: instance.Spec.DeploymentSpec,
	}
}
