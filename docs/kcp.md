
Project is based on [KCP](https://www.kcp.io/) allowing each tennat to be isolated
at the same time allowing centralized management.

Each tenant 'looks like' Kubernetes cluster. It has namespaces, events, rbac, etc.

# Why KCP?

Before we go into why, lets what is so great at not great about Kubernetes when
it comes to building SAAS platform.

## What is great about Kubernetes

Kubernetes ecosystem is amazing. Creating new APIS, integrations with external
services is as easy as it can get. Using CRDs and controllers you can build any
new API in minutes.
You get CLI, API, UI for free. You can use any of the existing tools to build
and interact with your platform.

People, who are familiar with Kubernetes, can easily understand your platform.
Other - can use abstractions.

You get rich RBAC, events, namespacing for free.

## What is not so great about Kubernetes

Kubernetes is not designed to be SAAS platform. It is designed to be a platform.
Namespaces are not fine grained enough. You have to limit users to Namespaces,
and CRD's itself are build for cluster scope.

To have fully SAAS version of Kubernetes you need to run APIServer per tenant.
Which is possible with projects like vCluster, but it is not easy to setup and
its very expensive to run.

When you need only APIServer and CRD's, you don't need Nodes, node controllers,
Deployments, Services, etc etc. All these things are not needed for SAAS platform
as such (with some exception based on what are you building).

## What is great about KCP

KCP is designed to be SAAS platform. It allows to expose Kubernetes-like-API
per tenant. So each tenant is `cluster-admin` and can do anything with its
resources. Bu at the same time you are still running single APIServer and
single etcd.

KCP is designed to be used with CRD's. It allows to single instance of CRD's
and "share" them between tenants. So you 'provider' tenant can share CRDs to
other tenants.

Same applies to controllers. You can have single instance of controller and
"share" it between tenants. So when tenant created CRD, controller will be
reconcile it, but tenants don't have to run their own controllers.

As there is no nodes, compute, storage, etc, you don't need to run any of the
controllers. All system is lighweight and easy to setup.

## Main benefits

- Single APIServer
- Single etcd
- Each new tenant/user is its own 'virtual' cluster with its own rbac and 'cluster-admin'
- Each tenant can have its own CRD's and controllers and users can write their own
controllers to interact with their own tenant
- You can have single operator to manage all tenants and all CRD's
- You can use all dynamic client/code generation tools to bootstrap your platform and features
- You can use all existing Kubernetes tools to build your platform

Overall its like Kubernetes, but with less complexity and more flexibility hence
"kubernetes-like". Less complexity, cheap to run, all the benefits of Kubernetes.
