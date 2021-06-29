package msg

import (
	m "eago/common/api-suite/message"
	"net/http"
)

var (
	Success = m.Message{0, "Success"}

	WarnLoginFailed    = m.Message{http.StatusUnauthorized, "Login failed"}
	WarnPermissionDeny = m.Message{http.StatusForbidden, "Permission deny"}
	WarnInvalidUri     = m.Message{http.StatusBadRequest, "Bad request, invalid Uri"}
	WarnInvalidParams  = m.Message{http.StatusBadRequest, "Bad request, invalid params"}
	WarnInvalidBody    = m.Message{http.StatusBadRequest, "Bad request, invalid request body"}

	ErrGenToken = m.Message{http.StatusInternalServerError, "An error occurred in generate token, please contact admin"}
	ErrGetToken = m.Message{http.StatusInternalServerError, "An error occurred in get token content, please contact admin"}
	ErrDatabase = m.Message{http.StatusInternalServerError, "An database error occurred, please contact admin"}
	ErrUnknown  = m.Message{http.StatusInternalServerError, "Unknown error, please contact admin"}
)
