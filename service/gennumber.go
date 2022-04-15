package service

import (
	"fmt"
	"github.com/paleblueyk/logger"
	"github.com/samber/lo"
	"math/rand"
	"sort"
	"time"
)

var seed int64
func init() {
	seed = time.Now().Unix()
	rand.Seed(seed)
	logger.Info("随机数种子: ", seed)
}

// GenNumber 生成随机号码
func GenNumber() string {
	return fmt.Sprintf("%s -- %d", redNumber(), blueNumber())
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
		num := rand.Intn(32) + 1
		// 防止数字重复
		if lo.Contains(redNumbers, num) {
			numRange++
			continue
		}
		redNumbers = append(redNumbers, num)
	}
	sort.Ints(redNumbers)
	var result string
	for idx, number := range redNumbers {
		if idx > 0 {
			result += ","
		}
		result += fmt.Sprintf("%d", number)
	}
	return result
}

// BlueNumber 篮球生成
func blueNumber() int {
	return rand.Intn(15) +1
}