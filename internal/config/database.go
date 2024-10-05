package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	var (
		host        = os.Getenv("MYSQL_HOST")
		user        = os.Getenv("MYSQL_USER")
		password    = os.Getenv("MYSQL_PASSWORD")
		dbname      = os.Getenv("MYSQL_DB")
		port_string = os.Getenv("MYSQL_PORT")
	)
	port, err := strconv.Atoi(port_string)

	if err != nil {
		log.Fatalf("Invalid port number: %v", err)
	}

	fmt.Printf("host: %s, user: %s, password: %s, dbname: %s, port: %d\n", host, user, password, dbname, port)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, dbname)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	// db.AutoMigrate(&model.User{})
	fmt.Println("Database connected successfully")
	return db
}
