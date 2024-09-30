package models

type Theme struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	SubjectID   int    `json:"subject_id"` // references id from subject table
}
