package model

import (
	"eago/common/utils"
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

type ListNodes struct {
	Id                int              `json:"id"`
	Name              string           `json:"name"`
	ParentId          *int             `json:"parent_id"`
	ParentName        string           `json:"parent_name"`
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
