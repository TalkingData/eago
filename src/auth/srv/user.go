package main

import (
	"context"
	"eago/auth/conf/msg"
	"eago/auth/model"
	"eago/auth/srv/proto"
	"eago/common/log"
)

// ListUsers RPC服务::列出所有用户
func (as *AuthService) ListUsers(ctx context.Context, req *auth.QueryWithPage, rsp *auth.PagedUsers) error {
	log.Info("srv.ListUsers called.")
	defer log.Info("srv.ListUsers end.")

	query := make(model.Query)
	for k, v := range query {
		query[k] = v
	}

	pagedData, ok := model.PagedListUsers(query, int(req.Page), int(req.PageSize))
	if !ok {
		m := msg.ErrDatabase.SetDetail("Error when AuthService.ListUsers called, in model.PagedListUsers.")
		log.Error(m.String())
		return m.Error()
	}

	rsp.Users = make([]*auth.User, 0)
	for _, u := range pagedData.Data.([]model.User) {
		newU := auth.User{}
		newU.Id = int32(u.Id)
		newU.Username = u.Username
		newU.Email = u.Email
		newU.Phone = u.Phone
		rsp.Users = append(rsp.Users, &newU)
	}

	rsp.Page = uint32(pagedData.Page)
	rsp.Pages = uint32(pagedData.Pages)
	rsp.PageSize = uint32(pagedData.PageSize)
	rsp.Total = uint32(pagedData.Total)
	return nil
}

// GetUserDepartment RPC服务::获得用户所在部门
func (as *AuthService) GetUserDepartment(ctx context.Context, in *auth.IdQuery, out *auth.Department) error {
	log.Info("srv.GetUserDepartment called.")
	defer log.Info("srv.GetUserDepartment end.")

	dept, ok := model.GetUserDepartment(int(in.Id))
	if !ok {
		m := msg.ErrDatabase.GenResponse("Error when AuthService.GetUserDepartment in model.GetUserDepartment.")
		log.Error(m.String())
		return m.Error()
	}

	log.Info("Making response.")
	if dept != nil {
		out.Id = int32(dept.Id)
		out.Name = dept.Name
		if dept.ParentId != nil {
			out.ParentId = int32(*dept.ParentId)
		}
	}

	return nil
}

// ListUserDepartmentUsers RPC服务::列出用户所在部门用户
func (as *AuthService) ListUserDepartmentUsers(ctx context.Context, in *auth.IdQuery, out *auth.MemberUsers) error {
	log.Info("srv.ListUserDepartmentUsers called.")
	defer log.Info("srv.ListUserDepartmentUsers end.")

	out.Users = make([]*auth.MemberUsers_MemberUser, 0)

	log.Info("Finding user department.")
	dept, ok := model.GetUserDepartment(int(in.Id))
	if !ok {
		m := msg.ErrDatabase.GenResponse("Error when AuthService.GetUserDepartment called, in model.GetUserDepartment.")
		log.Error(m.String())
		return m.Error()
	}
	if dept == nil {
		log.DebugWithFields(log.Fields{"department_id": in.Id}, "Department not found.")
		return nil
	}

	log.Info("Finding department users.")
	mem, ok := model.ListDepartmentUsers(dept.Id, model.Query{})
	if !ok {
		m := msg.ErrDatabase.GenResponse("Error when AuthService.GetUserDepartment called, in model.ListDepartmentUsers.")
		log.Error(m.String())
		return m.Error()
	}

	log.Info("Making response.")
	for _, u := range *mem {
		mUser := auth.MemberUsers_MemberUser{}
		mUser.Id = int32(u.Id)
		mUser.Username = u.Username
		mUser.IsOwner = u.IsOwner
		out.Users = append(out.Users, &mUser)
	}

	return nil
}
