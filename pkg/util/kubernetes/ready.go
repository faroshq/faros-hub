package kubernetes

import (
	"context"
	"time"

	kapiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	kapiextensions "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	kapiextensionsv1client "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/rest"
	apiregistrationv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"
	apiregistration "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset"
	apiregistrationv1client "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset/typed/apiregistration/v1"
)

// WaitForAllCRDsReady waits for all CRDs to be ready
func WaitForAllCRDReady(ctx context.Context, config *rest.Config) error {
	// TODO: handle connectivity retries (e.g. dial tcp 20.49.158.118:443: connect: connection refused)
	ae, err := kapiextensions.NewForConfig(config)
	if err != nil {
		return err
	}

	crds, err := ae.ApiextensionsV1().CustomResourceDefinitions().List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, crd := range crds.Items {
		err := waitForCRDReady(ctx, crd.Name, ae.ApiextensionsV1().CustomResourceDefinitions())
		if err != nil {
			return err
		}
	}

	return nil
}

// waitForCRDReady waits for a CustomResourceDefinition to be ready and registered
func waitForCRDReady(ctx context.Context, name string, cli kapiextensionsv1client.CustomResourceDefinitionInterface) error {
	if err := wait.PollImmediateInfinite(time.Second,
		checkCustomResourceDefinitionIsReady(ctx, cli, name),
	); err != nil {
		return err
	}

	return nil
}

// checkCustomResourceDefinitionIsReady returns a function which polls a
// CustomResourceDefinition and returns its readiness
func checkCustomResourceDefinitionIsReady(ctx context.Context, cli kapiextensionsv1client.CustomResourceDefinitionInterface, name string) func() (bool, error) {
	return func() (bool, error) {
		crd, err := cli.Get(ctx, name, metav1.GetOptions{})
		switch {
		case errors.IsNotFound(err):
			return false, nil
		case err != nil:
			return false, err
		}

		return customResourceDefinitionIsReady(crd), nil
	}
}

// customResourceDefinitionIsReady returns true if a CustomResourceDefinition is
// considered ready
func customResourceDefinitionIsReady(crd *kapiextensionsv1.CustomResourceDefinition) bool {
	for _, cond := range crd.Status.Conditions {
		if cond.Type == kapiextensionsv1.Established &&
			cond.Status == kapiextensionsv1.ConditionTrue {
			return true
		}
	}

	return false
}

// WaitForAllAPIServicesReady waits for all API services to be ready
func WaitForAllAPIServiceReady(ctx context.Context, config *rest.Config) error {
	ac, err := apiregistration.NewForConfig(config)
	if err != nil {
		return err
	}

	apis, err := ac.ApiregistrationV1().APIServices().List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, api := range apis.Items {
		err := waitForAPIServiceReady(ctx, api.Name, ac.ApiregistrationV1().APIServices())
		if err != nil {
			return err
		}
	}
	return nil
}

func waitForAPIServiceReady(ctx context.Context, name string, cli apiregistrationv1client.APIServiceInterface) error {
	if err := wait.PollImmediateInfinite(time.Second,
		checkAPIServiceIsReady(ctx, cli, name),
	); err != nil {
		return err
	}

	return nil
}

// checkAPIServiceIsReady returns a function which polls an APIService and
// returns its readiness
func checkAPIServiceIsReady(ctx context.Context, cli apiregistrationv1client.APIServiceInterface, name string) func() (bool, error) {
	return func() (bool, error) {
		svc, err := cli.Get(ctx, name, metav1.GetOptions{})
		switch {
		case errors.IsNotFound(err):
			return false, nil
		case err != nil:
			return false, err
		}

		return apiServiceIsReady(svc), nil
	}
}

// apiServiceIsReady returns true if an APIService is considered ready
func apiServiceIsReady(svc *apiregistrationv1.APIService) bool {
	for _, cond := range svc.Status.Conditions {
		if cond.Type == apiregistrationv1.Available &&
			cond.Status == apiregistrationv1.ConditionTrue {
			return true
		}
	}

	return false
}
