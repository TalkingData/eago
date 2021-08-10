package main

import (
	"context"
	"eago/auth/conf/msg"
	"eago/auth/model"
	"eago/auth/srv/proto"
	"eago/common/log"
)

// ListProducts RPC服务::列出所有产品线
func (as *AuthService) ListProducts(ctx context.Context, req *auth.QueryWithPage, rsp *auth.PagedProducts) error {
	log.Info("srv.ListProducts called.")
	defer log.Info("srv.ListProducts end.")

	query := make(model.Query)
	for k, v := range query {
		query[k] = v
	}

	log.Info("Finding products.")
	pagedData, ok := model.PagedListProducts(query, int(req.Page), int(req.PageSize))
	if !ok {
		m := msg.ErrDatabase.SetDetail("Error when AuthService.ListProducts called, in model.PagedListProducts.")
		log.Error(m.String())
		return m.Error()
	}

	log.Info("Making response.")
	rsp.Products = make([]*auth.Product, 0)
	for _, p := range pagedData.Data.([]model.Product) {
		newP := auth.Product{}
		newP.Id = int32(p.Id)
		newP.Name = p.Name
		newP.Alias = p.Alias
		newP.Disabled = *p.Disabled
		rsp.Products = append(rsp.Products, &newP)
	}

	log.Info("Making page info.")
	rsp.Page = uint32(pagedData.Page)
	rsp.Pages = uint32(pagedData.Pages)
	rsp.PageSize = uint32(pagedData.PageSize)
	rsp.Total = uint32(pagedData.Total)
	return nil
}

// ListProductUsers RPC服务::列出产品线下所有用户
func (as *AuthService) ListProductUsers(ctx context.Context, in *auth.IdQuery, out *auth.MemberUsers) error {
	log.InfoWithFields(log.Fields{"product_id": in.Id}, "srv.ListProductUsers called.")
	defer log.Info("srv.ListProductUsers end.")

	log.Info("Finding product users.")
	us, ok := model.ListProductUsers(int(in.Id), model.Query{})
	if !ok {
		m := msg.ErrDatabase.GenResponse("Error when AuthService.ListProductUsers called, in model.ListProductUsers.")
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
