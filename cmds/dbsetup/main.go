package main

import (
	"database/sql"
	"log"

	"github.com/flyfy1/diarier/orm"
	_ "github.com/mattn/go-sqlite3"
)

func InitDB() *sql.DB {
	db := orm.GetDB()

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
			username TEXT,
			name TEXT,
			password_hash TEXT,
			running_task_id INTEGER,
			session_id INTEGER NOT NULL DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

			FOREIGN KEY (running_task_id) REFERENCES tasks (id)
		);`, []string{
			"CREATE INDEX idx_users_email ON users(email);",
			"CREATE INDEX idx_users_username ON users(username);",
			"CREATE INDEX idx_users_created_at ON users(created_at);",
			"CREATE INDEX idx_users_updated_at ON users(updated_at);",
		},
		},
		{"create tasks table", `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		name TEXT NOT NULL,
		description TEXT,
		status INTEGER NOT NULL DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		completed_at TIMESTAMP,
		due_date DATE,

		FOREIGN KEY (user_id) REFERENCES users (id)
	);`, []string{
			"CREATE INDEX idx_tasks_user_id ON tasks(user_id);",
			"CREATE INDEX idx_tasks_due_date ON tasks(due_date);",
			"CREATE INDEX idx_tasks_time_created ON tasks(created_at);",
			"CREATE INDEX idx_tasks_time_updated ON tasks(updated_at);",
			"CREATE INDEX idx_tasks_time_completed ON tasks(completed_at);",
		}}, {"events table", `
		CREATE TABLE IF NOT EXISTS events (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			task_id INTEGER NOT NULL,
			timing TEXT,
			full_pomo BOOLEAN,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

			FOREIGN KEY (task_id) REFERENCES tasks (id)
		);
		`, []string{
			"CREATE INDEX idx_events_task_id ON events(task_id);",
			"CREATE INDEX idx_events_time_created ON events(created_at);",
			"CREATE INDEX idx_events_time_updated ON events(updated_at);",
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
	InitDB()
}
