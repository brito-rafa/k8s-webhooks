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

package v1

import (
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var rockbandlog = logf.Log.WithName("rockband-resource")

func (r *RockBand) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-music-example-io-v1-rockband,mutating=true,failurePolicy=fail,groups=music.example.io,resources=rockbands,verbs=create;update,versions=v1,name=mrockband.kb.io

var _ webhook.Defaulter = &RockBand{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *RockBand) Default() {
	rockbandlog.Info("mutator default", "name", r.Name, "namespace", r.Namespace)

	// TODO(user): fill in your defaulting logic.

	// LeadSinger is an optional field on RockBandv1
	// Adding "TBD" if it is empty
	if r.Spec.LeadSinger == "" {
		r.Spec.LeadSinger = "TBD"
	}

	// Silly mutation:
	// if the rockband name is beatles and leadSinger is John, set it as John Lennon
	if r.Name == "beatles" && r.Spec.LeadSinger == "John" {
		r.Spec.LeadSinger = "John Lennon"
	}
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// +kubebuilder:webhook:verbs=create;update,path=/validate-music-example-io-v1-rockband,mutating=false,failurePolicy=fail,groups=music.example.io,resources=rockbands,versions=v1,name=vrockband.kb.io

var _ webhook.Validator = &RockBand{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *RockBand) ValidateCreate() error {
	rockbandlog.Info("validate create", "name", r.Name, "namespace", r.Namespace, "lead singer", r.Spec.LeadSinger)

	// TODO(user): fill in your validation logic upon object creation.

	var allErrs field.ErrorList

	// Just an example of validation: one cannot create rockbands under kube-system namespace
	if r.Namespace == "kube-system" {
		err := field.Invalid(field.NewPath("metadata").Child("namespace"), r.Namespace, "is forbidden to have rockbands.")
		allErrs = append(allErrs, err)
	}

	if len(allErrs) == 0 {
		return nil
	}

	return apierrors.NewInvalid(
		schema.GroupKind{Group: "music.example.io", Kind: "RockBand"},
		r.Name, allErrs)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *RockBand) ValidateUpdate(old runtime.Object) error {
	rockbandlog.Info("validate update", "name", r.Name, "namespace", r.Namespace, "lead singer", r.Spec.LeadSinger)

	// TODO(user): fill in your validation logic upon object update.

	var allErrs field.ErrorList

	// Disclaimer: The following condition will never be met because of the Default mutation
	if r.Name == "beatles" && r.Spec.LeadSinger == "John" {
		err := field.Invalid(field.NewPath("spec").Child("leadSinger"), r.Spec.LeadSinger, "has the shortname of the singer.")
		allErrs = append(allErrs, err)
	}

	// Silly validation
	if r.Name == "beatles" && r.Spec.LeadSinger == "Ringo" {
		err := field.Invalid(field.NewPath("spec").Child("leadSinger"), r.Spec.LeadSinger, "was the drummer. Suggest you to pick John or Paul.")
		allErrs = append(allErrs, err)
	}

	if len(allErrs) == 0 {
		return nil
	}

	return apierrors.NewInvalid(
		schema.GroupKind{Group: "music.example.io", Kind: "RockBand"},
		r.Name, allErrs)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *RockBand) ValidateDelete() error {
	rockbandlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}
