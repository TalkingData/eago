package main

import (
	"context"
	"eago/auth/conf/msg"
	"eago/auth/dao"
	"eago/auth/model"
	"eago/auth/srv/proto"
	"eago/common/log"
)

// GetGroupById RPC服务::按Id查找组
func (as *AuthService) GetGroupById(ctx context.Context, req *auth.IdQuery, rsp *auth.Group) error {
	log.Info("srv.GetGroupById called.")
	defer log.Info("srv.GetGroupById end.")

	log.Info("Finding user.")
	g, ok := dao.GetGroup(dao.Query{"id=?": req.Id})
	if !ok {
		m := msg.UndefinedError.SetDetail("An error occurred while dao.GetGroup.")
		log.ErrorWithFields(m.LogFields())
		return m.RpcError()
	}

	if g == nil {
		log.WarnWithFields(log.Fields{
			"id": req.Id,
		}, "Warring, AuthService.GetGroupById got a nil group.")
		return nil
	}

	rsp.Id = int32(g.Id)
	rsp.Name = g.Name
	return nil
}

// PagedListGroups RPC服务::列出所有组-分页
func (as *AuthService) PagedListGroups(ctx context.Context, req *auth.QueryWithPage, rsp *auth.PagedGroups) error {
	log.Info("srv.PagedListGroups called.")
	defer log.Info("srv.PagedListGroups end.")

	query := make(dao.Query)
	for k, v := range req.Query {
		query[k] = v
	}

	log.Info("Finding groups.")
	pagedData, ok := dao.PagedListGroups(query, int(req.Page), int(req.PageSize))
	if !ok {
		m := msg.UndefinedError.SetDetail("An error occurred while dao.PagedListGroups.")
		log.ErrorWithFields(m.LogFields())
		return m.RpcError()
	}

	log.Info("Making response.")
	rsp.Groups = make([]*auth.Group, 0)
	for _, g := range *pagedData.Data.(*[]model.Group) {
		newG := auth.Group{}
		newG.Id = int32(g.Id)
		newG.Name = g.Name
		rsp.Groups = append(rsp.Groups, &newG)
	}

	log.Info("Making page info.")
	rsp.Page = uint32(pagedData.Page)
	rsp.Pages = uint32(pagedData.Pages)
	rsp.PageSize = uint32(pagedData.PageSize)
	rsp.Total = uint32(pagedData.Total)
	return nil
}

// ListGroupUsers RPC服务::列出组中所有用户
func (as *AuthService) ListGroupUsers(ctx context.Context, in *auth.IdQuery, out *auth.MemberUsers) error {
	log.InfoWithFields(log.Fields{"group_id": in.Id}, "srv.ListGroupUsers called.")
	defer log.Info("srv.ListGroupUsers end.")

	log.Info("Finding group users.")
	us, ok := dao.ListGroupUsers(int(in.Id), dao.Query{})
	if !ok {
		m := msg.UndefinedError.SetDetail("An error occurred while dao.ListGroupUsers.")
		log.ErrorWithFields(m.LogFields())
		return m.RpcError()
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
