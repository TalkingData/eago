package dto

import (
	authpb "eago/auth/proto"
	"sync"
)

type ProductInToken struct {
	Id       uint32 `json:"id"`
	Name     string `json:"name" binding:"required"`
	Alias    string `json:"alias" binding:"required"`
	Disabled bool   `json:"disabled"`
}

type GroupInToken struct {
	Id   uint32 `json:"id"`
	Name string `json:"name" binding:"required"`
}

// TokenContent Token详细内容
type TokenContent struct {
	UserId   uint32 `json:"user_id"`
	Username string `json:"username"`
	Phone    string `json:"phone"`

	IsSuperuser bool `json:"is_superuser"`

	Department  []*string         `json:"department"`
	Roles       []*string         `json:"roles"`
	Products    []*ProductInToken `json:"products"`
	OwnProducts []*ProductInToken `json:"own_products"`
	Groups      []*GroupInToken   `json:"groups"`
	OwnGroups   []*GroupInToken   `json:"own_groups"`
}

func (tc *TokenContent) Trans2AuthPb(in *authpb.TokenContent) {
	in.UserId = tc.UserId
	in.Username = tc.Username
	in.Phone = tc.Phone

	in.IsSuperuser = tc.IsSuperuser

	// 加载部门
	in.Department = make([]string, len(tc.Department))
	for idx, dept := range tc.Department {
		in.Department[idx] = *dept
	}

	// 加载角色
	in.Roles = make([]string, len(tc.Roles))
	for idx, role := range tc.Roles {
		in.Roles[idx] = *role
	}

	// 加载产品线
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		in.Products = make([]*authpb.Product, len(tc.Products))
		for idx, p := range tc.Products {
			in.Products[idx] = p.Trans2AuthPb()
		}

		in.OwnProducts = make([]*authpb.Product, len(tc.OwnProducts))
		for idx, op := range tc.OwnProducts {
			in.OwnProducts[idx] = op.Trans2AuthPb()
		}
	}(wg)

	// 加载组
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		in.Groups = make([]*authpb.Group, len(tc.Groups))
		for idx, g := range tc.Groups {
			in.Groups[idx] = g.Trans2AuthPb()
		}

		in.OwnGroups = make([]*authpb.Group, len(tc.OwnGroups))
		for idx, og := range tc.OwnGroups {
			in.OwnGroups[idx] = og.Trans2AuthPb()
		}
	}(wg)

	// 等待所有加载项结束
	wg.Wait()
}

func (pit *ProductInToken) Trans2AuthPb() *authpb.Product {
	return &authpb.Product{
		Id:       pit.Id,
		Name:     pit.Name,
		Alias:    pit.Alias,
		Disabled: pit.Disabled,
	}
}

func (git *GroupInToken) Trans2AuthPb() *authpb.Group {
	return &authpb.Group{
		Id:   git.Id,
		Name: git.Name,
	}
}
