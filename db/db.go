package db

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var mdb *gorm.DB

// EthAddressTb 地址表
type EthAddressTb struct {
	ID      int64     `gorm:"column:id"`
	Address string    `gorm:"column:address"` // 地址
	Pwd     string    `gorm:"column:pwd"`     // 加密私钥
	AddTime time.Time `gorm:"column:addtime"` // 创建时间
	Balance string    `gorm:"column:balance"` // 余额
	Msg     string    `gorm:"column:msg"`     // 备注信息

	// PType      string `json:"ptype"`      // 类型
	// PrivateKey []byte `json:"privatekey"` // 私钥
	// KeyType    string `json:"keytype"`    // 私钥类型
}

// EthTransferTb 待确认交易表
type EthTransferTb struct {
	ID          int64  `gorm:"column:id"`
	MTo         string `gorm:"column:mto"`         // to
	MFrom       string `gorm:"column:mfrom"`       // from
	Txhash      string `gorm:"column:txhash"`      // 交易哈希
	Value       string `gorm:"column:value"`       // 交易值
	Blocknumber uint64 `gorm:"column:blocknumber"` // blocknumber
}

// EthTransferdoneTb 交易完成表
type EthTransferdoneTb struct {
	ID     int64  `gorm:"column:id"`
	Txhash string `gorm:"column:txhash"` // 交易哈希
	State  int64  `gorm:"column:state"`  // 交易状态
	Msg    string `gorm:"column:msg"`    // 备注信息
}

var (
	sqlUseName        = ""
	sqlPassword       = ""
	sqlDatabase       = ""
	addressTable      = "eth_address_tbs"
	transferTable     = "eth_transfer_tbs"
	transferDoneTable = "eth_transferdone_tbs"
)

// const (
// 	HOST    = "localhost"
// 	PORT    = "3306"
// 	CHARSET = "utf8"
// )

// InitDB 初始化数据
func InitDB(use string, pass string, database string) {
	sqlUseName = use
	sqlPassword = pass
	sqlDatabase = database
}

// func openDB() *sql.DB {
// 	parm := fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True", sqlUseName, sqlPassword, sqlDatabase)
// 	// parm := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True", sqlUseName, sqlPassword, HOST, PORT, sqlDatabase, CHARSET)
// 	// parm := fmt.Sprintf("%s:%s@tcp(localhost)/%s?charset=utf8&parseTime=True&loc=Local", sqlUseName, sqlPassword, sqlDatabase)
// 	// log.Infoln("parm : ", parm)
// 	db, err := sql.Open("mysql", parm)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// See "Important settings" section.
// 	// db.SetConnMaxLifetime(time.Minute * 3)
// 	// db.SetMaxOpenConnsopenDB
// 	return db
// }

// GetDBConnection 获取数据库链接
func GetDBConnection() *gorm.DB {
	if mdb != nil {
		return mdb
	}

	dsn := fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True", sqlUseName, sqlPassword, sqlDatabase)
	mdb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// dbobj, err := mdb.DB()
	// dbobj.SetConnMaxLifetime()

	return mdb
}
