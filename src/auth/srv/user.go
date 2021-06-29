package main

import (
	"context"
	"eago/auth/conf/msg"
	"eago/auth/model"
	"eago/auth/srv/proto"
	"eago/common/log"
	"errors"
)

// GetUsers RPC服务::根据用户Id获得用户列表
func (as *AuthService) GetUsers(ctx context.Context, req *auth.Ids, rsp *auth.Users) error {
	uList, ok := model.ListUsers(model.Query{"id IN (?)": req.Ids})
	if !ok {
		m := msg.ErrDatabase.SetDetail("Error in AuthService.GetUsers.")
		log.Error(m.String())
		return errors.New(m.String())
	}

	rsp.Users = make([]*auth.User, 0)
	for _, u := range *uList {
		newU := auth.User{}
		newU.Id = int32(u.Id)
		newU.Username = u.Username
		newU.Email = u.Email
		newU.Phone = u.Phone
		rsp.Users = append(rsp.Users, &newU)
	}

	return nil
}
