package student

const (
	CreateQuery = `
		INSERT INTO students (user_uuid, university)
		VALUES ($1, $2)
		RETURNING uuid
	`

	GetByUUIDQuery = `SELECT uuid, user_uuid, university FROM students WHERE uuid = $1`

	GetByUserUUIDQuery = `SELECT uuid, user_uuid, university FROM students WHERE user_uuid = $1`

	UpdateQuery = `
		UPDATE students 
		SET university = $1
		WHERE uuid = $2
	`

	DeleteQuery = `DELETE FROM students WHERE uuid = $1`

	GetAllQuery = `SELECT uuid, user_uuid, university FROM students ORDER BY university`

	AddSkillQuery = `
		INSERT INTO student_skills (student_uuid, skill_uuid)
		VALUES ($1, $2)
		ON CONFLICT (student_uuid, skill_uuid) DO NOTHING
	`

	RemoveSkillQuery = `DELETE FROM student_skills WHERE student_uuid = $1 AND skill_uuid = $2`

	GetSkillsQuery = `
		SELECT s.uuid, s.name, s.description, s.created_at
		FROM skills s
		INNER JOIN student_skills ss ON s.uuid = ss.skill_uuid
		WHERE ss.student_uuid = $1
		ORDER BY s.name
	`
)
