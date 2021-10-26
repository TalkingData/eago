package utils

type Gorm struct {
	Id        int        `json:"id" gorm:"type:int(11) auto_increment;comment:主键id"`
	CreatedAt *LocalTime `json:"created_at" gorm:"type:datetime not null;comment:创建时间"`
	UpdatedAt *LocalTime `json:"updated_at" gorm:"comment:更新时间"`
}
