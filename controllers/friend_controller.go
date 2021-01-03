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
	"fmt"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	urlv1alpha1 "gytigyg.io/api/v1alpha1"
)

// FriendReconciler reconciles a Friend object
type FriendReconciler struct {
	client.Client
	Log        logr.Logger
	Scheme     *runtime.Scheme
	Properties FriendProperties
}
type FriendProperties struct {
	AgentVersion string
}

// +kubebuilder:rbac:groups=url.gytigyg.io,resources=friends,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=url.gytigyg.io,resources=friends/status,verbs=get;update;patch

func (f *FriendProperties) NewConfigMapForFriend(friend *urlv1alpha1.Friend) *corev1.ConfigMap {
	labels := map[string]string{
		"app": friend.Name,
	}

	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "gytigyg-" + friend.Name,
			Namespace: friend.Namespace,
			Labels:    labels,
		},
		Data: map[string]string{
			"uri": string(friend.Spec.Uri),
		},
	}
	return configMap
}

func (r *FriendReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("friend", req.NamespacedName)

	// your logic here
	var friend urlv1alpha1.Friend

	if err := r.Get(ctx, req.NamespacedName, &friend); err != nil {
		log.Error(err, "unable to fetch Friend")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// fmt.Println(friend.Spec.Uri)
	if friend.Spec.Uri != "https://url.gytigyg.io" {
		fmt.Println("URI doesn't match")
		// Configmap creation failed
		friend.Status.Active = "Failed"
		if err := r.Status().Update(ctx, &friend); err != nil {
			r.Log.Error(err, "unable to update friend status")
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	if friend.Spec.Uri == "https://url.gytigyg.io" {
		ctx := context.Background()
		configmapName := types.NamespacedName{Name: "gytigyg-" + friend.Name, Namespace: friend.Namespace}
		configmap := &corev1.ConfigMap{}

		err := r.Client.Get(ctx, configmapName, configmap)
		if err != nil && errors.IsNotFound(err) {

			r.Log.Info("Creating a new configmap", "Friend.Namespace", friend.Namespace, "friend.Name", "gytigyg-"+friend.Name)
			newConfigMap := r.Properties.NewConfigMapForFriend(&friend)
			if err := controllerutil.SetControllerReference(&friend, newConfigMap, r.Scheme); err != nil {
				return ctrl.Result{}, err
			}
			if err := r.Client.Create(ctx, newConfigMap); err != nil {
				return ctrl.Result{}, err
			}
		}
		if err != nil && !errors.IsNotFound(err) {
			friend.Status.Active = "Failed"
			if err := r.Status().Update(ctx, &friend); err != nil {
				r.Log.Error(err, "unable to update friend status")
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, err
		}

		// Configmap creation failed
		friend.Status.Active = "Success"
		if err := r.Status().Update(ctx, &friend); err != nil {
			r.Log.Error(err, "unable to update friend status")
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, nil
}

func (r *FriendReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&urlv1alpha1.Friend{}).
		Owns(&corev1.ConfigMap{}).
		Complete(r)
}
