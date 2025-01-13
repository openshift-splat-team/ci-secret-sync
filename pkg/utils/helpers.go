package utils

import (
	"bytes"
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/openshift-splat-team/ci-secret-sync/data"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func UpdateSecretKey(
	ctx context.Context,
	log logr.Logger,
	client client.Client,
	target *data.SyncItemTarget,
	sourceValue []byte) (bool, error) {
	key := target.Key
	var targetValue []byte
	update := false
	var ok bool

	namespacedName := types.NamespacedName{
		Namespace: target.Namespace,
		Name:      target.Name,
	}

	targetSecret := &corev1.Secret{}
	err := client.Get(ctx, namespacedName, targetSecret)
	if err != nil {
		return false, fmt.Errorf("unable to get target secret %s", target.Name)
	}

	if targetValue, ok = targetSecret.Data[key]; !ok {
		log.Info(fmt.Sprintf("key %s doesn't exist in secret %s, will add it", target.Key, target.Name))
		targetSecret.Data[key] = sourceValue
		update = true
	} else {
		if !bytes.Equal(targetValue, sourceValue) {
			log.Info(fmt.Sprintf("key %s in %s is out of sync with source secret, will update", target.Key, target.Name))
			targetSecret.Data[key] = sourceValue
			update = true
		}
	}

	if _config.DryRun {
		log.Info("dry run, not updating. set --dry-run=true to apply changes")
	} else if update {
		err = client.Update(ctx, targetSecret)
		if err != nil {
			log.Error(err, fmt.Sprintf("unable to update key %s in %s", target.Key, target.Name))
			return false, fmt.Errorf("unable to update key %s in %s. %v", target.Key, target.Name, err)
		}
	}

	return update, nil
}

func RolloutDaemonset(ctx context.Context,
	log logr.Logger,
	k8sclient client.Client,
	target *data.SyncItemTarget) error {

	namespacedName := types.NamespacedName{
		Namespace: target.Namespace,
		Name:      target.Name,
	}

	if _config.DryRun {
		log.Info("dry run, not rolling out daemonset. set --dry-run=true to apply changes")
		return nil
	}

	if target.Type == "daemonset" {
		targetDaemonset := appsv1.DaemonSet{}
		err := k8sclient.Get(ctx, namespacedName, &targetDaemonset)
		if err != nil {
			log.Error(err, fmt.Sprintf("unable to get daemonset %s", target.Key))
			return fmt.Errorf("unable to get daemonset %s, %v", target.Key, err)
		}

		labels := targetDaemonset.Spec.Template.ObjectMeta.Labels
		if len(labels) == 0 {
			log.Info("no labels found. can't proceed")
			return fmt.Errorf("no labels found. can't proceed")
		}

		err = DeletePodsInList(ctx, log, k8sclient, labels, target.Namespace)
		if err != nil {
			log.Error(err, fmt.Sprintf("unable to delete pods associated with %s", target.Name))
			return fmt.Errorf("unable to delete pods associated with %s, %v", target.Name, err)
		}
	} else {
		return fmt.Errorf("target %s is not a daemonset", target.Name)
	}
	return nil
}

func RolloutDeployment(ctx context.Context,
	log logr.Logger,
	k8sclient client.Client,
	target *data.SyncItemTarget) error {

	namespacedName := types.NamespacedName{
		Namespace: target.Namespace,
		Name:      target.Name,
	}

	if _config.DryRun {
		log.Info("dry run, not rolling out deployment. set --dry-run=true to apply changes")
		return nil
	}

	if target.Type == "deployment" {
		targetDeployment := appsv1.Deployment{}
		err := k8sclient.Get(ctx, namespacedName, &targetDeployment)
		if err != nil {
			log.Error(err, fmt.Sprintf("unable to get deployment %s", target.Key))
			return fmt.Errorf("unable to get deployment %s, %v", target.Key, err)
		}

		labels := targetDeployment.Spec.Template.ObjectMeta.Labels
		if len(labels) == 0 {
			log.Info("no labels found. can't proceed")
			return fmt.Errorf("no labels found. can't proceed")
		}

		err = DeletePodsInList(ctx, log, k8sclient, labels, target.Namespace)
		if err != nil {
			log.Error(err, fmt.Sprintf("unable to delete pods associated with %s", target.Name))
			return fmt.Errorf("unable to delete pods associated with %s, %v", target.Name, err)
		}
	} else {
		return fmt.Errorf("target %s is not a deployment", target.Name)
	}
	return nil
}

func DeletePodsInList(
	ctx context.Context,
	log logr.Logger,
	k8sclient client.Client,
	labels map[string]string,
	namespace string,
) error {
	associatedPods := corev1.PodList{
		TypeMeta: metav1.TypeMeta{},
	}

	for k, v := range labels {
		err := k8sclient.List(ctx, &associatedPods, []client.ListOption{
			client.InNamespace(namespace),
			client.MatchingLabels{k: v},
		}...)
		if err != nil {
			log.Error(err, "unable to get pods")
			return fmt.Errorf("unable to get pods %v", err)
		}

		if associatedPods.Items != nil {
			break
		}
	}

	if associatedPods.Items == nil {
		log.Info("no pods found for. unable to restart")
		return fmt.Errorf("no pods found for. unable to restart")
	}

	for _, pod := range associatedPods.Items {
		log.Info(fmt.Sprintf("deleting pod %s", pod.Name))
		err := k8sclient.Delete(ctx, &pod)
		if err != nil {
			log.Error(err, fmt.Sprintf("unable to delete pod %s", pod.Name))
			return fmt.Errorf("unable to delete pod %s, %v", pod.Name, err)
		}
	}

	return nil
}
