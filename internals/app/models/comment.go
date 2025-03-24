package models

import "time"

type Comment struct {
	ID      int    `json:"id" db:"id"`
	UserID  string `json:"user_id" db:"user_id"`
	Email   string `json:"email" db:"email"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	ThemeID int    `json:"theme_id" db:"theme_id" required`
	Content string `json:"content" db:"content" required`
}

type AdminReply struct {
	ID      int    `json:"id"`
	CommentID int    `json:"comment_id" required`
	UserID string `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	Content string `json:"content" required`
}

type CommentAdminReply struct {
	ID int `json:"id"`
	CommentID int `json:"comment_id"`
	Content string `json:"content"`
	UserID string `json:"user_id"`
	Email string `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
