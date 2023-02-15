package utils

import (
	"errors"
	"reflect"
)

// ReversSlice 反转切片
func ReversSlice(s *[]string) {
	for i, j := 0, len(*s)-1; i < j; i, j = i+1, j-1 {
		(*s)[i], (*s)[j] = (*s)[j], (*s)[i]
	}
}

// IsInSlice 判断needle是否是haystack中的一项
func IsInSlice(haystack, needle interface{}) (bool, error) {
	sVal := reflect.ValueOf(haystack)
	kind := sVal.Kind()
	if kind == reflect.Slice || kind == reflect.Array {
		for i := 0; i < sVal.Len(); i++ {
			if sVal.Index(i).Interface() == needle {
				return true, nil
			}
		}

		return false, nil
	}

	return false, errors.New("not supported haystack")
}

// MergeStringSlice 合并除空出串外的字符串切片
func MergeStringSlice(ss ...[]string) []string {
	var newS []string
	for _, s := range ss {
		for _, ele := range s {
			if ele != "" {
				newS = append(newS, ele)
			}
		}
	}

	return newS
}

// RemoveStringSliceElement 删除字符串切片中指定元素
func RemoveStringSliceElement(s []string, ele string) []string {
	var final []string
	for _, e := range s {
		if e == ele {
			continue
		}
		final = append(final, e)
	}

	return final
}
