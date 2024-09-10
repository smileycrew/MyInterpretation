package main

import (
	"MyInterpretation/handlers"
	"MyInterpretation/models"
	"MyInterpretation/views"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"os"
	"strconv"
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

		var page models.PageData

		page.Definitions = handlers.GetDefinitionsWithWord(cntx, dbPool)

		var wordId int

		for _, d := range page.Definitions {
			if d.Word.ID != wordId {
				wordId = d.Word.ID
				page.Words = append(page.Words, d.Word)
			}
		}

		return cntx.Render(200, "index", page)
	})

	app.DELETE("/definition/:id", func(cntx echo.Context) error {
		id, err := strconv.Atoi(cntx.Param("id"))

		if err != nil {
			log.Fatalf("Bad request: %v", err)
		}

		handlers.DeleteDefinitionById(cntx, dbPool, id)

		return cntx.NoContent(200)
	})

	app.POST("/definition", func(cntx echo.Context) error {
		text := cntx.FormValue("text")
		wordId, err := strconv.Atoi(cntx.FormValue("wordId"))

		if text == "" || err != nil {
			return cntx.Render(400, "footer", 1)
		}

		var newDefinition models.Definition

		if err := dbPool.QueryRow(context.Background(), "INSERT INTO Definition (text, wordId) Values ($1, $2) RETURNING id, text, wordId", text, wordId).Scan(&newDefinition.ID, &newDefinition.Text, &newDefinition.WordId); err != nil {
			return cntx.Render(404, "footer", 1)
		}

		newDefinition.Word = handlers.GetWordById(cntx, dbPool, newDefinition.WordId)

		// definitions := handlers.GetDefinitionsWithWord(cntx, dbPool)

		log.Print(newDefinition)
		return cntx.Render(200, "definition-list-item", newDefinition)
	})

	app.Logger.Fatal(app.Start(":8080"))
}
