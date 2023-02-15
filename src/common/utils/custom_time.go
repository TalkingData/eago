package utils

import (
	"database/sql/driver"
	"eago/common/global"
	"fmt"
	"time"
)

type CustomTime struct {
	time.Time
}

// NewCustomTimeByTimestamp 通过时间戳获取CustomTime
func NewCustomTimeByTimestamp(ts int64) *CustomTime {
	return &CustomTime{Time: time.Unix(ts, 0)}
}

func (t *CustomTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", t.Format(global.TimestampFormat))
	return []byte(formatted), nil
}

func (t *CustomTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t == nil {
		return nil, nil
	} else if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

func (t *CustomTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = CustomTime{Time: value}
		return nil
	}
	return fmt.Errorf("Can not convert %v to timestamp", v)
}

func (t *CustomTime) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+global.TimestampFormat+`"`, string(data), time.Local)
	*t = CustomTime{Time: now}
	return
}

func (t *CustomTime) String() string {
	return t.Format(global.TimestampFormat)
}
