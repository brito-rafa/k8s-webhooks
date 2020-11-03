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
	"strconv"
	"time"

	"github.com/go-logr/logr"
	k8serr "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	musicv1 "music/api/v1"
)

// RockBandReconciler reconciles a RockBand object
type RockBandReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=music.example.io,resources=rockbands,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=music.example.io,resources=rockbands/status,verbs=get;update;patch

func (r *RockBandReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("rockband", req.NamespacedName)

	// your logic here

	rb := &musicv1.RockBand{}
	if err := r.Client.Get(ctx, req.NamespacedName, rb); err != nil {
		// add some debug information if it's not a NotFound error
		if !k8serr.IsNotFound(err) {
			log.Error(err, "unable to fetch RockBand")
		}
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	msg := fmt.Sprintf("received reconcile request for %q (namespace: %q)", rb.GetName(), rb.GetNamespace())
	log.Info(msg)

	if rb.Status.LastPlayed == "" {
		year := time.Now().Year()
		// Adding the year in Status filed for now
		rb.Status.LastPlayed = strconv.Itoa(year)
		if err := r.Status().Update(ctx, rb); err != nil {
			log.Error(err, "unable to update RockBand status")
			return ctrl.Result{}, err
		}
	}
	return ctrl.Result{}, nil
}

func (r *RockBandReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&musicv1.RockBand{}).
		Complete(r)
}
