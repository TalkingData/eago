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

// GetGroupById 根据ID查询单个组
func (authSrv *AuthService) GetGroupById(ctx context.Context, req *commpb.IdQuery, rsp *authpb.Group) error {
	authSrv.logger.InfoWithFields(logger.Fields{
		"id": req.Value,
	}, "authSrv.GetGroupById called.")
	defer authSrv.logger.Info("authSrv.GetGroupById end.")

	g, err := authSrv.dao.GetGroup(ctx, orm.Query{"id=?": req.Value})
	if err != nil {
		cMsg := msg.MsgAuthDaoErr.SetError(err)
		authSrv.logger.ErrorWithFields(
			cMsg.ToLoggerFields().Append("id", req.Value),
			"An error occurred while dao.GetGroup in authSrv.GetGroupById.",
		)
		return cMsg.ToMicroErr()
	}

	if g == nil {
		authSrv.logger.WarnWithFields(logger.Fields{
			"id": req.Value,
		}, "Got a nil group object in authSrv.GetGroupById.")
		return nil
	}

	rsp.Id = g.Id
	rsp.Name = g.Name
	return nil
}

// PagedListGroups 分页查询组
func (authSrv *AuthService) PagedListGroups(
	ctx context.Context, req *commpb.QueryWithPage, rsp *authpb.PagedGroups,
) error {
	authSrv.logger.Info("authSrv.PagedListGroups called.")
	defer authSrv.logger.Info("authSrv.PagedListGroups end.")

	pagedData, err := authSrv.dao.PagedListGroups(
		ctx, orm.NewQueryByMapStrStr(req.Query), int(req.Page), int(req.PageSize),
	)
	if err != nil {
		cMsg := msg.MsgAuthDaoErr.SetError(err)
		authSrv.logger.ErrorWithFields(
			cMsg.ToLoggerFields().Append("query", req),
			"An error occurred while dao.PagedListGroups in authSrv.PagedListGroups.",
		)
		return cMsg.ToMicroErr()
	}

	rsp.Groups = make([]*authpb.Group, 0)
	for _, g := range *pagedData.Data.(*[]*model.Group) {
		rsp.Groups = append(rsp.Groups, &authpb.Group{
			Id:   g.Id,
			Name: g.Name,
		})
	}

	rsp.Page = uint32(pagedData.Page)
	rsp.Pages = uint32(pagedData.Pages)
	rsp.PageSize = uint32(pagedData.PageSize)
	rsp.Total = uint32(pagedData.Total)
	return nil
}

// ListGroupsUsers 列出指定组中用户
func (authSrv *AuthService) ListGroupsUsers(ctx context.Context, req *commpb.IdQuery, rsp *authpb.MemberUsers) error {
	authSrv.logger.InfoWithFields(logger.Fields{
		"id": req.Value,
	}, "authSrv.ListGroupsUsers called.")
	defer authSrv.logger.Info("authSrv.ListGroupsUsers end.")

	mUsers, err := authSrv.dao.ListGroupsUsers(ctx, req.Value, orm.Query{})
	if err != nil {
		cMsg := msg.MsgAuthDaoErr.SetError(err)
		authSrv.logger.ErrorWithFields(
			cMsg.ToLoggerFields().Append("id", req.Value),
			"An error occurred while dao.ListGroupsUsers in authSrv.ListGroupsUsers.",
		)
		return cMsg.ToMicroErr()
	}

	dto.CopyMemberUserGrpc(mUsers, rsp)

	return nil
}
