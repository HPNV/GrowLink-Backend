package developer

const (
	InsertDeveloperQuery = `
		INSERT INTO developers (uuid, name, email)
		VALUES ($1, $2, $3)
	`
)
