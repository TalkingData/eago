package handler

import (
	"eago/auth/conf/msg"
	"eago/auth/dao"
	"eago/auth/dto"
	w "eago/common/api-suite/writter"
	"eago/common/log"
	"github.com/gin-gonic/gin"
	"strconv"
)

// NewProduct 新建产品线
func NewProduct(c *gin.Context) {
	var npFrm dto.NewProduct
	// 序列化request body
	if err := c.ShouldBindJSON(&npFrm); err != nil {
		m := msg.SerializeFailed.SetError(err)
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}
	// 验证数据
	if m := npFrm.Validate(); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	p, err := dao.NewProduct(npFrm.Name, npFrm.Alias, *npFrm.Description, npFrm.Disabled)
	if p == nil {
		m := msg.UnknownError.SetError(err)
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "product", p)
}

// RemoveProduct 删除产品线
func RemoveProduct(c *gin.Context) {
	prdId, err := strconv.Atoi(c.Param("product_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "product_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	var rProdFrm dto.RemoveProduct
	// 验证数据
	if m := rProdFrm.Validate(prdId); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	if err := dao.RemoveProduct(prdId); err != nil {
		m := msg.UnknownError.SetError(err)
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccess(c)
}

// SetProduct 更新产品线
func SetProduct(c *gin.Context) {
	prdId, err := strconv.Atoi(c.Param("product_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "product_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	var spFrm dto.SetProduct
	// 序列化request body
	if err = c.ShouldBindJSON(&spFrm); err != nil {
		m := msg.SerializeFailed.SetError(err)
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}
	// 验证数据
	if m := spFrm.Validate(prdId); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	prod, err := dao.SetProduct(prdId, spFrm.Name, spFrm.Alias, *spFrm.Description, *spFrm.Disabled)
	if err != nil {
		m := msg.UnknownError.SetError(err)
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "product", prod)
}

// ListProducts 列出所有产品线
func ListProducts(c *gin.Context) {
	query := dao.Query{}
	// 设置查询filter
	lpq := dto.ListProductsQuery{}
	if c.ShouldBindQuery(&lpq) == nil {
		_ = lpq.UpdateQuery(query)
	}

	paged, ok := dao.PagedListProducts(
		query,
		c.GetInt("Page"),
		c.GetInt("PageSize"),
		c.GetStringSlice("OrderBy")...,
	)
	if !ok {
		m := msg.UnknownError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "products", paged)
}

// AddUser2Product 添加用户至产品线
func AddUser2Product(c *gin.Context) {
	prdId, err := strconv.Atoi(c.Param("product_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "product_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	var aumFrm dto.AddUser2Product
	// 序列化request body
	if err = c.ShouldBindJSON(&aumFrm); err != nil {
		m := msg.SerializeFailed.SetError(err)
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}
	// 验证数据
	if m := aumFrm.Validate(prdId); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	if !dao.AddProductUser(prdId, aumFrm.UserId, aumFrm.IsOwner) {
		m := msg.UnknownError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccess(c)
}

// RemoveProductUser 移除产品线中用户
func RemoveProductUser(c *gin.Context) {
	prdId, err := strconv.Atoi(c.Param("product_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "product_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	userId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "user_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	var rpuFrm dto.RemoveProductUser
	// 验证数据
	if m := rpuFrm.Validate(prdId, userId); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	if !dao.RemoveProductUser(prdId, userId) {
		m := msg.UnknownError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccess(c)
}

// SetUserIsProductOwner 设置用户是否是产品线Owner
func SetUserIsProductOwner(c *gin.Context) {
	prdId, err := strconv.Atoi(c.Param("product_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "product_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	userId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "user_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	var suoFrm dto.SetUserIsProductOwner
	// 序列化request body
	if err = c.ShouldBindJSON(&suoFrm); err != nil {
		m := msg.SerializeFailed.SetError(err)
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}
	// 验证数据
	if m := suoFrm.Validate(prdId, userId); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	if !dao.SetProductUserIsOwner(prdId, userId, suoFrm.IsOwner) {
		m := msg.UnknownError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccess(c)
}

// ListProductUsers 列出产品线所有用户
func ListProductUsers(c *gin.Context) {
	prdId, err := strconv.Atoi(c.Param("product_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "product_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	var lpuFrm dto.ListProductUsersQuery
	// 序列化request body
	_ = c.ShouldBindQuery(&lpuFrm)
	if m := lpuFrm.Validate(prdId); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	query := dao.Query{}
	// 设置查询filter
	_ = lpuFrm.UpdateQuery(query)
	u, ok := dao.ListProductUsers(prdId, query)
	if !ok {
		m := msg.UnknownError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "users", u)
}
