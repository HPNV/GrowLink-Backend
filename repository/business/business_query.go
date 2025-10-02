package business

const (
	CreateQuery = `
		INSERT INTO businesses (user_uuid, company_name)
		VALUES ($1, $2)
		RETURNING uuid
	`

	GetByUUIDQuery = `SELECT uuid, user_uuid, company_name FROM businesses WHERE uuid = $1`

	GetByUserUUIDQuery = `SELECT uuid, user_uuid, company_name FROM businesses WHERE user_uuid = $1`

	UpdateQuery = `
		UPDATE businesses 
		SET company_name = $1
		WHERE uuid = $2
	`

	DeleteQuery = `DELETE FROM businesses WHERE uuid = $1`

	GetAllQuery = `SELECT uuid, user_uuid, company_name FROM businesses ORDER BY company_name`
)
