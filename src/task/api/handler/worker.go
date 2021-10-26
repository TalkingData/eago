package handler

import (
	w "eago/common/api-suite/writter"
	"eago/common/tracer"
	"eago/task/cli"
	"github.com/gin-gonic/gin"
)

// ListWorkers 列出所有Worker
func ListWorkers(c *gin.Context) {
	sp, ctx := tracer.StartSpanFromContext(tracer.ExtractTraceContext(c))
	defer sp.Finish()

	w.WriteSuccessPayload(c, "workers", cli.WorkerClient.List(ctx))
}
