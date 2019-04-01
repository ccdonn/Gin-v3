package config

import (
	"database/sql"
	"log"

	"../constant"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func getDbConnection() {
	var err error
	db, err = sql.Open("mysql", constant.MySQLAccount+":"+constant.MySQLPassword+"@tcp("+constant.MySQLAddress+":"+constant.MySQLPort+")/promotion?parseTime=true")
	if err != nil {
		log.Panic(err)
	} else {
		db.SetMaxIdleConns(5)
		// db.SetMaxOpenConns(3)
		db.SetConnMaxLifetime(60)
	}
}

func GetDBConn() *sql.DB {

	if db == nil {
		getDbConnection()
	}

	return db
}
