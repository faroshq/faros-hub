/*
Copyright 2022.

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
// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/faroshq/faros-hub/pkg/apis/plugins/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// NotificationLister helps list Notifications.
// All objects returned here must be treated as read-only.
type NotificationLister interface {
	// List lists all Notifications in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.Notification, err error)
	// Notifications returns an object that can list and get Notifications.
	Notifications(namespace string) NotificationNamespaceLister
	NotificationListerExpansion
}

// notificationLister implements the NotificationLister interface.
type notificationLister struct {
	indexer cache.Indexer
}

// NewNotificationLister returns a new NotificationLister.
func NewNotificationLister(indexer cache.Indexer) NotificationLister {
	return &notificationLister{indexer: indexer}
}

// List lists all Notifications in the indexer.
func (s *notificationLister) List(selector labels.Selector) (ret []*v1alpha1.Notification, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Notification))
	})
	return ret, err
}

// Notifications returns an object that can list and get Notifications.
func (s *notificationLister) Notifications(namespace string) NotificationNamespaceLister {
	return notificationNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// NotificationNamespaceLister helps list and get Notifications.
// All objects returned here must be treated as read-only.
type NotificationNamespaceLister interface {
	// List lists all Notifications in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.Notification, err error)
	// Get retrieves the Notification from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.Notification, error)
	NotificationNamespaceListerExpansion
}

// notificationNamespaceLister implements the NotificationNamespaceLister
// interface.
type notificationNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Notifications in the indexer for a given namespace.
func (s notificationNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.Notification, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Notification))
	})
	return ret, err
}

// Get retrieves the Notification from the indexer for a given namespace and name.
func (s notificationNamespaceLister) Get(name string) (*v1alpha1.Notification, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("notification"), name)
	}
	return obj.(*v1alpha1.Notification), nil
}
