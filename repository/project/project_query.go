package project

const (
	CreateQuery = `
		INSERT INTO projects (name, description, duration, timeline, deliverables, created_by)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING uuid, created_at
	`

	GetByUUIDQuery = `SELECT uuid, name, description, status, created_by, created_at FROM projects WHERE uuid = $1`

	UpdateQuery = `
		UPDATE projects 
		SET name = $1, description = $2, status = $3
		WHERE uuid = $4
	`

	DeleteQuery = `DELETE FROM projects WHERE uuid = $1`

	GetAllQuery = `SELECT uuid, name, description, status, created_by, created_at FROM projects ORDER BY created_at DESC`

	GetByBusinessUUIDQuery = `SELECT uuid, name, description, status, created_by, created_at FROM projects WHERE created_by = $1 ORDER BY created_at DESC`

	AddSkillQuery = `
		INSERT INTO project_skills (project_uuid, skill_uuid)
		VALUES ($1, $2)
		ON CONFLICT (project_uuid, skill_uuid) DO NOTHING
	`

	RemoveSkillQuery = `DELETE FROM project_skills WHERE project_uuid = $1 AND skill_uuid = $2`

	GetSkillsQuery = `
		SELECT s.uuid, s.name, s.description, s.created_at
		FROM skills s
		INNER JOIN project_skills ps ON s.uuid = ps.skill_uuid
		WHERE ps.project_uuid = $1
		ORDER BY s.name
	`

	AddStudentQuery = `
		INSERT INTO student_projects (student_uuid, project_uuid)
		VALUES ($1, $2)
		ON CONFLICT (student_uuid, project_uuid) DO NOTHING
	`

	RemoveStudentQuery = `DELETE FROM student_projects WHERE project_uuid = $1 AND student_uuid = $2`

	GetStudentsQuery = `
		SELECT s.uuid, s.user_uuid, s.university
		FROM students s
		INNER JOIN student_projects sp ON s.uuid = sp.student_uuid
		WHERE sp.project_uuid = $1
		ORDER BY s.university
	`

	GetAllListQuery = `
		SELECT p.uuid, p.name, p.description, p.status, p.duration, p.timeline, p.deliverables, p.created_by, p.created_at
		FROM projects p
	`

	GetAllListCountQuery = `
		SELECT COUNT(*)
		FROM projects p
	`
)
