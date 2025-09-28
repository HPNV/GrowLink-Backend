package user

import (
	"context"
	"database/sql"

	constant "github.com/HPNV/growlink-backend/constant"
	modelDB "github.com/HPNV/growlink-backend/model/db"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type IUser interface {
	Login(ctx context.Context, email, password string) (*modelDB.User, error)
	Register(ctx context.Context, tx *sqlx.Tx, user *modelDB.User, plainPassword string) (*modelDB.User, error)
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
		Scan(&user.UUID, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt)
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
		user.Email, user.PasswordHash, user.Role).Scan(&user.UUID, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	user.PasswordHash = ""

	return user, nil
}
