/*
Copyright 2021 morningfish.

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
	"k8s.io/apimachinery/pkg/api/errors"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	belugav1 "github.com/morningfish/beluga/api/v1"
)

// BelugaReconciler reconciles a Beluga object
type BelugaReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=service.beluga.io,resources=belugas,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=service.beluga.io,resources=belugas/status,verbs=get;update;patch

func (r *BelugaReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("beluga", req.NamespacedName)
	instance := &belugav1.Beluga{}
	err := r.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			r.Log.Info("not exist")
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}
	err = InjectSidecar(instance)
	if err != nil {
		return ctrl.Result{}, err
	}
	err = InjectHost(instance)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *BelugaReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&belugav1.Beluga{}).
		Complete(r)
}
