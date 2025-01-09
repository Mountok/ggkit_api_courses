package models

type User struct {
	Id         string `json:"id" db:"id"`
	Email      string `json:"email" db:"email"`
	Password   string `json:"password" db:"password"`
	Role       string `json:"role" db:"role"`
	CreateDate string `json:"create_date" db:"create_date"`
}

type UserCreate struct {
	Id         string `json:"id"`
	Username   string `json:"username" binding:"required"`
	Email      string `json:"email" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Role       string `json:"role"`
	CreateDate string `json:"create_date"`
}
type UserSignIn struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ValidateUser struct {
	Id   string `json:"id"`
	Role string `json:"role"`
}

type RatingUser struct {
	Id       string `json:"id" db:"user_id"`
	Username string `json:"username" db:"full_name"`
	ImageUrl string `json:"image" db:"image"`
	Score    int    `json:"score" db:"score"`
}

type LastSubject struct {
	Id         int    `json:"id" db:"id"`
	UserId     string `json:"user_id" db:"user_id"`
	SubjectsId int    `json:"subjects_id" db:"subjects_id"`
}

type FindUser struct {
	Id string `json:"id" db:"user_id"`
	Image string `json:"image" db:"image"`
	FullName string `json:"full_name" db:"full_name"`
}