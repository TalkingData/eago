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

// GetUserById 根据ID查询单个用户
func (authSrv *AuthService) GetUserById(ctx context.Context, req *commpb.IdQuery, rsp *authpb.User) error {
	authSrv.logger.InfoWithFields(logger.Fields{
		"id": req.Value,
	}, "authSrv.GetUserById called.")
	defer authSrv.logger.Info("authSrv.GetUserById end.")

	user, err := authSrv.dao.GetUser(ctx, orm.Query{"id=?": req.Value})
	if err != nil {
		cMsg := msg.MsgAuthDaoErr.SetError(err)
		authSrv.logger.ErrorWithFields(
			cMsg.ToLoggerFields().Append("id", req.Value),
			"An error occurred while dao.GetUser in authSrv.GetUserById.",
		)
		return cMsg.ToMicroErr()
	}

	if user == nil {
		authSrv.logger.WarnWithFields(logger.Fields{
			"id": req.Value,
		}, "Got a nil user object dao.GetUser in authSrv.GetUserById.")
		return nil
	}

	rsp.Id = user.Id
	rsp.Username = user.Username
	rsp.Email = user.Email
	rsp.Phone = user.Phone

	return nil
}

// PagedListUsers 分页查询用户
func (authSrv *AuthService) PagedListUsers(
	ctx context.Context, req *commpb.QueryWithPage, rsp *authpb.PagedUsers,
) error {
	authSrv.logger.Info("authSrv.PagedListUsers called.")
	defer authSrv.logger.Info("authSrv.PagedListUsers end.")

	pagedData, err := authSrv.dao.PagedListUsers(
		ctx, orm.NewQueryByMapStrStr(req.Query), int(req.Page), int(req.PageSize),
	)
	if err != nil {
		cMsg := msg.MsgAuthDaoErr.SetError(err)
		authSrv.logger.ErrorWithFields(
			cMsg.ToLoggerFields().Append("query", req),
			"An error occurred while dao.PagedListUsers in authSrv.PagedListUsers.",
		)
		return cMsg.ToMicroErr()
	}

	rsp.Users = make([]*authpb.User, 0)
	for _, u := range *pagedData.Data.(*[]*model.User) {
		rsp.Users = append(rsp.Users, &authpb.User{
			Id:       u.Id,
			Username: u.Username,
			Email:    u.Email,
			Phone:    u.Phone,
		})
	}

	rsp.Page = uint32(pagedData.Page)
	rsp.Pages = uint32(pagedData.Pages)
	rsp.PageSize = uint32(pagedData.PageSize)
	rsp.Total = uint32(pagedData.Total)
	return nil
}

// GetUserDepartment 获得指定用户的部门
func (authSrv *AuthService) GetUsersDepartment(
	ctx context.Context, req *commpb.IdQuery, rsp *authpb.UsersDepartment,
) error {
	authSrv.logger.InfoWithFields(logger.Fields{
		"id": req.Value,
	}, "authSrv.GetUsersDepartment called.")
	defer authSrv.logger.Info("authSrv.GetUsersDepartment end.")

	dept, err := authSrv.dao.GetUsersDepartment(ctx, req.Value)
	if err != nil {
		cMsg := msg.MsgAuthDaoErr.SetError(err)
		authSrv.logger.ErrorWithFields(
			cMsg.ToLoggerFields().Append("id", req.Value),
			"An error occurred while dao.GetUsersDepartment in authSrv.GetUsersDepartment.",
		)
		return cMsg.ToMicroErr()
	}

	if dept != nil {
		rsp.Id = dept.Id
		rsp.Name = dept.Name
		if dept.ParentId != nil {
			rsp.ParentId = *dept.ParentId
		}
		rsp.IsOwner = dept.IsOwner
		if dept.JoinedAt != nil {
			rsp.JoinedAt = dept.JoinedAt.String()
		}
	}

	return nil
}

// ListUsersSameDepartmentUsers 列出与指定用户相同部门的所有用户
func (authSrv *AuthService) ListUsersSameDepartmentUsers(
	ctx context.Context, req *commpb.IdQuery, rsp *authpb.MemberUsers,
) error {
	authSrv.logger.InfoWithFields(logger.Fields{
		"id": req.Value,
	}, "authSrv.ListUsersDepartmentUsers called.")
	defer authSrv.logger.Info("authSrv.ListUsersDepartmentUsers end.")

	rsp.Users = make([]*authpb.MemberUsers_MemberUser, 0)

	authSrv.logger.Info("Finding user department.")
	dept, err := authSrv.dao.GetUsersDepartment(ctx, req.Value)
	if err != nil {
		cMsg := msg.MsgAuthDaoErr.SetError(err)
		authSrv.logger.ErrorWithFields(
			cMsg.ToLoggerFields().Append("id", req.Value),
			"An error occurred while dao.GetUsersDepartment in authSrv.ListUsersDepartmentUsers.",
		)
		return cMsg.ToMicroErr()
	}
	if dept == nil {
		authSrv.logger.DebugWithFields(logger.Fields{"department_id": req.Value}, "Department not found.")
		return nil
	}

	authSrv.logger.Info("Finding department users.")
	mUsers, err := authSrv.dao.ListDepartmentsUsers(ctx, dept.Id, orm.Query{})
	if err != nil {
		cMsg := msg.MsgAuthDaoErr.SetError(err)
		authSrv.logger.ErrorWithFields(
			cMsg.ToLoggerFields().Append("id", req.Value),
			"An error occurred while dao.ListDepartmentsUsers in authSrv.ListUsersDepartmentUsers.",
		)
		return cMsg.ToMicroErr()
	}

	dto.CopyMemberUserGrpc(mUsers, rsp)

	return nil
}
