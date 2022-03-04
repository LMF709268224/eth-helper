package db

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

// eth_address_key 表

// SaveNewAddress 保存新地址
func SaveNewAddress(infos []MAddressInfo) error {
	db := openDB()
	// 延迟到函数结束关闭链接
	defer db.Close()

	// TODO pwd 要去掉
	sqlStr := fmt.Sprintf("INSERT INTO %s(address,pwd,addTime) VALUES", addressTable)

	vals := []interface{}{}
	for _, data := range infos {
		sqlStr += "(?,?,Now()),"
		vals = append(vals, data.Address, data.Pwd)
	}

	// trim the last ,
	sqlStr = sqlStr[0 : len(sqlStr)-1]
	// prepare the statement
	log.Infoln("SaveNewAddress sqlStr: ", sqlStr)
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		log.Errorln("SaveNewAddress Prepare err : ", err)
		return err
	}

	// format all vals at once
	_, err = stmt.Exec(vals...)
	if err != nil {
		log.Errorln("SaveNewAddress Exec err : ", err)
	}

	return err
}

// GetAddressInfo 获取地址信息
// func GetAddressInfo(address string) (id int, privateKey []byte, keyType string) {
// 	db := openDB()
// 	// 延迟到函数结束关闭链接
// 	defer db.Close()

// 	// 执行单条查询
// 	skey := "address"
// 	ss := fmt.Sprintf("select id,privateKey,keyType,ptype from %s where %s = ?", addressTable, skey)
// 	rows := db.QueryRow(ss, address)

// 	rows.Scan(&id, &privateKey, &keyType)

// 	return id, privateKey, keyType
// }

// ChangeBalance 修改余额
func ChangeBalance(address string, balance int64) error {
	db := openDB()
	// 延迟到函数结束关闭链接
	defer db.Close()

	//-------修改数据--------
	ss := fmt.Sprintf("update %s set balance=%d, where address = ?", addressTable, balance)
	_, err := db.Exec(ss, address)
	if err != nil {
		log.Errorln("sql ChangeBalance err : ", err)
		// return err
	}

	return err
}

// GetAllAddressInfo 获取所有地址信息
func GetAllAddressInfo() []MAddressInfo {
	db := openDB()
	// 延迟到函数结束关闭链接
	defer db.Close()

	list := []MAddressInfo{}

	// 执行查询
	ss := fmt.Sprintf("select address,balance from %s;", addressTable)
	rows, err := db.Query(ss)
	if err != nil {
		log.Errorln("sql GetAllAddressInfo Query err : ", err)
		return list
	}

	for rows.Next() {
		var address string
		var balance int64

		err = rows.Scan(&address, &balance)
		if err != nil {
			log.Errorln("sql GetAllAddressInfo Scan err : ", err)
			continue
		}

		// fmt.Println(id, ip, time, methods, status, source)
		list = append(list, MAddressInfo{Address: address, Balance: balance})
	}

	return list
}
