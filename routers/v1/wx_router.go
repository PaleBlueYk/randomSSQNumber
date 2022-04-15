package v1

import (
	"github.com/PaleBlueYk/randomSSQNumber/api"
	"github.com/gin-gonic/gin"
)

func WxPusherRouter(router *gin.RouterGroup) {
	r := router.Group("/wxpusher")
	{
		r.POST("/upload_info", api.UserUploadInfo)
	}
}
