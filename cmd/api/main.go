package main

import (
	_ "REST-API-GO-GIN/docs"
	"REST-API-GO-GIN/internal/database"
	"REST-API-GO-GIN/internal/env"
	"database/sql"
	_ "github.com/joho/godotenv/autoload"
	_ "modernc.org/sqlite" // SQLite driver

	"log"
)

type application struct {
	port      int
	jwtSecret string
	models    database.Models
}

// @title Go Gin Rest API
// @version 1.0
// @description A rest API in Go using Gin framework.
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter your bearer token in the format **Bearer &lt;token&gt;**

// Apply the security definition to your endpoints
// @security BearerAuth

func main() {
	db, err := sql.Open("sqlite", "./data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	models := database.NewModels(db)
	app := &application{
		port:      env.GetEnvInt("PORT", 8080),
		models:    models,
		jwtSecret: env.GetEnvString("JWT_SECRET", "some_secret_12345"),
	}
	if err := app.serve(); err != nil {
		log.Fatal(err)
	}
}
