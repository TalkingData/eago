package handler

import (
	"database/sql"
	"eago-auth/api/form"
	"eago-auth/config/msg"
	db "eago-auth/database"
	"eago-common/log"
	"eago-common/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// NewDepartment 新建部门
// @Summary 新建部门
// @Tags 部门
// @Param token header string true "Token"
// @Param data body db.Department true "body"
// @Success 200 {string} string "{"code":0,"message":"Success","department":{"id":4,"name":"sub_dept3","parent_id":2,"created_at":"2021-01-21 15:11:00","updated_at":"2021-01-21 15:11:00"}}"
// @Router /departments [POST]
func NewDepartment(c *gin.Context) {
	var dForm db.Department

	// 序列化request body
	if err := c.ShouldBindJSON(&dForm); err != nil {
		m := msg.WarnInvalidBody.NewMsg("Field 'name', 'parent_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	dept := db.DepartmentModel.New(dForm.Name, dForm.ParentId)
	if dept == nil {
		m := msg.ErrDatabase.NewMsg("Error in db.DepartmentModel.NewProduct.")
		log.Error(m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	m := msg.Success.NewMsg().SetPayload(&gin.H{"department": dept})
	c.JSON(http.StatusOK, m.GinH())
}

// DeleteDepartment 删除部门
// @Summary 删除部门
// @Tags 部门
// @Param token header string true "Token"
// @Param department_id path string true "部门ID"
// @Success 200 {string} string "{"code":0,"message":"Success"}"
// @Router /departments/{department_id} [DELETE]
func DeleteDepartment(c *gin.Context) {
	deptId, err := strconv.Atoi(c.Param("department_id"))
	if err != nil {
		m := msg.WarnInvalidUri.NewMsg("Field 'department_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	if suc := db.DepartmentModel.Delete(deptId); !suc {
		m := msg.ErrDatabase.NewMsg("Error in db.DepartmentModel.Delete.")
		log.Error(m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	m := msg.Success.NewMsg()
	c.JSON(http.StatusOK, m.GinH())
}

// SetDepartment 更新部门
// @Summary 更新部门
// @Tags 部门
// @Param token header string true "Token"
// @Param department_id path string true "部门ID"
// @Param data body db.Department true "body"
// @Success 200 {string} string "{"code":0,"message":"Success","department":{"id":4,"name":"sub_dept3","parent_id":2,"created_at":"2021-01-21 15:11:00","updated_at":"2021-01-21 15:11:00"}}"
// @Router /departments [PUT]
func SetDepartment(c *gin.Context) {
	var dForm db.Department

	deptId, err := strconv.Atoi(c.Param("department_id"))
	if err != nil {
		m := msg.WarnInvalidUri.NewMsg("Field 'department_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	// 序列化request body
	if err := c.ShouldBindJSON(&dForm); err != nil {
		m := msg.WarnInvalidBody.NewMsg("Field 'name', 'parent_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	dept, suc := db.DepartmentModel.Set(deptId, dForm.Name, dForm.ParentId)
	if !suc {
		m := msg.ErrDatabase.NewMsg("Error in db.DepartmentModel.Set.")
		log.Error(m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	m := msg.Success.NewMsg().SetPayload(&gin.H{"department": dept})
	c.JSON(http.StatusOK, m.GinH())
}

// ListDepartments 列出所有部门
// @Summary 列出所有部门
// @Tags 部门
// @Param token header string true "Token"
// @Param query query string false "过滤条件"
// @Param page query string false "页数"
// @Param page_size query string false "页尺寸"
// @Success 200 {string} string "{"code":0,"departments":[{"id":2,"name":"root","parent_id":null,"created_at":"2021-01-21 15:10:26","updated_at":"2021-01-21 15:10:26"},{"id":5,"name":"sub2","parent_id":2,"created_at":"2021-01-21 15:11:07","updated_at":"2021-01-21 15:11:07"}],"message":"Success","page":1,"page_size":50,"pages":1,"total":2}"
// @Router /departments [GET]
func ListDepartments(c *gin.Context) {
	var query db.Query

	q := c.GetString("Query")
	if q != "" {
		likeQuery := fmt.Sprintf("%%%s%%", q)
		query = db.Query{"name LIKE @query OR alias LIKE @query id LIKE @query": sql.Named("query", likeQuery)}
	}

	paged, suc := db.DepartmentModel.PagedList(
		&query,
		c.GetInt("Page"),
		c.GetInt("PageSize"),
		c.GetStringSlice("OrderBy")...,
	)
	if !suc {
		m := msg.ErrDatabase.NewMsg("Error in db.DepartmentModel.PageList.")
		log.Error(m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	m := msg.Success.NewMsg().SetPagedPayload(paged, "departments")
	c.JSON(http.StatusOK, m.GinH())
}

// ListDepartmentsTree 以树结构列出所有部门
// @Summary 以树结构列出所有部门
// @Tags 部门
// @Param token header string true "Token"
// @Success 200 {string} string "{"code":0,"message":"Success","tree":{"id":2,"name":"root","sub_department":[{"id":4,"name":"sub1","sub_department":[],"created_at":"2021-01-21 15:11:00","updated_at":"2021-01-21 15:11:00"},{"id":5,"name":"sub2","sub_department":[],"created_at":"2021-01-21 15:11:07","updated_at":"2021-01-21 15:11:07"},{"id":6,"name":"sub3","sub_department":[{"id":9,"name":"sub3_1","sub_department":[],"created_at":"2021-02-19 03:17:30","updated_at":"2021-02-19 03:17:33"}],"created_at":"2021-01-21 15:11:10","updated_at":"2021-01-21 15:11:10"},{"id":8,"name":"sub1_1","sub_department":[],"created_at":"2021-01-21 15:11:34","updated_at":"2021-01-22 10:53:13"}],"created_at":"2021-01-21 15:10:26","updated_at":"2021-01-21 15:10:26"}}"
// @Router /departments/tree [GET]
func ListDepartmentsTree(c *gin.Context) {
	var root *db.DepartmentTree

	// 查找根部门
	dept, suc := db.DepartmentModel.Get(&db.Query{"parent_id": nil})
	if !suc {
		m := msg.ErrDatabase.NewMsg("Error in db.DepartmentModel.Get.")
		log.Error(m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}
	// 找不到根部门则直接返回空
	if dept == nil {
		m := msg.Success.NewMsg()
		m = m.SetPayload(&gin.H{"tree": gin.H{}})
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	// 列出所有部门
	deptList, suc := db.DepartmentModel.List(&db.Query{})
	if !suc {
		m := msg.ErrDatabase.NewMsg("Error in db.DepartmentModel.List.")
		log.Error(m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	// 将根部门转化为树结构
	root = db.DepartmentModel.Department2Tree(dept)
	db.DepartmentModel.List2Tree(root, deptList)

	m := msg.Success.NewMsg()
	m = m.SetPayload(&gin.H{"tree": root})
	c.JSON(http.StatusOK, m.GinH())
}

// ListDepartmentTree 列出指定部门子树
// @Summary 列出指定部门子树
// @Tags 部门
// @Param token header string true "Token"
// @Param department_id path string true "部门ID"
// @Success 200 {string} string "{"code":0,"message":"Success","tree":{"id":6,"name":"sub3","sub_department":[{"id":9,"name":"sub3_1","sub_department":[],"created_at":"2021-02-19 03:17:30","updated_at":"2021-02-19 03:17:33"}],"created_at":"2021-01-21 15:11:10","updated_at":"2021-01-21 15:11:10"}}"
// @Router /departments/{department_id}/tree [GET]
func ListDepartmentTree(c *gin.Context) {
	var root *db.DepartmentTree

	deptId, err := strconv.Atoi(c.Param("department_id"))
	if err != nil {
		m := msg.WarnInvalidUri.NewMsg("Field 'department_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	// 查找根部门
	dept, suc := db.DepartmentModel.Get(&db.Query{"id=?": deptId})
	if !suc {
		m := msg.ErrDatabase.NewMsg("Error in db.DepartmentModel.Get.")
		log.Error(m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}
	// 找不到根部门则直接返回空
	if dept == nil {
		m := msg.Success.NewMsg()
		m = m.SetPayload(&gin.H{"tree": gin.H{}})
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	// 列出所有部门
	deptList, suc := db.DepartmentModel.List(&db.Query{})
	if !suc {
		m := msg.ErrDatabase.NewMsg("Error in db.DepartmentModel.List.")
		log.Error(m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	// 将根部门转化为树结构
	root = db.DepartmentModel.Department2Tree(dept)
	db.DepartmentModel.List2Tree(root, deptList)

	m := msg.Success.NewMsg()
	m = m.SetPayload(&gin.H{"tree": root})
	c.JSON(http.StatusOK, m.GinH())
}

// AddUser2Department 添加用户至部门
// @Summary 添加用户至部门
// @Tags 部门
// @Param token header string true "Token"
// @Param department_id path string true "部门ID"
// @Param data body db.UserDepartment true "body"
// @Success 200 {string} string "{"code":0,"message":"Success"}"
// @Router /departments/{department_id}/users [POST]
func AddUser2Department(c *gin.Context) {
	var uDept db.UserDepartment

	deptId, err := strconv.Atoi(c.Param("department_id"))
	if err != nil {
		m := msg.WarnInvalidUri.NewMsg("Field 'department_id' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	// 序列化request body
	if err := c.ShouldBindJSON(&uDept); err != nil {
		m := msg.WarnInvalidBody.NewMsg("Field 'user_id', 'is_owner' required.")
		log.WarnWithFields(log.Fields{
			"error": err.Error(),
		}, m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	if !db.DepartmentModel.AddUser(uDept.UserId, deptId, *uDept.IsOwner) {
		m := msg.ErrDatabase.NewMsg("Error in db.DepartmentModel.AddUser.")
		log.Error(m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	m := msg.Success.NewMsg()
	c.JSON(http.StatusOK, m.GinH())
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
		m := msg.WarnInvalidUri.NewMsg("Field 'department_id' required.")
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

	if !db.DepartmentModel.RemoveUser(userId, deptId) {
		m := msg.ErrDatabase.NewMsg("Error in db.DepartmentModel.RemoveUser.")
		log.Error(m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	m := msg.Success.NewMsg()
	c.JSON(http.StatusOK, m.GinH())
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
		m := msg.WarnInvalidUri.NewMsg("Field 'department_id' required.")
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

	if !db.DepartmentModel.SetUserIsOwner(userId, deptId, *fm.IsOwner) {
		m := msg.ErrDatabase.NewMsg("Error in db.DepartmentModel.SetUserIsOwner.")
		log.Error(m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	m := msg.Success.NewMsg()
	c.JSON(http.StatusOK, m.GinH())
}

// ListProductUsers 列出部门中所有用户
// @Summary 列出部门中所有用户
// @Tags 部门
// @Param token header string true "Token"
// @Param department_id path string true "部门ID"
// @Success 200 {string} string "{"code":0,"message":"Success","users":[{"id":4,"username":"test2","is_owner":false,"joined_at":"2021-01-20 11:01:16"},{"id":3,"username":"test","is_owner":true,"joined_at":"2021-01-20 11:01:32"}]}"
// @Router /departments/{department_id}/users [GET]
func ListDepartmentUsers(c *gin.Context) {
	var query = db.Query{}

	deptId, err := strconv.Atoi(c.Param("department_id"))
	if err != nil {
		m := msg.WarnInvalidUri.NewMsg("Field 'department_id' required.")
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

	u, suc := db.DepartmentModel.ListUsers(deptId, &query)
	if !suc {
		m := msg.ErrDatabase.NewMsg("Error in db.DepartmentModel.ListUsers.")
		log.Error(m.String())
		c.JSON(http.StatusOK, m.GinH())
		return
	}

	m := msg.Success.NewMsg().SetPayload(&gin.H{"users": u})
	c.JSON(http.StatusOK, m.GinH())
}
