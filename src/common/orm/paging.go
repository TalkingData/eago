package orm

import (
	"eago/common/global"
	"eago/common/utils"
	"gorm.io/gorm"
	"math"
	"sync"
)

// Paginator 分页器结构
type Paginator struct {
	Page     int         `json:"page"`
	Pages    int         `json:"pages"`
	PageSize int         `json:"page_size"`
	Total    int64       `json:"total"`
	Data     interface{} `json:"data"`
	Offset   int
}

func PagingQuery(db *gorm.DB, _page, _pageSize int, _result interface{}, _orderBy ...string) (*Paginator, error) {
	// 处理OrderBy
	for _, o := range _orderBy {
		if len(o) < 1 {
			continue
		}
		db = db.Order(o)
		if db.Error != nil {
			return nil, db.Error
		}
	}

	paginator := &Paginator{}
	// 设置PageSize，每页的对象数量
	paginator.PageSize = utils.IntMin(utils.IntMax(_pageSize, 1), global.MaxPageSize)
	// 设置Offset，查询偏移量
	paginator.Offset = (paginator.Page - 1) * paginator.PageSize

	// 分页查询
	res := db.Session(&gorm.Session{}).
		Limit(paginator.PageSize).Offset(paginator.Offset).
		Count(&paginator.Total).
		Find(_result)
	if res.Error != nil {
		return nil, res.Error
	}

	// 设置Pages，总共有多少页
	paginator.Pages = int(math.Ceil(float64(paginator.Total) / float64(paginator.PageSize)))
	// 设置Page，当前查询第几页
	paginator.Page = utils.IntMin(utils.IntMax(_page, 1), paginator.Pages)
	// 设置Data，查询结果
	paginator.Data = _result

	return paginator, nil
}

// PagingQueryCoroutine 分页查询器并发版，测试中
func PagingQueryCoroutine(
	db *gorm.DB, _page, _pageSize int, _result interface{}, _orderBy ...string,
) (*Paginator, error) {
	// 处理OrderBy
	for _, o := range _orderBy {
		if len(o) < 1 {
			continue
		}
		db = db.Order(o)
		if db.Error != nil {
			return nil, db.Error
		}
	}

	paginator := &Paginator{}

	tx := db.Session(&gorm.Session{})

	wg := &sync.WaitGroup{}
	wg.Add(1)
	// 取得总数
	go func() {
		defer wg.Done()
		res := tx.Count(&paginator.Total)
		if res.Error != nil {
			paginator.Total = -1
			return
		}
	}()

	// 设置PageSize，每页的对象数量
	paginator.PageSize = utils.IntMin(utils.IntMax(_pageSize, 1), global.MaxPageSize)
	// 设置Offset，查询偏移量
	paginator.Offset = (paginator.Page - 1) * paginator.PageSize

	// 分页查询
	res := tx.Limit(paginator.PageSize).Offset(paginator.Offset).Find(_result)
	if res.Error != nil {
		return nil, tx.Error
	}
	// 设置Data，查询结果
	paginator.Data = _result

	wg.Wait()
	// 设置Pages，总共有多少页
	paginator.Pages = int(math.Ceil(float64(paginator.Total) / float64(paginator.PageSize)))
	// 设置Page，当前查询第几页
	paginator.Page = utils.IntMin(utils.IntMax(_page, 1), paginator.Pages)

	return paginator, nil
}
