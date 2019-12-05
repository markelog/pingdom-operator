package checks

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/markelog/pingdom-operator/pkg/apis/pingdom/v1alpha1"
	checks1alpha1 "github.com/markelog/pingdom-operator/pkg/apis/pingdom/v1alpha1"
)

var (
	namespace = "pingdom"
)

func createRequest(check *checks1alpha1.Checks) (*ReconcileChecks, reconcile.Request) {
	objs := []runtime.Object{check}

	scheme.Scheme.AddKnownTypes(checks1alpha1.SchemeGroupVersion, check)

	client := fake.NewFakeClient(objs...)
	checks := &ReconcileChecks{client: client, scheme: scheme.Scheme}
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      Name,
			Namespace: namespace,
		},
	}

	return checks, req
}

var createAndUpdateStub = httptest.NewServer(
	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var detailedCheckJSON = `
{
   "check":{
      "id":123456,
      "name":"Stuff yo!"
   }
}
`
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, detailedCheckJSON)
	}),
)

var deleteStub = httptest.NewServer(
	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}),
)

func TestChecksControllerCreate(t *testing.T) {
	check := &checks1alpha1.Checks{
		ObjectMeta: metav1.ObjectMeta{
			Name:      Name,
			Namespace: namespace,
		},
		Spec: checks1alpha1.ChecksSpec{
			User:     "killa",
			Password: "gorilla",
			Key:      "in-the-house",
			BaseURL:  createAndUpdateStub.URL,

			HTTP: &checks1alpha1.HTTPCheck{
				Name:       "eclectica.io-check",
				Hostname:   "eclectica.sh",
				Resolution: 1,
			},
		},
	}

	checks, req := createRequest(check)
	_, err := checks.Reconcile(req)
	if err != nil {
		t.Fatalf("Reconcile: (%v)", err)
	}

	test := &v1alpha1.Checks{}
	err = checks.client.Get(context.TODO(), req.NamespacedName, test)
	if err != nil {
		t.Fatalf("Get check: (%v)", err)
	}

	if test.Status.ID != 123456 {
		t.Fatal("Pingdom ID is wrong")
	}
}

func TestChecksControllerUpdate(t *testing.T) {
	var (
		namespace = "pingdom"
	)

	check := &checks1alpha1.Checks{
		ObjectMeta: metav1.ObjectMeta{
			Name:      Name,
			Namespace: namespace,
		},
		Spec: checks1alpha1.ChecksSpec{
			User:     "killa",
			Password: "gorilla",
			Key:      "in-the-house",
			BaseURL:  createAndUpdateStub.URL,

			HTTP: &checks1alpha1.HTTPCheck{
				Name:       "eclectica.io-check",
				Hostname:   "eclectica.sh",
				Resolution: 1,
			},
		},
		Status: checks1alpha1.ChecksStatus{
			ID: 1,
		},
	}

	checks, req := createRequest(check)
	_, err := checks.Reconcile(req)
	if err != nil {
		t.Fatalf("Reconcile: (%v)", err)
	}

	test := &v1alpha1.Checks{}
	err = checks.client.Get(context.TODO(), req.NamespacedName, test)
	if err != nil {
		t.Fatalf("Get check: (%v)", err)
	}

	if test.Status.ID != 1 {
		t.Fatal("Pingdom ID is wrong")
	}
}
func TestChecksControllerDelete(t *testing.T) {
	var (
		namespace = "pingdom"
		now       = metav1.NewTime(time.Now())
	)

	check := &checks1alpha1.Checks{
		ObjectMeta: metav1.ObjectMeta{
			Name:              Name,
			Namespace:         namespace,
			DeletionTimestamp: &now,
			Finalizers:        []string{"cleanup.pingdom"},
		},
		Status: checks1alpha1.ChecksStatus{
			ID: 1,
		},
		Spec: checks1alpha1.ChecksSpec{
			User:     "killa",
			Password: "gorilla",
			Key:      "in-the-house",
			BaseURL:  deleteStub.URL,

			HTTP: &checks1alpha1.HTTPCheck{
				Name:       "eclectica.io-check",
				Hostname:   "eclectica.sh",
				Resolution: 1,
			},
		},
	}

	checks, req := createRequest(check)
	_, err := checks.Reconcile(req)
	if err != nil {
		t.Fatalf("Reconcile: (%v)", err)
	}

	test := &v1alpha1.Checks{}
	err = checks.client.Get(context.TODO(), req.NamespacedName, test)
	if err != nil {
		t.Fatalf("Get check: (%v)", err)
	}

	if test.Status.ID == 0 {
		t.Fatal("Pingdom ID is not removed")
	}
}

func TestChecksControllerDeleteWithoutFinalizer(t *testing.T) {
	var (
		namespace = "pingdom"
		now       = metav1.NewTime(time.Now())
	)

	check := &checks1alpha1.Checks{
		ObjectMeta: metav1.ObjectMeta{
			Name:              Name,
			Namespace:         namespace,
			DeletionTimestamp: &now,

			// Empty
			Finalizers: []string{},
		},
		Status: checks1alpha1.ChecksStatus{
			ID: 1,
		},
		Spec: checks1alpha1.ChecksSpec{
			User:     "killa",
			Password: "gorilla",
			Key:      "in-the-house",
			BaseURL:  deleteStub.URL,

			HTTP: &checks1alpha1.HTTPCheck{
				Name:       "eclectica.io-check",
				Hostname:   "eclectica.sh",
				Resolution: 1,
			},
		},
	}

	checks, req := createRequest(check)
	_, err := checks.Reconcile(req)
	if err != nil {
		t.Fatalf("Reconcile: (%v)", err)
	}

	test := &v1alpha1.Checks{}
	err = checks.client.Get(context.TODO(), req.NamespacedName, test)
	if err != nil {
		t.Fatalf("Get check: (%v)", err)
	}

	if test.Status.ID != 1 {
		t.Fatal("Pingdom ID was removed")
	}
}
