/*
Copyright 2026.

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

package controller

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	appsv1alpha1 "github.com/Dylanadg/hello-operator/api/v1alpha1"
)

// HelloAppReconciler reconciles a HelloApp object
type HelloAppReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=apps.intern.dev,resources=helloapps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps.intern.dev,resources=helloapps/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps.intern.dev,resources=helloapps/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the HelloApp object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *HelloAppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = logf.FromContext(ctx)

	helloApp := &appsv1alpha1.HelloApp{}
	if err := r.Get(ctx, req.NamespacedName, helloApp); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	existing := &appsv1.Deployment{}
	err := r.Get(ctx, client.ObjectKey{Name: helloApp.Name, Namespace: helloApp.Namespace}, existing)

	if errors.IsNotFound(err) {
		dep := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      helloApp.Name,
				Namespace: helloApp.Namespace,
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: &helloApp.Spec.Replicas,
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"app": helloApp.Name},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{"app": helloApp.Name},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:    "busybox",
								Image:   "busybox:1.36",
								Command: []string{"sh", "-c", "echo $HELLO_MSG && sleep 3600"},
								Env: []corev1.EnvVar{
									{Name: "HELLO_MSG", Value: helloApp.Spec.Message},
								},
							},
						},
					},
				},
			},
		}
		ctrl.SetControllerReference(helloApp, dep, r.Scheme)
		return ctrl.Result{}, r.Create(ctx, dep)
	}

	if err != nil {
		return ctrl.Result{}, err
	}

	existing.Spec.Replicas = &helloApp.Spec.Replicas
	existing.Spec.Template.Spec.Containers[0].Env = []corev1.EnvVar{
		{Name: "HELLO_MSG", Value: helloApp.Spec.Message},
	}

	if err := r.Update(ctx, existing); err != nil {
		return ctrl.Result{}, err
	}

	helloApp.Status.AvailableReplicas = existing.Status.AvailableReplicas
	if err := r.Status().Update(ctx, helloApp); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *HelloAppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1alpha1.HelloApp{}).
		Owns(&appsv1.Deployment{}).
		Named("helloapp").
		Complete(r)
}
