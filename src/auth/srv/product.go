package main

import (
	"context"
	"eago/auth/conf"
	"eago/auth/conf/msg"
	"eago/auth/dao"
	"eago/auth/model"
	"eago/auth/srv/proto"
	"eago/common/log"
)

// GetProductById RPC服务::按Id查找产品线
func (as *AuthService) GetProductById(ctx context.Context, req *auth.IdQuery, rsp *auth.Product) error {
	log.Info("srv.GetProductById called.")
	defer log.Info("srv.GetProductById end.")

	log.Info("Finding product.")
	prod, ok := dao.GetProduct(dao.Query{"id=?": req.Id})
	if !ok {
		m := msg.UndefinedError.SetDetail("An error occurred while dao.GetProduct.")
		log.ErrorWithFields(m.LogFields())
		return m.RpcError()
	}

	if prod == nil {
		log.WarnWithFields(log.Fields{
			"id": req.Id,
		}, "Warring, AuthService.GetProductById got a nil product.")
		return nil
	}

	rsp.Id = int32(prod.Id)
	rsp.Name = prod.Name
	rsp.Alias = prod.Alias
	rsp.Disabled = *prod.Disabled
	return nil
}

// ListProducts RPC服务::列出所有产品线
func (as *AuthService) ListProducts(ctx context.Context, req *auth.QueryWithPage, rsp *auth.PagedProducts) error {
	log.Info("srv.ListProducts called.")
	defer log.Info("srv.ListProducts end.")

	query := make(dao.Query)
	for k, v := range req.Query {
		query[k] = v
	}

	log.Info("Finding products.")
	pagedData, ok := dao.PagedListProducts(query, int(req.Page), int(req.PageSize))
	if !ok {
		m := msg.UndefinedError.SetDetail("An error occurred while dao.PagedListProducts.")
		log.ErrorWithFields(m.LogFields())
		return m.RpcError()
	}

	log.Info("Making response.")
	rsp.Products = make([]*auth.Product, 0)
	for _, p := range *pagedData.Data.(*[]model.Product) {
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
	us, ok := dao.ListProductUsers(int(in.Id), dao.Query{})
	if !ok {
		m := msg.UndefinedError.SetDetail("An error occurred while dao.ListProductUsers.")
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
