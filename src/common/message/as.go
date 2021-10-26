package message

import (
	"strconv"
	"strings"
)

// As 尝试将字符串转换为Message
func As(s string, m *Message) bool {
	split := strings.Split(s, RPC_ERROR_SPLITOR)
	if len(split) != 2 {
		return false
	}
	if len(split[1]) < 1 {
		return false
	}

	code, err := strconv.Atoi(split[0])
	if err != nil {
		return false
	}
	if code < 100 {
		return false
	}

	m.code = code
	m.msg = split[1]
	return true
}
