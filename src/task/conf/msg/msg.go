package msg

import (
	m "eago/common/api-suite/message"
	"net/http"
)

var (
	Success = m.Message{0, "Success"}

	WarnInvalidUri     = m.Message{http.StatusBadRequest, "Bad request, invalid Uri."}
	WarnInvalidParams  = m.Message{http.StatusBadRequest, "Bad request, invalid params."}
	WarnInvalidBody    = m.Message{http.StatusBadRequest, "Bad request, invalid request body."}
	WarnPermissionDeny = m.Message{http.StatusForbidden, "Permission deny."}
	WarnNotFound       = m.Message{http.StatusNotFound, "Data not found."}

	ErrGetToken = m.Message{http.StatusInternalServerError, "An error occurred in get token content, please contact admin."}
	ErrDatabase = m.Message{http.StatusInternalServerError, "An database error occurred, please contact admin."}
	ErrCallTask = m.Message{http.StatusInternalServerError, "An error occurred in call task."}
	ErrUnknown  = m.Message{http.StatusInternalServerError, "Unknown error, please contact admin."}
)
