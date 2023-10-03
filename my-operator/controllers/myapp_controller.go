/*
Copyright 2023.

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
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appv1alpha1 "github.com/example/my-operator/api/v1alpha1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MyAppReconciler reconciles a MyApp object
type MyAppReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=app.example.com,resources=myapps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=app.example.com,resources=myapps/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=app.example.com,resources=myapps/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the MyApp object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *MyAppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	log.Info("Reconcile called")

	app := &appv1alpha1.MyApp{}

	err := r.Get(ctx, req.NamespacedName, app)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Myapp resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed")
		return ctrl.Result{}, err
	}

	log.Info("MyApp details", "Spec", app.Spec)

	// check namespace exist
	namespaceName := app.Spec.Namespace
	log.Info("namespace name : " + namespaceName)
	namespaceFound := &corev1.Namespace{}
	err = r.Client.Get(context.TODO(), types.NamespacedName{
		Name: namespaceName,
	}, namespaceFound)
	log.Info("namespace check completed")

	// Create the namespace if not exist
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating a new namespace", "name", namespaceName)
		err = r.Client.Create(context.TODO(), &corev1.Namespace{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Namespace",
				APIVersion: corev1.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: namespaceName,
			},
		})

		if err != nil {
			// namespace creation failed
			log.Error(err, "Failed to create new namespace", "name", namespaceName)
			return ctrl.Result{}, err
		}
		log.Info("Sucessfully created Namespace")
	}

	//  check cronjob exist
	cronjobName := app.Spec.CronJobConfig.Name
	cronjobFound := &batchv1.CronJob{}
	err = r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      cronjobName,
		Namespace: app.Spec.Namespace,
	}, cronjobFound)
	log.Info("Cronjob check completed", "name", cronjobName, "namespace", app.Spec.Namespace)

	// Create the cronjob if not exist
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating a new cronjob", "name", cronjobName, "image", app.Spec.CronJobConfig.Image)
		err = r.Client.Create(context.TODO(), &batchv1.CronJob{
			TypeMeta: metav1.TypeMeta{
				APIVersion: batchv1.SchemeGroupVersion.String(),
				Kind:       "CronJob",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      cronjobName,
				Namespace: app.Spec.Namespace,
			},
			Spec: batchv1.CronJobSpec{
				Schedule:                app.Spec.CronJobConfig.Schedule,
				ConcurrencyPolicy:       "Forbid",
				StartingDeadlineSeconds: &[]int64{100}[0],
				JobTemplate: batchv1.JobTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Name: cronjobName,
					},
					Spec: batchv1.JobSpec{
						Template: corev1.PodTemplateSpec{
							Spec: corev1.PodSpec{
								RestartPolicy: corev1.RestartPolicyNever,
								Containers: []corev1.Container{
									{
										Name:            cronjobName,
										Image:           app.Spec.CronJobConfig.Image,
										ImagePullPolicy: corev1.PullPolicy(corev1.PullAlways),
										Command:         []string{"python"},
										Args:            []string{"-m", "src.driver"},
									},
								},
							},
						},
					},
				},
			},
		})

		if err != nil {
			// cronjob creation failed
			log.Error(err, "Failed to create Cronjob")
			return ctrl.Result{}, err
		}
		log.Info("Sucessfully created Cronjob")
	}

	log.Info("Reconcile logic executed")

	// Call Reconsiler logic every 30 seconds so it will maintain required state( Ex: If you delete namespace/cronjob accidently then in 30 min when Reconsilder called at that time it will create required resources again)
	return ctrl.Result{RequeueAfter: time.Duration(30 * time.Second)}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MyAppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appv1alpha1.MyApp{}).
		Complete(r)
}
