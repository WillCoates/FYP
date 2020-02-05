// Code generated by protoc-gen-go. DO NOT EDIT.
// source: sensors.proto

package sensors

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

type SensorInfo struct {
	Unit                 string   `protobuf:"bytes,1,opt,name=unit,proto3" json:"unit,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Sensor               string   `protobuf:"bytes,3,opt,name=sensor,proto3" json:"sensor,omitempty"`
	Measurementname      string   `protobuf:"bytes,4,opt,name=measurementname,proto3" json:"measurementname,omitempty"`
	Measurementunit      string   `protobuf:"bytes,5,opt,name=measurementunit,proto3" json:"measurementunit,omitempty"`
	Hidden               bool     `protobuf:"varint,6,opt,name=hidden,proto3" json:"hidden,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SensorInfo) Reset()         { *m = SensorInfo{} }
func (m *SensorInfo) String() string { return proto.CompactTextString(m) }
func (*SensorInfo) ProtoMessage()    {}
func (*SensorInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_b96d375cab19813b, []int{0}
}

func (m *SensorInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SensorInfo.Unmarshal(m, b)
}
func (m *SensorInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SensorInfo.Marshal(b, m, deterministic)
}
func (m *SensorInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SensorInfo.Merge(m, src)
}
func (m *SensorInfo) XXX_Size() int {
	return xxx_messageInfo_SensorInfo.Size(m)
}
func (m *SensorInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_SensorInfo.DiscardUnknown(m)
}

var xxx_messageInfo_SensorInfo proto.InternalMessageInfo

func (m *SensorInfo) GetUnit() string {
	if m != nil {
		return m.Unit
	}
	return ""
}

func (m *SensorInfo) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *SensorInfo) GetSensor() string {
	if m != nil {
		return m.Sensor
	}
	return ""
}

func (m *SensorInfo) GetMeasurementname() string {
	if m != nil {
		return m.Measurementname
	}
	return ""
}

func (m *SensorInfo) GetMeasurementunit() string {
	if m != nil {
		return m.Measurementunit
	}
	return ""
}

func (m *SensorInfo) GetHidden() bool {
	if m != nil {
		return m.Hidden
	}
	return false
}

type SensorData struct {
	Sensor               string   `protobuf:"bytes,1,opt,name=sensor,proto3" json:"sensor,omitempty"`
	Reading              string   `protobuf:"bytes,2,opt,name=reading,proto3" json:"reading,omitempty"`
	Measurementname      string   `protobuf:"bytes,3,opt,name=measurementname,proto3" json:"measurementname,omitempty"`
	Measurementunit      string   `protobuf:"bytes,4,opt,name=measurementunit,proto3" json:"measurementunit,omitempty"`
	Timestamp            int64    `protobuf:"varint,5,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Unit                 string   `protobuf:"bytes,6,opt,name=unit,proto3" json:"unit,omitempty"`
	UnitName             string   `protobuf:"bytes,7,opt,name=unitName,proto3" json:"unitName,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SensorData) Reset()         { *m = SensorData{} }
func (m *SensorData) String() string { return proto.CompactTextString(m) }
func (*SensorData) ProtoMessage()    {}
func (*SensorData) Descriptor() ([]byte, []int) {
	return fileDescriptor_b96d375cab19813b, []int{1}
}

func (m *SensorData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SensorData.Unmarshal(m, b)
}
func (m *SensorData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SensorData.Marshal(b, m, deterministic)
}
func (m *SensorData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SensorData.Merge(m, src)
}
func (m *SensorData) XXX_Size() int {
	return xxx_messageInfo_SensorData.Size(m)
}
func (m *SensorData) XXX_DiscardUnknown() {
	xxx_messageInfo_SensorData.DiscardUnknown(m)
}

var xxx_messageInfo_SensorData proto.InternalMessageInfo

func (m *SensorData) GetSensor() string {
	if m != nil {
		return m.Sensor
	}
	return ""
}

func (m *SensorData) GetReading() string {
	if m != nil {
		return m.Reading
	}
	return ""
}

func (m *SensorData) GetMeasurementname() string {
	if m != nil {
		return m.Measurementname
	}
	return ""
}

func (m *SensorData) GetMeasurementunit() string {
	if m != nil {
		return m.Measurementunit
	}
	return ""
}

func (m *SensorData) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *SensorData) GetUnit() string {
	if m != nil {
		return m.Unit
	}
	return ""
}

func (m *SensorData) GetUnitName() string {
	if m != nil {
		return m.UnitName
	}
	return ""
}

