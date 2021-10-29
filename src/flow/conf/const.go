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

// 审批人分隔符
const ASSIGNEES_SPILT_TAG = ","

// 流程实例状态
const (
	INSTANCE_PANIC_END_STATUS    = 200 // 系统异常
	INSTANCE_REJECTED_END_STATUS = -1  // 被驳回
	INSTANCE_APPROVED_END_STATUS = 0   // 审批通过
	INSTANCE_PENDING_STATUS      = 1   // 系统处理中
	INSTANCE_RUNNING_STATUS      = 2   // 流转中
)

// AssigneeConditions 取值枚举范围
const (
	INITIATOR                           = "initiator"                          // 发起人
	INITIATORS_DEPARTMENTS_OWNER        = "initiators_department_owner"        // 发起人所属部门Owner
	INITIATORS_PARENT_DEPARTMENTS_OWNER = "initiators_parent_department_owner" // 发起人所属上级部门Owner
	SPECIFIED_USERS                     = "specified_users"                    // 指定用户（多个）
	SPECIFIED_PRODUCT_OWNER             = "specified_product_owner"            // 指定产品线Owner
	SPECIFIED_GROUP_OWNER               = "specified_group_owner"              // 指定组Owner
	SPECIFIED_DEPARTMENT_OWNER          = "specified_department_owner"         // 指定部门Owner
	SPECIFIED_ROLE                      = "specified_role"                     // 指定角色
)

// Category 取值枚举范围
const (
	ANY    = iota // 或签
	ALL           // 会签
	INFORM        // 知会
)
