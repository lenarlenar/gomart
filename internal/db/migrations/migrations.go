package migrations

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	
	_ "github.com/lib/pq"
)


var migrations = []string {
	`CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		login TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		created_at TIMESTAMPTZ DEFAULT NOW()
	);`,
}


func Apply(db *sql.DB) {
	for _, query := range migrations {
		_, err := db.ExecContext(context.Background(), query)
		if err != nil {
			log.Fatalf("Ошибка применения миграции: %v", err)
		}
	}

	fmt.Println("Все миграции успешно применены")
}

