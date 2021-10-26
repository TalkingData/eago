package utils

import (
	"bytes"
	"fmt"
	"strconv"
	"unicode"
)

func Camel2Case(name string) string {
	buf := newBuffer()
	for i, r := range name {
		if unicode.IsUpper(r) {
			if i != 0 {
				buf.Append('_')
			}
			buf.Append(unicode.ToLower(r))
		} else {
			buf.Append(r)
		}
	}
	return buf.String()
}

func newBuffer() *buffer {
	return &buffer{Buffer: new(bytes.Buffer)}
}

type buffer struct {
	*bytes.Buffer
}

func (b *buffer) Append(i interface{}) *buffer {
	switch val := i.(type) {
	case int:
		b.append(strconv.Itoa(val))
	case int64:
		b.append(strconv.FormatInt(val, 10))
	case uint:
		b.append(strconv.FormatUint(uint64(val), 10))
	case uint64:
		b.append(strconv.FormatUint(val, 10))
	case string:
		b.append(val)
	case []byte:
		_, _ = b.Write(val)
	case rune:
		_, _ = b.WriteRune(val)
	}
	return b
}

func (b *buffer) append(s string) *buffer {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("out of memory")
		}
	}()
	_, _ = b.WriteString(s)
	return b
}
