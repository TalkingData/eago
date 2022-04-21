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

// NewDepartment 新建部门
func NewDepartment(c *gin.Context) {
	var deptFrm dto.NewDepartment
	// 序列化request body
	if err := c.ShouldBindJSON(&deptFrm); err != nil {
		m := msg.SerializeFailed.SetError(err)
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}
	// 验证数据
	if m := deptFrm.Validate(); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	// 新建
	dept, err := dao.NewDepartment(deptFrm.Name, deptFrm.ParentId)
	// 新建失败
	if dept == nil {
		m := msg.UndefinedError.SetError(err)
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "department", dept)
}

// RemoveDepartment 删除部门
func RemoveDepartment(c *gin.Context) {
	deptId, err := strconv.Atoi(c.Param("department_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "department_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	var rdFrm dto.RemoveDepartment
	// 验证数据
	if m := rdFrm.Validate(deptId); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	if !dao.RemoveDepartment(deptId) {
		m := msg.UndefinedError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccess(c)
}

// SetDepartment 更新部门
func SetDepartment(c *gin.Context) {
	deptId, err := strconv.Atoi(c.Param("department_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "department_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	var deptFrm dto.SetDepartment
	// 序列化request body
	if err = c.ShouldBindJSON(&deptFrm); err != nil {
		m := msg.SerializeFailed.SetError(err)
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}
	// 验证数据
	if m := deptFrm.Validate(deptId); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	dept, err := dao.SetDepartment(deptId, deptFrm.Name, deptFrm.ParentId)
	if err != nil {
		m := msg.UndefinedError.SetError(err)
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "department", dept)
}

// PagedListDepartments 列出所有部门-分页
func PagedListDepartments(c *gin.Context) {
	query := dao.Query{}
	// 设置查询filter
	ldq := dto.PagedListDepartmentsQuery{}
	if c.ShouldBindQuery(&ldq) == nil {
		_ = ldq.UpdateQuery(query)
	}

	paged, ok := dao.PagedListDepartments(
		query,
		c.GetInt("Page"),
		c.GetInt("PageSize"),
		c.GetStringSlice("OrderBy")...,
	)
	if !ok {
		m := msg.UndefinedError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "departments", paged)
}

// ListDepartmentsTree 以树结构列出所有部门
func ListDepartmentsTree(c *gin.Context) {
	// 查找根部门
	dept, ok := dao.GetDepartment(dao.Query{"parent_id": nil})
	if !ok {
		m := msg.UndefinedError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}
	// 找不到根部门则直接返回空
	if dept == nil {
		w.WriteSuccessPayload(c, "tree", make(map[string]interface{}))
		return
	}

	// 列出所有部门
	deptList, ok := dao.ListDepartments(dao.Query{})
	if !ok {
		m := msg.UndefinedError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	// 将根部门转化为树结构
	root := dao.Department2Tree(dept)
	dao.ListDepartment2Tree(root, deptList)

	w.WriteSuccessPayload(c, "tree", root)
}

// ListDepartmentTree 列出指定部门子树
func ListDepartmentTree(c *gin.Context) {
	deptId, err := strconv.Atoi(c.Param("department_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "department_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	// 查找根部门
	dept, ok := dao.GetDepartment(dao.Query{"id=?": deptId})
	if !ok {
		m := msg.UndefinedError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}
	// 找不到根部门则直接返回空
	if dept == nil {
		w.WriteSuccessPayload(c, "tree", make(map[string]interface{}))
		return
	}

	// 列出所有部门
	deptList, ok := dao.ListDepartments(dao.Query{})
	if !ok {
		m := msg.UndefinedError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	// 将根部门转化为树结构
	root := dao.Department2Tree(dept)
	dao.ListDepartment2Tree(root, deptList)

	w.WriteSuccessPayload(c, "tree", root)
}

// AddUser2Department 添加用户至部门
func AddUser2Department(c *gin.Context) {
	deptId, err := strconv.Atoi(c.Param("department_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "department_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	var audFrm dto.AddUser2Department
	// 序列化request body
	if err = c.ShouldBindJSON(&audFrm); err != nil {
		m := msg.SerializeFailed
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	// 验证数据
	if m := audFrm.Validate(deptId); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	if !dao.AddDepartmentUser(audFrm.UserId, deptId, audFrm.IsOwner) {
		m := msg.UndefinedError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccess(c)
}

// RemoveDepartmentUser 移除部门中用户
func RemoveDepartmentUser(c *gin.Context) {
	deptId, err := strconv.Atoi(c.Param("department_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "department_id")
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

	var rduFrm dto.RemoveDepartmentUser
	// 验证数据
	if m := rduFrm.Validate(deptId, userId); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	if !dao.RemoveDepartmentUser(userId, deptId) {
		m := msg.UndefinedError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccess(c)
}

// SetUserIsDepartmentOwner 设置用户是否是部门Owner
func SetUserIsDepartmentOwner(c *gin.Context) {
	deptId, err := strconv.Atoi(c.Param("department_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "department_id")
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

	var suoFrm dto.SetUserIsDepartmentOwner
	// 序列化request body
	if err = c.ShouldBindJSON(&suoFrm); err != nil {
		m := msg.SerializeFailed.SetError(err)
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}
	// 验证数据
	if m := suoFrm.Validate(deptId, userId); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	if !dao.SetDepartmentUserIsOwner(deptId, userId, suoFrm.IsOwner) {
		m := msg.UndefinedError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccess(c)
}

// ListDepartmentUsers 列出部门中所有用户
func ListDepartmentUsers(c *gin.Context) {
	deptId, err := strconv.Atoi(c.Param("department_id"))
	if err != nil {
		m := msg.InvalidUriFailed.SetError(err, "department_id")
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	var lduFrm dto.ListDepartmentUsersQuery
	// 序列化request body
	_ = c.ShouldBindQuery(&lduFrm)
	if m := lduFrm.Validate(deptId); m != nil {
		// 数据验证未通过
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	query := dao.Query{}
	// 设置查询filter
	_ = lduFrm.UpdateQuery(query)

	u, ok := dao.ListDepartmentUsers(deptId, query)
	if !ok {
		m := msg.UndefinedError
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	w.WriteSuccessPayload(c, "users", u)
}
