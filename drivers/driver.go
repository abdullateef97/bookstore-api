package drivers

import (
	"database/sql"
	"os"
	"github.com/lib/pq"
	"go-api/helpers"
)

var db *sql.DB

//ConnectDB connects to mysql db
func ConnectDB() *sql.DB {
	_, err := pq.ParseURL(os.Getenv("DB_URL"))
	helpers.LogFatal(err)

	db, err = sql.Open("postgres", os.Getenv("DB_URL"))
	helpers.LogFatal(err)

	err = db.Ping()
	helpers.LogFatal(err)

	return db
}