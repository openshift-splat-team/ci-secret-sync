package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/openshift-splat-team/ci-secret-sync/pkg/controllers"
	"github.com/openshift-splat-team/ci-secret-sync/pkg/utils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func init() {
}
func StartControllers() {
	scheme := runtime.NewScheme()
	setupLog := ctrl.Log.WithName("setup")

	var metricsAddr string
	var namespace string
	var enableLeaderElection bool
	var probeAddr string
	var config utils.Config
	ctx := context.TODO()

	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&namespace, "namespace", "vsphere-infra-helpers", "The namespace where ")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         false,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	client := mgr.GetClient()
	err = corev1.AddToScheme(mgr.GetScheme())
	if err != nil {
		setupLog.Error(err, "unable to add corev1 to scheme")
		os.Exit(1)
	}
	err = appsv1.AddToScheme(mgr.GetScheme())
	if err != nil {
		setupLog.Error(err, "unable to add appsv1 to scheme")
		os.Exit(1)
	}

	sync := controllers.SyncController{
		Client:  client,
		Config:  config.Get(),
		Context: ctx,
		Logger:  ctrl.Log,
	}

	sync.SetupWithManager(mgr)

	// for {
	// 	time.Sleep(10 * time.Second)
	// }

	mgr.Start(ctx)
}

func main() {
	fmt.Printf("Starting application...")
	StartControllers()
}
