package dto

type ProductInToken struct {
	Id       int    `json:"id"`
	Name     string `json:"name" binding:"required"`
	Alias    string `json:"alias" binding:"required"`
	Disabled bool   `json:"disabled"`
}

type GroupInToken struct {
	Id   int    `json:"id"`
	Name string `json:"name" binding:"required"`
}

type TokenContent struct {
	UserId   int    `json:"user_id"`
	Username string `json:"username"`
	Phone    string `json:"phone"`

	IsSuperuser bool `json:"is_superuser"`

	Department  *[]string         `json:"department"`
	Roles       *[]string         `json:"roles"`
	Products    *[]ProductInToken `json:"products"`
	OwnProducts *[]ProductInToken `json:"own_products"`
	Groups      *[]GroupInToken   `json:"groups"`
	OwnGroups   *[]GroupInToken   `json:"own_groups"`
}
