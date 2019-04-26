// Code generated by protoc-gen-go. DO NOT EDIT.
// source: auth.proto

package devroute_auth

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

type LogInRequest struct {
	Username             string   `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Password             string   `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LogInRequest) Reset()         { *m = LogInRequest{} }
func (m *LogInRequest) String() string { return proto.CompactTextString(m) }
func (*LogInRequest) ProtoMessage()    {}
func (*LogInRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8bbd6f3875b0e874, []int{0}
}

func (m *LogInRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LogInRequest.Unmarshal(m, b)
}
func (m *LogInRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LogInRequest.Marshal(b, m, deterministic)
}
func (m *LogInRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LogInRequest.Merge(m, src)
}
func (m *LogInRequest) XXX_Size() int {
	return xxx_messageInfo_LogInRequest.Size(m)
}
func (m *LogInRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_LogInRequest.DiscardUnknown(m)
}

var xxx_messageInfo_LogInRequest proto.InternalMessageInfo

func (m *LogInRequest) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *LogInRequest) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

type LogInResponse struct {
	Token                string   `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LogInResponse) Reset()         { *m = LogInResponse{} }
func (m *LogInResponse) String() string { return proto.CompactTextString(m) }
func (*LogInResponse) ProtoMessage()    {}
func (*LogInResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_8bbd6f3875b0e874, []int{1}
}

func (m *LogInResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LogInResponse.Unmarshal(m, b)
}
func (m *LogInResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LogInResponse.Marshal(b, m, deterministic)
}
func (m *LogInResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LogInResponse.Merge(m, src)
}
func (m *LogInResponse) XXX_Size() int {
	return xxx_messageInfo_LogInResponse.Size(m)
}
func (m *LogInResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_LogInResponse.DiscardUnknown(m)
}

var xxx_messageInfo_LogInResponse proto.InternalMessageInfo

func (m *LogInResponse) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func init() {
	proto.RegisterType((*LogInRequest)(nil), "devroute.auth.LogInRequest")
	proto.RegisterType((*LogInResponse)(nil), "devroute.auth.LogInResponse")
}

func init() { proto.RegisterFile("auth.proto", fileDescriptor_8bbd6f3875b0e874) }

var fileDescriptor_8bbd6f3875b0e874 = []byte{
	// 167 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4a, 0x2c, 0x2d, 0xc9,
	0xd0, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x4d, 0x49, 0x2d, 0x2b, 0xca, 0x2f, 0x2d, 0x49,
	0xd5, 0x03, 0x09, 0x2a, 0xb9, 0x71, 0xf1, 0xf8, 0xe4, 0xa7, 0x7b, 0xe6, 0x05, 0xa5, 0x16, 0x96,
	0xa6, 0x16, 0x97, 0x08, 0x49, 0x71, 0x71, 0x94, 0x16, 0xa7, 0x16, 0xe5, 0x25, 0xe6, 0xa6, 0x4a,
	0x30, 0x2a, 0x30, 0x6a, 0x70, 0x06, 0xc1, 0xf9, 0x20, 0xb9, 0x82, 0xc4, 0xe2, 0xe2, 0xf2, 0xfc,
	0xa2, 0x14, 0x09, 0x26, 0x88, 0x1c, 0x8c, 0xaf, 0xa4, 0xca, 0xc5, 0x0b, 0x35, 0xa7, 0xb8, 0x20,
	0x3f, 0xaf, 0x38, 0x55, 0x48, 0x84, 0x8b, 0xb5, 0x24, 0x3f, 0x3b, 0x35, 0x0f, 0x6a, 0x0a, 0x84,
	0x63, 0x14, 0xc8, 0xc5, 0x1d, 0x5a, 0x9c, 0x5a, 0x14, 0x9c, 0x5a, 0x54, 0x96, 0x99, 0x9c, 0x2a,
	0xe4, 0xc4, 0xc5, 0x0a, 0xd6, 0x25, 0x24, 0xad, 0x87, 0xe2, 0x2c, 0x3d, 0x64, 0x37, 0x49, 0xc9,
	0x60, 0x97, 0x84, 0x58, 0x94, 0xc4, 0x06, 0xf6, 0x97, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff, 0x0e,
	0x08, 0x3f, 0xbe, 0xe5, 0x00, 0x00, 0x00,
}
