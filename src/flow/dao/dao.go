package dao

import (
	"context"
	"eago/common/logger"
	"eago/flow/conf"
	"gorm.io/gorm"
)

type Dao struct {
	db *gorm.DB

	conf *conf.Conf

	lg *logger.Logger
}

func NewDao(d *gorm.DB, _conf *conf.Conf, lg *logger.Logger) *Dao {
	return &Dao{
		db: d,

		conf: _conf,

		lg: lg,
	}
}

func (d *Dao) Close() {
	if d == nil {
		return
	}

	db, _ := d.db.DB()
	if db != nil {
		_ = db.Close()
	}
}

func (d *Dao) getDb() *gorm.DB {
	return d.db
}

func (d *Dao) getDbWithCtx(ctx context.Context) *gorm.DB {
	return d.db.WithContext(ctx)
}
