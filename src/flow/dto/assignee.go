package dto

// AssigneesSpiltTag 审批人分隔符
const AssigneesSpiltTag = ","

// AssigneeCondition 取值枚举范围
const (
	AssigneeConditionInitiator                        = "initiator"                          // 发起人
	AssigneeConditionInitiatorsDepartmentsOwner       = "initiators_department_owner"        // 发起人所属部门Owner
	AssigneeConditionInitiatorsParentDepartmentsOwner = "initiators_parent_department_owner" // 发起人所属上级部门Owner
	AssigneeConditionSpecifiedUsers                   = "specified_users"                    // 指定用户（多个）
	AssigneeConditionSpecifiedProductOwner            = "specified_product_owner"            // 指定产品线Owner
	AssigneeConditionSpecifiedGroupOwner              = "specified_group_owner"              // 指定组Owner
	AssigneeConditionSpecifiedDepartmentOwner         = "specified_department_owner"         // 指定部门Owner
	AssigneeConditionSpecifiedRole                    = "specified_role"                     // 指定角色
)

var AssigneeConditionsAllowed = map[string]struct{}{
	AssigneeConditionInitiator:                        activeEmptyStruct,
	AssigneeConditionInitiatorsDepartmentsOwner:       activeEmptyStruct,
	AssigneeConditionInitiatorsParentDepartmentsOwner: activeEmptyStruct,
	AssigneeConditionSpecifiedUsers:                   activeEmptyStruct,
	AssigneeConditionSpecifiedProductOwner:            activeEmptyStruct,
	AssigneeConditionSpecifiedGroupOwner:              activeEmptyStruct,
	AssigneeConditionSpecifiedDepartmentOwner:         activeEmptyStruct,
	AssigneeConditionSpecifiedRole:                    activeEmptyStruct,
}

type AssigneeCondition struct {
	Condition string `json:"condition"`
	Getter    string `json:"getter"`
	Data      string `json:"data"`
}
