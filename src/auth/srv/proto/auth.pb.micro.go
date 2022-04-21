// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/auth.proto

package auth

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "github.com/micro/go-micro/v2/api"
	client "github.com/micro/go-micro/v2/client"
	server "github.com/micro/go-micro/v2/server"
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

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for AuthService service

func NewAuthServiceEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for AuthService service

type AuthService interface {
	// Token
	VerifyToken(ctx context.Context, in *Token, opts ...client.CallOption) (*BoolMsg, error)
	GetTokenContent(ctx context.Context, in *Token, opts ...client.CallOption) (*TokenContent, error)
	// User
	GetUserById(ctx context.Context, in *IdQuery, opts ...client.CallOption) (*User, error)
	PagedListUsers(ctx context.Context, in *QueryWithPage, opts ...client.CallOption) (*PagedUsers, error)
	GetUserDepartment(ctx context.Context, in *IdQuery, opts ...client.CallOption) (*UserDepartment, error)
	ListUserDepartmentUsers(ctx context.Context, in *IdQuery, opts ...client.CallOption) (*MemberUsers, error)
	MakeUserHandover(ctx context.Context, in *HandoverRequest, opts ...client.CallOption) (*BoolMsg, error)
	// Product
	GetProductById(ctx context.Context, in *IdQuery, opts ...client.CallOption) (*Product, error)
	PagedListProducts(ctx context.Context, in *QueryWithPage, opts ...client.CallOption) (*PagedProducts, error)
	ListProductUsers(ctx context.Context, in *IdQuery, opts ...client.CallOption) (*MemberUsers, error)
	// Department
	GetDepartmentById(ctx context.Context, in *IdQuery, opts ...client.CallOption) (*Department, error)
	ListDepartmentUsers(ctx context.Context, in *IdQuery, opts ...client.CallOption) (*MemberUsers, error)
	ListParentDepartmentUsers(ctx context.Context, in *IdQuery, opts ...client.CallOption) (*MemberUsers, error)
	// Group
	GetGroupById(ctx context.Context, in *IdQuery, opts ...client.CallOption) (*Group, error)
	PagedListGroups(ctx context.Context, in *QueryWithPage, opts ...client.CallOption) (*PagedGroups, error)
	ListGroupUsers(ctx context.Context, in *IdQuery, opts ...client.CallOption) (*MemberUsers, error)
	// Role
	ListRoleUsers(ctx context.Context, in *NameQuery, opts ...client.CallOption) (*RoleMemberUsers, error)
}

type authService struct {
	c    client.Client
	name string
}

func NewAuthService(name string, c client.Client) AuthService {
	return &authService{
		c:    c,
		name: name,
	}
}

