package main

import (
	"MyInterpretation/models"
	"MyInterpretation/views"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"os"
)

func main() {
	// load env files
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env files: %v", err)
	}

	connStr := os.Getenv("DB_CONNECTION_STRING")

	// create a new connection pool to the database
	dbPool, err := pgxpool.New(context.Background(), connStr)

	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)

		os.Exit(1)
	}

	defer dbPool.Close()

	if err := dbPool.Ping(context.Background()); err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	// initialize echo
	app := echo.New()
	app.Use(middleware.Logger())

	// get your templates and static files
	app.Renderer = views.GetTemplates()
	// app.Static("/styles/", "styles")

	app.GET("/", func(cntx echo.Context) error {
		rows, err := dbPool.Query(context.Background(), "SELECT * FROM Definition JOIN Word ON Word.id = Definition.wordId")

		if err != nil {
			return cntx.Render(500, "index", 1)
		}

		defer rows.Close()

		var definitions []models.Definition

		for rows.Next() {
			var definition models.Definition
			var word models.Word

			if err := rows.Scan(&definition.ID, &definition.Text, &definition.WordId, &word.ID, &word.Text); err != nil {
				return cntx.Render(400, "index", 1)
			}

			definition.Word = word

			definitions = append(definitions, definition)
		}

		return cntx.Render(200, "index", definitions)
	})

	app.Logger.Fatal(app.Start(":8080"))
}
