package models

type User struct {
	Nickname string `db:"nickname"`
	Email    string `db:"email"`
	ID       string `db:"id"`
	Password []byte `db:"password"`
}
