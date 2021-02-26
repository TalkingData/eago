package msg

import (
	"eago-common/api-suite/message"
	"net/http"
)

var (
	Success = message.BaseMsg{0, "Success"}

	// Invalid
	WarnInvalidUri    = message.BaseMsg{http.StatusBadRequest, "Bad request, invalid Uri."}
	WarnInvalidParams = message.BaseMsg{http.StatusBadRequest, "Bad request, invalid request body."}
	WarnInvalidBody   = message.BaseMsg{http.StatusBadRequest, "Bad request, invalid params."}

	// Token
	ErrGenToken = message.BaseMsg{http.StatusInternalServerError, "An error occurred in generate token, please contact admin."}
	ErrGetToken = message.BaseMsg{http.StatusInternalServerError, "An error occurred in get token content, please contact admin."}

	// Login
	WarnLoginFailed = message.BaseMsg{http.StatusUnauthorized, "Login failed."}

	// Permission
	WarnPermissionDeny = message.BaseMsg{http.StatusForbidden, "Permission deny."}

	// Database
	ErrDatabase = message.BaseMsg{http.StatusInternalServerError, "An database error occurred, please contact admin."}

	// Others
	ErrUnknown = message.BaseMsg{http.StatusInternalServerError, "Unknown error, please contact admin."}
)
