package domain

type User struct {
	Id       string `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
	Score    string `json:"score" db:"score"`
}
