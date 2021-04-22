package handler

import (
	"database/sql"
	"eago-auth/api/form"
	"eago-auth/conf/msg"
	db "eago-auth/database"
	"eago-common/log"
	"eago-common/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// NewProduct 新建产品线
// @Summary 新建产品线
// @Tags 产品线
// @Param token header string true "Token"
// @Param data body db.Product true "body"
// @Success 200 {string} string "{"code":0,"message":"Success","product":{"id":2,"name":"new_role"}}"
// @Router /products [POST]
func NewProduct(c *gin.Context) {
	var prod db.Product

	// 序列化request body
	if err := c.ShouldBindJSON(&prod); err != nil {
		m := msg.WarnInvalidBody.NewMsg("Field 'name', 'alias', 'disabled', 'description' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	p := db.ProductModel.New(prod.Name, prod.Alias, prod.Description, *prod.Disabled)
	if p == nil {
		m := msg.ErrDatabase.NewMsg("Error in db.ProductModel.New.")
		log.Error(m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	m := msg.Success.NewMsg().SetPayload(&gin.H{"product": p})
	c.JSON(http.StatusOK, m.GinH())
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
		m := msg.WarnInvalidUri.NewMsg("Field 'product_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	if suc := db.ProductModel.Remove(prodId); !suc {
		m := msg.ErrDatabase.NewMsg("Error in db.ProductModel.Remove.")
		log.Error(m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	m := msg.Success.NewMsg()
	c.JSON(http.StatusOK, m.GinH())
}

// SetProduct 更新产品线
// @Summary 更新产品线
// @Tags 产品线
// @Param token header string true "Token"
// @Param product_id path string true "产品线ID"
// @Param data body db.Product true "body"
// @Success 200 {string} string "{"code":0,"message":"Success","product":{"id":2,"name":"new_role"}}"
// @Router /products/{product_id} [PUT]
func SetProduct(c *gin.Context) {
	var prodFm db.Product

	prodId, err := strconv.Atoi(c.Param("product_id"))
	if err != nil {
		m := msg.WarnInvalidUri.NewMsg("Field 'product_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	// 序列化request body
	if err := c.ShouldBindJSON(&prodFm); err != nil {
		m := msg.WarnInvalidBody.NewMsg("Field 'name' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	prod, suc := db.ProductModel.Set(prodId, prodFm.Name, prodFm.Alias, prodFm.Description, *prodFm.Disabled)
	if !suc {
		m := msg.ErrDatabase.NewMsg("Error in db.ProductModel.Set.")
		log.Error(m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	m := msg.Success.NewMsg().SetPayload(&gin.H{"product": prod})
	c.JSON(http.StatusOK, m.GinH())
}

// ListProducts 列出所有产品线
// @Summary 列出所有产品线
// @Tags 产品线
// @Param token header string true "Token"
// @Param query query string false "过滤条件"
// @Param page query string false "页数"
// @Param page_size query string false "页尺寸"
// @Success 200 {string} string "{"code":0,"message":"Success","page":1,"page_size":50,"pages":1,"products":[{"id":1,"name":"prodtct2","alias":"p2","disabled":false,"description":"1233","created_at":"2021-01-19 15:10:35","updated_at":"2021-01-19 15:10:35"}],"total":1}"
// @Router /products [GET]
func ListProducts(c *gin.Context) {
	var query db.Query

	q := c.GetString("Query")
	if q != "" {
		likeQuery := fmt.Sprintf("%%%s%%", q)
		query = db.Query{"name LIKE @query OR alias LIKE @query id LIKE @query": sql.Named("query", likeQuery)}
	}

	paged, suc := db.ProductModel.PagedList(
		&query,
		c.GetInt("Page"),
		c.GetInt("PageSize"),
		c.GetStringSlice("OrderBy")...,
	)
	if !suc {
		m := msg.ErrDatabase.NewMsg("Error in db.ProductModel.PageList.")
		log.Error(m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	m := msg.Success.NewMsg().SetPagedPayload(paged, "products")
	c.JSON(http.StatusOK, m.GinH())
}

// AddUser2Product 添加用户至产品线
// @Summary 添加用户至产品线
// @Tags 产品线
// @Param token header string true "Token"
// @Param product_id path string true "产品线ID"
// @Param data body db.UserProduct true "body"
// @Success 200 {string} string "{"code":0,"message":"Success"}"
// @Router /products/{product_id}/users [POST]
func AddUser2Product(c *gin.Context) {
	var uProd db.UserProduct

	prodId, err := strconv.Atoi(c.Param("product_id"))
	if err != nil {
		m := msg.WarnInvalidUri.NewMsg("Field 'product_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	// 序列化request body
	if err := c.ShouldBindJSON(&uProd); err != nil {
		m := msg.WarnInvalidBody.NewMsg("Field 'user_id', 'is_owner' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	if !db.ProductModel.AddUser(uProd.UserId, prodId, *uProd.IsOwner) {
		m := msg.ErrDatabase.NewMsg("Error in db.ProductModel.AddUser.")
		log.Error(m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	m := msg.Success.NewMsg()
	c.JSON(http.StatusOK, m.GinH())
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
		m := msg.WarnInvalidUri.NewMsg("Field 'product_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	userId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		m := msg.WarnInvalidUri.NewMsg("Field 'user_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	if !db.ProductModel.RemoveUser(userId, prodId) {
		m := msg.ErrDatabase.NewMsg("Error in db.ProductModel.RemoveUser.")
		log.Error(m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	m := msg.Success.NewMsg()
	c.JSON(http.StatusOK, m.GinH())
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
		m := msg.WarnInvalidUri.NewMsg("Field 'product_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	userId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		m := msg.WarnInvalidUri.NewMsg("Field 'user_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	// 序列化request body
	if err := c.ShouldBindJSON(&fm); err != nil {
		m := msg.WarnInvalidBody.NewMsg("Field 'is_owner' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	if !db.ProductModel.SetUserIsOwner(userId, prodId, *fm.IsOwner) {
		m := msg.ErrDatabase.NewMsg("Error in db.ProductModel.SetUserIsOwner.")
		log.Error(m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	m := msg.Success.NewMsg()
	c.JSON(http.StatusOK, m.GinH())
}

// ListProductUsers 列出角色所有用户
// @Summary 列出角色所有用户
// @Tags 产品线
// @Param token header string true "Token"
// @Param product_id path string true "产品线ID"
// @Success 200 {string} string "{"code":0,"message":"Success","users":[{"id":4,"username":"test2","is_owner":false,"joined_at":"2021-01-20 11:01:16"},{"id":3,"username":"test","is_owner":true,"joined_at":"2021-01-20 11:01:32"}]}"
// @Router /products/{product_id}/users [GET]
func ListProductUsers(c *gin.Context) {
	var query = db.Query{}

	prodId, err := strconv.Atoi(c.Param("product_id"))
	if err != nil {
		m := msg.WarnInvalidUri.NewMsg("Field 'product_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	// TODO(lai.li)
	// 方法中is_owner传值只能是0 or 1，待将来解决
	isOwner, err := strconv.Atoi(c.DefaultQuery("is_owner", "-1"))
	if err != nil {
		m := msg.WarnInvalidUri.NewMsg("Field 'is_owner' required, and must integer 0 or 1.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}
	if isOwner >= 0 {
		query["is_owner"] = tools.IntMin(isOwner, 1)
	}

	u, suc := db.ProductModel.ListUsers(prodId, &query)
	if !suc {
		m := msg.ErrDatabase.NewMsg("Error in db.ProductModel.ListUsers.")
		log.Error(m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	m := msg.Success.NewMsg().SetPayload(&gin.H{"users": u})
	c.JSON(http.StatusOK, m.GinH())
}
