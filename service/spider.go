package service

import (
	"bytes"
	"fmt"
	"github.com/PaleBlueYk/randomSSQNumber/model"
	"github.com/PaleBlueYk/randomSSQNumber/pkg/site"
	"github.com/PaleBlueYk/randomSSQNumber/pkg/utils"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"github.com/paleblueyk/logger"
	"strconv"
	"time"
)

// GetNextNum 获取下一期期号
func GetNextNum() int {
	return GetPrizeInformation()[0].Num + 1
}

// GetPrizeInformation 获取开奖信息
func GetPrizeInformation() []model.PrizeInformation {
	doc := catCP500()
	var result []model.PrizeInformation
	doc.Find("tbody#tdata tr").EachWithBreak(func(i int, selection *goquery.Selection) bool {
		var item model.PrizeInformation
		selection.Find("td").EachWithBreak(func(i int, td *goquery.Selection) bool {
			switch i {
			case 0:
				item.Num, _ = strconv.Atoi(td.Text())
			case 1,2,3,4,5,6:
				item.RedNum = append(item.RedNum, td.Text())
			case 7:
				item.BlueNum = td.Text()
			case 9:
				item.MoneyPool = fmt.Sprintf(td.Text() + "元")
			case 10:
				item.Prize1Item, _ = strconv.Atoi(td.Text())
			case 11:
				item.Prize1Money = td.Text() + "元"
			case 12:
				item.Prize2Item, _ = strconv.Atoi(td.Text())
			case 13:
				item.Prize2Money = td.Text() + "元"
			case 14:
				item.TotalSales = td.Text() + "元"
			case 15:
				item.PrizeTime, _ = time.Parse("2006-01-02", td.Text())
			}
			return true
		})
		result = append(result, item)
		return true
	})

	return result
}

// 抓取彩票50
func catCP500() (doc *goquery.Document) {
	resp, err := resty.New().R().Get(site.CP500)
	if err != nil {
		logger.Error(err)
		return
	}
	//docStr, err := iconv.ConvertString(resp.String(), "gb2312","utf-8")
	docBs, err := utils.GbkToUtf8(resp.Body())
	if err != nil {
		logger.Error(err)
	}
	doc, err = goquery.NewDocumentFromReader(bytes.NewReader(docBs))
	if err != nil {
		logger.Error(err)
		return
	}
	return
}