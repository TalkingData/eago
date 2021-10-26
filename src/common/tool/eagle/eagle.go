package eagle

import (
	"eago/common/log"
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type eagle struct {
	db *gorm.DB
}

// 创建Eagle工具
func NewEagle(address, user, password, dbName string) *eagle {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user,
		password,
		address,
		dbName,
	)

	m := mysql.New(mysql.Config{
		DSN:                      dsn,
		DefaultStringSize:        200,
		DisableDatetimePrecision: true,
	})

	db, err := gorm.Open(m, &gorm.Config{})
	if err != nil {
		panic(err.Error())
		return nil
	}

	return &eagle{db: db}
}

// ListProducts 列出所有产品线
func (e *eagle) ListProducts() (*[]Product, error) {
	prods := make([]Product, 0)

	rawSql := "SELECT name, alias, description " +
		"FROM management_product " +
		"WHERE enabled = 1;"
	if res := e.db.Raw(rawSql).Scan(&prods); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"error": res.Error.Error(),
			}, "Record not found in ListProducts.")
			return nil, nil
		}
		log.ErrorWithFields(log.Fields{
			"error": res.Error.Error(),
		}, "Error in ListProducts.")
		return nil, res.Error
	}

	return &prods, nil
}

// ListProductOwners 列出所有产品线Owner
func (e *eagle) ListProductsOwners() (*[]ProductMember, error) {
	pms := make([]ProductMember, 0)

	rawSql := "SELECT username AS username, p.`name` AS product " +
		"FROM management_approval AS ap " +
		"LEFT JOIN auth_user AS u ON u.id = ap.user_id " +
		"LEFT JOIN management_product AS p ON p.id = ap.product_id " +
		"WHERE u.is_active = 1 AND p.enabled = 1 AND ap.role_codename = 'leader' AND username LIKE '%@tendcloud.com';"
	if res := e.db.Raw(rawSql).Scan(&pms); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"error": res.Error.Error(),
			}, "Record not found in ListProductsOwners.")
			return nil, nil
		}
		log.ErrorWithFields(log.Fields{
			"error": res.Error.Error(),
		}, "Error in ListProductsOwners.")
		return nil, res.Error
	}

	return &pms, nil
}

// ListProductsMembers 列出所有产品线成员
func (e *eagle) ListProductsMembers() (*[]ProductMember, error) {
	pms := make([]ProductMember, 0)

	rawSql := "SELECT u.username AS username, p.name AS product " +
		"FROM management_product AS p " +
		"LEFT JOIN auth_user_groups AS g ON g.group_id = p.members_group_id " +
		"LEFT JOIN auth_user AS u ON u.id = g.user_id " +
		"WHERE u.is_active = 1 AND p.enabled = 1 AND username LIKE '%@tendcloud.com';"

	if res := e.db.Raw(rawSql).Scan(&pms); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.WarnWithFields(log.Fields{
				"error": res.Error.Error(),
			}, "Record not found in ListProductsOwners.")
			return nil, nil
		}
		log.ErrorWithFields(log.Fields{
			"error": res.Error.Error(),
		}, "Error in ListProductsOwners.")
		return nil, res.Error
	}

	return &pms, nil
}

func (e *eagle) Close() {
	if e.db == nil {
		return
	}

	d, _ := e.db.DB()
	_ = d.Close()
}
