package dto

// 流程创建者（发起人）信息在FormData中存储的Key
const (
	InitiatorKeyUserId      = "__user_id"
	InitiatorKeyUsernameKey = "__username"
	InitiatorKeyPhone       = "__phone"
)

// Getter 取值器
const (
	GetterDirect = "direct" // 直接取值
	GetterField  = "field"  // 从FormData的指定字段取值
)

var GetterAllowed = map[string]struct{}{
	GetterDirect: activeEmptyStruct,
	GetterField:  activeEmptyStruct,
}
