package db

import (
	"database/sql"
	"embed"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"
)

type DataBase struct {
	DB  *sql.DB
	DSN string
}

func Open(dsn string) (*DataBase, error) {
	db, err := sql.Open("postgres", dsn)

	if err != nil {
		return nil, err
	}
	return &DataBase{DB: db, DSN: dsn}, nil
}

//go:embed migrations/*
var migrationsFS embed.FS

func (db *DataBase) RunMigrations() error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("не удалось создать драйвер: %w", err)
	}
	source, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("не удалось создать источник миграций: %w", err)
	}

	migrations, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
	if err != nil {
		return fmt.Errorf("не удалось инициализировать миграции: %w", err)
	}

	err = migrations.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			log.Println("Новых миграций не найдено")
			return nil
		}
		return fmt.Errorf("ошибка при выполнении миграций: %w", err)
	}

	log.Println("Миграции успешно применены")
	return nil
}
