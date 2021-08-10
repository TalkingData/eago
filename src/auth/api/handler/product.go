package handler

import (
	"database/sql"
	"eago/auth/api/form"
	"eago/auth/conf/msg"
	"eago/auth/model"
	"eago/common/log"
	"eago/common/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

// NewProduct 新建产品线
// @Summary 新建产品线
// @Tags 产品线
// @Param token header string true "Token"
// @Param data body model.Product true "body"
// @Success 200 {string} string "{"code":0,"message":"Success","product":{"id":2,"name":"new_role"}}"
// @Router /products [POST]
func NewProduct(c *gin.Context) {
	var prod model.Product

	// 序列化request body
	if err := c.ShouldBindJSON(&prod); err != nil {
		resp := msg.WarnInvalidBody.GenResponse("Field 'name', 'alias', 'disabled', 'description' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())
		resp.Write(c)
		return
	}

	p := model.NewProduct(prod.Name, prod.Alias, *prod.Description, prod.Disabled)
	if p == nil {
		resp := msg.ErrDatabase.GenResponse("Error when NewProduct.")
		log.Error(resp.String())
		resp.Write(c)
		return
	}

	resp := msg.Success.GenResponse().SetPayload("product", p)
	resp.Write(c)
}

// RemoveProduct 删除产品线
// @Summary 删除产品线
// @Tags 产品线
// @Param token header string true "Token"
// @Param product_id path string true "产品线ID"
// @Success 200 {string} string "{"code":0,"message":"Success"}"
// @Router /products/{product_id} [DELETE]
func RemoveProduct(c *gin.Context) {
	prodId, err := strconv.Atoi(c.Param("product_id"))
	if err != nil {
		resp := msg.WarnInvalidUri.GenResponse("Field 'product_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())
		resp.Write(c)
		return
	}

	if ok := model.RemoveProduct(prodId); !ok {
		resp := msg.ErrDatabase.GenResponse("Error when RemoveProduct.")
		log.Error(resp.String())
		resp.Write(c)
		return
	}

	resp := msg.Success.GenResponse()
	resp.Write(c)
}

// SetProduct 更新产品线
// @Summary 更新产品线
// @Tags 产品线
// @Param token header string true "Token"
// @Param product_id path string true "产品线ID"
// @Param data body model.Product true "body"
// @Success 200 {string} string "{"code":0,"message":"Success","product":{"id":2,"name":"new_role"}}"
// @Router /products/{product_id} [PUT]
func SetProduct(c *gin.Context) {
	var prodFm model.Product

	prodId, err := strconv.Atoi(c.Param("product_id"))
	if err != nil {
		resp := msg.WarnInvalidUri.GenResponse("Field 'product_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())
		resp.Write(c)
		return
	}

	// 序列化request body
	if err := c.ShouldBindJSON(&prodFm); err != nil {
		resp := msg.WarnInvalidBody.GenResponse("Field 'name' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())
		resp.Write(c)
		return
	}

	prod, ok := model.SetProduct(prodId, prodFm.Name, prodFm.Alias, *prodFm.Description, *prodFm.Disabled)
	if !ok {
		resp := msg.ErrDatabase.GenResponse("Error when SetProduct.")
		log.Error(resp.String())
		resp.Write(c)
		return
	}

	resp := msg.Success.GenResponse().SetPayload("product", prod)
	resp.Write(c)
}

// ListProducts 列出所有产品线
// @Summary 列出所有产品线
// @Tags 产品线
// @Param token header string true "Token"
// @Param query query string false "过滤条件"
// @Param order_by query string false "排序字段(多个间逗号分割)"
// @Param page query string false "页数"
// @Param page_size query string false "页尺寸"
// @Success 200 {string} string "{"code":0,"message":"Success","page":1,"page_size":50,"pages":1,"products":[{"id":1,"name":"prodtct2","alias":"p2","disabled":false,"description":"1233","created_at":"2021-01-19 15:10:35","updated_at":"2021-01-19 15:10:35"}],"total":1}"
// @Router /products [GET]
func ListProducts(c *gin.Context) {
	var query model.Query

	q := c.GetString("Query")
	if q != "" {
		likeQuery := fmt.Sprintf("%%%s%%", q)
		query = model.Query{"name LIKE @query OR alias LIKE @query OR id LIKE @query": sql.Named("query", likeQuery)}
	}

	paged, ok := model.PagedListProducts(
		query,
		c.GetInt("Page"),
		c.GetInt("PageSize"),
		c.GetStringSlice("OrderBy")...,
	)
	if !ok {
		resp := msg.ErrDatabase.GenResponse("Error when PageListProducts.")
		log.Error(resp.String())
		resp.Write(c)
		return
	}

	resp := msg.Success.GenResponse().SetPagedPayload(paged, "products")
	resp.Write(c)
}

// AddUser2Product 添加用户至产品线
// @Summary 添加用户至产品线
// @Tags 产品线
// @Param token header string true "Token"
// @Param product_id path string true "产品线ID"
// @Param data body model.UserProduct true "body"
// @Success 200 {string} string "{"code":0,"message":"Success"}"
// @Router /products/{product_id}/users [POST]
func AddUser2Product(c *gin.Context) {
	var uProd model.UserProduct

	prodId, err := strconv.Atoi(c.Param("product_id"))
	if err != nil {
		resp := msg.WarnInvalidUri.GenResponse("Field 'product_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())
		resp.Write(c)
		return
	}

	// 序列化request body
	if err := c.ShouldBindJSON(&uProd); err != nil {
		resp := msg.WarnInvalidBody.GenResponse("Field 'user_id', 'is_owner' required, and 'user_id' must greater than 0.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())
		resp.Write(c)
		return
	}

	if !model.AddProductUser(uProd.UserId, prodId, *uProd.IsOwner) {
		resp := msg.ErrDatabase.GenResponse("Error when AddProductUser.")
		log.Error(resp.String())
		resp.Write(c)
		return
	}

	resp := msg.Success.GenResponse()
	resp.Write(c)
}

// RemoveProductUser 移除产品线中用户
// @Summary 移除产品线中用户
// @Tags 产品线
// @Param token header string true "Token"
// @Param product_id path string true "产品线ID"
// @Param user_id path string true "用户ID"
// @Success 200 {string} string "{"code":0,"message":"Success"}"
// @Router /products/{product_id}/users/{user_id} [DELETE]
func RemoveProductUser(c *gin.Context) {
	prodId, err := strconv.Atoi(c.Param("product_id"))
	if err != nil {
		resp := msg.WarnInvalidUri.GenResponse("Field 'product_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())
		resp.Write(c)
		return
	}

	userId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		resp := msg.WarnInvalidUri.GenResponse("Field 'user_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())
		resp.Write(c)
		return
	}

	if !model.RemoveProductUser(userId, prodId) {
		resp := msg.ErrDatabase.GenResponse("Error when RemoveProductUser.")
		log.Error(resp.String())
		resp.Write(c)
		return
	}

	resp := msg.Success.GenResponse()
	resp.Write(c)
}

// SetUserIsProductOwner 设置用户是否是产品线Owner
// @Summary 设置用户是否是产品线Owner
// @Tags 产品线
// @Param token header string true "Token"
// @Param product_id path string true "产品线ID"
// @Param user_id path string true "用户ID"
// @Param data body form.IsOwnerForm true "body"
// @Success 200 {string} string "{"code":0,"message":"Success"}"
// @Router /products/{product_id}/users/{user_id} [PUT]
func SetUserIsProductOwner(c *gin.Context) {
	var fm = form.IsOwnerForm{}

	prodId, err := strconv.Atoi(c.Param("product_id"))
	if err != nil {
		resp := msg.WarnInvalidUri.GenResponse("Field 'product_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())
		resp.Write(c)
		return
	}

	userId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		resp := msg.WarnInvalidUri.GenResponse("Field 'user_id' required")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())
		resp.Write(c)
		return
	}

	// 序列化request body
	if err := c.ShouldBindJSON(&fm); err != nil {
		resp := msg.WarnInvalidBody.GenResponse("Field 'is_owner' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())
		resp.Write(c)
		return
	}

	if !model.SetProductUserIsOwner(userId, prodId, *fm.IsOwner) {
		resp := msg.ErrDatabase.GenResponse("Error when SetProductUserIsOwner.")
		log.Error(resp.String())
		resp.Write(c)
		return
	}

	resp := msg.Success.GenResponse()
	resp.Write(c)
}

// ListProductUsers 列出产品线所有用户
// @Summary 列出产品线所有用户
// @Tags 产品线
// @Param token header string true "Token"
// @Param product_id path string true "产品线ID"
// @Success 200 {string} string "{"code":0,"message":"Success","users":[{"id":4,"username":"test2","is_owner":false,"joined_at":"2021-01-20 11:01:16"},{"id":3,"username":"test","is_owner":true,"joined_at":"2021-01-20 11:01:32"}]}"
// @Router /products/{product_id}/users [GET]
func ListProductUsers(c *gin.Context) {
	var query = model.Query{}

	prodId, err := strconv.Atoi(c.Param("product_id"))
	if err != nil {
		resp := msg.WarnInvalidUri.GenResponse("Field 'product_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())
		resp.Write(c)
		return
	}

	// TODO(lai.li)
	// 方法中is_owner传值只能是0 or 1，待将来解决
	isOwner, err := strconv.Atoi(c.DefaultQuery("is_owner", "-1"))
	if err != nil {
		resp := msg.WarnInvalidUri.GenResponse("Field 'is_owner' required, and must integer 0 or 1.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())
		resp.Write(c)
		return
	}
	if isOwner >= 0 {
		query["is_owner"] = utils.IntMin(isOwner, 1)
	}

	u, ok := model.ListProductUsers(prodId, query)
	if !ok {
		resp := msg.ErrDatabase.GenResponse("Error in ListProductUsers.")
		log.Error(resp.String())
		resp.Write(c)
		return
	}

	resp := msg.Success.GenResponse().SetPayload("users", u)
	resp.Write(c)
}
