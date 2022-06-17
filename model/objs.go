package model

import (
	"gorm.io/gorm"
	"time"
)

type PrizeInformation struct {
	Num         int       `json:"num"`           // 期数
	RedNum      []string  `json:"red_num"`       // 红球
	BlueNum     string    `json:"blue_num"`      // 篮球
	Prize1Item  int       `json:"prize_1_item"`  // 1等奖注数
	Prize1Money string    `json:"prize_1_money"` // 1等奖中奖金额
	Prize2Item  int       `json:"prize_2_item"`  // 2等奖注数
	Prize2Money string    `json:"prize_2_money"` // 2等奖中奖金额
	MoneyPool   string    `json:"money_pool"`    // 奖池
	TotalSales  string    `json:"total_sales"`   // 销售额
	PrizeTime   time.Time `json:"prize_time"`    // 开奖日期 2006-01-02
}

type Prize struct {
	Num     string   `json:"num"`
	RedNum  []string `json:"red_num"`  // 红球
	BlueNum string   `json:"blue_num"` // 篮球
}

type GenNum struct {
	RedNum  []string `json:"red_num"`  // 红球
	BlueNum string   `json:"blue_num"` // 篮球
}

type NumSaveData struct {
	gorm.Model
	RdbID      string `json:"rdb_id" gorm:"unique;comment:'redis id'"`
	Uid        string `json:"uid" gorm:"size:50"`        //用户uid
	Num        int    `json:"num"`                       // 第几期
	List       string `json:"list" gorm:"type:longtext"` // 选择的数字列表
	BingoInfo  string `json:"bingo_level"`               // 中奖情况
	BingoMoney uint   `json:"bingo_money"`               // 中奖金额 元
}

type CheckNum struct {
	gorm.Model

	Num       int  `json:"num" gorm:"index;comment:'期号'"` // 期号
	NeedCheck bool `json:"need_check" gorm:"comment:'是否需要检查'"`
	Checked   bool `json:"checked" gorm:"comment:'是否已经检查'"`
}

// TmpSaveNum redis 临时存储
type TmpSaveNum struct {
	ID   string   `json:"id"`   // 随机id redis key
	UID  string   `json:"uid"`  // 用户id
	Num  int      `json:"num"`  // 期号
	List []GenNum `json:"list"` // 号码
}

type BingoInfo struct {
	BingoInfo  string `json:"bingo_level"` // 中奖情况
	BingoMoney uint   `json:"bingo_money"` // 中奖金额 元
}
