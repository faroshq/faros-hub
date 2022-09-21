package bootstrap

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/ghodss/yaml"
	apisv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/apis/v1alpha1"
	schedulingv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/scheduling/v1alpha1"
	pluginhelpers "github.com/kcp-dev/kcp/pkg/cliplugins/helpers"
	utilyaml "github.com/mjudeikis/kcp-example/pkg/util/yaml"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/kube-openapi/pkg/util/sets"
)

func (b *bootstrap) bootstrapSyncTargets(ctx context.Context, name, workspace string, destinationRestConfig *rest.Config, labels map[string]string) error {
	fmt.Printf("Bootstrapping compute %s in workspace %s \n", name, workspace)
	client, rest, err := b.getWorkspaceClient(ctx, b.kcpClient, b.rest, workspace)
	if err != nil {
		return err
	}

	configURL, currentClusterName, err := pluginhelpers.ParseClusterURL(rest.Host)
	if err != nil {
		return fmt.Errorf("current URL %q does not point to cluster workspace", rest.Host)
	}

	// Configure synctargets on kcp side

	token, syncerID, syncTargetUID, err := b.enableSyncerForWorkspace(ctx, currentClusterName, name, client, rest, labels)
	if err != nil {
		return err
	}

	// render templates for target cluster

	serverURL := configURL.Scheme + "://" + configURL.Host

	requiredResourcesToSync := sets.NewString("deployments.apps", "secrets", "configmaps", "serviceaccounts")
	resourcesToSync := sets.NewString(b.config.Server.ComputeSyncerResourcesToSync...).Union(requiredResourcesToSync).List()

	input := syncerTemplateInput{
		ServerURL:          serverURL,
		CAData:             base64.StdEncoding.EncodeToString(rest.CAData),
		Token:              token,
		KCPNamespace:       b.config.Server.ComputeSyncerNamespace,
		Namespace:          syncerID,
		LogicalCluster:     currentClusterName.String(),
		SyncTarget:         name,
		SyncTargetUID:      syncTargetUID,
		Image:              b.config.Server.ComputeSyncerImage,
		Replicas:           1,
		ResourcesToSync:    resourcesToSync,
		QPS:                20,
		Burst:              30,
		FeatureGatesString: b.config.Server.ComputeSyncerFeatureGates,
	}

	resources, err := renderSyncerResources(input, syncerID)
	if err != nil {
		return err
	}

	documents, err := utilyaml.YAMLDocuments(bytes.NewBuffer(resources))
	if err != nil {
		return err
	}

	// get destination cluster and write
	dynamicClients, err := dynamic.NewForConfig(destinationRestConfig)
	if err != nil {
		return err
	}

	for _, document := range documents {
		var resource unstructured.Unstructured

		err := yaml.Unmarshal([]byte(document), &resource)
		if err != nil {
			return err
		}
		fmt.Println("Creating resource", resource.GetKind(), resource.GetName())
		var isUpdate bool
		gvr := schema.GroupVersionResource{
			Group:    resource.GroupVersionKind().Group,
			Version:  resource.GroupVersionKind().Version,
			Resource: strings.ToLower(resource.GroupVersionKind().Kind + "s"),
		}
		current, err := dynamicClients.Resource(gvr).Namespace(resource.GetNamespace()).Get(ctx, resource.GetName(), metav1.GetOptions{})
		if err == nil {
			isUpdate = true
		}
		if err != nil && !apierrors.IsNotFound(err) {
			return err
		}
		if isUpdate {
			resource.SetResourceVersion(current.GetResourceVersion())
			_, err = dynamicClients.Resource(gvr).Namespace(resource.GetNamespace()).Update(ctx, &resource, metav1.UpdateOptions{})
		} else {
			_, err = dynamicClients.Resource(gvr).Namespace(resource.GetNamespace()).Create(ctx, &resource, metav1.CreateOptions{})
		}
		if err != nil {
			return err
		}

	}

	return nil

}

