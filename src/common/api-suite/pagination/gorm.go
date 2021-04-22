package pagination

import (
	"eago-common/tools"
	"errors"
	"gorm.io/gorm"
	"math"
)

type GormParams struct {
	Db       *gorm.DB
	Page     int
	PageSize int
	OrderBy  []string
}

// GormPaging gorm分页处理
func GormPaging(p *GormParams, result interface{}) (*Paginator, error) {
	var (
		paginator Paginator
		offset    int
	)

	db := p.Db

	// 设置Page
	if p.Page < 1 {
		p.Page = 1
		offset = 0
	} else {
		offset = (p.Page - 1) * p.PageSize
	}

	// 设置PageSize
	if p.PageSize < 1 || p.PageSize > MAX_PAGE_SIZE {
		p.PageSize = DEFAULT_PAGE_SIZE
	}
	paginator.PageSize = p.PageSize

	// 处理OrderBy
	if len(p.OrderBy) > 0 {
		for _, o := range p.OrderBy {
			if o == "" {
				continue
			}
			db = db.Order(o)
			if db.Error != nil {
				return nil, db.Error
			}
		}
	}

	// 取得总数
	res := db.Count(&paginator.Total)
	if res.Error != nil {
		return nil, db.Error
	}

	// 分页查询
	res = db.Limit(p.PageSize).Offset(offset).Find(result)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			paginator.Data = result
			paginator.Pages = 0
			return &paginator, nil
		}
		return nil, db.Error
	}
	paginator.Data = result
	paginator.Pages = int(math.Ceil(float64(paginator.Total) / float64(p.PageSize)))
	paginator.Page = tools.IntMin(p.Page, paginator.Pages)

	return &paginator, nil
}
