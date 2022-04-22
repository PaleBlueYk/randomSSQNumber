package api

import (
	"github.com/PaleBlueYk/randomSSQNumber/config"
	"github.com/PaleBlueYk/randomSSQNumber/service"
	"github.com/gin-gonic/gin"
	"github.com/paleblueyk/logger"
	"github.com/wxpusher/wxpusher-sdk-go"
	"github.com/wxpusher/wxpusher-sdk-go/model"
	"strconv"
)

// UserUploadInfo 处理用户的上行消息
func UserUploadInfo(c *gin.Context) {
	var requestParams struct {
		Action string `json:"action"` // //动作，send_up_cmd 表示上行消息回调，后期可能会添加其他动作，请做好兼容。
		Data   struct {
			Uid         string      `json:"uid"`         //用户uid
			AppId       int         `json:"appId"`       //应用id
			AppKey      interface{} `json:"appKey"`      //应用名称
			AppName     string      `json:"appName"`     // 废弃
			UserName    string      `json:"userName"`    // 新用户无
			UserHeadImg string      `json:"userHeadImg"` // 新用户无
			Time        int64       `json:"time"`        //发生时间
			Content     string      `json:"content"`     //用户发送的内容
		} `json:"data"`
	}

	if err := c.ShouldBindJSON(&requestParams); err != nil {
		logger.Error(err)
		return
	}
	if requestParams.Action != "send_up_cmd" {
		logger.Error("此消息不是上行消息")
		return
	}
	logger.Info("用户uid: %s, 应用id： %d, 应用名称: %v, content: %s", requestParams.Data.Uid, requestParams.Data.AppId, requestParams.Data.AppName, requestParams.Data.Content)
	count, err := strconv.Atoi(requestParams.Data.Content)
	if err != nil {
		logger.Error(err)
		msgArr, err := wxpusher.SendMessage(model.NewMessage(config.AppConf.WxPusher.AppToken).SetSummary("处理错误").SetContent(err.Error()).AddUId(requestParams.Data.Uid))
		if err != nil {
			logger.Error(err)
			return
		}
		logger.Info(msgArr)
		return
	}
	msg := service.GetPage(count, requestParams.Data.Uid)
	resultList, err := wxpusher.SendMessage(model.NewMessage(config.AppConf.WxPusher.AppToken).SetSummary("双色球随机数").SetContentType(2).SetContent(msg).AddUId(requestParams.Data.Uid))
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Info(resultList)
}
