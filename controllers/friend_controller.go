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
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"


	urlv1alpha1 "gytigyg.io/api/v1alpha1"
)

// FriendReconciler reconciles a Friend object
type FriendReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=url.gytigyg.io,resources=friends,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=url.gytigyg.io,resources=friends/status,verbs=get;update;patch

func (r *FriendReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("friend", req.NamespacedName)

	// your logic here
	fmt.Println(req.NamespacedName)
	var friend urlv1alpha1.Friend

	if err := r.Get(ctx, req.NamespacedName, &friend); err != nil {
		log.Error(err, "unable to fetch Store")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	fmt.Println(friend.Spec.Uri)
	if friend.Spec.Uri != "https://url.gytigyg.io/amine"{
		fmt.Println("URI doesn't match")
	}

	if friend.Spec.Uri == "https://url.gytigyg.io/amine" {
		labels := map[string]string{
			"app": req.NamespacedName.Name,
		}

		configMap := &v1.ConfigMap{
			ObjectMeta:  metav1.ObjectMeta {
				Name:      "gytigyg-"+friend.Name,
				Namespace: req.Namespace,
				Labels:    labels,
			},
			Data: map[string]string{
				"uri": string(friend.Spec.Uri),
			},
		}
		if err := r.Client.Create(ctx, configMap); err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *FriendReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&urlv1alpha1.Friend{}).
		Complete(r)
}