// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/example/example.proto

package go_micro_srv_PostHouses

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
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Request struct {
	Body                 []byte   `protobuf:"bytes,1,opt,name=Body,proto3" json:"Body,omitempty"`
	SessionId            string   `protobuf:"bytes,2,opt,name=SessionId,proto3" json:"SessionId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Request) Reset()         { *m = Request{} }
func (m *Request) String() string { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()    {}
func (*Request) Descriptor() ([]byte, []int) {
	return fileDescriptor_097b3f5db5cf5789, []int{0}
}

func (m *Request) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Request.Unmarshal(m, b)
}
func (m *Request) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Request.Marshal(b, m, deterministic)
}
func (m *Request) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Request.Merge(m, src)
}
func (m *Request) XXX_Size() int {
	return xxx_messageInfo_Request.Size(m)
}
func (m *Request) XXX_DiscardUnknown() {
	xxx_messageInfo_Request.DiscardUnknown(m)
}

var xxx_messageInfo_Request proto.InternalMessageInfo

func (m *Request) GetBody() []byte {
	if m != nil {
		return m.Body
	}
	return nil
}

func (m *Request) GetSessionId() string {
	if m != nil {
		return m.SessionId
	}
	return ""
}

type Response struct {
	Errno                string   `protobuf:"bytes,1,opt,name=Errno,proto3" json:"Errno,omitempty"`
	Errmsg               string   `protobuf:"bytes,2,opt,name=Errmsg,proto3" json:"Errmsg,omitempty"`
	HouseId              string   `protobuf:"bytes,3,opt,name=House_id,json=HouseId,proto3" json:"House_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Response) Reset()         { *m = Response{} }
func (m *Response) String() string { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()    {}
func (*Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_097b3f5db5cf5789, []int{1}
}

func (m *Response) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Response.Unmarshal(m, b)
}
func (m *Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Response.Marshal(b, m, deterministic)
}
func (m *Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Response.Merge(m, src)
}
func (m *Response) XXX_Size() int {
	return xxx_messageInfo_Response.Size(m)
}
func (m *Response) XXX_DiscardUnknown() {
	xxx_messageInfo_Response.DiscardUnknown(m)
}

var xxx_messageInfo_Response proto.InternalMessageInfo

func (m *Response) GetErrno() string {
	if m != nil {
		return m.Errno
	}
	return ""
}

func (m *Response) GetErrmsg() string {
	if m != nil {
		return m.Errmsg
	}
	return ""
}

func (m *Response) GetHouseId() string {
	if m != nil {
		return m.HouseId
	}
	return ""
}

func init() {
	proto.RegisterType((*Request)(nil), "go.micro.srv.PostHouses.Request")
	proto.RegisterType((*Response)(nil), "go.micro.srv.PostHouses.Response")
}

func init() { proto.RegisterFile("proto/example/example.proto", fileDescriptor_097b3f5db5cf5789) }

var fileDescriptor_097b3f5db5cf5789 = []byte{
	// 213 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x2e, 0x28, 0xca, 0x2f,
	0xc9, 0xd7, 0x4f, 0xad, 0x48, 0xcc, 0x2d, 0xc8, 0x49, 0x85, 0xd1, 0x7a, 0x60, 0x51, 0x21, 0xf1,
	0xf4, 0x7c, 0xbd, 0xdc, 0xcc, 0xe4, 0xa2, 0x7c, 0xbd, 0xe2, 0xa2, 0x32, 0xbd, 0x80, 0xfc, 0xe2,
	0x12, 0x8f, 0xfc, 0xd2, 0xe2, 0xd4, 0x62, 0x25, 0x6b, 0x2e, 0xf6, 0xa0, 0xd4, 0xc2, 0xd2, 0xd4,
	0xe2, 0x12, 0x21, 0x21, 0x2e, 0x16, 0xa7, 0xfc, 0x94, 0x4a, 0x09, 0x46, 0x05, 0x46, 0x0d, 0x9e,
	0x20, 0x30, 0x5b, 0x48, 0x86, 0x8b, 0x33, 0x38, 0xb5, 0xb8, 0x38, 0x33, 0x3f, 0xcf, 0x33, 0x45,
	0x82, 0x49, 0x81, 0x51, 0x83, 0x33, 0x08, 0x21, 0xa0, 0x14, 0xcc, 0xc5, 0x11, 0x94, 0x5a, 0x5c,
	0x90, 0x9f, 0x57, 0x9c, 0x2a, 0x24, 0xc2, 0xc5, 0xea, 0x5a, 0x54, 0x94, 0x97, 0x0f, 0xd6, 0xce,
	0x19, 0x04, 0xe1, 0x08, 0x89, 0x71, 0xb1, 0xb9, 0x16, 0x15, 0xe5, 0x16, 0xa7, 0x43, 0x35, 0x43,
	0x79, 0x42, 0x92, 0x5c, 0x1c, 0x60, 0x07, 0xc4, 0x67, 0xa6, 0x48, 0x30, 0x83, 0x65, 0xd8, 0xc1,
	0x7c, 0xcf, 0x14, 0xa3, 0x38, 0x2e, 0x76, 0x57, 0x88, 0xdb, 0x85, 0x82, 0xb9, 0xb8, 0x10, 0x4e,
	0x15, 0x52, 0xd0, 0xc3, 0xe1, 0x09, 0x3d, 0xa8, 0x0f, 0xa4, 0x14, 0xf1, 0xa8, 0x80, 0x38, 0x53,
	0x89, 0x21, 0x89, 0x0d, 0x1c, 0x22, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0x91, 0x8b, 0x1e,
	0x80, 0x30, 0x01, 0x00, 0x00,
}
