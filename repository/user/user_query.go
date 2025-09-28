package user

const (
	getUserByEmailQuery = `
		SELECT 
			uuid, 
			email, 
			password_hash, 
			role, 
			created_at
		FROM users 
		WHERE email=$1
	`
	createUserQuery = `
		INSERT INTO users (email, password_hash, role) 
		VALUES ($1, $2, $3) 
		RETURNING uuid, email, password_hash, role, created_at
	`
	getAllUsersQuery = `
		SELECT 
			uuid, 
			email, 
			password_hash, 
			role, 
			created_at
		FROM users 
		ORDER BY created_at DESC
	`
)
