package db

// eth_scanningblock_tbs 表

// TODO

// SaveBlockNumber 保存当前扫块高度
func SaveBlockNumber(num int) error {
	// db := GetDBConnection()

	// tx := db.Create(&infos)
	// if tx.Error != nil {
	// 	log.Errorf("SaveNewAddress Error :%v", tx.Error)
	// }

	return nil
}

// GetBlockNumber 获取当前扫快高度
func GetBlockNumber() int {
	// db := GetDBConnection()

	// var row EthAddressTb
	// // select
	// tx := db.Where("address = ?", address).Find(&row)
	// if tx.Error != nil {
	// 	log.Errorf("GetAddressInfo Error :%v", tx.Error)
	// }

	return 0
}
