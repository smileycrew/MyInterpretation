package handlers

import (
	"MyInterpretation/models"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"log"
)

func GetWordById(cntx echo.Context, dbPool *pgxpool.Pool, wordId int) models.Word {
	rows, err := dbPool.Query(context.Background(), "SELECT * FROM Word WHERE id = $1", wordId)

	if err != nil {
		log.Fatalf("Error querying the database: %v", err)
	}
	defer rows.Close()

	var word models.Word

	for rows.Next() {
		if err := rows.Scan(&word.ID, &word.Text); err != nil {
			log.Fatalf("Error scanning the rows: %v", err)
		}
	}

	return word
}
