package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

const (
	port = "7540"
)

var DB *sql.DB

func main() {
	// Проверка наличия БД
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")
	log.Println(dbFile)
	_, err = os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}
	// если install равен true, после открытия БД требуется выполнить
	// sql-запрос с CREATE TABLE и CREATE INDEX
	DB, err = sql.Open("sqlite", "scheduler.db")
	if err != nil {
		log.Fatal(err)
	}
	// создаем таблицу и индекс
	if install {

		_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date CHAR(8) NOT NULL DEFAULT '',
			title VARCHAR(128) NOT NULL DEFAULT '',
			comment VARCHAR(256) NOT NULL DEFAULT '',
			repeat VARCHAR(128) NOT NULL DEFAULT '')`, `CREATE INDEX date_index on scheduler(date)`)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("БД создана")
	} else {
		log.Println("БД была создана ранее")
	}

	webDir := "./web"
	fileSrv := http.FileServer(http.Dir(webDir))
	http.Handle("/", fileSrv)

	log.Printf("Starting server on port: %s", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Server run error: %v", err)
	}
}
