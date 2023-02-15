package menu

type Button struct {
	Id   string    `json:"id"`
	Name string    `json:"name"`
	Perm *itemPerm `json:"-"`
}
