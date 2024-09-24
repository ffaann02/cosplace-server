package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/ffaann02/cosplace-server/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	var (
		host        = os.Getenv("POSTGRES_HOST")
		user        = os.Getenv("POSTGRES_USER")
		password    = os.Getenv("POSTGRES_PASSWORD")
		dbname      = os.Getenv("POSTGRES_DB")
		port_string = os.Getenv("POSTGRES_PORT")
	)
	port, err := strconv.Atoi(port_string)

	if err != nil {
		log.Fatalf("Invalid port number: %v", err)
	}

	fmt.Printf("host: %s, user: %s, password: %s, dbname: %s, port: %d\n", host, user, password, dbname, port)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai", host, user, password, dbname, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&model.User{})
	fmt.Println("Database connected successfully")
	return db
}
