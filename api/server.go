package api

import (
	"context"
	"github.com/PaleBlueYk/randomSSQNumber/db"
	"github.com/PaleBlueYk/randomSSQNumber/model"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/paleblueyk/logger"
	"gorm.io/gorm"
	"net/http"
)

// SubmitMySSQ 提交双色球号码
func SubmitMySSQ(c *gin.Context) {
	id := c.Query("id")
	logger.Info("保存数据 redis id: ", id)
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "id不能为空",
		})
		c.Abort()
		return
	}
	rdResult, err := db.RDB.Get(context.Background(), id).Result()
	if err != nil {
		logger.Error(err)
		return
	}
	var tmpData model.TmpSaveNum
	if err = jsoniter.UnmarshalFromString(rdResult, &tmpData); err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "数据不存在",
		})
		c.Abort()
		return
	}
	var numData model.NumSaveData
	if err := db.Mysql.Where("rdb_id = ?", tmpData.ID).First(&numData).Error; err != nil && err != gorm.ErrRecordNotFound {
		logger.Error(err)
	}
	if numData.ID > 0 {
		logger.Error("数据已存在,id: %d", numData.ID)
		c.JSON(http.StatusAlreadyReported, gin.H{
			"msg": "数据已保存",
		})
		c.Abort()
		return
	}
	lisStr, _ := jsoniter.MarshalToString(tmpData.List)
	if err = db.Mysql.Save(&model.NumSaveData{
		RdbID: tmpData.ID,
		Uid:   tmpData.UID,
		Num:   tmpData.Num,
		List:  lisStr,
	}).Error; err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "保存数据失败",
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "提交成功!",
	})
}
