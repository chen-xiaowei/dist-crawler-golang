package main

import (
	"database/sql"
	"fmt"
	_ "mysql"
	"strings"
	"time"
)

var DB *sql.DB = initDB()

const (
	userName = "root"
	password = "123456"
	ip       = "127.0.0.1"
	port     = "3306"
	dbName   = "crawler"
	params   = "?charset=utf8"
)

func initDB() *sql.DB {
	path := strings.Join([]string{userName, ":", password, "@tcp(", ip, ":", port, ")/", dbName, params}, "")
	db, _ := sql.Open("mysql", path)
	db.SetConnMaxLifetime(100)
	db.SetMaxIdleConns(10)
	if err := db.Ping(); err != nil {
		fmt.Println("DB连接异常|" + err.Error())
	}
	return db
}

func init() {
	// sql.Register("mysql", &mysql.MySQLDriver{})
}

func save(house House) bool {
	// tx, err := DB.Begin()
	insert := "insert into house " +
		"(city, region, sub_region, address, build_year, building_count, household_count, acreage, unit_price, publish_time)" +
		" values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	stmt, _ := DB.Prepare(insert)
	_, err := stmt.Exec(house.City, house.Region, house.SubRegion, house.Address, house.BuildYear, house.BuildingCount, house.HouseholdCount, house.Acreage, house.UnitPrice, house.PublishTime)
	if err != nil {
		Err(err.Error() + "|URL:" + house.Link + "|Region:" + house.Region + "|SubRegion:" + house.SubRegion + "|UnitPrice:" + house.UnitPrice)
		return false
	}
	// tx.Commit()
	return true
}

func getMaxIdOfHouse() string {
	var maxId string //  count(*)
	err := DB.QueryRow("select max(id) from house").Scan(&maxId)
	if err != nil {
		// fmt.Println(err.Error())
		return "0"
	}
	return maxId
}

func recordLog(log string) bool {
	// tx, err := DB.Begin()
	insert := "insert into err_log (log, time) values (?, ?)"
	stmt, _ := DB.Prepare(insert)
	_, err := stmt.Exec(log, time.DateTime)
	if err != nil {
		return false
	}
	// tx.Commit()
	return true
}
