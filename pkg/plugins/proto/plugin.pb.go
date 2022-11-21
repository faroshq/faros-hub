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

type GetAPIExportSchemaResponse struct {
	Schema               []byte   `protobuf:"bytes,1,opt,name=schema,proto3" json:"schema,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetAPIExportSchemaResponse) Reset()         { *m = GetAPIExportSchemaResponse{} }
func (m *GetAPIExportSchemaResponse) String() string { return proto.CompactTextString(m) }
func (*GetAPIExportSchemaResponse) ProtoMessage()    {}
func (*GetAPIExportSchemaResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_22a625af4bc1cc87, []int{2}
}

func (m *GetAPIExportSchemaResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetAPIExportSchemaResponse.Unmarshal(m, b)
}
func (m *GetAPIExportSchemaResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetAPIExportSchemaResponse.Marshal(b, m, deterministic)
}
func (m *GetAPIExportSchemaResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetAPIExportSchemaResponse.Merge(m, src)
}
func (m *GetAPIExportSchemaResponse) XXX_Size() int {
	return xxx_messageInfo_GetAPIExportSchemaResponse.Size(m)
}
func (m *GetAPIExportSchemaResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetAPIExportSchemaResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetAPIExportSchemaResponse proto.InternalMessageInfo

func (m *GetAPIExportSchemaResponse) GetSchema() []byte {
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
	return fileDescriptor_22a625af4bc1cc87, []int{3}
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
	return fileDescriptor_22a625af4bc1cc87, []int{4}
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
	proto.RegisterType((*GetAPIExportSchemaResponse)(nil), "proto.GetAPIExportSchemaResponse")
	proto.RegisterType((*InitRequest)(nil), "proto.InitRequest")
	proto.RegisterType((*Empty)(nil), "proto.Empty")
}

func init() {
	proto.RegisterFile("plugin.proto", fileDescriptor_22a625af4bc1cc87)
}

var fileDescriptor_22a625af4bc1cc87 = []byte{
	// 288 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x90, 0x51, 0x4b, 0xc3, 0x30,
	0x14, 0x85, 0xd9, 0xd6, 0x6d, 0xec, 0xae, 0x30, 0x08, 0x32, 0x4a, 0x19, 0xe8, 0x2a, 0xc2, 0x9e,
	0x26, 0xa8, 0xf8, 0x2e, 0x52, 0x46, 0x5f, 0x64, 0x64, 0x3f, 0x40, 0x6b, 0xb9, 0x9b, 0x05, 0x9b,
	0xc4, 0xe6, 0x16, 0xf4, 0x07, 0xfa, 0xbf, 0xa4, 0x77, 0xc5, 0xd9, 0x5a, 0x65, 0x4f, 0x37, 0xf7,
	0xe4, 0x9c, 0x93, 0xf0, 0x81, 0x6b, 0x5e, 0x8b, 0x5d, 0xaa, 0x96, 0x26, 0xd7, 0xa4, 0x45, 0x9f,
	0x47, 0x70, 0x01, 0x93, 0x15, 0xd2, 0x43, 0x9c, 0xa1, 0x44, 0x6b, 0xb4, 0xb2, 0x28, 0x04, 0x38,
	0x2a, 0xce, 0xd0, 0xeb, 0x9c, 0x75, 0x16, 0x23, 0xc9, 0xe7, 0xe0, 0x16, 0x66, 0x2b, 0xa4, 0xbb,
	0x75, 0x24, 0xd1, 0xea, 0x22, 0x4f, 0x70, 0x93, 0xbc, 0x60, 0x16, 0x7f, 0x67, 0xa6, 0x30, 0xb0,
	0xac, 0x70, 0xca, 0x95, 0xd5, 0x16, 0xdc, 0x80, 0xbf, 0xcf, 0x85, 0xef, 0x46, 0xe7, 0x74, 0x64,
	0xea, 0x09, 0xc6, 0x91, 0x4a, 0x49, 0xe2, 0x5b, 0x81, 0x96, 0xda, 0x3e, 0x24, 0x66, 0x30, 0x2a,
	0xa7, 0x35, 0x71, 0x82, 0x5e, 0x97, 0x2f, 0x0e, 0x82, 0x38, 0x85, 0x71, 0x8e, 0x96, 0x1e, 0x13,
	0xad, 0xb6, 0xe9, 0xce, 0xeb, 0x71, 0x3b, 0x94, 0xd2, 0x3d, 0x2b, 0xc1, 0x10, 0xfa, 0x61, 0x66,
	0xe8, 0xe3, 0xea, 0xb3, 0x0b, 0x93, 0x35, 0x73, 0x89, 0x14, 0x61, 0xbe, 0x2d, 0xd3, 0x97, 0x30,
	0xac, 0x98, 0x08, 0x77, 0x4f, 0x6b, 0xc9, 0x66, 0x7f, 0x5a, 0x6d, 0x4d, 0x62, 0x11, 0x9c, 0xb4,
	0xd1, 0x69, 0xa4, 0xcf, 0x0f, 0xe9, 0xbf, 0x41, 0x86, 0x20, 0x7e, 0x03, 0x6b, 0x14, 0xcd, 0x6b,
	0x45, 0xad, 0x64, 0x17, 0xe0, 0x94, 0x04, 0x85, 0xa8, 0xac, 0x3f, 0x70, 0xfa, 0xb5, 0x32, 0x31,
	0x87, 0x9e, 0x2c, 0x54, 0xe3, 0x85, 0xba, 0x25, 0x00, 0x67, 0x43, 0xda, 0xfc, 0xe7, 0x79, 0x1e,
	0xf0, 0x72, 0xfd, 0x15, 0x00, 0x00, 0xff, 0xff, 0x4b, 0xa3, 0xcc, 0x68, 0x65, 0x02, 0x00, 0x00,
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
	GetAPIExportSchema(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*GetAPIExportSchemaResponse, error)
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

func (c *pluginInterfaceClient) GetAPIExportSchema(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*GetAPIExportSchemaResponse, error) {
	out := new(GetAPIExportSchemaResponse)
	err := c.cc.Invoke(ctx, "/proto.PluginInterface/GetAPIExportSchema", in, out, opts...)
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
	GetAPIExportSchema(context.Context, *Empty) (*GetAPIExportSchemaResponse, error)
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
func (*UnimplementedPluginInterfaceServer) GetAPIExportSchema(ctx context.Context, req *Empty) (*GetAPIExportSchemaResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAPIExportSchema not implemented")
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

func _PluginInterface_GetAPIExportSchema_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginInterfaceServer).GetAPIExportSchema(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.PluginInterface/GetAPIExportSchema",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginInterfaceServer).GetAPIExportSchema(ctx, req.(*Empty))
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
			MethodName: "GetAPIExportSchema",
			Handler:    _PluginInterface_GetAPIExportSchema_Handler,
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
