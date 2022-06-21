/*

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
	"reflect"

	"github.com/go-logr/logr"
	predictorv1 "github.com/kserve/modelmesh-serving/apis/serving/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// const (
// 	AnnotationInjectOAuth = "Predictors.opendatahub.io/inject-oauth"
// )

// OpenshiftPredictorReconciler holds the controller configuration.
type OpenshiftPredictorReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

// ClusterRole permissions

// +kubebuilder:rbac:groups=serving.kserve.io,resources=predictors,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=services;serviceaccounts;secrets,verbs=get;list;watch;create;update;patch

// OAuthInjectionIsEnabled returns true if the oauth sidecar injection
// annotation is present in the Predictor
// func OAuthInjectionIsEnabled(Predictor nbv1.Predictor) bool {
// 	if Predictor.Annotations[AnnotationInjectOAuth] != "" {
// 		result, _ := strconv.ParseBool(Predictor.Annotations[AnnotationInjectOAuth])
// 		return result
// 	} else {
// 		return false
// 	}
// }

// ComparePredictors checks if two predictors are equal, if not return false
func ComparePredictors(pr1 predictorv1.Predictor, pr2 predictorv1.Predictor) bool {
	return reflect.DeepEqual(pr1.ObjectMeta.Labels, pr2.ObjectMeta.Labels) &&
		reflect.DeepEqual(pr1.ObjectMeta.Annotations, pr2.ObjectMeta.Annotations) &&
		reflect.DeepEqual(pr1.Spec, pr2.Spec)
}

// Reconcile performs the reconciling of the Openshift objects for a Kubeflow
// Predictor.
func (r *OpenshiftPredictorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// Initialize logger format
	log := r.Log.WithValues("Predictor", req.Name, "namespace", req.Namespace)

	// Get the Predictor object when a reconciliation event is triggered (create,
	// update, delete)
	Predictor := &predictorv1.Predictor{}
	err := r.Get(ctx, req.NamespacedName, Predictor)
	if err != nil && apierrs.IsNotFound(err) {
		log.Info("Stop Predictor reconciliation")
		return ctrl.Result{}, nil
	} else if err != nil {
		log.Error(err, "Unable to fetch the Predictor")
		return ctrl.Result{}, err
	}
	log.Info("Noticed a predictor")

	// // Create the objects required by the OAuth proxy sidecar (see
	// // Predictor_oauth.go file)
	// if OAuthInjectionIsEnabled(*Predictor) {
	// 	// Call the OAuth Service Account reconciler
	// 	err = r.ReconcileOAuthServiceAccount(Predictor, ctx)
	// 	if err != nil {
	// 		return ctrl.Result{}, err
	// 	}

	// 	// Call the OAuth Service reconciler
	// 	err = r.ReconcileOAuthService(Predictor, ctx)
	// 	if err != nil {
	// 		return ctrl.Result{}, err
	// 	}

	// 	// Call the OAuth Secret reconciler
	// 	err = r.ReconcileOAuthSecret(Predictor, ctx)
	// 	if err != nil {
	// 		return ctrl.Result{}, err
	// 	}

	// 	// Call the OAuth Route reconciler
	// 	err = r.ReconcileOAuthRoute(Predictor, ctx)
	// 	if err != nil {
	// 		return ctrl.Result{}, err
	// 	}
	// } else {
	// 	// Call the route reconciler (see Predictor_route.go file)
	// 	err = r.ReconcileRoute(Predictor, ctx)
	// 	if err != nil {
	// 		return ctrl.Result{}, err
	// 	}
	// }

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *OpenshiftPredictorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	builder := ctrl.NewControllerManagedBy(mgr).
		For(&predictorv1.Predictor{}).
		//		Owns(&routev1.Route{}).
		Owns(&corev1.ServiceAccount{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.Secret{})

	err := builder.Complete(r)
	if err != nil {
		return err
	}

	return nil
}