package db

type Business struct {
	UUID        string `db:"uuid"`
	UserUUID    string `db:"user_uuid"`
	CompanyName string `db:"company_name"`
}