func (c *authService) VerifyToken(ctx context.Context, in *Token, opts ...client.CallOption) (*BoolMsg, error) {
	req := c.c.NewRequest(c.name, "AuthService.VerifyToken", in)
	out := new(BoolMsg)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authService) GetTokenContent(ctx context.Context, in *Token, opts ...client.CallOption) (*TokenContent, error) {
	req := c.c.NewRequest(c.name, "AuthService.GetTokenContent", in)
	out := new(TokenContent)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authService) GetUserById(ctx context.Context, in *IdQuery, opts ...client.CallOption) (*User, error) {
	req := c.c.NewRequest(c.name, "AuthService.GetUserById", in)
	out := new(User)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authService) PagedListUsers(ctx context.Context, in *QueryWithPage, opts ...client.CallOption) (*PagedUsers, error) {
	req := c.c.NewRequest(c.name, "AuthService.PagedListUsers", in)
	out := new(PagedUsers)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authService) GetUserDepartment(ctx context.Context, in *IdQuery, opts ...client.CallOption) (*UserDepartment, error) {
	req := c.c.NewRequest(c.name, "AuthService.GetUserDepartment", in)
	out := new(UserDepartment)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authService) ListUserDepartmentUsers(ctx context.Context, in *IdQuery, opts ...client.CallOption) (*MemberUsers, error) {
	req := c.c.NewRequest(c.name, "AuthService.ListUserDepartmentUsers", in)
	out := new(MemberUsers)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authService) MakeUserHandover(ctx context.Context, in *HandoverRequest, opts ...client.CallOption) (*BoolMsg, error) {
	req := c.c.NewRequest(c.name, "AuthService.MakeUserHandover", in)
	out := new(BoolMsg)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authService) GetProductById(ctx context.Context, in *IdQuery, opts ...client.CallOption) (*Product, error) {
	req := c.c.NewRequest(c.name, "AuthService.GetProductById", in)
	out := new(Product)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authService) PagedListProducts(ctx context.Context, in *QueryWithPage, opts ...client.CallOption) (*PagedProducts, error) {
	req := c.c.NewRequest(c.name, "AuthService.PagedListProducts", in)
	out := new(PagedProducts)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authService) ListProductUsers(ctx context.Context, in *IdQuery, opts ...client.CallOption) (*MemberUsers, error) {
	req := c.c.NewRequest(c.name, "AuthService.ListProductUsers", in)
	out := new(MemberUsers)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authService) GetDepartmentById(ctx context.Context, in *IdQuery, opts ...client.CallOption) (*Department, error) {
	req := c.c.NewRequest(c.name, "AuthService.GetDepartmentById", in)
	out := new(Department)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authService) ListDepartmentUsers(ctx context.Context, in *IdQuery, opts ...client.CallOption) (*MemberUsers, error) {
	req := c.c.NewRequest(c.name, "AuthService.ListDepartmentUsers", in)
	out := new(MemberUsers)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authService) ListParentDepartmentUsers(ctx context.Context, in *IdQuery, opts ...client.CallOption) (*MemberUsers, error) {
	req := c.c.NewRequest(c.name, "AuthService.ListParentDepartmentUsers", in)
	out := new(MemberUsers)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authService) GetGroupById(ctx context.Context, in *IdQuery, opts ...client.CallOption) (*Group, error) {
	req := c.c.NewRequest(c.name, "AuthService.GetGroupById", in)
	out := new(Group)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authService) PagedListGroups(ctx context.Context, in *QueryWithPage, opts ...client.CallOption) (*PagedGroups, error) {
	req := c.c.NewRequest(c.name, "AuthService.PagedListGroups", in)
	out := new(PagedGroups)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authService) ListGroupUsers(ctx context.Context, in *IdQuery, opts ...client.CallOption) (*MemberUsers, error) {
	req := c.c.NewRequest(c.name, "AuthService.ListGroupUsers", in)
	out := new(MemberUsers)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authService) ListRoleUsers(ctx context.Context, in *NameQuery, opts ...client.CallOption) (*RoleMemberUsers, error) {
	req := c.c.NewRequest(c.name, "AuthService.ListRoleUsers", in)
	out := new(RoleMemberUsers)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for AuthService service

type AuthServiceHandler interface {
	// Token
	VerifyToken(context.Context, *Token, *BoolMsg) error
	GetTokenContent(context.Context, *Token, *TokenContent) error
	// User
	GetUserById(context.Context, *IdQuery, *User) error
	PagedListUsers(context.Context, *QueryWithPage, *PagedUsers) error
	GetUserDepartment(context.Context, *IdQuery, *UserDepartment) error
	ListUserDepartmentUsers(context.Context, *IdQuery, *MemberUsers) error
	MakeUserHandover(context.Context, *HandoverRequest, *BoolMsg) error
	// Product
	GetProductById(context.Context, *IdQuery, *Product) error
	PagedListProducts(context.Context, *QueryWithPage, *PagedProducts) error
	ListProductUsers(context.Context, *IdQuery, *MemberUsers) error
	// Department
	GetDepartmentById(context.Context, *IdQuery, *Department) error
	ListDepartmentUsers(context.Context, *IdQuery, *MemberUsers) error
	ListParentDepartmentUsers(context.Context, *IdQuery, *MemberUsers) error
	// Group
	GetGroupById(context.Context, *IdQuery, *Group) error
	PagedListGroups(context.Context, *QueryWithPage, *PagedGroups) error
	ListGroupUsers(context.Context, *IdQuery, *MemberUsers) error
	// Role
	ListRoleUsers(context.Context, *NameQuery, *RoleMemberUsers) error
}

func RegisterAuthServiceHandler(s server.Server, hdlr AuthServiceHandler, opts ...server.HandlerOption) error {
	type authService interface {
		VerifyToken(ctx context.Context, in *Token, out *BoolMsg) error
		GetTokenContent(ctx context.Context, in *Token, out *TokenContent) error
		GetUserById(ctx context.Context, in *IdQuery, out *User) error
		PagedListUsers(ctx context.Context, in *QueryWithPage, out *PagedUsers) error
		GetUserDepartment(ctx context.Context, in *IdQuery, out *UserDepartment) error
		ListUserDepartmentUsers(ctx context.Context, in *IdQuery, out *MemberUsers) error
		MakeUserHandover(ctx context.Context, in *HandoverRequest, out *BoolMsg) error
		GetProductById(ctx context.Context, in *IdQuery, out *Product) error
		PagedListProducts(ctx context.Context, in *QueryWithPage, out *PagedProducts) error
		ListProductUsers(ctx context.Context, in *IdQuery, out *MemberUsers) error
		GetDepartmentById(ctx context.Context, in *IdQuery, out *Department) error
		ListDepartmentUsers(ctx context.Context, in *IdQuery, out *MemberUsers) error
		ListParentDepartmentUsers(ctx context.Context, in *IdQuery, out *MemberUsers) error
		GetGroupById(ctx context.Context, in *IdQuery, out *Group) error
		PagedListGroups(ctx context.Context, in *QueryWithPage, out *PagedGroups) error
		ListGroupUsers(ctx context.Context, in *IdQuery, out *MemberUsers) error
		ListRoleUsers(ctx context.Context, in *NameQuery, out *RoleMemberUsers) error
	}
	type AuthService struct {
		authService
	}
	h := &authServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&AuthService{h}, opts...))
}

