package sqlc

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	dsn := "root:root@tcp(127.0.0.1:3306)/my_bank?charset=utf8mb4&parseTime=True&loc=Local"
	var err error

	testDB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Không thể kết nối database:", err)
	}

	if err = testDB.Ping(); err != nil {
		log.Fatal("Không thể ping database:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
