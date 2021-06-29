package main

import (
	"context"
	"eago/auth/conf/msg"
	"eago/auth/model"
	"eago/auth/srv/proto"
	"eago/common/log"
	"errors"
)

// GetGroups RPC服务::根据组Id获得组信息
func (as *AuthService) GetGroups(ctx context.Context, req *auth.Ids, rsp *auth.Groups) error {
	gList, ok := model.ListGroups(model.Query{"id IN (?)": req.Ids})
	if !ok {
		m := msg.ErrDatabase.SetDetail("Error in AuthService.GetGroups.")
		log.Error(m.String())
		return errors.New(m.String())
	}

	rsp.Groups = make([]*auth.Group, 0)
	for _, g := range *gList {
		newG := auth.Group{}
		newG.Id = int32(g.Id)
		newG.Name = g.Name
		rsp.Groups = append(rsp.Groups, &newG)
	}

	return nil
}
