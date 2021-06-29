package local

import (
	"eago/auth/conf"
	"eago/auth/model"
	"eago/auth/srv/dto"
	"eago/common/log"
	"eago/common/redis"
	"eago/common/utils"
	"encoding/json"
	"sync"
	"time"
)

// NewToken 本地服务::生成Token
func NewToken(userObj *model.User) string {
	tc := dto.TokenContent{}
	tc.UserId = userObj.Id
	tc.Username = userObj.Username
	tc.Phone = userObj.Phone

	tc.IsSuperuser = userObj.IsSuperuser

	// 填入角色信息
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		rolesStr := make([]string, 0)
		roles, ok := model.ListUserRoles(userObj.Id)
		if !ok {
			log.Error("Can not load roles, Error in new token model.ListUserProducts.")
		} else {
			for _, r := range *roles {
				rolesStr = append(rolesStr, r.Name)
			}
		}
		tc.Roles = &rolesStr

	}(wg)

	// 填入产品线信息
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		products := make([]dto.ProductInToken, 0)
		ownProducts := make([]dto.ProductInToken, 0)

		prods, ok := model.ListUserProducts(userObj.Id)
		if !ok {
			log.Error("Can not load products, Error in new token model.ListUserProducts.")
		} else {
			for _, p := range *prods {
				newProd := dto.ProductInToken{
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

	}(wg)

	// 填入组信息
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		groups := make([]dto.GroupInToken, 0)
		ownGroups := make([]dto.GroupInToken, 0)

		gps, ok := model.ListUserGroups(userObj.Id)
		if !ok {
			log.Error("Can not load groups, Error in new token model.ListUserGroups.")
		} else {
			for _, g := range *gps {
				newGp := dto.GroupInToken{
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

	}(wg)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		dep := make([]string, 0)
		tc.Department = &dep
	}(wg)

	// 计算token值
	currTime := time.Now().Format(conf.TIMESTAMP_FORMAT)
	baseToken := utils.GenSha256HashCode(userObj.Username + currTime)
	token := utils.GenSha256HashCode(baseToken + conf.Config.SecretKey)
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
func GetTokenContent(token string) (*dto.TokenContent, bool) {
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

	tc := dto.TokenContent{}
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