func (b *bootstrap) bootstrapLocations(ctx context.Context, name, workspace string, labels map[string]string) error {
	fmt.Printf("Bootstrapping location %s in workspace %s \n", name, workspace)
	client, rest, err := b.getWorkspaceClient(ctx, b.kcpClient, b.rest, workspace)
	if err != nil {
		return err
	}

	_, currentClusterName, err := pluginhelpers.ParseClusterURL(rest.Host)
	if err != nil {
		return fmt.Errorf("current URL %q does not point to cluster workspace", rest.Host)
	}

	// Configure synctargets on kcp side
	var isUpdate bool
	current, err := client.Cluster(currentClusterName).SchedulingV1alpha1().Locations().Get(ctx, name, metav1.GetOptions{})
	if err == nil {
		isUpdate = true
	}
	if err != nil && !apierrors.IsNotFound(err) {
		return err
	}

	location := &schedulingv1alpha1.Location{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Annotations: map[string]string{
				"kcp.dev/cluster": currentClusterName.String(),
				"name":            name,
			},
			Labels: labels,
		},
		Spec: schedulingv1alpha1.LocationSpec{
			InstanceSelector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Resource: schedulingv1alpha1.GroupVersionResource{
				Group:    "workload.kcp.dev",
				Resource: "synctargets",
				Version:  "v1alpha1",
			},
		},
	}

	if isUpdate {
		location.ResourceVersion = current.ResourceVersion
		_, err = client.Cluster(currentClusterName).SchedulingV1alpha1().Locations().Update(ctx, location, metav1.UpdateOptions{})
	} else {
		_, err = client.Cluster(currentClusterName).SchedulingV1alpha1().Locations().Create(ctx, location, metav1.CreateOptions{})
	}
	if err != nil {
		return err
	}

	return nil
}

// bootstrapPlacement will take location name from workspace and create placement in target workspace with given name.
func (b *bootstrap) bootstrapPlacements(ctx context.Context, fromWorkspace, toWorkspace, name string, labels map[string]string) error {
	fmt.Printf("Bootstrapping placement %s -> %s (%s) \n", fromWorkspace, toWorkspace, name)
	client, rest, err := b.getWorkspaceClient(ctx, b.kcpClient, b.rest, toWorkspace)
	if err != nil {
		return err
	}

	_, currentClusterName, err := pluginhelpers.ParseClusterURL(rest.Host)
	if err != nil {
		return fmt.Errorf("current URL %q does not point to cluster workspace", rest.Host)
	}

	var isUpdate bool
	current, err := client.Cluster(currentClusterName).SchedulingV1alpha1().Placements().Get(ctx, name, metav1.GetOptions{})
	if err == nil {
		isUpdate = true
	}
	if err != nil && !apierrors.IsNotFound(err) {
		return err
	}

	placement := &schedulingv1alpha1.Placement{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Annotations: map[string]string{
				"name": name,
			},
			Labels: labels,
		},
		Spec: schedulingv1alpha1.PlacementSpec{
			LocationSelectors: []metav1.LabelSelector{
				{
					MatchLabels: labels,
				},
			},
			LocationWorkspace: fromWorkspace,
			LocationResource: schedulingv1alpha1.GroupVersionResource{
				Group:    "workload.kcp.dev",
				Resource: "synctargets",
				Version:  "v1alpha1",
			},
			NamespaceSelector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
		},
	}

	if isUpdate {
		placement.ResourceVersion = current.ResourceVersion
		_, err = client.Cluster(currentClusterName).SchedulingV1alpha1().Placements().Update(ctx, placement, metav1.UpdateOptions{})
	} else {
		_, err = client.Cluster(currentClusterName).SchedulingV1alpha1().Placements().Create(ctx, placement, metav1.CreateOptions{})
	}
	if err != nil {
		return err
	}

	return nil
}

// bootstrapPlacement will take location name from workspace and create placement in target workspace with given name.
func (b *bootstrap) bootstrapBinding(ctx context.Context, fromWorkspace, fromExport, toWorkspace, name string) error {
	fmt.Printf("Bootstrapping APIBinding %s -> %s (%s) \n", fromWorkspace, toWorkspace, name)
	client, rest, err := b.getWorkspaceClient(ctx, b.kcpClient, b.rest, toWorkspace)
	if err != nil {
		return err
	}

	_, currentClusterName, err := pluginhelpers.ParseClusterURL(rest.Host)
	if err != nil {
		return fmt.Errorf("current URL %q does not point to cluster workspace", rest.Host)
	}

	var isUpdate bool
	current, err := client.Cluster(currentClusterName).ApisV1alpha1().APIBindings().Get(ctx, name, metav1.GetOptions{})
	if err == nil {
		isUpdate = true
	}
	if err != nil && !apierrors.IsNotFound(err) {
		return err
	}

	binding := &apisv1alpha1.APIBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Annotations: map[string]string{
				"name": name,
			},
		},
		Spec: apisv1alpha1.APIBindingSpec{
			Reference: apisv1alpha1.ExportReference{
				Workspace: &apisv1alpha1.WorkspaceExportReference{
					Path:       fromWorkspace,
					ExportName: fromExport,
				},
			},
		},
	}

	if isUpdate {
		binding.ResourceVersion = current.ResourceVersion
		_, err = client.Cluster(currentClusterName).ApisV1alpha1().APIBindings().Update(ctx, binding, metav1.UpdateOptions{})
	} else {
		_, err = client.Cluster(currentClusterName).ApisV1alpha1().APIBindings().Create(ctx, binding, metav1.CreateOptions{})
	}
	if err != nil {
		return err
	}

	return nil
}
