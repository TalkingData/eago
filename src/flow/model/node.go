package model

import (
	"eago/common/utils"
)

// AssigneeConditions 取值枚举范围
const (
	INITIATOR                           = "initiator"                          // 发起人
	INITIATORS_DEPARTMENTS_OWNER        = "initiators_department_owner"        // 发起人所属部门Owner
	INITIATORS_PARENT_DEPARTMENTS_OWNER = "initiators_parent_department_owner" // 发起人所属上级部门Owner
	SPECIFIED_USERS                     = "specified_users"                    // 指定用户（多个）
	SPECIFIED_PRODUCT_OWNER             = "specified_product_owner"            // 指定产品线Owner
	SPECIFIED_GROUP_OWNER               = "specified_group_owner"              // 指定组Owner
	SPECIFIED_DEPARTMENT_OWNER          = "specified_department_owner"         // 指定部门Owner
	SPECIFIED_ROLE                      = "specified_role"                     // 指定角色
)

// Category 取值枚举范围
const (
	ANY    = iota // 或签
	ALL           // 会签
	INFORM        // 知会
)

type Node struct {
	Id                int              `json:"id"`
	Name              string           `json:"name"`
	ParentId          *int             `json:"parent_id"`
	Category          int              `json:"category"`
	EntryCondition    *string          `json:"entry_condition" gorm:"default:'{}'"`
	AssigneeCondition *string          `json:"assignee_condition" gorm:"default:'{}'"`
	VisibleFields     string           `json:"visible_fields"`
	EditableFields    string           `json:"editable_fields"`
	CreatedAt         *utils.LocalTime `json:"created_at"`
	CreatedBy         string           `json:"created_by"`
	UpdatedAt         *utils.LocalTime `json:"updated_at"`
	UpdatedBy         *string          `json:"updated_by" gorm:"default:''"`
}

type NodeChain struct {
	Id                int            `json:"id"`
	Name              string         `json:"name"`
	ParentId          *int           `json:"parent_id"`
	Category          int            `json:"category"`
	EntryCondition    string         `json:"entry_condition"`
	AssigneeCondition string         `json:"assignee_condition"`
	Assignees         []string       `json:"assignees"`
	Triggers          []NodesTrigger `json:"triggers"`
	VisibleFields     string         `json:"visible_fields"`
	EditableFields    string         `json:"editable_fields"`
	SubNode           *NodeChain     `json:"sub_node"`
}

type NodeTrigger struct {
	Id        int              `json:"id"`
	NodeId    int              `json:"node_id"`
	TriggerId int              `json:"trigger_id"`
	CreatedAt *utils.LocalTime `json:"joined_at" gorm:"type:datetime;not null;autoCreateTime"`
	CreatedBy string           `json:"created_by"`
}

type NodesTrigger struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	TaskCodename string `json:"task_codename"`
	Arguments    string `json:"arguments"`
}
