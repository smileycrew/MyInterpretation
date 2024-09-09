package models

type Definition struct {
	ID     int    `json:"id"`
	Text   string `json:"text"`
	WordId int    `json:"wordId"`
	Word   Word   `json:"word"`
}
