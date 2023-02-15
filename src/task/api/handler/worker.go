package handler

import (
	"eago/common/api/ext"
	"eago/common/tracer"
	"github.com/gin-gonic/gin"
)

// ListWorkers 列出所有Worker
func (th *TaskHandler) ListWorkers(c *gin.Context) {
	ext.WriteSuccessPayload(c, "workers", th.workerCli.List(tracer.ExtractTraceCtxFromGin(c)))
}
