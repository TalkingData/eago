package eagle

type Product struct {
	Name        string `json:"name"`
	Alias       string `json:"alias"`
	Description string `json:"description"`
}

type ProductMember struct {
	Username string `json:"username"`
	Product  string `json:"product"`
}
