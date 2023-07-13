package puller

import (
	"context"
	"encoding/json"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/retry"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	pullerv1alpha1 "github.com/puller-io/puller/pkg/apis/puller/v1alpha1"
)

const (
	ControllerName = "puller-controller"
	SecretName     = "puller-config"
	SecretLabelKey = "puller.io/name"
)

type Controller struct {
	client.Client
	KubeClient    kubernetes.Interface
	EventRecorder record.EventRecorder
}

// Reconcile performs a full reconciliation for the object referred to by the Request.
func (c *Controller) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.V(4).Info("Reconciling puller", "name", req.NamespacedName.Name)

	puller := pullerv1alpha1.Puller{}
	err := c.Client.Get(ctx, req.NamespacedName, &puller)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{Requeue: true}, err
	}

	var nsList *corev1.NamespaceList
	if puller.Spec.NamespaceAffinity == nil {
		nsList, err = c.KubeClient.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
		if err != nil {
			logger.Error(err, "failed to list all namespace")
			return ctrl.Result{Requeue: true}, err
		}
	} else {
		nsList, err = c.KubeClient.CoreV1().Namespaces().List(ctx, metav1.ListOptions{
			LabelSelector: puller.Spec.NamespaceAffinity.String(),
		})
		if err != nil {
			logger.Error(err, "failed to list namespace")
			return ctrl.Result{Requeue: true}, err
		}
	}

	var errs []error
	for _, ns := range nsList.Items {
		secret, err := newDockerSecret(SecretName, puller.Name, puller.Spec.Registries)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		secret.SetNamespace(ns.Name)
		if err := c.ensureSecret(ctx, secret); err != nil {
			errs = append(errs, err)
			continue
		}
		if err := c.ensurerServiceAccount(ctx, ns.Name, SecretName); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) != 0 {
		return ctrl.Result{Requeue: true}, utilerrors.NewAggregate(errs)
	}

	return ctrl.Result{}, nil
}

func (c *Controller) ensureSecret(ctx context.Context, secret *corev1.Secret) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		got, err := c.KubeClient.CoreV1().Secrets(secret.Namespace).Get(ctx, secret.Name, metav1.GetOptions{})
		if err != nil && apierrors.IsNotFound(err) {
			_, err := c.KubeClient.CoreV1().Secrets(secret.Namespace).Create(ctx, secret, metav1.CreateOptions{})
			if err != nil {
				return err
			}
			return nil
		} else if err != nil {
			return err
		}
		secret.SetResourceVersion(got.GetResourceVersion())
		_, err = c.KubeClient.CoreV1().Secrets(got.Namespace).Update(ctx, secret, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
		return nil
	})
}

func (c *Controller) ensurerServiceAccount(ctx context.Context, namespace string, name string) error {
	saList, err := c.KubeClient.CoreV1().ServiceAccounts(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	var errs []error
	for _, sa := range saList.Items {
		exists := sets.Set[string]{}
		for _, im := range sa.ImagePullSecrets {
			exists.Insert(im.Name)
		}
		if exists.Has(name) {
			continue
		}
		sa.ImagePullSecrets = append(sa.ImagePullSecrets, corev1.LocalObjectReference{
			Name: name,
		})

		err = retry.RetryOnConflict(retry.DefaultRetry, func() error {
			got, err := c.KubeClient.CoreV1().ServiceAccounts(sa.Namespace).Get(ctx, sa.Name, metav1.GetOptions{})
			if err != nil && apierrors.IsNotFound(err) {
				return nil
			} else if err != nil {
				return err
			}
			sa.SetResourceVersion(got.ResourceVersion)
			_, err = c.KubeClient.CoreV1().ServiceAccounts(sa.Namespace).Update(ctx, &sa, metav1.UpdateOptions{})
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			errs = append(errs, err)
		}
	}
	return utilerrors.NewAggregate(errs)
}

func newDockerSecret(name, pullerName string, registries []pullerv1alpha1.Registry) (*corev1.Secret, error) {
	content, err := buildDockerConfigJSON(registries)
	if err != nil {
		return nil, err
	}
	return &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: corev1.SchemeGroupVersion.String(),
			Kind:       "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				SecretLabelKey: pullerName,
			},
		},
		Type: corev1.SecretTypeDockerConfigJson,
		Data: map[string][]byte{
			corev1.DockerConfigJsonKey: content,
		},
	}, nil
}

func buildDockerConfigJSON(registries []pullerv1alpha1.Registry) ([]byte, error) {
	data := make(map[string]pullerv1alpha1.Registry)
	for _, r := range registries {
		data[r.Server] = pullerv1alpha1.Registry{
			Username: r.Username,
			Password: r.Password,
			Email:    r.Email,
			Auth:     r.Auth,
		}
	}
	content, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return content, nil
}

// SetupWithManager sets up the controller with the Manager.
func (c *Controller) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).For(&pullerv1alpha1.Puller{}).
		Watches(&corev1.Namespace{}, &handler.Funcs{
			CreateFunc: func(ctx context.Context, createEvent event.CreateEvent, limitingInterface workqueue.RateLimitingInterface) {
				c.namespaceWatcherFunc(ctx, createEvent.Object, limitingInterface)
			},
		}).
		Watches(&corev1.Namespace{}, &handler.Funcs{
			DeleteFunc: func(ctx context.Context, deleteEvent event.DeleteEvent, limitingInterface workqueue.RateLimitingInterface) {
				c.secretWatcherFunc(ctx, deleteEvent.Object, limitingInterface)
			},
		}).
		Complete(c)
}

func (c *Controller) namespaceWatcherFunc(ctx context.Context, obj client.Object, limitingInterface workqueue.RateLimitingInterface) {
	pullerList := pullerv1alpha1.PullerList{}
	if err := c.Client.List(ctx, &pullerList); err != nil {
		return
	}
	for _, puller := range pullerList.Items {
		limitingInterface.Add(reconcile.Request{NamespacedName: types.NamespacedName{
			Name:      puller.GetName(),
			Namespace: puller.GetNamespace(),
		}})
	}
}

func (c *Controller) secretWatcherFunc(ctx context.Context, obj client.Object, limitingInterface workqueue.RateLimitingInterface) {
	val, ok := obj.GetLabels()[SecretLabelKey]
	if !ok {
		return
	}
	puller := pullerv1alpha1.Puller{}
	if err := c.Client.Get(ctx, types.NamespacedName{Namespace: obj.GetNamespace(), Name: val}, &puller); err != nil {
		return
	}
	limitingInterface.Add(reconcile.Request{NamespacedName: types.NamespacedName{
		Name:      puller.GetName(),
		Namespace: puller.GetNamespace(),
	}})
}
