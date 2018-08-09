// Code generated by protoc-gen-go. DO NOT EDIT.
// source: job.proto

package job

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Id struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Id) Reset()         { *m = Id{} }
func (m *Id) String() string { return proto.CompactTextString(m) }
func (*Id) ProtoMessage()    {}
func (*Id) Descriptor() ([]byte, []int) {
	return fileDescriptor_job_570c70dcec76efeb, []int{0}
}
func (m *Id) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Id.Unmarshal(m, b)
}
func (m *Id) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Id.Marshal(b, m, deterministic)
}
func (dst *Id) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Id.Merge(dst, src)
}
func (m *Id) XXX_Size() int {
	return xxx_messageInfo_Id.Size(m)
}
func (m *Id) XXX_DiscardUnknown() {
	xxx_messageInfo_Id.DiscardUnknown(m)
}

var xxx_messageInfo_Id proto.InternalMessageInfo

func (m *Id) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type IdList struct {
	Ids                  []*Id    `protobuf:"bytes,1,rep,name=ids,proto3" json:"ids,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *IdList) Reset()         { *m = IdList{} }
func (m *IdList) String() string { return proto.CompactTextString(m) }
func (*IdList) ProtoMessage()    {}
func (*IdList) Descriptor() ([]byte, []int) {
	return fileDescriptor_job_570c70dcec76efeb, []int{1}
}
func (m *IdList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IdList.Unmarshal(m, b)
}
func (m *IdList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IdList.Marshal(b, m, deterministic)
}
func (dst *IdList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IdList.Merge(dst, src)
}
func (m *IdList) XXX_Size() int {
	return xxx_messageInfo_IdList.Size(m)
}
func (m *IdList) XXX_DiscardUnknown() {
	xxx_messageInfo_IdList.DiscardUnknown(m)
}

var xxx_messageInfo_IdList proto.InternalMessageInfo

func (m *IdList) GetIds() []*Id {
	if m != nil {
		return m.Ids
	}
	return nil
}

type Job struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Body                 string   `protobuf:"bytes,2,opt,name=body,proto3" json:"body,omitempty"`
	Unique               bool     `protobuf:"varint,3,opt,name=unique,proto3" json:"unique,omitempty"`
	When                 int32    `protobuf:"varint,4,opt,name=when,proto3" json:"when,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Job) Reset()         { *m = Job{} }
func (m *Job) String() string { return proto.CompactTextString(m) }
func (*Job) ProtoMessage()    {}
func (*Job) Descriptor() ([]byte, []int) {
	return fileDescriptor_job_570c70dcec76efeb, []int{2}
}
func (m *Job) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Job.Unmarshal(m, b)
}
func (m *Job) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Job.Marshal(b, m, deterministic)
}
func (dst *Job) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Job.Merge(dst, src)
}
func (m *Job) XXX_Size() int {
	return xxx_messageInfo_Job.Size(m)
}
func (m *Job) XXX_DiscardUnknown() {
	xxx_messageInfo_Job.DiscardUnknown(m)
}

var xxx_messageInfo_Job proto.InternalMessageInfo

func (m *Job) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Job) GetBody() string {
	if m != nil {
		return m.Body
	}
	return ""
}

func (m *Job) GetUnique() bool {
	if m != nil {
		return m.Unique
	}
	return false
}

func (m *Job) GetWhen() int32 {
	if m != nil {
		return m.When
	}
	return 0
}

