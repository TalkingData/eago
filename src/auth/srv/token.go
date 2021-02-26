package srv

import (
	"context"
	"eago-auth/config"
	db "eago-auth/database"
	"eago-auth/srv/proto"
	"eago-common/log"
	"eago-common/redis"
	"eago-common/tools"
	"encoding/json"
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

// RPC服务::验证Token是否有效
func (as *AuthService) VerifyToken(ctx context.Context, req *auth.Token, res *auth.BoolResponse) error {
	log.InfoWithFields(log.Fields{"token": req.Token}, "Gor rpc call verify token.")
	res.Ok = VerifyToken(req.Token)
	return nil
}

// RPC服务::通过Token获得TokenContent
func (as *AuthService) GetTokenContent(ctx context.Context, req *auth.Token, res *auth.TokenContent) error {
	log.InfoWithFields(
		log.Fields{"token": req.Token},
		"Got rpc call get token content.",
	)
	tc, suc := GetTokenContent(req.Token)
	if !suc {
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
	loadProduct := make(chan bool, 1)
	go func(done chan bool) {
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

		done <- true
	}(loadProduct)

	// 加载组
	loadGroup := make(chan bool, 1)
	go func(done chan bool) {
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

		done <- true
	}(loadGroup)

	<-loadProduct
	<-loadGroup

	return nil
}

func intCopy2Int32(from *[]int, to *[]int32, done chan bool) {
	if from != nil {
		t := *to
		for _, f := range *from {
			t = append(t, int32(f))
		}
	}

	done <- true
}

// 本地服务::生成Token
func NewToken(userObj *db.User) string {
	var (
		tc       = TokenContent{}
		currTime = time.Now().Format(config.DEFAULT_TIMESTAMP_FORMAT)
	)

	tc.UserId = userObj.Id
	tc.Username = userObj.Username

	tc.IsSuperuser = userObj.IsSuperuser

	// 填入角色信息
	loadRoles := make(chan bool, 1)
	go func(done chan bool) {
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

		done <- true
	}(loadRoles)

	// 填入产品线信息
	loadProduct := make(chan bool, 1)
	go func(done chan bool) {
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

		done <- true
	}(loadProduct)

	// 填入组信息
	loadGroup := make(chan bool, 1)
	go func(done chan bool) {
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

		done <- true
	}(loadGroup)

	// 计算token值
	baseToken := tools.GenSha256HashCode(userObj.Username + currTime)
	token := tools.GenSha256HashCode(baseToken + config.Config.SecretKey)
	// 生成TokenKey
	tokenKey := genTokenKey(token)

	// 等待所有信息填入结束
	<-loadRoles
	<-loadProduct
	<-loadGroup
	// TokenContent存到redis
	tokenContent, _ := json.Marshal(tc)
	if err := redis.Redis.Set(tokenKey, string(tokenContent), config.Config.TokenTtl); err != nil {
		log.ErrorWithFields(log.Fields{
			"user_id":  userObj.Id,
			"username": userObj.Username,
			"error":    err.Error(),
		}, "Error in NewToken.")
		return ""
	}
	return token
}

// 本地服务::删除Token
func DeleteToken(token string) {
	redis.Redis.Del(genTokenKey(token))
}

// 本地服务::续期Token
func RenewalToken(token string) {
	if err := redis.Redis.Expire(genTokenKey(token), config.Config.TokenTtl); err != nil {
		log.ErrorWithFields(log.Fields{
			"token": token,
			"error": err.Error(),
		}, "Error in RenewalToken.")
	}
}

// 本地服务::验证Token是否有效
func VerifyToken(token string) bool {
	return redis.Redis.HasKey(genTokenKey(token))
}

// 本地服务::通过Token获得TokenContent
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

// 生成TokenKey
func genTokenKey(token string) string {
	return "token/" + token
}
