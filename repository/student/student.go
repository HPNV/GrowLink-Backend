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
	return tx.QueryRow(CreateQuery, student.UserUUID, student.University).Scan(&student.UUID)
}

func (s *Student) GetByUUID(uuid string) (*db.Student, error) {
	student := &db.Student{}
	err := s.db.Get(student, GetByUUIDQuery, uuid)
	return student, err
}

func (s *Student) GetByUserUUID(userUUID string) (*db.Student, error) {
	student := &db.Student{}
	err := s.db.Get(student, GetByUserUUIDQuery, userUUID)
	return student, err
}

func (s *Student) Update(tx *sqlx.Tx, student *db.Student) error {
	_, err := tx.Exec(UpdateQuery, student.University, student.UUID)
	return err
}

func (s *Student) Delete(tx *sqlx.Tx, uuid string) error {
	_, err := tx.Exec(DeleteQuery, uuid)
	return err
}

func (s *Student) GetAll() ([]*db.Student, error) {
	var students []*db.Student
	err := s.db.Select(&students, GetAllQuery)
	return students, err
}

func (s *Student) AddSkill(tx *sqlx.Tx, studentUUID, skillUUID string) error {
	_, err := tx.Exec(AddSkillQuery, studentUUID, skillUUID)
	return err
}

func (s *Student) RemoveSkill(tx *sqlx.Tx, studentUUID, skillUUID string) error {
	_, err := tx.Exec(RemoveSkillQuery, studentUUID, skillUUID)
	return err
}

func (s *Student) GetSkills(studentUUID string) ([]*db.Skill, error) {
	var skills []*db.Skill
	err := s.db.Select(&skills, GetSkillsQuery, studentUUID)
	return skills, err
}
