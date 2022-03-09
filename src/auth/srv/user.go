package main

import (
	"context"
	"eago/auth/conf"
	"eago/auth/dao"
	"eago/auth/model"
	"eago/auth/srv/proto"
	"eago/common/log"
	"eago/task/conf/msg"
)

// GetUserById RPC服务::按Id查找用户
func (as *AuthService) GetUserById(ctx context.Context, req *auth.IdQuery, rsp *auth.User) error {
	log.Info("srv.GetUserById called.")
	defer log.Info("srv.GetUserById end.")

	log.Info("Finding user.")
	user, ok := dao.GetUser(dao.Query{"id=?": req.Id})
	if !ok {
		m := msg.UnknownError.SetDetail("An error occurred while dao.GetUser.")
		log.ErrorWithFields(m.LogFields())
		return m.RpcError()
	}

	if user == nil {
		log.WarnWithFields(log.Fields{
			"id": req.Id,
		}, "Warring, AuthService.GetUserById got a nil user.")
		return nil
	}

	rsp.Id = int32(user.Id)
	rsp.Username = user.Username
	rsp.Email = user.Email
	rsp.Phone = user.Phone

	return nil
}

// ListUsers RPC服务::列出所有用户
func (as *AuthService) ListUsers(ctx context.Context, req *auth.QueryWithPage, rsp *auth.PagedUsers) error {
	log.Info("srv.ListUsers called.")
	defer log.Info("srv.ListUsers end.")

	query := make(dao.Query)
	for k, v := range req.Query {
		query[k] = v
	}

	pagedData, ok := dao.PagedListUsers(query, int(req.Page), int(req.PageSize))
	if !ok {
		m := msg.UnknownError.SetDetail("An error occurred while dao.PagedListUsers.")
		log.ErrorWithFields(m.LogFields())
		return m.RpcError()
	}

	rsp.Users = make([]*auth.User, 0)
	for _, u := range *pagedData.Data.(*[]model.User) {
		newU := &auth.User{
			Id:       int32(u.Id),
			Username: u.Username,
			Email:    u.Email,
			Phone:    u.Phone,
		}
		rsp.Users = append(rsp.Users, newU)
	}

	rsp.Page = uint32(pagedData.Page)
	rsp.Pages = uint32(pagedData.Pages)
	rsp.PageSize = uint32(pagedData.PageSize)
	rsp.Total = uint32(pagedData.Total)
	return nil
}

// GetUserDepartment RPC服务::获得用户所在部门
func (as *AuthService) GetUserDepartment(ctx context.Context, in *auth.IdQuery, out *auth.UserDepartment) error {
	log.Info("srv.GetUserDepartment called.")
	defer log.Info("srv.GetUserDepartment end.")

	dept, ok := dao.GetUserDepartment(int(in.Id))
	if !ok {
		m := msg.UnknownError.SetDetail("An error occurred while dao.GetUserDepartment.")
		log.ErrorWithFields(m.LogFields())
		return m.RpcError()
	}

	log.Info("Making response.")
	if dept != nil {
		out.Id = int32(dept.Id)
		out.Name = dept.Name
		if dept.ParentId != nil {
			out.ParentId = int32(*dept.ParentId)
		}
		out.IsOwner = dept.IsOwner
		out.JoinedAt = dept.JoinedAt.Format(conf.TIMESTAMP_FORMAT)
	}

	return nil
}

// ListUserDepartmentUsers RPC服务::列出用户所在部门用户
func (as *AuthService) ListUserDepartmentUsers(ctx context.Context, in *auth.IdQuery, out *auth.MemberUsers) error {
	log.Info("srv.ListUserDepartmentUsers called.")
	defer log.Info("srv.ListUserDepartmentUsers end.")

	out.Users = make([]*auth.MemberUsers_MemberUser, 0)

	log.Info("Finding user department.")
	dept, ok := dao.GetUserDepartment(int(in.Id))
	if !ok {
		m := msg.UnknownError.SetDetail("An error occurred while dao.GetUserDepartment.")
		log.ErrorWithFields(m.LogFields())
		return m.RpcError()
	}
	if dept == nil {
		log.DebugWithFields(log.Fields{"department_id": in.Id}, "Department not found.")
		return nil
	}

	log.Info("Finding department users.")
	mem, ok := dao.ListDepartmentUsers(dept.Id, dao.Query{})
	if !ok {
		m := msg.UnknownError.SetDetail("An error occurred while dao.ListDepartmentUsers.")
		log.ErrorWithFields(m.LogFields())
		return m.RpcError()
	}

	log.Info("Making response.")
	for _, u := range *mem {
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

// MakeUserHandover 用户交接
func (as *AuthService) MakeUserHandover(ctx context.Context, in *auth.HandoverRequest, out *auth.BoolMsg) error {
	// 获得交接用户
	user, ok := dao.GetUser(dao.Query{"id=?": in.UserId})
	if !ok {
		m := msg.UnknownError.SetDetail("An error occurred while dao.GetUser.")
		log.ErrorWithFields(log.Fields{
			"user_id": in.UserId,
		}, m.String())
		return m.RpcError()
	}
	// 找不到用户
	if user == nil || user.Id < 1 {
		m := msg.UnknownError.SetDetail("An nil object is returned after calling dao.GetUser.")
		log.ErrorWithFields(log.Fields{
			"user_id": in.UserId,
		}, m.String())
		return m.RpcError()
	}

	// 获得交接目标用户
	tgtUser, ok := dao.GetUser(dao.Query{"id=?": in.TargetUserId})
	if !ok {
		m := msg.UnknownError.SetDetail("An error occurred while dao.GetUser for target user.")
		log.ErrorWithFields(log.Fields{
			"target_user_id": in.TargetUserId,
		}, m.String())
		return m.RpcError()
	}
	// 找不到目标用户
	if tgtUser == nil || tgtUser.Id < 1 {
		m := msg.UnknownError.SetDetail("An nil object is returned after calling dao.GetUser for target user.")
		log.ErrorWithFields(log.Fields{
			"target_user_id": in.TargetUserId,
		}, m.String())
		return m.RpcError()
	}

	// 执行交接
	if err := dao.MakeUserHandover(int(in.UserId), int(in.TargetUserId)); err != nil {
		m := msg.UnknownError.SetError(err, "An error occurred while dao.MakeUserHandover.")
		log.ErrorWithFields(log.Fields{
			"user_id":        in.UserId,
			"target_user_id": in.TargetUserId,
			"error":          err,
		}, m.String())
		return m.RpcError()
	}

	bd := map[string]interface{}{
		// 交接用户
		"from": map[string]interface{}{
			"id":       user.Id,
			"username": user.Username,
			"email":    user.Email,
			"phone":    user.Phone,
		},
		// 交接给
		"to": map[string]interface{}{
			"id":       tgtUser.Id,
			"username": tgtUser.Username,
			"email":    tgtUser.Email,
			"phone":    tgtUser.Phone,
		},
	}

	// 发送消息用户交接消息
	// Topic：eago-auth.topic.user.MakeUserHandover
	if err := Publisher.Publish(ctx, "user", "MakeUserHandover", bd); err != nil {
		log.ErrorWithFields(log.Fields{
			"user_id":         in.UserId,
			"username":        user.Username,
			"target_user_id":  in.TargetUserId,
			"target_username": tgtUser.Username,
			"error":           err,
		}, "Failed when broker.Publisher.Publish.")
	}

	out.Ok = true
	return nil
}
