//+kubebuilder:object:generate=true
//+groupName=plugins.faros.sh
package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var GroupName = "plugins.faros.sh"

const (
	// PluginNetworkKind is the kind for a Network plugins
	PluginNetworkKind = "Network"
	// PluginAccessKind is the kind for an Access plugins
	PluginAccessKind = "Access"
	// PluginContainerRuntimeKind is the kind for a ContainerRuntime plugins
	PluginContainerRuntimeKind = "ContainerRuntime"
	// PluginMonitoringKind is the kind for a Monitoring plugins
	PluginMonitoringKind = "Monitoring"
	// PluginNotificationKind is the kind for a Notification plugins
	PluginNotificationKind = "Notification"
)

// SchemeGroupVersion is group version used to register these objects
var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1alpha1"}

// Kind takes an unqualified kind and returns back a Group qualified GroupKind
func Kind(kind string) schema.GroupKind {
	return SchemeGroupVersion.WithKind(kind).GroupKind()
}

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

var (
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme   = SchemeBuilder.AddToScheme
)

// Adds the list of known types to Scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&Network{},
		&NetworkList{},
		&Monitoring{},
		&MonitoringList{},
		&Notification{},
		&NotificationList{},
		&ContainerRuntime{},
		&ContainerRuntimeList{},
		&Access{},
		&AccessList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
