package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// PriceDivBy100 将字符串类型的整数的小数点往前移动两位
func PriceDivBy100(price string) (*string, error) {
	_, err := strconv.Atoi(price)
	if err != nil {
		return nil, errors.New("format error")
	}

	ret := ""
	length := len(price)
	if length > 2 {
		ret = price[:length-2] + "." + price[length-2:]
	} else {
		tmp := fmt.Sprintf("%03s%s", price, "")[:3]
		ret = tmp[:1] + "." + tmp[1:]
	}

	return &ret, nil
}

// PriceMulBy100 将浮点类型的字符串的小数点往后移动两位
func PriceMulBy100(price string) (*string, error) {
	if price == "" {
		return nil, errors.New("format error")
	}
	// 获取.出现的次数
	count := 0
	for _, item := range price {
		if item == 46 {
			count++
			continue
		}
		if item > 57 || item < 48 {
			return nil, errors.New("format error")
		}
	}

	ret := ""

	if count == 0 {
		ret = price + "00"
	} else if count == 1 {
		tmp := strings.Split(price, ".")
		right := fmt.Sprintf("%s%04s", tmp[1], "")[:4]
		ret = tmp[0] + right[0:2]
	} else {
		return nil, errors.New("format error")
	}

	return &ret, nil
}
