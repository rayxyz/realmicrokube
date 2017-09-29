package db

import (
	"database/sql"
	"fmt"
	"log"

	micro "realmicrokube/micro"

	_ "github.com/go-sql-driver/mysql"
)

var (
	dsn = "%s:%s@tcp(%s:%d)/%s?strict=true&sql_notes=false&parseTime=true&loc=Local&charset=utf8mb4,utf8"
)

type DB struct {
	*sql.DB
}

func NewDB() *DB {
	config := micro.DBConfig{
		Username: "root",
		Password: "root",
		Host:     "www.ray-xyz.com",
		Port:     3306,
		Database: "test",
	}
	dsn = fmt.Sprintf(dsn, config.Username, config.Password, config.Host, config.Port, config.Database)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	return &DB{
		DB: db,
	}
}

func (db *DB) Close() error {
	return db.DB.Close()
}

const queryUserCountSQL = `select count(*) from user where name is not null;`

func (db *DB) QueryUserCount() (int, error) {
	parsedSQL := fmt.Sprintf(queryUserCountSQL)
	rows, err := db.Query(parsedSQL)
	if err != nil {
		return 0, err
	}
	defer close(rows)
	var count int
	if rows.Next() {
		if err = rows.Scan(&count); err != nil {
			return 0, err
		}
		return count, nil
	}
	return 0, nil
}

func close(rows *sql.Rows) {
	if rows != nil {
		rows.Close()
	}
}
