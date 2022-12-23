# Bridge

Bridge is minimal controller like component that listens on SQL database for events
on service API and translates them to K8S objects. Objects itself are  created by
API server. Later k8s layer management is handled by workspaces, users controllers.

TODO: Add channel to bridge back status to sql database for from facing user API.
