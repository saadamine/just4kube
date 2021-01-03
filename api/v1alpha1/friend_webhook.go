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

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var friendlog = logf.Log.WithName("friend-resource")

func (r *Friend) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-batch-v1alpha1-friend,mutating=true,failurePolicy=fail,groups=batch,resources=friends,verbs=create;update,versions=v1alpha1,name=mfriend.kb.io

var _ webhook.Defaulter = &Friend{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Friend) Default() {
	friendlog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// +kubebuilder:webhook:verbs=create;update,path=/validate-batch-v1alpha1-friend,mutating=false,failurePolicy=fail,groups=batch,resources=friends,versions=v1alpha1,name=vfriend.kb.io

var _ webhook.Validator = &Friend{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Friend) ValidateCreate() error {
	friendlog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return r.validateFriend()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Friend) ValidateUpdate(old runtime.Object) error {
	friendlog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return r.validateFriend()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Friend) ValidateDelete() error {
	friendlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}

func (r *Friend) validateFriend() error {
	var allErrs field.ErrorList
	if err := r.validateFriendUri(); err != nil {
		allErrs = append(allErrs, err)
	}

	if len(allErrs) == 0 {
		return nil
	}

	return apierrors.NewInvalid(
		schema.GroupKind{Group: "batch.tutorial.kubebuilder.io", Kind: "Friend"},
		r.Name, allErrs)
}


func (r *Friend) validateFriendUri() *field.Error {
	if r.ObjectMeta.Name != "sample-friend" {
		return field.Invalid(field.NewPath("metadata").Child("name"), r.Name, "must be named 'sample-friend'")
	}
	return nil
}
