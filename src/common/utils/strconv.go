package utils

import "strconv"

// Force Type convert
func Str2Int(s string) (int, error) {
	return strconv.Atoi(s)
}

func Str2Int32(s string) (int32, error) {
	v, err := Str2Int(s)
	return int32(v), err
}

func Str2Uint(s string) (uint, error) {
	v, err := strconv.ParseUint(s, 10, 64)
	return uint(v), err
}

func Str2Uint32(s string) (uint32, error) {
	v, err := Str2Int(s)
	return uint32(v), err
}

func Str2Int64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func Str2Uint64(s string) (uint64, error) {
	v, err := Str2Int64(s)
	return uint64(v), err
}

func Str2Float32(s string) (float32, error) {
	v, err := strconv.ParseFloat(s, 32)
	return float32(v), err
}

func Str2Float64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}
