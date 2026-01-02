package database

import (
	"log"
	"os"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type postgresDatabase struct {
	Db *gorm.DB
}

var (
	once       sync.Once
	dbInstance *postgresDatabase
)

func NewPostgresDatabase() Database {
	once.Do(func() {
		dsn := os.Getenv("DATABASE_URL")
		if dsn == "" {
			log.Fatal("DATABASE_URL não encontrada no ambiente")
		}

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("Erro ao conectar no banco: %v", err)
		}

		dbInstance = &postgresDatabase{Db: db}
	})

	return dbInstance
}

func (p *postgresDatabase) GetDb() *gorm.DB {
	return p.Db
}
