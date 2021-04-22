package srv

import (
	"context"
	"eago-auth/conf"
	db "eago-auth/database"
	"eago-auth/srv/proto"
	"eago-common/log"
	"eago-common/redis"
	"eago-common/tools"
	"encoding/json"
	"sync"
	"time"
)

type ProductInToken struct {
	Id       int    `json:"id"`
	Name     string `json:"name" binding:"required"`
	Alias    string `json:"alias" binding:"required"`
	Disabled bool   `json:"disabled"`
}

type GroupInToken struct {
	Id   int    `json:"id"`
	Name string `json:"name" binding:"required"`
}

type TokenContent struct {
	UserId   int    `json:"user_id"`
	Username string `json:"username"`

	IsSuperuser bool              `json:"is_superuser"`
	Roles       *[]string         `json:"roles"`
	Products    *[]ProductInToken `json:"products"`
	OwnProducts *[]ProductInToken `json:"own_products"`
	Groups      *[]GroupInToken   `json:"groups"`
	OwnGroups   *[]GroupInToken   `json:"own_groups"`
}

// VerifyToken RPC服务::验证Token是否有效
func (as *AuthService) VerifyToken(ctx context.Context, req *auth.Token, res *auth.BoolResponse) error {
	log.InfoWithFields(log.Fields{"token": req.Token}, "Gor rpc call verify token.")
	res.Ok = VerifyToken(req.Token)
	return nil
}

// GetTokenContent RPC服务::通过Token获得TokenContent
func (as *AuthService) GetTokenContent(ctx context.Context, req *auth.Token, res *auth.TokenContent) error {
	log.InfoWithFields(
		log.Fields{"token": req.Token},
		"Got rpc call get token content.",
	)
	tc, suc := GetTokenContent(req.Token)
	if !suc || tc == nil {
		res.Ok = false
		return nil
	}

	res.Ok = true
	res.UserId = int32(tc.UserId)
	res.Username = tc.Username

	res.IsSuperuser = tc.IsSuperuser

	// 加载角色
	res.Roles = *tc.Roles

	// 加载产品线
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		res.Products = make([]*auth.Product, 0)
		res.OwnProducts = make([]*auth.Product, 0)

		for _, p := range *tc.Products {
			newP := auth.Product{}
			newP.Id = int32(p.Id)
			newP.Name = p.Name
			newP.Alias = p.Alias
			newP.Disabled = p.Disabled
			res.Products = append(res.Products, &newP)
		}

		for _, p := range *tc.OwnProducts {
			newP := auth.Product{}
			newP.Id = int32(p.Id)
			newP.Name = p.Name
			newP.Alias = p.Alias
			newP.Disabled = p.Disabled
			res.OwnProducts = append(res.OwnProducts, &newP)
		}

	}(&wg)

	// 加载组
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		res.Groups = make([]*auth.Group, 0)
		res.OwnGroups = make([]*auth.Group, 0)

		for _, g := range *tc.Groups {
			newG := auth.Group{}
			newG.Id = int32(g.Id)
			newG.Name = g.Name
			res.Groups = append(res.Groups, &newG)
		}

		for _, g := range *tc.OwnGroups {
			newG := auth.Group{}
			newG.Id = int32(g.Id)
			newG.Name = g.Name
			res.OwnGroups = append(res.OwnGroups, &newG)
		}

	}(&wg)

	// 等待所有加载项结束
	wg.Wait()

	return nil
}

