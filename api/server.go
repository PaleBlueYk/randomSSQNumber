package api

import (
	"context"
	"fmt"
	"github.com/PaleBlueYk/randomSSQNumber/config"
	"github.com/PaleBlueYk/randomSSQNumber/db"
	"github.com/PaleBlueYk/randomSSQNumber/model"
	"github.com/PaleBlueYk/randomSSQNumber/service"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/paleblueyk/logger"
	"github.com/wxpusher/wxpusher-sdk-go"
	model2 "github.com/wxpusher/wxpusher-sdk-go/model"
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

	if err = db.Mysql.FirstOrCreate(&model.CheckNum{Num: tmpData.Num, NeedCheck: true}).Error; err != nil {
		logger.Error("保存", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "提交成功!",
	})
}

// BingoCheck 中奖检查
func BingoCheck(c *gin.Context) {
	rid := c.Query("id")
	uid := c.Query("uid")
	var prizeInfo model.PrizeInformation
	prizeInfoList := service.GetPrizeInformation()
	var saveData model.NumSaveData
	if err := db.Mysql.Where("rdb_id = ?", rid).First(&saveData).Error; err != nil {
		logger.Error(err)
	}
	for _, information := range prizeInfoList {
		if information.Num == saveData.Num {
			prizeInfo = information
		}
	}
	result, err := service.BingoCheck(model.Prize{
		Num:     fmt.Sprintf("%d", prizeInfo.Num),
		RedNum:  prizeInfo.RedNum,
		BlueNum: prizeInfo.BlueNum,
	}, []model.NumSaveData{saveData})
	if err != nil && err.Error() == "期号不对应" {
		_, err := wxpusher.SendMessage(model2.NewMessage(config.AppConf.WxPusher.AppToken).SetSummary(fmt.Sprintf("第%d期双色球中奖通知", prizeInfo.Num)).SetContentType(1).SetContent(fmt.Sprintf("第%d期结果未出，请耐心等待", prizeInfo.Num)).AddUId(uid))
		if err != nil {
			logger.Error(err)
		}
	}
	if err != nil {
		logger.Error(err)
		return
	}
	for _, data := range result {
		html := service.NoticePage(data)
		_, err := wxpusher.SendMessage(model2.NewMessage(config.AppConf.WxPusher.AppToken).SetSummary(fmt.Sprintf("第%d期双色球中奖通知", prizeInfo.Num)).SetContentType(2).SetContent(html).AddUId(data.Uid))
		if err != nil {
			logger.Error(err)
		}
	}
}
