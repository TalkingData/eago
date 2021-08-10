package main

import (
	"context"
	"eago/auth/conf/msg"
	"eago/auth/model"
	"eago/auth/srv/proto"
	"eago/common/log"
)

// ListDepartmentUsers RPC服务::列出指定部门下所有用户
func (as *AuthService) ListDepartmentUsers(ctx context.Context, in *auth.IdQuery, out *auth.MemberUsers) error {
	log.InfoWithFields(log.Fields{"department_id": in.Id}, "srv.ListDepartmentUsers called.")
	defer log.Info("srv.ListDepartmentUsers end.")

	log.Info("Finding department users.")
	us, ok := model.ListDepartmentUsers(int(in.Id), model.Query{})
	if !ok {
		m := msg.ErrDatabase.GenResponse("Error when AuthService.ListDepartmentUsers called, in model.ListDepartmentUsers.")
		log.Error(m.String())
		return m.Error()
	}

	log.Info("Making response.")
	out.Users = make([]*auth.MemberUsers_MemberUser, 0)
	for _, u := range *us {
		mUser := auth.MemberUsers_MemberUser{}
		mUser.Id = int32(u.Id)
		mUser.Username = u.Username
		mUser.IsOwner = u.IsOwner
		out.Users = append(out.Users, &mUser)
	}

	return nil
}

// ListParentDepartmentUsers RPC服务::列出指定部门的父级部门下所有用户
func (as *AuthService) ListParentDepartmentUsers(ctx context.Context, in *auth.IdQuery, out *auth.MemberUsers) error {
	log.InfoWithFields(log.Fields{"department_id": in.Id}, "srv.ListDepartmentUsers called.")
	defer log.Info("srv.ListDepartmentUsers end.")

	out.Users = make([]*auth.MemberUsers_MemberUser, 0)

	log.Info("Finding department.")
	dept, ok := model.GetDepartment(model.Query{"id=?": string(in.Id)})
	if !ok {
		m := msg.ErrDatabase.GenResponse("Error when AuthService.ListParentDepartmentUsers called, in model.GetDepartment.")
		log.Error(m.String())
		return m.Error()
	}

	if dept == nil || dept.ParentId == nil {
		log.DebugWithFields(log.Fields{"department_id": in.Id}, "Department not found.")
		return nil
	}

	log.Info("Finding department users.")
	us, ok := model.ListDepartmentUsers(*dept.ParentId, model.Query{})
	if !ok {
		m := msg.ErrDatabase.GenResponse("Error when AuthService.ListParentDepartmentUsers called, in model.ListDepartmentUsers.")
		log.Error(m.String())
		return m.Error()
	}

	log.Info("Making response.")
	for _, u := range *us {
		mUser := auth.MemberUsers_MemberUser{}
		mUser.Id = int32(u.Id)
		mUser.Username = u.Username
		mUser.IsOwner = u.IsOwner
		out.Users = append(out.Users, &mUser)
	}

	return nil
}
