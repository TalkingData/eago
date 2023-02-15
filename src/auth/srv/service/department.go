package service

import (
	"context"
	"eago/auth/conf/msg"
	"eago/auth/dto"
	authpb "eago/auth/proto"
	"eago/common/logger"
	"eago/common/orm"
	commpb "eago/common/proto"
)

// GetDepartmentById 根据ID查询单个部门
func (authSrv *AuthService) GetDepartmentById(ctx context.Context, req *commpb.IdQuery, rsp *authpb.Department) error {
	authSrv.logger.InfoWithFields(logger.Fields{
		"id": req.Value,
	}, "authSrv.GetDepartmentById called.")
	defer authSrv.logger.Info("authSrv.GetDepartmentById end.")

	dept, err := authSrv.dao.GetDepartment(ctx, orm.Query{"id=?": req.Value})
	if err != nil {
		cMsg := msg.MsgAuthDaoErr.SetError(err)
		authSrv.logger.ErrorWithFields(
			cMsg.ToLoggerFields().Append("id", req.Value),
			"An error occurred while dao.GetDepartment in authSrv.GetDepartmentById.",
		)
		return cMsg.ToMicroErr()
	}

	if dept == nil {
		authSrv.logger.WarnWithFields(logger.Fields{
			"id": req.Value,
		}, "Got a nil department object in authSrv.GetDepartmentById.")
		return nil
	}

	rsp.Id = dept.Id
	rsp.Name = dept.Name
	if dept.ParentId != nil {
		rsp.ParentId = *dept.ParentId
	}
	return nil
}

// ListDepartmentsUsers 列出指定部门中用户
func (authSrv *AuthService) ListDepartmentsUsers(
	ctx context.Context, req *commpb.IdQuery, rsp *authpb.MemberUsers,
) error {
	authSrv.logger.InfoWithFields(logger.Fields{
		"id": req.Value,
	}, "authSrv.ListDepartmentsUsers called.")
	defer authSrv.logger.Info("authSrv.ListDepartmentsUsers end.")

	mUsers, err := authSrv.dao.ListDepartmentsUsers(ctx, req.Value, orm.Query{})
	if err != nil {
		cMsg := msg.MsgAuthDaoErr.SetError(err)
		authSrv.logger.ErrorWithFields(
			cMsg.ToLoggerFields().Append("id", req.Value),
			"An error occurred while dao.ListDepartmentsUsers in authSrv.ListDepartmentsUsers.",
		)
		return cMsg.ToMicroErr()
	}

	dto.CopyMemberUserGrpc(mUsers, rsp)

	return nil
}

// ListParentDepartmentUsers 列出指定部门的父部门中用户
func (authSrv *AuthService) ListParentDepartmentUsers(
	ctx context.Context, req *commpb.IdQuery, rsp *authpb.MemberUsers,
) error {
	authSrv.logger.InfoWithFields(logger.Fields{
		"id": req.Value,
	}, "authSrv.ListDepartmentsUsers called.")
	defer authSrv.logger.Info("authSrv.ListDepartmentsUsers end.")

	authSrv.logger.Info("Finding department.")
	dept, err := authSrv.dao.GetDepartment(ctx, orm.Query{"id=?": req.Value})
	if err != nil {
		cMsg := msg.MsgAuthDaoErr.SetError(err)
		authSrv.logger.ErrorWithFields(
			cMsg.ToLoggerFields().Append("id", req.Value),
			"",
		)
		return cMsg.ToMicroErr()
	}

	// 找不到当前部门
	if dept == nil {
		authSrv.logger.WarnWithFields(logger.Fields{
			"id": req.Value,
		}, "Got a nil department object dao.GetDepartment in authSrv.ListParentDepartmentUsers.")
		return nil
	}

	// 当前部门没有父部门
	if dept.ParentId == nil {
		authSrv.logger.WarnWithFields(logger.Fields{
			"department_id": req.Value,
		}, "Department has no parent department.")
		return nil
	}

	authSrv.logger.Info("Finding department users.")
	mUsers, err := authSrv.dao.ListDepartmentsUsers(ctx, *dept.ParentId, orm.Query{})
	if err != nil {
		cMsg := msg.MsgAuthDaoErr.SetError(err)
		authSrv.logger.ErrorWithFields(
			cMsg.ToLoggerFields().Append("id", req.Value),
			"",
		)
		return cMsg.ToMicroErr()
	}

	dto.CopyMemberUserGrpc(mUsers, rsp)

	return nil
}
