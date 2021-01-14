package controllers

import (
	"context"
	belugav1 "github.com/morningfish/beluga/api/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// 配置 Deployment
func (r *BelugaReconciler) reconcileService(instance *belugav1.Beluga) (reconcile.Result, error) {
	r.Log.Info("Reconciling AvatarAgency Deployment")
	var err error
	// Define a new Deployment object
	service := newServiceForCR(instance)
	// Set AvatarAgency instance as the owner and controller
	if err = controllerutil.SetControllerReference(instance, service, r.Scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this deployment already exists
	found := &corev1.Service{}
	err = r.Get(context.TODO(), types.NamespacedName{Name: service.Name, Namespace: service.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		r.Log.Info("Creating a new deployment", "service.Namespace", service.Namespace, "service.Name", service.Name)
		err = r.Create(context.TODO(), service)
		if err != nil {
			return reconcile.Result{}, err
		}
	} else if err != nil {
		return reconcile.Result{}, err
	} else {
		err = r.Update(context.TODO(), service)
		if err != nil {
			return reconcile.Result{}, err
		}
		r.Log.Info("Skip reconcile: service already exists", "service.Namespace", found.Namespace, "service.Name", found.Name)
	}

	return reconcile.Result{}, nil
}
func newServiceForCR(instance *belugav1.Beluga) *corev1.Service {
	service := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
			Labels:    instance.Labels,
		},
		Spec: instance.Spec.ServiceSpec,
	}
	return &service
}