type GetSensorsRequest struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	IncludeHidden        bool     `protobuf:"varint,2,opt,name=includeHidden,proto3" json:"includeHidden,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetSensorsRequest) Reset()         { *m = GetSensorsRequest{} }
func (m *GetSensorsRequest) String() string { return proto.CompactTextString(m) }
func (*GetSensorsRequest) ProtoMessage()    {}
func (*GetSensorsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_b96d375cab19813b, []int{2}
}

func (m *GetSensorsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetSensorsRequest.Unmarshal(m, b)
}
func (m *GetSensorsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetSensorsRequest.Marshal(b, m, deterministic)
}
func (m *GetSensorsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetSensorsRequest.Merge(m, src)
}
func (m *GetSensorsRequest) XXX_Size() int {
	return xxx_messageInfo_GetSensorsRequest.Size(m)
}
func (m *GetSensorsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetSensorsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetSensorsRequest proto.InternalMessageInfo

func (m *GetSensorsRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *GetSensorsRequest) GetIncludeHidden() bool {
	if m != nil {
		return m.IncludeHidden
	}
	return false
}

type GetSensorReadingsRequest struct {
	Unit                 []string `protobuf:"bytes,1,rep,name=unit,proto3" json:"unit,omitempty"`
	Sensor               []string `protobuf:"bytes,2,rep,name=sensor,proto3" json:"sensor,omitempty"`
	Since                int64    `protobuf:"varint,3,opt,name=since,proto3" json:"since,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetSensorReadingsRequest) Reset()         { *m = GetSensorReadingsRequest{} }
func (m *GetSensorReadingsRequest) String() string { return proto.CompactTextString(m) }
func (*GetSensorReadingsRequest) ProtoMessage()    {}
func (*GetSensorReadingsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_b96d375cab19813b, []int{3}
}

