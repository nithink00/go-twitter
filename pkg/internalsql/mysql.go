package internalsql

import (
	"database/sql"
	"fmt"
	"go-twitter/internal/config"

	_ "github.com/go-sql-driver/mysql"
)


func ConnectMySQL(cfg *config.Config) (*sql.DB, error) {
	dataSourcName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=%s", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, "Local")
	db, err := sql.Open("mysql", dataSourcName)
	if err != nil {
		return nil, err
	}

	fmt.Println("Database loaded successfully")
	return db, nil
}