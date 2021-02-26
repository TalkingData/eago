package form

type IsOwnerForm struct {
	IsOwner *bool `json:"is_owner" binding:"required"`
}

type LoginForm struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
