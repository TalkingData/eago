package dto

import (
	"eago/auth/model"
	"eago/common/utils"
)

type DepartmentTree struct {
	Id            uint32            `json:"id"`
	Name          string            `json:"name"`
	SubDepartment []*DepartmentTree `json:"sub_department"`
	CreatedAt     *utils.CustomTime `json:"created_at"`
	UpdatedAt     *utils.CustomTime `json:"updated_at"`
}

// TransDepartment2Tree 将部门转化为树结构的一个节点
func TransDepartment2Tree(dept *model.Department) *DepartmentTree {
	return &DepartmentTree{
		Id:            dept.Id,
		Name:          dept.Name,
		CreatedAt:     dept.CreatedAt,
		UpdatedAt:     dept.UpdatedAt,
		SubDepartment: make([]*DepartmentTree, 0),
	}
}
