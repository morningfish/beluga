package controllers

import (
	"context"
	belugav1 "github.com/morningfish/beluga/api/v1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// 配置 Deployment
func (r *BelugaReconciler) reconcileIngress(instance *belugav1.Beluga) (reconcile.Result, error) {
	r.Log.Info("Reconciling Beluga Deployment")
	var err error
	// Define a new Deployment object
	ingress := newIngressForCR(instance)
	if ingress == nil {
		return reconcile.Result{}, nil
	}
	// Set Beluga instance as the owner and controller
	if err = controllerutil.SetControllerReference(instance, ingress, r.Scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this deployment already exists
	found := &networkingv1beta1.Ingress{}
	err = r.Get(context.TODO(), types.NamespacedName{Name: ingress.Name, Namespace: ingress.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		r.Log.Info("Creating a new deployment", "ingress.Namespace", ingress.Namespace, "ingress.Name", ingress.Name)
		err = r.Create(context.TODO(), ingress)
		if err != nil {
			return reconcile.Result{}, err
		}
	} else if err != nil {
		return reconcile.Result{}, err
	} else {
		err = r.Update(context.TODO(), ingress)
		if err != nil {
			return reconcile.Result{}, err
		}
		r.Log.Info("Skip reconcile: service already exists", "ingress.Namespace", found.Namespace, "ingress.Name", found.Name)
	}

	return reconcile.Result{}, nil
}
func newIngressForCR(instance *belugav1.Beluga) *networkingv1beta1.Ingress {
	if reflect.DeepEqual(instance.Spec.IngressSpec, networkingv1beta1.IngressSpec{}) {
		return nil
	}
	return &networkingv1beta1.Ingress{
		ObjectMeta: instance.ObjectMeta,
		Spec:       instance.Spec.IngressSpec,
	}
}