type authServiceHandler struct {
	AuthServiceHandler
}

func (h *authServiceHandler) VerifyToken(ctx context.Context, in *Token, out *BoolMsg) error {
	return h.AuthServiceHandler.VerifyToken(ctx, in, out)
}

func (h *authServiceHandler) GetTokenContent(ctx context.Context, in *Token, out *TokenContent) error {
	return h.AuthServiceHandler.GetTokenContent(ctx, in, out)
}

func (h *authServiceHandler) GetUserById(ctx context.Context, in *IdQuery, out *User) error {
	return h.AuthServiceHandler.GetUserById(ctx, in, out)
}

func (h *authServiceHandler) PagedListUsers(ctx context.Context, in *QueryWithPage, out *PagedUsers) error {
	return h.AuthServiceHandler.PagedListUsers(ctx, in, out)
}

func (h *authServiceHandler) GetUserDepartment(ctx context.Context, in *IdQuery, out *UserDepartment) error {
	return h.AuthServiceHandler.GetUserDepartment(ctx, in, out)
}

func (h *authServiceHandler) ListUserDepartmentUsers(ctx context.Context, in *IdQuery, out *MemberUsers) error {
	return h.AuthServiceHandler.ListUserDepartmentUsers(ctx, in, out)
}

func (h *authServiceHandler) MakeUserHandover(ctx context.Context, in *HandoverRequest, out *BoolMsg) error {
	return h.AuthServiceHandler.MakeUserHandover(ctx, in, out)
}

func (h *authServiceHandler) GetProductById(ctx context.Context, in *IdQuery, out *Product) error {
	return h.AuthServiceHandler.GetProductById(ctx, in, out)
}

func (h *authServiceHandler) PagedListProducts(ctx context.Context, in *QueryWithPage, out *PagedProducts) error {
	return h.AuthServiceHandler.PagedListProducts(ctx, in, out)
}

func (h *authServiceHandler) ListProductUsers(ctx context.Context, in *IdQuery, out *MemberUsers) error {
	return h.AuthServiceHandler.ListProductUsers(ctx, in, out)
}

func (h *authServiceHandler) GetDepartmentById(ctx context.Context, in *IdQuery, out *Department) error {
	return h.AuthServiceHandler.GetDepartmentById(ctx, in, out)
}

func (h *authServiceHandler) ListDepartmentUsers(ctx context.Context, in *IdQuery, out *MemberUsers) error {
	return h.AuthServiceHandler.ListDepartmentUsers(ctx, in, out)
}

func (h *authServiceHandler) ListParentDepartmentUsers(ctx context.Context, in *IdQuery, out *MemberUsers) error {
	return h.AuthServiceHandler.ListParentDepartmentUsers(ctx, in, out)
}

func (h *authServiceHandler) GetGroupById(ctx context.Context, in *IdQuery, out *Group) error {
	return h.AuthServiceHandler.GetGroupById(ctx, in, out)
}

func (h *authServiceHandler) PagedListGroups(ctx context.Context, in *QueryWithPage, out *PagedGroups) error {
	return h.AuthServiceHandler.PagedListGroups(ctx, in, out)
}

func (h *authServiceHandler) ListGroupUsers(ctx context.Context, in *IdQuery, out *MemberUsers) error {
	return h.AuthServiceHandler.ListGroupUsers(ctx, in, out)
}

func (h *authServiceHandler) ListRoleUsers(ctx context.Context, in *NameQuery, out *RoleMemberUsers) error {
	return h.AuthServiceHandler.ListRoleUsers(ctx, in, out)
}
