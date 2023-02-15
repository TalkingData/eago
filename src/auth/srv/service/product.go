package service

import (
	"context"
	"eago/auth/conf/msg"
	"eago/auth/dto"
	"eago/auth/model"
	authpb "eago/auth/proto"
	"eago/common/logger"
	"eago/common/orm"
	commpb "eago/common/proto"
)

// GetProductById 根据ID查询单个产品线
func (authSrv *AuthService) GetProductById(ctx context.Context, req *commpb.IdQuery, rsp *authpb.Product) error {
	authSrv.logger.InfoWithFields(logger.Fields{
		"id": req.Value,
	}, "authSrv.GetProductById called.")
	defer authSrv.logger.Info("authSrv.GetProductById end.")

	authSrv.logger.Info("Finding product.")
	prod, err := authSrv.dao.GetProduct(ctx, orm.Query{"id=?": req.Value})
	if err != nil {
		cMsg := msg.MsgAuthDaoErr.SetError(err)
		authSrv.logger.ErrorWithFields(
			cMsg.ToLoggerFields().Append("id", req.Value),
			"An error occurred while dao.GetProduct in authSrv.GetProductById.",
		)
		return cMsg.ToMicroErr()
	}

	if prod == nil {
		authSrv.logger.WarnWithFields(logger.Fields{
			"id": req.Value,
		}, "Got a nil group object in authSrv.GetProductById.")
		return nil
	}

	rsp.Id = prod.Id
	rsp.Name = prod.Name
	rsp.Alias = prod.Alias
	rsp.Disabled = *prod.Disabled
	return nil
}

// PagedListProducts 分页查询产品线
func (authSrv *AuthService) PagedListProducts(
	ctx context.Context, req *commpb.QueryWithPage, rsp *authpb.PagedProducts,
) error {
	authSrv.logger.Info("authSrv.PagedListProducts called.")
	defer authSrv.logger.Info("authSrv.PagedListProducts end.")

	authSrv.logger.Info("Finding products.")
	pagedData, err := authSrv.dao.PagedListProducts(
		ctx, orm.NewQueryByMapStrStr(req.Query), int(req.Page), int(req.PageSize),
	)
	if err != nil {
		cMsg := msg.MsgAuthDaoErr.SetError(err)
		authSrv.logger.ErrorWithFields(
			cMsg.ToLoggerFields().Append("query", req),
			"An error occurred while dao.PagedListProducts in authSrv.PagedListProducts.",
		)
		return cMsg.ToMicroErr()
	}

	authSrv.logger.Debug("Making response.")
	rsp.Products = make([]*authpb.Product, 0)
	for _, p := range *pagedData.Data.(*[]*model.Product) {
		rsp.Products = append(rsp.Products, &authpb.Product{
			Id:       p.Id,
			Name:     p.Name,
			Alias:    p.Alias,
			Disabled: *p.Disabled,
		})
	}

	authSrv.logger.Debug("Making page info.")
	rsp.Page = uint32(pagedData.Page)
	rsp.Pages = uint32(pagedData.Pages)
	rsp.PageSize = uint32(pagedData.PageSize)
	rsp.Total = uint32(pagedData.Total)
	return nil
}

// ListProductsUsers 列出指定产品线中用户
func (authSrv *AuthService) ListProductsUsers(
	ctx context.Context, req *commpb.IdQuery, rsp *authpb.MemberUsers,
) error {
	authSrv.logger.InfoWithFields(logger.Fields{
		"id": req.Value,
	}, "authSrv.ListProductsUsers called.")
	defer authSrv.logger.Info("authSrv.ListProductsUsers end.")

	authSrv.logger.Info("Finding product users.")
	mUsers, err := authSrv.dao.ListProductsUsers(ctx, req.Value, orm.Query{})
	if err != nil {
		cMsg := msg.MsgAuthDaoErr.SetError(err)
		authSrv.logger.ErrorWithFields(
			cMsg.ToLoggerFields().Append("id", req.Value),
			"An error occurred while dao.ListProductsUsers in authSrv.ListProductsUsers.",
		)
		return cMsg.ToMicroErr()
	}

	dto.CopyMemberUserGrpc(mUsers, rsp)

	return nil
}
