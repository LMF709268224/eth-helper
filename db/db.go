package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// MAddressInfo 地址
type MAddressInfo struct {
	ID      int64     `json:"id"`
	Address string    `json:"address"` // 地址
	Pwd     string    `json:"pwd"`     // 加密私钥
	AddTime time.Time `json:"addtime"` // 创建时间

	Balance    int64  `json:"balance"`    // 余额
	PType      string `json:"ptype"`      // 类型
	PrivateKey []byte `json:"privatekey"` // 私钥
	KeyType    string `json:"keytype"`    // 私钥类型
	Msg        string `json:"msg"`        // 备注信息
}

// MTransferInfo 地址
type MTransferInfo struct {
	ID          int64  `json:"id"`
	To          string `json:"to"`          // to
	From        string `json:"from"`        // from
	Txhash      string `json:"txhash"`      // 交易哈希
	Value       int64  `json:"value"`       // 交易值
	Blocknumber int64  `json:"blocknumber"` // blocknumber
}

var (
	sqlUseName    = ""
	sqlPassword   = ""
	sqlDatabase   = ""
	addressTable  = "eth_address_key"
	transferTable = "eth_transfer_key"
)

// InitDB 初始化数据
func InitDB(use string, pass string, database string) {
	sqlUseName = use
	sqlPassword = pass
	sqlDatabase = database
}

func openDB() *sql.DB {
	parm := fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True", sqlUseName, sqlPassword, sqlDatabase)

	db, err := sql.Open("mysql", parm)
	if err != nil {
		panic(err)
	}

	// See "Important settings" section.
	// db.SetConnMaxLifetime(time.Minute * 3)
	// db.SetMaxOpenConns(10)
	// db.SetMaxIdleConns(10)

	return db
}
