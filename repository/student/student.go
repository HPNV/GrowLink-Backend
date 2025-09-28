package student

import (
	"github.com/HPNV/growlink-backend/model/db"
	"github.com/jmoiron/sqlx"
)

type IStudent interface {
	Create(tx *sqlx.Tx, student *db.Student) error
	GetByUUID(uuid string) (*db.Student, error)
	GetByUserUUID(userUUID string) (*db.Student, error)
	Update(tx *sqlx.Tx, student *db.Student) error
	Delete(tx *sqlx.Tx, uuid string) error
	GetAll() ([]*db.Student, error)
	AddSkill(tx *sqlx.Tx, studentUUID, skillUUID string) error
	RemoveSkill(tx *sqlx.Tx, studentUUID, skillUUID string) error
	GetSkills(studentUUID string) ([]*db.Skill, error)
}

type Student struct {
	db *sqlx.DB
}

func NewStudent(db *sqlx.DB) IStudent {
	return &Student{
		db: db,
	}
}

func (s *Student) Create(tx *sqlx.Tx, student *db.Student) error {
	query := `
		INSERT INTO students (user_uuid, university)
		VALUES ($1, $2)
		RETURNING uuid
	`
	return tx.QueryRow(query, student.UserUUID, student.University).Scan(&student.UUID)
}

func (s *Student) GetByUUID(uuid string) (*db.Student, error) {
	student := &db.Student{}
	query := `SELECT uuid, user_uuid, university FROM students WHERE uuid = $1`
	err := s.db.Get(student, query, uuid)
	return student, err
}

func (s *Student) GetByUserUUID(userUUID string) (*db.Student, error) {
	student := &db.Student{}
	query := `SELECT uuid, user_uuid, university FROM students WHERE user_uuid = $1`
	err := s.db.Get(student, query, userUUID)
	return student, err
}

func (s *Student) Update(tx *sqlx.Tx, student *db.Student) error {
	query := `
		UPDATE students 
		SET university = $1
		WHERE uuid = $2
	`
	_, err := tx.Exec(query, student.University, student.UUID)
	return err
}

func (s *Student) Delete(tx *sqlx.Tx, uuid string) error {
	query := `DELETE FROM students WHERE uuid = $1`
	_, err := tx.Exec(query, uuid)
	return err
}

func (s *Student) GetAll() ([]*db.Student, error) {
	var students []*db.Student
	query := `SELECT uuid, user_uuid, university FROM students ORDER BY university`
	err := s.db.Select(&students, query)
	return students, err
}

func (s *Student) AddSkill(tx *sqlx.Tx, studentUUID, skillUUID string) error {
	query := `
		INSERT INTO student_skills (student_uuid, skill_uuid)
		VALUES ($1, $2)
		ON CONFLICT (student_uuid, skill_uuid) DO NOTHING
	`
	_, err := tx.Exec(query, studentUUID, skillUUID)
	return err
}

func (s *Student) RemoveSkill(tx *sqlx.Tx, studentUUID, skillUUID string) error {
	query := `DELETE FROM student_skills WHERE student_uuid = $1 AND skill_uuid = $2`
	_, err := tx.Exec(query, studentUUID, skillUUID)
	return err
}

func (s *Student) GetSkills(studentUUID string) ([]*db.Skill, error) {
	var skills []*db.Skill
	query := `
		SELECT s.uuid, s.name, s.description, s.created_at
		FROM skills s
		INNER JOIN student_skills ss ON s.uuid = ss.skill_uuid
		WHERE ss.student_uuid = $1
		ORDER BY s.name
	`
	err := s.db.Select(&skills, query, studentUUID)
	return skills, err
}
