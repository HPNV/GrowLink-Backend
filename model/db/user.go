package db

type User struct {
	UUID         string `db:"uuid"`
	Email        string `db:"email"`
	Name         string `db:"name"`
	PasswordHash string `db:"password_hash"`
	Role         string `db:"role"`
	CreatedAt    string `db:"created_at"`
}
