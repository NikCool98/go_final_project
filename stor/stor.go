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

	// Проверка на дату
	if t.Date == "" {
		t.Date = time.Now().Format(config.DateFormat)
	}

	_, err = time.Parse(config.DateFormat, t.Date)
	if err != nil {
		return "", fmt.Errorf(`{"error":"Некорректный формат даты"}`)
	}
	// Если дата меньше сегодняшней, устанавливаем следующую дату
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

	// Возвращаем шв добавленной задачи
	id, err := result.LastInsertId()
	if err != nil {
		return "", fmt.Errorf(`{"error":"Не удалось вернуть id новой задачи"}`)
	}
	return fmt.Sprintf("%d", id), nil
}

func (s *Stor) GetTasks(search string) ([]config.Task, error) {
	var t config.Task
	var tasks []config.Task
	var rows *sql.Rows
	var err error
	if search == "" {
		rows, err = s.db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?", config.MaxTasks)
	} else if date, error := time.Parse("02.01.2006", search); error == nil {
		query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date LIMIT ?"
		rows, err = s.db.Query(query, date.Format(config.DateFormat), config.MaxTasks)
	} else {
		search = "%%%" + search + "%%%"
		query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?"
		rows, err = s.db.Query(query, search, search, config.MaxTasks)
	}
	if err != nil {
		return []config.Task{}, fmt.Errorf(`{"error":"ошибка запроса"}`)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err = rows.Err(); err != nil {
			return []config.Task{}, fmt.Errorf(`{"error":"Ошибка распознавания данных"}`)
		}
		tasks = append(tasks, t)
	}
	if len(tasks) == 0 {
		tasks = []config.Task{}
	}

	return tasks, nil
}

func (s *Stor) GetTask(id string) (config.Task, error) {
	var t config.Task
	if id == "" {
		return config.Task{}, fmt.Errorf(`{"error":"Не указан id"}`)
	}
	row := s.db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id)
	err := row.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		return config.Task{}, fmt.Errorf(`{"error":"Задача не найдена"}`)
	}
	return t, nil
}

func (s *Stor) UpdateTask(t config.Task) error {
	// Проверка на пустой id
	if t.ID == "" {
		return fmt.Errorf(`{"error":"Не указан id"}`)
	}
	// Проверка на Title
	if t.Title == "" {
		return fmt.Errorf(`{"error":"Не указан заголовок задачи"}`)
	}
	// Проверяем дату
	if t.Date == "" {
		t.Date = time.Now().Format(config.DateFormat)
	}

	_, err := time.Parse(config.DateFormat, t.Date)
	if err != nil {
		return fmt.Errorf(`{"error":"Некорректный формат даты"}`)
	}

	// Если дата меньше сегодняшней, устанавливаем следующую дату
	if t.Date < time.Now().Format(config.DateFormat) {
		if t.Repeat != "" {
			nextDate, err := taskrepeater.NextDate(time.Now(), t.Date, t.Repeat)
			if err != nil {

				return fmt.Errorf(`{"error":"Некорректное правило повторения"}`)
			}
			t.Date = nextDate
		} else {
			t.Date = time.Now().Format(config.DateFormat)
		}
	}

	// Обновляем задачу в базе
	query := `UPDATE scheduler SET date=?, title=?, comment=?, repeat=? WHERE id=?`
	result, err := s.db.Exec(query, t.Date, t.Title, t.Comment, t.Repeat, t.ID)
	if err != nil {

		return fmt.Errorf(`{"error":"Задача с таким id не найдена"}`) // вот эта вот история не работает, посчитаем измененные ряды
	}

	//измененные строки
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf(`{"error":"Не удалось посчитать измененные строки"}`)
	}

	if rowsAffected == 0 {
		return fmt.Errorf(`{"error":"Задача с таким id не найдена"}`)
	}
	return nil
}

func (s *Stor) DeleteTask(id string) error {
	// Проерка на пустой id
	if id == "" {
		return fmt.Errorf(`{"error":"Не указан id"}`)
	}
	query := "DELETE FROM scheduler WHERE id = ?"
	result, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf(`{"error":"Не удалось удалить задачу"}`)
	}
	//измененные строки
	rowsAffected, err := result.RowsAffected()
	if err != nil {

		return fmt.Errorf(`{"error":"Не удалось посчитать измененные строки"}`)
	}

	if rowsAffected == 0 {

		return fmt.Errorf(`{"error":"Задача с таким id не найдена"}`)
	}
	return nil
}

func (s *Stor) TaskDone(id string) error {
	var t config.Task

	t, err := s.GetTask(id)
	if err != nil {
		return err
	}
	if t.Repeat == "" {

		err := s.DeleteTask(id)
		if err != nil {
			return err
		}

	} else {
		next, err := taskrepeater.NextDate(time.Now(), t.Date, t.Repeat)
		if err != nil {
			return err
		}
		t.Date = next
		err = s.UpdateTask(t)
		if err != nil {
			return err
		}
	}

	return nil
}
