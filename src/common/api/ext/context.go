package ext

import (
	"eago/common/utils"
	"github.com/gin-gonic/gin"
)

func ParamUint64(c *gin.Context, key string) (uint64, error) {
	return utils.Str2Uint64(c.Param(key))
}

func ParamUint32(c *gin.Context, key string) (uint32, error) {
	return utils.Str2Uint32(c.Param(key))
}

func ParamUint(c *gin.Context, key string) (uint, error) {
	return utils.Str2Uint(c.Param(key))
}

func ParamInt64(c *gin.Context, key string) (int64, error) {
	return utils.Str2Int64(c.Param(key))
}

func ParamInt32(c *gin.Context, key string) (int32, error) {
	return utils.Str2Int32(c.Param(key))
}

func ParamInt(c *gin.Context, key string) (int, error) {
	return utils.Str2Int(c.Param(key))
}

func ParamFloat64(c *gin.Context, key string) (float64, error) {
	return utils.Str2Float64(c.Param(key))
}

func ParamFloat(c *gin.Context, key string) (float32, error) {
	return utils.Str2Float32(c.Param(key))
}
