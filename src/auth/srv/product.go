package main

import (
	"context"
	"eago/auth/conf/msg"
	"eago/auth/model"
	"eago/auth/srv/proto"
	"eago/common/log"
	"errors"
)

// GetProducts RPC服务::根据产品线Id获得产品线信息
func (as *AuthService) GetProducts(ctx context.Context, req *auth.Ids, rsp *auth.Products) error {
	pList, ok := model.ListProducts(model.Query{"id IN (?)": req.Ids})
	if !ok {
		m := msg.ErrDatabase.SetDetail("Error in AuthService.GetProducts.")
		log.Error(m.String())
		return errors.New(m.String())
	}

	rsp.Products = make([]*auth.Product, 0)
	for _, p := range *pList {
		newP := auth.Product{}
		newP.Id = int32(p.Id)
		newP.Name = p.Name
		newP.Alias = p.Alias
		newP.Disabled = *p.Disabled
		rsp.Products = append(rsp.Products, &newP)
	}

	return nil
}
