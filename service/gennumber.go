package service

import (
	"encoding/json"
	"fmt"
	"github.com/PaleBlueYk/randomSSQNumber/config"
	"github.com/PaleBlueYk/randomSSQNumber/model"
	"github.com/paleblueyk/logger"
	"github.com/samber/lo"
	"io/ioutil"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"
)

var seed int64

// GetPage 获取页面
func GetPage(count int, uid string) string {
	var result string
	var numModeList []model.GenNum
	for i := 0; i < count; i++ {
		numModeList = append(numModeList, getGenNum())
	}
	list := GenNumberHtml(numModeList)
	nextNum := GetNextNum()
	htmlTemplate, err := ioutil.ReadFile("./source/web/template.html")
	if err != nil {
		logger.Error(err)
		return list
	}
	dataMap := make(map[string]interface{})
	dataMap = map[string]interface{}{
		"uid":  uid,
		"num":  nextNum,
		"list": numModeList,
	}
	submitData, err := json.Marshal(&dataMap)
	if err != nil {
		logger.Error(err)
		return list
	}

	result = strings.ReplaceAll(string(htmlTemplate), "{{htmlContent}}", list) // 显示数组
	result = strings.ReplaceAll(result, "{{Num}}", strconv.Itoa(nextNum))      // 显示数组
	result = strings.ReplaceAll(result, "{{BaseUrl}}", config.AppConf.BaseUrl) // 网站部署地址
	result = strings.ReplaceAll(result, "{{submitData}}", string(submitData))          // 提交数据

	return result
}

// GenNumberHtml 生成随机号码
func GenNumberHtml(numModel []model.GenNum) string {
	seed = time.Now().Unix()
	rand.Seed(seed)
	logger.Info("随机数种子: ", seed)
	var result string
	for _, num := range numModel {
		var redStr string
		for _, red := range num.RedNum {
			redStr += fmt.Sprintf("<span class='redBall'>%s</span>", red)
		}
		result += fmt.Sprintf("<li>%s=><span class='blueBall'>%s</span></li>", redStr, num.BlueNum)
	}

	return result
}

/**
生成规则:
http://www.cwl.gov.cn/c/2018/10/12/417937.shtml
双色球投注区分为红色球号码区和蓝色球号码区，红色球号码区由1-33共三十三个号码组成，蓝色球号码区由1-16共十六个号码组成。
*/

func getGenNum() (result model.GenNum) {
	result.RedNum = redNumber()
	result.BlueNum = blueNumber()
	return
}

// RedNumber 红球生成
func redNumber() []string {
	var redNumbers []int
	numRange := 6
	for i := 0; i < numRange; i++ {
		num := rand.Intn(33) + 1
		// 防止数字重复
		if lo.Contains(redNumbers, num) {
			numRange++
			continue
		}
		redNumbers = append(redNumbers, num)
	}
	sort.Ints(redNumbers)
	var result []string
	for _, number := range redNumbers {
		num := strconv.Itoa(number)
		if len(num) == 1 {
			num = "0" + num
		}
		//result += fmt.Sprintf("<span class='redBall'>%s</span>", num)
		result = append(result, num)
	}
	return result
}

// BlueNumber 篮球生成
func blueNumber() string {
	num := rand.Intn(16) + 1
	numStr := strconv.Itoa(num)
	if len(numStr) == 1 {
		numStr = "0" + numStr
	}
	//return fmt.Sprintf("<span class='blueBall'>%s</span>", numStr)
	return numStr
}
