package service

import (
	"encoding/json"
	"github.com/PaleBlueYk/randomSSQNumber/model"
	"github.com/PaleBlueYk/randomSSQNumber/pkg/site"
	"github.com/PaleBlueYk/randomSSQNumber/pkg/utils"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"github.com/paleblueyk/logger"
	"strconv"
	"strings"
	"time"
)

// GetNextNum 获取下一期期号
func GetNextNum() int {
	return GetPrizeInformation()[0].Num + 1
}

// GetPrizeInformation 获取开奖信息
func GetPrizeInformation() []model.PrizeInformation {
	result, err := getSSQHistory()
	if err != nil {
		logger.Error(err)
		return nil
	}
	var resultList []model.PrizeInformation
	for _, r := range result.Result {
		num, _ := strconv.Atoi(r.Code)
		var prize1Item, prize2Item int
		var prize1Money, prize2Money string

		for _, prizegrade := range r.Prizegrades {
			switch prizegrade.Type {
			case 1:
				prize1Item, _ = strconv.Atoi(prizegrade.Typenum)
				prize1Money = prizegrade.Typemoney
			case 2:
				prize2Item, _ = strconv.Atoi(prizegrade.Typenum)
				prize2Money = prizegrade.Typemoney
			}
		}

		t := strings.ReplaceAll(r.Date, "(二)", "")
		t = strings.ReplaceAll(r.Date, "(四)", "")
		t = strings.ReplaceAll(r.Date, "(日)", "")
		prizeTime, _ := time.Parse("2006-01-02", t)
		resultList = append(resultList, model.PrizeInformation{
			Num:     num,
			RedNum:  utils.StrList2code(r.Red),
			BlueNum: r.Blue,
			Prize1Item: prize1Item,
			Prize1Money: prize1Money,
			Prize2Item: prize2Item,
			Prize2Money: prize2Money,
			MoneyPool: r.Poolmoney,
			TotalSales: r.Sales,
			PrizeTime: prizeTime,
		})
	}

	return resultList
}

// GetNewPrize 获取本期开奖号码
func GetNewPrize() (model.Prize, error) {
	//doc := catGovSite()
	//num := doc.Find("div.ssqQh-dom").Text()
	//logger.Info("第%s期", num)
	//redNum := doc.Find(".ssqRed-dom").Text()
	//logger.Info("红球: ", redNum)
	//blueNum := doc.Find(".ssqBlue-dom").Text()
	//logger.Info("篮球: ", blueNum)
	//var (
	//	redNumList  []string
	//	blueNumList []string
	//)
	//redNumList = utils.StrList2code(redNum)
	//blueNumList = utils.StrList2code(blueNum)
	//sort.Strings(redNumList)
	resultList := GetPrizeInformation()

	return model.Prize{
		Num:     strconv.Itoa(resultList[0].Num),
		RedNum:  resultList[0].RedNum,
		BlueNum: resultList[0].BlueNum,
	}, nil
}

type SSQHistoryResult struct {
	State     int    `json:"state"`
	Message   string `json:"message"`
	PageCount int    `json:"pageCount"`
	CountNum  int    `json:"countNum"`
	Tflag     int    `json:"Tflag"`
	Result    []struct {
		Name        string `json:"name"`
		Code        string `json:"code"`
		DetailsLink string `json:"detailsLink"`
		VideoLink   string `json:"videoLink"`
		Date        string `json:"date"`
		Week        string `json:"week"`
		Red         string `json:"red"`
		Blue        string `json:"blue"`
		Blue2       string `json:"blue2"`
		Sales       string `json:"sales"`
		Poolmoney   string `json:"poolmoney"`
		Content     string `json:"content"`
		Addmoney    string `json:"addmoney"`
		Addmoney2   string `json:"addmoney2"`
		Msg         string `json:"msg"`
		Z2Add       string `json:"z2add"`
		M2Add       string `json:"m2add"`
		Prizegrades []struct {
			Type      int    `json:"type"`
			Typenum   string `json:"typenum"`
			Typemoney string `json:"typemoney"`
		} `json:"prizegrades"`
	} `json:"result"`
}

func getSSQHistory() (SSQHistoryResult, error) {
	resp, err := resty.New().R().Get(site.QuerySite)
	if err != nil {
		logger.Error(err)
		return SSQHistoryResult{}, err
	}
	var result SSQHistoryResult
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		logger.Error(err)
		return SSQHistoryResult{}, err
	}
	return result, nil
}

// 抓取官网
func catGovSite() (doc *goquery.Document) {
	resp, err := resty.New().R().Get(site.GovSite)
	if err != nil {
		logger.Error(err)
		return
	}
	doc, err = goquery.NewDocumentFromReader(strings.NewReader(resp.String()))
	if err != nil {
		logger.Error(err)
		return
	}
	return
}
