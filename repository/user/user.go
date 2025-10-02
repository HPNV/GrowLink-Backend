package user

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	constant "github.com/HPNV/growlink-backend/constant"
	modelDB "github.com/HPNV/growlink-backend/model/db"
	"github.com/HPNV/growlink-backend/model/dto"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type IUser interface {
	Login(ctx context.Context, email, password string) (*modelDB.User, error)
	Register(ctx context.Context, tx *sqlx.Tx, user *modelDB.User, plainPassword string) (*modelDB.User, error)
	GetAll() ([]*modelDB.User, error)
	GetByUUID(ctx context.Context, uuid string) (*modelDB.User, error)
	GetStudentList(queryParam *dto.StudentListRequest) ([]*modelDB.User, []*modelDB.Student, int, error)
}

type User struct {
	db *sqlx.DB
}

func NewUser(db *sqlx.DB) IUser {
	return &User{
		db: db,
	}
}

func (u *User) Login(ctx context.Context, email, plainPassword string) (*modelDB.User, error) {
	var user modelDB.User

	err := u.db.QueryRowContext(ctx, getUserByEmailQuery, email).
		Scan(&user.UUID, &user.Email, &user.Name, &user.PasswordHash, &user.Role, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, constant.ErrUserNotFound
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(plainPassword)); err != nil {
		return nil, constant.ErrInvalidCredentials
	}

	user.PasswordHash = ""

	return &user, nil
}

func (u *User) Register(ctx context.Context, tx *sqlx.Tx, user *modelDB.User, plainPassword string) (*modelDB.User, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = string(hashedBytes)

	err = tx.QueryRowContext(ctx, createUserQuery,
		user.Email, user.Name, user.PasswordHash, user.Role).Scan(&user.UUID, &user.Email, &user.Name, &user.PasswordHash, &user.Role, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	user.PasswordHash = ""

	return user, nil
}

func (u *User) GetAll() ([]*modelDB.User, error) {
	var users []*modelDB.User

	rows, err := u.db.Query(getAllUsersQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user modelDB.User
		err := rows.Scan(&user.UUID, &user.Email, &user.Name, &user.PasswordHash, &user.Role, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		// Don't return password hash for security
		user.PasswordHash = ""
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (u *User) GetByUUID(ctx context.Context, uuid string) (*modelDB.User, error) {
	var user modelDB.User
	err := u.db.QueryRowContext(ctx, getUserByUUIDQuery, uuid).Scan(&user.UUID, &user.Email, &user.Name, &user.Role, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, constant.ErrUserNotFound
		}
		return nil, err
	}
	user.PasswordHash = ""
	return &user, nil
}

func (u *User) GetStudentList(queryParam *dto.StudentListRequest) ([]*modelDB.User, []*modelDB.Student, int, error) {
	var users []*modelDB.User
	var students []*modelDB.Student
	var args []interface{}
	var whereConditions []string
	argIndex := 1

	// Build WHERE conditions
	if queryParam.Name != nil && *queryParam.Name != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("u.email ILIKE $%d", argIndex))
		args = append(args, "%"+*queryParam.Name+"%")
		argIndex++
	}

	if queryParam.University != nil && *queryParam.University != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("s.university ILIKE $%d", argIndex))
		args = append(args, "%"+*queryParam.University+"%")
		argIndex++
	}

	if queryParam.Skill != nil && *queryParam.Skill != "" {
		whereConditions = append(whereConditions, fmt.Sprintf(`EXISTS (
			SELECT 1 FROM student_skills ss 
			JOIN skills sk ON ss.skill_uuid = sk.uuid 
			WHERE ss.student_uuid = s.uuid AND sk.name ILIKE $%d
		)`, argIndex))
		args = append(args, "%"+*queryParam.Skill+"%")
		argIndex++
	}

	// Build the complete queries
	userQuery := getStudentListUsersQuery
	studentQuery := getStudentListStudentsQuery
	countQuery := getStudentListCountQuery

	if len(whereConditions) > 0 {
		whereClause := " AND " + strings.Join(whereConditions, " AND ")
		userQuery += whereClause
		studentQuery += whereClause
		countQuery += whereClause
	}

	// Get total count
	var totalCount int
	err := u.db.Get(&totalCount, countQuery, args...)
	if err != nil {
		return nil, nil, 0, err
	}

	// Add ORDER BY and pagination
	userQuery += " ORDER BY u.email"
	studentQuery += " ORDER BY u.email"

	if queryParam.Limit > 0 {
		userQuery += fmt.Sprintf(" LIMIT $%d", argIndex)
		studentQuery += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, queryParam.Limit)
		argIndex++
	}

	if queryParam.Page > 0 {
		offset := (queryParam.Page - 1) * queryParam.Limit
		userQuery += fmt.Sprintf(" OFFSET $%d", argIndex)
		studentQuery += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, offset)
	}

	// Execute the queries
	err = u.db.Select(&users, userQuery, args...)
	if err != nil {
		return nil, nil, 0, err
	}

	err = u.db.Select(&students, studentQuery, args...)
	if err != nil {
		return nil, nil, 0, err
	}

	return users, students, totalCount, nil
}
