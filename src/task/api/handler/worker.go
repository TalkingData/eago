package handler

import (
	"eago/task/cli"
	"eago/task/conf/msg"
	"github.com/gin-gonic/gin"
)

// ListWorkers 列出所有Worker
// @Summary 列出所有Worker
// @Tags Worker
// @Param token header string true "Token"
// @Success 200 {string} string "{"code":0,"message":"Success","workers":[{"modular":"task","address":"172.30.105.34:46565","worker_id":"task.worker-33b9a8bd-dd6e-4eb9-92ce-6b51a09b9abe","start_time":"2021-05-14 15:20:41"},{"modular":"task","address":"172.30.105.34:41684","worker_id":"task.worker-579864f5-d0f0-49b4-bad8-a1993ca1700c","start_time":"2021-05-14 15:12:04"}]}"
// @Router /workers [GET]
func ListWorkers(c *gin.Context) {
	resp := msg.Success.GenResponse().SetPayload("workers", cli.WorkerClient.List())
	resp.Write(c)
}