func (m *GetSensorReadingsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetSensorReadingsRequest.Unmarshal(m, b)
}
func (m *GetSensorReadingsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetSensorReadingsRequest.Marshal(b, m, deterministic)
}
func (m *GetSensorReadingsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetSensorReadingsRequest.Merge(m, src)
}
func (m *GetSensorReadingsRequest) XXX_Size() int {
	return xxx_messageInfo_GetSensorReadingsRequest.Size(m)
}
func (m *GetSensorReadingsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetSensorReadingsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetSensorReadingsRequest proto.InternalMessageInfo

func (m *GetSensorReadingsRequest) GetUnit() []string {
	if m != nil {
		return m.Unit
	}
	return nil
}

func (m *GetSensorReadingsRequest) GetSensor() []string {
	if m != nil {
		return m.Sensor
	}
	return nil
}

func (m *GetSensorReadingsRequest) GetSince() int64 {
	if m != nil {
		return m.Since
	}
	return 0
}

func init() {
	proto.RegisterType((*SensorInfo)(nil), "SensorInfo")
	proto.RegisterType((*SensorData)(nil), "SensorData")
	proto.RegisterType((*GetSensorsRequest)(nil), "GetSensorsRequest")
	proto.RegisterType((*GetSensorReadingsRequest)(nil), "GetSensorReadingsRequest")
}

func init() { proto.RegisterFile("sensors.proto", fileDescriptor_b96d375cab19813b) }

var fileDescriptor_b96d375cab19813b = []byte{
	// 351 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0xbf, 0x4e, 0xc3, 0x30,
	0x10, 0xc6, 0xe5, 0xa6, 0x4d, 0xdb, 0x83, 0x82, 0xb0, 0x50, 0x15, 0x2a, 0x86, 0x2a, 0x62, 0xc8,
	0x14, 0x2a, 0x98, 0xd9, 0x90, 0x80, 0x01, 0x86, 0x54, 0x6c, 0x2c, 0xa6, 0x39, 0xc0, 0x12, 0x71,
	0x4a, 0xec, 0xf0, 0x50, 0x3c, 0x01, 0x8f, 0xc3, 0xa3, 0xa0, 0x9c, 0xd3, 0xc4, 0xfd, 0x83, 0xc4,
	0x94, 0xdc, 0xa7, 0xf3, 0xdd, 0xef, 0xbe, 0x3b, 0x18, 0x69, 0x54, 0x3a, 0x2f, 0x74, 0xbc, 0x2c,
	0x72, 0x93, 0x87, 0xdf, 0x0c, 0x60, 0x4e, 0xca, 0x9d, 0x7a, 0xc9, 0x39, 0x87, 0x6e, 0xa9, 0xa4,
	0x09, 0xd8, 0x94, 0x45, 0xc3, 0x84, 0xfe, 0x2b, 0x4d, 0x89, 0x0c, 0x83, 0x8e, 0xd5, 0xaa, 0x7f,
	0x3e, 0x06, 0xdf, 0xd6, 0x09, 0x3c, 0x52, 0xeb, 0x88, 0x47, 0x70, 0x98, 0xa1, 0xd0, 0x65, 0x81,
	0x19, 0x2a, 0x43, 0xcf, 0xba, 0x94, 0xb0, 0x29, 0x6f, 0x64, 0x52, 0xd3, 0xde, 0x56, 0x26, 0xf5,
	0x1f, 0x83, 0xff, 0x26, 0xd3, 0x14, 0x55, 0xe0, 0x4f, 0x59, 0x34, 0x48, 0xea, 0x28, 0xfc, 0x69,
	0xd0, 0xaf, 0x85, 0x11, 0x0e, 0x12, 0x5b, 0x43, 0x0a, 0xa0, 0x5f, 0xa0, 0x48, 0xa5, 0x7a, 0xad,
	0x27, 0x58, 0x85, 0xbb, 0x60, 0xbd, 0x7f, 0xc3, 0x76, 0x77, 0xc3, 0x9e, 0xc2, 0xd0, 0xc8, 0x0c,
	0xb5, 0x11, 0xd9, 0x92, 0x06, 0xf2, 0x92, 0x56, 0x68, 0xec, 0xf5, 0x1d, 0x7b, 0x27, 0x30, 0xa8,
	0xbe, 0x0f, 0x55, 0xfb, 0x3e, 0xe9, 0x4d, 0x1c, 0xde, 0xc3, 0xd1, 0x0d, 0x1a, 0x3b, 0xa4, 0x4e,
	0xf0, 0xa3, 0x44, 0xdd, 0xee, 0x83, 0x39, 0xfb, 0x38, 0x83, 0x91, 0x54, 0x8b, 0xf7, 0x32, 0xc5,
	0x5b, 0x6b, 0x55, 0x87, 0xac, 0x5a, 0x17, 0xc3, 0x27, 0x08, 0x9a, 0x72, 0x89, 0x35, 0xc1, 0xad,
	0x5a, 0x6f, 0xde, 0x6b, 0xd0, 0x5a, 0x4b, 0x3b, 0xa4, 0xae, 0x2c, 0x3d, 0x86, 0x9e, 0x96, 0x6a,
	0x61, 0xed, 0xf2, 0x12, 0x1b, 0x5c, 0x7c, 0x31, 0x38, 0xa8, 0x51, 0xe7, 0x58, 0x7c, 0xca, 0x05,
	0xf2, 0x73, 0x80, 0x96, 0x9f, 0xf3, 0x78, 0x6b, 0x98, 0xc9, 0x5e, 0xdc, 0x5e, 0xdf, 0x8c, 0xf1,
	0x2b, 0x67, 0xe0, 0x15, 0x21, 0x3f, 0x89, 0xff, 0xa2, 0x6e, 0x9e, 0x57, 0x17, 0x30, 0x63, 0x3c,
	0x82, 0xfd, 0xc7, 0x65, 0x2a, 0x0c, 0x5a, 0x95, 0xbb, 0xd5, 0xd7, 0x5a, 0x3d, 0xfb, 0x74, 0xfe,
	0x97, 0xbf, 0x01, 0x00, 0x00, 0xff, 0xff, 0xcd, 0x3b, 0xda, 0xec, 0x0f, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// SensorsServiceClient is the client API for SensorsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type SensorsServiceClient interface {
	GetSensors(ctx context.Context, in *GetSensorsRequest, opts ...grpc.CallOption) (SensorsService_GetSensorsClient, error)
	GetSensorReadings(ctx context.Context, in *GetSensorReadingsRequest, opts ...grpc.CallOption) (SensorsService_GetSensorReadingsClient, error)
	UpdateSensor(ctx context.Context, in *SensorInfo, opts ...grpc.CallOption) (*SensorInfo, error)
}

type sensorsServiceClient struct {
	cc *grpc.ClientConn
}

func NewSensorsServiceClient(cc *grpc.ClientConn) SensorsServiceClient {
	return &sensorsServiceClient{cc}
}

func (c *sensorsServiceClient) GetSensors(ctx context.Context, in *GetSensorsRequest, opts ...grpc.CallOption) (SensorsService_GetSensorsClient, error) {
	stream, err := c.cc.NewStream(ctx, &_SensorsService_serviceDesc.Streams[0], "/SensorsService/GetSensors", opts...)
	if err != nil {
		return nil, err
	}
	x := &sensorsServiceGetSensorsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type SensorsService_GetSensorsClient interface {
	Recv() (*SensorInfo, error)
	grpc.ClientStream
}

type sensorsServiceGetSensorsClient struct {
	grpc.ClientStream
}

func (x *sensorsServiceGetSensorsClient) Recv() (*SensorInfo, error) {
	m := new(SensorInfo)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *sensorsServiceClient) GetSensorReadings(ctx context.Context, in *GetSensorReadingsRequest, opts ...grpc.CallOption) (SensorsService_GetSensorReadingsClient, error) {
	stream, err := c.cc.NewStream(ctx, &_SensorsService_serviceDesc.Streams[1], "/SensorsService/GetSensorReadings", opts...)
	if err != nil {
		return nil, err
	}
	x := &sensorsServiceGetSensorReadingsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type SensorsService_GetSensorReadingsClient interface {
	Recv() (*SensorData, error)
	grpc.ClientStream
}

type sensorsServiceGetSensorReadingsClient struct {
	grpc.ClientStream
}

func (x *sensorsServiceGetSensorReadingsClient) Recv() (*SensorData, error) {
	m := new(SensorData)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *sensorsServiceClient) UpdateSensor(ctx context.Context, in *SensorInfo, opts ...grpc.CallOption) (*SensorInfo, error) {
	out := new(SensorInfo)
	err := c.cc.Invoke(ctx, "/SensorsService/UpdateSensor", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SensorsServiceServer is the server API for SensorsService service.
type SensorsServiceServer interface {
	GetSensors(*GetSensorsRequest, SensorsService_GetSensorsServer) error
	GetSensorReadings(*GetSensorReadingsRequest, SensorsService_GetSensorReadingsServer) error
	UpdateSensor(context.Context, *SensorInfo) (*SensorInfo, error)
}

// UnimplementedSensorsServiceServer can be embedded to have forward compatible implementations.
type UnimplementedSensorsServiceServer struct {
}

func (*UnimplementedSensorsServiceServer) GetSensors(req *GetSensorsRequest, srv SensorsService_GetSensorsServer) error {
	return status.Errorf(codes.Unimplemented, "method GetSensors not implemented")
}
func (*UnimplementedSensorsServiceServer) GetSensorReadings(req *GetSensorReadingsRequest, srv SensorsService_GetSensorReadingsServer) error {
	return status.Errorf(codes.Unimplemented, "method GetSensorReadings not implemented")
}
func (*UnimplementedSensorsServiceServer) UpdateSensor(ctx context.Context, req *SensorInfo) (*SensorInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateSensor not implemented")
}

func RegisterSensorsServiceServer(s *grpc.Server, srv SensorsServiceServer) {
	s.RegisterService(&_SensorsService_serviceDesc, srv)
}

func _SensorsService_GetSensors_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GetSensorsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SensorsServiceServer).GetSensors(m, &sensorsServiceGetSensorsServer{stream})
}

type SensorsService_GetSensorsServer interface {
	Send(*SensorInfo) error
	grpc.ServerStream
}

type sensorsServiceGetSensorsServer struct {
	grpc.ServerStream
}

func (x *sensorsServiceGetSensorsServer) Send(m *SensorInfo) error {
	return x.ServerStream.SendMsg(m)
}

func _SensorsService_GetSensorReadings_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GetSensorReadingsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SensorsServiceServer).GetSensorReadings(m, &sensorsServiceGetSensorReadingsServer{stream})
}

type SensorsService_GetSensorReadingsServer interface {
	Send(*SensorData) error
	grpc.ServerStream
}

type sensorsServiceGetSensorReadingsServer struct {
	grpc.ServerStream
}

func (x *sensorsServiceGetSensorReadingsServer) Send(m *SensorData) error {
	return x.ServerStream.SendMsg(m)
}

func _SensorsService_UpdateSensor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SensorInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SensorsServiceServer).UpdateSensor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/SensorsService/UpdateSensor",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SensorsServiceServer).UpdateSensor(ctx, req.(*SensorInfo))
	}
	return interceptor(ctx, in, info, handler)
}

var _SensorsService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "SensorsService",
	HandlerType: (*SensorsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UpdateSensor",
			Handler:    _SensorsService_UpdateSensor_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetSensors",
			Handler:       _SensorsService_GetSensors_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "GetSensorReadings",
			Handler:       _SensorsService_GetSensorReadings_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "sensors.proto",
}
