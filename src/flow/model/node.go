package model

import (
	"eago/common/utils"
)

type Node struct {
	Id uint32 `json:"id"`

	Name     string  `json:"name"`
	ParentId *uint32 `json:"parent_id"`
	Category int32   `json:"category"`

	EntryCondition    *string `json:"entry_condition" gorm:"default:'{}'"`
	AssigneeCondition *string `json:"assignee_condition" gorm:"default:'{}'"`

	VisibleFields  string `json:"visible_fields"`
	EditableFields string `json:"editable_fields"`

	CreatedAt *utils.CustomTime `json:"created_at"`
	CreatedBy string            `json:"created_by"`
	UpdatedAt *utils.CustomTime `json:"updated_at"`
	UpdatedBy *string           `json:"updated_by" gorm:"default:''"`
}

type ListNodes struct {
	Id uint32 `json:"id"`

	Name       string  `json:"name"`
	ParentId   *uint32 `json:"parent_id"`
	ParentName string  `json:"parent_name"`
	Category   int32   `json:"category"`

	EntryCondition    *string `json:"entry_condition" gorm:"default:'{}'"`
	AssigneeCondition *string `json:"assignee_condition" gorm:"default:'{}'"`

	VisibleFields  string `json:"visible_fields"`
	EditableFields string `json:"editable_fields"`

	CreatedAt *utils.CustomTime `json:"created_at"`
	CreatedBy string            `json:"created_by"`
	UpdatedAt *utils.CustomTime `json:"updated_at"`
	UpdatedBy *string           `json:"updated_by" gorm:"default:''"`
}

type NodeChain struct {
	Id uint32 `json:"id"`

	Name     string  `json:"name"`
	ParentId *uint32 `json:"parent_id"`
	Category int32   `json:"category"`

	EntryCondition    string `json:"entry_condition"`
	AssigneeCondition string `json:"assignee_condition"`

	Assignees []string        `json:"assignees"`
	Triggers  []*NodeTriggers `json:"triggers"`

	VisibleFields  string `json:"visible_fields"`
	EditableFields string `json:"editable_fields"`

	SubNode *NodeChain `json:"sub_node"`
}

type NodeTrigger struct {
	Id uint32 `json:"id"`

	NodeId    uint32 `json:"node_id"`
	TriggerId uint32 `json:"trigger_id"`

	CreatedAt *utils.CustomTime `json:"joined_at" gorm:"type:datetime;not null;autoCreateTime"`
	CreatedBy string            `json:"created_by"`
}

type NodeTriggers struct {
	Id uint32 `json:"id"`

	Name         string `json:"name"`
	Description  string `json:"description"`
	TaskCodename string `json:"task_codename"`
	Arguments    string `json:"arguments"`
}
