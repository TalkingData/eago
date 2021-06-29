package main

import (
	"context"
	"eago/auth/srv/local"
	"eago/auth/srv/proto"
	"eago/common/log"
	"sync"
)

// VerifyToken RPC服务::验证Token是否有效
func (as *AuthService) VerifyToken(ctx context.Context, req *auth.Token, res *auth.BoolMsg) error {
	log.InfoWithFields(log.Fields{"token": req.Token}, "Got rpc call verify token.")
	res.Ok = local.VerifyToken(req.Token)
	return nil
}

// GetTokenContent RPC服务::通过Token获得TokenContent
func (as *AuthService) GetTokenContent(ctx context.Context, req *auth.Token, rsp *auth.TokenContent) error {
	log.InfoWithFields(
		log.Fields{"token": req.Token},
		"Got rpc call get token content.",
	)
	tc, ok := local.GetTokenContent(req.Token)
	if !ok || tc == nil {
		rsp.Ok = false
		return nil
	}

	rsp.Ok = true
	rsp.UserId = int32(tc.UserId)
	rsp.Username = tc.Username
	rsp.Phone = tc.Phone

	rsp.IsSuperuser = tc.IsSuperuser

	// 加载部门
	if tc.Department != nil {
		rsp.Department = *tc.Department
	}
	// 加载角色
	rsp.Roles = *tc.Roles

	// 加载产品线
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		rsp.Products = make([]*auth.Product, 0)
		rsp.OwnProducts = make([]*auth.Product, 0)

		for _, p := range *tc.Products {
			newP := auth.Product{}
			newP.Id = int32(p.Id)
			newP.Name = p.Name
			newP.Alias = p.Alias
			newP.Disabled = p.Disabled
			rsp.Products = append(rsp.Products, &newP)
		}

		for _, p := range *tc.OwnProducts {
			newP := auth.Product{}
			newP.Id = int32(p.Id)
			newP.Name = p.Name
			newP.Alias = p.Alias
			newP.Disabled = p.Disabled
			rsp.OwnProducts = append(rsp.OwnProducts, &newP)
		}

	}(wg)

	// 加载组
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		rsp.Groups = make([]*auth.Group, 0)
		rsp.OwnGroups = make([]*auth.Group, 0)

		for _, g := range *tc.Groups {
			newG := auth.Group{}
			newG.Id = int32(g.Id)
			newG.Name = g.Name
			rsp.Groups = append(rsp.Groups, &newG)
		}

		for _, g := range *tc.OwnGroups {
			newG := auth.Group{}
			newG.Id = int32(g.Id)
			newG.Name = g.Name
			rsp.OwnGroups = append(rsp.OwnGroups, &newG)
		}

	}(wg)

	// 等待所有加载项结束
	wg.Wait()

	return nil
}
