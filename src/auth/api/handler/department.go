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

// NewDepartment 新建部门
// @Summary 新建部门
// @Tags 部门
// @Param token header string true "Token"
// @Param data body model.Department true "body"
// @Success 200 {string} string "{"code":0,"message":"Success","department":{"id":4,"name":"sub_dept3","parent_id":2,"created_at":"2021-01-21 15:11:00","updated_at":"2021-01-21 15:11:00"}}"
// @Router /departments [POST]
func NewDepartment(c *gin.Context) {
	var dForm model.Department

	// 序列化request body
	if err := c.ShouldBindJSON(&dForm); err != nil {
		resp := msg.WarnInvalidBody.GenResponse("Field 'name', 'parent_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())
		resp.Write(c)
		return
	}

	dept := model.NewDepartment(dForm.Name, dForm.ParentId)
	if dept == nil {
		resp := msg.ErrDatabase.GenResponse("Error when NewDepartment.")
		log.Error(resp.String())
		resp.Write(c)
		return
	}

	resp := msg.Success.GenResponse().SetPayload("department", dept)
	resp.Write(c)
}

// RemoveDepartment 删除部门
// @Summary 删除部门
// @Tags 部门
// @Param token header string true "Token"
// @Param department_id path string true "部门ID"
// @Success 200 {string} string "{"code":0,"message":"Success"}"
// @Router /departments/{department_id} [DELETE]
func RemoveDepartment(c *gin.Context) {
	deptId, err := strconv.Atoi(c.Param("department_id"))
	if err != nil {
		resp := msg.WarnInvalidUri.GenResponse("Field 'department_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())
		resp.Write(c)
		return
	}

	if ok := model.RemoveDepartment(deptId); !ok {
		resp := msg.ErrDatabase.GenResponse("Error when RemoveDepartment.")
		log.Error(resp.String())
		resp.Write(c)
		return
	}

	msg.Success.GenResponse().Write(c)
}

// SetDepartment 更新部门
// @Summary 更新部门
// @Tags 部门
// @Param token header string true "Token"
// @Param department_id path string true "部门ID"
// @Param data body model.Department true "body"
// @Success 200 {string} string "{"code":0,"message":"Success","department":{"id":4,"name":"sub_dept3","parent_id":2,"created_at":"2021-01-21 15:11:00","updated_at":"2021-01-21 15:11:00"}}"
// @Router /departments [PUT]
func SetDepartment(c *gin.Context) {
	var dForm model.Department

	deptId, err := strconv.Atoi(c.Param("department_id"))
	if err != nil {
		resp := msg.WarnInvalidUri.GenResponse("Field 'department_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())
		resp.Write(c)
		return
	}

	// 序列化request body
	if err := c.ShouldBindJSON(&dForm); err != nil {
		resp := msg.WarnInvalidBody.GenResponse("Field 'name', 'parent_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())
		resp.Write(c)
		return
	}

	dept, ok := model.SetDepartment(deptId, dForm.Name, dForm.ParentId)
	if !ok {
		resp := msg.ErrDatabase.GenResponse("Error when SetDepartment.")
		log.Error(resp.String())
		resp.Write(c)
		return
	}

	resp := msg.Success.GenResponse().SetPayload("department", dept)
	resp.Write(c)
}

// ListDepartments 列出所有部门
// @Summary 列出所有部门
// @Tags 部门
// @Param token header string true "Token"
// @Param query query string false "过滤条件"
// @Param order_by query string false "排序字段(多个间逗号分割)"
// @Param page query string false "页数"
// @Param page_size query string false "页尺寸"
// @Success 200 {string} string "{"code":0,"departments":[{"id":2,"name":"root","parent_id":null,"created_at":"2021-01-21 15:10:26","updated_at":"2021-01-21 15:10:26"},{"id":5,"name":"sub2","parent_id":2,"created_at":"2021-01-21 15:11:07","updated_at":"2021-01-21 15:11:07"}],"message":"Success","page":1,"page_size":50,"pages":1,"total":2}"
// @Router /departments [GET]
func ListDepartments(c *gin.Context) {
	var query model.Query

	q := c.GetString("Query")
	if q != "" {
		likeQuery := fmt.Sprintf("%%%s%%", q)
		query = model.Query{"name LIKE @query OR id LIKE @query": sql.Named("query", likeQuery)}
	}

	paged, ok := model.PagedListDepartments(
		query,
		c.GetInt("Page"),
		c.GetInt("PageSize"),
		c.GetStringSlice("OrderBy")...,
	)
	if !ok {
		resp := msg.ErrDatabase.GenResponse("Error when PageListDepartments.")
		log.Error(resp.String())
		resp.Write(c)
		return
	}

	resp := msg.Success.GenResponse().SetPagedPayload(paged, "departments")
	resp.Write(c)
}

// ListDepartmentsTree 以树结构列出所有部门
// @Summary 以树结构列出所有部门
// @Tags 部门
// @Param token header string true "Token"
// @Success 200 {string} string "{"code":0,"message":"Success","tree":{"id":2,"name":"root","sub_department":[{"id":4,"name":"sub1","sub_department":[],"created_at":"2021-01-21 15:11:00","updated_at":"2021-01-21 15:11:00"},{"id":5,"name":"sub2","sub_department":[],"created_at":"2021-01-21 15:11:07","updated_at":"2021-01-21 15:11:07"},{"id":6,"name":"sub3","sub_department":[{"id":9,"name":"sub3_1","sub_department":[],"created_at":"2021-02-19 03:17:30","updated_at":"2021-02-19 03:17:33"}],"created_at":"2021-01-21 15:11:10","updated_at":"2021-01-21 15:11:10"},{"id":8,"name":"sub1_1","sub_department":[],"created_at":"2021-01-21 15:11:34","updated_at":"2021-01-22 10:53:13"}],"created_at":"2021-01-21 15:10:26","updated_at":"2021-01-21 15:10:26"}}"
// @Router /departments/tree [GET]
func ListDepartmentsTree(c *gin.Context) {
	// 查找根部门
	dept, ok := model.GetDepartment(model.Query{"parent_id": nil})
	if !ok {
		resp := msg.ErrDatabase.GenResponse("Error when GetDepartment.")
		log.Error(resp.String())
		resp.Write(c)
		return
	}
	// 找不到根部门则直接返回空
	if dept == nil {
		resp := msg.Success.GenResponse().SetPayload("tree", make(map[string]interface{}))
		resp.Write(c)
		return
	}

	// 列出所有部门
	deptList, ok := model.ListDepartments(model.Query{})
	if !ok {
		resp := msg.ErrDatabase.GenResponse("Error when ListDepartments.")
		log.Error(resp.String())
		resp.Write(c)
		return
	}

	// 将根部门转化为树结构
	root := model.Department2Tree(dept)
	model.ListDepartment2Tree(root, deptList)

	resp := msg.Success.GenResponse().SetPayload("tree", root)
	resp.Write(c)
}

// ListDepartmentTree 列出指定部门子树
// @Summary 列出指定部门子树
// @Tags 部门
// @Param token header string true "Token"
// @Param department_id path string true "部门ID"
// @Success 200 {string} string "{"code":0,"message":"Success","tree":{"id":6,"name":"sub3","sub_department":[{"id":9,"name":"sub3_1","sub_department":[],"created_at":"2021-02-19 03:17:30","updated_at":"2021-02-19 03:17:33"}],"created_at":"2021-01-21 15:11:10","updated_at":"2021-01-21 15:11:10"}}"
// @Router /departments/{department_id}/tree [GET]
func ListDepartmentTree(c *gin.Context) {
	deptId, err := strconv.Atoi(c.Param("department_id"))
	if err != nil {
		resp := msg.WarnInvalidUri.GenResponse("Field 'department_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())
		resp.Write(c)
		return
	}

	// 查找根部门
	dept, ok := model.GetDepartment(model.Query{"id=?": deptId})
	if !ok {
		resp := msg.ErrDatabase.GenResponse("Error when GetDepartment.")
		log.Error(resp.String())
		resp.Write(c)
		return
	}
	// 找不到根部门则直接返回空
	if dept == nil {
		resp := msg.Success.GenResponse().SetPayload("tree", make(map[string]interface{}))
		resp.Write(c)
		return
	}

	// 列出所有部门
	deptList, ok := model.ListDepartments(model.Query{})
	if !ok {
		resp := msg.ErrDatabase.GenResponse("Error when ListDepartments.")
		log.Error(resp.String())
		resp.Write(c)
		return
	}

	// 将根部门转化为树结构
	root := model.Department2Tree(dept)
	model.ListDepartment2Tree(root, deptList)

	resp := msg.Success.GenResponse().SetPayload("tree", root)
	resp.Write(c)
}

// AddUser2Department 添加用户至部门
// @Summary 添加用户至部门
// @Tags 部门
// @Param token header string true "Token"
// @Param department_id path string true "部门ID"
// @Param data body model.UserDepartment true "body"
// @Success 200 {string} string "{"code":0,"message":"Success"}"
// @Router /departments/{department_id}/users [POST]
func AddUser2Department(c *gin.Context) {
	var uDept model.UserDepartment

	deptId, err := strconv.Atoi(c.Param("department_id"))
	if err != nil {
		resp := msg.WarnInvalidUri.GenResponse("Field 'department_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())
		resp.Write(c)
		return
	}

	// 序列化request body
	if err := c.ShouldBindJSON(&uDept); err != nil {
		resp := msg.WarnInvalidBody.GenResponse("Field 'user_id', 'is_owner' required, and 'user_id' must greater than 0.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())
		resp.Write(c)
		return
	}

	if !model.AddDepartmentUser(uDept.UserId, deptId, *uDept.IsOwner) {
		resp := msg.ErrDatabase.GenResponse("Error when AddDepartmentUser.")
		log.Error(resp.String())
		resp.Write(c)
		return
	}

	resp := msg.Success.GenResponse()
	resp.Write(c)
}

// RemoveDepartmentUser 移除部门中用户
// @Summary 移除部门中用户
// @Tags 部门
// @Param token header string true "Token"
// @Param department_id path string true "部门ID"
// @Param user_id path string true "用户ID"
// @Success 200 {string} string "{"code":0,"message":"Success"}"
// @Router /departments/{department_id}/users/{user_id} [DELETE]
func RemoveDepartmentUser(c *gin.Context) {
	deptId, err := strconv.Atoi(c.Param("department_id"))
	if err != nil {
		resp := msg.WarnInvalidUri.GenResponse("Field 'department_id' required.")
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

	if !model.RemoveDepartmentUser(userId, deptId) {
		resp := msg.ErrDatabase.GenResponse("Error when RemoveDepartmentUser.")
		log.Error(resp.String())
		resp.Write(c)
		return
	}

	resp := msg.Success.GenResponse()
	resp.Write(c)
}

// SetUserIsDepartmentOwner 设置用户是否是部门Owner
// @Summary 设置用户是否是部门Owner
// @Tags 部门
// @Param token header string true "Token"
// @Param department_id path string true "部门ID"
// @Param user_id path string true "用户ID"
// @Param data body form.IsOwnerForm true "body"
// @Success 200 {string} string "{"code":0,"message":"Success"}"
// @Router /departments/{department_id}/users/{user_id}  [PUT]
func SetUserIsDepartmentOwner(c *gin.Context) {
	var fm = form.IsOwnerForm{}

	deptId, err := strconv.Atoi(c.Param("department_id"))
	if err != nil {
		resp := msg.WarnInvalidUri.GenResponse("Field 'department_id' required.")
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

	// 序列化request body
	if err := c.ShouldBindJSON(&fm); err != nil {
		resp := msg.WarnInvalidBody.GenResponse("Field 'is_owner' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, resp.String())
		resp.Write(c)
		return
	}

	if !model.SetDepartmentUserIsOwner(userId, deptId, *fm.IsOwner) {
		resp := msg.ErrDatabase.GenResponse("Error when SetDepartmentUserIsOwner.")
		log.Error(resp.String())
		resp.Write(c)
		return
	}

	resp := msg.Success.GenResponse()
	resp.Write(c)
}

// ListProductUsers 列出部门中所有用户
// @Summary 列出部门中所有用户
// @Tags 部门
// @Param token header string true "Token"
// @Param department_id path string true "部门ID"
// @Success 200 {string} string "{"code":0,"message":"Success","users":[{"id":4,"username":"test2","is_owner":false,"joined_at":"2021-01-20 11:01:16"},{"id":3,"username":"test","is_owner":true,"joined_at":"2021-01-20 11:01:32"}]}"
// @Router /departments/{department_id}/users [GET]
func ListDepartmentUsers(c *gin.Context) {
	var query = model.Query{}

	deptId, err := strconv.Atoi(c.Param("department_id"))
	if err != nil {
		resp := msg.WarnInvalidUri.GenResponse("Field 'department_id' required.")
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

	u, ok := model.ListDepartmentUsers(deptId, query)
	if !ok {
		resp := msg.ErrDatabase.GenResponse("Error when ListDepartmentUsers.")
		log.Error(resp.String())
		resp.Write(c)
		return
	}

	resp := msg.Success.GenResponse().SetPayload("users", u)
	resp.Write(c)
}
