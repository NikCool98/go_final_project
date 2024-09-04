package stor

import (
	"database/sql"
	"fmt"
	"github.com/NikCool98/go_final_project/config"
	"github.com/NikCool98/go_final_project/taskrepeater"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Stor struct {
	db *sql.DB
}

func OpenDb() *sql.DB {
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
	DB, err := sql.Open("sqlite", "scheduler.db")
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
	return DB
}
func NewStore(db *sql.DB) Stor {
	return Stor{db: db}
}

// Добавление задачи
func (s *Stor) CreateTask(t config.Task) (string, error) {
	var err error

	if t.Title == "" {
		return "", fmt.Errorf(`{"error":"Не указан заголовок задачи"}`)
	}

	// Проверяем наличие даты
	if t.Date == "" {
		t.Date = time.Now().Format(config.DateFormat)
	}

	_, err = time.Parse(config.DateFormat, t.Date)
	if err != nil {
		return "", fmt.Errorf(`{"error":"Некорректный формат даты"}`)
	}
	// Если дата меньше сегодняшней, устанавливаем следующую дату по правилу
	if t.Date < time.Now().Format(config.DateFormat) {
		if t.Repeat != "" {
			nextDate, err := taskrepeater.NextDate(time.Now(), t.Date, t.Repeat)
			if err != nil {
				return "", fmt.Errorf(`{"error":"Некорректное правило повторения"}`)
			}
			t.Date = nextDate
		} else {
			t.Date = time.Now().Format(config.DateFormat)
		}
	}

	// Добавляем задачу в базу
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	result, err := s.db.Exec(query, t.Date, t.Title, t.Comment, t.Repeat)
	if err != nil {
		return "", fmt.Errorf(`{"error":"Не удалось добавить задачу"}`)
	}

	// Возвращаем идентификатор добавленной задачи
	id, err := result.LastInsertId()
	if err != nil {
		return "", fmt.Errorf(`{"error":"Не удалось вернуть id новой задачи"}`)
	}
	return fmt.Sprintf("%d", id), nil
}
