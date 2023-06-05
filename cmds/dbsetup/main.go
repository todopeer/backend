package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		log.Fatal(err)
	}

	if db == nil {
		log.Fatal("db nil")
	}

	createTables(db)

	return db
}

func createTables(db *sql.DB) {
	queries := []struct {
		name       string
		query      string
		extensions []string
	}{
		{"create users table", `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			email TEXT NOT NULL UNIQUE,
			name TEXT,
			password_hash TEXT,
			running_task_id INTEGER,

			FOREIGN KEY (running_task_id) REFERENCES tasks (id)
		);`, []string{"CREATE INDEX idx_users_email ON users(email);"},
		},
		{"create tasks table", `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		name TEXT NOT NULL,
		description TEXT,
		status INTEGER NOT NULL DEFAULT 0,
		time_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		time_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		time_completed TIMESTAMP,
		due_date DATE,

		FOREIGN KEY (user_id) REFERENCES users (id)
	);`, []string{
			"CREATE INDEX idx_tasks_user_id ON tasks(user_id);",
			"CREATE INDEX idx_tasks_due_date ON tasks(due_date);",
			"CREATE INDEX idx_tasks_time_created ON tasks(time_created);",
			"CREATE INDEX idx_tasks_time_updated ON tasks(time_updated);",
			"CREATE INDEX idx_tasks_time_completed ON tasks(time_completed);",
		}}, {"events table", `
		CREATE TABLE IF NOT EXISTS events (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			task_id INTEGER NOT NULL,
			timing TEXT,
			full_pomo BOOLEAN,
			time_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			time_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

			FOREIGN KEY (task_id) REFERENCES tasks (id)
		);
		`, []string{
			"CREATE INDEX idx_events_task_id ON events(task_id);",
			"CREATE INDEX idx_events_time_created ON events(time_created);",
			"CREATE INDEX idx_events_time_updated ON events(time_updated);",
		},
		},
	}

	for _, q := range queries {
		log.Println(q.name)
		log.Println(q.query)
		_, err := db.Exec(q.query)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		log.Println("\textensions: ")
		for _, ext := range q.extensions {
			log.Println(ext)
			_, err := db.Exec(ext)
			if err != nil {
				log.Fatalf("Error: %v", err)
			}
		}
	}

	log.Println("Successfully created tables.")
}

func main() {
	InitDB("todo.db")
}
