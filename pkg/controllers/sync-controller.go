package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"github.com/openshift-splat-team/ci-secret-sync/data"
	"github.com/openshift-splat-team/ci-secret-sync/pkg/schema"
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
	interval := 300 * time.Second
	if s.Config.Sync.RefreshPeriodSeconds != 0 {
		interval = time.Duration(s.Config.Sync.RefreshPeriodSeconds) * time.Second
	}

	s.Logger.Info(fmt.Sprintf("refresh period %d seconds", time.Duration(interval)/time.Second))

	// setup a periodic callback
	ticker := time.NewTicker(interval)
	done := make(chan bool)

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

		sourceSchema := action.Source.Schema
		var schemaImpl schema.SchemaInterface

		if len(sourceSchema) > 0 {
			switch data.TriggerType(sourceSchema) {
			case data.SYNC_SCHEMA_REGISTRY:
				schemaImpl = &schema.RegistryCredentials{
					Data:   sourceValue,
					Config: &action.Source,
				}
				s.Logger.Info("loading fields with the registry schema")

			default:
				schemaImpl = &schema.Generic{
					Data: sourceValue,
				}
				s.Logger.Info("loading fields with the generic schema")
			}
		} else {
			schemaImpl = &schema.Generic{
				Data: sourceValue,
			}
			s.Logger.Info("loading fields with the generic schema")
		}

		update := false
		// deal with non-deployment type targets first to ensure dependencies are updated before rolling out
		for _, target := range action.Targets {
			fieldIndex := target.SourceFieldIndex

			value, err := schemaImpl.GetField(fieldIndex)
			s.Logger.Info(fmt.Sprintf("field value: %s", value))
			if err != nil {
				return ctrl.Result{}, fmt.Errorf("unable to get field value: %v", err)
			}
			s.Logger.Info(fmt.Sprintf("checking to see if %s differs from targets", action.Source.Name))
			if target.Type == "secret" {
				_update, err := utils.UpdateSecretKey(s.Context, s.Logger, s.Client, &target, []byte(value))
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
			s.Logger.Info(fmt.Sprintf("update detected in %s and has been mirrored, will roll out pods", action.Source.Name))
			for _, target := range action.Targets {
				if target.Type == "daemonset" {
					err = utils.RolloutDaemonset(s.Context, s.Logger, s.Client, &target)
					if err != nil {
						s.Logger.Error(err, fmt.Sprintf("unable to roll out daemonset %s. %v", target.Name, err))
						return ctrl.Result{}, fmt.Errorf("unable to roll out daemonset %s. %v", target.Name, err)
					}
				} else if target.Type == "deployment" {
					err = utils.RolloutDeployment(s.Context, s.Logger, s.Client, &target)
					if err != nil {
						s.Logger.Error(err, fmt.Sprintf("unable to roll out deployment %s. %v", target.Name, err))
						return ctrl.Result{}, fmt.Errorf("unable to roll out deployment %s. %v", target.Name, err)
					}
				}

			}
		}
	}
	return ctrl.Result{}, nil
}
