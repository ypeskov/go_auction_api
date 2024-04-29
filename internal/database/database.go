package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"ypeskov/go_hillel_9/internal/config"
	log "ypeskov/go_hillel_9/internal/log"
)

//type Database struct {
//	Db *sqlx.DB
//}

type Database = *sqlx.DB

func GetDB(cfg *config.Config, log *log.Logger) Database {
	dbConnStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		cfg.DB_USER, cfg.DB_PASS, cfg.DB_HOST, cfg.DB_PORT, cfg.DB_NAME)

	db, err := sqlx.Connect("postgres", dbConnStr)
	if err != nil {
		log.Panicf("Error connecting to the database: %v", err)
		panic(err)
	}

	return db
}
