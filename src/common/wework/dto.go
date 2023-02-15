package wework

type WeworkUsersDepartment struct {
	Username     string
	DisplayName  string
	DepartmentId int
}

type WeworkDepartment struct {
	Id       int
	Name     string
	ParentId int
}
