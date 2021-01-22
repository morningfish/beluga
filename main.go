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

package main

import (
	"flag"
	"github.com/morningfish/beluga/api/config"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"net/http"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	servicev1 "github.com/morningfish/beluga/api/v1"
	"github.com/morningfish/beluga/controllers"
	// +kubebuilder:scaffold:imports
)

var (
	scheme               = runtime.NewScheme()
	setupLog             = ctrl.Log.WithName("setup")
	metricsAddr          string
	enableLeaderElection bool
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)

	_ = servicev1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

func main() {
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.StringVar(&config.RuleFile, "rule-file-path", "/rule.rule", "rule file path")
	flag.StringVar(&config.BindHostFile, "bind-host-file", "", "bind host file path, separated by commas")
	flag.StringVar(&config.SubServerListenAddress, "server-listen-address", ":8080", "sub server listen address")
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))
	err := controllers.InitBindHost()
	if err != nil {
		setupLog.Error(err, "get bind host error")
		os.Exit(1)
	}
	StartSubServer()
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		Port:               9443,
		LeaderElection:     enableLeaderElection,
		LeaderElectionID:   "9c6273fb.beluga.io",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}
	controllers.BelugaReconcile = &controllers.BelugaReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("Beluga"),
		Scheme: mgr.GetScheme(),
	}
	if err = controllers.BelugaReconcile.SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Beluga")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

func StartSubServer() {
	// 监听订阅地址
	_ = os.Mkdir(config.SubFilePath, 0755)
	go http.ListenAndServe(config.SubServerListenAddress, http.FileServer(http.Dir(config.SubFilePath)))
}
