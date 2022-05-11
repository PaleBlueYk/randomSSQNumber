package v1

import (
	"github.com/PaleBlueYk/randomSSQNumber/api"
	"github.com/gin-gonic/gin"
)

func SubmitNum(router *gin.RouterGroup)  {
	r := router.Group("/num")
	{
		r.GET("", api.SubmitMySSQ)
		r.GET("/bingoCheck", api.BingoCheck)
	}
}