// NewToken 本地服务::生成Token
func NewToken(userObj *db.User) string {
	var (
		tc       = TokenContent{}
		currTime = time.Now().Format(conf.TIMESTAMP_FORMAT)
	)

	tc.UserId = userObj.Id
	tc.Username = userObj.Username

	tc.IsSuperuser = userObj.IsSuperuser

	// 填入角色信息
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		rolesStr := make([]string, 0)
		roles, suc := db.UserModel.ListRoles(userObj.Id)
		if !suc {
			log.Error("Can not load roles, Error in new token db.UserModel.ListProducts.")
		} else {
			for _, r := range *roles {
				rolesStr = append(rolesStr, r.Name)
			}
		}
		tc.Roles = &rolesStr

	}(&wg)

	// 填入产品线信息
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		products := make([]ProductInToken, 0)
		ownProducts := make([]ProductInToken, 0)

		prods, suc := db.UserModel.ListProducts(userObj.Id)
		if !suc {
			log.Error("Can not load products, Error in new token db.UserModel.ListProducts.")
		} else {
			for _, p := range *prods {
				newProd := ProductInToken{
					p.Id,
					p.Name,
					p.Alias,
					p.Disabled,
				}
				products = append(products, newProd)
				if p.IsOwner {
					ownProducts = append(ownProducts, newProd)
				}
			}
		}
		tc.Products = &products
		tc.OwnProducts = &ownProducts

	}(&wg)

	// 填入组信息
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		groups := make([]GroupInToken, 0)
		ownGroups := make([]GroupInToken, 0)

		gps, suc := db.UserModel.ListGroups(userObj.Id)
		if !suc {
			log.Error("Can not load groups, Error in new token db.UserModel.ListGroups.")
		} else {
			for _, g := range *gps {
				newGp := GroupInToken{
					g.Id,
					g.Name,
				}
				groups = append(groups, newGp)
				if g.IsOwner {
					ownGroups = append(ownGroups, newGp)
				}
			}
		}
		tc.Groups = &groups
		tc.OwnGroups = &ownGroups

	}(&wg)

	// 计算token值
	baseToken := tools.GenSha256HashCode(userObj.Username + currTime)
	token := tools.GenSha256HashCode(baseToken + conf.Config.SecretKey)
	// 生成TokenKey
	tokenKey := genTokenKey(token)

	// 等待所有信息填入结束
	wg.Wait()

	// TokenContent存到redis
	tokenContent, _ := json.Marshal(tc)
	if err := redis.Redis.Set(tokenKey, string(tokenContent), conf.Config.TokenTtl); err != nil {
		log.ErrorWithFields(log.Fields{
			"user_id":  userObj.Id,
			"username": userObj.Username,
			"error":    err.Error(),
		}, "Error in NewToken.")
		return ""
	}
	return token
}

// RemoveToken 本地服务::删除Token
func RemoveToken(token string) {
	if err := redis.Redis.Del(genTokenKey(token)); err != nil {
		log.ErrorWithFields(log.Fields{
			"token": token,
			"error": err.Error(),
		}, "Error in RemoveToken.")
	}
}

// RenewalToken 本地服务::续期Token
func RenewalToken(token string) {
	if err := redis.Redis.Expire(genTokenKey(token), conf.Config.TokenTtl); err != nil {
		log.ErrorWithFields(log.Fields{
			"token": token,
			"error": err.Error(),
		}, "Error in RenewalToken.")
	}
}

// VerifyToken 本地服务::验证Token是否有效
func VerifyToken(token string) bool {
	return redis.Redis.HasKey(genTokenKey(token))
}

// GetTokenContent 本地服务::通过Token获得TokenContent
func GetTokenContent(token string) (*TokenContent, bool) {
	if !redis.Redis.HasKey(genTokenKey(token)) {
		log.WarnWithFields(log.Fields{
			"token": token,
		}, "Redis key not found in GetTokenContent.")
		return nil, true
	}
	content, err := redis.Redis.Get(genTokenKey(token))
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"token": token,
			"error": err.Error(),
		}, "Error in GetTokenContent.")
		return nil, false
	}

	tc := TokenContent{}
	if err := json.Unmarshal([]byte(content), &tc); err != nil {
		log.ErrorWithFields(log.Fields{
			"token": token,
			"error": err.Error(),
		}, "Error in GetTokenContent.")
		return nil, false
	}

	return &tc, true
}

// genTokenKey 生成TokenKey
func genTokenKey(token string) string {
	return "token/" + token
}
