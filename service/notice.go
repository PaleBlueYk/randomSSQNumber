package service

import (
	"errors"
	"fmt"
	"github.com/PaleBlueYk/randomSSQNumber/config"
	"github.com/PaleBlueYk/randomSSQNumber/db"
	"github.com/PaleBlueYk/randomSSQNumber/model"
	jsoniter "github.com/json-iterator/go"
	"github.com/paleblueyk/logger"
	"github.com/robfig/cron/v3"
	"github.com/wxpusher/wxpusher-sdk-go"
	model2 "github.com/wxpusher/wxpusher-sdk-go/model"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

// Notice2User 通知用户
func Notice2User() {
	c := cron.New()
	enID, err := c.AddFunc("30 21 ? * 0,2,4", func() {
		logger.Info("开始广播中奖通知")
	RESTART:
		prizeInfo, _ := GetNewPrize()
		var checkNum model.CheckNum
		if err := db.Mysql.Where("need_check = 1 AND checked = 0").Last(&checkNum).Error; err != nil {
			logger.Warn("当前无需通知", err)
			return
		}
		logger.Info("期号: ", checkNum)
		result, err := BingoCheck(prizeInfo, NoticeUserList(strconv.Itoa(checkNum.Num)))
		if err != nil && err.Error() == "期号不对应" {
			// 期号不对应等5分钟后检查
			time.Sleep(5 * time.Minute)
			logger.Info("等待5分钟...")
			goto RESTART
		}
		if err != nil {
			logger.Error(err)
			return
		}
		for _, data := range result {
			html := NoticePage(data)
			_, err := wxpusher.SendMessage(model2.NewMessage(config.AppConf.WxPusher.AppToken).SetSummary(fmt.Sprintf("第%s期双色球中奖通知", prizeInfo.Num)).SetContentType(2).SetContent(html).AddUId(data.Uid))
			if err != nil {
				logger.Error(err)
			}
		}
	})
	if err != nil {
		logger.Error(err)
		return
	}
	c.Start()
	logger.Info("定时任务开启 EntryID: ", enID)
}

// NoticeUserList 获取需要通知的用户列表
func NoticeUserList(num string) (saveData []model.NumSaveData) {
	n, _ := strconv.Atoi(num)
	if err := db.Mysql.Where("num = ?", n).Find(&saveData).Error; err != nil {
		logger.Error(err)
		return
	}
	return
}

// BingoCheck 中奖计算
func BingoCheck(prizeInfo model.Prize, dataList []model.NumSaveData) (result []model.NumSaveData, err error) {
	for idx, data := range dataList {
		if data.List == "" {
			continue
		}
		if strconv.Itoa(data.Num) != prizeInfo.Num {
			logger.Error("期号不对应: %s - %s", data.Num, prizeInfo.Num)
			err = errors.New("期号不对应")
			return
		}
		var (
			list []model.GenNum
		)
		_ = jsoniter.UnmarshalFromString(data.List, &list)
		for _, chooseNum := range list {
			// 篮球检查
			var blueBingo bool
			if chooseNum.BlueNum == prizeInfo.BlueNum {
				blueBingo = true
			}
			// 红球检查
			redCount := redBingoCheck(chooseNum.RedNum, prizeInfo.RedNum)
			switch redCount {
			case 0, 1, 2:
				if blueBingo {
					dataList[idx].BingoInfo = "恭喜你,六等奖"
					dataList[idx].BingoMoney += 5
				}
			case 3:
				if blueBingo {
					dataList[idx].BingoInfo = "恭喜你,五等奖"
					dataList[idx].BingoMoney += 10
				}
			case 4:
				if blueBingo {
					dataList[idx].BingoInfo = "恭喜你,四等奖"
					dataList[idx].BingoMoney += 200
				} else {
					dataList[idx].BingoInfo = "恭喜你,五等奖"
					dataList[idx].BingoMoney += 10
				}
			case 5:
				if blueBingo {
					dataList[idx].BingoInfo = "恭喜你,三等奖"
					dataList[idx].BingoMoney += 3000
				} else {
					dataList[idx].BingoInfo = "恭喜你,四等奖"
					dataList[idx].BingoMoney += 200
				}
			case 6:
				if blueBingo {
					dataList[idx].BingoInfo = "恭喜你,一等奖(一等奖奖池浮动，以官方发布为准)"
					dataList[idx].BingoMoney += 5000000
				} else {
					dataList[idx].BingoInfo = "恭喜你,二等奖(二等奖奖池浮动，以官方发布为准)"
					dataList[idx].BingoMoney += 5000000 / 0.25
				}
			default:
				dataList[idx].BingoInfo = "很遗憾,未中奖"
				dataList[idx].BingoMoney = 0
			}
		}
	}
	result = dataList
	return
}

// 红球中奖检查
func redBingoCheck(chooseNum []string, prizeNum []string) (count uint) {
	for _, red := range chooseNum {
		for _, bingoRed := range prizeNum {
			if bingoRed == red {
				count++
			}
		}
	}
	return
}

// NoticePage 制作中奖通知页面
func NoticePage(data model.NumSaveData) string {
	var result string
	html, err := ioutil.ReadFile("./source/web/notice.html")
	if err != nil {
		logger.Error(err)
		return ""
	}
	result = strings.ReplaceAll(string(html), "{{BingoInfo}}", data.BingoInfo)
	result = strings.ReplaceAll(result, "{{BingoMoney}}", strconv.Itoa(int(data.BingoMoney)))
	return result
}
