package dto

import (
	"github.com/beego/beego/v2/core/validation"
)

// Login struct
type Login struct {
	Username string `json:"username" valid:"Required;MinSize(3);MaxSize(100)"`
	Password string `json:"password" valid:"Required;MinSize(6);MaxSize(100)"`
}

// Validate
func (l *Login) Validate() (err interface{}) {
	valid := validation.Validation{}
	// 验证数据
	ok, err := valid.Valid(l)
	if err != nil {
		return
	}
	// 数据验证未通过
	if !ok {
		return valid.Errors
	}

	return
}
