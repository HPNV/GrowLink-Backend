package user

const (
	getUserByEmailQuery = `
		SELECT 
			uuid, 
			email, 
			name,
			password_hash, 
			role, 
			created_at
		FROM users 
		WHERE email=$1
	`
	createUserQuery = `
		INSERT INTO users (email, name, password_hash, role) 
		VALUES ($1, $2, $3, $4) 
		RETURNING uuid, email, name, password_hash, role, created_at
	`
	getAllUsersQuery = `
		SELECT 
			uuid, 
			email, 
			name,
			password_hash, 
			role, 
			created_at
		FROM users 
		ORDER BY created_at DESC
	`

	getUserByUUIDQuery = `
		SELECT 
			u.uuid,
			u.email,
			u.name,
			u.role,
			u.created_at
		FROM users u
		WHERE u.uuid = $1
	`

	getStudentListUsersQuery = `
		SELECT 
			u.uuid,
			u.email,
			u.name,
			u.role,
			u.created_at
		FROM users u
		INNER JOIN students s ON u.uuid = s.user_uuid
		WHERE u.role = 'student'
	`

	getStudentListStudentsQuery = `
		SELECT 
			s.uuid,
			s.user_uuid,
			s.university
		FROM students s
		INNER JOIN users u ON s.user_uuid = u.uuid
		WHERE u.role = 'student'
	`

	getStudentListCountQuery = `
		SELECT COUNT(*)
		FROM users u
		INNER JOIN students s ON u.uuid = s.user_uuid
		WHERE u.role = 'student'
	`
)
