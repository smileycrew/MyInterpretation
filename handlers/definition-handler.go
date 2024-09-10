package handlers

import (
	"MyInterpretation/models"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"log"
)

func DeleteDefinitionById(cntx echo.Context, dbPool *pgxpool.Pool, definitionId int) {
	rows, err := dbPool.Exec(context.Background(), "DELETE FROM Definition WHERE id = $1", definitionId)

	if err != nil {
		log.Fatalf("Could not execute delete from database: %v", err)
	}

	log.Print("//////////")
	log.Print(rows)
}

func GetDefinitionsWithWord(cntx echo.Context, dbPool *pgxpool.Pool) []models.Definition {
	rows, err := dbPool.Query(context.Background(), "SELECT * FROM Definition JOIN Word ON Word.id = Definition.wordId")

	if err != nil {
		log.Fatalf("Could not query the database: %v", err)
	}

	defer rows.Close()

	var definitions []models.Definition

	for rows.Next() {
		var definition models.Definition
		var word models.Word

		if err := rows.Scan(&definition.ID, &definition.Text, &definition.WordId, &word.ID, &word.Text); err != nil {
			log.Fatalf("Could not scan rows: %v", err)
		}

		definition.Word = word

		definitions = append(definitions, definition)
	}

	return definitions
}
