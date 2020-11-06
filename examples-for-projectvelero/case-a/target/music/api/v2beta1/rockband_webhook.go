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

package v2beta1

import (
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var rockbandlog = logf.Log.WithName("rockband-resource-v2beta1")

func (r *RockBand) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-music-example-io-v2beta1-rockband,mutating=true,failurePolicy=fail,groups=music.example.io,resources=rockbands,verbs=create;update,versions=v2beta1,name=mrockband.kb.io

var _ webhook.Defaulter = &RockBand{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *RockBand) Default() {
	rockbandlog.Info("mutator v2beta1", "name", r.Name, "namespace", r.Namespace, "lead guitar", r.Spec.LeadGuitar)

	// TODO(user): fill in your defaulting logic.

	// LeadSinger is an optional field on RockBandv1
	// Adding "TBD" if it is empty
	if r.Spec.LeadSinger == "" {
		r.Spec.LeadSinger = "TBD from v2beta1 webhook"
		rockbandlog.Info("mutator v2beta1 object created without leadSinger, setting as TBD from v2beta1 webhook", "name", r.Name, "namespace", r.Namespace, "lead singer", r.Spec.LeadSinger)
	}

	// LeadGuitar is an optional field since RockBandv2beta1
	// Adding "TBD" if it is empty
	if r.Spec.LeadGuitar == "" {
		r.Spec.LeadGuitar = "TBD from v2beta1 webhook"
		rockbandlog.Info("mutator v2beta1 object created without leadGuitar, setting as TBD from v2beta1 webhook", "name", r.Name, "namespace", r.Namespace, "lead guitar", r.Spec.LeadGuitar)
	}

	// Silly mutation:
	// if the rockband name is beatles and leadSinger is John, set it as John Lennon
	if r.Name == "beatles" && r.Spec.LeadSinger == "John" {
		r.Spec.LeadSinger = "John Lennon"
	}
	// if the rockband name is ledzeppelin and leadGuitar is Jimmy, set it as Jimmy Page
	if r.Name == "ledzeppelin" && r.Spec.LeadGuitar == "Jimmy" {
		r.Spec.LeadGuitar = "Jimmy Page"
	}

	// if the rockband name is beatles and leadGuitar is George, set it as George Harrison
	if r.Name == "beatles" && r.Spec.LeadGuitar == "George" {
		r.Spec.LeadGuitar = "George Harrison"
	}

	rockbandlog.Info("mutator default final", "name", r.Name, "namespace", r.Namespace, "lead guitar", r.Spec.LeadGuitar)

}
