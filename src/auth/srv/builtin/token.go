package builtin

import (
	"eago/auth/conf"
	"eago/auth/dao"
	"eago/auth/model"
	"eago/common/log"
	"eago/common/redis"
	"eago/common/utils"
	"encoding/json"
	"sync"
	"time"
)

// NewToken 本地服务::生成Token
func NewToken(userObj *model.User) string {
	log.Info("builtin.NewToken called.")
	defer log.Info("builtin.NewToken end.")

	tc := model.TokenContent{}
	tc.UserId = userObj.Id
	tc.Username = userObj.Username
	tc.Phone = userObj.Phone

	tc.IsSuperuser = userObj.IsSuperuser

	// 填入角色信息
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		rolesStr := make([]string, 0)

		log.InfoWithFields(log.Fields{"user_id": userObj.Id}, "Loading user roles.")
		roles, ok := dao.ListUserRoles(userObj.Id)
		if !ok {
			log.Error("Can not load roles, Error in new token dao.ListUserRoles.")
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

		products := make([]model.ProductInToken, 0)
		ownProducts := make([]model.ProductInToken, 0)

		log.InfoWithFields(log.Fields{"user_id": userObj.Id}, "Loading user products.")
		prods, ok := dao.ListUserProducts(userObj.Id)
		if !ok {
			log.Error("Can not load products, Error in new token dao.ListUserProducts.")
		} else {
			for _, p := range *prods {
				newProd := model.ProductInToken{
					Id:       p.Id,
					Name:     p.Name,
					Alias:    p.Alias,
					Disabled: p.Disabled,
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

		groups := make([]model.GroupInToken, 0)
		ownGroups := make([]model.GroupInToken, 0)

		log.InfoWithFields(log.Fields{"user_id": userObj.Id}, "Loading user groups.")
		gps, ok := dao.ListUserGroups(userObj.Id)
		if !ok {
			log.Error("Can not load groups, Error in new token dao.ListUserGroups.")
		} else {
			for _, g := range *gps {
				newGp := model.GroupInToken{
					Id:   g.Id,
					Name: g.Name,
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

		log.InfoWithFields(log.Fields{"user_id": userObj.Id}, "Loading user department.")
		dep := make([]string, 0)
		tc.Department = &dep
	}(wg)

	// 计算token值
	currTime := time.Now().Format(conf.TIMESTAMP_FORMAT)
	baseToken := utils.GenSha256HashCode(userObj.Username + currTime)
	token := utils.GenSha256HashCode(baseToken + conf.Conf.SecretKey)
	// 生成TokenKey
	tokenKey := genTokenKey(token)

	// 等待所有信息填入结束
	wg.Wait()

	// TokenContent存到redis
	log.Info(log.Fields{
		"user_id":   userObj.Id,
		"token_ket": tokenKey,
	}, "Write token to redis.")
	tokenContent, _ := json.Marshal(tc)
	if err := redis.Redis.Set(tokenKey, string(tokenContent), conf.Conf.TokenTtl); err != nil {
		log.ErrorWithFields(log.Fields{
			"user_id":  userObj.Id,
			"username": userObj.Username,
			"error":    err,
		}, "An error occurred while NewToken.")
		return ""
	}
	return token
}

// RemoveToken 本地服务::删除Token
func RemoveToken(token string) {
	log.Info("builtin.RemoveToken called.")
	defer log.Info("builtin.RemoveToken end.")

	log.DebugWithFields(log.Fields{"token": token}, "Remove token from redis.")
	if err := redis.Redis.Del(genTokenKey(token)); err != nil {
		log.ErrorWithFields(log.Fields{
			"token": token,
			"error": err,
		}, "An error occurred while RemoveToken.")
	}
}

// RenewalToken 本地服务::续期Token
func RenewalToken(token string) {
	log.Info("builtin.RenewalToken called.")
	defer log.Info("builtin.RenewalToken end.")

	if err := redis.Redis.Expire(genTokenKey(token), conf.Conf.TokenTtl); err != nil {
		log.ErrorWithFields(log.Fields{
			"token": token,
			"error": err,
		}, "An error occurred while RenewalToken.")
	}
}

// VerifyToken 本地服务::验证Token是否有效
func VerifyToken(token string) bool {
	return redis.Redis.HasKey(genTokenKey(token))
}

// GetTokenContent 本地服务::通过Token获得TokenContent
func GetTokenContent(token string) (*model.TokenContent, bool) {
	log.Info("builtin.GetTokenContent called.")
	defer log.Info("builtin.GetTokenContent end.")

	log.InfoWithFields(log.Fields{"token": token}, "Loading token from redis.")
	if !redis.Redis.HasKey(genTokenKey(token)) {
		log.WarnWithFields(log.Fields{
			"token": token,
		}, "Redis key not found in GetTokenContent.")
		return nil, false
	}
	content, err := redis.Redis.Get(genTokenKey(token))
	if err != nil {
		log.ErrorWithFields(log.Fields{
			"token": token,
			"error": err,
		}, "An error occurred while GetTokenContent.")
		return nil, false
	}

	log.Info("Unmarshal token content.")
	tc := model.TokenContent{}
	if err := json.Unmarshal([]byte(content), &tc); err != nil {
		log.ErrorWithFields(log.Fields{
			"token": token,
			"error": err,
		}, "An error occurred while GetTokenContent.")
		return nil, false
	}

	return &tc, true
}

// genTokenKey 生成TokenKey
func genTokenKey(token string) string {
	return "token/" + token
}
