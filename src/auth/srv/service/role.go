package service

import (
	"context"
	"eago/auth/conf/msg"
	authpb "eago/auth/proto"
	"eago/common/logger"
	"eago/common/orm"
	commpb "eago/common/proto"
)

func (authSrv *AuthService) ListRolesUsers(
	ctx context.Context, req *commpb.NameQuery, rsp *authpb.RolesMemberUsers,
) error {
	authSrv.logger.Info(logger.Fields{
		"name": req.Value,
	}, "authSrv.ListRolesUsers called.")
	defer authSrv.logger.Info("authSrv.ListRolesUsers end.")

	authSrv.logger.Info("Finding role.")
	r, err := authSrv.dao.GetRole(ctx, orm.Query{"name=?": req.Value})
	if err != nil {
		cMsg := msg.MsgAuthDaoErr.SetError(err)
		authSrv.logger.ErrorWithFields(
			cMsg.ToLoggerFields().Append("name", req.Value),
			"An error occurred while dao.GetRole in authSrv.ListRolesUsers.",
		)
		return cMsg.ToMicroErr()
	}

	rsp.Users = make([]*authpb.RolesMemberUsers_MemberUser, 0)

	// 查找不到角色时直接返回空结果
	if r == nil {
		return nil
	}

	authSrv.logger.Info("Finding role users.")
	us, err := authSrv.dao.ListRolesUsers(ctx, r.Id)
	if err != nil {
		cMsg := msg.MsgAuthDaoErr.SetError(err)
		authSrv.logger.ErrorWithFields(
			cMsg.ToLoggerFields().Append("id", r.Id),
			"An error occurred while dao.ListRolesUsers in authSrv.ListRolesUsers.",
		)
		return cMsg.ToMicroErr()
	}

	authSrv.logger.Info("Making response.")
	for _, u := range us {
		rsp.Users = append(rsp.Users, &authpb.RolesMemberUsers_MemberUser{
			Id:       u.Id,
			Username: u.Username,
		})
	}

	return nil
}
