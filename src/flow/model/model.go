package model

type AssigneeCondition struct {
	Condition string `json:"condition"`
	Getter    string `json:"getter"`
	Data      string `json:"data"`
}
