package worker

import (
	"encoding/json"
)

// WriteResponse
func WriteResponse(code int64, msg string) []byte {
	rspBody := WorkerResponse{
		Code:    code,
		Message: msg,
	}
	b, _ := json.Marshal(rspBody)
	return b
}
