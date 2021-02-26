package srv

import (
	"context"
	"eago-auth/config/msg"
	db "eago-auth/database"
	"eago-auth/srv/proto"
	"eago-common/log"
	"errors"
)

// RPC服务::根据组Id获得组信息
func (as *AuthService) GetGroups(ctx context.Context, req *auth.Ids, res *auth.Groups) error {
	var q = db.Query{"id IN (?)": req.Ids}

	gList, suc := db.GroupModel.List(&q)
	if !suc {
		m := msg.ErrDatabase.NewMsg("Error in AuthService.GetGroups.")
		log.Error(m.String())
		return errors.New(m.String())
	}

	res.Groups = make([]*auth.Group, 0)
	for _, g := range *gList {
		newG := auth.Group{}
		newG.Id = int32(g.Id)
		newG.Name = g.Name
		res.Groups = append(res.Groups, &newG)
	}

	return nil
}
