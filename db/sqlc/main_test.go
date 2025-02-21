package sqlc

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/hoangphuc3064/MyBank/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB
var testStore *Store

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("Không thể load config:", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Không thể kết nối database:", err)
	}

	testQueries = New(testDB)
	testStore = NewStore(testDB)

	os.Exit(m.Run())
}
