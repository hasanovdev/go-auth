package config

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func DBConn() (db *sql.DB, err error) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "Hasanov@0301"
	dbName := "go_auth"

	db, err = sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	return
}
