// Code generated by protoc-gen-go. DO NOT EDIT.
// source: plugin.proto

package proto

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type GetNameResponse struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetNameResponse) Reset()         { *m = GetNameResponse{} }
func (m *GetNameResponse) String() string { return proto.CompactTextString(m) }
func (*GetNameResponse) ProtoMessage()    {}
func (*GetNameResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_22a625af4bc1cc87, []int{0}
}

func (m *GetNameResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetNameResponse.Unmarshal(m, b)
}
func (m *GetNameResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetNameResponse.Marshal(b, m, deterministic)
}
func (m *GetNameResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetNameResponse.Merge(m, src)
}
func (m *GetNameResponse) XXX_Size() int {
	return xxx_messageInfo_GetNameResponse.Size(m)
}
func (m *GetNameResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetNameResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetNameResponse proto.InternalMessageInfo

func (m *GetNameResponse) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type GetAPIResourceSchemaResponse struct {
	Schema               []byte   `protobuf:"bytes,1,opt,name=schema,proto3" json:"schema,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetAPIResourceSchemaResponse) Reset()         { *m = GetAPIResourceSchemaResponse{} }
func (m *GetAPIResourceSchemaResponse) String() string { return proto.CompactTextString(m) }
func (*GetAPIResourceSchemaResponse) ProtoMessage()    {}
func (*GetAPIResourceSchemaResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_22a625af4bc1cc87, []int{1}
}

func (m *GetAPIResourceSchemaResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetAPIResourceSchemaResponse.Unmarshal(m, b)
}
func (m *GetAPIResourceSchemaResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetAPIResourceSchemaResponse.Marshal(b, m, deterministic)
}
func (m *GetAPIResourceSchemaResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetAPIResourceSchemaResponse.Merge(m, src)
}
func (m *GetAPIResourceSchemaResponse) XXX_Size() int {
	return xxx_messageInfo_GetAPIResourceSchemaResponse.Size(m)
}
func (m *GetAPIResourceSchemaResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetAPIResourceSchemaResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetAPIResourceSchemaResponse proto.InternalMessageInfo

func (m *GetAPIResourceSchemaResponse) GetSchema() []byte {
	if m != nil {
		return m.Schema
	}
	return nil
}

type InitRequest struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Namespace            string   `protobuf:"bytes,2,opt,name=namespace,proto3" json:"namespace,omitempty"`
	RestConfig           []byte   `protobuf:"bytes,3,opt,name=rest_config,json=restConfig,proto3" json:"rest_config,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InitRequest) Reset()         { *m = InitRequest{} }
func (m *InitRequest) String() string { return proto.CompactTextString(m) }
func (*InitRequest) ProtoMessage()    {}
func (*InitRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_22a625af4bc1cc87, []int{2}
}

func (m *InitRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InitRequest.Unmarshal(m, b)
}
func (m *InitRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InitRequest.Marshal(b, m, deterministic)
}
func (m *InitRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InitRequest.Merge(m, src)
}
func (m *InitRequest) XXX_Size() int {
	return xxx_messageInfo_InitRequest.Size(m)
}
func (m *InitRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_InitRequest.DiscardUnknown(m)
}

var xxx_messageInfo_InitRequest proto.InternalMessageInfo

func (m *InitRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *InitRequest) GetNamespace() string {
	if m != nil {
		return m.Namespace
	}
	return ""
}

func (m *InitRequest) GetRestConfig() []byte {
	if m != nil {
		return m.RestConfig
	}
	return nil
}

type Empty struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Empty) Reset()         { *m = Empty{} }
func (m *Empty) String() string { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()    {}
func (*Empty) Descriptor() ([]byte, []int) {
	return fileDescriptor_22a625af4bc1cc87, []int{3}
}

func (m *Empty) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Empty.Unmarshal(m, b)
}
func (m *Empty) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Empty.Marshal(b, m, deterministic)
}
func (m *Empty) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Empty.Merge(m, src)
}
func (m *Empty) XXX_Size() int {
	return xxx_messageInfo_Empty.Size(m)
}
func (m *Empty) XXX_DiscardUnknown() {
	xxx_messageInfo_Empty.DiscardUnknown(m)
}

