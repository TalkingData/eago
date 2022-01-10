package conf

const (
	SERVICE_NAME     = "eago-flow"
	RPC_REGISTER_KEY = "eago.srv.flow"
	API_REGISTER_KEY = "eago.api.flow"

	AUTH_RPC_REGISTER_KEY = "eago.srv.auth"
	TASK_RPC_REGISTER_KEY = "eago.srv.task"

	ADMIN_ROLE_NAME = "flow_admin"

	TIMESTAMP_FORMAT = "2006-01-02 15:04:05"

	CONFIG_FILE_PATHNAME = "../conf/eago_flow.conf"
)

// 流程实例名称最大长度
const INSTANCE_NAME_MAX_LENGTH = 200

// 流程实例状态
const (
	INSTANCE_STATUS_PANIC_END    = -200 // 系统异常
	INSTANCE_STATUS_REJECTED_END = -1   // 被驳回
	INSTANCE_STATUS_APPROVED_END = 0    // 审批通过
	INSTANCE_STATUS_PENDING      = 1    // 系统处理中
	INSTANCE_STATUS_RUNNING      = 2    // 流转中
)

// 审批人分隔符
const ASSIGNEES_SPILT_TAG = ","

// AssigneeCondition 取值枚举范围
const (
	AC_INITIATOR                           = "initiator"                          // 发起人
	AC_INITIATORS_DEPARTMENTS_OWNER        = "initiators_department_owner"        // 发起人所属部门Owner
	AC_INITIATORS_PARENT_DEPARTMENTS_OWNER = "initiators_parent_department_owner" // 发起人所属上级部门Owner
	AC_SPECIFIED_USERS                     = "specified_users"                    // 指定用户（多个）
	AC_SPECIFIED_PRODUCT_OWNER             = "specified_product_owner"            // 指定产品线Owner
	AC_SPECIFIED_GROUP_OWNER               = "specified_group_owner"              // 指定组Owner
	AC_SPECIFIED_DEPARTMENT_OWNER          = "specified_department_owner"         // 指定部门Owner
	AC_SPECIFIED_ROLE                      = "specified_role"                     // 指定角色
)

// 流程创建者（发起人）信息在FormData中存储的Key
const (
	INITIATOR_USER_ID_KEY  = "_user_id"
	INITIATOR_USERNAME_KEY = "_username"
	INITIATOR_PHONE_KEY    = "_phone"
)

// Getter 取值器
const (
	GETTER_DIRECT = "direct" // 直接取值
	GETTER_FIELD  = "field"  // 从FormData的指定字段取值
)

// NodeCategory 取值枚举范围
const (
	NODE_CATEGORY_FIRST  = -1 // 首节点
	NODE_CATEGORY_ANY    = 1  // 或签
	NODE_CATEGORY_ALL    = 2  // 会签
	NODE_CATEGORY_INFORM = 3  // 知会
)
