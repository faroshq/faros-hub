package bootstrap

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ghodss/yaml"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/kyaml/filesys"
)

const retryLimit = 5

func (b *bootstrap) deployComponents(ctx context.Context, workspace, dir string) error {
	fmt.Printf("Bootstrapping resources %s in workspace %s \n", dir, workspace)
	_, rest, err := b.clientFactory.GetWorkspaceKCPClient(ctx, workspace)
	if err != nil {
		return err
	}

	dynamicClients, err := dynamic.NewForConfig(rest)
	if err != nil {
		return err
	}

	fs := filesys.MakeFsOnDisk()
	kr := krusty.MakeKustomizer(krusty.MakeDefaultOptions())

	resMap, err := kr.Run(fs, dir)
	if err != nil {
		return fmt.Errorf("failed to run kustomize: %w", err)
	}

	for _, resource := range resMap.Resources() {
		retryCounter := 0
	retry:
		var isUpdate bool
		gvr := schema.GroupVersionResource{
			Group:    resource.GetGvk().Group,
			Version:  resource.GetGvk().Version,
			Resource: strings.ToLower(resource.GetKind() + "s"),
		}

		current, err := dynamicClients.Resource(gvr).Namespace(resource.GetNamespace()).Get(ctx, resource.GetName(), metav1.GetOptions{})
		if err == nil {
			isUpdate = true
		}
		if err != nil && !apierrors.IsNotFound(err) {
			return err
		}

		var result unstructured.Unstructured
		data, err := resource.AsYAML()
		if err != nil {
			return err
		}

		err = yaml.Unmarshal(data, &result)
		if err != nil {
			return err
		}

		if isUpdate {
			fmt.Printf("Updating (retry %d) %s - %s \n", retryCounter, resource.GetGvk(), resource.GetName())
			result.SetResourceVersion(current.GetResourceVersion())
			_, err = dynamicClients.Resource(gvr).Namespace(resource.GetNamespace()).Update(ctx, &result, metav1.UpdateOptions{})
		} else {
			fmt.Printf("Creating (retry %d) %s - %s \n", retryCounter, resource.GetGvk(), resource.GetName())
			_, err = dynamicClients.Resource(gvr).Namespace(resource.GetNamespace()).Create(ctx, &result, metav1.CreateOptions{})
		}
		if err != nil {
			retryCounter++
			if retryCounter < retryLimit {
				time.Sleep(time.Second)
				goto retry
			}
			return err
		}
	}

	return nil
}
