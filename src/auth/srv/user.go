package srv

import (
	"context"
	"eago-auth/conf/msg"
	db "eago-auth/database"
	"eago-auth/srv/proto"
	"eago-common/log"
	"errors"
)

// GetUsers RPC服务::根据用户Id获得用户列表
func (as *AuthService) GetUsers(ctx context.Context, req *auth.Ids, res *auth.Users) error {
	var q = db.Query{"id IN (?)": req.Ids}

	uList, suc := db.UserModel.List(&q)
	if !suc {
		m := msg.ErrDatabase.NewMsg("Error in AuthService.GetUsers.")
		log.Error(m.String())
		return errors.New(m.String())
	}

	res.Users = make([]*auth.User, 0)
	for _, u := range *uList {
		newU := auth.User{}
		newU.Id = int32(u.Id)
		newU.Username = u.Username
		newU.Email = u.Email
		newU.Phone = u.Phone
		res.Users = append(res.Users, &newU)
	}

	return nil
}
