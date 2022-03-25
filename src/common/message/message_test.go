package message

import (
	"eago/common/api-suite/writter"
	"eago/common/log"
	"fmt"
	"github.com/beego/beego/v2/core/validation"
	"github.com/gin-gonic/gin"
	"strings"
	"testing"
)

const (
	CODE_COMMON_UNKNOWN_ERROR = iota + 990000
	CODE_COMMON_BADREQUEST_URI
	CODE_COMMON_BADREQUEST_BODY
)

var (
	UnknownErr  = NewMessage(CODE_COMMON_UNKNOWN_ERROR, "Unknown error.")
	InvalidBody = NewMessage(CODE_COMMON_BADREQUEST_BODY, "Bad request, invalid request body.")
)

func TestMessage(t *testing.T) {
	eg := gin.Default()
	eg.POST("/users", NewUserHandler)

	fmt.Println("========================================")
	fmt.Println("For valid failed test:")
	fmt.Println("curl -X POST 'http://localhost:8888/users' --header 'Content-Type: application/json' --data-raw '{\"name\":\"admin\"}'")
	fmt.Println("For valid success test:")
	fmt.Println("curl -X POST 'http://localhost:8888/users' --header 'Content-Type: application/json' --data-raw '{\"name\":\"BeeUser\",\"age\":25,\"email\":\"test@test.com\",\"description\":\"\"}'")
	fmt.Println("========================================")

	fmt.Println("Starting server...")

	_ = eg.Run(":8888")
}

func NewUserHandler(c *gin.Context) {
	var userFrm User

	if err := c.BindJSON(&userFrm); err != nil {
		m := InvalidBody.SetError(err)
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	valid := validation.Validation{}
	ok, err := valid.Valid(userFrm)
	if err != nil {
		m := InvalidBody.SetDetail("An error occurred when valid.Valid.").SetError(err)
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}
	if !ok {
		m := InvalidBody.SetError(valid.Errors)
		log.WarnWithFields(m.LogFields())
		m.WriteRest(c)
		return
	}

	// Do new user in model
	// ...

	writter.WriteSuccessPayload(c, "user", userFrm)
}

// 验证函数写在 "valid" tag 的标签里
// 各个函数之间用分号 ";" 分隔，分号后面可以有空格
// 参数用括号 "()" 括起来，多个参数之间用逗号 "," 分开，逗号后面可以有空格
// 正则函数(Match)的匹配模式用两斜杠 "/" 括起来
// 各个函数的结果的 key 值为字段名.验证函数名
type User struct {
	Id          int
	Name        string  `json:"name" valid:"Required;Match(/^Bee.*/)"`        // Name 不能为空并且以 Bee 开头
	Age         int     `json:"age" valid:"Range(1, 140)"`                    // 1 <= Age <= 140，超出此范围即为不合法
	Email       string  `json:"email" valid:"Email; MaxSize(100)"`            // Email 字段需要符合邮箱格式，并且最大长度不能大于 100 个字符
	Description *string `json:"description" valid:"MinSize(0); MaxSize(500)"` // Description 字段必须存在，并且最大长度不能大于 500 个字符
}

// 如果你的 struct 实现了接口 validation.ValidFormer
// 当 StructTag 中的测试都成功时，将会执行 Valid 函数进行自定义验证
func (u *User) Valid(v *validation.Validation) {
	if strings.Index(u.Name, "admin") != -1 {
		// 通过 SetError 设置 Name 的错误信息，HasErrors 将会返回 true
		_ = v.SetError("Name", "名称里不能含有 admin")
	}
}

func init() {
	// 加载日志设置
	err := log.InitLog(
		"./logs",
		"eago-test",
		"debug",
	)
	if err != nil {
		fmt.Println("Failed to init logging, error:", err.Error())
		panic(err)
	}
}
