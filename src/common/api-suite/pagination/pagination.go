package pagination

const (
	DEFAULT_PAGE_SIZE = 50
	MAX_PAGE_SIZE     = 500
)

type Paginator struct {
	Page     int         `json:"page"`
	Pages    int         `json:"pages"`
	PageSize int         `json:"page_size"`
	Total    int64       `json:"total"`
	Data     interface{} `json:"data"`
}
