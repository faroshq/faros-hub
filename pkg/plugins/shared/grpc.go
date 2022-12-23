package shared

import (
	"context"

	"github.com/faroshq/faros-hub/pkg/plugins"
	"github.com/faroshq/faros-hub/pkg/plugins/proto"
	plugin "github.com/hashicorp/go-plugin"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/yaml"
)

// GRPCClient is an implementation of Interface that talks over RPC.
type GRPCClient struct {
	broker *plugin.GRPCBroker
	client proto.PluginInterfaceClient
}

func (m *GRPCClient) GetName(ctx context.Context) (string, error) {
	resp, err := m.client.GetName(ctx, &proto.Empty{})
	if err != nil {
		return "", err
	}

	return resp.Name, nil
}

func (m *GRPCClient) GetVersion(ctx context.Context) (string, error) {
	resp, err := m.client.GetVersion(ctx, &proto.Empty{})
	if err != nil {
		return "", err
	}

	return resp.Version, nil
}

func (m *GRPCClient) GetAPIResourceSchema(ctx context.Context) ([]byte, error) {
	resp, err := m.client.GetAPIResourceSchema(ctx, &proto.Empty{})
	if err != nil {
		return nil, err
	}

	return resp.Schema, nil
}

func (m *GRPCClient) GetAPIExportSchema(ctx context.Context) ([]byte, error) {
	resp, err := m.client.GetAPIExportSchema(ctx, &proto.Empty{})
	if err != nil {
		return nil, err
	}

	return resp.Schema, nil
}

func (m *GRPCClient) Init(ctx context.Context, name, namespace string, config *rest.Config) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	_, err = m.client.Init(ctx, &proto.InitRequest{
		Name:       name,
		Namespace:  namespace,
		RestConfig: data,
	})
	if err != nil {
		return err
	}

	return nil
}

func (m *GRPCClient) Run(ctx context.Context) error {
	_, err := m.client.Run(ctx, &proto.Empty{})
	if err != nil {
		return err
	}
	return nil
}

func (m *GRPCClient) Stop(ctx context.Context) error {
	_, err := m.client.Stop(ctx, &proto.Empty{})
	if err != nil {
		return err
	}
	return nil
}

// Here is the gRPC server that GRPCClient talks to.
type GRPCServer struct {
	// This is the real implementation
	Impl plugins.Interface

	broker *plugin.GRPCBroker
}

func (m *GRPCServer) GetName(ctx context.Context, req *proto.Empty) (*proto.GetNameResponse, error) {
	name, err := m.Impl.GetName(ctx)
	return &proto.GetNameResponse{Name: name}, err
}

func (m *GRPCServer) GetVersion(ctx context.Context, req *proto.Empty) (*proto.GetVersionResponse, error) {
	version, err := m.Impl.GetVersion(ctx)
	return &proto.GetVersionResponse{Version: version}, err
}

func (m *GRPCServer) GetAPIResourceSchema(ctx context.Context, req *proto.Empty) (*proto.GetAPIResourceSchemaResponse, error) {
	data, err := m.Impl.GetAPIResourceSchema(ctx)
	return &proto.GetAPIResourceSchemaResponse{Schema: data}, err
}

func (m *GRPCServer) GetAPIExportSchema(ctx context.Context, req *proto.Empty) (*proto.GetAPIExportSchemaResponse, error) {
	data, err := m.Impl.GetAPIExportSchema(ctx)
	return &proto.GetAPIExportSchemaResponse{Schema: data}, err
}

func (m *GRPCServer) Init(ctx context.Context, req *proto.InitRequest) (*proto.Empty, error) {
	var config rest.Config
	err := yaml.Unmarshal(req.RestConfig, &config)
	if err != nil {
		return &proto.Empty{}, err
	}
	err = m.Impl.Init(ctx, req.Name, req.Namespace, &config)
	if err != nil {
		return &proto.Empty{}, err
	}
	return &proto.Empty{}, nil
}

func (m *GRPCServer) Run(ctx context.Context, req *proto.Empty) (*proto.Empty, error) {
	err := m.Impl.Run(ctx)
	if err != nil {
		return &proto.Empty{}, err
	}
	return &proto.Empty{}, nil
}

func (m *GRPCServer) Stop(ctx context.Context, req *proto.Empty) (*proto.Empty, error) {
	err := m.Impl.Stop(ctx)
	if err != nil {
		return &proto.Empty{}, err
	}
	return &proto.Empty{}, nil
}