var xxx_messageInfo_Empty proto.InternalMessageInfo

func init() {
	proto.RegisterType((*GetNameResponse)(nil), "proto.GetNameResponse")
	proto.RegisterType((*GetAPIResourceSchemaResponse)(nil), "proto.GetAPIResourceSchemaResponse")
	proto.RegisterType((*InitRequest)(nil), "proto.InitRequest")
	proto.RegisterType((*Empty)(nil), "proto.Empty")
}

func init() {
	proto.RegisterFile("plugin.proto", fileDescriptor_22a625af4bc1cc87)
}

var fileDescriptor_22a625af4bc1cc87 = []byte{
	// 263 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x90, 0x41, 0x4b, 0xc3, 0x40,
	0x10, 0x85, 0x89, 0x4d, 0x5b, 0x3a, 0x0d, 0x14, 0x06, 0x29, 0xa1, 0x14, 0xd4, 0x15, 0xa1, 0xa7,
	0x0a, 0x0a, 0xde, 0x45, 0xa4, 0xe4, 0x22, 0x65, 0xfb, 0x03, 0x34, 0x86, 0x69, 0x0d, 0x98, 0xdd,
	0x35, 0x3b, 0x39, 0xf8, 0xe7, 0x45, 0x32, 0x06, 0x6b, 0x43, 0xec, 0x69, 0x32, 0x2f, 0xef, 0xbd,
	0xe5, 0x1b, 0x88, 0xdc, 0x7b, 0xb5, 0xcb, 0xcd, 0xd2, 0x95, 0x96, 0x2d, 0xf6, 0x65, 0xa8, 0x2b,
	0x98, 0xac, 0x88, 0x9f, 0xd2, 0x82, 0x34, 0x79, 0x67, 0x8d, 0x27, 0x44, 0x08, 0x4d, 0x5a, 0x50,
	0x1c, 0x9c, 0x07, 0x8b, 0x91, 0x96, 0x6f, 0x75, 0x07, 0xf3, 0x15, 0xf1, 0xfd, 0x3a, 0xd1, 0xe4,
	0x6d, 0x55, 0x66, 0xb4, 0xc9, 0xde, 0xa8, 0x48, 0x7f, 0x33, 0x53, 0x18, 0x78, 0x51, 0x24, 0x15,
	0xe9, 0x66, 0x53, 0x2f, 0x30, 0x4e, 0x4c, 0xce, 0x9a, 0x3e, 0x2a, 0xf2, 0xdc, 0x55, 0x8d, 0x73,
	0x18, 0xd5, 0xd3, 0xbb, 0x34, 0xa3, 0xf8, 0x44, 0x7e, 0xec, 0x05, 0x3c, 0x83, 0x71, 0x49, 0x9e,
	0x9f, 0x33, 0x6b, 0xb6, 0xf9, 0x2e, 0xee, 0x49, 0x3b, 0xd4, 0xd2, 0x83, 0x28, 0x6a, 0x08, 0xfd,
	0xc7, 0xc2, 0xf1, 0xe7, 0xcd, 0x57, 0x00, 0x93, 0xb5, 0x10, 0x26, 0x86, 0xa9, 0xdc, 0xd6, 0xe9,
	0x6b, 0x18, 0x36, 0x74, 0x18, 0xfd, 0x70, 0x2f, 0xc5, 0x3c, 0x9b, 0x36, 0x5b, 0x9b, 0x3d, 0x81,
	0xd3, 0x2e, 0xce, 0x56, 0xfa, 0x72, 0x9f, 0xfe, 0xff, 0x24, 0x0b, 0x08, 0x6b, 0x74, 0xc4, 0xc6,
	0xfc, 0xe7, 0x0e, 0xb3, 0x83, 0x3a, 0xbc, 0x80, 0x9e, 0xae, 0x4c, 0xeb, 0x8d, 0x43, 0x8b, 0x82,
	0x70, 0xc3, 0xd6, 0x1d, 0xf3, 0xbc, 0x0e, 0x64, 0xb9, 0xfd, 0x0e, 0x00, 0x00, 0xff, 0xff, 0x9d,
	0x94, 0xec, 0x1a, 0xe8, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// PluginInterfaceClient is the client API for PluginInterface service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type PluginInterfaceClient interface {
	GetName(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*GetNameResponse, error)
	GetAPIResourceSchema(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*GetAPIResourceSchemaResponse, error)
	Init(ctx context.Context, in *InitRequest, opts ...grpc.CallOption) (*Empty, error)
	Run(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	Stop(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
}

type pluginInterfaceClient struct {
	cc grpc.ClientConnInterface
}

func NewPluginInterfaceClient(cc grpc.ClientConnInterface) PluginInterfaceClient {
	return &pluginInterfaceClient{cc}
}

func (c *pluginInterfaceClient) GetName(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*GetNameResponse, error) {
	out := new(GetNameResponse)
	err := c.cc.Invoke(ctx, "/proto.PluginInterface/GetName", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pluginInterfaceClient) GetAPIResourceSchema(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*GetAPIResourceSchemaResponse, error) {
	out := new(GetAPIResourceSchemaResponse)
	err := c.cc.Invoke(ctx, "/proto.PluginInterface/GetAPIResourceSchema", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pluginInterfaceClient) Init(ctx context.Context, in *InitRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/proto.PluginInterface/Init", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pluginInterfaceClient) Run(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/proto.PluginInterface/Run", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pluginInterfaceClient) Stop(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/proto.PluginInterface/Stop", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PluginInterfaceServer is the server API for PluginInterface service.
type PluginInterfaceServer interface {
	GetName(context.Context, *Empty) (*GetNameResponse, error)
	GetAPIResourceSchema(context.Context, *Empty) (*GetAPIResourceSchemaResponse, error)
	Init(context.Context, *InitRequest) (*Empty, error)
	Run(context.Context, *Empty) (*Empty, error)
	Stop(context.Context, *Empty) (*Empty, error)
}

// UnimplementedPluginInterfaceServer can be embedded to have forward compatible implementations.
type UnimplementedPluginInterfaceServer struct {
}

func (*UnimplementedPluginInterfaceServer) GetName(ctx context.Context, req *Empty) (*GetNameResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetName not implemented")
}
func (*UnimplementedPluginInterfaceServer) GetAPIResourceSchema(ctx context.Context, req *Empty) (*GetAPIResourceSchemaResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAPIResourceSchema not implemented")
}
func (*UnimplementedPluginInterfaceServer) Init(ctx context.Context, req *InitRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Init not implemented")
}
func (*UnimplementedPluginInterfaceServer) Run(ctx context.Context, req *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Run not implemented")
}
func (*UnimplementedPluginInterfaceServer) Stop(ctx context.Context, req *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Stop not implemented")
}

func RegisterPluginInterfaceServer(s *grpc.Server, srv PluginInterfaceServer) {
	s.RegisterService(&_PluginInterface_serviceDesc, srv)
}

func _PluginInterface_GetName_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginInterfaceServer).GetName(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.PluginInterface/GetName",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginInterfaceServer).GetName(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _PluginInterface_GetAPIResourceSchema_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginInterfaceServer).GetAPIResourceSchema(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.PluginInterface/GetAPIResourceSchema",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginInterfaceServer).GetAPIResourceSchema(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _PluginInterface_Init_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InitRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginInterfaceServer).Init(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.PluginInterface/Init",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginInterfaceServer).Init(ctx, req.(*InitRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PluginInterface_Run_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginInterfaceServer).Run(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.PluginInterface/Run",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginInterfaceServer).Run(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _PluginInterface_Stop_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginInterfaceServer).Stop(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.PluginInterface/Stop",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginInterfaceServer).Stop(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _PluginInterface_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.PluginInterface",
	HandlerType: (*PluginInterfaceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetName",
			Handler:    _PluginInterface_GetName_Handler,
		},
		{
			MethodName: "GetAPIResourceSchema",
			Handler:    _PluginInterface_GetAPIResourceSchema_Handler,
		},
		{
			MethodName: "Init",
			Handler:    _PluginInterface_Init_Handler,
		},
		{
			MethodName: "Run",
			Handler:    _PluginInterface_Run_Handler,
		},
		{
			MethodName: "Stop",
			Handler:    _PluginInterface_Stop_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "plugin.proto",
}
