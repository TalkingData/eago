package srv

import (
	"context"
	"eago-auth/conf/msg"
	db "eago-auth/database"
	"eago-auth/srv/proto"
	"eago-common/log"
	"errors"
)

// GetProducts RPC服务::根据产品线Id获得产品线信息
func (as *AuthService) GetProducts(ctx context.Context, req *auth.Ids, res *auth.Products) error {
	var q = db.Query{"id IN (?)": req.Ids}

	pList, suc := db.ProductModel.List(&q)
	if !suc {
		m := msg.ErrDatabase.NewMsg("Error in AuthService.GetProducts.")
		log.Error(m.String())
		return errors.New(m.String())
	}

	res.Products = make([]*auth.Product, 0)
	for _, p := range *pList {
		newP := auth.Product{}
		newP.Id = int32(p.Id)
		newP.Name = p.Name
		newP.Alias = p.Alias
		newP.Disabled = *p.Disabled
		res.Products = append(res.Products, &newP)
	}

	return nil
}
