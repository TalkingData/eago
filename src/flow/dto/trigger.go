package dto

type TriggersNode struct {
	Id       uint32  `json:"id"`
	Name     string  `json:"name"`
	ParentId *uint32 `json:"parent_id"`
}
