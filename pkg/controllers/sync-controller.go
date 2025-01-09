package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"github.com/openshift-splat-team/ci-secret-sync/data"
	"github.com/openshift-splat-team/ci-secret-sync/pkg/utils"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type SyncController struct {
	Context context.Context
	Client  client.Client
	Config  *data.SyncConfig
	Logger  logr.Logger
}

func (s *SyncController) SetupWithManager(mgr ctrl.Manager) error {
	// Define the interval
	interval := 3 * time.Second

	// Create a new ticker
	ticker := time.NewTicker(interval)

	// Create a channel to signal the end of the program (optional)
	done := make(chan bool)

	// Start a goroutine that executes periodically
	go func() {
		for {
			select {
			case <-ticker.C:
				s.Reconcile()
			case <-done:
				return
			}
		}
	}()

	return nil
}

func (s *SyncController) Reconcile() (ctrl.Result, error) {
	for _, action := range s.Config.Sync.Actions {
		namespacedName := types.NamespacedName{
			Namespace: action.Source.Namespace,
			Name:      action.Source.Name,
		}

		sourceSecret := &corev1.Secret{}
		err := s.Client.Get(s.Context, namespacedName, sourceSecret)
		if err != nil {
			s.Logger.Error(err, fmt.Sprintf("unable to get source secret: %v. can't proceed", namespacedName))
			return ctrl.Result{}, fmt.Errorf("unable to get source secret: %v. can't proceed. %v", namespacedName, err)
		}

		var sourceValue []byte
		var ok bool
		if sourceValue, ok = sourceSecret.Data[action.Source.Key]; !ok {
			s.Logger.Error(err, fmt.Sprintf("unable to get key %s in source secret: %v. can't proceed", action.Source.Key, namespacedName))
			return ctrl.Result{}, fmt.Errorf("unable to get key %s in source secret: %v. can't proceed. %v", action.Source.Key, namespacedName, err)
		}

		update := false
		// deal with non-deployment type targets first to ensure dependencies are updated before rolling out
		for _, target := range action.Targets {
			if target.Type == "secret" {
				_update, err := utils.UpdateSecretKey(s.Context, s.Logger, s.Client, &target, sourceValue)
				if err != nil {
					s.Logger.Error(err, fmt.Sprintf("unable to update secret %s. %v", target.Name, err))
					return ctrl.Result{}, fmt.Errorf("unable to update secret %s. %v", target.Name, err)
				}
				if _update {
					update = true
				}
			}
		}

		if update {
			for _, target := range action.Targets {
				if target.Type == "daemonset" {
					err = utils.RolloutDaemonset(s.Context, s.Logger, s.Client, &target)
					if err != nil {
						s.Logger.Error(err, fmt.Sprintf("unable to roll out data set %s. %v", target.Name, err))
						return ctrl.Result{}, fmt.Errorf("unable to roll out data set %s. %v", target.Name, err)
					}
				}
			}
		}
	}
	return ctrl.Result{}, nil
}
