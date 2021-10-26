package main

import (
	"context"
	"eago/auth/conf"
	"eago/auth/conf/msg"
	"eago/auth/dao"
	"eago/auth/srv/proto"
	"eago/common/log"
)

// GetDepartmentById RPC服务::按Id查找部门
func (as *AuthService) GetDepartmentById(ctx context.Context, req *auth.IdQuery, rsp *auth.Department) error {
	log.Info("srv.GetDepartmentById called.")
	defer log.Info("srv.GetDepartmentById end.")

	log.Info("Finding department.")
	dept, ok := dao.GetDepartment(dao.Query{"id=?": req.Id})
	if !ok {
		m := msg.UnknownError.SetDetail("An error occurred while dao.GetDepartment.")
		log.ErrorWithFields(m.LogFields())
		return m.RpcError()
	}

	if dept == nil {
		log.WarnWithFields(log.Fields{
			"id": req.Id,
		}, "Warring, AuthService.GetDepartmentById got a nil department.")
		return nil
	}

	rsp.Id = int32(dept.Id)
	rsp.Name = dept.Name
	if dept.ParentId != nil {
		rsp.ParentId = int32(*dept.ParentId)
	}
	return nil
}

// ListDepartmentUsers RPC服务::列出指定部门下所有用户
func (as *AuthService) ListDepartmentUsers(ctx context.Context, in *auth.IdQuery, out *auth.MemberUsers) error {
	log.InfoWithFields(log.Fields{"department_id": in.Id}, "srv.ListDepartmentUsers called.")
	defer log.Info("srv.ListDepartmentUsers end.")

	log.Info("Finding department users.")
	us, ok := dao.ListDepartmentUsers(int(in.Id), dao.Query{})
	if !ok {
		m := msg.UnknownError.SetDetail("An error occurred while dao.ListDepartmentUsers.")
		log.ErrorWithFields(m.LogFields())
		return m.RpcError()
	}

	log.Info("Making response.")
	out.Users = make([]*auth.MemberUsers_MemberUser, 0)
	for _, u := range *us {
		mUser := &auth.MemberUsers_MemberUser{
			Id:       int32(u.Id),
			Username: u.Username,
			IsOwner:  u.IsOwner,
			JoinedAt: u.JoinedAt.Format(conf.TIMESTAMP_FORMAT),
		}
		out.Users = append(out.Users, mUser)
	}

	return nil
}

// ListParentDepartmentUsers RPC服务::列出指定部门的父级部门下所有用户
func (as *AuthService) ListParentDepartmentUsers(ctx context.Context, in *auth.IdQuery, out *auth.MemberUsers) error {
	log.InfoWithFields(log.Fields{"department_id": in.Id}, "srv.ListDepartmentUsers called.")
	defer log.Info("srv.ListDepartmentUsers end.")

	out.Users = make([]*auth.MemberUsers_MemberUser, 0)

	log.Info("Finding department.")
	dept, ok := dao.GetDepartment(dao.Query{"id=?": in.Id})
	if !ok {
		m := msg.UnknownError.SetDetail("An error occurred while dao.GetDepartment.")
		log.ErrorWithFields(m.LogFields())
		return m.RpcError()
	}

	// 找不到当前部门
	if dept == nil {
		log.WarnWithFields(log.Fields{
			"id": in.Id,
		}, "Warring, AuthService.ListParentDepartmentUsers got a nil department.")
		return nil
	}

	// 当前部门没有父部门
	if dept.ParentId == nil {
		log.DebugWithFields(log.Fields{"department_id": in.Id}, "Department has no parent department.")
		return nil
	}

	log.Info("Finding department users.")
	us, ok := dao.ListDepartmentUsers(*dept.ParentId, dao.Query{})
	if !ok {
		m := msg.UnknownError.SetDetail("An error occurred while dao.ListDepartmentUsers.")
		log.ErrorWithFields(m.LogFields())
		return m.RpcError()
	}

	log.Info("Making response.")
	for _, u := range *us {
		mUser := &auth.MemberUsers_MemberUser{
			Id:       int32(u.Id),
			Username: u.Username,
			IsOwner:  u.IsOwner,
			JoinedAt: u.JoinedAt.Format(conf.TIMESTAMP_FORMAT),
		}
		out.Users = append(out.Users, mUser)
	}

	return nil
}
