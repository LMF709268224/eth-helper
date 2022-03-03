package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// MAddressInfo 地址
type MAddressInfo struct {
	ID         int64     `json:"id"`
	Address    string    `json:"address"`    // 地址
	AddTime    time.Time `json:"addtime"`    // 创建时间
	PType      string    `json:"ptype"`      // 类型
	PrivateKey []byte    `json:"privatekey"` // 私钥
	KeyType    string    `json:"keytype"`    // 私钥类型
	Msg        string    `json:"msg"`        // 备注信息
	Balance    int64     `json:"balance"`    // 余额
}

var (
	sqlUseName   = ""
	sqlPassword  = ""
	sqlDatabase  = ""
	addressTable = "addressTable"
)

// var addressTableMap = map[structconst.BlockchainType]string{structconst.BlockchainTypeFilecoin: addressFilecoinTable}

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