package main

import (
	"MyInterpretation/models"
	"MyInterpretation/views"
	"context"
	"log"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

		var page models.PageData

		page.Definitions = definitions

		log.Print(page)

		return cntx.Render(200, "index", page)
	})

	app.POST("/definition", func(cntx echo.Context) error {
		text := cntx.FormValue("text")
		wordId, err := strconv.Atoi(cntx.FormValue("wordId"))

		if text == "" || err != nil {
			return cntx.Render(400, "footer", 1)
		}

		if err := dbPool.QueryRow(context.Background(), "INSERT INTO Definition (text, wordId) Values ($1, $2)", text, wordId); err != nil {
			return cntx.Render(404, "footer", 1)
		}

		return cntx.Render(200, "footer", 1)
	})

	app.Logger.Fatal(app.Start(":8080"))
}