type JobList struct {
	Jobs                 []*Job   `protobuf:"bytes,1,rep,name=jobs,proto3" json:"jobs,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *JobList) Reset()         { *m = JobList{} }
func (m *JobList) String() string { return proto.CompactTextString(m) }
func (*JobList) ProtoMessage()    {}
func (*JobList) Descriptor() ([]byte, []int) {
	return fileDescriptor_job_570c70dcec76efeb, []int{3}
}
func (m *JobList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_JobList.Unmarshal(m, b)
}
func (m *JobList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_JobList.Marshal(b, m, deterministic)
}
func (dst *JobList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_JobList.Merge(dst, src)
}
func (m *JobList) XXX_Size() int {
	return xxx_messageInfo_JobList.Size(m)
}
func (m *JobList) XXX_DiscardUnknown() {
	xxx_messageInfo_JobList.DiscardUnknown(m)
}

var xxx_messageInfo_JobList proto.InternalMessageInfo

func (m *JobList) GetJobs() []*Job {
	if m != nil {
		return m.Jobs
	}
	return nil
}

type Void struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Void) Reset()         { *m = Void{} }
func (m *Void) String() string { return proto.CompactTextString(m) }
func (*Void) ProtoMessage()    {}
func (*Void) Descriptor() ([]byte, []int) {
	return fileDescriptor_job_570c70dcec76efeb, []int{4}
}
func (m *Void) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Void.Unmarshal(m, b)
}
func (m *Void) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Void.Marshal(b, m, deterministic)
}
func (dst *Void) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Void.Merge(dst, src)
}
func (m *Void) XXX_Size() int {
	return xxx_messageInfo_Void.Size(m)
}
func (m *Void) XXX_DiscardUnknown() {
	xxx_messageInfo_Void.DiscardUnknown(m)
}

var xxx_messageInfo_Void proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Id)(nil), "Id")
	proto.RegisterType((*IdList)(nil), "IdList")
	proto.RegisterType((*Job)(nil), "Job")
	proto.RegisterType((*JobList)(nil), "JobList")
	proto.RegisterType((*Void)(nil), "Void")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// JobsClient is the client API for Jobs service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type JobsClient interface {
	Push(ctx context.Context, in *JobList, opts ...grpc.CallOption) (*IdList, error)
	Remove(ctx context.Context, in *IdList, opts ...grpc.CallOption) (*Void, error)
}

type jobsClient struct {
	cc *grpc.ClientConn
}

func NewJobsClient(cc *grpc.ClientConn) JobsClient {
	return &jobsClient{cc}
}

func (c *jobsClient) Push(ctx context.Context, in *JobList, opts ...grpc.CallOption) (*IdList, error) {
	out := new(IdList)
	err := c.cc.Invoke(ctx, "/Jobs/Push", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *jobsClient) Remove(ctx context.Context, in *IdList, opts ...grpc.CallOption) (*Void, error) {
	out := new(Void)
	err := c.cc.Invoke(ctx, "/Jobs/Remove", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// JobsServer is the server API for Jobs service.
type JobsServer interface {
	Push(context.Context, *JobList) (*IdList, error)
	Remove(context.Context, *IdList) (*Void, error)
}

func RegisterJobsServer(s *grpc.Server, srv JobsServer) {
	s.RegisterService(&_Jobs_serviceDesc, srv)
}

func _Jobs_Push_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JobList)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JobsServer).Push(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Jobs/Push",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JobsServer).Push(ctx, req.(*JobList))
	}
	return interceptor(ctx, in, info, handler)
}

func _Jobs_Remove_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IdList)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JobsServer).Remove(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Jobs/Remove",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JobsServer).Remove(ctx, req.(*IdList))
	}
	return interceptor(ctx, in, info, handler)
}

var _Jobs_serviceDesc = grpc.ServiceDesc{
	ServiceName: "Jobs",
	HandlerType: (*JobsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Push",
			Handler:    _Jobs_Push_Handler,
		},
		{
			MethodName: "Remove",
			Handler:    _Jobs_Remove_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "job.proto",
}

func init() { proto.RegisterFile("job.proto", fileDescriptor_job_570c70dcec76efeb) }

var fileDescriptor_job_570c70dcec76efeb = []byte{
	// 216 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x8f, 0xc1, 0x4a, 0x03, 0x31,
	0x10, 0x86, 0xc9, 0x26, 0x4d, 0xdb, 0x11, 0x3c, 0x0c, 0x2a, 0xb1, 0x17, 0x97, 0x78, 0xd9, 0x53,
	0x0e, 0xf5, 0xe8, 0x13, 0xec, 0xe2, 0x41, 0x02, 0xfa, 0x00, 0xc3, 0x04, 0x9a, 0x82, 0x1d, 0x35,
	0x5d, 0xc5, 0xb7, 0x97, 0x0d, 0xbb, 0x97, 0xde, 0xfe, 0xff, 0x23, 0x1f, 0x7f, 0x06, 0xb6, 0x47,
	0xa1, 0xf0, 0xf9, 0x2d, 0x67, 0xf1, 0x37, 0xd0, 0xf4, 0x8c, 0xd7, 0xd0, 0x64, 0x76, 0xaa, 0x55,
	0xdd, 0x36, 0x36, 0x99, 0xfd, 0x03, 0xd8, 0x9e, 0x5f, 0x72, 0x39, 0xe3, 0x2d, 0xe8, 0xcc, 0xc5,
	0xa9, 0x56, 0x77, 0x57, 0x7b, 0x1d, 0x7a, 0x8e, 0x53, 0xf7, 0x6f, 0xa0, 0x07, 0xa1, 0x4b, 0x0f,
	0x11, 0x0c, 0x09, 0xff, 0xb9, 0xa6, 0x92, 0x9a, 0xf1, 0x0e, 0xec, 0x78, 0xca, 0x5f, 0x63, 0x72,
	0xba, 0x55, 0xdd, 0x26, 0xce, 0x6d, 0x7a, 0xfb, 0x7b, 0x48, 0x27, 0x67, 0x5a, 0xd5, 0xad, 0x62,
	0xcd, 0xfe, 0x11, 0xd6, 0x83, 0x50, 0x1d, 0x76, 0x60, 0x8e, 0x42, 0xcb, 0xb2, 0x09, 0x83, 0x50,
	0xac, 0xc4, 0x5b, 0x30, 0xef, 0x92, 0x79, 0xff, 0x0c, 0x66, 0x10, 0x2a, 0x78, 0x0f, 0xe6, 0x75,
	0x2c, 0x07, 0xdc, 0x84, 0xd9, 0xdd, 0xad, 0xc3, 0xfc, 0x7b, 0x07, 0x36, 0xa6, 0x0f, 0xf9, 0x49,
	0xb8, 0xa0, 0xdd, 0x2a, 0x4c, 0x32, 0xd9, 0x7a, 0xfe, 0xd3, 0x7f, 0x00, 0x00, 0x00, 0xff, 0xff,
	0xc3, 0xcd, 0x7b, 0x34, 0x0b, 0x01, 0x00, 0x00,
}
