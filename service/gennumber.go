package service

import (
	"fmt"
	"github.com/paleblueyk/logger"
	"github.com/samber/lo"
	"math/rand"
	"sort"
	"strconv"
	"time"
)

var seed int64

// GenNumber 生成随机号码
func GenNumber(count int) string {
	seed = time.Now().Unix()
	rand.Seed(seed)
	logger.Info("随机数种子: ", seed)
	var result string
	for i := 0; i < count; i++ {
		result += fmt.Sprintf("<li>%s=><span style='color:white;background-color: blue; width: 50px;line-height: 50px; border-radius: 50%%; padding: 5px;'>%s</span></li>", redNumber(), blueNumber())
	}

	return result
}

/**
生成规则:
http://www.cwl.gov.cn/c/2018/10/12/417937.shtml
双色球投注区分为红色球号码区和蓝色球号码区，红色球号码区由1-33共三十三个号码组成，蓝色球号码区由1-16共十六个号码组成。
*/

// RedNumber 红球生成
func redNumber() string {
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
	var result string
	for _, number := range redNumbers {
		num := strconv.Itoa(number)
		if len(num) == 1 {
			num = "0" + num
		}
		result += fmt.Sprintf("<span style='color:white;background-color: red; width: 50px;line-height: 50px; border-radius: 50%%; padding: 5px;'>%s</span>", num)
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
	return numStr
}
