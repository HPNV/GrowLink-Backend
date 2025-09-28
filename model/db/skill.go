package db

type Skill struct {
	UUID        string `db:"uuid"`
	Name        string `db:"name"`
	Description string `db:"description"`
	CreatedAt   string `db:"created_at"`
}
