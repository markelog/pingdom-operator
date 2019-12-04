package checks

import (
	"context"
	"strings"

	pingdomv1alpha1 "github.com/markelog/pingdom-operator/pkg/apis/pingdom/v1alpha1"
	"github.com/redhat-cop/operator-utils/pkg/util"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const (
	// Name is the name of the controller
	Name = "controller_checks"

	// Finalizer is the name of the controllers finalizer
	Finalizer = "cleanup.pingdom"
)

var log = logf.Log.WithName(Name)

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Checks Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileChecks{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("checks-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Checks
	err = c.Watch(&source.Kind{Type: &pingdomv1alpha1.Checks{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource Pods and requeue the owner Checks
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &pingdomv1alpha1.Checks{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileChecks implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileChecks{}

// ReconcileChecks reconciles a Checks object
type ReconcileChecks struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Checks object and makes changes based on the state read
func (r *ReconcileChecks) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues(
		"Request.Namespace", request.Namespace,
		"Request.Name", request.Name,
	)

	reqLogger.Info("Reconciling Checks")

	// Fetch the Checks instance
	instance := &pingdomv1alpha1.Checks{}

	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Remove the check
	if util.IsBeingDeleted(instance) {

		// Remove it, yeah, but only if we have an Finalizer
		if !util.HasFinalizer(instance, Finalizer) {
			return reconcile.Result{}, nil
		}

		reqLogger.Info("Remove HTTP check")
		if err := instance.DeleteHTTP(); err != nil {
			// "Invalid check identifier" error means we have deleted
			// the check with other means
			if !strings.Contains(err.Error(), "Invalid check identifier") {
				reqLogger.Error(err, "Couldn't remove checker")
				return reconcile.Result{}, nil
			}

			reqLogger.Info("Seems you have removed the check already?")
		}

		util.RemoveFinalizer(instance, Finalizer)

		err := r.client.Update(context.TODO(), instance)
		if err != nil {
			return reconcile.Result{}, err
		}
	} else {
		reqLogger.Info("Setup HTTP check")
		if err := instance.SetupHTTP(); err != nil {
			return reconcile.Result{}, err
		}

		err = r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Error(err, "Failed to update pingdom ID")
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}
