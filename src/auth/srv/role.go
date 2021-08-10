package main

import (
	"context"
	"eago/auth/conf/msg"
	"eago/auth/model"
	"eago/auth/srv/proto"
	"eago/common/log"
)

// ListRoleUsers RPC服务::列出角色中所有用户
func (as *AuthService) ListRoleUsers(ctx context.Context, in *auth.NameQuery, out *auth.RoleMemberUsers) error {
	log.Info("srv.ListRoleUsers called.")
	defer log.Info("srv.ListRoleUsers end.")

	log.Info("Finding role.")
	r, ok := model.GetRole(model.Query{"name=?": in.Name})
	if !ok {
		m := msg.ErrDatabase.GenResponse("Error when AuthService.ListRoleUsers called, in model.GetRole.")
		log.Error(m.String())
		return m.Error()
	}

	log.Info("Finding role users.")
	us, ok := model.ListRoleUsers(r.Id)
	if !ok {
		m := msg.ErrDatabase.GenResponse("Error when AuthService.ListRoleUsers called, in model.ListRoleUsers.")
		log.Error(m.String())
		return m.Error()
	}

	log.Info("Making response.")
	out.Users = make([]*auth.RoleMemberUsers_MemberUser, 0)
	for _, u := range *us {
		mUser := auth.RoleMemberUsers_MemberUser{}
		mUser.Id = int32(u.Id)
		mUser.Username = u.Username
		out.Users = append(out.Users, &mUser)
	}

	return nil
}
