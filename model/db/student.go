package db

type Student struct {
	UUID       string `db:"uuid"`
	UserUUID   string `db:"user_uuid"`
	University string `db:"university"`
}
