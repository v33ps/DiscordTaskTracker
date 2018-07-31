package main

import (
	"database/sql"
	"log"
	"time"
)

// Tasks stores running tasks
type Tasks struct {
	ID           string
	username     string
	task         string
	reminder     string
	creationTime time.Time
	accomplished int
}

func initDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		log.Fatal(err)
	}
	if db == nil {
		panic("db nil oh no")
	}
	return db
}

func createTable(db *sql.DB) {
	sqlTable := `
	CREATE TABLE IF NOT EXISTS tasks(
		id TEXT NOT NULL PRIMARY KEY,
		username TEXT,
		task TEXT,
		reminder TEXT NOT NULL,
		creationTime DATETIME,
		accomplished INT
	);
	`
	_, err := db.Exec(sqlTable)
	if err != nil {
		log.Fatal(err)
	}
}

func addTask(db *sql.DB, tasks []Tasks) {
	sqlAddItem := `
	INSERT INTO tasks(
		Id,
		username,
		task,
		reminder,
		creationTime,
		accomplished
	) values(?, ?, ?, ?, CURRENT_TIMESTAMP, ?)
	`
	stmt, err := db.Prepare(sqlAddItem)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	for _, task := range tasks {
		_, err2 := stmt.Exec(task.ID, task.username, task.task, task.reminder, task.accomplished)
		if err2 != nil {
			panic(err2)
		}
	}
}
