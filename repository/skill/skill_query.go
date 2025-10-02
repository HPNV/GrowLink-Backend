package skill

const (
	CreateSkillQuery = `INSERT INTO skills (name) VALUES ($1) RETURNING uuid`

	CreateQuery = `
		INSERT INTO skills (name, description)
		VALUES ($1, $2)
		RETURNING uuid, created_at
	`

	GetByUUIDQuery = `SELECT uuid, name, description, created_at FROM skills WHERE uuid = $1`

	GetByNameQuery = `SELECT uuid, name, description, created_at FROM skills WHERE name = $1`

	UpdateQuery = `
		UPDATE skills 
		SET name = $1, description = $2
		WHERE uuid = $3
	`

	DeleteQuery = `DELETE FROM skills WHERE uuid = $1`

	GetAllQuery = `SELECT uuid, name, description, created_at FROM skills ORDER BY name`

	GetByProjectUUIDQuery = `
		SELECT s.uuid, s.name, s.description, s.created_at
		FROM skills s
		JOIN project_skills ps ON s.uuid = ps.skill_uuid
		WHERE ps.project_uuid = $1
	`

	GetByStudentUUIDQuery = `
		SELECT s.uuid, s.name, s.description, s.created_at
		FROM skills s
		JOIN student_skills ss ON s.uuid = ss.skill_uuid
		WHERE ss.student_uuid = $1
	`
)
