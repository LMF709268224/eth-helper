package test

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	mdb "eth-helper/db"

	log "github.com/sirupsen/logrus"
)

func openDB() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True", "gouse", "123456", "test")
	// dsn := "gouse:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}

func AddData() {
	db := openDB()
	// users := []mdb.EthTransferTb{{
	// 	MTo:         "454656465",
	// 	MFrom:       "515151515",
	// 	Txhash:      "jinzhu6",
	// 	Value:       "2",
	// 	Blocknumber: 65,
	// }, {
	// 	MTo:         "454656465",
	// 	MFrom:       "515151515",
	// 	Txhash:      "jinzhu1",
	// 	Value:       "59659",
	// 	Blocknumber: 69,
	// }, {
	// 	MTo:         "454656465",
	// 	MFrom:       "515151515",
	// 	Txhash:      "jinzhu3",
	// 	Value:       "555",
	// 	Blocknumber: 56,
	// }, {
	// 	MTo:         "454656465",
	// 	MFrom:       "515151515",
	// 	Txhash:      "jinzhu7",
	// 	Value:       "12",
	// 	Blocknumber: 5161,
	// }, {
	// 	MTo:         "454656465",
	// 	MFrom:       "515151515",
	// 	Txhash:      "jinzhu8",
	// 	Value:       "12",
	// 	Blocknumber: 56,
	// }}
	users := []mdb.EthAddressTb{{
		Address: "0xddfa8C217a0Fb51a6319e2D863743807d07C9e81",
		AddTime: time.Now(),
		Balance: "2",
	}, {
		Address: "0x6A859c1BD1D0722D7Bc3Df05a96FF7684dCA30eD",
		AddTime: time.Now(),
		Balance: "6",
	}}
	tx := db.Create(&users)
	if tx.Error != nil {
		log.Errorf("Error :%v", tx.Error)
	}

	// for _, user := range users {
	// 	log.Infoln("user ID: ", user.ID)
	// 	log.Infoln("user Address: ", user.Address)
	// }
}

func AddData1() {
	db := openDB()

	users := mdb.EthAddressTb{
		Address: "jinzhu5",
		AddTime: time.Now(),
		Balance: "123",
	}
	tx := db.Create(&users)
	if tx.Error != nil {
		log.Errorf("Error :%v", tx.Error)
	}

	// for _, user := range users {
	// 	log.Infoln("user ID: ", user.ID)
	// 	log.Infoln("user Address: ", user.Address)
	// }
}

// GetAddressInfo 获取地址信息
func GetAddressInfo() {
	db := openDB()

	var rows mdb.EthAddressTb
	// select
	ss := db.Where("address = ?", "jinzhu5").Find(&rows)
	if ss.Error != nil {
		fmt.Println("err 不是空")
	}

	fmt.Printf("rows id:%d,address:%s,err:%v\n", rows.ID, rows.Address, ss.Error)

	return
}

// GetAddressInfos 获取地址信息
func GetAddressInfos() {
	db := openDB()

	var rows []mdb.EthAddressTb
	// select
	ss := db.Find(&rows)
	fmt.Printf("err : %v", ss.Error)

	for _, row := range rows {
		fmt.Printf("rows id:%d,address:%s,err:%v\n", row.ID, row.Address, ss.Error)
	}
	return
}

func Dele() {
	db := openDB()
	var row []mdb.EthAddressTb

	ss := db.Where("address = ?", "jinzhu2").Delete(&row)
	fmt.Printf("err : %v", ss.Error)
	fmt.Printf("row : %v", row)
}

func getMin() {
	db := openDB()

	ss := fmt.Sprintf("select MIN(blocknumber) from  %s", "eth_transfer_tbs")
	// 	// ss := fmt.Sprintf("select blocknumber from %s limit 1", transferTable)
	row := db.Raw(ss)
	blocknumber := 0
	err := row.Scan(&blocknumber)
	if err.Error != nil {
		fmt.Printf("GetMinBlocknumber Scan err : %v", err.Error)
	}
	fmt.Printf("blocknumber : %v", blocknumber)
}

func TestG() {
	AddData()
}
