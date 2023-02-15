package biz

import (
	"context"
	"eago/auth/dto"
	"eago/auth/model"
	"eago/common/global"
	"eago/common/logger"
	"eago/common/utils"
	"encoding/json"
	"sync"
	"time"
)

// NewToken 生成Token
func (b *Biz) NewToken(ctx context.Context, userObj *model.User) string {
	b.logger.InfoWithFields(logger.Fields{
		"user_id": userObj.Id,
	}, "biz.NewToken called.")
	defer b.logger.Info("biz.NewToken end.")

	// 计算token值
	currTime := time.Now().Format(global.TimestampFormat)
	baseToken := utils.GenSha256HashCode(userObj.Username + currTime)
	token := utils.GenSha256HashCode(baseToken + b.conf.SecretKey)
	// 生成TokenKey
	tokenKey := genTokenKey(token)

	// TokenContent存到redis
	b.logger.Info(logger.Fields{
		"user_id":   userObj.Id,
		"token_ket": tokenKey,
	}, "Writing token to redis.")
	tokenContent, _ := json.Marshal(b.genTokenContent(ctx, userObj))
	if err := b.redis.Set(ctx, tokenKey, string(tokenContent), b.conf.TokenTtl); err != nil {
		b.logger.ErrorWithFields(logger.Fields{
			"user_id":  userObj.Id,
			"username": userObj.Username,
			"error":    err,
		}, "An error occurred while redis.Set in biz.NewToken.")
		return ""
	}
	return token
}

// RemoveToken 删除Token
func (b *Biz) RemoveToken(ctx context.Context, token string) {
	b.logger.InfoWithFields(logger.Fields{"token": token}, "biz.RemoveToken called.")
	defer b.logger.Info("biz.RemoveToken end.")

	b.logger.DebugWithFields(logger.Fields{"token": token}, "Remove token from redis.")
	if err := b.redis.Del(ctx, genTokenKey(token)); err != nil {
		b.logger.WarnWithFields(logger.Fields{
			"token": token,
			"error": err,
		}, "An error occurred while redis.Del in biz.RemoveToken, but skipped.")
	}
}

// RenewalToken 续期Token
func (b *Biz) RenewalToken(ctx context.Context, token string) {
	b.logger.DebugWithFields(logger.Fields{"token": token}, "biz.RenewalToken called.")
	defer b.logger.Debug("biz.RenewalToken end.")

	if err := b.redis.Expire(ctx, genTokenKey(token), b.conf.TokenTtl); err != nil {
		b.logger.WarnWithFields(logger.Fields{
			"token": token,
			"error": err,
		}, "An error occurred while redis.Expire in biz.RenewalToken, but skipped.")
	}
}

// VerifyToken 验证Token是否有效
func (b *Biz) VerifyToken(ctx context.Context, token string) bool {
	return b.redis.Exist(ctx, genTokenKey(token))
}

// GetTokenContent 通过Token获得TokenContent
func (b *Biz) GetTokenContent(ctx context.Context, token string) (tc *dto.TokenContent, err error) {
	b.logger.InfoWithFields(logger.Fields{
		"token": "token",
	}, "biz.GetTokenContent called.")
	defer b.logger.Info("biz.GetTokenContent end.")

	if !b.redis.Exist(ctx, genTokenKey(token)) {
		b.logger.WarnWithFields(logger.Fields{
			"token": token,
		}, "Redis key not found in biz.GetTokenContent.")
		return nil, nil
	}
	tcStr, err := b.redis.Get(ctx, genTokenKey(token))
	if err != nil {
		b.logger.ErrorWithFields(logger.Fields{
			"token": token,
			"error": err,
		}, "An error occurred while biz.GetTokenContent.")
		return nil, err
	}

	b.logger.Debug("Unmarshal token content.")
	if err = json.Unmarshal([]byte(tcStr), &tc); err != nil {
		b.logger.ErrorWithFields(logger.Fields{
			"token": token,
			"error": err,
		}, "An error occurred while biz.GetTokenContent.")
		return nil, err
	}

	return tc, nil
}

// genTokenContent 生成TokenContent
func (b *Biz) genTokenContent(ctx context.Context, userObj *model.User) *dto.TokenContent {
	b.logger.Info("biz.writeTokenContent called.")
	defer b.logger.Info("biz.writeTokenContent end.")

	tc := &dto.TokenContent{
		UserId:   userObj.Id,
		Username: userObj.Username,
		Phone:    userObj.Phone,

		IsSuperuser: userObj.IsSuperuser,
	}

	// 填入角色信息
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		b.logger.InfoWithFields(logger.Fields{"user_id": userObj.Id}, "Loading user roles.")
		roles, err := b.dao.ListUsersRoles(ctx, userObj.Id)
		if err != nil {
			b.logger.ErrorWithFields(logger.Fields{
				"user_id": userObj.Id,
			}, "An error occurred while dao.ListUsersRoles in biz.GetTokenContent.")
		}

		tc.Roles = make([]*string, len(roles))
		for idx, r := range roles {
			rName := r.Name
			tc.Roles[idx] = &rName
		}
	}(wg)

	// 填入产品线信息
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		b.logger.InfoWithFields(logger.Fields{"user_id": userObj.Id}, "Loading user products.")
		prods, err := b.dao.ListUsersProducts(ctx, userObj.Id)
		if err != nil {
			b.logger.ErrorWithFields(logger.Fields{
				"user_id": userObj.Id,
			}, "An error occurred while dao.ListUsersProducts in biz.GetTokenContent.")
		}

		tc.Products = make([]*dto.ProductInToken, 0)
		tc.OwnProducts = make([]*dto.ProductInToken, 0)
		for _, p := range prods {
			newProd := dto.ProductInToken{
				Id:       p.Id,
				Name:     p.Name,
				Alias:    p.Alias,
				Disabled: p.Disabled,
			}
			tc.Products = append(tc.Products, &newProd)
			if p.IsOwner {
				tc.OwnProducts = append(tc.OwnProducts, &newProd)
			}
		}
	}(wg)

	// 填入组信息
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		b.logger.InfoWithFields(logger.Fields{"user_id": userObj.Id}, "Loading user groups.")
		gps, err := b.dao.ListUsersGroups(ctx, userObj.Id)
		if err != nil {
			b.logger.ErrorWithFields(logger.Fields{
				"user_id": userObj.Id,
			}, "An error occurred while dao.ListUsersGroups in biz.GetTokenContent.")
		}

		tc.Groups = make([]*dto.GroupInToken, 0)
		tc.OwnGroups = make([]*dto.GroupInToken, 0)
		for _, g := range gps {
			newGp := dto.GroupInToken{
				Id:   g.Id,
				Name: g.Name,
			}
			tc.Groups = append(tc.Groups, &newGp)
			if g.IsOwner {
				tc.OwnGroups = append(tc.OwnGroups, &newGp)
			}
		}
	}(wg)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		b.logger.InfoWithFields(logger.Fields{"user_id": userObj.Id}, "Loading user department.")
		tc.Department = make([]*string, 0)
		deptTree := b.dao.GetUsersDepartmentChain(ctx, userObj.Id)
		for _, dept := range deptTree {
			deptName := dept.Name
			tc.Department = append(tc.Department, &deptName)
		}
	}(wg)

	// 等待所有信息填入结束
	wg.Wait()
	return tc
}

// genTokenKey 生成TokenKey
func genTokenKey(token string) string {
	return "token/" + token
}
