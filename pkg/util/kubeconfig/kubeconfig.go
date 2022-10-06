package kubeconfig

import (
	"encoding/json"

	clientcmdv1 "k8s.io/client-go/tools/clientcmd/api/v1"
)

func MakeKubeconfig(server, token, cacrt string) ([]byte, error) {
	return json.MarshalIndent(&clientcmdv1.Config{
		APIVersion: "v1",
		Kind:       "Config",
		Clusters: []clientcmdv1.NamedCluster{
			{
				Name: "cluster",
				Cluster: clientcmdv1.Cluster{
					Server:                   server,
					CertificateAuthorityData: []byte(cacrt),
				},
			},
		},
		AuthInfos: []clientcmdv1.NamedAuthInfo{
			{
				Name: "user",
				AuthInfo: clientcmdv1.AuthInfo{
					Token: token,
				},
			},
		},
		Contexts: []clientcmdv1.NamedContext{
			{
				Name: "cluster",
				Context: clientcmdv1.Context{
					Cluster:   "cluster",
					Namespace: "default",
					AuthInfo:  "user",
				},
			},
		},
		CurrentContext: "cluster",
	}, "", "    ")
}
