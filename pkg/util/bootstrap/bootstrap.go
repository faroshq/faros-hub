package bootstrap

// based on: https://github.com/kcp-dev/kcp/blob/715502c0274688b02e7310c4c32114a1c0c3a5c4/config/helpers/bootstrap.go

import (
	"bufio"
	"bytes"
	"context"
	"embed"
	"errors"
	"fmt"
	"io"
	"text/template"
	"time"

	tenancyhelper "github.com/kcp-dev/kcp/pkg/apis/tenancy/v1alpha1/helper"
	extensionsapiserver "k8s.io/apiextensions-apiserver/pkg/apiserver"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	apimachineryerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/wait"
	kubeyaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/discovery"
	memory "k8s.io/client-go/discovery/cached"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
	"k8s.io/klog"
)

// Bootstrap creates resources in a package's fs by
// continuously retrying the list. This is blocking, i.e. it only returns (with error)
// when the context is closed or with nil when the bootstrapping is successfully completed.
func Bootstrap(ctx context.Context, discoveryClient discovery.DiscoveryInterface, dynamicClient dynamic.Interface, fs embed.FS, opts ...Option) error {
	cache := memory.NewMemCacheClient(discoveryClient)
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(cache)

	// bootstrap non-crd resources
	var transformers []TransformFileFunc
	for _, opt := range opts {
		transformers = append(transformers, opt.TransformFile)
	}

	return wait.PollImmediateInfiniteWithContext(ctx, time.Second, func(ctx context.Context) (bool, error) {
		if err := CreateResourcesFromFS(ctx, dynamicClient, mapper, fs, transformers...); err != nil {
			klog.Infof("Failed to bootstrap resources, retrying: %v", err)
			// invalidate cache if resources not found
			// xref: https://github.com/kcp-dev/kcp/issues/655
			cache.Invalidate()
			return false, nil
		}
		return true, nil
	})
}

// CreateResourcesFromFS creates all resources from a filesystem.
func CreateResourcesFromFS(ctx context.Context, client dynamic.Interface, mapper meta.RESTMapper, fs embed.FS, transformers ...TransformFileFunc) error {
	files, err := fs.ReadDir(".")
	if err != nil {
		return err
	}

	var errs []error
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if err := CreateResourceFromFS(ctx, client, mapper, f.Name(), fs, transformers...); err != nil {
			errs = append(errs, err)
		}
	}
	return apimachineryerrors.NewAggregate(errs)
}

// CreateResourceFromFS creates given resource file.
func CreateResourceFromFS(ctx context.Context, client dynamic.Interface, mapper meta.RESTMapper, filename string, fs embed.FS, transformers ...TransformFileFunc) error {
	raw, err := fs.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("could not read %s: %w", filename, err)
	}

	if len(raw) == 0 {
		return nil // ignore empty files
	}

	d := kubeyaml.NewYAMLReader(bufio.NewReader(bytes.NewReader(raw)))
	var errs []error
	for i := 1; ; i++ {
		doc, err := d.Read()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return err
		}
		if len(bytes.TrimSpace(doc)) == 0 {
			continue
		}

		for _, transformer := range transformers {
			doc, err = transformer(doc)
			if err != nil {
				return err
			}
		}

		if err := createResourceFromFS(ctx, client, mapper, doc); err != nil {
			errs = append(errs, fmt.Errorf("failed to create resource %s doc %d: %w", filename, i, err))
		}
	}
	return apimachineryerrors.NewAggregate(errs)
}

func createResourceFromFS(ctx context.Context, client dynamic.Interface, mapper meta.RESTMapper, raw []byte) error {
	type Input struct {
		Batteries map[string]bool
	}
	input := Input{
		Batteries: map[string]bool{},
	}
	tmpl, err := template.New("manifest").Parse(string(raw))
	if err != nil {
		return fmt.Errorf("failed to parse manifest: %w", err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, input); err != nil {
		return fmt.Errorf("failed to execute manifest: %w", err)
	}

	obj, gvk, err := extensionsapiserver.Codecs.UniversalDeserializer().Decode(buf.Bytes(), nil, &unstructured.Unstructured{})
	if err != nil {
		return fmt.Errorf("could not decode raw: %w", err)
	}
	u, ok := obj.(*unstructured.Unstructured)
	if !ok {
		return fmt.Errorf("decoded into incorrect type, got %T, wanted %T", obj, &unstructured.Unstructured{})
	}

	m, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return fmt.Errorf("could not get REST mapping for %s: %w", gvk, err)
	}

	upserted, err := client.Resource(m.Resource).Namespace(u.GetNamespace()).Create(ctx, u, metav1.CreateOptions{})
	if err != nil {
		if apierrors.IsAlreadyExists(err) {
			existing, err := client.Resource(m.Resource).Namespace(u.GetNamespace()).Get(ctx, u.GetName(), metav1.GetOptions{})
			if err != nil {
				return err
			}

			u.SetResourceVersion(existing.GetResourceVersion())
			if _, err = client.Resource(m.Resource).Namespace(u.GetNamespace()).Update(ctx, u, metav1.UpdateOptions{}); err != nil {
				return fmt.Errorf("could not update %s %s: %w", gvk.Kind, tenancyhelper.QualifiedObjectName(existing), err)
			} else {
				klog.Infof("Updated %s %s", gvk, tenancyhelper.QualifiedObjectName(existing))
				return nil
			}
		}
		return err
	}

	klog.Infof("Bootstrapped %s %s", gvk.Kind, tenancyhelper.QualifiedObjectName(upserted))

	return nil
}

// TransformFileFunc transforms a resource file before being applied to the cluster.
type TransformFileFunc func(bs []byte) ([]byte, error)

// Option allows to customize the bootstrap process.
type Option struct {
	// TransformFileFunc is a function that transforms a resource file before being applied to the cluster.
	TransformFile TransformFileFunc
}

// ReplaceOption allows to customize the bootstrap process.
func ReplaceOption(pairs ...string) Option {
	return Option{
		TransformFile: func(bs []byte) ([]byte, error) {
			if len(pairs)%2 != 0 {
				return nil, fmt.Errorf("odd number of arguments: %v", pairs)
			}
			for i := 0; i < len(pairs); i += 2 {
				bs = bytes.ReplaceAll(bs, []byte(pairs[i]), []byte(pairs[i+1]))
			}
			return bs, nil
		},
	}
}
