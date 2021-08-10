package main

import (
	"context"
	"eago/auth/conf/msg"
	"eago/auth/model"
	"eago/auth/srv/proto"
	"eago/common/log"
)

// ListGroups RPC服务::列出所有组
func (as *AuthService) ListGroups(ctx context.Context, req *auth.QueryWithPage, rsp *auth.PagedGroups) error {
	log.Info("srv.ListGroups called.")
	defer log.Info("srv.ListGroups end.")

	query := make(model.Query)
	for k, v := range query {
		query[k] = v
	}

	log.Info("Finding groups.")
	pagedData, ok := model.PagedListGroups(query, int(req.Page), int(req.PageSize))
	if !ok {
		m := msg.ErrDatabase.SetDetail("Error when AuthService.ListGroups called, in model.PagedListGroups.")
		log.Error(m.String())
		return m.Error()
	}

	log.Info("Making response.")
	rsp.Groups = make([]*auth.Group, 0)
	for _, g := range pagedData.Data.([]model.Group) {
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
	us, ok := model.ListGroupUsers(int(in.Id), model.Query{})
	if !ok {
		m := msg.ErrDatabase.GenResponse("Error when AuthService.ListGroupUsers called, in model.ListGroupUsers.")
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
