package api

import (
	"github.com/PaleBlueYk/randomSSQNumber/db"
	model2 "github.com/PaleBlueYk/randomSSQNumber/model"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/paleblueyk/logger"
	"net/http"
)

// SubmitMySSQ 提交双色球号码
func SubmitMySSQ(c *gin.Context) {
	var requestParams struct {
		Uid  string          `json:"uid" binding:"required"` //用户uid
		Num  int             `json:"num"`                    // 第几期
		List []model2.GenNum `json:"list"`                   // 选择的数字列表
	}
	if err := c.ShouldBindJSON(&requestParams); err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "数据结构错误!",
		})
		c.Abort()
		return
	}
	listData, _ := jsoniter.MarshalToString(requestParams.List)
	if err := db.Mysql.Create(&model2.NumSaveData{
		Uid:  requestParams.Uid,
		Num:  requestParams.Num,
		List: listData,
	}).Error; err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "数据库保存失败!",
		})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "提交成功!",
	})
}
