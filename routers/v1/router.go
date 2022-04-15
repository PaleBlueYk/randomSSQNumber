package v1

import "github.com/gin-gonic/gin"

func Routers(r *gin.Engine) {
	v1 := r.Group("/v1")
	WxPusherRouter(v1)
}
